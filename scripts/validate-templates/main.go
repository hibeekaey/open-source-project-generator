package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// TemplateValidator validates template compilation
type TemplateValidator struct {
	templatesDir string
	outputDir    string
	testData     *models.ProjectConfig
}

// ValidationReport contains the results of template validation
type ValidationReport struct {
	TotalTemplates    int
	ValidTemplates    int
	InvalidTemplates  int
	CompilationErrors []CompilationError
	Warnings          []string
}

// CompilationError represents a compilation error in a generated file
type CompilationError struct {
	TemplatePath  string
	GeneratedFile string
	Error         string
	Line          int
	Column        int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <templates-directory-or-file> [output-directory]")
		os.Exit(1)
	}

	templatesPath := os.Args[1]
	outputDir := "validation-output"
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	// Determine if we're validating a single file or directory
	info, err := os.Stat(templatesPath)
	if err != nil {
		log.Fatalf("Failed to access path %s: %v", templatesPath, err)
	}

	var templatesDir string
	if info.IsDir() {
		templatesDir = templatesPath
	} else {
		// For single files, use the directory containing the file
		templatesDir = filepath.Dir(templatesPath)
	}

	validator := NewTemplateValidator(templatesDir, outputDir)

	fmt.Println("🔍 Starting template validation...")

	var report *ValidationReport
	if info.IsDir() {
		report, err = validator.ValidateAllTemplates()
	} else {
		report, err = validator.ValidateSingleTemplate(templatesPath)
	}

	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	validator.PrintReport(report)

	if report.InvalidTemplates > 0 {
		os.Exit(1)
	}
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator(templatesDir, outputDir string) *TemplateValidator {
	return &TemplateValidator{
		templatesDir: templatesDir,
		outputDir:    outputDir,
		testData:     createTestData(),
	}
}

// createTestData creates sample data for template variables
func createTestData() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:         "testproject",
		Organization: "testorg",
		Description:  "A test project for template validation",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/testorg/testproject",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    true,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.0.0",
			Go:     "1.22",
			Kotlin: "1.9.0",
			Swift:  "5.9",
			NextJS: "14.0.0",
			React:  "18.0.0",
			Packages: map[string]string{
				"typescript": "5.0.0",
				"eslint":     "8.0.0",
			},
			UpdatedAt: time.Now(),
		},
		CustomVars: map[string]string{
			"DATABASE_URL": "postgresql://localhost:5432/testdb",
			"REDIS_URL":    "redis://localhost:6379",
		},
		OutputPath:       "output",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}
}

// ValidateAllTemplates validates all template files in the templates directory
func (v *TemplateValidator) ValidateAllTemplates() (*ValidationReport, error) {
	report := &ValidationReport{
		CompilationErrors: make([]CompilationError, 0),
		Warnings:          make([]string, 0),
	}

	// Clean and create output directory
	if err := os.RemoveAll(v.outputDir); err != nil {
		return nil, fmt.Errorf("failed to clean output directory: %w", err)
	}
	if err := os.MkdirAll(v.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Walk through all template files
	err := filepath.Walk(v.templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Skip non-Go template files for compilation validation
		if !v.isGoTemplate(path) {
			report.TotalTemplates++
			report.ValidTemplates++ // Non-Go templates are considered valid for this validation
			return nil
		}

		report.TotalTemplates++

		fmt.Printf("📝 Validating: %s\n", path)

		if err := v.validateGoTemplate(path, report); err != nil {
			report.InvalidTemplates++
			fmt.Printf("❌ Failed: %s - %v\n", path, err)
		} else {
			report.ValidTemplates++
			fmt.Printf("✅ Valid: %s\n", path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}

	return report, nil
}

// isGoTemplate checks if a template file is a Go source file
func (v *TemplateValidator) isGoTemplate(templatePath string) bool {
	// Remove .tmpl extension and check if it's a Go file
	baseName := strings.TrimSuffix(templatePath, ".tmpl")
	return strings.HasSuffix(baseName, ".go") || strings.HasSuffix(baseName, ".mod")
}

// validateGoTemplate validates a single Go template file
func (v *TemplateValidator) validateGoTemplate(templatePath string, report *ValidationReport) error {
	// Generate the Go file from template
	generatedFile, err := v.generateFromTemplate(templatePath)
	if err != nil {
		report.CompilationErrors = append(report.CompilationErrors, CompilationError{
			TemplatePath:  templatePath,
			GeneratedFile: "",
			Error:         fmt.Sprintf("Template generation failed: %v", err),
		})
		return err
	}

	// Validate the generated Go file (syntax + compilation)
	if err := v.validateGoFileWithCompilation(generatedFile); err != nil {
		report.CompilationErrors = append(report.CompilationErrors, CompilationError{
			TemplatePath:  templatePath,
			GeneratedFile: generatedFile,
			Error:         err.Error(),
		})
		return err
	}

	return nil
}

// generateFromTemplate generates a Go file from a template
func (v *TemplateValidator) generateFromTemplate(templatePath string) (string, error) {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file path
	relPath, err := filepath.Rel(v.templatesDir, templatePath)
	if err != nil {
		// If we can't get relative path, use the base name
		relPath = filepath.Base(templatePath)
	}

	// Remove .tmpl extension and create output path
	outputPath := filepath.Join(v.outputDir, strings.TrimSuffix(relPath, ".tmpl"))

	// Create output directory
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Execute template with test data
	if err := tmpl.Execute(outputFile, v.testData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outputPath, nil
}

// validateGoFile validates a generated Go file by attempting to parse it
func (v *TemplateValidator) validateGoFile(filePath string) error {
	// For go.mod files, we just check if they're valid module files
	if strings.HasSuffix(filePath, "go.mod") {
		return v.validateGoMod(filePath)
	}

	// For .go files, we use go/parser to validate syntax
	return v.validateGoSource(filePath)
}

// validateGoMod validates a go.mod file
func (v *TemplateValidator) validateGoMod(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %w", err)
	}

	// Basic validation - check if it starts with module declaration
	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "module ") {
		return fmt.Errorf("invalid go.mod file: missing module declaration")
	}

	return nil
}

// validateGoSource validates a Go source file using go/parser
func (v *TemplateValidator) validateGoSource(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read Go file: %w", err)
	}

	// Use go/parser to validate syntax
	_, err = parseGoFile(filePath, content)
	if err != nil {
		return fmt.Errorf("Go syntax error: %w", err)
	}

	return nil
}

// PrintReport prints the validation report
func (v *TemplateValidator) PrintReport(report *ValidationReport) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 TEMPLATE VALIDATION REPORT")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("📁 Total Templates: %d\n", report.TotalTemplates)
	fmt.Printf("✅ Valid Templates: %d\n", report.ValidTemplates)
	fmt.Printf("❌ Invalid Templates: %d\n", report.InvalidTemplates)

	if len(report.CompilationErrors) > 0 {
		fmt.Println("\n🚨 COMPILATION ERRORS:")
		fmt.Println(strings.Repeat("-", 40))

		for i, err := range report.CompilationErrors {
			fmt.Printf("\n%d. Template: %s\n", i+1, err.TemplatePath)
			if err.GeneratedFile != "" {
				fmt.Printf("   Generated: %s\n", err.GeneratedFile)
			}
			fmt.Printf("   Error: %s\n", err.Error)
		}
	}

	if len(report.Warnings) > 0 {
		fmt.Println("\n⚠️  WARNINGS:")
		fmt.Println(strings.Repeat("-", 40))

		for i, warning := range report.Warnings {
			fmt.Printf("%d. %s\n", i+1, warning)
		}
	}

	fmt.Println(strings.Repeat("=", 60))

	if report.InvalidTemplates == 0 {
		fmt.Println("🎉 All templates are valid!")
	} else {
		fmt.Printf("⚠️  Found %d invalid templates that need fixing.\n", report.InvalidTemplates)
	}
}

// ValidateSingleTemplate validates a single template file
func (v *TemplateValidator) ValidateSingleTemplate(templatePath string) (*ValidationReport, error) {
	report := &ValidationReport{
		CompilationErrors: make([]CompilationError, 0),
		Warnings:          make([]string, 0),
	}

	// Clean and create output directory
	if err := os.RemoveAll(v.outputDir); err != nil {
		return nil, fmt.Errorf("failed to clean output directory: %w", err)
	}
	if err := os.MkdirAll(v.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if it's a template file
	if !strings.HasSuffix(templatePath, ".tmpl") {
		return nil, fmt.Errorf("file %s is not a template file (.tmpl)", templatePath)
	}

	// Skip non-Go template files for compilation validation
	if !v.isGoTemplate(templatePath) {
		report.TotalTemplates = 1
		report.ValidTemplates = 1
		fmt.Printf("📝 Validating: %s (non-Go template)\n", templatePath)
		fmt.Printf("✅ Valid: %s\n", templatePath)
		return report, nil
	}

	report.TotalTemplates = 1

	fmt.Printf("📝 Validating: %s\n", templatePath)

	if err := v.validateGoTemplate(templatePath, report); err != nil {
		report.InvalidTemplates = 1
		fmt.Printf("❌ Failed: %s - %v\n", templatePath, err)
	} else {
		report.ValidTemplates = 1
		fmt.Printf("✅ Valid: %s\n", templatePath)
	}

	return report, nil
}
