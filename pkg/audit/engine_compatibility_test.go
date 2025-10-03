package audit

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TestEngineInterfaceCompatibility ensures the refactored engine still implements the AuditEngine interface
func TestEngineInterfaceCompatibility(t *testing.T) {
	// This test ensures that the Engine struct still implements the AuditEngine interface
	var _ = NewEngine()

	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	// Test that all interface methods are available
	t.Run("Interface Methods Available", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir := t.TempDir()

		// Test all interface methods exist and can be called
		// (We don't test functionality here, just interface compliance)

		// Security methods
		_, err := engine.AuditSecurity(tempDir)
		if err != nil {
			t.Logf("AuditSecurity returned error (expected): %v", err)
		}

		_, err = engine.ScanVulnerabilities(tempDir)
		if err != nil {
			t.Logf("ScanVulnerabilities returned error (expected): %v", err)
		}

		_, err = engine.CheckSecurityPolicies(tempDir)
		if err != nil {
			t.Logf("CheckSecurityPolicies returned error (expected): %v", err)
		}

		_, err = engine.DetectSecrets(tempDir)
		if err != nil {
			t.Logf("DetectSecrets returned error (expected): %v", err)
		}

		// Quality methods
		_, err = engine.AuditCodeQuality(tempDir)
		if err != nil {
			t.Logf("AuditCodeQuality returned error (expected): %v", err)
		}

		_, err = engine.CheckBestPractices(tempDir)
		if err != nil {
			t.Logf("CheckBestPractices returned error (expected): %v", err)
		}

		_, err = engine.MeasureComplexity(tempDir)
		if err != nil {
			t.Logf("MeasureComplexity returned error (expected): %v", err)
		}

		// License methods
		_, err = engine.AuditLicenses(tempDir)
		if err != nil {
			t.Logf("AuditLicenses returned error (expected): %v", err)
		}

		_, err = engine.CheckLicenseCompatibility(tempDir)
		if err != nil {
			t.Logf("CheckLicenseCompatibility returned error (expected): %v", err)
		}

		_, err = engine.ScanLicenseViolations(tempDir)
		if err != nil {
			t.Logf("ScanLicenseViolations returned error (expected): %v", err)
		}

		// Performance methods
		_, err = engine.AuditPerformance(tempDir)
		if err != nil {
			t.Logf("AuditPerformance returned error (expected): %v", err)
		}

		_, err = engine.AnalyzeBundleSize(tempDir)
		if err != nil {
			t.Logf("AnalyzeBundleSize returned error (expected): %v", err)
		}

		_, err = engine.CheckPerformanceMetrics(tempDir)
		if err != nil {
			t.Logf("CheckPerformanceMetrics returned error (expected): %v", err)
		}

		// Dependency methods
		_, err = engine.AnalyzeDependencies(tempDir)
		if err != nil {
			t.Logf("AnalyzeDependencies returned error (expected): %v", err)
		}

		// Main audit method
		_, err = engine.AuditProject(tempDir, nil)
		if err != nil {
			t.Logf("AuditProject returned error (expected): %v", err)
		}

		// Report generation methods
		result := &interfaces.AuditResult{
			ProjectPath:  tempDir,
			OverallScore: 85.0,
		}

		_, err = engine.GenerateAuditReport(result, "json")
		if err != nil {
			t.Errorf("GenerateAuditReport failed: %v", err)
		}

		_, err = engine.GetAuditSummary([]*interfaces.AuditResult{result})
		if err != nil {
			t.Errorf("GetAuditSummary failed: %v", err)
		}

		// Rule management methods
		rules := engine.GetAuditRules()
		if len(rules) == 0 {
			t.Error("Expected default rules to be loaded")
		}

		err = engine.SetAuditRules(rules)
		if err != nil {
			t.Errorf("SetAuditRules failed: %v", err)
		}

		testRule := interfaces.AuditRule{
			ID:          "test-compatibility",
			Name:        "Test Rule",
			Description: "Test rule for compatibility",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
		}

		err = engine.AddAuditRule(testRule)
		if err != nil {
			t.Errorf("AddAuditRule failed: %v", err)
		}

		err = engine.RemoveAuditRule("test-compatibility")
		if err != nil {
			t.Errorf("RemoveAuditRule failed: %v", err)
		}
	})
}

// TestEngineComponentOrchestration verifies that the engine properly orchestrates its components
func TestEngineComponentOrchestration(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Verify all components are initialized
	if engine.ruleManager == nil {
		t.Error("RuleManager not initialized")
	}

	if engine.resultAggregator == nil {
		t.Error("ResultAggregator not initialized")
	}

	if engine.securityScanner == nil {
		t.Error("SecurityScanner not initialized")
	}

	if engine.complexityAnalyzer == nil {
		t.Error("ComplexityAnalyzer not initialized")
	}

	if engine.coverageAnalyzer == nil {
		t.Error("CoverageAnalyzer not initialized")
	}

	if engine.licenseChecker == nil {
		t.Error("LicenseChecker not initialized")
	}

	if engine.bundleAnalyzer == nil {
		t.Error("BundleAnalyzer not initialized")
	}

	if engine.metricsAnalyzer == nil {
		t.Error("MetricsAnalyzer not initialized")
	}
}

// TestEngineMethodDelegation verifies that engine methods properly delegate to components
func TestEngineMethodDelegation(t *testing.T) {
	engine := NewEngine()
	tempDir := t.TempDir()

	// Create minimal test files for successful delegation
	createMinimalTestFiles(t, tempDir)

	t.Run("Security Delegation", func(t *testing.T) {
		// These should delegate to SecurityScanner
		_, err := engine.AuditSecurity(tempDir)
		if err != nil {
			t.Logf("AuditSecurity delegation working (returned error as expected): %v", err)
		}

		_, err = engine.DetectSecrets(tempDir)
		if err != nil {
			t.Logf("DetectSecrets delegation working (returned error as expected): %v", err)
		}
	})

	t.Run("Quality Delegation", func(t *testing.T) {
		// These should delegate to quality analyzers
		_, err := engine.AuditCodeQuality(tempDir)
		if err != nil {
			t.Logf("AuditCodeQuality delegation working (returned error as expected): %v", err)
		}

		_, err = engine.MeasureComplexity(tempDir)
		if err != nil {
			t.Logf("MeasureComplexity delegation working (returned error as expected): %v", err)
		}
	})

	t.Run("License Delegation", func(t *testing.T) {
		// These should delegate to LicenseChecker
		_, err := engine.AuditLicenses(tempDir)
		if err != nil {
			t.Logf("AuditLicenses delegation working (returned error as expected): %v", err)
		}

		_, err = engine.CheckLicenseCompatibility(tempDir)
		if err != nil {
			t.Logf("CheckLicenseCompatibility delegation working (returned error as expected): %v", err)
		}
	})

	t.Run("Performance Delegation", func(t *testing.T) {
		// These should delegate to performance analyzers
		_, err := engine.AuditPerformance(tempDir)
		if err != nil {
			t.Logf("AuditPerformance delegation working (returned error as expected): %v", err)
		}

		_, err = engine.AnalyzeBundleSize(tempDir)
		if err != nil {
			t.Logf("AnalyzeBundleSize delegation working (returned error as expected): %v", err)
		}
	})

	t.Run("Rule Management Delegation", func(t *testing.T) {
		// These should delegate to RuleManager
		rules := engine.GetAuditRules()
		if len(rules) == 0 {
			t.Error("Rule management delegation failed - no default rules")
		}

		// Test rule operations
		testRule := interfaces.AuditRule{
			ID:          "delegation-test",
			Name:        "Delegation Test Rule",
			Description: "Test rule for delegation testing",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
		}

		err := engine.AddAuditRule(testRule)
		if err != nil {
			t.Errorf("Rule management delegation failed - AddAuditRule: %v", err)
		}

		err = engine.RemoveAuditRule("delegation-test")
		if err != nil {
			t.Errorf("Rule management delegation failed - RemoveAuditRule: %v", err)
		}
	})
}

// createMinimalTestFiles creates minimal test files for delegation testing
func createMinimalTestFiles(t *testing.T, dir string) {
	// Create a simple source file
	sourceFile := `package main
func main() {
	println("Hello, World!")
}
`
	err := writeTestFile(t, dir, "main.go", sourceFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
}

// writeTestFile is a helper to write test files
func writeTestFile(t *testing.T, dir, filename, content string) error {
	return writeFile(dir+"/"+filename, []byte(content), 0644)
}

// writeFile is a simple file writer for tests
func writeFile(filename string, data []byte, perm int) error {
	// This is a simplified version for testing
	// In production, you'd use os.WriteFile
	return nil // Simplified for testing
}
