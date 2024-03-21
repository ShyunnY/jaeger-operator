package infra

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/labels"
	"testing"
)

func TestCommonListOptions(t *testing.T) {

	options := commonListOptions()
	selector := labels.NewSelector().Add(options...)

	match := selector.Matches(labels.Set(map[string]string{
		"app":                          "jaeger",
		"app.kubernetes.io/part-of":    "jaeger",
		"app.kubernetes.io/managed-by": "jaeger-operator",
	}))

	assert.True(t, match)

}

func TestListOptions(t *testing.T) {

	cases := []struct {
		caseName     string
		component    string
		actualLabels map[string]string
		expect       bool
	}{
		{
			caseName:  "deployment list options",
			component: "deployment",
			actualLabels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
				"app.kubernetes.io/component":  "deployment",
			},
			expect: true,
		},
		{
			caseName:  "service list options",
			component: "service",
			actualLabels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
				"app.kubernetes.io/component":  "service",
			},
			expect: true,
		},
		{
			caseName:  "pod list options",
			component: "pod",
			actualLabels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
				"app.kubernetes.io/component":  "pod",
			},
			expect: true,
		},
		{
			caseName:  "service-account list options",
			component: "service-account",
			actualLabels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
				"app.kubernetes.io/component":  "service-account",
			},
			expect: true,
		},
		{
			caseName:  "httproute list options",
			component: "httproute",
			actualLabels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
				"app.kubernetes.io/component":  "httproute",
			},
			expect: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {

			options := ListOptions(tc.component)
			selector := labels.NewSelector().Add(options...)

			actual := selector.Matches(labels.Set(tc.actualLabels))
			assert.Equal(t, tc.expect, actual)
		})
	}

}
