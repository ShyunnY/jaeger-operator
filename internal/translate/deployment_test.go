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
		caseName             string
		render               StrategyRender
		expectDeploymentYaml string
	}{
		{
			caseName: "Render normal Deployment",
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
			expectDeploymentYaml: "all-in-one-normal",
		},
		{
			caseName: "Render normal Deployment with custom labels and annotations",
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
			expectDeploymentYaml: "all-in-one-with-custom-metadata",
		},
	}

	for _, tc := range cases {

		t.Run(tc.caseName, func(t *testing.T) {
			deploy, err := tc.render.Deployment()
			assert.NoError(t, err)

			expectDeploy, err := loadDeployment(tc.expectDeploymentYaml)
			assert.NoError(t, err)
			assert.Equal(t, expectDeploy, deploy)
		})

	}

}

func loadDeployment(deployYaml string) (*appsv1.Deployment, error) {
	file, err := os.ReadFile(fmt.Sprintf("testdata/out/deployment/%s.yaml", deployYaml))
	if err != nil {
		return nil, err
	}
	deploy := &appsv1.Deployment{}

	err = yaml.Unmarshal(file, deploy)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
