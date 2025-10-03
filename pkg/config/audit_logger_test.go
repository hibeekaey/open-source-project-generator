package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestAuditLogger(t *testing.T) (*ConfigAuditLogger, string) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test_audit.log")

	logger := &ConfigAuditLogger{
		logFile: logFile,
		logger:  nil, // Use nil logger for testing
	}

	return logger, tmpDir
}

func TestLogAction(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Test logging an action
	details := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	logger.LogAction("test_action", "test_resource", details)

	// Verify the log entry was written
	entries, err := logger.readLogEntries()
	if err != nil {
		t.Fatalf("Failed to read log entries: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Action != "test_action" {
		t.Errorf("Expected action 'test_action', got '%s'", entry.Action)
	}

	if entry.Resource != "test_resource" {
		t.Errorf("Expected resource 'test_resource', got '%s'", entry.Resource)
	}

	if !entry.Success {
		t.Error("Expected success to be true")
	}

	if entry.Details["key1"] != "value1" {
		t.Errorf("Expected details key1 to be 'value1', got '%v'", entry.Details["key1"])
	}
}

func TestLogError(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Test logging an error
	testErr := os.ErrNotExist
	details := map[string]interface{}{
		"context": "test context",
	}

	logger.LogError("test_action", "test_resource", testErr, details)

	// Verify the log entry was written
	entries, err := logger.readLogEntries()
	if err != nil {
		t.Fatalf("Failed to read log entries: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Action != "test_action" {
		t.Errorf("Expected action 'test_action', got '%s'", entry.Action)
	}

	if entry.Success {
		t.Error("Expected success to be false")
	}

	if entry.Error != testErr.Error() {
		t.Errorf("Expected error '%s', got '%s'", testErr.Error(), entry.Error)
	}
}

func TestGetAuditLogs(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Log multiple entries
	logger.LogAction("action1", "resource1", nil)
	logger.LogAction("action2", "resource2", nil)
	logger.LogError("action3", "resource3", os.ErrNotExist, nil)

	// Test getting all logs
	entries, err := logger.GetAuditLogs(nil)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 log entries, got %d", len(entries))
	}

	// Test filtering by action
	filter := &AuditLogFilter{
		Actions: []string{"action1"},
	}

	entries, err = logger.GetAuditLogs(filter)
	if err != nil {
		t.Fatalf("Failed to get filtered audit logs: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 filtered entry, got %d", len(entries))
	}

	if entries[0].Action != "action1" {
		t.Errorf("Expected action 'action1', got '%s'", entries[0].Action)
	}
}

func TestGetCurrentUser(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Save original environment
	originalUser := os.Getenv("USER")
	originalUsername := os.Getenv("USERNAME")

	// Test with USER environment variable
	if err := os.Setenv("USER", "testuser"); err != nil {
		t.Fatalf("Failed to set USER environment variable: %v", err)
	}
	if err := os.Unsetenv("USERNAME"); err != nil {
		t.Fatalf("Failed to unset USERNAME environment variable: %v", err)
	}

	user := logger.getCurrentUser()
	if user != "testuser" {
		t.Errorf("Expected 'testuser', got '%s'", user)
	}

	// Test with USERNAME environment variable
	if err := os.Unsetenv("USER"); err != nil {
		t.Fatalf("Failed to unset USER environment variable: %v", err)
	}
	if err := os.Setenv("USERNAME", "testusername"); err != nil {
		t.Fatalf("Failed to set USERNAME environment variable: %v", err)
	}

	user = logger.getCurrentUser()
	if user != "testusername" {
		t.Errorf("Expected 'testusername', got '%s'", user)
	}

	// Test with no environment variables
	if err := os.Unsetenv("USER"); err != nil {
		t.Fatalf("Failed to unset USER environment variable: %v", err)
	}
	if err := os.Unsetenv("USERNAME"); err != nil {
		t.Fatalf("Failed to unset USERNAME environment variable: %v", err)
	}

	user = logger.getCurrentUser()
	if user != "unknown" {
		t.Errorf("Expected 'unknown', got '%s'", user)
	}

	// Restore original environment
	if originalUser != "" {
		if err := os.Setenv("USER", originalUser); err != nil {
			t.Fatalf("Failed to restore USER environment variable: %v", err)
		}
	}
	if originalUsername != "" {
		if err := os.Setenv("USERNAME", originalUsername); err != nil {
			t.Fatalf("Failed to restore USERNAME environment variable: %v", err)
		}
	}
}

func TestValidateLogEntry(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Test valid entry
	entry := &AuditLogEntry{
		Timestamp: time.Now(),
		Action:    "test_action",
		Resource:  "test_resource",
	}

	err := logger.validateLogEntry(entry)
	if err != nil {
		t.Errorf("Expected valid entry to pass validation, got error: %v", err)
	}

	// Test entry with missing timestamp
	entry = &AuditLogEntry{
		Action:   "test_action",
		Resource: "test_resource",
	}

	err = logger.validateLogEntry(entry)
	if err == nil {
		t.Error("Expected error for missing timestamp")
	}

	// Test entry with missing action
	entry = &AuditLogEntry{
		Timestamp: time.Now(),
		Resource:  "test_resource",
	}

	err = logger.validateLogEntry(entry)
	if err == nil {
		t.Error("Expected error for missing action")
	}

	// Test entry with missing resource
	entry = &AuditLogEntry{
		Timestamp: time.Now(),
		Action:    "test_action",
	}

	err = logger.validateLogEntry(entry)
	if err == nil {
		t.Error("Expected error for missing resource")
	}
}

func TestApplyFilters(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Create test entries
	now := time.Now()
	entries := []*AuditLogEntry{
		{
			Timestamp: now.Add(-2 * time.Hour),
			Action:    "action1",
			Resource:  "resource1",
			User:      "user1",
			Success:   true,
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Action:    "action2",
			Resource:  "resource2",
			User:      "user2",
			Success:   false,
		},
		{
			Timestamp: now,
			Action:    "action1",
			Resource:  "resource1",
			User:      "user1",
			Success:   true,
		},
	}

	// Test time range filter
	startTime := now.Add(-90 * time.Minute)
	filter := &AuditLogFilter{
		StartTime: &startTime,
	}

	filtered := logger.applyFilters(entries, filter)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries after time filter, got %d", len(filtered))
	}

	// Test action filter
	filter = &AuditLogFilter{
		Actions: []string{"action1"},
	}

	filtered = logger.applyFilters(entries, filter)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries after action filter, got %d", len(filtered))
	}

	// Test success filter
	successFilter := true
	filter = &AuditLogFilter{
		Success: &successFilter,
	}

	filtered = logger.applyFilters(entries, filter)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 entries after success filter, got %d", len(filtered))
	}
}

func TestApplyPagination(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Create test entries
	entries := make([]*AuditLogEntry, 10)
	for i := 0; i < 10; i++ {
		entries[i] = &AuditLogEntry{
			Action:   "action",
			Resource: "resource",
		}
	}

	// Test pagination
	paginated := logger.applyPagination(entries, 2, 3)
	if len(paginated) != 3 {
		t.Errorf("Expected 3 entries after pagination, got %d", len(paginated))
	}

	// Test pagination beyond bounds
	paginated = logger.applyPagination(entries, 15, 5)
	if len(paginated) != 0 {
		t.Errorf("Expected 0 entries for offset beyond bounds, got %d", len(paginated))
	}

	// Test pagination with no limit
	paginated = logger.applyPagination(entries, 5, 0)
	if len(paginated) != 5 {
		t.Errorf("Expected 5 entries with no limit, got %d", len(paginated))
	}
}

func TestGetTopErrors(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	errorCounts := map[string]int{
		"error1": 5,
		"error2": 3,
		"error3": 8,
		"error4": 1,
		"error5": 2,
	}

	topErrors := logger.getTopErrors(errorCounts, 3)

	if len(topErrors) != 3 {
		t.Errorf("Expected 3 top errors, got %d", len(topErrors))
	}

	// Should be sorted by count (descending)
	expectedOrder := []string{"error3", "error1", "error2"}
	for i, expected := range expectedOrder {
		if topErrors[i] != expected {
			t.Errorf("Expected error %d to be '%s', got '%s'", i, expected, topErrors[i])
		}
	}
}

func TestLogFileOperations(t *testing.T) {
	logger, _ := createTestAuditLogger(t)

	// Test that log file is created when logging
	logger.LogAction("test", "test", nil)

	if _, err := os.Stat(logger.logFile); os.IsNotExist(err) {
		t.Error("Log file should be created after logging")
	}

	// Test file size calculation
	size := logger.getLogFileSize()
	if size <= 0 {
		t.Error("Log file size should be greater than 0")
	}
}
