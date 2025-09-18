// Package interfaces defines audit interfaces.
//
// This file contains interfaces for comprehensive project auditing
// including security, quality, license, and performance auditing.
package interfaces

import (
	"time"
)

// AuditEngine defines the interface for comprehensive project auditing operations.
//
// This interface provides enterprise-grade auditing capabilities including:
//   - Security vulnerability scanning and policy compliance
//   - Code quality analysis and best practices checking
//   - License compliance and compatibility checking
//   - Performance analysis and optimization recommendations
//   - Comprehensive reporting and scoring
type AuditEngine interface {
	// Security auditing
	AuditSecurity(path string) (*SecurityAuditResult, error)
	ScanVulnerabilities(path string) (*VulnerabilityReport, error)
	CheckSecurityPolicies(path string) (*PolicyComplianceResult, error)
	DetectSecrets(path string) (*SecretScanResult, error)

	// Quality auditing
	AuditCodeQuality(path string) (*QualityAuditResult, error)
	CheckBestPractices(path string) (*BestPracticesResult, error)
	AnalyzeDependencies(path string) (*DependencyAnalysisResult, error)
	MeasureComplexity(path string) (*ComplexityAnalysisResult, error)

	// License auditing
	AuditLicenses(path string) (*LicenseAuditResult, error)
	CheckLicenseCompatibility(path string) (*LicenseCompatibilityResult, error)
	ScanLicenseViolations(path string) (*LicenseViolationResult, error)

	// Performance auditing
	AuditPerformance(path string) (*PerformanceAuditResult, error)
	AnalyzeBundleSize(path string) (*BundleAnalysisResult, error)
	CheckPerformanceMetrics(path string) (*PerformanceMetricsResult, error)

	// Comprehensive auditing
	AuditProject(path string, options *AuditOptions) (*AuditResult, error)
	GenerateAuditReport(result *AuditResult, format string) ([]byte, error)
	GetAuditSummary(results []*AuditResult) (*AuditSummary, error)

	// Audit configuration
	SetAuditRules(rules []AuditRule) error
	GetAuditRules() []AuditRule
	AddAuditRule(rule AuditRule) error
	RemoveAuditRule(ruleID string) error
}

// VulnerabilityReport contains detailed vulnerability scan results
type VulnerabilityReport struct {
	ScanTime        time.Time            `json:"scan_time"`
	Vulnerabilities []Vulnerability      `json:"vulnerabilities"`
	Summary         VulnerabilitySummary `json:"summary"`
	Recommendations []string             `json:"recommendations"`
}

// VulnerabilitySummary contains vulnerability statistics
type VulnerabilitySummary struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Fixed    int `json:"fixed"`
	Ignored  int `json:"ignored"`
}

// PolicyComplianceResult contains security policy compliance results
type PolicyComplianceResult struct {
	Compliant  bool                    `json:"compliant"`
	Policies   []PolicyCheck           `json:"policies"`
	Violations []PolicyViolation       `json:"violations"`
	Score      float64                 `json:"score"`
	Summary    PolicyComplianceSummary `json:"summary"`
}

// PolicyCheck represents a security policy check
type PolicyCheck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Compliant   bool   `json:"compliant"`
	Details     string `json:"details"`
}

// PolicyComplianceSummary contains policy compliance statistics
type PolicyComplianceSummary struct {
	TotalPolicies      int `json:"total_policies"`
	CompliantPolicies  int `json:"compliant_policies"`
	Violations         int `json:"violations"`
	CriticalViolations int `json:"critical_violations"`
}

// SecretScanResult contains secret detection results
type SecretScanResult struct {
	ScanTime time.Time         `json:"scan_time"`
	Secrets  []SecretDetection `json:"secrets"`
	Summary  SecretScanSummary `json:"summary"`
}

// SecretScanSummary contains secret scan statistics
type SecretScanSummary struct {
	TotalSecrets     int `json:"total_secrets"`
	HighConfidence   int `json:"high_confidence"`
	MediumConfidence int `json:"medium_confidence"`
	LowConfidence    int `json:"low_confidence"`
	FilesScanned     int `json:"files_scanned"`
}

// BestPracticesResult contains best practices analysis results
type BestPracticesResult struct {
	Score      float64                 `json:"score"`
	Practices  []BestPracticeCheck     `json:"practices"`
	Violations []BestPracticeViolation `json:"violations"`
	Summary    BestPracticesSummary    `json:"summary"`
}

// BestPracticeCheck represents a best practice check
type BestPracticeCheck struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Technology  string  `json:"technology"`
	Compliant   bool    `json:"compliant"`
	Score       float64 `json:"score"`
	Details     string  `json:"details"`
}

// BestPracticeViolation represents a best practice violation
type BestPracticeViolation struct {
	Practice    string `json:"practice"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

// BestPracticesSummary contains best practices statistics
type BestPracticesSummary struct {
	TotalPractices     int     `json:"total_practices"`
	CompliantPractices int     `json:"compliant_practices"`
	Violations         int     `json:"violations"`
	OverallScore       float64 `json:"overall_score"`
}

// DependencyAnalysisResult contains dependency analysis results
type DependencyAnalysisResult struct {
	Dependencies    []DependencyInfo          `json:"dependencies"`
	Vulnerabilities []DependencyVulnerability `json:"vulnerabilities"`
	Licenses        []DependencyLicense       `json:"licenses"`
	Outdated        []OutdatedDependency      `json:"outdated"`
	Summary         DependencyAnalysisSummary `json:"summary"`
}

// DependencyInfo contains information about a dependency
type DependencyInfo struct {
	Name           string    `json:"name"`
	Version        string    `json:"version"`
	Type           string    `json:"type"` // direct, transitive
	License        string    `json:"license"`
	Repository     string    `json:"repository"`
	Homepage       string    `json:"homepage"`
	LastUpdated    time.Time `json:"last_updated"`
	Maintainers    []string  `json:"maintainers"`
	SecurityIssues int       `json:"security_issues"`
	QualityScore   float64   `json:"quality_score"`
}

// DependencyLicense contains license information for a dependency
type DependencyLicense struct {
	Dependency string `json:"dependency"`
	License    string `json:"license"`
	SPDXID     string `json:"spdx_id"`
	Compatible bool   `json:"compatible"`
	Risk       string `json:"risk"` // low, medium, high
}

// DependencyAnalysisSummary contains dependency analysis statistics
type DependencyAnalysisSummary struct {
	TotalDependencies  int     `json:"total_dependencies"`
	DirectDependencies int     `json:"direct_dependencies"`
	Vulnerabilities    int     `json:"vulnerabilities"`
	OutdatedCount      int     `json:"outdated_count"`
	LicenseIssues      int     `json:"license_issues"`
	AverageAge         float64 `json:"average_age_days"`
}

// ComplexityAnalysisResult contains code complexity analysis results
type ComplexityAnalysisResult struct {
	Files     []FileComplexity          `json:"files"`
	Functions []FunctionComplexity      `json:"functions"`
	Summary   ComplexityAnalysisSummary `json:"summary"`
}

// FileComplexity contains complexity metrics for a file
type FileComplexity struct {
	Path                 string  `json:"path"`
	Lines                int     `json:"lines"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	CognitiveComplexity  int     `json:"cognitive_complexity"`
	Maintainability      float64 `json:"maintainability"`
	TechnicalDebt        string  `json:"technical_debt"`
}

// FunctionComplexity contains complexity metrics for a function
type FunctionComplexity struct {
	Name                 string `json:"name"`
	File                 string `json:"file"`
	Line                 int    `json:"line"`
	CyclomaticComplexity int    `json:"cyclomatic_complexity"`
	CognitiveComplexity  int    `json:"cognitive_complexity"`
	Parameters           int    `json:"parameters"`
	Lines                int    `json:"lines"`
}

// ComplexityAnalysisSummary contains complexity analysis statistics
type ComplexityAnalysisSummary struct {
	TotalFiles          int     `json:"total_files"`
	TotalFunctions      int     `json:"total_functions"`
	AverageComplexity   float64 `json:"average_complexity"`
	HighComplexityFiles int     `json:"high_complexity_files"`
	TechnicalDebtHours  float64 `json:"technical_debt_hours"`
}

// LicenseCompatibilityResult contains license compatibility analysis results
type LicenseCompatibilityResult struct {
	Compatible      bool                        `json:"compatible"`
	ProjectLicense  string                      `json:"project_license"`
	Dependencies    []DependencyLicense         `json:"dependencies"`
	Conflicts       []LicenseConflict           `json:"conflicts"`
	Recommendations []LicenseRecommendation     `json:"recommendations"`
	Summary         LicenseCompatibilitySummary `json:"summary"`
}

// LicenseConflict represents a license compatibility conflict
type LicenseConflict struct {
	Dependency1 string `json:"dependency1"`
	License1    string `json:"license1"`
	Dependency2 string `json:"dependency2"`
	License2    string `json:"license2"`
	Reason      string `json:"reason"`
	Severity    string `json:"severity"`
	Resolution  string `json:"resolution"`
}

// LicenseRecommendation represents a license recommendation
type LicenseRecommendation struct {
	Type        string `json:"type"` // change, add, remove
	Target      string `json:"target"`
	Current     string `json:"current"`
	Recommended string `json:"recommended"`
	Reason      string `json:"reason"`
	Impact      string `json:"impact"`
}

// LicenseCompatibilitySummary contains license compatibility statistics
type LicenseCompatibilitySummary struct {
	TotalLicenses      int    `json:"total_licenses"`
	CompatibleLicenses int    `json:"compatible_licenses"`
	Conflicts          int    `json:"conflicts"`
	RiskLevel          string `json:"risk_level"`
}

// LicenseViolationResult contains license violation scan results
type LicenseViolationResult struct {
	Violations []LicenseViolation      `json:"violations"`
	Summary    LicenseViolationSummary `json:"summary"`
}

// LicenseViolation represents a license violation
type LicenseViolation struct {
	Type       string `json:"type"` // missing, incompatible, restricted
	Dependency string `json:"dependency"`
	License    string `json:"license"`
	Violation  string `json:"violation"`
	Severity   string `json:"severity"`
	File       string `json:"file,omitempty"`
	Line       int    `json:"line,omitempty"`
	Resolution string `json:"resolution"`
}

// LicenseViolationSummary contains license violation statistics
type LicenseViolationSummary struct {
	TotalViolations    int `json:"total_violations"`
	CriticalViolations int `json:"critical_violations"`
	HighViolations     int `json:"high_violations"`
	MediumViolations   int `json:"medium_violations"`
	LowViolations      int `json:"low_violations"`
}

// BundleAnalysisResult contains bundle size analysis results
type BundleAnalysisResult struct {
	TotalSize   int64                 `json:"total_size"`
	GzippedSize int64                 `json:"gzipped_size"`
	Assets      []BundleAsset         `json:"assets"`
	Chunks      []BundleChunk         `json:"chunks"`
	Summary     BundleAnalysisSummary `json:"summary"`
}

// BundleAsset represents an asset in the bundle
type BundleAsset struct {
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	GzippedSize int64   `json:"gzipped_size"`
	Type        string  `json:"type"`
	Percentage  float64 `json:"percentage"`
}

// BundleChunk represents a chunk in the bundle
type BundleChunk struct {
	Name    string   `json:"name"`
	Size    int64    `json:"size"`
	Files   []string `json:"files"`
	Entry   bool     `json:"entry"`
	Initial bool     `json:"initial"`
}

// BundleAnalysisSummary contains bundle analysis statistics
type BundleAnalysisSummary struct {
	TotalAssets      int      `json:"total_assets"`
	TotalChunks      int      `json:"total_chunks"`
	CompressionRatio float64  `json:"compression_ratio"`
	LargestAsset     string   `json:"largest_asset"`
	Recommendations  []string `json:"recommendations"`
}

// PerformanceMetricsResult contains performance metrics analysis results
type PerformanceMetricsResult struct {
	LoadTime              time.Duration             `json:"load_time"`
	FirstPaint            time.Duration             `json:"first_paint"`
	FirstContentful       time.Duration             `json:"first_contentful"`
	LargestContentful     time.Duration             `json:"largest_contentful"`
	TimeToInteractive     time.Duration             `json:"time_to_interactive"`
	CumulativeLayoutShift float64                   `json:"cumulative_layout_shift"`
	Issues                []PerformanceIssue        `json:"issues"`
	Summary               PerformanceMetricsSummary `json:"summary"`
}

// PerformanceMetricsSummary contains performance metrics statistics
type PerformanceMetricsSummary struct {
	OverallScore     float64  `json:"overall_score"`
	PerformanceGrade string   `json:"performance_grade"`
	IssueCount       int      `json:"issue_count"`
	Recommendations  []string `json:"recommendations"`
}

// AuditSummary contains summary information across multiple audits
type AuditSummary struct {
	TotalProjects    int           `json:"total_projects"`
	AverageScore     float64       `json:"average_score"`
	SecurityScore    float64       `json:"security_score"`
	QualityScore     float64       `json:"quality_score"`
	LicenseScore     float64       `json:"license_score"`
	PerformanceScore float64       `json:"performance_score"`
	CommonIssues     []CommonIssue `json:"common_issues"`
	Trends           []AuditTrend  `json:"trends"`
}

// CommonIssue represents a commonly found issue across projects
type CommonIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Frequency   int    `json:"frequency"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
}

// AuditTrend represents trends in audit results over time
type AuditTrend struct {
	Metric    string    `json:"metric"`
	Direction string    `json:"direction"` // improving, declining, stable
	Change    float64   `json:"change"`
	Period    string    `json:"period"`
	Timestamp time.Time `json:"timestamp"`
}

// AuditRule defines an audit rule configuration
type AuditRule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Type        string         `json:"type"` // security, quality, license, performance
	Severity    string         `json:"severity"`
	Enabled     bool           `json:"enabled"`
	Config      map[string]any `json:"config,omitempty"`
	Pattern     string         `json:"pattern,omitempty"`
	FileTypes   []string       `json:"file_types,omitempty"`
}

// AuditCategory defines categories for audit rules
const (
	AuditCategorySecurity      = "security"
	AuditCategoryQuality       = "quality"
	AuditCategoryLicense       = "license"
	AuditCategoryPerformance   = "performance"
	AuditCategoryCompliance    = "compliance"
	AuditCategoryBestPractices = "best_practices"
)

// AuditSeverity defines severity levels for audit issues
const (
	AuditSeverityCritical = "critical"
	AuditSeverityHigh     = "high"
	AuditSeverityMedium   = "medium"
	AuditSeverityLow      = "low"
	AuditSeverityInfo     = "info"
)

// SecurityAuditResult contains security audit results
type SecurityAuditResult struct {
	Score            float64           `json:"score"`
	Vulnerabilities  []Vulnerability   `json:"vulnerabilities"`
	PolicyViolations []PolicyViolation `json:"policy_violations"`
	Recommendations  []string          `json:"recommendations"`
}

// QualityAuditResult contains quality audit results
type QualityAuditResult struct {
	Score           float64       `json:"score"`
	CodeSmells      []CodeSmell   `json:"code_smells"`
	Duplications    []Duplication `json:"duplications"`
	TestCoverage    float64       `json:"test_coverage"`
	Recommendations []string      `json:"recommendations"`
}

// LicenseAuditResult contains license audit results
type LicenseAuditResult struct {
	Score           float64       `json:"score"`
	Compatible      bool          `json:"compatible"`
	Licenses        []LicenseInfo `json:"licenses"`
	Conflicts       []LicenseInfo `json:"conflicts"`
	Recommendations []string      `json:"recommendations"`
}

// PerformanceAuditResult contains performance audit results
type PerformanceAuditResult struct {
	Score           float64            `json:"score"`
	BundleSize      int64              `json:"bundle_size"`
	LoadTime        time.Duration      `json:"load_time"`
	Issues          []PerformanceIssue `json:"issues"`
	Recommendations []string           `json:"recommendations"`
}

// AuditResult contains the result of project auditing
type AuditResult struct {
	ProjectPath     string                  `json:"project_path"`
	AuditTime       time.Time               `json:"audit_time"`
	Security        *SecurityAuditResult    `json:"security,omitempty"`
	Quality         *QualityAuditResult     `json:"quality,omitempty"`
	Licenses        *LicenseAuditResult     `json:"licenses,omitempty"`
	Performance     *PerformanceAuditResult `json:"performance,omitempty"`
	OverallScore    float64                 `json:"overall_score"`
	Recommendations []string                `json:"recommendations"`
}

// AuditOptions defines options for project auditing
type AuditOptions struct {
	Security     bool   `json:"security"`
	Quality      bool   `json:"quality"`
	Licenses     bool   `json:"licenses"`
	Performance  bool   `json:"performance"`
	OutputFormat string `json:"output_format"`
	OutputFile   string `json:"output_file"`
	Detailed     bool   `json:"detailed"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Package     string `json:"package"`
	Version     string `json:"version"`
	FixedIn     string `json:"fixed_in,omitempty"`
}

// PolicyViolation represents a security policy violation
type PolicyViolation struct {
	Policy      string `json:"policy"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

// CodeSmell represents a code quality issue
type CodeSmell struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

// Duplication represents code duplication
type Duplication struct {
	Files      []string `json:"files"`
	Lines      int      `json:"lines"`
	Tokens     int      `json:"tokens"`
	Percentage float64  `json:"percentage"`
}

// LicenseInfo represents license information
type LicenseInfo struct {
	Name       string `json:"name"`
	SPDXID     string `json:"spdx_id"`
	Package    string `json:"package"`
	Compatible bool   `json:"compatible"`
}

// PerformanceIssue represents a performance issue
type PerformanceIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	File        string `json:"file,omitempty"`
}

// SecretDetection represents a detected secret
type SecretDetection struct {
	Type       string  `json:"type"`
	File       string  `json:"file"`
	Line       int     `json:"line"`
	Column     int     `json:"column"`
	Secret     string  `json:"secret"`
	Confidence float64 `json:"confidence"`
	Rule       string  `json:"rule"`
	Pattern    string  `json:"pattern"`
	Masked     string  `json:"masked"`
}

// OutdatedDependency represents an outdated dependency
type OutdatedDependency struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Type           string `json:"type"`
	Severity       string `json:"severity"`
	UpdateType     string `json:"update_type"` // major, minor, patch
	Breaking       bool   `json:"breaking"`
}

// DependencyVulnerability represents a vulnerability in a dependency
type DependencyVulnerability struct {
	Dependency  string  `json:"dependency"`
	Version     string  `json:"version"`
	CVEID       string  `json:"cve_id"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	FixedIn     string  `json:"fixed_in,omitempty"`
	CVSS        float64 `json:"cvss"`
}
