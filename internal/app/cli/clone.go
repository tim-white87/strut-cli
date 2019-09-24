package cli

import (
	"fmt"
	"os/exec"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/command"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var cloneCmd = cli.Command{
	Name:      "clone",
	Category:  "Setup",
	Usage:     "Clones the specified applications to their respective local config paths",
	ArgsUsage: "[applications]",
	Action:    clone,
}

func clone(c *cli.Context) error {
	color.HiBlack("Cloning...")
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	cmds := make([]*exec.Cmd, 0)
	for _, app := range product.Applications {
		if app.Repository != nil {
			fmt.Println(app.Repository.URL)
			cmd := exec.Command("git", "clone", app.Repository.URL, app.LocalConfig.Path)
			cmd.Dir = app.LocalConfig.Path
			cmds = append(cmds, cmd)
		}
	}
	command.SpawnGroup(cmds)
	return nil
}
