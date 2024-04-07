package consts

// jaeger port
const (
	OltpGrpcPort = 4317
	OltpHTTPPort = 4318

	AdminPort = 14269
)

// jaeger component
const (
	// TODO: Adding more components later?
	CollectorComponent = "collector"
	QueryComponent     = "query"
	AllIneOneComponent = "allinone"
)

// jaeger service type
const (
	QueryServiceType             = "query-svc"
	CollectorServiceType         = "collector-svc"
	CollectorServiceHeadlessType = "collector-headless-svc"
)

// jaeger metrics
const (
	MetricsHost = "0.0.0.0"
	MetricsPort = "10424"

	OtelHost   = "localhost"
	OtelPort   = "4318"
	OtelProtol = "http"
)

type LogLevel string

const (

	// LogLevelDebug define debug level logging
	LogLevelDebug LogLevel = "debug"

	// LogLevelDInfo define info level logging
	LogLevelDInfo LogLevel = "info"

	// LogLevelWarn define warn level logging
	LogLevelWarn LogLevel = "warn"

	// LogLevelError define error level logging
	LogLevelError LogLevel = "error"
)

const (
	ReconciliationTracer string = "operator/reconciliation"
)

const (
	OperatorName        = "jaeger-operator"
	DefaultAllNamespace = "all_namespace"
)
