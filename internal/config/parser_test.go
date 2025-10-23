package config
// Test file - gosec warnings suppressed for test utilities
//nolint:gosec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
	assert.NotEmpty(t, parser.SupportedFormats)
	assert.Contains(t, parser.SupportedFormats, "yaml")
	assert.Contains(t, parser.SupportedFormats, "yml")
	assert.Contains(t, parser.SupportedFormats, "json")
}

func TestParser_ParseString_YAML(t *testing.T) {
	parser := NewParser()

	yamlContent := `
name: test-project
description: A test project
output_dir: ./output
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true
integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
options:
  use_external_tools: true
  dry_run: false
  verbose: false
`

	config, err := parser.ParseString(yamlContent, "yaml")
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "test-project", config.Name)
	assert.Equal(t, "A test project", config.Description)
	assert.Equal(t, "./output", config.OutputDir)
	assert.Len(t, config.Components, 1)
	assert.Equal(t, "nextjs", config.Components[0].Type)
	assert.Equal(t, "web-app", config.Components[0].Name)
	assert.True(t, config.Components[0].Enabled)
	assert.True(t, config.Integration.GenerateDockerCompose)
	assert.True(t, config.Options.UseExternalTools)
}

func TestParser_ParseString_JSON(t *testing.T) {
	parser := NewParser()

	jsonContent := `{
  "name": "test-project",
  "description": "A test project",
  "output_dir": "./output",
  "components": [
    {
      "type": "go-backend",
      "name": "api-server",
      "enabled": true,
      "config": {
        "module": "github.com/user/test",
        "framework": "gin"
      }
    }
  ],
  "integration": {
    "generate_docker_compose": true,
    "generate_scripts": false
  },
  "options": {
    "use_external_tools": true,
    "dry_run": false
  }
}`

	config, err := parser.ParseString(jsonContent, "json")
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "test-project", config.Name)
	assert.Equal(t, "A test project", config.Description)
	assert.Len(t, config.Components, 1)
	assert.Equal(t, "go-backend", config.Components[0].Type)
	assert.Equal(t, "api-server", config.Components[0].Name)
}

func TestParser_ParseString_InvalidYAML(t *testing.T) {
	parser := NewParser()

	invalidYAML := `
name: test-project
components:
  - type: nextjs
    name: web-app
    invalid indentation
`

	_, err := parser.ParseString(invalidYAML, "yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML")
}

func TestParser_ParseString_InvalidJSON(t *testing.T) {
	parser := NewParser()

	invalidJSON := `{
  "name": "test-project",
  "components": [
    {
      "type": "nextjs"
      "name": "web-app"
    }
  ]
}`

	_, err := parser.ParseString(invalidJSON, "json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestParser_ParseString_UnsupportedFormat(t *testing.T) {
	parser := NewParser()

	_, err := parser.ParseString("some content", "xml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestParser_ParseFile_YAML(t *testing.T) {
	parser := NewParser()

	// Create temporary YAML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
name: file-test-project
description: Test from file
output_dir: ./output
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
integration:
  generate_docker_compose: true
options:
  use_external_tools: true
`

	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Parse the file
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "file-test-project", config.Name)
	assert.Equal(t, "Test from file", config.Description)
}

func TestParser_ParseFile_JSON(t *testing.T) {
	parser := NewParser()

	// Create temporary JSON file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	jsonContent := `{
  "name": "json-test-project",
  "description": "JSON test",
  "output_dir": "./output",
  "components": [],
  "integration": {},
  "options": {}
}`

	err := os.WriteFile(configPath, []byte(jsonContent), 0644)
	require.NoError(t, err)

	// Parse the file
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "json-test-project", config.Name)
}

func TestParser_ParseFile_NotFound(t *testing.T) {
	parser := NewParser()

	_, err := parser.ParseFile("/nonexistent/config.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestParser_ParseFile_UnsupportedExtension(t *testing.T) {
	parser := NewParser()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.txt")

	err := os.WriteFile(configPath, []byte("content"), 0644)
	require.NoError(t, err)

	_, err = parser.ParseFile(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file format")
}

func TestParser_DetectFormat(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		path     string
		expected string
	}{
		{"config.yaml", "yaml"},
		{"config.yml", "yml"},
		{"config.json", "json"},
		{"CONFIG.YAML", "yaml"},
		{"config.txt", ""},
		{"config", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			format := parser.detectFormat(tt.path)
			assert.Equal(t, tt.expected, format)
		})
	}
}

func TestParser_IsSupportedFormat(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		format   string
		expected bool
	}{
		{"yaml", true},
		{"yml", true},
		{"json", true},
		{"YAML", true},
		{"JSON", true},
		{"xml", false},
		{"txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := parser.IsSupportedFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_GenerateTemplate_YAML(t *testing.T) {
	parser := NewParser()

	template, err := parser.GenerateTemplate("yaml")
	require.NoError(t, err)
	assert.NotEmpty(t, template)

	// Verify it's valid YAML by parsing it back
	config, err := parser.ParseString(template, "yaml")
	require.NoError(t, err)
	assert.Equal(t, "my-project", config.Name)
	assert.NotEmpty(t, config.Components)
}

func TestParser_GenerateTemplate_JSON(t *testing.T) {
	parser := NewParser()

	template, err := parser.GenerateTemplate("json")
	require.NoError(t, err)
	assert.NotEmpty(t, template)

	// Verify it's valid JSON by parsing it back
	config, err := parser.ParseString(template, "json")
	require.NoError(t, err)
	assert.Equal(t, "my-project", config.Name)
	assert.NotEmpty(t, config.Components)
}

func TestParser_GenerateTemplate_UnsupportedFormat(t *testing.T) {
	parser := NewParser()

	_, err := parser.GenerateTemplate("xml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestParser_WriteTemplate(t *testing.T) {
	parser := NewParser()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := parser.WriteTemplate(configPath, "yaml")
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Verify content is valid
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, "my-project", config.Name)
}

func TestParser_WriteTemplate_FileExists(t *testing.T) {
	parser := NewParser()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create file first
	err := os.WriteFile(configPath, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Try to write template
	err = parser.WriteTemplate(configPath, "yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestParser_WriteTemplateForce(t *testing.T) {
	parser := NewParser()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create file first
	err := os.WriteFile(configPath, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Force write template
	err = parser.WriteTemplateForce(configPath, "yaml")
	require.NoError(t, err)

	// Verify content was overwritten
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, "my-project", config.Name)
}

func TestParser_ValidateFormat(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		path      string
		expectErr bool
	}{
		{"config.yaml", false},
		{"config.yml", false},
		{"config.json", false},
		{"config.txt", true},
		{"config", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := parser.ValidateFormat(tt.path)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParser_ParseFileWithValidation(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Valid configuration
	validYAML := `
name: valid-project
output_dir: ./output
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      name: web-app
      typescript: true
integration:
  generate_docker_compose: true
options:
  use_external_tools: true
`

	err := os.WriteFile(configPath, []byte(validYAML), 0644)
	require.NoError(t, err)

	config, err := parser.ParseFileWithValidation(configPath, validator)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "valid-project", config.Name)
}

func TestParser_ParseFileWithValidation_Invalid(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Invalid configuration (missing required fields)
	invalidYAML := `
name: ""
output_dir: ./output
components: []
`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = parser.ParseFileWithValidation(configPath, validator)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestParser_ParseWithValidation(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	validYAML := `
name: valid-project
output_dir: ./output
components:
  - type: go-backend
    name: api-server
    enabled: true
    config:
      name: api-server
      module: github.com/user/project
integration:
  generate_docker_compose: false
options:
  use_external_tools: true
`

	reader := strings.NewReader(validYAML)
	config, err := parser.ParseWithValidation(reader, "yaml", validator)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "valid-project", config.Name)
}

func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ParseError
		expected string
	}{
		{
			name: "with line and column",
			err: &ParseError{
				Path:    "config.yaml",
				Line:    10,
				Column:  5,
				Message: "invalid syntax",
			},
			expected: "parse error in config.yaml at line 10, column 5: invalid syntax",
		},
		{
			name: "with path only",
			err: &ParseError{
				Path:    "config.yaml",
				Message: "invalid syntax",
			},
			expected: "parse error in config.yaml: invalid syntax",
		},
		{
			name: "without path",
			err: &ParseError{
				Message: "invalid syntax",
			},
			expected: "parse error: invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestParser_ComplexConfiguration(t *testing.T) {
	parser := NewParser()

	complexYAML := `
name: complex-project
description: A complex multi-component project
output_dir: ./complex-output
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      name: web-app
      typescript: true
      tailwind: true
      app_router: true
      eslint: true
  - type: go-backend
    name: api-server
    enabled: true
    config:
      name: api-server
      module: github.com/user/complex-project
      framework: gin
      port: 8080
  - type: android
    name: mobile-android
    enabled: true
    config:
      name: mobile-android
      package: com.example.complex
      min_sdk: 24
      target_sdk: 34
      language: kotlin
  - type: ios
    name: mobile-ios
    enabled: true
    config:
      name: mobile-ios
      bundle_id: com.example.complex
      deployment_target: "15.0"
      language: swift
integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
    websocket: ws://localhost:8080/ws
  shared_environment:
    NODE_ENV: production
    LOG_LEVEL: info
    API_TIMEOUT: "30"
options:
  use_external_tools: true
  dry_run: false
  verbose: true
  create_backup: true
  force_overwrite: false
`

	config, err := parser.ParseString(complexYAML, "yaml")
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "complex-project", config.Name)
	assert.Len(t, config.Components, 4)
	assert.True(t, config.Integration.GenerateDockerCompose)
	assert.True(t, config.Integration.GenerateScripts)
	assert.Len(t, config.Integration.APIEndpoints, 2)
	assert.Len(t, config.Integration.SharedEnvironment, 3)
	assert.True(t, config.Options.Verbose)
	assert.True(t, config.Options.CreateBackup)
}
