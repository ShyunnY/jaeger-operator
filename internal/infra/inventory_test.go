package infra

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestComputeServiceAccount(t *testing.T) {

	exist := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-in-one-default-dec10161",
			Namespace: "default",
			Labels: map[string]string{
				"app":                                 "jaeger",
				"app.kubernetes.io/component":         "service-account",
				"app.kubernetes.io/managed-by":        "jaeger-operator",
				"app.kubernetes.io/name":              "all-in-one",
				"app.kubernetes.io/part-of":           "jaeger",
				"jaegertracing.orange.io/operated-by": "jaeger-operator",
				"tracing.orange.io/strategy":          "allInOne",
			},
		},
	}

	desire := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-in-one-default-dec1016123",
			Namespace: "default",
			Labels: map[string]string{
				"app":                                 "jaeger",
				"app.kubernetes.io/component":         "service-account",
				"app.kubernetes.io/managed-by":        "jaeger-operator",
				"app.kubernetes.io/name":              "all-in-one",
				"app.kubernetes.io/part-of":           "jaeger",
				"jaegertracing.orange.io/operated-by": "jaeger-operator",
				"tracing.orange.io/strategy":          "allInOne",
			},
		},
	}

	fakeClient := fake.NewFakeClient(exist)
	it := InventoryComputer{
		cli: fakeClient,
	}

	_, err := it.ComputeServiceAccount(context.TODO(), desire)
	assert.NoError(t, err)
}
