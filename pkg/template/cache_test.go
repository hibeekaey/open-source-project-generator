package template

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// createTestTemplates creates sample templates for testing
func createTestTemplates() []*models.TemplateInfo {
	return []*models.TemplateInfo{
		{
			Name:        "go-gin",
			DisplayName: "Go Gin API",
			Description: "Go API server with Gin framework",
			Category:    "backend",
			Technology:  "Go",
			Version:     "1.0.0",
			Tags:        []string{"api", "backend", "go"},
			Source:      "embedded",
			Size:        1024,
			FileCount:   10,
		},
		{
			Name:        "nextjs-app",
			DisplayName: "Next.js App",
			Description: "Next.js application template",
			Category:    "frontend",
			Technology:  "Next.js",
			Version:     "1.0.0",
			Tags:        []string{"frontend", "react", "nextjs"},
			Source:      "embedded",
			Size:        2048,
			FileCount:   20,
		},
		{
			Name:        "android-kotlin",
			DisplayName: "Android Kotlin",
			Description: "Android app with Kotlin",
			Category:    "mobile",
			Technology:  "Kotlin",
			Version:     "1.0.0",
			Tags:        []string{"mobile", "android", "kotlin"},
			Source:      "embedded",
			Size:        4096,
			FileCount:   30,
		},
	}
}

// mockDiscoveryFunc creates a mock discovery function for testing
func mockDiscoveryFunc(templates []*models.TemplateInfo, shouldError bool) func() ([]*models.TemplateInfo, error) {
	return func() ([]*models.TemplateInfo, error) {
		if shouldError {
			return nil, fmt.Errorf("mock discovery error")
		}
		return templates, nil
	}
}

func TestNewTemplateCache(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)

	// Test with default config
	cache := NewTemplateCache(nil, discoveryFunc)
	if cache == nil {
		t.Fatal("Expected cache to be created, got nil")
	}

	if cache.cacheTTL != 5*time.Minute {
		t.Errorf("Expected default TTL of 5 minutes, got %v", cache.cacheTTL)
	}

	// Test with custom config
	config := &CacheConfig{
		TTL:               10 * time.Minute,
		MaxSize:           100,
		EnableAutoRefresh: false,
	}

	cache = NewTemplateCache(config, discoveryFunc)
	if cache.cacheTTL != 10*time.Minute {
		t.Errorf("Expected TTL of 10 minutes, got %v", cache.cacheTTL)
	}
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()
	if config == nil {
		t.Fatal("Expected config to be created, got nil")
	}

	if config.TTL != 5*time.Minute {
		t.Errorf("Expected default TTL of 5 minutes, got %v", config.TTL)
	}

	if config.MaxSize != 0 {
		t.Errorf("Expected default MaxSize of 0, got %d", config.MaxSize)
	}

	if config.EnableAutoRefresh {
		t.Error("Expected default EnableAutoRefresh to be false")
	}
}

func TestTemplateCache_SetAndGet(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test Set and Get
	cache.Set(templates)

	// Test getting existing template
	template, found := cache.Get("go-gin")
	if !found {
		t.Error("Expected to find template 'go-gin'")
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	if template.Name != "go-gin" {
		t.Errorf("Expected template name 'go-gin', got '%s'", template.Name)
	}

	// Test getting non-existent template
	_, found = cache.Get("non-existent")
	if found {
		t.Error("Expected not to find non-existent template")
	}

	// Test that returned template is a copy (modifications don't affect cache)
	originalName := template.Name
	template.Name = "modified"

	template2, _ := cache.Get("go-gin")
	if template2.Name != originalName {
		t.Error("Expected template in cache to be unmodified")
	}
}

func TestTemplateCache_GetAll(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test GetAll with empty cache
	_, found := cache.GetAll()
	if found {
		t.Error("Expected not to find templates in empty cache")
	}

	// Test GetAll with populated cache
	cache.Set(templates)
	allTemplates, found := cache.GetAll()
	if !found {
		t.Error("Expected to find templates in populated cache")
	}

	if len(allTemplates) != len(templates) {
		t.Errorf("Expected %d templates, got %d", len(templates), len(allTemplates))
	}

	// Verify all templates are present
	templateNames := make(map[string]bool)
	for _, tmpl := range allTemplates {
		templateNames[tmpl.Name] = true
	}

	for _, original := range templates {
		if !templateNames[original.Name] {
			t.Errorf("Expected to find template '%s' in results", original.Name)
		}
	}
}

func TestTemplateCache_Put(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test putting single template
	template := templates[0]
	cache.Put(template)

	// Verify template was stored
	retrieved, found := cache.Get(template.Name)
	if !found {
		t.Error("Expected to find put template")
	}

	if retrieved.Name != template.Name {
		t.Errorf("Expected template name '%s', got '%s'", template.Name, retrieved.Name)
	}

	// Test that stored template is a copy
	originalName := templates[0].Name
	template.Name = "modified"
	retrieved2, found2 := cache.Get(originalName)
	if !found2 {
		t.Error("Expected to find template after modification")
	}
	if retrieved2 == nil {
		t.Fatal("Expected retrieved template not to be nil")
	}
	if retrieved2.Name != originalName {
		t.Error("Expected template in cache to be unmodified")
	}
}

func TestTemplateCache_Remove(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	cache.Set(templates)

	// Test removing existing template
	removed := cache.Remove("go-gin")
	if !removed {
		t.Error("Expected template to be removed")
	}

	// Verify template was removed
	_, found := cache.Get("go-gin")
	if found {
		t.Error("Expected template to be removed from cache")
	}

	// Test removing non-existent template
	removed = cache.Remove("non-existent")
	if removed {
		t.Error("Expected false when removing non-existent template")
	}
}

func TestTemplateCache_Clear(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	cache.Set(templates)

	// Verify cache has templates
	if cache.Size() == 0 {
		t.Error("Expected cache to have templates before clear")
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	if cache.Size() != 0 {
		t.Error("Expected cache to be empty after clear")
	}

	// Verify templates are not found
	_, found := cache.Get("go-gin")
	if found {
		t.Error("Expected not to find templates after clear")
	}
}

func TestTemplateCache_Refresh(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test refresh with empty cache
	err := cache.Refresh()
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	// Verify templates were loaded
	if cache.Size() != len(templates) {
		t.Errorf("Expected %d templates after refresh, got %d", len(templates), cache.Size())
	}

	// Test refresh with error
	errorDiscoveryFunc := mockDiscoveryFunc(nil, true)
	cache2 := NewTemplateCache(nil, errorDiscoveryFunc)

	err = cache2.Refresh()
	if err == nil {
		t.Error("Expected error from refresh with failing discovery function")
	}

	// Test refresh with no discovery function
	cache3 := NewTemplateCache(nil, nil)
	err = cache3.Refresh()
	if err == nil {
		t.Error("Expected error from refresh with no discovery function")
	}
}

func TestTemplateCache_IsValid(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)

	// Test with short TTL for expiration testing
	config := &CacheConfig{
		TTL: 100 * time.Millisecond,
	}
	cache := NewTemplateCache(config, discoveryFunc)

	// Empty cache should not be valid
	if cache.IsValid() {
		t.Error("Expected empty cache to be invalid")
	}

	// Populated cache should be valid
	cache.Set(templates)
	if !cache.IsValid() {
		t.Error("Expected populated cache to be valid")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	if cache.IsValid() {
		t.Error("Expected expired cache to be invalid")
	}
}

func TestTemplateCache_IsExpired(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)

	// Test with short TTL for expiration testing
	config := &CacheConfig{
		TTL: 100 * time.Millisecond,
	}
	cache := NewTemplateCache(config, discoveryFunc)

	// New cache should be expired (no cache time set)
	if !cache.IsExpired() {
		t.Error("Expected new cache to be expired")
	}

	// Set templates and check expiration
	cache.Set(templates)
	if cache.IsExpired() {
		t.Error("Expected fresh cache not to be expired")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	if !cache.IsExpired() {
		t.Error("Expected cache to be expired after TTL")
	}
}

func TestTemplateCache_Size(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Empty cache
	if cache.Size() != 0 {
		t.Errorf("Expected empty cache size 0, got %d", cache.Size())
	}

	// Populated cache
	cache.Set(templates)
	if cache.Size() != len(templates) {
		t.Errorf("Expected cache size %d, got %d", len(templates), cache.Size())
	}

	// After removing one
	cache.Remove("go-gin")
	if cache.Size() != len(templates)-1 {
		t.Errorf("Expected cache size %d, got %d", len(templates)-1, cache.Size())
	}
}

func TestTemplateCache_LastUpdated(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Empty cache should have zero time
	if !cache.LastUpdated().IsZero() {
		t.Error("Expected empty cache to have zero last updated time")
	}

	// Set templates and check time
	before := time.Now()
	cache.Set(templates)
	after := time.Now()

	lastUpdated := cache.LastUpdated()
	if lastUpdated.Before(before) || lastUpdated.After(after) {
		t.Error("Expected last updated time to be between before and after times")
	}
}

func TestTemplateCache_GetCacheStats(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test stats for empty cache
	stats := cache.GetCacheStats()
	if stats.Size != 0 {
		t.Errorf("Expected stats size 0, got %d", stats.Size)
	}
	if stats.IsValid {
		t.Error("Expected empty cache stats to show invalid")
	}
	if !stats.IsExpired {
		t.Error("Expected empty cache stats to show expired")
	}

	// Test stats for populated cache
	cache.Set(templates)
	stats = cache.GetCacheStats()
	if stats.Size != len(templates) {
		t.Errorf("Expected stats size %d, got %d", len(templates), stats.Size)
	}
	if !stats.IsValid {
		t.Error("Expected populated cache stats to show valid")
	}
	if stats.IsExpired {
		t.Error("Expected fresh cache stats to show not expired")
	}
}

func TestTemplateCache_InvalidateTemplate(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	cache.Set(templates)

	// Invalidate specific template
	cache.InvalidateTemplate("go-gin")

	// Verify template was removed
	_, found := cache.Get("go-gin")
	if found {
		t.Error("Expected invalidated template to be removed")
	}

	// Verify other templates still exist
	_, found = cache.Get("nextjs-app")
	if !found {
		t.Error("Expected other templates to remain after invalidation")
	}
}

func TestTemplateCache_InvalidateAll(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	cache.Set(templates)

	// Invalidate all templates
	cache.InvalidateAll()

	// Verify cache is empty
	if cache.Size() != 0 {
		t.Error("Expected cache to be empty after invalidating all")
	}

	// Verify no templates are found
	_, found := cache.Get("go-gin")
	if found {
		t.Error("Expected no templates after invalidating all")
	}
}

func TestTemplateCache_SetTTL(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test setting new TTL
	newTTL := 10 * time.Minute
	cache.SetTTL(newTTL)

	if cache.GetTTL() != newTTL {
		t.Errorf("Expected TTL %v, got %v", newTTL, cache.GetTTL())
	}
}

func TestTemplateCache_GetCachedTemplateNames(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Empty cache
	names := cache.GetCachedTemplateNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 names from empty cache, got %d", len(names))
	}

	// Populated cache
	cache.Set(templates)
	names = cache.GetCachedTemplateNames()
	if len(names) != len(templates) {
		t.Errorf("Expected %d names, got %d", len(templates), len(names))
	}

	// Verify all template names are present
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	for _, template := range templates {
		if !nameSet[template.Name] {
			t.Errorf("Expected to find template name '%s'", template.Name)
		}
	}
}

func TestTemplateCache_HasTemplate(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)

	// Test with short TTL for expiration testing
	config := &CacheConfig{
		TTL: 100 * time.Millisecond,
	}
	cache := NewTemplateCache(config, discoveryFunc)

	// Empty cache
	if cache.HasTemplate("go-gin") {
		t.Error("Expected empty cache not to have template")
	}

	// Populated cache
	cache.Set(templates)
	if !cache.HasTemplate("go-gin") {
		t.Error("Expected populated cache to have template")
	}

	if cache.HasTemplate("non-existent") {
		t.Error("Expected cache not to have non-existent template")
	}

	// Expired cache
	time.Sleep(150 * time.Millisecond)
	if cache.HasTemplate("go-gin") {
		t.Error("Expected expired cache not to have template")
	}
}

func TestTemplateCache_GetOrRefresh(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test with empty cache (should refresh)
	template, err := cache.GetOrRefresh("go-gin")
	if err != nil {
		t.Fatalf("GetOrRefresh failed: %v", err)
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	if template.Name != "go-gin" {
		t.Errorf("Expected template name 'go-gin', got '%s'", template.Name)
	}

	// Test with populated cache (should get from cache)
	template2, err := cache.GetOrRefresh("nextjs-app")
	if err != nil {
		t.Fatalf("GetOrRefresh failed: %v", err)
	}

	if template2.Name != "nextjs-app" {
		t.Errorf("Expected template name 'nextjs-app', got '%s'", template2.Name)
	}

	// Test with non-existent template
	_, err = cache.GetOrRefresh("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}

	// Test with failing discovery function
	errorDiscoveryFunc := mockDiscoveryFunc(nil, true)
	cache2 := NewTemplateCache(nil, errorDiscoveryFunc)

	_, err = cache2.GetOrRefresh("go-gin")
	if err == nil {
		t.Error("Expected error from failing discovery function")
	}
}

func TestTemplateCache_GetAllOrRefresh(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test with empty cache (should refresh)
	allTemplates, err := cache.GetAllOrRefresh()
	if err != nil {
		t.Fatalf("GetAllOrRefresh failed: %v", err)
	}

	if len(allTemplates) != len(templates) {
		t.Errorf("Expected %d templates, got %d", len(templates), len(allTemplates))
	}

	// Test with populated cache (should get from cache)
	allTemplates2, err := cache.GetAllOrRefresh()
	if err != nil {
		t.Fatalf("GetAllOrRefresh failed: %v", err)
	}

	if len(allTemplates2) != len(templates) {
		t.Errorf("Expected %d templates, got %d", len(templates), len(allTemplates2))
	}

	// Test with failing discovery function
	errorDiscoveryFunc := mockDiscoveryFunc(nil, true)
	cache2 := NewTemplateCache(nil, errorDiscoveryFunc)

	_, err = cache2.GetAllOrRefresh()
	if err == nil {
		t.Error("Expected error from failing discovery function")
	}
}

func TestTemplateCache_ConvertToInterfaceTemplateInfos(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	// Test conversion
	interfaceTemplates := cache.ConvertToInterfaceTemplateInfos(templates)

	if len(interfaceTemplates) != len(templates) {
		t.Errorf("Expected %d interface templates, got %d", len(templates), len(interfaceTemplates))
	}

	// Verify conversion
	for i, original := range templates {
		converted := interfaceTemplates[i]
		if converted.Name != original.Name {
			t.Errorf("Expected name '%s', got '%s'", original.Name, converted.Name)
		}
		if converted.Category != original.Category {
			t.Errorf("Expected category '%s', got '%s'", original.Category, converted.Category)
		}
		if converted.Technology != original.Technology {
			t.Errorf("Expected technology '%s', got '%s'", original.Technology, converted.Technology)
		}
	}
}

func TestTemplateCache_ConcurrentAccess(t *testing.T) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	cache.Set(templates)

	// Test concurrent reads and writes
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_, _ = cache.Get("go-gin")
				_ = cache.Size()
				_ = cache.IsValid()
				_ = cache.GetCachedTemplateNames()
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				testTemplate := &models.TemplateInfo{
					Name:        fmt.Sprintf("test-template-%d-%d", id, j),
					DisplayName: "Test Template",
					Category:    "test",
					Technology:  "Test",
					Version:     "1.0.0",
				}
				cache.Put(testTemplate)
				cache.Remove(testTemplate.Name)
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional
	template, found := cache.Get("go-gin")
	if !found {
		t.Error("Expected to find original template after concurrent operations")
	}
	if template.Name != "go-gin" {
		t.Error("Expected original template to be intact after concurrent operations")
	}
}

func TestCacheStats_String(t *testing.T) {
	stats := &CacheStats{
		Size:        5,
		LastUpdated: time.Now(),
		TTL:         5 * time.Minute,
		IsValid:     true,
		IsExpired:   false,
	}

	str := stats.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Should contain key information
	if !containsString(str, "Size: 5") {
		t.Error("Expected string to contain size information")
	}
	if !containsString(str, "IsValid: true") {
		t.Error("Expected string to contain validity information")
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubStr(s, substr)))
}

func containsSubStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkTemplateCache_Get(b *testing.B) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)
	cache.Set(templates)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("go-gin")
	}
}

func BenchmarkTemplateCache_Set(b *testing.B) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(templates)
	}
}

func BenchmarkTemplateCache_GetAll(b *testing.B) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)
	cache.Set(templates)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.GetAll()
	}
}

func BenchmarkTemplateCache_ConcurrentReads(b *testing.B) {
	templates := createTestTemplates()
	discoveryFunc := mockDiscoveryFunc(templates, false)
	cache := NewTemplateCache(nil, discoveryFunc)
	cache.Set(templates)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get("go-gin")
		}
	})
}
