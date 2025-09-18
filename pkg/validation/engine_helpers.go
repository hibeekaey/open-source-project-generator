package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// validateProjectStructureBasic performs basic project structure validation
func (e *Engine) validateProjectStructureBasic(path string, result *models.ValidationResult) error {
	// Check for README file
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "README"}
	readmeFound := false
	for _, readme := range readmeFiles {
		if _, err := os.Stat(filepath.Join(path, readme)); err == nil {
			readmeFound = true
			break
		}
	}
	if !readmeFound {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "README file is missing",
			File:     path,
			Rule:     "structure.readme.required",
			Fixable:  true,
		})
	}

	// Check for LICENSE file
	licenseFiles := []string{"LICENSE", "LICENSE.txt", "LICENSE.md", "COPYING"}
	licenseFound := false
	for _, license := range licenseFiles {
		if _, err := os.Stat(filepath.Join(path, license)); err == nil {
			licenseFound = true
			break
		}
	}
	if !licenseFound {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "LICENSE file is missing",
			File:     path,
			Rule:     "structure.license.required",
			Fixable:  true,
		})
	}

	return nil
}

// validateProjectDependenciesBasic performs basic dependency validation
func (e *Engine) validateProjectDependenciesBasic(path string, result *models.ValidationResult) error {
	// Validate package.json if it exists
	packageJsonPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		if err := e.ValidatePackageJSON(packageJsonPath); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid package.json: %v", err),
				File:     packageJsonPath,
				Rule:     "dependencies.package_json.valid",
				Fixable:  false,
			})
		}
	}

	// Validate go.mod if it exists
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		if err := e.ValidateGoMod(goModPath); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid go.mod: %v", err),
				File:     goModPath,
				Rule:     "dependencies.go_mod.valid",
				Fixable:  false,
			})
		}
	}

	return nil
}

// validateProjectConfigurationFiles validates configuration files
func (e *Engine) validateProjectConfigurationFiles(path string, result *models.ValidationResult) error {
	// Find and validate YAML files
	yamlFiles := []string{}
	jsonFiles := []string{}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		switch ext {
		case ".yaml", ".yml":
			yamlFiles = append(yamlFiles, filePath)
		case ".json":
			jsonFiles = append(jsonFiles, filePath)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	// Validate YAML files
	for _, yamlFile := range yamlFiles {
		if err := e.ValidateYAML(yamlFile); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid YAML file: %v", err),
				File:     yamlFile,
				Rule:     "configuration.yaml.syntax",
				Fixable:  false,
			})
		}
	}

	// Validate JSON files
	for _, jsonFile := range jsonFiles {
		if err := e.ValidateJSON(jsonFile); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid JSON file: %v", err),
				File:     jsonFile,
				Rule:     "configuration.json.syntax",
				Fixable:  false,
			})
		}
	}

	return nil
}

// validateFile validates a single file
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateFile(filePath string, required bool) interfaces.FileValidationResult {
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
			Message:  fmt.Sprintf("Failed to stat file: %v", err),
			File:     filePath,
			Rule:     "structure.file_access",
			Fixable:  false,
		})
		return result
	}

	result.Exists = true
	result.Size = info.Size()
	result.Mode = info.Mode().String()

	// File exists, so it's valid unless there are other issues
	result.Valid = len(result.Issues) == 0

	return result
}

// validateNamingConventions validates naming conventions
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateNamingConventions(path string, result *interfaces.StructureValidationResult) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Check for spaces in names (generally not recommended)
		if strings.Contains(info.Name(), " ") {
			issue := interfaces.NamingValidationIssue{
				Path:       filePath,
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

		return nil
	})

	return err
}

// validateFilePermissions validates file permissions
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateFilePermissions(path string, result *interfaces.StructureValidationResult) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check for overly permissive permissions
		mode := info.Mode()
		if mode&0o077 != 0 && !info.IsDir() {
			// File is readable/writable by group or others
			issue := interfaces.PermissionIssue{
				Path:     filePath,
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

		return nil
	})

	return err
}

// validatePackageJSONDependencies validates package.json dependencies
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validatePackageJSONDependencies(path string, result *interfaces.DependencyValidationResult) error {
	packageJsonPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return nil // No package.json to validate
	}

	content, err := utils.SafeReadFile(packageJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("invalid JSON in package.json: %w", err)
	}

	// Validate dependencies
	if deps, exists := pkg["dependencies"]; exists {
		if depsMap, ok := deps.(map[string]interface{}); ok {
			for name, version := range depsMap {
				versionStr, ok := version.(string)
				if !ok {
					continue
				}

				dep := interfaces.DependencyValidation{
					Name:          name,
					Version:       versionStr,
					Type:          "direct",
					Valid:         true,
					Available:     true, // Assume available for now
					LatestVersion: versionStr,
				}

				result.Dependencies = append(result.Dependencies, dep)
				result.Summary.TotalDependencies++
				result.Summary.ValidDependencies++
			}
		}
	}

	return nil
}

// validateGoModDependencies validates go.mod dependencies
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateGoModDependencies(path string, result *interfaces.DependencyValidationResult) error {
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil // No go.mod to validate
	}

	content, err := utils.SafeReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		if inRequireBlock || strings.HasPrefix(line, "require ") {
			// Parse dependency line
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				if strings.HasPrefix(name, "require") && len(parts) >= 3 {
					name = parts[1]
				}

				version := parts[len(parts)-1]

				dep := interfaces.DependencyValidation{
					Name:          name,
					Version:       version,
					Type:          "direct",
					Valid:         true,
					Available:     true, // Assume available for now
					LatestVersion: version,
				}

				result.Dependencies = append(result.Dependencies, dep)
				result.Summary.TotalDependencies++
				result.Summary.ValidDependencies++
			}
		}
	}

	return nil
}

// checkDependencyConflicts checks for dependency conflicts
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) checkDependencyConflicts(result *interfaces.DependencyValidationResult) error {
	// Simple conflict detection - check for duplicate dependencies with different versions
	depVersions := make(map[string][]string)

	for _, dep := range result.Dependencies {
		depVersions[dep.Name] = append(depVersions[dep.Name], dep.Version)
	}

	for name, versions := range depVersions {
		if len(versions) > 1 {
			// Check if all versions are the same
			firstVersion := versions[0]
			hasConflict := false
			for _, version := range versions[1:] {
				if version != firstVersion {
					hasConflict = true
					break
				}
			}

			if hasConflict {
				conflict := interfaces.DependencyConflict{
					Dependency1: name,
					Version1:    versions[0],
					Dependency2: name,
					Version2:    versions[1],
					Reason:      "Multiple versions of the same dependency",
					Severity:    interfaces.ValidationSeverityWarning,
				}
				result.Conflicts = append(result.Conflicts, conflict)
				result.Summary.ConflictCount++
			}
		}
	}

	return nil
}

// scanForSecrets scans for potential secrets in files
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) scanForSecrets(path string, result *interfaces.SecurityValidationResult) error {
	secretPatterns := map[string]*regexp.Regexp{
		"api_key":     regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?([a-zA-Z0-9]{20,})['"]?`),
		"password":    regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*['"]?([^'"\s]{8,})['"]?`),
		"token":       regexp.MustCompile(`(?i)(token|auth[_-]?token)\s*[:=]\s*['"]?([a-zA-Z0-9]{20,})['"]?`),
		"private_key": regexp.MustCompile(`-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`),
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip binary files and certain extensions
		ext := strings.ToLower(filepath.Ext(info.Name()))
		skipExtensions := []string{".exe", ".bin", ".jpg", ".png", ".gif", ".pdf", ".zip", ".tar", ".gz"}
		for _, skipExt := range skipExtensions {
			if ext == skipExt {
				return nil
			}
		}

		// Read file content
		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil // Skip files that can't be read
		}

		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")

		for lineNum, line := range lines {
			for secretType, pattern := range secretPatterns {
				if matches := pattern.FindStringSubmatch(line); matches != nil {
					secret := interfaces.SecretDetection{
						Type:       secretType,
						File:       filePath,
						Line:       lineNum + 1,
						Column:     strings.Index(line, matches[0]) + 1,
						Pattern:    pattern.String(),
						Confidence: 0.8,
						Masked:     maskSecret(matches[0]),
					}
					result.Secrets = append(result.Secrets, secret)
					result.Summary.SecretsFound++
				}
			}
		}

		return nil
	})

	return err
}

// validateSecurityConfigurations validates security configurations
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateSecurityConfigurations(path string, result *interfaces.SecurityValidationResult) error {
	// Check for common security configuration files
	securityFiles := map[string]func(string) error{
		".env":               e.validateEnvFile,
		"docker-compose.yml": e.validateDockerComposeFile,
		"Dockerfile":         e.validateDockerfileSecurityFile,
	}

	for filename, validator := range securityFiles {
		filePath := filepath.Join(path, filename)
		if _, err := os.Stat(filePath); err == nil {
			if err := validator(filePath); err != nil {
				issue := interfaces.SecurityIssue{
					Type:        "configuration",
					Severity:    interfaces.ValidationSeverityWarning,
					Title:       fmt.Sprintf("Security issue in %s", filename),
					Description: err.Error(),
					File:        filePath,
					Rule:        "security.configuration",
					Fixable:     false,
				}
				result.SecurityIssues = append(result.SecurityIssues, issue)
				result.Summary.TotalIssues++
				result.Summary.MediumSeverity++
			}
		}
	}

	return nil
}

// validateSecurityPermissions validates security-related file permissions
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateSecurityPermissions(path string, result *interfaces.SecurityValidationResult) error {
	sensitiveFiles := []string{".env", "config.json", "secrets.yaml", "private.key"}

	for _, filename := range sensitiveFiles {
		filePath := filepath.Join(path, filename)
		if info, err := os.Stat(filePath); err == nil {
			mode := info.Mode()
			if mode&0o044 != 0 {
				// File is readable by group or others
				issue := interfaces.PermissionIssue{
					Path:     filePath,
					Current:  mode.String(),
					Expected: "rw-------",
					Type:     "file",
					Security: true,
					Severity: interfaces.ValidationSeverityError,
					Fixable:  true,
				}
				result.Permissions = append(result.Permissions, issue)
				result.Summary.TotalIssues++
				result.Summary.HighSeverity++
			}
		}
	}

	return nil
}

// analyzeCodeSmells analyzes code for quality issues
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) analyzeCodeSmells(path string, result *interfaces.QualityValidationResult) error {
	// Simple code smell detection
	codeExtensions := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c"}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		isCodeFile := false
		for _, codeExt := range codeExtensions {
			if ext == codeExt {
				isCodeFile = true
				break
			}
		}

		if !isCodeFile {
			return nil
		}

		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			// Check for long lines
			if len(line) > 120 {
				smell := interfaces.CodeSmell{
					Type:        "long_line",
					File:        filePath,
					Line:        lineNum + 1,
					Description: fmt.Sprintf("Line too long (%d characters)", len(line)),
					Severity:    interfaces.ValidationSeverityWarning,
				}
				result.CodeSmells = append(result.CodeSmells, smell)
				result.Summary.CodeSmells++
			}

			// Check for TODO comments
			if strings.Contains(strings.ToUpper(line), "TODO") {
				smell := interfaces.CodeSmell{
					Type:        "todo_comment",
					File:        filePath,
					Line:        lineNum + 1,
					Description: "TODO comment found",
					Severity:    interfaces.ValidationSeverityInfo,
				}
				result.CodeSmells = append(result.CodeSmells, smell)
				result.Summary.CodeSmells++
			}
		}

		return nil
	})

	return err
}

// analyzeComplexity analyzes code complexity
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) analyzeComplexity(path string, result *interfaces.QualityValidationResult) error {
	// Simple complexity analysis - count nested blocks
	codeExtensions := []string{".go", ".js", ".ts", ".py", ".java"}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		isCodeFile := false
		for _, codeExt := range codeExtensions {
			if ext == codeExt {
				isCodeFile = true
				break
			}
		}

		if !isCodeFile {
			return nil
		}

		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		nestingLevel := 0
		maxNesting := 0
		currentFunction := ""

		for lineNum, line := range lines {
			trimmed := strings.TrimSpace(line)

			// Simple function detection
			if strings.Contains(trimmed, "func ") || strings.Contains(trimmed, "function ") {
				currentFunction = trimmed
				nestingLevel = 0
				maxNesting = 0
			}

			// Count opening braces
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			nestingLevel += openBraces - closeBraces

			if nestingLevel > maxNesting {
				maxNesting = nestingLevel
			}

			// If nesting is too deep, report complexity issue
			if nestingLevel > 4 {
				complexity := interfaces.ComplexityIssue{
					Type:       "cyclomatic",
					File:       filePath,
					Function:   currentFunction,
					Line:       lineNum + 1,
					Complexity: nestingLevel,
					Threshold:  4,
					Severity:   interfaces.ValidationSeverityWarning,
				}
				result.Complexity = append(result.Complexity, complexity)
				result.Summary.ComplexityIssues++
			}
		}

		return nil
	})

	return err
}

// detectDuplication detects code duplication
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) detectDuplication(path string, result *interfaces.QualityValidationResult) error {
	// Simple duplication detection - find identical lines
	fileContents := make(map[string][]string)
	codeExtensions := []string{".go", ".js", ".ts", ".py", ".java"}

	// Read all code files
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		isCodeFile := false
		for _, codeExt := range codeExtensions {
			if ext == codeExt {
				isCodeFile = true
				break
			}
		}

		if !isCodeFile {
			return nil
		}

		content, err := utils.SafeReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		fileContents[filePath] = lines

		return nil
	})

	if err != nil {
		return err
	}

	// Simple duplication detection - find blocks of 5+ identical consecutive lines
	const minDuplicationLines = 5

	for file1, lines1 := range fileContents {
		for file2, lines2 := range fileContents {
			if file1 >= file2 { // Avoid duplicate comparisons
				continue
			}

			// Compare lines between files
			for i := 0; i <= len(lines1)-minDuplicationLines; i++ {
				for j := 0; j <= len(lines2)-minDuplicationLines; j++ {
					matchCount := 0
					for k := 0; k < minDuplicationLines && i+k < len(lines1) && j+k < len(lines2); k++ {
						if strings.TrimSpace(lines1[i+k]) == strings.TrimSpace(lines2[j+k]) &&
							strings.TrimSpace(lines1[i+k]) != "" {
							matchCount++
						} else {
							break
						}
					}

					if matchCount >= minDuplicationLines {
						duplication := interfaces.Duplication{
							Files:      []string{file1, file2},
							Lines:      matchCount,
							Tokens:     matchCount * 10, // Rough estimate
							Percentage: 100.0,
						}
						result.Duplication = append(result.Duplication, duplication)
						result.Summary.DuplicationIssues++
					}
				}
			}
		}
	}

	return nil
}

// calculateQualityScore calculates overall quality score
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) calculateQualityScore(result *interfaces.QualityValidationResult) {
	totalIssues := result.Summary.CodeSmells + result.Summary.ComplexityIssues + result.Summary.DuplicationIssues
	result.Summary.TotalIssues = totalIssues

	// Simple scoring algorithm
	baseScore := 100.0
	penalty := float64(totalIssues) * 2.0

	result.Summary.QualityScore = baseScore - penalty
	if result.Summary.QualityScore < 0 {
		result.Summary.QualityScore = 0
	}

	// Determine maintainability grade
	if result.Summary.QualityScore >= 90 {
		result.Summary.Maintainability = "A"
	} else if result.Summary.QualityScore >= 80 {
		result.Summary.Maintainability = "B"
	} else if result.Summary.QualityScore >= 70 {
		result.Summary.Maintainability = "C"
	} else if result.Summary.QualityScore >= 60 {
		result.Summary.Maintainability = "D"
	} else {
		result.Summary.Maintainability = "F"
	}
}

// Helper functions for configuration validation
func (e *Engine) validateRequiredConfigFields(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	result.Summary.TotalProperties++

	if config.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "name",
			Value:    "",
			Type:     "required",
			Message:  "Project name is required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.name.required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	} else {
		result.Summary.ValidProperties++
	}

	result.Summary.TotalProperties++
	if config.Organization == "" {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "organization",
			Value:    "",
			Type:     "required",
			Message:  "Organization is required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.organization.required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	} else {
		result.Summary.ValidProperties++
	}

	// OutputPath is optional for basic config validation
	// It's only required when actually generating a project
	result.Summary.TotalProperties++
	if config.OutputPath != "" {
		result.Summary.ValidProperties++
	}
	// Note: OutputPath validation can be added separately for generation-time validation
}

// validateConfigFieldFormats validates configuration field formats
func (e *Engine) validateConfigFieldFormats(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	// Validate project name format
	if config.Name != "" {
		nameRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9]$`)
		if !nameRegex.MatchString(config.Name) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "name",
				Value:      config.Name,
				Type:       "format",
				Message:    "Project name should contain only alphanumeric characters, hyphens, and underscores",
				Suggestion: "Use a name like 'my-project' or 'my_project'",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.name.format",
			})
			result.Summary.WarningCount++
		}
	}

	// Validate output path format
	if config.OutputPath != "" {
		if strings.Contains(config.OutputPath, "..") {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    "output_path",
				Value:    config.OutputPath,
				Type:     "security",
				Message:  "Output path cannot contain '..' for security reasons",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "config.output_path.security",
			})
			result.Valid = false
			result.Summary.ErrorCount++
		}
	}
}

// validateComponentConfiguration validates component configuration
func (e *Engine) validateComponentConfiguration(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	// Basic component validation - check if any components are configured
	if !config.Components.Frontend.NextJS.App && !config.Components.Backend.GoGin {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "components",
			Value:      "",
			Type:       "missing",
			Message:    "No components are enabled in the configuration",
			Suggestion: "Enable at least one component (frontend or backend)",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "config.components.empty",
		})
		result.Summary.WarningCount++
	}
}

// validatePropertyValue validates a property value against its schema
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validatePropertyValue(key string, value interface{}, schema interfaces.PropertySchema) error {
	// Type validation
	switch schema.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
		strValue := value.(string)

		// Length validation
		if schema.MinLength != nil && len(strValue) < *schema.MinLength {
			return fmt.Errorf("string too short, minimum length is %d", *schema.MinLength)
		}
		if schema.MaxLength != nil && len(strValue) > *schema.MaxLength {
			return fmt.Errorf("string too long, maximum length is %d", *schema.MaxLength)
		}

		// Pattern validation
		if schema.Pattern != "" {
			regex, err := regexp.Compile(schema.Pattern)
			if err != nil {
				return fmt.Errorf("invalid pattern in schema: %w", err)
			}
			if !regex.MatchString(strValue) {
				return fmt.Errorf("string does not match pattern %s", schema.Pattern)
			}
		}

		// Enum validation
		if len(schema.Enum) > 0 {
			valid := false
			for _, enumValue := range schema.Enum {
				if strValue == enumValue {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("value must be one of: %v", schema.Enum)
			}
		}

	case "number":
		var numValue float64
		switch v := value.(type) {
		case float64:
			numValue = v
		case int:
			numValue = float64(v)
		default:
			return fmt.Errorf("expected number, got %T", value)
		}

		// Range validation
		if schema.Minimum != nil && numValue < *schema.Minimum {
			return fmt.Errorf("number too small, minimum is %f", *schema.Minimum)
		}
		if schema.Maximum != nil && numValue > *schema.Maximum {
			return fmt.Errorf("number too large, maximum is %f", *schema.Maximum)
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	}

	return nil
}

// validateVariableValidation validates template variable validation rules
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateVariableValidation(name string, validation *interfaces.VariableValidation) error {
	if validation.Pattern != "" {
		if _, err := regexp.Compile(validation.Pattern); err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	}

	if validation.MinLength != nil && *validation.MinLength < 0 {
		return fmt.Errorf("minimum length cannot be negative")
	}

	if validation.MaxLength != nil && *validation.MaxLength < 0 {
		return fmt.Errorf("maximum length cannot be negative")
	}

	if validation.MinLength != nil && validation.MaxLength != nil && *validation.MinLength > *validation.MaxLength {
		return fmt.Errorf("minimum length cannot be greater than maximum length")
	}

	return nil
}

// validateTemplateStructureInternal validates template structure internally
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateTemplateStructureInternal(path string, result *interfaces.TemplateValidationResult) error {
	// Check for template metadata file
	metadataFiles := []string{"template.yaml", "template.yml", "metadata.yaml", "metadata.yml"}
	metadataFound := false

	for _, metadataFile := range metadataFiles {
		metadataPath := filepath.Join(path, metadataFile)
		if _, err := os.Stat(metadataPath); err == nil {
			metadataFound = true
			break
		}
	}

	if !metadataFound {
		result.Valid = false
		result.Issues = append(result.Issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: interfaces.ValidationSeverityError,
			Message:  "Template metadata file (template.yaml) is missing",
			File:     path,
			Rule:     "template.metadata.required",
			Fixable:  true,
		})
		result.Summary.ErrorCount++
	}

	return nil
}

// validateTemplateMetadataFile validates template metadata file
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateTemplateMetadataFile(path string, result *interfaces.TemplateValidationResult) error {
	metadataFiles := []string{"template.yaml", "template.yml"}

	for _, metadataFile := range metadataFiles {
		metadataPath := filepath.Join(path, metadataFile)
		if _, err := os.Stat(metadataPath); err == nil {
			content, err := utils.SafeReadFile(metadataPath)
			if err != nil {
				result.Issues = append(result.Issues, interfaces.ValidationIssue{
					Type:     "error",
					Severity: interfaces.ValidationSeverityError,
					Message:  fmt.Sprintf("Failed to read metadata file: %v", err),
					File:     metadataPath,
					Rule:     "template.metadata.readable",
					Fixable:  false,
				})
				result.Summary.ErrorCount++
				continue
			}

			var metadata interfaces.TemplateMetadata
			if err := yaml.Unmarshal(content, &metadata); err != nil {
				result.Issues = append(result.Issues, interfaces.ValidationIssue{
					Type:     "error",
					Severity: interfaces.ValidationSeverityError,
					Message:  fmt.Sprintf("Invalid YAML in metadata file: %v", err),
					File:     metadataPath,
					Rule:     "template.metadata.syntax",
					Fixable:  false,
				})
				result.Summary.ErrorCount++
				continue
			}

			// Validate metadata content
			if err := e.ValidateTemplateMetadata(&metadata); err != nil {
				result.Issues = append(result.Issues, interfaces.ValidationIssue{
					Type:     "error",
					Severity: interfaces.ValidationSeverityError,
					Message:  fmt.Sprintf("Invalid metadata: %v", err),
					File:     metadataPath,
					Rule:     "template.metadata.content",
					Fixable:  false,
				})
				result.Summary.ErrorCount++
			}

			break
		}
	}

	return nil
}

// validateTemplateFiles validates template files
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateTemplateFiles(path string, result *interfaces.TemplateValidationResult) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		result.Summary.TotalFiles++

		// Check if template files have .tmpl extension
		if strings.Contains(filePath, "templates/") && !strings.HasSuffix(info.Name(), ".tmpl") {
			result.Warnings = append(result.Warnings, interfaces.ValidationIssue{
				Type:     "warning",
				Severity: interfaces.ValidationSeverityWarning,
				Message:  "Template file should have .tmpl extension",
				File:     filePath,
				Rule:     "template.file.extension",
				Fixable:  true,
			})
			result.Summary.WarningCount++
		} else {
			result.Summary.ValidFiles++
		}

		return nil
	})

	return err
}

// validateTemplateFileExtensions validates template file extensions
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateTemplateFileExtensions(path string, result *interfaces.StructureValidationResult) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if files in template directories have .tmpl extension
		if strings.Contains(filePath, "templates/") && !strings.HasSuffix(info.Name(), ".tmpl") {
			issue := interfaces.NamingValidationIssue{
				Path:       filePath,
				Type:       "file",
				Current:    info.Name(),
				Expected:   info.Name() + ".tmpl",
				Convention: "template_extension",
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

// Helper functions for security validation
//
//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateEnvFile(filePath string) error {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for unquoted values that might contain secrets
		if strings.Contains(line, "=") && !strings.Contains(line, "\"") && !strings.Contains(line, "'") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && len(parts[1]) > 20 {
				return fmt.Errorf("line %d: potentially sensitive value should be quoted", lineNum+1)
			}
		}
	}

	return nil
}

//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateDockerComposeFile(filePath string) error {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read docker-compose file: %w", err)
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(content, &compose); err != nil {
		return fmt.Errorf("invalid YAML in docker-compose file: %w", err)
	}

	// Check for security issues like privileged containers
	if services, ok := compose["services"].(map[string]interface{}); ok {
		for serviceName, service := range services {
			if serviceMap, ok := service.(map[string]interface{}); ok {
				if privileged, exists := serviceMap["privileged"]; exists {
					if privilegedBool, ok := privileged.(bool); ok && privilegedBool {
						return fmt.Errorf("service '%s' runs in privileged mode, which is a security risk", serviceName)
					}
				}
			}
		}
	}

	return nil
}

//nolint:unused // This function is part of the validation system and may be used by other components
func (e *Engine) validateDockerfileSecurityFile(filePath string) error {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Check for running as root
		if strings.HasPrefix(line, "USER root") {
			return fmt.Errorf("line %d: running as root user is a security risk", lineNum+1)
		}

		// Check for ADD instead of COPY
		if strings.HasPrefix(line, "ADD ") && !strings.Contains(line, ".tar") {
			return fmt.Errorf("line %d: use COPY instead of ADD for better security", lineNum+1)
		}
	}

	return nil
}

// maskSecret masks sensitive information for display
//
//nolint:unused // This function is part of the validation system and may be used by other components
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return strings.Repeat("*", len(secret))
	}
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}
