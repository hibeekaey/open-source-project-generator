package template

import (
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TestCacheIntegrationWithManager tests that the cache works correctly with the template manager
func TestCacheIntegrationWithManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine).(*Manager)

	// Test that cache is initialized
	if manager.cache == nil {
		t.Fatal("Expected template cache to be initialized")
	}

	// Test cache operations through manager interface
	t.Run("CacheTemplate", func(t *testing.T) {
		err := manager.CacheTemplate("go-gin")
		if err != nil {
			t.Fatalf("CacheTemplate failed: %v", err)
		}
	})

	t.Run("GetCachedTemplates", func(t *testing.T) {
		templates, err := manager.GetCachedTemplates()
		if err != nil {
			t.Fatalf("GetCachedTemplates failed: %v", err)
		}

		if len(templates) == 0 {
			t.Error("Expected cached templates to be returned")
		}

		// Verify we get interface types
		for _, tmpl := range templates {
			if tmpl.Name == "" {
				t.Error("Expected template to have a name")
			}
			if tmpl.Category == "" {
				t.Error("Expected template to have a category")
			}
		}
	})

	t.Run("ClearTemplateCache", func(t *testing.T) {
		err := manager.ClearTemplateCache()
		if err != nil {
			t.Fatalf("ClearTemplateCache failed: %v", err)
		}

		// Verify cache is empty
		if manager.cache.Size() != 0 {
			t.Error("Expected cache to be empty after clear")
		}
	})

	t.Run("RefreshTemplateCache", func(t *testing.T) {
		err := manager.RefreshTemplateCache()
		if err != nil {
			t.Fatalf("RefreshTemplateCache failed: %v", err)
		}

		// Verify cache has templates after refresh
		if manager.cache.Size() == 0 {
			t.Error("Expected cache to have templates after refresh")
		}
	})
}

// TestCachePerformanceWithManager tests cache performance improvements
func TestCachePerformanceWithManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// First call should populate cache
	start := time.Now()
	templates1, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("First ListTemplates failed: %v", err)
	}
	firstCallDuration := time.Since(start)

	// Second call should use cache and be faster
	start = time.Now()
	templates2, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("Second ListTemplates failed: %v", err)
	}
	secondCallDuration := time.Since(start)

	// Verify results are the same
	if len(templates1) != len(templates2) {
		t.Errorf("Expected same number of templates, got %d and %d", len(templates1), len(templates2))
	}

	// Second call should be faster (cache hit)
	// Note: This is a rough check since embedded templates are fast anyway
	t.Logf("First call: %v, Second call: %v", firstCallDuration, secondCallDuration)

	// Verify cache is being used by checking cache stats
	managerImpl := manager.(*Manager)
	stats := managerImpl.cache.GetCacheStats()
	if !stats.IsValid {
		t.Error("Expected cache to be valid after ListTemplates calls")
	}
	if stats.Size == 0 {
		t.Error("Expected cache to have templates")
	}
}

// TestCacheExpirationWithManager tests cache expiration behavior
func TestCacheExpirationWithManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine).(*Manager)

	// Set a very short TTL for testing
	manager.cache.SetTTL(50 * time.Millisecond)

	// Populate cache
	_, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	// Verify cache is valid
	if !manager.cache.IsValid() {
		t.Error("Expected cache to be valid initially")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Verify cache is expired
	if manager.cache.IsValid() {
		t.Error("Expected cache to be expired after TTL")
	}

	// Next call should refresh cache
	_, err = manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("ListTemplates after expiration failed: %v", err)
	}

	// Verify cache is valid again
	if !manager.cache.IsValid() {
		t.Error("Expected cache to be valid after refresh")
	}
}

// TestCacheThreadSafetyWithManager tests concurrent access through manager
func TestCacheThreadSafetyWithManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Run concurrent operations
	done := make(chan bool, 10)

	// Concurrent readers
	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				_, _ = manager.ListTemplates(interfaces.TemplateFilter{})
				_, _ = manager.GetTemplateInfo("go-gin")
				_, _ = manager.GetCachedTemplates()
			}
		}()
	}

	// Concurrent cache operations
	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				_ = manager.CacheTemplate("go-gin")
				_ = manager.RefreshTemplateCache()
				if j%3 == 0 {
					_ = manager.ClearTemplateCache()
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify manager is still functional
	templates, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("ListTemplates after concurrent access failed: %v", err)
	}

	// Should have templates (cache may have been cleared and refreshed)
	if len(templates) == 0 {
		// Try to refresh if empty
		err = manager.RefreshTemplateCache()
		if err != nil {
			t.Fatalf("RefreshTemplateCache failed: %v", err)
		}

		templates, err = manager.ListTemplates(interfaces.TemplateFilter{})
		if err != nil {
			t.Fatalf("ListTemplates after refresh failed: %v", err)
		}
	}

	t.Logf("Final template count after concurrent access: %d", len(templates))
}

// TestCacheInvalidationWithManager tests cache invalidation scenarios
func TestCacheInvalidationWithManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine).(*Manager)

	// Populate cache
	_, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	initialSize := manager.cache.Size()
	if initialSize == 0 {
		t.Fatal("Expected cache to have templates")
	}

	// Test invalidating specific template
	manager.cache.InvalidateTemplate("go-gin")

	// Cache should still have other templates
	if manager.cache.Size() >= initialSize {
		t.Error("Expected cache size to decrease after invalidating template")
	}

	// Test invalidating all templates
	manager.cache.InvalidateAll()

	// Cache should be empty
	if manager.cache.Size() != 0 {
		t.Error("Expected cache to be empty after invalidating all")
	}

	// Next operation should repopulate cache
	_, err = manager.GetTemplateInfo("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateInfo after invalidation failed: %v", err)
	}

	// Cache should have templates again
	if manager.cache.Size() == 0 {
		t.Error("Expected cache to be repopulated after GetTemplateInfo")
	}
}
