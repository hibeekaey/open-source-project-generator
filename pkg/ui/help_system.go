// Package ui provides context-sensitive help system for interactive CLI generation.
//
// This file implements the HelpSystem which provides context-aware help information,
// error recovery options, and completion summaries throughout the interactive
// project generation workflow.
package ui

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// HelpSystem manages context-sensitive help and error recovery
type HelpSystem struct {
	ui            interfaces.InteractiveUIInterface
	logger        interfaces.Logger
	config        *HelpConfig
	helpContexts  map[string]*HelpContext
	errorHandlers map[string]ErrorHandler
}

// HelpConfig defines configuration for the help system
type HelpConfig struct {
	ShowContextualHelp    bool `json:"show_contextual_help"`
	ShowExamples          bool `json:"show_examples"`
	ShowTroubleshooting   bool `json:"show_troubleshooting"`
	ShowRecoveryOptions   bool `json:"show_recovery_options"`
	AutoShowOnError       bool `json:"auto_show_on_error"`
	DetailedErrorMessages bool `json:"detailed_error_messages"`
}

// HelpContext defines help information for a specific context
type HelpContext struct {
	Name            string                        `json:"name"`
	Title           string                        `json:"title"`
	Description     string                        `json:"description"`
	Instructions    []string                      `json:"instructions"`
	Examples        []HelpExample                 `json:"examples"`
	Shortcuts       []interfaces.KeyboardShortcut `json:"shortcuts"`
	Troubleshooting []TroubleshootingItem         `json:"troubleshooting"`
	RelatedTopics   []string                      `json:"related_topics"`
	Metadata        map[string]interface{}        `json:"metadata"`
}

// HelpExample represents a help example
type HelpExample struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Notes       string `json:"notes,omitempty"`
}

// TroubleshootingItem represents a troubleshooting entry
type TroubleshootingItem struct {
	Problem    string   `json:"problem"`
	Cause      string   `json:"cause"`
	Solutions  []string `json:"solutions"`
	Prevention string   `json:"prevention,omitempty"`
}

// ErrorHandler defines how to handle specific types of errors
type ErrorHandler struct {
	ErrorType       string                      `json:"error_type"`
	RecoveryOptions []interfaces.RecoveryOption `json:"recovery_options"`
	HelpText        string                      `json:"help_text"`
	AutoRecover     bool                        `json:"auto_recover"`
}

// CompletionSummary represents a completion summary
type CompletionSummary struct {
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	GeneratedItems []GeneratedItem        `json:"generated_items"`
	NextSteps      []NextStep             `json:"next_steps"`
	AdditionalInfo []string               `json:"additional_info"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// GeneratedItem represents an item that was generated
type GeneratedItem struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Size        string `json:"size,omitempty"`
}

// NextStep represents a suggested next step
type NextStep struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Optional    bool   `json:"optional"`
}

// NewHelpSystem creates a new help system
func NewHelpSystem(ui interfaces.InteractiveUIInterface, logger interfaces.Logger, config *HelpConfig) *HelpSystem {
	if config == nil {
		config = &HelpConfig{
			ShowContextualHelp:    true,
			ShowExamples:          true,
			ShowTroubleshooting:   true,
			ShowRecoveryOptions:   true,
			AutoShowOnError:       true,
			DetailedErrorMessages: true,
		}
	}

	hs := &HelpSystem{
		ui:            ui,
		logger:        logger,
		config:        config,
		helpContexts:  make(map[string]*HelpContext),
		errorHandlers: make(map[string]ErrorHandler),
	}

	hs.setupDefaultHelpContexts()
	hs.setupDefaultErrorHandlers()
	return hs
}

// setupDefaultHelpContexts initializes default help contexts
func (hs *HelpSystem) setupDefaultHelpContexts() {
	contexts := []*HelpContext{
		{
			Name:        "menu",
			Title:       "Menu Navigation Help",
			Description: "Learn how to navigate and select from menus",
			Instructions: []string{
				"Use arrow keys (↑/↓) or j/k to navigate between options",
				"Press Enter to select the highlighted option",
				"Type a number to jump directly to that option",
				"Use shortcut keys when available (shown in brackets)",
				"Press 'h' for help, 'b' to go back, 'q' to quit",
			},
			Examples: []HelpExample{
				{
					Title:       "Selecting an Option",
					Description: "Navigate to an option and press Enter",
					Input:       "Use ↓ to highlight 'Backend Templates', then press Enter",
					Output:      "Opens the backend templates submenu",
				},
				{
					Title:       "Using Shortcuts",
					Description: "Use shortcut keys for quick selection",
					Input:       "Press 'b' for Backend Templates",
					Output:      "Directly opens backend templates without navigation",
				},
			},
			Shortcuts: []interfaces.KeyboardShortcut{
				{Key: "↑/↓", Description: "Navigate options", Action: "navigate"},
				{Key: "j/k", Description: "Navigate options (vim-style)", Action: "navigate"},
				{Key: "Enter", Description: "Select option", Action: "select"},
				{Key: "1-9", Description: "Jump to option", Action: "jump"},
			},
		},
		{
			Name:        "multiselect",
			Title:       "Multi-Selection Help",
			Description: "Learn how to select multiple items from a list",
			Instructions: []string{
				"Use arrow keys (↑/↓) to navigate between options",
				"Press Space to toggle selection of the current option",
				"Press Enter when you're done selecting",
				"Type '/' to search and filter options",
				"Type 'clear' to clear the search filter",
			},
			Examples: []HelpExample{
				{
					Title:       "Selecting Multiple Templates",
					Description: "Select both frontend and backend templates",
					Input:       "Navigate to 'React App', press Space, navigate to 'Go API', press Space, then Enter",
					Output:      "Both React App and Go API templates are selected",
				},
				{
					Title:       "Using Search",
					Description: "Find specific templates quickly",
					Input:       "Type '/' then 'react' to filter React-related templates",
					Output:      "Only React templates are shown in the list",
				},
			},
		},
		{
			Name:        "text_input",
			Title:       "Text Input Help",
			Description: "Learn how to enter and validate text input",
			Instructions: []string{
				"Type your input and press Enter to submit",
				"Leave empty to use the default value (if shown in brackets)",
				"Required fields are marked with an asterisk (*)",
				"Follow any format requirements shown in the prompt",
			},
			Examples: []HelpExample{
				{
					Title:       "Project Name",
					Description: "Enter a valid project name",
					Input:       "my-awesome-project",
					Output:      "Project name set to 'my-awesome-project'",
					Notes:       "Use lowercase letters, numbers, and hyphens only",
				},
			},
			Troubleshooting: []TroubleshootingItem{
				{
					Problem: "Input validation failed",
					Cause:   "Input doesn't meet the required format",
					Solutions: []string{
						"Check the format requirements in the prompt",
						"Remove special characters if not allowed",
						"Ensure required fields are not empty",
					},
				},
			},
		},
		{
			Name:        "directory_selection",
			Title:       "Directory Selection Help",
			Description: "Learn how to select and configure output directories",
			Instructions: []string{
				"Enter the full path where you want to generate the project",
				"Use absolute paths (starting with /) or relative paths",
				"The directory will be created if it doesn't exist",
				"Existing directories will show a warning before overwriting",
			},
			Examples: []HelpExample{
				{
					Title:       "Absolute Path",
					Description: "Use a full system path",
					Input:       "/home/user/projects/my-project",
					Output:      "Project will be generated in /home/user/projects/my-project",
				},
				{
					Title:       "Relative Path",
					Description: "Use a path relative to current directory",
					Input:       "./output/my-project",
					Output:      "Project will be generated in ./output/my-project",
				},
			},
		},
		{
			Name:        "preview",
			Title:       "Project Preview Help",
			Description: "Learn how to review and navigate the project structure preview",
			Instructions: []string{
				"Review the directory structure that will be created",
				"Check file counts and estimated project size",
				"Verify that selected templates are included",
				"Use navigation options to modify selections if needed",
			},
			Examples: []HelpExample{
				{
					Title:       "Reviewing Structure",
					Description: "Check the generated project structure",
					Input:       "Review the tree structure shown",
					Output:      "Understand what files and directories will be created",
				},
			},
		},
	}

	for _, context := range contexts {
		hs.helpContexts[context.Name] = context
	}
}

// setupDefaultErrorHandlers initializes default error handlers
func (hs *HelpSystem) setupDefaultErrorHandlers() {
	handlers := []ErrorHandler{
		{
			ErrorType: "validation_error",
			RecoveryOptions: []interfaces.RecoveryOption{
				{
					Label:       "Retry Input",
					Description: "Try entering the input again with corrections",
					Safe:        true,
				},
				{
					Label:       "Use Default",
					Description: "Use the default value if available",
					Safe:        true,
				},
				{
					Label:       "Skip Field",
					Description: "Skip this field and continue (if optional)",
					Safe:        false,
				},
			},
			HelpText: "Input validation failed. Please check the format requirements and try again.",
		},
		{
			ErrorType: "file_system_error",
			RecoveryOptions: []interfaces.RecoveryOption{
				{
					Label:       "Retry Operation",
					Description: "Try the file operation again",
					Safe:        true,
				},
				{
					Label:       "Choose Different Path",
					Description: "Select a different directory or file path",
					Safe:        true,
				},
				{
					Label:       "Create Directory",
					Description: "Create the missing directory structure",
					Safe:        true,
				},
			},
			HelpText: "File system operation failed. Check permissions and path validity.",
		},
		{
			ErrorType: "template_error",
			RecoveryOptions: []interfaces.RecoveryOption{
				{
					Label:       "Retry Template Processing",
					Description: "Try processing the template again",
					Safe:        true,
				},
				{
					Label:       "Skip Template",
					Description: "Skip this template and continue with others",
					Safe:        false,
				},
				{
					Label:       "Choose Different Template",
					Description: "Go back and select a different template",
					Safe:        true,
				},
			},
			HelpText: "Template processing failed. This may be due to template corruption or missing dependencies.",
		},
	}

	for _, handler := range handlers {
		hs.errorHandlers[handler.ErrorType] = handler
	}
}

// ShowContextHelp displays help for a specific context
func (hs *HelpSystem) ShowContextHelp(ctx context.Context, contextName string, additionalInfo ...string) error {
	if !hs.config.ShowContextualHelp {
		return nil
	}

	helpContext, exists := hs.helpContexts[contextName]
	if !exists {
		return hs.showGenericHelp(ctx, contextName, additionalInfo...)
	}

	return hs.displayHelpContext(ctx, helpContext, additionalInfo...)
}

// displayHelpContext displays a help context
func (hs *HelpSystem) displayHelpContext(ctx context.Context, helpContext *HelpContext, additionalInfo ...string) error {
	// Clear screen and show title
	fmt.Printf("\n%s\n", hs.colorize(helpContext.Title, "bold"))
	fmt.Printf("%s\n\n", hs.colorize(helpContext.Description, "gray"))

	// Show additional info if provided
	for _, info := range additionalInfo {
		fmt.Printf("%s\n", info)
	}
	if len(additionalInfo) > 0 {
		fmt.Println()
	}

	// Show instructions
	if len(helpContext.Instructions) > 0 {
		fmt.Printf("%s\n", hs.colorize("Instructions:", "bold"))
		for _, instruction := range helpContext.Instructions {
			fmt.Printf("  • %s\n", instruction)
		}
		fmt.Println()
	}

	// Show examples if enabled
	if hs.config.ShowExamples && len(helpContext.Examples) > 0 {
		fmt.Printf("%s\n", hs.colorize("Examples:", "bold"))
		for _, example := range helpContext.Examples {
			fmt.Printf("  %s\n", hs.colorize(example.Title, "cyan"))
			fmt.Printf("    %s\n", example.Description)
			if example.Input != "" {
				fmt.Printf("    Input: %s\n", hs.colorize(example.Input, "yellow"))
			}
			if example.Output != "" {
				fmt.Printf("    Result: %s\n", hs.colorize(example.Output, "green"))
			}
			if example.Notes != "" {
				fmt.Printf("    Note: %s\n", hs.colorize(example.Notes, "gray"))
			}
			fmt.Println()
		}
	}

	// Show keyboard shortcuts
	if len(helpContext.Shortcuts) > 0 {
		fmt.Printf("%s\n", hs.colorize("Keyboard Shortcuts:", "bold"))
		for _, shortcut := range helpContext.Shortcuts {
			fmt.Printf("  %s: %s\n", hs.colorize(shortcut.Key, "cyan"), shortcut.Description)
		}
		fmt.Println()
	}

	// Show troubleshooting if enabled
	if hs.config.ShowTroubleshooting && len(helpContext.Troubleshooting) > 0 {
		fmt.Printf("%s\n", hs.colorize("Troubleshooting:", "bold"))
		for _, item := range helpContext.Troubleshooting {
			fmt.Printf("  %s\n", hs.colorize("Problem: "+item.Problem, "red"))
			fmt.Printf("  %s\n", hs.colorize("Cause: "+item.Cause, "yellow"))
			fmt.Printf("  %s\n", hs.colorize("Solutions:", "green"))
			for _, solution := range item.Solutions {
				fmt.Printf("    • %s\n", solution)
			}
			if item.Prevention != "" {
				fmt.Printf("  %s\n", hs.colorize("Prevention: "+item.Prevention, "blue"))
			}
			fmt.Println()
		}
	}

	// Show related topics
	if len(helpContext.RelatedTopics) > 0 {
		fmt.Printf("%s\n", hs.colorize("Related Topics:", "bold"))
		for _, topic := range helpContext.RelatedTopics {
			fmt.Printf("  • %s\n", topic)
		}
		fmt.Println()
	}

	fmt.Print("Press Enter to continue...")
	_, err := fmt.Scanln()
	return err
}

// showGenericHelp displays generic help when no specific context is found
func (hs *HelpSystem) showGenericHelp(ctx context.Context, contextName string, additionalInfo ...string) error {
	fmt.Printf("\n%s\n", hs.colorize("Help", "bold"))
	fmt.Printf("Help for context: %s\n\n", contextName)

	for _, info := range additionalInfo {
		fmt.Printf("%s\n", info)
	}

	fmt.Printf("%s\n", hs.colorize("General Navigation:", "bold"))
	fmt.Println("  • h: Show help")
	fmt.Println("  • b: Go back")
	fmt.Println("  • q: Quit")
	fmt.Println("  • Enter: Confirm/Select")
	fmt.Println("  • Esc: Cancel/Back")
	fmt.Println("  • Ctrl+C: Cancel operation")
	fmt.Println()

	fmt.Print("Press Enter to continue...")
	_, err := fmt.Scanln()
	return err
}

// HandleError handles errors with recovery options
func (hs *HelpSystem) HandleError(ctx context.Context, err error, errorType string) (*interfaces.ErrorResult, error) {
	if !hs.config.AutoShowOnError {
		return nil, err
	}

	handler, exists := hs.errorHandlers[errorType]
	if !exists {
		return hs.handleGenericError(ctx, err)
	}

	return hs.handleSpecificError(ctx, err, handler)
}

// handleSpecificError handles errors with specific error handlers
func (hs *HelpSystem) handleSpecificError(ctx context.Context, err error, handler ErrorHandler) (*interfaces.ErrorResult, error) {
	errorConfig := interfaces.ErrorConfig{
		Title:           "Error Occurred",
		Message:         err.Error(),
		ErrorType:       handler.ErrorType,
		RecoveryOptions: handler.RecoveryOptions,
		ShowStack:       false,
		AllowRetry:      true,
		AllowIgnore:     false,
		AllowBack:       true,
		AllowQuit:       true,
	}

	if hs.config.DetailedErrorMessages && handler.HelpText != "" {
		errorConfig.Details = handler.HelpText
	}

	return hs.ui.ShowError(ctx, errorConfig)
}

// handleGenericError handles errors without specific handlers
func (hs *HelpSystem) handleGenericError(ctx context.Context, err error) (*interfaces.ErrorResult, error) {
	errorConfig := interfaces.ErrorConfig{
		Title:     "Error Occurred",
		Message:   err.Error(),
		ErrorType: "generic_error",
		RecoveryOptions: []interfaces.RecoveryOption{
			{
				Label:       "Retry",
				Description: "Try the operation again",
				Safe:        true,
			},
			{
				Label:       "Continue",
				Description: "Continue despite the error",
				Safe:        false,
			},
		},
		ShowStack:   false,
		AllowRetry:  true,
		AllowIgnore: true,
		AllowBack:   true,
		AllowQuit:   true,
	}

	return hs.ui.ShowError(ctx, errorConfig)
}

// ShowCompletionSummary displays a completion summary
func (hs *HelpSystem) ShowCompletionSummary(ctx context.Context, summary *CompletionSummary) error {
	fmt.Printf("\n%s\n", hs.colorize(summary.Title, "bold"))
	fmt.Printf("%s\n\n", hs.colorize(summary.Description, "green"))

	// Show generated items
	if len(summary.GeneratedItems) > 0 {
		fmt.Printf("%s\n", hs.colorize("Generated Items:", "bold"))
		for _, item := range summary.GeneratedItems {
			icon := hs.getItemIcon(item.Type)
			fmt.Printf("  %s %s\n", icon, hs.colorize(item.Name, "cyan"))
			fmt.Printf("    Path: %s\n", item.Path)
			if item.Description != "" {
				fmt.Printf("    %s\n", item.Description)
			}
			if item.Size != "" {
				fmt.Printf("    Size: %s\n", hs.colorize(item.Size, "gray"))
			}
		}
		fmt.Println()
	}

	// Show next steps
	if len(summary.NextSteps) > 0 {
		fmt.Printf("%s\n", hs.colorize("Next Steps:", "bold"))
		for i, step := range summary.NextSteps {
			prefix := fmt.Sprintf("%d.", i+1)
			if step.Optional {
				prefix = "•"
			}
			fmt.Printf("  %s %s\n", prefix, hs.colorize(step.Title, "yellow"))
			fmt.Printf("     %s\n", step.Description)
			if step.Command != "" {
				fmt.Printf("     Command: %s\n", hs.colorize(step.Command, "cyan"))
			}
		}
		fmt.Println()
	}

	// Show additional information
	if len(summary.AdditionalInfo) > 0 {
		fmt.Printf("%s\n", hs.colorize("Additional Information:", "bold"))
		for _, info := range summary.AdditionalInfo {
			fmt.Printf("  • %s\n", info)
		}
		fmt.Println()
	}

	return nil
}

// getItemIcon returns an icon for the item type
func (hs *HelpSystem) getItemIcon(itemType string) string {
	icons := map[string]string{
		"directory": "📁",
		"file":      "📄",
		"config":    "⚙️",
		"script":    "📜",
		"template":  "📋",
		"docs":      "📚",
		"test":      "🧪",
		"build":     "🔨",
		"deploy":    "🚀",
	}

	if icon, exists := icons[itemType]; exists {
		return icon
	}
	return "📄"
}

// AddHelpContext adds a new help context
func (hs *HelpSystem) AddHelpContext(context *HelpContext) {
	hs.helpContexts[context.Name] = context
}

// AddErrorHandler adds a new error handler
func (hs *HelpSystem) AddErrorHandler(handler ErrorHandler) {
	hs.errorHandlers[handler.ErrorType] = handler
}

// GetHelpContext retrieves a help context by name
func (hs *HelpSystem) GetHelpContext(name string) (*HelpContext, bool) {
	context, exists := hs.helpContexts[name]
	return context, exists
}

// colorize applies color formatting to text (placeholder implementation)
func (hs *HelpSystem) colorize(text, color string) string {
	// This would use the same colorization logic as InteractiveUI
	// For now, return text as-is
	return text
}
