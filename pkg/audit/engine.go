// Package audit provides auditing functionality for the
// Open Source Project Generator.
//
// Security Note: This package contains an audit engine that legitimately needs
// to read files for security, quality, and compliance analysis. The G304 warnings
// from gosec are false positives in this context as file reading is the core
// functionality of an audit tool.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/audit/license"
	"github.com/cuesoftinc/open-source-project-generator/pkg/audit/performance"
	"github.com/cuesoftinc/open-source-project-generator/pkg/audit/quality"
	"github.com/cuesoftinc/open-source-project-generator/pkg/audit/security"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Engine implements the AuditEngine interface for project auditing.
// It orchestrates various specialized audit components to provide comprehensive
// project analysis including security, quality, license, and performance auditing.
type Engine struct {
	// Core audit components
	ruleManager      *RuleManager
	resultAggregator *ResultAggregator

	// Specialized audit analyzers
	securityScanner    *security.SecurityScanner
	complexityAnalyzer *quality.ComplexityAnalyzer
	coverageAnalyzer   *quality.CoverageAnalyzer
	licenseChecker     *license.LicenseChecker
	bundleAnalyzer     *performance.BundleAnalyzer
	metricsAnalyzer    *performance.MetricsAnalyzer
}

// NewEngine creates a new audit engine instance with all specialized components.
func NewEngine() interfaces.AuditEngine {
	return &Engine{
		ruleManager:        NewRuleManager(),
		resultAggregator:   NewResultAggregator(),
		securityScanner:    security.NewSecurityScanner(),
		complexityAnalyzer: quality.NewComplexityAnalyzer(),
		coverageAnalyzer:   quality.NewCoverageAnalyzer(),
		licenseChecker:     license.NewLicenseChecker(),
		bundleAnalyzer:     performance.NewBundleAnalyzer(),
		metricsAnalyzer:    performance.NewMetricsAnalyzer(),
	}
}

// AuditSecurity performs security auditing on a project by delegating to the SecurityScanner.
func (e *Engine) AuditSecurity(path string) (*interfaces.SecurityAuditResult, error) {
	return e.securityScanner.AuditSecurity(path)
}

// ScanVulnerabilities scans for security vulnerabilities by delegating to the SecurityScanner.
func (e *Engine) ScanVulnerabilities(path string) (*interfaces.VulnerabilityReport, error) {
	return e.securityScanner.ScanVulnerabilities(path)
}

// CheckSecurityPolicies checks security policy compliance by delegating to the SecurityScanner.
func (e *Engine) CheckSecurityPolicies(path string) (*interfaces.PolicyComplianceResult, error) {
	return e.securityScanner.CheckSecurityPolicies(path)
}

// AuditCodeQuality performs code quality auditing by orchestrating quality analyzers.
func (e *Engine) AuditCodeQuality(path string) (*interfaces.QualityAuditResult, error) {
	result := &interfaces.QualityAuditResult{
		Score:           100.0,
		CodeSmells:      []interfaces.CodeSmell{},
		Duplications:    []interfaces.Duplication{},
		TestCoverage:    0.0,
		Recommendations: []string{},
	}

	// Analyze code smells using the coverage analyzer
	codeSmells, err := e.coverageAnalyzer.AnalyzeCodeSmells(path)
	if err != nil {
		return nil, fmt.Errorf("code smell analysis failed: %w", err)
	}
	result.CodeSmells = codeSmells

	// Analyze code duplications using the coverage analyzer
	duplications, err := e.coverageAnalyzer.AnalyzeDuplications(path)
	if err != nil {
		return nil, fmt.Errorf("duplication analysis failed: %w", err)
	}
	result.Duplications = duplications

	// Analyze test coverage using the coverage analyzer
	coverageResult, err := e.coverageAnalyzer.AnalyzeTestCoverage(path)
	if err != nil {
		// Test coverage analysis is optional, continue without error
		result.TestCoverage = 0.0
	} else {
		result.TestCoverage = coverageResult.OverallCoverage
	}

	// Calculate quality score using the result aggregator
	result.Score = e.resultAggregator.CalculateQualityScore(result)

	// Generate quality recommendations using the result aggregator
	result.Recommendations = e.resultAggregator.GenerateQualityRecommendations(result)

	return result, nil
}

// CheckBestPractices checks best practices compliance by orchestrating various analyzers.
func (e *Engine) CheckBestPractices(path string) (*interfaces.BestPracticesResult, error) {
	result := &interfaces.BestPracticesResult{
		Score:      100.0,
		Practices:  []interfaces.BestPracticeCheck{},
		Violations: []interfaces.BestPracticeViolation{},
		Summary: interfaces.BestPracticesSummary{
			TotalPractices:     0,
			CompliantPractices: 0,
			Violations:         0,
			OverallScore:       100.0,
		},
	}

	// Define best practices to check
	practices := []struct {
		id          string
		name        string
		description string
		checkFunc   func(string) bool
	}{
		{
			id:          "BP-001",
			name:        "README file exists",
			description: "Project should have a README file",
			checkFunc:   e.hasReadmeFile,
		},
		{
			id:          "BP-002",
			name:        "License file exists",
			description: "Project should have a LICENSE file",
			checkFunc:   e.licenseChecker.HasLicenseFile,
		},
		{
			id:          "BP-003",
			name:        "Gitignore file exists",
			description: "Project should have a .gitignore file",
			checkFunc:   func(p string) bool { return e.hasFile(p, ".gitignore") },
		},
	}

	// Check each best practice
	for _, practice := range practices {
		result.Summary.TotalPractices++

		practiceCheck := interfaces.BestPracticeCheck{
			ID:          practice.id,
			Name:        practice.name,
			Description: practice.description,
			Category:    "general",
			Technology:  "general",
			Compliant:   practice.checkFunc(path),
			Score:       0,
		}

		if practiceCheck.Compliant {
			practiceCheck.Score = 100.0
			result.Summary.CompliantPractices++
		} else {
			result.Violations = append(result.Violations, interfaces.BestPracticeViolation{
				Practice:    practice.name,
				File:        path,
				Line:        0,
				Description: fmt.Sprintf("%s not found", practice.name),
				Severity:    "medium",
				Suggestion:  fmt.Sprintf("Add %s to the project", practice.name),
			})
		}

		result.Practices = append(result.Practices, practiceCheck)
	}

	// Calculate overall score
	if result.Summary.TotalPractices > 0 {
		result.Score = float64(result.Summary.CompliantPractices) / float64(result.Summary.TotalPractices) * 100
		result.Summary.OverallScore = result.Score
	}

	return result, nil
}

// Helper methods for best practices checking
func (e *Engine) hasReadmeFile(path string) bool {
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "README", "readme.md", "readme.txt"}
	for _, readme := range readmeFiles {
		if e.hasFile(path, readme) {
			return true
		}
	}
	return false
}

func (e *Engine) hasFile(path, filename string) bool {
	_, err := os.Stat(filepath.Join(path, filename))
	return err == nil
}

// AnalyzeDependencies analyzes project dependencies by delegating to the SecurityScanner.
func (e *Engine) AnalyzeDependencies(path string) (*interfaces.DependencyAnalysisResult, error) {
	// Delegate to the SecurityScanner which handles dependency analysis
	return e.securityScanner.AnalyzeDependencies(path)
}

// AuditLicenses performs license auditing by delegating to the LicenseChecker.
func (e *Engine) AuditLicenses(path string) (*interfaces.LicenseAuditResult, error) {
	result, err := e.licenseChecker.AuditLicenses(path)
	if err != nil {
		return nil, err
	}

	// Generate recommendations using the result aggregator
	result.Recommendations = e.resultAggregator.GenerateLicenseRecommendations(result)

	return result, nil
}

// CheckLicenseCompatibility checks license compatibility by delegating to the LicenseChecker.
func (e *Engine) CheckLicenseCompatibility(path string) (*interfaces.LicenseCompatibilityResult, error) {
	return e.licenseChecker.CheckLicenseCompatibility(path)
}

// AuditPerformance performs performance auditing by orchestrating performance analyzers.
func (e *Engine) AuditPerformance(path string) (*interfaces.PerformanceAuditResult, error) {
	result := &interfaces.PerformanceAuditResult{
		Score:           100.0,
		BundleSize:      0,
		LoadTime:        0,
		Issues:          []interfaces.PerformanceIssue{},
		Recommendations: []string{},
	}

	// Analyze bundle size using the bundle analyzer
	bundleResult, err := e.bundleAnalyzer.AnalyzeBundleSize(path)
	if err == nil {
		result.BundleSize = bundleResult.TotalSize
	}

	// Check performance metrics using the metrics analyzer
	metricsResult, err := e.metricsAnalyzer.CheckPerformanceMetrics(path)
	if err == nil {
		result.LoadTime = metricsResult.LoadTime
	}

	// Analyze performance issues using the bundle analyzer
	issues, err := e.bundleAnalyzer.AnalyzePerformanceIssues(path)
	if err != nil {
		return nil, fmt.Errorf("performance analysis failed: %w", err)
	}
	result.Issues = issues

	// Calculate performance score using the result aggregator
	result.Score = e.resultAggregator.CalculatePerformanceScore(result)

	// Generate recommendations using the result aggregator
	result.Recommendations = e.resultAggregator.GeneratePerformanceRecommendations(result)

	return result, nil
}

// AnalyzeBundleSize analyzes bundle size and performance by delegating to the BundleAnalyzer.
func (e *Engine) AnalyzeBundleSize(path string) (*interfaces.BundleAnalysisResult, error) {
	return e.bundleAnalyzer.AnalyzeBundleSize(path)
}

// DetectSecrets detects secrets in the project by delegating to the SecurityScanner.
func (e *Engine) DetectSecrets(path string) (*interfaces.SecretScanResult, error) {
	return e.securityScanner.DetectSecrets(path)
}

// MeasureComplexity measures code complexity by delegating to the ComplexityAnalyzer.
func (e *Engine) MeasureComplexity(path string) (*interfaces.ComplexityAnalysisResult, error) {
	return e.complexityAnalyzer.MeasureComplexity(path)
}

// ScanLicenseViolations scans for license violations by delegating to the LicenseChecker.
func (e *Engine) ScanLicenseViolations(path string) (*interfaces.LicenseViolationResult, error) {
	return e.licenseChecker.ScanLicenseViolations(path)
}

// CheckPerformanceMetrics checks performance metrics by delegating to the MetricsAnalyzer.
func (e *Engine) CheckPerformanceMetrics(path string) (*interfaces.PerformanceMetricsResult, error) {
	return e.metricsAnalyzer.CheckPerformanceMetrics(path)
}

// AuditProject performs comprehensive project auditing by orchestrating all audit components.
// This is the main entry point for complete project analysis.
func (e *Engine) AuditProject(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	// Set default options if none provided
	if options == nil {
		options = &interfaces.AuditOptions{
			Security:    true,
			Quality:     true,
			Licenses:    true,
			Performance: true,
		}
	}

	result := &interfaces.AuditResult{
		ProjectPath: path,
		AuditTime:   time.Now(),
	}

	// Perform security audit if requested
	if options.Security {
		securityResult, err := e.securityScanner.AuditSecurity(path)
		if err != nil {
			return nil, fmt.Errorf("security audit failed: %w", err)
		}
		result.Security = securityResult
	}

	// Perform quality audit if requested
	if options.Quality {
		qualityResult, err := e.AuditCodeQuality(path)
		if err != nil {
			return nil, fmt.Errorf("code quality audit failed: %w", err)
		}
		result.Quality = qualityResult
	}

	// Perform license audit if requested
	if options.Licenses {
		licenseResult, err := e.licenseChecker.AuditLicenses(path)
		if err != nil {
			return nil, fmt.Errorf("license audit failed: %w", err)
		}
		// Generate recommendations using the result aggregator
		licenseResult.Recommendations = e.resultAggregator.GenerateLicenseRecommendations(licenseResult)
		result.Licenses = licenseResult
	}

	// Perform performance audit if requested
	if options.Performance {
		performanceResult, err := e.AuditPerformance(path)
		if err != nil {
			return nil, fmt.Errorf("performance audit failed: %w", err)
		}
		result.Performance = performanceResult
	}

	// Calculate overall score using the result aggregator
	result.OverallScore = e.resultAggregator.CalculateOverallScore(result)

	// Generate comprehensive recommendations using the result aggregator
	result.Recommendations = e.resultAggregator.GenerateRecommendations(result)

	return result, nil
}

// GenerateAuditReport generates an audit report by delegating to the ResultAggregator.
func (e *Engine) GenerateAuditReport(result *interfaces.AuditResult, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(result, "", "  ")
	case "html":
		return e.resultAggregator.GenerateHTMLReport(result)
	case "markdown":
		return e.resultAggregator.GenerateMarkdownReport(result)
	default:
		return nil, fmt.Errorf("unsupported report format '%s'. Available formats: json, html, markdown", format)
	}
}

// GetAuditSummary gets audit summary by delegating to the ResultAggregator.
func (e *Engine) GetAuditSummary(results []*interfaces.AuditResult) (*interfaces.AuditSummary, error) {
	return e.resultAggregator.AggregateResults(results)
}

// SetAuditRules sets audit rules by delegating to the RuleManager.
func (e *Engine) SetAuditRules(rules []interfaces.AuditRule) error {
	return e.ruleManager.SetRules(rules)
}

// GetAuditRules gets audit rules by delegating to the RuleManager.
func (e *Engine) GetAuditRules() []interfaces.AuditRule {
	return e.ruleManager.GetRules()
}

// AddAuditRule adds an audit rule by delegating to the RuleManager.
func (e *Engine) AddAuditRule(rule interfaces.AuditRule) error {
	return e.ruleManager.AddRule(rule)
}

// RemoveAuditRule removes an audit rule by delegating to the RuleManager.
func (e *Engine) RemoveAuditRule(ruleID string) error {
	return e.ruleManager.RemoveRule(ruleID)
}
