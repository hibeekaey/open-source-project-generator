// Package ui provides preview navigation functionality for interactive CLI generation.
package ui

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// PreviewNavigationManager handles navigation and modification of project structure previews
type PreviewNavigationManager struct {
	ui                 interfaces.InteractiveUIInterface
	previewGenerator   *ProjectStructurePreviewGenerator
	logger             interfaces.Logger
	currentPreview     *ProjectStructurePreview
	originalSelections []TemplateSelection
	originalConfig     *models.ProjectConfig
	originalOutputDir  string
}

// PreviewNavigationResult represents the result of preview navigation
type PreviewNavigationResult struct {
	Action             string                `json:"action"` // "confirm", "modify", "back", "quit"
	ModifiedSelections []TemplateSelection   `json:"modified_selections,omitempty"`
	ModifiedConfig     *models.ProjectConfig `json:"modified_config,omitempty"`
	ModifiedOutputDir  string                `json:"modified_output_dir,omitempty"`
	Confirmed          bool                  `json:"confirmed"`
	Cancelled          bool                  `json:"cancelled"`
}

// NewPreviewNavigationManager creates a new preview navigation manager
func NewPreviewNavigationManager(ui interfaces.InteractiveUIInterface, previewGenerator *ProjectStructurePreviewGenerator, logger interfaces.Logger) *PreviewNavigationManager {
	return &PreviewNavigationManager{
		ui:               ui,
		previewGenerator: previewGenerator,
		logger:           logger,
	}
}

// ShowPreviewWithNavigation displays the project structure preview with navigation options
func (pnm *PreviewNavigationManager) ShowPreviewWithNavigation(ctx context.Context, config *models.ProjectConfig, selections []TemplateSelection, outputDir string) (*PreviewNavigationResult, error) {
	// Store original values for potential modifications
	pnm.originalConfig = config
	pnm.originalSelections = selections
	pnm.originalOutputDir = outputDir

	// Generate initial preview
	preview, err := pnm.previewGenerator.GenerateProjectStructurePreview(ctx, config, selections, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project structure preview: %w", err)
	}

	pnm.currentPreview = preview

	// Display preview and handle navigation
	return pnm.handlePreviewNavigation(ctx)
}

// handlePreviewNavigation manages the preview navigation loop
func (pnm *PreviewNavigationManager) handlePreviewNavigation(ctx context.Context) (*PreviewNavigationResult, error) {
	for {
		// Display the current preview
		err := pnm.previewGenerator.DisplayProjectStructurePreview(ctx, pnm.currentPreview)
		if err != nil {
			return nil, fmt.Errorf("failed to display preview: %w", err)
		}

		// Show navigation menu
		result, err := pnm.showNavigationMenu(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to show navigation menu: %w", err)
		}

		// Handle the selected action
		switch result.Action {
		case "confirm":
			return &PreviewNavigationResult{
				Action:    "confirm",
				Confirmed: true,
			}, nil

		case "modify_templates":
			modified, err := pnm.handleTemplateModification(ctx)
			if err != nil {
				pnm.logger.WarnWithFields("Template modification failed", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			if modified {
				// Regenerate preview with modified selections
				err = pnm.regeneratePreview(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to regenerate preview: %w", err)
				}
			}

		case "modify_config":
			modified, err := pnm.handleConfigModification(ctx)
			if err != nil {
				pnm.logger.WarnWithFields("Config modification failed", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			if modified {
				// Regenerate preview with modified config
				err = pnm.regeneratePreview(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to regenerate preview: %w", err)
				}
			}

		case "modify_output":
			modified, err := pnm.handleOutputDirModification(ctx)
			if err != nil {
				pnm.logger.WarnWithFields("Output directory modification failed", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			if modified {
				// Regenerate preview with modified output directory
				err = pnm.regeneratePreview(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to regenerate preview: %w", err)
				}
			}

		case "back":
			return &PreviewNavigationResult{
				Action:    "back",
				Cancelled: true,
			}, nil

		case "quit":
			return &PreviewNavigationResult{
				Action:    "quit",
				Cancelled: true,
			}, nil

		default:
			pnm.logger.WarnWithFields("Unknown navigation action", map[string]interface{}{
				"action": result.Action,
			})
		}
	}
}

// showNavigationMenu displays the navigation options menu
func (pnm *PreviewNavigationManager) showNavigationMenu(ctx context.Context) (*interfaces.MenuResult, error) {
	menuConfig := interfaces.MenuConfig{
		Title:       "Preview Navigation",
		Description: "Choose an action to proceed with the project generation",
		Options: []interfaces.MenuOption{
			{
				Label:       "✅ Proceed with Generation",
				Description: "Confirm the preview and start generating the project",
				Value:       "confirm",
				Icon:        "✅",
				Shortcut:    "c",
			},
			{
				Label:       "🔄 Modify Template Selection",
				Description: "Go back to modify selected templates",
				Value:       "modify_templates",
				Icon:        "🔄",
				Shortcut:    "t",
			},
			{
				Label:       "⚙️ Modify Project Configuration",
				Description: "Go back to modify project settings",
				Value:       "modify_config",
				Icon:        "⚙️",
				Shortcut:    "p",
			},
			{
				Label:       "📁 Modify Output Directory",
				Description: "Change the output directory path",
				Value:       "modify_output",
				Icon:        "📁",
				Shortcut:    "o",
			},
			{
				Label:       "🔙 Back to Previous Step",
				Description: "Return to the previous configuration step",
				Value:       "back",
				Icon:        "🔙",
				Shortcut:    "b",
			},
			{
				Label:       "❌ Quit",
				Description: "Exit the project generation process",
				Value:       "quit",
				Icon:        "❌",
				Shortcut:    "q",
			},
		},
		DefaultItem: 0,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Use arrow keys to navigate, Enter to select, or use shortcut keys. You can modify any aspect of the project configuration before proceeding with generation.",
	}

	return pnm.ui.ShowMenu(ctx, menuConfig)
}

// handleTemplateModification handles modification of template selections
func (pnm *PreviewNavigationManager) handleTemplateModification(ctx context.Context) (bool, error) {
	// Show confirmation dialog
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Modify Template Selection",
		Description:  "This will take you back to the template selection step. Any current preview will be regenerated.",
		DefaultValue: false,
		YesLabel:     "Modify Templates",
		NoLabel:      "Keep Current",
		AllowBack:    true,
		ShowHelp:     true,
		HelpText:     "Choose 'Modify Templates' to go back and change your template selection, or 'Keep Current' to stay with the current preview.",
	}

	result, err := pnm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get confirmation: %w", err)
	}

	if result.Cancelled || !result.Confirmed {
		return false, nil
	}

	// For now, we'll just return true to indicate modification was requested
	// In a full implementation, this would integrate with the template selection UI
	pnm.logger.InfoWithFields("Template modification requested", map[string]interface{}{
		"current_templates": len(pnm.originalSelections),
	})

	return true, nil
}

// handleConfigModification handles modification of project configuration
func (pnm *PreviewNavigationManager) handleConfigModification(ctx context.Context) (bool, error) {
	// Show confirmation dialog
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Modify Project Configuration",
		Description:  "This will take you back to the project configuration step. The preview will be regenerated with new settings.",
		DefaultValue: false,
		YesLabel:     "Modify Config",
		NoLabel:      "Keep Current",
		AllowBack:    true,
		ShowHelp:     true,
		HelpText:     "Choose 'Modify Config' to go back and change project settings like name, description, etc., or 'Keep Current' to stay with the current configuration.",
	}

	result, err := pnm.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get confirmation: %w", err)
	}

	if result.Cancelled || !result.Confirmed {
		return false, nil
	}

	// For now, we'll just return true to indicate modification was requested
	// In a full implementation, this would integrate with the config collection UI
	pnm.logger.InfoWithFields("Config modification requested", map[string]interface{}{
		"current_project": pnm.originalConfig.Name,
	})

	return true, nil
}

// handleOutputDirModification handles modification of output directory
func (pnm *PreviewNavigationManager) handleOutputDirModification(ctx context.Context) (bool, error) {
	textConfig := interfaces.TextPromptConfig{
		Prompt:       "Output Directory",
		Description:  "Enter the new output directory path for project generation",
		DefaultValue: pnm.originalOutputDir,
		Required:     true,
		AllowBack:    true,
		ShowHelp:     true,
		HelpText:     "Specify the directory where the project files will be generated. Use an absolute path or relative to the current directory.",
		Validator: func(input string) error {
			if input == "" {
				return fmt.Errorf("output directory cannot be empty")
			}
			// Additional validation could be added here
			return nil
		},
	}

	result, err := pnm.ui.PromptText(ctx, textConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get new output directory: %w", err)
	}

	if result.Cancelled {
		return false, nil
	}

	// Check if the directory actually changed
	if result.Value == pnm.originalOutputDir {
		return false, nil
	}

	// Update the output directory
	pnm.originalOutputDir = result.Value
	pnm.logger.InfoWithFields("Output directory modified", map[string]interface{}{
		"new_output_dir": result.Value,
	})

	return true, nil
}

// regeneratePreview regenerates the preview with current settings
func (pnm *PreviewNavigationManager) regeneratePreview(ctx context.Context) error {
	pnm.logger.InfoWithFields("Regenerating project structure preview", map[string]interface{}{
		"project_name":   pnm.originalConfig.Name,
		"output_dir":     pnm.originalOutputDir,
		"template_count": len(pnm.originalSelections),
	})

	preview, err := pnm.previewGenerator.GenerateProjectStructurePreview(ctx, pnm.originalConfig, pnm.originalSelections, pnm.originalOutputDir)
	if err != nil {
		return fmt.Errorf("failed to regenerate preview: %w", err)
	}

	pnm.currentPreview = preview
	return nil
}

// GetCurrentPreview returns the current preview
func (pnm *PreviewNavigationManager) GetCurrentPreview() *ProjectStructurePreview {
	return pnm.currentPreview
}

// GetModifiedSelections returns the current template selections
func (pnm *PreviewNavigationManager) GetModifiedSelections() []TemplateSelection {
	return pnm.originalSelections
}

// GetModifiedConfig returns the current project configuration
func (pnm *PreviewNavigationManager) GetModifiedConfig() *models.ProjectConfig {
	return pnm.originalConfig
}

// GetModifiedOutputDir returns the current output directory
func (pnm *PreviewNavigationManager) GetModifiedOutputDir() string {
	return pnm.originalOutputDir
}
