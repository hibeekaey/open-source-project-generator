// Package interfaces defines workflow interfaces for end-to-end operations.
//
// This file contains interfaces for comprehensive workflow management including
// project generation, validation, audit, and configuration workflows.
package interfaces

import (
	"context"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// WorkflowManager defines the interface for managing end-to-end workflows.
//
// This interface provides comprehensive workflow orchestration capabilities including:
//   - Project generation workflows with validation and customization
//   - Validation and audit integration workflows
//   - Configuration management workflows
//   - Offline mode workflows with cache integration
//   - Progress tracking and error recovery
type WorkflowManager interface {
	// Project generation workflows
	CreateProjectWorkflow(config *models.ProjectConfig, options *ProjectWorkflowOptions) (ProjectWorkflow, error)
	ExecuteProjectGeneration(ctx context.Context, config *models.ProjectConfig, options *ProjectWorkflowOptions) (*ProjectWorkflowResult, error)

	// Validation and audit workflows
	CreateValidationWorkflow(projectPath string, options *ValidationWorkflowOptions) (ValidationWorkflow, error)
	CreateAuditWorkflow(projectPath string, options *AuditWorkflowOptions) (AuditWorkflow, error)
	CreateValidationAuditWorkflow(projectPath string, options *ValidationAuditWorkflowOptions) (ValidationAuditWorkflow, error)

	// Configuration workflows
	CreateConfigurationWorkflow(options *ConfigurationWorkflowOptions) (ConfigurationWorkflow, error)
	ExportConfiguration(ctx context.Context, config *models.ProjectConfig, outputPath string) error
	ImportConfiguration(ctx context.Context, configPath string) (*models.ProjectConfig, error)
	ValidateConfiguration(ctx context.Context, config *models.ProjectConfig) (*ConfigValidationResult, error)
	MergeConfigurations(ctx context.Context, configs []*models.ProjectConfig) (*models.ProjectConfig, error)

	// Offline workflows
	CreateOfflineWorkflow(options *OfflineWorkflowOptions) (OfflineWorkflow, error)
	SyncCache(ctx context.Context) (*CacheStats, error)
	ValidateCache(ctx context.Context) error
	GenerateOffline(ctx context.Context, config *models.ProjectConfig, outputPath string) (*ProjectWorkflowResult, error)

	// Workflow management
	GetActiveWorkflows() []WorkflowInfo
	CancelWorkflow(workflowID string) error
	GetWorkflowStatus(workflowID string) (*WorkflowStatus, error)
	GetWorkflowHistory() []WorkflowInfo
}

// ProjectWorkflow defines the interface for project generation workflows.
type ProjectWorkflow interface {
	// Workflow execution
	Execute(ctx context.Context) (*ProjectWorkflowResult, error)
	Cancel() error

	// Progress tracking
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))

	// Configuration
	GetConfiguration() *models.ProjectConfig
	UpdateConfiguration(config *models.ProjectConfig) error
	GetOptions() *ProjectWorkflowOptions
	UpdateOptions(options *ProjectWorkflowOptions) error

	// Status
	GetStatus() WorkflowStatus
	GetID() string
}

// ValidationWorkflow defines the interface for validation workflows.
type ValidationWorkflow interface {
	Execute(ctx context.Context) (*ValidationWorkflowResult, error)
	Cancel() error
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))
	GetStatus() WorkflowStatus
	GetID() string
}

// AuditWorkflow defines the interface for audit workflows.
type AuditWorkflow interface {
	Execute(ctx context.Context) (*AuditWorkflowResult, error)
	Cancel() error
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))
	GetStatus() WorkflowStatus
	GetID() string
}

// ValidationAuditWorkflow defines the interface for combined validation and audit workflows.
type ValidationAuditWorkflow interface {
	Execute(ctx context.Context) (*ValidationAuditWorkflowResult, error)
	Cancel() error
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))
	GetStatus() WorkflowStatus
	GetID() string
}

// ConfigurationWorkflow defines the interface for configuration workflows.
type ConfigurationWorkflow interface {
	Execute(ctx context.Context) (*ConfigurationWorkflowResult, error)
	Cancel() error
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))
	GetStatus() WorkflowStatus
	GetID() string
}

// OfflineWorkflow defines the interface for offline workflows.
type OfflineWorkflow interface {
	Execute(ctx context.Context) (*OfflineWorkflowResult, error)
	Cancel() error
	GetProgress() *WorkflowProgress
	SetProgressCallback(callback func(*WorkflowProgress))
	GetStatus() WorkflowStatus
	GetID() string
}

// Workflow options structures

// ProjectWorkflowOptions defines options for project generation workflows.
type ProjectWorkflowOptions struct {
	OutputPath        string                  `json:"output_path"`
	DryRun            bool                    `json:"dry_run"`
	Offline           bool                    `json:"offline"`
	Force             bool                    `json:"force"`
	BackupExisting    bool                    `json:"backup_existing"`
	ValidateAfter     bool                    `json:"validate_after"`
	AuditAfter        bool                    `json:"audit_after"`
	GenerateReport    bool                    `json:"generate_report"`
	ValidationOptions *ValidationOptions      `json:"validation_options,omitempty"`
	AuditOptions      *AuditOptions           `json:"audit_options,omitempty"`
	ProgressCallback  func(*WorkflowProgress) `json:"-"`
	Timeout           time.Duration           `json:"timeout"`
	Parallel          bool                    `json:"parallel"`
}

// ValidationWorkflowOptions defines options for validation workflows.
type ValidationWorkflowOptions struct {
	ProjectPath       string                  `json:"project_path"`
	ValidationOptions *ValidationOptions      `json:"validation_options,omitempty"`
	FixIssues         bool                    `json:"fix_issues"`
	GenerateReport    bool                    `json:"generate_report"`
	OutputFormat      string                  `json:"output_format"`
	OutputFile        string                  `json:"output_file"`
	ProgressCallback  func(*WorkflowProgress) `json:"-"`
	Timeout           time.Duration           `json:"timeout"`
}

// AuditWorkflowOptions defines options for audit workflows.
type AuditWorkflowOptions struct {
	ProjectPath      string                  `json:"project_path"`
	AuditOptions     *AuditOptions           `json:"audit_options,omitempty"`
	GenerateReport   bool                    `json:"generate_report"`
	OutputFormat     string                  `json:"output_format"`
	OutputFile       string                  `json:"output_file"`
	ProgressCallback func(*WorkflowProgress) `json:"-"`
	Timeout          time.Duration           `json:"timeout"`
}

// ValidationAuditWorkflowOptions defines options for combined validation and audit workflows.
type ValidationAuditWorkflowOptions struct {
	ProjectPath       string                  `json:"project_path"`
	ValidationEnabled bool                    `json:"validation_enabled"`
	AuditEnabled      bool                    `json:"audit_enabled"`
	ValidationOptions *ValidationOptions      `json:"validation_options,omitempty"`
	AuditOptions      *AuditOptions           `json:"audit_options,omitempty"`
	FixIssues         bool                    `json:"fix_issues"`
	GenerateReport    bool                    `json:"generate_report"`
	OutputFormat      string                  `json:"output_format"`
	OutputFile        string                  `json:"output_file"`
	ProgressCallback  func(*WorkflowProgress) `json:"-"`
	Timeout           time.Duration           `json:"timeout"`
}

// ConfigurationWorkflowOptions defines options for configuration workflows.
type ConfigurationWorkflowOptions struct {
	Operation        string                  `json:"operation"` // export, import, validate, merge
	ConfigPath       string                  `json:"config_path"`
	OutputPath       string                  `json:"output_path"`
	Format           string                  `json:"format"`
	Validate         bool                    `json:"validate"`
	Merge            bool                    `json:"merge"`
	Sources          []string                `json:"sources,omitempty"`
	ProgressCallback func(*WorkflowProgress) `json:"-"`
	Timeout          time.Duration           `json:"timeout"`
}

// OfflineWorkflowOptions defines options for offline workflows.
type OfflineWorkflowOptions struct {
	Operation        string                  `json:"operation"` // sync, validate, generate
	CacheUpdate      bool                    `json:"cache_update"`
	ValidateCache    bool                    `json:"validate_cache"`
	RepairCache      bool                    `json:"repair_cache"`
	ProjectConfig    *models.ProjectConfig   `json:"project_config,omitempty"`
	OutputPath       string                  `json:"output_path,omitempty"`
	ProgressCallback func(*WorkflowProgress) `json:"-"`
	Timeout          time.Duration           `json:"timeout"`
}

// Workflow result structures

// ProjectWorkflowResult contains the result of project generation workflow.
type ProjectWorkflowResult struct {
	Success          bool              `json:"success"`
	WorkflowID       string            `json:"workflow_id"`
	ProjectPath      string            `json:"project_path"`
	GeneratedFiles   []string          `json:"generated_files"`
	ValidationResult *ValidationResult `json:"validation_result,omitempty"`
	AuditResult      *AuditResult      `json:"audit_result,omitempty"`
	ReportPath       string            `json:"report_path,omitempty"`
	Duration         time.Duration     `json:"duration"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          time.Time         `json:"end_time"`
	Progress         *WorkflowProgress `json:"progress"`
	Errors           []string          `json:"errors,omitempty"`
	Warnings         []string          `json:"warnings,omitempty"`
}

// ValidationWorkflowResult contains the result of validation workflow.
type ValidationWorkflowResult struct {
	Success          bool              `json:"success"`
	WorkflowID       string            `json:"workflow_id"`
	ProjectPath      string            `json:"project_path"`
	ValidationResult *ValidationResult `json:"validation_result"`
	FixesApplied     []Fix             `json:"fixes_applied,omitempty"`
	ReportPath       string            `json:"report_path,omitempty"`
	Duration         time.Duration     `json:"duration"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          time.Time         `json:"end_time"`
	Progress         *WorkflowProgress `json:"progress"`
	Errors           []string          `json:"errors,omitempty"`
	Warnings         []string          `json:"warnings,omitempty"`
}

// AuditWorkflowResult contains the result of audit workflow.
type AuditWorkflowResult struct {
	Success     bool              `json:"success"`
	WorkflowID  string            `json:"workflow_id"`
	ProjectPath string            `json:"project_path"`
	AuditResult *AuditResult      `json:"audit_result"`
	ReportPath  string            `json:"report_path,omitempty"`
	Duration    time.Duration     `json:"duration"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Progress    *WorkflowProgress `json:"progress"`
	Errors      []string          `json:"errors,omitempty"`
	Warnings    []string          `json:"warnings,omitempty"`
}

// ValidationAuditWorkflowResult contains the result of combined validation and audit workflow.
type ValidationAuditWorkflowResult struct {
	Success          bool              `json:"success"`
	WorkflowID       string            `json:"workflow_id"`
	ProjectPath      string            `json:"project_path"`
	ValidationResult *ValidationResult `json:"validation_result,omitempty"`
	AuditResult      *AuditResult      `json:"audit_result,omitempty"`
	FixesApplied     []Fix             `json:"fixes_applied,omitempty"`
	ReportPath       string            `json:"report_path,omitempty"`
	Duration         time.Duration     `json:"duration"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          time.Time         `json:"end_time"`
	Progress         *WorkflowProgress `json:"progress"`
	Errors           []string          `json:"errors,omitempty"`
	Warnings         []string          `json:"warnings,omitempty"`
}

// ConfigurationWorkflowResult contains the result of configuration workflow.
type ConfigurationWorkflowResult struct {
	Success          bool                    `json:"success"`
	WorkflowID       string                  `json:"workflow_id"`
	Operation        string                  `json:"operation"`
	ConfigPath       string                  `json:"config_path"`
	OutputPath       string                  `json:"output_path"`
	Configuration    *models.ProjectConfig   `json:"configuration,omitempty"`
	ValidationResult *ConfigValidationResult `json:"validation_result,omitempty"`
	Duration         time.Duration           `json:"duration"`
	StartTime        time.Time               `json:"start_time"`
	EndTime          time.Time               `json:"end_time"`
	Progress         *WorkflowProgress       `json:"progress"`
	Errors           []string                `json:"errors,omitempty"`
	Warnings         []string                `json:"warnings,omitempty"`
}

// OfflineWorkflowResult contains the result of offline workflow.
type OfflineWorkflowResult struct {
	Success     bool              `json:"success"`
	WorkflowID  string            `json:"workflow_id"`
	Operation   string            `json:"operation"`
	CacheStats  *CacheStats       `json:"cache_stats,omitempty"`
	ProjectPath string            `json:"project_path,omitempty"`
	Duration    time.Duration     `json:"duration"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Progress    *WorkflowProgress `json:"progress"`
	Errors      []string          `json:"errors,omitempty"`
	Warnings    []string          `json:"warnings,omitempty"`
}

// Workflow management structures

// WorkflowProgress tracks the progress of workflow execution.
type WorkflowProgress struct {
	WorkflowID   string                 `json:"workflow_id"`
	Stage        string                 `json:"stage"`
	Step         string                 `json:"step"`
	Progress     float64                `json:"progress"` // 0.0 to 100.0
	Message      string                 `json:"message"`
	StartTime    time.Time              `json:"start_time"`
	ElapsedTime  time.Duration          `json:"elapsed_time"`
	EstimatedETA time.Duration          `json:"estimated_eta,omitempty"`
	Errors       []string               `json:"errors,omitempty"`
	Warnings     []string               `json:"warnings,omitempty"`
	Details      map[string]interface{} `json:"details,omitempty"`
}

// WorkflowStatus represents the current status of a workflow.
type WorkflowStatus struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"` // pending, running, completed, failed, cancelled
	Progress  *WorkflowProgress      `json:"progress"`
	StartTime time.Time              `json:"start_time"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowInfo contains information about a workflow.
type WorkflowInfo struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	ProjectPath string                 `json:"project_path,omitempty"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Workflow status constants
const (
	WorkflowStatusPending   = "pending"
	WorkflowStatusRunning   = "running"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusFailed    = "failed"
	WorkflowStatusCancelled = "cancelled"
)

// Workflow type constants
const (
	WorkflowTypeProjectGeneration = "project_generation"
	WorkflowTypeValidation        = "validation"
	WorkflowTypeAudit             = "audit"
	WorkflowTypeValidationAudit   = "validation_audit"
	WorkflowTypeConfiguration     = "configuration"
	WorkflowTypeOffline           = "offline"
)

// Workflow stage constants
const (
	WorkflowStageInitialization = "initialization"
	WorkflowStageValidation     = "validation"
	WorkflowStageTemplates      = "templates"
	WorkflowStageGeneration     = "generation"
	WorkflowStageProcessing     = "processing"
	WorkflowStageAudit          = "audit"
	WorkflowStageReporting      = "reporting"
	WorkflowStageCompletion     = "completion"
)
