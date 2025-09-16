package app

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp("test-version", "test-commit", "test-time")
	if err != nil {
		t.Fatalf("NewApp() returned error: %v", err)
	}

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	if app.configManager == nil {
		t.Error("App config manager not initialized")
	}

	if app.validator == nil {
		t.Error("App validator not initialized")
	}

	if app.cli == nil {
		t.Error("App CLI not initialized")
	}

	if app.generator == nil {
		t.Error("App generator not initialized")
	}

	if app.templateEngine == nil {
		t.Error("App template engine not initialized")
	}

	if app.versionManager == nil {
		t.Error("App version manager not initialized")
	}
}
