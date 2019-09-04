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

// FileService Manages the product file
type FileService struct{}

// CreateFile Creates the product file in JSON or YAML
func (fs *FileService) CreateFile(model model) {
	var fileName = fmt.Sprintf("strut.%s", model.fileType.Extension)
	switch model.fileType {
	case file.Types.YAML:
		{
			yamlData, err := yaml.Marshal(model.Product)
			if err != nil {
				log.Fatal(err)
			} else {
				err = ioutil.WriteFile(fileName, yamlData, 0644)
			}
		}
	case file.Types.JSON:
		{
			jsonData, err := json.MarshalIndent(model.Product, "", "  ")
			if err != nil {
				log.Fatal(err)

			} else {
				err = ioutil.WriteFile(fileName, jsonData, 0644)
			}
		}
	}
}

// ReadFile Loads the product file from the CWD
func (fs *FileService) ReadFile(model model) (*Product, error) {
	fileData, err := os.Open(fmt.Sprintf("strut.%s", model.fileType.Extension))
	defer fileData.Close()
	data, _ := ioutil.ReadAll(fileData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	switch model.fileType {
	case file.Types.JSON:
		err = json.Unmarshal(data, model.Product)
	case file.Types.YAML:
		err = yaml.Unmarshal(data, model.Product)
		// default: return new()
	}
	return model.Product, err
}

// UpdateFile Updates the product file
func (fs *FileService) UpdateFile(model model) (*Product, error) {
	return nil, fmt.Errorf("error")
}
