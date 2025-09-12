package version

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// TemplateUpdater implements the TemplateUpdater interface
type TemplateUpdater struct {
	backupDir string
}

// NewTemplateUpdater creates a new template updater
func NewTemplateUpdater(backupDir string) *TemplateUpdater {
	return &TemplateUpdater{
		backupDir: backupDir,
	}
}

// UpdateTemplate updates a single template file with new version information
func (tu *TemplateUpdater) UpdateTemplate(templatePath string, versions map[string]*models.VersionInfo) error {
	// Check if template path exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template path does not exist: %s", templatePath)
	}

	// Walk through all files in the template directory
	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process template files and configuration files
		if tu.shouldProcessFile(path) {
			return tu.updateFile(path, versions)
		}

		return nil
	})
}

// UpdateAllTemplates updates all templates with new version information
func (tu *TemplateUpdater) UpdateAllTemplates(versions map[string]*models.VersionInfo) error {
	templateDirs := []string{
		"templates/frontend/nextjs-app",
		"templates/frontend/nextjs-home",
		"templates/frontend/nextjs-admin",
		"templates/backend/go-gin",
		"templates/mobile/android-kotlin",
		"templates/mobile/ios-swift",
	}

	for _, templateDir := range templateDirs {
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			fmt.Printf("Warning: Template directory not found: %s\n", templateDir)
			continue
		}

		if err := tu.UpdateTemplate(templateDir, versions); err != nil {
			return fmt.Errorf("failed to update template %s: %w", templateDir, err)
		}
	}

	return nil
}

// ValidateTemplate validates that a template can be updated
func (tu *TemplateUpdater) ValidateTemplate(templatePath string) error {
	// Check if template path exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template path does not exist: %s", templatePath)
	}

	// Check if it's a directory
	info, err := os.Stat(templatePath)
	if err != nil {
		return fmt.Errorf("failed to get template info: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("template path must be a directory: %s", templatePath)
	}

	// Check for required template files
	requiredFiles := []string{"package.json.tmpl", "README.md.tmpl"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(templatePath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Not all templates have all files, so just warn
			fmt.Printf("Warning: Template file not found: %s\n", filePath)
		}
	}

	return nil
}

// GetAffectedTemplates returns templates that would be affected by version changes
func (tu *TemplateUpdater) GetAffectedTemplates(versions map[string]*models.VersionInfo) ([]string, error) {
	var affectedTemplates []string

	templateDirs := []string{
		"templates/frontend/nextjs-app",
		"templates/frontend/nextjs-home",
		"templates/frontend/nextjs-admin",
		"templates/backend/go-gin",
		"templates/mobile/android-kotlin",
		"templates/mobile/ios-swift",
	}

	for _, templateDir := range templateDirs {
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			continue
		}

		// Check if this template would be affected by any version changes
		affected, err := tu.isTemplateAffected(templateDir, versions)
		if err != nil {
			return nil, fmt.Errorf("failed to check if template %s is affected: %w", templateDir, err)
		}

		if affected {
			affectedTemplates = append(affectedTemplates, templateDir)
		}
	}

	return affectedTemplates, nil
}

// BackupTemplates creates backups of templates before updating
func (tu *TemplateUpdater) BackupTemplates(templatePaths []string) error {
	// Ensure backup directory exists
	if err := os.MkdirAll(tu.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	for _, templatePath := range templatePaths {
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			fmt.Printf("Warning: Template path not found for backup: %s\n", templatePath)
			continue
		}

		// Create backup path
		templateName := filepath.Base(templatePath)
		backupPath := filepath.Join(tu.backupDir, fmt.Sprintf("%s_backup_%s", templateName, timestamp))

		// Copy template directory to backup location
		if err := tu.copyDirectory(templatePath, backupPath); err != nil {
			return fmt.Errorf("failed to backup template %s: %w", templatePath, err)
		}

		fmt.Printf("Created backup: %s -> %s\n", templatePath, backupPath)
	}

	return nil
}

// RestoreTemplates restores templates from backup
func (tu *TemplateUpdater) RestoreTemplates(templatePaths []string) error {
	// This would implement restoration from the most recent backup
	// For now, just log what would be restored
	for _, templatePath := range templatePaths {
		fmt.Printf("Would restore template: %s\n", templatePath)
	}

	fmt.Println("Note: Template restoration will be fully implemented in a future version")
	return nil
}

// Helper methods

func (tu *TemplateUpdater) shouldProcessFile(filePath string) bool {
	// Process template files and configuration files
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)

	// Template files
	if strings.HasSuffix(filePath, ".tmpl") {
		return true
	}

	// Configuration files that might contain version references
	configFiles := []string{
		"package.json",
		"go.mod",
		"build.gradle",
		"pom.xml",
		"Dockerfile",
		"docker-compose.yml",
		"versions.md",
	}

	for _, configFile := range configFiles {
		if base == configFile || strings.HasPrefix(base, configFile) {
			return true
		}
	}

	// Files with specific extensions
	processableExts := []string{".json", ".yaml", ".yml", ".toml", ".md", ".txt"}
	for _, processableExt := range processableExts {
		if ext == processableExt {
			return true
		}
	}

	return false
}

func (tu *TemplateUpdater) updateFile(filePath string, versions map[string]*models.VersionInfo) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	originalContent := string(content)
	updatedContent := originalContent

	// Apply version updates based on file type
	if strings.Contains(filePath, "package.json") {
		updatedContent = tu.updatePackageJSON(updatedContent, versions)
	} else if strings.Contains(filePath, "go.mod") {
		updatedContent = tu.updateGoMod(updatedContent, versions)
	} else if strings.Contains(filePath, "Dockerfile") {
		updatedContent = tu.updateDockerfile(updatedContent, versions)
	} else {
		// Generic version replacement
		updatedContent = tu.updateGenericVersions(updatedContent, versions)
	}

	// Only write if content changed
	if updatedContent != originalContent {
		if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
			return fmt.Errorf("failed to write updated file %s: %w", filePath, err)
		}
		fmt.Printf("Updated: %s\n", filePath)
	}

	return nil
}

func (tu *TemplateUpdater) updatePackageJSON(content string, versions map[string]*models.VersionInfo) string {
	// Update Node.js version references
	if nodeInfo, exists := versions["nodejs"]; exists {
		// Update engines.node field
		nodeVersionRegex := regexp.MustCompile(`"node":\s*"[^"]*"`)
		content = nodeVersionRegex.ReplaceAllString(content, fmt.Sprintf(`"node": ">=%s"`, nodeInfo.LatestVersion))
	}

	// Update framework versions
	if nextjsInfo, exists := versions["nextjs"]; exists {
		nextjsRegex := regexp.MustCompile(`"next":\s*"[^"]*"`)
		content = nextjsRegex.ReplaceAllString(content, fmt.Sprintf(`"next": "^%s"`, nextjsInfo.LatestVersion))
	}

	if reactInfo, exists := versions["react"]; exists {
		reactRegex := regexp.MustCompile(`"react":\s*"[^"]*"`)
		content = reactRegex.ReplaceAllString(content, fmt.Sprintf(`"react": "^%s"`, reactInfo.LatestVersion))

		reactDomRegex := regexp.MustCompile(`"react-dom":\s*"[^"]*"`)
		content = reactDomRegex.ReplaceAllString(content, fmt.Sprintf(`"react-dom": "^%s"`, reactInfo.LatestVersion))
	}

	// Update common packages
	packageUpdates := map[string]string{
		"typescript":   "typescript",
		"tailwindcss":  "tailwindcss",
		"eslint":       "eslint",
		"prettier":     "prettier",
		"jest":         "jest",
		"@types/node":  "@types/node",
		"@types/react": "@types/react",
	}

	for packageName, versionKey := range packageUpdates {
		if packageInfo, exists := versions[versionKey]; exists {
			packageRegex := regexp.MustCompile(fmt.Sprintf(`"%s":\s*"[^"]*"`, regexp.QuoteMeta(packageName)))
			content = packageRegex.ReplaceAllString(content, fmt.Sprintf(`"%s": "^%s"`, packageName, packageInfo.LatestVersion))
		}
	}

	return content
}

func (tu *TemplateUpdater) updateGoMod(content string, versions map[string]*models.VersionInfo) string {
	// Update Go version
	if goInfo, exists := versions["go"]; exists {
		goVersionRegex := regexp.MustCompile(`go\s+\d+\.\d+(\.\d+)?`)
		content = goVersionRegex.ReplaceAllString(content, fmt.Sprintf("go %s", goInfo.LatestVersion))
	}

	return content
}

func (tu *TemplateUpdater) updateDockerfile(content string, versions map[string]*models.VersionInfo) string {
	// Update Node.js base image
	if nodeInfo, exists := versions["nodejs"]; exists {
		nodeImageRegex := regexp.MustCompile(`FROM\s+node:\d+[^\s]*`)
		content = nodeImageRegex.ReplaceAllString(content, fmt.Sprintf("FROM node:%s", nodeInfo.LatestVersion))
	}

	// Update Go base image
	if goInfo, exists := versions["go"]; exists {
		goImageRegex := regexp.MustCompile(`FROM\s+golang:\d+[^\s]*`)
		content = goImageRegex.ReplaceAllString(content, fmt.Sprintf("FROM golang:%s", goInfo.LatestVersion))
	}

	return content
}

func (tu *TemplateUpdater) updateGenericVersions(content string, versions map[string]*models.VersionInfo) string {
	// Generic version replacement for template variables
	for name, versionInfo := range versions {
		// Replace template variables like {{.NodeVersion}}, {{.GoVersion}}, etc.
		templateVarRegex := regexp.MustCompile(fmt.Sprintf(`\{\{\s*\.%sVersion\s*\}\}`, strings.Title(name)))
		content = templateVarRegex.ReplaceAllString(content, versionInfo.LatestVersion)

		// Replace version references in comments and documentation
		versionRefRegex := regexp.MustCompile(fmt.Sprintf(`%s\s+v?\d+\.\d+(\.\d+)?`, name))
		content = versionRefRegex.ReplaceAllString(content, fmt.Sprintf("%s %s", name, versionInfo.LatestVersion))
	}

	return content
}

func (tu *TemplateUpdater) isTemplateAffected(templatePath string, versions map[string]*models.VersionInfo) (bool, error) {
	// Check if any files in the template would be affected by version changes
	affected := false

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if tu.shouldProcessFile(path) {
			// Read file and check if it contains version references
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for name := range versions {
				// Check for various version reference patterns
				patterns := []string{
					fmt.Sprintf(`"%s":\s*"[^"]*"`, name),
					fmt.Sprintf(`%s:\s*\d+\.\d+`, name),
					fmt.Sprintf(`\{\{\s*\.%sVersion\s*\}\}`, strings.Title(name)),
					fmt.Sprintf(`FROM\s+%s:\d+`, name),
				}

				for _, pattern := range patterns {
					if matched, _ := regexp.Match(pattern, content); matched {
						affected = true
						return filepath.SkipDir // Found a match, no need to check more files
					}
				}
			}
		}

		return nil
	})

	return affected, err
}

func (tu *TemplateUpdater) copyDirectory(src, dst string) error {
	// Simple directory copy implementation
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
		srcFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, srcFile, info.Mode())
	})
}

// Ensure TemplateUpdater implements the interface
var _ interfaces.TemplateUpdater = (*TemplateUpdater)(nil)
