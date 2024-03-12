package message

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/telepresenceio/watchable"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type IRMessage struct {
	watchable.Map[string, *jaegerv1a1.Jaeger]
}

type InfraIRMaps struct {
	// key=jaeger name,val=ir
	watchable.Map[string, *InfraIR]
}

type InfraIR struct {
	InstanceName      string
	InstanceNamespace string
	Strategy          string

	Deployment     *appsv1.Deployment
	ConfigMap      *corev1.ConfigMap
	Service        []*corev1.Service
	ServiceAccount *corev1.ServiceAccount
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
	case *appsv1.Deployment:
		ir.Deployment = resource
	case *corev1.ConfigMap:
		ir.ConfigMap = resource
	case *corev1.ServiceAccount:
		ir.ServiceAccount = resource
	case []*corev1.Service:
		ir.Service = resource
	default:
		panic("undefined resources")
	}

}

type StatusIRMaps struct {
	watchable.Map[string, *jaegerv1a1.JaegerStatus]
}
