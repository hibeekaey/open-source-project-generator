package filesystem

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func createTestProjectConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project description",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Home:  true,
					Admin: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "18.17.0",
			Go:   "1.22.0",
			Packages: map[string]string{
				"kotlin": "1.9.0",
				"swift":  "5.9.0",
				"next":   "14.0.0",
				"react":  "18.2.0",
			},
		},
		OutputPath:       "/tmp",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}
}

func TestNewProjectGenerator(t *testing.T) {
	pg := NewProjectGenerator()
	if pg == nil {
		t.Fatal("NewProjectGenerator() returned nil")
	}

	if pg.fsGen == nil {
		t.Fatal("NewProjectGenerator() did not initialize filesystem generator")
	}

	if pg.structure == nil {
		t.Fatal("NewProjectGenerator() did not initialize project structure")
	}
}

func TestNewDryRunProjectGenerator(t *testing.T) {
	pg := NewDryRunProjectGenerator()
	if pg == nil {
		t.Fatal("NewDryRunProjectGenerator() returned nil")
	}

	if pg.fsGen == nil {
		t.Fatal("NewDryRunProjectGenerator() did not initialize filesystem generator")
	}

	if !pg.fsGen.dryRun {
		t.Fatal("NewDryRunProjectGenerator() did not set dry-run mode")
	}
}

func TestGetStandardProjectStructure(t *testing.T) {
	structure := GetStandardProjectStructure()
	if structure == nil {
		t.Fatal("GetStandardProjectStructure() returned nil")
	}

	// Verify root directories are defined
	if len(structure.RootDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty RootDirs")
	}

	// Verify component directories are defined
	if len(structure.FrontendDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty FrontendDirs")
	}

	if len(structure.BackendDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty BackendDirs")
	}

	if len(structure.MobileDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty MobileDirs")
	}

	if len(structure.InfraDirs) == 0 {
		t.Error("GetStandardProjectStructure() returned empty InfraDirs")
	}

	// Verify expected directories exist
	expectedRootDirs := []string{"docs", "scripts", ".github/workflows"}
	for _, expectedDir := range expectedRootDirs {
		found := false
		for _, dir := range structure.RootDirs {
			if dir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetStandardProjectStructure() missing expected root directory: %s", expectedDir)
		}
	}
}

func TestGenerateProjectStructure(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		outputPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project structure generation",
			config:      config,
			outputPath:  tempDir,
			expectError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			outputPath:  tempDir,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
		{
			name:        "empty output path",
			config:      config,
			outputPath:  "",
			expectError: true,
			errorMsg:    "output path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewProjectGenerator()
			err := pg.GenerateProjectStructure(tt.config, tt.outputPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateProjectStructure() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("GenerateProjectStructure() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("GenerateProjectStructure() unexpected error = %v", err)
				}

				// Verify project directory was created
				projectPath := filepath.Join(tt.outputPath, tt.config.Name)
				if !pg.fsGen.FileExists(projectPath) {
					t.Errorf("GenerateProjectStructure() did not create project directory")
				}

				// Verify root directories were created
				for _, dir := range pg.structure.RootDirs {
					dirPath := filepath.Join(projectPath, dir)
					if !pg.fsGen.FileExists(dirPath) {
						t.Errorf("GenerateProjectStructure() did not create root directory: %s", dir)
					}
				}

				// Verify component directories were created based on config
				if tt.config.Components.Frontend.NextJS.App {
					for _, dir := range pg.structure.FrontendDirs {
						dirPath := filepath.Join(projectPath, dir)
						if !pg.fsGen.FileExists(dirPath) {
							t.Errorf("GenerateProjectStructure() did not create frontend directory: %s", dir)
						}
					}
				}

				if tt.config.Components.Backend.GoGin {
					for _, dir := range pg.structure.BackendDirs {
						dirPath := filepath.Join(projectPath, dir)
						if !pg.fsGen.FileExists(dirPath) {
							t.Errorf("GenerateProjectStructure() did not create backend directory: %s", dir)
						}
					}
				}
			}
		})
	}
}

func TestGenerateComponentFiles(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	// First create the project structure
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		outputPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid component files generation",
			config:      config,
			outputPath:  tempDir,
			expectError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			outputPath:  tempDir,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pg.GenerateComponentFiles(tt.config, tt.outputPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateComponentFiles() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("GenerateComponentFiles() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("GenerateComponentFiles() unexpected error = %v", err)
				}

				projectPath := filepath.Join(tt.outputPath, tt.config.Name)

				// Verify root files were created
				rootFiles := []string{"Makefile", "README.md", "docker-compose.yml", ".gitignore"}
				for _, file := range rootFiles {
					filePath := filepath.Join(projectPath, file)
					if !pg.fsGen.FileExists(filePath) {
						t.Errorf("GenerateComponentFiles() did not create root file: %s", file)
					}
				}

				// Verify component-specific files were created
				if tt.config.Components.Frontend.NextJS.App {
					packageJsonPath := filepath.Join(projectPath, "App/package.json")
					if !pg.fsGen.FileExists(packageJsonPath) {
						t.Errorf("GenerateComponentFiles() did not create App/package.json")
					}
				}

				if tt.config.Components.Backend.GoGin {
					goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
					if !pg.fsGen.FileExists(goModPath) {
						t.Errorf("GenerateComponentFiles() did not create CommonServer/go.mod")
					}
				}

				if tt.config.Components.Mobile.Android {
					buildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
					if !pg.fsGen.FileExists(buildGradlePath) {
						t.Errorf("GenerateComponentFiles() did not create Mobile/Android/build.gradle")
					}
				}

				if tt.config.Components.Mobile.IOS {
					packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
					if !pg.fsGen.FileExists(packageSwiftPath) {
						t.Errorf("GenerateComponentFiles() did not create Mobile/iOS/Package.swift")
					}
				}
			}
		})
	}
}

func TestValidateProjectStructure(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	// Create a complete project structure
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	tests := []struct {
		name        string
		projectPath string
		config      *models.ProjectConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project structure validation",
			projectPath: projectPath,
			config:      config,
			expectError: false,
		},
		{
			name:        "empty project path",
			projectPath: "",
			config:      config,
			expectError: true,
			errorMsg:    "project path cannot be empty",
		},
		{
			name:        "nil config",
			projectPath: projectPath,
			config:      nil,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
		{
			name:        "non-existent project path",
			projectPath: filepath.Join(tempDir, "non-existent"),
			config:      config,
			expectError: true,
			errorMsg:    "required root directory missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pg.ValidateProjectStructure(tt.projectPath, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateProjectStructure() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateProjectStructure() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateProjectStructure() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateCrossReferences(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	// Create a complete project with files
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
		t.Fatalf("Failed to create component files: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	tests := []struct {
		name        string
		projectPath string
		config      *models.ProjectConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid cross-reference validation",
			projectPath: projectPath,
			config:      config,
			expectError: false,
		},
		{
			name:        "empty project path",
			projectPath: "",
			config:      config,
			expectError: true,
			errorMsg:    "project path cannot be empty",
		},
		{
			name:        "nil config",
			projectPath: projectPath,
			config:      nil,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pg.ValidateCrossReferences(tt.projectPath, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateCrossReferences() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateCrossReferences() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCrossReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComponentSelectionGeneration(t *testing.T) {
	tempDir := t.TempDir()
	pg := NewProjectGenerator()

	tests := []struct {
		name       string
		components models.Components
		checkDirs  []string
		checkFiles []string
	}{
		{
			name: "frontend only",
			components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
			checkDirs:  []string{"App/src/components/ui", "App/src/hooks"},
			checkFiles: []string{"App/package.json"},
		},
		{
			name: "backend only",
			components: models.Components{
				Backend: models.BackendComponents{GoGin: true},
			},
			checkDirs:  []string{"CommonServer/internal/controllers", "CommonServer/pkg/auth"},
			checkFiles: []string{"CommonServer/go.mod"},
		},
		{
			name: "mobile only",
			components: models.Components{
				Mobile: models.MobileComponents{Android: true, IOS: true},
			},
			checkDirs:  []string{"Mobile/Android/app/src/main/java", "Mobile/iOS/Sources"},
			checkFiles: []string{"Mobile/Android/build.gradle", "Mobile/iOS/Package.swift"},
		},
		{
			name: "infrastructure only",
			components: models.Components{
				Infrastructure: models.InfrastructureComponents{Terraform: true},
			},
			checkDirs:  []string{"Deploy/terraform/modules", "Deploy/kubernetes/base"},
			checkFiles: []string{"Deploy/terraform/main.tf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := createTestProjectConfig()
			config.Components = tt.components
			config.Name = "test-" + strings.ReplaceAll(tt.name, " ", "-")

			// Generate project structure
			if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
				t.Fatalf("Failed to generate project structure: %v", err)
			}

			// Generate component files
			if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
				t.Fatalf("Failed to generate component files: %v", err)
			}

			projectPath := filepath.Join(tempDir, config.Name)

			// Check that expected directories exist
			for _, dir := range tt.checkDirs {
				dirPath := filepath.Join(projectPath, dir)
				if !pg.fsGen.FileExists(dirPath) {
					t.Errorf("Expected directory not created: %s", dir)
				}
			}

			// Check that expected files exist
			for _, file := range tt.checkFiles {
				filePath := filepath.Join(projectPath, file)
				if !pg.fsGen.FileExists(filePath) {
					t.Errorf("Expected file not created: %s", file)
				}
			}

			// Validate the project structure
			if err := pg.ValidateProjectStructure(projectPath, config); err != nil {
				t.Errorf("Project structure validation failed: %v", err)
			}

			// Validate cross-references
			if err := pg.ValidateCrossReferences(projectPath, config); err != nil {
				t.Errorf("Cross-reference validation failed: %v", err)
			}
		})
	}
}

func TestProjectGeneratorDryRunMode(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewDryRunProjectGenerator()

	// These operations should not fail in dry-run mode
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Errorf("GenerateProjectStructure() in dry-run mode failed: %v", err)
	}

	if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
		t.Errorf("GenerateComponentFiles() in dry-run mode failed: %v", err)
	}

	// Verify that nothing was actually created
	projectPath := filepath.Join(tempDir, config.Name)
	if pg.fsGen.FileExists(projectPath) {
		t.Errorf("GenerateProjectStructure() in dry-run mode actually created directory")
	}
}

func TestContentGeneration(t *testing.T) {
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	tests := []struct {
		name     string
		method   func(*models.ProjectConfig) string
		contains []string
	}{
		{
			name:   "makefile content",
			method: pg.generateMakefileContent,
			contains: []string{
				config.Name,
				"help:",
				"setup:",
				"dev:",
				"test:",
				"build:",
			},
		},
		{
			name:   "readme content",
			method: pg.generateReadmeContent,
			contains: []string{
				config.Name,
				config.Description,
				config.Organization,
				"## Getting Started",
				"## Project Structure",
			},
		},
		{
			name:   "docker-compose content",
			method: pg.generateDockerComposeContent,
			contains: []string{
				"version: '3.8'",
				"services:",
				config.Name,
			},
		},
		{
			name:   "gitignore content",
			method: pg.generateGitignoreContent,
			contains: []string{
				"node_modules/",
				"*.log",
				".env",
				".DS_Store",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := tt.method(config)
			if content == "" {
				t.Errorf("%s generated empty content", tt.name)
			}

			for _, expected := range tt.contains {
				if !strings.Contains(content, expected) {
					t.Errorf("%s content does not contain expected string: %s", tt.name, expected)
				}
			}
		})
	}
}

// Integration tests for complete project generation

func TestCompleteProjectGeneration(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	// Test complete project generation workflow
	t.Run("complete project generation workflow", func(t *testing.T) {
		// Step 1: Generate project structure
		if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
			t.Fatalf("Failed to generate project structure: %v", err)
		}

		// Step 2: Generate component files
		if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
			t.Fatalf("Failed to generate component files: %v", err)
		}

		projectPath := filepath.Join(tempDir, config.Name)

		// Step 3: Validate project structure
		if err := pg.ValidateProjectStructure(projectPath, config); err != nil {
			t.Fatalf("Project structure validation failed: %v", err)
		}

		// Step 4: Validate cross-references
		if err := pg.ValidateCrossReferences(projectPath, config); err != nil {
			t.Fatalf("Cross-reference validation failed: %v", err)
		}

		// Verify all expected files and directories exist
		expectedStructure := map[string]bool{
			// Root files
			"Makefile":           true,
			"README.md":          true,
			"CONTRIBUTING.md":    true,
			"SECURITY.md":        true,
			"docker-compose.yml": true,
			".gitignore":         true,

			// Root directories
			"docs":                          true,
			"scripts":                       true,
			".github/workflows":             true,
			".github/ISSUE_TEMPLATE":        true,
			".github/PULL_REQUEST_TEMPLATE": true,

			// Frontend files and directories
			"App/package.json":       true,
			"App/next.config.js":     true,
			"App/tailwind.config.js": true,
			"App/src/components/ui":  true,
			"App/src/hooks":          true,
			"Home/package.json":      true,
			"Admin/package.json":     true,

			// Backend files and directories
			"CommonServer/go.mod":               true,
			"CommonServer/internal/controllers": true,
			"CommonServer/internal/models":      true,
			"CommonServer/pkg/auth":             true,

			// Mobile files and directories
			"Mobile/Android/build.gradle":      true,
			"Mobile/iOS/Package.swift":         true,
			"Mobile/Android/app/src/main/java": true,
			"Mobile/iOS/Sources":               true,

			// Infrastructure files and directories
			"Deploy/terraform/main.tf": true,
			"Deploy/terraform/modules": true,
			"Deploy/kubernetes/base":   true,

			// CI/CD files
			".github/workflows/ci.yml":       true,
			".github/workflows/security.yml": true,
			".github/dependabot.yml":         true,

			// Test directories
			"Tests/integration": true,
			"Tests/e2e":         true,
		}

		for expectedPath, shouldExist := range expectedStructure {
			fullPath := filepath.Join(projectPath, expectedPath)
			exists := pg.fsGen.FileExists(fullPath)

			if shouldExist && !exists {
				t.Errorf("Expected path does not exist: %s", expectedPath)
			} else if !shouldExist && exists {
				t.Errorf("Unexpected path exists: %s", expectedPath)
			}
		}
	})
}

func TestProjectGenerationWithDifferentComponentCombinations(t *testing.T) {
	tempDir := t.TempDir()
	pg := NewProjectGenerator()

	testCases := []struct {
		name          string
		components    models.Components
		expectedFiles []string
		expectedDirs  []string
	}{
		{
			name: "minimal frontend only",
			components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
			expectedFiles: []string{
				"App/package.json",
				"App/next.config.js",
				"Makefile",
				"README.md",
			},
			expectedDirs: []string{
				"App/src/components/ui",
				"App/src/hooks",
				"docs",
				"scripts",
			},
		},
		{
			name: "full stack with mobile",
			components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App:   true,
						Home:  true,
						Admin: true,
					},
				},
				Backend: models.BackendComponents{GoGin: true},
				Mobile: models.MobileComponents{
					Android: true,
					IOS:     true,
				},
				Infrastructure: models.InfrastructureComponents{
					Terraform:  true,
					Kubernetes: true,
					Docker:     true,
				},
			},
			expectedFiles: []string{
				"App/package.json",
				"Home/package.json",
				"Admin/package.json",
				"CommonServer/go.mod",
				"Mobile/Android/build.gradle",
				"Mobile/iOS/Package.swift",
				"Deploy/terraform/main.tf",
			},
			expectedDirs: []string{
				"App/src/components/ui",
				"CommonServer/internal/controllers",
				"Mobile/Android/app/src/main/java",
				"Mobile/iOS/Sources",
				"Deploy/terraform/modules",
				"Deploy/kubernetes/base",
			},
		},
		{
			name: "backend and infrastructure only",
			components: models.Components{
				Backend: models.BackendComponents{GoGin: true},
				Infrastructure: models.InfrastructureComponents{
					Terraform: true,
					Docker:    true,
				},
			},
			expectedFiles: []string{
				"CommonServer/go.mod",
				"Deploy/terraform/main.tf",
				"Makefile",
				"README.md",
			},
			expectedDirs: []string{
				"CommonServer/internal/controllers",
				"CommonServer/pkg/auth",
				"Deploy/terraform/modules",
				"Deploy/docker",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := createTestProjectConfig()
			config.Name = "test-" + strings.ReplaceAll(tc.name, " ", "-")
			config.Components = tc.components

			// Generate project
			if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
				t.Fatalf("Failed to generate project structure: %v", err)
			}

			if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
				t.Fatalf("Failed to generate component files: %v", err)
			}

			projectPath := filepath.Join(tempDir, config.Name)

			// Validate structure
			if err := pg.ValidateProjectStructure(projectPath, config); err != nil {
				t.Fatalf("Project structure validation failed: %v", err)
			}

			if err := pg.ValidateCrossReferences(projectPath, config); err != nil {
				t.Fatalf("Cross-reference validation failed: %v", err)
			}

			// Check expected files exist
			for _, expectedFile := range tc.expectedFiles {
				filePath := filepath.Join(projectPath, expectedFile)
				if !pg.fsGen.FileExists(filePath) {
					t.Errorf("Expected file does not exist: %s", expectedFile)
				}
			}

			// Check expected directories exist
			for _, expectedDir := range tc.expectedDirs {
				dirPath := filepath.Join(projectPath, expectedDir)
				if !pg.fsGen.FileExists(dirPath) {
					t.Errorf("Expected directory does not exist: %s", expectedDir)
				}
			}
		})
	}
}

func TestCrossReferenceValidation(t *testing.T) {
	tempDir := t.TempDir()
	config := createTestProjectConfig()
	pg := NewProjectGenerator()

	// Generate complete project
	if err := pg.GenerateProjectStructure(config, tempDir); err != nil {
		t.Fatalf("Failed to generate project structure: %v", err)
	}

	if err := pg.GenerateComponentFiles(config, tempDir); err != nil {
		t.Fatalf("Failed to generate component files: %v", err)
	}

	projectPath := filepath.Join(tempDir, config.Name)

	t.Run("validate required root files exist", func(t *testing.T) {
		requiredFiles := []string{
			"Makefile",
			"README.md",
			"docker-compose.yml",
			".gitignore",
		}

		for _, file := range requiredFiles {
			filePath := filepath.Join(projectPath, file)
			if !pg.fsGen.FileExists(filePath) {
				t.Errorf("Required root file missing: %s", file)
			}
		}
	})

	t.Run("validate component-specific cross-references", func(t *testing.T) {
		// Frontend cross-references
		if config.Components.Frontend.NextJS.App {
			packageJsonPath := filepath.Join(projectPath, "App/package.json")
			if !pg.fsGen.FileExists(packageJsonPath) {
				t.Error("Main app package.json missing")
			}

			nextConfigPath := filepath.Join(projectPath, "App/next.config.js")
			if !pg.fsGen.FileExists(nextConfigPath) {
				t.Error("Main app next.config.js missing")
			}
		}

		// Backend cross-references
		if config.Components.Backend.GoGin {
			goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
			if !pg.fsGen.FileExists(goModPath) {
				t.Error("Backend go.mod missing")
			}
		}

		// Mobile cross-references
		if config.Components.Mobile.Android {
			buildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
			if !pg.fsGen.FileExists(buildGradlePath) {
				t.Error("Android build.gradle missing")
			}
		}

		if config.Components.Mobile.IOS {
			packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
			if !pg.fsGen.FileExists(packageSwiftPath) {
				t.Error("iOS Package.swift missing")
			}
		}

		// Infrastructure cross-references
		if config.Components.Infrastructure.Terraform {
			terraformMainPath := filepath.Join(projectPath, "Deploy/terraform/main.tf")
			if !pg.fsGen.FileExists(terraformMainPath) {
				t.Error("Terraform main.tf missing")
			}
		}
	})

	t.Run("validate CI/CD files exist", func(t *testing.T) {
		cicdFiles := []string{
			".github/workflows/ci.yml",
			".github/workflows/security.yml",
			".github/dependabot.yml",
		}

		for _, file := range cicdFiles {
			filePath := filepath.Join(projectPath, file)
			if !pg.fsGen.FileExists(filePath) {
				t.Errorf("CI/CD file missing: %s", file)
			}
		}
	})

	t.Run("validate documentation files exist", func(t *testing.T) {
		docFiles := []string{
			"README.md",
			"CONTRIBUTING.md",
			"SECURITY.md",
		}

		for _, file := range docFiles {
			filePath := filepath.Join(projectPath, file)
			if !pg.fsGen.FileExists(filePath) {
				t.Errorf("Documentation file missing: %s", file)
			}
		}
	})
}

func TestProjectGenerationErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	pg := NewProjectGenerator()

	t.Run("invalid project configuration", func(t *testing.T) {
		invalidConfigs := []*models.ProjectConfig{
			nil,
			{Name: "", Organization: "test-org"},
			{Name: "test", Organization: ""},
		}

		for i, config := range invalidConfigs {
			err := pg.GenerateProjectStructure(config, tempDir)
			if err == nil {
				t.Errorf("Test case %d: expected error for invalid config, got none", i)
			}
		}
	})

	t.Run("invalid output path", func(t *testing.T) {
		config := createTestProjectConfig()

		err := pg.GenerateProjectStructure(config, "")
		if err == nil {
			t.Error("Expected error for empty output path, got none")
		}
	})

	t.Run("validation with missing directories", func(t *testing.T) {
		config := createTestProjectConfig()
		nonExistentPath := filepath.Join(tempDir, "non-existent-project")

		err := pg.ValidateProjectStructure(nonExistentPath, config)
		if err == nil {
			t.Error("Expected error for non-existent project path, got none")
		}
	})
}
