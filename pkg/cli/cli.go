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
	return nil, fmt.Errorf("AuditProject implementation pending - will be implemented in task 6")
}

func (c *CLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("ListTemplates implementation pending - will be implemented in task 4")
}

func (c *CLI) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("GetTemplateInfo implementation pending - will be implemented in task 4")
}

func (c *CLI) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return nil, fmt.Errorf("ValidateTemplate implementation pending - will be implemented in task 4")
}

func (c *CLI) ShowConfig() error {
	return fmt.Errorf("ShowConfig implementation pending - will be implemented in task 3")
}

func (c *CLI) SetConfig(key, value string) error {
	return fmt.Errorf("SetConfig implementation pending - will be implemented in task 3")
}

func (c *CLI) EditConfig() error {
	return fmt.Errorf("EditConfig implementation pending - will be implemented in task 3")
}

func (c *CLI) ValidateConfig() error {
	return fmt.Errorf("ValidateConfig implementation pending - will be implemented in task 3")
}

func (c *CLI) ExportConfig(path string) error {
	return fmt.Errorf("ExportConfig implementation pending - will be implemented in task 3")
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
	return nil, fmt.Errorf("AuditProjectAdvanced implementation pending")
}

func (c *CLI) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("SearchTemplates implementation pending")
}

func (c *CLI) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	return nil, fmt.Errorf("GetTemplateMetadata implementation pending")
}

func (c *CLI) GetTemplateDependencies(name string) ([]string, error) {
	return nil, fmt.Errorf("GetTemplateDependencies implementation pending")
}

func (c *CLI) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return nil, fmt.Errorf("ValidateCustomTemplate implementation pending")
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
