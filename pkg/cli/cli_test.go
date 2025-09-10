package cli

import (
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// MockConfigManager implements the ConfigManager interface for testing
type MockConfigManager struct {
	versions *models.VersionConfig
	err      error
}

func (m *MockConfigManager) LoadDefaults() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		License: "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
			Backend: models.BackendComponents{
				API: true,
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

func (m *MockConfigManager) GetLatestVersions() (*models.VersionConfig, error) {
	if m.versions != nil {
		return m.versions, m.err
	}
	return &models.VersionConfig{
		Node:      "20.11.0",
		Go:        "1.22.0",
		NextJS:    "15.0.0",
		React:     "18.2.0",
		UpdatedAt: time.Now(),
	}, m.err
}

func (m *MockConfigManager) MergeConfigs(base, override *models.ProjectConfig) *models.ProjectConfig {
	return base
}

func (m *MockConfigManager) SaveConfig(config *models.ProjectConfig, path string) error {
	return m.err
}

func (m *MockConfigManager) LoadConfig(path string) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{}, m.err
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
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}, m.err
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

	if cli.configManager == nil {
		t.Error("CLI configManager not set")
	}

	if cli.validator == nil {
		t.Error("CLI validator not set")
	}
}

func TestSetSelectedComponents(t *testing.T) {
	cli := &CLI{}
	config := &models.ProjectConfig{}

	testCases := []struct {
		name     string
		selected []string
		expected models.Components
	}{
		{
			name: "frontend components",
			selected: []string{
				"frontend.main_app - Main Next.js application",
				"frontend.admin - Admin dashboard application",
			},
			expected: models.Components{
				Frontend: models.FrontendComponents{
					MainApp: true,
					Admin:   true,
				},
			},
		},
		{
			name: "backend components",
			selected: []string{
				"backend.api - Go API server with Gin framework",
			},
			expected: models.Components{
				Backend: models.BackendComponents{
					API: true,
				},
			},
		},
		{
			name: "mobile components",
			selected: []string{
				"mobile.android - Android Kotlin application",
				"mobile.ios - iOS Swift application",
			},
			expected: models.Components{
				Mobile: models.MobileComponents{
					Android: true,
					IOS:     true,
				},
			},
		},
		{
			name: "infrastructure components",
			selected: []string{
				"infrastructure.docker - Docker configurations",
				"infrastructure.kubernetes - Kubernetes manifests",
				"infrastructure.terraform - Terraform configurations",
			},
			expected: models.Components{
				Infrastructure: models.InfrastructureComponents{
					Docker:     true,
					Kubernetes: true,
					Terraform:  true,
				},
			},
		},
		{
			name: "mixed components",
			selected: []string{
				"frontend.main_app - Main Next.js application",
				"backend.api - Go API server with Gin framework",
				"infrastructure.docker - Docker configurations",
			},
			expected: models.Components{
				Frontend: models.FrontendComponents{
					MainApp: true,
				},
				Backend: models.BackendComponents{
					API: true,
				},
				Infrastructure: models.InfrastructureComponents{
					Docker: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Components = models.Components{} // Reset components

			err := cli.setSelectedComponents(config, tc.selected)
			if err != nil {
				t.Fatalf("setSelectedComponents returned error: %v", err)
			}

			// Check frontend components
			if config.Components.Frontend.MainApp != tc.expected.Frontend.MainApp {
				t.Errorf("Frontend.MainApp: expected %v, got %v", tc.expected.Frontend.MainApp, config.Components.Frontend.MainApp)
			}
			if config.Components.Frontend.Home != tc.expected.Frontend.Home {
				t.Errorf("Frontend.Home: expected %v, got %v", tc.expected.Frontend.Home, config.Components.Frontend.Home)
			}
			if config.Components.Frontend.Admin != tc.expected.Frontend.Admin {
				t.Errorf("Frontend.Admin: expected %v, got %v", tc.expected.Frontend.Admin, config.Components.Frontend.Admin)
			}

			// Check backend components
			if config.Components.Backend.API != tc.expected.Backend.API {
				t.Errorf("Backend.API: expected %v, got %v", tc.expected.Backend.API, config.Components.Backend.API)
			}

			// Check mobile components
			if config.Components.Mobile.Android != tc.expected.Mobile.Android {
				t.Errorf("Mobile.Android: expected %v, got %v", tc.expected.Mobile.Android, config.Components.Mobile.Android)
			}
			if config.Components.Mobile.IOS != tc.expected.Mobile.IOS {
				t.Errorf("Mobile.IOS: expected %v, got %v", tc.expected.Mobile.IOS, config.Components.Mobile.IOS)
			}

			// Check infrastructure components
			if config.Components.Infrastructure.Docker != tc.expected.Infrastructure.Docker {
				t.Errorf("Infrastructure.Docker: expected %v, got %v", tc.expected.Infrastructure.Docker, config.Components.Infrastructure.Docker)
			}
			if config.Components.Infrastructure.Kubernetes != tc.expected.Infrastructure.Kubernetes {
				t.Errorf("Infrastructure.Kubernetes: expected %v, got %v", tc.expected.Infrastructure.Kubernetes, config.Components.Infrastructure.Kubernetes)
			}
			if config.Components.Infrastructure.Terraform != tc.expected.Infrastructure.Terraform {
				t.Errorf("Infrastructure.Terraform: expected %v, got %v", tc.expected.Infrastructure.Terraform, config.Components.Infrastructure.Terraform)
			}
		})
	}
}

func TestGetDefaultVersions(t *testing.T) {
	cli := &CLI{}
	versions := cli.getDefaultVersions()

	if versions == nil {
		t.Fatal("getDefaultVersions returned nil")
	}

	if versions.Node == "" {
		t.Error("Node version should not be empty")
	}

	if versions.Go == "" {
		t.Error("Go version should not be empty")
	}

	if versions.NextJS == "" {
		t.Error("NextJS version should not be empty")
	}

	if versions.React == "" {
		t.Error("React version should not be empty")
	}

	if versions.Packages == nil {
		t.Error("Packages map should be initialized")
	}

	// Check that UpdatedAt is recent (within last minute)
	if time.Since(versions.UpdatedAt) > time.Minute {
		t.Error("UpdatedAt should be recent")
	}
}

func TestCheckOutputPath(t *testing.T) {
	cli := &CLI{}

	// Test with non-existent path (should not error)
	err := cli.CheckOutputPath("/tmp/non-existent-test-path-12345")
	if err != nil {
		t.Errorf("CheckOutputPath should not error for non-existent path: %v", err)
	}
}

func TestProgressAndStatusMethods(t *testing.T) {
	cli := &CLI{}

	// These methods should not panic
	cli.ShowProgress("Test progress")
	cli.ShowSuccess("Test success")
	cli.ShowError("Test error")
	cli.ShowWarning("Test warning")
}

func TestConfirmGeneration(t *testing.T) {
	cli := &CLI{}
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test description",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/test/repo",
		OutputPath:   "/tmp/test-output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			Go:     "1.22.0",
			NextJS: "15.0.0",
			React:  "18.2.0",
		},
	}

	// This test just ensures the method doesn't panic
	// In a real scenario, this would require user interaction
	// For automated testing, we would need to mock the survey input

	// We can't easily test the interactive confirmation without mocking survey
	// but we can test that the method exists and doesn't panic with valid input
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConfirmGeneration panicked: %v", r)
		}
	}()

	// This will attempt to show the configuration but won't wait for user input in tests
	// The method will return false since there's no interactive terminal
	result := cli.ConfirmGeneration(config)
	_ = result // We can't predict the result in a test environment
}
