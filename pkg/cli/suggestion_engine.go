package cli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
)

// SuggestionEngine generates contextual recovery suggestions for errors
type SuggestionEngine struct {
	verbose bool
}

// NewSuggestionEngine creates a new suggestion engine
func NewSuggestionEngine(verbose bool) *SuggestionEngine {
	return &SuggestionEngine{
		verbose: verbose,
	}
}

// GenerateSuggestions analyzes an error and returns actionable suggestions
func (se *SuggestionEngine) GenerateSuggestions(err error) []string {
	if err == nil {
		return nil
	}

	suggestions := make([]string, 0)

	// Check for CLIError types
	if cliErr, ok := err.(*CLIError); ok {
		// Return existing suggestions if available
		if len(cliErr.Suggestions) > 0 {
			return cliErr.Suggestions
		}

		// Generate suggestions based on category
		switch cliErr.Category {
		case "configuration":
			suggestions = se.getConfigSuggestions(cliErr)
		case "tool":
			suggestions = se.getToolSuggestions(cliErr)
		case "generation":
			suggestions = se.getGenerationSuggestions(cliErr)
		case "filesystem":
			suggestions = se.getFileSystemSuggestions(cliErr)
		case "user":
			suggestions = se.getUserSuggestions(cliErr)
		default:
			suggestions = se.getGenericSuggestions(cliErr)
		}

		return suggestions
	}

	// Check for GenerationError types
	if genErr, ok := err.(*orchestrator.GenerationError); ok {
		// Return existing suggestions if available
		if len(genErr.Suggestions) > 0 {
			return genErr.Suggestions
		}

		// Generate suggestions based on category
		switch genErr.Category {
		case orchestrator.ErrCategoryToolNotFound:
			suggestions = se.getToolNotFoundSuggestions(genErr)
		case orchestrator.ErrCategoryToolExecution:
			suggestions = se.getToolExecutionSuggestions(genErr)
		case orchestrator.ErrCategoryInvalidConfig:
			suggestions = se.getInvalidConfigSuggestions(genErr)
		case orchestrator.ErrCategoryFileSystem:
			suggestions = se.getFileSystemErrorSuggestions(genErr)
		case orchestrator.ErrCategorySecurity:
			suggestions = se.getSecuritySuggestions(genErr)
		case orchestrator.ErrCategoryIntegration:
			suggestions = se.getIntegrationSuggestions(genErr)
		case orchestrator.ErrCategoryValidation:
			suggestions = se.getValidationSuggestions(genErr)
		case orchestrator.ErrCategoryTimeout:
			suggestions = se.getTimeoutSuggestions(genErr)
		default:
			suggestions = se.getGenericGenerationSuggestions(genErr)
		}

		return suggestions
	}

	// Generic error handling
	errMsg := err.Error()

	// Check for common error patterns
	if strings.Contains(errMsg, "network") || strings.Contains(errMsg, "connection") {
		suggestions = append(suggestions, "Check your network connection")
		suggestions = append(suggestions, "Try enabling offline mode with cached tools")
		suggestions = append(suggestions, "Run 'generator cache-tools --refresh' when online")
	} else if strings.Contains(errMsg, "permission denied") {
		suggestions = append(suggestions, "Check file system permissions")
		suggestions = append(suggestions, "Try running with appropriate permissions")
		suggestions = append(suggestions, "Verify the output directory is writable")
	} else if strings.Contains(errMsg, "not found") {
		suggestions = append(suggestions, "Verify the path or resource exists")
		suggestions = append(suggestions, "Check for typos in the configuration")
		suggestions = append(suggestions, "Run 'generator check-tools' to verify tool availability")
	}

	// Add verbose mode suggestion if not already verbose
	if !se.verbose && len(suggestions) > 0 {
		suggestions = append(suggestions, "Run with --verbose flag for detailed diagnostic information")
	}

	// If no specific suggestions, provide generic ones
	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Check the error message above for details")
		suggestions = append(suggestions, "Run with --verbose flag for more information")
		suggestions = append(suggestions, "Consult the documentation for troubleshooting guidance")
	}

	return suggestions
}

// GetToolInstallSuggestion returns installation suggestion for a missing tool
func (se *SuggestionEngine) GetToolInstallSuggestion(toolName string) string {
	os := runtime.GOOS

	switch toolName {
	case "npx", "npm":
		return se.getNodeInstallSuggestion(os)
	case "go":
		return se.getGoInstallSuggestion(os)
	case "gradle":
		return se.getGradleInstallSuggestion(os)
	case "xcodebuild":
		return se.getXcodeInstallSuggestion(os)
	case "docker":
		return se.getDockerInstallSuggestion(os)
	case "kubectl":
		return se.getKubectlInstallSuggestion(os)
	default:
		return fmt.Sprintf("Install %s on your system. Refer to the official documentation for installation instructions.", toolName)
	}
}

// GetConfigFixSuggestion returns suggestion for fixing config errors
func (se *SuggestionEngine) GetConfigFixSuggestion(field string, expectedType string) string {
	return fmt.Sprintf("Field '%s' should be of type %s. Check your configuration file and correct the value.", field, expectedType)
}

// getConfigSuggestions returns suggestions for configuration errors
func (se *SuggestionEngine) getConfigSuggestions(err *CLIError) []string {
	suggestions := []string{
		"Review your configuration file for syntax errors",
		"Use 'generator init-config --example <type>' to generate a valid template",
		"Refer to the configuration schema documentation",
	}

	if se.verbose {
		suggestions = append(suggestions, "Check the detailed error message above for the specific field causing issues")
	}

	return suggestions
}

// getToolSuggestions returns suggestions for tool errors
func (se *SuggestionEngine) getToolSuggestions(err *CLIError) []string {
	suggestions := []string{}

	// Extract tool name from error message
	toolName := se.extractToolName(err.Message)
	if toolName != "" {
		installSuggestion := se.GetToolInstallSuggestion(toolName)
		suggestions = append(suggestions, installSuggestion)
	}

	suggestions = append(suggestions,
		"Run 'generator check-tools' to see all required tools and their status",
		"Use --no-external-tools flag to force fallback generation without bootstrap tools",
	)

	return suggestions
}

// getGenerationSuggestions returns suggestions for generation errors
func (se *SuggestionEngine) getGenerationSuggestions(err *CLIError) []string {
	suggestions := []string{
		"Check the error messages above for specific component failures",
		"Try running with --verbose flag for detailed execution logs",
		"Use --no-external-tools to attempt fallback generation",
	}

	if strings.Contains(err.Message, "bootstrap") {
		suggestions = append(suggestions, "Verify the bootstrap tool is properly installed and configured")
	}

	if strings.Contains(err.Message, "fallback") {
		suggestions = append(suggestions, "Check that fallback templates are available")
		suggestions = append(suggestions, "Verify the component type is supported")
	}

	return suggestions
}

// getFileSystemSuggestions returns suggestions for file system errors
func (se *SuggestionEngine) getFileSystemSuggestions(err *CLIError) []string {
	suggestions := []string{}

	if strings.Contains(err.Message, "permission") {
		suggestions = append(suggestions,
			"Check file system permissions for the output directory",
			"Ensure you have write access to the target location",
		)

		if runtime.GOOS != "windows" {
			suggestions = append(suggestions, "Try running 'chmod +w <directory>' to grant write permissions")
		}
	}

	if strings.Contains(err.Message, "space") || strings.Contains(err.Message, "disk") {
		suggestions = append(suggestions,
			"Check available disk space with 'df -h' (Unix) or 'dir' (Windows)",
			"Free up disk space and try again",
		)
	}

	if strings.Contains(err.Message, "not found") || strings.Contains(err.Message, "no such") {
		suggestions = append(suggestions,
			"Verify the path exists and is accessible",
			"Check for typos in the path specification",
		)
	}

	// Generic file system suggestions
	if len(suggestions) == 0 {
		suggestions = append(suggestions,
			"Verify the output directory is accessible",
			"Check file system permissions and available space",
			"Ensure no other process is locking the files",
		)
	}

	return suggestions
}

// getUserSuggestions returns suggestions for user-related errors
func (se *SuggestionEngine) getUserSuggestions(err *CLIError) []string {
	return []string{
		"Run the command again when ready",
		"Use non-interactive mode with a configuration file for automation",
	}
}

// getGenericSuggestions returns generic suggestions for CLI errors
func (se *SuggestionEngine) getGenericSuggestions(err *CLIError) []string {
	return []string{
		"Check the error message for specific details",
		"Run with --verbose flag for more diagnostic information",
		"Consult the documentation for troubleshooting guidance",
	}
}

// getToolNotFoundSuggestions returns suggestions for tool not found errors
func (se *SuggestionEngine) getToolNotFoundSuggestions(err *orchestrator.GenerationError) []string {
	suggestions := []string{}

	// Extract tool name from error message
	toolName := se.extractToolName(err.Message)
	if toolName != "" {
		installSuggestion := se.GetToolInstallSuggestion(toolName)
		suggestions = append(suggestions, installSuggestion)
	}

	suggestions = append(suggestions,
		"Run 'generator check-tools' to see installation instructions for all required tools",
		"Use --no-external-tools flag to force fallback generation",
	)

	if err.Component != "" {
		suggestions = append(suggestions, fmt.Sprintf("Fallback generation is available for '%s' component", err.Component))
	}

	return suggestions
}

// getToolExecutionSuggestions returns suggestions for tool execution errors
func (se *SuggestionEngine) getToolExecutionSuggestions(err *orchestrator.GenerationError) []string {
	suggestions := []string{
		"Check the tool output above for specific error messages",
		"Verify the tool is properly installed and configured",
		"Try running the tool manually to diagnose the issue",
	}

	if se.verbose {
		suggestions = append(suggestions, "Review the detailed execution logs above")
	} else {
		suggestions = append(suggestions, "Run with --verbose flag for detailed execution logs")
	}

	suggestions = append(suggestions, "Use --no-external-tools to try fallback generation")

	return suggestions
}

// getInvalidConfigSuggestions returns suggestions for invalid config errors
func (se *SuggestionEngine) getInvalidConfigSuggestions(err *orchestrator.GenerationError) []string {
	suggestions := []string{
		"Review your configuration file for the reported field",
		"Check the configuration schema documentation for valid values",
		"Use 'generator init-config --example <type>' to see a valid example",
	}

	// Extract field name from error message
	if strings.Contains(err.Message, "field '") {
		suggestions = append(suggestions, "Correct the field value according to the validation rules")
	}

	return suggestions
}

// getFileSystemErrorSuggestions returns suggestions for file system errors
func (se *SuggestionEngine) getFileSystemErrorSuggestions(err *orchestrator.GenerationError) []string {
	return se.getFileSystemSuggestions(&CLIError{Message: err.Message})
}

// getSecuritySuggestions returns suggestions for security errors
func (se *SuggestionEngine) getSecuritySuggestions(err *orchestrator.GenerationError) []string {
	return []string{
		"Review the input for potentially dangerous patterns",
		"Ensure paths do not contain directory traversal attempts (../ or ..\\)",
		"Verify all inputs are properly sanitized",
		"Check that file paths are within the allowed output directory",
	}
}

// getIntegrationSuggestions returns suggestions for integration errors
func (se *SuggestionEngine) getIntegrationSuggestions(err *orchestrator.GenerationError) []string {
	return []string{
		"Verify all components were generated successfully",
		"Check that component configurations are compatible",
		"Review integration logs for specific errors",
		"Try generating components individually to isolate the issue",
	}
}

// getValidationSuggestions returns suggestions for validation errors
func (se *SuggestionEngine) getValidationSuggestions(err *orchestrator.GenerationError) []string {
	return []string{
		"Check the generated project structure",
		"Verify all required files were created",
		"Review validation logs for specific missing files or directories",
		"Try regenerating the project with --verbose flag",
	}
}

// getTimeoutSuggestions returns suggestions for timeout errors
func (se *SuggestionEngine) getTimeoutSuggestions(err *orchestrator.GenerationError) []string {
	suggestions := []string{
		"The operation took longer than expected",
		"Check your network connection if downloading dependencies",
		"Verify system resources (CPU, memory) are not constrained",
	}

	if err.Component != "" {
		suggestions = append(suggestions, fmt.Sprintf("Try generating the '%s' component separately", err.Component))
	}

	suggestions = append(suggestions, "Consider increasing timeout duration in configuration")

	return suggestions
}

// getGenericGenerationSuggestions returns generic suggestions for generation errors
func (se *SuggestionEngine) getGenericGenerationSuggestions(err *orchestrator.GenerationError) []string {
	suggestions := []string{
		"Check the error message for specific details",
	}

	if !se.verbose {
		suggestions = append(suggestions, "Run with --verbose flag for detailed diagnostic information")
	}

	suggestions = append(suggestions,
		"Consult the documentation for troubleshooting guidance",
		"Try running 'generator check-tools' to verify tool availability",
	)

	return suggestions
}

// extractToolName attempts to extract a tool name from an error message
func (se *SuggestionEngine) extractToolName(message string) string {
	// Common tool names to look for
	tools := []string{"npx", "npm", "node", "go", "gradle", "xcodebuild", "docker", "kubectl", "terraform"}

	lowerMsg := strings.ToLower(message)
	for _, tool := range tools {
		if strings.Contains(lowerMsg, tool) {
			return tool
		}
	}

	// Try to extract from patterns like "tool 'name'"
	if strings.Contains(message, "tool '") {
		start := strings.Index(message, "tool '") + 6
		end := strings.Index(message[start:], "'")
		if end > 0 {
			return message[start : start+end]
		}
	}

	return ""
}

// getNodeInstallSuggestion returns Node.js installation suggestion
func (se *SuggestionEngine) getNodeInstallSuggestion(os string) string {
	switch os {
	case "darwin":
		return "Install Node.js: brew install node"
	case "linux":
		return "Install Node.js: Use your package manager (apt install nodejs npm, yum install nodejs, etc.) or download from https://nodejs.org"
	case "windows":
		return "Install Node.js: Download the installer from https://nodejs.org or use 'choco install nodejs' (with Chocolatey)"
	default:
		return "Install Node.js from https://nodejs.org"
	}
}

// getGoInstallSuggestion returns Go installation suggestion
func (se *SuggestionEngine) getGoInstallSuggestion(os string) string {
	switch os {
	case "darwin":
		return "Install Go: brew install go"
	case "linux":
		return "Install Go: Download from https://golang.org/dl/ or use your package manager"
	case "windows":
		return "Install Go: Download the installer from https://golang.org/dl/ or use 'choco install golang' (with Chocolatey)"
	default:
		return "Install Go from https://golang.org/dl/"
	}
}

// getGradleInstallSuggestion returns Gradle installation suggestion
func (se *SuggestionEngine) getGradleInstallSuggestion(os string) string {
	switch os {
	case "darwin":
		return "Install Gradle: brew install gradle"
	case "linux":
		return "Install Gradle: Use SDKMAN (sdk install gradle) or download from https://gradle.org"
	case "windows":
		return "Install Gradle: Use 'choco install gradle' (with Chocolatey) or download from https://gradle.org"
	default:
		return "Install Gradle from https://gradle.org"
	}
}

// getXcodeInstallSuggestion returns Xcode installation suggestion
func (se *SuggestionEngine) getXcodeInstallSuggestion(os string) string {
	if os == "darwin" {
		return "Install Xcode: Download from the Mac App Store or run 'xcode-select --install' for command line tools"
	}
	return "Xcode is only available on macOS. iOS development requires a Mac."
}

// getDockerInstallSuggestion returns Docker installation suggestion
func (se *SuggestionEngine) getDockerInstallSuggestion(os string) string {
	switch os {
	case "darwin":
		return "Install Docker: Download Docker Desktop from https://www.docker.com/products/docker-desktop or use 'brew install --cask docker'"
	case "linux":
		return "Install Docker: Follow instructions at https://docs.docker.com/engine/install/"
	case "windows":
		return "Install Docker: Download Docker Desktop from https://www.docker.com/products/docker-desktop"
	default:
		return "Install Docker from https://www.docker.com/get-started"
	}
}

// getKubectlInstallSuggestion returns kubectl installation suggestion
func (se *SuggestionEngine) getKubectlInstallSuggestion(os string) string {
	switch os {
	case "darwin":
		return "Install kubectl: brew install kubectl"
	case "linux":
		return "Install kubectl: Follow instructions at https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/"
	case "windows":
		return "Install kubectl: Use 'choco install kubernetes-cli' or follow instructions at https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/"
	default:
		return "Install kubectl from https://kubernetes.io/docs/tasks/tools/"
	}
}
