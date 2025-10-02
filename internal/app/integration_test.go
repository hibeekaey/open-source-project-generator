package app

import (
	"testing"
)

func TestApp_Integration(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that all components are accessible
	if configManager, err := app.GetConfigManager(); err != nil || configManager == nil {
		t.Errorf("ConfigManager should be accessible: %v", err)
	}

	if validator, err := app.GetValidator(); err != nil || validator == nil {
		t.Errorf("Validator should be accessible: %v", err)
	}

	if templateManager, err := app.GetTemplateManager(); err != nil || templateManager == nil {
		t.Errorf("TemplateManager should be accessible: %v", err)
	}

	if cacheManager, err := app.GetCacheManager(); err != nil || cacheManager == nil {
		t.Errorf("CacheManager should be accessible: %v", err)
	}

	if versionManager, err := app.GetVersionManager(); err != nil || versionManager == nil {
		t.Errorf("VersionManager should be accessible: %v", err)
	}

	if auditEngine, err := app.GetAuditEngine(); err != nil || auditEngine == nil {
		t.Errorf("AuditEngine should be accessible: %v", err)
	}

	if securityManager, err := app.GetSecurityManager(); err != nil || securityManager == nil {
		t.Errorf("SecurityManager should be accessible: %v", err)
	}

	if cli, err := app.GetCLI(); err != nil || cli == nil {
		t.Errorf("CLI should be accessible: %v", err)
	}

	if generator, err := app.GetGenerator(); err != nil || generator == nil {
		t.Errorf("Generator should be accessible: %v", err)
	}

	if templateEngine, err := app.GetTemplateEngine(); err != nil || templateEngine == nil {
		t.Errorf("TemplateEngine should be accessible: %v", err)
	}

	if logger, err := app.GetLogger(); err != nil || logger == nil {
		t.Errorf("Logger should be accessible: %v", err)
	}
}
