package product

import (
	"fmt"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

// Model The product model
type Model interface {
	writeFile(model model) (*Product, error)
	readFile(model model) (*Product, error)
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
func New(fileType *file.Type) Model {
	return model{
		&fileService{},
		&productService{},
		&Product{},
		fileType,
	}
}

type productService struct{}

func (ps *productService) LoadProduct() (*Product, error) {
	return nil, fmt.Errorf("error")
}

func (ps *productService) SaveProduct() (*Product, error) {
	return nil, fmt.Errorf("error")
}

// AddApplication Adds an application to product and updates the file
func (ps *productService) AddApplication(application *Application) (*Product, error) {
	return nil, fmt.Errorf("error")
}
