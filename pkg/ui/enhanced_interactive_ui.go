// Package ui provides interactive UI with comprehensive navigation and help systems.
//
// This file implements the InteractiveUI which integrates the NavigationSystem
// and HelpSystem to provide a complete interactive experience with breadcrumbs,
// context-sensitive help, and comprehensive error handling.
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// AdvancedInteractiveUI extends the basic InteractiveUI with navigation and help systems
type AdvancedInteractiveUI struct {
	*InteractiveUI
	navigationSystem *NavigationSystem
	helpSystem       *HelpSystem
	config           *AdvancedUIConfig
}

// AdvancedUIConfig defines configuration for the advanced interactive UI
type AdvancedUIConfig struct {
	*UIConfig
	NavigationConfig  *NavigationConfig `json:"navigation_config"`
	HelpConfig        *HelpConfig       `json:"help_config"`
	EnableBreadcrumbs bool              `json:"enable_breadcrumbs"`
	EnableStepCounter bool              `json:"enable_step_counter"`
	EnableContextHelp bool              `json:"enable_context_help"`
}

// NewAdvancedInteractiveUI creates a new advanced interactive UI
func NewAdvancedInteractiveUI(logger interfaces.Logger, config *AdvancedUIConfig) interfaces.InteractiveUIInterface {
	if config == nil {
		config = &AdvancedUIConfig{
			UIConfig: &UIConfig{
				EnableColors:    true,
				EnableUnicode:   true,
				PageSize:        10,
				ShowBreadcrumbs: true,
				ShowShortcuts:   true,
			},
			EnableBreadcrumbs: true,
			EnableStepCounter: true,
			EnableContextHelp: true,
		}
	}

	// Create base interactive UI
	baseUI := NewInteractiveUI(logger, config.UIConfig).(*InteractiveUI)

	// Create advanced interactive UI
	advancedUI := &AdvancedInteractiveUI{
		InteractiveUI: baseUI,
		config:        config,
	}

	// Initialize navigation and help systems
	advancedUI.navigationSystem = NewNavigationSystem(advancedUI, logger, config.NavigationConfig)
	advancedUI.helpSystem = NewHelpSystem(advancedUI, logger, config.HelpConfig)

	return advancedUI
}

// ShowMenu displays a menu with navigation and help
func (ui *AdvancedInteractiveUI) ShowMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	// Set up navigation context
	ui.setupNavigationForMenu(config)

	// Show navigation header
	if err := ui.showNavigationHeader(ctx, config.Title, config.Description); err != nil {
		return nil, fmt.Errorf("failed to show navigation header: %w", err)
	}

	// Call base menu implementation
	result, err := ui.InteractiveUI.ShowMenu(ctx, config)
	if err != nil {
		return nil, err
	}

	// Handle navigation actions
	if result.Action == "help" {
		if err := ui.helpSystem.ShowContextHelp(ctx, "menu", config.HelpText); err != nil {
			ui.logger.WarnWithFields("Failed to show context help", map[string]interface{}{
				"error": err.Error(),
			})
		}
		// Recursively show menu again after help
		return ui.ShowMenu(ctx, config)
	}

	// Add to navigation history
	ui.navigationSystem.AddToHistory("menu_selection", map[string]interface{}{
		"selected_index": result.SelectedIndex,
		"selected_value": result.SelectedValue,
		"action":         result.Action,
	})

	return result, nil
}

// ShowMultiSelect displays a multi-select with navigation and help
func (ui *AdvancedInteractiveUI) ShowMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	// Set up navigation context
	ui.setupNavigationForMultiSelect(config)

	// Show navigation header
	if err := ui.showNavigationHeader(ctx, config.Title, config.Description); err != nil {
		return nil, fmt.Errorf("failed to show navigation header: %w", err)
	}

	// Call base multi-select implementation
	result, err := ui.InteractiveUI.ShowMultiSelect(ctx, config)
	if err != nil {
		return nil, err
	}

	// Handle navigation actions
	if result.Action == "help" {
		if err := ui.helpSystem.ShowContextHelp(ctx, "multiselect", config.HelpText); err != nil {
			ui.logger.WarnWithFields("Failed to show context help", map[string]interface{}{
				"error": err.Error(),
			})
		}
		// Recursively show multi-select again after help
		return ui.ShowMultiSelect(ctx, config)
	}

	// Add to navigation history
	ui.navigationSystem.AddToHistory("multiselect_selection", map[string]interface{}{
		"selected_indices": result.SelectedIndices,
		"selected_count":   len(result.SelectedIndices),
		"action":           result.Action,
		"search_query":     result.SearchQuery,
	})

	return result, nil
}

// PromptText displays a text prompt with navigation and help
func (ui *AdvancedInteractiveUI) PromptText(ctx context.Context, config interfaces.TextPromptConfig) (*interfaces.TextResult, error) {
	// Set up navigation context
	ui.setupNavigationForTextInput(config)

	// Show navigation header
	if err := ui.showNavigationHeader(ctx, config.Prompt, config.Description); err != nil {
		return nil, fmt.Errorf("failed to show navigation header: %w", err)
	}

	// Call base text prompt implementation
	result, err := ui.InteractiveUI.PromptText(ctx, config)
	if err != nil {
		// Handle error with help system
		if errorResult, helpErr := ui.helpSystem.HandleError(ctx, err, "validation_error"); helpErr == nil && errorResult != nil {
			switch errorResult.Action {
			case "retry":
				return ui.PromptText(ctx, config)
			case "back":
				result = &interfaces.TextResult{Action: "back", Cancelled: true}
				return result, nil
			case "quit":
				result = &interfaces.TextResult{Action: "quit", Cancelled: true}
				return result, nil
			}
		}
		return nil, err
	}

	// Handle navigation actions
	if result.Action == "help" {
		if err := ui.helpSystem.ShowContextHelp(ctx, "text_input", config.HelpText); err != nil {
			ui.logger.WarnWithFields("Failed to show context help", map[string]interface{}{
				"error": err.Error(),
			})
		}
		// Recursively show text prompt again after help
		return ui.PromptText(ctx, config)
	}

	// Add to navigation history
	ui.navigationSystem.AddToHistory("text_input", map[string]interface{}{
		"prompt": config.Prompt,
		"value":  result.Value,
		"action": result.Action,
	})

	return result, nil
}

// PromptConfirm displays a confirmation prompt with navigation and help
func (ui *AdvancedInteractiveUI) PromptConfirm(ctx context.Context, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, error) {
	// Set up navigation context
	ui.setupNavigationForConfirm(config)

	// Show navigation header
	if err := ui.showNavigationHeader(ctx, config.Prompt, config.Description); err != nil {
		return nil, fmt.Errorf("failed to show navigation header: %w", err)
	}

	// Call base confirm implementation
	result, err := ui.InteractiveUI.PromptConfirm(ctx, config)
	if err != nil {
		return nil, err
	}

	// Handle navigation actions
	if result.Action == "help" {
		if err := ui.helpSystem.ShowContextHelp(ctx, "confirm", config.HelpText); err != nil {
			ui.logger.WarnWithFields("Failed to show context help", map[string]interface{}{
				"error": err.Error(),
			})
		}
		// Recursively show confirm again after help
		return ui.PromptConfirm(ctx, config)
	}

	// Add to navigation history
	ui.navigationSystem.AddToHistory("confirm", map[string]interface{}{
		"prompt":    config.Prompt,
		"confirmed": result.Confirmed,
		"action":    result.Action,
	})

	return result, nil
}

// ShowError displays an error with recovery options
func (ui *AdvancedInteractiveUI) ShowError(ctx context.Context, config interfaces.ErrorConfig) (*interfaces.ErrorResult, error) {
	// Set up navigation context for error handling
	ui.setupNavigationForError(config)

	// Show navigation header
	if err := ui.showNavigationHeader(ctx, config.Title, "An error occurred. Please choose how to proceed."); err != nil {
		return nil, fmt.Errorf("failed to show navigation header: %w", err)
	}

	// Display error information
	fmt.Printf("%s %s\n", ui.colorize("Error:", "red"), config.Message)
	if config.Details != "" {
		fmt.Printf("%s\n", ui.colorize(config.Details, "gray"))
	}
	fmt.Println()

	// Show suggestions if available
	if len(config.Suggestions) > 0 {
		fmt.Printf("%s\n", ui.colorize("Suggestions:", "yellow"))
		for _, suggestion := range config.Suggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
		fmt.Println()
	}

	// Show recovery options if available
	if len(config.RecoveryOptions) > 0 {
		fmt.Printf("%s\n", ui.colorize("Recovery Options:", "cyan"))
		for i, option := range config.RecoveryOptions {
			safetyIndicator := ""
			if option.Safe {
				safetyIndicator = ui.colorize(" (Safe)", "green")
			} else {
				safetyIndicator = ui.colorize(" (Caution)", "yellow")
			}
			fmt.Printf("  %d. %s%s\n", i+1, option.Label, safetyIndicator)
			fmt.Printf("     %s\n", option.Description)
		}
		fmt.Println()
	}

	// Show available actions
	actions := []string{}
	if config.AllowRetry {
		actions = append(actions, "r: Retry")
	}
	if config.AllowIgnore {
		actions = append(actions, "i: Ignore")
	}
	if len(config.RecoveryOptions) > 0 {
		actions = append(actions, "1-9: Select recovery option")
	}
	if config.AllowBack {
		actions = append(actions, "b: Back")
	}
	if config.AllowQuit {
		actions = append(actions, "q: Quit")
	}

	if len(actions) > 0 {
		fmt.Printf("%s\n", ui.colorize(fmt.Sprintf("Actions: %s", fmt.Sprintf("%v", actions)), "gray"))
	}

	// Get user input
	fmt.Print("Choose an action: ")
	input, err := ui.readInput()
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Process input
	result := &interfaces.ErrorResult{}
	switch input {
	case "r", "retry":
		if config.AllowRetry {
			result.Action = "retry"
		}
	case "i", "ignore":
		if config.AllowIgnore {
			result.Action = "ignore"
		}
	case "b", "back":
		if config.AllowBack {
			result.Action = "back"
			result.Cancelled = true
		}
	case "q", "quit":
		if config.AllowQuit {
			result.Action = "quit"
			result.Cancelled = true
		}
	default:
		// Check for recovery option selection
		if len(config.RecoveryOptions) > 0 {
			if optionIndex := ui.parseRecoveryOptionIndex(input, len(config.RecoveryOptions)); optionIndex >= 0 {
				result.Action = "recovery"
				result.RecoverySelected = optionIndex
			}
		}
	}

	// Add to navigation history
	ui.navigationSystem.AddToHistory("error_handling", map[string]interface{}{
		"error_type": config.ErrorType,
		"action":     result.Action,
		"recovery":   result.RecoverySelected,
	})

	return result, nil
}

// ShowHelp displays context-sensitive help
func (ui *AdvancedInteractiveUI) ShowHelp(ctx context.Context, helpContext string) error {
	return ui.helpSystem.ShowContextHelp(ctx, helpContext)
}

// ShowBreadcrumb displays breadcrumb navigation
func (ui *AdvancedInteractiveUI) ShowBreadcrumb(ctx context.Context, path []string) error {
	if !ui.config.EnableBreadcrumbs || len(path) == 0 {
		return nil
	}

	breadcrumbText := strings.Join(path, " > ")
	fmt.Printf("%s %s\n", ui.colorize("📍", "blue"), ui.colorize(breadcrumbText, "blue"))
	return nil
}

// Helper methods for setting up navigation contexts

func (ui *AdvancedInteractiveUI) setupNavigationForMenu(config interfaces.MenuConfig) {
	actions := []NavigationAction{
		{Key: "h", Label: "Help", Description: "Show help", Available: config.ShowHelp, Global: true},
		{Key: "b", Label: "Back", Description: "Go back", Available: config.AllowBack, Global: false},
		{Key: "q", Label: "Quit", Description: "Quit", Available: config.AllowQuit, Global: false},
	}
	ui.navigationSystem.SetAvailableActions(actions)
}

func (ui *AdvancedInteractiveUI) setupNavigationForMultiSelect(config interfaces.MultiSelectConfig) {
	actions := []NavigationAction{
		{Key: "h", Label: "Help", Description: "Show help", Available: config.ShowHelp, Global: true},
		{Key: "b", Label: "Back", Description: "Go back", Available: config.AllowBack, Global: false},
		{Key: "q", Label: "Quit", Description: "Quit", Available: config.AllowQuit, Global: false},
	}
	ui.navigationSystem.SetAvailableActions(actions)
}

func (ui *AdvancedInteractiveUI) setupNavigationForTextInput(config interfaces.TextPromptConfig) {
	actions := []NavigationAction{
		{Key: "h", Label: "Help", Description: "Show help", Available: config.ShowHelp, Global: true},
		{Key: "b", Label: "Back", Description: "Go back", Available: config.AllowBack, Global: false},
		{Key: "q", Label: "Quit", Description: "Quit", Available: config.AllowQuit, Global: false},
	}
	ui.navigationSystem.SetAvailableActions(actions)
}

func (ui *AdvancedInteractiveUI) setupNavigationForConfirm(config interfaces.ConfirmConfig) {
	actions := []NavigationAction{
		{Key: "h", Label: "Help", Description: "Show help", Available: config.ShowHelp, Global: true},
		{Key: "b", Label: "Back", Description: "Go back", Available: config.AllowBack, Global: false},
		{Key: "q", Label: "Quit", Description: "Quit", Available: config.AllowQuit, Global: false},
	}
	ui.navigationSystem.SetAvailableActions(actions)
}

func (ui *AdvancedInteractiveUI) setupNavigationForError(config interfaces.ErrorConfig) {
	actions := []NavigationAction{
		{Key: "r", Label: "Retry", Description: "Retry operation", Available: config.AllowRetry, Global: false},
		{Key: "i", Label: "Ignore", Description: "Ignore error", Available: config.AllowIgnore, Global: false},
		{Key: "b", Label: "Back", Description: "Go back", Available: config.AllowBack, Global: false},
		{Key: "q", Label: "Quit", Description: "Quit", Available: config.AllowQuit, Global: false},
	}
	ui.navigationSystem.SetAvailableActions(actions)
}

func (ui *AdvancedInteractiveUI) showNavigationHeader(ctx context.Context, title, description string) error {
	// Show breadcrumbs
	if err := ui.navigationSystem.ShowBreadcrumbs(ctx); err != nil {
		return err
	}

	// Show navigation header
	if err := ui.navigationSystem.ShowNavigationHeader(ctx, title, description); err != nil {
		return err
	}

	return nil
}

func (ui *AdvancedInteractiveUI) parseRecoveryOptionIndex(input string, maxOptions int) int {
	// Try to parse as number
	if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
		index := int(input[0] - '1')
		if index < maxOptions {
			return index
		}
	}
	return -1
}

// Navigation system access methods

// SetNavigationStep sets the current navigation step
func (ui *AdvancedInteractiveUI) SetNavigationStep(stepName string, stepIndex, totalSteps int) {
	ui.navigationSystem.SetCurrentStep(stepName, stepIndex, totalSteps)
}

// GetNavigationSystem returns the navigation system
func (ui *AdvancedInteractiveUI) GetNavigationSystem() *NavigationSystem {
	return ui.navigationSystem
}

// GetHelpSystem returns the help system
func (ui *AdvancedInteractiveUI) GetHelpSystem() *HelpSystem {
	return ui.helpSystem
}

// ShowCompletionSummary displays a completion summary
func (ui *AdvancedInteractiveUI) ShowCompletionSummary(ctx context.Context, summary *CompletionSummary) error {
	return ui.helpSystem.ShowCompletionSummary(ctx, summary)
}
