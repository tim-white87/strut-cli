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

// fileService Manages the product file
type fileService struct{}

// createFile Creates the product file in JSON or YAML
func (fs *fileService) writeFile(model model) (*Product, error) {
	var fileName = fmt.Sprintf("%s.%s", ProductFileName, model.fileType.Extension)
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
	return fs.readFile(model)
}

// readFile Loads the product file from the CWD
func (fs *fileService) readFile(model model) (*Product, error) {
	fileData, err := os.Open(fmt.Sprintf("%s.%s", ProductFileName, model.fileType.Extension))
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
