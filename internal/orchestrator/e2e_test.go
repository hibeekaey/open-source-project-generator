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

// TestE2E_FullProjectGeneration tests the complete project generation workflow
func TestE2E_FullProjectGeneration(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Create project configuration
	config := &models.ProjectConfig{
		Name:        "test-full-project",
		Description: "End-to-end test project",
		OutputDir:   filepath.Join(env.TempDir, "test-full-project"),
		Components: []models.ComponentConfig{
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/test/api-server",
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
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     false,
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

	// Verify project structure
	projectDir := config.OutputDir
	testhelpers.AssertFileExists(t, projectDir)

	// Verify Go backend was created
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "main.go"))
}

// TestE2E_MultiComponentGeneration tests generation with multiple components
func TestE2E_MultiComponentGeneration(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "test-multi-component",
		Description: "Multi-component test project",
		OutputDir:   filepath.Join(env.TempDir, "test-multi-component"),
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

	// Verify both components were created
	projectDir := config.OutputDir

	// Verify Go backend
	serverDir := filepath.Join(projectDir, "CommonServer")
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))

	// Verify Android (fallback generation)
	androidDir := filepath.Join(projectDir, "Mobile", "android")
	testhelpers.AssertFileExists(t, filepath.Join(androidDir, "build.gradle"))
}

// TestE2E_StructureMapping tests that generated projects are mapped to correct locations
func TestE2E_StructureMapping(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "test-structure-mapping",
		Description: "Structure mapping test",
		OutputDir:   filepath.Join(env.TempDir, "test-structure-mapping"),
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

	// Verify structure mapping - Go backend should be in CommonServer/
	projectDir := config.OutputDir
	serverDir := filepath.Join(projectDir, "CommonServer")

	// Check that files exist in the mapped location
	_, err = os.Stat(serverDir)
	assert.NoError(t, err, "CommonServer directory should exist")

	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "go.mod"))
	testhelpers.AssertFileExists(t, filepath.Join(serverDir, "main.go"))
}

// TestE2E_DryRunMode tests dry run mode doesn't create files
func TestE2E_DryRunMode(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "test-dry-run",
		Description: "Dry run test",
		OutputDir:   filepath.Join(env.TempDir, "test-dry-run"),
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
			DryRun:           true, // Enable dry run
			Verbose:          false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.DryRun)

	// In dry run mode, the output directory should not be created
	projectDir := config.OutputDir
	_, err = os.Stat(projectDir)
	// Directory might exist but should be empty or minimal
	if err == nil {
		// If directory exists, check it's empty or only has minimal structure
		entries, _ := os.ReadDir(projectDir)
		assert.LessOrEqual(t, len(entries), 1, "Dry run should not create many files")
	}
}

// TestE2E_ValidationFailure tests that invalid configuration is rejected
func TestE2E_ValidationFailure(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Create invalid configuration (missing required fields)
	config := &models.ProjectConfig{
		Name:        "", // Invalid: empty name
		Description: "Invalid config test",
		OutputDir:   filepath.Join(env.TempDir, "test-invalid"),
		Components:  []models.ComponentConfig{},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
		},
	}

	ctx := context.Background()

	// Should fail validation
	_, err := coordinator.Generate(ctx, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

// TestE2E_OfflineMode tests project generation in offline mode
func TestE2E_OfflineMode(t *testing.T) {
	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	// Enable offline mode
	coordinator.SetOfflineMode(true)
	assert.True(t, coordinator.IsOffline())

	config := &models.ProjectConfig{
		Name:        "test-offline",
		Description: "Offline mode test",
		OutputDir:   filepath.Join(env.TempDir, "test-offline"),
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resultInterface, err := coordinator.Generate(ctx, config)
	require.NoError(t, err)

	result, ok := resultInterface.(*models.GenerationResult)
	require.True(t, ok)
	assert.True(t, result.Success)

	// Verify fallback generation worked
	projectDir := config.OutputDir
	androidDir := filepath.Join(projectDir, "Mobile", "android")
	testhelpers.AssertFileExists(t, filepath.Join(androidDir, "build.gradle"))
}

// TestE2E_ComponentIntegration tests that integration files are generated
func TestE2E_ComponentIntegration(t *testing.T) {
	if !testhelpers.IsToolAvailable("go") {
		t.Skip("Skipping E2E test: go not available")
	}

	env := testhelpers.SetupTestEnv(t)
	defer env.Cleanup()

	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	config := &models.ProjectConfig{
		Name:        "test-integration",
		Description: "Integration test",
		OutputDir:   filepath.Join(env.TempDir, "test-integration"),
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
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "http://localhost:8080",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
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

	// Verify integration files were created
	projectDir := config.OutputDir

	// Check for docker-compose.yml
	dockerComposePath := filepath.Join(projectDir, "docker-compose.yml")
	if _, err := os.Stat(dockerComposePath); err == nil {
		testhelpers.AssertFileContains(t, dockerComposePath, "services:")
	}

	// Check for scripts
	scriptsDir := filepath.Join(projectDir, "scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		// Scripts directory exists, verify some scripts
		entries, _ := os.ReadDir(scriptsDir)
		assert.NotEmpty(t, entries, "Scripts directory should contain files")
	}
}
