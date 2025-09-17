package config

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestConfigValidator_ValidateProjectConfig(t *testing.T) {
	// Create a test schema
	schema := &interfaces.ConfigSchema{
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:      "string",
				Required:  true,
				Pattern:   "^[a-zA-Z0-9_-]+$",
				MinLength: &[]int{1}[0],
				MaxLength: &[]int{100}[0],
			},
			"organization": {
				Type:     "string",
				Required: true,
			},
		},
		Required: []string{"name", "organization"},
	}

	validator := NewConfigValidator(schema)

	tests := []struct {
		name           string
		config         *models.ProjectConfig
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name:           "nil config",
			config:         nil,
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "valid config",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				License:      "MIT",
				OutputPath:   "./output",
				Components: models.Components{
					Backend: models.BackendComponents{
						GoGin: true,
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "missing required fields",
			config: &models.ProjectConfig{
				Description: "A test project",
			},
			expectValid:    false,
			expectErrors:   3, // name, organization, output_path
			expectWarnings: 1, // no components selected
		},
		{
			name: "invalid name pattern",
			config: &models.ProjectConfig{
				Name:         "test project!", // spaces and special chars not allowed
				Organization: "test-org",
				OutputPath:   "./output",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 1, // no components selected
		},
		{
			name: "invalid email",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Email:        "invalid-email",
				OutputPath:   "./output",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 1, // no components selected
		},
		{
			name: "unsupported license",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				License:      "CUSTOM-LICENSE",
				OutputPath:   "./output",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 2, // license warning + no components selected
		},
		{
			name: "dangerous output path",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				OutputPath:   "/usr/local",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 1, // no components selected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateProjectConfig(tt.config)
			if err != nil {
				t.Fatalf("ValidateProjectConfig() error = %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("ValidateProjectConfig() valid = %v, want %v", result.Valid, tt.expectValid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("ValidateProjectConfig() errors = %d, want %d", len(result.Errors), tt.expectErrors)
				for _, err := range result.Errors {
					t.Logf("Error: %s - %s", err.Field, err.Message)
				}
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("ValidateProjectConfig() warnings = %d, want %d", len(result.Warnings), tt.expectWarnings)
				for _, warn := range result.Warnings {
					t.Logf("Warning: %s - %s", warn.Field, warn.Message)
				}
			}
		})
	}
}

func TestConfigValidator_isValidSemVer(t *testing.T) {
	validator := &ConfigValidator{}

	tests := []struct {
		version string
		valid   bool
	}{
		{"1.0.0", true},
		{"1.2.3", true},
		{"10.20.30", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0+build.1", true},
		{"1.0.0-alpha+build.1", true},
		{"1.0", false},
		{"1", false},
		{"1.0.0.0", false},
		{"v1.0.0", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := validator.isValidSemVer(tt.version)
			if result != tt.valid {
				t.Errorf("isValidSemVer(%q) = %v, want %v", tt.version, result, tt.valid)
			}
		})
	}
}
func TestConfigValidator_isDangerousPath(t *testing.T) {
	validator := &ConfigValidator{}

	tests := []struct {
		path      string
		dangerous bool
	}{
		// Safe paths
		{"./my-project", false},
		{"/home/user/projects", false},
		{"/Users/user/projects", false},
		{"/var/folders/xx/temp", false}, // macOS temp
		{"/var/tmp/temp", false},        // Linux temp
		{"/tmp/temp", false},

		// Dangerous paths
		{"/", true},
		{"/usr", true},
		{"/usr/local", true},
		{"/etc", true},
		{"/etc/config", true},
		{"/bin", true},
		{"/bin/bash", true},
		{"/sbin", true},
		{"/boot", true},
		{"/var", true},     // /var itself is dangerous
		{"/var/log", true}, // /var/log is dangerous
		{"/var/lib", true}, // /var/lib is dangerous
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := validator.isDangerousPath(tt.path)
			if result != tt.dangerous {
				t.Errorf("isDangerousPath(%q) = %v, want %v", tt.path, result, tt.dangerous)
			}
		})
	}
}
