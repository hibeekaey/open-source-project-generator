// Package errors provides comprehensive error logging and reporting
package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ErrorLogger provides comprehensive error logging capabilities
type ErrorLogger struct {
	logFile    *os.File
	logPath    string
	mutex      sync.RWMutex
	level      LogLevel
	jsonFormat bool
	context    map[string]interface{}
}

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

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Error     *ErrorLogData          `json:"error,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Caller    *CallerInfo            `json:"caller,omitempty"`
}

// ErrorLogData contains detailed error information for logging
type ErrorLogData struct {
	Type        string                 `json:"type"`
	Code        int                    `json:"code"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Context     *ErrorContext          `json:"context,omitempty"`
	Cause       string                 `json:"cause,omitempty"`
	Stack       string                 `json:"stack,omitempty"`
	Severity    string                 `json:"severity"`
	Recoverable bool                   `json:"recoverable"`
}

// CallerInfo contains information about the caller
type CallerInfo struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// NewErrorLogger creates a new error logger
func NewErrorLogger(logPath string, level LogLevel, jsonFormat bool) (*ErrorLogger, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(logPath)
	if err := secureMkdirAll(logDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file with secure path validation
	logFile, err := secureOpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &ErrorLogger{
		logFile:    logFile,
		logPath:    logPath,
		level:      level,
		jsonFormat: jsonFormat,
		context:    make(map[string]interface{}),
	}, nil
}

// Close closes the error logger
func (el *ErrorLogger) Close() error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	if el.logFile != nil {
		return el.logFile.Close()
	}
	return nil
}

// SetLevel sets the logging level
func (el *ErrorLogger) SetLevel(level LogLevel) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.level = level
}

// SetJSONFormat sets whether to use JSON format
func (el *ErrorLogger) SetJSONFormat(jsonFormat bool) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.jsonFormat = jsonFormat
}

// SetContext sets global context for all log entries
func (el *ErrorLogger) SetContext(key string, value interface{}) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.context[key] = value
}

// LogError logs a CLI error with full context
func (el *ErrorLogger) LogError(err *CLIError) {
	if err == nil {
		return
	}

	// Determine log level based on error severity
	var logLevel LogLevel
	switch err.Severity {
	case SeverityLow:
		logLevel = LogLevelInfo
	case SeverityMedium:
		logLevel = LogLevelWarn
	case SeverityHigh:
		logLevel = LogLevelError
	case SeverityCritical:
		logLevel = LogLevelFatal
	default:
		logLevel = LogLevelError
	}

	// Create error log data
	errorData := &ErrorLogData{
		Type:        err.Type,
		Code:        err.Code,
		Message:     err.Message,
		Details:     err.Details,
		Suggestions: err.Suggestions,
		Context:     err.Context,
		Severity:    string(err.Severity),
		Recoverable: err.Recoverable,
	}

	if err.Cause != nil {
		errorData.Cause = err.Cause.Error()
	}

	if err.Stack != "" {
		errorData.Stack = err.Stack
	}

	// Log the error
	el.logWithLevel(logLevel, fmt.Sprintf("CLI Error: %s", err.Message), errorData, nil)
}

// LogRecoveryAttempt logs a recovery attempt
func (el *ErrorLogger) LogRecoveryAttempt(attempt *RecoveryAttempt) {
	context := map[string]interface{}{
		"strategy": attempt.Strategy,
		"duration": attempt.Duration.String(),
		"success":  attempt.Result.Success,
		"actions":  len(attempt.Result.Actions),
	}

	if attempt.Result.Success {
		el.logWithLevel(LogLevelInfo, fmt.Sprintf("Recovery successful: %s", attempt.Result.Message), nil, context)
	} else {
		el.logWithLevel(LogLevelWarn, fmt.Sprintf("Recovery failed: %s", attempt.Result.Message), nil, context)
	}
}

// Debug logs a debug message
func (el *ErrorLogger) Debug(message string, context map[string]interface{}) {
	el.logWithLevel(LogLevelDebug, message, nil, context)
}

// Info logs an info message
func (el *ErrorLogger) Info(message string, context map[string]interface{}) {
	el.logWithLevel(LogLevelInfo, message, nil, context)
}

// Warn logs a warning message
func (el *ErrorLogger) Warn(message string, context map[string]interface{}) {
	el.logWithLevel(LogLevelWarn, message, nil, context)
}

// Error logs an error message
func (el *ErrorLogger) Error(message string, context map[string]interface{}) {
	el.logWithLevel(LogLevelError, message, nil, context)
}

// Fatal logs a fatal message
func (el *ErrorLogger) Fatal(message string, context map[string]interface{}) {
	el.logWithLevel(LogLevelFatal, message, nil, context)
}

// logWithLevel logs a message at the specified level
func (el *ErrorLogger) logWithLevel(level LogLevel, message string, errorData *ErrorLogData, context map[string]interface{}) {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	// Check if we should log at this level
	if level < el.level {
		return
	}

	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   message,
		Error:     errorData,
		Context:   el.mergeContext(context),
	}

	// Add caller information for errors and above
	if level >= LogLevelError {
		file, line := GetCallerInfo(2)
		if file != "" {
			entry.Caller = &CallerInfo{
				File: file,
				Line: line,
			}
		}
	}

	// Write to log file
	el.writeLogEntry(&entry)

	// Also write to stderr for errors and above
	if level >= LogLevelError {
		el.writeToStderr(&entry)
	}
}

// mergeContext merges global context with entry-specific context
func (el *ErrorLogger) mergeContext(entryContext map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Add global context
	for k, v := range el.context {
		merged[k] = v
	}

	// Add entry-specific context (overwrites global)
	for k, v := range entryContext {
		merged[k] = v
	}

	return merged
}

// writeLogEntry writes a log entry to the log file
func (el *ErrorLogger) writeLogEntry(entry *LogEntry) {
	if el.logFile == nil {
		return
	}

	var output string

	if el.jsonFormat {
		if jsonData, err := json.Marshal(entry); err == nil {
			output = string(jsonData) + "\n"
		} else {
			// Fallback to simple format if JSON marshaling fails
			output = el.formatSimpleEntry(entry)
		}
	} else {
		output = el.formatSimpleEntry(entry)
	}

	if _, err := el.logFile.WriteString(output); err != nil {
		// Log to stderr as fallback if file write fails
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
	}
	if err := el.logFile.Sync(); err != nil {
		// Log to stderr as fallback if sync fails
		fmt.Fprintf(os.Stderr, "Failed to sync log file: %v\n", err)
	}
}

// writeToStderr writes a log entry to stderr
func (el *ErrorLogger) writeToStderr(entry *LogEntry) {
	output := el.formatSimpleEntry(entry)
	fmt.Fprint(os.Stderr, output)
}

// formatSimpleEntry formats a log entry in simple text format
func (el *ErrorLogger) formatSimpleEntry(entry *LogEntry) string {
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	var parts []string
	parts = append(parts, fmt.Sprintf("[%s] %s: %s", timestamp, entry.Level, entry.Message))

	// Add error details if present
	if entry.Error != nil {
		parts = append(parts, fmt.Sprintf("  Error Type: %s", entry.Error.Type))
		parts = append(parts, fmt.Sprintf("  Error Code: %d", entry.Error.Code))
		parts = append(parts, fmt.Sprintf("  Severity: %s", entry.Error.Severity))
		parts = append(parts, fmt.Sprintf("  Recoverable: %t", entry.Error.Recoverable))

		if entry.Error.Cause != "" {
			parts = append(parts, fmt.Sprintf("  Cause: %s", entry.Error.Cause))
		}

		if len(entry.Error.Details) > 0 {
			parts = append(parts, "  Details:")
			for k, v := range entry.Error.Details {
				parts = append(parts, fmt.Sprintf("    %s: %v", k, v))
			}
		}

		if len(entry.Error.Suggestions) > 0 {
			parts = append(parts, "  Suggestions:")
			for _, suggestion := range entry.Error.Suggestions {
				parts = append(parts, fmt.Sprintf("    - %s", suggestion))
			}
		}
	}

	// Add context if present
	if len(entry.Context) > 0 {
		parts = append(parts, "  Context:")
		for k, v := range entry.Context {
			parts = append(parts, fmt.Sprintf("    %s: %v", k, v))
		}
	}

	// Add caller info if present
	if entry.Caller != nil {
		parts = append(parts, fmt.Sprintf("  Caller: %s:%d", entry.Caller.File, entry.Caller.Line))
	}

	return strings.Join(parts, "\n") + "\n"
}

// ErrorReporter provides error reporting capabilities for debugging and support
type ErrorReporter struct {
	logger       *ErrorLogger
	reportDir    string
	maxReports   int
	reportFormat ReportFormat
}

// ReportFormat represents the format for error reports
type ReportFormat string

const (
	ReportFormatText ReportFormat = "text"
	ReportFormatJSON ReportFormat = "json"
	ReportFormatHTML ReportFormat = "html"
)

// ErrorReport represents a comprehensive error report
type ErrorReport struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Error       *CLIError              `json:"error"`
	Environment *EnvironmentInfo       `json:"environment"`
	System      *SystemInfo            `json:"system"`
	Recovery    *RecoveryAttempt       `json:"recovery,omitempty"`
	Logs        []LogEntry             `json:"logs,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// EnvironmentInfo contains environment information for reports
type EnvironmentInfo struct {
	GeneratorVersion string            `json:"generator_version"`
	GoVersion        string            `json:"go_version"`
	OS               string            `json:"os"`
	Arch             string            `json:"arch"`
	WorkingDir       string            `json:"working_dir"`
	Environment      map[string]string `json:"environment"`
	CI               *CIEnvironment    `json:"ci,omitempty"`
}

// SystemInfo contains system information for reports
type SystemInfo struct {
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	ProcessID   int    `json:"process_id"`
	ParentPID   int    `json:"parent_pid"`
	MemoryUsage int64  `json:"memory_usage"`
	DiskSpace   int64  `json:"disk_space"`
	CPUCount    int    `json:"cpu_count"`
}

// NewErrorReporter creates a new error reporter
func NewErrorReporter(logger *ErrorLogger, reportDir string, maxReports int, format ReportFormat) *ErrorReporter {
	return &ErrorReporter{
		logger:       logger,
		reportDir:    reportDir,
		maxReports:   maxReports,
		reportFormat: format,
	}
}

// GenerateReport generates a comprehensive error report
func (er *ErrorReporter) GenerateReport(err *CLIError, recovery *RecoveryAttempt, context map[string]interface{}) (*ErrorReport, error) {
	report := &ErrorReport{
		ID:          er.generateReportID(),
		Timestamp:   time.Now(),
		Error:       err,
		Environment: er.collectEnvironmentInfo(),
		System:      er.collectSystemInfo(),
		Recovery:    recovery,
		Context:     context,
	}

	// Collect recent logs
	if logs, logErr := er.collectRecentLogs(); logErr == nil {
		report.Logs = logs
	}

	return report, nil
}

// SaveReport saves an error report to disk
func (er *ErrorReporter) SaveReport(report *ErrorReport) (string, error) {
	// Ensure report directory exists
	if err := secureMkdirAll(er.reportDir, 0750); err != nil {
		return "", fmt.Errorf("failed to create report directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("error-report-%s.%s", report.ID, er.reportFormat)
	reportPath := filepath.Join(er.reportDir, filename)

	// Create report file with secure path validation
	file, err := secureCreateFile(reportPath)
	if err != nil {
		return "", fmt.Errorf("failed to create report file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close report file: %v\n", err)
		}
	}()

	// Write report based on format
	switch er.reportFormat {
	case ReportFormatJSON:
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			return "", fmt.Errorf("failed to encode JSON report: %w", err)
		}
	case ReportFormatHTML:
		if err := er.writeHTMLReport(file, report); err != nil {
			return "", fmt.Errorf("failed to write HTML report: %w", err)
		}
	default: // ReportFormatText
		if err := er.writeTextReport(file, report); err != nil {
			return "", fmt.Errorf("failed to write text report: %w", err)
		}
	}

	// Clean up old reports
	er.cleanupOldReports()

	return reportPath, nil
}

// generateReportID generates a unique report ID
func (er *ErrorReporter) generateReportID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// collectEnvironmentInfo collects environment information
func (er *ErrorReporter) collectEnvironmentInfo() *EnvironmentInfo {
	env := &EnvironmentInfo{
		Environment: make(map[string]string),
	}

	// Collect relevant environment variables
	relevantEnvVars := []string{
		"PATH", "HOME", "USER", "SHELL", "TERM",
		"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL",
		"GO_VERSION", "NODE_VERSION", "PYTHON_VERSION",
	}

	for _, envVar := range relevantEnvVars {
		if value := os.Getenv(envVar); value != "" {
			env.Environment[envVar] = value
		}
	}

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		env.WorkingDir = wd
	}

	// Detect CI environment
	collector := NewErrorContextCollector()
	env.CI = collector.detectCIEnvironment()

	return env
}

// collectSystemInfo collects system information
func (er *ErrorReporter) collectSystemInfo() *SystemInfo {
	info := &SystemInfo{
		ProcessID: os.Getpid(),
		ParentPID: os.Getppid(),
	}

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	// Get username
	if user := os.Getenv("USER"); user != "" {
		info.Username = user
	} else if user := os.Getenv("USERNAME"); user != "" {
		info.Username = user
	}

	return info
}

// collectRecentLogs collects recent log entries
func (er *ErrorReporter) collectRecentLogs() ([]LogEntry, error) {
	// This would read recent entries from the log file
	// For now, return empty slice
	return []LogEntry{}, nil
}

// writeTextReport writes a report in text format
func (er *ErrorReporter) writeTextReport(w io.Writer, report *ErrorReport) error {
	writeFunc := func(format string, args ...interface{}) error {
		_, err := fmt.Fprintf(w, format, args...)
		return err
	}

	if err := writeFunc("Error Report: %s\n", report.ID); err != nil {
		return err
	}
	if err := writeFunc("Generated: %s\n\n", report.Timestamp.Format(time.RFC3339)); err != nil {
		return err
	}

	// Error information
	if err := writeFunc("ERROR DETAILS\n"); err != nil {
		return err
	}
	if err := writeFunc("============\n"); err != nil {
		return err
	}
	if err := writeFunc("Type: %s\n", report.Error.Type); err != nil {
		return err
	}
	if err := writeFunc("Code: %d\n", report.Error.Code); err != nil {
		return err
	}
	if err := writeFunc("Message: %s\n", report.Error.Message); err != nil {
		return err
	}
	if err := writeFunc("Severity: %s\n", report.Error.Severity); err != nil {
		return err
	}
	if err := writeFunc("Recoverable: %t\n", report.Error.Recoverable); err != nil {
		return err
	}

	if report.Error.Cause != nil {
		if err := writeFunc("Cause: %s\n", report.Error.Cause.Error()); err != nil {
			return err
		}
	}

	if len(report.Error.Details) > 0 {
		if err := writeFunc("\nDetails:\n"); err != nil {
			return err
		}
		for k, v := range report.Error.Details {
			if err := writeFunc("  %s: %v\n", k, v); err != nil {
				return err
			}
		}
	}

	if len(report.Error.Suggestions) > 0 {
		if err := writeFunc("\nSuggestions:\n"); err != nil {
			return err
		}
		for _, suggestion := range report.Error.Suggestions {
			if err := writeFunc("  - %s\n", suggestion); err != nil {
				return err
			}
		}
	}

	// Environment information
	if err := writeFunc("\nENVIRONMENT\n"); err != nil {
		return err
	}
	if err := writeFunc("===========\n"); err != nil {
		return err
	}
	if err := writeFunc("Generator Version: %s\n", report.Environment.GeneratorVersion); err != nil {
		return err
	}
	if err := writeFunc("Go Version: %s\n", report.Environment.GoVersion); err != nil {
		return err
	}
	if err := writeFunc("OS: %s\n", report.Environment.OS); err != nil {
		return err
	}
	if err := writeFunc("Architecture: %s\n", report.Environment.Arch); err != nil {
		return err
	}
	if err := writeFunc("Working Directory: %s\n", report.Environment.WorkingDir); err != nil {
		return err
	}

	if report.Environment.CI != nil && report.Environment.CI.IsCI {
		if err := writeFunc("CI Provider: %s\n", report.Environment.CI.Provider); err != nil {
			return err
		}
	}

	// System information
	if err := writeFunc("\nSYSTEM\n"); err != nil {
		return err
	}
	if err := writeFunc("======\n"); err != nil {
		return err
	}
	if err := writeFunc("Hostname: %s\n", report.System.Hostname); err != nil {
		return err
	}
	if err := writeFunc("Username: %s\n", report.System.Username); err != nil {
		return err
	}
	if err := writeFunc("Process ID: %d\n", report.System.ProcessID); err != nil {
		return err
	}

	// Recovery information
	if report.Recovery != nil {
		if err := writeFunc("\nRECOVERY ATTEMPT\n"); err != nil {
			return err
		}
		if err := writeFunc("================\n"); err != nil {
			return err
		}
		if err := writeFunc("Strategy: %s\n", report.Recovery.Strategy); err != nil {
			return err
		}
		if err := writeFunc("Success: %t\n", report.Recovery.Result.Success); err != nil {
			return err
		}
		if err := writeFunc("Message: %s\n", report.Recovery.Result.Message); err != nil {
			return err
		}
		if err := writeFunc("Duration: %s\n", report.Recovery.Duration); err != nil {
			return err
		}
	}

	return nil
}

// writeHTMLReport writes a report in HTML format
func (er *ErrorReporter) writeHTMLReport(w io.Writer, report *ErrorReport) error {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Error Report: %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .section { margin-bottom: 20px; }
        .section h2 { color: #333; border-bottom: 2px solid #ccc; }
        .error { background-color: #ffe6e6; padding: 10px; border-left: 4px solid #ff0000; }
        .details { background-color: #f5f5f5; padding: 10px; }
        .suggestions { background-color: #e6f3ff; padding: 10px; }
        pre { background-color: #f0f0f0; padding: 10px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>Error Report: %s</h1>
    <p><strong>Generated:</strong> %s</p>
    
    <div class="section">
        <h2>Error Details</h2>
        <div class="error">
            <p><strong>Type:</strong> %s</p>
            <p><strong>Code:</strong> %d</p>
            <p><strong>Message:</strong> %s</p>
            <p><strong>Severity:</strong> %s</p>
            <p><strong>Recoverable:</strong> %t</p>
        </div>
    </div>
</body>
</html>`

	_, err := fmt.Fprintf(w, html,
		report.ID, report.ID,
		report.Timestamp.Format(time.RFC3339),
		report.Error.Type,
		report.Error.Code,
		report.Error.Message,
		report.Error.Severity,
		report.Error.Recoverable,
	)

	return err
}

// cleanupOldReports removes old report files to maintain the maximum count
func (er *ErrorReporter) cleanupOldReports() {
	// This would implement cleanup logic to maintain maxReports
	// For now, it's a placeholder
}

// GetReportPath returns the path where reports are stored
func (er *ErrorReporter) GetReportPath() string {
	return er.reportDir
}

// ListReports returns a list of available error reports
func (er *ErrorReporter) ListReports() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(er.reportDir, "error-report-*"))
	if err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}

	return files, nil
}
