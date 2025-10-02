// Package performance provides lazy loading implementations for expensive operations
package performance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TemplateLazyLoader provides lazy loading for template operations
type TemplateLazyLoader struct {
	templateManager interfaces.TemplateManager
	cache           interfaces.CacheManager
	loaded          bool
	data            interface{}
	mutex           sync.RWMutex
	cacheKey        string
	loadFunc        func(ctx context.Context) (interface{}, error)
}

// NewTemplateLazyLoader creates a new template lazy loader
func NewTemplateLazyLoader(templateManager interfaces.TemplateManager, cache interfaces.CacheManager) *TemplateLazyLoader {
	return &TemplateLazyLoader{
		templateManager: templateManager,
		cache:           cache,
		cacheKey:        "lazy:templates:list",
		loadFunc: func(ctx context.Context) (interface{}, error) {
			return templateManager.ListTemplates(interfaces.TemplateFilter{})
		},
	}
}

// Load loads the templates if not already loaded
func (tll *TemplateLazyLoader) Load(ctx context.Context) (interface{}, error) {
	tll.mutex.Lock()
	defer tll.mutex.Unlock()

	if tll.loaded && tll.data != nil {
		return tll.data, nil
	}

	// Try cache first
	if tll.cache != nil {
		if cached, err := tll.cache.Get(tll.cacheKey); err == nil {
			tll.data = cached
			tll.loaded = true
			return tll.data, nil
		}
	}

	// Load from source
	data, err := tll.loadFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	tll.data = data
	tll.loaded = true

	// Cache the result
	if tll.cache != nil {
		tll.cache.Set(tll.cacheKey, data, 30*time.Minute)
	}

	return tll.data, nil
}

// IsLoaded returns whether the data has been loaded
func (tll *TemplateLazyLoader) IsLoaded() bool {
	tll.mutex.RLock()
	defer tll.mutex.RUnlock()
	return tll.loaded
}

// GetCacheKey returns the cache key for this loader
func (tll *TemplateLazyLoader) GetCacheKey() string {
	return tll.cacheKey
}

// Invalidate clears the loaded data and cache
func (tll *TemplateLazyLoader) Invalidate() error {
	tll.mutex.Lock()
	defer tll.mutex.Unlock()

	tll.loaded = false
	tll.data = nil

	if tll.cache != nil {
		return tll.cache.Delete(tll.cacheKey)
	}

	return nil
}

// VersionLazyLoader provides lazy loading for version information
type VersionLazyLoader struct {
	versionManager interfaces.VersionManager
	cache          interfaces.CacheManager
	loaded         bool
	data           interface{}
	mutex          sync.RWMutex
	cacheKey       string
}

// NewVersionLazyLoader creates a new version lazy loader
func NewVersionLazyLoader(versionManager interfaces.VersionManager, cache interfaces.CacheManager) *VersionLazyLoader {
	return &VersionLazyLoader{
		versionManager: versionManager,
		cache:          cache,
		cacheKey:       "lazy:version:info",
	}
}

// Load loads the version information if not already loaded
func (vll *VersionLazyLoader) Load(ctx context.Context) (interface{}, error) {
	vll.mutex.Lock()
	defer vll.mutex.Unlock()

	if vll.loaded && vll.data != nil {
		return vll.data, nil
	}

	// Try cache first
	if vll.cache != nil {
		if cached, err := vll.cache.Get(vll.cacheKey); err == nil {
			vll.data = cached
			vll.loaded = true
			return vll.data, nil
		}
	}

	// Load from version manager
	versionInfo := vll.versionManager.GetCurrentVersion()

	vll.data = versionInfo
	vll.loaded = true

	// Cache the result
	if vll.cache != nil {
		vll.cache.Set(vll.cacheKey, versionInfo, 1*time.Hour)
	}

	return vll.data, nil
}

// IsLoaded returns whether the data has been loaded
func (vll *VersionLazyLoader) IsLoaded() bool {
	vll.mutex.RLock()
	defer vll.mutex.RUnlock()
	return vll.loaded
}

// GetCacheKey returns the cache key for this loader
func (vll *VersionLazyLoader) GetCacheKey() string {
	return vll.cacheKey
}

// ConfigLazyLoader provides lazy loading for configuration data
type ConfigLazyLoader struct {
	configManager interfaces.ConfigManager
	cache         interfaces.CacheManager
	loaded        bool
	data          interface{}
	mutex         sync.RWMutex
	cacheKey      string
	configPath    string
}

// NewConfigLazyLoader creates a new config lazy loader
func NewConfigLazyLoader(configManager interfaces.ConfigManager, cache interfaces.CacheManager, configPath string) *ConfigLazyLoader {
	return &ConfigLazyLoader{
		configManager: configManager,
		cache:         cache,
		configPath:    configPath,
		cacheKey:      fmt.Sprintf("lazy:config:%s", configPath),
	}
}

// Load loads the configuration if not already loaded
func (cll *ConfigLazyLoader) Load(ctx context.Context) (interface{}, error) {
	cll.mutex.Lock()
	defer cll.mutex.Unlock()

	if cll.loaded && cll.data != nil {
		return cll.data, nil
	}

	// Try cache first
	if cll.cache != nil {
		if cached, err := cll.cache.Get(cll.cacheKey); err == nil {
			cll.data = cached
			cll.loaded = true
			return cll.data, nil
		}
	}

	// Load configuration
	config, err := cll.configManager.LoadConfig(cll.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", cll.configPath, err)
	}

	cll.data = config
	cll.loaded = true

	// Cache the result
	if cll.cache != nil {
		cll.cache.Set(cll.cacheKey, config, 15*time.Minute)
	}

	return cll.data, nil
}

// IsLoaded returns whether the data has been loaded
func (cll *ConfigLazyLoader) IsLoaded() bool {
	cll.mutex.RLock()
	defer cll.mutex.RUnlock()
	return cll.loaded
}

// GetCacheKey returns the cache key for this loader
func (cll *ConfigLazyLoader) GetCacheKey() string {
	return cll.cacheKey
}

// AuditLazyLoader provides lazy loading for audit operations
type AuditLazyLoader struct {
	auditEngine interfaces.AuditEngine
	cache       interfaces.CacheManager
	loaded      bool
	data        interface{}
	mutex       sync.RWMutex
	cacheKey    string
	projectPath string
	auditType   string
}

// NewAuditLazyLoader creates a new audit lazy loader
func NewAuditLazyLoader(auditEngine interfaces.AuditEngine, cache interfaces.CacheManager, projectPath, auditType string) *AuditLazyLoader {
	return &AuditLazyLoader{
		auditEngine: auditEngine,
		cache:       cache,
		projectPath: projectPath,
		auditType:   auditType,
		cacheKey:    fmt.Sprintf("lazy:audit:%s:%s", auditType, projectPath),
	}
}

// Load loads the audit results if not already loaded
func (all *AuditLazyLoader) Load(ctx context.Context) (interface{}, error) {
	all.mutex.Lock()
	defer all.mutex.Unlock()

	if all.loaded && all.data != nil {
		return all.data, nil
	}

	// Try cache first
	if all.cache != nil {
		if cached, err := all.cache.Get(all.cacheKey); err == nil {
			all.data = cached
			all.loaded = true
			return all.data, nil
		}
	}

	// Perform audit based on type
	var result interface{}
	var err error

	switch all.auditType {
	case "security":
		result, err = all.auditEngine.AuditSecurity(all.projectPath)
	case "quality":
		result, err = all.auditEngine.AuditCodeQuality(all.projectPath)
	case "performance":
		result, err = all.auditEngine.AuditPerformance(all.projectPath)
	case "licenses":
		result, err = all.auditEngine.AuditLicenses(all.projectPath)
	default:
		return nil, fmt.Errorf("unknown audit type: %s", all.auditType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to perform %s audit: %w", all.auditType, err)
	}

	all.data = result
	all.loaded = true

	// Cache the result
	if all.cache != nil {
		all.cache.Set(all.cacheKey, result, 10*time.Minute)
	}

	return all.data, nil
}

// IsLoaded returns whether the data has been loaded
func (all *AuditLazyLoader) IsLoaded() bool {
	all.mutex.RLock()
	defer all.mutex.RUnlock()
	return all.loaded
}

// GetCacheKey returns the cache key for this loader
func (all *AuditLazyLoader) GetCacheKey() string {
	return all.cacheKey
}

// LazyLoaderManager manages multiple lazy loaders
type LazyLoaderManager struct {
	loaders map[string]LazyLoader
	mutex   sync.RWMutex
}

// NewLazyLoaderManager creates a new lazy loader manager
func NewLazyLoaderManager() *LazyLoaderManager {
	return &LazyLoaderManager{
		loaders: make(map[string]LazyLoader),
	}
}

// Register registers a lazy loader with a key
func (llm *LazyLoaderManager) Register(key string, loader LazyLoader) {
	llm.mutex.Lock()
	defer llm.mutex.Unlock()
	llm.loaders[key] = loader
}

// Get retrieves a lazy loader by key
func (llm *LazyLoaderManager) Get(key string) (LazyLoader, bool) {
	llm.mutex.RLock()
	defer llm.mutex.RUnlock()
	loader, exists := llm.loaders[key]
	return loader, exists
}

// LoadAll loads all registered lazy loaders
func (llm *LazyLoaderManager) LoadAll(ctx context.Context) error {
	llm.mutex.RLock()
	loaders := make(map[string]LazyLoader)
	for k, v := range llm.loaders {
		loaders[k] = v
	}
	llm.mutex.RUnlock()

	var errors []error
	for key, loader := range loaders {
		if !loader.IsLoaded() {
			if _, err := loader.Load(ctx); err != nil {
				errors = append(errors, fmt.Errorf("failed to load %s: %w", key, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to load some lazy loaders: %v", errors)
	}

	return nil
}

// InvalidateAll invalidates all registered lazy loaders
func (llm *LazyLoaderManager) InvalidateAll() error {
	llm.mutex.RLock()
	loaders := make(map[string]LazyLoader)
	for k, v := range llm.loaders {
		loaders[k] = v
	}
	llm.mutex.RUnlock()

	var errors []error
	for key, loader := range loaders {
		// Check if loader has Invalidate method
		if invalidator, ok := loader.(interface{ Invalidate() error }); ok {
			if err := invalidator.Invalidate(); err != nil {
				errors = append(errors, fmt.Errorf("failed to invalidate %s: %w", key, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to invalidate some lazy loaders: %v", errors)
	}

	return nil
}

// GetLoadedCount returns the number of loaded lazy loaders
func (llm *LazyLoaderManager) GetLoadedCount() int {
	llm.mutex.RLock()
	defer llm.mutex.RUnlock()

	count := 0
	for _, loader := range llm.loaders {
		if loader.IsLoaded() {
			count++
		}
	}

	return count
}

// GetTotalCount returns the total number of registered lazy loaders
func (llm *LazyLoaderManager) GetTotalCount() int {
	llm.mutex.RLock()
	defer llm.mutex.RUnlock()
	return len(llm.loaders)
}
