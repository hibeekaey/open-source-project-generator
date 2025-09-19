// Package ui provides the main interactive template selection interface.
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InteractiveTemplateSelection provides the main interface for interactive template selection
type InteractiveTemplateSelection struct {
	ui                     interfaces.InteractiveUIInterface
	templateManager        interfaces.TemplateManager
	templateSelector       *TemplateSelector
	templatePreviewManager *TemplatePreviewManager
	logger                 interfaces.Logger
}

// TemplateSelectionResult contains the result of interactive template selection
type TemplateSelectionResult struct {
	Selections      []TemplateSelection
	PreviewAccepted bool
	UserCancelled   bool
}

// NewInteractiveTemplateSelection creates a new interactive template selection interface
func NewInteractiveTemplateSelection(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *InteractiveTemplateSelection {
	templateSelector := NewTemplateSelector(ui, templateManager, logger)
	templatePreviewManager := NewTemplatePreviewManager(ui, templateManager, logger)

	return &InteractiveTemplateSelection{
		ui:                     ui,
		templateManager:        templateManager,
		templateSelector:       templateSelector,
		templatePreviewManager: templatePreviewManager,
		logger:                 logger,
	}
}

// SelectTemplatesInteractively runs the complete interactive template selection flow
func (its *InteractiveTemplateSelection) SelectTemplatesInteractively(ctx context.Context, config *models.ProjectConfig) (*TemplateSelectionResult, error) {
	// Start UI session
	sessionConfig := interfaces.SessionConfig{
		Title:       "Interactive Template Selection",
		Description: "Select and configure templates for your project",
		AutoSave:    true,
		Metadata: map[string]interface{}{
			"project_name": config.Name,
			"step":         "template_selection",
		},
	}

	session, err := its.ui.StartSession(ctx, sessionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start UI session: %w", err)
	}
	defer func() {
		if endErr := its.ui.EndSession(ctx, session); endErr != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to end UI session: %v\n", endErr)
		}
	}()

	// Main selection loop
	for {
		// Step 1: Template Selection
		selections, err := its.templateSelector.SelectTemplatesInteractively(ctx)
		if err != nil {
			return nil, fmt.Errorf("template selection failed: %w", err)
		}

		if len(selections) == 0 {
			return &TemplateSelectionResult{
				UserCancelled: true,
			}, nil
		}

		// Step 2: Validate selections
		if err := its.templateSelector.ValidateTemplateSelections(selections); err != nil {
			// Show validation error and allow user to retry
			errorConfig := interfaces.ErrorConfig{
				Title:     "Template Selection Error",
				Message:   err.Error(),
				ErrorType: "validation",
				Suggestions: []string{
					"Review your template selections",
					"Check template dependencies",
					"Select compatible templates",
				},
				AllowRetry: true,
				AllowBack:  true,
				AllowQuit:  true,
			}

			errorResult, err := its.ui.ShowError(ctx, errorConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to show validation error: %w", err)
			}

			if errorResult.Cancelled || errorResult.Action == "quit" {
				return &TemplateSelectionResult{UserCancelled: true}, nil
			}

			if errorResult.Action == "back" || errorResult.Action == "retry" {
				continue // Go back to template selection
			}
		}

		// Step 3: Show individual template previews (optional)
		showPreviews, err := its.promptForIndividualPreviews(ctx, selections)
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for previews: %w", err)
		}

		if showPreviews {
			for _, selection := range selections {
				err := its.templatePreviewManager.PreviewIndividualTemplate(ctx, selection.Template, config)
				if err != nil {
					its.logger.WarnWithFields("Failed to show individual template preview", map[string]interface{}{
						"template": selection.Template.Name,
						"error":    err.Error(),
					})
				}
			}
		}

		// Step 4: Show combined preview
		combinedPreview, err := its.templatePreviewManager.PreviewCombinedTemplates(ctx, selections, config)
		if err != nil {
			return nil, fmt.Errorf("failed to generate combined preview: %w", err)
		}

		// Step 5: Confirm selections
		_, action, err := its.confirmTemplateSelections(ctx, selections, combinedPreview)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm selections: %w", err)
		}

		switch action {
		case "confirm":
			return &TemplateSelectionResult{
				Selections:      selections,
				PreviewAccepted: true,
			}, nil

		case "back", "modify":
			continue // Go back to template selection

		case "quit", "cancel":
			return &TemplateSelectionResult{
				UserCancelled: true,
			}, nil

		default:
			continue // Default to going back
		}
	}
}

// promptForIndividualPreviews asks if user wants to see individual template previews
func (its *InteractiveTemplateSelection) promptForIndividualPreviews(ctx context.Context, selections []TemplateSelection) (bool, error) {
	if len(selections) <= 1 {
		return false, nil // No need for individual previews with single template
	}

	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Would you like to preview individual templates before seeing the combined result?",
		Description:  "This will show you what each template generates separately",
		DefaultValue: false,
		AllowBack:    false,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText:     "Individual previews show what each template generates on its own, while the combined preview shows how they work together.",
	}

	result, err := its.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, err
	}

	if result.Cancelled {
		return false, fmt.Errorf("user cancelled")
	}

	return result.Confirmed, nil
}

// confirmTemplateSelections asks user to confirm their template selections
func (its *InteractiveTemplateSelection) confirmTemplateSelections(ctx context.Context, selections []TemplateSelection, preview *CombinedTemplatePreview) (bool, string, error) {
	// Build confirmation message
	message := fmt.Sprintf("You have selected %d template(s):", len(selections))
	for _, selection := range selections {
		message += fmt.Sprintf("\n  • %s (%s)", selection.Template.DisplayName, selection.Template.Category)
	}

	message += fmt.Sprintf("\n\nThis will generate approximately %d files (%s)",
		preview.TotalFiles, formatBytes(preview.EstimatedSize))

	// Add warnings if any
	if len(preview.Warnings) > 0 {
		message += "\n\nWarnings:"
		for _, warning := range preview.Warnings {
			message += fmt.Sprintf("\n  ⚠️ %s", warning)
		}
	}

	// Add conflicts if any
	errorConflicts := 0
	for _, conflict := range preview.FileConflicts {
		if conflict.Severity == "error" {
			errorConflicts++
		}
	}

	if errorConflicts > 0 {
		message += fmt.Sprintf("\n\n❌ %d file conflicts require manual resolution", errorConflicts)
	}

	// Show menu with options
	options := []interfaces.MenuOption{
		{
			Label:       "Proceed with Generation",
			Description: "Continue with the selected templates",
			Value:       "confirm",
			Icon:        "✅",
		},
		{
			Label:       "Modify Selection",
			Description: "Go back and change template selection",
			Value:       "modify",
			Icon:        "🔄",
		},
		{
			Label:       "Show Detailed Preview",
			Description: "View detailed file structure and conflicts",
			Value:       "preview",
			Icon:        "🔍",
		},
		{
			Label:       "Cancel",
			Description: "Cancel template selection",
			Value:       "cancel",
			Icon:        "❌",
		},
	}

	menuConfig := interfaces.MenuConfig{
		Title:       "Confirm Template Selection",
		Description: message,
		Options:     options,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Review your template selection and choose how to proceed. You can modify your selection, view more details, or proceed with generation.",
	}

	for {
		result, err := its.ui.ShowMenu(ctx, menuConfig)
		if err != nil {
			return false, "", err
		}

		if result.Cancelled {
			return false, "cancel", nil
		}

		action, ok := result.SelectedValue.(string)
		if !ok {
			continue
		}

		switch action {
		case "confirm":
			return true, "confirm", nil

		case "modify":
			return false, "modify", nil

		case "cancel":
			return false, "cancel", nil

		case "preview":
			// Show detailed preview and continue loop
			err := its.showDetailedPreview(ctx, preview)
			if err != nil {
				its.logger.WarnWithFields("Failed to show detailed preview", map[string]interface{}{
					"error": err.Error(),
				})
			}
			continue

		default:
			continue
		}
	}
}

// showDetailedPreview shows a detailed preview with file conflicts and structure
func (its *InteractiveTemplateSelection) showDetailedPreview(ctx context.Context, preview *CombinedTemplatePreview) error {
	// Show conflicts in detail
	if len(preview.FileConflicts) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("FILE CONFLICTS ANALYSIS")
		fmt.Println(strings.Repeat("=", 60))

		for _, conflict := range preview.FileConflicts {
			icon := "⚠️"
			switch conflict.Severity {
			case "error":
				icon = "❌"
			case "info":
				icon = "ℹ️"
			}

			fmt.Printf("\n%s %s\n", icon, conflict.Path)
			fmt.Printf("   Templates: %s\n", strings.Join(conflict.Templates, ", "))
			fmt.Printf("   Issue: %s\n", conflict.Message)
			if conflict.Resolvable {
				fmt.Printf("   Resolution: Automatic\n")
			} else {
				fmt.Printf("   Resolution: Manual intervention required\n")
			}
		}
	}

	// Show dependency analysis
	if len(preview.Dependencies) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("DEPENDENCY ANALYSIS")
		fmt.Println(strings.Repeat("=", 60))

		selectedTemplates := make(map[string]bool)
		for _, selection := range preview.Templates {
			selectedTemplates[selection.Template.Name] = true
		}

		for _, dep := range preview.Dependencies {
			status := "✅ Satisfied"
			if !selectedTemplates[dep] {
				status = "❌ Missing"
			}
			fmt.Printf("  %s: %s\n", dep, status)
		}
	}

	// Show size breakdown
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PROJECT SIZE BREAKDOWN")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Files: %d\n", preview.TotalFiles)
	fmt.Printf("Estimated Size: %s\n", formatBytes(preview.EstimatedSize))
	fmt.Printf("Templates: %d\n", len(preview.Templates))

	fmt.Print("\nPress Enter to continue...")
	var input string
	_, _ = fmt.Scanln(&input)

	return nil
}

// GetTemplateInfo provides detailed information about a specific template
func (its *InteractiveTemplateSelection) GetTemplateInfo(ctx context.Context, templateName string) error {
	templateInfo, err := its.templateManager.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template info: %w", err)
	}

	return its.templateSelector.ShowTemplateInfo(ctx, *templateInfo)
}

// ListAvailableTemplates shows all available templates organized by category
func (its *InteractiveTemplateSelection) ListAvailableTemplates(ctx context.Context) error {
	allTemplates, err := its.templateManager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	categories := its.templateSelector.organizeTemplatesByCategory(allTemplates)

	// Display as a table
	var headers []string
	var rows [][]string

	headers = []string{"Category", "Template", "Technology", "Description"}

	for _, category := range categories {
		for _, template := range category.Templates {
			row := []string{
				category.DisplayName,
				template.DisplayName,
				template.Technology,
				template.Description,
			}
			rows = append(rows, row)
		}
	}

	tableConfig := interfaces.TableConfig{
		Title:      "Available Templates",
		Headers:    headers,
		Rows:       rows,
		Pagination: true,
		PageSize:   20,
		Sortable:   true,
		Searchable: true,
	}

	return its.ui.ShowTable(ctx, tableConfig)
}
