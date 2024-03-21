package translate

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
)

type DistributeRender struct {
	instance *jaegerv1a1.Jaeger
}

func (d *DistributeRender) Deployment() ([]*appsv1.Deployment, error) {

	buildOptions(d.instance)

	queryDeploy := QueryDeploy(d.instance)
	collectorDeploy := CollectorDeploy(d.instance)

	return []*appsv1.Deployment{queryDeploy, collectorDeploy}, nil
}

func (d *DistributeRender) ServiceAccount() (*corev1.ServiceAccount, error) {
	return GenericServiceAccount(d.instance), nil
}

func (d *DistributeRender) ConfigMap() (*corev1.ConfigMap, error) {
	return nil, nil
}

func (d *DistributeRender) Service() ([]*corev1.Service, error) {

	services := []*corev1.Service{}
	selector := ComponentLabels(
		"pod",
		ComponentName(NamespacedName(d.instance), consts.CollectorComponent),
		d.instance,
	)

	// In distribute deployment mode, we build three Service: Query,Collector,Collector-Headless
	// Indeed, Collector and Query each select a Pod under a different Deployment
	for _, collectorSvc := range CollectorServices(d.instance) {
		collectorSvc := collectorSvc
		collectorSvc.Spec.Selector = selector
		services = append(services, collectorSvc)
	}

	queryService := QueryService(d.instance)
	queryService.Spec.Selector = ComponentLabels(
		"pod",
		ComponentName(NamespacedName(d.instance), consts.QueryComponent),
		d.instance,
	)

	return append(services, queryService), nil
}
