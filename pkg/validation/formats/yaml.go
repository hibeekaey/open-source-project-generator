package formats

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// YAMLValidator provides specialized YAML configuration file validation
type YAMLValidator struct {
	schemas map[string]*interfaces.ConfigSchema
}

// NewYAMLValidator creates a new YAML configuration validator
func NewYAMLValidator() *YAMLValidator {
	validator := &YAMLValidator{
		schemas: make(map[string]*interfaces.ConfigSchema),
	}
	validator.initializeYAMLSchemas()
	return validator
}

// ValidateYAMLFile validates a YAML configuration file
func (yv *YAMLValidator) ValidateYAMLFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	// Read and parse YAML
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

	var data interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "syntax",
			Value:    string(content),
			Type:     "syntax_error",
			Message:  fmt.Sprintf("Invalid YAML syntax: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.yaml.syntax",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	// Validate against schema if available
	fileName := filepath.Base(filePath)
	if schema, exists := yv.schemas[fileName]; exists {
		if err := yv.validateAgainstSchema(data, schema, result); err != nil {
			return nil, fmt.Errorf("schema validation failed: %w", err)
		}
	}

	// Perform specific validations based on file name
	switch fileName {
	case "docker-compose.yml", "docker-compose.yaml":
		yv.validateDockerCompose(data, result)
	case ".github/workflows/ci.yml", ".github/workflows/ci.yaml":
		yv.validateGitHubWorkflow(data, result)
	case "k8s.yml", "k8s.yaml", "kubernetes.yml", "kubernetes.yaml":
		yv.validateKubernetesManifest(data, result)
	}

	// Check for common YAML issues
	yv.validateYAMLStructure(string(content), result)

	return result, nil
}

// validateDockerCompose validates Docker Compose configuration
func (yv *YAMLValidator) validateDockerCompose(data interface{}, result *interfaces.ConfigValidationResult) {
	compose, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for version
	result.Summary.TotalProperties++
	if version, exists := compose["version"]; exists {
		result.Summary.ValidProperties++
		if versionStr, ok := version.(string); ok {
			if versionStr == "2" || strings.HasPrefix(versionStr, "2.") {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      "version",
					Value:      versionStr,
					Type:       "deprecated",
					Message:    "Docker Compose version 2.x is deprecated",
					Suggestion: "Consider upgrading to version 3.x",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker_compose.version",
				})
				result.Summary.WarningCount++
			}
		}
	} else {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "version",
			Value:      "",
			Type:       "missing_recommended",
			Message:    "Docker Compose version is recommended",
			Suggestion: "Add version field to specify compose file format",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker_compose.version_missing",
		})
		result.Summary.WarningCount++
	}

	// Check for services
	result.Summary.TotalProperties++
	if services, exists := compose["services"]; exists {
		result.Summary.ValidProperties++

		if servicesMap, ok := services.(map[string]interface{}); ok {
			for serviceName, service := range servicesMap {
				if serviceMap, ok := service.(map[string]interface{}); ok {
					yv.validateDockerService(serviceName, serviceMap, result)
				}
			}
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "services",
			Value:    "",
			Type:     "missing_required",
			Message:  "Docker Compose file must contain services",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "docker_compose.services_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}
}

// validateDockerService validates individual Docker service configuration
func (yv *YAMLValidator) validateDockerService(serviceName string, service map[string]interface{}, result *interfaces.ConfigValidationResult) {
	// Check for privileged mode
	if privileged, exists := service["privileged"]; exists {
		if privilegedBool, ok := privileged.(bool); ok && privilegedBool {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("services.%s.privileged", serviceName),
				Value:      "true",
				Type:       "security",
				Message:    "Privileged mode is a security risk",
				Suggestion: "Avoid using privileged mode unless absolutely necessary",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "docker_compose.privileged",
			})
			result.Summary.WarningCount++
		}
	}

	// Check for image or build
	hasImage := false
	hasBuild := false
	if _, exists := service["image"]; exists {
		hasImage = true
	}
	if _, exists := service["build"]; exists {
		hasBuild = true
	}

	if !hasImage && !hasBuild {
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    fmt.Sprintf("services.%s", serviceName),
			Value:    serviceName,
			Type:     "missing_required",
			Message:  "Service must have either 'image' or 'build' specified",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "docker_compose.image_or_build",
		})
		result.Summary.ErrorCount++
		result.Valid = false
	}

	// Check for restart policy
	if restart, exists := service["restart"]; exists {
		if restartStr, ok := restart.(string); ok {
			validRestartPolicies := []string{"no", "always", "on-failure", "unless-stopped"}
			isValid := false
			for _, policy := range validRestartPolicies {
				if restartStr == policy {
					isValid = true
					break
				}
			}
			if !isValid {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      fmt.Sprintf("services.%s.restart", serviceName),
					Value:      restartStr,
					Type:       "invalid_value",
					Message:    "Invalid restart policy",
					Suggestion: fmt.Sprintf("Use one of: %v", validRestartPolicies),
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker_compose.restart_policy",
				})
				result.Summary.WarningCount++
			}
		}
	}
}

// validateGitHubWorkflow validates GitHub Actions workflow
func (yv *YAMLValidator) validateGitHubWorkflow(data interface{}, result *interfaces.ConfigValidationResult) {
	workflow, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for name
	result.Summary.TotalProperties++
	if _, exists := workflow["name"]; exists {
		result.Summary.ValidProperties++
	} else {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "name",
			Value:      "",
			Type:       "missing_recommended",
			Message:    "Workflow name is recommended for clarity",
			Suggestion: "Add a descriptive name for the workflow",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "github_workflow.name",
		})
		result.Summary.WarningCount++
	}

	// Check for on (triggers)
	result.Summary.TotalProperties++
	if triggers, exists := workflow["on"]; exists {
		result.Summary.ValidProperties++
		yv.validateWorkflowTriggers(triggers, result)
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "on",
			Value:    "",
			Type:     "missing_required",
			Message:  "Workflow triggers ('on') are required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "github_workflow.triggers",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	// Check for jobs
	result.Summary.TotalProperties++
	if jobs, exists := workflow["jobs"]; exists {
		result.Summary.ValidProperties++
		if jobsMap, ok := jobs.(map[string]interface{}); ok {
			yv.validateWorkflowJobs(jobsMap, result)
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "jobs",
			Value:    "",
			Type:     "missing_required",
			Message:  "Workflow must contain jobs",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "github_workflow.jobs_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}
}

// validateWorkflowTriggers validates GitHub workflow triggers
func (yv *YAMLValidator) validateWorkflowTriggers(triggers interface{}, result *interfaces.ConfigValidationResult) {
	switch t := triggers.(type) {
	case string:
		// Single trigger
		yv.validateSingleTrigger(t, result)
	case []interface{}:
		// Array of triggers
		for _, trigger := range t {
			if triggerStr, ok := trigger.(string); ok {
				yv.validateSingleTrigger(triggerStr, result)
			}
		}
	case map[string]interface{}:
		// Complex trigger configuration
		for triggerName := range t {
			yv.validateSingleTrigger(triggerName, result)
		}
	}
}

// validateSingleTrigger validates a single workflow trigger
func (yv *YAMLValidator) validateSingleTrigger(trigger string, result *interfaces.ConfigValidationResult) {
	validTriggers := []string{
		"push", "pull_request", "schedule", "workflow_dispatch",
		"release", "create", "delete", "fork", "gollum", "issue_comment",
		"issues", "label", "milestone", "page_build", "project",
		"project_card", "project_column", "public", "pull_request_review",
		"pull_request_review_comment", "pull_request_target", "registry_package",
		"repository_dispatch", "status", "watch", "workflow_call", "workflow_run",
	}

	isValid := false
	for _, validTrigger := range validTriggers {
		if trigger == validTrigger {
			isValid = true
			break
		}
	}

	if !isValid {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "on",
			Value:      trigger,
			Type:       "unknown_trigger",
			Message:    fmt.Sprintf("Unknown workflow trigger: %s", trigger),
			Suggestion: fmt.Sprintf("Use one of the valid triggers: %v", validTriggers[:5]),
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "github_workflow.valid_triggers",
		})
		result.Summary.WarningCount++
	}
}

// validateWorkflowJobs validates GitHub workflow jobs
func (yv *YAMLValidator) validateWorkflowJobs(jobs map[string]interface{}, result *interfaces.ConfigValidationResult) {
	for jobName, job := range jobs {
		if jobMap, ok := job.(map[string]interface{}); ok {
			// Check for runs-on
			if _, exists := jobMap["runs-on"]; !exists {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    fmt.Sprintf("jobs.%s.runs-on", jobName),
					Value:    "",
					Type:     "missing_required",
					Message:  "Job must specify 'runs-on'",
					Severity: interfaces.ValidationSeverityError,
					Rule:     "github_workflow.runs_on_required",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}

			// Check for steps
			if _, exists := jobMap["steps"]; !exists {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    fmt.Sprintf("jobs.%s.steps", jobName),
					Value:    "",
					Type:     "missing_required",
					Message:  "Job must contain steps",
					Severity: interfaces.ValidationSeverityError,
					Rule:     "github_workflow.steps_required",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}
		}
	}
}

// validateKubernetesManifest validates Kubernetes YAML manifests
func (yv *YAMLValidator) validateKubernetesManifest(data interface{}, result *interfaces.ConfigValidationResult) {
	manifest, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for apiVersion
	result.Summary.TotalProperties++
	if apiVersion, exists := manifest["apiVersion"]; exists {
		result.Summary.ValidProperties++
		if apiVersionStr, ok := apiVersion.(string); ok {
			yv.validateKubernetesAPIVersion(apiVersionStr, result)
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "apiVersion",
			Value:    "",
			Type:     "missing_required",
			Message:  "Kubernetes manifest must specify apiVersion",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "kubernetes.api_version_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	// Check for kind
	result.Summary.TotalProperties++
	if kind, exists := manifest["kind"]; exists {
		result.Summary.ValidProperties++
		if kindStr, ok := kind.(string); ok {
			yv.validateKubernetesKind(kindStr, result)
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "kind",
			Value:    "",
			Type:     "missing_required",
			Message:  "Kubernetes manifest must specify kind",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "kubernetes.kind_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	// Check for metadata
	result.Summary.TotalProperties++
	if metadata, exists := manifest["metadata"]; exists {
		result.Summary.ValidProperties++
		if metadataMap, ok := metadata.(map[string]interface{}); ok {
			yv.validateKubernetesMetadata(metadataMap, result)
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "metadata",
			Value:    "",
			Type:     "missing_required",
			Message:  "Kubernetes manifest must contain metadata",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "kubernetes.metadata_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}
}

// validateKubernetesAPIVersion validates Kubernetes API version
func (yv *YAMLValidator) validateKubernetesAPIVersion(apiVersion string, result *interfaces.ConfigValidationResult) {
	deprecatedVersions := map[string]string{
		"extensions/v1beta1": "apps/v1",
		"apps/v1beta1":       "apps/v1",
		"apps/v1beta2":       "apps/v1",
	}

	if replacement, isDeprecated := deprecatedVersions[apiVersion]; isDeprecated {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "apiVersion",
			Value:      apiVersion,
			Type:       "deprecated",
			Message:    fmt.Sprintf("API version %s is deprecated", apiVersion),
			Suggestion: fmt.Sprintf("Use %s instead", replacement),
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "kubernetes.deprecated_api_version",
		})
		result.Summary.WarningCount++
	}
}

// validateKubernetesKind validates Kubernetes resource kind
func (yv *YAMLValidator) validateKubernetesKind(kind string, result *interfaces.ConfigValidationResult) {
	validKinds := []string{
		"Pod", "Service", "Deployment", "ReplicaSet", "StatefulSet",
		"DaemonSet", "Job", "CronJob", "ConfigMap", "Secret",
		"PersistentVolume", "PersistentVolumeClaim", "Ingress",
		"NetworkPolicy", "ServiceAccount", "Role", "RoleBinding",
		"ClusterRole", "ClusterRoleBinding", "Namespace",
	}

	isValid := false
	for _, validKind := range validKinds {
		if kind == validKind {
			isValid = true
			break
		}
	}

	if !isValid {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "kind",
			Value:      kind,
			Type:       "unknown_kind",
			Message:    fmt.Sprintf("Unknown Kubernetes kind: %s", kind),
			Suggestion: "Verify the kind is correct for your Kubernetes version",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "kubernetes.unknown_kind",
		})
		result.Summary.WarningCount++
	}
}

// validateKubernetesMetadata validates Kubernetes metadata section
func (yv *YAMLValidator) validateKubernetesMetadata(metadata map[string]interface{}, result *interfaces.ConfigValidationResult) {
	// Check for name
	if _, exists := metadata["name"]; !exists {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "metadata.name",
			Value:    "",
			Type:     "missing_required",
			Message:  "Kubernetes resource must have a name",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "kubernetes.name_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	// Validate labels format
	if labels, exists := metadata["labels"]; exists {
		if labelsMap, ok := labels.(map[string]interface{}); ok {
			yv.validateKubernetesLabels(labelsMap, result)
		}
	}
}

// validateKubernetesLabels validates Kubernetes labels
func (yv *YAMLValidator) validateKubernetesLabels(labels map[string]interface{}, result *interfaces.ConfigValidationResult) {
	for key, value := range labels {
		// Validate label key format
		if len(key) > 63 {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("metadata.labels.%s", key),
				Value:      key,
				Type:       "format_warning",
				Message:    "Label key exceeds 63 characters",
				Suggestion: "Use shorter label keys",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "kubernetes.label_key_length",
			})
			result.Summary.WarningCount++
		}

		// Validate label value format
		if valueStr, ok := value.(string); ok {
			if len(valueStr) > 63 {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      fmt.Sprintf("metadata.labels.%s", key),
					Value:      valueStr,
					Type:       "format_warning",
					Message:    "Label value exceeds 63 characters",
					Suggestion: "Use shorter label values",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "kubernetes.label_value_length",
				})
				result.Summary.WarningCount++
			}
		}
	}
}

// validateYAMLStructure validates general YAML structure and common issues
func (yv *YAMLValidator) validateYAMLStructure(content string, result *interfaces.ConfigValidationResult) {
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// Check for tabs (YAML should use spaces)
		if strings.Contains(line, "\t") {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      line,
				Type:       "format_warning",
				Message:    "YAML should use spaces, not tabs for indentation",
				Suggestion: "Replace tabs with spaces",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "yaml.no_tabs",
			})
			result.Summary.WarningCount++
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
				Rule:       "yaml.trailing_whitespace",
			})
			result.Summary.WarningCount++
		}
	}
}

// Helper methods

// validateAgainstSchema validates data against a configuration schema
func (yv *YAMLValidator) validateAgainstSchema(data interface{}, schema *interfaces.ConfigSchema, result *interfaces.ConfigValidationResult) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("data must be an object")
	}

	// Check required properties
	for _, required := range schema.Required {
		result.Summary.TotalProperties++
		if _, exists := dataMap[required]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    required,
				Value:    "",
				Type:     "missing_required",
				Message:  fmt.Sprintf("Required property '%s' is missing", required),
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.required_property",
			})
			result.Summary.ErrorCount++
			result.Summary.MissingRequired++
		} else {
			result.Summary.ValidProperties++
		}
	}

	return nil
}

// initializeYAMLSchemas initializes default YAML configuration schemas
func (yv *YAMLValidator) initializeYAMLSchemas() {
	// Docker Compose schema
	yv.schemas["docker-compose.yml"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Docker Compose Schema",
		Description: "Schema for Docker Compose files",
		Required:    []string{"services"},
		Properties: map[string]interfaces.PropertySchema{
			"version": {
				Type:        "string",
				Description: "Compose file format version",
			},
			"services": {
				Type:        "object",
				Description: "Service definitions",
			},
		},
	}

	// GitHub Workflow schema
	yv.schemas[".github/workflows/ci.yml"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "GitHub Workflow Schema",
		Description: "Schema for GitHub Actions workflows",
		Required:    []string{"on", "jobs"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Workflow name",
			},
			"on": {
				Type:        "object",
				Description: "Workflow triggers",
			},
			"jobs": {
				Type:        "object",
				Description: "Job definitions",
			},
		},
	}
}
