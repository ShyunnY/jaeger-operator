package translate

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

func TestAllInOneServiceAccount(t *testing.T) {

	cases := []struct {
		caseName                 string
		render                   StrategyRender
		expectServiceAccountYaml string
	}{
		{
			caseName: "Render normal Service Account",
			render: &AllInOneRender{
				instance: &jaegerv1a1.Jaeger{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "tracing.orange.io/v1alpha1",
						Kind:       "Jaeger",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "all-in-one",
						Namespace: "default",
						UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
					},
					Spec: jaegerv1a1.JaegerSpec{
						Type: jaegerv1a1.AllInOneType,
					},
				},
			},
			expectServiceAccountYaml: "all-in-one-normal",
		},
		{
			caseName: "Render normal Service Account with custom labels and annotations",
			render: &AllInOneRender{
				instance: &jaegerv1a1.Jaeger{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "tracing.orange.io/v1alpha1",
						Kind:       "Jaeger",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "all-in-one",
						Namespace: "default",
						UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
					},
					Spec: jaegerv1a1.JaegerSpec{
						Type: jaegerv1a1.AllInOneType,
					},
				},
				labels: map[string]string{
					"custom-labels1": "label1",
					"custom-labels2": "label2",
				},
				annotations: map[string]string{
					"custom-annotation1": "annotation1",
					"custom-annotation2": "annotation2",
				},
			},
			expectServiceAccountYaml: "all-in-one-with-custom-metadata",
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			sa, err := tc.render.ServiceAccount()
			assert.NoError(t, err)

			expectServiceAccount, err := loadServiceAccount(tc.expectServiceAccountYaml)
			assert.NoError(t, err)
			assert.Equal(t, expectServiceAccount, sa)
		})

	}

}

func loadServiceAccount(serviceAccountYaml string) (*corev1.ServiceAccount, error) {
	file, err := os.ReadFile(fmt.Sprintf("testdata/out/serviceaccount/%s.yaml", serviceAccountYaml))
	if err != nil {
		return nil, err
	}
	sa := &corev1.ServiceAccount{}

	err = yaml.Unmarshal(file, sa)
	if err != nil {
		return nil, err
	}

	return sa, nil
}
