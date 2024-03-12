package translate

import (
	"fmt"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
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

	if instance.Spec.Type == "" {
		infraIR.Strategy = string(jaegerv1a1.AllInOneType)
	} else {
		infraIR.Strategy = string(instance.Spec.Type)
	}

	// Compute commonSpec and incorporate common metadata
	mergeCommonLabels := utils.MergeCommonMap(instance.Labels, instance.Spec.CommonSpec.Metadata.Labels)
	mergeCommonAnnotations := utils.MergeCommonMap(instance.Annotations, instance.Spec.CommonSpec.Metadata.Annotations)

	infraIR.InstanceName = instance.Name
	infraIR.InstanceNamespace = instance.Namespace

	var strRender StrategyRender
	switch infraIR.Strategy {
	case string(jaegerv1a1.AllInOneType):
		strRender = &AllInOneRender{
			instance:    instance,
			labels:      mergeCommonLabels,
			annotations: mergeCommonAnnotations,
		}
	case string(jaegerv1a1.Distribute):
		// TODO: production resources render
	default:
		t.Logger.Info("failed to get strategy render")

		return fmt.Errorf("current deployment type is not supported: %s", infraIR.Strategy)
	}

	// render strategy resources
	// TODO: If there is an error in rendering the resource, an error message will be added to the condition

	if sa, err := strRender.ServiceAccount(); err != nil {
		t.Logger.Error(err, "failed to render service account resource",
			"instance", instance.Name,
		)
		return err
	} else {
		infraIR.AddResources(sa)
	}

	if cm, err := strRender.ConfigMap(); err != nil {
		t.Logger.Error(err, "failed to render configmap resource",
			"instance", instance.Name,
		)
		return nil
	} else {
		infraIR.AddResources(cm)
	}

	if deploy, err := strRender.Deployment(); err != nil {
		t.Logger.Error(err, "failed to render deployment resource",
			"instance", instance.Name,
		)
		return nil
	} else {
		infraIR.AddResources(deploy)
	}

	// TODO: need to deal with services with different strategy and with multiple Service resources
	if services, err := strRender.Service(); err != nil {
		t.Logger.Error(err, "failed to render service resource",
			"instance", instance.Name,
		)
		return nil
	} else if len(services) != 0 {
		infraIR.AddResources(services)
	}

	t.InfraIRMap.Store(instance.Name, infraIR)

	return nil
}
