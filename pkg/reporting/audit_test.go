package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAuditTrail_LogVersionUpdate(t *testing.T) {
	// Create temporary file for audit log
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Test successful version update
	err = audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	// Test failed version update
	testErr := fmt.Errorf("update failed")
	err = audit.LogVersionUpdate("nextjs", "14.0.0", "15.0.0", false, testErr)
	if err != nil {
		t.Fatalf("Failed to log failed version update: %v", err)
	}

	// Verify log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Audit log file was not created")
	}

	// Read and verify log entries
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 audit events, got %d", len(events))
	}

	// Verify first event (successful)
	successEvent := events[0]
	if successEvent.EventType != "version_update" {
		t.Errorf("Expected event type 'version_update', got '%s'", successEvent.EventType)
	}
	if successEvent.Resource != "react" {
		t.Errorf("Expected resource 'react', got '%s'", successEvent.Resource)
	}
	if !successEvent.Success {
		t.Errorf("Expected successful event")
	}
	if successEvent.OldValue != "18.0.0" {
		t.Errorf("Expected old value '18.0.0', got '%s'", successEvent.OldValue)
	}
	if successEvent.NewValue != "19.0.0" {
		t.Errorf("Expected new value '19.0.0', got '%s'", successEvent.NewValue)
	}

	// Verify second event (failed)
	failEvent := events[1]
	if failEvent.Success {
		t.Errorf("Expected failed event")
	}
	if failEvent.Error != "update failed" {
		t.Errorf("Expected error 'update failed', got '%s'", failEvent.Error)
	}
}

func TestAuditTrail_LogTemplateUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	versionChanges := map[string]string{
		"react":  "18.0.0 -> 19.0.0",
		"nextjs": "14.0.0 -> 15.0.0",
	}

	// Test template update
	err = audit.LogTemplateUpdate("templates/frontend/nextjs-app", versionChanges, true, nil)
	if err != nil {
		t.Fatalf("Failed to log template update: %v", err)
	}

	// Verify log entry
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "template_update")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 audit event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "template_update" {
		t.Errorf("Expected event type 'template_update', got '%s'", event.EventType)
	}
	if event.Resource != "templates/frontend/nextjs-app" {
		t.Errorf("Expected resource 'templates/frontend/nextjs-app', got '%s'", event.Resource)
	}

	// Verify metadata
	if changesCount, ok := event.Metadata["changes_count"].(float64); !ok || int(changesCount) != 2 {
		t.Errorf("Expected changes_count metadata to be 2, got %v", event.Metadata["changes_count"])
	}
}

func TestAuditTrail_LogSecurityScan(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Test security scan
	err = audit.LogSecurityScan(50, 3, true, nil)
	if err != nil {
		t.Fatalf("Failed to log security scan: %v", err)
	}

	// Verify log entry
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "security_scan")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 audit event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "security_scan" {
		t.Errorf("Expected event type 'security_scan', got '%s'", event.EventType)
	}

	// Verify metadata
	if packagesScanned, ok := event.Metadata["packages_scanned"].(float64); !ok || int(packagesScanned) != 50 {
		t.Errorf("Expected packages_scanned to be 50, got %v", event.Metadata["packages_scanned"])
	}
	if issuesFound, ok := event.Metadata["issues_found"].(float64); !ok || int(issuesFound) != 3 {
		t.Errorf("Expected issues_found to be 3, got %v", event.Metadata["issues_found"])
	}
}

func TestAuditTrail_GetAuditHistory_TimeFilter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Log events at different times
	now := time.Now()

	// Log an old event (should be filtered out)
	oldEvent := AuditEvent{
		ID:        "old_event",
		Timestamp: now.Add(-2 * time.Hour),
		EventType: "version_update",
		Success:   true,
	}
	audit.writeEvent(oldEvent)

	// Log a recent event (should be included)
	recentEvent := AuditEvent{
		ID:        "recent_event",
		Timestamp: now.Add(-30 * time.Minute),
		EventType: "version_update",
		Success:   true,
	}
	audit.writeEvent(recentEvent)

	// Get events from last hour
	events, err := audit.GetAuditHistory(now.Add(-1*time.Hour), now.Add(1*time.Hour), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event within time range, got %d", len(events))
	}

	if events[0].ID != "recent_event" {
		t.Errorf("Expected recent_event, got %s", events[0].ID)
	}
}

func TestAuditTrail_GetAuditHistory_EventTypeFilter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Log different types of events
	err = audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	err = audit.LogSecurityScan(10, 0, true, nil)
	if err != nil {
		t.Fatalf("Failed to log security scan: %v", err)
	}

	// Get only version_update events
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "version_update")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 version_update event, got %d", len(events))
	}

	if events[0].EventType != "version_update" {
		t.Errorf("Expected version_update event, got %s", events[0].EventType)
	}
}

func TestAuditTrail_GetAuditSummary(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Log various events
	err = audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	err = audit.LogVersionUpdate("nextjs", "14.0.0", "15.0.0", false, fmt.Errorf("failed"))
	if err != nil {
		t.Fatalf("Failed to log failed version update: %v", err)
	}

	err = audit.LogSecurityScan(10, 2, true, nil)
	if err != nil {
		t.Fatalf("Failed to log security scan: %v", err)
	}

	// Get summary
	summary, err := audit.GetAuditSummary(time.Now().Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("Failed to get audit summary: %v", err)
	}

	if summary.TotalEvents != 3 {
		t.Errorf("Expected 3 total events, got %d", summary.TotalEvents)
	}

	if summary.EventTypes["version_update"] != 2 {
		t.Errorf("Expected 2 version_update events, got %d", summary.EventTypes["version_update"])
	}

	if summary.EventTypes["security_scan"] != 1 {
		t.Errorf("Expected 1 security_scan event, got %d", summary.EventTypes["security_scan"])
	}

	// Success rate should be 66.67% (2 out of 3 successful)
	expectedSuccessRate := 66.66666666666667
	if summary.SuccessRate < expectedSuccessRate-0.1 || summary.SuccessRate > expectedSuccessRate+0.1 {
		t.Errorf("Expected success rate around %.2f%%, got %.2f%%", expectedSuccessRate, summary.SuccessRate)
	}
}

func TestAuditTrail_RotateLog(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_audit")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	// Create a log file with some content
	err = audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	// Test rotation with a very small size limit (should trigger rotation)
	err = audit.RotateLog(1) // 1 byte limit
	if err != nil {
		t.Fatalf("Failed to rotate log: %v", err)
	}

	// Check that original log file no longer exists
	if _, err := os.Stat(logFile); !os.IsNotExist(err) {
		t.Errorf("Original log file should not exist after rotation")
	}

	// Check that rotated file exists
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	rotatedFileFound := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "audit.log.") {
			rotatedFileFound = true
			break
		}
	}

	if !rotatedFileFound {
		t.Errorf("Rotated log file not found")
	}
}

func TestAuditTrail_NonExistentLogFile(t *testing.T) {
	audit := NewAuditTrail("/non/existent/path/audit.log")

	// Should handle non-existent log file gracefully
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now(), "")
	if err != nil {
		t.Fatalf("Failed to handle non-existent log file: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events from non-existent log file, got %d", len(events))
	}
}
