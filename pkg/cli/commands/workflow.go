// Package commands provides workflow command implementations for the CLI interface.
//
// This module contains workflow command handlers that integrate end-to-end
// workflows into the CLI interface for comprehensive project operations.
package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/spf13/cobra"
)

// WorkflowCLI defines the CLI methods needed by WorkflowCommand.
type WorkflowCLI interface {
	// Workflow management
	GetWorkflowManager() interfaces.WorkflowManager

	// Output methods
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	QuietOutput(format string, args ...interface{})
	ErrorOutput(format string, args ...interface{})
	WarningOutput(format string, args ...interface{})
	SuccessOutput(format string, args ...interface{})

	// Formatting methods
	Error(text string) string
	Warning(text string) string
	Info(text string) string
	Success(text string) string
	Highlight(text string) string
	Dim(text string) string

	// Utility methods
	IsQuietMode() bool
	OutputMachineReadable(data interface{}, format string) error
	CreateWorkflowError(message string, workflowType string) error
}

// WorkflowCommand handles workflow-related command functionality.
type WorkflowCommand struct {
	cli WorkflowCLI
}

// NewWorkflowCommand creates a new WorkflowCommand instance.
func NewWorkflowCommand(cli WorkflowCLI) *WorkflowCommand {
	return &WorkflowCommand{
		cli: cli,
	}
}

// ExecuteProjectWorkflow handles the project workflow command execution.
func (wc *WorkflowCommand) ExecuteProjectWorkflow(cmd *cobra.Command, args []string) error {
	// Get flags
	configPath, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	offline, _ := cmd.Flags().GetBool("offline")
	force, _ := cmd.Flags().GetBool("force")
	backupExisting, _ := cmd.Flags().GetBool("backup-existing")
	validateAfter, _ := cmd.Flags().GetBool("validate-after")
	auditAfter, _ := cmd.Flags().GetBool("audit-after")
	generateReport, _ := cmd.Flags().GetBool("generate-report")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	wc.cli.VerboseOutput("🚀 Starting project generation workflow")

	// Load project configuration
	var config *models.ProjectConfig
	if configPath != "" {
		wc.cli.VerboseOutput("📋 Loading configuration from %s", configPath)
		// Load configuration from file (implementation would go here)
		config = &models.ProjectConfig{
			Name:        "example-project",
			Description: "Example project generated from workflow",
		}
	} else {
		// Use default configuration or prompt for interactive mode
		config = &models.ProjectConfig{
			Name:        "new-project",
			Description: "New project generated from workflow",
		}
	}

	// Create workflow options
	options := &interfaces.ProjectWorkflowOptions{
		OutputPath:     outputPath,
		DryRun:         dryRun,
		Offline:        offline,
		Force:          force,
		BackupExisting: backupExisting,
		ValidateAfter:  validateAfter,
		AuditAfter:     auditAfter,
		GenerateReport: generateReport,
		Timeout:        timeout,
	}

	// Set progress callback if not in quiet mode
	if !wc.cli.IsQuietMode() {
		options.ProgressCallback = wc.createProgressCallback()
	}

	// Get workflow manager
	workflowManager := wc.cli.GetWorkflowManager()
	if workflowManager == nil {
		return wc.cli.CreateWorkflowError("Workflow manager not available", "project")
	}

	// Create and execute workflow
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result, err := workflowManager.ExecuteProjectGeneration(ctx, config, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Project generation workflow failed."),
			wc.cli.Info(err.Error()))
	}

	// Output results
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return wc.cli.OutputMachineReadable(result, outputFormat)
	}

	wc.outputProjectWorkflowResult(result)
	return nil
}

// ExecuteValidationWorkflow handles the validation workflow command execution.
func (wc *WorkflowCommand) ExecuteValidationWorkflow(cmd *cobra.Command, args []string) error {
	// Get project path from args or flags
	var projectPath string
	if len(args) > 0 {
		projectPath = args[0]
	} else {
		projectPath, _ = cmd.Flags().GetString("project-path")
	}

	if projectPath == "" {
		projectPath = "."
	}

	// Get flags
	fixIssues, _ := cmd.Flags().GetBool("fix-issues")
	generateReport, _ := cmd.Flags().GetBool("generate-report")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	wc.cli.VerboseOutput("🔍 Starting validation workflow for %s", projectPath)

	// Create workflow options
	options := &interfaces.ValidationWorkflowOptions{
		ProjectPath:    projectPath,
		FixIssues:      fixIssues,
		GenerateReport: generateReport,
		OutputFormat:   outputFormat,
		OutputFile:     outputFile,
		Timeout:        timeout,
	}

	// Set progress callback if not in quiet mode
	if !wc.cli.IsQuietMode() {
		options.ProgressCallback = wc.createProgressCallback()
	}

	// Get workflow manager
	workflowManager := wc.cli.GetWorkflowManager()
	if workflowManager == nil {
		return wc.cli.CreateWorkflowError("Workflow manager not available", "validation")
	}

	// Create workflow
	workflow, err := workflowManager.CreateValidationWorkflow(projectPath, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Failed to create validation workflow."),
			wc.cli.Info(err.Error()))
	}

	// Execute workflow
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Validation workflow failed."),
			wc.cli.Info(err.Error()))
	}

	// Output results
	if outputFormat == "json" || outputFormat == "yaml" {
		return wc.cli.OutputMachineReadable(result, outputFormat)
	}

	wc.outputValidationWorkflowResult(result)
	return nil
}

// ExecuteAuditWorkflow handles the audit workflow command execution.
func (wc *WorkflowCommand) ExecuteAuditWorkflow(cmd *cobra.Command, args []string) error {
	// Get project path from args or flags
	var projectPath string
	if len(args) > 0 {
		projectPath = args[0]
	} else {
		projectPath, _ = cmd.Flags().GetString("project-path")
	}

	if projectPath == "" {
		projectPath = "."
	}

	// Get flags
	generateReport, _ := cmd.Flags().GetBool("generate-report")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	wc.cli.VerboseOutput("🔍 Starting audit workflow for %s", projectPath)

	// Create workflow options
	options := &interfaces.AuditWorkflowOptions{
		ProjectPath:    projectPath,
		GenerateReport: generateReport,
		OutputFormat:   outputFormat,
		OutputFile:     outputFile,
		Timeout:        timeout,
	}

	// Set progress callback if not in quiet mode
	if !wc.cli.IsQuietMode() {
		options.ProgressCallback = wc.createProgressCallback()
	}

	// Get workflow manager
	workflowManager := wc.cli.GetWorkflowManager()
	if workflowManager == nil {
		return wc.cli.CreateWorkflowError("Workflow manager not available", "audit")
	}

	// Create workflow
	workflow, err := workflowManager.CreateAuditWorkflow(projectPath, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Failed to create audit workflow."),
			wc.cli.Info(err.Error()))
	}

	// Execute workflow
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Audit workflow failed."),
			wc.cli.Info(err.Error()))
	}

	// Output results
	if outputFormat == "json" || outputFormat == "yaml" {
		return wc.cli.OutputMachineReadable(result, outputFormat)
	}

	wc.outputAuditWorkflowResult(result)
	return nil
}

// ExecuteValidationAuditWorkflow handles the combined validation and audit workflow.
func (wc *WorkflowCommand) ExecuteValidationAuditWorkflow(cmd *cobra.Command, args []string) error {
	// Get project path from args or flags
	var projectPath string
	if len(args) > 0 {
		projectPath = args[0]
	} else {
		projectPath, _ = cmd.Flags().GetString("project-path")
	}

	if projectPath == "" {
		projectPath = "."
	}

	// Get flags
	validationEnabled, _ := cmd.Flags().GetBool("validation")
	auditEnabled, _ := cmd.Flags().GetBool("audit")
	fixIssues, _ := cmd.Flags().GetBool("fix-issues")
	generateReport, _ := cmd.Flags().GetBool("generate-report")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	// Default to both validation and audit if neither is specified
	if !validationEnabled && !auditEnabled {
		validationEnabled = true
		auditEnabled = true
	}

	wc.cli.VerboseOutput("🔍 Starting validation and audit workflow for %s", projectPath)

	// Create workflow options
	options := &interfaces.ValidationAuditWorkflowOptions{
		ProjectPath:       projectPath,
		ValidationEnabled: validationEnabled,
		AuditEnabled:      auditEnabled,
		FixIssues:         fixIssues,
		GenerateReport:    generateReport,
		OutputFormat:      outputFormat,
		OutputFile:        outputFile,
		Timeout:           timeout,
	}

	// Set progress callback if not in quiet mode
	if !wc.cli.IsQuietMode() {
		options.ProgressCallback = wc.createProgressCallback()
	}

	// Get workflow manager
	workflowManager := wc.cli.GetWorkflowManager()
	if workflowManager == nil {
		return wc.cli.CreateWorkflowError("Workflow manager not available", "validation-audit")
	}

	// Create workflow
	workflow, err := workflowManager.CreateValidationAuditWorkflow(projectPath, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Failed to create validation audit workflow."),
			wc.cli.Info(err.Error()))
	}

	// Execute workflow
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result, err := workflow.Execute(ctx)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			wc.cli.Error("Validation audit workflow failed."),
			wc.cli.Info(err.Error()))
	}

	// Output results
	if outputFormat == "json" || outputFormat == "yaml" {
		return wc.cli.OutputMachineReadable(result, outputFormat)
	}

	wc.outputValidationAuditWorkflowResult(result)
	return nil
}

// ExecuteWorkflowStatus handles the workflow status command execution.
func (wc *WorkflowCommand) ExecuteWorkflowStatus(cmd *cobra.Command, args []string) error {
	// Get flags
	workflowID, _ := cmd.Flags().GetString("workflow-id")
	listActive, _ := cmd.Flags().GetBool("active")
	listHistory, _ := cmd.Flags().GetBool("history")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Get workflow manager
	workflowManager := wc.cli.GetWorkflowManager()
	if workflowManager == nil {
		return wc.cli.CreateWorkflowError("Workflow manager not available", "status")
	}

	// Handle specific workflow status
	if workflowID != "" {
		status, err := workflowManager.GetWorkflowStatus(workflowID)
		if err != nil {
			return fmt.Errorf("🚫 %s %s",
				wc.cli.Error("Failed to get workflow status."),
				wc.cli.Info(err.Error()))
		}

		if outputFormat == "json" || outputFormat == "yaml" {
			return wc.cli.OutputMachineReadable(status, outputFormat)
		}

		wc.outputWorkflowStatus(status)
		return nil
	}

	// Handle active workflows list
	if listActive {
		workflows := workflowManager.GetActiveWorkflows()
		if outputFormat == "json" || outputFormat == "yaml" {
			return wc.cli.OutputMachineReadable(workflows, outputFormat)
		}

		wc.outputActiveWorkflows(workflows)
		return nil
	}

	// Handle workflow history
	if listHistory {
		workflows := workflowManager.GetWorkflowHistory()
		if outputFormat == "json" || outputFormat == "yaml" {
			return wc.cli.OutputMachineReadable(workflows, outputFormat)
		}

		wc.outputWorkflowHistory(workflows)
		return nil
	}

	// Default: show active workflows
	workflows := workflowManager.GetActiveWorkflows()
	if outputFormat == "json" || outputFormat == "yaml" {
		return wc.cli.OutputMachineReadable(workflows, outputFormat)
	}

	wc.outputActiveWorkflows(workflows)
	return nil
}

// Helper methods for output formatting

// createProgressCallback creates a progress callback for workflows.
func (wc *WorkflowCommand) createProgressCallback() func(*interfaces.WorkflowProgress) {
	return func(progress *interfaces.WorkflowProgress) {
		if progress == nil {
			return
		}

		// Format progress message
		progressBar := wc.formatProgressBar(progress.Progress)
		wc.cli.QuietOutput("📊 %s [%s] %.1f%% - %s",
			wc.cli.Info(progress.Stage),
			progressBar,
			progress.Progress,
			progress.Message)

		// Show warnings if any
		if len(progress.Warnings) > 0 {
			for _, warning := range progress.Warnings {
				wc.cli.WarningOutput("⚠️  %s", warning)
			}
		}

		// Show errors if any
		if len(progress.Errors) > 0 {
			for _, error := range progress.Errors {
				wc.cli.ErrorOutput("❌ %s", error)
			}
		}
	}
}

// formatProgressBar creates a visual progress bar.
func (wc *WorkflowCommand) formatProgressBar(progress float64) string {
	const barWidth = 20
	filled := int(progress * barWidth / 100)
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return bar
}

// outputProjectWorkflowResult outputs the result of a project workflow.
func (wc *WorkflowCommand) outputProjectWorkflowResult(result *interfaces.ProjectWorkflowResult) {
	if result.Success {
		wc.cli.SuccessOutput("✅ %s", wc.cli.Success("Project generation workflow completed successfully!"))
		wc.cli.QuietOutput("📁 Project path: %s", wc.cli.Highlight(result.ProjectPath))
		wc.cli.QuietOutput("📄 Generated files: %d", len(result.GeneratedFiles))
		wc.cli.QuietOutput("⏱️  Duration: %s", result.Duration)

		if result.ValidationResult != nil {
			wc.cli.QuietOutput("🔍 Validation: %s (%d issues)",
				wc.getValidationStatus(result.ValidationResult.Valid),
				len(result.ValidationResult.Issues))
		}

		if result.AuditResult != nil {
			wc.cli.QuietOutput("🔍 Audit score: %.1f", result.AuditResult.OverallScore)
		}

		if result.ReportPath != "" {
			wc.cli.QuietOutput("📊 Report: %s", result.ReportPath)
		}
	} else {
		wc.cli.ErrorOutput("❌ %s", wc.cli.Error("Project generation workflow failed"))
		for _, error := range result.Errors {
			wc.cli.ErrorOutput("   %s", error)
		}
	}

	// Show warnings if any
	for _, warning := range result.Warnings {
		wc.cli.WarningOutput("⚠️  %s", warning)
	}
}

// outputValidationWorkflowResult outputs the result of a validation workflow.
func (wc *WorkflowCommand) outputValidationWorkflowResult(result *interfaces.ValidationWorkflowResult) {
	if result.Success {
		wc.cli.SuccessOutput("✅ %s", wc.cli.Success("Validation workflow completed successfully!"))

		if result.ValidationResult != nil {
			wc.cli.QuietOutput("🔍 Validation: %s (%d issues)",
				wc.getValidationStatus(result.ValidationResult.Valid),
				len(result.ValidationResult.Issues))
		}

		if len(result.FixesApplied) > 0 {
			wc.cli.QuietOutput("🔧 Fixes applied: %d", len(result.FixesApplied))
		}

		if result.ReportPath != "" {
			wc.cli.QuietOutput("📊 Report: %s", result.ReportPath)
		}

		wc.cli.QuietOutput("⏱️  Duration: %s", result.Duration)
	} else {
		wc.cli.ErrorOutput("❌ %s", wc.cli.Error("Validation workflow failed"))
		for _, error := range result.Errors {
			wc.cli.ErrorOutput("   %s", error)
		}
	}
}

// outputAuditWorkflowResult outputs the result of an audit workflow.
func (wc *WorkflowCommand) outputAuditWorkflowResult(result *interfaces.AuditWorkflowResult) {
	if result.Success {
		wc.cli.SuccessOutput("✅ %s", wc.cli.Success("Audit workflow completed successfully!"))

		if result.AuditResult != nil {
			wc.cli.QuietOutput("🔍 Audit score: %.1f", result.AuditResult.OverallScore)
		}

		if result.ReportPath != "" {
			wc.cli.QuietOutput("📊 Report: %s", result.ReportPath)
		}

		wc.cli.QuietOutput("⏱️  Duration: %s", result.Duration)
	} else {
		wc.cli.ErrorOutput("❌ %s", wc.cli.Error("Audit workflow failed"))
		for _, error := range result.Errors {
			wc.cli.ErrorOutput("   %s", error)
		}
	}
}

// outputValidationAuditWorkflowResult outputs the result of a combined validation and audit workflow.
func (wc *WorkflowCommand) outputValidationAuditWorkflowResult(result *interfaces.ValidationAuditWorkflowResult) {
	if result.Success {
		wc.cli.SuccessOutput("✅ %s", wc.cli.Success("Validation and audit workflow completed successfully!"))

		if result.ValidationResult != nil {
			wc.cli.QuietOutput("🔍 Validation: %s (%d issues)",
				wc.getValidationStatus(result.ValidationResult.Valid),
				len(result.ValidationResult.Issues))
		}

		if result.AuditResult != nil {
			wc.cli.QuietOutput("🔍 Audit score: %.1f", result.AuditResult.OverallScore)
		}

		if len(result.FixesApplied) > 0 {
			wc.cli.QuietOutput("🔧 Fixes applied: %d", len(result.FixesApplied))
		}

		if result.ReportPath != "" {
			wc.cli.QuietOutput("📊 Report: %s", result.ReportPath)
		}

		wc.cli.QuietOutput("⏱️  Duration: %s", result.Duration)
	} else {
		wc.cli.ErrorOutput("❌ %s", wc.cli.Error("Validation and audit workflow failed"))
		for _, error := range result.Errors {
			wc.cli.ErrorOutput("   %s", error)
		}
	}
}

// outputWorkflowStatus outputs the status of a specific workflow.
func (wc *WorkflowCommand) outputWorkflowStatus(status *interfaces.WorkflowStatus) {
	wc.cli.QuietOutput("🔄 Workflow Status")
	wc.cli.QuietOutput("================")
	wc.cli.QuietOutput("ID: %s", status.ID)
	wc.cli.QuietOutput("Type: %s", status.Type)
	wc.cli.QuietOutput("Status: %s", wc.getStatusColor(status.Status))
	wc.cli.QuietOutput("Start Time: %s", status.StartTime.Format(time.RFC3339))

	if status.EndTime != nil {
		wc.cli.QuietOutput("End Time: %s", status.EndTime.Format(time.RFC3339))
	}

	wc.cli.QuietOutput("Duration: %s", status.Duration)

	if status.Error != "" {
		wc.cli.ErrorOutput("Error: %s", status.Error)
	}

	if status.Progress != nil {
		wc.cli.QuietOutput("Progress: %.1f%% - %s", status.Progress.Progress, status.Progress.Message)
	}
}

// outputActiveWorkflows outputs the list of active workflows.
func (wc *WorkflowCommand) outputActiveWorkflows(workflows []interfaces.WorkflowInfo) {
	if len(workflows) == 0 {
		wc.cli.QuietOutput("📋 No active workflows")
		return
	}

	wc.cli.QuietOutput("🔄 Active Workflows (%d)", len(workflows))
	wc.cli.QuietOutput("==================")

	for _, workflow := range workflows {
		wc.cli.QuietOutput("%s | %s | %s | %s",
			wc.cli.Dim(workflow.ID[:8]),
			workflow.Type,
			wc.getStatusColor(workflow.Status),
			workflow.StartTime.Format("15:04:05"))
	}
}

// outputWorkflowHistory outputs the workflow history.
func (wc *WorkflowCommand) outputWorkflowHistory(workflows []interfaces.WorkflowInfo) {
	if len(workflows) == 0 {
		wc.cli.QuietOutput("📋 No workflow history")
		return
	}

	wc.cli.QuietOutput("📚 Workflow History (%d)", len(workflows))
	wc.cli.QuietOutput("==================")

	for _, workflow := range workflows {
		successIcon := "❌"
		if workflow.Success {
			successIcon = "✅"
		}

		wc.cli.QuietOutput("%s %s | %s | %s | %s",
			successIcon,
			wc.cli.Dim(workflow.ID[:8]),
			workflow.Type,
			wc.getStatusColor(workflow.Status),
			workflow.StartTime.Format("2006-01-02 15:04:05"))
	}
}

// Helper methods for formatting

// getValidationStatus returns a colored validation status.
func (wc *WorkflowCommand) getValidationStatus(valid bool) string {
	if valid {
		return wc.cli.Success("PASSED")
	}
	return wc.cli.Error("FAILED")
}

// getStatusColor returns a colored status string.
func (wc *WorkflowCommand) getStatusColor(status string) string {
	switch status {
	case interfaces.WorkflowStatusCompleted:
		return wc.cli.Success(status)
	case interfaces.WorkflowStatusFailed:
		return wc.cli.Error(status)
	case interfaces.WorkflowStatusCancelled:
		return wc.cli.Warning(status)
	case interfaces.WorkflowStatusRunning:
		return wc.cli.Info(status)
	default:
		return wc.cli.Dim(status)
	}
}
