package product

import (
	"os"
	"testing"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

var name = "Foobar"

func init() {
	os.Chdir("../../../test/testdata")
	// put this teardown maybe? os.Chdir("..")
}

func TestWriteFile(t *testing.T) {
	// Arrange
	yamlProductModel := New(name, file.Types.YAML).(model)
	jsonProductModel := New(name, file.Types.JSON).(model)
	// Act
	yamlProduct, yamlErr := yamlProductModel.writeFile(yamlProductModel)
	jsonProduct, jsonErr := jsonProductModel.writeFile(jsonProductModel)

	// Assert
	if _, err := os.Stat("./strut.yaml"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.yaml file to exist.")
	}
	if _, err := os.Stat("./strut.json"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.json file to exist.")
	}
	if yamlErr != nil {
		t.Fatalf("Expected no error %s", yamlErr)
	}
	if yamlProduct.Name != yamlProductModel.Product.Name {
		t.Fatalf("Expected YAML product name: %s to match match strut file name: %s", yamlProduct.Name, yamlProductModel.Product.Name)
	}
	if jsonErr != nil {
		t.Fatalf("Expected no error %s", jsonErr)
	}
	if jsonProduct.Name != jsonProductModel.Product.Name {
		t.Fatalf("Expected YAML product name: %s to match match strut file name: %s", jsonProduct.Name, jsonProductModel.Product.Name)
	}
}

func TestReadFile(t *testing.T) {
	// Arrange
	yamlProductModel := New(name, file.Types.YAML).(model)
	jsonProductModel := New(name, file.Types.JSON).(model)

	// Act
	yamlProduct, err := yamlProductModel.readFile(yamlProductModel)
	jsonProduct, err := yamlProductModel.readFile(jsonProductModel)

	// Assert
	if err != nil {
		t.Fatalf("Expected strut file to parse to model")
	}
	if yamlProductModel.Product.Name != yamlProduct.Name {
		t.Fatalf("Expected YAML product name: %s to match match strut file name: %s", name, yamlProduct.Name)
	}
	if jsonProductModel.Product.Name != jsonProduct.Name {
		t.Fatalf("Expected YAML product name: %s to match match strut file name: %s", name, jsonProduct.Name)
	}
}
