package cleanup

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// CodeRemover handles removal of unused code elements
type CodeRemover struct {
	fileSet *token.FileSet
	dryRun  bool
}

// RemovalPlan represents a plan for removing unused code
type RemovalPlan struct {
	UnusedItems []UnusedCodeItem
	FilesToFix  map[string][]UnusedCodeItem // file -> items to remove
	Summary     RemovalSummary
}

// RemovalSummary provides statistics about the removal plan
type RemovalSummary struct {
	TotalItems    int
	ImportItems   int
	FunctionItems int
	VariableItems int
	TypeItems     int
	ConstantItems int
	FilesAffected int
}

// NewCodeRemover creates a new code remover
func NewCodeRemover(dryRun bool) *CodeRemover {
	return &CodeRemover{
		fileSet: token.NewFileSet(),
		dryRun:  dryRun,
	}
}

// CreateRemovalPlan creates a plan for removing unused code
func (cr *CodeRemover) CreateRemovalPlan(unusedItems []UnusedCodeItem) *RemovalPlan {
	plan := &RemovalPlan{
		UnusedItems: unusedItems,
		FilesToFix:  make(map[string][]UnusedCodeItem),
	}

	// Group items by file
	for _, item := range unusedItems {
		plan.FilesToFix[item.File] = append(plan.FilesToFix[item.File], item)
	}

	// Calculate summary
	plan.Summary = RemovalSummary{
		TotalItems:    len(unusedItems),
		FilesAffected: len(plan.FilesToFix),
	}

	for _, item := range unusedItems {
		switch item.Type {
		case "import":
			plan.Summary.ImportItems++
		case "function":
			plan.Summary.FunctionItems++
		case "variable":
			plan.Summary.VariableItems++
		case "type":
			plan.Summary.TypeItems++
		case "constant":
			plan.Summary.ConstantItems++
		}
	}

	return plan
}

// ExecuteRemovalPlan executes the removal plan
func (cr *CodeRemover) ExecuteRemovalPlan(plan *RemovalPlan) error {
	if cr.dryRun {
		return cr.previewRemoval(plan)
	}

	for filePath, items := range plan.FilesToFix {
		if err := cr.removeUnusedFromFile(filePath, items); err != nil {
			return fmt.Errorf("failed to remove unused items from %s: %w", filePath, err)
		}
	}

	return nil
}

// previewRemoval shows what would be removed without making changes
func (cr *CodeRemover) previewRemoval(plan *RemovalPlan) error {
	fmt.Printf("=== Removal Plan Preview ===\n\n")
	fmt.Printf("Summary:\n")
	fmt.Printf("  Total items to remove: %d\n", plan.Summary.TotalItems)
	fmt.Printf("  Files affected: %d\n", plan.Summary.FilesAffected)
	fmt.Printf("  Imports: %d\n", plan.Summary.ImportItems)
	fmt.Printf("  Functions: %d\n", plan.Summary.FunctionItems)
	fmt.Printf("  Variables: %d\n", plan.Summary.VariableItems)
	fmt.Printf("  Types: %d\n", plan.Summary.TypeItems)
	fmt.Printf("  Constants: %d\n", plan.Summary.ConstantItems)
	fmt.Printf("\n")

	// Show details for a few files
	count := 0
	for filePath, items := range plan.FilesToFix {
		if count >= 5 { // Limit output
			fmt.Printf("... and %d more files\n", len(plan.FilesToFix)-5)
			break
		}

		fmt.Printf("File: %s (%d items)\n", filePath, len(items))
		for i, item := range items {
			if i >= 3 { // Limit items per file
				fmt.Printf("  ... and %d more items\n", len(items)-3)
				break
			}
			fmt.Printf("  - Line %d: %s %s (%s)\n", item.Line, item.Type, item.Name, item.Reason)
		}
		fmt.Printf("\n")
		count++
	}

	return nil
}

// removeUnusedFromFile removes unused items from a specific file
func (cr *CodeRemover) removeUnusedFromFile(filePath string, items []UnusedCodeItem) error {
	// Read and parse the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	file, err := parser.ParseFile(cr.fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// Remove unused imports first (they're easier to handle)
	modified := false
	for _, item := range items {
		if item.Type == "import" {
			if cr.removeUnusedImport(file, item.Name) {
				modified = true
			}
		}
	}

	// For other types of unused code, we'll be more conservative
	// and only remove simple cases to avoid breaking the code
	for _, item := range items {
		if item.Type != "import" {
			// Only remove if it's clearly safe
			if cr.isSafeToRemove(item) {
				if cr.removeUnusedDeclaration(file, item) {
					modified = true
				}
			}
		}
	}

	// Write the modified file back if changes were made
	if modified {
		var buf strings.Builder
		if err := format.Node(&buf, cr.fileSet, file); err != nil {
			return fmt.Errorf("failed to format modified file: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(buf.String()), 0644); err != nil {
			return fmt.Errorf("failed to write modified file: %w", err)
		}

		fmt.Printf("Cleaned up unused code in: %s\n", filePath)
	}

	return nil
}

// removeUnusedImport removes an unused import from the file
func (cr *CodeRemover) removeUnusedImport(file *ast.File, importPath string) bool {
	for i, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, `"`) == importPath {
			// Remove the import
			file.Imports = append(file.Imports[:i], file.Imports[i+1:]...)

			// Also remove from the GenDecl if it's in an import block
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
					for j, spec := range genDecl.Specs {
						if impSpec, ok := spec.(*ast.ImportSpec); ok {
							if strings.Trim(impSpec.Path.Value, `"`) == importPath {
								genDecl.Specs = append(genDecl.Specs[:j], genDecl.Specs[j+1:]...)
								return true
							}
						}
					}
				}
			}
			return true
		}
	}
	return false
}

// removeUnusedDeclaration removes an unused declaration from the file
func (cr *CodeRemover) removeUnusedDeclaration(file *ast.File, item UnusedCodeItem) bool {
	// This is a simplified implementation that only handles basic cases
	// In a production system, you'd want more sophisticated AST manipulation

	for i, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name != nil && d.Name.Name == item.Name {
				// Remove the function declaration
				file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
				return true
			}
		case *ast.GenDecl:
			// Handle variable, type, and constant declarations
			for j, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for k, name := range s.Names {
						if name.Name == item.Name {
							// Remove this name from the spec
							s.Names = append(s.Names[:k], s.Names[k+1:]...)
							if len(s.Names) == 0 {
								// Remove the entire spec if no names left
								d.Specs = append(d.Specs[:j], d.Specs[j+1:]...)
								if len(d.Specs) == 0 {
									// Remove the entire declaration if no specs left
									file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
								}
							}
							return true
						}
					}
				case *ast.TypeSpec:
					if s.Name.Name == item.Name {
						// Remove the type spec
						d.Specs = append(d.Specs[:j], d.Specs[j+1:]...)
						if len(d.Specs) == 0 {
							// Remove the entire declaration if no specs left
							file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
						}
						return true
					}
				}
			}
		}
	}
	return false
}

// isSafeToRemove determines if it's safe to remove a declaration
func (cr *CodeRemover) isSafeToRemove(item UnusedCodeItem) bool {
	// Be conservative - only remove items that are clearly safe

	// Don't remove exported items (they might be used by external packages)
	if len(item.Name) > 0 && item.Name[0] >= 'A' && item.Name[0] <= 'Z' {
		return false
	}

	// Don't remove main or init functions
	if item.Name == "main" || item.Name == "init" {
		return false
	}

	// Don't remove test functions
	if strings.HasPrefix(item.Name, "Test") || strings.HasPrefix(item.Name, "Benchmark") || strings.HasPrefix(item.Name, "Example") {
		return false
	}

	// Only remove if the reason indicates it's never used
	if strings.Contains(item.Reason, "never used") {
		return true
	}

	// For now, be conservative with other cases
	return false
}
