// Package interfaces defines the core contracts and interfaces for the
// Open Source Project Generator components.
package interfaces

import "time"

// AuditEngine defines the contract for project auditing operations.
//
// This interface abstracts auditing functionality to enable different
// audit implementations and security scanners.
type AuditEngine interface {
	// Security auditing
	AuditSecurity(path string) (*SecurityAuditResult, error)
	ScanVulnerabilities(path string) (*VulnerabilityReport, error)
	CheckSecurityPolicies(path string) (*PolicyComplianceResult, error)

	// Quality auditing
	AuditCodeQuality(path string) (*QualityAuditResult, error)
	CheckBestPractices(path string) (*BestPracticesResult, error)
	AnalyzeDependencies(path string) (*DependencyAnalysisResult, error)

	// License auditing
	AuditLicenses(path string) (*LicenseAuditResult, error)
	CheckLicenseCompatibility(path string) (*LicenseCompatibilityResult, error)

	// Performance auditing
	AuditPerformance(path string) (*PerformanceAuditResult, error)
	AnalyzeBundleSize(path string) (*BundleAnalysisResult, error)
}

// VulnerabilityReport contains vulnerability scan results
type VulnerabilityReport struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         VulnSummary     `json:"summary"`
	ScanTime        time.Time       `json:"scan_time"`
}

// VulnSummary contains vulnerability summary statistics
type VulnSummary struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

// PolicyComplianceResult contains security policy compliance results
type PolicyComplianceResult struct {
	Compliant  bool              `json:"compliant"`
	Violations []PolicyViolation `json:"violations"`
	Score      float64           `json:"score"`
	CheckedAt  time.Time         `json:"checked_at"`
}

// BestPracticesResult contains best practices analysis results
type BestPracticesResult struct {
	Score           float64     `json:"score"`
	Recommendations []string    `json:"recommendations"`
	Issues          []CodeSmell `json:"issues"`
	CheckedAt       time.Time   `json:"checked_at"`
}

// DependencyAnalysisResult contains dependency analysis results
type DependencyAnalysisResult struct {
	TotalDependencies int               `json:"total_dependencies"`
	OutdatedPackages  []OutdatedPackage `json:"outdated_packages"`
	SecurityIssues    []DependencyIssue `json:"security_issues"`
	LicenseIssues     []LicenseIssue    `json:"license_issues"`
	Recommendations   []string          `json:"recommendations"`
	AnalyzedAt        time.Time         `json:"analyzed_at"`
}

// OutdatedPackage represents an outdated dependency
type OutdatedPackage struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Severity       string `json:"severity"`
}

// DependencyIssue represents a security issue in a dependency
type DependencyIssue struct {
	Package  string `json:"package"`
	Version  string `json:"version"`
	Issue    string `json:"issue"`
	Severity string `json:"severity"`
	FixedIn  string `json:"fixed_in,omitempty"`
}

// LicenseIssue represents a license compatibility issue
type LicenseIssue struct {
	Package  string `json:"package"`
	License  string `json:"license"`
	Issue    string `json:"issue"`
	Severity string `json:"severity"`
}

// LicenseCompatibilityResult contains license compatibility results
type LicenseCompatibilityResult struct {
	Compatible      bool          `json:"compatible"`
	ProjectLicense  string        `json:"project_license"`
	ConflictingDeps []LicenseInfo `json:"conflicting_deps"`
	Recommendations []string      `json:"recommendations"`
	CheckedAt       time.Time     `json:"checked_at"`
}

// BundleAnalysisResult contains bundle size analysis results
type BundleAnalysisResult struct {
	TotalSize       int64       `json:"total_size"`
	CompressedSize  int64       `json:"compressed_size"`
	LargestAssets   []AssetInfo `json:"largest_assets"`
	Recommendations []string    `json:"recommendations"`
	AnalyzedAt      time.Time   `json:"analyzed_at"`
}

// AssetInfo contains information about a bundle asset
type AssetInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}
