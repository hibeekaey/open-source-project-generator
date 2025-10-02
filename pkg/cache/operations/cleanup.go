// Package operations provides cache operation implementations.
package operations

import (
	"sort"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CleanupOperation handles cache cleanup operations.
type CleanupOperation struct {
	callbacks *OperationCallbacks
}

// NewCleanupOperation creates a new cleanup operation handler.
func NewCleanupOperation(callbacks *OperationCallbacks) *CleanupOperation {
	return &CleanupOperation{
		callbacks: callbacks,
	}
}

// Execute performs cache cleanup by removing expired entries.
func (c *CleanupOperation) Execute(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) []string {
	now := time.Now()
	expiredKeys := make([]string, 0)

	// Find expired entries
	for key, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired entries
	for _, key := range expiredKeys {
		entry := entries[key]
		delete(entries, key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if c.callbacks != nil && c.callbacks.OnEviction != nil {
			c.callbacks.OnEviction(key, interfaces.EvictionReasonTTL)
		}
	}

	metrics.CurrentEntries = len(entries)
	return expiredKeys
}

// Compact removes expired entries and optimizes cache structure.
func (c *CleanupOperation) Compact(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// Clean expired entries first
	c.Execute(entries, metrics)

	// Recalculate metrics to ensure accuracy
	totalSize := int64(0)
	for _, entry := range entries {
		totalSize += entry.Size
	}

	metrics.CurrentSize = totalSize
	metrics.CurrentEntries = len(entries)

	return nil
}

// GetExpiredKeys returns all expired keys without removing them.
func (c *CleanupOperation) GetExpiredKeys(entries map[string]*interfaces.CacheEntry) []string {
	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	sort.Strings(expiredKeys)
	return expiredKeys
}
