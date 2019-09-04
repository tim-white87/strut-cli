package product

import (
	"github.com/cecotw/strut-cli/internal/pkg/file"
)

// Model The product model
type Model interface {
	CreateFile(model model)
	ReadFile(model model) (*Product, error)
	UpdateFile()
	AddApplication()
}

type model struct {
	*FileService
	Product  *Product
	FileType *file.Type
}

// New product model constructor
func New(name string, fileType *file.Type) Model {
	return model{
		&FileService{},
		&Product{Name: name},
		fileType,
	}
}

// Product Product anemic model
type Product struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Description  string        `json:"description,omitempty"`
	Dependencies []Dependency  `json:"dependencies"`
	Applications []Application `json:"applications"`
}

// Dependency Product tool/language/installation dependency
type Dependency struct{}

// Application Application defined in the product
type Application struct {
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Repository  *Repository  `json:"repository"`
	LocalConfig *LocalConfig `json:"localConfig"`
	Providers   *[]Provider  `json:"providers"`
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

// Provider Hosted application provider
type Provider struct {
	Name      string      `json:"name"`
	Resources *[]Resource `json:"resources"`
	*ResourceCommands
}

// Resource Provider resource (i.e. cloudformation)
type Resource struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Body string `json:"body"`
	*ResourceCommands
}

// ResourceCommands Custom resource commands
type ResourceCommands struct {
	PreProvision  []string `json:"preProvision"`
	PostProvision []string `json:"postProvision"`
}
