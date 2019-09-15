package cli

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/cecotw/strut-cli/internal/pkg/provider"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var addCmd = cli.Command{
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
}

func addApplication(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
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

func addApplicationPrompt() *product.Application {
	color.Yellow("Lets add an application to your product.")
	answers := &product.Application{
		Version:     "0.0.0",
		LocalConfig: &product.LocalConfig{Commands: &product.LocalCommandRegistry{}},
	}
	err := survey.AskOne(&survey.Input{
		Message: "Enter application name:",
	}, &answers.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
	hasRepo := false
	err = survey.AskOne(&survey.Confirm{Message: "Include Repo?"}, hasRepo)
	if err != nil {
		fmt.Println(err.Error())
	}
	if hasRepo {
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
	err = survey.AskOne(&survey.Input{
		Message: "Provide the local path to the application code:",
	}, &answers.LocalConfig.Path)
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
	selectedApp := selectApplicationPrompt(applications)

	var providerOptions []string
	for _, providerType := range provider.TypeList {
		providerOptions = append(providerOptions, providerType.Name)
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
			if hasProvider {
				color.Red("%s: already has provider: %s", app.Name, selectedProviderType.Name)
			} else {
				app.Providers = append(app.Providers, &provider.Provider{
					Type: selectedProviderType,
				})
				break
			}
		}
	}

	return applications
}

func selectApplicationPrompt(applications []*product.Application) *product.Application {
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
