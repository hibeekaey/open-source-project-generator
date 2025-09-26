package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestStandardizedStructureGenerator_GenerateStandardizedStructure(t *testing.T) {
	// Skip test if templates directory doesn't exist
	if _, err := os.Stat("templates/base"); os.IsNotExist(err) {
		t.Skip("Skipping test: templates directory not found")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "standardized_structure_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "testorg",
		Description:  "Test project description",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/testorg/test-project",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   true,
					Admin:  true,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "18.0.0",
			Go:   "1.22.0",
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.0.0",
			},
		},
	}

	// Create generator
	generator := NewStandardizedStructureGenerator()

	// Generate structure
	err = generator.GenerateStandardizedStructure(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to generate standardized structure: %v", err)
	}

	// Verify project root directory exists
	projectPath := filepath.Join(tempDir, config.Name)
	if !generator.fsGen.FileExists(projectPath) {
		t.Errorf("Project root directory not created: %s", projectPath)
	}

	// Test that base template files were created
	t.Run("Base Template Files", func(t *testing.T) {
		testBaseTemplateFiles(t, projectPath, generator)
	})

	// Test component selection logic
	t.Run("Component Selection", func(t *testing.T) {
		testComponentSelection(t, generator, config)
	})
}

func testBaseTemplateFiles(t *testing.T, projectPath string, generator *StandardizedStructureGenerator) {
	// Test that base template files were created (these come from the base template)
	baseFiles := []string{
		"README.md",
		"CONTRIBUTING.md",
		"LICENSE",
		".gitignore",
		"Makefile",
	}

	for _, file := range baseFiles {
		filePath := filepath.Join(projectPath, file)
		if !generator.fsGen.FileExists(filePath) {
			t.Errorf("Base template file not created: %s", file)
		}
	}

	// Test that GitHub workflows directory was created
	githubDir := filepath.Join(projectPath, ".github/workflows")
	if !generator.fsGen.FileExists(githubDir) {
		t.Errorf("GitHub workflows directory not created")
	}
}

func testComponentSelection(t *testing.T, generator *StandardizedStructureGenerator, config *models.ProjectConfig) {
	// Test helper methods for component selection
	if !generator.hasFrontendComponents(config) {
		t.Errorf("hasFrontendComponents should return true when frontend components are selected")
	}

	if !generator.hasMobileComponents(config) {
		t.Errorf("hasMobileComponents should return true when mobile components are selected")
	}

	if !generator.hasInfrastructureComponents(config) {
		t.Errorf("hasInfrastructureComponents should return true when infrastructure components are selected")
	}

	// Test with empty config
	emptyConfig := &models.ProjectConfig{}
	if generator.hasFrontendComponents(emptyConfig) {
		t.Errorf("hasFrontendComponents should return false when no frontend components are selected")
	}

	if generator.hasMobileComponents(emptyConfig) {
		t.Errorf("hasMobileComponents should return false when no mobile components are selected")
	}

	if generator.hasInfrastructureComponents(emptyConfig) {
		t.Errorf("hasInfrastructureComponents should return false when no infrastructure components are selected")
	}
}

func TestStandardizedStructureGenerator_GenerateStandardProjectFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "standardized_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "testorg",
		Description:  "Test project description",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/testorg/test-project",
	}

	// Create generator
	generator := NewStandardizedStructureGenerator()

	// Create project directory
	projectPath := filepath.Join(tempDir, config.Name)
	err = generator.fsGen.EnsureDirectory(projectPath)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Generate standard project files (now handled by templates)
	err = generator.GenerateStandardProjectFiles(projectPath, config)
	if err != nil {
		t.Fatalf("Failed to generate standard project files: %v", err)
	}

	// The method should complete without error (files are generated by templates)
}

func TestStandardizedStructureGenerator_ConditionalGeneration(t *testing.T) {
	// Skip test if templates directory doesn't exist
	if _, err := os.Stat("templates/base"); os.IsNotExist(err) {
		t.Skip("Skipping test: templates directory not found")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "conditional_generation_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to remove temp directory: %v", err)
		}
	}()

	// Test with only frontend components
	config := &models.ProjectConfig{
		Name:         "frontend-only-project",
		Organization: "testorg",
		Description:  "Frontend only test project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	// Create generator
	generator := NewStandardizedStructureGenerator()

	// Generate structure
	err = generator.GenerateStandardizedStructure(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to generate standardized structure: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	// Verify project directory exists
	if !generator.fsGen.FileExists(projectPath) {
		t.Errorf("Project directory not created")
	}

	// Test component selection logic
	if !generator.hasFrontendComponents(config) {
		t.Errorf("Should detect frontend components")
	}

	if generator.hasMobileComponents(config) {
		t.Errorf("Should not detect mobile components")
	}

	if generator.hasInfrastructureComponents(config) {
		t.Errorf("Should not detect infrastructure components")
	}
}
