package product

import (
	"os"
	"testing"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

func init() {
	os.Chdir("../../../test/testdata")
	// put this teardown maybe? os.Chdir("..")
}

func TestCreate(t *testing.T) {
	// Arrange
	productModel := New()

	// Act
	productModel.CreateFile(productModel, file.Types.YAML)

	// Assert
	if _, err := os.Stat("./strut.yaml"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.yaml file to exist.")
	}
}

func TestRead(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestAddApplication(t *testing.T) {}
