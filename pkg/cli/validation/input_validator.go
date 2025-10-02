// Package validation provides CLI-specific validation functionality.
//
// This module contains input validation logic for CLI operations including
// configuration validation, option validation, and input format validation.
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Pre-compiled regular expressions for input validation
var (
	validProjectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`)
)

// InputValidator handles CLI-specific input validation operations.
//
// The InputValidator provides methods for:
//   - Project configuration validation
//   - Command option validation
//   - Input format validation
//   - File system validation
type InputValidator struct {
	logger interfaces.Logger
	output OutputInterface
}

// OutputInterface defines the output methods needed by the validator
type OutputInterface interface {
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	WarningOutput(format string, args ...interface{})
	Error(text string) string
	Info(text string) string
	Warning(text string) string
}

// NewInputValidator creates a new input validator instance.
func NewInputValidator(logger interfaces.Logger, output OutputInterface) *InputValidator {
	return &InputValidator{
		logger: logger,
		output: output,
	}
}

// ValidateGenerateConfiguration validates the configuration for generation
func (iv *InputValidator) ValidateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	if config == nil {
		return fmt.Errorf("project configuration is required")
	}

	// Validate project name
	if err := iv.validateProjectName(config.Name); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Validate license
	if err := iv.validateLicense(config.License); err != nil {
		return fmt.Errorf("invalid license: %w", err)
	}

	// Validate organization if provided
	if config.Organization != "" {
		if err := iv.validateOrganization(config.Organization); err != nil {
			return fmt.Errorf("invalid organization: %w", err)
		}
	}

	// Validate component selections
	if err := iv.validateComponentSelections(config); err != nil {
		return fmt.Errorf("invalid component selection: %w", err)
	}

	// Validate output path if provided
	if config.OutputPath != "" {
		if err := iv.validateOutputPath(config.OutputPath); err != nil {
			return fmt.Errorf("invalid output path: %w", err)
		}
	}

	iv.output.DebugOutput("Configuration validation passed for project: %s", config.Name)
	return nil
}

// ValidateGenerateOptions validates the generate command options
func (iv *InputValidator) ValidateGenerateOptions(options interfaces.GenerateOptions) error {
	// Validate output path if provided
	if options.OutputPath != "" {
		if err := iv.validateOutputPath(options.OutputPath); err != nil {
			return fmt.Errorf("invalid output path in options: %w", err)
		}
	}

	// Note: Interactive mode is the default, NonInteractive is the flag
	// No conflicting options to validate for this field

	if options.Offline && options.UpdateVersions {
		return fmt.Errorf("cannot update versions in offline mode")
	}

	if options.DryRun && options.Force {
		iv.output.WarningOutput("⚠️  --force flag has no effect in dry run mode")
	}

	iv.output.DebugOutput("Generate options validation passed")
	return nil
}

// validateProjectName validates project name format
func (iv *InputValidator) validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("project name cannot exceed 100 characters")
	}

	// Project names should contain only letters, numbers, hyphens, and underscores
	// Should not start or end with hyphens or underscores using pre-compiled regex
	if !validProjectNameRegex.MatchString(name) {
		return fmt.Errorf("project name must contain only letters, numbers, hyphens, and underscores, and cannot start or end with hyphens or underscores")
	}

	// Check for reserved names
	reservedNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	lowerName := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			return fmt.Errorf("project name '%s' is reserved and cannot be used", name)
		}
	}

	return nil
}

// validateLicense validates license identifier
func (iv *InputValidator) validateLicense(license string) error {
	if license == "" {
		return fmt.Errorf("license cannot be empty")
	}

	// Common SPDX license identifiers
	commonLicenses := map[string]bool{
		"MIT":          true,
		"Apache-2.0":   true,
		"GPL-3.0":      true,
		"GPL-2.0":      true,
		"BSD-3-Clause": true,
		"BSD-2-Clause": true,
		"ISC":          true,
		"MPL-2.0":      true,
		"LGPL-3.0":     true,
		"LGPL-2.1":     true,
		"Unlicense":    true,
		"CC0-1.0":      true,
		"AGPL-3.0":     true,
		"Proprietary":  true,
	}

	if !commonLicenses[license] {
		iv.output.WarningOutput("⚠️  License '%s' is not a common SPDX identifier", license)
		iv.output.VerboseOutput("Common licenses: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, ISC, MPL-2.0")
	}

	return nil
}

// validateOrganization validates organization name
func (iv *InputValidator) validateOrganization(organization string) error {
	if len(organization) > 200 {
		return fmt.Errorf("organization name cannot exceed 200 characters")
	}

	// Organization names should not contain control characters
	if strings.ContainsAny(organization, "\n\r\t\v\f") {
		return fmt.Errorf("organization name cannot contain control characters")
	}

	return nil
}

// validateComponentSelections validates component selections
func (iv *InputValidator) validateComponentSelections(config *models.ProjectConfig) error {
	hasAnyComponent := config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared ||
		config.Components.Backend.GoGin ||
		config.Components.Mobile.Android ||
		config.Components.Mobile.IOS ||
		config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform

	if !hasAnyComponent {
		return fmt.Errorf("at least one component must be selected")
	}

	return nil
}

// validateOutputPath validates output path
func (iv *InputValidator) validateOutputPath(outputPath string) error {
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Check if path is absolute or relative
	if filepath.IsAbs(outputPath) {
		// For absolute paths, check if parent directory exists
		parentDir := filepath.Dir(outputPath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			return fmt.Errorf("parent directory does not exist: %s", parentDir)
		}
	}

	// Check for invalid characters in path
	if strings.ContainsAny(outputPath, "<>:\"|?*") {
		return fmt.Errorf("output path contains invalid characters")
	}

	return nil
}

// PerformPreGenerationChecks performs checks before generation
func (iv *InputValidator) PerformPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	iv.output.VerboseOutput("🔍 Performing pre-generation checks...")

	// Check if output directory exists
	if _, err := os.Stat(outputPath); err == nil {
		if !options.Force {
			return fmt.Errorf("🚫 %s %s %s",
				iv.output.Error("Output directory already exists:"),
				iv.output.Info(fmt.Sprintf("'%s'.", outputPath)),
				iv.output.Info("Use --force to overwrite or choose a different location"))
		}
		iv.output.WarningOutput("⚠️  Output directory exists and will be overwritten")
	}

	// Check write permissions for parent directory
	parentDir := filepath.Dir(outputPath)
	if err := iv.checkWritePermissions(parentDir); err != nil {
		return fmt.Errorf("🚫 %s %s %s",
			iv.output.Error("Cannot write to output directory:"),
			iv.output.Info(fmt.Sprintf("'%s'.", parentDir)),
			iv.output.Info("Check directory permissions"))
	}

	// Check available disk space (basic check)
	if err := iv.checkDiskSpace(parentDir); err != nil {
		return fmt.Errorf("🚫 %s %s",
			iv.output.Error("Insufficient disk space:"),
			iv.output.Info("Ensure at least 100MB of free space is available"))
	}

	iv.output.VerboseOutput("✅ Pre-generation checks passed")
	return nil
}

// CheckWritePermissions checks if we have write permissions to a directory
func (iv *InputValidator) CheckWritePermissions(path string) error {
	return iv.checkWritePermissions(path)
}

// checkWritePermissions checks if we have write permissions to a directory
func (iv *InputValidator) checkWritePermissions(path string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	// Try to create a temporary file to test write permissions
	tempFile := filepath.Join(path, ".generator_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("no write permission: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close test file: %w", err)
	}

	// Clean up the temporary file
	if err := os.Remove(tempFile); err != nil {
		iv.output.VerboseOutput("Warning: could not clean up temporary file: %v", err)
	}

	return nil
}

// checkDiskSpace performs a basic disk space check
func (iv *InputValidator) checkDiskSpace(path string) error {
	// This is a simplified check - in a real implementation,
	// you would use platform-specific APIs to check actual disk space

	// For now, just check if we can create a small test file
	testFile := filepath.Join(path, ".generator_space_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot create test file: %w", err)
	}

	// Write a small amount of data
	testData := make([]byte, 1024) // 1KB
	_, err = file.Write(testData)
	if closeErr := file.Close(); closeErr != nil {
		return fmt.Errorf("failed to close test file: %w", closeErr)
	}

	if err != nil {
		if removeErr := os.Remove(testFile); removeErr != nil {
			return fmt.Errorf("write failed and cleanup failed: %w (cleanup error: %v)", err, removeErr)
		}
		return fmt.Errorf("cannot write test data: %w", err)
	}

	// Clean up
	if err := os.Remove(testFile); err != nil {
		iv.output.VerboseOutput("Warning: could not clean up test file: %v", err)
	}

	return nil
}

// IsValidProjectName validates project name format (public method)
func (iv *InputValidator) IsValidProjectName(name string) bool {
	return iv.validateProjectName(name) == nil
}

// IsValidLicense validates license identifier (public method)
func (iv *InputValidator) IsValidLicense(license string) bool {
	return iv.validateLicense(license) == nil
}

// ValidateConfigurationFile validates a configuration file
func (iv *InputValidator) ValidateConfigurationFile(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("configuration file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("🚫 %s %s",
			iv.output.Error("Configuration file not found:"),
			iv.output.Info(fmt.Sprintf("'%s'", configPath)))
	}

	// Check if file is readable
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			iv.output.Error("Unable to read configuration file."),
			iv.output.Info("Check file permissions and ensure it's not corrupted"))
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close configuration file: %w", err)
	}

	iv.output.DebugOutput("✅ Configuration file looks good: %s", configPath)
	return nil
}
