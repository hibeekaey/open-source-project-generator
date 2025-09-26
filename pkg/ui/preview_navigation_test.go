package ui

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestPreviewNavigationManager_Basic(t *testing.T) {
	// Create a basic preview navigation manager
	manager := &PreviewNavigationManager{}

	// Test that the manager can be created
	// Note: manager is created with &PreviewNavigationManager{} so it cannot be nil

	// Test setting original values
	config := &models.ProjectConfig{
		Name:        "test-project",
		Description: "A test project",
	}

	selections := []TemplateSelection{
		{
			Template: interfaces.TemplateInfo{
				Name:        "nextjs-app",
				DisplayName: "Next.js Application",
				Category:    "frontend",
			},
			Selected: true,
		},
	}

	outputDir := "/tmp/test-output"

	manager.originalConfig = config
	manager.originalSelections = selections
	manager.originalOutputDir = outputDir

	// Test getters
	if manager.GetModifiedConfig() != config {
		t.Error("Expected modified config to match original")
	}

	if len(manager.GetModifiedSelections()) != len(selections) {
		t.Error("Expected modified selections to match original")
	}

	if manager.GetModifiedOutputDir() != outputDir {
		t.Error("Expected modified output dir to match original")
	}
}

func TestPreviewNavigationResult_Structure(t *testing.T) {
	// Test the preview navigation result structure
	result := &PreviewNavigationResult{
		Action:    "confirm",
		Confirmed: true,
		Cancelled: false,
	}

	if result.Action != "confirm" {
		t.Errorf("Expected action 'confirm', got '%s'", result.Action)
	}

	if !result.Confirmed {
		t.Error("Expected confirmed to be true")
	}

	if result.Cancelled {
		t.Error("Expected cancelled to be false")
	}
}

func TestPreviewNavigationManager_GetCurrentPreview(t *testing.T) {
	// Test getting current preview
	manager := &PreviewNavigationManager{}

	// Initially should be nil
	if manager.GetCurrentPreview() != nil {
		t.Error("Expected current preview to be nil initially")
	}

	// Set a preview
	preview := &ProjectStructurePreview{
		ProjectName:     "test-project",
		OutputDirectory: "/tmp/test",
	}

	manager.currentPreview = preview

	// Should return the set preview
	if manager.GetCurrentPreview() != preview {
		t.Error("Expected current preview to match set preview")
	}

	if manager.GetCurrentPreview().ProjectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", manager.GetCurrentPreview().ProjectName)
	}
}

func TestPreviewNavigationManager_ModificationTracking(t *testing.T) {
	// Test that modifications are tracked correctly
	manager := &PreviewNavigationManager{}

	// Set initial values
	originalConfig := &models.ProjectConfig{
		Name: "original-project",
	}

	originalSelections := []TemplateSelection{
		{
			Template: interfaces.TemplateInfo{
				Name: "original-template",
			},
			Selected: true,
		},
	}

	originalOutputDir := "/original/path"

	manager.originalConfig = originalConfig
	manager.originalSelections = originalSelections
	manager.originalOutputDir = originalOutputDir

	// Verify original values are preserved
	if manager.GetModifiedConfig().Name != "original-project" {
		t.Error("Expected original config to be preserved")
	}

	if len(manager.GetModifiedSelections()) != 1 {
		t.Error("Expected original selections to be preserved")
	}

	if manager.GetModifiedOutputDir() != "/original/path" {
		t.Error("Expected original output dir to be preserved")
	}

	// Simulate modification
	newConfig := &models.ProjectConfig{
		Name: "modified-project",
	}

	manager.originalConfig = newConfig

	// Verify modification is reflected
	if manager.GetModifiedConfig().Name != "modified-project" {
		t.Error("Expected modified config to be updated")
	}
}
