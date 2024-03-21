package translate

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/consts"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
)

func allInOneDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {

	// merge ports
	ports := append(getCollectorPort(true), getQueryPort()...)

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  consts.AllIneOneComponent,
		Image: "jaegertracing/all-in-one:1.55.0",
		Args:  instance.Spec.Components.AllInOne.Args,
		Env:   utils.ConvertEnvVar(instance.Spec.Components.AllInOne.Envs),
		Ports: ports,
	}

	deployName := ComponentName(NamespacedName(instance), consts.AllIneOneComponent)
	deploy := expectDeploySpec(deployName, instance, container)

	return deploy
}

func QueryDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {

	// merge ports
	ports := getQueryPort()

	// set default envs
	envs := utils.MergePodEnv(
		[]corev1.EnvVar{
			{
				Name:  "JAEGER_SERVICE_NAME",
				Value: NamespacedName(instance),
			},
			{
				Name:  "JAEGER_PROPAGATION",
				Value: "JAEGER_PROPAGATION",
			},
		},
		instance.Spec.Components.Query.Envs...,
	)

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  consts.QueryComponent,
		Image: "jaegertracing/jaeger-query:1.55.0",
		Args:  instance.Spec.Components.Query.Args,
		Env:   utils.ConvertEnvVar(envs),
		Ports: ports,
	}
	deployName := ComponentName(NamespacedName(instance), consts.QueryComponent)
	deploy := expectDeploySpec(deployName, instance, container)

	return deploy
}

func CollectorDeploy(instance *jaegerv1a1.Jaeger) *appsv1.Deployment {

	// merge ports
	ports := getCollectorPort(true)

	container := &corev1.Container{
		// TODO: add more settings?
		Name:  consts.CollectorComponent,
		Image: "jaegertracing/jaeger-collector:1.55.0",
		Args:  instance.Spec.Components.Collector.Args,
		Env:   utils.ConvertEnvVar(instance.Spec.Components.Collector.Envs),
		Ports: ports,
	}
	deployName := ComponentName(NamespacedName(instance), consts.CollectorComponent)
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
		// container.LivenessProbe = livenessProbe()
	}

	deployLabels := ComponentLabels("deployment", name, instance)
	podLabels := ComponentLabels("pod", name, instance)

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

// TODO: 不同的Jaeger组件有不同的admin端口, 我们应该根据外部传入端口进行构建
func livenessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/",
				Port: intstr.FromInt32(consts.AdminPort),
			},
		},
		InitialDelaySeconds: 5,
		PeriodSeconds:       15,
		FailureThreshold:    5,
	}
}
