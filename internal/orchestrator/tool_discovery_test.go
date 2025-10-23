package orchestrator

import (
	"runtime"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToolDiscovery(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	assert.NotNil(t, td)
	assert.NotNil(t, td.registry)
	assert.NotNil(t, td.cache)
	assert.Equal(t, 5*time.Minute, td.cacheTTL)

	// Verify known tools are registered
	tools := td.ListRegisteredTools()
	assert.Contains(t, tools, "npx")
	assert.Contains(t, tools, "go")
	assert.Contains(t, tools, "gradle")
	assert.Contains(t, tools, "xcodebuild")
	assert.Contains(t, tools, "docker")
	assert.Contains(t, tools, "terraform")
}

func TestToolDiscovery_IsAvailable(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test with 'go' which should be available in test environment
	available, err := td.IsAvailable("go")
	assert.NoError(t, err)
	assert.True(t, available)

	// Test with non-existent tool
	available, err = td.IsAvailable("nonexistenttool12345")
	assert.NoError(t, err)
	assert.False(t, available)
}

func TestToolDiscovery_GetVersion(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test with 'go' which should be available
	version, err := td.GetVersion("go")
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
	assert.Contains(t, version, "go")
}

func TestToolDiscovery_CheckRequirements(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test with 'go' which should be available
	result, err := td.CheckRequirements([]string{"go"})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	checkResult, ok := result.(*models.ToolCheckResult)
	require.True(t, ok)
	assert.True(t, checkResult.AllAvailable)
	assert.Empty(t, checkResult.Missing)
	assert.Contains(t, checkResult.Tools, "go")

	// Test with mix of available and unavailable tools
	result, err = td.CheckRequirements([]string{"go", "nonexistenttool12345"})
	assert.NoError(t, err)

	checkResult, ok = result.(*models.ToolCheckResult)
	require.True(t, ok)
	assert.False(t, checkResult.AllAvailable)
	assert.Contains(t, checkResult.Missing, "nonexistenttool12345")
}

func TestToolDiscovery_GetInstallInstructions(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test with registered tool
	instructions := td.GetInstallInstructions("go", runtime.GOOS)
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "go")
	assert.Contains(t, instructions, "Installation instructions")

	// Test with unregistered tool
	instructions = td.GetInstallInstructions("unknowntool", runtime.GOOS)
	assert.Contains(t, instructions, "not registered")
}

func TestToolDiscovery_Cache(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)
	td.SetCacheTTL(1 * time.Second)

	// First check - should hit the actual system
	available1, err := td.IsAvailable("go")
	assert.NoError(t, err)

	// Second check - should hit cache
	available2, err := td.IsAvailable("go")
	assert.NoError(t, err)
	assert.Equal(t, available1, available2)

	// Verify cache is working by checking stats
	stats := td.GetCacheStats()
	assert.NotNil(t, stats)

	// Clear cache
	td.ClearCache()

	// After clearing, cache should be empty
	stats = td.GetCacheStats()
	assert.NotNil(t, stats)
}

func TestToolDiscovery_GetToolMetadata(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test with registered tool
	metadata, err := td.GetToolMetadata("go")
	assert.NoError(t, err)
	assert.NotNil(t, metadata)
	assert.Equal(t, "go", metadata.Name)
	assert.Equal(t, "go", metadata.Command)
	assert.NotEmpty(t, metadata.InstallDocs)

	// Test with unregistered tool
	metadata, err = td.GetToolMetadata("unknowntool")
	assert.Error(t, err)
	assert.Nil(t, metadata)
}

func TestToolDiscovery_GetToolsForComponent(t *testing.T) {
	log := logger.NewLogger()
	td := NewToolDiscovery(log)

	// Test Next.js component
	tools := td.GetToolsForComponent("nextjs")
	assert.Contains(t, tools, "npx")

	// Test Go backend component
	tools = td.GetToolsForComponent("go-backend")
	assert.Contains(t, tools, "go")

	// Test Android component
	tools = td.GetToolsForComponent("android")
	assert.Contains(t, tools, "gradle")

	// Test iOS component
	tools = td.GetToolsForComponent("ios")
	assert.Contains(t, tools, "xcodebuild")
}

func TestNormalizeOS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"darwin", "darwin"},
		{"macos", "darwin"},
		{"osx", "darwin"},
		{"mac", "darwin"},
		{"linux", "linux"},
		{"unix", "linux"},
		{"windows", "windows"},
		{"win", "windows"},
		{"", runtime.GOOS},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeOS(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
