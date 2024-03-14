package status

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
)

func TestApply(t *testing.T) {

	instance := &jaegerv1a1.Jaeger{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tracing.orange.io/v1alpha1",
			Kind:       "Jaeger",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-in-one",
			Namespace: "default",
		},
		Spec: jaegerv1a1.JaegerSpec{
			Type: jaegerv1a1.AllInOneType,
		},
		Status: jaegerv1a1.JaegerStatus{
			Phase: "one",
			Conditions: []metav1.Condition{
				{
					Type:               "Success",
					Status:             metav1.ConditionTrue,
					ObservedGeneration: 123456789,
					LastTransitionTime: metav1.NewTime(time.Now()),
					Reason:             "Test-1",
					Message:            "normal",
				},
			},
		},
	}

	instance2 := &jaegerv1a1.Jaeger{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tracing.orange.io/v1alpha1",
			Kind:       "Jaeger",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-in-one",
			Namespace: "default",
		},
		Spec: jaegerv1a1.JaegerSpec{
			Type: jaegerv1a1.AllInOneType,
		},
		Status: jaegerv1a1.JaegerStatus{
			Phase: "tow",
			Conditions: []metav1.Condition{
				{
					Type:               "Success",
					Status:             metav1.ConditionTrue,
					ObservedGeneration: 123456789,
					LastTransitionTime: metav1.NewTime(time.Now()),
					Reason:             "Test-2",
					Message:            "error",
				},
			},
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(jaeger.GetScheme()).
		WithStatusSubresource(instance).WithObjects(instance).Build()
	uh := NewUpdateHandler(fakeClient, logging.DefaultLogger())
	uh.apply(Update{
		NamespacedName: types.NamespacedName{
			Namespace: instance2.Namespace,
			Name:      instance2.Name,
		},
		Object: new(jaegerv1a1.Jaeger),
		Mutator: func(oldObj client.Object) client.Object {
			obj := oldObj.(*jaegerv1a1.Jaeger)

			dp := obj.DeepCopy()
			dp.Status.Conditions = MergeCondition(dp.Status.Conditions, instance2.Status.Conditions...)
			dp.Status.Phase = instance2.Status.Phase
			return dp
		},
	})

	assert.NotNil(t, instance2)
}
