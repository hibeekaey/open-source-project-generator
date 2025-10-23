package orchestrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCoordinator_ValidateFinalStructure(t *testing.T) {
	log := logger.NewLogger()
	coordinator := NewProjectCoordinator(log)

	t.Run("validates successful structure with checkmark", func(t *testing.T) {
		// Create temporary directory with valid structure
		tempDir := t.TempDir()

		// Create valid Next.js structure
		appDir := filepath.Join(tempDir, "App")
		err := os.MkdirAll(filepath.Join(appDir, "app"), 0755)
		require.NoError(t, err)

		packageJSON := filepath.Join(appDir, "package.json")
		err = os.WriteFile(packageJSON, []byte(`{"name": "test"}`), 0644)
		require.NoError(t, err)

		nextConfig := filepath.Join(appDir, "next.config.js")
		err = os.WriteFile(nextConfig, []byte(`module.exports = {}`), 0644)
		require.NoError(t, err)

		// Create valid Go structure
		serverDir := filepath.Join(tempDir, "CommonServer")
		err = os.MkdirAll(serverDir, 0755)
		require.NoError(t, err)

		goMod := filepath.Join(serverDir, "go.mod")
		err = os.WriteFile(goMod, []byte(`module test`), 0644)
		require.NoError(t, err)

		mainGo := filepath.Join(serverDir, "main.go")
		err = os.WriteFile(mainGo, []byte(`package main`), 0644)
		require.NoError(t, err)

		// Create config and results
		config := &models.ProjectConfig{
			OutputDir: tempDir,
		}

		results := []*models.ComponentResult{
			{
				Type:    "nextjs",
				Name:    "frontend",
				Success: true,
			},
			{
				Type:    "go-backend",
				Name:    "backend",
				Success: true,
			},
		}

		// Validate structure
		err = coordinator.validateFinalStructure(config, results)
		assert.NoError(t, err)
	})

	t.Run("reports missing files and directories", func(t *testing.T) {
		// Create temporary directory with incomplete structure
		tempDir := t.TempDir()

		// Create only App directory without required files
		appDir := filepath.Join(tempDir, "App")
		err := os.MkdirAll(appDir, 0755)
		require.NoError(t, err)

		// Create config and results
		config := &models.ProjectConfig{
			OutputDir: tempDir,
		}

		results := []*models.ComponentResult{
			{
				Type:    "nextjs",
				Name:    "frontend",
				Success: true,
			},
		}

		// Validate structure - should fail
		err = coordinator.validateFinalStructure(config, results)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "structure validation failed")
	})

	t.Run("skips validation for failed components", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create config and results with failed component
		config := &models.ProjectConfig{
			OutputDir: tempDir,
		}

		results := []*models.ComponentResult{
			{
				Type:    "nextjs",
				Name:    "frontend",
				Success: false, // Failed component
			},
		}

		// Validate structure - should pass because failed components are skipped
		err := coordinator.validateFinalStructure(config, results)
		assert.NoError(t, err)
	})
}
