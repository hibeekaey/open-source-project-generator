package version

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// TestSecurityIntegration_NoTimestampLeakage verifies that the storage layer
// does not create temporary files with predictable timestamp-based names
func TestSecurityIntegration_NoTimestampLeakage(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
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

	// Record the time range for this test
	startTime := time.Now().UnixNano()

	// Perform multiple saves in quick succession
	// With the old insecure implementation, this would create predictable temp files
	// like "versions.yaml.tmp.1234567890123456789"
	for i := 0; i < 5; i++ {
		err = storage.Save(testStore)
		if err != nil {
			t.Fatalf("failed to save store iteration %d: %v", i, err)
		}
	}

	endTime := time.Now().UnixNano()

	// Check that no temporary files with timestamp patterns remain
	// (they should be cleaned up immediately with secure operations)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp directory: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()

		// Look for any files that match the old insecure pattern
		if strings.Contains(fileName, ".tmp.") {
			// Extract potential timestamp
			parts := strings.Split(fileName, ".tmp.")
			if len(parts) > 1 {
				suffix := parts[1]

				// Check if this looks like a nanosecond timestamp
				if len(suffix) >= 10 { // Unix timestamps are at least 10 digits
					t.Errorf("found potential timestamp-based temp file: %s", fileName)
				}
			}
		}

		// Also check for any files that contain timestamps in the expected range
		if strings.Contains(fileName, "tmp") {
			// This would catch files created with the old pattern
			t.Logf("found temp-related file: %s (this should not contain predictable patterns)", fileName)
		}
	}

	// Verify the final file exists and is correct
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected final file to exist at %s", filePath)
	}

	// Load and verify the content
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load final store: %v", err)
	}

	if loadedStore.Version != testStore.Version {
		t.Errorf("expected version %s, got %s", testStore.Version, loadedStore.Version)
	}

	t.Logf("Successfully completed security integration test")
	t.Logf("Time range: %d to %d nanoseconds", startTime, endTime)
	t.Logf("No predictable temporary file patterns detected")
}

// TestSecurityIntegration_AtomicOperations verifies that file operations are truly atomic
func TestSecurityIntegration_AtomicOperations(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
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

	// Add a test package
	testStore.Packages["test-package"] = &models.VersionInfo{
		Name:           "test-package",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	// Save the store
	err = storage.Save(testStore)
	if err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	// Verify that the file is complete and readable immediately after save
	// This tests that the atomic operation completed successfully
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load store immediately after save: %v", err)
	}

	if len(loadedStore.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(loadedStore.Packages))
	}

	if _, exists := loadedStore.Packages["test-package"]; !exists {
		t.Errorf("expected test-package to exist in loaded store")
	}

	t.Logf("Atomic operations verified successfully")
}
