package translate

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

// QueryService TODO: Supports gateway-api
func QueryService(instance *jaegerv1a1.Jaeger) *corev1.Service {

	// TODO: Do you need more ports?
	queryPorts := []corev1.ServicePort{
		// grpc protocol port
		{
			Name:       "grpc-query",
			Protocol:   corev1.ProtocolTCP,
			Port:       16685,
			TargetPort: intstr.FromInt32(16685),
		},
		// http protocol port
		{
			Name:       "http-query",
			Protocol:   corev1.ProtocolTCP,
			Port:       16686,
			TargetPort: intstr.FromInt32(16686),
		},
		// admin protocol port
		{
			Name:       "admin",
			Protocol:   corev1.ProtocolTCP,
			Port:       16687,
			TargetPort: intstr.FromInt32(16687),
		},
	}

	queryService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: queryPorts,
		},
	}

	return queryService
}
