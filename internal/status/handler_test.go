package status

import (
	"context"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
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
		Status: jaegerv1a1.JaegerStatus{},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(jaeger.GetScheme()).WithObjects(instance).Build()
	uh := NewUpdateHandler(fakeClient, logging.DefaultLogger())

	instance.Status.Phase = "YYes"
	err := fakeClient.Update(context.TODO(), instance)
	assert.NoError(t, err)

	uh.apply(Update{
		NamespacedName: types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		},
		Object: new(jaegerv1a1.Jaeger),
		Mutator: func(oldObj client.Object) client.Object {
			obj := oldObj.(*jaegerv1a1.Jaeger)

			dp := obj.DeepCopy()
			dp.Status = jaegerv1a1.JaegerStatus{Phase: "YY语音es"}

			return dp
		},
	})

	assert.NotNil(t, instance)
}
