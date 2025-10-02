// Package cli provides command registration and management functionality for the CLI.
//
// This module handles the registration and setup of all CLI commands,
// organizing command definitions and their configuration in a centralized way.
package cli

import (
	"github.com/spf13/cobra"
)

// CommandRegistry manages the registration and setup of all CLI commands.
//
// The CommandRegistry provides centralized command management including:
//   - Command registration and setup
//   - Command organization and orchestration
//   - Consistent command configuration
type CommandRegistry struct {
	cli *CLI
}

// NewCommandRegistry creates a new CommandRegistry instance.
func NewCommandRegistry(cli *CLI) *CommandRegistry {
	return &CommandRegistry{
		cli: cli,
	}
}

// RegisterAllCommands registers all CLI commands with the root command.
//
// This method orchestrates the registration of all available commands
// in the proper order and ensures consistent setup.
func (cr *CommandRegistry) RegisterAllCommands() {
	// Add all commands to the root command
	cr.setupGenerateCommand()
	cr.setupValidateCommand()
	cr.setupAuditCommand()
	cr.setupVersionCommand()
	cr.setupConfigCommand()
	cr.setupListTemplatesCommand()
	cr.setupTemplateCommand()
	cr.setupUpdateCommand()
	cr.setupCacheCommand()
	cr.setupLogsCommand()
}

// setupGenerateCommand sets up the generate command with all documented flags
func (cr *CommandRegistry) setupGenerateCommand() {
	generateCmd := &cobra.Command{
		Use:   "generate [flags]",
		Short: "Generate a new project from templates with modern best practices",
		Long: `Generate production-ready projects with modern best practices.

Supports Go, Next.js, React, Android, iOS, Docker, Kubernetes, and Terraform.
Use interactive mode or provide a configuration file.

FLAG USAGE NOTES:
• Output modes: --verbose, --quiet, and --debug are mutually exclusive
• Generation modes: --interactive and --non-interactive cannot be used together
• Mode flags: Don't combine mode flags (--interactive) with --mode parameter
• Enhanced validation provides specific suggestions for flag conflicts`,
		RunE: cr.cli.commandHandlers.runGenerate,
		Example: `  # Interactive project creation
  generator generate
  
  # Generate from configuration file
  generator generate --config project.yaml
  
  # Use specific template with verbose output
  generator generate --template go-gin --verbose
  
  # Non-interactive mode for automation (quiet output)
  generator generate --config project.yaml --non-interactive --quiet
  
  # Force interactive mode with debug information
  generator generate --force-interactive --debug
  
  # AVOID: These combinations will cause conflicts
  # generator generate --verbose --quiet          # ❌ Output mode conflict
  # generator generate --interactive --mode=auto  # ❌ Mode specification conflict
  # generator generate --non-interactive --force-interactive  # ❌ Mode conflict`,
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
	generateCmd.Flags().Bool("interactive", false, "Force interactive mode for project configuration")
	generateCmd.Flags().String("preset", "", "Use predefined configuration preset")

	// Mode-specific flags
	generateCmd.Flags().Bool("force-interactive", false, "Force interactive mode even in CI/automated environments")
	generateCmd.Flags().Bool("force-non-interactive", false, "Force non-interactive mode even in terminal environments")
	generateCmd.Flags().String("mode", "", "Explicitly set generation mode (interactive, non-interactive, config-file)")

	cr.cli.rootCmd.AddCommand(generateCmd)
}

// setupValidateCommand sets up the validate command with all documented flags
func (cr *CommandRegistry) setupValidateCommand() {
	validateCmd := &cobra.Command{
		Use:   "validate [path] [flags]",
		Short: "Validate project structure, configuration, and dependencies",
		Long: `Validate project structure, configuration files, and dependencies.

Checks code quality, security, and best practices. Can automatically fix common issues.`,
		RunE: cr.cli.commandHandlers.runValidate,
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

	cr.cli.rootCmd.AddCommand(validateCmd)
}

// setupAuditCommand sets up the audit command with all documented flags
func (cr *CommandRegistry) setupAuditCommand() {
	auditCmd := &cobra.Command{
		Use:   "audit [path] [flags]",
		Short: "Comprehensive security, quality, and compliance auditing",
		Long: `Audit project security, code quality, license compliance, and performance.

Provides detailed reports with scores and recommendations for improvement.`,
		RunE: cr.cli.commandHandlers.runAudit,
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

	cr.cli.rootCmd.AddCommand(auditCmd)
}

// setupVersionCommand sets up the version command with all documented flags
func (cr *CommandRegistry) setupVersionCommand() {
	versionCmd := &cobra.Command{
		Use:   "version [flags]",
		Short: "Display comprehensive version information",
		Long: `Display comprehensive version information for the generator and supported technologies.

This command provides detailed version information including:
- Generator version and build information
- Latest versions of supported packages and technologies
- Update availability and compatibility information
- Build metadata and system information

ENHANCED JSON OUTPUT:
The --json flag now produces properly structured JSON with comprehensive
version information, build metadata, and system details. The JSON output
is validated and formatted for easy parsing by automation tools.

The command supports multiple output formats and can check for updates
both for the generator itself and for supported technology packages.`,
		RunE: cr.cli.commandHandlers.runVersion,
		Example: `  # Basic version information
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
  generator version --packages --verbose --debug`,
	}

	// Add additional flags for the main CLI
	versionCmd.Flags().Bool("packages", false, "Show latest package versions for all supported technologies")
	versionCmd.Flags().Bool("check-updates", false, "Check for generator updates")
	versionCmd.Flags().Bool("build-info", false, "Show detailed build information")
	versionCmd.Flags().Bool("short", false, "Show only version number")
	versionCmd.Flags().String("format", "text", "Output format (text, json, yaml)")
	versionCmd.Flags().Bool("json", false, "Output version information in JSON format")
	versionCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	versionCmd.Flags().String("check-package", "", "Check version for specific package")

	cr.cli.rootCmd.AddCommand(versionCmd)
}

// setupConfigCommand sets up the config command with subcommands
func (cr *CommandRegistry) setupConfigCommand() {
	// The config command is already implemented in config_commands.go
	// We just need to call the existing setupConfigCommand method
	cr.cli.setupConfigCommand()
}

// setupListTemplatesCommand sets up the list-templates command
func (cr *CommandRegistry) setupListTemplatesCommand() {
	listTemplatesCmd := &cobra.Command{
		Use:   "list-templates [flags]",
		Short: "List and discover available project templates",
		Long: `List available project templates with filtering and search.

Browse templates for frontend, backend, mobile, and infrastructure projects.`,
		RunE: cr.cli.commandHandlers.runListTemplates,
		Example: `  # List all templates
  generator list-templates
  
  # Filter by category
  generator list-templates --category backend
  
  # Search for templates
  generator list-templates --search api
  
  # Show detailed information
  generator list-templates --detailed`,
	}

	// Set up flags directly instead of using templateCmd
	listTemplatesCmd.Flags().String("category", "", "Filter by category (frontend, backend, mobile, infrastructure)")
	listTemplatesCmd.Flags().String("technology", "", "Filter by technology (go, nodejs, react, etc.)")
	listTemplatesCmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	listTemplatesCmd.Flags().String("search", "", "Search templates by name or description")
	listTemplatesCmd.Flags().Bool("detailed", false, "Show detailed template information")

	cr.cli.rootCmd.AddCommand(listTemplatesCmd)
}

// setupTemplateCommand sets up the template command with subcommands
func (cr *CommandRegistry) setupTemplateCommand() {
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
dependencies, configuration options, and usage examples.`,
		RunE: cr.cli.commandHandlers.runTemplateInfo,
		Args: cobra.ExactArgs(1),
		Example: `  # Show template information
  generator template info go-gin
  
  # Show detailed information
  generator template info nextjs-app --detailed
  
  # Show template variables
  generator template info go-gin --variables`,
	}
	// Set up flags directly instead of using templateCmd
	templateInfoCmd.Flags().Bool("detailed", false, "Show detailed template information")
	templateInfoCmd.Flags().Bool("variables", false, "Show template variables")
	templateInfoCmd.Flags().Bool("dependencies", false, "Show template dependencies")
	templateInfoCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	templateCmd.AddCommand(templateInfoCmd)

	// template validate
	templateValidateCmd := &cobra.Command{
		Use:   "validate <template-path> [flags]",
		Short: "Validate custom template structure, syntax, and compliance",
		Long: `Validate custom template directories including structure, metadata, syntax, and best practices.

Provides detailed feedback and auto-fix capabilities for common issues.`,
		RunE: cr.cli.commandHandlers.runTemplateValidate,
		Args: cobra.ExactArgs(1),
		Example: `  # Validate template directory
  generator template validate ./my-template
  
  # Validate with detailed output
  generator template validate ./my-template --detailed
  
  # Validate and auto-fix issues
  generator template validate ./my-template --fix`,
	}
	// Set up flags directly instead of using templateCmd
	templateValidateCmd.Flags().Bool("detailed", false, "Show detailed validation results")
	templateValidateCmd.Flags().Bool("fix", false, "Attempt to fix validation issues")
	templateValidateCmd.Flags().String("output-format", "text", "Output format (text, json)")
	templateCmd.AddCommand(templateValidateCmd)

	cr.cli.rootCmd.AddCommand(templateCmd)
}

// setupUpdateCommand sets up the update command
func (cr *CommandRegistry) setupUpdateCommand() {
	updateCmd := &cobra.Command{
		Use:   "update [flags]",
		Short: "Update generator, templates, and package information",
		Long: `Update generator, templates, and package information.

Includes safety checks, rollback capabilities, and multiple update channels.`,
		RunE: cr.cli.commandHandlers.runUpdate,
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

	cr.cli.rootCmd.AddCommand(updateCmd)
}

// setupCacheCommand sets up the cache command with all subcommands
func (cr *CommandRegistry) setupCacheCommand() {
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
		RunE:  cr.cli.commandHandlers.runCacheShow,
	}
	cacheCmd.AddCommand(cacheShowCmd)

	// cache clear
	cacheClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cached data",
		Long:  "Remove all cached templates, versions, and other data",
		RunE:  cr.cli.commandHandlers.runCacheClear,
	}
	cacheClearCmd.Flags().Bool("force", false, "Clear cache without confirmation")
	cacheCmd.AddCommand(cacheClearCmd)

	// cache clean
	cacheCleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove expired and invalid cache entries",
		Long:  "Clean up expired cache entries and repair corrupted cache data",
		RunE:  cr.cli.commandHandlers.runCacheClean,
	}
	cacheCmd.AddCommand(cacheCleanCmd)

	// cache validate
	cacheValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate cache integrity",
		Long:  "Check cache integrity and report any issues",
		RunE:  cr.cli.commandHandlers.runCacheValidate,
	}
	cacheCmd.AddCommand(cacheValidateCmd)

	// cache repair
	cacheRepairCmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair corrupted cache data",
		Long:  "Attempt to repair corrupted cache entries and fix cache issues",
		RunE:  cr.cli.commandHandlers.runCacheRepair,
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
		RunE:  cr.cli.commandHandlers.runCacheOfflineEnable,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineEnableCmd)

	// cache offline disable
	cacheOfflineDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable offline mode",
		Long:  "Disable offline mode to allow network access",
		RunE:  cr.cli.commandHandlers.runCacheOfflineDisable,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineDisableCmd)

	// cache offline status
	cacheOfflineStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show offline mode status",
		Long:  "Display current offline mode status and readiness",
		RunE:  cr.cli.commandHandlers.runCacheOfflineStatus,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineStatusCmd)

	cacheCmd.AddCommand(cacheOfflineCmd)
	cr.cli.rootCmd.AddCommand(cacheCmd)
}

// setupLogsCommand sets up the logs command
func (cr *CommandRegistry) setupLogsCommand() {
	logsCmd := &cobra.Command{
		Use:   "logs [flags]",
		Short: "View and manage application logs",
		Long: `View application logs with filtering and formatting options.

Provides access to detailed logs for debugging and monitoring purposes.`,
		RunE: cr.cli.commandHandlers.runLogs,
		Example: `  # View recent logs
  generator logs
  
  # View logs with specific level
  generator logs --level error
  
  # Follow logs in real-time
  generator logs --follow
  
  # View logs from specific time
  generator logs --since 1h`,
	}

	logsCmd.Flags().String("level", "", "Filter by log level (debug, info, warn, error)")
	logsCmd.Flags().String("since", "", "Show logs since timestamp (e.g., 1h, 30m, 2006-01-02T15:04:05Z)")
	logsCmd.Flags().String("until", "", "Show logs until timestamp")
	logsCmd.Flags().Int("tail", 100, "Number of lines to show from the end")
	logsCmd.Flags().Bool("follow", false, "Follow log output")
	logsCmd.Flags().Bool("timestamps", true, "Show timestamps")
	logsCmd.Flags().String("format", "text", "Output format (text, json)")
	logsCmd.Flags().String("component", "", "Filter by component name")

	cr.cli.rootCmd.AddCommand(logsCmd)
}
