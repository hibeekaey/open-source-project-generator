package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCacheStorage(t *testing.T) {
	cacheDir := "/tmp/test-cache"
	config := interfaces.DefaultCacheConfig()

	storage := NewCacheStorage(cacheDir, config)

	assert.NotNil(t, storage)
	assert.Equal(t, cacheDir, storage.cacheDir)
	assert.Equal(t, config, storage.config)
}

func TestCacheStorage_Initialize(t *testing.T) {
	tests := []struct {
		name      string
		cacheDir  string
		wantError bool
	}{
		{
			name:      "valid directory",
			cacheDir:  filepath.Join(os.TempDir(), "test-cache-init"),
			wantError: false,
		},
		{
			name:      "invalid directory",
			cacheDir:  "/invalid/path/that/cannot/be/created",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := interfaces.DefaultCacheConfig()
			storage := NewCacheStorage(tt.cacheDir, config)

			err := storage.Initialize()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify directory was created
				_, err := os.Stat(tt.cacheDir)
				assert.NoError(t, err)

				// Cleanup
				if err := os.RemoveAll(tt.cacheDir); err != nil {
					t.Errorf("Failed to remove cache directory: %v", err)
				}
			}
		})
	}
}

func TestCacheStorage_LoadCache(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-load")
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	// Initialize storage
	require.NoError(t, storage.Initialize())
	defer os.RemoveAll(cacheDir)

	t.Run("no existing cache file", func(t *testing.T) {
		entries, metrics, err := storage.LoadCache()

		assert.NoError(t, err)
		assert.NotNil(t, entries)
		assert.NotNil(t, metrics)
		assert.Empty(t, entries)
	})

	t.Run("existing cache file", func(t *testing.T) {
		// Create test cache data
		now := time.Now()
		testEntries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:        "key1",
				Value:      "value1",
				Size:       6,
				CreatedAt:  now,
				UpdatedAt:  now,
				AccessedAt: now,
			},
		}
		testMetrics := &interfaces.CacheMetrics{
			CurrentEntries: 1,
			CurrentSize:    6,
		}

		// Save test data
		err := storage.SaveCache(testEntries, testMetrics)
		require.NoError(t, err)

		// Load cache
		entries, metrics, err := storage.LoadCache()

		assert.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Contains(t, entries, "key1")
		assert.Equal(t, "value1", entries["key1"].Value)
		assert.Equal(t, int64(6), metrics.CurrentSize)
	})

	t.Run("expired entries filtered out", func(t *testing.T) {
		// Create test cache data with expired entry
		now := time.Now()
		expiredTime := now.Add(-time.Hour)
		testEntries := map[string]*interfaces.CacheEntry{
			"expired": {
				Key:        "expired",
				Value:      "expired_value",
				Size:       13,
				CreatedAt:  now.Add(-2 * time.Hour),
				UpdatedAt:  now.Add(-2 * time.Hour),
				AccessedAt: now.Add(-2 * time.Hour),
				ExpiresAt:  &expiredTime,
			},
			"valid": {
				Key:        "valid",
				Value:      "valid_value",
				Size:       11,
				CreatedAt:  now,
				UpdatedAt:  now,
				AccessedAt: now,
			},
		}
		testMetrics := &interfaces.CacheMetrics{
			CurrentEntries: 2,
			CurrentSize:    24,
		}

		// Save test data
		err := storage.SaveCache(testEntries, testMetrics)
		require.NoError(t, err)

		// Load cache
		entries, _, err := storage.LoadCache()

		assert.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Contains(t, entries, "valid")
		assert.NotContains(t, entries, "expired")
	})
}

func TestCacheStorage_SaveCache(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-save")
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	// Initialize storage
	require.NoError(t, storage.Initialize())
	defer os.RemoveAll(cacheDir)

	t.Run("persistence enabled", func(t *testing.T) {
		config.PersistToDisk = true
		storage.SetConfig(config)

		now := time.Now()
		testEntries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:        "key1",
				Value:      "value1",
				Size:       6,
				CreatedAt:  now,
				UpdatedAt:  now,
				AccessedAt: now,
			},
		}
		testMetrics := &interfaces.CacheMetrics{
			CurrentEntries: 1,
			CurrentSize:    6,
		}

		err := storage.SaveCache(testEntries, testMetrics)

		assert.NoError(t, err)

		// Verify file was created
		cacheFile := filepath.Join(cacheDir, "cache.json")
		_, err = os.Stat(cacheFile)
		assert.NoError(t, err)

		// Verify file content
		data, err := os.ReadFile(cacheFile)
		require.NoError(t, err)

		var cacheData CacheFile
		err = json.Unmarshal(data, &cacheData)
		require.NoError(t, err)

		assert.Equal(t, "1.0", cacheData.Version)
		assert.Len(t, cacheData.Entries, 1)
		assert.Contains(t, cacheData.Entries, "key1")
	})

	t.Run("persistence disabled", func(t *testing.T) {
		config.PersistToDisk = false
		storage.SetConfig(config)

		testEntries := map[string]*interfaces.CacheEntry{}
		testMetrics := &interfaces.CacheMetrics{}

		err := storage.SaveCache(testEntries, testMetrics)

		assert.NoError(t, err)

		// Verify no file was created (or existing file wasn't modified)
		cacheFile := filepath.Join(cacheDir, "cache.json")
		info1, _ := os.Stat(cacheFile)

		time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference

		err = storage.SaveCache(testEntries, testMetrics)
		assert.NoError(t, err)

		info2, _ := os.Stat(cacheFile)
		if info1 != nil && info2 != nil {
			assert.Equal(t, info1.ModTime(), info2.ModTime())
		}
	})
}

func TestCacheStorage_ClearCache(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-clear")
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	// Initialize storage
	require.NoError(t, storage.Initialize())
	defer os.RemoveAll(cacheDir)

	// Create a cache file
	testEntries := map[string]*interfaces.CacheEntry{
		"key1": {Key: "key1", Value: "value1"},
	}
	testMetrics := &interfaces.CacheMetrics{}

	err := storage.SaveCache(testEntries, testMetrics)
	require.NoError(t, err)

	// Verify file exists
	cacheFile := filepath.Join(cacheDir, "cache.json")
	_, err = os.Stat(cacheFile)
	require.NoError(t, err)

	// Clear cache
	err = storage.ClearCache()

	assert.NoError(t, err)

	// Verify file was removed
	_, err = os.Stat(cacheFile)
	assert.True(t, os.IsNotExist(err))
}

func TestCacheStorage_BackupCache(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-backup")
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	// Initialize storage
	require.NoError(t, storage.Initialize())
	defer os.RemoveAll(cacheDir)

	now := time.Now()
	testEntries := map[string]*interfaces.CacheEntry{
		"key1": {
			Key:        "key1",
			Value:      "value1",
			Size:       6,
			CreatedAt:  now,
			UpdatedAt:  now,
			AccessedAt: now,
		},
	}
	testMetrics := &interfaces.CacheMetrics{
		CurrentEntries: 1,
		CurrentSize:    6,
	}

	t.Run("valid backup path", func(t *testing.T) {
		backupPath := filepath.Join(os.TempDir(), "cache-backup.json")
		defer os.Remove(backupPath)

		err := storage.BackupCache(backupPath, testEntries, testMetrics)

		assert.NoError(t, err)

		// Verify backup file was created
		_, err = os.Stat(backupPath)
		assert.NoError(t, err)

		// Verify backup content
		data, err := os.ReadFile(backupPath)
		require.NoError(t, err)

		var backup CacheFile
		err = json.Unmarshal(data, &backup)
		require.NoError(t, err)

		assert.Equal(t, "1.0", backup.Version)
		assert.Len(t, backup.Entries, 1)
		assert.Contains(t, backup.Entries, "key1")
	})

	t.Run("invalid backup path", func(t *testing.T) {
		backupPath := "../invalid/path/backup.json"

		err := storage.BackupCache(backupPath, testEntries, testMetrics)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path traversal detected")
	})
}

func TestCacheStorage_RestoreCache(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-restore")
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	// Initialize storage
	require.NoError(t, storage.Initialize())
	defer os.RemoveAll(cacheDir)

	t.Run("valid restore", func(t *testing.T) {
		// Create backup file
		backupPath := filepath.Join(os.TempDir(), "restore-backup.json")
		defer os.Remove(backupPath)

		now := time.Now()
		backup := CacheFile{
			Version:   "1.0",
			CreatedAt: now,
			UpdatedAt: now,
			Entries: map[string]*interfaces.CacheEntry{
				"restored": {
					Key:        "restored",
					Value:      "restored_value",
					Size:       14,
					CreatedAt:  now,
					UpdatedAt:  now,
					AccessedAt: now,
				},
			},
			Metrics: &interfaces.CacheMetrics{
				CurrentEntries: 1,
				CurrentSize:    14,
			},
		}

		data, err := json.Marshal(backup)
		require.NoError(t, err)

		err = os.WriteFile(backupPath, data, 0600)
		require.NoError(t, err)

		// Restore cache
		entries, metrics, err := storage.RestoreCache(backupPath)

		assert.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Contains(t, entries, "restored")
		assert.Equal(t, "restored_value", entries["restored"].Value)
		assert.Equal(t, int64(14), metrics.CurrentSize)
	})

	t.Run("invalid restore path", func(t *testing.T) {
		restorePath := "../invalid/path/backup.json"

		_, _, err := storage.RestoreCache(restorePath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path traversal detected")
	})

	t.Run("non-existent file", func(t *testing.T) {
		restorePath := filepath.Join(os.TempDir(), "non-existent.json")

		_, _, err := storage.RestoreCache(restorePath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open backup file")
	})
}

func TestCacheStorage_ValidateStorage(t *testing.T) {
	t.Run("valid storage", func(t *testing.T) {
		cacheDir := filepath.Join(os.TempDir(), "test-cache-validate")
		config := interfaces.DefaultCacheConfig()
		storage := NewCacheStorage(cacheDir, config)

		// Initialize storage
		require.NoError(t, storage.Initialize())
		defer os.RemoveAll(cacheDir)

		err := storage.ValidateStorage()

		assert.NoError(t, err)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		cacheDir := "/non/existent/directory"
		config := interfaces.DefaultCacheConfig()
		storage := NewCacheStorage(cacheDir, config)

		err := storage.ValidateStorage()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache directory does not exist")
	})
}

func TestCacheStorage_GetCacheDir(t *testing.T) {
	cacheDir := "/test/cache/dir"
	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config)

	result := storage.GetCacheDir()

	assert.Equal(t, cacheDir, result)
}

func TestCacheStorage_SetConfig(t *testing.T) {
	cacheDir := "/test/cache/dir"
	config1 := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage(cacheDir, config1)

	config2 := &interfaces.CacheConfig{
		MaxSize: 2048,
	}

	storage.SetConfig(config2)

	assert.Equal(t, config2, storage.config)
	assert.Equal(t, int64(2048), storage.config.MaxSize)
}

func TestCacheStorage_ValidateStoragePath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			name:      "valid relative path",
			path:      "cache/backup.json",
			wantError: false,
		},
		{
			name:      "valid absolute path",
			path:      "/tmp/cache/backup.json",
			wantError: false,
		},
		{
			name:      "path traversal attempt",
			path:      "../../../etc/passwd",
			wantError: true,
		},
		{
			name:      "system directory access",
			path:      "/etc/passwd",
			wantError: true,
		},
		{
			name:      "proc directory access",
			path:      "/proc/version",
			wantError: true,
		},
		{
			name:      "sys directory access",
			path:      "/sys/kernel/version",
			wantError: true,
		},
	}

	config := interfaces.DefaultCacheConfig()
	storage := NewCacheStorage("/tmp", config)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.validateStoragePath(tt.path)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
