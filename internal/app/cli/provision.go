package cli

import (
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/provider"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var provisionCmd = cli.Command{
	Name:      "provision",
	Category:  "Cloud",
	Usage:     "Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers. Provision commands run concurrently. Specify a priority integer on dependencies to batch. The lowest number indicates highest priority. Resources with undefined priorities will run last in the final batch. Note - items with the same priority number will still be concurrent.",
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

	provider.ProvisionResources(product.ProvisionMap)
	return nil
}
