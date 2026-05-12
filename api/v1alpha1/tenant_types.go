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

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	Metadata `json:",inline"`

	// Configuring an identity provider allows for creating of identity groups for members of this project.
	// +optional
	IdentityProvider *IdentityProviderConfiguration `json:"identityProvider,omitempty"`

	// Members of this Tenant.
	// +required
	Members []TenantMember `json:"members,omitempty"`

	// Kubernetes infrastructure associated with this Tenant.
	// +optional
	Kubernetes TenantKubernetesSpec `json:"kubernetes,omitempty"`

	// Gitlab infrastructure associated with this Tenant.
	// +optional
	Gitlab TenantGitlabSpec `json:"gitlab,omitempty"`

	// Configuration for provisoning Backstage Software catalog Entities associated with this Tenant.
	// +required
	BackstageCatalog TenantBackstageCatalogSpec `json:"backstageCatalog,omitempty"`
}

type TenantGitlabSpec struct {
	// Gitlab groups to associate.
	// +optional
	Groups []TenantGitlabGroup `json:"groups,omitempty"`
}

type TenantGitlabGroup struct {
	// Name of the group.
	// +required
	Name string `json:"name"`
}

type TenantKubernetesSpec struct {
	// Namespaces that will be owned by this Tenant.
	// Owned namespaces will have special annotations set on them to mark the relationshop.
	// +optional
	Namespaces TenantKubernetesNamespace `json:"namespaces,omitempty"`
}

type TenantKubernetesNamespace struct {
	// Name of the namespace.
	// +required
	Name string `json:"name"`

	// Name of the cluster.
	// +required
	Cluster string `json:"cluster"`
}

type TenantBackstageCatalogSpec struct {
	BackstageCatalogSpec `json:",inline"`

	// Backstage System entities to create for this Tenant.
	// +optional
	Systems []BackstageCatalogSystem `json:"systems,omitempty"`
}

type BackstageCatalogSystem struct {
	// Name of the System.
	// +required
	Name string `json:"name"`

	// Metadata fields to set on the Backstage entity.
	// +optional
	Metadata BackstageCatalogSystemMetadata `json:"metadata,omitempty"`

	// System spec.
	// +required
	Spec BackstageCatalogSystemSpec `json:"spec"`
}

type BackstageCatalogSystemSpec struct {
	// Type of the System.
	// +required
	Type string `json:"type"`
}

// Subset of available metadata fields on Backstage entities.
type BackstageCatalogSystemMetadata struct {
	// Description for this System.
	// +optional
	Description string `json:"description,omitempty"`

	// Name of this entity that will be displayed to viewers.
	// +optional
	Title string `json:"title,omitempty"`

	// Links to include on the entity.
	// +optional
	Links []BackstageCatalogSystemLink `json:"links,omitempty"`
}

type BackstageCatalogSystemLink struct {
	// A url in a standard uri format (e.g. https://example.com/some/page)
	// +required
	URL string `json:"url"`

	// A user friendly display name for the link.
	// +optional
	Title string `json:"title,omitempty"`

	// A key representing a visual icon to be displayed in the UI.
	// +optional
	Icon string `json:"icon,omitempty"`

	// An optional value to categorize links into specific groups.
	// +optional
	Type string `json:"type,omitempty"`
}

type IdentityProviderConfiguration struct {
	// Configuration for Okta integration.
	// +optional
	Okta *OktaConfiguration `json:"okta,omitempty"`

	// Configuration for LDAP integration.
	// +optional
	LDAP *LDAPConfiguration `json:"ldap,omitempty"`
}

type OktaConfiguration struct {
	// Groups to create.
	// +required
	Groups []IdentityGroup `json:"groups,omitempty"`
}

type LDAPConfiguration struct {
	// Groups to create.
	// +required
	Groups []IdentityGroup `json:"groups,omitempty"`
}

type IdentityGroup struct {
	// Name of the group to create.
	// +required
	Name string `json:"name"`

	// Only members with the specified roles will be added to this group.
	// +optional
	Roles []string `json:"roles,omitempty"`
}

type TenantMember struct {
	// Email address of the user.
	// +required
	Email string `json:"email"`

	// Name of the User Entity in Backstage.
	// +required
	EntityName string `json:"entityName"`

	// Roles to assign to this member.
	// +optional
	Roles []string `json:"roles,omitempty"`
}

// TenantStatus defines the observed state of Tenant.
type TenantStatus struct {
	// Provisioned indicates whether the Tenant has been successfully provisioned.
	Provisioned bool `json:"provisioned,omitempty"`

	// Top level Backstage group created from this Tenant.
	BackstageGroup string `json:"backstageGroup,omitempty"`

	// conditions represent the current state of the Tenant resource.
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

// Tenant is the Schema for the tenants API
// +kubebuilder:printcolumn:name="Tenant name",type="string",JSONPath=".spec.tenantName"
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.provisioned"
// +kubebuilder:printcolumn:name="Backstage group",type="string",JSONPath=".status.backstageGroup"
// +kubebuilder:printcolumn:name="Description",type="string",JSONPath=".spec.description"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Tenant struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Tenant
	// +required
	Spec TenantSpec `json:"spec"`

	// status defines the observed state of Tenant
	// +optional
	Status TenantStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
