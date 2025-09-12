package standards

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TemplateUpdater interface defines methods for updating template files
type TemplateUpdater interface {
	UpdateTemplate(templatePath, templateType string, versions map[string]string) error
	UpdateAllTemplates(templatesDir string, versions map[string]string) error
	ApplyVersionUpdates(templatePath string, versions map[string]string) error
	BackupTemplate(templatePath string) (string, error)
	RestoreTemplate(templatePath, backupPath string) error
}

// StandardTemplateUpdater implements the TemplateUpdater interface
type StandardTemplateUpdater struct {
	generator *TemplateGenerator
	validator *TemplateValidator
	standards *FrontendStandards
}

// NewTemplateUpdater creates a new template updater
func NewTemplateUpdater() TemplateUpdater {
	return &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}
}

// UpdateTemplate updates a single template with standardized configurations and version updates
func (tu *StandardTemplateUpdater) UpdateTemplate(templatePath, templateType string, versions map[string]string) error {
	// Validate template type
	if !isValidTemplateType(templateType) {
		return fmt.Errorf("invalid template type: %s", templateType)
	}

	// Create backup before updating
	backupPath, err := tu.BackupTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("📦 Created backup at %s\n", backupPath)

	// Apply standardized configurations
	if err := tu.applyStandardizedConfigurations(templatePath, templateType); err != nil {
		// Restore from backup on failure
		if restoreErr := tu.RestoreTemplate(templatePath, backupPath); restoreErr != nil {
			return fmt.Errorf("update failed and restore failed: %w (original error: %v)", restoreErr, err)
		}
		return fmt.Errorf("failed to apply standardized configurations: %w", err)
	}

	// Apply version updates
	if err := tu.ApplyVersionUpdates(templatePath, versions); err != nil {
		// Restore from backup on failure
		if restoreErr := tu.RestoreTemplate(templatePath, backupPath); restoreErr != nil {
			return fmt.Errorf("version update failed and restore failed: %w (original error: %v)", restoreErr, err)
		}
		return fmt.Errorf("failed to apply version updates: %w", err)
	}

	// Validate the updated template
	result, err := tu.validator.ValidateTemplate(templatePath, templateType)
	if err != nil {
		return fmt.Errorf("failed to validate updated template: %w", err)
	}

	if !result.IsValid {
		fmt.Printf("⚠️  Template %s has validation issues after update:\n", templateType)
		for _, error := range result.Errors {
			fmt.Printf("  - %s: %s\n", error.Type, error.Description)
		}
	} else {
		fmt.Printf("✅ Template %s updated and validated successfully\n", templateType)
	}

	return nil
}

// UpdateAllTemplates updates all frontend templates
func (tu *StandardTemplateUpdater) UpdateAllTemplates(templatesDir string, versions map[string]string) error {
	frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	for _, templateType := range frontendTemplates {
		templatePath := filepath.Join(templatesDir, "frontend", templateType)

		fmt.Printf("🔄 Updating template: %s\n", templateType)

		if err := tu.UpdateTemplate(templatePath, templateType, versions); err != nil {
			return fmt.Errorf("failed to update template %s: %w", templateType, err)
		}

		fmt.Printf("✅ Successfully updated %s\n\n", templateType)
	}

	return nil
}

// ApplyVersionUpdates applies version updates to template files
func (tu *StandardTemplateUpdater) ApplyVersionUpdates(templatePath string, versions map[string]string) error {
	// Find all template files that need version updates
	templateFiles, err := findTemplateFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to find template files: %w", err)
	}

	for _, filePath := range templateFiles {
		if err := tu.updateVersionsInFile(filePath, versions); err != nil {
			return fmt.Errorf("failed to update versions in %s: %w", filePath, err)
		}
	}

	return nil
}

// BackupTemplate creates a backup of the template directory
func (tu *StandardTemplateUpdater) BackupTemplate(templatePath string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := templatePath + ".backup." + timestamp

	if err := copyDir(templatePath, backupPath); err != nil {
		return "", fmt.Errorf("failed to copy template directory: %w", err)
	}

	return backupPath, nil
}

// RestoreTemplate restores a template from backup
func (tu *StandardTemplateUpdater) RestoreTemplate(templatePath, backupPath string) error {
	// Remove current template directory
	if err := os.RemoveAll(templatePath); err != nil {
		return fmt.Errorf("failed to remove current template: %w", err)
	}

	// Restore from backup
	if err := copyDir(backupPath, templatePath); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	return nil
}

// applyStandardizedConfigurations applies standardized configurations to a template
func (tu *StandardTemplateUpdater) applyStandardizedConfigurations(templatePath, templateType string) error {
	// Generate standardized files
	files := tu.generator.GenerateAllStandardizedFiles(templateType)
	filePaths := GetStandardizedFilePaths()

	// Write each standardized file
	fileMap := map[string]string{
		"PackageJSON":    files.PackageJSON,
		"TSConfig":       files.TSConfig,
		"ESLintConfig":   files.ESLintConfig,
		"PrettierConfig": files.PrettierConfig,
		"VercelConfig":   files.VercelConfig,
		"TailwindConfig": files.TailwindConfig,
		"NextConfig":     files.NextConfig,
		"PostCSSConfig":  files.PostCSSConfig,
		"JestConfig":     files.JestConfig,
		"JestSetup":      files.JestSetup,
	}

	for fileType, content := range fileMap {
		fileName := filePaths[fileType]
		filePath := filepath.Join(templatePath, fileName)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fileName, err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}

		fmt.Printf("  ✅ Updated %s\n", fileName)
	}

	return nil
}

// updateVersionsInFile updates version placeholders in a template file
func (tu *StandardTemplateUpdater) updateVersionsInFile(filePath string, versions map[string]string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	updated := false

	// Replace version placeholders
	for versionKey, versionValue := range versions {
		placeholder := fmt.Sprintf("{{.Versions.%s}}", versionKey)
		if strings.Contains(contentStr, placeholder) {
			contentStr = strings.ReplaceAll(contentStr, placeholder, versionValue)
			updated = true
		}
	}

	// Write back if updated
	if updated {
		if err := os.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write updated file: %w", err)
		}
		fmt.Printf("  🔄 Updated versions in %s\n", filepath.Base(filePath))
	}

	return nil
}

// VersionUpdateRule defines a rule for updating versions
type VersionUpdateRule struct {
	Pattern     string   // Pattern to match (e.g., "{{.Versions.NextJS}}")
	Replacement string   // Replacement value
	FileTypes   []string // File types to apply to (e.g., [".json", ".js"])
}

// StandardizationRule defines a rule for standardizing configurations
type StandardizationRule struct {
	Name        string
	Description string
	FilePattern string
	Validator   func(content string) (bool, []string) // Returns (isValid, issues)
	Fixer       func(content string) (string, error)  // Returns fixed content
}

// GetStandardizationRules returns all standardization rules
func (tu *StandardTemplateUpdater) GetStandardizationRules() []StandardizationRule {
	return []StandardizationRule{
		{
			Name:        "PackageJSONScripts",
			Description: "Ensure package.json has all required scripts",
			FilePattern: "package.json.tmpl",
			Validator:   tu.validatePackageJSONScripts,
			Fixer:       tu.fixPackageJSONScripts,
		},
		{
			Name:        "TypeScriptConfig",
			Description: "Ensure tsconfig.json has standardized compiler options",
			FilePattern: "tsconfig.json.tmpl",
			Validator:   tu.validateTypeScriptConfig,
			Fixer:       tu.fixTypeScriptConfig,
		},
		{
			Name:        "ESLintConfig",
			Description: "Ensure .eslintrc.json has standardized rules",
			FilePattern: ".eslintrc.json.tmpl",
			Validator:   tu.validateESLintConfig,
			Fixer:       tu.fixESLintConfig,
		},
		{
			Name:        "PrettierConfig",
			Description: "Ensure .prettierrc has standardized formatting options",
			FilePattern: ".prettierrc.tmpl",
			Validator:   tu.validatePrettierConfig,
			Fixer:       tu.fixPrettierConfig,
		},
		{
			Name:        "VercelConfig",
			Description: "Ensure vercel.json has standardized deployment settings",
			FilePattern: "vercel.json.tmpl",
			Validator:   tu.validateVercelConfig,
			Fixer:       tu.fixVercelConfig,
		},
	}
}

// ApplyStandardizationRules applies all standardization rules to a template
func (tu *StandardTemplateUpdater) ApplyStandardizationRules(templatePath string) error {
	rules := tu.GetStandardizationRules()

	for _, rule := range rules {
		filePath := filepath.Join(templatePath, rule.FilePattern)

		if !fileExists(filePath) {
			fmt.Printf("  ⚠️  File %s not found, skipping rule %s\n", rule.FilePattern, rule.Name)
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filePath, err)
		}

		isValid, issues := rule.Validator(string(content))
		if !isValid {
			fmt.Printf("  🔧 Applying rule %s to %s\n", rule.Name, rule.FilePattern)
			for _, issue := range issues {
				fmt.Printf("    - %s\n", issue)
			}

			fixedContent, err := rule.Fixer(string(content))
			if err != nil {
				return fmt.Errorf("failed to fix %s with rule %s: %w", filePath, rule.Name, err)
			}

			if err := os.WriteFile(filePath, []byte(fixedContent), 0644); err != nil {
				return fmt.Errorf("failed to write fixed %s: %w", filePath, err)
			}

			fmt.Printf("  ✅ Applied rule %s to %s\n", rule.Name, rule.FilePattern)
		}
	}

	return nil
}

// Helper functions for validation and fixing

func (tu *StandardTemplateUpdater) validatePackageJSONScripts(content string) (bool, []string) {
	// Simple validation - check if content contains required scripts
	requiredScripts := []string{"dev", "build", "start", "lint", "type-check", "test"}
	var issues []string

	for _, script := range requiredScripts {
		if !strings.Contains(content, fmt.Sprintf(`"%s":`, script)) {
			issues = append(issues, fmt.Sprintf("Missing required script: %s", script))
		}
	}

	return len(issues) == 0, issues
}

func (tu *StandardTemplateUpdater) fixPackageJSONScripts(content string) (string, error) {
	// For now, return the standardized package.json content
	// In a real implementation, this would merge existing content with standards
	return tu.generator.GenerateStandardizedPackageJSON("nextjs-app"), nil
}

func (tu *StandardTemplateUpdater) validateTypeScriptConfig(content string) (bool, []string) {
	var issues []string

	requiredOptions := []string{"strict", "noEmit", "jsx", "moduleResolution"}
	for _, option := range requiredOptions {
		if !strings.Contains(content, fmt.Sprintf(`"%s":`, option)) {
			issues = append(issues, fmt.Sprintf("Missing compiler option: %s", option))
		}
	}

	return len(issues) == 0, issues
}

func (tu *StandardTemplateUpdater) fixTypeScriptConfig(content string) (string, error) {
	return tu.generator.GenerateStandardizedTSConfig(), nil
}

func (tu *StandardTemplateUpdater) validateESLintConfig(content string) (bool, []string) {
	var issues []string

	if !strings.Contains(content, "next/core-web-vitals") {
		issues = append(issues, "Missing Next.js ESLint config")
	}

	if !strings.Contains(content, "@typescript-eslint/recommended") {
		issues = append(issues, "Missing TypeScript ESLint config")
	}

	return len(issues) == 0, issues
}

func (tu *StandardTemplateUpdater) fixESLintConfig(content string) (string, error) {
	return tu.generator.GenerateStandardizedESLintConfig(), nil
}

func (tu *StandardTemplateUpdater) validatePrettierConfig(content string) (bool, []string) {
	var issues []string

	requiredOptions := []string{"semi", "singleQuote", "printWidth", "tabWidth"}
	for _, option := range requiredOptions {
		if !strings.Contains(content, fmt.Sprintf(`"%s":`, option)) {
			issues = append(issues, fmt.Sprintf("Missing Prettier option: %s", option))
		}
	}

	return len(issues) == 0, issues
}

func (tu *StandardTemplateUpdater) fixPrettierConfig(content string) (string, error) {
	return tu.generator.GenerateStandardizedPrettierConfig(), nil
}

func (tu *StandardTemplateUpdater) validateVercelConfig(content string) (bool, []string) {
	var issues []string

	if !strings.Contains(content, `"framework": "nextjs"`) {
		issues = append(issues, "Missing or incorrect framework setting")
	}

	if !strings.Contains(content, `"buildCommand": "npm run build"`) {
		issues = append(issues, "Missing or incorrect build command")
	}

	return len(issues) == 0, issues
}

func (tu *StandardTemplateUpdater) fixVercelConfig(content string) (string, error) {
	return tu.generator.GenerateStandardizedVercelConfig(), nil
}

// Helper functions

// isValidTemplateType checks if a template type is valid
func isValidTemplateType(templateType string) bool {
	validTypes := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	for _, validType := range validTypes {
		if templateType == validType {
			return true
		}
	}
	return false
}

// findTemplateFiles finds all template files in a directory
func findTemplateFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tmpl") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = srcFile.WriteTo(dstFile)
		return err
	})
}
