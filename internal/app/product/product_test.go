package product

import (
	"testing"

	"github.com/cecotw/strut-cli/internal/pkg/file"
)

func TestCreate(t *testing.T) {
	// Arrange
	productModel := New()

	// Act
	productModel.CreateFile(file.Types.YAML)

	// Assert
	// if os.
}

func TestRead(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestAddApplication(t *testing.T) {}
