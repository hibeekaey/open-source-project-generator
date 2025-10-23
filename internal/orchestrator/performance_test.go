package orchestrator

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestParallelComponentGeneration(t *testing.T) {
	// Create a logger
	log := logger.NewLogger()
	log.SetLevel(logger.InfoLevel)

	// Create coordinator
	coordinator := NewProjectCoordinator(log)

	// Create test config with multiple components
	config := &models.ProjectConfig{
		Name:      "test-parallel",
		OutputDir: t.TempDir(),
		Components: []models.ComponentConfig{
			{Type: "nextjs", Name: "web", Enabled: true, Config: make(map[string]interface{})},
			{Type: "go-backend", Name: "api", Enabled: true, Config: make(map[string]interface{})},
			{Type: "android", Name: "mobile-android", Enabled: true, Config: make(map[string]interface{})},
		},
		Options: models.ProjectOptions{
			UseExternalTools: false, // Use fallback to avoid needing actual tools
			DryRun:           false,
			Verbose:          false,
			DisableParallel:  false, // Enable parallel generation
		},
	}

	// Create mock tool check result
	toolCheckResult := &models.ToolCheckResult{
		AllAvailable: false,
		Tools:        make(map[string]*models.Tool),
		Missing:      []string{"npx", "go", "gradle"},
		Outdated:     []string{},
	}

	// Test parallel generation
	ctx := context.Background()
	startTime := time.Now()

	results, err := coordinator.generateComponentsParallel(ctx, config, toolCheckResult)

	duration := time.Since(startTime)

	// Verify results
	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// Verify all components were processed
	for i, result := range results {
		assert.NotNil(t, result, "Result %d should not be nil", i)
		assert.Equal(t, config.Components[i].Type, result.Type)
		assert.Equal(t, config.Components[i].Name, result.Name)
	}

	t.Logf("Parallel generation completed in %v", duration)
}

func TestSequentialComponentGeneration(t *testing.T) {
	// Create a logger
	log := logger.NewLogger()
	log.SetLevel(logger.InfoLevel)

	// Create coordinator
	coordinator := NewProjectCoordinator(log)

	// Create test config with multiple components
	config := &models.ProjectConfig{
		Name:      "test-sequential",
		OutputDir: t.TempDir(),
		Components: []models.ComponentConfig{
			{Type: "nextjs", Name: "web", Enabled: true, Config: make(map[string]interface{})},
			{Type: "go-backend", Name: "api", Enabled: true, Config: make(map[string]interface{})},
		},
		Options: models.ProjectOptions{
			UseExternalTools: false,
			DryRun:           false,
			Verbose:          false,
			DisableParallel:  true, // Disable parallel generation
		},
	}

	// Create mock tool check result
	toolCheckResult := &models.ToolCheckResult{
		AllAvailable: false,
		Tools:        make(map[string]*models.Tool),
		Missing:      []string{"npx", "go"},
		Outdated:     []string{},
	}

	// Test sequential generation
	ctx := context.Background()
	startTime := time.Now()

	results, err := coordinator.generateComponentsSequential(ctx, config, toolCheckResult)

	duration := time.Since(startTime)

	// Verify results
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify all components were processed
	for i, result := range results {
		assert.NotNil(t, result, "Result %d should not be nil", i)
		assert.Equal(t, config.Components[i].Type, result.Type)
		assert.Equal(t, config.Components[i].Name, result.Name)
	}

	t.Logf("Sequential generation completed in %v", duration)
}

func TestProgressIndicator(t *testing.T) {
	var buf bytes.Buffer

	// Test with spinner disabled
	progress := NewProgressIndicator(&buf, "Testing progress", false)
	progress.Start()
	time.Sleep(50 * time.Millisecond)
	progress.Update("Updated message")
	time.Sleep(50 * time.Millisecond)
	progress.Stop()

	output := buf.String()
	assert.Contains(t, output, "Testing progress")
}

func TestStreamingWriter(t *testing.T) {
	var buf bytes.Buffer

	writer := NewStreamingWriter(&buf, "test", true)

	// Write some lines
	writer.Write([]byte("Line 1\n"))
	writer.Write([]byte("Line 2\n"))
	writer.Write([]byte("Line 3"))
	writer.Flush()

	output := buf.String()
	assert.Contains(t, output, "[test] Line 1")
	assert.Contains(t, output, "[test] Line 2")
	assert.Contains(t, output, "[test] Line 3")
}

func TestProgressTracker(t *testing.T) {
	var buf bytes.Buffer

	tracker := NewProgressTracker(&buf, 3, true)

	tracker.Increment("Item 1")
	tracker.Increment("Item 2")
	tracker.Increment("Item 3")
	tracker.Complete()

	output := buf.String()
	assert.Contains(t, output, "[3/3]")
	assert.Contains(t, output, "100%")
}
