package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/ShyunnY/jaeger-operator/internal/config"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

var (
	tracingLog = logging.NewLogger(consts.LogLevelDInfo).WithName("tracing")
)

// New TODO: 需要完善trace, 并且添加日志记录等
func New(cfg *config.Server) error {

	if !enableTraceObservability(cfg) {
		return nil
	}

	if err := BuildTracer(
		cfg.JaegerOperatorName,
		consts.DefaultAllNamespace,
		*cfg.Observability.Trace.Endpoint,
	); err != nil {
		tracingLog.Error(err, "failed to build trace provider")

		return err
	}

	tracingLog.Info("success to build trace provider")
	return nil
}

// BuildTracer Building a Global tracer
func BuildTracer(instance, namespace, endpoint string) error {

	export, err := buildJaegerTraceExport(endpoint)
	if err != nil {
		return err
	}

	attr := []attribute.KeyValue{
		semconv.ServiceNameKey.String("jaeger-operator"),
		semconv.ServiceVersionKey.String("0.1.0"),
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

func enableTraceObservability(cfg *config.Server) bool {
	if cfg.Observability.Trace != nil &&
		cfg.Observability.Trace.Endpoint != nil &&
		len(*cfg.Observability.Trace.Endpoint) != 0 {
		return true
	}
	return false
}
