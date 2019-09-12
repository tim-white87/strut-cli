package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/urfave/cli"
)

// StartCli - gathers command line args
func StartCli(args []string) error {
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
			Name:      "add",
			Usage:     "Add to the product, must run subcommand for the item to be added",
			Category:  "Setup",
			ArgsUsage: "<type> [name] [value]",
			Subcommands: []cli.Command{
				{
					Name:   "application",
					Usage:  "Setup product application",
					Action: addApplication,
				},
				{
					Name:   "dependency",
					Usage:  "Setup product software dependency, such as git, docker, etc.",
					Action: addDependency,
				},
				{
					Name:  "provider",
					Usage: "Add provider to application(s)",
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

func create(c *cli.Context) error {
	if exists, _ := checkForProductFile(); exists {
		return cli.NewExitError("Product file already exists in folder.", 1)
	}
	name := ""
	if c != nil {
		name = c.Args().First()
	}
	p, ft := createPrompt(name)
	pm := product.NewProductModel(ft)
	pm.SaveProduct(p)
	return nil
}

const missingFileText = "Product file doesn't exist in folder."

func addApplication(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		return cli.NewExitError(missingFileText, 1)
	}
	a := addApplicationPrompt()
	pm := product.NewProductModel(ft)
	pm.LoadProduct()
	pm.AddApplication(a)
	return nil
}

func addDependency(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		return cli.NewExitError(missingFileText, 1)
	}
	d := addDependencyPrompt()
	pm := product.NewProductModel(ft)
	pm.LoadProduct()
	pm.AddDependency(d)
	return nil
}
