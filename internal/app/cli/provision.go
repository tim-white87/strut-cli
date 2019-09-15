package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

var provisionCmd = cli.Command{
	Name:      "provision",
	Category:  "Develop",
	Usage:     "Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers",
	ArgsUsage: "[applications] [providers]",
	Action: func(c *cli.Context) error {
		fmt.Println("first arg: ", c.Args().First())
		return nil
	},
}
