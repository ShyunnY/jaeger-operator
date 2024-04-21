package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

const (
	defaultConfigMap    = "jaeger-operator"
	defaultConfigMapKey = "jaeger-operator.yaml"
)

var (
	notReSync      = time.Duration(0)
	informerLogger = logging.DefaultLogger()
)

type CallBacks []func(*Server)

type ServerConfigInformer struct {
	name  string
	mutex sync.RWMutex

	started      bool
	callback     CallBacks
	informer     cache.SharedInformer
	registration cache.ResourceEventHandlerRegistration
}

func New(name string, r rest.Interface, cbs CallBacks) *ServerConfigInformer {

	lw := cache.NewListWatchFromClient(
		r,
		"configmaps",
		"",
		fields.Everything(),
	)
	sc := &ServerConfigInformer{
		name:     name,
		callback: cbs,
		mutex:    sync.RWMutex{},
		informer: cache.NewSharedInformer(
			lw,
			&corev1.ConfigMap{},
			notReSync,
		),
	}
	return sc
}

func (sc *ServerConfigInformer) Start(ctx context.Context) {
	func() {
		sc.mutex.Lock()
		defer sc.mutex.Unlock()

		sc.started = true
	}()

	var ch chan struct{}
	go func() {
		select {
		case <-ctx.Done():
			ch <- struct{}{}
		}
	}()

	sc.informer.Run(ch)
	sc.started = false
}

func (sc *ServerConfigInformer) addHandlers() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.started {
		return fmt.Errorf("configMap informer already started")
	}

	registration, err := sc.informer.AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: sc.filterServerConfigMap,
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					sc.EventHandlerFunc(obj, false)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					// we don't really care what old Obj is
					sc.EventHandlerFunc(newObj, false)
				},
				DeleteFunc: func(obj interface{}) {
					sc.EventHandlerFunc(obj, true)
				},
			},
		},
	)
	sc.registration = registration

	return err
}

func (sc *ServerConfigInformer) filterServerConfigMap(obj interface{}) bool {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return false
	}
	// TODO: Add more criteria: we need to check not only the name of the configMap, but also the namespace.
	if configMap.Name == defaultConfigMap {
		return true
	}
	return false
}

func (sc *ServerConfigInformer) EventHandlerFunc(obj interface{}, isDelete bool) {
	var reset bool
	var jaegerOperator *jaegerv1a1.JaegerOperator

	if !isDelete {
		// filterFunc helps us filter out `types != ConfigMap`,
		// so we can be sure this is a ConfigMap resource
		configMap := obj.(*corev1.ConfigMap)
		jaegerOperator = &jaegerv1a1.JaegerOperator{}

		// If an error occurs during the conversion,
		// don't worry. We will fall back to using the default values
		if jaegerOperator = convertToJaegerOperator(configMap); jaegerOperator == nil {
			reset = true
		}
	} else {
		reset = true
	}
	server := JaegerOperatorToServer(jaegerOperator, reset)
	for _, cb := range sc.callback {
		cb(server)
	}
}

// convertToJaegerOperator Get the configuration from the ConfigMap and deserialize it with the codec.
// The configuration will also be validated
func convertToJaegerOperator(configMap *corev1.ConfigMap) *jaegerv1a1.JaegerOperator {
	content, ok := configMap.Data[defaultConfigMapKey]
	if !ok {
		informerLogger.Info("configMap does not contain key " + defaultConfigMapKey)
		return nil
	}

	decoder := serializer.NewCodecFactory(jaeger.GetScheme()).UniversalDeserializer()
	object, gvk, err := decoder.Decode([]byte(content), nil, nil)
	if (err != nil) ||
		(gvk.Group != jaegerv1a1.GroupVersion.Group ||
			gvk.Version != jaegerv1a1.GroupVersion.Version ||
			gvk.Kind != jaegerv1a1.JaegerOperatorKind) {

		informerLogger.Info("failed to convert object to JaegerOperator type")
		return nil
	}
	jaegerOperator := object.(*jaegerv1a1.JaegerOperator)
	return jaegerOperator
}
