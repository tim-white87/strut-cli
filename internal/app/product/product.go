package product

import (
	"github.com/cecotw/strut-cli/internal/pkg/file"
)

// Model The product model
type Model struct {
	Product     *Product
	FileService *FileService
}

// Product Product anemic model
type Product struct {
	Name         string
	Version      string
	Dependencies []Dependency
	Applications []Application
}

// Dependency Product tool/language/installation dependency
type Dependency struct{}

// Application Application defined in the product
type Application struct{}

// FileService Manages the product file
type FileService struct{}

// CreateFile Creates the product file in JSON or YAML
func (fs *FileService) CreateFile(name string, filetype *file.TypeRegistry) {}

// ReadFile Loads the product file from the CWD
func (fs *FileService) ReadFile() {}

// UpdateFile Updates the product file
func (fs *FileService) UpdateFile() {}

// DeleteFile Deletes the product file
func (fs *FileService) DeleteFile() {}

// AddApplication Adds an application to product file
func (fs *FileService) AddApplication() {}
