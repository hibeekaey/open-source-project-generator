package standards

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewTemplateGenerator(t *testing.T) {
	generator := NewTemplateGenerator()
	if generator == nil {
		t.Fatal("NewTemplateGenerator() returned nil")
	}

	if generator.standards == nil {
		t.Error("Generator should have standards initialized")
	}
}

func TestGenerateStandardizedPackageJSON(t *testing.T) {
	generator := NewTemplateGenerator()

	testCases := []struct {
		templateType string
		expectedName string
	}{
		{"nextjs-app", "app"},
		{"nextjs-home", "home"},
		{"nextjs-admin", "admin"},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			result := generator.GenerateStandardizedPackageJSON(tc.templateType)

			if result == "" {
				t.Error("GenerateStandardizedPackageJSON should return non-empty string")
			}

			// Check that it contains expected template name
			if !strings.Contains(result, tc.expectedName) {
				t.Errorf("Generated package.json should contain template name '%s'", tc.expectedName)
			}

			// Check for required fields
			requiredFields := []string{
				"\"name\":",
				"\"version\":",
				"\"private\":",
				"\"scripts\":",
				"\"dependencies\":",
				"\"devDependencies\":",
				"\"engines\":",
			}

			for _, field := range requiredFields {
				if !strings.Contains(result, field) {
					t.Errorf("Generated package.json should contain field %s", field)
				}
			}

			// Check for required scripts
			requiredScripts := []string{
				"\"dev\":",
				"\"build\":",
				"\"start\":",
				"\"lint\":",
				"\"type-check\":",
				"\"test\":",
			}

			for _, script := range requiredScripts {
				if !strings.Contains(result, script) {
					t.Errorf("Generated package.json should contain script %s", script)
				}
			}

			// Check for required dependencies
			requiredDeps := []string{
				"\"next\":",
				"\"react\":",
				"\"react-dom\":",
				"\"typescript\":",
			}

			for _, dep := range requiredDeps {
				if !strings.Contains(result, dep) {
					t.Errorf("Generated package.json should contain dependency %s", dep)
				}
			}
		})
	}
}

func TestGenerateStandardizedTSConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedTSConfig()

	if result == "" {
		t.Error("GenerateStandardizedTSConfig should return non-empty string")
	}

	// Validate that it's valid JSON
	var tsconfig map[string]interface{}
	if err := json.Unmarshal([]byte(result), &tsconfig); err != nil {
		t.Fatalf("Generated tsconfig.json is not valid JSON: %v", err)
	}

	// Check for required fields
	requiredFields := []string{
		"compilerOptions",
		"include",
		"exclude",
	}

	for _, field := range requiredFields {
		if _, exists := tsconfig[field]; !exists {
			t.Errorf("Generated tsconfig.json should contain field '%s'", field)
		}
	}

	// Check compiler options
	compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{})
	if !ok {
		t.Fatal("compilerOptions should be an object")
	}

	requiredCompilerOptions := []string{
		"target",
		"lib",
		"strict",
		"jsx",
		"moduleResolution",
		"paths",
	}

	for _, option := range requiredCompilerOptions {
		if _, exists := compilerOptions[option]; !exists {
			t.Errorf("Generated tsconfig.json should contain compiler option '%s'", option)
		}
	}
}

func TestGenerateStandardizedESLintConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedESLintConfig()

	if result == "" {
		t.Error("GenerateStandardizedESLintConfig should return non-empty string")
	}

	// Validate that it's valid JSON
	var eslintConfig map[string]interface{}
	if err := json.Unmarshal([]byte(result), &eslintConfig); err != nil {
		t.Fatalf("Generated .eslintrc.json is not valid JSON: %v", err)
	}

	// Check for required fields
	requiredFields := []string{
		"extends",
		"parser",
		"plugins",
		"rules",
		"ignorePatterns",
	}

	for _, field := range requiredFields {
		if _, exists := eslintConfig[field]; !exists {
			t.Errorf("Generated .eslintrc.json should contain field '%s'", field)
		}
	}

	// Check extends
	extends, ok := eslintConfig["extends"].([]interface{})
	if !ok {
		t.Fatal("extends should be an array")
	}

	if len(extends) == 0 {
		t.Error("extends should not be empty")
	}

	// Check rules
	rules, ok := eslintConfig["rules"].(map[string]interface{})
	if !ok {
		t.Fatal("rules should be an object")
	}

	if len(rules) == 0 {
		t.Error("rules should not be empty")
	}
}

func TestGenerateStandardizedPrettierConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedPrettierConfig()

	if result == "" {
		t.Error("GenerateStandardizedPrettierConfig should return non-empty string")
	}

	// Validate that it's valid JSON
	var prettierConfig map[string]interface{}
	if err := json.Unmarshal([]byte(result), &prettierConfig); err != nil {
		t.Fatalf("Generated .prettierrc is not valid JSON: %v", err)
	}

	// Check for required fields
	requiredFields := []string{
		"semi",
		"trailingComma",
		"singleQuote",
		"printWidth",
		"tabWidth",
		"useTabs",
		"plugins",
	}

	for _, field := range requiredFields {
		if _, exists := prettierConfig[field]; !exists {
			t.Errorf("Generated .prettierrc should contain field '%s'", field)
		}
	}

	// Check plugins
	plugins, ok := prettierConfig["plugins"].([]interface{})
	if !ok {
		t.Fatal("plugins should be an array")
	}

	if len(plugins) == 0 {
		t.Error("plugins should not be empty")
	}
}

func TestGenerateStandardizedVercelConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedVercelConfig()

	if result == "" {
		t.Error("GenerateStandardizedVercelConfig should return non-empty string")
	}

	// Validate that it's valid JSON
	var vercelConfig map[string]interface{}
	if err := json.Unmarshal([]byte(result), &vercelConfig); err != nil {
		t.Fatalf("Generated vercel.json is not valid JSON: %v", err)
	}

	// Check for required fields
	requiredFields := []string{
		"buildCommand",
		"devCommand",
		"installCommand",
		"framework",
		"regions",
		"headers",
		"functions",
		"build",
		"env",
		"rewrites",
	}

	for _, field := range requiredFields {
		if _, exists := vercelConfig[field]; !exists {
			t.Errorf("Generated vercel.json should contain field '%s'", field)
		}
	}

	// Check framework
	framework, ok := vercelConfig["framework"].(string)
	if !ok || framework != "nextjs" {
		t.Error("framework should be 'nextjs'")
	}

	// Check headers
	headers, ok := vercelConfig["headers"].([]interface{})
	if !ok {
		t.Fatal("headers should be an array")
	}

	if len(headers) == 0 {
		t.Error("headers should not be empty")
	}
}

func TestGenerateStandardizedTailwindConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedTailwindConfig()

	if result == "" {
		t.Error("GenerateStandardizedTailwindConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"module.exports",
		"darkMode",
		"content",
		"theme",
		"extend",
		"colors",
		"plugins",
		"tailwindcss-animate",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Generated tailwind.config.js should contain '%s'", content)
		}
	}
}

func TestGenerateStandardizedNextConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedNextConfig()

	if result == "" {
		t.Error("GenerateStandardizedNextConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"nextConfig",
		"typescript",
		"eslint",
		"experimental",
		"images",
		"module.exports",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Generated next.config.js should contain '%s'", content)
		}
	}
}

func TestGenerateStandardizedPostCSSConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedPostCSSConfig()

	if result == "" {
		t.Error("GenerateStandardizedPostCSSConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"module.exports",
		"plugins",
		"tailwindcss",
		"autoprefixer",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Generated postcss.config.js should contain '%s'", content)
		}
	}
}

func TestGenerateStandardizedJestConfig(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedJestConfig()

	if result == "" {
		t.Error("GenerateStandardizedJestConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"nextJest",
		"createJestConfig",
		"setupFilesAfterEnv",
		"moduleNameMapping",
		"testEnvironment",
		"jest-environment-jsdom",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Generated jest.config.js should contain '%s'", content)
		}
	}
}

func TestGenerateStandardizedJestSetup(t *testing.T) {
	generator := NewTemplateGenerator()

	result := generator.GenerateStandardizedJestSetup()

	if result == "" {
		t.Error("GenerateStandardizedJestSetup should return non-empty string")
	}

	// Check for key content
	if !strings.Contains(result, "@testing-library/jest-dom") {
		t.Error("Generated jest.setup.js should import @testing-library/jest-dom")
	}
}

func TestGetScriptsForTemplate(t *testing.T) {
	generator := NewTemplateGenerator()

	testCases := []struct {
		templateType string
		expectPort   bool
		port         int
	}{
		{"nextjs-app", false, 0},
		{"nextjs-home", true, 3001},
		{"nextjs-admin", true, 3002},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			scripts := generator.getScriptsForTemplate(tc.templateType)

			if len(scripts) == 0 {
				t.Error("getScriptsForTemplate should return non-empty map")
			}

			// Check required scripts
			requiredScripts := []string{"dev", "build", "start", "lint", "type-check", "test"}
			for _, script := range requiredScripts {
				if _, exists := scripts[script]; !exists {
					t.Errorf("Required script '%s' is missing", script)
				}
			}

			// Check port-specific scripts
			if tc.expectPort {
				devScript := scripts["dev"]
				startScript := scripts["start"]

				// Simple check for port in script
				if !strings.Contains(devScript, "-p") {
					t.Errorf("Dev script should contain port flag for template type '%s'", tc.templateType)
				}

				if !strings.Contains(startScript, "-p") {
					t.Errorf("Start script should contain port flag for template type '%s'", tc.templateType)
				}
			}
		})
	}
}

func TestGetDependenciesForTemplate(t *testing.T) {
	generator := NewTemplateGenerator()

	testCases := []struct {
		templateType     string
		expectedBaseDeps []string
		expectedSpecific []string
	}{
		{
			templateType:     "nextjs-app",
			expectedBaseDeps: []string{"next", "react", "react-dom", "typescript"},
			expectedSpecific: []string{"@radix-ui/react-dialog", "@radix-ui/react-toast"},
		},
		{
			templateType:     "nextjs-home",
			expectedBaseDeps: []string{"next", "react", "react-dom", "typescript"},
			expectedSpecific: []string{"framer-motion", "react-intersection-observer"},
		},
		{
			templateType:     "nextjs-admin",
			expectedBaseDeps: []string{"next", "react", "react-dom", "typescript"},
			expectedSpecific: []string{"@tanstack/react-table", "react-hook-form", "zod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			deps := generator.getDependenciesForTemplate(tc.templateType)

			if len(deps) == 0 {
				t.Error("getDependenciesForTemplate should return non-empty map")
			}

			// Check base dependencies
			for _, dep := range tc.expectedBaseDeps {
				if _, exists := deps[dep]; !exists {
					t.Errorf("Base dependency '%s' is missing for template type '%s'", dep, tc.templateType)
				}
			}

			// Check template-specific dependencies
			for _, dep := range tc.expectedSpecific {
				if _, exists := deps[dep]; !exists {
					t.Errorf("Template-specific dependency '%s' is missing for template type '%s'", dep, tc.templateType)
				}
			}
		})
	}
}

func TestGenerateAllStandardizedFiles(t *testing.T) {
	generator := NewTemplateGenerator()

	templateTypes := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	for _, templateType := range templateTypes {
		t.Run(templateType, func(t *testing.T) {
			files := generator.GenerateAllStandardizedFiles(templateType)

			if files == nil {
				t.Fatal("GenerateAllStandardizedFiles should not return nil")
			}

			// Check that all files are generated
			if files.PackageJSON == "" {
				t.Error("PackageJSON should not be empty")
			}

			if files.TSConfig == "" {
				t.Error("TSConfig should not be empty")
			}

			if files.ESLintConfig == "" {
				t.Error("ESLintConfig should not be empty")
			}

			if files.PrettierConfig == "" {
				t.Error("PrettierConfig should not be empty")
			}

			if files.VercelConfig == "" {
				t.Error("VercelConfig should not be empty")
			}

			if files.TailwindConfig == "" {
				t.Error("TailwindConfig should not be empty")
			}

			if files.NextConfig == "" {
				t.Error("NextConfig should not be empty")
			}

			if files.PostCSSConfig == "" {
				t.Error("PostCSSConfig should not be empty")
			}

			if files.JestConfig == "" {
				t.Error("JestConfig should not be empty")
			}

			if files.JestSetup == "" {
				t.Error("JestSetup should not be empty")
			}
		})
	}
}

func TestGetStandardizedFilePaths(t *testing.T) {
	paths := GetStandardizedFilePaths()

	if len(paths) == 0 {
		t.Error("GetStandardizedFilePaths should return non-empty map")
	}

	expectedPaths := map[string]string{
		"PackageJSON":    "package.json.tmpl",
		"TSConfig":       "tsconfig.json.tmpl",
		"ESLintConfig":   ".eslintrc.json.tmpl",
		"PrettierConfig": ".prettierrc.tmpl",
		"VercelConfig":   "vercel.json.tmpl",
		"TailwindConfig": "tailwind.config.js.tmpl",
		"NextConfig":     "next.config.js.tmpl",
		"PostCSSConfig":  "postcss.config.js.tmpl",
		"JestConfig":     "jest.config.js.tmpl",
		"JestSetup":      "jest.setup.js.tmpl",
	}

	for key, expectedPath := range expectedPaths {
		if actualPath, exists := paths[key]; !exists {
			t.Errorf("Path for '%s' is missing", key)
		} else if actualPath != expectedPath {
			t.Errorf("Path for '%s' should be '%s' but is '%s'", key, expectedPath, actualPath)
		}
	}
}

func TestGetTemplateDescription(t *testing.T) {
	testCases := []struct {
		templateType string
		expected     string
	}{
		{"nextjs-app", "Main Application"},
		{"nextjs-home", "Landing Page"},
		{"nextjs-admin", "Admin Dashboard"},
		{"unknown-template", "Frontend Application"},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			result := getTemplateDescription(tc.templateType)
			if result != tc.expected {
				t.Errorf("getTemplateDescription('%s') = '%s', expected '%s'", tc.templateType, result, tc.expected)
			}
		})
	}
}
