package collectors

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// InputCollector handles collecting configuration input from users
type InputCollector struct {
	ui interfaces.InteractiveUIInterface
}

// NewInputCollector creates a new input collector
func NewInputCollector(ui interfaces.InteractiveUIInterface) *InputCollector {
	return &InputCollector{
		ui: ui,
	}
}

// ConfigurationSaveOptions contains options for saving configurations
type ConfigurationSaveOptions struct {
	Name        string
	Description string
	Tags        []string
	Overwrite   bool
}

// CollectSaveOptions collects configuration save options from the user
func (ic *InputCollector) CollectSaveOptions(ctx context.Context, validator func(string) error) (*ConfigurationSaveOptions, error) {
	options := &ConfigurationSaveOptions{}

	// Get configuration name
	nameConfig := interfaces.TextPromptConfig{
		Prompt:      "Configuration Name",
		Description: "Enter a name for this configuration",
		Required:    true,
		Validator:   validator,
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

	nameResult, err := ic.ui.PromptText(ctx, nameConfig)
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

	descResult, err := ic.ui.PromptText(ctx, descConfig)
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

	tagsResult, err := ic.ui.PromptText(ctx, tagsConfig)
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

// SelectConfiguration shows a menu to select from available configurations
func (ic *InputCollector) SelectConfiguration(ctx context.Context, configs []*config.SavedConfiguration) (*config.SavedConfiguration, error) {
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

	result, err := ic.ui.ShowMenu(ctx, menuConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to show configuration menu: %w", err)
	}

	if result.Cancelled || result.Action != "select" {
		return nil, nil // User cancelled
	}

	return result.SelectedValue.(*config.SavedConfiguration), nil
}

// CollectExportFormat prompts user to select export format
func (ic *InputCollector) CollectExportFormat(ctx context.Context) (string, error) {
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

	formatResult, err := ic.ui.ShowMenu(ctx, formatConfig)
	if err != nil {
		return "", fmt.Errorf("failed to select export format: %w", err)
	}

	if formatResult.Cancelled || formatResult.Action != "select" {
		return "", nil // User cancelled
	}

	return formatResult.SelectedValue.(string), nil
}

// CollectExportPath prompts user for export file path
func (ic *InputCollector) CollectExportPath(ctx context.Context, defaultFilename string) (string, error) {
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

	pathResult, err := ic.ui.PromptText(ctx, pathConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get export path: %w", err)
	}

	if pathResult.Cancelled || pathResult.Action != "submit" {
		return "", nil // User cancelled
	}

	return pathResult.Value, nil
}

// CollectImportPath prompts user for import file path
func (ic *InputCollector) CollectImportPath(ctx context.Context) (string, error) {
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

	pathResult, err := ic.ui.PromptText(ctx, pathConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get import path: %w", err)
	}

	if pathResult.Cancelled || pathResult.Action != "submit" {
		return "", nil // User cancelled
	}

	return pathResult.Value, nil
}

// CollectConfigurationName prompts user for configuration name during import
func (ic *InputCollector) CollectConfigurationName(ctx context.Context, validator func(string) error) (string, error) {
	nameConfig := interfaces.TextPromptConfig{
		Prompt:      "Configuration Name",
		Description: "Enter a name for the imported configuration",
		Required:    true,
		Validator:   validator,
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

	nameResult, err := ic.ui.PromptText(ctx, nameConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get configuration name: %w", err)
	}

	if nameResult.Cancelled || nameResult.Action != "submit" {
		return "", nil // User cancelled
	}

	return nameResult.Value, nil
}

// ShowManagementMenu displays the configuration management menu
func (ic *InputCollector) ShowManagementMenu(ctx context.Context) (string, error) {
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

	result, err := ic.ui.ShowMenu(ctx, menuConfig)
	if err != nil {
		return "", fmt.Errorf("failed to show management menu: %w", err)
	}

	if result.Cancelled {
		return "quit", nil
	}

	return result.SelectedValue.(string), nil
}
