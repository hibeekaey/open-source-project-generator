package app

import (
	"os"
	"testing"
)

// TestApp_ConfigurationInitialization tests configuration manager initialization
func TestApp_ConfigurationInitialization(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	configManager, err := app.GetConfigManager()
	if err != nil {
		t.Fatalf("Failed to get config manager: %v", err)
	}
	if configManager == nil {
		t.Fatal("ConfigManager should not be nil")
	}
}

// TestApp_CacheConfiguration tests cache manager configuration
func TestApp_CacheConfiguration(t *testing.T) {
	// Test with normal home directory
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	cacheManager, err := app.GetCacheManager()
	if err != nil {
		t.Fatalf("Failed to get cache manager: %v", err)
	}
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Test cache manager initialization with fallback home directory
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", "")
	defer func() { _ = os.Setenv("HOME", originalHome) }()

	// This should work with fallback when HOME is empty
	app2, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Errorf("Expected app to work with fallback when HOME is empty, but got error: %v", err)
	}
	if app2 == nil {
		t.Error("Expected app to be created with fallback when HOME is empty")
	}
}

// TestApp_WorkspaceConfiguration tests workspace directory configuration
func TestApp_WorkspaceConfiguration(t *testing.T) {
	// Test with current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	securityManager, err := app.GetSecurityManager()
	if err != nil {
		t.Fatalf("Failed to get security manager: %v", err)
	}
	if securityManager == nil {
		t.Fatal("SecurityManager should not be nil")
	}

	// Test with temporary directory
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	app2, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app in temp directory: %v", err)
	}

	securityManager2, err := app2.GetSecurityManager()
	if err != nil {
		t.Fatalf("Failed to get security manager: %v", err)
	}
	if securityManager2 == nil {
		t.Fatal("SecurityManager should not be nil in temp directory")
	}
}

// TestApp_LoggerConfiguration tests logger configuration
func TestApp_LoggerConfiguration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	logger, err := app.GetLogger()
	if err != nil {
		t.Fatalf("Failed to get logger: %v", err)
	}
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test logger functionality
	logger.Info("Test configuration message")
	logger.Debug("Debug configuration message")
	logger.Warn("Warning configuration message")
	logger.Error("Error configuration message")
}

// TestApp_ComponentConfiguration tests individual component configuration
func TestApp_ComponentConfiguration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test validator configuration
	validator, err := app.GetValidator()
	if err != nil {
		t.Fatalf("Failed to get validator: %v", err)
	}
	if validator == nil {
		t.Fatal("Validator should not be nil")
	}

	// Test template engine configuration
	templateEngine, err := app.GetTemplateEngine()
	if err != nil {
		t.Fatalf("Failed to get template engine: %v", err)
	}
	if templateEngine == nil {
		t.Fatal("TemplateEngine should not be nil")
	}

	// Test template manager configuration
	templateManager, err := app.GetTemplateManager()
	if err != nil {
		t.Fatalf("Failed to get template manager: %v", err)
	}
	if templateManager == nil {
		t.Fatal("TemplateManager should not be nil")
	}

	// Test version manager configuration
	versionManager, err := app.GetVersionManager()
	if err != nil {
		t.Fatalf("Failed to get version manager: %v", err)
	}
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	// Test audit engine configuration
	auditEngine, err := app.GetAuditEngine()
	if err != nil {
		t.Fatalf("Failed to get audit engine: %v", err)
	}
	if auditEngine == nil {
		t.Fatal("AuditEngine should not be nil")
	}

	// Test filesystem generator configuration
	generator, err := app.GetGenerator()
	if err != nil {
		t.Fatalf("Failed to get generator: %v", err)
	}
	if generator == nil {
		t.Fatal("Generator should not be nil")
	}
}

// TestApp_CLIConfiguration tests CLI configuration with all dependencies
func TestApp_CLIConfiguration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	cli, err := app.GetCLI()
	if err != nil {
		t.Fatalf("Failed to get CLI: %v", err)
	}
	if cli == nil {
		t.Fatal("CLI should not be nil")
	}

	// Test that CLI can show help
	err = app.Run([]string{"--help"})
	if err != nil {
		t.Errorf("CLI help should work: %v", err)
	}
}

// TestApp_VersionManagerConfiguration tests version manager setup with cache
func TestApp_VersionManagerConfiguration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	versionManager, err := app.GetVersionManager()
	if err != nil {
		t.Fatalf("Failed to get version manager: %v", err)
	}
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	cacheManager, err := app.GetCacheManager()
	if err != nil {
		t.Fatalf("Failed to get cache manager: %v", err)
	}
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Version manager should be created with the app version and cache manager
	version, _, _ := app.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
}

// TestApp_ConfigurationConsistency tests that configuration is consistent across components
func TestApp_ConfigurationConsistency(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that version information is consistent
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

	// Test that all components are accessible
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

// TestApp_ConfigurationIsolation tests that multiple app instances have isolated configuration
func TestApp_ConfigurationIsolation(t *testing.T) {
	app1, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create first app: %v", err)
	}

	app2, err := NewApp("2.0.0", "def456", "2023-02-01")
	if err != nil {
		t.Fatalf("Failed to create second app: %v", err)
	}

	// Test that configurations are isolated
	version1, gitCommit1, buildTime1 := app1.GetVersion()
	version2, gitCommit2, buildTime2 := app2.GetVersion()

	if version1 == version2 {
		t.Error("App instances should have different versions")
	}
	if gitCommit1 == gitCommit2 {
		t.Error("App instances should have different git commits")
	}
	if buildTime1 == buildTime2 {
		t.Error("App instances should have different build times")
	}

	// Test that component instances are different
	logger1, err := app1.GetLogger()
	if err != nil {
		t.Fatalf("Failed to get logger from app1: %v", err)
	}
	logger2, err := app2.GetLogger()
	if err != nil {
		t.Fatalf("Failed to get logger from app2: %v", err)
	}
	// Check if loggers are different instances by comparing their addresses
	if logger1 == logger2 {
		t.Error("App instances should have different logger instances")
	}
	// Additional check: they should be different pointers
	if &logger1 == &logger2 {
		t.Error("Logger pointers should be different")
	}

	cache1, err := app1.GetCacheManager()
	if err != nil {
		t.Fatalf("Failed to get cache manager from app1: %v", err)
	}
	cache2, err := app2.GetCacheManager()
	if err != nil {
		t.Fatalf("Failed to get cache manager from app2: %v", err)
	}
	if cache1 == cache2 {
		t.Error("App instances should have different cache manager instances")
	}
}
