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
	SelectComponents() ([]string, error)
	ConfirmGeneration(*models.ProjectConfig) bool

	// Non-interactive operations
	GenerateFromConfig(configPath string, options GenerateOptions) error
	ValidateProject(path string, options ValidationOptions) (*ValidationResult, error)
	AuditProject(path string, options AuditOptions) (*AuditResult, error)

	// Template operations
	ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
	GetTemplateInfo(name string) (*TemplateInfo, error)
	ValidateTemplate(path string) (*TemplateValidationResult, error)

	// Configuration operations
	ShowConfig() error
	SetConfig(key, value string) error
	EditConfig() error
	ValidateConfig() error
	ExportConfig(path string) error

	// Version and update operations
	ShowVersion(options VersionOptions) error
	CheckUpdates() (*UpdateInfo, error)
	InstallUpdates() error

	// Cache operations
	ShowCache() error
	ClearCache() error
	CleanCache() error

	// Utility operations
	ShowLogs() error
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

// AuditOptions defines options for project auditing
type AuditOptions struct {
	Security     bool   `json:"security"`
	Quality      bool   `json:"quality"`
	Licenses     bool   `json:"licenses"`
	Performance  bool   `json:"performance"`
	OutputFormat string `json:"output_format"`
	OutputFile   string `json:"output_file"`
	Detailed     bool   `json:"detailed"`
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
	Valid   bool              `json:"valid"`
	Issues  []ValidationIssue `json:"issues"`
	Summary ValidationSummary `json:"summary"`
}

// AuditResult contains the result of project auditing
type AuditResult struct {
	ProjectPath     string                  `json:"project_path"`
	AuditTime       time.Time               `json:"audit_time"`
	Security        *SecurityAuditResult    `json:"security,omitempty"`
	Quality         *QualityAuditResult     `json:"quality,omitempty"`
	Licenses        *LicenseAuditResult     `json:"licenses,omitempty"`
	Performance     *PerformanceAuditResult `json:"performance,omitempty"`
	OverallScore    float64                 `json:"overall_score"`
	Recommendations []string                `json:"recommendations"`
}

// SecurityAuditResult contains security audit results
type SecurityAuditResult struct {
	Score            float64           `json:"score"`
	Vulnerabilities  []Vulnerability   `json:"vulnerabilities"`
	PolicyViolations []PolicyViolation `json:"policy_violations"`
	Recommendations  []string          `json:"recommendations"`
}

// QualityAuditResult contains quality audit results
type QualityAuditResult struct {
	Score           float64       `json:"score"`
	CodeSmells      []CodeSmell   `json:"code_smells"`
	Duplications    []Duplication `json:"duplications"`
	TestCoverage    float64       `json:"test_coverage"`
	Recommendations []string      `json:"recommendations"`
}

// LicenseAuditResult contains license audit results
type LicenseAuditResult struct {
	Compatible      bool          `json:"compatible"`
	Licenses        []LicenseInfo `json:"licenses"`
	Conflicts       []LicenseInfo `json:"conflicts"`
	Recommendations []string      `json:"recommendations"`
}

// PerformanceAuditResult contains performance audit results
type PerformanceAuditResult struct {
	Score           float64            `json:"score"`
	BundleSize      int64              `json:"bundle_size"`
	LoadTime        time.Duration      `json:"load_time"`
	Issues          []PerformanceIssue `json:"issues"`
	Recommendations []string           `json:"recommendations"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Package     string `json:"package"`
	Version     string `json:"version"`
	FixedIn     string `json:"fixed_in,omitempty"`
}

// PolicyViolation represents a security policy violation
type PolicyViolation struct {
	Policy      string `json:"policy"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

// CodeSmell represents a code quality issue
type CodeSmell struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

// Duplication represents code duplication
type Duplication struct {
	Files      []string `json:"files"`
	Lines      int      `json:"lines"`
	Tokens     int      `json:"tokens"`
	Percentage float64  `json:"percentage"`
}

// LicenseInfo represents license information
type LicenseInfo struct {
	Name       string `json:"name"`
	SPDXID     string `json:"spdx_id"`
	Package    string `json:"package"`
	Compatible bool   `json:"compatible"`
}

// PerformanceIssue represents a performance issue
type PerformanceIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	File        string `json:"file,omitempty"`
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
}
