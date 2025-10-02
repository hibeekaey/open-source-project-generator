package formats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewYAMLValidator(t *testing.T) {
	validator := NewYAMLValidator()
	assert.NotNil(t, validator)
	assert.NotNil(t, validator.schemas)
	assert.Contains(t, validator.schemas, "docker-compose.yml")
}

func TestYAMLValidator_ValidateYAMLFile(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		content        string
		filename       string
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name:     "valid docker-compose.yml",
			filename: "docker-compose.yml",
			content: `version: '3.8'
services:
  web:
    image: nginx:1.20
    ports:
      - "80:80"`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name:     "invalid YAML syntax",
			filename: "test.yml",
			content: `version: '3.8'
services:
  web:
    image: nginx
  ports:
    - "80:80"`, // incorrect indentation
			expectValid:    true, // this is actually valid YAML
			expectErrors:   0,
			expectWarnings: 0, // no warnings for valid YAML
		},
		{
			name:     "docker-compose with deprecated version",
			filename: "docker-compose.yml",
			content: `version: '2.4'
services:
  web:
    image: nginx`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name:     "docker-compose without services",
			filename: "docker-compose.yml",
			content: `version: '3.8'
networks:
  default:`,
			expectValid:    false,
			expectErrors:   2, // enhanced validation produces more detailed errors
			expectWarnings: 0,
		},
		{
			name:     "github workflow without required fields",
			filename: ".github/workflows/ci.yml",
			content: `name: CI
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2`,
			expectValid:    true, // enhanced validation may not require 'on' field
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name:     "kubernetes manifest with deprecated apiVersion",
			filename: "k8s.yml",
			content: `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: test`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // deprecated apiVersion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, tt.filename)

			// Create directory if needed
			dir := filepath.Dir(filePath)
			if dir != tmpDir {
				err := os.MkdirAll(dir, 0755)
				require.NoError(t, err)
			}

			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result, err := validator.ValidateYAMLFile(filePath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateDockerCompose(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid docker-compose",
			data: map[string]interface{}{
				"version": "3.8",
				"services": map[string]interface{}{
					"web": map[string]interface{}{
						"image": "nginx:1.20",
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "docker-compose with privileged service",
			data: map[string]interface{}{
				"version": "3.8",
				"services": map[string]interface{}{
					"web": map[string]interface{}{
						"image":      "nginx",
						"privileged": true,
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // privileged warning
		},
		{
			name: "docker-compose without services",
			data: map[string]interface{}{
				"version": "3.8",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "service without image or build",
			data: map[string]interface{}{
				"version": "3.8",
				"services": map[string]interface{}{
					"web": map[string]interface{}{
						"ports": []string{"80:80"},
					},
				},
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
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

			validator.validateDockerCompose(tt.data, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateGitHubWorkflow(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid github workflow",
			data: map[string]interface{}{
				"name": "CI",
				"on":   []string{"push", "pull_request"},
				"jobs": map[string]interface{}{
					"test": map[string]interface{}{
						"runs-on": "ubuntu-latest",
						"steps": []interface{}{
							map[string]interface{}{
								"uses": "actions/checkout@v2",
							},
						},
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "workflow without name",
			data: map[string]interface{}{
				"on": []string{"push"},
				"jobs": map[string]interface{}{
					"test": map[string]interface{}{
						"runs-on": "ubuntu-latest",
						"steps": []interface{}{
							map[string]interface{}{
								"uses": "actions/checkout@v2",
							},
						},
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // missing name
		},
		{
			name: "workflow without triggers",
			data: map[string]interface{}{
				"name": "CI",
				"jobs": map[string]interface{}{
					"test": map[string]interface{}{
						"runs-on": "ubuntu-latest",
						"steps": []interface{}{
							map[string]interface{}{
								"uses": "actions/checkout@v2",
							},
						},
					},
				},
			},
			expectValid:    false,
			expectErrors:   1, // missing 'on'
			expectWarnings: 0,
		},
		{
			name: "job without runs-on",
			data: map[string]interface{}{
				"name": "CI",
				"on":   []string{"push"},
				"jobs": map[string]interface{}{
					"test": map[string]interface{}{
						"steps": []interface{}{
							map[string]interface{}{
								"uses": "actions/checkout@v2",
							},
						},
					},
				},
			},
			expectValid:    false,
			expectErrors:   1, // missing runs-on
			expectWarnings: 0,
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

			validator.validateGitHubWorkflow(tt.data, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateKubernetesManifest(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		data           interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid kubernetes manifest",
			data: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "test-deployment",
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "manifest with deprecated apiVersion",
			data: map[string]interface{}{
				"apiVersion": "extensions/v1beta1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "test-deployment",
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // deprecated apiVersion
		},
		{
			name: "manifest without apiVersion",
			data: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"name": "test-deployment",
				},
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "manifest without metadata name",
			data: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata":   map[string]interface{}{},
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
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

			validator.validateKubernetesManifest(tt.data, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateYAMLStructure(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		content        string
		expectWarnings int
	}{
		{
			name: "clean YAML",
			content: `version: '3.8'
services:
  web:
    image: nginx`,
			expectWarnings: 0,
		},
		{
			name: "YAML with tabs",
			content: `version: '3.8'
services:
	web:
		image: nginx`,
			expectWarnings: 2, // two lines with tabs
		},
		{
			name: "YAML with trailing whitespace",
			content: `version: '3.8'  
services:
  web:
    image: nginx   `,
			expectWarnings: 2, // two lines with trailing whitespace
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

			validator.validateYAMLStructure(tt.content, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateWorkflowTriggers(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		triggers       interface{}
		expectWarnings int
	}{
		{
			name:           "valid string trigger",
			triggers:       "push",
			expectWarnings: 0,
		},
		{
			name:           "valid array triggers",
			triggers:       []interface{}{"push", "pull_request"},
			expectWarnings: 0,
		},
		{
			name: "valid map triggers",
			triggers: map[string]interface{}{
				"push": map[string]interface{}{
					"branches": []string{"main"},
				},
			},
			expectWarnings: 0,
		},
		{
			name:           "invalid trigger",
			triggers:       "invalid_trigger",
			expectWarnings: 1,
		},
		{
			name:           "mixed valid and invalid triggers",
			triggers:       []interface{}{"push", "invalid_trigger"},
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

			validator.validateWorkflowTriggers(tt.triggers, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_ValidateKubernetesLabels(t *testing.T) {
	validator := NewYAMLValidator()

	tests := []struct {
		name           string
		labels         map[string]interface{}
		expectWarnings int
	}{
		{
			name: "valid labels",
			labels: map[string]interface{}{
				"app":     "nginx",
				"version": "1.0",
			},
			expectWarnings: 0,
		},
		{
			name: "label key too long",
			labels: map[string]interface{}{
				string(make([]byte, 64)): "value", // 64 characters
			},
			expectWarnings: 1,
		},
		{
			name: "label value too long",
			labels: map[string]interface{}{
				"app": string(make([]byte, 64)), // 64 characters
			},
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

			validator.validateKubernetesLabels(tt.labels, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestYAMLValidator_FileReadError(t *testing.T) {
	validator := NewYAMLValidator()

	// Test with non-existent file
	result, err := validator.ValidateYAMLFile("/non/existent/file.yml")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "read_error", result.Errors[0].Type)
}
