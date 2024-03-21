package translate

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
)

type AllInOneRender struct {
	instance *jaegerv1a1.Jaeger
}

// Deployment Returns the Deployment of the expected AllInOne Strategy
func (r *AllInOneRender) Deployment() ([]*appsv1.Deployment, error) {
	buildOptions(r.instance)
	return []*appsv1.Deployment{allInOneDeploy(r.instance)}, nil
}

// ConfigMap Returns the ConfigMap of the expected AllInOne Strategy
func (r *AllInOneRender) ConfigMap() (*corev1.ConfigMap, error) {
	return nil, nil
}

// Service Returns the Service of the expected AllInOne Strategy
func (r *AllInOneRender) Service() ([]*corev1.Service, error) {

	services := []*corev1.Service{}
	selector := ComponentLabels(
		"pod",
		ComponentName(NamespacedName(r.instance), consts.AllIneOneComponent),
		r.instance,
	)
	// In allInOne deployment mode, we build three Service: Query,Collector,Collector-Headless
	// In fact, all three services choose the Pod deployed by the same Deployment
	for _, collectorSvc := range CollectorServices(r.instance) {
		collectorSvc := collectorSvc
		collectorSvc.Spec.Selector = selector
		services = append(services, collectorSvc)
	}

	queryService := QueryService(r.instance)
	queryService.Spec.Selector = selector

	return append(services, queryService), nil
}

// ServiceAccount Returns the ServiceAccount of the expected AllInOne Strategy
func (r *AllInOneRender) ServiceAccount() (*corev1.ServiceAccount, error) {
	return GenericServiceAccount(r.instance), nil
}
