package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"

	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// BuildTracer Building a Global tracer
func BuildTracer(instance string, namespace string) error {

	export, err := buildJaegerTraceExport("endpoint")
	if err != nil {
		return err
	}

	attr := []attribute.KeyValue{
		semconv.ServiceNameKey.String("jaeger-operator"),
		semconv.ServiceVersionKey.String("version.Get().Operator"),
		semconv.ServiceNamespaceKey.String(namespace),
	}
	if len(instance) != 0 {
		attr = append(attr, semconv.ServiceInstanceIDKey.String(instance))
	}

	opts := []sdkTrace.TracerProviderOption{
		sdkTrace.WithBatcher(export),
		sdkTrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attr...)),
	}

	tracerProvider := sdkTrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tracerProvider)

	return nil
}

func buildJaegerTraceExport(endpoint string) (sdkTrace.SpanExporter, error) {

	if len(endpoint) == 0 {
		return nil, nil
	}

	jaegerExport, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)),
	)
	if err != nil {
		return nil, err
	}

	return jaegerExport, nil
}
