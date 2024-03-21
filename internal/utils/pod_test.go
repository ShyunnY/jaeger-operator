package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

func TestMergePodEnv(t *testing.T) {

	cases := []struct {
		caseName string
		exist    []corev1.EnvVar
		appended []jaegerv1a1.EnvSetting
		expected []jaegerv1a1.EnvSetting
	}{
		{
			caseName: "append is empty",
			exist: []corev1.EnvVar{
				{
					Name:  "env-a",
					Value: "a",
				},
				{
					Name:  "env-b",
					Value: "b",
				},
			},
			appended: nil,
			expected: []jaegerv1a1.EnvSetting{
				{
					Name:  "env-a",
					Value: "a",
				},
				{
					Name:  "env-b",
					Value: "b",
				},
			},
		},
		{
			caseName: "append is not-empty",
			exist: []corev1.EnvVar{
				{
					Name:  "env-a",
					Value: "a",
				},
				{
					Name:  "env-b",
					Value: "b",
				},
			},
			appended: []jaegerv1a1.EnvSetting{
				{
					Name:  "env-c",
					Value: "c",
				},
			},
			expected: []jaegerv1a1.EnvSetting{
				{
					Name:  "env-a",
					Value: "a",
				},
				{
					Name:  "env-b",
					Value: "b",
				},
				{
					Name:  "env-c",
					Value: "c",
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			actual := MergePodEnv(tc.exist, tc.appended...)
			assert.Equal(t, tc.expected, actual)
		})

	}

}

func TestMergePodArgs(t *testing.T) {

	cases := []struct {
		caseName string
		exist    []string
		appended []string
		expect   []string
	}{
		{
			caseName: "append is empty",
			exist: []string{
				"-a",
				"-b",
			},
			appended: nil,
			expect: []string{
				"-a",
				"-b",
			},
		},
		{
			caseName: "append is not-empty",
			exist: []string{
				"-a",
				"-b",
			},
			appended: []string{
				"-c",
			},
			expect: []string{
				"-a",
				"-b",
				"-c",
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			actual := MergePodArgs(tc.exist, tc.appended...)
			assert.Equal(t, tc.expect, actual)
		})

	}

}
