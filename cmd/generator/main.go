// Package main provides the entry point for the Open Source Project Generator CLI.
//
// The generator is a comprehensive tool for creating production-ready, enterprise-grade
// open source project structures following modern best practices. It supports multiple
// technology stacks including Go 1.25+, Node.js 20+, Next.js 15+, React 19+, and more.
//
// Usage:
//
//	generator generate                    # Interactive project generation
//	generator generate --config file.yml # Generate from configuration
//	generator validate [path]            # Validate project structure
//	generator audit [path]               # Audit existing codebase
//	generator version                    # Show version information
//
// The generator follows clean architecture principles with dependency injection
// and comprehensive error handling to ensure reliable project generation.
package main

import (
	"log"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/internal/app"
)

// Version information set by build
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// main is the entry point for the Open Source Project Generator CLI application.
// It initializes the dependency injection container, creates the application instance,
// and handles proper cleanup and error reporting.
func main() {
	// Create and configure the application
	// This sets up all CLI commands, flags, and comprehensive functionality
	application, err := app.NewApp(Version, GitCommit, BuildTime)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Execute the CLI application with command-line arguments
	// This processes user input, validates configuration, and generates projects
	if err := application.Run(os.Args[1:]); err != nil {
		// Only log detailed error in verbose mode, otherwise just exit
		// The CLI already printed a user-friendly error message
		os.Exit(1)
	}
}
