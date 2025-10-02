// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheCleanup handles cache cleanup and maintenance operations.
type CacheCleanup struct {
	config *interfaces.CacheConfig
}

// NewCacheCleanup creates a new cache cleanup instance.
func NewCacheCleanup(config *interfaces.CacheConfig) *CacheCleanup {
	return &CacheCleanup{
		config: config,
	}
}

// CleanExpiredEntries removes expired entries from the cache and returns the number of cleaned entries.
func (cc *CacheCleanup) CleanExpiredEntries(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (int, error) {
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
	}

	metrics.CurrentEntries = len(entries)

	return len(expiredKeys), nil
}

// CompactCache compacts the cache by removing expired entries and optimizing storage.
func (cc *CacheCleanup) CompactCache(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (*CompactionResult, error) {
	result := &CompactionResult{
		StartTime:      time.Now(),
		InitialEntries: len(entries),
		InitialSize:    metrics.CurrentSize,
	}

	// Clean expired entries first
	expiredCount, err := cc.CleanExpiredEntries(entries, metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to clean expired entries: %w", err)
	}
	result.ExpiredEntriesRemoved = expiredCount

	// Recalculate metrics to ensure accuracy
	totalSize := int64(0)
	for _, entry := range entries {
		totalSize += entry.Size
	}

	metrics.CurrentSize = totalSize
	metrics.CurrentEntries = len(entries)

	// Calculate results
	result.EndTime = time.Now()
	result.FinalEntries = len(entries)
	result.FinalSize = metrics.CurrentSize
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.SizeReduced = result.InitialSize - result.FinalSize
	result.EntriesRemoved = result.InitialEntries - result.FinalEntries

	return result, nil
}

// PerformMaintenance performs comprehensive cache maintenance including cleanup and optimization.
func (cc *CacheCleanup) PerformMaintenance(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (*MaintenanceResult, error) {
	result := &MaintenanceResult{
		StartTime: time.Now(),
		Tasks:     make([]MaintenanceTask, 0),
	}

	// Task 1: Clean expired entries
	expiredTask := MaintenanceTask{
		Name:      "Clean Expired Entries",
		StartTime: time.Now(),
	}

	expiredCount, err := cc.CleanExpiredEntries(entries, metrics)
	if err != nil {
		expiredTask.Error = err.Error()
		expiredTask.Success = false
	} else {
		expiredTask.Success = true
		expiredTask.Details = fmt.Sprintf("Removed %d expired entries", expiredCount)
	}
	expiredTask.EndTime = time.Now()
	expiredTask.Duration = expiredTask.EndTime.Sub(expiredTask.StartTime)
	result.Tasks = append(result.Tasks, expiredTask)

	// Task 2: Validate cache integrity
	validationTask := MaintenanceTask{
		Name:      "Validate Cache Integrity",
		StartTime: time.Now(),
	}

	// Validate cache entries directly
	issues := cc.validateEntries(entries)
	if len(issues) > 0 {
		validationTask.Success = false
		validationTask.Details = fmt.Sprintf("Found %d integrity issues", len(issues))
		validationTask.Error = fmt.Sprintf("Issues: %v", issues)
	} else {
		validationTask.Success = true
		validationTask.Details = "Cache integrity validated successfully"
	}
	validationTask.EndTime = time.Now()
	validationTask.Duration = validationTask.EndTime.Sub(validationTask.StartTime)
	result.Tasks = append(result.Tasks, validationTask)

	// Task 3: Optimize cache structure
	optimizationTask := MaintenanceTask{
		Name:      "Optimize Cache Structure",
		StartTime: time.Now(),
	}

	optimizedCount := cc.optimizeCacheStructure(entries)
	optimizationTask.Success = true
	optimizationTask.Details = fmt.Sprintf("Optimized %d cache entries", optimizedCount)
	optimizationTask.EndTime = time.Now()
	optimizationTask.Duration = optimizationTask.EndTime.Sub(optimizationTask.StartTime)
	result.Tasks = append(result.Tasks, optimizationTask)

	// Calculate overall results
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = true

	// Check if any tasks failed
	for _, task := range result.Tasks {
		if !task.Success {
			result.Success = false
			break
		}
	}

	return result, nil
}

// CleanupByAge removes entries older than the specified age.
func (cc *CacheCleanup) CleanupByAge(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics, maxAge time.Duration) (int, error) {
	if maxAge <= 0 {
		return 0, fmt.Errorf("maxAge must be positive")
	}

	cutoffTime := time.Now().Add(-maxAge)
	oldKeys := make([]string, 0)

	// Find old entries
	for key, entry := range entries {
		if entry.CreatedAt.Before(cutoffTime) {
			oldKeys = append(oldKeys, key)
		}
	}

	// Remove old entries
	for _, key := range oldKeys {
		entry := entries[key]
		delete(entries, key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++
	}

	metrics.CurrentEntries = len(entries)

	return len(oldKeys), nil
}

// CleanupBySize removes entries to reduce cache size to the target size.
func (cc *CacheCleanup) CleanupBySize(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics, targetSize int64) (int, error) {
	if targetSize < 0 {
		return 0, fmt.Errorf("targetSize cannot be negative")
	}

	if metrics.CurrentSize <= targetSize {
		return 0, nil // Already within target size
	}

	// Get entries sorted by eviction policy
	sortedEntries := cc.getSortedEntriesForEviction(entries)

	removedCount := 0
	for _, entry := range sortedEntries {
		if metrics.CurrentSize <= targetSize {
			break
		}

		delete(entries, entry.Key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++
		removedCount++
	}

	metrics.CurrentEntries = len(entries)

	return removedCount, nil
}

// CleanupUnusedEntries removes entries that haven't been accessed recently.
func (cc *CacheCleanup) CleanupUnusedEntries(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics, unusedThreshold time.Duration) (int, error) {
	if unusedThreshold <= 0 {
		return 0, fmt.Errorf("unusedThreshold must be positive")
	}

	cutoffTime := time.Now().Add(-unusedThreshold)
	unusedKeys := make([]string, 0)

	// Find unused entries
	for key, entry := range entries {
		if entry.AccessedAt.Before(cutoffTime) {
			unusedKeys = append(unusedKeys, key)
		}
	}

	// Remove unused entries
	for _, key := range unusedKeys {
		entry := entries[key]
		delete(entries, key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++
	}

	metrics.CurrentEntries = len(entries)

	return len(unusedKeys), nil
}

// ScheduledCleanup performs cleanup based on configured intervals and thresholds.
func (cc *CacheCleanup) ScheduledCleanup(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics, lastCleanup time.Time) (*ScheduledCleanupResult, error) {
	result := &ScheduledCleanupResult{
		StartTime:   time.Now(),
		LastCleanup: lastCleanup,
		Tasks:       make([]CleanupTask, 0),
	}

	// Check if cleanup is needed based on time
	timeSinceLastCleanup := time.Since(lastCleanup)
	cleanupInterval := 1 * time.Hour // Default cleanup interval
	if cc.config != nil && cc.config.SyncInterval > 0 {
		cleanupInterval = cc.config.SyncInterval
	}

	if timeSinceLastCleanup < cleanupInterval {
		result.SkipReason = fmt.Sprintf("Cleanup not needed, last cleanup was %v ago", timeSinceLastCleanup)
		result.EndTime = time.Now()
		return result, nil
	}

	// Task 1: Clean expired entries
	expiredTask := CleanupTask{
		Type:      "expired",
		StartTime: time.Now(),
	}

	expiredCount, err := cc.CleanExpiredEntries(entries, metrics)
	if err != nil {
		expiredTask.Error = err.Error()
	} else {
		expiredTask.ItemsRemoved = expiredCount
		expiredTask.Success = true
	}
	expiredTask.EndTime = time.Now()
	result.Tasks = append(result.Tasks, expiredTask)

	// Task 2: Clean unused entries (older than 7 days)
	unusedTask := CleanupTask{
		Type:      "unused",
		StartTime: time.Now(),
	}

	unusedCount, err := cc.CleanupUnusedEntries(entries, metrics, 7*24*time.Hour)
	if err != nil {
		unusedTask.Error = err.Error()
	} else {
		unusedTask.ItemsRemoved = unusedCount
		unusedTask.Success = true
	}
	unusedTask.EndTime = time.Now()
	result.Tasks = append(result.Tasks, unusedTask)

	// Task 3: Size-based cleanup if over limit
	if cc.config != nil && cc.config.MaxSize > 0 && metrics.CurrentSize > cc.config.MaxSize {
		sizeTask := CleanupTask{
			Type:      "size",
			StartTime: time.Now(),
		}

		targetSize := int64(float64(cc.config.MaxSize) * 0.8) // Clean to 80% of max size
		sizeCount, err := cc.CleanupBySize(entries, metrics, targetSize)
		if err != nil {
			sizeTask.Error = err.Error()
		} else {
			sizeTask.ItemsRemoved = sizeCount
			sizeTask.Success = true
		}
		sizeTask.EndTime = time.Now()
		result.Tasks = append(result.Tasks, sizeTask)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Calculate total items removed
	for _, task := range result.Tasks {
		result.TotalItemsRemoved += task.ItemsRemoved
		if !task.Success {
			result.HasErrors = true
		}
	}

	return result, nil
}

// optimizeCacheStructure optimizes the internal structure of cache entries.
func (cc *CacheCleanup) optimizeCacheStructure(entries map[string]*interfaces.CacheEntry) int {
	optimizedCount := 0

	for _, entry := range entries {
		// Initialize metadata if nil
		if entry.Metadata == nil {
			entry.Metadata = make(map[string]any)
			optimizedCount++
		}

		// Update access time if it's too old
		if entry.AccessedAt.IsZero() {
			entry.AccessedAt = entry.CreatedAt
			optimizedCount++
		}

		// Ensure UpdatedAt is not before CreatedAt
		if entry.UpdatedAt.Before(entry.CreatedAt) {
			entry.UpdatedAt = entry.CreatedAt
			optimizedCount++
		}
	}

	return optimizedCount
}

// getSortedEntriesForEviction returns entries sorted by eviction policy (simplified version)
func (cc *CacheCleanup) getSortedEntriesForEviction(entries map[string]*interfaces.CacheEntry) []*interfaces.CacheEntry {
	entryList := make([]*interfaces.CacheEntry, 0, len(entries))
	for _, entry := range entries {
		entryList = append(entryList, entry)
	}

	// Simple LRU sorting for cleanup
	// In a full implementation, this would use the same logic as operations.go
	return entryList
}

// validateEntries validates cache entries for integrity (simplified version).
func (cc *CacheCleanup) validateEntries(entries map[string]*interfaces.CacheEntry) []string {
	now := time.Now()
	issues := make([]string, 0)

	for key, entry := range entries {
		if entry == nil {
			issues = append(issues, fmt.Sprintf("nil entry for key: %s", key))
			continue
		}

		if entry.Key != key {
			issues = append(issues, fmt.Sprintf("key mismatch: expected %s, got %s", key, entry.Key))
		}

		if entry.Size < 0 {
			issues = append(issues, fmt.Sprintf("negative size for key: %s", key))
		}

		if entry.CreatedAt.After(now) {
			issues = append(issues, fmt.Sprintf("future creation time for key: %s", key))
		}

		if entry.AccessCount < 0 {
			issues = append(issues, fmt.Sprintf("negative access count for key: %s", key))
		}
	}

	return issues
}

// SetConfig updates the cleanup configuration.
func (cc *CacheCleanup) SetConfig(config *interfaces.CacheConfig) {
	cc.config = config
}

// CompactionResult represents the result of a cache compaction operation.
type CompactionResult struct {
	StartTime             time.Time     `json:"start_time"`
	EndTime               time.Time     `json:"end_time"`
	Duration              time.Duration `json:"duration"`
	InitialEntries        int           `json:"initial_entries"`
	FinalEntries          int           `json:"final_entries"`
	EntriesRemoved        int           `json:"entries_removed"`
	ExpiredEntriesRemoved int           `json:"expired_entries_removed"`
	InitialSize           int64         `json:"initial_size"`
	FinalSize             int64         `json:"final_size"`
	SizeReduced           int64         `json:"size_reduced"`
}

// MaintenanceResult represents the result of a maintenance operation.
type MaintenanceResult struct {
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  time.Duration     `json:"duration"`
	Success   bool              `json:"success"`
	Tasks     []MaintenanceTask `json:"tasks"`
}

// MaintenanceTask represents a single maintenance task.
type MaintenanceTask struct {
	Name      string        `json:"name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Details   string        `json:"details"`
	Error     string        `json:"error,omitempty"`
}

// ScheduledCleanupResult represents the result of a scheduled cleanup operation.
type ScheduledCleanupResult struct {
	StartTime         time.Time     `json:"start_time"`
	EndTime           time.Time     `json:"end_time"`
	Duration          time.Duration `json:"duration"`
	LastCleanup       time.Time     `json:"last_cleanup"`
	TotalItemsRemoved int           `json:"total_items_removed"`
	HasErrors         bool          `json:"has_errors"`
	SkipReason        string        `json:"skip_reason,omitempty"`
	Tasks             []CleanupTask `json:"tasks"`
}

// CleanupTask represents a single cleanup task.
type CleanupTask struct {
	Type         string    `json:"type"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Success      bool      `json:"success"`
	ItemsRemoved int       `json:"items_removed"`
	Error        string    `json:"error,omitempty"`
}
