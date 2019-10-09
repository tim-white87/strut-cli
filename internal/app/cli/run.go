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
	mapCmds := buildMapCmds(cmd, product.Applications)
	command.SpawnMapGroup(mapCmds)
	return nil
}

func buildMapCmds(cmd string, apps []*product.Application) map[string][]*exec.Cmd {
	mapCmds := make(map[string][]*exec.Cmd)
	for _, app := range apps {
		cmds := make([]*exec.Cmd, 0)
		if app.LocalConfig == nil || app.LocalConfig.Commands == nil {
			continue
		}
		appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)
		for _, appCmd := range appCmds {
			parts := strings.Fields(appCmd)
			args := make([]string, 0)
			for i := 0; i < len(parts); {
				arg := parts[i]
				if strings.HasPrefix(arg, "'") && strings.HasSuffix(parts[i+1], "'") {
					arg = strings.Replace(arg, "'", "", 1) + " " + strings.Replace(parts[i+1], "'", "", 1)
					args = append(args, arg)
					i += 2
				} else {
					args = append(args, arg)
					i++
				}
			}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = app.LocalConfig.Path
			cmds = append(cmds, cmd)
		}
		mapCmds[app.Name] = cmds
	}
	return mapCmds
}
