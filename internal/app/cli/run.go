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

var owd, _ = os.Getwd()

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

	appWg := &sync.WaitGroup{}
	appWg.Add(len(product.Applications))
	defer appWg.Wait()
	for _, app := range product.Applications {
		if app.LocalConfig == nil || app.LocalConfig.Commands == nil {
			continue
		}
		go runAppCommands(app, cmd, appWg)
	}
	return nil
}

func runAppCommands(app *product.Application, cmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	appCmds := reflect.ValueOf(app.LocalConfig.Commands).Elem().FieldByName(strings.Title(cmd)).Interface().([]string)
	cwg := &sync.WaitGroup{}
	cwg.Add(len(appCmds))
	defer cwg.Wait()
	for _, appCmd := range appCmds {
		go Execute(app.LocalConfig.Path, appCmd, cwg)
	}
}

// Execute executes a command and logs with scanner
func Execute(path string, cmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	parts := strings.Fields(cmd)
	command := exec.Command(parts[0], parts[1:]...)
	command.Dir = path
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()
	err := command.Start()
	if err != nil {
		color.Red(err.Error())
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		color.Red(scanner.Text())
	}
	command.Wait()
}
