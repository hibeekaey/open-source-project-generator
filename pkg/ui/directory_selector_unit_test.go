package ui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DirectoryMockInteractiveUI is a mock implementation for testing directory selector
type DirectoryMockInteractiveUI struct {
	textResponses    map[string]*interfaces.TextResult
	menuResponses    map[string]*interfaces.MenuResult
	confirmResponses map[string]*interfaces.ConfirmResult
	callCount        map[string]int
}

func NewDirectoryMockInteractiveUI() *DirectoryMockInteractiveUI {
	return &DirectoryMockInteractiveUI{
		textResponses:    make(map[string]*interfaces.TextResult),
		menuResponses:    make(map[string]*interfaces.MenuResult),
		confirmResponses: make(map[string]*interfaces.ConfirmResult),
		callCount:        make(map[string]int),
	}
}

func (m *DirectoryMockInteractiveUI) PromptText(ctx context.Context, config interfaces.TextPromptConfig) (*interfaces.TextResult, error) {
	m.callCount["PromptText"]++
	if response, exists := m.textResponses[config.Prompt]; exists {
		return response, nil
	}
	return &interfaces.TextResult{Value: config.DefaultValue, Action: "submit"}, nil
}

func (m *DirectoryMockInteractiveUI) ShowMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	m.callCount["ShowMenu"]++
	if response, exists := m.menuResponses[config.Title]; exists {
		return response, nil
	}
	return &interfaces.MenuResult{SelectedIndex: 0, SelectedValue: "cancel", Action: "select"}, nil
}

func (m *DirectoryMockInteractiveUI) PromptConfirm(ctx context.Context, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, error) {
	m.callCount["PromptConfirm"]++
	if response, exists := m.confirmResponses[config.Prompt]; exists {
		return response, nil
	}
	return &interfaces.ConfirmResult{Confirmed: false, Action: "confirm"}, nil
}

func (m *DirectoryMockInteractiveUI) ShowTable(ctx context.Context, config interfaces.TableConfig) error {
	m.callCount["ShowTable"]++
	return nil
}

// Implement other required methods with no-op implementations
func (m *DirectoryMockInteractiveUI) ShowMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) ShowCheckboxList(ctx context.Context, config interfaces.CheckboxConfig) (*interfaces.CheckboxResult, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) PromptSelect(ctx context.Context, config interfaces.SelectConfig) (*interfaces.SelectResult, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) ShowTree(ctx context.Context, config interfaces.TreeConfig) error {
	return nil
}
func (m *DirectoryMockInteractiveUI) ShowProgress(ctx context.Context, config interfaces.ProgressConfig) (interfaces.ProgressTracker, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) ShowBreadcrumb(ctx context.Context, path []string) error {
	return nil
}
func (m *DirectoryMockInteractiveUI) ShowHelp(ctx context.Context, helpContext string) error {
	return nil
}
func (m *DirectoryMockInteractiveUI) ShowError(ctx context.Context, config interfaces.ErrorConfig) (*interfaces.ErrorResult, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) StartSession(ctx context.Context, config interfaces.SessionConfig) (*interfaces.UISession, error) {
	return nil, nil
}
func (m *DirectoryMockInteractiveUI) EndSession(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (m *DirectoryMockInteractiveUI) SaveSessionState(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (m *DirectoryMockInteractiveUI) RestoreSessionState(ctx context.Context, sessionID string) (*interfaces.UISession, error) {
	return nil, nil
}

// DirectoryMockLogger is a mock implementation for testing directory selector
type DirectoryMockLogger struct{}

func (m *DirectoryMockLogger) Debug(format string, args ...interface{})                      {}
func (m *DirectoryMockLogger) Info(format string, args ...interface{})                       {}
func (m *DirectoryMockLogger) Warn(format string, args ...interface{})                       {}
func (m *DirectoryMockLogger) Error(format string, args ...interface{})                      {}
func (m *DirectoryMockLogger) Fatal(format string, args ...interface{})                      {}
func (m *DirectoryMockLogger) DebugWithFields(message string, fields map[string]interface{}) {}
func (m *DirectoryMockLogger) InfoWithFields(message string, fields map[string]interface{})  {}
func (m *DirectoryMockLogger) WarnWithFields(message string, fields map[string]interface{})  {}
func (m *DirectoryMockLogger) ErrorWithFields(message string, fields map[string]interface{}) {}
func (m *DirectoryMockLogger) FatalWithFields(message string, fields map[string]interface{}) {}
func (m *DirectoryMockLogger) SetLevel(level int)                                            {}
func (m *DirectoryMockLogger) SetJSONOutput(enabled bool)                                    {}
func (m *DirectoryMockLogger) SetCallerInfo(enabled bool)                                    {}
func (m *DirectoryMockLogger) IsDebugEnabled() bool                                          { return false }
func (m *DirectoryMockLogger) IsInfoEnabled() bool                                           { return false }
func (m *DirectoryMockLogger) StartOperation(operation string, metadata map[string]interface{}) *interfaces.OperationContext {
	return nil
}
func (m *DirectoryMockLogger) FinishOperation(ctx *interfaces.OperationContext, metadata map[string]interface{}) {
}
func (m *DirectoryMockLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, metadata map[string]interface{}) {
}
func (m *DirectoryMockLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
}
func (m *DirectoryMockLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {}
func (m *DirectoryMockLogger) GetLevel() int                                                       { return 0 }
func (m *DirectoryMockLogger) LogOperationStart(operation string, fields map[string]interface{})   {}
func (m *DirectoryMockLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
}
func (m *DirectoryMockLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {
}
func (m *DirectoryMockLogger) LogMemoryUsage(operation string)                  {}
func (m *DirectoryMockLogger) WithComponent(component string) interfaces.Logger { return m }
func (m *DirectoryMockLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return nil
}
func (m *DirectoryMockLogger) GetLogDir() string                                { return "" }
func (m *DirectoryMockLogger) GetRecentEntries(limit int) []interfaces.LogEntry { return nil }
func (m *DirectoryMockLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (m *DirectoryMockLogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (m *DirectoryMockLogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }
func (m *DirectoryMockLogger) Close() error                                { return nil }

func TestDirectorySelector_SelectOutputDirectory_NewDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "new-project")

	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}

	// Set up mock response for directory path input
	mockUI.textResponses["Output Directory"] = &interfaces.TextResult{
		Value:  testPath,
		Action: "submit",
	}

	selector := NewDirectorySelector(mockUI, mockLogger)
	ctx := context.Background()

	result, err := selector.SelectOutputDirectory(ctx, "output/generated")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Cancelled {
		t.Fatal("Expected result not to be cancelled")
	}

	if result.Path != testPath {
		t.Errorf("Expected path %s, got %s", testPath, result.Path)
	}

	if result.Exists {
		t.Error("Expected directory not to exist")
	}

	if !result.RequiresCreation {
		t.Error("Expected directory to require creation")
	}
}

func TestDirectorySelector_SelectOutputDirectory_ExistingEmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "existing-empty")

	// Create the directory
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}

	// Set up mock response for directory path input
	mockUI.textResponses["Output Directory"] = &interfaces.TextResult{
		Value:  testPath,
		Action: "submit",
	}

	selector := NewDirectorySelector(mockUI, mockLogger)
	ctx := context.Background()

	result, err := selector.SelectOutputDirectory(ctx, "output/generated")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Cancelled {
		t.Fatal("Expected result not to be cancelled")
	}

	if result.Path != testPath {
		t.Errorf("Expected path %s, got %s", testPath, result.Path)
	}

	if !result.Exists {
		t.Error("Expected directory to exist")
	}

	if result.RequiresCreation {
		t.Error("Expected directory not to require creation")
	}
}

func TestDirectorySelector_SelectOutputDirectory_ExistingNonEmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "existing-nonempty")

	// Create the directory and add a file
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testFile := filepath.Join(testPath, "existing-file.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}

	// Set up mock responses
	mockUI.textResponses["Output Directory"] = &interfaces.TextResult{
		Value:  testPath,
		Action: "submit",
	}

	mockUI.menuResponses["Directory Conflict"] = &interfaces.MenuResult{
		SelectedIndex: 3, // Cancel option
		SelectedValue: "cancel",
		Action:        "select",
	}

	selector := NewDirectorySelector(mockUI, mockLogger)
	ctx := context.Background()

	result, err := selector.SelectOutputDirectory(ctx, "output/generated")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Cancelled {
		t.Error("Expected result to be cancelled")
	}

	// Verify that ShowTable was called to display directory contents
	if mockUI.callCount["ShowTable"] == 0 {
		t.Error("Expected ShowTable to be called to display directory contents")
	}

	// Verify that ShowMenu was called for conflict resolution
	if mockUI.callCount["ShowMenu"] == 0 {
		t.Error("Expected ShowMenu to be called for conflict resolution")
	}
}

func TestDirectorySelector_SelectOutputDirectory_OverwriteWithBackup(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "existing-overwrite")

	// Create the directory and add a file
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testFile := filepath.Join(testPath, "existing-file.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}

	// Set up mock responses
	mockUI.textResponses["Output Directory"] = &interfaces.TextResult{
		Value:  testPath,
		Action: "submit",
	}

	mockUI.menuResponses["Directory Conflict"] = &interfaces.MenuResult{
		SelectedIndex: 0, // Overwrite option
		SelectedValue: "overwrite",
		Action:        "select",
	}

	mockUI.confirmResponses["Confirm Overwrite with Backup"] = &interfaces.ConfirmResult{
		Confirmed: true,
		Action:    "confirm",
	}

	selector := NewDirectorySelector(mockUI, mockLogger)
	ctx := context.Background()

	result, err := selector.SelectOutputDirectory(ctx, "output/generated")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Cancelled {
		t.Error("Expected result not to be cancelled")
	}

	if result.ConflictResolution != "overwrite" {
		t.Errorf("Expected conflict resolution 'overwrite', got %s", result.ConflictResolution)
	}

	if result.BackupPath == "" {
		t.Error("Expected backup path to be set")
	}

	// Verify that PromptConfirm was called for overwrite confirmation
	if mockUI.callCount["PromptConfirm"] == 0 {
		t.Error("Expected PromptConfirm to be called for overwrite confirmation")
	}
}

func TestDirectoryValidator_ValidateDirectoryPath(t *testing.T) {
	validator := &DirectoryValidator{}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid relative path",
			path:        "output/my-project",
			expectError: false,
		},
		{
			name:        "Valid absolute path",
			path:        "/home/user/projects/my-app",
			expectError: false,
		},
		{
			name:        "Empty path",
			path:        "",
			expectError: true,
			errorCode:   "required",
		},
		{
			name:        "Path with invalid character",
			path:        "output/my<project",
			expectError: true,
			errorCode:   "invalid_character",
		},
		{
			name:        "Path too long",
			path:        string(make([]byte, 501)), // 501 characters
			expectError: true,
			errorCode:   "max_length",
		},
		{
			name:        "Path with reserved name",
			path:        "output/CON/project",
			expectError: true,
			errorCode:   "reserved_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDirectoryPath(tt.path)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}

				if validationErr, ok := err.(*interfaces.ValidationError); ok {
					if validationErr.Code != tt.errorCode {
						t.Errorf("Expected error code %s, got %s", tt.errorCode, validationErr.Code)
					}
				} else {
					t.Errorf("Expected ValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestDirectoryValidator_ValidateParentDirectory(t *testing.T) {
	validator := &DirectoryValidator{}

	// Test with existing directory (temp dir)
	tempDir := t.TempDir()
	err := validator.ValidateParentDirectory(tempDir)
	if err != nil {
		t.Errorf("Expected no error for existing directory, got: %v", err)
	}

	// Test with non-existing directory
	nonExistentDir := filepath.Join(tempDir, "non-existent", "deeply", "nested")
	err = validator.ValidateParentDirectory(nonExistentDir)
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}

	if validationErr, ok := err.(*interfaces.ValidationError); ok {
		if validationErr.Code != "not_exists" {
			t.Errorf("Expected error code 'not_exists', got %s", validationErr.Code)
		}
	}
}

func TestDirectorySelector_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "new", "nested", "directory")

	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}
	selector := NewDirectorySelector(mockUI, mockLogger)

	err := selector.CreateDirectory(testPath)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatFileSize(tt.size)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
func TestDirectorySelector_CreateBackup(t *testing.T) {
	tempDir := t.TempDir()

	// Create source directory with some files
	sourceDir := filepath.Join(tempDir, "source")
	err := os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create test files
	testFile1 := filepath.Join(sourceDir, "file1.txt")
	err = os.WriteFile(testFile1, []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file1: %v", err)
	}

	testFile2 := filepath.Join(sourceDir, "file2.txt")
	err = os.WriteFile(testFile2, []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file2: %v", err)
	}

	// Create subdirectory with file
	subDir := filepath.Join(sourceDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	testFile3 := filepath.Join(subDir, "file3.txt")
	err = os.WriteFile(testFile3, []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file3: %v", err)
	}

	// Create backup
	backupDir := filepath.Join(tempDir, "backup")
	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}
	selector := NewDirectorySelector(mockUI, mockLogger)

	err = selector.CreateBackup(sourceDir, backupDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify backup was created correctly
	// Check that backup directory exists
	info, err := os.Stat(backupDir)
	if err != nil {
		t.Fatalf("Backup directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Backup path is not a directory")
	}

	// Check that files were copied
	backupFile1 := filepath.Join(backupDir, "file1.txt")
	content1, err := os.ReadFile(backupFile1)
	if err != nil {
		t.Fatalf("Failed to read backup file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("Expected content1, got %s", string(content1))
	}

	backupFile2 := filepath.Join(backupDir, "file2.txt")
	content2, err := os.ReadFile(backupFile2)
	if err != nil {
		t.Fatalf("Failed to read backup file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("Expected content2, got %s", string(content2))
	}

	// Check subdirectory was copied
	backupSubDir := filepath.Join(backupDir, "subdir")
	info, err = os.Stat(backupSubDir)
	if err != nil {
		t.Fatalf("Backup subdirectory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Backup subdirectory is not a directory")
	}

	backupFile3 := filepath.Join(backupSubDir, "file3.txt")
	content3, err := os.ReadFile(backupFile3)
	if err != nil {
		t.Fatalf("Failed to read backup file3: %v", err)
	}
	if string(content3) != "content3" {
		t.Errorf("Expected content3, got %s", string(content3))
	}
}

func TestDirectorySelector_GenerateBackupPath(t *testing.T) {
	mockUI := NewDirectoryMockInteractiveUI()
	mockLogger := &DirectoryMockLogger{}
	selector := NewDirectorySelector(mockUI, mockLogger)

	originalPath := "/home/user/projects/my-app"
	backupPath := selector.generateBackupPath(originalPath)

	// Check that backup path is in the same directory
	expectedDir := "/home/user/projects"
	actualDir := filepath.Dir(backupPath)
	if actualDir != expectedDir {
		t.Errorf("Expected backup in directory %s, got %s", expectedDir, actualDir)
	}

	// Check that backup path contains the original name
	backupName := filepath.Base(backupPath)
	if !strings.Contains(backupName, "my-app") {
		t.Errorf("Expected backup name to contain 'my-app', got %s", backupName)
	}

	// Check that backup path contains "backup"
	if !strings.Contains(backupName, "backup") {
		t.Errorf("Expected backup name to contain 'backup', got %s", backupName)
	}
}
