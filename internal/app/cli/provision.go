package cli

import (
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var provisionCmd = cli.Command{
	Name:      "provision",
	Category:  "Cloud",
	Usage:     "Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers",
	ArgsUsage: "[applications] [providers]",
	Action:    provision,
}

func provision(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	for _, app := range product.Applications {
		color.Green(app.Name)
	}
	return nil
}
