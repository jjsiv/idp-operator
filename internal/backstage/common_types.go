package backstage

const (
	GroupKind    = "Group"
	SystemKind   = "System"
	ResourceKind = "Resource"
	LocationKind = "Location"

	BackstageV1Alpha1APIVersion = "backstage.io/v1alpha1"

	IDPOperatorAnnotationPrefix = "idp-operator.autopay.pl/"
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

// Validates and returns EntityReference from a string.
// TODO: validation
func EntityReferenceFromString(s string) (EntityReference, error) {
	return EntityReference(s), nil
}
