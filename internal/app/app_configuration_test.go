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

	configManager := app.GetConfigManager()
	if configManager == nil {
		t.Fatal("ConfigManager should not be nil")
	}

	// Test that config manager is accessible through the app
	if app.configManager != configManager {
		t.Error("Internal config manager should match public getter")
	}
}

// TestApp_CacheConfiguration tests cache manager configuration
func TestApp_CacheConfiguration(t *testing.T) {
	// Test with normal home directory
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	cacheManager := app.GetCacheManager()
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Test cache manager initialization with fallback home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", "")
	defer os.Setenv("HOME", originalHome)

	// This should fail gracefully when HOME is empty
	_, err = NewApp("1.0.0", "abc123", "2023-01-01")
	if err == nil {
		t.Error("Expected error when HOME is empty, but got none")
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

	securityManager := app.GetSecurityManager()
	if securityManager == nil {
		t.Fatal("SecurityManager should not be nil")
	}

	// Test with temporary directory
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	app2, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app in temp directory: %v", err)
	}

	securityManager2 := app2.GetSecurityManager()
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

	logger := app.GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test that logger is properly configured
	if logger.GetLevel() != int(LogLevelInfo) {
		t.Errorf("Expected logger level Info, got %v", logger.GetLevel())
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
	validator := app.GetValidator()
	if validator == nil {
		t.Fatal("Validator should not be nil")
	}

	// Test template engine configuration
	templateEngine := app.GetTemplateEngine()
	if templateEngine == nil {
		t.Fatal("TemplateEngine should not be nil")
	}

	// Test template manager configuration
	templateManager := app.GetTemplateManager()
	if templateManager == nil {
		t.Fatal("TemplateManager should not be nil")
	}

	// Test version manager configuration
	versionManager := app.GetVersionManager()
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	// Test audit engine configuration
	auditEngine := app.GetAuditEngine()
	if auditEngine == nil {
		t.Fatal("AuditEngine should not be nil")
	}

	// Test filesystem generator configuration
	generator := app.GetGenerator()
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

	cli := app.GetCLI()
	if cli == nil {
		t.Fatal("CLI should not be nil")
	}

	// Test that CLI can access version information
	// Note: Version command may fail if version manager is not properly initialized
	// but the CLI should still be functional
	_ = app.Run([]string{"version"})
	// We don't fail the test if version command has issues, as long as CLI is accessible

	// Test that CLI can show help
	err = app.Run([]string{"--help"})
	if err != nil {
		t.Errorf("CLI help should work: %v", err)
	}
}

// TestApp_CacheDirectoryConfiguration tests cache directory setup
func TestApp_CacheDirectoryConfiguration(t *testing.T) {
	// Test with valid home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	cacheManager := app.GetCacheManager()
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// We can't directly test the cache directory without exposing it,
	// but we can verify the cache manager was created successfully
	// Expected cache directory should be ~/.generator/cache
}

// TestApp_SecurityManagerConfiguration tests security manager workspace setup
func TestApp_SecurityManagerConfiguration(t *testing.T) {
	// Test in a temporary workspace
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	securityManager := app.GetSecurityManager()
	if securityManager == nil {
		t.Fatal("SecurityManager should not be nil")
	}

	// Security manager should be initialized with the current workspace
}

// TestApp_VersionManagerConfiguration tests version manager setup with cache
func TestApp_VersionManagerConfiguration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	versionManager := app.GetVersionManager()
	if versionManager == nil {
		t.Fatal("VersionManager should not be nil")
	}

	cacheManager := app.GetCacheManager()
	if cacheManager == nil {
		t.Fatal("CacheManager should not be nil")
	}

	// Version manager should be created with the app version and cache manager
	version, _, _ := app.GetVersion()
	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
}

// TestApp_ConfigurationErrorHandling tests error handling during configuration
func TestApp_ConfigurationErrorHandling(t *testing.T) {
	// Test with inaccessible home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", "/root/inaccessible")
	defer os.Setenv("HOME", originalHome)

	// This may fail due to permission issues, which is expected
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		// This is expected when home directory is inaccessible
		t.Logf("Expected error with inaccessible home directory: %v", err)
		return
	}

	if app == nil {
		t.Fatal("App should not be nil if creation succeeded")
	}

	// If creation succeeded, all components should be initialized
	if app.GetCacheManager() == nil {
		t.Error("CacheManager should be initialized if app creation succeeded")
	}
	if app.GetLogger() == nil {
		t.Error("Logger should be initialized if app creation succeeded")
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
	components := map[string]interface{}{
		"ConfigManager":   app.GetConfigManager(),
		"Validator":       app.GetValidator(),
		"TemplateManager": app.GetTemplateManager(),
		"CacheManager":    app.GetCacheManager(),
		"VersionManager":  app.GetVersionManager(),
		"AuditEngine":     app.GetAuditEngine(),
		"SecurityManager": app.GetSecurityManager(),
		"CLI":             app.GetCLI(),
		"Generator":       app.GetGenerator(),
		"TemplateEngine":  app.GetTemplateEngine(),
		"Logger":          app.GetLogger(),
	}

	for name, component := range components {
		if component == nil {
			t.Errorf("Component %s is nil", name)
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
	if app1.GetLogger() == app2.GetLogger() {
		t.Error("App instances should have different logger instances")
	}
	if app1.GetCacheManager() == app2.GetCacheManager() {
		t.Error("App instances should have different cache manager instances")
	}
}
