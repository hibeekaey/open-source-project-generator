package audit

import (
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResultAggregator(t *testing.T) {
	ra := NewResultAggregator()

	assert.NotNil(t, ra)
	assert.Equal(t, 0.3, ra.securityWeight)
	assert.Equal(t, 0.3, ra.qualityWeight)
	assert.Equal(t, 0.2, ra.licenseWeight)
	assert.Equal(t, 0.2, ra.performanceWeight)
}

func TestNewResultAggregatorWithWeights(t *testing.T) {
	ra := NewResultAggregatorWithWeights(0.4, 0.3, 0.2, 0.1)

	assert.NotNil(t, ra)
	assert.Equal(t, 0.4, ra.securityWeight)
	assert.Equal(t, 0.3, ra.qualityWeight)
	assert.Equal(t, 0.2, ra.licenseWeight)
	assert.Equal(t, 0.1, ra.performanceWeight)
}

func TestAggregateResults_EmptyResults(t *testing.T) {
	ra := NewResultAggregator()

	summary, err := ra.AggregateResults([]*interfaces.AuditResult{})

	require.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 0, summary.TotalProjects)
}

func TestAggregateResults_SingleResult(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Now(),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score: 90.0,
			Vulnerabilities: []interfaces.Vulnerability{
				{Severity: "high", Title: "Test vulnerability"},
			},
		},
		Quality: &interfaces.QualityAuditResult{
			Score: 80.0,
			CodeSmells: []interfaces.CodeSmell{
				{Type: "complexity", Severity: "medium"},
			},
		},
	}

	summary, err := ra.AggregateResults([]*interfaces.AuditResult{result})

	require.NoError(t, err)
	assert.Equal(t, 1, summary.TotalProjects)
	assert.Equal(t, 85.0, summary.AverageScore)
	assert.Equal(t, 90.0, summary.SecurityScore)
	assert.Equal(t, 80.0, summary.QualityScore)
	// With single project, no issues should be considered "common" (need >25% frequency)
	assert.Len(t, summary.CommonIssues, 0)
}

func TestAggregateResults_MultipleResults(t *testing.T) {
	ra := NewResultAggregator()

	results := []*interfaces.AuditResult{
		{
			ProjectPath:  "/test/project1",
			OverallScore: 85.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 90.0,
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high", Title: "Test vulnerability"},
				},
			},
		},
		{
			ProjectPath:  "/test/project2",
			OverallScore: 75.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 80.0,
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high", Title: "Another vulnerability"},
				},
			},
		},
	}

	summary, err := ra.AggregateResults(results)

	require.NoError(t, err)
	assert.Equal(t, 2, summary.TotalProjects)
	assert.Equal(t, 80.0, summary.AverageScore)  // (85 + 75) / 2
	assert.Equal(t, 85.0, summary.SecurityScore) // (90 + 80) / 2
	assert.Len(t, summary.CommonIssues, 1)       // High security issues appear in both projects
}

func TestCalculateOverallScore(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		Security:    &interfaces.SecurityAuditResult{Score: 90.0},
		Quality:     &interfaces.QualityAuditResult{Score: 80.0},
		Licenses:    &interfaces.LicenseAuditResult{Score: 85.0},
		Performance: &interfaces.PerformanceAuditResult{Score: 75.0},
	}

	score := ra.CalculateOverallScore(result)

	// Expected: (90*0.3 + 80*0.3 + 85*0.2 + 75*0.2) = 27 + 24 + 17 + 15 = 83
	assert.Equal(t, 83.0, score)
}

func TestCalculateOverallScore_PartialResults(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		Security: &interfaces.SecurityAuditResult{Score: 90.0},
		Quality:  &interfaces.QualityAuditResult{Score: 80.0},
		// No license or performance results
	}

	score := ra.CalculateOverallScore(result)

	// Expected: (90*0.3 + 80*0.3) / (0.3 + 0.3) = 51 / 0.6 = 85
	assert.Equal(t, 85.0, score)
}

func TestCalculateSecurityScore(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		name          string
		secResult     *interfaces.SecurityAuditResult
		secretResult  *interfaces.SecretScanResult
		expectedScore float64
	}{
		{
			name: "perfect security",
			secResult: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult:  &interfaces.SecretScanResult{Secrets: []interfaces.SecretDetection{}},
			expectedScore: 100.0,
		},
		{
			name: "critical vulnerability",
			secResult: &interfaces.SecurityAuditResult{
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "critical"},
				},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult:  &interfaces.SecretScanResult{Secrets: []interfaces.SecretDetection{}},
			expectedScore: 80.0, // 100 - 20
		},
		{
			name: "high confidence secret",
			secResult: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Secrets: []interfaces.SecretDetection{
					{Confidence: 0.9},
				},
			},
			expectedScore: 90.0, // 100 - 10
		},
		{
			name: "multiple issues",
			secResult: &interfaces.SecurityAuditResult{
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high"},
					{Severity: "medium"},
				},
				PolicyViolations: []interfaces.PolicyViolation{
					{Severity: "critical"},
				},
			},
			secretResult: &interfaces.SecretScanResult{
				Secrets: []interfaces.SecretDetection{
					{Confidence: 0.8},
				},
			},
			expectedScore: 60.0, // 100 - 10 - 5 - 15 - 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := ra.CalculateSecurityScore(tt.secResult, tt.secretResult)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestCalculateQualityScore(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		name          string
		result        *interfaces.QualityAuditResult
		expectedScore float64
	}{
		{
			name: "perfect quality",
			result: &interfaces.QualityAuditResult{
				CodeSmells:   []interfaces.CodeSmell{},
				Duplications: []interfaces.Duplication{},
				TestCoverage: 90.0,
			},
			expectedScore: 100.0,
		},
		{
			name: "low test coverage",
			result: &interfaces.QualityAuditResult{
				CodeSmells:   []interfaces.CodeSmell{},
				Duplications: []interfaces.Duplication{},
				TestCoverage: 40.0,
			},
			expectedScore: 80.0, // 100 - 20
		},
		{
			name: "high duplication",
			result: &interfaces.QualityAuditResult{
				CodeSmells: []interfaces.CodeSmell{},
				Duplications: []interfaces.Duplication{
					{Percentage: 15.0},
				},
				TestCoverage: 90.0,
			},
			expectedScore: 90.0, // 100 - 10
		},
		{
			name: "multiple issues",
			result: &interfaces.QualityAuditResult{
				CodeSmells: []interfaces.CodeSmell{
					{Severity: "high"},
					{Severity: "medium"},
				},
				Duplications: []interfaces.Duplication{
					{Percentage: 8.0},
				},
				TestCoverage: 70.0,
			},
			expectedScore: 77.0, // 100 - 5 - 3 - 5 - 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := ra.CalculateQualityScore(tt.result)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestCalculatePerformanceScore(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		name          string
		result        *interfaces.PerformanceAuditResult
		expectedScore float64
	}{
		{
			name: "perfect performance",
			result: &interfaces.PerformanceAuditResult{
				BundleSize: 1024 * 1024, // 1MB
				LoadTime:   500 * time.Millisecond,
				Issues:     []interfaces.PerformanceIssue{},
			},
			expectedScore: 100.0,
		},
		{
			name: "large bundle",
			result: &interfaces.PerformanceAuditResult{
				BundleSize: 12 * 1024 * 1024, // 12MB
				LoadTime:   500 * time.Millisecond,
				Issues:     []interfaces.PerformanceIssue{},
			},
			expectedScore: 70.0, // 100 - 30
		},
		{
			name: "slow load time",
			result: &interfaces.PerformanceAuditResult{
				BundleSize: 1024 * 1024, // 1MB
				LoadTime:   6 * time.Second,
				Issues:     []interfaces.PerformanceIssue{},
			},
			expectedScore: 60.0, // 100 - 40
		},
		{
			name: "performance issues",
			result: &interfaces.PerformanceAuditResult{
				BundleSize: 1024 * 1024, // 1MB
				LoadTime:   500 * time.Millisecond,
				Issues: []interfaces.PerformanceIssue{
					{Severity: "high"},
					{Severity: "medium"},
				},
			},
			expectedScore: 85.0, // 100 - 10 - 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := ra.CalculatePerformanceScore(tt.result)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestGenerateRecommendations(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		OverallScore: 65.0, // Below 70%
		Security: &interfaces.SecurityAuditResult{
			Score: 60.0,
			Vulnerabilities: []interfaces.Vulnerability{
				{Severity: "high"},
			},
		},
		Quality: &interfaces.QualityAuditResult{
			Score:        70.0,
			TestCoverage: 60.0, // Below 80%
		},
	}

	recommendations := ra.GenerateRecommendations(result)

	assert.NotEmpty(t, recommendations)

	// Should contain overall recommendation
	found := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "Overall project health is below 70%") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain overall health recommendation")
}

func TestGenerateSecurityRecommendations(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		name         string
		secResult    *interfaces.SecurityAuditResult
		secretResult *interfaces.SecretScanResult
		expectCount  int
	}{
		{
			name: "no issues",
			secResult: &interfaces.SecurityAuditResult{
				Score:            100.0,
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{HighConfidence: 0},
			},
			expectCount: 1, // Good security message
		},
		{
			name: "vulnerabilities found",
			secResult: &interfaces.SecurityAuditResult{
				Score: 80.0,
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high"},
				},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{HighConfidence: 0},
			},
			expectCount: 2, // Update dependencies + CI/CD
		},
		{
			name: "secrets found",
			secResult: &interfaces.SecurityAuditResult{
				Score:            90.0,
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{HighConfidence: 2},
			},
			expectCount: 2, // Remove secrets + use env vars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendations := ra.GenerateSecurityRecommendations(tt.secResult, tt.secretResult)
			assert.Len(t, recommendations, tt.expectCount)
		})
	}
}

func TestGenerateHTMLReport(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		OverallScore: 85.0,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{{Severity: "medium"}},
		},
		Recommendations: []string{"Test recommendation"},
	}

	htmlBytes, err := ra.GenerateHTMLReport(result)

	require.NoError(t, err)
	html := string(htmlBytes)

	assert.Contains(t, html, "<!DOCTYPE html>")
	assert.Contains(t, html, "Audit Report")
	assert.Contains(t, html, "project")
	assert.Contains(t, html, "85.0%")
	assert.Contains(t, html, "Security Audit")
	assert.Contains(t, html, "90.0%")
	assert.Contains(t, html, "Test recommendation")
}

func TestGenerateMarkdownReport(t *testing.T) {
	ra := NewResultAggregator()

	result := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		OverallScore: 85.0,
		Quality: &interfaces.QualityAuditResult{
			Score:        80.0,
			TestCoverage: 75.0,
			CodeSmells:   []interfaces.CodeSmell{{Type: "complexity"}},
		},
		Recommendations: []string{"Test recommendation"},
	}

	mdBytes, err := ra.GenerateMarkdownReport(result)

	require.NoError(t, err)
	md := string(mdBytes)

	assert.Contains(t, md, "# Audit Report")
	assert.Contains(t, md, "/test/project")
	assert.Contains(t, md, "85.0%")
	assert.Contains(t, md, "## Quality Audit")
	assert.Contains(t, md, "80.0%")
	assert.Contains(t, md, "## Recommendations")
	assert.Contains(t, md, "Test recommendation")
}

func TestGetIssueDescription(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		issueType    string
		expectedDesc string
	}{
		{"security-critical", "Critical security vulnerabilities found"},
		{"quality-high", "High severity code quality issues found"},
		{"license-conflict", "License compatibility conflicts found"},
		{"performance-large_file", "Large files affecting performance found"},
		{"unknown-type", "Issue type: unknown-type"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			desc := ra.getIssueDescription(tt.issueType)
			assert.Equal(t, tt.expectedDesc, desc)
		})
	}
}

func TestGetIssueSeverity(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		issueType        string
		expectedSeverity string
	}{
		{"security-critical", "critical"},
		{"security-high", "high"},
		{"quality-medium", "medium"},
		{"license-conflict", "medium"},
		{"performance-large_file", "medium"},
		{"unknown-type", "low"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			severity := ra.getIssueSeverity(tt.issueType)
			assert.Equal(t, tt.expectedSeverity, severity)
		})
	}
}

func TestGetIssueCategory(t *testing.T) {
	ra := NewResultAggregator()

	tests := []struct {
		issueType        string
		expectedCategory string
	}{
		{"security-critical", "security"},
		{"quality-high", "quality"},
		{"license-conflict", "license"},
		{"performance-large_file", "performance"},
		{"unknown-type", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			category := ra.getIssueCategory(tt.issueType)
			assert.Equal(t, tt.expectedCategory, category)
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkCalculateOverallScore(b *testing.B) {
	ra := NewResultAggregator()
	result := &interfaces.AuditResult{
		Security:    &interfaces.SecurityAuditResult{Score: 90.0},
		Quality:     &interfaces.QualityAuditResult{Score: 80.0},
		Licenses:    &interfaces.LicenseAuditResult{Score: 85.0},
		Performance: &interfaces.PerformanceAuditResult{Score: 75.0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ra.CalculateOverallScore(result)
	}
}

func BenchmarkGenerateRecommendations(b *testing.B) {
	ra := NewResultAggregator()
	result := &interfaces.AuditResult{
		OverallScore: 65.0,
		Security: &interfaces.SecurityAuditResult{
			Score: 60.0,
			Vulnerabilities: []interfaces.Vulnerability{
				{Severity: "high"}, {Severity: "medium"}, {Severity: "low"},
			},
		},
		Quality: &interfaces.QualityAuditResult{
			Score:        70.0,
			TestCoverage: 60.0,
			CodeSmells: []interfaces.CodeSmell{
				{Type: "complexity", Severity: "high"},
				{Type: "duplication", Severity: "medium"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ra.GenerateRecommendations(result)
	}
}
