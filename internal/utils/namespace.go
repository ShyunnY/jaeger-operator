package utils

import (
	"k8s.io/utils/set"
	"strings"
)

// ExtractNamespace Extract namespaces separated by ',' in the string
func ExtractNamespace(ns string) set.Set[string] {

	nsSet := make(set.Set[string])

	if len(ns) == 0 {
		return nsSet
	}
	for _, singleNs := range strings.Split(ns, ",") {
		nsSet.Insert(singleNs)
	}

	return nsSet
}
