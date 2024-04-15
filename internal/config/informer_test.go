package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

func TestConvertToJaegerOperator(t *testing.T) {

	cases := []struct {
		caseName  string
		configMap string
		expect    *jaegerv1a1.JaegerOperator
	}{
		{
			caseName: "normal",
			configMap: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-operator
  namespace: default
data:
  jaeger-operator.yaml: |
    apiVersion: tracing.orange.io/v1alpha1
    kind: JaegerOperator
    metadata: 
      name: jaeger-operator
    telemetry:
      logging:
        level: debug
`,
			expect: &jaegerv1a1.JaegerOperator{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "tracing.orange.io/v1alpha1",
					Kind:       "JaegerOperator",
				},
				Metadata: &jaegerv1a1.JaegerOperatorMetadata{
					Name: ptr.To[string]("jaeger-operator"),
				},
				Telemetry: &jaegerv1a1.JaegerOperatorTelemetry{
					Logging: &jaegerv1a1.JaegerOperatorLoggingSettings{
						Level: ptr.To[string]("debug"),
					},
				},
			},
		},
		{
			caseName: "kind or apiVersion does not match",
			configMap: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-operator
  namespace: default
data:
  jaeger-operator.yaml: |
    apiVersion: tracing.orange.io/v1alpha1222
    kind: JaegerOperator111
    metadata: 
      name: jaeger-operator
    telemetry:
      logging:
        level: debug
`,
			expect: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			configMap, err := unmarshalConfigMap(tc.configMap)
			assert.NoError(t, err)
			assert.NotNil(t, configMap)

			actual := convertToJaegerOperator(configMap)
			assert.Equal(t, tc.expect, actual)
		})
	}

}

func unmarshalConfigMap(content string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := yaml.Unmarshal([]byte(content), configMap)
	return configMap, err
}
