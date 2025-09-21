// Package ui provides interactive user interface components for the CLI generator.
//
// This package implements comprehensive interactive functionality including:
//   - Menu navigation with keyboard shortcuts
//   - Input collection with validation and error recovery
//   - Progress tracking and visual feedback
//   - Context-sensitive help system
//   - Session management and state persistence
package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// InteractiveUI implements the InteractiveUIInterface for terminal-based interaction
type InteractiveUI struct {
	reader         *bufio.Reader
	writer         *bufio.Writer
	shortcuts      map[string]interfaces.KeyboardShortcut
	currentSession *interfaces.UISession
	logger         interfaces.Logger
	config         *UIConfig
}

// UIConfig defines configuration for the interactive UI
type UIConfig struct {
	EnableColors    bool          `json:"enable_colors"`
	EnableUnicode   bool          `json:"enable_unicode"`
	PageSize        int           `json:"page_size"`
	Timeout         time.Duration `json:"timeout"`
	AutoSave        bool          `json:"auto_save"`
	ShowBreadcrumbs bool          `json:"show_breadcrumbs"`
	ShowShortcuts   bool          `json:"show_shortcuts"`
	ConfirmOnQuit   bool          `json:"confirm_on_quit"`
}

// NewInteractiveUI creates a new interactive UI instance
func NewInteractiveUI(logger interfaces.Logger, config *UIConfig) interfaces.InteractiveUIInterface {
	if config == nil {
		config = &UIConfig{
			EnableColors:    true,
			EnableUnicode:   true,
			PageSize:        10,
			Timeout:         30 * time.Minute,
			AutoSave:        true,
			ShowBreadcrumbs: true,
			ShowShortcuts:   true,
			ConfirmOnQuit:   true,
		}
	}

	ui := &InteractiveUI{
		reader:    bufio.NewReader(os.Stdin),
		writer:    bufio.NewWriter(os.Stdout),
		shortcuts: make(map[string]interfaces.KeyboardShortcut),
		logger:    logger,
		config:    config,
	}

	ui.setupDefaultShortcuts()
	return ui
}

// setupDefaultShortcuts initializes default keyboard shortcuts
func (ui *InteractiveUI) setupDefaultShortcuts() {
	shortcuts := []interfaces.KeyboardShortcut{
		{Key: "q", Description: "Quit", Action: "quit", Global: true},
		{Key: "h", Description: "Help", Action: "help", Global: true},
		{Key: "b", Description: "Back", Action: "back", Global: true},
		{Key: "ctrl+c", Description: "Cancel", Action: "cancel", Global: true},
		{Key: "enter", Description: "Select/Confirm", Action: "confirm", Global: true},
		{Key: "esc", Description: "Cancel/Back", Action: "back", Global: true},
		{Key: "tab", Description: "Next", Action: "next", Global: false},
		{Key: "shift+tab", Description: "Previous", Action: "previous", Global: false},
	}

	for _, shortcut := range shortcuts {
		ui.shortcuts[shortcut.Key] = shortcut
	}
}

// ShowMenu displays an interactive menu and handles user selection
func (ui *InteractiveUI) ShowMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("menu must have at least one option")
	}

	// Always use scrollable menu for proper arrow key handling
	scrollableMenu := NewScrollableMenu(ui)
	return scrollableMenu.ShowScrollableMenu(ctx, config)
}



// ShowMultiSelect displays a multi-selection interface
func (ui *InteractiveUI) ShowMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("multi-select must have at least one option")
	}

	// Always use scrollable multi-select for proper arrow key handling
	scrollableMenu := NewScrollableMenu(ui)
	return scrollableMenu.ShowScrollableMultiSelect(ctx, config)
}



// PromptText displays a text input prompt with validation
func (ui *InteractiveUI) PromptText(ctx context.Context, config interfaces.TextPromptConfig) (*interfaces.TextResult, error) {
	ui.clearScreen()
	ui.showPromptHeader(config.Prompt, config.Description)

	for {
		select {
		case <-ctx.Done():
			return &interfaces.TextResult{Cancelled: true}, ctx.Err()
		default:
		}

		ui.displayTextPrompt(config)

		input, err := ui.readInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := ui.handleTextInput(input, config)
		if shouldReturn {
			return result, nil
		}
	}
}

// displayTextPrompt renders the text input prompt
func (ui *InteractiveUI) displayTextPrompt(config interfaces.TextPromptConfig) {
	prompt := config.Prompt
	if config.Required {
		prompt += ui.colorize(" *", "red")
	}

	if config.DefaultValue != "" {
		prompt += ui.colorize(fmt.Sprintf(" [%s]", config.DefaultValue), "gray")
	}

	if config.Placeholder != "" {
		prompt += ui.colorize(fmt.Sprintf(" (%s)", config.Placeholder), "gray")
	}

	fmt.Printf("%s: ", prompt)
}

// handleTextInput processes user input for text prompts
func (ui *InteractiveUI) handleTextInput(input string, config interfaces.TextPromptConfig) (*interfaces.TextResult, bool) {
	input = strings.TrimSpace(input)
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "q", "quit":
		if config.AllowQuit {
			if ui.config.ConfirmOnQuit {
				if ui.confirmQuit() {
					return &interfaces.TextResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.TextResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			ui.showContextHelp(config.HelpText, "text")
			return nil, false
		}

	case "b", "back":
		if config.AllowBack {
			return &interfaces.TextResult{Action: "back", Cancelled: true}, true
		}
	}

	// Use default value if input is empty
	if input == "" && config.DefaultValue != "" {
		input = config.DefaultValue
	}

	// Validate required field
	if config.Required && input == "" {
		ui.showError("This field is required")
		return nil, false
	}

	// Validate length constraints
	if config.MinLength > 0 && len(input) < config.MinLength {
		ui.showError(fmt.Sprintf("Input must be at least %d characters", config.MinLength))
		return nil, false
	}

	if config.MaxLength > 0 && len(input) > config.MaxLength {
		ui.showError(fmt.Sprintf("Input must be at most %d characters", config.MaxLength))
		return nil, false
	}

	// Run custom validator
	if config.Validator != nil {
		if err := config.Validator(input); err != nil {
			if validationErr, ok := err.(*interfaces.ValidationError); ok {
				ui.showValidationError(validationErr)
			} else {
				ui.showError(err.Error())
			}
			return nil, false
		}
	}

	return &interfaces.TextResult{
		Value:  input,
		Action: "submit",
	}, true
}

// Utility methods for UI rendering and interaction

// clearScreen clears the terminal screen
func (ui *InteractiveUI) clearScreen() {
	if ui.config.EnableColors {
		fmt.Print("\033[2J\033[H")
	} else {
		fmt.Print("\n" + strings.Repeat("=", 80) + "\n")
	}
}

// showHeader displays a formatted header
func (ui *InteractiveUI) showHeader(title, description string) {
	if title != "" {
		fmt.Println(ui.colorize(title, "bold"))
		if description != "" {
			fmt.Println(ui.colorize(description, "gray"))
		}
		fmt.Println()
	}
}

// showPromptHeader displays a formatted prompt header
func (ui *InteractiveUI) showPromptHeader(prompt, description string) {
	if prompt != "" {
		fmt.Println(ui.colorize(prompt, "bold"))
		if description != "" {
			fmt.Println(ui.colorize(description, "gray"))
		}
		fmt.Println()
	}
}

// colorize applies color formatting to text
func (ui *InteractiveUI) colorize(text, color string) string {
	if !ui.config.EnableColors {
		return text
	}

	colors := map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"blue":   "\033[34m",
		"purple": "\033[35m",
		"cyan":   "\033[36m",
		"white":  "\033[37m",
		"gray":   "\033[90m",
		"bold":   "\033[1m",
		"reset":  "\033[0m",
	}

	if colorCode, exists := colors[color]; exists {
		return colorCode + text + colors["reset"]
	}
	return text
}

// readInput reads user input from stdin
func (ui *InteractiveUI) readInput() (string, error) {
	input, err := ui.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}

// showError displays an error message
func (ui *InteractiveUI) showError(message string) {
	fmt.Printf("\n%s %s\n\n", ui.colorize("Error:", "red"), message)
	fmt.Print("Press Enter to continue...")
	if _, err := ui.readInput(); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("\nError reading input: %v\n", err)
	}
}

// showValidationError displays a validation error with recovery options
func (ui *InteractiveUI) showValidationError(err *interfaces.ValidationError) {
	fmt.Printf("\n%s %s\n", ui.colorize("Validation Error:", "red"), err.Message)

	if len(err.Suggestions) > 0 {
		fmt.Println(ui.colorize("Suggestions:", "yellow"))
		for _, suggestion := range err.Suggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
	}

	fmt.Print("\nPress Enter to continue...")
	if _, readErr := ui.readInput(); readErr != nil {
		// Log error but don't fail the operation
		fmt.Printf("\nError reading input: %v\n", readErr)
	}
}

// confirmQuit asks for confirmation before quitting
func (ui *InteractiveUI) confirmQuit() bool {
	fmt.Printf("\n%s ", ui.colorize("Are you sure you want to quit? (y/N):", "yellow"))
	input, _ := ui.readInput()
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// Additional methods will be implemented in the next part...
// showNavigationHelp displays navigation help for menus
func (ui *InteractiveUI) showNavigationHelp(allowBack, allowQuit, showHelp bool) {
	if !ui.config.ShowShortcuts {
		return
	}

	help := []string{}
	help = append(help, "↑/↓ or j/k: Navigate")
	help = append(help, "Enter: Select")

	if allowBack {
		help = append(help, "b: Back")
	}
	if allowQuit {
		help = append(help, "q: Quit")
	}
	if showHelp {
		help = append(help, "h: Help")
	}

	fmt.Printf("%s\n", ui.colorize(strings.Join(help, " | "), "gray"))
}


// showContextHelp displays context-sensitive help
func (ui *InteractiveUI) showContextHelp(helpText, context string) {
	ui.clearScreen()
	fmt.Println(ui.colorize("Help", "bold"))
	fmt.Println()

	if helpText != "" {
		fmt.Println(helpText)
		fmt.Println()
	}

	// Show context-specific help
	switch context {
	case "menu":
		fmt.Println("Menu Navigation:")
		fmt.Println("  • Use arrow keys (↑/↓) or j/k to navigate")
		fmt.Println("  • Press Enter to select an option")
		fmt.Println("  • Type a number to jump to that option")
		fmt.Println("  • Use shortcut keys if available")

	case "multiselect":
		fmt.Println("Multi-Selection:")
		fmt.Println("  • Use arrow keys (↑/↓) to navigate")
		fmt.Println("  • Press Space to toggle selection")
		fmt.Println("  • Press Enter when done selecting")
		fmt.Println("  • Type / to search options")
		fmt.Println("  • Type 'clear' to clear search")

	case "text":
		fmt.Println("Text Input:")
		fmt.Println("  • Type your input and press Enter")
		fmt.Println("  • Leave empty to use default value (if available)")
		fmt.Println("  • Required fields are marked with *")
	}

	fmt.Println()
	fmt.Println("Global Commands:")
	fmt.Println("  • h: Show this help")
	fmt.Println("  • b: Go back (if available)")
	fmt.Println("  • q: Quit application (if available)")
	fmt.Println("  • Ctrl+C: Cancel current operation")

	fmt.Print("\nPress Enter to continue...")
	if _, err := ui.readInput(); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("\nError reading input: %v\n", err)
	}
}


// filterOptions filters options based on search query
func (ui *InteractiveUI) filterOptions(options []interfaces.SelectOption, query string) []int {
	if query == "" {
		indices := make([]int, len(options))
		for i := range indices {
			indices[i] = i
		}
		return indices
	}

	query = strings.ToLower(query)
	var filtered []int

	for i, option := range options {
		matched := false

		// Check label and description
		if strings.Contains(strings.ToLower(option.Label), query) ||
			strings.Contains(strings.ToLower(option.Description), query) {
			matched = true
		}

		// Also search in tags if not already matched
		if !matched {
			for _, tag := range option.Tags {
				if strings.Contains(strings.ToLower(tag), query) {
					matched = true
					break
				}
			}
		}

		if matched {
			filtered = append(filtered, i)
		}
	}

	return filtered
}

// PromptConfirm displays a confirmation prompt
func (ui *InteractiveUI) PromptConfirm(ctx context.Context, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, error) {
	ui.clearScreen()
	ui.showPromptHeader(config.Prompt, config.Description)

	for {
		select {
		case <-ctx.Done():
			return &interfaces.ConfirmResult{Cancelled: true}, ctx.Err()
		default:
		}

		ui.displayConfirmPrompt(config)

		input, err := ui.readInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := ui.handleConfirmInput(input, config)
		if shouldReturn {
			return result, nil
		}
	}
}

// displayConfirmPrompt renders the confirmation prompt
func (ui *InteractiveUI) displayConfirmPrompt(config interfaces.ConfirmConfig) {
	yesLabel := config.YesLabel
	if yesLabel == "" {
		yesLabel = "Yes"
	}

	noLabel := config.NoLabel
	if noLabel == "" {
		noLabel = "No"
	}

	defaultIndicator := ""
	if config.DefaultValue {
		defaultIndicator = fmt.Sprintf(" [%s/%s]", ui.colorize(yesLabel, "green"), strings.ToLower(noLabel))
	} else {
		defaultIndicator = fmt.Sprintf(" [%s/%s]", strings.ToLower(yesLabel), ui.colorize(noLabel, "red"))
	}

	fmt.Printf("%s%s: ", config.Prompt, defaultIndicator)
}

// handleConfirmInput processes user input for confirmation prompts
func (ui *InteractiveUI) handleConfirmInput(input string, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "q", "quit":
		if config.AllowQuit {
			if ui.config.ConfirmOnQuit {
				if ui.confirmQuit() {
					return &interfaces.ConfirmResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.ConfirmResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			ui.showContextHelp(config.HelpText, "confirm")
			return nil, false
		}

	case "b", "back":
		if config.AllowBack {
			return &interfaces.ConfirmResult{Action: "back", Cancelled: true}, true
		}

	case "y", "yes", "true", "1":
		return &interfaces.ConfirmResult{
			Confirmed: true,
			Action:    "confirm",
		}, true

	case "n", "no", "false", "0":
		return &interfaces.ConfirmResult{
			Confirmed: false,
			Action:    "confirm",
		}, true

	case "":
		// Use default value
		return &interfaces.ConfirmResult{
			Confirmed: config.DefaultValue,
			Action:    "confirm",
		}, true
	}

	ui.showError("Please enter 'y' for yes or 'n' for no")
	return nil, false
}

// ShowCheckboxList displays a checkbox list interface
func (ui *InteractiveUI) ShowCheckboxList(ctx context.Context, config interfaces.CheckboxConfig) (*interfaces.CheckboxResult, error) {
	if len(config.Items) == 0 {
		return nil, fmt.Errorf("checkbox list must have at least one item")
	}

	ui.clearScreen()
	ui.showHeader(config.Title, config.Description)

	currentIndex := 0

	for {
		select {
		case <-ctx.Done():
			return &interfaces.CheckboxResult{Cancelled: true}, ctx.Err()
		default:
		}

		ui.displayCheckboxList(config, currentIndex)
		ui.showNavigationHelp(config.AllowBack, config.AllowQuit, config.ShowHelp)

		input, err := ui.readInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := ui.handleCheckboxInput(input, &config, &currentIndex)
		if shouldReturn {
			return result, nil
		}
	}
}

// displayCheckboxList renders the checkbox list
func (ui *InteractiveUI) displayCheckboxList(config interfaces.CheckboxConfig, currentIndex int) {
	ui.clearScreen()
	ui.showHeader(config.Title, config.Description)

	fmt.Println()
	for i, item := range config.Items {
		prefix := "  "
		if i == currentIndex {
			prefix = ui.colorize("> ", "cyan")
		}

		checkbox := "☐"
		if item.Checked {
			checkbox = ui.colorize("☑", "green")
		}

		label := item.Label
		if item.Required {
			label += ui.colorize(" *", "red")
		}
		if item.Disabled {
			label = ui.colorize(label, "gray")
		} else if i == currentIndex {
			label = ui.colorize(label, "white")
		}

		fmt.Printf("%s%s %s\n", prefix, checkbox, label)

		if item.Description != "" && i == currentIndex {
			fmt.Printf("    %s\n", ui.colorize(item.Description, "gray"))
		}
	}
	fmt.Println()
}

// handleCheckboxInput processes user input for checkbox lists
func (ui *InteractiveUI) handleCheckboxInput(input string, config *interfaces.CheckboxConfig, currentIndex *int) (*interfaces.CheckboxResult, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "q", "quit":
		if config.AllowQuit {
			if ui.config.ConfirmOnQuit {
				if ui.confirmQuit() {
					return &interfaces.CheckboxResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.CheckboxResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			ui.showContextHelp(config.HelpText, "checkbox")
			return nil, false
		}

	case "b", "back":
		if config.AllowBack {
			return &interfaces.CheckboxResult{Action: "back", Cancelled: true}, true
		}

	case "", "enter", "done":
		// Validate required items
		for _, item := range config.Items {
			if item.Required && !item.Checked {
				ui.showError(fmt.Sprintf("'%s' is required", item.Label))
				return nil, false
			}
		}

		checkedIndices := []int{}
		checkedValues := []interface{}{}

		for i, item := range config.Items {
			if item.Checked {
				checkedIndices = append(checkedIndices, i)
				checkedValues = append(checkedValues, item.Value)
			}
		}

		return &interfaces.CheckboxResult{
			CheckedIndices: checkedIndices,
			CheckedValues:  checkedValues,
			Action:         "confirm",
		}, true

	case "up", "k":
		if *currentIndex > 0 {
			*currentIndex--
		} else {
			*currentIndex = len(config.Items) - 1
		}

	case "down", "j":
		if *currentIndex < len(config.Items)-1 {
			*currentIndex++
		} else {
			*currentIndex = 0
		}

	case " ", "space":
		if *currentIndex >= 0 && *currentIndex < len(config.Items) {
			item := &config.Items[*currentIndex]
			if !item.Disabled {
				item.Checked = !item.Checked
			}
		}
	}

	return nil, false
}

// Remaining interface methods (ShowTable, ShowTree, etc.) will be implemented
// in separate files to keep this file manageable

// StartSession starts a new UI session
func (ui *InteractiveUI) StartSession(ctx context.Context, config interfaces.SessionConfig) (*interfaces.UISession, error) {
	sessionID := config.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session_%d", time.Now().Unix())
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	if config.Timeout > 0 {
		sessionCtx, cancel = context.WithTimeout(ctx, config.Timeout)
	}

	session := &interfaces.UISession{
		ID:         sessionID,
		Title:      config.Title,
		StartTime:  time.Now(),
		LastActive: time.Now(),
		State:      make(map[string]interface{}),
		History:    []interfaces.SessionAction{},
		Context:    sessionCtx,
		CancelFunc: cancel,
	}

	// Copy metadata
	for k, v := range config.Metadata {
		session.State[k] = v
	}

	ui.currentSession = session

	if ui.logger != nil {
		ui.logger.InfoWithFields("Started UI session", map[string]interface{}{
			"session_id": sessionID,
			"title":      config.Title,
		})
	}

	return session, nil
}

// EndSession ends the current UI session
func (ui *InteractiveUI) EndSession(ctx context.Context, session *interfaces.UISession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	if session.CancelFunc != nil {
		session.CancelFunc()
	}

	if ui.config.AutoSave {
		if err := ui.SaveSessionState(ctx, session); err != nil {
			if ui.logger != nil {
				ui.logger.WarnWithFields("Failed to auto-save session state", map[string]interface{}{
					"session_id": session.ID,
					"error":      err.Error(),
				})
			}
		}
	}

	if ui.currentSession != nil && ui.currentSession.ID == session.ID {
		ui.currentSession = nil
	}

	if ui.logger != nil {
		duration := time.Since(session.StartTime)
		ui.logger.InfoWithFields("Ended UI session", map[string]interface{}{
			"session_id": session.ID,
			"duration":   duration.String(),
			"actions":    len(session.History),
		})
	}

	return nil
}

// SaveSessionState saves the current session state
func (ui *InteractiveUI) SaveSessionState(ctx context.Context, session *interfaces.UISession) error {
	// This would typically save to a file or database
	// For now, we'll just log the action
	if ui.logger != nil {
		ui.logger.DebugWithFields("Saving session state", map[string]interface{}{
			"session_id": session.ID,
			"state_keys": len(session.State),
		})
	}
	return nil
}

// RestoreSessionState restores a session from saved state
func (ui *InteractiveUI) RestoreSessionState(ctx context.Context, sessionID string) (*interfaces.UISession, error) {
	// This would typically load from a file or database
	// For now, we'll return an error indicating not implemented
	return nil, fmt.Errorf("session restoration not implemented")
}

// Additional utility methods for the remaining interface methods will be
// implemented in separate files to maintain code organization

// PromptSelect displays a single selection prompt
func (ui *InteractiveUI) PromptSelect(ctx context.Context, config interfaces.SelectConfig) (*interfaces.SelectResult, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("select prompt must have at least one option")
	}

	ui.clearScreen()
	ui.showPromptHeader(config.Prompt, config.Description)

	currentIndex := config.DefaultItem
	if currentIndex < 0 || currentIndex >= len(config.Options) {
		currentIndex = 0
	}

	for {
		select {
		case <-ctx.Done():
			return &interfaces.SelectResult{Cancelled: true}, ctx.Err()
		default:
		}

		ui.displaySelectPrompt(config, currentIndex)

		input, err := ui.readInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := ui.handleSelectInput(input, config, &currentIndex)
		if shouldReturn {
			return result, nil
		}
	}
}

// displaySelectPrompt renders the selection prompt
func (ui *InteractiveUI) displaySelectPrompt(config interfaces.SelectConfig, currentIndex int) {
	ui.clearScreen()
	ui.showPromptHeader(config.Prompt, config.Description)

	fmt.Println()
	for i, option := range config.Options {
		prefix := "  "
		if i == currentIndex {
			prefix = ui.colorize("> ", "cyan")
		}

		label := option
		if i == currentIndex {
			label = ui.colorize(label, "white")
		}

		fmt.Printf("%s%d. %s\n", prefix, i+1, label)
	}
	fmt.Println()

	ui.showSelectHelp(config.AllowBack, config.AllowQuit, config.ShowHelp)
}

// handleSelectInput processes user input for selection prompts
func (ui *InteractiveUI) handleSelectInput(input string, config interfaces.SelectConfig, currentIndex *int) (*interfaces.SelectResult, bool) {
	input = strings.TrimSpace(input)
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "q", "quit":
		if config.AllowQuit {
			if ui.config.ConfirmOnQuit {
				if ui.confirmQuit() {
					return &interfaces.SelectResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.SelectResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			ui.showContextHelp(config.HelpText, "select")
			return nil, false
		}

	case "b", "back":
		if config.AllowBack {
			return &interfaces.SelectResult{Action: "back", Cancelled: true}, true
		}

	case "", "enter":
		if *currentIndex >= 0 && *currentIndex < len(config.Options) {
			return &interfaces.SelectResult{
				SelectedIndex: *currentIndex,
				SelectedValue: config.Options[*currentIndex],
				Action:        "select",
			}, true
		}

	case "up", "k":
		if *currentIndex > 0 {
			*currentIndex--
		} else {
			*currentIndex = len(config.Options) - 1
		}

	case "down", "j":
		if *currentIndex < len(config.Options)-1 {
			*currentIndex++
		} else {
			*currentIndex = 0
		}

	default:
		// Check for numeric input
		if num, err := strconv.Atoi(input); err == nil {
			if num > 0 && num <= len(config.Options) {
				*currentIndex = num - 1
				return &interfaces.SelectResult{
					SelectedIndex: *currentIndex,
					SelectedValue: config.Options[*currentIndex],
					Action:        "select",
				}, true
			}
		}
	}

	return nil, false
}

// showSelectHelp displays help for selection prompts
func (ui *InteractiveUI) showSelectHelp(allowBack, allowQuit, showHelp bool) {
	if !ui.config.ShowShortcuts {
		return
	}

	help := []string{}
	help = append(help, "↑/↓ or j/k: Navigate")
	help = append(help, "Enter or #: Select")

	if allowBack {
		help = append(help, "b: Back")
	}
	if allowQuit {
		help = append(help, "q: Quit")
	}
	if showHelp {
		help = append(help, "h: Help")
	}

	fmt.Printf("%s\n", ui.colorize(strings.Join(help, " | "), "gray"))
}
