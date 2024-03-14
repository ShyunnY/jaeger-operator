package translate

import (
	"fmt"

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
			Name:      serviceName(instance.Name, "query"),
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

// CollectorServices
// We provide two kinds of collect service: HeadlessService and normal Service
func CollectorServices(instance *jaegerv1a1.Jaeger) []*corev1.Service {

	var retServices []*corev1.Service

	ports := []corev1.ServicePort{
		{
			Name:       "zipkin",
			Protocol:   corev1.ProtocolTCP,
			Port:       9411,
			TargetPort: intstr.FromInt32(9411),
		},
		{
			Name:       "grpc",
			Protocol:   corev1.ProtocolTCP,
			Port:       14250,
			TargetPort: intstr.FromInt32(14250),
		},
		{
			Name:       "c-tchan-trft",
			Protocol:   corev1.ProtocolTCP,
			Port:       14267,
			TargetPort: intstr.FromInt32(14267),
		},
		{
			Name:       "c-binary-trft",
			Protocol:   corev1.ProtocolTCP,
			Port:       14268,
			TargetPort: intstr.FromInt32(14268),
		},
	}

	clusterService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(instance.Name, "collect"),
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeClusterIP,
			Ports: ports,
		},
	}
	retServices = append(retServices, clusterService)

	headlessService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(instance.Name, "collect-headless"),
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			Ports:     ports,
			ClusterIP: "None",
		},
	}
	retServices = append(retServices, headlessService)

	return retServices
}

func serviceName(instanceName, suffix string) string {
	return fmt.Sprintf("%s-%s", instanceName, suffix)
}
