package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

// Manager implements the ConfigManager interface
type Manager struct {
	cacheDir     string
	defaultsPath string
}

// NewManager creates a new configuration manager
func NewManager(cacheDir, defaultsPath string) interfaces.ConfigManager {
	return &Manager{
		cacheDir:     cacheDir,
		defaultsPath: defaultsPath,
	}
}

// LoadDefaults loads default configuration values
func (m *Manager) LoadDefaults() (*models.ProjectConfig, error) {
	// Try to load from defaults file first
	if m.defaultsPath != "" && fileExists(m.defaultsPath) {
		config, err := m.LoadConfig(m.defaultsPath)
		if err == nil {
			return config, nil
		}
		// If loading fails, fall back to hardcoded defaults
	}

	// Return hardcoded defaults
	return &models.ProjectConfig{
		License: "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   false,
					Admin:  false,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Mobile: models.MobileComponents{
				Android: false,
				IOS:     false,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: false,
				Terraform:  false,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"react":      "18.2.0",
				"next":       "13.4.0",
				"typescript": "5.0.0",
			},
		},
	}, nil
}

// LoadConfig loads configuration from a file
func (m *Manager) LoadConfig(configPath string) (*models.ProjectConfig, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(configPath); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	content, err := utils.SafeReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.ProjectConfig
	ext := strings.ToLower(filepath.Ext(configPath))

	switch ext {
	case ".json":
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func (m *Manager) SaveConfig(config *models.ProjectConfig, configPath string) error {
	ext := strings.ToLower(filepath.Ext(configPath))
	var content []byte
	var err error

	switch ext {
	case ".json":
		content, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}
	case ".yaml", ".yml":
		content, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err := utils.SafeWriteFile(configPath, content); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateConfig validates a project configuration
func (m *Manager) ValidateConfig(config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if config.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if config.Organization == "" {
		return fmt.Errorf("organization is required")
	}

	if config.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	// Normalize license - set default if empty or unsupported
	m.normalizeLicense(config)

	return nil
}

// normalizeLicense normalizes the license field, setting default for unsupported types
func (m *Manager) normalizeLicense(config *models.ProjectConfig) {
	supportedLicenses := []string{
		"MIT",
		"Apache-2.0",
		"GPL-3.0",
		"BSD-3-Clause",
	}

	// Set default if empty
	if config.License == "" {
		config.License = "MIT"
		return
	}

	// Check if license is supported
	for _, supported := range supportedLicenses {
		if config.License == supported {
			return // License is supported
		}
	}

	// License not supported - set to default
	config.License = "MIT"
}

// Helper function to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
