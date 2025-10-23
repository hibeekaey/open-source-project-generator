package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	assert.NotNil(t, validator)
	assert.NotNil(t, validator.schema)
}

func TestValidator_Validate_ValidConfig(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "web-app",
					"typescript": true,
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
		},
	}

	err := validator.Validate(config)
	assert.NoError(t, err)
}

func TestValidator_Validate_MissingName(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config:  map[string]interface{}{},
			},
		},
	}

	err := validator.Validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestValidator_Validate_MissingOutputDir(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config:  map[string]interface{}{},
			},
		},
	}

	err := validator.Validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output")
}

func TestValidator_Validate_NoEnabledComponents(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: false,
				Config:  map[string]interface{}{},
			},
		},
	}

	err := validator.Validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one component must be enabled")
}

func TestValidator_ValidateComponent_UnsupportedType(t *testing.T) {
	validator := NewValidator()

	comp := &models.ComponentConfig{
		Type:    "unsupported-type",
		Name:    "test",
		Enabled: true,
		Config:  map[string]interface{}{},
	}

	err := validator.ValidateComponent(comp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported component type")
}

func TestValidator_ValidateComponent_MissingName(t *testing.T) {
	validator := NewValidator()

	comp := &models.ComponentConfig{
		Type:    "nextjs",
		Name:    "",
		Enabled: true,
		Config:  map[string]interface{}{},
	}

	err := validator.ValidateComponent(comp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestValidator_ValidateComponent_InvalidName(t *testing.T) {
	validator := NewValidator()

	comp := &models.ComponentConfig{
		Type:    "nextjs",
		Name:    "invalid name with spaces",
		Enabled: true,
		Config:  map[string]interface{}{},
	}

	err := validator.ValidateComponent(comp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid component name")
}

func TestValidator_ValidateComponent_NextJS(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		config    map[string]interface{}
		expectErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"name":       "web-app",
				"typescript": true,
				"tailwind":   true,
			},
			expectErr: false,
		},
		{
			name: "missing required name",
			config: map[string]interface{}{
				"typescript": true,
			},
			expectErr: true,
		},
		{
			name: "invalid name",
			config: map[string]interface{}{
				"name": "invalid name!",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &models.ComponentConfig{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config:  tt.config,
			}

			err := validator.ValidateComponent(comp)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateComponent_GoBackend(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		config    map[string]interface{}
		expectErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"name":      "api-server",
				"module":    "github.com/user/project",
				"framework": "gin",
				"port":      8080,
			},
			expectErr: false,
		},
		{
			name: "missing module",
			config: map[string]interface{}{
				"name": "api-server",
			},
			expectErr: true,
		},
		{
			name: "invalid module",
			config: map[string]interface{}{
				"name":   "api-server",
				"module": "invalid",
			},
			expectErr: true,
		},
		{
			name: "invalid port",
			config: map[string]interface{}{
				"name":   "api-server",
				"module": "github.com/user/project",
				"port":   99999,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &models.ComponentConfig{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config:  tt.config,
			}

			err := validator.ValidateComponent(comp)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateComponent_Android(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		config    map[string]interface{}
		expectErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"name":       "mobile-android",
				"package":    "com.example.app",
				"min_sdk":    24,
				"target_sdk": 34,
			},
			expectErr: false,
		},
		{
			name: "missing package",
			config: map[string]interface{}{
				"name": "mobile-android",
			},
			expectErr: true,
		},
		{
			name: "invalid package",
			config: map[string]interface{}{
				"name":    "mobile-android",
				"package": "invalid",
			},
			expectErr: true,
		},
		{
			name: "package with invalid segment",
			config: map[string]interface{}{
				"name":    "mobile-android",
				"package": "com.123invalid.app",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &models.ComponentConfig{
				Type:    "android",
				Name:    "mobile-android",
				Enabled: true,
				Config:  tt.config,
			}

			err := validator.ValidateComponent(comp)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateComponent_iOS(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		config    map[string]interface{}
		expectErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"name":              "mobile-ios",
				"bundle_id":         "com.example.app",
				"deployment_target": "15.0",
			},
			expectErr: false,
		},
		{
			name: "missing bundle_id",
			config: map[string]interface{}{
				"name": "mobile-ios",
			},
			expectErr: true,
		},
		{
			name: "invalid bundle_id",
			config: map[string]interface{}{
				"name":      "mobile-ios",
				"bundle_id": "invalid",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &models.ComponentConfig{
				Type:    "ios",
				Name:    "mobile-ios",
				Enabled: true,
				Config:  tt.config,
			}

			err := validator.ValidateComponent(comp)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateIntegration(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		integration models.IntegrationConfig
		expectErr   bool
	}{
		{
			name: "valid integration",
			integration: models.IntegrationConfig{
				GenerateDockerCompose: true,
				APIEndpoints: map[string]string{
					"backend": "http://localhost:8080",
				},
				SharedEnvironment: map[string]string{
					"NODE_ENV": "development",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid endpoint",
			integration: models.IntegrationConfig{
				APIEndpoints: map[string]string{
					"backend": "invalid-url",
				},
			},
			expectErr: true,
		},
		{
			name: "invalid environment variable key",
			integration: models.IntegrationConfig{
				SharedEnvironment: map[string]string{
					"invalid-key": "value",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateIntegration(&tt.integration)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateOptions(t *testing.T) {
	validator := NewValidator()

	// All option combinations should be valid (just logical checks, no errors)
	options := models.ProjectOptions{
		UseExternalTools: true,
		DryRun:           true,
		Verbose:          true,
		CreateBackup:     true,
		ForceOverwrite:   true,
	}

	err := validator.ValidateOptions(&options)
	assert.NoError(t, err)
}

func TestValidator_CheckDuplicateNames(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		components []models.ComponentConfig
		expectErr  bool
	}{
		{
			name: "no duplicates",
			components: []models.ComponentConfig{
				{Name: "web-app", Enabled: true},
				{Name: "api-server", Enabled: true},
			},
			expectErr: false,
		},
		{
			name: "duplicate names",
			components: []models.ComponentConfig{
				{Name: "web-app", Enabled: true},
				{Name: "web-app", Enabled: true},
			},
			expectErr: true,
		},
		{
			name: "duplicate but one disabled",
			components: []models.ComponentConfig{
				{Name: "web-app", Enabled: true},
				{Name: "web-app", Enabled: false},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkDuplicateNames(tt.components)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ApplyDefaults(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
					// typescript, tailwind, etc. should be added as defaults
				},
			},
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":   "api-server",
					"module": "github.com/user/project",
					// framework and port should be added as defaults
				},
			},
		},
	}

	err := validator.ApplyDefaults(config)
	require.NoError(t, err)

	// Check that defaults were applied
	assert.Equal(t, true, config.Components[0].Config["typescript"])
	assert.Equal(t, true, config.Components[0].Config["tailwind"])
	assert.Equal(t, "gin", config.Components[1].Config["framework"])
	assert.Equal(t, 8080, config.Components[1].Config["port"])
}

func TestValidator_ValidateAndApplyDefaults(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
				},
			},
		},
	}

	err := validator.ValidateAndApplyDefaults(config)
	require.NoError(t, err)

	// Check that defaults were applied
	assert.NotNil(t, config.Components[0].Config["typescript"])
}

func TestValidator_GetSupportedComponentTypes(t *testing.T) {
	validator := NewValidator()

	types := validator.GetSupportedComponentTypes()
	assert.NotEmpty(t, types)
	assert.Contains(t, types, "nextjs")
	assert.Contains(t, types, "go-backend")
	assert.Contains(t, types, "android")
	assert.Contains(t, types, "ios")
}

func TestValidator_GetComponentSchema(t *testing.T) {
	validator := NewValidator()

	schema, err := validator.GetComponentSchema("nextjs")
	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "nextjs", schema.Type)

	_, err = validator.GetComponentSchema("unsupported")
	assert.Error(t, err)
}

func TestValidationReport_AddError(t *testing.T) {
	report := &ValidationReport{
		Valid: true,
	}

	report.AddError("field1", "error message")
	assert.False(t, report.Valid)
	assert.Len(t, report.Errors, 1)
	assert.Equal(t, "field1", report.Errors[0].Field)
}

func TestValidationReport_AddWarning(t *testing.T) {
	report := &ValidationReport{
		Valid: true,
	}

	report.AddWarning("warning message")
	assert.True(t, report.Valid) // Warnings don't affect validity
	assert.Len(t, report.Warnings, 1)
}

func TestValidationReport_HasErrors(t *testing.T) {
	report := &ValidationReport{}
	assert.False(t, report.HasErrors())

	report.AddError("field", "error")
	assert.True(t, report.HasErrors())
}

func TestValidationReport_HasWarnings(t *testing.T) {
	report := &ValidationReport{}
	assert.False(t, report.HasWarnings())

	report.AddWarning("warning")
	assert.True(t, report.HasWarnings())
}

func TestValidationReport_String(t *testing.T) {
	report := &ValidationReport{
		Valid: false,
	}
	report.AddError("field1", "error 1")
	report.AddWarning("warning 1")

	str := report.String()
	assert.Contains(t, str, "Validation failed")
	assert.Contains(t, str, "error 1")
	assert.Contains(t, str, "warning 1")
}

func TestValidator_ValidateWithReport(t *testing.T) {
	validator := NewValidator()

	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
				},
			},
		},
	}

	report := validator.ValidateWithReport(config)
	assert.NotNil(t, report)
	assert.True(t, report.Valid)
	assert.False(t, report.HasErrors())
}

func TestValidator_ValidateWithReport_Warnings(t *testing.T) {
	validator := NewValidator()

	// Config with force overwrite but no backup
	config := &models.ProjectConfig{
		Name:      "test-project",
		OutputDir: "./output",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
				},
			},
		},
		Options: models.ProjectOptions{
			ForceOverwrite: true,
			CreateBackup:   false,
		},
	}

	report := validator.ValidateWithReport(config)
	assert.NotNil(t, report)
	assert.True(t, report.Valid)
	assert.True(t, report.HasWarnings())
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"valid name", "my-project", false},
		{"with underscore", "my_project", false},
		{"with dash", "my-project", false},
		{"alphanumeric", "project123", false},
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 101), true},
		{"with spaces", "my project", true},
		{"with special chars", "my-project!", true},
		{"starts with dash", "-project", true},
		{"ends with dash", "project-", true},
		{"not a string", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGoModule(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"valid module", "github.com/user/project", false},
		{"valid with subdirs", "github.com/user/project/subdir", false},
		{"empty string", "", true},
		{"no slash", "invalid", true},
		{"with spaces", "github.com/user/my project", true},
		{"not a string", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGoModule(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"valid port", 8080, false},
		{"valid port float", 8080.0, false},
		{"min port", 1, false},
		{"max port", 65535, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"too large", 99999, true},
		{"not a number", "8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAndroidPackage(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"valid package", "com.example.app", false},
		{"valid with underscores", "com.example.my_app", false},
		{"empty string", "", true},
		{"no dot", "invalid", true},
		{"single segment", "com", true},
		{"starts with number", "com.123example.app", true},
		{"empty segment", "com..app", true},
		{"not a string", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAndroidPackage(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBundleID(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"valid bundle id", "com.example.app", false},
		{"valid with hyphens", "com.example.my-app", false},
		{"empty string", "", true},
		{"no dot", "invalid", true},
		{"single segment", "com", true},
		{"starts with number", "com.123example.app", true},
		{"empty segment", "com..app", true},
		{"not a string", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBundleID(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		endpoint  string
		expectErr bool
	}{
		{"valid http", "http://localhost:8080", false},
		{"valid https", "https://api.example.com", false},
		{"empty", "", true},
		{"no protocol", "localhost:8080", true},
		{"with spaces", "http://local host:8080", true},
		{"with invalid chars", "http://localhost<>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEndpoint("test", tt.endpoint)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		expectErr bool
	}{
		{"valid", "NODE_ENV", "development", false},
		{"with numbers", "LOG_LEVEL_1", "info", false},
		{"empty key", "", "value", true},
		{"lowercase", "node_env", "value", true},
		{"starts with number", "1_VAR", "value", true},
		{"with dash", "NODE-ENV", "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnvironmentVariable(tt.key, tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputDir(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{"valid relative", "./output", false},
		{"valid absolute", "/tmp/output", false},
		{"empty", "", true},
		{"path traversal", "../../../etc", true},
		{"with invalid chars", "output*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutputDir(tt.path)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
