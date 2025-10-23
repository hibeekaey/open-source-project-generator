package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentConfigValidator_NextJS(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("nextjs", NewNextJSConfigValidator())

	tests := []struct {
		name      string
		config    map[string]interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid nextjs config",
			config: map[string]interface{}{
				"name":       "my-app",
				"typescript": true,
				"tailwind":   true,
				"app_router": true,
				"eslint":     true,
			},
			wantError: false,
		},
		{
			name: "typescript not boolean",
			config: map[string]interface{}{
				"name":       "my-app",
				"typescript": "yes",
			},
			wantError: true,
			errorMsg:  "typescript",
		},
		{
			name: "tailwind not boolean",
			config: map[string]interface{}{
				"name":     "my-app",
				"tailwind": 1,
			},
			wantError: true,
			errorMsg:  "tailwind",
		},
		{
			name: "app_router not boolean",
			config: map[string]interface{}{
				"name":       "my-app",
				"app_router": "true",
			},
			wantError: true,
			errorMsg:  "app_router",
		},
		{
			name: "eslint not boolean",
			config: map[string]interface{}{
				"name":   "my-app",
				"eslint": 1.5,
			},
			wantError: true,
			errorMsg:  "eslint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate("nextjs", tt.config)
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComponentConfigValidator_GoBackend(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("go-backend", NewGoConfigValidator())

	tests := []struct {
		name      string
		config    map[string]interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid go-backend config",
			config: map[string]interface{}{
				"name":      "my-api",
				"module":    "github.com/user/my-api",
				"framework": "gin",
				"port":      8080,
			},
			wantError: false,
		},
		{
			name: "invalid module path",
			config: map[string]interface{}{
				"name":   "my-api",
				"module": "invalid",
			},
			wantError: true,
			errorMsg:  "module",
		},
		{
			name: "port out of range",
			config: map[string]interface{}{
				"name":   "my-api",
				"module": "github.com/user/my-api",
				"port":   70000,
			},
			wantError: true,
			errorMsg:  "port",
		},
		{
			name: "invalid framework",
			config: map[string]interface{}{
				"name":      "my-api",
				"module":    "github.com/user/my-api",
				"framework": "express",
			},
			wantError: true,
			errorMsg:  "framework",
		},
		{
			name: "port as float64",
			config: map[string]interface{}{
				"name":   "my-api",
				"module": "github.com/user/my-api",
				"port":   8080.0,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate("go-backend", tt.config)
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComponentConfigValidator_Android(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("android", NewAndroidConfigValidator())

	tests := []struct {
		name      string
		config    map[string]interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid android config",
			config: map[string]interface{}{
				"name":       "MyApp",
				"package":    "com.example.myapp",
				"min_sdk":    24,
				"target_sdk": 34,
				"language":   "kotlin",
			},
			wantError: false,
		},
		{
			name: "invalid package name",
			config: map[string]interface{}{
				"name":    "MyApp",
				"package": "invalid",
			},
			wantError: true,
			errorMsg:  "package",
		},
		{
			name: "min_sdk greater than target_sdk",
			config: map[string]interface{}{
				"name":       "MyApp",
				"package":    "com.example.myapp",
				"min_sdk":    30,
				"target_sdk": 24,
			},
			wantError: true,
			errorMsg:  "min_sdk",
		},
		{
			name: "invalid language",
			config: map[string]interface{}{
				"name":     "MyApp",
				"package":  "com.example.myapp",
				"language": "swift",
			},
			wantError: true,
			errorMsg:  "language",
		},
		{
			name: "sdk as float64",
			config: map[string]interface{}{
				"name":       "MyApp",
				"package":    "com.example.myapp",
				"min_sdk":    24.0,
				"target_sdk": 34.0,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate("android", tt.config)
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComponentConfigValidator_iOS(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("ios", NewIOSConfigValidator())

	tests := []struct {
		name      string
		config    map[string]interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid ios config",
			config: map[string]interface{}{
				"name":              "MyApp",
				"bundle_id":         "com.example.myapp",
				"deployment_target": "15.0",
				"language":          "swift",
			},
			wantError: false,
		},
		{
			name: "invalid bundle_id",
			config: map[string]interface{}{
				"name":      "MyApp",
				"bundle_id": "invalid",
			},
			wantError: true,
			errorMsg:  "bundle_id",
		},
		{
			name: "invalid deployment_target",
			config: map[string]interface{}{
				"name":              "MyApp",
				"bundle_id":         "com.example.myapp",
				"deployment_target": "invalid",
			},
			wantError: true,
			errorMsg:  "deployment_target",
		},
		{
			name: "invalid language",
			config: map[string]interface{}{
				"name":      "MyApp",
				"bundle_id": "com.example.myapp",
				"language":  "kotlin",
			},
			wantError: true,
			errorMsg:  "language",
		},
		{
			name: "deployment_target as float",
			config: map[string]interface{}{
				"name":              "MyApp",
				"bundle_id":         "com.example.myapp",
				"deployment_target": 15.0,
			},
			wantError: false,
		},
		{
			name: "deployment_target as int",
			config: map[string]interface{}{
				"name":              "MyApp",
				"bundle_id":         "com.example.myapp",
				"deployment_target": 15,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate("ios", tt.config)
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComponentConfigValidator_GetRequiredFields(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("nextjs", NewNextJSConfigValidator())
	validator.RegisterValidator("go-backend", NewGoConfigValidator())
	validator.RegisterValidator("android", NewAndroidConfigValidator())
	validator.RegisterValidator("ios", NewIOSConfigValidator())

	tests := []struct {
		componentType string
		expected      []string
	}{
		{
			componentType: "nextjs",
			expected:      []string{"name"},
		},
		{
			componentType: "go-backend",
			expected:      []string{"name", "module"},
		},
		{
			componentType: "android",
			expected:      []string{"name", "package"},
		},
		{
			componentType: "ios",
			expected:      []string{"name", "bundle_id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.componentType, func(t *testing.T) {
			fields := validator.GetRequiredFields(tt.componentType)
			assert.ElementsMatch(t, tt.expected, fields)
		})
	}
}

func TestComponentConfigValidator_GetFieldDescription(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("nextjs", NewNextJSConfigValidator())

	desc := validator.GetFieldDescription("nextjs", "typescript")
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "TypeScript")
}

func TestComponentConfigValidator_ValidateWithDetails(t *testing.T) {
	validator := NewComponentConfigValidator()
	validator.RegisterValidator("nextjs", NewNextJSConfigValidator())

	config := map[string]interface{}{
		"typescript": "not-a-boolean",
	}

	result := validator.ValidateWithDetails("nextjs", config)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Field, "name")
}
