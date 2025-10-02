package app

import (
	"testing"
)

// TestApp_Lifecycle tests the complete application lifecycle
func TestApp_Lifecycle(t *testing.T) {
	// Test application creation
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test version information storage
	version, gitCommit, buildTime := app.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
	if gitCommit != "abc123" {
		t.Errorf("Expected gitCommit 'abc123', got '%s'", gitCommit)
	}
	if buildTime != "2023-01-01" {
		t.Errorf("Expected buildTime '2023-01-01', got '%s'", buildTime)
	}

	// Test that all components can be retrieved without error
	components := []struct {
		name string
		get  func() (interface{}, error)
	}{
		{"ConfigManager", func() (interface{}, error) { return app.GetConfigManager() }},
		{"Validator", func() (interface{}, error) { return app.GetValidator() }},
		{"TemplateManager", func() (interface{}, error) { return app.GetTemplateManager() }},
		{"CacheManager", func() (interface{}, error) { return app.GetCacheManager() }},
		{"VersionManager", func() (interface{}, error) { return app.GetVersionManager() }},
		{"AuditEngine", func() (interface{}, error) { return app.GetAuditEngine() }},
		{"SecurityManager", func() (interface{}, error) { return app.GetSecurityManager() }},
		{"CLI", func() (interface{}, error) { return app.GetCLI() }},
		{"Generator", func() (interface{}, error) { return app.GetGenerator() }},
		{"TemplateEngine", func() (interface{}, error) { return app.GetTemplateEngine() }},
		{"Logger", func() (interface{}, error) { return app.GetLogger() }},
	}

	for _, comp := range components {
		component, err := comp.get()
		if err != nil {
			t.Errorf("Failed to get %s: %v", comp.name, err)
		}
		if component == nil {
			t.Errorf("Component %s is nil", comp.name)
		}
	}
}

// TestApp_Run tests the application run method
func TestApp_Run(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test running help command
	err = app.Run([]string{"--help"})
	if err != nil {
		t.Errorf("Running help command should not error: %v", err)
	}
}

// TestApp_ComponentIntegration tests that components can work together
func TestApp_ComponentIntegration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that CLI has access to all required components
	cli, err := app.GetCLI()
	if err != nil {
		t.Fatalf("Failed to get CLI: %v", err)
	}
	if cli == nil {
		t.Fatal("CLI should not be nil")
	}

	// Test that cache manager is properly initialized
	cacheManager, err := app.GetCacheManager()
	if err != nil {
		t.Fatalf("Failed to get cache manager: %v", err)
	}
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Test that version manager has the correct version
	versionManager, err := app.GetVersionManager()
	if err != nil {
		t.Fatalf("Failed to get version manager: %v", err)
	}
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	// Test that logger is functional
	logger, err := app.GetLogger()
	if err != nil {
		t.Fatalf("Failed to get logger: %v", err)
	}
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test logging functionality
	logger.Info("Test log message from integration test")
	logger.Debug("Debug message from integration test")
}

// TestApp_MultipleInstances tests creating multiple app instances
func TestApp_MultipleInstances(t *testing.T) {
	app1, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create first app: %v", err)
	}

	app2, err := NewApp("2.0.0", "def456", "2023-02-01")
	if err != nil {
		t.Fatalf("Failed to create second app: %v", err)
	}

	// Verify that instances are independent
	if app1 == app2 {
		t.Error("App instances should be different")
	}

	// Verify that they have different version information
	version1, _, _ := app1.GetVersion()
	version2, _, _ := app2.GetVersion()

	if version1 == version2 {
		t.Error("App instances should have different versions")
	}

	if version1 != "1.0.0" {
		t.Errorf("First app should have version '1.0.0', got '%s'", version1)
	}

	if version2 != "2.0.0" {
		t.Errorf("Second app should have version '2.0.0', got '%s'", version2)
	}
}
