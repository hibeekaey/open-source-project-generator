package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

func TestNewManager(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
}

func TestLoadDefaults(t *testing.T) {
	manager := NewManager("", "")

	config, err := manager.LoadDefaults()
	if err != nil {
		t.Fatalf("LoadDefaults failed: %v", err)
	}

	if config == nil {
		t.Fatal("LoadDefaults returned nil")
	}

	// Verify default values
	if config.License != "MIT" {
		t.Errorf("Expected default license to be MIT, got %s", config.License)
	}

	if !config.Components.Frontend.NextJS.App {
		t.Error("Expected NextJS App to be true by default")
	}

	if !config.Components.Backend.GoGin {
		t.Error("Expected GoGin to be true by default")
	}

	if !config.Components.Infrastructure.Docker {
		t.Error("Expected Docker to be true by default")
	}

	if config.Versions == nil {
		t.Fatal("Expected versions to be non-nil")
	}
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	// Create a test config file
	testConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project from defaults",
		License:      "Apache-2.0",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    false,
					Home:   true,
					Admin:  true,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
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
			Node: "18.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "17.0.0",
			},
		},
		OutputPath: "./custom-output",
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	configPath := filepath.Join(tempDir, "test-config.yaml")
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load the config
	config, err := manager.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig returned nil")
	}

	// Verify loaded values
	if config.License != "Apache-2.0" {
		t.Errorf("Expected license to be Apache-2.0, got %s", config.License)
	}

	if config.Components.Frontend.NextJS.App {
		t.Error("Expected NextJS App to be false from defaults file")
	}

	if !config.Components.Frontend.NextJS.Home {
		t.Error("Expected NextJS Home to be true from defaults file")
	}

	if config.Versions.Node != "18.0.0" {
		t.Errorf("Expected Node version to be 18.0.0, got %s", config.Versions.Node)
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, "")

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.22.0",
			Packages: map[string]string{
				"next": "15.0.0",
			},
		},
		OutputPath: "./test-output",
	}

	configPath := filepath.Join(tempDir, "test-config.yaml")
	err := manager.SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify the saved config
	loadedConfig, err := manager.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.Name != config.Name {
		t.Errorf("Expected name to be %s, got %s", config.Name, loadedConfig.Name)
	}

	if loadedConfig.License != config.License {
		t.Errorf("Expected license to be %s, got %s", config.License, loadedConfig.License)
	}
}

func TestValidateConfig(t *testing.T) {
	manager := NewManager("", "")

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
				OutputPath:   "./test-output",
				Components: models.Components{
					Frontend: models.FrontendComponents{
						NextJS: models.NextJSComponents{
							App: true,
						},
					},
					Backend: models.BackendComponents{
						GoGin: true,
					},
				},
			},
			expectError: false,
		},
		{
			name: "config with empty name",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "test-org",
				Description:  "A test project with empty name",
				License:      "MIT",
			},
			expectError: true,
		},
		{
			name: "config with empty organization",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "",
				Description:  "A test project with empty organization",
				License:      "MIT",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := manager.ValidateConfig(tc.config)

			if tc.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}
