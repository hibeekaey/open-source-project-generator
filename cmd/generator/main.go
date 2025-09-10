package main

import (
	"log"
	"os"

	"github.com/open-source-template-generator/internal/app"
	"github.com/open-source-template-generator/internal/container"
)

func main() {
	// Initialize dependency injection container
	c := container.NewContainer()

	// Create and configure the application
	application := app.NewApp(c)

	// Ensure proper cleanup
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("Error closing application: %v", err)
		}
	}()

	// Execute the CLI application
	if err := application.Execute(); err != nil {
		log.Printf("Error executing application: %v", err)
		os.Exit(1)
	}
}
