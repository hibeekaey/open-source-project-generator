// Package validation provides cache validation and integrity checking.
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Validator handles cache validation and integrity checking.
type Validator struct {
	cacheDir string
	config   *interfaces.CacheConfig
}

// NewValidator creates a new cache validator.
func NewValidator(cacheDir string, config *interfaces.CacheConfig) *Validator {
	return &Validator{
		cacheDir: cacheDir,
		config:   config,
	}
}

// ValidateCache validates the complete cache including directory, permissions, and entries.
func (v *Validator) ValidateCache(entries map[string]*interfaces.CacheEntry) error {
	// Validate directory
	if err := v.validateDirectory(); err != nil {
		return fmt.Errorf("directory validation failed: %w", err)
	}

	// Validate configuration
	if err := v.validateConfiguration(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Validate entries
	if err := v.validateEntries(entries); err != nil {
		return fmt.Errorf("entry validation failed: %w", err)
	}

	return nil
}

// validateDirectory validates the cache directory exists and is accessible.
func (v *Validator) validateDirectory() error {
	// Check if cache directory exists
	info, err := os.Stat(v.cacheDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", v.cacheDir)
	}
	if err != nil {
		return fmt.Errorf("cannot access cache directory: %w", err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("cache path is not a directory: %s", v.cacheDir)
	}

	// Check permissions
	testFile := filepath.Join(v.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	return nil
}

// validateConfiguration validates the cache configuration.
func (v *Validator) validateConfiguration() error {
	if v.config == nil {
		return fmt.Errorf("cache configuration is nil")
	}

	// Validate size limits
	if v.config.MaxSize < 0 {
		return fmt.Errorf("MaxSize cannot be negative: %d", v.config.MaxSize)
	}

	if v.config.MaxEntries < 0 {
		return fmt.Errorf("MaxEntries cannot be negative: %d", v.config.MaxEntries)
	}

	// Validate eviction ratio
	if v.config.EvictionRatio < 0 || v.config.EvictionRatio > 1 {
		return fmt.Errorf("EvictionRatio must be between 0 and 1: %f", v.config.EvictionRatio)
	}

	// Validate eviction policy
	validPolicies := map[string]bool{
		"lru":  true,
		"lfu":  true,
		"fifo": true,
		"ttl":  true,
	}
	if !validPolicies[v.config.EvictionPolicy] {
		return fmt.Errorf("invalid eviction policy: %s", v.config.EvictionPolicy)
	}

	// Validate compression settings
	if v.config.EnableCompression {
		if v.config.CompressionLevel < 1 || v.config.CompressionLevel > 9 {
			return fmt.Errorf("CompressionLevel must be between 1 and 9: %d", v.config.CompressionLevel)
		}

		validCompressionTypes := map[string]bool{
			"gzip": true,
		}
		if !validCompressionTypes[v.config.CompressionType] {
			return fmt.Errorf("invalid compression type: %s", v.config.CompressionType)
		}
	}

	// Validate TTL
	if v.config.DefaultTTL < 0 {
		return fmt.Errorf("DefaultTTL cannot be negative: %v", v.config.DefaultTTL)
	}

	return nil
}

// validateEntries validates individual cache entries for integrity.
func (v *Validator) validateEntries(entries map[string]*interfaces.CacheEntry) error {
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
		if err := v.validateEntryMetadata(key, entry); err != nil {
			issues = append(issues, err.Error())
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("validation issues found: %s", strings.Join(issues, "; "))
	}

	return nil
}

// validateEntryMetadata validates the metadata of a cache entry.
func (v *Validator) validateEntryMetadata(key string, entry *interfaces.CacheEntry) error {
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

// RepairEntries repairs corrupted cache entries.
func (v *Validator) RepairEntries(entries map[string]*interfaces.CacheEntry) map[string]*interfaces.CacheEntry {
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
			repairedEntry.Size = v.calculateSize(repairedEntry.Value)
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

// CheckHealth performs a comprehensive health check of the cache.
func (v *Validator) CheckHealth(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (*HealthReport, error) {
	report := &HealthReport{
		Timestamp:       time.Now(),
		OverallHealth:   "healthy",
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	// Check directory health
	if err := v.validateDirectory(); err != nil {
		report.Issues = append(report.Issues, fmt.Sprintf("Directory issue: %v", err))
		report.OverallHealth = "unhealthy"
	}

	// Check configuration health
	if err := v.validateConfiguration(); err != nil {
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

// HealthReport represents a cache health check report.
type HealthReport struct {
	Timestamp        time.Time `json:"timestamp"`
	OverallHealth    string    `json:"overall_health"`
	Issues           []string  `json:"issues"`
	Recommendations  []string  `json:"recommendations"`
	ExpiredEntries   int       `json:"expired_entries"`
	CorruptedEntries int       `json:"corrupted_entries"`
	TotalEntries     int       `json:"total_entries"`
}

// SetConfig updates the validator configuration.
func (v *Validator) SetConfig(config *interfaces.CacheConfig) {
	v.config = config
}

// GetCacheDir returns the cache directory path.
func (v *Validator) GetCacheDir() string {
	return v.cacheDir
}

// calculateSize estimates the size of a value in bytes.
func (v *Validator) calculateSize(value any) int64 {
	// This is a rough estimation
	switch val := value.(type) {
	case string:
		return int64(len(val))
	case []byte:
		return int64(len(val))
	case int, int32, int64, float32, float64:
		return 8
	case bool:
		return 1
	default:
		// For complex types, use JSON marshaling to estimate size
		if data, err := json.Marshal(val); err == nil {
			return int64(len(data))
		}
		return 100 // Default estimate
	}
}
