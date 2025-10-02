// Package operations provides cache operation implementations.
package operations

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DeleteOperation handles cache delete operations.
type DeleteOperation struct {
	callbacks *OperationCallbacks
}

// NewDeleteOperation creates a new delete operation handler.
func NewDeleteOperation(callbacks *OperationCallbacks) *DeleteOperation {
	return &DeleteOperation{
		callbacks: callbacks,
	}
}

// Execute performs a cache delete operation.
func (d *DeleteOperation) Execute(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	entry, exists := entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	delete(entries, key)
	metrics.CurrentSize -= entry.Size
	metrics.CurrentEntries = len(entries)
	metrics.Deletes++

	if d.callbacks != nil && d.callbacks.OnEviction != nil {
		d.callbacks.OnEviction(key, interfaces.EvictionReasonManual)
	}

	return nil
}

// Clear removes all entries from the cache.
func (d *DeleteOperation) Clear(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) {
	// Clear all entries
	for key := range entries {
		delete(entries, key)
	}

	metrics.CurrentSize = 0
	metrics.CurrentEntries = 0
}
