// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
//
// This file contains CLI-specific validation logic including configuration validation,
// pre-generation checks, and option validation.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// CLIValidator handles CLI-specific validation operations.
//
// The CLIValidator provides methods for:
//   - Project configuration validation
//   - Pre-generation checks and directory validation
//   - Command option validation
//   - File system permission checks
//   - Input format validation

// Pre-compiled regular expressions for CLI validation
var (
	projectNameRegex  = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	templateNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

type CLIValidator struct {
	cli           *CLI
	outputManager *OutputManager
	logger        interfaces.Logger
}

// NewCLIValidator creates a new CLI validator instance.
//
// Parameters:
//   - cli: The main CLI instance for accessing dependencies
//   - outputManager: For formatted output and messaging
//   - logger: For logging validation operations
//
// Returns:
//   - *CLIValidator: New validator instance ready for use
func NewCLIValidator(cli *CLI, outputManager *OutputManager, logger interfaces.Logger) *CLIValidator {
	return &CLIValidator{
		cli:           cli,
		outputManager: outputManager,
		logger:        logger,
	}
}

// ValidateGenerateConfiguration validates the configuration for project generation.
//
// This method performs comprehensive validation of project configuration including:
//   - Project name validation (required and format)
//   - License validation (if specified)
//   - Basic configuration completeness checks
//
// Parameters:
//   - config: The project configuration to validate
//   - options: Generation options that may affect validation
//
// Returns:
//   - error: Validation error with suggestions, or nil if valid
func (v *CLIValidator) ValidateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	if config.Name == "" {
		err := v.cli.createConfigurationError("🚫 Your project needs a name", "")
		err = err.WithSuggestions("Set project name in configuration file or GENERATOR_PROJECT_NAME environment variable")
		return err
	}

	// Validate project name format
	if !v.IsValidProjectName(config.Name) {
		err := v.cli.createConfigurationError(fmt.Sprintf("🚫 '%s' isn't a valid project name", config.Name), "")
		err = err.WithSuggestions("Project names can only contain letters, numbers, hyphens, and underscores")
		return err
	}

	// Validate license if specified
	if config.License != "" && !v.IsValidLicense(config.License) {
		err := v.cli.createConfigurationError(fmt.Sprintf("🚫 '%s' isn't a valid license", config.License), "")
		err = err.WithSuggestions("Use a valid SPDX license identifier like MIT, Apache-2.0, or GPL-3.0")
		return err
	}

	return nil
}

// PerformPreGenerationChecks performs comprehensive checks before project generation.
//
// This method handles:
//   - Output directory existence and overwrite confirmation
//   - Directory creation and cleanup
//   - Write permission validation
//   - User interaction for confirmations (in interactive mode)
//
// Parameters:
//   - outputPath: The target directory for project generation
//   - options: Generation options including force and non-interactive flags
//
// Returns:
//   - error: Pre-generation check error, or nil if all checks pass
func (v *CLIValidator) PerformPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	// Check if output directory exists
	if _, err := os.Stat(outputPath); err == nil {
		if !options.Force && !options.NonInteractive {
			// Ask user for confirmation
			fmt.Printf("\n⚠️  Directory %s already exists.\n", v.cli.Highlight(fmt.Sprintf("'%s'", outputPath)))
			fmt.Printf("Do you want to overwrite it? %s: ", v.cli.Dim("(y/N)"))

			var response string
			_, _ = fmt.Scanln(&response)
			response = strings.ToLower(strings.TrimSpace(response))

			if response != "y" && response != "yes" {
				return fmt.Errorf("🚫 %s %s",
					v.cli.Error("Project generation cancelled by user."),
					v.cli.Info("Run again with --force to automatically overwrite existing directories"))
			}
		}

		v.cli.VerboseOutput("🗑️  Removing existing directory: %s", v.cli.Highlight(outputPath))
		if err := os.RemoveAll(outputPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				v.cli.Error("Unable to remove existing directory."),
				v.cli.Info("Check directory permissions and ensure no files are in use"))
		}
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, 0750); err != nil {
			return fmt.Errorf("🚫 %s %s",
				v.cli.Error("Unable to create output directory."),
				v.cli.Info("Check parent directory permissions and available disk space"))
		}
		v.cli.VerboseOutput("📁 Created output directory: %s", v.cli.Info(outputPath))
	}

	// Check write permissions on the output directory
	if err := v.CheckWritePermissions(outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			v.cli.Error("No write permission for output directory."),
			v.cli.Info("Check directory permissions or run with appropriate privileges"))
	}

	return nil
}

// ValidateGenerateOptions validates the generate command options.
//
// This method performs comprehensive validation of generation options including:
//   - Output path validation and normalization
//   - Template name format validation
//   - Conflicting option detection
//   - Warning generation for potentially problematic combinations
//
// Parameters:
//   - options: The generation options to validate
//
// Returns:
//   - error: Validation error with details, or nil if valid
func (v *CLIValidator) ValidateGenerateOptions(options interfaces.GenerateOptions) error {
	var validationErrors []string

	// Validate output path
	if options.OutputPath != "" {
		if !filepath.IsAbs(options.OutputPath) && !strings.HasPrefix(options.OutputPath, "./") && !strings.HasPrefix(options.OutputPath, "../") {
			// Relative path without ./ prefix - this is okay, but we'll make it explicit
			options.OutputPath = "./" + options.OutputPath
		}

		// Check for invalid characters in path
		if strings.ContainsAny(options.OutputPath, "<>:\"|?*") {
			validationErrors = append(validationErrors, "output path contains invalid characters")
		}
	}

	// Validate template names
	for _, template := range options.Templates {
		if template == "" {
			validationErrors = append(validationErrors, "empty template name specified")
			continue
		}

		// Validate template name format
		if !v.IsValidTemplateName(template) {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid template name '%s' - must contain only letters, numbers, hyphens, and underscores", template))
		}
	}

	// Validate conflicting options
	if options.Offline && options.UpdateVersions {
		validationErrors = append(validationErrors, "cannot use --offline and --update-versions together")
	}

	if options.Minimal && options.IncludeExamples {
		v.cli.WarningOutput("Using --minimal with --include-examples may result in minimal examples only")
	}

	// Validate dry-run with force
	if options.DryRun && options.Force {
		v.cli.WarningOutput("--force flag has no effect in dry-run mode")
	}

	if len(validationErrors) > 0 {
		return &interfaces.CLIError{
			Type:        interfaces.ErrorTypeValidation,
			Message:     "generate options validation failed",
			Code:        interfaces.ErrorCodeValidationFailed,
			Details:     map[string]any{"errors": validationErrors},
			Suggestions: []string{"Fix the validation errors and try again"},
		}
	}

	return nil
}

// CheckWritePermissions checks if we have write permissions to a directory.
//
// This method tests write permissions by attempting to create and remove
// a temporary file in the target directory.
//
// Parameters:
//   - path: The directory path to check for write permissions
//
// Returns:
//   - error: Permission error if write access is denied, or nil if accessible
func (v *CLIValidator) CheckWritePermissions(path string) error {
	// Try to create a temporary file to test permissions
	// Use a secure temporary file name with random suffix
	tempFile := filepath.Join(path, ".generator-permission-test-"+strconv.FormatInt(time.Now().UnixNano(), 36))

	// #nosec G304 - This is a controlled temporary file creation for permission testing
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("no write permission to directory %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		v.cli.WarningOutput("📄 Couldn't close temporary file: %v", err)
	}
	if err := os.Remove(tempFile); err != nil {
		v.cli.WarningOutput("🗑️  Couldn't remove temporary file: %v", err)
	}
	return nil
}

// IsValidProjectName validates project name format.
//
// Project names must:
//   - Contain only letters, numbers, hyphens, and underscores
//   - Be between 1 and 100 characters long
//   - Match the pattern ^[a-zA-Z0-9_-]+$
//
// Parameters:
//   - name: The project name to validate
//
// Returns:
//   - bool: true if the name is valid, false otherwise
func (v *CLIValidator) IsValidProjectName(name string) bool {
	// Allow letters, numbers, hyphens, and underscores using pre-compiled regex
	return projectNameRegex.MatchString(name) && len(name) > 0 && len(name) <= 100
}

// IsValidLicense validates license identifier against common SPDX identifiers.
//
// Supported licenses include:
//   - MIT, Apache-2.0, GPL-2.0, GPL-3.0
//   - LGPL-2.1, LGPL-3.0, BSD-2-Clause, BSD-3-Clause
//   - ISC, MPL-2.0, UNLICENSED
//
// Parameters:
//   - license: The license identifier to validate
//
// Returns:
//   - bool: true if the license is a valid SPDX identifier, false otherwise
func (v *CLIValidator) IsValidLicense(license string) bool {
	// Common SPDX license identifiers
	validLicenses := []string{
		"MIT", "Apache-2.0", "GPL-2.0", "GPL-3.0", "LGPL-2.1", "LGPL-3.0",
		"BSD-2-Clause", "BSD-3-Clause", "ISC", "MPL-2.0", "UNLICENSED",
	}

	for _, valid := range validLicenses {
		if strings.EqualFold(license, valid) {
			return true
		}
	}
	return false
}

// IsValidTemplateName validates template name format.
//
// Template names must:
//   - Contain only letters, numbers, hyphens, underscores, and dots
//   - Be between 1 and 50 characters long
//   - Match the pattern ^[a-zA-Z0-9_.-]+$
//
// Parameters:
//   - name: The template name to validate
//
// Returns:
//   - bool: true if the name is valid, false otherwise
func (v *CLIValidator) IsValidTemplateName(name string) bool {
	// Allow letters, numbers, hyphens, underscores, and dots using pre-compiled regex
	return templateNameRegex.MatchString(name) && len(name) > 0 && len(name) <= 50
}
