package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// compileGoFile attempts to validate Go file compilation focusing on standard library imports
func (v *TemplateValidator) compileGoFile(filePath string) error {
	// Skip go.mod files as they don't need compilation
	if strings.HasSuffix(filePath, "go.mod") {
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Check if this file has only standard library imports
	hasOnlyStdLib := v.hasOnlyStandardLibraryImports(string(content))
	if !hasOnlyStdLib {
		// Skip compilation for files with non-standard imports
		// as we can't resolve project-specific dependencies
		return nil
	}

	// Create a temporary directory for compilation
	tempDir := filepath.Join(v.outputDir, "compile-test")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Copy the file to temp directory with a simple name
	tempFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Create a simple go.mod for the temp directory
	goModContent := `module temp-validation
go 1.22
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	// Try to compile the file
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()

	// Clean up temp directory
	os.RemoveAll(tempDir)

	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	return nil
}

// hasOnlyStandardLibraryImports checks if a Go file only imports standard library packages
func (v *TemplateValidator) hasOnlyStandardLibraryImports(content string) bool {
	lines := strings.Split(content, "\n")
	inImportBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Handle single-line imports
		if strings.HasPrefix(line, "import \"") {
			importPath := strings.Trim(strings.TrimPrefix(line, "import "), "\"")
			if !v.isStandardLibraryImport(importPath) {
				return false
			}
			continue
		}

		// Handle import blocks
		if strings.HasPrefix(line, "import (") {
			inImportBlock = true
			continue
		}

		if inImportBlock {
			if line == ")" {
				inImportBlock = false
				continue
			}

			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}

			// Extract import path - handle both quoted and unquoted
			// Remove any leading/trailing whitespace and quotes
			importPath := strings.Trim(line, " \t\"")

			// Handle aliased imports (alias "path")
			parts := strings.Fields(importPath)
			if len(parts) > 1 {
				// Last part is the import path, remove quotes
				importPath = strings.Trim(parts[len(parts)-1], "\"")
			}

			// Skip empty import paths
			if importPath == "" {
				continue
			}

			if !v.isStandardLibraryImport(importPath) {
				return false
			}
		}
	}

	return true
}

// isStandardLibraryImport checks if an import path is from the standard library
func (v *TemplateValidator) isStandardLibraryImport(importPath string) bool {
	// Standard library packages don't contain dots (except for some special cases)
	// and don't start with known third-party prefixes

	// Common third-party prefixes
	thirdPartyPrefixes := []string{
		"github.com/",
		"gitlab.com/",
		"bitbucket.org/",
		"golang.org/x/",
		"google.golang.org/",
		"gopkg.in/",
		"go.uber.org/",
	}

	for _, prefix := range thirdPartyPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return false
		}
	}

	// If it contains a dot and doesn't start with known third-party prefixes,
	// it's likely a project-specific import
	if strings.Contains(importPath, ".") && !strings.HasPrefix(importPath, "golang.org/") {
		return false
	}

	// Standard library packages
	stdLibPackages := []string{
		"bufio", "bytes", "context", "crypto", "database", "encoding", "errors",
		"fmt", "go", "hash", "html", "image", "io", "log", "math", "mime", "net",
		"os", "path", "reflect", "regexp", "runtime", "sort", "strconv", "strings",
		"sync", "syscall", "testing", "text", "time", "unicode", "unsafe",
	}

	// Check if it's a known standard library package or subpackage
	for _, pkg := range stdLibPackages {
		if importPath == pkg || strings.HasPrefix(importPath, pkg+"/") {
			return true
		}
	}

	// If it contains slashes but isn't in our standard library list, it's likely project-specific
	if strings.Contains(importPath, "/") {
		return false
	}

	// If it doesn't contain dots or slashes and isn't in our known list, it's likely standard library
	return !strings.Contains(importPath, ".")
}

// Enhanced validateGoFile that includes compilation check
func (v *TemplateValidator) validateGoFileWithCompilation(filePath string) error {
	// First, do syntax validation
	if err := v.validateGoFile(filePath); err != nil {
		return err
	}

	// Then, try to compile the file
	if err := v.compileGoFile(filePath); err != nil {
		return err
	}

	return nil
}
