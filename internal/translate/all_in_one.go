package translate

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
)

type AllInOneRender struct {
	instance *jaegerv1a1.Jaeger

	labels      map[string]string
	annotations map[string]string
}

// Deployment Returns the Deployment of the expected AllInOne Strategy
func (r *AllInOneRender) Deployment() (*appsv1.Deployment, error) {

	envs := []corev1.EnvVar{
		{
			Name:  "SPAN_STORAGE_TYPE",
			Value: "memory",
		},
		{
			Name:  "METRICS_STORAGE_TYPE",
			Value: "",
		},
		{
			Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
			Value: ":9411",
		},
		{
			Name:  "JAEGER_DISABLED",
			Value: "true",
		},
	}

	ports := []corev1.ContainerPort{
		{
			ContainerPort: 5775,
			Name:          "zk-compact-trft",
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 5778,
			Name:          "config-rest",
		},
		{
			ContainerPort: 6831,
			Name:          "jg-compact-trft",
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 6832,
			Name:          "jg-binary-trft",
			Protocol:      corev1.ProtocolUDP,
		},
		{
			ContainerPort: 9411,
			Name:          "zipkin",
		},
		{
			ContainerPort: 14267,
			Name:          "c-tchan-trft",
		},
		{
			ContainerPort: 14268,
			Name:          "c-binary-trft",
		},
		{
			ContainerPort: 16685,
			Name:          "grpc-query",
		},
		{
			ContainerPort: 16686,
			Name:          "query",
		},
		{
			ContainerPort: 14269,
			Name:          "admin-http",
		},
		{
			ContainerPort: 14250,
			Name:          "grpc",
		},
	}

	depLabels := utils.MergeCommonMap(utils.Labels(r.instance.Name, "deployment", string(r.GetStrategy())), r.labels)
	podLabels := utils.MergeCommonMap(utils.Labels(r.instance.Name, "pod", string(r.GetStrategy())), r.labels)

	var replicas int32 = 1
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        r.Name(),
			Namespace:   r.instance.Namespace,
			Labels:      depLabels,
			Annotations: r.annotations,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(r.instance),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: corev1.PodTemplateSpec{
				// Pod的 metadata元数据
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: r.annotations,
				},
				Spec: corev1.PodSpec{
					// ImagePullSecrets: commonSpec.ImagePullSecrets,
					Containers: []corev1.Container{
						{
							Image: "jaegertracing/all-in-one:1.54.0",
							Name:  "jaeger",
							// Args:  options,
							// TODO: 需要添加otlp env/port
							Env:   envs,
							Ports: ports,
							// LivenessProbe: livenessProbe,
						},
					},
					ServiceAccountName: r.Name(),
				},
			},
		},
	}, nil
}

// ConfigMap Returns the ConfigMap of the expected AllInOne Strategy
func (r *AllInOneRender) ConfigMap() (*corev1.ConfigMap, error) {
	return nil, nil
}

// Service Returns the Service of the expected AllInOne Strategy
func (r *AllInOneRender) Service() ([]*corev1.Service, error) {
	// services: agent,collect,query
	services := []*corev1.Service{}
	selector := utils.MergeCommonMap(utils.Labels(r.instance.Name, "pod", string(r.GetStrategy())), r.labels)

	queryService := QueryService(r.instance)
	queryService.Spec.Selector = selector
	queryService.Labels = utils.MergeCommonMap(utils.Labels(r.instance.Name, "service", string(r.GetStrategy())), r.labels)
	queryService.Annotations = r.annotations
	services = append(services, queryService)

	collectorServices := CollectorServices(r.instance)
	for i := range collectorServices {
		collectorServices[i].Spec.Selector = selector
		collectorServices[i].Labels = utils.MergeCommonMap(utils.Labels(r.instance.Name, "service", string(r.GetStrategy())), r.labels)
		collectorServices[i].Annotations = r.annotations
	}
	services = append(services, collectorServices...)

	return services, nil
}

// ServiceAccount Returns the ServiceAccount of the expected AllInOne Strategy
func (r *AllInOneRender) ServiceAccount() (*corev1.ServiceAccount, error) {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name(),
			Namespace: r.instance.Namespace,
			Labels: utils.MergeCommonMap(utils.Labels(r.instance.Name, "service-account",
				string(r.GetStrategy())), r.labels),
			Annotations: r.annotations,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(r.instance),
			},
		},
	}, nil
}

func (r *AllInOneRender) GetStrategy() jaegerv1a1.DeploymentType {
	return jaegerv1a1.AllInOneType
}

func (r *AllInOneRender) Name() string {
	return utils.GetHashedName(
		types.NamespacedName{
			Namespace: r.instance.Name,
			Name:      r.instance.Namespace,
		},
	)
}
