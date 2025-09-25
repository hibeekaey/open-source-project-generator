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
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Command execution methods

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
		return fmt.Errorf("🚫 %s %s",
			c.error("Conflicting mode flags detected."),
			c.info("Please use only one mode flag at a time"))
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
			return fmt.Errorf("🚫 %s %s",
				c.error("Configuration validation failed."),
				c.info("Please check your settings and try again"))
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
		_ = c.isNonInteractiveMode()
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
		return fmt.Errorf("🚫 %s %s",
			c.error("Project validation encountered an issue."),
			c.info("Try running with --verbose to see more details"))
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
			c.QuietOutput("%s %s",
				c.error("❌ Found some issues that need attention."),
				c.info("See details below"))
		}
		c.QuietOutput("📊 Issues: %s", c.error(fmt.Sprintf("%d", len(result.Issues))))
		c.QuietOutput("⚠️  Warnings: %s", c.warning(fmt.Sprintf("%d", len(result.Warnings))))

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
			c.QuietOutput("\n%s", c.warning("⚠️  Things to consider:"))
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
	if report {
		return c.generateValidationReport(result, reportFormat, outputFile)
	}

	// Set exit code based on validation results
	if !result.Valid {
		c.SetExitCode(1)
		return fmt.Errorf("🚫 Found %s validation issues that need your attention", c.error(fmt.Sprintf("%d", len(result.Issues))))
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

	// Additional audit flags
	_, _ = cmd.Flags().GetBool("fail-on-high")
	_, _ = cmd.Flags().GetBool("fail-on-medium")
	_, _ = cmd.Flags().GetFloat64("min-score")
	_, _ = cmd.Flags().GetStringSlice("exclude-categories")
	_, _ = cmd.Flags().GetBool("summary-only")

	// Get global flags
	_, _ = cmd.Flags().GetBool("verbose")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		_ = c.isNonInteractiveMode()
	}

	// Create audit options
	options := &interfaces.AuditOptions{
		Security:     security,
		Quality:      quality,
		Licenses:     licenses,
		Performance:  performance,
		Detailed:     detailed,
		OutputFormat: outputFormat,
		OutputFile:   outputFile,
	}

	// Perform audit
	result, err := c.AuditProject(path, *options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Project audit encountered an issue."),
			c.info("Try running with --verbose to see more details"))
	}

	// Output results based on format and mode
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		// Machine-readable output for automation
		return c.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	if !c.quietMode {
		c.QuietOutput("🔍 Audit completed for: %s", path)
		c.QuietOutput("🎉 Overall Score: %.1f/100", result.OverallScore)
		if result.Security != nil {
			c.QuietOutput("🔒 Security Score: %.1f/100", result.Security.Score)
		}
		if result.Quality != nil {
			c.QuietOutput("✨ Quality Score: %.1f/100", result.Quality.Score)
		}
		if result.Licenses != nil {
			c.QuietOutput("License Compatible: %t", result.Licenses.Compatible)
		}
		if result.Performance != nil {
			c.QuietOutput("⚡ Performance Score: %.1f/100", result.Performance.Score)
		}

		if len(result.Recommendations) > 0 {
			c.QuietOutput("\nRecommendations:")
			for _, rec := range result.Recommendations {
				c.QuietOutput("  - %s", rec)
			}
		}
	}

	// Generate report if requested
	if outputFile != "" {
		return c.generateAuditReport(result, outputFormat, outputFile)
	}

	// Set exit code based on audit results
	if result.OverallScore < 70.0 { // Use a default minimum score
		c.SetExitCode(1)
		return fmt.Errorf("🚫 Audit score %.1f is below minimum required score %.1f", result.OverallScore, 70.0)
	}

	return nil
}

func (c *CLI) runListTemplates(cmd *cobra.Command, args []string) error {
	// Get flags
	category, _ := cmd.Flags().GetString("category")
	technology, _ := cmd.Flags().GetString("technology")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	_, _ = cmd.Flags().GetString("search")
	_, _ = cmd.Flags().GetBool("detailed")

	// Get global flags
	_, _ = cmd.Flags().GetBool("verbose")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		_ = c.isNonInteractiveMode()
	}

	// Create filter
	filter := interfaces.TemplateFilter{
		Category:   category,
		Technology: technology,
		Tags:       tags,
		// Search and Detailed fields don't exist in the interface
	}

	// List templates
	templates, err := c.ListTemplates(filter)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to list templates."),
			c.info("Try running with --verbose to see more details"))
	}

	// Output results
	if !c.quietMode {
		c.QuietOutput("📦 Available Templates (%d found)", len(templates))
		c.QuietOutput("")

		// Group templates by category
		categories := make(map[string][]interfaces.TemplateInfo)
		for _, template := range templates {
			categories[template.Category] = append(categories[template.Category], template)
		}

		// Display templates by category
		for category, templates := range categories {
			c.QuietOutput("🎨  %s Templates:", cases.Title(language.English).String(category))
			for _, template := range templates {
				c.QuietOutput("  • %s - %s", template.Name, template.Description)
			}
			c.QuietOutput("")
		}

		c.QuietOutput("💡 Use --detailed for more information")
		c.QuietOutput("🔍 Use --search <term> to find specific templates")
		c.QuietOutput("Listed %d templates", len(templates))
	}

	return nil
}

func (c *CLI) runUpdate(cmd *cobra.Command, args []string) error {
	// Get flags
	_, _ = cmd.Flags().GetBool("check")
	_, _ = cmd.Flags().GetBool("install")
	_, _ = cmd.Flags().GetBool("templates")
	_, _ = cmd.Flags().GetBool("force")
	_, _ = cmd.Flags().GetBool("compatibility")
	_, _ = cmd.Flags().GetBool("release-notes")
	_, _ = cmd.Flags().GetString("channel")
	_, _ = cmd.Flags().GetBool("backup")
	_, _ = cmd.Flags().GetBool("verify")
	_, _ = cmd.Flags().GetString("version")

	// Get global flags
	_, _ = cmd.Flags().GetBool("verbose")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		_ = c.isNonInteractiveMode()
	}

	// Check for updates
	result, err := c.CheckUpdates()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Update encountered an issue."),
			c.info("Try running with --verbose to see more details"))
	}

	// Output results
	if !c.quietMode {
		if result.UpdateAvailable {
			c.QuietOutput("🔄 Updates available!")
			c.QuietOutput("Current version: %s", result.CurrentVersion)
			c.QuietOutput("Latest version: %s", result.LatestVersion)
			if result.ReleaseNotes != "" {
				c.QuietOutput("Release notes: %s", result.ReleaseNotes)
			}
		} else {
			c.QuietOutput("✅ You're up to date!")
		}
	}

	return nil
}

// Cache command execution methods

func (c *CLI) runCacheShow(cmd *cobra.Command, args []string) error {
	// Show cache status
	status, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to get cache status."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("📊 Cache Status:")
		c.QuietOutput("Location: %s", status.CacheLocation)
		c.QuietOutput("Size: %d bytes", status.TotalSize)
		c.QuietOutput("Health: %s", status.CacheHealth)
		c.QuietOutput("Last Updated: %s", status.LastCleanup.Format(time.RFC3339))
	}

	return nil
}

func (c *CLI) runCacheClear(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	if !force && !c.isNonInteractiveMode() {
		c.QuietOutput("⚠️  This will clear all cached data. Continue? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			c.QuietOutput("❌ Cache clear cancelled")
			return nil
		}
	}

	// Clear cache
	err := c.cacheManager.Clear()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to clear cache."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Cache cleared successfully")
	}

	return nil
}

func (c *CLI) runCacheClean(cmd *cobra.Command, args []string) error {
	// Clean cache
	err := c.cacheManager.Clean()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to clean cache."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Cache cleaned successfully")
	}

	return nil
}

func (c *CLI) runCacheValidate(cmd *cobra.Command, args []string) error {
	// Validate cache
	err := c.cacheManager.ValidateCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Cache validation failed."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Cache is valid")
	}

	return nil
}

func (c *CLI) runCacheRepair(cmd *cobra.Command, args []string) error {
	// Repair cache
	err := c.cacheManager.RepairCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to repair cache."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Cache repaired successfully")
	}

	return nil
}

func (c *CLI) runCacheOfflineEnable(cmd *cobra.Command, args []string) error {
	// Enable offline mode
	err := c.cacheManager.EnableOfflineMode()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to enable offline mode."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Offline mode enabled")
	}

	return nil
}

func (c *CLI) runCacheOfflineDisable(cmd *cobra.Command, args []string) error {
	// Disable offline mode
	err := c.cacheManager.DisableOfflineMode()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to disable offline mode."),
			c.info("Try running with --verbose to see more details"))
	}

	if !c.quietMode {
		c.QuietOutput("✅ Offline mode disabled")
	}

	return nil
}

func (c *CLI) runCacheOfflineStatus(cmd *cobra.Command, args []string) error {
	// Get offline mode status
	enabled := c.cacheManager.IsOfflineMode()

	if !c.quietMode {
		if enabled {
			c.QuietOutput("📴 Offline mode is enabled")
		} else {
			c.QuietOutput("🌐 Offline mode is disabled")
		}
	}

	return nil
}

func (c *CLI) runLogs(cmd *cobra.Command, args []string) error {
	// Get flags
	level, _ := cmd.Flags().GetString("level")
	lines, _ := cmd.Flags().GetInt("lines")
	_, _ = cmd.Flags().GetString("component")
	_, _ = cmd.Flags().GetString("since")
	_, _ = cmd.Flags().GetBool("follow")
	_, _ = cmd.Flags().GetString("format")
	_, _ = cmd.Flags().GetBool("timestamps")
	_, _ = cmd.Flags().GetBool("no-color")
	_, _ = cmd.Flags().GetBool("locations")

	// Get global flags
	_, _ = cmd.Flags().GetBool("verbose")
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		_ = c.isNonInteractiveMode()
	}

	// Show logs
	err := c.ShowRecentLogs(lines, level)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to get logs."),
			c.info("Try running with --verbose to see more details"))
	}

	// Logs are displayed by ShowRecentLogs method

	return nil
}

// Template command execution methods

func (c *CLI) runTemplateInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	variables, _ := cmd.Flags().GetBool("variables")
	dependencies, _ := cmd.Flags().GetBool("dependencies")
	compatibility, _ := cmd.Flags().GetBool("compatibility")

	// Get template info
	info, err := c.templateManager.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to get template information."),
			c.info("Try running with --verbose to see more details"))
	}

	// Output template information
	if !c.quietMode {
		c.QuietOutput("📋 Template Information: %s", info.Name)
		c.QuietOutput("Description: %s", info.Description)
		c.QuietOutput("Category: %s", info.Category)
		c.QuietOutput("Technology: %s", info.Technology)

		if detailed {
			c.QuietOutput("Version: %s", info.Version)
			c.QuietOutput("Author: %s", info.Metadata.Author)
			c.QuietOutput("License: %s", info.Metadata.License)
		}

		if variables {
			c.QuietOutput("\nVariables: Not available in current template info")
		}

		if dependencies {
			c.QuietOutput("\nDependencies:")
			for _, dep := range info.Dependencies {
				c.QuietOutput("  %s", dep)
			}
		}

		if compatibility {
			c.QuietOutput("\nCompatibility: Not available in current template info")
		}
	}

	return nil
}

func (c *CLI) runTemplateValidate(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	_, _ = cmd.Flags().GetBool("fix")
	_, _ = cmd.Flags().GetString("output-format")

	// Validate template
	result, err := c.templateManager.ValidateCustomTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.error("Failed to validate template."),
			c.info("Try running with --verbose to see more details"))
	}

	// Output validation results
	if !c.quietMode {
		c.QuietOutput("🔍 Template Validation Results:")
		c.QuietOutput("Valid: %t", result.Valid)
		c.QuietOutput("Issues: %d", len(result.Issues))

		if detailed && len(result.Issues) > 0 {
			c.QuietOutput("\nIssues:")
			for _, issue := range result.Issues {
				c.QuietOutput("  - %s: %s", issue.Severity, issue.Message)
			}
		}
	}

	return nil
}
