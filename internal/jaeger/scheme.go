package jaeger

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := jaegerv1a1.AddToScheme(scheme); err != nil {
		panic(err)
	}
}

func GetScheme() *runtime.Scheme {
	return scheme
}
