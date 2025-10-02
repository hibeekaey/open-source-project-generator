// Package workflow provides concrete implementations of workflow interfaces.
//
// This file contains the concrete implementations of all workflow types
// including project generation, validation, audit, configuration, and offline workflows.
package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// projectWorkflowImpl implements the ProjectWorkflow interface.
type projectWorkflowImpl struct {
	id       string
	manager  *Manager
	config   *models.ProjectConfig
	options  *WorkflowOptions
	status   string
	pipeline *ProjectGenerationWorkflow
	progress *interfaces.WorkflowProgress
}

// Execute runs the project generation workflow.
func (w *projectWorkflowImpl) Execute(ctx context.Context) (*interfaces.ProjectWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	// Execute the pipeline workflow
	result, err := w.pipeline.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("project generation workflow failed: %w", err)
	}

	// Convert internal result to interface result
	interfaceResult := &interfaces.ProjectWorkflowResult{
		Success:          result.Success,
		WorkflowID:       w.id,
		ProjectPath:      result.ProjectPath,
		GeneratedFiles:   result.GeneratedFiles,
		ValidationResult: result.ValidationResult,
		AuditResult:      result.AuditResult,
		Duration:         result.Duration,
		StartTime:        time.Now().Add(-result.Duration),
		EndTime:          time.Now(),
		Progress:         w.convertProgress(result.Progress),
		Errors:           result.Errors,
		Warnings:         result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, interfaceResult.Progress, nil)

	return interfaceResult, nil
}

// Cancel cancels the workflow.
func (w *projectWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *projectWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *projectWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := w.convertProgress(progress)
			callback(interfaceProgress)
		}
	}
}

// GetConfiguration returns the project configuration.
func (w *projectWorkflowImpl) GetConfiguration() *models.ProjectConfig {
	return w.config
}

// UpdateConfiguration updates the project configuration.
func (w *projectWorkflowImpl) UpdateConfiguration(config *models.ProjectConfig) error {
	w.config = config
	return nil
}

// GetOptions returns the workflow options.
func (w *projectWorkflowImpl) GetOptions() *interfaces.ProjectWorkflowOptions {
	return &interfaces.ProjectWorkflowOptions{
		OutputPath:     w.options.OutputPath,
		DryRun:         w.options.DryRun,
		Offline:        w.options.Offline,
		Force:          w.options.Force,
		BackupExisting: w.options.BackupExisting,
		ValidateAfter:  w.options.ValidateAfter,
		AuditAfter:     w.options.AuditAfter,
	}
}

// UpdateOptions updates the workflow options.
func (w *projectWorkflowImpl) UpdateOptions(options *interfaces.ProjectWorkflowOptions) error {
	w.options.OutputPath = options.OutputPath
	w.options.DryRun = options.DryRun
	w.options.Offline = options.Offline
	w.options.Force = options.Force
	w.options.BackupExisting = options.BackupExisting
	w.options.ValidateAfter = options.ValidateAfter
	w.options.AuditAfter = options.AuditAfter
	return nil
}

// GetStatus returns the workflow status.
func (w *projectWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeProjectGeneration,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *projectWorkflowImpl) GetID() string {
	return w.id
}

// convertProgress converts internal progress to interface progress.
func (w *projectWorkflowImpl) convertProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}

// validationWorkflowImpl implements the ValidationWorkflow interface.
type validationWorkflowImpl struct {
	id          string
	manager     *Manager
	projectPath string
	options     *interfaces.ValidationWorkflowOptions
	status      string
	progress    *interfaces.WorkflowProgress
}

// Execute runs the validation workflow.
func (w *validationWorkflowImpl) Execute(ctx context.Context) (*interfaces.ValidationWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	startTime := time.Now()

	// Create validation workflow options
	validationAuditOptions := &ValidationAuditOptions{
		ProjectPath:       w.projectPath,
		ValidationEnabled: true,
		AuditEnabled:      false,
		ValidationOptions: w.options.ValidationOptions,
		FixIssues:         w.options.FixIssues,
		GenerateReport:    w.options.GenerateReport,
		OutputFormat:      w.options.OutputFormat,
		OutputFile:        w.options.OutputFile,
	}

	// Execute validation using the pipeline
	validationAuditWorkflow := w.manager.pipeline.CreateValidationAuditWorkflow(validationAuditOptions)
	result, err := validationAuditWorkflow.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("validation workflow failed: %w", err)
	}

	// Convert to validation workflow result
	validationResult := &interfaces.ValidationWorkflowResult{
		Success:          result.Success,
		WorkflowID:       w.id,
		ProjectPath:      result.ProjectPath,
		ValidationResult: result.ValidationResult,
		FixesApplied:     result.FixesApplied,
		ReportPath:       result.ReportPath,
		Duration:         time.Since(startTime),
		StartTime:        startTime,
		EndTime:          time.Now(),
		Progress:         w.convertValidationProgress(result.Progress),
		Errors:           result.Errors,
		Warnings:         result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, validationResult.Progress, nil)

	return validationResult, nil
}

// Cancel cancels the validation workflow.
func (w *validationWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *validationWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *validationWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = callback
	}
}

// GetStatus returns the workflow status.
func (w *validationWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeValidation,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *validationWorkflowImpl) GetID() string {
	return w.id
}

// convertValidationProgress converts internal progress to interface progress.
func (w *validationWorkflowImpl) convertValidationProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}

// auditWorkflowImpl implements the AuditWorkflow interface.
type auditWorkflowImpl struct {
	id          string
	manager     *Manager
	projectPath string
	options     *interfaces.AuditWorkflowOptions
	status      string
	progress    *interfaces.WorkflowProgress
}

// Execute runs the audit workflow.
func (w *auditWorkflowImpl) Execute(ctx context.Context) (*interfaces.AuditWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	startTime := time.Now()

	// Create audit workflow options
	validationAuditOptions := &ValidationAuditOptions{
		ProjectPath:       w.projectPath,
		ValidationEnabled: false,
		AuditEnabled:      true,
		AuditOptions:      w.options.AuditOptions,
		GenerateReport:    w.options.GenerateReport,
		OutputFormat:      w.options.OutputFormat,
		OutputFile:        w.options.OutputFile,
	}

	// Execute audit using the pipeline
	validationAuditWorkflow := w.manager.pipeline.CreateValidationAuditWorkflow(validationAuditOptions)
	result, err := validationAuditWorkflow.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("audit workflow failed: %w", err)
	}

	// Convert to audit workflow result
	auditResult := &interfaces.AuditWorkflowResult{
		Success:     result.Success,
		WorkflowID:  w.id,
		ProjectPath: result.ProjectPath,
		AuditResult: result.AuditResult,
		ReportPath:  result.ReportPath,
		Duration:    time.Since(startTime),
		StartTime:   startTime,
		EndTime:     time.Now(),
		Progress:    w.convertAuditProgress(result.Progress),
		Errors:      result.Errors,
		Warnings:    result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, auditResult.Progress, nil)

	return auditResult, nil
}

// Cancel cancels the audit workflow.
func (w *auditWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *auditWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *auditWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = callback
	}
}

// GetStatus returns the workflow status.
func (w *auditWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeAudit,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *auditWorkflowImpl) GetID() string {
	return w.id
}

// convertAuditProgress converts internal progress to interface progress.
func (w *auditWorkflowImpl) convertAuditProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}

// validationAuditWorkflowImpl implements the ValidationAuditWorkflow interface.
type validationAuditWorkflowImpl struct {
	id       string
	manager  *Manager
	options  *ValidationAuditOptions
	status   string
	pipeline *ValidationAuditWorkflow
	progress *interfaces.WorkflowProgress
}

// Execute runs the combined validation and audit workflow.
func (w *validationAuditWorkflowImpl) Execute(ctx context.Context) (*interfaces.ValidationAuditWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	// Execute the pipeline workflow
	result, err := w.pipeline.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("validation audit workflow failed: %w", err)
	}

	// Convert internal result to interface result
	interfaceResult := &interfaces.ValidationAuditWorkflowResult{
		Success:          result.Success,
		WorkflowID:       w.id,
		ProjectPath:      result.ProjectPath,
		ValidationResult: result.ValidationResult,
		AuditResult:      result.AuditResult,
		FixesApplied:     result.FixesApplied,
		ReportPath:       result.ReportPath,
		Duration:         result.Duration,
		StartTime:        time.Now().Add(-result.Duration),
		EndTime:          time.Now(),
		Progress:         w.convertValidationAuditProgress(result.Progress),
		Errors:           result.Errors,
		Warnings:         result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, interfaceResult.Progress, nil)

	return interfaceResult, nil
}

// Cancel cancels the workflow.
func (w *validationAuditWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *validationAuditWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *validationAuditWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := w.convertValidationAuditProgress(progress)
			callback(interfaceProgress)
		}
	}
}

// GetStatus returns the workflow status.
func (w *validationAuditWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeValidationAudit,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *validationAuditWorkflowImpl) GetID() string {
	return w.id
}

// convertValidationAuditProgress converts internal progress to interface progress.
func (w *validationAuditWorkflowImpl) convertValidationAuditProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}

// configurationWorkflowImpl implements the ConfigurationWorkflow interface.
type configurationWorkflowImpl struct {
	id       string
	manager  *Manager
	options  *ConfigurationOptions
	status   string
	pipeline *ConfigurationWorkflow
	progress *interfaces.WorkflowProgress
}

// Execute runs the configuration workflow.
func (w *configurationWorkflowImpl) Execute(ctx context.Context) (*interfaces.ConfigurationWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	// Execute the pipeline workflow
	result, err := w.pipeline.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("configuration workflow failed: %w", err)
	}

	// Convert internal result to interface result
	interfaceResult := &interfaces.ConfigurationWorkflowResult{
		Success:          result.Success,
		WorkflowID:       w.id,
		Operation:        result.Operation,
		ConfigPath:       result.ConfigPath,
		OutputPath:       result.OutputPath,
		Configuration:    result.Configuration,
		ValidationResult: result.ValidationResult,
		Duration:         result.Duration,
		StartTime:        time.Now().Add(-result.Duration),
		EndTime:          time.Now(),
		Progress:         w.convertConfigProgress(result.Progress),
		Errors:           result.Errors,
		Warnings:         result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, interfaceResult.Progress, nil)

	return interfaceResult, nil
}

// Cancel cancels the workflow.
func (w *configurationWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *configurationWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *configurationWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := w.convertConfigProgress(progress)
			callback(interfaceProgress)
		}
	}
}

// GetStatus returns the workflow status.
func (w *configurationWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeConfiguration,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *configurationWorkflowImpl) GetID() string {
	return w.id
}

// convertConfigProgress converts internal progress to interface progress.
func (w *configurationWorkflowImpl) convertConfigProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}

// offlineWorkflowImpl implements the OfflineWorkflow interface.
type offlineWorkflowImpl struct {
	id       string
	manager  *Manager
	options  *OfflineOptions
	status   string
	pipeline *OfflineWorkflow
	progress *interfaces.WorkflowProgress
}

// Execute runs the offline workflow.
func (w *offlineWorkflowImpl) Execute(ctx context.Context) (*interfaces.OfflineWorkflowResult, error) {
	w.status = interfaces.WorkflowStatusRunning
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)

	// Execute the pipeline workflow
	result, err := w.pipeline.Execute(ctx)
	if err != nil {
		w.status = interfaces.WorkflowStatusFailed
		w.manager.updateWorkflowStatus(w.id, w.status, w.progress, err)
		return nil, fmt.Errorf("offline workflow failed: %w", err)
	}

	// Convert internal result to interface result
	interfaceResult := &interfaces.OfflineWorkflowResult{
		Success:     result.Success,
		WorkflowID:  w.id,
		Operation:   result.Operation,
		CacheStats:  result.CacheStats,
		ProjectPath: result.ProjectPath,
		Duration:    result.Duration,
		StartTime:   time.Now().Add(-result.Duration),
		EndTime:     time.Now(),
		Progress:    w.convertOfflineProgress(result.Progress),
		Errors:      result.Errors,
		Warnings:    result.Warnings,
	}

	w.status = interfaces.WorkflowStatusCompleted
	w.manager.updateWorkflowStatus(w.id, w.status, interfaceResult.Progress, nil)

	return interfaceResult, nil
}

// Cancel cancels the workflow.
func (w *offlineWorkflowImpl) Cancel() error {
	w.status = interfaces.WorkflowStatusCancelled
	w.manager.updateWorkflowStatus(w.id, w.status, w.progress, nil)
	return nil
}

// GetProgress returns the current progress.
func (w *offlineWorkflowImpl) GetProgress() *interfaces.WorkflowProgress {
	return w.progress
}

// SetProgressCallback sets the progress callback.
func (w *offlineWorkflowImpl) SetProgressCallback(callback func(*interfaces.WorkflowProgress)) {
	if w.options != nil {
		w.options.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := w.convertOfflineProgress(progress)
			callback(interfaceProgress)
		}
	}
}

// GetStatus returns the workflow status.
func (w *offlineWorkflowImpl) GetStatus() interfaces.WorkflowStatus {
	return interfaces.WorkflowStatus{
		ID:       w.id,
		Type:     interfaces.WorkflowTypeOffline,
		Status:   w.status,
		Progress: w.progress,
	}
}

// GetID returns the workflow ID.
func (w *offlineWorkflowImpl) GetID() string {
	return w.id
}

// convertOfflineProgress converts internal progress to interface progress.
func (w *offlineWorkflowImpl) convertOfflineProgress(progress *WorkflowProgress) *interfaces.WorkflowProgress {
	if progress == nil {
		return nil
	}

	return &interfaces.WorkflowProgress{
		WorkflowID:  w.id,
		Stage:       progress.Stage,
		Step:        progress.Step,
		Progress:    progress.Progress,
		Message:     progress.Message,
		StartTime:   progress.StartTime,
		ElapsedTime: progress.ElapsedTime,
		Errors:      progress.Errors,
		Warnings:    progress.Warnings,
	}
}
