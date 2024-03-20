package translate

import (
	"fmt"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

var _ StrategyRender = (*AllInOneRender)(nil)
var _ StrategyRender = (*DistributeRender)(nil)

type StrategyRender interface {
	Deployment() ([]*appsv1.Deployment, error)
	ServiceAccount() (*corev1.ServiceAccount, error)
	ConfigMap() (*corev1.ConfigMap, error)
	Service() ([]*corev1.Service, error)
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

func NamespacedName(instance *jaegerv1a1.Jaeger) string {
	ret := strings.Split(types.NamespacedName{
		Namespace: instance.Name,
		Name:      instance.Namespace,
	}.String(), "/")

	return fmt.Sprintf("%s-%s", ret[0], ret[1])
}

func GetStrategy(instance *jaegerv1a1.Jaeger) string {
	return string(instance.GetDeploymentType())
}

func ComponentLabels(component string, componentName string, instance *jaegerv1a1.Jaeger) map[string]string {
	return utils.MergeCommonMap(
		utils.Labels(
			componentName, component,
			GetStrategy(instance)),
		instance.GetCommonSpecLabels(),
	)
}

func ComponentName(instanceName, suffix string) string {
	return fmt.Sprintf("%s-%s", instanceName, suffix)
}
