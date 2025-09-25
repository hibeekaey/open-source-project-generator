// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
//
// This package handles comprehensive CLI operations including:
//   - Advanced project generation with multiple options
//   - Project validation and auditing
//   - Template management and discovery
//   - Configuration management
//   - Version and update management
//   - Cache management
//   - Logging and debugging support
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CLI implements the CLIInterface for comprehensive CLI operations.
//
// The CLI struct provides methods for:
//   - All documented CLI commands and flags
//   - Advanced project generation options
//   - Template management and validation
//   - Configuration management
//   - Auditing and validation capabilities
//   - Cache and offline mode support
//   - Interactive UI for guided project generation
type CLI struct {
	configManager          interfaces.ConfigManager
	validator              interfaces.ValidationEngine
	templateManager        interfaces.TemplateManager
	cacheManager           interfaces.CacheManager
	versionManager         interfaces.VersionManager
	auditEngine            interfaces.AuditEngine
	logger                 interfaces.Logger
	interactiveUI          interfaces.InteractiveUIInterface
	interactiveFlowManager *InteractiveFlowManager
	generatorVersion       string
	rootCmd                *cobra.Command
	verboseMode            bool
	quietMode              bool
	debugMode              bool
	exitCode               int
}

// NewCLI creates a new CLI instance with all required dependencies.
//
// Parameters:
//   - configManager: Handles configuration loading and validation
//   - validator: Provides input and project validation
//   - templateManager: Manages template operations
//   - cacheManager: Handles caching and offline mode
//   - versionManager: Manages version information and updates
//   - auditEngine: Provides auditing capabilities
//   - logger: Provides logging functionality
//   - version: Generator version string
//
// Returns:
//   - interfaces.CLIInterface: New CLI instance ready for use
func NewCLI(
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	templateManager interfaces.TemplateManager,
	cacheManager interfaces.CacheManager,
	versionManager interfaces.VersionManager,
	auditEngine interfaces.AuditEngine,
	logger interfaces.Logger,
	version string,
) interfaces.CLIInterface {
	// Create interactive UI with default configuration
	uiConfig := &ui.UIConfig{
		EnableColors:    true,
		EnableUnicode:   true,
		PageSize:        10,
		Timeout:         30 * time.Minute,
		AutoSave:        true,
		ShowBreadcrumbs: true,
		ShowShortcuts:   true,
		ConfirmOnQuit:   true,
	}
	interactiveUI := ui.NewInteractiveUI(logger, uiConfig)

	cli := &CLI{
		configManager:    configManager,
		validator:        validator,
		templateManager:  templateManager,
		cacheManager:     cacheManager,
		versionManager:   versionManager,
		auditEngine:      auditEngine,
		logger:           logger,
		interactiveUI:    interactiveUI,
		generatorVersion: version,
	}

	// Initialize interactive flow manager
	cli.interactiveFlowManager = NewInteractiveFlowManager(
		cli,
		templateManager,
		configManager,
		validator,
		logger,
		interactiveUI,
	)

	cli.setupCommands()
	return cli
}

// setupCommands initializes all CLI commands and their flags
func (c *CLI) setupCommands() {
	c.rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "Open Source Project Generator - Create production-ready projects with modern best practices",
		Long: `Generate production-ready projects with modern best practices.

Supports Go, Next.js, React, Android, iOS, Docker, Kubernetes, and Terraform.

Quick start:
  generator generate              # Interactive project creation
  generator list-templates        # Browse available templates
  generator validate <path>       # Check project structure
  generator audit <path>          # Security and quality analysis`,
		SilenceUsage:  true,
		SilenceErrors: false,
		Example: `  # Interactive project creation
  generator generate

  # Generate from config file
  generator generate --config project.yaml

  # Validate project
  generator validate ./my-project

  # Audit project security
  generator audit ./my-project`,
	}

	// Add global flags
	c.setupGlobalFlags()

	// Add all commands
	c.setupGenerateCommand()
	c.setupValidateCommand()
	c.setupAuditCommand()
	c.setupVersionCommand()
	c.setupConfigCommand()
	c.setupListTemplatesCommand()
	c.setupTemplateCommand()
	c.setupUpdateCommand()
	c.setupCacheCommand()
	c.setupLogsCommand()
}

// setupGlobalFlags adds global flags that apply to all commands
func (c *CLI) setupGlobalFlags() {
	c.rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging with detailed operation information")
	c.rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output (quiet mode)")
	c.rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging with performance metrics")
	c.rootCmd.PersistentFlags().String("log-level", "info", "Set log level (debug, info, warn, error, fatal)")
	c.rootCmd.PersistentFlags().Bool("log-json", false, "Output logs in JSON format")
	c.rootCmd.PersistentFlags().Bool("log-caller", false, "Include caller information in logs")
	c.rootCmd.PersistentFlags().Bool("non-interactive", false, "Run in non-interactive mode")
	c.rootCmd.PersistentFlags().String("output-format", "text", "Output format (text, json, yaml)")
}

// Run executes the CLI application with the provided arguments
func (c *CLI) Run(args []string) error {
	c.rootCmd.SetArgs(args)

	// Handle global flags before command execution
	if err := c.handleGlobalFlags(c.rootCmd); err != nil {
		return err
	}

	// Execute the command
	return c.rootCmd.Execute()
}

// handleGlobalFlags processes global flags and sets up the CLI state
func (c *CLI) handleGlobalFlags(cmd *cobra.Command) error {
	// Get global flags
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")
	debug, _ := cmd.Flags().GetBool("debug")
	logLevel, _ := cmd.Flags().GetString("log-level")
	logJSON, _ := cmd.Flags().GetBool("log-json")
	logCaller, _ := cmd.Flags().GetBool("log-caller")

	// Set CLI state
	c.verboseMode = verbose
	c.quietMode = quiet
	c.debugMode = debug

	// Configure logger if available
	if c.logger != nil {
		// Set log level
		switch logLevel {
		case "debug":
			// Logger will handle debug level
		case "info":
			// Logger will handle info level
		case "warn":
			// Logger will handle warn level
		case "error":
			// Logger will handle error level
		case "fatal":
			// Logger will handle fatal level
		}

		// Set JSON output
		if logJSON {
			c.logger.SetJSONOutput(true)
		}

		// Set caller info
		if logCaller {
			c.logger.SetCallerInfo(true)
		}
	}

	return nil
}

// PromptProjectDetails collects basic project configuration from user input.
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	if c.isNonInteractiveMode() {
		return nil, fmt.Errorf("🚫 %s %s",
			c.error("Interactive prompts not available in non-interactive mode."),
			c.info("Use environment variables or a configuration file instead"))
	}

	c.QuietOutput("Project Configuration")
	c.QuietOutput("====================")

	config := &models.ProjectConfig{}

	// Get project name
	fmt.Print("Project name: ")
	var name string
	if _, err := fmt.Scanln(&name); err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			c.error("Unable to read project name."),
			c.info("Please try typing it again or check your input"))
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

	// Get email (optional)
	fmt.Print("Email (optional): ")
	var email string
	_, _ = fmt.Scanln(&email) // Ignore error for optional input
	config.Email = strings.TrimSpace(email)

	// Get license (optional)
	fmt.Print("License (optional, default: Apache-2.0): ")
	var license string
	_, _ = fmt.Scanln(&license) // Ignore error for optional input
	config.License = strings.TrimSpace(license)
	if config.License == "" {
		config.License = "Apache-2.0"
	}

	return config, nil
}

// ConfirmGeneration asks the user to confirm project generation
func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	if c.isNonInteractiveMode() {
		return true // Auto-confirm in non-interactive mode
	}

	c.QuietOutput("\n📋 Project Summary:")
	c.QuietOutput("==================")
	c.QuietOutput("Name: %s", c.highlight(config.Name))
	c.QuietOutput("Description: %s", c.dim(config.Description))
	c.QuietOutput("License: %s", c.info(config.License))

	c.QuietOutput("\n🧩 Components:")
	// This would be populated based on the selected components
	c.QuietOutput("  (Components will be listed here)")

	fmt.Print("\nProceed with generation? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

// Helper methods

// Interface implementation methods (delegated to appropriate managers)

func (c *CLI) GenerateProject(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	// Delegate to template manager
	return c.templateManager.ProcessTemplate("", config, options.OutputPath)
}

func (c *CLI) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	// Delegate to validator
	result, err := c.validator.ValidateProject(path)
	if err != nil {
		return nil, err
	}

	// Convert issues to interface type
	issues := make([]interfaces.ValidationIssue, len(result.Issues))
	for i, issue := range result.Issues {
		issues[i] = interfaces.ValidationIssue{
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
		}
	}

	return &interfaces.ValidationResult{
		Valid:    result.Valid,
		Issues:   issues,
		Warnings: []interfaces.ValidationIssue{}, // No warnings field in models
	}, nil
}

func (c *CLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	// Delegate to audit engine
	return c.auditEngine.AuditProject(path, &options)
}

func (c *CLI) AuditProjectAdvanced(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	// Delegate to audit engine
	return c.auditEngine.AuditProject(path, options)
}

func (c *CLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	// Delegate to template manager
	return c.templateManager.ListTemplates(filter)
}

func (c *CLI) UpdateGenerator(options interface{}) (interface{}, error) {
	// Delegate to version manager
	return nil, fmt.Errorf("UpdateGenerator not implemented")
}

func (c *CLI) GetLogs(options interface{}) ([]string, error) {
	// Delegate to logger
	if c.logger == nil {
		return []string{}, nil
	}

	entries := c.logger.GetRecentEntries(50)
	logs := make([]string, len(entries))
	for i, entry := range entries {
		logs[i] = fmt.Sprintf("[%s] %s: %s", entry.Level, entry.Component, entry.Message)
	}

	return logs, nil
}

func (c *CLI) CheckCompatibility(version string) (*interfaces.CompatibilityResult, error) {
	// Delegate to version manager
	return c.versionManager.CheckCompatibility(version)
}

func (c *CLI) CheckUpdates() (*interfaces.UpdateInfo, error) {
	// Delegate to version manager
	return nil, fmt.Errorf("CheckUpdates not implemented")
}

func (c *CLI) CleanCache() error {
	// Delegate to cache manager
	return c.cacheManager.Clean()
}

func (c *CLI) ClearCache() error {
	// Delegate to cache manager
	return c.cacheManager.Clear()
}

func (c *CLI) ConfirmAdvancedGeneration(config *models.ProjectConfig, options *interfaces.AdvancedOptions) bool {
	// Use the same confirmation logic as regular generation
	return c.ConfirmGeneration(config)
}

func (c *CLI) DisableOfflineMode() error {
	// Delegate to cache manager
	return fmt.Errorf("DisableOfflineMode not implemented")
}

func (c *CLI) EditConfig() error {
	// Delegate to config manager
	return fmt.Errorf("EditConfig not implemented")
}

func (c *CLI) EnableOfflineMode() error {
	// Delegate to cache manager
	return fmt.Errorf("EnableOfflineMode not implemented")
}

func (c *CLI) ExportConfig(name string) error {
	c.VerboseOutput("📤 Exporting configuration to: %s", name)

	// Load current configuration
	config, err := c.configManager.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load default configuration: %w", err)
	}

	// Try to merge with environment variables
	envConfig, err := c.configManager.LoadFromEnvironment()
	if err == nil {
		config = c.configManager.MergeConfigurations(config, envConfig)
	}

	// Determine format from file extension
	format := "yaml"
	if strings.HasSuffix(strings.ToLower(name), ".json") {
		format = "json"
	} else if strings.HasSuffix(strings.ToLower(name), ".toml") {
		format = "toml"
	}

	// Generate configuration data
	var data []byte
	switch format {
	case "json":
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	case "toml":
		// For TOML, we'd need a TOML library
		return fmt.Errorf("TOML export not implemented yet")
	default: // yaml
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	}

	// Write to file
	if err := os.WriteFile(name, data, 0600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	c.VerboseOutput("✅ Configuration exported successfully")
	return nil
}

func (c *CLI) GenerateFromConfig(configPath string, options interfaces.GenerateOptions) error {
	// Delegate to template manager
	return fmt.Errorf("GenerateFromConfig not implemented")
}

func (c *CLI) GenerateReport(result, format, outputFile string) error {
	// Delegate to appropriate manager
	return fmt.Errorf("GenerateReport not implemented")
}

func (c *CLI) GenerateWithAdvancedOptions(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	// Delegate to template manager
	return fmt.Errorf("GenerateWithAdvancedOptions not implemented")
}

func (c *CLI) GetCacheStats() (*interfaces.CacheStats, error) {
	// Delegate to cache manager
	return nil, fmt.Errorf("GetCacheStats not implemented")
}

func (c *CLI) GetConfigurationSources() ([]interfaces.ConfigSource, error) {
	// Delegate to config manager
	return nil, fmt.Errorf("GetConfigurationSources not implemented")
}

func (c *CLI) GetLatestPackageVersions() (map[string]string, error) {
	// Delegate to version manager
	return nil, fmt.Errorf("GetLatestPackageVersions not implemented")
}

func (c *CLI) GetLogFileLocations() ([]string, error) {
	// Delegate to logger
	return nil, fmt.Errorf("GetLogFileLocations not implemented")
}

func (c *CLI) GetLogLevel() string {
	// Delegate to logger
	return "info"
}

func (c *CLI) GetPackageVersions() (map[string]string, error) {
	// Delegate to version manager
	return nil, fmt.Errorf("GetPackageVersions not implemented")
}

func (c *CLI) GetTemplateDependencies(templateName string) ([]string, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("GetTemplateDependencies not implemented")
}

func (c *CLI) GetTemplateInfo(templateName string) (*interfaces.TemplateInfo, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("GetTemplateInfo not implemented")
}

func (c *CLI) GetTemplateMetadata(templateName string) (*interfaces.TemplateMetadata, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("GetTemplateMetadata not implemented")
}

func (c *CLI) InstallUpdates() error {
	// Delegate to version manager
	return fmt.Errorf("InstallUpdates not implemented")
}

func (c *CLI) LoadConfiguration(paths []string) (*models.ProjectConfig, error) {
	// Delegate to config manager
	return nil, fmt.Errorf("LoadConfiguration not implemented")
}

func (c *CLI) MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error) {
	// Delegate to config manager
	return nil, fmt.Errorf("MergeConfigurations not implemented")
}

func (c *CLI) PromptAdvancedOptions() (*interfaces.AdvancedOptions, error) {
	// Delegate to interactive UI
	return nil, fmt.Errorf("PromptAdvancedOptions not implemented")
}

func (c *CLI) RepairCache() error {
	// Delegate to cache manager
	return fmt.Errorf("RepairCache not implemented")
}

func (c *CLI) RunNonInteractive(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	// Delegate to CLI
	return fmt.Errorf("RunNonInteractive not implemented")
}

func (c *CLI) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("SearchTemplates not implemented")
}

func (c *CLI) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	// Delegate to interactive UI
	return nil, fmt.Errorf("SelectTemplateInteractively not implemented")
}

func (c *CLI) SetConfig(key, value string) error {
	// Delegate to config manager
	return fmt.Errorf("SetConfig not implemented")
}

func (c *CLI) SetLogLevel(level string) error {
	// Delegate to logger
	return fmt.Errorf("SetLogLevel not implemented")
}

func (c *CLI) ShowCache() error {
	// Get cache statistics
	status, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	fmt.Println("Cache Information:")
	fmt.Println("=================")
	fmt.Printf("Location: %s\n", status.CacheLocation)
	fmt.Printf("Total Size: %d bytes\n", status.TotalSize)
	fmt.Printf("Health: %s\n", status.CacheHealth)
	fmt.Printf("Last Cleanup: %s\n", status.LastCleanup.Format(time.RFC3339))

	// Show cache validation
	fmt.Println("\nCache Validation:")
	fmt.Println("=================")
	if err := c.cacheManager.ValidateCache(); err != nil {
		fmt.Printf("❌ Cache validation failed: %v\n", err)
		fmt.Println("💡 Try running 'generator cache repair' to fix issues")
	} else {
		fmt.Println("✅ Cache is valid and healthy")
	}

	// Show offline mode status
	offlineMode := c.cacheManager.IsOfflineMode()
	fmt.Printf("\nOffline Mode: %t\n", offlineMode)
	if offlineMode {
		fmt.Println("📡 Running in offline mode - using cached templates and versions")
	} else {
		fmt.Println("🌐 Online mode - can fetch latest templates and versions")
	}

	return nil
}

func (c *CLI) ShowConfig() error {
	// Get configuration sources
	sources, err := c.configManager.GetConfigSources()
	if err != nil {
		return fmt.Errorf("failed to get configuration sources: %w", err)
	}

	fmt.Println("Configuration Sources:")
	fmt.Println("=====================")
	for _, source := range sources {
		status := "✓"
		if !source.Valid {
			status = "✗"
		}
		fmt.Printf("%s [%s] %s (priority: %d)\n", status, source.Type, source.Location, source.Priority)
	}

	// Load and display current configuration
	fmt.Println("\nCurrent Configuration:")
	fmt.Println("=====================")

	// Try to load defaults first
	config, err := c.configManager.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Try to merge with environment variables
	envConfig, err := c.configManager.LoadFromEnvironment()
	if err == nil {
		config = c.configManager.MergeConfigurations(config, envConfig)
	}

	// Display configuration values
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Organization: %s\n", config.Organization)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("License: %s\n", config.License)
	fmt.Printf("Author: %s\n", config.Author)
	fmt.Printf("Email: %s\n", config.Email)
	fmt.Printf("Repository: %s\n", config.Repository)
	fmt.Printf("Output Path: %s\n", config.OutputPath)

	// Display components
	fmt.Println("\nComponents:")
	fmt.Printf("  Frontend - NextJS App: %t\n", config.Components.Frontend.NextJS.App)
	fmt.Printf("  Frontend - NextJS Home: %t\n", config.Components.Frontend.NextJS.Home)
	fmt.Printf("  Frontend - NextJS Admin: %t\n", config.Components.Frontend.NextJS.Admin)
	fmt.Printf("  Frontend - NextJS Shared: %t\n", config.Components.Frontend.NextJS.Shared)
	fmt.Printf("  Backend - Go Gin: %t\n", config.Components.Backend.GoGin)
	fmt.Printf("  Mobile - Android: %t\n", config.Components.Mobile.Android)
	fmt.Printf("  Mobile - iOS: %t\n", config.Components.Mobile.IOS)
	fmt.Printf("  Infrastructure - Docker: %t\n", config.Components.Infrastructure.Docker)
	fmt.Printf("  Infrastructure - Kubernetes: %t\n", config.Components.Infrastructure.Kubernetes)
	fmt.Printf("  Infrastructure - Terraform: %t\n", config.Components.Infrastructure.Terraform)

	// Display versions if available
	if config.Versions != nil {
		fmt.Println("\nVersions:")
		fmt.Printf("  Node.js: %s\n", config.Versions.Node)
		fmt.Printf("  Go: %s\n", config.Versions.Go)
		if len(config.Versions.Packages) > 0 {
			fmt.Println("  Packages:")
			for pkg, version := range config.Versions.Packages {
				fmt.Printf("    %s: %s\n", pkg, version)
			}
		}
	}

	// Display environment variables
	fmt.Println("\nEnvironment Variables:")
	envVars := c.configManager.LoadEnvironmentVariables()
	if len(envVars) > 0 {
		for key, value := range envVars {
			fmt.Printf("  %s: %s\n", key, value)
		}
	} else {
		fmt.Println("  No relevant environment variables set")
	}

	return nil
}

func (c *CLI) ShowLogs() error {
	// Delegate to logger
	return fmt.Errorf("ShowLogs not implemented")
}

func (c *CLI) ShowRecentLogs(lines int, level string) error {
	// Delegate to logger
	return fmt.Errorf("ShowRecentLogs not implemented")
}

func (c *CLI) ShowVersion(options interfaces.VersionOptions) error {
	// Get current version info
	currentVersion := c.versionManager.GetCurrentVersion()

	// Basic version display
	fmt.Printf("Generator Version: %s\n", currentVersion)

	// Show build info if requested
	if options.ShowBuildInfo {
		if latestInfo, err := c.versionManager.GetLatestVersion(); err == nil {
			fmt.Printf("Latest Version: %s\n", latestInfo.Version)
			if !latestInfo.BuildDate.IsZero() {
				fmt.Printf("Build Date: %s\n", latestInfo.BuildDate.Format(time.RFC3339))
			}
		}
	}

	// Show update info if requested
	if options.CheckUpdates {
		updateInfo, err := c.versionManager.CheckForUpdates()
		if err == nil {
			if updateInfo.UpdateAvailable {
				fmt.Printf("Update Available: %s\n", updateInfo.LatestVersion)
				fmt.Printf("Current Version: %s\n", updateInfo.CurrentVersion)
			} else {
				fmt.Printf("You are running the latest version\n")
			}
		}
	}

	return nil
}

func (c *CLI) ValidateCache() error {
	fmt.Println("Validating cache...")

	// Validate cache integrity
	if err := c.cacheManager.ValidateCache(); err != nil {
		fmt.Printf("❌ Cache validation failed: %v\n", err)
		fmt.Println("💡 Try running 'generator cache repair' to fix issues")
		return err
	}

	fmt.Println("✅ Cache is valid and healthy")
	return nil
}

func (c *CLI) ValidateConfig() error {
	configLocation := c.configManager.GetConfigLocation()
	if configLocation == "" {
		return fmt.Errorf("no configuration file found")
	}

	// Validate the configuration file
	result, err := c.configManager.ValidateConfigFromFile(configLocation)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	// Display validation results
	fmt.Printf("Configuration file: %s\n", configLocation)
	fmt.Printf("Valid: %t\n", result.Valid)

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, validationError := range result.Errors {
			fmt.Printf("  ✗ %s: %s\n", validationError.Field, validationError.Message)
			if validationError.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", validationError.Suggestion)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠ %s: %s\n", warning.Field, warning.Message)
			if warning.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", warning.Suggestion)
			}
		}
	}

	// Display summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total properties: %d\n", result.Summary.TotalProperties)
	fmt.Printf("  Valid properties: %d\n", result.Summary.ValidProperties)
	fmt.Printf("  Errors: %d\n", result.Summary.ErrorCount)
	fmt.Printf("  Warnings: %d\n", result.Summary.WarningCount)
	fmt.Printf("  Missing required: %d\n", result.Summary.MissingRequired)

	if !result.Valid {
		return fmt.Errorf("configuration validation failed with %d errors", result.Summary.ErrorCount)
	}

	return nil
}

func (c *CLI) ValidateConfigurationSchema(config *models.ProjectConfig) error {
	// Delegate to validator
	return fmt.Errorf("ValidateConfigurationSchema not implemented")
}

func (c *CLI) ValidateCustomTemplate(templatePath string) (*interfaces.TemplateValidationResult, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("ValidateCustomTemplate not implemented")
}

func (c *CLI) ValidateProjectAdvanced(path string, options *interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	// Delegate to validator
	return nil, fmt.Errorf("ValidateProjectAdvanced not implemented")
}

func (c *CLI) ValidateTemplate(templateName string) (*interfaces.TemplateValidationResult, error) {
	// Delegate to template manager
	return nil, fmt.Errorf("ValidateTemplate not implemented")
}

// Additional methods needed for execution.go

func (c *CLI) validateGenerateOptions(options interfaces.GenerateOptions) error {
	var validationErrors []string

	// Validate output path
	if options.OutputPath != "" {
		if !filepath.IsAbs(options.OutputPath) && !strings.HasPrefix(options.OutputPath, "./") && !strings.HasPrefix(options.OutputPath, "../") {
			// Relative path without ./ prefix - this is okay, but we'll make it explicit
			options.OutputPath = "./" + options.OutputPath
		}

		// Check for invalid characters in path
		if strings.ContainsAny(options.OutputPath, "<>:\"|?*") {
			validationErrors = append(validationErrors, "output path contains invalid characters")
		}
	}

	// Validate template names
	for _, template := range options.Templates {
		if template == "" {
			validationErrors = append(validationErrors, "empty template name specified")
			continue
		}

		// Validate template name format
		if !isValidTemplateName(template) {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid template name '%s' - must contain only letters, numbers, hyphens, and underscores", template))
		}
	}

	// Validate conflicting options
	if options.Offline && options.UpdateVersions {
		validationErrors = append(validationErrors, "cannot use --offline and --update-versions together")
	}

	if options.Minimal && options.IncludeExamples {
		c.WarningOutput("Using --minimal with --include-examples may result in minimal examples only")
	}

	// Validate dry-run with force
	if options.DryRun && options.Force {
		c.WarningOutput("--force flag has no effect in dry-run mode")
	}

	if len(validationErrors) > 0 {
		return &interfaces.CLIError{
			Type:        interfaces.ErrorTypeValidation,
			Message:     "generate options validation failed",
			Code:        interfaces.ErrorCodeValidationFailed,
			Details:     map[string]any{"errors": validationErrors},
			Suggestions: []string{"Fix the validation errors and try again"},
		}
	}

	return nil
}

// isValidTemplateName validates template name format
func isValidTemplateName(name string) bool {
	// Allow letters, numbers, hyphens, underscores, and dots
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 50
}

func (c *CLI) routeToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🚀 Starting %s generation...", mode)

	switch mode {
	case "config-file":
		return c.handleConfigFileGeneration(configPath, options)
	case "non-interactive":
		return c.handleNonInteractiveGeneration(options)
	case "interactive":
		return c.handleInteractiveGeneration(options)
	default:
		return fmt.Errorf("unsupported generation mode: %s", mode)
	}
}

func (c *CLI) generateValidationReport(result *interfaces.ValidationResult, format, outputFile string) error {
	c.VerboseOutput("📄 Generating validation report...")

	// Determine output file
	if outputFile == "" {
		outputFile = "validation-report." + format
	}

	// Generate report based on format
	switch format {
	case "json":
		return c.generateJSONReport(result, outputFile)
	case "yaml":
		return c.generateYAMLReport(result, outputFile)
	case "html":
		return c.generateHTMLReport(result, outputFile)
	default:
		return c.generateTextReport(result, outputFile)
	}
}

// Handler methods for different generation modes
func (c *CLI) handleConfigFileGeneration(configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("📄 Loading your project configuration...")
	c.DebugOutput("📄 Using configuration: %s", configPath)

	// Validate configuration file exists and is readable
	if err := c.validateConfigurationFile(configPath); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Load configuration from file
	config, err := c.loadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration from %s: %w", configPath, err)
	}

	// Execute project generation
	return c.executeGenerationWorkflow(config, options)
}

func (c *CLI) handleNonInteractiveGeneration(options interfaces.GenerateOptions) error {
	c.VerboseOutput("Starting non-interactive generation")

	// Check for required environment variables
	if err := c.validateNonInteractiveEnvironment(); err != nil {
		return fmt.Errorf("non-interactive environment validation failed: %w", err)
	}

	// Load configuration from environment variables
	config, err := c.loadConfigFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load environment configuration: %w", err)
	}

	// Execute project generation
	return c.executeGenerationWorkflow(config, options)
}

func (c *CLI) handleInteractiveGeneration(options interfaces.GenerateOptions) error {
	c.VerboseOutput("Starting interactive generation")

	// Validate interactive environment
	if err := c.validateInteractiveEnvironment(); err != nil {
		return fmt.Errorf("interactive environment validation failed: %w", err)
	}

	// Use the interactive flow manager
	ctx := context.Background()
	return c.interactiveFlowManager.RunInteractiveFlow(ctx, options)
}

func (c *CLI) runInteractiveProjectConfiguration(ctx context.Context) (*models.ProjectConfig, error) {
	// Interactive project configuration
	return nil, fmt.Errorf("runInteractiveProjectConfiguration not implemented")
}

// Helper methods for generation workflow
func (c *CLI) validateConfigurationFile(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("configuration file path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("🚫 %s %s %s",
			c.error("Configuration file"),
			c.highlight(fmt.Sprintf("'%s'", configPath)),
			c.info("not found. Check the file path and try again"))
	}

	// Check if file is readable
	file, err := os.Open(configPath) // #nosec G304 - configPath is validated before use
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Unable to read configuration file."),
			c.info("Check file permissions and ensure it's not corrupted"))
	}
	_ = file.Close()

	c.DebugOutput("✅ Configuration file looks good: %s", configPath) // #nosec G304 - configPath is validated and only used for logging
	return nil
}

func (c *CLI) loadConfigFromFile(configPath string) (*models.ProjectConfig, error) {
	if c.configManager == nil {
		return nil, fmt.Errorf("configuration manager not initialized")
	}

	config, err := c.configManager.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from file: %w", err)
	}

	c.VerboseOutput("✅ Configuration loaded from file: %s", configPath)
	return config, nil
}

func (c *CLI) loadConfigFromEnvironment() (*models.ProjectConfig, error) {
	c.VerboseOutput("🌍 Loading configuration from environment variables...")

	// Load environment configuration
	envConfig, err := c.loadEnvironmentConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment configuration: %w", err)
	}

	// Convert environment config to project config
	config, err := c.convertEnvironmentConfigToProjectConfig(envConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert environment configuration: %w", err)
	}

	// Validate required fields for non-interactive mode
	if config.Name == "" {
		return nil, fmt.Errorf("project name is required for non-interactive mode (set GENERATOR_PROJECT_NAME)")
	}

	c.VerboseOutput("✅ Configuration loaded from environment variables")
	return config, nil
}

func (c *CLI) executeGenerationWorkflow(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🚀 Starting project generation for: %s", config.Name)

	// Validate configuration if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("🔍 Validating project configuration...")
		if err := c.validateGenerateConfiguration(config, options); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
		c.VerboseOutput("✅ Configuration validation passed")
	}

	// Set output path from options or config
	outputPath := c.determineOutputPath(config, options)
	c.VerboseOutput("📁 Output directory: %s", outputPath)

	// Handle offline mode
	if options.Offline {
		c.VerboseOutput("📡 Running in offline mode - using cached templates and versions")
		if c.cacheManager != nil {
			if err := c.cacheManager.EnableOfflineMode(); err != nil {
				c.WarningOutput("⚠️  Couldn't enable offline mode: %v", err)
			}
		}
	}

	// Handle version updates
	if options.UpdateVersions && !options.Offline {
		c.VerboseOutput("📦 Fetching latest package versions...")
		if err := c.updatePackageVersions(config); err != nil {
			c.WarningOutput("⚠️  Couldn't update package versions: %v", err)
		}
	}

	// Log CI environment information if detected
	ci := c.detectCIEnvironment()
	if ci.IsCI {
		c.VerboseOutput("🤖 Detected CI environment: %s", ci.Provider)
		if ci.BuildID != "" {
			c.VerboseOutput("   Build ID: %s", ci.BuildID)
		}
		if ci.Branch != "" {
			c.VerboseOutput("   Branch: %s", ci.Branch)
		}
	}

	// Pre-generation checks
	if err := c.performPreGenerationChecks(outputPath, options); err != nil {
		return fmt.Errorf("pre-generation checks failed: %w", err)
	}

	// Handle dry run mode
	if options.DryRun {
		c.QuietOutput("🔍 %s - would generate project %s in directory %s",
			c.warning("Dry run mode"),
			c.highlight(fmt.Sprintf("'%s'", config.Name)),
			c.info(fmt.Sprintf("'%s'", outputPath)))
		c.displayProjectSummary(config)
		return nil
	}

	// Execute the actual project generation
	c.VerboseOutput("🏗️  Generating project structure...")
	if err := c.generateProjectFromComponents(config, outputPath, options); err != nil {
		return fmt.Errorf("project generation failed: %w", err)
	}

	// Post-generation tasks
	if err := c.performPostGenerationTasks(config, outputPath, options); err != nil {
		c.WarningOutput("⚠️  Some post-generation tasks failed: %v", err)
	}

	c.SuccessOutput("🎉 Project %s %s!", c.highlight(fmt.Sprintf("'%s'", config.Name)), c.success("generated successfully"))
	c.displayGenerationSummary(config, outputPath)

	return nil
}

// Additional helper methods needed for the generation workflow
func (c *CLI) validateNonInteractiveEnvironment() error {
	// Check for required environment variables
	requiredEnvVars := []string{"GENERATOR_PROJECT_NAME"}
	missingVars := []string{}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("required environment variables missing: %s", strings.Join(missingVars, ", "))
	}

	c.DebugOutput("Non-interactive environment validation passed")
	return nil
}

func (c *CLI) validateInteractiveEnvironment() error {
	// Check if we're in a terminal
	if c.isNonInteractiveMode() {
		return fmt.Errorf("interactive mode not available in non-interactive environment (CI, piped input, etc.)")
	}

	// Check if interactive UI is available
	if c.interactiveUI == nil {
		return fmt.Errorf("interactive UI not initialized")
	}

	// Check if interactive flow manager is available
	if c.interactiveFlowManager == nil {
		return fmt.Errorf("interactive flow manager not initialized")
	}

	c.DebugOutput("Interactive environment validation passed")
	return nil
}

func (c *CLI) determineOutputPath(config *models.ProjectConfig, options interfaces.GenerateOptions) string {
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = config.OutputPath
	}
	if outputPath == "" {
		outputPath = "./output/generated"
	}

	// Always append project name to the output path
	return filepath.Join(outputPath, config.Name)
}

func (c *CLI) displayProjectSummary(config *models.ProjectConfig) {
	c.QuietOutput("\n%s", c.highlight("📋 Project Summary:"))
	c.QuietOutput("%s", c.dim("=================="))
	c.QuietOutput("Name: %s", c.success(config.Name))
	if config.Organization != "" {
		c.QuietOutput("Organization: %s", c.info(config.Organization))
	}
	if config.Description != "" {
		c.QuietOutput("Description: %s", c.dim(config.Description))
	}
	c.QuietOutput("License: %s", c.info(config.License))

	c.QuietOutput("\n%s", c.highlight("🧩 Components:"))
	if c.hasFrontendComponents(config) {
		c.QuietOutput("  %s %s", c.success("✅"), c.info("Frontend (Next.js)"))
	}
	if c.hasBackendComponents(config) {
		c.QuietOutput("  %s %s", c.success("✅"), c.info("Backend (Go Gin)"))
	}
	if c.hasMobileComponents(config) {
		c.QuietOutput("  %s %s", c.success("✅"), c.info("Mobile"))
	}
	if c.hasInfrastructureComponents(config) {
		c.QuietOutput("  %s %s", c.success("✅"), c.info("Infrastructure"))
	}
}

func (c *CLI) displayGenerationSummary(config *models.ProjectConfig, outputPath string) {
	c.QuietOutput("\n%s", c.highlight("📊 Generation Summary:"))
	c.QuietOutput("%s", c.dim("====================="))
	c.QuietOutput("Project: %s", c.success(config.Name))
	c.QuietOutput("Location: %s", c.info(outputPath))
	c.QuietOutput("Components generated: %s", c.success(fmt.Sprintf("%d", c.countSelectedComponents(config))))

	c.QuietOutput("\n%s", c.highlight("🚀 Next Steps:"))
	c.QuietOutput("%s. Navigate to your project: %s", c.info("1"), c.colorize(ColorCyan, fmt.Sprintf("cd %s", outputPath)))
	c.QuietOutput("%s. Review the generated %s for setup instructions", c.info("2"), c.highlight("README.md"))
	c.QuietOutput("%s. Install dependencies and start development", c.info("3"))
}

// Additional helper methods for component detection and counting
func (c *CLI) hasFrontendComponents(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared
}

func (c *CLI) hasBackendComponents(config *models.ProjectConfig) bool {
	return config.Components.Backend.GoGin
}

func (c *CLI) hasMobileComponents(config *models.ProjectConfig) bool {
	return config.Components.Mobile.Android ||
		config.Components.Mobile.IOS ||
		config.Components.Mobile.Shared
}

func (c *CLI) hasInfrastructureComponents(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

func (c *CLI) countSelectedComponents(config *models.ProjectConfig) int {
	count := 0
	if c.hasFrontendComponents(config) {
		count++
	}
	if c.hasBackendComponents(config) {
		count++
	}
	if c.hasMobileComponents(config) {
		count++
	}
	if c.hasInfrastructureComponents(config) {
		count++
	}
	return count
}

// Placeholder methods for missing functionality
func (c *CLI) validateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	// Basic validation - can be expanded later
	if config.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

func (c *CLI) updatePackageVersions(config *models.ProjectConfig) error {
	// Placeholder - would update package versions
	return fmt.Errorf("updatePackageVersions not implemented")
}

func (c *CLI) performPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	// Placeholder - would perform pre-generation checks
	return nil
}

func (c *CLI) generateProjectFromComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🏗️  Generating project structure for: %s", config.Name)

	// Create output directory
	if err := os.MkdirAll(outputPath, 0750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate README.md
	if err := c.generateREADME(config, outputPath); err != nil {
		c.WarningOutput("⚠️  Failed to generate README: %v", err)
	}

	// Generate frontend components
	if c.hasFrontendComponents(config) {
		if err := c.generateFrontendComponents(config, outputPath, options); err != nil {
			c.WarningOutput("⚠️  Failed to generate frontend components: %v", err)
		}
	}

	// Generate backend components
	if c.hasBackendComponents(config) {
		if err := c.generateBackendComponents(config, outputPath, options); err != nil {
			c.WarningOutput("⚠️  Failed to generate backend components: %v", err)
		}
	}

	// Generate mobile components
	if c.hasMobileComponents(config) {
		if err := c.generateMobileComponents(config, outputPath, options); err != nil {
			c.WarningOutput("⚠️  Failed to generate mobile components: %v", err)
		}
	}

	// Generate infrastructure components
	if c.hasInfrastructureComponents(config) {
		if err := c.generateInfrastructureComponents(config, outputPath, options); err != nil {
			c.WarningOutput("⚠️  Failed to generate infrastructure components: %v", err)
		}
	}

	c.VerboseOutput("✅ Project structure generated successfully")
	return nil
}

func (c *CLI) performPostGenerationTasks(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	// Initialize git repository if not in minimal mode
	if !options.Minimal {
		if err := c.initializeGitRepository(outputPath); err != nil {
			c.VerboseOutput("⚠️  Git initialization skipped: %v", err)
		}
	}

	// Set file permissions
	if err := c.setFilePermissions(outputPath); err != nil {
		c.VerboseOutput("⚠️  File permission setup skipped: %v", err)
	}

	return nil
}

// Helper methods for generating different components
func (c *CLI) generateREADME(config *models.ProjectConfig, outputPath string) error {
	readmePath := filepath.Join(outputPath, "README.md")

	readmeContent := fmt.Sprintf(`# %s

%s

## Project Overview

- **Name**: %s
- **Organization**: %s
- **License**: %s
- **Author**: %s

## Components

%s

## Getting Started

1. Install dependencies
2. Run the development server
3. Build for production

## Development

This project was generated using the Open Source Project Generator.

## License

%s
`,
		config.Name,
		config.Description,
		config.Name,
		config.Organization,
		config.License,
		config.Author,
		c.generateComponentList(config),
		config.License,
	)

	return os.WriteFile(readmePath, []byte(readmeContent), 0644) // #nosec G306 -- Project files need to be readable
}

func (c *CLI) generateComponentList(config *models.ProjectConfig) string {
	var components []string

	if c.hasFrontendComponents(config) {
		components = append(components, "- Frontend (Next.js)")
	}
	if c.hasBackendComponents(config) {
		components = append(components, "- Backend (Go Gin)")
	}
	if c.hasMobileComponents(config) {
		components = append(components, "- Mobile (Android/iOS)")
	}
	if c.hasInfrastructureComponents(config) {
		components = append(components, "- Infrastructure (Docker, Kubernetes, Terraform)")
	}

	if len(components) == 0 {
		return "- No components selected"
	}

	return strings.Join(components, "\n")
}

func (c *CLI) generateFrontendComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	frontendPath := filepath.Join(outputPath, "frontend")
	if err := os.MkdirAll(frontendPath, 0750); err != nil {
		return err
	}

	// Generate package.json
	packageJson := `{
  "name": "` + config.Name + `-frontend",
  "version": "1.0.0",
  "description": "` + config.Description + `",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  }
}`

	return os.WriteFile(filepath.Join(frontendPath, "package.json"), []byte(packageJson), 0644) // #nosec G306 -- Project files need to be readable
}

func (c *CLI) generateBackendComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	backendPath := filepath.Join(outputPath, "backend")
	if err := os.MkdirAll(backendPath, 0750); err != nil {
		return err
	}

	// Generate go.mod
	goMod := `module ` + config.Name + `-backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`

	return os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644) // #nosec G306 -- Project files need to be readable
}

func (c *CLI) generateMobileComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	mobilePath := filepath.Join(outputPath, "mobile")
	if err := os.MkdirAll(mobilePath, 0750); err != nil {
		return err
	}

	// Generate basic mobile structure
	readmePath := filepath.Join(mobilePath, "README.md")
	readmeContent := `# Mobile App

This directory contains the mobile application components.

## Android
- Native Android development with Kotlin
- Modern Android architecture patterns

## iOS  
- Native iOS development with Swift
- SwiftUI and UIKit support

## Shared
- Shared business logic
- Common utilities and models
`

	return os.WriteFile(readmePath, []byte(readmeContent), 0644) // #nosec G306 -- Project files need to be readable
}

func (c *CLI) generateInfrastructureComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	infraPath := filepath.Join(outputPath, "infrastructure")
	if err := os.MkdirAll(infraPath, 0750); err != nil {
		return err
	}

	// Generate Docker configuration
	if config.Components.Infrastructure.Docker {
		dockerfile := `FROM node:18-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ .
RUN npm run build

FROM golang:1.21-alpine AS backend
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=frontend /app/out ./frontend
COPY --from=backend /app/main .
EXPOSE 8080
CMD ["./main"]
`
		// #nosec G306 -- Project files need to be readable
		if err := os.WriteFile(filepath.Join(infraPath, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
			return fmt.Errorf("failed to write Dockerfile: %w", err)
		}
	}

	// Generate Kubernetes configuration
	if config.Components.Infrastructure.Kubernetes {
		k8sPath := filepath.Join(infraPath, "k8s")
		if err := os.MkdirAll(k8sPath, 0750); err != nil {
			return fmt.Errorf("failed to create k8s directory: %w", err)
		}

		deployment := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + config.Name + `
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ` + config.Name + `
  template:
    metadata:
      labels:
        app: ` + config.Name + `
    spec:
      containers:
      - name: ` + config.Name + `
        image: ` + config.Name + `:latest
        ports:
        - containerPort: 8080
`
		// #nosec G306 -- Project files need to be readable
		if err := os.WriteFile(filepath.Join(k8sPath, "deployment.yaml"), []byte(deployment), 0644); err != nil {
			return fmt.Errorf("failed to write deployment.yaml: %w", err)
		}
	}

	// Generate Terraform configuration
	if config.Components.Infrastructure.Terraform {
		tfPath := filepath.Join(infraPath, "terraform")
		if err := os.MkdirAll(tfPath, 0750); err != nil {
			return fmt.Errorf("failed to create terraform directory: %w", err)
		}

		mainTf := `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}

resource "aws_ecs_cluster" "main" {
  name = "` + config.Name + `-cluster"
}
`
		// #nosec G306 -- Project files need to be readable
		if err := os.WriteFile(filepath.Join(tfPath, "main.tf"), []byte(mainTf), 0644); err != nil {
			return fmt.Errorf("failed to write main.tf: %w", err)
		}
	}

	return nil
}

func (c *CLI) initializeGitRepository(outputPath string) error {
	c.VerboseOutput("📝 Initializing git repository...")

	cmd := exec.Command("git", "init")
	cmd.Dir = outputPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	c.VerboseOutput("✅ Git repository initialized")
	return nil
}

func (c *CLI) setFilePermissions(outputPath string) error {
	c.VerboseOutput("🔒 Setting file permissions...")

	// Make script files executable
	scriptsDir := filepath.Join(outputPath, "Scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		err := filepath.Walk(scriptsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".sh") || strings.HasSuffix(path, ".py")) {
				return os.Chmod(path, 0750) // #nosec G302 -- Scripts need to be executable
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to set script permissions: %w", err)
		}
	}

	c.VerboseOutput("✅ File permissions set")
	return nil
}

// Report generation helper methods
func (c *CLI) generateJSONReport(result *interfaces.ValidationResult, outputFile string) error {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(outputFile, jsonData, 0600)
}

func (c *CLI) generateYAMLReport(result *interfaces.ValidationResult, outputFile string) error {
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return os.WriteFile(outputFile, yamlData, 0600)
}

func (c *CLI) generateHTMLReport(result *interfaces.ValidationResult, outputFile string) error {
	statusClass := "error"
	statusText := "Invalid"
	if result.Valid {
		statusClass = "success"
		statusText = "Valid"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Validation Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f4f4f4; padding: 20px; border-radius: 5px; }
        .error { color: #d32f2f; }
        .warning { color: #f57c00; }
        .success { color: #388e3c; }
        .issue { margin: 10px 0; padding: 10px; border-left: 4px solid #ccc; }
        .error-issue { border-left-color: #d32f2f; }
        .warning-issue { border-left-color: #f57c00; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Validation Report</h1>
        <p><strong>Status:</strong> <span class="%s">%s</span></p>
        <p><strong>Total Issues:</strong> %d</p>
        <p><strong>Warnings:</strong> %d</p>
    </div>
    
    <h2>Issues</h2>
    %s
    
    <h2>Warnings</h2>
    %s
</body>
</html>`,
		statusClass,
		statusText,
		len(result.Issues),
		len(result.Warnings),
		c.generateHTMLIssues(result.Issues, "error-issue"),
		c.generateHTMLIssues(result.Warnings, "warning-issue"),
	)

	return os.WriteFile(outputFile, []byte(html), 0600)
}

func (c *CLI) generateTextReport(result *interfaces.ValidationResult, outputFile string) error {
	var report strings.Builder

	report.WriteString("Validation Report\n")
	report.WriteString("================\n\n")
	statusText := "Invalid"
	if result.Valid {
		statusText = "Valid"
	}
	report.WriteString(fmt.Sprintf("Status: %s\n", statusText))
	report.WriteString(fmt.Sprintf("Total Issues: %d\n", len(result.Issues)))
	report.WriteString(fmt.Sprintf("Warnings: %d\n\n", len(result.Warnings)))

	if len(result.Issues) > 0 {
		report.WriteString("Issues:\n")
		report.WriteString("-------\n")
		for _, issue := range result.Issues {
			report.WriteString(fmt.Sprintf("- %s: %s\n", issue.Severity, issue.Message))
			if issue.File != "" {
				report.WriteString(fmt.Sprintf("  File: %s:%d:%d\n", issue.File, issue.Line, issue.Column))
			}
		}
		report.WriteString("\n")
	}

	if len(result.Warnings) > 0 {
		report.WriteString("Warnings:\n")
		report.WriteString("---------\n")
		for _, warning := range result.Warnings {
			report.WriteString(fmt.Sprintf("- %s: %s\n", warning.Severity, warning.Message))
			if warning.File != "" {
				report.WriteString(fmt.Sprintf("  File: %s:%d:%d\n", warning.File, warning.Line, warning.Column))
			}
		}
	}

	return os.WriteFile(outputFile, []byte(report.String()), 0600)
}

func (c *CLI) generateHTMLIssues(issues []interfaces.ValidationIssue, cssClass string) string {
	if len(issues) == 0 {
		return "<p>None</p>"
	}

	var html strings.Builder
	for _, issue := range issues {
		html.WriteString(fmt.Sprintf(`<div class="issue %s">
            <strong>%s:</strong> %s
            <br><small>File: %s:%d:%d</small>
        </div>`,
			cssClass,
			issue.Severity,
			issue.Message,
			issue.File,
			issue.Line,
			issue.Column,
		))
	}

	return html.String()
}

// Audit report generation
func (c *CLI) generateAuditReport(result *interfaces.AuditResult, format, outputFile string) error {
	c.VerboseOutput("📄 Generating audit report...")

	// Determine output file
	if outputFile == "" {
		outputFile = "audit-report." + format
	}

	// Generate report based on format
	switch format {
	case "json":
		return c.generateAuditJSONReport(result, outputFile)
	case "yaml":
		return c.generateAuditYAMLReport(result, outputFile)
	case "html":
		return c.generateAuditHTMLReport(result, outputFile)
	default:
		return c.generateAuditTextReport(result, outputFile)
	}
}

func (c *CLI) generateAuditJSONReport(result *interfaces.AuditResult, outputFile string) error {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal audit JSON: %w", err)
	}

	return os.WriteFile(outputFile, jsonData, 0600)
}

func (c *CLI) generateAuditYAMLReport(result *interfaces.AuditResult, outputFile string) error {
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal audit YAML: %w", err)
	}

	return os.WriteFile(outputFile, yamlData, 0600)
}

func (c *CLI) generateAuditHTMLReport(result *interfaces.AuditResult, outputFile string) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Audit Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f4f4f4; padding: 20px; border-radius: 5px; }
        .score { font-size: 24px; font-weight: bold; }
        .score-high { color: #388e3c; }
        .score-medium { color: #f57c00; }
        .score-low { color: #d32f2f; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .recommendation { margin: 10px 0; padding: 10px; background: #f9f9f9; border-left: 4px solid #2196f3; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Audit Report</h1>
        <p><strong>Overall Score:</strong> <span class="score %s">%.1f/100</span></p>
    </div>
    
    <div class="section">
        <h2>Security</h2>
        <p><strong>Score:</strong> %.1f/100</p>
        <p><strong>Compatible:</strong> %t</p>
    </div>
    
    <div class="section">
        <h2>Quality</h2>
        <p><strong>Score:</strong> %.1f/100</p>
    </div>
    
    <div class="section">
        <h2>Licenses</h2>
        <p><strong>Compatible:</strong> %t</p>
    </div>
    
    <div class="section">
        <h2>Performance</h2>
        <p><strong>Score:</strong> %.1f/100</p>
    </div>
    
    <div class="section">
        <h2>Recommendations</h2>
        %s
    </div>
</body>
</html>`,
		c.getScoreClass(result.OverallScore),
		result.OverallScore,
		c.getScoreValue(result.Security),
		c.getCompatibleValue(result.Security),
		c.getScoreValue(result.Quality),
		c.getCompatibleValue(result.Licenses),
		c.getScoreValue(result.Performance),
		c.generateAuditRecommendations(result.Recommendations),
	)

	return os.WriteFile(outputFile, []byte(html), 0600)
}

func (c *CLI) generateAuditTextReport(result *interfaces.AuditResult, outputFile string) error {
	var report strings.Builder

	report.WriteString("Audit Report\n")
	report.WriteString("============\n\n")
	report.WriteString(fmt.Sprintf("Overall Score: %.1f/100\n\n", result.OverallScore))

	if result.Security != nil {
		report.WriteString(fmt.Sprintf("Security Score: %.1f/100\n\n", result.Security.Score))
	}

	if result.Quality != nil {
		report.WriteString(fmt.Sprintf("Quality Score: %.1f/100\n\n", result.Quality.Score))
	}

	if result.Licenses != nil {
		report.WriteString(fmt.Sprintf("License Compatible: %t\n\n", result.Licenses.Compatible))
	}

	if result.Performance != nil {
		report.WriteString(fmt.Sprintf("Performance Score: %.1f/100\n\n", result.Performance.Score))
	}

	if len(result.Recommendations) > 0 {
		report.WriteString("Recommendations:\n")
		report.WriteString("----------------\n")
		for _, rec := range result.Recommendations {
			report.WriteString(fmt.Sprintf("- %s\n", rec))
		}
	}

	return os.WriteFile(outputFile, []byte(report.String()), 0600)
}

func (c *CLI) getScoreClass(score float64) string {
	if score >= 80 {
		return "score-high"
	} else if score >= 60 {
		return "score-medium"
	}
	return "score-low"
}

func (c *CLI) getScoreValue(section interface{}) float64 {
	if section == nil {
		return 0.0
	}

	// Type assertion to get the score
	switch s := section.(type) {
	case *interfaces.SecurityAuditResult:
		return s.Score
	case *interfaces.QualityAuditResult:
		return s.Score
	case *interfaces.LicenseAuditResult:
		return s.Score
	case *interfaces.PerformanceAuditResult:
		return s.Score
	default:
		return 0.0
	}
}

func (c *CLI) getCompatibleValue(section interface{}) bool {
	if section == nil {
		return false
	}

	// Type assertion to get the compatible field
	switch s := section.(type) {
	case *interfaces.LicenseAuditResult:
		return s.Compatible
	default:
		return false
	}
}

func (c *CLI) generateAuditRecommendations(recommendations []string) string {
	if len(recommendations) == 0 {
		return "<p>None</p>"
	}

	var html strings.Builder
	for _, rec := range recommendations {
		html.WriteString(fmt.Sprintf(`<div class="recommendation">%s</div>`, rec))
	}

	return html.String()
}

func (c *CLI) validateAndNormalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "interactive", "i":
		return "interactive"
	case "non-interactive", "noninteractive", "ni", "auto":
		return "non-interactive"
	case "config-file", "config", "file", "cf":
		return "config-file"
	default:
		c.WarningOutput("Unknown mode '%s', defaulting to interactive", mode)
		return "interactive"
	}
}

func (c *CLI) detectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string {
	// Priority 1: Explicit mode flag (highest priority)
	if explicitMode != "" {
		c.DebugOutput("🎯 You specified %s mode explicitly", explicitMode)
		return c.validateAndNormalizeMode(explicitMode)
	}

	// Priority 2: Configuration file mode
	if configPath != "" {
		c.DebugOutput("📄 Found configuration file: %s", configPath)
		return "config-file"
	}

	// Priority 3: Explicit non-interactive flag (overrides auto-detection)
	if nonInteractive {
		c.DebugOutput("🤖 Non-interactive mode requested")
		return "non-interactive"
	}

	// Priority 4: Auto-detect non-interactive environment (CI, piped input, etc.)
	if c.isNonInteractiveMode() {
		c.DebugOutput("🤖 Detected automated environment (CI/scripts)")
		return "non-interactive"
	}

	// Priority 5: Explicit interactive flag or default
	if interactive {
		c.DebugOutput("👤 Interactive mode selected")
		return "interactive"
	}

	// Fallback: Interactive mode (should not reach here with current logic)
	c.DebugOutput("👤 Defaulting to interactive mode")
	return "interactive"
}

func (c *CLI) validateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	conflictCount := 0

	if nonInteractive {
		conflictCount++
	}
	if forceInteractive {
		conflictCount++
	}
	if forceNonInteractive {
		conflictCount++
	}
	if explicitMode != "" {
		conflictCount++
	}

	if conflictCount > 1 {
		return fmt.Errorf("🚫 %s %s",
			c.error("Multiple mode flags detected."),
			c.info("Please use only one mode flag at a time"))
	}

	// Validate explicit mode value
	if explicitMode != "" {
		validModes := []string{"interactive", "non-interactive", "config-file"}
		isValid := false
		for _, mode := range validModes {
			if explicitMode == mode {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("🚫 %s %s %s",
				c.error(fmt.Sprintf("'%s' is not a valid mode.", explicitMode)),
				c.info("Available modes:"),
				c.highlight(strings.Join(validModes, ", ")))
		}
	}

	return nil
}

func (c *CLI) applyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool) {
	// Handle explicit mode
	if explicitMode != "" {
		switch explicitMode {
		case "non-interactive":
			return true, false
		case "interactive":
			return false, true
		case "config-file":
			// Config file mode doesn't change interactive flags
			return nonInteractive, interactive
		}
	}

	// Handle force flags
	if forceInteractive {
		return false, true
	}
	if forceNonInteractive {
		return true, false
	}

	// Return original values if no overrides
	return nonInteractive, interactive
}
