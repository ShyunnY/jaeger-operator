package config

import (
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"k8s.io/utils/set"
)

type Server struct {
	Logger       logging.Logger
	NamespaceSet set.Set[string]

	JaegerOperatorName string
}
