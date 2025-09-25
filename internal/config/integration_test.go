package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestConfigManagerIntegration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create config manager
	configPath := filepath.Join(tempDir, "config.yaml")
	manager := NewManager(tempDir, "")

	t.Run("CreateDefaultConfig", func(t *testing.T) {
		err := manager.CreateDefaultConfig(configPath)
		if err != nil {
			t.Fatalf("CreateDefaultConfig failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatalf("Config file was not created")
		}
	})

	t.Run("LoadAndValidateConfig", func(t *testing.T) {
		// Ensure config file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err := manager.CreateDefaultConfig(configPath)
			if err != nil {
				t.Fatalf("CreateDefaultConfig failed: %v", err)
			}
		}

		config, err := manager.LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Validate the loaded config
		err = manager.ValidateConfig(config)
		if err != nil {
			t.Fatalf("ValidateConfig failed: %v", err)
		}
	})

	t.Run("ValidateConfigFromFile", func(t *testing.T) {
		// Ensure config file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err := manager.CreateDefaultConfig(configPath)
			if err != nil {
				t.Fatalf("CreateDefaultConfig failed: %v", err)
			}
		}

		result, err := manager.ValidateConfigFromFile(configPath)
		if err != nil {
			t.Fatalf("ValidateConfigFromFile failed: %v", err)
		}

		if !result.Valid {
			t.Fatalf("Config should be valid, but got %d errors", len(result.Errors))
		}
	})

	t.Run("BackupAndRestore", func(t *testing.T) {
		// Create backup
		err := manager.BackupConfig(configPath)
		if err != nil {
			t.Fatalf("BackupConfig failed: %v", err)
		}

		// Modify the original file
		config := &models.ProjectConfig{
			Name:         "modified-project",
			Organization: "modified-org",
			Author:       "Modified Author",
			License:      "Apache-2.0",
			OutputPath:   "./modified-output",
		}
		err = manager.SaveConfig(config, configPath)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Verify modification
		loadedConfig, err := manager.LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if loadedConfig.Name != "modified-project" {
			t.Fatalf("Config was not modified correctly")
		}

		// Find backup file
		backupFiles, err := filepath.Glob(configPath + ".backup_*")
		if err != nil || len(backupFiles) == 0 {
			t.Fatalf("Backup file not found")
		}

		// Restore from backup
		err = manager.RestoreConfig(backupFiles[0])
		if err != nil {
			t.Fatalf("RestoreConfig failed: %v", err)
		}

		// Verify restoration
		restoredConfig, err := manager.LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig after restore failed: %v", err)
		}
		if restoredConfig.Name == "modified-project" {
			t.Fatalf("Config was not restored correctly")
		}
	})

	t.Run("SettingsManagement", func(t *testing.T) {
		// Set required settings first
		err := manager.SetSetting("name", "test-project")
		if err != nil {
			t.Fatalf("SetSetting for name failed: %v", err)
		}

		err = manager.SetSetting("organization", "test-org")
		if err != nil {
			t.Fatalf("SetSetting for organization failed: %v", err)
		}

		// Set a custom setting
		err = manager.SetSetting("test.key", "test.value")
		if err != nil {
			t.Fatalf("SetSetting failed: %v", err)
		}

		// Get the setting
		value, err := manager.GetSetting("test.key")
		if err != nil {
			t.Fatalf("GetSetting failed: %v", err)
		}

		if value != "test.value" {
			t.Fatalf("Expected 'test.value', got '%v'", value)
		}

		// Validate settings
		err = manager.ValidateSettings()
		if err != nil {
			t.Fatalf("ValidateSettings failed: %v", err)
		}
	})

	t.Run("EnvironmentVariables", func(t *testing.T) {
		// Set environment variables
		if err := os.Setenv("GENERATOR_NAME", "env-project"); err != nil {
			t.Fatalf("Failed to set GENERATOR_NAME: %v", err)
		}
		if err := os.Setenv("GENERATOR_ORGANIZATION", "env-org"); err != nil {
			t.Fatalf("Failed to set GENERATOR_ORGANIZATION: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("GENERATOR_NAME"); err != nil {
				t.Logf("Failed to unset GENERATOR_NAME: %v", err)
			}
			if err := os.Unsetenv("GENERATOR_ORGANIZATION"); err != nil {
				t.Logf("Failed to unset GENERATOR_ORGANIZATION: %v", err)
			}
		}()

		// Load from environment
		envConfig, err := manager.LoadFromEnvironment()
		if err != nil {
			t.Fatalf("LoadFromEnvironment failed: %v", err)
		}

		if envConfig.Name != "env-project" {
			t.Fatalf("Expected 'env-project', got '%s'", envConfig.Name)
		}
		if envConfig.Organization != "env-org" {
			t.Fatalf("Expected 'env-org', got '%s'", envConfig.Organization)
		}

		// Test environment variable loading
		envVars := manager.LoadEnvironmentVariables()
		if len(envVars) == 0 {
			t.Fatalf("Expected environment variables to be loaded")
		}
	})

	t.Run("ConfigurationMerging", func(t *testing.T) {
		config1 := &models.ProjectConfig{
			Name:         "project1",
			Organization: "org1",
			License:      "MIT",
		}

		config2 := &models.ProjectConfig{
			Name:        "project2", // This should override
			Description: "A test project",
			Author:      "Test Author",
		}

		merged := manager.MergeConfigurations(config1, config2)

		// Check that config2 values override config1
		if merged.Name != "project2" {
			t.Fatalf("Expected 'project2', got '%s'", merged.Name)
		}
		// Check that config1 values are preserved when not overridden
		if merged.Organization != "org1" {
			t.Fatalf("Expected 'org1', got '%s'", merged.Organization)
		}
		if merged.License != "MIT" {
			t.Fatalf("Expected 'MIT', got '%s'", merged.License)
		}
		// Check that new values from config2 are added
		if merged.Description != "A test project" {
			t.Fatalf("Expected 'A test project', got '%s'", merged.Description)
		}
		if merged.Author != "Test Author" {
			t.Fatalf("Expected 'Test Author', got '%s'", merged.Author)
		}
	})
}
