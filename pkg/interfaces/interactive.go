// Package interfaces defines interactive UI contracts for the CLI generator.
//
// This file contains interface definitions for interactive user interface components
// that enable guided project generation with menus, selections, and input validation.
package interfaces

import (
	"context"
	"time"
)

// InteractiveUIInterface defines the contract for interactive user interface operations.
//
// This interface provides comprehensive interactive functionality including:
//   - Menu navigation and selection
//   - Input collection with validation
//   - Progress tracking and feedback
//   - Error handling with recovery options
//   - Context-sensitive help system
type InteractiveUIInterface interface {
	// Menu and selection operations
	ShowMenu(ctx context.Context, config MenuConfig) (*MenuResult, error)
	ShowMultiSelect(ctx context.Context, config MultiSelectConfig) (*MultiSelectResult, error)
	ShowCheckboxList(ctx context.Context, config CheckboxConfig) (*CheckboxResult, error)

	// Input collection operations
	PromptText(ctx context.Context, config TextPromptConfig) (*TextResult, error)
	PromptConfirm(ctx context.Context, config ConfirmConfig) (*ConfirmResult, error)
	PromptSelect(ctx context.Context, config SelectConfig) (*SelectResult, error)

	// Display and formatting operations
	ShowTable(ctx context.Context, config TableConfig) error
	ShowTree(ctx context.Context, config TreeConfig) error
	ShowProgress(ctx context.Context, config ProgressConfig) (ProgressTracker, error)

	// Navigation and help operations
	ShowBreadcrumb(ctx context.Context, path []string) error
	ShowHelp(ctx context.Context, helpContext string) error
	ShowError(ctx context.Context, config ErrorConfig) (*ErrorResult, error)

	// Session management
	StartSession(ctx context.Context, config SessionConfig) (*UISession, error)
	EndSession(ctx context.Context, session *UISession) error
	SaveSessionState(ctx context.Context, session *UISession) error
	RestoreSessionState(ctx context.Context, sessionID string) (*UISession, error)
}

// MenuConfig defines configuration for menu display
type MenuConfig struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Options     []MenuOption `json:"options"`
	DefaultItem int          `json:"default_item,omitempty"`
	AllowBack   bool         `json:"allow_back"`
	AllowQuit   bool         `json:"allow_quit"`
	ShowHelp    bool         `json:"show_help"`
	HelpText    string       `json:"help_text,omitempty"`
	MaxHeight   int          `json:"max_height,omitempty"`
}

// MenuOption represents a single menu option
type MenuOption struct {
	Label       string                 `json:"label"`
	Description string                 `json:"description,omitempty"`
	Value       interface{}            `json:"value"`
	Disabled    bool                   `json:"disabled"`
	Icon        string                 `json:"icon,omitempty"`
	Shortcut    string                 `json:"shortcut,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MenuResult contains the result of menu selection
type MenuResult struct {
	SelectedIndex int         `json:"selected_index"`
	SelectedValue interface{} `json:"selected_value"`
	Action        string      `json:"action"` // "select", "back", "quit", "help"
	Cancelled     bool        `json:"cancelled"`
}

// MultiSelectConfig defines configuration for multi-selection interface
type MultiSelectConfig struct {
	Title         string           `json:"title"`
	Description   string           `json:"description,omitempty"`
	Options       []SelectOption   `json:"options"`
	MinSelection  int              `json:"min_selection,omitempty"`
	MaxSelection  int              `json:"max_selection,omitempty"`
	AllowBack     bool             `json:"allow_back"`
	AllowQuit     bool             `json:"allow_quit"`
	ShowHelp      bool             `json:"show_help"`
	HelpText      string           `json:"help_text,omitempty"`
	SearchEnabled bool             `json:"search_enabled"`
	Categories    map[string][]int `json:"categories,omitempty"`
}

// SelectOption represents a selectable option
type SelectOption struct {
	Label       string                 `json:"label"`
	Description string                 `json:"description,omitempty"`
	Value       interface{}            `json:"value"`
	Selected    bool                   `json:"selected"`
	Disabled    bool                   `json:"disabled"`
	Icon        string                 `json:"icon,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MultiSelectResult contains the result of multi-selection
type MultiSelectResult struct {
	SelectedIndices []int         `json:"selected_indices"`
	SelectedValues  []interface{} `json:"selected_values"`
	Action          string        `json:"action"` // "confirm", "back", "quit", "help"
	Cancelled       bool          `json:"cancelled"`
	SearchQuery     string        `json:"search_query,omitempty"`
}

// CheckboxConfig defines configuration for checkbox list
type CheckboxConfig struct {
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Items       []CheckboxItem `json:"items"`
	AllowBack   bool           `json:"allow_back"`
	AllowQuit   bool           `json:"allow_quit"`
	ShowHelp    bool           `json:"show_help"`
	HelpText    string         `json:"help_text,omitempty"`
}

// CheckboxItem represents a checkbox item
type CheckboxItem struct {
	Label       string                 `json:"label"`
	Description string                 `json:"description,omitempty"`
	Value       interface{}            `json:"value"`
	Checked     bool                   `json:"checked"`
	Disabled    bool                   `json:"disabled"`
	Required    bool                   `json:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CheckboxResult contains the result of checkbox selection
type CheckboxResult struct {
	CheckedIndices []int         `json:"checked_indices"`
	CheckedValues  []interface{} `json:"checked_values"`
	Action         string        `json:"action"` // "confirm", "back", "quit", "help"
	Cancelled      bool          `json:"cancelled"`
}

// TextPromptConfig defines configuration for text input
type TextPromptConfig struct {
	Prompt       string             `json:"prompt"`
	Description  string             `json:"description,omitempty"`
	DefaultValue string             `json:"default_value,omitempty"`
	Placeholder  string             `json:"placeholder,omitempty"`
	Required     bool               `json:"required"`
	Multiline    bool               `json:"multiline"`
	Masked       bool               `json:"masked"` // for passwords
	Validator    func(string) error `json:"-"`
	Suggestions  []string           `json:"suggestions,omitempty"`
	AllowBack    bool               `json:"allow_back"`
	AllowQuit    bool               `json:"allow_quit"`
	ShowHelp     bool               `json:"show_help"`
	HelpText     string             `json:"help_text,omitempty"`
	MaxLength    int                `json:"max_length,omitempty"`
	MinLength    int                `json:"min_length,omitempty"`
}

// TextResult contains the result of text input
type TextResult struct {
	Value     string `json:"value"`
	Action    string `json:"action"` // "submit", "back", "quit", "help"
	Cancelled bool   `json:"cancelled"`
}

// ConfirmConfig defines configuration for confirmation prompt
type ConfirmConfig struct {
	Prompt       string `json:"prompt"`
	Description  string `json:"description,omitempty"`
	DefaultValue bool   `json:"default_value"`
	YesLabel     string `json:"yes_label,omitempty"`
	NoLabel      string `json:"no_label,omitempty"`
	AllowBack    bool   `json:"allow_back"`
	AllowQuit    bool   `json:"allow_quit"`
	ShowHelp     bool   `json:"show_help"`
	HelpText     string `json:"help_text,omitempty"`
}

// ConfirmResult contains the result of confirmation
type ConfirmResult struct {
	Confirmed bool   `json:"confirmed"`
	Action    string `json:"action"` // "confirm", "back", "quit", "help"
	Cancelled bool   `json:"cancelled"`
}

// SelectConfig defines configuration for single selection
type SelectConfig struct {
	Prompt      string   `json:"prompt"`
	Description string   `json:"description,omitempty"`
	Options     []string `json:"options"`
	DefaultItem int      `json:"default_item,omitempty"`
	AllowBack   bool     `json:"allow_back"`
	AllowQuit   bool     `json:"allow_quit"`
	ShowHelp    bool     `json:"show_help"`
	HelpText    string   `json:"help_text,omitempty"`
}

// SelectResult contains the result of selection
type SelectResult struct {
	SelectedIndex int    `json:"selected_index"`
	SelectedValue string `json:"selected_value"`
	Action        string `json:"action"` // "select", "back", "quit", "help"
	Cancelled     bool   `json:"cancelled"`
}

// TableConfig defines configuration for table display
type TableConfig struct {
	Title      string     `json:"title,omitempty"`
	Headers    []string   `json:"headers"`
	Rows       [][]string `json:"rows"`
	MaxWidth   int        `json:"max_width,omitempty"`
	Pagination bool       `json:"pagination"`
	PageSize   int        `json:"page_size,omitempty"`
	Sortable   bool       `json:"sortable"`
	Searchable bool       `json:"searchable"`
}

// TreeConfig defines configuration for tree display
type TreeConfig struct {
	Title      string   `json:"title,omitempty"`
	Root       TreeNode `json:"root"`
	Expandable bool     `json:"expandable"`
	ShowIcons  bool     `json:"show_icons"`
	MaxDepth   int      `json:"max_depth,omitempty"`
}

// TreeNode represents a node in a tree structure
type TreeNode struct {
	Label      string                 `json:"label"`
	Value      interface{}            `json:"value,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Children   []TreeNode             `json:"children,omitempty"`
	Expanded   bool                   `json:"expanded"`
	Selectable bool                   `json:"selectable"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ProgressConfig defines configuration for progress tracking
type ProgressConfig struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Steps       []string `json:"steps,omitempty"`
	ShowPercent bool     `json:"show_percent"`
	ShowETA     bool     `json:"show_eta"`
	Cancellable bool     `json:"cancellable"`
}

// ProgressTracker interface for progress tracking operations
type ProgressTracker interface {
	// Update progress (0.0 to 1.0)
	SetProgress(progress float64) error

	// Set current step
	SetCurrentStep(step int, description string) error

	// Add log message
	AddLog(message string) error

	// Mark as complete
	Complete() error

	// Mark as failed
	Fail(err error) error

	// Check if cancelled
	IsCancelled() bool

	// Close the progress tracker
	Close() error
}

// ErrorConfig defines configuration for error display
type ErrorConfig struct {
	Title           string           `json:"title,omitempty"`
	Message         string           `json:"message"`
	Details         string           `json:"details,omitempty"`
	ErrorType       string           `json:"error_type,omitempty"`
	Suggestions     []string         `json:"suggestions,omitempty"`
	RecoveryOptions []RecoveryOption `json:"recovery_options,omitempty"`
	ShowStack       bool             `json:"show_stack"`
	AllowRetry      bool             `json:"allow_retry"`
	AllowIgnore     bool             `json:"allow_ignore"`
	AllowBack       bool             `json:"allow_back"`
	AllowQuit       bool             `json:"allow_quit"`
}

// RecoveryOption represents a recovery action for errors
type RecoveryOption struct {
	Label       string                 `json:"label"`
	Description string                 `json:"description"`
	Action      func() error           `json:"-"`
	Safe        bool                   `json:"safe"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ErrorResult contains the result of error handling
type ErrorResult struct {
	Action           string `json:"action"` // "retry", "ignore", "back", "quit", "recovery"
	RecoverySelected int    `json:"recovery_selected,omitempty"`
	Cancelled        bool   `json:"cancelled"`
}

// SessionConfig defines configuration for UI session
type SessionConfig struct {
	SessionID   string                 `json:"session_id,omitempty"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Timeout     time.Duration          `json:"timeout,omitempty"`
	AutoSave    bool                   `json:"auto_save"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UISession represents an interactive UI session
type UISession struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	StartTime  time.Time              `json:"start_time"`
	LastActive time.Time              `json:"last_active"`
	State      map[string]interface{} `json:"state"`
	History    []SessionAction        `json:"history"`
	Context    context.Context        `json:"-"`
	CancelFunc context.CancelFunc     `json:"-"`
}

// SessionAction represents an action taken during a session
type SessionAction struct {
	Timestamp time.Time              `json:"timestamp"`
	Action    string                 `json:"action"`
	Component string                 `json:"component"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// KeyboardShortcut defines a keyboard shortcut
type KeyboardShortcut struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Action      string `json:"action"`
	Global      bool   `json:"global"` // Available in all contexts
}

// NavigationAction represents navigation actions
type NavigationAction string

const (
	NavigationActionBack   NavigationAction = "back"
	NavigationActionNext   NavigationAction = "next"
	NavigationActionQuit   NavigationAction = "quit"
	NavigationActionHelp   NavigationAction = "help"
	NavigationActionRetry  NavigationAction = "retry"
	NavigationActionIgnore NavigationAction = "ignore"
	NavigationActionCancel NavigationAction = "cancel"
)

// ValidationError represents a validation error with recovery options
type ValidationError struct {
	Field           string           `json:"field"`
	Value           string           `json:"value"`
	Message         string           `json:"message"`
	Code            string           `json:"code"`
	Suggestions     []string         `json:"suggestions,omitempty"`
	RecoveryOptions []RecoveryOption `json:"recovery_options,omitempty"`
}

// Error implements the error interface
func (v *ValidationError) Error() string {
	return v.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, value, message, code string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Code:    code,
	}
}

// WithSuggestions adds suggestions to a validation error
func (v *ValidationError) WithSuggestions(suggestions ...string) *ValidationError {
	v.Suggestions = append(v.Suggestions, suggestions...)
	return v
}

// WithRecoveryOptions adds recovery options to a validation error
func (v *ValidationError) WithRecoveryOptions(options ...RecoveryOption) *ValidationError {
	v.RecoveryOptions = append(v.RecoveryOptions, options...)
	return v
}

// TemplateSelection represents a selected template with options
type TemplateSelection struct {
	Template TemplateInfo           `json:"template"`
	Selected bool                   `json:"selected"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ProjectStructurePreview represents a preview of the project structure
type ProjectStructurePreview struct {
	RootDirectory string             `json:"root_directory"`
	Structure     []DirectoryNode    `json:"structure"`
	FileCount     int                `json:"file_count"`
	EstimatedSize int64              `json:"estimated_size"`
	Components    []ComponentSummary `json:"components"`
}

// DirectoryNode represents a node in the directory structure
type DirectoryNode struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"` // "directory" or "file"
	Children    []DirectoryNode `json:"children,omitempty"`
	Source      string          `json:"source,omitempty"` // which template generated this
	Description string          `json:"description,omitempty"`
}

// ComponentSummary represents a summary of a project component
type ComponentSummary struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Files        []string `json:"files"`
	Dependencies []string `json:"dependencies"`
}
