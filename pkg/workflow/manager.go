// Package workflow provides workflow management and orchestration.
//
// This file implements the WorkflowManager interface and provides comprehensive
// workflow management capabilities for all types of operations.
package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Manager implements the WorkflowManager interface.
type Manager struct {
	pipeline        *Pipeline
	activeWorkflows map[string]interfaces.WorkflowStatus
	workflowHistory []interfaces.WorkflowInfo
	mutex           sync.RWMutex
	logger          interfaces.Logger
}

// NewManager creates a new workflow manager.
func NewManager(
	generator interfaces.FileSystemGenerator,
	templateManager interfaces.TemplateManager,
	validator interfaces.ValidationEngine,
	auditor interfaces.AuditEngine,
	cacheManager interfaces.CacheManager,
	configManager interfaces.ConfigManager,
	logger interfaces.Logger,
) *Manager {
	pipeline := NewPipeline(
		generator,
		templateManager,
		validator,
		auditor,
		cacheManager,
		configManager,
		logger,
	)

	return &Manager{
		pipeline:        pipeline,
		activeWorkflows: make(map[string]interfaces.WorkflowStatus),
		workflowHistory: make([]interfaces.WorkflowInfo, 0),
		logger:          logger,
	}
}

// CreateProjectWorkflow creates a new project generation workflow.
func (m *Manager) CreateProjectWorkflow(
	config *models.ProjectConfig,
	options *interfaces.ProjectWorkflowOptions,
) (interfaces.ProjectWorkflow, error) {
	if config == nil {
		return nil, fmt.Errorf("project configuration is required")
	}

	if options == nil {
		options = &interfaces.ProjectWorkflowOptions{
			OutputPath:     ".",
			ValidateAfter:  true,
			AuditAfter:     false,
			GenerateReport: false,
		}
	}

	// Convert to internal workflow options
	workflowOptions := &WorkflowOptions{
		OutputPath:        options.OutputPath,
		DryRun:            options.DryRun,
		Offline:           options.Offline,
		ValidateAfter:     options.ValidateAfter,
		AuditAfter:        options.AuditAfter,
		BackupExisting:    options.BackupExisting,
		Force:             options.Force,
		ValidationOptions: options.ValidationOptions,
		AuditOptions:      options.AuditOptions,
	}

	// Convert progress callback if provided
	if options.ProgressCallback != nil {
		workflowOptions.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := &interfaces.WorkflowProgress{
				Stage:       progress.Stage,
				Step:        progress.Step,
				Progress:    progress.Progress,
				Message:     progress.Message,
				StartTime:   progress.StartTime,
				ElapsedTime: progress.ElapsedTime,
				Errors:      progress.Errors,
				Warnings:    progress.Warnings,
			}
			options.ProgressCallback(interfaceProgress)
		}
	}

	workflow := &projectWorkflowImpl{
		id:       generateWorkflowID(),
		manager:  m,
		config:   config,
		options:  workflowOptions,
		status:   interfaces.WorkflowStatusPending,
		pipeline: m.pipeline.CreateProjectGenerationWorkflow(config, workflowOptions),
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeProjectGeneration, config.Name)

	return workflow, nil
}

// ExecuteProjectGeneration executes a project generation workflow directly.
func (m *Manager) ExecuteProjectGeneration(
	ctx context.Context,
	config *models.ProjectConfig,
	options *interfaces.ProjectWorkflowOptions,
) (*interfaces.ProjectWorkflowResult, error) {
	workflow, err := m.CreateProjectWorkflow(config, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create project workflow: %w", err)
	}

	return workflow.Execute(ctx)
}

// CreateValidationWorkflow creates a new validation workflow.
func (m *Manager) CreateValidationWorkflow(
	projectPath string,
	options *interfaces.ValidationWorkflowOptions,
) (interfaces.ValidationWorkflow, error) {
	if projectPath == "" {
		return nil, fmt.Errorf("project path is required")
	}

	if options == nil {
		options = &interfaces.ValidationWorkflowOptions{
			ProjectPath:    projectPath,
			FixIssues:      false,
			GenerateReport: false,
			OutputFormat:   "text",
		}
	}

	workflow := &validationWorkflowImpl{
		id:          generateWorkflowID(),
		manager:     m,
		projectPath: projectPath,
		options:     options,
		status:      interfaces.WorkflowStatusPending,
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeValidation, projectPath)

	return workflow, nil
}

// CreateAuditWorkflow creates a new audit workflow.
func (m *Manager) CreateAuditWorkflow(
	projectPath string,
	options *interfaces.AuditWorkflowOptions,
) (interfaces.AuditWorkflow, error) {
	if projectPath == "" {
		return nil, fmt.Errorf("project path is required")
	}

	if options == nil {
		options = &interfaces.AuditWorkflowOptions{
			ProjectPath:    projectPath,
			GenerateReport: false,
			OutputFormat:   "text",
		}
	}

	workflow := &auditWorkflowImpl{
		id:          generateWorkflowID(),
		manager:     m,
		projectPath: projectPath,
		options:     options,
		status:      interfaces.WorkflowStatusPending,
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeAudit, projectPath)

	return workflow, nil
}

// CreateValidationAuditWorkflow creates a new combined validation and audit workflow.
func (m *Manager) CreateValidationAuditWorkflow(
	projectPath string,
	options *interfaces.ValidationAuditWorkflowOptions,
) (interfaces.ValidationAuditWorkflow, error) {
	if projectPath == "" {
		return nil, fmt.Errorf("project path is required")
	}

	if options == nil {
		options = &interfaces.ValidationAuditWorkflowOptions{
			ProjectPath:       projectPath,
			ValidationEnabled: true,
			AuditEnabled:      true,
			FixIssues:         false,
			GenerateReport:    false,
			OutputFormat:      "text",
		}
	}

	// Convert to internal options
	internalOptions := &ValidationAuditOptions{
		ProjectPath:       projectPath,
		ValidationEnabled: options.ValidationEnabled,
		AuditEnabled:      options.AuditEnabled,
		ValidationOptions: options.ValidationOptions,
		AuditOptions:      options.AuditOptions,
		OutputFormat:      options.OutputFormat,
		OutputFile:        options.OutputFile,
		FixIssues:         options.FixIssues,
		GenerateReport:    options.GenerateReport,
	}

	// Convert progress callback if provided
	if options.ProgressCallback != nil {
		internalOptions.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := &interfaces.WorkflowProgress{
				Stage:       progress.Stage,
				Step:        progress.Step,
				Progress:    progress.Progress,
				Message:     progress.Message,
				StartTime:   progress.StartTime,
				ElapsedTime: progress.ElapsedTime,
				Errors:      progress.Errors,
				Warnings:    progress.Warnings,
			}
			options.ProgressCallback(interfaceProgress)
		}
	}

	workflow := &validationAuditWorkflowImpl{
		id:       generateWorkflowID(),
		manager:  m,
		options:  internalOptions,
		status:   interfaces.WorkflowStatusPending,
		pipeline: m.pipeline.CreateValidationAuditWorkflow(internalOptions),
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeValidationAudit, projectPath)

	return workflow, nil
}

// CreateConfigurationWorkflow creates a new configuration workflow.
func (m *Manager) CreateConfigurationWorkflow(
	options *interfaces.ConfigurationWorkflowOptions,
) (interfaces.ConfigurationWorkflow, error) {
	if options == nil {
		return nil, fmt.Errorf("configuration workflow options are required")
	}

	// Convert to internal options
	internalOptions := &ConfigurationOptions{
		Operation:  options.Operation,
		ConfigPath: options.ConfigPath,
		OutputPath: options.OutputPath,
		Format:     options.Format,
		Validate:   options.Validate,
		Merge:      options.Merge,
		Sources:    options.Sources,
	}

	// Convert progress callback if provided
	if options.ProgressCallback != nil {
		internalOptions.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := &interfaces.WorkflowProgress{
				Stage:       progress.Stage,
				Step:        progress.Step,
				Progress:    progress.Progress,
				Message:     progress.Message,
				StartTime:   progress.StartTime,
				ElapsedTime: progress.ElapsedTime,
				Errors:      progress.Errors,
				Warnings:    progress.Warnings,
			}
			options.ProgressCallback(interfaceProgress)
		}
	}

	workflow := &configurationWorkflowImpl{
		id:       generateWorkflowID(),
		manager:  m,
		options:  internalOptions,
		status:   interfaces.WorkflowStatusPending,
		pipeline: m.pipeline.CreateConfigurationWorkflow(internalOptions),
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeConfiguration, options.ConfigPath)

	return workflow, nil
}

// CreateOfflineWorkflow creates a new offline workflow.
func (m *Manager) CreateOfflineWorkflow(
	options *interfaces.OfflineWorkflowOptions,
) (interfaces.OfflineWorkflow, error) {
	if options == nil {
		return nil, fmt.Errorf("offline workflow options are required")
	}

	// Convert to internal options
	internalOptions := &OfflineOptions{
		Operation:     options.Operation,
		CacheUpdate:   options.CacheUpdate,
		ValidateCache: options.ValidateCache,
		RepairCache:   options.RepairCache,
		ProjectConfig: options.ProjectConfig,
		OutputPath:    options.OutputPath,
	}

	// Convert progress callback if provided
	if options.ProgressCallback != nil {
		internalOptions.ProgressCallback = func(progress *WorkflowProgress) {
			interfaceProgress := &interfaces.WorkflowProgress{
				Stage:       progress.Stage,
				Step:        progress.Step,
				Progress:    progress.Progress,
				Message:     progress.Message,
				StartTime:   progress.StartTime,
				ElapsedTime: progress.ElapsedTime,
				Errors:      progress.Errors,
				Warnings:    progress.Warnings,
			}
			options.ProgressCallback(interfaceProgress)
		}
	}

	workflow := &offlineWorkflowImpl{
		id:       generateWorkflowID(),
		manager:  m,
		options:  internalOptions,
		status:   interfaces.WorkflowStatusPending,
		pipeline: m.pipeline.CreateOfflineWorkflow(internalOptions),
	}

	// Register workflow
	m.registerWorkflow(workflow.id, interfaces.WorkflowTypeOffline, options.Operation)

	return workflow, nil
}

// Configuration management methods

// ExportConfiguration exports project configuration.
func (m *Manager) ExportConfiguration(
	ctx context.Context,
	config *models.ProjectConfig,
	outputPath string,
) error {
	options := &interfaces.ConfigurationWorkflowOptions{
		Operation:  "export",
		OutputPath: outputPath,
	}

	workflow, err := m.CreateConfigurationWorkflow(options)
	if err != nil {
		return fmt.Errorf("failed to create configuration workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("configuration export failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("configuration export failed: %v", result.Errors)
	}

	return nil
}

// ImportConfiguration imports project configuration.
func (m *Manager) ImportConfiguration(
	ctx context.Context,
	configPath string,
) (*models.ProjectConfig, error) {
	options := &interfaces.ConfigurationWorkflowOptions{
		Operation:  "import",
		ConfigPath: configPath,
	}

	workflow, err := m.CreateConfigurationWorkflow(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("configuration import failed: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("configuration import failed: %v", result.Errors)
	}

	return result.Configuration, nil
}

// ValidateConfiguration validates project configuration.
func (m *Manager) ValidateConfiguration(
	ctx context.Context,
	config *models.ProjectConfig,
) (*interfaces.ConfigValidationResult, error) {
	options := &interfaces.ConfigurationWorkflowOptions{
		Operation: "validate",
		Validate:  true,
	}

	workflow, err := m.CreateConfigurationWorkflow(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return result.ValidationResult, nil
}

// MergeConfigurations merges multiple configurations.
func (m *Manager) MergeConfigurations(
	ctx context.Context,
	configs []*models.ProjectConfig,
) (*models.ProjectConfig, error) {
	// Convert configs to sources (simplified implementation)
	sources := make([]string, len(configs))
	for i := range configs {
		sources[i] = fmt.Sprintf("config_%d", i)
	}

	options := &interfaces.ConfigurationWorkflowOptions{
		Operation: "merge",
		Merge:     true,
		Sources:   sources,
	}

	workflow, err := m.CreateConfigurationWorkflow(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("configuration merge failed: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("configuration merge failed: %v", result.Errors)
	}

	return result.Configuration, nil
}

// Offline workflow methods

// SyncCache synchronizes the cache with remote sources.
func (m *Manager) SyncCache(ctx context.Context) (*interfaces.CacheStats, error) {
	options := &interfaces.OfflineWorkflowOptions{
		Operation:   "sync",
		CacheUpdate: true,
	}

	workflow, err := m.CreateOfflineWorkflow(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create offline workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("cache sync failed: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("cache sync failed: %v", result.Errors)
	}

	return result.CacheStats, nil
}

// ValidateCache validates the cache integrity.
func (m *Manager) ValidateCache(ctx context.Context) error {
	options := &interfaces.OfflineWorkflowOptions{
		Operation:     "validate",
		ValidateCache: true,
	}

	workflow, err := m.CreateOfflineWorkflow(options)
	if err != nil {
		return fmt.Errorf("failed to create offline workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("cache validation failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("cache validation failed: %v", result.Errors)
	}

	return nil
}

// GenerateOffline generates a project using only cached resources.
func (m *Manager) GenerateOffline(
	ctx context.Context,
	config *models.ProjectConfig,
	outputPath string,
) (*interfaces.ProjectWorkflowResult, error) {
	options := &interfaces.OfflineWorkflowOptions{
		Operation:     "generate",
		ProjectConfig: config,
		OutputPath:    outputPath,
	}

	workflow, err := m.CreateOfflineWorkflow(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create offline workflow: %w", err)
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("offline generation failed: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("offline generation failed: %v", result.Errors)
	}

	// Convert to project workflow result
	return &interfaces.ProjectWorkflowResult{
		Success:     result.Success,
		WorkflowID:  result.WorkflowID,
		ProjectPath: result.ProjectPath,
		Duration:    result.Duration,
		StartTime:   result.StartTime,
		EndTime:     result.EndTime,
		Progress:    result.Progress,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
	}, nil
}

// Workflow management methods

// GetActiveWorkflows returns all currently active workflows.
func (m *Manager) GetActiveWorkflows() []interfaces.WorkflowInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	workflows := make([]interfaces.WorkflowInfo, 0, len(m.activeWorkflows))
	for _, status := range m.activeWorkflows {
		info := interfaces.WorkflowInfo{
			ID:          status.ID,
			Type:        status.Type,
			Status:      status.Status,
			ProjectPath: getProjectPathFromMetadata(status.Metadata),
			StartTime:   status.StartTime,
			EndTime:     status.EndTime,
			Duration:    status.Duration,
			Success:     status.Status == interfaces.WorkflowStatusCompleted,
			Error:       status.Error,
			Metadata:    status.Metadata,
		}
		workflows = append(workflows, info)
	}

	return workflows
}

// CancelWorkflow cancels an active workflow.
func (m *Manager) CancelWorkflow(workflowID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	status, exists := m.activeWorkflows[workflowID]
	if !exists {
		return fmt.Errorf("workflow %s not found", workflowID)
	}

	if status.Status != interfaces.WorkflowStatusRunning {
		return fmt.Errorf("workflow %s is not running (status: %s)", workflowID, status.Status)
	}

	// Update status to cancelled
	status.Status = interfaces.WorkflowStatusCancelled
	endTime := time.Now()
	status.EndTime = &endTime
	status.Duration = endTime.Sub(status.StartTime)
	m.activeWorkflows[workflowID] = status

	// Move to history
	m.moveToHistory(workflowID)

	m.logger.InfoWithFields("Workflow cancelled", map[string]interface{}{
		"workflow_id": workflowID,
		"type":        status.Type,
	})

	return nil
}

// GetWorkflowStatus returns the status of a specific workflow.
func (m *Manager) GetWorkflowStatus(workflowID string) (*interfaces.WorkflowStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	status, exists := m.activeWorkflows[workflowID]
	if !exists {
		// Check history
		for _, info := range m.workflowHistory {
			if info.ID == workflowID {
				return &interfaces.WorkflowStatus{
					ID:        info.ID,
					Type:      info.Type,
					Status:    info.Status,
					StartTime: info.StartTime,
					EndTime:   info.EndTime,
					Duration:  info.Duration,
					Error:     info.Error,
					Metadata:  info.Metadata,
				}, nil
			}
		}
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}

	return &status, nil
}

// GetWorkflowHistory returns the history of completed workflows.
func (m *Manager) GetWorkflowHistory() []interfaces.WorkflowInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the history
	history := make([]interfaces.WorkflowInfo, len(m.workflowHistory))
	copy(history, m.workflowHistory)
	return history
}

// Helper methods

// registerWorkflow registers a new workflow.
func (m *Manager) registerWorkflow(id, workflowType, projectPath string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	status := interfaces.WorkflowStatus{
		ID:        id,
		Type:      workflowType,
		Status:    interfaces.WorkflowStatusPending,
		StartTime: time.Now(),
		Metadata: map[string]interface{}{
			"project_path": projectPath,
		},
	}

	m.activeWorkflows[id] = status

	m.logger.InfoWithFields("Workflow registered", map[string]interface{}{
		"workflow_id":  id,
		"type":         workflowType,
		"project_path": projectPath,
	})
}

// updateWorkflowStatus updates the status of a workflow.
func (m *Manager) updateWorkflowStatus(id, status string, progress *interfaces.WorkflowProgress, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	workflowStatus, exists := m.activeWorkflows[id]
	if !exists {
		return
	}

	workflowStatus.Status = status
	workflowStatus.Progress = progress
	workflowStatus.Duration = time.Since(workflowStatus.StartTime)

	if err != nil {
		workflowStatus.Error = err.Error()
	}

	if status == interfaces.WorkflowStatusCompleted || status == interfaces.WorkflowStatusFailed || status == interfaces.WorkflowStatusCancelled {
		endTime := time.Now()
		workflowStatus.EndTime = &endTime
		workflowStatus.Duration = endTime.Sub(workflowStatus.StartTime)
	}

	m.activeWorkflows[id] = workflowStatus

	// Move completed workflows to history
	if status == interfaces.WorkflowStatusCompleted || status == interfaces.WorkflowStatusFailed || status == interfaces.WorkflowStatusCancelled {
		m.moveToHistory(id)
	}
}

// moveToHistory moves a workflow from active to history.
func (m *Manager) moveToHistory(id string) {
	status, exists := m.activeWorkflows[id]
	if !exists {
		return
	}

	info := interfaces.WorkflowInfo{
		ID:          status.ID,
		Type:        status.Type,
		Status:      status.Status,
		ProjectPath: getProjectPathFromMetadata(status.Metadata),
		StartTime:   status.StartTime,
		EndTime:     status.EndTime,
		Duration:    status.Duration,
		Success:     status.Status == interfaces.WorkflowStatusCompleted,
		Error:       status.Error,
		Metadata:    status.Metadata,
	}

	m.workflowHistory = append(m.workflowHistory, info)
	delete(m.activeWorkflows, id)

	// Keep only the last 100 workflows in history
	if len(m.workflowHistory) > 100 {
		m.workflowHistory = m.workflowHistory[len(m.workflowHistory)-100:]
	}
}

// generateWorkflowID generates a unique workflow ID.
func generateWorkflowID() string {
	return fmt.Sprintf("workflow_%d", time.Now().UnixNano())
}

// getProjectPathFromMetadata extracts project path from workflow metadata.
func getProjectPathFromMetadata(metadata map[string]interface{}) string {
	if metadata == nil {
		return ""
	}

	if path, ok := metadata["project_path"].(string); ok {
		return path
	}

	return ""
}
