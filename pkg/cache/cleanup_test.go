package cache

import (
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCacheCleanup(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	assert.NotNil(t, cleanup)
	assert.Equal(t, config, cleanup.config)
}

func TestCacheCleanup_CleanExpiredEntries(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	t.Run("CleanExpiredEntries", func(t *testing.T) {
		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"valid": {
				Key:       "valid",
				Value:     "value1",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(1 * time.Hour)}[0], // Valid
			},
			"expired": {
				Key:       "expired",
				Value:     "value2",
				Size:      6,
				CreatedAt: now.Add(-2 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-1 * time.Hour)}[0], // Expired
			},
			"no_expiry": {
				Key:       "no_expiry",
				Value:     "value3",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				// No ExpiresAt - should not be cleaned
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    18,
			CurrentEntries: 3,
		}

		cleanedCount, err := cleanup.CleanExpiredEntries(entries, metrics)
		require.NoError(t, err)

		assert.Equal(t, 1, cleanedCount)
		assert.Equal(t, 2, len(entries))
		assert.Contains(t, entries, "valid")
		assert.Contains(t, entries, "no_expiry")
		assert.NotContains(t, entries, "expired")
		assert.Equal(t, int64(12), metrics.CurrentSize)
		assert.Equal(t, 2, metrics.CurrentEntries)
		assert.Equal(t, int64(1), metrics.Evictions)
	})

	t.Run("NoExpiredEntries", func(t *testing.T) {
		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"valid": {
				Key:       "valid",
				Value:     "value1",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(1 * time.Hour)}[0],
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    6,
			CurrentEntries: 1,
		}

		cleanedCount, err := cleanup.CleanExpiredEntries(entries, metrics)
		require.NoError(t, err)

		assert.Equal(t, 0, cleanedCount)
		assert.Equal(t, 1, len(entries))
		assert.Equal(t, int64(6), metrics.CurrentSize)
	})
}

func TestCacheCleanup_CompactCache(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:       "valid",
			Value:     "value1",
			Size:      6,
			CreatedAt: now.Add(-1 * time.Hour),
		},
		"expired": {
			Key:       "expired",
			Value:     "value2",
			Size:      6,
			CreatedAt: now.Add(-2 * time.Hour),
			ExpiresAt: &[]time.Time{now.Add(-1 * time.Hour)}[0], // Expired
		},
	}

	metrics := &interfaces.CacheMetrics{
		CurrentSize:    12,
		CurrentEntries: 2,
	}

	result, err := cleanup.CompactCache(entries, metrics)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, 2, result.InitialEntries)
	assert.Equal(t, 1, result.FinalEntries)
	assert.Equal(t, 1, result.EntriesRemoved)
	assert.Equal(t, 1, result.ExpiredEntriesRemoved)
	assert.Equal(t, int64(12), result.InitialSize)
	assert.Equal(t, int64(6), result.FinalSize)
	assert.Equal(t, int64(6), result.SizeReduced)
	assert.True(t, result.Duration > 0)

	// Check that expired entry was removed
	assert.Equal(t, 1, len(entries))
	assert.Contains(t, entries, "valid")
	assert.NotContains(t, entries, "expired")
}

func TestCacheCleanup_PerformMaintenance(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"valid": {
			Key:         "valid",
			Value:       "value1",
			Size:        6,
			CreatedAt:   now.Add(-1 * time.Hour),
			UpdatedAt:   now.Add(-1 * time.Hour),
			AccessedAt:  now.Add(-30 * time.Minute),
			AccessCount: 5,
			Metadata:    make(map[string]any),
		},
		"expired": {
			Key:       "expired",
			Value:     "value2",
			Size:      6,
			CreatedAt: now.Add(-2 * time.Hour),
			ExpiresAt: &[]time.Time{now.Add(-1 * time.Hour)}[0], // Expired
		},
	}

	metrics := &interfaces.CacheMetrics{
		CurrentSize:    12,
		CurrentEntries: 2,
	}

	result, err := cleanup.PerformMaintenance(entries, metrics)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, 3, len(result.Tasks)) // Clean, Validate, Optimize
	assert.True(t, result.Duration > 0)

	// Check individual tasks
	cleanTask := result.Tasks[0]
	assert.Equal(t, "Clean Expired Entries", cleanTask.Name)
	assert.True(t, cleanTask.Success)
	assert.Contains(t, cleanTask.Details, "Removed 1 expired entries")

	validateTask := result.Tasks[1]
	assert.Equal(t, "Validate Cache Integrity", validateTask.Name)
	assert.True(t, validateTask.Success)

	optimizeTask := result.Tasks[2]
	assert.Equal(t, "Optimize Cache Structure", optimizeTask.Name)
	assert.True(t, optimizeTask.Success)
}

func TestCacheCleanup_CleanupByAge(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	t.Run("CleanupOldEntries", func(t *testing.T) {
		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"recent": {
				Key:       "recent",
				Value:     "value1",
				Size:      6,
				CreatedAt: now.Add(-1 * time.Hour), // Recent
			},
			"old": {
				Key:       "old",
				Value:     "value2",
				Size:      6,
				CreatedAt: now.Add(-25 * time.Hour), // Old
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    12,
			CurrentEntries: 2,
		}

		maxAge := 24 * time.Hour
		cleanedCount, err := cleanup.CleanupByAge(entries, metrics, maxAge)
		require.NoError(t, err)

		assert.Equal(t, 1, cleanedCount)
		assert.Equal(t, 1, len(entries))
		assert.Contains(t, entries, "recent")
		assert.NotContains(t, entries, "old")
		assert.Equal(t, int64(6), metrics.CurrentSize)
	})

	t.Run("InvalidMaxAge", func(t *testing.T) {
		entries := make(map[string]*interfaces.CacheEntry)
		metrics := &interfaces.CacheMetrics{}

		_, err := cleanup.CleanupByAge(entries, metrics, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maxAge must be positive")
	})
}

func TestCacheCleanup_CleanupBySize(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	t.Run("CleanupToTargetSize", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"entry1": {
				Key:   "entry1",
				Value: "value1",
				Size:  10,
			},
			"entry2": {
				Key:   "entry2",
				Value: "value2",
				Size:  10,
			},
			"entry3": {
				Key:   "entry3",
				Value: "value3",
				Size:  10,
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    30,
			CurrentEntries: 3,
		}

		targetSize := int64(15)
		cleanedCount, err := cleanup.CleanupBySize(entries, metrics, targetSize)
		require.NoError(t, err)

		assert.True(t, cleanedCount > 0)
		assert.True(t, metrics.CurrentSize <= targetSize)
	})

	t.Run("AlreadyWithinTargetSize", func(t *testing.T) {
		entries := map[string]*interfaces.CacheEntry{
			"entry1": {
				Key:   "entry1",
				Value: "value1",
				Size:  5,
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    5,
			CurrentEntries: 1,
		}

		targetSize := int64(10)
		cleanedCount, err := cleanup.CleanupBySize(entries, metrics, targetSize)
		require.NoError(t, err)

		assert.Equal(t, 0, cleanedCount)
		assert.Equal(t, int64(5), metrics.CurrentSize)
	})

	t.Run("NegativeTargetSize", func(t *testing.T) {
		entries := make(map[string]*interfaces.CacheEntry)
		metrics := &interfaces.CacheMetrics{}

		_, err := cleanup.CleanupBySize(entries, metrics, -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "targetSize cannot be negative")
	})
}

func TestCacheCleanup_CleanupUnusedEntries(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	t.Run("CleanupUnusedEntries", func(t *testing.T) {
		now := time.Now()
		entries := map[string]*interfaces.CacheEntry{
			"recent": {
				Key:        "recent",
				Value:      "value1",
				Size:       6,
				AccessedAt: now.Add(-1 * time.Hour), // Recently accessed
			},
			"unused": {
				Key:        "unused",
				Value:      "value2",
				Size:       6,
				AccessedAt: now.Add(-25 * time.Hour), // Not accessed recently
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    12,
			CurrentEntries: 2,
		}

		unusedThreshold := 24 * time.Hour
		cleanedCount, err := cleanup.CleanupUnusedEntries(entries, metrics, unusedThreshold)
		require.NoError(t, err)

		assert.Equal(t, 1, cleanedCount)
		assert.Equal(t, 1, len(entries))
		assert.Contains(t, entries, "recent")
		assert.NotContains(t, entries, "unused")
		assert.Equal(t, int64(6), metrics.CurrentSize)
	})

	t.Run("InvalidThreshold", func(t *testing.T) {
		entries := make(map[string]*interfaces.CacheEntry)
		metrics := &interfaces.CacheMetrics{}

		_, err := cleanup.CleanupUnusedEntries(entries, metrics, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unusedThreshold must be positive")
	})
}

func TestCacheCleanup_ScheduledCleanup(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	config.SyncInterval = 30 * time.Minute
	cleanup := NewCacheCleanup(config)

	t.Run("CleanupNeeded", func(t *testing.T) {
		now := time.Now()
		lastCleanup := now.Add(-2 * time.Hour) // Long time ago

		entries := map[string]*interfaces.CacheEntry{
			"valid": {
				Key:        "valid",
				Value:      "value1",
				Size:       6,
				CreatedAt:  now.Add(-1 * time.Hour),
				AccessedAt: now.Add(-30 * time.Minute),
			},
			"expired": {
				Key:       "expired",
				Value:     "value2",
				Size:      6,
				CreatedAt: now.Add(-2 * time.Hour),
				ExpiresAt: &[]time.Time{now.Add(-1 * time.Hour)}[0], // Expired
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    12,
			CurrentEntries: 2,
		}

		result, err := cleanup.ScheduledCleanup(entries, metrics, lastCleanup)
		require.NoError(t, err)

		assert.NotNil(t, result)
		assert.Empty(t, result.SkipReason)
		assert.True(t, len(result.Tasks) >= 2) // At least expired and unused cleanup
		assert.True(t, result.TotalItemsRemoved > 0)
		assert.False(t, result.HasErrors)
	})

	t.Run("CleanupNotNeeded", func(t *testing.T) {
		now := time.Now()
		lastCleanup := now.Add(-10 * time.Minute) // Recent cleanup

		entries := make(map[string]*interfaces.CacheEntry)
		metrics := &interfaces.CacheMetrics{}

		result, err := cleanup.ScheduledCleanup(entries, metrics, lastCleanup)
		require.NoError(t, err)

		assert.NotNil(t, result)
		assert.NotEmpty(t, result.SkipReason)
		assert.Contains(t, result.SkipReason, "Cleanup not needed")
		assert.Empty(t, result.Tasks)
	})

	t.Run("SizeBasedCleanup", func(t *testing.T) {
		config := interfaces.DefaultCacheConfig()
		config.MaxSize = 10                    // Small max size
		config.SyncInterval = 30 * time.Minute // Set sync interval
		cleanup := NewCacheCleanup(config)

		now := time.Now()
		lastCleanup := now.Add(-2 * time.Hour)

		entries := map[string]*interfaces.CacheEntry{
			"entry1": {
				Key:        "entry1",
				Value:      "value1",
				Size:       8,
				CreatedAt:  now.Add(-1 * time.Hour),
				AccessedAt: now.Add(-1 * time.Hour), // Recently accessed, won't be cleaned as unused
			},
			"entry2": {
				Key:        "entry2",
				Value:      "value2",
				Size:       8,
				CreatedAt:  now.Add(-1 * time.Hour),
				AccessedAt: now.Add(-1 * time.Hour), // Recently accessed, won't be cleaned as unused
			},
		}

		metrics := &interfaces.CacheMetrics{
			CurrentSize:    16, // Exceeds max size of 10
			CurrentEntries: 2,
		}

		result, err := cleanup.ScheduledCleanup(entries, metrics, lastCleanup)
		require.NoError(t, err)

		// Should include size-based cleanup task
		sizeTaskFound := false
		for _, task := range result.Tasks {
			if task.Type == "size" {
				sizeTaskFound = true
				break
			}
		}
		assert.True(t, sizeTaskFound)
	})
}

func TestCacheCleanup_OptimizeCacheStructure(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"needs_metadata": {
			Key:       "needs_metadata",
			Value:     "value1",
			CreatedAt: now,
			Metadata:  nil, // Will be initialized
		},
		"needs_access_time": {
			Key:        "needs_access_time",
			Value:      "value2",
			CreatedAt:  now,
			AccessedAt: time.Time{}, // Zero time, will be fixed
			Metadata:   make(map[string]any),
		},
		"needs_update_time": {
			Key:       "needs_update_time",
			Value:     "value3",
			CreatedAt: now,
			UpdatedAt: now.Add(-1 * time.Hour), // Before created time, will be fixed
			Metadata:  make(map[string]any),
		},
		"already_optimized": {
			Key:        "already_optimized",
			Value:      "value4",
			CreatedAt:  now,
			UpdatedAt:  now,
			AccessedAt: now,
			Metadata:   make(map[string]any),
		},
	}

	optimizedCount := cleanup.optimizeCacheStructure(entries)

	assert.True(t, optimizedCount >= 3) // At least three optimizations were made

	// Check that metadata was initialized
	assert.NotNil(t, entries["needs_metadata"].Metadata)

	// Check that access time was set
	assert.False(t, entries["needs_access_time"].AccessedAt.IsZero())

	// Check that update time was fixed
	assert.False(t, entries["needs_update_time"].UpdatedAt.Before(entries["needs_update_time"].CreatedAt))
}

func TestCacheCleanup_SetConfig(t *testing.T) {
	config := interfaces.DefaultCacheConfig()
	cleanup := NewCacheCleanup(config)

	newConfig := &interfaces.CacheConfig{
		MaxSize: 2000,
	}

	cleanup.SetConfig(newConfig)
	assert.Equal(t, newConfig, cleanup.config)
}

func TestCompactionResult(t *testing.T) {
	result := &CompactionResult{
		StartTime:             time.Now(),
		EndTime:               time.Now().Add(1 * time.Second),
		InitialEntries:        10,
		FinalEntries:          8,
		ExpiredEntriesRemoved: 2,
		InitialSize:           1000,
		FinalSize:             800,
	}

	result.Duration = result.EndTime.Sub(result.StartTime)
	result.EntriesRemoved = result.InitialEntries - result.FinalEntries
	result.SizeReduced = result.InitialSize - result.FinalSize

	assert.Equal(t, 2, result.EntriesRemoved)
	assert.Equal(t, int64(200), result.SizeReduced)
	assert.True(t, result.Duration > 0)
}

func TestMaintenanceResult(t *testing.T) {
	result := &MaintenanceResult{
		StartTime: time.Now(),
		Success:   true,
		Tasks: []MaintenanceTask{
			{
				Name:    "Test Task",
				Success: true,
				Details: "Task completed successfully",
			},
		},
	}

	assert.True(t, result.Success)
	assert.Equal(t, 1, len(result.Tasks))
	assert.Equal(t, "Test Task", result.Tasks[0].Name)
}

func TestScheduledCleanupResult(t *testing.T) {
	result := &ScheduledCleanupResult{
		StartTime:         time.Now(),
		TotalItemsRemoved: 5,
		HasErrors:         false,
		Tasks: []CleanupTask{
			{
				Type:         "expired",
				Success:      true,
				ItemsRemoved: 3,
			},
			{
				Type:         "unused",
				Success:      true,
				ItemsRemoved: 2,
			},
		},
	}

	assert.Equal(t, 5, result.TotalItemsRemoved)
	assert.False(t, result.HasErrors)
	assert.Equal(t, 2, len(result.Tasks))
}
