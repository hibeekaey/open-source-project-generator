// Package security provides security auditing functionality for the
// Open Source Project Generator.
//
// Security Note: This package contains security audit functionality that legitimately needs
// to read files for security analysis. The G304 warnings from gosec are false positives
// in this context as file reading is the core functionality of a security audit tool.
package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// SecurityScanner implements security auditing functionality.
type SecurityScanner struct {
	// secretDetector handles secret detection
	secretDetector *SecretDetector
	// dependencyChecker handles dependency vulnerability scanning
	dependencyChecker *DependencyChecker
}

// NewSecurityScanner creates a new security scanner instance.
func NewSecurityScanner() *SecurityScanner {
	return &SecurityScanner{
		secretDetector:    NewSecretDetector(),
		dependencyChecker: NewDependencyChecker(),
	}
}

// AuditSecurity performs security auditing on a project
func (s *SecurityScanner) AuditSecurity(path string) (*interfaces.SecurityAuditResult, error) {
	if err := s.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.SecurityAuditResult{
		Score:            100.0,
		Vulnerabilities:  []interfaces.Vulnerability{},
		PolicyViolations: []interfaces.PolicyViolation{},
		Recommendations:  []string{},
	}

	// Scan for vulnerabilities
	vulnReport, err := s.ScanVulnerabilities(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Vulnerability scan failed.",
			"Unable to analyze project for security vulnerabilities")
	}
	result.Vulnerabilities = vulnReport.Vulnerabilities

	// Check security policies
	policyResult, err := s.CheckSecurityPolicies(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Security policy check failed.",
			"Unable to validate project security policies")
	}
	result.PolicyViolations = policyResult.Violations

	// Detect secrets using the dedicated secret detector
	secretResult, err := s.secretDetector.DetectSecrets(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Secret detection failed.",
			"Unable to scan project for exposed secrets")
	}

	// Calculate security score based on findings
	result.Score = s.calculateSecurityScore(result, secretResult)

	// Generate security recommendations
	result.Recommendations = s.generateSecurityRecommendations(result, secretResult)

	return result, nil
}

// ScanVulnerabilities scans for security vulnerabilities
func (s *SecurityScanner) ScanVulnerabilities(path string) (*interfaces.VulnerabilityReport, error) {
	report := &interfaces.VulnerabilityReport{
		ScanTime:        time.Now(),
		Vulnerabilities: []interfaces.Vulnerability{},
		Summary: interfaces.VulnerabilitySummary{
			Total:    0,
			Critical: 0,
			High:     0,
			Medium:   0,
			Low:      0,
			Fixed:    0,
			Ignored:  0,
		},
		Recommendations: []string{},
	}

	// Scan dependency files for known vulnerabilities using the dedicated dependency checker
	vulnerabilities, err := s.dependencyChecker.ScanDependencyVulnerabilities(path)
	if err != nil {
		return nil, fmt.Errorf("dependency vulnerability scan failed: %w", err)
	}
	report.Vulnerabilities = append(report.Vulnerabilities, vulnerabilities...)

	// Update summary
	for _, vuln := range report.Vulnerabilities {
		report.Summary.Total++
		switch vuln.Severity {
		case "critical":
			report.Summary.Critical++
		case "high":
			report.Summary.High++
		case "medium":
			report.Summary.Medium++
		case "low":
			report.Summary.Low++
		}
	}

	// Generate recommendations
	if report.Summary.Critical > 0 {
		report.Recommendations = append(report.Recommendations, "Immediately address critical vulnerabilities")
	}
	if report.Summary.High > 0 {
		report.Recommendations = append(report.Recommendations, "Address high severity vulnerabilities as soon as possible")
	}
	if report.Summary.Total > 0 {
		report.Recommendations = append(report.Recommendations, "Consider using automated dependency scanning tools")
	}

	return report, nil
}

// CheckSecurityPolicies checks security policy compliance
func (s *SecurityScanner) CheckSecurityPolicies(path string) (*interfaces.PolicyComplianceResult, error) {
	result := &interfaces.PolicyComplianceResult{
		Compliant:  true,
		Policies:   []interfaces.PolicyCheck{},
		Violations: []interfaces.PolicyViolation{},
		Score:      100.0,
		Summary: interfaces.PolicyComplianceSummary{
			TotalPolicies:      0,
			CompliantPolicies:  0,
			Violations:         0,
			CriticalViolations: 0,
		},
	}

	// Define security policies to check
	policies := []interfaces.PolicyCheck{
		{
			ID:          "SEC-001",
			Name:        "No hardcoded secrets",
			Description: "Source code should not contain hardcoded secrets",
			Category:    "secrets",
			Severity:    "critical",
			Compliant:   true,
		},
		{
			ID:          "SEC-002",
			Name:        "Secure dependencies",
			Description: "Dependencies should not have known vulnerabilities",
			Category:    "dependencies",
			Severity:    "high",
			Compliant:   true,
		},
		{
			ID:          "SEC-003",
			Name:        "Secure configuration",
			Description: "Configuration files should follow security best practices",
			Category:    "configuration",
			Severity:    "medium",
			Compliant:   true,
		},
	}

	// Check each policy
	for _, policy := range policies {
		result.Summary.TotalPolicies++

		// Check policy compliance based on type
		switch policy.ID {
		case "SEC-001":
			violations, err := s.checkHardcodedSecrets(path)
			if err != nil {
				return nil, fmt.Errorf("failed to check hardcoded secrets: %w", err)
			}
			if len(violations) > 0 {
				policy.Compliant = false
				result.Compliant = false
				result.Violations = append(result.Violations, violations...)
				result.Summary.Violations += len(violations)
				for _, v := range violations {
					if v.Severity == "critical" {
						result.Summary.CriticalViolations++
					}
				}
			}
		case "SEC-002":
			violations, err := s.dependencyChecker.CheckDependencyVulnerabilities(path)
			if err != nil {
				return nil, fmt.Errorf("failed to check dependency vulnerabilities: %w", err)
			}
			if len(violations) > 0 {
				policy.Compliant = false
				result.Compliant = false
				result.Violations = append(result.Violations, violations...)
				result.Summary.Violations += len(violations)
			}
		case "SEC-003":
			violations, err := s.checkSecureConfiguration(path)
			if err != nil {
				return nil, fmt.Errorf("failed to check secure configuration: %w", err)
			}
			if len(violations) > 0 {
				policy.Compliant = false
				result.Compliant = false
				result.Violations = append(result.Violations, violations...)
				result.Summary.Violations += len(violations)
			}
		}

		if policy.Compliant {
			result.Summary.CompliantPolicies++
		}
		result.Policies = append(result.Policies, policy)
	}

	// Calculate compliance score
	if result.Summary.TotalPolicies > 0 {
		result.Score = float64(result.Summary.CompliantPolicies) / float64(result.Summary.TotalPolicies) * 100
	}

	return result, nil
}

// DetectSecrets detects secrets in the project using the dedicated SecretDetector
func (s *SecurityScanner) DetectSecrets(path string) (*interfaces.SecretScanResult, error) {
	return s.secretDetector.DetectSecrets(path)
}

// AnalyzeDependencies analyzes project dependencies for security issues
func (s *SecurityScanner) AnalyzeDependencies(path string) (*interfaces.DependencyAnalysisResult, error) {
	return s.dependencyChecker.AnalyzeDependencies(path)
}

// Helper methods

// projectExists checks if the project path exists and is a directory
func (s *SecurityScanner) projectExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("project path does not exist: %s", path)
		}
		return fmt.Errorf("🚫 %s %s",
			"Unable to access project path.",
			"Check if the directory exists and has proper permissions")
	}

	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", path)
	}

	return nil
}

// checkHardcodedSecrets checks for hardcoded secrets in source files
func (s *SecurityScanner) checkHardcodedSecrets(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	// Walk through source files
	walkErr := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || s.shouldSkipFile(filePath) {
			return nil
		}

		// Scan file for secrets using the secret detector
		secrets, err := s.secretDetector.scanFileForSecrets(filePath)
		if err != nil {
			return nil // Continue on error
		}

		// Convert secrets to policy violations
		for _, secret := range secrets {
			if secret.Confidence >= 0.7 { // Only high confidence secrets
				violations = append(violations, interfaces.PolicyViolation{
					Policy:      "SEC-001",
					Severity:    "critical",
					Description: fmt.Sprintf("Potential %s found in source code", secret.Type),
					File:        secret.File,
					Line:        secret.Line,
				})
			}
		}

		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("failed to scan for hardcoded secrets: %w", walkErr)
	}

	return violations, nil
}

// checkSecureConfiguration checks for secure configuration practices
func (s *SecurityScanner) checkSecureConfiguration(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	// Check common configuration files for security issues
	configFiles := map[string][]string{
		"docker-compose.yml": {
			`privileged:\s*true`,
			`--privileged`,
			`security_opt:\s*-\s*seccomp:unconfined`,
		},
		"Dockerfile": {
			`USER\s+root`,
			`--privileged`,
			`COPY\s+.*\s+/`,
		},
		".env": {
			`DEBUG\s*=\s*true`,
			`ENVIRONMENT\s*=\s*development`,
		},
	}

	for configFile, patterns := range configFiles {
		filePath := filepath.Join(path, configFile)
		if _, err := os.Stat(filePath); err == nil {
			content, err := os.ReadFile(filePath) // #nosec G304 - This is an audit tool that needs to read files
			if err != nil {
				continue
			}

			contentStr := string(content)
			lines := strings.Split(contentStr, "\n")

			for _, pattern := range patterns {
				regex, err := regexp.Compile(pattern)
				if err != nil {
					continue
				}

				for lineNum, line := range lines {
					if regex.MatchString(line) {
						violations = append(violations, interfaces.PolicyViolation{
							Policy:      "SEC-003",
							Severity:    "medium",
							Description: fmt.Sprintf("Insecure configuration pattern found: %s", pattern),
							File:        filePath,
							Line:        lineNum + 1,
						})
					}
				}
			}
		}
	}

	return violations, nil
}

// shouldSkipFile determines if a file should be skipped during scanning
func (s *SecurityScanner) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".class", ".jar", ".war", ".ear",
		".o", ".obj", ".lib", ".a",
	}

	// Skip directories that typically don't contain source code
	skipDirs := []string{
		"node_modules", "vendor", "target", "build", "dist",
		".git", ".svn", ".hg", ".bzr",
		"__pycache__", ".pytest_cache",
		"coverage", ".coverage",
		"logs", "log",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Check if any part of the path contains skip directories
	pathParts := strings.Split(filePath, string(filepath.Separator))
	for _, part := range pathParts {
		for _, skipDir := range skipDirs {
			if part == skipDir {
				return true
			}
		}
	}

	return false
}

// maskSecret masks a secret for safe display
func (s *SecurityScanner) maskSecret(secret string) string {
	if len(secret) <= 4 {
		return strings.Repeat("*", len(secret))
	}
	return secret[:2] + strings.Repeat("*", len(secret)-4) + secret[len(secret)-2:]
}

// calculateSecurityScore calculates the security score based on findings
func (s *SecurityScanner) calculateSecurityScore(result *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) float64 {
	score := 100.0

	// Deduct points for vulnerabilities
	for _, vuln := range result.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			score -= 25.0
		case "high":
			score -= 15.0
		case "medium":
			score -= 10.0
		case "low":
			score -= 5.0
		}
	}

	// Deduct points for policy violations
	for _, violation := range result.PolicyViolations {
		switch violation.Severity {
		case "critical":
			score -= 20.0
		case "high":
			score -= 10.0
		case "medium":
			score -= 5.0
		case "low":
			score -= 2.0
		}
	}

	// Deduct points for secrets
	if secretResult != nil {
		score -= float64(secretResult.Summary.HighConfidence) * 15.0
		score -= float64(secretResult.Summary.MediumConfidence) * 8.0
		score -= float64(secretResult.Summary.LowConfidence) * 3.0
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// generateSecurityRecommendations generates security recommendations
func (s *SecurityScanner) generateSecurityRecommendations(result *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) []string {
	var recommendations []string

	// Recommendations based on vulnerabilities
	if len(result.Vulnerabilities) > 0 {
		recommendations = append(recommendations, "Update vulnerable dependencies to their latest secure versions")
		recommendations = append(recommendations, "Implement automated dependency scanning in your CI/CD pipeline")
	}

	// Recommendations based on policy violations
	if len(result.PolicyViolations) > 0 {
		recommendations = append(recommendations, "Review and fix security policy violations")
		recommendations = append(recommendations, "Implement security linting tools in your development workflow")
	}

	// Recommendations based on secrets
	if secretResult != nil && secretResult.Summary.TotalSecrets > 0 {
		recommendations = append(recommendations, "Remove hardcoded secrets from source code")
		recommendations = append(recommendations, "Use environment variables or secure secret management systems")
		recommendations = append(recommendations, "Implement pre-commit hooks to prevent secret commits")
	}

	// General security recommendations
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Continue following security best practices")
		recommendations = append(recommendations, "Regularly update dependencies and scan for vulnerabilities")
	}

	return recommendations
}
