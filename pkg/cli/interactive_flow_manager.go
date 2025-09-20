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
	// Start interactive UI session
	sessionConfig := interfaces.SessionConfig{
		Title:       "Project Generator",
		Description: "Interactive project configuration and generation",
		Timeout:     30 * 60, // 30 minutes
		AutoSave:    true,
	}

	session, err := ifm.ui.StartSession(ctx, sessionConfig)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Unable to start interactive session.",
			"Check if your terminal supports interactive mode")
	}
	defer func() {
		if endErr := ifm.ui.EndSession(ctx, session); endErr != nil {
			ifm.logger.Error("🖥️  Couldn't end UI session properly", "error", endErr)
		}
	}()

	// Step 1: Template Selection
	ifm.cli.VerboseOutput("🎯 Let's choose your project templates...")
	selectedTemplates, err := ifm.runTemplateSelection(ctx)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Template selection failed.",
			"Check if templates are available and accessible")
	}

	// Step 2: Project Configuration Collection
	ifm.cli.VerboseOutput("⚙️  Now let's configure your project details...")
	config, err := ifm.runProjectConfiguration(ctx, selectedTemplates)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Project configuration failed.",
			"Check your input values and try again")
	}

	// Step 3: Output Directory Selection
	ifm.cli.VerboseOutput("📁 Where would you like to create your project?")
	outputPath, err := ifm.runDirectorySelection(ctx, options.OutputPath, config.Name)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Directory selection failed.",
			"Check if the directory path is valid and accessible")
	}

	// Step 4: Project Structure Preview
	ifm.cli.VerboseOutput("👀 Preparing a preview of your project structure...")
	preview, err := ifm.runStructurePreview(ctx, config, selectedTemplates, outputPath)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Preview generation failed.",
			"Check if templates are valid and accessible")
	}

	// Step 5: Final Confirmation
	ifm.cli.VerboseOutput("✅ Ready to generate! Let's confirm everything looks good...")
	confirmed, err := ifm.runFinalConfirmation(ctx, config, preview, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			"Confirmation process failed.",
			"Check your terminal input capabilities")
	}

	if !confirmed {
		ifm.cli.QuietOutput("❌ Project generation cancelled")
		return nil
	}

	// Step 6: Configuration Persistence (optional)
	if err := ifm.runConfigurationPersistence(ctx, config, selectedTemplates); err != nil {
		ifm.logger.WarnWithFields("💾 Couldn't save your configuration", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Step 7: Project Generation
	ifm.cli.VerboseOutput("🚀 Generating your project now...")
	return ifm.runProjectGeneration(ctx, config, selectedTemplates, outputPath, options)
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
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Directory Selection"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Determine default path
	if defaultPath == "" {
		if projectName != "" {
			defaultPath = "output/generated/" + projectName
		} else {
			defaultPath = "output/generated"
		}
	}

	// Use existing directory selection implementation
	return ifm.cli.runInteractiveDirectorySelection(ctx, defaultPath)
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
func (ifm *InteractiveFlowManager) runFinalConfirmation(ctx context.Context, config *models.ProjectConfig, preview *interfaces.ProjectStructurePreview, options interfaces.GenerateOptions) (bool, error) {
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Final Confirmation"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Use existing confirmation method
	return ifm.cli.runInteractiveConfirmation(ctx, config, options), nil
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
	if err := ifm.ui.ShowBreadcrumb(ctx, []string{"Generator", "Generation"}); err != nil {
		ifm.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Use existing interactive generation method
	templateName := "go-gin"
	if len(templates) > 0 {
		templateName = templates[0].Template.Name
	}

	return ifm.cli.runInteractiveGeneration(ctx, templateName, config, outputPath)
}

// Helper methods

// formatSize formats a byte size into a human-readable string
// TODO: This method will be used for displaying file sizes in progress tracking
//
//nolint:unused // Will be used in future progress tracking implementation
func (ifm *InteractiveFlowManager) formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// generateProgressSteps creates progress steps based on selected templates
// TODO: This method will be used for progress tracking during generation
//
//nolint:unused // Will be used in future progress tracking implementation
func (ifm *InteractiveFlowManager) generateProgressSteps(templates []interfaces.TemplateSelection) []string {
	steps := []string{
		"Initializing generation",
		"Creating directory structure",
		"Processing templates",
		"Generating configuration files",
		"Setting up project files",
		"Finalizing project",
	}

	// Add template-specific steps
	for _, template := range templates {
		steps = append(steps, fmt.Sprintf("Processing %s template", template.Template.DisplayName))
	}

	return steps
}
