// Package operations provides cache operation implementations.
package operations

import (
	"regexp"
	"sort"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// OperationCallbacks holds callback functions for cache events.
type OperationCallbacks struct {
	OnHit      func(key string)
	OnMiss     func(key string)
	OnEviction func(key string, reason string)
}

// CacheOperations coordinates all cache operations.
type CacheOperations struct {
	config    *interfaces.CacheConfig
	callbacks *OperationCallbacks
	get       *GetOperation
	set       *SetOperation
	delete    *DeleteOperation
	cleanup   *CleanupOperation
}

// NewCacheOperations creates a new cache operations coordinator.
func NewCacheOperations(config *interfaces.CacheConfig) *CacheOperations {
	callbacks := &OperationCallbacks{}

	return &CacheOperations{
		config:    config,
		callbacks: callbacks,
		get:       NewGetOperation(callbacks),
		set:       NewSetOperation(config, callbacks),
		delete:    NewDeleteOperation(callbacks),
		cleanup:   NewCleanupOperation(callbacks),
	}
}

// Get retrieves a value from the cache.
func (co *CacheOperations) Get(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (any, error) {
	return co.get.Execute(key, entries, metrics)
}

// Set stores a value in the cache.
func (co *CacheOperations) Set(key string, value any, ttl time.Duration, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	return co.set.Execute(key, value, ttl, entries, metrics)
}

// Delete removes a value from the cache.
func (co *CacheOperations) Delete(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	return co.delete.Execute(key, entries, metrics)
}

// Clear removes all values from the cache.
func (co *CacheOperations) Clear(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) {
	co.delete.Clear(entries, metrics)
}

// Clean removes expired entries from the cache.
func (co *CacheOperations) Clean(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) []string {
	return co.cleanup.Execute(entries, metrics)
}

// Compact optimizes the cache structure.
func (co *CacheOperations) Compact(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	return co.cleanup.Compact(entries, metrics)
}

// Exists checks if a key exists in the cache.
func (co *CacheOperations) Exists(key string, entries map[string]*interfaces.CacheEntry) bool {
	return co.get.Exists(key, entries)
}

// GetKeys returns all non-expired cache keys.
func (co *CacheOperations) GetKeys(entries map[string]*interfaces.CacheEntry) []string {
	keys := make([]string, 0, len(entries))
	now := time.Now()

	for key, entry := range entries {
		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// GetKeysByPattern returns cache keys matching a pattern.
func (co *CacheOperations) GetKeysByPattern(pattern string, entries map[string]*interfaces.CacheEntry) ([]string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	now := time.Now()

	for key, entry := range entries {
		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}

		if regex.MatchString(key) {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	return keys, nil
}

// GetExpiredKeys returns all expired keys.
func (co *CacheOperations) GetExpiredKeys(entries map[string]*interfaces.CacheEntry) []string {
	return co.cleanup.GetExpiredKeys(entries)
}

// SetCallbacks sets the callback functions for cache events.
func (co *CacheOperations) SetCallbacks(onHit, onMiss func(string), onEviction func(string, string)) {
	co.callbacks.OnHit = onHit
	co.callbacks.OnMiss = onMiss
	co.callbacks.OnEviction = onEviction
}

// SetConfig updates the operations configuration.
func (co *CacheOperations) SetConfig(config *interfaces.CacheConfig) {
	co.config = config
	co.set.config = config
}

// getSortedEntriesForEviction returns entries sorted by eviction policy (used by set operation)
func (co *CacheOperations) getSortedEntriesForEviction(entries map[string]*interfaces.CacheEntry) []*interfaces.CacheEntry {
	entryList := make([]*interfaces.CacheEntry, 0, len(entries))
	for _, entry := range entries {
		entryList = append(entryList, entry)
	}

	switch co.config.EvictionPolicy {
	case "lru": // Least Recently Used
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessedAt.Before(entryList[j].AccessedAt)
		})
	case "lfu": // Least Frequently Used
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessCount < entryList[j].AccessCount
		})
	case "fifo": // First In, First Out
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].CreatedAt.Before(entryList[j].CreatedAt)
		})
	case "ttl": // Shortest TTL first
		sort.Slice(entryList, func(i, j int) bool {
			if entryList[i].ExpiresAt == nil && entryList[j].ExpiresAt == nil {
				return entryList[i].CreatedAt.Before(entryList[j].CreatedAt)
			}
			if entryList[i].ExpiresAt == nil {
				return false
			}
			if entryList[j].ExpiresAt == nil {
				return true
			}
			return entryList[i].ExpiresAt.Before(*entryList[j].ExpiresAt)
		})
	default: // Default to LRU
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessedAt.Before(entryList[j].AccessedAt)
		})
	}

	return entryList
}
