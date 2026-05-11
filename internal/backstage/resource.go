package backstage

import (
	"encoding/json"
	"fmt"

	idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"
)

// A Resource describes the infrastructure a system needs to operate, like BigTable databases, Pub/Sub topics, S3 buckets or CDNs.
// Modelling them together with components and systems allows to visualize resource footprint, and create tooling around them.
type Resource struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta   `json:"metadata"`
	Spec     ResourceSpec `json:"spec"`
}

// Returns a Resource with TypeMeta pre-filled.
func NewResource() *Resource {
	return &Resource{
		TypeMeta: TypeMeta{
			APIVersion: BackstageV1Alpha1APIVersion,
			Kind:       ResourceKind,
		},
	}
}

func (r *Resource) FromIDPResource(idpResource *idpv1alpha1.Resource) (*Resource, error) {
	params, err := json.Marshal(idpResource.Spec.TemplateRef.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource parameters: %s", err.Error())
	}

	r.Metadata = ObjectMeta{
		Name:        idpResource.Spec.TemplateRef.Name + "-" + idpResource.Name,
		Title:       idpResource.Spec.DisplayName,
		Description: idpResource.Spec.Description,
		Annotations: map[string]string{
			"idp-operator.autopay.pl/resource":     idpResource.Namespace + "/" + idpResource.Name,
			"idp-operator.autopay.pl/template":     idpResource.Spec.TemplateRef.Namespace + "/" + idpResource.Spec.TemplateRef.Name,
			"idp-operator.autopay.pl/owner-tenant": idpResource.Spec.OwnerRef.Name,
			"idp-operator.autopay.pl/parameters":   string(params),
		},
	}

	ownerRef, _ := EntityReferenceFromString(idpResource.Spec.OwnerRef.Name)

	// Currently these have the same name as Tenant, but eventually they should be retrieved from status
	// Type should also be set based on another field on the template (which currently does not exist)
	r.Spec = ResourceSpec{
		Owner:  ownerRef,
		Type:   idpResource.Spec.TemplateRef.Name,
		System: ownerRef,
	}

	return r, nil
}

type ResourceSpec struct {
	// An Entity reference to owning Group.
	Owner EntityReference `json:"owner"`

	// The type of the Resource.
	Type string `json:"type"`

	// An entity reference to the system that the resource belongs.
	System EntityReference `json:"system,omitempty"`

	// An array of entity references to the components and resources that the resource depends on.
	DependsOn []EntityReference `json:"dependsOn,omitempty"`

	// An array of entity references to the components and resources that the resource is a dependency of.
	DependencyOf []EntityReference `json:"dependencyOf,omitempty"`
}
