package cleanup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ValidationFramework ensures no functionality is lost during cleanup
type ValidationFramework struct {
	projectRoot string
	testTimeout time.Duration
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Success      bool
	TestsPassed  bool
	BuildSuccess bool
	Errors       []ValidationError
	Warnings     []string
	Duration     time.Duration
}

// ValidationError represents a validation error
type ValidationError struct {
	Type       string
	Message    string
	File       string
	Line       int
	Suggestion string
}

// TestResults represents test execution results
type TestResults struct {
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Coverage     float64
	Duration     time.Duration
	Failures     []TestFailure
}

// TestFailure represents a failed test
type TestFailure struct {
	TestName string
	Package  string
	Error    string
	Output   string
}

// NewValidationFramework creates a new validation framework
func NewValidationFramework(projectRoot string) *ValidationFramework {
	return &ValidationFramework{
		projectRoot: projectRoot,
		testTimeout: 10 * time.Minute,
	}
}

// ValidateProject performs comprehensive project validation
func (vf *ValidationFramework) ValidateProject() (*ValidationResult, error) {
	start := time.Now()

	result := &ValidationResult{
		Success:      true,
		TestsPassed:  true,
		BuildSuccess: true,
		Errors:       []ValidationError{},
		Warnings:     []string{},
	}

	// 1. Validate Go syntax and compilation
	if err := vf.validateBuild(result); err != nil {
		result.Success = false
		result.BuildSuccess = false
	}

	// 2. Run test suite
	if err := vf.runTests(result); err != nil {
		result.Success = false
		result.TestsPassed = false
	}

	// 3. Validate go.mod consistency
	if err := vf.validateGoMod(result); err != nil {
		result.Success = false
	}

	// 4. Check for basic code quality issues
	if err := vf.validateCodeQuality(result); err != nil {
		// Code quality issues are warnings, not failures
		result.Warnings = append(result.Warnings, err.Error())
	}

	result.Duration = time.Since(start)
	return result, nil
}

// validateBuild checks if the project builds successfully
func (vf *ValidationFramework) validateBuild(result *ValidationResult) error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = vf.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:       "build",
			Message:    "Build failed",
			Suggestion: "Fix compilation errors before proceeding",
		})

		// Parse build errors for more specific information
		vf.parseBuildErrors(string(output), result)
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}

// runTests executes the test suite and collects results
func (vf *ValidationFramework) runTests(result *ValidationResult) error {
	// Run tests with verbose output and coverage
	cmd := exec.Command("go", "test", "-v", "-cover", "./...")
	cmd.Dir = vf.projectRoot

	output, err := cmd.CombinedOutput()
	testOutput := string(output)

	// Parse test results
	testResults := vf.parseTestResults(testOutput)

	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:       "test",
			Message:    "Test suite failed",
			Suggestion: "Fix failing tests before proceeding with cleanup",
		})

		// Add specific test failures
		for _, failure := range testResults.Failures {
			result.Errors = append(result.Errors, ValidationError{
				Type:       "test_failure",
				Message:    fmt.Sprintf("Test %s failed: %s", failure.TestName, failure.Error),
				Suggestion: "Review and fix the failing test",
			})
		}

		return fmt.Errorf("tests failed: %w", err)
	}

	// Check coverage threshold
	if testResults.Coverage < 50.0 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Test coverage is low: %.1f%%", testResults.Coverage))
	}

	return nil
}

// validateGoMod checks go.mod consistency
func (vf *ValidationFramework) validateGoMod(result *ValidationResult) error {
	// Run go mod tidy to check for inconsistencies
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = vf.projectRoot

	if err := cmd.Run(); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:       "gomod",
			Message:    "go.mod has inconsistencies",
			Suggestion: "Run 'go mod tidy' to fix module dependencies",
		})
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Check if go.mod was modified
	cmd = exec.Command("git", "diff", "--name-only", "go.mod", "go.sum")
	cmd.Dir = vf.projectRoot

	output, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		result.Warnings = append(result.Warnings,
			"go.mod or go.sum was modified by 'go mod tidy'")
	}

	return nil
}

// validateCodeQuality performs basic code quality checks
func (vf *ValidationFramework) validateCodeQuality(result *ValidationResult) error {
	// Run go vet
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = vf.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		vetOutput := string(output)
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("go vet found issues: %s", vetOutput))
	}

	// Check for gofmt issues
	cmd = exec.Command("gofmt", "-l", ".")
	cmd.Dir = vf.projectRoot

	output, err = cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		result.Warnings = append(result.Warnings,
			"Some files are not properly formatted with gofmt")
	}

	return nil
}

// ValidateAfterChanges validates the project after making changes
func (vf *ValidationFramework) ValidateAfterChanges(changedFiles []string) (*ValidationResult, error) {
	// For now, run full validation
	// In a more sophisticated implementation, we could run targeted validation
	return vf.ValidateProject()
}

// CreateValidationCheckpoint creates a validation checkpoint
func (vf *ValidationFramework) CreateValidationCheckpoint() (*ValidationResult, error) {
	return vf.ValidateProject()
}

// CompareValidationResults compares two validation results
func (vf *ValidationFramework) CompareValidationResults(before, after *ValidationResult) []string {
	var differences []string

	if before.Success && !after.Success {
		differences = append(differences, "Validation status changed from success to failure")
	}

	if before.TestsPassed && !after.TestsPassed {
		differences = append(differences, "Tests were passing before but are now failing")
	}

	if before.BuildSuccess && !after.BuildSuccess {
		differences = append(differences, "Build was successful before but is now failing")
	}

	if len(after.Errors) > len(before.Errors) {
		differences = append(differences,
			fmt.Sprintf("Number of errors increased from %d to %d",
				len(before.Errors), len(after.Errors)))
	}

	return differences
}

// Helper methods

func (vf *ValidationFramework) parseBuildErrors(output string, result *ValidationResult) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") && (strings.Contains(line, "error") || strings.Contains(line, "undefined")) {
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 3 {
				result.Errors = append(result.Errors, ValidationError{
					Type:    "syntax",
					File:    parts[0],
					Message: strings.TrimSpace(parts[2]),
				})
			}
		}
	}
}

func (vf *ValidationFramework) parseTestResults(output string) *TestResults {
	results := &TestResults{
		Failures: []TestFailure{},
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Parse test results - simplified implementation
		if strings.Contains(line, "PASS") {
			results.PassedTests++
		} else if strings.Contains(line, "FAIL") {
			results.FailedTests++
			// Extract test failure information
			if strings.Contains(line, "---") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					results.Failures = append(results.Failures, TestFailure{
						TestName: parts[2],
						Error:    line,
					})
				}
			}
		} else if strings.Contains(line, "coverage:") {
			// Parse coverage percentage
			// This is a simplified parser
			if strings.Contains(line, "%") {
				// Extract coverage percentage
				results.Coverage = 0.0 // Simplified
			}
		}
	}

	results.TotalTests = results.PassedTests + results.FailedTests + results.SkippedTests
	return results
}

// EnsureProjectIntegrity performs integrity checks
func (vf *ValidationFramework) EnsureProjectIntegrity() error {
	// Check that essential files exist
	essentialFiles := []string{
		"go.mod",
		"README.md",
	}

	for _, file := range essentialFiles {
		path := filepath.Join(vf.projectRoot, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("essential file missing: %s", file)
		}
	}

	// Check that essential directories exist
	essentialDirs := []string{
		"internal",
		"pkg",
		"cmd",
	}

	for _, dir := range essentialDirs {
		path := filepath.Join(vf.projectRoot, dir)
		if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
			return fmt.Errorf("essential directory missing or not a directory: %s", dir)
		}
	}

	// Check that at least one main.go exists in cmd subdirectories
	cmdDir := filepath.Join(vf.projectRoot, "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return fmt.Errorf("failed to read cmd directory: %w", err)
	}

	hasMainFile := false
	for _, entry := range entries {
		if entry.IsDir() {
			mainPath := filepath.Join(cmdDir, entry.Name(), "main.go")
			if _, err := os.Stat(mainPath); err == nil {
				hasMainFile = true
				break
			}
		}
	}

	if !hasMainFile {
		return fmt.Errorf("no main.go file found in cmd subdirectories")
	}

	return nil
}
