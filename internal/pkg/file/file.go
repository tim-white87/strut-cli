package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

// Types file types
var Types = &TypeRegistry{
	JSON: &Type{"json"},
	YAML: &Type{"yaml"},
}

// TypeList list of types
var TypeList = []*Type{Types.JSON, Types.YAML}

// TypeRegistry Struct for various file types
type TypeRegistry struct {
	JSON *Type
	YAML *Type
}

// Type file type
type Type struct {
	Extension string
}

// Model File model
type Model struct {
	FileType *Type
}

// WriteFile Writes the product file in JSON or YAML to the CWD
func (m *Model) WriteFile(name string, fileData interface{}) ([]byte, error) {
	var fileName = fmt.Sprintf("%s.%s", name, m.FileType.Extension)
	var data []byte

	switch m.FileType {
	case Types.YAML:
		data, _ = yaml.Marshal(fileData)
	case Types.JSON:
		data, _ = json.MarshalIndent(fileData, "", "  ")
	}
	ioutil.WriteFile(fileName, data, 0644)
	return m.ReadFile(name)
}

// ReadFile Loads the product file from the CWD
func (m *Model) ReadFile(name string) ([]byte, error) {
	fileData, err := os.Open(fmt.Sprintf("%s.%s", name, m.FileType.Extension))
	defer fileData.Close()
	data, err := ioutil.ReadAll(fileData)
	return data, err
}

// ParseFile Parses file contents back to object
func (m *Model) ParseFile(data []byte, parsedData interface{}) {
	switch m.FileType {
	case Types.JSON:
		json.Unmarshal(data, parsedData)
	case Types.YAML:
		yaml.Unmarshal(data, parsedData)
	}
}
