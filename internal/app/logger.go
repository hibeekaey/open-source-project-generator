package app

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// Re-export comprehensive logging types and functions

// LogLevel represents the logging level
type LogLevel = logger.LogLevel

const (
	LogLevelDebug = logger.LevelDebug
	LogLevelInfo  = logger.LevelInfo
	LogLevelWarn  = logger.LevelWarn
	LogLevelError = logger.LevelError
	LogLevelFatal = logger.LevelFatal
)

// LogEntry represents a structured log entry
type LogEntry = logger.LogEntry

// LogRotationConfig holds configuration for log rotation
type LogRotationConfig = logger.LogRotationConfig

// DefaultLogRotationConfig returns default log rotation settings
func DefaultLogRotationConfig() *LogRotationConfig {
	return logger.DefaultLogRotationConfig()
}

// Logger provides structured logging for the application
type Logger = logger.Logger

// NewLogger creates a new logger instance
func NewLogger(component string, level LogLevel, enableJSON, enableCaller, enableColors bool, logFile string) *Logger {
	config := &logger.LoggerConfig{
		Level:        level,
		Component:    component,
		EnableJSON:   enableJSON,
		EnableCaller: enableCaller,
		EnableColors: enableColors,
		LogFile:      logFile,
	}

	loggerInstance, err := logger.NewLogger(config)
	if err != nil {
		// Fallback to basic logger
		config := logger.DefaultLoggerConfig()
		config.Component = component
		loggerInstance, _ = logger.NewLogger(config)
	}

	return loggerInstance
}

// NewLoggerWithConfig creates a new logger with configuration
func NewLoggerWithConfig(config *LoggerConfig) (*Logger, error) {
	loggerConfig := &logger.LoggerConfig{
		Level:        config.Level,
		Component:    config.Component,
		EnableJSON:   config.EnableJSON,
		EnableCaller: config.EnableCaller,
		EnableColors: config.EnableColors,
		LogFile:      config.LogFile,
		MaxEntries:   config.MaxEntries,
	}

	return logger.NewLogger(loggerConfig)
}

// LoggerConfig holds configuration for the logger
type LoggerConfig = logger.LoggerConfig

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return logger.DefaultLoggerConfig()
}

// OperationContext holds context for tracking operations
type OperationContext struct {
	Operation string
	StartTime time.Time
	Fields    map[string]interface{}
}

// Convenience functions that delegate to the comprehensive logger

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	logger.Fatal(msg, args...)
}
