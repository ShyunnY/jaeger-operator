package v1alpha1

// TODO: 将其配置成可以通过ConfigMap动态配置的.

// JaegerOperator Define the configuration of the Jaeger Operator
type JaegerOperator struct {
	// Telemetry Define the configuration of the observability of the Jaeger Operator
	Telemetry *JaegerOperatorTelemetry `json:"telemetry,omitempty"`
}

type JaegerOperatorTelemetry struct {
	// Metrics Define the metrics configuration in the Jaeger Operator
	Metrics *JaegerOperatorMetricsSettings `json:"metrics,omitempty"`
}

type JaegerOperatorMetricsSettings struct {
	// DisablePrometheus Define explicitly disable Prometheus
	DisablePrometheus bool `json:"disablePrometheus"`

	// Sink Define the receiver configuration for the Open Telemetry Sink
	Sink []OpenTelemetrySink `json:"sink,omitempty"`
}

type OpenTelemetrySink struct {
	URL *string `json:"url,omitempty"`
}
