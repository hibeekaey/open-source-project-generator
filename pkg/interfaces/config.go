package interfaces

import "github.com/cuesoftinc/open-source-project-generator/pkg/models"

// ConfigManager defines the contract for configuration management operations
type ConfigManager interface {
	// LoadDefaults loads default configuration values
	LoadDefaults() (*models.ProjectConfig, error)

	// ValidateConfig validates the provided project configuration
	ValidateConfig(*models.ProjectConfig) error

	// SaveConfig saves configuration to a file
	SaveConfig(config *models.ProjectConfig, path string) error

	// LoadConfig loads configuration from a file
	LoadConfig(path string) (*models.ProjectConfig, error)
}
