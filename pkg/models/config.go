package models

import (
	"time"
)

// ProjectConfig represents the basic project configuration
type ProjectConfig struct {
	// Basic project information
	Name         string `yaml:"name" json:"name"`
	Organization string `yaml:"organization" json:"organization"`
	Description  string `yaml:"description" json:"description"`
	License      string `yaml:"license" json:"license"`
	Author       string `yaml:"author" json:"author"`
	Email        string `yaml:"email" json:"email"`
	Repository   string `yaml:"repository" json:"repository"`

	// Component selection
	Components Components `yaml:"components" json:"components"`

	// Version configuration
	Versions *VersionConfig `yaml:"versions" json:"versions"`

	// Output configuration
	OutputPath string `yaml:"output_path" json:"output_path"`

	// Generation metadata
	GeneratedAt      time.Time `yaml:"generated_at" json:"generated_at"`
	GeneratorVersion string    `yaml:"generator_version" json:"generator_version"`
}

// Components defines which components to include in the generated project
type Components struct {
	Frontend       FrontendComponents       `yaml:"frontend" json:"frontend"`
	Backend        BackendComponents        `yaml:"backend" json:"backend"`
	Mobile         MobileComponents         `yaml:"mobile" json:"mobile"`
	Infrastructure InfrastructureComponents `yaml:"infrastructure" json:"infrastructure"`
}

// FrontendComponents defines frontend application options
type FrontendComponents struct {
	// New format (documented)
	MainApp bool `yaml:"main_app" json:"main_app"`
	Home    bool `yaml:"home" json:"home"`
	Admin   bool `yaml:"admin" json:"admin"`

	// Legacy format (for backward compatibility)
	NextJS NextJSComponents `yaml:"nextjs" json:"nextjs"`
}

// NextJSComponents defines Next.js specific options
type NextJSComponents struct {
	App    bool `yaml:"app" json:"app"`
	Home   bool `yaml:"home" json:"home"`
	Admin  bool `yaml:"admin" json:"admin"`
	Shared bool `yaml:"shared" json:"shared"`
}

// BackendComponents defines backend application options
type BackendComponents struct {
	// New format (documented)
	API bool `yaml:"api" json:"api"`

	// Legacy format (for backward compatibility)
	GoGin bool `yaml:"go_gin" json:"go_gin"`
}

// MobileComponents defines mobile application options
type MobileComponents struct {
	Android bool `yaml:"android" json:"android"`
	IOS     bool `yaml:"ios" json:"ios"`
}

// InfrastructureComponents defines infrastructure options
type InfrastructureComponents struct {
	Docker     bool `yaml:"docker" json:"docker"`
	Kubernetes bool `yaml:"kubernetes" json:"kubernetes"`
	Terraform  bool `yaml:"terraform" json:"terraform"`
}
