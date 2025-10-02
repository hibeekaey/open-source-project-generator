// Package cli provides interface method implementations for backward compatibility.
package cli

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Legacy methods for backward compatibility (to be removed in future refactoring)
func (c *CLI) validateGenerateOptions(options interfaces.GenerateOptions) error {
	return c.inputValidator.ValidateGenerateOptions(options)
}

func (c *CLI) detectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string {
	// Priority 1: Explicit mode flag (highest priority)
	if explicitMode != "" {
		return c.validateAndNormalizeMode(explicitMode)
	}

	// Priority 2: Direct mode flags
	if nonInteractive {
		return "non-interactive"
	}
	if interactive {
		return "interactive"
	}

	// Priority 3: Configuration file presence
	if configPath != "" {
		return "config-file"
	}

	// Priority 4: Environment detection
	if c.helper.DetectNonInteractiveMode(c.rootCmd) {
		return "non-interactive"
	}

	// Default: Interactive mode
	return "interactive"
}

func (c *CLI) routeToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🚀 Starting %s generation...", mode)

	switch mode {
	case "interactive":
		return c.executeInteractiveGeneration(options)
	case "non-interactive":
		return c.executeNonInteractiveGeneration(configPath, options)
	case "config-file":
		return c.executeConfigFileGeneration(configPath, options)
	default:
		return fmt.Errorf("🚫 %s %s",
			c.Error(fmt.Sprintf("Unknown generation mode: %s", mode)),
			c.Info("Valid modes are: interactive, non-interactive, config-file"))
	}
}

func (c *CLI) applyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool) {
	// Handle explicit mode
	if explicitMode != "" {
		normalizedMode := c.validateAndNormalizeMode(explicitMode)
		switch normalizedMode {
		case "interactive":
			return false, true
		case "non-interactive":
			return true, false
		case "config-file":
			return true, false // Config file mode is non-interactive
		}
	}

	// Handle force flags (highest priority)
	if forceNonInteractive {
		return true, false
	}
	if forceInteractive {
		return false, true
	}

	// Return original flags if no overrides
	return nonInteractive, interactive
}

// Template management methods
func (c *CLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	if c.templateManager == nil {
		return nil, fmt.Errorf("template manager not initialized")
	}

	return c.templateManager.ListTemplates(filter)
}

func (c *CLI) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	if c.templateManager == nil {
		return nil, fmt.Errorf("template manager not initialized")
	}

	return c.templateManager.GetTemplateInfo(name)
}

func (c *CLI) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	if c.templateManager == nil {
		return nil, fmt.Errorf("template manager not initialized")
	}

	return c.templateManager.ValidateTemplate(path)
}

func (c *CLI) CheckUpdates() (*interfaces.UpdateInfo, error) {
	if c.versionManager == nil {
		return nil, fmt.Errorf("version manager not initialized")
	}
	return c.versionManager.CheckForUpdates()
}

func (c *CLI) InstallUpdates() error {
	if c.versionManager == nil {
		return fmt.Errorf("version manager not initialized")
	}

	// Check for updates first
	updateInfo, err := c.versionManager.CheckForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !updateInfo.UpdateAvailable {
		c.QuietOutput("✅ %s", c.Success("Already running the latest version"))
		return nil
	}

	c.VerboseOutput("📦 Installing update from %s to %s...",
		updateInfo.CurrentVersion, updateInfo.LatestVersion)

	// Install the update
	if err := c.versionManager.InstallUpdate(updateInfo.LatestVersion); err != nil {
		return fmt.Errorf("failed to install update: %w", err)
	}

	c.QuietOutput("✅ %s", c.Success("Update installed successfully!"))
	c.QuietOutput("🔄 Please restart the generator to use the new version")

	return nil
}

func (c *CLI) ShowVersion(options interfaces.VersionOptions) error {
	if c.versionManager == nil {
		return fmt.Errorf("🚫 %s %s",
			c.Error("Version manager not initialized."),
			c.Info("This is an internal error - please report this issue"))
	}

	// Get current version info
	currentVersion := c.versionManager.GetCurrentVersion()

	// Handle different output formats
	if options.OutputFormat == "json" {
		return c.showVersionJSON(options)
	}

	// Show basic version information
	c.QuietOutput("🚀 %s %s", c.Highlight("Open Source Project Generator"), c.Success(fmt.Sprintf("v%s", currentVersion)))

	// Show build information if requested
	if options.ShowBuildInfo {
		version, gitCommit, buildTime := c.GetBuildInfo()
		c.QuietOutput("Version: %s", c.Info(version))
		c.QuietOutput("Git Commit: %s", c.Dim(gitCommit))
		c.QuietOutput("Build Time: %s", c.Dim(buildTime))
	}

	// Show package versions if requested
	if options.ShowPackages {
		c.VerboseOutput("📦 Fetching package versions...")
		packages, err := c.versionManager.GetAllPackageVersions()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				c.Error("Failed to fetch package versions."),
				c.Info("Check your internet connection or try --offline mode"))
		}

		c.QuietOutput("\n%s", c.Highlight("📦 Package Versions:"))
		c.QuietOutput("%s", c.Dim("=================="))

		// Group packages by type
		nodePackages := make(map[string]string)
		goPackages := make(map[string]string)
		systemPackages := make(map[string]string)

		for pkg, version := range packages {
			if pkg == "node" || pkg == "go" {
				systemPackages[pkg] = version
			} else if strings.Contains(pkg, "/") || strings.Contains(pkg, ".") {
				goPackages[pkg] = version
			} else {
				nodePackages[pkg] = version
			}
		}

		// Show system packages
		if len(systemPackages) > 0 {
			c.QuietOutput("\n%s", c.Info("System:"))
			for pkg, version := range systemPackages {
				c.QuietOutput("  %s: %s", pkg, c.Success(version))
			}
		}

		// Show Node.js packages
		if len(nodePackages) > 0 {
			c.QuietOutput("\n%s", c.Info("Node.js Packages:"))
			for pkg, version := range nodePackages {
				c.QuietOutput("  %s: %s", pkg, c.Success(version))
			}
		}

		// Show Go packages
		if len(goPackages) > 0 {
			c.QuietOutput("\n%s", c.Info("Go Modules:"))
			for pkg, version := range goPackages {
				c.QuietOutput("  %s: %s", pkg, c.Success(version))
			}
		}
	}

	// Check for updates if requested
	if options.CheckUpdates {
		c.VerboseOutput("🔍 Checking for updates...")
		updateInfo, err := c.versionManager.CheckForUpdates()
		if err != nil {
			c.WarningOutput("⚠️  Could not check for updates: %v", err)
			return nil
		}

		if updateInfo.UpdateAvailable {
			c.QuietOutput("\n🎉 %s", c.Success("Update available!"))
			c.QuietOutput("Current: %s → Latest: %s",
				c.Dim(updateInfo.CurrentVersion),
				c.Highlight(updateInfo.LatestVersion))

			if updateInfo.Security {
				c.QuietOutput("🔒 %s", c.Error("This update includes security fixes"))
			}
			if updateInfo.Breaking {
				c.QuietOutput("⚠️  %s", c.Warning("This update includes breaking changes"))
			}
			if updateInfo.Recommended {
				c.QuietOutput("✨ %s", c.Success("This update is recommended"))
			}
		} else {
			c.QuietOutput("\n✅ %s", c.Success("You're running the latest version"))
		}
	}

	return nil
}

func (c *CLI) GetPackageVersions() (map[string]string, error) {
	if c.versionManager == nil {
		return nil, fmt.Errorf("version manager not initialized")
	}
	return c.versionManager.GetAllPackageVersions()
}

func (c *CLI) GetLatestPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) CheckCompatibility(path string) (*interfaces.CompatibilityResult, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ExportConfig(path string) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) PromptAdvancedOptions() (*interfaces.AdvancedOptions, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ConfirmAdvancedGeneration(config *models.ProjectConfig, options *interfaces.AdvancedOptions) bool {
	return false
}

func (c *CLI) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ShowConfig() error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) SetConfig(key, value string) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) EditConfig() error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) ValidateConfig() error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) RunNonInteractive(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) GenerateReport(reportType string, format string, outputFile string) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ValidateConfigurationSchema(config *models.ProjectConfig) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) GetConfigurationSources() ([]interfaces.ConfigSource, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) executeGenerationWorkflow(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	return c.workflowHandler.ExecuteGenerationWorkflow(config, options)
}

// executeInteractiveGeneration handles interactive project generation
func (c *CLI) executeInteractiveGeneration(options interfaces.GenerateOptions) error {
	c.VerboseOutput("🎯 Starting interactive project configuration...")

	// Collect project configuration interactively
	config, err := c.PromptProjectDetails()
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.Error("Failed to collect project configuration."),
			c.Info("Please check your input and try again"))
	}

	// Confirm generation with user
	if !c.ConfirmGeneration(config) {
		c.QuietOutput("Project generation cancelled by user.")
		return nil
	}

	// Execute generation workflow
	return c.executeGenerationWorkflow(config, options)
}

// executeNonInteractiveGeneration handles non-interactive project generation
func (c *CLI) executeNonInteractiveGeneration(configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("🤖 Starting non-interactive project generation...")

	var config *models.ProjectConfig
	var err error

	if configPath != "" {
		// Load configuration from file
		config, err = c.workflowHandler.LoadConfigFromFile(configPath)
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				c.Error("Failed to load configuration file."),
				c.Info("Please check the file path and format"))
		}
	} else {
		// Load configuration from environment or defaults
		config, err = c.workflowHandler.LoadConfigFromEnvironment()
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				c.Error("Failed to load configuration from environment."),
				c.Info("Please provide a configuration file or set environment variables"))
		}
	}

	// Execute generation workflow
	return c.executeGenerationWorkflow(config, options)
}

// executeConfigFileGeneration handles config-file based project generation
func (c *CLI) executeConfigFileGeneration(configPath string, options interfaces.GenerateOptions) error {
	c.VerboseOutput("📄 Starting config-file based project generation...")

	if configPath == "" {
		return fmt.Errorf("🚫 %s %s",
			c.Error("Configuration file path is required for config-file mode."),
			c.Info("Please provide a configuration file using --config flag"))
	}

	// Load configuration from file
	config, err := c.workflowHandler.LoadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.Error("Failed to load configuration file."),
			c.Info("Please check the file path and format"))
	}

	// Execute generation workflow
	return c.executeGenerationWorkflow(config, options)
}

func (c *CLI) ShowRecentLogs(lines int, level string) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) GetLogFileLocations() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) GetTemplateDependencies(name string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) LoadConfiguration(sources []string) (*models.ProjectConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

// Additional missing methods
func (c *CLI) SaveConfiguration(config *models.ProjectConfig, name string) error {
	return fmt.Errorf("not implemented")
}
func (c *CLI) DeleteConfiguration(name string) error { return fmt.Errorf("not implemented") }
func (c *CLI) ListConfigurations() ([]string, error) { return nil, fmt.Errorf("not implemented") }
func (c *CLI) ImportConfiguration(path string) error { return fmt.Errorf("not implemented") }
func (c *CLI) CreateTemplate(name string, config interface{}) error {
	return fmt.Errorf("not implemented")
}
func (c *CLI) UpdateTemplate(name string, config interface{}) error {
	return fmt.Errorf("not implemented")
}
func (c *CLI) DeleteTemplate(name string) error    { return fmt.Errorf("not implemented") }
func (c *CLI) InstallTemplate(source string) error { return fmt.Errorf("not implemented") }
func (c *CLI) UninstallTemplate(name string) error { return fmt.Errorf("not implemented") }
func (c *CLI) RefreshTemplates() error             { return fmt.Errorf("not implemented") }
func (c *CLI) GetTemplateSource(name string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
func (c *CLI) SetTemplateSource(name string, source string) error {
	return fmt.Errorf("not implemented")
}
func (c *CLI) ValidateTemplateSource(source string) error { return fmt.Errorf("not implemented") }
func (c *CLI) GetCacheSize() (int64, error)               { return 0, fmt.Errorf("not implemented") }
func (c *CLI) GetCacheLocation() (string, error)          { return "", fmt.Errorf("not implemented") }
func (c *CLI) SetCacheLocation(path string) error         { return fmt.Errorf("not implemented") }
func (c *CLI) GetOfflineMode() (bool, error)              { return false, fmt.Errorf("not implemented") }
func (c *CLI) SetOfflineMode(enabled bool) error          { return fmt.Errorf("not implemented") }
func (c *CLI) SyncCache() error                           { return fmt.Errorf("not implemented") }
func (c *CLI) BackupCache() error                         { return fmt.Errorf("not implemented") }
func (c *CLI) RestoreCache(path string) error             { return fmt.Errorf("not implemented") }

// showVersionJSON shows version information in JSON format
func (c *CLI) showVersionJSON(options interfaces.VersionOptions) error {
	version, gitCommit, buildTime := c.GetBuildInfo()

	versionData := map[string]interface{}{
		"version":    version,
		"git_commit": gitCommit,
		"build_time": buildTime,
	}

	// Add package versions if requested
	if options.ShowPackages {
		packages, err := c.versionManager.GetAllPackageVersions()
		if err != nil {
			return fmt.Errorf("failed to fetch package versions: %w", err)
		}
		versionData["packages"] = packages
	}

	// Add update information if requested
	if options.CheckUpdates {
		updateInfo, err := c.versionManager.CheckForUpdates()
		if err != nil {
			c.WarningOutput("⚠️  Could not check for updates: %v", err)
		} else {
			versionData["update_info"] = updateInfo
		}
	}

	return c.outputMachineReadable(versionData, "json")
}
