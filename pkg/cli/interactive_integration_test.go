package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

func TestCLI_runInteractiveDirectorySelection(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test-project")

	// Create a mock CLI with minimal dependencies
	cli := &CLI{
		logger: &MockCLILogger{},
	}
	// Initialize the OutputFormatter
	cli.outputFormatter = NewOutputFormatter(false, false, false, cli.logger)

	// Create a mock UI that will return our test path
	mockUI := &MockCLIInteractiveUI{
		textResponse: testPath,
	}
	cli.interactiveUI = mockUI

	ctx := context.Background()
	defaultPath := "output/generated"

	result, err := cli.runInteractiveDirectorySelection(ctx, defaultPath)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != testPath {
		t.Errorf("Expected path %s, got %s", testPath, result)
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

// MockCLILogger is a minimal mock logger for CLI testing
type MockCLILogger struct{}

func (m *MockCLILogger) Debug(msg string, args ...interface{})                               {}
func (m *MockCLILogger) Info(msg string, args ...interface{})                                {}
func (m *MockCLILogger) Warn(msg string, args ...interface{})                                {}
func (m *MockCLILogger) Error(msg string, args ...interface{})                               {}
func (m *MockCLILogger) Fatal(msg string, args ...interface{})                               {}
func (m *MockCLILogger) DebugWithFields(message string, fields map[string]interface{})       {}
func (m *MockCLILogger) InfoWithFields(message string, fields map[string]interface{})        {}
func (m *MockCLILogger) WarnWithFields(message string, fields map[string]interface{})        {}
func (m *MockCLILogger) ErrorWithFields(message string, fields map[string]interface{})       {}
func (m *MockCLILogger) FatalWithFields(message string, fields map[string]interface{})       {}
func (m *MockCLILogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {}
func (m *MockCLILogger) SetLevel(level int)                                                  {}
func (m *MockCLILogger) GetLevel() int                                                       { return 0 }
func (m *MockCLILogger) SetJSONOutput(enabled bool)                                          {}
func (m *MockCLILogger) SetCallerInfo(enabled bool)                                          {}
func (m *MockCLILogger) IsDebugEnabled() bool                                                { return false }
func (m *MockCLILogger) IsInfoEnabled() bool                                                 { return false }
func (m *MockCLILogger) StartOperation(operation string, metadata map[string]interface{}) *interfaces.OperationContext {
	return nil
}
func (m *MockCLILogger) LogOperationStart(operation string, fields map[string]interface{}) {}
func (m *MockCLILogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
}
func (m *MockCLILogger) LogOperationError(operation string, err error, fields map[string]interface{}) {
}
func (m *MockCLILogger) FinishOperation(ctx *interfaces.OperationContext, metadata map[string]interface{}) {
}
func (m *MockCLILogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, metadata map[string]interface{}) {
}
func (m *MockCLILogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {}
func (m *MockCLILogger) LogMemoryUsage(operation string)                                        {}
func (m *MockCLILogger) WithComponent(component string) interfaces.Logger                       { return m }
func (m *MockCLILogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return nil
}
func (m *MockCLILogger) GetLogDir() string                                { return "" }
func (m *MockCLILogger) GetRecentEntries(limit int) []interfaces.LogEntry { return nil }
func (m *MockCLILogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (m *MockCLILogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (m *MockCLILogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }
func (m *MockCLILogger) Close() error                                { return nil }

// MockCLIInteractiveUI is a minimal mock UI for CLI testing
type MockCLIInteractiveUI struct {
	textResponse string
}

func (m *MockCLIInteractiveUI) PromptText(ctx context.Context, config interfaces.TextPromptConfig) (*interfaces.TextResult, error) {
	return &interfaces.TextResult{Value: m.textResponse, Action: "submit"}, nil
}

func (m *MockCLIInteractiveUI) ShowMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	return &interfaces.MenuResult{SelectedIndex: 0, SelectedValue: "cancel", Action: "select"}, nil
}

func (m *MockCLIInteractiveUI) PromptConfirm(ctx context.Context, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, error) {
	return &interfaces.ConfirmResult{Confirmed: false, Action: "confirm"}, nil
}

func (m *MockCLIInteractiveUI) ShowTable(ctx context.Context, config interfaces.TableConfig) error {
	return nil
}

// Implement other required methods with no-op implementations
func (m *MockCLIInteractiveUI) ShowMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) ShowCheckboxList(ctx context.Context, config interfaces.CheckboxConfig) (*interfaces.CheckboxResult, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) PromptSelect(ctx context.Context, config interfaces.SelectConfig) (*interfaces.SelectResult, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) ShowTree(ctx context.Context, config interfaces.TreeConfig) error {
	return nil
}
func (m *MockCLIInteractiveUI) ShowProgress(ctx context.Context, config interfaces.ProgressConfig) (interfaces.ProgressTracker, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) ShowBreadcrumb(ctx context.Context, path []string) error {
	return nil
}
func (m *MockCLIInteractiveUI) ShowHelp(ctx context.Context, helpContext string) error {
	return nil
}
func (m *MockCLIInteractiveUI) ShowError(ctx context.Context, config interfaces.ErrorConfig) (*interfaces.ErrorResult, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) StartSession(ctx context.Context, config interfaces.SessionConfig) (*interfaces.UISession, error) {
	return nil, nil
}
func (m *MockCLIInteractiveUI) EndSession(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (m *MockCLIInteractiveUI) SaveSessionState(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (m *MockCLIInteractiveUI) RestoreSessionState(ctx context.Context, sessionID string) (*interfaces.UISession, error) {
	return nil, nil
}
