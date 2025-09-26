package models

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid   bool              `json:"valid"`
	Issues  []ValidationIssue `json:"issues"`
	Summary string            `json:"summary"`
}

// ValidationIssue represents a single validation issue
type ValidationIssue struct {
	Type     string            `json:"type"`
	Severity string            `json:"severity"`
	Message  string            `json:"message"`
	File     string            `json:"file,omitempty"`
	Line     int               `json:"line,omitempty"`
	Column   int               `json:"column,omitempty"`
	Rule     string            `json:"rule"`
	Fixable  bool              `json:"fixable"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ConfigValidator handles basic configuration validation
type ConfigValidator struct {
	// Basic validator with no complex features
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateProjectConfig validates a project configuration
func (cv *ConfigValidator) ValidateProjectConfig(config *ProjectConfig) *ValidationResult {
	result := &ValidationResult{
		Valid:   true,
		Issues:  []ValidationIssue{},
		Summary: "Configuration validation completed",
	}

	// Basic validation
	if config == nil {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Message: "Configuration cannot be nil",
		})
		return result
	}

	if config.Name == "" {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Message: "Project name is required",
		})
	}

	if config.Organization == "" {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Message: "Organization is required",
		})
	}

	if config.OutputPath == "" {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:    "error",
			Message: "Output path is required",
		})
	}

	return result
}
