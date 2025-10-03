package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MockLogger implements a simple logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, args ...interface{})                               {}
func (m *MockLogger) Info(msg string, args ...interface{})                                {}
func (m *MockLogger) Warn(msg string, args ...interface{})                                {}
func (m *MockLogger) Error(msg string, args ...interface{})                               {}
func (m *MockLogger) Fatal(msg string, args ...interface{})                               {}
func (m *MockLogger) DebugWithFields(msg string, fields map[string]interface{})           {}
func (m *MockLogger) InfoWithFields(msg string, fields map[string]interface{})            {}
func (m *MockLogger) WarnWithFields(msg string, fields map[string]interface{})            {}
func (m *MockLogger) ErrorWithFields(msg string, fields map[string]interface{})           {}
func (m *MockLogger) FatalWithFields(msg string, fields map[string]interface{})           {}
func (m *MockLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {}
func (m *MockLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	return nil
}
func (m *MockLogger) LogOperationStart(operation string, fields map[string]interface{}) {}
func (m *MockLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
}
func (m *MockLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {}
func (m *MockLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
}
func (m *MockLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
}
func (m *MockLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {}
func (m *MockLogger) LogMemoryUsage(operation string)                                        {}
func (m *MockLogger) SetLevel(level int)                                                     {}
func (m *MockLogger) GetLevel() int                                                          { return 0 }
func (m *MockLogger) SetJSONOutput(enable bool)                                              {}
func (m *MockLogger) SetCallerInfo(enable bool)                                              {}
func (m *MockLogger) IsDebugEnabled() bool                                                   { return false }
func (m *MockLogger) IsInfoEnabled() bool                                                    { return true }
func (m *MockLogger) WithComponent(component string) interfaces.Logger                       { return m }
func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &MockLoggerContext{}
}
func (m *MockLogger) GetLogDir() string { return "" }
func (m *MockLogger) GetRecentEntries(limit int) []interfaces.LogEntry {
	return nil
}
func (m *MockLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (m *MockLogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (m *MockLogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }
func (m *MockLogger) Close() error                                { return nil }

type MockLoggerContext struct{}

func (m *MockLoggerContext) Debug(msg string, args ...interface{}) {}
func (m *MockLoggerContext) Info(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Warn(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Error(msg string, args ...interface{}) {}
func (m *MockLoggerContext) ErrorWithError(msg string, err error)  {}

func TestNewManager(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	logger := &MockLogger{}

	// Create manager
	manager, err := NewManager(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	if manager == nil {
		t.Fatal("Manager is nil")
	}
}

func TestManagerLoadDefaults(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	logger := &MockLogger{}

	// Create manager
	manager, err := NewManager(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Load defaults
	config, err := manager.LoadDefaults()
	if err != nil {
		t.Fatalf("Failed to load defaults: %v", err)
	}

	if config == nil {
		t.Fatal("Config is nil")
	}

	if config.Name != "my-project" {
		t.Errorf("Expected name 'my-project', got '%s'", config.Name)
	}

	if config.Organization != "my-org" {
		t.Errorf("Expected organization 'my-org', got '%s'", config.Organization)
	}

	if config.License != "MIT" {
		t.Errorf("Expected license 'MIT', got '%s'", config.License)
	}
}

func TestManagerValidateConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	logger := &MockLogger{}

	// Create manager
	manager, err := NewManager(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test valid configuration
	validConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		License:      "MIT",
		OutputPath:   "./output",
	}

	err = manager.ValidateConfig(validConfig)
	if err != nil {
		t.Errorf("Valid config should not have errors: %v", err)
	}

	// Test invalid configuration (missing required fields)
	invalidConfig := &models.ProjectConfig{
		License:    "MIT",
		OutputPath: "./output",
	}

	err = manager.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Invalid config should have errors")
	}
}

func TestManagerSaveAndLoadConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	logger := &MockLogger{}

	// Create manager
	manager, err := NewManager(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create test configuration
	testConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "A test project",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/test/repo",
		OutputPath:   "./output",
		Features:     []string{"feature1", "feature2"},
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"react": "18.2.0",
				"next":  "13.4.0",
			},
		},
	}

	// Test YAML format
	yamlPath := filepath.Join(tempDir, "test-config.yaml")
	err = manager.SaveConfig(testConfig, yamlPath)
	if err != nil {
		t.Fatalf("Failed to save YAML config: %v", err)
	}

	// Load YAML config
	loadedConfig, err := manager.LoadConfig(yamlPath)
	if err != nil {
		t.Fatalf("Failed to load YAML config: %v", err)
	}

	// Verify loaded config
	if loadedConfig.Name != testConfig.Name {
		t.Errorf("Expected name '%s', got '%s'", testConfig.Name, loadedConfig.Name)
	}

	if loadedConfig.Organization != testConfig.Organization {
		t.Errorf("Expected organization '%s', got '%s'", testConfig.Organization, loadedConfig.Organization)
	}

	// Test JSON format
	jsonPath := filepath.Join(tempDir, "test-config.json")
	err = manager.SaveConfig(testConfig, jsonPath)
	if err != nil {
		t.Fatalf("Failed to save JSON config: %v", err)
	}

	// Load JSON config
	loadedJSONConfig, err := manager.LoadConfig(jsonPath)
	if err != nil {
		t.Fatalf("Failed to load JSON config: %v", err)
	}

	// Verify loaded JSON config
	if loadedJSONConfig.Name != testConfig.Name {
		t.Errorf("Expected name '%s', got '%s'", testConfig.Name, loadedJSONConfig.Name)
	}
}

func TestConfigValidator(t *testing.T) {
	schema := createConfigSchema()
	validator := &ConfigValidator{schema: schema}

	// Test valid project config
	validConfig := &models.ProjectConfig{
		Name:         "valid-project",
		Organization: "valid-org",
		License:      "MIT",
		Email:        "test@example.com",
		OutputPath:   "./output",
	}

	result, err := validator.ValidateProjectConfig(validConfig)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Valid config should pass validation. Errors: %v", result.Errors)
	}

	// Test invalid project config
	invalidConfig := &models.ProjectConfig{
		Name:         "", // Missing required field
		Organization: "valid-org",
		Email:        "invalid-email", // Invalid email format
	}

	result, err = validator.ValidateProjectConfig(invalidConfig)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if result.Valid {
		t.Error("Invalid config should fail validation")
	}

	if len(result.Errors) == 0 {
		t.Error("Invalid config should have errors")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test isValidProjectName
	validNames := []string{"my-project", "test_app", "Project123", "a"}
	for _, name := range validNames {
		if !isValidProjectName(name) {
			t.Errorf("'%s' should be a valid project name", name)
		}
	}

	invalidNames := []string{"", "my project", "test@app", "very-long-project-name-that-exceeds-the-maximum-allowed-length-of-one-hundred-characters-and-should-be-rejected"}
	for _, name := range invalidNames {
		if isValidProjectName(name) {
			t.Errorf("'%s' should be an invalid project name", name)
		}
	}

	// Test isValidEmail
	validEmails := []string{"test@example.com", "user.name@domain.org", "a@b.co"}
	for _, email := range validEmails {
		if !isValidEmail(email) {
			t.Errorf("'%s' should be a valid email", email)
		}
	}

	invalidEmails := []string{"", "invalid", "test@", "@domain.com", "test@domain", "test.domain.com"}
	for _, email := range invalidEmails {
		if isValidEmail(email) {
			t.Errorf("'%s' should be an invalid email", email)
		}
	}
}
