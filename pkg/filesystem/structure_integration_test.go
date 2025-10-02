package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestStructureManagerIntegrationWithProjectGenerator tests that the StructureManager
// integrates correctly with the ProjectGenerator
func TestStructureManagerIntegrationWithProjectGenerator(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "structure_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a project generator (which uses StructureManager internally)
	pg := NewProjectGenerator()

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "integration-test-project",
		Organization: "test-org",
		Description:  "Integration test project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Admin: true,
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
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Go: "1.22",
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.0.0",
			},
		},
	}

	// Test project structure generation
	err = pg.GenerateProjectStructure(config, tempDir)
	if err != nil {
		t.Fatalf("GenerateProjectStructure failed: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	// Test project structure validation
	err = pg.ValidateProjectStructure(projectPath, config)
	if err != nil {
		t.Fatalf("ValidateProjectStructure failed: %v", err)
	}

	// Verify that the StructureManager is working correctly by checking specific directories
	structure := pg.structureManager.GetStructure()

	// Check that root directories exist
	for _, dir := range structure.RootDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Root directory not created: %s", dir)
		}
	}

	// Check that frontend directories exist (since frontend is enabled)
	for _, dir := range structure.FrontendDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Frontend directory not created: %s", dir)
		}
	}

	// Check that backend directories exist (since backend is enabled)
	for _, dir := range structure.BackendDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Backend directory not created: %s", dir)
		}
	}

	// Check that mobile directories exist (since mobile is enabled)
	for _, dir := range structure.MobileDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Mobile directory not created: %s", dir)
		}
	}

	// Check that infrastructure directories exist (since infrastructure is enabled)
	for _, dir := range structure.InfraDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Infrastructure directory not created: %s", dir)
		}
	}

	t.Logf("✓ Integration test passed: StructureManager works correctly with ProjectGenerator")
}

// TestStructureManagerCustomizationIntegration tests that structure customization
// works correctly with the ProjectGenerator
func TestStructureManagerCustomizationIntegration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "structure_customization_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a project generator
	pg := NewProjectGenerator()

	// Customize the structure
	customization := &StructureCustomization{
		AdditionalRootDirs: []string{"custom-docs", "custom-tools"},
		ExcludedDirs:       []string{"scripts"}, // Exclude the default scripts directory
	}

	err = pg.structureManager.CustomizeStructure(customization)
	if err != nil {
		t.Fatalf("CustomizeStructure failed: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "customization-test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	// Generate project structure with customization
	err = pg.GenerateProjectStructure(config, tempDir)
	if err != nil {
		t.Fatalf("GenerateProjectStructure with customization failed: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	// Verify custom directories were created
	customDocsPath := filepath.Join(projectPath, "custom-docs")
	if _, err := os.Stat(customDocsPath); os.IsNotExist(err) {
		t.Error("Custom directory 'custom-docs' was not created")
	}

	customToolsPath := filepath.Join(projectPath, "custom-tools")
	if _, err := os.Stat(customToolsPath); os.IsNotExist(err) {
		t.Error("Custom directory 'custom-tools' was not created")
	}

	// Verify excluded directory was not created
	scriptsPath := filepath.Join(projectPath, "scripts")
	if _, err := os.Stat(scriptsPath); !os.IsNotExist(err) {
		t.Error("Excluded directory 'scripts' was created when it should have been excluded")
	}

	t.Logf("✓ Customization integration test passed: Structure customization works correctly")
}

// TestBackwardCompatibility tests that the refactored code maintains backward compatibility
func TestBackwardCompatibility(t *testing.T) {
	// Test that GetStandardProjectStructure still works
	structure := GetStandardProjectStructure()
	if structure == nil {
		t.Fatal("GetStandardProjectStructure returned nil")
	}

	// Test that the structure has the expected directories
	if len(structure.RootDirs) == 0 {
		t.Error("Standard structure has no root directories")
	}

	if len(structure.FrontendDirs) == 0 {
		t.Error("Standard structure has no frontend directories")
	}

	if len(structure.BackendDirs) == 0 {
		t.Error("Standard structure has no backend directories")
	}

	// Test that NewProjectGenerator still works
	pg := NewProjectGenerator()
	if pg == nil {
		t.Fatal("NewProjectGenerator returned nil")
	}

	if pg.structureManager == nil {
		t.Error("ProjectGenerator does not have a structure manager")
	}

	// Test that NewDryRunProjectGenerator still works
	dryRunPg := NewDryRunProjectGenerator()
	if dryRunPg == nil {
		t.Fatal("NewDryRunProjectGenerator returned nil")
	}

	if dryRunPg.structureManager == nil {
		t.Error("DryRun ProjectGenerator does not have a structure manager")
	}

	t.Logf("✓ Backward compatibility test passed: All existing APIs work correctly")
}
