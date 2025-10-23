package orchestrator

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// IntegrationManager handles post-generation integration of components
type IntegrationManager struct {
	projectRoot string
	verbose     bool
}

// NewIntegrationManager creates a new integration manager
func NewIntegrationManager(projectRoot string, verbose bool) *IntegrationManager {
	return &IntegrationManager{
		projectRoot: projectRoot,
		verbose:     verbose,
	}
}

// Integrate configures generated components to work together
func (im *IntegrationManager) Integrate(ctx context.Context, components []*models.Component, config *models.IntegrationConfig) error {
	if config.GenerateDockerCompose {
		if _, err := im.GenerateDockerCompose(components); err != nil {
			return fmt.Errorf("failed to generate docker-compose: %w", err)
		}
	}

	if err := im.ConfigureEnvironment(components, config); err != nil {
		return fmt.Errorf("failed to configure environment: %w", err)
	}

	if config.GenerateScripts {
		if err := im.GenerateScripts(components); err != nil {
			return fmt.Errorf("failed to generate scripts: %w", err)
		}
	}

	if err := im.GenerateDocumentation(components, config); err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	return nil
}

// GenerateDockerCompose creates a Docker Compose file for all components
func (im *IntegrationManager) GenerateDockerCompose(components []*models.Component) (string, error) {
	var services []string
	var networks []string
	var volumes []string

	// Track unique networks and volumes
	networkSet := make(map[string]bool)
	volumeSet := make(map[string]bool)

	for _, comp := range components {
		service, err := im.generateServiceDefinition(comp)
		if err != nil {
			return "", fmt.Errorf("failed to generate service for %s: %w", comp.Name, err)
		}
		services = append(services, service)

		// Add default network
		networkSet["app-network"] = true

		// Add component-specific volumes
		volumeName := fmt.Sprintf("%s-data", comp.Name)
		volumeSet[volumeName] = true
	}

	// Convert sets to slices
	for network := range networkSet {
		networks = append(networks, fmt.Sprintf("  %s:\n    driver: bridge", network))
	}
	for volume := range volumeSet {
		volumes = append(volumes, fmt.Sprintf("  %s:", volume))
	}

	// Build complete docker-compose.yml
	compose := fmt.Sprintf(`version: '3.8'

services:
%s

networks:
%s

volumes:
%s
`,
		strings.Join(services, "\n\n"),
		strings.Join(networks, "\n"),
		strings.Join(volumes, "\n"))

	return compose, nil
}

// generateServiceDefinition creates a Docker service definition for a component
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateServiceDefinition(comp *models.Component) (string, error) {
	switch comp.Type {
	case "nextjs":
		return im.generateNextJSService(comp), nil
	case "go-backend":
		return im.generateGoBackendService(comp), nil
	case "android", "ios":
		// Mobile apps don't typically run in Docker during development
		return "", nil
	default:
		return im.generateGenericService(comp), nil
	}
}

// generateNextJSService creates a Docker service for Next.js
func (im *IntegrationManager) generateNextJSService(comp *models.Component) string {
	serviceName := sanitizeServiceName(comp.Name)
	relPath := filepath.Base(comp.Path)

	return fmt.Sprintf(`  %s:
    build:
      context: ./%s
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development
      - NEXT_PUBLIC_API_URL=${API_URL:-http://backend:8080}
    volumes:
      - ./%s:/app
      - /app/node_modules
      - /app/.next
    networks:
      - app-network
    depends_on:
      - backend
    command: npm run dev`, serviceName, relPath, relPath)
}

// generateGoBackendService creates a Docker service for Go backend
func (im *IntegrationManager) generateGoBackendService(comp *models.Component) string {
	serviceName := sanitizeServiceName(comp.Name)
	relPath := filepath.Base(comp.Path)

	return fmt.Sprintf(`  %s:
    build:
      context: ./%s
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - GO_ENV=development
      - PORT=8080
      - DATABASE_URL=${DATABASE_URL:-}
    volumes:
      - ./%s:/app
      - %s-data:/data
    networks:
      - app-network
    command: go run main.go`, serviceName, relPath, relPath, serviceName)
}

// generateGenericService creates a generic Docker service
func (im *IntegrationManager) generateGenericService(comp *models.Component) string {
	serviceName := sanitizeServiceName(comp.Name)
	relPath := filepath.Base(comp.Path)

	return fmt.Sprintf(`  %s:
    build:
      context: ./%s
      dockerfile: Dockerfile
    volumes:
      - ./%s:/app
    networks:
      - app-network`, serviceName, relPath, relPath)
}

// ConfigureEnvironment generates shared environment configuration
func (im *IntegrationManager) ConfigureEnvironment(components []*models.Component, config *models.IntegrationConfig) error {
	// Generate root .env file
	rootEnv, err := im.generateRootEnvFile(components, config)
	if err != nil {
		return fmt.Errorf("failed to generate root .env: %w", err)
	}

	// Store the generated env content (in real implementation, this would write to file)
	if im.verbose {
		fmt.Printf("Generated root .env file:\n%s\n", rootEnv)
	}

	// Generate component-specific environment configurations
	for _, comp := range components {
		if err := im.configureComponentEnvironment(comp, config); err != nil {
			return fmt.Errorf("failed to configure environment for %s: %w", comp.Name, err)
		}
	}

	return nil
}

// generateRootEnvFile creates a shared .env file with common configuration
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateRootEnvFile(components []*models.Component, config *models.IntegrationConfig) (string, error) {
	var envVars []string

	// Add header comment
	envVars = append(envVars, "# Shared Environment Configuration")
	envVars = append(envVars, "# Generated by Open Source Project Generator")
	envVars = append(envVars, "")

	// Add API endpoints from config
	if len(config.APIEndpoints) > 0 {
		envVars = append(envVars, "# API Endpoints")
		for key, value := range config.APIEndpoints {
			envKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_")) + "_URL"
			envVars = append(envVars, fmt.Sprintf("%s=%s", envKey, value))
		}
		envVars = append(envVars, "")
	}

	// Add shared environment variables from config
	if len(config.SharedEnvironment) > 0 {
		envVars = append(envVars, "# Shared Environment Variables")
		for key, value := range config.SharedEnvironment {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
		envVars = append(envVars, "")
	}

	// Add component-specific environment variables
	envVars = append(envVars, "# Component Configuration")
	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			envVars = append(envVars, fmt.Sprintf("# %s (Next.js)", comp.Name))
			envVars = append(envVars, "NEXT_PUBLIC_API_URL=http://localhost:8080")
			envVars = append(envVars, "NODE_ENV=development")
		case "go-backend":
			envVars = append(envVars, fmt.Sprintf("# %s (Go Backend)", comp.Name))
			envVars = append(envVars, "PORT=8080")
			envVars = append(envVars, "GO_ENV=development")
			envVars = append(envVars, "DATABASE_URL=")
		}
		envVars = append(envVars, "")
	}

	return strings.Join(envVars, "\n"), nil
}

// configureComponentEnvironment creates component-specific environment files
func (im *IntegrationManager) configureComponentEnvironment(comp *models.Component, config *models.IntegrationConfig) error {
	switch comp.Type {
	case "nextjs":
		return im.configureNextJSEnvironment(comp, config)
	case "go-backend":
		return im.configureGoBackendEnvironment(comp, config)
	default:
		// No specific environment configuration needed
		return nil
	}
}

// configureNextJSEnvironment creates Next.js-specific environment configuration
func (im *IntegrationManager) configureNextJSEnvironment(comp *models.Component, config *models.IntegrationConfig) error {
	var envVars []string

	envVars = append(envVars, "# Next.js Environment Configuration")
	envVars = append(envVars, fmt.Sprintf("# Component: %s", comp.Name))
	envVars = append(envVars, "")

	// Add API URL
	apiURL := "http://localhost:8080"
	if url, ok := config.APIEndpoints["backend"]; ok {
		apiURL = url
	}
	envVars = append(envVars, fmt.Sprintf("NEXT_PUBLIC_API_URL=%s", apiURL))

	// Add common Next.js variables
	envVars = append(envVars, "NODE_ENV=development")
	envVars = append(envVars, "")

	// Add any shared environment variables
	for key, value := range config.SharedEnvironment {
		if strings.HasPrefix(key, "NEXT_PUBLIC_") {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
	}

	content := strings.Join(envVars, "\n")
	if im.verbose {
		fmt.Printf("Generated .env.local for %s:\n%s\n", comp.Name, content)
	}

	return nil
}

// configureGoBackendEnvironment creates Go backend-specific environment configuration
func (im *IntegrationManager) configureGoBackendEnvironment(comp *models.Component, config *models.IntegrationConfig) error {
	var envVars []string

	envVars = append(envVars, "# Go Backend Environment Configuration")
	envVars = append(envVars, fmt.Sprintf("# Component: %s", comp.Name))
	envVars = append(envVars, "")

	// Add common backend variables
	envVars = append(envVars, "PORT=8080")
	envVars = append(envVars, "GO_ENV=development")
	envVars = append(envVars, "")

	// Add database configuration
	envVars = append(envVars, "# Database Configuration")
	envVars = append(envVars, "DATABASE_URL=")
	envVars = append(envVars, "DB_HOST=localhost")
	envVars = append(envVars, "DB_PORT=5432")
	envVars = append(envVars, "DB_NAME=")
	envVars = append(envVars, "DB_USER=")
	envVars = append(envVars, "DB_PASSWORD=")
	envVars = append(envVars, "")

	// Add any shared environment variables
	for key, value := range config.SharedEnvironment {
		if !strings.HasPrefix(key, "NEXT_PUBLIC_") {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
	}

	content := strings.Join(envVars, "\n")
	if im.verbose {
		fmt.Printf("Generated .env for %s:\n%s\n", comp.Name, content)
	}

	return nil
}

// GenerateScripts creates build and run scripts
func (im *IntegrationManager) GenerateScripts(components []*models.Component) error {
	// Generate build script
	buildScript, err := im.generateBuildScript(components)
	if err != nil {
		return fmt.Errorf("failed to generate build script: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated build.sh:\n%s\n", buildScript)
	}

	// Generate development run script
	devScript, err := im.generateDevScript(components)
	if err != nil {
		return fmt.Errorf("failed to generate dev script: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated dev.sh:\n%s\n", devScript)
	}

	// Generate production run script
	prodScript, err := im.generateProdScript(components)
	if err != nil {
		return fmt.Errorf("failed to generate prod script: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated prod.sh:\n%s\n", prodScript)
	}

	// Generate Docker-based execution script
	dockerScript, err := im.generateDockerScript(components)
	if err != nil {
		return fmt.Errorf("failed to generate docker script: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated docker.sh:\n%s\n", dockerScript)
	}

	return nil
}

// generateBuildScript creates a script to build all components
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateBuildScript(components []*models.Component) (string, error) {
	var lines []string

	lines = append(lines, "#!/bin/bash")
	lines = append(lines, "# Build script for all project components")
	lines = append(lines, "# Generated by Open Source Project Generator")
	lines = append(lines, "")
	lines = append(lines, "set -e  # Exit on error")
	lines = append(lines, "")
	lines = append(lines, "echo \"Building all components...\"")
	lines = append(lines, "")

	for _, comp := range components {
		lines = append(lines, fmt.Sprintf("# Build %s (%s)", comp.Name, comp.Type))

		switch comp.Type {
		case "nextjs":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("echo \"Building %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s", relPath))
			lines = append(lines, "npm install")
			lines = append(lines, "npm run build")
			lines = append(lines, "cd ..")
		case "go-backend":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("echo \"Building %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s", relPath))
			lines = append(lines, "go mod download")
			lines = append(lines, "go build -o bin/server .")
			lines = append(lines, "cd ..")
		case "android":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("echo \"Building %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s", relPath))
			lines = append(lines, "./gradlew build")
			lines = append(lines, "cd ..")
		case "ios":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("echo \"Building %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s", relPath))
			lines = append(lines, "xcodebuild -scheme MyApp -configuration Release")
			lines = append(lines, "cd ..")
		}
		lines = append(lines, "")
	}

	lines = append(lines, "echo \"All components built successfully!\"")

	return strings.Join(lines, "\n"), nil
}

// generateDevScript creates a script to run all components in development mode
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateDevScript(components []*models.Component) (string, error) {
	var lines []string

	lines = append(lines, "#!/bin/bash")
	lines = append(lines, "# Development run script for all project components")
	lines = append(lines, "# Generated by Open Source Project Generator")
	lines = append(lines, "")
	lines = append(lines, "echo \"Starting all components in development mode...\"")
	lines = append(lines, "")

	// Start components in background
	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("# Start %s", comp.Name))
			lines = append(lines, fmt.Sprintf("echo \"Starting %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s && npm run dev &", relPath))
			lines = append(lines, fmt.Sprintf("%s_PID=$!", strings.ToUpper(sanitizeServiceName(comp.Name))))
			lines = append(lines, "cd ..")
		case "go-backend":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("# Start %s", comp.Name))
			lines = append(lines, fmt.Sprintf("echo \"Starting %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s && go run main.go &", relPath))
			lines = append(lines, fmt.Sprintf("%s_PID=$!", strings.ToUpper(sanitizeServiceName(comp.Name))))
			lines = append(lines, "cd ..")
		}
		lines = append(lines, "")
	}

	// Add cleanup handler
	lines = append(lines, "# Cleanup function")
	lines = append(lines, "cleanup() {")
	lines = append(lines, "    echo \"Stopping all components...\"")
	for _, comp := range components {
		if comp.Type == "nextjs" || comp.Type == "go-backend" {
			pidVar := strings.ToUpper(sanitizeServiceName(comp.Name)) + "_PID"
			lines = append(lines, fmt.Sprintf("    kill $%s 2>/dev/null || true", pidVar))
		}
	}
	lines = append(lines, "    exit 0")
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, "# Set up signal handlers")
	lines = append(lines, "trap cleanup SIGINT SIGTERM")
	lines = append(lines, "")
	lines = append(lines, "echo \"All components started. Press Ctrl+C to stop.\"")
	lines = append(lines, "wait")

	return strings.Join(lines, "\n"), nil
}

// generateProdScript creates a script to run all components in production mode
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateProdScript(components []*models.Component) (string, error) {
	var lines []string

	lines = append(lines, "#!/bin/bash")
	lines = append(lines, "# Production run script for all project components")
	lines = append(lines, "# Generated by Open Source Project Generator")
	lines = append(lines, "")
	lines = append(lines, "set -e  # Exit on error")
	lines = append(lines, "")
	lines = append(lines, "echo \"Starting all components in production mode...\"")
	lines = append(lines, "")

	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("# Start %s", comp.Name))
			lines = append(lines, fmt.Sprintf("echo \"Starting %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s && npm start &", relPath))
			lines = append(lines, fmt.Sprintf("%s_PID=$!", strings.ToUpper(sanitizeServiceName(comp.Name))))
			lines = append(lines, "cd ..")
		case "go-backend":
			relPath := filepath.Base(comp.Path)
			lines = append(lines, fmt.Sprintf("# Start %s", comp.Name))
			lines = append(lines, fmt.Sprintf("echo \"Starting %s...\"", comp.Name))
			lines = append(lines, fmt.Sprintf("cd %s && ./bin/server &", relPath))
			lines = append(lines, fmt.Sprintf("%s_PID=$!", strings.ToUpper(sanitizeServiceName(comp.Name))))
			lines = append(lines, "cd ..")
		}
		lines = append(lines, "")
	}

	// Add cleanup handler
	lines = append(lines, "# Cleanup function")
	lines = append(lines, "cleanup() {")
	lines = append(lines, "    echo \"Stopping all components...\"")
	for _, comp := range components {
		if comp.Type == "nextjs" || comp.Type == "go-backend" {
			pidVar := strings.ToUpper(sanitizeServiceName(comp.Name)) + "_PID"
			lines = append(lines, fmt.Sprintf("    kill $%s 2>/dev/null || true", pidVar))
		}
	}
	lines = append(lines, "    exit 0")
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, "trap cleanup SIGINT SIGTERM")
	lines = append(lines, "")
	lines = append(lines, "echo \"All components started. Press Ctrl+C to stop.\"")
	lines = append(lines, "wait")

	return strings.Join(lines, "\n"), nil
}

// generateDockerScript creates a script to run all components using Docker
func (im *IntegrationManager) generateDockerScript(components []*models.Component) (string, error) {
	var lines []string

	lines = append(lines, "#!/bin/bash")
	lines = append(lines, "# Docker-based execution script")
	lines = append(lines, "# Generated by Open Source Project Generator")
	lines = append(lines, "")
	lines = append(lines, "set -e  # Exit on error")
	lines = append(lines, "")
	lines = append(lines, "# Check if docker-compose is installed")
	lines = append(lines, "if ! command -v docker-compose &> /dev/null; then")
	lines = append(lines, "    echo \"Error: docker-compose is not installed\"")
	lines = append(lines, "    exit 1")
	lines = append(lines, "fi")
	lines = append(lines, "")
	lines = append(lines, "# Parse command line arguments")
	lines = append(lines, "COMMAND=${1:-up}")
	lines = append(lines, "")
	lines = append(lines, "case $COMMAND in")
	lines = append(lines, "    up)")
	lines = append(lines, "        echo \"Starting all services with Docker Compose...\"")
	lines = append(lines, "        docker-compose up -d")
	lines = append(lines, "        echo \"Services started. Use 'docker-compose logs -f' to view logs.\"")
	lines = append(lines, "        ;;")
	lines = append(lines, "    down)")
	lines = append(lines, "        echo \"Stopping all services...\"")
	lines = append(lines, "        docker-compose down")
	lines = append(lines, "        ;;")
	lines = append(lines, "    logs)")
	lines = append(lines, "        docker-compose logs -f")
	lines = append(lines, "        ;;")
	lines = append(lines, "    build)")
	lines = append(lines, "        echo \"Building all services...\"")
	lines = append(lines, "        docker-compose build")
	lines = append(lines, "        ;;")
	lines = append(lines, "    restart)")
	lines = append(lines, "        echo \"Restarting all services...\"")
	lines = append(lines, "        docker-compose restart")
	lines = append(lines, "        ;;")
	lines = append(lines, "    *)")
	lines = append(lines, "        echo \"Usage: $0 {up|down|logs|build|restart}\"")
	lines = append(lines, "        echo \"  up      - Start all services\"")
	lines = append(lines, "        echo \"  down    - Stop all services\"")
	lines = append(lines, "        echo \"  logs    - View service logs\"")
	lines = append(lines, "        echo \"  build   - Build all services\"")
	lines = append(lines, "        echo \"  restart - Restart all services\"")
	lines = append(lines, "        exit 1")
	lines = append(lines, "        ;;")
	lines = append(lines, "esac")

	return strings.Join(lines, "\n"), nil
}

// GenerateDocumentation creates root-level documentation
func (im *IntegrationManager) GenerateDocumentation(components []*models.Component, config *models.IntegrationConfig) error {
	// Generate main README
	readme, err := im.generateMainReadme(components, config)
	if err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated README.md:\n%s\n", readme)
	}

	// Generate troubleshooting guide
	troubleshooting, err := im.generateTroubleshootingGuide(components)
	if err != nil {
		return fmt.Errorf("failed to generate troubleshooting guide: %w", err)
	}

	if im.verbose {
		fmt.Printf("Generated TROUBLESHOOTING.md:\n%s\n", troubleshooting)
	}

	return nil
}

// generateMainReadme creates a comprehensive README for the project
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateMainReadme(components []*models.Component, config *models.IntegrationConfig) (string, error) {
	var sections []string

	// Header
	sections = append(sections, "# Project Overview")
	sections = append(sections, "")
	sections = append(sections, "This project was generated using the Open Source Project Generator.")
	sections = append(sections, "")

	// Table of Contents
	sections = append(sections, "## Table of Contents")
	sections = append(sections, "")
	sections = append(sections, "- [Project Structure](#project-structure)")
	sections = append(sections, "- [Components](#components)")
	sections = append(sections, "- [Getting Started](#getting-started)")
	sections = append(sections, "- [Development](#development)")
	sections = append(sections, "- [Production](#production)")
	sections = append(sections, "- [Docker](#docker)")
	sections = append(sections, "- [Troubleshooting](#troubleshooting)")
	sections = append(sections, "")

	// Project Structure
	sections = append(sections, "## Project Structure")
	sections = append(sections, "")
	sections = append(sections, "```")
	sections = append(sections, ".")
	for _, comp := range components {
		relPath := filepath.Base(comp.Path)
		sections = append(sections, fmt.Sprintf("├── %s/          # %s (%s)", relPath, comp.Name, comp.Type))
	}
	sections = append(sections, "├── .env              # Environment configuration")
	if config.GenerateDockerCompose {
		sections = append(sections, "├── docker-compose.yml # Docker orchestration")
	}
	if config.GenerateScripts {
		sections = append(sections, "├── build.sh          # Build all components")
		sections = append(sections, "├── dev.sh            # Run in development mode")
		sections = append(sections, "├── prod.sh           # Run in production mode")
		sections = append(sections, "├── docker.sh         # Docker management")
	}
	sections = append(sections, "└── README.md         # This file")
	sections = append(sections, "```")
	sections = append(sections, "")

	// Components
	sections = append(sections, "## Components")
	sections = append(sections, "")
	for _, comp := range components {
		sections = append(sections, fmt.Sprintf("### %s", comp.Name))
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("**Type:** %s", comp.Type))
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("**Location:** `%s`", filepath.Base(comp.Path)))
		sections = append(sections, "")

		switch comp.Type {
		case "nextjs":
			sections = append(sections, "A Next.js frontend application with TypeScript and Tailwind CSS.")
			sections = append(sections, "")
			sections = append(sections, "**Key Features:**")
			sections = append(sections, "- Server-side rendering")
			sections = append(sections, "- TypeScript for type safety")
			sections = append(sections, "- Tailwind CSS for styling")
			sections = append(sections, "- App Router architecture")
		case "go-backend":
			sections = append(sections, "A Go backend API server using the Gin framework.")
			sections = append(sections, "")
			sections = append(sections, "**Key Features:**")
			sections = append(sections, "- RESTful API endpoints")
			sections = append(sections, "- Gin web framework")
			sections = append(sections, "- Structured logging")
			sections = append(sections, "- Environment-based configuration")
		case "android":
			sections = append(sections, "An Android mobile application built with Kotlin.")
			sections = append(sections, "")
			sections = append(sections, "**Key Features:**")
			sections = append(sections, "- Native Android development")
			sections = append(sections, "- Kotlin programming language")
			sections = append(sections, "- Material Design components")
		case "ios":
			sections = append(sections, "An iOS mobile application built with Swift.")
			sections = append(sections, "")
			sections = append(sections, "**Key Features:**")
			sections = append(sections, "- Native iOS development")
			sections = append(sections, "- Swift programming language")
			sections = append(sections, "- SwiftUI for modern UI")
		}
		sections = append(sections, "")
	}

	// Getting Started
	sections = append(sections, "## Getting Started")
	sections = append(sections, "")
	sections = append(sections, "### Prerequisites")
	sections = append(sections, "")

	prereqs := make(map[string]bool)
	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			prereqs["Node.js (v18 or later)"] = true
			prereqs["npm or yarn"] = true
		case "go-backend":
			prereqs["Go (v1.21 or later)"] = true
		case "android":
			prereqs["Android Studio"] = true
			prereqs["JDK 11 or later"] = true
		case "ios":
			prereqs["Xcode (macOS only)"] = true
		}
	}
	if config.GenerateDockerCompose {
		prereqs["Docker and Docker Compose"] = true
	}

	for prereq := range prereqs {
		sections = append(sections, fmt.Sprintf("- %s", prereq))
	}
	sections = append(sections, "")

	// Installation
	sections = append(sections, "### Installation")
	sections = append(sections, "")
	sections = append(sections, "1. Clone the repository")
	sections = append(sections, "2. Copy `.env.example` to `.env` and configure environment variables")
	sections = append(sections, "3. Install dependencies for each component (see component-specific instructions below)")
	sections = append(sections, "")

	for _, comp := range components {
		relPath := filepath.Base(comp.Path)
		switch comp.Type {
		case "nextjs":
			sections = append(sections, fmt.Sprintf("**%s:**", comp.Name))
			sections = append(sections, "```bash")
			sections = append(sections, fmt.Sprintf("cd %s", relPath))
			sections = append(sections, "npm install")
			sections = append(sections, "```")
			sections = append(sections, "")
		case "go-backend":
			sections = append(sections, fmt.Sprintf("**%s:**", comp.Name))
			sections = append(sections, "```bash")
			sections = append(sections, fmt.Sprintf("cd %s", relPath))
			sections = append(sections, "go mod download")
			sections = append(sections, "```")
			sections = append(sections, "")
		}
	}

	// Development
	sections = append(sections, "## Development")
	sections = append(sections, "")
	if config.GenerateScripts {
		sections = append(sections, "### Quick Start")
		sections = append(sections, "")
		sections = append(sections, "Run all components in development mode:")
		sections = append(sections, "```bash")
		sections = append(sections, "./dev.sh")
		sections = append(sections, "```")
		sections = append(sections, "")
	}

	sections = append(sections, "### Running Components Individually")
	sections = append(sections, "")
	for _, comp := range components {
		relPath := filepath.Base(comp.Path)
		sections = append(sections, fmt.Sprintf("**%s:**", comp.Name))
		sections = append(sections, "```bash")
		sections = append(sections, fmt.Sprintf("cd %s", relPath))
		switch comp.Type {
		case "nextjs":
			sections = append(sections, "npm run dev")
			sections = append(sections, "# Access at http://localhost:3000")
		case "go-backend":
			sections = append(sections, "go run main.go")
			sections = append(sections, "# Access at http://localhost:8080")
		case "android":
			sections = append(sections, "# Open in Android Studio and run")
		case "ios":
			sections = append(sections, "# Open in Xcode and run")
		}
		sections = append(sections, "```")
		sections = append(sections, "")
	}

	// Production
	sections = append(sections, "## Production")
	sections = append(sections, "")
	if config.GenerateScripts {
		sections = append(sections, "### Build All Components")
		sections = append(sections, "")
		sections = append(sections, "```bash")
		sections = append(sections, "./build.sh")
		sections = append(sections, "```")
		sections = append(sections, "")
		sections = append(sections, "### Run in Production Mode")
		sections = append(sections, "")
		sections = append(sections, "```bash")
		sections = append(sections, "./prod.sh")
		sections = append(sections, "```")
		sections = append(sections, "")
	}

	// Docker
	if config.GenerateDockerCompose {
		sections = append(sections, "## Docker")
		sections = append(sections, "")
		sections = append(sections, "### Using Docker Compose")
		sections = append(sections, "")
		sections = append(sections, "Start all services:")
		sections = append(sections, "```bash")
		sections = append(sections, "./docker.sh up")
		sections = append(sections, "```")
		sections = append(sections, "")
		sections = append(sections, "Stop all services:")
		sections = append(sections, "```bash")
		sections = append(sections, "./docker.sh down")
		sections = append(sections, "```")
		sections = append(sections, "")
		sections = append(sections, "View logs:")
		sections = append(sections, "```bash")
		sections = append(sections, "./docker.sh logs")
		sections = append(sections, "```")
		sections = append(sections, "")
	}

	// Troubleshooting
	sections = append(sections, "## Troubleshooting")
	sections = append(sections, "")
	sections = append(sections, "See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues and solutions.")
	sections = append(sections, "")

	// Footer
	sections = append(sections, "## Additional Resources")
	sections = append(sections, "")
	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			sections = append(sections, "- [Next.js Documentation](https://nextjs.org/docs)")
		case "go-backend":
			sections = append(sections, "- [Gin Framework Documentation](https://gin-gonic.com/docs/)")
		case "android":
			sections = append(sections, "- [Android Developer Documentation](https://developer.android.com/docs)")
		case "ios":
			sections = append(sections, "- [iOS Developer Documentation](https://developer.apple.com/documentation/)")
		}
	}
	sections = append(sections, "")

	return strings.Join(sections, "\n"), nil
}

// generateTroubleshootingGuide creates a troubleshooting guide
//
//nolint:unparam // error return reserved for future validation
func (im *IntegrationManager) generateTroubleshootingGuide(components []*models.Component) (string, error) {
	var sections []string

	sections = append(sections, "# Troubleshooting Guide")
	sections = append(sections, "")
	sections = append(sections, "This guide covers common issues and their solutions.")
	sections = append(sections, "")

	// Common Issues
	sections = append(sections, "## Common Issues")
	sections = append(sections, "")

	// Port conflicts
	sections = append(sections, "### Port Already in Use")
	sections = append(sections, "")
	sections = append(sections, "**Problem:** Error message indicating a port is already in use.")
	sections = append(sections, "")
	sections = append(sections, "**Solution:**")
	sections = append(sections, "1. Check which process is using the port:")
	sections = append(sections, "   ```bash")
	sections = append(sections, "   # On macOS/Linux")
	sections = append(sections, "   lsof -i :3000")
	sections = append(sections, "   lsof -i :8080")
	sections = append(sections, "   ```")
	sections = append(sections, "2. Kill the process or change the port in your configuration")
	sections = append(sections, "")

	// Environment variables
	sections = append(sections, "### Environment Variables Not Loading")
	sections = append(sections, "")
	sections = append(sections, "**Problem:** Application can't find environment variables.")
	sections = append(sections, "")
	sections = append(sections, "**Solution:**")
	sections = append(sections, "1. Ensure `.env` file exists in the project root")
	sections = append(sections, "2. Check that environment variables are properly formatted (no spaces around `=`)")
	sections = append(sections, "3. Restart the application after changing `.env`")
	sections = append(sections, "")

	// Component-specific issues
	for _, comp := range components {
		switch comp.Type {
		case "nextjs":
			sections = append(sections, "### Next.js Build Errors")
			sections = append(sections, "")
			sections = append(sections, "**Problem:** Build fails with TypeScript or dependency errors.")
			sections = append(sections, "")
			sections = append(sections, "**Solution:**")
			sections = append(sections, "1. Delete `node_modules` and `.next` directories")
			sections = append(sections, "2. Run `npm install` again")
			sections = append(sections, "3. Clear npm cache: `npm cache clean --force`")
			sections = append(sections, "")

		case "go-backend":
			sections = append(sections, "### Go Module Issues")
			sections = append(sections, "")
			sections = append(sections, "**Problem:** Go can't find modules or dependencies.")
			sections = append(sections, "")
			sections = append(sections, "**Solution:**")
			sections = append(sections, "1. Run `go mod tidy` to clean up dependencies")
			sections = append(sections, "2. Run `go mod download` to re-download modules")
			sections = append(sections, "3. Check your `GOPATH` and `GOMODCACHE` settings")
			sections = append(sections, "")

		case "android":
			sections = append(sections, "### Android Build Failures")
			sections = append(sections, "")
			sections = append(sections, "**Problem:** Gradle build fails.")
			sections = append(sections, "")
			sections = append(sections, "**Solution:**")
			sections = append(sections, "1. Clean the project: `./gradlew clean`")
			sections = append(sections, "2. Invalidate caches in Android Studio")
			sections = append(sections, "3. Check that ANDROID_HOME is set correctly")
			sections = append(sections, "")

		case "ios":
			sections = append(sections, "### iOS Build Failures")
			sections = append(sections, "")
			sections = append(sections, "**Problem:** Xcode build fails.")
			sections = append(sections, "")
			sections = append(sections, "**Solution:**")
			sections = append(sections, "1. Clean build folder: Product > Clean Build Folder in Xcode")
			sections = append(sections, "2. Delete derived data")
			sections = append(sections, "3. Ensure code signing is configured correctly")
			sections = append(sections, "")
		}
	}

	// Docker issues
	sections = append(sections, "### Docker Issues")
	sections = append(sections, "")
	sections = append(sections, "**Problem:** Docker containers won't start or crash.")
	sections = append(sections, "")
	sections = append(sections, "**Solution:**")
	sections = append(sections, "1. Check Docker logs: `docker-compose logs`")
	sections = append(sections, "2. Rebuild containers: `docker-compose build --no-cache`")
	sections = append(sections, "3. Remove old containers: `docker-compose down -v`")
	sections = append(sections, "")

	// Getting Help
	sections = append(sections, "## Getting Help")
	sections = append(sections, "")
	sections = append(sections, "If you continue to experience issues:")
	sections = append(sections, "")
	sections = append(sections, "1. Check the component-specific documentation")
	sections = append(sections, "2. Review application logs for error messages")
	sections = append(sections, "3. Ensure all prerequisites are installed and up to date")
	sections = append(sections, "")

	return strings.Join(sections, "\n"), nil
}

// sanitizeServiceName converts a component name to a valid Docker service name
func sanitizeServiceName(name string) string {
	// Replace invalid characters with hyphens
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	return name
}
