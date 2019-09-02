package product

import (
	"os"
	"testing"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

func TestCreate(t *testing.T) {
	// Arrange
	os.Chdir("../../../test/testdata")
	productModel := New()

	// Act
	productModel.CreateFile(file.Types.YAML)

	// Assert
	if _, err := os.Stat("./strut.yaml"); os.IsNotExist(err) {
		t.Fatalf("Expected ./strut.yaml file to exist")
	}
	os.Chdir("..")
}

func TestRead(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestAddApplication(t *testing.T) {}
