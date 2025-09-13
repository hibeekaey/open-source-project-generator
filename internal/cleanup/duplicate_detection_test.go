package cleanup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicateCode(t *testing.T) {
	// Create a temporary directory with test files
	tempDir := t.TempDir()

	// Create test files with duplicate code
	file1Content := `package test

func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	if len(input) > 100 {
		return fmt.Errorf("input too long")
	}
	return nil
}

func HelperFunction() string {
	return "helper"
}
`

	file2Content := `package test2

func ValidateData(data string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("data cannot be empty")
	}
	if len(data) > 100 {
		return fmt.Errorf("data too long")
	}
	return nil
}

func HelperFunction() string {
	return "helper"
}
`

	// Write test files
	file1Path := filepath.Join(tempDir, "file1.go")
	file2Path := filepath.Join(tempDir, "file2.go")

	if err := os.WriteFile(file1Path, []byte(file1Content), 0644); err != nil {
		t.Fatalf("Failed to write test file1: %v", err)
	}

	if err := os.WriteFile(file2Path, []byte(file2Content), 0644); err != nil {
		t.Fatalf("Failed to write test file2: %v", err)
	}

	// Test duplicate detection
	analyzer := NewCodeAnalyzer()
	duplicates, err := analyzer.FindDuplicateCode(tempDir)
	if err != nil {
		t.Fatalf("FindDuplicateCode failed: %v", err)
	}

	// Verify results
	if len(duplicates) == 0 {
		t.Error("Expected to find duplicate code, but found none")
	}

	// Check that we found the duplicate helper functions
	foundHelperDuplicate := false
	for _, dup := range duplicates {
		if len(dup.Files) >= 2 {
			foundHelperDuplicate = true
			break
		}
	}

	if !foundHelperDuplicate {
		t.Error("Expected to find duplicate helper functions")
	}

	t.Logf("Found %d duplicate code blocks", len(duplicates))
	for i, dup := range duplicates {
		t.Logf("Duplicate %d: %s (similarity: %.2f)", i+1, dup.Content, dup.Similarity)
		t.Logf("  Files: %v", dup.Files)
		t.Logf("  Suggestion: %s", dup.Suggestion)
	}
}

func TestFindDuplicateValidationLogic(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with similar validation patterns
	validationFile1 := `package validation1

func Validate1(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("empty input")
	}
	return nil
}
`

	validationFile2 := `package validation2

func Validate2(data string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("empty data")
	}
	return nil
}
`

	validationFile3 := `package validation3

func Validate3(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("empty value")
	}
	return nil
}
`

	// Write test files
	files := map[string]string{
		"val1.go": validationFile1,
		"val2.go": validationFile2,
		"val3.go": validationFile3,
	}

	for filename, content := range files {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	analyzer := NewCodeAnalyzer()
	duplicates := analyzer.findDuplicateValidationLogic(tempDir)

	if len(duplicates) == 0 {
		t.Error("Expected to find duplicate validation patterns, but found none")
	}

	t.Logf("Found %d duplicate validation patterns", len(duplicates))
	for i, dup := range duplicates {
		t.Logf("Pattern %d: %s", i+1, dup.Content)
		t.Logf("  Files: %v", dup.Files)
	}
}
