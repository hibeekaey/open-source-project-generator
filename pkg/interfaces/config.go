package interfaces

import "github.com/cuesoftinc/open-source-project-generator/pkg/models"

// ConfigManager defines the contract for comprehensive configuration management operations
type ConfigManager interface {
	// Basic configuration operations
	LoadDefaults() (*models.ProjectConfig, error)
	ValidateConfig(*models.ProjectConfig) error
	SaveConfig(config *models.ProjectConfig, path string) error
	LoadConfig(path string) (*models.ProjectConfig, error)

	// Settings management
	GetSetting(key string) (any, error)
	SetSetting(key string, value any) error
	ValidateSettings() error

	// Configuration sources
	LoadFromFile(path string) (*models.ProjectConfig, error)
	LoadFromEnvironment() (*models.ProjectConfig, error)
	MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig

	// Configuration validation
	GetConfigSchema() *ConfigSchema
	ValidateConfigFromFile(path string) (*ConfigValidationResult, error)

	// Configuration management
	GetConfigSources() ([]ConfigSource, error)
	GetConfigLocation() string
	CreateDefaultConfig(path string) error
	BackupConfig(path string) error
	RestoreConfig(backupPath string) error

	// Environment integration
	LoadEnvironmentVariables() map[string]string
	SetEnvironmentDefaults() error
	GetEnvironmentPrefix() string
}

// Comprehensive configuration types and structures

// ConfigurationOptions defines options for configuration operations
type ConfigurationOptions struct {
	// Loading options
	Sources       []string `json:"sources"`
	IgnoreMissing bool     `json:"ignore_missing"`
	IgnoreInvalid bool     `json:"ignore_invalid"`
	MergeStrategy string   `json:"merge_strategy"` // override, merge, append

	// Validation options
	StrictValidation bool     `json:"strict_validation"`
	AllowUnknown     bool     `json:"allow_unknown"`
	ValidateSchema   bool     `json:"validate_schema"`
	Rules            []string `json:"rules"`

	// Output options
	Format       string `json:"format"` // yaml, json, toml
	Indent       int    `json:"indent"`
	SortKeys     bool   `json:"sort_keys"`
	IncludeEmpty bool   `json:"include_empty"`

	// Environment options
	EnvPrefix    string `json:"env_prefix"`
	EnvSeparator string `json:"env_separator"`
	EnvTransform string `json:"env_transform"` // upper, lower, none
}

// DefaultConfigurationOptions returns default configuration options
func DefaultConfigurationOptions() *ConfigurationOptions {
	return &ConfigurationOptions{
		Sources:          []string{"file", "environment", "defaults"},
		IgnoreMissing:    false,
		IgnoreInvalid:    false,
		MergeStrategy:    "override",
		StrictValidation: true,
		AllowUnknown:     false,
		ValidateSchema:   true,
		Format:           "yaml",
		Indent:           2,
		SortKeys:         true,
		IncludeEmpty:     false,
		EnvPrefix:        "GENERATOR",
		EnvSeparator:     "_",
		EnvTransform:     "upper",
	}
}
