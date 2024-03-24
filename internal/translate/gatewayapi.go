package translate

import (
	"fmt"
	"github.com/ShyunnY/jaeger-operator/internal/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gtwapi "sigs.k8s.io/gateway-api/apis/v1"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
)

func processHTTPRoute(instance *jaegerv1a1.Jaeger, services []*corev1.Service) ([]*gtwapi.HTTPRoute, error) {

	var httpRoutes []*gtwapi.HTTPRoute

	for _, route := range instance.Spec.Extensions.HTTPRoute {

		var ref *gtwapi.ParentReference
		var backendRefs []gtwapi.HTTPBackendRef

		ref = route.ParentRef
		if ref.Port == nil && ref.SectionName == nil {
			return nil, fmt.Errorf("failed to create HTTPRoute, the sectionName or port must be specified")
		}

		r := gtwapi.ParentReference{
			Namespace:   ref.Namespace,
			Name:        ref.Name,
			SectionName: ref.SectionName,
			Port:        ref.Port,
		}

		// Set different HTTPRoute depending on the service target
		var service *corev1.Service
		if service = GetServiceByTarget(route.Target, services); service != nil {
			backendRefs = buildBackendRef(service, route.TargetPort)
		}

		httpRoute := &gtwapi.HTTPRoute{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "gateway.networking.k8s.io/v1",
				Kind:       "HTTPRoute",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      service.Name,
				Namespace: instance.Namespace,
				Labels:    utils.Labels(instance.Name, "httproute", string(instance.GetDeploymentType())),
				OwnerReferences: []metav1.OwnerReference{
					GetOwnerRef(instance),
				},
			},
			Spec: gtwapi.HTTPRouteSpec{
				CommonRouteSpec: gtwapi.CommonRouteSpec{
					ParentRefs: []gtwapi.ParentReference{
						r,
					},
				},
				Hostnames: route.Hostnames,
				Rules: []gtwapi.HTTPRouteRule{
					{
						BackendRefs: backendRefs,
					},
				},
			},
		}

		httpRoutes = append(httpRoutes, httpRoute)
	}

	return httpRoutes, nil
}

func GetServiceByTarget(target jaegerv1a1.JaegerServiceTarget, services []*corev1.Service) *corev1.Service {
	for _, service := range services {
		service := service

		if service.Labels == nil {
			return nil
		}

		if val, ok := service.Labels[jaegerv1a1.ServiceTargetLabelKey]; !ok {
			return nil
		} else if val == string(target) {
			return service
		}
	}

	return nil
}

func buildBackendRef(service *corev1.Service, port *int) []gtwapi.HTTPBackendRef {

	ns := gtwapi.Namespace(service.Namespace)
	portNumber := gtwapi.PortNumber(*port)

	if service.Spec.ClusterIP == "None" {
		return nil
	}

	backendRef := gtwapi.HTTPBackendRef{
		BackendRef: gtwapi.BackendRef{
			BackendObjectReference: gtwapi.BackendObjectReference{
				Name:      gtwapi.ObjectName(service.Name),
				Namespace: &ns,
				Port:      &portNumber,
			},
		},
	}

	return []gtwapi.HTTPBackendRef{backendRef}
}
