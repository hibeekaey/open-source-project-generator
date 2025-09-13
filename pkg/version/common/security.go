package common

import "fmt"

// VulnerabilityDB defines the interface for vulnerability databases
type VulnerabilityDB interface {
	CheckVulnerabilities(packageName, version string) ([]SecurityIssue, error)
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	ID          string
	Severity    string
	Title       string
	Description string
	FixedIn     string
	CVSS        float64
}

// CheckSecurity performs security checks on a package version
func CheckSecurity(packageName, version string, vulnerabilityDB VulnerabilityDB) ([]SecurityIssue, error) {
	if vulnerabilityDB == nil {
		return nil, fmt.Errorf("vulnerability database is required")
	}

	issues, err := vulnerabilityDB.CheckVulnerabilities(packageName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to check vulnerabilities for %s@%s: %w", packageName, version, err)
	}

	return issues, nil
}

// FilterSecurityIssues filters security issues by severity threshold
func FilterSecurityIssues(issues []SecurityIssue, minSeverity string) []SecurityIssue {
	severityLevels := map[string]int{
		"low":      1,
		"moderate": 2,
		"high":     3,
		"critical": 4,
	}

	minLevel, exists := severityLevels[minSeverity]
	if !exists {
		return issues // Return all if invalid severity
	}

	var filtered []SecurityIssue
	for _, issue := range issues {
		if level, exists := severityLevels[issue.Severity]; exists && level >= minLevel {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// HasCriticalIssues checks if there are any critical security issues
func HasCriticalIssues(issues []SecurityIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "critical" {
			return true
		}
	}
	return false
}
