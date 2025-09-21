// Package cli provides interactive flow management for the CLI generator.
//
// This file implements the InteractiveFlowManager that orchestrates the complete
// interactive project generation workflow including template selection, configuration
// collection, directory selection, preview generation, and confirmation.
package cli

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InteractiveFlowManager orchestrates the complete interactive generation workflow
type InteractiveFlowManager struct {
	cli             *CLI
	templateManager interfaces.TemplateManager
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	logger          interfaces.Logger
	ui              interfaces.InteractiveUIInterface
}

// NewInteractiveFlowManager creates a new interactive flow manager
func NewInteractiveFlowManager(
	cli *CLI,
	templateManager interfaces.TemplateManager,
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	logger interfaces.Logger,
	ui interfaces.InteractiveUIInterface,
) *InteractiveFlowManager {
	return &InteractiveFlowManager{
		cli:             cli,
		templateManager: templateManager,
		configManager:   configManager,
		validator:       validator,
		logger:          logger,
		ui:              ui,
	}
}

// RunInteractiveFlow executes the complete interactive generation workflow
func (ifm *InteractiveFlowManager) RunInteractiveFlow(ctx context.Context, options interfaces.GenerateOptions) error {
	fmt.Println("🚀 Project Generator")
	fmt.Println()

	// Step 1: Project Configuration Collection
	config, err := ifm.runProjectConfiguration(ctx, nil)
	if err != nil {
		return fmt.Errorf("project configuration failed: %w", err)
	}

	// Step 2: Output Directory Selection
	outputPath, err := ifm.runDirectorySelection(ctx, options.OutputPath, config.Name)
	if err != nil {
		return fmt.Errorf("directory selection failed: %w", err)
	}

	// Step 3: Final Confirmation
	confirmed := ifm.runFinalConfirmation(ctx, config, nil, options)
	if !confirmed {
		fmt.Println("❌ Project generation cancelled")
		return nil
	}

	// Step 4: Project Generation
	return ifm.runProjectGeneration(ctx, config, nil, outputPath, options)
}

// runTemplateSelection handles interactive template selection
func (ifm *InteractiveFlowManager) runTemplateSelection(ctx context.Context) ([]interfaces.TemplateSelection, error) {
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Template Selection"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// For now, return a default template selection
	// This will be enhanced in future tasks when template selection UI is implemented
	defaultTemplate := interfaces.TemplateSelection{
		Template: interfaces.TemplateInfo{
			Name:        "go-gin",
			DisplayName: "Go Gin API",
			Description: "RESTful API server using Gin framework",
			Category:    "backend",
			Technology:  "go",
			Version:     "1.0.0",
		},
		Selected: true,
		Options:  make(map[string]interface{}),
	}

	return []interfaces.TemplateSelection{defaultTemplate}, nil
}

// runProjectConfiguration handles interactive project configuration collection
func (ifm *InteractiveFlowManager) runProjectConfiguration(ctx context.Context, templates []interfaces.TemplateSelection) (*models.ProjectConfig, error) {
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Project Configuration"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Use existing interactive project configuration method
	return ifm.cli.runInteractiveProjectConfiguration(ctx)
}

// runDirectorySelection handles interactive output directory selection
func (ifm *InteractiveFlowManager) runDirectorySelection(ctx context.Context, defaultPath, projectName string) (string, error) {
	// Determine default path - just the base directory, not including project name
	if defaultPath == "" {
		defaultPath = "output/generated"
	}

	// Simple text prompt for output directory
	dirConfig := interfaces.TextPromptConfig{
		Prompt:       "Output Directory",
		Description:  "Enter the base path where your project should be generated",
		Required:     false,
		DefaultValue: defaultPath,
	}

	dirResult, err := ifm.ui.PromptText(ctx, dirConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get output directory: %w", err)
	}
	if dirResult.Cancelled {
		return "", fmt.Errorf("directory selection cancelled")
	}

	outputPath := dirResult.Value
	if outputPath == "" {
		outputPath = defaultPath
	}

	return outputPath, nil
}

// runStructurePreview handles project structure preview generation and display
func (ifm *InteractiveFlowManager) runStructurePreview(ctx context.Context, config *models.ProjectConfig, templates []interfaces.TemplateSelection, outputPath string) (*interfaces.ProjectStructurePreview, error) {
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Project Preview"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Create a basic preview structure
	// This will be enhanced in future tasks when preview generation is fully implemented
	preview := &interfaces.ProjectStructurePreview{
		RootDirectory: outputPath,
		Structure: []interfaces.DirectoryNode{
			{
				Name: config.Name,
				Type: "directory",
				Children: []interfaces.DirectoryNode{
					{Name: "cmd", Type: "directory", Source: "go-gin"},
					{Name: "internal", Type: "directory", Source: "go-gin"},
					{Name: "pkg", Type: "directory", Source: "go-gin"},
					{Name: "README.md", Type: "file", Source: "base"},
					{Name: "go.mod", Type: "file", Source: "go-gin"},
				},
			},
		},
		FileCount:     15,
		EstimatedSize: 1024 * 50, // 50KB estimate
		Components: []interfaces.ComponentSummary{
			{
				Name:        "Go Gin API",
				Type:        "backend",
				Description: "RESTful API server",
				Files:       []string{"main.go", "go.mod", "internal/", "pkg/"},
			},
		},
	}

	// Display preview using table
	tableConfig := interfaces.TableConfig{
		Title:   "Project Structure Preview",
		Headers: []string{"Component", "Type", "Files"},
		Rows: [][]string{
			{"Go Gin API", "Backend", "15 files"},
			{"Base Files", "Infrastructure", "README.md, LICENSE"},
		},
		MaxWidth:   80,
		Pagination: false,
	}

	if err := ifm.ui.ShowTable(ctx, tableConfig); err != nil {
		ifm.logger.WarnWithFields("📋 Couldn't display preview table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return preview, nil
}

// runFinalConfirmation handles final confirmation before generation
func (ifm *InteractiveFlowManager) runFinalConfirmation(ctx context.Context, config *models.ProjectConfig, preview *interfaces.ProjectStructurePreview, options interfaces.GenerateOptions) bool {
	// Show configuration summary
	fmt.Println("\n📋 Project Summary")
	fmt.Println("==================")
	fmt.Printf("Name: %s\n", config.Name)
	if config.Organization != "" {
		fmt.Printf("Organization: %s\n", config.Organization)
	}
	if config.Description != "" {
		fmt.Printf("Description: %s\n", config.Description)
	}
	if config.Author != "" {
		fmt.Printf("Author: %s\n", config.Author)
	}
	fmt.Printf("License: %s\n", config.License)

	// Show selected components
	fmt.Println("\nComponents:")
	if config.Components.Backend.GoGin {
		fmt.Println("  ✅ Go Gin API")
	}
	if config.Components.Frontend.NextJS.App {
		fmt.Println("  ✅ Next.js Frontend")
	}
	if config.Components.Database.PostgreSQL {
		fmt.Println("  ✅ PostgreSQL Database")
	}
	if config.Components.Cache.Redis {
		fmt.Println("  ✅ Redis Cache")
	}
	if config.Components.Infrastructure.Docker {
		fmt.Println("  ✅ Docker Configuration")
	}
	if config.Components.Infrastructure.Kubernetes {
		fmt.Println("  ✅ Kubernetes Manifests")
	}

	// Simple confirmation prompt
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Generate Project",
		Description:  "Proceed with generating the project?",
		DefaultValue: true,
		YesLabel:     "Generate",
		NoLabel:      "Cancel",
	}

	confirmResult, err := ifm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false
	}

	return confirmResult.Confirmed && !confirmResult.Cancelled
}

// runConfigurationPersistence handles saving configuration for reuse
func (ifm *InteractiveFlowManager) runConfigurationPersistence(ctx context.Context, config *models.ProjectConfig, templates []interfaces.TemplateSelection) error {
	// For now, skip configuration persistence
	// This will be implemented in future tasks when configuration management is enhanced
	ifm.cli.VerboseOutput("💾 Configuration persistence will be available in a future update")
	return nil
}

// runProjectGeneration handles the actual project generation with progress tracking
func (ifm *InteractiveFlowManager) runProjectGeneration(ctx context.Context, config *models.ProjectConfig, templates []interfaces.TemplateSelection, outputPath string, options interfaces.GenerateOptions) error {
	// Set the output path in options
	options.OutputPath = outputPath

	// Use the same core generation workflow as non-interactive mode
	return ifm.cli.executeGenerationWorkflow(config, options)
}

// Helper methods
