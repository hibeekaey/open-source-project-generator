// Package cli provides flag management functionality for the CLI.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// FlagHandler manages CLI flag setup, validation, and processing.
type FlagHandler struct {
	cli           *CLI
	outputManager *OutputManager
	logger        interfaces.Logger
}

// NewFlagHandler creates a new FlagHandler instance.
func NewFlagHandler(cli *CLI, outputManager *OutputManager, logger interfaces.Logger) *FlagHandler {
	return &FlagHandler{
		cli:           cli,
		outputManager: outputManager,
		logger:        logger,
	}
}

// SetupGlobalFlags adds global flags that apply to all commands.
func (fh *FlagHandler) SetupGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging with detailed operation information")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output (quiet mode)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging with performance metrics")
	rootCmd.PersistentFlags().String("log-level", "info", "Set log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool("log-json", false, "Output logs in JSON format")
	rootCmd.PersistentFlags().Bool("log-caller", false, "Include caller information in logs")
	rootCmd.PersistentFlags().Bool("non-interactive", false, "Run in non-interactive mode")
	rootCmd.PersistentFlags().String("output-format", "text", "Output format (text, json, yaml)")
}

// HandleGlobalFlags processes global flags that apply to all commands.
func (fh *FlagHandler) HandleGlobalFlags(cmd *cobra.Command) error {
	// Add null safety check
	if fh == nil {
		return fmt.Errorf("flag handler not initialized")
	}
	if cmd == nil {
		return fmt.Errorf("command not provided")
	}

	// Get global flags with error handling (check both regular and persistent flags)
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		// Try persistent flags if regular flags don't have it
		verbose, err = cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			verbose = false // Use default if flag not found
		}
	}

	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		quiet, err = cmd.PersistentFlags().GetBool("quiet")
		if err != nil {
			quiet = false
		}
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		debug, err = cmd.PersistentFlags().GetBool("debug")
		if err != nil {
			debug = false
		}
	}

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		logLevel, err = cmd.PersistentFlags().GetString("log-level")
		if err != nil {
			logLevel = "info" // Use default
		}
	}

	logJSON, err := cmd.Flags().GetBool("log-json")
	if err != nil {
		logJSON, err = cmd.PersistentFlags().GetBool("log-json")
		if err != nil {
			logJSON = false
		}
	}

	logCaller, err := cmd.Flags().GetBool("log-caller")
	if err != nil {
		logCaller, err = cmd.PersistentFlags().GetBool("log-caller")
		if err != nil {
			logCaller = false
		}
	}

	nonInteractive, err := cmd.Flags().GetBool("non-interactive")
	if err != nil {
		nonInteractive, err = cmd.PersistentFlags().GetBool("non-interactive")
		if err != nil {
			nonInteractive = false
		}
	}

	outputFormat, err := cmd.Flags().GetString("output-format")
	if err != nil {
		outputFormat, err = cmd.PersistentFlags().GetString("output-format")
		if err != nil {
			outputFormat = "text"
		}
	}

	// Auto-detect non-interactive mode if not explicitly set
	if !nonInteractive {
		nonInteractive = fh.IsNonInteractiveMode(cmd)
		if nonInteractive && fh.cli != nil {
			fh.cli.VerboseOutput("Auto-detected non-interactive mode (CI environment or piped input)")
		}
	}

	// Validate conflicting flags
	if err := fh.ValidateConflictingFlags(verbose, quiet, debug); err != nil {
		return err
	}

	// Set log level based on flags (priority: debug > verbose > explicit level > quiet)
	if debug {
		logLevel = "debug"
		fh.cli.debugMode = true
		fh.outputManager.SetDebugMode(true)
	} else if verbose {
		logLevel = "debug"
		fh.cli.verboseMode = true
		fh.outputManager.SetVerboseMode(true)
	} else if quiet {
		logLevel = "error"
		fh.cli.quietMode = true
		fh.outputManager.SetQuietMode(true)
	}

	// Validate log level
	if err := fh.ValidateLogLevel(logLevel); err != nil {
		return err
	}

	// Validate output format
	if err := fh.ValidateOutputFormat(outputFormat); err != nil {
		return err
	}

	// Configure logger based on flags
	if fh.logger != nil {
		fh.configureLogger(logLevel, logJSON, logCaller, cmd, verbose, quiet, debug, nonInteractive, outputFormat)
	}

	// Store global settings for use in commands
	fh.displayModeMessages(cmd, verbose, quiet, debug, nonInteractive)

	return nil
}

// ValidateConflictingFlags checks for conflicting flag combinations using the enhanced conflict detection system.
func (fh *FlagHandler) ValidateConflictingFlags(verbose, quiet, debug bool) error {
	// Add null safety check
	if fh == nil {
		return fmt.Errorf("flag handler not initialized")
	}

	// Create flag state map for enhanced validation
	flagState := map[string]bool{
		"--verbose": verbose,
		"--quiet":   quiet,
		"--debug":   debug,
	}

	// Use enhanced conflict detection
	return fh.ValidateConflictingFlagsEnhanced(flagState)
}

// ValidateConflictingFlagsEnhanced performs comprehensive flag conflict detection using the conflict matrix.
func (fh *FlagHandler) ValidateConflictingFlagsEnhanced(flagState map[string]bool) error {
	// Add null safety check
	if fh == nil {
		return fmt.Errorf("flag handler not initialized")
	}

	conflictMatrix := NewFlagConflictMatrix()
	var detectedConflicts []ConflictRule

	// Check each conflict rule against current flag state
	for _, rule := range conflictMatrix.Rules {
		if fh.checkConflictRule(rule, flagState) {
			detectedConflicts = append(detectedConflicts, rule)
		}
	}

	// If conflicts found, generate enhanced error message
	if len(detectedConflicts) > 0 {
		return fh.generateConflictError(detectedConflicts)
	}

	return nil
}

// checkConflictRule checks if a specific conflict rule is violated by the current flag state.
func (fh *FlagHandler) checkConflictRule(rule ConflictRule, flagState map[string]bool) bool {
	activeFlags := 0

	// Count how many conflicting flags are active
	for _, flag := range rule.Flags {
		// Handle special cases for mode flags and explicit values
		if fh.isFlagActive(flag, flagState) {
			activeFlags++
		}
	}

	// Conflict detected if more than one conflicting flag is active
	return activeFlags > 1
}

// isFlagActive checks if a flag is active, handling special cases like --mode flags.
func (fh *FlagHandler) isFlagActive(flag string, flagState map[string]bool) bool {
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

	// Handle output format flags
	if strings.HasPrefix(flag, "--output-format=") {
		if formatValue, exists := flagState["--output-format"]; exists && formatValue {
			return true
		}
	}

	return false
}

// generateConflictError creates a comprehensive error message for detected conflicts.
func (fh *FlagHandler) generateConflictError(conflicts []ConflictRule) error {
	if len(conflicts) == 0 {
		return nil
	}

	var errorMsg strings.Builder

	// Header
	if fh.cli != nil {
		errorMsg.WriteString(fmt.Sprintf("🚫 %s\n\n", fh.cli.Error("Flag conflicts detected")))
	} else {
		errorMsg.WriteString("🚫 Flag conflicts detected\n\n")
	}

	// List each conflict with details
	for i, conflict := range conflicts {
		if i > 0 {
			errorMsg.WriteString("\n")
		}

		// Conflict description
		if fh.cli != nil {
			errorMsg.WriteString(fmt.Sprintf("%s %s: %s\n",
				fh.cli.Info("Conflict"),
				fh.cli.Highlight(fmt.Sprintf("#%d", i+1)),
				conflict.Description))
			errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
				fh.cli.Info("Conflicting flags"),
				fh.cli.Dim(strings.Join(conflict.Flags, ", "))))
			errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
				fh.cli.Info("Suggestion"),
				fh.cli.Dim(conflict.Suggestion)))

			if len(conflict.Examples) > 0 {
				errorMsg.WriteString(fmt.Sprintf("%s: %s\n",
					fh.cli.Info("Examples"),
					fh.cli.Dim(strings.Join(conflict.Examples, ", "))))
			}
		} else {
			errorMsg.WriteString(fmt.Sprintf("Conflict #%d: %s\n", i+1, conflict.Description))
			errorMsg.WriteString(fmt.Sprintf("Conflicting flags: %s\n", strings.Join(conflict.Flags, ", ")))
			errorMsg.WriteString(fmt.Sprintf("Suggestion: %s\n", conflict.Suggestion))

			if len(conflict.Examples) > 0 {
				errorMsg.WriteString(fmt.Sprintf("Examples: %s\n", strings.Join(conflict.Examples, ", ")))
			}
		}
	}

	return fmt.Errorf("%s", errorMsg.String())
}

// ValidateAllFlags performs comprehensive validation of all flag combinations.
func (fh *FlagHandler) ValidateAllFlags(cmd *cobra.Command) error {
	// Add null safety check
	if fh == nil {
		return fmt.Errorf("flag handler not initialized")
	}
	if cmd == nil {
		return fmt.Errorf("command not provided")
	}

	// Collect all flag states
	flagState := fh.collectFlagState(cmd)

	// Perform enhanced conflict validation
	return fh.ValidateConflictingFlagsEnhanced(flagState)
}

// collectFlagState safely collects the current state of all flags.
func (fh *FlagHandler) collectFlagState(cmd *cobra.Command) map[string]bool {
	flagState := make(map[string]bool)

	// Safely get flag values with error handling
	if verbose, err := cmd.Flags().GetBool("verbose"); err == nil {
		flagState["--verbose"] = verbose
	}
	if quiet, err := cmd.Flags().GetBool("quiet"); err == nil {
		flagState["--quiet"] = quiet
	}
	if debug, err := cmd.Flags().GetBool("debug"); err == nil {
		flagState["--debug"] = debug
	}
	if interactive, err := cmd.Flags().GetBool("interactive"); err == nil {
		flagState["--interactive"] = interactive
	}
	if nonInteractive, err := cmd.Flags().GetBool("non-interactive"); err == nil {
		flagState["--non-interactive"] = nonInteractive
	}
	if forceInteractive, err := cmd.Flags().GetBool("force-interactive"); err == nil {
		flagState["--force-interactive"] = forceInteractive
	}
	if forceNonInteractive, err := cmd.Flags().GetBool("force-non-interactive"); err == nil {
		flagState["--force-non-interactive"] = forceNonInteractive
	}
	if mode, err := cmd.Flags().GetString("mode"); err == nil && mode != "" {
		flagState["--mode"] = true
		flagState[fmt.Sprintf("--mode=%s", mode)] = true
	}
	if outputFormat, err := cmd.Flags().GetString("output-format"); err == nil && outputFormat != "" {
		flagState["--output-format"] = true
		flagState[fmt.Sprintf("--output-format=%s", outputFormat)] = true
	}

	return flagState
}

// ValidateLogLevel validates the provided log level.
func (fh *FlagHandler) ValidateLogLevel(logLevel string) error {
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	for _, level := range validLogLevels {
		if logLevel == level {
			return nil
		}
	}
	return fmt.Errorf("🚫 %s isn't a valid log level. %s: %s",
		fh.cli.Error(fmt.Sprintf("'%s'", logLevel)),
		fh.cli.Info("Available options"),
		fh.cli.Highlight(strings.Join(validLogLevels, ", ")))
}

// ValidateOutputFormat validates the provided output format.
func (fh *FlagHandler) ValidateOutputFormat(outputFormat string) error {
	validOutputFormats := []string{"text", "json", "yaml"}
	for _, format := range validOutputFormats {
		if outputFormat == format {
			return nil
		}
	}
	return fmt.Errorf("🚫 %s isn't a valid output format. %s: %s",
		fh.cli.Error(fmt.Sprintf("'%s'", outputFormat)),
		fh.cli.Info("Available options"),
		fh.cli.Highlight(strings.Join(validOutputFormats, ", ")))
}

// ValidateModeFlags checks for conflicting mode flags using the enhanced conflict detection system.
func (fh *FlagHandler) ValidateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	// Add null safety check
	if fh == nil {
		return fmt.Errorf("flag handler not initialized")
	}

	// Create flag state map for mode flags
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
		if err := fh.validateExplicitMode(explicitMode); err != nil {
			return err
		}
	}

	// Use enhanced conflict detection for mode flags
	return fh.ValidateConflictingFlagsEnhanced(flagState)
}

// ValidateModeFlagsWithGracefulFallback validates mode flags with graceful fallback handling.
func (fh *FlagHandler) ValidateModeFlagsWithGracefulFallback(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	// First try normal validation
	if err := fh.ValidateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode); err != nil {
		// Check if this is a recoverable conflict
		if fh.isRecoverableConflict(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode) {
			return fh.handleRecoverableConflict(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)
		}
		return err
	}
	return nil
}

// isRecoverableConflict checks if a mode flag conflict can be automatically resolved.
func (fh *FlagHandler) isRecoverableConflict(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) bool {
	// Count active flags
	activeCount := 0
	if nonInteractive {
		activeCount++
	}
	if interactive {
		activeCount++
	}
	if forceInteractive {
		activeCount++
	}
	if forceNonInteractive {
		activeCount++
	}
	if explicitMode != "" {
		activeCount++
	}

	// Only recoverable if exactly 2 flags are set and one is a "force" flag
	if activeCount == 2 {
		return forceInteractive || forceNonInteractive
	}

	return false
}

// handleRecoverableConflict attempts to resolve a recoverable mode flag conflict.
func (fh *FlagHandler) handleRecoverableConflict(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) error {
	var resolvedMode string
	var overriddenFlag string

	// Force flags take precedence
	if forceInteractive {
		resolvedMode = "interactive"
		if nonInteractive {
			overriddenFlag = "--non-interactive"
		} else if explicitMode != "" && explicitMode != "interactive" {
			overriddenFlag = fmt.Sprintf("--mode=%s", explicitMode)
		}
	} else if forceNonInteractive {
		resolvedMode = "non-interactive"
		if interactive {
			overriddenFlag = "--interactive"
		} else if explicitMode != "" && explicitMode != "non-interactive" {
			overriddenFlag = fmt.Sprintf("--mode=%s", explicitMode)
		}
	}

	// Log the resolution if we have a CLI instance
	if fh.cli != nil && resolvedMode != "" && overriddenFlag != "" {
		fh.cli.VerboseOutput("⚠️  Resolved mode conflict: Using %s mode (overriding %s)", resolvedMode, overriddenFlag)
	}

	return nil
}

// validateExplicitMode validates the explicit mode value
func (fh *FlagHandler) validateExplicitMode(explicitMode string) error {
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

	// Mode is invalid, provide helpful error
	var errorMsg string
	if fh.cli != nil {
		errorMsg = fmt.Sprintf("🚫 %s %s\n%s: %s\n%s: %s",
			fh.cli.Error(fmt.Sprintf("'%s' is not a valid mode", explicitMode)),
			fh.cli.Info("Invalid mode specified"),
			fh.cli.Info("Available modes"),
			fh.cli.Highlight(strings.Join(validModes, ", ")),
			fh.cli.Info("Example"),
			fh.cli.Dim("--mode=interactive"))
	} else {
		errorMsg = fmt.Sprintf("'%s' is not a valid mode. Available modes: %s", explicitMode, strings.Join(validModes, ", "))
	}
	return fmt.Errorf("%s", errorMsg)
}

// IsNonInteractiveMode checks if the CLI should run in non-interactive mode.
func (fh *FlagHandler) IsNonInteractiveMode(cmd *cobra.Command) bool {
	// Add null safety check
	if fh == nil {
		return false
	}

	// Check explicit flag
	if cmd != nil {
		if nonInteractive, err := cmd.PersistentFlags().GetBool("non-interactive"); err == nil && nonInteractive {
			return true
		}
	}

	// Check environment variable
	if fh.parseBoolEnv("GENERATOR_NON_INTERACTIVE", false) {
		return true
	}

	// Check CI environment (with null safety)
	if fh.cli != nil {
		ci := fh.cli.detectCIEnvironment()
		if ci.IsCI {
			return true
		}
	}

	// Check if stdin is not a terminal (piped input)
	if !fh.isTerminal() {
		return true
	}

	return false
}

// configureLogger sets up the logger based on flag values.
func (fh *FlagHandler) configureLogger(logLevel string, logJSON, logCaller bool, cmd *cobra.Command, verbose, quiet, debug, nonInteractive bool, outputFormat string) {
	// Set log level
	switch logLevel {
	case "debug":
		fh.logger.SetLevel(0) // LogLevelDebug
	case "info":
		fh.logger.SetLevel(1) // LogLevelInfo
	case "warn":
		fh.logger.SetLevel(2) // LogLevelWarn
	case "error":
		fh.logger.SetLevel(3) // LogLevelError
	case "fatal":
		fh.logger.SetLevel(4) // LogLevelFatal
	}

	// Configure JSON output
	fh.logger.SetJSONOutput(logJSON)

	// Configure caller information
	fh.logger.SetCallerInfo(logCaller)

	// Log configuration changes in debug mode
	if debug {
		fh.logger.DebugWithFields("🔧 Setting up CLI configuration", map[string]interface{}{
			"command":         cmd.Name(),
			"verbose":         verbose,
			"quiet":           quiet,
			"debug":           debug,
			"log_level":       logLevel,
			"log_json":        logJSON,
			"log_caller":      logCaller,
			"non_interactive": nonInteractive,
			"output_format":   outputFormat,
		})
	}

	// Log operation start for verbose mode
	if verbose || debug {
		fh.logger.InfoWithFields("🚀 Starting your command", map[string]interface{}{
			"command": cmd.Name(),
			"args":    cmd.Flags().Args(),
		})
	}
}

// displayModeMessages shows appropriate messages based on the current mode.
func (fh *FlagHandler) displayModeMessages(cmd *cobra.Command, verbose, quiet, debug, nonInteractive bool) {
	if verbose && !quiet {
		fmt.Printf("🔍 Running '%s' with detailed output\n", cmd.Name())
	}

	if debug && !quiet {
		fmt.Printf("🐛 Running '%s' with debug logging and performance metrics\n", cmd.Name())
	}

	if nonInteractive && (verbose || debug) && !quiet {
		fmt.Printf("🤖 Running in non-interactive mode\n")
	}
}

// Helper methods for environment detection

// parseBoolEnv parses a boolean environment variable with a default value.
func (fh *FlagHandler) parseBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// isTerminal checks if the current process is running in a terminal.
func (fh *FlagHandler) isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// GetModeFromFlags safely extracts mode information from command flags
func (fh *FlagHandler) GetModeFromFlags(cmd *cobra.Command) (GenerationMode, error) {
	// Add null safety checks
	if fh == nil {
		return ModeAuto, fmt.Errorf("flag handler not initialized")
	}
	if cmd == nil {
		return ModeAuto, fmt.Errorf("command not provided")
	}

	// Get mode flags with error handling
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	interactive, _ := cmd.Flags().GetBool("interactive")
	forceInteractive, _ := cmd.Flags().GetBool("force-interactive")
	forceNonInteractive, _ := cmd.Flags().GetBool("force-non-interactive")
	explicitMode, _ := cmd.Flags().GetString("mode")

	// Validate mode flags first
	if err := fh.ValidateModeFlags(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode); err != nil {
		return ModeAuto, err
	}

	// Determine mode based on flags
	if explicitMode != "" {
		return fh.parseExplicitMode(explicitMode), nil
	}
	if forceNonInteractive || nonInteractive {
		return ModeNonInteractive, nil
	}
	if forceInteractive || interactive {
		return ModeInteractive, nil
	}

	// Auto-detect mode
	if fh.IsNonInteractiveMode(cmd) {
		return ModeNonInteractive, nil
	}

	return ModeInteractive, nil
}

// parseExplicitMode converts explicit mode string to GenerationMode
func (fh *FlagHandler) parseExplicitMode(mode string) GenerationMode {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "non-interactive", "noninteractive", "ni", "auto":
		return ModeNonInteractive
	case "interactive", "i":
		return ModeInteractive
	case "config-file", "config", "file", "cf":
		return ModeConfig
	default:
		return ModeAuto
	}
}

// GenerationMode represents the different generation modes
type GenerationMode int

const (
	ModeAuto GenerationMode = iota
	ModeInteractive
	ModeNonInteractive
	ModeConfig
)

// ConflictRule represents a flag conflict rule with detailed information
type ConflictRule struct {
	Flags       []string `json:"flags"`
	Description string   `json:"description"`
	Suggestion  string   `json:"suggestion"`
	Examples    []string `json:"examples"`
	Severity    string   `json:"severity"` // "error", "warning", "info"
}

// FlagConflictMatrix defines all possible flag conflicts and their rules
type FlagConflictMatrix struct {
	Rules []ConflictRule `json:"rules"`
}

// NewFlagConflictMatrix creates a comprehensive conflict rule matrix
func NewFlagConflictMatrix() *FlagConflictMatrix {
	return &FlagConflictMatrix{
		Rules: []ConflictRule{
			// Output mode conflicts
			{
				Flags:       []string{"--verbose", "--quiet"},
				Description: "Verbose and quiet modes are mutually exclusive",
				Suggestion:  "Choose either verbose output for detailed information OR quiet mode for minimal output",
				Examples:    []string{"--verbose", "--quiet", "--debug (implies verbose)"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--debug", "--quiet"},
				Description: "Debug and quiet modes are mutually exclusive",
				Suggestion:  "Choose either debug mode for detailed debugging OR quiet mode for minimal output",
				Examples:    []string{"--debug", "--quiet"},
				Severity:    "error",
			},
			// Generation mode conflicts
			{
				Flags:       []string{"--interactive", "--non-interactive"},
				Description: "Interactive and non-interactive modes cannot be used together",
				Suggestion:  "Choose either interactive mode for guided setup OR non-interactive for automated generation",
				Examples:    []string{"--interactive", "--non-interactive", "--mode=interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--force-interactive", "--force-non-interactive"},
				Description: "Force interactive and force non-interactive modes are mutually exclusive",
				Suggestion:  "Choose either force-interactive to override detection OR force-non-interactive for automation",
				Examples:    []string{"--force-interactive", "--force-non-interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--interactive", "--force-non-interactive"},
				Description: "Interactive mode conflicts with forced non-interactive mode",
				Suggestion:  "Use either --interactive for guided setup OR --force-non-interactive for automation",
				Examples:    []string{"--interactive", "--force-non-interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--non-interactive", "--force-interactive"},
				Description: "Non-interactive mode conflicts with forced interactive mode",
				Suggestion:  "Use either --non-interactive for automation OR --force-interactive for guided setup",
				Examples:    []string{"--non-interactive", "--force-interactive"},
				Severity:    "error",
			},
			// Mode flag with explicit mode conflicts
			{
				Flags:       []string{"--interactive", "--mode"},
				Description: "Interactive flag conflicts with explicit mode specification",
				Suggestion:  "Use either --interactive flag OR --mode=interactive, not both",
				Examples:    []string{"--interactive", "--mode=interactive", "--mode=non-interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--non-interactive", "--mode"},
				Description: "Non-interactive flag conflicts with explicit mode specification",
				Suggestion:  "Use either --non-interactive flag OR --mode=non-interactive, not both",
				Examples:    []string{"--non-interactive", "--mode=non-interactive", "--mode=interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--force-interactive", "--mode"},
				Description: "Force-interactive flag conflicts with explicit mode specification",
				Suggestion:  "Use either --force-interactive flag OR --mode=interactive, not both",
				Examples:    []string{"--force-interactive", "--mode=interactive"},
				Severity:    "error",
			},
			{
				Flags:       []string{"--force-non-interactive", "--mode"},
				Description: "Force-non-interactive flag conflicts with explicit mode specification",
				Suggestion:  "Use either --force-non-interactive flag OR --mode=non-interactive, not both",
				Examples:    []string{"--force-non-interactive", "--mode=non-interactive"},
				Severity:    "error",
			},
			// Output format conflicts (potential future conflicts)
			{
				Flags:       []string{"--output-format=json", "--quiet"},
				Description: "JSON output format may conflict with quiet mode for readability",
				Suggestion:  "Consider using JSON format without quiet mode for better structured output",
				Examples:    []string{"--output-format=json", "--output-format=yaml --verbose"},
				Severity:    "warning",
			},
		},
	}
}
