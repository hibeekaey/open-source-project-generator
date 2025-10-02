package validators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ConfigValidator handles validation of configuration data
type ConfigValidator struct{}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateConfigurationName validates a configuration name
func (cv *ConfigValidator) ValidateConfigurationName(name string) error {
	if name == "" {
		return interfaces.NewValidationError("name", name, "Configuration name is required", "required").
			WithSuggestions("Enter a descriptive name for your configuration")
	}

	if len(name) < 2 {
		return interfaces.NewValidationError("name", name, "Configuration name must be at least 2 characters long", "min_length").
			WithSuggestions("Use a longer, more descriptive name")
	}

	if len(name) > 50 {
		return interfaces.NewValidationError("name", name, "Configuration name must be at most 50 characters long", "max_length").
			WithSuggestions("Use a shorter, more concise name")
	}

	// Check for invalid characters
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return interfaces.NewValidationError("name", name, "Configuration name contains invalid characters", "invalid_chars").
			WithSuggestions("Use only alphanumeric characters, spaces, hyphens, and underscores")
	}

	return nil
}

// ValidateExportPath validates an export file path
func (cv *ConfigValidator) ValidateExportPath(path string) error {
	if path == "" {
		return fmt.Errorf("export path is required")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check if path is valid
	if !filepath.IsAbs(cleanPath) && !strings.HasPrefix(cleanPath, ".") {
		return fmt.Errorf("export path must be absolute or relative to current directory")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext != ".yaml" && ext != ".yml" && ext != ".json" {
		return fmt.Errorf("export file must have .yaml, .yml, or .json extension")
	}

	return nil
}

// ValidateImportPath validates an import file path
func (cv *ConfigValidator) ValidateImportPath(path string) error {
	if path == "" {
		return fmt.Errorf("import path is required")
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	// For security, require absolute paths for imports
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("import file path must be absolute: %s", path)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext != ".yaml" && ext != ".yml" && ext != ".json" {
		return fmt.Errorf("import file must have .yaml, .yml, or .json extension")
	}

	return nil
}

// ValidateConfigurationDescription validates a configuration description
func (cv *ConfigValidator) ValidateConfigurationDescription(description string) error {
	if len(description) > 200 {
		return interfaces.NewValidationError("description", description, "Description must be at most 200 characters long", "max_length").
			WithSuggestions("Use a shorter, more concise description")
	}

	return nil
}

// ValidateConfigurationTags validates configuration tags
func (cv *ConfigValidator) ValidateConfigurationTags(tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("maximum of 10 tags allowed")
	}

	for i, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("tag %d is empty", i+1)
		}

		if len(tag) > 20 {
			return fmt.Errorf("tag '%s' is too long (maximum 20 characters)", tag)
		}

		// Check for invalid characters in tags
		if strings.ContainsAny(tag, "/\\:*?\"<>|,") {
			return fmt.Errorf("tag '%s' contains invalid characters", tag)
		}
	}

	return nil
}

// ValidateTagsString validates a comma-separated tags string
func (cv *ConfigValidator) ValidateTagsString(tagsStr string) error {
	if tagsStr == "" {
		return nil // Empty tags string is valid
	}

	if len(tagsStr) > 100 {
		return interfaces.NewValidationError("tags", tagsStr, "Tags string must be at most 100 characters long", "max_length").
			WithSuggestions("Use shorter tag names or fewer tags")
	}

	// Parse and validate individual tags
	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return cv.ValidateConfigurationTags(tags)
}
