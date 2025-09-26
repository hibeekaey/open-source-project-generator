package version

import (
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// MockCacheManager is a simple mock implementation for testing
type MockCacheManager struct {
	data map[string]interface{}
}

func NewMockCacheManager() *MockCacheManager {
	return &MockCacheManager{
		data: make(map[string]interface{}),
	}
}

func (m *MockCacheManager) Get(key string) (interface{}, error) {
	if val, exists := m.data[key]; exists {
		return val, nil
	}
	return nil, interfaces.NewCLIError("cache", "key not found", 404)
}

func (m *MockCacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockCacheManager) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *MockCacheManager) Exists(key string) bool {
	_, exists := m.data[key]
	return exists
}

func (m *MockCacheManager) Clear() error {
	m.data = make(map[string]interface{})
	return nil
}

func (m *MockCacheManager) Clean() error {
	return nil
}

func (m *MockCacheManager) GetStats() (*interfaces.CacheStats, error) {
	return &interfaces.CacheStats{
		TotalEntries:  len(m.data),
		TotalSize:     1024,
		HitRate:       0.8,
		CacheLocation: "/tmp/test-cache",
		OfflineMode:   false,
		CacheHealth:   "healthy",
	}, nil
}

func (m *MockCacheManager) GetSize() (int64, error) {
	return 1024, nil
}

func (m *MockCacheManager) GetLocation() string {
	return "/tmp/test-cache"
}

func (m *MockCacheManager) GetKeys() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *MockCacheManager) GetKeysByPattern(pattern string) ([]string, error) {
	return m.GetKeys()
}

func (m *MockCacheManager) ValidateCache() error {
	return nil
}

func (m *MockCacheManager) RepairCache() error {
	return nil
}

func (m *MockCacheManager) CompactCache() error {
	return nil
}

func (m *MockCacheManager) BackupCache(path string) error {
	return nil
}

func (m *MockCacheManager) RestoreCache(path string) error {
	return nil
}

func (m *MockCacheManager) EnableOfflineMode() error {
	return nil
}

func (m *MockCacheManager) DisableOfflineMode() error {
	return nil
}

func (m *MockCacheManager) IsOfflineMode() bool {
	return false
}

func (m *MockCacheManager) SyncCache() error {
	return nil
}

func (m *MockCacheManager) SetTTL(key string, ttl time.Duration) error {
	return nil
}

func (m *MockCacheManager) GetTTL(key string) (time.Duration, error) {
	return time.Hour, nil
}

func (m *MockCacheManager) RefreshTTL(key string) error {
	return nil
}

func (m *MockCacheManager) GetExpiredKeys() ([]string, error) {
	return []string{}, nil
}

func (m *MockCacheManager) SetCacheConfig(config *interfaces.CacheConfig) error {
	return nil
}

func (m *MockCacheManager) GetCacheConfig() (*interfaces.CacheConfig, error) {
	return interfaces.DefaultCacheConfig(), nil
}

func (m *MockCacheManager) SetMaxSize(size int64) error {
	return nil
}

func (m *MockCacheManager) SetDefaultTTL(ttl time.Duration) error {
	return nil
}

func (m *MockCacheManager) OnCacheHit(callback func(key string)) {
}

func (m *MockCacheManager) OnCacheMiss(callback func(key string)) {
}

func (m *MockCacheManager) OnCacheEviction(callback func(key string, reason string)) {
}

func (m *MockCacheManager) GetHitRate() float64 {
	return 0.8
}

func (m *MockCacheManager) GetMissRate() float64 {
	return 0.2
}

func TestManager_GetCurrentVersion(t *testing.T) {
	manager := NewManagerWithVersion("1.0.0")

	version := manager.GetCurrentVersion()
	if version == "" {
		t.Error("Expected non-empty version")
	}

	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
}

func TestManager_WithCache_GetCurrentVersion(t *testing.T) {
	cacheManager := NewMockCacheManager()
	manager := NewManagerWithVersionAndCache("1.0.0", cacheManager)

	version := manager.GetCurrentVersion()
	if version == "" {
		t.Error("Expected non-empty version")
	}

	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
}

func TestManager_CheckForUpdates(t *testing.T) {
	manager := NewManager()

	updateInfo, err := manager.CheckForUpdates()
	if err != nil {
		// This might fail due to network issues in tests, which is acceptable
		t.Logf("CheckForUpdates failed (expected in test environment): %v", err)
		return
	}

	if updateInfo == nil {
		t.Error("Expected non-nil update info")
		return
	}

	if updateInfo.CurrentVersion == "" {
		t.Error("Expected non-empty current version")
	}
}

func TestManager_GetVersionConfig(t *testing.T) {
	manager := NewManager()

	config, err := manager.GetVersionConfig()
	if err != nil {
		t.Errorf("GetVersionConfig failed: %v", err)
		return
	}

	if config == nil {
		t.Error("Expected non-nil config")
		return
	}

	if config.UpdateChannel == "" {
		t.Error("Expected non-empty update channel")
	}
}

func TestManager_SetVersionConfig(t *testing.T) {
	manager := NewManager()

	config := DefaultVersionConfig()
	config.AutoUpdate = true
	config.UpdateChannel = "beta"

	err := manager.SetVersionConfig(config)
	if err != nil {
		t.Errorf("SetVersionConfig failed: %v", err)
		return
	}

	// Verify the config was set
	retrievedConfig, err := manager.GetVersionConfig()
	if err != nil {
		t.Errorf("GetVersionConfig failed after set: %v", err)
		return
	}

	if retrievedConfig.AutoUpdate != true {
		t.Error("Expected AutoUpdate to be true")
	}

	if retrievedConfig.UpdateChannel != "beta" {
		t.Errorf("Expected UpdateChannel to be 'beta', got '%s'", retrievedConfig.UpdateChannel)
	}
}

func TestManager_CheckCompatibility(t *testing.T) {
	manager := NewManager()

	result, err := manager.CheckCompatibility(".")
	if err != nil {
		t.Errorf("CheckCompatibility failed: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected non-nil compatibility result")
		return
	}

	// Basic validation of result structure
	if result.GeneratorVersion == "" {
		t.Error("Expected non-empty generator version")
	}
}

func TestManager_WithCache_CacheOperations(t *testing.T) {
	cacheManager := NewMockCacheManager()
	manager := NewManagerWithCache(cacheManager)

	// Test caching version info
	versionInfo := &interfaces.VersionInfo{
		Version:   "1.2.3",
		BuildDate: time.Now(),
		Metadata:  make(map[string]string),
	}

	err := manager.CacheVersionInfo(versionInfo)
	if err != nil {
		t.Errorf("CacheVersionInfo failed: %v", err)
	}

	// Test cache clearing
	err = manager.ClearVersionCache()
	if err != nil {
		t.Errorf("ClearVersionCache failed: %v", err)
	}
}

func TestManager_VersionComparison(t *testing.T) {
	manager := NewManager()

	// Test version comparison
	result := manager.compareVersions("1.0.0", "1.1.0")
	if result != -1 {
		t.Errorf("Expected -1 for 1.0.0 < 1.1.0, got %d", result)
	}

	result = manager.compareVersions("1.1.0", "1.0.0")
	if result != 1 {
		t.Errorf("Expected 1 for 1.1.0 > 1.0.0, got %d", result)
	}

	result = manager.compareVersions("1.0.0", "1.0.0")
	if result != 0 {
		t.Errorf("Expected 0 for 1.0.0 == 1.0.0, got %d", result)
	}
}

func TestManager_UpdateAvailability(t *testing.T) {
	manager := NewManager()

	// Test update availability
	available := manager.isUpdateAvailable("1.0.0", "1.1.0")
	if !available {
		t.Error("Expected update to be available for 1.0.0 -> 1.1.0")
	}

	available = manager.isUpdateAvailable("1.1.0", "1.0.0")
	if available {
		t.Error("Expected no update available for 1.1.0 -> 1.0.0")
	}
}

func TestManager_BreakingUpdateDetection(t *testing.T) {
	manager := NewManager()

	// Test breaking update detection
	breaking := manager.isBreakingUpdate("1.0.0", "2.0.0")
	if !breaking {
		t.Error("Expected breaking update for 1.0.0 -> 2.0.0")
	}

	breaking = manager.isBreakingUpdate("1.0.0", "1.1.0")
	if breaking {
		t.Error("Expected non-breaking update for 1.0.0 -> 1.1.0")
	}
}

func TestManager_UpdateType(t *testing.T) {
	manager := NewManager()

	// Test update type detection
	updateType := manager.getUpdateType("1.0.0", "2.0.0")
	if updateType != "major" {
		t.Errorf("Expected 'major' update type for 1.0.0 -> 2.0.0, got '%s'", updateType)
	}

	updateType = manager.getUpdateType("1.0.0", "1.1.0")
	if updateType != "minor" {
		t.Errorf("Expected 'minor' update type for 1.0.0 -> 1.1.0, got '%s'", updateType)
	}

	updateType = manager.getUpdateType("1.0.0", "1.0.1")
	if updateType != "patch" {
		t.Errorf("Expected 'patch' update type for 1.0.0 -> 1.0.1, got '%s'", updateType)
	}
}

// Additional enhanced tests

func TestManager_SetUpdateChannel(t *testing.T) {
	manager := NewManager()

	// Test setting valid channel
	err := manager.SetUpdateChannel("beta")
	if err != nil {
		t.Fatalf("SetUpdateChannel failed: %v", err)
	}

	channel := manager.GetUpdateChannel()
	if channel != "beta" {
		t.Errorf("Expected channel 'beta', got '%s'", channel)
	}

	// Reset to stable
	err = manager.SetUpdateChannel("stable")
	if err != nil {
		t.Fatalf("SetUpdateChannel reset failed: %v", err)
	}
}

func TestManager_SetAutoUpdate(t *testing.T) {
	manager := NewManager()

	// Enable auto update
	err := manager.SetAutoUpdate(true)
	if err != nil {
		t.Fatalf("SetAutoUpdate failed: %v", err)
	}

	config, err := manager.GetVersionConfig()
	if err != nil {
		t.Fatalf("GetVersionConfig failed: %v", err)
	}

	if !config.AutoUpdate {
		t.Error("Expected AutoUpdate to be enabled")
	}

	// Disable auto update
	err = manager.SetAutoUpdate(false)
	if err != nil {
		t.Fatalf("SetAutoUpdate disable failed: %v", err)
	}

	config, err = manager.GetVersionConfig()
	if err != nil {
		t.Fatalf("GetVersionConfig after disable failed: %v", err)
	}

	if config.AutoUpdate {
		t.Error("Expected AutoUpdate to be disabled")
	}
}

func TestManager_GetUpdateChannel(t *testing.T) {
	manager := NewManager()

	channel := manager.GetUpdateChannel()
	if channel == "" {
		t.Error("Expected update channel to be non-empty")
	}

	// Should be a valid channel
	validChannels := []string{"stable", "beta", "alpha", "nightly"}
	found := false
	for _, validChannel := range validChannels {
		if channel == validChannel {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected valid update channel, got '%s'", channel)
	}
}
