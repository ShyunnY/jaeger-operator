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

const (
	QueryServiceType             = "query-svc"
	CollectorServiceType         = "collector-svc"
	CollectorServiceHeadlessType = "collector-headless-svc"
)
