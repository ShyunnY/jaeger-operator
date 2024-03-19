/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gtwapi "sigs.k8s.io/gateway-api/apis/v1"
)

// DeploymentType Define the type of Jaeger deployment
// +kubebuilder:validation:Enum={allInOne,distribute}
type DeploymentType string

var (
	// AllInOneType The deployment type is allInOne
	AllInOneType DeploymentType = "allInOne"

	// Distribute The deployment type is distribute
	Distribute DeploymentType = "distribute"
)

// JaegerSpec Define the desired state of Jaeger
type JaegerSpec struct {

	// +kubebuilder:default=allInOne
	// +optional
	// Type Define the type of Jaeger deployment
	Type DeploymentType `json:"type,omitempty"`

	// Components Define the subComponents of Jaeger
	// +optional
	Components JaegerComponent `json:"components,omitempty"`

	// +optional
	CommonSpec CommonSpec `json:"commonSpec,omitempty"`
}

// JaegerStatus Define the observed state of Jaeger
type JaegerStatus struct {

	// +kubebuilder:default=Unknown
	// Phase Define the component phase of Jaeger
	Phase string `json:"phase"`

	// Conditions  Define the conditions of Jaeger
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STRATEGY",type="string",JSONPath=".spec.type",description="Jaeger deploy strategy"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.phase",description="Jaeger instance's status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Jaeger is the Schema for the jaegers API
type Jaeger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JaegerSpec   `json:"spec,omitempty"`
	Status JaegerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JaegerList contains a list of Jaeger
type JaegerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Jaeger `json:"items"`
}

// JaegerComponent Defining individual components
type JaegerComponent struct {

	// AllInOne Define all-in-one component
	// +optional
	AllInOne AllInOneComponent `json:"allInOne,omitempty"`

	// Collector Define collector component
	// +optional
	Collector CollectorComponent

	// QueryComponent Define collector component
	// +optional
	Query QueryComponent

	// Storage Define backend storage component
	// +optional
	Storage StorageComponent
}

type CollectorComponent struct {
	ComponentSettings `json:"setting,omitempty"`
}

type QueryComponent struct {
	ComponentSettings `json:"setting,omitempty"`
}

// StorageComponent Define storage component
type StorageComponent struct {

	// Type Define backend storage type
	Type StorageType

	// Type Define backend storage options
	Options []string
}

type StorageType string

var (
	MemoryStorageType StorageType = "memory"

	ElasticSarchStorage StorageType = "elasticsearch"
)

// CommonSpec Define Generic configuration of Kubernetes components
type CommonSpec struct {

	// Metadata Define metadata configuration of the component
	Metadata CommonMetadata `json:"metadata,omitempty"`

	// Service Define configuration of the kubernetes Services
	Service ServiceSettings `json:"service,omitempty"`

	// Deployment Define configuration of the kubernetes Deployments
	Deployment DeploymentSettings `json:"deployment"`

	// +optional
	// HTTPRoute Define the configuration of GatewayAPI routes
	HTTPRoute []HTTPRoute `json:"httpRoutes,omitempty"`
}

type DeploymentSettings struct {
	Replicas *int32
}

// HTTPRoute Define the HTTPRoute configuration for the Gateway API
type HTTPRoute struct {

	// +kubebuilder:validation:Required
	// Target Define the Service target to which routes need to be added
	Target JaegerServiceTarget `json:"target,omitempty"`

	// +kubebuilder:validation:Required
	// TargetPort Define the Service target port to which routes need to be added
	TargetPort *int `json:"targetPort,omitempty"`

	// +optional
	// Hostnames Define a set of hostnames that should match against the HTTP Host
	// header to select a HTTPRoute used to process the request.
	Hostnames []gtwapi.Hostname `json:"hostnames,omitempty"`

	// +kubebuilder:validation:Required
	// ParentRefs references the resources (usually Gateways) that a Route wants
	// to be attached to.
	ParentRef *gtwapi.ParentReference `json:"parentRef,omitempty"`
}

// JaegerServiceTarget Define target of Jaeger Kubernetes Services target
// +kubebuilder:validation:Enum={collector,query}
type JaegerServiceTarget string

var (
	CollectorServiceTarget JaegerServiceTarget = "collector"

	QueryServiceTarget JaegerServiceTarget = "query"
)

// ServiceType Define type of components Services
// +kubebuilder:validation:Enum={ClusterIP,NodePort,LoadBalancer}
type ServiceType string

var (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the cluster IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"

	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
)

// ServiceSettings Define personalized configuration of Jaeger component Service
type ServiceSettings struct {

	// Service Type string describes ingress methods for a service
	// +kubebuilder:default=ClusterIP
	Type ServiceType `json:"type,omitempty"`
}

// CommonMetadata Define Metadata common to all components and infrastructure resources
type CommonMetadata struct {

	// Labels Define labels setting for metadata on the resource
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations Define annotations setting for metadata on the resource
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AllInOneComponent Define the AllInOne Jaeger component under the type=allInOne
type AllInOneComponent struct {
	ComponentSettings `json:"setting,omitempty"`
}

// ComponentSettings Define common Settings between components
type ComponentSettings struct {

	// Args Defined cmd args for Jaeger components
	// +optional
	Args []string `json:"args,omitempty"`

	// Envs Defined env for Jaeger components
	// +optional
	Envs []EnvSetting `json:"envs,omitempty"`
}

type EnvSetting struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// Define Env name
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// Define Env value
	Value string `json:"value,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Jaeger{}, &JaegerList{})
}
