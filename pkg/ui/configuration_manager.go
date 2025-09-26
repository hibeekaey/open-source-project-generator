// Package ui provides interactive configuration management for the CLI generator.
//
// This file implements the InteractiveConfigurationManager which handles saving,
// loading, and managing interactive configurations through the UI system.
package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InteractiveConfigurationManager handles interactive configuration operations
type InteractiveConfigurationManager struct {
	ui          interfaces.InteractiveUIInterface
	persistence *config.ConfigurationPersistence
	logger      interfaces.Logger
}

// ConfigurationSaveOptions contains options for saving configurations
type ConfigurationSaveOptions struct {
	Name        string
	Description string
	Tags        []string
	Overwrite   bool
}

// ConfigurationLoadOptions contains options for loading configurations
type ConfigurationLoadOptions struct {
	AllowModification bool
	ShowPreview       bool
	ConfirmLoad       bool
}

// NewInteractiveConfigurationManager creates a new interactive configuration manager
func NewInteractiveConfigurationManager(
	ui interfaces.InteractiveUIInterface,
	configDir string,
	logger interfaces.Logger,
) *InteractiveConfigurationManager {
	persistence := config.NewConfigurationPersistence(configDir, logger)

	return &InteractiveConfigurationManager{
		ui:          ui,
		persistence: persistence,
		logger:      logger,
	}
}

// SaveConfigurationInteractively prompts the user to save a configuration
func (icm *InteractiveConfigurationManager) SaveConfigurationInteractively(
	ctx context.Context,
	projectConfig *models.ProjectConfig,
	selectedTemplates []TemplateSelection,
) (string, error) {
	// Ask if user wants to save configuration
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Save Configuration",
		Description:  "Would you like to save this configuration for future use?",
		DefaultValue: false,
		AllowBack:    false,
		AllowQuit:    false,
		ShowHelp:     true,
		HelpText: `Saving Configuration:
• Saved configurations can be reused for similar projects
• You can load saved configurations to skip the interactive setup
• Configurations include project metadata and template selections
• Output directory is not saved (you'll be prompted each time)
• You can modify saved configurations before using them`,
	}

	confirmResult, err := icm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get save confirmation: %w", err)
	}

	if !confirmResult.Confirmed || confirmResult.Action != "confirm" {
		return "", nil // User chose not to save
	}

	// Get configuration details
	saveOptions, err := icm.collectSaveOptions(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to collect save options: %w", err)
	}

	if saveOptions == nil {
		return "", nil // User cancelled
	}

	// Check if configuration already exists
	if icm.persistence.ConfigurationExists(saveOptions.Name) && !saveOptions.Overwrite {
		overwrite, err := icm.confirmOverwrite(ctx, saveOptions.Name)
		if err != nil {
			return "", fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !overwrite {
			return "", nil // User chose not to overwrite
		}
	}

	// Create saved configuration
	savedConfig := &config.SavedConfiguration{
		Name:              saveOptions.Name,
		Description:       saveOptions.Description,
		Tags:              saveOptions.Tags,
		ProjectConfig:     projectConfig,
		SelectedTemplates: icm.convertTemplateSelections(selectedTemplates),
		GenerationSettings: &config.GenerationSettings{
			IncludeExamples:   true,
			IncludeTests:      true,
			IncludeDocs:       true,
			UpdateVersions:    false,
			MinimalMode:       false,
			BackupExisting:    true,
			OverwriteExisting: false,
		},
		UserPreferences: &config.UserPreferences{
			DefaultLicense:      projectConfig.License,
			DefaultAuthor:       projectConfig.Author,
			DefaultEmail:        projectConfig.Email,
			DefaultOrganization: projectConfig.Organization,
			PreferredFormat:     "yaml",
		},
	}

	// Save configuration
	if err := icm.persistence.SaveConfiguration(saveOptions.Name, savedConfig); err != nil {
		return "", fmt.Errorf("failed to save configuration: %w", err)
	}

	// Show success message
	icm.showSaveSuccess(ctx, saveOptions.Name)

	return saveOptions.Name, nil
}

// LoadConfigurationInteractively prompts the user to load a saved configuration
func (icm *InteractiveConfigurationManager) LoadConfigurationInteractively(
	ctx context.Context,
	options *ConfigurationLoadOptions,
) (*LoadedConfiguration, error) {
	if options == nil {
		options = &ConfigurationLoadOptions{
			AllowModification: true,
			ShowPreview:       true,
			ConfirmLoad:       true,
		}
	}

	// List available configurations
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.showNoConfigurationsMessage(ctx)
		return nil, nil
	}

	// Show configuration selection menu
	selectedConfig, err := icm.selectConfiguration(ctx, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil, nil // User cancelled selection
	}

	// Show configuration preview if requested
	if options.ShowPreview {
		if err := icm.showConfigurationPreview(ctx, selectedConfig); err != nil {
			return nil, fmt.Errorf("failed to show configuration preview: %w", err)
		}
	}

	// Confirm loading if requested
	if options.ConfirmLoad {
		confirmed, err := icm.confirmLoad(ctx, selectedConfig.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm load: %w", err)
		}
		if !confirmed {
			return nil, nil // User chose not to load
		}
	}

	// Convert to loaded configuration
	loadedConfig := &LoadedConfiguration{
		Name:               selectedConfig.Name,
		ProjectConfig:      selectedConfig.ProjectConfig,
		SelectedTemplates:  icm.convertFromSavedTemplates(selectedConfig.SelectedTemplates),
		GenerationSettings: selectedConfig.GenerationSettings,
		UserPreferences:    selectedConfig.UserPreferences,
		LoadedAt:           time.Now(),
		AllowModification:  options.AllowModification,
	}

	icm.logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"config_name":  selectedConfig.Name,
		"project_name": selectedConfig.ProjectConfig.Name,
	})

	return loadedConfig, nil
}

// ManageConfigurationsInteractively provides an interactive configuration management interface
func (icm *InteractiveConfigurationManager) ManageConfigurationsInteractively(ctx context.Context) error {
	for {
		action, err := icm.showManagementMenu(ctx)
		if err != nil {
			return fmt.Errorf("failed to show management menu: %w", err)
		}

		if action == "quit" {
			break
		}

		switch action {
		case "list":
			if err := icm.listConfigurationsInteractively(ctx); err != nil {
				icm.showError(ctx, "Failed to list configurations", err)
			}
		case "view":
			if err := icm.viewConfigurationInteractively(ctx); err != nil {
				icm.showError(ctx, "Failed to view configuration", err)
			}
		case "delete":
			if err := icm.deleteConfigurationInteractively(ctx); err != nil {
				icm.showError(ctx, "Failed to delete configuration", err)
			}
		case "export":
			if err := icm.exportConfigurationInteractively(ctx); err != nil {
				icm.showError(ctx, "Failed to export configuration", err)
			}
		case "import":
			if err := icm.importConfigurationInteractively(ctx); err != nil {
				icm.showError(ctx, "Failed to import configuration", err)
			}
		}
	}

	return nil
}

// LoadedConfiguration represents a loaded configuration with metadata
type LoadedConfiguration struct {
	Name               string                     `json:"name"`
	ProjectConfig      *models.ProjectConfig      `json:"project_config"`
	SelectedTemplates  []TemplateSelection        `json:"selected_templates"`
	GenerationSettings *config.GenerationSettings `json:"generation_settings"`
	UserPreferences    *config.UserPreferences    `json:"user_preferences"`
	LoadedAt           time.Time                  `json:"loaded_at"`
	AllowModification  bool                       `json:"allow_modification"`
}

// Helper methods

// collectSaveOptions collects configuration save options from the user
func (icm *InteractiveConfigurationManager) collectSaveOptions(ctx context.Context) (*ConfigurationSaveOptions, error) {
	options := &ConfigurationSaveOptions{}

	// Get configuration name
	nameConfig := interfaces.TextPromptConfig{
		Prompt:      "Configuration Name",
		Description: "Enter a name for this configuration",
		Required:    true,
		Validator:   icm.validateConfigurationName,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Name Guidelines:
• Use a descriptive name that helps you identify the configuration
• Names should be unique within your saved configurations
• Use alphanumeric characters, hyphens, underscores, and spaces
• Examples: "Go API Project", "nextjs-frontend", "full_stack_app"`,
		MaxLength: 50,
		MinLength: 2,
	}

	nameResult, err := icm.ui.PromptText(ctx, nameConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration name: %w", err)
	}

	if nameResult.Cancelled || nameResult.Action != "submit" {
		return nil, nil // User cancelled
	}

	options.Name = nameResult.Value

	// Get optional description
	descConfig := interfaces.TextPromptConfig{
		Prompt:      "Description (optional)",
		Description: "Enter a description for this configuration",
		Required:    false,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Description:
• Optional brief description of what this configuration is for
• Helps you remember the purpose when selecting configurations later
• Examples: "Standard Go API with PostgreSQL", "React frontend with TypeScript"`,
		MaxLength: 200,
	}

	descResult, err := icm.ui.PromptText(ctx, descConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration description: %w", err)
	}

	if descResult.Cancelled {
		return nil, nil // User cancelled
	}

	options.Description = descResult.Value

	// Get optional tags
	tagsConfig := interfaces.TextPromptConfig{
		Prompt:      "Tags (optional)",
		Description: "Enter comma-separated tags for this configuration",
		Required:    false,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Tags:
• Optional tags to help organize and filter configurations
• Use comma-separated values
• Examples: "backend,api,go", "frontend,react,typescript", "fullstack"`,
		MaxLength: 100,
	}

	tagsResult, err := icm.ui.PromptText(ctx, tagsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration tags: %w", err)
	}

	if tagsResult.Cancelled {
		return nil, nil // User cancelled
	}

	if tagsResult.Value != "" {
		tags := strings.Split(tagsResult.Value, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		options.Tags = tags
	}

	return options, nil
}

// selectConfiguration shows a menu to select from available configurations
func (icm *InteractiveConfigurationManager) selectConfiguration(ctx context.Context, configs []*config.SavedConfiguration) (*config.SavedConfiguration, error) {
	var menuOptions []interfaces.MenuOption

	for _, config := range configs {
		description := config.Description
		if description == "" {
			description = fmt.Sprintf("Created: %s", config.CreatedAt.Format("2006-01-02 15:04"))
		}

		// Add tags to description if available
		if len(config.Tags) > 0 {
			description += fmt.Sprintf(" [Tags: %s]", strings.Join(config.Tags, ", "))
		}

		menuOptions = append(menuOptions, interfaces.MenuOption{
			Label:       config.Name,
			Description: description,
			Value:       config,
		})
	}

	menuConfig := interfaces.MenuConfig{
		Title:       "Select Configuration",
		Description: "Choose a saved configuration to load",
		Options:     menuOptions,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Selection:
• Select a previously saved configuration to load
• Configurations include project settings and template selections
• You can modify loaded configurations before using them
• Use arrow keys to navigate, Enter to select`,
	}

	result, err := icm.ui.ShowMenu(ctx, menuConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to show configuration menu: %w", err)
	}

	if result.Cancelled || result.Action != "select" {
		return nil, nil // User cancelled
	}

	return result.SelectedValue.(*config.SavedConfiguration), nil
}

// showConfigurationPreview displays a preview of the selected configuration
func (icm *InteractiveConfigurationManager) showConfigurationPreview(ctx context.Context, config *config.SavedConfiguration) error {
	// Create preview table
	headers := []string{"Setting", "Value"}
	rows := [][]string{
		{"Configuration Name", config.Name},
		{"Project Name", config.ProjectConfig.Name},
		{"Organization", config.ProjectConfig.Organization},
		{"Author", config.ProjectConfig.Author},
		{"License", config.ProjectConfig.License},
		{"Created", config.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", config.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if config.Description != "" {
		rows = append(rows, []string{"Description", config.Description})
	}

	if len(config.Tags) > 0 {
		rows = append(rows, []string{"Tags", strings.Join(config.Tags, ", ")})
	}

	// Add template information
	templateNames := make([]string, 0, len(config.SelectedTemplates))
	for _, template := range config.SelectedTemplates {
		if template.Selected {
			templateNames = append(templateNames, template.TemplateName)
		}
	}
	if len(templateNames) > 0 {
		rows = append(rows, []string{"Templates", strings.Join(templateNames, ", ")})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Configuration Preview",
		Headers: headers,
		Rows:    rows,
	}

	return icm.ui.ShowTable(ctx, tableConfig)
}

// showManagementMenu displays the configuration management menu
func (icm *InteractiveConfigurationManager) showManagementMenu(ctx context.Context) (string, error) {
	menuOptions := []interfaces.MenuOption{
		{Label: "List Configurations", Description: "View all saved configurations", Value: "list"},
		{Label: "View Configuration", Description: "View details of a specific configuration", Value: "view"},
		{Label: "Delete Configuration", Description: "Delete a saved configuration", Value: "delete"},
		{Label: "Export Configuration", Description: "Export a configuration to file", Value: "export"},
		{Label: "Import Configuration", Description: "Import a configuration from file", Value: "import"},
		{Label: "Quit", Description: "Return to main menu", Value: "quit"},
	}

	menuConfig := interfaces.MenuConfig{
		Title:       "Configuration Management",
		Description: "Manage your saved configurations",
		Options:     menuOptions,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Management:
• List: View all your saved configurations
• View: See detailed information about a configuration
• Delete: Remove a configuration you no longer need
• Export: Save a configuration to a file for sharing or backup
• Import: Load a configuration from a file`,
	}

	result, err := icm.ui.ShowMenu(ctx, menuConfig)
	if err != nil {
		return "", fmt.Errorf("failed to show management menu: %w", err)
	}

	if result.Cancelled {
		return "quit", nil
	}

	return result.SelectedValue.(string), nil
}

// listConfigurationsInteractively shows a list of all configurations
func (icm *InteractiveConfigurationManager) listConfigurationsInteractively(ctx context.Context) error {
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.showNoConfigurationsMessage(ctx)
		return nil
	}

	// Create table with configuration information
	headers := []string{"Name", "Project", "Templates", "Created", "Updated"}
	rows := make([][]string, 0, len(configs))

	for _, config := range configs {
		templateCount := 0
		for _, template := range config.SelectedTemplates {
			if template.Selected {
				templateCount++
			}
		}

		rows = append(rows, []string{
			config.Name,
			config.ProjectConfig.Name,
			fmt.Sprintf("%d", templateCount),
			config.CreatedAt.Format("2006-01-02"),
			config.UpdatedAt.Format("2006-01-02"),
		})
	}

	tableConfig := interfaces.TableConfig{
		Title:      "Saved Configurations",
		Headers:    headers,
		Rows:       rows,
		Pagination: true,
		PageSize:   10,
		Sortable:   true,
	}

	return icm.ui.ShowTable(ctx, tableConfig)
}

// Additional helper methods would be implemented here...

// validateConfigurationName validates a configuration name
func (icm *InteractiveConfigurationManager) validateConfigurationName(name string) error {
	if name == "" {
		return interfaces.NewValidationError("name", name, "Configuration name is required", "required").
			WithSuggestions("Enter a descriptive name for your configuration")
	}

	if len(name) < 2 {
		return interfaces.NewValidationError("name", name, "Configuration name must be at least 2 characters long", "min_length").
			WithSuggestions("Use a longer, more descriptive name")
	}

	if len(name) > 50 {
		return interfaces.NewValidationError("name", name, "Configuration name must be at most 50 characters long", "max_length").
			WithSuggestions("Use a shorter, more concise name")
	}

	return nil
}

// convertTemplateSelections converts UI template selections to saved format
func (icm *InteractiveConfigurationManager) convertTemplateSelections(selections []TemplateSelection) []config.TemplateSelection {
	result := make([]config.TemplateSelection, len(selections))
	for i, sel := range selections {
		result[i] = config.TemplateSelection{
			TemplateName: sel.Template.Name,
			Category:     sel.Template.Category,
			Technology:   sel.Template.Technology,
			Version:      sel.Template.Version,
			Selected:     sel.Selected,
			Options:      sel.Options,
		}
	}
	return result
}

// convertFromSavedTemplates converts saved template selections to UI format
func (icm *InteractiveConfigurationManager) convertFromSavedTemplates(saved []config.TemplateSelection) []TemplateSelection {
	result := make([]TemplateSelection, len(saved))
	for i, sel := range saved {
		result[i] = TemplateSelection{
			Template: interfaces.TemplateInfo{
				Name:       sel.TemplateName,
				Category:   sel.Category,
				Technology: sel.Technology,
				Version:    sel.Version,
			},
			Selected: sel.Selected,
			Options:  sel.Options,
		}
	}
	return result
}

// Placeholder implementations for remaining methods
func (icm *InteractiveConfigurationManager) confirmOverwrite(ctx context.Context, name string) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Overwrite Configuration '%s'", name),
		Description:  "A configuration with this name already exists. Do you want to overwrite it?",
		DefaultValue: false,
		AllowBack:    true,
		AllowQuit:    true,
	}

	result, err := icm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

func (icm *InteractiveConfigurationManager) confirmLoad(ctx context.Context, name string) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Load Configuration '%s'", name),
		Description:  "Load this configuration and use it for project generation?",
		DefaultValue: true,
		AllowBack:    true,
		AllowQuit:    true,
	}

	result, err := icm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

func (icm *InteractiveConfigurationManager) showSaveSuccess(ctx context.Context, name string) {
	// This would show a success message - placeholder implementation
	fmt.Printf("Configuration '%s' saved successfully!\n", name)
}

func (icm *InteractiveConfigurationManager) showNoConfigurationsMessage(ctx context.Context) {
	// This would show a message about no configurations - placeholder implementation
	fmt.Println("No saved configurations found.")
}

func (icm *InteractiveConfigurationManager) showError(ctx context.Context, title string, err error) {
	// This would show an error message - placeholder implementation
	fmt.Printf("Error: %s - %v\n", title, err)
}

// viewConfigurationInteractively shows detailed view of a configuration
func (icm *InteractiveConfigurationManager) viewConfigurationInteractively(ctx context.Context) error {
	// Get list of configurations
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.showNoConfigurationsMessage(ctx)
		return nil
	}

	// Select configuration to view
	selectedConfig, err := icm.selectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	// Show detailed view
	return icm.showDetailedConfigurationView(ctx, selectedConfig)
}

func (icm *InteractiveConfigurationManager) deleteConfigurationInteractively(ctx context.Context) error {
	// Get list of configurations
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.showNoConfigurationsMessage(ctx)
		return nil
	}

	// Select configuration to delete
	selectedConfig, err := icm.selectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	// Confirm deletion
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Delete Configuration '%s'", selectedConfig.Name),
		Description:  "Are you sure you want to delete this configuration? This action cannot be undone.",
		DefaultValue: false,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Delete Configuration:
• This will permanently remove the configuration from your saved configurations
• The action cannot be undone
• Consider exporting the configuration first if you might need it later`,
	}

	confirmResult, err := icm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return fmt.Errorf("failed to get deletion confirmation: %w", err)
	}

	if !confirmResult.Confirmed || confirmResult.Action != "confirm" {
		return nil // User chose not to delete
	}

	// Delete configuration
	if err := icm.persistence.DeleteConfiguration(selectedConfig.Name); err != nil {
		return fmt.Errorf("failed to delete configuration: %w", err)
	}

	// Show success message
	fmt.Printf("Configuration '%s' deleted successfully.\n", selectedConfig.Name)
	return nil
}

func (icm *InteractiveConfigurationManager) exportConfigurationInteractively(ctx context.Context) error {
	// Get list of configurations
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.showNoConfigurationsMessage(ctx)
		return nil
	}

	// Select configuration to export
	selectedConfig, err := icm.selectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	// Select export format
	formatOptions := []interfaces.MenuOption{
		{Label: "YAML", Description: "Human-readable YAML format", Value: "yaml"},
		{Label: "JSON", Description: "Machine-readable JSON format", Value: "json"},
	}

	formatConfig := interfaces.MenuConfig{
		Title:       "Export Format",
		Description: "Choose the format for exporting the configuration",
		Options:     formatOptions,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Export Formats:
• YAML: Human-readable format, good for editing and version control
• JSON: Machine-readable format, good for automation and APIs`,
	}

	formatResult, err := icm.ui.ShowMenu(ctx, formatConfig)
	if err != nil {
		return fmt.Errorf("failed to select export format: %w", err)
	}

	if formatResult.Cancelled || formatResult.Action != "select" {
		return nil // User cancelled
	}

	format := formatResult.SelectedValue.(string)

	// Get export file path
	defaultFilename := fmt.Sprintf("%s.%s", selectedConfig.Name, format)
	pathConfig := interfaces.TextPromptConfig{
		Prompt:       "Export File Path",
		Description:  "Enter the path where you want to save the exported configuration",
		DefaultValue: defaultFilename,
		Required:     true,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Export File Path:
• Specify the full path where you want to save the configuration
• Include the file extension (.yaml or .json)
• The directory will be created if it doesn't exist`,
		MaxLength: 200,
	}

	pathResult, err := icm.ui.PromptText(ctx, pathConfig)
	if err != nil {
		return fmt.Errorf("failed to get export path: %w", err)
	}

	if pathResult.Cancelled || pathResult.Action != "submit" {
		return nil // User cancelled
	}

	exportPath := pathResult.Value

	// Export configuration
	data, err := icm.persistence.ExportConfiguration(selectedConfig.Name, format)
	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}

	// Write to file
	if err := icm.writeExportFile(exportPath, data); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	// Show success message
	fmt.Printf("Configuration '%s' exported to '%s' successfully.\n", selectedConfig.Name, exportPath)
	return nil
}

func (icm *InteractiveConfigurationManager) importConfigurationInteractively(ctx context.Context) error {
	// Get import file path
	pathConfig := interfaces.TextPromptConfig{
		Prompt:      "Import File Path",
		Description: "Enter the path to the configuration file you want to import",
		Required:    true,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Import File Path:
• Specify the full path to the configuration file
• Supported formats: .yaml, .yml, .json
• The file should be a previously exported configuration`,
		MaxLength: 200,
	}

	pathResult, err := icm.ui.PromptText(ctx, pathConfig)
	if err != nil {
		return fmt.Errorf("failed to get import path: %w", err)
	}

	if pathResult.Cancelled || pathResult.Action != "submit" {
		return nil // User cancelled
	}

	importPath := pathResult.Value

	// Determine format from file extension
	format := "yaml"
	if strings.HasSuffix(strings.ToLower(importPath), ".json") {
		format = "json"
	}

	// Read import file
	data, err := icm.readImportFile(importPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Get configuration name for import
	nameConfig := interfaces.TextPromptConfig{
		Prompt:      "Configuration Name",
		Description: "Enter a name for the imported configuration",
		Required:    true,
		Validator:   icm.validateConfigurationName,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Configuration Name:
• Choose a unique name for the imported configuration
• This name will be used to identify the configuration in your saved list
• Use descriptive names to help you remember what the configuration is for`,
		MaxLength: 50,
		MinLength: 2,
	}

	nameResult, err := icm.ui.PromptText(ctx, nameConfig)
	if err != nil {
		return fmt.Errorf("failed to get configuration name: %w", err)
	}

	if nameResult.Cancelled || nameResult.Action != "submit" {
		return nil // User cancelled
	}

	configName := nameResult.Value

	// Check if configuration already exists
	if icm.persistence.ConfigurationExists(configName) {
		overwrite, err := icm.confirmOverwrite(ctx, configName)
		if err != nil {
			return fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !overwrite {
			return nil // User chose not to overwrite
		}
	}

	// Import configuration
	if err := icm.persistence.ImportConfiguration(configName, data, format); err != nil {
		return fmt.Errorf("failed to import configuration: %w", err)
	}

	// Show success message
	fmt.Printf("Configuration imported as '%s' successfully.\n", configName)
	return nil
}

// showDetailedConfigurationView shows a detailed view of a configuration
func (icm *InteractiveConfigurationManager) showDetailedConfigurationView(ctx context.Context, config *config.SavedConfiguration) error {
	// Show basic information
	if err := icm.showConfigurationPreview(ctx, config); err != nil {
		return fmt.Errorf("failed to show configuration preview: %w", err)
	}

	// Show template details
	if err := icm.showTemplateDetails(ctx, config); err != nil {
		return fmt.Errorf("failed to show template details: %w", err)
	}

	// Show generation settings if available
	if config.GenerationSettings != nil {
		if err := icm.showGenerationSettings(ctx, config.GenerationSettings); err != nil {
			return fmt.Errorf("failed to show generation settings: %w", err)
		}
	}

	return nil
}

// showTemplateDetails shows detailed template information
func (icm *InteractiveConfigurationManager) showTemplateDetails(ctx context.Context, config *config.SavedConfiguration) error {
	if len(config.SelectedTemplates) == 0 {
		return nil
	}

	headers := []string{"Template", "Category", "Technology", "Version", "Selected"}
	rows := make([][]string, 0, len(config.SelectedTemplates))

	for _, template := range config.SelectedTemplates {
		selected := "No"
		if template.Selected {
			selected = "Yes"
		}

		rows = append(rows, []string{
			template.TemplateName,
			template.Category,
			template.Technology,
			template.Version,
			selected,
		})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Template Details",
		Headers: headers,
		Rows:    rows,
	}

	return icm.ui.ShowTable(ctx, tableConfig)
}

// showGenerationSettings shows generation settings information
func (icm *InteractiveConfigurationManager) showGenerationSettings(ctx context.Context, settings *config.GenerationSettings) error {
	headers := []string{"Setting", "Value"}
	rows := [][]string{
		{"Include Examples", fmt.Sprintf("%t", settings.IncludeExamples)},
		{"Include Tests", fmt.Sprintf("%t", settings.IncludeTests)},
		{"Include Documentation", fmt.Sprintf("%t", settings.IncludeDocs)},
		{"Update Versions", fmt.Sprintf("%t", settings.UpdateVersions)},
		{"Minimal Mode", fmt.Sprintf("%t", settings.MinimalMode)},
		{"Backup Existing", fmt.Sprintf("%t", settings.BackupExisting)},
		{"Overwrite Existing", fmt.Sprintf("%t", settings.OverwriteExisting)},
	}

	if len(settings.ExcludePatterns) > 0 {
		rows = append(rows, []string{"Exclude Patterns", strings.Join(settings.ExcludePatterns, ", ")})
	}

	if len(settings.IncludeOnlyPaths) > 0 {
		rows = append(rows, []string{"Include Only Paths", strings.Join(settings.IncludeOnlyPaths, ", ")})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Generation Settings",
		Headers: headers,
		Rows:    rows,
	}

	return icm.ui.ShowTable(ctx, tableConfig)
}

// Helper methods for file operations
func (icm *InteractiveConfigurationManager) writeExportFile(path string, data []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (icm *InteractiveConfigurationManager) readImportFile(path string) ([]byte, error) {
	// Validate and clean path to prevent directory traversal
	path = filepath.Clean(path)

	// Ensure path is absolute to prevent traversal attacks
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("import file path must be absolute: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}
