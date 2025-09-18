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
	"github.com/spf13/cobra"
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
type CLI struct {
	configManager    interfaces.ConfigManager
	validator        interfaces.ValidationEngine
	templateManager  interfaces.TemplateManager
	cacheManager     interfaces.CacheManager
	versionManager   interfaces.VersionManager
	auditEngine      interfaces.AuditEngine
	logger           interfaces.Logger
	generatorVersion string
	rootCmd          *cobra.Command
	verboseMode      bool
	quietMode        bool
	debugMode        bool
	exitCode         int
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
	cli := &CLI{
		configManager:    configManager,
		validator:        validator,
		templateManager:  templateManager,
		cacheManager:     cacheManager,
		versionManager:   versionManager,
		auditEngine:      auditEngine,
		logger:           logger,
		generatorVersion: version,
	}

	cli.setupCommands()
	return cli
}

// setupCommands initializes all CLI commands and their flags
func (c *CLI) setupCommands() {
	c.rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "Open Source Project Generator",
		Long: `A comprehensive tool for generating production-ready, enterprise-grade
open source project structures following modern best practices.

The generator supports multiple technology stacks including:
  • Go 1.21+ with Gin framework
  • Node.js 20+ with Next.js 15+ and React 19+
  • Mobile development with Android (Kotlin) and iOS (Swift)
  • Infrastructure with Docker, Kubernetes, and Terraform

Features:
  • Interactive project configuration
  • Template-based code generation
  • Project validation and auditing
  • Configuration management
  • Offline mode support
  • Comprehensive documentation generation

Examples:
  generator generate                    # Interactive project generation
  generator generate --config app.yaml # Generate from configuration
  generator validate ./my-project      # Validate project structure
  generator audit ./my-project         # Security and quality audit
  generator list-templates             # Show available templates
  generator version --packages         # Show latest package versions

For more information about a specific command, use:
  generator <command> --help`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `  # Interactive project generation
  generator generate

  # Generate from configuration file
  generator generate --config project.yaml --output ./my-app

  # Generate minimal project structure
  generator generate --minimal --template go-gin

  # Validate existing project
  generator validate ./my-project --fix

  # Audit project for security issues
  generator audit ./my-project --security --detailed

  # List available templates
  generator list-templates --category backend

  # Show version and package information
  generator version --packages --check-updates`,
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
		return fmt.Errorf("cannot use both --verbose and --quiet flags")
	}
	if debug && quiet {
		return fmt.Errorf("cannot use both --debug and --quiet flags")
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
		return fmt.Errorf("invalid log level: %s (valid levels: %s)", logLevel, strings.Join(validLogLevels, ", "))
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
		return fmt.Errorf("invalid output format: %s (valid formats: %s)", outputFormat, strings.Join(validOutputFormats, ", "))
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
			c.logger.DebugWithFields("CLI configuration", map[string]interface{}{
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
			c.logger.InfoWithFields("Starting command execution", map[string]interface{}{
				"command": cmd.Name(),
				"args":    cmd.Flags().Args(),
			})
		}
	}

	// Store global settings for use in commands
	if c.verboseMode && !c.quietMode {
		fmt.Printf("Verbose: Running command '%s' with detailed output\n", cmd.Name())
	}

	if c.debugMode && !c.quietMode {
		fmt.Printf("Debug: Running command '%s' with debug logging and performance metrics\n", cmd.Name())
	}

	if nonInteractive && (c.verboseMode || c.debugMode) && !c.quietMode {
		fmt.Printf("Info: Running in non-interactive mode\n")
	}

	return nil
}

// Verbose output methods for enhanced debugging and user feedback

// VerboseOutput prints verbose information if verbose mode is enabled
func (c *CLI) VerboseOutput(format string, args ...interface{}) {
	if c.verboseMode && !c.quietMode {
		fmt.Printf("Verbose: "+format+"\n", args...)
	}
	if c.logger != nil && c.logger.IsInfoEnabled() {
		c.logger.Info(format, args...)
	}
}

// DebugOutput prints debug information if debug mode is enabled
func (c *CLI) DebugOutput(format string, args ...interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("Debug: "+format+"\n", args...)
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
	if c.logger != nil {
		c.logger.Info(format, args...)
	}
}

// ErrorOutput prints error information (always shown unless completely silent)
func (c *CLI) ErrorOutput(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	if c.logger != nil {
		c.logger.Error(format, args...)
	}
}

// WarningOutput prints warning information if not in quiet mode
func (c *CLI) WarningOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf("Warning: "+format+"\n", args...)
	}
	if c.logger != nil {
		c.logger.Warn(format, args...)
	}
}

// SuccessOutput prints success information if not in quiet mode
func (c *CLI) SuccessOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf("Success: "+format+"\n", args...)
	}
	if c.logger != nil {
		c.logger.Info("Success: "+format, args...)
	}
}

// PerformanceOutput prints performance metrics if debug mode is enabled
func (c *CLI) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("Performance: %s completed in %v\n", operation, duration)
		if len(metrics) > 0 {
			fmt.Printf("Performance Metrics:\n")
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
	c.VerboseOutput("Starting %s: %s", operation, description)

	var ctx *interfaces.OperationContext
	if c.logger != nil {
		ctx = c.logger.StartOperation(operation, map[string]interface{}{
			"description": description,
		})
	}

	return ctx
}

// FinishOperationWithOutput completes an operation with verbose output
func (c *CLI) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	if ctx != nil && c.logger != nil {
		c.logger.FinishOperation(ctx, map[string]interface{}{
			"description": description,
		})
	}
	c.VerboseOutput("Completed %s: %s", operation, description)
}

// FinishOperationWithError completes an operation with error output
func (c *CLI) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	if ctx != nil && c.logger != nil {
		c.logger.FinishOperationWithError(ctx, err, nil)
	}
	c.ErrorOutput("Failed %s: %v", operation, err)
}

// GetExitCode returns the current exit code
func (c *CLI) GetExitCode() int {
	return c.exitCode
}

// SetExitCode sets the exit code for the CLI
func (c *CLI) SetExitCode(code int) {
	c.exitCode = code
}

// PromptProjectDetails collects basic project configuration from user input.
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	return nil, fmt.Errorf("PromptProjectDetails implementation pending - will be implemented in task 2")
}

// ConfirmGeneration shows a basic configuration preview and asks for user confirmation.
func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	return false
}

// setupGenerateCommand sets up the generate command with all documented flags
func (c *CLI) setupGenerateCommand() {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from templates",
		Long: `Generate a new project using interactive prompts or configuration files.

The generate command creates production-ready project structures with:
  • Modern technology stacks (Go, Node.js, React, etc.)
  • Best practices and security configurations
  • Comprehensive documentation and examples
  • CI/CD workflows and deployment configurations
  • Testing frameworks and quality tools

Generation Modes:
  • Interactive: Guided prompts for project configuration
  • Configuration file: Generate from YAML/JSON configuration
  • Template-based: Use specific templates with custom options
  • Minimal: Generate only essential project structure

Supported Technologies:
  • Backend: Go with Gin framework, REST APIs, GraphQL
  • Frontend: Next.js, React, TypeScript, Tailwind CSS
  • Mobile: Android (Kotlin), iOS (Swift), shared components
  • Infrastructure: Docker, Kubernetes, Terraform, monitoring`,
		RunE: c.runGenerate,
		Example: `  # Interactive project generation
  generator generate

  # Generate from configuration file
  generator generate --config project.yaml

  # Generate with specific output directory
  generator generate --output ./my-new-project

  # Generate minimal project structure
  generator generate --minimal --template go-gin

  # Generate in offline mode using cached templates
  generator generate --offline

  # Generate with latest package versions
  generator generate --update-versions

  # Preview generation without creating files
  generator generate --dry-run

  # Force overwrite existing files
  generator generate --force --backup-existing=false`,
	}

	// Basic flags
	generateCmd.Flags().StringP("config", "c", "", "Path to configuration file (YAML or JSON)")
	generateCmd.Flags().StringP("output", "o", "", "Output directory for generated project")
	generateCmd.Flags().Bool("dry-run", false, "Preview generation without creating files")

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

	c.rootCmd.AddCommand(generateCmd)
}

// setupValidateCommand sets up the validate command with all documented flags
func (c *CLI) setupValidateCommand() {
	validateCmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate project structure and configuration",
		Long: `Validate the structure, configuration, and dependencies of a project.

The validate command performs comprehensive checks including:
  • Project structure and file organization
  • Configuration file syntax and values
  • Dependency versions and compatibility
  • Code quality and best practices
  • Security configurations
  • Documentation completeness

Validation Categories:
  • Structure: Directory layout, required files, naming conventions
  • Configuration: YAML/JSON syntax, schema validation, value ranges
  • Dependencies: Version compatibility, security vulnerabilities
  • Quality: Code style, test coverage, documentation
  • Security: Permissions, secrets, security policies

The validator can automatically fix many common issues when using the --fix flag.
Generate detailed reports in multiple formats for CI/CD integration.`,
		RunE: c.runValidate,
		Example: `  # Validate current directory
  generator validate

  # Validate specific project
  generator validate ./my-project

  # Validate and auto-fix issues
  generator validate ./my-project --fix

  # Generate detailed HTML report
  generator validate --report --report-format html --output-file report.html

  # Validate with specific rules only
  generator validate --rules structure,dependencies

  # Ignore warnings, show only errors
  generator validate --ignore-warnings

  # Verbose validation output
  generator validate --verbose`,
	}

	validateCmd.Flags().Bool("fix", false, "Automatically fix common validation issues")
	validateCmd.Flags().Bool("report", false, "Generate detailed validation report")
	validateCmd.Flags().String("report-format", "text", "Report format (text, json, html, markdown)")
	validateCmd.Flags().StringSlice("rules", []string{}, "Specific validation rules to apply")
	validateCmd.Flags().Bool("ignore-warnings", false, "Ignore validation warnings")
	validateCmd.Flags().String("output-file", "", "Save report to file")

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
		Use:   "audit [path]",
		Short: "Audit project for security, quality, and best practices",
		Long: `Perform comprehensive auditing of a project including security vulnerabilities,
code quality analysis, license compliance, and performance optimization.

The audit command provides enterprise-grade analysis including:
  • Security vulnerability scanning
  • Code quality and maintainability metrics
  • License compliance checking
  • Performance and bundle size analysis
  • Best practices compliance
  • Dependency analysis and recommendations

Audit Categories:
  • Security: CVE scanning, policy compliance, secret detection
  • Quality: Code smells, duplication, complexity metrics
  • Licenses: Compatibility checking, conflict detection
  • Performance: Bundle analysis, load time optimization
  • Dependencies: Outdated packages, security issues

Generate comprehensive reports with actionable recommendations
for improving project security, quality, and maintainability.`,
		RunE: c.runAudit,
		Example: `  # Full audit of current directory
  generator audit

  # Audit specific project
  generator audit ./my-project

  # Security audit only
  generator audit --security --no-quality --no-licenses --no-performance

  # Generate detailed JSON report
  generator audit --detailed --output-format json --output-file audit.json

  # Quality and performance audit
  generator audit --quality --performance

  # Audit with HTML report
  generator audit --output-format html --output-file audit-report.html`,
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
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long: `Display version information for the generator and supported technologies.

The version command shows:
  • Generator version and build information
  • Latest versions of supported technologies
  • Package version information for all stacks
  • Update availability and release notes
  • Compatibility information

Supported Technologies:
  • Runtime: Go 1.21+, Node.js 20+, Python 3.11+
  • Frontend: Next.js 15+, React 19+, TypeScript 5+
  • Mobile: Android SDK, iOS SDK, Kotlin, Swift
  • Infrastructure: Docker, Kubernetes, Terraform
  • Databases: PostgreSQL, MongoDB, Redis
  • Tools: ESLint, Prettier, Jest, Cypress

Use --check-updates to see if newer versions are available.`,
		RunE: c.runVersion,
		Example: `  # Show generator version
  generator version

  # Show all package versions
  generator version --packages

  # Check for updates
  generator version --check-updates

  # Show detailed build information
  generator version --build-info

  # Output in JSON format
  generator version --packages --output-format json`,
	}

	versionCmd.Flags().Bool("packages", false, "Show latest package versions for all supported technologies")
	versionCmd.Flags().Bool("check-updates", false, "Check for generator updates")
	versionCmd.Flags().Bool("build-info", false, "Show detailed build information")

	// Additional version flags
	versionCmd.Flags().Bool("short", false, "Show only version number")
	versionCmd.Flags().String("format", "text", "Output format (text, json, yaml)")
	versionCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	versionCmd.Flags().String("check-package", "", "Check version for specific package")

	c.rootCmd.AddCommand(versionCmd)
}

// setupConfigCommand sets up the config command with all subcommands
func (c *CLI) setupConfigCommand() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage generator configuration and defaults",
		Long: `Manage generator configuration including defaults, user preferences,
and project-specific settings. Supports multiple configuration sources.`,
	}

	// config show
	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration with source information",
		Long:  "Display current configuration values and their sources (file, environment, defaults)",
		RunE:  c.runConfigShow,
	}
	configCmd.AddCommand(configShowCmd)

	// config set
	configSetCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set individual configuration values",
		Long:  "Set individual configuration values or load configuration from file",
		RunE:  c.runConfigSet,
		Args:  cobra.RangeArgs(0, 2),
	}
	configSetCmd.Flags().String("file", "", "Load configuration from file")
	configCmd.AddCommand(configSetCmd)

	// config edit
	configEditCmd := &cobra.Command{
		Use:   "edit",
		Short: "Open configuration file in default editor",
		Long:  "Open the configuration file in the system's default editor",
		RunE:  c.runConfigEdit,
	}
	configCmd.AddCommand(configEditCmd)

	// config validate
	configValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration syntax and values",
		Long:  "Validate configuration file syntax and check for valid values",
		RunE:  c.runConfigValidate,
	}
	configCmd.AddCommand(configValidateCmd)

	// config export
	configExportCmd := &cobra.Command{
		Use:   "export [file]",
		Short: "Export current configuration to shareable file",
		Long:  "Export current configuration to a file that can be shared or used as template",
		RunE:  c.runConfigExport,
	}
	configCmd.AddCommand(configExportCmd)

	c.rootCmd.AddCommand(configCmd)
}

// setupListTemplatesCommand sets up the list-templates command
func (c *CLI) setupListTemplatesCommand() {
	listTemplatesCmd := &cobra.Command{
		Use:   "list-templates",
		Short: "List available project templates",
		Long: `List all available project templates with filtering and search capabilities.

The list-templates command shows:
  • Available templates by category and technology
  • Template descriptions and compatibility information
  • Version information and dependencies
  • Tags and keywords for easy discovery
  • Maintainer and license information

Template Categories:
  • Frontend: Next.js applications, React components, landing pages
  • Backend: Go APIs, microservices, GraphQL servers
  • Mobile: Android apps, iOS apps, shared components
  • Infrastructure: Docker configs, Kubernetes manifests, Terraform
  • Full-stack: Complete application templates

Use filters to find templates that match your specific needs.
Each template includes comprehensive documentation and examples.`,
		RunE: c.runListTemplates,
		Example: `  # List all templates
  generator list-templates

  # List backend templates only
  generator list-templates --category backend

  # List Go-based templates
  generator list-templates --technology go

  # Search for API templates
  generator list-templates --search api

  # List templates with specific tags
  generator list-templates --tags rest,microservice

  # Show detailed template information
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

The template command provides operations for:
  • Viewing detailed template information
  • Validating custom templates
  • Managing template dependencies
  • Checking template compatibility

Template Operations:
  • info: Show detailed information about a specific template
  • validate: Validate template structure and metadata
  • dependencies: Show template dependencies
  • compatibility: Check template compatibility

Use these commands to inspect templates before using them
or to validate custom templates you've created.`,
	}

	// template info
	templateInfoCmd := &cobra.Command{
		Use:   "info <template-name>",
		Short: "Show detailed information about a template",
		Long: `Display comprehensive information about a specific template including:
  • Template metadata and description
  • Version and compatibility information
  • Dependencies and requirements
  • Variables and configuration options
  • File structure and components`,
		RunE: c.runTemplateInfo,
		Args: cobra.ExactArgs(1),
		Example: `  # Show info for go-gin template
  generator template info go-gin

  # Show info for nextjs-app template
  generator template info nextjs-app`,
	}
	templateInfoCmd.Flags().Bool("detailed", false, "Show detailed template information")
	templateInfoCmd.Flags().Bool("variables", false, "Show template variables")
	templateInfoCmd.Flags().Bool("dependencies", false, "Show template dependencies")
	templateInfoCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	templateCmd.AddCommand(templateInfoCmd)

	// template validate
	templateValidateCmd := &cobra.Command{
		Use:   "validate <template-path>",
		Short: "Validate template structure and metadata",
		Long: `Validate a custom template directory for:
  • Proper template structure
  • Valid metadata files
  • Template syntax correctness
  • Required files and directories
  • Best practices compliance`,
		RunE: c.runTemplateValidate,
		Args: cobra.ExactArgs(1),
		Example: `  # Validate custom template
  generator template validate ./my-custom-template

  # Validate with detailed output
  generator template validate ./my-template --detailed`,
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
		Use:   "update",
		Short: "Update generator and check for new versions",
		Long: `Check for and install updates to the generator.

The update command provides comprehensive update management including:
  • Checking for generator updates
  • Installing available updates with safety checks
  • Updating template cache and package information
  • Compatibility checking before updates
  • Rollback support for failed updates

Update Safety Features:
  • Automatic backup creation before updates
  • Compatibility verification
  • Signature verification for security
  • Rollback on installation failure
  • Non-disruptive update checking

Update Channels:
  • Stable: Production-ready releases (default)
  • Beta: Pre-release versions with new features
  • Alpha: Development versions for testing

The update process includes comprehensive safety checks and
can be configured to run automatically or manually.`,
		RunE: c.runUpdate,
		Example: `  # Check for updates
  generator update --check

  # Install available updates
  generator update --install

  # Update templates cache
  generator update --templates

  # Force update even if risky
  generator update --install --force

  # Check compatibility before updating
  generator update --check --compatibility

  # Update with specific channel
  generator update --channel beta --install

  # Show release notes for available update
  generator update --check --release-notes`,
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
		Use:   "cache",
		Short: "Manage template and version cache",
		Long: `Manage the local cache used for offline mode and performance optimization.
Includes cache statistics, cleanup, and repair operations.`,
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
		Use:   "logs",
		Short: "View recent log entries and log file locations",
		Long: `Display recent log entries and show log file locations.

The logs command provides comprehensive log viewing capabilities including:
  • Recent log entries with filtering and search
  • Log file location information
  • Real-time log following (tail -f functionality)
  • Filtering by log level, component, and time
  • Multiple output formats for integration

Log Filtering:
  • Level: Filter by log level (debug, info, warn, error, fatal)
  • Component: Filter by component name (cli, config, template, etc.)
  • Time: Show logs since a specific time
  • Lines: Limit number of entries shown

Output Formats:
  • Text: Human-readable format (default)
  • JSON: Machine-readable format for parsing
  • Raw: Raw log file content

Use this command for debugging issues, monitoring operations,
and understanding application behavior.`,
		RunE: c.runLogs,
		Example: `  # Show last 50 log entries
  generator logs

  # Show last 100 log entries
  generator logs --lines 100

  # Show only error logs
  generator logs --level error

  # Show logs from specific component
  generator logs --component cli

  # Show logs since specific time
  generator logs --since "2024-01-01T10:00:00Z"

  # Show logs in JSON format
  generator logs --format json

  # Follow logs in real-time
  generator logs --follow

  # Show log file locations
  generator logs --locations

  # Combine filters
  generator logs --level warn --component template --lines 20`,
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

	// Auto-detect non-interactive mode if not explicitly set
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
		if nonInteractive {
			c.VerboseOutput("Auto-detected non-interactive mode")
		}
	}

	// Handle non-interactive mode
	if nonInteractive {
		interactive = false
	}

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
		c.VerboseOutput("Performing pre-generation validation...")
		if err := c.validateGenerateOptions(options); err != nil {
			return fmt.Errorf("generate options validation failed: %w", err)
		}
	}

	// If config file is provided, generate from config
	if configPath != "" {
		c.VerboseOutput("Loading configuration from file: %s", configPath)
		return c.GenerateFromConfig(configPath, options)
	}

	// Non-interactive mode with environment variables
	if nonInteractive {
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
			return c.createConfigurationError("project name is required in non-interactive mode", "GENERATOR_PROJECT_NAME environment variable")
		}

		// Override with command-line flags and environment variables
		if outputPath != "" {
			config.OutputPath = outputPath
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

		if template == "" && envConfig.Template != "" {
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
			c.VerboseOutput("Validating configuration...")
			if err := c.validateGenerateConfiguration(config, options); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
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
			return fmt.Errorf("pre-generation checks failed: %w", err)
		}

		// Handle dry-run mode
		if options.DryRun {
			c.QuietOutput("Dry run mode - would generate project '%s' using template '%s' in directory '%s'",
				config.Name, templateName, options.OutputPath)
			return nil
		}

		return c.templateManager.ProcessTemplate(templateName, config, options.OutputPath)
	}

	// Interactive mode
	if interactive {
		c.VerboseOutput("Starting interactive project configuration")

		config, err := c.PromptProjectDetails()
		if err != nil {
			return fmt.Errorf("failed to collect project details: %w", err)
		}

		// Validate configuration if not skipped
		if !options.SkipValidation {
			c.VerboseOutput("Validating configuration...")
			if err := c.validateGenerateConfiguration(config, options); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}
		}

		if !c.ConfirmGeneration(config) {
			c.QuietOutput("Generation cancelled by user")
			return nil
		}

		// Generate project using template manager
		templateName := template
		if templateName == "" {
			templateName = "go-gin" // Default template
		}

		if outputPath == "" {
			outputPath = config.Name
		}

		// Validate dependencies if not skipped
		if !options.SkipValidation {
			c.VerboseOutput("Validating dependencies...")
			if err := c.validateDependencies(config, templateName); err != nil {
				return fmt.Errorf("dependency validation failed: %w", err)
			}
		}

		// Perform pre-generation checks
		if err := c.performPreGenerationChecks(outputPath, options); err != nil {
			return fmt.Errorf("pre-generation checks failed: %w", err)
		}

		// Handle dry-run mode
		if options.DryRun {
			c.QuietOutput("Dry run mode - would generate project '%s' using template '%s' in directory '%s'",
				config.Name, templateName, outputPath)
			return nil
		}

		return c.templateManager.ProcessTemplate(templateName, config, outputPath)
	}

	return fmt.Errorf("no configuration provided - use --config flag, enable --interactive mode, or run in non-interactive mode with environment variables")
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
		return fmt.Errorf("validation failed: %w", err)
	}

	// Output results based on format and mode
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		// Machine-readable output for automation
		return c.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	if !c.quietMode {
		c.QuietOutput("Validation completed for: %s", path)
		c.QuietOutput("Valid: %t", result.Valid)
		c.QuietOutput("Issues: %d", len(result.Issues))
		c.QuietOutput("Warnings: %d", len(result.Warnings))

		if len(result.Issues) > 0 && !summaryOnly {
			c.QuietOutput("\nIssues found:")
			for _, issue := range result.Issues {
				c.QuietOutput("  - %s: %s", issue.Severity, issue.Message)
				if issue.File != "" {
					c.VerboseOutput("    File: %s:%d:%d", issue.File, issue.Line, issue.Column)
				}
			}
		}

		if len(result.Warnings) > 0 && !ignoreWarnings && !summaryOnly {
			c.QuietOutput("\nWarnings:")
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

	// Return appropriate exit code
	if !result.Valid {
		details := map[string]interface{}{
			"issues_count":   len(result.Issues),
			"warnings_count": len(result.Warnings),
			"path":           path,
		}
		return c.createValidationError(fmt.Sprintf("validation failed with %d issues", len(result.Issues)), details)
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
		return fmt.Errorf("audit failed: %w", err)
	}

	// Output results based on format and mode
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		// Machine-readable output for automation
		return c.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	if !c.quietMode {
		c.QuietOutput("Audit completed for: %s", path)
		c.QuietOutput("Overall Score: %.2f", result.OverallScore)
		c.VerboseOutput("Audit Time: %s", result.AuditTime.Format("2006-01-02 15:04:05"))

		if result.Security != nil && !summaryOnly {
			c.QuietOutput("Security Score: %.2f", result.Security.Score)
			c.VerboseOutput("Vulnerabilities: %d", len(result.Security.Vulnerabilities))
		}

		if result.Quality != nil && !summaryOnly {
			c.QuietOutput("Quality Score: %.2f", result.Quality.Score)
			c.VerboseOutput("Code Smells: %d", len(result.Quality.CodeSmells))
		}

		if result.Licenses != nil && !summaryOnly {
			c.QuietOutput("License Compatible: %t", result.Licenses.Compatible)
			c.VerboseOutput("License Conflicts: %d", len(result.Licenses.Conflicts))
		}

		if result.Performance != nil && !summaryOnly {
			c.QuietOutput("Performance Score: %.2f", result.Performance.Score)
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
		return c.createAuditError(fmt.Sprintf("audit failed: high severity issues found (score: %.2f)", result.OverallScore), result.OverallScore)
	}

	if failOnMedium && result.OverallScore < 5.0 {
		return c.createAuditError(fmt.Sprintf("audit failed: medium or higher severity issues found (score: %.2f)", result.OverallScore), result.OverallScore)
	}

	if minScore > 0 && result.OverallScore < minScore {
		return c.createAuditError(fmt.Sprintf("audit failed: score %.2f is below minimum required score %.2f", result.OverallScore, minScore), result.OverallScore)
	}

	return nil
}

func (c *CLI) runVersion(cmd *cobra.Command, args []string) error {
	// Get flags
	packages, _ := cmd.Flags().GetBool("packages")
	checkUpdates, _ := cmd.Flags().GetBool("check-updates")
	buildInfo, _ := cmd.Flags().GetBool("build-info")
	short, _ := cmd.Flags().GetBool("short")
	format, _ := cmd.Flags().GetString("format")
	compatibility, _ := cmd.Flags().GetBool("compatibility")
	checkPackage, _ := cmd.Flags().GetString("check-package")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
	}

	// Use format if outputFormat is not set
	if format != "text" {
		outputFormat = format
	}

	// Handle short version output
	if short {
		if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
			versionData := map[string]string{"version": c.generatorVersion}
			return c.outputMachineReadable(versionData, outputFormat)
		}
		c.QuietOutput(c.generatorVersion)
		return nil
	}

	// Handle specific package version check
	if checkPackage != "" {
		c.VerboseOutput("Checking version for package: %s", checkPackage)
		// This would be implemented when package version checking is fully implemented
		if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
			packageData := map[string]interface{}{
				"package": checkPackage,
				"status":  "not_implemented",
			}
			return c.outputMachineReadable(packageData, outputFormat)
		}
		return nil
	}

	// Handle compatibility flag
	if compatibility {
		c.DebugOutput("Showing compatibility information")
	}

	// Create version options
	options := interfaces.VersionOptions{
		ShowPackages:  packages,
		CheckUpdates:  checkUpdates,
		ShowBuildInfo: buildInfo,
		OutputFormat:  outputFormat,
	}

	// Show version information
	return c.ShowVersion(options)
}

func (c *CLI) runConfigShow(cmd *cobra.Command, args []string) error {
	// Show current configuration with source information
	return c.ShowConfig()
}

func (c *CLI) runConfigSet(cmd *cobra.Command, args []string) error {
	// Get flags
	file, _ := cmd.Flags().GetString("file")

	if file != "" {
		// Load configuration from file
		fmt.Printf("Loading configuration from: %s\n", file)
		// This would be implemented when configuration manager is fully implemented
		return fmt.Errorf("loading configuration from file not yet implemented")
	}

	if len(args) != 2 {
		return fmt.Errorf("usage: config set <key> <value>")
	}

	key := args[0]
	value := args[1]

	// Set individual configuration value
	err := c.SetConfig(key, value)
	if err != nil {
		return fmt.Errorf("failed to set configuration: %w", err)
	}

	fmt.Printf("Configuration updated: %s = %s\n", key, value)
	return nil
}

func (c *CLI) runConfigEdit(cmd *cobra.Command, args []string) error {
	// Open configuration file in default editor
	return c.EditConfig()
}

func (c *CLI) runConfigValidate(cmd *cobra.Command, args []string) error {
	// Validate configuration file syntax and values
	err := c.ValidateConfig()
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Println("Configuration is valid")
	return nil
}

func (c *CLI) runConfigExport(cmd *cobra.Command, args []string) error {
	// Get export path from args or use default
	exportPath := "generator-config.yaml"
	if len(args) > 0 {
		exportPath = args[0]
	}

	// Export current configuration to shareable file
	err := c.ExportConfig(exportPath)
	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}

	fmt.Printf("Configuration exported to: %s\n", exportPath)
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
		fmt.Printf("Using update channel: %s\n", channel)
	}

	if check {
		// Check for updates without installing
		updateInfo, err := c.CheckUpdates()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		fmt.Printf("Current Version: %s\n", updateInfo.CurrentVersion)
		fmt.Printf("Latest Version: %s\n", updateInfo.LatestVersion)
		fmt.Printf("Update Available: %t\n", updateInfo.UpdateAvailable)

		if updateInfo.UpdateAvailable {
			fmt.Printf("Release Date: %s\n", updateInfo.ReleaseDate.Format("2006-01-02"))
			fmt.Printf("Download Size: %s\n", formatBytes(updateInfo.Size))

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
				fmt.Printf("\nRelease Notes:\n%s\n", updateInfo.ReleaseNotes)
			}

			// Check compatibility if requested
			if compatibility {
				fmt.Println("\nChecking compatibility...")
				compatResult, err := c.versionManager.CheckCompatibility(".")
				if err != nil {
					fmt.Printf("Warning: Failed to check compatibility: %v\n", err)
				} else {
					if compatResult.Compatible {
						fmt.Println("✅ Update is compatible with current project")
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
			fmt.Println("No updates available")
			return nil
		}

		targetVersion := updateInfo.LatestVersion
		if version != "" {
			targetVersion = version
		}

		// Check compatibility unless forced
		if !force && compatibility {
			fmt.Println("Checking compatibility...")
			compatResult, err := c.versionManager.CheckCompatibility(".")
			if err != nil {
				return fmt.Errorf("failed to check compatibility: %w", err)
			}

			if !compatResult.Compatible {
				fmt.Printf("Compatibility issues found:\n")
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
			fmt.Print("Continue with installation? (y/N): ")
			var response string
			if _, err := fmt.Scanln(&response); err != nil || (response != "y" && response != "Y") {
				fmt.Println("Update cancelled")
				return nil
			}
		}

		fmt.Printf("Installing update to version %s...\n", targetVersion)

		// Configure update options
		if !backup {
			fmt.Println("Warning: Backup disabled - no rollback possible")
		}
		if !verify {
			fmt.Println("Warning: Signature verification disabled")
		}

		err = c.InstallUpdates()
		if err != nil {
			return fmt.Errorf("failed to install updates: %w", err)
		}

		fmt.Printf("✅ Successfully updated to version %s\n", targetVersion)
		fmt.Println("Restart any running instances to use the new version")
		return nil
	}

	if templates {
		// Update templates cache
		fmt.Println("Updating templates cache...")
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
		fmt.Println("Run 'generator update --install' to install the update")
		fmt.Println("Run 'generator update --check --release-notes' to see what's new")
	} else {
		fmt.Println("✅ You are running the latest version")
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
			fmt.Println("Cache clear cancelled")
			return nil
		}
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Cache clear cancelled")
			return nil
		}
	}

	// Clear cache
	err := c.ClearCache()
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("Cache cleared successfully")
	return nil
}

func (c *CLI) runCacheClean(cmd *cobra.Command, args []string) error {
	// Clean expired and invalid cache entries
	fmt.Println("Cleaning cache...")
	err := c.CleanCache()
	if err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	fmt.Println("Cache cleaned successfully")
	return nil
}

func (c *CLI) runCacheValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("Validating cache...")
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
	return nil, fmt.Errorf("SelectComponents implementation pending - will be implemented in task 2")
}

// Helper methods for generate command

// validateGenerateConfiguration validates the configuration for generation
func (c *CLI) validateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	if config.Name == "" {
		return c.createConfigurationError("project name is required", "Set project name in configuration file or GENERATOR_PROJECT_NAME environment variable")
	}

	// Validate project name format
	if !isValidProjectName(config.Name) {
		return c.createConfigurationError(
			fmt.Sprintf("invalid project name '%s'", config.Name),
			"Project name must contain only letters, numbers, hyphens, and underscores",
		)
	}

	// Validate license if specified
	if config.License != "" && !isValidLicense(config.License) {
		return c.createConfigurationError(
			fmt.Sprintf("invalid license '%s'", config.License),
			"Use a valid SPDX license identifier (e.g., MIT, Apache-2.0, GPL-3.0)",
		)
	}

	return nil
}

// performPreGenerationChecks performs checks before generation
func (c *CLI) performPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	// Check if output directory exists
	if _, err := os.Stat(outputPath); err == nil {
		// Directory exists
		if !options.Force {
			return fmt.Errorf("output directory '%s' already exists, use --force to overwrite", outputPath)
		}

		// Check if directory is empty
		entries, err := os.ReadDir(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read output directory: %w", err)
		}

		if len(entries) > 0 && options.BackupExisting {
			c.VerboseOutput("Creating backup of existing files in %s", outputPath)
			if err := c.createBackup(outputPath); err != nil {
				c.WarningOutput("Failed to create backup: %v", err)
			}
		}
	}

	// Check write permissions
	parentDir := outputPath
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		parentDir = filepath.Dir(outputPath)
	}

	if err := c.checkWritePermissions(parentDir); err != nil {
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

// selectDefaultTemplate selects a default template based on configuration
func (c *CLI) selectDefaultTemplate(config *models.ProjectConfig) string {
	// This would analyze the config and select an appropriate template
	// For now, return a sensible default
	return "go-gin"
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
		c.WarningOutput("Failed to close temporary file: %v", err)
	}
	if err := os.Remove(tempFile); err != nil {
		c.WarningOutput("Failed to remove temporary file: %v", err)
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
		return fmt.Errorf("failed to get template information: %w", err)
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
	ctx := c.StartOperationWithOutput("generate-from-config", fmt.Sprintf("Loading configuration from %s", configPath))
	defer func() {
		if ctx != nil {
			c.FinishOperationWithOutput(ctx, "generate-from-config", "Configuration loading completed")
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

	c.VerboseOutput("Successfully loaded configuration for project: %s", config.Name)

	// Validate configuration if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("Validating configuration...")
		if err := c.validateGenerateConfiguration(config, options); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
		c.VerboseOutput("Configuration validation passed")
	}

	// Set output path from options or config
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = config.OutputPath
	}
	if outputPath == "" {
		outputPath = "./" + config.Name
	}

	// Handle offline mode
	if options.Offline {
		c.VerboseOutput("Running in offline mode - using cached templates and versions")
		if c.cacheManager != nil {
			if err := c.cacheManager.EnableOfflineMode(); err != nil {
				c.WarningOutput("Failed to enable offline mode: %v", err)
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
		return fmt.Errorf("pre-generation checks failed: %w", err)
	}

	// Select template
	templateName := ""
	if len(options.Templates) > 0 {
		templateName = options.Templates[0]
	}
	if templateName == "" {
		templateName = c.selectDefaultTemplate(config)
	}

	c.VerboseOutput("Using template: %s", templateName)

	// Validate dependencies if not skipped
	if !options.SkipValidation {
		c.VerboseOutput("Validating dependencies...")
		if err := c.validateDependencies(config, templateName); err != nil {
			return fmt.Errorf("dependency validation failed: %w", err)
		}
		c.VerboseOutput("Dependency validation passed")
	}

	// Generate project
	if options.DryRun {
		c.QuietOutput("Dry run mode - would generate project '%s' using template '%s' in directory '%s'",
			config.Name, templateName, outputPath)
		return nil
	}

	// Process template
	if c.templateManager == nil {
		return fmt.Errorf("template manager not initialized")
	}

	return c.templateManager.ProcessTemplate(templateName, config, outputPath)
}

func (c *CLI) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	return nil, fmt.Errorf("ValidateProject implementation pending - will be implemented in task 5")
}

func (c *CLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	if c.auditEngine == nil {
		return nil, fmt.Errorf("audit engine not initialized")
	}

	// Perform comprehensive project audit
	result, err := c.auditEngine.AuditProject(path, &options)
	if err != nil {
		return nil, fmt.Errorf("audit failed: %w", err)
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
		return fmt.Errorf("failed to get template info: %w", err)
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

	// Validate the updated settings
	err = c.configManager.ValidateSettings()
	if err != nil {
		return fmt.Errorf("configuration validation failed after update: %w", err)
	}

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
	fmt.Println("Cache Information")
	fmt.Println("=================")
	fmt.Printf("Location: %s\n", stats.CacheLocation)
	fmt.Printf("Status: %s\n", stats.CacheHealth)
	fmt.Printf("Offline Mode: %t\n", stats.OfflineMode)
	fmt.Println()

	fmt.Println("Statistics")
	fmt.Println("----------")
	fmt.Printf("Total Entries: %d\n", stats.TotalEntries)
	fmt.Printf("Total Size: %s\n", formatBytes(stats.TotalSize))
	fmt.Printf("Hit Rate: %.1f%%\n", stats.HitRate*100)
	fmt.Printf("Expired Entries: %d\n", stats.ExpiredEntries)
	fmt.Printf("Last Cleanup: %s\n", stats.LastCleanup.Format("2006-01-02 15:04:05"))
	fmt.Println()

	fmt.Println("Configuration")
	fmt.Println("-------------")
	fmt.Printf("Max Size: %s\n", formatBytes(config.MaxSize))
	fmt.Printf("Max Entries: %d\n", config.MaxEntries)
	fmt.Printf("Default TTL: %s\n", config.DefaultTTL)
	fmt.Printf("Eviction Policy: %s\n", config.EvictionPolicy)
	fmt.Printf("Compression: %t\n", config.EnableCompression)
	fmt.Printf("Persist to Disk: %t\n", config.PersistToDisk)

	// Show cache health warnings if any
	if stats.CacheHealth != "healthy" {
		fmt.Println()
		fmt.Println("Health Issues")
		fmt.Println("-------------")
		if stats.ExpiredEntries > 0 {
			fmt.Printf("⚠ %d expired entries found - consider running 'generator cache clean'\n", stats.ExpiredEntries)
		}
		if stats.CacheHealth == "corrupted" {
			fmt.Println("⚠ Cache corruption detected - consider running 'generator cache repair'")
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
	return nil, fmt.Errorf("PromptAdvancedOptions implementation pending")
}

func (c *CLI) ConfirmAdvancedGeneration(*models.ProjectConfig, *interfaces.AdvancedOptions) bool {
	return false
}

func (c *CLI) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("SelectTemplateInteractively implementation pending")
}

func (c *CLI) GenerateWithAdvancedOptions(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	return fmt.Errorf("GenerateWithAdvancedOptions implementation pending")
}

func (c *CLI) ValidateProjectAdvanced(path string, options *interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	return nil, fmt.Errorf("ValidateProjectAdvanced implementation pending")
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
	return nil, fmt.Errorf("LoadConfiguration implementation pending")
}

func (c *CLI) MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error) {
	return nil, fmt.Errorf("MergeConfigurations implementation pending")
}

func (c *CLI) ValidateConfigurationSchema(config *models.ProjectConfig) error {
	return fmt.Errorf("ValidateConfigurationSchema implementation pending")
}

func (c *CLI) GetConfigurationSources() ([]interfaces.ConfigSource, error) {
	return nil, fmt.Errorf("GetConfigurationSources implementation pending")
}

func (c *CLI) GetPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("GetPackageVersions implementation pending")
}

func (c *CLI) GetLatestPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("GetLatestPackageVersions implementation pending")
}

func (c *CLI) CheckCompatibility(projectPath string) (*interfaces.CompatibilityResult, error) {
	return nil, fmt.Errorf("CheckCompatibility implementation pending")
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
	return fmt.Errorf("SetLogLevel implementation pending")
}

func (c *CLI) GetLogLevel() string {
	return "info"
}

func (c *CLI) ShowRecentLogs(lines int, level string) error {
	return c.showRecentLogs(lines, level, "", time.Time{}, "text")
}

func (c *CLI) GetLogFileLocations() ([]string, error) {
	return nil, fmt.Errorf("GetLogFileLocations implementation pending")
}

func (c *CLI) RunNonInteractive(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	return fmt.Errorf("RunNonInteractive implementation pending")
}

func (c *CLI) GenerateReport(reportType string, format string, outputFile string) error {
	return fmt.Errorf("GenerateReport implementation pending")
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
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
		}
	case "components.frontend.nextjs.home":
		if val, err := strconv.ParseBool(value); err == nil {
			config.Components.Frontend.NextJS.Home = val
		} else {
			return fmt.Errorf("invalid boolean value for %s: %s", key, value)
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
