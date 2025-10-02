// Package interfaces defines the core contracts and interfaces for the
// Open Source Project Generator components.
//
// This package contains interface definitions that enable dependency injection,
// testing, and modular architecture throughout the application.
package interfaces

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// CLIInterface defines the contract for command-line interface operations.
//
// This interface abstracts CLI functionality to enable testing and different
// CLI implementations. It covers comprehensive user interaction workflow for
// project configuration collection, generation, validation, and management.
//
// Implementations should provide:
//   - Interactive project configuration
//   - Template-based code generation
//   - Project validation and auditing
//   - Configuration management
//   - Offline mode support
//   - Comprehensive documentation generation
type CLIInterface interface {
	// Core operations
	Run(args []string) error

	// Interactive operations
	PromptProjectDetails() (*models.ProjectConfig, error)
	ConfirmGeneration(*models.ProjectConfig) bool

	// Advanced interactive operations
	PromptAdvancedOptions() (*AdvancedOptions, error)
	ConfirmAdvancedGeneration(*models.ProjectConfig, *AdvancedOptions) bool
	SelectTemplateInteractively(filter TemplateFilter) (*TemplateInfo, error)

	// Non-interactive operations
	GenerateFromConfig(path string, options GenerateOptions) error
	ValidateProject(path string, options ValidationOptions) (*ValidationResult, error)
	AuditProject(path string, options AuditOptions) (*AuditResult, error)

	// Advanced non-interactive operations
	GenerateWithAdvancedOptions(config *models.ProjectConfig, options *AdvancedOptions) error
	ValidateProjectAdvanced(path string, options *ValidationOptions) (*ValidationResult, error)
	AuditProjectAdvanced(path string, options *AuditOptions) (*AuditResult, error)

	// Template operations
	ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
	GetTemplateInfo(name string) (*TemplateInfo, error)
	ValidateTemplate(path string) (*TemplateValidationResult, error)

	// Template management operations
	SearchTemplates(query string) ([]TemplateInfo, error)
	GetTemplateMetadata(name string) (*TemplateMetadata, error)
	GetTemplateDependencies(name string) ([]string, error)
	ValidateCustomTemplate(path string) (*TemplateValidationResult, error)

	// Configuration operations
	ShowConfig() error
	SetConfig(key, value string) error
	EditConfig() error
	ValidateConfig() error
	ExportConfig(path string) error

	// Configuration management operations
	LoadConfiguration(sources []string) (*models.ProjectConfig, error)
	MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error)
	ValidateConfigurationSchema(config *models.ProjectConfig) error
	GetConfigurationSources() ([]ConfigSource, error)

	// Version and update operations
	ShowVersion(options VersionOptions) error
	CheckUpdates() (*UpdateInfo, error)
	InstallUpdates() error

	// Advanced version operations
	GetPackageVersions() (map[string]string, error)
	GetLatestPackageVersions() (map[string]string, error)
	CheckCompatibility(path string) (*CompatibilityResult, error)

	// Cache operations
	ShowCache() error
	ClearCache() error
	CleanCache() error

	// Cache management operations
	GetCacheStats() (*CacheStats, error)
	ValidateCache() error
	RepairCache() error
	EnableOfflineMode() error
	DisableOfflineMode() error

	// Utility operations
	ShowLogs() error

	// Logging and debugging operations
	SetLogLevel(level string) error
	GetLogLevel() string
	ShowRecentLogs(lines int, level string) error
	GetLogFileLocations() ([]string, error)

	// Automation and integration operations
	RunNonInteractive(config *models.ProjectConfig, options *AdvancedOptions) error
	GenerateReport(reportType string, format string, outputFile string) error
	GetExitCode() int
	SetExitCode(code int)

	// Component access operations
	GetVersionManager() VersionManager
	GetBuildInfo() (version, gitCommit, buildTime string)
}

// GenerateOptions defines options for project generation
type GenerateOptions struct {
	Force           bool     `json:"force" yaml:"force"`
	Minimal         bool     `json:"minimal" yaml:"minimal"`
	Offline         bool     `json:"offline" yaml:"offline"`
	UpdateVersions  bool     `json:"update_versions" yaml:"update_versions"`
	SkipValidation  bool     `json:"skip_validation" yaml:"skip_validation"`
	BackupExisting  bool     `json:"backup_existing" yaml:"backup_existing"`
	IncludeExamples bool     `json:"include_examples" yaml:"include_examples"`
	Templates       []string `json:"templates" yaml:"templates"`
	OutputPath      string   `json:"output_path" yaml:"output_path"`
	DryRun          bool     `json:"dry_run" yaml:"dry_run"`
	NonInteractive  bool     `json:"non_interactive" yaml:"non_interactive"`
}

// ValidationOptions defines options for project validation
type ValidationOptions struct {
	Verbose        bool     `json:"verbose"`
	Fix            bool     `json:"fix"`
	Report         bool     `json:"report"`
	ReportFormat   string   `json:"report_format"`
	Rules          []string `json:"rules"`
	IgnoreWarnings bool     `json:"ignore_warnings"`
	OutputFile     string   `json:"output_file"`
}

// VersionOptions defines options for version display
type VersionOptions struct {
	ShowPackages  bool   `json:"show_packages"`
	CheckUpdates  bool   `json:"check_updates"`
	ShowBuildInfo bool   `json:"show_build_info"`
	OutputFormat  string `json:"output_format"`
}

// TemplateFilter defines filtering options for template listing
type TemplateFilter struct {
	Category   string   `json:"category,omitempty"`
	Technology string   `json:"technology,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	MinVersion string   `json:"min_version,omitempty"`
	MaxVersion string   `json:"max_version,omitempty"`
}

// TemplateInfo contains information about a template
type TemplateInfo struct {
	Name         string           `json:"name"`
	DisplayName  string           `json:"display_name"`
	Description  string           `json:"description"`
	Category     string           `json:"category"`
	Technology   string           `json:"technology"`
	Version      string           `json:"version"`
	Tags         []string         `json:"tags"`
	Dependencies []string         `json:"dependencies"`
	Metadata     TemplateMetadata `json:"metadata"`
}

// TemplateMetadata contains metadata about a template
type TemplateMetadata struct {
	Author      string            `json:"author"`
	License     string            `json:"license"`
	Repository  string            `json:"repository"`
	Homepage    string            `json:"homepage"`
	Keywords    []string          `json:"keywords"`
	Maintainers []string          `json:"maintainers"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
	Variables   map[string]string `json:"variables"`
}

// ValidationResult contains the result of project validation
type ValidationResult struct {
	Valid          bool              `json:"valid"`
	Issues         []ValidationIssue `json:"issues"`
	Warnings       []ValidationIssue `json:"warnings"`
	Summary        ValidationSummary `json:"summary"`
	FixSuggestions []FixSuggestion   `json:"fix_suggestions"`
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
	Type     string            `json:"type"`
	Severity string            `json:"severity"`
	Message  string            `json:"message"`
	File     string            `json:"file,omitempty"`
	Line     int               `json:"line,omitempty"`
	Column   int               `json:"column,omitempty"`
	Rule     string            `json:"rule"`
	Fixable  bool              `json:"fixable"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ValidationSummary contains validation statistics
type ValidationSummary struct {
	TotalFiles   int `json:"total_files"`
	ValidFiles   int `json:"valid_files"`
	ErrorCount   int `json:"error_count"`
	WarningCount int `json:"warning_count"`
	FixableCount int `json:"fixable_count"`
}

// FixSuggestion represents a suggested fix for a validation issue
type FixSuggestion struct {
	Issue       ValidationIssue `json:"issue"`
	Description string          `json:"description"`
	AutoFixable bool            `json:"auto_fixable"`
	Command     string          `json:"command,omitempty"`
}

// TemplateValidationResult contains the result of template validation
type TemplateValidationResult struct {
	Valid    bool              `json:"valid"`
	Issues   []ValidationIssue `json:"issues"`
	Warnings []ValidationIssue `json:"warnings"`
	Summary  ValidationSummary `json:"summary"`
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	CurrentVersion  string    `json:"current_version"`
	LatestVersion   string    `json:"latest_version"`
	UpdateAvailable bool      `json:"update_available"`
	ReleaseNotes    string    `json:"release_notes"`
	DownloadURL     string    `json:"download_url"`
	ReleaseDate     time.Time `json:"release_date"`
	Breaking        bool      `json:"breaking"`
	Security        bool      `json:"security"`
	Recommended     bool      `json:"recommended"`
	Size            int64     `json:"size"`
	Checksum        string    `json:"checksum"`
	SignatureURL    string    `json:"signature_url"`
}

// AdvancedOptions contains advanced options for project generation
type AdvancedOptions struct {
	GenerateOptions

	// Security options
	EnableSecurityScanning bool     `json:"enable_security_scanning" yaml:"enable_security_scanning"`
	SecurityPolicies       []string `json:"security_policies" yaml:"security_policies"`

	// Quality options
	EnableQualityChecks bool     `json:"enable_quality_checks" yaml:"enable_quality_checks"`
	QualityRules        []string `json:"quality_rules" yaml:"quality_rules"`

	// Performance options
	EnablePerformanceOptimization bool `json:"enable_performance_optimization" yaml:"enable_performance_optimization"`
	BundleOptimization            bool `json:"bundle_optimization" yaml:"bundle_optimization"`

	// Documentation options
	GenerateDocumentation bool     `json:"generate_documentation" yaml:"generate_documentation"`
	DocumentationFormats  []string `json:"documentation_formats" yaml:"documentation_formats"`

	// CI/CD options
	EnableCICD     bool     `json:"enable_cicd" yaml:"enable_cicd"`
	CICDProviders  []string `json:"cicd_providers" yaml:"cicd_providers"`
	DeploymentType string   `json:"deployment_type" yaml:"deployment_type"`

	// Monitoring options
	EnableMonitoring  bool     `json:"enable_monitoring" yaml:"enable_monitoring"`
	MonitoringTools   []string `json:"monitoring_tools" yaml:"monitoring_tools"`
	LoggingFrameworks []string `json:"logging_frameworks" yaml:"logging_frameworks"`
}

// ConfigSource represents a configuration source
type ConfigSource struct {
	Type     string `json:"type"`     // file, environment, defaults
	Location string `json:"location"` // file path or environment variable name
	Priority int    `json:"priority"` // higher number = higher priority
	Valid    bool   `json:"valid"`    // whether the source is valid
}

// CompatibilityResult contains compatibility check results
type CompatibilityResult struct {
	Compatible       bool                    `json:"compatible"`
	GeneratorVersion string                  `json:"generator_version"`
	ProjectVersion   string                  `json:"project_version"`
	Issues           []CompatibilityIssue    `json:"issues"`
	Recommendations  []string                `json:"recommendations"`
	PackageVersions  map[string]VersionCheck `json:"package_versions"`
}

// CompatibilityIssue represents a compatibility issue
type CompatibilityIssue struct {
	Type        string `json:"type"`        // version, dependency, configuration
	Severity    string `json:"severity"`    // error, warning, info
	Component   string `json:"component"`   // component name
	Current     string `json:"current"`     // current version
	Required    string `json:"required"`    // required version
	Description string `json:"description"` // issue description
	Fixable     bool   `json:"fixable"`     // whether automatically fixable
}

// VersionCheck represents a version compatibility check
type VersionCheck struct {
	Current    string `json:"current"`
	Latest     string `json:"latest"`
	Compatible bool   `json:"compatible"`
	UpdateType string `json:"update_type"` // major, minor, patch
}

// CLIError represents a CLI error with detailed information
type CLIError struct {
	Type        string         `json:"type"`
	Message     string         `json:"message"`
	Code        int            `json:"code"`
	Details     map[string]any `json:"details,omitempty"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Context     *ErrorContext  `json:"context,omitempty"`
}

// ErrorContext provides context information for errors
type ErrorContext struct {
	Command     string            `json:"command"`
	Arguments   []string          `json:"arguments"`
	Flags       map[string]string `json:"flags"`
	WorkingDir  string            `json:"working_dir"`
	Timestamp   time.Time         `json:"timestamp"`
	Environment map[string]string `json:"environment,omitempty"`
}

// Error types for categorization
const (
	ErrorTypeValidation    = "validation"
	ErrorTypeConfiguration = "configuration"
	ErrorTypeTemplate      = "template"
	ErrorTypeNetwork       = "network"
	ErrorTypeFileSystem    = "filesystem"
	ErrorTypePermission    = "permission"
	ErrorTypeCache         = "cache"
	ErrorTypeVersion       = "version"
	ErrorTypeAudit         = "audit"
	ErrorTypeGeneration    = "generation"
	ErrorTypeInternal      = "internal"
)

// Error codes for programmatic handling
const (
	ErrorCodeSuccess              = 0
	ErrorCodeGeneral              = 1
	ErrorCodeValidationFailed     = 2
	ErrorCodeConfigurationInvalid = 3
	ErrorCodeTemplateNotFound     = 4
	ErrorCodeNetworkError         = 5
	ErrorCodeFileSystemError      = 6
	ErrorCodePermissionDenied     = 7
	ErrorCodeCacheError           = 8
	ErrorCodeVersionError         = 9
	ErrorCodeAuditFailed          = 10
	ErrorCodeGenerationFailed     = 11
	ErrorCodeInternalError        = 99
)

// Implement error interface for CLIError
func (e *CLIError) Error() string {
	return e.Message
}

// NewCLIError creates a new CLI error with the specified type and message
func NewCLIError(errorType, message string, code int) *CLIError {
	return &CLIError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Details: make(map[string]any),
		Context: &ErrorContext{
			Timestamp: time.Now(),
		},
	}
}

// WithDetails adds details to a CLI error
func (e *CLIError) WithDetails(key string, value any) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// WithSuggestions adds suggestions to a CLI error
func (e *CLIError) WithSuggestions(suggestions ...string) *CLIError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithContext adds context to a CLI error
func (e *CLIError) WithContext(ctx *ErrorContext) *CLIError {
	e.Context = ctx
	return e
}
