package product

import (
	"fmt"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

// Model The product model
type Model interface {
	createFile(model model)
	readFile(model model) (*Product, error)
	updateFile(model model) (*Product, error)
	SaveProduct() (*Product, error)
	AddApplication(application *Application) (*Product, error)
}

type model struct {
	*fileService
	*productService
	Product  *Product
	fileType *file.Type
}

// New product model constructor
func New(name string, fileType *file.Type) Model {
	return model{
		&fileService{},
		&productService{},
		&Product{Name: name},
		fileType,
	}
}

type productService struct{}

func (ps *productService) SaveProduct() (*Product, error) {
	return nil, fmt.Errorf("error")
}

// AddApplication Adds an application to product and updates the file
func (ps *productService) AddApplication(application *Application) (*Product, error) {
	return nil, fmt.Errorf("error")
}
