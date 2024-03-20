package utils

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
)

func MergePodArgs(existArgs []string, args ...string) []string {

	if len(args) == 0 {
		return existArgs
	}

	argsMap := argToMap(args)
	for _, arg := range existArgs {
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

func MergePodEnv(existEnvs []corev1.EnvVar, envs ...jaegerv1a1.EnvSetting) []jaegerv1a1.EnvSetting {

	if len(envs) == 0 {
		return ConvertEnvSettings(existEnvs)
	}

	envsMap := envToMap(envs)
	for _, v := range existEnvs {
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

func mapToEnv(envsMap map[string]corev1.EnvVar) []jaegerv1a1.EnvSetting {

	envs := make([]corev1.EnvVar, 0, len(envsMap))
	for _, envVar := range envsMap {
		envs = append(envs, envVar)
	}

	return ConvertEnvSettings(envs)
}

func ConvertEnvVar(envs []jaegerv1a1.EnvSetting) []corev1.EnvVar {

	if len(envs) == 0 {
		return nil
	}

	ret := make([]corev1.EnvVar, 0, len(envs))
	for _, env := range envs {
		envVar := corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		}

		ret = append(ret, envVar)
	}

	return ret
}

func ConvertEnvSettings(envs []corev1.EnvVar) []jaegerv1a1.EnvSetting {

	if len(envs) == 0 {
		return nil
	}

	ret := make([]jaegerv1a1.EnvSetting, 0, len(envs))
	for _, env := range envs {
		envVar := jaegerv1a1.EnvSetting{
			Name:  env.Name,
			Value: env.Value,
		}

		ret = append(ret, envVar)
	}

	return ret
}
