package cli

import (
	"os"
	"testing"

	"github.com/urfave/cli"
)

const TestDataFolder = "../../../test/testdata"

func init() {
	os.Chdir(TestDataFolder)
	// put this teardown maybe? os.Chdir("..")
}

func TestCreate(t *testing.T) {
	err := create(new(cli.Context))

	if err == nil {
		t.Fatalf("Expected exit error since if there already is a product file.")
	}
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
