// Package interfaces defines the core contracts and interfaces for the
// Open Source Project Generator components.
package interfaces

// TemplateManager defines the contract for template management operations.
//
// This interface abstracts template discovery, validation, and processing
// to enable different template sources and implementations.
type TemplateManager interface {
	// Template discovery
	ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
	GetTemplateInfo(name string) (*TemplateInfo, error)
	SearchTemplates(query string) ([]TemplateInfo, error)

	// Template validation
	ValidateTemplate(path string) (*TemplateValidationResult, error)
	ValidateTemplateStructure(template *TemplateInfo) error

	// Template processing
	ProcessTemplate(templateName string, config interface{}, outputPath string) error
	ProcessCustomTemplate(templatePath string, config interface{}, outputPath string) error

	// Template metadata
	GetTemplateMetadata(name string) (*TemplateMetadata, error)
	GetTemplateDependencies(name string) ([]string, error)
	GetTemplateCompatibility(name string) (*CompatibilityInfo, error)
}

// CompatibilityInfo contains template compatibility information
type CompatibilityInfo struct {
	MinGeneratorVersion string            `json:"min_generator_version"`
	MaxGeneratorVersion string            `json:"max_generator_version"`
	SupportedPlatforms  []string          `json:"supported_platforms"`
	RequiredFeatures    []string          `json:"required_features"`
	Dependencies        map[string]string `json:"dependencies"`
}
