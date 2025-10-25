package main

import (
	"flag"
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
		generateCmd := flag.NewFlagSet("generate", flag.ContinueOnError)
		generateCmd.SetOutput(os.Stderr)
		configFile := generateCmd.String("config-file", "", "Path to project configuration file")
		generateCmd.StringVar(configFile, "c", "", "Path to project configuration file (shorthand)")
		generateCmd.Usage = printGenerateHelp

		if err := generateCmd.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}

		if generateCmd.NArg() > 0 {
			fmt.Fprintf(os.Stderr, output.ColorRed+"Error: unexpected argument '%s'\n"+output.ColorReset, generateCmd.Arg(0))
			fmt.Fprintln(os.Stderr, "Run 'generator generate --help' for usage.")
			os.Exit(1)
		}

		runGenerate(*configFile)
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
