package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

func TestManager_LoadDefaults(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	config, err := manager.LoadDefaults()
	if err != nil {
		t.Fatalf("LoadDefaults failed: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be non-nil")
	}

	// Verify default values
	if config.License != "MIT" {
		t.Errorf("Expected default license to be MIT, got %s", config.License)
	}

	if !config.Components.Frontend.MainApp {
		t.Error("Expected MainApp to be true by default")
	}

	if !config.Components.Backend.API {
		t.Error("Expected API to be true by default")
	}

	if !config.Components.Infrastructure.Docker {
		t.Error("Expected Docker to be true by default")
	}

	if config.Versions == nil {
		t.Fatal("Expected versions to be non-nil")
	}

	if config.Versions.Node == "" {
		t.Error("Expected Node version to be set")
	}

	if config.Versions.Go == "" {
		t.Error("Expected Go version to be set")
	}
}

func TestManager_LoadDefaultsFromFile(t *testing.T) {
	tempDir := t.TempDir()
	defaultsFile := filepath.Join(tempDir, "defaults.yaml")

	// Create a defaults file
	defaultConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project from defaults",
		License:      "Apache-2.0",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: false,
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
			Node:   "18.0.0",
			Go:     "1.21.0",
			NextJS: "14.0.0",
			React:  "17.0.0",
		},
		OutputPath: "./custom-output",
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		t.Fatalf("Failed to marshal default config: %v", err)
	}

	if err := os.WriteFile(defaultsFile, data, 0644); err != nil {
		t.Fatalf("Failed to write defaults file: %v", err)
	}

	manager := NewManager(tempDir, defaultsFile)
	config, err := manager.LoadDefaults()
	if err != nil {
		t.Fatalf("LoadDefaults failed: %v", err)
	}

	// Verify loaded values
	if config.License != "Apache-2.0" {
		t.Errorf("Expected license to be Apache-2.0, got %s", config.License)
	}

	if config.Components.Frontend.MainApp {
		t.Error("Expected MainApp to be false from defaults file")
	}

	if !config.Components.Frontend.Home {
		t.Error("Expected Home to be true from defaults file")
	}

	if config.Versions.Node != "18.0.0" {
		t.Errorf("Expected Node version to be 18.0.0, got %s", config.Versions.Node)
	}
}

func TestManager_ValidateConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "A valid test project",
				License:      "MIT",
				Components: models.Components{
					Frontend: models.FrontendComponents{MainApp: true},
					Backend:  models.BackendComponents{},
					Mobile:   models.MobileComponents{},
					Infrastructure: models.InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &models.VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			},
			expectError: false,
		},
		{
			name: "invalid config - missing name",
			config: &models.ProjectConfig{
				Organization: "test-org",
				Description:  "A test project",
				License:      "MIT",
				Components: models.Components{
					Frontend: models.FrontendComponents{MainApp: true},
					Backend:  models.BackendComponents{},
					Mobile:   models.MobileComponents{},
					Infrastructure: models.InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &models.VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateConfig(tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected validation error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestManager_GetLatestVersions(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	versions, err := manager.GetLatestVersions()
	if err != nil {
		t.Fatalf("GetLatestVersions failed: %v", err)
	}

	if versions == nil {
		t.Fatal("Expected versions to be non-nil")
	}

	if versions.Node == "" {
		t.Error("Expected Node version to be set")
	}

	if versions.Go == "" {
		t.Error("Expected Go version to be set")
	}

	if versions.Packages == nil {
		t.Error("Expected packages map to be non-nil")
	}

	// Test caching
	versions2, err := manager.GetLatestVersions()
	if err != nil {
		t.Fatalf("Second GetLatestVersions call failed: %v", err)
	}

	if versions2.Node != versions.Node {
		t.Error("Expected cached versions to match")
	}
}

func TestManager_MergeConfigs(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	base := &models.ProjectConfig{
		Name:         "base-project",
		Organization: "base-org",
		Description:  "Base project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   false,
			},
			Backend: models.BackendComponents{
				API: true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "18.0.0",
			Go:     "1.21.0",
			NextJS: "14.0.0",
		},
		CustomVars: map[string]string{
			"base_var": "base_value",
		},
		OutputPath: "./base-output",
	}

	override := &models.ProjectConfig{
		Name:        "override-project",
		Description: "Override project description",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				Home:  true,
				Admin: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.0.0",
			NextJS: "15.0.0",
			React:  "18.0.0",
		},
		CustomVars: map[string]string{
			"override_var": "override_value",
		},
		OutputPath: "./override-output",
	}

	merged := manager.MergeConfigs(base, override)

	// Test basic field merging
	if merged.Name != "override-project" {
		t.Errorf("Expected name to be 'override-project', got %s", merged.Name)
	}

	if merged.Organization != "base-org" {
		t.Errorf("Expected organization to remain 'base-org', got %s", merged.Organization)
	}

	if merged.Description != "Override project description" {
		t.Errorf("Expected description to be overridden, got %s", merged.Description)
	}

	// Test component merging
	if !merged.Components.Frontend.MainApp {
		t.Error("Expected MainApp to remain true from base")
	}

	if !merged.Components.Frontend.Home {
		t.Error("Expected Home to be true from override")
	}

	if !merged.Components.Frontend.Admin {
		t.Error("Expected Admin to be true from override")
	}

	if !merged.Components.Backend.API {
		t.Error("Expected API to remain true from base")
	}

	if !merged.Components.Mobile.Android {
		t.Error("Expected Android to be true from override")
	}

	// Test version merging
	if merged.Versions.Node != "20.0.0" {
		t.Errorf("Expected Node version to be overridden to '20.0.0', got %s", merged.Versions.Node)
	}

	if merged.Versions.Go != "1.21.0" {
		t.Errorf("Expected Go version to remain '1.21.0', got %s", merged.Versions.Go)
	}

	if merged.Versions.React != "18.0.0" {
		t.Errorf("Expected React version to be '18.0.0' from override, got %s", merged.Versions.React)
	}

	// Test custom vars merging
	if merged.CustomVars["base_var"] != "base_value" {
		t.Error("Expected base_var to be preserved")
	}

	if merged.CustomVars["override_var"] != "override_value" {
		t.Error("Expected override_var to be added")
	}

	if merged.OutputPath != "./override-output" {
		t.Errorf("Expected output path to be overridden, got %s", merged.OutputPath)
	}
}

func TestManager_SaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "A test project for save/load",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     false,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  false,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.0.0",
			Go:     "1.22.0",
			Kotlin: "2.0.0",
			NextJS: "15.0.0",
			React:  "18.0.0",
			Packages: map[string]string{
				"express": "4.18.0",
				"lodash":  "4.17.21",
			},
			UpdatedAt: time.Now(),
		},
		CustomVars: map[string]string{
			"custom_var1": "value1",
			"custom_var2": "value2",
		},
		OutputPath:       "./test-output",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}

	tests := []struct {
		name     string
		filename string
	}{
		{"YAML format", "config.yaml"},
		{"JSON format", "config.json"},
		{"YML format", "config.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, tt.filename)

			// Test saving
			err := manager.SaveConfig(config, configPath)
			if err != nil {
				t.Fatalf("SaveConfig failed: %v", err)
			}

			// Verify file exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Fatal("Config file was not created")
			}

			// Test loading
			loadedConfig, err := manager.LoadConfig(configPath)
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}

			// Verify loaded config matches original
			if loadedConfig.Name != config.Name {
				t.Errorf("Expected name %s, got %s", config.Name, loadedConfig.Name)
			}

			if loadedConfig.Organization != config.Organization {
				t.Errorf("Expected organization %s, got %s", config.Organization, loadedConfig.Organization)
			}

			if loadedConfig.Components.Frontend.MainApp != config.Components.Frontend.MainApp {
				t.Errorf("Expected MainApp %v, got %v", config.Components.Frontend.MainApp, loadedConfig.Components.Frontend.MainApp)
			}

			if loadedConfig.Versions.Node != config.Versions.Node {
				t.Errorf("Expected Node version %s, got %s", config.Versions.Node, loadedConfig.Versions.Node)
			}

			if len(loadedConfig.CustomVars) != len(config.CustomVars) {
				t.Errorf("Expected %d custom vars, got %d", len(config.CustomVars), len(loadedConfig.CustomVars))
			}

			for k, v := range config.CustomVars {
				if loadedConfig.CustomVars[k] != v {
					t.Errorf("Expected custom var %s=%s, got %s", k, v, loadedConfig.CustomVars[k])
				}
			}
		})
	}
}

func TestManager_LoadConfigNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	_, err := manager.LoadConfig(filepath.Join(tempDir, "nonexistent.yaml"))
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}

func TestManager_MergeConfigsNilHandling(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	config := &models.ProjectConfig{
		Name: "test",
	}

	// Test nil base
	result := manager.MergeConfigs(nil, config)
	if result != config {
		t.Error("Expected override config when base is nil")
	}

	// Test nil override
	result = manager.MergeConfigs(config, nil)
	if result != config {
		t.Error("Expected base config when override is nil")
	}

	// Test both nil
	result = manager.MergeConfigs(nil, nil)
	if result != nil {
		t.Error("Expected nil when both configs are nil")
	}
}

func TestManager_VersionCaching(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	// Create a cached version file
	cacheFile := filepath.Join(tempDir, "versions.json")
	cachedVersions := &models.VersionConfig{
		Node:      "19.0.0",
		Go:        "1.21.5",
		UpdatedAt: time.Now().Add(-1 * time.Hour), // 1 hour old
	}

	data, err := json.MarshalIndent(cachedVersions, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal cached versions: %v", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	// Get versions (should use cache)
	versions, err := manager.GetLatestVersions()
	if err != nil {
		t.Fatalf("GetLatestVersions failed: %v", err)
	}

	if versions.Node != "19.0.0" {
		t.Errorf("Expected cached Node version 19.0.0, got %s", versions.Node)
	}

	// Test expired cache
	expiredVersions := &models.VersionConfig{
		Node:      "18.0.0",
		Go:        "1.20.0",
		UpdatedAt: time.Now().Add(-25 * time.Hour), // 25 hours old (expired)
	}

	data, err = json.MarshalIndent(expiredVersions, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal expired versions: %v", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		t.Fatalf("Failed to write expired cache file: %v", err)
	}

	// Get versions (should fetch new ones, not use expired cache)
	versions, err = manager.GetLatestVersions()
	if err != nil {
		t.Fatalf("GetLatestVersions failed: %v", err)
	}

	// Should get default versions, not expired cached ones
	if versions.Node == "18.0.0" {
		t.Error("Should not use expired cache")
	}
}

func TestManager_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("invalid cache directory permissions", func(t *testing.T) {
		// Create a directory with no write permissions
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := os.MkdirAll(restrictedDir, 0444); err != nil {
			t.Fatalf("Failed to create restricted directory: %v", err)
		}
		defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup

		manager := NewManager(restrictedDir, "")

		// Should still work but may not cache successfully
		versions, err := manager.GetLatestVersions()
		if err != nil {
			t.Errorf("GetLatestVersions should work even with cache write issues: %v", err)
		}
		if versions == nil {
			t.Error("Expected versions to be returned even with cache issues")
		}
	})

	t.Run("corrupted cache file", func(t *testing.T) {
		cacheDir := filepath.Join(tempDir, "corrupted")
		os.MkdirAll(cacheDir, 0755)

		// Create corrupted cache file
		cacheFile := filepath.Join(cacheDir, "versions.json")
		os.WriteFile(cacheFile, []byte("corrupted json {"), 0644)

		manager := NewManager(cacheDir, "")

		// Should fallback to default versions
		versions, err := manager.GetLatestVersions()
		if err != nil {
			t.Errorf("GetLatestVersions should handle corrupted cache gracefully: %v", err)
		}
		if versions == nil {
			t.Error("Expected default versions when cache is corrupted")
		}
	})

	t.Run("empty cache directory", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty")
		os.MkdirAll(emptyDir, 0755)

		manager := NewManager(emptyDir, "")

		versions, err := manager.GetLatestVersions()
		if err != nil {
			t.Errorf("GetLatestVersions should work with empty cache dir: %v", err)
		}
		if versions == nil {
			t.Error("Expected versions to be returned")
		}
	})

	t.Run("invalid defaults file path", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "non-existent", "defaults.yaml")
		manager := NewManager(tempDir, invalidPath)

		// Should fallback to hardcoded defaults
		config, err := manager.LoadDefaults()
		if err != nil {
			t.Errorf("LoadDefaults should fallback gracefully: %v", err)
		}
		if config == nil {
			t.Error("Expected default config to be returned")
		}
		if config.License != "MIT" {
			t.Error("Expected hardcoded default license")
		}
	})

	t.Run("malformed defaults file", func(t *testing.T) {
		defaultsFile := filepath.Join(tempDir, "malformed.yaml")
		os.WriteFile(defaultsFile, []byte("invalid: yaml: content: ["), 0644)

		manager := NewManager(tempDir, defaultsFile)

		// Should fallback to hardcoded defaults
		config, err := manager.LoadDefaults()
		if err != nil {
			t.Errorf("LoadDefaults should handle malformed file gracefully: %v", err)
		}
		if config.License != "MIT" {
			t.Error("Expected fallback to hardcoded defaults")
		}
	})
}

func TestManager_ConfigValidationEdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	t.Run("config with unicode characters", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "测试项目",
			Organization: "тест-орг",
			Description:  "Проект с unicode символами 🚀",
			License:      "MIT",
			Components: models.Components{
				Frontend:       models.FrontendComponents{MainApp: true},
				Infrastructure: models.InfrastructureComponents{Docker: true},
			},
			Versions: &models.VersionConfig{
				Node: "20.0.0",
				Go:   "1.22.0",
			},
			OutputPath: "./output",
		}

		err := manager.ValidateConfig(config)
		// Should handle unicode gracefully (may pass or fail depending on validation rules)
		if err != nil {
			t.Logf("Unicode config validation result: %v", err)
		}
	})

	t.Run("config with extremely long values", func(t *testing.T) {
		longString := strings.Repeat("a", 1000)
		config := &models.ProjectConfig{
			Name:         longString,
			Organization: longString,
			Description:  longString,
			License:      "MIT",
			Components: models.Components{
				Frontend:       models.FrontendComponents{MainApp: true},
				Infrastructure: models.InfrastructureComponents{Docker: true},
			},
			Versions: &models.VersionConfig{
				Node: "20.0.0",
				Go:   "1.22.0",
			},
			OutputPath: "./output",
		}

		err := manager.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation to fail for extremely long values")
		}
	})

	t.Run("config with special characters in paths", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "test-project",
			Organization: "test-org",
			Description:  "Test project",
			License:      "MIT",
			Components: models.Components{
				Frontend:       models.FrontendComponents{MainApp: true},
				Infrastructure: models.InfrastructureComponents{Docker: true},
			},
			Versions: &models.VersionConfig{
				Node: "20.0.0",
				Go:   "1.22.0",
			},
			OutputPath: "./output/../../../etc/passwd",
		}

		err := manager.ValidateConfig(config)
		// Note: Current validation may not catch path traversal - this is a test for future enhancement
		if err != nil {
			t.Logf("Path validation result: %v", err)
		} else {
			t.Logf("Path validation passed - consider adding path safety validation")
		}
	})
}

func TestManager_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	// Test concurrent access to version cache
	t.Run("concurrent version fetching", func(t *testing.T) {
		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := manager.GetLatestVersions()
				results <- err
			}()
		}

		for i := 0; i < numGoroutines; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent version fetching failed: %v", err)
			}
		}
	})

	// Test concurrent config operations
	t.Run("concurrent config save/load", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "concurrent-test",
			Organization: "test-org",
			Description:  "Concurrent access test",
			License:      "MIT",
			Components: models.Components{
				Frontend:       models.FrontendComponents{MainApp: true},
				Infrastructure: models.InfrastructureComponents{Docker: true},
			},
			Versions: &models.VersionConfig{
				Node: "20.0.0",
				Go:   "1.22.0",
			},
			OutputPath: "./output",
		}

		const numGoroutines = 5
		results := make(chan error, numGoroutines*2)

		// Concurrent saves
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				configPath := filepath.Join(tempDir, fmt.Sprintf("config-%d.yaml", index))
				results <- manager.SaveConfig(config, configPath)
			}(i)
		}

		// Concurrent loads (after a brief delay)
		time.Sleep(10 * time.Millisecond)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				configPath := filepath.Join(tempDir, fmt.Sprintf("config-%d.yaml", index))
				_, err := manager.LoadConfig(configPath)
				results <- err
			}(i)
		}

		for i := 0; i < numGoroutines*2; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent config operation failed: %v", err)
			}
		}
	})
}
