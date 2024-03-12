package utils

import (
	"crypto/sha256"
	"fmt"
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
