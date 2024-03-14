package utils

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
)

func MergePodArgs(args ...string) []string {

	// default args
	defaultArgs := []string{
		"--memory.max-traces=100000",
	}
	if len(args) == 0 {
		return defaultArgs
	}

	argsMap := argToMap(args)
	for _, arg := range defaultArgs {
		split := strings.Split(arg, "=")

		if _, ok := argsMap[split[0]]; !ok {
			argsMap[split[0]] = arg
		}
	}

	return mapToArgs(argsMap)
}

func argToMap(args []string) map[string]string {

	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		split := strings.Split(arg, "=")

		argsMap[split[0]] = arg
	}

	return argsMap
}

func mapToArgs(argsMap map[string]string) []string {

	args := make([]string, 0, len(argsMap))

	for _, arg := range argsMap {
		args = append(args, arg)
	}

	sort.Strings(args)
	return args
}

func MergePodEnv(envs ...jaegerv1a1.EnvSetting) []corev1.EnvVar {

	// default env
	defaultEnvs := []corev1.EnvVar{
		// By default, memory is used to store trace information
		{
			Name:  "SPAN_STORAGE_TYPE",
			Value: "memory",
		},
		{
			Name:  "COLLECTOR_ZIPKIN_HOST_PORT",
			Value: ":9411",
		},
		{
			Name:  "JAEGER_DISABLED",
			Value: "false",
		},
		{
			Name:  "COLLECTOR_OTLP_ENABLED",
			Value: "true",
		},
	}

	if len(envs) == 0 {
		return defaultEnvs
	}

	envsMap := envToMap(envs)
	for _, v := range defaultEnvs {
		if _, ok := envsMap[v.Name]; !ok {
			envsMap[v.Name] = v
		}
	}

	return mapToEnv(envsMap)
}

func envToMap(envs []jaegerv1a1.EnvSetting) map[string]corev1.EnvVar {
	retMap := make(map[string]corev1.EnvVar, len(envs))

	for _, v := range envs {
		env := corev1.EnvVar{
			Name:  v.Name,
			Value: v.Value,
		}
		retMap[env.Name] = env
	}

	return retMap
}

func mapToEnv(envsMap map[string]corev1.EnvVar) []corev1.EnvVar {

	envs := make([]corev1.EnvVar, 0, len(envsMap))
	for _, envVar := range envsMap {
		envs = append(envs, envVar)
	}

	return envs
}
