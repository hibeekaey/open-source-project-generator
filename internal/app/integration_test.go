package app

import (
	"testing"
)

// TestAppIntegration tests that all components are properly integrated in the App struct
func TestAppIntegration(t *testing.T) {
	// Create a new app instance
	app, err := NewApp("test-version", "test-commit", "test-time")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that all managers are properly initialized
	if app.GetConfigManager() == nil {
		t.Error("ConfigManager is nil")
	}

	if app.GetValidator() == nil {
		t.Error("Validator is nil")
	}

	if app.GetTemplateManager() == nil {
		t.Error("TemplateManager is nil")
	}

	if app.GetCacheManager() == nil {
		t.Error("CacheManager is nil")
	}

	if app.GetVersionManager() == nil {
		t.Error("VersionManager is nil")
	}

	if app.GetAuditEngine() == nil {
		t.Error("AuditEngine is nil")
	}

	if app.GetSecurityManager() == nil {
		t.Error("SecurityManager is nil")
	}

	if app.GetCLI() == nil {
		t.Error("CLI is nil")
	}

	if app.GetGenerator() == nil {
		t.Error("Generator is nil")
	}

	if app.GetTemplateEngine() == nil {
		t.Error("TemplateEngine is nil")
	}

	if app.GetLogger() == nil {
		t.Error("Logger is nil")
	}

	// Test version information
	version, gitCommit, buildTime := app.GetVersion()
	if version != "test-version" {
		t.Errorf("Expected version 'test-version', got '%s'", version)
	}
	if gitCommit != "test-commit" {
		t.Errorf("Expected gitCommit 'test-commit', got '%s'", gitCommit)
	}
	if buildTime != "test-time" {
		t.Errorf("Expected buildTime 'test-time', got '%s'", buildTime)
	}
}

// TestAppComponentsIntegration tests that components can interact with each other
func TestAppComponentsIntegration(t *testing.T) {
	app, err := NewApp("test-version", "test-commit", "test-time")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that CLI can access all managers (this would fail if CLI constructor doesn't match)
	cli := app.GetCLI()
	if cli == nil {
		t.Fatal("CLI is nil")
	}

	// Test that we can run help command without errors
	err = app.Run([]string{"--help"})
	if err != nil {
		t.Errorf("Failed to run help command: %v", err)
	}
}

// TestAppManagersNotNil ensures all managers are properly initialized and not nil
func TestAppManagersNotNil(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	managers := map[string]interface{}{
		"ConfigManager":   app.configManager,
		"Validator":       app.validator,
		"TemplateManager": app.templateManager,
		"CacheManager":    app.cacheManager,
		"VersionManager":  app.versionManager,
		"AuditEngine":     app.auditEngine,
		"SecurityManager": app.securityManager,
		"CLI":             app.cli,
		"Generator":       app.generator,
		"TemplateEngine":  app.templateEngine,
		"Logger":          app.logger,
	}

	for name, manager := range managers {
		if manager == nil {
			t.Errorf("%s is nil", name)
		}
	}
}
