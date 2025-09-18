package cache

import (
	"fmt"
	"os"
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
