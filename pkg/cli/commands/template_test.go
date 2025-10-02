package commands

import (
	"errors"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTemplateCLI is a mock implementation of TemplateCLI for testing.
type MockTemplateCLI struct {
	mock.Mock
}

func (m *MockTemplateCLI) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	args := m.Called(filter)
	return args.Get(0).([]interfaces.TemplateInfo), args.Error(1)
}

func (m *MockTemplateCLI) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateInfo), args.Error(1)
}

func (m *MockTemplateCLI) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateValidationResult), args.Error(1)
}

func (m *MockTemplateCLI) VerboseOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) DebugOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) QuietOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) ErrorOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) WarningOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) SuccessOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockTemplateCLI) Error(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) Warning(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) Info(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) Success(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) Highlight(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) Dim(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockTemplateCLI) IsQuietMode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockTemplateCLI) OutputMachineReadable(data interface{}, format string) error {
	args := m.Called(data, format)
	return args.Error(0)
}

func (m *MockTemplateCLI) CreateTemplateError(message string, templateName string) error {
	args := m.Called(message, templateName)
	return args.Error(0)
}

func (m *MockTemplateCLI) OutputSuccess(message string, data interface{}, operation string, args []string) error {
	mockArgs := m.Called(message, data, operation, args)
	return mockArgs.Error(0)
}

func (m *MockTemplateCLI) IsNonInteractiveMode(cmd *cobra.Command) bool {
	args := m.Called(cmd)
	return args.Bool(0)
}

// Test data
func createTestTemplateInfo() *interfaces.TemplateInfo {
	return &interfaces.TemplateInfo{
		Name:         "test-template",
		DisplayName:  "Test Template",
		Description:  "A test template for unit testing",
		Category:     "backend",
		Technology:   "go",
		Version:      "1.0.0",
		Tags:         []string{"test", "example"},
		Dependencies: []string{"go", "gin"},
		Metadata: interfaces.TemplateMetadata{
			Author:     "Test Author",
			License:    "MIT",
			Repository: "https://github.com/test/template",
			Homepage:   "https://test.com",
			Keywords:   []string{"test", "template"},
			Created:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Updated:    time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			Variables: map[string]string{
				"ProjectName": "Name of the project",
				"Author":      "Project author",
			},
		},
	}
}

func createTestTemplates() []interfaces.TemplateInfo {
	return []interfaces.TemplateInfo{
		{
			Name:        "go-gin",
			DisplayName: "Go Gin API",
			Description: "REST API with Go and Gin framework",
			Category:    "backend",
			Technology:  "go",
			Version:     "1.0.0",
			Tags:        []string{"api", "rest"},
		},
		{
			Name:        "nextjs-app",
			DisplayName: "Next.js Application",
			Description: "Modern React application with Next.js",
			Category:    "frontend",
			Technology:  "nodejs",
			Version:     "1.0.0",
			Tags:        []string{"react", "nextjs"},
		},
	}
}

func createTestValidationResult(valid bool) *interfaces.TemplateValidationResult {
	result := &interfaces.TemplateValidationResult{
		Valid: valid,
		Summary: interfaces.ValidationSummary{
			TotalFiles: 10,
		},
		Issues:   []interfaces.ValidationIssue{},
		Warnings: []interfaces.ValidationIssue{},
	}

	if !valid {
		result.Issues = []interfaces.ValidationIssue{
			{
				Severity: "error",
				Message:  "Missing template.yaml file",
				File:     "template.yaml",
				Line:     0,
				Column:   0,
			},
		}
		result.Warnings = []interfaces.ValidationIssue{
			{
				Severity: "warning",
				Message:  "No README.md found",
				File:     "README.md",
				Line:     0,
				Column:   0,
			},
		}
	}

	return result
}

func TestNewTemplateCommands(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	assert.NotNil(t, tc)
	assert.Equal(t, mockCLI, tc.cli)
}

func TestTemplateCommands_ExecuteList_Success(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	templates := createTestTemplates()
	filter := interfaces.TemplateFilter{
		Tags: []string{}, // Explicitly set empty slice instead of nil
	}

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ListTemplates", filter).Return(templates, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("OutputSuccess", mock.AnythingOfType("string"), mock.Anything, "list-templates", []string{}).Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String("category", "", "")
	cmd.Flags().String("technology", "", "")
	cmd.Flags().StringSlice("tags", []string{}, "")
	cmd.Flags().String("search", "", "")
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	err := tc.ExecuteList(cmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteList_WithSearch(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	templates := createTestTemplates()
	filter := interfaces.TemplateFilter{
		Tags: []string{}, // Explicitly set empty slice instead of nil
	}

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ListTemplates", filter).Return(templates, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("OutputSuccess", mock.AnythingOfType("string"), mock.Anything, "list-templates", []string{}).Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String("category", "", "")
	cmd.Flags().String("technology", "", "")
	cmd.Flags().StringSlice("tags", []string{}, "")
	cmd.Flags().String("search", "", "")
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	// Set search flag
	if err := cmd.Flags().Set("search", "go"); err != nil {
		t.Fatalf("Failed to set search flag: %v", err)
	}

	err := tc.ExecuteList(cmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteList_NoTemplatesFound(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	filter := interfaces.TemplateFilter{
		Tags: []string{}, // Explicitly set empty slice instead of nil
	}

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ListTemplates", filter).Return([]interfaces.TemplateInfo{}, nil)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("OutputSuccess", "No templates found", mock.Anything, "list-templates", []string{}).Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String("category", "", "")
	cmd.Flags().String("technology", "", "")
	cmd.Flags().StringSlice("tags", []string{}, "")
	cmd.Flags().String("search", "", "")
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	err := tc.ExecuteList(cmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteList_Error(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	filter := interfaces.TemplateFilter{
		Tags: []string{}, // Explicitly set empty slice instead of nil
	}
	expectedError := errors.New("template service unavailable")

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ListTemplates", filter).Return([]interfaces.TemplateInfo{}, expectedError)
	mockCLI.On("CreateTemplateError", "failed to list templates", "").Return(expectedError)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String("category", "", "")
	cmd.Flags().String("technology", "", "")
	cmd.Flags().StringSlice("tags", []string{}, "")
	cmd.Flags().String("search", "", "")
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	err := tc.ExecuteList(cmd, []string{})

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteList_MachineReadable(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	templates := createTestTemplates()
	filter := interfaces.TemplateFilter{
		Tags: []string{}, // Explicitly set empty slice instead of nil
	}

	// Setup expectations
	mockCLI.On("ListTemplates", filter).Return(templates, nil)
	mockCLI.On("OutputMachineReadable", mock.Anything, "json").Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().String("category", "", "")
	cmd.Flags().String("technology", "", "")
	cmd.Flags().StringSlice("tags", []string{}, "")
	cmd.Flags().String("search", "", "")
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("non-interactive", true, "")
	cmd.Flags().String("output-format", "json", "")

	// Set the flags to trigger machine-readable mode
	cmd.Flags().Set("non-interactive", "true")
	cmd.Flags().Set("output-format", "json")

	err := tc.ExecuteList(cmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteInfo_Success(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	templateInfo := createTestTemplateInfo()

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("GetTemplateInfo", "test-template").Return(templateInfo, nil)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("variables", false, "")
	cmd.Flags().Bool("dependencies", false, "")
	cmd.Flags().Bool("compatibility", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	err := tc.ExecuteInfo(cmd, []string{"test-template"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteInfo_NotFound(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	expectedError := errors.New("template not found")

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("GetTemplateInfo", "nonexistent").Return((*interfaces.TemplateInfo)(nil), expectedError)
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("variables", false, "")
	cmd.Flags().Bool("dependencies", false, "")
	cmd.Flags().Bool("compatibility", false, "")
	cmd.Flags().Bool("non-interactive", false, "")
	cmd.Flags().String("output-format", "text", "")

	err := tc.ExecuteInfo(cmd, []string{"nonexistent"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "🚫")
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteInfo_MachineReadable(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	templateInfo := createTestTemplateInfo()

	// Setup expectations
	mockCLI.On("GetTemplateInfo", "test-template").Return(templateInfo, nil)
	mockCLI.On("OutputMachineReadable", templateInfo, "json").Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("variables", false, "")
	cmd.Flags().Bool("dependencies", false, "")
	cmd.Flags().Bool("compatibility", false, "")
	cmd.Flags().Bool("non-interactive", true, "")
	cmd.Flags().String("output-format", "json", "")

	// Set the flags to trigger machine-readable mode
	cmd.Flags().Set("non-interactive", "true")
	cmd.Flags().Set("output-format", "json")

	err := tc.ExecuteInfo(cmd, []string{"test-template"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteValidate_Success(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	validationResult := createTestValidationResult(true)

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ValidateTemplate", "./test-template").Return(validationResult, nil)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Success", mock.AnythingOfType("string")).Return("SUCCESS")

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("fix", false, "")
	cmd.Flags().String("output-format", "text", "")
	cmd.Flags().Bool("non-interactive", false, "")

	err := tc.ExecuteValidate(cmd, []string{"./test-template"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteValidate_Invalid(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	validationResult := createTestValidationResult(false)

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ValidateTemplate", "./invalid-template").Return(validationResult, nil)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("fix", false, "")
	cmd.Flags().String("output-format", "text", "")
	cmd.Flags().Bool("non-interactive", false, "")

	err := tc.ExecuteValidate(cmd, []string{"./invalid-template"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "🚫")
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteValidate_Error(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	expectedError := errors.New("validation service error")

	// Setup expectations
	mockCLI.On("IsNonInteractiveMode", mock.AnythingOfType("*cobra.Command")).Return(false)
	mockCLI.On("ValidateTemplate", "./error-template").Return((*interfaces.TemplateValidationResult)(nil), expectedError)
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("fix", false, "")
	cmd.Flags().String("output-format", "text", "")
	cmd.Flags().Bool("non-interactive", false, "")

	err := tc.ExecuteValidate(cmd, []string{"./error-template"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "🚫")
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_ExecuteValidate_MachineReadable(t *testing.T) {
	mockCLI := &MockTemplateCLI{}
	tc := NewTemplateCommands(mockCLI)

	validationResult := createTestValidationResult(true)

	// Setup expectations
	mockCLI.On("ValidateTemplate", "./test-template").Return(validationResult, nil)
	mockCLI.On("OutputMachineReadable", validationResult, "json").Return(nil)

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.Flags().Bool("detailed", false, "")
	cmd.Flags().Bool("fix", false, "")
	cmd.Flags().String("output-format", "json", "")
	cmd.Flags().Bool("non-interactive", true, "")

	// Set the flags to trigger machine-readable mode
	cmd.Flags().Set("non-interactive", "true")
	cmd.Flags().Set("output-format", "json")

	err := tc.ExecuteValidate(cmd, []string{"./test-template"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestTemplateCommands_SetupListFlags(t *testing.T) {
	tc := &TemplateCommands{}
	cmd := &cobra.Command{}

	tc.SetupListFlags(cmd)

	// Verify flags are set up correctly
	assert.True(t, cmd.Flags().HasAvailableFlags())

	categoryFlag := cmd.Flags().Lookup("category")
	assert.NotNil(t, categoryFlag)
	assert.Equal(t, "", categoryFlag.DefValue)

	technologyFlag := cmd.Flags().Lookup("technology")
	assert.NotNil(t, technologyFlag)
	assert.Equal(t, "", technologyFlag.DefValue)

	tagsFlag := cmd.Flags().Lookup("tags")
	assert.NotNil(t, tagsFlag)

	searchFlag := cmd.Flags().Lookup("search")
	assert.NotNil(t, searchFlag)
	assert.Equal(t, "", searchFlag.DefValue)

	detailedFlag := cmd.Flags().Lookup("detailed")
	assert.NotNil(t, detailedFlag)
	assert.Equal(t, "false", detailedFlag.DefValue)
}

func TestTemplateCommands_SetupInfoFlags(t *testing.T) {
	tc := &TemplateCommands{}
	cmd := &cobra.Command{}

	tc.SetupInfoFlags(cmd)

	// Verify flags are set up correctly
	assert.True(t, cmd.Flags().HasAvailableFlags())

	detailedFlag := cmd.Flags().Lookup("detailed")
	assert.NotNil(t, detailedFlag)
	assert.Equal(t, "false", detailedFlag.DefValue)

	variablesFlag := cmd.Flags().Lookup("variables")
	assert.NotNil(t, variablesFlag)
	assert.Equal(t, "false", variablesFlag.DefValue)

	dependenciesFlag := cmd.Flags().Lookup("dependencies")
	assert.NotNil(t, dependenciesFlag)
	assert.Equal(t, "false", dependenciesFlag.DefValue)

	compatibilityFlag := cmd.Flags().Lookup("compatibility")
	assert.NotNil(t, compatibilityFlag)
	assert.Equal(t, "false", compatibilityFlag.DefValue)
}

func TestTemplateCommands_SetupValidateFlags(t *testing.T) {
	tc := &TemplateCommands{}
	cmd := &cobra.Command{}

	tc.SetupValidateFlags(cmd)

	// Verify flags are set up correctly
	assert.True(t, cmd.Flags().HasAvailableFlags())

	detailedFlag := cmd.Flags().Lookup("detailed")
	assert.NotNil(t, detailedFlag)
	assert.Equal(t, "false", detailedFlag.DefValue)

	fixFlag := cmd.Flags().Lookup("fix")
	assert.NotNil(t, fixFlag)
	assert.Equal(t, "false", fixFlag.DefValue)

	outputFormatFlag := cmd.Flags().Lookup("output-format")
	assert.NotNil(t, outputFormatFlag)
	assert.Equal(t, "text", outputFormatFlag.DefValue)
}
