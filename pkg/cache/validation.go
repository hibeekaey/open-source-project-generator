// Package cache provides cache validation and repair functionality.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheValidator handles cache validation and repair operations
type CacheValidator struct {
	cacheManager interfaces.CacheManager
	cacheDir     string
}

// NewCacheValidator creates a new cache validator instance
func NewCacheValidator(cacheManager interfaces.CacheManager) *CacheValidator {
	return &CacheValidator{
		cacheManager: cacheManager,
		cacheDir:     cacheManager.GetLocation(),
	}
}

// ValidateCache performs comprehensive cache validation
func (cv *CacheValidator) ValidateCache() (*interfaces.CacheHealth, error) {
	health := &interfaces.CacheHealth{
		Status:          interfaces.CacheStatusHealthy,
		LastCheck:       time.Now(),
		Issues:          make([]interfaces.CacheIssue, 0),
		Warnings:        make([]interfaces.CacheWarning, 0),
		Recommendations: make([]string, 0),
	}

	// Check cache directory accessibility
	if err := cv.validateCacheDirectory(health); err != nil {
		health.Status = interfaces.CacheStatusUnhealthy
		health.Issues = append(health.Issues, interfaces.CacheIssue{
			Type:        interfaces.CacheIssueTypePermission,
			Severity:    "error",
			Description: fmt.Sprintf("Cache directory not accessible: %v", err),
			DetectedAt:  time.Now(),
			Resolution:  "Check directory permissions and disk space",
			Fixable:     true,
		})
	}

	// Check cache file integrity
	if err := cv.validateCacheFile(health); err != nil {
		if health.Status == interfaces.CacheStatusHealthy {
			health.Status = interfaces.CacheStatusDegraded
		}
		health.Issues = append(health.Issues, interfaces.CacheIssue{
			Type:        interfaces.CacheIssueTypeCorruption,
			Severity:    "warning",
			Description: fmt.Sprintf("Cache file integrity issue: %v", err),
			DetectedAt:  time.Now(),
			Resolution:  "Run cache repair or clear cache",
			Fixable:     true,
		})
	}

	// Check disk space
	if err := cv.validateDiskSpace(health); err != nil {
		health.Warnings = append(health.Warnings, interfaces.CacheWarning{
			Type:        interfaces.CacheWarningTypeSize,
			Description: fmt.Sprintf("Low disk space: %v", err),
			Suggestion:  "Clean cache or increase available disk space",
		})
	}

	// Check cache performance
	if err := cv.validateCachePerformance(health); err != nil {
		health.Warnings = append(health.Warnings, interfaces.CacheWarning{
			Type:        interfaces.CacheWarningTypePerformance,
			Description: fmt.Sprintf("Cache performance issue: %v", err),
			Suggestion:  "Consider compacting or rebuilding cache",
		})
	}

	// Check for expired entries
	if err := cv.validateExpiredEntries(health); err != nil {
		health.Warnings = append(health.Warnings, interfaces.CacheWarning{
			Type:        interfaces.CacheWarningTypeExpiration,
			Description: fmt.Sprintf("Many expired entries found: %v", err),
			Suggestion:  "Run cache cleanup to remove expired entries",
		})
	}

	// Generate recommendations
	cv.generateRecommendations(health)

	return health, nil
}

// validateCacheDirectory checks if cache directory is accessible and writable
func (cv *CacheValidator) validateCacheDirectory(health *interfaces.CacheHealth) error {
	// Check if directory exists
	if _, err := os.Stat(cv.cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", cv.cacheDir)
	}

	// Check if directory is writable
	testFile := filepath.Join(cv.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}

	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	return nil
}

// validateCacheFile checks cache file integrity
func (cv *CacheValidator) validateCacheFile(health *interfaces.CacheHealth) error {
	cacheFilePath := filepath.Join(cv.cacheDir, "cache.json")

	// Check if cache file exists
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil // No cache file is not an error
	}

	// Try to read and parse cache file
	// #nosec G304 - cacheFilePath is constructed internally and safe
	file, err := os.Open(cacheFilePath)
	if err != nil {
		return fmt.Errorf("cannot open cache file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close cache file: %v\n", err)
		}
	}()

	var cacheData cacheFile
	if err := json.NewDecoder(file).Decode(&cacheData); err != nil {
		return fmt.Errorf("cache file is corrupted: %w", err)
	}

	// Validate cache data structure
	if cacheData.Version == "" {
		return fmt.Errorf("cache file missing version information")
	}

	if cacheData.Entries == nil {
		return fmt.Errorf("cache file missing entries")
	}

	return nil
}

// validateDiskSpace checks available disk space
func (cv *CacheValidator) validateDiskSpace(health *interfaces.CacheHealth) error {
	// Get cache directory info
	var stat os.FileInfo
	var err error

	if stat, err = os.Stat(cv.cacheDir); err != nil {
		return fmt.Errorf("cannot stat cache directory: %w", err)
	}

	// This is a simplified check - in production you might want to use syscalls
	// to get actual disk space information
	_ = stat

	// For now, just check if we can write a small file
	testFile := filepath.Join(cv.cacheDir, ".space_test")
	testData := make([]byte, 1024) // 1KB test

	if err := os.WriteFile(testFile, testData, 0600); err != nil {
		return fmt.Errorf("insufficient disk space: %w", err)
	}

	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}
	return nil
}

// validateCachePerformance checks cache performance metrics
func (cv *CacheValidator) validateCachePerformance(health *interfaces.CacheHealth) error {
	stats, err := cv.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("cannot get cache stats: %w", err)
	}

	// Check hit rate
	if stats.HitRate < 0.5 { // Less than 50% hit rate
		return fmt.Errorf("low cache hit rate: %.2f%%", stats.HitRate*100)
	}

	// Check cache size efficiency
	if stats.TotalSize > 0 && stats.TotalEntries > 0 {
		avgEntrySize := stats.TotalSize / int64(stats.TotalEntries)
		if avgEntrySize > 1024*1024 { // Average entry > 1MB
			return fmt.Errorf("large average entry size: %d bytes", avgEntrySize)
		}
	}

	return nil
}

// validateExpiredEntries checks for expired entries
func (cv *CacheValidator) validateExpiredEntries(health *interfaces.CacheHealth) error {
	stats, err := cv.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("cannot get cache stats: %w", err)
	}

	// Check if more than 25% of entries are expired
	if stats.TotalEntries > 0 {
		expiredRatio := float64(stats.ExpiredEntries) / float64(stats.TotalEntries)
		if expiredRatio > 0.25 {
			return fmt.Errorf("%.1f%% of entries are expired", expiredRatio*100)
		}
	}

	return nil
}

// generateRecommendations generates recommendations based on validation results
func (cv *CacheValidator) generateRecommendations(health *interfaces.CacheHealth) {
	// Recommendations based on issues
	for _, issue := range health.Issues {
		switch issue.Type {
		case interfaces.CacheIssueTypeCorruption:
			health.Recommendations = append(health.Recommendations, "Run 'generator cache repair' to fix corrupted cache data")
		case interfaces.CacheIssueTypePermission:
			health.Recommendations = append(health.Recommendations, "Check file system permissions for cache directory")
		case interfaces.CacheIssueTypeDiskSpace:
			health.Recommendations = append(health.Recommendations, "Free up disk space or run 'generator cache clean'")
		}
	}

	// Recommendations based on warnings
	for _, warning := range health.Warnings {
		switch warning.Type {
		case interfaces.CacheWarningTypeSize:
			health.Recommendations = append(health.Recommendations, "Run 'generator cache clean' to remove expired entries")
		case interfaces.CacheWarningTypeHitRate:
			health.Recommendations = append(health.Recommendations, "Consider preloading frequently used data")
		case interfaces.CacheWarningTypePerformance:
			health.Recommendations = append(health.Recommendations, "Run 'generator cache compact' to optimize cache performance")
		case interfaces.CacheWarningTypeExpiration:
			health.Recommendations = append(health.Recommendations, "Run 'generator cache clean' to remove expired entries")
		}
	}

	// General recommendations
	if len(health.Issues) == 0 && len(health.Warnings) == 0 {
		health.Recommendations = append(health.Recommendations, "Cache is healthy - no action required")
	}
}

// RepairCache attempts to repair cache issues
func (cv *CacheValidator) RepairCache() error {
	// Validate first to identify issues
	health, err := cv.ValidateCache()
	if err != nil {
		return fmt.Errorf("failed to validate cache before repair: %w", err)
	}

	// Repair each identified issue
	for _, issue := range health.Issues {
		if !issue.Fixable {
			continue
		}

		switch issue.Type {
		case interfaces.CacheIssueTypeCorruption:
			if err := cv.repairCorruption(); err != nil {
				return fmt.Errorf("failed to repair corruption: %w", err)
			}
		case interfaces.CacheIssueTypePermission:
			if err := cv.repairPermissions(); err != nil {
				return fmt.Errorf("failed to repair permissions: %w", err)
			}
		case interfaces.CacheIssueTypeDiskSpace:
			if err := cv.repairDiskSpace(); err != nil {
				return fmt.Errorf("failed to repair disk space: %w", err)
			}
		}
	}

	return nil
}

// repairCorruption repairs corrupted cache data
func (cv *CacheValidator) repairCorruption() error {
	// Use the cache manager's repair functionality
	return cv.cacheManager.RepairCache()
}

// repairPermissions attempts to fix permission issues
func (cv *CacheValidator) repairPermissions() error {
	// Try to create cache directory with proper permissions
	if err := os.MkdirAll(cv.cacheDir, 0750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Try to fix permissions on existing cache file
	cacheFilePath := filepath.Join(cv.cacheDir, "cache.json")
	if _, err := os.Stat(cacheFilePath); err == nil {
		if err := os.Chmod(cacheFilePath, 0600); err != nil {
			return fmt.Errorf("failed to fix cache file permissions: %w", err)
		}
	}

	return nil
}

// repairDiskSpace attempts to free up disk space
func (cv *CacheValidator) repairDiskSpace() error {
	// Clean expired entries
	if err := cv.cacheManager.Clean(); err != nil {
		return fmt.Errorf("failed to clean expired entries: %w", err)
	}

	// Compact cache
	if err := cv.cacheManager.CompactCache(); err != nil {
		return fmt.Errorf("failed to compact cache: %w", err)
	}

	return nil
}
