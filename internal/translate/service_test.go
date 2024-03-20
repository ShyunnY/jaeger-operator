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

func TestServices(t *testing.T) {

	instance := &jaegerv1a1.Jaeger{
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
	}

	cases := []struct {
		caseName      string
		actualService *corev1.Service
	}{
		{
			caseName:      "query-service",
			actualService: QueryService(instance),
		},
		{
			caseName:      "collect-service",
			actualService: CollectorServices(instance)[0], // We establish that the first one is always cluster svc
		},
		{
			caseName:      "collect-headless-service",
			actualService: CollectorServices(instance)[1], // We establish that the second one is always headless svc
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			expected, err := loadService(tc.caseName)
			assert.NoError(t, err)

			if true {
				outYaml, err := yaml.Marshal(tc.actualService)
				assert.NoError(t, err)

				err = os.WriteFile(fmt.Sprintf("testdata/out/service/%s.yaml", tc.caseName), outYaml, 0644)
				assert.NoError(t, err)
				return
			}

			assert.Equal(t, expected, tc.actualService)
		})
	}

}

func TestAllInOneService(t *testing.T) {
	cases := []struct {
		caseName string
		render   StrategyRender
	}{
		{
			caseName: "all-in-one-normal",
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
			caseName: "distribute-normal",
			render: &DistributeRender{
				instance: &jaegerv1a1.Jaeger{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "tracing.orange.io/v1alpha1",
						Kind:       "Jaeger",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "production",
						Namespace: "default",
						UID:       types.UID("a98d5c73-8656-4035-be2f-0930f58bc89d"),
					},
					Spec: jaegerv1a1.JaegerSpec{
						Type: jaegerv1a1.Distribute,
						Components: jaegerv1a1.JaegerComponent{
							Storage: jaegerv1a1.StorageComponent{
								Type: jaegerv1a1.ElasticSearchStorage,
								Es: jaegerv1a1.EsStorage{
									URL: "127.0.0.1:9200",
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
			services, err := tc.render.Service()
			assert.NoError(t, err)

			if true {
				outYaml, err := yaml.Marshal(services)
				assert.NoError(t, err)

				err = os.WriteFile(fmt.Sprintf("testdata/out/service/%s.yaml", tc.caseName), outYaml, 0644)
				assert.NoError(t, err)
				return
			}

			expectServices, err := loadService(tc.caseName)
			assert.NoError(t, err)
			assert.Equal(t, len(expectServices), len(services))
		})

	}
}

// TODO: In order to compare the traversals in the right order, we need to add a sort method
func loadService(caseName string) ([]*corev1.Service, error) {
	file, err := os.ReadFile(fmt.Sprintf("testdata/out/service/%s.yaml", caseName))
	if err != nil {
		return nil, err
	}
	svc := []*corev1.Service{}

	err = yaml.Unmarshal(file, &svc)
	if err != nil {
		return nil, err
	}

	return svc, nil
}
