package standards

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTemplateValidator(t *testing.T) {
	validator := NewTemplateValidator()
	if validator == nil {
		t.Fatal("NewTemplateValidator() returned nil")
	}

	if validator.standards == nil {
		t.Error("Validator should have standards initialized")
	}
}

func TestValidateTemplate(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validator := NewTemplateValidator()

	// Test with missing template directory
	result, err := validator.ValidateTemplate("/nonexistent", "nextjs-app")
	if err != nil {
		t.Fatalf("ValidateTemplate should not return error for missing directory: %v", err)
	}

	if result.IsValid {
		t.Error("Validation should fail for missing template files")
	}

	if len(result.Errors) == 0 {
		t.Error("Should have validation errors for missing files")
	}
}

func TestValidatePackageJSON(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "package-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validator := NewTemplateValidator()

	testCases := []struct {
		name         string
		packageJSON  map[string]interface{}
		templateType string
		expectErrors bool
	}{
		{
			name: "valid package.json",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "0.1.0",
				"private": true,
				"scripts": map[string]interface{}{
					"dev":           "next dev",
					"build":         "next build",
					"start":         "next start",
					"lint":          "next lint",
					"lint:fix":      "next lint --fix",
					"type-check":    "tsc --noEmit",
					"test":          "jest",
					"test:watch":    "jest --watch",
					"test:coverage": "jest --coverage",
					"format":        "prettier --write .",
					"format:check":  "prettier --check .",
					"clean":         "rm -rf .next out dist",
				},
				"dependencies": map[string]interface{}{
					"next":                     "{{.Versions.NextJS}}",
					"react":                    "{{.Versions.React}}",
					"react-dom":                "{{.Versions.React}}",
					"typescript":               "^5.3.0",
					"@types/node":              "^20.10.0",
					"@types/react":             "^18.2.0",
					"@types/react-dom":         "^18.2.0",
					"tailwindcss":              "^3.4.0",
					"autoprefixer":             "^10.4.0",
					"postcss":                  "^8.4.0",
					"clsx":                     "^2.0.0",
					"class-variance-authority": "^0.7.0",
					"tailwind-merge":           "^2.2.0",
					"lucide-react":             "^0.300.0",
					"@radix-ui/react-slot":     "^1.0.0",
					"tailwindcss-animate":      "^1.0.7",
				},
				"devDependencies": map[string]interface{}{
					"eslint":                           "^8.55.0",
					"eslint-config-next":               "{{.Versions.NextJS}}",
					"@typescript-eslint/eslint-plugin": "^6.15.0",
					"@typescript-eslint/parser":        "^6.15.0",
					"prettier":                         "^3.1.0",
					"prettier-plugin-tailwindcss":      "^0.5.0",
					"jest":                             "^29.7.0",
					"jest-environment-jsdom":           "^29.7.0",
					"@testing-library/react":           "^14.1.0",
					"@testing-library/jest-dom":        "^6.1.0",
					"@testing-library/user-event":      "^14.5.0",
					"@types/jest":                      "^29.5.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
			templateType: "nextjs-app",
			expectErrors: false,
		},
		{
			name: "missing scripts",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "0.1.0",
				"private": true,
				"scripts": map[string]interface{}{
					"dev": "next dev",
				},
				"dependencies": map[string]interface{}{
					"next":  "{{.Versions.NextJS}}",
					"react": "{{.Versions.React}}",
				},
				"devDependencies": map[string]interface{}{
					"eslint": "^8.55.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
			templateType: "nextjs-app",
			expectErrors: true,
		},
		{
			name: "missing dependencies",
			packageJSON: map[string]interface{}{
				"name":    "test-app",
				"version": "0.1.0",
				"private": true,
				"scripts": map[string]interface{}{
					"dev":           "next dev",
					"build":         "next build",
					"start":         "next start",
					"lint":          "next lint",
					"lint:fix":      "next lint --fix",
					"type-check":    "tsc --noEmit",
					"test":          "jest",
					"test:watch":    "jest --watch",
					"test:coverage": "jest --coverage",
					"format":        "prettier --write .",
					"format:check":  "prettier --check .",
					"clean":         "rm -rf .next out dist",
				},
				"dependencies": map[string]interface{}{
					"next": "{{.Versions.NextJS}}",
				},
				"devDependencies": map[string]interface{}{
					"eslint": "^8.55.0",
				},
				"engines": map[string]interface{}{
					"node": ">=22.0.0",
					"npm":  ">=10.0.0",
				},
			},
			templateType: "nextjs-app",
			expectErrors: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create package.json file
			packagePath := filepath.Join(tempDir, "package.json.tmpl")
			packageBytes, err := json.MarshalIndent(tc.packageJSON, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal package.json: %v", err)
			}

			if err := os.WriteFile(packagePath, packageBytes, 0644); err != nil {
				t.Fatalf("Failed to write package.json: %v", err)
			}

			result := &ValidationResult{
				IsValid:      true,
				Errors:       []ValidationError{},
				Warnings:     []ValidationWarning{},
				Suggestions:  []ValidationSuggestion{},
				TemplateName: tc.templateType,
			}

			err = validator.validatePackageJSON(tempDir, tc.templateType, result)
			if err != nil {
				t.Fatalf("validatePackageJSON failed: %v", err)
			}

			if tc.expectErrors && len(result.Errors) == 0 {
				t.Error("Expected validation errors but got none")
			}

			if !tc.expectErrors && len(result.Errors) > 0 {
				t.Errorf("Expected no validation errors but got %d: %v", len(result.Errors), result.Errors)
			}

			// Clean up
			os.Remove(packagePath)
		})
	}
}

func TestValidateTSConfig(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "tsconfig-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validator := NewTemplateValidator()

	testCases := []struct {
		name         string
		tsconfig     map[string]interface{}
		expectErrors bool
	}{
		{
			name: "valid tsconfig.json",
			tsconfig: map[string]interface{}{
				"compilerOptions": map[string]interface{}{
					"target":            "es5",
					"lib":               []string{"dom", "dom.iterable", "es6"},
					"allowJs":           true,
					"skipLibCheck":      true,
					"strict":            true,
					"noEmit":            true,
					"esModuleInterop":   true,
					"module":            "esnext",
					"moduleResolution":  "bundler",
					"resolveJsonModule": true,
					"isolatedModules":   true,
					"jsx":               "preserve",
					"incremental":       true,
					"baseUrl":           ".",
					"paths": map[string]interface{}{
						"@/*": []string{"./src/*"},
					},
				},
				"include": []string{"next-env.d.ts", "**/*.ts", "**/*.tsx"},
				"exclude": []string{"node_modules"},
			},
			expectErrors: false,
		},
		{
			name: "missing compilerOptions",
			tsconfig: map[string]interface{}{
				"include": []string{"**/*.ts"},
				"exclude": []string{"node_modules"},
			},
			expectErrors: true,
		},
		{
			name: "invalid JSON",
			tsconfig: map[string]interface{}{
				"compilerOptions": "invalid",
			},
			expectErrors: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create tsconfig.json file
			tsconfigPath := filepath.Join(tempDir, "tsconfig.json.tmpl")
			tsconfigBytes, err := json.MarshalIndent(tc.tsconfig, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal tsconfig.json: %v", err)
			}

			if err := os.WriteFile(tsconfigPath, tsconfigBytes, 0644); err != nil {
				t.Fatalf("Failed to write tsconfig.json: %v", err)
			}

			result := &ValidationResult{
				IsValid:      true,
				Errors:       []ValidationError{},
				Warnings:     []ValidationWarning{},
				Suggestions:  []ValidationSuggestion{},
				TemplateName: "nextjs-app",
			}

			err = validator.validateTSConfig(tempDir, result)
			if err != nil {
				t.Fatalf("validateTSConfig failed: %v", err)
			}

			if tc.expectErrors && len(result.Errors) == 0 {
				t.Error("Expected validation errors but got none")
			}

			if !tc.expectErrors && len(result.Errors) > 0 {
				t.Errorf("Expected no validation errors but got %d: %v", len(result.Errors), result.Errors)
			}

			// Clean up
			os.Remove(tsconfigPath)
		})
	}
}

func TestValidateESLintConfig(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "eslint-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validator := NewTemplateValidator()

	testCases := []struct {
		name         string
		eslintConfig map[string]interface{}
		expectErrors bool
	}{
		{
			name: "valid .eslintrc.json",
			eslintConfig: map[string]interface{}{
				"extends": []string{
					"next/core-web-vitals",
					"@typescript-eslint/recommended",
				},
				"parser":  "@typescript-eslint/parser",
				"plugins": []string{"@typescript-eslint"},
				"rules": map[string]interface{}{
					"@typescript-eslint/no-unused-vars": "error",
				},
				"ignorePatterns": []string{"node_modules/", ".next/"},
			},
			expectErrors: false,
		},
		{
			name: "missing extends",
			eslintConfig: map[string]interface{}{
				"parser":  "@typescript-eslint/parser",
				"plugins": []string{"@typescript-eslint"},
				"rules": map[string]interface{}{
					"@typescript-eslint/no-unused-vars": "error",
				},
			},
			expectErrors: false, // This should only generate suggestions, not errors
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create .eslintrc.json file
			eslintPath := filepath.Join(tempDir, ".eslintrc.json.tmpl")
			eslintBytes, err := json.MarshalIndent(tc.eslintConfig, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal .eslintrc.json: %v", err)
			}

			if err := os.WriteFile(eslintPath, eslintBytes, 0644); err != nil {
				t.Fatalf("Failed to write .eslintrc.json: %v", err)
			}

			result := &ValidationResult{
				IsValid:      true,
				Errors:       []ValidationError{},
				Warnings:     []ValidationWarning{},
				Suggestions:  []ValidationSuggestion{},
				TemplateName: "nextjs-app",
			}

			err = validator.validateESLintConfig(tempDir, result)
			if err != nil {
				t.Fatalf("validateESLintConfig failed: %v", err)
			}

			if tc.expectErrors && len(result.Errors) == 0 {
				t.Error("Expected validation errors but got none")
			}

			if !tc.expectErrors && len(result.Errors) > 0 {
				t.Errorf("Expected no validation errors but got %d: %v", len(result.Errors), result.Errors)
			}

			// Clean up
			os.Remove(eslintPath)
		})
	}
}

func TestCompareValues(t *testing.T) {
	testCases := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"equal bools", true, true, true},
		{"different bools", true, false, false},
		{"equal numbers", 42.0, 42.0, true},
		{"different numbers", 42.0, 43.0, false},
		{"equal arrays", []interface{}{"a", "b"}, []interface{}{"a", "b"}, true},
		{"different arrays", []interface{}{"a", "b"}, []interface{}{"a", "c"}, false},
		{"different array lengths", []interface{}{"a"}, []interface{}{"a", "b"}, false},
		{"equal objects", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}, true},
		{"different objects", map[string]interface{}{"key": "value1"}, map[string]interface{}{"key": "value2"}, false},
		{"different types", "hello", 42, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareValues(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("compareValues(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test existing file
	if !fileExists(tempFile.Name()) {
		t.Error("fileExists should return true for existing file")
	}

	// Test non-existing file
	if fileExists("/nonexistent/file") {
		t.Error("fileExists should return false for non-existing file")
	}
}
