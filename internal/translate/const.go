package translate

import (
	"fmt"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	corev1 "k8s.io/api/core/v1"
)

// jaeger port
const (
	oltpGrpcPort = 4317
	oltpHTTPPort = 4318

	adminPort = 14269
)

// jaeger component
const (
	// TODO: Adding more components later?
	collectorComponent = "collector"
	queryComponent     = "query"
	allIneOneComponent = "allinone"
)

const (
	queryServiceType             = "query-svc"
	collectorServiceType         = "collector-svc"
	collectorServiceHeadlessType = "collector-headless-svc"
)

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
			ContainerPort: oltpGrpcPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "oltp-http",
			ContainerPort: oltpHTTPPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}
