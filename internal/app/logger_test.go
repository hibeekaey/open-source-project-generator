package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevelFatal, "FATAL"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultLogRotationConfig(t *testing.T) {
	config := DefaultLogRotationConfig()
	if config == nil {
		t.Fatal("DefaultLogRotationConfig returned nil")
	}

	if config.MaxSize != 10*1024*1024 {
		t.Errorf("Expected MaxSize 10MB, got %d", config.MaxSize)
	}

	if config.MaxAge != 7*24*time.Hour {
		t.Errorf("Expected MaxAge 7 days, got %v", config.MaxAge)
	}

	if config.MaxBackups != 5 {
		t.Errorf("Expected MaxBackups 5, got %d", config.MaxBackups)
	}

	if !config.Compress {
		t.Error("Expected Compress to be true")
	}
}

func TestDefaultLoggerConfig(t *testing.T) {
	config := DefaultLoggerConfig()
	if config == nil {
		t.Fatal("DefaultLoggerConfig returned nil")
	}

	if config.Level != LogLevelInfo {
		t.Errorf("Expected Level Info, got %v", config.Level)
	}

	if !config.LogToFile {
		t.Error("Expected LogToFile to be true")
	}

	if config.Component != "app" {
		t.Errorf("Expected Component 'app', got '%s'", config.Component)
	}

	if config.EnableJSON {
		t.Error("Expected EnableJSON to be false")
	}

	if config.EnableCaller {
		t.Error("Expected EnableCaller to be false")
	}

	if config.MaxEntries != 1000 {
		t.Errorf("Expected MaxEntries 1000, got %d", config.MaxEntries)
	}
}

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if logger.level != LogLevelInfo {
		t.Errorf("Expected level Info, got %v", logger.level)
	}

	if logger.logFile != nil {
		t.Error("Expected logFile to be nil when logToFile is false")
	}
}

func TestNewLoggerWithComponent(t *testing.T) {
	logger, err := NewLoggerWithComponent(LogLevelDebug, false, "test-component")
	if err != nil {
		t.Fatalf("NewLoggerWithComponent failed: %v", err)
	}

	if logger.component != "test-component" {
		t.Errorf("Expected component 'test-component', got '%s'", logger.component)
	}

	if logger.level != LogLevelDebug {
		t.Errorf("Expected level Debug, got %v", logger.level)
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	config := &LoggerConfig{
		Level:        LogLevelWarn,
		LogToFile:    false,
		Component:    "test",
		EnableJSON:   true,
		EnableCaller: true,
		MaxEntries:   500,
	}

	logger, err := NewLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}

	if logger.level != LogLevelWarn {
		t.Errorf("Expected level Warn, got %v", logger.level)
	}

	if logger.component != "test" {
		t.Errorf("Expected component 'test', got '%s'", logger.component)
	}

	if !logger.enableJSON {
		t.Error("Expected enableJSON to be true")
	}

	if !logger.enableCaller {
		t.Error("Expected enableCaller to be true")
	}

	if logger.maxEntries != 500 {
		t.Errorf("Expected maxEntries 500, got %d", logger.maxEntries)
	}
}

func TestLogger_Close(t *testing.T) {
	// Test with file logging disabled
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	err = logger.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Test with file logging enabled (creates temporary file)
	tempDir := t.TempDir()
	config := &LoggerConfig{
		Level:     LogLevelInfo,
		LogToFile: true,
		Component: "test",
	}

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger, err = NewLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}

	err = logger.Close()
	if err != nil {
		t.Errorf("Close() with file failed: %v", err)
	}
}

func TestLogger_LoggingMethods(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Test basic logging methods
	logger.Debug("Debug message %s", "test")
	logger.Info("Info message %s", "test")
	logger.Warn("Warn message %s", "test")
	logger.Error("Error message %s", "test")

	// Test logging with fields
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	logger.DebugWithFields("Debug with fields", fields)
	logger.InfoWithFields("Info with fields", fields)
	logger.WarnWithFields("Warn with fields", fields)
	logger.ErrorWithFields("Error with fields", fields)

	// Test error logging with error object
	testErr := errors.New("test error")
	logger.ErrorWithError("Error with error object", testErr, fields)
}

func TestLogger_LogLevels(t *testing.T) {
	// Test that messages below the log level are not logged
	logger, err := NewLogger(LogLevelError, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// These should not be logged (level too low)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")

	// This should be logged
	logger.Error("Error message")
}

func TestLogger_OperationLogging(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Test operation context
	fields := map[string]interface{}{"test": "value"}
	ctx := logger.StartOperation("test-operation", fields)

	if ctx.Operation != "test-operation" {
		t.Errorf("Expected operation 'test-operation', got '%s'", ctx.Operation)
	}

	if ctx.Fields["test"] != "value" {
		t.Error("Operation context fields not set correctly")
	}

	// Test operation completion
	logger.FinishOperation(ctx, map[string]interface{}{"result": "success"})

	// Test operation with error
	testErr := errors.New("operation failed")
	logger.FinishOperationWithError(ctx, testErr, map[string]interface{}{"error_code": 500})

	// Test individual operation logging methods
	logger.LogOperationStart("start-test", fields)
	logger.LogOperationSuccess("success-test", time.Millisecond*100, fields)
	logger.LogOperationError("error-test", testErr, fields)
}

func TestLogger_PerformanceLogging(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Test performance metrics logging
	metrics := map[string]interface{}{
		"duration_ms": 150,
		"memory_mb":   64,
		"cpu_percent": 25.5,
	}
	logger.LogPerformanceMetrics("test-operation", metrics)

	// Test memory usage logging
	logger.LogMemoryUsage("memory-test")
}

func TestLogger_LevelManagement(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Test level getters/setters
	if logger.GetLevel() != int(LogLevelInfo) {
		t.Errorf("Expected level %d, got %d", int(LogLevelInfo), logger.GetLevel())
	}

	logger.SetLevel(int(LogLevelDebug))
	if logger.GetLevel() != int(LogLevelDebug) {
		t.Errorf("Expected level %d after SetLevel, got %d", int(LogLevelDebug), logger.GetLevel())
	}

	// Test LogLevel type methods
	if logger.GetLogLevel() != LogLevelDebug {
		t.Errorf("Expected LogLevel %v, got %v", LogLevelDebug, logger.GetLogLevel())
	}

	logger.SetLevelByLogLevel(LogLevelWarn)
	if logger.GetLogLevel() != LogLevelWarn {
		t.Errorf("Expected LogLevel %v after SetLevelByLogLevel, got %v", LogLevelWarn, logger.GetLogLevel())
	}

	// Test level check methods
	logger.SetLevelByLogLevel(LogLevelDebug)
	if !logger.IsDebugEnabled() {
		t.Error("IsDebugEnabled should return true for Debug level")
	}
	if !logger.IsInfoEnabled() {
		t.Error("IsInfoEnabled should return true for Debug level")
	}

	logger.SetLevelByLogLevel(LogLevelError)
	if logger.IsDebugEnabled() {
		t.Error("IsDebugEnabled should return false for Error level")
	}
	if logger.IsInfoEnabled() {
		t.Error("IsInfoEnabled should return false for Error level")
	}
}

func TestLogger_Configuration(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Test JSON output configuration
	logger.SetJSONOutput(true)
	if !logger.enableJSON {
		t.Error("SetJSONOutput(true) should enable JSON output")
	}

	logger.SetJSONOutput(false)
	if logger.enableJSON {
		t.Error("SetJSONOutput(false) should disable JSON output")
	}

	// Test caller info configuration
	logger.SetCallerInfo(true)
	if !logger.enableCaller {
		t.Error("SetCallerInfo(true) should enable caller info")
	}

	logger.SetCallerInfo(false)
	if logger.enableCaller {
		t.Error("SetCallerInfo(false) should disable caller info")
	}
}

func TestLogger_WithComponent(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	newLogger := logger.WithComponent("new-component")
	if newLogger == nil {
		t.Fatal("WithComponent returned nil")
	}

	// Should return a different logger instance
	if newLogger == logger {
		t.Error("WithComponent should return a new logger instance")
	}
}

func TestLogger_WithFields(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	ctx := logger.WithFields(fields)
	if ctx == nil {
		t.Fatal("WithFields returned nil")
	}

	// Test logging with context
	ctx.Debug("Debug message")
	ctx.Info("Info message")
	ctx.Warn("Warn message")
	ctx.Error("Error message")

	testErr := errors.New("test error")
	ctx.ErrorWithError("Error with error", testErr)
}

func TestLogger_GetRecentEntries(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Log some entries
	logger.Info("Entry 1")
	logger.Warn("Entry 2")
	logger.Error("Entry 3")

	// Get recent entries
	entries := logger.GetRecentEntries(2)
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Test with limit larger than available entries
	entries = logger.GetRecentEntries(10)
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Test with zero limit
	entries = logger.GetRecentEntries(0)
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries with zero limit, got %d", len(entries))
	}
}

func TestLogger_FilterEntries(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Log entries with different levels and components
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	// Test filtering by level
	entries := logger.FilterEntries("INFO", "", time.Time{}, 0)
	if len(entries) != 1 {
		t.Errorf("Expected 1 INFO entry, got %d", len(entries))
	}

	// Test filtering by component
	entries = logger.FilterEntries("", "app", time.Time{}, 0)
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries with component 'app', got %d", len(entries))
	}

	// Test filtering by time
	now := time.Now()
	entries = logger.FilterEntries("", "", now.Add(time.Hour), 0)
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries after future time, got %d", len(entries))
	}

	// Test filtering with limit
	entries = logger.FilterEntries("", "", time.Time{}, 2)
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries with limit, got %d", len(entries))
	}
}

func TestLogger_GetLogDir(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	logDir := logger.GetLogDir()
	// Should be empty for non-file logger
	if logDir != "" {
		t.Errorf("Expected empty log dir for non-file logger, got '%s'", logDir)
	}
}

func TestLogger_GetLogFiles(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Should return error for non-file logger
	files, err := logger.GetLogFiles()
	if err == nil {
		t.Error("Expected error for non-file logger")
	}
	if files != nil {
		t.Error("Expected nil files for non-file logger")
	}
}

func TestLogger_ReadLogFile(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Should return error for non-file logger
	content, err := logger.ReadLogFile("test.log")
	if err == nil {
		t.Error("Expected error for non-file logger")
	}
	if content != nil {
		t.Error("Expected nil content for non-file logger")
	}

	// Test with file logger
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	config := &LoggerConfig{
		Level:     LogLevelInfo,
		LogToFile: true,
		Component: "test",
	}

	fileLogger, err := NewLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}
	defer fileLogger.Close()

	// Test reading invalid file name
	_, err = fileLogger.ReadLogFile("../invalid.log")
	if err == nil {
		t.Error("Expected error for invalid file name")
	}

	// Test reading non-log file
	_, err = fileLogger.ReadLogFile("not-a-log.txt")
	if err == nil {
		t.Error("Expected error for non-log file")
	}
}

func TestLogLevelFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LogLevelDebug},
		{"DEBUG", LogLevelDebug},
		{"info", LogLevelInfo},
		{"INFO", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"WARN", LogLevelWarn},
		{"warning", LogLevelWarn},
		{"error", LogLevelError},
		{"ERROR", LogLevelError},
		{"fatal", LogLevelFatal},
		{"FATAL", LogLevelFatal},
		{"unknown", LogLevelInfo}, // Default
		{"", LogLevelInfo},        // Default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := LogLevelFromString(tt.input); got != tt.expected {
				t.Errorf("LogLevelFromString(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLogger_formatTextEntry(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Component: "test",
		Message:   "test message",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Caller: "test.go:123",
		Error:  "test error",
	}

	formatted := logger.formatTextEntry(entry)

	expectedParts := []string{
		"[INFO]",
		"component=test",
		"message=\"test message\"",
		"key1=value1",
		"key2=42",
		"caller=test.go:123",
		"error=\"test error\"",
	}

	for _, part := range expectedParts {
		if !strings.Contains(formatted, part) {
			t.Errorf("Formatted message should contain '%s', got: %s", part, formatted)
		}
	}
}

func TestLogger_addEntry(t *testing.T) {
	config := &LoggerConfig{
		Level:      LogLevelInfo,
		LogToFile:  false,
		Component:  "test",
		MaxEntries: 2, // Small limit for testing
	}

	logger, err := NewLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("NewLoggerWithConfig failed: %v", err)
	}

	// Add entries beyond the limit
	entry1 := LogEntry{Level: "INFO", Message: "Entry 1"}
	entry2 := LogEntry{Level: "INFO", Message: "Entry 2"}
	entry3 := LogEntry{Level: "INFO", Message: "Entry 3"}

	logger.addEntry(entry1)
	logger.addEntry(entry2)
	logger.addEntry(entry3)

	// Should only keep the most recent entries
	entries := logger.GetRecentEntries(10)
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries after limit exceeded, got %d", len(entries))
	}

	// Should have the most recent entries
	if entries[0].Message != "Entry 2" || entries[1].Message != "Entry 3" {
		t.Error("Should keep the most recent entries")
	}
}

func TestLogger_rotateIfNeeded(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Should not error for logger without file
	err = logger.rotateIfNeeded()
	if err != nil {
		t.Errorf("rotateIfNeeded should not error for non-file logger: %v", err)
	}
}

func TestLogger_compressLogFile(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Create a temporary file for testing
	tempFile := filepath.Join(t.TempDir(), "test.log")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()

	// Should not panic
	logger.compressLogFile(tempFile)

	// Check if file was "compressed" (renamed with .gz extension)
	compressedFile := tempFile + ".gz"
	if _, err := os.Stat(compressedFile); os.IsNotExist(err) {
		t.Error("Compressed file should exist")
	}
}

func TestLogger_cleanupOldBackups(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Should not panic for logger without log directory
	logger.cleanupOldBackups()

	// Test with log directory set
	logger.logDir = t.TempDir()
	logger.rotationConfig = &LogRotationConfig{MaxBackups: 2}

	// Create some backup files
	for i := 0; i < 5; i++ {
		filename := filepath.Join(logger.logDir, "generator-test.log.backup"+string(rune('0'+i)))
		file, err := os.Create(filename)
		if err != nil {
			t.Fatalf("Failed to create backup file: %v", err)
		}
		file.Close()
	}

	// Should not panic
	logger.cleanupOldBackups()
}
