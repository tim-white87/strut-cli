package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/urfave/cli"
)

var createCmd = cli.Command{
	Name:      "create",
	Usage:     "Create a new strut product",
	Category:  "Setup",
	ArgsUsage: "[name]",
	Action:    create,
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

func createPrompt(name string) (*product.Product, *file.Type) {
	if name == "" {
		survey.AskOne(&survey.Input{
			Message: "Enter new strut product name:",
		}, name)
	}

	answers := struct {
		Description string
		Extension   string
	}{}
	prompt := []*survey.Question{
		{
			Name:   "description",
			Prompt: &survey.Multiline{Message: "Enter new product description:"},
		},
		{
			Name: "extension",
			Prompt: &survey.Select{
				Message: "Select file type:",
				Options: []string{
					file.Types.YAML.Extension,
					file.Types.JSON.Extension,
				},
				Default: file.Types.YAML.Extension,
			},
		},
	}

	err := survey.Ask(prompt, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}
	ft := file.Types.YAML
	for _, fileType := range file.TypeList {
		if fileType.Extension == answers.Extension {
			ft = fileType
			break
		}
	}
	return &product.Product{
		Name:        name,
		Description: answers.Description,
	}, ft
}
