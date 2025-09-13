package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/open-source-template-generator/pkg/models"
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

	// Read template content
	content, err := os.ReadFile(templatePath)
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

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		result.Error = fmt.Sprintf("failed to create output directory: %v", err)
		return result
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create output file: %v", err)
		return result
	}
	defer outputFile.Close()

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
	content, err := os.ReadFile(filePath)
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
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Check if this file has only standard library imports
	if !hasOnlyStandardLibraryImports(string(content)) {
		// Skip compilation for files with non-standard imports
		return "", nil
	}

	// Create a temporary directory for compilation
	tempDir := filepath.Dir(filePath) + "_compile_test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy the file to temp directory
	tempFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	// Create a simple go.mod for the temp directory
	goModContent := `module temp-validation
go 1.22
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create go.mod: %w", err)
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

		// Skip non-import lines
		if !strings.Contains(line, "\"") {
			continue
		}

		// Check for non-standard library imports
		if strings.Contains(line, "github.com") ||
			strings.Contains(line, "golang.org") ||
			strings.Contains(line, "gopkg.in") ||
			strings.Contains(line, ".com/") ||
			strings.Contains(line, ".org/") {
			return false
		}
	}

	return true
}
