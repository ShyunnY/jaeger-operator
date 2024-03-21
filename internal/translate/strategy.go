package translate

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
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

func buildOptions(instance *jaegerv1a1.Jaeger) {

	defaultArgs := []string{}
	defaultEnvs := []corev1.EnvVar{}

	switch instance.GetDeploymentType() {
	case jaegerv1a1.AllInOneType:
		allInOneArgs := []string{
			"--memory.max-traces=100000",
		}
		defaultArgs = append(defaultArgs, allInOneArgs...)

		allInOneEnvs := []corev1.EnvVar{
			{
				Name:  "SPAN_STORAGE_TYPE",
				Value: "memory",
			},
			{
				Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
				Value: ":9411",
			},
			{
				Name:  "JAEGER_DISABLED",
				Value: "false",
			},
			{
				Name:  "COLLECTOR_OTLP_ENABLED",
				Value: "true",
			},
		}
		defaultEnvs = append(defaultEnvs, allInOneEnvs...)

		instance.Spec.Components.AllInOne.Args = utils.MergePodArgs(
			defaultArgs,
			instance.Spec.Components.AllInOne.Args...,
		)
		instance.Spec.Components.AllInOne.Envs = utils.MergePodEnv(
			defaultEnvs,
			instance.Spec.Components.AllInOne.Envs...,
		)
	case jaegerv1a1.Distribute:
		defaultArgs, defaultEnvs = buildStorageOptions(instance)

		instance.Spec.Components.Collector.Args = utils.MergePodArgs(
			defaultArgs,
			instance.Spec.Components.Collector.Args...,
		)
		instance.Spec.Components.Collector.Envs = utils.MergePodEnv(
			defaultEnvs,
			instance.Spec.Components.Collector.Envs...,
		)

		instance.Spec.Components.Query.Args = utils.MergePodArgs(
			defaultArgs,
			instance.Spec.Components.Query.Args...,
		)
		instance.Spec.Components.Query.Envs = utils.MergePodEnv(
			defaultEnvs,
			instance.Spec.Components.Query.Envs...,
		)
	}

}

func buildStorageOptions(instance *jaegerv1a1.Jaeger) ([]string, []corev1.EnvVar) {

	storageArgs := []string{}
	storageEnvs := []corev1.EnvVar{}

	switch instance.GetStorageType() {
	case jaegerv1a1.MemoryStorageType:
		return storageArgs, storageEnvs
	case jaegerv1a1.ElasticSearchStorage:

		var url = "http://elasticsearch:9200"
		if len(instance.Spec.Components.Storage.Es.URL) != 0 {
			url = instance.Spec.Components.Storage.Es.URL
		}

		storageArgs = append(storageArgs, []string{
			fmt.Sprintf("--es.server-urls=%s", url),
		}...)

		storageEnvs = append(storageEnvs, []corev1.EnvVar{
			{
				Name:  "SPAN_STORAGE_TYPE",
				Value: "elasticsearch",
			},
			{
				Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
				Value: ":9411",
			},
			{
				Name:  "JAEGER_DISABLED",
				Value: "false",
			},
			{
				Name:  "COLLECTOR_OTLP_ENABLED",
				Value: "true",
			},
		}...)
	}

	return storageArgs, storageEnvs
}

func getQueryPort() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "grpc-query",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 16685,
		},
		{
			Name:          "http-query",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 16686,
		},
		{
			Name:          "admin",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 16687,
		},
	}
}

func getCollectorPort(enableOTLP bool) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "zipkin",
			ContainerPort: 9411,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "binary-thrift",
			ContainerPort: 14268,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "admin-http",
			ContainerPort: 14269,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	if enableOTLP {
		ports = append(ports, getOTLPPort()...)
	}

	return ports
}

func getOTLPPort() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "oltp-grpc",
			ContainerPort: consts.OltpGrpcPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "oltp-http",
			ContainerPort: consts.OltpHTTPPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}
