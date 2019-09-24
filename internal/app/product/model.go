package product

import (
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/cecotw/strut-cli/internal/pkg/provider"
)

// ProductFileName name of the product file
const ProductFileName = "strut"

// Model The product model
type Model interface {
	LoadProduct() *Product
	SaveProduct(product *Product)
	AddApplication(application *Application)
	AddDependency(dependency *Dependency)
	UpdateApplications(applications []*Application)
}

type model struct {
	Product *Product
	*file.Model
}

// NewProductModel product model constructor
func NewProductModel(fileType *file.Type) Model {
	return &model{
		&Product{},
		&file.Model{fileType},
	}
}

func (m *model) LoadProduct() *Product {
	data, err := m.ReadFile(ProductFileName)
	if err == nil {
		m.ParseFile(data, m.Product)
	} else {
		m.SaveProduct(m.Product)
	}
	if m.Product.Applications != nil {
		m.mapProviderResources()
	}
	return m.Product
}

func (m *model) SaveProduct(product *Product) {
	m.Product = product
	data, _ := m.WriteFile(ProductFileName, m.Product)
	m.ParseFile(data, m.Product)
}

// AddApplication Adds an application to product and updates the file
func (m *model) AddApplication(application *Application) {
	m.Product.Applications = append(m.Product.Applications, application)
	m.SaveProduct(m.Product)
}

// AddDependency Adds a dependency to the product
func (m *model) AddDependency(dependency *Dependency) {
	m.Product.Dependencies = append(m.Product.Dependencies, dependency)
	m.SaveProduct(m.Product)
}

// UpdateApplications Updates the applications in the product
func (m *model) UpdateApplications(applications []*Application) {
	m.Product.Applications = applications
	m.SaveProduct(m.Product)
}

func (m *model) mapProviderResources() {
	m.Product.ProvisionMap = make(map[int][]*provider.Resource)
	for _, app := range m.Product.Applications {
		for _, r := range app.Resources {
			m.Product.ProvisionMap[r.Priority] = append(m.Product.ProvisionMap[r.Priority], r)
		}
	}
}
