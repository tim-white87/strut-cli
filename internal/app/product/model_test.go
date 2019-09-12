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
	// Arrange
	var name = "Derp"
	m := NewProductModel(file.Types.YAML).(*model)

	// Act
	m.LoadProduct()
	m.Product.Name = name
	m.LoadProduct()

	// Assert
	if m.Product.Name == name {
		t.Fatalf("Expected name to get loaded from file")
	}
}

func TestSaveProduct(t *testing.T) {
	// Arrange
	var version1 = "1.0.0"
	var version2 = "2.0.0"
	ym := NewProductModel(file.Types.YAML).(*model)

	// Act
	ym.SaveProduct(&Product{Version: version1})

	// Assert
	if ym.Product.Version != version1 {
		t.Fatalf("Expected version to be set to: %s", version1)
	}
	ym.SaveProduct(&Product{Version: version2})
	if ym.Product.Version != version2 {
		t.Fatalf("Expected version to be changed to: %s", version2)
	}
}

func TestAddApplication(t *testing.T) {
	// Arrange
	m := NewProductModel(file.Types.YAML).(*model)
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

func TestAddDependency(t *testing.T) {
	// Arrange
	m := NewProductModel(file.Types.YAML).(*model)
	var name = "Derp"
	var exists = false

	// Act
	m.AddDependency(&Dependency{
		Name: name,
	})

	// Assert
	for _, p := range m.Product.Dependencies {
		if p.Name == name {
			exists = true
		}
	}
	if !exists {
		t.Fatalf("Expected application to exist in applications list.")
	}
}

func TestWriteFile(t *testing.T) {
	// Arrange
	yamlProductModel := NewProductModel(file.Types.YAML).(*model)
	jsonProductModel := NewProductModel(file.Types.JSON).(*model)

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
	yamlProductModel := NewProductModel(file.Types.YAML).(*model)
	jsonProductModel := NewProductModel(file.Types.JSON).(*model)

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
	yamlProductModel := NewProductModel(file.Types.YAML).(*model)
	jsonProductModel := NewProductModel(file.Types.JSON).(*model)

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
