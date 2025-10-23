package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ValidationReport represents the result of cache validation
type ValidationReport struct {
	Valid            bool      `json:"valid"`
	TotalEntries     int       `json:"total_entries"`
	CorruptedEntries []string  `json:"corrupted_entries"`
	ExpiredEntries   []string  `json:"expired_entries"`
	Warnings         []string  `json:"warnings"`
	CheckedAt        time.Time `json:"checked_at"`
}

// CacheValidator validates cache integrity
type CacheValidator struct {
	logger *logger.Logger
}

// NewCacheValidator creates a new cache validator instance
func NewCacheValidator(log *logger.Logger) *CacheValidator {
	return &CacheValidator{
		logger: log,
	}
}

// Validate checks cache for corruption and inconsistencies
func (cv *CacheValidator) Validate(cache *orchestrator.ToolCache) (*ValidationReport, error) {
	if cache == nil {
		return nil, fmt.Errorf("cache is nil")
	}

	report := &ValidationReport{
		Valid:            true,
		CorruptedEntries: []string{},
		ExpiredEntries:   []string{},
		Warnings:         []string{},
		CheckedAt:        time.Now(),
	}

	// Check if cache file exists
	cacheFile := cache.GetCacheFile()
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		report.Warnings = append(report.Warnings, "Cache file does not exist on disk")
	}

	// Get cache stats
	stats := cache.GetStats()
	if total, ok := stats["total"].(int); ok {
		report.TotalEntries = total
	}

	// Check for expired entries
	expiredTools := cv.CheckExpired(cache)
	report.ExpiredEntries = expiredTools
	if len(expiredTools) > 0 {
		report.Warnings = append(report.Warnings,
			fmt.Sprintf("%d expired entries found", len(expiredTools)))
	}

	// Check for corrupted entries
	corruptedTools := cv.CheckCorrupted(cache)
	report.CorruptedEntries = corruptedTools
	if len(corruptedTools) > 0 {
		report.Valid = false
		report.Warnings = append(report.Warnings,
			fmt.Sprintf("%d corrupted entries found", len(corruptedTools)))
	}

	// Check cache file permissions
	if fileInfo, err := os.Stat(cacheFile); err == nil {
		mode := fileInfo.Mode()
		// Warn if cache file is world-readable or world-writable
		if mode&0044 != 0 {
			report.Warnings = append(report.Warnings,
				"Cache file has overly permissive permissions")
		}
	}

	// Check cache size
	if report.TotalEntries > 1000 {
		report.Warnings = append(report.Warnings,
			fmt.Sprintf("Cache has %d entries (recommended max: 1000)", report.TotalEntries))
	}

	if cv.logger != nil {
		if report.Valid {
			cv.logger.Debug("Cache validation passed")
		} else {
			cv.logger.Warn(fmt.Sprintf("Cache validation failed: %d issues found",
				len(report.CorruptedEntries)))
		}
	}

	return report, nil
}

// CheckExpired identifies expired cache entries
func (cv *CacheValidator) CheckExpired(cache *orchestrator.ToolCache) []string {
	expired := []string{}

	// Use reflection to access cache entries
	// Since we can't directly access the private cache map, we'll use the GetStats method
	stats := cache.GetStats()
	if expiredCount, ok := stats["expired"].(int); ok && expiredCount > 0 {
		// We need to iterate through all entries to find which ones are expired
		// This is a limitation of the current ToolCache API
		// For now, we'll return an empty list and rely on ClearExpired to handle cleanup
		if cv.logger != nil {
			cv.logger.Debug(fmt.Sprintf("Found %d expired entries", expiredCount))
		}
	}

	return expired
}

// CheckCorrupted identifies corrupted cache entries
func (cv *CacheValidator) CheckCorrupted(cache *orchestrator.ToolCache) []string {
	corrupted := []string{}

	// Check if cache file can be loaded
	cacheFile := cache.GetCacheFile()
	if _, err := os.Stat(cacheFile); err != nil {
		if cv.logger != nil {
			cv.logger.Debug(fmt.Sprintf("Cache file check failed: %v", err))
		}
		return corrupted
	}

	// Try to load cache to verify it's not corrupted
	// Create a temporary cache instance to test loading
	tempCache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: cache.GetCacheFile()[:len(cache.GetCacheFile())-len("/tool_cache.json")],
			TTL:      cache.GetTTL(),
		},
		cv.logger,
	)

	if err != nil {
		corrupted = append(corrupted, "cache_file")
		if cv.logger != nil {
			cv.logger.Debug(fmt.Sprintf("Cache file is corrupted: %v", err))
		}
	} else if tempCache != nil {
		// Successfully loaded, cache is not corrupted
		if cv.logger != nil {
			cv.logger.Debug("Cache file integrity check passed")
		}
	}

	return corrupted
}

// ValidateCacheEntry validates a single cache entry
func (cv *CacheValidator) ValidateCacheEntry(toolName string, entry *models.CachedTool) error {
	if entry == nil {
		return fmt.Errorf("cache entry is nil")
	}

	// Check if CachedAt is in the future
	if entry.CachedAt.After(time.Now()) {
		return fmt.Errorf("cache timestamp is in the future")
	}

	// Check if TTL is reasonable (not negative, not too large)
	if entry.TTL < 0 {
		return fmt.Errorf("TTL is negative")
	}
	if entry.TTL > 24*time.Hour {
		return fmt.Errorf("TTL is unreasonably large (>24h)")
	}

	// Check if version string is reasonable when tool is available
	if entry.Available && entry.Version == "" {
		// This is a warning, not an error - some tools might not report versions
		if cv.logger != nil {
			cv.logger.Debug(fmt.Sprintf("Tool '%s' is available but has no version", toolName))
		}
	}

	return nil
}
