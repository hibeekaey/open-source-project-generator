package interfaces

import "github.com/open-source-template-generator/pkg/models"

// ConfigManager defines the contract for configuration management operations
type ConfigManager interface {
	// LoadDefaults loads default configuration values
	LoadDefaults() (*models.ProjectConfig, error)

	// ValidateConfig validates the provided project configuration
	ValidateConfig(*models.ProjectConfig) error

	// GetLatestVersions fetches the latest versions of packages and frameworks
	GetLatestVersions() (*models.VersionConfig, error)

	// MergeConfigs merges base configuration with override values
	MergeConfigs(base, override *models.ProjectConfig) *models.ProjectConfig

	// SaveConfig saves configuration to a file
	SaveConfig(config *models.ProjectConfig, path string) error

	// LoadConfig loads configuration from a file
	LoadConfig(path string) (*models.ProjectConfig, error)
}
