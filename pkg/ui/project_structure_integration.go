// Package ui provides integration functionality for project structure preview and navigation.
package ui

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectStructureIntegration provides a complete integration of project structure preview and navigation
type ProjectStructureIntegration struct {
	ui                interfaces.InteractiveUIInterface
	templateManager   interfaces.TemplateManager
	logger            interfaces.Logger
	previewGenerator  *ProjectStructurePreviewGenerator
	navigationManager *PreviewNavigationManager
}

// ProjectStructureResult represents the final result of the project structure workflow
type ProjectStructureResult struct {
	Preview         *ProjectStructurePreview `json:"preview"`
	FinalConfig     *models.ProjectConfig    `json:"final_config"`
	FinalSelections []TemplateSelection      `json:"final_selections"`
	FinalOutputDir  string                   `json:"final_output_dir"`
	UserConfirmed   bool                     `json:"user_confirmed"`
	Action          string                   `json:"action"` // "confirmed", "cancelled", "back"
}

// NewProjectStructureIntegration creates a new project structure integration
func NewProjectStructureIntegration(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *ProjectStructureIntegration {
	previewGenerator := NewProjectStructurePreviewGenerator(ui, templateManager, logger)
	navigationManager := NewPreviewNavigationManager(ui, previewGenerator, logger)

	return &ProjectStructureIntegration{
		ui:                ui,
		templateManager:   templateManager,
		logger:            logger,
		previewGenerator:  previewGenerator,
		navigationManager: navigationManager,
	}
}

// ShowProjectStructurePreviewWithNavigation displays the complete project structure preview workflow
func (psi *ProjectStructureIntegration) ShowProjectStructurePreviewWithNavigation(ctx context.Context, config *models.ProjectConfig, selections []TemplateSelection, outputDir string) (*ProjectStructureResult, error) {
	psi.logger.InfoWithFields("Starting project structure preview workflow", map[string]interface{}{
		"project_name":     config.Name,
		"template_count":   len(selections),
		"output_directory": outputDir,
	})

	// Show the preview with navigation
	navigationResult, err := psi.navigationManager.ShowPreviewWithNavigation(ctx, config, selections, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to show preview with navigation: %w", err)
	}

	// Build the final result
	result := &ProjectStructureResult{
		Preview:         psi.navigationManager.GetCurrentPreview(),
		FinalConfig:     psi.navigationManager.GetModifiedConfig(),
		FinalSelections: psi.navigationManager.GetModifiedSelections(),
		FinalOutputDir:  psi.navigationManager.GetModifiedOutputDir(),
		UserConfirmed:   navigationResult.Confirmed,
		Action:          navigationResult.Action,
	}

	psi.logger.InfoWithFields("Project structure preview workflow completed", map[string]interface{}{
		"action":       result.Action,
		"confirmed":    result.UserConfirmed,
		"final_config": result.FinalConfig.Name,
		"final_output": result.FinalOutputDir,
	})

	return result, nil
}

// GenerateProjectStructurePreviewOnly generates a preview without navigation (for display-only purposes)
func (psi *ProjectStructureIntegration) GenerateProjectStructurePreviewOnly(ctx context.Context, config *models.ProjectConfig, selections []TemplateSelection, outputDir string) (*ProjectStructurePreview, error) {
	psi.logger.InfoWithFields("Generating project structure preview only", map[string]interface{}{
		"project_name":     config.Name,
		"template_count":   len(selections),
		"output_directory": outputDir,
	})

	preview, err := psi.previewGenerator.GenerateProjectStructurePreview(ctx, config, selections, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project structure preview: %w", err)
	}

	return preview, nil
}

// DisplayProjectStructurePreviewOnly displays a preview without navigation options
func (psi *ProjectStructureIntegration) DisplayProjectStructurePreviewOnly(ctx context.Context, preview *ProjectStructurePreview) error {
	psi.logger.InfoWithFields("Displaying project structure preview", map[string]interface{}{
		"project_name": preview.ProjectName,
		"directories":  preview.Summary.TotalDirectories,
		"files":        preview.Summary.TotalFiles,
	})

	err := psi.previewGenerator.DisplayProjectStructurePreview(ctx, preview)
	if err != nil {
		return fmt.Errorf("failed to display project structure preview: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates the project structure for potential issues
func (psi *ProjectStructureIntegration) ValidateProjectStructure(ctx context.Context, preview *ProjectStructurePreview) (*ProjectStructureValidationResult, error) {
	psi.logger.InfoWithFields("Validating project structure", map[string]interface{}{
		"project_name": preview.ProjectName,
		"conflicts":    len(preview.Conflicts),
		"warnings":     len(preview.Warnings),
	})

	validation := &ProjectStructureValidationResult{
		Valid:       true,
		Issues:      []ValidationIssue{},
		Warnings:    []ValidationIssue{},
		Suggestions: []string{},
	}

	// Check for critical issues
	for _, conflict := range preview.Conflicts {
		if conflict.Severity == "error" {
			validation.Valid = false
			validation.Issues = append(validation.Issues, ValidationIssue{
				Type:       "conflict",
				Severity:   "error",
				Message:    conflict.Message,
				Path:       conflict.Path,
				Resolvable: conflict.Resolvable,
			})
		} else {
			validation.Warnings = append(validation.Warnings, ValidationIssue{
				Type:       "conflict",
				Severity:   conflict.Severity,
				Message:    conflict.Message,
				Path:       conflict.Path,
				Resolvable: conflict.Resolvable,
			})
		}
	}

	// Add warnings as validation warnings
	for _, warning := range preview.Warnings {
		validation.Warnings = append(validation.Warnings, ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  warning,
		})
	}

	// Generate suggestions
	if preview.Summary.TotalFiles > 1000 {
		validation.Suggestions = append(validation.Suggestions, "Consider using fewer templates to reduce project complexity")
	}

	if preview.Summary.EstimatedSize > 100*1024*1024 { // 100MB
		validation.Suggestions = append(validation.Suggestions, "Large project size detected - ensure adequate disk space")
	}

	if len(preview.Conflicts) > 0 {
		validation.Suggestions = append(validation.Suggestions, "Review file conflicts before generation to avoid issues")
	}

	psi.logger.InfoWithFields("Project structure validation completed", map[string]interface{}{
		"valid":       validation.Valid,
		"issues":      len(validation.Issues),
		"warnings":    len(validation.Warnings),
		"suggestions": len(validation.Suggestions),
	})

	return validation, nil
}

// ProjectStructureValidationResult represents the result of project structure validation
type ProjectStructureValidationResult struct {
	Valid       bool              `json:"valid"`
	Issues      []ValidationIssue `json:"issues"`
	Warnings    []ValidationIssue `json:"warnings"`
	Suggestions []string          `json:"suggestions"`
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
	Type       string `json:"type"`     // "conflict", "warning", "error"
	Severity   string `json:"severity"` // "error", "warning", "info"
	Message    string `json:"message"`
	Path       string `json:"path,omitempty"`
	Resolvable bool   `json:"resolvable"`
}

// GetPreviewGenerator returns the preview generator for advanced usage
func (psi *ProjectStructureIntegration) GetPreviewGenerator() *ProjectStructurePreviewGenerator {
	return psi.previewGenerator
}

// GetNavigationManager returns the navigation manager for advanced usage
func (psi *ProjectStructureIntegration) GetNavigationManager() *PreviewNavigationManager {
	return psi.navigationManager
}

// Example usage function for documentation purposes
func ExampleProjectStructureIntegrationUsage(ctx context.Context, ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) error {
	// Create the integration
	integration := NewProjectStructureIntegration(ui, templateManager, logger)

	// Example project configuration
	config := &models.ProjectConfig{
		Name:        "my-awesome-project",
		Description: "An awesome full-stack project",
		Author:      "Developer Name",
		License:     "MIT",
	}

	// Example template selections
	selections := []TemplateSelection{
		{
			Template: interfaces.TemplateInfo{
				Name:        "nextjs-app",
				DisplayName: "Next.js Application",
				Category:    "frontend",
				Technology:  "React",
			},
			Selected: true,
		},
		{
			Template: interfaces.TemplateInfo{
				Name:        "go-gin",
				DisplayName: "Go Gin API",
				Category:    "backend",
				Technology:  "Go",
			},
			Selected: true,
		},
	}

	outputDir := "./output/my-awesome-project"

	// Show the complete workflow
	result, err := integration.ShowProjectStructurePreviewWithNavigation(ctx, config, selections, outputDir)
	if err != nil {
		return fmt.Errorf("failed to show project structure preview: %w", err)
	}

	// Handle the result
	switch result.Action {
	case "confirm":
		logger.InfoWithFields("User confirmed project structure", map[string]interface{}{
			"project_name": result.FinalConfig.Name,
			"output_dir":   result.FinalOutputDir,
		})
		// Proceed with project generation...

	case "back":
		logger.Info("User chose to go back to previous step")
		// Return to previous workflow step...

	case "quit":
		logger.Info("User chose to quit")
		// Handle quit...

	default:
		logger.WarnWithFields("Unknown action", map[string]interface{}{
			"action": result.Action,
		})
	}

	return nil
}
