package cli

import (
	"os"
	"testing"
)

const TestDataFolder = "../../../test/testdata"

func init() {
	os.Chdir(TestDataFolder)
	// put this teardown maybe? os.Chdir("..")
}

func TestCreate(t *testing.T) {
	t.Fatalf("Expected implementation")
}

func TestCheckForProductFile(t *testing.T) {
	exists := checkForProductFile()

	if !exists {
		t.Fatalf("Expected true if product file exists.")
	}

	os.Chdir("../")
	exists = checkForProductFile()

	if exists {
		t.Fatalf("Expected false if there are no product files.")
	}

}
