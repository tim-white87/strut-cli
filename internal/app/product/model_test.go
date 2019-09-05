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
	t.Fatalf("Expected implementation")
}

func TestAddApplication(t *testing.T) {
	t.Fatalf("Expected implementation")
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
