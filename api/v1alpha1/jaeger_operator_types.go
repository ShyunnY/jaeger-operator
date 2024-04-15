package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	JaegerOperatorKind = "JaegerOperator"
)

// +kubebuilder:object:root=true

// JaegerOperator Define the configuration of the Jaeger Operator
type JaegerOperator struct {
	metav1.TypeMeta

	// Metadata Define Jaeger Operator related metadata
	Metadata *JaegerOperatorMetadata `json:"metadata"`

	// Telemetry Define the configuration of the observability of the Jaeger Operator
	Telemetry *JaegerOperatorTelemetry `json:"telemetry,omitempty"`
}

type JaegerOperatorMetadata struct {

	// Name Define the name of the Jaeger Operator
	Name *string `json:"name,omitempty"`
}

type JaegerOperatorTelemetry struct {
	// Metrics Define the metrics configuration in the Jaeger Operator
	Metrics *JaegerOperatorMetricsSettings `json:"metrics,omitempty"`

	// Logging Define the logging configuration in the Jaeger Operator
	Logging *JaegerOperatorLoggingSettings `json:"logging,omitempty"`
}

type JaegerOperatorLoggingSettings struct {
	Level *string `json:"level,omitempty"`
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

func init() {
	SchemeBuilder.Register(&JaegerOperator{})
}
