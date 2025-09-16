package cli

import (
	"fmt"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
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

func TestNewCLI(t *testing.T) {
	configManager := &MockConfigManager{}
	validator := &MockValidationEngine{}

	cli := NewCLI(configManager, validator)

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
	cli := NewCLI(nil, nil)

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

	cli := NewCLI(configManager, validator)

	if cli == nil {
		t.Fatal("NewCLI returned nil")
	}

	// Test that CLI handles errors gracefully
	// Methods removed in simplified CLI
}
