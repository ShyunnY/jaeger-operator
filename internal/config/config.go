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
	Host     string
	Port     string
	Protocol string

	DisablePrometheus *bool
}

type Traces struct {
}

func (s *Server) DisablePrometheus() bool {
	if s.Metric != nil &&
		s.Metric.DisablePrometheus != nil &&
		*s.Metric.DisablePrometheus == true {
		return true
	}

	return false
}
