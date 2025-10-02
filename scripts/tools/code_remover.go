package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CodeRemover handles safe removal of unused code
type CodeRemover struct {
	fileSet   *token.FileSet
	backupDir string
	dryRun    bool
	verbose   bool
}

// NewCodeRemover creates a new code remover
func NewCodeRemover(backupDir string, dryRun, verbose bool) *CodeRemover {
	return &CodeRemover{
		fileSet:   token.NewFileSet(),
		backupDir: backupDir,
		dryRun:    dryRun,
		verbose:   verbose,
	}
}

// RemovalResult represents the result of a removal operation
type RemovalResult struct {
	File       string
	ItemType   string
	ItemName   string
	Success    bool
	Error      error
	BackupFile string
}

// RemoveUnusedImports removes unused imports from a Go file
func (cr *CodeRemover) RemoveUnusedImports(filename string) (*RemovalResult, error) {
	result := &RemovalResult{
		File:     filename,
		ItemType: "imports",
		ItemName: "unused imports",
	}

	if cr.verbose {
		fmt.Printf("Processing unused imports in: %s\n", filename)
	}

	// Create backup
	backupFile, err := cr.createBackup(filename)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.BackupFile = backupFile

	// Read and parse the file
	src, err := os.ReadFile(filename)
	if err != nil {
		result.Error = err
		return result, err
	}

	file, err := parser.ParseFile(cr.fileSet, filename, src, parser.ParseComments)
	if err != nil {
		result.Error = err
		return result, err
	}

	// Remove unused imports
	modified := cr.removeUnusedImportsFromAST(file)

	if !modified {
		if cr.verbose {
			fmt.Printf("No unused imports found in: %s\n", filename)
		}
		result.Success = true
		return result, nil
	}

	if cr.dryRun {
		fmt.Printf("[DRY RUN] Would remove unused imports from: %s\n", filename)
		result.Success = true
		return result, nil
	}

	// Write the modified file
	err = cr.writeASTToFile(file, filename)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Success = true
	if cr.verbose {
		fmt.Printf("Removed unused imports from: %s\n", filename)
	}

	return result, nil
}

// RemoveCommentedCodeBlock removes a specific commented code block
func (cr *CodeRemover) RemoveCommentedCodeBlock(filename string, startLine, endLine int) (*RemovalResult, error) {
	result := &RemovalResult{
		File:     filename,
		ItemType: "commented_code",
		ItemName: fmt.Sprintf("lines %d-%d", startLine, endLine),
	}

	if cr.verbose {
		fmt.Printf("Removing commented code block from %s (lines %d-%d)\n", filename, startLine, endLine)
	}

	// Create backup
	backupFile, err := cr.createBackup(filename)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.BackupFile = backupFile

	if cr.dryRun {
		fmt.Printf("[DRY RUN] Would remove commented code block from %s (lines %d-%d)\n", filename, startLine, endLine)
		result.Success = true
		return result, nil
	}

	// Read file lines
	lines, err := cr.readFileLines(filename)
	if err != nil {
		result.Error = err
		return result, err
	}

	// Validate line numbers
	if startLine < 1 || endLine > len(lines) || startLine > endLine {
		result.Error = fmt.Errorf("invalid line range: %d-%d (file has %d lines)", startLine, endLine, len(lines))
		return result, result.Error
	}

	// Remove the lines (convert to 0-based indexing)
	newLines := append(lines[:startLine-1], lines[endLine:]...)

	// Write back to file
	err = cr.writeFileLines(filename, newLines)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Success = true
	if cr.verbose {
		fmt.Printf("Removed commented code block from %s (lines %d-%d)\n", filename, startLine, endLine)
	}

	return result, nil
}

// RemoveUnusedFunction removes an unused function from a Go file
func (cr *CodeRemover) RemoveUnusedFunction(filename, functionName string, lineNumber int) (*RemovalResult, error) {
	result := &RemovalResult{
		File:     filename,
		ItemType: "function",
		ItemName: functionName,
	}

	if cr.verbose {
		fmt.Printf("Removing unused function %s from %s (around line %d)\n", functionName, filename, lineNumber)
	}

	// Create backup
	backupFile, err := cr.createBackup(filename)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.BackupFile = backupFile

	// Read and parse the file
	src, err := os.ReadFile(filename)
	if err != nil {
		result.Error = err
		return result, err
	}

	file, err := parser.ParseFile(cr.fileSet, filename, src, parser.ParseComments)
	if err != nil {
		result.Error = err
		return result, err
	}

	if cr.dryRun {
		fmt.Printf("[DRY RUN] Would remove function %s from %s\n", functionName, filename)
		result.Success = true
		return result, nil
	}

	// Remove the function from AST
	modified := cr.removeFunctionFromAST(file, functionName)
	if !modified {
		result.Error = fmt.Errorf("function %s not found in %s", functionName, filename)
		return result, result.Error
	}

	// Write the modified file
	err = cr.writeASTToFile(file, filename)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Success = true
	if cr.verbose {
		fmt.Printf("Removed function %s from %s\n", functionName, filename)
	}

	return result, nil
}

// createBackup creates a backup of the file
func (cr *CodeRemover) createBackup(filename string) (string, error) {
	if cr.dryRun {
		return "", nil
	}

	// Ensure backup directory exists
	err := os.MkdirAll(cr.backupDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup filename
	baseName := filepath.Base(filename)
	backupName := fmt.Sprintf("%s.backup", baseName)
	backupPath := filepath.Join(cr.backupDir, backupName)

	// Copy file to backup location
	src, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	err = os.WriteFile(backupPath, src, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}

// removeUnusedImportsFromAST removes unused imports from an AST
func (cr *CodeRemover) removeUnusedImportsFromAST(file *ast.File) bool {
	if len(file.Imports) == 0 {
		return false
	}

	// Collect all identifiers used in the file
	usedIdents := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			usedIdents[ident.Name] = true
		}
		return true
	})

	// Check each import
	var newImports []*ast.ImportSpec
	modified := false

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)

		// Get package name
		var pkgName string
		if imp.Name != nil {
			pkgName = imp.Name.Name
		} else {
			parts := strings.Split(importPath, "/")
			pkgName = parts[len(parts)-1]
		}

		// Skip dot imports and blank imports
		if pkgName == "." || pkgName == "_" {
			newImports = append(newImports, imp)
			continue
		}

		// Check if package is used
		if usedIdents[pkgName] {
			newImports = append(newImports, imp)
		} else {
			modified = true
			if cr.verbose {
				fmt.Printf("  Removing unused import: %s\n", importPath)
			}
		}
	}

	if modified {
		file.Imports = newImports

		// Also need to update the GenDecl that contains the imports
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				var newSpecs []ast.Spec
				for _, spec := range genDecl.Specs {
					if impSpec, ok := spec.(*ast.ImportSpec); ok {
						for _, newImp := range newImports {
							if impSpec == newImp {
								newSpecs = append(newSpecs, spec)
								break
							}
						}
					}
				}
				genDecl.Specs = newSpecs
			}
		}
	}

	return modified
}

// removeFunctionFromAST removes a function from an AST
func (cr *CodeRemover) removeFunctionFromAST(file *ast.File, functionName string) bool {
	var newDecls []ast.Decl
	modified := false

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name != nil && funcDecl.Name.Name == functionName {
				modified = true
				if cr.verbose {
					fmt.Printf("  Removing function: %s\n", functionName)
				}
				continue // Skip this declaration
			}
		}
		newDecls = append(newDecls, decl)
	}

	if modified {
		file.Decls = newDecls
	}

	return modified
}

// writeASTToFile writes an AST back to a file
func (cr *CodeRemover) writeASTToFile(file *ast.File, filename string) error {
	// Format the AST
	var buf strings.Builder
	err := format.Node(&buf, cr.fileSet, file)
	if err != nil {
		return fmt.Errorf("failed to format AST: %w", err)
	}

	// Write to file
	err = os.WriteFile(filename, []byte(buf.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// readFileLines reads a file and returns its lines
func (cr *CodeRemover) readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeFileLines writes lines to a file
func (cr *CodeRemover) writeFileLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// Example usage function
func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run code_remover.go <operation> <file> <args...>")
		fmt.Println("Operations:")
		fmt.Println("  imports <file>                    - Remove unused imports")
		fmt.Println("  comments <file> <start> <end>     - Remove commented code block")
		fmt.Println("  function <file> <name> <line>     - Remove unused function")
		os.Exit(1)
	}

	operation := os.Args[1]
	filename := os.Args[2]

	remover := NewCodeRemover(".dead_code_backups", false, true)

	var result *RemovalResult
	var err error

	switch operation {
	case "imports":
		result, err = remover.RemoveUnusedImports(filename)
	case "comments":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run code_remover.go comments <file> <start_line> <end_line>")
			os.Exit(1)
		}
		var startLine, endLine int
		fmt.Sscanf(os.Args[3], "%d", &startLine)
		fmt.Sscanf(os.Args[4], "%d", &endLine)
		result, err = remover.RemoveCommentedCodeBlock(filename, startLine, endLine)
	case "function":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run code_remover.go function <file> <function_name> <line_number>")
			os.Exit(1)
		}
		functionName := os.Args[3]
		var lineNumber int
		fmt.Sscanf(os.Args[4], "%d", &lineNumber)
		result, err = remover.RemoveUnusedFunction(filename, functionName, lineNumber)
	default:
		fmt.Printf("Unknown operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if result.Success {
		fmt.Printf("Successfully removed %s from %s\n", result.ItemName, result.File)
		if result.BackupFile != "" {
			fmt.Printf("Backup created: %s\n", result.BackupFile)
		}
	} else {
		fmt.Printf("Failed to remove %s from %s: %v\n", result.ItemName, result.File, result.Error)
		os.Exit(1)
	}
}
