package kubernetes

import (
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type GenerationChanger struct {
	predicate.GenerationChangedPredicate
}

func (p GenerationChanger) Update(e event.UpdateEvent) bool {
	if !p.GenerationChangedPredicate.Update(e) {
		return false
	}

	oldJaeger := e.ObjectOld.(*jaegerv1a1.Jaeger)
	newJaeger := e.ObjectOld.(*jaegerv1a1.Jaeger)

	// We do not handle cases where the status changes
	return !cmp.Equal(oldJaeger.Status, newJaeger.Status)
}
