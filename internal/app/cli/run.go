package cli

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var runCmd = cli.Command{
	Name:      "run",
	Category:  "Develop",
	Usage:     "Runs command for applications that have it defined in local config",
	ArgsUsage: "<cmd> [applications]",
	Action:    runCommand,
}

func runCommand(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		return cli.NewExitError(missingFileText, 1)
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	cmd := c.Args().First()
	if cmd == "" {
		color.Red("Error >>> Specify command")
		return nil
	}
	for _, app := range product.Applications {
		if app.LocalConfig.Commands == nil {
			continue
		}
		err := os.Chdir(app.LocalConfig.Path)
		if err != nil {
			color.Red("Error >>> app: %s, local path: %s", app.Name, app.LocalConfig.Path)
			color.Red(err.Error())
			return nil
		}
		appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)

		for _, appCmd := range appCmds {
			parts := strings.Fields(appCmd)
			data, err := exec.Command(parts[0], parts[1:]...).Output()
			if err != nil {
				color.Red(err.Error())
				return nil
			}
			fmt.Println(string(data))
		}
	}
	return nil
}
