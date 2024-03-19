package translate

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func allInOneDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {

	// merge Envs and Args
	envs := utils.MergePodEnv(instance.Spec.Components.AllInOne.Envs...)
	args := utils.MergePodArgs(instance.Spec.Components.AllInOne.Args...)
	// merge ports
	ports := append(getCollectorPort(true), getQueryPort()...)

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  allIneOneComponent,
		Image: "jaegertracing/all-in-one:1.55.0",
		Args:  args,
		Env:   envs,
		Ports: ports,
	}

	deployName := ComponentName(NamespacedName(instance), allIneOneComponent)
	deploy := expectDeploySpec(deployName, instance, container)

	return deploy
}

func queryDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {
	// merge Envs and Args
	args := utils.MergePodArgs(instance.Spec.Components.Query.Args...)
	envs := utils.MergePodEnv(instance.Spec.Components.Query.Envs...)
	// merge ports
	ports := getQueryPort()

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  queryComponent,
		Image: "jaegertracing/jaeger-query:1.55.0",
		Args:  args,
		Env:   envs,
		Ports: ports,
	}
	deployName := ComponentName(NamespacedName(instance), queryComponent)
	deploy := expectDeploySpec(deployName, instance, container)

	return deploy
}

func collectorDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {

	// merge Envs and Args
	args := utils.MergePodArgs(instance.Spec.Components.Collector.Args...)
	envs := utils.MergePodEnv(instance.Spec.Components.Collector.Envs...)
	// merge ports
	ports := getCollectorPort(true)

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  collectorComponent,
		Image: "jaegertracing/jaeger-collector:1.55.0",
		Args:  args,
		Env:   envs,
		Ports: ports,
	}
	deployName := ComponentName(NamespacedName(instance), collectorComponent)
	deploy := expectDeploySpec(deployName, instance, container)

	return deploy
}

func expectDeploySpec(name string, instance *jaegerv1a1.Jaeger, container *corev1.Container) *appsv1.Deployment {

	var replicas int32
	if instance.Spec.CommonSpec.Deployment.Replicas == nil ||
		*instance.Spec.CommonSpec.Deployment.Replicas == 0 {
		replicas = 1
	} else {
		replicas = *instance.Spec.CommonSpec.Deployment.Replicas
	}

	if container.LivenessProbe == nil {
		container.LivenessProbe = livenessProbe()
	}

	deployLabels := ComponentLabels("deployment", instance)
	podLabels := ComponentLabels("pod", instance)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   instance.Namespace,
			Labels:      deployLabels,
			Annotations: instance.GetCommonSpecAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(instance),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: instance.GetCommonSpecAnnotations(),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*container,
					},
					ServiceAccountName: NamespacedName(instance),
				},
			},
		},
	}
}

func livenessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/",
				Port: intstr.FromInt32(adminPort),
			},
		},
		InitialDelaySeconds: 5,
		PeriodSeconds:       15,
		FailureThreshold:    5,
	}
}
