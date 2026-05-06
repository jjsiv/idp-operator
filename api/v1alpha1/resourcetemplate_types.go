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

// ResourceTemplateSpec defines the desired state of ResourceTemplate
type ResourceTemplateSpec struct {
	// +optional
	Metadata `json:",inline"`

	// List of parameters for this template.
	// These will be made available in provisioning templates through "{{ .parameters.<name> }}".
	// +optional
	Parameters []TemplateParameter `json:"parameters,omitempty"`

	// Defines how resources will be provisioned from this template.
	// +required
	Provisioning ResourceProvisioning `json:"provisioning"`
}

// TemplateParameter is spec for template parameters.
type TemplateParameter struct {
	// Name of the parameter. Must be unique within the template.
	// +required
	Name string `json:"name"`

	// Schema for this parameter.
	// +required
	Schema TemplateParameterSchema `json:"schema"`
}

// TemplateParameterSchema defines parameter schema.
type TemplateParameterSchema struct {
	// Parameter type.
	// One of: string | number | enum
	// +required
	Type ParameterType `json:"type"`

	// If true, value of this parameter on the Resource cannot be changed after creation.
	// +optional
	Immutable bool `json:"immutable,omitempty"`

	// If true, setting this parameter on the Resource will not be required.
	// +optional
	Optonal bool `json:"optional,omitempty"`

	// Default value for the parameter. Must match Type.
	// +optional
	Default string `json:"default,omitempty"`

	// Only for Type=enum.
	// List of values to choose from.
	// +optional
	Values []string `json:"values,omitempty"`
}

// ResourceProvisioning defines the provisioning configuration for a ResourceTemplate.
type ResourceProvisioning struct {
	// Git defines a strategy of provisioning resource by creating files in Git repositories.
	Git []*ResourceProvisioningGit `json:"git,omitempty"`
}

type ResourceProvisioningGit struct {
	GitProvisioning `json:",inline"`

	// Files to create in the repository.
	// +required
	Templates []FileTemplate `json:"templates"`
}

// GitProvisioning defines a provisioning method that creates files in Git repositories.
type GitProvisioning struct {
	// URL of the repository to create files in.
	// +required
	RepositoryURL string `json:"repositoryURL"`

	// Branch to push to. If unset, default will be assumed.
	// +optional
	Branch string `json:"branch,omitempty"`

	// Configuration for repository connection.
	// +required
	Config GitConfig `json:"config"`

	// Message for the commit that will be created.
	// +required
	CommitMessage string `json:"commitMessage"`
}

// GitConfig defines configuration for repository connection.
type GitConfig struct {
	// Information about commit author.
	// +required
	Author GitAuthorConfig `json:"author"`

	// Authentication information for cloning and pushing to the repository.
	// +required
	Auth GitAuthConfig `json:"auth"`
}

// GitAuthorConfig defines information about commit author.
type GitAuthorConfig struct {
	// Commit author's name.
	// +required
	Name string `json:"name"`

	// Commit author's email.
	// +required
	Email string `json:"email"`
}

// GitAuthConfig defines authentication information for cloning and pushing to the repository.
type GitAuthConfig struct {
	// Reference to a Kubernetes secret. Must be found in the same namespace.
	// +required
	SecretRef SecretReference `json:"secretRef"`
}

// SecretReference is a reference to a Kubernetes secret.
type SecretReference struct {
	// Name of the secret.
	// +required
	Name string `json:"name"`

	// Namespace of the secret.
	// If unset, template's namespace will be used.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Key containing the private key for authentication.
	// +required
	Key string `json:"key"`
}

// FileTemplate defines a template for a file to create.
type FileTemplate struct {
	// Path at which the file will be created (without file name).
	// Supports Go templating.
	// +required
	Path string `json:"path"`

	// Name of the file to create.
	// Supports Go templating.
	// +required
	Filename string `json:"filename"`

	// Content of the file.
	// Supports Go templating.
	// +required
	Content string `json:"content"`
}

type ParameterType string

const (
	ParameterTypeString ParameterType = "string"
	ParameterTypeNumber ParameterType = "number"
	ParameterTypeEnum   ParameterType = "enum"
)

// ResourceTemplateStatus defines the observed state of ResourceTemplate.
type ResourceTemplateStatus struct {
	// conditions represent the current state of the ResourceTemplate resource.
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
// +kubebuilder:printcolumn:name="Template name",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ResourceTemplate is the Schema for the resourcetemplates API
type ResourceTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ResourceTemplate
	// +required
	Spec ResourceTemplateSpec `json:"spec"`

	// status defines the observed state of ResourceTemplate
	// +optional
	Status ResourceTemplateStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ResourceTemplateList contains a list of ResourceTemplate
type ResourceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ResourceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceTemplate{}, &ResourceTemplateList{})
}
