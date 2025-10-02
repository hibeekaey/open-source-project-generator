package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestGetStandardProjectStructure(t *testing.T) {
	structure := GetStandardProjectStructure()

	if structure == nil {
		t.Fatal("GetStandardProjectStructure() returned nil")
	}

	// Verify root directories are defined
	if len(structure.RootDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty RootDirs")
	}

	// Verify component directories are defined
	if len(structure.FrontendDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty FrontendDirs")
	}

	if len(structure.BackendDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty BackendDirs")
	}

	if len(structure.MobileDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty MobileDirs")
	}

	if len(structure.InfraDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty InfraDirs")
	}

	// Verify expected directories are present
	expectedRootDirs := []string{"docs", "scripts", ".github/workflows"}
	for _, expectedDir := range expectedRootDirs {
		found := false
		for _, dir := range structure.RootDirs {
			if dir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetStandardProjectStructure() missing expected root directory: %s", expectedDir)
		}
	}

	// Verify frontend directories contain expected paths
	expectedFrontendDirs := []string{"App/src/components/ui", "Home/src/components", "Admin/src/components"}
	for _, expectedDir := range expectedFrontendDirs {
		found := false
		for _, dir := range structure.FrontendDirs {
			if dir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetStandardProjectStructure() missing expected frontend directory: %s", expectedDir)
		}
	}

	// Verify backend directories contain expected paths
	expectedBackendDirs := []string{"CommonServer/cmd", "CommonServer/internal/controllers", "CommonServer/pkg/auth"}
	for _, expectedDir := range expectedBackendDirs {
		found := false
		for _, dir := range structure.BackendDirs {
			if dir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetStandardProjectStructure() missing expected backend directory: %s", expectedDir)
		}
	}
}

func TestNewStructureManager(t *testing.T) {
	sm := NewStructureManager()

	if sm == nil {
		t.Fatal("NewStructureManager() returned nil")
	}

	if sm.structure == nil {
		t.Error("StructureManager structure is nil")
	}

	if sm.fsGen == nil {
		t.Error("StructureManager fsGen is nil")
	}
}

func TestNewDryRunStructureManager(t *testing.T) {
	sm := NewDryRunStructureManager()

	if sm == nil {
		t.Fatal("NewDryRunStructureManager() returned nil")
	}

	if sm.structure == nil {
		t.Error("DryRun StructureManager structure is nil")
	}

	if sm.fsGen == nil {
		t.Error("DryRun StructureManager fsGen is nil")
	}
}

func TestStructureManager_GetStructure(t *testing.T) {
	sm := NewStructureManager()
	structure := sm.GetStructure()

	if structure == nil {
		t.Error("GetStructure() returned nil")
		return
	}

	if len(structure.RootDirs) == 0 {
		t.Error("GetStructure() returned structure with empty RootDirs")
	}
}

func TestStructureManager_SetStructure(t *testing.T) {
	sm := NewStructureManager()

	// Test setting valid structure
	customStructure := &ProjectStructure{
		RootDirs:     []string{"custom-docs", "custom-scripts"},
		FrontendDirs: []string{"web/src"},
		BackendDirs:  []string{"api/src"},
		MobileDirs:   []string{"mobile/src"},
		InfraDirs:    []string{"infra/terraform"},
		CommonDirs:   []string{"tests"},
	}

	err := sm.SetStructure(customStructure)
	if err != nil {
		t.Errorf("SetStructure() failed with valid structure: %v", err)
	}

	// Verify structure was set
	if sm.GetStructure() != customStructure {
		t.Error("SetStructure() did not set the structure correctly")
	}

	// Test setting nil structure
	err = sm.SetStructure(nil)
	if err == nil {
		t.Error("SetStructure() should fail with nil structure")
	}

	// Test setting invalid structure (empty root dirs)
	invalidStructure := &ProjectStructure{
		RootDirs: []string{},
	}

	err = sm.SetStructure(invalidStructure)
	if err == nil {
		t.Error("SetStructure() should fail with invalid structure")
	}
}

func TestStructureManager_ValidateStructure(t *testing.T) {
	sm := NewStructureManager()

	// Test valid structure
	validStructure := &ProjectStructure{
		RootDirs:     []string{"docs", "scripts"},
		FrontendDirs: []string{"web/src"},
		BackendDirs:  []string{"api/src"},
		MobileDirs:   []string{"mobile/src"},
		InfraDirs:    []string{"infra/terraform"},
		CommonDirs:   []string{"tests"},
	}

	err := sm.ValidateStructure(validStructure)
	if err != nil {
		t.Errorf("ValidateStructure() failed with valid structure: %v", err)
	}

	// Test nil structure
	err = sm.ValidateStructure(nil)
	if err == nil {
		t.Error("ValidateStructure() should fail with nil structure")
	}

	// Test structure with empty root dirs
	invalidStructure := &ProjectStructure{
		RootDirs: []string{},
	}

	err = sm.ValidateStructure(invalidStructure)
	if err == nil {
		t.Error("ValidateStructure() should fail with empty root dirs")
	}

	// Test structure with empty directory path
	invalidStructure = &ProjectStructure{
		RootDirs: []string{"docs", ""},
	}

	err = sm.ValidateStructure(invalidStructure)
	if err == nil {
		t.Error("ValidateStructure() should fail with empty directory path")
	}

	// Test structure with absolute path
	invalidStructure = &ProjectStructure{
		RootDirs: []string{"/absolute/path"},
	}

	err = sm.ValidateStructure(invalidStructure)
	if err == nil {
		t.Error("ValidateStructure() should fail with absolute path")
	}

	// Test structure with path traversal
	invalidStructure = &ProjectStructure{
		RootDirs: []string{"../traversal"},
	}

	err = sm.ValidateStructure(invalidStructure)
	if err == nil {
		t.Error("ValidateStructure() should fail with path traversal")
	}
}

func TestStructureManager_CreateDirectories(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "structure_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewStructureManager()
	projectPath := filepath.Join(tempDir, "test-project")

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
	}

	// Test creating directories
	err = sm.CreateDirectories(projectPath, config)
	if err != nil {
		t.Errorf("CreateDirectories() failed: %v", err)
	}

	// Verify root directories were created
	for _, dir := range sm.structure.RootDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Root directory not created: %s", dir)
		}
	}

	// Verify frontend directories were created (since frontend is enabled)
	for _, dir := range sm.structure.FrontendDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Frontend directory not created: %s", dir)
		}
	}

	// Verify backend directories were created (since backend is enabled)
	for _, dir := range sm.structure.BackendDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Backend directory not created: %s", dir)
		}
	}

	// Test with empty project path
	err = sm.CreateDirectories("", config)
	if err == nil {
		t.Error("CreateDirectories() should fail with empty project path")
	}

	// Test with nil config
	err = sm.CreateDirectories(projectPath, nil)
	if err == nil {
		t.Error("CreateDirectories() should fail with nil config")
	}
}

func TestStructureManager_ValidateProjectStructure(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "structure_validate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewStructureManager()
	projectPath := filepath.Join(tempDir, "test-project")

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	// First create the directories
	err = sm.CreateDirectories(projectPath, config)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Test validation of existing structure
	err = sm.ValidateProjectStructure(projectPath, config)
	if err != nil {
		t.Errorf("ValidateProjectStructure() failed with valid structure: %v", err)
	}

	// Test with empty project path
	err = sm.ValidateProjectStructure("", config)
	if err == nil {
		t.Error("ValidateProjectStructure() should fail with empty project path")
	}

	// Test with nil config
	err = sm.ValidateProjectStructure(projectPath, nil)
	if err == nil {
		t.Error("ValidateProjectStructure() should fail with nil config")
	}

	// Test with non-existent project path
	err = sm.ValidateProjectStructure("/non/existent/path", config)
	if err == nil {
		t.Error("ValidateProjectStructure() should fail with non-existent path")
	}
}

func TestStructureManager_CustomizeStructure(t *testing.T) {
	sm := NewStructureManager()
	originalRootDirsCount := len(sm.structure.RootDirs)

	// Test adding additional directories
	customization := &StructureCustomization{
		AdditionalRootDirs:     []string{"custom-docs", "custom-scripts"},
		AdditionalFrontendDirs: []string{"web/custom"},
		AdditionalBackendDirs:  []string{"api/custom"},
		AdditionalMobileDirs:   []string{"mobile/custom"},
		AdditionalInfraDirs:    []string{"infra/custom"},
	}

	err := sm.CustomizeStructure(customization)
	if err != nil {
		t.Errorf("CustomizeStructure() failed: %v", err)
	}

	// Verify additional directories were added
	if len(sm.structure.RootDirs) != originalRootDirsCount+2 {
		t.Errorf("Expected %d root dirs, got %d", originalRootDirsCount+2, len(sm.structure.RootDirs))
	}

	// Verify specific directories were added
	found := false
	for _, dir := range sm.structure.RootDirs {
		if dir == "custom-docs" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Custom root directory 'custom-docs' was not added")
	}

	// Test excluding directories
	customization = &StructureCustomization{
		ExcludedDirs: []string{"custom-docs"},
	}

	err = sm.CustomizeStructure(customization)
	if err != nil {
		t.Errorf("CustomizeStructure() failed when excluding directories: %v", err)
	}

	// Verify excluded directory was removed
	for _, dir := range sm.structure.RootDirs {
		if dir == "custom-docs" {
			t.Error("Excluded directory 'custom-docs' was not removed")
		}
	}

	// Test with nil customization
	err = sm.CustomizeStructure(nil)
	if err == nil {
		t.Error("CustomizeStructure() should fail with nil customization")
	}
}

func TestStructureManager_filterDirectories(t *testing.T) {
	sm := NewStructureManager()

	dirs := []string{"dir1", "dir2", "dir3", "dir4"}
	excluded := []string{"dir2", "dir4"}

	filtered := sm.filterDirectories(dirs, excluded)

	expectedFiltered := []string{"dir1", "dir3"}
	if len(filtered) != len(expectedFiltered) {
		t.Errorf("Expected %d filtered directories, got %d", len(expectedFiltered), len(filtered))
	}

	for i, dir := range filtered {
		if dir != expectedFiltered[i] {
			t.Errorf("Expected filtered directory %s, got %s", expectedFiltered[i], dir)
		}
	}

	// Test with empty excluded list
	filtered = sm.filterDirectories(dirs, []string{})
	if len(filtered) != len(dirs) {
		t.Errorf("Expected %d directories when no exclusions, got %d", len(dirs), len(filtered))
	}

	// Test with nil excluded list
	filtered = sm.filterDirectories(dirs, nil)
	if len(filtered) != len(dirs) {
		t.Errorf("Expected %d directories when nil exclusions, got %d", len(dirs), len(filtered))
	}
}

func TestStructureManager_createComponentDirectories(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "component_dirs_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewStructureManager()
	projectPath := filepath.Join(tempDir, "test-project")

	// Test with frontend enabled
	config := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	err = sm.createComponentDirectories(projectPath, config)
	if err != nil {
		t.Errorf("createComponentDirectories() failed with frontend config: %v", err)
	}

	// Verify frontend directories were created
	for _, dir := range sm.structure.FrontendDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Frontend directory not created: %s", dir)
		}
	}

	// Test with backend enabled
	config = &models.ProjectConfig{
		Components: models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
	}

	err = sm.createComponentDirectories(projectPath, config)
	if err != nil {
		t.Errorf("createComponentDirectories() failed with backend config: %v", err)
	}

	// Test with mobile enabled
	config = &models.ProjectConfig{
		Components: models.Components{
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
		},
	}

	err = sm.createComponentDirectories(projectPath, config)
	if err != nil {
		t.Errorf("createComponentDirectories() failed with mobile config: %v", err)
	}

	// Test with infrastructure enabled
	config = &models.ProjectConfig{
		Components: models.Components{
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
	}

	err = sm.createComponentDirectories(projectPath, config)
	if err != nil {
		t.Errorf("createComponentDirectories() failed with infrastructure config: %v", err)
	}
}

func TestStructureManager_validateComponentDirectories(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "validate_component_dirs_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewStructureManager()
	projectPath := filepath.Join(tempDir, "test-project")

	// Create test configuration with frontend enabled
	config := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	// First create the directories
	err = sm.createComponentDirectories(projectPath, config)
	if err != nil {
		t.Fatalf("Failed to create component directories: %v", err)
	}

	// Test validation of existing directories
	err = sm.validateComponentDirectories(projectPath, config)
	if err != nil {
		t.Errorf("validateComponentDirectories() failed with valid directories: %v", err)
	}

	// Test validation with missing directories (remove one directory)
	firstFrontendDir := sm.structure.FrontendDirs[0]
	dirToRemove := filepath.Join(projectPath, firstFrontendDir)
	os.RemoveAll(dirToRemove)

	err = sm.validateComponentDirectories(projectPath, config)
	if err == nil {
		t.Error("validateComponentDirectories() should fail with missing directory")
	}
}

func TestStructureCustomization(t *testing.T) {
	customization := &StructureCustomization{
		AdditionalRootDirs:     []string{"custom1", "custom2"},
		AdditionalFrontendDirs: []string{"web/custom"},
		AdditionalBackendDirs:  []string{"api/custom"},
		AdditionalMobileDirs:   []string{"mobile/custom"},
		AdditionalInfraDirs:    []string{"infra/custom"},
		ExcludedDirs:           []string{"docs"},
	}

	if len(customization.AdditionalRootDirs) != 2 {
		t.Error("StructureCustomization AdditionalRootDirs not set correctly")
	}

	if len(customization.ExcludedDirs) != 1 {
		t.Error("StructureCustomization ExcludedDirs not set correctly")
	}
}
