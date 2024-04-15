package config

import (
	"k8s.io/utils/ptr"
	"k8s.io/utils/set"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

const (
	defaultJaegerOperatorName = "jaeger-operator"
)

type Server struct {
	JaegerOperatorName string
	Logger             logging.Logger

	// TODO: 我们应该使用NamespaceSelector
	NamespaceSet  set.Set[string]
	Observability ObservabilityConfig
}

type ObservabilityConfig struct {
	// Jaeger-Operator metrics config
	Metric *Metrics

	// Jaeger-Operator trace config
	Trace *Traces
}

type Metrics struct {
	Sink *OpenTelemetrySink

	DisablePrometheus *bool
}

type OpenTelemetrySink struct {
	Host     string
	Port     string
	Protocol string
}

type Traces struct {
	Endpoint *string
}

func (s *Server) DisablePrometheus() bool {
	if s.Observability.Metric != nil &&
		s.Observability.Metric.DisablePrometheus != nil &&
		*s.Observability.Metric.DisablePrometheus == true {
		return true
	}

	return false
}

func JaegerOperatorToServer(operator *jaegerv1a1.JaegerOperator, reset bool) *Server {
	// Create a Server with default values
	server := setDefaultServer()
	if reset {
		return server
	}

	// metadata setting
	var logLevel = consts.LogLevelDInfo
	if operator.Metadata != nil {
		if operator.Metadata.Name != nil && len(*operator.Metadata.Name) != 0 {
			server.JaegerOperatorName = *operator.Metadata.Name
		}
	}

	// telemetry setting
	if operator.Telemetry != nil {
		if operator.Telemetry.Metrics != nil {
			server.Observability.Metric.DisablePrometheus = ptr.To[bool](operator.Telemetry.Metrics.DisablePrometheus)
			// TODO: sink
		}

		if operator.Telemetry.Logging != nil {
			if operator.Telemetry.Logging.Level != nil && len(*operator.Telemetry.Logging.Level) != 0 {
				logLevel = consts.LogLevel(*operator.Telemetry.Logging.Level)
			}
		}
	}
	server.Logger = logging.NewLogger(logLevel).WithName(server.JaegerOperatorName)

	return server
}

func setDefaultServer() *Server {
	server := new(Server)
	server.JaegerOperatorName = defaultJaegerOperatorName
	server.Logger = logging.NewLogger(consts.LogLevelDInfo).WithName(server.JaegerOperatorName)
	server.Observability = ObservabilityConfig{
		Metric: &Metrics{},
		Trace:  &Traces{},
	}
	return server
}
