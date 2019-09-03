package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/ghodss/yaml"
)

// Model The product model
type Model interface {
	CreateFile(data interface{}, fileType *file.Type)
	ReadFile()
	UpdateFile()
	AddApplication()
}

type model struct {
	*Product
	*FileService
}

// New product model constructor
func New() Model {
	return model{
		&Product{},
		&FileService{},
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
	Name    string `json:"name"`
	Version string `json:"version"`
	// repository: {
	//   type: null,
	//   url: null
	// },
	// localConfig,
	// providers: {}
}

// FileService Manages the product file
type FileService struct{}

// CreateFile Creates the product file in JSON or YAML
func (fs *FileService) CreateFile(data interface{}, fileType *file.Type) {
	var fileName = fmt.Sprintf("strut.%s", fileType.Extension)
	switch fileType {
	case file.Types.YAML:
		{
			yamlData, err := yaml.Marshal(data)
			if err != nil {
				log.Fatal(err)
			} else {
				err = ioutil.WriteFile(fileName, yamlData, 0644)
			}
		}
	case file.Types.JSON:
		{
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Fatal(err)

			} else {
				err = ioutil.WriteFile(fileName, jsonData, 0644)
			}
		}
	}
}

// ReadFile Loads the product file from the CWD
func (fs *FileService) ReadFile() {}

// UpdateFile Updates the product file
func (fs *FileService) UpdateFile() {}

// AddApplication Adds an application to product file
func (fs *FileService) AddApplication() {}
