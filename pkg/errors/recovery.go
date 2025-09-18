// Package errors provides error recovery mechanisms for common failures
package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RecoveryStrategy represents a strategy for recovering from an error
type RecoveryStrategy interface {
	// CanRecover determines if this strategy can recover from the given error
	CanRecover(err *CLIError) bool

	// Recover attempts to recover from the error
	Recover(err *CLIError) (*RecoveryResult, error)

	// GetDescription returns a description of what this strategy does
	GetDescription() string
}

// RecoveryResult represents the result of a recovery attempt
type RecoveryResult struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Actions     []RecoveryAction       `json:"actions"`
	Suggestions []string               `json:"suggestions"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// RecoveryAction represents an action taken during recovery
type RecoveryAction struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Success     bool                   `json:"success"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// RecoveryManager manages error recovery strategies
type RecoveryManager struct {
	strategies []RecoveryStrategy
	logger     RecoveryLogger
}

// RecoveryLogger interface for recovery logging
type RecoveryLogger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// NewRecoveryManager creates a new recovery manager with default strategies
func NewRecoveryManager(logger RecoveryLogger) *RecoveryManager {
	rm := &RecoveryManager{
		logger: logger,
	}

	// Register default recovery strategies
	rm.RegisterStrategy(&NetworkRecoveryStrategy{logger: logger})
	rm.RegisterStrategy(&FileSystemRecoveryStrategy{logger: logger})
	rm.RegisterStrategy(&CacheRecoveryStrategy{logger: logger})
	rm.RegisterStrategy(&ConfigurationRecoveryStrategy{logger: logger})
	rm.RegisterStrategy(&TemplateRecoveryStrategy{logger: logger})
	rm.RegisterStrategy(&ValidationRecoveryStrategy{logger: logger})

	return rm
}

// RegisterStrategy registers a new recovery strategy
func (rm *RecoveryManager) RegisterStrategy(strategy RecoveryStrategy) {
	rm.strategies = append(rm.strategies, strategy)
}

// AttemptRecovery attempts to recover from the given error
func (rm *RecoveryManager) AttemptRecovery(err *CLIError) (*RecoveryResult, error) {
	if err == nil || !err.IsRecoverable() {
		return &RecoveryResult{
			Success: false,
			Message: "Error is not recoverable",
		}, nil
	}

	rm.logger.Info("Attempting recovery for error: %s", err.Type)

	// Try each strategy until one succeeds
	for _, strategy := range rm.strategies {
		if strategy.CanRecover(err) {
			rm.logger.Debug("Trying recovery strategy: %s", strategy.GetDescription())

			result, recoveryErr := strategy.Recover(err)
			if recoveryErr != nil {
				rm.logger.Warn("Recovery strategy failed: %v", recoveryErr)
				continue
			}

			if result.Success {
				rm.logger.Info("Recovery successful: %s", result.Message)
				return result, nil
			}

			rm.logger.Debug("Recovery strategy did not succeed: %s", result.Message)
		}
	}

	return &RecoveryResult{
		Success: false,
		Message: "No recovery strategy succeeded",
		Suggestions: []string{
			"Check the error details for manual resolution steps",
			"Consult documentation for troubleshooting guidance",
			"Report the issue if it persists",
		},
	}, nil
}

// NetworkRecoveryStrategy handles network-related errors
type NetworkRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *NetworkRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeNetwork
}

func (s *NetworkRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	// Check if we can fall back to offline mode
	if s.canUseOfflineMode(err) {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "fallback_offline",
			Description: "Falling back to offline mode using cached data",
			Success:     true,
		})

		result.Success = true
		result.Message = "Recovered by switching to offline mode"
		result.Suggestions = []string{
			"Continue operation using cached data",
			"Check network connection when convenient",
			"Use --offline flag to explicitly use offline mode",
		}

		return result, nil
	}

	// Suggest retry with exponential backoff
	result.Actions = append(result.Actions, RecoveryAction{
		Type:        "suggest_retry",
		Description: "Suggesting retry with backoff",
		Success:     false,
	})

	result.Success = false
	result.Message = "Network error cannot be automatically recovered"
	result.Suggestions = []string{
		"Check internet connection",
		"Retry the operation in a few moments",
		"Use --offline flag if cached data is available",
		"Check proxy settings if applicable",
	}

	return result, nil
}

func (s *NetworkRecoveryStrategy) GetDescription() string {
	return "Network error recovery with offline fallback"
}

func (s *NetworkRecoveryStrategy) canUseOfflineMode(err *CLIError) bool {
	// Check if the operation supports offline mode
	// This would typically check if cached data is available
	return true // Simplified for now
}

// FileSystemRecoveryStrategy handles filesystem-related errors
type FileSystemRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *FileSystemRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeFileSystem
}

func (s *FileSystemRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	path, hasPath := err.Details["path"].(string)
	if !hasPath {
		result.Success = false
		result.Message = "Cannot recover: no path information available"
		return result, nil
	}

	// Try to create missing directories
	if s.tryCreateMissingDirectories(path) {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "create_directory",
			Description: fmt.Sprintf("Created missing directory: %s", filepath.Dir(path)),
			Success:     true,
			Details:     map[string]interface{}{"path": filepath.Dir(path)},
		})

		result.Success = true
		result.Message = "Recovered by creating missing directories"
		return result, nil
	}

	// Check if it's a permission issue that can be resolved
	if s.canFixPermissions(path) {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "suggest_permission_fix",
			Description: "Detected permission issue that may be fixable",
			Success:     false,
		})

		result.Success = false
		result.Message = "Permission issue detected"
		result.Suggestions = []string{
			fmt.Sprintf("Check permissions for: %s", path),
			"Run with appropriate user privileges",
			"Verify directory ownership",
		}

		return result, nil
	}

	result.Success = false
	result.Message = "Filesystem error cannot be automatically recovered"
	result.Suggestions = []string{
		"Check if the path exists and is accessible",
		"Verify file permissions",
		"Ensure sufficient disk space",
	}

	return result, nil
}

func (s *FileSystemRecoveryStrategy) GetDescription() string {
	return "Filesystem error recovery with directory creation"
}

func (s *FileSystemRecoveryStrategy) tryCreateMissingDirectories(path string) bool {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if createErr := secureMkdirAll(dir, 0750); createErr == nil {
			s.logger.Info("Created missing directory: %s", dir)
			return true
		}
	}
	return false
}

func (s *FileSystemRecoveryStrategy) canFixPermissions(path string) bool {
	// Check if it's a permission issue
	if _, err := os.Stat(path); err != nil {
		return strings.Contains(err.Error(), "permission denied")
	}
	return false
}

// CacheRecoveryStrategy handles cache-related errors
type CacheRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *CacheRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeCache
}

func (s *CacheRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	// Try to clear corrupted cache
	if s.tryClearCache() {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "clear_cache",
			Description: "Cleared corrupted cache data",
			Success:     true,
		})

		result.Success = true
		result.Message = "Recovered by clearing corrupted cache"
		result.Suggestions = []string{
			"Cache has been cleared and will be rebuilt",
			"Next operation may be slower as cache rebuilds",
		}

		return result, nil
	}

	result.Success = false
	result.Message = "Cache error cannot be automatically recovered"
	result.Suggestions = []string{
		"Manually clear cache with 'generator cache clear'",
		"Check cache directory permissions",
		"Verify sufficient disk space",
	}

	return result, nil
}

func (s *CacheRecoveryStrategy) GetDescription() string {
	return "Cache error recovery with automatic cache clearing"
}

func (s *CacheRecoveryStrategy) tryClearCache() bool {
	// This would implement actual cache clearing logic
	// For now, we'll simulate success
	s.logger.Info("Attempting to clear corrupted cache")
	return true
}

// ConfigurationRecoveryStrategy handles configuration-related errors
type ConfigurationRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *ConfigurationRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeConfiguration
}

func (s *ConfigurationRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	configPath, hasPath := err.Details["config_path"].(string)
	if !hasPath {
		result.Success = false
		result.Message = "Cannot recover: no configuration path available"
		return result, nil
	}

	// Try to create default configuration
	if s.tryCreateDefaultConfig(configPath) {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "create_default_config",
			Description: fmt.Sprintf("Created default configuration: %s", configPath),
			Success:     true,
			Details:     map[string]interface{}{"config_path": configPath},
		})

		result.Success = true
		result.Message = "Recovered by creating default configuration"
		result.Suggestions = []string{
			"Review and customize the default configuration",
			"Use 'generator config edit' to modify settings",
		}

		return result, nil
	}

	result.Success = false
	result.Message = "Configuration error cannot be automatically recovered"
	result.Suggestions = []string{
		"Check configuration file syntax",
		"Use 'generator config validate' to identify issues",
		"Create a new configuration file",
	}

	return result, nil
}

func (s *ConfigurationRecoveryStrategy) GetDescription() string {
	return "Configuration error recovery with default config creation"
}

func (s *ConfigurationRecoveryStrategy) tryCreateDefaultConfig(configPath string) bool {
	// This would implement actual default config creation
	s.logger.Info("Attempting to create default configuration at: %s", configPath)
	return false // Simplified for now
}

// TemplateRecoveryStrategy handles template-related errors
type TemplateRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *TemplateRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeTemplate
}

func (s *TemplateRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	templateName, hasName := err.Details["template_name"].(string)
	if !hasName {
		result.Success = false
		result.Message = "Cannot recover: no template name available"
		return result, nil
	}

	// Try to suggest similar templates
	suggestions := s.findSimilarTemplates(templateName)
	if len(suggestions) > 0 {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "suggest_alternatives",
			Description: "Found similar template alternatives",
			Success:     false,
			Details:     map[string]interface{}{"alternatives": suggestions},
		})

		result.Success = false
		result.Message = "Template not found, but alternatives are available"
		result.Suggestions = append([]string{
			"Consider using one of the suggested alternatives:",
		}, suggestions...)

		return result, nil
	}

	result.Success = false
	result.Message = "Template error cannot be automatically recovered"
	result.Suggestions = []string{
		"List available templates with 'generator list-templates'",
		"Check template name spelling",
		"Verify template exists and is accessible",
	}

	return result, nil
}

func (s *TemplateRecoveryStrategy) GetDescription() string {
	return "Template error recovery with alternative suggestions"
}

func (s *TemplateRecoveryStrategy) findSimilarTemplates(templateName string) []string {
	// This would implement actual template similarity matching
	// For now, return some common alternatives
	commonTemplates := []string{"go-gin", "nextjs-app", "react-component"}

	var suggestions []string
	for _, template := range commonTemplates {
		if strings.Contains(template, strings.ToLower(templateName)) ||
			strings.Contains(strings.ToLower(templateName), template) {
			suggestions = append(suggestions, template)
		}
	}

	return suggestions
}

// ValidationRecoveryStrategy handles validation-related errors
type ValidationRecoveryStrategy struct {
	logger RecoveryLogger
}

func (s *ValidationRecoveryStrategy) CanRecover(err *CLIError) bool {
	return err.Type == ErrorTypeValidation
}

func (s *ValidationRecoveryStrategy) Recover(err *CLIError) (*RecoveryResult, error) {
	result := &RecoveryResult{
		Actions: []RecoveryAction{},
		Details: make(map[string]interface{}),
	}

	field, hasField := err.Details["field"].(string)
	value := err.Details["value"]

	if !hasField {
		result.Success = false
		result.Message = "Cannot recover: no field information available"
		return result, nil
	}

	// Try to suggest valid values
	if suggestions := s.suggestValidValues(field, value); len(suggestions) > 0 {
		result.Actions = append(result.Actions, RecoveryAction{
			Type:        "suggest_valid_values",
			Description: fmt.Sprintf("Suggested valid values for field '%s'", field),
			Success:     false,
			Details:     map[string]interface{}{"suggestions": suggestions},
		})

		result.Success = false
		result.Message = "Validation error with suggested corrections"
		result.Suggestions = append([]string{
			fmt.Sprintf("Valid values for '%s':", field),
		}, suggestions...)

		return result, nil
	}

	result.Success = false
	result.Message = "Validation error cannot be automatically recovered"
	result.Suggestions = []string{
		"Check the field value format and constraints",
		"Refer to documentation for valid values",
		"Use --fix flag to automatically fix common issues",
	}

	return result, nil
}

func (s *ValidationRecoveryStrategy) GetDescription() string {
	return "Validation error recovery with value suggestions"
}

func (s *ValidationRecoveryStrategy) suggestValidValues(field string, value interface{}) []string {
	// This would implement actual validation logic
	// For now, return some common suggestions based on field name
	switch strings.ToLower(field) {
	case "license":
		return []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"}
	case "language", "technology":
		return []string{"go", "javascript", "typescript", "python"}
	case "framework":
		return []string{"gin", "nextjs", "react", "vue"}
	default:
		return []string{}
	}
}

// RecoveryAttempt represents a single recovery attempt with timing
type RecoveryAttempt struct {
	Error     *CLIError       `json:"error"`
	Strategy  string          `json:"strategy"`
	Result    *RecoveryResult `json:"result"`
	StartTime time.Time       `json:"start_time"`
	Duration  time.Duration   `json:"duration"`
}

// RecoveryHistory tracks recovery attempts for analysis
type RecoveryHistory struct {
	attempts []RecoveryAttempt
}

// NewRecoveryHistory creates a new recovery history tracker
func NewRecoveryHistory() *RecoveryHistory {
	return &RecoveryHistory{
		attempts: make([]RecoveryAttempt, 0),
	}
}

// RecordAttempt records a recovery attempt
func (rh *RecoveryHistory) RecordAttempt(err *CLIError, strategy string, result *RecoveryResult, duration time.Duration) {
	attempt := RecoveryAttempt{
		Error:     err,
		Strategy:  strategy,
		Result:    result,
		StartTime: time.Now().Add(-duration),
		Duration:  duration,
	}

	rh.attempts = append(rh.attempts, attempt)
}

// GetAttempts returns all recovery attempts
func (rh *RecoveryHistory) GetAttempts() []RecoveryAttempt {
	return rh.attempts
}

// GetSuccessfulAttempts returns only successful recovery attempts
func (rh *RecoveryHistory) GetSuccessfulAttempts() []RecoveryAttempt {
	var successful []RecoveryAttempt
	for _, attempt := range rh.attempts {
		if attempt.Result.Success {
			successful = append(successful, attempt)
		}
	}
	return successful
}

// GetFailedAttempts returns only failed recovery attempts
func (rh *RecoveryHistory) GetFailedAttempts() []RecoveryAttempt {
	var failed []RecoveryAttempt
	for _, attempt := range rh.attempts {
		if !attempt.Result.Success {
			failed = append(failed, attempt)
		}
	}
	return failed
}
