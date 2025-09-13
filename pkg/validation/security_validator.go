package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// SecurityValidator handles security vulnerability checking
type SecurityValidator struct {
	client            *http.Client
	npmAuditURL       string
	goVulnDBURL       string
	githubAdvisoryURL string
	cache             map[string]*SecurityResult
}

// SecurityResult represents the result of a security check
type SecurityResult struct {
	Package         string              `json:"package"`
	Version         string              `json:"version"`
	Vulnerabilities []VulnerabilityInfo `json:"vulnerabilities"`
	LastChecked     time.Time           `json:"last_checked"`
}

// VulnerabilityInfo represents information about a vulnerability
type VulnerabilityInfo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	References  []string `json:"references"`
	FixedIn     string   `json:"fixed_in,omitempty"`
}

// NPMAuditResponse represents the response from npm audit
type NPMAuditResponse struct {
	Advisories map[string]NPMAdvisory `json:"advisories"`
	Metadata   NPMMetadata            `json:"metadata"`
}

// NPMAdvisory represents an npm security advisory
type NPMAdvisory struct {
	ID                 int      `json:"id"`
	Title              string   `json:"title"`
	ModuleName         string   `json:"module_name"`
	VulnerableVersions string   `json:"vulnerable_versions"`
	PatchedVersions    string   `json:"patched_versions"`
	Severity           string   `json:"severity"`
	Overview           string   `json:"overview"`
	References         []string `json:"references"`
}

// NPMMetadata represents metadata from npm audit
type NPMMetadata struct {
	Vulnerabilities NPMVulnCounts `json:"vulnerabilities"`
}

// NPMVulnCounts represents vulnerability counts by severity
type NPMVulnCounts struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		npmAuditURL:       "https://registry.npmjs.org/-/npm/v1/security/audits",
		goVulnDBURL:       "https://vuln.go.dev",
		githubAdvisoryURL: "https://api.github.com/advisories",
		cache:             make(map[string]*SecurityResult),
	}
}

// ValidateSecurityVulnerabilities validates packages for security vulnerabilities
func (sv *SecurityValidator) ValidateSecurityVulnerabilities(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check Node.js packages
	if err := sv.validateNPMPackages(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate npm packages: %w", err)
	}

	// Check Go modules
	if err := sv.validateGoModules(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate go modules: %w", err)
	}

	return result, nil
}

// ValidatePackageVersionSecurity validates a specific package version for security issues
func (sv *SecurityValidator) ValidatePackageVersionSecurity(packageName, version, ecosystem string) (*SecurityResult, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", ecosystem, packageName, version)

	// Check cache first
	if cached, exists := sv.cache[cacheKey]; exists {
		if time.Since(cached.LastChecked) < 24*time.Hour {
			return cached, nil
		}
	}

	var result *SecurityResult
	var err error

	switch ecosystem {
	case "npm":
		result, err = sv.checkNPMPackageSecurity(packageName, version)
	case "go":
		result, err = sv.checkGoModuleSecurity(packageName, version)
	default:
		return nil, fmt.Errorf("unsupported ecosystem: %s", ecosystem)
	}

	if err != nil {
		return nil, err
	}

	// Cache the result
	sv.cache[cacheKey] = result
	return result, nil
}

// PrioritizeSecurityUpdates analyzes vulnerabilities and prioritizes updates
func (sv *SecurityValidator) PrioritizeSecurityUpdates(vulnerabilities []VulnerabilityInfo) []VulnerabilityInfo {
	// Sort by severity: critical > high > moderate > low > info
	severityOrder := map[string]int{
		"critical": 5,
		"high":     4,
		"moderate": 3,
		"low":      2,
		"info":     1,
	}

	// Create a copy to avoid modifying the original slice
	prioritized := make([]VulnerabilityInfo, len(vulnerabilities))
	copy(prioritized, vulnerabilities)

	// Simple bubble sort by severity
	for i := 0; i < len(prioritized)-1; i++ {
		for j := 0; j < len(prioritized)-i-1; j++ {
			severity1 := severityOrder[strings.ToLower(prioritized[j].Severity)]
			severity2 := severityOrder[strings.ToLower(prioritized[j+1].Severity)]

			if severity1 < severity2 {
				prioritized[j], prioritized[j+1] = prioritized[j+1], prioritized[j]
			}
		}
	}

	return prioritized
}

// validateNPMPackages validates npm packages for security vulnerabilities
func (sv *SecurityValidator) validateNPMPackages(projectPath string, result *models.ValidationResult) error {
	// Find all package.json files
	packageJSONFiles, err := sv.findPackageJSONFiles(projectPath)
	if err != nil {
		return err
	}

	for _, packageJSONPath := range packageJSONFiles {
		if err := sv.validateSinglePackageJSON(packageJSONPath, result); err != nil {
			return err
		}
	}

	return nil
}

// validateGoModules validates Go modules for security vulnerabilities
func (sv *SecurityValidator) validateGoModules(projectPath string, result *models.ValidationResult) error {
	// Find all go.mod files
	goModFiles, err := sv.findGoModFiles(projectPath)
	if err != nil {
		return err
	}

	for _, goModPath := range goModFiles {
		if err := sv.validateSingleGoMod(goModPath, result); err != nil {
			return err
		}
	}

	return nil
}

// validateSinglePackageJSON validates a single package.json file
func (sv *SecurityValidator) validateSinglePackageJSON(packageJSONPath string, result *models.ValidationResult) error {
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Check dependencies
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		sv.validateDependencyMap(deps, "prod", result)
	}

	// Check devDependencies
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		sv.validateDependencyMap(devDeps, "dev", result)
	}

	return nil
}

// validateDependencyMap validates a map of dependencies (either dependencies or devDependencies)
func (sv *SecurityValidator) validateDependencyMap(deps map[string]interface{}, depType string, result *models.ValidationResult) {
	for packageName, version := range deps {
		if versionStr, ok := version.(string); ok {
			cleanVersion := sv.cleanVersion(versionStr)
			securityResult, err := sv.ValidatePackageVersionSecurity(packageName, cleanVersion, "npm")
			if err != nil {
				// Log warning but don't fail validation
				message := fmt.Sprintf("Failed to check security for %s@%s: %s", packageName, cleanVersion, err.Error())
				if depType == "dev" {
					message = fmt.Sprintf("Failed to check security for dev dependency %s@%s: %s", packageName, cleanVersion, err.Error())
				}
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "SecurityCheck",
					Message: message,
				})
				continue
			}

			if len(securityResult.Vulnerabilities) > 0 {
				sv.addVulnerabilityWarnings(packageName, cleanVersion, securityResult.Vulnerabilities, result)
			}
		}
	}
}

// validateSingleGoMod validates a single go.mod file
func (sv *SecurityValidator) validateSingleGoMod(goModPath string, result *models.ValidationResult) error {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require ") {
			// Parse require line: "require github.com/example/package v1.2.3"
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				packageName := parts[1]
				version := parts[2]

				securityResult, err := sv.ValidatePackageVersionSecurity(packageName, version, "go")
				if err != nil {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "SecurityCheck",
						Message: fmt.Sprintf("Failed to check security for Go module %s@%s: %s", packageName, version, err.Error()),
					})
					continue
				}

				if len(securityResult.Vulnerabilities) > 0 {
					sv.addVulnerabilityWarnings(packageName, version, securityResult.Vulnerabilities, result)
				}
			}
		}
	}

	return nil
}

// checkNPMPackageSecurity checks npm package security using a mock implementation
func (sv *SecurityValidator) checkNPMPackageSecurity(packageName, version string) (*SecurityResult, error) {
	// This is a simplified implementation for demonstration
	// In a real implementation, this would call npm audit API or security databases

	result := &SecurityResult{
		Package:         packageName,
		Version:         version,
		Vulnerabilities: []VulnerabilityInfo{},
		LastChecked:     time.Now(),
	}

	// Mock some known vulnerable packages for testing
	knownVulnerablePackages := map[string][]VulnerabilityInfo{
		"lodash": {
			{
				ID:          "CVE-2021-23337",
				Title:       "Command Injection in lodash",
				Severity:    "high",
				Description: "lodash versions prior to 4.17.21 are vulnerable to Command Injection via template.",
				References:  []string{"https://nvd.nist.gov/vuln/detail/CVE-2021-23337"},
				FixedIn:     "4.17.21",
			},
		},
		"axios": {
			{
				ID:          "CVE-2021-3749",
				Title:       "Regular Expression Denial of Service in axios",
				Severity:    "moderate",
				Description: "axios is vulnerable to Inefficient Regular Expression Complexity.",
				References:  []string{"https://nvd.nist.gov/vuln/detail/CVE-2021-3749"},
				FixedIn:     "0.21.2",
			},
		},
	}

	if vulns, exists := knownVulnerablePackages[packageName]; exists {
		// Check if the current version is affected
		for _, vuln := range vulns {
			if sv.isVersionAffected(version, vuln.FixedIn) {
				result.Vulnerabilities = append(result.Vulnerabilities, vuln)
			}
		}
	}

	return result, nil
}

// checkGoModuleSecurity checks Go module security using a mock implementation
func (sv *SecurityValidator) checkGoModuleSecurity(packageName, version string) (*SecurityResult, error) {
	// This is a simplified implementation for demonstration
	// In a real implementation, this would call Go vulnerability database

	result := &SecurityResult{
		Package:         packageName,
		Version:         version,
		Vulnerabilities: []VulnerabilityInfo{},
		LastChecked:     time.Now(),
	}

	// Mock some known vulnerable Go modules for testing
	knownVulnerableModules := map[string][]VulnerabilityInfo{
		"github.com/gin-gonic/gin": {
			{
				ID:          "GO-2021-0052",
				Title:       "Improper input validation in github.com/gin-gonic/gin",
				Severity:    "moderate",
				Description: "A maliciously crafted HTTP request can cause gin to panic.",
				References:  []string{"https://pkg.go.dev/vuln/GO-2021-0052"},
				FixedIn:     "v1.7.0",
			},
		},
	}

	if vulns, exists := knownVulnerableModules[packageName]; exists {
		for _, vuln := range vulns {
			if sv.isVersionAffected(version, vuln.FixedIn) {
				result.Vulnerabilities = append(result.Vulnerabilities, vuln)
			}
		}
	}

	return result, nil
}

// addVulnerabilityWarnings adds vulnerability warnings to the result
func (sv *SecurityValidator) addVulnerabilityWarnings(packageName, version string, vulnerabilities []VulnerabilityInfo, result *models.ValidationResult) {
	prioritized := sv.PrioritizeSecurityUpdates(vulnerabilities)

	for _, vuln := range prioritized {
		severity := strings.ToUpper(vuln.Severity)
		message := fmt.Sprintf("Security vulnerability in %s@%s: %s (%s severity)",
			packageName, version, vuln.Title, severity)

		if vuln.FixedIn != "" {
			message += fmt.Sprintf(" - Fixed in %s", vuln.FixedIn)
		}

		if severity == "CRITICAL" || severity == "HIGH" {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "SecurityVulnerability",
				Tag:     "security",
				Value:   fmt.Sprintf("%s@%s", packageName, version),
				Message: message,
			})
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "SecurityVulnerability",
				Message: message,
			})
		}
	}
}

// Helper functions

// findPackageJSONFiles finds all package.json files in the project
func (sv *SecurityValidator) findPackageJSONFiles(projectPath string) ([]string, error) {
	var packageJSONFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "package.json" {
			// Skip node_modules directories
			if !strings.Contains(path, "node_modules") {
				packageJSONFiles = append(packageJSONFiles, path)
			}
		}

		return nil
	})

	return packageJSONFiles, err
}

// findGoModFiles finds all go.mod files in the project
func (sv *SecurityValidator) findGoModFiles(projectPath string) ([]string, error) {
	var goModFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "go.mod" {
			goModFiles = append(goModFiles, path)
		}

		return nil
	})

	return goModFiles, err
}

// cleanVersion removes version prefixes like ^, ~, >=, etc.
func (sv *SecurityValidator) cleanVersion(version string) string {
	// Remove common npm version prefixes
	prefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}

	for _, prefix := range prefixes {
		if strings.HasPrefix(version, prefix) {
			return strings.TrimPrefix(version, prefix)
		}
	}

	return version
}

// isVersionAffected checks if a version is affected by a vulnerability
func (sv *SecurityValidator) isVersionAffected(currentVersion, fixedVersion string) bool {
	// This is a simplified version comparison
	// In a real implementation, you would use proper semantic version comparison

	if fixedVersion == "" {
		return true // No fix available, assume affected
	}

	// Remove 'v' prefix if present
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	fixedVersion = strings.TrimPrefix(fixedVersion, "v")

	// Simple string comparison (not semantically correct, but good enough for demo)
	return currentVersion < fixedVersion
}
