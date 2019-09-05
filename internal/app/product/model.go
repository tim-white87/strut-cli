package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/ghodss/yaml"
)

// ProductFileName name of the product file
const ProductFileName = "strut"

// Model The product model
type Model interface {
	LoadProduct()
	SaveProduct()
	AddApplication(application *Application)
	writeFile() ([]byte, error)
	readFile() ([]byte, error)
	parseFile([]byte)
}

type model struct {
	Product  *Product
	fileType *file.Type
}

// New product model constructor
func New(fileType *file.Type) Model {
	return &model{
		&Product{},
		fileType,
	}
}

func (m *model) LoadProduct() {

}

func (m *model) SaveProduct() {
	data, err := m.writeFile()
	m.parseFile(data)
	if err != nil {
		log.Fatal(err)
	}
}

// AddApplication Adds an application to product and updates the file
func (m *model) AddApplication(application *Application) {
	m.Product.Applications = append(m.Product.Applications, *application)
	m.SaveProduct()
}

// writeFile Writes the product file in JSON or YAML to the CWD
func (m *model) writeFile() ([]byte, error) {
	var fileName = fmt.Sprintf("%s.%s", ProductFileName, m.fileType.Extension)
	var data []byte
	var err error
	switch m.fileType {
	case file.Types.YAML:
		data, err = yaml.Marshal(m.Product)
	case file.Types.JSON:
		data, err = json.MarshalIndent(m.Product, "", "  ")
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return m.readFile()
}

// readFile Loads the product file from the CWD
func (m *model) readFile() ([]byte, error) {
	fileData, err := os.Open(fmt.Sprintf("%s.%s", ProductFileName, m.fileType.Extension))
	defer fileData.Close()
	data, err := ioutil.ReadAll(fileData)
	if err != nil {
		log.Fatal(err)
	}
	return data, err
}

func (m *model) parseFile(data []byte) {
	var err error
	switch m.fileType {
	case file.Types.JSON:
		err = json.Unmarshal(data, m.Product)
	case file.Types.YAML:
		err = yaml.Unmarshal(data, m.Product)
	}
	if err != nil {
		log.Fatal(err)
	}
}
