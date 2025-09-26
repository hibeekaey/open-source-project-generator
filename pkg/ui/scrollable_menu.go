// Package ui provides scrollable menu functionality for better navigation in long lists.
//
// This file implements scrollable menu components that handle proper keyboard navigation
// and viewport management for lists that exceed the terminal height.
package ui

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"golang.org/x/term"
)

// ScrollableMenu provides a scrollable menu interface with proper keyboard navigation
type ScrollableMenu struct {
	ui       *InteractiveUI
	viewport ViewportManager
}

// ViewportManager manages the visible portion of a long list
type ViewportManager struct {
	totalItems     int
	visibleItems   int
	currentIndex   int
	scrollOffset   int
	terminalHeight int
}

// NewScrollableMenu creates a new scrollable menu instance
func NewScrollableMenu(ui *InteractiveUI) *ScrollableMenu {
	height := getTerminalHeight()
	// Reserve space for header, footer, and navigation help
	visibleItems := height - 8
	if visibleItems < 5 {
		visibleItems = 5 // Minimum visible items
	}

	return &ScrollableMenu{
		ui: ui,
		viewport: ViewportManager{
			visibleItems:   visibleItems,
			terminalHeight: height,
		},
	}
}

// ShowScrollableMenu displays a scrollable menu with proper navigation
func (sm *ScrollableMenu) ShowScrollableMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("menu must have at least one option")
	}

	sm.viewport.totalItems = len(config.Options)
	sm.viewport.currentIndex = config.DefaultItem
	if sm.viewport.currentIndex < 0 || sm.viewport.currentIndex >= sm.viewport.totalItems {
		sm.viewport.currentIndex = 0
	}

	for {
		select {
		case <-ctx.Done():
			return &interfaces.MenuResult{Cancelled: true}, ctx.Err()
		default:
		}

		sm.displayScrollableMenu(config)
		sm.showScrollableNavigationHelp(config.AllowBack, config.AllowQuit, config.ShowHelp)

		input, err := sm.readRawInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := sm.handleScrollableMenuInput(input, &config)
		if shouldReturn {
			return result, nil
		}
	}
}

// displayScrollableMenu renders the scrollable menu with viewport management
func (sm *ScrollableMenu) displayScrollableMenu(config interfaces.MenuConfig) {
	sm.ui.clearScreen()
	sm.ui.showHeader(config.Title, config.Description)

	// Update scroll offset based on current index
	sm.updateScrollOffset()

	// Calculate visible range
	startIdx := sm.viewport.scrollOffset
	endIdx := startIdx + sm.viewport.visibleItems
	if endIdx > sm.viewport.totalItems {
		endIdx = sm.viewport.totalItems
	}

	// Show scroll indicator if needed
	if sm.viewport.totalItems > sm.viewport.visibleItems {
		sm.showScrollIndicator()
	}

	fmt.Println()

	// Display visible items
	for i := startIdx; i < endIdx; i++ {
		option := config.Options[i]
		prefix := "  "
		if i == sm.viewport.currentIndex {
			prefix = sm.ui.colorize("> ", "cyan")
		}

		icon := ""
		if option.Icon != "" {
			icon = option.Icon + " "
		}

		label := option.Label
		if option.Disabled {
			label = sm.ui.colorize(label, "gray")
		} else if i == sm.viewport.currentIndex {
			label = sm.ui.colorize(label, "white")
		}

		shortcut := ""
		if option.Shortcut != "" {
			shortcut = sm.ui.colorize(fmt.Sprintf(" [%s]", option.Shortcut), "gray")
		}

		fmt.Printf("%s%s%s%s\n", prefix, icon, label, shortcut)

		// Show description for current item
		if option.Description != "" && i == sm.viewport.currentIndex {
			fmt.Printf("    %s\n", sm.ui.colorize(option.Description, "gray"))
		}
	}

	fmt.Println()
}

// updateScrollOffset updates the scroll offset to keep current item visible
func (sm *ScrollableMenu) updateScrollOffset() {
	// If current item is above visible area, scroll up
	if sm.viewport.currentIndex < sm.viewport.scrollOffset {
		sm.viewport.scrollOffset = sm.viewport.currentIndex
	}

	// If current item is below visible area, scroll down
	if sm.viewport.currentIndex >= sm.viewport.scrollOffset+sm.viewport.visibleItems {
		sm.viewport.scrollOffset = sm.viewport.currentIndex - sm.viewport.visibleItems + 1
	}

	// Ensure scroll offset is within bounds
	if sm.viewport.scrollOffset < 0 {
		sm.viewport.scrollOffset = 0
	}
	maxOffset := sm.viewport.totalItems - sm.viewport.visibleItems
	if maxOffset < 0 {
		maxOffset = 0
	}
	if sm.viewport.scrollOffset > maxOffset {
		sm.viewport.scrollOffset = maxOffset
	}
}

// showScrollIndicator displays scroll position indicator
func (sm *ScrollableMenu) showScrollIndicator() {
	totalPages := (sm.viewport.totalItems + sm.viewport.visibleItems - 1) / sm.viewport.visibleItems
	currentPage := (sm.viewport.currentIndex / sm.viewport.visibleItems) + 1

	indicator := fmt.Sprintf("Page %d of %d", currentPage, totalPages)
	if sm.viewport.scrollOffset > 0 {
		indicator += " ↑"
	}
	if sm.viewport.scrollOffset+sm.viewport.visibleItems < sm.viewport.totalItems {
		indicator += " ↓"
	}

	fmt.Printf("%s\n", sm.ui.colorize(indicator, "gray"))
}

// readRawInput reads raw keyboard input including escape sequences
func (sm *ScrollableMenu) readRawInput() (string, error) {
	// Set terminal to raw mode to capture arrow keys
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// Fallback to regular input if raw mode fails
		return sm.ui.readInput()
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return "", err
	}

	input := string(buf[:n])

	// Handle escape sequences (arrow keys)
	if input == "\x1b[A" {
		return "up", nil
	}
	if input == "\x1b[B" {
		return "down", nil
	}
	if input == "\x1b[C" {
		return "right", nil
	}
	if input == "\x1b[D" {
		return "left", nil
	}

	// Handle other special keys
	if input == "\r" || input == "\n" {
		return "enter", nil
	}
	if input == "\x1b" {
		return "esc", nil
	}
	if input == "\x03" { // Ctrl+C
		return "ctrl+c", nil
	}
	if input == "\x7f" || input == "\b" {
		return "backspace", nil
	}

	// Handle regular characters (don't trim space as it's used for selection)
	return input, nil
}

// handleScrollableMenuInput processes user input for scrollable menu navigation
func (sm *ScrollableMenu) handleScrollableMenuInput(input string, config *interfaces.MenuConfig) (*interfaces.MenuResult, bool) {
	switch input {
	case "q", "quit":
		if config.AllowQuit {
			if sm.ui.config.ConfirmOnQuit {
				if sm.ui.confirmQuit() {
					return &interfaces.MenuResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.MenuResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			sm.ui.showContextHelp(config.HelpText, "menu")
			return nil, false
		}

	case "b", "back", "esc":
		if config.AllowBack {
			return &interfaces.MenuResult{Action: "back", Cancelled: true}, true
		}

	case "enter", " ":
		if sm.viewport.currentIndex >= 0 && sm.viewport.currentIndex < len(config.Options) {
			option := config.Options[sm.viewport.currentIndex]
			if !option.Disabled {
				return &interfaces.MenuResult{
					SelectedIndex: sm.viewport.currentIndex,
					SelectedValue: option.Value,
					Action:        "select",
				}, true
			}
		}

	case "up", "k":
		sm.navigateUp()

	case "down", "j":
		sm.navigateDown()

	case "ctrl+c":
		return &interfaces.MenuResult{Action: "cancel", Cancelled: true}, true

	default:
		// Check for numeric input
		if num, err := strconv.Atoi(input); err == nil {
			if num > 0 && num <= len(config.Options) {
				sm.viewport.currentIndex = num - 1
				option := config.Options[sm.viewport.currentIndex]
				if !option.Disabled {
					return &interfaces.MenuResult{
						SelectedIndex: sm.viewport.currentIndex,
						SelectedValue: option.Value,
						Action:        "select",
					}, true
				}
			}
		}

		// Check for shortcut keys
		for i, option := range config.Options {
			if option.Shortcut != "" && strings.EqualFold(option.Shortcut, input) {
				if !option.Disabled {
					return &interfaces.MenuResult{
						SelectedIndex: i,
						SelectedValue: option.Value,
						Action:        "select",
					}, true
				}
			}
		}
	}

	return nil, false
}

// navigateUp moves selection up with proper scrolling
func (sm *ScrollableMenu) navigateUp() {
	if sm.viewport.currentIndex > 0 {
		sm.viewport.currentIndex--
	} else {
		// Wrap to bottom
		sm.viewport.currentIndex = sm.viewport.totalItems - 1
	}
}

// navigateDown moves selection down with proper scrolling
func (sm *ScrollableMenu) navigateDown() {
	if sm.viewport.currentIndex < sm.viewport.totalItems-1 {
		sm.viewport.currentIndex++
	} else {
		// Wrap to top
		sm.viewport.currentIndex = 0
	}
}

// showScrollableNavigationHelp displays navigation help for scrollable menus
func (sm *ScrollableMenu) showScrollableNavigationHelp(allowBack, allowQuit, showHelp bool) {
	if !sm.ui.config.ShowShortcuts {
		return
	}

	help := []string{}
	help = append(help, "↑/↓ or j/k: Navigate")
	help = append(help, "Enter: Select")

	if sm.viewport.totalItems > sm.viewport.visibleItems {
		help = append(help, "Page Up/Down: Scroll")
	}

	if allowBack {
		help = append(help, "b/Esc: Back")
	}
	if allowQuit {
		help = append(help, "q: Quit")
	}
	if showHelp {
		help = append(help, "h: Help")
	}

	fmt.Printf("%s\n", sm.ui.colorize(strings.Join(help, " | "), "gray"))
}

// getTerminalHeight returns the current terminal height
func getTerminalHeight() int {
	if width, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		_ = width // We only need height
		return height
	}
	return 24 // Default fallback
}

// ShowScrollableMultiSelect displays a scrollable multi-selection interface
func (sm *ScrollableMenu) ShowScrollableMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("multi-select must have at least one option")
	}

	sm.viewport.totalItems = len(config.Options)
	sm.viewport.currentIndex = 0

	searchQuery := ""
	filteredIndices := make([]int, len(config.Options))
	for i := range filteredIndices {
		filteredIndices[i] = i
	}

	for {
		select {
		case <-ctx.Done():
			return &interfaces.MultiSelectResult{Cancelled: true}, ctx.Err()
		default:
		}

		sm.displayScrollableMultiSelect(config, searchQuery, filteredIndices)
		sm.showScrollableMultiSelectHelp(config.AllowBack, config.AllowQuit, config.ShowHelp, config.SearchEnabled)

		input, err := sm.readRawInput()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		result, shouldReturn := sm.handleScrollableMultiSelectInput(input, &config, &searchQuery, &filteredIndices)
		if shouldReturn {
			return result, nil
		}
	}
}

// displayScrollableMultiSelect renders the scrollable multi-selection interface
func (sm *ScrollableMenu) displayScrollableMultiSelect(config interfaces.MultiSelectConfig, searchQuery string, filteredIndices []int) {
	sm.ui.clearScreen()
	sm.ui.showHeader(config.Title, config.Description)

	if config.SearchEnabled && searchQuery != "" {
		fmt.Printf("Search: %s\n\n", sm.ui.colorize(searchQuery, "cyan"))
	}

	selectedCount := 0
	for _, option := range config.Options {
		if option.Selected {
			selectedCount++
		}
	}

	fmt.Printf("Selected: %d", selectedCount)
	if config.MinSelection > 0 {
		fmt.Printf(" (min: %d)", config.MinSelection)
	}
	if config.MaxSelection > 0 {
		fmt.Printf(" (max: %d)", config.MaxSelection)
	}
	fmt.Println()

	// Update viewport for filtered items
	sm.viewport.totalItems = len(filteredIndices)
	sm.updateScrollOffset()

	// Calculate visible range
	startIdx := sm.viewport.scrollOffset
	endIdx := startIdx + sm.viewport.visibleItems
	if endIdx > len(filteredIndices) {
		endIdx = len(filteredIndices)
	}

	// Show scroll indicator if needed
	if len(filteredIndices) > sm.viewport.visibleItems {
		sm.showScrollIndicator()
	}

	fmt.Println()

	// Display visible items
	for i := startIdx; i < endIdx; i++ {
		if i >= len(filteredIndices) {
			continue
		}

		idx := filteredIndices[i]
		if idx >= len(config.Options) {
			continue
		}

		option := config.Options[idx]
		prefix := "  "
		if i == sm.viewport.currentIndex {
			prefix = sm.ui.colorize("> ", "cyan")
		}

		checkbox := "☐"
		if option.Selected {
			checkbox = sm.ui.colorize("☑", "green")
		}

		icon := ""
		if option.Icon != "" {
			icon = option.Icon + " "
		}

		label := option.Label
		if option.Disabled {
			label = sm.ui.colorize(label, "gray")
		} else if i == sm.viewport.currentIndex {
			label = sm.ui.colorize(label, "white")
		}

		fmt.Printf("%s%s %s%s\n", prefix, checkbox, icon, label)

		if option.Description != "" && i == sm.viewport.currentIndex {
			fmt.Printf("    %s\n", sm.ui.colorize(option.Description, "gray"))
		}
	}
	fmt.Println()
}

// handleScrollableMultiSelectInput processes user input for scrollable multi-selection
func (sm *ScrollableMenu) handleScrollableMultiSelectInput(input string, config *interfaces.MultiSelectConfig, searchQuery *string, filteredIndices *[]int) (*interfaces.MultiSelectResult, bool) {
	switch input {
	case "q", "quit":
		if config.AllowQuit {
			if sm.ui.config.ConfirmOnQuit {
				if sm.ui.confirmQuit() {
					return &interfaces.MultiSelectResult{Action: "quit", Cancelled: true}, true
				}
				return nil, false
			}
			return &interfaces.MultiSelectResult{Action: "quit", Cancelled: true}, true
		}

	case "h", "help":
		if config.ShowHelp {
			sm.ui.showContextHelp(config.HelpText, "multiselect")
			return nil, false
		}

	case "b", "back", "esc":
		if config.AllowBack {
			return &interfaces.MultiSelectResult{Action: "back", Cancelled: true}, true
		}

	case "enter":
		selectedIndices := []int{}
		selectedValues := []interface{}{}

		for i, option := range config.Options {
			if option.Selected {
				selectedIndices = append(selectedIndices, i)
				selectedValues = append(selectedValues, option.Value)
			}
		}

		// Validate selection count
		if config.MinSelection > 0 && len(selectedIndices) < config.MinSelection {
			sm.ui.showError(fmt.Sprintf("Please select at least %d options", config.MinSelection))
			return nil, false
		}

		if config.MaxSelection > 0 && len(selectedIndices) > config.MaxSelection {
			sm.ui.showError(fmt.Sprintf("Please select at most %d options", config.MaxSelection))
			return nil, false
		}

		return &interfaces.MultiSelectResult{
			SelectedIndices: selectedIndices,
			SelectedValues:  selectedValues,
			Action:          "confirm",
			SearchQuery:     *searchQuery,
		}, true

	case "up", "k":
		sm.navigateUp()

	case "down", "j":
		sm.navigateDown()

	case " ":
		if sm.viewport.currentIndex >= 0 && sm.viewport.currentIndex < len(*filteredIndices) {
			idx := (*filteredIndices)[sm.viewport.currentIndex]
			if idx < len(config.Options) && !config.Options[idx].Disabled {
				config.Options[idx].Selected = !config.Options[idx].Selected
			}
		}

	case "/":
		if config.SearchEnabled {
			*searchQuery = sm.promptSearch()
			*filteredIndices = sm.filterOptions(config.Options, *searchQuery)
			sm.viewport.currentIndex = 0
			sm.viewport.totalItems = len(*filteredIndices)
		}

	case "ctrl+c":
		return &interfaces.MultiSelectResult{Action: "cancel", Cancelled: true}, true

	default:
		// Handle search input if search is enabled
		if config.SearchEnabled && len(input) > 0 && !strings.HasPrefix(strings.ToLower(input), "/") {
			*searchQuery = input
			*filteredIndices = sm.filterOptions(config.Options, *searchQuery)
			sm.viewport.currentIndex = 0
			sm.viewport.totalItems = len(*filteredIndices)
		}
	}

	return nil, false
}

// showScrollableMultiSelectHelp displays navigation help for scrollable multi-select
func (sm *ScrollableMenu) showScrollableMultiSelectHelp(allowBack, allowQuit, showHelp, searchEnabled bool) {
	if !sm.ui.config.ShowShortcuts {
		return
	}

	help := []string{}
	help = append(help, "↑/↓: Navigate")
	help = append(help, "Space: Toggle")
	help = append(help, "Enter: Confirm")

	if searchEnabled {
		help = append(help, "/: Search")
	}
	if allowBack {
		help = append(help, "b/Esc: Back")
	}
	if allowQuit {
		help = append(help, "q: Quit")
	}
	if showHelp {
		help = append(help, "h: Help")
	}

	fmt.Printf("%s\n", sm.ui.colorize(strings.Join(help, " | "), "gray"))
}

// promptSearch prompts for search input
func (sm *ScrollableMenu) promptSearch() string {
	fmt.Print("Search: ")
	input, _ := sm.ui.readInput()
	return strings.TrimSpace(input)
}

// filterOptions filters options based on search query
func (sm *ScrollableMenu) filterOptions(options []interfaces.SelectOption, query string) []int {
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
