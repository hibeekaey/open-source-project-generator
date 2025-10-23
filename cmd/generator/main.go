package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli"
	interactivemode "github.com/cuesoftinc/open-source-project-generator/pkg/cli/interactive"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

var (
	// Version information (set during build)
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"

	// Global flags
	configFile      string
	outputDir       string
	verbose         bool
	dryRun          bool
	noExternalTools bool
	createBackup    bool
	forceOverwrite  bool
	interactive     bool
	offlineMode     bool
	streamOutput    bool
	noRollback      bool
)

func main() {
	// Initialize logger for exit code handler
	log := logger.NewLogger()
	exitHandler := cli.NewExitCodeHandler(log)
	suggestionEngine := cli.NewSuggestionEngine(verbose)
	diagnostics := cli.NewDiagnosticsCollector(log, verbose)

	if err := rootCmd.Execute(); err != nil {
		// Determine appropriate exit code based on error
		exitCode := exitHandler.DetermineExitCode(err)

		// Display error message
		fmt.Fprintf(os.Stderr, "\n%s\n", log.GetFormatter().Error("Error: "+err.Error()))

		// In verbose mode, display full diagnostic information
		if verbose {
			context := map[string]interface{}{
				"command":   os.Args,
				"exit_code": exitCode,
			}
			diagnostics.LogDiagnostics(err, context)

			// Display formatted diagnostics
			diagnosticOutput := diagnostics.FormatVerboseError(err, context)
			if diagnosticOutput != "" {
				fmt.Fprintln(os.Stderr, diagnosticOutput)
			}
		}

		// Generate and display suggestions
		suggestions := suggestionEngine.GenerateSuggestions(err)
		if len(suggestions) > 0 {
			fmt.Fprintf(os.Stderr, "\n%s\n", log.GetFormatter().Section("Next Steps:"))
			for i, suggestion := range suggestions {
				fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, suggestion)
			}
			fmt.Fprintln(os.Stderr)
		}

		// Exit with appropriate code
		os.Exit(int(exitCode))
	}
}

var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "Open Source Project Generator - Modern tool-orchestration architecture",
	Long: `A CLI tool that generates production-ready project scaffolding across multiple
technology stacks by orchestrating industry-standard bootstrap tools.

Supports: Next.js, Go/Gin, Android/Kotlin, iOS/Swift, Docker, Kubernetes, Terraform`,
	Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitCommit),
}

func init() {
	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(checkToolsCmd)
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.AddCommand(cacheToolsCmd)
	rootCmd.AddCommand(versionCmd)
}

var checkToolsCmd = &cobra.Command{
	Use:   "check-tools",
	Short: "Check availability of required bootstrap tools",
	Long: `Check which bootstrap tools are available on your system and display
installation instructions for any missing tools.

This command validates your environment before project generation and helps
you install any missing dependencies.

Examples:
  # Check all registered tools
  generator check-tools

  # Check with verbose output
  generator check-tools --verbose

  # Check specific tools
  generator check-tools npx go gradle`,
	RunE: runCheckTools,
}

func init() {
	checkToolsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runCheckTools(cmd *cobra.Command, args []string) error {
	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)
	}

	// Create suggestion engine
	suggestionEngine := cli.NewSuggestionEngine(verbose)

	// Create tool discovery
	toolDiscovery := orchestrator.NewToolDiscovery(log)

	// Determine which tools to check
	var toolsToCheck []string
	if len(args) > 0 {
		// Check specific tools provided as arguments
		toolsToCheck = args
	} else {
		// Check all registered tools
		toolsToCheck = toolDiscovery.ListRegisteredTools()
	}

	log.Info(fmt.Sprintf("Checking %d tools...", len(toolsToCheck)))

	// Check tool requirements
	result, err := toolDiscovery.CheckRequirements(toolsToCheck)
	if err != nil {
		return cli.NewToolError("unknown", "tool check failed", err).
			WithSuggestions(suggestionEngine.GenerateSuggestions(err)...)
	}

	toolCheckResult := result.(*models.ToolCheckResult)

	// Display results
	displayToolCheckResults(toolCheckResult, toolDiscovery, log, suggestionEngine)

	// Return error if not all tools are available (non-zero exit code)
	if !toolCheckResult.AllAvailable {
		missingTools := strings.Join(toolCheckResult.Missing, ", ")
		return cli.NewToolError(missingTools, "some required tools are not available", nil).
			WithSuggestions(
				"Install the missing tools listed above",
				"Use --no-external-tools flag to force fallback generation",
				"Run 'generator check-tools' again after installation to verify",
			)
	}

	return nil
}

func displayToolCheckResults(result *models.ToolCheckResult, discovery *orchestrator.ToolDiscovery, log *logger.Logger, suggestionEngine *cli.SuggestionEngine) {
	fmt.Println("\n" + separator("="))
	fmt.Println("TOOL AVAILABILITY CHECK")
	fmt.Println(separator("="))

	availableCount := len(result.Tools) - len(result.Missing)
	fmt.Printf("\nStatus: %d of %d tools available\n", availableCount, len(result.Tools))
	fmt.Printf("Checked at: %s\n", result.CheckedAt.Format(time.RFC3339))

	// Display available tools
	if availableCount > 0 {
		fmt.Println("\n✓ Available Tools:")
		fmt.Println(separator("-"))
		for name, tool := range result.Tools {
			if tool.Available {
				versionInfo := ""
				if tool.InstalledVersion != "" {
					versionInfo = fmt.Sprintf(" (version: %s)", tool.InstalledVersion)
				}
				fmt.Printf("  ✓ %s%s\n", name, versionInfo)

				// Show component types this tool supports
				if metadata, err := discovery.GetToolMetadata(name); err == nil && len(metadata.ComponentTypes) > 0 {
					fmt.Printf("    Supports: %v\n", metadata.ComponentTypes)
				}
			}
		}
	}

	// Display missing tools
	if len(result.Missing) > 0 {
		fmt.Println("\n✗ Missing Tools:")
		fmt.Println(separator("-"))
		for _, name := range result.Missing {
			fmt.Printf("  ✗ %s\n", name)

			// Get metadata to show what it's used for
			if metadata, err := discovery.GetToolMetadata(name); err == nil {
				if len(metadata.ComponentTypes) > 0 {
					fmt.Printf("    Required for: %v\n", metadata.ComponentTypes)
				}
				if metadata.FallbackAvailable {
					fmt.Println("    Note: Fallback generation available")
				}
			}

			// Show installation instructions using suggestion engine
			installSuggestion := suggestionEngine.GetToolInstallSuggestion(name)
			fmt.Printf("\n    Installation:\n")
			fmt.Printf("    %s\n", installSuggestion)
		}
	}

	// Display outdated tools
	if len(result.Outdated) > 0 {
		fmt.Println("\n⚠ Outdated Tools:")
		fmt.Println(separator("-"))
		for _, name := range result.Outdated {
			if tool, exists := result.Tools[name]; exists {
				fmt.Printf("  ⚠ %s (installed: %s, minimum: %s)\n",
					name, tool.InstalledVersion, tool.MinVersion)
				fmt.Printf("    Update recommended for best compatibility\n")
			}
		}
	}

	fmt.Println("\n" + separator("="))
	if result.AllAvailable {
		fmt.Println("✓ All required tools are available!")
		fmt.Println("You can proceed with project generation.")
	} else {
		fmt.Println("✗ Some tools are missing.")
		fmt.Println("\nNext Steps:")
		fmt.Println("  1. Install the missing tools using the instructions above")
		fmt.Println("  2. Run 'generator check-tools' again to verify installation")
		fmt.Println("  3. Or use --no-external-tools flag for fallback generation")
	}
	fmt.Println(separator("=") + "\n")
}

// indentText indents each line of text with the given prefix
func indentText(text, indent string) string {
	lines := splitLines(text)
	result := ""
	for _, line := range lines {
		result += indent + line + "\n"
	}
	return result
}

// splitLines splits text into lines
func splitLines(text string) []string {
	var lines []string
	current := ""
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

var initConfigCmd = &cobra.Command{
	Use:   "init-config [output-file]",
	Short: "Generate a configuration template file",
	Long: `Generate a configuration template file with all available options documented.

This command creates a YAML configuration file that you can customize for your
project. The template includes examples for all supported component types and
integration options with inline comments explaining each field.

Examples:
  # Generate config in current directory
  generator init-config

  # Generate config with custom name
  generator init-config my-project.yaml

  # Generate minimal config (only required fields)
  generator init-config --minimal

  # Generate full-stack example
  generator init-config --example fullstack

  # Generate frontend-only example
  generator init-config --example frontend

  # Generate backend-only example
  generator init-config --example backend

  # Generate mobile-only example
  generator init-config --example mobile

  # Generate microservice example
  generator init-config --example microservice

  # Force overwrite existing file
  generator init-config --force my-project.yaml`,
	RunE: runInitConfig,
}

var (
	minimalConfig      bool
	exampleType        string
	forceInitOverwrite bool
)

func init() {
	initConfigCmd.Flags().BoolVar(&minimalConfig, "minimal", false, "Generate minimal configuration with only required fields")
	initConfigCmd.Flags().StringVar(&exampleType, "example", "", "Generate example configuration (fullstack, frontend, backend, mobile, microservice)")
	initConfigCmd.Flags().BoolVar(&forceInitOverwrite, "force", false, "Force overwrite existing configuration file")
}

func runInitConfig(cmd *cobra.Command, args []string) error {
	// Determine output file
	outputFile := "project-config.yaml"
	if len(args) > 0 {
		outputFile = args[0]
	}

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		if !forceInitOverwrite {
			return cli.NewFileSystemError("create", outputFile, fmt.Errorf("file already exists")).
				WithSuggestions(
					"Remove the existing file first: rm "+outputFile,
					"Choose a different output filename",
					"Use --force flag to overwrite: generator init-config --force "+outputFile,
				)
		}
		// File exists but --force flag is set, so we'll overwrite
		fmt.Printf("⚠ Overwriting existing file: %s\n", outputFile)
	}

	// Generate configuration based on flags
	var config *models.ProjectConfig
	if exampleType != "" {
		config = generateExampleConfig(exampleType)
		if config == nil {
			return cli.NewConfigError("invalid example type: "+exampleType, nil).
				WithSuggestions(
					"Valid example types: fullstack, frontend, backend, mobile, microservice",
					"Use 'generator init-config --help' to see all options",
				)
		}
	} else if minimalConfig {
		config = generateMinimalConfig()
	} else {
		config = generateFullConfig()
	}

	// Marshal to YAML with inline documentation
	data, err := marshalConfigWithComments(config)
	if err != nil {
		return cli.NewConfigError("failed to marshal configuration", err).
			WithSuggestions(
				"This is an internal error - please report it",
				"Try using a different example type",
			)
	}

	// Add header comments
	header := generateConfigHeader()
	fullContent := header + "\n" + data

	// Write to file
	if err := os.WriteFile(outputFile, []byte(fullContent), 0644); err != nil {
		return cli.NewFileSystemError("write", outputFile, err).
			WithSuggestions(
				"Check that you have write permissions in the current directory",
				"Verify that disk space is available",
				"Try specifying a different output path",
			)
	}

	fmt.Printf("✓ Configuration template created: %s\n", outputFile)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the configuration file to match your project requirements")
	fmt.Println("2. Enable/disable components as needed")
	fmt.Println("3. Run: generator generate --config " + outputFile)

	return nil
}

func generateConfigHeader() string {
	return `# Open Source Project Generator Configuration
# 
# This file defines the project structure and components to generate.
# Edit this file to customize your project.
#
# Configuration Structure:
# - name: Project name (used for directory names and documentation)
# - description: Brief description of the project
# - output_dir: Where to generate the project (relative or absolute path)
# - components: List of components to generate (frontend, backend, mobile, etc.)
# - integration: Settings for integrating components together
# - options: Generation options (dry-run, verbose, backups, etc.)
#
# Component Types:
# - nextjs: Next.js frontend application with TypeScript and Tailwind CSS
# - go-backend: Go backend API server with Gin framework
# - android: Android mobile app with Kotlin
# - ios: iOS mobile app with Swift
#
# Common Use Cases:
# - Fullstack: Enable nextjs + go-backend components
# - Frontend Only: Enable only nextjs component
# - Backend Only: Enable only go-backend component
# - Mobile: Enable android + ios components
# - Microservice: Single go-backend with Docker/K8s support
#
# Examples:
# - Fullstack: generator init-config --example fullstack
# - Frontend: generator init-config --example frontend
# - Backend: generator init-config --example backend
# - Mobile: generator init-config --example mobile
# - Microservice: generator init-config --example microservice
#
# Documentation: https://github.com/cuesoftinc/open-source-project-generator
# Issues: https://github.com/cuesoftinc/open-source-project-generator/issues

`
}

func generateMinimalConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-project",
		Description: "A new project",
		OutputDir:   "./my-project",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints:          map[string]string{},
			SharedEnvironment:     map[string]string{},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

func generateFullConfig() *models.ProjectConfig {
	return generateFullstackConfig()
}

// generateFullstackConfig generates a configuration for fullstack projects
// with frontend, backend, and optional mobile components
func generateFullstackConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-fullstack-project",
		Description: "A full-stack project with frontend, backend, and mobile apps",
		OutputDir:   "./my-fullstack-project",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "web-app",
					"typescript": true,
					"tailwind":   true,
					"app_router": true,
					"eslint":     true,
				},
			},
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/user/my-fullstack-project",
					"framework": "gin",
					"port":      8080,
				},
			},
			{
				Type:    "android",
				Name:    "mobile-android",
				Enabled: false,
				Config: map[string]interface{}{
					"name":       "mobile-android",
					"package":    "com.example.myapp",
					"min_sdk":    24,
					"target_sdk": 34,
					"language":   "kotlin",
				},
			},
			{
				Type:    "ios",
				Name:    "mobile-ios",
				Enabled: false,
				Config: map[string]interface{}{
					"name":              "mobile-ios",
					"bundle_id":         "com.example.myapp",
					"deployment_target": "15.0",
					"language":          "swift",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "http://localhost:8080",
			},
			SharedEnvironment: map[string]string{
				"NODE_ENV":    "development",
				"API_URL":     "http://localhost:8080",
				"API_TIMEOUT": "30000",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

// generateFrontendConfig generates a configuration for frontend-only projects
func generateFrontendConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-frontend-project",
		Description: "A frontend-only project with Next.js",
		OutputDir:   "./my-frontend-project",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "web-app",
					"typescript": true,
					"tailwind":   true,
					"app_router": true,
					"eslint":     true,
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "https://api.example.com",
			},
			SharedEnvironment: map[string]string{
				"NEXT_PUBLIC_API_URL": "https://api.example.com",
				"NEXT_PUBLIC_ENV":     "production",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

// generateBackendConfig generates a configuration for backend-only projects
func generateBackendConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-backend-project",
		Description: "A backend-only project with Go API server",
		OutputDir:   "./my-backend-project",
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/user/my-backend-project",
					"framework": "gin",
					"port":      8080,
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "http://localhost:8080",
			},
			SharedEnvironment: map[string]string{
				"PORT":      "8080",
				"HOST":      "0.0.0.0",
				"ENV":       "development",
				"LOG_LEVEL": "info",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

// generateMobileConfig generates a configuration for mobile-only projects
func generateMobileConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-mobile-project",
		Description: "A mobile-only project with Android and iOS apps",
		OutputDir:   "./my-mobile-project",
		Components: []models.ComponentConfig{
			{
				Type:    "android",
				Name:    "mobile-android",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "mobile-android",
					"package":    "com.example.myapp",
					"min_sdk":    24,
					"target_sdk": 34,
					"language":   "kotlin",
				},
			},
			{
				Type:    "ios",
				Name:    "mobile-ios",
				Enabled: true,
				Config: map[string]interface{}{
					"name":              "mobile-ios",
					"bundle_id":         "com.example.myapp",
					"deployment_target": "15.0",
					"language":          "swift",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: false,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "https://api.example.com",
			},
			SharedEnvironment: map[string]string{
				"API_URL":     "https://api.example.com",
				"API_TIMEOUT": "30000",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

// generateMicroserviceConfig generates a configuration optimized for microservices
func generateMicroserviceConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-microservice",
		Description: "A microservice with Go backend, Docker, and Kubernetes support",
		OutputDir:   "./my-microservice",
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "service",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "service",
					"module":    "github.com/user/my-microservice",
					"framework": "gin",
					"port":      8081,
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"service":  "http://localhost:8081",
				"database": "postgres://localhost:5432/servicedb",
				"cache":    "redis://localhost:6379",
			},
			SharedEnvironment: map[string]string{
				"SERVICE_NAME": "my-microservice",
				"SERVICE_PORT": "8081",
				"DB_HOST":      "postgres",
				"DB_PORT":      "5432",
				"REDIS_HOST":   "redis",
				"REDIS_PORT":   "6379",
				"LOG_LEVEL":    "info",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

func generateExampleConfig(exampleType string) *models.ProjectConfig {
	switch exampleType {
	case "fullstack":
		return generateFullstackConfig()
	case "frontend":
		return generateFrontendConfig()
	case "backend":
		return generateBackendConfig()
	case "mobile":
		return generateMobileConfig()
	case "microservice":
		return generateMicroserviceConfig()
	default:
		// Return nil for invalid types - caller will handle error
		return nil
	}
}

// marshalConfigWithComments marshals a ProjectConfig to YAML with inline comments
func marshalConfigWithComments(config *models.ProjectConfig) (string, error) {
	// Marshal to basic YAML first
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	// Add inline comments to explain each section
	lines := strings.Split(string(data), "\n")
	var result []string

	for i, line := range lines {
		// Add comments for major sections
		if strings.HasPrefix(line, "name:") {
			result = append(result, "# Project name - used for directory names and documentation")
		} else if strings.HasPrefix(line, "description:") {
			result = append(result, "# Brief description of what this project does")
		} else if strings.HasPrefix(line, "output_dir:") {
			result = append(result, "# Output directory where the project will be generated")
		} else if strings.HasPrefix(line, "components:") {
			result = append(result, "")
			result = append(result, "# Components to generate - enable/disable as needed")
		} else if strings.HasPrefix(line, "- type: nextjs") {
			result = append(result, "")
			result = append(result, "  # Next.js Frontend Component")
			result = append(result, "  # Generates a modern React application with TypeScript and Tailwind CSS")
		} else if strings.HasPrefix(line, "- type: go-backend") {
			result = append(result, "")
			result = append(result, "  # Go Backend Component")
			result = append(result, "  # Generates a REST API server using the Gin web framework")
		} else if strings.HasPrefix(line, "- type: android") {
			result = append(result, "")
			result = append(result, "  # Android Mobile Component")
			result = append(result, "  # Generates a native Android app with Kotlin")
		} else if strings.HasPrefix(line, "- type: ios") {
			result = append(result, "")
			result = append(result, "  # iOS Mobile Component")
			result = append(result, "  # Generates a native iOS app with Swift")
		} else if strings.HasPrefix(line, "  enabled:") {
			result = append(result, line+" # Set to false to skip this component")
			continue
		} else if strings.HasPrefix(line, "  config:") {
			result = append(result, line+" # Component-specific configuration")
			continue
		} else if strings.HasPrefix(line, "    typescript:") {
			result = append(result, line+" # Use TypeScript instead of JavaScript")
			continue
		} else if strings.HasPrefix(line, "    tailwind:") {
			result = append(result, line+" # Include Tailwind CSS for styling")
			continue
		} else if strings.HasPrefix(line, "    app_router:") {
			result = append(result, line+" # Use Next.js App Router (recommended)")
			continue
		} else if strings.HasPrefix(line, "    eslint:") {
			result = append(result, line+" # Include ESLint for code quality")
			continue
		} else if strings.HasPrefix(line, "    module:") {
			result = append(result, line+" # Go module path (e.g., github.com/user/project)")
			continue
		} else if strings.HasPrefix(line, "    framework:") {
			result = append(result, line+" # Web framework (gin, echo, fiber)")
			continue
		} else if strings.HasPrefix(line, "    port:") {
			result = append(result, line+" # Server port number")
			continue
		} else if strings.HasPrefix(line, "    package:") {
			result = append(result, line+" # Android package name (e.g., com.example.app)")
			continue
		} else if strings.HasPrefix(line, "    min_sdk:") {
			result = append(result, line+" # Minimum Android SDK version")
			continue
		} else if strings.HasPrefix(line, "    target_sdk:") {
			result = append(result, line+" # Target Android SDK version")
			continue
		} else if strings.HasPrefix(line, "    language:") && i > 0 && strings.Contains(lines[i-1], "android") {
			result = append(result, line+" # Programming language (kotlin, java)")
			continue
		} else if strings.HasPrefix(line, "    bundle_id:") {
			result = append(result, line+" # iOS bundle identifier (e.g., com.example.app)")
			continue
		} else if strings.HasPrefix(line, "    deployment_target:") {
			result = append(result, line+" # Minimum iOS version")
			continue
		} else if strings.HasPrefix(line, "    language:") && i > 0 && strings.Contains(lines[i-1], "ios") {
			result = append(result, line+" # Programming language (swift, objective-c)")
			continue
		} else if strings.HasPrefix(line, "integration:") {
			result = append(result, "")
			result = append(result, "# Integration settings - how components work together")
		} else if strings.HasPrefix(line, "  generate_docker_compose:") {
			result = append(result, line+" # Generate docker-compose.yml for local development")
			continue
		} else if strings.HasPrefix(line, "  generate_scripts:") {
			result = append(result, line+" # Generate build and run scripts")
			continue
		} else if strings.HasPrefix(line, "  api_endpoints:") {
			result = append(result, line+" # API endpoint URLs for component communication")
			continue
		} else if strings.HasPrefix(line, "  shared_environment:") {
			result = append(result, line+" # Environment variables shared across components")
			continue
		} else if strings.HasPrefix(line, "options:") {
			result = append(result, "")
			result = append(result, "# Generation options")
		} else if strings.HasPrefix(line, "  use_external_tools:") {
			result = append(result, line+" # Use bootstrap tools (create-next-app, go mod init)")
			continue
		} else if strings.HasPrefix(line, "  dry_run:") {
			result = append(result, line+" # Preview without creating files")
			continue
		} else if strings.HasPrefix(line, "  verbose:") {
			result = append(result, line+" # Enable detailed logging")
			continue
		} else if strings.HasPrefix(line, "  create_backup:") {
			result = append(result, line+" # Backup existing files before overwriting")
			continue
		} else if strings.HasPrefix(line, "  force_overwrite:") {
			result = append(result, line+" # Overwrite existing files without prompting")
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n"), nil
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new project from configuration",
	Long: `Generate a new project using the specified configuration file.

The generate command orchestrates bootstrap tools (like create-next-app, go mod init)
to create project components, then maps them to a standardized directory structure
and integrates them together.

Automatic rollback is enabled by default - if generation fails, all changes will be
reverted automatically. Use --no-rollback to disable this behavior and leave files
in place for inspection.

Examples:
  # Interactive mode (prompts for configuration)
  generator generate --interactive

  # Generate from config file
  generator generate --config project.yaml

  # Generate with custom output directory
  generator generate --config project.yaml --output ./my-project

  # Dry run (preview without creating files)
  generator generate --config project.yaml --dry-run

  # Force fallback generation (don't use external tools)
  generator generate --config project.yaml --no-external-tools

  # Disable automatic rollback on failure
  generator generate --config project.yaml --no-rollback

  # Verbose output for debugging
  generator generate --config project.yaml --verbose`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to project configuration file (YAML/JSON)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated project (overrides config)")
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview what would be generated without creating files")
	generateCmd.Flags().BoolVar(&noExternalTools, "no-external-tools", false, "Force fallback generation (don't use external bootstrap tools)")
	generateCmd.Flags().BoolVar(&offlineMode, "offline", false, "Force offline mode (use cached tools and fallback generators)")
	generateCmd.Flags().BoolVar(&createBackup, "backup", true, "Create backup before overwriting existing directory")
	generateCmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite existing directory")
	generateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode (prompts for configuration)")
	generateCmd.Flags().BoolVar(&streamOutput, "stream-output", false, "Stream real-time output from bootstrap tools")
	generateCmd.Flags().BoolVar(&noRollback, "no-rollback", false, "Disable automatic rollback on failure (leave generated files for inspection)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal, cancelling generation...")
		cancel()
	}()

	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)

		// Log system information in verbose mode
		diagnostics := cli.NewDiagnosticsCollector(log, verbose)
		diagnostics.LogSystemInfo()

		log.Debug("Verbose mode enabled - detailed logging active")
		log.Debug(fmt.Sprintf("Command arguments: %v", os.Args))
	}

	// Load or create configuration
	var config *models.ProjectConfig
	var err error

	if interactive {
		// Interactive mode
		log.Debug("Starting interactive mode...")
		config, err = runInteractiveMode(log)
		if err != nil {
			// Check if user cancelled
			if err == cli.ErrUserCancelled {
				return cli.NewUserCancelledError("interactive mode cancelled by user")
			}

			if verbose {
				log.Debug(fmt.Sprintf("Interactive mode error: %v", err))
			}

			return cli.NewConfigError("interactive mode failed", err).
				WithSuggestions(
					"Try using a configuration file instead: generator generate --config project.yaml",
					"Use 'generator init-config' to create a template configuration",
					"Run with --verbose flag for more details",
				)
		}
		log.Debug("Interactive mode completed successfully")
	} else {
		// Non-interactive mode - require config file
		if configFile == "" {
			return cli.NewConfigError("config file is required in non-interactive mode", nil).
				WithSuggestions(
					"Provide a config file: generator generate --config project.yaml",
					"Use interactive mode: generator generate --interactive",
					"Create a config template: generator init-config",
				)
		}

		log.Debug(fmt.Sprintf("Loading configuration from: %s", configFile))
		config, err = loadConfigFile(configFile)
		if err != nil {
			if verbose {
				log.Debug(fmt.Sprintf("Config load error: %v", err))
			}

			return cli.NewConfigError("failed to load config file: "+configFile, err).
				WithSuggestions(
					"Verify the config file exists and is readable",
					"Check the YAML syntax is valid",
					"Use 'generator init-config' to create a valid template",
					"Run with --verbose flag for detailed error information",
				)
		}
		log.Debug("Configuration loaded successfully")
	}

	// Log configuration details in verbose mode
	if verbose {
		log.Debug(fmt.Sprintf("Project Name: %s", config.Name))
		log.Debug(fmt.Sprintf("Output Directory: %s", config.OutputDir))
		log.Debug(fmt.Sprintf("Number of Components: %d", len(config.Components)))
		for i, comp := range config.Components {
			log.Debug(fmt.Sprintf("  Component %d: %s (%s) - Enabled: %t", i+1, comp.Name, comp.Type, comp.Enabled))
		}
	}

	// Apply command-line overrides
	if outputDir != "" {
		config.OutputDir = outputDir
	}
	if dryRun {
		config.Options.DryRun = true
	}
	if noExternalTools {
		config.Options.UseExternalTools = false
	} else {
		config.Options.UseExternalTools = true
	}
	if createBackup {
		config.Options.CreateBackup = true
	}
	if forceOverwrite {
		config.Options.ForceOverwrite = true
	}
	if verbose {
		config.Options.Verbose = true
	}
	if streamOutput {
		config.Options.StreamOutput = true
	}

	// Create project coordinator
	coordinator := orchestrator.NewProjectCoordinator(log)

	// Set offline mode if requested
	if offlineMode {
		log.Info("Forcing offline mode - will use cached tools and fallback generators")
		coordinator.SetOfflineMode(true)
	}

	// Disable automatic rollback if requested
	if noRollback {
		log.Info("Automatic rollback disabled - generated files will be left in place on failure")
		coordinator.SetAutoRollbackEnabled(false)
	}

	// Check and display offline status
	if coordinator.IsOffline() {
		fmt.Println(coordinator.GetOfflineMessage())
	}

	// Execute generation
	startTime := time.Now()

	if verbose {
		log.Debug("Starting generation execution...")
		log.Debug(fmt.Sprintf("Dry Run: %t", dryRun))
		log.Debug(fmt.Sprintf("Use External Tools: %t", config.Options.UseExternalTools))
		log.Debug(fmt.Sprintf("Offline Mode: %t", offlineMode))
		log.Debug(fmt.Sprintf("Stream Output: %t", streamOutput))
	}

	if dryRun {
		log.Info("Running in dry-run mode (no files will be created)")
		log.Debug("Executing dry-run...")

		result, err := coordinator.DryRun(ctx, config)
		if err != nil {
			if verbose {
				log.Debug(fmt.Sprintf("Dry-run failed: %v", err))
			}

			// Check for specific error types and provide better messages
			if genErr, ok := err.(*orchestrator.GenerationError); ok {
				return cli.NewGenerationError("dry-run", genErr.Message, genErr).
					WithSuggestions(genErr.Suggestions...)
			}
			return cli.NewGenerationError("dry-run", "dry-run failed", err).
				WithSuggestions(
					"Check the configuration for errors",
					"Verify all component configurations are valid",
					"Run with --verbose flag for detailed diagnostic information",
				)
		}

		log.Debug("Dry-run completed successfully")

		// Display preview
		displayPreview(result.(*models.PreviewResult), log)
	} else {
		log.Info(fmt.Sprintf("Starting project generation: %s", config.Name))
		log.Debug(fmt.Sprintf("Generation started at: %s", startTime.Format(time.RFC3339)))

		result, err := coordinator.Generate(ctx, config)
		if err != nil {
			if verbose {
				log.Debug(fmt.Sprintf("Generation failed: %v", err))
				log.Debug(fmt.Sprintf("Error type: %T", err))
			}

			// Check for specific error types and provide better messages
			if genErr, ok := err.(*orchestrator.GenerationError); ok {
				if verbose {
					log.Debug(fmt.Sprintf("Generation error category: %s", genErr.Category))
					log.Debug(fmt.Sprintf("Component: %s", genErr.Component))
					log.Debug(fmt.Sprintf("Recoverable: %t", genErr.Recoverable))
				}

				cliErr := cli.NewGenerationError(genErr.Component, genErr.Message, genErr)

				// Add context-specific suggestions
				if len(genErr.Suggestions) > 0 {
					cliErr = cliErr.WithSuggestions(genErr.Suggestions...)
				} else {
					// Add default suggestions based on category
					switch genErr.Category {
					case orchestrator.ErrCategoryToolNotFound:
						cliErr = cliErr.WithSuggestions(
							"Install the required tool or use --no-external-tools for fallback generation",
							"Run 'generator check-tools' to see installation instructions",
						)
					case orchestrator.ErrCategoryToolExecution:
						cliErr = cliErr.WithSuggestions(
							"Check the tool output above for specific errors",
							"Try running with --verbose flag for detailed logs",
							"Use --no-external-tools to try fallback generation",
						)
					case orchestrator.ErrCategoryFileSystem:
						cliErr = cliErr.WithSuggestions(
							"Check file system permissions and available disk space",
							"Verify the output directory is writable",
							"Try a different output directory",
						)
					case orchestrator.ErrCategoryInvalidConfig:
						cliErr = cliErr.WithSuggestions(
							"Review the configuration file for errors",
							"Use 'generator init-config --example <type>' for a valid template",
							"Check the documentation for configuration schema",
						)
					default:
						cliErr = cliErr.WithSuggestions(
							"Run with --verbose flag for detailed diagnostic information",
							"Check the error message above for specific details",
						)
					}
				}

				return cliErr
			}

			// Generic generation error
			return cli.NewGenerationError("unknown", "project generation failed", err).
				WithSuggestions(
					"Check the error messages above for specific component failures",
					"Run with --verbose flag for detailed execution logs",
					"Try using --no-external-tools for fallback generation",
					"Verify all required tools are installed: generator check-tools",
				)
		}

		log.Debug("Generation completed successfully")

		// Display results
		displayResults(result.(*models.GenerationResult), log)
	}

	duration := time.Since(startTime)
	log.Info(fmt.Sprintf("Completed in %v", duration))

	if verbose {
		log.Debug(fmt.Sprintf("Total execution time: %v", duration))
		log.Debug(fmt.Sprintf("Finished at: %s", time.Now().Format(time.RFC3339)))
	}

	return nil
}

// loadConfigFile loads a project configuration from a YAML or JSON file
func loadConfigFile(path string) (*models.ProjectConfig, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, cli.NewFileSystemError("read", path, err).
				WithSuggestions(
					"Verify the file path is correct",
					"Create a config file: generator init-config "+path,
					"Use absolute path if the file is in a different directory",
				)
		}
		if os.IsPermission(err) {
			return nil, cli.NewFileSystemError("read", path, err).
				WithSuggestions(
					"Check file permissions: ls -l "+path,
					"Grant read permissions: chmod +r "+path,
				)
		}
		return nil, cli.NewFileSystemError("read", path, err).
			WithSuggestions(
				"Verify the file exists and is readable",
				"Check file system permissions",
			)
	}

	// Parse YAML/JSON
	var config models.ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, cli.NewConfigError("failed to parse config file", err).
			WithSuggestions(
				"Check YAML syntax is valid (indentation, colons, etc.)",
				"Validate YAML online: https://www.yamllint.com/",
				"Use 'generator init-config' to see a valid example",
				"Common issues: incorrect indentation, missing colons, invalid characters",
			)
	}

	// Validate required fields
	if config.Name == "" {
		return nil, cli.NewConfigError("project name is required", nil).
			WithSuggestions(
				"Add 'name: your-project-name' to the config file",
				"See example: generator init-config --example fullstack",
			)
	}

	// Set default output directory if not specified
	if config.OutputDir == "" {
		config.OutputDir = filepath.Join(".", config.Name)
	}

	return &config, nil
}

// runInteractiveMode runs the interactive configuration wizard
func runInteractiveMode(log *logger.Logger) (*models.ProjectConfig, error) {
	// Create the interactive wizard
	wizard := interactivemode.NewInteractiveWizard(log)

	// Run the wizard
	ctx := context.Background()
	config, err := wizard.Run(ctx)
	if err != nil {
		// Check if user cancelled
		if strings.Contains(err.Error(), "cancelled") || strings.Contains(err.Error(), "interrupt") {
			return nil, cli.ErrUserCancelled
		}
		return nil, fmt.Errorf("config error: %w", err)
	}

	return config, nil
}

// displayPreview displays the dry-run preview results
func displayPreview(preview *models.PreviewResult, log *logger.Logger) {
	fmt.Println("\n" + separator("="))
	fmt.Println("DRY RUN PREVIEW")
	fmt.Println(separator("="))

	fmt.Printf("\nProject Root: %s\n", preview.ProjectRoot)

	fmt.Println("\nComponents:")
	fmt.Println(separator("-"))
	for i, comp := range preview.Components {
		fmt.Printf("\n%d. %s (%s)\n", i+1, comp.Name, comp.Type)
		fmt.Printf("   Target Path: %s\n", comp.TargetPath)
		fmt.Printf("   Method: %s\n", comp.Method)
		fmt.Printf("   Tool: %s\n", comp.ToolUsed)

		if len(comp.Files) > 0 {
			fmt.Printf("   Expected Files: %d files\n", len(comp.Files))
		}

		if len(comp.Warnings) > 0 {
			fmt.Println("   Warnings:")
			for _, warning := range comp.Warnings {
				fmt.Printf("     ⚠ %s\n", warning)
			}
		}
	}

	if len(preview.Structure) > 0 {
		fmt.Println("\nDirectory Structure:")
		fmt.Println(separator("-"))
		for _, dir := range preview.Structure {
			fmt.Printf("  📁 %s\n", dir)
		}
	}

	if len(preview.Files) > 0 {
		fmt.Println("\nIntegration Files:")
		fmt.Println(separator("-"))
		for _, file := range preview.Files {
			fmt.Printf("  📄 %s\n", file)
		}
	}

	if len(preview.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		fmt.Println(separator("-"))
		for _, warning := range preview.Warnings {
			fmt.Printf("  ⚠ %s\n", warning)
		}
	}

	fmt.Println("\n" + separator("="))
	fmt.Println("This was a dry run. No files were created.")
	fmt.Println("Run without --dry-run to generate the project.")
	fmt.Println(separator("=") + "\n")
}

// displayResults displays the generation results
func displayResults(result *models.GenerationResult, log *logger.Logger) {
	fmt.Println("\n" + separator("="))
	if result.Success {
		fmt.Println("✓ PROJECT GENERATION SUCCESSFUL")
	} else {
		fmt.Println("✗ PROJECT GENERATION FAILED")
	}
	fmt.Println(separator("="))

	fmt.Printf("\nProject Root: %s\n", result.ProjectRoot)
	fmt.Printf("Duration: %v\n", result.Duration)

	fmt.Println("\nComponents:")
	fmt.Println(separator("-"))
	for i, comp := range result.Components {
		status := "✓"
		if !comp.Success {
			status = "✗"
		}

		fmt.Printf("\n%s %d. %s (%s)\n", status, i+1, comp.Name, comp.Type)
		fmt.Printf("   Method: %s\n", comp.Method)
		fmt.Printf("   Tool: %s\n", comp.ToolUsed)
		fmt.Printf("   Output: %s\n", comp.OutputPath)
		fmt.Printf("   Duration: %v\n", comp.Duration)

		if comp.Error != nil {
			fmt.Printf("   Error: %v\n", comp.Error)
		}

		if len(comp.ManualSteps) > 0 {
			fmt.Println("   Manual Steps Required:")
			for j, step := range comp.ManualSteps {
				fmt.Printf("     %d. %s\n", j+1, step)
			}
		}

		if len(comp.Warnings) > 0 {
			fmt.Println("   Warnings:")
			for _, warning := range comp.Warnings {
				fmt.Printf("     ⚠ %s\n", warning)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		fmt.Println(separator("-"))
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		fmt.Println(separator("-"))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
	}

	if result.Success {
		fmt.Println("\n" + separator("="))
		fmt.Println("Next Steps:")
		fmt.Println(separator("-"))
		fmt.Printf("1. Navigate to project: cd %s\n", result.ProjectRoot)
		fmt.Println("2. Read the README.md for setup instructions")
		fmt.Println("3. Review any manual steps listed above")
		fmt.Println("4. Start development!")
		fmt.Println(separator("=") + "\n")
	}
}

// separator creates a separator line
func separator(char string) string {
	return repeatString(char, 70)
}

// repeatString repeats a string n times
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

var cacheToolsCmd = &cobra.Command{
	Use:   "cache-tools",
	Short: "Manage tool cache for offline use",
	Long: `Manage the tool availability cache for offline project generation.

This command allows you to:
- View cache statistics and status
- Clear the cache to force fresh tool checks
- Save current tool availability for offline use
- Display information about which tools can be cached

The cache stores tool availability and version information to speed up
subsequent tool checks and enable offline operation.

Examples:
  # View cache statistics
  generator cache-tools --stats

  # Clear the cache
  generator cache-tools --clear

  # Save current tool availability
  generator cache-tools --save

  # Show cache file location
  generator cache-tools --info`,
	RunE: runCacheTools,
}

var (
	showCacheStats bool
	clearCache     bool
	saveCache      bool
	showCacheInfo  bool
	validateCache  bool
	refreshCache   bool
	exportCache    string
	importCache    string
)

func init() {
	cacheToolsCmd.Flags().BoolVar(&showCacheStats, "stats", false, "Show cache statistics")
	cacheToolsCmd.Flags().BoolVar(&clearCache, "clear", false, "Clear the tool cache")
	cacheToolsCmd.Flags().BoolVar(&saveCache, "save", false, "Save current tool availability to cache")
	cacheToolsCmd.Flags().BoolVar(&showCacheInfo, "info", false, "Show cache information and location")
	cacheToolsCmd.Flags().BoolVar(&validateCache, "validate", false, "Validate cache integrity")
	cacheToolsCmd.Flags().BoolVar(&refreshCache, "refresh", false, "Refresh all cached tools")
	cacheToolsCmd.Flags().StringVar(&exportCache, "export", "", "Export cache to file")
	cacheToolsCmd.Flags().StringVar(&importCache, "import", "", "Import cache from file")
	cacheToolsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runCacheTools(cmd *cobra.Command, args []string) error {
	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)
	}

	// Create tool discovery with cache
	toolDiscovery := orchestrator.NewToolDiscovery(log)

	// Get the cache from tool discovery
	toolCache, err := orchestrator.NewToolCache(orchestrator.DefaultToolCacheConfig(), log)
	if err != nil {
		return cli.NewFileSystemError("initialize", "cache", err).
			WithSuggestions(
				"Check that the cache directory is writable",
				"Verify disk space is available",
				"Try clearing the cache: rm -rf ~/.cache/generator",
			)
	}

	// Create cache manager
	cacheManager := cache.NewCacheManager(toolCache, log)

	// If no flags specified, show stats by default
	if !showCacheStats && !clearCache && !saveCache && !showCacheInfo && !validateCache && !refreshCache && exportCache == "" && importCache == "" {
		showCacheStats = true
	}

	// Handle import cache
	if importCache != "" {
		log.Info(fmt.Sprintf("Importing cache from %s...", importCache))
		if err := cacheManager.Import(importCache); err != nil {
			return cli.NewFileSystemError("import", importCache, err).
				WithSuggestions(
					"Verify the import file exists and is readable",
					"Check that the file is a valid cache export",
					"Try exporting from another system: generator cache-tools --export cache.json",
				)
		}
		fmt.Printf("✓ Cache imported successfully from %s\n", importCache)
		return nil
	}

	// Handle export cache
	if exportCache != "" {
		log.Info(fmt.Sprintf("Exporting cache to %s...", exportCache))
		if err := cacheManager.Export(exportCache); err != nil {
			return cli.NewFileSystemError("export", exportCache, err).
				WithSuggestions(
					"Check that you have write permissions in the target directory",
					"Verify disk space is available",
					"Try a different output path",
				)
		}
		fmt.Printf("✓ Cache exported successfully to %s\n", exportCache)
		return nil
	}

	// Handle validate cache
	if validateCache {
		log.Info("Validating cache integrity...")
		report, err := cacheManager.Validate()
		if err != nil {
			return cli.NewCLIError("cache", "cache validation failed", err).
				WithSuggestions(
					"Try clearing the cache: generator cache-tools --clear",
					"Re-save the cache: generator cache-tools --save",
					"Check cache file permissions",
				)
		}

		fmt.Println("\n" + separator("="))
		fmt.Println("CACHE VALIDATION REPORT")
		fmt.Println(separator("="))

		if report.Valid {
			fmt.Println("\n✓ Cache is valid")
		} else {
			fmt.Println("\n✗ Cache validation failed")
		}

		fmt.Printf("\nTotal Entries: %d\n", report.TotalEntries)
		fmt.Printf("Corrupted Entries: %d\n", len(report.CorruptedEntries))
		fmt.Printf("Expired Entries: %d\n", len(report.ExpiredEntries))
		fmt.Printf("Checked At: %s\n", report.CheckedAt.Format(time.RFC3339))

		if len(report.Warnings) > 0 {
			fmt.Println("\nWarnings:")
			for _, warning := range report.Warnings {
				fmt.Printf("  - %s\n", warning)
			}
		}

		if len(report.CorruptedEntries) > 0 {
			fmt.Println("\nCorrupted Entries:")
			for _, entry := range report.CorruptedEntries {
				fmt.Printf("  - %s\n", entry)
			}
		}

		fmt.Println("\n" + separator("=") + "\n")
		return nil
	}

	// Handle refresh cache
	if refreshCache {
		log.Info("Refreshing cache...")
		if err := cacheManager.Refresh(toolDiscovery); err != nil {
			return cli.NewCLIError("cache", "failed to refresh cache", err).
				WithSuggestions(
					"Check your network connection if tools need to be checked",
					"Try clearing and re-saving the cache",
					"Run with --verbose flag for detailed error information",
				)
		}
		fmt.Println("✓ Cache refreshed successfully")

		// Show updated stats
		stats := cacheManager.GetStats()
		fmt.Printf("\nTotal Entries: %d\n", stats.TotalEntries)
		fmt.Printf("Available Tools: %d\n", stats.AvailableTools)
		fmt.Printf("Unavailable Tools: %d\n", stats.UnavailableTools)
		return nil
	}

	// Handle clear cache
	if clearCache {
		log.Info("Clearing tool cache...")
		toolDiscovery.ClearCache()
		fmt.Println("✓ Tool cache cleared successfully")
		return nil
	}

	// Handle save cache
	if saveCache {
		log.Info("Checking and caching all registered tools...")

		// Get all registered tools
		allTools := toolDiscovery.ListRegisteredTools()

		// Check each tool to populate cache
		for _, toolName := range allTools {
			log.Debug(fmt.Sprintf("Checking tool: %s", toolName))
			available, _ := toolDiscovery.IsAvailable(toolName)
			if available {
				// Also get version to cache it
				toolDiscovery.GetVersion(toolName)
			}
		}

		// Save cache to disk
		if err := toolDiscovery.SaveCache(); err != nil {
			return cli.NewFileSystemError("save", "cache", err).
				WithSuggestions(
					"Check that the cache directory is writable",
					"Verify disk space is available",
					"Check file system permissions",
				)
		}

		fmt.Printf("✓ Tool availability cached for %d tools\n", len(allTools))
		fmt.Println("\nCached tools can be used for offline project generation.")
		fmt.Println("Note: External bootstrap tools themselves are not cached,")
		fmt.Println("only their availability status. For true offline operation,")
		fmt.Println("ensure tools are installed before going offline.")
		return nil
	}

	// Handle show info
	if showCacheInfo {
		stats := cacheManager.GetStats()

		fmt.Println("\n" + separator("="))
		fmt.Println("TOOL CACHE INFORMATION")
		fmt.Println(separator("="))

		fmt.Printf("\nCache File: %s\n", stats.CacheFile)
		fmt.Printf("Cache TTL: %s\n", stats.TTL.String())
		fmt.Printf("Last Saved: %s\n", stats.LastSaved.Format(time.RFC3339))

		fmt.Println("\nCacheable Information:")
		fmt.Println("  - Tool availability (whether tool is in PATH)")
		fmt.Println("  - Tool versions")
		fmt.Println("  - Last check timestamp")

		fmt.Println("\nNot Cached:")
		fmt.Println("  - Tool executables themselves")
		fmt.Println("  - Tool dependencies")
		fmt.Println("  - Network-based resources")

		fmt.Println("\nOffline Operation:")
		fmt.Println("  - Cache stores tool availability for quick checks")
		fmt.Println("  - Fallback generators work fully offline")
		fmt.Println("  - Bootstrap tools require prior installation")

		fmt.Println("\n" + separator("=") + "\n")
		return nil
	}

	// Handle show stats (default)
	if showCacheStats {
		stats := cacheManager.GetStats()

		fmt.Println("\n" + separator("="))
		fmt.Println("TOOL CACHE STATISTICS")
		fmt.Println(separator("="))

		fmt.Printf("\nTotal Entries: %d\n", stats.TotalEntries)
		fmt.Printf("Available Tools: %d\n", stats.AvailableTools)
		fmt.Printf("Unavailable Tools: %d\n", stats.UnavailableTools)
		fmt.Printf("Expired Entries: %d\n", stats.ExpiredEntries)

		fmt.Printf("\nCache File: %s\n", stats.CacheFile)
		fmt.Printf("Cache TTL: %s\n", stats.TTL.String())
		fmt.Printf("Last Saved: %s\n", stats.LastSaved.Format(time.RFC3339))
		fmt.Printf("Time Since Check: %s\n", stats.TimeSinceCheck.String())

		fmt.Println("\nCommands:")
		fmt.Println("  --clear      Clear the cache")
		fmt.Println("  --save       Cache all registered tools")
		fmt.Println("  --info       Show cache information")
		fmt.Println("  --validate   Validate cache integrity")
		fmt.Println("  --refresh    Refresh all cached tools")
		fmt.Println("  --export     Export cache to file")
		fmt.Println("  --import     Import cache from file")

		fmt.Println("\n" + separator("=") + "\n")
	}

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Open Source Project Generator\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Built: %s\n", buildTime)
		fmt.Printf("Commit: %s\n", gitCommit)
	},
}
