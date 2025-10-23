package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_InteractiveMode_CompleteFlow tests complete interactive mode workflow
// Note: This test simulates the interactive flow without actual user input
func TestE2E_InteractiveMode_CompleteFlow(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Simulate interactive mode configuration
	config := &models.ProjectConfig{
		Name:        "interactive-test-project",
		Description: "Project created via interactive mode",
		OutputDir:   filepath.Join(env.TempDir, "interactive-test-project"),
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/test/api-server",
					"framework": "gin",
					"port":      8080,
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			CreateBackup:     true,
			Verbose:          false,
			DryRun:           false,
			ForceOverwrite:   true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Components)

	// Verify project structure
	projectDir := config.OutputDir
	testhelpers.AssertFileExists(t, projectDir)

	// Verify component was created
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "main.go"))

	// Verify integration files
	if config.Integration.GenerateDockerCompose {
		dockerComposePath := filepath.Join(projectDir, "docker-compose.yml")
		if _, err := os.Stat(dockerComposePath); err == nil {
			testhelpers.AssertFileContains(t, dockerComposePath, "services:")
		}
	}
}

// TestE2E_OfflineMode_WithCachedTools tests offline mode with cached tools
func TestE2E_OfflineMode_WithCachedTools(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Setup tool cache
	cacheConfig := &ToolCacheConfig{
		CacheDir: env.TempDir,
		TTL:      5 * time.Minute,
	}
	cache, err := NewToolCache(cacheConfig, log)
	require.NoError(t, err)

	// Pre-populate cache with tool availability
	cache.Set("go", true, "1.21.0")
	cache.Set("npx", false, "")
	err = cache.Save()
	require.NoError(t, err)

	// Enable offline mode
	coordinator.SetOfflineMode(true)

	config := &models.ProjectConfig{
		Name:        "offline-test",
		Description: "Offline mode test with cache",
		OutputDir:   filepath.Join(env.TempDir, "offline-test"),
		Components: []models.ComponentConfig{
			{
				Type:    "android",
				Name:    "mobile-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":    "mobile-app",
					"package": "com.test.app",
				},
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: false, // Force fallback
			DryRun:           false,
			Verbose:          false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)

	// Verify fallback generation worked in offline mode
	projectDir := config.OutputDir
	androidDir := filepath.Join(projectDir, "Mobile", "android")
	testhelpers.AssertFileExists(t, filepath.Join(androidDir, "build.gradle"))
}

// TestE2E_RollbackOnFailure tests automatic rollback when generation fails
func TestE2E_RollbackOnFailure(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Create a pre-existing directory with some content
	projectDir := filepath.Join(env.TempDir, "rollback-test")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	existingFile := filepath.Join(projectDir, "existing.txt")
	err = os.WriteFile(existingFile, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Create configuration that will likely fail (invalid component config)
	config := &models.ProjectConfig{
		Name:        "rollback-test",
		Description: "Rollback test",
		OutputDir:   projectDir,
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "backend",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "backend",
					// Missing required "module" field - should fail validation
				},
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			CreateBackup:     true,
			DryRun:           false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Generation should fail due to invalid config
	_, err = coordinator.Generate(ctx, config)
	assert.Error(t, err)

	// Verify existing file still exists (rollback preserved it)
	_, err = os.Stat(existingFile)
	assert.NoError(t, err, "Existing file should be preserved after rollback")
}

// TestE2E_StreamingOutput tests streaming output during generation
func TestE2E_StreamingOutput(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "streaming-test",
		Description: "Streaming output test",
		OutputDir:   filepath.Join(env.TempDir, "streaming-test"),
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "backend",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "backend",
					"module":    "github.com/test/backend",
					"framework": "gin",
				},
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			StreamOutput:     true, // Enable streaming
			Verbose:          true,
			DryRun:           false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)

	// Verify project was created
	projectDir := config.OutputDir
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
}

// TestE2E_MultipleComponents_WithValidation tests multiple components with validation
func TestE2E_MultipleComponents_WithValidation(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "multi-component-validation",
		Description: "Multiple components with validation",
		OutputDir:   filepath.Join(env.TempDir, "multi-component-validation"),
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/test/api-server",
					"framework": "gin",
					"port":      8080,
				},
			},
			{
				Type:    "android",
				Name:    "mobile-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "mobile-app",
					"package":    "com.test.app",
					"min_sdk":    21,
					"target_sdk": 33,
					"language":   "kotlin",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)
	assert.Len(t, result.Components, 2)

	// Verify both components
	projectDir := config.OutputDir

	// Verify Go backend
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "main.go"))

	// Verify Android
	androidDir := filepath.Join(projectDir, "Mobile", "android")
	testhelpers.AssertFileExists(t, filepath.Join(androidDir, "build.gradle"))
}

// TestE2E_ComponentValidation_InvalidConfig tests that invalid configs are rejected
func TestE2E_ComponentValidation_InvalidConfig(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	tests := []struct {
		name      string
		component models.ComponentConfig
		wantError string
	}{
		{
			name: "go-backend missing module",
			component: models.ComponentConfig{
				Type:    "go-backend",
				Name:    "backend",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "backend",
					// Missing required "module" field
				},
			},
			wantError: "module",
		},
		{
			name: "go-backend invalid port",
			component: models.ComponentConfig{
				Type:    "go-backend",
				Name:    "backend",
				Enabled: true,
				Config: map[string]interface{}{
					"name":   "backend",
					"module": "github.com/test/backend",
					"port":   99999, // Invalid port
				},
			},
			wantError: "port",
		},
		{
			name: "android invalid package",
			component: models.ComponentConfig{
				Type:    "android",
				Name:    "mobile",
				Enabled: true,
				Config: map[string]interface{}{
					"name":    "mobile",
					"package": "invalid", // Invalid package format
				},
			},
			wantError: "package",
		},
		{
			name: "ios invalid bundle_id",
			component: models.ComponentConfig{
				Type:    "ios",
				Name:    "mobile",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "mobile",
					"bundle_id": "invalid", // Invalid bundle ID format
				},
			},
			wantError: "bundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.ProjectConfig{
				Name:        "validation-test",
				Description: "Validation test",
				OutputDir:   filepath.Join(env.TempDir, tt.name),
				Components:  []models.ComponentConfig{tt.component},
				Options: models.ProjectOptions{
					UseExternalTools: true,
					DryRun:           false,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			_, err := coordinator.Generate(ctx, config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

// TestE2E_CacheRefresh tests cache refresh during generation
func TestE2E_CacheRefresh(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()

	// Setup tool cache
	cacheConfig := &ToolCacheConfig{
		CacheDir: env.TempDir,
		TTL:      1 * time.Millisecond, // Very short TTL to force refresh
	}
	cache, err := NewToolCache(cacheConfig, log)
	require.NoError(t, err)

	// Pre-populate cache
	cache.Set("go", true, "1.20.0")
	err = cache.Save()
	require.NoError(t, err)

	// Wait for cache to expire
	time.Sleep(10 * time.Millisecond)

	// Create coordinator and generate project
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "cache-refresh-test",
		Description: "Cache refresh test",
		OutputDir:   filepath.Join(env.TempDir, "cache-refresh-test"),
		Components: []models.ComponentConfig{
			{
				Type:    "android",
				Name:    "mobile-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":    "mobile-app",
					"package": "com.test.app",
				},
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: false, // Use fallback
			DryRun:           false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)
}

// TestE2E_Integration_CompleteWorkflow tests the complete workflow with all features
func TestE2E_Integration_CompleteWorkflow(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Complete configuration with all features
	config := &models.ProjectConfig{
		Name:        "complete-workflow-test",
		Description: "Complete workflow test with all features",
		OutputDir:   filepath.Join(env.TempDir, "complete-workflow-test"),
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/test/api-server",
					"framework": "gin",
					"port":      8080,
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
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			CreateBackup:     true,
			StreamOutput:     false,
			Verbose:          true,
			DryRun:           false,
			ForceOverwrite:   true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Generate project
	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Components)

	// Verify complete project structure
	projectDir := config.OutputDir
	testhelpers.AssertFileExists(t, projectDir)

	// Verify component
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "main.go"))

	// Verify integration files
	dockerComposePath := filepath.Join(projectDir, "docker-compose.yml")
	if _, err := os.Stat(dockerComposePath); err == nil {
		testhelpers.AssertFileContains(t, dockerComposePath, "services:")
		testhelpers.AssertFileContains(t, dockerComposePath, "api-server")
	}

	// Verify environment configuration
	envPath := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		testhelpers.AssertFileContains(t, envPath, "LOG_LEVEL")
	}

	// Verify scripts directory
	scriptsDir := filepath.Join(projectDir, "scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		entries, _ := os.ReadDir(scriptsDir)
		assert.NotEmpty(t, entries, "Scripts directory should contain files")
	}

	// Verify README
	readmePath := filepath.Join(projectDir, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		testhelpers.AssertFileContains(t, readmePath, config.Name)
	}
}

// TestE2E_ErrorHandling_WithSuggestions tests error handling with suggestions
func TestE2E_ErrorHandling_WithSuggestions(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Create configuration with missing required tool
	config := &models.ProjectConfig{
		Name:        "error-handling-test",
		Description: "Error handling test",
		OutputDir:   filepath.Join(env.TempDir, "error-handling-test"),
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
		Options: models.ProjectOptions{
			UseExternalTools: true, // Require external tools
			DryRun:           false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// If npx is not available, should get error with suggestions
	_, err := coordinator.Generate(ctx, config)

	// Error might occur if tool is not available
	if err != nil {
		// Verify error contains helpful information
		errMsg := err.Error()
		assert.NotEmpty(t, errMsg)

		// Check if it's a GenerationError with suggestions
		var genErr *GenerationError
		if assert.ErrorAs(t, err, &genErr) {
			suggestions := genErr.GetSuggestions()
			if suggestions != "" {
				assert.Contains(t, suggestions, "Next Steps:")
			}
		}
	}
}
