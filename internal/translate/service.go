package translate

import (
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

// QueryService Build a Service for accessing Jaeger queries
func QueryService(instance *jaegerv1a1.Jaeger) *corev1.Service {

	queryPorts := getQueryPort()
	name := ComponentName(instance.Name, queryServiceType)
	queryService := expectServiceSpec(
		name,
		instance,
		convertServicePort(queryPorts),
	)

	queryService.Labels = utils.MergeCommonMap(queryService.Labels, instance.GetCommonSpecLabels())
	queryService.Annotations = utils.MergeCommonMap(queryService.Annotations, instance.GetAnnotations())

	return queryService
}

// CollectorServices
// We provide two kinds of collect service: HeadlessService and normal Service
func CollectorServices(instance *jaegerv1a1.Jaeger) []*corev1.Service {

	var retServices []*corev1.Service
	ports := getCollectorPort(true)

	collectorName := ComponentName(instance.Name, collectorServiceType)
	collectorSvc := expectServiceSpec(
		collectorName,
		instance,
		convertServicePort(ports),
	)
	labels := utils.MergeCommonMap(collectorSvc.Labels, instance.GetCommonSpecLabels())
	annotations := utils.MergeCommonMap(collectorSvc.Annotations, instance.GetAnnotations())

	collectorSvc.Labels = labels
	collectorSvc.Annotations = annotations
	retServices = append(retServices, collectorSvc)

	collectorHeadlessName := ComponentName(instance.Name, collectorServiceHeadlessType)
	collectorHeadlessSvc := expectServiceSpec(
		collectorHeadlessName,
		instance,
		convertServicePort(ports),
	)
	collectorHeadlessSvc.Labels = labels
	collectorHeadlessSvc.Annotations = annotations

	// We set the cluster IP too NONE to provide a headless service
	collectorHeadlessSvc.Spec.ClusterIP = "NONE"

	retServices = append(retServices, collectorHeadlessSvc)

	return retServices
}

func expectServiceSpec(name string, instance *jaegerv1a1.Jaeger, ports []corev1.ServicePort) *corev1.Service {

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				jaegerv1a1.ServiceTargetLabelKey: string(jaegerv1a1.CollectorServiceTarget),
			},
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  serviceType(instance),
			Ports: ports,
		},
	}
}

func serviceType(instance *jaegerv1a1.Jaeger) corev1.ServiceType {
	svcType := jaegerv1a1.ServiceTypeClusterIP
	if len(instance.Spec.CommonSpec.Service.Type) != 0 {
		svcType = instance.Spec.CommonSpec.Service.Type
	}

	return corev1.ServiceType(svcType)
}

func convertServicePort(ports []corev1.ContainerPort) []corev1.ServicePort {

	svcPorts := []corev1.ServicePort{}
	for _, cp := range ports {
		svcPort := corev1.ServicePort{
			Name:       cp.Name,
			Port:       cp.ContainerPort,
			TargetPort: intstr.FromInt32(cp.ContainerPort),
			Protocol:   cp.Protocol,
		}

		svcPorts = append(svcPorts, svcPort)
	}

	return svcPorts
}
