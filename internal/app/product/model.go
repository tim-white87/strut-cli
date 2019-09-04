package product

import (
	"fmt"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

// Model The product model
type Model interface {
	CreateFile(model model)
	ReadFile(model model) (*Product, error)
	UpdateFile(model model) (*Product, error)
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

// AddApplication Adds an application to product and updates the file
func (m *model) AddApplication(application *Application) (*Product, error) {
	return nil, fmt.Errorf("error")
}
