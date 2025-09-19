// Package ui provides interactive flow management for the CLI generator.
//
// This file implements the InteractiveFlowManager which orchestrates the complete
// interactive project generation workflow including template selection, configuration
// collection, directory selection, and preview generation.
package ui

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InteractiveFlowManager orchestrates the interactive project generation workflow
type InteractiveFlowManager struct {
	ui              interfaces.InteractiveUIInterface
	templateManager interfaces.TemplateManager
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	logger          interfaces.Logger

	// Component managers
	templateSelector     *TemplateSelector
	configCollector      *ProjectConfigCollector
	directorySelector    *DirectorySelector
	previewGenerator     *PreviewGenerator
	configurationManager *InteractiveConfigurationManager
}

// InteractiveGenerationConfig contains configuration for interactive generation
type InteractiveGenerationConfig struct {
	DefaultOutputPath string
	AllowBack         bool
	AllowQuit         bool
	ShowHelp          bool
	AutoSave          bool
	SkipPreview       bool
}

// InteractiveGenerationResult contains the result of interactive generation
type InteractiveGenerationResult struct {
	ProjectConfig      *models.ProjectConfig
	SelectedTemplates  []TemplateSelection
	OutputDirectory    *DirectorySelectionResult
	PreviewAccepted    bool
	Cancelled          bool
	SavedConfiguration string
}

// PreviewGenerator handles project structure preview generation
type PreviewGenerator struct {
	ui     interfaces.InteractiveUIInterface
	logger interfaces.Logger
}

// NewInteractiveFlowManager creates a new interactive flow manager
func NewInteractiveFlowManager(
	ui interfaces.InteractiveUIInterface,
	templateManager interfaces.TemplateManager,
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	logger interfaces.Logger,
) *InteractiveFlowManager {
	// Get configuration directory from config manager or use default
	configDir := ""
	if configManager != nil {
		configDir = configManager.GetConfigLocation()
	}

	return &InteractiveFlowManager{
		ui:                   ui,
		templateManager:      templateManager,
		configManager:        configManager,
		validator:            validator,
		logger:               logger,
		templateSelector:     NewTemplateSelector(ui, templateManager, logger),
		configCollector:      NewProjectConfigCollector(ui, logger),
		directorySelector:    NewDirectorySelector(ui, logger),
		previewGenerator:     NewPreviewGenerator(ui, logger),
		configurationManager: NewInteractiveConfigurationManager(ui, configDir, logger),
	}
}

// RunInteractiveGeneration runs the complete interactive generation workflow
func (ifm *InteractiveFlowManager) RunInteractiveGeneration(ctx context.Context, config *InteractiveGenerationConfig) (*InteractiveGenerationResult, error) {
	if config == nil {
		config = &InteractiveGenerationConfig{
			DefaultOutputPath: "output/generated",
			AllowBack:         true,
			AllowQuit:         true,
			ShowHelp:          true,
			AutoSave:          true,
			SkipPreview:       false,
		}
	}

	result := &InteractiveGenerationResult{}

	// Step 1: Template Selection (placeholder - already implemented in other tasks)
	ifm.logger.InfoWithFields("Starting template selection", map[string]interface{}{
		"step": "template_selection",
	})

	templates, err := ifm.selectTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("template selection failed: %w", err)
	}
	if templates == nil {
		result.Cancelled = true
		return result, nil
	}
	result.SelectedTemplates = templates

	// Step 2: Project Configuration Collection
	ifm.logger.InfoWithFields("Starting project configuration collection", map[string]interface{}{
		"step": "config_collection",
	})

	projectConfig, err := ifm.configCollector.CollectProjectConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("project configuration collection failed: %w", err)
	}
	if projectConfig == nil {
		result.Cancelled = true
		return result, nil
	}
	result.ProjectConfig = projectConfig

	// Step 3: Output Directory Selection
	ifm.logger.InfoWithFields("Starting output directory selection", map[string]interface{}{
		"step":         "directory_selection",
		"default_path": config.DefaultOutputPath,
	})

	directoryResult, err := ifm.directorySelector.SelectOutputDirectory(ctx, config.DefaultOutputPath)
	if err != nil {
		return nil, fmt.Errorf("output directory selection failed: %w", err)
	}
	if directoryResult.Cancelled {
		result.Cancelled = true
		return result, nil
	}
	result.OutputDirectory = directoryResult

	// Step 4: Project Structure Preview (if not skipped)
	if !config.SkipPreview {
		ifm.logger.InfoWithFields("Starting project structure preview", map[string]interface{}{
			"step": "preview_generation",
		})

		previewAccepted, err := ifm.generateAndShowPreview(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("preview generation failed: %w", err)
		}
		if !previewAccepted {
			result.Cancelled = true
			return result, nil
		}
		result.PreviewAccepted = previewAccepted
	}

	// Step 5: Configuration Saving (if enabled)
	if config.AutoSave {
		ifm.logger.InfoWithFields("Offering to save configuration", map[string]interface{}{
			"step": "config_saving",
		})

		savedConfig, err := ifm.configurationManager.SaveConfigurationInteractively(
			ctx,
			result.ProjectConfig,
			result.SelectedTemplates,
		)
		if err != nil {
			ifm.logger.WarnWithFields("Failed to save configuration", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			result.SavedConfiguration = savedConfig
		}
	}

	ifm.logger.InfoWithFields("Interactive generation completed successfully", map[string]interface{}{
		"project_name":     result.ProjectConfig.Name,
		"output_directory": result.OutputDirectory.Path,
		"templates_count":  len(result.SelectedTemplates),
		"saved_config":     result.SavedConfiguration,
	})

	return result, nil
}

// selectTemplates handles template selection (placeholder implementation)
func (ifm *InteractiveFlowManager) selectTemplates(ctx context.Context) ([]TemplateSelection, error) {
	// This is a placeholder implementation since template selection is handled in other tasks
	// For now, we'll return a default selection to satisfy the interface

	// In a real implementation, this would call the template selector
	// return ifm.templateSelector.SelectTemplatesInteractively(ctx)

	// Placeholder: return a basic Go Gin template selection
	defaultTemplate := TemplateSelection{
		Template: interfaces.TemplateInfo{
			Name:        "go-gin",
			DisplayName: "Go Gin API",
			Description: "RESTful API server using Go and Gin framework",
			Category:    "backend",
			Technology:  "go",
			Version:     "1.0.0",
		},
		Selected: true,
		Options:  make(map[string]interface{}),
	}

	return []TemplateSelection{defaultTemplate}, nil
}

// generateAndShowPreview generates and displays the project structure preview
func (ifm *InteractiveFlowManager) generateAndShowPreview(ctx context.Context, result *InteractiveGenerationResult) (bool, error) {
	return ifm.previewGenerator.GenerateAndShowPreview(ctx, result)
}

// LoadConfigurationInteractively loads a saved configuration interactively
func (ifm *InteractiveFlowManager) LoadConfigurationInteractively(ctx context.Context) (*LoadedConfiguration, error) {
	return ifm.configurationManager.LoadConfigurationInteractively(ctx, &ConfigurationLoadOptions{
		AllowModification: true,
		ShowPreview:       true,
		ConfirmLoad:       true,
	})
}

// ManageConfigurationsInteractively provides configuration management interface
func (ifm *InteractiveFlowManager) ManageConfigurationsInteractively(ctx context.Context) error {
	return ifm.configurationManager.ManageConfigurationsInteractively(ctx)
}

// RunInteractiveGenerationFromConfig runs generation using a loaded configuration
func (ifm *InteractiveFlowManager) RunInteractiveGenerationFromConfig(
	ctx context.Context,
	loadedConfig *LoadedConfiguration,
	generationConfig *InteractiveGenerationConfig,
) (*InteractiveGenerationResult, error) {
	if generationConfig == nil {
		generationConfig = &InteractiveGenerationConfig{
			DefaultOutputPath: "output/generated",
			AllowBack:         true,
			AllowQuit:         true,
			ShowHelp:          true,
			AutoSave:          false, // Don't auto-save when loading from config
			SkipPreview:       false,
		}
	}

	result := &InteractiveGenerationResult{
		ProjectConfig:     loadedConfig.ProjectConfig,
		SelectedTemplates: loadedConfig.SelectedTemplates,
	}

	// Allow modification of loaded configuration if enabled
	if loadedConfig.AllowModification {
		modified, err := ifm.offerConfigurationModification(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("failed to offer configuration modification: %w", err)
		}
		if modified {
			ifm.logger.InfoWithFields("Configuration modified by user", map[string]interface{}{
				"config_name": loadedConfig.Name,
			})
		}
	}

	// Continue with normal flow from output directory selection
	directoryResult, err := ifm.directorySelector.SelectOutputDirectory(ctx, generationConfig.DefaultOutputPath)
	if err != nil {
		return nil, fmt.Errorf("output directory selection failed: %w", err)
	}
	if directoryResult.Cancelled {
		result.Cancelled = true
		return result, nil
	}
	result.OutputDirectory = directoryResult

	// Show preview if not skipped
	if !generationConfig.SkipPreview {
		previewAccepted, err := ifm.generateAndShowPreview(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("preview generation failed: %w", err)
		}
		if !previewAccepted {
			result.Cancelled = true
			return result, nil
		}
		result.PreviewAccepted = previewAccepted
	}

	ifm.logger.InfoWithFields("Generation from loaded configuration completed", map[string]interface{}{
		"config_name":      loadedConfig.Name,
		"project_name":     result.ProjectConfig.Name,
		"output_directory": result.OutputDirectory.Path,
	})

	return result, nil
}

// offerConfigurationModification offers to modify a loaded configuration
func (ifm *InteractiveFlowManager) offerConfigurationModification(ctx context.Context, result *InteractiveGenerationResult) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Modify Configuration",
		Description:  "Would you like to modify the loaded configuration before proceeding?",
		DefaultValue: false,
		AllowBack:    false,
		AllowQuit:    false,
		ShowHelp:     true,
		HelpText: `Configuration Modification:
• You can modify project settings like name, description, author, etc.
• Template selections can also be changed
• Original saved configuration will not be affected
• Choose 'No' to use the configuration as-is`,
	}

	confirmResult, err := ifm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get modification confirmation: %w", err)
	}

	if !confirmResult.Confirmed || confirmResult.Action != "confirm" {
		return false, nil // User chose not to modify
	}

	// Allow modification of project configuration
	modifiedConfig, err := ifm.configCollector.CollectProjectConfiguration(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to collect modified configuration: %w", err)
	}

	if modifiedConfig != nil {
		result.ProjectConfig = modifiedConfig
		return true, nil
	}

	return false, nil
}

// NewPreviewGenerator creates a new preview generator
func NewPreviewGenerator(ui interfaces.InteractiveUIInterface, logger interfaces.Logger) *PreviewGenerator {
	return &PreviewGenerator{
		ui:     ui,
		logger: logger,
	}
}

// GenerateAndShowPreview generates and displays the project structure preview
func (pg *PreviewGenerator) GenerateAndShowPreview(ctx context.Context, result *InteractiveGenerationResult) (bool, error) {
	// Generate preview structure
	preview := pg.generateProjectStructurePreview(result)

	// Show preview in tree format
	treeConfig := interfaces.TreeConfig{
		Title:      "Project Structure Preview",
		Root:       preview,
		Expandable: true,
		ShowIcons:  true,
		MaxDepth:   5,
	}

	if err := pg.ui.ShowTree(ctx, treeConfig); err != nil {
		return false, fmt.Errorf("failed to show preview tree: %w", err)
	}

	// Show summary information
	if err := pg.showPreviewSummary(ctx, result); err != nil {
		return false, fmt.Errorf("failed to show preview summary: %w", err)
	}

	// Confirm with user
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Proceed with Generation",
		Description:  "Review the project structure above. Proceed with generating the project?",
		DefaultValue: true,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Project Structure Preview:
• The tree above shows the complete directory structure that will be created
• Files and directories are organized according to best practices
• You can go back to modify your selections if needed
• Proceeding will create all the files and directories shown`,
	}

	confirmResult, err := pg.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get preview confirmation: %w", err)
	}

	return confirmResult.Confirmed && confirmResult.Action == "confirm", nil
}

// generateProjectStructurePreview generates a tree structure preview
func (pg *PreviewGenerator) generateProjectStructurePreview(result *InteractiveGenerationResult) interfaces.TreeNode {
	projectName := result.ProjectConfig.Name
	if projectName == "" {
		projectName = "my-project"
	}

	root := interfaces.TreeNode{
		Label:      projectName + "/",
		Icon:       "📁",
		Expanded:   true,
		Selectable: false,
		Children:   []interfaces.TreeNode{},
	}

	// Add standard directories based on selected templates
	for _, template := range result.SelectedTemplates {
		if !template.Selected {
			continue
		}

		switch template.Template.Category {
		case "backend":
			root.Children = append(root.Children, pg.generateBackendStructure())
		case "frontend":
			root.Children = append(root.Children, pg.generateFrontendStructure())
		case "mobile":
			root.Children = append(root.Children, pg.generateMobileStructure())
		case "infrastructure":
			root.Children = append(root.Children, pg.generateInfrastructureStructure())
		}
	}

	// Add common directories
	root.Children = append(root.Children, pg.generateCommonStructure()...)

	return root
}

// generateBackendStructure generates backend directory structure
func (pg *PreviewGenerator) generateBackendStructure() interfaces.TreeNode {
	return interfaces.TreeNode{
		Label:    "CommonServer/",
		Icon:     "🐹",
		Expanded: true,
		Children: []interfaces.TreeNode{
			{Label: "cmd/", Icon: "📁", Children: []interfaces.TreeNode{
				{Label: "server/", Icon: "📁", Children: []interfaces.TreeNode{
					{Label: "main.go", Icon: "📄"},
				}},
			}},
			{Label: "internal/", Icon: "📁", Children: []interfaces.TreeNode{
				{Label: "config/", Icon: "📁"},
				{Label: "handlers/", Icon: "📁"},
				{Label: "middleware/", Icon: "📁"},
				{Label: "models/", Icon: "📁"},
				{Label: "services/", Icon: "📁"},
			}},
			{Label: "pkg/", Icon: "📁"},
			{Label: "migrations/", Icon: "📁"},
			{Label: "go.mod", Icon: "📄"},
			{Label: "go.sum", Icon: "📄"},
		},
	}
}

// generateFrontendStructure generates frontend directory structure
func (pg *PreviewGenerator) generateFrontendStructure() interfaces.TreeNode {
	return interfaces.TreeNode{
		Label:    "App/",
		Icon:     "⚛️",
		Expanded: true,
		Children: []interfaces.TreeNode{
			{Label: "main/", Icon: "📁", Children: []interfaces.TreeNode{
				{Label: "src/", Icon: "📁"},
				{Label: "package.json", Icon: "📄"},
			}},
			{Label: "home/", Icon: "📁"},
			{Label: "admin/", Icon: "📁"},
			{Label: "shared-components/", Icon: "📁"},
		},
	}
}

// generateMobileStructure generates mobile directory structure
func (pg *PreviewGenerator) generateMobileStructure() interfaces.TreeNode {
	return interfaces.TreeNode{
		Label:    "Mobile/",
		Icon:     "📱",
		Expanded: true,
		Children: []interfaces.TreeNode{
			{Label: "android/", Icon: "🤖"},
			{Label: "ios/", Icon: "🍎"},
			{Label: "shared/", Icon: "📁"},
		},
	}
}

// generateInfrastructureStructure generates infrastructure directory structure
func (pg *PreviewGenerator) generateInfrastructureStructure() interfaces.TreeNode {
	return interfaces.TreeNode{
		Label:    "Deploy/",
		Icon:     "🚀",
		Expanded: true,
		Children: []interfaces.TreeNode{
			{Label: "docker/", Icon: "🐳"},
			{Label: "k8s/", Icon: "☸️"},
			{Label: "terraform/", Icon: "🏗️"},
			{Label: "monitoring/", Icon: "📊"},
		},
	}
}

// generateCommonStructure generates common project directories
func (pg *PreviewGenerator) generateCommonStructure() []interfaces.TreeNode {
	return []interfaces.TreeNode{
		{
			Label: "Docs/",
			Icon:  "📚",
			Children: []interfaces.TreeNode{
				{Label: "README.md", Icon: "📄"},
				{Label: "CONTRIBUTING.md", Icon: "📄"},
				{Label: "API.md", Icon: "📄"},
			},
		},
		{
			Label: "Scripts/",
			Icon:  "📜",
			Children: []interfaces.TreeNode{
				{Label: "build.sh", Icon: "📄"},
				{Label: "test.sh", Icon: "📄"},
				{Label: "deploy.sh", Icon: "📄"},
			},
		},
		{
			Label: ".github/",
			Icon:  "🐙",
			Children: []interfaces.TreeNode{
				{Label: "workflows/", Icon: "📁"},
				{Label: "ISSUE_TEMPLATE/", Icon: "📁"},
			},
		},
		{Label: ".gitignore", Icon: "📄"},
		{Label: "LICENSE", Icon: "📄"},
		{Label: "Makefile", Icon: "📄"},
	}
}

// showPreviewSummary shows a summary of what will be generated
func (pg *PreviewGenerator) showPreviewSummary(ctx context.Context, result *InteractiveGenerationResult) error {
	// Count files and directories
	fileCount := pg.countFiles(result)
	templateCount := len(result.SelectedTemplates)

	// Create summary table
	headers := []string{"Item", "Value"}
	rows := [][]string{
		{"Project Name", result.ProjectConfig.Name},
		{"Output Directory", result.OutputDirectory.Path},
		{"Templates Selected", fmt.Sprintf("%d", templateCount)},
		{"Estimated Files", fmt.Sprintf("~%d", fileCount)},
		{"License", result.ProjectConfig.License},
		{"Author", result.ProjectConfig.Author},
	}

	if result.ProjectConfig.Description != "" {
		rows = append(rows, []string{"Description", result.ProjectConfig.Description})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Generation Summary",
		Headers: headers,
		Rows:    rows,
	}

	return pg.ui.ShowTable(ctx, tableConfig)
}

// countFiles estimates the number of files that will be generated
func (pg *PreviewGenerator) countFiles(result *InteractiveGenerationResult) int {
	// This is a rough estimate based on typical template file counts
	baseFiles := 10 // README, LICENSE, .gitignore, etc.

	for _, template := range result.SelectedTemplates {
		if !template.Selected {
			continue
		}

		switch template.Template.Category {
		case "backend":
			baseFiles += 25 // Go project files
		case "frontend":
			baseFiles += 30 // Node.js/React project files
		case "mobile":
			baseFiles += 20 // Mobile project files
		case "infrastructure":
			baseFiles += 15 // Docker, K8s, Terraform files
		}
	}

	return baseFiles
}
