package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// TestCompleteProjectGenerationIntegration tests the complete project generation workflow
// This test demonstrates the full functionality of task 5.2
func TestCompleteProjectGenerationIntegration(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a comprehensive project configuration
	config := &models.ProjectConfig{
		Name:         "awesome-project",
		Organization: "awesome-org",
		Description:  "An awesome open source project with all components",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@awesome-org.com",
		Repository:   "https://github.com/awesome-org/awesome-project",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    true,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Node:      "18.17.0",
			Go:        "1.22.0",
			Kotlin:    "1.9.0",
			Swift:     "5.9.0",
			NextJS:    "14.0.0",
			React:     "18.2.0",
			UpdatedAt: time.Now(),
		},
		OutputPath:       tempDir,
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}

	// Initialize the project generator
	pg := NewProjectGenerator()

	t.Run("Step 1: Generate complete project directory structure", func(t *testing.T) {
		err := pg.GenerateProjectStructure(config, tempDir)
		if err != nil {
			t.Fatalf("Failed to generate project structure: %v", err)
		}

		// Verify the project root directory was created
		projectPath := filepath.Join(tempDir, config.Name)
		if !pg.fsGen.FileExists(projectPath) {
			t.Fatal("Project root directory was not created")
		}

		t.Logf("✓ Project structure generated successfully at: %s", projectPath)
	})

	t.Run("Step 2: Generate component-specific files", func(t *testing.T) {
		err := pg.GenerateComponentFiles(config, tempDir)
		if err != nil {
			t.Fatalf("Failed to generate component files: %v", err)
		}

		projectPath := filepath.Join(tempDir, config.Name)

		// Verify key files were created
		keyFiles := []string{
			"Makefile",
			"README.md",
			"CONTRIBUTING.md",
			"SECURITY.md",
			"docker-compose.yml",
			".gitignore",
			"App/package.json",
			"Home/package.json",
			"Admin/package.json",
			"CommonServer/go.mod",
			"Mobile/Android/build.gradle",
			"Mobile/iOS/Package.swift",
			"Deploy/terraform/main.tf",
			".github/workflows/ci.yml",
			".github/workflows/security.yml",
			".github/dependabot.yml",
		}

		for _, file := range keyFiles {
			filePath := filepath.Join(projectPath, file)
			if !pg.fsGen.FileExists(filePath) {
				t.Errorf("Expected file not created: %s", file)
			}
		}

		t.Logf("✓ Component files generated successfully")
	})

	t.Run("Step 3: Validate project structure", func(t *testing.T) {
		projectPath := filepath.Join(tempDir, config.Name)

		err := pg.ValidateProjectStructure(projectPath, config)
		if err != nil {
			t.Fatalf("Project structure validation failed: %v", err)
		}

		t.Logf("✓ Project structure validation passed")
	})

	t.Run("Step 4: Validate cross-references between generated files", func(t *testing.T) {
		projectPath := filepath.Join(tempDir, config.Name)

		err := pg.ValidateCrossReferences(projectPath, config)
		if err != nil {
			t.Fatalf("Cross-reference validation failed: %v", err)
		}

		t.Logf("✓ Cross-reference validation passed")
	})

	t.Run("Step 5: Verify complete project structure", func(t *testing.T) {
		projectPath := filepath.Join(tempDir, config.Name)

		// Count total files and directories created
		fileCount := 0
		dirCount := 0

		err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
			return nil
		})

		if err != nil {
			t.Fatalf("Failed to walk project directory: %v", err)
		}

		t.Logf("✓ Project generation complete:")
		t.Logf("  - Project name: %s", config.Name)
		t.Logf("  - Organization: %s", config.Organization)
		t.Logf("  - Total directories created: %d", dirCount)
		t.Logf("  - Total files created: %d", fileCount)
		t.Logf("  - Project path: %s", projectPath)

		// Verify minimum expected structure
		if dirCount < 20 {
			t.Errorf("Expected at least 20 directories, got %d", dirCount)
		}

		if fileCount < 15 {
			t.Errorf("Expected at least 15 files, got %d", fileCount)
		}
	})

	t.Run("Step 6: Verify component-specific generation based on selection", func(t *testing.T) {
		projectPath := filepath.Join(tempDir, config.Name)

		// Test that only selected components were generated
		componentTests := map[string]struct {
			enabled   bool
			checkPath string
		}{
			"Frontend Main App": {
				enabled:   config.Components.Frontend.MainApp,
				checkPath: "App/package.json",
			},
			"Frontend Home": {
				enabled:   config.Components.Frontend.Home,
				checkPath: "Home/package.json",
			},
			"Frontend Admin": {
				enabled:   config.Components.Frontend.Admin,
				checkPath: "Admin/package.json",
			},
			"Backend API": {
				enabled:   config.Components.Backend.API,
				checkPath: "CommonServer/go.mod",
			},
			"Mobile Android": {
				enabled:   config.Components.Mobile.Android,
				checkPath: "Mobile/Android/build.gradle",
			},
			"Mobile iOS": {
				enabled:   config.Components.Mobile.IOS,
				checkPath: "Mobile/iOS/Package.swift",
			},
			"Infrastructure Terraform": {
				enabled:   config.Components.Infrastructure.Terraform,
				checkPath: "Deploy/terraform/main.tf",
			},
		}

		for componentName, test := range componentTests {
			fullPath := filepath.Join(projectPath, test.checkPath)
			exists := pg.fsGen.FileExists(fullPath)

			if test.enabled && !exists {
				t.Errorf("Component %s is enabled but file %s does not exist", componentName, test.checkPath)
			} else if test.enabled && exists {
				t.Logf("✓ Component %s: file %s exists", componentName, test.checkPath)
			}
		}
	})
}

// TestProjectGenerationWithMinimalComponents tests generation with minimal component selection
func TestProjectGenerationWithMinimalComponents(t *testing.T) {
	tempDir := t.TempDir()

	// Create a minimal configuration (only frontend main app)
	config := &models.ProjectConfig{
		Name:         "minimal-project",
		Organization: "test-org",
		Description:  "A minimal project with only frontend",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   false,
			},
			Backend: models.BackendComponents{
				API: false,
			},
			Mobile: models.MobileComponents{
				Android: false,
				IOS:     false,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  false,
				Kubernetes: false,
				Docker:     false,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "18.17.0",
			Go:     "1.22.0",
			NextJS: "14.0.0",
			React:  "18.2.0",
		},
		OutputPath: tempDir,
	}

	pg := NewProjectGenerator()

	// Generate project
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Failed to generate project structure: %v", err)
	}

	if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
		t.Fatalf("Failed to generate component files: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	// Validate structure
	if err := pg.ValidateProjectStructure(projectPath, config); err != nil {
		t.Fatalf("Project structure validation failed: %v", err)
	}

	if err := pg.ValidateCrossReferences(projectPath, config); err != nil {
		t.Fatalf("Cross-reference validation failed: %v", err)
	}

	// Verify that only selected components exist
	shouldExist := []string{
		"App/package.json",
		"App/next.config.js",
		"Makefile",
		"README.md",
	}

	shouldNotExist := []string{
		"Home/package.json",
		"Admin/package.json",
		"CommonServer/go.mod",
		"Mobile/Android/build.gradle",
		"Mobile/iOS/Package.swift",
		"Deploy/terraform/main.tf",
	}

	for _, file := range shouldExist {
		filePath := filepath.Join(projectPath, file)
		if !pg.fsGen.FileExists(filePath) {
			t.Errorf("Expected file does not exist: %s", file)
		}
	}

	for _, file := range shouldNotExist {
		filePath := filepath.Join(projectPath, file)
		if pg.fsGen.FileExists(filePath) {
			t.Errorf("Unexpected file exists: %s", file)
		}
	}

	t.Logf("✓ Minimal project generation completed successfully")
}

// TestProjectGenerationDryRun tests the dry-run functionality
func TestProjectGenerationDryRun(t *testing.T) {
	tempDir := t.TempDir()

	config := &models.ProjectConfig{
		Name:         "dry-run-project",
		Organization: "test-org",
		Description:  "A project for testing dry-run mode",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{MainApp: true},
			Backend:  models.BackendComponents{API: true},
		},
		Versions: &models.VersionConfig{
			Node:   "18.17.0",
			Go:     "1.22.0",
			NextJS: "14.0.0",
			React:  "18.2.0",
		},
		OutputPath: tempDir,
	}

	// Use dry-run generator
	pg := NewDryRunProjectGenerator()

	// These operations should succeed without creating actual files
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Dry-run project structure generation failed: %v", err)
	}

	if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
		t.Fatalf("Dry-run component files generation failed: %v", err)
	}

	// Verify that no files were actually created
	projectPath := filepath.Join(tempDir, config.Name)
	if pg.fsGen.FileExists(projectPath) {
		t.Errorf("Dry-run mode should not create actual files, but project directory exists")
	}

	t.Logf("✓ Dry-run mode completed successfully without creating files")
}
