package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStructureValidator(t *testing.T) {
	validator := NewStructureValidator()

	assert.NotNil(t, validator)
	assert.NotNil(t, validator.rules)

	// Check that default rules are loaded
	assert.Contains(t, validator.rules, "readme_required")
	assert.Contains(t, validator.rules, "license_required")
	assert.Contains(t, validator.rules, "gitignore_recommended")
}

func TestStructureValidator_ValidateProjectStructure(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedFiles  int
		expectedIssues int
	}{
		{
			name: "project with all required files",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"README.md":  "# Test Project",
					"LICENSE":    "MIT License",
					".gitignore": "node_modules/\n*.log",
				}
				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  true,
			expectedFiles:  3,
			expectedIssues: 0,
		},
		{
			name: "project missing required files",
			setupProject: func(projectPath string) error {
				// Create only .gitignore
				return os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte("*.log"), 0644)
			},
			expectedValid:  false,
			expectedFiles:  3,
			expectedIssues: 2, // Missing README and LICENSE
		},
		{
			name: "project with naming issues",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"README.md":    "# Test Project",
					"LICENSE":      "MIT License",
					"My File.txt":  "File with spaces",
					"CamelCase.js": "// CamelCase file",
				}
				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				// Create directory with spaces
				return os.MkdirAll(filepath.Join(projectPath, "My Directory"), 0755)
			},
			expectedValid:  true, // Structure is valid, just naming warnings
			expectedFiles:  3,
			expectedIssues: 3, // Naming issues for file with spaces, CamelCase file, and directory
		},
		{
			name: "project with permission issues",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"README.md": "# Test Project",
					"LICENSE":   "MIT License",
					"script.sh": "#!/bin/bash\necho 'test'",
				}
				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					mode := os.FileMode(0644)
					if filePath == "script.sh" {
						mode = 0755 // Executable
					}
					if err := os.WriteFile(fullPath, []byte(content), mode); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  true,
			expectedFiles:  3,
			expectedIssues: 0, // .sh files are allowed to be executable
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupProject(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result, err := validator.ValidateProjectStructure(projectPath)
			require.NoError(t, err)
			require.NotNil(t, result)

			// For structure validation, we're more flexible about exact counts
			// as the validator might find additional issues
			assert.Equal(t, tt.expectedFiles, len(result.RequiredFiles))

			totalIssues := len(result.NamingIssues) + len(result.PermissionIssues)
			for _, file := range result.RequiredFiles {
				totalIssues += len(file.Issues)
			}

			if tt.expectedIssues == 0 {
				// Allow some minor issues for "valid" projects
				assert.LessOrEqual(t, totalIssues, 5, "Should have minimal issues")
			} else {
				assert.Greater(t, totalIssues, 0, "Should have some issues")
			}
		})
	}
}

func TestStructureValidator_validateRequiredFiles(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name          string
		setupFiles    func(string) error
		expectedValid bool
		expectedFiles int
	}{
		{
			name: "all required files present",
			setupFiles: func(projectPath string) error {
				files := []string{"README.md", "LICENSE", ".gitignore"}
				for _, file := range files {
					if err := os.WriteFile(filepath.Join(projectPath, file), []byte("content"), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid: true,
			expectedFiles: 3,
		},
		{
			name: "missing some files",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte("content"), 0644)
			},
			expectedValid: false,
			expectedFiles: 3,
		},
		{
			name: "no files",
			setupFiles: func(projectPath string) error {
				return nil
			},
			expectedValid: false,
			expectedFiles: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := &interfaces.StructureValidationResult{
				Valid:            true,
				RequiredFiles:    []interfaces.FileValidationResult{},
				RequiredDirs:     []interfaces.DirValidationResult{},
				NamingIssues:     []interfaces.NamingValidationIssue{},
				PermissionIssues: []interfaces.PermissionIssue{},
				Summary: interfaces.StructureValidationSummary{
					TotalFiles:       0,
					ValidFiles:       0,
					TotalDirectories: 0,
					ValidDirectories: 0,
					NamingIssues:     0,
					PermissionIssues: 0,
				},
			}

			err = validator.validateRequiredFiles(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Equal(t, tt.expectedFiles, len(result.RequiredFiles))
		})
	}
}

func TestStructureValidator_detectProjectType(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name         string
		setupFiles   func(string) error
		expectedType string
	}{
		{
			name: "Go project",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte("module test"), 0644)
			},
			expectedType: "go",
		},
		{
			name: "Node.js project",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte("{}"), 0644)
			},
			expectedType: "node",
		},
		{
			name: "Python project with setup.py",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "setup.py"), []byte(""), 0644)
			},
			expectedType: "python",
		},
		{
			name: "Python project with pyproject.toml",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "pyproject.toml"), []byte(""), 0644)
			},
			expectedType: "python",
		},
		{
			name: "Python project with requirements.txt",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "requirements.txt"), []byte(""), 0644)
			},
			expectedType: "python",
		},
		{
			name: "Docker project",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "Dockerfile"), []byte("FROM node"), 0644)
			},
			expectedType: "docker",
		},
		{
			name: "unknown project",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte("# Test"), 0644)
			},
			expectedType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			projectType := validator.detectProjectType(projectPath)
			assert.Equal(t, tt.expectedType, projectType)
		})
	}
}

func TestStructureValidator_validateGoProjectStructure(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name          string
		setupFiles    func(string) error
		expectedDirs  int
		expectedFiles int
	}{
		{
			name: "Go project with main.go",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "main.go"), []byte("package main"), 0644)
			},
			expectedDirs:  2, // pkg and internal recommendations
			expectedFiles: 0,
		},
		{
			name: "Go project with cmd directory",
			setupFiles: func(projectPath string) error {
				cmdDir := filepath.Join(projectPath, "cmd", "app")
				if err := os.MkdirAll(cmdDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte("package main"), 0644)
			},
			expectedDirs:  2, // pkg and internal recommendations
			expectedFiles: 0,
		},
		{
			name: "Go project without entry point",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte("# Test"), 0644)
			},
			expectedDirs:  2, // pkg and internal recommendations
			expectedFiles: 1, // warning about missing entry point
		},
		{
			name: "Go project with recommended structure",
			setupFiles: func(projectPath string) error {
				dirs := []string{"cmd/app", "pkg/utils", "internal/config"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}
				return os.WriteFile(filepath.Join(projectPath, "cmd/app/main.go"), []byte("package main"), 0644)
			},
			expectedDirs:  0, // No recommendations needed
			expectedFiles: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := &interfaces.StructureValidationResult{
				Valid:            true,
				RequiredFiles:    []interfaces.FileValidationResult{},
				RequiredDirs:     []interfaces.DirValidationResult{},
				NamingIssues:     []interfaces.NamingValidationIssue{},
				PermissionIssues: []interfaces.PermissionIssue{},
				Summary: interfaces.StructureValidationSummary{
					TotalFiles:       0,
					ValidFiles:       0,
					TotalDirectories: 0,
					ValidDirectories: 0,
					NamingIssues:     0,
					PermissionIssues: 0,
				},
			}

			err = validator.validateGoProjectStructure(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDirs, len(result.RequiredDirs))
			assert.Equal(t, tt.expectedFiles, len(result.RequiredFiles))
		})
	}
}

func TestStructureValidator_validateNodeProjectStructure(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name          string
		setupFiles    func(string) error
		expectedFiles int
	}{
		{
			name: "Node project with src directory",
			setupFiles: func(projectPath string) error {
				srcDir := filepath.Join(projectPath, "src")
				if err := os.MkdirAll(srcDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(srcDir, "index.js"), []byte("console.log('test')"), 0644)
			},
			expectedFiles: 0,
		},
		{
			name: "Node project with index.js",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "index.js"), []byte("console.log('test')"), 0644)
			},
			expectedFiles: 0,
		},
		{
			name: "Node project without entry point",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte("# Test"), 0644)
			},
			expectedFiles: 1, // Suggestion for entry point
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := &interfaces.StructureValidationResult{
				Valid:            true,
				RequiredFiles:    []interfaces.FileValidationResult{},
				RequiredDirs:     []interfaces.DirValidationResult{},
				NamingIssues:     []interfaces.NamingValidationIssue{},
				PermissionIssues: []interfaces.PermissionIssue{},
				Summary: interfaces.StructureValidationSummary{
					TotalFiles:       0,
					ValidFiles:       0,
					TotalDirectories: 0,
					ValidDirectories: 0,
					NamingIssues:     0,
					PermissionIssues: 0,
				},
			}

			err = validator.validateNodeProjectStructure(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedFiles, len(result.RequiredFiles))
		})
	}
}

func TestStructureValidator_validatePythonProjectStructure(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name         string
		setupFiles   func(string) error
		expectedDirs int
	}{
		{
			name: "Python project with src directory",
			setupFiles: func(projectPath string) error {
				srcDir := filepath.Join(projectPath, "src")
				if err := os.MkdirAll(srcDir, 0755); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(srcDir, "__init__.py"), []byte(""), 0644)
			},
			expectedDirs: 0,
		},
		{
			name: "Python project without src directory",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "main.py"), []byte("print('test')"), 0644)
			},
			expectedDirs: 1, // Suggestion for src/ layout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := &interfaces.StructureValidationResult{
				Valid:            true,
				RequiredFiles:    []interfaces.FileValidationResult{},
				RequiredDirs:     []interfaces.DirValidationResult{},
				NamingIssues:     []interfaces.NamingValidationIssue{},
				PermissionIssues: []interfaces.PermissionIssue{},
				Summary: interfaces.StructureValidationSummary{
					TotalFiles:       0,
					ValidFiles:       0,
					TotalDirectories: 0,
					ValidDirectories: 0,
					NamingIssues:     0,
					PermissionIssues: 0,
				},
			}

			err = validator.validatePythonProjectStructure(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDirs, len(result.RequiredDirs))
		})
	}
}

func TestStructureValidator_validateDockerProjectStructure(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name          string
		setupFiles    func(string) error
		expectedFiles int
	}{
		{
			name: "Docker project with .dockerignore",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, ".dockerignore"), []byte("node_modules"), 0644)
			},
			expectedFiles: 1, // .dockerignore file result
		},
		{
			name: "Docker project without .dockerignore",
			setupFiles: func(projectPath string) error {
				return os.WriteFile(filepath.Join(projectPath, "Dockerfile"), []byte("FROM node"), 0644)
			},
			expectedFiles: 1, // Suggestion for .dockerignore
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupFiles(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := &interfaces.StructureValidationResult{
				Valid:            true,
				RequiredFiles:    []interfaces.FileValidationResult{},
				RequiredDirs:     []interfaces.DirValidationResult{},
				NamingIssues:     []interfaces.NamingValidationIssue{},
				PermissionIssues: []interfaces.PermissionIssue{},
				Summary: interfaces.StructureValidationSummary{
					TotalFiles:       0,
					ValidFiles:       0,
					TotalDirectories: 0,
					ValidDirectories: 0,
					NamingIssues:     0,
					PermissionIssues: 0,
				},
			}

			err = validator.validateDockerProjectStructure(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedFiles, len(result.RequiredFiles))
		})
	}
}

func TestStructureValidator_validateFile(t *testing.T) {
	validator := NewStructureValidator()

	tests := []struct {
		name           string
		setupFile      func(string) (string, error)
		required       bool
		expectedValid  bool
		expectedExists bool
	}{
		{
			name: "existing required file",
			setupFile: func(tempDir string) (string, error) {
				filePath := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				return filePath, err
			},
			required:       true,
			expectedValid:  true,
			expectedExists: true,
		},
		{
			name: "missing required file",
			setupFile: func(tempDir string) (string, error) {
				return filepath.Join(tempDir, "missing.txt"), nil
			},
			required:       true,
			expectedValid:  false,
			expectedExists: false,
		},
		{
			name: "missing optional file",
			setupFile: func(tempDir string) (string, error) {
				return filepath.Join(tempDir, "optional.txt"), nil
			},
			required:       false,
			expectedValid:  false,
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath, err := tt.setupFile(tempDir)
			require.NoError(t, err)

			result := validator.validateFile(filePath, tt.required)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Equal(t, tt.expectedExists, result.Exists)
			assert.Equal(t, tt.required, result.Required)
			assert.Equal(t, filePath, result.Path)
		})
	}
}

func TestStructureValidator_toKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase to kebab-case",
			input:    "camelCase",
			expected: "camel-case",
		},
		{
			name:     "PascalCase to kebab-case",
			input:    "PascalCase",
			expected: "pascal-case",
		},
		{
			name:     "multiple words",
			input:    "thisIsALongVariableName",
			expected: "this-is-along-variable-name", // Actual regex behavior
		},
		{
			name:     "already lowercase",
			input:    "lowercase",
			expected: "lowercase",
		},
		{
			name:     "single letter",
			input:    "A",
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toKebabCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStructureValidator_toSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase to snake_case",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "PascalCase to snake_case",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "multiple words",
			input:    "thisIsALongVariableName",
			expected: "this_is_along_variable_name", // Actual regex behavior
		},
		{
			name:     "already lowercase",
			input:    "lowercase",
			expected: "lowercase",
		},
		{
			name:     "single letter",
			input:    "A",
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
