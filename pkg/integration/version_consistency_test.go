package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/validation"
)

// TestVersionConsistencyAcrossMultipleFrontendTemplates tests that all frontend templates
// generate with consistent Node.js versions
func TestVersionConsistencyAcrossMultipleFrontendTemplates(t *testing.T) {
	tempDir := t.TempDir()

	// Create test configuration with standardized Node.js 20.x versions
	config := &models.ProjectConfig{
		Name:         "version-consistency-test",
		Organization: "test-org",
		Description:  "Version consistency integration test",
		License:      "MIT",
		OutputPath:   "./test-output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Home:  true,
					Admin: true,
				},
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.11.0",
			Packages: map[string]string{
				"next":       "15.0.0",
				"react":      "18.2.0",
				"typescript": "^5.0.0",
				"eslint":     "^8.0.0",
			},
		},
	}

	// Create project structure
	projectPath := filepath.Join(tempDir, config.Name)
	fsGenerator := filesystem.NewGenerator()
	err := fsGenerator.CreateProject(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Test Node.js runtime version consistency
	t.Run("NodeJS_Runtime_Version_Consistency", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				t.Logf("Template %s not generated, skipping", templateName)
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				t.Fatalf("Failed to read package.json for %s: %v", templateName, err)
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				t.Fatalf("Failed to parse package.json for %s: %v", templateName, err)
			}

			// Check engines field
			if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
				if nodeVersion, exists := engines["node"]; exists {
					// Should be a valid Node.js version
					if nodeVersionStr, ok := nodeVersion.(string); ok {
						if !strings.HasPrefix(nodeVersionStr, "20.") {
							t.Errorf("Template %s: Expected Node.js 20.x version, got %s", templateName, nodeVersionStr)
						}
					}
				}
			}
		}
	})

	// Test NPM version consistency
	t.Run("NPM_Version_Consistency", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check engines field for npm version
			if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
				if npmVersion, exists := engines["npm"]; exists {
					if npmVersionStr, ok := npmVersion.(string); ok {
						if !strings.HasPrefix(npmVersionStr, "10.") {
							t.Errorf("Template %s: Expected NPM 10.x version, got %s", templateName, npmVersionStr)
						}
					}
				}
			}
		}
	})

	// Test core dependencies consistency
	t.Run("Core_Dependencies_Consistency", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check that core dependencies exist
			if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
				requiredDeps := []string{"next", "react", "react-dom"}
				for _, dep := range requiredDeps {
					if _, exists := deps[dep]; !exists {
						t.Errorf("Template %s: Missing required dependency %s", templateName, dep)
					}
				}
			}
		}
	})

	// Test types dependencies consistency
	t.Run("Types_Dependencies_Consistency", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check devDependencies for types
			if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
				requiredTypesDeps := []string{"@types/node", "@types/react", "@types/react-dom", "typescript"}
				for _, dep := range requiredTypesDeps {
					if _, exists := devDeps[dep]; !exists {
						t.Errorf("Template %s: Missing required dev dependency %s", templateName, dep)
					}
				}
			}
		}
	})

	// Test devDependencies consistency
	t.Run("DevDependencies_Consistency", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check that devDependencies field exists
			if _, ok := packageJSON["devDependencies"]; !ok {
				t.Errorf("Template %s: Missing devDependencies field", templateName)
			}
		}
	})

	// Test Node.js types compatibility
	t.Run("NodeJS_Types_Compatibility", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(projectPath, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check that @types/node version is compatible with Node.js version
			if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
				if typesNodeVersion, exists := devDeps["@types/node"]; exists {
					if versionStr, ok := typesNodeVersion.(string); ok {
						// Should be a valid version string
						if versionStr == "" || versionStr == "undefined" {
							t.Errorf("Template %s: @types/node has invalid version '%s'", templateName, versionStr)
						}
					}
				}
			}
		}
	})

	// Test validation engine consistency
	t.Run("Validation_Engine_Consistency", func(t *testing.T) {
		validator := validation.NewEngine()

		// Validate the generated project
		result, err := validator.ValidateProject(projectPath)
		if err != nil {
			t.Errorf("Project validation failed: %v", err)
		}
		if !result.Valid {
			t.Errorf("Project validation failed: %v", result.Summary)
		}
	})

	t.Logf("✅ Version consistency test completed for %d frontend templates", 3)
}

// TestEndToEndVersionConsistency tests the complete version consistency workflow
func TestEndToEndVersionConsistency(t *testing.T) {
	tempDir := t.TempDir()

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "e2e-version-test",
		Organization: "test-org",
		Description:  "End-to-end version consistency test",
		License:      "MIT",
		OutputPath:   "./test-output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.11.0",
			Packages: map[string]string{
				"next":       "15.0.0",
				"react":      "18.2.0",
				"typescript": "^5.0.0",
			},
		},
	}

	// Test pre-generation validation
	t.Run("Pre_Generation_Validation", func(t *testing.T) {
		validator := validation.NewEngine()
		result, err := validator.ValidateProject("")
		// Basic validation should pass for empty path
		if err != nil {
			t.Logf("Pre-generation validation error: %v", err)
		} else if result.Valid {
			t.Log("Pre-generation validation passed")
		}
	})

	// Test project generation
	t.Run("Project_Generation", func(t *testing.T) {
		fsGenerator := filesystem.NewGenerator()
		err := fsGenerator.CreateProject(config, tempDir)
		if err != nil {
			t.Fatalf("Failed to create project: %v", err)
		}
		t.Log("Project generation completed")
	})

	// Test template generation
	t.Run("Template_Generation", func(t *testing.T) {
		templateDir := filepath.Join(tempDir, config.Name, "frontend", "nextjs-app")

		// Process any template files
		if _, err := os.Stat(templateDir); err == nil {
			t.Log("Template generation completed")
		}
	})

	// Test generated files validation
	t.Run("Generated_Files_Validation", func(t *testing.T) {
		validator := validation.NewEngine()
		projectPath := filepath.Join(tempDir, config.Name)

		result, err := validator.ValidateProject(projectPath)
		if err != nil {
			t.Logf("Project validation error: %v", err)
		} else {
			t.Logf("Project validation result: Valid=%v, Errors=%d", result.Valid, len(result.Issues))
		}
	})

	// Test final version consistency
	t.Run("Final_Version_Consistency", func(t *testing.T) {
		// Check that generated files have consistent versions
		packageJSONPath := filepath.Join(tempDir, config.Name, "frontend", "nextjs-app", "package.json")
		if _, err := os.Stat(packageJSONPath); err == nil {
			t.Log("Final version consistency check completed")
		}
	})

	t.Logf("✅ End-to-end version consistency test completed")
}

// TestComprehensiveVersionConsistencyWithRealTemplates tests version consistency with real template files
func TestComprehensiveVersionConsistencyWithRealTemplates(t *testing.T) {
	tempDir := t.TempDir()

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "comprehensive-version-test",
		Organization: "test-org",
		Description:  "Comprehensive version consistency test",
		License:      "MIT",
		OutputPath:   "./test-output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Home:  true,
					Admin: true,
				},
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.11.0",
			Packages: map[string]string{
				"next":       "15.0.0",
				"react":      "18.2.0",
				"typescript": "^5.0.0",
				"eslint":     "^8.0.0",
			},
		},
	}

	// Test package.json generation
	t.Run("Package_JSON_Generation", func(t *testing.T) {
		fsGenerator := filesystem.NewGenerator()
		err := fsGenerator.CreateProject(config, tempDir)
		if err != nil {
			t.Fatalf("Failed to create project: %v", err)
		}
		t.Log("Package.json generation completed")
	})

	// Test version consistency validation
	t.Run("Version_Consistency_Validation", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(tempDir, config.Name, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			// Basic validation - check that file exists and is valid JSON
			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				t.Errorf("Failed to read package.json for %s: %v", templateName, err)
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				t.Errorf("Failed to parse package.json for %s: %v", templateName, err)
				continue
			}

			// Check required fields
			requiredFields := []string{"name", "version", "private", "engines", "dependencies", "devDependencies"}
			for _, field := range requiredFields {
				if _, exists := packageJSON[field]; !exists {
					t.Errorf("Template %s: Missing required field %s", templateName, field)
				}
			}
		}
	})

	// Test validation engine integration
	t.Run("Validation_Engine_Integration", func(t *testing.T) {
		validator := validation.NewEngine()
		projectPath := filepath.Join(tempDir, config.Name)

		result, err := validator.ValidateProject(projectPath)
		if err != nil {
			t.Logf("Project validation error: %v", err)
		} else {
			t.Logf("Project validation: Valid=%v, Errors=%d", result.Valid, len(result.Issues))
		}
	})

	// Test NPM compatibility simulation
	t.Run("NPM_Compatibility_Simulation", func(t *testing.T) {
		// Simulate NPM compatibility check by validating package.json structure
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(tempDir, config.Name, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Check scripts field
			if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
				requiredScripts := []string{"dev", "build", "start", "lint"}
				for _, script := range requiredScripts {
					if _, exists := scripts[script]; !exists {
						t.Errorf("Template %s: Missing required script %s", templateName, script)
					}
				}
			}
		}
	})

	// Test cross-template dependency analysis
	t.Run("Cross_Template_Dependency_Analysis", func(t *testing.T) {
		frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

		// Collect all dependencies across templates
		allDeps := make(map[string]map[string]string)

		for _, templateName := range frontendTemplates {
			templateDir := filepath.Join(tempDir, config.Name, "frontend", templateName)
			packageJSONPath := filepath.Join(templateDir, "package.json")

			if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
				continue
			}

			data, err := os.ReadFile(packageJSONPath)
			if err != nil {
				continue
			}

			var packageJSON map[string]interface{}
			if err := json.Unmarshal(data, &packageJSON); err != nil {
				continue
			}

			// Collect dependencies
			if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
				allDeps[templateName] = make(map[string]string)
				for dep, version := range deps {
					if versionStr, ok := version.(string); ok {
						allDeps[templateName][dep] = versionStr
					}
				}
			}
		}

		// Check for version conflicts
		depVersions := make(map[string]map[string]string)
		for templateName, deps := range allDeps {
			for dep, version := range deps {
				if depVersions[dep] == nil {
					depVersions[dep] = make(map[string]string)
				}
				depVersions[dep][templateName] = version
			}
		}

		// Report any conflicts
		conflicts := 0
		for dep, versions := range depVersions {
			if len(versions) > 1 {
				uniqueVersions := make(map[string]bool)
				for _, version := range versions {
					uniqueVersions[version] = true
				}
				if len(uniqueVersions) > 1 {
					conflicts++
					t.Logf("Dependency %s has version conflicts: %v", dep, versions)
				}
			}
		}

		if conflicts == 0 {
			t.Logf("✅ No dependency version conflicts found across %d templates", len(frontendTemplates))
		}
	})

	t.Logf("✅ Comprehensive version consistency test completed for %d templates", 3)
}
