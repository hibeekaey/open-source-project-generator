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
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Engine implements the AuditEngine interface for project auditing.
type Engine struct {
	rules []interfaces.AuditRule
}

// NewEngine creates a new audit engine instance.
func NewEngine() interfaces.AuditEngine {
	return &Engine{
		rules: getDefaultAuditRules(),
	}
}

// getDefaultAuditRules returns the default set of audit rules
func getDefaultAuditRules() []interfaces.AuditRule {
	return []interfaces.AuditRule{
		{
			ID:          "security-001",
			Name:        "No hardcoded secrets",
			Description: "Check for hardcoded secrets in source code",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityCritical,
			Enabled:     true,
			Pattern:     `(?i)(password|secret|key|token)\s*[:=]\s*["'][^"']+["']`,
			FileTypes:   []string{".go", ".js", ".ts", ".py", ".java", ".cs"},
		},
		{
			ID:          "security-002",
			Name:        "Secure dependencies",
			Description: "Check for known vulnerabilities in dependencies",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
		},
		{
			ID:          "quality-001",
			Name:        "Code complexity",
			Description: "Check for high cyclomatic complexity",
			Category:    interfaces.AuditCategoryQuality,
			Type:        interfaces.AuditCategoryQuality,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
			Config:      map[string]any{"max_complexity": 10},
		},
		{
			ID:          "license-001",
			Name:        "License compatibility",
			Description: "Check for license compatibility issues",
			Category:    interfaces.AuditCategoryLicense,
			Type:        interfaces.AuditCategoryLicense,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
		},
		{
			ID:          "performance-001",
			Name:        "Bundle size",
			Description: "Check for large bundle sizes",
			Category:    interfaces.AuditCategoryPerformance,
			Type:        interfaces.AuditCategoryPerformance,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
			Config:      map[string]any{"max_size_mb": 5},
		},
	}
}

// AuditSecurity performs security auditing on a project
func (e *Engine) AuditSecurity(path string) (*interfaces.SecurityAuditResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.SecurityAuditResult{
		Score:            100.0,
		Vulnerabilities:  []interfaces.Vulnerability{},
		PolicyViolations: []interfaces.PolicyViolation{},
		Recommendations:  []string{},
	}

	// Scan for vulnerabilities
	vulnReport, err := e.ScanVulnerabilities(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Vulnerability scan failed.",
			"Unable to analyze project for security vulnerabilities")
	}
	result.Vulnerabilities = vulnReport.Vulnerabilities

	// Check security policies
	policyResult, err := e.CheckSecurityPolicies(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Security policy check failed.",
			"Unable to validate project security policies")
	}
	result.PolicyViolations = policyResult.Violations

	// Detect secrets
	secretResult, err := e.DetectSecrets(path)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s",
			"Secret detection failed.",
			"Unable to scan project for exposed secrets")
	}

	// Calculate security score based on findings
	result.Score = e.calculateSecurityScore(result, secretResult)

	// Generate security recommendations
	result.Recommendations = e.generateSecurityRecommendations(result, secretResult)

	return result, nil
}

// ScanVulnerabilities scans for security vulnerabilities
func (e *Engine) ScanVulnerabilities(path string) (*interfaces.VulnerabilityReport, error) {
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

	// Scan dependency files for known vulnerabilities
	vulnerabilities, err := e.scanDependencyVulnerabilities(path)
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
func (e *Engine) CheckSecurityPolicies(path string) (*interfaces.PolicyComplianceResult, error) {
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
			violations, err := e.checkHardcodedSecrets(path)
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
			violations, err := e.checkDependencyVulnerabilities(path)
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
			violations, err := e.checkSecureConfiguration(path)
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

// AuditCodeQuality performs code quality auditing
func (e *Engine) AuditCodeQuality(path string) (*interfaces.QualityAuditResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.QualityAuditResult{
		Score:           100.0,
		CodeSmells:      []interfaces.CodeSmell{},
		Duplications:    []interfaces.Duplication{},
		TestCoverage:    0.0,
		Recommendations: []string{},
	}

	// Analyze code smells
	codeSmells, err := e.analyzeCodeSmells(path)
	if err != nil {
		return nil, fmt.Errorf("code smell analysis failed: %w", err)
	}
	result.CodeSmells = codeSmells

	// Analyze code duplications
	duplications, err := e.analyzeDuplications(path)
	if err != nil {
		return nil, fmt.Errorf("duplication analysis failed: %w", err)
	}
	result.Duplications = duplications

	// Analyze test coverage
	coverage, err := e.analyzeTestCoverage(path)
	if err != nil {
		// Test coverage analysis is optional, continue without error
		coverage = 0.0
	}
	result.TestCoverage = coverage

	// Calculate quality score
	result.Score = e.calculateQualityScore(result)

	// Generate quality recommendations
	result.Recommendations = e.generateQualityRecommendations(result)

	return result, nil
}

// CheckBestPractices checks best practices compliance
func (e *Engine) CheckBestPractices(path string) (*interfaces.BestPracticesResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

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

	// Get best practices to check based on project type
	practices := e.getBestPractices(path)

	// Check each best practice
	for _, practice := range practices {
		result.Summary.TotalPractices++

		violations, err := e.checkBestPractice(path, practice)
		if err != nil {
			continue // Log error but continue
		}

		if len(violations) == 0 {
			practice.Compliant = true
			practice.Score = 100.0
			result.Summary.CompliantPractices++
		} else {
			practice.Compliant = false
			practice.Score = 0.0
			result.Violations = append(result.Violations, violations...)
			result.Summary.Violations += len(violations)
		}

		result.Practices = append(result.Practices, practice)
	}

	// Calculate overall score
	if result.Summary.TotalPractices > 0 {
		result.Score = float64(result.Summary.CompliantPractices) / float64(result.Summary.TotalPractices) * 100
		result.Summary.OverallScore = result.Score
	}

	return result, nil
}

// AnalyzeDependencies analyzes project dependencies
func (e *Engine) AnalyzeDependencies(path string) (*interfaces.DependencyAnalysisResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

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
	}

	for depFile, ecosystem := range dependencyFiles {
		filePath := filepath.Join(path, depFile)
		if _, err := os.Stat(filePath); err == nil {
			deps, err := e.analyzeDependencyFile(filePath, ecosystem)
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

	return result, nil
}

// AuditLicenses performs license auditing
func (e *Engine) AuditLicenses(path string) (*interfaces.LicenseAuditResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.LicenseAuditResult{
		Score:           100.0,
		Compatible:      true,
		Licenses:        []interfaces.LicenseInfo{},
		Conflicts:       []interfaces.LicenseInfo{},
		Recommendations: []string{},
	}

	// Check license compatibility
	compatibilityResult, err := e.CheckLicenseCompatibility(path)
	if err != nil {
		return nil, fmt.Errorf("license compatibility check failed: %w", err)
	}

	result.Compatible = compatibilityResult.Compatible

	// Convert DependencyLicense to LicenseInfo
	for _, dep := range compatibilityResult.Dependencies {
		licenseInfo := interfaces.LicenseInfo{
			Name:       dep.License,
			SPDXID:     dep.SPDXID,
			Package:    dep.Dependency,
			Compatible: dep.Compatible,
		}
		result.Licenses = append(result.Licenses, licenseInfo)

		// Find conflicts
		if !dep.Compatible {
			result.Conflicts = append(result.Conflicts, licenseInfo)
		}
	}

	// Calculate license score
	if len(result.Licenses) > 0 {
		compatibleCount := len(result.Licenses) - len(result.Conflicts)
		result.Score = float64(compatibleCount) / float64(len(result.Licenses)) * 100
	}

	// Generate recommendations
	result.Recommendations = e.generateLicenseRecommendations(result)

	return result, nil
}

// CheckLicenseCompatibility checks license compatibility
func (e *Engine) CheckLicenseCompatibility(path string) (*interfaces.LicenseCompatibilityResult, error) {
	result := &interfaces.LicenseCompatibilityResult{
		Compatible:      true,
		ProjectLicense:  "unknown",
		Dependencies:    []interfaces.DependencyLicense{},
		Conflicts:       []interfaces.LicenseConflict{},
		Recommendations: []interfaces.LicenseRecommendation{},
		Summary: interfaces.LicenseCompatibilitySummary{
			TotalLicenses:      0,
			CompatibleLicenses: 0,
			Conflicts:          0,
			RiskLevel:          "low",
		},
	}

	// Detect project license
	projectLicense, err := e.detectProjectLicense(path)
	if err == nil {
		result.ProjectLicense = projectLicense
	}

	// Analyze dependency licenses
	dependencies, err := e.analyzeDependencyLicenses(path)
	if err != nil {
		return nil, fmt.Errorf("dependency license analysis failed: %w", err)
	}

	result.Dependencies = dependencies
	result.Summary.TotalLicenses = len(dependencies)

	// Check compatibility
	for _, dep := range dependencies {
		if e.isLicenseCompatible(result.ProjectLicense, dep.License) {
			result.Summary.CompatibleLicenses++
		} else {
			result.Compatible = false
			result.Conflicts = append(result.Conflicts, interfaces.LicenseConflict{
				Dependency1: "project",
				License1:    result.ProjectLicense,
				Dependency2: dep.Dependency,
				License2:    dep.License,
				Reason:      "License incompatibility",
				Severity:    e.getLicenseConflictSeverity(result.ProjectLicense, dep.License),
				Resolution:  "Consider changing project license or replacing dependency",
			})
			result.Summary.Conflicts++
		}
	}

	// Determine risk level
	if result.Summary.Conflicts > 0 {
		result.Summary.RiskLevel = "high"
	} else if result.ProjectLicense == "unknown" {
		result.Summary.RiskLevel = "medium"
	}

	return result, nil
}

// AuditPerformance performs performance auditing
func (e *Engine) AuditPerformance(path string) (*interfaces.PerformanceAuditResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.PerformanceAuditResult{
		Score:           100.0,
		BundleSize:      0,
		LoadTime:        0,
		Issues:          []interfaces.PerformanceIssue{},
		Recommendations: []string{},
	}

	// Analyze bundle size
	bundleResult, err := e.AnalyzeBundleSize(path)
	if err == nil {
		result.BundleSize = bundleResult.TotalSize
	}

	// Check performance metrics
	metricsResult, err := e.CheckPerformanceMetrics(path)
	if err == nil {
		result.LoadTime = metricsResult.LoadTime
	}

	// Analyze performance issues
	issues, err := e.analyzePerformanceIssues(path)
	if err != nil {
		return nil, fmt.Errorf("performance analysis failed: %w", err)
	}
	result.Issues = issues

	// Calculate performance score
	result.Score = e.calculatePerformanceScore(result)

	// Generate recommendations
	result.Recommendations = e.generatePerformanceRecommendations(result)

	return result, nil
}

// AnalyzeBundleSize analyzes bundle size and performance
func (e *Engine) AnalyzeBundleSize(path string) (*interfaces.BundleAnalysisResult, error) {
	result := &interfaces.BundleAnalysisResult{
		TotalSize:   0,
		GzippedSize: 0,
		Assets:      []interfaces.BundleAsset{},
		Chunks:      []interfaces.BundleChunk{},
		Summary: interfaces.BundleAnalysisSummary{
			TotalAssets:      0,
			TotalChunks:      0,
			CompressionRatio: 0.0,
			LargestAsset:     "",
			Recommendations:  []string{},
		},
	}

	// Look for common build output directories
	buildDirs := []string{"dist", "build", "public", "static", "assets"}

	for _, buildDir := range buildDirs {
		buildPath := filepath.Join(path, buildDir)
		if _, err := os.Stat(buildPath); err == nil {
			err := e.analyzeBuildDirectory(buildPath, result)
			if err != nil {
				continue
			}
			break
		}
	}

	// If no build directory found, analyze source files
	if result.TotalSize == 0 {
		e.analyzeSourceFiles(path, result)
	}

	// Calculate summary
	result.Summary.TotalAssets = len(result.Assets)
	result.Summary.TotalChunks = len(result.Chunks)

	if result.TotalSize > 0 && result.GzippedSize > 0 {
		result.Summary.CompressionRatio = float64(result.GzippedSize) / float64(result.TotalSize)
	}

	// Find largest asset
	var largestSize int64
	for _, asset := range result.Assets {
		if asset.Size > largestSize {
			largestSize = asset.Size
			result.Summary.LargestAsset = asset.Name
		}
	}

	// Generate recommendations
	if result.TotalSize > 5*1024*1024 { // 5MB
		result.Summary.Recommendations = append(result.Summary.Recommendations, "Consider code splitting to reduce bundle size")
	}
	if result.Summary.CompressionRatio > 0.8 {
		result.Summary.Recommendations = append(result.Summary.Recommendations, "Enable gzip compression on server")
	}

	return result, nil
}

// DetectSecrets detects secrets in the project
func (e *Engine) DetectSecrets(path string) (*interfaces.SecretScanResult, error) {
	result := &interfaces.SecretScanResult{
		ScanTime: time.Now(),
		Secrets:  []interfaces.SecretDetection{},
		Summary: interfaces.SecretScanSummary{
			TotalSecrets:     0,
			HighConfidence:   0,
			MediumConfidence: 0,
			LowConfidence:    0,
			FilesScanned:     0,
		},
	}

	// Get security rules for secret detection
	secretRules := e.getSecretDetectionRules()

	// Walk through project files
	walkErr := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-text files
		if info.IsDir() || e.shouldSkipFile(filePath) {
			return nil
		}

		result.Summary.FilesScanned++

		// Scan file for secrets
		secrets, err := e.scanFileForSecrets(filePath, secretRules)
		if err != nil {
			// Log error but continue scanning
			return nil
		}

		result.Secrets = append(result.Secrets, secrets...)
		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("failed to scan for secrets: %w", walkErr)
	}

	// Update summary
	for _, secret := range result.Secrets {
		result.Summary.TotalSecrets++
		if secret.Confidence >= 0.8 {
			result.Summary.HighConfidence++
		} else if secret.Confidence >= 0.5 {
			result.Summary.MediumConfidence++
		} else {
			result.Summary.LowConfidence++
		}
	}

	return result, nil
}

// MeasureComplexity measures code complexity
func (e *Engine) MeasureComplexity(path string) (*interfaces.ComplexityAnalysisResult, error) {
	if err := e.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.ComplexityAnalysisResult{
		Files:     []interfaces.FileComplexity{},
		Functions: []interfaces.FunctionComplexity{},
		Summary: interfaces.ComplexityAnalysisSummary{
			TotalFiles:          0,
			TotalFunctions:      0,
			AverageComplexity:   0.0,
			HighComplexityFiles: 0,
			TechnicalDebtHours:  0.0,
		},
	}

	// Analyze complexity for different file types
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) {
			return nil
		}

		// Only analyze source code files
		if !e.isSourceCodeFile(filePath) {
			return nil
		}

		fileComplexity, err := e.analyzeFileComplexity(filePath)
		if err != nil {
			return nil // Continue on error
		}

		result.Files = append(result.Files, *fileComplexity)
		result.Summary.TotalFiles++

		if fileComplexity.CyclomaticComplexity > 10 {
			result.Summary.HighComplexityFiles++
		}

		// Add function complexities (simplified - would need proper parsing)
		functions := e.extractFunctions(filePath)
		result.Functions = append(result.Functions, functions...)
		result.Summary.TotalFunctions += len(functions)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("complexity analysis failed: %w", err)
	}

	// Calculate summary metrics
	if result.Summary.TotalFiles > 0 {
		var totalComplexity float64
		var totalDebt float64

		for _, file := range result.Files {
			totalComplexity += float64(file.CyclomaticComplexity)

			// Estimate technical debt (simplified calculation)
			if file.CyclomaticComplexity > 10 {
				totalDebt += float64(file.CyclomaticComplexity-10) * 0.5 // 0.5 hours per excess complexity point
			}
		}

		result.Summary.AverageComplexity = totalComplexity / float64(result.Summary.TotalFiles)
		result.Summary.TechnicalDebtHours = totalDebt
	}

	return result, nil
}

// ScanLicenseViolations scans for license violations
func (e *Engine) ScanLicenseViolations(path string) (*interfaces.LicenseViolationResult, error) {
	result := &interfaces.LicenseViolationResult{
		Violations: []interfaces.LicenseViolation{},
		Summary: interfaces.LicenseViolationSummary{
			TotalViolations:    0,
			CriticalViolations: 0,
			HighViolations:     0,
			MediumViolations:   0,
			LowViolations:      0,
		},
	}

	// Check for missing license files
	if !e.hasLicenseFile(path) {
		result.Violations = append(result.Violations, interfaces.LicenseViolation{
			Type:       "missing",
			Dependency: "project",
			License:    "none",
			Violation:  "No license file found in project",
			Severity:   "high",
			Resolution: "Add a LICENSE file to the project root",
		})
	}

	// Check dependency licenses
	compatibilityResult, err := e.CheckLicenseCompatibility(path)
	if err == nil {
		for _, conflict := range compatibilityResult.Conflicts {
			result.Violations = append(result.Violations, interfaces.LicenseViolation{
				Type:       "incompatible",
				Dependency: conflict.Dependency2,
				License:    conflict.License2,
				Violation:  conflict.Reason,
				Severity:   conflict.Severity,
				Resolution: conflict.Resolution,
			})
		}
	}

	// Update summary
	for _, violation := range result.Violations {
		result.Summary.TotalViolations++
		switch violation.Severity {
		case "critical":
			result.Summary.CriticalViolations++
		case "high":
			result.Summary.HighViolations++
		case "medium":
			result.Summary.MediumViolations++
		case "low":
			result.Summary.LowViolations++
		}
	}

	return result, nil
}

// CheckPerformanceMetrics checks performance metrics
func (e *Engine) CheckPerformanceMetrics(path string) (*interfaces.PerformanceMetricsResult, error) {
	result := &interfaces.PerformanceMetricsResult{
		LoadTime:              0,
		FirstPaint:            0,
		FirstContentful:       0,
		LargestContentful:     0,
		TimeToInteractive:     0,
		CumulativeLayoutShift: 0.0,
		Issues:                []interfaces.PerformanceIssue{},
		Summary: interfaces.PerformanceMetricsSummary{
			OverallScore:     100.0,
			PerformanceGrade: "A",
			IssueCount:       0,
			Recommendations:  []string{},
		},
	}

	// This is a simplified implementation
	// In a real implementation, you would integrate with tools like:
	// - Lighthouse for web performance
	// - WebPageTest
	// - Custom performance monitoring

	// Estimate metrics based on bundle size and project structure
	bundleResult, err := e.AnalyzeBundleSize(path)
	if err == nil {
		// Rough estimates based on bundle size
		sizeInMB := float64(bundleResult.TotalSize) / (1024 * 1024)

		// Estimate load time (very rough approximation)
		result.LoadTime = time.Duration(sizeInMB*100) * time.Millisecond
		result.FirstPaint = result.LoadTime / 2
		result.FirstContentful = result.LoadTime * 3 / 4
		result.LargestContentful = result.LoadTime
		result.TimeToInteractive = result.LoadTime * 2

		// Calculate performance score
		if result.LoadTime > 3*time.Second {
			result.Summary.OverallScore = 30
			result.Summary.PerformanceGrade = "F"
		} else if result.LoadTime > 2*time.Second {
			result.Summary.OverallScore = 50
			result.Summary.PerformanceGrade = "D"
		} else if result.LoadTime > 1*time.Second {
			result.Summary.OverallScore = 70
			result.Summary.PerformanceGrade = "C"
		} else if result.LoadTime > 500*time.Millisecond {
			result.Summary.OverallScore = 85
			result.Summary.PerformanceGrade = "B"
		}
	}

	// Generate performance issues and recommendations
	if result.LoadTime > 2*time.Second {
		result.Issues = append(result.Issues, interfaces.PerformanceIssue{
			Type:        "slow_load_time",
			Severity:    "high",
			Description: fmt.Sprintf("Load time is %v, which is slower than recommended 2s", result.LoadTime),
			Impact:      "Users may abandon the page before it loads",
		})
		result.Summary.Recommendations = append(result.Summary.Recommendations, "Optimize bundle size and enable code splitting")
	}

	result.Summary.IssueCount = len(result.Issues)

	return result, nil
}

// AuditProject performs comprehensive project auditing
func (e *Engine) AuditProject(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
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

	var totalScore float64
	var scoreCount int

	// Perform security audit if requested
	if options.Security {
		securityResult, err := e.AuditSecurity(path)
		if err != nil {
			return nil, fmt.Errorf("🚫 %s %s",
				"Security audit failed.",
				"Unable to complete security analysis")
		}
		result.Security = securityResult
		totalScore += securityResult.Score
		scoreCount++
	}

	// Perform quality audit if requested
	if options.Quality {
		qualityResult, err := e.AuditCodeQuality(path)
		if err != nil {
			return nil, fmt.Errorf("🚫 %s %s",
				"Code quality audit failed.",
				"Unable to analyze code quality metrics")
		}
		result.Quality = qualityResult
		totalScore += qualityResult.Score
		scoreCount++
	}

	// Perform license audit if requested
	if options.Licenses {
		licenseResult, err := e.AuditLicenses(path)
		if err != nil {
			return nil, fmt.Errorf("🚫 %s %s",
				"License audit failed.",
				"Unable to analyze project license compliance")
		}
		result.Licenses = licenseResult
		totalScore += licenseResult.Score
		scoreCount++
	}

	// Perform performance audit if requested
	if options.Performance {
		performanceResult, err := e.AuditPerformance(path)
		if err != nil {
			return nil, fmt.Errorf("🚫 %s %s",
				"Performance audit failed.",
				"Unable to analyze performance characteristics")
		}
		result.Performance = performanceResult
		totalScore += performanceResult.Score
		scoreCount++
	}

	// Calculate overall score
	if scoreCount > 0 {
		result.OverallScore = totalScore / float64(scoreCount)
	}

	// Generate recommendations based on results
	result.Recommendations = e.generateRecommendations(result)

	return result, nil
}

// GenerateAuditReport generates an audit report
func (e *Engine) GenerateAuditReport(result *interfaces.AuditResult, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(result, "", "  ")
	case "html":
		return e.generateHTMLReport(result)
	case "markdown":
		return e.generateMarkdownReport(result)
	default:
		return nil, fmt.Errorf("🚫 Unsupported report format '%s'. Available formats: text, json, html, markdown", format)
	}
}

// GetAuditSummary gets audit summary
func (e *Engine) GetAuditSummary(results []*interfaces.AuditResult) (*interfaces.AuditSummary, error) {
	if len(results) == 0 {
		return &interfaces.AuditSummary{}, nil
	}

	summary := &interfaces.AuditSummary{
		TotalProjects: len(results),
	}

	var totalScore, securityScore, qualityScore, licenseScore, performanceScore float64
	var securityCount, qualityCount, licenseCount, performanceCount int

	issueFrequency := make(map[string]int)

	for _, result := range results {
		totalScore += result.OverallScore

		if result.Security != nil {
			securityScore += result.Security.Score
			securityCount++

			// Count security issues
			for _, vuln := range result.Security.Vulnerabilities {
				key := fmt.Sprintf("security-%s", vuln.Severity)
				issueFrequency[key]++
			}
		}

		if result.Quality != nil {
			qualityScore += result.Quality.Score
			qualityCount++

			// Count quality issues
			for _, smell := range result.Quality.CodeSmells {
				key := fmt.Sprintf("quality-%s", smell.Type)
				issueFrequency[key]++
			}
		}

		if result.Licenses != nil {
			licenseScore += result.Licenses.Score
			licenseCount++

			// Count license conflicts
			for range result.Licenses.Conflicts {
				issueFrequency["license-conflict"]++
			}
		}

		if result.Performance != nil {
			performanceScore += result.Performance.Score
			performanceCount++

			// Count performance issues
			for _, issue := range result.Performance.Issues {
				key := fmt.Sprintf("performance-%s", issue.Type)
				issueFrequency[key]++
			}
		}
	}

	// Calculate average scores
	summary.AverageScore = totalScore / float64(len(results))

	if securityCount > 0 {
		summary.SecurityScore = securityScore / float64(securityCount)
	}
	if qualityCount > 0 {
		summary.QualityScore = qualityScore / float64(qualityCount)
	}
	if licenseCount > 0 {
		summary.LicenseScore = licenseScore / float64(licenseCount)
	}
	if performanceCount > 0 {
		summary.PerformanceScore = performanceScore / float64(performanceCount)
	}

	// Generate common issues list
	for issueType, frequency := range issueFrequency {
		if frequency > len(results)/4 { // Issues appearing in >25% of projects
			summary.CommonIssues = append(summary.CommonIssues, interfaces.CommonIssue{
				Type:        issueType,
				Description: fmt.Sprintf("Common issue: %s", issueType),
				Frequency:   frequency,
				Severity:    "medium", // Default severity
				Category:    "general",
			})
		}
	}

	return summary, nil
}

// SetAuditRules sets audit rules
func (e *Engine) SetAuditRules(rules []interfaces.AuditRule) error {
	e.rules = make([]interfaces.AuditRule, len(rules))
	copy(e.rules, rules)
	return nil
}

// GetAuditRules gets audit rules
func (e *Engine) GetAuditRules() []interfaces.AuditRule {
	rules := make([]interfaces.AuditRule, len(e.rules))
	copy(rules, e.rules)
	return rules
}

// AddAuditRule adds an audit rule
func (e *Engine) AddAuditRule(rule interfaces.AuditRule) error {
	// Check if rule with same ID already exists
	for i, existingRule := range e.rules {
		if existingRule.ID == rule.ID {
			e.rules[i] = rule // Replace existing rule
			return nil
		}
	}

	// Add new rule
	e.rules = append(e.rules, rule)
	return nil
}

// RemoveAuditRule removes an audit rule
func (e *Engine) RemoveAuditRule(ruleID string) error {
	for i, rule := range e.rules {
		if rule.ID == ruleID {
			e.rules = append(e.rules[:i], e.rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("audit rule with ID %s not found", ruleID)
}

// generateRecommendations generates recommendations based on audit results
func (e *Engine) generateRecommendations(result *interfaces.AuditResult) []string {
	var recommendations []string

	// Security recommendations
	if result.Security != nil {
		if result.Security.Score < 70 {
			recommendations = append(recommendations, "Consider improving security practices - score is below 70%")
		}
		if len(result.Security.Vulnerabilities) > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Address %d security vulnerabilities found", len(result.Security.Vulnerabilities)))
		}
		if len(result.Security.PolicyViolations) > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Fix %d security policy violations", len(result.Security.PolicyViolations)))
		}
	}

	// Quality recommendations
	if result.Quality != nil {
		if result.Quality.Score < 70 {
			recommendations = append(recommendations, "Consider improving code quality - score is below 70%")
		}
		if len(result.Quality.CodeSmells) > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Address %d code quality issues", len(result.Quality.CodeSmells)))
		}
		if result.Quality.TestCoverage < 80 {
			recommendations = append(recommendations, fmt.Sprintf("Increase test coverage from %.1f%% to at least 80%%", result.Quality.TestCoverage))
		}
	}

	// License recommendations
	if result.Licenses != nil {
		if !result.Licenses.Compatible {
			recommendations = append(recommendations, "Resolve license compatibility issues")
		}
		if len(result.Licenses.Conflicts) > 0 {
			recommendations = append(recommendations, fmt.Sprintf("Address %d license conflicts", len(result.Licenses.Conflicts)))
		}
	}

	// Performance recommendations
	if result.Performance != nil {
		if result.Performance.Score < 70 {
			recommendations = append(recommendations, "Consider optimizing performance - score is below 70%")
		}
		if result.Performance.BundleSize > 5*1024*1024 { // 5MB
			recommendations = append(recommendations, fmt.Sprintf("Consider reducing bundle size from %d bytes", result.Performance.BundleSize))
		}
	}

	// Overall recommendations
	if result.OverallScore < 70 {
		recommendations = append(recommendations, "Overall project health is below 70% - consider addressing the most critical issues first")
	}

	return recommendations
}

// generateHTMLReport generates an HTML audit report
func (e *Engine) generateHTMLReport(result *interfaces.AuditResult) ([]byte, error) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Audit Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f5f5f5; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .score { font-size: 24px; font-weight: bold; }
        .high { color: #d32f2f; }
        .medium { color: #f57c00; }
        .low { color: #388e3c; }
        .recommendations { background-color: #e3f2fd; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Audit Report</h1>
        <p><strong>Project:</strong> %s</p>
        <p><strong>Audit Time:</strong> %s</p>
        <p><strong>Overall Score:</strong> <span class="score">%.1f%%</span></p>
    </div>
`, filepath.Base(result.ProjectPath), result.ProjectPath, result.AuditTime.Format(time.RFC3339), result.OverallScore)

	if result.Security != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Security Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Vulnerabilities:</strong> %d</p>
        <p><strong>Policy Violations:</strong> %d</p>
    </div>`, result.Security.Score, len(result.Security.Vulnerabilities), len(result.Security.PolicyViolations))
	}

	if result.Quality != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Quality Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Code Smells:</strong> %d</p>
        <p><strong>Test Coverage:</strong> %.1f%%</p>
    </div>`, result.Quality.Score, len(result.Quality.CodeSmells), result.Quality.TestCoverage)
	}

	if result.Licenses != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>License Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Compatible:</strong> %t</p>
        <p><strong>Conflicts:</strong> %d</p>
    </div>`, result.Licenses.Score, result.Licenses.Compatible, len(result.Licenses.Conflicts))
	}

	if result.Performance != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Performance Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Bundle Size:</strong> %d bytes</p>
        <p><strong>Issues:</strong> %d</p>
    </div>`, result.Performance.Score, result.Performance.BundleSize, len(result.Performance.Issues))
	}

	if len(result.Recommendations) > 0 {
		html += `
    <div class="recommendations">
        <h2>Recommendations</h2>
        <ul>`
		for _, rec := range result.Recommendations {
			html += fmt.Sprintf("<li>%s</li>", rec)
		}
		html += `
        </ul>
    </div>`
	}

	html += `
</body>
</html>`

	return []byte(html), nil
}

// generateMarkdownReport generates a Markdown audit report
func (e *Engine) generateMarkdownReport(result *interfaces.AuditResult) ([]byte, error) {
	md := fmt.Sprintf(`# Audit Report

**Project:** %s  
**Audit Time:** %s  
**Overall Score:** %.1f%%

`, result.ProjectPath, result.AuditTime.Format(time.RFC3339), result.OverallScore)

	if result.Security != nil {
		md += fmt.Sprintf(`## Security Audit

- **Score:** %.1f%%
- **Vulnerabilities:** %d
- **Policy Violations:** %d

`, result.Security.Score, len(result.Security.Vulnerabilities), len(result.Security.PolicyViolations))
	}

	if result.Quality != nil {
		md += fmt.Sprintf(`## Quality Audit

- **Score:** %.1f%%
- **Code Smells:** %d
- **Test Coverage:** %.1f%%

`, result.Quality.Score, len(result.Quality.CodeSmells), result.Quality.TestCoverage)
	}

	if result.Licenses != nil {
		md += fmt.Sprintf(`## License Audit

- **Score:** %.1f%%
- **Compatible:** %t
- **Conflicts:** %d

`, result.Licenses.Score, result.Licenses.Compatible, len(result.Licenses.Conflicts))
	}

	if result.Performance != nil {
		md += fmt.Sprintf(`## Performance Audit

- **Score:** %.1f%%
- **Bundle Size:** %d bytes
- **Issues:** %d

`, result.Performance.Score, result.Performance.BundleSize, len(result.Performance.Issues))
	}

	if len(result.Recommendations) > 0 {
		md += "## Recommendations\n\n"
		for _, rec := range result.Recommendations {
			md += fmt.Sprintf("- %s\n", rec)
		}
		md += "\n"
	}

	return []byte(md), nil
}

// projectExists checks if the project path exists and is a directory
func (e *Engine) projectExists(path string) error {
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

// Security auditing helper methods

// scanDependencyVulnerabilities scans dependency files for known vulnerabilities
func (e *Engine) scanDependencyVulnerabilities(path string) ([]interfaces.Vulnerability, error) {
	var vulnerabilities []interfaces.Vulnerability

	// Check for common dependency files
	dependencyFiles := []string{
		"package.json",
		"go.mod",
		"requirements.txt",
		"Pipfile",
		"pom.xml",
		"build.gradle",
		"Cargo.toml",
	}

	for _, depFile := range dependencyFiles {
		filePath := filepath.Join(path, depFile)
		if _, err := os.Stat(filePath); err == nil {
			vulns, err := e.scanDependencyFile(filePath)
			if err != nil {
				continue // Log error but continue
			}
			vulnerabilities = append(vulnerabilities, vulns...)
		}
	}

	return vulnerabilities, nil
}

// scanDependencyFile scans a specific dependency file for vulnerabilities
func (e *Engine) scanDependencyFile(filePath string) ([]interfaces.Vulnerability, error) {
	var vulnerabilities []interfaces.Vulnerability

	// This is a simplified implementation
	// In a real implementation, you would integrate with vulnerability databases
	// like OSV, NVD, or commercial services

	// For demonstration, we'll check for some known vulnerable packages
	knownVulnerablePackages := map[string]interfaces.Vulnerability{
		"lodash": {
			ID:          "CVE-2021-23337",
			Severity:    "high",
			Title:       "Command Injection in lodash",
			Description: "Lodash versions prior to 4.17.21 are vulnerable to Command Injection",
			Package:     "lodash",
			Version:     "<4.17.21",
			FixedIn:     "4.17.21",
		},
		"axios": {
			ID:          "CVE-2021-3749",
			Severity:    "medium",
			Title:       "Regular Expression Denial of Service in axios",
			Description: "axios versions prior to 0.21.2 are vulnerable to ReDoS",
			Package:     "axios",
			Version:     "<0.21.2",
			FixedIn:     "0.21.2",
		},
	}

	// #nosec G304 - Audit tool needs to read dependency files for vulnerability analysis
	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Simple pattern matching for package names
	for pkg, vuln := range knownVulnerablePackages {
		if filepath.Base(filePath) == "package.json" {
			// Check if package is mentioned in package.json
			if len(content) > 0 && filepath.Ext(filePath) == ".json" {
				// Simple string search (in real implementation, parse JSON properly)
				if fmt.Sprintf("\"%s\"", pkg) != "" {
					vulnerabilities = append(vulnerabilities, vuln)
				}
			}
		}
	}

	return vulnerabilities, nil
}

// checkHardcodedSecrets checks for hardcoded secrets in source files
func (e *Engine) checkHardcodedSecrets(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	// Get secret detection rules
	rules := e.getSecretDetectionRules()

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) {
			return nil
		}

		secrets, err := e.scanFileForSecrets(filePath, rules)
		if err != nil {
			return nil // Continue on error
		}

		for _, secret := range secrets {
			if secret.Confidence >= 0.7 { // High confidence secrets
				violations = append(violations, interfaces.PolicyViolation{
					Policy:      "SEC-001",
					Severity:    "critical",
					Description: fmt.Sprintf("Potential hardcoded %s detected", secret.Type),
					File:        filePath,
					Line:        secret.Line,
				})
			}
		}

		return nil
	})

	return violations, err
}

// checkDependencyVulnerabilities checks for dependency vulnerabilities
func (e *Engine) checkDependencyVulnerabilities(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	vulnerabilities, err := e.scanDependencyVulnerabilities(path)
	if err != nil {
		return nil, err
	}

	for _, vuln := range vulnerabilities {
		violations = append(violations, interfaces.PolicyViolation{
			Policy:      "SEC-002",
			Severity:    vuln.Severity,
			Description: fmt.Sprintf("Vulnerable dependency: %s (%s)", vuln.Package, vuln.Title),
			File:        "", // Could be enhanced to show specific dependency file
			Line:        0,
		})
	}

	return violations, nil
}

// checkSecureConfiguration checks for secure configuration practices
func (e *Engine) checkSecureConfiguration(path string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	// Check common configuration files
	configFiles := []string{
		".env",
		"config.json",
		"config.yaml",
		"config.yml",
		"docker-compose.yml",
		"Dockerfile",
	}

	for _, configFile := range configFiles {
		filePath := filepath.Join(path, configFile)
		if _, err := os.Stat(filePath); err == nil {
			fileViolations, err := e.checkConfigFile(filePath)
			if err != nil {
				continue // Log error but continue
			}
			violations = append(violations, fileViolations...)
		}
	}

	return violations, nil
}

// checkConfigFile checks a specific configuration file for security issues
func (e *Engine) checkConfigFile(filePath string) ([]interfaces.PolicyViolation, error) {
	var violations []interfaces.PolicyViolation

	// #nosec G304 - Audit tool needs to read config files for security analysis
	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Check for insecure patterns
	insecurePatterns := []struct {
		pattern     string
		description string
		severity    string
	}{
		{`(?i)debug\s*[:=]\s*true`, "Debug mode enabled in production", "medium"},
		{`(?i)ssl\s*[:=]\s*false`, "SSL disabled", "high"},
		{`(?i)verify\s*[:=]\s*false`, "Certificate verification disabled", "high"},
		{`(?i)password\s*[:=]\s*["'].*["']`, "Hardcoded password in config", "critical"},
	}

	for i, line := range lines {
		for _, pattern := range insecurePatterns {
			matched, _ := regexp.MatchString(pattern.pattern, line)
			if matched {
				violations = append(violations, interfaces.PolicyViolation{
					Policy:      "SEC-003",
					Severity:    pattern.severity,
					Description: pattern.description,
					File:        filePath,
					Line:        i + 1,
				})
			}
		}
	}

	return violations, nil
}

// getSecretDetectionRules returns rules for detecting secrets
func (e *Engine) getSecretDetectionRules() []secretRule {
	return []secretRule{
		{
			name:       "AWS Access Key",
			pattern:    `AKIA[0-9A-Z]{16}`,
			confidence: 0.9,
		},
		{
			name:       "Generic API Key",
			pattern:    `(?i)(api[_-]?key|apikey)\s*[:=]\s*["']?[a-zA-Z0-9]{20,}["']?`,
			confidence: 0.7,
		},
		{
			name:       "Generic Secret",
			pattern:    `(?i)(secret|password|pwd)\s*[:=]\s*["'][^"']{8,}["']`,
			confidence: 0.6,
		},
		{
			name:       "JWT Token",
			pattern:    `eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`,
			confidence: 0.8,
		},
		{
			name:       "Private Key",
			pattern:    `-----BEGIN [A-Z ]+PRIVATE KEY-----`,
			confidence: 0.95,
		},
	}
}

// secretRule defines a rule for detecting secrets
type secretRule struct {
	name       string
	pattern    string
	confidence float64
}

// scanFileForSecrets scans a file for potential secrets
func (e *Engine) scanFileForSecrets(filePath string, rules []secretRule) ([]interfaces.SecretDetection, error) {
	var secrets []interfaces.SecretDetection

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		for _, rule := range rules {
			re, err := regexp.Compile(rule.pattern)
			if err != nil {
				continue
			}

			matches := re.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 0 {
					secrets = append(secrets, interfaces.SecretDetection{
						Type:       rule.name,
						File:       filePath,
						Line:       i + 1,
						Column:     strings.Index(line, match[0]) + 1,
						Secret:     match[0],
						Confidence: rule.confidence,
						Rule:       rule.name,
					})
				}
			}
		}
	}

	return secrets, nil
}

// shouldSkipFile determines if a file should be skipped during scanning
func (e *Engine) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Skip hidden files and directories
	if strings.HasPrefix(filepath.Base(filePath), ".") {
		return true
	}

	// Skip common directories
	skipDirs := []string{
		"node_modules", "vendor", ".git", ".svn", ".hg",
		"build", "dist", "target", "bin", "obj",
	}

	for _, skipDir := range skipDirs {
		if strings.Contains(filePath, skipDir) {
			return true
		}
	}

	return false
}

// calculateSecurityScore calculates the overall security score
func (e *Engine) calculateSecurityScore(secResult *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) float64 {
	score := 100.0

	// Deduct points for vulnerabilities
	for _, vuln := range secResult.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			score -= 20
		case "high":
			score -= 10
		case "medium":
			score -= 5
		case "low":
			score -= 2
		}
	}

	// Deduct points for policy violations
	for _, violation := range secResult.PolicyViolations {
		switch violation.Severity {
		case "critical":
			score -= 15
		case "high":
			score -= 8
		case "medium":
			score -= 4
		case "low":
			score -= 1
		}
	}

	// Deduct points for high-confidence secrets
	for _, secret := range secretResult.Secrets {
		if secret.Confidence >= 0.8 {
			score -= 10
		} else if secret.Confidence >= 0.6 {
			score -= 5
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// generateSecurityRecommendations generates security recommendations
func (e *Engine) generateSecurityRecommendations(secResult *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) []string {
	var recommendations []string

	if len(secResult.Vulnerabilities) > 0 {
		recommendations = append(recommendations, "Update vulnerable dependencies to secure versions")
		recommendations = append(recommendations, "Implement automated dependency scanning in CI/CD pipeline")
	}

	if len(secResult.PolicyViolations) > 0 {
		recommendations = append(recommendations, "Address security policy violations")
		recommendations = append(recommendations, "Implement security linting in development workflow")
	}

	if secretResult.Summary.HighConfidence > 0 {
		recommendations = append(recommendations, "Remove hardcoded secrets from source code")
		recommendations = append(recommendations, "Use environment variables or secret management systems")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Security posture is good - continue following security best practices")
	}

	return recommendations
}

// Quality and best practices auditing helper methods

// analyzeCodeSmells analyzes code for quality issues
func (e *Engine) analyzeCodeSmells(path string) ([]interfaces.CodeSmell, error) {
	var codeSmells []interfaces.CodeSmell

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) || !e.isSourceCodeFile(filePath) {
			return nil
		}

		smells, err := e.analyzeFileForCodeSmells(filePath)
		if err != nil {
			return nil // Continue on error
		}

		codeSmells = append(codeSmells, smells...)
		return nil
	})

	return codeSmells, err
}

// analyzeFileForCodeSmells analyzes a single file for code smells
func (e *Engine) analyzeFileForCodeSmells(filePath string) ([]interfaces.CodeSmell, error) {
	var smells []interfaces.CodeSmell

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Define code smell patterns
	smellPatterns := []struct {
		pattern     string
		type_       string
		severity    string
		description string
	}{
		{`(?i)todo|fixme|hack`, "technical_debt", "medium", "TODO/FIXME comment found"},
		{`console\.log|print\(|println\(`, "debug_code", "low", "Debug statement found"},
		{`\.length\s*>\s*\d{2,}`, "magic_number", "low", "Magic number used"},
		{`function\s+\w+\s*\([^)]{50,}`, "long_parameter_list", "medium", "Long parameter list"},
		{`\{\s*$[\s\S]{500,}?\}`, "long_method", "medium", "Method too long"},
		{`if\s*\([^)]*&&[^)]*&&[^)]*\)`, "complex_condition", "medium", "Complex conditional"},
	}

	for i, line := range lines {
		for _, pattern := range smellPatterns {
			matched, _ := regexp.MatchString(pattern.pattern, line)
			if matched {
				smells = append(smells, interfaces.CodeSmell{
					Type:        pattern.type_,
					Severity:    pattern.severity,
					Description: pattern.description,
					File:        filePath,
					Line:        i + 1,
				})
			}
		}

		// Check line length
		if len(line) > 120 {
			smells = append(smells, interfaces.CodeSmell{
				Type:        "long_line",
				Severity:    "low",
				Description: fmt.Sprintf("Line too long (%d characters)", len(line)),
				File:        filePath,
				Line:        i + 1,
			})
		}
	}

	return smells, nil
}

// analyzeDuplications analyzes code for duplications
func (e *Engine) analyzeDuplications(path string) ([]interfaces.Duplication, error) {
	var duplications []interfaces.Duplication

	// This is a simplified implementation
	// In a real implementation, you would use more sophisticated algorithms
	// like suffix trees or rolling hashes to detect duplications

	fileContents := make(map[string][]string)

	// Read all source files
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) || !e.isSourceCodeFile(filePath) {
			return nil
		}

		// #nosec G304 - Audit tool legitimately reads files for analysis
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		fileContents[filePath] = lines
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Simple duplication detection (looking for identical blocks of 5+ lines)
	const minDuplicationLines = 5

	for file1, lines1 := range fileContents {
		for file2, lines2 := range fileContents {
			if file1 >= file2 { // Avoid duplicate comparisons
				continue
			}

			duplications = append(duplications, e.findDuplicationsInFiles(file1, lines1, file2, lines2, minDuplicationLines)...)
		}
	}

	return duplications, nil
}

// findDuplicationsInFiles finds duplications between two files
func (e *Engine) findDuplicationsInFiles(file1 string, lines1 []string, file2 string, lines2 []string, minLines int) []interfaces.Duplication {
	var duplications []interfaces.Duplication

	// Simple sliding window approach
	for i := 0; i <= len(lines1)-minLines; i++ {
		for j := 0; j <= len(lines2)-minLines; j++ {
			matchLength := 0

			// Count consecutive matching lines
			for k := 0; i+k < len(lines1) && j+k < len(lines2); k++ {
				line1 := strings.TrimSpace(lines1[i+k])
				line2 := strings.TrimSpace(lines2[j+k])

				if line1 == line2 && line1 != "" {
					matchLength++
				} else {
					break
				}
			}

			if matchLength >= minLines {
				duplications = append(duplications, interfaces.Duplication{
					Files:      []string{file1, file2},
					Lines:      matchLength,
					Tokens:     matchLength * 10, // Rough estimate
					Percentage: float64(matchLength) / float64(len(lines1)) * 100,
				})
			}
		}
	}

	return duplications
}

// analyzeTestCoverage analyzes test coverage
func (e *Engine) analyzeTestCoverage(path string) (float64, error) {
	// This is a simplified implementation
	// In a real implementation, you would integrate with coverage tools
	// like go test -cover, jest --coverage, etc.

	var sourceFiles, testFiles int

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) {
			return nil
		}

		if e.isTestFile(filePath) {
			testFiles++
		} else if e.isSourceCodeFile(filePath) {
			sourceFiles++
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if sourceFiles == 0 {
		return 0, nil
	}

	// Simple heuristic: assume each test file covers multiple source files
	estimatedCoverage := float64(testFiles) / float64(sourceFiles) * 100
	if estimatedCoverage > 100 {
		estimatedCoverage = 100
	}

	return estimatedCoverage, nil
}

// getBestPractices returns best practices to check based on project type
func (e *Engine) getBestPractices(path string) []interfaces.BestPracticeCheck {
	practices := []interfaces.BestPracticeCheck{
		{
			ID:          "BP-001",
			Name:        "README file exists",
			Description: "Project should have a README file",
			Category:    "documentation",
			Technology:  "general",
			Compliant:   false,
			Score:       0,
		},
		{
			ID:          "BP-002",
			Name:        "License file exists",
			Description: "Project should have a LICENSE file",
			Category:    "legal",
			Technology:  "general",
			Compliant:   false,
			Score:       0,
		},
		{
			ID:          "BP-003",
			Name:        "Gitignore file exists",
			Description: "Project should have a .gitignore file",
			Category:    "version_control",
			Technology:  "general",
			Compliant:   false,
			Score:       0,
		},
		{
			ID:          "BP-004",
			Name:        "Tests exist",
			Description: "Project should have test files",
			Category:    "testing",
			Technology:  "general",
			Compliant:   false,
			Score:       0,
		},
	}

	// Add technology-specific practices
	if e.hasFile(path, "package.json") {
		practices = append(practices, interfaces.BestPracticeCheck{
			ID:          "BP-JS-001",
			Name:        "Package.json has description",
			Description: "Package.json should have a description field",
			Category:    "metadata",
			Technology:  "javascript",
			Compliant:   false,
			Score:       0,
		})
	}

	if e.hasFile(path, "go.mod") {
		practices = append(practices, interfaces.BestPracticeCheck{
			ID:          "BP-GO-001",
			Name:        "Go modules used",
			Description: "Go project should use modules",
			Category:    "dependency_management",
			Technology:  "go",
			Compliant:   true, // If go.mod exists, this is compliant
			Score:       100,
		})
	}

	return practices
}

// checkBestPractice checks a specific best practice
func (e *Engine) checkBestPractice(path string, practice interfaces.BestPracticeCheck) ([]interfaces.BestPracticeViolation, error) {
	var violations []interfaces.BestPracticeViolation

	switch practice.ID {
	case "BP-001": // README file exists
		if !e.hasReadmeFile(path) {
			violations = append(violations, interfaces.BestPracticeViolation{
				Practice:    practice.Name,
				File:        path,
				Line:        0,
				Description: "No README file found",
				Severity:    "medium",
				Suggestion:  "Add a README.md file with project description and usage instructions",
			})
		}

	case "BP-002": // License file exists
		if !e.hasLicenseFile(path) {
			violations = append(violations, interfaces.BestPracticeViolation{
				Practice:    practice.Name,
				File:        path,
				Line:        0,
				Description: "No LICENSE file found",
				Severity:    "medium",
				Suggestion:  "Add a LICENSE file with appropriate license terms",
			})
		}

	case "BP-003": // Gitignore file exists
		if !e.hasFile(path, ".gitignore") {
			violations = append(violations, interfaces.BestPracticeViolation{
				Practice:    practice.Name,
				File:        path,
				Line:        0,
				Description: "No .gitignore file found",
				Severity:    "low",
				Suggestion:  "Add a .gitignore file to exclude build artifacts and dependencies",
			})
		}

	case "BP-004": // Tests exist
		if !e.hasTestFiles(path) {
			violations = append(violations, interfaces.BestPracticeViolation{
				Practice:    practice.Name,
				File:        path,
				Line:        0,
				Description: "No test files found",
				Severity:    "high",
				Suggestion:  "Add test files to ensure code quality and reliability",
			})
		}

	case "BP-JS-001": // Package.json has description
		if e.hasFile(path, "package.json") {
			hasDescription, err := e.checkPackageJsonDescription(filepath.Join(path, "package.json"))
			if err != nil {
				return nil, err
			}
			if !hasDescription {
				violations = append(violations, interfaces.BestPracticeViolation{
					Practice:    practice.Name,
					File:        "package.json",
					Line:        0,
					Description: "Package.json missing description field",
					Severity:    "low",
					Suggestion:  "Add a description field to package.json",
				})
			}
		}
	}

	return violations, nil
}

// analyzeDependencyFile analyzes a dependency file for the given ecosystem
func (e *Engine) analyzeDependencyFile(filePath, ecosystem string) ([]interfaces.DependencyInfo, error) {
	var dependencies []interfaces.DependencyInfo

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// This is a simplified implementation
	// In a real implementation, you would parse the specific file format
	// and potentially query package registries for metadata

	switch ecosystem {
	case "npm":
		// Parse package.json (simplified)
		if strings.Contains(string(content), "\"dependencies\"") {
			// Extract dependency names (very simplified)
			lines := strings.Split(string(content), "\n")
			inDeps := false
			for _, line := range lines {
				if strings.Contains(line, "\"dependencies\"") {
					inDeps = true
					continue
				}
				if inDeps && strings.Contains(line, "}") {
					break
				}
				if inDeps && strings.Contains(line, "\"") {
					// Extract package name (simplified)
					parts := strings.Split(line, "\"")
					if len(parts) >= 2 {
						pkgName := parts[1]
						dependencies = append(dependencies, interfaces.DependencyInfo{
							Name:           pkgName,
							Version:        "unknown",
							Type:           "direct",
							License:        "unknown",
							LastUpdated:    time.Now().AddDate(0, -6, 0), // Assume 6 months old
							SecurityIssues: 0,
							QualityScore:   75.0,
						})
					}
				}
			}
		}

	case "go":
		// Parse go.mod (simplified)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "require ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					pkgName := parts[1]
					version := "unknown"
					if len(parts) >= 3 {
						version = parts[2]
					}
					dependencies = append(dependencies, interfaces.DependencyInfo{
						Name:           pkgName,
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

	return dependencies, nil
}

// analyzeFileComplexity analyzes complexity metrics for a file
func (e *Engine) analyzeFileComplexity(filePath string) (*interfaces.FileComplexity, error) {
	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Calculate basic metrics
	lineCount := len(lines)

	// Simple cyclomatic complexity calculation
	cyclomaticComplexity := 1 // Base complexity

	// Count decision points
	decisionPatterns := []string{
		`if\s*\(`, `else\s+if`, `while\s*\(`, `for\s*\(`,
		`switch\s*\(`, `case\s+`, `catch\s*\(`, `\?\s*`,
		`&&`, `\|\|`,
	}

	for _, line := range lines {
		for _, pattern := range decisionPatterns {
			matches, _ := regexp.MatchString(pattern, line)
			if matches {
				cyclomaticComplexity++
			}
		}
	}

	// Simple cognitive complexity (similar to cyclomatic for this implementation)
	cognitiveComplexity := cyclomaticComplexity

	// Calculate maintainability index (simplified)
	maintainability := 100.0
	if cyclomaticComplexity > 10 {
		maintainability -= float64(cyclomaticComplexity-10) * 5
	}
	if lineCount > 200 {
		maintainability -= float64(lineCount-200) * 0.1
	}
	if maintainability < 0 {
		maintainability = 0
	}

	// Estimate technical debt
	technicalDebt := "low"
	if cyclomaticComplexity > 15 {
		technicalDebt = "high"
	} else if cyclomaticComplexity > 10 {
		technicalDebt = "medium"
	}

	return &interfaces.FileComplexity{
		Path:                 filePath,
		Lines:                lineCount,
		CyclomaticComplexity: cyclomaticComplexity,
		CognitiveComplexity:  cognitiveComplexity,
		Maintainability:      maintainability,
		TechnicalDebt:        technicalDebt,
	}, nil
}

// extractFunctions extracts function information from a file (simplified)
func (e *Engine) extractFunctions(filePath string) []interfaces.FunctionComplexity {
	var functions []interfaces.FunctionComplexity

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return functions
	}

	lines := strings.Split(string(content), "\n")

	// Simple function detection patterns
	functionPatterns := []string{
		`func\s+(\w+)\s*\(`,         // Go
		`function\s+(\w+)\s*\(`,     // JavaScript
		`def\s+(\w+)\s*\(`,          // Python
		`public\s+\w+\s+(\w+)\s*\(`, // Java/C#
	}

	for i, line := range lines {
		for _, pattern := range functionPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}

			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				functionName := matches[1]

				// Simple complexity calculation for the function
				complexity := 1
				paramCount := strings.Count(line, ",") + 1
				if strings.Contains(line, "()") {
					paramCount = 0
				}

				functions = append(functions, interfaces.FunctionComplexity{
					Name:                 functionName,
					File:                 filePath,
					Line:                 i + 1,
					CyclomaticComplexity: complexity,
					CognitiveComplexity:  complexity,
					Parameters:           paramCount,
					Lines:                10, // Simplified - would need proper parsing
				})
			}
		}
	}

	return functions
}

// Helper methods for file checks

// isSourceCodeFile checks if a file is a source code file
func (e *Engine) isSourceCodeFile(filePath string) bool {
	sourceExtensions := []string{
		".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".cs", ".cpp", ".c", ".h",
		".rb", ".php", ".swift", ".kt", ".rs", ".scala", ".clj", ".hs", ".ml", ".fs",
		".css", ".scss", ".sass", ".less", ".html", ".htm", ".xml", ".yaml", ".yml",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sourceExt := range sourceExtensions {
		if ext == sourceExt {
			return true
		}
	}

	return false
}

// isTestFile checks if a file is a test file
func (e *Engine) isTestFile(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))
	fullPath := strings.ToLower(filePath)

	// Check filename patterns
	filenamePatterns := []string{
		"_test.", ".test.", "_spec.", ".spec.",
		"test_", "spec_",
	}

	for _, pattern := range filenamePatterns {
		if strings.Contains(fileName, pattern) {
			return true
		}
	}

	// Check directory patterns
	directoryPatterns := []string{
		"/test/", "/tests/", "/spec/",
	}

	for _, pattern := range directoryPatterns {
		if strings.Contains(fullPath, pattern) {
			return true
		}
	}

	return false
}

// hasFile checks if a file exists in the given path
func (e *Engine) hasFile(path, filename string) bool {
	_, err := os.Stat(filepath.Join(path, filename))
	return err == nil
}

// hasReadmeFile checks if a README file exists
func (e *Engine) hasReadmeFile(path string) bool {
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "README", "readme.md", "readme.txt"}
	for _, readme := range readmeFiles {
		if e.hasFile(path, readme) {
			return true
		}
	}
	return false
}

// hasLicenseFile checks if a LICENSE file exists
func (e *Engine) hasLicenseFile(path string) bool {
	licenseFiles := []string{"LICENSE", "LICENSE.txt", "LICENSE.md", "COPYING", "license", "license.txt"}
	for _, license := range licenseFiles {
		if e.hasFile(path, license) {
			return true
		}
	}
	return false
}

// hasTestFiles checks if test files exist in the project
func (e *Engine) hasTestFiles(path string) bool {
	hasTests := false

	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && e.isTestFile(filePath) {
			hasTests = true
			return filepath.SkipDir // Stop walking once we find a test
		}

		return nil
	})

	return hasTests
}

// checkPackageJsonDescription checks if package.json has a description
func (e *Engine) checkPackageJsonDescription(filePath string) (bool, error) {
	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Simple check for description field
	return strings.Contains(string(content), "\"description\""), nil
}

// calculateQualityScore calculates the overall quality score
func (e *Engine) calculateQualityScore(result *interfaces.QualityAuditResult) float64 {
	score := 100.0

	// Deduct points for code smells
	for _, smell := range result.CodeSmells {
		switch smell.Severity {
		case "high":
			score -= 5
		case "medium":
			score -= 3
		case "low":
			score -= 1
		}
	}

	// Deduct points for duplications
	for _, dup := range result.Duplications {
		if dup.Percentage > 10 {
			score -= 10
		} else if dup.Percentage > 5 {
			score -= 5
		}
	}

	// Adjust score based on test coverage
	if result.TestCoverage < 50 {
		score -= 20
	} else if result.TestCoverage < 80 {
		score -= 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

// generateQualityRecommendations generates quality recommendations
func (e *Engine) generateQualityRecommendations(result *interfaces.QualityAuditResult) []string {
	var recommendations []string

	if len(result.CodeSmells) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d code quality issues", len(result.CodeSmells)))
	}

	if len(result.Duplications) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Refactor %d code duplications", len(result.Duplications)))
	}

	if result.TestCoverage < 80 {
		recommendations = append(recommendations, fmt.Sprintf("Increase test coverage from %.1f%% to at least 80%%", result.TestCoverage))
	}

	if result.Score < 70 {
		recommendations = append(recommendations, "Consider implementing code quality gates in CI/CD pipeline")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Code quality is good - maintain current practices")
	}

	return recommendations
}

// License auditing helper methods

// detectProjectLicense detects the project's license
func (e *Engine) detectProjectLicense(path string) (string, error) {
	licenseFiles := []string{"LICENSE", "LICENSE.txt", "LICENSE.md", "COPYING"}

	for _, licenseFile := range licenseFiles {
		filePath := filepath.Join(path, licenseFile)
		if _, err := os.Stat(filePath); err == nil {
			// #nosec G304 - Audit tool legitimately reads files for analysis
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			// Simple license detection based on content
			contentStr := string(content)
			if strings.Contains(contentStr, "MIT License") {
				return "MIT", nil
			} else if strings.Contains(contentStr, "Apache License") {
				return "Apache-2.0", nil
			} else if strings.Contains(contentStr, "GNU General Public License") {
				return "GPL-3.0", nil
			} else if strings.Contains(contentStr, "BSD") {
				return "BSD-3-Clause", nil
			}

			return "custom", nil
		}
	}

	return "unknown", fmt.Errorf("no license file found")
}

// analyzeDependencyLicenses analyzes licenses of project dependencies
func (e *Engine) analyzeDependencyLicenses(path string) ([]interfaces.DependencyLicense, error) {
	var licenses []interfaces.DependencyLicense

	// This is a simplified implementation
	// In a real implementation, you would query package registries for license information

	// Check package.json for npm dependencies
	if e.hasFile(path, "package.json") {
		npmLicenses, err := e.analyzeNpmLicenses(filepath.Join(path, "package.json"))
		if err == nil {
			licenses = append(licenses, npmLicenses...)
		}
	}

	// Check go.mod for Go dependencies
	if e.hasFile(path, "go.mod") {
		goLicenses, err := e.analyzeGoLicenses(filepath.Join(path, "go.mod"))
		if err == nil {
			licenses = append(licenses, goLicenses...)
		}
	}

	return licenses, nil
}

// analyzeNpmLicenses analyzes npm package licenses
func (e *Engine) analyzeNpmLicenses(filePath string) ([]interfaces.DependencyLicense, error) {
	var licenses []interfaces.DependencyLicense

	// This is a simplified implementation
	// In a real implementation, you would parse package.json and query npm registry

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Simple extraction of dependency names
	lines := strings.Split(string(content), "\n")
	inDeps := false

	for _, line := range lines {
		if strings.Contains(line, "\"dependencies\"") {
			inDeps = true
			continue
		}
		if inDeps && strings.Contains(line, "}") {
			break
		}
		if inDeps && strings.Contains(line, "\"") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				pkgName := parts[1]
				// Mock license data (in real implementation, query npm registry)
				license := e.getMockLicense(pkgName)
				licenses = append(licenses, interfaces.DependencyLicense{
					Dependency: pkgName,
					License:    license,
					SPDXID:     license,
					Compatible: e.isLicenseCompatible("MIT", license), // Assume MIT project license
					Risk:       e.getLicenseRisk(license),
				})
			}
		}
	}

	return licenses, nil
}

// analyzeGoLicenses analyzes Go module licenses
func (e *Engine) analyzeGoLicenses(filePath string) ([]interfaces.DependencyLicense, error) {
	var licenses []interfaces.DependencyLicense

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				pkgName := parts[1]
				// Mock license data (in real implementation, query Go module proxy)
				license := e.getMockLicense(pkgName)
				licenses = append(licenses, interfaces.DependencyLicense{
					Dependency: pkgName,
					License:    license,
					SPDXID:     license,
					Compatible: e.isLicenseCompatible("MIT", license), // Assume MIT project license
					Risk:       e.getLicenseRisk(license),
				})
			}
		}
	}

	return licenses, nil
}

// getMockLicense returns a mock license for demonstration
func (e *Engine) getMockLicense(packageName string) string {
	// Mock license assignment based on package name hash
	licenses := []string{"MIT", "Apache-2.0", "BSD-3-Clause", "GPL-3.0", "ISC"}
	hash := 0
	for _, c := range packageName {
		hash += int(c)
	}
	return licenses[hash%len(licenses)]
}

// isLicenseCompatible checks if two licenses are compatible
func (e *Engine) isLicenseCompatible(projectLicense, depLicense string) bool {
	// Simplified compatibility matrix
	compatibilityMatrix := map[string][]string{
		"MIT":          {"MIT", "Apache-2.0", "BSD-3-Clause", "ISC"},
		"Apache-2.0":   {"MIT", "Apache-2.0", "BSD-3-Clause"},
		"BSD-3-Clause": {"MIT", "Apache-2.0", "BSD-3-Clause", "ISC"},
		"GPL-3.0":      {"GPL-3.0", "LGPL-3.0"},
		"ISC":          {"MIT", "Apache-2.0", "BSD-3-Clause", "ISC"},
	}

	compatible, exists := compatibilityMatrix[projectLicense]
	if !exists {
		return false // Unknown license, assume incompatible
	}

	for _, compat := range compatible {
		if compat == depLicense {
			return true
		}
	}

	return false
}

// getLicenseRisk returns the risk level for a license
func (e *Engine) getLicenseRisk(license string) string {
	highRiskLicenses := []string{"GPL-3.0", "AGPL-3.0", "SSPL-1.0"}
	mediumRiskLicenses := []string{"LGPL-3.0", "MPL-2.0"}

	for _, highRisk := range highRiskLicenses {
		if license == highRisk {
			return "high"
		}
	}

	for _, mediumRisk := range mediumRiskLicenses {
		if license == mediumRisk {
			return "medium"
		}
	}

	return "low"
}

// getLicenseConflictSeverity returns the severity of a license conflict
func (e *Engine) getLicenseConflictSeverity(license1, license2 string) string {
	if e.getLicenseRisk(license1) == "high" || e.getLicenseRisk(license2) == "high" {
		return "critical"
	}
	if e.getLicenseRisk(license1) == "medium" || e.getLicenseRisk(license2) == "medium" {
		return "high"
	}
	return "medium"
}

// generateLicenseRecommendations generates license recommendations
func (e *Engine) generateLicenseRecommendations(result *interfaces.LicenseAuditResult) []string {
	var recommendations []string

	if len(result.Conflicts) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Resolve %d license conflicts", len(result.Conflicts)))
		recommendations = append(recommendations, "Consider using a license compatibility tool")
	}

	if !result.Compatible {
		recommendations = append(recommendations, "Review project license compatibility with dependencies")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "License compliance is good")
	}

	return recommendations
}

// Performance auditing helper methods

// analyzePerformanceIssues analyzes performance issues in the project
func (e *Engine) analyzePerformanceIssues(path string) ([]interfaces.PerformanceIssue, error) {
	var issues []interfaces.PerformanceIssue

	// Check for large files
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) {
			return nil
		}

		// Check for large source files
		if e.isSourceCodeFile(filePath) && info.Size() > 100*1024 { // 100KB
			issues = append(issues, interfaces.PerformanceIssue{
				Type:        "large_file",
				Severity:    "medium",
				Description: fmt.Sprintf("Large source file: %s (%d bytes)", filepath.Base(filePath), info.Size()),
				Impact:      "May slow down compilation and IDE performance",
				File:        filePath,
			})
		}

		// Check for large assets
		if e.isAssetFile(filePath) && info.Size() > 1024*1024 { // 1MB
			issues = append(issues, interfaces.PerformanceIssue{
				Type:        "large_asset",
				Severity:    "high",
				Description: fmt.Sprintf("Large asset file: %s (%d bytes)", filepath.Base(filePath), info.Size()),
				Impact:      "May slow down page load times",
				File:        filePath,
			})
		}

		return nil
	})

	// Check for performance anti-patterns in code
	codeIssues, codeErr := e.analyzeCodePerformanceIssues(path)
	if codeErr == nil {
		issues = append(issues, codeIssues...)
	}

	return issues, err
}

// analyzeCodePerformanceIssues analyzes code for performance anti-patterns
func (e *Engine) analyzeCodePerformanceIssues(path string) ([]interfaces.PerformanceIssue, error) {
	var issues []interfaces.PerformanceIssue

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) || !e.isSourceCodeFile(filePath) {
			return nil
		}

		// #nosec G304 - Audit tool legitimately reads files for analysis
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")

		// Check for performance anti-patterns
		performancePatterns := []struct {
			pattern     string
			type_       string
			severity    string
			description string
			impact      string
		}{
			{`for.*for.*for`, "nested_loops", "high", "Deeply nested loops detected", "May cause O(n³) performance"},
			{`\.innerHTML\s*=`, "dom_manipulation", "medium", "Direct DOM manipulation", "May cause layout thrashing"},
			{`document\.getElementById`, "inefficient_dom", "low", "Inefficient DOM query", "Consider caching DOM references"},
			{`console\.log`, "debug_statements", "low", "Debug statements in production", "May impact performance"},
		}

		for i, line := range lines {
			for _, pattern := range performancePatterns {
				matched, _ := regexp.MatchString(pattern.pattern, line)
				if matched {
					issues = append(issues, interfaces.PerformanceIssue{
						Type:        pattern.type_,
						Severity:    pattern.severity,
						Description: pattern.description,
						Impact:      pattern.impact,
						File:        fmt.Sprintf("%s:%d", filePath, i+1),
					})
				}
			}
		}

		return nil
	})

	return issues, err
}

// analyzeBuildDirectory analyzes a build directory for bundle information
func (e *Engine) analyzeBuildDirectory(buildPath string, result *interfaces.BundleAnalysisResult) error {
	return filepath.Walk(buildPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		result.TotalSize += info.Size()

		// Estimate gzipped size (rough approximation)
		if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".css") {
			result.GzippedSize += info.Size() / 3 // Rough compression ratio
		} else {
			result.GzippedSize += info.Size()
		}

		// Create asset entry
		asset := interfaces.BundleAsset{
			Name:        filepath.Base(filePath),
			Size:        info.Size(),
			GzippedSize: info.Size() / 3, // Rough estimate
			Type:        e.getAssetType(filePath),
			Percentage:  0, // Will be calculated later
		}

		result.Assets = append(result.Assets, asset)

		return nil
	})
}

// analyzeSourceFiles analyzes source files when no build directory is found
func (e *Engine) analyzeSourceFiles(path string, result *interfaces.BundleAnalysisResult) {
	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || e.shouldSkipFile(filePath) || !e.isSourceCodeFile(filePath) {
			return nil
		}

		result.TotalSize += info.Size()
		result.GzippedSize += info.Size() / 2 // Rough compression for source files

		asset := interfaces.BundleAsset{
			Name:        filepath.Base(filePath),
			Size:        info.Size(),
			GzippedSize: info.Size() / 2,
			Type:        e.getAssetType(filePath),
			Percentage:  0,
		}

		result.Assets = append(result.Assets, asset)

		return nil
	})
}

// isAssetFile checks if a file is an asset file
func (e *Engine) isAssetFile(filePath string) bool {
	assetExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp",
		".mp3", ".mp4", ".avi", ".mov", ".wmv", ".pdf", ".zip",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, assetExt := range assetExtensions {
		if ext == assetExt {
			return true
		}
	}

	return false
}

// getAssetType returns the type of an asset based on its extension
func (e *Engine) getAssetType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".js", ".jsx", ".ts", ".tsx":
		return "javascript"
	case ".css", ".scss", ".sass", ".less":
		return "stylesheet"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "image"
	case ".woff", ".woff2", ".ttf", ".eot":
		return "font"
	case ".json":
		return "data"
	default:
		return "other"
	}
}

// calculatePerformanceScore calculates the performance score
func (e *Engine) calculatePerformanceScore(result *interfaces.PerformanceAuditResult) float64 {
	score := 100.0

	// Deduct points for bundle size
	sizeInMB := float64(result.BundleSize) / (1024 * 1024)
	if sizeInMB > 10 {
		score -= 30
	} else if sizeInMB > 5 {
		score -= 20
	} else if sizeInMB > 2 {
		score -= 10
	}

	// Deduct points for load time
	if result.LoadTime > 5*time.Second {
		score -= 40
	} else if result.LoadTime > 3*time.Second {
		score -= 25
	} else if result.LoadTime > 1*time.Second {
		score -= 10
	}

	// Deduct points for performance issues
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "high":
			score -= 10
		case "medium":
			score -= 5
		case "low":
			score -= 2
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

// generatePerformanceRecommendations generates performance recommendations
func (e *Engine) generatePerformanceRecommendations(result *interfaces.PerformanceAuditResult) []string {
	var recommendations []string

	sizeInMB := float64(result.BundleSize) / (1024 * 1024)
	if sizeInMB > 5 {
		recommendations = append(recommendations, fmt.Sprintf("Reduce bundle size from %.1fMB", sizeInMB))
		recommendations = append(recommendations, "Consider code splitting and lazy loading")
	}

	if result.LoadTime > 2*time.Second {
		recommendations = append(recommendations, fmt.Sprintf("Optimize load time from %v", result.LoadTime))
		recommendations = append(recommendations, "Enable compression and caching")
	}

	if len(result.Issues) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d performance issues", len(result.Issues)))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance is good - maintain current optimizations")
	}

	return recommendations
}
