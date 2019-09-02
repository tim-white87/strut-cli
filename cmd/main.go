package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var productModel = new(product.Model)

func main() {
	color.Blue("Welcome to Strut!")
	app := cli.NewApp()
	app.Name = "strut"
	app.Description = "Utility for setting up and running strut based products"

	app.Commands = []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new strut product",
			ArgsUsage: "[name]",
			Action: func(c *cli.Context) error {
				fmt.Println("first arg: ", c.Args().First())
				return nil
			},
		},
		{
			Name:      "add",
			Usage:     "Add an <application|provider> to the product",
			ArgsUsage: "<type> [name] [value]",
			Action: func(c *cli.Context) error {
				fmt.Println("first arg: ", c.Args().First())
				return nil
			},
		},
		{
			Name:      "run",
			Usage:     "Runs the specified command <install|build|start> for the product applications (separated with a comma). Default will run all apps.",
			ArgsUsage: "<cmd> [applications]",
			Action: func(c *cli.Context) error {
				fmt.Println("first arg: ", c.Args().First())
				return nil
			},
		},
		{
			Name:      "provision",
			Usage:     "Provisions the defined infrastructure for the applications to the specified provider. Defaults to all applications deployed to all providers",
			ArgsUsage: "[applications] [providers]",
			Action: func(c *cli.Context) error {
				fmt.Println("first arg: ", c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
