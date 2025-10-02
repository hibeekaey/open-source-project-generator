package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuditCLI implements AuditCLI interface for testing
type MockAuditCLI struct {
	mock.Mock
}

func (m *MockAuditCLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.AuditResult), args.Error(1)
}

func (m *MockAuditCLI) QuietOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockAuditCLI) VerboseOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockAuditCLI) DebugOutput(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockAuditCLI) Error(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockAuditCLI) Warning(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockAuditCLI) Info(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockAuditCLI) IsQuietMode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockAuditCLI) CreateAuditError(message string, score float64) error {
	args := m.Called(message, score)
	return args.Error(0)
}

func (m *MockAuditCLI) OutputMachineReadable(data interface{}, format string) error {
	args := m.Called(data, format)
	return args.Error(0)
}

func TestNewAuditCommand(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	assert.NotNil(t, cmd)
	assert.Equal(t, mockCLI, cmd.cli)
}

func TestAuditCommand_Execute_Success(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	// Mock successful audit
	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Quality: &interfaces.QualityAuditResult{
			Score:      80.0,
			CodeSmells: []interfaces.CodeSmell{},
		},
		Licenses: &interfaces.LicenseAuditResult{
			Score:      95.0,
			Compatible: true,
			Conflicts:  []interfaces.LicenseInfo{},
		},
		Performance: &interfaces.PerformanceAuditResult{
			Score:      75.0,
			BundleSize: 1024000,
		},
		Recommendations: []string{"Consider updating dependencies"},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_WithPath(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	auditResult := &interfaces.AuditResult{
		ProjectPath:  "/test/path",
		AuditTime:    time.Now(),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Quality: &interfaces.QualityAuditResult{
			Score:      80.0,
			CodeSmells: []interfaces.CodeSmell{},
		},
		Licenses: &interfaces.LicenseAuditResult{
			Score:      95.0,
			Compatible: true,
			Conflicts:  []interfaces.LicenseInfo{},
		},
		Performance: &interfaces.PerformanceAuditResult{
			Score:      75.0,
			BundleSize: 1024000,
		},
		Recommendations: []string{"Consider updating dependencies"},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", "/test/path", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()

	err := cmd.Execute(cobraCmd, []string{"/test/path"})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_WithVulnerabilities(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	// Mock audit with vulnerabilities
	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 45.0,
		Security: &interfaces.SecurityAuditResult{
			Score: 30.0,
			Vulnerabilities: []interfaces.Vulnerability{
				{
					ID:          "CVE-2023-1234",
					Severity:    "critical",
					Title:       "Critical vulnerability in package",
					Description: "This is a critical security vulnerability",
					Package:     "vulnerable-package",
					Version:     "1.0.0",
					FixedIn:     "1.0.1",
				},
				{
					ID:          "CVE-2023-5678",
					Severity:    "high",
					Title:       "High severity vulnerability",
					Description: "This is a high severity vulnerability",
					Package:     "another-package",
					Version:     "2.0.0",
					FixedIn:     "2.1.0",
				},
			},
		},
		Quality: &interfaces.QualityAuditResult{
			Score: 60.0,
			CodeSmells: []interfaces.CodeSmell{
				{
					Type:        "complexity",
					Severity:    "medium",
					Description: "Function is too complex",
					File:        "main.go",
					Line:        42,
				},
			},
		},
		Licenses: &interfaces.LicenseAuditResult{
			Score:      50.0,
			Compatible: false,
			Conflicts: []interfaces.LicenseInfo{
				{
					Name:       "GPL-3.0",
					SPDXID:     "GPL-3.0",
					Package:    "gpl-package",
					Compatible: false,
				},
			},
		},
		Performance: &interfaces.PerformanceAuditResult{
			Score:      40.0,
			BundleSize: 5024000,
		},
		Recommendations: []string{
			"Update vulnerable-package to version 1.0.1",
			"Update another-package to version 2.1.0",
			"Refactor complex functions",
			"Review license compatibility",
		},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_FailOnHigh(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with fail-on-high flag
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", true, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	// Mock audit with low score
	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 65.0, // Below 70.0 threshold
		Security: &interfaces.SecurityAuditResult{
			Score:           65.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Recommendations: []string{},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("CreateAuditError", mock.AnythingOfType("string"), 65.0).Return(assert.AnError)

	err := cmd.Execute(cobraCmd, []string{})

	assert.Error(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_FailOnMedium(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with fail-on-medium flag
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", true, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	// Mock audit with low score
	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 45.0, // Below 50.0 threshold
		Security: &interfaces.SecurityAuditResult{
			Score:           45.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Recommendations: []string{},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("CreateAuditError", mock.AnythingOfType("string"), 45.0).Return(assert.AnError)

	err := cmd.Execute(cobraCmd, []string{})

	assert.Error(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_MinScore(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with min-score flag
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "text", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 80.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	// Mock audit with score below minimum
	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 75.0, // Below 80.0 minimum
		Security: &interfaces.SecurityAuditResult{
			Score:           75.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Recommendations: []string{},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "text",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("CreateAuditError", mock.AnythingOfType("string"), 75.0).Return(assert.AnError)

	err := cmd.Execute(cobraCmd, []string{})

	assert.Error(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_MachineReadableOutput(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with JSON output and non-interactive mode
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "json", "")
	cobraCmd.Flags().String("output-file", "", "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", true, "")

	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Recommendations: []string{},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "json",
		OutputFile:   "",
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("OutputMachineReadable", auditResult, "json").Return(nil)

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_Execute_WithReportGeneration(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create temporary file for report
	tmpDir := t.TempDir()
	reportFile := filepath.Join(tmpDir, "audit-report.json")

	// Create a test command with output file
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", true, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", true, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "json", "")
	cobraCmd.Flags().String("output-file", reportFile, "")
	cobraCmd.Flags().Bool("detailed", false, "")
	cobraCmd.Flags().Bool("fail-on-high", false, "")
	cobraCmd.Flags().Bool("fail-on-medium", false, "")
	cobraCmd.Flags().Float64("min-score", 0.0, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{}, "")
	cobraCmd.Flags().Bool("summary-only", false, "")
	cobraCmd.Flags().Bool("non-interactive", false, "")

	auditResult := &interfaces.AuditResult{
		ProjectPath:  ".",
		AuditTime:    time.Now(),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Recommendations: []string{},
	}

	expectedOptions := interfaces.AuditOptions{
		Security:     true,
		Quality:      true,
		Licenses:     true,
		Performance:  true,
		OutputFormat: "json",
		OutputFile:   reportFile,
		Detailed:     false,
	}

	mockCLI.On("AuditProject", ".", expectedOptions).Return(auditResult, nil)
	mockCLI.On("IsQuietMode").Return(false)
	mockCLI.On("QuietOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("VerboseOutput", mock.AnythingOfType("string"), mock.Anything).Return()
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")

	err := cmd.Execute(cobraCmd, []string{})

	assert.NoError(t, err)
	mockCLI.AssertExpectations(t)

	// Verify report file was created
	assert.FileExists(t, reportFile)

	// Verify report content
	content, err := os.ReadFile(reportFile)
	assert.NoError(t, err)

	var reportData interfaces.AuditResult
	err = json.Unmarshal(content, &reportData)
	assert.NoError(t, err)
	assert.Equal(t, auditResult.OverallScore, reportData.OverallScore)
}

func TestAuditCommand_parseFlags(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	// Create a test command with various flags
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().Bool("security", false, "")
	cobraCmd.Flags().Bool("quality", true, "")
	cobraCmd.Flags().Bool("licenses", false, "")
	cobraCmd.Flags().Bool("performance", true, "")
	cobraCmd.Flags().String("output-format", "json", "")
	cobraCmd.Flags().String("output-file", "report.json", "")
	cobraCmd.Flags().Bool("detailed", true, "")
	cobraCmd.Flags().StringSlice("exclude-categories", []string{"security", "licenses"}, "")
	cobraCmd.Flags().Bool("summary-only", true, "")

	mockCLI.On("DebugOutput", mock.AnythingOfType("string"), mock.Anything).Return()

	options, err := cmd.parseFlags(cobraCmd)

	assert.NoError(t, err)
	assert.NotNil(t, options)
	assert.False(t, options.Security)
	assert.True(t, options.Quality)
	assert.False(t, options.Licenses)
	assert.True(t, options.Performance)
	assert.Equal(t, "json", options.OutputFormat)
	assert.Equal(t, "report.json", options.OutputFile)
	assert.True(t, options.Detailed)

	mockCLI.AssertExpectations(t)
}

func TestAuditCommand_generateTextReport(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	auditResult := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Date(2023, 12, 1, 10, 30, 0, 0, time.UTC),
		OverallScore: 75.5,
		Security: &interfaces.SecurityAuditResult{
			Score: 80.0,
			Vulnerabilities: []interfaces.Vulnerability{
				{
					ID:          "CVE-2023-1234",
					Severity:    "high",
					Title:       "Test vulnerability",
					Description: "Test description",
					Package:     "test-package",
					Version:     "1.0.0",
					FixedIn:     "1.0.1",
				},
			},
		},
		Quality: &interfaces.QualityAuditResult{
			Score: 70.0,
			CodeSmells: []interfaces.CodeSmell{
				{
					Type:        "complexity",
					Severity:    "medium",
					Description: "Function too complex",
					File:        "main.go",
					Line:        42,
				},
			},
		},
		Recommendations: []string{
			"Update test-package to 1.0.1",
			"Refactor complex functions",
		},
	}

	report := cmd.generateTextReport(auditResult)

	reportStr := string(report)
	assert.Contains(t, reportStr, "🔍 Audit Report")
	assert.Contains(t, reportStr, "/test/project")
	assert.Contains(t, reportStr, "2023-12-01 10:30:00")
	assert.Contains(t, reportStr, "75.5/100")
	assert.Contains(t, reportStr, "🔒 Security Audit")
	assert.Contains(t, reportStr, "CVE-2023-1234")
	assert.Contains(t, reportStr, "✨ Quality Audit")
	assert.Contains(t, reportStr, "Function too complex")
	assert.Contains(t, reportStr, "💡 Recommendations:")
	assert.Contains(t, reportStr, "Update test-package to 1.0.1")
}

func TestAuditCommand_handleAuditError(t *testing.T) {
	mockCLI := &MockAuditCLI{}
	cmd := NewAuditCommand(mockCLI)

	originalErr := assert.AnError
	mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
	mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")

	err := cmd.handleAuditError(originalErr)

	assert.Error(t, err)
	// The error message will contain the mocked "ERROR" and "INFO" strings
	assert.Contains(t, err.Error(), "🚫")
	mockCLI.AssertExpectations(t)
}
