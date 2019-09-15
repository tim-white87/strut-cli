package cli

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/fatih/color"
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
					Name:   "provider",
					Usage:  "Add provider to application(s)",
					Action: addProvider,
				},
			},
		},
		{
			Name:      "run",
			Category:  "Develop",
			Usage:     "Runs command for applications that have it defined",
			ArgsUsage: "<cmd> [applications]",
			Action:    runCommand,
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

func addProvider(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		return cli.NewExitError(missingFileText, 1)
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	applications := addProviderPrompt(product.Applications)
	pm.UpdateApplications(applications)
	return nil
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
		color.Red("Specify command")
		return cli.NewExitError("", 1)
	}
	for _, app := range product.Applications {
		if app.LocalConfig.Commands == nil {
			continue
		}
		err := os.Chdir(app.LocalConfig.Path)
		if err != nil {
			color.Red("Error >>> app: %s, local path: %s", app.Name, app.LocalConfig.Path)
			color.Red("%s", err)
			return nil
		}
		appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)

		for _, appCmd := range appCmds {
			parts := strings.Fields(appCmd)
			data, err := exec.Command(parts[0], parts[1:]...).Output()
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		}
	}
	return nil
}
