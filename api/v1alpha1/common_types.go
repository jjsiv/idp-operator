package v1alpha1

// Metadata is a common struct for all resources to provide human readable information about the resource.
type Metadata struct {
	// Human readalbe name of the resource.
	DisplayName string `json:"displayName,omitempty"`

	// Description about this resource.
	Description string `json:"description,omitempty"`
}
