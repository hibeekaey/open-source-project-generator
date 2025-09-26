package interfaces

import "github.com/cuesoftinc/open-source-project-generator/pkg/models"

// ValidationEngine defines the contract for comprehensive project validation operations.
//
// The ValidationEngine interface provides essential validation capabilities for generated projects,
// templates, and configurations. It ensures that generated projects meet quality standards
// and are free from common configuration issues.
//
// The validation engine covers:
//   - Project structure and dependency validation
//   - Configuration validation with schema checking
//   - Template validation and structure checking
//   - Auto-fix capabilities for common issues
//   - Validation rule management and customization
//
// Implementations should provide:
//   - Comprehensive validation results with actionable feedback
//   - Validation rules for core functionality and best practices
//   - Auto-fix capabilities where possible
type ValidationEngine interface {
	// Basic validation methods
	ValidateProject(path string) (*models.ValidationResult, error)
	ValidatePackageJSON(path string) error
	ValidateGoMod(path string) error
	ValidateDockerfile(path string) error
	ValidateYAML(path string) error
	ValidateJSON(path string) error
	ValidateTemplate(path string) error

	// Comprehensive project validation
	ValidateProjectStructure(path string) (*StructureValidationResult, error)
	ValidateProjectDependencies(path string) (*DependencyValidationResult, error)
	ValidateProjectSecurity(path string) (*SecurityValidationResult, error)
	ValidateProjectQuality(path string) (*QualityValidationResult, error)

	// Configuration validation
	ValidateConfiguration(config *models.ProjectConfig) (*ConfigValidationResult, error)
	ValidateConfigurationSchema(config any, schema *ConfigSchema) error
	ValidateConfigurationValues(config *models.ProjectConfig) (*ConfigValidationResult, error)

	// Template validation (comprehensive versions)
	ValidateTemplateAdvanced(path string) (*TemplateValidationResult, error)
	ValidateTemplateMetadata(metadata *TemplateMetadata) error
	ValidateTemplateStructure(path string) (*StructureValidationResult, error)
	ValidateTemplateVariables(variables map[string]TemplateVariable) error

	// Validation options
	SetValidationRules(rules []ValidationRule) error
	GetValidationRules() []ValidationRule
	AddValidationRule(rule ValidationRule) error
	RemoveValidationRule(ruleID string) error

	// Auto-fix capabilities
	FixValidationIssues(path string, issues []ValidationIssue) (*FixResult, error)
	GetFixableIssues(issues []ValidationIssue) []ValidationIssue
	PreviewFixes(path string, issues []ValidationIssue) (*FixPreview, error)
	ApplyFix(path string, fix Fix) error

	// Validation reporting
	GenerateValidationReport(result *ValidationResult, format string) ([]byte, error)
	GetValidationSummary(results []*ValidationResult) (*ValidationSummary, error)
}

// Comprehensive validation types and structures

// ValidationRule defines a validation rule with its configuration
type ValidationRule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Severity    string         `json:"severity"`
	Enabled     bool           `json:"enabled"`
	Config      map[string]any `json:"config,omitempty"`
	Pattern     string         `json:"pattern,omitempty"`
	FileTypes   []string       `json:"file_types,omitempty"`
	Fixable     bool           `json:"fixable"`
}

// StructureValidationResult contains project structure validation results
type StructureValidationResult struct {
	Valid            bool                       `json:"valid"`
	RequiredFiles    []FileValidationResult     `json:"required_files"`
	RequiredDirs     []DirValidationResult      `json:"required_dirs"`
	NamingIssues     []NamingValidationIssue    `json:"naming_issues"`
	PermissionIssues []PermissionIssue          `json:"permission_issues"`
	Summary          StructureValidationSummary `json:"summary"`
}

// FileValidationResult represents validation result for a file
type FileValidationResult struct {
	Path     string            `json:"path"`
	Required bool              `json:"required"`
	Exists   bool              `json:"exists"`
	Valid    bool              `json:"valid"`
	Issues   []ValidationIssue `json:"issues"`
	Size     int64             `json:"size"`
	Mode     string            `json:"mode"`
}

// DirValidationResult represents validation result for a directory
type DirValidationResult struct {
	Path      string            `json:"path"`
	Required  bool              `json:"required"`
	Exists    bool              `json:"exists"`
	Valid     bool              `json:"valid"`
	Issues    []ValidationIssue `json:"issues"`
	FileCount int               `json:"file_count"`
	Mode      string            `json:"mode"`
}

// NamingValidationIssue represents a naming convention issue
type NamingValidationIssue struct {
	Path       string `json:"path"`
	Type       string `json:"type"` // file, directory, variable, function
	Current    string `json:"current"`
	Expected   string `json:"expected"`
	Convention string `json:"convention"`
	Severity   string `json:"severity"`
	Fixable    bool   `json:"fixable"`
}

// PermissionIssue represents a file permission issue
type PermissionIssue struct {
	Path     string `json:"path"`
	Current  string `json:"current"`
	Expected string `json:"expected"`
	Type     string `json:"type"` // file, directory
	Security bool   `json:"security"`
	Severity string `json:"severity"`
	Fixable  bool   `json:"fixable"`
}

// StructureValidationSummary contains structure validation statistics
type StructureValidationSummary struct {
	TotalFiles       int `json:"total_files"`
	ValidFiles       int `json:"valid_files"`
	TotalDirectories int `json:"total_directories"`
	ValidDirectories int `json:"valid_directories"`
	NamingIssues     int `json:"naming_issues"`
	PermissionIssues int `json:"permission_issues"`
}

// DependencyValidationResult contains dependency validation results
type DependencyValidationResult struct {
	Valid           bool                        `json:"valid"`
	Dependencies    []DependencyValidation      `json:"dependencies"`
	Vulnerabilities []DependencyVulnerability   `json:"vulnerabilities"`
	Outdated        []OutdatedDependency        `json:"outdated"`
	Conflicts       []DependencyConflict        `json:"conflicts"`
	Summary         DependencyValidationSummary `json:"summary"`
}

// DependencyValidation represents validation result for a dependency
type DependencyValidation struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Type           string `json:"type"` // direct, transitive
	Valid          bool   `json:"valid"`
	Available      bool   `json:"available"`
	LatestVersion  string `json:"latest_version"`
	SecurityIssues int    `json:"security_issues"`
	LicenseIssues  int    `json:"license_issues"`
}

// DependencyConflict represents a dependency conflict
type DependencyConflict struct {
	Dependency1 string `json:"dependency1"`
	Version1    string `json:"version1"`
	Dependency2 string `json:"dependency2"`
	Version2    string `json:"version2"`
	Reason      string `json:"reason"`
	Severity    string `json:"severity"`
}

// DependencyValidationSummary contains dependency validation statistics
type DependencyValidationSummary struct {
	TotalDependencies int `json:"total_dependencies"`
	ValidDependencies int `json:"valid_dependencies"`
	Vulnerabilities   int `json:"vulnerabilities"`
	OutdatedCount     int `json:"outdated_count"`
	ConflictCount     int `json:"conflict_count"`
}

// SecurityValidationResult contains security validation results
type SecurityValidationResult struct {
	Valid          bool                      `json:"valid"`
	SecurityIssues []SecurityIssue           `json:"security_issues"`
	Secrets        []SecretDetection         `json:"secrets"`
	Permissions    []PermissionIssue         `json:"permissions"`
	Configurations []SecurityConfig          `json:"configurations"`
	Summary        SecurityValidationSummary `json:"summary"`
}

// SecurityIssue represents a security issue
type SecurityIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Rule        string `json:"rule"`
	CWE         string `json:"cwe,omitempty"`
	Fixable     bool   `json:"fixable"`
}

// SecurityConfig represents security configuration validation
type SecurityConfig struct {
	Component   string `json:"component"`
	Setting     string `json:"setting"`
	Current     string `json:"current"`
	Recommended string `json:"recommended"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Fixable     bool   `json:"fixable"`
}

// SecurityValidationSummary contains security validation statistics
type SecurityValidationSummary struct {
	TotalIssues    int `json:"total_issues"`
	HighSeverity   int `json:"high_severity"`
	MediumSeverity int `json:"medium_severity"`
	LowSeverity    int `json:"low_severity"`
	SecretsFound   int `json:"secrets_found"`
	ConfigIssues   int `json:"config_issues"`
}

// QualityValidationResult contains code quality validation results
type QualityValidationResult struct {
	Valid       bool                     `json:"valid"`
	CodeSmells  []CodeSmell              `json:"code_smells"`
	Complexity  []ComplexityIssue        `json:"complexity"`
	Duplication []Duplication            `json:"duplication"`
	Coverage    *CoverageInfo            `json:"coverage,omitempty"`
	Summary     QualityValidationSummary `json:"summary"`
}

// ComplexityIssue represents a code complexity issue
type ComplexityIssue struct {
	Type       string `json:"type"` // cyclomatic, cognitive, npath
	File       string `json:"file"`
	Function   string `json:"function"`
	Line       int    `json:"line"`
	Complexity int    `json:"complexity"`
	Threshold  int    `json:"threshold"`
	Severity   string `json:"severity"`
}

// CoverageInfo represents test coverage information
type CoverageInfo struct {
	LinesCovered     int     `json:"lines_covered"`
	LinesTotal       int     `json:"lines_total"`
	LineCoverage     float64 `json:"line_coverage"`
	BranchesCovered  int     `json:"branches_covered"`
	BranchesTotal    int     `json:"branches_total"`
	BranchCoverage   float64 `json:"branch_coverage"`
	FunctionsCovered int     `json:"functions_covered"`
	FunctionsTotal   int     `json:"functions_total"`
	FunctionCoverage float64 `json:"function_coverage"`
}

// QualityValidationSummary contains quality validation statistics
type QualityValidationSummary struct {
	TotalIssues       int     `json:"total_issues"`
	CodeSmells        int     `json:"code_smells"`
	ComplexityIssues  int     `json:"complexity_issues"`
	DuplicationIssues int     `json:"duplication_issues"`
	QualityScore      float64 `json:"quality_score"`
	Maintainability   string  `json:"maintainability"`
}

// ConfigValidationResult contains the result of configuration validation
type ConfigValidationResult struct {
	Valid    bool                    `json:"valid"`
	Errors   []ConfigValidationError `json:"errors"`
	Warnings []ConfigValidationError `json:"warnings"`
	Summary  ConfigValidationSummary `json:"summary"`
}

// ConfigValidationError represents a configuration validation error
type ConfigValidationError struct {
	Field      string `json:"field"`
	Value      string `json:"value"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
	Severity   string `json:"severity"`
	Rule       string `json:"rule"`
}

// ConfigValidationSummary contains validation statistics
type ConfigValidationSummary struct {
	TotalProperties int `json:"total_properties"`
	ValidProperties int `json:"valid_properties"`
	ErrorCount      int `json:"error_count"`
	WarningCount    int `json:"warning_count"`
	MissingRequired int `json:"missing_required"`
}

// ConfigSchema defines the structure and validation rules for configuration
type ConfigSchema struct {
	Properties  map[string]PropertySchema `json:"properties"`
	Required    []string                  `json:"required"`
	Version     string                    `json:"version"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
}

// PropertySchema defines validation rules for individual configuration properties
type PropertySchema struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
	MinLength   *int     `json:"min_length,omitempty"`
	MaxLength   *int     `json:"max_length,omitempty"`
	Minimum     *float64 `json:"minimum,omitempty"`
	Maximum     *float64 `json:"maximum,omitempty"`
	Required    bool     `json:"required"`
	Deprecated  bool     `json:"deprecated"`
	Examples    []string `json:"examples,omitempty"`
}

// TemplateVariable defines a template variable with validation rules
type TemplateVariable struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Default     any                 `json:"default,omitempty"`
	Required    bool                `json:"required"`
	Validation  *VariableValidation `json:"validation,omitempty"`
	Examples    []string            `json:"examples,omitempty"`
	Deprecated  bool                `json:"deprecated"`
}

// VariableValidation defines validation rules for template variables
type VariableValidation struct {
	Pattern   string   `json:"pattern,omitempty"`
	MinLength *int     `json:"min_length,omitempty"`
	MaxLength *int     `json:"max_length,omitempty"`
	Minimum   *float64 `json:"minimum,omitempty"`
	Maximum   *float64 `json:"maximum,omitempty"`
	Enum      []string `json:"enum,omitempty"`
}

// FixResult contains the result of applying fixes
type FixResult struct {
	Applied []Fix        `json:"applied"`
	Failed  []FixFailure `json:"failed"`
	Skipped []Fix        `json:"skipped"`
	Summary FixSummary   `json:"summary"`
}

// Fix represents a fix that can be applied
type Fix struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	File        string         `json:"file"`
	Line        int            `json:"line"`
	Column      int            `json:"column"`
	Action      string         `json:"action"` // replace, insert, delete, rename
	Content     string         `json:"content"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Automatic   bool           `json:"automatic"`
}

// FixFailure represents a failed fix attempt
type FixFailure struct {
	Fix   Fix    `json:"fix"`
	Error string `json:"error"`
}

// FixPreview contains a preview of fixes to be applied
type FixPreview struct {
	Fixes   []Fix        `json:"fixes"`
	Changes []FileChange `json:"changes"`
	Summary FixSummary   `json:"summary"`
}

// FileChange represents a change to be made to a file
type FileChange struct {
	File        string `json:"file"`
	Action      string `json:"action"`
	LinesBefore int    `json:"lines_before"`
	LinesAfter  int    `json:"lines_after"`
	Preview     string `json:"preview"`
}

// FixSummary contains fix statistics
type FixSummary struct {
	TotalFixes    int `json:"total_fixes"`
	AppliedFixes  int `json:"applied_fixes"`
	FailedFixes   int `json:"failed_fixes"`
	SkippedFixes  int `json:"skipped_fixes"`
	FilesModified int `json:"files_modified"`
}

// DuplicationFile represents a file involved in duplication
type DuplicationFile struct {
	Path      string `json:"path"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

// ValidationRuleCategory defines categories for validation rules
const (
	ValidationCategoryStructure     = "structure"
	ValidationCategoryDependencies  = "dependencies"
	ValidationCategorySecurity      = "security"
	ValidationCategoryQuality       = "quality"
	ValidationCategoryConfiguration = "configuration"
	ValidationCategoryTemplate      = "template"
	ValidationCategoryNaming        = "naming"
	ValidationCategoryPermissions   = "permissions"
)

// ValidationSeverity defines severity levels for validation issues
const (
	ValidationSeverityError   = "error"
	ValidationSeverityWarning = "warning"
	ValidationSeverityInfo    = "info"
)

// FixAction defines types of fix actions
const (
	FixActionReplace = "replace"
	FixActionInsert  = "insert"
	FixActionDelete  = "delete"
	FixActionRename  = "rename"
	FixActionMove    = "move"
	FixActionCreate  = "create"
)
