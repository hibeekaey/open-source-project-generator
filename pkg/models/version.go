package models

import (
	"time"
)

// VersionInfo represents detailed version information for a package or language
type VersionInfo struct {
	// Package identification
	Name     string `yaml:"name" json:"name" validate:"required"`
	Language string `yaml:"language" json:"language" validate:"required"`
	Type     string `yaml:"type" json:"type" validate:"required,oneof=language framework package"`

	// Version tracking
	CurrentVersion  string `yaml:"current_version" json:"current_version" validate:"required,semver"`
	LatestVersion   string `yaml:"latest_version" json:"latest_version" validate:"required,semver"`
	PreviousVersion string `yaml:"previous_version,omitempty" json:"previous_version,omitempty" validate:"omitempty,semver"`

	// Security information
	SecurityIssues []SecurityIssue `yaml:"security_issues,omitempty" json:"security_issues,omitempty"`
	IsSecure       bool            `yaml:"is_secure" json:"is_secure"`

	// Update metadata
	UpdatedAt    time.Time `yaml:"updated_at" json:"updated_at"`
	CheckedAt    time.Time `yaml:"checked_at" json:"checked_at"`
	UpdateSource string    `yaml:"update_source" json:"update_source" validate:"required"`

	// Registry information
	RegistryURL string            `yaml:"registry_url,omitempty" json:"registry_url,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	ID          string    `yaml:"id" json:"id" validate:"required"`
	Severity    string    `yaml:"severity" json:"severity" validate:"required,oneof=low medium high critical"`
	Description string    `yaml:"description" json:"description" validate:"required"`
	FixedIn     string    `yaml:"fixed_in,omitempty" json:"fixed_in,omitempty" validate:"omitempty,semver"`
	ReportedAt  time.Time `yaml:"reported_at" json:"reported_at"`
	URL         string    `yaml:"url,omitempty" json:"url,omitempty" validate:"omitempty,url"`
}

// VersionStore represents a collection of version information
type VersionStore struct {
	// Store metadata
	LastUpdated time.Time `yaml:"last_updated" json:"last_updated"`
	Version     string    `yaml:"version" json:"version" validate:"required"`

	// Language versions
	Languages map[string]*VersionInfo `yaml:"languages" json:"languages" validate:"required"`

	// Framework versions
	Frameworks map[string]*VersionInfo `yaml:"frameworks" json:"frameworks" validate:"required"`

	// Package versions
	Packages map[string]*VersionInfo `yaml:"packages" json:"packages" validate:"required"`

	// Update policy
	UpdatePolicy UpdatePolicy `yaml:"update_policy" json:"update_policy"`
}

// UpdatePolicy defines how versions should be updated
type UpdatePolicy struct {
	AutoUpdate             bool          `yaml:"auto_update" json:"auto_update"`
	SecurityPriority       bool          `yaml:"security_priority" json:"security_priority"`
	BreakingChangeApproval bool          `yaml:"breaking_change_approval" json:"breaking_change_approval"`
	UpdateSchedule         string        `yaml:"update_schedule" json:"update_schedule" validate:"required"`
	MaxAge                 time.Duration `yaml:"max_age" json:"max_age"`
}

// VersionUpdateRequest represents a request to update version information
type VersionUpdateRequest struct {
	Name           string            `json:"name" validate:"required"`
	TargetVersion  string            `json:"target_version" validate:"required,semver"`
	Force          bool              `json:"force"`
	SkipValidation bool              `json:"skip_validation"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// VersionUpdateResult represents the result of a version update operation
type VersionUpdateResult struct {
	Success         bool              `json:"success"`
	PreviousVersion string            `json:"previous_version"`
	NewVersion      string            `json:"new_version"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Errors          []string          `json:"errors,omitempty"`
	Warnings        []string          `json:"warnings,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// VersionQuery represents a query for version information
type VersionQuery struct {
	Name     string `json:"name,omitempty"`
	Language string `json:"language,omitempty"`
	Type     string `json:"type,omitempty" validate:"omitempty,oneof=language framework package"`
	Outdated bool   `json:"outdated,omitempty"`
	Insecure bool   `json:"insecure,omitempty"`
}

// VersionReport represents a summary report of version status
type VersionReport struct {
	GeneratedAt     time.Time                 `json:"generated_at"`
	TotalPackages   int                       `json:"total_packages"`
	OutdatedCount   int                       `json:"outdated_count"`
	SecurityIssues  int                       `json:"security_issues"`
	LastUpdateCheck time.Time                 `json:"last_update_check"`
	Summary         map[string]VersionSummary `json:"summary"`
	Details         map[string]*VersionInfo   `json:"details"`
	Recommendations []UpdateRecommendation    `json:"recommendations"`
}

// VersionSummary provides a high-level summary for a category
type VersionSummary struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	Outdated int `json:"outdated"`
	Insecure int `json:"insecure"`
}

// UpdateRecommendation suggests version updates
type UpdateRecommendation struct {
	Name               string `json:"name"`
	CurrentVersion     string `json:"current_version"`
	RecommendedVersion string `json:"recommended_version"`
	Priority           string `json:"priority" validate:"oneof=low medium high critical"`
	Reason             string `json:"reason"`
	BreakingChange     bool   `json:"breaking_change"`
}

// DashboardData represents comprehensive dashboard information
type DashboardData struct {
	GeneratedAt       time.Time                  `json:"generated_at"`
	VersionReport     *VersionReport             `json:"version_report"`
	TemplateStatus    *TemplateConsistencyStatus `json:"template_status"`
	ValidationResults *ValidationResults         `json:"validation_results"`
	Metadata          map[string]string          `json:"metadata"`
}

// TemplateConsistencyStatus represents template consistency check results
type TemplateConsistencyStatus struct {
	CheckedAt time.Time          `json:"checked_at"`
	Status    string             `json:"status" validate:"oneof=consistent issues_found critical unknown"`
	Issues    []ConsistencyIssue `json:"issues"`
	Summary   ConsistencySummary `json:"summary"`
}

// ConsistencyIssue represents a template consistency issue
type ConsistencyIssue struct {
	Type        string `json:"type" validate:"required"`
	Severity    string `json:"severity" validate:"required,oneof=low medium high critical"`
	Template    string `json:"template" validate:"required"`
	Description string `json:"description" validate:"required"`
	File        string `json:"file,omitempty"`
}

// ConsistencySummary provides summary statistics for template consistency
type ConsistencySummary struct {
	TotalTemplates      int `json:"total_templates"`
	ConsistentTemplates int `json:"consistent_templates"`
	IssuesFound         int `json:"issues_found"`
}

// ValidationResults represents comprehensive validation results
type ValidationResults struct {
	CheckedAt time.Time                           `json:"checked_at"`
	Status    string                              `json:"status" validate:"oneof=passed warnings failed unknown"`
	Results   map[string]DetailedValidationResult `json:"results"`
	Summary   ValidationSummary                   `json:"summary"`
}

// DetailedValidationResult represents a single validation check result with detailed information
type DetailedValidationResult struct {
	Name      string            `json:"name" validate:"required"`
	Status    string            `json:"status" validate:"required,oneof=passed warnings failed"`
	CheckedAt time.Time         `json:"checked_at"`
	Errors    []string          `json:"errors"`
	Warnings  []string          `json:"warnings"`
	Details   map[string]string `json:"details"`
}

// ValidationSummary provides summary statistics for validation results
type ValidationSummary struct {
	TotalChecks  int `json:"total_checks"`
	PassedChecks int `json:"passed_checks"`
	FailedChecks int `json:"failed_checks"`
	Warnings     int `json:"warnings"`
}

// UpdateReport represents a comprehensive update report
type UpdateReport struct {
	GeneratedAt     time.Time              `json:"generated_at"`
	ReportID        string                 `json:"report_id"`
	Type            string                 `json:"type" validate:"oneof=version_update template_update"`
	Summary         string                 `json:"summary"`
	Details         string                 `json:"details,omitempty"`
	Recommendations []UpdateRecommendation `json:"recommendations,omitempty"`
	TemplateUpdates []TemplateUpdate       `json:"template_updates,omitempty"`
	Metadata        map[string]string      `json:"metadata"`
}

// TemplateUpdate represents the result of updating a template
type TemplateUpdate struct {
	TemplatePath   string            `json:"template_path"`
	Success        bool              `json:"success"`
	UpdatedAt      time.Time         `json:"updated_at"`
	VersionChanges map[string]string `json:"version_changes"`
	Error          string            `json:"error,omitempty"`
	BackupPath     string            `json:"backup_path,omitempty"`
}

// SecurityReport represents a security-focused report
type SecurityReport struct {
	GeneratedAt    time.Time             `json:"generated_at"`
	ReportID       string                `json:"report_id"`
	TotalIssues    int                   `json:"total_issues"`
	CriticalIssues int                   `json:"critical_issues"`
	HighIssues     int                   `json:"high_issues"`
	Issues         []SecurityIssueDetail `json:"issues"`
	Metadata       map[string]string     `json:"metadata"`
}

// SecurityIssueDetail provides detailed information about a security issue
type SecurityIssueDetail struct {
	PackageName    string        `json:"package_name"`
	CurrentVersion string        `json:"current_version"`
	SecurityIssue  SecurityIssue `json:"security_issue"`
	RecommendedFix string        `json:"recommended_fix,omitempty"`
}

// ReportInfo provides metadata about generated reports
type ReportInfo struct {
	Filename    string    `json:"filename"`
	Type        string    `json:"type"`
	GeneratedAt time.Time `json:"generated_at"`
	Size        int64     `json:"size"`
}

// AuditSummary provides summary statistics for audit events
type AuditSummary struct {
	Period      string         `json:"period"`
	TotalEvents int            `json:"total_events"`
	EventTypes  map[string]int `json:"event_types"`
	Actions     map[string]int `json:"actions"`
	SuccessRate float64        `json:"success_rate"`
}

// NodeVersionConfig represents Node.js version configuration for templates
type NodeVersionConfig struct {
	// Runtime version requirement (e.g., ">=20.0.0")
	Runtime string `yaml:"runtime" json:"runtime" validate:"required"`

	// TypeScript Node.js types version (e.g., "^20.17.0")
	TypesPackage string `yaml:"types" json:"types" validate:"required"`

	// NPM version requirement (e.g., ">=10.0.0")
	NPMVersion string `yaml:"npm" json:"npm" validate:"required"`

	// Docker base image version (e.g., "node:20-alpine")
	DockerImage string `yaml:"docker" json:"docker" validate:"required"`

	// Additional metadata
	LTSStatus   bool              `yaml:"lts_status" json:"lts_status"`
	EOLDate     *time.Time        `yaml:"eol_date,omitempty" json:"eol_date,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// VersionCompatibilityMatrix defines compatibility rules between different version components
type VersionCompatibilityMatrix struct {
	NodeJS       NodeVersionConfig            `yaml:"nodejs" json:"nodejs"`
	Frameworks   map[string]FrameworkVersion  `yaml:"frameworks,omitempty" json:"frameworks,omitempty"`
	Dependencies map[string]DependencyVersion `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	UpdatedAt    time.Time                    `yaml:"updated_at" json:"updated_at"`
}

// FrameworkVersion represents version requirements for a specific framework
type FrameworkVersion struct {
	Name           string            `yaml:"name" json:"name" validate:"required"`
	Version        string            `yaml:"version" json:"version" validate:"required"`
	MinNodeVersion string            `yaml:"min_node_version" json:"min_node_version" validate:"required"`
	MaxNodeVersion string            `yaml:"max_node_version,omitempty" json:"max_node_version,omitempty"`
	TypesPackage   string            `yaml:"types_package,omitempty" json:"types_package,omitempty"`
	Metadata       map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// DependencyVersion represents version requirements for a specific dependency
type DependencyVersion struct {
	Name           string            `yaml:"name" json:"name" validate:"required"`
	Version        string            `yaml:"version" json:"version" validate:"required"`
	MinNodeVersion string            `yaml:"min_node_version" json:"min_node_version" validate:"required"`
	MaxNodeVersion string            `yaml:"max_node_version,omitempty" json:"max_node_version,omitempty"`
	Metadata       map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// VersionValidationResult represents the result of version validation
type VersionValidationResult struct {
	Valid       bool                       `json:"valid"`
	Errors      []VersionValidationError   `json:"errors,omitempty"`
	Warnings    []VersionValidationWarning `json:"warnings,omitempty"`
	Suggestions []VersionSuggestion        `json:"suggestions,omitempty"`
	ValidatedAt time.Time                  `json:"validated_at"`
}

// VersionValidationError represents a version validation error
type VersionValidationError struct {
	Field    string `json:"field" validate:"required"`
	Value    string `json:"value"`
	Expected string `json:"expected,omitempty"`
	Message  string `json:"message" validate:"required"`
	Severity string `json:"severity" validate:"required,oneof=error critical"`
	Code     string `json:"code,omitempty"`
}

// VersionValidationWarning represents a version validation warning
type VersionValidationWarning struct {
	Field   string `json:"field" validate:"required"`
	Value   string `json:"value"`
	Message string `json:"message" validate:"required"`
	Code    string `json:"code,omitempty"`
}

// VersionSuggestion represents a suggested version fix
type VersionSuggestion struct {
	Field          string `json:"field" validate:"required"`
	CurrentValue   string `json:"current_value"`
	SuggestedValue string `json:"suggested_value" validate:"required"`
	Reason         string `json:"reason" validate:"required"`
	Priority       string `json:"priority" validate:"required,oneof=low medium high"`
	BreakingChange bool   `json:"breaking_change"`
}
