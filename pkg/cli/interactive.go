// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
//
// This file contains the InteractiveManager which handles all interactive user
// prompts and confirmations for project generation.
package cli

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InteractiveManager handles all interactive user prompts and confirmations.
//
// This struct provides methods for:
//   - Collecting project configuration from user input
//   - Confirming generation with configuration preview
//   - Advanced options prompting
//   - Template selection
//   - Non-interactive mode detection and handling
type InteractiveManager struct {
	cli           *CLI
	outputManager *OutputManager
	flagHandler   *FlagHandler
	logger        interfaces.Logger
}

// NewInteractiveManager creates a new interactive manager instance.
//
// Parameters:
//   - cli: Reference to the main CLI instance
//   - outputManager: For formatted output and color management
//   - flagHandler: For checking non-interactive mode
//   - logger: For logging operations
//
// Returns:
//   - *InteractiveManager: New interactive manager ready for use
func NewInteractiveManager(
	cli *CLI,
	outputManager *OutputManager,
	flagHandler *FlagHandler,
	logger interfaces.Logger,
) *InteractiveManager {
	return &InteractiveManager{
		cli:           cli,
		outputManager: outputManager,
		flagHandler:   flagHandler,
		logger:        logger,
	}
}

// PromptProjectDetails collects basic project configuration from user input.
//
// This method prompts the user for essential project information including:
//   - Project name (required)
//   - Organization (optional)
//   - Description (optional)
//   - Author (optional)
//   - License (defaults to MIT)
//
// Returns default component configuration with Go Gin backend, Next.js frontend,
// and Docker infrastructure enabled.
//
// Returns:
//   - *models.ProjectConfig: Collected project configuration
//   - error: If non-interactive mode or input reading fails
func (im *InteractiveManager) PromptProjectDetails() (*models.ProjectConfig, error) {
	if im.isNonInteractiveMode() {
		return nil, fmt.Errorf("🚫 %s %s",
			im.outputManager.GetColorManager().Error("Interactive prompts not available in non-interactive mode."),
			im.outputManager.GetColorManager().Info("Use environment variables or a configuration file instead"))
	}

	im.outputManager.QuietOutput("Project Configuration")
	im.outputManager.QuietOutput("====================")

	config := &models.ProjectConfig{}

	// Get project name
	fmt.Print("Project name: ")
	var name string
	if _, err := fmt.Scanln(&name); err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			im.outputManager.GetColorManager().Error("Unable to read project name."),
			im.outputManager.GetColorManager().Info("Please try typing it again or check your input"))
	}
	config.Name = strings.TrimSpace(name)

	// Get organization (optional)
	fmt.Print("Organization (optional): ")
	var org string
	_, _ = fmt.Scanln(&org) // Ignore error for optional input
	config.Organization = strings.TrimSpace(org)

	// Get description (optional)
	fmt.Print("Description (optional): ")
	var desc string
	_, _ = fmt.Scanln(&desc) // Ignore error for optional input
	config.Description = strings.TrimSpace(desc)

	// Get author (optional)
	fmt.Print("Author (optional): ")
	var author string
	_, _ = fmt.Scanln(&author) // Ignore error for optional input
	config.Author = strings.TrimSpace(author)

	// Get license (default: MIT)
	fmt.Print("License (default: MIT): ")
	var license string
	_, _ = fmt.Scanln(&license) // Ignore error for optional input
	if strings.TrimSpace(license) == "" {
		license = "MIT"
	}
	config.License = strings.TrimSpace(license)

	// Set default components
	config.Components = models.Components{
		Backend: models.BackendComponents{
			GoGin: true,
		},
		Frontend: models.FrontendComponents{
			NextJS: models.NextJSComponents{
				App: true,
			},
		},
		Infrastructure: models.InfrastructureComponents{
			Docker: true,
		},
	}

	return config, nil
}

// ConfirmGeneration shows a basic configuration preview and asks for user confirmation.
//
// This method displays a summary of the project configuration including:
//   - Project name, organization, description, author, license
//   - Enabled components (backend, frontend, infrastructure)
//
// In non-interactive mode, automatically returns true.
//
// Parameters:
//   - config: Project configuration to preview
//
// Returns:
//   - bool: True if user confirms generation, false otherwise
func (im *InteractiveManager) ConfirmGeneration(config *models.ProjectConfig) bool {
	if im.isNonInteractiveMode() {
		return true // Auto-confirm in non-interactive mode
	}

	im.outputManager.QuietOutput("\nProject Configuration Preview:")
	im.outputManager.QuietOutput("==============================")
	im.outputManager.QuietOutput("Name: %s", config.Name)
	if config.Organization != "" {
		im.outputManager.QuietOutput("Organization: %s", config.Organization)
	}
	if config.Description != "" {
		im.outputManager.QuietOutput("Description: %s", config.Description)
	}
	if config.Author != "" {
		im.outputManager.QuietOutput("Author: %s", config.Author)
	}
	im.outputManager.QuietOutput("License: %s", config.License)

	im.outputManager.QuietOutput("\nComponents:")
	if config.Components.Backend.GoGin {
		im.outputManager.QuietOutput("  - Go Gin API")
	}
	if config.Components.Frontend.NextJS.App {
		im.outputManager.QuietOutput("  - Next.js App")
	}
	if config.Components.Infrastructure.Docker {
		im.outputManager.QuietOutput("  - Docker configuration")
	}

	fmt.Print("\nProceed with generation? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

// PromptAdvancedOptions collects advanced configuration options from user input.
//
// This method prompts for advanced options including:
//   - Security scanning (default: enabled)
//   - Quality checks (default: enabled)
//   - Performance optimization (default: disabled)
//   - Documentation generation (default: enabled)
//   - CI/CD setup (default: enabled with GitHub Actions)
//
// In non-interactive mode, returns sensible defaults.
//
// Returns:
//   - *interfaces.AdvancedOptions: Collected advanced options
//   - error: If non-interactive mode detection fails
func (im *InteractiveManager) PromptAdvancedOptions() (*interfaces.AdvancedOptions, error) {
	if im.isNonInteractiveMode() {
		// Return default advanced options in non-interactive mode
		return &interfaces.AdvancedOptions{
			EnableSecurityScanning:        true,
			EnableQualityChecks:           true,
			EnablePerformanceOptimization: false,
			GenerateDocumentation:         true,
			EnableCICD:                    true,
			CICDProviders:                 []string{"github-actions"},
			EnableMonitoring:              false,
		}, nil
	}

	im.outputManager.QuietOutput("Advanced Options")
	im.outputManager.QuietOutput("================")

	options := &interfaces.AdvancedOptions{}

	// Security options
	fmt.Print("Enable security scanning? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	options.EnableSecurityScanning = strings.ToLower(strings.TrimSpace(response)) != "n"

	// Quality options
	fmt.Print("Enable quality checks? (Y/n): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	options.EnableQualityChecks = strings.ToLower(strings.TrimSpace(response)) != "n"

	// Performance options
	fmt.Print("Enable performance optimization? (y/N): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	options.EnablePerformanceOptimization = strings.ToLower(strings.TrimSpace(response)) == "y"

	// Documentation options
	fmt.Print("Generate documentation? (Y/n): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	options.GenerateDocumentation = strings.ToLower(strings.TrimSpace(response)) != "n"

	// CI/CD options
	fmt.Print("Enable CI/CD? (Y/n): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	options.EnableCICD = strings.ToLower(strings.TrimSpace(response)) != "n"

	if options.EnableCICD {
		options.CICDProviders = []string{"github-actions"}
	}

	return options, nil
}

// ConfirmAdvancedGeneration shows advanced configuration preview and asks for confirmation.
//
// This method displays a summary of the advanced options including:
//   - Security scanning status
//   - Quality checks status
//   - Performance optimization status
//   - Documentation generation status
//   - CI/CD configuration
//
// In non-interactive mode, automatically returns true.
//
// Parameters:
//   - config: Project configuration (for context)
//   - options: Advanced options to preview
//
// Returns:
//   - bool: True if user confirms advanced generation, false otherwise
func (im *InteractiveManager) ConfirmAdvancedGeneration(config *models.ProjectConfig, options *interfaces.AdvancedOptions) bool {
	if im.isNonInteractiveMode() {
		return true // Auto-confirm in non-interactive mode
	}

	im.outputManager.QuietOutput("\nAdvanced Configuration Preview:")
	im.outputManager.QuietOutput("===============================")
	im.outputManager.QuietOutput("Security Scanning: %t", options.EnableSecurityScanning)
	im.outputManager.QuietOutput("Quality Checks: %t", options.EnableQualityChecks)
	im.outputManager.QuietOutput("Performance Optimization: %t", options.EnablePerformanceOptimization)
	im.outputManager.QuietOutput("Generate Documentation: %t", options.GenerateDocumentation)
	im.outputManager.QuietOutput("Enable CI/CD: %t", options.EnableCICD)
	if options.EnableCICD {
		im.outputManager.QuietOutput("CI/CD Providers: %v", options.CICDProviders)
	}

	fmt.Print("\nProceed with advanced generation? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

// SelectTemplateInteractively presents available templates and allows user selection.
//
// This method:
//   - Lists available templates matching the filter
//   - Shows template details (name, description, category, technology)
//   - Prompts user to select a template by number
//   - Validates the selection
//
// In non-interactive mode, returns an error as template selection requires user input.
//
// Parameters:
//   - filter: Template filter criteria
//
// Returns:
//   - *interfaces.TemplateInfo: Selected template information
//   - error: If non-interactive mode, no templates found, or invalid selection
func (im *InteractiveManager) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	if im.isNonInteractiveMode() {
		return nil, fmt.Errorf("interactive template selection not available in non-interactive mode")
	}

	// Get available templates
	templates, err := im.cli.templateManager.ListTemplates(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found matching the criteria")
	}

	im.outputManager.QuietOutput("Available Templates:")
	im.outputManager.QuietOutput("===================")

	for i, template := range templates {
		im.outputManager.QuietOutput("%d. %s - %s", i+1, template.DisplayName, template.Description)
		im.outputManager.VerboseOutput("   Category: %s, Technology: %s", template.Category, template.Technology)
	}

	fmt.Printf("\nSelect template (1-%d): ", len(templates))
	var selection int
	if _, err := fmt.Scanln(&selection); err != nil {
		return nil, fmt.Errorf("failed to read template selection: %w", err)
	}

	if selection < 1 || selection > len(templates) {
		return nil, fmt.Errorf("🚫 %s %s",
			im.outputManager.GetColorManager().Error("Invalid selection."),
			im.outputManager.GetColorManager().Info(fmt.Sprintf("Please choose a number between 1 and %d", len(templates))))
	}

	return &templates[selection-1], nil
}

// isNonInteractiveMode checks if the CLI is running in non-interactive mode.
//
// This is a helper method that delegates to the flag handler to determine
// if interactive prompts should be disabled.
//
// Returns:
//   - bool: True if non-interactive mode is enabled
func (im *InteractiveManager) isNonInteractiveMode() bool {
	return im.flagHandler.IsNonInteractiveMode(im.cli.rootCmd)
}
