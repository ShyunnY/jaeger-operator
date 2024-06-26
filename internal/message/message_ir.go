package message

import (
	"context"

	"github.com/telepresenceio/watchable"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gtwapi "sigs.k8s.io/gateway-api/apis/v1"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

type IRMessage struct {
	watchable.Map[string, *JaegerWithCtx]
}

type InfraIRMaps struct {
	watchable.Map[string, *InfraIR]
}

type InfraIR struct {
	Ctx            context.Context
	Strategy       string
	InstanceMedata metav1.ObjectMeta

	Deployment     []*appsv1.Deployment
	ConfigMap      *corev1.ConfigMap
	Service        []*corev1.Service
	ServiceAccount *corev1.ServiceAccount

	// extension resources
	HTTPRoutes []*gtwapi.HTTPRoute
}

func (ir *InfraIR) DeepCopy() *InfraIR {
	if ir == nil {
		return nil
	}
	infraIr := new(InfraIR)
	infraIr = ir

	return infraIr
}

func (ir *InfraIR) AddResources(obj any) {

	if obj == nil {
		return
	}

	switch resource := obj.(type) {
	case []*appsv1.Deployment:
		ir.Deployment = resource
	case *corev1.ConfigMap:
		ir.ConfigMap = resource
	case *corev1.ServiceAccount:
		ir.ServiceAccount = resource
	case []*corev1.Service:
		ir.Service = resource
	case []*gtwapi.HTTPRoute:
		ir.HTTPRoutes = resource
	default:
		panic("undefined resources")
	}

}

type StatusIRMaps struct {
	watchable.Map[types.NamespacedName, *jaegerv1a1.JaegerStatus]
}
