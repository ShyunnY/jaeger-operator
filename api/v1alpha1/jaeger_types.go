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

// JaegerSpec defines the desired state of Jaeger
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

// JaegerStatus defines the observed state of Jaeger
type JaegerStatus struct {

	// Conditions  Define the conditions of Jaeger
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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

type JaegerComponent struct {
}

type CommonSpec struct {
	Metadata CommonMetadata `json:"metadata,omitempty"`
}

type CommonMetadata struct {
	Labels map[string]string `json:"labels,omitempty"`

	Annotations map[string]string `json:"annotations,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Jaeger{}, &JaegerList{})
}
