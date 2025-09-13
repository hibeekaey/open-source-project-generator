package cleanup

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanupInfrastructureIntegration(t *testing.T) {
	// Create temporary project directory
	tempDir, err := os.MkdirTemp("", "cleanup_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up minimal project structure
	setupTestProject(t, tempDir)

	// Test cleanup manager initialization
	config := &Config{
		DryRun:          true,
		Verbose:         false,
		ValidationLevel: ValidationBasic,
	}

	manager, err := NewManager(tempDir, config)
	if err != nil {
		t.Fatalf("Failed to create cleanup manager: %v", err)
	}
	defer manager.Shutdown()

	// Test initialization
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize cleanup manager: %v", err)
	}

	// Test project analysis
	analysis, err := manager.AnalyzeProject()
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Verify analysis results
	if analysis == nil {
		t.Fatal("Analysis result is nil")
	}

	if len(analysis.TODOs) == 0 {
		t.Error("Expected to find TODO comments in test project")
	}

	// Test backup creation
	testFiles := []string{
		filepath.Join(tempDir, "main.go"),
		filepath.Join(tempDir, "test.go"),
	}

	backup, err := manager.CreateBackup(testFiles)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	if backup == nil {
		t.Fatal("Backup is nil")
	}

	// Test validation
	result, err := manager.ValidateProject()
	if err != nil {
		t.Fatalf("Failed to validate project: %v", err)
	}

	if result == nil {
		t.Fatal("Validation result is nil")
	}

	t.Logf("Integration test completed successfully")
	t.Logf("Analysis summary: %s", analysis.GetSummary())
	t.Logf("Validation success: %t", result.Success)
}

func setupTestProject(t *testing.T, projectDir string) {
	// Create essential directories
	dirs := []string{"internal", "pkg", "cmd/test"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create essential files
	files := map[string]string{
		"go.mod": `module test-project

go 1.23.0
`,
		"README.md": `# Test Project

This is a test project for cleanup infrastructure.
`,
		"cmd/test/main.go": `package main

import "fmt"

// TODO: Add proper error handling
func main() {
	// FIXME: This is a security issue
	fmt.Println("Hello, World!")
	
	// HACK: Temporary workaround
	doSomething()
}

// XXX: This function needs refactoring
func doSomething() {
	// NOTE: This is just a note
}
`,
		"test.go": `package main

import (
	"fmt"
	"os"
	"unused/package"
)

// TODO: Implement proper testing
func TestFunction() {
	fmt.Println("Test")
}

// UnusedFunction is not used anywhere
func UnusedFunction() {
	// This function is never called
}
`,
		"internal/app.go": `package internal

// TODO: Add application logic
type App struct {
	Name string
}

func NewApp() *App {
	return &App{Name: "test"}
}
`,
		"pkg/utils.go": `package pkg

import "fmt"

// TODO: Add utility functions
func PrintMessage(msg string) {
	fmt.Println(msg)
}
`,
	}

	for filename, content := range files {
		path := filepath.Join(projectDir, filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", filename, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}
}

func TestCleanupManagerLifecycle(t *testing.T) {
	// Test manager creation and shutdown
	tempDir, err := os.MkdirTemp("", "cleanup_lifecycle_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	setupTestProject(t, tempDir)

	config := DefaultConfig()
	config.DryRun = true

	manager, err := NewManager(tempDir, config)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test initialization
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// Test shutdown
	if err := manager.Shutdown(); err != nil {
		t.Fatalf("Failed to shutdown: %v", err)
	}
}

func TestBackupManagerIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "backup_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	backupDir := filepath.Join(tempDir, "backups")
	bm := NewBackupManager(backupDir)

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.go")
	testFile2 := filepath.Join(tempDir, "test2.go")

	content1 := "package main\n\nfunc main() {}\n"
	content2 := "package test\n\nfunc Test() {}\n"

	if err := os.WriteFile(testFile1, []byte(content1), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte(content2), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test backup creation
	files := []string{testFile1, testFile2}
	backup, err := bm.CreateBackup(files)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Verify backup
	if len(backup.Files) != 2 {
		t.Errorf("Expected 2 files in backup, got %d", len(backup.Files))
	}

	// Test backup listing
	backups, err := bm.ListBackups()
	if err != nil {
		t.Fatalf("Failed to list backups: %v", err)
	}

	if len(backups) == 0 {
		t.Error("Expected to find at least one backup")
	}

	// Test cleanup (with very short age to clean up our test backup)
	if err := bm.CleanupOldBackups(1 * time.Nanosecond); err != nil {
		t.Fatalf("Failed to cleanup old backups: %v", err)
	}
}

func TestValidationFrameworkIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "validation_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	setupTestProject(t, tempDir)

	vf := NewValidationFramework(tempDir)

	// Test project integrity check
	if err := vf.EnsureProjectIntegrity(); err != nil {
		t.Fatalf("Project integrity check failed: %v", err)
	}

	// Test validation checkpoint creation
	checkpoint, err := vf.CreateValidationCheckpoint()
	if err != nil {
		t.Fatalf("Failed to create validation checkpoint: %v", err)
	}

	if checkpoint == nil {
		t.Fatal("Validation checkpoint is nil")
	}

	// Test validation after changes (simulate)
	result, err := vf.ValidateAfterChanges([]string{"test.go"})
	if err != nil {
		t.Fatalf("Failed to validate after changes: %v", err)
	}

	if result == nil {
		t.Fatal("Validation result is nil")
	}

	// Test comparison
	differences := vf.CompareValidationResults(checkpoint, result)
	t.Logf("Found %d differences between validations", len(differences))
}

func TestCodeAnalyzerIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "analyzer_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	setupTestProject(t, tempDir)

	analyzer := NewCodeAnalyzer()

	// Test TODO analysis
	todos, err := analyzer.AnalyzeTODOComments(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze TODOs: %v", err)
	}

	if len(todos) == 0 {
		t.Error("Expected to find TODO comments")
	}

	// Verify TODO categories and priorities
	foundSecurity := false
	for _, todo := range todos {
		if todo.Category == CategorySecurity {
			foundSecurity = true
		}
		if todo.Priority == PriorityHigh && todo.Type == "FIXME" {
			t.Logf("Found high priority FIXME: %s", todo.Message)
		}
	}

	if !foundSecurity {
		t.Error("Expected to find security-related TODO")
	}

	// Test duplicate code detection
	duplicates, err := analyzer.FindDuplicateCode(tempDir)
	if err != nil {
		t.Fatalf("Failed to find duplicates: %v", err)
	}

	t.Logf("Found %d potential duplicate code blocks", len(duplicates))

	// Test unused code detection
	unused, err := analyzer.IdentifyUnusedCode(tempDir)
	if err != nil {
		t.Fatalf("Failed to identify unused code: %v", err)
	}

	if len(unused) == 0 {
		t.Error("Expected to find unused code items")
	}

	// Test import validation
	importIssues, err := analyzer.ValidateImportOrganization(tempDir)
	if err != nil {
		t.Fatalf("Failed to validate imports: %v", err)
	}

	t.Logf("Found %d import organization issues", len(importIssues))
}
