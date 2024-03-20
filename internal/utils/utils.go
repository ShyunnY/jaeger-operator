package utils

import (
	"crypto/sha256"
	"fmt"
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

func GetHashedName(nsName types.NamespacedName) string {

	nsString := nsName.String()
	prefix := HashString(nsString)

	nsString = strings.ReplaceAll(nsString, "/", "-")
	resourceName := fmt.Sprintf("%s-%s", nsString, prefix[0:8])

	return resourceName
}

func HashString(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return strings.ToLower(fmt.Sprintf("%x", h.Sum(nil)))
}

// Labels Returns the generic labels
func Labels(name, component, strategy string) map[string]string {
	return map[string]string{
		"app":                          "jaeger",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/part-of":    "jaeger",
		"app.kubernetes.io/managed-by": "jaeger-operator",
		"tracing.orange.io/strategy":   strategy,
	}
}

// MergeCommonMap Combining the data of two maps, the append map overwrites the value of the exist map
func MergeCommonMap(exist map[string]string, append map[string]string) map[string]string {

	if append == nil || len(append) == 0 {
		return exist
	}

	retMap := make(map[string]string)
	for existK, existV := range exist {
		if val, ok := append[existK]; !ok {
			retMap[existK] = existV
		} else {
			retMap[existK] = val
			delete(append, existK)
		}
	}

	for k, v := range append {
		retMap[k] = v
	}

	return retMap
}

func AppendEnvs(envs []corev1.EnvVar, instance *jaegerv1a1.Jaeger) {

	if len(envs) == 0 {
		return
	}

	adds := make([]jaegerv1a1.EnvSetting, 0, len(envs))
	for _, env := range envs {

		envSetting := jaegerv1a1.EnvSetting{
			Name:  env.Name,
			Value: env.Value,
		}

		adds = append(adds, envSetting)
	}

	instance.Spec.Components.Collector.Envs = append(adds, instance.Spec.Components.Collector.Envs...)
}
