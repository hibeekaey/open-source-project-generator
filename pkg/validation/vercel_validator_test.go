package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewVercelValidator(t *testing.T) {
	validator := NewVercelValidator()
	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}
	if validator.standardConfig == nil {
		t.Fatal("Expected standardConfig to be initialized")
	}
}

func TestValidateVercelConfig(t *testing.T) {
	validator := NewVercelValidator()

	tests := []struct {
		name           string
		vercelConfig   map[string]interface{}
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid vercel config",
			vercelConfig: map[string]interface{}{
				"framework":    "nextjs",
				"buildCommand": "npm run build",
				"headers": []interface{}{
					map[string]interface{}{
						"source": "/(.*)",
						"headers": []interface{}{
							map[string]interface{}{
								"key":   "X-Frame-Options",
								"value": "DENY",
							},
							map[string]interface{}{
								"key":   "X-Content-Type-Options",
								"value": "nosniff",
							},
							map[string]interface{}{
								"key":   "Referrer-Policy",
								"value": "strict-origin-when-cross-origin",
							},
						},
					},
				},
				"functions": map[string]interface{}{
					"api/**/*.ts": map[string]interface{}{
						"maxDuration": 30,
						"memory":      512,
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "config with warnings",
			vercelConfig: map[string]interface{}{
				"framework":    "unknown-framework",
				"buildCommand": "",
				"functions": map[string]interface{}{
					"api/**/*.ts": map[string]interface{}{
						"maxDuration": 400,  // > 300s
						"memory":      2048, // > 1024MB
					},
				},
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // unusual framework, empty build command, no headers, high duration, high memory
		},
		{
			name: "minimal config",
			vercelConfig: map[string]interface{}{
				"buildCommand": "npm run build",
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 2, // no framework, no headers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			vercelConfigPath := filepath.Join(tmpDir, "vercel.json")

			data, err := json.Marshal(tt.vercelConfig)
			if err != nil {
				t.Fatalf("Failed to marshal test data: %v", err)
			}

			if err := os.WriteFile(vercelConfigPath, data, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := validator.ValidateVercelConfig(vercelConfigPath)
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

func TestValidateVercelCompatibility(t *testing.T) {
	validator := NewVercelValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid Next.js project",
			setupProject: func(projectPath string) error {
				// Create package.json
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
				if err := os.WriteFile(filepath.Join(projectPath, "package.json"), data, 0644); err != nil {
					return err
				}

				// Create vercel.json
				vercelConfig := map[string]interface{}{
					"framework":    "nextjs",
					"buildCommand": "npm run build",
				}

				data, _ = json.Marshal(vercelConfig)
				if err := os.WriteFile(filepath.Join(projectPath, "vercel.json"), data, 0644); err != nil {
					return err
				}

				// Create public directory
				if err := os.MkdirAll(filepath.Join(projectPath, "public"), 0755); err != nil {
					return err
				}

				// Create src/app directory
				if err := os.MkdirAll(filepath.Join(projectPath, "src", "app"), 0755); err != nil {
					return err
				}

				return nil
			},
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // missing .vercelignore
		},
		{
			name: "missing required scripts",
			setupProject: func(projectPath string) error {
				packageJSON := map[string]interface{}{
					"name":    "test-app",
					"version": "1.0.0",
					"scripts": map[string]interface{}{
						"dev": "next dev",
						// missing build and start
					},
					"dependencies": map[string]interface{}{
						"next": "15.5.2",
					},
				}

				data, _ := json.Marshal(packageJSON)
				return os.WriteFile(filepath.Join(projectPath, "package.json"), data, 0644)
			},
			expectValid:    false,
			expectErrors:   2, // missing build and start scripts
			expectWarnings: 4, // no vercel.json, no node version, no public dir, no pages/app dir
		},
		{
			name: "missing package.json",
			setupProject: func(projectPath string) error {
				// Don't create package.json
				return nil
			},
			expectValid:    false,
			expectErrors:   1, // missing package.json
			expectWarnings: 3, // no vercel.json, no public dir, no pages/app dir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setupProject(tmpDir); err != nil {
				t.Fatalf("Failed to setup test project: %v", err)
			}

			result, err := validator.ValidateVercelCompatibility(tmpDir)
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

func TestValidateEnvironmentVariablesConsistency(t *testing.T) {
	validator := NewVercelValidator()

	// Create temporary template structure
	tmpDir := t.TempDir()
	frontendDir := filepath.Join(tmpDir, "frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		t.Fatalf("Failed to create frontend directory: %v", err)
	}

	// Create test templates with different env vars
	templates := []struct {
		name         string
		vercelConfig map[string]interface{}
		envExample   string
	}{
		{
			name: "nextjs-app",
			vercelConfig: map[string]interface{}{
				"env": map[string]interface{}{
					"NEXT_PUBLIC_APP_NAME": "{{.Name}}",
					"DATABASE_URL":         "{{.DatabaseURL}}",
				},
			},
			envExample: "NEXT_PUBLIC_APP_NAME=myapp\nDATABASE_URL=postgres://localhost\nAPI_KEY=secret\n",
		},
		{
			name: "nextjs-home",
			vercelConfig: map[string]interface{}{
				"env": map[string]interface{}{
					"NEXT_PUBLIC_APP_NAME": "{{.Name}}",
					// missing DATABASE_URL
				},
			},
			envExample: "NEXT_PUBLIC_APP_NAME=myapp\nAPI_KEY=secret\n",
		},
	}

	// Create template directories and files
	for _, tmpl := range templates {
		templateDir := filepath.Join(frontendDir, tmpl.name)
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		// Create vercel.json.tmpl
		vercelConfigPath := filepath.Join(templateDir, "vercel.json.tmpl")
		data, err := json.Marshal(tmpl.vercelConfig)
		if err != nil {
			t.Fatalf("Failed to marshal vercel config: %v", err)
		}

		if err := os.WriteFile(vercelConfigPath, data, 0644); err != nil {
			t.Fatalf("Failed to write vercel config: %v", err)
		}

		// Create .env.example.tmpl
		envExamplePath := filepath.Join(templateDir, ".env.example.tmpl")
		if err := os.WriteFile(envExamplePath, []byte(tmpl.envExample), 0644); err != nil {
			t.Fatalf("Failed to write env example: %v", err)
		}
	}

	result, err := validator.ValidateEnvironmentVariablesConsistency(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected validation to pass, but got errors: %v", result.Errors)
	}

	// Should have warnings about inconsistent env vars
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings about inconsistent environment variables")
	}
}

func TestValidateSecurityHeaders(t *testing.T) {
	validator := NewVercelValidator()

	tests := []struct {
		name           string
		headers        []interface{}
		expectWarnings int
	}{
		{
			name: "all security headers present",
			headers: []interface{}{
				map[string]interface{}{
					"source": "/(.*)",
					"headers": []interface{}{
						map[string]interface{}{
							"key":   "X-Frame-Options",
							"value": "DENY",
						},
						map[string]interface{}{
							"key":   "X-Content-Type-Options",
							"value": "nosniff",
						},
						map[string]interface{}{
							"key":   "Referrer-Policy",
							"value": "strict-origin-when-cross-origin",
						},
					},
				},
			},
			expectWarnings: 0,
		},
		{
			name: "missing security headers",
			headers: []interface{}{
				map[string]interface{}{
					"source": "/(.*)",
					"headers": []interface{}{
						map[string]interface{}{
							"key":   "X-Frame-Options",
							"value": "DENY",
						},
						// missing X-Content-Type-Options and Referrer-Policy
					},
				},
			},
			expectWarnings: 2,
		},
		{
			name:           "no headers",
			headers:        []interface{}{},
			expectWarnings: 3, // all required headers missing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			validator.validateSecurityHeaders(tt.headers, result)

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

func TestValidateFunctionsConfig(t *testing.T) {
	validator := NewVercelValidator()

	tests := []struct {
		name           string
		functions      map[string]interface{}
		expectWarnings int
	}{
		{
			name: "valid functions config",
			functions: map[string]interface{}{
				"api/**/*.ts": map[string]interface{}{
					"maxDuration": 30,
					"memory":      512,
				},
			},
			expectWarnings: 0,
		},
		{
			name: "functions exceeding limits",
			functions: map[string]interface{}{
				"api/**/*.ts": map[string]interface{}{
					"maxDuration": 400,  // > 300s
					"memory":      2048, // > 1024MB
				},
			},
			expectWarnings: 2,
		},
		{
			name: "mixed functions",
			functions: map[string]interface{}{
				"api/fast/*.ts": map[string]interface{}{
					"maxDuration": 30,
					"memory":      256,
				},
				"api/slow/*.ts": map[string]interface{}{
					"maxDuration": 500,  // > 300s
					"memory":      1536, // > 1024MB
				},
			},
			expectWarnings: 2, // one function exceeds both limits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			validator.validateFunctionsConfig(tt.functions, result)

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tt.expectWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

func TestIsValidNodeVersionForVercel(t *testing.T) {
	validator := NewVercelValidator()

	tests := []struct {
		version string
		valid   bool
	}{
		{">=18.0.0", true},
		{">=20.0.0", true},
		{">=22.0.0", true},
		{"18.17.0", true},
		{"20.10.0", true},
		{"22.0.0", true},
		{">=16.0.0", false},
		{"14.21.0", false},
		{">=24.0.0", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := validator.isValidNodeVersionForVercel(tt.version)
			if result != tt.valid {
				t.Errorf("Expected %v for version %q, got %v", tt.valid, tt.version, result)
			}
		})
	}
}

func TestExtractEnvVarsFromVercelConfig(t *testing.T) {
	validator := NewVercelValidator()

	vercelConfig := map[string]interface{}{
		"env": map[string]interface{}{
			"NEXT_PUBLIC_APP_NAME": "myapp",
			"DATABASE_URL":         "postgres://localhost",
		},
		"build": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV":                "production",
				"NEXT_TELEMETRY_DISABLED": "1",
			},
		},
	}

	tmpDir := t.TempDir()
	vercelConfigPath := filepath.Join(tmpDir, "vercel.json")

	data, err := json.Marshal(vercelConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(vercelConfigPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	envVars, err := validator.extractEnvVarsFromVercelConfig(vercelConfigPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedVars := []string{"NEXT_PUBLIC_APP_NAME", "DATABASE_URL", "NODE_ENV", "NEXT_TELEMETRY_DISABLED"}
	if len(envVars) != len(expectedVars) {
		t.Errorf("Expected %d env vars, got %d: %v", len(expectedVars), len(envVars), envVars)
	}

	// Check that all expected vars are present
	varMap := make(map[string]bool)
	for _, v := range envVars {
		varMap[v] = true
	}

	for _, expected := range expectedVars {
		if !varMap[expected] {
			t.Errorf("Expected env var %s not found in result", expected)
		}
	}
}

func TestExtractEnvVarsFromEnvFile(t *testing.T) {
	validator := NewVercelValidator()

	envContent := `# This is a comment
NEXT_PUBLIC_APP_NAME=myapp
DATABASE_URL=postgres://localhost:5432/mydb
API_KEY=secret123

# Another comment
REDIS_URL=redis://localhost:6379
`

	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env.example")

	if err := os.WriteFile(envFilePath, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	envVars, err := validator.extractEnvVarsFromEnvFile(envFilePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedVars := []string{"NEXT_PUBLIC_APP_NAME", "DATABASE_URL", "API_KEY", "REDIS_URL"}
	if len(envVars) != len(expectedVars) {
		t.Errorf("Expected %d env vars, got %d: %v", len(expectedVars), len(envVars), envVars)
	}

	// Check that all expected vars are present
	varMap := make(map[string]bool)
	for _, v := range envVars {
		varMap[v] = true
	}

	for _, expected := range expectedVars {
		if !varMap[expected] {
			t.Errorf("Expected env var %s not found in result", expected)
		}
	}
}
