package main

import (
	"fmt"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		runGenerate()
	case "help", "--help", "-h":
		if len(os.Args) > 2 {
			printCommandHelp(os.Args[2])
		} else {
			printHelp()
		}
	case "version", "--version", "-v":
		fmt.Println(Version)
	default:
		fmt.Fprintf(os.Stderr, output.ColorRed+"Error: unknown command '%s'\n"+output.ColorReset, command)
		fmt.Fprintln(os.Stderr, "Run 'generator --help' for usage.")
		os.Exit(1)
	}
}
