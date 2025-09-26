// Package security provides comprehensive security management for the Open Source Project Generator.
package security

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Manager implements the SecurityManager interface
type Manager struct {
	inputSanitizer      *InputSanitizer
	templateSecurity    *TemplateSecurityManager
	fileOperations      *SecureFileOperations
	backupManager       *BackupManager
	dryRunManager       *DryRunManager
	confirmationManager *ConfirmationManager
	confirmationHistory *ConfirmationHistory
	config              map[string]interface{}
}

// NewSecurityManager creates a new security manager with all components
func NewSecurityManager(workspaceDir string) interfaces.SecurityManager {
	// Create backup directory
	backupDir := filepath.Join(workspaceDir, ".generator", "backups")

	// Initialize all security components with proper allowed paths
	allowedPaths := []string{workspaceDir, backupDir}
	manager := &Manager{
		inputSanitizer:      NewInputSanitizer(),
		templateSecurity:    NewTemplateSecurityManager(),
		fileOperations:      NewSecureFileOperations(allowedPaths),
		backupManager:       NewBackupManager(backupDir),
		dryRunManager:       NewDryRunManager(),
		confirmationManager: NewConfirmationManager(),
		confirmationHistory: NewConfirmationHistory(),
		config:              make(map[string]interface{}),
	}

	// Set default configuration
	manager.config["workspace_dir"] = workspaceDir
	manager.config["backup_dir"] = backupDir
	manager.config["backup_enabled"] = true
	manager.config["dry_run_enabled"] = false
	manager.config["non_interactive"] = false
	manager.config["max_file_size"] = int64(100 * 1024 * 1024) // 100MB
	manager.config["max_backups"] = 10

	return manager
}

// Input sanitization methods

func (m *Manager) SanitizeString(input string, fieldName string) *interfaces.SanitizationResult {
	result := m.inputSanitizer.SanitizeString(input, fieldName)
	return &interfaces.SanitizationResult{
		Original:    result.Original,
		Sanitized:   result.Sanitized,
		IsValid:     result.IsValid,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		WasModified: result.WasModified,
	}
}

func (m *Manager) SanitizeProjectName(name string) *interfaces.SanitizationResult {
	result := m.inputSanitizer.SanitizeProjectName(name)
	return &interfaces.SanitizationResult{
		Original:    result.Original,
		Sanitized:   result.Sanitized,
		IsValid:     result.IsValid,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		WasModified: result.WasModified,
	}
}

func (m *Manager) SanitizeFilePath(path string, fieldName string) *interfaces.SanitizationResult {
	result := m.inputSanitizer.SanitizeFilePath(path, fieldName)
	return &interfaces.SanitizationResult{
		Original:    result.Original,
		Sanitized:   result.Sanitized,
		IsValid:     result.IsValid,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		WasModified: result.WasModified,
	}
}

func (m *Manager) SanitizeURL(urlStr string, fieldName string) *interfaces.SanitizationResult {
	result := m.inputSanitizer.SanitizeURL(urlStr, fieldName)
	return &interfaces.SanitizationResult{
		Original:    result.Original,
		Sanitized:   result.Sanitized,
		IsValid:     result.IsValid,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		WasModified: result.WasModified,
	}
}

func (m *Manager) SanitizeEmail(email string, fieldName string) *interfaces.SanitizationResult {
	result := m.inputSanitizer.SanitizeEmail(email, fieldName)
	return &interfaces.SanitizationResult{
		Original:    result.Original,
		Sanitized:   result.Sanitized,
		IsValid:     result.IsValid,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		WasModified: result.WasModified,
	}
}

func (m *Manager) ValidateAndSanitizeMap(data map[string]interface{}, prefix string) map[string]*interfaces.SanitizationResult {
	results := m.inputSanitizer.ValidateAndSanitizeMap(data, prefix)
	interfaceResults := make(map[string]*interfaces.SanitizationResult)

	for key, result := range results {
		interfaceResults[key] = &interfaces.SanitizationResult{
			Original:    result.Original,
			Sanitized:   result.Sanitized,
			IsValid:     result.IsValid,
			Errors:      result.Errors,
			Warnings:    result.Warnings,
			WasModified: result.WasModified,
		}
	}

	return interfaceResults
}

// Template security methods

func (m *Manager) ValidateTemplateContent(content string, filePath string) *interfaces.TemplateSecurityResult {
	result := m.templateSecurity.ValidateTemplateContent(content, filePath)
	return &interfaces.TemplateSecurityResult{
		IsSecure:        result.IsSecure,
		SecurityIssues:  result.SecurityIssues,
		Warnings:        result.Warnings,
		BlockedPatterns: result.BlockedPatterns,
		FilePath:        result.FilePath,
		FileSize:        result.FileSize,
	}
}

func (m *Manager) ValidateTemplateFile(filePath string) (*interfaces.TemplateSecurityResult, error) {
	result, err := m.templateSecurity.ValidateTemplateFile(filePath)
	if err != nil {
		return nil, err
	}

	return &interfaces.TemplateSecurityResult{
		IsSecure:        result.IsSecure,
		SecurityIssues:  result.SecurityIssues,
		Warnings:        result.Warnings,
		BlockedPatterns: result.BlockedPatterns,
		FilePath:        result.FilePath,
		FileSize:        result.FileSize,
	}, nil
}

func (m *Manager) ScanTemplateDirectory(dirPath string) (map[string]*interfaces.TemplateSecurityResult, error) {
	results, err := m.templateSecurity.ScanTemplateDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	interfaceResults := make(map[string]*interfaces.TemplateSecurityResult)
	for path, result := range results {
		interfaceResults[path] = &interfaces.TemplateSecurityResult{
			IsSecure:        result.IsSecure,
			SecurityIssues:  result.SecurityIssues,
			Warnings:        result.Warnings,
			BlockedPatterns: result.BlockedPatterns,
			FilePath:        result.FilePath,
			FileSize:        result.FileSize,
		}
	}

	return interfaceResults, nil
}

// File operations security methods

func (m *Manager) SecureReadFile(path string) (*interfaces.FileOperationResult, []byte, error) {
	result, data, err := m.fileOperations.SecureReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	interfaceResult := &interfaces.FileOperationResult{
		Success:      result.Success,
		FilePath:     result.FilePath,
		Operation:    result.Operation,
		BytesWritten: result.BytesWritten,
		BytesRead:    result.BytesRead,
		Checksum:     result.Checksum,
		Permissions:  result.Permissions,
		Error:        result.Error,
		Warnings:     result.Warnings,
		Timestamp:    result.Timestamp,
	}

	return interfaceResult, data, nil
}

func (m *Manager) SecureWriteFile(path string, data []byte, perm os.FileMode) (*interfaces.FileOperationResult, error) {
	// Record dry-run operation if enabled
	if m.dryRunManager.IsEnabled() {
		m.dryRunManager.RecordFileWrite(path, data, true)
		return &interfaces.FileOperationResult{
			Success:   true,
			FilePath:  path,
			Operation: "write (dry-run)",
			Timestamp: time.Now(),
		}, nil
	}

	// Create backup if file exists and backup is enabled
	if m.backupManager.IsEnabled() {
		if _, err := os.Stat(path); err == nil {
			if _, backupErr := m.backupManager.BackupFile(path); backupErr != nil {
				fmt.Printf("Warning: failed to create backup for %s: %v\n", path, backupErr)
			}
		}
	}

	result, err := m.fileOperations.SecureWriteFile(path, data, perm)
	if err != nil {
		return nil, err
	}

	return &interfaces.FileOperationResult{
		Success:      result.Success,
		FilePath:     result.FilePath,
		Operation:    result.Operation,
		BytesWritten: result.BytesWritten,
		BytesRead:    result.BytesRead,
		Checksum:     result.Checksum,
		Permissions:  result.Permissions,
		Error:        result.Error,
		Warnings:     result.Warnings,
		Timestamp:    result.Timestamp,
	}, nil
}

func (m *Manager) SecureCopyFile(srcPath, dstPath string) (*interfaces.FileOperationResult, error) {
	// Record dry-run operation if enabled
	if m.dryRunManager.IsEnabled() {
		m.dryRunManager.RecordFileCopy(srcPath, dstPath, true)
		return &interfaces.FileOperationResult{
			Success:   true,
			FilePath:  fmt.Sprintf("%s -> %s", srcPath, dstPath),
			Operation: "copy (dry-run)",
			Timestamp: time.Now(),
		}, nil
	}

	// Create backup if destination exists and backup is enabled
	if m.backupManager.IsEnabled() {
		if _, err := os.Stat(dstPath); err == nil {
			if _, backupErr := m.backupManager.BackupFile(dstPath); backupErr != nil {
				fmt.Printf("Warning: failed to create backup for %s: %v\n", dstPath, backupErr)
			}
		}
	}

	result, err := m.fileOperations.SecureCopyFile(srcPath, dstPath)
	if err != nil {
		return nil, err
	}

	return &interfaces.FileOperationResult{
		Success:      result.Success,
		FilePath:     result.FilePath,
		Operation:    result.Operation,
		BytesWritten: result.BytesWritten,
		BytesRead:    result.BytesRead,
		Checksum:     result.Checksum,
		Permissions:  result.Permissions,
		Error:        result.Error,
		Warnings:     result.Warnings,
		Timestamp:    result.Timestamp,
	}, nil
}

func (m *Manager) SecureCreateDirectory(path string, perm os.FileMode) (*interfaces.FileOperationResult, error) {
	// Record dry-run operation if enabled
	if m.dryRunManager.IsEnabled() {
		m.dryRunManager.RecordDirectoryCreate(path)
		return &interfaces.FileOperationResult{
			Success:   true,
			FilePath:  path,
			Operation: "mkdir (dry-run)",
			Timestamp: time.Now(),
		}, nil
	}

	result, err := m.fileOperations.SecureCreateDirectory(path, perm)
	if err != nil {
		return nil, err
	}

	return &interfaces.FileOperationResult{
		Success:      result.Success,
		FilePath:     result.FilePath,
		Operation:    result.Operation,
		BytesWritten: result.BytesWritten,
		BytesRead:    result.BytesRead,
		Checksum:     result.Checksum,
		Permissions:  result.Permissions,
		Error:        result.Error,
		Warnings:     result.Warnings,
		Timestamp:    result.Timestamp,
	}, nil
}

func (m *Manager) ValidateFilePath(path string, operation string) error {
	return m.fileOperations.ValidateFilePath(path, operation)
}

func (m *Manager) GetFilePermissions(path string) (map[string]interface{}, error) {
	return m.fileOperations.GetFilePermissions(path)
}

// Backup management methods

func (m *Manager) BackupFile(filePath string) (*interfaces.BackupResult, error) {
	result, err := m.backupManager.BackupFile(filePath)
	if err != nil {
		return nil, err
	}

	return &interfaces.BackupResult{
		OriginalPath: result.OriginalPath,
		BackupPath:   result.BackupPath,
		Success:      result.Success,
		Error:        result.Error,
		Timestamp:    result.Timestamp,
		FileSize:     result.FileSize,
		Checksum:     result.Checksum,
	}, nil
}

func (m *Manager) BackupDirectory(dirPath string) (map[string]*interfaces.BackupResult, error) {
	results, err := m.backupManager.BackupDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	interfaceResults := make(map[string]*interfaces.BackupResult)
	for path, result := range results {
		interfaceResults[path] = &interfaces.BackupResult{
			OriginalPath: result.OriginalPath,
			BackupPath:   result.BackupPath,
			Success:      result.Success,
			Error:        result.Error,
			Timestamp:    result.Timestamp,
			FileSize:     result.FileSize,
			Checksum:     result.Checksum,
		}
	}

	return interfaceResults, nil
}

func (m *Manager) RestoreFile(originalPath string, backupTimestamp time.Time) (*interfaces.BackupResult, error) {
	result, err := m.backupManager.RestoreFile(originalPath, backupTimestamp)
	if err != nil {
		return nil, err
	}

	return &interfaces.BackupResult{
		OriginalPath: result.OriginalPath,
		BackupPath:   result.BackupPath,
		Success:      result.Success,
		Error:        result.Error,
		Timestamp:    result.Timestamp,
		FileSize:     result.FileSize,
		Checksum:     result.Checksum,
	}, nil
}

func (m *Manager) ListBackups(originalPath string) ([]interfaces.BackupInfo, error) {
	backups, err := m.backupManager.ListBackups(originalPath)
	if err != nil {
		return nil, err
	}

	interfaceBackups := make([]interfaces.BackupInfo, len(backups))
	for i, backup := range backups {
		interfaceBackups[i] = interfaces.BackupInfo{
			OriginalPath: backup.OriginalPath,
			BackupPath:   backup.BackupPath,
			Timestamp:    backup.Timestamp,
			Size:         backup.Size,
		}
	}

	return interfaceBackups, nil
}

func (m *Manager) SetBackupEnabled(enabled bool) {
	m.backupManager.SetEnabled(enabled)
	m.config["backup_enabled"] = enabled
}

func (m *Manager) IsBackupEnabled() bool {
	return m.backupManager.IsEnabled()
}

// Dry-run methods

func (m *Manager) SetDryRunMode(enabled bool) {
	m.dryRunManager.SetEnabled(enabled)
	m.config["dry_run_enabled"] = enabled
}

func (m *Manager) IsDryRunMode() bool {
	return m.dryRunManager.IsEnabled()
}

func (m *Manager) RecordFileWrite(path string, data []byte, overwrite bool) {
	m.dryRunManager.RecordFileWrite(path, data, overwrite)
}

func (m *Manager) RecordFileDelete(path string) {
	m.dryRunManager.RecordFileDelete(path)
}

func (m *Manager) RecordDirectoryCreate(path string) {
	m.dryRunManager.RecordDirectoryCreate(path)
}

func (m *Manager) RecordDirectoryDelete(path string) {
	m.dryRunManager.RecordDirectoryDelete(path)
}

func (m *Manager) RecordFileCopy(srcPath, dstPath string, overwrite bool) {
	m.dryRunManager.RecordFileCopy(srcPath, dstPath, overwrite)
}

func (m *Manager) RecordTemplateProcess(templatePath, outputPath string, variables map[string]interface{}) {
	m.dryRunManager.RecordTemplateProcess(templatePath, outputPath, variables)
}

func (m *Manager) GetDryRunOperations() []interfaces.DryRunOperation {
	operations := m.dryRunManager.GetOperations()
	interfaceOperations := make([]interfaces.DryRunOperation, len(operations))

	for i, op := range operations {
		interfaceOperations[i] = interfaces.DryRunOperation{
			Type:        op.Type,
			Description: op.Description,
			Path:        op.Path,
			Details:     op.Details,
			Timestamp:   op.Timestamp,
			Impact:      op.Impact,
			Size:        op.Size,
		}
	}

	return interfaceOperations
}

func (m *Manager) GetDryRunSummary() map[string]interface{} {
	return m.dryRunManager.GetSummary()
}

func (m *Manager) GenerateDryRunReport(format string) (string, error) {
	return m.dryRunManager.GenerateReport(format)
}

// User confirmation methods

func (m *Manager) SetNonInteractive(nonInteractive bool) {
	m.confirmationManager.SetNonInteractive(nonInteractive)
	m.config["non_interactive"] = nonInteractive
}

func (m *Manager) IsNonInteractive() bool {
	return m.confirmationManager.IsNonInteractive()
}

func (m *Manager) ConfirmFileOverwrite(filePath string, fileSize int64) (*interfaces.ConfirmationResult, error) {
	result, err := m.confirmationManager.ConfirmFileOverwrite(filePath, fileSize)
	if err != nil {
		return nil, err
	}

	interfaceResult := &interfaces.ConfirmationResult{
		Confirmed:      result.Confirmed,
		UserInput:      result.UserInput,
		Timestamp:      result.Timestamp,
		NonInteractive: result.NonInteractive,
		TimedOut:       result.TimedOut,
		DefaultUsed:    result.DefaultUsed,
	}

	return interfaceResult, nil
}

func (m *Manager) ConfirmDirectoryDelete(dirPath string, fileCount int, totalSize int64) (*interfaces.ConfirmationResult, error) {
	result, err := m.confirmationManager.ConfirmDirectoryDelete(dirPath, fileCount, totalSize)
	if err != nil {
		return nil, err
	}

	return &interfaces.ConfirmationResult{
		Confirmed:      result.Confirmed,
		UserInput:      result.UserInput,
		Timestamp:      result.Timestamp,
		NonInteractive: result.NonInteractive,
		TimedOut:       result.TimedOut,
		DefaultUsed:    result.DefaultUsed,
	}, nil
}

func (m *Manager) ConfirmBulkOperation(operationType string, itemCount int, details []string) (*interfaces.ConfirmationResult, error) {
	result, err := m.confirmationManager.ConfirmBulkOperation(operationType, itemCount, details)
	if err != nil {
		return nil, err
	}

	return &interfaces.ConfirmationResult{
		Confirmed:      result.Confirmed,
		UserInput:      result.UserInput,
		Timestamp:      result.Timestamp,
		NonInteractive: result.NonInteractive,
		TimedOut:       result.TimedOut,
		DefaultUsed:    result.DefaultUsed,
	}, nil
}

func (m *Manager) ConfirmSecurityRisk(riskDescription string, riskLevel string, details []string) (*interfaces.ConfirmationResult, error) {
	result, err := m.confirmationManager.ConfirmSecurityRisk(riskDescription, riskLevel, details)
	if err != nil {
		return nil, err
	}

	return &interfaces.ConfirmationResult{
		Confirmed:      result.Confirmed,
		UserInput:      result.UserInput,
		Timestamp:      result.Timestamp,
		NonInteractive: result.NonInteractive,
		TimedOut:       result.TimedOut,
		DefaultUsed:    result.DefaultUsed,
	}, nil
}

func (m *Manager) ConfirmWithDryRun(dryRunSummary map[string]interface{}) (*interfaces.ConfirmationResult, error) {
	result, err := m.confirmationManager.ConfirmWithDryRun(dryRunSummary)
	if err != nil {
		return nil, err
	}

	return &interfaces.ConfirmationResult{
		Confirmed:      result.Confirmed,
		UserInput:      result.UserInput,
		Timestamp:      result.Timestamp,
		NonInteractive: result.NonInteractive,
		TimedOut:       result.TimedOut,
		DefaultUsed:    result.DefaultUsed,
	}, nil
}

// Configuration methods

func (m *Manager) SetSecurityConfig(config map[string]interface{}) error {
	// Update internal configuration
	for key, value := range config {
		m.config[key] = value
	}

	// Apply configuration to components
	if err := m.fileOperations.SetSecurityConfig(config); err != nil {
		return fmt.Errorf("failed to configure file operations: %w", err)
	}

	// Configure backup manager
	if maxBackups, ok := config["max_backups"].(int); ok {
		m.backupManager.SetMaxBackups(maxBackups)
	}

	if backupEnabled, ok := config["backup_enabled"].(bool); ok {
		m.backupManager.SetEnabled(backupEnabled)
	}

	// Configure confirmation manager
	if nonInteractive, ok := config["non_interactive"].(bool); ok {
		m.confirmationManager.SetNonInteractive(nonInteractive)
	}

	if defaultAnswer, ok := config["default_answer"].(bool); ok {
		m.confirmationManager.SetDefaultAnswer(defaultAnswer)
	}

	// Configure dry-run manager
	if dryRunEnabled, ok := config["dry_run_enabled"].(bool); ok {
		m.dryRunManager.SetEnabled(dryRunEnabled)
	}

	return nil
}

func (m *Manager) GetSecurityConfig() map[string]interface{} {
	// Return a copy of the configuration
	config := make(map[string]interface{})
	for key, value := range m.config {
		config[key] = value
	}
	return config
}
