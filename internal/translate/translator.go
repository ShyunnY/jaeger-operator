package translate

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Translator struct {
	Logger logging.Logger

	InfraIRMap  *message.InfraIRMaps
	StatusIRMap *message.StatusIRMaps
}

// Translate
// Jaeger resources are translated into kubernetes resources according to strategy,
// which are passed to infrastructure as infraIR for construction
func (t *Translator) Translate(instance *jaegerv1a1.Jaeger) error {

	infraIR := new(message.InfraIR)
	instance.Status.Phase = "Failed"
	defer func() {
		t.StatusIRMap.Map.Store(types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		}, &instance.Status)
	}()

	infraIR.InstanceMedata = instance.ObjectMeta
	infraIR.Strategy = string(instance.Spec.Type)

	var strRender StrategyRender
	switch infraIR.Strategy {
	case string(jaegerv1a1.AllInOneType):
		strRender = &AllInOneRender{
			instance: instance,
		}
	case string(jaegerv1a1.Distribute):
		strRender = &DistributeRender{
			instance: instance,
		}
	default:
		t.Logger.Info("unsupported deployment strategy")

		status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "unsupported deployment strategy")
		return fmt.Errorf("current deployment type is not supported: %s", infraIR.Strategy)
	}

	// render strategy resources
	if sa, err := strRender.ServiceAccount(); err != nil {
		t.Logger.Error(err, "failed to render service account resource",
			"instance", instance.Name,
		)

		status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "failed to render service account resource")
	} else {
		infraIR.AddResources(sa)
	}

	if cm, err := strRender.ConfigMap(); err != nil {
		t.Logger.Error(err, "failed to render configmap resource",
			"instance", instance.Name,
		)

		status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "failed to render configmap resource")
	} else {
		infraIR.AddResources(cm)
	}

	if deploy, err := strRender.Deployment(); err != nil {
		t.Logger.Error(err, "failed to render deployment resource",
			"instance", instance.Name,
		)

		status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "failed to render deployment resource")
	} else {
		infraIR.AddResources(deploy)
	}

	// TODO: need to deal with services with different strategy and with multiple Service resources
	var services []*corev1.Service
	var err error
	if services, err = strRender.Service(); err != nil {
		t.Logger.Error(err, "failed to render service resource",
			"instance", instance.Name,
		)

		status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "failed to render deployment resource")
	} else if len(services) != 0 {
		infraIR.AddResources(services)
	}

	if instance.EnableHTTPRoute() {
		if httpRoute, err := processHTTPRoute(instance, services); err != nil {
			t.Logger.Error(err, "failed to render httpRoute resource",
				"instance", instance.Name,
			)

			status.SetJaegerCondition(instance, "Error", metav1.ConditionFalse, "Translate", "failed to render httpRoute resource")
		} else if len(httpRoute) != 0 {
			infraIR.AddResources(httpRoute)
		}
	}

	if instance.Status.Conditions == nil || len(instance.Status.Conditions) == 0 {
		status.SetJaegerCondition(instance, "Success", metav1.ConditionTrue, "Translate", "success to translate resource")
		instance.Status.Phase = "Success"
	}

	// push ir
	t.InfraIRMap.Store(instance.Name, infraIR)
	return nil
}
