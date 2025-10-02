// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheValidator handles cache validation and integrity checking.
type CacheValidator struct {
	cacheDir string
	config   *interfaces.CacheConfig
}

// NewCacheValidator creates a new cache validator instance.
func NewCacheValidator(cacheDir string, config *interfaces.CacheConfig) *CacheValidator {
	return &CacheValidator{
		cacheDir: cacheDir,
		config:   config,
	}
}

// ValidateCache validates the cache integrity including directory, permissions, and entries.
func (cv *CacheValidator) ValidateCache(entries map[string]*interfaces.CacheEntry) error {
	// Check if cache directory exists
	if _, err := os.Stat(cv.cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", cv.cacheDir)
	}

	// Check if cache directory is writable
	testFile := filepath.Join(cv.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	// Validate cache entries
	if err := cv.validateEntries(entries); err != nil {
		return fmt.Errorf("cache entry validation failed: %w", err)
	}

	return nil
}

// validateEntries validates individual cache entries for integrity.
func (cv *CacheValidator) validateEntries(entries map[string]*interfaces.CacheEntry) error {
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

		// Validate metadata if present
		if err := cv.validateEntryMetadata(key, entry); err != nil {
			issues = append(issues, err.Error())
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("validation issues found: %s", strings.Join(issues, "; "))
	}

	return nil
}

// validateEntryMetadata validates the metadata of a cache entry.
func (cv *CacheValidator) validateEntryMetadata(key string, entry *interfaces.CacheEntry) error {
	if entry.Metadata == nil {
		return nil // No metadata to validate
	}

	// Validate compression metadata if entry is compressed
	if entry.Compressed {
		if _, exists := entry.Metadata["compression_type"]; !exists {
			return fmt.Errorf("compressed entry missing compression_type metadata for key: %s", key)
		}

		if _, exists := entry.Metadata["original_size"]; !exists {
			return fmt.Errorf("compressed entry missing original_size metadata for key: %s", key)
		}

		// Validate original size is reasonable
		if originalSize, ok := entry.Metadata["original_size"].(int); ok {
			if originalSize < 0 {
				return fmt.Errorf("negative original_size in metadata for key: %s", key)
			}
		}
	}

	return nil
}

// RepairCache repairs corrupted cache data and returns the repaired entries and metrics.
func (cv *CacheValidator) RepairCache(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cv.cacheDir, 0750); err != nil {
		return nil, nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Clean up corrupted entries
	repairedEntries := cv.repairEntries(entries)

	// Recalculate metrics based on repaired entries
	repairedMetrics := cv.recalculateMetrics(repairedEntries, metrics)

	return repairedEntries, repairedMetrics, nil
}

// repairEntries repairs individual cache entries.
func (cv *CacheValidator) repairEntries(entries map[string]*interfaces.CacheEntry) map[string]*interfaces.CacheEntry {
	now := time.Now()
	repairedEntries := make(map[string]*interfaces.CacheEntry)

	for key, entry := range entries {
		if entry == nil {
			continue // Skip nil entries
		}

		// Create a copy to avoid modifying the original
		repairedEntry := *entry

		// Fix key mismatch
		if repairedEntry.Key != key {
			repairedEntry.Key = key
		}

		// Fix negative size
		if repairedEntry.Size < 0 {
			repairedEntry.Size = cv.calculateSize(repairedEntry.Value)
		}

		// Fix future timestamps
		if repairedEntry.CreatedAt.After(now) {
			repairedEntry.CreatedAt = now
		}
		if repairedEntry.UpdatedAt.After(now) {
			repairedEntry.UpdatedAt = now
		}
		if repairedEntry.AccessedAt.After(now) {
			repairedEntry.AccessedAt = now
		}

		// Fix negative access count
		if repairedEntry.AccessCount < 0 {
			repairedEntry.AccessCount = 0
		}

		// Initialize metadata if nil
		if repairedEntry.Metadata == nil {
			repairedEntry.Metadata = make(map[string]any)
		}

		// Skip expired entries during repair
		if repairedEntry.ExpiresAt != nil && repairedEntry.ExpiresAt.Before(now) {
			continue
		}

		repairedEntries[key] = &repairedEntry
	}

	return repairedEntries
}

// recalculateMetrics recalculates cache metrics based on current entries.
func (cv *CacheValidator) recalculateMetrics(entries map[string]*interfaces.CacheEntry, currentMetrics *interfaces.CacheMetrics) *interfaces.CacheMetrics {
	totalSize := int64(0)
	for _, entry := range entries {
		totalSize += entry.Size
	}

	// Preserve historical metrics but update current state
	metrics := &interfaces.CacheMetrics{
		CurrentSize:    totalSize,
		MaxSize:        cv.config.MaxSize,
		CurrentEntries: len(entries),
		MaxEntries:     cv.config.MaxEntries,
	}

	// Preserve historical counters if available
	if currentMetrics != nil {
		metrics.Hits = currentMetrics.Hits
		metrics.Misses = currentMetrics.Misses
		metrics.Gets = currentMetrics.Gets
		metrics.Sets = currentMetrics.Sets
		metrics.Deletes = currentMetrics.Deletes
		metrics.Evictions = currentMetrics.Evictions
	}

	return metrics
}

// ValidateConfiguration validates the cache configuration.
func (cv *CacheValidator) ValidateConfiguration() error {
	if cv.config == nil {
		return fmt.Errorf("cache configuration is nil")
	}

	// Validate size limits
	if cv.config.MaxSize < 0 {
		return fmt.Errorf("MaxSize cannot be negative: %d", cv.config.MaxSize)
	}

	if cv.config.MaxEntries < 0 {
		return fmt.Errorf("MaxEntries cannot be negative: %d", cv.config.MaxEntries)
	}

	// Validate eviction ratio
	if cv.config.EvictionRatio < 0 || cv.config.EvictionRatio > 1 {
		return fmt.Errorf("EvictionRatio must be between 0 and 1: %f", cv.config.EvictionRatio)
	}

	// Validate eviction policy
	validPolicies := map[string]bool{
		"lru":  true,
		"lfu":  true,
		"fifo": true,
		"ttl":  true,
	}
	if !validPolicies[cv.config.EvictionPolicy] {
		return fmt.Errorf("invalid eviction policy: %s", cv.config.EvictionPolicy)
	}

	// Validate compression settings
	if cv.config.EnableCompression {
		if cv.config.CompressionLevel < 1 || cv.config.CompressionLevel > 9 {
			return fmt.Errorf("CompressionLevel must be between 1 and 9: %d", cv.config.CompressionLevel)
		}

		validCompressionTypes := map[string]bool{
			"gzip": true,
		}
		if !validCompressionTypes[cv.config.CompressionType] {
			return fmt.Errorf("invalid compression type: %s", cv.config.CompressionType)
		}
	}

	// Validate TTL
	if cv.config.DefaultTTL < 0 {
		return fmt.Errorf("DefaultTTL cannot be negative: %v", cv.config.DefaultTTL)
	}

	return nil
}

// CacheHealthReport represents a cache health check report.
type CacheHealthReport struct {
	Timestamp        time.Time `json:"timestamp"`
	OverallHealth    string    `json:"overall_health"`
	Issues           []string  `json:"issues"`
	Recommendations  []string  `json:"recommendations"`
	ExpiredEntries   int       `json:"expired_entries"`
	CorruptedEntries int       `json:"corrupted_entries"`
	TotalEntries     int       `json:"total_entries"`
}

// CheckCacheHealth performs a comprehensive health check of the cache.
func (cv *CacheValidator) CheckCacheHealth(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (*CacheHealthReport, error) {
	report := &CacheHealthReport{
		Timestamp:       time.Now(),
		OverallHealth:   "healthy",
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	// Check directory health
	if err := cv.validateCacheDirectory(); err != nil {
		report.Issues = append(report.Issues, fmt.Sprintf("Directory issue: %v", err))
		report.OverallHealth = "unhealthy"
	}

	// Check configuration health
	if err := cv.ValidateConfiguration(); err != nil {
		report.Issues = append(report.Issues, fmt.Sprintf("Configuration issue: %v", err))
		report.OverallHealth = "degraded"
	}

	// Check entries health
	now := time.Now()
	expiredCount := 0
	corruptedCount := 0

	for key, entry := range entries {
		if entry == nil {
			corruptedCount++
			continue
		}

		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredCount++
		}

		// Check for basic corruption
		if entry.Size < 0 || entry.AccessCount < 0 || entry.Key != key {
			corruptedCount++
		}
	}

	// Assess health based on expired and corrupted entries
	totalEntries := len(entries)
	if totalEntries > 0 {
		expiredRatio := float64(expiredCount) / float64(totalEntries)
		corruptedRatio := float64(corruptedCount) / float64(totalEntries)

		if expiredRatio > 0.5 {
			report.Issues = append(report.Issues, fmt.Sprintf("High expired entry ratio: %.2f%%", expiredRatio*100))
			report.Recommendations = append(report.Recommendations, "Consider running cache cleanup")
			if report.OverallHealth == "healthy" {
				report.OverallHealth = "degraded"
			}
		}

		if corruptedRatio > 0.1 {
			report.Issues = append(report.Issues, fmt.Sprintf("Corrupted entries detected: %.2f%%", corruptedRatio*100))
			report.Recommendations = append(report.Recommendations, "Consider running cache repair")
			report.OverallHealth = "unhealthy"
		}
	}

	// Check metrics health
	if metrics != nil {
		if metrics.CurrentSize > metrics.MaxSize && metrics.MaxSize > 0 {
			report.Issues = append(report.Issues, "Cache size exceeds maximum limit")
			report.Recommendations = append(report.Recommendations, "Consider increasing MaxSize or running cleanup")
		}

		if metrics.CurrentEntries > metrics.MaxEntries && metrics.MaxEntries > 0 {
			report.Issues = append(report.Issues, "Cache entry count exceeds maximum limit")
			report.Recommendations = append(report.Recommendations, "Consider increasing MaxEntries or running cleanup")
		}

		// Check hit rate
		if metrics.Gets > 0 {
			hitRate := float64(metrics.Hits) / float64(metrics.Gets)
			if hitRate < 0.5 {
				report.Issues = append(report.Issues, fmt.Sprintf("Low cache hit rate: %.2f%%", hitRate*100))
				report.Recommendations = append(report.Recommendations, "Consider reviewing cache strategy or increasing cache size")
			}
		}
	}

	report.ExpiredEntries = expiredCount
	report.CorruptedEntries = corruptedCount
	report.TotalEntries = totalEntries

	return report, nil
}

// validateCacheDirectory validates the cache directory exists and is accessible.
func (cv *CacheValidator) validateCacheDirectory() error {
	// Check if cache directory exists
	info, err := os.Stat(cv.cacheDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", cv.cacheDir)
	}
	if err != nil {
		return fmt.Errorf("cannot access cache directory: %w", err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("cache path is not a directory: %s", cv.cacheDir)
	}

	// Check permissions
	testFile := filepath.Join(cv.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	return nil
}

// SetConfig updates the validator configuration.
func (cv *CacheValidator) SetConfig(config *interfaces.CacheConfig) {
	cv.config = config
}

// GetCacheDir returns the cache directory path.
func (cv *CacheValidator) GetCacheDir() string {
	return cv.cacheDir
}

// calculateSize estimates the size of a value in bytes (same as operations.go)
func (cv *CacheValidator) calculateSize(value any) int64 {
	// This is a rough estimation
	switch v := value.(type) {
	case string:
		return int64(len(v))
	case []byte:
		return int64(len(v))
	case int, int32, int64, float32, float64:
		return 8
	case bool:
		return 1
	default:
		// For complex types, use JSON marshaling to estimate size
		if data, err := json.Marshal(v); err == nil {
			return int64(len(data))
		}
		return 100 // Default estimate
	}
}
