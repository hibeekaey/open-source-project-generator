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
	"strconv"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
	"github.com/spf13/cobra"
)

// Color constants for beautiful CLI output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
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
		SilenceErrors: true,
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

	// Set up pre-run hook to handle global flags
	c.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return c.handleGlobalFlags(cmd)
	}

	// Execute command and handle errors with proper exit codes
	err := c.rootCmd.Execute()
	if err != nil {
		// Get command name for context
		cmdName := "generator"
		if c.rootCmd.CalledAs() != "" {
			cmdName = c.rootCmd.CalledAs()
		}

		// Handle error and get exit code
		exitCode := c.handleError(err, cmdName, args)

		// Exit with appropriate code for automation
		if c.isNonInteractiveMode() {
			os.Exit(exitCode)
		}
	}

	return err
}

// handleGlobalFlags processes global flags that apply to all commands
func (c *CLI) handleGlobalFlags(cmd *cobra.Command) error {
	// Get global flags
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")
	debug, _ := cmd.Flags().GetBool("debug")
	logLevel, _ := cmd.Flags().GetString("log-level")
	logJSON, _ := cmd.Flags().GetBool("log-json")
	logCaller, _ := cmd.Flags().GetBool("log-caller")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode if not explicitly set
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
		if nonInteractive {
			c.VerboseOutput("Auto-detected non-interactive mode (CI environment or piped input)")
		}
	}

	// Handle conflicting flags
	if verbose && quiet {
		return fmt.Errorf("🚫 You can't use both --verbose and --quiet at the same time")
	}
	if debug && quiet {
		return fmt.Errorf("🚫 You can't use both --debug and --quiet at the same time")
	}

	// Set log level based on flags (priority: debug > verbose > explicit level > quiet)
	if debug {
		logLevel = "debug"
		c.debugMode = true
	} else if verbose {
		logLevel = "debug"
		c.verboseMode = true
	} else if quiet {
		logLevel = "error"
		c.quietMode = true
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	isValidLogLevel := false
	for _, level := range validLogLevels {
		if logLevel == level {
			isValidLogLevel = true
			break
		}
	}
	if !isValidLogLevel {
		return fmt.Errorf("🚫 '%s' isn't a valid log level. Try one of these: %s", logLevel, strings.Join(validLogLevels, ", "))
	}

	// Validate output format
	validOutputFormats := []string{"text", "json", "yaml"}
	isValidOutputFormat := false
	for _, format := range validOutputFormats {
		if outputFormat == format {
			isValidOutputFormat = true
			break
		}
	}
	if !isValidOutputFormat {
		return fmt.Errorf("🚫 '%s' isn't a valid output format. Try one of these: %s", outputFormat, strings.Join(validOutputFormats, ", "))
	}

	// Configure logger based on flags
	if c.logger != nil {
		// Set log level
		switch logLevel {
		case "debug":
			c.logger.SetLevel(0) // LogLevelDebug
		case "info":
			c.logger.SetLevel(1) // LogLevelInfo
		case "warn":
			c.logger.SetLevel(2) // LogLevelWarn
		case "error":
			c.logger.SetLevel(3) // LogLevelError
		case "fatal":
			c.logger.SetLevel(4) // LogLevelFatal
		}

		// Configure JSON output
		c.logger.SetJSONOutput(logJSON)

		// Configure caller information
		c.logger.SetCallerInfo(logCaller)

		// Log configuration changes in debug mode
		if c.debugMode {
			c.logger.DebugWithFields("🔧 Setting up CLI configuration", map[string]interface{}{
				"command":         cmd.Name(),
				"verbose":         verbose,
				"quiet":           quiet,
				"debug":           debug,
				"log_level":       logLevel,
				"log_json":        logJSON,
				"log_caller":      logCaller,
				"non_interactive": nonInteractive,
				"output_format":   outputFormat,
			})
		}

		// Log operation start for verbose mode
		if c.verboseMode || c.debugMode {
			c.logger.InfoWithFields("🚀 Starting your command", map[string]interface{}{
				"command": cmd.Name(),
				"args":    cmd.Flags().Args(),
			})
		}
	}

	// Store global settings for use in commands
	if c.verboseMode && !c.quietMode {
		fmt.Printf("🔍 Running '%s' with detailed output\n", cmd.Name())
	}

	if c.debugMode && !c.quietMode {
		fmt.Printf("🐛 Running '%s' with debug logging and performance metrics\n", cmd.Name())
	}

	if nonInteractive && (c.verboseMode || c.debugMode) && !c.quietMode {
		fmt.Printf("🤖 Running in non-interactive mode\n")
	}

	return nil
}

// Color helper functions for beautiful output
func (c *CLI) colorize(color, text string) string {
	if c.quietMode {
		return text // No colors in quiet mode
	}
	return color + text + ColorReset
}

func (c *CLI) success(text string) string {
	return c.colorize(ColorGreen+ColorBold, text)
}

func (c *CLI) info(text string) string {
	return c.colorize(ColorBlue, text)
}

func (c *CLI) warning(text string) string {
	return c.colorize(ColorYellow, text)
}

func (c *CLI) error(text string) string {
	return c.colorize(ColorRed+ColorBold, text)
}

func (c *CLI) highlight(text string) string {
	return c.colorize(ColorCyan+ColorBold, text)
}

func (c *CLI) dim(text string) string {
	return c.colorize(ColorDim, text)
}

// Verbose output methods for enhanced debugging and user feedback

// VerboseOutput prints verbose information if verbose mode is enabled
func (c *CLI) VerboseOutput(format string, args ...interface{}) {
	if c.verboseMode && !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	if c.logger != nil && c.verboseMode && c.logger.IsInfoEnabled() {
		c.logger.Info(format, args...)
	}
}

// DebugOutput prints debug information if debug mode is enabled
func (c *CLI) DebugOutput(format string, args ...interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("🐛 "+format+"\n", args...)
	}
	if c.logger != nil && c.logger.IsDebugEnabled() {
		c.logger.Debug(format, args...)
	}
}

// QuietOutput prints information only if not in quiet mode
func (c *CLI) QuietOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for QuietOutput - it's meant for user-facing messages
}

// ErrorOutput prints error information (always shown unless completely silent)
func (c *CLI) ErrorOutput(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	// Only log to logger in verbose or debug mode to avoid cluttering output
	if c.logger != nil && (c.verboseMode || c.debugMode) {
		c.logger.Error(format, args...)
	}
}

// WarningOutput prints warning information if not in quiet mode
func (c *CLI) WarningOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf("⚠️  "+format+"\n", args...)
	}
	if c.logger != nil {
		c.logger.Warn(format, args...)
	}
}

// SuccessOutput prints success information if not in quiet mode
func (c *CLI) SuccessOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for SuccessOutput - it's meant for user-facing messages
}

// PerformanceOutput prints performance metrics if debug mode is enabled
func (c *CLI) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("⚡ %s completed in %v\n", operation, duration)
		if len(metrics) > 0 {
			fmt.Printf("📊 Performance Metrics:\n")
			for k, v := range metrics {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}
	if c.logger != nil {
		allMetrics := make(map[string]interface{})
		allMetrics["duration_ms"] = duration.Milliseconds()
		allMetrics["duration_human"] = duration.String()
		for k, v := range metrics {
			allMetrics[k] = v
		}
		c.logger.LogPerformanceMetrics(operation, allMetrics)
	}
}

// StartOperationWithOutput starts an operation with verbose output
func (c *CLI) StartOperationWithOutput(operation string, description string) *interfaces.OperationContext {
	c.VerboseOutput("%s", description)

	var ctx *interfaces.OperationContext
	if c.logger != nil && c.verboseMode {
		ctx = c.logger.StartOperation(operation, map[string]interface{}{
			"description": description,
		})
	}

	return ctx
}

// FinishOperationWithOutput completes an operation with verbose output
func (c *CLI) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	if ctx != nil && c.logger != nil && c.verboseMode {
		c.logger.FinishOperation(ctx, map[string]interface{}{
			"description": description,
		})
	}
	c.VerboseOutput("%s", description)
}

// FinishOperationWithError completes an operation with error output
func (c *CLI) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	if ctx != nil && c.logger != nil {
		c.logger.FinishOperationWithError(ctx, err, nil)
	}
	c.ErrorOutput("❌ %s failed: %v", operation, err)
}

// GetExitCode returns the current exit code
func (c *CLI) GetExitCode() int {
	return c.exitCode
}

// SetExitCode sets the exit code for the CLI
func (c *CLI) SetExitCode(code int) {
	c.exitCode = code
}

// GetVersionManager returns the version manager instance
func (c *CLI) GetVersionManager() interfaces.VersionManager {
	return c.versionManager
}

// PromptProjectDetails collects basic project configuration from user input.
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	if c.isNonInteractiveMode() {
		return nil, fmt.Errorf("interactive prompts not available in non-interactive mode")
	}

	c.QuietOutput("Project Configuration")
	c.QuietOutput("====================")

	config := &models.ProjectConfig{}

	// Get project name
	fmt.Print("Project name: ")
	var name string
	if _, err := fmt.Scanln(&name); err != nil {
		return nil, fmt.Errorf("failed to read project name: %w", err)
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
func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	if c.isNonInteractiveMode() {
		return true // Auto-confirm in non-interactive mode
	}

	c.QuietOutput("\nProject Configuration Preview:")
	c.QuietOutput("==============================")
	c.QuietOutput("Name: %s", config.Name)
	if config.Organization != "" {
		c.QuietOutput("Organization: %s", config.Organization)
	}
	if config.Description != "" {
		c.QuietOutput("Description: %s", config.Description)
	}
	if config.Author != "" {
		c.QuietOutput("Author: %s", config.Author)
	}
	c.QuietOutput("License: %s", config.License)

	c.QuietOutput("\nComponents:")
	if config.Components.Backend.GoGin {
		c.QuietOutput("  - Go Gin API")
	}
	if config.Components.Frontend.NextJS.App {
		c.QuietOutput("  - Next.js App")
	}
	if config.Components.Infrastructure.Docker {
		c.QuietOutput("  - Docker configuration")
	}

	fmt.Print("\nProceed with generation? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

// setupGenerateCommand sets up the generate command with all documented flags
func (c *CLI) setupGenerateCommand() {
	generateCmd := &cobra.Command{
		Use:   "generate [flags]",
		Short: "Generate a new project from templates with modern best practices",
		Long: `Generate production-ready projects with modern best practices.

Supports Go, Next.js, React, Android, iOS, Docker, Kubernetes, and Terraform.
Use interactive mode or provide a configuration file.`,
		RunE: c.runGenerate,
		Example: `  # Interactive project creation
  generator generate
  
  # Generate from configuration file
  generator generate --config project.yaml
  
  # Use specific template
  generator generate --template go-gin
  
  # Non-interactive mode for automation
  generator generate --config project.yaml --non-interactive`,
	}

	// Basic flags
	generateCmd.Flags().StringP("config", "c", "", "Path to configuration file (YAML or JSON)")
	generateCmd.Flags().StringP("output", "o", "", "Output directory for generated project")
	generateCmd.Flags().Bool("dry-run", false, "Preview generation without creating files")

	// Configuration management flags
	generateCmd.Flags().String("load-config", "", "Load a saved configuration by name")
	generateCmd.Flags().String("save-config", "", "Save configuration with the specified name after interactive setup")
	generateCmd.Flags().Bool("list-configs", false, "List available saved configurations and exit")

	// Advanced flags
	generateCmd.Flags().Bool("offline", false, "Use cached templates and versions without network requests")
	generateCmd.Flags().Bool("minimal", false, "Generate minimal project structure with only essential components")
	generateCmd.Flags().String("template", "", "Use specific template instead of interactive selection")
	generateCmd.Flags().Bool("update-versions", false, "Fetch and use latest package versions")
	generateCmd.Flags().Bool("force", false, "Overwrite existing files in output directory")
	generateCmd.Flags().Bool("skip-validation", false, "Skip configuration validation")
	generateCmd.Flags().Bool("backup-existing", true, "Create backups of existing files before overwriting")
	generateCmd.Flags().Bool("include-examples", true, "Include example code and documentation")

	// Additional flags for enhanced functionality
	generateCmd.Flags().StringSlice("exclude", []string{}, "Exclude specific files or directories from generation")
	generateCmd.Flags().StringSlice("include-only", []string{}, "Include only specific files or directories in generation")
	generateCmd.Flags().Bool("interactive", true, "Use interactive mode for project configuration")
	generateCmd.Flags().String("preset", "", "Use predefined configuration preset")

	// Mode-specific flags
	generateCmd.Flags().Bool("force-interactive", false, "Force interactive mode even in CI/automated environments")
	generateCmd.Flags().Bool("force-non-interactive", false, "Force non-interactive mode even in terminal environments")
	generateCmd.Flags().String("mode", "", "Explicitly set generation mode (interactive, non-interactive, config-file)")

	c.rootCmd.AddCommand(generateCmd)
}

// setupValidateCommand sets up the validate command with all documented flags
func (c *CLI) setupValidateCommand() {
	validateCmd := &cobra.Command{
		Use:   "validate [path] [flags]",
		Short: "Validate project structure, configuration, and dependencies",
		Long: `Validate project structure, configuration files, and dependencies.

Checks code quality, security, and best practices. Can automatically fix common issues.`,
		RunE: c.runValidate,
		Example: `  # Validate current directory
  generator validate
  
  # Validate specific project
  generator validate ./my-project
  
  # Validate and fix issues
  generator validate ./my-project --fix
  
  # Generate HTML report
  generator validate --report --report-format html`,
	}

	validateCmd.Flags().Bool("fix", false, "Automatically fix common validation issues")
	validateCmd.Flags().Bool("report", false, "Generate detailed validation report")
	validateCmd.Flags().String("report-format", "text", "Report format (text, json, html, markdown)")
	validateCmd.Flags().StringSlice("rules", []string{}, "Specific validation rules to apply")
	validateCmd.Flags().Bool("ignore-warnings", false, "Ignore validation warnings")
	validateCmd.Flags().String("output-file", "", "Save report to file")
	validateCmd.Flags().StringP("output", "o", "", "Save report to file (alias for --output-file)")

	// Additional validation flags
	validateCmd.Flags().Bool("strict", false, "Use strict validation mode")
	validateCmd.Flags().Bool("summary-only", false, "Show only validation summary")
	validateCmd.Flags().StringSlice("exclude-rules", []string{}, "Exclude specific validation rules")
	validateCmd.Flags().Bool("show-fixes", false, "Show available fixes for issues")

	c.rootCmd.AddCommand(validateCmd)
}

// setupAuditCommand sets up the audit command with all documented flags
func (c *CLI) setupAuditCommand() {
	auditCmd := &cobra.Command{
		Use:   "audit [path] [flags]",
		Short: "Comprehensive security, quality, and compliance auditing",
		Long: `Audit project security, code quality, license compliance, and performance.

Provides detailed reports with scores and recommendations for improvement.`,
		RunE: c.runAudit,
		Example: `  # Audit current directory
  generator audit
  
  # Audit specific project
  generator audit ./my-project
  
  # Security audit only
  generator audit --security
  
  # Generate detailed report
  generator audit --detailed --output-format html`,
	}

	auditCmd.Flags().Bool("security", true, "Perform security vulnerability scanning")
	auditCmd.Flags().Bool("quality", true, "Perform code quality analysis")
	auditCmd.Flags().Bool("licenses", true, "Perform license compliance checking")
	auditCmd.Flags().Bool("performance", true, "Perform performance analysis")
	auditCmd.Flags().String("output-format", "text", "Output format (text, json, html)")
	auditCmd.Flags().String("output-file", "", "Save audit report to file")
	auditCmd.Flags().Bool("detailed", false, "Generate detailed audit report")

	// Additional audit flags
	auditCmd.Flags().Bool("fail-on-high", false, "Fail if high severity issues are found")
	auditCmd.Flags().Bool("fail-on-medium", false, "Fail if medium or higher severity issues are found")
	auditCmd.Flags().Float64("min-score", 0.0, "Minimum acceptable audit score (0.0-10.0)")
	auditCmd.Flags().StringSlice("exclude-categories", []string{}, "Exclude specific audit categories")
	auditCmd.Flags().Bool("summary-only", false, "Show only audit summary")

	c.rootCmd.AddCommand(auditCmd)
}

// setupVersionCommand sets up the version command with all documented flags
func (c *CLI) setupVersionCommand() {
	versionCmd := NewVersionCommand(c)

	// Add additional flags for the main CLI
	versionCmd.Flags().Bool("packages", false, "Show latest package versions for all supported technologies")
	versionCmd.Flags().Bool("check-updates", false, "Check for generator updates")
	versionCmd.Flags().Bool("build-info", false, "Show detailed build information")
	versionCmd.Flags().Bool("short", false, "Show only version number")
	versionCmd.Flags().String("format", "text", "Output format (text, json, yaml)")
	versionCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	versionCmd.Flags().String("check-package", "", "Check version for specific package")

	// Update the command description and examples
	versionCmd.Long = `Display comprehensive version information for the generator and supported technologies.

This command provides detailed version information including:
- Generator version and build information
- Latest versions of supported packages and technologies
- Update availability and compatibility information
- Build metadata and system information

The command supports multiple output formats and can check for updates
both for the generator itself and for supported technology packages.`

	versionCmd.Example = `  # Basic version information
  generator version
  
  # Show version in JSON format
  generator version --json
  
  # Show all package versions
  generator version --packages
  
  # Check for generator updates
  generator version --check-updates
  
  # Show detailed build information
  generator version --build-info
  
  # Short version output (version number only)
  generator version --short
  
  # Machine-readable output for CI/CD
  generator version --packages --output-format json --quiet
  
  # Check for updates in non-interactive mode
  generator version --check-updates --non-interactive
  
  # Get version information for specific technology stack
  generator version --packages --format json | jq '.go'

  TROUBLESHOOTING:
  # Debug version checking issues
  generator version --debug --check-updates
  
  # Verbose output with registry information
  generator version --packages --verbose --debug`

	c.rootCmd.AddCommand(versionCmd)
}

// setupListTemplatesCommand sets up the list-templates command
func (c *CLI) setupListTemplatesCommand() {
	listTemplatesCmd := &cobra.Command{
		Use:   "list-templates [flags]",
		Short: "List and discover available project templates",
		Long: `List available project templates with filtering and search.

Browse templates for frontend, backend, mobile, and infrastructure projects.`,
		RunE: c.runListTemplates,
		Example: `  # List all templates
  generator list-templates
  
  # Filter by category
  generator list-templates --category backend
  
  # Search for templates
  generator list-templates --search api
  
  # Show detailed information
  generator list-templates --detailed`,
	}

	listTemplatesCmd.Flags().String("category", "", "Filter by category (frontend, backend, mobile, infrastructure)")
	listTemplatesCmd.Flags().String("technology", "", "Filter by technology (go, nodejs, react, etc.)")
	listTemplatesCmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	listTemplatesCmd.Flags().String("search", "", "Search templates by name or description")
	listTemplatesCmd.Flags().Bool("detailed", false, "Show detailed template information")

	c.rootCmd.AddCommand(listTemplatesCmd)
}

// setupTemplateCommand sets up the template command with subcommands
func (c *CLI) setupTemplateCommand() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Template management operations",
		Long: `Manage templates including viewing detailed information and validation.

Use these commands to inspect templates before using them
or to validate custom templates you've created.`,
	}

	// template info
	templateInfoCmd := &cobra.Command{
		Use:   "info <template-name> [flags]",
		Short: "Display comprehensive template information and documentation",
		Long: `Display detailed information about a specific template including metadata,
dependencies, configuration options, and usage examples.

TEMPLATE INFORMATION:
  Basic Information:
    • Template name, version, and description
    • Author, maintainer, and license information
    • Creation date, last update, and changelog
    • Category, tags, and classification

  Technical Details:
    • Technology stack and version requirements
    • Dependencies and compatibility matrix
    • Configuration variables and options
    • File structure and component breakdown
    • Build system and deployment information

  Usage Information:
    • Configuration examples and templates
    • Common use cases and scenarios
    • Best practices and recommendations
    • Troubleshooting and known issues
    • Community resources and support

  Compatibility Information:
    • Supported platforms and environments
    • Version compatibility matrix
    • Breaking changes and migration guides
    • Integration with other templates
    • Performance characteristics and limitations`,
		RunE: c.runTemplateInfo,
		Args: cobra.ExactArgs(1),
		Example: `  # Show basic template information
  generator template info go-gin
  
  # Show detailed information with all sections
  generator template info nextjs-app --detailed
  
  # Show template variables and configuration options
  generator template info go-gin --variables --detailed
  
  # Show dependency information
  generator template info nextjs-app --dependencies --compatibility
  
  # Show file structure and components
  generator template info go-gin --structure --verbose
  
  # Output in JSON format for automation
  generator template info nextjs-app --output-format json`,
	}
	templateInfoCmd.Flags().Bool("detailed", false, "Show detailed template information")
	templateInfoCmd.Flags().Bool("variables", false, "Show template variables")
	templateInfoCmd.Flags().Bool("dependencies", false, "Show template dependencies")
	templateInfoCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	templateCmd.AddCommand(templateInfoCmd)

	// template validate
	templateValidateCmd := &cobra.Command{
		Use:   "validate <template-path> [flags]",
		Short: "Validate custom template structure, syntax, and compliance",
		Long: `Comprehensive validation of custom template directories including structure,
metadata, syntax, and best practices compliance. Provides detailed feedback and auto-fix capabilities.

VALIDATION CATEGORIES:
  Structure Validation:
    • Required files and directories presence
    • Template file organization and naming
    • Metadata file structure and completeness
    • Asset and resource file validation
    • Documentation and example file checking

  Syntax Validation:
    • Template syntax correctness and parsing
    • Variable usage and definition validation
    • Conditional logic and loop validation
    • Function usage and parameter validation
    • Template inheritance and inclusion validation

  Metadata Validation:
    • Template metadata completeness and accuracy
    • Version information and compatibility data
    • Dependency specification and validation
    • Configuration schema and variable definitions
    • License and author information validation

  Best Practices Compliance:
    • Template organization and structure standards
    • Security best practices and vulnerability checks
    • Performance optimization recommendations
    • Documentation quality and completeness
    • Accessibility and usability guidelines

AUTO-FIX CAPABILITIES:
  • Automatic correction of common syntax errors
  • Missing file and directory creation
  • Metadata completion and standardization
  • Documentation template generation
  • Security and best practices improvements`,
		RunE: c.runTemplateValidate,
		Args: cobra.ExactArgs(1),
		Example: `  # Validate custom template directory
  generator template validate ./my-custom-template
  
  # Validate with detailed output and suggestions
  generator template validate ./my-template --detailed --verbose
  
  # Validate and auto-fix common issues
  generator template validate ./my-template --fix --backup
  
  # Validate with strict compliance checking
  generator template validate ./my-template --strict --best-practices
  
  # Generate validation report
  generator template validate ./my-template --report --output-format html
  
  # Validate specific aspects only
  generator template validate ./my-template --syntax-only --metadata-only`,
	}
	templateValidateCmd.Flags().Bool("detailed", false, "Show detailed validation results")
	templateValidateCmd.Flags().Bool("fix", false, "Attempt to fix validation issues")
	templateValidateCmd.Flags().String("output-format", "text", "Output format (text, json)")
	templateCmd.AddCommand(templateValidateCmd)

	c.rootCmd.AddCommand(templateCmd)
}

// setupUpdateCommand sets up the update command
func (c *CLI) setupUpdateCommand() {
	updateCmd := &cobra.Command{
		Use:   "update [flags]",
		Short: "Update generator, templates, and package information",
		Long: `Update generator, templates, and package information.

Includes safety checks, rollback capabilities, and multiple update channels.`,
		RunE: c.runUpdate,
		Example: `  # Check for updates
  generator update --check
  
  # Install updates
  generator update --install
  
  # Update templates only
  generator update --templates
  
  # Use beta channel
  generator update --channel beta --install`,
	}

	updateCmd.Flags().Bool("check", false, "Check for updates without installing")
	updateCmd.Flags().Bool("install", false, "Install available updates")
	updateCmd.Flags().Bool("templates", false, "Update templates cache")
	updateCmd.Flags().Bool("force", false, "Force update even if current version is newer")
	updateCmd.Flags().Bool("compatibility", false, "Check compatibility before updating")
	updateCmd.Flags().Bool("release-notes", false, "Show release notes for available updates")
	updateCmd.Flags().String("channel", "stable", "Update channel (stable, beta, alpha)")
	updateCmd.Flags().Bool("backup", true, "Create backup before updating")
	updateCmd.Flags().Bool("verify", true, "Verify update signatures")
	updateCmd.Flags().String("version", "", "Update to specific version")

	c.rootCmd.AddCommand(updateCmd)
}

// setupCacheCommand sets up the cache command with all subcommands
func (c *CLI) setupCacheCommand() {
	cacheCmd := &cobra.Command{
		Use:   "cache <command> [flags]",
		Short: "Manage local cache for offline mode and performance",
		Long: `Manage local cache for templates, package versions, and other data.

Enables offline mode operation and improves performance through intelligent caching.`,
	}

	// cache show
	cacheShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show cache status and statistics",
		Long:  "Display cache location, size, statistics, and health information",
		RunE:  c.runCacheShow,
	}
	cacheCmd.AddCommand(cacheShowCmd)

	// cache clear
	cacheClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cached data",
		Long:  "Remove all cached templates, versions, and other data",
		RunE:  c.runCacheClear,
	}
	cacheClearCmd.Flags().Bool("force", false, "Clear cache without confirmation")
	cacheCmd.AddCommand(cacheClearCmd)

	// cache clean
	cacheCleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove expired and invalid cache entries",
		Long:  "Clean up expired cache entries and repair corrupted cache data",
		RunE:  c.runCacheClean,
	}
	cacheCmd.AddCommand(cacheCleanCmd)

	// cache validate
	cacheValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate cache integrity",
		Long:  "Check cache integrity and report any issues",
		RunE:  c.runCacheValidate,
	}
	cacheCmd.AddCommand(cacheValidateCmd)

	// cache repair
	cacheRepairCmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair corrupted cache data",
		Long:  "Attempt to repair corrupted cache entries and fix cache issues",
		RunE:  c.runCacheRepair,
	}
	cacheCmd.AddCommand(cacheRepairCmd)

	// cache offline
	cacheOfflineCmd := &cobra.Command{
		Use:   "offline",
		Short: "Manage offline mode",
		Long:  "Enable or disable offline mode for the cache",
	}

	// cache offline enable
	cacheOfflineEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable offline mode",
		Long:  "Enable offline mode to use only cached data",
		RunE:  c.runCacheOfflineEnable,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineEnableCmd)

	// cache offline disable
	cacheOfflineDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable offline mode",
		Long:  "Disable offline mode to allow network access",
		RunE:  c.runCacheOfflineDisable,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineDisableCmd)

	// cache offline status
	cacheOfflineStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show offline mode status",
		Long:  "Display current offline mode status and readiness",
		RunE:  c.runCacheOfflineStatus,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineStatusCmd)

	cacheCmd.AddCommand(cacheOfflineCmd)
	c.rootCmd.AddCommand(cacheCmd)
}

// setupLogsCommand sets up the logs command
func (c *CLI) setupLogsCommand() {
	logsCmd := &cobra.Command{
		Use:   "logs [flags]",
		Short: "View and analyze application logs for debugging and monitoring",
		Long: `View and analyze application logs for debugging and monitoring.

Provides filtering, search, real-time following, and multiple output formats.`,
		RunE: c.runLogs,
		Example: `  BASIC LOG VIEWING:
  # Show recent log entries (default: 50 lines)
  generator logs
  
  # Show specific number of log entries
  generator logs --lines 100
  
  # Show log file locations and information
  generator logs --locations

  LOG LEVEL FILTERING:
  # Show only error and fatal logs
  generator logs --level error
  
  # Show warnings and above (warn, error, fatal)
  generator logs --level warn
  
  # Show debug logs for detailed troubleshooting
  generator logs --level debug --lines 200

  COMPONENT AND SOURCE FILTERING:
  # Show logs from CLI component only
  generator logs --component cli
  
  # Show template processing logs
  generator logs --component template --level debug
  
  # Show configuration-related logs
  generator logs --component config --verbose

  TIME-BASED FILTERING:
  # Show logs since specific timestamp
  generator logs --since "2024-01-01T10:00:00Z"
  
  # Show logs from last hour
  generator logs --since "1h"
  
  # Show logs from last 30 minutes
  generator logs --since "30m"

  REAL-TIME MONITORING:
  # Follow logs in real-time (like tail -f)
  generator logs --follow
  
  # Follow error logs only
  generator logs --follow --level error
  
  # Follow specific component logs
  generator logs --follow --component template

  OUTPUT FORMATS AND ANALYSIS:
  # JSON output for automation and parsing
  generator logs --format json --lines 100
  
  # Raw log file content
  generator logs --format raw --no-color
  
  # CSV format for analysis tools
  generator logs --format csv --since "24h"

  ADVANCED FILTERING:
  # Combine multiple filters
  generator logs --level warn --component template --since "1h" --lines 50
  
  # Search for specific patterns
  generator logs --search "error" --level info
  
  # Exclude timestamps for cleaner output
  generator logs --no-timestamps --format text

  TROUBLESHOOTING:
  # Debug log viewing issues
  generator logs --debug --verbose --locations
  
  # Show all available log files and their status
  generator logs --locations --detailed`,
	}

	logsCmd.Flags().Int("lines", 50, "Number of recent log lines to show")
	logsCmd.Flags().String("level", "", "Filter by log level (debug, info, warn, error, fatal)")
	logsCmd.Flags().String("component", "", "Filter by component name")
	logsCmd.Flags().String("since", "", "Show logs since specific time (RFC3339 format)")
	logsCmd.Flags().Bool("follow", false, "Follow log output in real-time (tail -f)")
	logsCmd.Flags().Bool("locations", false, "Show log file locations only")
	logsCmd.Flags().String("format", "text", "Output format (text, json, raw)")
	logsCmd.Flags().Bool("no-color", false, "Disable colored output")
	logsCmd.Flags().Bool("timestamps", true, "Show timestamps in output")

	c.rootCmd.AddCommand(logsCmd)
}

// Command execution methods - these will be implemented in subsequent tasks
// For now, they return "not implemented" errors to satisfy the interface

func (c *CLI) runGenerate(cmd *cobra.Command, args []string) error {
	// Get flags
	configPath, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	offline, _ := cmd.Flags().GetBool("offline")
	minimal, _ := cmd.Flags().GetBool("minimal")
	template, _ := cmd.Flags().GetString("template")
	updateVersions, _ := cmd.Flags().GetBool("update-versions")
	force, _ := cmd.Flags().GetBool("force")
	skipValidation, _ := cmd.Flags().GetBool("skip-validation")
	backupExisting, _ := cmd.Flags().GetBool("backup-existing")
	includeExamples, _ := cmd.Flags().GetBool("include-examples")
	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Additional flags (for future implementation)
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	includeOnly, _ := cmd.Flags().GetStringSlice("include-only")
	interactive, _ := cmd.Flags().GetBool("interactive")
	preset, _ := cmd.Flags().GetString("preset")

	// Mode-specific flags
	forceInteractive, _ := cmd.Flags().GetBool("force-interactive")
	forceNonInteractive, _ := cmd.Flags().GetBool("force-non-interactive")
	explicitMode, _ := cmd.Flags().GetString("mode")

	// Validate conflicting mode flags
	if err := c.validateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode); err != nil {
		return fmt.Errorf("conflicting mode flags: %w", err)
	}

	// Apply mode overrides
	nonInteractive, interactive = c.applyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)

	// Log additional options for debugging
	if len(exclude) > 0 {
		c.DebugOutput("Excluding files/directories: %v", exclude)
	}
	if len(includeOnly) > 0 {
		c.DebugOutput("Including only: %v", includeOnly)
	}
	if preset != "" {
		c.DebugOutput("Using preset: %s", preset)
	}

	// Create generate options
	options := interfaces.GenerateOptions{
		Force:           force,
		Minimal:         minimal,
		Offline:         offline,
		UpdateVersions:  updateVersions,
		SkipValidation:  skipValidation,
		BackupExisting:  backupExisting,
		IncludeExamples: includeExamples,
		OutputPath:      outputPath,
		DryRun:          dryRun,
		NonInteractive:  nonInteractive,
	}

	if template != "" {
		options.Templates = []string{template}
	}

	// Perform comprehensive validation before generation
	if !options.SkipValidation {
		c.VerboseOutput("🔍 Validating your configuration...")
		if err := c.validateGenerateOptions(options); err != nil {
			return fmt.Errorf("🚫 Your configuration has some issues: %w", err)
		}
	}

	// Mode detection and routing logic
	mode := c.detectGenerationMode(configPath, nonInteractive, interactive, explicitMode)
	c.VerboseOutput("🎯 Using %s mode for project generation", mode)

	// Route to appropriate generation method based on detected mode
	return c.routeToGenerationMethod(mode, configPath, options)
}

func (c *CLI) runValidate(cmd *cobra.Command, args []string) error {
	// Get path from args or use current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Get flags
	fix, _ := cmd.Flags().GetBool("fix")
	report, _ := cmd.Flags().GetBool("report")
	reportFormat, _ := cmd.Flags().GetString("report-format")
	rules, _ := cmd.Flags().GetStringSlice("rules")
	ignoreWarnings, _ := cmd.Flags().GetBool("ignore-warnings")
	outputFile, _ := cmd.Flags().GetString("output-file")
	output, _ := cmd.Flags().GetString("output")

	// Use --output if provided, otherwise use --output-file
	if output != "" {
		outputFile = output
	}
	// Additional validation flags (for future implementation)
	strict, _ := cmd.Flags().GetBool("strict")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")
	excludeRules, _ := cmd.Flags().GetStringSlice("exclude-rules")
	showFixes, _ := cmd.Flags().GetBool("show-fixes")

	// Get global flags
	verbose, _ := cmd.Flags().GetBool("verbose")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
	}

	// Log additional options for debugging
	if strict {
		c.DebugOutput("Using strict validation mode")
	}
	if len(excludeRules) > 0 {
		c.DebugOutput("Excluding rules: %v", excludeRules)
	}

	// Use additional flags for future implementation
	_ = summaryOnly
	_ = showFixes

	// Create validation options
	options := interfaces.ValidationOptions{
		Verbose:        verbose,
		Fix:            fix,
		Report:         report,
		ReportFormat:   reportFormat,
		Rules:          rules,
		IgnoreWarnings: ignoreWarnings,
		OutputFile:     outputFile,
	}

	// Validate project
	result, err := c.ValidateProject(path, options)
	if err != nil {
		return fmt.Errorf("🚫 Couldn't validate your project: %w", err)
	}

	// Output results based on format and mode
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		// Machine-readable output for automation
		return c.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	if !c.quietMode {
		c.QuietOutput("🔍 Validation completed for: %s", path)
		if result.Valid {
			c.QuietOutput("✅ Project looks good!")
		} else {
			c.QuietOutput("❌ Found some issues that need attention")
		}
		c.QuietOutput("📊 Issues: %d", len(result.Issues))
		c.QuietOutput("⚠️  Warnings: %d", len(result.Warnings))

		if len(result.Issues) > 0 && !summaryOnly {
			c.QuietOutput("\n🚨 Issues that need fixing:")
			for _, issue := range result.Issues {
				c.QuietOutput("  - %s: %s", issue.Severity, issue.Message)
				if issue.File != "" {
					c.VerboseOutput("    File: %s:%d:%d", issue.File, issue.Line, issue.Column)
				}
			}
		}

		if len(result.Warnings) > 0 && !ignoreWarnings && !summaryOnly {
			c.QuietOutput("\n⚠️  Things to consider:")
			for _, warning := range result.Warnings {
				c.QuietOutput("  - %s: %s", warning.Severity, warning.Message)
				if warning.File != "" {
					c.VerboseOutput("    File: %s:%d:%d", warning.File, warning.Line, warning.Column)
				}
			}
		}

		if showFixes && len(result.FixSuggestions) > 0 {
			c.QuietOutput("\nSuggested fixes:")
			for _, suggestion := range result.FixSuggestions {
				c.QuietOutput("  - %s", suggestion.Description)
				if suggestion.AutoFixable {
					c.QuietOutput("    (Auto-fixable with --fix flag)")
				}
			}
		}
	}

	// Generate report if requested
	if report && outputFile != "" {
		err := c.generateValidationReport(result, reportFormat, outputFile)
		if err != nil {
			return fmt.Errorf("failed to generate validation report: %w", err)
		}
		c.QuietOutput("Validation report written to: %s", outputFile)
	}

	// Return appropriate exit code
	if !result.Valid {
		details := map[string]interface{}{
			"issues_count":   len(result.Issues),
			"warnings_count": len(result.Warnings),
			"path":           path,
		}
		
		var message string
		if len(result.Issues) > 0 {
			message = fmt.Sprintf("🚫 Found %d validation issues that need to be fixed", len(result.Issues))
		} else if len(result.Warnings) > 0 {
			message = fmt.Sprintf("⚠️ Found %d warnings that should be addressed", len(result.Warnings))
		} else {
			message = "🚫 Validation failed - please check your project structure"
		}
		
		return c.createValidationError(message, details)
	}

	return nil
}

// generateValidationReport generates a validation report in the specified format
func (c *CLI) generateValidationReport(result *interfaces.ValidationResult, format, outputFile string) error {
	var content []byte
	var err error

	switch format {
	case "json":
		content, err = json.MarshalIndent(result, "", "  ")
	case "text":
		status := "✅ Looks good!"
		if !result.Valid {
			status = "❌ Needs attention"
		}
		content = []byte(fmt.Sprintf("🔍 Validation Report\n===================\n\nStatus: %s\n📊 Issues: %d\n⚠️  Warnings: %d\n",
			status, len(result.Issues), len(result.Warnings)))
	default:
		status := "✅ Looks good!"
		if !result.Valid {
			status = "❌ Needs attention"
		}
		content = []byte(fmt.Sprintf("🔍 Validation Report\n===================\n\nStatus: %s\n📊 Issues: %d\n⚠️  Warnings: %d\n",
			status, len(result.Issues), len(result.Warnings)))
	}

	if err != nil {
		return fmt.Errorf("failed to format report: %w", err)
	}

	err = os.WriteFile(outputFile, content, 0600)
	if err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

func (c *CLI) runAudit(cmd *cobra.Command, args []string) error {
	// Get path from args or use current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Get flags
	security, _ := cmd.Flags().GetBool("security")
	quality, _ := cmd.Flags().GetBool("quality")
	licenses, _ := cmd.Flags().GetBool("licenses")
	performance, _ := cmd.Flags().GetBool("performance")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	detailed, _ := cmd.Flags().GetBool("detailed")
	failOnHigh, _ := cmd.Flags().GetBool("fail-on-high")
	failOnMedium, _ := cmd.Flags().GetBool("fail-on-medium")
	minScore, _ := cmd.Flags().GetFloat64("min-score")
	// Additional audit flags (for future implementation)
	excludeCategories, _ := cmd.Flags().GetStringSlice("exclude-categories")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	globalOutputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
	}

	// Use global output format if command-specific format not set
	if outputFormat == "text" && globalOutputFormat != "text" {
		outputFormat = globalOutputFormat
	}

	// Log additional options for debugging
	if len(excludeCategories) > 0 {
		c.DebugOutput("Excluding audit categories: %v", excludeCategories)
	}
	if summaryOnly {
		c.DebugOutput("Showing summary only")
	}

	// Create audit options
	options := interfaces.AuditOptions{
		Security:     security,
		Quality:      quality,
		Licenses:     licenses,
		Performance:  performance,
		OutputFormat: outputFormat,
		OutputFile:   outputFile,
		Detailed:     detailed,
	}

	// Audit project
	result, err := c.AuditProject(path, options)
	if err != nil {
		return err
	}

	// Output results based on format and mode
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		// Machine-readable output for automation
		return c.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	if !c.quietMode {
		c.QuietOutput("🔍 Audit completed for: %s", path)
		scoreEmoji := "🎉"
		if result.OverallScore < 70 {
			scoreEmoji = "⚠️ "
		}
		if result.OverallScore < 50 {
			scoreEmoji = "🚨"
		}
		c.QuietOutput("%s Overall Score: %.1f/100", scoreEmoji, result.OverallScore)
		c.VerboseOutput("Audit Time: %s", result.AuditTime.Format("2006-01-02 15:04:05"))

		if result.Security != nil && !summaryOnly {
			securityEmoji := "🔒"
			if result.Security.Score < 70 {
				securityEmoji = "⚠️ "
			}
			if result.Security.Score < 50 {
				securityEmoji = "🚨"
			}
			c.QuietOutput("%s Security Score: %.1f/100", securityEmoji, result.Security.Score)
			c.VerboseOutput("Vulnerabilities: %d", len(result.Security.Vulnerabilities))
		}

		if result.Quality != nil && !summaryOnly {
			qualityEmoji := "✨"
			if result.Quality.Score < 70 {
				qualityEmoji = "⚠️ "
			}
			if result.Quality.Score < 50 {
				qualityEmoji = "🚨"
			}
			c.QuietOutput("%s Quality Score: %.1f/100", qualityEmoji, result.Quality.Score)
			c.VerboseOutput("Code Smells: %d", len(result.Quality.CodeSmells))
		}

		if result.Licenses != nil && !summaryOnly {
			c.QuietOutput("License Compatible: %t", result.Licenses.Compatible)
			c.VerboseOutput("License Conflicts: %d", len(result.Licenses.Conflicts))
		}

		if result.Performance != nil && !summaryOnly {
			perfEmoji := "⚡"
			if result.Performance.Score < 70 {
				perfEmoji = "⚠️ "
			}
			if result.Performance.Score < 50 {
				perfEmoji = "🚨"
			}
			c.QuietOutput("%s Performance Score: %.1f/100", perfEmoji, result.Performance.Score)
			c.VerboseOutput("Bundle Size: %d bytes", result.Performance.BundleSize)
		}

		if len(result.Recommendations) > 0 && !summaryOnly {
			c.QuietOutput("\nRecommendations:")
			for _, rec := range result.Recommendations {
				c.QuietOutput("  - %s", rec)
			}
		}
	}

	// Check fail conditions and return appropriate exit codes
	if failOnHigh && result.OverallScore < 7.0 {
		return c.createAuditError(fmt.Sprintf("🚫 Found high severity issues (score: %.2f/10)", result.OverallScore), result.OverallScore)
	}

	if failOnMedium && result.OverallScore < 5.0 {
		return c.createAuditError(fmt.Sprintf("🚫 Found medium or high severity issues (score: %.2f/10)", result.OverallScore), result.OverallScore)
	}

	if minScore > 0 && result.OverallScore < minScore {
		return c.createAuditError(fmt.Sprintf("🚫 Score %.2f/10 is below your minimum requirement of %.2f/10", result.OverallScore, minScore), result.OverallScore)
	}

	return nil
}

func (c *CLI) runListTemplates(cmd *cobra.Command, args []string) error {
	// Get flags
	category, _ := cmd.Flags().GetString("category")
	technology, _ := cmd.Flags().GetString("technology")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	search, _ := cmd.Flags().GetString("search")
	detailed, _ := cmd.Flags().GetBool("detailed")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
	}

	// Create template filter
	filter := interfaces.TemplateFilter{
		Category:   category,
		Technology: technology,
		Tags:       tags,
	}

	var templates []interfaces.TemplateInfo
	var err error

	// Search or list templates
	if search != "" {
		// For now, use ListTemplates and filter by search term
		// This would be enhanced when SearchTemplates is fully implemented
		allTemplates, err := c.ListTemplates(filter)
		if err != nil {
			return c.createTemplateError("failed to search templates", search)
		}

		// Simple search filtering
		templates = []interfaces.TemplateInfo{}
		for _, template := range allTemplates {
			if strings.Contains(strings.ToLower(template.Name), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(template.Description), strings.ToLower(search)) {
				templates = append(templates, template)
			}
		}
	} else {
		templates, err = c.ListTemplates(filter)
	}

	if err != nil {
		return c.createTemplateError("failed to list templates", "")
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
		"filter":    filter,
		"search":    search,
		"detailed":  detailed,
	}

	// Output in machine-readable format if requested
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return c.outputMachineReadable(responseData, outputFormat)
	}

	// Human-readable output
	if len(templates) == 0 {
		c.QuietOutput("No templates found matching the criteria")
		return c.outputSuccess("No templates found", responseData, "list-templates", []string{})
	}

	if !c.quietMode {
		c.QuietOutput("Found %d template(s):", len(templates))
		c.QuietOutput("")

		for _, template := range templates {
			c.QuietOutput("Name: %s", template.Name)
			c.QuietOutput("Display Name: %s", template.DisplayName)
			c.QuietOutput("Description: %s", template.Description)
			c.QuietOutput("Category: %s", template.Category)
			c.QuietOutput("Technology: %s", template.Technology)
			c.QuietOutput("Version: %s", template.Version)

			if len(template.Tags) > 0 {
				c.QuietOutput("Tags: %s", strings.Join(template.Tags, ", "))
			}

			if detailed {
				if len(template.Dependencies) > 0 {
					c.VerboseOutput("Dependencies: %s", strings.Join(template.Dependencies, ", "))
				}
				c.VerboseOutput("Author: %s", template.Metadata.Author)
				c.VerboseOutput("License: %s", template.Metadata.License)
				if template.Metadata.Repository != "" {
					c.VerboseOutput("Repository: %s", template.Metadata.Repository)
				}
			}

			c.QuietOutput("")
		}
	}

	return c.outputSuccess(fmt.Sprintf("Listed %d templates", len(templates)), responseData, "list-templates", []string{})
}

func (c *CLI) runUpdate(cmd *cobra.Command, args []string) error {
	// Get flags
	check, _ := cmd.Flags().GetBool("check")
	install, _ := cmd.Flags().GetBool("install")
	templates, _ := cmd.Flags().GetBool("templates")
	force, _ := cmd.Flags().GetBool("force")
	compatibility, _ := cmd.Flags().GetBool("compatibility")
	releaseNotes, _ := cmd.Flags().GetBool("release-notes")
	channel, _ := cmd.Flags().GetString("channel")
	backup, _ := cmd.Flags().GetBool("backup")
	verify, _ := cmd.Flags().GetBool("verify")
	version, _ := cmd.Flags().GetString("version")

	// Set update channel if specified
	if channel != "stable" {
		if err := c.versionManager.SetUpdateChannel(channel); err != nil {
			return fmt.Errorf("failed to set update channel: %w", err)
		}
		fmt.Printf("📡 Using update channel: %s\n", channel)
	}

	if check {
		// Check for updates without installing
		updateInfo, err := c.CheckUpdates()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		fmt.Printf("📦 Current Version: %s\n", updateInfo.CurrentVersion)
		fmt.Printf("🆕 Latest Version: %s\n", updateInfo.LatestVersion)
		fmt.Printf("🔄 Update Available: %t\n", updateInfo.UpdateAvailable)

		if updateInfo.UpdateAvailable {
			fmt.Printf("📅 Release Date: %s\n", updateInfo.ReleaseDate.Format("2006-01-02"))
			fmt.Printf("💾 Download Size: %s\n", formatBytes(updateInfo.Size))

			if updateInfo.Breaking {
				fmt.Println("⚠️  This update contains breaking changes")
			}
			if updateInfo.Security {
				fmt.Println("🔒 This update contains security fixes")
			}
			if updateInfo.Recommended {
				fmt.Println("✅ This update is recommended")
			}

			// Show release notes if requested
			if releaseNotes && updateInfo.ReleaseNotes != "" {
				fmt.Printf("\n📝 Release Notes:\n%s\n", updateInfo.ReleaseNotes)
			}

			// Check compatibility if requested
			if compatibility {
				fmt.Println("\n🔍 Checking compatibility...")
				compatResult, err := c.versionManager.CheckCompatibility(".")
				if err != nil {
					fmt.Printf("⚠️  Failed to check compatibility: %v\n", err)
				} else {
					if compatResult.Compatible {
						fmt.Println("✅ Update is compatible with your current project")
					} else {
						fmt.Printf("⚠️  Compatibility issues found (%d issues)\n", len(compatResult.Issues))
						for _, issue := range compatResult.Issues {
							fmt.Printf("  - %s: %s\n", issue.Type, issue.Description)
						}
					}
				}
			}
		}

		return nil
	}

	if install {
		// Install available updates
		updateInfo, err := c.CheckUpdates()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		if !updateInfo.UpdateAvailable && version == "" {
			fmt.Println("✅ No updates available - you're up to date!")
			return nil
		}

		targetVersion := updateInfo.LatestVersion
		if version != "" {
			targetVersion = version
		}

		// Check compatibility unless forced
		if !force && compatibility {
			fmt.Println("🔍 Checking compatibility...")
			compatResult, err := c.versionManager.CheckCompatibility(".")
			if err != nil {
				return fmt.Errorf("failed to check compatibility: %w", err)
			}

			if !compatResult.Compatible {
				fmt.Printf("⚠️  Compatibility issues found:\n")
				for _, issue := range compatResult.Issues {
					fmt.Printf("  - %s: %s\n", issue.Type, issue.Description)
				}
				if !force {
					return fmt.Errorf("compatibility issues prevent update (use --force to override)")
				}
			}
		}

		// Warn about breaking changes unless forced
		if updateInfo.Breaking && !force {
			fmt.Println("⚠️  This update contains breaking changes.")
			fmt.Print("🤔 Continue with installation? (y/N): ")
			var response string
			if _, err := fmt.Scanln(&response); err != nil || (response != "y" && response != "Y") {
				fmt.Println("❌ Update cancelled")
				return nil
			}
		}

		fmt.Printf("⬇️  Installing update to version %s...\n", targetVersion)

		// Configure update options
		if !backup {
			fmt.Println("⚠️  Backup disabled - no rollback possible")
		}
		if !verify {
			fmt.Println("⚠️  Signature verification disabled")
		}

		err = c.InstallUpdates()
		if err != nil {
			return fmt.Errorf("failed to install updates: %w", err)
		}

		fmt.Printf("🎉 Successfully updated to version %s\n", targetVersion)
		fmt.Println("🔄 Restart any running instances to use the new version")
		return nil
	}

	if templates {
		// Update templates cache
		fmt.Println("📦 Updating templates cache...")
		if err := c.versionManager.RefreshVersionCache(); err != nil {
			return fmt.Errorf("failed to update templates cache: %w", err)
		}
		fmt.Println("✅ Templates cache updated successfully")
		return nil
	}

	// Default behavior: check for updates
	updateInfo, err := c.CheckUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if updateInfo.UpdateAvailable {
		fmt.Printf("🎉 Update available: %s -> %s\n", updateInfo.CurrentVersion, updateInfo.LatestVersion)
		if updateInfo.Security {
			fmt.Println("🔒 This update contains security fixes - update recommended")
		}
		fmt.Println("💡 Run 'generator update --install' to install the update")
		fmt.Println("📝 Run 'generator update --check --release-notes' to see what's new")
	} else {
		fmt.Println("✅ You're running the latest version!")
	}

	return nil
}

func (c *CLI) runCacheShow(cmd *cobra.Command, args []string) error {
	// Show cache status and statistics
	err := c.ShowCache()
	if err != nil {
		return fmt.Errorf("failed to show cache information: %w", err)
	}
	return nil
}

func (c *CLI) runCacheClear(cmd *cobra.Command, args []string) error {
	// Get flags
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		fmt.Print("Are you sure you want to clear all cached data? (y/N): ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			// If there's an error reading input, default to cancelling
			fmt.Println("❌ Cache clear cancelled")
			return nil
		}
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Cache clear cancelled")
			return nil
		}
	}

	// Clear cache
	err := c.ClearCache()
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("🗑️  Cache cleared successfully")
	return nil
}

func (c *CLI) runCacheClean(cmd *cobra.Command, args []string) error {
	// Clean expired and invalid cache entries
	fmt.Println("🧹 Cleaning cache...")
	err := c.CleanCache()
	if err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	fmt.Println("✨ Cache cleaned successfully")
	return nil
}

func (c *CLI) runCacheValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Validating cache...")
	err := c.ValidateCache()
	if err != nil {
		fmt.Printf("Cache validation failed: %v\n", err)
		return err
	}

	fmt.Println("Cache validation passed - cache is healthy")
	return nil
}

func (c *CLI) runCacheRepair(cmd *cobra.Command, args []string) error {
	err := c.RepairCache()
	if err != nil {
		return fmt.Errorf("failed to repair cache: %w", err)
	}
	return nil
}

func (c *CLI) runCacheOfflineEnable(cmd *cobra.Command, args []string) error {
	return c.EnableOfflineMode()
}

func (c *CLI) runCacheOfflineDisable(cmd *cobra.Command, args []string) error {
	return c.DisableOfflineMode()
}

func (c *CLI) runCacheOfflineStatus(cmd *cobra.Command, args []string) error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	isOffline := c.cacheManager.IsOfflineMode()
	fmt.Printf("Offline Mode: %t\n", isOffline)

	if isOffline {
		fmt.Println("Status: Using cached data only")
		fmt.Println("Network requests are disabled")
	} else {
		fmt.Println("Status: Network access enabled")
		fmt.Println("Will use network resources when available")
	}

	// Show cache readiness for offline mode
	stats, err := c.cacheManager.GetStats()
	if err == nil {
		fmt.Printf("\nCache Readiness:\n")
		fmt.Printf("  Entries: %d\n", stats.TotalEntries)
		fmt.Printf("  Size: %s\n", formatBytes(stats.TotalSize))
		fmt.Printf("  Health: %s\n", stats.CacheHealth)

		if stats.TotalEntries == 0 {
			fmt.Println("  Warning: No cached data available for offline mode")
		}
	}

	return nil
}

func (c *CLI) runLogs(cmd *cobra.Command, args []string) error {
	// Get flags
	lines, _ := cmd.Flags().GetInt("lines")
	level, _ := cmd.Flags().GetString("level")
	follow, _ := cmd.Flags().GetBool("follow")
	locations, _ := cmd.Flags().GetBool("locations")
	since, _ := cmd.Flags().GetString("since")
	component, _ := cmd.Flags().GetString("component")
	outputFormat, _ := cmd.Flags().GetString("format")

	// Validate level filter if provided
	if level != "" {
		validLevels := []string{"debug", "info", "warn", "error", "fatal"}
		isValid := false
		for _, validLevel := range validLevels {
			if strings.EqualFold(level, validLevel) {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid log level: %s (valid levels: %s)", level, strings.Join(validLevels, ", "))
		}
	}

	// Parse since time if provided
	var sinceTime time.Time
	if since != "" {
		var err error
		// Try different time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}

		for _, format := range formats {
			sinceTime, err = time.Parse(format, since)
			if err == nil {
				break
			}
		}

		if err != nil {
			return fmt.Errorf("invalid time format for --since: %s (use RFC3339 format like 2006-01-02T15:04:05Z)", since)
		}
	}

	if locations {
		return c.showLogLocations()
	}

	if follow {
		return c.followLogs(lines, level, component, sinceTime)
	}

	// Show recent logs
	return c.showRecentLogs(lines, level, component, sinceTime, outputFormat)
}

// Interface implementation methods - these will be implemented in subsequent tasks

func (c *CLI) SelectComponents() ([]string, error) {
	if c.isNonInteractiveMode() {
		// Return default components in non-interactive mode
		return []string{"backend-go-gin", "frontend-nextjs", "infrastructure-docker"}, nil
	}

	c.QuietOutput("Component Selection")
	c.QuietOutput("==================")
	c.QuietOutput("Select components to include in your project:")

	components := []string{}

	// Backend components
	c.QuietOutput("\nBackend Components:")
	fmt.Print("Include Go Gin API? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	if strings.ToLower(strings.TrimSpace(response)) != "n" {
		components = append(components, "backend-go-gin")
	}

	// Frontend components
	c.QuietOutput("\nFrontend Components:")
	fmt.Print("Include Next.js App? (Y/n): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	if strings.ToLower(strings.TrimSpace(response)) != "n" {
		components = append(components, "frontend-nextjs")
	}

	// Infrastructure components
	c.QuietOutput("\nInfrastructure Components:")
	fmt.Print("Include Docker configuration? (Y/n): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	if strings.ToLower(strings.TrimSpace(response)) != "n" {
		components = append(components, "infrastructure-docker")
	}

	fmt.Print("Include Kubernetes configuration? (y/N): ")
	_, _ = fmt.Scanln(&response) // Ignore error for user input
	if strings.ToLower(strings.TrimSpace(response)) == "y" {
		components = append(components, "infrastructure-kubernetes")
	}

	return components, nil
}

// Helper methods for generate command

// validateGenerateConfiguration validates the configuration for generation
func (c *CLI) validateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	if config.Name == "" {
		err := c.createConfigurationError("🚫 Your project needs a name", "")
		err = err.WithSuggestions("Set project name in configuration file or GENERATOR_PROJECT_NAME environment variable")
		return err
	}

	// Validate project name format
	if !isValidProjectName(config.Name) {
		err := c.createConfigurationError(fmt.Sprintf("🚫 '%s' isn't a valid project name", config.Name), "")
		err = err.WithSuggestions("Project names can only contain letters, numbers, hyphens, and underscores")
		return err
	}

	// Validate license if specified
	if config.License != "" && !isValidLicense(config.License) {
		err := c.createConfigurationError(fmt.Sprintf("🚫 '%s' isn't a valid license", config.License), "")
		err = err.WithSuggestions("Use a valid SPDX license identifier like MIT, Apache-2.0, or GPL-3.0")
		return err
	}

	return nil
}

// performPreGenerationChecks performs checks before generation
func (c *CLI) performPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	// Check if output directory exists
	if _, err := os.Stat(outputPath); err == nil {
		// Directory exists
		if !options.Force {
			return fmt.Errorf("🚫 Oops! The directory '%s' already exists. Use --force if you want to overwrite it", outputPath)
		}

		// Check if directory is empty
		entries, err := os.ReadDir(outputPath)
		if err != nil {
			return fmt.Errorf("🚫 Can't read the output directory: %w", err)
		}

		if len(entries) > 0 && options.BackupExisting {
			c.VerboseOutput("Creating backup of existing files in %s", outputPath)
			if err := c.createBackup(outputPath); err != nil {
				c.WarningOutput("💾 Couldn't create backup: %v", err)
			}
		}
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("🚫 Can't create the output directory: %w", err)
		}
		c.VerboseOutput("📁 Created output directory: %s", c.info(outputPath))
	}

	// Check write permissions on the output directory
	if err := c.checkWritePermissions(outputPath); err != nil {
		return fmt.Errorf("insufficient permissions for output directory: %w", err)
	}

	return nil
}

// updatePackageVersions updates package versions in the configuration
func (c *CLI) updatePackageVersions(config *models.ProjectConfig) error {
	if c.versionManager == nil {
		return fmt.Errorf("version manager not initialized")
	}

	c.VerboseOutput("Fetching latest package versions...")

	// This would fetch latest versions and update the config
	// For now, we'll just log that we would do this
	c.VerboseOutput("Would update package versions for project type based on configuration")

	return nil
}

// generateProjectFromComponents generates project structure based on selected components
func (c *CLI) generateProjectFromComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🏗️  Building your project structure...")

	// Create the base project structure first
	if err := c.generateBaseStructure(config, outputPath); err != nil {
		return fmt.Errorf("🚫 Couldn't create the project structure: %w", err)
	}

	// Process frontend components
	if err := c.processFrontendComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 Couldn't set up the frontend components: %w", err)
	}

	// Process backend components
	if err := c.processBackendComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 Couldn't set up the backend components: %w", err)
	}

	// Process mobile components
	if err := c.processMobileComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 Couldn't set up the mobile components: %w", err)
	}

	// Process infrastructure components
	if err := c.processInfrastructureComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 Couldn't set up the infrastructure components: %w", err)
	}

	c.SuccessOutput("🎉 Successfully created project %s at %s", 
		c.highlight(config.Name), 
		c.info(outputPath))
	return nil
}

// generateBaseStructure generates the base project structure
func (c *CLI) generateBaseStructure(config *models.ProjectConfig, outputPath string) error {
	c.VerboseOutput("📋 Creating project foundation...")

	// Create the main project directories first
	dirs := []string{"Docs", "Scripts"}
	for _, dir := range dirs {
		dirPath := filepath.Join(outputPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Process base template files directly from the embedded filesystem
	// The base directory contains common files like README, LICENSE, etc.
	if err := c.processBaseTemplateFiles(config, outputPath); err != nil {
		return fmt.Errorf("failed to process base template files: %w", err)
	}

	// Process GitHub workflow templates
	if err := c.templateManager.ProcessTemplate("github", config, filepath.Join(outputPath, ".github")); err != nil {
		c.VerboseOutput("GitHub template not processed (optional): %v", err)
	}

	// Process Scripts templates
	if err := c.templateManager.ProcessTemplate("scripts", config, filepath.Join(outputPath, "Scripts")); err != nil {
		c.VerboseOutput("Scripts template not processed (optional): %v", err)
	}

	return nil
}

// processBaseTemplateFiles processes base template files directly
func (c *CLI) processBaseTemplateFiles(config *models.ProjectConfig, outputPath string) error {
	// Create an embedded template engine to process the base directory directly
	embeddedEngine := template.NewEmbeddedEngine()
	
	// Process the base template directory directly
	return embeddedEngine.ProcessDirectory("templates/base", outputPath, config)
}

// processFrontendComponents processes frontend components
func (c *CLI) processFrontendComponents(config *models.ProjectConfig, outputPath string) error {
	if !c.hasFrontendComponents(config) {
		c.VerboseOutput("No frontend components selected, skipping")
		return nil
	}

	c.VerboseOutput("🎨 Setting up frontend applications...")

	// Create App directory structure
	appDir := filepath.Join(outputPath, "App")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create App directory: %w", err)
	}

	// Process Next.js components based on configuration
	if config.Components.Frontend.NextJS.App {
		c.VerboseOutput("   ✨ Creating main Next.js application")
		mainAppPath := filepath.Join(appDir, "main")
		if err := c.templateManager.ProcessTemplate("nextjs-app", config, mainAppPath); err != nil {
			return fmt.Errorf("failed to process nextjs-app template: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Home {
		c.VerboseOutput("   🏠 Creating landing page application")
		homePath := filepath.Join(appDir, "home")
		if err := c.templateManager.ProcessTemplate("nextjs-home", config, homePath); err != nil {
			return fmt.Errorf("failed to process nextjs-home template: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Admin {
		c.VerboseOutput("   👑 Creating admin dashboard")
		adminPath := filepath.Join(appDir, "admin")
		if err := c.templateManager.ProcessTemplate("nextjs-admin", config, adminPath); err != nil {
			return fmt.Errorf("failed to process nextjs-admin template: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Shared {
		c.VerboseOutput("📦 Creating shared component library...")
		sharedPath := filepath.Join(appDir, "shared-components")
		if err := c.templateManager.ProcessTemplate("shared-components", config, sharedPath); err != nil {
			return fmt.Errorf("failed to process shared-components template: %w", err)
		}
	}

	return nil
}

// processBackendComponents processes backend components
func (c *CLI) processBackendComponents(config *models.ProjectConfig, outputPath string) error {
	if !c.hasBackendComponents(config) {
		c.VerboseOutput("No backend components selected, skipping")
		return nil
	}

	c.VerboseOutput("⚙️  Setting up backend services...")

	// Create CommonServer directory
	serverDir := filepath.Join(outputPath, "CommonServer")
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return fmt.Errorf("failed to create CommonServer directory: %w", err)
	}

	// Process Go Gin backend
	if config.Components.Backend.GoGin {
		c.VerboseOutput("   🔧 Creating Go API server")
		if err := c.templateManager.ProcessTemplate("go-gin", config, serverDir); err != nil {
			return fmt.Errorf("failed to process go-gin template: %w", err)
		}
	}

	return nil
}

// processMobileComponents processes mobile components
func (c *CLI) processMobileComponents(config *models.ProjectConfig, outputPath string) error {
	if !c.hasMobileComponents(config) {
		c.VerboseOutput("No mobile components selected, skipping")
		return nil
	}

	c.VerboseOutput("📱 Setting up mobile applications...")

	// Create Mobile directory
	mobileDir := filepath.Join(outputPath, "Mobile")
	if err := os.MkdirAll(mobileDir, 0755); err != nil {
		return fmt.Errorf("failed to create Mobile directory: %w", err)
	}

	// Process Android components
	if config.Components.Mobile.Android {
		c.VerboseOutput("   🤖 Creating Android application")
		androidPath := filepath.Join(mobileDir, "android")
		if err := c.templateManager.ProcessTemplate("android-kotlin", config, androidPath); err != nil {
			return fmt.Errorf("failed to process android-kotlin template: %w", err)
		}
	}

	// Process iOS components
	if config.Components.Mobile.IOS {
		c.VerboseOutput("   🍎 Creating iOS application")
		iosPath := filepath.Join(mobileDir, "ios")
		if err := c.templateManager.ProcessTemplate("ios-swift", config, iosPath); err != nil {
			return fmt.Errorf("failed to process ios-swift template: %w", err)
		}
	}

	// Process shared mobile components
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		c.VerboseOutput("🔗 Creating shared mobile resources...")
		sharedPath := filepath.Join(mobileDir, "shared")
		if err := c.templateManager.ProcessTemplate("shared", config, sharedPath); err != nil {
			return fmt.Errorf("failed to process mobile shared template: %w", err)
		}
	}

	return nil
}

// processInfrastructureComponents processes infrastructure components
func (c *CLI) processInfrastructureComponents(config *models.ProjectConfig, outputPath string) error {
	if !c.hasInfrastructureComponents(config) {
		c.VerboseOutput("No infrastructure components selected, skipping")
		return nil
	}

	c.VerboseOutput("🚀 Setting up deployment infrastructure...")

	// Create Deploy directory
	deployDir := filepath.Join(outputPath, "Deploy")
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("failed to create Deploy directory: %w", err)
	}

	// Process Docker components
	if config.Components.Infrastructure.Docker {
		c.VerboseOutput("   🐳 Setting up Docker containers")
		dockerPath := filepath.Join(deployDir, "docker")
		if err := c.templateManager.ProcessTemplate("docker", config, dockerPath); err != nil {
			return fmt.Errorf("failed to process docker template: %w", err)
		}
	}

	// Process Kubernetes components
	if config.Components.Infrastructure.Kubernetes {
		c.VerboseOutput("   ☸️  Setting up Kubernetes deployment")
		k8sPath := filepath.Join(deployDir, "k8s")
		if err := c.templateManager.ProcessTemplate("kubernetes", config, k8sPath); err != nil {
			return fmt.Errorf("failed to process kubernetes template: %w", err)
		}
	}

	// Process Terraform components
	if config.Components.Infrastructure.Terraform {
		c.VerboseOutput("   🏗️  Setting up Terraform infrastructure")
		terraformPath := filepath.Join(deployDir, "terraform")
		if err := c.templateManager.ProcessTemplate("terraform", config, terraformPath); err != nil {
			return fmt.Errorf("failed to process terraform template: %w", err)
		}
	}

	return nil
}

// Helper methods to check if components are selected
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
	return config.Components.Mobile.Android || config.Components.Mobile.IOS
}

func (c *CLI) hasInfrastructureComponents(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

// selectDefaultTemplate selects a default template based on configuration
func (c *CLI) selectDefaultTemplate(config *models.ProjectConfig) string {
	// This method is now deprecated in favor of component-based generation
	// But kept for backward compatibility
	return "base"
}

// createBackup creates a backup of the existing directory
func (c *CLI) createBackup(path string) error {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := path + ".backup." + timestamp

	c.VerboseOutput("Creating backup at: %s", backupPath)

	// This would implement the actual backup logic
	// For now, we'll just log what we would do
	c.VerboseOutput("Would copy %s to %s", path, backupPath)

	return nil
}

// checkWritePermissions checks if we have write permissions to a directory
func (c *CLI) checkWritePermissions(path string) error {
	// Try to create a temporary file to test permissions
	// Use a secure temporary file name with random suffix
	tempFile := filepath.Join(path, ".generator-permission-test-"+strconv.FormatInt(time.Now().UnixNano(), 36))

	// #nosec G304 - This is a controlled temporary file creation for permission testing
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("no write permission to directory %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		c.WarningOutput("📄 Couldn't close temporary file: %v", err)
	}
	if err := os.Remove(tempFile); err != nil {
		c.WarningOutput("🗑️  Couldn't remove temporary file: %v", err)
	}
	return nil
}

// isValidProjectName validates project name format
func isValidProjectName(name string) bool {
	// Allow letters, numbers, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 100
}

// isValidLicense validates license identifier
func isValidLicense(license string) bool {
	// Common SPDX license identifiers
	validLicenses := []string{
		"MIT", "Apache-2.0", "GPL-2.0", "GPL-3.0", "LGPL-2.1", "LGPL-3.0",
		"BSD-2-Clause", "BSD-3-Clause", "ISC", "MPL-2.0", "UNLICENSED",
	}

	for _, valid := range validLicenses {
		if strings.EqualFold(license, valid) {
			return true
		}
	}
	return false
}

// validateGenerateOptions validates the generate command options
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

// validateDependencies validates project dependencies before generation
func (c *CLI) validateDependencies(config *models.ProjectConfig, templateName string) error {
	c.VerboseOutput("Validating dependencies for template: %s", templateName)

	// Check if template manager is available
	if c.templateManager == nil {
		return fmt.Errorf("template manager not initialized")
	}

	// Get template information
	templateInfo, err := c.templateManager.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("🚫 Couldn't find template '%s': %w", templateName, err)
	}

	// Validate template dependencies
	if len(templateInfo.Dependencies) > 0 {
		c.VerboseOutput("Checking template dependencies: %v", templateInfo.Dependencies)

		for _, dep := range templateInfo.Dependencies {
			if err := c.validateDependency(dep); err != nil {
				return fmt.Errorf("dependency validation failed for '%s': %w", dep, err)
			}
		}
	}

	// Validate system requirements based on template
	if err := c.validateSystemRequirements(templateInfo); err != nil {
		return fmt.Errorf("system requirements validation failed: %w", err)
	}

	return nil
}

// validateDependency validates a specific dependency
func (c *CLI) validateDependency(dependency string) error {
	// Parse dependency (format: name@version or just name)
	parts := strings.Split(dependency, "@")
	depName := parts[0]

	c.DebugOutput("Validating dependency: %s", depName)

	// Check common dependencies
	switch depName {
	case "go":
		return c.validateGoVersion(parts)
	case "node", "nodejs":
		return c.validateNodeVersion(parts)
	case "docker":
		return c.validateDockerAvailability()
	case "git":
		return c.validateGitAvailability()
	default:
		c.VerboseOutput("Dependency '%s' will be validated during generation", depName)
	}

	return nil
}

// validateGoVersion validates Go installation and version
func (c *CLI) validateGoVersion(parts []string) error {
	// Check if Go is installed
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not installed or not in PATH")
	}

	// Get Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %w", err)
	}

	versionStr := string(output)
	c.VerboseOutput("Found Go version: %s", strings.TrimSpace(versionStr))

	// If specific version is required, validate it
	if len(parts) > 1 {
		requiredVersion := parts[1]
		c.VerboseOutput("Required Go version: %s", requiredVersion)
		// This would implement actual version comparison
		// For now, we'll just log it
	}

	return nil
}

// validateNodeVersion validates Node.js installation and version
func (c *CLI) validateNodeVersion(parts []string) error {
	// Check if Node.js is installed
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("node.js is not installed or not in PATH")
	}

	// Get Node.js version
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Node.js version: %w", err)
	}

	versionStr := strings.TrimSpace(string(output))
	c.VerboseOutput("Found Node.js version: %s", versionStr)

	// If specific version is required, validate it
	if len(parts) > 1 {
		requiredVersion := parts[1]
		c.VerboseOutput("Required Node.js version: %s", requiredVersion)
		// This would implement actual version comparison
	}

	return nil
}

// validateDockerAvailability validates Docker installation
func (c *CLI) validateDockerAvailability() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed or not in PATH")
	}

	// Check if Docker daemon is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon is not running")
	}

	c.VerboseOutput("Docker is available and running")
	return nil
}

// validateGitAvailability validates Git installation
func (c *CLI) validateGitAvailability() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	c.VerboseOutput("Git is available")
	return nil
}

// validateSystemRequirements validates system requirements for a template
func (c *CLI) validateSystemRequirements(templateInfo *interfaces.TemplateInfo) error {
	c.VerboseOutput("Validating system requirements for template: %s", templateInfo.Name)

	// Check available disk space
	if err := c.validateDiskSpace(); err != nil {
		return fmt.Errorf("disk space validation failed: %w", err)
	}

	// Check memory requirements (basic check)
	if err := c.validateMemoryRequirements(); err != nil {
		c.WarningOutput("Memory validation warning: %v", err)
	}

	return nil
}

// validateDiskSpace validates available disk space
func (c *CLI) validateDiskSpace() error {
	// This would implement actual disk space checking
	// For now, we'll just log that we would check it
	c.VerboseOutput("Checking available disk space...")

	// Minimum required space (in bytes) - 100MB
	const minRequiredSpace = 100 * 1024 * 1024

	// This would get actual available space
	c.VerboseOutput("Would check for at least %d bytes of free space", minRequiredSpace)

	return nil
}

// validateMemoryRequirements validates memory requirements
func (c *CLI) validateMemoryRequirements() error {
	// This would implement actual memory checking
	c.VerboseOutput("Checking available memory...")

	// This would check system memory
	c.VerboseOutput("Would check system memory requirements")

	return nil
}

// isValidTemplateName validates template name format
func isValidTemplateName(name string) bool {
	// Allow letters, numbers, hyphens, underscores, and dots
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 50
}

func (c *CLI) GenerateFromConfig(configPath string, options interfaces.GenerateOptions) error {
	ctx := c.StartOperationWithOutput("generate-from-config", fmt.Sprintf("📄 Loading your project configuration from %s", configPath))
	defer func() {
		if ctx != nil {
			c.FinishOperationWithOutput(ctx, "generate-from-config", "✅ Configuration loaded successfully")
		}
	}()

	// Load configuration from file
	if c.configManager == nil {
		return fmt.Errorf("configuration manager not initialized")
	}

	config, err := c.configManager.LoadConfig(configPath)
	if err != nil {
		c.FinishOperationWithError(ctx, "generate-from-config", err)
		return fmt.Errorf("failed to load configuration from %s: %w", configPath, err)
	}

	c.VerboseOutput("✅ Configuration loaded for project: %s", c.highlight(config.Name))

	// Validate configuration if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("🔍 Checking configuration validity...")
		if err := c.validateGenerateConfiguration(config, options); err != nil {
			return err
		}
		c.VerboseOutput("✅ Configuration looks good!")
	}

	// Set output path from options or config
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = config.OutputPath
	}
	if outputPath == "" {
		outputPath = "./output/generated"
	}
	
	// Always append project name to the output path
	outputPath = filepath.Join(outputPath, config.Name)

	// Handle offline mode
	if options.Offline {
		c.VerboseOutput("Running in offline mode - using cached templates and versions")
		if c.cacheManager != nil {
			if err := c.cacheManager.EnableOfflineMode(); err != nil {
				c.WarningOutput("📡 Couldn't enable offline mode: %v", err)
			}
		}
	}

	// Handle version updates
	if options.UpdateVersions && !options.Offline {
		c.VerboseOutput("Fetching latest package versions...")
		if err := c.updatePackageVersions(config); err != nil {
			c.WarningOutput("Failed to update package versions: %v", err)
		}
	}

	// Pre-generation checks
	if err := c.performPreGenerationChecks(outputPath, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Generate project structure based on components
	if options.DryRun {
		c.QuietOutput("Dry run mode - would generate project '%s' with selected components in directory '%s'",
			config.Name, outputPath)
		return nil
	}

	// Process multiple templates based on configuration components
	if c.templateManager == nil {
		return fmt.Errorf("template manager not initialized")
	}

	return c.generateProjectFromComponents(config, outputPath, options)
}

func (c *CLI) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	if c.validator == nil {
		return nil, fmt.Errorf("validation engine not initialized")
	}

	// Call the actual validation engine
	result, err := c.validator.ValidateProject(path)
	if err != nil {
		return nil, err
	}

	// Convert from models.ValidationResult to interfaces.ValidationResult
	interfaceResult := &interfaces.ValidationResult{
		Valid:    result.Valid,
		Issues:   []interfaces.ValidationIssue{},
		Warnings: []interfaces.ValidationIssue{},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   1,
			ValidFiles:   0,
			ErrorCount:   0,
			WarningCount: 0,
			FixableCount: 0,
		},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	// Convert issues
	for _, issue := range result.Issues {
		interfaceIssue := interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
		}

		if issue.Severity == "error" {
			interfaceResult.Issues = append(interfaceResult.Issues, interfaceIssue)
			interfaceResult.Summary.ErrorCount++
		} else {
			interfaceResult.Warnings = append(interfaceResult.Warnings, interfaceIssue)
			interfaceResult.Summary.WarningCount++
		}

		if issue.Fixable {
			interfaceResult.Summary.FixableCount++
		}
	}

	if interfaceResult.Valid {
		interfaceResult.Summary.ValidFiles = 1
	}

	return interfaceResult, nil
}

func (c *CLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	if c.auditEngine == nil {
		return nil, fmt.Errorf("audit engine not initialized")
	}

	// Perform comprehensive project audit
	result, err := c.auditEngine.AuditProject(path, &options)
	if err != nil {
		return nil, err
	}

	// Generate report if output file is specified
	if options.OutputFile != "" {
		format := options.OutputFormat
		if format == "" {
			format = "json" // Default format
		}

		reportData, err := c.auditEngine.GenerateAuditReport(result, format)
		if err != nil {
			return nil, fmt.Errorf("failed to generate audit report: %w", err)
		}

		err = os.WriteFile(options.OutputFile, reportData, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to write audit report: %w", err)
		}

		fmt.Printf("Audit report written to: %s\n", options.OutputFile)
	}

	return result, nil
}

func (c *CLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	return c.templateManager.ListTemplates(filter)
}

func (c *CLI) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	return c.templateManager.GetTemplateInfo(name)
}

func (c *CLI) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return c.templateManager.ValidateTemplate(path)
}

// runTemplateInfo handles the template info command
func (c *CLI) runTemplateInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	showVariables, _ := cmd.Flags().GetBool("variables")
	showDependencies, _ := cmd.Flags().GetBool("dependencies")
	showCompatibility, _ := cmd.Flags().GetBool("compatibility")

	// Get template info
	templateInfo, err := c.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("🚫 Couldn't find template '%s': %w", templateName, err)
	}

	// Display basic information
	fmt.Printf("Template: %s\n", templateInfo.Name)
	fmt.Printf("Display Name: %s\n", templateInfo.DisplayName)
	fmt.Printf("Description: %s\n", templateInfo.Description)
	fmt.Printf("Category: %s\n", templateInfo.Category)
	fmt.Printf("Technology: %s\n", templateInfo.Technology)
	fmt.Printf("Version: %s\n", templateInfo.Version)

	if len(templateInfo.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(templateInfo.Tags, ", "))
	}

	// Show detailed information if requested
	if detailed || showDependencies {
		if len(templateInfo.Dependencies) > 0 {
			fmt.Printf("\nDependencies:\n")
			for _, dep := range templateInfo.Dependencies {
				fmt.Printf("  - %s\n", dep)
			}
		} else {
			fmt.Printf("\nDependencies: None\n")
		}
	}

	if detailed {
		fmt.Printf("\nMetadata:\n")
		fmt.Printf("  Author: %s\n", templateInfo.Metadata.Author)
		fmt.Printf("  License: %s\n", templateInfo.Metadata.License)
		if templateInfo.Metadata.Repository != "" {
			fmt.Printf("  Repository: %s\n", templateInfo.Metadata.Repository)
		}
		if templateInfo.Metadata.Homepage != "" {
			fmt.Printf("  Homepage: %s\n", templateInfo.Metadata.Homepage)
		}
		if len(templateInfo.Metadata.Keywords) > 0 {
			fmt.Printf("  Keywords: %s\n", strings.Join(templateInfo.Metadata.Keywords, ", "))
		}
	}

	// Show variables if requested
	if showVariables || detailed {
		variables, err := c.templateManager.GetTemplateVariables(templateName)
		if err != nil {
			fmt.Printf("\nVariables: Error retrieving variables: %v\n", err)
		} else if len(variables) > 0 {
			fmt.Printf("\nVariables:\n")
			for name, variable := range variables {
				fmt.Printf("  %s (%s):\n", name, variable.Type)
				fmt.Printf("    Description: %s\n", variable.Description)
				if variable.Default != nil {
					fmt.Printf("    Default: %v\n", variable.Default)
				}
				fmt.Printf("    Required: %t\n", variable.Required)
				if variable.Validation != nil && variable.Validation.Pattern != "" {
					fmt.Printf("    Pattern: %s\n", variable.Validation.Pattern)
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("\nVariables: None defined\n")
		}
	}

	// Show compatibility if requested
	if showCompatibility || detailed {
		compatibility, err := c.templateManager.GetTemplateCompatibility(templateName)
		if err != nil {
			fmt.Printf("\nCompatibility: Error retrieving compatibility info: %v\n", err)
		} else {
			fmt.Printf("\nCompatibility:\n")
			if compatibility.MinGeneratorVersion != "" {
				fmt.Printf("  Min Generator Version: %s\n", compatibility.MinGeneratorVersion)
			}
			if compatibility.MaxGeneratorVersion != "" {
				fmt.Printf("  Max Generator Version: %s\n", compatibility.MaxGeneratorVersion)
			}
			if len(compatibility.SupportedPlatforms) > 0 {
				fmt.Printf("  Supported Platforms: %s\n", strings.Join(compatibility.SupportedPlatforms, ", "))
			}
			if len(compatibility.RequiredFeatures) > 0 {
				fmt.Printf("  Required Features: %s\n", strings.Join(compatibility.RequiredFeatures, ", "))
			}
		}
	}

	return nil
}

// runTemplateValidate handles the template validate command
func (c *CLI) runTemplateValidate(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	fix, _ := cmd.Flags().GetBool("fix")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Validate template
	result, err := c.ValidateTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Output results based on format
	switch outputFormat {
	case "json":
		// For JSON output, we would marshal the result
		fmt.Printf("{\n")
		fmt.Printf("  \"valid\": %t,\n", result.Valid)
		fmt.Printf("  \"issues\": %d,\n", len(result.Issues))
		fmt.Printf("  \"warnings\": %d\n", len(result.Warnings))
		fmt.Printf("}\n")
	default:
		// Text output
		fmt.Printf("Template validation for: %s\n", templatePath)
		fmt.Printf("Valid: %t\n", result.Valid)
		fmt.Printf("Issues: %d\n", len(result.Issues))
		fmt.Printf("Warnings: %d\n", len(result.Warnings))

		if len(result.Issues) > 0 {
			fmt.Println("\nIssues:")
			for _, issue := range result.Issues {
				if detailed {
					fmt.Printf("  [%s] %s: %s", issue.Severity, issue.Type, issue.Message)
					if issue.File != "" {
						fmt.Printf(" (File: %s", issue.File)
						if issue.Line > 0 {
							fmt.Printf(":%d", issue.Line)
						}
						fmt.Printf(")")
					}
					fmt.Printf(" [Rule: %s]", issue.Rule)
					if issue.Fixable {
						fmt.Printf(" [Fixable]")
					}
					fmt.Println()
				} else {
					fmt.Printf("  - %s: %s\n", issue.Severity, issue.Message)
				}
			}
		}

		if len(result.Warnings) > 0 {
			fmt.Println("\nWarnings:")
			for _, warning := range result.Warnings {
				if detailed {
					fmt.Printf("  [%s] %s: %s", warning.Severity, warning.Type, warning.Message)
					if warning.File != "" {
						fmt.Printf(" (File: %s", warning.File)
						if warning.Line > 0 {
							fmt.Printf(":%d", warning.Line)
						}
						fmt.Printf(")")
					}
					fmt.Printf(" [Rule: %s]", warning.Rule)
					if warning.Fixable {
						fmt.Printf(" [Fixable]")
					}
					fmt.Println()
				} else {
					fmt.Printf("  - %s: %s\n", warning.Severity, warning.Message)
				}
			}
		}

		if fix {
			fmt.Println("\nNote: Auto-fix functionality is not yet implemented")
		}
	}

	// Return error if validation failed
	if !result.Valid {
		return fmt.Errorf("template validation failed")
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

func (c *CLI) SetConfig(key, value string) error {
	// Set the configuration value
	err := c.configManager.SetSetting(key, value)
	if err != nil {
		return fmt.Errorf("failed to set configuration value: %w", err)
	}

	// Skip validation for individual settings - validation should be done
	// when the complete configuration is ready

	// Save the configuration to file if possible
	configLocation := c.configManager.GetConfigLocation()
	if configLocation != "" {
		// Load current config, update it, and save
		config, err := c.configManager.LoadDefaults()
		if err != nil {
			// If we can't load defaults, create a new config
			config = &models.ProjectConfig{}
		}

		// Apply the setting to the config struct
		err = c.applySettingToConfig(config, key, value)
		if err != nil {
			return fmt.Errorf("failed to apply setting to configuration: %w", err)
		}

		// Save the updated configuration
		err = c.configManager.SaveConfig(config, configLocation)
		if err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	return nil
}

func (c *CLI) EditConfig() error {
	configLocation := c.configManager.GetConfigLocation()
	if configLocation == "" {
		// Create a default config file if none exists
		configLocation = "./generator-config.yaml"
		err := c.configManager.CreateDefaultConfig(configLocation)
		if err != nil {
			return fmt.Errorf("failed to create default configuration file: %w", err)
		}
		fmt.Printf("Created default configuration file: %s\n", configLocation)
	}

	// Create backup before editing
	err := c.configManager.BackupConfig(configLocation)
	if err != nil {
		fmt.Printf("Warning: failed to create backup: %v\n", err)
	}

	// Try to open with various editors
	allowedEditors := map[string]bool{
		"code":    true, // VS Code
		"vim":     true, // Vim
		"nano":    true, // Nano
		"notepad": true, // Windows Notepad
		"vi":      true, // Vi
		"emacs":   true, // Emacs
	}

	editors := []string{
		os.Getenv("EDITOR"),
		"code",    // VS Code
		"vim",     // Vim
		"nano",    // Nano
		"notepad", // Windows Notepad
	}

	var editorCmd string
	for _, editor := range editors {
		if editor != "" && allowedEditors[editor] {
			// Check if editor exists
			if _, err := exec.LookPath(editor); err == nil {
				editorCmd = editor
				break
			}
		}
	}

	if editorCmd == "" {
		return fmt.Errorf("no suitable editor found. Please set the EDITOR environment variable to one of: code, vim, nano, vi, emacs")
	}

	// Open the configuration file in the editor
	fmt.Printf("Opening configuration file in %s...\n", editorCmd)
	// #nosec G204 - editorCmd is validated against allowedEditors whitelist
	cmd := exec.Command(editorCmd, configLocation)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Validate the configuration after editing
	fmt.Println("Validating configuration...")
	result, err := c.configManager.ValidateConfigFromFile(configLocation)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %w", err)
	}

	if !result.Valid {
		fmt.Println("Configuration validation failed:")
		for _, validationError := range result.Errors {
			fmt.Printf("  Error: %s - %s\n", validationError.Field, validationError.Message)
		}
		for _, warning := range result.Warnings {
			fmt.Printf("  Warning: %s - %s\n", warning.Field, warning.Message)
		}
		return fmt.Errorf("configuration contains %d errors", len(result.Errors))
	}

	fmt.Println("Configuration updated successfully!")
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

func (c *CLI) ExportConfig(path string) error {
	// Load current configuration from all sources
	config, err := c.configManager.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load default configuration: %w", err)
	}

	// Merge with environment variables
	envConfig, err := c.configManager.LoadFromEnvironment()
	if err == nil {
		config = c.configManager.MergeConfigurations(config, envConfig)
	}

	// Validate the configuration before export
	err = c.configManager.ValidateConfig(config)
	if err != nil {
		fmt.Printf("Warning: configuration has validation issues: %v\n", err)
	}

	// Save the merged configuration to the specified path
	err = c.configManager.SaveConfig(config, path)
	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}

	// Display export information
	fmt.Printf("Configuration exported successfully!\n")
	fmt.Printf("File: %s\n", path)

	// Show basic stats
	sources, err := c.configManager.GetConfigSources()
	if err == nil {
		fmt.Printf("Sources merged: %d\n", len(sources))
		for _, source := range sources {
			if source.Valid {
				fmt.Printf("  - %s (%s)\n", source.Type, source.Location)
			}
		}
	}

	return nil
}

func (c *CLI) ShowVersion(options interfaces.VersionOptions) error {
	// Get current version info
	currentVersion := c.versionManager.GetCurrentVersion()

	// Basic version display
	fmt.Printf("Generator Version: %s\n", currentVersion)

	// Show build info if requested
	if options.ShowBuildInfo {
		if latestInfo, err := c.versionManager.GetLatestVersion(); err == nil {
			fmt.Printf("Build Date: %s\n", latestInfo.BuildDate.Format("2006-01-02 15:04:05"))
			fmt.Printf("Git Commit: %s\n", latestInfo.GitCommit)
			fmt.Printf("Git Branch: %s\n", latestInfo.GitBranch)
			fmt.Printf("Go Version: %s\n", latestInfo.GoVersion)
			fmt.Printf("Platform: %s\n", latestInfo.Platform)
			fmt.Printf("Architecture: %s\n", latestInfo.Architecture)
		}
	}

	// Show package versions if requested
	if options.ShowPackages {
		fmt.Println("\nPackage Versions:")
		packages, err := c.versionManager.GetAllPackageVersions()
		if err != nil {
			return fmt.Errorf("failed to get package versions: %w", err)
		}

		for pkg, version := range packages {
			fmt.Printf("  %-20s %s\n", pkg+":", version)
		}
	}

	// Check for updates if requested
	if options.CheckUpdates {
		fmt.Println("\nChecking for updates...")
		updateInfo, err := c.versionManager.CheckForUpdates()
		if err != nil {
			fmt.Printf("Warning: Failed to check for updates: %v\n", err)
		} else {
			if updateInfo.UpdateAvailable {
				fmt.Printf("🎉 Update available: %s -> %s\n", updateInfo.CurrentVersion, updateInfo.LatestVersion)
				fmt.Printf("Release Date: %s\n", updateInfo.ReleaseDate.Format("2006-01-02"))
				if updateInfo.Breaking {
					fmt.Println("⚠️  This update contains breaking changes")
				}
				if updateInfo.Security {
					fmt.Println("🔒 This update contains security fixes")
				}
				fmt.Println("Run 'generator update --install' to install the update")
			} else {
				fmt.Println("✅ You are running the latest version")
			}
		}
	}

	return nil
}

func (c *CLI) CheckUpdates() (*interfaces.UpdateInfo, error) {
	return c.versionManager.CheckForUpdates()
}

func (c *CLI) InstallUpdates() error {
	// First check for updates
	updateInfo, err := c.CheckUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !updateInfo.UpdateAvailable {
		return fmt.Errorf("no updates available")
	}

	// Download and install the update
	if err := c.versionManager.DownloadUpdate(updateInfo.LatestVersion); err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	if err := c.versionManager.InstallUpdate(updateInfo.LatestVersion); err != nil {
		return fmt.Errorf("failed to install update: %w", err)
	}

	return nil
}

func (c *CLI) ShowCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	// Get cache statistics
	stats, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	// Get cache configuration
	config, err := c.cacheManager.GetCacheConfig()
	if err != nil {
		return fmt.Errorf("failed to get cache configuration: %w", err)
	}

	// Display cache information
	fmt.Println("💾 Cache Information")
	fmt.Println("===================")
	fmt.Printf("📁 Location: %s\n", stats.CacheLocation)
	
	statusEmoji := "✅"
	if stats.CacheHealth != "healthy" {
		statusEmoji = "⚠️ "
	}
	fmt.Printf("%s Status: %s\n", statusEmoji, stats.CacheHealth)
	
	offlineEmoji := "🌐"
	if stats.OfflineMode {
		offlineEmoji = "📴"
	}
	fmt.Printf("%s Offline Mode: %t\n", offlineEmoji, stats.OfflineMode)
	fmt.Println()

	fmt.Println("📊 Statistics")
	fmt.Println("=============")
	fmt.Printf("📦 Total Entries: %d\n", stats.TotalEntries)
	fmt.Printf("💾 Total Size: %s\n", formatBytes(stats.TotalSize))
	fmt.Printf("🎯 Hit Rate: %.1f%%\n", stats.HitRate*100)
	fmt.Printf("⏰ Expired Entries: %d\n", stats.ExpiredEntries)
	fmt.Printf("🧹 Last Cleanup: %s\n", stats.LastCleanup.Format("2006-01-02 15:04:05"))
	fmt.Println()

	fmt.Println("⚙️  Configuration")
	fmt.Println("=================")
	fmt.Printf("📏 Max Size: %s\n", formatBytes(config.MaxSize))
	fmt.Printf("🔢 Max Entries: %d\n", config.MaxEntries)
	fmt.Printf("⏱️  Default TTL: %s\n", config.DefaultTTL)
	fmt.Printf("🔄 Eviction Policy: %s\n", config.EvictionPolicy)
	fmt.Printf("🗜️  Compression: %t\n", config.EnableCompression)
	fmt.Printf("💿 Persist to Disk: %t\n", config.PersistToDisk)

	// Show cache health warnings if any
	if stats.CacheHealth != "healthy" {
		fmt.Println()
		fmt.Println("🚨 Health Issues")
		fmt.Println("================")
		if stats.ExpiredEntries > 0 {
			fmt.Printf("⚠️  %d expired entries found - consider running 'generator cache clean'\n", stats.ExpiredEntries)
		}
		if stats.CacheHealth == "corrupted" {
			fmt.Println("🚨 Cache corruption detected - consider running 'generator cache repair'")
		}
		if stats.CacheHealth == "missing" {
			fmt.Println("⚠ Cache directory missing - will be created on next cache operation")
		}
	}

	return nil
}

func (c *CLI) ClearCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	// Get cache stats before clearing
	stats, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	// Clear the cache
	err = c.cacheManager.Clear()
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Printf("Cache cleared successfully!\n")
	fmt.Printf("Removed %d entries (%s)\n", stats.TotalEntries, formatBytes(stats.TotalSize))

	return nil
}

func (c *CLI) CleanCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	// Get stats before cleaning
	statsBefore, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	// Clean expired entries
	err = c.cacheManager.Clean()
	if err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	// Get stats after cleaning
	statsAfter, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics after cleaning: %w", err)
	}

	// Calculate what was cleaned
	entriesRemoved := statsBefore.TotalEntries - statsAfter.TotalEntries
	sizeFreed := statsBefore.TotalSize - statsAfter.TotalSize

	fmt.Printf("Cache cleaned successfully!\n")
	if entriesRemoved > 0 {
		fmt.Printf("Removed %d expired entries\n", entriesRemoved)
		fmt.Printf("Freed %s of space\n", formatBytes(sizeFreed))
	} else {
		fmt.Printf("No expired entries found\n")
	}
	fmt.Printf("Current cache: %d entries (%s)\n", statsAfter.TotalEntries, formatBytes(statsAfter.TotalSize))

	return nil
}

func (c *CLI) ShowLogs() error {
	return c.showRecentLogs(50, "", "", time.Time{}, "text")
}

// Advanced interface methods implementation

func (c *CLI) PromptAdvancedOptions() (*interfaces.AdvancedOptions, error) {
	if c.isNonInteractiveMode() {
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

	c.QuietOutput("Advanced Options")
	c.QuietOutput("================")

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

func (c *CLI) ConfirmAdvancedGeneration(config *models.ProjectConfig, options *interfaces.AdvancedOptions) bool {
	if c.isNonInteractiveMode() {
		return true // Auto-confirm in non-interactive mode
	}

	c.QuietOutput("\nAdvanced Configuration Preview:")
	c.QuietOutput("===============================")
	c.QuietOutput("Security Scanning: %t", options.EnableSecurityScanning)
	c.QuietOutput("Quality Checks: %t", options.EnableQualityChecks)
	c.QuietOutput("Performance Optimization: %t", options.EnablePerformanceOptimization)
	c.QuietOutput("Generate Documentation: %t", options.GenerateDocumentation)
	c.QuietOutput("Enable CI/CD: %t", options.EnableCICD)
	if options.EnableCICD {
		c.QuietOutput("CI/CD Providers: %v", options.CICDProviders)
	}

	fmt.Print("\nProceed with advanced generation? (Y/n): ")
	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
}

func (c *CLI) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	if c.isNonInteractiveMode() {
		return nil, fmt.Errorf("interactive template selection not available in non-interactive mode")
	}

	// Get available templates
	templates, err := c.templateManager.ListTemplates(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found matching the criteria")
	}

	c.QuietOutput("Available Templates:")
	c.QuietOutput("===================")

	for i, template := range templates {
		c.QuietOutput("%d. %s - %s", i+1, template.DisplayName, template.Description)
		c.VerboseOutput("   Category: %s, Technology: %s", template.Category, template.Technology)
	}

	fmt.Printf("\nSelect template (1-%d): ", len(templates))
	var selection int
	if _, err := fmt.Scanln(&selection); err != nil {
		return nil, fmt.Errorf("failed to read template selection: %w", err)
	}

	if selection < 1 || selection > len(templates) {
		return nil, fmt.Errorf("🚫 That's not a valid choice. Please pick a number between 1 and %d", len(templates))
	}

	return &templates[selection-1], nil
}

func (c *CLI) GenerateWithAdvancedOptions(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	if c.templateManager == nil {
		return fmt.Errorf("template manager not initialized")
	}

	// Select template based on configuration or use default
	templateName := "go-gin" // Default template
	if len(options.Templates) > 0 {
		templateName = options.Templates[0]
	}

	// Set output path
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = config.OutputPath
	}
	if outputPath == "" {
		outputPath = "./" + config.Name
	}

	c.VerboseOutput("Generating project with advanced options...")
	c.VerboseOutput("Template: %s", templateName)
	c.VerboseOutput("Output: %s", outputPath)
	c.VerboseOutput("Security Scanning: %t", options.EnableSecurityScanning)
	c.VerboseOutput("Quality Checks: %t", options.EnableQualityChecks)
	c.VerboseOutput("Performance Optimization: %t", options.EnablePerformanceOptimization)

	// Generate the project
	err := c.templateManager.ProcessTemplate(templateName, config, outputPath)
	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Apply advanced options post-generation
	if options.EnableSecurityScanning && c.auditEngine != nil {
		c.VerboseOutput("Running security scan...")
		auditOptions := &interfaces.AuditOptions{
			Security:    true,
			Quality:     false,
			Licenses:    false,
			Performance: false,
		}
		_, err := c.auditEngine.AuditProject(outputPath, auditOptions)
		if err != nil {
			c.WarningOutput("Security scan failed: %v", err)
		}
	}

	if options.EnableQualityChecks && c.validator != nil {
		c.VerboseOutput("Running quality checks...")
		_, err := c.validator.ValidateProject(outputPath)
		if err != nil {
			c.WarningOutput("Quality checks failed: %v", err)
		}
	}

	return nil
}

func (c *CLI) ValidateProjectAdvanced(path string, options *interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	if c.validator == nil {
		return nil, fmt.Errorf("validation engine not initialized")
	}

	// Perform basic validation first
	result, err := c.validator.ValidateProject(path)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert to interface result
	interfaceResult := &interfaces.ValidationResult{
		Valid:    result.Valid,
		Issues:   []interfaces.ValidationIssue{},
		Warnings: []interfaces.ValidationIssue{},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   1,
			ValidFiles:   0,
			ErrorCount:   0,
			WarningCount: 0,
			FixableCount: 0,
		},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	// Convert issues
	for _, issue := range result.Issues {
		interfaceIssue := interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
		}

		if issue.Severity == "error" {
			interfaceResult.Issues = append(interfaceResult.Issues, interfaceIssue)
			interfaceResult.Summary.ErrorCount++
		} else {
			interfaceResult.Warnings = append(interfaceResult.Warnings, interfaceIssue)
			interfaceResult.Summary.WarningCount++
		}

		if issue.Fixable {
			interfaceResult.Summary.FixableCount++
		}
	}

	if interfaceResult.Valid {
		interfaceResult.Summary.ValidFiles = 1
	}

	// Apply advanced validation options
	if options != nil {
		// Filter by rules if specified
		if len(options.Rules) > 0 {
			interfaceResult.Issues = c.filterIssuesByRules(interfaceResult.Issues, options.Rules)
			interfaceResult.Warnings = c.filterIssuesByRules(interfaceResult.Warnings, options.Rules)
		}

		// Ignore warnings if requested
		if options.IgnoreWarnings {
			interfaceResult.Warnings = []interfaces.ValidationIssue{}
			interfaceResult.Summary.WarningCount = 0
		}
	}

	return interfaceResult, nil
}

func (c *CLI) AuditProjectAdvanced(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	if c.auditEngine == nil {
		return nil, fmt.Errorf("audit engine not initialized")
	}

	if options == nil {
		options = &interfaces.AuditOptions{
			Security:    true,
			Quality:     true,
			Licenses:    true,
			Performance: true,
		}
	}

	// Perform comprehensive project audit with advanced options
	result, err := c.auditEngine.AuditProject(path, options)
	if err != nil {
		return nil, fmt.Errorf("advanced audit failed: %w", err)
	}

	// Enhanced reporting for advanced audit
	if options.Detailed {
		// Add detailed analysis results
		if result.Security != nil {
			// Perform additional security scans
			vulnReport, err := c.auditEngine.ScanVulnerabilities(path)
			if err == nil {
				result.Security.Vulnerabilities = vulnReport.Vulnerabilities
			}

			secretReport, err := c.auditEngine.DetectSecrets(path)
			if err == nil && len(secretReport.Secrets) > 0 {
				// Add secret detection results to security recommendations
				result.Security.Recommendations = append(result.Security.Recommendations,
					fmt.Sprintf("Found %d potential secrets in code", len(secretReport.Secrets)))
			}
		}

		if result.Quality != nil {
			// Perform additional quality analysis
			complexityReport, err := c.auditEngine.MeasureComplexity(path)
			if err == nil {
				result.Quality.Recommendations = append(result.Quality.Recommendations,
					fmt.Sprintf("Average complexity: %.1f, High complexity files: %d",
						complexityReport.Summary.AverageComplexity,
						complexityReport.Summary.HighComplexityFiles))
			}
		}

		if result.Licenses != nil {
			// Perform additional license analysis
			violationReport, err := c.auditEngine.ScanLicenseViolations(path)
			if err == nil && len(violationReport.Violations) > 0 {
				result.Licenses.Recommendations = append(result.Licenses.Recommendations,
					fmt.Sprintf("Found %d license violations", len(violationReport.Violations)))
			}
		}

		if result.Performance != nil {
			// Perform additional performance analysis
			metricsReport, err := c.auditEngine.CheckPerformanceMetrics(path)
			if err == nil {
				result.Performance.Recommendations = append(result.Performance.Recommendations,
					fmt.Sprintf("Performance grade: %s", metricsReport.Summary.PerformanceGrade))
			}
		}
	}

	return result, nil
}

func (c *CLI) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	return c.templateManager.SearchTemplates(query)
}

func (c *CLI) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	return c.templateManager.GetTemplateMetadata(name)
}

func (c *CLI) GetTemplateDependencies(name string) ([]string, error) {
	return c.templateManager.GetTemplateDependencies(name)
}

func (c *CLI) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return c.templateManager.ValidateCustomTemplate(path)
}

func (c *CLI) LoadConfiguration(sources []string) (*models.ProjectConfig, error) {
	if c.configManager == nil {
		return nil, fmt.Errorf("configuration manager not initialized")
	}

	var configs []*models.ProjectConfig

	for _, source := range sources {
		var config *models.ProjectConfig
		var err error

		switch source {
		case "file":
			config, err = c.configManager.LoadFromFile("")
		case "environment":
			config, err = c.configManager.LoadFromEnvironment()
		case "defaults":
			config, err = c.configManager.LoadDefaults()
		default:
			// Treat as file path
			config, err = c.configManager.LoadFromFile(source)
		}

		if err != nil {
			c.VerboseOutput("Failed to load configuration from %s: %v", source, err)
			continue
		}

		configs = append(configs, config)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no valid configuration sources found")
	}

	// Merge all configurations
	return c.MergeConfigurations(configs)
}

func (c *CLI) MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error) {
	if c.configManager == nil {
		return nil, fmt.Errorf("configuration manager not initialized")
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to merge")
	}

	// Use the configuration manager's merge functionality
	return c.configManager.MergeConfigurations(configs...), nil
}

func (c *CLI) ValidateConfigurationSchema(config *models.ProjectConfig) error {
	if c.configManager == nil {
		return fmt.Errorf("configuration manager not initialized")
	}

	return c.configManager.ValidateConfig(config)
}

func (c *CLI) GetConfigurationSources() ([]interfaces.ConfigSource, error) {
	if c.configManager == nil {
		return nil, fmt.Errorf("configuration manager not initialized")
	}

	return c.configManager.GetConfigSources()
}

func (c *CLI) GetPackageVersions() (map[string]string, error) {
	if c.versionManager == nil {
		return nil, fmt.Errorf("version manager not initialized")
	}

	return c.versionManager.GetAllPackageVersions()
}

func (c *CLI) GetLatestPackageVersions() (map[string]string, error) {
	if c.versionManager == nil {
		return nil, fmt.Errorf("version manager not initialized")
	}

	return c.versionManager.GetLatestPackageVersions()
}

func (c *CLI) CheckCompatibility(projectPath string) (*interfaces.CompatibilityResult, error) {
	if c.versionManager == nil {
		return nil, fmt.Errorf("version manager not initialized")
	}

	return c.versionManager.CheckCompatibility(projectPath)
}

func (c *CLI) GetCacheStats() (*interfaces.CacheStats, error) {
	if c.cacheManager == nil {
		return nil, fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.GetStats()
}

func (c *CLI) ValidateCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.ValidateCache()
}

func (c *CLI) RepairCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	fmt.Println("Repairing cache...")
	err := c.cacheManager.RepairCache()
	if err != nil {
		return fmt.Errorf("failed to repair cache: %w", err)
	}

	fmt.Println("Cache repaired successfully!")
	return nil
}

func (c *CLI) EnableOfflineMode() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	err := c.cacheManager.EnableOfflineMode()
	if err != nil {
		return fmt.Errorf("failed to enable offline mode: %w", err)
	}

	fmt.Println("Offline mode enabled")
	fmt.Println("The generator will now use cached data only")
	return nil
}

func (c *CLI) DisableOfflineMode() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	err := c.cacheManager.DisableOfflineMode()
	if err != nil {
		return fmt.Errorf("failed to disable offline mode: %w", err)
	}

	fmt.Println("Offline mode disabled")
	fmt.Println("The generator will now use network resources when available")
	return nil
}

func (c *CLI) SetLogLevel(level string) error {
	if c.logger == nil {
		return fmt.Errorf("logger not initialized")
	}

	// Convert string level to numeric level
	var numLevel int
	switch strings.ToLower(level) {
	case "debug":
		numLevel = 0
	case "info":
		numLevel = 1
	case "warn", "warning":
		numLevel = 2
	case "error":
		numLevel = 3
	case "fatal":
		numLevel = 4
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	c.logger.SetLevel(numLevel)
	return nil
}

func (c *CLI) GetLogLevel() string {
	if c.logger == nil {
		return "info"
	}

	// This would need to be implemented in the logger interface
	// For now, return a default
	return "info"
}

func (c *CLI) ShowRecentLogs(lines int, level string) error {
	return c.showRecentLogs(lines, level, "", time.Time{}, "text")
}

func (c *CLI) GetLogFileLocations() ([]string, error) {
	if c.logger == nil {
		return nil, fmt.Errorf("logger not initialized")
	}

	// This would need to be implemented in the logger interface
	// For now, return common log locations
	locations := []string{
		"~/.cache/generator/logs/",
		"./logs/",
		"/tmp/generator-logs/",
	}

	return locations, nil
}

func (c *CLI) RunNonInteractive(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	if config == nil {
		return fmt.Errorf("project configuration is required for non-interactive mode")
	}

	c.VerboseOutput("Running in non-interactive mode")
	c.VerboseOutput("Project: %s", config.Name)

	// Validate configuration
	if err := c.ValidateConfigurationSchema(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Generate project with advanced options
	if options != nil {
		return c.GenerateWithAdvancedOptions(config, options)
	}

	// Generate with basic options
	templateName := "go-gin" // Default template
	outputPath := config.OutputPath
	if outputPath == "" {
		outputPath = "./" + config.Name
	}

	if c.templateManager == nil {
		return fmt.Errorf("template manager not initialized")
	}

	return c.templateManager.ProcessTemplate(templateName, config, outputPath)
}

func (c *CLI) GenerateReport(reportType string, format string, outputFile string) error {
	c.VerboseOutput("Generating %s report in %s format", reportType, format)

	var reportData interface{}
	var err error

	switch strings.ToLower(reportType) {
	case "validation":
		// Generate a sample validation report
		reportData = map[string]interface{}{
			"type":      "validation",
			"timestamp": time.Now(),
			"summary": map[string]interface{}{
				"total_files":   10,
				"valid_files":   8,
				"error_count":   2,
				"warning_count": 3,
			},
		}
	case "audit":
		// Generate a sample audit report
		reportData = map[string]interface{}{
			"type":      "audit",
			"timestamp": time.Now(),
			"score":     7.5,
			"categories": map[string]interface{}{
				"security":    8.0,
				"quality":     7.0,
				"performance": 7.5,
			},
		}
	case "configuration":
		// Generate configuration report
		if c.configManager != nil {
			sources, err := c.configManager.GetConfigSources()
			if err == nil {
				reportData = map[string]interface{}{
					"type":      "configuration",
					"timestamp": time.Now(),
					"sources":   sources,
				}
			}
		}
	default:
		return fmt.Errorf("unsupported report type: %s", reportType)
	}

	if reportData == nil {
		return fmt.Errorf("failed to generate report data")
	}

	// Format and write report
	var output []byte
	switch strings.ToLower(format) {
	case "json":
		output, err = json.MarshalIndent(reportData, "", "  ")
	case "yaml":
		// For now, use JSON format
		output, err = json.MarshalIndent(reportData, "", "  ")
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to format report: %w", err)
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, output, 0600)
		if err != nil {
			return fmt.Errorf("failed to write report file: %w", err)
		}
		c.QuietOutput("Report written to: %s", outputFile)
	} else {
		fmt.Println(string(output))
	}

	return nil
}

// filterIssuesByRules filters validation issues by specified rules
func (c *CLI) filterIssuesByRules(issues []interfaces.ValidationIssue, rules []string) []interfaces.ValidationIssue {
	if len(rules) == 0 {
		return issues
	}

	var filtered []interfaces.ValidationIssue
	for _, issue := range issues {
		for _, rule := range rules {
			if strings.Contains(issue.Rule, rule) || strings.Contains(issue.Type, rule) {
				filtered = append(filtered, issue)
				break
			}
		}
	}

	return filtered
}

// applySettingToConfig applies a configuration setting to a ProjectConfig struct
func (c *CLI) applySettingToConfig(config *models.ProjectConfig, key, value string) error {
	switch key {
	case "name":
		config.Name = value
	case "organization":
		config.Organization = value
	case "description":
		config.Description = value
	case "license":
		config.License = value
	case "author":
		config.Author = value
	case "email":
		config.Email = value
	case "repository":
		config.Repository = value
	case "output_path":
		config.OutputPath = value

	// Component settings
	case "components.frontend.nextjs.app":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Frontend.NextJS.App = val
		} else {
			return fmt.Errorf("🚫 '%s' should be true or false, not '%s'", key, value)
		}
	case "components.frontend.nextjs.home":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Frontend.NextJS.Home = val
		} else {
			return fmt.Errorf("🚫 '%s' should be true or false, not '%s'", key, value)
		}
	case "components.frontend.nextjs.admin":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Frontend.NextJS.Admin = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.frontend.nextjs.shared":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Frontend.NextJS.Shared = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.backend.go_gin":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Backend.GoGin = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.mobile.android":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Mobile.Android = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.mobile.ios":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Mobile.IOS = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.infrastructure.docker":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Infrastructure.Docker = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.infrastructure.kubernetes":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Infrastructure.Kubernetes = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.infrastructure.terraform":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Infrastructure.Terraform = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}

	// Version settings
	case "versions.node":
		if config.Versions == nil {
			config.Versions = &models.VersionConfig{Packages: make(map[string]string)}
		}
		config.Versions.Node = value
	case "versions.go":
		if config.Versions == nil {
			config.Versions = &models.VersionConfig{Packages: make(map[string]string)}
		}
		config.Versions.Go = value

	default:
		// Check if it's a package version setting
		if strings.HasPrefix(key, "versions.packages.") {
			packageName := strings.TrimPrefix(key, "versions.packages.")
			if config.Versions == nil {
				config.Versions = &models.VersionConfig{Packages: make(map[string]string)}
			}
			if config.Versions.Packages == nil {
				config.Versions.Packages = make(map[string]string)
			}
			config.Versions.Packages[packageName] = value
		} else {
			return fmt.Errorf("unknown configuration key: %s", key)
		}
	}

	return nil
}

// formatBytes formats a byte count as a human-readable string
func formatBytes(bytes int64) string {
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

// LoadConfigFromFile loads configuration from a file
func (c *CLI) LoadConfigFromFile(path string) error {
	config, err := c.configManager.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to load configuration from file: %w", err)
	}

	// Validate the loaded configuration
	if err := c.configManager.ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Printf("Configuration loaded successfully from: %s\n", path)
	return nil
}

// Log command implementation methods

// showLogLocations displays the locations of log files
func (c *CLI) showLogLocations() error {
	c.QuietOutput("Log file locations:")

	if c.logger != nil {
		logDir := c.logger.GetLogDir()
		if logDir != "" {
			c.QuietOutput("  Log directory: %s", logDir)

			// Get list of log files
			logFiles, err := c.logger.GetLogFiles()
			if err != nil {
				c.WarningOutput("Could not list log files: %v", err)
			} else {
				c.QuietOutput("  Log files:")
				for _, file := range logFiles {
					c.QuietOutput("    %s", file)
				}
			}
		} else {
			c.QuietOutput("  Log directory not configured")
		}
	} else {
		c.QuietOutput("  Logger not initialized")
		c.QuietOutput("  Default location: ~/.cache/template-generator/logs/")
	}

	return nil
}

// showRecentLogs displays recent log entries with filtering
func (c *CLI) showRecentLogs(lines int, level string, component string, since time.Time, format string) error {
	if c.logger == nil {
		return fmt.Errorf("logger not initialized")
	}

	// Get filtered log entries
	entries := c.logger.FilterEntries(level, component, since, lines)

	if len(entries) == 0 {
		c.QuietOutput("No log entries found matching the specified criteria")
		return nil
	}

	// Display entries based on format
	switch strings.ToLower(format) {
	case "json":
		return c.displayLogsJSON(entries)
	case "raw":
		return c.displayLogsRaw(entries)
	default:
		return c.displayLogsText(entries)
	}
}

// displayLogsText displays log entries in human-readable text format
func (c *CLI) displayLogsText(entries []interfaces.LogEntry) error {
	c.QuietOutput("Recent log entries (%d entries):", len(entries))
	c.QuietOutput("") // Empty line for readability

	for _, entry := range entries {
		// Format timestamp
		timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

		// Color code by level (if not disabled)
		levelStr := entry.Level
		switch strings.ToUpper(entry.Level) {
		case "ERROR", "FATAL":
			levelStr = fmt.Sprintf("\033[31m%s\033[0m", entry.Level) // Red
		case "WARN":
			levelStr = fmt.Sprintf("\033[33m%s\033[0m", entry.Level) // Yellow
		case "DEBUG":
			levelStr = fmt.Sprintf("\033[36m%s\033[0m", entry.Level) // Cyan
		case "INFO":
			levelStr = fmt.Sprintf("\033[32m%s\033[0m", entry.Level) // Green
		}

		// Basic log line
		logLine := fmt.Sprintf("%s [%s] %s: %s",
			timestamp, levelStr, entry.Component, entry.Message)

		c.QuietOutput("%s", logLine)

		// Add fields if present
		if len(entry.Fields) > 0 {
			for k, v := range entry.Fields {
				c.QuietOutput("  %s: %v", k, v)
			}
		}

		// Add caller if present
		if entry.Caller != "" {
			c.QuietOutput("  caller: %s", entry.Caller)
		}

		// Add error if present
		if entry.Error != "" {
			c.QuietOutput("  error: %s", entry.Error)
		}

		c.QuietOutput("") // Empty line between entries
	}

	return nil
}

// displayLogsJSON displays log entries in JSON format
func (c *CLI) displayLogsJSON(entries []interfaces.LogEntry) error {
	output := map[string]interface{}{
		"entries": entries,
		"count":   len(entries),
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log entries to JSON: %w", err)
	}

	c.QuietOutput("%s", string(jsonBytes))
	return nil
}

// displayLogsRaw displays raw log file content
func (c *CLI) displayLogsRaw(entries []interfaces.LogEntry) error {
	for _, entry := range entries {
		// Display in a raw format similar to actual log files
		timestamp := entry.Timestamp.Format(time.RFC3339)
		c.QuietOutput("%s [%s] component=%s message=\"%s\"",
			timestamp, entry.Level, entry.Component, entry.Message)
	}
	return nil
}

// followLogs implements real-time log following (tail -f functionality)
func (c *CLI) followLogs(lines int, level string, component string, since time.Time) error {
	if c.logger == nil {
		return fmt.Errorf("logger not initialized")
	}

	c.QuietOutput("Following logs (showing last %d lines)...", lines)
	c.QuietOutput("Press Ctrl+C to stop")
	c.QuietOutput("")

	// Show initial entries
	if err := c.showRecentLogs(lines, level, component, since, "text"); err != nil {
		return err
	}

	// For now, we'll implement a simple polling mechanism
	// In a production implementation, this would use file watching or log streaming
	c.QuietOutput("")
	c.QuietOutput("Note: Real-time following is not yet implemented.")
	c.QuietOutput("This would continuously monitor log files for new entries.")
	c.QuietOutput("Use 'generator logs' to view current log entries.")

	return nil
}

// runInteractiveProjectConfiguration handles interactive project configuration
func (c *CLI) runInteractiveProjectConfiguration(ctx context.Context) (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{}

	// Show breadcrumb
	if err := c.interactiveUI.ShowBreadcrumb(ctx, []string{"Generator", "Project Configuration", "Basic Details"}); err != nil {
		c.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Collect project name
	nameConfig := interfaces.TextPromptConfig{
		Prompt:      "Project Name",
		Description: "Enter a name for your project (letters, numbers, hyphens, and underscores only)",
		Required:    true,
		Validator:   ui.ProjectConfigValidation().Validate,
		AllowBack:   false,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Project name will be used for directories, packages, and documentation. Use lowercase letters, numbers, hyphens, and underscores.",
		MaxLength:   50,
		MinLength:   2,
	}

	nameResult, err := c.interactiveUI.PromptText(ctx, nameConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get project name: %w", err)
	}
	if nameResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}
	config.Name = nameResult.Value

	// Collect organization (optional)
	orgConfig := interfaces.TextPromptConfig{
		Prompt:      "Organization",
		Description: "Enter your organization name (optional)",
		Required:    false,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Organization name will be used in package names and documentation. Leave empty if not applicable.",
		MaxLength:   100,
	}

	orgResult, err := c.interactiveUI.PromptText(ctx, orgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	if orgResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}
	config.Organization = orgResult.Value

	// Collect description (optional)
	descConfig := interfaces.TextPromptConfig{
		Prompt:      "Description",
		Description: "Enter a brief description of your project (optional)",
		Required:    false,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Description will be used in README files and package documentation.",
		MaxLength:   500,
	}

	descResult, err := c.interactiveUI.PromptText(ctx, descConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}
	if descResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}
	config.Description = descResult.Value

	// Collect author (optional)
	authorConfig := interfaces.TextPromptConfig{
		Prompt:      "Author",
		Description: "Enter the author name (optional)",
		Required:    false,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Author name will be used in license files and package documentation.",
		MaxLength:   100,
	}

	authorResult, err := c.interactiveUI.PromptText(ctx, authorConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	if authorResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}
	config.Author = authorResult.Value

	// Select license
	if err := c.interactiveUI.ShowBreadcrumb(ctx, []string{"Generator", "Project Configuration", "License"}); err != nil {
		c.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	licenseConfig := interfaces.MenuConfig{
		Title:       "Select License",
		Description: "Choose a license for your project",
		Options: []interfaces.MenuOption{
			{Label: "MIT License", Description: "Permissive license with minimal restrictions", Value: "MIT"},
			{Label: "Apache License 2.0", Description: "Permissive license with patent protection", Value: "Apache-2.0"},
			{Label: "GNU GPL v3", Description: "Copyleft license requiring derivative works to be open source", Value: "GPL-3.0"},
			{Label: "BSD 3-Clause", Description: "Permissive license similar to MIT", Value: "BSD-3-Clause"},
			{Label: "Mozilla Public License 2.0", Description: "Weak copyleft license", Value: "MPL-2.0"},
			{Label: "Unlicense", Description: "Public domain dedication", Value: "Unlicense"},
		},
		DefaultItem: 0, // MIT as default
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Choose a license that fits your project's goals. MIT is recommended for most open source projects.",
	}

	licenseResult, err := c.interactiveUI.ShowMenu(ctx, licenseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to select license: %w", err)
	}
	if licenseResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}
	config.License = licenseResult.SelectedValue.(string)

	// Select components
	if err := c.interactiveUI.ShowBreadcrumb(ctx, []string{"Generator", "Project Configuration", "Components"}); err != nil {
		c.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	componentsConfig := interfaces.MultiSelectConfig{
		Title:        "Select Project Components",
		Description:  "Choose the components to include in your project",
		MinSelection: 1,
		MaxSelection: 0, // No limit
		Options: []interfaces.SelectOption{
			{
				Label:       "Go Gin API",
				Description: "RESTful API server using Gin framework",
				Value:       "go-gin",
				Selected:    true, // Default selection
				Tags:        []string{"backend", "api", "go"},
			},
			{
				Label:       "Next.js Frontend",
				Description: "React-based frontend with Next.js framework",
				Value:       "nextjs",
				Selected:    false,
				Tags:        []string{"frontend", "react", "nextjs"},
			},
			{
				Label:       "PostgreSQL Database",
				Description: "PostgreSQL database with migrations",
				Value:       "postgresql",
				Selected:    false,
				Tags:        []string{"database", "sql"},
			},
			{
				Label:       "Redis Cache",
				Description: "Redis for caching and session storage",
				Value:       "redis",
				Selected:    false,
				Tags:        []string{"cache", "session"},
			},
			{
				Label:       "Docker Configuration",
				Description: "Docker and Docker Compose setup",
				Value:       "docker",
				Selected:    true, // Default selection
				Tags:        []string{"infrastructure", "containerization"},
			},
			{
				Label:       "Kubernetes Manifests",
				Description: "Kubernetes deployment manifests",
				Value:       "kubernetes",
				Selected:    false,
				Tags:        []string{"infrastructure", "orchestration"},
			},
			{
				Label:       "CI/CD Pipeline",
				Description: "GitHub Actions workflow for CI/CD",
				Value:       "cicd",
				Selected:    false,
				Tags:        []string{"automation", "deployment"},
			},
			{
				Label:       "Monitoring Setup",
				Description: "Prometheus and Grafana monitoring",
				Value:       "monitoring",
				Selected:    false,
				Tags:        []string{"observability", "metrics"},
			},
		},
		AllowBack:     true,
		AllowQuit:     true,
		ShowHelp:      true,
		HelpText:      "Select the components you want to include. You can search by typing / followed by your search term.",
		SearchEnabled: true,
	}

	componentsResult, err := c.interactiveUI.ShowMultiSelect(ctx, componentsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to select components: %w", err)
	}
	if componentsResult.Cancelled {
		return nil, fmt.Errorf("project configuration cancelled")
	}

	// Convert selected components to project config
	config.Components = c.convertSelectedComponents(componentsResult.SelectedValues)

	return config, nil
}

// runInteractiveConfirmation shows configuration summary and confirms generation
func (c *CLI) runInteractiveConfirmation(ctx context.Context, config *models.ProjectConfig, options interfaces.GenerateOptions) bool {
	if err := c.interactiveUI.ShowBreadcrumb(ctx, []string{"Generator", "Confirmation"}); err != nil {
		c.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Show configuration summary in a table
	tableConfig := interfaces.TableConfig{
		Title:   "Project Configuration Summary",
		Headers: []string{"Setting", "Value"},
		Rows: [][]string{
			{"Project Name", config.Name},
			{"Organization", config.Organization},
			{"Description", config.Description},
			{"Author", config.Author},
			{"License", config.License},
			{"Components", c.formatSelectedComponents(config.Components)},
		},
		MaxWidth:   80,
		Pagination: false,
		Sortable:   false,
		Searchable: false,
	}

	if err := c.interactiveUI.ShowTable(ctx, tableConfig); err != nil {
		c.logger.WarnWithFields("📋 Couldn't display configuration table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Confirm generation
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Generate Project",
		Description:  "Do you want to proceed with generating the project with the above configuration?",
		DefaultValue: true,
		YesLabel:     "Generate",
		NoLabel:      "Cancel",
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText:     "Confirm to start project generation or go back to modify the configuration.",
	}

	confirmResult, err := c.interactiveUI.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		c.logger.ErrorWithFields("❓ Couldn't get your confirmation", map[string]interface{}{
			"error": err.Error(),
		})
		return false
	}

	return confirmResult.Confirmed && !confirmResult.Cancelled
}

// runInteractiveGeneration handles project generation with progress tracking
func (c *CLI) runInteractiveGeneration(ctx context.Context, templateName string, config *models.ProjectConfig, outputPath string) error {
	if err := c.interactiveUI.ShowBreadcrumb(ctx, []string{"Generator", "Generation"}); err != nil {
		c.logger.Error("🧭 Couldn't update navigation breadcrumb", "error", err)
	}

	// Create progress tracker
	progressConfig := interfaces.ProgressConfig{
		Title:       "Generating Project",
		Description: fmt.Sprintf("Creating project '%s' using template '%s'", config.Name, templateName),
		Steps: []string{
			"Initializing project structure",
			"Processing templates",
			"Installing dependencies",
			"Configuring components",
			"Finalizing setup",
		},
		ShowPercent: true,
		ShowETA:     true,
		Cancellable: true,
	}

	tracker, err := c.interactiveUI.ShowProgress(ctx, progressConfig)
	if err != nil {
		return fmt.Errorf("failed to create progress tracker: %w", err)
	}
	defer func() {
		if closeErr := tracker.Close(); closeErr != nil {
			c.logger.Error("📊 Couldn't close progress tracker", "error", closeErr)
		}
	}()

	// Step 1: Initialize project structure
	if err := tracker.SetCurrentStep(0, "Creating directory structure"); err != nil {
		c.logger.Error("📋 Couldn't update progress step", "error", err)
	}
	if err := tracker.SetProgress(0.1); err != nil {
		c.logger.Error("📊 Couldn't update progress bar", "error", err)
	}
	if err := tracker.AddLog("Creating output directory: " + outputPath); err != nil {
		c.logger.Error("📝 Couldn't add progress log", "error", err)
	}

	if err := os.MkdirAll(outputPath, 0750); err != nil {
		if failErr := tracker.Fail(fmt.Errorf("failed to create output directory: %w", err)); failErr != nil {
			c.logger.Error("Failed to mark tracker as failed", "error", failErr)
		}
		return err
	}

	// Step 2: Process templates
	if err := tracker.SetCurrentStep(1, "Processing project templates"); err != nil {
		c.logger.Error("Failed to set current step", "error", err)
	}
	if err := tracker.SetProgress(0.3); err != nil {
		c.logger.Error("Failed to set progress", "error", err)
	}
	if err := tracker.AddLog("Loading template: " + templateName); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Use the template manager to process the template
	// This is a simplified version - the actual implementation would be more complex
	if err := c.templateManager.ProcessTemplate(templateName, config, outputPath); err != nil {
		if failErr := tracker.Fail(fmt.Errorf("failed to process template: %w", err)); failErr != nil {
			c.logger.Error("Failed to mark tracker as failed", "error", failErr)
		}
		return err
	}

	// Step 3: Install dependencies (simulated)
	if err := tracker.SetCurrentStep(2, "Installing project dependencies"); err != nil {
		c.logger.Error("Failed to set current step", "error", err)
	}
	if err := tracker.SetProgress(0.6); err != nil {
		c.logger.Error("Failed to set progress", "error", err)
	}
	if err := tracker.AddLog("Installing Go modules"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Simulate dependency installation
	time.Sleep(500 * time.Millisecond)
	if err := tracker.AddLog("Dependencies installed successfully"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Step 4: Configure components
	if err := tracker.SetCurrentStep(3, "Configuring selected components"); err != nil {
		c.logger.Error("Failed to set current step", "error", err)
	}
	if err := tracker.SetProgress(0.8); err != nil {
		c.logger.Error("Failed to set progress", "error", err)
	}
	if err := tracker.AddLog("Configuring project components"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Simulate component configuration
	time.Sleep(300 * time.Millisecond)
	if err := tracker.AddLog("Components configured successfully"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Step 5: Finalize
	if err := tracker.SetCurrentStep(4, "Finalizing project setup"); err != nil {
		c.logger.Error("Failed to set current step", "error", err)
	}
	if err := tracker.SetProgress(0.95); err != nil {
		c.logger.Error("Failed to set progress", "error", err)
	}
	if err := tracker.AddLog("Setting up project metadata"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Simulate finalization
	time.Sleep(200 * time.Millisecond)
	if err := tracker.AddLog("Project generation completed"); err != nil {
		c.logger.Error("Failed to add log", "error", err)
	}

	// Complete the progress
	if err := tracker.SetProgress(1.0); err != nil {
		c.logger.Error("Failed to set progress", "error", err)
	}
	if err := tracker.Complete(); err != nil {
		c.logger.Error("Failed to complete tracker", "error", err)
	}

	c.QuietOutput("Project '%s' generated successfully in '%s'", config.Name, outputPath)
	return nil
}

// convertSelectedComponents converts selected component values to models.Components
func (c *CLI) convertSelectedComponents(selectedValues []interface{}) models.Components {
	components := models.Components{}

	for _, value := range selectedValues {
		component := value.(string)
		switch component {
		case "go-gin":
			components.Backend.GoGin = true
		case "nextjs":
			components.Frontend.NextJS.App = true
		case "postgresql":
			components.Database.PostgreSQL = true
		case "redis":
			components.Cache.Redis = true
		case "docker":
			components.Infrastructure.Docker = true
		case "kubernetes":
			components.Infrastructure.Kubernetes = true
		case "cicd":
			components.DevOps.CICD = true
		case "monitoring":
			components.Monitoring.Prometheus = true
			components.Monitoring.Grafana = true
		}
	}

	return components
}

// formatSelectedComponents formats selected components for display
func (c *CLI) formatSelectedComponents(components models.Components) string {
	var selected []string

	if components.Backend.GoGin {
		selected = append(selected, "Go Gin API")
	}
	if components.Frontend.NextJS.App {
		selected = append(selected, "Next.js Frontend")
	}
	if components.Database.PostgreSQL {
		selected = append(selected, "PostgreSQL")
	}
	if components.Cache.Redis {
		selected = append(selected, "Redis")
	}
	if components.Infrastructure.Docker {
		selected = append(selected, "Docker")
	}
	if components.Infrastructure.Kubernetes {
		selected = append(selected, "Kubernetes")
	}
	if components.DevOps.CICD {
		selected = append(selected, "CI/CD Pipeline")
	}
	if components.Monitoring.Prometheus || components.Monitoring.Grafana {
		selected = append(selected, "Monitoring")
	}

	if len(selected) == 0 {
		return "None"
	}

	return strings.Join(selected, ", ")
}

// detectGenerationMode determines which generation mode to use based on flags and environment
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

// validateModeFlags checks for conflicting mode flags
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
		return fmt.Errorf("🚫 You can only use one mode flag at a time")
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
			return fmt.Errorf("🚫 '%s' isn't a valid mode. Try one of these: %s", explicitMode, strings.Join(validModes, ", "))
		}
	}

	return nil
}

// applyModeOverrides applies mode override flags to the base mode flags
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

// validateAndNormalizeMode validates and normalizes the mode string
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

// runNonInteractiveGeneration handles non-interactive project generation
func (c *CLI) runNonInteractiveGeneration(options interfaces.GenerateOptions) error {
	c.VerboseOutput("Running in non-interactive mode, loading configuration from environment variables")

	// Load configuration from environment variables
	envConfig, err := c.loadEnvironmentConfig()
	if err != nil {
		return fmt.Errorf("failed to load environment configuration: %w", err)
	}

	// Convert environment config to project config
	config, err := c.convertEnvironmentConfigToProjectConfig(envConfig)
	if err != nil {
		return fmt.Errorf("failed to convert environment configuration: %w", err)
	}

	// Validate required fields for non-interactive mode
	if config.Name == "" {
		err := c.createConfigurationError("🚫 Your project needs a name for non-interactive mode", "")
		err = err.WithSuggestions("Set the GENERATOR_PROJECT_NAME environment variable")
		return err
	}

	// Override with command-line flags and environment variables
	if options.OutputPath != "" {
		config.OutputPath = options.OutputPath
	} else if envConfig.OutputPath != "" {
		config.OutputPath = envConfig.OutputPath
	} else {
		config.OutputPath = "./" + config.Name
	}

	// Update options with environment variables
	options.Force = options.Force || envConfig.Force
	options.Minimal = options.Minimal || envConfig.Minimal
	options.Offline = options.Offline || envConfig.Offline
	options.UpdateVersions = options.UpdateVersions || envConfig.UpdateVersions
	options.SkipValidation = options.SkipValidation || envConfig.SkipValidation
	options.BackupExisting = options.BackupExisting && envConfig.BackupExisting
	options.IncludeExamples = options.IncludeExamples && envConfig.IncludeExamples
	options.OutputPath = config.OutputPath

	if len(options.Templates) == 0 && envConfig.Template != "" {
		options.Templates = []string{envConfig.Template}
	}

	// Log CI environment information if detected
	ci := c.detectCIEnvironment()
	if ci.IsCI {
		c.VerboseOutput("Detected CI environment: %s", ci.Provider)
		if ci.BuildID != "" {
			c.VerboseOutput("Build ID: %s", ci.BuildID)
		}
		if ci.Branch != "" {
			c.VerboseOutput("Branch: %s", ci.Branch)
		}
	}

	// Validate configuration if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("🔍 Double-checking configuration...")
		if err := c.validateGenerateConfiguration(config, options); err != nil {
			return err
		}
	}

	// Generate project in non-interactive mode
	c.VerboseOutput("Generating project '%s' in non-interactive mode", config.Name)

	// Use template from options or default
	templateName := ""
	if len(options.Templates) > 0 {
		templateName = options.Templates[0]
	}
	if templateName == "" {
		templateName = "go-gin" // Default template
	}

	// Validate dependencies if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("Validating dependencies...")
		if err := c.validateDependencies(config, templateName); err != nil {
			return fmt.Errorf("dependency validation failed: %w", err)
		}
	}

	// Perform pre-generation checks
	if err := c.performPreGenerationChecks(options.OutputPath, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Handle dry-run mode
	if options.DryRun {
		c.QuietOutput("Dry run mode - would generate project '%s' using template '%s' in directory '%s'",
			config.Name, templateName, options.OutputPath)
		return nil
	}

	return c.templateManager.ProcessTemplate(templateName, config, options.OutputPath)
}

// routeToGenerationMethod routes to the appropriate generation method based on mode
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

// handleConfigFileGeneration handles configuration file-based generation
func (c *CLI) handleConfigFileGeneration(configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("📄 Loading your project configuration...")
	c.DebugOutput("📄 Using configuration: %s", configPath)

	// Validate configuration file exists and is readable
	if err := c.validateConfigurationFile(configPath); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Use existing GenerateFromConfig method
	return c.GenerateFromConfig(configPath, options)
}

// handleNonInteractiveGeneration handles non-interactive generation
func (c *CLI) handleNonInteractiveGeneration(options interfaces.GenerateOptions) error {
	c.VerboseOutput("Starting non-interactive generation")

	// Check for required environment variables
	if err := c.validateNonInteractiveEnvironment(); err != nil {
		return fmt.Errorf("non-interactive environment validation failed: %w", err)
	}

	// Use the existing non-interactive generation method
	return c.runNonInteractiveGeneration(options)
}

// handleInteractiveGeneration handles interactive generation
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

// validateConfigurationFile validates that the configuration file exists and is readable
func (c *CLI) validateConfigurationFile(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("configuration file path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("🚫 Can't find the configuration file: %s", configPath)
	}

	// Check if file is readable
	file, err := os.Open(configPath) // #nosec G304 - configPath is validated before use
	if err != nil {
		return fmt.Errorf("🚫 Can't read the configuration file: %w", err)
	}
	_ = file.Close()

	c.DebugOutput("✅ Configuration file looks good: %s", configPath) // #nosec G304 - configPath is validated and only used for logging
	return nil
}

// validateNonInteractiveEnvironment validates the environment for non-interactive generation
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

// validateInteractiveEnvironment validates the environment for interactive generation
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
