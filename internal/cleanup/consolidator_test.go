//go:build !ci

package cleanup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateConsolidationPlan(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files with duplicate validation patterns
	testFiles := map[string]string{
		"file1.go": `package test1

func Validate1(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}`,
		"file2.go": `package test2

func Validate2(data string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("data cannot be empty")
	}
	return nil
}`,
		"file3.go": `package test3

func Validate3(items []string) error {
	if len(items) == 0 {
		return fmt.Errorf("items cannot be empty")
	}
	return nil
}`,
	}

	// Write test files
	for filename, content := range testFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	// Create consolidator and plan
	consolidator := NewCodeConsolidator(tempDir)
	plan, err := consolidator.CreateConsolidationPlan()
	if err != nil {
		t.Fatalf("Failed to create consolidation plan: %v", err)
	}

	// Verify plan contains expected utilities
	if len(plan.SharedUtilities) == 0 {
		t.Error("Expected shared utilities to be created")
	}

	// Preview the plan
	if err := consolidator.ExecuteConsolidationPlan(plan, true); err != nil {
		t.Fatalf("Failed to preview consolidation plan: %v", err)
	}

	t.Logf("Created consolidation plan with %d utilities, %d refactorings, %d validation fixes",
		len(plan.SharedUtilities), len(plan.Refactorings), len(plan.ValidationFixes))
}

func TestCreateSharedUtility(t *testing.T) {
	tempDir := t.TempDir()

	consolidator := NewCodeConsolidator(tempDir)

	utility := SharedUtility{
		Name:        "ValidateNonEmptyString",
		Package:     "utils",
		FilePath:    filepath.Join(tempDir, "utils", "validation.go"),
		Function:    consolidator.generateValidationFunction("ValidateNonEmptyString"),
		Description: "Test validation utility",
	}

	// Create the utility
	if err := consolidator.createSharedUtility(utility); err != nil {
		t.Fatalf("Failed to create shared utility: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(utility.FilePath); os.IsNotExist(err) {
		t.Error("Utility file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(utility.FilePath)
	if err != nil {
		t.Fatalf("Failed to read utility file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "ValidateNonEmptyString") {
		t.Error("Utility function not found in file")
	}

	t.Logf("Created utility file with content:\n%s", contentStr)
}
