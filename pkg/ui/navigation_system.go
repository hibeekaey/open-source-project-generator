// Package ui provides comprehensive navigation system for interactive CLI generation.
//
// This file implements the NavigationSystem which provides breadcrumb navigation,
// keyboard shortcuts, help system, and navigation state management throughout
// the interactive project generation workflow.
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// NavigationSystem manages navigation state and provides navigation functionality
type NavigationSystem struct {
	ui        interfaces.InteractiveUIInterface
	logger    interfaces.Logger
	config    *NavigationConfig
	state     *NavigationState
	history   []NavigationStep
	shortcuts map[string]interfaces.KeyboardShortcut
}

// NavigationConfig defines configuration for the navigation system
type NavigationConfig struct {
	ShowBreadcrumbs   bool `json:"show_breadcrumbs"`
	ShowShortcuts     bool `json:"show_shortcuts"`
	ShowStepCounter   bool `json:"show_step_counter"`
	ShowProgress      bool `json:"show_progress"`
	EnableHistory     bool `json:"enable_history"`
	MaxHistorySize    int  `json:"max_history_size"`
	AutoShowHelp      bool `json:"auto_show_help"`
	ConfirmNavigation bool `json:"confirm_navigation"`
}

// NavigationState tracks the current navigation state
type NavigationState struct {
	CurrentStep      string                 `json:"current_step"`
	StepIndex        int                    `json:"step_index"`
	TotalSteps       int                    `json:"total_steps"`
	Breadcrumbs      []string               `json:"breadcrumbs"`
	AvailableActions []NavigationAction     `json:"available_actions"`
	Context          map[string]interface{} `json:"context"`
	CanGoBack        bool                   `json:"can_go_back"`
	CanGoNext        bool                   `json:"can_go_next"`
	CanQuit          bool                   `json:"can_quit"`
	CanShowHelp      bool                   `json:"can_show_help"`
}

// NavigationStep represents a step in the navigation history
type NavigationStep struct {
	StepName    string                 `json:"step_name"`
	StepIndex   int                    `json:"step_index"`
	Timestamp   string                 `json:"timestamp"`
	Action      string                 `json:"action"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Breadcrumbs []string               `json:"breadcrumbs"`
}

// NavigationAction represents available navigation actions
type NavigationAction struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	Global      bool   `json:"global"`
}

// NewNavigationSystem creates a new navigation system
func NewNavigationSystem(ui interfaces.InteractiveUIInterface, logger interfaces.Logger, config *NavigationConfig) *NavigationSystem {
	if config == nil {
		config = &NavigationConfig{
			ShowBreadcrumbs:   true,
			ShowShortcuts:     true,
			ShowStepCounter:   true,
			ShowProgress:      true,
			EnableHistory:     true,
			MaxHistorySize:    50,
			AutoShowHelp:      false,
			ConfirmNavigation: false,
		}
	}

	ns := &NavigationSystem{
		ui:        ui,
		logger:    logger,
		config:    config,
		state:     &NavigationState{},
		history:   make([]NavigationStep, 0),
		shortcuts: make(map[string]interfaces.KeyboardShortcut),
	}

	ns.setupDefaultShortcuts()
	return ns
}

// setupDefaultShortcuts initializes default keyboard shortcuts
func (ns *NavigationSystem) setupDefaultShortcuts() {
	shortcuts := []interfaces.KeyboardShortcut{
		{Key: "h", Description: "Show help", Action: "help", Global: true},
		{Key: "b", Description: "Go back", Action: "back", Global: true},
		{Key: "n", Description: "Next step", Action: "next", Global: true},
		{Key: "q", Description: "Quit", Action: "quit", Global: true},
		{Key: "r", Description: "Retry", Action: "retry", Global: false},
		{Key: "s", Description: "Skip", Action: "skip", Global: false},
		{Key: "ctrl+c", Description: "Cancel", Action: "cancel", Global: true},
		{Key: "?", Description: "Show shortcuts", Action: "shortcuts", Global: true},
		{Key: "enter", Description: "Confirm", Action: "confirm", Global: true},
		{Key: "esc", Description: "Cancel/Back", Action: "back", Global: true},
	}

	for _, shortcut := range shortcuts {
		ns.shortcuts[shortcut.Key] = shortcut
	}
}

// SetCurrentStep sets the current navigation step
func (ns *NavigationSystem) SetCurrentStep(stepName string, stepIndex, totalSteps int) {
	ns.state.CurrentStep = stepName
	ns.state.StepIndex = stepIndex
	ns.state.TotalSteps = totalSteps

	// Update breadcrumbs
	ns.updateBreadcrumbs(stepName)

	// Log navigation step
	if ns.logger != nil {
		ns.logger.DebugWithFields("Navigation step changed", map[string]interface{}{
			"step_name":   stepName,
			"step_index":  stepIndex,
			"total_steps": totalSteps,
		})
	}
}

// updateBreadcrumbs updates the breadcrumb trail
func (ns *NavigationSystem) updateBreadcrumbs(stepName string) {
	// Add current step to breadcrumbs if not already present
	if len(ns.state.Breadcrumbs) == 0 || ns.state.Breadcrumbs[len(ns.state.Breadcrumbs)-1] != stepName {
		ns.state.Breadcrumbs = append(ns.state.Breadcrumbs, stepName)
	}

	// Limit breadcrumb length to prevent overflow
	maxBreadcrumbs := 5
	if len(ns.state.Breadcrumbs) > maxBreadcrumbs {
		ns.state.Breadcrumbs = ns.state.Breadcrumbs[len(ns.state.Breadcrumbs)-maxBreadcrumbs:]
	}
}

// SetAvailableActions sets the available navigation actions for the current step
func (ns *NavigationSystem) SetAvailableActions(actions []NavigationAction) {
	ns.state.AvailableActions = actions

	// Update state flags based on available actions
	ns.state.CanGoBack = false
	ns.state.CanGoNext = false
	ns.state.CanQuit = false
	ns.state.CanShowHelp = false

	for _, action := range actions {
		switch action.Key {
		case "b", "back":
			ns.state.CanGoBack = action.Available
		case "n", "next":
			ns.state.CanGoNext = action.Available
		case "q", "quit":
			ns.state.CanQuit = action.Available
		case "h", "help":
			ns.state.CanShowHelp = action.Available
		}
	}
}

// ShowBreadcrumbs displays the current breadcrumb navigation
func (ns *NavigationSystem) ShowBreadcrumbs(ctx context.Context) error {
	if !ns.config.ShowBreadcrumbs || len(ns.state.Breadcrumbs) == 0 {
		return nil
	}

	return ns.ui.ShowBreadcrumb(ctx, ns.state.Breadcrumbs)
}

// ShowNavigationHeader displays the navigation header with breadcrumbs and step info
func (ns *NavigationSystem) ShowNavigationHeader(ctx context.Context, title, description string) error {
	// Show breadcrumbs
	if err := ns.ShowBreadcrumbs(ctx); err != nil {
		return fmt.Errorf("failed to show breadcrumbs: %w", err)
	}

	// Show step counter if enabled
	if ns.config.ShowStepCounter && ns.state.TotalSteps > 0 {
		stepInfo := fmt.Sprintf("Step %d of %d", ns.state.StepIndex+1, ns.state.TotalSteps)
		fmt.Printf("%s\n", ns.colorize(stepInfo, "gray"))
	}

	// Show title and description
	if title != "" {
		fmt.Printf("%s\n", ns.colorize(title, "bold"))
		if description != "" {
			fmt.Printf("%s\n", ns.colorize(description, "gray"))
		}
		fmt.Println()
	}

	return nil
}

// ShowNavigationFooter displays available navigation actions and shortcuts
func (ns *NavigationSystem) ShowNavigationFooter(ctx context.Context) error {
	if !ns.config.ShowShortcuts {
		return nil
	}

	shortcuts := ns.getAvailableShortcuts()
	if len(shortcuts) == 0 {
		return nil
	}

	shortcutTexts := make([]string, 0, len(shortcuts))
	for _, shortcut := range shortcuts {
		text := fmt.Sprintf("%s: %s", shortcut.Key, shortcut.Description)
		shortcutTexts = append(shortcutTexts, text)
	}

	fmt.Printf("\n%s\n", ns.colorize(strings.Join(shortcutTexts, " | "), "gray"))
	return nil
}

// getAvailableShortcuts returns shortcuts available for the current context
func (ns *NavigationSystem) getAvailableShortcuts() []interfaces.KeyboardShortcut {
	var available []interfaces.KeyboardShortcut

	for _, action := range ns.state.AvailableActions {
		if shortcut, exists := ns.shortcuts[action.Key]; exists && action.Available {
			available = append(available, shortcut)
		}
	}

	// Add global shortcuts that are always available
	for _, shortcut := range ns.shortcuts {
		if shortcut.Global {
			// Check if not already added
			found := false
			for _, existing := range available {
				if existing.Key == shortcut.Key {
					found = true
					break
				}
			}
			if !found {
				available = append(available, shortcut)
			}
		}
	}

	return available
}

// HandleNavigationInput processes navigation input and returns the action
func (ns *NavigationSystem) HandleNavigationInput(input string) (interfaces.NavigationAction, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Check global shortcuts first
	if shortcut, exists := ns.shortcuts[input]; exists && shortcut.Global {
		return ns.convertToNavigationAction(shortcut.Action), nil
	}

	// Check if input matches any available shortcuts
	for _, action := range ns.state.AvailableActions {
		if action.Key == input && action.Available {
			// Map the key to the appropriate action
			if shortcut, exists := ns.shortcuts[action.Key]; exists {
				return ns.convertToNavigationAction(shortcut.Action), nil
			}
		}
	}

	// Handle common navigation patterns
	switch input {
	case "back", "b":
		if ns.state.CanGoBack {
			return interfaces.NavigationActionBack, nil
		}
	case "next", "n":
		if ns.state.CanGoNext {
			return interfaces.NavigationActionNext, nil
		}
	case "quit", "q":
		if ns.state.CanQuit {
			return interfaces.NavigationActionQuit, nil
		}
	case "help", "h", "?":
		if ns.state.CanShowHelp {
			return interfaces.NavigationActionHelp, nil
		}
	case "retry", "r":
		return interfaces.NavigationActionRetry, nil
	case "cancel", "esc", "ctrl+c":
		return interfaces.NavigationActionCancel, nil
	}

	return "", fmt.Errorf("unknown or unavailable navigation action: %s", input)
}

// convertToNavigationAction converts string action to NavigationAction
func (ns *NavigationSystem) convertToNavigationAction(action string) interfaces.NavigationAction {
	switch action {
	case "back":
		return interfaces.NavigationActionBack
	case "next":
		return interfaces.NavigationActionNext
	case "quit":
		return interfaces.NavigationActionQuit
	case "help":
		return interfaces.NavigationActionHelp
	case "retry":
		return interfaces.NavigationActionRetry
	case "ignore":
		return interfaces.NavigationActionIgnore
	case "cancel":
		return interfaces.NavigationActionCancel
	default:
		return interfaces.NavigationAction(action)
	}
}

// AddToHistory adds a navigation step to the history
func (ns *NavigationSystem) AddToHistory(action string, data map[string]interface{}) {
	if !ns.config.EnableHistory {
		return
	}

	step := NavigationStep{
		StepName:    ns.state.CurrentStep,
		StepIndex:   ns.state.StepIndex,
		Timestamp:   fmt.Sprintf("%d", len(ns.history)),
		Action:      action,
		Data:        data,
		Breadcrumbs: make([]string, len(ns.state.Breadcrumbs)),
	}
	copy(step.Breadcrumbs, ns.state.Breadcrumbs)

	ns.history = append(ns.history, step)

	// Limit history size
	if len(ns.history) > ns.config.MaxHistorySize {
		ns.history = ns.history[1:]
	}

	if ns.logger != nil {
		ns.logger.DebugWithFields("Added navigation history entry", map[string]interface{}{
			"step_name": step.StepName,
			"action":    action,
		})
	}
}

// GetHistory returns the navigation history
func (ns *NavigationSystem) GetHistory() []NavigationStep {
	return ns.history
}

// GetCurrentState returns the current navigation state
func (ns *NavigationSystem) GetCurrentState() *NavigationState {
	return ns.state
}

// CanNavigateBack returns true if back navigation is available
func (ns *NavigationSystem) CanNavigateBack() bool {
	return ns.state.CanGoBack
}

// CanNavigateNext returns true if next navigation is available
func (ns *NavigationSystem) CanNavigateNext() bool {
	return ns.state.CanGoNext
}

// CanQuit returns true if quit is available
func (ns *NavigationSystem) CanQuit() bool {
	return ns.state.CanQuit
}

// CanShowHelp returns true if help is available
func (ns *NavigationSystem) CanShowHelp() bool {
	return ns.state.CanShowHelp
}

// SetContext sets context data for the current navigation state
func (ns *NavigationSystem) SetContext(key string, value interface{}) {
	if ns.state.Context == nil {
		ns.state.Context = make(map[string]interface{})
	}
	ns.state.Context[key] = value
}

// GetContext gets context data from the current navigation state
func (ns *NavigationSystem) GetContext(key string) (interface{}, bool) {
	if ns.state.Context == nil {
		return nil, false
	}
	value, exists := ns.state.Context[key]
	return value, exists
}

// Reset resets the navigation system to initial state
func (ns *NavigationSystem) Reset() {
	ns.state = &NavigationState{}
	ns.history = make([]NavigationStep, 0)
}

// colorize applies color formatting to text (placeholder implementation)
func (ns *NavigationSystem) colorize(text, color string) string {
	// This would use the same colorization logic as InteractiveUI
	// For now, return text as-is
	return text
}

// ShowNavigationHelp displays help for navigation commands
func (ns *NavigationSystem) ShowNavigationHelp(ctx context.Context) error {
	helpText := `Navigation Help:

Available Commands:
• h or ?: Show this help
• b: Go back to previous step
• n: Go to next step (if available)
• q: Quit the application
• Enter: Confirm current selection
• Esc: Cancel or go back
• Ctrl+C: Cancel current operation

Navigation:
• Use arrow keys (↑/↓) or j/k to navigate menus
• Use Tab/Shift+Tab to move between form fields
• Type numbers to jump to menu items
• Use shortcut keys when available

Tips:
• Breadcrumbs show your current location
• Step counter shows progress through the workflow
• Available actions are shown at the bottom of each screen
• You can always go back to modify previous selections`

	return ns.ui.ShowHelp(ctx, helpText)
}
