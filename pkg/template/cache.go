// Package template provides template caching functionality for the
// Open Source Project Generator.
package template

import (
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateCache manages template caching operations including storage,
// retrieval, invalidation, and refresh mechanisms.
type TemplateCache struct {
	// cache stores template information by name
	cache map[string]*models.TemplateInfo

	// cacheTime tracks when the cache was last updated
	cacheTime time.Time

	// cacheTTL defines how long cache entries remain valid
	cacheTTL time.Duration

	// mutex protects concurrent access to cache
	mutex sync.RWMutex

	// discoveryFunc is used to refresh cache when needed
	discoveryFunc func() ([]*models.TemplateInfo, error)
}

// CacheConfig defines configuration options for template caching
type CacheConfig struct {
	// TTL defines cache time-to-live duration
	TTL time.Duration

	// MaxSize defines maximum number of cached templates (0 = unlimited)
	MaxSize int

	// EnableAutoRefresh enables automatic cache refresh
	EnableAutoRefresh bool

	// RefreshInterval defines how often to auto-refresh cache
	RefreshInterval time.Duration
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		TTL:               5 * time.Minute,
		MaxSize:           0, // unlimited
		EnableAutoRefresh: false,
		RefreshInterval:   10 * time.Minute,
	}
}

// NewTemplateCache creates a new template cache instance
func NewTemplateCache(config *CacheConfig, discoveryFunc func() ([]*models.TemplateInfo, error)) *TemplateCache {
	if config == nil {
		config = DefaultCacheConfig()
	}

	cache := &TemplateCache{
		cache:         make(map[string]*models.TemplateInfo),
		cacheTTL:      config.TTL,
		discoveryFunc: discoveryFunc,
	}

	// Start auto-refresh if enabled
	if config.EnableAutoRefresh && config.RefreshInterval > 0 {
		go cache.startAutoRefresh(config.RefreshInterval)
	}

	return cache
}

// Get retrieves a template from cache by name
func (tc *TemplateCache) Get(name string) (*models.TemplateInfo, bool) {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	template, exists := tc.cache[name]
	if !exists {
		return nil, false
	}

	// Check if cache is still valid
	if tc.isExpired() {
		return nil, false
	}

	// Return a copy to prevent external modifications
	templateCopy := *template
	return &templateCopy, true
}

// GetAll retrieves all cached templates
func (tc *TemplateCache) GetAll() ([]*models.TemplateInfo, bool) {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	// Check if cache is valid
	if tc.isExpired() || len(tc.cache) == 0 {
		return nil, false
	}

	// Return copies of all templates
	templates := make([]*models.TemplateInfo, 0, len(tc.cache))
	for _, template := range tc.cache {
		templateCopy := *template
		templates = append(templates, &templateCopy)
	}

	return templates, true
}

// Set stores templates in cache
func (tc *TemplateCache) Set(templates []*models.TemplateInfo) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// Clear existing cache
	tc.cache = make(map[string]*models.TemplateInfo)

	// Store new templates
	for _, template := range templates {
		// Store a copy to prevent external modifications
		templateCopy := *template
		tc.cache[template.Name] = &templateCopy
	}

	// Update cache time
	tc.cacheTime = time.Now()
}

// Put stores a single template in cache
func (tc *TemplateCache) Put(template *models.TemplateInfo) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// Store a copy to prevent external modifications
	templateCopy := *template
	tc.cache[template.Name] = &templateCopy

	// Update cache time if this is the first entry
	if len(tc.cache) == 1 {
		tc.cacheTime = time.Now()
	}
}

// Remove removes a template from cache by name
func (tc *TemplateCache) Remove(name string) bool {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	_, exists := tc.cache[name]
	if exists {
		delete(tc.cache, name)
	}

	return exists
}

// Clear removes all templates from cache
func (tc *TemplateCache) Clear() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.cache = make(map[string]*models.TemplateInfo)
	tc.cacheTime = time.Time{}
}

// Refresh refreshes the cache by calling the discovery function
func (tc *TemplateCache) Refresh() error {
	if tc.discoveryFunc == nil {
		return fmt.Errorf("no discovery function configured")
	}

	// Discover templates
	templates, err := tc.discoveryFunc()
	if err != nil {
		return fmt.Errorf("failed to refresh cache: %w", err)
	}

	// Update cache
	tc.Set(templates)

	return nil
}

// IsValid checks if the cache is still valid (not expired)
func (tc *TemplateCache) IsValid() bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return !tc.isExpired() && len(tc.cache) > 0
}

// IsExpired checks if the cache has expired
func (tc *TemplateCache) IsExpired() bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return tc.isExpired()
}

// Size returns the number of cached templates
func (tc *TemplateCache) Size() int {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return len(tc.cache)
}

// LastUpdated returns when the cache was last updated
func (tc *TemplateCache) LastUpdated() time.Time {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return tc.cacheTime
}

// GetCacheStats returns cache statistics
func (tc *TemplateCache) GetCacheStats() *CacheStats {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return &CacheStats{
		Size:        len(tc.cache),
		LastUpdated: tc.cacheTime,
		TTL:         tc.cacheTTL,
		IsValid:     !tc.isExpired() && len(tc.cache) > 0,
		IsExpired:   tc.isExpired(),
	}
}

// InvalidateTemplate invalidates a specific template in cache
func (tc *TemplateCache) InvalidateTemplate(name string) {
	tc.Remove(name)
}

// InvalidateAll invalidates all cached templates
func (tc *TemplateCache) InvalidateAll() {
	tc.Clear()
}

// SetTTL updates the cache TTL
func (tc *TemplateCache) SetTTL(ttl time.Duration) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.cacheTTL = ttl
}

// GetTTL returns the current cache TTL
func (tc *TemplateCache) GetTTL() time.Duration {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return tc.cacheTTL
}

// GetCachedTemplateNames returns names of all cached templates
func (tc *TemplateCache) GetCachedTemplateNames() []string {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	names := make([]string, 0, len(tc.cache))
	for name := range tc.cache {
		names = append(names, name)
	}

	return names
}

// HasTemplate checks if a template is cached
func (tc *TemplateCache) HasTemplate(name string) bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	_, exists := tc.cache[name]
	return exists && !tc.isExpired()
}

// GetOrRefresh gets a template from cache or refreshes if not found/expired
func (tc *TemplateCache) GetOrRefresh(name string) (*models.TemplateInfo, error) {
	// Try to get from cache first
	if template, found := tc.Get(name); found {
		return template, nil
	}

	// Cache miss or expired, refresh cache
	if err := tc.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh cache: %w", err)
	}

	// Try again after refresh
	if template, found := tc.Get(name); found {
		return template, nil
	}

	return nil, fmt.Errorf("template '%s' not found even after cache refresh", name)
}

// GetAllOrRefresh gets all templates from cache or refreshes if expired
func (tc *TemplateCache) GetAllOrRefresh() ([]*models.TemplateInfo, error) {
	// Try to get from cache first
	if templates, found := tc.GetAll(); found {
		return templates, nil
	}

	// Cache miss or expired, refresh cache
	if err := tc.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh cache: %w", err)
	}

	// Try again after refresh
	if templates, found := tc.GetAll(); found {
		return templates, nil
	}

	return nil, fmt.Errorf("no templates found even after cache refresh")
}

// ConvertToInterfaceTemplateInfos converts models.TemplateInfo to interfaces.TemplateInfo
func (tc *TemplateCache) ConvertToInterfaceTemplateInfos(templates []*models.TemplateInfo) []interfaces.TemplateInfo {
	result := make([]interfaces.TemplateInfo, len(templates))
	for i, tmpl := range templates {
		result[i] = tc.convertToInterfaceTemplateInfo(tmpl)
	}
	return result
}

// convertToInterfaceTemplateInfo converts a single models.TemplateInfo to interfaces.TemplateInfo
func (tc *TemplateCache) convertToInterfaceTemplateInfo(tmpl *models.TemplateInfo) interfaces.TemplateInfo {
	return interfaces.TemplateInfo{
		Name:         tmpl.Name,
		DisplayName:  tmpl.DisplayName,
		Description:  tmpl.Description,
		Category:     tmpl.Category,
		Technology:   tmpl.Technology,
		Version:      tmpl.Version,
		Tags:         tmpl.Tags,
		Dependencies: tmpl.Dependencies,
		Metadata: interfaces.TemplateMetadata{
			Author:     tmpl.Metadata.Author,
			License:    tmpl.Metadata.License,
			Repository: tmpl.Metadata.Repository,
			Homepage:   tmpl.Metadata.Homepage,
			Keywords:   tmpl.Metadata.Keywords,
			Created:    tmpl.Metadata.CreatedAt,
			Updated:    tmpl.Metadata.UpdatedAt,
		},
	}
}

// isExpired checks if cache has expired (must be called with lock held)
func (tc *TemplateCache) isExpired() bool {
	if tc.cacheTime.IsZero() {
		return true
	}
	return time.Since(tc.cacheTime) >= tc.cacheTTL
}

// startAutoRefresh starts automatic cache refresh in a goroutine
func (tc *TemplateCache) startAutoRefresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := tc.Refresh(); err != nil {
			// Log error but continue (in a real implementation, use proper logging)
			fmt.Printf("Auto-refresh failed: %v\n", err)
		}
	}
}

// CacheStats contains cache statistics
type CacheStats struct {
	Size        int           `json:"size"`
	LastUpdated time.Time     `json:"last_updated"`
	TTL         time.Duration `json:"ttl"`
	IsValid     bool          `json:"is_valid"`
	IsExpired   bool          `json:"is_expired"`
}

// String returns a string representation of cache stats
func (cs *CacheStats) String() string {
	return fmt.Sprintf("CacheStats{Size: %d, LastUpdated: %v, TTL: %v, IsValid: %t, IsExpired: %t}",
		cs.Size, cs.LastUpdated, cs.TTL, cs.IsValid, cs.IsExpired)
}
