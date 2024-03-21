package infra

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func ListOptions(component string) []labels.Requirement {

	componentReq, _ := labels.NewRequirement(
		"app.kubernetes.io/component",
		selection.Equals,
		[]string{component},
	)

	return append(commonListOptions(), *componentReq)
}

func commonListOptions() []labels.Requirement {

	ret := []labels.Requirement{}

	appReq, _ := labels.NewRequirement(
		"app",
		selection.DoubleEquals,
		[]string{"jaeger"},
	)
	ret = append(ret, *appReq)

	partReq, _ := labels.NewRequirement(
		"app.kubernetes.io/part-of",
		selection.DoubleEquals,
		[]string{"jaeger"},
	)
	ret = append(ret, *partReq)

	managerByReq, _ := labels.NewRequirement(
		"app.kubernetes.io/managed-by",
		selection.DoubleEquals,
		[]string{"jaeger-operator"},
	)
	ret = append(ret, *managerByReq)

	return ret
}
