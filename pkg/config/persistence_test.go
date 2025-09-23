package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

func TestConfigurationPersistence(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create persistence manager
	persistence := NewConfigurationPersistence(tempDir, nil)

	// Create test configuration
	testConfig := &SavedConfiguration{
		Name:        "test-config",
		Description: "Test configuration for unit tests",
		Version:     "1.0.0",
		Tags:        []string{"test", "unit"},
		ProjectConfig: &models.ProjectConfig{
			Name:         "test-project",
			Organization: "test-org",
			Author:       "Test Author",
			Email:        "test@example.com",
			License:      "MIT",
			Description:  "A test project",
		},
		SelectedTemplates: []TemplateSelection{
			{
				TemplateName: "go-gin",
				Category:     "backend",
				Technology:   "go",
				Version:      "1.0.0",
				Selected:     true,
				Options:      map[string]interface{}{"database": "postgresql"},
			},
		},
		GenerationSettings: &GenerationSettings{
			IncludeExamples:   true,
			IncludeTests:      true,
			IncludeDocs:       true,
			UpdateVersions:    false,
			MinimalMode:       false,
			BackupExisting:    true,
			OverwriteExisting: false,
		},
	}

	// Test saving configuration
	t.Run("SaveConfiguration", func(t *testing.T) {
		err := persistence.SaveConfiguration("test-config", testConfig)
		if err != nil {
			t.Fatalf("Failed to save configuration: %v", err)
		}

		// Check if file was created
		expectedFile := filepath.Join(tempDir, "test-config.yaml")
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Fatalf("Configuration file was not created: %s", expectedFile)
		}
	})

	// Test loading configuration
	t.Run("LoadConfiguration", func(t *testing.T) {
		loadedConfig, err := persistence.LoadConfiguration("test-config")
		if err != nil {
			t.Fatalf("Failed to load configuration: %v", err)
		}

		// Verify loaded configuration matches saved configuration
		if loadedConfig.Name != testConfig.Name {
			t.Errorf("Expected name %s, got %s", testConfig.Name, loadedConfig.Name)
		}
		if loadedConfig.ProjectConfig.Name != testConfig.ProjectConfig.Name {
			t.Errorf("Expected project name %s, got %s", testConfig.ProjectConfig.Name, loadedConfig.ProjectConfig.Name)
		}
		if len(loadedConfig.SelectedTemplates) != len(testConfig.SelectedTemplates) {
			t.Errorf("Expected %d templates, got %d", len(testConfig.SelectedTemplates), len(loadedConfig.SelectedTemplates))
		}
	})

	// Test configuration exists
	t.Run("ConfigurationExists", func(t *testing.T) {
		if !persistence.ConfigurationExists("test-config") {
			t.Error("Configuration should exist")
		}
		if persistence.ConfigurationExists("non-existent-config") {
			t.Error("Non-existent configuration should not exist")
		}
	})

	// Test listing configurations
	t.Run("ListConfigurations", func(t *testing.T) {
		configs, err := persistence.ListConfigurations(nil)
		if err != nil {
			t.Fatalf("Failed to list configurations: %v", err)
		}

		if len(configs) != 1 {
			t.Errorf("Expected 1 configuration, got %d", len(configs))
		}

		if len(configs) > 0 && configs[0].Name != "test-config" {
			t.Errorf("Expected configuration name 'test-config', got '%s'", configs[0].Name)
		}
	})

	// Test export configuration
	t.Run("ExportConfiguration", func(t *testing.T) {
		data, err := persistence.ExportConfiguration("test-config", "yaml")
		if err != nil {
			t.Fatalf("Failed to export configuration: %v", err)
		}

		if len(data) == 0 {
			t.Error("Exported data should not be empty")
		}
	})

	// Test delete configuration
	t.Run("DeleteConfiguration", func(t *testing.T) {
		err := persistence.DeleteConfiguration("test-config")
		if err != nil {
			t.Fatalf("Failed to delete configuration: %v", err)
		}

		// Verify configuration no longer exists
		if persistence.ConfigurationExists("test-config") {
			t.Error("Configuration should have been deleted")
		}
	})
}

func TestConfigurationValidator(t *testing.T) {
	validator := &ConfigurationValidator{
		nameRegex:    utils.ProjectNamePattern,
		versionRegex: utils.VersionPattern,
	}

	// Test valid configuration
	t.Run("ValidConfiguration", func(t *testing.T) {
		validConfig := &SavedConfiguration{
			Name:    "valid-config",
			Version: "1.0.0",
			ProjectConfig: &models.ProjectConfig{
				Name: "valid-project",
			},
			SelectedTemplates: []TemplateSelection{
				{
					TemplateName: "test-template",
					Category:     "backend",
					Selected:     true,
				},
			},
		}

		err := validator.ValidateConfiguration(validConfig)
		if err != nil {
			t.Errorf("Valid configuration should not have validation errors: %v", err)
		}
	})

	// Test invalid configuration
	t.Run("InvalidConfiguration", func(t *testing.T) {
		invalidConfig := &SavedConfiguration{
			Name:    "", // Empty name should be invalid
			Version: "invalid-version",
		}

		err := validator.ValidateConfiguration(invalidConfig)
		if err == nil {
			t.Error("Invalid configuration should have validation errors")
		}
	})
}

func TestConfigurationNameValidation(t *testing.T) {
	persistence := NewConfigurationPersistence("", nil)

	testCases := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{"Valid name", "my-config", false},
		{"Valid name with spaces", "My Config", false},
		{"Valid name with underscores", "my_config", false},
		{"Empty name", "", true},
		{"Too short", "a", true},
		{"Too long", "this-is-a-very-long-configuration-name-that-exceeds-the-maximum-allowed-length", true},
		{"Invalid characters", "my-config@#$", true},
		{"Reserved name", "default", true},
		{"Reserved name case insensitive", "DEFAULT", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := persistence.validateConfigurationName(tc.input)
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for input '%s', but got none", tc.input)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for input '%s', but got: %v", tc.input, err)
			}
		})
	}
}
