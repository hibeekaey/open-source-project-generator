// Package errors provides enhanced logging with structured output and diagnostics
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

// EnhancedLogger provides comprehensive logging with structured output and diagnostics
type EnhancedLogger struct {
	baseLogger      *ErrorLogger
	diagnosticMode  bool
	traceMode       bool
	performanceMode bool
	structuredMode  bool
	config          *EnhancedLoggingConfig
	writers         []io.Writer
	mutex           sync.RWMutex
	sessionID       string
	startTime       time.Time
	operations      map[string]*OperationContext
	metrics         *LoggingMetrics
}

// EnhancedLoggingConfig contains configuration for enhanced logging
type EnhancedLoggingConfig struct {
	Level             LogLevel               `json:"level"`
	Format            LogFormat              `json:"format"`
	OutputPaths       []string               `json:"output_paths"`
	EnableDiagnostics bool                   `json:"enable_diagnostics"`
	EnableTracing     bool                   `json:"enable_tracing"`
	EnablePerformance bool                   `json:"enable_performance"`
	EnableStructured  bool                   `json:"enable_structured"`
	BufferSize        int                    `json:"buffer_size"`
	FlushInterval     time.Duration          `json:"flush_interval"`
	MaxFileSize       int64                  `json:"max_file_size"`
	MaxBackups        int                    `json:"max_backups"`
	Compress          bool                   `json:"compress"`
	Fields            map[string]interface{} `json:"fields"`
}

// LogFormat represents different log output formats
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
	LogFormatCSV  LogFormat = "csv"
)

// OperationContext tracks operation-specific logging context
type OperationContext struct {
	ID        string                 `json:"id"`
	Operation string                 `json:"operation"`
	StartTime time.Time              `json:"start_time"`
	Context   map[string]interface{} `json:"context"`
	Events    []LogEvent             `json:"events"`
	Metrics   map[string]interface{} `json:"metrics"`
	mutex     sync.RWMutex
}

// LogEvent represents a single log event within an operation
type LogEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    *CallerInfo            `json:"caller,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
}

// LoggingMetrics tracks logging performance and statistics
type LoggingMetrics struct {
	TotalLogs       int64              `json:"total_logs"`
	LogsByLevel     map[LogLevel]int64 `json:"logs_by_level"`
	LogsByOperation map[string]int64   `json:"logs_by_operation"`
	AverageLatency  time.Duration      `json:"average_latency"`
	ErrorCount      int64              `json:"error_count"`
	DroppedLogs     int64              `json:"dropped_logs"`
	LastFlush       time.Time          `json:"last_flush"`
	BufferUsage     int                `json:"buffer_usage"`
}

// DefaultEnhancedLoggingConfig returns default configuration for enhanced logging
func DefaultEnhancedLoggingConfig() *EnhancedLoggingConfig {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".generator", "logs")

	return &EnhancedLoggingConfig{
		Level:             LogLevelInfo,
		Format:            LogFormatText,
		OutputPaths:       []string{filepath.Join(logDir, "enhanced.log")},
		EnableDiagnostics: true,
		EnableTracing:     false,
		EnablePerformance: true,
		EnableStructured:  false,
		BufferSize:        1000,
		FlushInterval:     5 * time.Second,
		MaxFileSize:       100 * 1024 * 1024, // 100MB
		MaxBackups:        5,
		Compress:          true,
		Fields:            make(map[string]interface{}),
	}
}

// NewEnhancedLogger creates a new enhanced logger
func NewEnhancedLogger(config *EnhancedLoggingConfig) (*EnhancedLogger, error) {
	if config == nil {
		config = DefaultEnhancedLoggingConfig()
	}

	// Create base logger
	baseLogger, err := NewErrorLogger(config.OutputPaths[0], config.Level, config.Format == LogFormatJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create base logger: %w", err)
	}

	// Generate session ID
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	logger := &EnhancedLogger{
		baseLogger:      baseLogger,
		diagnosticMode:  config.EnableDiagnostics,
		traceMode:       config.EnableTracing,
		performanceMode: config.EnablePerformance,
		structuredMode:  config.EnableStructured,
		config:          config,
		sessionID:       sessionID,
		startTime:       time.Now(),
		operations:      make(map[string]*OperationContext),
		metrics: &LoggingMetrics{
			LogsByLevel:     make(map[LogLevel]int64),
			LogsByOperation: make(map[string]int64),
		},
	}

	// Initialize writers
	if err := logger.initializeWriters(); err != nil {
		return nil, fmt.Errorf("failed to initialize writers: %w", err)
	}

	// Start background flusher if buffering is enabled
	if config.BufferSize > 0 {
		go logger.backgroundFlusher()
	}

	// Log session start
	logger.logSessionStart()

	return logger, nil
}

// initializeWriters initializes log writers for different output paths
func (el *EnhancedLogger) initializeWriters() error {
	el.writers = make([]io.Writer, 0, len(el.config.OutputPaths))

	for _, path := range el.config.OutputPaths {
		switch path {
		case "stdout":
			el.writers = append(el.writers, os.Stdout)
		case "stderr":
			el.writers = append(el.writers, os.Stderr)
		default:
			// Ensure directory exists
			dir := filepath.Dir(path)
			if err := secureMkdirAll(dir, 0750); err != nil {
				return fmt.Errorf("failed to create log directory %s: %w", dir, err)
			}

			// Open file
			file, err := secureOpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if err != nil {
				return fmt.Errorf("failed to open log file %s: %w", path, err)
			}

			el.writers = append(el.writers, file)
		}
	}

	return nil
}

// logSessionStart logs the start of a logging session
func (el *EnhancedLogger) logSessionStart() {
	fields := map[string]interface{}{
		"session_id": el.sessionID,
		"pid":        os.Getpid(),
		"version":    "enhanced-logger-1.0",
	}

	// Add global fields
	for k, v := range el.config.Fields {
		fields[k] = v
	}

	el.logWithFields(LogLevelInfo, "Enhanced logging session started", fields)
}

// StartOperation starts tracking a new operation
func (el *EnhancedLogger) StartOperation(operation string, context map[string]interface{}) *OperationContext {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	opID := fmt.Sprintf("%s_%d", operation, time.Now().UnixNano())

	opCtx := &OperationContext{
		ID:        opID,
		Operation: operation,
		StartTime: time.Now(),
		Context:   context,
		Events:    make([]LogEvent, 0),
		Metrics:   make(map[string]interface{}),
	}

	el.operations[opID] = opCtx

	// Log operation start
	fields := map[string]interface{}{
		"operation_id": opID,
		"operation":    operation,
	}
	for k, v := range context {
		fields[k] = v
	}

	el.logWithFields(LogLevelInfo, fmt.Sprintf("Operation started: %s", operation), fields)

	return opCtx
}

// FinishOperation finishes tracking an operation
func (el *EnhancedLogger) FinishOperation(opCtx *OperationContext, result map[string]interface{}) {
	if opCtx == nil {
		return
	}

	el.mutex.Lock()
	defer el.mutex.Unlock()

	duration := time.Since(opCtx.StartTime)

	// Update operation metrics
	opCtx.mutex.Lock()
	opCtx.Metrics["duration"] = duration
	opCtx.Metrics["end_time"] = time.Now()
	for k, v := range result {
		opCtx.Metrics[k] = v
	}
	opCtx.mutex.Unlock()

	// Log operation completion
	fields := map[string]interface{}{
		"operation_id": opCtx.ID,
		"operation":    opCtx.Operation,
		"duration":     duration,
		"duration_ms":  duration.Milliseconds(),
	}
	for k, v := range result {
		fields[k] = v
	}

	el.logWithFields(LogLevelInfo, fmt.Sprintf("Operation completed: %s", opCtx.Operation), fields)

	// Update metrics
	el.metrics.LogsByOperation[opCtx.Operation]++

	// Clean up operation context
	delete(el.operations, opCtx.ID)
}

// FinishOperationWithError finishes an operation with an error
func (el *EnhancedLogger) FinishOperationWithError(opCtx *OperationContext, err error, context map[string]interface{}) {
	if opCtx == nil {
		return
	}

	result := map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	}
	for k, v := range context {
		result[k] = v
	}

	el.FinishOperation(opCtx, result)
}

// LogPerformanceMetrics logs performance metrics for an operation
func (el *EnhancedLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
	if !el.performanceMode {
		return
	}

	fields := map[string]interface{}{
		"operation": operation,
		"type":      "performance_metrics",
	}
	for k, v := range metrics {
		fields[k] = v
	}

	el.logWithFields(LogLevelDebug, fmt.Sprintf("Performance metrics for %s", operation), fields)
}

// LogDiagnostic logs diagnostic information
func (el *EnhancedLogger) LogDiagnostic(component string, diagnostic map[string]interface{}) {
	if !el.diagnosticMode {
		return
	}

	fields := map[string]interface{}{
		"component": component,
		"type":      "diagnostic",
	}
	for k, v := range diagnostic {
		fields[k] = v
	}

	el.logWithFields(LogLevelDebug, fmt.Sprintf("Diagnostic info for %s", component), fields)
}

// LogTrace logs trace information for debugging
func (el *EnhancedLogger) LogTrace(operation string, step string, context map[string]interface{}) {
	if !el.traceMode {
		return
	}

	fields := map[string]interface{}{
		"operation": operation,
		"step":      step,
		"type":      "trace",
	}
	for k, v := range context {
		fields[k] = v
	}

	// Add caller information
	if file, line := GetCallerInfo(1); file != "" {
		fields["caller_file"] = file
		fields["caller_line"] = line
	}

	el.logWithFields(LogLevelDebug, fmt.Sprintf("Trace: %s - %s", operation, step), fields)
}

// Debug logs a debug message with enhanced context
func (el *EnhancedLogger) Debug(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	el.logWithFields(LogLevelDebug, message, nil)
}

// DebugWithFields logs a debug message with structured fields
func (el *EnhancedLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	el.logWithFields(LogLevelDebug, msg, fields)
}

// Info logs an info message with enhanced context
func (el *EnhancedLogger) Info(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	el.logWithFields(LogLevelInfo, message, nil)
}

// InfoWithFields logs an info message with structured fields
func (el *EnhancedLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	el.logWithFields(LogLevelInfo, msg, fields)
}

// Warn logs a warning message with enhanced context
func (el *EnhancedLogger) Warn(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	el.logWithFields(LogLevelWarn, message, nil)
}

// WarnWithFields logs a warning message with structured fields
func (el *EnhancedLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	el.logWithFields(LogLevelWarn, msg, fields)
}

// Error logs an error message with enhanced context
func (el *EnhancedLogger) Error(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	el.logWithFields(LogLevelError, message, nil)
	el.metrics.ErrorCount++
}

// ErrorWithFields logs an error message with structured fields
func (el *EnhancedLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	el.logWithFields(LogLevelError, msg, fields)
	el.metrics.ErrorCount++
}

// Fatal logs a fatal message with enhanced context
func (el *EnhancedLogger) Fatal(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	el.logWithFields(LogLevelFatal, message, nil)
}

// FatalWithFields logs a fatal message with structured fields
func (el *EnhancedLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	el.logWithFields(LogLevelFatal, msg, fields)
}

// logWithFields logs a message with structured fields
func (el *EnhancedLogger) logWithFields(level LogLevel, message string, fields map[string]interface{}) {
	startTime := time.Now()

	// Check log level
	if level < el.config.Level {
		return
	}

	// Create log event
	event := LogEvent{
		Timestamp: startTime,
		Level:     level,
		Message:   message,
		Fields:    make(map[string]interface{}),
	}

	// Add global fields
	for k, v := range el.config.Fields {
		event.Fields[k] = v
	}

	// Add session fields
	event.Fields["session_id"] = el.sessionID
	event.Fields["session_duration"] = time.Since(el.startTime)

	// Add provided fields
	for k, v := range fields {
		event.Fields[k] = v
	}

	// Add caller information if enabled
	if el.traceMode || level >= LogLevelError {
		if file, line := GetCallerInfo(2); file != "" {
			event.Caller = &CallerInfo{
				File: file,
				Line: line,
			}
		}
	}

	// Write log entry
	el.writeLogEvent(&event)

	// Update metrics
	el.mutex.Lock()
	el.metrics.TotalLogs++
	el.metrics.LogsByLevel[level]++

	// Update average latency
	latency := time.Since(startTime)
	if el.metrics.TotalLogs == 1 {
		el.metrics.AverageLatency = latency
	} else {
		el.metrics.AverageLatency = (el.metrics.AverageLatency*time.Duration(el.metrics.TotalLogs-1) + latency) / time.Duration(el.metrics.TotalLogs)
	}
	el.mutex.Unlock()
}

// writeLogEvent writes a log event to all configured writers
func (el *EnhancedLogger) writeLogEvent(event *LogEvent) {
	var output string

	switch el.config.Format {
	case LogFormatJSON:
		output = el.formatJSONEvent(event)
	case LogFormatCSV:
		output = el.formatCSVEvent(event)
	default:
		output = el.formatTextEvent(event)
	}

	// Write to all writers
	for _, writer := range el.writers {
		if _, err := writer.Write([]byte(output + "\n")); err != nil {
			// Log write error to stderr as fallback
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		}
	}
}

// formatTextEvent formats a log event as text
func (el *EnhancedLogger) formatTextEvent(event *LogEvent) string {
	var parts []string

	// Timestamp and level
	timestamp := event.Timestamp.Format("2006-01-02 15:04:05.000")
	parts = append(parts, fmt.Sprintf("[%s] %s: %s", timestamp, event.Level.String(), event.Message))

	// Add caller information
	if event.Caller != nil {
		parts = append(parts, fmt.Sprintf("  Caller: %s:%d", event.Caller.File, event.Caller.Line))
	}

	// Add fields
	if len(event.Fields) > 0 {
		parts = append(parts, "  Fields:")
		for k, v := range event.Fields {
			parts = append(parts, fmt.Sprintf("    %s: %v", k, v))
		}
	}

	return strings.Join(parts, "\n")
}

// formatJSONEvent formats a log event as JSON
func (el *EnhancedLogger) formatJSONEvent(event *LogEvent) string {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal log event: %v"}`, err)
	}
	return string(jsonData)
}

// formatCSVEvent formats a log event as CSV
func (el *EnhancedLogger) formatCSVEvent(event *LogEvent) string {
	// Simple CSV format: timestamp,level,message,fields_json
	fieldsJSON, _ := json.Marshal(event.Fields)
	return fmt.Sprintf("%s,%s,%q,%s",
		event.Timestamp.Format("2006-01-02 15:04:05.000"),
		event.Level.String(),
		event.Message,
		string(fieldsJSON))
}

// backgroundFlusher periodically flushes buffered logs
func (el *EnhancedLogger) backgroundFlusher() {
	ticker := time.NewTicker(el.config.FlushInterval)
	defer ticker.Stop()

	for range ticker.C {
		el.flush()
	}
}

// flush flushes any buffered logs
func (el *EnhancedLogger) flush() {
	el.mutex.Lock()
	el.metrics.LastFlush = time.Now()
	el.mutex.Unlock()

	// Flush file writers
	for _, writer := range el.writers {
		if file, ok := writer.(*os.File); ok {
			if err := file.Sync(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to flush log file: %v\n", err)
			}
		}
	}
}

// SetLevel sets the logging level
func (el *EnhancedLogger) SetLevel(level int) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.config.Level = LogLevel(level)
}

// SetJSONOutput enables or disables JSON output
func (el *EnhancedLogger) SetJSONOutput(enable bool) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	if enable {
		el.config.Format = LogFormatJSON
	} else {
		el.config.Format = LogFormatText
	}
}

// SetCallerInfo enables or disables caller information
func (el *EnhancedLogger) SetCallerInfo(enable bool) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.traceMode = enable
}

// IsDebugEnabled returns whether debug logging is enabled
func (el *EnhancedLogger) IsDebugEnabled() bool {
	return el.config.Level <= LogLevelDebug
}

// IsInfoEnabled returns whether info logging is enabled
func (el *EnhancedLogger) IsInfoEnabled() bool {
	return el.config.Level <= LogLevelInfo
}

// GetMetrics returns current logging metrics
func (el *EnhancedLogger) GetMetrics() *LoggingMetrics {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	// Return a copy
	metrics := *el.metrics
	return &metrics
}

// GetActiveOperations returns currently active operations
func (el *EnhancedLogger) GetActiveOperations() map[string]*OperationContext {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	// Return a copy
	operations := make(map[string]*OperationContext)
	for k, v := range el.operations {
		operations[k] = v
	}
	return operations
}

// GenerateReport generates a comprehensive logging report
func (el *EnhancedLogger) GenerateReport() *LoggingReport {
	metrics := el.GetMetrics()
	operations := el.GetActiveOperations()

	return &LoggingReport{
		GeneratedAt:      time.Now(),
		SessionID:        el.sessionID,
		SessionDuration:  time.Since(el.startTime),
		Metrics:          metrics,
		ActiveOperations: len(operations),
		Configuration:    el.config,
		Recommendations:  el.generateRecommendations(metrics),
	}
}

// LoggingReport contains comprehensive logging information
type LoggingReport struct {
	GeneratedAt      time.Time              `json:"generated_at"`
	SessionID        string                 `json:"session_id"`
	SessionDuration  time.Duration          `json:"session_duration"`
	Metrics          *LoggingMetrics        `json:"metrics"`
	ActiveOperations int                    `json:"active_operations"`
	Configuration    *EnhancedLoggingConfig `json:"configuration"`
	Recommendations  []string               `json:"recommendations"`
}

// generateRecommendations generates logging recommendations based on metrics
func (el *EnhancedLogger) generateRecommendations(metrics *LoggingMetrics) []string {
	var recommendations []string

	// Check error rate
	if metrics.TotalLogs > 0 {
		errorRate := float64(metrics.ErrorCount) / float64(metrics.TotalLogs) * 100
		if errorRate > 10 {
			recommendations = append(recommendations, "High error rate detected - investigate error patterns")
		}
	}

	// Check average latency
	if metrics.AverageLatency > 10*time.Millisecond {
		recommendations = append(recommendations, "High logging latency - consider async logging or reducing log volume")
	}

	// Check dropped logs
	if metrics.DroppedLogs > 0 {
		recommendations = append(recommendations, "Dropped logs detected - increase buffer size or reduce log volume")
	}

	// Check log volume
	if metrics.TotalLogs > 10000 {
		recommendations = append(recommendations, "High log volume - consider increasing log level or adding filters")
	}

	return recommendations
}

// Close closes the enhanced logger and releases resources
func (el *EnhancedLogger) Close() error {
	// Flush any remaining logs
	el.flush()

	// Close file writers
	for _, writer := range el.writers {
		if file, ok := writer.(*os.File); ok && file != os.Stdout && file != os.Stderr {
			if err := file.Close(); err != nil {
				return err
			}
		}
	}

	// Close base logger
	if el.baseLogger != nil {
		return el.baseLogger.Close()
	}

	return nil
}
