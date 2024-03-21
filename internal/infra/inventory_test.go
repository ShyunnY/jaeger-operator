package infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/translate"
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

func TestComputeService(t *testing.T) {

	instance := &jaegerv1a1.Jaeger{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tracing.orange.io/v1alpha1",
			Kind:       "Jaeger",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prod",
			Namespace: "default",
			UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
		},
		Spec: jaegerv1a1.JaegerSpec{
			Type: jaegerv1a1.AllInOneType,
			CommonSpec: jaegerv1a1.CommonSpec{
				Metadata: jaegerv1a1.CommonMetadata{
					Labels: map[string]string{
						"label-1": "l-1",
					},
					Annotations: map[string]string{
						"annotation-1": "a-1",
					},
				},
			},
		},
	}

	querySvc := translate.QueryService(instance)

	fakeCli := fake.NewClientBuilder().WithObjects(querySvc).Build()
	ic := InventoryComputer{
		cli: fakeCli,
	}

	deepCopy := querySvc.DeepCopy()
	deepCopy.Spec.Type = corev1.ServiceTypeNodePort
	service, err := ic.ComputeService(context.TODO(), []*corev1.Service{deepCopy})
	assert.NoError(t, err)
	assert.NotEmpty(t, service)

}
