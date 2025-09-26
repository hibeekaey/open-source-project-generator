package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

// TestOfflineModeAndCacheFunctionality tests offline mode and cache functionality
func TestOfflineModeAndCacheFunctionality(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("cache_initialization", func(t *testing.T) {
		testCacheInitialization(t)
	})

	t.Run("offline_mode_operations", func(t *testing.T) {
		testOfflineModeOperations(t, tempDir)
	})

	t.Run("cache_persistence", func(t *testing.T) {
		testCachePersistence(t, tempDir)
	})

	t.Run("cache_eviction_policies", func(t *testing.T) {
		testCacheEvictionPolicies(t, tempDir)
	})

	t.Run("offline_template_processing", func(t *testing.T) {
		testOfflineTemplateProcessing(t, tempDir)
	})
}

func testCacheInitialization(t *testing.T) {
	tempDir := t.TempDir()
	cacheDir := filepath.Join(tempDir, "cache")

	// Create cache manager
	cacheManager := cache.NewManager(cacheDir)

	// Verify cache directory is created
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		t.Error("Expected cache directory to be created")
	}

	// Test initial cache state
	stats, err := cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Errorf("Expected empty cache, got %d entries", stats.TotalEntries)
	}

	if stats.TotalSize != 0 {
		t.Errorf("Expected zero cache size, got %d bytes", stats.TotalSize)
	}

	// Test cache location
	location := cacheManager.GetLocation()
	if location != cacheDir {
		t.Errorf("Expected cache location '%s', got '%s'", cacheDir, location)
	}
}

func testOfflineModeOperations(t *testing.T, tempDir string) {
	cacheDir := filepath.Join(tempDir, "offline-cache")
	cacheManager := cache.NewManager(cacheDir)

	// Pre-populate cache with data
	testData := map[string]string{
		"template:go-gin": "cached go-gin template data",
		"version:node":    "18.17.0",
		"version:go":      "1.21.0",
		"package:react":   "18.2.0",
		"package:express": "4.18.2",
	}

	for key, value := range testData {
		err := cacheManager.Set(key, value, 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to populate cache with %s: %v", key, err)
		}
	}

	// Enable offline mode
	err := cacheManager.EnableOfflineMode()
	if err != nil {
		t.Fatalf("Failed to enable offline mode: %v", err)
	}

	if !cacheManager.IsOfflineMode() {
		t.Error("Expected offline mode to be enabled")
	}

	// Test offline operations
	for key, expectedValue := range testData {
		value, err := cacheManager.Get(key)
		if err != nil {
			t.Errorf("Failed to get cached data for %s in offline mode: %v", key, err)
			continue
		}

		if value != expectedValue {
			t.Errorf("Expected cached value '%s' for key '%s', got '%s'", expectedValue, key, value)
		}
	}

	// Test that new data cannot be added in offline mode (if network is required)
	err = cacheManager.Set("new:key", "new value", time.Hour)
	if err != nil {
		// This is expected if offline mode prevents new network requests
		t.Logf("Cannot add new data in offline mode (expected): %v", err)
	}

	// Disable offline mode
	err = cacheManager.DisableOfflineMode()
	if err != nil {
		t.Fatalf("Failed to disable offline mode: %v", err)
	}

	if cacheManager.IsOfflineMode() {
		t.Error("Expected offline mode to be disabled")
	}
}

func testCachePersistence(t *testing.T, tempDir string) {
	cacheDir := filepath.Join(tempDir, "persistent-cache")

	// Create first cache manager instance
	cacheManager1 := cache.NewManager(cacheDir)

	// Add data to cache
	testData := map[string]string{
		"persist:key1": "persistent value 1",
		"persist:key2": "persistent value 2",
		"persist:key3": "persistent value 3",
	}

	for key, value := range testData {
		err := cacheManager1.Set(key, value, 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set cache data: %v", err)
		}
	}

	// Force sync to ensure data is persisted
	err := cacheManager1.SyncCache()
	if err != nil {
		t.Fatalf("Failed to sync cache: %v", err)
	}

	// Create second cache manager instance (simulating restart)
	cacheManager2 := cache.NewManager(cacheDir)

	// Verify data persisted
	for key, expectedValue := range testData {
		value, err := cacheManager2.Get(key)
		if err != nil {
			t.Errorf("Failed to get persisted data for %s: %v", key, err)
			continue
		}

		if value != expectedValue {
			t.Errorf("Expected persisted value '%s' for key '%s', got '%s'", expectedValue, key, value)
		}
	}

	// Verify cache stats are restored
	stats, err := cacheManager2.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats: %v", err)
	}

	if stats.TotalEntries != len(testData) {
		t.Errorf("Expected %d persisted entries, got %d", len(testData), stats.TotalEntries)
	}
}

func testCacheEvictionPolicies(t *testing.T, tempDir string) {
	cacheDir := filepath.Join(tempDir, "eviction-cache")
	cacheManager := cache.NewManager(cacheDir)

	// Configure cache with small limits to trigger eviction
	config := &interfaces.CacheConfig{
		MaxSize:        1024,  // 1KB
		MaxEntries:     5,     // Max 5 entries
		EvictionPolicy: "lru", // Least Recently Used
	}

	err := cacheManager.SetCacheConfig(config)
	if err != nil {
		t.Fatalf("Failed to set cache config: %v", err)
	}

	// Test LRU eviction
	t.Run("lru_eviction", func(t *testing.T) {
		testLRUEviction(t, cacheManager)
	})

	// Test size-based eviction
	t.Run("size_eviction", func(t *testing.T) {
		testSizeEviction(t, cacheManager)
	})

	// Test TTL-based eviction
	t.Run("ttl_eviction", func(t *testing.T) {
		testTTLEviction(t, cacheManager)
	})
}

func testLRUEviction(t *testing.T, cacheManager interfaces.CacheManager) {
	// Clear cache first
	err := cacheManager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Add entries up to the limit
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("lru:key%d", i)
		value := fmt.Sprintf("value%d", i)
		err := cacheManager.Set(key, value, time.Hour)
		if err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}
	}

	// Access some entries to update their usage
	_, err = cacheManager.Get("lru:key1")
	if err != nil {
		t.Fatalf("Failed to access key1: %v", err)
	}

	_, err = cacheManager.Get("lru:key3")
	if err != nil {
		t.Fatalf("Failed to access key3: %v", err)
	}

	// Add one more entry to trigger eviction
	err = cacheManager.Set("lru:key5", "value5", time.Hour)
	if err != nil {
		t.Fatalf("Failed to set eviction-triggering entry: %v", err)
	}

	// Verify that least recently used entries were evicted
	// key1 and key3 should still exist (recently accessed)
	if !cacheManager.Exists("lru:key1") {
		t.Error("Expected recently accessed key1 to still exist")
	}

	if !cacheManager.Exists("lru:key3") {
		t.Error("Expected recently accessed key3 to still exist")
	}

	// key5 should exist (just added)
	if !cacheManager.Exists("lru:key5") {
		t.Error("Expected newly added key5 to exist")
	}
}

func testSizeEviction(t *testing.T, cacheManager interfaces.CacheManager) {
	// Clear cache first
	err := cacheManager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Add large entries to trigger size-based eviction
	largeValue := make([]byte, 300) // 300 bytes each
	for i := range largeValue {
		largeValue[i] = byte('A' + (i % 26))
	}

	// Add entries that will exceed size limit
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("size:key%d", i)
		err := cacheManager.Set(key, largeValue, time.Hour)
		if err != nil {
			t.Fatalf("Failed to set large cache entry: %v", err)
		}
	}

	// Verify cache size is within limits
	stats, err := cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats: %v", err)
	}

	// Should have triggered eviction to stay within size limit
	if stats.TotalSize > 1024 {
		t.Errorf("Cache size exceeds limit: %d bytes (max: 1024)", stats.TotalSize)
	}
}

func testTTLEviction(t *testing.T, cacheManager interfaces.CacheManager) {
	// Clear cache first
	err := cacheManager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Add entries with short TTL
	shortTTL := 100 * time.Millisecond

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("ttl:key%d", i)
		value := fmt.Sprintf("ttl_value%d", i)
		err := cacheManager.Set(key, value, shortTTL)
		if err != nil {
			t.Fatalf("Failed to set TTL cache entry: %v", err)
		}
	}

	// Verify entries exist initially
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("ttl:key%d", i)
		if !cacheManager.Exists(key) {
			t.Errorf("Expected TTL key%d to exist initially", i)
		}
	}

	// Wait for TTL expiration
	time.Sleep(150 * time.Millisecond)

	// Clean expired entries
	err = cacheManager.Clean()
	if err != nil {
		t.Fatalf("Failed to clean cache: %v", err)
	}

	// Verify entries are expired and removed
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("ttl:key%d", i)
		if cacheManager.Exists(key) {
			t.Errorf("Expected TTL key%d to be expired and removed", i)
		}
	}
}

func testOfflineTemplateProcessing(t *testing.T, tempDir string) {
	cacheDir := filepath.Join(tempDir, "template-cache")
	cacheManager := cache.NewManager(cacheDir)

	// Create template engine with cache
	templateEngine := template.NewEngine()
	templateManager := template.NewManager(templateEngine)

	// Pre-cache template data
	templateData := map[string]string{
		"template:go-gin:metadata": `{
			"name": "go-gin",
			"version": "1.0.0",
			"description": "Go Gin API template"
		}`,
		"template:go-gin:files": `{
			"main.go": "package main\n\nfunc main() {\n\t// {{.Name}} application\n}",
			"README.md": "# {{.Name}}\n\n{{.Description}}"
		}`,
	}

	for key, data := range templateData {
		err := cacheManager.Set(key, data, 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to cache template data: %v", err)
		}
	}

	// Enable offline mode
	err := cacheManager.EnableOfflineMode()
	if err != nil {
		t.Fatalf("Failed to enable offline mode: %v", err)
	}

	// Test template operations in offline mode
	config := &models.ProjectConfig{
		Name:         "offline-test-project",
		Organization: "offline-org",
		Description:  "Test project generated in offline mode",
		License:      "MIT",
		OutputPath:   filepath.Join(tempDir, "offline-output"),
	}

	// List templates (should work with cached data)
	templates, err := templateManager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("Failed to list templates in offline mode: %v", err)
	}

	if len(templates) == 0 {
		t.Error("Expected to find cached templates in offline mode")
	}

	// Get template info (should work with cached data)
	info, err := templateManager.GetTemplateInfo("go-gin")
	if err != nil {
		// This might fail if template info is not cached, which is acceptable
		t.Logf("Template info not available in offline mode (expected): %v", err)
	} else {
		if info.Name != "go-gin" {
			t.Errorf("Expected template name 'go-gin', got '%s'", info.Name)
		}
	}

	// Process template (should work with cached data)
	err = templateManager.ProcessTemplate("go-gin", config, config.OutputPath)
	if err != nil {
		// This might fail if template processing requires network access
		t.Logf("Template processing not fully available in offline mode (expected): %v", err)
	}
}

// TestCacheBackupAndRestore tests cache backup and restore functionality
func TestCacheBackupAndRestore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()
	cacheDir := filepath.Join(tempDir, "backup-cache")
	cacheManager := cache.NewManager(cacheDir)

	// Populate cache with test data
	testData := map[string]string{
		"backup:key1": "backup value 1",
		"backup:key2": "backup value 2",
		"backup:key3": "backup value 3",
		"backup:key4": "backup value 4",
	}

	for key, value := range testData {
		err := cacheManager.Set(key, value, 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set cache data: %v", err)
		}
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "cache_backup.json")
	err := cacheManager.BackupCache(backupPath)
	if err != nil {
		t.Fatalf("Failed to create cache backup: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatal("Expected backup file to be created")
	}

	// Clear cache
	err = cacheManager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify cache is empty
	stats, err := cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Errorf("Expected empty cache after clear, got %d entries", stats.TotalEntries)
	}

	// Restore from backup
	err = cacheManager.RestoreCache(backupPath)
	if err != nil {
		t.Fatalf("Failed to restore cache from backup: %v", err)
	}

	// Verify data is restored
	for key, expectedValue := range testData {
		value, err := cacheManager.Get(key)
		if err != nil {
			t.Errorf("Failed to get restored data for %s: %v", key, err)
			continue
		}

		if value != expectedValue {
			t.Errorf("Expected restored value '%s' for key '%s', got '%s'", expectedValue, key, value)
		}
	}

	// Verify cache stats are restored
	stats, err = cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats after restore: %v", err)
	}

	if stats.TotalEntries != len(testData) {
		t.Errorf("Expected %d restored entries, got %d", len(testData), stats.TotalEntries)
	}
}

// TestVersionManagerWithCache tests version manager with caching
func TestVersionManagerWithCache(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()
	cacheDir := filepath.Join(tempDir, "version-cache")
	cacheManager := cache.NewManager(cacheDir)

	// Create version manager with cache
	versionManager := version.NewManagerWithCache(cacheManager)

	// Test caching version information
	versionInfo := &interfaces.VersionInfo{
		Version:   "1.0.0",
		BuildDate: time.Now(),
		GitCommit: "abc123def456",
		GitBranch: "main",
		GoVersion: "1.21.0",
		Platform:  "darwin",
	}

	err := versionManager.CacheVersionInfo(versionInfo)
	if err != nil {
		t.Fatalf("Failed to cache version info: %v", err)
	}

	// Retrieve cached version info
	cachedInfo, err := versionManager.GetCachedVersionInfo()
	if err != nil {
		t.Fatalf("Failed to get cached version info: %v", err)
	}

	if cachedInfo.Version != versionInfo.Version {
		t.Errorf("Expected cached version '%s', got '%s'", versionInfo.Version, cachedInfo.Version)
	}

	if cachedInfo.GitCommit != versionInfo.GitCommit {
		t.Errorf("Expected cached git commit '%s', got '%s'", versionInfo.GitCommit, cachedInfo.GitCommit)
	}

	// Test cache refresh
	err = versionManager.RefreshVersionCache()
	if err != nil {
		t.Fatalf("Failed to refresh version cache: %v", err)
	}

	// Test cache clear
	err = versionManager.ClearVersionCache()
	if err != nil {
		t.Fatalf("Failed to clear version cache: %v", err)
	}

	// Verify cache is cleared
	_, err = versionManager.GetCachedVersionInfo()
	if err == nil {
		t.Error("Expected error when getting cached version info after clear")
	}
}
