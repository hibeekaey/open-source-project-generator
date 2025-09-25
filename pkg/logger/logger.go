// Package logger provides a unified logging system for the entire application.
//
// This package consolidates all logging logic from various packages into a single,
// comprehensive logging system that can be used across the application.
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel parses a string to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelInfo
	}
}

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

// LogRotationConfig holds configuration for log rotation
type LogRotationConfig struct {
	MaxSize    int64         `json:"max_size"`
	MaxAge     time.Duration `json:"max_age"`
	MaxBackups int           `json:"max_backups"`
	Compress   bool          `json:"compress"`
}

// DefaultLogRotationConfig returns default log rotation settings
func DefaultLogRotationConfig() *LogRotationConfig {
	return &LogRotationConfig{
		MaxSize:    10 * 1024 * 1024,   // 10MB
		MaxAge:     7 * 24 * time.Hour, // 7 days
		MaxBackups: 5,
		Compress:   true,
	}
}

// UnifiedLogger provides comprehensive logging capabilities
type UnifiedLogger struct {
	level          LogLevel
	component      string
	enableJSON     bool
	enableCaller   bool
	enableColors   bool
	rotationConfig *LogRotationConfig

	// Output destinations
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errLogger   *log.Logger
	fatalLogger *log.Logger
	debugLogger *log.Logger

	// File logging
	logFile *os.File
	logPath string

	// In-memory storage
	entries      []LogEntry
	maxEntries   int
	entriesMutex sync.RWMutex

	// Performance tracking
	operationStart map[string]time.Time
	operationMutex sync.RWMutex
}

// LoggerConfig holds configuration for the unified logger
type LoggerConfig struct {
	Level          LogLevel           `json:"level"`
	Component      string             `json:"component"`
	EnableJSON     bool               `json:"enable_json"`
	EnableCaller   bool               `json:"enable_caller"`
	EnableColors   bool               `json:"enable_colors"`
	LogFile        string             `json:"log_file"`
	MaxEntries     int                `json:"max_entries"`
	RotationConfig *LogRotationConfig `json:"rotation_config"`
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          LevelInfo,
		Component:      "app",
		EnableJSON:     false,
		EnableCaller:   false,
		EnableColors:   true,
		LogFile:        "",
		MaxEntries:     1000,
		RotationConfig: DefaultLogRotationConfig(),
	}
}

// NewUnifiedLogger creates a new unified logger
func NewUnifiedLogger(config *LoggerConfig) (*UnifiedLogger, error) {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	logger := &UnifiedLogger{
		level:          config.Level,
		component:      config.Component,
		enableJSON:     config.EnableJSON,
		enableCaller:   config.EnableCaller,
		enableColors:   config.EnableColors,
		rotationConfig: config.RotationConfig,
		maxEntries:     config.MaxEntries,
		entries:        make([]LogEntry, 0),
		operationStart: make(map[string]time.Time),
	}

	// Set up output destinations
	logger.setupOutputs(config.LogFile)

	return logger, nil
}

// setupOutputs sets up the output destinations
func (l *UnifiedLogger) setupOutputs(logFile string) {
	// Default to stdout/stderr
	l.infoLogger = log.New(os.Stdout, "", 0)
	l.warnLogger = log.New(os.Stderr, "", 0)
	l.errLogger = log.New(os.Stderr, "", 0)
	l.fatalLogger = log.New(os.Stderr, "", 0)
	l.debugLogger = log.New(os.Stdout, "", 0)

	// Set up file logging if specified
	if logFile != "" {
		if err := l.setupFileLogging(logFile); err != nil {
			// Fall back to stdout/stderr if file logging fails
			fmt.Fprintf(os.Stderr, "Failed to setup file logging: %v\n", err)
		}
	}
}

// setupFileLogging sets up file-based logging
func (l *UnifiedLogger) setupFileLogging(logFile string) error {
	// Validate and sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(logFile)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid log file path: %s", logFile)
	}

	// Ensure directory exists
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	file, err := os.OpenFile(cleanPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	l.logFile = file
	l.logPath = logFile

	// Update loggers to write to file
	l.infoLogger = log.New(file, "", 0)
	l.warnLogger = log.New(file, "", 0)
	l.errLogger = log.New(file, "", 0)
	l.fatalLogger = log.New(file, "", 0)
	l.debugLogger = log.New(file, "", 0)

	return nil
}

// Basic logging methods

// Debug logs a debug message
func (l *UnifiedLogger) Debug(msg string, args ...interface{}) {
	l.logWithLevel(LevelDebug, msg, nil, args...)
}

// Info logs an info message
func (l *UnifiedLogger) Info(msg string, args ...interface{}) {
	l.logWithLevel(LevelInfo, msg, nil, args...)
}

// Warn logs a warning message
func (l *UnifiedLogger) Warn(msg string, args ...interface{}) {
	l.logWithLevel(LevelWarn, msg, nil, args...)
}

// Error logs an error message
func (l *UnifiedLogger) Error(msg string, args ...interface{}) {
	l.logWithLevel(LevelError, msg, nil, args...)
}

// Fatal logs a fatal message and exits
func (l *UnifiedLogger) Fatal(msg string, args ...interface{}) {
	l.logWithLevel(LevelFatal, msg, nil, args...)
	os.Exit(1)
}

// Structured logging methods

// DebugWithFields logs a debug message with fields
func (l *UnifiedLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LevelDebug, msg, fields)
}

// InfoWithFields logs an info message with fields
func (l *UnifiedLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LevelInfo, msg, fields)
}

// WarnWithFields logs a warning message with fields
func (l *UnifiedLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LevelWarn, msg, fields)
}

// ErrorWithFields logs an error message with fields
func (l *UnifiedLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LevelError, msg, fields)
}

// FatalWithFields logs a fatal message with fields and exits
func (l *UnifiedLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LevelFatal, msg, fields)
	os.Exit(1)
}

// Error logging with error objects

// ErrorWithError logs an error message with an error object
func (l *UnifiedLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	l.logWithLevel(LevelError, msg, fields)
}

// Operation tracking

// StartOperation starts tracking an operation
func (l *UnifiedLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	l.operationMutex.Lock()
	defer l.operationMutex.Unlock()

	startTime := time.Now()
	l.operationStart[operation] = startTime

	ctx := &interfaces.OperationContext{
		Operation: operation,
		StartTime: startTime,
		Fields:    fields,
	}

	l.InfoWithFields(fmt.Sprintf("Starting operation: %s", operation), fields)
	return ctx
}

// LogOperationStart logs the start of an operation
func (l *UnifiedLogger) LogOperationStart(operation string, fields map[string]interface{}) {
	l.StartOperation(operation, fields)
}

// LogOperationSuccess logs successful completion of an operation
func (l *UnifiedLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["duration"] = duration.String()
	fields["status"] = "success"

	l.InfoWithFields(fmt.Sprintf("Operation completed successfully: %s", operation), fields)
}

// LogOperationError logs an error during an operation
func (l *UnifiedLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	fields["status"] = "error"

	l.ErrorWithFields(fmt.Sprintf("Operation failed: %s", operation), fields)
}

// FinishOperation finishes tracking an operation
func (l *UnifiedLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
	duration := time.Since(ctx.StartTime)

	fields := ctx.Fields
	for k, v := range additionalFields {
		fields[k] = v
	}

	l.LogOperationSuccess(ctx.Operation, duration, fields)

	// Clean up
	l.operationMutex.Lock()
	delete(l.operationStart, ctx.Operation)
	l.operationMutex.Unlock()
}

// FinishOperationWithError finishes tracking an operation with an error
func (l *UnifiedLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
	fields := ctx.Fields
	for k, v := range additionalFields {
		fields[k] = v
	}

	l.LogOperationError(ctx.Operation, err, fields)

	// Clean up
	l.operationMutex.Lock()
	delete(l.operationStart, ctx.Operation)
	l.operationMutex.Unlock()
}

// Performance logging

// LogPerformanceMetrics logs performance metrics
func (l *UnifiedLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
	fields := make(map[string]interface{})
	fields["operation"] = operation
	fields["metrics"] = metrics

	l.InfoWithFields("Performance metrics", fields)
}

// OperationContext holds context for tracking operations
type OperationContext struct {
	Operation string
	StartTime time.Time
	Fields    map[string]interface{}
}

// Core logging method

// logWithLevel provides consistent structured logging format
func (l *UnifiedLogger) logWithLevel(level LogLevel, msg string, fields map[string]interface{}, args ...interface{}) {
	if l.level > level {
		return
	}

	// Format the message
	formattedMsg := fmt.Sprintf(msg, args...)

	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level.String(),
		Component: l.component,
		Message:   formattedMsg,
		Fields:    fields,
	}

	// Add caller information if enabled
	if l.enableCaller {
		if _, file, line, ok := runtime.Caller(3); ok {
			entry.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	// Add to in-memory store
	l.addEntry(entry)

	// Format and log the entry
	var logOutput string
	if l.enableJSON {
		if jsonBytes, err := json.Marshal(entry); err == nil {
			logOutput = string(jsonBytes)
		} else {
			logOutput = l.formatTextEntry(entry)
		}
	} else {
		logOutput = l.formatTextEntry(entry)
	}

	// Select appropriate logger based on level
	var logger *log.Logger
	switch level {
	case LevelDebug:
		logger = l.debugLogger
	case LevelInfo:
		logger = l.infoLogger
	case LevelWarn:
		logger = l.warnLogger
	case LevelError:
		logger = l.errLogger
	case LevelFatal:
		logger = l.fatalLogger
	default:
		logger = l.infoLogger
	}

	logger.Println(logOutput)

	// Rotate log if needed (check periodically)
	if l.logFile != nil {
		go func() {
			if err := l.rotateIfNeeded(); err != nil {
				fmt.Fprintf(os.Stderr, "Log rotation error: %v\n", err)
			}
		}()
	}
}

// formatTextEntry formats a log entry as human-readable text
func (l *UnifiedLogger) formatTextEntry(entry LogEntry) string {
	var parts []string

	// Add timestamp
	parts = append(parts, fmt.Sprintf("[%s]", entry.Timestamp.Format("2006-01-02 15:04:05")))

	// Add level with color if enabled
	levelStr := entry.Level
	if l.enableColors {
		levelStr = l.colorizeLevel(entry.Level)
	}
	parts = append(parts, fmt.Sprintf("[%s]", levelStr))

	// Add component
	parts = append(parts, fmt.Sprintf("[%s]", entry.Component))

	// Add caller if available
	if entry.Caller != "" {
		parts = append(parts, fmt.Sprintf("[%s]", entry.Caller))
	}

	// Add message
	parts = append(parts, entry.Message)

	// Add fields if present
	if len(entry.Fields) > 0 {
		fieldParts := make([]string, 0, len(entry.Fields))
		for k, v := range entry.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, ", ")))
	}

	return strings.Join(parts, " ")
}

// colorizeLevel adds color to log level
func (l *UnifiedLogger) colorizeLevel(level string) string {
	if !l.enableColors {
		return level
	}

	switch level {
	case "DEBUG":
		return "\033[36m" + level + "\033[0m" // Cyan
	case "INFO":
		return "\033[32m" + level + "\033[0m" // Green
	case "WARN":
		return "\033[33m" + level + "\033[0m" // Yellow
	case "ERROR":
		return "\033[31m" + level + "\033[0m" // Red
	case "FATAL":
		return "\033[35m" + level + "\033[0m" // Magenta
	default:
		return level
	}
}

// addEntry adds an entry to the in-memory store
func (l *UnifiedLogger) addEntry(entry LogEntry) {
	l.entriesMutex.Lock()
	defer l.entriesMutex.Unlock()

	l.entries = append(l.entries, entry)

	// Trim entries if we exceed maxEntries
	if len(l.entries) > l.maxEntries {
		l.entries = l.entries[len(l.entries)-l.maxEntries:]
	}
}

// GetEntries returns all log entries
func (l *UnifiedLogger) GetEntries() []LogEntry {
	l.entriesMutex.RLock()
	defer l.entriesMutex.RUnlock()

	// Return a copy
	entries := make([]LogEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

// GetEntriesByLevel returns log entries for a specific level
func (l *UnifiedLogger) GetEntriesByLevel(level LogLevel) []LogEntry {
	l.entriesMutex.RLock()
	defer l.entriesMutex.RUnlock()

	var entries []LogEntry
	for _, entry := range l.entries {
		if entry.Level == level.String() {
			entries = append(entries, entry)
		}
	}
	return entries
}

// ClearEntries clears all log entries
func (l *UnifiedLogger) ClearEntries() {
	l.entriesMutex.Lock()
	defer l.entriesMutex.Unlock()
	l.entries = make([]LogEntry, 0)
}

// SetLevel sets the logging level
func (l *UnifiedLogger) SetLevel(level int) {
	l.level = LogLevel(level)
}

// SetLevelInt sets the logging level (interface compatibility)
func (l *UnifiedLogger) SetLevelInt(level int) {
	l.level = LogLevel(level)
}

// GetLevel returns the current logging level
func (l *UnifiedLogger) GetLevel() int {
	return int(l.level)
}

// GetLevelInt returns the current logging level as int (interface compatibility)
func (l *UnifiedLogger) GetLevelInt() int {
	return int(l.level)
}

// Close closes the logger and any open files
func (l *UnifiedLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// rotateIfNeeded checks if log rotation is needed and performs it
func (l *UnifiedLogger) rotateIfNeeded() error {
	if l.logFile == nil || l.rotationConfig == nil {
		return nil
	}

	// Check file size
	info, err := l.logFile.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= l.rotationConfig.MaxSize {
		return l.rotateLog()
	}

	return nil
}

// rotateLog performs log rotation
func (l *UnifiedLogger) rotateLog() error {
	// Close current file
	if err := l.logFile.Close(); err != nil {
		return err
	}

	// Create rotated filename with timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	rotatedPath := fmt.Sprintf("%s.%s", l.logPath, timestamp)

	// Rename current file
	if err := os.Rename(l.logPath, rotatedPath); err != nil {
		return err
	}

	// Open new log file
	file, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	l.logFile = file

	// Update loggers
	l.infoLogger = log.New(file, "", 0)
	l.warnLogger = log.New(file, "", 0)
	l.errLogger = log.New(file, "", 0)
	l.fatalLogger = log.New(file, "", 0)
	l.debugLogger = log.New(file, "", 0)

	return nil
}

// Global logger instance
var (
	globalLogger *UnifiedLogger
	globalOnce   sync.Once
)

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *UnifiedLogger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *UnifiedLogger {
	if globalLogger == nil {
		globalOnce.Do(func() {
			config := DefaultLoggerConfig()
			config.Component = "global"
			logger, err := NewUnifiedLogger(config)
			if err != nil {
				// Fallback to basic logger
				globalLogger = &UnifiedLogger{
					level:       LevelInfo,
					component:   "global",
					infoLogger:  log.New(os.Stdout, "", 0),
					warnLogger:  log.New(os.Stderr, "", 0),
					errLogger:   log.New(os.Stderr, "", 0),
					fatalLogger: log.New(os.Stderr, "", 0),
					debugLogger: log.New(os.Stdout, "", 0),
				}
			} else {
				globalLogger = logger
			}
		})
	}
	return globalLogger
}

// Convenience functions for global logger

// Debug logs a debug message using the global logger
func Debug(msg string, args ...interface{}) {
	GetGlobalLogger().Debug(msg, args...)
}

// Info logs an info message using the global logger
func Info(msg string, args ...interface{}) {
	GetGlobalLogger().Info(msg, args...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, args ...interface{}) {
	GetGlobalLogger().Warn(msg, args...)
}

// Error logs an error message using the global logger
func Error(msg string, args ...interface{}) {
	GetGlobalLogger().Error(msg, args...)
}

// Fatal logs a fatal message using the global logger and exits
func Fatal(msg string, args ...interface{}) {
	GetGlobalLogger().Fatal(msg, args...)
}

// Missing methods to implement interfaces.Logger interface

// GetLogDir returns the log directory
func (l *UnifiedLogger) GetLogDir() string {
	if l.logPath != "" {
		return filepath.Dir(l.logPath)
	}
	return ""
}

// GetRecentEntries returns recent log entries
func (l *UnifiedLogger) GetRecentEntries(limit int) []interfaces.LogEntry {
	entries := l.GetEntries()
	if limit <= 0 || limit >= len(entries) {
		// Convert to interfaces.LogEntry
		result := make([]interfaces.LogEntry, len(entries))
		for i, entry := range entries {
			result[i] = interfaces.LogEntry{
				Timestamp: entry.Timestamp,
				Level:     entry.Level,
				Component: entry.Component,
				Message:   entry.Message,
				Fields:    entry.Fields,
				Caller:    entry.Caller,
				Error:     entry.Error,
			}
		}
		return result
	}

	// Convert to interfaces.LogEntry
	result := make([]interfaces.LogEntry, limit)
	for i, entry := range entries[len(entries)-limit:] {
		result[i] = interfaces.LogEntry{
			Timestamp: entry.Timestamp,
			Level:     entry.Level,
			Component: entry.Component,
			Message:   entry.Message,
			Fields:    entry.Fields,
			Caller:    entry.Caller,
			Error:     entry.Error,
		}
	}
	return result
}

// FilterEntries filters log entries based on criteria
func (l *UnifiedLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	entries := l.GetEntries()
	var filtered []LogEntry

	for _, entry := range entries {
		// Filter by level
		if level != "" && entry.Level != level {
			continue
		}

		// Filter by component
		if component != "" && entry.Component != component {
			continue
		}

		// Filter by time
		if !since.IsZero() && entry.Timestamp.Before(since) {
			continue
		}

		filtered = append(filtered, entry)

		// Apply limit
		if limit > 0 && len(filtered) >= limit {
			break
		}
	}

	// Convert to interfaces.LogEntry
	result := make([]interfaces.LogEntry, len(filtered))
	for i, entry := range filtered {
		result[i] = interfaces.LogEntry{
			Timestamp: entry.Timestamp,
			Level:     entry.Level,
			Component: entry.Component,
			Message:   entry.Message,
			Fields:    entry.Fields,
			Caller:    entry.Caller,
			Error:     entry.Error,
		}
	}

	return result
}

// GetLogFiles returns available log files
func (l *UnifiedLogger) GetLogFiles() ([]string, error) {
	if l.logPath == "" {
		return []string{}, nil
	}

	logDir := filepath.Dir(l.logPath)
	files, err := filepath.Glob(filepath.Join(logDir, "*.log*"))
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ReadLogFile reads a log file
func (l *UnifiedLogger) ReadLogFile(filename string) ([]byte, error) {
	// Validate and sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(filename)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
		return nil, fmt.Errorf("invalid filename: %s", filename)
	}
	return os.ReadFile(cleanPath)
}

// LogMemoryUsage logs memory usage metrics
func (l *UnifiedLogger) LogMemoryUsage(operation string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
		"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
		"sys_mb":         float64(m.Sys) / 1024 / 1024,
		"num_gc":         m.NumGC,
	}

	l.LogPerformanceMetrics(operation, metrics)
}

// SetJSONOutput enables/disables JSON output
func (l *UnifiedLogger) SetJSONOutput(enable bool) {
	l.enableJSON = enable
}

// SetCallerInfo enables/disables caller information
func (l *UnifiedLogger) SetCallerInfo(enable bool) {
	l.enableCaller = enable
}

// IsDebugEnabled checks if debug logging is enabled
func (l *UnifiedLogger) IsDebugEnabled() bool {
	return l.level <= LevelDebug
}

// IsInfoEnabled checks if info logging is enabled
func (l *UnifiedLogger) IsInfoEnabled() bool {
	return l.level <= LevelInfo
}

// WithComponent creates a logger with a specific component
func (l *UnifiedLogger) WithComponent(component string) interfaces.Logger {
	return &UnifiedLogger{
		level:          l.level,
		component:      component,
		enableJSON:     l.enableJSON,
		enableCaller:   l.enableCaller,
		enableColors:   l.enableColors,
		rotationConfig: l.rotationConfig,
		maxEntries:     l.maxEntries,
		entries:        l.entries,
		operationStart: l.operationStart,
		infoLogger:     l.infoLogger,
		warnLogger:     l.warnLogger,
		errLogger:      l.errLogger,
		fatalLogger:    l.fatalLogger,
		debugLogger:    l.debugLogger,
		logFile:        l.logFile,
		logPath:        l.logPath,
	}
}

// WithFields creates a logger context with additional fields
func (l *UnifiedLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &LoggerContextImpl{
		logger: l,
		fields: fields,
	}
}

// LoggerContextImpl implements LoggerContext
type LoggerContextImpl struct {
	logger *UnifiedLogger
	fields map[string]interface{}
}

// Debug logs a debug message with context fields
func (ctx *LoggerContextImpl) Debug(msg string, args ...interface{}) {
	ctx.logger.DebugWithFields(msg, ctx.fields)
}

// Info logs an info message with context fields
func (ctx *LoggerContextImpl) Info(msg string, args ...interface{}) {
	ctx.logger.InfoWithFields(msg, ctx.fields)
}

// Warn logs a warning message with context fields
func (ctx *LoggerContextImpl) Warn(msg string, args ...interface{}) {
	ctx.logger.WarnWithFields(msg, ctx.fields)
}

// Error logs an error message with context fields
func (ctx *LoggerContextImpl) Error(msg string, args ...interface{}) {
	ctx.logger.ErrorWithFields(msg, ctx.fields)
}

// ErrorWithError logs an error message with error object and context fields
func (ctx *LoggerContextImpl) ErrorWithError(msg string, err error) {
	fields := make(map[string]interface{})
	for k, v := range ctx.fields {
		fields[k] = v
	}
	fields["error"] = err.Error()
	ctx.logger.ErrorWithFields(msg, fields)
}
