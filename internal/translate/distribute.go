package translate

import (
	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	"github.com/ShyunnY/jaeger-operator/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DistributeRender struct {
	instance *jaegerv1a1.Jaeger

	labels      map[string]string
	annotations map[string]string
}

func (d *DistributeRender) Deployment() (*appsv1.Deployment, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DistributeRender) ServiceAccount() (*corev1.ServiceAccount, error) {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      NamespacedName(d.instance),
			Namespace: d.instance.Namespace,
			Labels: utils.MergeCommonMap(utils.Labels(d.instance.Name, "service-account",
				string(d.GetStrategy())), d.labels),
			Annotations: d.annotations,
			OwnerReferences: []metav1.OwnerReference{
				GetOwnerRef(d.instance),
			},
		},
	}, nil
}

func (d *DistributeRender) ConfigMap() (*corev1.ConfigMap, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DistributeRender) Service() ([]*corev1.Service, error) {

	services := []*corev1.Service{}
	selector := utils.MergeCommonMap(utils.Labels(d.instance.Name, "pod", string(d.GetStrategy())), d.labels)

	queryService := QueryService(d.instance)
	queryService.Spec.Selector = selector
	queryService.Labels = utils.MergeCommonMap(queryService.Labels, utils.MergeCommonMap(
		utils.Labels(d.instance.Name, "service", string(d.GetStrategy())), d.labels))
	queryService.Annotations = d.annotations
	services = append(services, queryService)

	collectorServices := CollectorServices(d.instance)
	for i := range collectorServices {
		collectorServices[i].Spec.Selector = selector
		collectorServices[i].Labels = utils.MergeCommonMap(collectorServices[i].Labels,
			utils.MergeCommonMap(utils.Labels(d.instance.Name, "service", string(d.GetStrategy())), d.labels))
		collectorServices[i].Annotations = d.annotations
	}
	services = append(services, collectorServices...)

	return services, nil
}

func (d *DistributeRender) GetStrategy() jaegerv1a1.DeploymentType {
	//TODO implement me
	panic("implement me")
}
