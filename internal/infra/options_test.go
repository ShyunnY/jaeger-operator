package infra

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/labels"
	"testing"
)

func TestCommonListOptions(t *testing.T) {

	options := commonListOptions()
	selector := labels.NewSelector()

	selector = selector.Add(options...)
	match := selector.Matches(labels.Set(map[string]string{
		"app":                          "jaeger",
		"app.kubernetes.io/part-of":    "jaeger",
		"app.kubernetes.io/managed-by": "jaeger-operator",
	}))

	assert.True(t, match)

}

func TestListOptions(t *testing.T) {

}
