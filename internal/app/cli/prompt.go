package cli

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/cecotw/strut-cli/internal/pkg/provider"
	"github.com/fatih/color"
)

func createPrompt(name string) (*product.Product, *file.Type) {
	answers := struct {
		*product.Product
		Extension string
	}{}
	prompt := []*survey.Question{}
	if name == "" {
		prompt = append(prompt, &survey.Question{
			Name:   "name",
			Prompt: &survey.Input{Message: "Enter new strut product name:"},
		})
	} else {
		answers.Name = name
	}

	prompt = append(prompt, []*survey.Question{
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
	}...)

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
		Name:        answers.Name,
		Description: answers.Description,
	}, ft
}

func addApplicationPrompt() *product.Application {
	color.Yellow("Lets add an application to your product.")
	answers := &product.Application{
		Version:     "0.0.0",
		LocalConfig: &product.LocalConfig{},
	}
	err := survey.Ask([]*survey.Question{
		{
			Name:   "name",
			Prompt: &survey.Input{Message: "Enter application name:"},
		},
	}, answers)
	if err != nil {
		fmt.Println(err.Error())
	}
	include := struct{ repo bool }{}
	err = survey.Ask([]*survey.Question{
		{
			Name:   "repo",
			Prompt: &survey.Confirm{Message: "Include Repo?"},
		},
	}, include)
	if err != nil {
		fmt.Println(err.Error())
	}
	if include.repo {
		err = survey.Ask([]*survey.Question{
			{
				Name:   "url",
				Prompt: &survey.Input{Message: "Provide the remote URL to the app code:"},
			},
			{
				Name: "type",
				Prompt: &survey.Select{
					Message: "Please select your VCS:",
					Options: []string{"git", "SVN", "mercurial"},
				},
			},
		}, answers.Repository)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	err = survey.Ask([]*survey.Question{
		{
			Name:   "path",
			Prompt: &survey.Input{Message: "Provide the local path to the application code:"},
		},
	}, &answers.LocalConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	return answers
}

func addDependencyPrompt() *product.Dependency {
	color.Yellow("Lets add a product software dependency.")
	answers := struct {
		Name    string
		Install string
	}{}
	err := survey.Ask([]*survey.Question{
		{
			Name:   "name",
			Prompt: &survey.Input{Message: "What is the name of this dependency?"},
		},
		{
			Name:   "install",
			Prompt: &survey.Input{Message: "Please add the install commands in a comma separated list:"},
		},
	}, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &product.Dependency{
		Name:    answers.Name,
		Install: strings.Split(answers.Install, ","),
	}
}

func addProviderPrompt(applications []*product.Application) []*product.Application {
	color.Yellow("Well get a provider added to an applicationn")
	selectedApp := selectApplication(applications)

	var providerOptions []string
	for _, provider := range provider.TypeList {
		providerOptions = append(providerOptions, provider.Name)
	}
	var providerIndex int
	prompt := &survey.Select{
		Message: "Select provider type:",
		Options: providerOptions,
	}
	err := survey.AskOne(prompt, &providerIndex)
	if err != nil {
		fmt.Println(err.Error())
	}
	selectedProviderType := provider.TypeList[providerIndex]
	for _, app := range applications {
		if app.Name == selectedApp.Name {
			hasProvider := false
			for _, provider := range app.Providers {
				if provider.Type.Name == selectedProviderType.Name {
					hasProvider = true
					break
				}
			}
			if !hasProvider {
				color.Red("%s: already has provider: %s", app.Name, selectedProviderType.Name)
			} else {
				app.Providers = append(app.Providers, &product.Provider{
					Type: selectedProviderType,
				})
				break
			}
		}
	}

	return applications
}

func selectApplication(applications []*product.Application) *product.Application {
	var appIndex int
	var appOptions []string
	for _, app := range applications {
		appOptions = append(appOptions, app.Name)
	}
	prompt := &survey.Select{
		Message: "Which application should this be added to?",
		Options: appOptions,
	}
	err := survey.AskOne(prompt, appIndex)
	if err != nil {
		fmt.Println(err.Error())
	}
	return applications[appIndex]
}
