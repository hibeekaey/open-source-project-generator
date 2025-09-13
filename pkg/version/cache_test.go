package version

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(1 * time.Hour)

	t.Run("set and get", func(t *testing.T) {
		err := cache.Set("test-key", "1.0.0")
		if err != nil {
			t.Errorf("Unexpected error setting cache: %v", err)
		}

		value, found := cache.Get("test-key")
		if !found {
			t.Errorf("Expected to find cached value")
		}
		if value != "1.0.0" {
			t.Errorf("Expected value 1.0.0, got %s", value)
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		_, found := cache.Get("non-existent")
		if found {
			t.Errorf("Expected not to find non-existent key")
		}
	})

	t.Run("delete key", func(t *testing.T) {
		cache.Set("delete-me", "1.0.0")

		err := cache.Delete("delete-me")
		if err != nil {
			t.Errorf("Unexpected error deleting key: %v", err)
		}

		_, found := cache.Get("delete-me")
		if found {
			t.Errorf("Expected key to be deleted")
		}
	})

	t.Run("clear cache", func(t *testing.T) {
		cache.Set("key1", "1.0.0")
		cache.Set("key2", "2.0.0")

		err := cache.Clear()
		if err != nil {
			t.Errorf("Unexpected error clearing cache: %v", err)
		}

		_, found1 := cache.Get("key1")
		_, found2 := cache.Get("key2")
		if found1 || found2 {
			t.Errorf("Expected all keys to be cleared")
		}
	})

	t.Run("keys", func(t *testing.T) {
		cache.Clear()
		cache.Set("key1", "1.0.0")
		cache.Set("key2", "2.0.0")

		keys := cache.Keys()
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		// Check that both keys are present
		keyMap := make(map[string]bool)
		for _, key := range keys {
			keyMap[key] = true
		}
		if !keyMap["key1"] || !keyMap["key2"] {
			t.Errorf("Expected keys key1 and key2, got %v", keys)
		}
	})

	t.Run("TTL expiration", func(t *testing.T) {
		shortCache := NewMemoryCache(10 * time.Millisecond)

		shortCache.Set("expire-me", "1.0.0")

		// Should be available immediately
		value, found := shortCache.Get("expire-me")
		if !found || value != "1.0.0" {
			t.Errorf("Expected to find value immediately after setting")
		}

		// Wait for expiration
		time.Sleep(20 * time.Millisecond)

		_, found = shortCache.Get("expire-me")
		if found {
			t.Errorf("Expected value to be expired")
		}
	})
}

func TestFileCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "version-cache-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cache, err := NewFileCache(tempDir, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create file cache: %v", err)
	}

	t.Run("set and get", func(t *testing.T) {
		err := cache.Set("test-key", "1.0.0")
		if err != nil {
			t.Errorf("Unexpected error setting cache: %v", err)
		}

		value, found := cache.Get("test-key")
		if !found {
			t.Errorf("Expected to find cached value")
		}
		if value != "1.0.0" {
			t.Errorf("Expected value 1.0.0, got %s", value)
		}
	})

	t.Run("persistence", func(t *testing.T) {
		// Set a value
		cache.Set("persist-key", "2.0.0")

		// Wait a bit for async save
		time.Sleep(100 * time.Millisecond)

		// Create a new cache instance with the same directory
		newCache, err := NewFileCache(tempDir, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to create new cache instance: %v", err)
		}

		// Should load the persisted value
		value, found := newCache.Get("persist-key")
		if !found {
			t.Errorf("Expected to find persisted value")
		}
		if value != "2.0.0" {
			t.Errorf("Expected persisted value 2.0.0, got %s", value)
		}
	})

	t.Run("clean expired", func(t *testing.T) {
		shortCache, err := NewFileCache(tempDir, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to create short TTL cache: %v", err)
		}

		shortCache.Set("expire1", "1.0.0")
		shortCache.Set("expire2", "2.0.0")

		// Wait for expiration
		time.Sleep(20 * time.Millisecond)

		removed := shortCache.CleanExpired()
		if removed != 2 {
			t.Errorf("Expected to remove 2 expired entries, removed %d", removed)
		}
	})

	t.Run("invalid cache directory", func(t *testing.T) {
		// Try to create cache in a file (should fail)
		invalidPath := filepath.Join(tempDir, "file.txt")
		if err := os.WriteFile(invalidPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err := NewFileCache(invalidPath, 1*time.Hour)
		if err == nil {
			t.Errorf("Expected error when creating cache with invalid directory")
		}
	})
}

func TestFileCacheLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "version-cache-load-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("load non-existent file", func(t *testing.T) {
		// Should not fail when cache file doesn't exist
		cache, err := NewFileCache(tempDir, 1*time.Hour)
		if err != nil {
			t.Errorf("Unexpected error creating cache with non-existent file: %v", err)
		}

		keys := cache.Keys()
		if len(keys) != 0 {
			t.Errorf("Expected empty cache, got %d keys", len(keys))
		}
	})

	t.Run("load invalid JSON", func(t *testing.T) {
		// Create invalid JSON file
		cacheFile := filepath.Join(tempDir, "version_cache.json")
		if err := os.WriteFile(cacheFile, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("Failed to create invalid JSON file: %v", err)
		}

		// Should not fail but should start with empty cache
		cache, err := NewFileCache(tempDir, 1*time.Hour)
		if err != nil {
			t.Errorf("Unexpected error creating cache with invalid JSON: %v", err)
		}

		keys := cache.Keys()
		if len(keys) != 0 {
			t.Errorf("Expected empty cache after loading invalid JSON, got %d keys", len(keys))
		}
	})
}

func TestCacheEdgeCases(t *testing.T) {
	t.Run("memory cache with zero TTL", func(t *testing.T) {
		cache := NewMemoryCache(0) // Should use default TTL

		err := cache.Set("test", "1.0.0")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		value, found := cache.Get("test")
		if !found || value != "1.0.0" {
			t.Error("Expected to find cached value with default TTL")
		}
	})

	t.Run("memory cache with negative TTL", func(t *testing.T) {
		cache := NewMemoryCache(-1 * time.Hour) // Should use default TTL

		err := cache.Set("test", "1.0.0")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		value, found := cache.Get("test")
		if !found || value != "1.0.0" {
			t.Error("Expected to find cached value with default TTL")
		}
	})

	t.Run("cache with empty key", func(t *testing.T) {
		cache := NewMemoryCache(1 * time.Hour)

		err := cache.Set("", "1.0.0")
		// Current implementation may allow empty keys - this is a test for future enhancement
		if err != nil {
			t.Logf("Empty key validation: %v", err)
		} else {
			t.Logf("Empty key allowed - consider adding validation")
		}

		_, found := cache.Get("")
		// Behavior may vary depending on implementation
		t.Logf("Empty key lookup result: found=%v", found)
	})

	t.Run("cache with nil value", func(t *testing.T) {
		cache := NewMemoryCache(1 * time.Hour)

		err := cache.Set("test", "")
		if err != nil {
			t.Errorf("Should allow empty string value: %v", err)
		}

		value, found := cache.Get("test")
		if !found || value != "" {
			t.Error("Should find empty string value")
		}
	})

	t.Run("file cache with read-only directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "readonly-cache-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Make directory read-only
		if err := os.Chmod(tempDir, 0444); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}
		defer os.Chmod(tempDir, 0755) // Restore for cleanup

		_, err = NewFileCache(tempDir, 1*time.Hour)
		// File cache may still work in read-only mode if it can't write cache files
		if err != nil {
			t.Logf("Read-only directory handling: %v", err)
		} else {
			t.Logf("File cache created in read-only directory - may work in read-only mode")
		}
	})

	t.Run("concurrent cache operations", func(t *testing.T) {
		cache := NewMemoryCache(1 * time.Hour)
		const numGoroutines = 100

		// Concurrent writes
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				key := fmt.Sprintf("key-%d", index)
				value := fmt.Sprintf("value-%d", index)
				cache.Set(key, value)
			}(i)
		}

		// Concurrent reads
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				key := fmt.Sprintf("key-%d", index)
				cache.Get(key)
			}(i)
		}

		// Give goroutines time to complete
		time.Sleep(100 * time.Millisecond)

		// Verify some data was written
		keys := cache.Keys()
		if len(keys) == 0 {
			t.Error("Expected some keys to be written by concurrent operations")
		}
	})
}

func TestCachePerformance(t *testing.T) {
	t.Run("memory cache performance", func(t *testing.T) {
		cache := NewMemoryCache(1 * time.Hour)
		const numOperations = 10000

		// Measure write performance
		start := time.Now()
		for i := 0; i < numOperations; i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			cache.Set(key, value)
		}
		writeTime := time.Since(start)

		// Measure read performance
		start = time.Now()
		for i := 0; i < numOperations; i++ {
			key := fmt.Sprintf("key-%d", i)
			cache.Get(key)
		}
		readTime := time.Since(start)

		t.Logf("Memory cache: %d writes in %v, %d reads in %v",
			numOperations, writeTime, numOperations, readTime)

		// Performance should be reasonable
		if writeTime > 100*time.Millisecond {
			t.Errorf("Memory cache writes too slow: %v for %d operations", writeTime, numOperations)
		}
		if readTime > 50*time.Millisecond {
			t.Errorf("Memory cache reads too slow: %v for %d operations", readTime, numOperations)
		}
	})

	t.Run("file cache performance", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "perf-cache-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		cache, err := NewFileCache(tempDir, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to create file cache: %v", err)
		}

		const numOperations = 1000 // Fewer operations for file cache

		// Measure write performance
		start := time.Now()
		for i := 0; i < numOperations; i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			cache.Set(key, value)
		}
		writeTime := time.Since(start)

		// Wait for async saves
		time.Sleep(100 * time.Millisecond)

		// Measure read performance
		start = time.Now()
		for i := 0; i < numOperations; i++ {
			key := fmt.Sprintf("key-%d", i)
			cache.Get(key)
		}
		readTime := time.Since(start)

		t.Logf("File cache: %d writes in %v, %d reads in %v",
			numOperations, writeTime, numOperations, readTime)

		// File cache will be slower than memory cache
		if writeTime > 5*time.Second {
			t.Errorf("File cache writes too slow: %v for %d operations", writeTime, numOperations)
		}
		if readTime > 1*time.Second {
			t.Errorf("File cache reads too slow: %v for %d operations", readTime, numOperations)
		}
	})
}

func TestCacheDefaultTTL(t *testing.T) {
	t.Run("memory cache default TTL", func(t *testing.T) {
		cache := NewMemoryCache(0) // Should use default TTL
		if cache.ttl != 24*time.Hour {
			t.Errorf("Expected default TTL of 24 hours, got %v", cache.ttl)
		}
	})

	t.Run("file cache default TTL", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "version-cache-ttl-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		cache, err := NewFileCache(tempDir, 0) // Should use default TTL
		if err != nil {
			t.Fatalf("Failed to create cache: %v", err)
		}

		if cache.ttl != 24*time.Hour {
			t.Errorf("Expected default TTL of 24 hours, got %v", cache.ttl)
		}
	})
}
