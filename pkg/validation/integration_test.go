package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/interfaces"
)

func TestValidationEngineIntegration(t *testing.T) {
	// Create a validation engine
	engine := NewEngine()

	// Verify it implements the interface
	var _ interfaces.ValidationEngine = engine

	// Test template consistency validation
	t.Run("template consistency validation", func(t *testing.T) {
		tmpDir := t.TempDir()
		frontendDir := filepath.Join(tmpDir, "frontend")
		if err := os.MkdirAll(frontendDir, 0755); err != nil {
			t.Fatalf("Failed to create frontend directory: %v", err)
		}

		// Create a simple template
		templateDir := filepath.Join(frontendDir, "nextjs-app")
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		packageJSON := map[string]interface{}{
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
		}

		data, _ := json.Marshal(packageJSON)
		packageJSONPath := filepath.Join(templateDir, "package.json.tmpl")
		if err := os.WriteFile(packageJSONPath, data, 0644); err != nil {
			t.Fatalf("Failed to write package.json: %v", err)
		}

		result, err := engine.ValidateTemplateConsistency(tmpDir)
		if err != nil {
			t.Fatalf("Template consistency validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected template consistency validation to pass, got errors: %v", result.Errors)
		}
	})

	// Test Vercel compatibility validation
	t.Run("vercel compatibility validation", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a Next.js project structure
		packageJSON := map[string]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
			"scripts": map[string]interface{}{
				"build": "next build",
				"start": "next start",
				"dev":   "next dev",
			},
			"dependencies": map[string]interface{}{
				"next":  "15.5.2",
				"react": "19.1.0",
			},
			"engines": map[string]interface{}{
				"node": ">=22.0.0",
			},
		}

		data, _ := json.Marshal(packageJSON)
		if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), data, 0644); err != nil {
			t.Fatalf("Failed to write package.json: %v", err)
		}

		// Create src/app directory
		if err := os.MkdirAll(filepath.Join(tmpDir, "src", "app"), 0755); err != nil {
			t.Fatalf("Failed to create src/app directory: %v", err)
		}

		result, err := engine.ValidateVercelCompatibility(tmpDir)
		if err != nil {
			t.Fatalf("Vercel compatibility validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected Vercel compatibility validation to pass, got errors: %v", result.Errors)
		}
	})

	// Test security vulnerability validation
	t.Run("security vulnerability validation", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a project with safe packages
		packageJSON := map[string]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
			"dependencies": map[string]interface{}{
				"lodash": "4.17.21", // safe version
			},
		}

		data, _ := json.Marshal(packageJSON)
		if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), data, 0644); err != nil {
			t.Fatalf("Failed to write package.json: %v", err)
		}

		result, err := engine.ValidateSecurityVulnerabilities(tmpDir)
		if err != nil {
			t.Fatalf("Security vulnerability validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected security validation to pass for safe packages, got errors: %v", result.Errors)
		}
	})

	// Test package.json structure validation
	t.Run("package.json structure validation", func(t *testing.T) {
		tmpDir := t.TempDir()

		packageJSON := map[string]interface{}{
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
			"dependencies": map[string]interface{}{
				"next": "15.5.2",
			},
			"devDependencies": map[string]interface{}{
				"typescript": "^5.3.0",
			},
			"engines": map[string]interface{}{
				"node": ">=22.0.0",
				"npm":  ">=10.0.0",
			},
		}

		data, _ := json.Marshal(packageJSON)
		packageJSONPath := filepath.Join(tmpDir, "package.json")
		if err := os.WriteFile(packageJSONPath, data, 0644); err != nil {
			t.Fatalf("Failed to write package.json: %v", err)
		}

		result, err := engine.ValidatePackageJSONStructure(packageJSONPath)
		if err != nil {
			t.Fatalf("Package.json structure validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected package.json structure validation to pass, got errors: %v", result.Errors)
		}
	})

	// Test TypeScript config validation
	t.Run("typescript config validation", func(t *testing.T) {
		tmpDir := t.TempDir()

		tsconfig := map[string]interface{}{
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
		}

		data, _ := json.Marshal(tsconfig)
		tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
		if err := os.WriteFile(tsconfigPath, data, 0644); err != nil {
			t.Fatalf("Failed to write tsconfig.json: %v", err)
		}

		result, err := engine.ValidateTypeScriptConfig(tsconfigPath)
		if err != nil {
			t.Fatalf("TypeScript config validation failed: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected TypeScript config validation to pass, got errors: %v", result.Errors)
		}
	})
}

func TestValidationEngineMethodsExist(t *testing.T) {
	engine := NewEngine()

	// Test that all new methods are accessible
	methods := []struct {
		name string
		test func() error
	}{
		{
			name: "ValidateTemplateConsistency",
			test: func() error {
				_, err := engine.ValidateTemplateConsistency("/tmp")
				return err
			},
		},
		{
			name: "ValidatePackageJSONStructure",
			test: func() error {
				tmpDir := t.TempDir()
				packageJSONPath := filepath.Join(tmpDir, "package.json")
				if err := os.WriteFile(packageJSONPath, []byte("{}"), 0644); err != nil {
					return err
				}
				_, err := engine.ValidatePackageJSONStructure(packageJSONPath)
				return err
			},
		},
		{
			name: "ValidateTypeScriptConfig",
			test: func() error {
				tmpDir := t.TempDir()
				tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
				if err := os.WriteFile(tsconfigPath, []byte("{}"), 0644); err != nil {
					return err
				}
				_, err := engine.ValidateTypeScriptConfig(tsconfigPath)
				return err
			},
		},
		{
			name: "ValidateVercelCompatibility",
			test: func() error {
				tmpDir := t.TempDir()
				_, err := engine.ValidateVercelCompatibility(tmpDir)
				return err
			},
		},
		{
			name: "ValidateSecurityVulnerabilities",
			test: func() error {
				tmpDir := t.TempDir()
				_, err := engine.ValidateSecurityVulnerabilities(tmpDir)
				return err
			},
		},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			// We don't care about the specific result, just that the method exists and can be called
			_ = method.test()
		})
	}
}
