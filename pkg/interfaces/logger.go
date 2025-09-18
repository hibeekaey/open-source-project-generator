package interfaces

import (
	"time"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// OperationContext holds context for tracking operations
type OperationContext struct {
	Operation string
	StartTime time.Time
	Fields    map[string]interface{}
}

// Logger provides comprehensive logging functionality
type Logger interface {
	// Basic logging methods
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	// Structured logging methods
	DebugWithFields(msg string, fields map[string]interface{})
	InfoWithFields(msg string, fields map[string]interface{})
	WarnWithFields(msg string, fields map[string]interface{})
	ErrorWithFields(msg string, fields map[string]interface{})
	FatalWithFields(msg string, fields map[string]interface{})

	// Error logging with error objects
	ErrorWithError(msg string, err error, fields map[string]interface{})

	// Operation tracking
	StartOperation(operation string, fields map[string]interface{}) *OperationContext
	LogOperationStart(operation string, fields map[string]interface{})
	LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{})
	LogOperationError(operation string, err error, fields map[string]interface{})
	FinishOperation(ctx *OperationContext, additionalFields map[string]interface{})
	FinishOperationWithError(ctx *OperationContext, err error, additionalFields map[string]interface{})

	// Performance logging
	LogPerformanceMetrics(operation string, metrics map[string]interface{})
	LogMemoryUsage(operation string)

	// Configuration methods
	SetLevel(level int)
	GetLevel() int
	SetJSONOutput(enable bool)
	SetCallerInfo(enable bool)
	IsDebugEnabled() bool
	IsInfoEnabled() bool

	// Context methods
	WithComponent(component string) Logger
	WithFields(fields map[string]interface{}) LoggerContext

	// Log management
	GetLogDir() string
	GetRecentEntries(limit int) []LogEntry
	FilterEntries(level string, component string, since time.Time, limit int) []LogEntry
	GetLogFiles() ([]string, error)
	ReadLogFile(filename string) ([]byte, error)

	// Lifecycle
	Close() error
}

// LoggerContext provides temporary logging context with additional fields
type LoggerContext interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	ErrorWithError(msg string, err error)
}
