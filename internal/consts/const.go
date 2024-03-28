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
	OtelPort   = "4317"
	OtelProtol = "4317"
)
