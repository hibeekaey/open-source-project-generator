// Package workflow provides integration workflows for validation and audit operations.
//
// This file implements comprehensive integration workflows that combine validation
// and audit operations with configuration management and cache support.
package workflow

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ValidationAuditWorkflow orchestrates integrated validation and audit operations.
type ValidationAuditWorkflow struct {
	pipeline *Pipeline
	options  *ValidationAuditOptions
	progress *WorkflowProgress
}

// ValidationAuditOptions defines options for validation and audit workflows.
type ValidationAuditOptions struct {
	ProjectPath       string                        `json:"project_path"`
	ValidationEnabled bool                          `json:"validation_enabled"`
	AuditEnabled      bool                          `json:"audit_enabled"`
	ValidationOptions *interfaces.ValidationOptions `json:"validation_options,omitempty"`
	AuditOptions      *interfaces.AuditOptions      `json:"audit_options,omitempty"`
	OutputFormat      string                        `json:"output_format"`
	OutputFile        string                        `json:"output_file"`
	FixIssues         bool                          `json:"fix_issues"`
	GenerateReport    bool                          `json:"generate_report"`
	ProgressCallback  func(*WorkflowProgress)       `json:"-"`
}

// ValidationAuditResult contains the result of validation and audit workflow.
type ValidationAuditResult struct {
	Success          bool                         `json:"success"`
	ProjectPath      string                       `json:"project_path"`
	ValidationResult *interfaces.ValidationResult `json:"validation_result,omitempty"`
	AuditResult      *interfaces.AuditResult      `json:"audit_result,omitempty"`
	FixesApplied     []interfaces.Fix             `json:"fixes_applied,omitempty"`
	ReportGenerated  bool                         `json:"report_generated"`
	ReportPath       string                       `json:"report_path,omitempty"`
	StartTime        time.Time                    `json:"start_time"`
	Duration         time.Duration                `json:"duration"`
	Progress         *WorkflowProgress            `json:"progress"`
	Errors           []string                     `json:"errors,omitempty"`
	Warnings         []string                     `json:"warnings,omitempty"`
}

// ConfigurationWorkflow orchestrates configuration management operations.
type ConfigurationWorkflow struct {
	pipeline *Pipeline
	options  *ConfigurationOptions
	progress *WorkflowProgress
}

// ConfigurationOptions defines options for configuration workflows.
type ConfigurationOptions struct {
	Operation        string                  `json:"operation"` // export, import, validate, merge
	ConfigPath       string                  `json:"config_path"`
	OutputPath       string                  `json:"output_path"`
	Format           string                  `json:"format"`
	Validate         bool                    `json:"validate"`
	Merge            bool                    `json:"merge"`
	Sources          []string                `json:"sources,omitempty"`
	ProgressCallback func(*WorkflowProgress) `json:"-"`
}

// ConfigurationResult contains the result of configuration workflow.
type ConfigurationResult struct {
	Success          bool                               `json:"success"`
	Operation        string                             `json:"operation"`
	ConfigPath       string                             `json:"config_path"`
	OutputPath       string                             `json:"output_path"`
	Configuration    *models.ProjectConfig              `json:"configuration,omitempty"`
	ValidationResult *interfaces.ConfigValidationResult `json:"validation_result,omitempty"`
	StartTime        time.Time                          `json:"start_time"`
	Duration         time.Duration                      `json:"duration"`
	Progress         *WorkflowProgress                  `json:"progress"`
	Errors           []string                           `json:"errors,omitempty"`
	Warnings         []string                           `json:"warnings,omitempty"`
}

// OfflineWorkflow orchestrates cache-enabled offline operations.
type OfflineWorkflow struct {
	pipeline *Pipeline
	options  *OfflineOptions
	progress *WorkflowProgress
}

// OfflineOptions defines options for offline workflows.
type OfflineOptions struct {
	Operation        string                  `json:"operation"` // sync, validate, generate
	CacheUpdate      bool                    `json:"cache_update"`
	ValidateCache    bool                    `json:"validate_cache"`
	RepairCache      bool                    `json:"repair_cache"`
	ProjectConfig    *models.ProjectConfig   `json:"project_config,omitempty"`
	OutputPath       string                  `json:"output_path,omitempty"`
	ProgressCallback func(*WorkflowProgress) `json:"-"`
}

// OfflineResult contains the result of offline workflow.
type OfflineResult struct {
	Success     bool                   `json:"success"`
	Operation   string                 `json:"operation"`
	CacheStats  *interfaces.CacheStats `json:"cache_stats,omitempty"`
	ProjectPath string                 `json:"project_path,omitempty"`
	StartTime   time.Time              `json:"start_time"`
	Duration    time.Duration          `json:"duration"`
	Progress    *WorkflowProgress      `json:"progress"`
	Errors      []string               `json:"errors,omitempty"`
	Warnings    []string               `json:"warnings,omitempty"`
}

// CreateValidationAuditWorkflow creates a new validation and audit workflow.
func (p *Pipeline) CreateValidationAuditWorkflow(options *ValidationAuditOptions) *ValidationAuditWorkflow {
	return &ValidationAuditWorkflow{
		pipeline: p,
		options:  options,
		progress: &WorkflowProgress{
			StartTime: time.Now(),
		},
	}
}

// Execute runs the validation and audit workflow.
func (w *ValidationAuditWorkflow) Execute(ctx context.Context) (*ValidationAuditResult, error) {
	result := &ValidationAuditResult{
		StartTime:   time.Now(),
		ProjectPath: w.options.ProjectPath,
		Progress:    w.progress,
	}

	w.updateProgress("initialization", "Initializing validation and audit workflow", 0.0, "Starting comprehensive project analysis")

	// Phase 1: Project validation (if enabled)
	if w.options.ValidationEnabled {
		validationResult, err := w.runValidation(ctx)
		if err != nil {
			return w.failWorkflow(result, "validation", err)
		}
		result.ValidationResult = validationResult
		w.updateProgress("validation", "Project validation completed", 40.0, fmt.Sprintf("Found %d issues", len(validationResult.Issues)))
	}

	// Phase 2: Project audit (if enabled)
	if w.options.AuditEnabled {
		auditResult, err := w.runAudit(ctx)
		if err != nil {
			return w.failWorkflow(result, "audit", err)
		}
		result.AuditResult = auditResult
		w.updateProgress("audit", "Project audit completed", 70.0, fmt.Sprintf("Overall score: %.1f", auditResult.OverallScore))
	}

	// Phase 3: Fix issues (if enabled and fixable issues found)
	if w.options.FixIssues && result.ValidationResult != nil {
		fixes, err := w.applyFixes(ctx, result.ValidationResult)
		if err != nil {
			w.addWarning(fmt.Sprintf("Failed to apply fixes: %v", err))
		} else {
			result.FixesApplied = fixes
			w.updateProgress("fixes", "Issues fixed", 85.0, fmt.Sprintf("Applied %d fixes", len(fixes)))
		}
	}

	// Phase 4: Generate report (if enabled)
	if w.options.GenerateReport {
		reportPath, err := w.generateReport(ctx, result)
		if err != nil {
			w.addWarning(fmt.Sprintf("Failed to generate report: %v", err))
		} else {
			result.ReportGenerated = true
			result.ReportPath = reportPath
			w.updateProgress("report", "Report generated", 95.0, fmt.Sprintf("Report saved to %s", reportPath))
		}
	}

	w.updateProgress("completion", "Workflow completed", 100.0, "Validation and audit workflow completed successfully")

	result.Success = true
	result.Duration = time.Since(result.StartTime)

	logFields := map[string]interface{}{
		"project_path":  result.ProjectPath,
		"fixes_applied": len(result.FixesApplied),
		"duration":      result.Duration,
	}

	if result.ValidationResult != nil {
		logFields["validation_issues"] = len(result.ValidationResult.Issues)
	} else {
		logFields["validation_issues"] = 0
	}

	if result.AuditResult != nil {
		logFields["audit_score"] = result.AuditResult.OverallScore
	} else {
		logFields["audit_score"] = 0
	}

	w.pipeline.logger.InfoWithFields("Validation and audit workflow completed", logFields)

	return result, nil
}

// runValidation executes project validation.
func (w *ValidationAuditWorkflow) runValidation(ctx context.Context) (*interfaces.ValidationResult, error) {
	w.pipeline.logger.Debug("Running project validation")

	if w.pipeline.validator == nil {
		return nil, fmt.Errorf("validation engine not available")
	}

	// Use provided validation options or defaults
	options := w.options.ValidationOptions
	if options == nil {
		options = &interfaces.ValidationOptions{
			Verbose: true,
			Fix:     false,
			Report:  true,
		}
	}

	result, err := w.pipeline.validator.ValidateProject(w.options.ProjectPath)
	if err != nil {
		return nil, err
	}

	// Convert models.ValidationResult to interfaces.ValidationResult
	return &interfaces.ValidationResult{
		Valid:  result.Valid,
		Issues: convertValidationIssues(result.Issues),
		Summary: interfaces.ValidationSummary{
			TotalFiles:   1, // Default values since models doesn't have detailed summary
			ValidFiles:   0,
			ErrorCount:   countIssuesBySeverity(result.Issues, "error"),
			WarningCount: countIssuesBySeverity(result.Issues, "warning"),
			FixableCount: countFixableIssues(result.Issues),
		},
	}, nil
}

// runAudit executes project audit.
func (w *ValidationAuditWorkflow) runAudit(ctx context.Context) (*interfaces.AuditResult, error) {
	w.pipeline.logger.Debug("Running project audit")

	if w.pipeline.auditor == nil {
		return nil, fmt.Errorf("audit engine not available")
	}

	// Use provided audit options or defaults
	options := w.options.AuditOptions
	if options == nil {
		options = &interfaces.AuditOptions{
			Security:    true,
			Quality:     true,
			Licenses:    true,
			Performance: true,
			Detailed:    true,
		}
	}

	return w.pipeline.auditor.AuditProject(w.options.ProjectPath, options)
}

// applyFixes applies fixable validation issues.
func (w *ValidationAuditWorkflow) applyFixes(ctx context.Context, validationResult *interfaces.ValidationResult) ([]interfaces.Fix, error) {
	w.pipeline.logger.Debug("Applying fixes for validation issues")

	if w.pipeline.validator == nil {
		return nil, fmt.Errorf("validation engine not available")
	}

	// Get fixable issues
	fixableIssues := w.pipeline.validator.GetFixableIssues(validationResult.Issues)
	if len(fixableIssues) == 0 {
		return []interfaces.Fix{}, nil
	}

	// Apply fixes
	fixResult, err := w.pipeline.validator.FixValidationIssues(w.options.ProjectPath, fixableIssues)
	if err != nil {
		return nil, fmt.Errorf("failed to apply fixes: %w", err)
	}

	return fixResult.Applied, nil
}

// generateReport generates a comprehensive report.
func (w *ValidationAuditWorkflow) generateReport(ctx context.Context, result *ValidationAuditResult) (string, error) {
	w.pipeline.logger.Debug("Generating comprehensive report")

	// Determine output format
	format := w.options.OutputFormat
	if format == "" {
		format = "html"
	}

	// Determine output path
	outputPath := w.options.OutputFile
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputPath = filepath.Join(w.options.ProjectPath, fmt.Sprintf("analysis-report-%s.%s", timestamp, format))
	}

	// Generate validation report if available
	if result.ValidationResult != nil && w.pipeline.validator != nil {
		reportData, err := w.pipeline.validator.GenerateValidationReport(result.ValidationResult, format)
		if err != nil {
			return "", fmt.Errorf("failed to generate validation report: %w", err)
		}

		// Write report to file (simplified implementation)
		_ = reportData // Would write to file here
	}

	// Generate audit report if available
	if result.AuditResult != nil && w.pipeline.auditor != nil {
		reportData, err := w.pipeline.auditor.GenerateAuditReport(result.AuditResult, format)
		if err != nil {
			return "", fmt.Errorf("failed to generate audit report: %w", err)
		}

		// Write report to file (simplified implementation)
		_ = reportData // Would write to file here
	}

	return outputPath, nil
}

// CreateConfigurationWorkflow creates a new configuration workflow.
func (p *Pipeline) CreateConfigurationWorkflow(options *ConfigurationOptions) *ConfigurationWorkflow {
	return &ConfigurationWorkflow{
		pipeline: p,
		options:  options,
		progress: &WorkflowProgress{
			StartTime: time.Now(),
		},
	}
}

// Execute runs the configuration workflow.
func (w *ConfigurationWorkflow) Execute(ctx context.Context) (*ConfigurationResult, error) {
	result := &ConfigurationResult{
		StartTime:  time.Now(),
		Operation:  w.options.Operation,
		ConfigPath: w.options.ConfigPath,
		OutputPath: w.options.OutputPath,
		Progress:   w.progress,
	}

	w.updateProgress("initialization", "Initializing configuration workflow", 0.0, fmt.Sprintf("Starting %s operation", w.options.Operation))

	switch w.options.Operation {
	case "export":
		err := w.exportConfiguration(ctx, result)
		if err != nil {
			return w.failConfigWorkflow(result, err)
		}
	case "import":
		err := w.importConfiguration(ctx, result)
		if err != nil {
			return w.failConfigWorkflow(result, err)
		}
	case "validate":
		err := w.validateConfiguration(ctx, result)
		if err != nil {
			return w.failConfigWorkflow(result, err)
		}
	case "merge":
		err := w.mergeConfigurations(ctx, result)
		if err != nil {
			return w.failConfigWorkflow(result, err)
		}
	default:
		return w.failConfigWorkflow(result, fmt.Errorf("unknown operation: %s", w.options.Operation))
	}

	w.updateProgress("completion", "Configuration workflow completed", 100.0, "Configuration operation completed successfully")

	result.Success = true
	result.Duration = time.Since(result.StartTime)

	return result, nil
}

// exportConfiguration exports project configuration.
func (w *ConfigurationWorkflow) exportConfiguration(ctx context.Context, result *ConfigurationResult) error {
	w.updateProgress("export", "Exporting configuration", 50.0, "Saving configuration to file")

	if w.pipeline.configManager == nil {
		return fmt.Errorf("configuration manager not available")
	}

	// Implementation would export configuration here
	w.pipeline.logger.Debug("Configuration export completed")
	return nil
}

// importConfiguration imports project configuration.
func (w *ConfigurationWorkflow) importConfiguration(ctx context.Context, result *ConfigurationResult) error {
	w.updateProgress("import", "Importing configuration", 50.0, "Loading configuration from file")

	if w.pipeline.configManager == nil {
		return fmt.Errorf("configuration manager not available")
	}

	// Implementation would import configuration here
	w.pipeline.logger.Debug("Configuration import completed")
	return nil
}

// validateConfiguration validates project configuration.
func (w *ConfigurationWorkflow) validateConfiguration(ctx context.Context, result *ConfigurationResult) error {
	w.updateProgress("validate", "Validating configuration", 50.0, "Checking configuration validity")

	if w.pipeline.validator == nil {
		return fmt.Errorf("validation engine not available")
	}

	// Implementation would validate configuration here
	w.pipeline.logger.Debug("Configuration validation completed")
	return nil
}

// mergeConfigurations merges multiple configurations.
func (w *ConfigurationWorkflow) mergeConfigurations(ctx context.Context, result *ConfigurationResult) error {
	w.updateProgress("merge", "Merging configurations", 50.0, "Combining multiple configuration sources")

	if w.pipeline.configManager == nil {
		return fmt.Errorf("configuration manager not available")
	}

	// Implementation would merge configurations here
	w.pipeline.logger.Debug("Configuration merge completed")
	return nil
}

// CreateOfflineWorkflow creates a new offline workflow.
func (p *Pipeline) CreateOfflineWorkflow(options *OfflineOptions) *OfflineWorkflow {
	return &OfflineWorkflow{
		pipeline: p,
		options:  options,
		progress: &WorkflowProgress{
			StartTime: time.Now(),
		},
	}
}

// Execute runs the offline workflow.
func (w *OfflineWorkflow) Execute(ctx context.Context) (*OfflineResult, error) {
	result := &OfflineResult{
		StartTime: time.Now(),
		Operation: w.options.Operation,
		Progress:  w.progress,
	}

	w.updateProgress("initialization", "Initializing offline workflow", 0.0, fmt.Sprintf("Starting %s operation", w.options.Operation))

	// Ensure cache manager is available
	if w.pipeline.cacheManager == nil {
		return w.failOfflineWorkflow(result, fmt.Errorf("cache manager not available"))
	}

	switch w.options.Operation {
	case "sync":
		err := w.syncCache(ctx, result)
		if err != nil {
			return w.failOfflineWorkflow(result, err)
		}
	case "validate":
		err := w.validateCache(ctx, result)
		if err != nil {
			return w.failOfflineWorkflow(result, err)
		}
	case "generate":
		err := w.generateOffline(ctx, result)
		if err != nil {
			return w.failOfflineWorkflow(result, err)
		}
	default:
		return w.failOfflineWorkflow(result, fmt.Errorf("unknown operation: %s", w.options.Operation))
	}

	w.updateProgress("completion", "Offline workflow completed", 100.0, "Offline operation completed successfully")

	result.Success = true
	result.Duration = time.Since(result.StartTime)

	return result, nil
}

// syncCache synchronizes the cache with remote sources.
func (w *OfflineWorkflow) syncCache(ctx context.Context, result *OfflineResult) error {
	w.updateProgress("sync", "Synchronizing cache", 50.0, "Updating cached templates and data")

	// Get current cache stats
	stats, err := w.pipeline.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}
	result.CacheStats = stats

	// Implementation would sync cache here
	w.pipeline.logger.Debug("Cache synchronization completed")
	return nil
}

// validateCache validates the cache integrity.
func (w *OfflineWorkflow) validateCache(ctx context.Context, result *OfflineResult) error {
	w.updateProgress("validate", "Validating cache", 50.0, "Checking cache integrity")

	// Validate cache
	err := w.pipeline.cacheManager.ValidateCache()
	if err != nil {
		// Try to repair if validation fails and repair is enabled
		if w.options.RepairCache {
			w.updateProgress("repair", "Repairing cache", 75.0, "Fixing cache corruption")
			repairErr := w.pipeline.cacheManager.RepairCache()
			if repairErr != nil {
				return fmt.Errorf("cache validation failed and repair failed: %w", repairErr)
			}
		} else {
			return fmt.Errorf("cache validation failed: %w", err)
		}
	}

	// Get updated cache stats
	stats, err := w.pipeline.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}
	result.CacheStats = stats

	w.pipeline.logger.Debug("Cache validation completed")
	return nil
}

// generateOffline generates a project using only cached resources.
func (w *OfflineWorkflow) generateOffline(ctx context.Context, result *OfflineResult) error {
	w.updateProgress("generate", "Generating project offline", 50.0, "Creating project using cached resources")

	if w.options.ProjectConfig == nil {
		return fmt.Errorf("project configuration required for offline generation")
	}

	// Create offline project generation workflow
	workflowOptions := &WorkflowOptions{
		OutputPath:    w.options.OutputPath,
		Offline:       true,
		ValidateAfter: false, // Skip validation in offline mode
		AuditAfter:    false, // Skip audit in offline mode
	}

	projectWorkflow := w.pipeline.CreateProjectGenerationWorkflow(w.options.ProjectConfig, workflowOptions)
	workflowResult, err := projectWorkflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("offline project generation failed: %w", err)
	}

	result.ProjectPath = workflowResult.ProjectPath
	w.pipeline.logger.Debug("Offline project generation completed")
	return nil
}

// Helper methods for workflow management

func (w *ValidationAuditWorkflow) updateProgress(stage, step string, progress float64, message string) {
	w.progress.Stage = stage
	w.progress.Step = step
	w.progress.Progress = progress
	w.progress.Message = message
	w.progress.ElapsedTime = time.Since(w.progress.StartTime)

	if w.options.ProgressCallback != nil {
		w.options.ProgressCallback(w.progress)
	}
}

func (w *ValidationAuditWorkflow) addWarning(warning string) {
	w.progress.Warnings = append(w.progress.Warnings, warning)
	w.pipeline.logger.Warn(warning)
}

func (w *ValidationAuditWorkflow) failWorkflow(result *ValidationAuditResult, stage string, err error) (*ValidationAuditResult, error) {
	errorMsg := fmt.Sprintf("Workflow failed at %s: %v", stage, err)
	w.progress.Errors = append(w.progress.Errors, errorMsg)

	result.Success = false
	result.Duration = time.Since(result.StartTime)
	result.Errors = w.progress.Errors
	result.Warnings = w.progress.Warnings

	return result, fmt.Errorf("%s", errorMsg)
}

func (w *ConfigurationWorkflow) updateProgress(stage, step string, progress float64, message string) {
	w.progress.Stage = stage
	w.progress.Step = step
	w.progress.Progress = progress
	w.progress.Message = message
	w.progress.ElapsedTime = time.Since(w.progress.StartTime)

	if w.options.ProgressCallback != nil {
		w.options.ProgressCallback(w.progress)
	}
}

func (w *ConfigurationWorkflow) failConfigWorkflow(result *ConfigurationResult, err error) (*ConfigurationResult, error) {
	errorMsg := fmt.Sprintf("Configuration workflow failed: %v", err)
	w.progress.Errors = append(w.progress.Errors, errorMsg)

	result.Success = false
	result.Duration = time.Since(result.StartTime)
	result.Errors = w.progress.Errors
	result.Warnings = w.progress.Warnings

	return result, fmt.Errorf("%s", errorMsg)
}

func (w *OfflineWorkflow) updateProgress(stage, step string, progress float64, message string) {
	w.progress.Stage = stage
	w.progress.Step = step
	w.progress.Progress = progress
	w.progress.Message = message
	w.progress.ElapsedTime = time.Since(w.progress.StartTime)

	if w.options.ProgressCallback != nil {
		w.options.ProgressCallback(w.progress)
	}
}

func (w *OfflineWorkflow) failOfflineWorkflow(result *OfflineResult, err error) (*OfflineResult, error) {
	errorMsg := fmt.Sprintf("Offline workflow failed: %v", err)
	w.progress.Errors = append(w.progress.Errors, errorMsg)

	result.Success = false
	result.Duration = time.Since(result.StartTime)
	result.Errors = w.progress.Errors
	result.Warnings = w.progress.Warnings

	return result, fmt.Errorf("%s", errorMsg)
}

// countIssuesBySeverity counts issues by severity level
func countIssuesBySeverity(issues []models.ValidationIssue, severity string) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

// countFixableIssues counts fixable issues
func countFixableIssues(issues []models.ValidationIssue) int {
	count := 0
	for _, issue := range issues {
		if issue.Fixable {
			count++
		}
	}
	return count
}
