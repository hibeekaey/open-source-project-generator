package container

import (
	"testing"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("NewContainer should not return nil")
	}

	// All components should be nil initially since no services are registered
	if cli, _ := container.GetCLI(); cli != nil {
		t.Error("CLI should be nil initially")
	}
	if engine, _ := container.GetTemplateEngine(); engine != nil {
		t.Error("TemplateEngine should be nil initially")
	}
	if manager, _ := container.GetConfigManager(); manager != nil {
		t.Error("ConfigManager should be nil initially")
	}
	if generator, _ := container.GetFileSystemGenerator(); generator != nil {
		t.Error("FileSystemGenerator should be nil initially")
	}
	if versionManager, _ := container.GetVersionManager(); versionManager != nil {
		t.Error("VersionManager should be nil initially")
	}
	if validator, _ := container.GetValidator(); validator != nil {
		t.Error("Validator should be nil initially")
	}
}

// TODO: Add proper container tests using RegisterService pattern
// The previous tests were using a different API that no longer exists
func TestContainerServiceRegistration(t *testing.T) {
	container := NewContainer()

	// Test that we can register a service factory
	err := container.RegisterService("test", func() (interface{}, error) {
		return "test-service", nil
	})

	if err != nil {
		t.Errorf("Failed to register service: %v", err)
	}

	// Test that we can retrieve the service
	service, err := container.GetService("test")
	if err != nil {
		t.Errorf("Failed to get service: %v", err)
	}

	if service != "test-service" {
		t.Errorf("Expected 'test-service', got %v", service)
	}
}
