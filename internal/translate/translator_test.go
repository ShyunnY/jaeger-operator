package translate

import (
	"github.com/stretchr/testify/assert"
	"testing"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/message"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTranslateNormal(t *testing.T) {

	irMaps := new(message.InfraIRMaps)
	translator := Translator{
		InfraIRMap: irMaps,
	}
	instance := &jaegerv1a1.Jaeger{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tracing.orange.io/v1alpha1",
			Kind:       "Jaeger",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-in-one",
			Namespace: "default",
		},
		Spec: jaegerv1a1.JaegerSpec{
			Type: jaegerv1a1.AllInOneType,
		},
	}

	err := translator.Translate(instance)
	assert.NoError(t, err)

	value, exist := irMaps.Load(instance.Name)
	assert.True(t, exist)
	assert.NotEmpty(t, value)
}
