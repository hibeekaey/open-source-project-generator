// Package cache provides tool cache management and export/import functionality
// #nosec G304 - Cache file operations use paths from application configuration
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ExportFormat defines the portable cache format
type ExportFormat struct {
	Version    string                        `json:"version"`
	ExportedAt time.Time                     `json:"exported_at"`
	Platform   string                        `json:"platform"`
	Entries    map[string]*models.CachedTool `json:"entries"`
}

// CacheExporter handles cache import/export
type CacheExporter struct {
	logger *logger.Logger
}

// NewCacheExporter creates a new cache exporter instance
func NewCacheExporter(log *logger.Logger) *CacheExporter {
	return &CacheExporter{
		logger: log,
	}
}

// Export exports cache to JSON format
func (ce *CacheExporter) Export(cache *orchestrator.ToolCache, outputPath string) error {
	if cache == nil {
		return fmt.Errorf("cache is nil")
	}

	if ce.logger != nil {
		ce.logger.Info(fmt.Sprintf("Exporting cache to %s", outputPath))
	}

	// Create export format
	exportData := &ExportFormat{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Platform:   runtime.GOOS,
		Entries:    make(map[string]*models.CachedTool),
	}

	// We need to access cache entries
	// Since ToolCache doesn't expose its internal map, we'll need to save and read it
	// First, ensure cache is saved
	if err := cache.Save(); err != nil {
		return fmt.Errorf("failed to save cache before export: %w", err)
	}

	// Read the cache file
	cacheFile := cache.GetCacheFile()
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	// Unmarshal cache entries
	entries := make(map[string]*models.CachedTool)
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	exportData.Entries = entries

	// Marshal export data
	exportJSON, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(outputPath, exportJSON, 0600); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	if ce.logger != nil {
		ce.logger.Info(fmt.Sprintf("Cache exported successfully: %d entries", len(entries)))
	}

	return nil
}

// Import imports cache from JSON format
func (ce *CacheExporter) Import(cache *orchestrator.ToolCache, inputPath string) error {
	if cache == nil {
		return fmt.Errorf("cache is nil")
	}

	if ce.logger != nil {
		ce.logger.Info(fmt.Sprintf("Importing cache from %s", inputPath))
	}

	// Read import file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Unmarshal import data
	var importData ExportFormat
	if err := json.Unmarshal(data, &importData); err != nil {
		return fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	// Validate import data
	if err := ce.validateImportData(&importData); err != nil {
		return fmt.Errorf("import validation failed: %w", err)
	}

	// Warn if platform mismatch
	if importData.Platform != runtime.GOOS && ce.logger != nil {
		ce.logger.Warn(fmt.Sprintf("Platform mismatch: exported on %s, importing on %s",
			importData.Platform, runtime.GOOS))
	}

	// Import entries into cache
	imported := 0
	for toolName, entry := range importData.Entries {
		// Update the cached timestamp to now to avoid immediate expiration
		cache.Set(toolName, entry.Available, entry.Version)
		imported++
	}

	// Save cache to disk
	if err := cache.Save(); err != nil {
		return fmt.Errorf("failed to save cache after import: %w", err)
	}

	if ce.logger != nil {
		ce.logger.Info(fmt.Sprintf("Cache imported successfully: %d entries", imported))
	}

	return nil
}

// validateImportData validates imported cache data
func (ce *CacheExporter) validateImportData(data *ExportFormat) error {
	if data == nil {
		return fmt.Errorf("import data is nil")
	}

	// Check version
	if data.Version == "" {
		return fmt.Errorf("missing version in import data")
	}

	// Check if version is supported
	if data.Version != "1.0" {
		return fmt.Errorf("unsupported import version: %s", data.Version)
	}

	// Check if exported timestamp is reasonable
	if data.ExportedAt.IsZero() {
		return fmt.Errorf("missing export timestamp")
	}

	// Warn if export is very old (more than 30 days)
	if time.Since(data.ExportedAt) > 30*24*time.Hour && ce.logger != nil {
		ce.logger.Warn(fmt.Sprintf("Import data is %d days old",
			int(time.Since(data.ExportedAt).Hours()/24)))
	}

	// Check if entries exist
	if data.Entries == nil {
		return fmt.Errorf("missing entries in import data")
	}

	// Validate each entry
	for toolName, entry := range data.Entries {
		if entry == nil {
			return fmt.Errorf("nil entry for tool '%s'", toolName)
		}

		// Check if CachedAt is reasonable
		if entry.CachedAt.IsZero() {
			return fmt.Errorf("missing cache timestamp for tool '%s'", toolName)
		}

		// Check if TTL is reasonable
		if entry.TTL < 0 {
			return fmt.Errorf("negative TTL for tool '%s'", toolName)
		}
	}

	return nil
}

// GetExportStats returns statistics about an export file
func (ce *CacheExporter) GetExportStats(exportPath string) (*ExportFormat, error) {
	// Read export file
	data, err := os.ReadFile(exportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read export file: %w", err)
	}

	// Unmarshal export data
	var exportData ExportFormat
	if err := json.Unmarshal(data, &exportData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal export data: %w", err)
	}

	return &exportData, nil
}
