package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/ghodss/yaml"
)

// ProductFileName name of the product file
const ProductFileName = "strut"

// Model The product model
type Model interface {
	LoadProduct() *Product
	SaveProduct(product *Product)
	AddApplication(application *Application)
	AddDependency(dependency *Dependency)
	writeFile() ([]byte, error)
	readFile() ([]byte, error)
	parseFile([]byte)
}

type model struct {
	Product  *Product
	fileType *file.Type
}

// NewProductModel product model constructor
func NewProductModel(fileType *file.Type) Model {
	return &model{
		&Product{},
		fileType,
	}
}

func (m *model) LoadProduct() *Product {
	data, err := m.readFile()
	if err == nil {
		m.parseFile(data)
	} else {
		m.SaveProduct(m.Product)
	}
	return m.Product
}

func (m *model) SaveProduct(product *Product) {
	m.Product = product
	data, _ := m.writeFile()
	m.parseFile(data)
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

// writeFile Writes the product file in JSON or YAML to the CWD
func (m *model) writeFile() ([]byte, error) {
	var fileName = fmt.Sprintf("%s.%s", ProductFileName, m.fileType.Extension)
	var data []byte

	switch m.fileType {
	case file.Types.YAML:
		data, _ = yaml.Marshal(m.Product)
	case file.Types.JSON:
		data, _ = json.MarshalIndent(m.Product, "", "  ")
	}
	ioutil.WriteFile(fileName, data, 0644)
	return m.readFile()
}

// readFile Loads the product file from the CWD
func (m *model) readFile() ([]byte, error) {
	fileData, err := os.Open(fmt.Sprintf("%s.%s", ProductFileName, m.fileType.Extension))
	defer fileData.Close()
	data, err := ioutil.ReadAll(fileData)
	return data, err
}

func (m *model) parseFile(data []byte) {
	switch m.fileType {
	case file.Types.JSON:
		json.Unmarshal(data, m.Product)
	case file.Types.YAML:
		yaml.Unmarshal(data, m.Product)
	}
}
