package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
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
	MaxSize    int64         // Maximum size in bytes before rotation
	MaxAge     time.Duration // Maximum age before rotation
	MaxBackups int           // Maximum number of backup files to keep
	Compress   bool          // Whether to compress rotated files
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

// Logger provides structured logging for the application
type Logger struct {
	level          LogLevel
	infoLogger     *log.Logger
	warnLogger     *log.Logger
	errLogger      *log.Logger
	debugLogger    *log.Logger
	fatalLogger    *log.Logger
	logFile        *os.File
	component      string
	enableJSON     bool
	enableCaller   bool
	rotationConfig *LogRotationConfig
	logDir         string
	mu             sync.RWMutex
	entries        []LogEntry // Recent log entries for logs command
	maxEntries     int        // Maximum entries to keep in memory
}

// LoggerConfig holds configuration for logger creation
type LoggerConfig struct {
	Level          LogLevel
	LogToFile      bool
	Component      string
	EnableJSON     bool
	EnableCaller   bool
	RotationConfig *LogRotationConfig
	MaxEntries     int
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          LogLevelInfo,
		LogToFile:      true,
		Component:      "app",
		EnableJSON:     false,
		EnableCaller:   false,
		RotationConfig: DefaultLogRotationConfig(),
		MaxEntries:     1000,
	}
}

// NewLogger creates a new logger instance with default configuration
func NewLogger(level LogLevel, logToFile bool) (*Logger, error) {
	config := DefaultLoggerConfig()
	config.Level = level
	config.LogToFile = logToFile
	return NewLoggerWithConfig(config)
}

// NewLoggerWithComponent creates a new logger instance with a specific component name
func NewLoggerWithComponent(level LogLevel, logToFile bool, component string) (*Logger, error) {
	config := DefaultLoggerConfig()
	config.Level = level
	config.LogToFile = logToFile
	config.Component = component
	return NewLoggerWithConfig(config)
}

// NewLoggerWithConfig creates a new logger instance with full configuration
func NewLoggerWithConfig(config *LoggerConfig) (*Logger, error) {
	logger := &Logger{
		level:          config.Level,
		component:      config.Component,
		enableJSON:     config.EnableJSON,
		enableCaller:   config.EnableCaller,
		rotationConfig: config.RotationConfig,
		maxEntries:     config.MaxEntries,
		entries:        make([]LogEntry, 0, config.MaxEntries),
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Add file logging if requested
	if config.LogToFile {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		logDir := filepath.Join(homeDir, ".cache", "template-generator", "logs")
		logger.logDir = logDir
		if mkdirErr := os.MkdirAll(logDir, 0750); mkdirErr != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", mkdirErr)
		}

		logFile := filepath.Join(logDir, fmt.Sprintf("generator-%s.log", time.Now().Format("2006-01-02")))
		file, err := utils.SafeOpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.logFile = file
		writers = append(writers, file)

		// Perform log rotation if needed
		if err := logger.rotateIfNeeded(); err != nil {
			// Log rotation failure shouldn't prevent logger creation
			fmt.Fprintf(os.Stderr, "Warning: log rotation failed: %v\n", err)
		}
	}

	multiWriter := io.MultiWriter(writers...)

	// Create loggers for different levels with consistent formatting
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	logger.infoLogger = log.New(multiWriter, "", flags)
	logger.warnLogger = log.New(multiWriter, "", flags)
	logger.errLogger = log.New(multiWriter, "", flags)
	logger.debugLogger = log.New(multiWriter, "", flags)
	logger.fatalLogger = log.New(multiWriter, "", flags)

	return logger, nil
}

// Close closes the logger and any open file handles
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// rotateIfNeeded checks if log rotation is needed and performs it
func (l *Logger) rotateIfNeeded() error {
	if l.logFile == nil || l.rotationConfig == nil {
		return nil
	}

	stat, err := l.logFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	// Check if rotation is needed based on size
	if stat.Size() >= l.rotationConfig.MaxSize {
		return l.rotateLog()
	}

	// Check if rotation is needed based on age
	if time.Since(stat.ModTime()) >= l.rotationConfig.MaxAge {
		return l.rotateLog()
	}

	return nil
}

// rotateLog performs log file rotation
func (l *Logger) rotateLog() error {
	if l.logFile == nil {
		return nil
	}

	// Close current log file
	if err := l.logFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	currentPath := l.logFile.Name()
	backupPath := fmt.Sprintf("%s.%s", currentPath, timestamp)

	// Rename current log file to backup
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// Compress backup if configured
	if l.rotationConfig.Compress {
		go l.compressLogFile(backupPath)
	}

	// Clean up old backups
	go l.cleanupOldBackups()

	// Create new log file
	file, err := utils.SafeOpenFile(currentPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	l.logFile = file

	// Update all loggers to use new file
	var writers []io.Writer
	writers = append(writers, os.Stdout, file)
	multiWriter := io.MultiWriter(writers...)

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	l.infoLogger = log.New(multiWriter, "", flags)
	l.warnLogger = log.New(multiWriter, "", flags)
	l.errLogger = log.New(multiWriter, "", flags)
	l.debugLogger = log.New(multiWriter, "", flags)
	l.fatalLogger = log.New(multiWriter, "", flags)

	return nil
}

// compressLogFile compresses a log file (placeholder implementation)
func (l *Logger) compressLogFile(filePath string) {
	// This is a placeholder - in a real implementation, you would use
	// a compression library like gzip to compress the file
	// For now, we'll just add a .gz extension to indicate it should be compressed
	compressedPath := filePath + ".gz"
	if err := os.Rename(filePath, compressedPath); err != nil {
		l.Error("Failed to compress log file: %v", err)
	}
}

// cleanupOldBackups removes old backup files based on configuration
func (l *Logger) cleanupOldBackups() {
	if l.logDir == "" || l.rotationConfig == nil {
		return
	}

	pattern := filepath.Join(l.logDir, "generator-*.log.*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		l.Error("Failed to find backup log files: %v", err)
		return
	}

	// Sort by modification time and keep only the most recent backups
	if len(matches) > l.rotationConfig.MaxBackups {
		// This is a simplified cleanup - in a real implementation,
		// you would sort by modification time and remove the oldest files
		for i := l.rotationConfig.MaxBackups; i < len(matches); i++ {
			if err := os.Remove(matches[i]); err != nil {
				l.Error("Failed to remove old backup log file %s: %v", matches[i], err)
			}
		}
	}
}

// GetLogDir returns the log directory path
func (l *Logger) GetLogDir() string {
	return l.logDir
}

// GetRecentEntries returns recent log entries for the logs command
func (l *Logger) GetRecentEntries(limit int) []interfaces.LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > len(l.entries) {
		limit = len(l.entries)
	}

	// Return the most recent entries
	start := len(l.entries) - limit
	if start < 0 {
		start = 0
	}

	result := make([]interfaces.LogEntry, limit)
	for i, entry := range l.entries[start:] {
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

// addEntry adds a log entry to the in-memory store
func (l *Logger) addEntry(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, entry)

	// Keep only the most recent entries
	if len(l.entries) > l.maxEntries {
		// Remove oldest entries
		copy(l.entries, l.entries[len(l.entries)-l.maxEntries:])
		l.entries = l.entries[:l.maxEntries]
	}
}

// Debug logs a debug message with structured format
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logWithLevel(LogLevelDebug, msg, nil, args...)
	}
}

// Info logs an info message with structured format
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logWithLevel(LogLevelInfo, msg, nil, args...)
	}
}

// Warn logs a warning message with structured format
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logWithLevel(LogLevelWarn, msg, nil, args...)
	}
}

// Error logs an error message with structured format
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.logWithLevel(LogLevelError, msg, nil, args...)
	}
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logWithLevel(LogLevelFatal, msg, nil, args...)
	os.Exit(1)
}

// logWithLevel provides consistent structured logging format
func (l *Logger) logWithLevel(level LogLevel, msg string, fields map[string]interface{}, args ...interface{}) {
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
	case LogLevelDebug:
		logger = l.debugLogger
	case LogLevelInfo:
		logger = l.infoLogger
	case LogLevelWarn:
		logger = l.warnLogger
	case LogLevelError:
		logger = l.errLogger
	case LogLevelFatal:
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
func (l *Logger) formatTextEntry(entry LogEntry) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", entry.Level))
	parts = append(parts, fmt.Sprintf("component=%s", entry.Component))
	parts = append(parts, fmt.Sprintf("message=\"%s\"", entry.Message))

	// Add fields
	if entry.Fields != nil {
		for k, v := range entry.Fields {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
	}

	// Add caller if present
	if entry.Caller != "" {
		parts = append(parts, fmt.Sprintf("caller=%s", entry.Caller))
	}

	// Add error if present
	if entry.Error != "" {
		parts = append(parts, fmt.Sprintf("error=\"%s\"", entry.Error))
	}

	return strings.Join(parts, " ")
}

// Structured logging methods for better consistency

// DebugWithFields logs a debug message with structured fields
func (l *Logger) DebugWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelDebug {
		l.logWithLevel(LogLevelDebug, msg, fields)
	}
}

// InfoWithFields logs an info message with structured fields
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelInfo {
		l.logWithLevel(LogLevelInfo, msg, fields)
	}
}

// WarnWithFields logs a warning message with structured fields
func (l *Logger) WarnWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelWarn {
		l.logWithLevel(LogLevelWarn, msg, fields)
	}
}

// ErrorWithFields logs an error message with structured fields
func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
	if l.level <= LogLevelError {
		l.logWithLevel(LogLevelError, msg, fields)
	}
}

// FatalWithFields logs a fatal message with structured fields and exits
func (l *Logger) FatalWithFields(msg string, fields map[string]interface{}) {
	l.logWithLevel(LogLevelFatal, msg, fields)
	os.Exit(1)
}

// ErrorWithError logs an error message with an error object
func (l *Logger) ErrorWithError(msg string, err error, fields map[string]interface{}) {
	if l.level > LogLevelError {
		return
	}

	if fields == nil {
		fields = make(map[string]interface{})
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     LogLevelError.String(),
		Component: l.component,
		Message:   msg,
		Fields:    fields,
		Error:     err.Error(),
	}

	// Add caller information if enabled
	if l.enableCaller {
		if _, file, line, ok := runtime.Caller(2); ok {
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

	l.errLogger.Println(logOutput)
}

// Operation logging methods for tracking key operations

// OperationContext holds context for tracking operations
type OperationContext struct {
	Operation string
	StartTime time.Time
	Fields    map[string]interface{}
}

// StartOperation begins tracking an operation and returns a context
func (l *Logger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	ctx := &interfaces.OperationContext{
		Operation: operation,
		StartTime: time.Now(),
		Fields:    fields,
	}

	logFields := make(map[string]interface{})
	for k, v := range fields {
		logFields[k] = v
	}
	logFields["operation"] = operation
	logFields["status"] = "started"
	logFields["start_time"] = ctx.StartTime.Format(time.RFC3339)

	l.InfoWithFields("Operation started", logFields)
	return ctx
}

// LogOperationStart logs the start of an operation (legacy method)
func (l *Logger) LogOperationStart(operation string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["operation"] = operation
	fields["status"] = "started"
	fields["start_time"] = time.Now().Format(time.RFC3339)
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
	fields["duration_human"] = duration.String()
	l.InfoWithFields("Operation completed successfully", fields)
}

// LogOperationError logs an error during an operation
func (l *Logger) LogOperationError(operation string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["operation"] = operation
	fields["status"] = "error"
	l.ErrorWithError("Operation failed", err, fields)
}

// FinishOperation completes an operation context with success
func (l *Logger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
	duration := time.Since(ctx.StartTime)

	fields := make(map[string]interface{})
	for k, v := range ctx.Fields {
		fields[k] = v
	}
	for k, v := range additionalFields {
		fields[k] = v
	}

	l.LogOperationSuccess(ctx.Operation, duration, fields)
}

// FinishOperationWithError completes an operation context with an error
func (l *Logger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
	duration := time.Since(ctx.StartTime)

	fields := make(map[string]interface{})
	for k, v := range ctx.Fields {
		fields[k] = v
	}
	for k, v := range additionalFields {
		fields[k] = v
	}
	fields["duration_ms"] = duration.Milliseconds()
	fields["duration_human"] = duration.String()

	l.LogOperationError(ctx.Operation, err, fields)
}

// LogPerformanceMetrics logs performance-related metrics
func (l *Logger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
	fields := make(map[string]interface{})
	fields["operation"] = operation
	fields["type"] = "performance_metrics"

	for k, v := range metrics {
		fields[k] = v
	}

	l.InfoWithFields("Performance metrics", fields)
}

// LogMemoryUsage logs current memory usage statistics
func (l *Logger) LogMemoryUsage(operation string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fields := map[string]interface{}{
		"operation":      operation,
		"type":           "memory_usage",
		"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
		"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
		"sys_mb":         float64(m.Sys) / 1024 / 1024,
		"num_gc":         m.NumGC,
		"goroutines":     runtime.NumGoroutine(),
	}

	l.DebugWithFields("Memory usage", fields)
}

// LogLevelFromString converts a string to LogLevel
func LogLevelFromString(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

// SetLevel updates the logger's log level (interface compatibility)
func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = LogLevel(level)
}

// GetLevel returns the current log level (interface compatibility)
func (l *Logger) GetLevel() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return int(l.level)
}

// SetLevelByLogLevel updates the logger's log level using LogLevel type
func (l *Logger) SetLevelByLogLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLogLevel returns the current log level as LogLevel type
func (l *Logger) GetLogLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetJSONOutput enables or disables JSON output format
func (l *Logger) SetJSONOutput(enable bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enableJSON = enable
}

// SetCallerInfo enables or disables caller information in logs
func (l *Logger) SetCallerInfo(enable bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enableCaller = enable
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= LogLevelDebug
}

// IsInfoEnabled returns true if info logging is enabled
func (l *Logger) IsInfoEnabled() bool {
	return l.level <= LogLevelInfo
}

// WithComponent creates a new logger instance with a different component name
func (l *Logger) WithComponent(component string) interfaces.Logger {
	config := &LoggerConfig{
		Level:          l.level,
		LogToFile:      l.logFile != nil,
		Component:      component,
		EnableJSON:     l.enableJSON,
		EnableCaller:   l.enableCaller,
		RotationConfig: l.rotationConfig,
		MaxEntries:     l.maxEntries,
	}

	newLogger, err := NewLoggerWithConfig(config)
	if err != nil {
		// Fallback to current logger if creation fails
		return l
	}

	return newLogger
}

// WithFields creates a temporary logger context with additional fields
type LoggerContext struct {
	logger *Logger
	fields map[string]interface{}
}

// WithFields returns a logger context with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &LoggerContext{
		logger: l,
		fields: fields,
	}
}

// Debug logs a debug message with the context fields
func (lc *LoggerContext) Debug(msg string, args ...interface{}) {
	lc.logger.logWithLevel(LogLevelDebug, fmt.Sprintf(msg, args...), lc.fields)
}

// Info logs an info message with the context fields
func (lc *LoggerContext) Info(msg string, args ...interface{}) {
	lc.logger.logWithLevel(LogLevelInfo, fmt.Sprintf(msg, args...), lc.fields)
}

// Warn logs a warning message with the context fields
func (lc *LoggerContext) Warn(msg string, args ...interface{}) {
	lc.logger.logWithLevel(LogLevelWarn, fmt.Sprintf(msg, args...), lc.fields)
}

// Error logs an error message with the context fields
func (lc *LoggerContext) Error(msg string, args ...interface{}) {
	lc.logger.logWithLevel(LogLevelError, fmt.Sprintf(msg, args...), lc.fields)
}

// ErrorWithError logs an error message with an error object and context fields
func (lc *LoggerContext) ErrorWithError(msg string, err error) {
	lc.logger.ErrorWithError(msg, err, lc.fields)
}

// GetLogFiles returns a list of log files in the log directory
func (l *Logger) GetLogFiles() ([]string, error) {
	if l.logDir == "" {
		return nil, fmt.Errorf("log directory not configured")
	}

	pattern := filepath.Join(l.logDir, "generator-*.log*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find log files: %w", err)
	}

	return matches, nil
}

// ReadLogFile reads the contents of a specific log file
func (l *Logger) ReadLogFile(filename string) ([]byte, error) {
	if l.logDir == "" {
		return nil, fmt.Errorf("log directory not configured")
	}

	// Clean the filename to prevent directory traversal
	cleanFilename := filepath.Base(filename)

	// Additional security check: only allow log files with expected patterns
	if !strings.HasPrefix(cleanFilename, "generator-") || !strings.HasSuffix(cleanFilename, ".log") {
		return nil, fmt.Errorf("invalid log file name: must match pattern generator-*.log")
	}

	// Ensure the file is in the log directory for security
	fullPath := filepath.Join(l.logDir, cleanFilename)

	// Resolve any symbolic links and ensure the path is within logDir
	resolvedPath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve log file path: %w", err)
	}

	resolvedLogDir, err := filepath.Abs(l.logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve log directory path: %w", err)
	}

	if !strings.HasPrefix(resolvedPath, resolvedLogDir+string(filepath.Separator)) {
		return nil, fmt.Errorf("invalid log file path: outside log directory")
	}

	// #nosec G304 - Path is validated above to prevent directory traversal
	return os.ReadFile(resolvedPath)
}

// FilterEntries filters log entries based on criteria
func (l *Logger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var filtered []interfaces.LogEntry
	for _, entry := range l.entries {
		// Filter by level
		if level != "" && !strings.EqualFold(entry.Level, level) {
			continue
		}

		// Filter by component
		if component != "" && !strings.Contains(strings.ToLower(entry.Component), strings.ToLower(component)) {
			continue
		}

		// Filter by time
		if !since.IsZero() && entry.Timestamp.Before(since) {
			continue
		}

		// Convert internal LogEntry to interface LogEntry
		interfaceEntry := interfaces.LogEntry{
			Timestamp: entry.Timestamp,
			Level:     entry.Level,
			Component: entry.Component,
			Message:   entry.Message,
			Fields:    entry.Fields,
			Caller:    entry.Caller,
			Error:     entry.Error,
		}
		filtered = append(filtered, interfaceEntry)
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}

	return filtered
}
