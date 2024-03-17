package translate

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gtwapi "sigs.k8s.io/gateway-api/apis/v1"
	"testing"
)

func TestProcessHTTPRoute(t *testing.T) {

	sectionName := gtwapi.SectionName("section-1")
	port := gtwapi.PortNumber(8080)

	cases := []struct {
		caseName     string
		instance     *jaegerv1a1.Jaeger
		servicesFunc func(jaeger *jaegerv1a1.Jaeger) []*corev1.Service
	}{
		{
			caseName: "query HTTPRoute target",
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
						HTTPRoute: []jaegerv1a1.HTTPRoute{
							{
								Target: jaegerv1a1.QueryServiceTarget,
								TargetPort: func() *int {
									port := 16686
									return &port
								}(),
								ParentRef: &gtwapi.ParentReference{
									Name:        "eg-gateway",
									SectionName: &sectionName,
									Port:        &port,
								},
							},
						},
					},
				},
			},
			servicesFunc: func(instance *jaegerv1a1.Jaeger) []*corev1.Service {
				return []*corev1.Service{
					QueryService(instance),
				}
			},
		},
		{
			caseName: "collector HTTPRoute target",
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
						HTTPRoute: []jaegerv1a1.HTTPRoute{
							{
								Target: jaegerv1a1.QueryServiceTarget,
								TargetPort: func() *int {
									port := 14268
									return &port
								}(),
								ParentRef: &gtwapi.ParentReference{
									Name:        "eg-gateway",
									SectionName: &sectionName,
									Port:        &port,
								},
							},
							{
								Target: jaegerv1a1.QueryServiceTarget,
								TargetPort: func() *int {
									port := 16686
									return &port
								}(),
								ParentRef: &gtwapi.ParentReference{
									Name:        "eg-gateway",
									SectionName: &sectionName,
									Port:        &port,
								},
							},
						},
					},
				},
			},
			servicesFunc: func(instance *jaegerv1a1.Jaeger) []*corev1.Service {
				return []*corev1.Service{
					QueryService(instance),
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {

			services := tc.servicesFunc(tc.instance)
			assert.NotZero(t, len(services))

			route, err := processHTTPRoute(tc.instance, services)
			assert.NoError(t, err)
			assert.NotZero(t, len(route))

		})
	}

}
