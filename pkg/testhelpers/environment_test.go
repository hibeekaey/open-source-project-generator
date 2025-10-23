package testhelpers

// Test file - gosec warnings suppressed for test utilities
//nolint:gosec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestEnv(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	assert.NotEmpty(t, env.TempDir)
	assert.NotNil(t, env.MockTools)
	assert.NotEmpty(t, env.MockBinDir)

	// Verify temp directory exists
	_, err := os.Stat(env.TempDir)
	assert.NoError(t, err)

	// Verify mock bin directory exists
	_, err = os.Stat(env.MockBinDir)
	assert.NoError(t, err)
}

func TestMockToolAvailable(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	err := env.MockToolAvailable("test-tool", "1.0.0")
	require.NoError(t, err)

	// Verify mock tool was created
	mockTool, exists := env.MockTools["test-tool"]
	assert.True(t, exists)
	assert.Equal(t, "test-tool", mockTool.Name)
	assert.Equal(t, "1.0.0", mockTool.Version)
	assert.True(t, mockTool.Available)

	// Verify script file exists
	_, err = os.Stat(mockTool.ScriptPath)
	assert.NoError(t, err)
}

func TestMockToolUnavailable(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	env.MockToolUnavailable("unavailable-tool")

	mockTool, exists := env.MockTools["unavailable-tool"]
	assert.True(t, exists)
	assert.False(t, mockTool.Available)
}

func TestMockToolWithBehavior(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	err := env.MockToolWithBehavior("custom-tool", 1, "output", "error")
	require.NoError(t, err)

	mockTool, exists := env.MockTools["custom-tool"]
	assert.True(t, exists)
	assert.Equal(t, 1, mockTool.ExitCode)
	assert.Equal(t, "output", mockTool.Stdout)
	assert.Equal(t, "error", mockTool.Stderr)
}

func TestCreateProjectStructure(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	projectDir, err := env.CreateProjectStructure("test-project")
	require.NoError(t, err)
	assert.NotEmpty(t, projectDir)

	// Verify directories were created
	expectedDirs := []string{
		filepath.Join(projectDir, "App"),
		filepath.Join(projectDir, "CommonServer"),
		filepath.Join(projectDir, "Mobile", "android"),
		filepath.Join(projectDir, "Mobile", "ios"),
		filepath.Join(projectDir, "Deploy"),
	}

	for _, dir := range expectedDirs {
		_, err := os.Stat(dir)
		assert.NoError(t, err, "Directory should exist: %s", dir)
	}
}

func TestCreateTestConfig(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	components := []string{"nextjs", "go-backend"}
	config := env.CreateTestConfig("test-project", components)

	assert.Equal(t, "test-project", config.Name)
	assert.Equal(t, "Test project", config.Description)
	assert.Len(t, config.Components, 2)

	// Verify nextjs component
	nextjsComp := config.Components[0]
	assert.Equal(t, "nextjs", nextjsComp.Type)
	assert.True(t, nextjsComp.Enabled)
	assert.Equal(t, true, nextjsComp.Config["typescript"])

	// Verify go-backend component
	goComp := config.Components[1]
	assert.Equal(t, "go-backend", goComp.Type)
	assert.True(t, goComp.Enabled)
	assert.Contains(t, goComp.Config["module"], "test-project")
}

func TestCreateMockNextJSProject(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	projectDir := filepath.Join(env.TempDir, "nextjs-project")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockNextJSProject(projectDir)
	require.NoError(t, err)

	// Verify essential files exist
	expectedFiles := []string{
		"package.json",
		"next.config.js",
		"app/page.tsx",
		"app/layout.tsx",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectDir, file)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "File should exist: %s", file)
	}

	// Verify package.json content
	content, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "next")
}

func TestCreateMockGoProject(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	projectDir := filepath.Join(env.TempDir, "go-project")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockGoProject(projectDir, "github.com/test/app")
	require.NoError(t, err)

	// Verify essential files exist
	expectedFiles := []string{
		"go.mod",
		"main.go",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectDir, file)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "File should exist: %s", file)
	}

	// Verify go.mod content
	content, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "github.com/test/app")
	assert.Contains(t, string(content), "gin")
}

func TestCreateMockAndroidProject(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	projectDir := filepath.Join(env.TempDir, "android-project")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockAndroidProject(projectDir)
	require.NoError(t, err)

	// Verify essential files exist
	expectedFiles := []string{
		"settings.gradle",
		"build.gradle",
		"app/build.gradle",
		"app/src/main/AndroidManifest.xml",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectDir, file)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "File should exist: %s", file)
	}
}

func TestCreateMockIOSProject(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	projectDir := filepath.Join(env.TempDir, "ios-project")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = env.CreateMockIOSProject(projectDir, "TestApp")
	require.NoError(t, err)

	// Verify essential files exist
	expectedFiles := []string{
		"Package.swift",
		"TestApp/TestAppApp.swift",
		"TestApp/ContentView.swift",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectDir, file)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "File should exist: %s", file)
	}

	// Verify Package.swift content
	content, err := os.ReadFile(filepath.Join(projectDir, "Package.swift"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "TestApp")
}

func TestAssertFileExists(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// Create a test file
	testFile := filepath.Join(env.TempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	// This should not fail
	AssertFileExists(t, testFile)
}

func TestAssertFileContains(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// Create a test file with content
	testFile := filepath.Join(env.TempDir, "test.txt")
	content := "Hello World"
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// This should not fail
	AssertFileContains(t, testFile, "Hello")
	AssertFileContains(t, testFile, "World")
}

func TestIsToolAvailable(t *testing.T) {
	// Test with a tool that should always be available
	available := IsToolAvailable("sh")
	assert.True(t, available, "sh should be available on Unix systems")

	// Test with a tool that should not exist
	available = IsToolAvailable("nonexistent-tool-xyz")
	assert.False(t, available)
}
