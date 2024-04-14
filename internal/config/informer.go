package config

import (
	"sync"
	"time"

	"github.com/ShyunnY/jaeger-operator/internal/logging"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	notReSync      = time.Duration(0)
	informerLogger = logging.DefaultLogger()
)

type ServerConfigInformer struct {
	name  string
	mutex sync.RWMutex

	callback     func(*Server)
	informer     cache.SharedInformer
	registration cache.ResourceEventHandlerRegistration
}

func New() *ServerConfigInformer {

	informer := cache.NewSharedInformer(
		nil,
		&corev1.ConfigMap{},
		notReSync,
	)

	sc := &ServerConfigInformer{
		informer: informer,
	}

	return sc
}

func (sc *ServerConfigInformer) addHandlers() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	registration, err := sc.informer.AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: sc.filterServerConfigMap,
			Handler: cache.ResourceEventHandlerFuncs{
				// Add时需要将Config设置到Server中.
				AddFunc: func(obj interface{}) {
					// TODO: 序列化cm,
					// filterFunc helps us filter out `types != ConfigMap`,
					// so we can be sure this is a ConfigMap resource
					configMap := obj.(*corev1.ConfigMap)
				},
				// Update时需要将Config设置到Server中.
				UpdateFunc: func(oldObj, newObj interface{}) {

				},
				// Delete需要将默认的Config设置到Server中.
				DeleteFunc: func(obj interface{}) {

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

	// TODO: 我们不仅需要判断configMap的name, 还需要判断namespace.
	if configMap.Name == "" {
		return true
	}
	return false
}

// TODO: 1. 通过codeFactory -> 2. 判断gvk -> 3. 默认值设置
func handlerConfigMap(configMap *corev1.ConfigMap) (*Server, error) {

	return nil, nil
}
