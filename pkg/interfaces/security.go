// Package interfaces defines security-related interfaces for the Open Source Project Generator.
package interfaces

import (
	"os"
	"time"
)

// SecurityManager provides comprehensive security management capabilities
type SecurityManager interface {
	// Input sanitization
	SanitizeString(input string, fieldName string) *SanitizationResult
	SanitizeProjectName(name string) *SanitizationResult
	SanitizeFilePath(path string, fieldName string) *SanitizationResult
	SanitizeURL(urlStr string, fieldName string) *SanitizationResult
	SanitizeEmail(email string, fieldName string) *SanitizationResult
	ValidateAndSanitizeMap(data map[string]interface{}, prefix string) map[string]*SanitizationResult

	// Template security
	ValidateTemplateContent(content string, filePath string) *TemplateSecurityResult
	ValidateTemplateFile(filePath string) (*TemplateSecurityResult, error)
	ScanTemplateDirectory(dirPath string) (map[string]*TemplateSecurityResult, error)

	// File operations security
	SecureReadFile(path string) (*FileOperationResult, []byte, error)
	SecureWriteFile(path string, data []byte, perm os.FileMode) (*FileOperationResult, error)
	SecureCopyFile(srcPath, dstPath string) (*FileOperationResult, error)
	SecureCreateDirectory(path string, perm os.FileMode) (*FileOperationResult, error)
	ValidateFilePath(path string, operation string) error
	GetFilePermissions(path string) (map[string]interface{}, error)

	// Backup management
	BackupFile(filePath string) (*BackupResult, error)
	BackupDirectory(dirPath string) (map[string]*BackupResult, error)
	RestoreFile(originalPath string, backupTimestamp time.Time) (*BackupResult, error)
	ListBackups(originalPath string) ([]BackupInfo, error)
	SetBackupEnabled(enabled bool)
	IsBackupEnabled() bool

	// Dry-run operations
	SetDryRunMode(enabled bool)
	IsDryRunMode() bool
	RecordFileWrite(path string, data []byte, overwrite bool)
	RecordFileDelete(path string)
	RecordDirectoryCreate(path string)
	RecordDirectoryDelete(path string)
	RecordFileCopy(srcPath, dstPath string, overwrite bool)
	RecordTemplateProcess(templatePath, outputPath string, variables map[string]interface{})
	GetDryRunOperations() []DryRunOperation
	GetDryRunSummary() map[string]interface{}
	GenerateDryRunReport(format string) (string, error)

	// User confirmation
	SetNonInteractive(nonInteractive bool)
	IsNonInteractive() bool
	ConfirmFileOverwrite(filePath string, fileSize int64) (*ConfirmationResult, error)
	ConfirmDirectoryDelete(dirPath string, fileCount int, totalSize int64) (*ConfirmationResult, error)
	ConfirmBulkOperation(operationType string, itemCount int, details []string) (*ConfirmationResult, error)
	ConfirmSecurityRisk(riskDescription string, riskLevel string, details []string) (*ConfirmationResult, error)
	ConfirmWithDryRun(dryRunSummary map[string]interface{}) (*ConfirmationResult, error)

	// Security configuration
	SetSecurityConfig(config map[string]interface{}) error
	GetSecurityConfig() map[string]interface{}
}

// SanitizationResult contains the result of input sanitization
type SanitizationResult struct {
	Original    string   `json:"original"`
	Sanitized   string   `json:"sanitized"`
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	WasModified bool     `json:"was_modified"`
}

// TemplateSecurityResult contains the result of template security validation
type TemplateSecurityResult struct {
	IsSecure        bool     `json:"is_secure"`
	SecurityIssues  []string `json:"security_issues"`
	Warnings        []string `json:"warnings"`
	BlockedPatterns []string `json:"blocked_patterns"`
	FilePath        string   `json:"file_path"`
	FileSize        int64    `json:"file_size"`
}

// FileOperationResult contains the result of a secure file operation
type FileOperationResult struct {
	Success      bool      `json:"success"`
	FilePath     string    `json:"file_path"`
	Operation    string    `json:"operation"`
	BytesWritten int64     `json:"bytes_written,omitempty"`
	BytesRead    int64     `json:"bytes_read,omitempty"`
	Checksum     string    `json:"checksum,omitempty"`
	Permissions  string    `json:"permissions,omitempty"`
	Error        string    `json:"error,omitempty"`
	Warnings     []string  `json:"warnings,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// BackupResult contains information about a backup operation
type BackupResult struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	Success      bool      `json:"success"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	FileSize     int64     `json:"file_size"`
	Checksum     string    `json:"checksum,omitempty"`
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	Timestamp    time.Time `json:"timestamp"`
	Size         int64     `json:"size"`
}

// DryRunOperation represents a planned operation in dry-run mode
type DryRunOperation struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Path        string                 `json:"path"`
	Details     map[string]interface{} `json:"details"`
	Timestamp   time.Time              `json:"timestamp"`
	Impact      string                 `json:"impact"` // "safe", "warning", "destructive"
	Size        int64                  `json:"size,omitempty"`
}

// ConfirmationResult represents the result of a confirmation request
type ConfirmationResult struct {
	Confirmed      bool      `json:"confirmed"`
	UserInput      string    `json:"user_input"`
	Timestamp      time.Time `json:"timestamp"`
	NonInteractive bool      `json:"non_interactive"`
	TimedOut       bool      `json:"timed_out"`
	DefaultUsed    bool      `json:"default_used"`
}
