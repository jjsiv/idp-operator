package backstage

const (
	GroupKind    = "Group"
	SystemKind   = "System"
	ResourceKind = "Resource"

	BackstageV1Alpha1APIVersion = "backstage.io/v1alpha1"
)

type TypeMeta struct {
	// The apiVersion is the version of specification format for that particular entity that the specification is made against.
	APIVersion string `json:"apiVersion"`

	// The kind is the high level entity type being described.
	Kind string `json:"kind"`
}

type ObjectMeta struct {
	// The name of the entity. This name is both meant for human eyes to recognize the entity,
	// and for machines and other components to reference the entity (e.g. in URLs or from other entity specification files).
	Name string `json:"name"`

	// The ID of a namespace that the entity belongs to.
	// This field is optional, and has no special semantics apart from bounding the name uniqueness constraint if specified.
	Namespace string `json:"namespace,omitempty"`

	// A display name of the entity, to be presented in user interfaces instead of the name property above, when available.
	Title string `json:"title,omitempty"`

	// A human readable description of the entity, to be shown in Backstage.
	// Should be kept short and informative, suitable to give an overview of the entity's purpose at a glance.
	// More detailed explanations and documentation should be placed elsewhere.
	Description string `json:"description,omitempty"`

	// Labels are optional key/value pairs of that are attached to the entity, and their use is identical to Kubernetes object labels.
	Labels map[string]string `json:"labels,omitempty"`

	// An object with arbitrary non-identifying metadata attached to the entity, identical in use to Kubernetes object annotations.
	Annotations map[string]string `json:"annotations,omitempty"`

	// A list of single-valued strings, for example to classify catalog entities in various ways.
	// This is different to the labels in metadata, as labels are key-value pairs.
	Tags []string `json:"tags,omitempty"`

	// A list of external hyperlinks related to the entity.
	// Links can provide additional contextual information that may be located outside of Backstage itself.
	// For example, an admin dashboard or external CMS page.
	Links []Link `json:"links,omitempty"`
}

type Link struct {
	// [Required] A url in a standard uri format (e.g. https://example.com/some/page)
	URL string `json:"url"`

	// [Optional] A user friendly display name for the link.
	Title string `json:"title,omitempty"`

	// [Optional] A key representing a visual icon to be displayed in the UI.
	Icon string `json:"icon,omitempty"`

	// [Optional] An optional value to categorize links into specific groups.
	Type string `json:"type,omitempty"`
}

// EntityReference is a Backstage entity reference in format "<type>:<namespace>/<entity>"
type EntityReference string

// A group describes an organizational entity, such as for example a team, a business unit, or a loose collection of people in an interest group.
// Members of these groups are modeled in the catalog as kind User.
type Group struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:",inline"`
	Spec       GroupSpec `json:"spec"`
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

// A system is a collection of resources and components. The system may expose or consume one or several APIs.
// It is viewed as abstraction level that provides potential consumers insights into exposed features without needing a too detailed view into the details of all components.
// This also gives the owning team the possibility to decide about published artifacts and APIs.
type System struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:",inline"`
	Spec       SystemSpec `json:"spec"`
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

type SystemSpec struct {
	// An Entity reference to owning Group.
	Owner EntityReference `json:"owner"`

	// An Entity reference to Domain the system belongs to.
	Domain EntityReference `json:"domain"`

	// The type of the system.
	Type string `json:"type,omitempty"`
}

// A Resource describes the infrastructure a system needs to operate, like BigTable databases, Pub/Sub topics, S3 buckets or CDNs.
// Modelling them together with components and systems allows to visualize resource footprint, and create tooling around them.
type Resource struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:",inline"`
	Spec       ResourceSpec `json:"spec"`
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
