package standards

import (
	"encoding/json"
	"testing"
)

func TestGetFrontendStandards(t *testing.T) {
	standards := GetFrontendStandards()

	// Test that standards are not nil
	if standards == nil {
		t.Fatal("GetFrontendStandards() returned nil")
	}

	// Test PackageJSON standards
	if len(standards.PackageJSON.Scripts) == 0 {
		t.Error("PackageJSON.Scripts should not be empty")
	}

	if len(standards.PackageJSON.Dependencies) == 0 {
		t.Error("PackageJSON.Dependencies should not be empty")
	}

	if len(standards.PackageJSON.DevDeps) == 0 {
		t.Error("PackageJSON.DevDeps should not be empty")
	}

	// Test required scripts
	requiredScripts := []string{"dev", "build", "start", "lint", "type-check", "test"}
	for _, script := range requiredScripts {
		if _, exists := standards.PackageJSON.Scripts[script]; !exists {
			t.Errorf("Required script '%s' is missing", script)
		}
	}

	// Test required dependencies
	requiredDeps := []string{"next", "react", "react-dom", "typescript"}
	for _, dep := range requiredDeps {
		if _, exists := standards.PackageJSON.Dependencies[dep]; !exists {
			t.Errorf("Required dependency '%s' is missing", dep)
		}
	}

	// Test engines
	if nodeVersion, exists := standards.PackageJSON.Engines["node"]; !exists || nodeVersion == "" {
		t.Error("Node.js engine version should be specified")
	}

	if npmVersion, exists := standards.PackageJSON.Engines["npm"]; !exists || npmVersion == "" {
		t.Error("npm engine version should be specified")
	}
}

func TestGetTemplateSpecificDependencies(t *testing.T) {
	testCases := []struct {
		templateType string
		expectedDeps []string
	}{
		{
			templateType: "nextjs-app",
			expectedDeps: []string{"@radix-ui/react-dialog", "@radix-ui/react-dropdown-menu", "@radix-ui/react-toast"},
		},
		{
			templateType: "nextjs-home",
			expectedDeps: []string{"@radix-ui/react-accordion", "framer-motion", "react-intersection-observer"},
		},
		{
			templateType: "nextjs-admin",
			expectedDeps: []string{"@tanstack/react-table", "react-hook-form", "zod", "recharts"},
		},
		{
			templateType: "unknown-template",
			expectedDeps: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			deps := GetTemplateSpecificDependencies(tc.templateType)

			if len(tc.expectedDeps) == 0 {
				if len(deps) != 0 {
					t.Errorf("Expected no dependencies for %s, but got %d", tc.templateType, len(deps))
				}
				return
			}

			for _, expectedDep := range tc.expectedDeps {
				if _, exists := deps[expectedDep]; !exists {
					t.Errorf("Expected dependency '%s' not found for template type '%s'", expectedDep, tc.templateType)
				}
			}
		})
	}
}

func TestGeneratePackageJSON(t *testing.T) {
	standards := GetFrontendStandards()

	testCases := []struct {
		templateType string
		name         string
		description  string
	}{
		{"nextjs-app", "test-project", "Test Project Description"},
		{"nextjs-home", "landing-site", "Landing Site Description"},
		{"nextjs-admin", "admin-panel", "Admin Panel Description"},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			jsonBytes, err := standards.GeneratePackageJSON(tc.templateType, tc.name, tc.description)
			if err != nil {
				t.Fatalf("GeneratePackageJSON failed: %v", err)
			}

			// Validate that it's valid JSON
			var pkg map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &pkg); err != nil {
				t.Fatalf("Generated package.json is not valid JSON: %v", err)
			}

			// Check basic fields
			if name, ok := pkg["name"].(string); !ok || name == "" {
				t.Error("package.json should have a non-empty name field")
			}

			if version, ok := pkg["version"].(string); !ok || version == "" {
				t.Error("package.json should have a non-empty version field")
			}

			if private, ok := pkg["private"].(bool); !ok || !private {
				t.Error("package.json should have private set to true")
			}

			// Check scripts
			scripts, ok := pkg["scripts"].(map[string]interface{})
			if !ok {
				t.Fatal("package.json should have scripts field")
			}

			requiredScripts := []string{"dev", "build", "start", "lint"}
			for _, script := range requiredScripts {
				if _, exists := scripts[script]; !exists {
					t.Errorf("Required script '%s' is missing", script)
				}
			}

			// Check dependencies
			dependencies, ok := pkg["dependencies"].(map[string]interface{})
			if !ok {
				t.Fatal("package.json should have dependencies field")
			}

			requiredDeps := []string{"next", "react", "react-dom"}
			for _, dep := range requiredDeps {
				if _, exists := dependencies[dep]; !exists {
					t.Errorf("Required dependency '%s' is missing", dep)
				}
			}
		})
	}
}

func TestGenerateTSConfig(t *testing.T) {
	standards := GetFrontendStandards()

	jsonBytes, err := standards.GenerateTSConfig()
	if err != nil {
		t.Fatalf("GenerateTSConfig failed: %v", err)
	}

	// Validate that it's valid JSON
	var tsconfig map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &tsconfig); err != nil {
		t.Fatalf("Generated tsconfig.json is not valid JSON: %v", err)
	}

	// Check compiler options
	compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{})
	if !ok {
		t.Fatal("tsconfig.json should have compilerOptions field")
	}

	// Check key compiler options
	requiredOptions := []string{"target", "lib", "strict", "jsx", "moduleResolution"}
	for _, option := range requiredOptions {
		if _, exists := compilerOptions[option]; !exists {
			t.Errorf("Required compiler option '%s' is missing", option)
		}
	}

	// Check paths
	paths, ok := compilerOptions["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("compilerOptions should have paths field")
	}

	if _, exists := paths["@/*"]; !exists {
		t.Error("Path alias '@/*' should be defined")
	}
}

func TestGenerateESLintConfig(t *testing.T) {
	standards := GetFrontendStandards()

	jsonBytes, err := standards.GenerateESLintConfig()
	if err != nil {
		t.Fatalf("GenerateESLintConfig failed: %v", err)
	}

	// Validate that it's valid JSON
	var eslintConfig map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &eslintConfig); err != nil {
		t.Fatalf("Generated .eslintrc.json is not valid JSON: %v", err)
	}

	// Check extends
	extends, ok := eslintConfig["extends"].([]interface{})
	if !ok {
		t.Fatal(".eslintrc.json should have extends field")
	}

	if len(extends) == 0 {
		t.Error("extends should not be empty")
	}

	// Check parser
	parser, ok := eslintConfig["parser"].(string)
	if !ok || parser == "" {
		t.Error(".eslintrc.json should have a non-empty parser field")
	}

	// Check rules
	rules, ok := eslintConfig["rules"].(map[string]interface{})
	if !ok {
		t.Fatal(".eslintrc.json should have rules field")
	}

	if len(rules) == 0 {
		t.Error("rules should not be empty")
	}
}

func TestGeneratePrettierConfig(t *testing.T) {
	standards := GetFrontendStandards()

	jsonBytes, err := standards.GeneratePrettierConfig()
	if err != nil {
		t.Fatalf("GeneratePrettierConfig failed: %v", err)
	}

	// Validate that it's valid JSON
	var prettierConfig map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &prettierConfig); err != nil {
		t.Fatalf("Generated .prettierrc is not valid JSON: %v", err)
	}

	// Check key options
	requiredOptions := []string{"semi", "trailingComma", "singleQuote", "printWidth", "tabWidth"}
	for _, option := range requiredOptions {
		if _, exists := prettierConfig[option]; !exists {
			t.Errorf("Required Prettier option '%s' is missing", option)
		}
	}

	// Check plugins
	plugins, ok := prettierConfig["plugins"].([]interface{})
	if !ok {
		t.Fatal(".prettierrc should have plugins field")
	}

	if len(plugins) == 0 {
		t.Error("plugins should not be empty")
	}
}

func TestGenerateVercelConfig(t *testing.T) {
	standards := GetFrontendStandards()

	jsonBytes, err := standards.GenerateVercelConfig("nextjs-app")
	if err != nil {
		t.Fatalf("GenerateVercelConfig failed: %v", err)
	}

	// Validate that it's valid JSON
	var vercelConfig map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &vercelConfig); err != nil {
		t.Fatalf("Generated vercel.json is not valid JSON: %v", err)
	}

	// Check framework
	framework, ok := vercelConfig["framework"].(string)
	if !ok || framework != "nextjs" {
		t.Error("framework should be 'nextjs'")
	}

	// Check build command
	buildCommand, ok := vercelConfig["buildCommand"].(string)
	if !ok || buildCommand == "" {
		t.Error("buildCommand should be specified")
	}

	// Check headers
	headers, ok := vercelConfig["headers"].([]interface{})
	if !ok {
		t.Fatal("vercel.json should have headers field")
	}

	if len(headers) == 0 {
		t.Error("headers should not be empty")
	}
}

func TestGenerateTailwindConfig(t *testing.T) {
	standards := GetFrontendStandards()

	config := standards.GenerateTailwindConfig()

	if config == "" {
		t.Error("GenerateTailwindConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"tailwindcss",
		"darkMode",
		"content",
		"theme",
		"plugins",
		"tailwindcss-animate",
	}

	for _, content := range expectedContent {
		if !contains(config, content) {
			t.Errorf("Tailwind config should contain '%s'", content)
		}
	}
}

func TestGenerateNextConfig(t *testing.T) {
	standards := GetFrontendStandards()

	config := standards.GenerateNextConfig()

	if config == "" {
		t.Error("GenerateNextConfig should return non-empty string")
	}

	// Check for key content
	expectedContent := []string{
		"nextConfig",
		"typescript",
		"eslint",
		"experimental",
		"images",
	}

	for _, content := range expectedContent {
		if !contains(config, content) {
			t.Errorf("Next.js config should contain '%s'", content)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
