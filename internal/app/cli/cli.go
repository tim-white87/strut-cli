package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/urfave/cli"
)

// StartCLI - gathers command line args
func StartCLI(args []string) error {
	app := cli.NewApp()
	app.Name = "strut"
	app.Description = "Utility for setting up and running strut based products"

	app.Commands = []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new strut product",
			Category:  "Setup",
			ArgsUsage: "[name]",
			Action:    create,
		},
		{
			Name:      "application",
			Usage:     "Setup product application",
			Category:  "Setup",
			ArgsUsage: "<type> [name] [value]",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Add application",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "provider",
					Usage: "Add application",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:      "run",
			Category:  "Develop",
			Usage:     "Runs defined application commands",
			ArgsUsage: "<cmd> [applications]",
			// TODO lets read the commands from the strut file and set this up dynamically
			Subcommands: []cli.Command{
				{
					Name:  "install",
					Usage: "Runs defined install commands",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "build",
					Usage: "Runs defined build commands",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Runs defined start commands",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:      "provision",
			Category:  "Develop",
			Usage:     "Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers",
			ArgsUsage: "[applications] [providers]",
			Action: func(c *cli.Context) error {
				fmt.Println("first arg: ", c.Args().First())
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(os.Args)
}

func create(c *cli.Context) error {
	if checkForProductFile() {
		return cli.NewExitError("Product file already exists in folder.", 1)
	}
	var pm = product.New(file.Types.YAML)
	var name string
	if c != nil {
		name = c.Args().First()
	}
	if name == "" {
		// TODO prompt for user input to setup other product attributes
	}
	pm.SaveProduct(&product.Product{Name: name})
	return nil
}

func checkForProductFile() bool {
	for _, fileType := range file.TypeList {
		var filePath = fmt.Sprintf("./%s.%s", product.ProductFileName, fileType.Extension)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return true
		}
	}
	return false
}
