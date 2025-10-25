package main

import (
	"fmt"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

func printUsage() {
	fmt.Println(output.ColorCyan + "generator" + output.ColorReset + " - Open Source Project Generator")
	fmt.Println()
	fmt.Println(output.ColorYellow + "USAGE" + output.ColorReset)
	fmt.Println("  generator <command> [options]")
	fmt.Println()
	fmt.Println(output.ColorYellow + "COMMANDS" + output.ColorReset)
	fmt.Println("  generate    Generate a new project (interactive or from config file)")
	fmt.Println()
	fmt.Println(output.ColorYellow + "OPTIONS" + output.ColorReset)
	fmt.Println("  -h, --help     Show help")
	fmt.Println("  -v, --version  Show version")
	fmt.Println()
	fmt.Println("Run 'generator help <command>' for more information on a command.")
}

func printHelp() {
	printUsage()
	fmt.Println()
	fmt.Println(output.ColorYellow + "EXAMPLES" + output.ColorReset)
	fmt.Println("  generator generate")
	fmt.Println("  generator --version")
	fmt.Println("  generator help generate")
}

func printGenerateHelp() {
	fmt.Println(output.ColorCyan + "generator generate" + output.ColorReset)
	fmt.Println()
	fmt.Println("Generate a new project with interactive prompts.")
	fmt.Println()
	fmt.Println(output.ColorYellow + "USAGE" + output.ColorReset)
	fmt.Println("  generator generate [flags]")
	fmt.Println()
	fmt.Println(output.ColorYellow + "FLAGS" + output.ColorReset)
	fmt.Println("  -c, --config-file <path>    Path to project configuration file")
	fmt.Println()
	fmt.Println(output.ColorYellow + "DESCRIPTION" + output.ColorReset)
	fmt.Println("  Without flags, this command starts an interactive wizard that guides")
	fmt.Println("  you through creating a new project. You'll be prompted to:")
	fmt.Println("    - Enter a project name")
	fmt.Println("    - Select components to create")
	fmt.Println("    - Choose frontend apps to generate")
	fmt.Println("    - Specify an output directory")
	fmt.Println()
	fmt.Println("  With --config-file, the project is generated from a YAML configuration")
	fmt.Println("  file. See configs/project.yaml for the format.")
	fmt.Println()
	fmt.Println(output.ColorYellow + "EXAMPLES" + output.ColorReset)
	fmt.Println("  # Interactive mode")
	fmt.Println("  generator generate")
	fmt.Println()
	fmt.Println("  # Using config file")
	fmt.Println("  generator generate --config-file project.yaml")
	fmt.Println("  generator generate -c project.yaml")
}

func printCommandHelp(command string) {
	switch command {
	case "generate":
		printGenerateHelp()
	default:
		fmt.Fprintf(os.Stderr, output.ColorRed+"Error: unknown command '%s'\n"+output.ColorReset, command)
		fmt.Fprintln(os.Stderr, "Run 'generator --help' for usage.")
		os.Exit(1)
	}
}
