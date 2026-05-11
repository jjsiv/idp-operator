/*
Copyright 2026.

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

// ResourceSpec defines the desired state of Resource
type ResourceSpec struct {
	// +optional
	Metadata `json:",inline"`

	// Reference to a Tenant that will be set as the owner of this Resource.
	// +required
	OwnerRef OwnerReference `json:"ownerRef"`

	// Reference to a ResourceTemplate that this Resource will be provisioned from.
	// +required
	TemplateRef TemplateReference `json:"templateRef"`

	// Specification for Backstage Software Catalog Resource to create from this Resource.
	// Resource owner will be set to the Tenant Backstage Group.
	// Type will be set to ResourceTemplate's name.
	// Parameters passed to the template as well as other values will be exposed as annotations.
	Backstage BackstageResourceSpec `json:"backstage,omitempty"`
}

// Reference to an owner resource.
type OwnerReference struct {
	// Reference to owner's Kind.
	// Only "Tenant" is supported.
	// +required
	Kind string `json:"kind"`

	// Name of the Tenant resource in the same namespace.
	// // +required
	Name string `json:"name"`
}

// TemplateReference defines a reference and parameters for a ResourceTemplate.
type TemplateReference struct {
	// Name of the ResourceTemplate.
	// +required
	Name string `json:"name"`

	// Namespace of the ResourceTemplate. If not specified, it is assumed to be in the same namespace as the Resource.
	Namespace string `json:"namespace,omitempty"`

	// Parameters to pass to the template.
	Parameters []ResourceParameter `json:"parameters,omitempty"`
}

// ResourceParameter is a spec for parameters to pass to a template.
type ResourceParameter struct {
	// Name of the parameter.
	// +required
	Name string `json:"name"`

	// Value for the parameter.
	// +required
	Value string `json:"value"`
}

// BackstageResourceSpec is a spec for Backstage Software Catalog Resource to create from this Resource.
type BackstageResourceSpec struct {
	// Create defines whether the Backstage Resource should be created.
	Create bool `json:"create,omitempty"`
}

// ResourceStatus defines the observed state of Resource.
type ResourceStatus struct {
	// Provisioned indicates whether the Resource has been successfully provisioned.
	Provisioned bool `json:"provisioned,omitempty"`

	// conditions represent the current state of the Resource resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Resource is the Schema for the resources API
//
// +kubebuilder:printcolumn:name="Resource name",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Owner",type="string",JSONPath=".spec.ownerRef.name"
// +kubebuilder:printcolumn:name="Template",type="string",JSONPath=".spec.templateRef.name"
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.provisioned"
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Resource struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Resource
	// +required
	Spec ResourceSpec `json:"spec"`

	// status defines the observed state of Resource
	// +optional
	Status ResourceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ResourceList contains a list of Resource
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resource{}, &ResourceList{})
}
