//go:build !ci

package cleanup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIdentifyUnusedCode(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file with unused code
	testContent := `package test

import (
	"fmt"
	"strings" // unused import
	"os"      // used import
)

// UnusedFunction is never called
func UnusedFunction() {
	fmt.Println("unused")
}

// UsedFunction is called from main
func UsedFunction() {
	fmt.Println("used")
}

var unusedVar = "never used"
var usedVar = "used in main"

type UnusedType struct {
	Field string
}

const unusedConst = "never used"

func main() {
	UsedFunction()
	fmt.Println(usedVar)
	os.Exit(0)
}
`

	testFile := filepath.Join(tempDir, "test.go")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	unused, err := analyzer.IdentifyUnusedCode(tempDir)
	if err != nil {
		t.Fatalf("IdentifyUnusedCode failed: %v", err)
	}

	// Verify we found some unused items
	if len(unused) == 0 {
		t.Error("Expected to find unused code items")
	}

	// Check for specific unused items
	foundUnusedImport := false
	foundUnusedFunction := false

	for _, item := range unused {
		if item.Type == "import" && strings.Contains(item.Name, "strings") {
			foundUnusedImport = true
		}
		if item.Type == "function" && item.Name == "UnusedFunction" {
			foundUnusedFunction = true
		}
	}

	if !foundUnusedImport {
		t.Error("Expected to find unused import 'strings'")
	}
	if !foundUnusedFunction {
		t.Error("Expected to find unused function 'UnusedFunction'")
	}

	t.Logf("Found %d unused items", len(unused))
	for _, item := range unused {
		t.Logf("  %s:%d - %s %s (%s)", item.File, item.Line, item.Type, item.Name, item.Reason)
	}
}

func TestCreateRemovalPlan(t *testing.T) {
	unusedItems := []UnusedCodeItem{
		{File: "test1.go", Line: 10, Type: "import", Name: "unused/package", Reason: "Import not used"},
		{File: "test1.go", Line: 20, Type: "function", Name: "unusedFunc", Reason: "Function never used"},
		{File: "test2.go", Line: 15, Type: "variable", Name: "unusedVar", Reason: "Variable never used"},
	}

	remover := NewCodeRemover(true)
	plan := remover.CreateRemovalPlan(unusedItems)

	if plan.Summary.TotalItems != 3 {
		t.Errorf("Expected 3 total items, got %d", plan.Summary.TotalItems)
	}

	if plan.Summary.ImportItems != 1 {
		t.Errorf("Expected 1 import item, got %d", plan.Summary.ImportItems)
	}

	if plan.Summary.FunctionItems != 1 {
		t.Errorf("Expected 1 function item, got %d", plan.Summary.FunctionItems)
	}

	if plan.Summary.VariableItems != 1 {
		t.Errorf("Expected 1 variable item, got %d", plan.Summary.VariableItems)
	}

	if plan.Summary.FilesAffected != 2 {
		t.Errorf("Expected 2 files affected, got %d", plan.Summary.FilesAffected)
	}

	// Test preview (dry run)
	if err := remover.ExecuteRemovalPlan(plan); err != nil {
		t.Errorf("ExecuteRemovalPlan (dry run) failed: %v", err)
	}
}

func TestIsSafeToRemove(t *testing.T) {
	remover := NewCodeRemover(true)

	tests := []struct {
		name string
		item UnusedCodeItem
		want bool
	}{
		{
			name: "unexported function never used",
			item: UnusedCodeItem{Name: "unusedFunc", Type: "function", Reason: "function 'unusedFunc' is never used"},
			want: true,
		},
		{
			name: "exported function",
			item: UnusedCodeItem{Name: "ExportedFunc", Type: "function", Reason: "function 'ExportedFunc' is never used"},
			want: false,
		},
		{
			name: "main function",
			item: UnusedCodeItem{Name: "main", Type: "function", Reason: "function 'main' is never used"},
			want: false,
		},
		{
			name: "test function",
			item: UnusedCodeItem{Name: "TestSomething", Type: "function", Reason: "function 'TestSomething' is never used"},
			want: false,
		},
		{
			name: "function only used locally",
			item: UnusedCodeItem{Name: "localFunc", Type: "function", Reason: "function 'localFunc' is only used in the same file"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := remover.isSafeToRemove(tt.item)
			if got != tt.want {
				t.Errorf("isSafeToRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}
