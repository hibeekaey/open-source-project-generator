// Package cli provides integration for interactive directory selection in the CLI.
//
// This file integrates the directory selection functionality into the existing
// CLI generate command workflow.
package cli

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
)

// runInteractiveDirectorySelection handles interactive directory selection
func (c *CLI) runInteractiveDirectorySelection(ctx context.Context, defaultPath string) (string, error) {
	c.VerboseOutput("Starting interactive directory selection")

	// Create directory selector
	directorySelector := ui.NewDirectorySelector(c.interactiveUI, c.logger)

	// Select output directory
	result, err := directorySelector.SelectOutputDirectory(ctx, defaultPath)
	if err != nil {
		return "", fmt.Errorf("directory selection failed: %w", err)
	}

	if result.Cancelled {
		return "", fmt.Errorf("directory selection cancelled by user")
	}

	// Handle directory preparation
	if result.RequiresCreation {
		c.VerboseOutput("Creating output directory: %s", result.Path)
		if err := directorySelector.CreateDirectory(result.Path); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Handle backup if needed
	if result.ConflictResolution == "overwrite" && result.BackupPath != "" {
		c.VerboseOutput("Creating backup: %s", result.BackupPath)
		if err := directorySelector.CreateBackup(result.Path, result.BackupPath); err != nil {
			return "", fmt.Errorf("failed to create backup: %w", err)
		}
		c.QuietOutput("Backup created: %s", result.BackupPath)
	}

	return result.Path, nil
}

// Note: runInteractiveProjectConfiguration already exists in cli.go
// This file only adds the directory selection functionality
