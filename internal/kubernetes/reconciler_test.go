package kubernetes

import (
	"context"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/jaeger"
	"github.com/ShyunnY/jaeger-operator/internal/logging"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestReconcile(t *testing.T) {
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
			Type: "allInOne",
		},
	}

	fakeClient := fake.NewClientBuilder().WithObjects(instance).WithScheme(jaeger.GetScheme()).Build()
	r := jaegerReconciler{
		name:       "jaeger-operator",
		namespaces: nil,
		logger:     logging.DefaultLogger(),
		client:     fakeClient,
		irMessage:  new(message.IRMessage),
	}

	r.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "all-in-one",
			Namespace: "default",
		},
	})
}
