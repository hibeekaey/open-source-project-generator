package bootstrap

import (
	"testing"
)

// TestExecutorGetDefaultFlags tests that all executors return appropriate default flags
func TestExecutorGetDefaultFlags(t *testing.T) {
	tests := []struct {
		name     string
		executor interface {
			GetDefaultFlags(componentType string) []string
			SupportsComponent(componentType string) bool
		}
		componentType string
		wantNonEmpty  bool
	}{
		{
			name:          "NextJS returns default flags for nextjs",
			executor:      NewNextJSExecutor(),
			componentType: "nextjs",
			wantNonEmpty:  true,
		},
		{
			name:          "NextJS returns empty for unsupported type",
			executor:      NewNextJSExecutor(),
			componentType: "android",
			wantNonEmpty:  false,
		},
		{
			name:          "Go returns default flags for go-backend",
			executor:      NewGoExecutor(),
			componentType: "go-backend",
			wantNonEmpty:  true,
		},
		{
			name:          "Go returns empty for unsupported type",
			executor:      NewGoExecutor(),
			componentType: "nextjs",
			wantNonEmpty:  false,
		},
		{
			name:          "Android returns default flags for android",
			executor:      NewAndroidExecutor(nil),
			componentType: "android",
			wantNonEmpty:  true,
		},
		{
			name:          "Android returns empty for unsupported type",
			executor:      NewAndroidExecutor(nil),
			componentType: "ios",
			wantNonEmpty:  false,
		},
		{
			name:          "iOS returns default flags for ios",
			executor:      NewiOSExecutor(nil),
			componentType: "ios",
			wantNonEmpty:  true,
		},
		{
			name:          "iOS returns empty for unsupported type",
			executor:      NewiOSExecutor(nil),
			componentType: "android",
			wantNonEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := tt.executor.GetDefaultFlags(tt.componentType)

			if tt.wantNonEmpty && len(flags) == 0 {
				t.Errorf("GetDefaultFlags() returned empty flags, want non-empty for supported component type")
			}

			if !tt.wantNonEmpty && len(flags) != 0 {
				t.Errorf("GetDefaultFlags() returned %d flags, want empty for unsupported component type", len(flags))
			}
		})
	}
}

// TestExecutorValidateConfig tests that all executors validate configuration correctly
func TestExecutorValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		executor interface {
			ValidateConfig(config map[string]interface{}) error
		}
		config  map[string]interface{}
		wantErr bool
	}{
		// NextJS tests
		{
			name:     "NextJS valid config",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":       "test-app",
				"typescript": true,
				"tailwind":   true,
			},
			wantErr: false,
		},
		{
			name:     "NextJS missing name",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"typescript": true,
			},
			wantErr: true,
		},
		{
			name:     "NextJS invalid typescript type",
			executor: NewNextJSExecutor(),
			config: map[string]interface{}{
				"name":       "test-app",
				"typescript": "yes",
			},
			wantErr: true,
		},
		// Go tests
		{
			name:     "Go valid config",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "test-backend",
				"module": "github.com/example/test",
				"port":   8080,
			},
			wantErr: false,
		},
		{
			name:     "Go missing module",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name": "test-backend",
			},
			wantErr: true,
		},
		{
			name:     "Go invalid port",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":   "test-backend",
				"module": "github.com/example/test",
				"port":   99999,
			},
			wantErr: true,
		},
		{
			name:     "Go invalid framework",
			executor: NewGoExecutor(),
			config: map[string]interface{}{
				"name":      "test-backend",
				"module":    "github.com/example/test",
				"framework": "invalid",
			},
			wantErr: true,
		},
		// Android tests
		{
			name:     "Android valid config",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":       "test-android",
				"package":    "com.example.test",
				"min_sdk":    21,
				"target_sdk": 33,
			},
			wantErr: false,
		},
		{
			name:     "Android invalid package",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":    "test-android",
				"package": "invalid",
			},
			wantErr: true,
		},
		{
			name:     "Android invalid min_sdk",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":    "test-android",
				"package": "com.example.test",
				"min_sdk": 10,
			},
			wantErr: true,
		},
		{
			name:     "Android invalid language",
			executor: NewAndroidExecutor(nil),
			config: map[string]interface{}{
				"name":     "test-android",
				"package":  "com.example.test",
				"language": "python",
			},
			wantErr: true,
		},
		// iOS tests
		{
			name:     "iOS valid config",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":              "test-ios",
				"bundle_id":         "com.example.test",
				"deployment_target": "14.0",
			},
			wantErr: false,
		},
		{
			name:     "iOS invalid bundle_id",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":      "test-ios",
				"bundle_id": "invalid",
			},
			wantErr: true,
		},
		{
			name:     "iOS invalid deployment_target",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":              "test-ios",
				"bundle_id":         "com.example.test",
				"deployment_target": "invalid",
			},
			wantErr: true,
		},
		{
			name:     "iOS invalid language",
			executor: NewiOSExecutor(nil),
			config: map[string]interface{}{
				"name":      "test-ios",
				"bundle_id": "com.example.test",
				"language":  "python",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.executor.ValidateConfig(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
