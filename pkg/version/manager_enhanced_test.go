package version

import (
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// MockVersionStorage implements VersionStorage for testing
type MockVersionStorage struct {
	store *models.VersionStore
}

func NewMockVersionStorage() *MockVersionStorage {
	return &MockVersionStorage{
		store: &models.VersionStore{
			LastUpdated: time.Now(),
			Version:     "1.0.0",
			Languages:   make(map[string]*models.VersionInfo),
			Frameworks:  make(map[string]*models.VersionInfo),
			Packages:    make(map[string]*models.VersionInfo),
			UpdatePolicy: models.UpdatePolicy{
				AutoUpdate:             true,
				SecurityPriority:       true,
				BreakingChangeApproval: true,
				UpdateSchedule:         "daily",
				MaxAge:                 24 * time.Hour,
			},
		},
	}
}

func (m *MockVersionStorage) Load() (*models.VersionStore, error) {
	return m.store, nil
}

func (m *MockVersionStorage) Save(store *models.VersionStore) error {
	m.store = store
	return nil
}

func (m *MockVersionStorage) GetVersionInfo(name string) (*models.VersionInfo, error) {
	if info, exists := m.store.Languages[name]; exists {
		return info, nil
	}
	if info, exists := m.store.Frameworks[name]; exists {
		return info, nil
	}
	if info, exists := m.store.Packages[name]; exists {
		return info, nil
	}
	return nil, nil
}

func (m *MockVersionStorage) SetVersionInfo(name string, info *models.VersionInfo) error {
	switch info.Type {
	case "language":
		m.store.Languages[name] = info
	case "framework":
		m.store.Frameworks[name] = info
	case "package":
		m.store.Packages[name] = info
	}
	return nil
}

func (m *MockVersionStorage) DeleteVersionInfo(name string) error {
	delete(m.store.Languages, name)
	delete(m.store.Frameworks, name)
	delete(m.store.Packages, name)
	return nil
}

func (m *MockVersionStorage) ListVersions() (map[string]*models.VersionInfo, error) {
	result := make(map[string]*models.VersionInfo)
	for name, info := range m.store.Languages {
		result[name] = info
	}
	for name, info := range m.store.Frameworks {
		result[name] = info
	}
	for name, info := range m.store.Packages {
		result[name] = info
	}
	return result, nil
}

func (m *MockVersionStorage) Query(query *models.VersionQuery) (map[string]*models.VersionInfo, error) {
	return m.ListVersions()
}

func (m *MockVersionStorage) Backup() error {
	return nil
}

func (m *MockVersionStorage) Restore(backupPath string) error {
	return nil
}

func (m *MockVersionStorage) GetLastUpdated() (time.Time, error) {
	return m.store.LastUpdated, nil
}

func (m *MockVersionStorage) SetLastUpdated(timestamp time.Time) error {
	m.store.LastUpdated = timestamp
	return nil
}

// MockVersionRegistry implements VersionRegistry for testing
type MockVersionRegistry struct {
	versions map[string]*models.VersionInfo
}

func NewMockVersionRegistry() *MockVersionRegistry {
	return &MockVersionRegistry{
		versions: make(map[string]*models.VersionInfo),
	}
}

func (m *MockVersionRegistry) GetLatestVersion(packageName string) (*models.VersionInfo, error) {
	if info, exists := m.versions[packageName]; exists {
		return info, nil
	}

	// Return a default version for testing
	return &models.VersionInfo{
		Name:           packageName,
		Language:       "javascript",
		Type:           "package",
		LatestVersion:  "2.0.0",
		IsSecure:       true,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "mock",
		SecurityIssues: make([]models.SecurityIssue, 0),
	}, nil
}

func (m *MockVersionRegistry) GetVersionHistory(packageName string, limit int) ([]*models.VersionInfo, error) {
	latest, err := m.GetLatestVersion(packageName)
	if err != nil {
		return nil, err
	}
	return []*models.VersionInfo{latest}, nil
}

func (m *MockVersionRegistry) CheckSecurity(packageName, version string) ([]models.SecurityIssue, error) {
	return []models.SecurityIssue{}, nil
}

func (m *MockVersionRegistry) GetRegistryInfo() interfaces.RegistryInfo {
	return interfaces.RegistryInfo{
		Name:        "Mock Registry",
		URL:         "https://mock.registry",
		Type:        "mock",
		Description: "Mock registry for testing",
		Supported:   []string{"javascript"},
	}
}

func (m *MockVersionRegistry) IsAvailable() bool {
	return true
}

func (m *MockVersionRegistry) GetSupportedPackages() ([]string, error) {
	packages := make([]string, 0, len(m.versions))
	for name := range m.versions {
		packages = append(packages, name)
	}
	return packages, nil
}

func (m *MockVersionRegistry) SetVersion(name string, info *models.VersionInfo) {
	m.versions[name] = info
}

func TestManagerWithStorage_CheckLatestVersions(t *testing.T) {
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	manager := NewManagerWithStorage(cache, storage)

	// Add mock registry
	mockRegistry := NewMockVersionRegistry()
	manager.registries["npm"] = mockRegistry

	// Add test data to storage
	testInfo := &models.VersionInfo{
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
	storage.SetVersionInfo("react", testInfo)

	// Set up mock registry to return newer version
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

	report, err := manager.CheckLatestVersions()
	if err != nil {
		t.Fatalf("CheckLatestVersions failed: %v", err)
	}

	if report == nil {
		t.Fatal("Expected report, got nil")
	}

	if report.TotalPackages != 1 {
		t.Errorf("Expected 1 total package, got %d", report.TotalPackages)
	}

	if len(report.Recommendations) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(report.Recommendations))
	}

	if report.Recommendations[0].RecommendedVersion != "2.0.0" {
		t.Errorf("Expected recommended version 2.0.0, got %s", report.Recommendations[0].RecommendedVersion)
	}
}

func TestManagerWithStorage_UpdateVersionInfo(t *testing.T) {
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	manager := NewManagerWithStorage(cache, storage)

	// Add test data to storage
	testInfo := &models.VersionInfo{
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
	storage.SetVersionInfo("react", testInfo)

	// Update to newer version
	result, err := manager.UpdateVersionInfo("react", "2.0.0", false)
	if err != nil {
		t.Fatalf("UpdateVersionInfo failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected update to succeed")
	}

	if result.PreviousVersion != "1.0.0" {
		t.Errorf("Expected previous version 1.0.0, got %s", result.PreviousVersion)
	}

	if result.NewVersion != "2.0.0" {
		t.Errorf("Expected new version 2.0.0, got %s", result.NewVersion)
	}

	// Verify storage was updated
	updatedInfo, err := storage.GetVersionInfo("react")
	if err != nil {
		t.Fatalf("Failed to get updated version info: %v", err)
	}

	if updatedInfo.CurrentVersion != "2.0.0" {
		t.Errorf("Expected current version 2.0.0, got %s", updatedInfo.CurrentVersion)
	}

	if updatedInfo.PreviousVersion != "1.0.0" {
		t.Errorf("Expected previous version 1.0.0, got %s", updatedInfo.PreviousVersion)
	}
}

func TestManagerWithStorage_CompareVersions(t *testing.T) {
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	manager := NewManagerWithStorage(cache, storage)

	tests := []struct {
		version1 string
		version2 string
		expected int
	}{
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.1.0", "1.0.0", 1},
	}

	for _, test := range tests {
		result, err := manager.CompareVersions(test.version1, test.version2)
		if err != nil {
			t.Errorf("CompareVersions(%s, %s) failed: %v", test.version1, test.version2, err)
			continue
		}

		if result != test.expected {
			t.Errorf("CompareVersions(%s, %s) = %d, expected %d", test.version1, test.version2, result, test.expected)
		}
	}
}

func TestManagerWithStorage_DetectVersionUpdates(t *testing.T) {
	cache := NewMemoryCache(24 * time.Hour)
	storage := NewMockVersionStorage()
	manager := NewManagerWithStorage(cache, storage)

	// Add mock registry
	mockRegistry := NewMockVersionRegistry()
	manager.registries["npm"] = mockRegistry

	// Add test data to storage
	testInfo := &models.VersionInfo{
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
	storage.SetVersionInfo("react", testInfo)

	// Set up mock registry to return newer version
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

	updates, err := manager.DetectVersionUpdates()
	if err != nil {
		t.Fatalf("DetectVersionUpdates failed: %v", err)
	}

	if len(updates) != 1 {
		t.Errorf("Expected 1 update, got %d", len(updates))
	}

	if update, exists := updates["react"]; exists {
		if update.LatestVersion != "2.0.0" {
			t.Errorf("Expected latest version 2.0.0, got %s", update.LatestVersion)
		}
	} else {
		t.Error("Expected update for react package")
	}
}

// Ensure interfaces are implemented
var _ interfaces.VersionStorage = (*MockVersionStorage)(nil)
var _ interfaces.VersionRegistry = (*MockVersionRegistry)(nil)
