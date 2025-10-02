package formatters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileFormatter handles file operations for configuration management
type FileFormatter struct{}

// NewFileFormatter creates a new file formatter
func NewFileFormatter() *FileFormatter {
	return &FileFormatter{}
}

// WriteExportFile writes configuration data to an export file
func (ff *FileFormatter) WriteExportFile(path string, data []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ReadImportFile reads configuration data from an import file
func (ff *FileFormatter) ReadImportFile(path string) ([]byte, error) {
	// Validate and clean path to prevent directory traversal
	path = filepath.Clean(path)

	// Ensure path is absolute to prevent traversal attacks
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("import file path must be absolute: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// DetermineFormatFromPath determines the file format from the file path
func (ff *FileFormatter) DetermineFormatFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "yaml" // Default to YAML
	}
}

// GenerateDefaultFilename generates a default filename for export
func (ff *FileFormatter) GenerateDefaultFilename(configName, format string) string {
	// Sanitize config name for filename
	sanitized := strings.ReplaceAll(configName, " ", "_")
	sanitized = strings.ToLower(sanitized)

	// Remove any characters that might be problematic in filenames
	sanitized = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, sanitized)

	return fmt.Sprintf("%s.%s", sanitized, format)
}

// ValidateFilePath validates that a file path is safe and accessible
func (ff *FileFormatter) ValidateFilePath(path string) error {
	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path contains directory traversal: %s", path)
	}

	// Check if the directory exists or can be created
	dir := filepath.Dir(cleanPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Try to create the directory
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
	}

	return nil
}
