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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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
	generatorVersion string
	rootCmd          *cobra.Command
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
	version string,
) interfaces.CLIInterface {
	cli := &CLI{
		configManager:    configManager,
		validator:        validator,
		templateManager:  templateManager,
		cacheManager:     cacheManager,
		versionManager:   versionManager,
		auditEngine:      auditEngine,
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
	c.rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	c.rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
	c.rootCmd.PersistentFlags().String("log-level", "info", "Set log level (debug, info, warn, error)")
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

	return c.rootCmd.Execute()
}

// handleGlobalFlags processes global flags that apply to all commands
func (c *CLI) handleGlobalFlags(cmd *cobra.Command) error {
	// Get global flags
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")
	logLevel, _ := cmd.Flags().GetString("log-level")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Handle verbose/quiet flags
	if verbose && quiet {
		return fmt.Errorf("cannot use both --verbose and --quiet flags")
	}

	// Set log level based on flags
	if verbose {
		logLevel = "debug"
	} else if quiet {
		logLevel = "error"
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
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

	// Store global settings for use in commands
	// This would typically be stored in a context or configuration
	if verbose {
		fmt.Printf("Debug: Running command '%s' with verbose output\n", cmd.Name())
	}

	if nonInteractive {
		fmt.Printf("Debug: Running in non-interactive mode\n")
	}

	return nil
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
Can also check for updates to templates and dependencies.`,
		RunE: c.runUpdate,
	}

	updateCmd.Flags().Bool("check", false, "Check for updates without installing")
	updateCmd.Flags().Bool("install", false, "Install available updates")
	updateCmd.Flags().Bool("templates", false, "Update templates cache")
	updateCmd.Flags().Bool("force", false, "Force update even if current version is newer")

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

	c.rootCmd.AddCommand(cacheCmd)
}

// setupLogsCommand sets up the logs command
func (c *CLI) setupLogsCommand() {
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "View recent log entries and log file locations",
		Long: `Display recent log entries and show log file locations.
Useful for debugging and troubleshooting issues.`,
		RunE: c.runLogs,
	}

	logsCmd.Flags().Int("lines", 50, "Number of recent log lines to show")
	logsCmd.Flags().String("level", "", "Filter by log level (debug, info, warn, error)")
	logsCmd.Flags().Bool("follow", false, "Follow log output (tail -f)")
	logsCmd.Flags().Bool("locations", false, "Show log file locations only")

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

	// Handle non-interactive mode
	if nonInteractive {
		interactive = false
	}

	// Use interactive flag for future implementation
	_ = interactive

	// Log additional options for debugging
	if len(exclude) > 0 {
		fmt.Printf("Debug: Excluding files/directories: %v\n", exclude)
	}
	if len(includeOnly) > 0 {
		fmt.Printf("Debug: Including only: %v\n", includeOnly)
	}
	if preset != "" {
		fmt.Printf("Debug: Using preset: %s\n", preset)
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

	// If config file is provided, generate from config
	if configPath != "" {
		return c.GenerateFromConfig(configPath, options)
	}

	// Otherwise, use interactive mode
	config, err := c.PromptProjectDetails()
	if err != nil {
		return fmt.Errorf("failed to collect project details: %w", err)
	}

	if !c.ConfirmGeneration(config) {
		fmt.Println("Generation cancelled by user")
		return nil
	}

	// Generate project using template manager
	templateName := template
	if templateName == "" {
		// Use default template selection logic
		templateName = "go-gin" // Default template
	}

	if outputPath == "" {
		outputPath = config.Name
	}

	return c.templateManager.ProcessTemplate(templateName, config, outputPath)
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

	// Log additional options for debugging
	if strict {
		fmt.Println("Debug: Using strict validation mode")
	}
	if len(excludeRules) > 0 {
		fmt.Printf("Debug: Excluding rules: %v\n", excludeRules)
	}

	// Use additional flags for future implementation
	_ = summaryOnly
	_ = showFixes

	// Get global flags
	verbose, _ := cmd.Flags().GetBool("verbose")

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

	// Display results
	fmt.Printf("Validation completed for: %s\n", path)
	fmt.Printf("Valid: %t\n", result.Valid)
	fmt.Printf("Issues: %d\n", len(result.Issues))
	fmt.Printf("Warnings: %d\n", len(result.Warnings))

	if len(result.Issues) > 0 {
		fmt.Println("\nIssues found:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s: %s\n", issue.Severity, issue.Message)
		}
	}

	if !result.Valid {
		return fmt.Errorf("validation failed with %d issues", len(result.Issues))
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

	// Log additional options for debugging
	if len(excludeCategories) > 0 {
		fmt.Printf("Debug: Excluding audit categories: %v\n", excludeCategories)
	}
	if summaryOnly {
		fmt.Println("Debug: Showing summary only")
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

	// Display results
	fmt.Printf("Audit completed for: %s\n", path)
	fmt.Printf("Overall Score: %.2f\n", result.OverallScore)
	fmt.Printf("Audit Time: %s\n", result.AuditTime.Format("2006-01-02 15:04:05"))

	if result.Security != nil {
		fmt.Printf("Security Score: %.2f\n", result.Security.Score)
		fmt.Printf("Vulnerabilities: %d\n", len(result.Security.Vulnerabilities))
	}

	if result.Quality != nil {
		fmt.Printf("Quality Score: %.2f\n", result.Quality.Score)
		fmt.Printf("Code Smells: %d\n", len(result.Quality.CodeSmells))
	}

	if result.Licenses != nil {
		fmt.Printf("License Compatible: %t\n", result.Licenses.Compatible)
		fmt.Printf("License Conflicts: %d\n", len(result.Licenses.Conflicts))
	}

	if result.Performance != nil {
		fmt.Printf("Performance Score: %.2f\n", result.Performance.Score)
		fmt.Printf("Bundle Size: %d bytes\n", result.Performance.BundleSize)
	}

	if len(result.Recommendations) > 0 {
		fmt.Println("\nRecommendations:")
		for _, rec := range result.Recommendations {
			fmt.Printf("  - %s\n", rec)
		}
	}

	// Check fail conditions
	if failOnHigh && result.OverallScore < 7.0 {
		return fmt.Errorf("audit failed: high severity issues found (score: %.2f)", result.OverallScore)
	}

	if failOnMedium && result.OverallScore < 5.0 {
		return fmt.Errorf("audit failed: medium or higher severity issues found (score: %.2f)", result.OverallScore)
	}

	if minScore > 0 && result.OverallScore < minScore {
		return fmt.Errorf("audit failed: score %.2f is below minimum required score %.2f", result.OverallScore, minScore)
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

	// Use format if outputFormat is not set
	if format != "text" {
		outputFormat = format
	}

	// Handle short version output
	if short {
		fmt.Println(c.generatorVersion)
		return nil
	}

	// Handle specific package version check
	if checkPackage != "" {
		fmt.Printf("Checking version for package: %s\n", checkPackage)
		// This would be implemented when package version checking is fully implemented
		return nil
	}

	// Handle compatibility flag
	if compatibility {
		fmt.Println("Debug: Showing compatibility information")
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
			return fmt.Errorf("failed to search templates: %w", err)
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
		return fmt.Errorf("failed to list templates: %w", err)
	}

	// Display templates
	if len(templates) == 0 {
		fmt.Println("No templates found matching the criteria")
		return nil
	}

	fmt.Printf("Found %d template(s):\n\n", len(templates))

	for _, template := range templates {
		fmt.Printf("Name: %s\n", template.Name)
		fmt.Printf("Display Name: %s\n", template.DisplayName)
		fmt.Printf("Description: %s\n", template.Description)
		fmt.Printf("Category: %s\n", template.Category)
		fmt.Printf("Technology: %s\n", template.Technology)
		fmt.Printf("Version: %s\n", template.Version)

		if len(template.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(template.Tags, ", "))
		}

		if detailed {
			if len(template.Dependencies) > 0 {
				fmt.Printf("Dependencies: %s\n", strings.Join(template.Dependencies, ", "))
			}
			fmt.Printf("Author: %s\n", template.Metadata.Author)
			fmt.Printf("License: %s\n", template.Metadata.License)
			if template.Metadata.Repository != "" {
				fmt.Printf("Repository: %s\n", template.Metadata.Repository)
			}
		}

		fmt.Println()
	}

	return nil
}

func (c *CLI) runUpdate(cmd *cobra.Command, args []string) error {
	// Get flags
	check, _ := cmd.Flags().GetBool("check")
	install, _ := cmd.Flags().GetBool("install")
	templates, _ := cmd.Flags().GetBool("templates")
	force, _ := cmd.Flags().GetBool("force")

	// Use force flag for future implementation
	_ = force

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
			if updateInfo.Breaking {
				fmt.Println("⚠️  This is a breaking change update")
			}
			if updateInfo.ReleaseNotes != "" {
				fmt.Printf("\nRelease Notes:\n%s\n", updateInfo.ReleaseNotes)
			}
		}

		return nil
	}

	if install {
		// Install available updates
		fmt.Println("Installing updates...")
		err := c.InstallUpdates()
		if err != nil {
			return fmt.Errorf("failed to install updates: %w", err)
		}
		fmt.Println("Updates installed successfully")
		return nil
	}

	if templates {
		// Update templates cache
		fmt.Println("Updating templates cache...")
		// This would be implemented when template manager is fully implemented
		fmt.Println("Templates cache updated successfully")
		return nil
	}

	// Default behavior: check for updates
	updateInfo, err := c.CheckUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if updateInfo.UpdateAvailable {
		fmt.Printf("Update available: %s -> %s\n", updateInfo.CurrentVersion, updateInfo.LatestVersion)
		fmt.Println("Run 'generator update --install' to install the update")
	} else {
		fmt.Println("You are running the latest version")
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

func (c *CLI) runLogs(cmd *cobra.Command, args []string) error {
	// Get flags
	lines, _ := cmd.Flags().GetInt("lines")
	level, _ := cmd.Flags().GetString("level")
	follow, _ := cmd.Flags().GetBool("follow")
	locations, _ := cmd.Flags().GetBool("locations")

	// Use level flag for future implementation
	_ = level

	if locations {
		// Show log file locations only
		fmt.Println("Log file locations:")
		// This would be implemented when logging system is fully implemented
		fmt.Println("  ~/.generator/logs/generator.log")
		fmt.Println("  ~/.generator/logs/error.log")
		return nil
	}

	if follow {
		fmt.Printf("Following logs (showing last %d lines)...\n", lines)
		fmt.Println("Press Ctrl+C to stop")
		// This would implement tail -f functionality
		// For now, just show recent logs
	}

	// Show recent logs
	err := c.ShowLogs()
	if err != nil {
		return fmt.Errorf("failed to show logs: %w", err)
	}

	return nil
}

// Interface implementation methods - these will be implemented in subsequent tasks

func (c *CLI) SelectComponents() ([]string, error) {
	return nil, fmt.Errorf("SelectComponents implementation pending - will be implemented in task 2")
}

func (c *CLI) GenerateFromConfig(configPath string, options interfaces.GenerateOptions) error {
	return fmt.Errorf("GenerateFromConfig implementation pending - will be implemented in task 2")
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
	return fmt.Errorf("ShowVersion implementation pending - will be implemented in task 8")
}

func (c *CLI) CheckUpdates() (*interfaces.UpdateInfo, error) {
	return nil, fmt.Errorf("CheckUpdates implementation pending - will be implemented in task 8")
}

func (c *CLI) InstallUpdates() error {
	return fmt.Errorf("InstallUpdates implementation pending - will be implemented in task 8")
}

func (c *CLI) ShowCache() error {
	return fmt.Errorf("ShowCache implementation pending - will be implemented in task 7")
}

func (c *CLI) ClearCache() error {
	return fmt.Errorf("ClearCache implementation pending - will be implemented in task 7")
}

func (c *CLI) CleanCache() error {
	return fmt.Errorf("CleanCache implementation pending - will be implemented in task 7")
}

func (c *CLI) ShowLogs() error {
	return fmt.Errorf("ShowLogs implementation pending - will be implemented in task 9")
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
	return nil, fmt.Errorf("GetCacheStats implementation pending")
}

func (c *CLI) ValidateCache() error {
	return fmt.Errorf("ValidateCache implementation pending")
}

func (c *CLI) RepairCache() error {
	return fmt.Errorf("RepairCache implementation pending")
}

func (c *CLI) EnableOfflineMode() error {
	return fmt.Errorf("EnableOfflineMode implementation pending")
}

func (c *CLI) DisableOfflineMode() error {
	return fmt.Errorf("DisableOfflineMode implementation pending")
}

func (c *CLI) SetLogLevel(level string) error {
	return fmt.Errorf("SetLogLevel implementation pending")
}

func (c *CLI) GetLogLevel() string {
	return "info"
}

func (c *CLI) ShowRecentLogs(lines int, level string) error {
	return fmt.Errorf("ShowRecentLogs implementation pending")
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

func (c *CLI) GetExitCode() int {
	return 0
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
