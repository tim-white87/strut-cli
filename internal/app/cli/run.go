package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"

	"github.com/cecotw/strut-cli/internal/app/product"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var runCmd = cli.Command{
	Name:      "run",
	Category:  "Develop",
	Usage:     "Runs command for applications that have it defined in local config",
	ArgsUsage: "<cmd> [applications]",
	Action:    runCommand,
}

func runCommand(c *cli.Context) error {
	exists, ft := checkForProductFile()
	if !exists {
		color.Red(missingFileText)
		return nil
	}
	pm := product.NewProductModel(ft)
	product := pm.LoadProduct()
	cmd := c.Args().First()
	if cmd == "" {
		color.Red("Error >>> Specify command")
		return nil
	}
	owd, _ := os.Getwd()
	for _, app := range product.Applications {
		if app.LocalConfig == nil || app.LocalConfig.Commands == nil {
			continue
		}
		os.Chdir(owd)
		err := os.Chdir(app.LocalConfig.Path)
		if err != nil {
			color.Red("Error >>> app: %s, local path: %s", app.Name, app.LocalConfig.Path)
			color.Red(err.Error())
			return nil
		}
		appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)
		wg := &sync.WaitGroup{}
		wg.Add(len(appCmds))
		defer wg.Wait()
		for _, appCmd := range appCmds {
			go execute(appCmd, wg)
		}
	}
	return nil
}

func execute(appCmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	parts := strings.Fields(appCmd)
	command := exec.Command(parts[0], parts[1:]...)
	stdout, _ := command.StdoutPipe()
	err := command.Start()
	if err != nil {
		color.Red(err.Error())
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	command.Wait()
}
