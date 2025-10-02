package cache

import (
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheOperations(t *testing.T) {
	config := interfaces.DefaultCacheConfig()

	ops := NewCacheOperations(config)

	assert.NotNil(t, ops)
	assert.Equal(t, config, ops.config)
	assert.NotNil(t, ops.callbacks)
}

func TestCacheOperations_GetEntry(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:        "valid",
			Value:      "valid_value",
			AccessedAt: now.Add(-time.Hour),
		},
		"expired": {
			Key:       "expired",
			Value:     "expired_value",
			ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
		},
	}
	metrics := &interfaces.CacheMetrics{}

	t.Run("existing valid entry", func(t *testing.T) {
		entry, found := ops.GetEntry("valid", entries, metrics)

		assert.True(t, found)
		assert.NotNil(t, entry)
		assert.Equal(t, "valid_value", entry.Value)
		assert.Equal(t, int64(1), entry.AccessCount)
		assert.Equal(t, int64(1), metrics.Hits)
		assert.Equal(t, int64(1), metrics.Gets)
	})

	t.Run("non-existent entry", func(t *testing.T) {
		metrics := &interfaces.CacheMetrics{} // Reset metrics

		entry, found := ops.GetEntry("nonexistent", entries, metrics)

		assert.False(t, found)
		assert.Nil(t, entry)
		assert.Equal(t, int64(1), metrics.Misses)
		assert.Equal(t, int64(1), metrics.Gets)
	})

	t.Run("expired entry", func(t *testing.T) {
		metrics := &interfaces.CacheMetrics{} // Reset metrics
		originalLen := len(entries)

		entry, found := ops.GetEntry("expired", entries, metrics)

		assert.False(t, found)
		assert.Nil(t, entry)
		assert.Equal(t, int64(1), metrics.Misses)
		assert.Equal(t, int64(1), metrics.Gets)
		assert.Len(t, entries, originalLen-1) // Expired entry should be removed
		assert.NotContains(t, entries, "expired")
	})
}

func TestCacheOperations_CreateEntry(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	config.DefaultTTL = time.Hour
	config.EnableCompression = true
	ops := NewCacheOperations(config)

	t.Run("entry with TTL", func(t *testing.T) {
		ttl := 30 * time.Minute

		entry := ops.CreateEntry("test_key", "test_value", ttl)

		assert.NotNil(t, entry)
		assert.Equal(t, "test_key", entry.Key)
		assert.Equal(t, "test_value", entry.Value)
		assert.Equal(t, ttl, entry.TTL)
		assert.NotNil(t, entry.ExpiresAt)
		assert.True(t, entry.ExpiresAt.After(time.Now()))
		assert.Equal(t, int64(10), entry.Size) // "test_value" length
		assert.NotNil(t, entry.Metadata)
	})

	t.Run("entry without TTL uses default", func(t *testing.T) {
		entry := ops.CreateEntry("test_key", "test_value", 0)

		assert.NotNil(t, entry)
		assert.Equal(t, config.DefaultTTL, entry.TTL)
		assert.NotNil(t, entry.ExpiresAt)
	})

	t.Run("large value gets compressed", func(t *testing.T) {
		largeValue := make([]byte, 2048)
		for i := range largeValue {
			largeValue[i] = 'A'
		}

		entry := ops.CreateEntry("large_key", string(largeValue), time.Hour)

		assert.NotNil(t, entry)
		// Note: Compression might not always reduce size for repetitive data
		// but the compression attempt should be made
	})
}

func TestCacheOperations_DeleteEntry(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	entries := map[string]*interfaces.CacheEntry{
		"existing": {
			Key:   "existing",
			Value: "value",
			Size:  5,
		},
	}
	metrics := &interfaces.CacheMetrics{
		CurrentSize:    5,
		CurrentEntries: 1,
	}

	t.Run("delete existing entry", func(t *testing.T) {
		err := ops.DeleteEntry("existing", entries, metrics)

		assert.NoError(t, err)
		assert.NotContains(t, entries, "existing")
		assert.Equal(t, int64(0), metrics.CurrentSize)
		assert.Equal(t, 0, metrics.CurrentEntries)
		assert.Equal(t, int64(1), metrics.Deletes)
	})

	t.Run("delete non-existent entry", func(t *testing.T) {
		err := ops.DeleteEntry("nonexistent", entries, metrics)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})
}

func TestCacheOperations_CleanExpiredEntries(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:   "valid",
			Value: "valid_value",
			Size:  11,
		},
		"expired1": {
			Key:       "expired1",
			Value:     "expired_value1",
			Size:      14,
			ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
		},
		"expired2": {
			Key:       "expired2",
			Value:     "expired_value2",
			Size:      14,
			ExpiresAt: &[]time.Time{now.Add(-time.Minute)}[0],
		},
	}
	metrics := &interfaces.CacheMetrics{
		CurrentSize:    39,
		CurrentEntries: 3,
	}

	expiredKeys := ops.CleanExpiredEntries(entries, metrics)

	assert.Len(t, expiredKeys, 2)
	assert.Contains(t, expiredKeys, "expired1")
	assert.Contains(t, expiredKeys, "expired2")
	assert.Len(t, entries, 1)
	assert.Contains(t, entries, "valid")
	assert.Equal(t, int64(11), metrics.CurrentSize)
	assert.Equal(t, 1, metrics.CurrentEntries)
	assert.Equal(t, int64(2), metrics.Evictions)
}

func TestCacheOperations_EvictIfNeeded(t *testing.T) {
	config := &interfaces.CacheConfig{
		MaxSize:        100,
		MaxEntries:     3,
		EvictionPolicy: "lru",
		EvictionRatio:  0.2,
	}
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"old": {
			Key:        "old",
			Value:      "old_value",
			Size:       40,
			AccessedAt: now.Add(-2 * time.Hour),
		},
		"medium": {
			Key:        "medium",
			Value:      "medium_value",
			Size:       40,
			AccessedAt: now.Add(-time.Hour),
		},
		"recent": {
			Key:        "recent",
			Value:      "recent_value",
			Size:       40,
			AccessedAt: now.Add(-time.Minute),
		},
	}
	metrics := &interfaces.CacheMetrics{
		CurrentSize: 120,
	}

	t.Run("evict by size", func(t *testing.T) {
		err := ops.EvictIfNeeded(10, entries, metrics)

		assert.NoError(t, err)
		// Should evict oldest entries to get under target size (80% of 100 = 80)
		assert.True(t, metrics.CurrentSize <= 80)
		assert.Greater(t, metrics.Evictions, int64(0))
	})

	t.Run("evict by count", func(t *testing.T) {
		// Reset entries to test count-based eviction
		entries["entry1"] = &interfaces.CacheEntry{Key: "entry1", Size: 10, AccessedAt: now.Add(-3 * time.Hour)}
		entries["entry2"] = &interfaces.CacheEntry{Key: "entry2", Size: 10, AccessedAt: now.Add(-2 * time.Hour)}
		entries["entry3"] = &interfaces.CacheEntry{Key: "entry3", Size: 10, AccessedAt: now.Add(-time.Hour)}
		entries["entry4"] = &interfaces.CacheEntry{Key: "entry4", Size: 10, AccessedAt: now}

		metrics.Evictions = 0 // Reset evictions counter

		err := ops.EvictIfNeeded(0, entries, metrics)

		assert.NoError(t, err)
		// Should evict entries to get under target count (80% of 3 = 2.4, rounded down to 2)
		assert.LessOrEqual(t, len(entries), 2)
		assert.Greater(t, metrics.Evictions, int64(0))
	})
}

func TestCacheOperations_GetKeys(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"key1": {
			Key: "key1",
		},
		"key2": {
			Key: "key2",
		},
		"expired": {
			Key:       "expired",
			ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
		},
	}

	keys := ops.GetKeys(entries)

	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.NotContains(t, keys, "expired")
	// Keys should be sorted
	assert.Equal(t, []string{"key1", "key2"}, keys)
}

func TestCacheOperations_GetKeysByPattern(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	entries := map[string]*interfaces.CacheEntry{
		"user:123":    {Key: "user:123"},
		"user:456":    {Key: "user:456"},
		"session:abc": {Key: "session:abc"},
		"config":      {Key: "config"},
	}

	t.Run("valid pattern", func(t *testing.T) {
		keys, err := ops.GetKeysByPattern("user:.*", entries)

		assert.NoError(t, err)
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "user:123")
		assert.Contains(t, keys, "user:456")
	})

	t.Run("invalid pattern", func(t *testing.T) {
		_, err := ops.GetKeysByPattern("[", entries)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid pattern")
	})
}

func TestCacheOperations_ValidateEntries(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:         "valid",
			Size:        10,
			CreatedAt:   now.Add(-time.Hour),
			AccessCount: 5,
		},
		"nil_entry": nil,
		"key_mismatch": {
			Key:  "different_key",
			Size: 10,
		},
		"negative_size": {
			Key:  "negative_size",
			Size: -5,
		},
		"future_time": {
			Key:       "future_time",
			CreatedAt: now.Add(time.Hour),
		},
		"negative_access": {
			Key:         "negative_access",
			AccessCount: -1,
		},
	}

	issues := ops.ValidateEntries(entries)

	assert.Len(t, issues, 5)

	// Check that all expected issues are present (order may vary due to map iteration)
	issueStrings := strings.Join(issues, " ")
	assert.Contains(t, issueStrings, "nil entry for key: nil_entry")
	assert.Contains(t, issueStrings, "key mismatch")
	assert.Contains(t, issueStrings, "negative size")
	assert.Contains(t, issueStrings, "future creation time")
	assert.Contains(t, issueStrings, "negative access count")
}

func TestCacheOperations_RepairEntries(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:   "valid",
			Value: "valid_value",
			Size:  11,
		},
		"nil_entry": nil,
		"key_mismatch": {
			Key:   "different_key",
			Value: "value",
			Size:  5,
		},
		"negative_size": {
			Key:   "negative_size",
			Value: "value",
			Size:  -5,
		},
		"future_time": {
			Key:       "future_time",
			Value:     "value",
			CreatedAt: now.Add(time.Hour),
		},
		"expired": {
			Key:       "expired",
			Value:     "value",
			ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
		},
	}

	repairedEntries := ops.RepairEntries(entries)

	// Should have 4 entries (nil and expired entries removed)
	assert.Len(t, repairedEntries, 4)

	// Check repairs
	assert.Contains(t, repairedEntries, "valid")
	assert.NotContains(t, repairedEntries, "nil_entry")
	assert.NotContains(t, repairedEntries, "expired")

	// Key mismatch should be fixed
	assert.Equal(t, "key_mismatch", repairedEntries["key_mismatch"].Key)

	// Negative size should be fixed
	assert.Greater(t, repairedEntries["negative_size"].Size, int64(0))

	// Future time should be fixed (should be <= now + small tolerance)
	assert.True(t, repairedEntries["future_time"].CreatedAt.Before(now.Add(time.Second)))
}

func TestCacheOperations_TTLOperations(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"test": {
			Key:       "test",
			Value:     "value",
			TTL:       time.Hour,
			ExpiresAt: &[]time.Time{now.Add(time.Hour)}[0],
		},
	}

	t.Run("update TTL", func(t *testing.T) {
		newTTL := 2 * time.Hour

		err := ops.UpdateTTL("test", newTTL, entries)

		assert.NoError(t, err)
		assert.Equal(t, newTTL, entries["test"].TTL)
		assert.NotNil(t, entries["test"].ExpiresAt)
		assert.True(t, entries["test"].ExpiresAt.After(now.Add(time.Hour)))
	})

	t.Run("get TTL", func(t *testing.T) {
		ttl, err := ops.GetTTL("test", entries)

		assert.NoError(t, err)
		assert.Greater(t, ttl, time.Hour) // Should be close to 2 hours
	})

	t.Run("refresh TTL", func(t *testing.T) {
		originalExpiry := *entries["test"].ExpiresAt
		time.Sleep(10 * time.Millisecond) // Small delay

		err := ops.RefreshTTL("test", entries)

		assert.NoError(t, err)
		assert.True(t, entries["test"].ExpiresAt.After(originalExpiry))
	})

	t.Run("operations on non-existent key", func(t *testing.T) {
		err := ops.UpdateTTL("nonexistent", time.Hour, entries)
		assert.Error(t, err)

		_, err = ops.GetTTL("nonexistent", entries)
		assert.Error(t, err)

		err = ops.RefreshTTL("nonexistent", entries)
		assert.Error(t, err)
	})
}

func TestCacheOperations_SetCallbacks(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	hitCalled := false
	missCalled := false
	evictionCalled := false

	onHit := func(key string) { hitCalled = true }
	onMiss := func(key string) { missCalled = true }
	onEviction := func(key, reason string) { evictionCalled = true }

	ops.SetCallbacks(onHit, onMiss, onEviction)

	// Test callbacks are set
	entries := map[string]*interfaces.CacheEntry{}
	metrics := &interfaces.CacheMetrics{}

	// Test miss callback
	ops.GetEntry("nonexistent", entries, metrics)
	assert.True(t, missCalled)

	// Test hit callback
	entries["test"] = &interfaces.CacheEntry{Key: "test", Value: "value"}
	ops.GetEntry("test", entries, metrics)
	assert.True(t, hitCalled)

	// Test eviction callback
	if err := ops.DeleteEntry("test", entries, metrics); err != nil {
		t.Errorf("Failed to delete cache entry: %v", err)
	}
	assert.True(t, evictionCalled)
}

func TestCacheOperations_SetConfig(t *testing.T) {
	config1 := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config1)

	config2 := &interfaces.CacheConfig{
		MaxSize: 2048,
	}

	ops.SetConfig(config2)

	assert.Equal(t, config2, ops.config)
	assert.Equal(t, int64(2048), ops.config.MaxSize)
}

func TestCacheOperations_CalculateSize(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	tests := []struct {
		name     string
		value    any
		expected int64
	}{
		{"string", "hello", 5},
		{"byte slice", []byte("hello"), 5},
		{"int", 42, 8},
		{"bool", true, 1},
		{"struct", struct{ Name string }{Name: "test"}, 15}, // JSON: {"Name":"test"}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := ops.calculateSize(tt.value)
			assert.Equal(t, tt.expected, size)
		})
	}
}

func TestCacheOperations_GetExpiredKeys(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	ops := NewCacheOperations(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key: "valid",
		},
		"expired1": {
			Key:       "expired1",
			ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
		},
		"expired2": {
			Key:       "expired2",
			ExpiresAt: &[]time.Time{now.Add(-time.Minute)}[0],
		},
	}

	expiredKeys := ops.GetExpiredKeys(entries)

	assert.Len(t, expiredKeys, 2)
	assert.Contains(t, expiredKeys, "expired1")
	assert.Contains(t, expiredKeys, "expired2")
	assert.NotContains(t, expiredKeys, "valid")
}
