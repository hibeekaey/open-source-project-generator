package validation

import (
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewSchemaManager(t *testing.T) {
	sm := NewSchemaManager()

	assert.NotNil(t, sm)
	assert.NotNil(t, sm.schemas)
	assert.NotNil(t, sm.validationRules)

	// Check that default schemas are loaded
	schemas := sm.ListSchemas()
	assert.Contains(t, schemas, "package.json")
	assert.Contains(t, schemas, "tsconfig.json")
	assert.Contains(t, schemas, ".eslintrc.json")
}

func TestSchemaManager_GetSchema(t *testing.T) {
	sm := NewSchemaManager()

	// Test existing schema
	schema, exists := sm.GetSchema("package.json")
	assert.True(t, exists)
	assert.NotNil(t, schema)
	assert.Equal(t, "Package.json Schema", schema.Title)

	// Test non-existing schema
	_, exists = sm.GetSchema("nonexistent.json")
	assert.False(t, exists)
}

func TestSchemaManager_AddRemoveSchema(t *testing.T) {
	sm := NewSchemaManager()

	testSchema := &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Test Schema",
		Description: "Test schema for testing",
		Required:    []string{"test"},
		Properties: map[string]interfaces.PropertySchema{
			"test": {
				Type:        "string",
				Description: "Test property",
			},
		},
	}

	// Add schema
	sm.AddSchema("test.json", testSchema)

	// Verify it was added
	schema, exists := sm.GetSchema("test.json")
	assert.True(t, exists)
	assert.Equal(t, "Test Schema", schema.Title)

	// Remove schema
	sm.RemoveSchema("test.json")

	// Verify it was removed
	_, exists = sm.GetSchema("test.json")
	assert.False(t, exists)
}

func TestSchemaManager_ValidateAgainstSchema(t *testing.T) {
	sm := NewSchemaManager()

	schema := &interfaces.ConfigSchema{
		Required: []string{"name", "version"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:      "string",
				MinLength: &[]int{1}[0],
				MaxLength: &[]int{50}[0],
			},
			"version": {
				Type:    "string",
				Pattern: `^\d+\.\d+\.\d+$`,
			},
		},
	}

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test valid data
	validData := map[string]interface{}{
		"name":    "test-package",
		"version": "1.0.0",
	}

	err := sm.ValidateAgainstSchema(validData, schema, result)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)
}
func TestSchemaManager_ValidateAgainstSchema_MissingRequired(t *testing.T) {
	sm := NewSchemaManager()

	schema := &interfaces.ConfigSchema{
		Required: []string{"name", "version"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type: "string",
			},
			"version": {
				Type: "string",
			},
		},
	}

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test missing required field
	invalidData := map[string]interface{}{
		"name": "test-package",
		// missing version
	}

	err := sm.ValidateAgainstSchema(invalidData, schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
	assert.Equal(t, 1, result.Summary.MissingRequired)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "version", result.Errors[0].Field)
	assert.Equal(t, "missing_required", result.Errors[0].Type)
}

func TestSchemaManager_ValidatePropertyAgainstSchema_String(t *testing.T) {
	sm := NewSchemaManager()

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test string length validation
	schema := interfaces.PropertySchema{
		Type:      "string",
		MinLength: &[]int{5}[0],
		MaxLength: &[]int{10}[0],
	}

	// Test too short
	err := sm.validatePropertyAgainstSchema("test", "abc", schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)

	// Reset result
	result.Valid = true
	result.Errors = []interfaces.ConfigValidationError{}
	result.Summary.ErrorCount = 0

	// Test too long
	err = sm.validatePropertyAgainstSchema("test", "this is too long", schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
}

func TestSchemaManager_ValidatePropertyAgainstSchema_Pattern(t *testing.T) {
	sm := NewSchemaManager()

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test pattern validation
	schema := interfaces.PropertySchema{
		Type:    "string",
		Pattern: `^\d+\.\d+\.\d+$`, // semantic version pattern
	}

	// Test valid pattern
	err := sm.validatePropertyAgainstSchema("version", "1.0.0", schema, result)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)

	// Test invalid pattern
	err = sm.validatePropertyAgainstSchema("version", "invalid-version", schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
}

func TestSchemaManager_ValidatePropertyAgainstSchema_Enum(t *testing.T) {
	sm := NewSchemaManager()

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test enum validation
	schema := interfaces.PropertySchema{
		Type: "string",
		Enum: []string{"MIT", "Apache-2.0", "GPL-3.0"},
	}

	// Test valid enum value
	err := sm.validatePropertyAgainstSchema("license", "MIT", schema, result)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)

	// Test invalid enum value
	err = sm.validatePropertyAgainstSchema("license", "InvalidLicense", schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
}

func TestSchemaManager_ValidatePropertyAgainstSchema_Number(t *testing.T) {
	sm := NewSchemaManager()

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test number range validation
	schema := interfaces.PropertySchema{
		Type:    "number",
		Minimum: &[]float64{0}[0],
		Maximum: &[]float64{100}[0],
	}

	// Test valid number
	err := sm.validatePropertyAgainstSchema("port", float64(80), schema, result)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)

	// Test number too small
	err = sm.validatePropertyAgainstSchema("port", float64(-1), schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)

	// Reset result
	result.Valid = true
	result.Errors = []interfaces.ConfigValidationError{}
	result.Summary.ErrorCount = 0

	// Test number too large
	err = sm.validatePropertyAgainstSchema("port", float64(101), schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
}

func TestSchemaManager_ValidatePropertyAgainstSchema_TypeMismatch(t *testing.T) {
	sm := NewSchemaManager()

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Test type mismatch - expecting string but got number
	schema := interfaces.PropertySchema{
		Type: "string",
	}

	err := sm.validatePropertyAgainstSchema("name", 123, schema, result)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Equal(t, 1, result.Summary.ErrorCount)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "type_error", result.Errors[0].Type)
}

func TestSchemaManager_ValidationRules(t *testing.T) {
	sm := NewSchemaManager()

	// Test getting validation rules
	rules := sm.GetValidationRules("package.json")
	assert.NotEmpty(t, rules)

	// Test adding validation rule
	newRule := ValidationSchemaRule{
		ID:          "test_rule",
		Name:        "Test Rule",
		Description: "Test validation rule",
		Category:    "test",
		Severity:    "warning",
		Enabled:     true,
		FileTypes:   []string{"test.json"},
	}

	sm.AddValidationRule("test.json", newRule)
	testRules := sm.GetValidationRules("test.json")
	assert.Len(t, testRules, 1)
	assert.Equal(t, "test_rule", testRules[0].ID)

	// Test removing validation rule
	sm.RemoveValidationRule("test.json", "test_rule")
	testRules = sm.GetValidationRules("test.json")
	assert.Empty(t, testRules)
}

func TestSchemaManager_ValidatePackageName(t *testing.T) {
	sm := NewSchemaManager()

	// Test valid package names
	validNames := []string{
		"my-package",
		"package.name",
		"package_name",
		"package123",
	}

	for _, name := range validNames {
		err := sm.ValidatePackageName(name)
		assert.NoError(t, err, "Expected %s to be valid", name)
	}

	// Test invalid package names
	invalidNames := []string{
		"My-Package",             // uppercase
		"package name",           // space
		"package@name",           // invalid character
		strings.Repeat("a", 215), // too long
	}

	for _, name := range invalidNames {
		err := sm.ValidatePackageName(name)
		assert.Error(t, err, "Expected %s to be invalid", name)
	}
}

func TestSchemaManager_ValidateEnvKey(t *testing.T) {
	sm := NewSchemaManager()

	// Test valid environment variable names
	validKeys := []string{
		"API_KEY",
		"DATABASE_URL",
		"PORT",
		"NODE_ENV",
	}

	for _, key := range validKeys {
		err := sm.ValidateEnvKey(key)
		assert.NoError(t, err, "Expected %s to be valid", key)
	}

	// Test invalid environment variable names
	invalidKeys := []string{
		"api_key", // lowercase
		"API-KEY", // hyphen
		"123_KEY", // starts with number
		"API KEY", // space
	}

	for _, key := range invalidKeys {
		err := sm.ValidateEnvKey(key)
		assert.Error(t, err, "Expected %s to be invalid", key)
	}
}

func TestSchemaManager_IsPotentialSecret(t *testing.T) {
	sm := NewSchemaManager()

	// Test potential secrets
	secretPairs := []struct {
		key      string
		value    string
		isSecret bool
	}{
		{"API_KEY", "sk-1234567890abcdef", true},
		{"PASSWORD", "supersecretpassword", true},
		{"DATABASE_TOKEN", "token_1234567890", true},
		{"PORT", "3000", false},
		{"NODE_ENV", "production", false},
		{"SECRET_KEY", "short", false}, // too short
	}

	for _, pair := range secretPairs {
		result := sm.IsPotentialSecret(pair.key, pair.value)
		assert.Equal(t, pair.isSecret, result,
			"Expected %s=%s to be secret: %v", pair.key, pair.value, pair.isSecret)
	}
}

func TestSchemaManager_DefaultSchemas(t *testing.T) {
	sm := NewSchemaManager()

	// Test that all expected default schemas are present
	expectedSchemas := []string{
		"package.json",
		"tsconfig.json",
		".eslintrc.json",
		"docker-compose.yml",
		"workflow.yml",
	}

	for _, schemaName := range expectedSchemas {
		schema, exists := sm.GetSchema(schemaName)
		assert.True(t, exists, "Expected schema %s to exist", schemaName)
		assert.NotNil(t, schema)
		assert.NotEmpty(t, schema.Title)
		assert.NotEmpty(t, schema.Description)
	}
}
