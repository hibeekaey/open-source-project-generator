package test_helpers

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// MockLogger provides a no-op implementation of the Logger interface for testing
type MockLogger struct{}

// All logging methods are no-ops for testing
func (m *MockLogger) Debug(format string, args ...interface{})                            {}
func (m *MockLogger) Info(format string, args ...interface{})                             {}
func (m *MockLogger) Warn(format string, args ...interface{})                             {}
func (m *MockLogger) Error(format string, args ...interface{})                            {}
func (m *MockLogger) Fatal(format string, args ...interface{})                            {}
func (m *MockLogger) DebugWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) InfoWithFields(message string, fields map[string]interface{})        {}
func (m *MockLogger) WarnWithFields(message string, fields map[string]interface{})        {}
func (m *MockLogger) ErrorWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) FatalWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {}
func (m *MockLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	return &interfaces.OperationContext{}
}
func (m *MockLogger) LogOperationStart(operation string, fields map[string]interface{}) {}
func (m *MockLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
}
func (m *MockLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {}
func (m *MockLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
}
func (m *MockLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
}
func (m *MockLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {}
func (m *MockLogger) LogMemoryUsage(operation string)                                        {}
func (m *MockLogger) SetLevel(level int)                                                     {}
func (m *MockLogger) GetLevel() int                                                          { return 0 }
func (m *MockLogger) SetJSONOutput(enabled bool)                                             {}
func (m *MockLogger) SetCallerInfo(enabled bool)                                             {}
func (m *MockLogger) IsDebugEnabled() bool                                                   { return true }
func (m *MockLogger) IsInfoEnabled() bool                                                    { return true }
func (m *MockLogger) WithComponent(component string) interfaces.Logger                       { return m }
func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &MockLoggerContext{}
}
func (m *MockLogger) GetLogDir() string                                { return "/tmp" }
func (m *MockLogger) GetRecentEntries(limit int) []interfaces.LogEntry { return nil }
func (m *MockLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (m *MockLogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (m *MockLogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }
func (m *MockLogger) Close() error                                { return nil }

// MockLoggerContext provides a no-op implementation of LoggerContext for testing
type MockLoggerContext struct{}

func (m *MockLoggerContext) Debug(msg string, args ...interface{}) {}
func (m *MockLoggerContext) Info(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Warn(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Error(msg string, args ...interface{}) {}
func (m *MockLoggerContext) ErrorWithError(msg string, err error)  {}

// NewMockLogger creates a new mock logger instance
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
