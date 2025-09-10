package cli

import (
	"fmt"
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
func TestPreviewConfiguration(t *testing.T) {
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
				Home:    true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			Go:     "1.22.0",
			NextJS: "15.0.0",
			React:  "18.2.0",
			Kotlin: "2.0.0",
			Swift:  "5.9.0",
		},
	}

	// This test ensures the method doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PreviewConfiguration panicked: %v", r)
		}
	}()

	cli.PreviewConfiguration(config)
}

func TestSetSelectedComponentsInvalidInput(t *testing.T) {
	cli := &CLI{}
	config := &models.ProjectConfig{}

	// Test with invalid component format
	invalidSelections := []string{
		"invalid.component",
		"frontend.invalid",
		"backend.invalid",
		"mobile.invalid",
		"infrastructure.invalid",
		"completely.invalid.format",
		"",
	}

	for _, selection := range invalidSelections {
		t.Run("invalid_"+selection, func(t *testing.T) {
			config.Components = models.Components{} // Reset components

			err := cli.setSelectedComponents(config, []string{selection})
			// Invalid components should be ignored, not cause errors
			if err != nil {
				t.Logf("setSelectedComponents with invalid input '%s' returned error: %v", selection, err)
			}
		})
	}
}

func TestGetDefaultVersionsConsistency(t *testing.T) {
	cli := &CLI{}

	// Call getDefaultVersions multiple times and ensure consistency
	versions1 := cli.getDefaultVersions()
	versions2 := cli.getDefaultVersions()

	if versions1.Node != versions2.Node {
		t.Error("Node version should be consistent across calls")
	}
	if versions1.Go != versions2.Go {
		t.Error("Go version should be consistent across calls")
	}
	if versions1.NextJS != versions2.NextJS {
		t.Error("NextJS version should be consistent across calls")
	}
	if versions1.React != versions2.React {
		t.Error("React version should be consistent across calls")
	}
}

func TestCLIWithNilDependencies(t *testing.T) {
	// Test CLI with nil dependencies
	cli := NewCLI(nil, nil)

	if cli == nil {
		t.Fatal("NewCLI with nil dependencies returned nil")
	}

	// Methods should handle nil dependencies gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CLI methods panicked with nil dependencies: %v", r)
		}
	}()

	cli.ShowProgress("Test")
	cli.ShowSuccess("Test")
	cli.ShowError("Test")
	cli.ShowWarning("Test")

	versions := cli.getDefaultVersions()
	if versions == nil {
		t.Error("getDefaultVersions should not return nil even with nil dependencies")
	}
}

func TestCLIErrorHandling(t *testing.T) {
	// Test CLI with mock that returns errors
	testErr := fmt.Errorf("test error")
	configManager := &MockConfigManager{err: testErr}
	validator := &MockValidationEngine{err: testErr}

	cli := NewCLI(configManager, validator)

	// Test that CLI handles errors gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CLI panicked when dependencies return errors: %v", r)
		}
	}()

	// These operations might fail but shouldn't panic
	cli.ShowProgress("Test with error dependencies")
}

func TestSetSelectedComponentsEdgeCases(t *testing.T) {
	cli := &CLI{}
	config := &models.ProjectConfig{}

	testCases := []struct {
		name     string
		selected []string
	}{
		{
			name:     "empty selection",
			selected: []string{},
		},
		{
			name:     "nil selection",
			selected: nil,
		},
		{
			name: "duplicate selections",
			selected: []string{
				"frontend.main_app - Main Next.js application",
				"frontend.main_app - Main Next.js application",
			},
		},
		{
			name: "mixed valid and invalid",
			selected: []string{
				"frontend.main_app - Main Next.js application",
				"invalid.component",
				"backend.api - Go API server with Gin framework",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Components = models.Components{} // Reset components

			err := cli.setSelectedComponents(config, tc.selected)
			if err != nil {
				t.Logf("setSelectedComponents returned error for %s: %v", tc.name, err)
			}
		})
	}
}

func TestCheckOutputPathEdgeCases(t *testing.T) {
	cli := &CLI{}

	testCases := []struct {
		name string
		path string
	}{
		{
			name: "empty path",
			path: "",
		},
		{
			name: "current directory",
			path: ".",
		},
		{
			name: "parent directory",
			path: "..",
		},
		{
			name: "root directory",
			path: "/",
		},
		{
			name: "relative path",
			path: "./test/path",
		},
		{
			name: "path with spaces",
			path: "/tmp/path with spaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cli.CheckOutputPath(tc.path)
			// Most paths should not error, but we're testing that the method doesn't panic
			if err != nil {
				t.Logf("CheckOutputPath returned error for %s: %v", tc.name, err)
			}
		})
	}
}

func TestCLIMethodsWithComplexConfig(t *testing.T) {
	cli := &CLI{}

	// Create a complex configuration
	config := &models.ProjectConfig{
		Name:         "complex-project-name-with-dashes",
		Organization: "complex-org-name",
		Description:  "A very long description that might contain special characters !@#$%^&*()_+{}|:<>?[]\\;'\",./ and unicode characters: 你好世界",
		License:      "Apache-2.0",
		Author:       "Author Name with Spaces",
		Email:        "complex.email+tag@subdomain.example.com",
		Repository:   "https://github.com/complex-org/complex-repo-name",
		OutputPath:   "/very/long/path/to/output/directory/that/might/not/exist",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    true,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			Go:     "1.22.0",
			NextJS: "15.0.0",
			React:  "18.2.0",
			Kotlin: "2.0.0",
			Swift:  "5.9.0",
			Packages: map[string]string{
				"typescript":    "5.3.0",
				"tailwindcss":   "3.4.0",
				"gin-gonic/gin": "1.9.1",
				"gorm.io/gorm":  "1.25.5",
			},
			UpdatedAt: time.Now(),
		},
		CustomVars: map[string]string{
			"custom_var_1": "value1",
			"custom_var_2": "value2",
		},
	}

	// Test that all CLI methods handle complex configuration without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CLI methods panicked with complex config: %v", r)
		}
	}()

	cli.PreviewConfiguration(config)
	cli.ConfirmGeneration(config)
	cli.CheckOutputPath(config.OutputPath)
}

func TestCLIWithMockDependenciesReturningData(t *testing.T) {
	// Test CLI with mocks that return actual data
	versions := &models.VersionConfig{
		Node:   "18.19.0",
		Go:     "1.21.5",
		NextJS: "14.1.0",
		React:  "18.2.0",
		Packages: map[string]string{
			"typescript": "5.3.0",
		},
		UpdatedAt: time.Now(),
	}

	configManager := &MockConfigManager{versions: versions}

	validationResult := &models.ValidationResult{
		Valid: false,
		Errors: []models.ValidationError{
			{Field: "name", Message: "Name is required"},
			{Field: "organization", Message: "Organization is required"},
		},
		Warnings: []models.ValidationWarning{
			{Field: "description", Message: "Description is recommended"},
		},
	}

	validator := &MockValidationEngine{result: validationResult}

	cli := NewCLI(configManager, validator)

	// Test that CLI works with mocks returning data
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CLI panicked with mock data: %v", r)
		}
	}()

	cli.ShowProgress("Testing with mock data")
}
