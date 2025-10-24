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
	fmt.Println("  generate    Generate a new project interactively")
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

func printCommandHelp(command string) {
	switch command {
	case "generate":
		fmt.Println(output.ColorCyan + "generator generate" + output.ColorReset)
		fmt.Println()
		fmt.Println("Generate a new project with interactive prompts.")
		fmt.Println()
		fmt.Println(output.ColorYellow + "USAGE" + output.ColorReset)
		fmt.Println("  generator generate")
		fmt.Println()
		fmt.Println(output.ColorYellow + "DESCRIPTION" + output.ColorReset)
		fmt.Println("  This command starts an interactive wizard that guides you through")
		fmt.Println("  creating a new project. You'll be prompted to:")
		fmt.Println("    - Enter a project name")
		fmt.Println("    - Select folders to create")
		fmt.Println("    - Choose Next.js apps to generate")
		fmt.Println("    - Specify an output directory")
	default:
		fmt.Fprintf(os.Stderr, output.ColorRed+"Error: unknown command '%s'\n"+output.ColorReset, command)
		os.Exit(1)
	}
}
