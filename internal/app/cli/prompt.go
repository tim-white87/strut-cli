package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
)

func createPrompt(name string) (*product.Product, *file.Type) {
	prompt := []*survey.Question{}
	if name == "" {
		prompt = append(prompt, &survey.Question{
			Name:   "name",
			Prompt: &survey.Input{Message: "Enter new strut product name:"},
		})
	}

	prompt = append(prompt, &survey.Question{
		Name: "extension",
		Prompt: &survey.Select{
			Message: "Select file type:",
			Options: []string{
				file.Types.YAML.Extension,
				file.Types.JSON.Extension,
			},
			Default: file.Types.YAML.Extension,
		},
	})
	answers := struct {
		Name      string
		Extension string
	}{}
	err := survey.Ask(prompt, answers)
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
	return &product.Product{Name: answers.Name}, ft
}
