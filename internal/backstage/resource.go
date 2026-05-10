package backstage

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
