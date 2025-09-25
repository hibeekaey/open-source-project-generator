package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Manager implements the ConfigManager interface using unified config management
type Manager struct {
	unified *config.UnifiedConfigManager
}

// NewManager creates a new configuration manager
func NewManager(cacheDir, defaultsPath string) interfaces.ConfigManager {
	unified := config.NewUnifiedConfigManager(cacheDir, defaultsPath)
	return &Manager{
		unified: unified,
	}
}

// LoadDefaults loads default configuration values
func (m *Manager) LoadDefaults() (*models.ProjectConfig, error) {
	return m.unified.LoadDefaults()
}

// ValidateConfig validates a configuration
func (m *Manager) ValidateConfig(config *models.ProjectConfig) error {
	result, err := m.unified.ValidateConfig(config)
	if err != nil {
		return err
	}
	if !result.Valid {
		return fmt.Errorf("configuration validation failed: %d errors", result.Summary.ErrorCount)
	}
	return nil
}

// SaveConfig saves a configuration to a file
func (m *Manager) SaveConfig(config *models.ProjectConfig, path string) error {
	return m.unified.SaveToFile(config, path)
}

// LoadConfig loads a configuration from a file
func (m *Manager) LoadConfig(path string) (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{}
	err := m.unified.LoadFromFile(path, config)
	return config, err
}

// GetSetting gets a setting value
func (m *Manager) GetSetting(key string) (interface{}, error) {
	return m.unified.GetSetting(key)
}

// SetSetting sets a setting value
func (m *Manager) SetSetting(key string, value interface{}) error {
	return m.unified.SetSetting(key, value)
}

// ValidateSettings validates all settings
func (m *Manager) ValidateSettings() error {
	// This would validate all current settings
	// For now, just return nil as the unified system handles validation
	return nil
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

// LoadFromEnvironment loads configuration from environment variables
func (m *Manager) LoadFromEnvironment() (*models.ProjectConfig, error) {
	config := m.unified.LoadFromEnvironment()
	return config, nil
}

// MergeConfigurations merges multiple configurations
func (m *Manager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
	return m.unified.MergeConfigurations(configs...)
}

// GetConfigSchema returns the configuration schema
func (m *Manager) GetConfigSchema() *interfaces.ConfigSchema {
	// Return a basic schema for now
	return &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Project Configuration",
		Description: "Configuration schema for project generation",
		Required:    []string{"name", "organization", "author"},
		Properties:  make(map[string]interfaces.PropertySchema),
	}
}

// ValidateConfigFromFile validates a configuration from a file
func (m *Manager) ValidateConfigFromFile(path string) (*interfaces.ConfigValidationResult, error) {
	config, err := m.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	return m.unified.ValidateConfig(config)
}

// GetConfigSources returns available configuration sources
func (m *Manager) GetConfigSources() ([]interfaces.ConfigSource, error) {
	// Return basic config sources
	return []interfaces.ConfigSource{
		{Type: "file"},
		{Type: "environment"},
	}, nil
}

// GetConfigLocation returns the configuration directory
func (m *Manager) GetConfigLocation() string {
	return m.unified.GetConfigLocation()
}

// CreateDefaultConfig creates a default configuration file
func (m *Manager) CreateDefaultConfig(path string) error {
	return m.unified.CreateDefaultConfig(path)
}

// BackupConfig backs up the current configuration
func (m *Manager) BackupConfig(path string) error {
	// Load current config from the path
	config, err := m.LoadConfig(path)
	if err != nil {
		return err
	}

	// Create backup filename with timestamp
	backupPath := path + ".backup_" + fmt.Sprintf("%d", time.Now().Unix()) + ".yaml"
	return m.SaveConfig(config, backupPath)
}

// RestoreConfig restores configuration from backup
func (m *Manager) RestoreConfig(backupPath string) error {
	config, err := m.LoadConfig(backupPath)
	if err != nil {
		return err
	}
	// Extract original path from backup path
	// backupPath format: original.yaml.backup_timestamp.yaml
	// We need to extract: original.yaml
	parts := strings.Split(backupPath, ".backup_")
	if len(parts) < 2 {
		return fmt.Errorf("invalid backup path format")
	}
	originalPath := parts[0]
	return m.SaveConfig(config, originalPath)
}

// LoadEnvironmentVariables loads environment variables
func (m *Manager) LoadEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)
	prefix := m.GetEnvironmentPrefix()

	// Load environment variables with the prefix
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix+"_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], prefix+"_")
				envVars[key] = parts[1]
			}
		}
	}

	return envVars
}

// SetEnvironmentDefaults sets environment variable defaults
func (m *Manager) SetEnvironmentDefaults() error {
	// This would set environment defaults
	// For now, just return nil
	return nil
}

// GetEnvironmentPrefix returns the environment variable prefix
func (m *Manager) GetEnvironmentPrefix() string {
	return "GENERATOR"
}
