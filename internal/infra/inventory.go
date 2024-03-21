package infra

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gtwapi "sigs.k8s.io/gateway-api/apis/v1"
)

var _ Inventory = (*InventoryComputer)(nil)

type Inventory interface {
	ComputeServiceAccount(ctx context.Context, desire *corev1.ServiceAccount) (*InventoryObject, error)
	ComputeService(ctx context.Context, desire []*corev1.Service) (*InventoryObject, error)
	ComputeDeployment(ctx context.Context, desire []*appsv1.Deployment) (*InventoryObject, error)
	ComputeHTTPRoutes(ctx context.Context, desire []*gtwapi.HTTPRoute) (*InventoryObject, error)
}

type InventoryObject struct {
	CreateObjects []client.Object
	DeleteObjects []client.Object
	UpdateObjects []client.Object
}

type Deployment struct {
}

// InventoryComputer TODO: 资源热更新问题:  apply之后, 资源并没有按照预期的更改
type InventoryComputer struct {
	instanceNamespace string
	instanceName      string

	cli client.Client
}

func (ic *InventoryComputer) ComputeHTTPRoutes(ctx context.Context, desires []*gtwapi.HTTPRoute) (*InventoryObject, error) {

	// Lists the HTTPRoutes managed by the current operator and whose instance is the current Jaeger
	list := &gtwapi.HTTPRouteList{}
	selector := labels.NewSelector().Add(ListOptions("httproute")...)
	if err := ic.cli.List(
		ctx,
		list,
		client.InNamespace(ic.instanceNamespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	); err != nil {
		return nil, err
	}

	updates := []client.Object{}
	mcreate := make(map[string]*gtwapi.HTTPRoute, len(desires))
	mdelete := make(map[string]*gtwapi.HTTPRoute, len(desires))
	for i := range desires {
		desire := desires[i]
		mcreate[toNsName(desire)] = desire
		mdelete[toNsName(desire)] = desire
	}

	if len(list.Items) == 0 {
		clear(mdelete)
	} else {
		for i := range list.Items {
			exist := list.Items[i]
			if desire, ok := mcreate[toNsName(&exist)]; ok {
				dp := exist.DeepCopy()

				dp.SetLabels(map[string]string{})
				for k, v := range desire.Labels {
					dp.Labels[k] = v
				}
				dp.SetAnnotations(map[string]string{})
				for k, v := range desire.Annotations {
					dp.Annotations[k] = v
				}
				dp.OwnerReferences = desire.OwnerReferences

				dp.Spec = *desire.Spec.DeepCopy()

				updates = append(updates, dp)
				delete(mcreate, toNsName(&exist))
				delete(mdelete, toNsName(&exist))
			}
		}
	}

	createObjs := []client.Object{}
	deleteObjs := []client.Object{}
	for _, service := range mcreate {
		createObjs = append(createObjs, service)
	}
	for _, service := range mdelete {
		deleteObjs = append(deleteObjs, service)
	}
	return &InventoryObject{
		CreateObjects: createObjs,
		UpdateObjects: updates,
		DeleteObjects: deleteObjs,
	}, nil
}

func (ic *InventoryComputer) ComputeService(ctx context.Context, desires []*corev1.Service) (*InventoryObject, error) {

	// Lists the Services managed by the current operator and whose instance is the current Jaeger
	list := &corev1.ServiceList{}
	selector := labels.NewSelector().Add(ListOptions("service")...)
	if err := ic.cli.List(
		ctx,
		list,
		client.InNamespace(ic.instanceNamespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	); err != nil {
		return nil, err
	}

	updates := []client.Object{}
	mcreate := make(map[string]*corev1.Service, len(desires))
	mdelete := make(map[string]*corev1.Service, len(desires))
	for i := range desires {
		desire := desires[i]
		mcreate[toNsName(desire)] = desire
		mdelete[toNsName(desire)] = desire
	}

	if len(list.Items) == 0 {
		clear(mdelete)
	} else {
		for i := range list.Items {
			exist := list.Items[i]

			if desire, ok := mcreate[toNsName(&exist)]; ok {
				dp := exist.DeepCopy()

				dp.SetLabels(map[string]string{})
				for k, v := range desire.Labels {
					dp.Labels[k] = v
				}
				dp.SetAnnotations(map[string]string{})
				for k, v := range desire.Annotations {
					dp.Annotations[k] = v
				}
				dp.OwnerReferences = desire.OwnerReferences

				dp.Spec = desire.Spec
				// We assign existing Service ip to Service ip that need to be updated
				if len(exist.Spec.ClusterIP) != 0 && len(desire.Spec.ClusterIP) == 0 {
					dp.Spec.ClusterIP = exist.Spec.ClusterIP
				}

				updates = append(updates, dp)
				delete(mcreate, toNsName(&exist))
				delete(mdelete, toNsName(&exist))
			}
		}
	}

	createObjs := []client.Object{}
	deleteObjs := []client.Object{}
	for _, service := range mcreate {
		createObjs = append(createObjs, service)
	}
	for _, service := range mdelete {
		deleteObjs = append(deleteObjs, service)
	}
	return &InventoryObject{
		CreateObjects: createObjs,
		UpdateObjects: updates,
		DeleteObjects: deleteObjs,
	}, nil
}

func (ic *InventoryComputer) ComputeServiceAccount(ctx context.Context, desire *corev1.ServiceAccount) (*InventoryObject, error) {

	// Lists the ServiceAccount managed by the current operator and whose instance is the current Jaeger
	list := &corev1.ServiceAccountList{}
	selector := labels.NewSelector().Add(ListOptions("service-account")...)
	if err := ic.cli.List(
		ctx,
		list,
		client.InNamespace(ic.instanceNamespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	); err != nil {

		return nil, err
	}

	createObjs := []client.Object{desire}
	updateObjs := []client.Object{}
	deleteObjs := []client.Object{desire}

	if len(list.Items) == 0 {
		deleteObjs = nil
	} else {
		for i := range list.Items {
			item := list.Items[i]

			// Check the Namespaced names of both
			if toNsName(&item) == toNsName(desire) {
				dp := item.DeepCopy()
				dp.SetLabels(map[string]string{})
				dp.SetAnnotations(map[string]string{})

				dp.OwnerReferences = desire.OwnerReferences
				for k, v := range desire.Labels {
					dp.Labels[k] = v
				}
				for k, v := range desire.Annotations {
					dp.Annotations[k] = v
				}

				updateObjs = append(updateObjs, dp)
				createObjs = nil
				updateObjs = nil
			}
		}
	}

	return &InventoryObject{
		CreateObjects: createObjs,
		UpdateObjects: updateObjs,
		DeleteObjects: deleteObjs,
	}, nil
}

func (ic *InventoryComputer) ComputeDeployment(ctx context.Context, desires []*appsv1.Deployment) (*InventoryObject, error) {

	// Lists the Deployment managed by the current operator and whose instance is the current Jaeger
	list := &appsv1.DeploymentList{}
	selector := labels.NewSelector().Add(ListOptions("deployment")...)
	if err := ic.cli.List(
		ctx,
		list,
		client.InNamespace(ic.instanceNamespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	); err != nil {

		return nil, err
	}

	updates := []client.Object{}
	mcreate := make(map[string]*appsv1.Deployment, len(desires))
	mdelete := make(map[string]*appsv1.Deployment, len(desires))
	for i := range desires {
		desire := desires[i]
		mcreate[toNsName(desire)] = desire
		mdelete[toNsName(desire)] = desire
	}

	if len(list.Items) == 0 {
		clear(mdelete)
	} else {
		for i := range list.Items {
			exist := list.Items[i]

			if desire, ok := mcreate[toNsName(&exist)]; ok {
				dp := desire.DeepCopy()
				dp.SetLabels(map[string]string{})
				dp.SetAnnotations(map[string]string{})

				// We don't overwrite with replicas, maybe replicas are written by the HPA
				dp.Spec = desire.Spec
				dp.Spec.Selector = desire.Spec.Selector
				dp.OwnerReferences = desire.OwnerReferences
				for k, v := range desire.Labels {
					dp.Labels[k] = v
				}
				for k, v := range desire.Annotations {
					dp.Annotations[k] = v
				}

				updates = append(updates, dp)
				delete(mcreate, toNsName(&exist))
				delete(mdelete, toNsName(&exist))
			}
		}
	}

	createObjs := []client.Object{}
	deleteObjs := []client.Object{}
	for _, service := range mcreate {
		createObjs = append(createObjs, service)
	}
	for _, service := range mdelete {
		deleteObjs = append(deleteObjs, service)
	}
	return &InventoryObject{
		CreateObjects: createObjs,
		UpdateObjects: updates,
		DeleteObjects: deleteObjs,
	}, nil

}

func toNsName(obj client.Object) string {
	return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
}
