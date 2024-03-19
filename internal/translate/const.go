package translate

import (
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
	allIneOneComponent = "allInOne"
)

const (
	queryServiceType             = "query-svc"
	collectorServiceType         = "collector-svc"
	collectorServiceHeadlessType = "collector-headless-svc"
)

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
