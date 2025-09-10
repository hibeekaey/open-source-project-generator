package container

import (
	"testing"

	"github.com/open-source-template-generator/pkg/cli"
	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
	"github.com/open-source-template-generator/pkg/validation"
	"github.com/open-source-template-generator/pkg/version"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("NewContainer() returned nil")
	}

	// All components should be nil initially
	if container.GetCLI() != nil {
		t.Error("CLI should be nil initially")
	}
	if container.GetTemplateEngine() != nil {
		t.Error("TemplateEngine should be nil initially")
	}
	if container.GetConfigManager() != nil {
		t.Error("ConfigManager should be nil initially")
	}
	if container.GetFileSystemGenerator() != nil {
		t.Error("FileSystemGenerator should be nil initially")
	}
	if container.GetVersionManager() != nil {
		t.Error("VersionManager should be nil initially")
	}
	if container.GetValidator() != nil {
		t.Error("Validator should be nil initially")
	}
}

func TestContainerCLI(t *testing.T) {
	container := NewContainer()

	// Create a mock CLI (we'll use the real CLI for testing)
	mockCLI := cli.NewCLI(nil, nil)

	// Test SetCLI and GetCLI
	container.SetCLI(mockCLI)
	retrievedCLI := container.GetCLI()

	if retrievedCLI != mockCLI {
		t.Error("SetCLI/GetCLI not working correctly")
	}
}

func TestContainerTemplateEngine(t *testing.T) {
	container := NewContainer()

	// Create a template engine
	templateEngine := template.NewEngine()

	// Test SetTemplateEngine and GetTemplateEngine
	container.SetTemplateEngine(templateEngine)
	retrievedEngine := container.GetTemplateEngine()

	if retrievedEngine != templateEngine {
		t.Error("SetTemplateEngine/GetTemplateEngine not working correctly")
	}
}

func TestContainerConfigManager(t *testing.T) {
	container := NewContainer()

	// We can't easily create a config manager without dependencies,
	// so we'll test with nil and ensure the container handles it
	container.SetConfigManager(nil)
	retrievedManager := container.GetConfigManager()

	if retrievedManager != nil {
		t.Error("SetConfigManager/GetConfigManager not working correctly with nil")
	}
}

func TestContainerFileSystemGenerator(t *testing.T) {
	container := NewContainer()

	// Create a filesystem generator
	fsGenerator := filesystem.NewGenerator()

	// Test SetFileSystemGenerator and GetFileSystemGenerator
	container.SetFileSystemGenerator(fsGenerator)
	retrievedGenerator := container.GetFileSystemGenerator()

	if retrievedGenerator != fsGenerator {
		t.Error("SetFileSystemGenerator/GetFileSystemGenerator not working correctly")
	}
}

func TestContainerVersionManager(t *testing.T) {
	container := NewContainer()

	// Create a version manager with memory cache
	cache := version.NewMemoryCache(0) // No TTL for testing
	versionManager := version.NewManager(cache)

	// Test SetVersionManager and GetVersionManager
	container.SetVersionManager(versionManager)
	retrievedManager := container.GetVersionManager()

	if retrievedManager != versionManager {
		t.Error("SetVersionManager/GetVersionManager not working correctly")
	}
}

func TestContainerValidator(t *testing.T) {
	container := NewContainer()

	// Create a validation engine
	validator := validation.NewEngine()

	// Test SetValidator and GetValidator
	container.SetValidator(validator)
	retrievedValidator := container.GetValidator()

	if retrievedValidator != validator {
		t.Error("SetValidator/GetValidator not working correctly")
	}
}

func TestContainerMultipleComponents(t *testing.T) {
	container := NewContainer()

	// Set multiple components
	templateEngine := template.NewEngine()
	fsGenerator := filesystem.NewGenerator()
	validator := validation.NewEngine()
	cache := version.NewMemoryCache(0)
	versionManager := version.NewManager(cache)

	container.SetTemplateEngine(templateEngine)
	container.SetFileSystemGenerator(fsGenerator)
	container.SetValidator(validator)
	container.SetVersionManager(versionManager)

	// Verify all components are set correctly
	if container.GetTemplateEngine() != templateEngine {
		t.Error("TemplateEngine not set correctly")
	}
	if container.GetFileSystemGenerator() != fsGenerator {
		t.Error("FileSystemGenerator not set correctly")
	}
	if container.GetValidator() != validator {
		t.Error("Validator not set correctly")
	}
	if container.GetVersionManager() != versionManager {
		t.Error("VersionManager not set correctly")
	}
}

func TestContainerOverwrite(t *testing.T) {
	container := NewContainer()

	// Set initial components
	templateEngine1 := template.NewEngine()
	templateEngine2 := template.NewEngine()

	container.SetTemplateEngine(templateEngine1)
	if container.GetTemplateEngine() != templateEngine1 {
		t.Error("Initial TemplateEngine not set correctly")
	}

	// Overwrite with new component
	container.SetTemplateEngine(templateEngine2)
	if container.GetTemplateEngine() != templateEngine2 {
		t.Error("TemplateEngine not overwritten correctly")
	}

	// Ensure it's not the old one
	if container.GetTemplateEngine() == templateEngine1 {
		t.Error("TemplateEngine still references old instance")
	}
}

func TestContainerNilHandling(t *testing.T) {
	container := NewContainer()

	// Test setting nil values
	container.SetCLI(nil)
	container.SetTemplateEngine(nil)
	container.SetConfigManager(nil)
	container.SetFileSystemGenerator(nil)
	container.SetVersionManager(nil)
	container.SetValidator(nil)

	// All should return nil
	if container.GetCLI() != nil {
		t.Error("CLI should be nil after setting to nil")
	}
	if container.GetTemplateEngine() != nil {
		t.Error("TemplateEngine should be nil after setting to nil")
	}
	if container.GetConfigManager() != nil {
		t.Error("ConfigManager should be nil after setting to nil")
	}
	if container.GetFileSystemGenerator() != nil {
		t.Error("FileSystemGenerator should be nil after setting to nil")
	}
	if container.GetVersionManager() != nil {
		t.Error("VersionManager should be nil after setting to nil")
	}
	if container.GetValidator() != nil {
		t.Error("Validator should be nil after setting to nil")
	}
}

// Mock implementations for testing
type mockCLI struct{}

func (m *mockCLI) Run() error                                           { return nil }
func (m *mockCLI) PromptProjectDetails() (*models.ProjectConfig, error) { return nil, nil }
func (m *mockCLI) SelectComponents() ([]string, error)                  { return nil, nil }
func (m *mockCLI) ConfirmGeneration(*models.ProjectConfig) bool         { return true }
func (m *mockCLI) ShowProgress(string)                                  {}
func (m *mockCLI) ShowError(string)                                     {}
func (m *mockCLI) ShowSuccess(string)                                   {}
func (m *mockCLI) ShowWarning(string)                                   {}
func (m *mockCLI) PreviewConfiguration(*models.ProjectConfig)           {}
func (m *mockCLI) CheckOutputPath(string) error                         { return nil }

func TestContainerWithMockCLI(t *testing.T) {
	container := NewContainer()

	// Test with mock CLI
	mockCLI := &mockCLI{}
	container.SetCLI(mockCLI)

	retrievedCLI := container.GetCLI()
	if retrievedCLI == nil {
		t.Error("Mock CLI not set correctly")
	}

	// Test that we can call methods on the retrieved CLI
	err := retrievedCLI.Run()
	if err != nil {
		t.Errorf("Mock CLI Run() returned error: %v", err)
	}
}
