// Package main provides the entry point for the Open Source Template Generator CLI.
//
// The generator is a comprehensive tool for creating production-ready, enterprise-grade
// open source project structures following modern best practices. It supports multiple
// technology stacks including Go 1.24+, Node.js 20+, Next.js 15+, React 19+, and more.
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

	"github.com/open-source-template-generator/internal/app"
	"github.com/open-source-template-generator/internal/container"
)

// Version information set by build
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// main is the entry point for the Open Source Template Generator CLI application.
// It initializes the dependency injection container, creates the application instance,
// and handles proper cleanup and error reporting.
func main() {
	// Initialize dependency injection container with all required services
	// including CLI handlers, template processors, and validation engines
	c := container.NewContainer()

	// Create and configure the application with the initialized container
	// This sets up all CLI commands, flags, and validation logic
	application := app.NewAppWithVersion(c, Version, GitCommit, BuildTime)

	// Ensure proper cleanup of resources when the application exits
	// This includes closing file handles, network connections, and temporary files
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("Error closing application: %v", err)
		}
	}()

	// Execute the CLI application with command-line arguments
	// This processes user input, validates configuration, and generates projects
	if err := application.Execute(); err != nil {
		log.Printf("Error executing application: %v", err)
		os.Exit(1)
	}
}
