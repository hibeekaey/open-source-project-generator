package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger provides structured logging for the application
type Logger struct {
	level       LogLevel
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errLogger   *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
	component   string // Component name for contextual logging
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, logToFile bool) (*Logger, error) {
	return NewLoggerWithComponent(level, logToFile, "app")
}

// NewLoggerWithComponent creates a new logger instance with a specific component name
func NewLoggerWithComponent(level LogLevel, logToFile bool, component string) (*Logger, error) {
	logger := &Logger{
		level:     level,
		component: component,
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Add file logging if requested
	if logToFile {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		logDir := filepath.Join(homeDir, ".cache", "template-generator", "logs")
		if mkdirErr := os.MkdirAll(logDir, 0750); mkdirErr != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", mkdirErr)
		}

		logFile := filepath.Join(logDir, fmt.Sprintf("generator-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.logFile = file
		writers = append(writers, file)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create loggers for different levels with consistent formatting
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	logger.infoLogger = log.New(multiWriter, "", flags)
	logger.warnLogger = log.New(multiWriter, "", flags)
	logger.errLogger = log.New(multiWriter, "", flags)
	logger.debugLogger = log.New(multiWriter, "", flags)

	return logger, nil
}

// Close closes the logger and any open file handles
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Debug logs a debug message with structured format
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logWithLevel("DEBUG", msg, args...)
	}
}

// Info logs an info message with structured format
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logWithLevel("INFO", msg, args...)
	}
}

// Warn logs a warning message with structured format
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logWithLevel("WARN", msg, args...)
	}
}

// Error logs an error message with structured format
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.logWithLevel("ERROR", msg, args...)
	}
}

// logWithLevel provides consistent structured logging format
func (l *Logger) logWithLevel(level, msg string, args ...interface{}) {
	// Format the message
	formattedMsg := fmt.Sprintf(msg, args...)

	// Create structured log entry
	logEntry := fmt.Sprintf("[%s] component=%s message=\"%s\"", level, l.component, formattedMsg)

	// Select appropriate logger based on level
	var logger *log.Logger
	switch level {
	case "DEBUG":
		logger = l.debugLogger
	case "INFO":
		logger = l.infoLogger
	case "WARN":
		logger = l.warnLogger
	case "ERROR":
		logger = l.errLogger
	default:
		logger = l.infoLogger
	}

	logger.Println(logEntry)
}

// Structured logging methods for better consistency

// DebugWithFields logs a debug message with structured fields
func (l *Logger) DebugWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelDebug {
		l.logWithFields("DEBUG", msg, fields)
	}
}

// InfoWithFields logs an info message with structured fields
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelInfo {
		l.logWithFields("INFO", msg, fields)
	}
}

// WarnWithFields logs a warning message with structured fields
func (l *Logger) WarnWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelWarn {
		l.logWithFields("WARN", msg, fields)
	}
}

// ErrorWithFields logs an error message with structured fields
func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelError {
		l.logWithFields("ERROR", msg, fields)
	}
}

// logWithFields provides structured logging with key-value pairs
func (l *Logger) logWithFields(level, msg string, fields map[string]interface{}) {
	// Build structured log entry
	var parts []string
	parts = append(parts, fmt.Sprintf("component=%s", l.component))
	parts = append(parts, fmt.Sprintf("message=\"%s\"", msg))

	// Add fields
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	logEntry := fmt.Sprintf("[%s] %s", level, strings.Join(parts, " "))

	// Select appropriate logger
	var logger *log.Logger
	switch level {
	case "DEBUG":
		logger = l.debugLogger
	case "INFO":
		logger = l.infoLogger
	case "WARN":
		logger = l.warnLogger
	case "ERROR":
		logger = l.errLogger
	default:
		logger = l.infoLogger
	}

	logger.Println(logEntry)
}

// Operation logging methods for tracking key operations

// LogOperationStart logs the start of an operation
func (l *Logger) LogOperationStart(operation string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["operation"] = operation
	fields["status"] = "started"
	l.InfoWithFields("Operation started", fields)
}

// LogOperationSuccess logs successful completion of an operation
func (l *Logger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["operation"] = operation
	fields["status"] = "success"
	fields["duration_ms"] = duration.Milliseconds()
	l.InfoWithFields("Operation completed successfully", fields)
}

// LogOperationError logs an error during an operation
func (l *Logger) LogOperationError(operation string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["operation"] = operation
	fields["status"] = "error"
	fields["error"] = err.Error()
	l.ErrorWithFields("Operation failed", fields)
}

// LogLevelFromString converts a string to LogLevel
func LogLevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}
