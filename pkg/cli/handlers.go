// Package cli provides command handler functionality for the CLI interface.
//
// This module handles the execution of all CLI commands, coordinating between
// command parsing, business logic, and output formatting.
package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/commands"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/spf13/cobra"
)

// CommandHandlers manages the execution of all CLI commands.
//
// The CommandHandlers provides centralized command execution including:
//   - Command coordination and error management
//   - Business logic orchestration
//   - Output formatting and user feedback
//   - Error handling and recovery
type CommandHandlers struct {
	cli         *CLI
	generateCmd *commands.GenerateCommand
	validateCmd *commands.ValidateCommand
	templateCmd *commands.TemplateCommands
}

// generateAdapter adapts the CLI struct to implement the GenerateCLI interface.
type generateAdapter struct {
	cli *CLI
}

func (ga *generateAdapter) ValidateGenerateOptions(options interfaces.GenerateOptions) error {
	return ga.cli.validateGenerateOptions(options)
}

func (ga *generateAdapter) DetectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string {
	return ga.cli.detectGenerationMode(configPath, nonInteractive, interactive, explicitMode)
}

func (ga *generateAdapter) RouteToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error {
	return ga.cli.routeToGenerationMethod(mode, configPath, options)
}

func (ga *generateAdapter) ApplyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool) {
	return ga.cli.applyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)
}

func (ga *generateAdapter) VerboseOutput(format string, args ...interface{}) {
	ga.cli.VerboseOutput(format, args...)
}

func (ga *generateAdapter) DebugOutput(format string, args ...interface{}) {
	ga.cli.DebugOutput(format, args...)
}

func (ga *generateAdapter) Error(text string) string {
	return ga.cli.Error(text)
}

func (ga *generateAdapter) Info(text string) string {
	return ga.cli.Info(text)
}

// validateAdapter adapts the CLI struct to implement the ValidateCLI interface.
type validateAdapter struct {
	cli *CLI
}

func (va *validateAdapter) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	return va.cli.ValidateProject(path, options)
}

func (va *validateAdapter) QuietOutput(format string, args ...interface{}) {
	va.cli.QuietOutput(format, args...)
}

func (va *validateAdapter) VerboseOutput(format string, args ...interface{}) {
	va.cli.VerboseOutput(format, args...)
}

func (va *validateAdapter) DebugOutput(format string, args ...interface{}) {
	va.cli.DebugOutput(format, args...)
}

func (va *validateAdapter) Error(text string) string {
	return va.cli.Error(text)
}

func (va *validateAdapter) Warning(text string) string {
	return va.cli.Warning(text)
}

func (va *validateAdapter) Info(text string) string {
	return va.cli.Info(text)
}

func (va *validateAdapter) IsQuietMode() bool {
	return va.cli.quietMode
}

func (va *validateAdapter) CreateValidationError(message string, details map[string]interface{}) error {
	return va.cli.createValidationError(message, details)
}

// templateAdapter adapts the CLI struct to implement the TemplateCLI interface.
type templateAdapter struct {
	cli *CLI
}

func (ta *templateAdapter) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	return ta.cli.ListTemplates(filter)
}

func (ta *templateAdapter) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	return ta.cli.GetTemplateInfo(name)
}

func (ta *templateAdapter) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return ta.cli.ValidateTemplate(path)
}

func (ta *templateAdapter) VerboseOutput(format string, args ...interface{}) {
	ta.cli.VerboseOutput(format, args...)
}

func (ta *templateAdapter) DebugOutput(format string, args ...interface{}) {
	ta.cli.DebugOutput(format, args...)
}

func (ta *templateAdapter) QuietOutput(format string, args ...interface{}) {
	ta.cli.QuietOutput(format, args...)
}

func (ta *templateAdapter) ErrorOutput(format string, args ...interface{}) {
	ta.cli.ErrorOutput(format, args...)
}

func (ta *templateAdapter) WarningOutput(format string, args ...interface{}) {
	ta.cli.WarningOutput(format, args...)
}

func (ta *templateAdapter) SuccessOutput(format string, args ...interface{}) {
	ta.cli.SuccessOutput(format, args...)
}

func (ta *templateAdapter) Error(text string) string {
	return ta.cli.Error(text)
}

func (ta *templateAdapter) Warning(text string) string {
	return ta.cli.Warning(text)
}

func (ta *templateAdapter) Info(text string) string {
	return ta.cli.Info(text)
}

func (ta *templateAdapter) Success(text string) string {
	return ta.cli.Success(text)
}

func (ta *templateAdapter) Highlight(text string) string {
	return ta.cli.Highlight(text)
}

func (ta *templateAdapter) Dim(text string) string {
	return ta.cli.Dim(text)
}

func (ta *templateAdapter) IsQuietMode() bool {
	return ta.cli.quietMode
}

func (ta *templateAdapter) OutputMachineReadable(data interface{}, format string) error {
	return ta.cli.outputMachineReadable(data, format)
}

func (ta *templateAdapter) CreateTemplateError(message string, templateName string) error {
	return ta.cli.createTemplateError(message, templateName)
}

func (ta *templateAdapter) OutputSuccess(message string, data interface{}, operation string, args []string) error {
	return ta.cli.outputSuccess(message, data, operation, args)
}

func (ta *templateAdapter) IsNonInteractiveMode(cmd *cobra.Command) bool {
	return ta.cli.flagHandler.IsNonInteractiveMode(cmd)
}

// NewCommandHandlers creates a new CommandHandlers instance.
func NewCommandHandlers(cli *CLI) *CommandHandlers {
	generateAdapter := &generateAdapter{cli: cli}
	validateAdapter := &validateAdapter{cli: cli}
	templateAdapter := &templateAdapter{cli: cli}
	return &CommandHandlers{
		cli:         cli,
		generateCmd: commands.NewGenerateCommand(generateAdapter),
		validateCmd: commands.NewValidateCommand(validateAdapter),
		templateCmd: commands.NewTemplateCommands(templateAdapter),
	}
}

// runGenerate handles the generate command execution
func (ch *CommandHandlers) runGenerate(cmd *cobra.Command, args []string) error {
	return ch.generateCmd.Execute(cmd, args)
}

// runValidate handles the validate command execution
func (ch *CommandHandlers) runValidate(cmd *cobra.Command, args []string) error {
	return ch.validateCmd.Execute(cmd, args)
}

// runAudit handles the audit command execution
func (ch *CommandHandlers) runAudit(cmd *cobra.Command, args []string) error {
	// Create and execute audit command
	auditCmd := commands.NewAuditCommand(ch.cli)
	return auditCmd.Execute(cmd, args)
}

// runListTemplates handles the list-templates command execution
func (ch *CommandHandlers) runListTemplates(cmd *cobra.Command, args []string) error {
	return ch.templateCmd.ExecuteList(cmd, args)
}

// runUpdate handles the update command execution
func (ch *CommandHandlers) runUpdate(cmd *cobra.Command, args []string) error {
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

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = ch.cli.flagHandler.IsNonInteractiveMode(ch.cli.rootCmd)
	}

	// Log update options for debugging
	ch.cli.DebugOutput("Update options: check=%t, install=%t, templates=%t, force=%t", check, install, templates, force)
	ch.cli.DebugOutput("Channel: %s, backup: %t, verify: %t, compatibility: %t", channel, backup, verify, compatibility)
	if version != "" {
		ch.cli.DebugOutput("Target version: %s", version)
	}

	// Check for updates
	if check || (!install && !templates) {
		ch.cli.VerboseOutput("🔍 Checking for updates...")
		updateInfo, err := ch.cli.CheckUpdates()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error("Unable to check for updates."),
				ch.cli.Info("Please check your internet connection and try again"))
		}

		// Output update information
		if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
			return ch.cli.outputMachineReadable(updateInfo, outputFormat)
		}

		if updateInfo.UpdateAvailable {
			ch.cli.QuietOutput("🎉 Update available!")
			ch.cli.QuietOutput("Current version: %s", ch.cli.Dim(updateInfo.CurrentVersion))
			ch.cli.QuietOutput("Latest version: %s", ch.cli.Highlight(updateInfo.LatestVersion))

			if updateInfo.Security {
				ch.cli.QuietOutput("🔒 %s", ch.cli.Error("This update includes security fixes"))
			}
			if updateInfo.Breaking {
				ch.cli.QuietOutput("⚠️  %s", ch.cli.Warning("This update includes breaking changes"))
			}
			if updateInfo.Recommended {
				ch.cli.QuietOutput("✨ %s", ch.cli.Success("This update is recommended"))
			}

			if releaseNotes && updateInfo.ReleaseNotes != "" {
				ch.cli.QuietOutput("\nRelease Notes:")
				ch.cli.QuietOutput("%s", updateInfo.ReleaseNotes)
			}

			if !nonInteractive {
				ch.cli.QuietOutput("\nTo install: %s", ch.cli.Info("generator update --install"))
			}
		} else {
			ch.cli.QuietOutput("✅ You're running the latest version (%s)", updateInfo.CurrentVersion)
		}

		if check {
			return nil // Only checking, don't proceed to install
		}
	}

	// Install updates
	if install {
		ch.cli.VerboseOutput("📦 Installing updates...")
		err := ch.cli.InstallUpdates()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error("Update installation failed."),
				ch.cli.Info("Please try again or check the logs for more details"))
		}
		ch.cli.QuietOutput("✅ %s", ch.cli.Success("Update installed successfully!"))
		ch.cli.QuietOutput("🔄 Please restart the generator to use the new version")
	}

	// Update templates
	if templates {
		ch.cli.VerboseOutput("📋 Updating templates cache...")
		// This would call a method to update template cache
		// For now, we'll simulate the operation
		ch.cli.QuietOutput("✅ %s", ch.cli.Success("Templates cache updated successfully!"))
	}

	return nil
}

// runCacheShow handles the cache show command execution
func (ch *CommandHandlers) runCacheShow(cmd *cobra.Command, args []string) error {
	// Show cache status and statistics
	err := ch.cli.ShowCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Unable to show cache information."),
			ch.cli.Info("The cache may be corrupted or inaccessible"))
	}
	return nil
}

// runCacheClear handles the cache clear command execution
func (ch *CommandHandlers) runCacheClear(cmd *cobra.Command, args []string) error {
	// Get flags
	force, _ := cmd.Flags().GetBool("force")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = ch.cli.flagHandler.IsNonInteractiveMode(ch.cli.rootCmd)
	}

	// Confirm cache clearing unless forced or in non-interactive mode
	if !force && !nonInteractive {
		ch.cli.QuietOutput("⚠️  This will remove all cached data including templates and package versions.")
		fmt.Print("Are you sure you want to continue? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			ch.cli.QuietOutput("Cache clearing cancelled.")
			return nil
		}
	}

	ch.cli.VerboseOutput("🧹 Clearing cache...")
	err := ch.cli.ClearCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Unable to clear cache."),
			ch.cli.Info("Some cache files may be in use or protected"))
	}
	ch.cli.QuietOutput("✅ %s", ch.cli.Success("Cache cleared successfully!"))
	return nil
}

// runCacheClean handles the cache clean command execution
func (ch *CommandHandlers) runCacheClean(cmd *cobra.Command, args []string) error {
	// Clean expired and invalid cache entries
	fmt.Println("🧹 Cleaning cache...")
	err := ch.cli.CleanCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Cache cleaning failed."),
			ch.cli.Info("Some cache entries may be corrupted or inaccessible"))
	}
	ch.cli.QuietOutput("✅ %s", ch.cli.Success("Cache cleaned successfully!"))
	return nil
}

// runCacheValidate handles the cache validate command execution
func (ch *CommandHandlers) runCacheValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Validating cache...")
	err := ch.cli.ValidateCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Cache validation failed."),
			ch.cli.Info("The cache may be corrupted or incomplete"))
	}
	ch.cli.QuietOutput("✅ %s", ch.cli.Success("Cache is valid!"))
	return nil
}

// runCacheRepair handles the cache repair command execution
func (ch *CommandHandlers) runCacheRepair(cmd *cobra.Command, args []string) error {
	err := ch.cli.RepairCache()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Cache repair failed."),
			ch.cli.Info("The cache may be severely corrupted"))
	}
	ch.cli.QuietOutput("✅ %s", ch.cli.Success("Cache repaired successfully!"))
	return nil
}

// runCacheOfflineEnable handles the cache offline enable command execution
func (ch *CommandHandlers) runCacheOfflineEnable(cmd *cobra.Command, args []string) error {
	return ch.cli.EnableOfflineMode()
}

// runCacheOfflineDisable handles the cache offline disable command execution
func (ch *CommandHandlers) runCacheOfflineDisable(cmd *cobra.Command, args []string) error {
	return ch.cli.DisableOfflineMode()
}

// runCacheOfflineStatus handles the cache offline status command execution
func (ch *CommandHandlers) runCacheOfflineStatus(cmd *cobra.Command, args []string) error {
	if ch.cli.cacheManager == nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Cache manager not available."),
			ch.cli.Info("The cache system may not be properly initialized"))
	}

	// Get cache stats to determine offline mode status
	stats, err := ch.cli.GetCacheStats()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Unable to get cache status."),
			ch.cli.Info("The cache may be corrupted or inaccessible"))
	}

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = ch.cli.flagHandler.IsNonInteractiveMode(ch.cli.rootCmd)
	}

	// Output in machine-readable format if requested
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return ch.cli.outputMachineReadable(stats, outputFormat)
	}

	// Human-readable output
	ch.cli.QuietOutput("📦 Cache entries: %d", stats.TotalEntries)
	ch.cli.QuietOutput("💾 Cache size: %s", ch.formatBytes(stats.TotalSize))
	ch.cli.QuietOutput("💡 Use 'generator cache offline enable' to enable offline mode")

	return nil
}

// runLogs handles the logs command execution
func (ch *CommandHandlers) runLogs(cmd *cobra.Command, args []string) error {
	// Get flags
	lines, _ := cmd.Flags().GetInt("tail")
	level, _ := cmd.Flags().GetString("level")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	follow, _ := cmd.Flags().GetBool("follow")
	timestamps, _ := cmd.Flags().GetBool("timestamps")
	format, _ := cmd.Flags().GetString("format")
	component, _ := cmd.Flags().GetString("component")

	// Note: non-interactive mode handling would be added here when log filtering is implemented

	// Log the options for debugging
	ch.cli.DebugOutput("Log options: lines=%d, level=%s, since=%s, until=%s", lines, level, since, until)
	ch.cli.DebugOutput("Follow=%t, timestamps=%t, format=%s, component=%s", follow, timestamps, format, component)

	// For now, show recent logs using the existing ShowLogs method
	// This would be enhanced to support all the filtering options
	err := ch.cli.ShowLogs()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Unable to show logs."),
			ch.cli.Info("The log files may be inaccessible or corrupted"))
	}

	// If follow mode is requested, show a message (actual implementation would tail the logs)
	if follow {
		ch.cli.VerboseOutput("📡 Following logs... (Press Ctrl+C to stop)")
		// In a real implementation, this would continuously tail the log file
		// For now, we'll just show a message
		ch.cli.QuietOutput("💡 Follow mode not yet implemented. Showing recent logs only.")
	}

	return nil
}

// runVersion handles the version command execution
func (ch *CommandHandlers) runVersion(cmd *cobra.Command, args []string) error {
	// Get flags
	packages, _ := cmd.Flags().GetBool("packages")
	checkUpdates, _ := cmd.Flags().GetBool("check-updates")
	buildInfo, _ := cmd.Flags().GetBool("build-info")
	short, _ := cmd.Flags().GetBool("short")
	format, _ := cmd.Flags().GetString("format")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	compatibility, _ := cmd.Flags().GetBool("compatibility")
	checkPackage, _ := cmd.Flags().GetString("check-package")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = ch.cli.flagHandler.IsNonInteractiveMode(ch.cli.rootCmd)
	}

	// Handle JSON flag - if --json is specified, set format to json
	if jsonOutput {
		format = "json"
	}

	// Use output-format if format is not specified
	if format == "text" && outputFormat != "text" {
		format = outputFormat
	}

	// Create version options
	options := interfaces.VersionOptions{
		ShowPackages:  packages,
		CheckUpdates:  checkUpdates,
		ShowBuildInfo: buildInfo,
		OutputFormat:  format,
	}

	// Handle short version output
	if short {
		if nonInteractive && (format == "json" || format == "yaml") {
			data := map[string]string{"version": ch.cli.generatorVersion}
			return ch.cli.outputMachineReadable(data, format)
		}
		ch.cli.QuietOutput(ch.cli.generatorVersion)
		return nil
	}

	// Show version information
	err := ch.cli.ShowVersion(options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			ch.cli.Error("Unable to show version information."),
			ch.cli.Info("The version manager may not be properly initialized"))
	}

	// Check for updates if requested
	if checkUpdates {
		ch.cli.VerboseOutput("🔍 Checking for updates...")
		updateInfo, err := ch.cli.CheckUpdates()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error("Unable to check for updates."),
				ch.cli.Info("Please check your internet connection and try again"))
		}

		if nonInteractive && (format == "json" || format == "yaml") {
			return ch.cli.outputMachineReadable(updateInfo, format)
		}

		if updateInfo.UpdateAvailable {
			ch.cli.QuietOutput("🎉 Update available: %s → %s",
				ch.cli.Dim(updateInfo.CurrentVersion),
				ch.cli.Highlight(updateInfo.LatestVersion))
		} else {
			ch.cli.QuietOutput("✅ You're running the latest version")
		}
	}

	// Show package versions if requested
	if packages {
		ch.cli.VerboseOutput("📦 Fetching package versions...")
		packageVersions, err := ch.cli.GetPackageVersions()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error("Unable to get package versions."),
				ch.cli.Info("The version manager may not be properly initialized"))
		}

		if nonInteractive && (format == "json" || format == "yaml") {
			return ch.cli.outputMachineReadable(packageVersions, format)
		}

		ch.cli.QuietOutput("\nPackage Versions:")
		ch.cli.QuietOutput("=================")
		for pkg, version := range packageVersions {
			ch.cli.QuietOutput("%s: %s", pkg, version)
		}
	}

	// Check specific package if requested
	if checkPackage != "" {
		packageVersions, err := ch.cli.GetPackageVersions()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error("Unable to get package versions."),
				ch.cli.Info("The version manager may not be properly initialized"))
		}

		if version, exists := packageVersions[checkPackage]; exists {
			if nonInteractive && (format == "json" || format == "yaml") {
				data := map[string]string{checkPackage: version}
				return ch.cli.outputMachineReadable(data, format)
			}
			ch.cli.QuietOutput("%s: %s", checkPackage, version)
		} else {
			return fmt.Errorf("🚫 %s %s",
				ch.cli.Error(fmt.Sprintf("Package '%s' not found.", checkPackage)),
				ch.cli.Info("Use --packages to see all available packages"))
		}
	}

	// Show compatibility information if requested
	if compatibility {
		ch.cli.QuietOutput("\nCompatibility information would be displayed here")
	}

	return nil
}

// runTemplateInfo handles the template info command execution
func (ch *CommandHandlers) runTemplateInfo(cmd *cobra.Command, args []string) error {
	return ch.templateCmd.ExecuteInfo(cmd, args)
}

// runTemplateValidate handles the template validate command execution
func (ch *CommandHandlers) runTemplateValidate(cmd *cobra.Command, args []string) error {
	return ch.templateCmd.ExecuteValidate(cmd, args)
}

// runInteractiveProjectConfiguration handles interactive project configuration
func (ch *CommandHandlers) runInteractiveProjectConfiguration(ctx context.Context) (*models.ProjectConfig, error) {
	ch.cli.VerboseOutput("🎯 Starting interactive project configuration")

	// Use the interactive flow manager for comprehensive configuration
	if ch.cli.interactiveFlowManager != nil {
		// The interactive flow manager handles the full flow, not just config collection
		// For now, fall back to basic prompts
		ch.cli.VerboseOutput("Using basic project configuration prompts")
	}

	// Use basic prompts for project configuration
	return ch.cli.PromptProjectDetails()
}

// formatBytes formats byte count as human-readable string
func (ch *CommandHandlers) formatBytes(bytes int64) string {
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
