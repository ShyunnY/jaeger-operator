package translate

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

func TestAllInOneDeployment(t *testing.T) {

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
			caseName: "all-in-one-custom",
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
						CommonSpec: jaegerv1a1.CommonSpec{
							Deployment: jaegerv1a1.DeploymentSettings{
								Replicas: func() *int32 {
									var replicas int32 = 3
									return &replicas
								}(),
							},
							Metadata: jaegerv1a1.CommonMetadata{
								Labels: map[string]string{
									"label-1": "l-1",
									"label-2": "l-2",
								},
								Annotations: map[string]string{
									"annotation-1": "a-1",
									"annotation-2": "a-2",
								},
							},
						},
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
			deploy, err := tc.render.Deployment()
			assert.NoError(t, err)

			if true {
				outYaml, err := yaml.Marshal(deploy)
				assert.NoError(t, err)

				err = os.WriteFile(fmt.Sprintf("testdata/out/deployment/%s.yaml", tc.caseName), outYaml, 0644)
				assert.NoError(t, err)
				return
			}

			expectDeploy, err := loadDeployment(tc.caseName)
			assert.NoError(t, err)
			assert.Equal(t, len(expectDeploy), len(deploy))
		})

	}

}

// TODO: In order to compare the traversals in the right order, we need to add a sort method
func loadDeployment(deployYaml string) ([]*appsv1.Deployment, error) {
	file, err := os.ReadFile(fmt.Sprintf("testdata/out/deployment/%s.yaml", deployYaml))
	if err != nil {
		return nil, err
	}
	deploy := []*appsv1.Deployment{}

	err = yaml.Unmarshal(file, &deploy)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
