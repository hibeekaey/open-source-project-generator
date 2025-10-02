package collectors

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ConfirmationCollector handles collecting confirmations from users
type ConfirmationCollector struct {
	ui interfaces.InteractiveUIInterface
}

// NewConfirmationCollector creates a new confirmation collector
func NewConfirmationCollector(ui interfaces.InteractiveUIInterface) *ConfirmationCollector {
	return &ConfirmationCollector{
		ui: ui,
	}
}

// ConfirmSave asks if user wants to save configuration
func (cc *ConfirmationCollector) ConfirmSave(ctx context.Context) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Save Configuration",
		Description:  "Would you like to save this configuration for future use?",
		DefaultValue: false,
		AllowBack:    false,
		AllowQuit:    false,
		ShowHelp:     true,
		HelpText: `Saving Configuration:
• Saved configurations can be reused for similar projects
• You can load saved configurations to skip the interactive setup
• Configurations include project metadata and template selections
• Output directory is not saved (you'll be prompted each time)
• You can modify saved configurations before using them`,
	}

	confirmResult, err := cc.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get save confirmation: %w", err)
	}

	return confirmResult.Confirmed && confirmResult.Action == "confirm", nil
}

// ConfirmOverwrite asks if user wants to overwrite existing configuration
func (cc *ConfirmationCollector) ConfirmOverwrite(ctx context.Context, name string) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Overwrite Configuration '%s'", name),
		Description:  "A configuration with this name already exists. Do you want to overwrite it?",
		DefaultValue: false,
		AllowBack:    true,
		AllowQuit:    true,
	}

	result, err := cc.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

// ConfirmLoad asks if user wants to load the selected configuration
func (cc *ConfirmationCollector) ConfirmLoad(ctx context.Context, name string) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Load Configuration '%s'", name),
		Description:  "Load this configuration and use it for project generation?",
		DefaultValue: true,
		AllowBack:    true,
		AllowQuit:    true,
	}

	result, err := cc.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

// ConfirmDelete asks if user wants to delete the selected configuration
func (cc *ConfirmationCollector) ConfirmDelete(ctx context.Context, name string) (bool, error) {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       fmt.Sprintf("Delete Configuration '%s'", name),
		Description:  "Are you sure you want to delete this configuration? This action cannot be undone.",
		DefaultValue: false,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Delete Configuration:
• This will permanently remove the configuration from your saved configurations
• The action cannot be undone
• Consider exporting the configuration first if you might need it later`,
	}

	confirmResult, err := cc.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to get deletion confirmation: %w", err)
	}

	return confirmResult.Confirmed && confirmResult.Action == "confirm", nil
}
