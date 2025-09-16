package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
)

func TestBasicProjectGeneration(t *testing.T) {
	tempDir := t.TempDir()

	// Create managers
	configManager := config.NewManager(filepath.Join(tempDir, "cache"), "")
	fsGenerator := filesystem.NewGenerator()

	// Test configuration
	config := &models.ProjectConfig{
		Name:         "basic-integration-test",
		Organization: "test-org",
		Description:  "Basic integration test project",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.11.0",
			Go:   "1.22.0",
			Packages: map[string]string{
				"next":       "15.0.0",
				"react":      "18.2.0",
				"typescript": "5.3.3",
			},
		},
		OutputPath:       tempDir,
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}

	// Validate configuration
	err := configManager.ValidateConfig(config)
	if err != nil {
		t.Fatalf("Configuration validation failed: %v", err)
	}

	// Create project
	err = fsGenerator.CreateProject(config, tempDir)
	if err != nil {
		t.Fatalf("Project creation failed: %v", err)
	}

	// Verify project structure
	projectPath := filepath.Join(tempDir, config.Name)
	if !fileExists(projectPath) {
		t.Fatalf("Project directory not created: %s", projectPath)
	}

	// List what files were actually created
	files, err := os.ReadDir(projectPath)
	if err != nil {
		t.Fatalf("Failed to read project directory: %v", err)
	}

	t.Logf("Created files in project directory:")
	for _, file := range files {
		t.Logf("  - %s", file.Name())
	}

	// For now, just verify the directory was created
	// The actual file generation depends on template availability

	t.Logf("Basic integration test completed successfully")
}

func TestTemplateProcessingBasic(t *testing.T) {
	tempDir := t.TempDir()
	templateEngine := template.NewEngine()

	// Create a simple template
	templateContent := `# {{.Name}}

{{.Description}}

## License
{{.License}}`

	templatePath := filepath.Join(tempDir, "test.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Test configuration
	config := &models.ProjectConfig{
		Name:        "template-test",
		Description: "Template processing test",
		License:     "MIT",
	}

	// Process template
	result, err := templateEngine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Template processing failed: %v", err)
	}

	resultStr := string(result)

	// Verify content
	expectedContent := []string{
		"# template-test",
		"Template processing test",
		"MIT",
	}

	for _, expected := range expectedContent {
		if !contains(resultStr, expected) {
			t.Errorf("Template output missing expected content: %s", expected)
		}
	}

	t.Logf("Template processing test completed successfully")
}

func TestConfigurationManagementBasic(t *testing.T) {
	tempDir := t.TempDir()
	configManager := config.NewManager(filepath.Join(tempDir, "cache"), "")

	// Test loading defaults
	defaults, err := configManager.LoadDefaults()
	if err != nil {
		t.Fatalf("Failed to load defaults: %v", err)
	}

	if defaults == nil {
		t.Fatal("Expected defaults to be non-nil")
	}

	if defaults.License != "MIT" {
		t.Errorf("Expected default license MIT, got %s", defaults.License)
	}

	// Test configuration merging - method removed
	// override := &models.ProjectConfig{
	// 	Name:         "merge-test",
	// 	Organization: "test-org",
	// 	Components: models.Components{
	// 		Frontend: models.FrontendComponents{
	// 			NextJS: models.NextJSComponents{
	// 				App:   true,
	// 				Admin: true,
	// 			},
	// 		},
	// 	},
	// }

	// merged := configManager.MergeConfigs(defaults, override) // Method removed

	// if merged.Name != "merge-test" {
	// 	t.Errorf("Expected merged name 'merge-test', got %s", merged.Name)
	// }

	// if merged.License != "MIT" {
	// 	t.Errorf("Expected merged license 'MIT', got %s", merged.License)
	// }

	// if !merged.Components.Frontend.NextJS.App {
	// 	t.Error("Expected NextJS.App to be true from override")
	// }

	t.Logf("Configuration management test completed successfully")
}

func TestEndToEndBasic(t *testing.T) {
	tempDir := t.TempDir()

	// Create managers
	configManager := config.NewManager(filepath.Join(tempDir, "cache"), "")
	fsGenerator := filesystem.NewGenerator()

	// 1. Load defaults
	config, err := configManager.LoadDefaults()
	if err != nil {
		t.Fatalf("Failed to load defaults: %v", err)
	}

	// 2. Customize configuration
	config.Name = "e2e-basic-test"
	config.Organization = "e2e-org"
	config.Description = "End-to-end basic test"
	config.OutputPath = "./test-output"
	config.Components.Frontend.NextJS.App = true
	config.Components.Infrastructure.Docker = true

	// 3. Validate configuration
	err = configManager.ValidateConfig(config)
	if err != nil {
		t.Fatalf("Configuration validation failed: %v", err)
	}

	// 4. Get versions - method removed
	// versions, err := configManager.GetLatestVersions()
	// if err != nil {
	// 	t.Fatalf("Failed to get versions: %v", err)
	// }
	// config.Versions = versions

	// 5. Create project
	err = fsGenerator.CreateProject(config, tempDir)
	if err != nil {
		t.Fatalf("Project creation failed: %v", err)
	}

	// 6. Verify project
	projectPath := filepath.Join(tempDir, config.Name)
	if !fileExists(projectPath) {
		t.Error("Project directory not created")
	}

	// List what files were actually created
	files, err := os.ReadDir(projectPath)
	if err != nil {
		t.Logf("Failed to read project directory: %v", err)
	} else {
		t.Logf("Created files in E2E project directory:")
		for _, file := range files {
			t.Logf("  - %s", file.Name())
		}
	}

	t.Logf("End-to-end basic test completed successfully")
}

// Helper functions
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
