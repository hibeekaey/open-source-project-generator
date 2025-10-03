package validation

import (
	"os"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

func TestValidator_ValidateCache(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := interfaces.DefaultCacheConfig()
	config.Location = tempDir

	validator := NewValidator(tempDir, config)

	// Test with valid entries
	entries := map[string]*interfaces.CacheEntry{
		"test1": {
			Key:         "test1",
			Value:       "value1",
			Size:        6,
			CreatedAt:   time.Now().Add(-time.Hour),
			UpdatedAt:   time.Now().Add(-time.Hour),
			AccessedAt:  time.Now().Add(-time.Minute),
			AccessCount: 5,
			Metadata:    make(map[string]any),
		},
	}

	err = validator.ValidateCache(entries)
	if err != nil {
		t.Errorf("ValidateCache failed with valid entries: %v", err)
	}
}

func TestValidator_RepairEntries(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := interfaces.DefaultCacheConfig()
	validator := NewValidator(tempDir, config)

	// Create corrupted entries
	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"corrupted1": {
			Key:         "wrong_key", // Key mismatch
			Value:       "value1",
			Size:        -10,                // Negative size
			CreatedAt:   now.Add(time.Hour), // Future timestamp
			UpdatedAt:   now,
			AccessedAt:  now,
			AccessCount: -5,  // Negative access count
			Metadata:    nil, // Nil metadata
		},
		"expired": {
			Key:        "expired",
			Value:      "expired_value",
			Size:       12,
			CreatedAt:  now.Add(-2 * time.Hour),
			UpdatedAt:  now.Add(-2 * time.Hour),
			AccessedAt: now.Add(-2 * time.Hour),
			ExpiresAt:  &[]time.Time{now.Add(-time.Hour)}[0], // Expired
			Metadata:   make(map[string]any),
		},
	}

	repairedEntries := validator.RepairEntries(entries)

	// Check that corrupted entry was repaired
	if repaired, exists := repairedEntries["corrupted1"]; exists {
		if repaired.Key != "corrupted1" {
			t.Errorf("Expected key to be fixed to 'corrupted1', got %s", repaired.Key)
		}
		if repaired.Size < 0 {
			t.Errorf("Expected size to be fixed, still negative: %d", repaired.Size)
		}
		if repaired.AccessCount < 0 {
			t.Errorf("Expected access count to be fixed, still negative: %d", repaired.AccessCount)
		}
		if repaired.CreatedAt.After(now) {
			t.Errorf("Expected creation time to be fixed, still in future")
		}
		if repaired.Metadata == nil {
			t.Errorf("Expected metadata to be initialized")
		}
	} else {
		t.Errorf("Expected corrupted entry to be repaired and present")
	}

	// Check that expired entry was removed
	if _, exists := repairedEntries["expired"]; exists {
		t.Errorf("Expected expired entry to be removed")
	}
}

func TestValidator_CheckHealth(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := interfaces.DefaultCacheConfig()
	validator := NewValidator(tempDir, config)

	// Create test entries with some issues
	now := time.Now()
	entries := map[string]*interfaces.CacheEntry{
		"healthy": {
			Key:         "healthy",
			Value:       "value",
			Size:        5,
			CreatedAt:   now.Add(-time.Hour),
			UpdatedAt:   now.Add(-time.Hour),
			AccessedAt:  now.Add(-time.Minute),
			AccessCount: 1,
			Metadata:    make(map[string]any),
		},
		"expired": {
			Key:        "expired",
			Value:      "value",
			Size:       5,
			CreatedAt:  now.Add(-2 * time.Hour),
			UpdatedAt:  now.Add(-2 * time.Hour),
			AccessedAt: now.Add(-2 * time.Hour),
			ExpiresAt:  &[]time.Time{now.Add(-time.Hour)}[0], // Expired
			Metadata:   make(map[string]any),
		},
	}

	metrics := &interfaces.CacheMetrics{
		Hits:           50,
		Misses:         50,
		Gets:           100,
		CurrentSize:    10,
		CurrentEntries: 2,
	}

	report, err := validator.CheckHealth(entries, metrics)
	if err != nil {
		t.Fatalf("CheckHealth failed: %v", err)
	}

	if report.TotalEntries != 2 {
		t.Errorf("Expected 2 total entries, got %d", report.TotalEntries)
	}

	if report.ExpiredEntries != 1 {
		t.Errorf("Expected 1 expired entry, got %d", report.ExpiredEntries)
	}

	if report.OverallHealth == "healthy" && report.ExpiredEntries > 0 {
		// Should be degraded due to expired entries
		t.Errorf("Expected health to be degraded due to expired entries, got %s", report.OverallHealth)
	}
}

func TestValidator_ValidateConfiguration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	validator := NewValidator(tempDir, nil)

	// Test with nil config
	err = validator.validateConfiguration()
	if err == nil {
		t.Errorf("Expected error with nil config")
	}

	// Test with invalid config
	invalidConfig := &interfaces.CacheConfig{
		MaxSize:           -1,        // Invalid
		EvictionRatio:     1.5,       // Invalid
		EvictionPolicy:    "invalid", // Invalid
		CompressionLevel:  10,        // Invalid
		EnableCompression: true,
	}
	validator.SetConfig(invalidConfig)

	err = validator.validateConfiguration()
	if err == nil {
		t.Errorf("Expected error with invalid config")
	}

	// Test with valid config
	validConfig := interfaces.DefaultCacheConfig()
	validator.SetConfig(validConfig)

	err = validator.validateConfiguration()
	if err != nil {
		t.Errorf("ValidateConfiguration failed with valid config: %v", err)
	}
}
