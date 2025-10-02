package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// Pre-compiled regular expressions for dependency validation
var (
	goVersionRegex      = regexp.MustCompile(`^1\.\d+(\.\d+)?$`)
	npmNameRegex        = regexp.MustCompile(`^(@[a-z0-9-~][a-z0-9-._~]*/)?[a-z0-9-~][a-z0-9-._~]*$`)
	versionExtractRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
	semverRegex         = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	pythonNameRegex     = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$`)
)

// DependencyValidator provides specialized dependency validation
type DependencyValidator struct {
	vulnerabilityDB map[string][]interfaces.DependencyVulnerability
	packageRegistry map[string]string // package name -> latest version
}

// NewDependencyValidator creates a new dependency validator
func NewDependencyValidator() *DependencyValidator {
	validator := &DependencyValidator{
		vulnerabilityDB: make(map[string][]interfaces.DependencyVulnerability),
		packageRegistry: make(map[string]string),
	}
	validator.initializeKnownVulnerabilities()
	return validator
}

// ValidateProjectDependencies validates all project dependencies
func (dv *DependencyValidator) ValidateProjectDependencies(projectPath string) (*interfaces.DependencyValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(projectPath); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	result := &interfaces.DependencyValidationResult{
		Valid:           true,
		Dependencies:    []interfaces.DependencyValidation{},
		Vulnerabilities: []interfaces.DependencyVulnerability{},
		Outdated:        []interfaces.OutdatedDependency{},
		Conflicts:       []interfaces.DependencyConflict{},
		Summary: interfaces.DependencyValidationSummary{
			TotalDependencies: 0,
			ValidDependencies: 0,
			Vulnerabilities:   0,
			OutdatedCount:     0,
			ConflictCount:     0,
		},
	}

	// Validate package.json dependencies
	if err := dv.validatePackageJSONDependencies(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate package.json dependencies: %w", err)
	}

	// Validate go.mod dependencies
	if err := dv.validateGoModDependencies(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate go.mod dependencies: %w", err)
	}

	// Validate requirements.txt dependencies (Python)
	if err := dv.validateRequirementsTxtDependencies(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate requirements.txt dependencies: %w", err)
	}

	// Check for dependency conflicts
	if err := dv.checkDependencyConflicts(result); err != nil {
		return nil, fmt.Errorf("failed to check dependency conflicts: %w", err)
	}

	// Check for vulnerabilities
	if err := dv.checkVulnerabilities(result); err != nil {
		return nil, fmt.Errorf("failed to check vulnerabilities: %w", err)
	}

	// Check for outdated dependencies
	if err := dv.checkOutdatedDependencies(result); err != nil {
		return nil, fmt.Errorf("failed to check outdated dependencies: %w", err)
	}

	return result, nil
}

// validatePackageJSONDependencies validates Node.js package.json dependencies
func (dv *DependencyValidator) validatePackageJSONDependencies(projectPath string, result *interfaces.DependencyValidationResult) error {
	packageJsonPath := filepath.Join(projectPath, "package.json")
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
		if err := dv.validateDependencyMap(deps, "production", result); err != nil {
			return fmt.Errorf("failed to validate dependencies: %w", err)
		}
	}

	// Validate devDependencies
	if devDeps, exists := pkg["devDependencies"]; exists {
		if err := dv.validateDependencyMap(devDeps, "development", result); err != nil {
			return fmt.Errorf("failed to validate devDependencies: %w", err)
		}
	}

	// Validate peerDependencies
	if peerDeps, exists := pkg["peerDependencies"]; exists {
		if err := dv.validateDependencyMap(peerDeps, "peer", result); err != nil {
			return fmt.Errorf("failed to validate peerDependencies: %w", err)
		}
	}

	// Validate package.json structure
	if err := dv.validatePackageJSONStructure(pkg, result); err != nil {
		return fmt.Errorf("failed to validate package.json structure: %w", err)
	}

	return nil
}

// validateGoModDependencies validates Go module dependencies
func (dv *DependencyValidator) validateGoModDependencies(projectPath string, result *interfaces.DependencyValidationResult) error {
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil // No go.mod to validate
	}

	content, err := utils.SafeReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	inRequireBlock := false
	moduleVersion := ""

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Parse module declaration
		if strings.HasPrefix(line, "module ") {
			// Module name validation is handled elsewhere
			continue
		}

		// Parse Go version
		if strings.HasPrefix(line, "go ") {
			moduleVersion = strings.TrimPrefix(line, "go ")
			if err := dv.validateGoVersion(moduleVersion, result); err != nil {
				return fmt.Errorf("invalid Go version at line %d: %w", lineNum+1, err)
			}
			continue
		}

		// Parse require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		// Parse individual require statements
		if inRequireBlock || strings.HasPrefix(line, "require ") {
			if err := dv.parseGoRequirement(line, result); err != nil {
				return fmt.Errorf("failed to parse requirement at line %d: %w", lineNum+1, err)
			}
		}
	}

	return nil
}

// validateRequirementsTxtDependencies validates Python requirements.txt dependencies
func (dv *DependencyValidator) validateRequirementsTxtDependencies(projectPath string, result *interfaces.DependencyValidationResult) error {
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		return nil // No requirements.txt to validate
	}

	content, err := utils.SafeReadFile(requirementsPath)
	if err != nil {
		return fmt.Errorf("failed to read requirements.txt: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := dv.parsePythonRequirement(line, result); err != nil {
			return fmt.Errorf("failed to parse requirement at line %d: %w", lineNum+1, err)
		}
	}

	return nil
}

// validateDependencyMap validates a dependency map from package.json
func (dv *DependencyValidator) validateDependencyMap(deps interface{}, depType string, result *interfaces.DependencyValidationResult) error {
	depsMap, ok := deps.(map[string]interface{})
	if !ok {
		return fmt.Errorf("dependencies must be an object")
	}

	for name, version := range depsMap {
		versionStr, ok := version.(string)
		if !ok {
			continue
		}

		dep := interfaces.DependencyValidation{
			Name:           name,
			Version:        versionStr,
			Type:           depType,
			Valid:          true,
			Available:      true,
			LatestVersion:  versionStr,
			SecurityIssues: 0,
			LicenseIssues:  0,
		}

		// Validate version format
		if err := dv.validateNpmVersionFormat(versionStr); err != nil {
			dep.Valid = false
		}

		// Check for known security issues
		if vulns, exists := dv.vulnerabilityDB[name]; exists {
			dep.SecurityIssues = len(vulns)
			for _, vuln := range vulns {
				if dv.versionAffectedByVulnerability(versionStr, vuln) {
					result.Vulnerabilities = append(result.Vulnerabilities, vuln)
					result.Summary.Vulnerabilities++
				}
			}
		}

		result.Dependencies = append(result.Dependencies, dep)
		result.Summary.TotalDependencies++

		if dep.Valid {
			result.Summary.ValidDependencies++
		} else {
			result.Valid = false
		}
	}

	return nil
}

// validatePackageJSONStructure validates the structure of package.json
func (dv *DependencyValidator) validatePackageJSONStructure(pkg map[string]interface{}, result *interfaces.DependencyValidationResult) error {
	// Check required fields
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		if _, exists := pkg[field]; !exists {
			result.Valid = false
			// Note: This would typically add to a validation issues array
		}
	}

	// Validate name format
	if name, exists := pkg["name"]; exists {
		if nameStr, ok := name.(string); ok {
			if err := dv.validateNpmPackageName(nameStr); err != nil {
				result.Valid = false
			}
		}
	}

	// Validate version format
	if version, exists := pkg["version"]; exists {
		if versionStr, ok := version.(string); ok {
			if err := dv.validateNpmVersionFormat(versionStr); err != nil {
				result.Valid = false
			}
		}
	}

	return nil
}

// validateGoVersion validates Go version format
func (dv *DependencyValidator) validateGoVersion(version string, result *interfaces.DependencyValidationResult) error {
	// Go version should be in format like "1.19", "1.20", etc.
	if !goVersionRegex.MatchString(version) {
		result.Valid = false
		return fmt.Errorf("invalid Go version format: %s", version)
	}

	// Check if version is too old
	if strings.HasPrefix(version, "1.1") && len(version) == 4 {
		// Versions like 1.10, 1.11, etc. are old
		majorMinor := version[2:]
		if len(majorMinor) == 2 && majorMinor < "18" {
			result.Valid = false
			return fmt.Errorf("go version %s is too old, consider upgrading to 1.18+", version)
		}
	}

	return nil
}

// parseGoRequirement parses a Go module requirement
func (dv *DependencyValidator) parseGoRequirement(line string, result *interfaces.DependencyValidationResult) error {
	// Remove "require " prefix if present
	line = strings.TrimPrefix(line, "require ")
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid require statement: %s", line)
	}

	name := parts[0]
	version := parts[1]

	// Handle indirect dependencies
	depType := "direct"
	if len(parts) > 2 && parts[2] == "// indirect" {
		depType = "indirect"
	}

	dep := interfaces.DependencyValidation{
		Name:           name,
		Version:        version,
		Type:           depType,
		Valid:          true,
		Available:      true,
		LatestVersion:  version,
		SecurityIssues: 0,
		LicenseIssues:  0,
	}

	// Validate version format (semantic versioning)
	if err := dv.validateSemanticVersion(version); err != nil {
		dep.Valid = false
	}

	// Check for known vulnerabilities
	if vulns, exists := dv.vulnerabilityDB[name]; exists {
		dep.SecurityIssues = len(vulns)
		for _, vuln := range vulns {
			if dv.versionAffectedByVulnerability(version, vuln) {
				result.Vulnerabilities = append(result.Vulnerabilities, vuln)
				result.Summary.Vulnerabilities++
			}
		}
	}

	result.Dependencies = append(result.Dependencies, dep)
	result.Summary.TotalDependencies++

	if dep.Valid {
		result.Summary.ValidDependencies++
	} else {
		result.Valid = false
	}

	return nil
}

// parsePythonRequirement parses a Python requirement
func (dv *DependencyValidator) parsePythonRequirement(line string, result *interfaces.DependencyValidationResult) error {
	// Parse requirements like "package==1.0.0", "package>=1.0.0", etc.
	operators := []string{"==", ">=", "<=", ">", "<", "~=", "!="}

	var name, version string

	for _, op := range operators {
		if strings.Contains(line, op) {
			parts := strings.Split(line, op)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				version = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if name == "" {
		// No version specified, just package name
		name = strings.TrimSpace(line)
		version = "latest"
	}

	dep := interfaces.DependencyValidation{
		Name:           name,
		Version:        version,
		Type:           "direct",
		Valid:          true,
		Available:      true,
		LatestVersion:  version,
		SecurityIssues: 0,
		LicenseIssues:  0,
	}

	// Validate package name
	if err := dv.validatePythonPackageName(name); err != nil {
		dep.Valid = false
	}

	// Validate version format if specified
	if version != "latest" && version != "" {
		if err := dv.validateSemanticVersion(version); err != nil {
			dep.Valid = false
		}
	}

	result.Dependencies = append(result.Dependencies, dep)
	result.Summary.TotalDependencies++

	if dep.Valid {
		result.Summary.ValidDependencies++
	} else {
		result.Valid = false
	}

	return nil
}

// checkDependencyConflicts checks for conflicts between dependencies
func (dv *DependencyValidator) checkDependencyConflicts(result *interfaces.DependencyValidationResult) error {
	// Group dependencies by name
	depsByName := make(map[string][]interfaces.DependencyValidation)

	for _, dep := range result.Dependencies {
		depsByName[dep.Name] = append(depsByName[dep.Name], dep)
	}

	// Check for version conflicts
	for name, deps := range depsByName {
		if len(deps) > 1 {
			// Check if all versions are compatible
			for i := 0; i < len(deps); i++ {
				for j := i + 1; j < len(deps); j++ {
					if deps[i].Version != deps[j].Version {
						conflict := interfaces.DependencyConflict{
							Dependency1: name,
							Version1:    deps[i].Version,
							Dependency2: name,
							Version2:    deps[j].Version,
							Reason:      "Multiple versions of the same dependency",
							Severity:    interfaces.ValidationSeverityWarning,
						}
						result.Conflicts = append(result.Conflicts, conflict)
						result.Summary.ConflictCount++
					}
				}
			}
		}
	}

	return nil
}

// checkVulnerabilities checks for known vulnerabilities
func (dv *DependencyValidator) checkVulnerabilities(result *interfaces.DependencyValidationResult) error {
	// Vulnerabilities are already checked during dependency parsing
	// This method can be extended to check against external vulnerability databases
	return nil
}

// checkOutdatedDependencies checks for outdated dependencies
func (dv *DependencyValidator) checkOutdatedDependencies(result *interfaces.DependencyValidationResult) error {
	for _, dep := range result.Dependencies {
		// Simple check - if we have a "latest" version in our registry
		if latestVersion, exists := dv.packageRegistry[dep.Name]; exists {
			if dep.Version != latestVersion && dv.isVersionOlder(dep.Version, latestVersion) {
				outdated := interfaces.OutdatedDependency{
					Name:           dep.Name,
					CurrentVersion: dep.Version,
					LatestVersion:  latestVersion,
					UpdateType:     dv.getUpdateType(dep.Version, latestVersion),
					Breaking:       dv.isBreakingUpdate(dep.Version, latestVersion),
				}
				result.Outdated = append(result.Outdated, outdated)
				result.Summary.OutdatedCount++
			}
		}
	}

	return nil
}

// Validation helper methods

// validateNpmPackageName validates NPM package name format
func (dv *DependencyValidator) validateNpmPackageName(name string) error {
	// NPM package names must be lowercase and can contain hyphens, dots, and underscores
	if !npmNameRegex.MatchString(name) {
		return fmt.Errorf("invalid NPM package name: %s", name)
	}
	return nil
}

// validateNpmVersionFormat validates NPM version format
func (dv *DependencyValidator) validateNpmVersionFormat(version string) error {
	// Handle version ranges and special characters
	if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") ||
		strings.HasPrefix(version, ">=") || strings.HasPrefix(version, "<=") ||
		strings.HasPrefix(version, ">") || strings.HasPrefix(version, "<") {
		// Extract the actual version part using pre-compiled regex
		versionPart := versionExtractRegex.FindString(version)
		if versionPart != "" {
			return dv.validateSemanticVersion(versionPart)
		}
	}

	return dv.validateSemanticVersion(version)
}

// validateSemanticVersion validates semantic version format
func (dv *DependencyValidator) validateSemanticVersion(version string) error {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	if !semverRegex.MatchString(version) {
		return fmt.Errorf("invalid semantic version: %s", version)
	}
	return nil
}

// validatePythonPackageName validates Python package name format
func (dv *DependencyValidator) validatePythonPackageName(name string) error {
	// Python package names can contain letters, numbers, hyphens, underscores, and dots
	if !pythonNameRegex.MatchString(name) {
		return fmt.Errorf("invalid Python package name: %s", name)
	}
	return nil
}

// versionAffectedByVulnerability checks if a version is affected by a vulnerability
func (dv *DependencyValidator) versionAffectedByVulnerability(version string, vuln interfaces.DependencyVulnerability) bool {
	// Simple check - if the vulnerability version matches exactly
	return version == vuln.Version
}

// isVersionOlder checks if version1 is older than version2
func (dv *DependencyValidator) isVersionOlder(version1, version2 string) bool {
	// Simple string comparison for now
	// In a real implementation, this would use proper semantic version comparison
	return version1 < version2
}

// getUpdateType determines the type of update (major, minor, patch)
func (dv *DependencyValidator) getUpdateType(currentVersion, latestVersion string) string {
	// Simple implementation - would need proper semver parsing in production
	if strings.HasPrefix(latestVersion, strings.Split(currentVersion, ".")[0]) {
		return "minor"
	}
	return "major"
}

// isBreakingUpdate checks if an update is potentially breaking
func (dv *DependencyValidator) isBreakingUpdate(currentVersion, latestVersion string) bool {
	// Major version changes are typically breaking
	currentMajor := strings.Split(currentVersion, ".")[0]
	latestMajor := strings.Split(latestVersion, ".")[0]
	return currentMajor != latestMajor
}

// initializeKnownVulnerabilities initializes the vulnerability database with known issues
func (dv *DependencyValidator) initializeKnownVulnerabilities() {
	// Example vulnerabilities - in production, this would be loaded from a real database
	dv.vulnerabilityDB["lodash"] = []interfaces.DependencyVulnerability{
		{
			Dependency:  "lodash",
			Version:     "4.17.15",
			CVEID:       "CVE-2020-8203",
			Severity:    "high",
			Description: "Prototype pollution vulnerability",
			FixedIn:     "4.17.19",
			CVSS:        7.4,
		},
	}

	dv.vulnerabilityDB["express"] = []interfaces.DependencyVulnerability{
		{
			Dependency:  "express",
			Version:     "4.16.0",
			CVEID:       "CVE-2022-24999",
			Severity:    "medium",
			Description: "Open redirect vulnerability",
			FixedIn:     "4.17.3",
			CVSS:        6.1,
		},
	}

	// Initialize package registry with some known latest versions
	dv.packageRegistry["lodash"] = "4.17.21"
	dv.packageRegistry["express"] = "4.18.2"
	dv.packageRegistry["react"] = "18.2.0"
	dv.packageRegistry["vue"] = "3.3.4"
}
