package models

import (
	"time"
)

// ProjectConfig represents the complete project configuration
type ProjectConfig struct {
	// Basic project information
	Name         string `yaml:"name" json:"name" validate:"required,min=1,max=50,alphanum"`
	Organization string `yaml:"organization" json:"organization" validate:"required,min=1,max=100"`
	Description  string `yaml:"description" json:"description" validate:"required,max=500"`
	License      string `yaml:"license" json:"license" validate:"required,oneof=MIT Apache-2.0 GPL-3.0 BSD-3-Clause"`
	Author       string `yaml:"author" json:"author" validate:"omitempty,max=100"`
	Email        string `yaml:"email" json:"email" validate:"omitempty,email"`
	Repository   string `yaml:"repository" json:"repository" validate:"omitempty,url"`

	// Component selection
	Components Components `yaml:"components" json:"components" validate:"required"`

	// Version configuration
	Versions *VersionConfig `yaml:"versions" json:"versions" validate:"required"`

	// Custom variables
	CustomVars map[string]string `yaml:"custom_vars" json:"custom_vars" validate:"omitempty"`

	// Output configuration
	OutputPath string `yaml:"output_path" json:"output_path" validate:"required,min=1"`

	// Generation metadata
	GeneratedAt      time.Time `yaml:"generated_at" json:"generated_at"`
	GeneratorVersion string    `yaml:"generator_version" json:"generator_version"`
}

// Components defines which components to include in the generated project
type Components struct {
	Frontend       FrontendComponents       `yaml:"frontend" json:"frontend" validate:"required"`
	Backend        BackendComponents        `yaml:"backend" json:"backend" validate:"required"`
	Mobile         MobileComponents         `yaml:"mobile" json:"mobile" validate:"required"`
	Infrastructure InfrastructureComponents `yaml:"infrastructure" json:"infrastructure" validate:"required"`
}

// FrontendComponents defines frontend application options
type FrontendComponents struct {
	MainApp bool `yaml:"main_app" json:"main_app"`
	Home    bool `yaml:"home" json:"home"`
	Admin   bool `yaml:"admin" json:"admin"`
}

// BackendComponents defines backend service options
type BackendComponents struct {
	API bool `yaml:"api" json:"api"`
}

// MobileComponents defines mobile application options
type MobileComponents struct {
	Android bool `yaml:"android" json:"android"`
	IOS     bool `yaml:"ios" json:"ios"`
}

// InfrastructureComponents defines infrastructure options
type InfrastructureComponents struct {
	Terraform  bool `yaml:"terraform" json:"terraform"`
	Kubernetes bool `yaml:"kubernetes" json:"kubernetes"`
	Docker     bool `yaml:"docker" json:"docker"`
}

// VersionConfig holds version information for packages and frameworks
type VersionConfig struct {
	Node      string            `yaml:"node" json:"node" validate:"required,semver"`
	Go        string            `yaml:"go" json:"go" validate:"required,semver"`
	Kotlin    string            `yaml:"kotlin" json:"kotlin" validate:"omitempty,semver"`
	Swift     string            `yaml:"swift" json:"swift" validate:"omitempty,semver"`
	NextJS    string            `yaml:"nextjs" json:"nextjs" validate:"omitempty,semver"`
	React     string            `yaml:"react" json:"react" validate:"omitempty,semver"`
	Packages  map[string]string `yaml:"packages" json:"packages" validate:"omitempty,dive,keys,required,endkeys,required,semver"`
	UpdatedAt time.Time         `yaml:"updated_at" json:"updated_at"`

	// Enhanced Node.js version configuration
	NodeJS *NodeVersionConfig `yaml:"nodejs,omitempty" json:"nodejs,omitempty"`
}

// ValidationResult represents the result of configuration validation
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ConfigDefaults holds default configuration values
type ConfigDefaults struct {
	License      string                 `yaml:"license" json:"license"`
	Components   Components             `yaml:"components" json:"components"`
	CustomVars   map[string]string      `yaml:"custom_vars" json:"custom_vars"`
	TemplateVars map[string]interface{} `yaml:"template_vars" json:"template_vars"`
}
