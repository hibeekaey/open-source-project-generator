package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

// TemplateMetadata represents metadata for a template
type TemplateMetadata struct {
	Name         string              `yaml:"name" json:"name"`
	Description  string              `yaml:"description" json:"description"`
	Version      string              `yaml:"version" json:"version"`
	Dependencies []string            `yaml:"dependencies" json:"dependencies"`
	Variables    []TemplateVariable  `yaml:"variables" json:"variables"`
	Conditions   []TemplateCondition `yaml:"conditions" json:"conditions"`
	Files        []TemplateFile      `yaml:"files" json:"files"`
	Directories  []string            `yaml:"directories" json:"directories"`
}

// TemplateVariable represents a template variable definition
type TemplateVariable struct {
	Name        string      `yaml:"name" json:"name"`
	Type        string      `yaml:"type" json:"type"`
	Default     interface{} `yaml:"default" json:"default"`
	Description string      `yaml:"description" json:"description"`
	Required    bool        `yaml:"required" json:"required"`
	Options     []string    `yaml:"options" json:"options"`
}

// TemplateCondition represents a condition for template rendering
type TemplateCondition struct {
	Name      string      `yaml:"name" json:"name"`
	Component string      `yaml:"component" json:"component"`
	Value     interface{} `yaml:"value" json:"value"`
	Operator  string      `yaml:"operator" json:"operator"` // eq, ne, gt, lt, contains, etc.
}

// TemplateFile represents a file in the template
type TemplateFile struct {
	Source      string              `yaml:"source" json:"source"`
	Destination string              `yaml:"destination" json:"destination"`
	Conditions  []TemplateCondition `yaml:"conditions" json:"conditions"`
	Binary      bool                `yaml:"binary" json:"binary"`
	Executable  bool                `yaml:"executable" json:"executable"`
}

// MetadataParser handles template metadata parsing and validation
type MetadataParser struct{}

// NewMetadataParser creates a new metadata parser
func NewMetadataParser() *MetadataParser {
	return &MetadataParser{}
}

// ParseMetadata parses template metadata from a file
func (p *MetadataParser) ParseMetadata(metadataPath string) (*TemplateMetadata, error) {
	// Read metadata file
	content, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file %s: %w", metadataPath, err)
	}

	// Parse YAML content
	var metadata TemplateMetadata
	if err := yaml.Unmarshal(content, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata YAML: %w", err)
	}

	// Validate metadata
	if err := p.ValidateMetadata(&metadata); err != nil {
		return nil, fmt.Errorf("metadata validation failed: %w", err)
	}

	return &metadata, nil
}

// ParseMetadataFromDirectory looks for and parses metadata from a template directory
func (p *MetadataParser) ParseMetadataFromDirectory(templateDir string) (*TemplateMetadata, error) {
	// Look for common metadata file names
	metadataFiles := []string{
		"template.yaml",
		"template.yml",
		"metadata.yaml",
		"metadata.yml",
		".template.yaml",
		".template.yml",
	}

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templateDir, filename)
		if _, err := os.Stat(metadataPath); err == nil {
			return p.ParseMetadata(metadataPath)
		}
	}

	// If no metadata file found, return default metadata
	return p.createDefaultMetadata(templateDir), nil
}

// ValidateMetadata validates template metadata
func (p *MetadataParser) ValidateMetadata(metadata *TemplateMetadata) error {
	if metadata.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if metadata.Version == "" {
		return fmt.Errorf("template version is required")
	}

	// Validate variables
	for i, variable := range metadata.Variables {
		if variable.Name == "" {
			return fmt.Errorf("variable at index %d is missing name", i)
		}

		if variable.Type == "" {
			metadata.Variables[i].Type = "string" // Default type
		}

		// Validate variable type
		validTypes := []string{"string", "int", "bool", "array", "object"}
		isValidType := false
		for _, validType := range validTypes {
			if metadata.Variables[i].Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return fmt.Errorf("invalid variable type '%s' for variable '%s'", metadata.Variables[i].Type, variable.Name)
		}
	}

	// Validate conditions
	for i, condition := range metadata.Conditions {
		if condition.Name == "" {
			return fmt.Errorf("condition at index %d is missing name", i)
		}

		if condition.Operator == "" {
			metadata.Conditions[i].Operator = "eq" // Default operator
		}

		// Validate operator
		validOperators := []string{"eq", "ne", "gt", "lt", "ge", "le", "contains", "startswith", "endswith"}
		isValidOperator := false
		for _, validOp := range validOperators {
			if metadata.Conditions[i].Operator == validOp {
				isValidOperator = true
				break
			}
		}
		if !isValidOperator {
			return fmt.Errorf("invalid operator '%s' for condition '%s'", metadata.Conditions[i].Operator, condition.Name)
		}
	}

	// Validate files
	for i, file := range metadata.Files {
		if file.Source == "" {
			return fmt.Errorf("file at index %d is missing source", i)
		}

		if file.Destination == "" {
			// Default destination is the same as source
			metadata.Files[i].Destination = file.Source
		}
	}

	return nil
}

// createDefaultMetadata creates default metadata for a template directory
func (p *MetadataParser) createDefaultMetadata(templateDir string) *TemplateMetadata {
	dirName := filepath.Base(templateDir)

	return &TemplateMetadata{
		Name:        dirName,
		Description: fmt.Sprintf("Template for %s", dirName),
		Version:     "1.0.0",
		Variables:   []TemplateVariable{},
		Conditions:  []TemplateCondition{},
		Files:       []TemplateFile{},
		Directories: []string{},
	}
}

// EvaluateConditions evaluates template conditions against project configuration
func (p *MetadataParser) EvaluateConditions(conditions []TemplateCondition, config interface{}) (bool, error) {
	if len(conditions) == 0 {
		return true, nil // No conditions means always include
	}

	// All conditions must be true (AND logic)
	for _, condition := range conditions {
		result, err := p.evaluateCondition(condition, config)
		if err != nil {
			return false, fmt.Errorf("failed to evaluate condition '%s': %w", condition.Name, err)
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

// evaluateCondition evaluates a single condition
func (p *MetadataParser) evaluateCondition(condition TemplateCondition, config interface{}) (bool, error) {
	// Handle component-based conditions
	if condition.Component != "" {
		return p.evaluateComponentCondition(condition, config)
	}

	// Handle simple value comparisons
	switch condition.Operator {
	case "eq":
		return fmt.Sprintf("%v", condition.Value) == fmt.Sprintf("%v", config), nil
	case "ne":
		return fmt.Sprintf("%v", condition.Value) != fmt.Sprintf("%v", config), nil
	case "contains":
		configStr := fmt.Sprintf("%v", config)
		valueStr := fmt.Sprintf("%v", condition.Value)
		return contains(configStr, valueStr), nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", condition.Operator)
	}
}

// evaluateComponentCondition evaluates component-specific conditions
func (p *MetadataParser) evaluateComponentCondition(condition TemplateCondition, config interface{}) (bool, error) {
	// Type assert to ProjectConfig
	projectConfig, ok := config.(*models.ProjectConfig)
	if !ok {
		return false, fmt.Errorf("config is not a ProjectConfig")
	}

	var actualValue bool

	// Extract component value based on component type
	switch condition.Component {
	case "frontend":
		actualValue = projectConfig.Components.Frontend.MainApp ||
			projectConfig.Components.Frontend.Home ||
			projectConfig.Components.Frontend.Admin
	case "backend":
		actualValue = projectConfig.Components.Backend.API
	case "mobile":
		actualValue = projectConfig.Components.Mobile.Android ||
			projectConfig.Components.Mobile.IOS
	case "infrastructure":
		actualValue = projectConfig.Components.Infrastructure.Docker ||
			projectConfig.Components.Infrastructure.Kubernetes ||
			projectConfig.Components.Infrastructure.Terraform
	case "always":
		actualValue = true
	default:
		return false, fmt.Errorf("unknown component type: %s", condition.Component)
	}

	expectedValue, ok := condition.Value.(bool)
	if !ok {
		return false, fmt.Errorf("condition value must be boolean for component conditions")
	}

	// Apply operator
	switch condition.Operator {
	case "eq":
		return actualValue == expectedValue, nil
	case "ne":
		return actualValue != expectedValue, nil
	default:
		return false, fmt.Errorf("unsupported operator for component condition: %s", condition.Operator)
	}
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
