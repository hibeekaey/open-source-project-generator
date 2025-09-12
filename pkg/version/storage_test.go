package version

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewFileStorage(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		filePath    string
		format      string
		expectError bool
	}{
		{
			name:        "valid yaml format",
			filePath:    filepath.Join(tempDir, "test.yaml"),
			format:      "yaml",
			expectError: false,
		},
		{
			name:        "valid json format",
			filePath:    filepath.Join(tempDir, "test.json"),
			format:      "json",
			expectError: false,
		},
		{
			name:        "invalid format",
			filePath:    filepath.Join(tempDir, "test.xml"),
			format:      "xml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewFileStorage(tt.filePath, tt.format)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if storage == nil {
				t.Errorf("expected storage instance but got nil")
				return
			}

			// Verify file was created
			if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
				t.Errorf("expected file to be created at %s", tt.filePath)
			}
		})
	}
}

func TestFileStorage_SaveAndLoad(t *testing.T) {
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
		Languages: map[string]*models.VersionInfo{
			"go": {
				Name:           "go",
				Language:       "go",
				Type:           "language",
				CurrentVersion: "1.21.0",
				LatestVersion:  "1.22.0",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
				CheckedAt:      time.Now(),
				UpdateSource:   "golang.org",
			},
		},
		Frameworks: map[string]*models.VersionInfo{
			"nextjs": {
				Name:           "nextjs",
				Language:       "javascript",
				Type:           "framework",
				CurrentVersion: "14.0.0",
				LatestVersion:  "15.0.0",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
				CheckedAt:      time.Now(),
				UpdateSource:   "npm",
			},
		},
		Packages: make(map[string]*models.VersionInfo),
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:       true,
			SecurityPriority: true,
			UpdateSchedule:   "daily",
		},
	}

	// Test Save
	err = storage.Save(testStore)
	if err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	// Test Load
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}

	// Verify loaded data
	if loadedStore.Version != testStore.Version {
		t.Errorf("expected version %s, got %s", testStore.Version, loadedStore.Version)
	}

	if len(loadedStore.Languages) != len(testStore.Languages) {
		t.Errorf("expected %d languages, got %d", len(testStore.Languages), len(loadedStore.Languages))
	}

	goInfo, exists := loadedStore.Languages["go"]
	if !exists {
		t.Errorf("expected 'go' language info to exist")
	} else {
		if goInfo.CurrentVersion != "1.21.0" {
			t.Errorf("expected go current version 1.21.0, got %s", goInfo.CurrentVersion)
		}
		if goInfo.LatestVersion != "1.22.0" {
			t.Errorf("expected go latest version 1.22.0, got %s", goInfo.LatestVersion)
		}
	}
}

func TestFileStorage_GetSetVersionInfo(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Test version info
	testInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "framework",
		CurrentVersion: "18.0.0",
		LatestVersion:  "19.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	// Test SetVersionInfo
	err = storage.SetVersionInfo("react", testInfo)
	if err != nil {
		t.Fatalf("failed to set version info: %v", err)
	}

	// Test GetVersionInfo
	retrievedInfo, err := storage.GetVersionInfo("react")
	if err != nil {
		t.Fatalf("failed to get version info: %v", err)
	}

	if retrievedInfo.Name != testInfo.Name {
		t.Errorf("expected name %s, got %s", testInfo.Name, retrievedInfo.Name)
	}
	if retrievedInfo.CurrentVersion != testInfo.CurrentVersion {
		t.Errorf("expected current version %s, got %s", testInfo.CurrentVersion, retrievedInfo.CurrentVersion)
	}
	if retrievedInfo.LatestVersion != testInfo.LatestVersion {
		t.Errorf("expected latest version %s, got %s", testInfo.LatestVersion, retrievedInfo.LatestVersion)
	}

	// Test GetVersionInfo for non-existent package
	_, err = storage.GetVersionInfo("nonexistent")
	if err == nil {
		t.Errorf("expected error for non-existent package")
	}
}

func TestFileStorage_DeleteVersionInfo(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Add test version info
	testInfo := &models.VersionInfo{
		Name:           "vue",
		Language:       "javascript",
		Type:           "framework",
		CurrentVersion: "3.0.0",
		LatestVersion:  "3.4.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	err = storage.SetVersionInfo("vue", testInfo)
	if err != nil {
		t.Fatalf("failed to set version info: %v", err)
	}

	// Verify it exists
	_, err = storage.GetVersionInfo("vue")
	if err != nil {
		t.Fatalf("version info should exist before deletion: %v", err)
	}

	// Delete it
	err = storage.DeleteVersionInfo("vue")
	if err != nil {
		t.Fatalf("failed to delete version info: %v", err)
	}

	// Verify it's gone
	_, err = storage.GetVersionInfo("vue")
	if err == nil {
		t.Errorf("expected error after deletion")
	}

	// Test deleting non-existent package
	err = storage.DeleteVersionInfo("nonexistent")
	if err == nil {
		t.Errorf("expected error when deleting non-existent package")
	}
}

func TestFileStorage_ListVersions(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Add multiple version infos
	testInfos := map[string]*models.VersionInfo{
		"go": {
			Name:           "go",
			Language:       "go",
			Type:           "language",
			CurrentVersion: "1.21.0",
			LatestVersion:  "1.22.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "golang.org",
		},
		"react": {
			Name:           "react",
			Language:       "javascript",
			Type:           "framework",
			CurrentVersion: "18.0.0",
			LatestVersion:  "19.0.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		},
		"lodash": {
			Name:           "lodash",
			Language:       "javascript",
			Type:           "package",
			CurrentVersion: "4.17.20",
			LatestVersion:  "4.17.21",
			IsSecure:       false,
			SecurityIssues: []models.SecurityIssue{
				{
					ID:          "CVE-2021-23337",
					Severity:    "high",
					Description: "Command injection vulnerability",
					FixedIn:     "4.17.21",
					ReportedAt:  time.Now(),
				},
			},
			UpdatedAt:    time.Now(),
			CheckedAt:    time.Now(),
			UpdateSource: "npm",
		},
	}

	for name, info := range testInfos {
		err = storage.SetVersionInfo(name, info)
		if err != nil {
			t.Fatalf("failed to set version info for %s: %v", name, err)
		}
	}

	// Test ListVersions
	allVersions, err := storage.ListVersions()
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(allVersions) != len(testInfos) {
		t.Errorf("expected %d versions, got %d", len(testInfos), len(allVersions))
	}

	for name, expectedInfo := range testInfos {
		actualInfo, exists := allVersions[name]
		if !exists {
			t.Errorf("expected version info for %s to exist", name)
			continue
		}
		if actualInfo.CurrentVersion != expectedInfo.CurrentVersion {
			t.Errorf("expected current version %s for %s, got %s",
				expectedInfo.CurrentVersion, name, actualInfo.CurrentVersion)
		}
	}
}

func TestFileStorage_Query(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Add test data
	testInfos := map[string]*models.VersionInfo{
		"go": {
			Name:           "go",
			Language:       "go",
			Type:           "language",
			CurrentVersion: "1.21.0",
			LatestVersion:  "1.22.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "golang.org",
		},
		"react": {
			Name:           "react",
			Language:       "javascript",
			Type:           "framework",
			CurrentVersion: "19.0.0",
			LatestVersion:  "19.0.0",
			IsSecure:       true,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		},
		"lodash": {
			Name:           "lodash",
			Language:       "javascript",
			Type:           "package",
			CurrentVersion: "4.17.20",
			LatestVersion:  "4.17.21",
			IsSecure:       false,
			UpdatedAt:      time.Now(),
			CheckedAt:      time.Now(),
			UpdateSource:   "npm",
		},
	}

	for name, info := range testInfos {
		err = storage.SetVersionInfo(name, info)
		if err != nil {
			t.Fatalf("failed to set version info for %s: %v", name, err)
		}
	}

	tests := []struct {
		name          string
		query         *models.VersionQuery
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "query by language",
			query:         &models.VersionQuery{Language: "javascript"},
			expectedCount: 2,
			expectedNames: []string{"react", "lodash"},
		},
		{
			name:          "query by type",
			query:         &models.VersionQuery{Type: "framework"},
			expectedCount: 1,
			expectedNames: []string{"react"},
		},
		{
			name:          "query outdated",
			query:         &models.VersionQuery{Outdated: true},
			expectedCount: 2,
			expectedNames: []string{"go", "lodash"},
		},
		{
			name:          "query insecure",
			query:         &models.VersionQuery{Insecure: true},
			expectedCount: 1,
			expectedNames: []string{"lodash"},
		},
		{
			name:          "query by name pattern",
			query:         &models.VersionQuery{Name: "re"},
			expectedCount: 1,
			expectedNames: []string{"react"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := storage.Query(tt.query)
			if err != nil {
				t.Fatalf("failed to query: %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
			}

			for _, expectedName := range tt.expectedNames {
				if _, exists := results[expectedName]; !exists {
					t.Errorf("expected %s in results", expectedName)
				}
			}
		})
	}
}

func TestFileStorage_BackupAndRestore(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.yaml")

	storage, err := NewFileStorage(filePath, "yaml")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Add test data
	testInfo := &models.VersionInfo{
		Name:           "test-package",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "2.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	err = storage.SetVersionInfo("test-package", testInfo)
	if err != nil {
		t.Fatalf("failed to set version info: %v", err)
	}

	// Test Backup
	err = storage.Backup()
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	// Verify backup file exists
	backupFiles, err := filepath.Glob(filepath.Join(storage.backupDir, "versions_backup_*.yaml"))
	if err != nil {
		t.Fatalf("failed to list backup files: %v", err)
	}
	if len(backupFiles) == 0 {
		t.Fatalf("no backup files found")
	}

	// Modify original data
	modifiedInfo := &models.VersionInfo{
		Name:           "test-package",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "2.0.0",
		LatestVersion:  "3.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
	}

	err = storage.SetVersionInfo("test-package", modifiedInfo)
	if err != nil {
		t.Fatalf("failed to modify version info: %v", err)
	}

	// Verify modification
	retrievedInfo, err := storage.GetVersionInfo("test-package")
	if err != nil {
		t.Fatalf("failed to get modified version info: %v", err)
	}
	if retrievedInfo.CurrentVersion != "2.0.0" {
		t.Errorf("expected modified current version 2.0.0, got %s", retrievedInfo.CurrentVersion)
	}

	// Test Restore
	err = storage.Restore(backupFiles[0])
	if err != nil {
		t.Fatalf("failed to restore from backup: %v", err)
	}

	// Verify restoration
	restoredInfo, err := storage.GetVersionInfo("test-package")
	if err != nil {
		t.Fatalf("failed to get restored version info: %v", err)
	}
	if restoredInfo.CurrentVersion != "1.0.0" {
		t.Errorf("expected restored current version 1.0.0, got %s", restoredInfo.CurrentVersion)
	}
}
