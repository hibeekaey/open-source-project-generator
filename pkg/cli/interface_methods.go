// Package cli provides interface method implementations for backward compatibility.
package cli

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Legacy methods for backward compatibility (to be removed in future refactoring)
func (c *CLI) validateGenerateOptions(options interfaces.GenerateOptions) error {
	return c.inputValidator.ValidateGenerateOptions(options)
}

func (c *CLI) detectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string {
	return "interactive"
}

func (c *CLI) routeToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) applyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool) {
	return nonInteractive, interactive
}

// Placeholder methods for interface compatibility (to be implemented in future tasks)
func (c *CLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) CheckUpdates() (*interfaces.UpdateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) InstallUpdates() error     { return fmt.Errorf("not implemented") }
func (c *CLI) ShowCache() error          { return fmt.Errorf("not implemented") }
func (c *CLI) ClearCache() error         { return fmt.Errorf("not implemented") }
func (c *CLI) CleanCache() error         { return fmt.Errorf("not implemented") }
func (c *CLI) ValidateCache() error      { return fmt.Errorf("not implemented") }
func (c *CLI) RepairCache() error        { return fmt.Errorf("not implemented") }
func (c *CLI) EnableOfflineMode() error  { return fmt.Errorf("not implemented") }
func (c *CLI) DisableOfflineMode() error { return fmt.Errorf("not implemented") }

func (c *CLI) GetCacheStats() (*interfaces.CacheStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) ShowLogs() error { return fmt.Errorf("not implemented") }

func (c *CLI) ShowVersion(options interfaces.VersionOptions) error {
	return fmt.Errorf("not implemented")
}

func (c *CLI) GetPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("not implemented")
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

func (c *CLI) GetLogLevel() string {
	return "info"
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
func (c *CLI) SetLogLevel(level string) error { return fmt.Errorf("not implemented") }
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
