package fallback

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if registry.generators == nil {
		t.Fatal("Registry generators map is nil")
	}
}

func TestDefaultRegistry(t *testing.T) {
	registry := DefaultRegistry()
	if registry == nil {
		t.Fatal("DefaultRegistry returned nil")
	}

	// Check that Android generator is registered
	if !registry.Supports("android") {
		t.Error("Android generator not registered")
	}

	// Check that iOS generator is registered
	if !registry.Supports("ios") {
		t.Error("iOS generator not registered")
	}

	// Check supported types
	types := registry.GetSupportedTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 supported types, got %d", len(types))
	}
}

func TestRegistryRegisterAndGet(t *testing.T) {
	registry := NewRegistry()
	androidGen := NewAndroidGenerator()

	registry.Register("android", androidGen)

	gen, err := registry.Get("android")
	if err != nil {
		t.Fatalf("Failed to get registered generator: %v", err)
	}
	if gen == nil {
		t.Fatal("Retrieved generator is nil")
	}
}

func TestRegistryGetNonExistent(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent generator, got nil")
	}
}

func TestAndroidGeneratorSupportsComponent(t *testing.T) {
	gen := NewAndroidGenerator()

	if !gen.SupportsComponent("android") {
		t.Error("Android generator should support 'android' component type")
	}

	if gen.SupportsComponent("ios") {
		t.Error("Android generator should not support 'ios' component type")
	}
}

func TestAndroidGeneratorGetRequiredManualSteps(t *testing.T) {
	gen := NewAndroidGenerator()

	steps := gen.GetRequiredManualSteps("android")
	if len(steps) == 0 {
		t.Error("Expected manual steps, got none")
	}

	// Check that steps mention Android Studio
	foundAndroidStudio := false
	for _, step := range steps {
		if len(step) > 0 {
			foundAndroidStudio = true
			break
		}
	}
	if !foundAndroidStudio {
		t.Error("Expected manual steps to include setup instructions")
	}
}

func TestAndroidGeneratorGenerate(t *testing.T) {
	gen := NewAndroidGenerator()
	ctx := context.Background()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "android-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	spec := &models.FallbackSpec{
		ComponentType: "android",
		TargetDir:     tmpDir,
		Config: map[string]interface{}{
			"package":  "com.test.app",
			"app_name": "TestApp",
		},
	}

	result, err := gen.Generate(ctx, spec)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Generation was not successful: %v", result.Error)
	}

	if result.Type != "android" {
		t.Errorf("Expected type 'android', got '%s'", result.Type)
	}

	if result.Method != "fallback" {
		t.Errorf("Expected method 'fallback', got '%s'", result.Method)
	}

	// Check that key files were created
	expectedFiles := []string{
		"settings.gradle",
		"build.gradle",
		"app/build.gradle",
		"app/src/main/AndroidManifest.xml",
		"README.md",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}
}

func TestIOSGeneratorSupportsComponent(t *testing.T) {
	gen := NewIOSGenerator()

	if !gen.SupportsComponent("ios") {
		t.Error("iOS generator should support 'ios' component type")
	}

	if gen.SupportsComponent("android") {
		t.Error("iOS generator should not support 'android' component type")
	}
}

func TestIOSGeneratorGetRequiredManualSteps(t *testing.T) {
	gen := NewIOSGenerator()

	steps := gen.GetRequiredManualSteps("ios")
	if len(steps) == 0 {
		t.Error("Expected manual steps, got none")
	}
}

func TestIOSGeneratorGenerate(t *testing.T) {
	gen := NewIOSGenerator()
	ctx := context.Background()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ios-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	spec := &models.FallbackSpec{
		ComponentType: "ios",
		TargetDir:     tmpDir,
		Config: map[string]interface{}{
			"bundle_id":    "com.test.app",
			"app_name":     "TestApp",
			"organization": "TestOrg",
		},
	}

	result, err := gen.Generate(ctx, spec)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Generation was not successful: %v", result.Error)
	}

	if result.Type != "ios" {
		t.Errorf("Expected type 'ios', got '%s'", result.Type)
	}

	if result.Method != "fallback" {
		t.Errorf("Expected method 'fallback', got '%s'", result.Method)
	}

	// Check that key files were created
	expectedFiles := []string{
		"TestApp.xcodeproj/project.pbxproj",
		"TestApp/TestAppApp.swift",
		"TestApp/ContentView.swift",
		"README.md",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", file)
		}
	}
}
