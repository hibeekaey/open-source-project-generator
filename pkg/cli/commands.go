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
	"github.com/spf13/cobra"
)

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

	// Additional flags for comprehensive functionality
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
dependencies, configuration options, and usage examples.`,
		RunE: c.runTemplateInfo,
		Args: cobra.ExactArgs(1),
		Example: `  # Show template information
  generator template info go-gin
  
  # Show detailed information
  generator template info nextjs-app --detailed
  
  # Show template variables
  generator template info go-gin --variables`,
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
		Long: `Validate custom template directories including structure, metadata, syntax, and best practices.

Provides detailed feedback and auto-fix capabilities for common issues.`,
		RunE: c.runTemplateValidate,
		Args: cobra.ExactArgs(1),
		Example: `  # Validate template directory
  generator template validate ./my-template
  
  # Validate with detailed output
  generator template validate ./my-template --detailed
  
  # Validate and auto-fix issues
  generator template validate ./my-template --fix`,
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
		Long:  "Disable offline mode to allow network requests",
		RunE:  c.runCacheOfflineDisable,
	}
	cacheOfflineCmd.AddCommand(cacheOfflineDisableCmd)

	// cache offline status
	cacheOfflineStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show offline mode status",
		Long:  "Display current offline mode status and configuration",
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
		Example: `  # Show recent logs
  generator logs
  
  # Show error logs only
  generator logs --level error
  
  # Follow logs in real-time
  generator logs --follow
  
  # Show logs from last hour
  generator logs --since "1h"`,
	}

	logsCmd.Flags().String("level", "", "Filter by log level (debug, info, warn, error, fatal)")
	logsCmd.Flags().String("component", "", "Filter by component name")
	logsCmd.Flags().Int("lines", 50, "Number of recent log lines to show")
	logsCmd.Flags().String("since", "", "Show logs since specific time (RFC3339 format)")
	logsCmd.Flags().Bool("follow", false, "Follow log output in real-time (tail -f)")
	logsCmd.Flags().String("format", "text", "Output format (text, json, raw)")
	logsCmd.Flags().Bool("timestamps", true, "Show timestamps in output")
	logsCmd.Flags().Bool("no-color", false, "Disable colored output")
	logsCmd.Flags().Bool("locations", false, "Show log file locations only")

	c.rootCmd.AddCommand(logsCmd)
}
