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
		caseName string
		render   StrategyRender
	}{
		{
			caseName: "service-account-normal",
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
		},
		{
			caseName: "service-account-difference-strategy",
			render: &AllInOneRender{
				instance: &jaegerv1a1.Jaeger{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "tracing.orange.io/v1alpha1",
						Kind:       "Jaeger",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "distribute",
						Namespace: "default",
						UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
					},
					Spec: jaegerv1a1.JaegerSpec{
						Type: jaegerv1a1.Distribute,
					},
				},
			},
		},
		{
			caseName: "service-account-custom",
			render: &AllInOneRender{
				instance: &jaegerv1a1.Jaeger{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "tracing.orange.io/v1alpha1",
						Kind:       "Jaeger",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "prod",
						Namespace: "default",
						UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
					},
					Spec: jaegerv1a1.JaegerSpec{
						Type: jaegerv1a1.Distribute,
						CommonSpec: jaegerv1a1.CommonSpec{
							Metadata: jaegerv1a1.CommonMetadata{
								Labels: map[string]string{
									"label-1": "l-1",
								},
								Annotations: map[string]string{
									"annotation-1": "a-1",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			sa, err := tc.render.ServiceAccount()
			assert.NoError(t, err)

			if false {
				outYaml, err := yaml.Marshal(sa)
				assert.NoError(t, err)

				err = os.WriteFile(fmt.Sprintf("testdata/out/serviceaccount/%s.yaml", tc.caseName), outYaml, 0644)
				assert.NoError(t, err)
				return
			}

			expectServiceAccount, err := loadServiceAccount(tc.caseName)
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
