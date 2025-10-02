package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockValidateCLI implements ValidateCLI interface for testing
type MockValidateCLI struct {
	mock.Mock
}

func (m *MockValidateCLI) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.ValidationResult), args.Error(1)
}

func (m *MockValidateCLI) QuietOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockValidateCLI) VerboseOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockValidateCLI) DebugOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockValidateCLI) Error(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockValidateCLI) Warning(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockValidateCLI) Info(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockValidateCLI) IsQuietMode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockValidateCLI) CreateValidationError(message string, details map[string]interface{}) error {
	args := m.Called(message, details)
	return args.Error(0)
}

func TestNewValidateCommand(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	assert.NotNil(t, cmd)
	assert.Equal(t, mockCLI, cmd.cli)
}

func TestValidateCommand_Execute_Success(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", false, "")
	cobraCmd.Flags().Bool("report", false, "")
	cobraCmd.Flags().String("report-format", "text", "")
	cobraCmd.Flags().StringSlice("rules", []string{}, "")
	cobraCmd.Flags().Bool("ignore-warnings", false, "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().String("output", "", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")
	cobraCmd.Flags().String("output-format", "text", "")

	// Mock successful validation
	validResult := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	expectedOptions := interfaces.ValidationOptions{
		Verbose:        false,
		Fix:            false,
		Report:         false,
		ReportFormat:   "text",
		Rules:          []string{},
		IgnoreWarnings: false,
		OutputFile:     "",
	}

	mockCLI.On("ValidateProject", ".", expectedOptions).Return(validResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Warning", mock.AnythingOfType("string")).Return("WARNING")

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestValidateCommand_Execute_WithPath(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", false, "")
	cobraCmd.Flags().Bool("report", false, "")
	cobraCmd.Flags().String("report-format", "text", "")
	cobraCmd.Flags().StringSlice("rules", []string{}, "")
	cobraCmd.Flags().Bool("ignore-warnings", false, "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().String("output", "", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")
	cobraCmd.Flags().String("output-format", "text", "")

	validResult := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	expectedOptions := interfaces.ValidationOptions{
		Verbose:        false,
		Fix:            false,
		Report:         false,
		ReportFormat:   "text",
		Rules:          []string{},
		IgnoreWarnings: false,
		OutputFile:     "",
	}

	mockCLI.On("ValidateProject", "/test/path", expectedOptions).Return(validResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Warning", mock.AnythingOfType("string")).Return("WARNING")

	err := cmd.Execute(cobraCmd, []string{"/test/path"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestValidateCommand_Execute_WithIssues(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", false, "")
	cobraCmd.Flags().Bool("report", false, "")
	cobraCmd.Flags().String("report-format", "text", "")
	cobraCmd.Flags().StringSlice("rules", []string{}, "")
	cobraCmd.Flags().Bool("ignore-warnings", false, "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().String("output", "", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")
	cobraCmd.Flags().String("output-format", "text", "")

	// Mock validation with issues
	invalidResult := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: "high",
				Message:  "Missing required file",
				File:     "package.json",
				Line:     0,
				Column:   0,
				Rule:     "required-files",
				Fixable:  true,
			},
		},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	expectedOptions := interfaces.ValidationOptions{
		Verbose:        false,
		Fix:            false,
		Report:         false,
		ReportFormat:   "text",
		Rules:          []string{},
		IgnoreWarnings: false,
		OutputFile:     "",
	}

	mockCLI.On("ValidateProject", ".", expectedOptions).Return(invalidResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")
	mockCLI.On("Warning", mock.AnythingOfType("string")).Return("WARNING")
	mockCLI.On("CreateValidationError", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(assert.AnError)

	err := cmd.Execute(cobraCmd, []string{})

	assert.Error(t, err)
	mockCLI.AssertExpectations(t)
}

func TestValidateCommand_Execute_WithReport(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create temporary file for report
	tempDir := t.TempDir()
	reportFile := filepath.Join(tempDir, "validation-report.json")

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", false, "")
	cobraCmd.Flags().Bool("report", true, "")
	cobraCmd.Flags().String("report-format", "json", "")
	cobraCmd.Flags().StringSlice("rules", []string{}, "")
	cobraCmd.Flags().Bool("ignore-warnings", false, "")
	cobraCmd.Flags().String("output-file", reportFile, "")
	cobraCmd.Flags().String("output", "", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")
	cobraCmd.Flags().String("output-format", "text", "")

	validResult := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	expectedOptions := interfaces.ValidationOptions{
		Verbose:        false,
		Fix:            false,
		Report:         true,
		ReportFormat:   "json",
		Rules:          []string{},
		IgnoreWarnings: false,
		OutputFile:     reportFile,
	}

	mockCLI.On("ValidateProject", ".", expectedOptions).Return(validResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Warning", mock.AnythingOfType("string")).Return("WARNING")

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)

	// Verify report file was created
	assert.FileExists(t, reportFile)

	// Verify report content
	content, err := os.ReadFile(reportFile)
	assert.NoError(t, err)

	var reportData interfaces.ValidationResult
	err = json.Unmarshal(content, &reportData)
	assert.NoError(t, err)
	assert.True(t, reportData.Valid)

	mockCLI.AssertExpectations(t)
}

func TestValidateCommand_parseFlags(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", true, "")
	cobraCmd.Flags().Bool("report", true, "")
	cobraCmd.Flags().String("report-format", "json", "")
	cobraCmd.Flags().StringSlice("rules", []string{"rule1", "rule2"}, "")
	cobraCmd.Flags().Bool("ignore-warnings", true, "")
	cobraCmd.Flags().String("output-file", "test.json", "")
	cobraCmd.Flags().String("output", "", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", true, "")

	// Set flag values
	cobraCmd.Flags().Set("fix", "true")
	cobraCmd.Flags().Set("report", "true")
	cobraCmd.Flags().Set("report-format", "json")
	cobraCmd.Flags().Set("rules", "rule1,rule2")
	cobraCmd.Flags().Set("ignore-warnings", "true")
	cobraCmd.Flags().Set("output-file", "test.json")
	cobraCmd.Flags().Set("verbose", "true")

	options, err := cmd.parseFlags(cobraCmd)

	assert.NoError(t, err)
	assert.NotNil(t, options)
	assert.True(t, options.Fix)
	assert.True(t, options.Report)
	assert.Equal(t, "json", options.ReportFormat)
	assert.Equal(t, []string{"rule1", "rule2"}, options.Rules)
	assert.True(t, options.IgnoreWarnings)
	assert.Equal(t, "test.json", options.OutputFile)
	assert.True(t, options.Verbose)
}

func TestValidateCommand_parseFlags_OutputAlias(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("fix", false, "")
	cobraCmd.Flags().Bool("report", false, "")
	cobraCmd.Flags().String("report-format", "text", "")
	cobraCmd.Flags().StringSlice("rules", []string{}, "")
	cobraCmd.Flags().Bool("ignore-warnings", false, "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().String("output", "test-output.json", "")
	cobraCmd.Flags().Bool("strict", false, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().StringSlice("exclude-rules", []string{}, "")
	cobraCmd.Flags().Bool("show-fixes", false, "")
	cobraCmd.Flags().Bool("verbose", false, "")

	// Set flag values
	cobraCmd.Flags().Set("output", "test-output.json")

	options, err := cmd.parseFlags(cobraCmd)

	assert.NoError(t, err)
	assert.NotNil(t, options)
	assert.Equal(t, "test-output.json", options.OutputFile)
}

func TestValidateCommand_outputMachineReadable_JSON(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	result := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.outputMachineReadable(result, "json")

	w.Close()
	os.Stdout = oldStdout

	assert.NoError(t, err)

	// Read captured output
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	output := buf.String()

	// Verify JSON output
	var outputResult interfaces.ValidationResult
	err = json.Unmarshal([]byte(output), &outputResult)
	assert.NoError(t, err)
	assert.True(t, outputResult.Valid)
}

func TestValidateCommand_outputMachineReadable_UnsupportedFormat(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	result := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	err := cmd.outputMachineReadable(result, "xml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}

func TestValidateCommand_generateTextReport(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: "high",
				Message:  "Missing required file",
				File:     "package.json",
				Line:     1,
				Column:   1,
				Rule:     "required-files",
				Fixable:  true,
			},
		},
		Warnings: []interfaces.ValidationIssue{
			{
				Type:     "warning",
				Severity: "medium",
				Message:  "Deprecated dependency",
				File:     "package.json",
				Line:     10,
				Column:   5,
				Rule:     "deprecated-deps",
				Fixable:  false,
			},
		},
		FixSuggestions: []interfaces.FixSuggestion{
			{
				Description: "Create missing package.json file",
				AutoFixable: true,
			},
		},
	}

	content := cmd.generateTextReport(result)

	report := string(content)
	assert.Contains(t, report, "🔍 Validation Report")
	assert.Contains(t, report, "❌ Needs attention")
	assert.Contains(t, report, "📊 Issues: 1")
	assert.Contains(t, report, "⚠️  Warnings: 1")
	assert.Contains(t, report, "Missing required file")
	assert.Contains(t, report, "Deprecated dependency")
	assert.Contains(t, report, "Create missing package.json file")
	assert.Contains(t, report, "Auto-fixable with --fix flag")
}

func TestValidateCommand_generateReport_JSON(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	tempDir := t.TempDir()
	reportFile := filepath.Join(tempDir, "test-report.json")

	result := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	err := cmd.generateReport(result, "json", reportFile)

	assert.NoError(t, err)
	assert.FileExists(t, reportFile)

	// Verify content
	content, err := os.ReadFile(reportFile)
	assert.NoError(t, err)

	var reportData interfaces.ValidationResult
	err = json.Unmarshal(content, &reportData)
	assert.NoError(t, err)
	assert.True(t, reportData.Valid)
}

func TestValidateCommand_generateReport_Text(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	tempDir := t.TempDir()
	reportFile := filepath.Join(tempDir, "test-report.txt")

	result := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	err := cmd.generateReport(result, "text", reportFile)

	assert.NoError(t, err)
	assert.FileExists(t, reportFile)

	// Verify content
	content, err := os.ReadFile(reportFile)
	assert.NoError(t, err)

	report := string(content)
	assert.Contains(t, report, "🔍 Validation Report")
	assert.Contains(t, report, "✅ Looks good!")
}

func TestValidateCommand_handleValidationResult_Success(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	result := &interfaces.ValidationResult{
		Valid:          true,
		Issues:         []interfaces.ValidationIssue{},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	err := cmd.handleValidationResult(result, "/test/path")

	assert.NoError(t, err)
}

func TestValidateCommand_handleValidationResult_WithIssues(t *testing.T) {
	mockCLI := &MockValidateCLI{}
	cmd := NewValidateCommand(mockCLI)

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: "high",
				Message:  "Test issue",
			},
		},
		Warnings:       []interfaces.ValidationIssue{},
		FixSuggestions: []interfaces.FixSuggestion{},
	}

	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("CreateValidationError", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(assert.AnError)

	err := cmd.handleValidationResult(result, "/test/path")

	assert.Error(t, err)
	mockCLI.AssertExpectations(t)
}
