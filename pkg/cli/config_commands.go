// Package cli provides configuration management commands for the CLI generator.
//
// This file implements CLI commands for managing saved configurations including
// listing, viewing, saving, loading, and deleting configurations.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
	"github.com/spf13/cobra"
)

// setupConfigCommand sets up the config command with all subcommands
func (c *CLI) setupConfigCommand() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage saved project configurations",
		Long: `Manage saved project configurations for reuse in future project generation.

Save, load, view, and manage project configurations to skip interactive setup.`,
		Example: `  # List configurations
  generator config list

  # View a configuration
  generator config view my-config

  # Delete a configuration
  generator config delete old-config

  # Export configuration
  generator config export my-config`,
	}

	// Add subcommands
	c.setupConfigListCommand(configCmd)
	c.setupConfigViewCommand(configCmd)
	c.setupConfigDeleteCommand(configCmd)
	c.setupConfigExportCommand(configCmd)
	c.setupConfigImportCommand(configCmd)
	c.setupConfigManageCommand(configCmd)

	c.rootCmd.AddCommand(configCmd)
}

// setupConfigListCommand sets up the config list subcommand
func (c *CLI) setupConfigListCommand(parent *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List all saved configurations",
		Long: `List all saved project configurations with optional filtering and sorting.
Shows basic information about each configuration including name, project name,
creation date, and associated templates.`,
		RunE: c.runConfigList,
		Example: `  # List all configurations
  generator config list

  # List configurations with JSON output
  generator config list --output-format json

  # List configurations sorted by creation date
  generator config list --sort-by created_at --sort-order desc

  # Search configurations by name or description
  generator config list --search "api"

  # Filter configurations by tags
  generator config list --tags backend,go

  # Limit number of results
  generator config list --limit 10`,
	}

	listCmd.Flags().String("sort-by", "updated_at", "Sort by field (name, created_at, updated_at)")
	listCmd.Flags().String("sort-order", "desc", "Sort order (asc, desc)")
	listCmd.Flags().String("search", "", "Search in configuration names and descriptions")
	listCmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	listCmd.Flags().Int("limit", 50, "Maximum number of configurations to show")
	listCmd.Flags().Int("offset", 0, "Number of configurations to skip (pagination)")

	parent.AddCommand(listCmd)
}

// setupConfigViewCommand sets up the config view subcommand
func (c *CLI) setupConfigViewCommand(parent *cobra.Command) {
	viewCmd := &cobra.Command{
		Use:   "view <config-name> [flags]",
		Short: "View detailed information about a configuration",
		Long: `Display detailed information about a specific saved configuration including
project metadata, template selections, generation settings, and file information.`,
		Args: cobra.ExactArgs(1),
		RunE: c.runConfigView,
		Example: `  # View configuration details
  generator config view my-api-config

  # View configuration with JSON output
  generator config view my-api-config --output-format json

  # View configuration with full template details
  generator config view my-api-config --show-templates`,
	}

	viewCmd.Flags().Bool("show-templates", true, "Show detailed template information")
	viewCmd.Flags().Bool("show-settings", true, "Show generation settings")

	parent.AddCommand(viewCmd)
}

// setupConfigDeleteCommand sets up the config delete subcommand
func (c *CLI) setupConfigDeleteCommand(parent *cobra.Command) {
	deleteCmd := &cobra.Command{
		Use:   "delete <config-name> [flags]",
		Short: "Delete a saved configuration",
		Long: `Delete a saved project configuration. This action cannot be undone.
Consider exporting the configuration first if you might need it later.`,
		Args: cobra.ExactArgs(1),
		RunE: c.runConfigDelete,
		Example: `  # Delete a configuration (with confirmation)
  generator config delete old-config

  # Delete without confirmation (non-interactive)
  generator config delete old-config --force

  # Delete multiple configurations
  generator config delete config1 config2 config3`,
	}

	deleteCmd.Flags().Bool("force", false, "Delete without confirmation")

	parent.AddCommand(deleteCmd)
}

// setupConfigExportCommand sets up the config export subcommand
func (c *CLI) setupConfigExportCommand(parent *cobra.Command) {
	exportCmd := &cobra.Command{
		Use:   "export [config-name-or-file] [flags]",
		Short: "Export a configuration to a file",
		Long: `Export a saved configuration or current configuration to a file for sharing, backup, or version control.
Supports YAML and JSON formats.`,
		Args: cobra.RangeArgs(0, 1),
		RunE: c.runConfigExport,
		Example: `  # Export configuration to YAML file
  generator config export my-config --output my-config.yaml

  # Export configuration to JSON file
  generator config export my-config --format json --output my-config.json

  # Export to stdout
  generator config export my-config --format yaml`,
	}

	exportCmd.Flags().String("format", "yaml", "Export format (yaml, json)")
	exportCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")

	parent.AddCommand(exportCmd)
}

// setupConfigImportCommand sets up the config import subcommand
func (c *CLI) setupConfigImportCommand(parent *cobra.Command) {
	importCmd := &cobra.Command{
		Use:   "import [flags]",
		Short: "Import a configuration from a file",
		Long: `Import a project configuration from a file. The file should be a previously
exported configuration in YAML or JSON format.`,
		RunE: c.runConfigImport,
		Example: `  # Import configuration from file
  generator config import --file my-config.yaml --name imported-config

  # Import with automatic name detection
  generator config import --file my-config.yaml

  # Import and overwrite existing configuration
  generator config import --file my-config.yaml --name existing-config --force`,
	}

	importCmd.Flags().StringP("file", "f", "", "Configuration file to import (required)")
	importCmd.Flags().String("name", "", "Name for imported configuration (auto-detected if not provided)")
	importCmd.Flags().Bool("force", false, "Overwrite existing configuration without confirmation")

	_ = importCmd.MarkFlagRequired("file")

	parent.AddCommand(importCmd)
}

// setupConfigManageCommand sets up the config manage subcommand
func (c *CLI) setupConfigManageCommand(parent *cobra.Command) {
	manageCmd := &cobra.Command{
		Use:   "manage",
		Short: "Interactive configuration management",
		Long: `Launch an interactive interface for managing saved configurations.
Provides a menu-driven interface for listing, viewing, deleting, exporting,
and importing configurations.`,
		RunE: c.runConfigManage,
		Example: `  # Launch interactive configuration management
  generator config manage`,
	}

	parent.AddCommand(manageCmd)
}

// Command implementations

// runConfigList executes the config list command
func (c *CLI) runConfigList(cmd *cobra.Command, args []string) error {
	// Get flags
	sortBy, _ := cmd.Flags().GetString("sort-by")
	sortOrder, _ := cmd.Flags().GetString("sort-order")
	search, _ := cmd.Flags().GetString("search")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Create configuration persistence
	persistence := c.createConfigurationPersistence()

	// Set up list options
	options := &config.ConfigurationListOptions{
		SortBy:      sortBy,
		SortOrder:   sortOrder,
		SearchQuery: search,
		FilterTags:  tags,
		Limit:       limit,
		Offset:      offset,
	}

	// Get configurations
	configs, err := persistence.ListConfigurations(options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Unable to access your saved configurations."),
			c.info("Check if the configuration directory exists and is readable"))
	}

	// Output results
	return c.outputConfigurationList(configs, outputFormat)
}

// runConfigView executes the config view command
func (c *CLI) runConfigView(cmd *cobra.Command, args []string) error {
	configName := args[0]
	showTemplates, _ := cmd.Flags().GetBool("show-templates")
	showSettings, _ := cmd.Flags().GetBool("show-settings")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Create configuration persistence
	persistence := c.createConfigurationPersistence()

	// Load configuration
	config, err := persistence.LoadConfiguration(configName)
	if err != nil {
		return fmt.Errorf("🚫 %s %s %s",
			c.error("Unable to load configuration"),
			c.highlight(fmt.Sprintf("'%s'.", configName)),
			c.info("Check if it exists and is readable"))
	}

	// Output configuration details
	return c.outputConfigurationDetails(config, showTemplates, showSettings, outputFormat)
}

// runConfigDelete executes the config delete command
func (c *CLI) runConfigDelete(cmd *cobra.Command, args []string) error {
	configName := args[0]
	force, _ := cmd.Flags().GetBool("force")

	// Create configuration persistence
	persistence := c.createConfigurationPersistence()

	// Check if configuration exists
	if !persistence.ConfigurationExists(configName) {
		return fmt.Errorf("🚫 %s %s %s",
			c.error("Configuration"),
			c.highlight(fmt.Sprintf("'%s'", configName)),
			c.info("doesn't exist. Use 'generator config list' to see available configurations"))
	}

	// Confirm deletion if not forced
	if !force && !c.isNonInteractiveMode() {
		fmt.Printf("🗑️  Are you sure you want to delete configuration '%s'? (y/N): ", configName)
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			c.QuietOutput("%s %s",
				c.warning("❌ Deletion cancelled."),
				c.info("Your configuration is safe"))
			return nil
		}
	}

	// Delete configuration
	if err := persistence.DeleteConfiguration(configName); err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Unable to delete the configuration."),
			c.info("Check file permissions and try again"))
	}

	c.SuccessOutput("🗑️  Configuration '%s' deleted successfully", configName)
	return nil
}

// runConfigExport executes the config export command
func (c *CLI) runConfigExport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	var data []byte
	var err error

	if len(args) == 0 {
		// Export current configuration to stdout or specified output file
		if output == "" {
			output = "generator-config.yaml"
		}

		// Export current configuration
		err = c.ExportConfig(output)
		if err != nil {
			return fmt.Errorf("🚫 Couldn't export the current configuration: %w", err)
		}

		c.SuccessOutput("📤 Configuration exported to '%s'", output)
		return nil
	}

	// Check if the argument is a file path (contains path separators or file extension)
	arg := args[0]
	if strings.Contains(arg, string(filepath.Separator)) || strings.Contains(arg, ".") {
		// Treat as output file path for current configuration
		err = c.ExportConfig(arg)
		if err != nil {
			return fmt.Errorf("failed to export current configuration: %w", err)
		}

		c.SuccessOutput("📤 Configuration exported to '%s'", arg)
		return nil
	}

	// Treat as saved configuration name
	configName := arg

	// Create configuration persistence
	persistence := c.createConfigurationPersistence()

	// Export saved configuration
	data, err = persistence.ExportConfiguration(configName, format)
	if err != nil {
		return fmt.Errorf("🚫 Couldn't export the configuration: %w", err)
	}

	// Output to file or stdout
	if output != "" {
		// Create directory if needed
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Write to file
		if err := os.WriteFile(output, data, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		c.SuccessOutput("📤 Configuration '%s' exported to '%s'", configName, output)
	} else {
		// Output to stdout
		fmt.Print(string(data))
	}

	return nil
}

// runConfigImport executes the config import command
func (c *CLI) runConfigImport(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")
	name, _ := cmd.Flags().GetString("name")
	force, _ := cmd.Flags().GetBool("force")

	// Read import file
	// Validate and clean path to prevent directory traversal
	file = filepath.Clean(file)

	// Ensure path is absolute to prevent traversal attacks
	if !filepath.IsAbs(file) {
		return fmt.Errorf("import file path must be absolute: %s", file)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Determine format from file extension
	format := "yaml"
	if strings.HasSuffix(strings.ToLower(file), ".json") {
		format = "json"
	}

	// Auto-detect name if not provided
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	}

	// Create configuration persistence
	persistence := c.createConfigurationPersistence()

	// Check if configuration exists
	if persistence.ConfigurationExists(name) && !force {
		if c.isNonInteractiveMode() {
			return fmt.Errorf("configuration '%s' already exists (use --force to overwrite)", name)
		}

		fmt.Printf("⚠️  Configuration '%s' already exists. Overwrite? (y/N): ", name)
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			c.QuietOutput("❌ Import cancelled")
			return nil
		}
	}

	// Import configuration
	if err := persistence.ImportConfiguration(name, data, format); err != nil {
		return fmt.Errorf("🚫 Couldn't import the configuration: %w", err)
	}

	c.SuccessOutput("📥 Configuration imported as '%s'", name)
	return nil
}

// runConfigManage executes the config manage command
func (c *CLI) runConfigManage(cmd *cobra.Command, args []string) error {
	if c.isNonInteractiveMode() {
		return fmt.Errorf("interactive configuration management not available in non-interactive mode")
	}

	ctx := context.Background()

	// Create interactive configuration manager
	configManager := ui.NewInteractiveConfigurationManager(
		c.interactiveUI,
		c.getConfigDirectory(),
		c.logger,
	)

	// Run interactive management
	return configManager.ManageConfigurationsInteractively(ctx)
}

// Helper methods

// createConfigurationPersistence creates a configuration persistence instance
func (c *CLI) createConfigurationPersistence() *config.ConfigurationPersistence {
	return config.NewConfigurationPersistence(c.getConfigDirectory(), c.logger)
}

// getConfigDirectory returns the configuration directory path
func (c *CLI) getConfigDirectory() string {
	if c.configManager != nil {
		return c.configManager.GetConfigLocation()
	}

	// Fallback to default location
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".generator", "configs")
}

// outputConfigurationList outputs a list of configurations
func (c *CLI) outputConfigurationList(configs []*config.SavedConfiguration, format string) error {
	if len(configs) == 0 {
		c.QuietOutput("📋 No saved configurations found")
		return nil
	}

	switch format {
	case "json":
		return c.outputConfigurationListJSON(configs)
	case "yaml":
		return c.outputConfigurationListYAML(configs)
	default:
		return c.outputConfigurationListTable(configs)
	}
}

// outputConfigurationListTable outputs configurations in table format
func (c *CLI) outputConfigurationListTable(configs []*config.SavedConfiguration) error {
	c.QuietOutput("📋 Saved Configurations:")
	c.QuietOutput("========================")
	c.QuietOutput("")

	for _, config := range configs {
		c.QuietOutput("Name: %s", config.Name)
		c.QuietOutput("Project: %s", config.ProjectConfig.Name)
		if config.Description != "" {
			c.QuietOutput("Description: %s", config.Description)
		}
		c.QuietOutput("Created: %s", config.CreatedAt.Format("2006-01-02 15:04:05"))
		c.QuietOutput("Updated: %s", config.UpdatedAt.Format("2006-01-02 15:04:05"))

		// Count selected templates
		templateCount := 0
		for _, template := range config.SelectedTemplates {
			if template.Selected {
				templateCount++
			}
		}
		c.QuietOutput("Templates: %d", templateCount)

		if len(config.Tags) > 0 {
			c.QuietOutput("Tags: %s", strings.Join(config.Tags, ", "))
		}

		c.QuietOutput("")
	}

	return nil
}

// outputConfigurationListJSON outputs configurations in JSON format
func (c *CLI) outputConfigurationListJSON(configs []*config.SavedConfiguration) error {
	// Create simplified output for JSON
	type ConfigSummary struct {
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		ProjectName string    `json:"project_name"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Tags        []string  `json:"tags,omitempty"`
		Templates   int       `json:"template_count"`
	}

	summaries := make([]ConfigSummary, len(configs))
	for i, config := range configs {
		templateCount := 0
		for _, template := range config.SelectedTemplates {
			if template.Selected {
				templateCount++
			}
		}

		summaries[i] = ConfigSummary{
			Name:        config.Name,
			Description: config.Description,
			ProjectName: config.ProjectConfig.Name,
			CreatedAt:   config.CreatedAt,
			UpdatedAt:   config.UpdatedAt,
			Tags:        config.Tags,
			Templates:   templateCount,
		}
	}

	return c.outputJSON(summaries)
}

// outputConfigurationListYAML outputs configurations in YAML format
func (c *CLI) outputConfigurationListYAML(configs []*config.SavedConfiguration) error {
	// Similar to JSON but output as YAML
	return fmt.Errorf("YAML output not implemented yet")
}

// outputConfigurationDetails outputs detailed configuration information
func (c *CLI) outputConfigurationDetails(config *config.SavedConfiguration, showTemplates, showSettings bool, format string) error {
	switch format {
	case "json":
		return c.outputJSON(config)
	case "yaml":
		return c.outputYAML(config)
	default:
		return c.outputConfigurationDetailsTable(config, showTemplates, showSettings)
	}
}

// outputConfigurationDetailsTable outputs configuration details in table format
func (c *CLI) outputConfigurationDetailsTable(config *config.SavedConfiguration, showTemplates, showSettings bool) error {
	c.QuietOutput("📋 Configuration Details: %s", config.Name)
	c.QuietOutput("===============================")
	c.QuietOutput("")

	// Basic information
	c.QuietOutput("Name: %s", config.Name)
	if config.Description != "" {
		c.QuietOutput("Description: %s", config.Description)
	}
	c.QuietOutput("Version: %s", config.Version)
	c.QuietOutput("Created: %s", config.CreatedAt.Format("2006-01-02 15:04:05"))
	c.QuietOutput("Updated: %s", config.UpdatedAt.Format("2006-01-02 15:04:05"))

	if len(config.Tags) > 0 {
		c.QuietOutput("Tags: %s", strings.Join(config.Tags, ", "))
	}

	c.QuietOutput("")

	// Project configuration
	c.QuietOutput("Project Configuration:")
	c.QuietOutput("  Name: %s", config.ProjectConfig.Name)
	c.QuietOutput("  Organization: %s", config.ProjectConfig.Organization)
	c.QuietOutput("  Author: %s", config.ProjectConfig.Author)
	c.QuietOutput("  Email: %s", config.ProjectConfig.Email)
	c.QuietOutput("  License: %s", config.ProjectConfig.License)
	if config.ProjectConfig.Description != "" {
		c.QuietOutput("  Description: %s", config.ProjectConfig.Description)
	}
	if config.ProjectConfig.Repository != "" {
		c.QuietOutput("  Repository: %s", config.ProjectConfig.Repository)
	}

	c.QuietOutput("")

	// Template information
	if showTemplates && len(config.SelectedTemplates) > 0 {
		c.QuietOutput("Selected Templates:")
		for _, template := range config.SelectedTemplates {
			if template.Selected {
				c.QuietOutput("  - %s (%s/%s) v%s", template.TemplateName, template.Category, template.Technology, template.Version)
			}
		}
		c.QuietOutput("")
	}

	// Generation settings
	if showSettings && config.GenerationSettings != nil {
		settings := config.GenerationSettings
		c.QuietOutput("Generation Settings:")
		c.QuietOutput("  Include Examples: %t", settings.IncludeExamples)
		c.QuietOutput("  Include Tests: %t", settings.IncludeTests)
		c.QuietOutput("  Include Documentation: %t", settings.IncludeDocs)
		c.QuietOutput("  Update Versions: %t", settings.UpdateVersions)
		c.QuietOutput("  Minimal Mode: %t", settings.MinimalMode)
		c.QuietOutput("  Backup Existing: %t", settings.BackupExisting)
		c.QuietOutput("  Overwrite Existing: %t", settings.OverwriteExisting)

		if len(settings.ExcludePatterns) > 0 {
			c.QuietOutput("  Exclude Patterns: %s", strings.Join(settings.ExcludePatterns, ", "))
		}
		if len(settings.IncludeOnlyPaths) > 0 {
			c.QuietOutput("  Include Only Paths: %s", strings.Join(settings.IncludeOnlyPaths, ", "))
		}
		c.QuietOutput("")
	}

	return nil
}

// outputJSON outputs data in JSON format
func (c *CLI) outputJSON(data interface{}) error {
	// This would use the existing JSON output functionality
	// Placeholder implementation
	fmt.Printf("🚧 JSON output is coming soon: %+v\n", data)
	return nil
}

// outputYAML outputs data in YAML format
func (c *CLI) outputYAML(data interface{}) error {
	// This would use YAML marshaling
	// Placeholder implementation
	fmt.Printf("YAML output not fully implemented: %+v\n", data)
	return nil
}
