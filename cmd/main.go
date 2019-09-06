package main

import (
	"log"
	"os"

	"github.com/cecotw/strut-cli/internal/app/cli"

	"github.com/fatih/color"
)

func main() {
	color.Blue("Welcome to Strut!")
	cliErr := cli.StartCLI(os.Args)
	if cliErr != nil {
		log.Fatal(cliErr)
	}
}
