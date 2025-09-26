package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// VerificationStatus represents the status of template verification
type VerificationStatus int

const (
	VerificationSuccess VerificationStatus = iota
	VerificationFailed
	VerificationSkipped
)

// TemplateVerificationResult represents the result of template verification
type TemplateVerificationResult struct {
	TemplatePath      string
	Status            VerificationStatus
	Error             string
	CompilationOutput string
	GeneratedFile     string
}

// verifyTemplateCompilation verifies that a template generates compilable Go code
func verifyTemplateCompilation(templatePath string, testData *models.ProjectConfig, outputDir string) TemplateVerificationResult {
	result := TemplateVerificationResult{
		TemplatePath: templatePath,
		Status:       VerificationFailed,
	}

	// Read template content with path validation
	content, err := utils.SafeReadFile(templatePath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read template: %v", err)
		return result
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		result.Error = fmt.Sprintf("failed to parse template: %v", err)
		return result
	}

	// Create output file path
	outputFileName := strings.TrimSuffix(filepath.Base(templatePath), ".tmpl")
	outputPath := filepath.Join(outputDir, outputFileName)
	result.GeneratedFile = outputPath

	// Create output directory with secure permissions
	if mkdirErr := utils.SafeMkdirAll(filepath.Dir(outputPath)); mkdirErr != nil {
		result.Error = fmt.Sprintf("failed to create output directory: %v", mkdirErr)
		return result
	}

	// Create output file with secure permissions
	outputFile, err := utils.SafeCreate(outputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create output file: %v", err)
		return result
	}
	defer func() { _ = outputFile.Close() }()

	// Execute template
	err = tmpl.Execute(outputFile, testData)
	if err != nil {
		result.Error = fmt.Sprintf("failed to execute template: %v", err)
		return result
	}

	// Verify the generated file
	if err := verifyGeneratedFile(outputPath); err != nil {
		result.Error = err.Error()
		return result
	}

	// Try to compile if it's a Go file with only standard library imports
	if strings.HasSuffix(outputPath, ".go") {
		if compilationOutput, err := tryCompileGoFile(outputPath); err != nil {
			result.Error = fmt.Sprintf("compilation failed: %v", err)
			result.CompilationOutput = compilationOutput
			return result
		}
	}

	result.Status = VerificationSuccess
	return result
}

// verifyGeneratedFile performs basic verification on the generated file
func verifyGeneratedFile(filePath string) error {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read generated file: %w", err)
	}

	contentStr := string(content)

	// Check for remaining template syntax
	if strings.Contains(contentStr, "{{") || strings.Contains(contentStr, "}}") {
		return fmt.Errorf("generated file contains unprocessed template syntax")
	}

	// For Go files, perform additional checks
	if strings.HasSuffix(filePath, ".go") {
		return verifyGoFile(contentStr)
	}

	// For go.mod files, check basic structure
	if strings.HasSuffix(filePath, "go.mod") {
		return verifyGoModFile(contentStr)
	}

	return nil
}

// verifyGoFile performs Go-specific verification
func verifyGoFile(content string) error {
	lines := strings.Split(content, "\n")

	// Check for package declaration
	hasPackage := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			hasPackage = true
			break
		}
	}

	if !hasPackage {
		return fmt.Errorf("generated Go file missing package declaration")
	}

	// Check for basic Go syntax issues
	if strings.Contains(content, "template_placeholder") {
		return fmt.Errorf("generated file contains template placeholders")
	}

	return nil
}

// verifyGoModFile performs go.mod-specific verification
func verifyGoModFile(content string) error {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "module ") {
		return fmt.Errorf("invalid go.mod file: missing module declaration")
	}
	return nil
}

// tryCompileGoFile attempts to compile a Go file
func tryCompileGoFile(filePath string) (string, error) {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Check if this file has only standard library imports
	if !hasOnlyStandardLibraryImports(string(content)) {
		// Skip compilation for files with non-standard imports
		return "", nil
	}

	// Create a temporary directory for compilation with secure permissions
	tempDir := filepath.Dir(filePath) + "_compile_test"
	if mkdirErr := utils.SafeMkdirAll(tempDir); mkdirErr != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", mkdirErr)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Copy the file to temp directory with secure permissions
	tempFile := filepath.Join(tempDir, "main.go")
	if writeErr := utils.SafeWriteFile(tempFile, content); writeErr != nil {
		return "", fmt.Errorf("failed to write temp file: %w", writeErr)
	}

	// Create a basic go.mod for the temp directory
	goModContent := `module temp-validation
go 1.25
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if goModErr := utils.SafeWriteFile(goModPath, []byte(goModContent)); goModErr != nil {
		return "", fmt.Errorf("failed to create go.mod: %w", goModErr)
	}

	// Try to compile the file
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("compilation failed: %s", string(output))
	}

	return string(output), nil
}

// hasOnlyStandardLibraryImports checks if a Go file only imports standard library packages
func hasOnlyStandardLibraryImports(content string) bool {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this is an import line
		if strings.HasPrefix(line, "import") || (strings.Contains(line, "import") && strings.Contains(line, "\"")) {
			continue // Skip import statement line itself
		}

		// Check for actual import paths in quotes
		if strings.Contains(line, "\"") && (strings.HasPrefix(strings.TrimSpace(line), "\"") || strings.Contains(line, " \"")) {
			// Extract the import path
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && start != end {
				importPath := line[start+1 : end]

				// Check for non-standard library imports
				if strings.Contains(importPath, "github.com") ||
					strings.Contains(importPath, "golang.org") ||
					strings.Contains(importPath, "gopkg.in") ||
					strings.Contains(importPath, ".com/") ||
					strings.Contains(importPath, ".org/") ||
					strings.Contains(importPath, "testproject/") ||
					strings.Contains(importPath, "project/") ||
					strings.Contains(importPath, "/internal/") ||
					(strings.Contains(importPath, "/") && !isStandardLibraryPackage(importPath)) {
					return false
				}
			}
		}
	}

	return true
}

// isStandardLibraryPackage checks if a package is part of Go's standard library
func isStandardLibraryPackage(pkg string) bool {
	// List of common standard library packages with slashes
	standardPackages := []string{
		"net/http", "encoding/json", "text/template", "html/template",
		"path/filepath", "crypto/rand", "crypto/md5", "crypto/sha256",
		"encoding/base64", "encoding/xml", "net/url", "os/exec",
		"database/sql", "context", "time", "fmt", "strings", "strconv",
		"io", "bufio", "bytes", "regexp", "sort", "math", "errors",
	}

	// If it has no slash, it's likely a standard package (fmt, os, etc.)
	if !strings.Contains(pkg, "/") {
		return true
	}

	// Check against known standard packages
	for _, stdPkg := range standardPackages {
		if pkg == stdPkg {
			return true
		}
	}

	// If it starts with known third-party patterns, it's not standard
	if strings.HasPrefix(pkg, "github.com") ||
		strings.HasPrefix(pkg, "golang.org") ||
		strings.HasPrefix(pkg, "gopkg.in") ||
		strings.Contains(pkg, ".com/") ||
		strings.Contains(pkg, ".org/") {
		return false
	}

	// For other patterns with slashes, assume they're project-specific
	return false
}
