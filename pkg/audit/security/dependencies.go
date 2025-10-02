// Package security provides security auditing functionality for the
// Open Source Project Generator.
//
// Security Note: This package contains security audit functionality that legitimately needs
// to read files for security analysis. The G304 warnings from gosec are false positives
// in this context as file reading is the core functionality of a security audit tool.
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DependencyChecker handles dependency vulnerability scanning and analysis.
type DependencyChecker struct {
	// vulnerabilityDatabase contains known vulnerability patterns
	vulnerabilityDatabase map[string][]interfaces.Vulnerability
	// ecosystemParsers contains parsers for different dependency ecosystems
	ecosystemParsers map[string]DependencyParser
}

// DependencyParser defines the interface for parsing dependency files
type DependencyParser interface {
	ParseDependencies(content []byte) ([]interfaces.DependencyInfo, error)
	GetEcosystem() string
}

// NewDependencyChecker creates a new dependency checker instance.
func NewDependencyChecker() *DependencyChecker {
	checker := &DependencyChecker{
		vulnerabilityDatabase: make(map[string][]interfaces.Vulnerability),
		ecosystemParsers:      make(map[string]DependencyParser),
	}

	// Initialize vulnerability database with known vulnerabilities
	checker.initializeVulnerabilityDatabase()

	// Initialize ecosystem parsers
	checker.initializeEcosystemParsers()

	return checker
}

// ScanDependencyVulnerabilities scans dependency files for known vulnerabilities
func (dc *DependencyChecker) ScanDependencyVulnerabilities(path string) ([]interfaces.Vulnerability, error) {
	var vulnerabilities []interfaces.Vulnerability

	// Check for common dependency files
	dependencyFiles := []string{
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"go.mod",
		"go.sum",
		"requirements.txt",
		"Pipfile",
		"Pipfile.lock",
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
		"Cargo.toml",
		"Cargo.lock",
		"composer.json",
		"composer.lock",
		"Gemfile",
		"Gemfile.lock",
	}

	for _, depFile := range dependencyFiles {
		filePath := filepath.Join(path, depFile)
		if _, err := os.Stat(filePath); err == nil {
			vulns, err := dc.scanDependencyFile(filePath)
			if err != nil {
				// Log error but continue scanning other files
				continue
			}
			vulnerabilities = append(vulnerabilities, vulns...)
		}
	}

	return vulnerabilities, nil
}

// AnalyzeDependencies analyzes project dependencies for security issues
func (dc *DependencyChecker) AnalyzeDependencies(path string) (*interfaces.DependencyAnalysisResult, error) {
	result := &interfaces.DependencyAnalysisResult{
		Dependencies:    []interfaces.DependencyInfo{},
		Vulnerabilities: []interfaces.DependencyVulnerability{},
		Licenses:        []interfaces.DependencyLicense{},
		Outdated:        []interfaces.OutdatedDependency{},
		Summary: interfaces.DependencyAnalysisSummary{
			TotalDependencies:  0,
			DirectDependencies: 0,
			Vulnerabilities:    0,
			OutdatedCount:      0,
			LicenseIssues:      0,
			AverageAge:         0.0,
		},
	}

	// Analyze different types of dependency files
	dependencyFiles := map[string]string{
		"package.json":     "npm",
		"go.mod":           "go",
		"requirements.txt": "pip",
		"Pipfile":          "pipenv",
		"pom.xml":          "maven",
		"build.gradle":     "gradle",
		"Cargo.toml":       "cargo",
		"composer.json":    "composer",
		"Gemfile":          "bundler",
	}

	for depFile, ecosystem := range dependencyFiles {
		filePath := filepath.Join(path, depFile)
		if _, err := os.Stat(filePath); err == nil {
			deps, err := dc.analyzeDependencyFile(filePath, ecosystem)
			if err != nil {
				continue // Log error but continue
			}
			result.Dependencies = append(result.Dependencies, deps...)
		}
	}

	// Update summary
	result.Summary.TotalDependencies = len(result.Dependencies)

	// Count direct dependencies and calculate metrics
	var totalAge float64
	for _, dep := range result.Dependencies {
		if dep.Type == "direct" {
			result.Summary.DirectDependencies++
		}

		// Calculate age in days
		age := time.Since(dep.LastUpdated).Hours() / 24
		totalAge += age

		// Check for vulnerabilities
		if dep.SecurityIssues > 0 {
			result.Summary.Vulnerabilities += dep.SecurityIssues
		}
	}

	if len(result.Dependencies) > 0 {
		result.Summary.AverageAge = totalAge / float64(len(result.Dependencies))
	}

	// Scan for vulnerabilities and convert to DependencyVulnerability
	vulnerabilities, err := dc.ScanDependencyVulnerabilities(path)
	if err == nil {
		for _, vuln := range vulnerabilities {
			depVuln := interfaces.DependencyVulnerability{
				Dependency:  vuln.Package,
				Version:     vuln.Version,
				CVEID:       vuln.ID,
				Severity:    vuln.Severity,
				Description: vuln.Description,
				FixedIn:     vuln.FixedIn,
				CVSS:        0.0, // Would be populated from vulnerability database
			}
			result.Vulnerabilities = append(result.Vulnerabilities, depVuln)
		}
	}

	return result, nil
}

// CheckDependencyVulnerabilities checks for dependency vulnerabilities and returns policy violations
func (dc *DependencyChecker) CheckDependencyVulnerabilities(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	vulnerabilities, err := dc.ScanDependencyVulnerabilities(path)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dependency vulnerabilities: %w", err)
	}

	for _, vuln := range vulnerabilities {
		violations = append(violations, interfaces.PolicyViolation{
			Policy:      "SEC-002",
			Severity:    vuln.Severity,
			Description: fmt.Sprintf("Vulnerable dependency: %s %s - %s", vuln.Package, vuln.Version, vuln.Title),
			File:        "", // Would need to determine which file contains the dependency
			Line:        0,
		})
	}

	return violations, nil
}

// scanDependencyFile scans a specific dependency file for vulnerabilities
func (dc *DependencyChecker) scanDependencyFile(filePath string) ([]interfaces.Vulnerability, error) {
	var vulnerabilities []interfaces.Vulnerability

	// Read the dependency file
	content, err := os.ReadFile(filePath) // #nosec G304 - This is an audit tool that needs to read files
	if err != nil {
		return nil, fmt.Errorf("failed to read dependency file %s: %w", filePath, err)
	}

	// Determine the ecosystem based on file name
	ecosystem := dc.getEcosystemFromFile(filePath)
	if ecosystem == "" {
		return vulnerabilities, nil // Unknown ecosystem, skip
	}

	// Parse dependencies using the appropriate parser
	dependencies, err := dc.parseDependencies(content, ecosystem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dependencies in %s: %w", filePath, err)
	}

	// Check each dependency against the vulnerability database
	for _, dep := range dependencies {
		vulns := dc.checkDependencyForVulnerabilities(dep.Name, dep.Version, ecosystem)
		vulnerabilities = append(vulnerabilities, vulns...)
	}

	// Also perform pattern-based scanning for additional vulnerabilities
	patternVulns := dc.scanForVulnerabilityPatterns(string(content), ecosystem)
	vulnerabilities = append(vulnerabilities, patternVulns...)

	return vulnerabilities, nil
}

// analyzeDependencyFile analyzes a dependency file for the given ecosystem
func (dc *DependencyChecker) analyzeDependencyFile(filePath, ecosystem string) ([]interfaces.DependencyInfo, error) {
	content, err := os.ReadFile(filePath) // #nosec G304 - This is an audit tool that needs to read files
	if err != nil {
		return nil, fmt.Errorf("failed to read dependency file %s: %w", filePath, err)
	}

	return dc.parseDependencies(content, ecosystem)
}

// parseDependencies parses dependencies from content based on ecosystem
func (dc *DependencyChecker) parseDependencies(content []byte, ecosystem string) ([]interfaces.DependencyInfo, error) {
	parser, exists := dc.ecosystemParsers[ecosystem]
	if !exists {
		// Fallback to simple parsing
		return dc.parseGenericDependencies(content, ecosystem)
	}

	return parser.ParseDependencies(content)
}

// parseGenericDependencies provides a fallback parser for unsupported ecosystems
func (dc *DependencyChecker) parseGenericDependencies(content []byte, ecosystem string) ([]interfaces.DependencyInfo, error) {
	var dependencies []interfaces.DependencyInfo
	contentStr := string(content)

	switch ecosystem {
	case "npm":
		dependencies = dc.parseNPMDependencies(contentStr)
	case "go":
		dependencies = dc.parseGoDependencies(contentStr)
	case "pip":
		dependencies = dc.parsePipDependencies(contentStr)
	case "maven":
		dependencies = dc.parseMavenDependencies(contentStr)
	case "gradle":
		dependencies = dc.parseGradleDependencies(contentStr)
	case "cargo":
		dependencies = dc.parseCargoDependencies(contentStr)
	case "composer":
		dependencies = dc.parseComposerDependencies(contentStr)
	case "bundler":
		dependencies = dc.parseBundlerDependencies(contentStr)
	}

	return dependencies, nil
}

// getEcosystemFromFile determines the ecosystem based on the file name
func (dc *DependencyChecker) getEcosystemFromFile(filePath string) string {
	fileName := filepath.Base(filePath)

	ecosystemMap := map[string]string{
		"package.json":      "npm",
		"package-lock.json": "npm",
		"yarn.lock":         "npm",
		"go.mod":            "go",
		"go.sum":            "go",
		"requirements.txt":  "pip",
		"Pipfile":           "pip",
		"Pipfile.lock":      "pip",
		"pom.xml":           "maven",
		"build.gradle":      "gradle",
		"build.gradle.kts":  "gradle",
		"Cargo.toml":        "cargo",
		"Cargo.lock":        "cargo",
		"composer.json":     "composer",
		"composer.lock":     "composer",
		"Gemfile":           "bundler",
		"Gemfile.lock":      "bundler",
	}

	return ecosystemMap[fileName]
}

// checkDependencyForVulnerabilities checks a specific dependency for known vulnerabilities
func (dc *DependencyChecker) checkDependencyForVulnerabilities(name, version, ecosystem string) []interfaces.Vulnerability {
	var vulnerabilities []interfaces.Vulnerability

	// Check against the vulnerability database
	if vulns, exists := dc.vulnerabilityDatabase[name]; exists {
		for _, vuln := range vulns {
			// Simple version matching (in a real implementation, use semantic versioning)
			if dc.isVersionVulnerable(version, vuln.Version) {
				vulnerabilities = append(vulnerabilities, vuln)
			}
		}
	}

	return vulnerabilities
}

// scanForVulnerabilityPatterns scans content for vulnerability patterns
func (dc *DependencyChecker) scanForVulnerabilityPatterns(content, ecosystem string) []interfaces.Vulnerability {
	var vulnerabilities []interfaces.Vulnerability

	// Define vulnerability patterns for different ecosystems
	patterns := dc.getVulnerabilityPatterns(ecosystem)

	for pattern, vuln := range patterns {
		matched, err := regexp.MatchString(pattern, content)
		if err != nil {
			continue
		}
		if matched {
			vulnerabilities = append(vulnerabilities, vuln)
		}
	}

	return vulnerabilities
}

// isVersionVulnerable checks if a version is vulnerable based on version constraint
func (dc *DependencyChecker) isVersionVulnerable(version, constraint string) bool {
	// This is a simplified version check
	// In a real implementation, you would use proper semantic versioning libraries

	if constraint == "" || version == "" {
		return false
	}

	// Handle simple constraints like "<4.17.21"
	if strings.HasPrefix(constraint, "<") {
		constraintVersion := strings.TrimPrefix(constraint, "<")
		return dc.compareVersions(version, constraintVersion) < 0
	}

	// Handle range constraints like ">=4.0.0 <4.17.21"
	if strings.Contains(constraint, " ") {
		// For now, just check if version contains the pattern
		return strings.Contains(version, strings.Fields(constraint)[0])
	}

	// Exact match
	return version == constraint
}

// compareVersions compares two version strings (simplified)
func (dc *DependencyChecker) compareVersions(v1, v2 string) int {
	// This is a very simplified version comparison
	// In a real implementation, use a proper semver library
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	}
	return 1
}

// initializeVulnerabilityDatabase initializes the vulnerability database with known vulnerabilities
func (dc *DependencyChecker) initializeVulnerabilityDatabase() {
	// Initialize with some common vulnerabilities
	// In a real implementation, this would be loaded from external databases

	dc.vulnerabilityDatabase["lodash"] = []interfaces.Vulnerability{
		{
			ID:          "CVE-2019-10744",
			Severity:    "high",
			Title:       "Prototype Pollution in lodash",
			Description: "Versions of lodash before 4.17.12 are vulnerable to Prototype Pollution",
			Package:     "lodash",
			Version:     "<4.17.12",
			FixedIn:     "4.17.12",
		},
		{
			ID:          "CVE-2021-23337",
			Severity:    "high",
			Title:       "Command Injection in lodash",
			Description: "Lodash versions prior to 4.17.21 are vulnerable to Command Injection",
			Package:     "lodash",
			Version:     "<4.17.21",
			FixedIn:     "4.17.21",
		},
	}

	dc.vulnerabilityDatabase["axios"] = []interfaces.Vulnerability{
		{
			ID:          "CVE-2021-3749",
			Severity:    "medium",
			Title:       "Regular Expression Denial of Service in axios",
			Description: "axios versions prior to 0.21.2 are vulnerable to ReDoS",
			Package:     "axios",
			Version:     "<0.21.2",
			FixedIn:     "0.21.2",
		},
	}

	dc.vulnerabilityDatabase["express"] = []interfaces.Vulnerability{
		{
			ID:          "CVE-2022-24999",
			Severity:    "medium",
			Title:       "Express.js vulnerability",
			Description: "Potential security issue in Express.js versions",
			Package:     "express",
			Version:     "<4.18.0",
			FixedIn:     "4.18.0",
		},
	}

	// Add Go vulnerabilities
	dc.vulnerabilityDatabase["github.com/gin-gonic/gin"] = []interfaces.Vulnerability{
		{
			ID:          "GO-2023-1234",
			Severity:    "medium",
			Title:       "Path traversal in Gin",
			Description: "Potential path traversal vulnerability in Gin framework",
			Package:     "github.com/gin-gonic/gin",
			Version:     "<1.9.0",
			FixedIn:     "1.9.0",
		},
	}
}

// initializeEcosystemParsers initializes parsers for different ecosystems
func (dc *DependencyChecker) initializeEcosystemParsers() {
	// For now, we'll use the generic parsers
	// In a real implementation, you might have more sophisticated parsers
}

// getVulnerabilityPatterns returns vulnerability patterns for a given ecosystem
func (dc *DependencyChecker) getVulnerabilityPatterns(ecosystem string) map[string]interfaces.Vulnerability {
	patterns := make(map[string]interfaces.Vulnerability)

	switch ecosystem {
	case "npm":
		patterns[`"lodash":\s*"[^"]*4\.17\.[0-9]"`] = interfaces.Vulnerability{
			ID:          "CVE-2019-10744",
			Severity:    "high",
			Title:       "Prototype Pollution in lodash",
			Description: "Versions of lodash before 4.17.12 are vulnerable to Prototype Pollution",
			Package:     "lodash",
			Version:     "4.17.x",
			FixedIn:     "4.17.12",
		}
		patterns[`"express":\s*"[^"]*4\.[0-9]+\.[0-9]+"`] = interfaces.Vulnerability{
			ID:          "CVE-2022-24999",
			Severity:    "medium",
			Title:       "Express.js vulnerability",
			Description: "Potential security issue in Express.js versions",
			Package:     "express",
			Version:     "4.x",
			FixedIn:     "4.18.0",
		}
	case "go":
		patterns[`github\.com/gin-gonic/gin\s+v1\.[0-8]\.`] = interfaces.Vulnerability{
			ID:          "GO-2023-1234",
			Severity:    "medium",
			Title:       "Path traversal in Gin",
			Description: "Potential path traversal vulnerability in Gin framework",
			Package:     "github.com/gin-gonic/gin",
			Version:     "v1.0-v1.8",
			FixedIn:     "v1.9.0",
		}
	}

	return patterns
}

// Ecosystem-specific parsers

// parseNPMDependencies parses NPM package.json dependencies
func (dc *DependencyChecker) parseNPMDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	// Try to parse as JSON first
	var packageJSON map[string]interface{}
	if err := json.Unmarshal([]byte(content), &packageJSON); err == nil {
		// Parse dependencies section
		if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
			for name, version := range deps {
				if versionStr, ok := version.(string); ok {
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           name,
						Version:        versionStr,
						Type:           "direct",
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -6, 0), // Assume 6 months old
						SecurityIssues: 0,
						QualityScore:   75.0,
					})
				}
			}
		}

		// Parse devDependencies section
		if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
			for name, version := range devDeps {
				if versionStr, ok := version.(string); ok {
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           name,
						Version:        versionStr,
						Type:           "dev",
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -6, 0),
						SecurityIssues: 0,
						QualityScore:   75.0,
					})
				}
			}
		}
	} else {
		// Fallback to regex parsing
		dependencies = dc.parseNPMDependenciesRegex(content)
	}

	return dependencies
}

// parseNPMDependenciesRegex parses NPM dependencies using regex (fallback)
func (dc *DependencyChecker) parseNPMDependenciesRegex(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	// Simple regex-based parsing for dependencies
	depRegex := regexp.MustCompile(`"([^"]+)":\s*"([^"]+)"`)
	matches := depRegex.FindAllStringSubmatch(content, -1)

	inDepsSection := false
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.Contains(line, `"dependencies"`) || strings.Contains(line, `"devDependencies"`) {
			inDepsSection = true
			continue
		}
		if inDepsSection && strings.Contains(line, "}") {
			inDepsSection = false
			continue
		}
		if inDepsSection {
			for _, match := range matches {
				if len(match) >= 3 && strings.Contains(line, match[0]) {
					depType := "direct"
					if strings.Contains(content[:strings.Index(content, line)], `"devDependencies"`) {
						depType = "dev"
					}

					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           match[1],
						Version:        match[2],
						Type:           depType,
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -6, 0),
						SecurityIssues: 0,
						QualityScore:   75.0,
					})
				}
			}
		}
	}

	return dependencies
}

// parseGoDependencies parses Go mod dependencies
func (dc *DependencyChecker) parseGoDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	lines := strings.Split(content, "\n")
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		if strings.HasPrefix(line, "require ") || inRequireBlock {
			// Parse require line
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				var name, version string
				if strings.HasPrefix(line, "require ") {
					name = parts[1]
					if len(parts) >= 3 {
						version = parts[2]
					}
				} else {
					name = parts[0]
					if len(parts) >= 2 {
						version = parts[1]
					}
				}

				if name != "" {
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           name,
						Version:        version,
						Type:           "direct",
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -3, 0), // Assume 3 months old
						SecurityIssues: 0,
						QualityScore:   80.0,
					})
				}
			}
		}
	}

	return dependencies
}

// parsePipDependencies parses Python requirements.txt dependencies
func (dc *DependencyChecker) parsePipDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse requirement line (package==version or package>=version)
		parts := regexp.MustCompile(`[><=!]+`).Split(line, 2)
		if len(parts) >= 1 {
			name := strings.TrimSpace(parts[0])
			version := "unknown"
			if len(parts) >= 2 {
				version = strings.TrimSpace(parts[1])
			}

			dependencies = append(dependencies, interfaces.DependencyInfo{
				Name:           name,
				Version:        version,
				Type:           "direct",
				License:        "unknown",
				LastUpdated:    time.Now().AddDate(0, -4, 0), // Assume 4 months old
				SecurityIssues: 0,
				QualityScore:   70.0,
			})
		}
	}

	return dependencies
}

// parseMavenDependencies parses Maven pom.xml dependencies (simplified)
func (dc *DependencyChecker) parseMavenDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	// Simple regex-based parsing for Maven dependencies
	depRegex := regexp.MustCompile(`<groupId>([^<]+)</groupId>\s*<artifactId>([^<]+)</artifactId>\s*<version>([^<]+)</version>`)
	matches := depRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			name := fmt.Sprintf("%s:%s", match[1], match[2])
			version := match[3]

			dependencies = append(dependencies, interfaces.DependencyInfo{
				Name:           name,
				Version:        version,
				Type:           "direct",
				License:        "unknown",
				LastUpdated:    time.Now().AddDate(0, -5, 0), // Assume 5 months old
				SecurityIssues: 0,
				QualityScore:   75.0,
			})
		}
	}

	return dependencies
}

// parseGradleDependencies parses Gradle build.gradle dependencies (simplified)
func (dc *DependencyChecker) parseGradleDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	// Simple regex-based parsing for Gradle dependencies
	depRegex := regexp.MustCompile(`(?:implementation|compile|api|testImplementation)\s+['"]([^'"]+)['"]`)
	matches := depRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			depString := match[1]
			parts := strings.Split(depString, ":")

			name := depString
			version := "unknown"
			if len(parts) >= 3 {
				name = fmt.Sprintf("%s:%s", parts[0], parts[1])
				version = parts[2]
			}

			dependencies = append(dependencies, interfaces.DependencyInfo{
				Name:           name,
				Version:        version,
				Type:           "direct",
				License:        "unknown",
				LastUpdated:    time.Now().AddDate(0, -5, 0), // Assume 5 months old
				SecurityIssues: 0,
				QualityScore:   75.0,
			})
		}
	}

	return dependencies
}

// parseCargoDependencies parses Rust Cargo.toml dependencies
func (dc *DependencyChecker) parseCargoDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	lines := strings.Split(content, "\n")
	inDepsSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "[dependencies]" {
			inDepsSection = true
			continue
		}
		if inDepsSection && strings.HasPrefix(line, "[") {
			inDepsSection = false
			continue
		}

		if inDepsSection && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) >= 2 {
				name := strings.TrimSpace(parts[0])
				version := strings.Trim(strings.TrimSpace(parts[1]), `"`)

				dependencies = append(dependencies, interfaces.DependencyInfo{
					Name:           name,
					Version:        version,
					Type:           "direct",
					License:        "unknown",
					LastUpdated:    time.Now().AddDate(0, -4, 0), // Assume 4 months old
					SecurityIssues: 0,
					QualityScore:   80.0,
				})
			}
		}
	}

	return dependencies
}

// parseComposerDependencies parses PHP composer.json dependencies
func (dc *DependencyChecker) parseComposerDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	// Try to parse as JSON
	var composerJSON map[string]interface{}
	if err := json.Unmarshal([]byte(content), &composerJSON); err == nil {
		// Parse require section
		if require, ok := composerJSON["require"].(map[string]interface{}); ok {
			for name, version := range require {
				if versionStr, ok := version.(string); ok {
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           name,
						Version:        versionStr,
						Type:           "direct",
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -6, 0),
						SecurityIssues: 0,
						QualityScore:   70.0,
					})
				}
			}
		}

		// Parse require-dev section
		if requireDev, ok := composerJSON["require-dev"].(map[string]interface{}); ok {
			for name, version := range requireDev {
				if versionStr, ok := version.(string); ok {
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           name,
						Version:        versionStr,
						Type:           "dev",
						License:        "unknown",
						LastUpdated:    time.Now().AddDate(0, -6, 0),
						SecurityIssues: 0,
						QualityScore:   70.0,
					})
				}
			}
		}
	}

	return dependencies
}

// parseBundlerDependencies parses Ruby Gemfile dependencies
func (dc *DependencyChecker) parseBundlerDependencies(content string) []interfaces.DependencyInfo {
	var dependencies []interfaces.DependencyInfo

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "gem ") {
			// Parse gem line: gem 'name', 'version' or gem 'name'
			// Remove 'gem ' prefix
			gemLine := strings.TrimPrefix(line, "gem ")

			// Split by comma to handle: gem 'rails', '~> 7.0.0'
			parts := strings.Split(gemLine, ",")
			if len(parts) >= 1 {
				name := strings.Trim(strings.TrimSpace(parts[0]), `'"`)
				version := "unknown"
				if len(parts) >= 2 {
					version = strings.Trim(strings.TrimSpace(parts[1]), `'",`)
				}

				dependencies = append(dependencies, interfaces.DependencyInfo{
					Name:           name,
					Version:        version,
					Type:           "direct",
					License:        "unknown",
					LastUpdated:    time.Now().AddDate(0, -5, 0), // Assume 5 months old
					SecurityIssues: 0,
					QualityScore:   75.0,
				})
			}
		}
	}

	return dependencies
}
