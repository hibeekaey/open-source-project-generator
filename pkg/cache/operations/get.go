// Package operations provides cache operation implementations.
package operations

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// GetOperation handles cache get operations.
type GetOperation struct {
	callbacks *OperationCallbacks
}

// NewGetOperation creates a new get operation handler.
func NewGetOperation(callbacks *OperationCallbacks) *GetOperation {
	return &GetOperation{
		callbacks: callbacks,
	}
}

// Execute performs a cache get operation.
func (g *GetOperation) Execute(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (any, error) {
	entry, exists := entries[key]
	if !exists {
		metrics.Misses++
		metrics.Gets++
		if g.callbacks != nil && g.callbacks.OnMiss != nil {
			g.callbacks.OnMiss(key)
		}
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Check if entry is expired
	now := time.Now()
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
		delete(entries, key)
		metrics.Misses++
		metrics.Gets++
		if g.callbacks != nil && g.callbacks.OnMiss != nil {
			g.callbacks.OnMiss(key)
		}
		return nil, fmt.Errorf("key expired: %s", key)
	}

	// Update access information
	entry.AccessedAt = now
	entry.AccessCount++

	metrics.Hits++
	metrics.Gets++
	if g.callbacks != nil && g.callbacks.OnHit != nil {
		g.callbacks.OnHit(key)
	}

	return entry.Value, nil
}

// Exists checks if a key exists and is not expired.
func (g *GetOperation) Exists(key string, entries map[string]*interfaces.CacheEntry) bool {
	entry, exists := entries[key]
	if !exists {
		return false
	}

	// Check if entry is expired
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}
