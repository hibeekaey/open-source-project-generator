package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAllExecutors_SupportsComponent tests SupportsComponent for all executors
func TestAllExecutors_SupportsComponent(t *testing.T) {
	tests := []struct {
		name          string
		executor      interface{ SupportsComponent(string) bool }
		componentType string
		expected      bool
	}{
		// NextJS Executor
		{
			name:          "NextJS supports nextjs",
			executor:      NewNextJSExecutor(),
			componentType: "nextjs",
			expected:      true,
		},
		{
			name:          "NextJS does not support go-backend",
			executor:      NewNextJSExecutor(),
			componentType: "go-backend",
			expected:      false,
		},
		{
			name:          "NextJS does not support android",
			executor:      NewNextJSExecutor(),
			componentType: "android",
			expected:      false,
		},
		{
			name:          "NextJS does not support ios",
			executor:      NewNextJSExecutor(),
			componentType: "ios",
			expected:      false,
		},
		// Go Executor
		{
			name:          "Go supports go-backend",
			executor:      NewGoExecutor(),
			componentType: "go-backend",
			expected:      true,
		},
		{
			name:          "Go does not support nextjs",
			executor:      NewGoExecutor(),
			componentType: "nextjs",
			expected:      false,
		},
		{
			name:          "Go does not support android",
			executor:      NewGoExecutor(),
			componentType: "android",
			expected:      false,
		},
		{
			name:          "Go does not support ios",
			executor:      NewGoExecutor(),
			componentType: "ios",
			expected:      false,
		},
		// Android Executor
		{
			name:          "Android supports android",
			executor:      NewAndroidExecutor(nil),
			componentType: "android",
			expected:      true,
		},
		{
			name:          "Android does not support nextjs",
			executor:      NewAndroidExecutor(nil),
			componentType: "nextjs",
			expected:      false,
		},
		{
			name:          "Android does not support go-backend",
			executor:      NewAndroidExecutor(nil),
			componentType: "go-backend",
			expected:      false,
		},
		{
			name:          "Android does not support ios",
			executor:      NewAndroidExecutor(nil),
			componentType: "ios",
			expected:      false,
		},
		// iOS Executor
		{
			name:          "iOS supports ios",
			executor:      NewiOSExecutor(nil),
			componentType: "ios",
			expected:      true,
		},
		{
			name:          "iOS does not support nextjs",
			executor:      NewiOSExecutor(nil),
			componentType: "nextjs",
			expected:      false,
		},
		{
			name:          "iOS does not support go-backend",
			executor:      NewiOSExecutor(nil),
			componentType: "go-backend",
			expected:      false,
		},
		{
			name:          "iOS does not support android",
			executor:      NewiOSExecutor(nil),
			componentType: "android",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.executor.SupportsComponent(tt.componentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAllExecutors_GetDefaultFlags tests GetDefaultFlags for all component types
func TestAllExecutors_GetDefaultFlags(t *testing.T) {
	tests := []struct {
		name          string
		executor      interface{ GetDefaultFlags(string) []string }
		componentType string
		wantNonEmpty  bool
		validateFlags func(*testing.T, []string)
	}{
		// NextJS
		{
			name:          "NextJS returns flags for nextjs",
			executor:      NewNextJSExecutor(),
			componentType: "nextjs",
			wantNonEmpty:  true,
			validateFlags: func(t *testing.T, flags []string) {
				assert.Contains(t, flags, "create-next-app@latest")
			},
		},
		{
			name:          "NextJS returns empty for unsupported",
			executor:      NewNextJSExecutor(),
			componentType: "android",
			wantNonEmpty:  false,
		},
		// Go
		{
			name:          "Go returns flags for go-backend",
			executor:      NewGoExecutor(),
			componentType: "go-backend",
			wantNonEmpty:  true,
			validateFlags: func(t *testing.T, flags []string) {
				assert.Contains(t, flags, "mod")
				assert.Contains(t, flags, "init")
			},
		},
		{
			name:          "Go returns empty for unsupported",
			executor:      NewGoExecutor(),
			componentType: "nextjs",
			wantNonEmpty:  false,
		},
		// Android
		{
			name:          "Android returns flags for android",
			executor:      NewAndroidExecutor(nil),
			componentType: "android",
			wantNonEmpty:  true,
			validateFlags: func(t *testing.T, flags []string) {
				assert.Contains(t, flags, "init")
			},
		},
		{
			name:          "Android returns empty for unsupported",
			executor:      NewAndroidExecutor(nil),
			componentType: "ios",
			wantNonEmpty:  false,
		},
		// iOS
		{
			name:          "iOS returns flags for ios",
			executor:      NewiOSExecutor(nil),
			componentType: "ios",
			wantNonEmpty:  true,
			validateFlags: func(t *testing.T, flags []string) {
				assert.Contains(t, flags, "package")
				assert.Contains(t, flags, "init")
			},
		},
		{
			name:          "iOS returns empty for unsupported",
			executor:      NewiOSExecutor(nil),
			componentType: "android",
			wantNonEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := tt.executor.GetDefaultFlags(tt.componentType)

			if tt.wantNonEmpty {
				assert.NotEmpty(t, flags, "Expected non-empty flags for supported component type")
				if tt.validateFlags != nil {
					tt.validateFlags(t, flags)
				}
			} else {
				assert.Empty(t, flags, "Expected empty flags for unsupported component type")
			}
		})
	}
}

// TestAllExecutors_ValidateConfig_ValidInputs tests validation with valid inputs
func TestAllExecutors_ValidateConfig_ValidInputs(t *testing.T) {
	tests := []struct {
		name     string
		executor interface {
			ValidateConfig(map[string]interface{}) error
		}
		config map[string]interface{}
	}{
		// NextJS valid configs
		{
			name:     "NextJS minimal valid config",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name": "my-app",
			},
		},
		{
			name:     "NextJS full valid config",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":       "my-app",
				"typescript": true,
				"tailwind":   true,
				"app_router": true,
				"eslint":     true,
			},
		},
		{
			name:     "NextJS with false booleans",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":       "my-app",
				"typescript": false,
				"tailwind":   false,
			},
		},
		// Go valid configs
		{
			name:     "Go minimal valid config",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "my-backend",
				"module": "github.com/user/project",
			},
		},
		{
			name:     "Go full valid config with gin",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":      "my-backend",
				"module":    "github.com/user/project",
				"framework": "gin",
				"port":      8080,
			},
		},
		{
			name:     "Go with echo framework",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":      "my-backend",
				"module":    "example.com/api",
				"framework": "echo",
				"port":      3000,
			},
		},
		{
			name:     "Go with fiber framework",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":      "my-backend",
				"module":    "myapp.io/server",
				"framework": "fiber",
				"port":      9000,
			},
		},
		// Android valid configs
		{
			name:     "Android minimal valid config",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":    "my-android-app",
				"package": "com.example.app",
			},
		},
		{
			name:     "Android full valid config",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":       "my-android-app",
				"package":    "com.example.myapp",
				"min_sdk":    21,
				"target_sdk": 33,
				"language":   "kotlin",
			},
		},
		{
			name:     "Android with java",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":     "my-android-app",
				"package":  "com.company.product",
				"language": "java",
			},
		},
		// iOS valid configs
		{
			name:     "iOS minimal valid config",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":      "my-ios-app",
				"bundle_id": "com.example.app",
			},
		},
		{
			name:     "iOS full valid config",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":              "my-ios-app",
				"bundle_id":         "com.example.myapp",
				"deployment_target": "14.0",
				"language":          "swift",
			},
		},
		{
			name:     "iOS with objective-c",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":              "my-ios-app",
				"bundle_id":         "com.company.product",
				"deployment_target": "13.0",
				"language":          "objective-c",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.executor.ValidateConfig(tt.config)
			assert.NoError(t, err)
		})
	}
}

// TestAllExecutors_ValidateConfig_InvalidInputs tests validation with invalid inputs
func TestAllExecutors_ValidateConfig_InvalidInputs(t *testing.T) {
	tests := []struct {
		name     string
		executor interface {
			ValidateConfig(map[string]interface{}) error
		}
		config map[string]interface{}
	}{
		// NextJS invalid configs
		{
			name:     "NextJS missing name",
			executor: NewNextJSExecutor(),
			config:   map[string]interface{}{},
		},
		{
			name:     "NextJS invalid typescript type",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":       "my-app",
				"typescript": "yes",
			},
		},
		{
			name:     "NextJS invalid tailwind type",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":     "my-app",
				"tailwind": 1,
			},
		},
		// Go invalid configs
		{
			name:     "Go missing name",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"module": "github.com/user/project",
			},
		},
		{
			name:     "Go missing module",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name": "my-backend",
			},
		},
		{
			name:     "Go invalid module format",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "my-backend",
				"module": "invalid module",
			},
		},
		{
			name:     "Go port too low",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "my-backend",
				"module": "github.com/user/project",
				"port":   0,
			},
		},
		{
			name:     "Go port too high",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "my-backend",
				"module": "github.com/user/project",
				"port":   70000,
			},
		},
		{
			name:     "Go invalid framework",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":      "my-backend",
				"module":    "github.com/user/project",
				"framework": "django",
			},
		},
		// Android invalid configs
		{
			name:     "Android missing name",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"package": "com.example.app",
			},
		},
		{
			name:     "Android missing package",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name": "my-app",
			},
		},
		{
			name:     "Android invalid package format",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":    "my-app",
				"package": "invalid",
			},
		},
		{
			name:     "Android min_sdk too low",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":    "my-app",
				"package": "com.example.app",
				"min_sdk": 10,
			},
		},
		{
			name:     "Android target_sdk too high",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":       "my-app",
				"package":    "com.example.app",
				"target_sdk": 50,
			},
		},
		{
			name:     "Android invalid language",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":     "my-app",
				"package":  "com.example.app",
				"language": "python",
			},
		},
		// iOS invalid configs
		{
			name:     "iOS missing name",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"bundle_id": "com.example.app",
			},
		},
		{
			name:     "iOS missing bundle_id",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name": "my-app",
			},
		},
		{
			name:     "iOS invalid bundle_id format",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":      "my-app",
				"bundle_id": "invalid",
			},
		},
		{
			name:     "iOS invalid deployment_target format",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":              "my-app",
				"bundle_id":         "com.example.app",
				"deployment_target": "invalid",
			},
		},
		{
			name:     "iOS invalid language",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":      "my-app",
				"bundle_id": "com.example.app",
				"language":  "rust",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.executor.ValidateConfig(tt.config)
			assert.Error(t, err, "Expected validation error for invalid config")
		})
	}
}

// TestExecutorRegistry_Integration tests the executor registry with all executors
func TestExecutorRegistry_Integration(t *testing.T) {
	// This test would require the ExecutorRegistry from orchestrator package
	// Testing that all executors can be registered and retrieved correctly

	executors := map[string]interface {
		SupportsComponent(string) bool
		GetDefaultFlags(string) []string
		ValidateConfig(map[string]interface{}) error
	}{
		"nextjs":     NewNextJSExecutor(),
		"go-backend": NewGoExecutor(),
		"android":    NewAndroidExecutor(nil),
		"ios":        NewiOSExecutor(nil),
	}

	for componentType, executor := range executors {
		t.Run("executor_"+componentType, func(t *testing.T) {
			// Test SupportsComponent
			assert.True(t, executor.SupportsComponent(componentType),
				"Executor should support its own component type")

			// Test GetDefaultFlags
			flags := executor.GetDefaultFlags(componentType)
			assert.NotEmpty(t, flags,
				"Executor should return default flags for its component type")

			// Test ValidateConfig with minimal valid config
			minimalConfig := map[string]interface{}{
				"name": "test-" + componentType,
			}

			// Add required fields based on component type
			switch componentType {
			case "go-backend":
				minimalConfig["module"] = "github.com/test/project"
			case "android":
				minimalConfig["package"] = "com.test.app"
			case "ios":
				minimalConfig["bundle_id"] = "com.test.app"
			}

			err := executor.ValidateConfig(minimalConfig)
			require.NoError(t, err,
				"Executor should accept minimal valid config")
		})
	}
}

// TestExecutor_CrossComponentValidation tests that executors reject configs for other component types
func TestExecutor_CrossComponentValidation(t *testing.T) {
	nextjsConfig := map[string]interface{}{
		"name":       "web-app",
		"typescript": true,
	}

	goConfig := map[string]interface{}{
		"name":   "api-server",
		"module": "github.com/test/api",
	}

	androidConfig := map[string]interface{}{
		"name":    "mobile-app",
		"package": "com.test.app",
	}

	iosConfig := map[string]interface{}{
		"name":      "mobile-app",
		"bundle_id": "com.test.app",
	}

	tests := []struct {
		name     string
		executor interface {
			ValidateConfig(map[string]interface{}) error
		}
		config  map[string]interface{}
		wantErr bool
	}{
		// NextJS executor with other configs
		{"NextJS with Go config", NewNextJSExecutor(), goConfig, true},
		{"NextJS with Android config", NewNextJSExecutor(), androidConfig, true},
		{"NextJS with iOS config", NewNextJSExecutor(), iosConfig, true},
		{"NextJS with NextJS config", NewNextJSExecutor(), nextjsConfig, false},

		// Go executor with other configs
		{"Go with NextJS config", NewGoExecutor(), nextjsConfig, true},
		{"Go with Android config", NewGoExecutor(), androidConfig, true},
		{"Go with iOS config", NewGoExecutor(), iosConfig, true},
		{"Go with Go config", NewGoExecutor(), goConfig, false},

		// Android executor with other configs
		{"Android with NextJS config", NewAndroidExecutor(nil), nextjsConfig, true},
		{"Android with Go config", NewAndroidExecutor(nil), goConfig, true},
		{"Android with iOS config", NewAndroidExecutor(nil), iosConfig, true},
		{"Android with Android config", NewAndroidExecutor(nil), androidConfig, false},

		// iOS executor with other configs
		{"iOS with NextJS config", NewiOSExecutor(nil), nextjsConfig, true},
		{"iOS with Go config", NewiOSExecutor(nil), goConfig, true},
		{"iOS with Android config", NewiOSExecutor(nil), androidConfig, true},
		{"iOS with iOS config", NewiOSExecutor(nil), iosConfig, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.executor.ValidateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
