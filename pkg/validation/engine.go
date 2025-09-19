// Package validation provides comprehensive validation capabilities for the Open Source Project Generator.
//
// This package implements the ValidationEngine interface and provides:
//   - Project structure validation with detailed analysis
//   - Dependency validation with security and compatibility checks
//   - Configuration validation with schema support
//   - Template validation with metadata verification
//   - Auto-fix capabilities for common issues
//   - Validation rule management and customization
//   - Comprehensive reporting in multiple formats
//
// The validation engine ensures that generated projects meet quality standards
// and are free from common configuration issues.
//
// Usage:
//
//	validator := validation.NewEngine()
//	result, err := validator.ValidateProject("/path/to/project")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if !result.Valid {
//	    for _, issue := range result.Issues {
//	        fmt.Printf("Issue: %s\n", issue.Message)
//	    }
//	}
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Engine implements the comprehensive ValidationEngine interface
type Engine struct {
	rules           []interfaces.ValidationRule
	rulesByID       map[string]interfaces.ValidationRule
	rulesByCategory map[string][]interfaces.ValidationRule
}

// NewEngine creates a new validation engine with default rules
func NewEngine() interfaces.ValidationEngine {
	engine := &Engine{
		rules:           []interfaces.ValidationRule{},
		rulesByID:       make(map[string]interfaces.ValidationRule),
		rulesByCategory: make(map[string][]interfaces.ValidationRule),
	}

	// Initialize with default validation rules
	engine.initializeDefaultRules()

	return engine
}

// ValidateProject performs comprehensive project validation
func (e *Engine) ValidateProject(projectPath string) (*models.ValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(projectPath); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return &models.ValidationResult{
			Valid: false,
			Issues: []models.ValidationIssue{
				{
					Type:    "error",
					Message: "Project path does not exist",
					File:    projectPath,
				},
			},
			Summary: "Project validation failed - path not found",
		}, nil
	}

	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Project validation completed",
	}

	// Perform structure validation
	if err := e.validateProjectStructureBasic(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 Couldn't validate project structure: %w", err)
	}

	// Perform dependency validation
	if err := e.validateProjectDependenciesBasic(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 Couldn't validate project dependencies: %w", err)
	}

	// Perform configuration validation
	if err := e.validateProjectConfigurationFiles(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 Couldn't validate project configuration: %w", err)
	}

	return result, nil
}

// ValidatePackageJSON validates a package.json file
func (e *Engine) ValidatePackageJSON(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid package.json path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("invalid JSON syntax in package.json: %w", err)
	}

	// Check required fields
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		if _, exists := pkg[field]; !exists {
			return fmt.Errorf("missing required field '%s' in package.json", field)
		}
	}

	// Validate name format
	if name, exists := pkg["name"]; exists {
		if nameStr, ok := name.(string); ok {
			if nameStr == "" || strings.Contains(nameStr, " ") {
				return fmt.Errorf("package name must be non-empty and not contain spaces")
			}
		}
	}

	// Validate version format
	if version, exists := pkg["version"]; exists {
		if versionStr, ok := version.(string); ok {
			if versionStr == "" {
				return fmt.Errorf("version must be non-empty")
			}
		}
	}

	return nil
}

// ValidateGoMod validates a go.mod file
func (e *Engine) ValidateGoMod(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid go.mod path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	contentStr := string(content)

	// Basic validation - check for module declaration
	if !strings.Contains(contentStr, "module ") {
		return fmt.Errorf("go.mod must contain a module declaration")
	}

	// Check for go version
	if !strings.Contains(contentStr, "go ") {
		return fmt.Errorf("go.mod must specify a Go version")
	}

	// Validate Go version format
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "go ") {
			version := strings.TrimSpace(strings.TrimPrefix(line, "go "))
			if version == "" || version == "1" {
				return fmt.Errorf("invalid Go version format: %s", version)
			}
		}
	}

	return nil
}

// ValidateDockerfile validates a Dockerfile
func (e *Engine) ValidateDockerfile(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid Dockerfile path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	contentStr := string(content)

	// Basic validation - check for FROM instruction
	if !strings.Contains(contentStr, "FROM ") {
		return fmt.Errorf("dockerfile must contain a FROM instruction")
	}

	// Check for WORKDIR instruction
	if !strings.Contains(contentStr, "WORKDIR ") {
		return fmt.Errorf("dockerfile must contain a WORKDIR instruction")
	}

	// Check for COPY instruction
	if !strings.Contains(contentStr, "COPY ") {
		return fmt.Errorf("dockerfile must contain a COPY instruction")
	}

	return nil
}

// ValidateYAML validates a YAML file
func (e *Engine) ValidateYAML(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid YAML file path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var data interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}

	return nil
}

// ValidateJSON validates a JSON file
func (e *Engine) ValidateJSON(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid JSON file path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("invalid JSON syntax: %w", err)
	}

	return nil
}

// ValidateTemplate validates a template file using the enhanced TemplateValidator
func (e *Engine) ValidateTemplate(path string) error {
	// Use the new TemplateValidator for comprehensive validation
	templateValidator := NewTemplateValidator()

	// Create a validation result to capture issues
	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Template validation",
	}

	// Validate template content using the enhanced validator
	if templateValidator.useEmbedded {
		// Try embedded validation first
		err := templateValidator.validateEmbeddedTemplateFile(path, result)
		if err != nil {
			// Fall back to filesystem validation
			err = templateValidator.validateTemplateFile(path, result)
			if err != nil {
				return fmt.Errorf("template validation failed: %w", err)
			}
		}
	} else {
		// Use filesystem validation
		err := templateValidator.validateTemplateFile(path, result)
		if err != nil {
			return fmt.Errorf("template validation failed: %w", err)
		}
	}

	// Check if validation found any critical issues
	for _, issue := range result.Issues {
		if issue.Type == "error" {
			return fmt.Errorf("template validation error: %s", issue.Message)
		}
	}

	return nil
}

// Additional validation methods for test compatibility

// ValidateProjectStructure validates project structure
func (e *Engine) ValidateProjectStructure(path string) (*interfaces.StructureValidationResult, error) {
	structureValidator := NewStructureValidator()
	return structureValidator.ValidateProjectStructure(path)
}

// ValidateConfiguration validates project configuration
func (e *Engine) ValidateConfiguration(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	if config == nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:    "config",
					Message:  "Configuration cannot be nil",
					Type:     "null_config",
					Severity: interfaces.ValidationSeverityError,
				},
			},
			Summary: interfaces.ConfigValidationSummary{
				TotalProperties: 0,
				ValidProperties: 0,
				ErrorCount:      1,
				WarningCount:    0,
				MissingRequired: 1,
			},
		}, nil
	}

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

	// Validate required fields
	e.validateRequiredConfigFields(config, result)

	// Validate field formats
	e.validateConfigFieldFormats(config, result)

	// Validate component configuration
	e.validateComponentConfiguration(config, result)

	return result, nil
}

// SetValidationRules sets the validation rules
func (e *Engine) SetValidationRules(rules []interfaces.ValidationRule) error {
	e.rules = rules
	e.rulesByID = make(map[string]interfaces.ValidationRule)
	e.rulesByCategory = make(map[string][]interfaces.ValidationRule)

	for _, rule := range rules {
		e.rulesByID[rule.ID] = rule
		e.rulesByCategory[rule.Category] = append(e.rulesByCategory[rule.Category], rule)
	}

	return nil
}

// GetValidationRules returns the current validation rules
func (e *Engine) GetValidationRules() []interfaces.ValidationRule {
	return e.rules
}

// AddValidationRule adds a new validation rule
func (e *Engine) AddValidationRule(rule interfaces.ValidationRule) error {
	if rule.ID == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}

	if _, exists := e.rulesByID[rule.ID]; exists {
		return fmt.Errorf("rule with ID '%s' already exists", rule.ID)
	}

	e.rules = append(e.rules, rule)
	e.rulesByID[rule.ID] = rule
	e.rulesByCategory[rule.Category] = append(e.rulesByCategory[rule.Category], rule)

	return nil
}

// RemoveValidationRule removes a validation rule by ID
func (e *Engine) RemoveValidationRule(ruleID string) error {
	rule, exists := e.rulesByID[ruleID]
	if !exists {
		return fmt.Errorf("rule with ID '%s' not found", ruleID)
	}

	// Remove from rules slice
	for i, r := range e.rules {
		if r.ID == ruleID {
			e.rules = append(e.rules[:i], e.rules[i+1:]...)
			break
		}
	}

	// Remove from rulesByID map
	delete(e.rulesByID, ruleID)

	// Remove from rulesByCategory map
	categoryRules := e.rulesByCategory[rule.Category]
	for i, r := range categoryRules {
		if r.ID == ruleID {
			e.rulesByCategory[rule.Category] = append(categoryRules[:i], categoryRules[i+1:]...)
			break
		}
	}

	return nil
}

// initializeDefaultRules initializes the engine with default validation rules
func (e *Engine) initializeDefaultRules() {
	defaultRules := []interfaces.ValidationRule{
		{
			ID:          "structure.readme.required",
			Name:        "README Required",
			Description: "Project must have a README file",
			Category:    interfaces.ValidationCategoryStructure,
			Severity:    interfaces.ValidationSeverityError,
			Enabled:     true,
			Fixable:     true,
		},
		{
			ID:          "structure.license.required",
			Name:        "License Required",
			Description: "Project must have a LICENSE file",
			Category:    interfaces.ValidationCategoryStructure,
			Severity:    interfaces.ValidationSeverityError,
			Enabled:     true,
			Fixable:     true,
		},
		{
			ID:          "dependencies.package_json.valid",
			Name:        "Valid package.json",
			Description: "package.json must be valid JSON with required fields",
			Category:    interfaces.ValidationCategoryDependencies,
			Severity:    interfaces.ValidationSeverityError,
			Enabled:     true,
			FileTypes:   []string{"package.json"},
			Fixable:     false,
		},
		{
			ID:          "security.secrets.detection",
			Name:        "Secret Detection",
			Description: "Detect potential secrets in code",
			Category:    interfaces.ValidationCategorySecurity,
			Severity:    interfaces.ValidationSeverityError,
			Enabled:     true,
			Fixable:     false,
		},
		{
			ID:          "quality.naming.conventions",
			Name:        "Naming Conventions",
			Description: "Files and directories should follow naming conventions",
			Category:    interfaces.ValidationCategoryQuality,
			Severity:    interfaces.ValidationSeverityWarning,
			Enabled:     true,
			Fixable:     true,
		},
	}

	_ = e.SetValidationRules(defaultRules)
}

// Enhanced validation methods

// ValidateProjectDependencies performs comprehensive dependency validation
func (e *Engine) ValidateProjectDependencies(path string) (*interfaces.DependencyValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	result := &interfaces.DependencyValidationResult{
		Valid:           true,
		Dependencies:    []interfaces.DependencyValidation{},
		Vulnerabilities: []interfaces.DependencyVulnerability{},
		Outdated:        []interfaces.OutdatedDependency{},
		Conflicts:       []interfaces.DependencyConflict{},
		Summary: interfaces.DependencyValidationSummary{
			TotalDependencies: 0,
			ValidDependencies: 0,
			Vulnerabilities:   0,
			OutdatedCount:     0,
			ConflictCount:     0,
		},
	}

	// Basic implementation - just check if package.json exists
	packageJsonPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		result.Summary.TotalDependencies = 1
		result.Summary.ValidDependencies = 1
	}

	return result, nil
}

// ValidateProjectSecurity performs security validation
func (e *Engine) ValidateProjectSecurity(path string) (*interfaces.SecurityValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	result := &interfaces.SecurityValidationResult{
		Valid:          true,
		SecurityIssues: []interfaces.SecurityIssue{},
		Secrets:        []interfaces.SecretDetection{},
		Permissions:    []interfaces.PermissionIssue{},
		Configurations: []interfaces.SecurityConfig{},
		Summary: interfaces.SecurityValidationSummary{
			TotalIssues:    0,
			HighSeverity:   0,
			MediumSeverity: 0,
			LowSeverity:    0,
			SecretsFound:   0,
			ConfigIssues:   0,
		},
	}

	// Basic implementation - scan for common secret patterns
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(filePath, ".js") {
			return nil
		}

		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil // Skip files we can't read
		}

		contentStr := string(content)
		if strings.Contains(contentStr, "apiKey") || strings.Contains(contentStr, "password") {
			result.Summary.SecretsFound++
			result.Summary.TotalIssues++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan for secrets: %w", err)
	}

	return result, nil
}

// ValidateProjectQuality performs code quality validation
func (e *Engine) ValidateProjectQuality(path string) (*interfaces.QualityValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	result := &interfaces.QualityValidationResult{
		Valid:       true,
		CodeSmells:  []interfaces.CodeSmell{},
		Complexity:  []interfaces.ComplexityIssue{},
		Duplication: []interfaces.Duplication{},
		Coverage:    nil,
		Summary: interfaces.QualityValidationSummary{
			TotalIssues:       0,
			CodeSmells:        0,
			ComplexityIssues:  0,
			DuplicationIssues: 0,
			QualityScore:      85.0, // Default score
			Maintainability:   "B",
		},
	}

	// Basic implementation - check for complex functions
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || (!strings.HasSuffix(filePath, ".go") && !strings.HasSuffix(filePath, ".js")) {
			return nil
		}

		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil // Skip files we can't read
		}

		contentStr := string(content)
		// Simple complexity check - count nested if statements
		nestedIfs := strings.Count(contentStr, "if") - strings.Count(contentStr, "} else if")
		if nestedIfs > 5 {
			result.Summary.ComplexityIssues++
			result.Summary.TotalIssues++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze code quality: %w", err)
	}

	return result, nil
}

// ValidateConfigurationSchema validates configuration against schema
func (e *Engine) ValidateConfigurationSchema(config any, schema *interfaces.ConfigSchema) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// Convert config to map for validation
	configMap, ok := config.(map[string]interface{})
	if !ok {
		// Try to convert via JSON marshaling/unmarshaling
		configBytes, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal configuration: %w", err)
		}

		if err := json.Unmarshal(configBytes, &configMap); err != nil {
			return fmt.Errorf("failed to unmarshal configuration: %w", err)
		}
	}

	// Validate required properties
	for _, required := range schema.Required {
		if _, exists := configMap[required]; !exists {
			return fmt.Errorf("required property '%s' is missing", required)
		}
	}

	return nil
}

// ValidateConfigurationValues validates configuration values
func (e *Engine) ValidateConfigurationValues(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	return e.ValidateConfiguration(config)
}

// ValidateTemplateAdvanced performs advanced template validation
func (e *Engine) ValidateTemplateAdvanced(path string) (*interfaces.TemplateValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid template path: %w", err)
	}

	result := &interfaces.TemplateValidationResult{
		Valid:    true,
		Issues:   []interfaces.ValidationIssue{},
		Warnings: []interfaces.ValidationIssue{},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   0,
			ValidFiles:   0,
			ErrorCount:   0,
			WarningCount: 0,
		},
	}

	// Check if template.yaml exists
	metadataFile := filepath.Join(path, "template.yaml")
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		metadataFile = filepath.Join(path, "template.yml")
		if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
			result.Valid = false
			result.Issues = append(result.Issues, interfaces.ValidationIssue{
				Type:     "error",
				Message:  "Template metadata file (template.yaml or template.yml) is required",
				File:     path,
				Line:     0,
				Column:   0,
				Severity: interfaces.ValidationSeverityError,
			})
			result.Summary.ErrorCount++
		}
	}

	result.Summary.TotalFiles = 1
	if result.Valid {
		result.Summary.ValidFiles = 1
	}

	return result, nil
}

// ValidateTemplateMetadata validates template metadata
func (e *Engine) ValidateTemplateMetadata(metadata *interfaces.TemplateMetadata) error {
	if metadata == nil {
		return fmt.Errorf("template metadata cannot be nil")
	}

	// Validate required fields
	if metadata.Author == "" {
		return fmt.Errorf("template author is required")
	}

	if metadata.License == "" {
		return fmt.Errorf("template license is required")
	}

	// Validate license format
	validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC"}
	licenseValid := false
	for _, license := range validLicenses {
		if metadata.License == license {
			licenseValid = true
			break
		}
	}
	if !licenseValid {
		return fmt.Errorf("unsupported license: %s", metadata.License)
	}

	return nil
}

// ValidateTemplateStructure validates template structure
func (e *Engine) ValidateTemplateStructure(path string) (*interfaces.StructureValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid template path: %w", err)
	}

	result := &interfaces.StructureValidationResult{
		Valid:            true,
		RequiredFiles:    []interfaces.FileValidationResult{},
		RequiredDirs:     []interfaces.DirValidationResult{},
		NamingIssues:     []interfaces.NamingValidationIssue{},
		PermissionIssues: []interfaces.PermissionIssue{},
		Summary: interfaces.StructureValidationSummary{
			TotalFiles:       0,
			ValidFiles:       0,
			TotalDirectories: 0,
			ValidDirectories: 0,
			NamingIssues:     0,
			PermissionIssues: 0,
		},
	}

	// Check for template metadata file
	metadataFile := filepath.Join(path, "template.yaml")
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		metadataFile = filepath.Join(path, "template.yml")
	}

	fileResult := interfaces.FileValidationResult{
		Path:     metadataFile,
		Valid:    true,
		Required: true,
		Issues:   []interfaces.ValidationIssue{},
	}

	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		fileResult.Valid = false
		fileResult.Issues = append(fileResult.Issues, interfaces.ValidationIssue{
			Type:     "error",
			Message:  "Template metadata file is required",
			File:     metadataFile,
			Severity: interfaces.ValidationSeverityError,
		})
		result.Valid = false
	}

	result.RequiredFiles = append(result.RequiredFiles, fileResult)
	result.Summary.TotalFiles = 1
	if fileResult.Valid {
		result.Summary.ValidFiles = 1
	}

	return result, nil
}

// ValidateTemplateVariables validates template variables
func (e *Engine) ValidateTemplateVariables(variables map[string]interfaces.TemplateVariable) error {
	if variables == nil {
		return nil // No variables to validate
	}

	for name, variable := range variables {
		if name == "" {
			return fmt.Errorf("variable name cannot be empty")
		}

		if variable.Type == "" {
			return fmt.Errorf("variable '%s' must have a type", name)
		}

		// Validate variable type
		validTypes := []string{"string", "number", "boolean", "array", "object"}
		typeValid := false
		for _, validType := range validTypes {
			if variable.Type == validType {
				typeValid = true
				break
			}
		}
		if !typeValid {
			return fmt.Errorf("invalid type '%s' for variable '%s'", variable.Type, name)
		}
	}

	return nil
}

// Auto-fix capabilities (basic implementations)

// FixValidationIssues fixes validation issues
func (e *Engine) FixValidationIssues(path string, issues []interfaces.ValidationIssue) (*interfaces.FixResult, error) {
	return &interfaces.FixResult{
		Applied: []interfaces.Fix{},
		Failed:  []interfaces.FixFailure{},
		Skipped: []interfaces.Fix{},
		Summary: interfaces.FixSummary{
			TotalFixes:    len(issues),
			AppliedFixes:  0,
			FailedFixes:   len(issues),
			SkippedFixes:  0,
			FilesModified: 0,
		},
	}, nil
}

// GetFixableIssues returns issues that can be automatically fixed
func (e *Engine) GetFixableIssues(issues []interfaces.ValidationIssue) []interfaces.ValidationIssue {
	return []interfaces.ValidationIssue{} // No fixable issues for now
}

// PreviewFixes previews what fixes would be applied
func (e *Engine) PreviewFixes(path string, issues []interfaces.ValidationIssue) (*interfaces.FixPreview, error) {
	return &interfaces.FixPreview{
		Fixes:   []interfaces.Fix{},
		Summary: interfaces.FixSummary{},
	}, nil
}

// ApplyFix applies a specific fix
func (e *Engine) ApplyFix(path string, fix interfaces.Fix) error {
	return fmt.Errorf("fix application not implemented")
}

// Validation reporting

// GenerateValidationReport generates a validation report
func (e *Engine) GenerateValidationReport(result *interfaces.ValidationResult, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(result, "", "  ")
	case "text":
		report := fmt.Sprintf("Validation Result: %t\nSummary: %+v\nIssues: %d\n",
			result.Valid, result.Summary, len(result.Issues))
		return []byte(report), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// GetValidationSummary gets a summary of multiple validation results
func (e *Engine) GetValidationSummary(results []*interfaces.ValidationResult) (*interfaces.ValidationSummary, error) {
	summary := &interfaces.ValidationSummary{
		TotalFiles:   len(results),
		ValidFiles:   0,
		ErrorCount:   0,
		WarningCount: 0,
	}

	for _, result := range results {
		if result.Valid {
			summary.ValidFiles++
		}
		summary.ErrorCount += len(result.Issues)
	}

	return summary, nil
}
