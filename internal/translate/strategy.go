package translate

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ StrategyRender = (*AllInOneRender)(nil)

type StrategyRender interface {
	Deployment() (*appsv1.Deployment, error)
	ServiceAccount() (*corev1.ServiceAccount, error)
	ConfigMap() (*corev1.ConfigMap, error)
	Service() ([]*corev1.Service, error)
	GetStrategy() jaegerv1a1.DeploymentType
	Name() string
}

func GetOwnerRef(instance *jaegerv1a1.Jaeger) metav1.OwnerReference {
	controlled := true
	return metav1.OwnerReference{
		APIVersion: instance.APIVersion,
		Kind:       instance.Kind,
		Name:       instance.Name,
		UID:        instance.UID,
		Controller: &controlled,
	}
}
