// Package test_helpers provides comprehensive test utilities for the entire application.
//
// This package consolidates all test helper functions and utilities from various packages
// into a single, comprehensive test system that can be used across the application.
package test_helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestHelper provides comprehensive test utilities
type TestHelper struct {
	tempFiles []string
	tempDirs  []string
	cleanup   []func() error
	config    *TestConfig
}

// TestConfig holds configuration for test utilities
type TestConfig struct {
	EnableDB       bool          `json:"enable_db"`
	EnableLogging  bool          `json:"enable_logging"`
	CleanupTimeout time.Duration `json:"cleanup_timeout"`
	MaxTempFiles   int           `json:"max_temp_files"`
	MaxTempDirs    int           `json:"max_temp_dirs"`
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		EnableDB:       true,
		EnableLogging:  false,
		CleanupTimeout: 30 * time.Second,
		MaxTempFiles:   100,
		MaxTempDirs:    50,
	}
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T, config *TestConfig) *TestHelper {
	if config == nil {
		config = DefaultTestConfig()
	}

	helper := &TestHelper{
		tempFiles: make([]string, 0),
		tempDirs:  make([]string, 0),
		cleanup:   make([]func() error, 0),
		config:    config,
	}

	// Set up cleanup
	t.Cleanup(func() {
		helper.Cleanup(t)
	})

	return helper
}

// CreateTempDir creates a temporary directory for testing
func (h *TestHelper) CreateTempDir(t *testing.T, prefix string) string {
	if len(h.tempDirs) >= h.config.MaxTempDirs {
		t.Fatalf("Maximum number of temp directories (%d) exceeded", h.config.MaxTempDirs)
	}

	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	h.tempDirs = append(h.tempDirs, dir)
	return dir
}

// CreateTempFile creates a temporary file for testing
func (h *TestHelper) CreateTempFile(t *testing.T, prefix, content string) string {
	if len(h.tempFiles) >= h.config.MaxTempFiles {
		t.Fatalf("Maximum number of temp files (%d) exceeded", h.config.MaxTempFiles)
	}

	file, err := os.CreateTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if content != "" {
		if _, err := file.WriteString(content); err != nil {
			_ = file.Close()
			_ = os.Remove(file.Name())
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}

	_ = file.Close()
	h.tempFiles = append(h.tempFiles, file.Name())
	return file.Name()
}

// SetupTestDB sets up a test database (simplified version without GORM)
func (h *TestHelper) SetupTestDB(t *testing.T) interface{} {
	if !h.config.EnableDB {
		t.Skip("Database testing is disabled")
	}

	// Return a mock database interface
	// In a real implementation, this would set up an actual database
	return &MockDatabase{}
}

// MockDatabase is a mock database for testing
type MockDatabase struct {
	Data map[string]interface{}
}

// CleanupTestDB cleans up test database
func (h *TestHelper) CleanupTestDB() error {
	// Mock cleanup - in real implementation would clean actual database
	return nil
}

// SeedTestData seeds the database with test data
func (h *TestHelper) SeedTestData(t *testing.T) {
	// Mock seed data - in real implementation would seed actual database
	t.Log("Mock: Seeding test data")
}

// CreateTestProjectConfig creates a test project configuration
func (h *TestHelper) CreateTestProjectConfig() *TestProjectConfig {
	return &TestProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "A test project",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/test-org/test-project",
		License:      "MIT",
		OutputPath:   h.CreateTempDir(nil, "test-output"),
		GeneratedAt:  time.Now(),
		Features:     []string{"backend", "frontend"},
	}
}

// TestProjectConfig represents a test project configuration
type TestProjectConfig struct {
	Name         string    `json:"name"`
	Organization string    `json:"organization"`
	Description  string    `json:"description"`
	Author       string    `json:"author"`
	Email        string    `json:"email"`
	Repository   string    `json:"repository"`
	License      string    `json:"license"`
	OutputPath   string    `json:"output_path"`
	GeneratedAt  time.Time `json:"generated_at"`
	Features     []string  `json:"features"`
}

// CreateTestConfigFile creates a test configuration file
func (h *TestHelper) CreateTestConfigFile(t *testing.T, config *TestProjectConfig) string {
	if config == nil {
		config = h.CreateTestProjectConfig()
	}

	// Create YAML content
	yamlContent := fmt.Sprintf(`
name: %s
organization: %s
description: %s
author: %s
email: %s
repository: %s
license: %s
output_path: %s
`, config.Name, config.Organization, config.Description, config.Author,
		config.Email, config.Repository, config.License, config.OutputPath)

	return h.CreateTempFile(t, "test-config-*.yaml", yamlContent)
}

// RegisterCleanup registers a cleanup function
func (h *TestHelper) RegisterCleanup(cleanup func() error) {
	h.cleanup = append(h.cleanup, cleanup)
}

// Cleanup performs all cleanup operations
func (h *TestHelper) Cleanup(t *testing.T) {
	// Execute custom cleanup functions
	for _, cleanup := range h.cleanup {
		if err := cleanup(); err != nil {
			t.Logf("Cleanup function failed: %v", err)
		}
	}

	// Clean up database
	if err := h.CleanupTestDB(); err != nil {
		t.Logf("Database cleanup failed: %v", err)
	}

	// Clean up temporary files
	for _, file := range h.tempFiles {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove temp file %s: %v", file, err)
		}
	}

	// Clean up temporary directories
	for _, dir := range h.tempDirs {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("Failed to remove temp dir %s: %v", dir, err)
		}
	}

	// Force garbage collection
	runtime.GC()
	runtime.GC() // Double GC to ensure cleanup
}

// AssertFileExists checks if a file exists
func (h *TestHelper) AssertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to exist, but it doesn't", path)
	}
}

// AssertFileNotExists checks if a file doesn't exist
func (h *TestHelper) AssertFileNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Expected file %s to not exist, but it does", path)
	}
}

// AssertDirExists checks if a directory exists
func (h *TestHelper) AssertDirExists(t *testing.T, path string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Fatalf("Expected directory %s to exist, but it doesn't", path)
	}
	if !info.IsDir() {
		t.Fatalf("Expected %s to be a directory, but it's not", path)
	}
}

// AssertFileContent checks if a file contains expected content
func (h *TestHelper) AssertFileContent(t *testing.T, path, expectedContent string) {
	// Validate and sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
		t.Fatalf("Invalid path: %s", path)
	}
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	if !strings.Contains(string(content), expectedContent) {
		t.Fatalf("Expected file %s to contain '%s', but it doesn't", path, expectedContent)
	}
}

// AssertFileNotContent checks if a file doesn't contain content
func (h *TestHelper) AssertFileNotContent(t *testing.T, path, unexpectedContent string) {
	// Validate and sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
		t.Fatalf("Invalid path: %s", path)
	}
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	if strings.Contains(string(content), unexpectedContent) {
		t.Fatalf("Expected file %s to not contain '%s', but it does", path, unexpectedContent)
	}
}

// CreateTestTree creates a directory tree for testing
func (h *TestHelper) CreateTestTree(t *testing.T, baseDir string, structure map[string]interface{}) {
	for name, content := range structure {
		path := filepath.Join(baseDir, name)

		switch v := content.(type) {
		case string:
			// It's a file
			if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
				t.Fatalf("Failed to create directory for %s: %v", path, err)
			}
			if err := os.WriteFile(path, []byte(v), 0600); err != nil {
				t.Fatalf("Failed to create file %s: %v", path, err)
			}
		case map[string]interface{}:
			// It's a directory
			if err := os.MkdirAll(path, 0750); err != nil {
				t.Fatalf("Failed to create directory %s: %v", path, err)
			}
			h.CreateTestTree(t, path, v)
		}
	}
}

// AssertTestTree checks if a directory tree matches expected structure
func (h *TestHelper) AssertTestTree(t *testing.T, baseDir string, expectedStructure map[string]interface{}) {
	for name, expectedContent := range expectedStructure {
		path := filepath.Join(baseDir, name)

		switch v := expectedContent.(type) {
		case string:
			// It should be a file
			h.AssertFileExists(t, path)
			h.AssertFileContent(t, path, v)
		case map[string]interface{}:
			// It should be a directory
			h.AssertDirExists(t, path)
			h.AssertTestTree(t, path, v)
		}
	}
}

// WaitForCondition waits for a condition to be true
func (h *TestHelper) WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Condition not met within timeout: %s", message)
}

// CaptureOutput captures stdout/stderr for testing
func (h *TestHelper) CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	// This is a simplified version - in a real implementation,
	// you'd use os.Pipe() to capture actual stdout/stderr
	// For now, we'll just run the function
	fn()
	return "", ""
}

// MockLogger creates a mock logger for testing
func (h *TestHelper) MockLogger(t *testing.T) *MockLogger {
	return &MockLogger{
		Logs: make([]LogEntry, 0),
	}
}

// MockLogger is a mock logger for testing
type MockLogger struct {
	Logs []LogEntry
}

// LogEntry represents a log entry
type LogEntry struct {
	Level   string
	Message string
	Fields  map[string]interface{}
}

// Debug logs a debug message
func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Logs = append(m.Logs, LogEntry{
		Level:   "debug",
		Message: fmt.Sprintf(msg, args...),
	})
}

// Info logs an info message
func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Logs = append(m.Logs, LogEntry{
		Level:   "info",
		Message: fmt.Sprintf(msg, args...),
	})
}

// Warn logs a warning message
func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Logs = append(m.Logs, LogEntry{
		Level:   "warn",
		Message: fmt.Sprintf(msg, args...),
	})
}

// Error logs an error message
func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Logs = append(m.Logs, LogEntry{
		Level:   "error",
		Message: fmt.Sprintf(msg, args...),
	})
}

// Fatal logs a fatal message
func (m *MockLogger) Fatal(msg string, args ...interface{}) {
	m.Logs = append(m.Logs, LogEntry{
		Level:   "fatal",
		Message: fmt.Sprintf(msg, args...),
	})
}

// GetLogs returns all logged messages
func (m *MockLogger) GetLogs() []LogEntry {
	return m.Logs
}

// ClearLogs clears all logged messages
func (m *MockLogger) ClearLogs() {
	m.Logs = make([]LogEntry, 0)
}

// HasLogLevel checks if any logs exist for a specific level
func (m *MockLogger) HasLogLevel(level string) bool {
	for _, log := range m.Logs {
		if log.Level == level {
			return true
		}
	}
	return false
}

// GetLogsByLevel returns logs for a specific level
func (m *MockLogger) GetLogsByLevel(level string) []LogEntry {
	var logs []LogEntry
	for _, log := range m.Logs {
		if log.Level == level {
			logs = append(logs, log)
		}
	}
	return logs
}
