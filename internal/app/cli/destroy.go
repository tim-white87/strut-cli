package cli

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/provider"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var destroyCmd = cli.Command{
	Name:      "destroy",
	Category:  "Cloud",
	Usage:     "Destroys the defined infrastructure. Destroy commands run concurrently in reverse priority from provisioning.",
	ArgsUsage: "[applications] [providers]",
	Action:    destroy,
}

func destroy(c *cli.Context) error {
	color.HiBlack("Destroying...")
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
	}
	if confrimDestroyPrompt() {
		pm := product.NewProductModel(ft)
		product := pm.LoadProduct()
		provider.Destroy(product.ProvisionMap)
	}
	return nil
}

func confrimDestroyPrompt() bool {
	reallyNuke := false
	survey.AskOne(&survey.Confirm{Message: "Are you sure? This will nuke your shit."}, reallyNuke)
	return reallyNuke
}
