package version

import (
	"fmt"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// MockTemplateUpdater implements TemplateUpdater for testing
type MockTemplateUpdater struct {
	updatedTemplates  []string
	affectedTemplates []string
	backupCreated     bool
	restoreCalled     bool
	shouldFail        bool
	failBackup        bool
}

func NewMockTemplateUpdater() *MockTemplateUpdater {
	return &MockTemplateUpdater{
		updatedTemplates:  make([]string, 0),
		affectedTemplates: []string{"templates/frontend/nextjs-app", "templates/backend/go-gin"},
	}
}

func (m *MockTemplateUpdater) UpdateTemplate(templatePath string, versions map[string]*models.VersionInfo) error {
	if m.shouldFail {
		return fmt.Errorf("mock template update failure")
	}
	m.updatedTemplates = append(m.updatedTemplates, templatePath)
	return nil
}

func (m *MockTemplateUpdater) UpdateAllTemplates(versions map[string]*models.VersionInfo) error {
	if m.shouldFail {
		return fmt.Errorf("mock template update failure")
	}
	m.updatedTemplates = append(m.updatedTemplates, m.affectedTemplates...)
	return nil
}

func (m *MockTemplateUpdater) ValidateTemplate(templatePath string) error {
	return nil
}

func (m *MockTemplateUpdater) GetAffectedTemplates(versions map[string]*models.VersionInfo) ([]string, error) {
	return m.affectedTemplates, nil
}

func (m *MockTemplateUpdater) BackupTemplates(templatePaths []string) error {
	if m.failBackup {
		return fmt.Errorf("mock backup failure")
	}
	m.backupCreated = true
	return nil
}

func (m *MockTemplateUpdater) RestoreTemplates(templatePaths []string) error {
	m.restoreCalled = true
	return nil
}

func TestUpdatePipelineExecute_NoUpdates(t *testing.T) {
	// Setup mocks
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	templateUpdater := NewMockTemplateUpdater()

	// Add current version that's already up to date
	currentInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "2.0.0",
		LatestVersion:  "2.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	}
	storage.SetVersionInfo("react", currentInfo)

	// Setup registries
	registries := map[string]interfaces.VersionRegistry{
		"mock": NewMockVersionRegistry(),
	}

	// Create pipeline
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		MaxRetries:             3,
		RetryDelay:             1 * time.Second,
	}

	manager := NewManagerWithStorage(cache, storage)
	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Error("Expected pipeline to succeed")
	}

	if result.UpdatesDetected != 0 {
		t.Errorf("Expected 0 updates detected, got %d", result.UpdatesDetected)
	}

	if result.UpdatesApplied != 0 {
		t.Errorf("Expected 0 updates applied, got %d", result.UpdatesApplied)
	}
}

func TestUpdatePipelineExecute_WithUpdates(t *testing.T) {
	// Setup mocks
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	templateUpdater := NewMockTemplateUpdater()

	// Add current version that needs updating
	currentInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	}
	storage.SetVersionInfo("react", currentInfo)

	// Setup registries with newer version
	mockRegistry := NewMockVersionRegistry()
	mockRegistry.SetVersion("react", &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		LatestVersion:  "2.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	})

	registries := map[string]interfaces.VersionRegistry{
		"mock": mockRegistry,
	}

	// Create pipeline
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		MaxRetries:             3,
		RetryDelay:             1 * time.Second,
	}

	manager := NewManagerWithStorage(cache, storage)
	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Error("Expected pipeline to succeed")
	}

	if result.UpdatesDetected != 1 {
		t.Errorf("Expected 1 update detected, got %d", result.UpdatesDetected)
	}

	if result.UpdatesApplied != 1 {
		t.Errorf("Expected 1 update applied, got %d", result.UpdatesApplied)
	}

	if result.TemplatesUpdated != 2 {
		t.Errorf("Expected 2 templates updated, got %d", result.TemplatesUpdated)
	}

	if !templateUpdater.backupCreated {
		t.Error("Expected backup to be created")
	}

	if len(templateUpdater.updatedTemplates) == 0 {
		t.Error("Expected templates to be updated")
	}
}

func TestUpdatePipelineExecute_WithSecurityUpdates(t *testing.T) {
	// Setup mocks
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	templateUpdater := NewMockTemplateUpdater()

	// Add current version with security issues
	currentInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.0.0",
		IsSecure:       false,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: []models.SecurityIssue{
			{
				ID:          "CVE-2023-1234",
				Severity:    "high",
				Description: "Test security issue",
				FixedIn:     "2.0.0",
				ReportedAt:  time.Now(),
			},
		},
	}
	storage.SetVersionInfo("react", currentInfo)

	// Setup registries with security fix
	mockRegistry := NewMockVersionRegistry()
	mockRegistry.SetVersion("react", &models.VersionInfo{
		Name:          "react",
		Language:      "javascript",
		Type:          "package",
		LatestVersion: "2.0.0",
		IsSecure:      true,
		UpdatedAt:     time.Now(),
		CheckedAt:     time.Now(),
		UpdateSource:  "mock",
		SecurityIssues: []models.SecurityIssue{
			{
				ID:          "CVE-2023-1234",
				Severity:    "high",
				Description: "Test security issue",
				FixedIn:     "2.0.0",
				ReportedAt:  time.Now(),
			},
		},
	})

	registries := map[string]interfaces.VersionRegistry{
		"mock": mockRegistry,
	}

	// Create pipeline with security priority
	config := &PipelineConfig{
		AutoUpdate:             false, // Disabled to test security priority
		SecurityPriority:       true,
		BreakingChangeApproval: true,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		MaxRetries:             3,
		RetryDelay:             1 * time.Second,
	}

	manager := NewManagerWithStorage(cache, storage)
	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results
	if !result.Success {
		t.Error("Expected pipeline to succeed")
	}

	if result.SecurityUpdates != 1 {
		t.Errorf("Expected 1 security update, got %d", result.SecurityUpdates)
	}

	if result.UpdatesApplied != 1 {
		t.Errorf("Expected 1 update applied (security), got %d", result.UpdatesApplied)
	}
}

func TestUpdatePipelineExecute_WithRollback(t *testing.T) {
	// Setup mocks
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	templateUpdater := NewMockTemplateUpdater()
	templateUpdater.shouldFail = true // Force failure to test rollback

	// Add current version that needs updating
	currentInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	}
	storage.SetVersionInfo("react", currentInfo)

	// Setup registries with newer version
	mockRegistry := NewMockVersionRegistry()
	mockRegistry.SetVersion("react", &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		LatestVersion:  "2.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	})

	registries := map[string]interfaces.VersionRegistry{
		"mock": mockRegistry,
	}

	// Create pipeline with rollback enabled
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: false,
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		MaxRetries:             2, // Reduced for faster test
		RetryDelay:             100 * time.Millisecond,
	}

	manager := NewManagerWithStorage(cache, storage)
	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute pipeline (should fail and rollback)
	result, err := pipeline.Execute()
	if err == nil {
		t.Error("Expected pipeline to fail due to mock failure")
	}

	// Verify rollback was performed
	if !result.RollbackPerformed {
		t.Error("Expected rollback to be performed")
	}

	if !templateUpdater.restoreCalled {
		t.Error("Expected restore to be called")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors to be recorded")
	}
}

func TestUpdatePipelineExecute_BreakingChangeApproval(t *testing.T) {
	// Setup mocks
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	templateUpdater := NewMockTemplateUpdater()

	// Add current version that would have breaking change
	currentInfo := &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		CurrentVersion: "1.9.0",
		LatestVersion:  "1.9.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	}
	storage.SetVersionInfo("react", currentInfo)

	// Setup registries with major version bump
	mockRegistry := NewMockVersionRegistry()
	mockRegistry.SetVersion("react", &models.VersionInfo{
		Name:           "react",
		Language:       "javascript",
		Type:           "package",
		LatestVersion:  "2.0.0", // Major version bump
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	})

	registries := map[string]interfaces.VersionRegistry{
		"mock": mockRegistry,
	}

	// Create pipeline with breaking change approval required
	config := &PipelineConfig{
		AutoUpdate:             true,
		SecurityPriority:       true,
		BreakingChangeApproval: true, // Require approval for breaking changes
		BackupEnabled:          true,
		RollbackOnFailure:      true,
		MaxRetries:             3,
		RetryDelay:             1 * time.Second,
	}

	manager := NewManagerWithStorage(cache, storage)
	pipeline := NewUpdatePipeline(manager, storage, templateUpdater, registries, config)

	// Execute pipeline
	result, err := pipeline.Execute()
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Verify results - breaking change should be detected but not applied
	if !result.Success {
		t.Error("Expected pipeline to succeed")
	}

	if result.BreakingChanges != 1 {
		t.Errorf("Expected 1 breaking change detected, got %d", result.BreakingChanges)
	}

	if result.UpdatesApplied != 0 {
		t.Errorf("Expected 0 updates applied (breaking change blocked), got %d", result.UpdatesApplied)
	}
}

// Ensure MockTemplateUpdater implements the interface
var _ interfaces.TemplateUpdater = (*MockTemplateUpdater)(nil)
