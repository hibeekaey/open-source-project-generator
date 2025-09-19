package ui

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestProjectStructureIntegration_Creation(t *testing.T) {
	// Test that the integration can be created
	integration := &ProjectStructureIntegration{}

	// Note: integration is created with &ProjectStructureIntegration{} so it cannot be nil

	// Test that components can be set
	previewGenerator := &ProjectStructurePreviewGenerator{}
	navigationManager := &PreviewNavigationManager{}

	integration.previewGenerator = previewGenerator
	integration.navigationManager = navigationManager

	if integration.GetPreviewGenerator() != previewGenerator {
		t.Error("Expected preview generator to match")
	}

	if integration.GetNavigationManager() != navigationManager {
		t.Error("Expected navigation manager to match")
	}
}

func TestProjectStructureResult_Structure(t *testing.T) {
	// Test the project structure result structure
	config := &models.ProjectConfig{
		Name: "test-project",
	}

	selections := []TemplateSelection{
		{
			Template: interfaces.TemplateInfo{
				Name: "test-template",
			},
			Selected: true,
		},
	}

	preview := &ProjectStructurePreview{
		ProjectName: "test-project",
	}

	result := &ProjectStructureResult{
		Preview:         preview,
		FinalConfig:     config,
		FinalSelections: selections,
		FinalOutputDir:  "/tmp/test",
		UserConfirmed:   true,
		Action:          "confirmed",
	}

	if result.Preview != preview {
		t.Error("Expected preview to match")
	}

	if result.FinalConfig != config {
		t.Error("Expected final config to match")
	}

	if len(result.FinalSelections) != 1 {
		t.Error("Expected 1 final selection")
	}

	if result.FinalOutputDir != "/tmp/test" {
		t.Error("Expected final output dir to match")
	}

	if !result.UserConfirmed {
		t.Error("Expected user confirmed to be true")
	}

	if result.Action != "confirmed" {
		t.Error("Expected action to be 'confirmed'")
	}
}

func TestProjectStructureValidationResult_Structure(t *testing.T) {
	// Test the validation result structure
	validation := &ProjectStructureValidationResult{
		Valid: true,
		Issues: []ValidationIssue{
			{
				Type:     "conflict",
				Severity: "error",
				Message:  "Test conflict",
				Path:     "/test/path",
			},
		},
		Warnings: []ValidationIssue{
			{
				Type:     "warning",
				Severity: "warning",
				Message:  "Test warning",
			},
		},
		Suggestions: []string{"Test suggestion"},
	}

	if !validation.Valid {
		t.Error("Expected validation to be valid")
	}

	if len(validation.Issues) != 1 {
		t.Error("Expected 1 issue")
	}

	if validation.Issues[0].Type != "conflict" {
		t.Error("Expected issue type to be 'conflict'")
	}

	if validation.Issues[0].Severity != "error" {
		t.Error("Expected issue severity to be 'error'")
	}

	if len(validation.Warnings) != 1 {
		t.Error("Expected 1 warning")
	}

	if len(validation.Suggestions) != 1 {
		t.Error("Expected 1 suggestion")
	}
}

func TestValidationIssue_Structure(t *testing.T) {
	// Test the validation issue structure
	issue := ValidationIssue{
		Type:       "conflict",
		Severity:   "error",
		Message:    "Test message",
		Path:       "/test/path",
		Resolvable: true,
	}

	if issue.Type != "conflict" {
		t.Errorf("Expected type 'conflict', got '%s'", issue.Type)
	}

	if issue.Severity != "error" {
		t.Errorf("Expected severity 'error', got '%s'", issue.Severity)
	}

	if issue.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", issue.Message)
	}

	if issue.Path != "/test/path" {
		t.Errorf("Expected path '/test/path', got '%s'", issue.Path)
	}

	if !issue.Resolvable {
		t.Error("Expected resolvable to be true")
	}
}
