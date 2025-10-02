package app

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp("1.0.0", "abc123", "2023-01-01")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	if app == nil {
		t.Fatal("App should not be nil")
	}

	// Test that all components can be retrieved
	if _, err := app.GetConfigManager(); err != nil {
		t.Errorf("Failed to get config manager: %v", err)
	}

	if _, err := app.GetValidator(); err != nil {
		t.Errorf("Failed to get validator: %v", err)
	}

	if _, err := app.GetTemplateManager(); err != nil {
		t.Errorf("Failed to get template manager: %v", err)
	}

	if _, err := app.GetGenerator(); err != nil {
		t.Errorf("Failed to get generator: %v", err)
	}

	if _, err := app.GetTemplateEngine(); err != nil {
		t.Errorf("Failed to get template engine: %v", err)
	}

	if _, err := app.GetVersionManager(); err != nil {
		t.Errorf("Failed to get version manager: %v", err)
	}

	// Test version information
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
}
