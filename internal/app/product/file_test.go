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
	productModel := New(name, file.Types.YAML).(model)

	// Act
	product, err := productModel.ReadFile()

	// Assert
	if err != nil {
		t.Fatalf("Expected strut file to parse to model")
	}
	if productModel.Name != product.Name {
		t.Fatalf("Expected product name: %s to match match strut file name: %s", name, product.Name)
	}
}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestAddApplication(t *testing.T) {}
