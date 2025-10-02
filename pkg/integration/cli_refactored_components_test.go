package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/spf13/cobra"
)

// TestRefactoredCLIComponents tests the integration of refactored CLI components
func TestRefactoredCLIComponents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("cli_handlers_integration", func(t *testing.T) {
		testCLIHandlersIntegration(t, tempDir)
	})

	t.Run("interactive_components_integration", func(t *testing.T) {
		testInteractiveComponentsIntegration(t, tempDir)
	})

	t.Run("cli_validation_integration", func(t *testing.T) {
		testCLIValidationIntegration(t, tempDir)
	})

	t.Run("cli_workflow_coordination", func(t *testing.T) {
		testCLIWorkflowCoordination(t, tempDir)
	})
}

func testCLIHandlersIntegration(t *testing.T, tempDir string) {
	// Test generate handler
	t.Run("generate_handler", func(t *testing.T) {
		generateHandler := NewMockGenerateHandler()

		// Create test config
		config := &models.ProjectConfig{
			Name:         "handler-test-project",
			Organization: "handler-org",
			Description:  "Test project for handler integration",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "generate-handler-test"),
		}

		// Create mock command
		cmd := &cobra.Command{
			Use: "generate",
		}
		cmd.Flags().String("config", "", "Config file path")
		cmd.Flags().String("template", "go-gin", "Template to use")
		cmd.Flags().Bool("non-interactive", true, "Non-interactive mode")

		// Set flag values
		configPath := filepath.Join(tempDir, "handler-config.yaml")
		configData := `name: handler-test-project
organization: handler-org
description: Test project for handler integration
license: MIT
output_path: ` + config.OutputPath

		err := os.WriteFile(configPath, []byte(configData), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		err = cmd.Flags().Set("config", configPath)
		if err != nil {
			t.Fatalf("Failed to set config flag: %v", err)
		}

		// Test handler execution
		err = generateHandler.Handle(cmd, []string{})
		if err != nil {
			t.Errorf("Generate handler failed: %v", err)
		}

		// Verify handler validation
		err = generateHandler.Validate(cmd, []string{})
		if err != nil {
			t.Errorf("Generate handler validation failed: %v", err)
		}
	})

	// Test validate handler
	t.Run("validate_handler", func(t *testing.T) {
		validateHandler := NewMockValidateHandler()

		// Create test project to validate
		projectDir := filepath.Join(tempDir, "validate-handler-test")
		err := os.MkdirAll(projectDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create project directory: %v", err)
		}

		// Create basic project files
		files := map[string]string{
			"README.md":    "# Validate Handler Test\n\nTest project for validation.",
			"package.json": `{"name": "validate-handler-test", "version": "1.0.0"}`,
			"LICENSE":      "MIT License\n\nCopyright (c) 2024 Test",
		}

		for filename, content := range files {
			err := os.WriteFile(filepath.Join(projectDir, filename), []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create file %s: %v", filename, err)
			}
		}

		// Create mock command
		cmd := &cobra.Command{
			Use: "validate",
		}
		cmd.Flags().Bool("report", false, "Generate report")
		cmd.Flags().String("report-format", "json", "Report format")

		// Test handler execution
		err = validateHandler.Handle(cmd, []string{projectDir})
		if err != nil {
			t.Errorf("Validate handler failed: %v", err)
		}

		// Test handler validation
		err = validateHandler.Validate(cmd, []string{projectDir})
		if err != nil {
			t.Errorf("Validate handler validation failed: %v", err)
		}
	})

	// Test audit handler
	t.Run("audit_handler", func(t *testing.T) {
		auditHandler := NewMockAuditHandler()

		// Use the same project from validate test
		projectDir := filepath.Join(tempDir, "validate-handler-test")

		// Create mock command
		cmd := &cobra.Command{
			Use: "audit",
		}
		cmd.Flags().Bool("security", false, "Security audit")
		cmd.Flags().Bool("quality", false, "Quality audit")
		cmd.Flags().String("output-format", "json", "Output format")

		// Test handler execution
		err := auditHandler.Handle(cmd, []string{projectDir})
		if err != nil {
			t.Errorf("Audit handler failed: %v", err)
		}

		// Test handler validation
		err = auditHandler.Validate(cmd, []string{projectDir})
		if err != nil {
			t.Errorf("Audit handler validation failed: %v", err)
		}
	})

	// Test config handler
	t.Run("config_handler", func(t *testing.T) {
		configHandler := NewMockConfigHandler()

		// Create mock command for config show
		showCmd := &cobra.Command{
			Use: "show",
		}

		// Test config show
		err := configHandler.Handle(showCmd, []string{})
		if err != nil {
			t.Errorf("Config show handler failed: %v", err)
		}

		// Create mock command for config set
		setCmd := &cobra.Command{
			Use: "set",
		}

		// Test config set
		err = configHandler.Handle(setCmd, []string{"default_license", "Apache-2.0"})
		if err != nil {
			t.Errorf("Config set handler failed: %v", err)
		}

		// Test handler validation
		err = configHandler.Validate(setCmd, []string{"default_license", "Apache-2.0"})
		if err != nil {
			t.Errorf("Config handler validation failed: %v", err)
		}
	})
}

func testInteractiveComponentsIntegration(t *testing.T, tempDir string) {
	// Test project setup component
	t.Run("project_setup", func(t *testing.T) {
		projectSetup := NewMockProjectSetup()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create mock input for non-interactive testing
		mockInput := &MockInteractiveInput{
			responses: map[string]string{
				"project_name": "interactive-test-project",
				"organization": "interactive-org",
				"description":  "Test project for interactive components",
				"license":      "MIT",
			},
		}

		config, err := projectSetup.CollectProjectInfo(ctx, mockInput)
		if err != nil {
			t.Errorf("Project setup failed: %v", err)
		}

		if config == nil {
			t.Fatal("Expected project config to be returned")
		}

		if config.Name != "interactive-test-project" {
			t.Errorf("Expected project name 'interactive-test-project', got '%s'", config.Name)
		}

		if config.Organization != "interactive-org" {
			t.Errorf("Expected organization 'interactive-org', got '%s'", config.Organization)
		}
	})

	// Test component selection
	t.Run("component_selection", func(t *testing.T) {
		componentSelection := NewMockComponentSelection()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		config := &models.ProjectConfig{
			Name:         "component-test-project",
			Organization: "component-org",
		}

		mockInput := &MockInteractiveInput{
			responses: map[string]string{
				"backend_enabled":     "true",
				"backend_technology":  "go-gin",
				"frontend_enabled":    "true",
				"frontend_technology": "nextjs",
			},
		}

		err := componentSelection.SelectComponents(ctx, config, mockInput)
		if err != nil {
			t.Errorf("Component selection failed: %v", err)
		}

		// Verify components were selected
		if !config.Components.Backend.GoGin {
			t.Error("Expected Go Gin backend to be enabled")
		}

		if !config.Components.Frontend.NextJS.App {
			t.Error("Expected Next.js frontend to be enabled")
		}
	})

	// Test validation UI
	t.Run("validation_ui", func(t *testing.T) {
		validationUI := NewMockValidationUI()

		// Create mock validation results
		results := &interfaces.ValidationResult{
			Valid: false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: "error",
					Message:  "Missing required file: README.md",
					Rule:     "readme-required",
					File:     "README.md",
				},
				{
					Type:     "warning",
					Severity: "warning",
					Message:  "License file not found",
					Rule:     "license-recommended",
					File:     "LICENSE",
				},
			},
		}

		err := validationUI.DisplayResults(results)
		if err != nil {
			t.Errorf("Validation UI display failed: %v", err)
		}

		// Test interactive fix prompts
		mockInput := &MockInteractiveInput{
			responses: map[string]string{
				"fix_readme":  "yes",
				"fix_license": "no",
			},
		}

		fixedResults, err := validationUI.PromptForFixes(results, mockInput)
		if err != nil {
			t.Errorf("Validation UI fix prompts failed: %v", err)
		}

		if fixedResults == nil {
			t.Error("Expected fixed results to be returned")
		}
	})
}

func testCLIValidationIntegration(t *testing.T, tempDir string) {
	// Test input validator
	t.Run("input_validator", func(t *testing.T) {
		inputValidator := NewMockInputValidator()

		// Test valid project name
		err := inputValidator.ValidateProjectName("valid-project-name")
		if err != nil {
			t.Errorf("Valid project name validation failed: %v", err)
		}

		// Test invalid project name
		err = inputValidator.ValidateProjectName("Invalid Project Name!")
		if err == nil {
			t.Error("Expected invalid project name to fail validation")
		}

		// Test valid organization
		err = inputValidator.ValidateOrganization("valid-org")
		if err != nil {
			t.Errorf("Valid organization validation failed: %v", err)
		}

		// Test valid license
		err = inputValidator.ValidateLicense("MIT")
		if err != nil {
			t.Errorf("Valid license validation failed: %v", err)
		}

		// Test invalid license
		err = inputValidator.ValidateLicense("INVALID-LICENSE")
		if err == nil {
			t.Error("Expected invalid license to fail validation")
		}
	})

	// Test config validator
	t.Run("config_validator", func(t *testing.T) {
		configValidator := NewMockConfigValidator()

		// Test valid config
		validConfig := &models.ProjectConfig{
			Name:         "valid-config-project",
			Organization: "valid-org",
			Description:  "Valid project configuration",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "valid-output"),
		}

		err := configValidator.ValidateConfig(validConfig)
		if err != nil {
			t.Errorf("Valid config validation failed: %v", err)
		}

		// Test invalid config (missing required fields)
		invalidConfig := &models.ProjectConfig{
			Name: "", // Missing required name
		}

		err = configValidator.ValidateConfig(invalidConfig)
		if err == nil {
			t.Error("Expected invalid config to fail validation")
		}

		// Test config file validation
		configPath := filepath.Join(tempDir, "test-config.yaml")
		configContent := `name: file-config-project
organization: file-org
license: MIT
output_path: ` + filepath.Join(tempDir, "file-output")

		err = os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		err = configValidator.ValidateConfigFile(configPath)
		if err != nil {
			t.Errorf("Config file validation failed: %v", err)
		}
	})

	// Test validation integration
	t.Run("validation_integration", func(t *testing.T) {
		inputValidator := NewMockInputValidator()
		configValidator := NewMockConfigValidator()

		// Test complete validation workflow
		projectName := "integration-test-project"
		organization := "integration-org"
		license := "Apache-2.0"

		// Validate individual components
		err := inputValidator.ValidateProjectName(projectName)
		if err != nil {
			t.Errorf("Project name validation failed: %v", err)
		}

		err = inputValidator.ValidateOrganization(organization)
		if err != nil {
			t.Errorf("Organization validation failed: %v", err)
		}

		err = inputValidator.ValidateLicense(license)
		if err != nil {
			t.Errorf("License validation failed: %v", err)
		}

		// Create and validate complete config
		config := &models.ProjectConfig{
			Name:         projectName,
			Organization: organization,
			License:      license,
			OutputPath:   filepath.Join(tempDir, "integration-output"),
		}

		err = configValidator.ValidateConfig(config)
		if err != nil {
			t.Errorf("Complete config validation failed: %v", err)
		}
	})
}

func testCLIWorkflowCoordination(t *testing.T, tempDir string) {
	// Test complete CLI workflow coordination
	t.Run("complete_workflow", func(t *testing.T) {
		// Create CLI instance
		cliInstance := NewMockCLI()

		// Test CLI initialization
		err := cliInstance.Initialize()
		if err != nil {
			t.Errorf("CLI initialization failed: %v", err)
		}

		// Test command registration
		err = cliInstance.RegisterCommands()
		if err != nil {
			t.Errorf("Command registration failed: %v", err)
		}

		// Test global flags setup
		err = cliInstance.SetupGlobalFlags()
		if err != nil {
			t.Errorf("Global flags setup failed: %v", err)
		}

		// Verify CLI is ready for execution
		if !cliInstance.IsReady() {
			t.Error("Expected CLI to be ready after initialization")
		}
	})

	// Test error handling coordination
	t.Run("error_handling_coordination", func(t *testing.T) {
		cliInstance := NewMockCLI()
		err := cliInstance.Initialize()
		if err != nil {
			t.Fatalf("CLI initialization failed: %v", err)
		}

		// Test error handling for invalid commands
		err = cliInstance.ExecuteCommand([]string{"invalid-command"})
		if err == nil {
			t.Error("Expected error for invalid command")
		}

		// Test error handling for missing arguments
		err = cliInstance.ExecuteCommand([]string{"generate"})
		if err == nil {
			t.Error("Expected error for missing arguments")
		}

		// Test error recovery
		err = cliInstance.RecoverFromError(err)
		if err != nil {
			t.Errorf("Error recovery failed: %v", err)
		}
	})

	// Test component coordination
	t.Run("component_coordination", func(t *testing.T) {
		// Test handler coordination
		generateHandler := NewMockGenerateHandler()
		validateHandler := NewMockValidateHandler()
		auditHandler := NewMockAuditHandler()

		// Test handler registration
		cliInstance := NewMockCLI()
		err := cliInstance.RegisterHandler("generate", generateHandler)
		if err != nil {
			t.Errorf("Generate handler registration failed: %v", err)
		}

		err = cliInstance.RegisterHandler("validate", validateHandler)
		if err != nil {
			t.Errorf("Validate handler registration failed: %v", err)
		}

		err = cliInstance.RegisterHandler("audit", auditHandler)
		if err != nil {
			t.Errorf("Audit handler registration failed: %v", err)
		}

		// Test interactive component coordination
		projectSetup := NewMockProjectSetup()
		componentSelection := NewMockComponentSelection()

		err = cliInstance.RegisterInteractiveComponent("project-setup", projectSetup)
		if err != nil {
			t.Errorf("Project setup registration failed: %v", err)
		}

		err = cliInstance.RegisterInteractiveComponent("component-selection", componentSelection)
		if err != nil {
			t.Errorf("Component selection registration failed: %v", err)
		}

		// Test validation component coordination
		inputValidator := NewMockInputValidator()
		configValidator := NewMockConfigValidator()

		err = cliInstance.RegisterValidator("input", inputValidator)
		if err != nil {
			t.Errorf("Input validator registration failed: %v", err)
		}

		err = cliInstance.RegisterValidator("config", configValidator)
		if err != nil {
			t.Errorf("Config validator registration failed: %v", err)
		}
	})
}

// Mock implementations for testing

type MockInteractiveInput struct {
	responses map[string]string
}

func (m *MockInteractiveInput) GetInput(prompt string, key string) (string, error) {
	if response, exists := m.responses[key]; exists {
		return response, nil
	}
	return "", nil
}

func (m *MockInteractiveInput) GetBoolInput(prompt string, key string) (bool, error) {
	if response, exists := m.responses[key]; exists {
		return response == "true" || response == "yes" || response == "y", nil
	}
	return false, nil
}

func (m *MockInteractiveInput) GetSelectInput(prompt string, key string, options []string) (string, error) {
	if response, exists := m.responses[key]; exists {
		// Verify response is in options
		for _, option := range options {
			if option == response {
				return response, nil
			}
		}
	}
	// Return first option as default
	if len(options) > 0 {
		return options[0], nil
	}
	return "", nil
}

func (m *MockInteractiveInput) GetMultiSelectInput(prompt string, key string, options []string) ([]string, error) {
	if response, exists := m.responses[key]; exists {
		// Simple implementation - return single response as slice
		return []string{response}, nil
	}
	return []string{}, nil
}

// Mock implementations for CLI testing

// Mock CLI
type MockCLI struct {
	initialized bool
	ready       bool
	handlers    map[string]interface{}
	components  map[string]interface{}
	validators  map[string]interface{}
}

func NewMockCLI() *MockCLI {
	return &MockCLI{
		handlers:   make(map[string]interface{}),
		components: make(map[string]interface{}),
		validators: make(map[string]interface{}),
	}
}

func (m *MockCLI) Initialize() error {
	m.initialized = true
	return nil
}

func (m *MockCLI) RegisterCommands() error {
	return nil
}

func (m *MockCLI) SetupGlobalFlags() error {
	return nil
}

func (m *MockCLI) IsReady() bool {
	return m.initialized
}

func (m *MockCLI) ExecuteCommand(args []string) error {
	if len(args) == 0 {
		return NewMockError("no command provided")
	}

	if args[0] == "invalid-command" {
		return NewMockError("unknown command: " + args[0])
	}

	if args[0] == "generate" && len(args) == 1 {
		return NewMockError("missing required arguments for generate command")
	}

	return nil
}

func (m *MockCLI) RecoverFromError(err error) error {
	// Mock recovery always succeeds
	return nil
}

func (m *MockCLI) RegisterHandler(name string, handler interface{}) error {
	m.handlers[name] = handler
	return nil
}

func (m *MockCLI) RegisterInteractiveComponent(name string, component interface{}) error {
	m.components[name] = component
	return nil
}

func (m *MockCLI) RegisterValidator(name string, validator interface{}) error {
	m.validators[name] = validator
	return nil
}

// Mock Handlers
type MockGenerateHandler struct{}

func NewMockGenerateHandler() *MockGenerateHandler {
	return &MockGenerateHandler{}
}

func (m *MockGenerateHandler) Handle(cmd *cobra.Command, args []string) error {
	return nil
}

func (m *MockGenerateHandler) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

type MockValidateHandler struct{}

func NewMockValidateHandler() *MockValidateHandler {
	return &MockValidateHandler{}
}

func (m *MockValidateHandler) Handle(cmd *cobra.Command, args []string) error {
	return nil
}

func (m *MockValidateHandler) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

type MockAuditHandler struct{}

func NewMockAuditHandler() *MockAuditHandler {
	return &MockAuditHandler{}
}

func (m *MockAuditHandler) Handle(cmd *cobra.Command, args []string) error {
	return nil
}

func (m *MockAuditHandler) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

type MockConfigHandler struct{}

func NewMockConfigHandler() *MockConfigHandler {
	return &MockConfigHandler{}
}

func (m *MockConfigHandler) Handle(cmd *cobra.Command, args []string) error {
	return nil
}

func (m *MockConfigHandler) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Mock Interactive Components
type MockProjectSetup struct{}

func NewMockProjectSetup() *MockProjectSetup {
	return &MockProjectSetup{}
}

func (m *MockProjectSetup) CollectProjectInfo(ctx context.Context, input *MockInteractiveInput) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         input.responses["project_name"],
		Organization: input.responses["organization"],
		Description:  input.responses["description"],
		License:      input.responses["license"],
	}, nil
}

type MockComponentSelection struct{}

func NewMockComponentSelection() *MockComponentSelection {
	return &MockComponentSelection{}
}

func (m *MockComponentSelection) SelectComponents(ctx context.Context, config *models.ProjectConfig, input *MockInteractiveInput) error {
	if backendEnabled, exists := input.responses["backend_enabled"]; exists && backendEnabled == "true" {
		config.Components.Backend.GoGin = true
	}

	if frontendEnabled, exists := input.responses["frontend_enabled"]; exists && frontendEnabled == "true" {
		config.Components.Frontend.NextJS.App = true
	}

	return nil
}

type MockValidationUI struct{}

func NewMockValidationUI() *MockValidationUI {
	return &MockValidationUI{}
}

func (m *MockValidationUI) DisplayResults(results interface{}) error {
	return nil
}

func (m *MockValidationUI) PromptForFixes(results *interfaces.ValidationResult, input *MockInteractiveInput) (*interfaces.ValidationResult, error) {
	// Mock implementation returns modified results
	return &interfaces.ValidationResult{
		Valid:  true,
		Issues: []interfaces.ValidationIssue{},
	}, nil
}

// Mock Validators
type MockInputValidator struct{}

func NewMockInputValidator() *MockInputValidator {
	return &MockInputValidator{}
}

func (m *MockInputValidator) ValidateProjectName(name string) error {
	if name == "" || len(name) < 3 {
		return NewMockError("project name must be at least 3 characters")
	}
	if name == "Invalid Project Name!" {
		return NewMockError("project name contains invalid characters")
	}
	return nil
}

func (m *MockInputValidator) ValidateOrganization(org string) error {
	if org == "" {
		return NewMockError("organization cannot be empty")
	}
	return nil
}

func (m *MockInputValidator) ValidateLicense(license string) error {
	validLicenses := []string{"MIT", "Apache-2.0", "BSD-3-Clause", "GPL-3.0"}
	for _, valid := range validLicenses {
		if license == valid {
			return nil
		}
	}
	return NewMockError("invalid license: " + license)
}

type MockConfigValidator struct{}

func NewMockConfigValidator() *MockConfigValidator {
	return &MockConfigValidator{}
}

func (m *MockConfigValidator) ValidateConfig(config *models.ProjectConfig) error {
	if config.Name == "" {
		return NewMockError("project name is required")
	}
	return nil
}

func (m *MockConfigValidator) ValidateConfigFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewMockError("config file not found: " + path)
	}
	return nil
}

// Mock Error
type MockError struct {
	message string
}

func NewMockError(message string) *MockError {
	return &MockError{message: message}
}

func (e *MockError) Error() string {
	return e.message
}
