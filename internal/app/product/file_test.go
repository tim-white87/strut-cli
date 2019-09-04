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

func TestCreate(t *testing.T) {
	// Arrange
	yamlProductModel := New(name, file.Types.YAML).(model)
	jsonProductModel := New(name, file.Types.JSON).(model)
	// Act
	yamlProductModel.CreateFile(yamlProductModel)
	jsonProductModel.CreateFile(jsonProductModel)

	// Assert
	if _, err := os.Stat("./strut.yaml"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.yaml file to exist.")
	}
	if _, err := os.Stat("./strut.json"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.json file to exist.")
	}
}

func TestRead(t *testing.T) {
	// Arrange
	yamlProductModel := New(name, file.Types.YAML).(model)
	jsonProductModel := New(name, file.Types.JSON).(model)

	// Act
	yamlProduct, err := yamlProductModel.ReadFile(yamlProductModel)
	jsonProduct, err := yamlProductModel.ReadFile(jsonProductModel)

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

func TestUpdate(t *testing.T) {
	// Arrange
	yamlProductModel := New(name, file.Types.YAML).(model)
	jsonProductModel := New(name, file.Types.JSON).(model)

	// Act
	yamlProductModel.Product.Version = "2.0.0"
	yamlProduct, err := yamlProductModel.UpdateFile(yamlProductModel)

	jsonProductModel.Product.Version = "2.0.0"
	jsonProduct, err := jsonProductModel.UpdateFile(jsonProductModel)

	// Assert
	if err != nil {
		t.Fatalf("Expected strut file to parse to model")
	}
	if yamlProduct.Version != yamlProductModel.Product.Version {
		t.Fatalf("Expected YAML product version update for version to change from %s to %s", yamlProduct.Version, yamlProductModel.Product.Version)
	}
	if jsonProduct.Version != jsonProductModel.Product.Version {
		t.Fatalf("Expected JSON product version update for version to change from %s to %s", jsonProduct.Version, jsonProductModel.Product.Version)
	}
}

func TestAddApplication(t *testing.T) {
	t.Fatalf("Expected implementation")
}
