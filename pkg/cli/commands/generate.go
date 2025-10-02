// Package commands provides individual command implementations for the CLI interface.
//
// This module contains the GenerateCommand struct and its associated functionality,
// extracted from the main CLI handlers to improve modularity and maintainability.
package commands

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// ConflictRule represents a flag conflict rule with detailed information
type ConflictRule struct {
	Flags       []string `json:"flags"`
	Description string   `json:"description"`
	Suggestion  string   `json:"suggestion"`
	Examples    []string `json:"examples"`
	Severity    string   `json:"severity"` // "error", "warning", "info"
}

// GenerateCLI defines the CLI methods needed by GenerateCommand.
type GenerateCLI interface {
	// Validation methods
	ValidateGenerateOptions(options interfaces.GenerateOptions) error

	// Mode detection and routing
	DetectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string
	RouteToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error

	// Flag handling
	ApplyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool)

	// Output methods
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	Error(text string) string
	Info(text string) string
}

// GenerateCommand handles the generate command functionality.
//
// The GenerateCommand provides centralized generate command execution including:
//   - Flag parsing and validation
//   - Generate options creation and validation
//   - Mode detection and routing
//   - Error handling and user feedback
type GenerateCommand struct {
	cli GenerateCLI
}

// NewGenerateCommand creates a new GenerateCommand instance.
func NewGenerateCommand(cli GenerateCLI) *GenerateCommand {
	return &GenerateCommand{
		cli: cli,
	}
}

// Execute handles the generate command execution.
func (gc *GenerateCommand) Execute(cmd *cobra.Command, args []string) error {
	// Get flags
	configPath, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	offline, _ := cmd.Flags().GetBool("offline")
	minimal, _ := cmd.Flags().GetBool("minimal")
	template, _ := cmd.Flags().GetString("template")
	updateVersions, _ := cmd.Flags().GetBool("update-versions")
	force, _ := cmd.Flags().GetBool("force")
	skipValidation, _ := cmd.Flags().GetBool("skip-validation")
	backupExisting, _ := cmd.Flags().GetBool("backup-existing")
	includeExamples, _ := cmd.Flags().GetBool("include-examples")
	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// Additional flags (for future implementation)
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	includeOnly, _ := cmd.Flags().GetStringSlice("include-only")
	interactive, _ := cmd.Flags().GetBool("interactive")
	preset, _ := cmd.Flags().GetString("preset")

	// Mode-specific flags
	forceInteractive, _ := cmd.Flags().GetBool("force-interactive")
	forceNonInteractive, _ := cmd.Flags().GetBool("force-non-interactive")
	explicitMode, _ := cmd.Flags().GetString("mode")

	// Validate conflicting mode flags with enhanced error handling
	if err := gc.validateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode); err != nil {
		return err // Return the enhanced error message directly
	}

	// Apply mode overrides
	nonInteractive, interactive = gc.cli.ApplyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)

	// Log additional options for debugging
	gc.logAdditionalOptions(exclude, includeOnly, preset)

	// Create generate options
	options := gc.createGenerateOptions(
		force, minimal, offline, updateVersions, skipValidation,
		backupExisting, includeExamples, outputPath, dryRun, nonInteractive, template,
	)

	// Perform comprehensive validation before generation
	if !options.SkipValidation {
		gc.cli.VerboseOutput("🔍 Validating your configuration...")
		if err := gc.cli.ValidateGenerateOptions(options); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gc.cli.Error("Configuration validation failed."),
				gc.cli.Info("Please check your settings and try again"))
		}
	}

	// Mode detection and routing logic
	mode := gc.cli.DetectGenerationMode(configPath, nonInteractive, interactive, explicitMode)
	gc.cli.VerboseOutput("🎯 Using %s mode for project generation", mode)

	// Route to appropriate generation method based on detected mode
	return gc.cli.RouteToGenerationMethod(mode, configPath, options)
}

// validateModeFlags validates conflicting mode flags using enhanced conflict detection.
func (gc *GenerateCommand) validateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	// Add null safety check
	if gc == nil || gc.cli == nil {
		return fmt.Errorf("generate command not properly initialized")
	}

	// Create enhanced flag conflict detector
	flagConflictDetector := NewGenerateFlagConflictDetector(gc.cli)

	// Validate mode flags using enhanced detection
	return flagConflictDetector.ValidateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)
}

// GenerateFlagConflictDetector provides enhanced flag conflict detection for the generate command.
type GenerateFlagConflictDetector struct {
	cli GenerateCLI
}

// NewGenerateFlagConflictDetector creates a new flag conflict detector for the generate command.
func NewGenerateFlagConflictDetector(cli GenerateCLI) *GenerateFlagConflictDetector {
	return &GenerateFlagConflictDetector{
		cli: cli,
	}
}

// ValidateModeFlags performs comprehensive mode flag validation with enhanced error messages.
func (detector *GenerateFlagConflictDetector) ValidateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	// Create flag state map
	flagState := map[string]bool{
		"--non-interactive":       nonInteractive,
		"--interactive":           interactive,
		"--force-interactive":     forceInteractive,
		"--force-non-interactive": forceNonInteractive,
	}

	// Add explicit mode if provided
	if explicitMode != "" {
		flagState["--mode"] = true
		flagState[fmt.Sprintf("--mode=%s", explicitMode)] = true

		// Validate explicit mode value first
		if err := detector.validateExplicitMode(explicitMode); err != nil {
			return err
		}
	}

	// Get conflict rules specific to generate command
	conflictRules := detector.getGenerateConflictRules()

	// Check for conflicts
	var detectedConflicts []ConflictRule
	for _, rule := range conflictRules {
		if detector.checkConflictRule(rule, flagState) {
			detectedConflicts = append(detectedConflicts, rule)
		}
	}

	// Generate error if conflicts found
	if len(detectedConflicts) > 0 {
		return detector.generateConflictError(detectedConflicts)
	}

	return nil
}

// getGenerateConflictRules returns conflict rules specific to the generate command.
func (detector *GenerateFlagConflictDetector) getGenerateConflictRules() []ConflictRule {
	return []ConflictRule{
		{
			Flags:       []string{"--interactive", "--non-interactive"},
			Description: "Interactive and non-interactive modes cannot be used together",
			Suggestion:  "Choose either interactive mode for guided setup OR non-interactive for automated generation",
			Examples:    []string{"generator generate --interactive", "generator generate --non-interactive", "generator generate --mode=interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--force-interactive", "--force-non-interactive"},
			Description: "Force interactive and force non-interactive modes are mutually exclusive",
			Suggestion:  "Choose either force-interactive to override detection OR force-non-interactive for automation",
			Examples:    []string{"generator generate --force-interactive", "generator generate --force-non-interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--interactive", "--force-non-interactive"},
			Description: "Interactive mode conflicts with forced non-interactive mode",
			Suggestion:  "Use either --interactive for guided setup OR --force-non-interactive for automation",
			Examples:    []string{"generator generate --interactive", "generator generate --force-non-interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--non-interactive", "--force-interactive"},
			Description: "Non-interactive mode conflicts with forced interactive mode",
			Suggestion:  "Use either --non-interactive for automation OR --force-interactive for guided setup",
			Examples:    []string{"generator generate --non-interactive", "generator generate --force-interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--interactive", "--mode"},
			Description: "Interactive flag conflicts with explicit mode specification",
			Suggestion:  "Use either --interactive flag OR --mode=interactive, not both",
			Examples:    []string{"generator generate --interactive", "generator generate --mode=interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--non-interactive", "--mode"},
			Description: "Non-interactive flag conflicts with explicit mode specification",
			Suggestion:  "Use either --non-interactive flag OR --mode=non-interactive, not both",
			Examples:    []string{"generator generate --non-interactive", "generator generate --mode=non-interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--force-interactive", "--mode"},
			Description: "Force-interactive flag conflicts with explicit mode specification",
			Suggestion:  "Use either --force-interactive flag OR --mode=interactive, not both",
			Examples:    []string{"generator generate --force-interactive", "generator generate --mode=interactive"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--force-non-interactive", "--mode"},
			Description: "Force-non-interactive flag conflicts with explicit mode specification",
			Suggestion:  "Use either --force-non-interactive flag OR --mode=non-interactive, not both",
			Examples:    []string{"generator generate --force-non-interactive", "generator generate --mode=non-interactive"},
			Severity:    "error",
		},
	}
}

// checkConflictRule checks if a specific conflict rule is violated.
func (detector *GenerateFlagConflictDetector) checkConflictRule(rule ConflictRule, flagState map[string]bool) bool {
	activeFlags := 0

	// Count how many conflicting flags are active
	for _, flag := range rule.Flags {
		if detector.isFlagActive(flag, flagState) {
			activeFlags++
		}
	}

	// Conflict detected if more than one conflicting flag is active
	return activeFlags > 1
}

// isFlagActive checks if a flag is active, handling special cases.
func (detector *GenerateFlagConflictDetector) isFlagActive(flag string, flagState map[string]bool) bool {
	// Direct flag check
	if active, exists := flagState[flag]; exists && active {
		return true
	}

	// Handle special cases for mode flags with values
	if strings.HasPrefix(flag, "--mode=") {
		if modeValue, exists := flagState["--mode"]; exists && modeValue {
			return true
		}
	}

	return false
}

// generateConflictError creates a comprehensive error message for detected conflicts.
func (detector *GenerateFlagConflictDetector) generateConflictError(conflicts []ConflictRule) error {
	if len(conflicts) == 0 {
		return nil
	}

	var errorMsg strings.Builder

	// Header
	errorMsg.WriteString(fmt.Sprintf("🚫 %s\n\n", detector.cli.Error("Flag conflicts detected")))

	// List each conflict with details
	for i, conflict := range conflicts {
		if i > 0 {
			errorMsg.WriteString("\n")
		}

		// Conflict description
		errorMsg.WriteString(fmt.Sprintf("%s %s: %s\n",
			detector.cli.Info("Conflict"),
			detector.cli.Info(fmt.Sprintf("#%d", i+1)),
			conflict.Description))
		errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
			detector.cli.Info("Conflicting flags"),
			strings.Join(conflict.Flags, ", ")))
		errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
			detector.cli.Info("Suggestion"),
			conflict.Suggestion))

		if len(conflict.Examples) > 0 {
			errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
				detector.cli.Info("Examples"),
				strings.Join(conflict.Examples, ", ")))
		}
	}

	return fmt.Errorf("%s", errorMsg.String())
}

// validateExplicitMode validates the explicit mode value with enhanced error messages.
func (detector *GenerateFlagConflictDetector) validateExplicitMode(explicitMode string) error {
	validModes := []string{"interactive", "non-interactive", "config-file"}
	normalizedMode := strings.ToLower(strings.TrimSpace(explicitMode))

	// Check for exact matches and common variations
	validVariations := map[string]string{
		"interactive":     "interactive",
		"i":               "interactive",
		"non-interactive": "non-interactive",
		"noninteractive":  "non-interactive",
		"ni":              "non-interactive",
		"auto":            "non-interactive",
		"config-file":     "config-file",
		"config":          "config-file",
		"file":            "config-file",
		"cf":              "config-file",
	}

	if _, exists := validVariations[normalizedMode]; exists {
		return nil
	}

	// Mode is invalid, provide helpful error with suggestions
	errorMsg := fmt.Sprintf("🚫 %s %s\n%s: %s\n%s: %s\n%s: %s",
		detector.cli.Error(fmt.Sprintf("'%s' is not a valid mode", explicitMode)),
		detector.cli.Info("Invalid mode specified"),
		detector.cli.Info("Available modes"),
		strings.Join(validModes, ", "),
		detector.cli.Info("Common variations"),
		"i, ni, auto, config, file, cf",
		detector.cli.Info("Example"),
		"generator generate --mode=interactive")

	return fmt.Errorf("%s", errorMsg)
}

// logAdditionalOptions logs additional options for debugging.
func (gc *GenerateCommand) logAdditionalOptions(exclude, includeOnly []string, preset string) {
	if len(exclude) > 0 {
		gc.cli.DebugOutput("Excluding files/directories: %v", exclude)
	}
	if len(includeOnly) > 0 {
		gc.cli.DebugOutput("Including only: %v", includeOnly)
	}
	if preset != "" {
		gc.cli.DebugOutput("Using preset: %s", preset)
	}
}

// createGenerateOptions creates GenerateOptions from command flags.
func (gc *GenerateCommand) createGenerateOptions(
	force, minimal, offline, updateVersions, skipValidation,
	backupExisting, includeExamples bool,
	outputPath string, dryRun, nonInteractive bool, template string,
) interfaces.GenerateOptions {
	options := interfaces.GenerateOptions{
		Force:           force,
		Minimal:         minimal,
		Offline:         offline,
		UpdateVersions:  updateVersions,
		SkipValidation:  skipValidation,
		BackupExisting:  backupExisting,
		IncludeExamples: includeExamples,
		OutputPath:      outputPath,
		DryRun:          dryRun,
		NonInteractive:  nonInteractive,
	}

	if template != "" {
		options.Templates = []string{template}
	}

	return options
}

// SetupFlags sets up the generate command flags.
func (gc *GenerateCommand) SetupFlags(cmd *cobra.Command) {
	// Core generation flags
	cmd.Flags().StringP("config", "c", "", "Path to configuration file")
	cmd.Flags().StringP("output", "o", ".", "Output directory for generated project")
	cmd.Flags().BoolP("dry-run", "n", false, "Show what would be generated without creating files")
	cmd.Flags().Bool("offline", false, "Use cached templates only (no network access)")
	cmd.Flags().Bool("minimal", false, "Generate minimal project structure")
	cmd.Flags().StringP("template", "t", "", "Specific template to use")
	cmd.Flags().Bool("update-versions", false, "Update to latest package versions")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	cmd.Flags().Bool("skip-validation", false, "Skip configuration validation")
	cmd.Flags().Bool("backup-existing", false, "Backup existing files before overwriting")
	cmd.Flags().Bool("include-examples", false, "Include example files and documentation")

	// Mode control flags
	cmd.Flags().Bool("interactive", false, "Force interactive mode")
	cmd.Flags().Bool("force-interactive", false, "Force interactive mode (override detection)")
	cmd.Flags().Bool("force-non-interactive", false, "Force non-interactive mode (override detection)")
	cmd.Flags().String("mode", "", "Explicit mode: interactive, non-interactive, config-based")

	// Advanced flags (for future implementation)
	cmd.Flags().StringSlice("exclude", []string{}, "Exclude specific files or directories")
	cmd.Flags().StringSlice("include-only", []string{}, "Include only specific files or directories")
	cmd.Flags().String("preset", "", "Use predefined configuration preset")
}
