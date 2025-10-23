package config

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidator_ValidConfiguration tests validation of valid configurations
func TestValidator_ValidConfiguration(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name   string
		config *models.ProjectConfig
	}{
		{
			name: "minimal valid configuration",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test project",
				OutputDir:   "/tmp/test-project",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
				},
				Options: models.ProjectOptions{
					UseExternalTools: true,
				},
			},
		},
		{
			name: "full configuration with all fields",
			config: &models.ProjectConfig{
				Name:        "full-project",
				Description: "Full test project",
				OutputDir:   "/tmp/full-project",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "web-app",
						Enabled: true,
						Config: map[string]interface{}{
							"name":       "web-app",
							"typescript": true,
							"tailwind":   true,
						},
					},
					{
						Type:    "go-backend",
						Name:    "api-server",
						Enabled: true,
						Config: map[string]interface{}{
							"name":      "api-server",
							"module":    "github.com/test/api",
							"framework": "gin",
						},
					},
				},
				Integration: models.IntegrationConfig{
					GenerateDockerCompose: true,
					GenerateScripts:       true,
					APIEndpoints: map[string]string{
						"backend": "http://localhost:8080",
					},
					SharedEnvironment: map[string]string{
						"LOG_LEVEL": "debug",
						"NODE_ENV":  "development",
					},
				},
				Options: models.ProjectOptions{
					UseExternalTools: true,
					DryRun:           false,
					Verbose:          true,
					CreateBackup:     true,
					ForceOverwrite:   false,
				},
			},
		},
		{
			name: "configuration with disabled component",
			config: &models.ProjectConfig{
				Name:        "partial-project",
				Description: "Project with disabled component",
				OutputDir:   "/tmp/partial-project",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
					{
						Type:    "android",
						Name:    "mobile",
						Enabled: false, // Disabled
						Config: map[string]interface{}{
							"name":    "mobile",
							"package": "com.test.mobile",
						},
					},
				},
				Options: models.ProjectOptions{
					UseExternalTools: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.config)
			assert.NoError(t, err)
		})
	}
}

// TestValidator_InvalidConfiguration tests validation of invalid configurations
func TestValidator_InvalidConfiguration(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectedErr string
	}{
		{
			name: "empty project name",
			config: &models.ProjectConfig{
				Name:        "",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components:  []models.ComponentConfig{},
			},
			expectedErr: "name",
		},
		{
			name: "unsupported component type",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "unsupported-type",
						Name:    "component",
						Enabled: true,
						Config:  map[string]interface{}{},
					},
				},
			},
			expectedErr: "unsupported component type",
		},
		{
			name: "duplicate component names",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
					{
						Type:    "go-backend",
						Name:    "frontend", // Duplicate name
						Enabled: true,
						Config: map[string]interface{}{
							"name":   "frontend",
							"module": "github.com/test/frontend",
						},
					},
				},
			},
			expectedErr: "duplicate component name",
		},
		{
			name: "empty component name",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "", // Empty name
						Enabled: true,
						Config:  map[string]interface{}{},
					},
				},
			},
			expectedErr: "component name is required",
		},
		{
			name: "invalid API endpoint",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
				},
				Integration: models.IntegrationConfig{
					APIEndpoints: map[string]string{
						"backend": "invalid-url", // Invalid URL
					},
				},
			},
			expectedErr: "endpoint must start with http",
		},
		{
			name: "invalid environment variable key",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
				},
				Integration: models.IntegrationConfig{
					SharedEnvironment: map[string]string{
						"invalid-key": "value", // Invalid key (lowercase with dash)
					},
				},
			},
			expectedErr: "contains invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// TestValidator_ComponentValidation tests component-specific validation
func TestValidator_ComponentValidation(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		component   models.ComponentConfig
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid nextjs component",
			component: models.ComponentConfig{
				Type:    "nextjs",
				Name:    "frontend",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "frontend",
					"typescript": true,
				},
			},
			shouldError: false,
		},
		{
			name: "valid go-backend component",
			component: models.ComponentConfig{
				Type:    "go-backend",
				Name:    "backend",
				Enabled: true,
				Config: map[string]interface{}{
					"name":   "backend",
					"module": "github.com/test/backend",
				},
			},
			shouldError: false,
		},
		{
			name: "invalid component type",
			component: models.ComponentConfig{
				Type:    "invalid-type",
				Name:    "component",
				Enabled: true,
				Config:  map[string]interface{}{},
			},
			shouldError: true,
			errorMsg:    "unsupported component type",
		},
		{
			name: "component with invalid name characters",
			component: models.ComponentConfig{
				Type:    "nextjs",
				Name:    "invalid name!", // Invalid characters
				Enabled: true,
				Config: map[string]interface{}{
					"name": "invalid name!",
				},
			},
			shouldError: true,
			errorMsg:    "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateComponent(&tt.component)
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidator_IntegrationValidation tests integration configuration validation
func TestValidator_IntegrationValidation(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		integration models.IntegrationConfig
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid integration config",
			integration: models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
				APIEndpoints: map[string]string{
					"backend": "http://localhost:8080",
					"api":     "https://api.example.com",
				},
				SharedEnvironment: map[string]string{
					"LOG_LEVEL": "debug",
					"NODE_ENV":  "development",
				},
			},
			shouldError: false,
		},
		{
			name: "empty integration config",
			integration: models.IntegrationConfig{
				GenerateDockerCompose: false,
				GenerateScripts:       false,
			},
			shouldError: false,
		},
		{
			name: "invalid endpoint URL",
			integration: models.IntegrationConfig{
				APIEndpoints: map[string]string{
					"backend": "not-a-url",
				},
			},
			shouldError: true,
			errorMsg:    "endpoint must start with http",
		},
		{
			name: "invalid environment variable",
			integration: models.IntegrationConfig{
				SharedEnvironment: map[string]string{
					"lowercase_key": "value",
				},
			},
			shouldError: true,
			errorMsg:    "contains invalid character",
		},
		{
			name: "environment variable starting with digit",
			integration: models.IntegrationConfig{
				SharedEnvironment: map[string]string{
					"1INVALID": "value",
				},
			},
			shouldError: true,
			errorMsg:    "cannot start with a digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateIntegration(&tt.integration)
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidator_ApplyDefaultsIntegration tests default value application
func TestValidator_ApplyDefaultsIntegration(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:        "test-project",
		Description: "Test",
		OutputDir:   "/tmp/test",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "frontend",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "frontend",
					// typescript and tailwind should get defaults
				},
			},
		},
	}

	err := validator.ApplyDefaults(config)
	require.NoError(t, err)

	// Check that defaults were applied
	frontendConfig := config.Components[0].Config
	assert.NotNil(t, frontendConfig["typescript"])
	assert.NotNil(t, frontendConfig["tailwind"])
}

// TestValidator_ValidateWithReportIntegration tests detailed validation reporting
func TestValidator_ValidateWithReportIntegration(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name           string
		config         *models.ProjectConfig
		expectValid    bool
		expectErrors   bool
		expectWarnings bool
	}{
		{
			name: "valid config with no warnings",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
				},
				Options: models.ProjectOptions{
					CreateBackup: true,
				},
			},
			expectValid:    true,
			expectErrors:   false,
			expectWarnings: false,
		},
		{
			name: "invalid config with errors",
			config: &models.ProjectConfig{
				Name:        "",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components:  []models.ComponentConfig{},
			},
			expectValid:    false,
			expectErrors:   true,
			expectWarnings: false,
		},
		{
			name: "valid config with warnings",
			config: &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "/tmp/test",
				Components: []models.ComponentConfig{
					{
						Type:    "nextjs",
						Name:    "frontend",
						Enabled: true,
						Config: map[string]interface{}{
							"name": "frontend",
						},
					},
				},
				Options: models.ProjectOptions{
					ForceOverwrite: true,
					CreateBackup:   false, // Warning: force overwrite without backup
				},
			},
			expectValid:    true,
			expectErrors:   false,
			expectWarnings: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := validator.ValidateWithReport(tt.config)

			assert.Equal(t, tt.expectValid, report.Valid)
			assert.Equal(t, tt.expectErrors, report.HasErrors())
			assert.Equal(t, tt.expectWarnings, report.HasWarnings())

			// Verify report string is not empty
			reportStr := report.String()
			assert.NotEmpty(t, reportStr)
		})
	}
}

// TestValidator_SchemaValidation tests schema-based validation
func TestValidator_SchemaValidation(t *testing.T) {
	validator := NewValidator()

	// Test that validator has access to schema
	supportedTypes := validator.GetSupportedComponentTypes()
	assert.NotEmpty(t, supportedTypes)
	assert.Contains(t, supportedTypes, "nextjs")
	assert.Contains(t, supportedTypes, "go-backend")

	// Test getting component schema
	schema, err := validator.GetComponentSchema("nextjs")
	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "nextjs", schema.Type)

	// Test getting schema for unsupported type
	_, err = validator.GetComponentSchema("unsupported")
	assert.Error(t, err)
}

// TestValidator_ComplexScenarios tests complex validation scenarios
func TestValidator_ComplexScenarios(t *testing.T) {
	validator := NewValidator()

	t.Run("multiple components with integration", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:        "complex-project",
			Description: "Complex test project",
			OutputDir:   "/tmp/complex",
			Components: []models.ComponentConfig{
				{
					Type:    "nextjs",
					Name:    "web-frontend",
					Enabled: true,
					Config: map[string]interface{}{
						"name": "web-frontend",
					},
				},
				{
					Type:    "go-backend",
					Name:    "api-backend",
					Enabled: true,
					Config: map[string]interface{}{
						"name":   "api-backend",
						"module": "github.com/test/api",
					},
				},
				{
					Type:    "android",
					Name:    "mobile-android",
					Enabled: true,
					Config: map[string]interface{}{
						"name":    "mobile-android",
						"package": "com.test.app",
					},
				},
			},
			Integration: models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
				APIEndpoints: map[string]string{
					"backend": "http://localhost:8080",
				},
			},
		}

		err := validator.Validate(config)
		assert.NoError(t, err)
	})

	t.Run("validate and apply defaults together", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:        "defaults-project",
			Description: "Test defaults",
			OutputDir:   "/tmp/defaults",
			Components: []models.ComponentConfig{
				{
					Type:    "nextjs",
					Name:    "frontend",
					Enabled: true,
					Config: map[string]interface{}{
						"name": "frontend",
					},
				},
			},
		}

		err := validator.ValidateAndApplyDefaults(config)
		assert.NoError(t, err)

		// Verify defaults were applied
		assert.NotNil(t, config.Components[0].Config["typescript"])
	})
}
