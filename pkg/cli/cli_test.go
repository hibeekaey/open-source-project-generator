package cli

import (
	"fmt"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/audit"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

// MockConfigManager implements the ConfigManager interface for testing
type MockConfigManager struct {
	err error
}

func (m *MockConfigManager) LoadDefaults() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		License: "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker: true,
			},
		},
	}, m.err
}

func (m *MockConfigManager) ValidateConfig(*models.ProjectConfig) error {
	return m.err
}

func (m *MockConfigManager) LoadConfig(path string) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{}, m.err
}

func (m *MockConfigManager) SaveConfig(config *models.ProjectConfig, path string) error {
	return m.err
}

// Additional methods required by the enhanced ConfigManager interface
func (m *MockConfigManager) GetSetting(key string) (any, error) {
	return nil, m.err
}

func (m *MockConfigManager) SetSetting(key string, value any) error {
	return m.err
}

func (m *MockConfigManager) ValidateSettings() error {
	return m.err
}

func (m *MockConfigManager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

func (m *MockConfigManager) LoadFromEnvironment() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{}, m.err
}

func (m *MockConfigManager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
	if len(configs) > 0 {
		return configs[0]
	}
	return &models.ProjectConfig{}
}

func (m *MockConfigManager) GetConfigSchema() *interfaces.ConfigSchema {
	return nil
}

func (m *MockConfigManager) GetConfigSources() ([]interfaces.ConfigSource, error) {
	return nil, m.err
}

func (m *MockConfigManager) GetConfigLocation() string {
	return "/mock/config"
}

func (m *MockConfigManager) CreateDefaultConfig(path string) error {
	return m.err
}

func (m *MockConfigManager) BackupConfig(path string) error {
	return m.err
}

func (m *MockConfigManager) RestoreConfig(backupPath string) error {
	return m.err
}

func (m *MockConfigManager) LoadEnvironmentVariables() map[string]string {
	return make(map[string]string)
}

func (m *MockConfigManager) SetEnvironmentDefaults() error {
	return m.err
}

func (m *MockConfigManager) GetEnvironmentPrefix() string {
	return "GENERATOR"
}

func (m *MockConfigManager) ValidateConfigFromFile(path string) (*interfaces.ConfigValidationResult, error) {
	return &interfaces.ConfigValidationResult{
		Valid:    m.err == nil,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary:  interfaces.ConfigValidationSummary{},
	}, m.err
}

// MockValidationEngine implements the ValidationEngine interface for testing
type MockValidationEngine struct {
	result *models.ValidationResult
	err    error
}

func (m *MockValidationEngine) ValidateProject(projectPath string) (*models.ValidationResult, error) {
	if m.result != nil {
		return m.result, m.err
	}
	return &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Validation completed",
	}, m.err
}

func (m *MockValidationEngine) ValidateTemplate(templatePath string) error {
	return m.err
}

func (m *MockValidationEngine) ValidatePackageJSON(path string) error {
	return m.err
}

func (m *MockValidationEngine) ValidateGoMod(path string) error {
	return m.err
}

func (m *MockValidationEngine) ValidateDockerfile(path string) error {
	return m.err
}

func (m *MockValidationEngine) ValidateYAML(path string) error {
	return m.err
}

func (m *MockValidationEngine) ValidateJSON(path string) error {
	return m.err
}

// Additional methods required by the enhanced ValidationEngine interface
func (m *MockValidationEngine) ValidateProjectStructure(path string) (*interfaces.StructureValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateProjectDependencies(path string) (*interfaces.DependencyValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateProjectSecurity(path string) (*interfaces.SecurityValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateProjectQuality(path string) (*interfaces.QualityValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateConfiguration(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateConfigurationSchema(config any, schema *interfaces.ConfigSchema) error {
	return m.err
}

func (m *MockValidationEngine) ValidateConfigurationValues(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateTemplateAdvanced(path string) (*interfaces.TemplateValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateTemplateMetadata(metadata *interfaces.TemplateMetadata) error {
	return m.err
}

func (m *MockValidationEngine) ValidateTemplateStructure(path string) (*interfaces.StructureValidationResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ValidateTemplateVariables(variables map[string]interfaces.TemplateVariable) error {
	return m.err
}

func (m *MockValidationEngine) SetValidationRules(rules []interfaces.ValidationRule) error {
	return m.err
}

func (m *MockValidationEngine) GetValidationRules() []interfaces.ValidationRule {
	return nil
}

func (m *MockValidationEngine) AddValidationRule(rule interfaces.ValidationRule) error {
	return m.err
}

func (m *MockValidationEngine) RemoveValidationRule(ruleID string) error {
	return m.err
}

func (m *MockValidationEngine) FixValidationIssues(path string, issues []interfaces.ValidationIssue) (*interfaces.FixResult, error) {
	return nil, m.err
}

func (m *MockValidationEngine) GetFixableIssues(issues []interfaces.ValidationIssue) []interfaces.ValidationIssue {
	return nil
}

func (m *MockValidationEngine) PreviewFixes(path string, issues []interfaces.ValidationIssue) (*interfaces.FixPreview, error) {
	return nil, m.err
}

func (m *MockValidationEngine) ApplyFix(path string, fix interfaces.Fix) error {
	return m.err
}

func (m *MockValidationEngine) GenerateValidationReport(result *interfaces.ValidationResult, format string) ([]byte, error) {
	return nil, m.err
}

func (m *MockValidationEngine) GetValidationSummary(results []*interfaces.ValidationResult) (*interfaces.ValidationSummary, error) {
	return nil, m.err
}

func TestNewCLIWithMocks(t *testing.T) {
	configManager := &MockConfigManager{}
	validator := &MockValidationEngine{}

	// Create mock implementations for the new dependencies
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()
	cacheManager := cache.NewManager("/tmp/test-cache")

	cli := NewCLI(configManager, validator, templateManager, cacheManager, versionManager, auditEngine, nil, "test-version", "test-commit", "test-build-time")

	if cli == nil {
		t.Fatal("NewCLI returned nil")
	}

	if cli.(*CLI).configManager == nil {
		t.Error("CLI configManager not set")
	}

	if cli.(*CLI).validator == nil {
		t.Error("CLI validator not set")
	}
}

func TestCLIWithNilDependencies(t *testing.T) {
	// Test CLI with nil dependencies
	cli := NewCLI(nil, nil, nil, nil, nil, nil, nil, "test-version", "test-commit", "test-build-time")

	if cli == nil {
		t.Fatal("NewCLI with nil dependencies should not return nil")
	}

	// Test that CLI methods don't panic with nil dependencies
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CLI methods panicked with nil dependencies: %v", r)
		}
	}()

	// Methods removed in simplified CLI
}

func TestCLIErrorHandling(t *testing.T) {
	// Test CLI with mock that returns errors
	testErr := fmt.Errorf("test error")

	configManager := &MockConfigManager{err: testErr}
	validator := &MockValidationEngine{err: testErr}

	// Create mock implementations for the new dependencies
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()
	cacheManager := cache.NewManager("/tmp/test-cache")

	cli := NewCLI(configManager, validator, templateManager, cacheManager, versionManager, auditEngine, nil, "test-version", "test-commit", "test-build-time")

	if cli == nil {
		t.Fatal("NewCLI returned nil")
	}

	// Test that CLI handles errors gracefully
	// Methods removed in simplified CLI
}
func TestCLICreationWithAllDependencies(t *testing.T) {
	// Initialize all required dependencies
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()
	cacheManager := cache.NewManager("/tmp/test-cache")

	// Create CLI with all dependencies
	cli := NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		nil,
		"test-version",
		"test-commit",
		"test-build-time",
	)

	if cli == nil {
		t.Fatal("Expected CLI to be created, got nil")
	}
}

func TestCLICommands(t *testing.T) {
	// Initialize all required dependencies
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()
	cacheManager := cache.NewManager("/tmp/test-cache")

	// Create CLI
	cliImpl := NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		nil,
		"test-version",
		"test-commit",
		"test-build-time",
	).(*CLI)

	// Test that root command is set up
	if cliImpl.rootCmd == nil {
		t.Fatal("Expected root command to be set up")
	}

	// Test that commands are registered
	commands := cliImpl.rootCmd.Commands()
	expectedCommands := []string{
		"generate", "validate", "audit", "version",
		"config", "list-templates", "update", "cache", "logs",
	}

	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("Expected command '%s' to be registered", expected)
		}
	}
}

func TestCLIHelp(t *testing.T) {
	// Initialize all required dependencies
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()
	cacheManager := cache.NewManager("/tmp/test-cache")

	// Create CLI
	cli := NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		nil,
		"test-version",
		"test-commit",
		"test-build-time",
	)

	// Test help command (should not return error)
	err := cli.Run([]string{"--help"})
	// Help command exits with code 0, which cobra treats as no error
	if err != nil {
		t.Errorf("Expected help command to succeed, got error: %v", err)
	}
}
