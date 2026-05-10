package backstage

import idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"

// A system is a collection of resources and components. The system may expose or consume one or several APIs.
// It is viewed as abstraction level that provides potential consumers insights into exposed features without needing a too detailed view into the details of all components.
// This also gives the owning team the possibility to decide about published artifacts and APIs.
type System struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta `json:"metadata"`
	Spec     SystemSpec `json:"spec"`
}

// Returns a System with TypeMeta pre-filled.
func NewSystem() *System {
	return &System{
		TypeMeta: TypeMeta{
			APIVersion: BackstageV1Alpha1APIVersion,
			Kind:       SystemKind,
		},
	}
}

func (s *System) FromTenant(tenant *idpv1alpha1.Tenant) (*System, error) {
	s.Metadata = ObjectMeta{
		Name:        tenant.Name,
		Title:       tenant.Spec.DisplayName,
		Description: tenant.Spec.Description,
		Annotations: map[string]string{
			"idp-operator.autopay.pl/tenant": tenant.Namespace + "/" + tenant.Name,
		},
	}

	ownerRef, err := EntityReferenceFromString(tenant.Name)
	if err != nil {
		return nil, err
	}

	s.Spec = SystemSpec{
		Type:  "product",
		Owner: ownerRef,
	}

	return s, nil
}

type SystemSpec struct {
	// An Entity reference to owning Group.
	Owner EntityReference `json:"owner"`

	// An Entity reference to Domain the system belongs to.
	Domain EntityReference `json:"domain,omitempty"`

	// The type of the system.
	Type string `json:"type,omitempty"`
}
