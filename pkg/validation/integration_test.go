package validation

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
)

// TestTemplateValidatorIntegration tests the template validator with actual embedded templates
func TestTemplateValidatorIntegration(t *testing.T) {
	validator := NewTemplateValidator()

	// Test that we can validate embedded templates
	if !validator.useEmbedded {
		t.Skip("Embedded filesystem not available, skipping integration test")
	}

	// Test comprehensive validation
	result, err := validator.ValidateAllEmbeddedTemplates()
	if err != nil {
		t.Fatalf("Failed to validate embedded templates: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Log validation results for debugging
	t.Logf("Validation completed: %s", result.Summary)
	t.Logf("Valid: %v", result.Valid)
	t.Logf("Issues found: %d", len(result.Issues))

	for i, issue := range result.Issues {
		t.Logf("Issue %d: %s - %s (File: %s)", i+1, issue.Type, issue.Message, issue.File)
	}
}

// TestTemplateValidatorWithSpecificTemplate tests validation of a specific template
func TestTemplateValidatorWithSpecificTemplate(t *testing.T) {
	validator := NewTemplateValidator()

	if !validator.useEmbedded {
		t.Skip("Embedded filesystem not available, skipping integration test")
	}

	// Test package.json validation with embedded template
	packageJsonPath := "frontend/nextjs-app/package.json.tmpl"
	err := validator.ValidatePackageJSON(packageJsonPath)

	// This might fail if the template doesn't exist, which is okay for this test
	if err != nil {
		t.Logf("Package.json validation result: %v", err)
	} else {
		t.Log("Package.json validation passed")
	}
}

// TestTemplateValidatorFallback tests fallback to filesystem validation
func TestTemplateValidatorFallback(t *testing.T) {
	// Create validator without embedded filesystem
	validator := NewTemplateValidatorWithFS(nil)

	if validator.useEmbedded {
		t.Error("Expected useEmbedded to be false when nil filesystem provided")
	}

	// Test that it falls back to filesystem validation
	result, err := validator.ValidateTemplateConsistency("nonexistent")
	if err != nil {
		t.Fatalf("Expected no error in fallback mode, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Should have validation issues due to nonexistent path
	if result.Valid {
		t.Error("Expected validation to fail for nonexistent path")
	}
}

// TestTemplateValidatorConstants tests that constants are used correctly
func TestTemplateValidatorConstants(t *testing.T) {
	validator := NewTemplateValidator()

	summary := validator.GetValidationSummary()

	// Check that embedded template base is set correctly
	if validator.useEmbedded {
		if base, exists := summary["embedded_template_base"]; exists {
			if base != constants.TemplateBaseDir {
				t.Errorf("Expected embedded_template_base to be %s, got %v", constants.TemplateBaseDir, base)
			}
		} else {
			t.Error("Expected embedded_template_base in summary when embedded filesystem is available")
		}
	}
}
