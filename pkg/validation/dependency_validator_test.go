package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDependencyValidator(t *testing.T) {
	validator := NewDependencyValidator()

	assert.NotNil(t, validator)
	assert.NotNil(t, validator.vulnerabilityDB)
	assert.NotNil(t, validator.packageRegistry)

	// Check that some known vulnerabilities are loaded
	assert.Contains(t, validator.vulnerabilityDB, "lodash")
	assert.Contains(t, validator.vulnerabilityDB, "express")
}

func TestDependencyValidator_ValidateProjectDependencies(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name          string
		setupProject  func(string) error
		expectedValid bool
		expectedDeps  int
		expectedVulns int
	}{
		{
			name: "project with valid package.json",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"name": "test-project",
					"version": "1.0.0",
					"dependencies": {
						"react": "^18.0.0",
						"express": "^4.18.0"
					},
					"devDependencies": {
						"jest": "^28.0.0"
					}
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedValid: true,
			expectedDeps:  3,
			expectedVulns: 0,
		},
		{
			name: "project with vulnerable dependencies",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"name": "test-project",
					"version": "1.0.0",
					"dependencies": {
						"lodash": "4.17.15",
						"express": "4.16.0"
					}
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedValid: true,
			expectedDeps:  2,
			expectedVulns: 2,
		},
		{
			name: "project with valid go.mod",
			setupProject: func(projectPath string) error {
				goMod := `module test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.0 // indirect
)`
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedValid: true,
			expectedDeps:  2,
			expectedVulns: 0,
		},
		{
			name: "project with requirements.txt",
			setupProject: func(projectPath string) error {
				requirements := `django==4.2.0
requests>=2.28.0
pytest==7.1.0`
				return os.WriteFile(filepath.Join(projectPath, "requirements.txt"), []byte(requirements), 0644)
			},
			expectedValid: true,
			expectedDeps:  3,
			expectedVulns: 0,
		},
		{
			name: "empty project",
			setupProject: func(projectPath string) error {
				return nil // No dependency files
			},
			expectedValid: true,
			expectedDeps:  0,
			expectedVulns: 0,
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

			result, err := validator.ValidateProjectDependencies(projectPath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Equal(t, tt.expectedDeps, len(result.Dependencies))
			assert.Equal(t, tt.expectedVulns, len(result.Vulnerabilities))
		})
	}
}

func TestDependencyValidator_validatePackageJSONDependencies(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name           string
		packageJSON    string
		expectedValid  bool
		expectedDeps   int
		expectedErrors int
	}{
		{
			name: "valid package.json with dependencies",
			packageJSON: `{
				"name": "test-project",
				"version": "1.0.0",
				"dependencies": {
					"react": "^18.0.0",
					"lodash": "^4.17.21"
				},
				"devDependencies": {
					"jest": "^28.0.0"
				}
			}`,
			expectedValid: true,
			expectedDeps:  3,
		},
		{
			name: "package.json with invalid JSON",
			packageJSON: `{
				"name": "test-project"
				"version": "1.0.0"
			}`,
			expectedValid:  false,
			expectedDeps:   0,
			expectedErrors: 1,
		},
		{
			name: "package.json missing required fields",
			packageJSON: `{
				"description": "Test project"
			}`,
			expectedValid: false,
			expectedDeps:  0,
		},
		{
			name: "package.json with invalid name",
			packageJSON: `{
				"name": "Invalid Package Name",
				"version": "1.0.0"
			}`,
			expectedValid: false,
			expectedDeps:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			packageJSONPath := filepath.Join(tempDir, "package.json")
			err := os.WriteFile(packageJSONPath, []byte(tt.packageJSON), 0644)
			require.NoError(t, err)

			result := &interfaces.DependencyValidationResult{
				Valid:           true,
				Dependencies:    []interfaces.DependencyValidation{},
				Vulnerabilities: []interfaces.DependencyVulnerability{},
				Outdated:        []interfaces.OutdatedDependency{},
				Conflicts:       []interfaces.DependencyConflict{},
				Summary: interfaces.DependencyValidationSummary{
					TotalDependencies: 0,
					ValidDependencies: 0,
					Vulnerabilities:   0,
					OutdatedCount:     0,
					ConflictCount:     0,
				},
			}

			err = validator.validatePackageJSONDependencies(tempDir, result)

			if tt.expectedErrors > 0 {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValid, result.Valid)
				assert.Equal(t, tt.expectedDeps, len(result.Dependencies))
			}
		})
	}
}

func TestDependencyValidator_validateGoModDependencies(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name         string
		goMod        string
		expectedDeps int
	}{
		{
			name: "valid go.mod with dependencies",
			goMod: `module test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.0 // indirect
)`,
			expectedDeps: 2,
		},
		{
			name: "go.mod with single line requires",
			goMod: `module test-project

go 1.21

require github.com/gin-gonic/gin v1.9.1
require github.com/stretchr/testify v1.8.0`,
			expectedDeps: 2,
		},
		{
			name: "minimal go.mod",
			goMod: `module test-project

go 1.21`,
			expectedDeps: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			goModPath := filepath.Join(tempDir, "go.mod")
			err := os.WriteFile(goModPath, []byte(tt.goMod), 0644)
			require.NoError(t, err)

			result := &interfaces.DependencyValidationResult{
				Valid:           true,
				Dependencies:    []interfaces.DependencyValidation{},
				Vulnerabilities: []interfaces.DependencyVulnerability{},
				Outdated:        []interfaces.OutdatedDependency{},
				Conflicts:       []interfaces.DependencyConflict{},
				Summary: interfaces.DependencyValidationSummary{
					TotalDependencies: 0,
					ValidDependencies: 0,
					Vulnerabilities:   0,
					OutdatedCount:     0,
					ConflictCount:     0,
				},
			}

			err = validator.validateGoModDependencies(tempDir, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDeps, len(result.Dependencies))
		})
	}
}

func TestDependencyValidator_validateRequirementsTxtDependencies(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name         string
		requirements string
		expectedDeps int
	}{
		{
			name: "valid requirements.txt",
			requirements: `django==4.2.0
requests>=2.28.0
pytest==7.1.0
# This is a comment
numpy~=1.24.0`,
			expectedDeps: 4,
		},
		{
			name: "requirements with various operators",
			requirements: `package1==1.0.0
package2>=2.0.0
package3<=3.0.0
package4>4.0.0
package5<5.0.0
package6~=6.0.0
package7!=7.0.0`,
			expectedDeps: 7,
		},
		{
			name: "requirements with package names only",
			requirements: `django
requests
pytest`,
			expectedDeps: 3,
		},
		{
			name: "empty requirements.txt",
			requirements: `# Just comments

# More comments`,
			expectedDeps: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			requirementsPath := filepath.Join(tempDir, "requirements.txt")
			err := os.WriteFile(requirementsPath, []byte(tt.requirements), 0644)
			require.NoError(t, err)

			result := &interfaces.DependencyValidationResult{
				Valid:           true,
				Dependencies:    []interfaces.DependencyValidation{},
				Vulnerabilities: []interfaces.DependencyVulnerability{},
				Outdated:        []interfaces.OutdatedDependency{},
				Conflicts:       []interfaces.DependencyConflict{},
				Summary: interfaces.DependencyValidationSummary{
					TotalDependencies: 0,
					ValidDependencies: 0,
					Vulnerabilities:   0,
					OutdatedCount:     0,
					ConflictCount:     0,
				},
			}

			err = validator.validateRequirementsTxtDependencies(tempDir, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDeps, len(result.Dependencies))
		})
	}
}

func TestDependencyValidator_validateGoVersion(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "valid current version",
			version:     "1.21",
			expectError: false,
		},
		{
			name:        "valid version with patch",
			version:     "1.21.5",
			expectError: false,
		},
		{
			name:        "valid newer version",
			version:     "1.22",
			expectError: false,
		},
		{
			name:        "old but acceptable version",
			version:     "1.19",
			expectError: false,
		},
		{
			name:        "very old version",
			version:     "1.16",
			expectError: true,
		},
		{
			name:        "invalid version format",
			version:     "1",
			expectError: true,
		},
		{
			name:        "invalid version format with letters",
			version:     "1.21a",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.DependencyValidationResult{
				Valid: true,
				Summary: interfaces.DependencyValidationSummary{
					TotalDependencies: 0,
					ValidDependencies: 0,
					Vulnerabilities:   0,
					OutdatedCount:     0,
					ConflictCount:     0,
				},
			}

			err := validator.validateGoVersion(tt.version, result)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDependencyValidator_validateNpmPackageName(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name        string
		packageName string
		expectError bool
	}{
		{
			name:        "valid package name",
			packageName: "my-package",
			expectError: false,
		},
		{
			name:        "valid scoped package",
			packageName: "@scope/package-name",
			expectError: false,
		},
		{
			name:        "valid package with dots",
			packageName: "package.name",
			expectError: false,
		},
		{
			name:        "invalid uppercase",
			packageName: "MyPackage",
			expectError: true,
		},
		{
			name:        "invalid spaces",
			packageName: "my package",
			expectError: true,
		},
		{
			name:        "invalid special characters",
			packageName: "package@name",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateNpmPackageName(tt.packageName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDependencyValidator_validateNpmVersionFormat(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "exact version",
			version:     "1.0.0",
			expectError: false,
		},
		{
			name:        "caret range",
			version:     "^1.0.0",
			expectError: false,
		},
		{
			name:        "tilde range",
			version:     "~1.0.0",
			expectError: false,
		},
		{
			name:        "greater than",
			version:     ">=1.0.0",
			expectError: false,
		},
		{
			name:        "less than",
			version:     "<=2.0.0",
			expectError: false,
		},
		{
			name:        "invalid version",
			version:     "invalid",
			expectError: true,
		},
		{
			name:        "incomplete version",
			version:     "1.0",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateNpmVersionFormat(tt.version)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDependencyValidator_validateSemanticVersion(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "valid semantic version",
			version:     "1.0.0",
			expectError: false,
		},
		{
			name:        "valid with pre-release",
			version:     "1.0.0-alpha.1",
			expectError: false,
		},
		{
			name:        "valid with build metadata",
			version:     "1.0.0+build.1",
			expectError: false,
		},
		{
			name:        "valid with v prefix",
			version:     "v1.0.0",
			expectError: false,
		},
		{
			name:        "invalid format",
			version:     "1.0",
			expectError: true,
		},
		{
			name:        "invalid characters",
			version:     "1.0.0a",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateSemanticVersion(tt.version)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDependencyValidator_validatePythonPackageName(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name        string
		packageName string
		expectError bool
	}{
		{
			name:        "valid package name",
			packageName: "django",
			expectError: false,
		},
		{
			name:        "valid with hyphens",
			packageName: "django-rest-framework",
			expectError: false,
		},
		{
			name:        "valid with underscores",
			packageName: "django_extensions",
			expectError: false,
		},
		{
			name:        "valid with dots",
			packageName: "zope.interface",
			expectError: false,
		},
		{
			name:        "valid with numbers",
			packageName: "python3-dev",
			expectError: false,
		},
		{
			name:        "invalid special characters",
			packageName: "package@name",
			expectError: true,
		},
		{
			name:        "invalid spaces",
			packageName: "package name",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validatePythonPackageName(tt.packageName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDependencyValidator_checkDependencyConflicts(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name              string
		dependencies      []interfaces.DependencyValidation
		expectedConflicts int
	}{
		{
			name: "no conflicts",
			dependencies: []interfaces.DependencyValidation{
				{Name: "react", Version: "18.0.0", Type: "production"},
				{Name: "lodash", Version: "4.17.21", Type: "production"},
			},
			expectedConflicts: 0,
		},
		{
			name: "version conflict",
			dependencies: []interfaces.DependencyValidation{
				{Name: "react", Version: "18.0.0", Type: "production"},
				{Name: "react", Version: "17.0.0", Type: "development"},
			},
			expectedConflicts: 1,
		},
		{
			name: "multiple conflicts",
			dependencies: []interfaces.DependencyValidation{
				{Name: "lodash", Version: "4.17.21", Type: "production"},
				{Name: "lodash", Version: "4.17.15", Type: "development"},
				{Name: "express", Version: "4.18.0", Type: "production"},
				{Name: "express", Version: "4.17.0", Type: "peer"},
			},
			expectedConflicts: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.DependencyValidationResult{
				Valid:        true,
				Dependencies: tt.dependencies,
				Conflicts:    []interfaces.DependencyConflict{},
				Summary: interfaces.DependencyValidationSummary{
					ConflictCount: 0,
				},
			}

			err := validator.checkDependencyConflicts(result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedConflicts, len(result.Conflicts))
			assert.Equal(t, tt.expectedConflicts, result.Summary.ConflictCount)
		})
	}
}

func TestDependencyValidator_checkOutdatedDependencies(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name             string
		dependencies     []interfaces.DependencyValidation
		expectedOutdated int
	}{
		{
			name: "up to date dependencies",
			dependencies: []interfaces.DependencyValidation{
				{Name: "lodash", Version: "4.17.21", Type: "production"},
				{Name: "express", Version: "4.18.2", Type: "production"},
			},
			expectedOutdated: 0,
		},
		{
			name: "outdated dependencies",
			dependencies: []interfaces.DependencyValidation{
				{Name: "lodash", Version: "4.17.15", Type: "production"},
				{Name: "express", Version: "4.16.0", Type: "production"},
			},
			expectedOutdated: 2,
		},
		{
			name: "unknown dependencies",
			dependencies: []interfaces.DependencyValidation{
				{Name: "unknown-package", Version: "1.0.0", Type: "production"},
			},
			expectedOutdated: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.DependencyValidationResult{
				Valid:        true,
				Dependencies: tt.dependencies,
				Outdated:     []interfaces.OutdatedDependency{},
				Summary: interfaces.DependencyValidationSummary{
					OutdatedCount: 0,
				},
			}

			err := validator.checkOutdatedDependencies(result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedOutdated, len(result.Outdated))
			assert.Equal(t, tt.expectedOutdated, result.Summary.OutdatedCount)
		})
	}
}

func TestDependencyValidator_versionAffectedByVulnerability(t *testing.T) {
	validator := NewDependencyValidator()

	vuln := interfaces.DependencyVulnerability{
		Dependency:  "lodash",
		Version:     "4.17.15",
		CVEID:       "CVE-2020-8203",
		Severity:    "high",
		Description: "Prototype pollution vulnerability",
		FixedIn:     "4.17.19",
		CVSS:        7.4,
	}

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "exact match",
			version:  "4.17.15",
			expected: true,
		},
		{
			name:     "different version",
			version:  "4.17.21",
			expected: false,
		},
		{
			name:     "fixed version",
			version:  "4.17.19",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.versionAffectedByVulnerability(tt.version, vuln)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDependencyValidator_isVersionOlder(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name     string
		version1 string
		version2 string
		expected bool
	}{
		{
			name:     "older version",
			version1: "1.0.0",
			version2: "2.0.0",
			expected: true,
		},
		{
			name:     "newer version",
			version1: "2.0.0",
			version2: "1.0.0",
			expected: false,
		},
		{
			name:     "same version",
			version1: "1.0.0",
			version2: "1.0.0",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isVersionOlder(tt.version1, tt.version2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDependencyValidator_getUpdateType(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name           string
		currentVersion string
		latestVersion  string
		expected       string
	}{
		{
			name:           "minor update",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			expected:       "minor",
		},
		{
			name:           "major update",
			currentVersion: "1.0.0",
			latestVersion:  "2.0.0",
			expected:       "major",
		},
		{
			name:           "patch update",
			currentVersion: "1.0.0",
			latestVersion:  "1.0.1",
			expected:       "minor", // Simple implementation treats all as minor if major matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.getUpdateType(tt.currentVersion, tt.latestVersion)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDependencyValidator_isBreakingUpdate(t *testing.T) {
	validator := NewDependencyValidator()

	tests := []struct {
		name           string
		currentVersion string
		latestVersion  string
		expected       bool
	}{
		{
			name:           "major version change",
			currentVersion: "1.0.0",
			latestVersion:  "2.0.0",
			expected:       true,
		},
		{
			name:           "minor version change",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			expected:       false,
		},
		{
			name:           "patch version change",
			currentVersion: "1.0.0",
			latestVersion:  "1.0.1",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isBreakingUpdate(tt.currentVersion, tt.latestVersion)
			assert.Equal(t, tt.expected, result)
		})
	}
}
