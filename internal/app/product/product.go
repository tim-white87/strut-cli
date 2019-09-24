package product

import "github.com/cecotw/strut-cli/internal/pkg/provider"

// Product Product anemic model
type Product struct {
	Name         string                       `json:"name"`
	Version      string                       `json:"version"`
	Description  string                       `json:"description,omitempty"`
	Dependencies []*Dependency                `json:"dependencies"`
	Applications []*Application               `json:"applications"`
	ProvisionMap map[int][]*provider.Resource `json:"-"`
}

// Dependency Product tool/language/installation dependency
type Dependency struct {
	Name    string   `json:"name"`
	Install []string `json:"install"`
}

// Application Application defined in the product
type Application struct {
	Name        string               `json:"name"`
	Version     string               `json:"version"`
	Repository  *Repository          `json:"repository"`
	LocalConfig *LocalConfig         `json:"localConfig"`
	Resources   []*provider.Resource `json:"resources"`
}

// Repository Application repository
type Repository struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// LocalConfig Local application config
type LocalConfig struct {
	Path      string                `json:"path"`
	Commands  *LocalCommandRegistry `json:"commands"`
	Artifacts []string              `json:"artifacts"`
}

// LocalCommandRegistry Local config command types
type LocalCommandRegistry struct {
	Install  []string `json:"install"`
	Validate []string `json:"validate"`
	Build    []string `json:"build"`
	Start    []string `json:"start"`
	Deploy   []string `json:"deploy"`
}
