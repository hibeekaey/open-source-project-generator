package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/internal/testutils"
)

// setupAuditTest creates a temporary directory and audit trail for testing
func setupAuditTest(t *testing.T) (string, *AuditTrail, func()) {
	suite := testutils.NewTestSuite()
	tempDir, err := suite.File.CreateTempDir("test_audit")
	suite.Assertions.AssertNoError(t, err, "Failed to create temp directory")

	logFile := filepath.Join(tempDir, "audit.log")
	audit := NewAuditTrail(logFile)

	cleanup := func() {
		suite.File.CleanupTestFiles(tempDir)
	}

	return tempDir, audit, cleanup
}

func TestAuditTrail_LogVersionUpdate(t *testing.T) {
	tempDir, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Test successful version update
	err := audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
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
	logFile := filepath.Join(tempDir, "audit.log")
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
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	versionChanges := map[string]string{
		"react":  "18.0.0 -> 19.0.0",
		"nextjs": "14.0.0 -> 15.0.0",
	}

	// Test template update
	err := audit.LogTemplateUpdate("templates/frontend/nextjs-app", versionChanges, true, nil)
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
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Test security scan
	err := audit.LogSecurityScan(50, 3, true, nil)
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

func TestAuditTrailGetAuditHistory_WithTimeFilter(t *testing.T) {
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

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

func TestAuditTrailGetAuditHistory_WithEventTypeFilter(t *testing.T) {
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Log different types of events
	err := audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
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
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Log various events
	err := audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
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
	tempDir, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Create a log file with some content
	err := audit.LogVersionUpdate("react", "18.0.0", "19.0.0", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	// Test rotation with a very small size limit (should trigger rotation)
	err = audit.RotateLog(1) // 1 byte limit
	if err != nil {
		t.Fatalf("Failed to rotate log: %v", err)
	}

	// Check that original log file no longer exists
	logFile := filepath.Join(tempDir, "audit.log")
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

func TestAuditTrail_SecureIDGeneration(t *testing.T) {
	_, audit, cleanup := setupAuditTest(t)
	defer cleanup()

	// Generate multiple audit events and collect their IDs
	var eventIDs []string
	for i := 0; i < 100; i++ {
		err := audit.LogVersionUpdate("test-package", "1.0.0", "1.0.1", true, nil)
		if err != nil {
			t.Fatalf("Failed to log version update: %v", err)
		}
	}

	// Retrieve all events and collect IDs
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 100 {
		t.Fatalf("Expected 100 events, got %d", len(events))
	}

	for _, event := range events {
		eventIDs = append(eventIDs, event.ID)
	}

	// Test 1: All IDs should be unique (collision resistance)
	idSet := make(map[string]bool)
	for _, id := range eventIDs {
		if idSet[id] {
			t.Errorf("Duplicate ID found: %s", id)
		}
		idSet[id] = true
	}

	// Test 2: All IDs should have the "audit_" prefix
	for _, id := range eventIDs {
		if !strings.HasPrefix(id, "audit_") {
			t.Errorf("ID does not have 'audit_' prefix: %s", id)
		}
	}

	// Test 3: IDs should not be predictable (not timestamp-based)
	// Check that IDs don't follow a sequential pattern
	for i := 1; i < len(eventIDs); i++ {
		// Extract the suffix after "audit_"
		suffix1 := strings.TrimPrefix(eventIDs[i-1], "audit_")
		suffix2 := strings.TrimPrefix(eventIDs[i], "audit_")

		// If these were timestamp-based, they would be sequential numbers
		// With secure random generation, they should be completely different
		if suffix1 == suffix2 {
			t.Errorf("Sequential IDs are identical, indicating predictable generation")
		}

		// Check that suffixes are not simple increments (would indicate timestamp-based)
		if len(suffix1) > 0 && len(suffix2) > 0 {
			// For hex strings, check if they're not sequential
			if isSequentialHex(suffix1, suffix2) {
				t.Errorf("IDs appear to be sequential, indicating predictable generation: %s -> %s", suffix1, suffix2)
			}
		}
	}

	// Test 4: IDs should have sufficient entropy (length check)
	for _, id := range eventIDs {
		suffix := strings.TrimPrefix(id, "audit_")
		if len(suffix) < 16 { // Should have at least 16 characters of randomness
			t.Errorf("ID suffix too short, may lack sufficient entropy: %s", suffix)
		}
	}
}

func TestAuditTrail_SecureIDGenerationWithCustomRandom(t *testing.T) {
	tempDir, _, cleanup := setupAuditTest(t)
	defer cleanup()

	logFile := filepath.Join(tempDir, "audit.log")

	// Create audit trail with custom secure random generator
	customRandom := &MockSecureRandom{}
	audit := NewAuditTrailWithSecureRandom(logFile, customRandom)

	// Log an event
	err := audit.LogVersionUpdate("test-package", "1.0.0", "1.0.1", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	// Verify that custom random generator was used
	if !customRandom.GenerateSecureIDCalled {
		t.Errorf("Custom secure random generator was not called")
	}

	// Retrieve event and check ID
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	expectedID := "audit_mock_secure_id"
	if events[0].ID != expectedID {
		t.Errorf("Expected ID '%s', got '%s'", expectedID, events[0].ID)
	}
}

func TestAuditTrail_SecureIDGenerationFallback(t *testing.T) {
	tempDir, _, cleanup := setupAuditTest(t)
	defer cleanup()

	logFile := filepath.Join(tempDir, "audit.log")

	// Create audit trail with failing secure random generator
	failingRandom := &FailingSecureRandom{}
	audit := NewAuditTrailWithSecureRandom(logFile, failingRandom)

	// Log an event (should use fallback ID generation)
	err := audit.LogVersionUpdate("test-package", "1.0.0", "1.0.1", true, nil)
	if err != nil {
		t.Fatalf("Failed to log version update: %v", err)
	}

	// Retrieve event and check that fallback ID was used
	events, err := audit.GetAuditHistory(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	// Fallback ID should have "audit_fallback_" prefix
	if !strings.HasPrefix(events[0].ID, "audit_fallback_") {
		t.Errorf("Expected fallback ID with 'audit_fallback_' prefix, got '%s'", events[0].ID)
	}
}

// Helper function to check if two hex strings are sequential
func isSequentialHex(hex1, hex2 string) bool {
	// This is a simple check - in practice, timestamp-based IDs would show
	// very small differences between consecutive calls
	if len(hex1) != len(hex2) {
		return false
	}

	// Convert to integers and check if difference is small (indicating timestamp-based)
	// This is a heuristic - real random hex strings should have large differences
	var diff int
	for i := 0; i < len(hex1) && i < len(hex2); i++ {
		if hex1[i] != hex2[i] {
			diff++
		}
	}

	// If only a few characters differ, it might be sequential
	// Random strings should differ in many positions
	return diff < 3
}

// Mock implementations for testing

type MockSecureRandom struct {
	GenerateSecureIDCalled bool
}

func (m *MockSecureRandom) GenerateRandomSuffix(length int) (string, error) {
	return "mock_suffix", nil
}

func (m *MockSecureRandom) GenerateSecureID(prefix string) (string, error) {
	m.GenerateSecureIDCalled = true
	return prefix + "_mock_secure_id", nil
}

func (m *MockSecureRandom) GenerateBytes(length int) ([]byte, error) {
	return make([]byte, length), nil
}

func (m *MockSecureRandom) GenerateHexString(length int) (string, error) {
	return strings.Repeat("a", length), nil
}

func (m *MockSecureRandom) GenerateBase64String(length int) (string, error) {
	return strings.Repeat("A", length), nil
}

func (m *MockSecureRandom) GenerateAlphanumeric(length int) (string, error) {
	return strings.Repeat("A", length), nil
}

type FailingSecureRandom struct{}

func (f *FailingSecureRandom) GenerateRandomSuffix(length int) (string, error) {
	return "", fmt.Errorf("random generation failed")
}

func (f *FailingSecureRandom) GenerateSecureID(prefix string) (string, error) {
	return "", fmt.Errorf("secure ID generation failed")
}

func (f *FailingSecureRandom) GenerateBytes(length int) ([]byte, error) {
	return nil, fmt.Errorf("byte generation failed")
}

func (f *FailingSecureRandom) GenerateHexString(length int) (string, error) {
	return "", fmt.Errorf("hex generation failed")
}

func (f *FailingSecureRandom) GenerateBase64String(length int) (string, error) {
	return "", fmt.Errorf("base64 generation failed")
}

func (f *FailingSecureRandom) GenerateAlphanumeric(length int) (string, error) {
	return "", fmt.Errorf("alphanumeric generation failed")
}
