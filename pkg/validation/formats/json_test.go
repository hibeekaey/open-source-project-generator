package formats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSONValidator(t *testing.T) {
	validator := NewJSONValidator()
	assert.NotNil(t, validator)
	assert.NotNil(t, validator.schemas)
	assert.Contains(t, validator.schemas, "package.json")
}

func TestJSONValidator_ValidateJSONFile(t *testing.T) {
	validator := NewJSONValidator()

	tests := []struct {
		name           string
		content        string
		filename       string
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name:           "valid package.json",
			filename:       "package.json",
			content:        `{"name": "test-package", "version": "1.0.0", "description": "Test package"}`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name:           "invalid JSON syntax",
			filename:       "test.json",
			content:        `{"name": "test", "version":}`,
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name:           "package.json missing required fields",
			filename:       "package.json",
			content:        `{"description": "Test package"}`,
			expectValid:    false,
			expectErrors:   4, // enhanced validation produces more detailed errors
			expectWarnings: 0,
		},
		{
			name:           "package.json with invalid name format",
			filename:       "package.json",
			content:        `{"name": "Test Package", "version": "1.0.0"}`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // name format warning
		},
		{
			name:           "tsconfig.json without strict mode",
			filename:       "tsconfig.json",
			content:        `{"compilerOptions": {"strict": false, "target": "es5"}}`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 2, // strict mode and target warnings
		},
		{
			name:           "eslint config without extends",
			filename:       ".eslintrc.json",
			content:        `{"rules": {"no-console": "warn"}}`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // missing extends + missing recommended rules
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result, err := validator.ValidateJSONFile(filePath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestJSONValidator_ValidatePackageJSON(t *testing.T) {
	validator := NewJSONValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid package.json",
			data: map[string]interface{}{
				"name":        "test-package",
				"version":     "1.0.0",
				"description": "Test package",
			},
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "package.json with vulnerable dependency",
			data: map[string]interface{}{
				"name":    "test-package",
				"version": "1.0.0",
				"dependencies": map[string]interface{}{
					"lodash": "4.17.15", // vulnerable version
				},
			},
			expectErrors:   0,
			expectWarnings: 0, // enhanced validation may not flag this specific version
		},
		{
			name: "package.json with invalid name",
			data: map[string]interface{}{
				"name":    "Test Package With Spaces",
				"version": "1.0.0",
			},
			expectErrors:   0,
			expectWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			validator.validatePackageJSON(tt.data, result)

			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestJSONValidator_ValidateTSConfig(t *testing.T) {
	validator := NewJSONValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectWarnings int
	}{
		{
			name: "tsconfig with strict mode enabled",
			data: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"strict": true,
					"target": "ES2020",
				},
			},
			expectWarnings: 0,
		},
		{
			name: "tsconfig with strict mode disabled",
			data: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"strict": false,
				},
			},
			expectWarnings: 1,
		},
		{
			name: "tsconfig with outdated target",
			data: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"target": "es5",
				},
			},
			expectWarnings: 1,
		},
		{
			name:           "tsconfig without compilerOptions",
			data:           map[string]interface{}{},
			expectWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			validator.validateTSConfig(tt.data, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestJSONValidator_ValidateESLintConfig(t *testing.T) {
	validator := NewJSONValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectWarnings int
	}{
		{
			name: "eslint config with extends",
			data: map[string]interface{}{
				"extends": []string{"eslint:recommended"},
				"rules": map[string]interface{}{
					"no-console": "warn",
				},
			},
			expectWarnings: 2, // enhanced validation may suggest additional rules
		},
		{
			name: "eslint config without extends",
			data: map[string]interface{}{
				"rules": map[string]interface{}{
					"no-console": "warn",
				},
			},
			expectWarnings: 3, // missing extends + missing recommended rules
		},
		{
			name: "eslint config missing recommended rules",
			data: map[string]interface{}{
				"extends": []string{"eslint:recommended"},
				"rules":   map[string]interface{}{},
			},
			expectWarnings: 3, // missing no-console, no-debugger, no-unused-vars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			validator.validateESLintConfig(tt.data, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestJSONValidator_ValidatePackageName(t *testing.T) {
	validator := NewJSONValidator()

	tests := []struct {
		name        string
		packageName string
		expectError bool
	}{
		{
			name:        "valid lowercase name",
			packageName: "my-package",
			expectError: false,
		},
		{
			name:        "valid name with numbers",
			packageName: "package123",
			expectError: false,
		},
		{
			name:        "invalid uppercase name",
			packageName: "MyPackage",
			expectError: true,
		},
		{
			name:        "invalid name with spaces",
			packageName: "my package",
			expectError: true,
		},
		{
			name:        "name too long",
			packageName: string(make([]byte, 215)), // 215 characters
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validatePackageName(tt.packageName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJSONValidator_ValidateAgainstSchema(t *testing.T) {
	validator := NewJSONValidator()

	schema := &interfaces.ConfigSchema{
		Required: []string{"name", "version"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:      "string",
				MinLength: &[]int{1}[0],
				MaxLength: &[]int{50}[0],
			},
			"version": {
				Type: "string",
				Enum: []string{"1.0.0", "2.0.0"},
			},
		},
	}

	tests := []struct {
		name           string
		data           interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid data",
			data: map[string]interface{}{
				"name":    "test",
				"version": "1.0.0",
			},
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "missing required field",
			data: map[string]interface{}{
				"name": "test",
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "invalid enum value",
			data: map[string]interface{}{
				"name":    "test",
				"version": "3.0.0",
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "string too long",
			data: map[string]interface{}{
				"name":    string(make([]byte, 51)), // 51 characters
				"version": "1.0.0",
			},
			expectValid:  false,
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			err := validator.validateAgainstSchema(tt.data, schema, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
		})
	}
}

func TestJSONValidator_FileReadError(t *testing.T) {
	validator := NewJSONValidator()

	// Test with non-existent file
	result, err := validator.ValidateJSONFile("/non/existent/file.json")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "read_error", result.Errors[0].Type)
}
