package components

import (
	"os"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MockFileSystemOperations implements FileSystemOperations for testing
type MockFileSystemOperations struct {
	files map[string][]byte
}

func NewMockFileSystemOperations() *MockFileSystemOperations {
	return &MockFileSystemOperations{
		files: make(map[string][]byte),
	}
}

func (m *MockFileSystemOperations) WriteFile(path string, content []byte, perm os.FileMode) error {
	m.files[path] = content
	return nil
}

func (m *MockFileSystemOperations) EnsureDirectory(path string) error {
	return nil
}

func (m *MockFileSystemOperations) FileExists(path string) bool {
	_, exists := m.files[path]
	return exists
}

func (m *MockFileSystemOperations) GetFileContent(path string) []byte {
	return m.files[path]
}

func TestFrontendGenerator_GenerateFiles(t *testing.T) {
	tests := []struct {
		name          string
		config        *models.ProjectConfig
		expectedFiles []string
		expectedError bool
	}{
		{
			name: "Generate all frontend components",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Frontend: models.FrontendComponents{
						NextJS: models.NextJSComponents{
							App:    true,
							Home:   true,
							Admin:  true,
							Shared: true,
						},
					},
				},
				Versions: &models.VersionConfig{
					Packages: map[string]string{
						"next":  "14.0.0",
						"react": "18.0.0",
					},
				},
			},
			expectedFiles: []string{
				"testproject/App/package.json",
				"testproject/App/next.config.js",
				"testproject/App/tailwind.config.js",
				"testproject/App/tsconfig.json",
				"testproject/App/.eslintrc.json",
				"testproject/Home/package.json",
				"testproject/Home/next.config.js",
				"testproject/Admin/package.json",
				"testproject/Admin/next.config.js",
				"testproject/Shared/package.json",
			},
			expectedError: false,
		},
		{
			name: "Generate only app component",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Frontend: models.FrontendComponents{
						NextJS: models.NextJSComponents{
							App: true,
						},
					},
				},
			},
			expectedFiles: []string{
				"testproject/App/package.json",
				"testproject/App/next.config.js",
				"testproject/App/tailwind.config.js",
				"testproject/App/tsconfig.json",
				"testproject/App/.eslintrc.json",
			},
			expectedError: false,
		},
		{
			name:          "Nil config should return error",
			config:        nil,
			expectedFiles: []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := NewMockFileSystemOperations()
			fg := NewFrontendGenerator(mockFS)

			err := fg.GenerateFiles("testproject", tt.config)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that expected files were created
			for _, expectedFile := range tt.expectedFiles {
				if !mockFS.FileExists(expectedFile) {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}

			// Verify file contents for package.json files
			if tt.config != nil && tt.config.Components.Frontend.NextJS.App {
				appPackageJsonPath := "testproject/App/package.json"
				content := string(mockFS.GetFileContent(appPackageJsonPath))

				if !strings.Contains(content, tt.config.Name+"-app") {
					t.Errorf("App package.json should contain project name")
				}

				if tt.config.Versions != nil && tt.config.Versions.Packages != nil {
					if nextVersion, ok := tt.config.Versions.Packages["next"]; ok {
						if !strings.Contains(content, nextVersion) {
							t.Errorf("App package.json should contain Next.js version %s", nextVersion)
						}
					}
				}
			}
		})
	}
}

func TestFrontendGenerator_generateAppPackageJson(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.0.0",
			},
		},
	}

	content := fg.generateAppPackageJson(config)

	// Check that content contains expected elements
	expectedElements := []string{
		`"name": "testapp-app"`,
		`"next": "14.0.0"`,
		`"react": "18.0.0"`,
		`"dev": "next dev"`,
		`"build": "next build"`,
		`"@tailwindcss/forms"`,
		`"typescript"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Package.json should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateNextConfig(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	content := fg.generateNextConfig()

	expectedElements := []string{
		"nextConfig",
		"experimental",
		"appDir: true",
		"images",
		"domains",
		"module.exports = nextConfig",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Next.js config should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateTailwindConfig(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	content := fg.generateTailwindConfig()

	expectedElements := []string{
		"module.exports",
		"content",
		"./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/components/**/*.{js,ts,jsx,tsx,mdx}",
		"theme",
		"extend",
		"colors",
		"primary",
		"plugins",
		"@tailwindcss/forms",
		"@tailwindcss/typography",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Tailwind config should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateTSConfig(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	content := fg.generateTSConfig()

	expectedElements := []string{
		`"compilerOptions"`,
		`"target": "es5"`,
		`"strict": true`,
		`"jsx": "preserve"`,
		`"moduleResolution": "bundler"`,
		`"@/*": ["./src/*"]`,
		`"include"`,
		`"exclude"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("TypeScript config should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateESLintConfig(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	content := fg.generateESLintConfig()

	expectedElements := []string{
		`"extends"`,
		`"next/core-web-vitals"`,
		`"rules"`,
		`"prefer-const"`,
		`"no-unused-vars"`,
		`"no-console"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("ESLint config should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateHomePackageJson(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.0.0",
			},
		},
	}

	content := fg.generateHomePackageJson(config)

	expectedElements := []string{
		`"name": "testapp-home"`,
		`"dev": "next dev -p 3001"`,
		`"start": "next start -p 3001"`,
		`"next": "14.0.0"`,
		`"react": "18.0.0"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Home package.json should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateAdminPackageJson(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.0.0",
			},
		},
	}

	content := fg.generateAdminPackageJson(config)

	expectedElements := []string{
		`"name": "testapp-admin"`,
		`"dev": "next dev -p 3002"`,
		`"start": "next start -p 3002"`,
		`"@headlessui/react"`,
		`"@heroicons/react"`,
		`"next": "14.0.0"`,
		`"react": "18.0.0"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Admin package.json should contain %s", element)
		}
	}
}

func TestFrontendGenerator_generateSharedPackageJson(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"react": "18.0.0",
			},
		},
	}

	content := fg.generateSharedPackageJson(config)

	expectedElements := []string{
		`"name": "testapp-shared"`,
		`"main": "dist/index.js"`,
		`"types": "dist/index.d.ts"`,
		`"peerDependencies"`,
		`"react": "18.0.0"`,
		`"typescript"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Shared package.json should contain %s", element)
		}
	}
}

func TestFrontendGenerator_WithoutVersions(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	fg := NewFrontendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		// No versions specified
	}

	content := fg.generateAppPackageJson(config)

	// Should use default versions
	expectedElements := []string{
		`"next": "14.0.0"`,
		`"react": "18.0.0"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Should use default versions when none specified: %s", element)
		}
	}
}
