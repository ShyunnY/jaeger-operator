package translate

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenericServiceAccount Generate a generic Kubernetes Service Account
func GenericServiceAccount(instance *jaegerv1a1.Jaeger) *corev1.ServiceAccount {
	saLabels := utils.MergeCommonMap(
		utils.Labels(instance.Name, "service-account",
			GetStrategy(instance)),
		instance.Labels,
	)

	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        NamespacedName(instance),
			Namespace:   instance.Namespace,
			Labels:      saLabels,
			Annotations: instance.GetAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
	}
}
