package cli

import (
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/command"
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

var owd, _ = os.Getwd()

func runCommand(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	cmd := c.Args().First()
	if cmd == "" {
		color.Red("Error >>> Specify command")
		return nil
	}
	cmds := buildMapCmds(cmd, product.Applications)
	command.ExecuteGroup(cmds)
	return nil
}

func buildMapCmds(cmd string, apps []*product.Application) []*exec.Cmd {
	cmds := make([]*exec.Cmd, 0)
	for _, app := range apps {
		if app.LocalConfig == nil || app.LocalConfig.Commands == nil {
			continue
		}
		appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)
		for _, appCmd := range appCmds {
			parts := strings.Fields(appCmd)
			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Dir = app.LocalConfig.Path
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}
