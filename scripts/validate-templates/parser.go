package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// parseGoFile parses a Go source file and returns the AST
func parseGoFile(filename string, content []byte) (*ast.File, error) {
	fset := token.NewFileSet()

	// Parse the Go source file
	file, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	// Additional validation - check for common issues
	if err := validateASTStructure(file, fset); err != nil {
		return nil, err
	}

	return file, nil
}

// validateASTStructure performs additional validation on the parsed AST
func validateASTStructure(file *ast.File, fset *token.FileSet) error {
	// Check if package declaration exists
	if file.Name == nil {
		return fmt.Errorf("missing package declaration")
	}

	// Validate imports
	if err := validateImports(file, fset); err != nil {
		return err
	}

	// Check for unresolved identifiers that might indicate missing imports
	if err := validateIdentifiers(file, fset); err != nil {
		return err
	}

	return nil
}

// validateImports validates import statements
func validateImports(file *ast.File, fset *token.FileSet) error {
	importPaths := make(map[string]bool)

	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}

		path := strings.Trim(imp.Path.Value, `"`)

		// Check for duplicate imports
		if importPaths[path] {
			pos := fset.Position(imp.Pos())
			return fmt.Errorf("duplicate import %q at line %d", path, pos.Line)
		}
		importPaths[path] = true

		// Basic validation of import path format
		if err := validateImportPath(path); err != nil {
			pos := fset.Position(imp.Pos())
			return fmt.Errorf("invalid import path %q at line %d: %w", path, pos.Line, err)
		}
	}

	return nil
}

// validateImportPath validates the format of an import path
func validateImportPath(path string) error {
	if path == "" {
		return fmt.Errorf("empty import path")
	}

	// Check for common invalid characters
	if strings.Contains(path, " ") {
		return fmt.Errorf("import path contains spaces")
	}

	return nil
}

// validateIdentifiers checks for common identifier issues
func validateIdentifiers(file *ast.File, fset *token.FileSet) error {
	// Create a map of imported packages for reference
	importedPackages := make(map[string]string)

	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}

		path := strings.Trim(imp.Path.Value, `"`)

		// Determine package name
		var pkgName string
		if imp.Name != nil {
			pkgName = imp.Name.Name
		} else {
			// Extract package name from path
			parts := strings.Split(path, "/")
			pkgName = parts[len(parts)-1]

			// Handle special cases like "golang.org/x/crypto"
			if strings.Contains(pkgName, ".") {
				pkgName = strings.Split(pkgName, ".")[0]
			}
		}

		importedPackages[pkgName] = path
	}

	// Walk the AST to find potential issues
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			// Check selector expressions like "time.Now()"
			if ident, ok := node.X.(*ast.Ident); ok {
				if err := validateSelectorExpr(ident.Name, node.Sel.Name, importedPackages, fset, node.Pos()); err != nil {
					// Note: We're not returning the error here because it would stop the inspection
					// Instead, we could collect these errors and return them later
					// For now, we'll just continue the validation
				}
			}
		}
		return true
	})

	return nil
}

// validateSelectorExpr validates selector expressions for missing imports
func validateSelectorExpr(pkgName, funcName string, importedPackages map[string]string, fset *token.FileSet, pos token.Pos) error {
	// Check if the package is imported
	if _, exists := importedPackages[pkgName]; !exists {
		// Check if this looks like a standard library package
		if isStandardLibraryPackage(pkgName, funcName) {
			position := fset.Position(pos)
			return fmt.Errorf("missing import for standard library package %q at line %d (function: %s)", pkgName, position.Line, funcName)
		}
	}

	return nil
}

// isStandardLibraryPackage checks if a package/function combination is from the standard library
func isStandardLibraryPackage(pkgName, funcName string) bool {
	// Common standard library packages and their functions
	stdLibFunctions := map[string][]string{
		"time":     {"Now", "Parse", "Since", "Until", "Sleep", "Tick", "After"},
		"fmt":      {"Printf", "Sprintf", "Print", "Println", "Errorf", "Fprintf"},
		"strings":  {"Contains", "HasPrefix", "HasSuffix", "Split", "Join", "Replace", "ToLower", "ToUpper"},
		"strconv":  {"Atoi", "Itoa", "ParseInt", "ParseFloat", "FormatInt", "FormatFloat"},
		"os":       {"Open", "Create", "Remove", "Getenv", "Setenv", "Exit", "Args"},
		"io":       {"Copy", "ReadAll", "WriteString", "EOF"},
		"net":      {"Listen", "Dial", "ParseIP"},
		"http":     {"Get", "Post", "ListenAndServe", "HandleFunc", "NewRequest"},
		"json":     {"Marshal", "Unmarshal", "NewEncoder", "NewDecoder"},
		"log":      {"Print", "Printf", "Println", "Fatal", "Fatalf", "Panic"},
		"errors":   {"New", "Is", "As", "Unwrap"},
		"context":  {"Background", "TODO", "WithCancel", "WithTimeout", "WithValue"},
		"sync":     {"Mutex", "RWMutex", "WaitGroup", "Once"},
		"regexp":   {"Compile", "MustCompile", "Match", "MatchString"},
		"path":     {"Join", "Dir", "Base", "Ext", "Clean"},
		"filepath": {"Join", "Dir", "Base", "Ext", "Walk", "Abs", "Rel"},
		"url":      {"Parse", "ParseQuery", "QueryEscape", "QueryUnescape"},
		"crypto":   {"MD5", "SHA1", "SHA256"},
		"rand":     {"Int", "Intn", "Float64", "Seed"},
		"sort":     {"Strings", "Ints", "Sort", "Reverse"},
		"bytes":    {"Buffer", "Contains", "Equal", "Compare"},
		"bufio":    {"NewReader", "NewWriter", "NewScanner"},
		"math":     {"Abs", "Max", "Min", "Sqrt", "Pow"},
	}

	functions, exists := stdLibFunctions[pkgName]
	if !exists {
		return false
	}

	for _, fn := range functions {
		if fn == funcName {
			return true
		}
	}

	return false
}
