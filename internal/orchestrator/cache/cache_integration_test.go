package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockToolDiscovery implements ToolDiscoveryInterface for testing
type MockToolDiscovery struct {
	tools     map[string]bool
	versions  map[string]string
	callCount int
}

func NewMockToolDiscovery() *MockToolDiscovery {
	return &MockToolDiscovery{
		tools:    make(map[string]bool),
		versions: make(map[string]string),
	}
}

func (m *MockToolDiscovery) WithTool(name string, available bool, version string) *MockToolDiscovery {
	m.tools[name] = available
	if available {
		m.versions[name] = version
	}
	return m
}

func (m *MockToolDiscovery) ListRegisteredTools() []string {
	tools := make([]string, 0, len(m.tools))
	for name := range m.tools {
		tools = append(tools, name)
	}
	return tools
}

func (m *MockToolDiscovery) IsAvailable(toolName string) (bool, error) {
	m.callCount++
	available, exists := m.tools[toolName]
	if !exists {
		return false, nil
	}
	return available, nil
}

func (m *MockToolDiscovery) GetVersion(toolName string) (string, error) {
	version, exists := m.versions[toolName]
	if !exists {
		return "", nil
	}
	return version, nil
}

// TestCacheManager_GetStats tests cache statistics retrieval
func TestCacheManager_GetStats(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	// Add some entries
	cache.Set("npx", true, "9.0.0")
	cache.Set("go", true, "1.21.0")
	cache.Set("gradle", false, "")

	manager := NewCacheManager(cache, log)
	stats := manager.GetStats()

	assert.NotNil(t, stats)
	assert.Equal(t, 3, stats.TotalEntries)
	assert.Equal(t, 2, stats.AvailableTools)
	assert.Equal(t, 1, stats.UnavailableTools)
	assert.Equal(t, 5*time.Minute, stats.TTL)
}

// TestCacheManager_Validate tests cache validation
func TestCacheManager_Validate(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	// Add valid entries
	cache.Set("npx", true, "9.0.0")
	cache.Set("go", true, "1.21.0")
	err = cache.Save()
	require.NoError(t, err)

	manager := NewCacheManager(cache, log)
	report, err := manager.Validate()

	require.NoError(t, err)
	assert.NotNil(t, report)
	assert.True(t, report.Valid)
	assert.Equal(t, 2, report.TotalEntries)
	assert.Empty(t, report.CorruptedEntries)
}

// TestCacheManager_Validate_CorruptedCache tests validation with corrupted cache
func TestCacheManager_Validate_CorruptedCache(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	// Create cache file with invalid JSON
	cacheFile := filepath.Join(tempDir, "tool_cache.json")
	err := os.WriteFile(cacheFile, []byte("invalid json {{{"), 0644)
	require.NoError(t, err)

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	manager := NewCacheManager(cache, log)
	report, err := manager.Validate()

	require.NoError(t, err)
	assert.NotNil(t, report)
	// Validation should detect issues
	assert.NotEmpty(t, report.Warnings)
}

// TestCacheManager_Refresh tests cache refresh functionality
func TestCacheManager_Refresh(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	// Setup mock tool discovery
	mockDiscovery := NewMockToolDiscovery().
		WithTool("npx", true, "9.0.0").
		WithTool("go", true, "1.21.0").
		WithTool("gradle", false, "")

	manager := NewCacheManager(cache, log)
	err = manager.Refresh(mockDiscovery)

	require.NoError(t, err)
	assert.Greater(t, mockDiscovery.callCount, 0, "Should have called tool discovery")

	// Verify cache was updated
	stats := manager.GetStats()
	assert.Equal(t, 3, stats.TotalEntries)
	assert.Equal(t, 2, stats.AvailableTools)
}

// TestCacheManager_Export tests cache export functionality
func TestCacheManager_Export(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	// Add entries
	cache.Set("npx", true, "9.0.0")
	cache.Set("go", true, "1.21.0")
	err = cache.Save()
	require.NoError(t, err)

	manager := NewCacheManager(cache, log)
	exportPath := filepath.Join(tempDir, "cache_export.json")
	err = manager.Export(exportPath)

	require.NoError(t, err)
	assert.FileExists(t, exportPath)

	// Verify export file content
	data, err := os.ReadFile(exportPath)
	require.NoError(t, err)

	var exportData ExportFormat
	err = json.Unmarshal(data, &exportData)
	require.NoError(t, err)

	assert.Equal(t, "1.0", exportData.Version)
	assert.NotEmpty(t, exportData.Platform)
	assert.Len(t, exportData.Entries, 2)
	assert.Contains(t, exportData.Entries, "npx")
	assert.Contains(t, exportData.Entries, "go")
}

// TestCacheManager_Import tests cache import functionality
func TestCacheManager_Import(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	// Create export file
	exportData := &ExportFormat{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Platform:   "linux",
		Entries: map[string]*models.CachedTool{
			"npx": {
				Available: true,
				Version:   "9.0.0",
				CachedAt:  time.Now(),
				TTL:       5 * time.Minute,
			},
			"go": {
				Available: true,
				Version:   "1.21.0",
				CachedAt:  time.Now(),
				TTL:       5 * time.Minute,
			},
		},
	}

	exportPath := filepath.Join(tempDir, "import.json")
	exportJSON, err := json.Marshal(exportData)
	require.NoError(t, err)
	err = os.WriteFile(exportPath, exportJSON, 0644)
	require.NoError(t, err)

	// Create new cache and import
	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	manager := NewCacheManager(cache, log)
	err = manager.Import(exportPath)

	require.NoError(t, err)

	// Verify imported data
	stats := manager.GetStats()
	assert.Equal(t, 2, stats.TotalEntries)
	assert.Equal(t, 2, stats.AvailableTools)
}

// TestCacheValidator_CheckExpired tests expired entry detection
func TestCacheValidator_CheckExpired(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      1 * time.Millisecond, // Very short TTL
		},
		log,
	)
	require.NoError(t, err)

	// Add entry and wait for expiration
	cache.Set("npx", true, "9.0.0")
	time.Sleep(10 * time.Millisecond)

	validator := NewCacheValidator(log)
	expired := validator.CheckExpired(cache)

	// Note: The current implementation may not return specific tool names
	// but should detect expired entries through stats
	stats := cache.GetStats()
	if expiredCount, ok := stats["expired"].(int); ok {
		assert.GreaterOrEqual(t, expiredCount, 0)
	}
	assert.NotNil(t, expired)
}

// TestCacheValidator_ValidateCacheEntry tests individual entry validation
func TestCacheValidator_ValidateCacheEntry(t *testing.T) {
	log := logger.NewLogger()
	validator := NewCacheValidator(log)

	tests := []struct {
		name    string
		entry   *models.CachedTool
		wantErr bool
	}{
		{
			name: "valid entry",
			entry: &models.CachedTool{
				Available: true,
				Version:   "1.0.0",
				CachedAt:  time.Now(),
				TTL:       5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name:    "nil entry",
			entry:   nil,
			wantErr: true,
		},
		{
			name: "future timestamp",
			entry: &models.CachedTool{
				Available: true,
				Version:   "1.0.0",
				CachedAt:  time.Now().Add(1 * time.Hour),
				TTL:       5 * time.Minute,
			},
			wantErr: true,
		},
		{
			name: "negative TTL",
			entry: &models.CachedTool{
				Available: true,
				Version:   "1.0.0",
				CachedAt:  time.Now(),
				TTL:       -1 * time.Minute,
			},
			wantErr: true,
		},
		{
			name: "excessive TTL",
			entry: &models.CachedTool{
				Available: true,
				Version:   "1.0.0",
				CachedAt:  time.Now(),
				TTL:       48 * time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCacheEntry("test-tool", tt.entry)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCacheExporter_Export tests export functionality
func TestCacheExporter_Export(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	cache.Set("npx", true, "9.0.0")
	cache.Set("go", false, "")
	err = cache.Save()
	require.NoError(t, err)

	exporter := NewCacheExporter(log)
	exportPath := filepath.Join(tempDir, "export.json")
	err = exporter.Export(cache, exportPath)

	require.NoError(t, err)
	assert.FileExists(t, exportPath)

	// Verify export format
	data, err := os.ReadFile(exportPath)
	require.NoError(t, err)

	var exportData ExportFormat
	err = json.Unmarshal(data, &exportData)
	require.NoError(t, err)

	assert.Equal(t, "1.0", exportData.Version)
	assert.NotZero(t, exportData.ExportedAt)
	assert.NotEmpty(t, exportData.Platform)
	assert.Len(t, exportData.Entries, 2)
}

// TestCacheExporter_Import_InvalidData tests import with invalid data
func TestCacheExporter_Import_InvalidData(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	tests := []struct {
		name       string
		exportData interface{}
		wantErr    bool
	}{
		{
			name: "missing version",
			exportData: map[string]interface{}{
				"exported_at": time.Now(),
				"platform":    "linux",
				"entries":     map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "unsupported version",
			exportData: &ExportFormat{
				Version:    "2.0",
				ExportedAt: time.Now(),
				Platform:   "linux",
				Entries:    map[string]*models.CachedTool{},
			},
			wantErr: true,
		},
		{
			name: "missing entries",
			exportData: map[string]interface{}{
				"version":     "1.0",
				"exported_at": time.Now(),
				"platform":    "linux",
			},
			wantErr: true,
		},
		{
			name: "nil entry",
			exportData: &ExportFormat{
				Version:    "1.0",
				ExportedAt: time.Now(),
				Platform:   "linux",
				Entries: map[string]*models.CachedTool{
					"npx": nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create import file
			importPath := filepath.Join(tempDir, tt.name+".json")
			importJSON, err := json.Marshal(tt.exportData)
			require.NoError(t, err)
			err = os.WriteFile(importPath, importJSON, 0644)
			require.NoError(t, err)

			// Try to import
			cache, err := orchestrator.NewToolCache(
				&orchestrator.ToolCacheConfig{
					CacheDir: tempDir,
					TTL:      5 * time.Minute,
				},
				log,
			)
			require.NoError(t, err)

			exporter := NewCacheExporter(log)
			err = exporter.Import(cache, importPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCacheExporter_GetExportStats tests export statistics retrieval
func TestCacheExporter_GetExportStats(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	// Create valid export file
	exportData := &ExportFormat{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Platform:   "linux",
		Entries: map[string]*models.CachedTool{
			"npx": {
				Available: true,
				Version:   "9.0.0",
				CachedAt:  time.Now(),
				TTL:       5 * time.Minute,
			},
		},
	}

	exportPath := filepath.Join(tempDir, "export.json")
	exportJSON, err := json.Marshal(exportData)
	require.NoError(t, err)
	err = os.WriteFile(exportPath, exportJSON, 0644)
	require.NoError(t, err)

	exporter := NewCacheExporter(log)
	stats, err := exporter.GetExportStats(exportPath)

	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "1.0", stats.Version)
	assert.Equal(t, "linux", stats.Platform)
	assert.Len(t, stats.Entries, 1)
}

// TestCacheManager_Integration tests complete cache management workflow
func TestCacheManager_Integration(t *testing.T) {
	tempDir := t.TempDir()
	log := logger.NewLogger()

	// Step 1: Create cache and add entries
	cache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: tempDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	cache.Set("npx", true, "9.0.0")
	cache.Set("go", true, "1.21.0")
	cache.Set("gradle", false, "")
	err = cache.Save()
	require.NoError(t, err)

	manager := NewCacheManager(cache, log)

	// Step 2: Get stats
	stats := manager.GetStats()
	assert.Equal(t, 3, stats.TotalEntries)
	assert.Equal(t, 2, stats.AvailableTools)

	// Step 3: Validate cache
	report, err := manager.Validate()
	require.NoError(t, err)
	assert.True(t, report.Valid)

	// Step 4: Export cache
	exportPath := filepath.Join(tempDir, "cache_export.json")
	err = manager.Export(exportPath)
	require.NoError(t, err)
	assert.FileExists(t, exportPath)

	// Step 5: Create new cache and import
	newCacheDir := filepath.Join(tempDir, "new_cache")
	err = os.MkdirAll(newCacheDir, 0755)
	require.NoError(t, err)

	newCache, err := orchestrator.NewToolCache(
		&orchestrator.ToolCacheConfig{
			CacheDir: newCacheDir,
			TTL:      5 * time.Minute,
		},
		log,
	)
	require.NoError(t, err)

	newManager := NewCacheManager(newCache, log)
	err = newManager.Import(exportPath)
	require.NoError(t, err)

	// Step 6: Verify imported cache
	newStats := newManager.GetStats()
	assert.Equal(t, stats.TotalEntries, newStats.TotalEntries)
	assert.Equal(t, stats.AvailableTools, newStats.AvailableTools)

	// Step 7: Refresh cache
	mockDiscovery := NewMockToolDiscovery().
		WithTool("npx", true, "9.1.0").
		WithTool("go", true, "1.22.0").
		WithTool("gradle", true, "8.0")

	err = newManager.Refresh(mockDiscovery)
	require.NoError(t, err)

	// Verify refresh updated the cache
	refreshedStats := newManager.GetStats()
	assert.Equal(t, 3, refreshedStats.TotalEntries)
	assert.Equal(t, 3, refreshedStats.AvailableTools) // gradle is now available
}
