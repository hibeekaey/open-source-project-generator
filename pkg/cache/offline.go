// Package cache provides offline mode functionality for the
// Open Source Project Generator.
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// OfflineManager handles offline mode operations and cached data usage
type OfflineManager struct {
	cacheManager interfaces.CacheManager
	isOffline    bool
	config       *interfaces.CacheConfig
}

// NewOfflineManager creates a new offline manager instance
func NewOfflineManager(cacheManager interfaces.CacheManager) *OfflineManager {
	config, _ := cacheManager.GetCacheConfig()
	return &OfflineManager{
		cacheManager: cacheManager,
		isOffline:    cacheManager.IsOfflineMode(),
		config:       config,
	}
}

// EnableOfflineMode enables offline mode and validates cached data
func (om *OfflineManager) EnableOfflineMode() error {
	// Enable offline mode in cache manager
	if err := om.cacheManager.EnableOfflineMode(); err != nil {
		return fmt.Errorf("failed to enable offline mode: %w", err)
	}

	om.isOffline = true

	// Validate that essential cached data is available
	if err := om.validateOfflineData(); err != nil {
		return fmt.Errorf("offline mode validation failed: %w", err)
	}

	return nil
}

// DisableOfflineMode disables offline mode
func (om *OfflineManager) DisableOfflineMode() error {
	if err := om.cacheManager.DisableOfflineMode(); err != nil {
		return fmt.Errorf("failed to disable offline mode: %w", err)
	}

	om.isOffline = false
	return nil
}

// IsOfflineMode returns whether offline mode is currently enabled
func (om *OfflineManager) IsOfflineMode() bool {
	return om.isOffline
}

// validateOfflineData validates that essential data is available in cache for offline operation
func (om *OfflineManager) validateOfflineData() error {
	// Check if cache is accessible
	if err := om.cacheManager.ValidateCache(); err != nil {
		return fmt.Errorf("cache validation failed: %w", err)
	}

	// Check for essential cached data
	essentialKeys := []string{
		"templates:list",
		"templates:metadata",
		"versions:latest",
		"config:defaults",
	}

	missingKeys := make([]string, 0)
	for _, key := range essentialKeys {
		if !om.cacheManager.Exists(key) {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing essential cached data for offline mode: %v", missingKeys)
	}

	return nil
}

// GetCachedTemplates retrieves cached template data for offline use
func (om *OfflineManager) GetCachedTemplates() ([]interfaces.TemplateInfo, error) {
	if !om.isOffline {
		return nil, fmt.Errorf("not in offline mode")
	}

	data, err := om.cacheManager.Get("templates:list")
	if err != nil {
		return nil, fmt.Errorf("failed to get cached templates: %w", err)
	}

	templates, ok := data.([]interfaces.TemplateInfo)
	if !ok {
		return nil, fmt.Errorf("invalid cached template data format")
	}

	return templates, nil
}

// GetCachedTemplateMetadata retrieves cached template metadata for offline use
func (om *OfflineManager) GetCachedTemplateMetadata(templateName string) (*interfaces.TemplateMetadata, error) {
	if !om.isOffline {
		return nil, fmt.Errorf("not in offline mode")
	}

	key := fmt.Sprintf("templates:metadata:%s", templateName)
	data, err := om.cacheManager.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached template metadata: %w", err)
	}

	metadata, ok := data.(*interfaces.TemplateMetadata)
	if !ok {
		return nil, fmt.Errorf("invalid cached template metadata format")
	}

	return metadata, nil
}

// GetCachedVersions retrieves cached version data for offline use
func (om *OfflineManager) GetCachedVersions() (map[string]string, error) {
	if !om.isOffline {
		return nil, fmt.Errorf("not in offline mode")
	}

	data, err := om.cacheManager.Get("versions:latest")
	if err != nil {
		return nil, fmt.Errorf("failed to get cached versions: %w", err)
	}

	versions, ok := data.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid cached version data format")
	}

	return versions, nil
}

// CacheTemplateData caches template data for offline use
func (om *OfflineManager) CacheTemplateData(templates []interfaces.TemplateInfo) error {
	ttl := om.config.OfflineTTL
	if ttl == 0 {
		ttl = 7 * 24 * time.Hour // Default 7 days
	}

	if err := om.cacheManager.Set("templates:list", templates, ttl); err != nil {
		return fmt.Errorf("failed to cache template list: %w", err)
	}

	// Cache individual template metadata
	for _, template := range templates {
		key := fmt.Sprintf("templates:metadata:%s", template.Name)
		if err := om.cacheManager.Set(key, template.Metadata, ttl); err != nil {
			return fmt.Errorf("failed to cache template metadata for %s: %w", template.Name, err)
		}
	}

	return nil
}

// CacheVersionData caches version data for offline use
func (om *OfflineManager) CacheVersionData(versions map[string]string) error {
	ttl := om.config.OfflineTTL
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours for versions
	}

	if err := om.cacheManager.Set("versions:latest", versions, ttl); err != nil {
		return fmt.Errorf("failed to cache version data: %w", err)
	}

	return nil
}

// CacheConfigDefaults caches default configuration for offline use
func (om *OfflineManager) CacheConfigDefaults(config interface{}) error {
	ttl := om.config.OfflineTTL
	if ttl == 0 {
		ttl = 7 * 24 * time.Hour // Default 7 days
	}

	if err := om.cacheManager.Set("config:defaults", config, ttl); err != nil {
		return fmt.Errorf("failed to cache config defaults: %w", err)
	}

	return nil
}

// GetOfflineStatus returns detailed offline mode status
func (om *OfflineManager) GetOfflineStatus() (*OfflineStatus, error) {
	stats, err := om.cacheManager.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Check availability of essential data
	essentialData := make(map[string]bool)
	essentialKeys := []string{
		"templates:list",
		"templates:metadata",
		"versions:latest",
		"config:defaults",
	}

	for _, key := range essentialKeys {
		essentialData[key] = om.cacheManager.Exists(key)
	}

	// Calculate offline readiness score
	availableCount := 0
	for _, available := range essentialData {
		if available {
			availableCount++
		}
	}
	readinessScore := float64(availableCount) / float64(len(essentialKeys)) * 100

	return &OfflineStatus{
		Enabled:        om.isOffline,
		CacheLocation:  om.cacheManager.GetLocation(),
		CacheSize:      stats.TotalSize,
		CacheEntries:   stats.TotalEntries,
		EssentialData:  essentialData,
		ReadinessScore: readinessScore,
		LastSync:       time.Now(), // TODO: Track actual last sync time
	}, nil
}

// PreloadEssentialData preloads essential data for offline operation
func (om *OfflineManager) PreloadEssentialData() error {
	// This would typically fetch data from network sources and cache it
	// For now, we'll just validate that the cache is ready
	return om.validateOfflineData()
}

// SyncOfflineData synchronizes offline data with remote sources
func (om *OfflineManager) SyncOfflineData() error {
	if om.isOffline {
		return fmt.Errorf("cannot sync while in offline mode")
	}

	// This would typically:
	// 1. Fetch latest template data
	// 2. Fetch latest version data
	// 3. Update cached configuration defaults
	// 4. Cache the data for offline use

	// For now, just sync the cache to disk
	return om.cacheManager.SyncCache()
}

// OfflineStatus represents the current offline mode status
type OfflineStatus struct {
	Enabled        bool            `json:"enabled"`
	CacheLocation  string          `json:"cache_location"`
	CacheSize      int64           `json:"cache_size"`
	CacheEntries   int             `json:"cache_entries"`
	EssentialData  map[string]bool `json:"essential_data"`
	ReadinessScore float64         `json:"readiness_score"`
	LastSync       time.Time       `json:"last_sync"`
}

// DetectOfflineMode detects if the system should operate in offline mode
func DetectOfflineMode() bool {
	// Check for network connectivity
	// This is a basic implementation - in production you might want more sophisticated detection

	// Check if we're in a CI environment (often offline for caching)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		return true
	}

	// Check for explicit offline flag in environment
	if os.Getenv("GENERATOR_OFFLINE") == "true" {
		return true
	}

	// Check if cache directory exists and has recent data
	cacheDir := os.Getenv("GENERATOR_CACHE_DIR")
	if cacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".generator", "cache")
	}

	cacheFile := filepath.Join(cacheDir, "cache.json")
	if info, err := os.Stat(cacheFile); err == nil {
		// If cache file is less than 24 hours old, we might be able to work offline
		if time.Since(info.ModTime()) < 24*time.Hour {
			return false // Recent cache, but don't force offline mode
		}
	}

	return false
}

// ValidateOfflineCapabilities validates that the system can operate offline
func ValidateOfflineCapabilities(cacheManager interfaces.CacheManager) error {
	// Check cache accessibility
	if err := cacheManager.ValidateCache(); err != nil {
		return fmt.Errorf("cache not accessible: %w", err)
	}

	// Check for minimum required cached data
	requiredKeys := []string{
		"templates:list",
		"config:defaults",
	}

	for _, key := range requiredKeys {
		if !cacheManager.Exists(key) {
			return fmt.Errorf("missing required cached data: %s", key)
		}
	}

	return nil
}
