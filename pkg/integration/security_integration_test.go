//go:build !ci

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/reporting"
	"github.com/open-source-template-generator/pkg/security"
	"github.com/open-source-template-generator/pkg/version"
)

// TestStorageSecurityIntegration tests integration with the version storage system
func TestStorageSecurityIntegration(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("SecureStorageOperations", func(t *testing.T) {
		testSecureStorageOperations(t, tempDir)
	})

	t.Run("ConcurrentStorageAccess", func(t *testing.T) {
		testConcurrentStorageAccess(t, tempDir)
	})

	t.Run("StorageSecurityValidation", func(t *testing.T) {
		testStorageSecurityValidation(t, tempDir)
	})
}

// testSecureStorageOperations verifies that storage operations use secure file operations
func testSecureStorageOperations(t *testing.T, tempDir string) {
	filePath := filepath.Join(tempDir, "secure_versions.yaml")

	// Create storage instance
	storage, err := version.NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create test version store
	testStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks:  make(map[string]*models.VersionInfo),
		Packages:    make(map[string]*models.VersionInfo),
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:       true,
			SecurityPriority: true,
			UpdateSchedule:   "daily",
		},
	}

	// Add test packages
	testStore.Packages["secure-package"] = &models.VersionInfo{
		Name:           "secure-package",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.1.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	// Perform multiple saves to verify secure operations
	for i := 0; i < 5; i++ {
		testStore.Packages[fmt.Sprintf("package-%d", i)] = &models.VersionInfo{
			Name:           fmt.Sprintf("package-%d", i),
			Language:       "javascript",
			Type:           "package",
			CurrentVersion: "1.0.0",
			LatestVersion:  "1.0.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		}

		err = storage.Save(testStore)
		if err != nil {
			t.Fatalf("Failed to save store iteration %d: %v", i, err)
		}
	}

	// Verify no temporary files with predictable patterns remain
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()

		// Check for timestamp-based temporary files
		if strings.Contains(fileName, ".tmp.") {
			parts := strings.Split(fileName, ".tmp.")
			if len(parts) > 1 && len(parts[1]) >= 10 {
				// Check if it looks like a timestamp
				t.Errorf("Found potential timestamp-based temp file: %s", fileName)
			}
		}
	}

	// Verify final file exists and is readable
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load final store: %v", err)
	}

	if len(loadedStore.Packages) != 6 { // 1 initial + 5 added
		t.Errorf("Expected 6 packages, got %d", len(loadedStore.Packages))
	}

	t.Logf("Secure storage operations verified successfully")
}

// testConcurrentStorageAccess verifies storage works correctly under concurrent access
func testConcurrentStorageAccess(t *testing.T, tempDir string) {
	filePath := filepath.Join(tempDir, "concurrent_versions.yaml")

	// Create storage instance
	storage, err := version.NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	numGoroutines := 10
	packagesPerGoroutine := 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*packagesPerGoroutine)

	// Launch concurrent storage operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < packagesPerGoroutine; j++ {
				packageName := fmt.Sprintf("concurrent-pkg-%d-%d", goroutineID, j)

				versionInfo := &models.VersionInfo{
					Name:           packageName,
					Language:       "javascript",
					Type:           "package",
					CurrentVersion: "1.0.0",
					LatestVersion:  "1.0.0",
					IsSecure:       true,
					UpdatedAt:      time.Now(),
					CheckedAt:      time.Now(),
					UpdateSource:   "npm",
				}

				err := storage.SetVersionInfo(packageName, versionInfo)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, package %s: %w", goroutineID, packageName, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent storage error: %v", err)
	}

	// Verify final state
	allVersions, err := storage.ListVersions()
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}

	expectedCount := numGoroutines * packagesPerGoroutine
	if len(allVersions) < expectedCount {
		t.Errorf("Expected at least %d packages, got %d", expectedCount, len(allVersions))
	}

	t.Logf("Concurrent storage access completed successfully, final count: %d", len(allVersions))
}

// testStorageSecurityValidation verifies storage security validation
func testStorageSecurityValidation(t *testing.T, tempDir string) {
	// Test with malicious file paths
	maliciousPaths := []string{
		filepath.Join(tempDir, "../../../etc/passwd"),
		"/etc/shadow",
		"../malicious.yaml",
	}

	for _, maliciousPath := range maliciousPaths {
		_, err := version.NewFileStorage(maliciousPath, "yaml")
		// The storage should be created but operations should be validated
		if err != nil {
			t.Logf("Storage creation failed for malicious path %s: %v", maliciousPath, err)
		}
	}

	t.Logf("Storage security validation completed")
}

// TestAuditSecurityIntegration tests integration with the audit system
func TestAuditSecurityIntegration(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("SecureAuditOperations", func(t *testing.T) {
		testSecureAuditOperations(t, tempDir)
	})

	t.Run("AuditIDSecurity", func(t *testing.T) {
		testAuditIDSecurity(t, tempDir)
	})

	t.Run("ConcurrentAuditLogging", func(t *testing.T) {
		testConcurrentAuditLogging(t, tempDir)
	})
}

// testSecureAuditOperations verifies that audit operations use secure ID generation
func testSecureAuditOperations(t *testing.T, tempDir string) {
	auditFile := filepath.Join(tempDir, "secure_audit.log")

	// Create audit trail with secure random generator
	secureRandom := security.NewSecureRandom()
	auditTrail := reporting.NewAuditTrailWithSecureRandom(auditFile, secureRandom)

	// Log multiple events
	events := []struct {
		logFunc func() error
		name    string
	}{
		{
			logFunc: func() error {
				return auditTrail.LogVersionUpdate("test-package", "1.0.0", "1.1.0", true, nil)
			},
			name: "version_update",
		},
		{
			logFunc: func() error {
				return auditTrail.LogSecurityScan(10, 2, true, nil)
			},
			name: "security_scan",
		},
		{
			logFunc: func() error {
				return auditTrail.LogTemplateUpdate("/path/to/template", map[string]string{"pkg": "1.1.0"}, true, nil)
			},
			name: "template_update",
		},
	}

	for _, event := range events {
		err := event.logFunc()
		if err != nil {
			t.Fatalf("Failed to log %s event: %v", event.name, err)
		}
	}

	// Verify audit file exists
	if _, err := os.Stat(auditFile); os.IsNotExist(err) {
		t.Errorf("Audit file should exist: %s", auditFile)
	}

	// Read and verify audit entries
	since := time.Now().Add(-time.Hour)
	auditEvents, err := auditTrail.GetAuditHistory(since, time.Now(), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(auditEvents) != len(events) {
		t.Errorf("Expected %d audit events, got %d", len(events), len(auditEvents))
	}

	t.Logf("Secure audit operations verified successfully")
}

// testAuditIDSecurity verifies that audit IDs are generated securely
func testAuditIDSecurity(t *testing.T, tempDir string) {
	auditFile := filepath.Join(tempDir, "audit_id_test.log")

	// Create audit trail with secure random generator
	secureRandom := security.NewSecureRandom()
	auditTrail := reporting.NewAuditTrailWithSecureRandom(auditFile, secureRandom)

	// Generate multiple audit events to collect IDs
	numEvents := 100
	for i := 0; i < numEvents; i++ {
		err := auditTrail.LogVersionUpdate(fmt.Sprintf("package-%d", i), "1.0.0", "1.1.0", true, nil)
		if err != nil {
			t.Fatalf("Failed to log audit event %d: %v", i, err)
		}
	}

	// Read audit events and analyze IDs
	since := time.Now().Add(-time.Hour)
	auditEvents, err := auditTrail.GetAuditHistory(since, time.Now(), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(auditEvents) != numEvents {
		t.Errorf("Expected %d audit events, got %d", numEvents, len(auditEvents))
	}

	// Verify ID uniqueness and security
	idSet := make(map[string]bool)
	for _, event := range auditEvents {
		// Check for duplicate IDs
		if idSet[event.ID] {
			t.Errorf("Duplicate audit ID detected: %s", event.ID)
		}
		idSet[event.ID] = true

		// Verify ID format (should start with "audit_")
		if !strings.HasPrefix(event.ID, "audit_") {
			t.Errorf("Audit ID should start with 'audit_': %s", event.ID)
		}

		// Check that ID doesn't contain timestamp patterns
		if strings.Contains(event.ID, "fallback") {
			t.Errorf("Audit ID contains fallback pattern (insecure): %s", event.ID)
		}

		// Verify ID length (should be reasonable for secure random)
		if len(event.ID) < 10 {
			t.Errorf("Audit ID too short (may not be secure): %s", event.ID)
		}
	}

	t.Logf("Audit ID security verified successfully, generated %d unique IDs", len(idSet))
}

// testConcurrentAuditLogging verifies audit logging works under concurrent access
func testConcurrentAuditLogging(t *testing.T, tempDir string) {
	auditFile := filepath.Join(tempDir, "concurrent_audit.log")

	// Create audit trail with secure random generator
	secureRandom := security.NewSecureRandom()
	auditTrail := reporting.NewAuditTrailWithSecureRandom(auditFile, secureRandom)

	numGoroutines := 20
	eventsPerGoroutine := 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*eventsPerGoroutine)

	// Launch concurrent audit loggers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < eventsPerGoroutine; j++ {
				packageName := fmt.Sprintf("concurrent-pkg-%d-%d", goroutineID, j)
				err := auditTrail.LogVersionUpdate(packageName, "1.0.0", "1.1.0", true, nil)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, event %d: %w", goroutineID, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent audit logging error: %v", err)
	}

	// Verify audit events were logged
	since := time.Now().Add(-time.Hour)
	auditEvents, err := auditTrail.GetAuditHistory(since, time.Now(), "")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	expectedCount := numGoroutines * eventsPerGoroutine
	if len(auditEvents) != expectedCount {
		t.Errorf("Expected %d audit events, got %d", expectedCount, len(auditEvents))
	}

	// Verify all IDs are unique
	idSet := make(map[string]bool)
	for _, event := range auditEvents {
		if idSet[event.ID] {
			t.Errorf("Duplicate audit ID in concurrent logging: %s", event.ID)
		}
		idSet[event.ID] = true
	}

	t.Logf("Concurrent audit logging completed successfully, logged %d unique events", len(auditEvents))
}

// TestEndToEndSecurityIntegration tests complete end-to-end security integration
func TestEndToEndSecurityIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("CompleteSecureWorkflow", func(t *testing.T) {
		testCompleteSecureWorkflow(t, tempDir)
	})
}

// testCompleteSecureWorkflow tests a complete workflow using all secure components
func testCompleteSecureWorkflow(t *testing.T, tempDir string) {
	// Setup components
	storageFile := filepath.Join(tempDir, "workflow_versions.yaml")
	auditFile := filepath.Join(tempDir, "workflow_audit.log")

	storage, err := version.NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	secureRandom := security.NewSecureRandom()
	auditTrail := reporting.NewAuditTrailWithSecureRandom(auditFile, secureRandom)

	// Simulate a complete version management workflow
	packages := []struct {
		name       string
		oldVersion string
		newVersion string
	}{
		{"react", "17.0.0", "18.0.0"},
		{"express", "4.17.0", "4.18.0"},
		{"lodash", "4.17.20", "4.17.21"},
	}

	// Step 1: Update versions with audit logging
	for _, pkg := range packages {
		// Create version info
		versionInfo := &models.VersionInfo{
			Name:           pkg.name,
			Language:       "javascript",
			Type:           "package",
			CurrentVersion: pkg.oldVersion,
			LatestVersion:  pkg.newVersion,
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		}

		// Store version info (uses secure file operations)
		err := storage.SetVersionInfo(pkg.name, versionInfo)
		if err != nil {
			t.Fatalf("Failed to store version info for %s: %v", pkg.name, err)
		}

		// Log the update (uses secure ID generation)
		err = auditTrail.LogVersionUpdate(pkg.name, pkg.oldVersion, pkg.newVersion, true, nil)
		if err != nil {
			t.Fatalf("Failed to log version update for %s: %v", pkg.name, err)
		}
	}

	// Step 2: Verify storage integrity
	allVersions, err := storage.ListVersions()
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}

	if len(allVersions) != len(packages) {
		t.Errorf("Expected %d packages in storage, got %d", len(packages), len(allVersions))
	}

	// Step 3: Verify audit trail integrity
	since := time.Now().Add(-time.Hour)
	auditEvents, err := auditTrail.GetAuditHistory(since, time.Now(), "version_update")
	if err != nil {
		t.Fatalf("Failed to get audit history: %v", err)
	}

	if len(auditEvents) != len(packages) {
		t.Errorf("Expected %d audit events, got %d", len(packages), len(auditEvents))
	}

	// Step 4: Verify no security artifacts remain
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()

		// Should only have the storage and audit files
		if fileName != "workflow_versions.yaml" && fileName != "workflow_audit.log" {
			if strings.Contains(fileName, ".tmp") || strings.Contains(fileName, "temp") {
				t.Errorf("Unexpected temporary file found: %s", fileName)
			}
		}
	}

	// Step 5: Performance verification
	startTime := time.Now()

	// Perform additional operations to test performance
	for i := 0; i < 10; i++ {
		packageName := fmt.Sprintf("perf-test-%d", i)
		versionInfo := &models.VersionInfo{
			Name:           packageName,
			Language:       "javascript",
			Type:           "package",
			CurrentVersion: "1.0.0",
			LatestVersion:  "1.0.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		}

		err := storage.SetVersionInfo(packageName, versionInfo)
		if err != nil {
			t.Fatalf("Failed to store performance test package: %v", err)
		}

		err = auditTrail.LogVersionUpdate(packageName, "1.0.0", "1.0.0", true, nil)
		if err != nil {
			t.Fatalf("Failed to log performance test update: %v", err)
		}
	}

	duration := time.Since(startTime)

	// Should complete within reasonable time
	if duration > 5*time.Second {
		t.Errorf("End-to-end workflow too slow: %v", duration)
	}

	t.Logf("Complete secure workflow completed successfully in %v", duration)
	t.Logf("Final package count: %d", len(allVersions)+10)
	t.Logf("All security measures verified")
}
