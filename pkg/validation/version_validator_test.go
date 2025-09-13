package validation

import (
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewVersionValidator(t *testing.T) {
	validator := NewVersionValidator()

	if validator == nil {
		t.Fatal("NewVersionValidator should not return nil")
	}

	if validator.compatibilityMatrix == nil {
		t.Fatal("VersionValidator should have a compatibility matrix")
	}
}

func TestValidateNodeVersionConfig(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name           string
		config         *models.NodeVersionConfig
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid configuration",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "invalid runtime version format",
			config: &models.NodeVersionConfig{
				Runtime:      "invalid-version",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    false,
			expectErrors:   2, // Format error + compatibility error
			expectWarnings: 0,
		},
		{
			name: "invalid types version format",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "invalid-types",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    false,
			expectErrors:   2, // Format error + compatibility error
			expectWarnings: 0,
		},
		{
			name: "invalid npm version format",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   "invalid-npm",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "invalid docker image format",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "invalid-image",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "incompatible runtime and types versions",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^18.0.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "empty runtime version",
			config: &models.NodeVersionConfig{
				Runtime:      "",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:    false,
			expectErrors:   2, // Format error + compatibility error
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateNodeVersionConfig(tt.config)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d warnings: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}

			// Validate that ValidatedAt is set
			if result.ValidatedAt.IsZero() {
				t.Error("ValidatedAt should be set")
			}
		})
	}
}

func TestValidateVersionCompatibility(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name           string
		configs        []*models.NodeVersionConfig
		expectValid    bool
		expectWarnings int
	}{
		{
			name: "single config - no comparison needed",
			configs: []*models.NodeVersionConfig{
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
			},
			expectValid:    true,
			expectWarnings: 0,
		},
		{
			name: "consistent configurations",
			configs: []*models.NodeVersionConfig{
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
			},
			expectValid:    true,
			expectWarnings: 0,
		},
		{
			name: "inconsistent runtime versions",
			configs: []*models.NodeVersionConfig{
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
				{
					Runtime:      ">=18.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
			},
			expectValid:    true,
			expectWarnings: 1,
		},
		{
			name: "inconsistent types versions",
			configs: []*models.NodeVersionConfig{
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^20.17.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
				{
					Runtime:      ">=20.0.0",
					TypesPackage: "^18.0.0",
					NPMVersion:   ">=10.0.0",
					DockerImage:  "node:20-alpine",
				},
			},
			expectValid:    true,
			expectWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateVersionCompatibility(tt.configs)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d warnings: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

func TestValidateAgainstLTS(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name              string
		config            *models.NodeVersionConfig
		expectValid       bool
		expectWarnings    int
		expectSuggestions int
	}{
		{
			name: "LTS version (Node 20)",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectValid:       true,
			expectWarnings:    0,
			expectSuggestions: 0,
		},
		{
			name: "LTS version (Node 18)",
			config: &models.NodeVersionConfig{
				Runtime:      ">=18.0.0",
				TypesPackage: "^18.17.0",
				NPMVersion:   ">=9.0.0",
				DockerImage:  "node:18-alpine",
			},
			expectValid:       true,
			expectWarnings:    0,
			expectSuggestions: 0,
		},
		{
			name: "Non-LTS version (Node 19)",
			config: &models.NodeVersionConfig{
				Runtime:      ">=19.0.0",
				TypesPackage: "^19.0.0",
				NPMVersion:   ">=9.0.0",
				DockerImage:  "node:19-alpine",
			},
			expectValid:       true,
			expectWarnings:    1,
			expectSuggestions: 1,
		},
		{
			name: "Non-LTS version (Node 21)",
			config: &models.NodeVersionConfig{
				Runtime:      ">=21.0.0",
				TypesPackage: "^21.0.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:21-alpine",
			},
			expectValid:       true,
			expectWarnings:    1,
			expectSuggestions: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateAgainstLTS(tt.config)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d warnings: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}

			if len(result.Suggestions) != tt.expectSuggestions {
				t.Errorf("Expected %d suggestions, got %d suggestions: %v", tt.expectSuggestions, len(result.Suggestions), result.Suggestions)
			}
		})
	}
}

func TestValidateVersionFormat(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name      string
		version   string
		field     string
		expectErr bool
	}{
		// Valid formats
		{"valid semver", "1.2.3", "runtime", false},
		{"valid with >=", ">=20.0.0", "runtime", false},
		{"valid with ^", "^20.17.0", "types", false},
		{"valid with ~", "~10.0.0", "npm", false},
		{"valid with prerelease", "20.0.0-beta.1", "runtime", false},
		{"valid with build", "20.0.0+build.1", "runtime", false},

		// Invalid formats
		{"empty version", "", "runtime", true},
		{"invalid format", "invalid", "runtime", true},
		{"missing patch", "20.0", "runtime", true},
		{"non-numeric", "v20.0.0", "runtime", true},
		{"invalid operator", "==20.0.0", "runtime", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateVersionFormat(tt.version, tt.field)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error for version %s, but got none", tt.version)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for version %s, but got: %v", tt.version, err)
			}
		})
	}
}

func TestValidateDockerImageFormat(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name      string
		image     string
		expectErr bool
	}{
		// Valid formats
		{"valid node image", "node:20-alpine", false},
		{"valid with registry", "docker.io/node:20-alpine", false},
		{"valid node LTS", "node:20.17.0-alpine", false},

		// Invalid formats
		{"empty image", "", true},
		{"no tag", "node", true},
		{"non-node image", "ubuntu:20.04", true},
		{"invalid format", "invalid-image-format", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateDockerImageFormat(tt.image)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error for image %s, but got none", tt.image)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for image %s, but got: %v", tt.image, err)
			}
		})
	}
}

func TestValidateRuntimeTypesCompatibility(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name      string
		runtime   string
		types     string
		expectErr bool
	}{
		// Compatible versions
		{"same major version", ">=20.0.0", "^20.17.0", false},
		{"types one version ahead", ">=20.0.0", "^21.0.0", false},
		{"types two versions ahead", ">=20.0.0", "^22.0.0", false},

		// Incompatible versions
		{"types too old", ">=20.0.0", "^18.0.0", true},
		{"types too new", ">=20.0.0", "^23.0.0", true},
		{"major version mismatch", ">=18.0.0", "^22.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateRuntimeTypesCompatibility(tt.runtime, tt.types)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error for runtime %s and types %s, but got none", tt.runtime, tt.types)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for runtime %s and types %s, but got: %v", tt.runtime, tt.types, err)
			}
		})
	}
}

func TestExtractMajorVersion(t *testing.T) {
	validator := NewVersionValidator()

	tests := []struct {
		name          string
		version       string
		expectedMajor int
		expectErr     bool
	}{
		{"simple version", "20.0.0", 20, false},
		{"with >= operator", ">=20.0.0", 20, false},
		{"with ^ operator", "^20.17.0", 20, false},
		{"with ~ operator", "~18.0.0", 18, false},
		{"with < operator", "<21.0.0", 21, false},
		{"invalid format", "invalid", 0, true},
		{"empty version", "", 0, true},
		{"non-numeric major", "v20.0.0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, err := validator.extractMajorVersion(tt.version)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error for version %s, but got none", tt.version)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for version %s, but got: %v", tt.version, err)
			}

			if !tt.expectErr && major != tt.expectedMajor {
				t.Errorf("Expected major version %d for %s, got %d", tt.expectedMajor, tt.version, major)
			}
		})
	}
}

func TestGenerateVersionSuggestions(t *testing.T) {
	validator := NewVersionValidator()

	// Test with a config that differs from recommendations
	config := &models.NodeVersionConfig{
		Runtime:      ">=18.0.0",
		TypesPackage: "^18.0.0",
		NPMVersion:   ">=9.0.0",
		DockerImage:  "node:18-alpine",
	}

	suggestions := validator.generateVersionSuggestions(config)

	// Should have suggestions for all fields that differ from recommended
	if len(suggestions) == 0 {
		t.Error("Expected suggestions for outdated configuration")
	}

	// Verify suggestion structure
	for _, suggestion := range suggestions {
		if suggestion.Field == "" {
			t.Error("Suggestion field should not be empty")
		}
		if suggestion.CurrentValue == "" {
			t.Error("Suggestion current value should not be empty")
		}
		if suggestion.SuggestedValue == "" {
			t.Error("Suggestion suggested value should not be empty")
		}
		if suggestion.Reason == "" {
			t.Error("Suggestion reason should not be empty")
		}
		if suggestion.Priority == "" {
			t.Error("Suggestion priority should not be empty")
		}
	}
}

func TestGetDefaultCompatibilityMatrix(t *testing.T) {
	matrix := getDefaultCompatibilityMatrix()

	if matrix == nil {
		t.Fatal("Default compatibility matrix should not be nil")
	}

	// Verify the default Node.js configuration
	nodeConfig := matrix.NodeJS
	if nodeConfig.Runtime == "" {
		t.Error("Default runtime version should not be empty")
	}
	if nodeConfig.TypesPackage == "" {
		t.Error("Default types package version should not be empty")
	}
	if nodeConfig.NPMVersion == "" {
		t.Error("Default NPM version should not be empty")
	}
	if nodeConfig.DockerImage == "" {
		t.Error("Default Docker image should not be empty")
	}

	// Verify LTS status
	if !nodeConfig.LTSStatus {
		t.Error("Default configuration should use LTS version")
	}

	// Verify UpdatedAt is set
	if matrix.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set in default matrix")
	}
}

// Benchmark tests for performance validation
func BenchmarkValidateNodeVersionConfig(b *testing.B) {
	validator := NewVersionValidator()
	config := &models.NodeVersionConfig{
		Runtime:      ">=20.0.0",
		TypesPackage: "^20.17.0",
		NPMVersion:   ">=10.0.0",
		DockerImage:  "node:20-alpine",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateNodeVersionConfig(config)
	}
}

func BenchmarkValidateVersionCompatibility(b *testing.B) {
	validator := NewVersionValidator()
	configs := []*models.NodeVersionConfig{
		{
			Runtime:      ">=20.0.0",
			TypesPackage: "^20.17.0",
			NPMVersion:   ">=10.0.0",
			DockerImage:  "node:20-alpine",
		},
		{
			Runtime:      ">=20.0.0",
			TypesPackage: "^20.17.0",
			NPMVersion:   ">=10.0.0",
			DockerImage:  "node:20-alpine",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateVersionCompatibility(configs)
	}
}
