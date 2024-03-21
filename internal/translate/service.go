package translate

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
)

// QueryService Build a Service for accessing Jaeger queries
func QueryService(instance *jaegerv1a1.Jaeger) *corev1.Service {

	queryPorts := getQueryPort()
	name := ComponentName(instance.Name, consts.QueryServiceType)
	queryService := expectServiceSpec(
		name,
		jaegerv1a1.QueryServiceTarget,
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

	collectorName := ComponentName(instance.Name, consts.CollectorServiceType)
	collectorSvc := expectServiceSpec(
		collectorName,
		jaegerv1a1.CollectorServiceTarget,
		instance,
		convertServicePort(ports),
	)
	labels := utils.MergeCommonMap(collectorSvc.Labels, instance.GetCommonSpecLabels())
	annotations := utils.MergeCommonMap(collectorSvc.Annotations, instance.GetAnnotations())

	collectorSvc.Labels = labels
	collectorSvc.Annotations = annotations
	retServices = append(retServices, collectorSvc)

	collectorHeadlessName := ComponentName(instance.Name, consts.CollectorServiceHeadlessType)
	collectorHeadlessSvc := expectServiceSpec(
		collectorHeadlessName,
		jaegerv1a1.CollectorServiceTarget,
		instance,
		convertServicePort(ports),
	)
	collectorHeadlessSvc.Labels = labels
	collectorHeadlessSvc.Annotations = annotations

	if serviceType(instance) != corev1.ServiceTypeClusterIP {
		collectorHeadlessSvc.Spec.Type = corev1.ServiceTypeClusterIP
	}
	// We set the cluster IP too None to provide a headless service
	collectorHeadlessSvc.Spec.ClusterIP = "None"

	retServices = append(retServices, collectorHeadlessSvc)

	return retServices
}

func expectServiceSpec(name string, target jaegerv1a1.JaegerServiceTarget, instance *jaegerv1a1.Jaeger, ports []corev1.ServicePort) *corev1.Service {

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels: utils.MergeCommonMap(map[string]string{
				jaegerv1a1.ServiceTargetLabelKey: string(target),
			}, ComponentLabels("service", name, instance)),
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
