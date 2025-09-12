package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewTemplateValidator(t *testing.T) {
	validator := NewTemplateValidator()
	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}
	if validator.standardConfigs == nil {
		t.Fatal("Expected standardConfigs to be initialized")
	}
}

func TestValidatePackageJSONStructure(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		packageJSON    map[string]interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid package.json",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "1.0.0",
				"scripts": map[string]interface{}{
					"dev":        "next dev",
					"build":      "next build",
					"start":      "next start",
					"lint":       "next lint",
					"type-check": "tsc --noEmit",
					"test":       "jest",
					"format":     "prettier --write .",
					"clean":      "rm -rf .next",
				},
				"dependencies": map[string]interface{}{
					"next":  "15.5.2",
					"react": "19.1.0",
				},
				"devDependencies": map[string]interface{}{
					"typescript": "^5.3.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "missing required fields",
			packageJSON: map[string]interface{}{
				"name": "test-app",
			},
			expectValid:    false,
			expectErrors:   5, // missing version, scripts, dependencies, devDependencies, engines
			expectWarnings: 0,
		},
		{
			name: "missing required scripts",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "1.0.0",
				"scripts": map[string]interface{}{
					"dev": "next dev",
				},
				"dependencies":    map[string]interface{}{},
				"devDependencies": map[string]interface{}{},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
			expectValid:    false,
			expectErrors:   5, // missing build, start, lint, type-check, test
			expectWarnings: 0,
		},
		{
			name: "missing required engines",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "1.0.0",
				"scripts": map[string]interface{}{
					"dev":        "next dev",
					"build":      "next build",
					"start":      "next start",
					"lint":       "next lint",
					"type-check": "tsc --noEmit",
					"test":       "jest",
				},
				"dependencies":    map[string]interface{}{},
				"devDependencies": map[string]interface{}{},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					// missing npm
				},
			},
			expectValid:    false,
			expectErrors:   1, // missing npm engine
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			packageJSONPath := filepath.Join(tmpDir, "package.json")

			data, err := json.Marshal(tt.packageJSON)
			if err != nil {
				t.Fatalf("Failed to marshal test data: %v", err)
			}

			if err := os.WriteFile(packageJSONPath, data, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := validator.ValidatePackageJSONStructure(packageJSONPath)
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

func TestValidateTypeScriptConfig(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		tsconfig       map[string]interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid tsconfig.json",
			tsconfig: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"target":           "es5",
					"lib":              []string{"dom", "dom.iterable", "es6"},
					"strict":           true,
					"esModuleInterop":  true,
					"moduleResolution": "bundler",
					"jsx":              "preserve",
				},
				"include": []string{"**/*.ts", "**/*.tsx"},
				"exclude": []string{"node_modules"},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "missing required sections",
			tsconfig: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"target": "es5",
				},
			},
			expectValid:    false,
			expectErrors:   1, // missing include
			expectWarnings: 4, // missing lib, strict, esModuleInterop, moduleResolution
		},
		{
			name: "missing compiler options",
			tsconfig: map[string]interface{}{
				"include": []string{"**/*.ts"},
			},
			expectValid:    false,
			expectErrors:   1, // missing compilerOptions
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")

			data, err := json.Marshal(tt.tsconfig)
			if err != nil {
				t.Fatalf("Failed to marshal test data: %v", err)
			}

			if err := os.WriteFile(tsconfigPath, data, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := validator.ValidateTypeScriptConfig(tsconfigPath)
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

func TestValidateTemplateConsistency(t *testing.T) {
	validator := NewTemplateValidator()

	// Create temporary template structure
	tmpDir := t.TempDir()
	frontendDir := filepath.Join(tmpDir, "frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		t.Fatalf("Failed to create frontend directory: %v", err)
	}

	// Create test templates
	templates := []struct {
		name        string
		packageJSON map[string]interface{}
	}{
		{
			name: "nextjs-app",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "1.0.0",
				"scripts": map[string]interface{}{
					"dev":        "next dev",
					"build":      "next build",
					"start":      "next start",
					"lint":       "next lint",
					"type-check": "tsc --noEmit",
					"test":       "jest",
					"format":     "prettier --write .",
					"clean":      "rm -rf .next",
				},
				"dependencies": map[string]interface{}{
					"next":  "15.5.2",
					"react": "19.1.0",
				},
				"devDependencies": map[string]interface{}{
					"typescript": "^5.3.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
		},
		{
			name: "nextjs-home",
			packageJSON: map[string]interface{}{
				"name":    "test-home",
				"version": "1.0.0",
				"scripts": map[string]interface{}{
					"dev":        "next dev -p 3001",
					"build":      "next build",
					"start":      "next start -p 3001",
					"lint":       "next lint",
					"type-check": "tsc --noEmit",
					"test":       "jest",
					"format":     "prettier --write .",
					"clean":      "rm -rf .next",
				},
				"dependencies": map[string]interface{}{
					"next":  "15.5.2",
					"react": "19.1.0",
				},
				"devDependencies": map[string]interface{}{
					"typescript": "^5.3.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
		},
	}

	// Create template directories and files
	for _, tmpl := range templates {
		templateDir := filepath.Join(frontendDir, tmpl.name)
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		packageJSONPath := filepath.Join(templateDir, "package.json.tmpl")
		data, err := json.Marshal(tmpl.packageJSON)
		if err != nil {
			t.Fatalf("Failed to marshal package.json: %v", err)
		}

		if err := os.WriteFile(packageJSONPath, data, 0644); err != nil {
			t.Fatalf("Failed to write package.json: %v", err)
		}
	}

	result, err := validator.ValidateTemplateConsistency(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected validation to pass, but got errors: %v", result.Errors)
	}
}

func TestValidateTemplateConsistency_MissingFrontendDir(t *testing.T) {
	validator := NewTemplateValidator()
	tmpDir := t.TempDir()

	result, err := validator.ValidateTemplateConsistency(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Valid {
		t.Error("Expected validation to fail when frontend directory is missing")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected at least one error when frontend directory is missing")
	}
}

func TestValidateScripts(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		scripts        map[string]interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "all required scripts present",
			scripts: map[string]interface{}{
				"dev":        "next dev",
				"build":      "next build",
				"start":      "next start",
				"lint":       "next lint",
				"type-check": "tsc --noEmit",
				"test":       "jest",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "missing required scripts",
			scripts: map[string]interface{}{
				"dev":   "next dev",
				"build": "next build",
			},
			expectValid:    false,
			expectErrors:   4, // missing start, lint, type-check, test
			expectWarnings: 0,
		},
		{
			name: "empty script command",
			scripts: map[string]interface{}{
				"dev":        "next dev",
				"build":      "next build",
				"start":      "next start",
				"lint":       "next lint",
				"type-check": "tsc --noEmit",
				"test":       "",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // empty test command
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			validator.validateScripts(tt.scripts, result)

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

func TestValidateEngines(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		engines        map[string]interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid engines",
			engines: map[string]interface{}{
				"node": ">=22.0.0",
				"npm":  ">=10.0.0",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "missing required engines",
			engines: map[string]interface{}{
				"node": ">=22.0.0",
			},
			expectValid:    false,
			expectErrors:   1, // missing npm
			expectWarnings: 0,
		},
		{
			name:           "no engines",
			engines:        map[string]interface{}{},
			expectValid:    false,
			expectErrors:   2, // missing node and npm
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			validator.validateEngines(tt.engines, result)

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

func TestIsValidVersionConstraint(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		version string
		valid   bool
	}{
		{"^1.0.0", true},
		{"~1.0.0", true},
		{">=1.0.0", true},
		{"<=1.0.0", true},
		{">1.0.0", true},
		{"<1.0.0", true},
		{"=1.0.0", true},
		{"1.0.0", true},
		{"1.0", true},
		{"", false},
		{"   ", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := validator.isValidVersionConstraint(tt.version)
			if result != tt.valid {
				t.Errorf("Expected %v for version %q, got %v", tt.valid, tt.version, result)
			}
		})
	}
}

func TestContains(t *testing.T) {
	validator := NewTemplateValidator()

	slice := []string{"apple", "banana", "cherry"}

	tests := []struct {
		item     string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := validator.contains(slice, tt.item)
			if result != tt.expected {
				t.Errorf("Expected %v for item %q, got %v", tt.expected, tt.item, result)
			}
		})
	}
}
