// Package ui provides error handling and recovery components for the interactive UI.
//
// This file contains implementations for error display, recovery options,
// and user-friendly error handling with actionable suggestions.
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ShowError displays an error with recovery options and handles user response
func (ui *InteractiveUI) ShowError(ctx context.Context, config interfaces.ErrorConfig) (*interfaces.ErrorResult, error) {
	ui.clearScreen()

	// Display error header
	title := config.Title
	if title == "" {
		title = "Error"
	}
	ui.showHeader(title, "")

	// Display error message
	fmt.Printf("%s %s\n", ui.colorize("✗", "red"), config.Message)

	// Display error details if provided
	if config.Details != "" {
		fmt.Printf("\n%s\n", ui.colorize("Details:", "yellow"))
		fmt.Printf("%s\n", config.Details)
	}

	// Display error type if provided
	if config.ErrorType != "" {
		fmt.Printf("\n%s %s\n", ui.colorize("Type:", "gray"), config.ErrorType)
	}

	// Display suggestions if provided
	if len(config.Suggestions) > 0 {
		fmt.Printf("\n%s\n", ui.colorize("Suggestions:", "cyan"))
		for i, suggestion := range config.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	// Display recovery options if provided
	if len(config.RecoveryOptions) > 0 {
		fmt.Printf("\n%s\n", ui.colorize("Recovery Options:", "green"))
		for i, option := range config.RecoveryOptions {
			safetyIndicator := ""
			if option.Safe {
				safetyIndicator = ui.colorize(" (safe)", "green")
			} else {
				safetyIndicator = ui.colorize(" (caution)", "yellow")
			}
			fmt.Printf("  %d. %s%s\n", i+1, option.Label, safetyIndicator)
			if option.Description != "" {
				fmt.Printf("     %s\n", ui.colorize(option.Description, "gray"))
			}
		}
	}

	fmt.Println()

	// Handle user input for error response
	for {
		select {
		case <-ctx.Done():
			return &interfaces.ErrorResult{Cancelled: true}, ctx.Err()
		default:
		}

		ui.showErrorHelp(config)

		input, err := ui.readInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := ui.handleErrorInput(input, config)
		if shouldReturn {
			return result, nil
		}
	}
}

// showErrorHelp displays available error handling options
func (ui *InteractiveUI) showErrorHelp(config interfaces.ErrorConfig) {
	if !ui.config.ShowShortcuts {
		return
	}

	help := []string{}

	if config.AllowRetry {
		help = append(help, "r: Retry")
	}
	if config.AllowIgnore {
		help = append(help, "i: Ignore")
	}
	if len(config.RecoveryOptions) > 0 {
		help = append(help, "#: Select recovery option")
	}
	if config.AllowBack {
		help = append(help, "b: Back")
	}
	if config.AllowQuit {
		help = append(help, "q: Quit")
	}

	if len(help) > 0 {
		fmt.Printf("%s\n", ui.colorize(strings.Join(help, " | "), "gray"))
	}

	fmt.Print("Choose an action: ")
}

// handleErrorInput processes user input for error handling
func (ui *InteractiveUI) handleErrorInput(input string, config interfaces.ErrorConfig) (*interfaces.ErrorResult, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "r", "retry":
		if config.AllowRetry {
			return &interfaces.ErrorResult{Action: "retry"}, true
		}

	case "i", "ignore":
		if config.AllowIgnore {
			return &interfaces.ErrorResult{Action: "ignore"}, true
		}

	case "b", "back":
		if config.AllowBack {
			return &interfaces.ErrorResult{Action: "back"}, true
		}

	case "q", "quit":
		if config.AllowQuit {
			if ui.config.ConfirmOnQuit {
				if ui.confirmQuit() {
					return &interfaces.ErrorResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.ErrorResult{Action: "quit", Cancelled: true}, true
		}

	default:
		// Check for recovery option selection
		if len(config.RecoveryOptions) > 0 {
			for i := range config.RecoveryOptions {
				optionNum := fmt.Sprintf("%d", i+1)
				if input == optionNum {
					// Execute recovery action if provided
					option := config.RecoveryOptions[i]
					if option.Action != nil {
						if err := option.Action(); err != nil {
							ui.showError(fmt.Sprintf("Recovery action failed: %s", err.Error()))
							return nil, false
						}
					}
					return &interfaces.ErrorResult{
						Action:           "recovery",
						RecoverySelected: i,
					}, true
				}
			}
		}
	}

	ui.showError("Invalid option. Please try again.")
	return nil, false
}

// ShowBreadcrumb displays navigation breadcrumbs
func (ui *InteractiveUI) ShowBreadcrumb(ctx context.Context, path []string) error {
	if !ui.config.ShowBreadcrumbs || len(path) == 0 {
		return nil
	}

	breadcrumb := strings.Join(path, ui.colorize(" > ", "gray"))
	fmt.Printf("%s %s\n\n", ui.colorize("Navigation:", "cyan"), breadcrumb)

	return nil
}

// ShowHelp displays context-sensitive help information
func (ui *InteractiveUI) ShowHelp(ctx context.Context, helpContext string) error {
	ui.clearScreen()
	ui.showHeader("Help", "")

	switch helpContext {
	case "menu":
		ui.showMenuHelp()
	case "multiselect":
		ui.showMultiSelectHelpDetailed()
	case "text":
		ui.showTextInputHelp()
	case "confirm":
		ui.showConfirmHelp()
	case "checkbox":
		ui.showCheckboxHelp()
	case "table":
		ui.showTableHelpDetailed()
	case "tree":
		ui.showTreeHelpDetailed()
	case "error":
		ui.showErrorHelpDetailed()
	default:
		ui.showGeneralHelp()
	}

	fmt.Print("\nPress Enter to continue...")
	if _, err := ui.readInput(); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("\nError reading input: %v\n", err)
	}
	return nil
}

// showMenuHelp displays detailed menu help
func (ui *InteractiveUI) showMenuHelp() {
	fmt.Println("Menu Navigation Help")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  • Use ↑/↓ arrow keys or j/k to move up and down")
	fmt.Println("  • Press Enter to select the highlighted option")
	fmt.Println("  • Type a number (1-9) to jump directly to that option")
	fmt.Println("  • Use shortcut keys shown in brackets [x] if available")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back to previous screen (if available)")
	fmt.Println("  • q: Quit the application (if available)")
	fmt.Println("  • Ctrl+C: Cancel current operation")
}

// showMultiSelectHelpDetailed displays detailed multi-select help
func (ui *InteractiveUI) showMultiSelectHelpDetailed() {
	fmt.Println("Multi-Selection Help")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  • Use ↑/↓ arrow keys to move up and down")
	fmt.Println("  • Press Space to toggle selection of current item")
	fmt.Println("  • Press Enter when finished selecting")
	fmt.Println()
	fmt.Println("Search:")
	fmt.Println("  • Type / to start searching")
	fmt.Println("  • Type search terms to filter options")
	fmt.Println("  • Type 'clear' to clear search and show all options")
	fmt.Println()
	fmt.Println("Selection:")
	fmt.Println("  • ☐ indicates unselected item")
	fmt.Println("  • ☑ indicates selected item")
	fmt.Println("  • Selection count and limits are shown at the top")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit (if available)")
}

// showTextInputHelp displays detailed text input help
func (ui *InteractiveUI) showTextInputHelp() {
	fmt.Println("Text Input Help")
	fmt.Println(strings.Repeat("=", 15))
	fmt.Println()
	fmt.Println("Input:")
	fmt.Println("  • Type your response and press Enter")
	fmt.Println("  • Leave empty to use default value (shown in brackets)")
	fmt.Println("  • Required fields are marked with * (red asterisk)")
	fmt.Println()
	fmt.Println("Validation:")
	fmt.Println("  • Input is validated before acceptance")
	fmt.Println("  • Error messages will show specific requirements")
	fmt.Println("  • Suggestions may be provided for invalid input")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit (if available)")
}

// showConfirmHelp displays detailed confirmation help
func (ui *InteractiveUI) showConfirmHelp() {
	fmt.Println("Confirmation Help")
	fmt.Println(strings.Repeat("=", 17))
	fmt.Println()
	fmt.Println("Responses:")
	fmt.Println("  • y, yes, true, 1: Confirm (yes)")
	fmt.Println("  • n, no, false, 0: Decline (no)")
	fmt.Println("  • Enter: Use default value (shown in brackets)")
	fmt.Println()
	fmt.Println("Default values:")
	fmt.Println("  • [Yes/no]: Default is Yes")
	fmt.Println("  • [yes/No]: Default is No")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit (if available)")
}

// showCheckboxHelp displays detailed checkbox help
func (ui *InteractiveUI) showCheckboxHelp() {
	fmt.Println("Checkbox List Help")
	fmt.Println(strings.Repeat("=", 18))
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  • Use ↑/↓ arrow keys or j/k to move up and down")
	fmt.Println("  • Press Space to toggle checkbox state")
	fmt.Println("  • Press Enter when finished")
	fmt.Println()
	fmt.Println("Checkboxes:")
	fmt.Println("  • ☐ indicates unchecked item")
	fmt.Println("  • ☑ indicates checked item")
	fmt.Println("  • Required items are marked with * (red asterisk)")
	fmt.Println()
	fmt.Println("Validation:")
	fmt.Println("  • All required items must be checked")
	fmt.Println("  • Error messages will indicate missing requirements")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit (if available)")
}

// showTableHelpDetailed displays detailed table help
func (ui *InteractiveUI) showTableHelpDetailed() {
	fmt.Println("Table Navigation Help")
	fmt.Println(strings.Repeat("=", 21))
	fmt.Println()
	fmt.Println("Pagination:")
	fmt.Println("  • n, next: Go to next page")
	fmt.Println("  • p, prev, previous: Go to previous page")
	fmt.Println("  • Type page number to jump to specific page")
	fmt.Println()
	fmt.Println("Search:")
	fmt.Println("  • Type / to search table contents")
	fmt.Println("  • Type 'clear' to clear search filter")
	fmt.Println("  • Search looks in all columns")
	fmt.Println()
	fmt.Println("Sorting:")
	fmt.Println("  • Type 'sort #' where # is column number (1-based)")
	fmt.Println("  • Sorting same column twice reverses order")
	fmt.Println("  • Sort indicators: ↑ (ascending) ↓ (descending)")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • q, quit, exit: Close table")
	fmt.Println("  • Enter: Close table")
}

// showTreeHelpDetailed displays detailed tree help
func (ui *InteractiveUI) showTreeHelpDetailed() {
	fmt.Println("Tree Navigation Help")
	fmt.Println(strings.Repeat("=", 20))
	fmt.Println()
	fmt.Println("Expansion:")
	fmt.Println("  • ▶ indicates collapsed node (has children)")
	fmt.Println("  • ▼ indicates expanded node (showing children)")
	fmt.Println("  • e, expand: Expand all nodes")
	fmt.Println("  • c, collapse: Collapse all nodes")
	fmt.Println("  • t, toggle: Toggle current node")
	fmt.Println()
	fmt.Println("Structure:")
	fmt.Println("  • Tree lines show parent-child relationships")
	fmt.Println("  • Icons may be shown for different node types")
	fmt.Println("  • Indentation indicates hierarchy level")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • q, quit, exit: Close tree view")
	fmt.Println("  • Enter: Close tree view")
}

// showErrorHelpDetailed displays detailed error handling help
func (ui *InteractiveUI) showErrorHelpDetailed() {
	fmt.Println("Error Handling Help")
	fmt.Println(strings.Repeat("=", 19))
	fmt.Println()
	fmt.Println("Error Information:")
	fmt.Println("  • Error message describes what went wrong")
	fmt.Println("  • Details provide additional context")
	fmt.Println("  • Error type categorizes the issue")
	fmt.Println()
	fmt.Println("Suggestions:")
	fmt.Println("  • Numbered list of potential solutions")
	fmt.Println("  • Try suggestions before using recovery options")
	fmt.Println()
	fmt.Println("Recovery Options:")
	fmt.Println("  • Numbered list of automated recovery actions")
	fmt.Println("  • (safe) indicates low-risk recovery")
	fmt.Println("  • (caution) indicates potentially risky recovery")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  • r, retry: Try the operation again")
	fmt.Println("  • i, ignore: Continue despite the error")
	fmt.Println("  • #: Execute recovery option (1-9)")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit (if available)")
}

// showGeneralHelp displays general application help
func (ui *InteractiveUI) showGeneralHelp() {
	fmt.Println("General Help")
	fmt.Println(strings.Repeat("=", 12))
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  • Use arrow keys, j/k, or numbers to navigate")
	fmt.Println("  • Press Enter to select or confirm")
	fmt.Println("  • Press Space to toggle selections")
	fmt.Println()
	fmt.Println("Global Commands:")
	fmt.Println("  • h: Show context-sensitive help")
	fmt.Println("  • b: Go back to previous screen")
	fmt.Println("  • q: Quit application")
	fmt.Println("  • Ctrl+C: Cancel current operation")
	fmt.Println()
	fmt.Println("Visual Indicators:")
	fmt.Println("  • > indicates current selection")
	fmt.Println("  • * indicates required fields")
	fmt.Println("  • ☐/☑ indicate checkbox states")
	fmt.Println("  • ▶/▼ indicate expandable items")
	fmt.Println("  • Colors indicate status (green=good, red=error, etc.)")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("  • Default values are shown in [brackets]")
	fmt.Println("  • Shortcut keys are shown in [brackets] after options")
	fmt.Println("  • Help is available in most screens with 'h'")
	fmt.Println("  • Use Tab completion where available")
}

// Additional error handling utilities

// CreateValidationError creates a validation error with suggestions and recovery options
func CreateValidationError(field, value, message string, suggestions []string, recoveryOptions []interfaces.RecoveryOption) *interfaces.ValidationError {
	err := interfaces.NewValidationError(field, value, message, "validation_failed")
	if suggestErr := err.WithSuggestions(suggestions...); suggestErr != nil {
		// Log error but continue
		fmt.Printf("Warning: Failed to add suggestions to validation error: %v\n", suggestErr)
	}
	if recoveryErr := err.WithRecoveryOptions(recoveryOptions...); recoveryErr != nil {
		// Log error but continue
		fmt.Printf("Warning: Failed to add recovery options to validation error: %v\n", recoveryErr)
	}
	return err
}

// CreateRecoveryOption creates a recovery option with the specified parameters
func CreateRecoveryOption(label, description string, action func() error, safe bool) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       label,
		Description: description,
		Action:      action,
		Safe:        safe,
	}
}

// HandleValidationError handles a validation error with user-friendly display
func (ui *InteractiveUI) HandleValidationError(ctx context.Context, err *interfaces.ValidationError) (*interfaces.ErrorResult, error) {
	config := interfaces.ErrorConfig{
		Title:           "Validation Error",
		Message:         err.Message,
		Details:         fmt.Sprintf("Field: %s, Value: %s", err.Field, err.Value),
		ErrorType:       "Validation",
		Suggestions:     err.Suggestions,
		RecoveryOptions: err.RecoveryOptions,
		AllowRetry:      true,
		AllowBack:       true,
		AllowQuit:       true,
	}

	return ui.ShowError(ctx, config)
}

// HandleGenericError handles a generic error with basic recovery options
func (ui *InteractiveUI) HandleGenericError(ctx context.Context, err error, allowRetry, allowIgnore bool) (*interfaces.ErrorResult, error) {
	config := interfaces.ErrorConfig{
		Title:       "Error",
		Message:     err.Error(),
		ErrorType:   "General",
		AllowRetry:  allowRetry,
		AllowIgnore: allowIgnore,
		AllowBack:   true,
		AllowQuit:   true,
	}

	return ui.ShowError(ctx, config)
}
