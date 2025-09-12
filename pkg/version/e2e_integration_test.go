package version

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// E2ETemplateUpdater implements TemplateUpdater interface for E2E testing
type E2ETemplateUpdater struct {
	affectedTemplates []string
	updateCalls       []map[string]*models.VersionInfo
	backupCalls       [][]string
	restoreCalls      [][]string
	validateCalls     []string
	shouldFailUpdate  bool
	shouldFailBackup  bool
	shouldFailRestore bool
}

func NewE2ETemplateUpdater() *E2ETemplateUpdater {
	return &E2ETemplateUpdater{
		affectedTemplates: []string{"templates/frontend/nextjs-app", "templates/frontend/nextjs-home"},
		updateCalls:       make([]map[string]*models.VersionInfo, 0),
		backupCalls:       make([][]string, 0),
		restoreCalls:      make([][]string, 0),
		validateCalls:     make([]string, 0),
	}
}

func (m *E2ETemplateUpdater) GetAffectedTemplates(updates map[string]*models.VersionInfo) ([]string, error) {
	return m.affectedTemplates, nil
}

func (m *E2ETemplateUpdater) UpdateAllTemplates(updates map[string]*models.VersionInfo) error {
	if m.shouldFailUpdate {
		return fmt.Errorf("mock template update failure")
	}
	m.updateCalls = append(m.updateCalls, updates)
	return nil
}

func (m *E2ETemplateUpdater) BackupTemplates(templatePaths []string) error {
	if m.shouldFailBackup {
		return fmt.Errorf("mock backup failure")
	}
	m.backupCalls = append(m.backupCalls, templatePaths)
	return nil
}

func (m *E2ETemplateUpdater) RestoreTemplates(backupPaths []string) error {
	if m.shouldFailRestore {
		return fmt.Errorf("mock restore failure")
	}
	m.restoreCalls = append(m.restoreCalls, backupPaths)
	return nil
}

func (m *E2ETemplateUpdater) ValidateTemplate(templatePath string) error {
	m.validateCalls = append(m.validateCalls, templatePath)
	return nil
}

func (m *E2ETemplateUpdater) UpdateTemplate(templatePath string, versions map[string]*models.VersionInfo) error {
	// Simple implementation for testing
	return nil
}

// E2EVersionRegistry implements VersionRegistry interface for E2E testing
type E2EVersionRegistry struct {
	name              string
	available         bool
	supportedPackages []string
	versions          map[string]*models.VersionInfo
	securityIssues    map[string][]models.SecurityIssue
}

func NewE2EVersionRegistry(name string) *E2EVersionRegistry {
	return &E2EVersionRegistry{
		name:              name,
		available:         true,
		supportedPackages: []string{"react", "next", "typescript"},
		versions:          make(map[string]*models.VersionInfo),
		securityIssues:    make(map[string][]models.SecurityIssue),
	}
}

func (m *E2EVersionRegistry) GetLatestVersion(packageName string) (*models.VersionInfo, error) {
	if info, exists := m.versions[packageName]; exists {
		return info, nil
	}
	return nil, fmt.Errorf("package not found: %s", packageName)
}

func (m *E2EVersionRegistry) CheckSecurity(packageName, version string) ([]models.SecurityIssue, error) {
	if issues, exists := m.securityIssues[packageName]; exists {
		return issues, nil
	}
	return []models.SecurityIssue{}, nil
}

func (m *E2EVersionRegistry) IsAvailable() bool {
	return m.available
}

func (m *E2EVersionRegistry) GetSupportedPackages() ([]string, error) {
	return m.supportedPackages, nil
}

func (m *E2EVersionRegistry) SetVersion(packageName string, info *models.VersionInfo) {
	m.versions[packageName] = info
}

func (m *E2EVersionRegistry) SetSecurityIssues(packageName string, issues []models.SecurityIssue) {
	m.securityIssues[packageName] = issues
}

func (m *E2EVersionRegistry) GetVersionHistory(packageName string, limit int) ([]*models.VersionInfo, error) {
	// Simple implementation for testing
	if info, exists := m.versions[packageName]; exists {
		return []*models.VersionInfo{info}, nil
	}
	return []*models.VersionInfo{}, nil
}

func (m *E2EVersionRegistry) GetRegistryInfo() interfaces.RegistryInfo {
	return interfaces.RegistryInfo{
		Name:        m.name,
		URL:         "https://test-registry.com",
		Type:        "test",
		Description: "Test registry for E2E testing",
		Supported:   []string{"javascript", "typescript"},
	}
}

// TestCompleteVersionUpdateWorkflow tests the complete end-to-end version update workflow
func TestCompleteVersionUpdateWorkflow(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	storageFile := filepath.Join(tempDir, "versions.yaml")

	// Initialize storage
	storage, err := NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize cache
	cache := NewMemoryCache(24 * time.Hour)

	// Create mock registries
	npmRegistry := NewE2EVersionRegistry("npm")
	npmRegistry.SetVersion("react", &models.VersionInfo{
		Name:           "react",
		CurrentVersion: "18.2.0",
		LatestVersion:  "19.1.0",
		Language:       "javascript",
		Type:           "framework",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})
	npmRegistry.SetVersion("next", &models.VersionInfo{
		Name:           "next",
		CurrentVersion: "14.0.0",
		LatestVersion:  "15.5.2",
		Language:       "javascript",
		Type:           "framework",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})
	npmRegistry.SetVersion("typescript", &models.VersionInfo{
		Name:           "typescript",
		CurrentVersion: "5.0.0",
		LatestVersion:  "5.3.3",
		Language:       "typescript",
		Type:           "package",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})

	registries := map[string]interfaces.VersionRegistry{
		"npm": npmRegistry,
	}

	// Create version manager
	manager := NewManagerWithStorage(cache, storage)

	// Initialize storage with current versions
	initialStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks: map[string]*models.VersionInfo{
			"react": {
				Name:           "react",
				CurrentVersion: "18.2.0",
				LatestVersion:  "18.2.0",
				Language:       "javascript",
				Type:           "framework",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
			"next": {
				Name:           "next",
				CurrentVersion: "14.0.0",
				LatestVersion:  "14.0.0",
				Language:       "javascript",
				Type:           "framework",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
		},
		Packages: map[string]*models.VersionInfo{
			"typescript": {
				Name:           "typescript",
				CurrentVersion: "5.0.0",
				LatestVersion:  "5.0.0",
				Language:       "typescript",
				Type:           "package",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
		},
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:             true,
			SecurityPriority:       true,
			BreakingChangeApproval: false, // Allow breaking changes for test
			UpdateSchedule:         "daily",
			MaxAge:                 24 * time.Hour,
		},
	}

	if err := storage.Save(initialStore); err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Create mock template updater
	templateUpdater := NewE2ETemplateUpdater()

	// Create update pipeline
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		UpdateSchedule:         "daily",
		MaxRetries:             3,
		RetryDelay:             100 * time.Millisecond,
		NotificationEnabled:    false,
	}

	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute the pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Errorf("Pipeline should have succeeded")
	}

	if result.UpdatesDetected != 3 {
		t.Errorf("Expected 3 updates detected, got %d", result.UpdatesDetected)
	}

	if result.UpdatesApplied != 3 {
		t.Errorf("Expected 3 updates applied, got %d", result.UpdatesApplied)
	}

	if result.TemplatesUpdated != 2 {
		t.Errorf("Expected 2 templates updated, got %d", result.TemplatesUpdated)
	}

	// Verify template updater was called correctly
	if len(templateUpdater.updateCalls) != 1 {
		t.Errorf("Expected 1 template update call, got %d", len(templateUpdater.updateCalls))
	}

	if len(templateUpdater.backupCalls) != 1 {
		t.Errorf("Expected 1 backup call, got %d", len(templateUpdater.backupCalls))
	}

	if len(templateUpdater.validateCalls) != 2 {
		t.Errorf("Expected 2 validation calls, got %d", len(templateUpdater.validateCalls))
	}

	// Verify template updater was called with correct versions
	if len(templateUpdater.updateCalls) == 0 {
		t.Error("Expected template updater to be called")
	} else {
		updates := templateUpdater.updateCalls[0]
		if updates["react"] == nil || updates["react"].LatestVersion != "19.1.0" {
			t.Errorf("React version not passed correctly to template updater")
		}
		if updates["next"] == nil || updates["next"].LatestVersion != "15.5.2" {
			t.Errorf("Next.js version not passed correctly to template updater")
		}
		if updates["typescript"] == nil || updates["typescript"].LatestVersion != "5.3.3" {
			t.Errorf("TypeScript version not passed correctly to template updater")
		}
	}

	t.Logf("✅ Complete version update workflow test passed")
}

// TestVersionUpdateWithSecurityIssues tests version updates with security vulnerabilities
func TestVersionUpdateWithSecurityIssues(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	storageFile := filepath.Join(tempDir, "versions.yaml")

	// Initialize storage
	storage, err := NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize cache
	cache := NewMemoryCache(24 * time.Hour)

	// Create mock registry with security issues
	npmRegistry := NewE2EVersionRegistry("npm")
	npmRegistry.supportedPackages = append(npmRegistry.supportedPackages, "vulnerable-package")
	npmRegistry.SetVersion("vulnerable-package", &models.VersionInfo{
		Name:           "vulnerable-package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.2.0",
		Language:       "javascript",
		Type:           "package",
		IsSecure:       false,
		SecurityIssues: []models.SecurityIssue{
			{
				ID:          "CVE-2023-1234",
				Severity:    "high",
				Description: "Remote code execution vulnerability",
				FixedIn:     "1.2.0",
			},
		},
		UpdatedAt: time.Now(),
	})

	npmRegistry.SetSecurityIssues("vulnerable-package", []models.SecurityIssue{
		{
			ID:          "CVE-2023-1234",
			Severity:    "high",
			Description: "Remote code execution vulnerability",
			FixedIn:     "1.2.0",
		},
	})

	registries := map[string]interfaces.VersionRegistry{
		"npm": npmRegistry,
	}

	// Create version manager
	manager := NewManagerWithStorage(cache, storage)

	// Initialize storage with vulnerable package
	initialStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks:  make(map[string]*models.VersionInfo),
		Packages: map[string]*models.VersionInfo{
			"vulnerable-package": {
				Name:           "vulnerable-package",
				CurrentVersion: "1.0.0",
				LatestVersion:  "1.0.0",
				Language:       "javascript",
				Type:           "package",
				IsSecure:       false,
				SecurityIssues: []models.SecurityIssue{
					{
						ID:          "CVE-2023-1234",
						Severity:    "high",
						Description: "Remote code execution vulnerability",
						FixedIn:     "1.2.0",
					},
				},
				UpdatedAt: time.Now(),
			},
		},
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:             false, // Disable auto-update
			SecurityPriority:       true,  // But enable security priority
			BreakingChangeApproval: true,
			UpdateSchedule:         "daily",
			MaxAge:                 24 * time.Hour,
		},
	}

	if err := storage.Save(initialStore); err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Create mock template updater
	templateUpdater := NewE2ETemplateUpdater()

	// Create update pipeline
	config := &PipelineConfig{
		AutoUpdate:             false, // Disabled
		SecurityPriority:       true,  // But security updates are prioritized
		BreakingChangeApproval: true,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		UpdateSchedule:         "daily",
		MaxRetries:             3,
		RetryDelay:             100 * time.Millisecond,
		NotificationEnabled:    false,
	}

	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute the pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Errorf("Pipeline should have succeeded")
	}

	if result.SecurityUpdates != 1 {
		t.Errorf("Expected 1 security update, got %d", result.SecurityUpdates)
	}

	if result.UpdatesApplied != 1 {
		t.Errorf("Expected 1 update applied (security), got %d", result.UpdatesApplied)
	}

	// Verify template updater was called with the security update
	if len(templateUpdater.updateCalls) == 0 {
		t.Error("Expected template updater to be called for security update")
	} else {
		updates := templateUpdater.updateCalls[0]
		if updates["vulnerable-package"] == nil || updates["vulnerable-package"].LatestVersion != "1.2.0" {
			t.Errorf("Vulnerable package not passed correctly to template updater")
		}
	}

	t.Logf("✅ Security update workflow test passed")
}

// TestVersionUpdateErrorHandlingAndRollback tests error handling and rollback functionality
func TestVersionUpdateErrorHandlingAndRollback(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	storageFile := filepath.Join(tempDir, "versions.yaml")

	// Initialize storage
	storage, err := NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize cache
	cache := NewMemoryCache(24 * time.Hour)

	// Create mock registry
	npmRegistry := NewE2EVersionRegistry("npm")
	npmRegistry.supportedPackages = append(npmRegistry.supportedPackages, "test-package")
	npmRegistry.SetVersion("test-package", &models.VersionInfo{
		Name:           "test-package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "2.0.0",
		Language:       "javascript",
		Type:           "package",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})

	registries := map[string]interfaces.VersionRegistry{
		"npm": npmRegistry,
	}

	// Create version manager
	manager := NewManagerWithStorage(cache, storage)

	// Initialize storage
	initialStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks:  make(map[string]*models.VersionInfo),
		Packages: map[string]*models.VersionInfo{
			"test-package": {
				Name:           "test-package",
				CurrentVersion: "1.0.0",
				LatestVersion:  "1.0.0",
				Language:       "javascript",
				Type:           "package",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
		},
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:             true,
			SecurityPriority:       true,
			BreakingChangeApproval: false,
			UpdateSchedule:         "daily",
			MaxAge:                 24 * time.Hour,
		},
	}

	if err := storage.Save(initialStore); err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Create mock template updater that will fail
	templateUpdater := NewE2ETemplateUpdater()
	templateUpdater.shouldFailUpdate = true // Force template update failure

	// Create update pipeline
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		UpdateSchedule:         "daily",
		MaxRetries:             2, // Reduced for faster test
		RetryDelay:             50 * time.Millisecond,
		NotificationEnabled:    false,
	}

	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute the pipeline (should fail and rollback)
	result, err := pipeline.Execute()
	if err == nil {
		t.Errorf("Expected pipeline to fail due to template update failure")
	}

	// Verify rollback was performed
	if !result.RollbackPerformed {
		t.Errorf("Expected rollback to be performed")
	}

	if len(result.Errors) == 0 {
		t.Errorf("Expected errors to be recorded")
	}

	// Verify backup was called
	if len(templateUpdater.backupCalls) != 1 {
		t.Errorf("Expected 1 backup call, got %d", len(templateUpdater.backupCalls))
	}

	// Verify restore was called
	if len(templateUpdater.restoreCalls) != 1 {
		t.Errorf("Expected 1 restore call, got %d", len(templateUpdater.restoreCalls))
	}

	t.Logf("✅ Error handling and rollback test passed")
}

// TestRegistryQueryToTemplateUpdatePipeline tests the complete pipeline from registry query to template update
func TestRegistryQueryToTemplateUpdatePipeline(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	storageFile := filepath.Join(tempDir, "versions.yaml")

	// Initialize storage
	storage, err := NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize cache
	cache := NewMemoryCache(24 * time.Hour)

	// Create multiple mock registries
	npmRegistry := NewE2EVersionRegistry("npm")
	npmRegistry.SetVersion("react", &models.VersionInfo{
		Name:           "react",
		CurrentVersion: "18.0.0",
		LatestVersion:  "19.1.0",
		Language:       "javascript",
		Type:           "framework",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})

	goRegistry := NewE2EVersionRegistry("go")
	goRegistry.supportedPackages = []string{"gin", "gorm"}
	goRegistry.SetVersion("gin", &models.VersionInfo{
		Name:           "gin",
		CurrentVersion: "1.9.0",
		LatestVersion:  "1.10.0",
		Language:       "go",
		Type:           "framework",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
	})

	registries := map[string]interfaces.VersionRegistry{
		"npm": npmRegistry,
		"go":  goRegistry,
	}

	// Create version manager
	manager := NewManagerWithStorage(cache, storage)

	// Initialize storage with packages from different registries
	initialStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks: map[string]*models.VersionInfo{
			"react": {
				Name:           "react",
				CurrentVersion: "18.0.0",
				LatestVersion:  "18.0.0",
				Language:       "javascript",
				Type:           "framework",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
			"gin": {
				Name:           "gin",
				CurrentVersion: "1.9.0",
				LatestVersion:  "1.9.0",
				Language:       "go",
				Type:           "framework",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
		},
		Packages: make(map[string]*models.VersionInfo),
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:             true,
			SecurityPriority:       true,
			BreakingChangeApproval: false,
			UpdateSchedule:         "daily",
			MaxAge:                 24 * time.Hour,
		},
	}

	if err := storage.Save(initialStore); err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Create mock template updater
	templateUpdater := NewE2ETemplateUpdater()
	templateUpdater.affectedTemplates = []string{
		"templates/frontend/nextjs-app",
		"templates/backend/go-gin",
	}

	// Create update pipeline
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		UpdateSchedule:         "daily",
		MaxRetries:             3,
		RetryDelay:             100 * time.Millisecond,
		NotificationEnabled:    false,
	}

	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute the pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify the complete pipeline worked
	if !result.Success {
		t.Errorf("Pipeline should have succeeded")
	}

	// Verify updates from multiple registries were detected
	if result.UpdatesDetected != 2 {
		t.Errorf("Expected 2 updates detected (from different registries), got %d", result.UpdatesDetected)
	}

	// Verify all updates were applied
	if result.UpdatesApplied != 2 {
		t.Errorf("Expected 2 updates applied, got %d", result.UpdatesApplied)
	}

	// Verify templates from different categories were updated
	if result.TemplatesUpdated != 2 {
		t.Errorf("Expected 2 templates updated, got %d", result.TemplatesUpdated)
	}

	// Verify the template updater received the correct updates
	if len(templateUpdater.updateCalls) != 1 {
		t.Errorf("Expected 1 template update call, got %d", len(templateUpdater.updateCalls))
	}

	updateCall := templateUpdater.updateCalls[0]
	if len(updateCall) != 2 {
		t.Errorf("Expected 2 packages in update call, got %d", len(updateCall))
	}

	// Verify specific packages were updated
	if _, exists := updateCall["react"]; !exists {
		t.Errorf("React update not found in template update call")
	}

	if _, exists := updateCall["gin"]; !exists {
		t.Errorf("Gin update not found in template update call")
	}

	t.Logf("✅ Registry query to template update pipeline test passed")
}

// TestConcurrentVersionUpdates tests concurrent version update operations
func TestConcurrentVersionUpdates(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	storageFile := filepath.Join(tempDir, "versions.yaml")

	// Initialize storage
	storage, err := NewFileStorage(storageFile, "yaml")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize cache
	cache := NewMemoryCache(24 * time.Hour)

	// Create version manager
	manager := NewManagerWithStorage(cache, storage)

	// Initialize storage with test packages
	initialStore := &models.VersionStore{
		LastUpdated: time.Now(),
		Version:     "1.0.0",
		Languages:   make(map[string]*models.VersionInfo),
		Frameworks:  make(map[string]*models.VersionInfo),
		Packages: map[string]*models.VersionInfo{
			"package1": {
				Name:           "package1",
				CurrentVersion: "1.0.0",
				LatestVersion:  "1.0.0",
				Language:       "javascript",
				Type:           "package",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
			"package2": {
				Name:           "package2",
				CurrentVersion: "2.0.0",
				LatestVersion:  "2.0.0",
				Language:       "javascript",
				Type:           "package",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
			},
		},
		UpdatePolicy: models.UpdatePolicy{
			AutoUpdate:             true,
			SecurityPriority:       true,
			BreakingChangeApproval: false,
			UpdateSchedule:         "daily",
			MaxAge:                 24 * time.Hour,
		},
	}

	if err := storage.Save(initialStore); err != nil {
		t.Fatalf("Failed to save initial store: %v", err)
	}

	// Test concurrent updates
	done := make(chan bool, 2)
	errors := make(chan error, 2)

	// Concurrent update 1
	go func() {
		result, err := manager.UpdateVersionInfo("package1", "1.1.0", false)
		if err != nil {
			errors <- err
			return
		}
		if !result.Success {
			errors <- fmt.Errorf("update 1 failed")
			return
		}
		done <- true
	}()

	// Concurrent update 2
	go func() {
		result, err := manager.UpdateVersionInfo("package2", "2.1.0", false)
		if err != nil {
			errors <- err
			return
		}
		if !result.Success {
			errors <- fmt.Errorf("update 2 failed")
			return
		}
		done <- true
	}()

	// Wait for both updates to complete
	completedCount := 0
	for completedCount < 2 {
		select {
		case <-done:
			completedCount++
		case err := <-errors:
			t.Fatalf("Concurrent update failed: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatalf("Concurrent updates timed out")
		}
	}

	// Verify both updates were applied
	updatedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load updated store: %v", err)
	}

	if updatedStore.Packages["package1"].CurrentVersion != "1.1.0" {
		t.Errorf("Package1 not updated correctly: %s", updatedStore.Packages["package1"].CurrentVersion)
	}

	if updatedStore.Packages["package2"].CurrentVersion != "2.1.0" {
		t.Errorf("Package2 not updated correctly: %s", updatedStore.Packages["package2"].CurrentVersion)
	}

	t.Logf("✅ Concurrent version updates test passed")
}
