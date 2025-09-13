package cleanup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTODOResolver_determineTODOAction(t *testing.T) {
	resolver := NewTODOResolver("/test", &TODOResolverConfig{})

	tests := []struct {
		name     string
		todo     TODOItem
		expected TODOAction
	}{
		{
			name: "false positive - documentation file",
			todo: TODOItem{
				File:    "/test/docs/README.md",
				Message: "implement feature",
			},
			expected: TODOActionIgnore,
		},
		{
			name: "false positive - spec file",
			todo: TODOItem{
				File:    "/test/.kiro/specs/project-cleanup/tasks.md",
				Message: "TODO comments",
			},
			expected: TODOActionIgnore,
		},
		{
			name: "legitimate code reference",
			todo: TODOItem{
				File:    "/test/pkg/template/import_detector.go",
				Context: `"context.TODO": "context",`,
			},
			expected: TODOActionIgnore,
		},
		{
			name: "template file TODO",
			todo: TODOItem{
				File:    "/test/templates/backend/go-gin/service.go.tmpl",
				Message: "implement email sending",
			},
			expected: TODOActionDocument,
		},
		{
			name: "obsolete security TODO",
			todo: TODOItem{
				File:     "/test/pkg/version/npm_registry.go",
				Message:  "implement security checking",
				Category: CategorySecurity,
			},
			expected: TODOActionRemove,
		},
		{
			name: "resolvable email TODO",
			todo: TODOItem{
				File:    "/test/internal/services/auth.go",
				Message: "send email with reset token",
			},
			expected: TODOActionResolve,
		},
		{
			name: "feature TODO for documentation",
			todo: TODOItem{
				File:     "/test/pkg/template/engine.go",
				Message:  "add caching support",
				Category: CategoryFeature,
			},
			expected: TODOActionDocument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := resolver.determineTODOAction(tt.todo)
			if action != tt.expected {
				t.Errorf("determineTODOAction() = %v, want %v", action, tt.expected)
			}
		})
	}
}

func TestTODOResolver_isFalsePositive(t *testing.T) {
	resolver := NewTODOResolver("/test", &TODOResolverConfig{})

	tests := []struct {
		name     string
		todo     TODOItem
		expected bool
	}{
		{
			name:     "documentation file",
			todo:     TODOItem{File: "/test/docs/README.md"},
			expected: true,
		},
		{
			name:     "markdown file",
			todo:     TODOItem{File: "/test/CONTRIBUTING.md"},
			expected: true,
		},
		{
			name:     "spec file",
			todo:     TODOItem{File: "/test/.kiro/specs/feature/design.md"},
			expected: true,
		},
		{
			name:     "script file",
			todo:     TODOItem{File: "/test/scripts/audit.sh"},
			expected: true,
		},
		{
			name:     "TODO comment about TODOs",
			todo:     TODOItem{Message: "Check for TODO comments"},
			expected: true,
		},
		{
			name:     "legitimate Go file",
			todo:     TODOItem{File: "/test/pkg/template/engine.go"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.isFalsePositive(tt.todo)
			if result != tt.expected {
				t.Errorf("isFalsePositive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTODOResolver_isLegitimateCodeReference(t *testing.T) {
	resolver := NewTODOResolver("/test", &TODOResolverConfig{})

	tests := []struct {
		name     string
		todo     TODOItem
		expected bool
	}{
		{
			name:     "context.TODO reference",
			todo:     TODOItem{Context: `"context.TODO": "context",`},
			expected: true,
		},
		{
			name:     "PR template TODO check",
			todo:     TODOItem{Context: "- [ ] No TODO comments without issues"},
			expected: true,
		},
		{
			name:     "actual TODO comment",
			todo:     TODOItem{Context: "// TODO: implement this feature"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.isLegitimateCodeReference(tt.todo)
			if result != tt.expected {
				t.Errorf("isLegitimateCodeReference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTODOResolver_isTemplateFile(t *testing.T) {
	resolver := NewTODOResolver("/test", &TODOResolverConfig{})

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "template file",
			filePath: "/test/templates/backend/go-gin/service.go.tmpl",
			expected: true,
		},
		{
			name:     "regular go file",
			filePath: "/test/pkg/template/engine.go",
			expected: false,
		},
		{
			name:     "template directory but not template file",
			filePath: "/test/templates/README.md",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.isTemplateFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("isTemplateFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTODOResolver_resolveEmailTODO(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	content := `package main

func resetPassword() {
	// NOTE: Email sending should be implemented based on your email service provider
	// In a real application, you would send an email here
}
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	resolver := NewTODOResolver(tmpDir, &TODOResolverConfig{})

	todo := TODOItem{
		File:    testFile,
		Line:    4,
		Message: "Send email with reset token",
		// NOTE: Email sending should be implemented based on your email service provider
	}

	resolved, err := resolver.resolveEmailTODO(todo)
	if err != nil {
		t.Fatalf("resolveEmailTODO() failed: %v", err)
	}

	if resolved.Action != "Replaced with documentation" {
		t.Errorf("Expected action 'Replaced with documentation', got %s", resolved.Action)
	}

	// Check that the file was modified
	modifiedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}

	if !strings.Contains(string(modifiedContent), "NOTE: Email sending should be implemented") {
		t.Error("File was not properly modified with documentation")
	}
}

func TestTODOResolver_getDocumentationReason(t *testing.T) {
	resolver := NewTODOResolver("/test", &TODOResolverConfig{})

	tests := []struct {
		name     string
		todo     TODOItem
		expected string
	}{
		{
			name: "template file",
			todo: TODOItem{
				File: "/test/templates/backend/service.go.tmpl",
			},
			expected: "Template placeholder - intentional TODO for generated projects",
		},
		{
			name: "security TODO",
			todo: TODOItem{
				Category: CategorySecurity,
			},
			expected: "Security enhancement - requires careful implementation and testing",
		},
		{
			name: "performance TODO",
			todo: TODOItem{
				Category: CategoryPerformance,
			},
			expected: "Performance optimization - requires benchmarking and analysis",
		},
		{
			name: "feature TODO",
			todo: TODOItem{
				Category: CategoryFeature,
			},
			expected: "Feature enhancement - requires design and implementation planning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.getDocumentationReason(tt.todo)
			if result != tt.expected {
				t.Errorf("getDocumentationReason() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTODOResolutionReport_GenerateReport(t *testing.T) {
	report := &TODOResolutionReport{
		TotalFound: 10,
		Resolved: []ResolvedTODO{
			{
				Original: TODOItem{
					File:    "/test/service.go",
					Line:    42,
					Message: "send email",
				},
				Action:      "Replaced with documentation",
				Description: "Added proper implementation guidance",
			},
		},
		Documented: []DocumentedTODO{
			{
				Original: TODOItem{
					File:     "/test/templates/service.tmpl",
					Line:     10,
					Message:  "implement feature",
					Category: CategoryFeature,
				},
				Reason:     "Template placeholder",
				Documented: true,
				FutureWork: true,
			},
		},
		Removed: []RemovedTODO{
			{
				Original: TODOItem{
					File:    "/test/old.go",
					Line:    5,
					Message: "obsolete feature",
				},
				Reason: "Feature already implemented",
			},
		},
		FalsePositives: []FalsePositiveTODO{
			{
				Original: TODOItem{
					File:    "/test/docs/README.md",
					Line:    1,
					Message: "documentation",
				},
				Reason: "Documentation file",
			},
		},
	}

	result := report.GenerateReport()

	// Check that all sections are present
	expectedSections := []string{
		"# TODO Resolution Report",
		"**Total TODOs Found:** 10",
		"## Summary",
		"- **Resolved:** 1",
		"- **Documented for Future Work:** 1",
		"- **Removed (Obsolete):** 1",
		"- **False Positives:** 1",
		"## Resolved TODOs",
		"## Documented for Future Work",
		"## Removed TODOs",
		"## False Positives",
	}

	for _, section := range expectedSections {
		if !strings.Contains(result, section) {
			t.Errorf("Report missing expected section: %s", section)
		}
	}
}
