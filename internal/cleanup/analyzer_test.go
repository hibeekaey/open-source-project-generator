package cleanup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCodeAnalyzer_AnalyzeTODOComments(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "analyzer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with TODO comments
	testFile := filepath.Join(tempDir, "test.go")
	content := `package main

import "fmt"

// TODO: Implement proper error handling
func main() {
	// FIXME: This is a security vulnerability
	fmt.Println("Hello, World!")
	
	// HACK: Temporary workaround
	doSomething()
}

// XXX: This function needs refactoring
func doSomething() {
	// NOTE: This is just a note
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Analyze TODO comments
	analyzer := NewCodeAnalyzer()
	todos, err := analyzer.AnalyzeTODOComments(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze TODO comments: %v", err)
	}

	// Verify results
	expectedTodos := 5 // TODO, FIXME, HACK, XXX, NOTE
	if len(todos) != expectedTodos {
		t.Errorf("Expected %d TODOs, got %d", expectedTodos, len(todos))
	}

	// Check specific TODOs
	todoTypes := make(map[string]bool)
	for _, todo := range todos {
		todoTypes[todo.Type] = true

		if todo.File != testFile {
			t.Errorf("Expected file %s, got %s", testFile, todo.File)
		}

		if todo.Line <= 0 {
			t.Errorf("Invalid line number: %d", todo.Line)
		}
	}

	expectedTypes := []string{"TODO", "FIXME", "HACK", "XXX", "NOTE"}
	for _, expectedType := range expectedTypes {
		if !todoTypes[expectedType] {
			t.Errorf("Expected to find %s comment", expectedType)
		}
	}
}

func TestCodeAnalyzer_ValidateImportOrganization(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "import_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with poorly organized imports
	testFile := filepath.Join(tempDir, "test.go")
	content := `package main

import (
	"github.com/external/package"
	"fmt"
	"os"
	"github.com/another/package"
)

func main() {
	fmt.Println("Hello")
	os.Exit(0)
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Analyze import organization
	analyzer := NewCodeAnalyzer()
	issues, err := analyzer.ValidateImportOrganization(tempDir)
	if err != nil {
		t.Fatalf("Failed to validate import organization: %v", err)
	}

	// Should find import organization issues
	if len(issues) == 0 {
		t.Error("Expected to find import organization issues")
	}

	for _, issue := range issues {
		if issue.Type != "grouping" {
			t.Errorf("Expected grouping issue, got %s", issue.Type)
		}

		if issue.File != testFile {
			t.Errorf("Expected file %s, got %s", testFile, issue.File)
		}
	}
}

func TestCodeAnalyzer_IdentifyUnusedCode(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "unused_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with unused import
	testFile := filepath.Join(tempDir, "test.go")
	content := `package main

import (
	"fmt"
	"os"
	"unused/package"
)

func main() {
	fmt.Println("Hello")
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Analyze unused code
	analyzer := NewCodeAnalyzer()
	unused, err := analyzer.IdentifyUnusedCode(tempDir)
	if err != nil {
		t.Fatalf("Failed to identify unused code: %v", err)
	}

	// Should find unused imports
	foundUnusedImport := false
	for _, item := range unused {
		if item.Type == "import" && item.Name == "unused/package" {
			foundUnusedImport = true
			break
		}
	}

	if !foundUnusedImport {
		t.Error("Expected to find unused import")
	}
}

func TestCodeAnalyzer_DeterminePriority(t *testing.T) {
	analyzer := NewCodeAnalyzer()

	tests := []struct {
		todoType string
		message  string
		expected Priority
	}{
		{"FIXME", "security issue", PriorityHigh},
		{"BUG", "critical bug", PriorityHigh},
		{"TODO", "security vulnerability", PriorityHigh},
		{"HACK", "performance issue", PriorityMedium},
		{"TODO", "performance optimization", PriorityMedium},
		{"TODO", "add feature", PriorityLow},
		{"NOTE", "documentation", PriorityLow},
	}

	for _, test := range tests {
		priority := analyzer.determinePriority(test.todoType, test.message)
		if priority != test.expected {
			t.Errorf("For %s '%s', expected priority %d, got %d",
				test.todoType, test.message, test.expected, priority)
		}
	}
}

func TestCodeAnalyzer_DetermineCategory(t *testing.T) {
	analyzer := NewCodeAnalyzer()

	tests := []struct {
		message  string
		expected Category
	}{
		{"security vulnerability", CategorySecurity},
		{"performance optimization", CategoryPerformance},
		{"update documentation", CategoryDocumentation},
		{"refactor this code", CategoryRefactor},
		{"fix this bug", CategoryBug},
		{"add new feature", CategoryFeature},
	}

	for _, test := range tests {
		category := analyzer.determineCategory(test.message)
		if category != test.expected {
			t.Errorf("For message '%s', expected category %d, got %d",
				test.message, test.expected, category)
		}
	}
}
