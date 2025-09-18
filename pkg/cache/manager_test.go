package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// setupTestCache creates a temporary directory and cache manager for testing
func setupTestCache(t *testing.T) (string, interfaces.CacheManager, func()) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temp directory: %v", err)
		}
	}

	manager := NewManager(tempDir)
	return tempDir, manager, cleanup
}

func TestCacheManager_BasicOperations(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Test Set and Get
	key := "test_key"
	value := "test_value"
	ttl := time.Hour

	err := manager.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set cache entry: %v", err)
	}

	retrievedValue, err := manager.Get(key)
	if err != nil {
		t.Fatalf("Failed to get cache entry: %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected %v, got %v", value, retrievedValue)
	}

	// Test Exists
	if !manager.Exists(key) {
		t.Error("Key should exist in cache")
	}

	// Test Delete
	err = manager.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete cache entry: %v", err)
	}

	if manager.Exists(key) {
		t.Error("Key should not exist after deletion")
	}
}

func TestCacheManager_TTL(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Set entry with short TTL
	key := "ttl_test"
	value := "test_value"
	ttl := 100 * time.Millisecond

	err := manager.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set cache entry: %v", err)
	}

	// Should exist immediately
	if !manager.Exists(key) {
		t.Error("Key should exist immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	if manager.Exists(key) {
		t.Error("Key should not exist after TTL expiration")
	}

	// Get should return error for expired key
	_, err = manager.Get(key)
	if err == nil {
		t.Error("Get should return error for expired key")
	}
}

func TestCacheManager_Stats(t *testing.T) {
	tempDir, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add some entries
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		err := manager.Set(key, value, time.Hour)
		if err != nil {
			t.Fatalf("Failed to set cache entry: %v", err)
		}
	}

	// Get stats
	stats, err := manager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get cache stats: %v", err)
	}

	if stats.TotalEntries != 5 {
		t.Errorf("Expected 5 entries, got %d", stats.TotalEntries)
	}

	if stats.TotalSize <= 0 {
		t.Error("Total size should be greater than 0")
	}

	if stats.CacheLocation != tempDir {
		t.Errorf("Expected cache location %s, got %s", tempDir, stats.CacheLocation)
	}
}

func TestCacheManager_Clean(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add entries with different TTLs
	err := manager.Set("persistent", "value1", time.Hour)
	if err != nil {
		t.Fatalf("Failed to set persistent entry: %v", err)
	}

	err = manager.Set("expiring", "value2", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to set expiring entry: %v", err)
	}

	// Wait for one to expire
	time.Sleep(100 * time.Millisecond)

	// Clean cache
	err = manager.Clean()
	if err != nil {
		t.Fatalf("Failed to clean cache: %v", err)
	}

	// Persistent entry should still exist
	if !manager.Exists("persistent") {
		t.Error("Persistent entry should still exist after cleaning")
	}

	// Expired entry should be removed
	if manager.Exists("expiring") {
		t.Error("Expired entry should be removed after cleaning")
	}
}

func TestCacheManager_Clear(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add some entries
	err := manager.Set("key1", "value1", time.Hour)
	if err != nil {
		t.Fatalf("Failed to set entry: %v", err)
	}

	err = manager.Set("key2", "value2", time.Hour)
	if err != nil {
		t.Fatalf("Failed to set entry: %v", err)
	}

	// Clear cache
	err = manager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// No entries should exist
	if manager.Exists("key1") || manager.Exists("key2") {
		t.Error("No entries should exist after clearing cache")
	}

	stats, err := manager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", stats.TotalEntries)
	}
}

func TestCacheManager_OfflineMode(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Initially not in offline mode
	if manager.IsOfflineMode() {
		t.Error("Should not be in offline mode initially")
	}

	// Enable offline mode
	err := manager.EnableOfflineMode()
	if err != nil {
		t.Fatalf("Failed to enable offline mode: %v", err)
	}

	if !manager.IsOfflineMode() {
		t.Error("Should be in offline mode after enabling")
	}

	// Disable offline mode
	err = manager.DisableOfflineMode()
	if err != nil {
		t.Fatalf("Failed to disable offline mode: %v", err)
	}

	if manager.IsOfflineMode() {
		t.Error("Should not be in offline mode after disabling")
	}
}

func TestCacheManager_Validation(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add some valid entries
	err := manager.Set("valid_key", "valid_value", time.Hour)
	if err != nil {
		t.Fatalf("Failed to set valid entry: %v", err)
	}

	// Validate cache
	err = manager.ValidateCache()
	if err != nil {
		t.Fatalf("Cache validation should pass: %v", err)
	}
}

func TestCacheManager_GetKeys(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add some entries
	expectedKeys := []string{"key1", "key2", "key3"}
	for _, key := range expectedKeys {
		err := manager.Set(key, "value", time.Hour)
		if err != nil {
			t.Fatalf("Failed to set entry: %v", err)
		}
	}

	// Get all keys
	keys, err := manager.GetKeys()
	if err != nil {
		t.Fatalf("Failed to get keys: %v", err)
	}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	// Check that all expected keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("Expected key %s not found", expectedKey)
		}
	}
}

func TestCacheManager_GetKeysByPattern(t *testing.T) {
	_, manager, cleanup := setupTestCache(t)
	defer cleanup()

	// Add entries with different patterns
	entries := map[string]string{
		"user:1":     "user1",
		"user:2":     "user2",
		"config:app": "config",
		"temp:data":  "temp",
	}

	for key, value := range entries {
		err := manager.Set(key, value, time.Hour)
		if err != nil {
			t.Fatalf("Failed to set entry: %v", err)
		}
	}

	// Get keys matching pattern
	userKeys, err := manager.GetKeysByPattern("^user:")
	if err != nil {
		t.Fatalf("Failed to get keys by pattern: %v", err)
	}

	if len(userKeys) != 2 {
		t.Errorf("Expected 2 user keys, got %d", len(userKeys))
	}

	// Check that returned keys match pattern
	for _, key := range userKeys {
		if !strings.HasPrefix(key, "user:") {
			t.Errorf("Key %s does not match pattern", key)
		}
	}
}

// Enhanced tests from manager_enhanced_test.go

func TestNewManager(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}

	// Verify cache directory is set
	if manager.GetLocation() != tempDir {
		t.Errorf("Expected cache location '%s', got '%s'", tempDir, manager.GetLocation())
	}
}

func TestManager_TTLHandling(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	key := "ttl-test"
	value := "ttl-value"
	shortTTL := 100 * time.Millisecond

	// Set with short TTL
	err := manager.Set(key, value, shortTTL)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should exist immediately
	if !manager.Exists(key) {
		t.Error("Expected key to exist immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, err = manager.Get(key)
	if err == nil {
		t.Error("Expected error when getting expired key")
	}

	if manager.Exists(key) {
		t.Error("Expected expired key to not exist")
	}
}

func TestManager_SetTTL(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	key := "ttl-update-test"
	value := "test-value"

	// Set without TTL
	err := manager.Set(key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Set TTL
	newTTL := 5 * time.Minute
	err = manager.SetTTL(key, newTTL)
	if err != nil {
		t.Fatalf("SetTTL failed: %v", err)
	}

	// Get TTL
	retrievedTTL, err := manager.GetTTL(key)
	if err != nil {
		t.Fatalf("GetTTL failed: %v", err)
	}

	// Should be close to the set TTL (allowing for small time differences)
	if retrievedTTL < 4*time.Minute || retrievedTTL > 5*time.Minute {
		t.Errorf("Expected TTL around 5 minutes, got %v", retrievedTTL)
	}

	// Test RefreshTTL
	err = manager.RefreshTTL(key)
	if err != nil {
		t.Fatalf("RefreshTTL failed: %v", err)
	}
}

func TestManager_GetExpiredKeys(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Set keys with different TTLs
	err := manager.Set("active", "value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set active failed: %v", err)
	}

	err = manager.Set("expired", "value", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Set expired failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Get expired keys
	expiredKeys, err := manager.GetExpiredKeys()
	if err != nil {
		t.Fatalf("GetExpiredKeys failed: %v", err)
	}

	if len(expiredKeys) != 1 {
		t.Errorf("Expected 1 expired key, got %d", len(expiredKeys))
	}

	if expiredKeys[0] != "expired" {
		t.Errorf("Expected expired key 'expired', got '%s'", expiredKeys[0])
	}
}

func TestManager_GetSize(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Initially should be empty
	size, err := manager.GetSize()
	if err != nil {
		t.Fatalf("GetSize failed: %v", err)
	}

	if size != 0 {
		t.Errorf("Expected initial size 0, got %d", size)
	}

	// Add some data
	err = manager.Set("key1", "small value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = manager.Set("key2", "another small value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Size should increase
	newSize, err := manager.GetSize()
	if err != nil {
		t.Fatalf("GetSize failed: %v", err)
	}

	if newSize <= size {
		t.Error("Expected size to increase after adding data")
	}
}

func TestManager_RepairCache(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Add some data
	err := manager.Set("key1", "value1", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Repair cache (should work even with valid cache)
	err = manager.RepairCache()
	if err != nil {
		t.Fatalf("RepairCache failed: %v", err)
	}

	// Verify data still exists
	if !manager.Exists("key1") {
		t.Error("Expected key to exist after repair")
	}
}

func TestManager_CompactCache(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Add data with different TTLs
	err := manager.Set("active", "value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set active failed: %v", err)
	}

	err = manager.Set("expired", "value", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Set expired failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Compact cache
	err = manager.CompactCache()
	if err != nil {
		t.Fatalf("CompactCache failed: %v", err)
	}

	// Active key should remain, expired should be removed
	if !manager.Exists("active") {
		t.Error("Expected active key to remain after compact")
	}

	if manager.Exists("expired") {
		t.Error("Expected expired key to be removed after compact")
	}
}

func TestManager_BackupAndRestore(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Add some data
	err := manager.Set("key1", "value1", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = manager.Set("key2", "value2", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "cache_backup.json")
	err = manager.BackupCache(backupPath)
	if err != nil {
		t.Fatalf("BackupCache failed: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Expected backup file to be created")
	}

	// Clear cache
	err = manager.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify cache is empty
	keys, err := manager.GetKeys()
	if err != nil {
		t.Fatalf("GetKeys failed: %v", err)
	}

	if len(keys) != 0 {
		t.Error("Expected cache to be empty after clear")
	}

	// Restore from backup
	err = manager.RestoreCache(backupPath)
	if err != nil {
		t.Fatalf("RestoreCache failed: %v", err)
	}

	// Verify data is restored
	if !manager.Exists("key1") {
		t.Error("Expected key1 to be restored")
	}

	if !manager.Exists("key2") {
		t.Error("Expected key2 to be restored")
	}
}

func TestManager_SyncCache(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Add some data
	err := manager.Set("sync-key", "sync-value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Sync cache
	err = manager.SyncCache()
	if err != nil {
		t.Fatalf("SyncCache failed: %v", err)
	}
}

func TestManager_CacheConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Get default config
	config, err := manager.GetCacheConfig()
	if err != nil {
		t.Fatalf("GetCacheConfig failed: %v", err)
	}

	if config == nil {
		t.Fatal("Expected cache config, got nil")
	}

	// Modify config
	newConfig := *config
	newConfig.MaxSize = 2048 * 1024 * 1024 // 2GB
	newConfig.DefaultTTL = 12 * time.Hour

	// Set new config
	err = manager.SetCacheConfig(&newConfig)
	if err != nil {
		t.Fatalf("SetCacheConfig failed: %v", err)
	}

	// Verify config was updated
	updatedConfig, err := manager.GetCacheConfig()
	if err != nil {
		t.Fatalf("GetCacheConfig after update failed: %v", err)
	}

	if updatedConfig.MaxSize != newConfig.MaxSize {
		t.Errorf("Expected MaxSize %d, got %d", newConfig.MaxSize, updatedConfig.MaxSize)
	}

	if updatedConfig.DefaultTTL != newConfig.DefaultTTL {
		t.Errorf("Expected DefaultTTL %v, got %v", newConfig.DefaultTTL, updatedConfig.DefaultTTL)
	}
}

func TestManager_SetMaxSize(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	newMaxSize := int64(1024 * 1024) // 1MB

	err := manager.SetMaxSize(newMaxSize)
	if err != nil {
		t.Fatalf("SetMaxSize failed: %v", err)
	}

	config, err := manager.GetCacheConfig()
	if err != nil {
		t.Fatalf("GetCacheConfig failed: %v", err)
	}

	if config.MaxSize != newMaxSize {
		t.Errorf("Expected MaxSize %d, got %d", newMaxSize, config.MaxSize)
	}
}

func TestManager_SetDefaultTTL(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	newTTL := 2 * time.Hour

	err := manager.SetDefaultTTL(newTTL)
	if err != nil {
		t.Fatalf("SetDefaultTTL failed: %v", err)
	}

	config, err := manager.GetCacheConfig()
	if err != nil {
		t.Fatalf("GetCacheConfig failed: %v", err)
	}

	if config.DefaultTTL != newTTL {
		t.Errorf("Expected DefaultTTL %v, got %v", newTTL, config.DefaultTTL)
	}
}

func TestManager_Callbacks(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	var hitKey, missKey, evictedKey string
	var evictionReason string

	// Set callbacks
	manager.OnCacheHit(func(key string) {
		hitKey = key
	})

	manager.OnCacheMiss(func(key string) {
		missKey = key
	})

	manager.OnCacheEviction(func(key string, reason string) {
		evictedKey = key
		evictionReason = reason
	})

	// Test hit callback
	err := manager.Set("hit-test", "value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	_, err = manager.Get("hit-test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if hitKey != "hit-test" {
		t.Errorf("Expected hit key 'hit-test', got '%s'", hitKey)
	}

	// Test miss callback
	_, err = manager.Get("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	if missKey != "non-existent" {
		t.Errorf("Expected miss key 'non-existent', got '%s'", missKey)
	}

	// Test eviction callback
	err = manager.Delete("hit-test")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if evictedKey != "hit-test" {
		t.Errorf("Expected evicted key 'hit-test', got '%s'", evictedKey)
	}

	if evictionReason != interfaces.EvictionReasonManual {
		t.Errorf("Expected eviction reason '%s', got '%s'", interfaces.EvictionReasonManual, evictionReason)
	}
}

func TestManager_GetHitAndMissRates(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Initially should be 0
	if manager.GetHitRate() != 0 {
		t.Error("Expected initial hit rate to be 0")
	}

	if manager.GetMissRate() != 0 {
		t.Error("Expected initial miss rate to be 0")
	}

	// Add data and access it
	err := manager.Set("test-key", "test-value", 5*time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Hit
	_, err = manager.Get("test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Miss
	_, err = manager.Get("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	// Check rates
	hitRate := manager.GetHitRate()
	missRate := manager.GetMissRate()

	if hitRate <= 0 {
		t.Error("Expected positive hit rate")
	}

	if missRate <= 0 {
		t.Error("Expected positive miss rate")
	}

	if hitRate+missRate != 1.0 {
		t.Errorf("Expected hit rate + miss rate = 1.0, got %f", hitRate+missRate)
	}
}

func TestManager_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Test nil config
	err := manager.SetCacheConfig(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	// Test invalid backup path
	err = manager.BackupCache("../../../invalid/path")
	if err == nil {
		t.Error("Expected error for invalid backup path")
	}

	// Test invalid restore path
	err = manager.RestoreCache("../../../invalid/path")
	if err == nil {
		t.Error("Expected error for invalid restore path")
	}

	// Test TTL operations on non-existent key
	_, err = manager.GetTTL("non-existent")
	if err == nil {
		t.Error("Expected error for GetTTL on non-existent key")
	}

	err = manager.SetTTL("non-existent", 5*time.Minute)
	if err == nil {
		t.Error("Expected error for SetTTL on non-existent key")
	}

	err = manager.RefreshTTL("non-existent")
	if err == nil {
		t.Error("Expected error for RefreshTTL on non-existent key")
	}
}

// Benchmark tests
func BenchmarkManager_Set(b *testing.B) {
	tempDir := b.TempDir()
	manager := NewManager(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := manager.Set(key, "benchmark-value", 5*time.Minute)
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}
}

func BenchmarkManager_Get(b *testing.B) {
	tempDir := b.TempDir()
	manager := NewManager(tempDir)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := manager.Set(key, "benchmark-value", 5*time.Minute)
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		_, err := manager.Get(key)
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkManager_GetStats(b *testing.B) {
	tempDir := b.TempDir()
	manager := NewManager(tempDir)

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := manager.Set(key, "benchmark-value", 5*time.Minute)
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GetStats()
		if err != nil {
			b.Fatalf("GetStats failed: %v", err)
		}
	}
}
