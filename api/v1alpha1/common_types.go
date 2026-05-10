package v1alpha1

// Metadata is a common struct for all resources to provide human readable information about the resource.
type Metadata struct {
	// Human readalbe name of the resource.
	DisplayName string `json:"displayName,omitempty"`

	// Description about this resource.
	Description string `json:"description,omitempty"`
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

type BackstageCatalogSpec struct {
	// Strategy for provisioning Backstage Software catalog entity data.
	// +required
	Provisioning BackstageCatalogProvisioningMethods `json:"provisioning"`
}

type BackstageCatalogProvisioningMethods struct {
	// Provision entity data as file in a Git repository.
	// This is the only supported option currently.
	// +required
	Git *BackstageCatalogGitProvisioning `json:"git"`
}

type BackstageCatalogGitProvisioning struct {
	GitProvisioning `json:",inline"`

	// Path (full directory + filename) to create catalog entities.
	// +required
	Path string `json:"path"`
}
