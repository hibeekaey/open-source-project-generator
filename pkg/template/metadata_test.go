package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	tempDir := t.TempDir()
	metadataPath := filepath.Join(tempDir, "template.yaml")

	metadataContent := `name: test-template
description: A test template
version: 1.0.0
dependencies:
  - base
  - common
variables:
  - name: project_name
    type: string
    required: true
    description: The name of the project
  - name: enable_auth
    type: bool
    default: false
    description: Enable authentication
conditions:
  - name: has_frontend
    component: frontend
    value: true
    operator: eq
files:
  - source: src/main.go.tmpl
    destination: src/main.go
  - source: assets/logo.png
    destination: assets/logo.png
    binary: true
directories:
  - src
  - assets`

	err := os.WriteFile(metadataPath, []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata file: %v", err)
	}

	parser := NewMetadataParser()
	metadata, err := parser.ParseMetadata(metadataPath)
	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	// Verify basic fields
	if metadata.Name != "test-template" {
		t.Errorf("Expected name 'test-template', got '%s'", metadata.Name)
	}
	if metadata.Description != "A test template" {
		t.Errorf("Expected description 'A test template', got '%s'", metadata.Description)
	}
	if metadata.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", metadata.Version)
	}

	// Verify dependencies
	if len(metadata.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(metadata.Dependencies))
	}
	if metadata.Dependencies[0] != "base" || metadata.Dependencies[1] != "common" {
		t.Errorf("Unexpected dependencies: %v", metadata.Dependencies)
	}

	// Verify variables
	if len(metadata.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(metadata.Variables))
	}

	projectNameVar := metadata.Variables[0]
	if projectNameVar.Name != "project_name" {
		t.Errorf("Expected variable name 'project_name', got '%s'", projectNameVar.Name)
	}
	if projectNameVar.Type != "string" {
		t.Errorf("Expected variable type 'string', got '%s'", projectNameVar.Type)
	}
	if !projectNameVar.Required {
		t.Error("Expected project_name variable to be required")
	}

	enableAuthVar := metadata.Variables[1]
	if enableAuthVar.Name != "enable_auth" {
		t.Errorf("Expected variable name 'enable_auth', got '%s'", enableAuthVar.Name)
	}
	if enableAuthVar.Type != "bool" {
		t.Errorf("Expected variable type 'bool', got '%s'", enableAuthVar.Type)
	}
	if enableAuthVar.Default != false {
		t.Errorf("Expected default value false, got %v", enableAuthVar.Default)
	}

	// Verify conditions
	if len(metadata.Conditions) != 1 {
		t.Errorf("Expected 1 condition, got %d", len(metadata.Conditions))
	}

	condition := metadata.Conditions[0]
	if condition.Name != "has_frontend" {
		t.Errorf("Expected condition name 'has_frontend', got '%s'", condition.Name)
	}
	if condition.Component != "frontend" {
		t.Errorf("Expected condition component 'frontend', got '%s'", condition.Component)
	}
	if condition.Operator != "eq" {
		t.Errorf("Expected condition operator 'eq', got '%s'", condition.Operator)
	}

	// Verify files
	if len(metadata.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(metadata.Files))
	}

	templateFile := metadata.Files[0]
	if templateFile.Source != "src/main.go.tmpl" {
		t.Errorf("Expected source 'src/main.go.tmpl', got '%s'", templateFile.Source)
	}
	if templateFile.Destination != "src/main.go" {
		t.Errorf("Expected destination 'src/main.go', got '%s'", templateFile.Destination)
	}

	binaryFile := metadata.Files[1]
	if !binaryFile.Binary {
		t.Error("Expected binary file to have Binary=true")
	}

	// Verify directories
	if len(metadata.Directories) != 2 {
		t.Errorf("Expected 2 directories, got %d", len(metadata.Directories))
	}
}

func TestParseMetadataFromDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Test with template.yaml
	metadataContent := `name: dir-template
description: Template from directory
version: 2.0.0`

	metadataPath := filepath.Join(tempDir, "template.yaml")
	err := os.WriteFile(metadataPath, []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata file: %v", err)
	}

	parser := NewMetadataParser()
	metadata, err := parser.ParseMetadataFromDirectory(tempDir)
	if err != nil {
		t.Fatalf("ParseMetadataFromDirectory failed: %v", err)
	}

	if metadata.Name != "dir-template" {
		t.Errorf("Expected name 'dir-template', got '%s'", metadata.Name)
	}
	if metadata.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", metadata.Version)
	}
}

func TestParseMetadataFromDirectoryNoFile(t *testing.T) {
	tempDir := t.TempDir()

	parser := NewMetadataParser()
	metadata, err := parser.ParseMetadataFromDirectory(tempDir)
	if err != nil {
		t.Fatalf("ParseMetadataFromDirectory failed: %v", err)
	}

	// Should return default metadata
	expectedName := filepath.Base(tempDir)
	if metadata.Name != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, metadata.Name)
	}
	if metadata.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got '%s'", metadata.Version)
	}
}

func TestValidateMetadata(t *testing.T) {
	parser := NewMetadataParser()

	// Test valid metadata
	validMetadata := &TemplateMetadata{
		Name:    "valid-template",
		Version: "1.0.0",
		Variables: []TemplateVariable{
			{
				Name: "test_var",
				Type: "string",
			},
		},
		Conditions: []TemplateCondition{
			{
				Name:     "test_condition",
				Operator: "eq",
			},
		},
		Files: []TemplateFile{
			{
				Source: "test.tmpl",
			},
		},
	}

	err := parser.ValidateMetadata(validMetadata)
	if err != nil {
		t.Errorf("ValidateMetadata failed for valid metadata: %v", err)
	}

	// Test missing name
	invalidMetadata := &TemplateMetadata{
		Version: "1.0.0",
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for missing name")
	}

	// Test missing version
	invalidMetadata = &TemplateMetadata{
		Name: "test",
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for missing version")
	}

	// Test invalid variable type
	invalidMetadata = &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Variables: []TemplateVariable{
			{
				Name: "test_var",
				Type: "invalid_type",
			},
		},
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for invalid variable type")
	}

	// Test invalid operator
	invalidMetadata = &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Conditions: []TemplateCondition{
			{
				Name:     "test_condition",
				Operator: "invalid_op",
			},
		},
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for invalid operator")
	}

	// Test missing variable name
	invalidMetadata = &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Variables: []TemplateVariable{
			{
				Type: "string",
			},
		},
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for missing variable name")
	}

	// Test missing condition name
	invalidMetadata = &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Conditions: []TemplateCondition{
			{
				Operator: "eq",
			},
		},
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for missing condition name")
	}

	// Test missing file source
	invalidMetadata = &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Files: []TemplateFile{
			{
				Destination: "output.txt",
			},
		},
	}
	err = parser.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("Expected validation error for missing file source")
	}
}

func TestValidateMetadataDefaults(t *testing.T) {
	parser := NewMetadataParser()

	// Test that defaults are applied during validation
	metadata := &TemplateMetadata{
		Name:    "test",
		Version: "1.0.0",
		Variables: []TemplateVariable{
			{
				Name: "test_var",
				// Type should default to "string"
			},
		},
		Conditions: []TemplateCondition{
			{
				Name: "test_condition",
				// Operator should default to "eq"
			},
		},
		Files: []TemplateFile{
			{
				Source: "test.tmpl",
				// Destination should default to source
			},
		},
	}

	err := parser.ValidateMetadata(metadata)
	if err != nil {
		t.Errorf("ValidateMetadata failed: %v", err)
	}

	// Check that defaults were applied
	if metadata.Variables[0].Type != "string" {
		t.Errorf("Expected default type 'string', got '%s'", metadata.Variables[0].Type)
	}
	if metadata.Conditions[0].Operator != "eq" {
		t.Errorf("Expected default operator 'eq', got '%s'", metadata.Conditions[0].Operator)
	}
	if metadata.Files[0].Destination != "test.tmpl" {
		t.Errorf("Expected default destination 'test.tmpl', got '%s'", metadata.Files[0].Destination)
	}
}

func TestEvaluateConditions(t *testing.T) {
	parser := NewMetadataParser()

	// Test empty conditions (should return true)
	result, err := parser.EvaluateConditions([]TemplateCondition{}, "test")
	if err != nil {
		t.Errorf("EvaluateConditions failed for empty conditions: %v", err)
	}
	if !result {
		t.Error("Expected true for empty conditions")
	}

	// Test single condition that matches
	conditions := []TemplateCondition{
		{
			Name:     "test_eq",
			Operator: "eq",
			Value:    "test",
		},
	}
	result, err = parser.EvaluateConditions(conditions, "test")
	if err != nil {
		t.Errorf("EvaluateConditions failed: %v", err)
	}
	if !result {
		t.Error("Expected true for matching condition")
	}

	// Test single condition that doesn't match
	result, err = parser.EvaluateConditions(conditions, "other")
	if err != nil {
		t.Errorf("EvaluateConditions failed: %v", err)
	}
	if result {
		t.Error("Expected false for non-matching condition")
	}

	// Test multiple conditions (all must match)
	conditions = []TemplateCondition{
		{
			Name:     "test_eq",
			Operator: "eq",
			Value:    "test",
		},
		{
			Name:     "test_ne",
			Operator: "ne",
			Value:    "other",
		},
	}
	result, err = parser.EvaluateConditions(conditions, "test")
	if err != nil {
		t.Errorf("EvaluateConditions failed: %v", err)
	}
	if !result {
		t.Error("Expected true when all conditions match")
	}

	// Test multiple conditions where one doesn't match
	conditions[1].Value = "test" // Now both conditions check for "test"
	result, err = parser.EvaluateConditions(conditions, "test")
	if err != nil {
		t.Errorf("EvaluateConditions failed: %v", err)
	}
	if result {
		t.Error("Expected false when not all conditions match")
	}
}

func TestEvaluateCondition(t *testing.T) {
	parser := NewMetadataParser()

	tests := []struct {
		name      string
		condition TemplateCondition
		config    interface{}
		expected  bool
		shouldErr bool
	}{
		{
			name: "eq_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "eq",
				Value:    "hello",
			},
			config:   "hello",
			expected: true,
		},
		{
			name: "eq_no_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "eq",
				Value:    "hello",
			},
			config:   "world",
			expected: false,
		},
		{
			name: "ne_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "ne",
				Value:    "hello",
			},
			config:   "world",
			expected: true,
		},
		{
			name: "ne_no_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "ne",
				Value:    "hello",
			},
			config:   "hello",
			expected: false,
		},
		{
			name: "contains_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "contains",
				Value:    "ell",
			},
			config:   "hello",
			expected: true,
		},
		{
			name: "contains_no_match",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "contains",
				Value:    "xyz",
			},
			config:   "hello",
			expected: false,
		},
		{
			name: "unsupported_operator",
			condition: TemplateCondition{
				Name:     "test",
				Operator: "invalid",
				Value:    "hello",
			},
			config:    "hello",
			expected:  false,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.evaluateCondition(tt.condition, tt.config)

			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCreateDefaultMetadata(t *testing.T) {
	parser := NewMetadataParser()

	tempDir := t.TempDir()
	metadata := parser.createDefaultMetadata(tempDir)

	expectedName := filepath.Base(tempDir)
	if metadata.Name != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, metadata.Name)
	}

	if metadata.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", metadata.Version)
	}

	expectedDescription := "Template for " + expectedName
	if metadata.Description != expectedDescription {
		t.Errorf("Expected description '%s', got '%s'", expectedDescription, metadata.Description)
	}

	if len(metadata.Variables) != 0 {
		t.Errorf("Expected empty variables, got %d", len(metadata.Variables))
	}

	if len(metadata.Conditions) != 0 {
		t.Errorf("Expected empty conditions, got %d", len(metadata.Conditions))
	}
}

func TestParseMetadataInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	metadataPath := filepath.Join(tempDir, "template.yaml")

	// Invalid YAML content
	invalidYAML := `name: test
description: [unclosed array
version: 1.0.0`

	err := os.WriteFile(metadataPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid metadata file: %v", err)
	}

	parser := NewMetadataParser()
	_, err = parser.ParseMetadata(metadataPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestParseMetadataNonExistentFile(t *testing.T) {
	parser := NewMetadataParser()
	_, err := parser.ParseMetadata("/non/existent/file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
