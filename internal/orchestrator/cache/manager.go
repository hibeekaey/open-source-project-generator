package cache

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// CacheStats represents detailed cache statistics
type CacheStats struct {
	TotalEntries     int           `json:"total_entries"`
	AvailableTools   int           `json:"available_tools"`
	UnavailableTools int           `json:"unavailable_tools"`
	ExpiredEntries   int           `json:"expired_entries"`
	CacheFile        string        `json:"cache_file"`
	TTL              time.Duration `json:"ttl"`
	LastSaved        time.Time     `json:"last_saved"`
	TimeSinceCheck   time.Duration `json:"time_since_check"`
}

// CacheManager provides advanced cache management functionality
type CacheManager struct {
	cache     *orchestrator.ToolCache
	validator *CacheValidator
	exporter  *CacheExporter
	logger    *logger.Logger
}

// NewCacheManager creates a new cache manager instance
func NewCacheManager(cache *orchestrator.ToolCache, log *logger.Logger) *CacheManager {
	return &CacheManager{
		cache:     cache,
		validator: NewCacheValidator(log),
		exporter:  NewCacheExporter(log),
		logger:    log,
	}
}

// GetStats returns detailed cache statistics
func (cm *CacheManager) GetStats() *CacheStats {
	if cm.cache == nil {
		return &CacheStats{}
	}

	// Get basic stats from cache
	basicStats := cm.cache.GetStats()

	stats := &CacheStats{
		CacheFile: cm.cache.GetCacheFile(),
		TTL:       cm.cache.GetTTL(),
	}

	// Extract values from map
	if total, ok := basicStats["total"].(int); ok {
		stats.TotalEntries = total
	}
	if available, ok := basicStats["available"].(int); ok {
		stats.AvailableTools = available
	}
	if unavailable, ok := basicStats["unavailable"].(int); ok {
		stats.UnavailableTools = unavailable
	}
	if expired, ok := basicStats["expired"].(int); ok {
		stats.ExpiredEntries = expired
	}
	if lastSaved, ok := basicStats["last_saved"].(time.Time); ok {
		stats.LastSaved = lastSaved
		stats.TimeSinceCheck = time.Since(lastSaved)
	}

	return stats
}

// Validate checks cache integrity and reports issues
func (cm *CacheManager) Validate() (*ValidationReport, error) {
	if cm.cache == nil {
		return nil, fmt.Errorf("cache not initialized")
	}

	return cm.validator.Validate(cm.cache)
}

// Refresh re-checks all cached tools and updates their status
func (cm *CacheManager) Refresh(toolDiscovery ToolDiscoveryInterface) error {
	if cm.cache == nil {
		return fmt.Errorf("cache not initialized")
	}

	if cm.logger != nil {
		cm.logger.Info("Refreshing tool cache...")
	}

	// Get all registered tools
	allTools := toolDiscovery.ListRegisteredTools()

	refreshed := 0
	for _, toolName := range allTools {
		if cm.logger != nil {
			cm.logger.Debug(fmt.Sprintf("Refreshing tool: %s", toolName))
		}

		// Check availability
		available, _ := toolDiscovery.IsAvailable(toolName)

		// Get version if available
		version := ""
		if available {
			if v, err := toolDiscovery.GetVersion(toolName); err == nil {
				version = v
			}
		}

		// Update cache
		cm.cache.Set(toolName, available, version)
		refreshed++
	}

	// Clear expired entries
	cm.cache.ClearExpired()

	// Save cache to disk
	if err := cm.cache.Save(); err != nil {
		return fmt.Errorf("failed to save cache after refresh: %w", err)
	}

	if cm.logger != nil {
		cm.logger.Info(fmt.Sprintf("Cache refreshed: %d tools updated", refreshed))
	}

	return nil
}

// Export exports cache data to a portable format
func (cm *CacheManager) Export(outputPath string) error {
	if cm.cache == nil {
		return fmt.Errorf("cache not initialized")
	}

	return cm.exporter.Export(cm.cache, outputPath)
}

// Import imports cache data from a file
func (cm *CacheManager) Import(inputPath string) error {
	if cm.cache == nil {
		return fmt.Errorf("cache not initialized")
	}

	return cm.exporter.Import(cm.cache, inputPath)
}

// ToolDiscoveryInterface defines the interface for tool discovery operations
type ToolDiscoveryInterface interface {
	ListRegisteredTools() []string
	IsAvailable(toolName string) (bool, error)
	GetVersion(toolName string) (string, error)
}
