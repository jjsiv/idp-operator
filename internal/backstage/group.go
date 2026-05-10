package backstage

import (
	"fmt"
	"slices"

	idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"
)

// A group describes an organizational entity, such as for example a team, a business unit, or a loose collection of people in an interest group.
// Members of these groups are modeled in the catalog as kind User.
type Group struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta `json:"metadata"`
	Spec     GroupSpec  `json:"spec"`
}

// Returns a Group with TypeMeta pre-filled.
func NewGroup() *Group {
	return &Group{
		TypeMeta: TypeMeta{
			APIVersion: BackstageV1Alpha1APIVersion,
			Kind:       GroupKind,
		},
	}
}

type GroupOption func(*Group, *idpv1alpha1.Tenant)

// WithGroupRole creates a role subgroup.
// Role subgroup will have:
//
// - its type set to <role>-role-group
//
// - name set to <tenant>-<role>
//
// - parent set to tenant name
//
// - only members with the role assigned to it
func WithGroupRole(role string) GroupOption {
	return func(g *Group, tenant *idpv1alpha1.Tenant) {
		originalName := g.Metadata.Name
		g.Metadata.Name = originalName + "-" + role
		g.Metadata.Title = fmt.Sprintf("%s role group for %s", role, originalName)
		g.Metadata.Annotations["idp-operator.autopay.pl/member-role"] = role

		g.Spec.Type = fmt.Sprintf("%s-role-group", role)
		// Should already be valid so we don't check error
		parentEntityRef, _ := EntityReferenceFromString(originalName)
		g.Spec.Parent = parentEntityRef

		var members []EntityReference
		for _, member := range tenant.Spec.Members {
			if slices.Contains(member.Roles, role) {
				// Should already be valid so we don't check error
				memberEntityRef, _ := EntityReferenceFromString(member.EntityName)
				members = append(members, memberEntityRef)
			}
		}

		g.Spec.Members = members
	}
}

// WithDescription sets a custom description on the Group.
func WithDescription(description string) GroupOption {
	return func(g *Group, tenant *idpv1alpha1.Tenant) {
		g.Metadata.Description = description
	}
}

func (g *Group) FromTenant(tenant *idpv1alpha1.Tenant, opts ...GroupOption) (*Group, error) {
	g.Metadata = ObjectMeta{
		Name:        tenant.Name,
		Title:       tenant.Spec.DisplayName,
		Description: tenant.Spec.Description,
		Annotations: map[string]string{
			"idp-operator.autopay.pl/tenant": tenant.Name + "/" + tenant.Namespace,
		},
	}

	g.Spec = GroupSpec{
		Type: "team",
	}

	for _, member := range tenant.Spec.Members {
		entityRef, err := EntityReferenceFromString(member.EntityName)
		if err != nil {
			return nil, err
		}

		g.Spec.Members = append(g.Spec.Members, entityRef)
	}

	for _, opt := range opts {
		opt(g, tenant)
	}

	return g, nil
}

// GroupSpec defines spec for Backstage Group.
type GroupSpec struct {
	// The type of group as a string, e.g. team.
	// There is currently no enforced set of values for this field,
	// so it is left up to the adopting organization to choose a nomenclature that matches their org hierarchy.
	Type string `json:"type"`

	// Optional profile information about the Group.
	Profile *GroupProfile `json:"profile,omitempty"`

	// The immediate parent group in the hierarchy, if any. Not all groups must have a parent; the catalog supports multi-root hierarchies.
	// Groups may however not have more than one parent.
	Parent EntityReference `json:"parent,omitempty"`

	// The immediate child groups of this group in the hierarchy (whose parent field points to this group).
	// The list must be present, but may be empty if there are no child groups.
	Children []EntityReference `json:"children"`

	// The users that are direct members of this group.
	Members []EntityReference `json:"members,omitempty"`
}

// Optional profile information about the group, mainly for display purposes. All fields of this structure are also optional.
// The email would be a group email of some form, that the group may wish to be used for contacting them.
// The picture is expected to be a URL pointing to an image that's representative of the group, and that a browser could fetch and render on a group page or similar.
type GroupProfile struct {
	// A human-readable name for the group.
	DisplayName string `json:"displayName,omitempty"`

	// An email the group may wish to be used for contacting them.
	Email string `json:"email,omitempty"`

	// A URL pointing to an image that's representative of the group.
	Picture string `json:"picture,omitempty"`
}
