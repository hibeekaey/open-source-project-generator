package models

import (
	"time"
)

// ProjectConfig represents the complete project configuration
type ProjectConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	OutputDir   string            `yaml:"output_dir" json:"output_dir"`
	Components  []ComponentConfig `yaml:"components" json:"components"`
	Integration IntegrationConfig `yaml:"integration" json:"integration"`
	Options     ProjectOptions    `yaml:"options" json:"options"`
}

// ComponentConfig represents configuration for a single project component
type ComponentConfig struct {
	Type    string                 `yaml:"type" json:"type"` // "nextjs", "go-backend", "android", "ios"
	Name    string                 `yaml:"name" json:"name"`
	Config  map[string]interface{} `yaml:"config" json:"config"`
	Enabled bool                   `yaml:"enabled" json:"enabled"`
}

// IntegrationConfig defines how components should be integrated
type IntegrationConfig struct {
	GenerateDockerCompose bool              `yaml:"generate_docker_compose" json:"generate_docker_compose"`
	GenerateScripts       bool              `yaml:"generate_scripts" json:"generate_scripts"`
	APIEndpoints          map[string]string `yaml:"api_endpoints" json:"api_endpoints"`
	SharedEnvironment     map[string]string `yaml:"shared_environment" json:"shared_environment"`
}

// ProjectOptions defines generation options
type ProjectOptions struct {
	UseExternalTools bool `yaml:"use_external_tools" json:"use_external_tools"`
	DryRun           bool `yaml:"dry_run" json:"dry_run"`
	Verbose          bool `yaml:"verbose" json:"verbose"`
	CreateBackup     bool `yaml:"create_backup" json:"create_backup"`
	ForceOverwrite   bool `yaml:"force_overwrite" json:"force_overwrite"`
	DisableParallel  bool `yaml:"disable_parallel" json:"disable_parallel"`
	StreamOutput     bool `yaml:"stream_output" json:"stream_output"`
}

// BootstrapSpec defines the specification for bootstrap tool execution
type BootstrapSpec struct {
	ComponentType string                 // "nextjs", "go-backend", "android", "ios"
	TargetDir     string                 // Where to generate
	Config        map[string]interface{} // Component-specific config
	Flags         []string               // Additional CLI flags
	Timeout       time.Duration          // Execution timeout
}

// FallbackSpec defines the specification for fallback generation
type FallbackSpec struct {
	ComponentType string                 // Component type to generate
	TargetDir     string                 // Where to generate
	Config        map[string]interface{} // Component-specific config
	TemplatePath  string                 // Path to templates
}

// Component represents a generated project component
type Component struct {
	Type        string                 // Component type
	Name        string                 // Component name
	Path        string                 // Generated path
	Config      map[string]interface{} // Configuration used
	Metadata    map[string]string      // Additional metadata
	GeneratedAt time.Time              // Generation timestamp
}
