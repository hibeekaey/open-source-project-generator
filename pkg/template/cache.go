package template

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sync"
	texttemplate "text/template"
	"time"
)

// TemplateCache provides caching for parsed templates to improve performance
type TemplateCache struct {
	cache   map[string]*CachedTemplate
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
}

// CachedTemplate represents a cached template with metadata
type CachedTemplate struct {
	Template  *texttemplate.Template
	Hash      string
	CreatedAt time.Time
	AccessAt  time.Time
	HitCount  int64
}

// NewTemplateCache creates a new template cache with specified max size and TTL
func NewTemplateCache(maxSize int, ttl time.Duration) *TemplateCache {
	return &TemplateCache{
		cache:   make(map[string]*CachedTemplate),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves a template from cache if available and not expired
func (tc *TemplateCache) Get(key string, contentHash string) (*texttemplate.Template, bool) {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	cached, exists := tc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if template has expired
	if time.Since(cached.CreatedAt) > tc.ttl {
		// Don't remove here to avoid write lock, let cleanup handle it
		return nil, false
	}

	// Check if content has changed
	if cached.Hash != contentHash {
		return nil, false
	}

	// Update access time and hit count (need write lock for this)
	tc.mutex.RUnlock()
	tc.mutex.Lock()
	cached.AccessAt = time.Now()
	cached.HitCount++
	tc.mutex.Unlock()
	tc.mutex.RLock()

	return cached.Template, true
}

// Put stores a template in the cache
func (tc *TemplateCache) Put(key string, template *texttemplate.Template, contentHash string) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// Check if we need to evict entries
	if len(tc.cache) >= tc.maxSize {
		tc.evictLRU()
	}

	now := time.Now()
	tc.cache[key] = &CachedTemplate{
		Template:  template,
		Hash:      contentHash,
		CreatedAt: now,
		AccessAt:  now,
		HitCount:  0,
	}
}

// evictLRU removes the least recently used entry
func (tc *TemplateCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range tc.cache {
		if oldestKey == "" || cached.AccessAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.AccessAt
		}
	}

	if oldestKey != "" {
		delete(tc.cache, oldestKey)
	}
}

// Cleanup removes expired entries from the cache
func (tc *TemplateCache) Cleanup() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	now := time.Now()
	for key, cached := range tc.cache {
		if now.Sub(cached.CreatedAt) > tc.ttl {
			delete(tc.cache, key)
		}
	}
}

// Stats returns cache statistics
func (tc *TemplateCache) Stats() CacheStats {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	var totalHits int64
	for _, cached := range tc.cache {
		totalHits += cached.HitCount
	}

	return CacheStats{
		Size:      len(tc.cache),
		MaxSize:   tc.maxSize,
		TotalHits: totalHits,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	Size      int
	MaxSize   int
	TotalHits int64
}

// Clear removes all entries from the cache
func (tc *TemplateCache) Clear() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	tc.cache = make(map[string]*CachedTemplate)
}

// HashContent creates a hash of template content for cache validation
func HashContent(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}

// RenderCache provides caching for rendered template output
type RenderCache struct {
	cache   map[string]*CachedRender
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
}

// CachedRender represents cached rendered output
type CachedRender struct {
	Content   []byte
	Hash      string
	CreatedAt time.Time
	AccessAt  time.Time
}

// NewRenderCache creates a new render cache
func NewRenderCache(maxSize int, ttl time.Duration) *RenderCache {
	return &RenderCache{
		cache:   make(map[string]*CachedRender),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves rendered content from cache
func (rc *RenderCache) Get(templateKey string, configHash string) ([]byte, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", templateKey, configHash)
	cached, exists := rc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if render has expired
	if time.Since(cached.CreatedAt) > rc.ttl {
		return nil, false
	}

	// Update access time
	rc.mutex.RUnlock()
	rc.mutex.Lock()
	cached.AccessAt = time.Now()
	rc.mutex.Unlock()
	rc.mutex.RLock()

	return cached.Content, true
}

// Put stores rendered content in cache
func (rc *RenderCache) Put(templateKey string, configHash string, content []byte) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", templateKey, configHash)

	// Check if we need to evict entries
	if len(rc.cache) >= rc.maxSize {
		rc.evictLRU()
	}

	now := time.Now()
	rc.cache[key] = &CachedRender{
		Content:   content,
		Hash:      configHash,
		CreatedAt: now,
		AccessAt:  now,
	}
}

// evictLRU removes the least recently used render entry
func (rc *RenderCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range rc.cache {
		if oldestKey == "" || cached.AccessAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.AccessAt
		}
	}

	if oldestKey != "" {
		delete(rc.cache, oldestKey)
	}
}

// Cleanup removes expired render entries
func (rc *RenderCache) Cleanup() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	now := time.Now()
	for key, cached := range rc.cache {
		if now.Sub(cached.CreatedAt) > rc.ttl {
			delete(rc.cache, key)
		}
	}
}

// HashConfig creates a hash of configuration for render caching
func HashConfig(config interface{}) (string, error) {
	var buf bytes.Buffer

	// Use a simple string representation for hashing
	// In production, you might want to use a more sophisticated approach
	configStr := fmt.Sprintf("%+v", config)
	buf.WriteString(configStr)

	hash := sha256.Sum256(buf.Bytes())
	return fmt.Sprintf("%x", hash), nil
}

// Clear removes all entries from the render cache
func (rc *RenderCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.cache = make(map[string]*CachedRender)
}
