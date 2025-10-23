package bootstrap

// Test file - gosec warnings suppressed for test utilities
//nolint:gosec

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNextJSExecutor_Integration tests Next.js executor with real npx if available
func TestNextJSExecutor_Integration(t *testing.T) {
	if !testhelpers.IsToolAvailable("npx") {
		t.Skip("Skipping Next.js integration test: npx not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	executor := NewNextJSExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	spec := &BootstrapSpec{
		ComponentType: "nextjs",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":       "test-nextjs-app",
			"typescript": true,
			"tailwind":   true,
			"app_router": true,
		},
		Timeout: 5 * time.Minute,
	}

	result, err := executor.Execute(ctx, spec)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "npx", result.ToolUsed)
	assert.NotEmpty(t, result.OutputDir)

	// Verify essential Next.js files were created
	projectDir := filepath.Join(env.TempDir, "test-nextjs-app")
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "package.json"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "next.config.js"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "tsconfig.json"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "tailwind.config.ts"))

	// Verify app directory structure
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "app", "layout.tsx"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "app", "page.tsx"))
}

// TestNextJSExecutor_WithMockTool tests Next.js executor with mocked npx
func TestNextJSExecutor_WithMockTool(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	// Mock npx to simulate successful execution
	err := env.MockToolAvailable("npx", "10.2.3")
	require.NoError(t, err)

	// Create mock Next.js project structure
	projectDir := filepath.Join(env.TempDir, "test-nextjs-app")
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockNextJSProject(projectDir)
	require.NoError(t, err)

	executor := NewNextJSExecutor()

	spec := &BootstrapSpec{
		ComponentType: "nextjs",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":       "test-nextjs-app",
			"typescript": true,
		},
	}

	// Verify command building works
	args, err := executor.buildNextJSCommand(spec)
	require.NoError(t, err)
	assert.Contains(t, args, "create-next-app@latest")
	assert.Contains(t, args, "test-nextjs-app")
	assert.Contains(t, args, "--typescript")
}

// TestGoExecutor_Integration tests Go executor with real go command if available
func TestGoExecutor_Integration(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping Go integration test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	executor := NewGoExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	spec := &BootstrapSpec{
		ComponentType: "go-backend",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":      "test-go-backend",
			"module":    "github.com/test/backend",
			"framework": "gin",
		},
		Timeout: 2 * time.Minute,
	}

	result, err := executor.Execute(ctx, spec)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "go", result.ToolUsed)
	assert.NotEmpty(t, result.OutputDir)

	// Verify essential Go files were created
	projectDir := filepath.Join(env.TempDir, "test-go-backend")
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "main.go"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, ".gitignore"))
	testhelpers.AssertFileExists(t, filepath.Join(projectDir, "README.md"))

	// Verify go.mod contains correct module name
	testhelpers.AssertFileContains(t, filepath.Join(projectDir, "go.mod"), "github.com/test/backend")

	// Verify main.go contains Gin setup
	testhelpers.AssertFileContains(t, filepath.Join(projectDir, "main.go"), "gin.Default()")
}

// TestGoExecutor_WithMockTool tests Go executor with mocked go command
func TestGoExecutor_WithMockTool(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	// Mock go command
	err := env.MockToolAvailable("go", "1.21.0")
	require.NoError(t, err)

	// Create mock Go project structure
	projectDir := filepath.Join(env.TempDir, "test-go-backend")
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockGoProject(projectDir, "github.com/test/backend")
	require.NoError(t, err)

	executor := NewGoExecutor()

	// Verify component support
	assert.True(t, executor.SupportsComponent("go-backend"))
	assert.True(t, executor.SupportsComponent("go"))
	assert.True(t, executor.SupportsComponent("backend"))
	assert.False(t, executor.SupportsComponent("nextjs"))

	// Verify manual steps are provided
	spec := &BootstrapSpec{
		Config: map[string]interface{}{
			"name":   "test-go-backend",
			"module": "github.com/test/backend",
		},
	}
	steps := executor.GetManualSteps(spec)
	assert.NotEmpty(t, steps)
}

// TestAndroidExecutor_FallbackRequired tests Android executor fallback behavior
func TestAndroidExecutor_FallbackRequired(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	// Create mock tool discovery that reports no tools available
	mockDiscovery := &MockToolDiscovery{
		tools: map[string]bool{
			"gradle":  false,
			"android": false,
		},
	}

	executor := NewAndroidExecutor(mockDiscovery)
	ctx := context.Background()

	spec := &BootstrapSpec{
		ComponentType: "android",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":    "test-android-app",
			"package": "com.test.app",
		},
	}

	// Should return error indicating fallback is required
	_, err := executor.Execute(ctx, spec)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fallback required")

	// Verify component support
	assert.True(t, executor.SupportsComponent("android"))
	assert.True(t, executor.SupportsComponent("mobile-android"))
	assert.False(t, executor.SupportsComponent("ios"))
}

// TestAndroidExecutor_WithGradle tests Android executor when gradle is available
func TestAndroidExecutor_WithGradle(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	// Create mock tool discovery that reports gradle available
	mockDiscovery := &MockToolDiscovery{
		tools: map[string]bool{
			"gradle": true,
		},
		versions: map[string]string{
			"gradle": "8.0.0",
		},
	}

	executor := NewAndroidExecutor(mockDiscovery)
	ctx := context.Background()

	spec := &BootstrapSpec{
		ComponentType: "android",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":    "test-android-app",
			"package": "com.test.app",
		},
	}

	// Even with gradle, should require fallback (no built-in project generator)
	_, err := executor.Execute(ctx, spec)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fallback required")

	// Verify manual steps are provided
	steps := executor.GetManualSteps(spec)
	assert.NotEmpty(t, steps)
	assert.Contains(t, steps[0], "Android Studio")
}

// TestIOSExecutor_FallbackRequired tests iOS executor fallback behavior
func TestIOSExecutor_FallbackRequired(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	// Create mock tool discovery
	mockDiscovery := &MockToolDiscovery{
		tools: map[string]bool{
			"swift":      false,
			"xcodebuild": false,
		},
	}

	executor := NewiOSExecutor(mockDiscovery)
	ctx := context.Background()

	spec := &BootstrapSpec{
		ComponentType: "ios",
		TargetDir:     env.TempDir,
		Config: map[string]interface{}{
			"name":      "test-ios-app",
			"bundle_id": "com.test.app",
		},
	}

	// Should return error indicating fallback is required
	_, err := executor.Execute(ctx, spec)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fallback required")

	// Verify component support
	assert.True(t, executor.SupportsComponent("ios"))
	assert.True(t, executor.SupportsComponent("mobile-ios"))
	assert.False(t, executor.SupportsComponent("android"))

	// Verify manual steps are provided
	steps := executor.GetManualSteps(spec)
	assert.NotEmpty(t, steps)
	assert.Contains(t, steps[0], "Xcode")
}

// TestIOSExecutor_RequiresFallback tests that iOS always requires fallback
func TestIOSExecutor_RequiresFallback(t *testing.T) {
	mockDiscovery := &MockToolDiscovery{
		tools: map[string]bool{
			"swift":      true,
			"xcodebuild": true,
		},
		versions: map[string]string{
			"swift":      "5.9.0",
			"xcodebuild": "15.0",
		},
	}

	executor := NewiOSExecutor(mockDiscovery)

	// iOS should always require fallback since there's no CLI project generator
	assert.True(t, executor.RequiresFallback())
}

// MockToolDiscovery is a mock implementation of ToolDiscovery for testing
type MockToolDiscovery struct {
	tools    map[string]bool
	versions map[string]string
}

func (m *MockToolDiscovery) IsAvailable(toolName string) (bool, error) {
	if available, ok := m.tools[toolName]; ok {
		return available, nil
	}
	return false, nil
}

func (m *MockToolDiscovery) GetVersion(toolName string) (string, error) {
	if version, ok := m.versions[toolName]; ok {
		return version, nil
	}
	return "", nil
}

// TestExecutorSupportsComponent tests component type support across executors
func TestExecutorSupportsComponent(t *testing.T) {
	tests := []struct {
		name          string
		executor      interface{ SupportsComponent(string) bool }
		componentType string
		expected      bool
	}{
		{"NextJS supports nextjs", NewNextJSExecutor(), "nextjs", true},
		{"NextJS supports next", NewNextJSExecutor(), "next", true},
		{"NextJS supports frontend", NewNextJSExecutor(), "frontend", true},
		{"NextJS does not support go", NewNextJSExecutor(), "go", false},
		{"Go supports go-backend", NewGoExecutor(), "go-backend", true},
		{"Go supports go", NewGoExecutor(), "go", true},
		{"Go supports backend", NewGoExecutor(), "backend", true},
		{"Go does not support nextjs", NewGoExecutor(), "nextjs", false},
		{"Android supports android", NewAndroidExecutor(nil), "android", true},
		{"Android supports mobile-android", NewAndroidExecutor(nil), "mobile-android", true},
		{"Android does not support ios", NewAndroidExecutor(nil), "ios", false},
		{"iOS supports ios", NewiOSExecutor(nil), "ios", true},
		{"iOS supports mobile-ios", NewiOSExecutor(nil), "mobile-ios", true},
		{"iOS does not support android", NewiOSExecutor(nil), "android", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.executor.SupportsComponent(tt.componentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestExecutorGetManualSteps tests manual steps generation
func TestExecutorGetManualSteps(t *testing.T) {
	tests := []struct {
		name     string
		executor interface {
			GetManualSteps(*BootstrapSpec) []string
		}
		spec     *BootstrapSpec
		contains string
	}{
		{
			name:     "NextJS manual steps",
			executor: NewNextJSExecutor(),
			spec: &BootstrapSpec{
				Config: map[string]interface{}{
					"name": "test-app",
				},
			},
			contains: "npm",
		},
		{
			name:     "Go manual steps",
			executor: NewGoExecutor(),
			spec: &BootstrapSpec{
				Config: map[string]interface{}{
					"name": "test-app",
				},
			},
			contains: "go run",
		},
		{
			name:     "Android manual steps",
			executor: NewAndroidExecutor(nil),
			spec: &BootstrapSpec{
				Config: map[string]interface{}{
					"name": "test-app",
				},
			},
			contains: "Android Studio",
		},
		{
			name:     "iOS manual steps",
			executor: NewiOSExecutor(nil),
			spec: &BootstrapSpec{
				Config: map[string]interface{}{
					"name": "test-app",
				},
			},
			contains: "Xcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steps := tt.executor.GetManualSteps(tt.spec)
			assert.NotEmpty(t, steps)

			// Check that at least one step contains the expected text
			found := false
			for _, step := range steps {
				if contains(step, tt.contains) {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected to find '%s' in manual steps", tt.contains)
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
