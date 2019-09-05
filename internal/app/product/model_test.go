package product

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/ghodss/yaml"
)

func init() {
	os.Chdir("../../../test/testdata")
	// put this teardown maybe? os.Chdir("..")
}

func TestLoadProduct(t *testing.T) {
	t.Fatalf("Expected implementation")
}

func TestSaveProduct(t *testing.T) {
	// Arrange
	var version = "2.0.0"
	ym := New(file.Types.YAML).(*model)
	ym2 := New(file.Types.YAML).(*model)

	// Act
	ym.Product.Version = version
	ym.SaveProduct()
	data, _ := ym2.readFile()
	ym2.parseFile(data)

	// Assert
	if ym.Product.Version != ym2.Product.Version {
		t.Fatalf("Expected version: %s to be updated to: %s", ym.Product.Version, ym2.Product.Version)
	}
}

func TestAddApplication(t *testing.T) {
	// Arrange
	m := New(file.Types.YAML).(*model)
	var appName = "Derp"
	var exists = false

	// Act
	m.AddApplication(&Application{Name: appName})

	// Assert
	for _, p := range m.Product.Applications {
		if p.Name == appName {
			exists = true
		}
	}
	if !exists {
		t.Fatalf("Expected application to exist in applications list.")
	}
}

func TestWriteFile(t *testing.T) {
	// Arrange
	yamlProductModel := New(file.Types.YAML).(*model)
	jsonProductModel := New(file.Types.JSON).(*model)

	// Act
	_, yamlErr := yamlProductModel.writeFile()
	_, jsonErr := jsonProductModel.writeFile()

	// Assert
	if _, err := os.Stat("./strut.yaml"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.yaml file to exist.")
	}
	if _, err := os.Stat("./strut.json"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.json file to exist.")
	}
	if yamlErr != nil {
		t.Fatalf("Expected no YAML error: %s", yamlErr)
	}
	if jsonErr != nil {
		t.Fatalf("Expected no JSON error: %s", jsonErr)
	}
}

func TestReadFile(t *testing.T) {
	// Arrange
	yamlProductModel := New(file.Types.YAML).(*model)
	jsonProductModel := New(file.Types.JSON).(*model)

	// Act
	yamlData, err := yamlProductModel.readFile()
	jsonData, err := jsonProductModel.readFile()

	// Assert
	if err != nil {
		t.Fatalf("Expected to read file: %s", err)
	}
	if err := yaml.Unmarshal(yamlData, yamlProductModel.Product); err != nil {
		t.Fatalf("Expected strut.yaml file to parse to model: %s", err)
	}
	if err := json.Unmarshal(jsonData, jsonProductModel.Product); err != nil {
		t.Fatalf("Expected strut.json file to parse to model: %s", err)
	}
}

func TestParseFile(t *testing.T) {
	// Arrange
	var name = "Derp"
	yamlProductModel := New(file.Types.YAML).(*model)
	jsonProductModel := New(file.Types.JSON).(*model)

	// Act
	yamlProductModel.Product.Name = name
	yamlData, _ := yamlProductModel.readFile()
	yamlProductModel.parseFile(yamlData)

	jsonProductModel.Product.Name = name
	jsonData, _ := jsonProductModel.readFile()
	jsonProductModel.parseFile(jsonData)

	// Assert
	if yamlProductModel.Product.Name == name {
		t.Fatalf("Expected product model to get parsed from YAML file")
	}
	if jsonProductModel.Product.Name == name {
		t.Fatalf("Expected product model to get parsed from JSON file")
	}
}
