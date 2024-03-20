package utils

import (
	"testing"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestMergePodEnv(t *testing.T) {

	cases := []struct {
		caseName string
		envs     []jaegerv1a1.EnvSetting
		expect   []corev1.EnvVar
	}{
		{
			caseName: "envs is empty",
			envs:     nil,
			expect: []corev1.EnvVar{
				{
					Name:  "SPAN_STORAGE_TYPE",
					Value: "memory",
				},
				{
					Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
					Value: ":9411",
				},
				{
					Name:  "JAEGER_DISABLED",
					Value: "false",
				},
				{
					Name:  "COLLECTOR_OTLP_ENABLED",
					Value: "true",
				},
			},
		},
		{
			caseName: "envs is not empty",
			envs: []jaegerv1a1.EnvSetting{
				{
					Name:  "JAEGER_DISABLED",
					Value: "true",
				},
			},
			expect: []corev1.EnvVar{
				{
					Name:  "SPAN_STORAGE_TYPE",
					Value: "memory",
				},
				{
					Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
					Value: ":9411",
				},
				{
					Name:  "JAEGER_DISABLED",
					Value: "true",
				},
				{
					Name:  "COLLECTOR_OTLP_ENABLED",
					Value: "true",
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			//actual := MergePodEnv(tc.envs...)
			//assert.True(t, checkEnv(actual, tc.expect))
		})

	}

}

func checkEnv(actual, expect []corev1.EnvVar) bool {

	actualMap := make(map[string]corev1.EnvVar)
	expectMap := make(map[string]corev1.EnvVar)

	for _, envVar := range actual {
		actualMap[envVar.Name] = envVar
	}

	for _, envVar := range expect {
		expectMap[envVar.Name] = envVar
	}

	for k, envVar := range actualMap {
		if expectEnv, ok := expectMap[k]; !ok || envVar.Name != expectEnv.Name || envVar.Value != expectEnv.Value {
			return false
		}
	}

	return true
}

func TestMergePodArgs(t *testing.T) {

	cases := []struct {
		caseName string
		args     []string
		expect   []string
	}{
		{
			caseName: "envs is empty",
			args:     nil,
			expect: []string{
				"--memory.max-traces=100000",
			},
		},
		{
			caseName: "envs is empty",
			args: []string{
				"--sampling.strategies-file=/etc/jaeger/sampling/sampling.json",
			},
			expect: []string{
				"--memory.max-traces=100000",
				"--sampling.strategies-file=/etc/jaeger/sampling/sampling.json",
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			//actual := MergePodArgs(tc.args...)
			//assert.Equal(t, tc.expect, actual)
		})

	}

}
