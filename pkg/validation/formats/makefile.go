package formats

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// MakefileValidator provides specialized Makefile validation
type MakefileValidator struct {
	commonTargets []string
	shellCommands map[string]string
}

// MakefileTarget represents a parsed Makefile target
type MakefileTarget struct {
	Name         string
	Dependencies []string
	Commands     []string
	LineNum      int
	IsPhony      bool
}

// MakefileVariable represents a Makefile variable
type MakefileVariable struct {
	Name    string
	Value   string
	LineNum int
	Type    string // =, :=, +=, ?=
}

// NewMakefileValidator creates a new Makefile validator
func NewMakefileValidator() *MakefileValidator {
	validator := &MakefileValidator{
		commonTargets: []string{"all", "build", "test", "clean", "install", "help"},
		shellCommands: make(map[string]string),
	}
	validator.initializeShellCommands()
	return validator
}

// ValidateMakefile validates Makefile configuration
func (mv *MakefileValidator) ValidateMakefile(filePath string) (*interfaces.ConfigValidationResult, error) {
	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "file",
			Value:    filePath,
			Type:     "read_error",
			Message:  fmt.Sprintf("Failed to read file: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.file_access",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	targets, variables := mv.parseMakefile(string(content))

	// Validate overall structure
	mv.validateMakefileStructure(targets, variables, result)

	// Validate individual targets
	for _, target := range targets {
		mv.validateTarget(target, result)
	}

	// Validate variables
	for _, variable := range variables {
		mv.validateVariable(variable, result)
	}

	// Validate syntax and best practices
	mv.validateSyntaxAndBestPractices(string(content), result)

	return result, nil
}

// parseMakefile parses Makefile content into targets and variables
func (mv *MakefileValidator) parseMakefile(content string) ([]MakefileTarget, []MakefileVariable) {
	var targets []MakefileTarget
	var variables []MakefileVariable

	lines := strings.Split(content, "\n")
	var currentTarget *MakefileTarget

	for lineNum, line := range lines {
		line = strings.TrimRight(line, " \t")

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Check for target definitions (lines ending with :)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "\t") {
			// Save previous target
			if currentTarget != nil {
				targets = append(targets, *currentTarget)
			}

			// Parse new target
			parts := strings.SplitN(line, ":", 2)
			targetName := strings.TrimSpace(parts[0])
			dependencies := []string{}

			if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
				depStr := strings.TrimSpace(parts[1])
				dependencies = strings.Fields(depStr)
			}

			currentTarget = &MakefileTarget{
				Name:         targetName,
				Dependencies: dependencies,
				Commands:     []string{},
				LineNum:      lineNum + 1,
				IsPhony:      false,
			}
			continue
		}

		// Check for command lines (should start with tab)
		if strings.HasPrefix(line, "\t") {
			if currentTarget != nil {
				command := strings.TrimPrefix(line, "\t")
				currentTarget.Commands = append(currentTarget.Commands, command)
			}
			continue
		}

		// Check for variable assignments
		if mv.isVariableAssignment(line) {
			variable := mv.parseVariable(line, lineNum+1)
			if variable != nil {
				variables = append(variables, *variable)
			}
			continue
		}

		// Check for .PHONY declarations
		if strings.HasPrefix(strings.TrimSpace(line), ".PHONY:") {
			phonyTargets := strings.Fields(strings.TrimPrefix(strings.TrimSpace(line), ".PHONY:"))
			// Mark targets as phony (would need to track this)
			_ = phonyTargets
			continue
		}
	}

	// Add the last target
	if currentTarget != nil {
		targets = append(targets, *currentTarget)
	}

	return targets, variables
}

// isVariableAssignment checks if a line is a variable assignment
func (mv *MakefileValidator) isVariableAssignment(line string) bool {
	// Check for various assignment operators
	assignmentOps := []string{"=", ":=", "+=", "?="}

	for _, op := range assignmentOps {
		if strings.Contains(line, op) && !strings.HasPrefix(strings.TrimSpace(line), "\t") {
			return true
		}
	}

	return false
}

// parseVariable parses a variable assignment line
func (mv *MakefileValidator) parseVariable(line string, lineNum int) *MakefileVariable {
	// Find the assignment operator
	assignmentOps := []string{":=", "+=", "?=", "="}

	for _, op := range assignmentOps {
		if strings.Contains(line, op) {
			parts := strings.SplitN(line, op, 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				return &MakefileVariable{
					Name:    name,
					Value:   value,
					LineNum: lineNum,
					Type:    op,
				}
			}
		}
	}

	return nil
}

// validateMakefileStructure validates the overall structure of the Makefile
func (mv *MakefileValidator) validateMakefileStructure(targets []MakefileTarget, variables []MakefileVariable, result *interfaces.ConfigValidationResult) {
	result.Summary.TotalProperties = len(targets) + len(variables)
	result.Summary.ValidProperties = len(targets) + len(variables)

	// Check for common targets
	targetNames := make(map[string]bool)
	for _, target := range targets {
		targetNames[target.Name] = true
	}

	// Check for recommended targets
	recommendedTargets := []string{"all", "clean", "help"}
	for _, recommended := range recommendedTargets {
		if !targetNames[recommended] {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "targets",
				Value:      recommended,
				Type:       "missing_recommended",
				Message:    fmt.Sprintf("Consider adding '%s' target", recommended),
				Suggestion: fmt.Sprintf("Add a '%s' target for better usability", recommended),
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "makefile.recommended_targets",
			})
			result.Summary.WarningCount++
		}
	}

	// Check if there's a default target (first non-special target)
	hasDefaultTarget := false
	for _, target := range targets {
		if !strings.HasPrefix(target.Name, ".") {
			hasDefaultTarget = true
			break
		}
	}

	if !hasDefaultTarget && len(targets) > 0 {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "structure",
			Value:      "no default target",
			Type:       "missing_recommended",
			Message:    "No default target found",
			Suggestion: "Add a default target (usually 'all') as the first target",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "makefile.default_target",
		})
		result.Summary.WarningCount++
	}
}

// validateTarget validates individual Makefile targets
func (mv *MakefileValidator) validateTarget(target MakefileTarget, result *interfaces.ConfigValidationResult) {
	// Check for empty targets
	if len(target.Commands) == 0 && len(target.Dependencies) == 0 {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("target_%s", target.Name),
			Value:      target.Name,
			Type:       "empty_target",
			Message:    "Target has no commands or dependencies",
			Suggestion: "Add commands or dependencies, or remove unused target",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "makefile.empty_target",
		})
		result.Summary.WarningCount++
	}

	// Validate target name
	if err := mv.validateTargetName(target.Name); err != nil {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", target.LineNum),
			Value:      target.Name,
			Type:       "naming_convention",
			Message:    fmt.Sprintf("Target name issue: %v", err),
			Suggestion: "Use lowercase letters, numbers, hyphens, and underscores",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "makefile.target_naming",
		})
		result.Summary.WarningCount++
	}

	// Validate commands
	for i, command := range target.Commands {
		mv.validateCommand(target.Name, command, i, result)
	}

	// Check for phony targets that should be declared
	if mv.shouldBePhony(target) && !target.IsPhony {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("target_%s", target.Name),
			Value:      target.Name,
			Type:       "missing_phony",
			Message:    "Target should be declared as .PHONY",
			Suggestion: fmt.Sprintf("Add '.PHONY: %s' to prevent conflicts with files", target.Name),
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "makefile.phony_declaration",
		})
		result.Summary.WarningCount++
	}
}

// validateVariable validates Makefile variables
func (mv *MakefileValidator) validateVariable(variable MakefileVariable, result *interfaces.ConfigValidationResult) {
	// Check variable naming convention
	if err := mv.validateVariableName(variable.Name); err != nil {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", variable.LineNum),
			Value:      variable.Name,
			Type:       "naming_convention",
			Message:    fmt.Sprintf("Variable name issue: %v", err),
			Suggestion: "Use uppercase letters, numbers, and underscores for variables",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "makefile.variable_naming",
		})
		result.Summary.WarningCount++
	}

	// Check for empty variables
	if strings.TrimSpace(variable.Value) == "" {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", variable.LineNum),
			Value:      variable.Name,
			Type:       "empty_variable",
			Message:    "Variable has empty value",
			Suggestion: "Provide a default value or remove unused variable",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "makefile.empty_variable",
		})
		result.Summary.WarningCount++
	}

	// Check assignment operator usage
	mv.validateAssignmentOperator(variable, result)

	// Check for hardcoded paths
	if mv.containsHardcodedPaths(variable.Value) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", variable.LineNum),
			Value:      variable.Value,
			Type:       "hardcoded_path",
			Message:    "Variable contains hardcoded paths",
			Suggestion: "Use relative paths or other variables for better portability",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "makefile.hardcoded_paths",
		})
		result.Summary.WarningCount++
	}
}

// validateCommand validates individual commands in targets
func (mv *MakefileValidator) validateCommand(targetName, command string, commandIndex int, result *interfaces.ConfigValidationResult) {
	// Check for dangerous commands
	dangerousCommands := []string{"rm -rf /", "chmod 777", "sudo rm"}
	for _, dangerous := range dangerousCommands {
		if strings.Contains(command, dangerous) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("target_%s_cmd_%d", targetName, commandIndex),
				Value:      command,
				Type:       "dangerous_command",
				Message:    fmt.Sprintf("Potentially dangerous command: %s", dangerous),
				Suggestion: "Review command for safety and add appropriate safeguards",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "makefile.dangerous_command",
			})
			result.Summary.WarningCount++
		}
	}

	// Check for commands that should use variables
	if mv.shouldUseVariable(command) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("target_%s_cmd_%d", targetName, commandIndex),
			Value:      command,
			Type:       "hardcoded_command",
			Message:    "Command uses hardcoded values that could be variables",
			Suggestion: "Consider using variables for repeated values",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "makefile.use_variables",
		})
		result.Summary.WarningCount++
	}

	// Check for missing error handling
	if mv.needsErrorHandling(command) && !strings.Contains(command, "||") && !strings.HasPrefix(command, "-") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("target_%s_cmd_%d", targetName, commandIndex),
			Value:      command,
			Type:       "missing_error_handling",
			Message:    "Command might need error handling",
			Suggestion: "Consider adding error handling with '||' or prefix with '-' to ignore errors",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "makefile.error_handling",
		})
		result.Summary.WarningCount++
	}

	// Check for shell-specific commands
	mv.validateShellCommand(targetName, command, commandIndex, result)
}

// validateSyntaxAndBestPractices validates general syntax and best practices
func (mv *MakefileValidator) validateSyntaxAndBestPractices(content string, result *interfaces.ConfigValidationResult) {
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Check for spaces instead of tabs in command lines
		if strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			// Check if this looks like a command (indented line after a target)
			if mv.looksLikeCommand(line) {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    fmt.Sprintf("line_%d", lineNum+1),
					Value:    line,
					Type:     "syntax_error",
					Message:  "Commands must be indented with tabs, not spaces",
					Severity: interfaces.ValidationSeverityError,
					Rule:     "makefile.tab_indentation",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}
		}

		// Check for trailing whitespace
		if len(line) > 0 && line != strings.TrimRight(line, " \t") {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      line,
				Type:       "format_warning",
				Message:    "Line has trailing whitespace",
				Suggestion: "Remove trailing whitespace",
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "makefile.trailing_whitespace",
			})
			result.Summary.WarningCount++
		}

		// Check for very long lines
		if len(line) > 120 {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      line,
				Type:       "format_warning",
				Message:    "Line is very long (>120 characters)",
				Suggestion: "Consider breaking long lines with backslash continuation",
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "makefile.long_lines",
			})
			result.Summary.WarningCount++
		}
	}
}

// Helper methods

// Pre-compiled regular expressions for makefile validation
var (
	validTargetNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	validVariableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// validateTargetName validates target naming conventions
func (mv *MakefileValidator) validateTargetName(name string) error {
	if name == "" {
		return fmt.Errorf("target name cannot be empty")
	}

	// Check for spaces
	if strings.Contains(name, " ") {
		return fmt.Errorf("target name cannot contain spaces")
	}

	// Check for invalid characters using pre-compiled regex
	if !validTargetNameRegex.MatchString(name) {
		return fmt.Errorf("target name contains invalid characters")
	}

	// Recommend lowercase
	if strings.ToLower(name) != name && !strings.HasPrefix(name, ".") {
		return fmt.Errorf("consider using lowercase for target names")
	}

	return nil
}

// validateVariableName validates variable naming conventions
func (mv *MakefileValidator) validateVariableName(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	// Check for spaces
	if strings.Contains(name, " ") {
		return fmt.Errorf("variable name cannot contain spaces")
	}

	// Check for invalid characters using pre-compiled regex
	if !validVariableNameRegex.MatchString(name) {
		return fmt.Errorf("variable name contains invalid characters")
	}

	// Recommend uppercase for variables
	if strings.ToUpper(name) != name {
		return fmt.Errorf("consider using uppercase for variable names")
	}

	return nil
}

// shouldBePhony determines if a target should be declared as .PHONY
func (mv *MakefileValidator) shouldBePhony(target MakefileTarget) bool {
	phonyTargets := []string{"all", "clean", "test", "install", "help", "build", "run", "start", "stop"}

	for _, phony := range phonyTargets {
		if target.Name == phony {
			return true
		}
	}

	return false
}

// validateAssignmentOperator validates the use of assignment operators
func (mv *MakefileValidator) validateAssignmentOperator(variable MakefileVariable, result *interfaces.ConfigValidationResult) {
	switch variable.Type {
	case "=":
		// Recursive assignment - warn if it might cause issues
		if strings.Contains(variable.Value, variable.Name) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", variable.LineNum),
				Value:      variable.Name,
				Type:       "recursive_assignment",
				Message:    "Recursive assignment might cause infinite recursion",
				Suggestion: "Consider using := for simple assignment",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "makefile.recursive_assignment",
			})
			result.Summary.WarningCount++
		}
	case ":=":
		// Simple assignment - generally good
	case "+=":
		// Append assignment - check if variable was defined first
	case "?=":
		// Conditional assignment - generally good for defaults
	}
}

// containsHardcodedPaths checks if a value contains hardcoded paths
func (mv *MakefileValidator) containsHardcodedPaths(value string) bool {
	hardcodedPatterns := []string{"/usr/local/", "/opt/", "/home/", "C:\\", "/Users/"}

	for _, pattern := range hardcodedPatterns {
		if strings.Contains(value, pattern) {
			return true
		}
	}

	return false
}

// shouldUseVariable checks if a command should use variables instead of hardcoded values
func (mv *MakefileValidator) shouldUseVariable(command string) bool {
	// Check for repeated compiler names, common paths, etc.
	repeatedValues := []string{"gcc", "g++", "clang", "go build", "npm", "yarn", "docker"}

	for _, value := range repeatedValues {
		if strings.Contains(command, value) {
			return true
		}
	}

	return false
}

// needsErrorHandling checks if a command might need error handling
func (mv *MakefileValidator) needsErrorHandling(command string) bool {
	riskyCommands := []string{"curl", "wget", "git", "docker", "npm install", "go get"}

	for _, risky := range riskyCommands {
		if strings.Contains(command, risky) {
			return true
		}
	}

	return false
}

// validateShellCommand validates shell-specific commands
func (mv *MakefileValidator) validateShellCommand(targetName, command string, commandIndex int, result *interfaces.ConfigValidationResult) {
	// Check for bash-specific features without setting SHELL
	bashFeatures := []string{"[[", "((", "source ", ". "}

	for _, feature := range bashFeatures {
		if strings.Contains(command, feature) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("target_%s_cmd_%d", targetName, commandIndex),
				Value:      command,
				Type:       "shell_compatibility",
				Message:    "Command uses bash-specific features",
				Suggestion: "Set SHELL := /bin/bash or use POSIX-compatible syntax",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "makefile.bash_features",
			})
			result.Summary.WarningCount++
		}
	}
}

// looksLikeCommand checks if an indented line looks like a command
func (mv *MakefileValidator) looksLikeCommand(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Check if it starts with common command patterns
	commandPatterns := []string{"echo", "cd", "mkdir", "rm", "cp", "mv", "ls", "cat", "grep", "sed", "awk", "go", "npm", "yarn", "docker", "git"}

	for _, pattern := range commandPatterns {
		if strings.HasPrefix(trimmed, pattern+" ") || trimmed == pattern {
			return true
		}
	}

	// Check if it looks like a shell command (contains common shell operators)
	shellOperators := []string{"|", "&&", "||", ">", "<", ">>"}
	for _, op := range shellOperators {
		if strings.Contains(trimmed, op) {
			return true
		}
	}

	return false
}

// initializeShellCommands initializes common shell commands and their descriptions
func (mv *MakefileValidator) initializeShellCommands() {
	mv.shellCommands = map[string]string{
		"rm -rf":    "Dangerous recursive delete",
		"chmod 777": "Overly permissive permissions",
		"sudo":      "Requires elevated privileges",
		"curl":      "Network operation that might fail",
		"wget":      "Network operation that might fail",
		"git clone": "Network operation that might fail",
		"docker":    "Docker operation that might fail",
	}
}
