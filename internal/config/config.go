package config

import (
	"k8s.io/utils/set"

	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

type Server struct {
	Logger       logging.Logger
	NamespaceSet set.Set[string]

	JaegerOperatorName string

	// observability config
	Observability
}

type Observability struct {
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
		s.Metric.DisablePrometheus != nil &&
		*s.Metric.DisablePrometheus == true {
		return true
	}

	return false
}
