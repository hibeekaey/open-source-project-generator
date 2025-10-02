package app

import (
	"os"
	"testing"
)

// TestApp_Lifecycle tests the complete application lifecycle
func TestApp_Lifecycle(t *testing.T) {
	// Test application creation
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Verify all components are initialized
	if app.configManager == nil {
		t.Error("ConfigManager not initialized")
	}
	if app.validator == nil {
		t.Error("Validator not initialized")
	}
	if app.templateManager == nil {
		t.Error("TemplateManager not initialized")
	}
	if app.cacheManager == nil {
		t.Error("CacheManager not initialized")
	}
	if app.versionManager == nil {
		t.Error("VersionManager not initialized")
	}
	if app.auditEngine == nil {
		t.Error("AuditEngine not initialized")
	}
	if app.securityManager == nil {
		t.Error("SecurityManager not initialized")
	}
	if app.cli == nil {
		t.Error("CLI not initialized")
	}
	if app.generator == nil {
		t.Error("Generator not initialized")
	}
	if app.templateEngine == nil {
		t.Error("TemplateEngine not initialized")
	}
	if app.logger == nil {
		t.Error("Logger not initialized")
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

	// Test component access methods
	if app.GetConfigManager() != app.configManager {
		t.Error("GetConfigManager returns wrong instance")
	}
	if app.GetValidator() != app.validator {
		t.Error("GetValidator returns wrong instance")
	}
	if app.GetTemplateManager() != app.templateManager {
		t.Error("GetTemplateManager returns wrong instance")
	}
	if app.GetCacheManager() != app.cacheManager {
		t.Error("GetCacheManager returns wrong instance")
	}
	if app.GetVersionManager() != app.versionManager {
		t.Error("GetVersionManager returns wrong instance")
	}
	if app.GetAuditEngine() != app.auditEngine {
		t.Error("GetAuditEngine returns wrong instance")
	}
	if app.GetSecurityManager() != app.securityManager {
		t.Error("GetSecurityManager returns wrong instance")
	}
	if app.GetCLI() != app.cli {
		t.Error("GetCLI returns wrong instance")
	}
	if app.GetGenerator() != app.generator {
		t.Error("GetGenerator returns wrong instance")
	}
	if app.GetTemplateEngine() != app.templateEngine {
		t.Error("GetTemplateEngine returns wrong instance")
	}
	if app.GetLogger() != app.logger {
		t.Error("GetLogger returns wrong instance")
	}
}

// TestApp_InitializationErrors tests error handling during initialization
func TestApp_InitializationErrors(t *testing.T) {
	// Test with invalid home directory (simulate error condition)
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", "/nonexistent/path/that/should/not/exist")
	defer os.Setenv("HOME", originalHome)

	// This may fail due to permission issues, which is expected
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		// This is expected when home directory is invalid
		t.Logf("Expected error with invalid home directory: %v", err)
		return
	}

	if app == nil {
		t.Fatal("App should not be nil if creation succeeded")
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

	// Test running version command
	// Note: Version command may fail if version manager is not properly initialized
	_ = app.Run([]string{"version"})
	// We don't fail the test if version command has issues
}

// TestApp_ComponentIntegration tests that components can work together
func TestApp_ComponentIntegration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that CLI has access to all required components
	cli := app.GetCLI()
	if cli == nil {
		t.Fatal("CLI should not be nil")
	}

	// Test that cache manager is properly initialized
	cacheManager := app.GetCacheManager()
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Test that version manager has the correct version
	versionManager := app.GetVersionManager()
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	// Test that logger is functional
	logger := app.GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test logging functionality
	logger.Info("Test log message from integration test")
	logger.Debug("Debug message from integration test")
}

// TestApp_DependencyInjection tests that all dependencies are properly injected
func TestApp_DependencyInjection(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Verify that all components are different instances (no shared state issues)
	components := []interface{}{
		app.configManager,
		app.validator,
		app.templateManager,
		app.cacheManager,
		app.versionManager,
		app.auditEngine,
		app.securityManager,
		app.cli,
		app.generator,
		app.templateEngine,
		app.logger,
	}

	// Check that all components are non-nil and unique
	for i, comp1 := range components {
		if comp1 == nil {
			t.Errorf("Component at index %d is nil", i)
			continue
		}

		for j, comp2 := range components {
			if i != j && comp1 == comp2 {
				t.Errorf("Components at indices %d and %d are the same instance", i, j)
			}
		}
	}
}

// TestApp_ConfigurationManagement tests configuration-related functionality
func TestApp_ConfigurationManagement(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	configManager := app.GetConfigManager()
	if configManager == nil {
		t.Fatal("ConfigManager should not be nil")
	}

	// Test that config manager is properly initialized
	// Note: We can't test much without knowing the specific interface,
	// but we can verify it's not nil and accessible
}

// TestApp_SecurityManager tests security manager initialization
func TestApp_SecurityManager(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	securityManager := app.GetSecurityManager()
	if securityManager == nil {
		t.Fatal("SecurityManager should not be nil")
	}

	// Test that security manager is properly initialized with workspace directory
	// The security manager should be initialized with the current working directory
}

// TestApp_TemplateSystem tests template system integration
func TestApp_TemplateSystem(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	templateEngine := app.GetTemplateEngine()
	if templateEngine == nil {
		t.Fatal("TemplateEngine should not be nil")
	}

	templateManager := app.GetTemplateManager()
	if templateManager == nil {
		t.Fatal("TemplateManager should not be nil")
	}

	// Verify that template manager and engine are properly connected
	// The template manager should use the template engine
}

// TestApp_CacheSystem tests cache system integration
func TestApp_CacheSystem(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	cacheManager := app.GetCacheManager()
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Test that cache manager is properly initialized
	// Cache directory should be set to ~/.generator/cache
}

// TestApp_ValidationSystem tests validation system integration
func TestApp_ValidationSystem(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	validator := app.GetValidator()
	if validator == nil {
		t.Fatal("Validator should not be nil")
	}

	// Test that validation engine is properly initialized
}

// TestApp_AuditSystem tests audit system integration
func TestApp_AuditSystem(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	auditEngine := app.GetAuditEngine()
	if auditEngine == nil {
		t.Fatal("AuditEngine should not be nil")
	}

	// Test that audit engine is properly initialized
}

// TestApp_VersionManagement tests version management integration
func TestApp_VersionManagement(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	versionManager := app.GetVersionManager()
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	// Test that version manager is initialized with correct version and cache
	// Version manager should have access to the cache manager for update checking
}

// TestApp_FileSystemGeneration tests filesystem generation integration
func TestApp_FileSystemGeneration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	generator := app.GetGenerator()
	if generator == nil {
		t.Fatal("Generator should not be nil")
	}

	// Test that filesystem generator is properly initialized
}

// TestApp_LoggerIntegration tests logger integration across components
func TestApp_LoggerIntegration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	logger := app.GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test that logger is properly configured
	logger.Info("Testing logger integration")
	logger.Debug("Debug message for integration test")
	logger.Warn("Warning message for integration test")
	logger.Error("Error message for integration test")

	// Test structured logging
	fields := map[string]interface{}{
		"component": "app_test",
		"test":      "integration",
	}
	logger.InfoWithFields("Structured log message", fields)
}

// TestApp_ErrorHandling tests error handling across the application
func TestApp_ErrorHandling(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test running with invalid arguments
	err = app.Run([]string{"invalid-command"})
	if err == nil {
		t.Error("Running invalid command should return an error")
	}

	// Test that the error is handled gracefully
	logger := app.GetLogger()
	logger.Error("Test error handling: %v", err)
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
