package mapper

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentMapper(t *testing.T) {
	t.Run("GetMapping returns correct path for known component", func(t *testing.T) {
		mapper := NewComponentMapper()

		path, err := mapper.GetMapping("nextjs")
		require.NoError(t, err)
		assert.Equal(t, "App/", path)

		path, err = mapper.GetMapping("go-backend")
		require.NoError(t, err)
		assert.Equal(t, "CommonServer/", path)
	})

	t.Run("GetMapping returns error for unknown component", func(t *testing.T) {
		mapper := NewComponentMapper()

		_, err := mapper.GetMapping("unknown")
		assert.Error(t, err)
	})

	t.Run("RegisterMapping adds new mapping", func(t *testing.T) {
		mapper := NewComponentMapper()

		err := mapper.RegisterMapping("custom", "Custom/")
		require.NoError(t, err)

		path, err := mapper.GetMapping("custom")
		require.NoError(t, err)
		assert.Equal(t, "Custom/", path)
	})

	t.Run("ListMappings returns all mappings", func(t *testing.T) {
		mapper := NewComponentMapper()

		mappings := mapper.ListMappings()
		assert.NotEmpty(t, mappings)
		assert.Contains(t, mappings, "nextjs")
		assert.Contains(t, mappings, "go-backend")
	})
}

func TestStructureMapper_GetTargetPath(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper)

	t.Run("returns correct target path", func(t *testing.T) {
		path := mapper.GetTargetPath("nextjs")
		assert.Equal(t, "App/", path)
	})

	t.Run("returns empty string for unknown component", func(t *testing.T) {
		path := mapper.GetTargetPath("unknown")
		assert.Equal(t, "", path)
	})
}

func TestStructureMapper_Map(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper)
	ctx := context.Background()

	t.Run("moves directory to target location", func(t *testing.T) {
		// Create temporary directories
		tempDir := t.TempDir()
		sourceDir := filepath.Join(tempDir, "source")
		targetRoot := filepath.Join(tempDir, "target")

		// Create source directory with a file
		err := os.MkdirAll(sourceDir, 0755)
		require.NoError(t, err)

		testFile := filepath.Join(sourceDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)

		// Map the directory
		err = mapper.Map(ctx, sourceDir, targetRoot, "nextjs")
		require.NoError(t, err)

		// Verify target exists
		targetPath := filepath.Join(targetRoot, "App")
		_, err = os.Stat(targetPath)
		assert.NoError(t, err)

		// Verify file was moved
		movedFile := filepath.Join(targetPath, "test.txt")
		content, err := os.ReadFile(movedFile)
		require.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		// Verify source was removed
		_, err = os.Stat(sourceDir)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("returns error for non-existent source", func(t *testing.T) {
		tempDir := t.TempDir()
		sourceDir := filepath.Join(tempDir, "nonexistent")
		targetRoot := filepath.Join(tempDir, "target")

		err := mapper.Map(ctx, sourceDir, targetRoot, "nextjs")
		assert.Error(t, err)
	})

	t.Run("returns error for unknown component type", func(t *testing.T) {
		tempDir := t.TempDir()
		sourceDir := filepath.Join(tempDir, "source")
		targetRoot := filepath.Join(tempDir, "target")

		err := os.MkdirAll(sourceDir, 0755)
		require.NoError(t, err)

		err = mapper.Map(ctx, sourceDir, targetRoot, "unknown")
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateStructure(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper)

	t.Run("validates existing structure", func(t *testing.T) {
		// Create temporary directory with expected structure
		tempDir := t.TempDir()

		// Create a valid Next.js structure
		appDir := filepath.Join(tempDir, "App")
		err := os.MkdirAll(filepath.Join(appDir, "app"), 0755)
		require.NoError(t, err)

		packageJSON := filepath.Join(appDir, "package.json")
		err = os.WriteFile(packageJSON, []byte(`{"name": "test"}`), 0644)
		require.NoError(t, err)

		nextConfig := filepath.Join(appDir, "next.config.js")
		err = os.WriteFile(nextConfig, []byte(`module.exports = {}`), 0644)
		require.NoError(t, err)

		// Create a valid Go structure
		serverDir := filepath.Join(tempDir, "CommonServer")
		err = os.MkdirAll(serverDir, 0755)
		require.NoError(t, err)

		goMod := filepath.Join(serverDir, "go.mod")
		err = os.WriteFile(goMod, []byte(`module test`), 0644)
		require.NoError(t, err)

		mainGo := filepath.Join(serverDir, "main.go")
		err = os.WriteFile(mainGo, []byte(`package main`), 0644)
		require.NoError(t, err)

		// Validate - should not error
		err = mapper.ValidateStructure(tempDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for non-existent root", func(t *testing.T) {
		err := mapper.ValidateStructure("/nonexistent/path")
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateNextJSStructure(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper).(*StructureMapper)

	t.Run("validates Next.js with app directory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create package.json
		packageJSON := filepath.Join(tempDir, "package.json")
		err := os.WriteFile(packageJSON, []byte(`{"name": "test"}`), 0644)
		require.NoError(t, err)

		// Create next.config.js
		nextConfig := filepath.Join(tempDir, "next.config.js")
		err = os.WriteFile(nextConfig, []byte(`module.exports = {}`), 0644)
		require.NoError(t, err)

		// Create app directory
		appDir := filepath.Join(tempDir, "app")
		err = os.MkdirAll(appDir, 0755)
		require.NoError(t, err)

		err = mapper.validateNextJSStructure(tempDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for missing essential files", func(t *testing.T) {
		tempDir := t.TempDir()

		err := mapper.validateNextJSStructure(tempDir)
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateGoStructure(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper).(*StructureMapper)

	t.Run("validates Go project with main.go", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create go.mod
		goMod := filepath.Join(tempDir, "go.mod")
		err := os.WriteFile(goMod, []byte(`module test`), 0644)
		require.NoError(t, err)

		// Create main.go
		mainGo := filepath.Join(tempDir, "main.go")
		err = os.WriteFile(mainGo, []byte(`package main`), 0644)
		require.NoError(t, err)

		err = mapper.validateGoStructure(tempDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for missing go.mod", func(t *testing.T) {
		tempDir := t.TempDir()

		err := mapper.validateGoStructure(tempDir)
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateAndroidStructure(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper).(*StructureMapper)

	t.Run("validates Android project", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create build.gradle
		buildGradle := filepath.Join(tempDir, "build.gradle")
		err := os.WriteFile(buildGradle, []byte(`// build file`), 0644)
		require.NoError(t, err)

		// Create settings.gradle
		settingsGradle := filepath.Join(tempDir, "settings.gradle")
		err = os.WriteFile(settingsGradle, []byte(`// settings`), 0644)
		require.NoError(t, err)

		// Create app directory
		appDir := filepath.Join(tempDir, "app")
		err = os.MkdirAll(appDir, 0755)
		require.NoError(t, err)

		err = mapper.validateAndroidStructure(tempDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for missing build.gradle", func(t *testing.T) {
		tempDir := t.TempDir()

		err := mapper.validateAndroidStructure(tempDir)
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateiOSStructure(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper).(*StructureMapper)

	t.Run("validates iOS project with Package.swift", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create Package.swift
		packageSwift := filepath.Join(tempDir, "Package.swift")
		err := os.WriteFile(packageSwift, []byte(`// swift package`), 0644)
		require.NoError(t, err)

		err = mapper.validateiOSStructure(tempDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for missing project files", func(t *testing.T) {
		tempDir := t.TempDir()

		err := mapper.validateiOSStructure(tempDir)
		assert.Error(t, err)
	})
}

func TestStructureMapper_ValidateStructureDetailed(t *testing.T) {
	componentMapper := NewComponentMapper()
	mapper := NewStructureMapper(componentMapper).(*StructureMapper)

	t.Run("returns detailed validation results", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a valid Next.js structure
		appDir := filepath.Join(tempDir, "App")
		err := os.MkdirAll(filepath.Join(appDir, "app"), 0755)
		require.NoError(t, err)

		packageJSON := filepath.Join(appDir, "package.json")
		err = os.WriteFile(packageJSON, []byte(`{"name": "test"}`), 0644)
		require.NoError(t, err)

		nextConfig := filepath.Join(appDir, "next.config.js")
		err = os.WriteFile(nextConfig, []byte(`module.exports = {}`), 0644)
		require.NoError(t, err)

		result := mapper.ValidateStructureDetailed(tempDir, []string{"nextjs"})
		assert.NotNil(t, result)
		assert.True(t, result.Valid)
		assert.Contains(t, result.Components, "nextjs")
		assert.True(t, result.Components["nextjs"].Valid)
	})
}
