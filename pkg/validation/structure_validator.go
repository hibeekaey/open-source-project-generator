package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// Pre-compiled regular expressions for structure validation
var (
	camelCaseRegex        = regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)+$`)
	readmePatternRegex    = regexp.MustCompile(`^README\.(md|txt|rst)$`)
	licensePatternRegex   = regexp.MustCompile(`^LICENSE(\.(txt|md))?$`)
	gitignorePatternRegex = regexp.MustCompile(`^\.gitignore$`)
	camelToKebabRegex     = regexp.MustCompile(`([a-z])([A-Z])`)
)

// StructureValidator provides specialized project structure validation
type StructureValidator struct {
	rules map[string]StructureRule
}

// StructureRule defines a rule for validating project structure
type StructureRule struct {
	Name        string
	Description string
	Required    bool
	Pattern     *regexp.Regexp
	Validator   func(path string) error
}

// NewStructureValidator creates a new structure validator
func NewStructureValidator() *StructureValidator {
	validator := &StructureValidator{
		rules: make(map[string]StructureRule),
	}
	validator.initializeDefaultRules()
	return validator
}

// ValidateProjectStructure validates the complete project structure
func (sv *StructureValidator) ValidateProjectStructure(projectPath string) (*interfaces.StructureValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(projectPath); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	result := &interfaces.StructureValidationResult{
		Valid:            true,
		RequiredFiles:    []interfaces.FileValidationResult{},
		RequiredDirs:     []interfaces.DirValidationResult{},
		NamingIssues:     []interfaces.NamingValidationIssue{},
		PermissionIssues: []interfaces.PermissionIssue{},
		Summary: interfaces.StructureValidationSummary{
			TotalFiles:       0,
			ValidFiles:       0,
			TotalDirectories: 0,
			ValidDirectories: 0,
			NamingIssues:     0,
			PermissionIssues: 0,
		},
	}

	// Validate required files
	if err := sv.validateRequiredFiles(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Unable to validate required project files.",
			"Some essential files may be missing or inaccessible")
	}

	// Validate directory structure
	if err := sv.validateDirectoryStructure(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Unable to validate project structure.",
			"Check if your project follows the expected directory layout")
	}

	// Validate naming conventions
	if err := sv.validateNamingConventions(projectPath, result); err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Unable to validate naming conventions.",
			"Some files or directories may not follow standard naming patterns")
	}

	// Validate file permissions
	if err := sv.validateFilePermissions(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate file permissions: %w", err)
	}

	// Validate project type specific structure
	if err := sv.validateProjectTypeStructure(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate project type structure: %w", err)
	}

	return result, nil
}

// validateRequiredFiles validates that required files exist
func (sv *StructureValidator) validateRequiredFiles(projectPath string, result *interfaces.StructureValidationResult) error {
	requiredFiles := []string{
		"README.md",
		"LICENSE",
		".gitignore",
	}

	for _, fileName := range requiredFiles {
		filePath := filepath.Join(projectPath, fileName)
		fileResult := sv.validateFile(filePath, true)
		result.RequiredFiles = append(result.RequiredFiles, fileResult)
		result.Summary.TotalFiles++

		if fileResult.Valid {
			result.Summary.ValidFiles++
		} else {
			result.Valid = false
		}
	}

	return nil
}

// validateDirectoryStructure validates the directory structure
func (sv *StructureValidator) validateDirectoryStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			result.Summary.TotalDirectories++

			// Validate directory naming
			if err := sv.validateDirectoryNaming(path, info, result); err != nil {
				return err
			}

			result.Summary.ValidDirectories++
		} else {
			result.Summary.TotalFiles++

			// Validate file naming
			if err := sv.validateFileNaming(path, info, result); err != nil {
				return err
			}

			result.Summary.ValidFiles++
		}

		return nil
	})

	return err
}

// validateNamingConventions validates naming conventions
func (sv *StructureValidator) validateNamingConventions(projectPath string, result *interfaces.StructureValidationResult) error {
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(info.Name(), ".") && info.Name() != ".gitignore" {
			return nil
		}

		// Check for spaces in names
		if strings.Contains(info.Name(), " ") {
			issue := interfaces.NamingValidationIssue{
				Path:       path,
				Type:       "file",
				Current:    info.Name(),
				Expected:   strings.ReplaceAll(info.Name(), " ", "_"),
				Convention: "no_spaces",
				Severity:   interfaces.ValidationSeverityWarning,
				Fixable:    true,
			}
			if info.IsDir() {
				issue.Type = "directory"
			}
			result.NamingIssues = append(result.NamingIssues, issue)
			result.Summary.NamingIssues++
		}

		// Check for uppercase in file names (except specific files)
		allowedUppercase := []string{"README.md", "LICENSE", "CHANGELOG.md", "CONTRIBUTING.md", "Dockerfile", "Makefile"}
		isAllowed := false
		for _, allowed := range allowedUppercase {
			if info.Name() == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed && !info.IsDir() && strings.ToLower(info.Name()) != info.Name() {
			issue := interfaces.NamingValidationIssue{
				Path:       path,
				Type:       "file",
				Current:    info.Name(),
				Expected:   strings.ToLower(info.Name()),
				Convention: "lowercase",
				Severity:   interfaces.ValidationSeverityWarning,
				Fixable:    true,
			}
			result.NamingIssues = append(result.NamingIssues, issue)
			result.Summary.NamingIssues++
		}

		return nil
	})

	return err
}

// validateFilePermissions validates file permissions
func (sv *StructureValidator) validateFilePermissions(projectPath string, result *interfaces.StructureValidationResult) error {
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		mode := info.Mode()

		// Check for overly permissive permissions
		if !info.IsDir() && mode&0o077 != 0 {
			issue := interfaces.PermissionIssue{
				Path:     path,
				Current:  mode.String(),
				Expected: "rw-r--r--",
				Type:     "file",
				Security: true,
				Severity: interfaces.ValidationSeverityWarning,
				Fixable:  true,
			}
			result.PermissionIssues = append(result.PermissionIssues, issue)
			result.Summary.PermissionIssues++
		}

		// Check for executable files that shouldn't be executable
		if !info.IsDir() && mode&0o111 != 0 {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			nonExecutableExts := []string{".md", ".txt", ".json", ".yaml", ".yml", ".xml", ".html", ".css", ".js", ".ts"}

			for _, nonExecExt := range nonExecutableExts {
				if ext == nonExecExt {
					issue := interfaces.PermissionIssue{
						Path:     path,
						Current:  mode.String(),
						Expected: "rw-r--r--",
						Type:     "file",
						Security: false,
						Severity: interfaces.ValidationSeverityInfo,
						Fixable:  true,
					}
					result.PermissionIssues = append(result.PermissionIssues, issue)
					result.Summary.PermissionIssues++
					break
				}
			}
		}

		return nil
	})

	return err
}

// validateProjectTypeStructure validates structure based on project type
func (sv *StructureValidator) validateProjectTypeStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Detect project type and validate accordingly
	projectType := sv.detectProjectType(projectPath)

	switch projectType {
	case "go":
		return sv.validateGoProjectStructure(projectPath, result)
	case "node":
		return sv.validateNodeProjectStructure(projectPath, result)
	case "python":
		return sv.validatePythonProjectStructure(projectPath, result)
	case "docker":
		return sv.validateDockerProjectStructure(projectPath, result)
	}

	return nil
}

// detectProjectType detects the type of project based on files present
func (sv *StructureValidator) detectProjectType(projectPath string) string {
	// Check for Go project
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return "go"
	}

	// Check for Node.js project
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return "node"
	}

	// Check for Python project
	pythonFiles := []string{"setup.py", "pyproject.toml", "requirements.txt"}
	for _, file := range pythonFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return "python"
		}
	}

	// Check for Docker project
	if _, err := os.Stat(filepath.Join(projectPath, "Dockerfile")); err == nil {
		return "docker"
	}

	return "unknown"
}

// validateGoProjectStructure validates Go project structure
func (sv *StructureValidator) validateGoProjectStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Check for main.go or cmd directory
	mainGoPath := filepath.Join(projectPath, "main.go")
	cmdDirPath := filepath.Join(projectPath, "cmd")

	hasMainGo := false
	hasCmdDir := false

	if _, err := os.Stat(mainGoPath); err == nil {
		hasMainGo = true
	}

	if info, err := os.Stat(cmdDirPath); err == nil && info.IsDir() {
		hasCmdDir = true
	}

	if !hasMainGo && !hasCmdDir {
		fileResult := interfaces.FileValidationResult{
			Path:     mainGoPath,
			Required: true,
			Exists:   false,
			Valid:    false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: interfaces.ValidationSeverityWarning,
					Message:  "Go project should have either main.go or cmd/ directory",
					File:     projectPath,
					Rule:     "go.entry_point",
					Fixable:  false,
				},
			},
		}
		result.RequiredFiles = append(result.RequiredFiles, fileResult)
	}

	// Check for recommended Go directories
	recommendedDirs := []string{"pkg", "internal"}
	for _, dir := range recommendedDirs {
		dirPath := filepath.Join(projectPath, dir)
		if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
			dirResult := interfaces.DirValidationResult{
				Path:     dirPath,
				Required: false,
				Exists:   false,
				Valid:    true, // Not required, so still valid
				Issues: []interfaces.ValidationIssue{
					{
						Type:     "info",
						Severity: interfaces.ValidationSeverityInfo,
						Message:  fmt.Sprintf("Consider adding %s/ directory for better organization", dir),
						File:     dirPath,
						Rule:     "go.recommended_structure",
						Fixable:  true,
					},
				},
			}
			result.RequiredDirs = append(result.RequiredDirs, dirResult)
		}
	}

	return nil
}

// validateNodeProjectStructure validates Node.js project structure
func (sv *StructureValidator) validateNodeProjectStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Check for src directory or index.js
	srcDirPath := filepath.Join(projectPath, "src")
	indexJsPath := filepath.Join(projectPath, "index.js")

	hasSrcDir := false
	hasIndexJs := false

	if info, err := os.Stat(srcDirPath); err == nil && info.IsDir() {
		hasSrcDir = true
	}

	if _, err := os.Stat(indexJsPath); err == nil {
		hasIndexJs = true
	}

	if !hasSrcDir && !hasIndexJs {
		fileResult := interfaces.FileValidationResult{
			Path:     indexJsPath,
			Required: false,
			Exists:   false,
			Valid:    true,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "info",
					Severity: interfaces.ValidationSeverityInfo,
					Message:  "Consider having either src/ directory or index.js as entry point",
					File:     projectPath,
					Rule:     "node.entry_point",
					Fixable:  true,
				},
			},
		}
		result.RequiredFiles = append(result.RequiredFiles, fileResult)
	}

	return nil
}

// validatePythonProjectStructure validates Python project structure
func (sv *StructureValidator) validatePythonProjectStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Check for src directory or main module
	srcDirPath := filepath.Join(projectPath, "src")

	if info, err := os.Stat(srcDirPath); err != nil || !info.IsDir() {
		dirResult := interfaces.DirValidationResult{
			Path:     srcDirPath,
			Required: false,
			Exists:   false,
			Valid:    true,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "info",
					Severity: interfaces.ValidationSeverityInfo,
					Message:  "Consider using src/ directory for Python packages",
					File:     srcDirPath,
					Rule:     "python.src_layout",
					Fixable:  true,
				},
			},
		}
		result.RequiredDirs = append(result.RequiredDirs, dirResult)
	}

	return nil
}

// validateDockerProjectStructure validates Docker project structure
func (sv *StructureValidator) validateDockerProjectStructure(projectPath string, result *interfaces.StructureValidationResult) error {
	// Check for .dockerignore
	dockerignorePath := filepath.Join(projectPath, ".dockerignore")
	fileResult := sv.validateFile(dockerignorePath, false)

	if !fileResult.Exists {
		fileResult.Issues = append(fileResult.Issues, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: interfaces.ValidationSeverityWarning,
			Message:  "Consider adding .dockerignore file for better Docker builds",
			File:     dockerignorePath,
			Rule:     "docker.dockerignore",
			Fixable:  true,
		})
	}

	result.RequiredFiles = append(result.RequiredFiles, fileResult)

	return nil
}

// validateFile validates a single file
func (sv *StructureValidator) validateFile(filePath string, required bool) interfaces.FileValidationResult {
	result := interfaces.FileValidationResult{
		Path:     filePath,
		Required: required,
		Exists:   false,
		Valid:    false,
		Issues:   []interfaces.ValidationIssue{},
		Size:     0,
		Mode:     "",
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if required {
				result.Issues = append(result.Issues, interfaces.ValidationIssue{
					Type:     "error",
					Severity: interfaces.ValidationSeverityError,
					Message:  "Required file does not exist",
					File:     filePath,
					Rule:     "structure.required_file",
					Fixable:  true,
				})
			}
			return result
		}
		result.Issues = append(result.Issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: interfaces.ValidationSeverityError,
			Message:  fmt.Sprintf("Failed to access file: %v", err),
			File:     filePath,
			Rule:     "structure.file_access",
			Fixable:  false,
		})
		return result
	}

	result.Exists = true
	result.Size = info.Size()
	result.Mode = info.Mode().String()
	result.Valid = len(result.Issues) == 0

	return result
}

// validateDirectoryNaming validates directory naming conventions
func (sv *StructureValidator) validateDirectoryNaming(path string, info os.FileInfo, result *interfaces.StructureValidationResult) error {
	name := info.Name()

	// Skip root directory and hidden directories
	if name == "." || strings.HasPrefix(name, ".") {
		return nil
	}

	// Check for camelCase in directory names (should be kebab-case or snake_case)
	if camelCaseRegex.MatchString(name) {
		issue := interfaces.NamingValidationIssue{
			Path:       path,
			Type:       "directory",
			Current:    name,
			Expected:   toKebabCase(name),
			Convention: "kebab_case",
			Severity:   interfaces.ValidationSeverityWarning,
			Fixable:    true,
		}
		result.NamingIssues = append(result.NamingIssues, issue)
		result.Summary.NamingIssues++
	}

	return nil
}

// validateFileNaming validates file naming conventions
func (sv *StructureValidator) validateFileNaming(path string, info os.FileInfo, result *interfaces.StructureValidationResult) error {
	name := info.Name()

	// Skip hidden files
	if strings.HasPrefix(name, ".") {
		return nil
	}

	// Get file extension
	ext := filepath.Ext(name)
	baseName := strings.TrimSuffix(name, ext)

	// Check for camelCase in file names (context-dependent)
	if camelCaseRegex.MatchString(baseName) {
		// Allow camelCase for certain file types (JavaScript, TypeScript)
		allowedExts := []string{".js", ".ts", ".jsx", ".tsx"}
		isAllowed := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			issue := interfaces.NamingValidationIssue{
				Path:       path,
				Type:       "file",
				Current:    name,
				Expected:   toSnakeCase(baseName) + ext,
				Convention: "snake_case",
				Severity:   interfaces.ValidationSeverityWarning,
				Fixable:    true,
			}
			result.NamingIssues = append(result.NamingIssues, issue)
			result.Summary.NamingIssues++
		}
	}

	return nil
}

// initializeDefaultRules initializes default structure validation rules
func (sv *StructureValidator) initializeDefaultRules() {
	sv.rules["readme_required"] = StructureRule{
		Name:        "README Required",
		Description: "Project must have a README file",
		Required:    true,
		Pattern:     readmePatternRegex,
	}

	sv.rules["license_required"] = StructureRule{
		Name:        "License Required",
		Description: "Project must have a LICENSE file",
		Required:    true,
		Pattern:     licensePatternRegex,
	}

	sv.rules["gitignore_recommended"] = StructureRule{
		Name:        "Gitignore Recommended",
		Description: "Project should have a .gitignore file",
		Required:    false,
		Pattern:     gitignorePatternRegex,
	}
}

// Helper functions for naming convention conversion
func toKebabCase(s string) string {
	// Convert camelCase to kebab-case using pre-compiled regex
	return strings.ToLower(camelToKebabRegex.ReplaceAllString(s, "${1}-${2}"))
}

func toSnakeCase(s string) string {
	// Convert camelCase to snake_case using pre-compiled regex
	return strings.ToLower(camelToKebabRegex.ReplaceAllString(s, "${1}_${2}"))
}
