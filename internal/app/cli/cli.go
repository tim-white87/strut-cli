package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	createCmd,
	addCmd,
	runCmd,
	provisionCmd,
	destroyCmd,
	cloneCmd,
}

// StartCli - gathers command line args
func StartCli(args []string) error {
	app := cli.NewApp()
	app.Name = "strut"
	app.Usage = "cli for managing strut apps"
	app.Description = "Utility for setting up and running strut based products"
	app.Commands = commands
	sort.Sort(cli.CommandsByName(app.Commands))
	return app.Run(args)
}

func checkForProductFile() (bool, *file.Type) {
	for _, fileType := range file.TypeList {
		var filePath = fmt.Sprintf("./%s.%s", product.ProductFileName, fileType.Extension)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return true, fileType
		}
	}
	return false, nil
}

const missingFileText = "Product file doesn't exist in folder."
