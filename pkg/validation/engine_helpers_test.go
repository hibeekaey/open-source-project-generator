package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine_validateProjectStructureBasic(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedIssues int
	}{
		{
			name: "project with all required files",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"README.md": "# Test Project",
					"LICENSE":   "MIT License",
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
			expectedIssues: 0,
		},
		{
			name: "project missing README",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"LICENSE": "MIT License",
				}
				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
		{
			name: "project missing LICENSE",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"README.md": "# Test Project",
				}
				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
		{
			name: "project missing both required files",
			setupProject: func(projectPath string) error {
				// Create empty project directory
				return nil
			},
			expectedValid:  false,
			expectedIssues: 2,
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

			result := &models.ValidationResult{
				Valid:   true,
				Issues:  []models.ValidationIssue{},
				Summary: "Test validation",
			}

			err = engine.validateProjectStructureBasic(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Issues, tt.expectedIssues)
		})
	}
}

func TestEngine_validateProjectDependenciesBasic(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedIssues int
	}{
		{
			name: "project with valid package.json",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"name": "test-project",
					"version": "1.0.0",
					"description": "Test project"
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
		{
			name: "project with invalid package.json",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"description": "Missing required fields"
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
		{
			name: "project with valid go.mod",
			setupProject: func(projectPath string) error {
				goMod := `module test-project

go 1.21

require (
	github.com/stretchr/testify v1.8.0
)`
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
		{
			name: "project with invalid go.mod",
			setupProject: func(projectPath string) error {
				goMod := `go 1.21` // Missing module declaration
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
		{
			name: "project with both valid files",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"name": "test-project",
					"version": "1.0.0"
				}`
				goMod := `module test-project

go 1.21`
				if err := os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedValid:  true,
			expectedIssues: 0,
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

			result := &models.ValidationResult{
				Valid:   true,
				Issues:  []models.ValidationIssue{},
				Summary: "Test validation",
			}

			err = engine.validateProjectDependenciesBasic(projectPath, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Issues, tt.expectedIssues)
		})
	}
}

func TestEngine_validateRequiredConfigFields(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name           string
		config         *models.ProjectConfig
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "config with all required fields",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				OutputPath:   "/tmp/test",
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "config missing name",
			config: &models.ProjectConfig{
				Organization: "test-org",
				OutputPath:   "/tmp/test",
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "config missing organization",
			config: &models.ProjectConfig{
				Name:       "test-project",
				OutputPath: "/tmp/test",
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "config missing both required fields",
			config: &models.ProjectConfig{
				OutputPath: "/tmp/test",
			},
			expectedValid:  false,
			expectedErrors: 2,
		},
		{
			name: "config with empty strings",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "",
				OutputPath:   "/tmp/test",
			},
			expectedValid:  false,
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			engine.validateRequiredConfigFields(tt.config, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Equal(t, tt.expectedErrors, len(result.Errors))
		})
	}
}

func TestEngine_validateConfigFieldFormats(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name             string
		config           *models.ProjectConfig
		expectedWarnings int
	}{
		{
			name: "config with valid formats",
			config: &models.ProjectConfig{
				Name:       "valid-project-name",
				OutputPath: "/valid/path",
			},
			expectedWarnings: 0,
		},
		{
			name: "config with invalid name format",
			config: &models.ProjectConfig{
				Name:       "Invalid Project Name!",
				OutputPath: "/valid/path",
			},
			expectedWarnings: 1,
		},
		{
			name: "config with dangerous output path",
			config: &models.ProjectConfig{
				Name:       "valid-name",
				OutputPath: "/path/with/../traversal",
			},
			expectedWarnings: 0, // Path traversal creates an error, not a warning
		},
		{
			name: "config with multiple format issues",
			config: &models.ProjectConfig{
				Name:       "Invalid Name!",
				OutputPath: "/path/../dangerous",
			},
			expectedWarnings: 1, // Only the name format creates a warning
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			engine.validateConfigFieldFormats(tt.config, result)

			assert.Equal(t, tt.expectedWarnings, len(result.Warnings))
		})
	}
}

func TestEngine_validateComponentConfiguration(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name             string
		config           *models.ProjectConfig
		expectedWarnings int
	}{
		{
			name: "config with frontend component enabled",
			config: &models.ProjectConfig{
				Components: models.Components{
					Frontend: models.FrontendComponents{
						NextJS: models.NextJSComponents{
							App: true,
						},
					},
				},
			},
			expectedWarnings: 0,
		},
		{
			name: "config with backend component enabled",
			config: &models.ProjectConfig{
				Components: models.Components{
					Backend: models.BackendComponents{
						GoGin: true,
					},
				},
			},
			expectedWarnings: 0,
		},
		{
			name: "config with no components enabled",
			config: &models.ProjectConfig{
				Components: models.Components{
					Frontend: models.FrontendComponents{
						NextJS: models.NextJSComponents{
							App: false,
						},
					},
					Backend: models.BackendComponents{
						GoGin: false,
					},
				},
			},
			expectedWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			engine.validateComponentConfiguration(tt.config, result)

			assert.Equal(t, tt.expectedWarnings, len(result.Warnings))
		})
	}
}
