package interfaces

import (
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateEngine defines the contract for low-level template processing operations.
//
// The TemplateEngine interface provides core template processing capabilities:
//   - Single template file processing with variable substitution
//   - Recursive directory processing for complete project generation
//   - Custom function registration for extended template functionality
//   - Template loading, parsing, and rendering
//
// Implementations should provide:
//   - Robust error handling and validation
//   - Security considerations for template processing
//   - Integration with version management systems
type TemplateEngine interface {
	// Core template processing
	ProcessTemplate(path string, config *models.ProjectConfig) ([]byte, error)
	ProcessDirectory(templatePath string, outputPath string, config *models.ProjectConfig) error

	// Template loading and rendering
	LoadTemplate(path string) (*template.Template, error)
	RenderTemplate(tmpl *template.Template, data any) ([]byte, error)
	RegisterFunctions(funcMap template.FuncMap)
}

// TemplateManager defines the interface for high-level template management operations.
//
// This interface provides comprehensive template management including:
//   - Template discovery and listing with filtering
//   - Template validation and structure checking
//   - Template metadata management and caching
//   - Custom template support and installation
//   - High-level template processing using TemplateEngine
//
// TemplateManager uses TemplateEngine for low-level processing while providing
// higher-level management capabilities like discovery, validation, and caching.
type TemplateManager interface {
	// Template processing (high-level, uses template names/IDs)
	ProcessTemplate(templateName string, config *models.ProjectConfig, outputPath string) error
	ProcessCustomTemplate(path string, config *models.ProjectConfig, outputPath string) error
	PreviewTemplate(templateName string, config *models.ProjectConfig) (*TemplatePreview, error)

	// Template discovery
	ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
	GetTemplateInfo(name string) (*TemplateInfo, error)
	SearchTemplates(query string) ([]TemplateInfo, error)
	GetTemplatesByCategory(category string) ([]TemplateInfo, error)
	GetTemplatesByTechnology(technology string) ([]TemplateInfo, error)

	// Template validation
	ValidateTemplate(path string) (*TemplateValidationResult, error)
	ValidateTemplateStructure(template *TemplateInfo) error
	ValidateTemplateMetadata(metadata *TemplateMetadata) error
	ValidateCustomTemplate(path string) (*TemplateValidationResult, error)

	// Enhanced template validation
	ValidateTemplateStructureAdvanced(template *TemplateInfo) (*TemplateValidationResult, error)
	ValidateTemplateDependencies(template *TemplateInfo) (*TemplateValidationResult, error)
	ValidateTemplateCompatibility(template *TemplateInfo) (*TemplateValidationResult, error)
	ValidateTemplateMetadataAdvanced(metadata *TemplateMetadata) (*TemplateValidationResult, error)
	ValidateTemplateComprehensive(templateName string) (*TemplateValidationResult, error)

	// Template metadata
	GetTemplateMetadata(name string) (*TemplateMetadata, error)
	GetTemplateDependencies(name string) ([]string, error)
	GetTemplateCompatibility(name string) (*CompatibilityInfo, error)
	GetTemplateVariables(name string) (map[string]TemplateVariable, error)

	// Template management
	InstallTemplate(source string, name string) error
	UninstallTemplate(name string) error
	UpdateTemplate(name string) error
	GetTemplateLocation(name string) (string, error)

	// Template caching
	CacheTemplate(name string) error
	GetCachedTemplates() ([]TemplateInfo, error)
	ClearTemplateCache() error
	RefreshTemplateCache() error
}

// Enhanced template types and structures

// TemplatePreview contains a preview of template processing results
type TemplatePreview struct {
	TemplateName string                 `json:"template_name"`
	OutputPath   string                 `json:"output_path"`
	Files        []TemplatePreviewFile  `json:"files"`
	Variables    map[string]any         `json:"variables"`
	Summary      TemplatePreviewSummary `json:"summary"`
}

// TemplatePreviewFile represents a file in the template preview
type TemplatePreviewFile struct {
	Path        string `json:"path"`
	Type        string `json:"type"` // file, directory, symlink
	Size        int64  `json:"size"`
	Templated   bool   `json:"templated"`
	Executable  bool   `json:"executable"`
	Content     string `json:"content,omitempty"` // for small files
	ContentHash string `json:"content_hash"`
}

// TemplatePreviewSummary contains summary information about the template preview
type TemplatePreviewSummary struct {
	TotalFiles       int   `json:"total_files"`
	TotalDirectories int   `json:"total_directories"`
	TotalSize        int64 `json:"total_size"`
	TemplatedFiles   int   `json:"templated_files"`
	ExecutableFiles  int   `json:"executable_files"`
}

// CompatibilityInfo contains template compatibility information
type CompatibilityInfo struct {
	MinGeneratorVersion string               `json:"min_generator_version"`
	MaxGeneratorVersion string               `json:"max_generator_version"`
	SupportedPlatforms  []string             `json:"supported_platforms"`
	RequiredFeatures    []string             `json:"required_features"`
	Dependencies        []TemplateDependency `json:"dependencies"`
	Conflicts           []string             `json:"conflicts"`
}

// TemplateDependency represents a template dependency
type TemplateDependency struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Type        string `json:"type"` // template, package, tool
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// TemplateProcessingOptions defines options for template processing
type TemplateProcessingOptions struct {
	// Processing options
	DryRun         bool           `json:"dry_run"`
	Force          bool           `json:"force"`
	BackupExisting bool           `json:"backup_existing"`
	Variables      map[string]any `json:"variables"`

	// Output options
	OutputPath     string `json:"output_path"`
	CreateDirs     bool   `json:"create_dirs"`
	PreservePerms  bool   `json:"preserve_perms"`
	FollowSymlinks bool   `json:"follow_symlinks"`

	// Filtering options
	IncludePatterns []string `json:"include_patterns"`
	ExcludePatterns []string `json:"exclude_patterns"`
	IgnoreFiles     []string `json:"ignore_files"`

	// Validation options
	ValidateOutput bool `json:"validate_output"`
	StrictMode     bool `json:"strict_mode"`
}

// TemplateInstallOptions defines options for template installation
type TemplateInstallOptions struct {
	Source       string `json:"source"`       // git, http, file, registry
	Version      string `json:"version"`      // specific version or branch
	Force        bool   `json:"force"`        // force reinstall
	Verify       bool   `json:"verify"`       // verify signature/checksum
	Cache        bool   `json:"cache"`        // cache for offline use
	Dependencies bool   `json:"dependencies"` // install dependencies
}

// DefaultTemplateProcessingOptions returns default template processing options
func DefaultTemplateProcessingOptions() *TemplateProcessingOptions {
	return &TemplateProcessingOptions{
		DryRun:         false,
		Force:          false,
		BackupExisting: true,
		CreateDirs:     true,
		PreservePerms:  true,
		FollowSymlinks: false,
		ValidateOutput: true,
		StrictMode:     false,
	}
}

// DefaultTemplateInstallOptions returns default template installation options
func DefaultTemplateInstallOptions() *TemplateInstallOptions {
	return &TemplateInstallOptions{
		Source:       "registry",
		Version:      "latest",
		Force:        false,
		Verify:       true,
		Cache:        true,
		Dependencies: true,
	}
}
