package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewSecurityValidator(t *testing.T) {
	validator := NewSecurityValidator()
	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}
	if validator.cache == nil {
		t.Fatal("Expected cache to be initialized")
	}
	if validator.client == nil {
		t.Fatal("Expected HTTP client to be initialized")
	}
}

func TestValidatePackageVersionSecurity(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name        string
		packageName string
		version     string
		ecosystem   string
		expectVulns int
		expectError bool
	}{
		{
			name:        "vulnerable lodash package",
			packageName: "lodash",
			version:     "4.17.20",
			ecosystem:   "npm",
			expectVulns: 1,
			expectError: false,
		},
		{
			name:        "safe lodash package",
			packageName: "lodash",
			version:     "4.17.21",
			ecosystem:   "npm",
			expectVulns: 0,
			expectError: false,
		},
		{
			name:        "vulnerable axios package",
			packageName: "axios",
			version:     "0.21.1",
			ecosystem:   "npm",
			expectVulns: 1,
			expectError: false,
		},
		{
			name:        "safe axios package",
			packageName: "axios",
			version:     "0.21.2",
			ecosystem:   "npm",
			expectVulns: 0,
			expectError: false,
		},
		{
			name:        "vulnerable go module",
			packageName: "github.com/gin-gonic/gin",
			version:     "v1.6.0",
			ecosystem:   "go",
			expectVulns: 1,
			expectError: false,
		},
		{
			name:        "safe go module",
			packageName: "github.com/gin-gonic/gin",
			version:     "v1.7.0",
			ecosystem:   "go",
			expectVulns: 0,
			expectError: false,
		},
		{
			name:        "unknown package",
			packageName: "unknown-package",
			version:     "1.0.0",
			ecosystem:   "npm",
			expectVulns: 0,
			expectError: false,
		},
		{
			name:        "unsupported ecosystem",
			packageName: "some-package",
			version:     "1.0.0",
			ecosystem:   "python",
			expectVulns: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidatePackageVersionSecurity(tt.packageName, tt.version, tt.ecosystem)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(result.Vulnerabilities) != tt.expectVulns {
				t.Errorf("Expected %d vulnerabilities, got %d", tt.expectVulns, len(result.Vulnerabilities))
			}

			if result.Package != tt.packageName {
				t.Errorf("Expected package name %s, got %s", tt.packageName, result.Package)
			}

			if result.Version != tt.version {
				t.Errorf("Expected version %s, got %s", tt.version, result.Version)
			}
		})
	}
}

func TestValidateSecurityVulnerabilities(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "project with vulnerable npm packages",
			setupProject: func(projectPath string) error {
				packageJSON := map[string]interface{}{
					"name":    "test-app",
					"version": "1.0.0",
					"dependencies": map[string]interface{}{
						"lodash": "4.17.20", // vulnerable version
						"axios":  "0.21.1",  // vulnerable version
					},
					"devDependencies": map[string]interface{}{
						"jest": "^29.0.0", // safe package
					},
				}

				data, _ := json.Marshal(packageJSON)
				return os.WriteFile(filepath.Join(projectPath, "package.json"), data, 0644)
			},
			expectValid:    false, // high severity vulnerabilities should fail validation
			expectErrors:   1,     // lodash has high severity
			expectWarnings: 1,     // axios has moderate severity
		},
		{
			name: "project with safe packages",
			setupProject: func(projectPath string) error {
				packageJSON := map[string]interface{}{
					"name":    "test-app",
					"version": "1.0.0",
					"dependencies": map[string]interface{}{
						"lodash": "4.17.21", // safe version
						"axios":  "0.21.2",  // safe version
					},
				}

				data, _ := json.Marshal(packageJSON)
				return os.WriteFile(filepath.Join(projectPath, "package.json"), data, 0644)
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "project with vulnerable go modules",
			setupProject: func(projectPath string) error {
				goMod := `module test-app

go 1.21

require github.com/gin-gonic/gin v1.6.0
`
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectValid:    true, // go module vulnerability is moderate severity
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "project with no dependencies",
			setupProject: func(projectPath string) error {
				packageJSON := map[string]interface{}{
					"name":    "test-app",
					"version": "1.0.0",
				}

				data, _ := json.Marshal(packageJSON)
				return os.WriteFile(filepath.Join(projectPath, "package.json"), data, 0644)
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setupProject(tmpDir); err != nil {
				t.Fatalf("Failed to setup test project: %v", err)
			}

			result, err := validator.ValidateSecurityVulnerabilities(tmpDir)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

func TestPrioritizeSecurityUpdates(t *testing.T) {
	validator := NewSecurityValidator()

	vulnerabilities := []VulnerabilityInfo{
		{
			ID:       "vuln-1",
			Severity: "low",
			Title:    "Low severity vulnerability",
		},
		{
			ID:       "vuln-2",
			Severity: "critical",
			Title:    "Critical vulnerability",
		},
		{
			ID:       "vuln-3",
			Severity: "moderate",
			Title:    "Moderate vulnerability",
		},
		{
			ID:       "vuln-4",
			Severity: "high",
			Title:    "High severity vulnerability",
		},
		{
			ID:       "vuln-5",
			Severity: "info",
			Title:    "Info level vulnerability",
		},
	}

	prioritized := validator.PrioritizeSecurityUpdates(vulnerabilities)

	expectedOrder := []string{"critical", "high", "moderate", "low", "info"}

	if len(prioritized) != len(vulnerabilities) {
		t.Fatalf("Expected %d vulnerabilities, got %d", len(vulnerabilities), len(prioritized))
	}

	for i, expected := range expectedOrder {
		if prioritized[i].Severity != expected {
			t.Errorf("Expected severity %s at position %d, got %s", expected, i, prioritized[i].Severity)
		}
	}
}

func TestCleanVersion(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		input    string
		expected string
	}{
		{"^1.0.0", "1.0.0"},
		{"~1.0.0", "1.0.0"},
		{">=1.0.0", "1.0.0"},
		{"<=1.0.0", "1.0.0"},
		{">1.0.0", "1.0.0"},
		{"<1.0.0", "1.0.0"},
		{"=1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := validator.cleanVersion(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsVersionAffected(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		currentVersion string
		fixedVersion   string
		expected       bool
	}{
		{"1.0.0", "1.0.1", true},   // current < fixed = affected
		{"1.0.1", "1.0.0", false},  // current > fixed = not affected
		{"1.0.0", "1.0.0", false},  // current = fixed = not affected
		{"v1.0.0", "v1.0.1", true}, // with v prefix
		{"1.0.0", "", true},        // no fix available = affected
	}

	for _, tt := range tests {
		t.Run(tt.currentVersion+"_vs_"+tt.fixedVersion, func(t *testing.T) {
			result := validator.isVersionAffected(tt.currentVersion, tt.fixedVersion)
			if result != tt.expected {
				t.Errorf("Expected %v for %s vs %s, got %v", tt.expected, tt.currentVersion, tt.fixedVersion, result)
			}
		})
	}
}

func TestFindPackageJSONFiles(t *testing.T) {
	validator := NewSecurityValidator()

	tmpDir := t.TempDir()

	// Create test structure
	testFiles := []string{
		"package.json",
		"frontend/package.json",
		"backend/package.json",
		"node_modules/some-package/package.json", // should be ignored
	}

	for _, file := range testFiles {
		dir := filepath.Dir(filepath.Join(tmpDir, file))
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		if err := os.WriteFile(filepath.Join(tmpDir, file), []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	files, err := validator.findPackageJSONFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should find 3 files (excluding node_modules)
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d: %v", expectedCount, len(files), files)
	}

	// Check that node_modules file is not included
	for _, file := range files {
		if strings.Contains(file, "node_modules") {
			t.Errorf("Found file in node_modules, should be excluded: %s", file)
		}
	}
}

func TestFindGoModFiles(t *testing.T) {
	validator := NewSecurityValidator()

	tmpDir := t.TempDir()

	// Create test structure
	testFiles := []string{
		"go.mod",
		"backend/go.mod",
		"services/auth/go.mod",
	}

	for _, file := range testFiles {
		dir := filepath.Dir(filepath.Join(tmpDir, file))
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		if err := os.WriteFile(filepath.Join(tmpDir, file), []byte("module test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	files, err := validator.findGoModFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d: %v", expectedCount, len(files), files)
	}
}

func TestSecurityResultCaching(t *testing.T) {
	validator := NewSecurityValidator()

	// First call should populate cache
	result1, err := validator.ValidatePackageVersionSecurity("lodash", "4.17.20", "npm")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Second call should use cache
	result2, err := validator.ValidatePackageVersionSecurity("lodash", "4.17.20", "npm")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Results should be identical
	if result1.Package != result2.Package {
		t.Error("Cached result package name differs")
	}

	if len(result1.Vulnerabilities) != len(result2.Vulnerabilities) {
		t.Error("Cached result vulnerabilities count differs")
	}

	// Check that cache was actually used (timestamps should be very close)
	timeDiff := result2.LastChecked.Sub(result1.LastChecked)
	if timeDiff > time.Millisecond {
		t.Error("Cache doesn't appear to be working - timestamps differ significantly")
	}
}
