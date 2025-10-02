package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCacheValidator(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator("/tmp/test-cache", config)

	assert.NotNil(t, validator)
	assert.Equal(t, "/tmp/test-cache", validator.cacheDir)
	assert.Equal(t, config, validator.config)
}

func TestCacheValidator_ValidateCache(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cache-validator-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator(tempDir, config)

	t.Run("ValidCache", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:         "key1",
				Value:       "value1",
				Size:        6,
				CreatedAt:   time.Now().Add(-1 * time.Hour),
				UpdatedAt:   time.Now().Add(-1 * time.Hour),
				AccessedAt:  time.Now().Add(-30 * time.Minute),
				AccessCount: 5,
				Metadata:    make(map[string]any),
			},
		}

		err := validator.ValidateCache(entries)
		assert.NoError(t, err)
	})

	t.Run("InvalidEntries", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:         "wrong-key", // Key mismatch
				Value:       "value1",
				Size:        -1,                            // Negative size
				CreatedAt:   time.Now().Add(1 * time.Hour), // Future time
				UpdatedAt:   time.Now(),
				AccessedAt:  time.Now(),
				AccessCount: -1, // Negative access count
				Metadata:    make(map[string]any),
			},
		}

		err := validator.ValidateCache(entries)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key mismatch")
		assert.Contains(t, err.Error(), "negative size")
		assert.Contains(t, err.Error(), "future creation time")
		assert.Contains(t, err.Error(), "negative access count")
	})

	t.Run("NilEntry", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"key1": nil,
		}

		err := validator.ValidateCache(entries)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil entry")
	})
}

func TestCacheValidator_ValidateEntryMetadata(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator("/tmp/test", config)

	t.Run("ValidCompressedEntry", func(t *testing.T) {
		entry := &interfaces.CacheEntry{
			Key:        "test",
			Compressed: true,
			Metadata: map[string]any{
				"compression_type": "gzip",
				"original_size":    100,
			},
		}

		err := validator.validateEntryMetadata("test", entry)
		assert.NoError(t, err)
	})

	t.Run("CompressedEntryMissingMetadata", func(t *testing.T) {
		entry := &interfaces.CacheEntry{
			Key:        "test",
			Compressed: true,
			Metadata:   map[string]any{},
		}

		err := validator.validateEntryMetadata("test", entry)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing compression_type metadata")
	})

	t.Run("NegativeOriginalSize", func(t *testing.T) {
		entry := &interfaces.CacheEntry{
			Key:        "test",
			Compressed: true,
			Metadata: map[string]any{
				"compression_type": "gzip",
				"original_size":    -1,
			},
		}

		err := validator.validateEntryMetadata("test", entry)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "negative original_size")
	})
}

func TestCacheValidator_RepairCache(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cache-repair-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator(tempDir, config)

	t.Run("RepairCorruptedEntries", func(t *testing.T) {
		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:         "wrong-key", // Will be fixed
				Value:       "value1",
				Size:        -1,                     // Will be fixed
				CreatedAt:   now.Add(1 * time.Hour), // Will be fixed
				UpdatedAt:   now,
				AccessedAt:  now,
				AccessCount: -5,  // Will be fixed
				Metadata:    nil, // Will be initialized
			},
			"key2": nil, // Will be removed
			"key3": {
				Key:       "key3",
				Value:     "value3",
				Size:      6,
				CreatedAt: now.Add(-2 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-1 * time.Hour)}[0], // Expired, will be removed
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    100,
			CurrentEntries: 3,
		}

		repairedEntries, repairedMetrics, err := validator.RepairCache(entries, metrics)
		require.NoError(t, err)

		// Check that corrupted entry was fixed
		assert.Contains(t, repairedEntries, "key1")
		entry1 := repairedEntries["key1"]
		assert.Equal(t, "key1", entry1.Key)
		assert.Equal(t, int64(6), entry1.Size) // "value1" = 6 bytes
		// The CreatedAt should have been fixed to not be in the future
		assert.True(t, entry1.CreatedAt.Before(now.Add(1*time.Second)) || entry1.CreatedAt.Equal(now))
		assert.Equal(t, int64(0), entry1.AccessCount)
		assert.NotNil(t, entry1.Metadata)

		// Check that nil entry was removed
		assert.NotContains(t, repairedEntries, "key2")

		// Check that expired entry was removed
		assert.NotContains(t, repairedEntries, "key3")

		// Check metrics were recalculated
		assert.Equal(t, 1, repairedMetrics.CurrentEntries)
		assert.Equal(t, int64(6), repairedMetrics.CurrentSize)
	})
}

func TestCacheValidator_ValidateConfiguration(t *testing.T) {
	validator := NewCacheValidator("/tmp/test", nil)

	t.Run("NilConfig", func(t *testing.T) {
		err := validator.ValidateConfiguration()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "configuration is nil")
	})

	t.Run("ValidConfig", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		validator.SetConfig(config)

		err := validator.ValidateConfiguration()
		assert.NoError(t, err)
	})

	t.Run("InvalidMaxSize", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		config.MaxSize = -1
		validator.SetConfig(config)

		err := validator.ValidateConfiguration()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MaxSize cannot be negative")
	})

	t.Run("InvalidEvictionRatio", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		config.EvictionRatio = 1.5
		validator.SetConfig(config)

		err := validator.ValidateConfiguration()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EvictionRatio must be between 0 and 1")
	})

	t.Run("InvalidEvictionPolicy", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		config.EvictionPolicy = "invalid"
		validator.SetConfig(config)

		err := validator.ValidateConfiguration()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid eviction policy")
	})

	t.Run("InvalidCompressionLevel", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		config.EnableCompression = true
		config.CompressionLevel = 10
		validator.SetConfig(config)

		err := validator.ValidateConfiguration()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CompressionLevel must be between 1 and 9")
	})
}

func TestCacheValidator_CheckCacheHealth(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cache-health-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator(tempDir, config)

	t.Run("HealthyCache", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:         "key1",
				Value:       "value1",
				Size:        6,
				CreatedAt:   time.Now().Add(-1 * time.Hour),
				UpdatedAt:   time.Now().Add(-1 * time.Hour),
				AccessedAt:  time.Now().Add(-30 * time.Minute),
				AccessCount: 5,
				Metadata:    make(map[string]any),
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    6,
			MaxSize:        1000,
			CurrentEntries: 1,
			MaxEntries:     100,
			Hits:           10,
			Gets:           15,
		}

		report, err := validator.CheckCacheHealth(entries, metrics)
		require.NoError(t, err)

		assert.Equal(t, "healthy", report.OverallHealth)
		assert.Equal(t, 0, report.ExpiredEntries)
		assert.Equal(t, 0, report.CorruptedEntries)
		assert.Equal(t, 1, report.TotalEntries)
		assert.Empty(t, report.Issues)
	})

	t.Run("DegradedCache", func(t *testing.T) {
		// Create separate temp dir for degraded cache test
		tempDir2, err := os.MkdirTemp("", "cache-degraded-test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir2)

		validator2 := NewCacheValidator(tempDir2, config)

		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"expired1": {
				Key:       "expired1",
				Value:     "value1",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-30 * time.Minute)}[0], // Expired
			},
			"expired2": {
				Key:       "expired2",
				Value:     "value2",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-20 * time.Minute)}[0], // Expired
			},
			"expired3": {
				Key:       "expired3",
				Value:     "value3",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-10 * time.Minute)}[0], // Expired
			},
			"valid": {
				Key:         "valid",
				Value:       "value4",
				Size:        6,
				CreatedAt:   now.Add(-1 * time.Hour),
				AccessCount: 5,
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    24,
			MaxSize:        1000,
			CurrentEntries: 4,
			MaxEntries:     100,
			Hits:           5,
			Gets:           15, // Low hit rate
		}

		report, err := validator2.CheckCacheHealth(entries, metrics)
		require.NoError(t, err)

		assert.Equal(t, "degraded", report.OverallHealth)
		assert.Equal(t, 3, report.ExpiredEntries)
		assert.Equal(t, 0, report.CorruptedEntries)
		assert.Equal(t, 4, report.TotalEntries)
		assert.NotEmpty(t, report.Issues)
		assert.Contains(t, report.Issues[0], "High expired entry ratio")
	})

	t.Run("UnhealthyCache", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"key1": {
				Key:         "wrong-key", // Corrupted
				Value:       "value1",
				Size:        -1, // Corrupted
				CreatedAt:   time.Now(),
				AccessCount: -1, // Corrupted
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    1500, // Exceeds max
			MaxSize:        1000,
			CurrentEntries: 1,
			MaxEntries:     100,
		}

		report, err := validator.CheckCacheHealth(entries, metrics)
		require.NoError(t, err)

		assert.Equal(t, "unhealthy", report.OverallHealth)
		assert.Equal(t, 1, report.CorruptedEntries)
		assert.NotEmpty(t, report.Issues)
		assert.Contains(t, report.Recommendations, "Consider running cache repair")
	})
}

func TestCacheValidator_ValidateCacheDirectory(t *testing.T) {
	t.Run("ValidDirectory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "cache-dir-test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		config := interfaces.DefaultCacheConfig()
		validator := NewCacheValidator(tempDir, config)

		err = validator.validateCacheDirectory()
		assert.NoError(t, err)
	})

	t.Run("NonExistentDirectory", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		validator := NewCacheValidator("/non/existent/path", config)

		err := validator.validateCacheDirectory()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("FileInsteadOfDirectory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "cache-file-test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a file instead of directory
		filePath := filepath.Join(tempDir, "cache-file")
		err = os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		config := interfaces.DefaultCacheConfig()
		validator := NewCacheValidator(filePath, config)

		err = validator.validateCacheDirectory()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})
}

func TestCacheValidator_CalculateSize(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator("/tmp/test", config)

	tests := []struct {
		name     string
		value    any
		expected int64
	}{
		{"String", "hello", 5},
		{"ByteSlice", []byte("hello"), 5},
		{"Int", 42, 8},
		{"Bool", true, 1},
		{"Struct", struct{ Name string }{Name: "test"}, 15}, // JSON: {"Name":"test"}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := validator.calculateSize(tt.value)
			assert.Equal(t, tt.expected, size)
		})
	}
}

func TestCacheValidator_SettersAndGetters(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	validator := NewCacheValidator("/tmp/test", config)

	// Test GetCacheDir
	assert.Equal(t, "/tmp/test", validator.GetCacheDir())

	// Test SetConfig
	newConfig := &interfaces.CacheConfig{
		MaxSize: 2000,
	}
	validator.SetConfig(newConfig)
	assert.Equal(t, newConfig, validator.config)
}
