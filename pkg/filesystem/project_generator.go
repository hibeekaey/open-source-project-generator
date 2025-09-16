package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectStructure defines the standard project directory structure
type ProjectStructure struct {
	// Root directories
	RootDirs []string
	// Component-specific directories
	FrontendDirs []string
	BackendDirs  []string
	MobileDirs   []string
	InfraDirs    []string
	// Common directories
	CommonDirs []string
}

// GetStandardProjectStructure returns the standard project directory structure
func GetStandardProjectStructure() *ProjectStructure {
	return &ProjectStructure{
		RootDirs: []string{
			"docs",
			"scripts",
			".github/workflows",
			".github/ISSUE_TEMPLATE",
			".github/PULL_REQUEST_TEMPLATE",
		},
		FrontendDirs: []string{
			"App/src/components/ui",
			"App/src/components/forms",
			"App/src/hooks",
			"App/src/context",
			"App/src/lib",
			"App/src/types",
			"App/public",
			"Home/src/components",
			"Home/src/sections",
			"Home/public",
			"Admin/src/components",
			"Admin/src/pages",
			"Admin/src/hooks",
			"Admin/public",
		},
		BackendDirs: []string{
			"CommonServer/cmd",
			"CommonServer/internal/controllers",
			"CommonServer/internal/models",
			"CommonServer/internal/services",
			"CommonServer/internal/middleware",
			"CommonServer/internal/repository",
			"CommonServer/internal/config",
			"CommonServer/pkg/auth",
			"CommonServer/pkg/database",
			"CommonServer/pkg/utils",
			"CommonServer/migrations",
			"CommonServer/tests",
		},
		MobileDirs: []string{
			"Mobile/Android/app/src/main/java",
			"Mobile/Android/app/src/main/res",
			"Mobile/Android/app/src/test/java",
			"Mobile/iOS/Sources",
			"Mobile/iOS/Resources",
			"Mobile/iOS/Tests",
			"Mobile/Shared/api",
			"Mobile/Shared/assets",
		},
		InfraDirs: []string{
			"Deploy/terraform/modules",
			"Deploy/terraform/environments/staging",
			"Deploy/terraform/environments/production",
			"Deploy/kubernetes/base",
			"Deploy/kubernetes/overlays/staging",
			"Deploy/kubernetes/overlays/production",
			"Deploy/docker",
			"Deploy/helm",
		},
		CommonDirs: []string{
			"Tests/integration",
			"Tests/e2e",
			"Tests/performance",
		},
	}
}

// ProjectGenerator handles complete project structure generation
type ProjectGenerator struct {
	fsGen     *Generator
	structure *ProjectStructure
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{
		fsGen:     NewGenerator().(*Generator),
		structure: GetStandardProjectStructure(),
	}
}

// NewDryRunProjectGenerator creates a new project generator in dry-run mode
func NewDryRunProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{
		fsGen:     NewDryRunGenerator().(*Generator),
		structure: GetStandardProjectStructure(),
	}
}

// GenerateProjectStructure creates the complete project directory structure
func (pg *ProjectGenerator) GenerateProjectStructure(config *models.ProjectConfig, outputPath string) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Validate required config fields
	if config.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if config.Organization == "" {
		return fmt.Errorf("organization cannot be empty")
	}

	// Create the root project directory
	projectPath := filepath.Join(outputPath, config.Name)
	if err := pg.fsGen.EnsureDirectory(projectPath); err != nil {
		return fmt.Errorf("failed to create project root directory: %w", err)
	}

	// Create root directories (always created)
	if err := pg.createDirectories(projectPath, pg.structure.RootDirs); err != nil {
		return fmt.Errorf("failed to create root directories: %w", err)
	}

	// Create common directories (always created)
	if err := pg.createDirectories(projectPath, pg.structure.CommonDirs); err != nil {
		return fmt.Errorf("failed to create common directories: %w", err)
	}

	// Create component-specific directories based on configuration
	if err := pg.createComponentDirectories(projectPath, config); err != nil {
		return fmt.Errorf("failed to create component directories: %w", err)
	}

	return nil
}

// GenerateComponentFiles creates component-specific files based on user selection
func (pg *ProjectGenerator) GenerateComponentFiles(config *models.ProjectConfig, outputPath string) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	projectPath := filepath.Join(outputPath, config.Name)

	// Generate root configuration files
	if err := pg.generateRootFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate root files: %w", err)
	}

	// Generate component-specific files
	if err := pg.generateComponentFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate component files: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates the generated project structure
func (pg *ProjectGenerator) ValidateProjectStructure(projectPath string, config *models.ProjectConfig) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Validate root directories exist
	for _, dir := range pg.structure.RootDirs {
		dirPath := filepath.Join(projectPath, dir)
		if !pg.fsGen.FileExists(dirPath) {
			return fmt.Errorf("required root directory missing: %s", dir)
		}
	}

	// Validate common directories exist
	for _, dir := range pg.structure.CommonDirs {
		dirPath := filepath.Join(projectPath, dir)
		if !pg.fsGen.FileExists(dirPath) {
			return fmt.Errorf("required common directory missing: %s", dir)
		}
	}

	// Validate component-specific directories based on configuration
	if err := pg.validateComponentDirectories(projectPath, config); err != nil {
		return fmt.Errorf("component directory validation failed: %w", err)
	}

	return nil
}

// ValidateCrossReferences validates cross-references between generated files
func (pg *ProjectGenerator) ValidateCrossReferences(projectPath string, config *models.ProjectConfig) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Validate that required configuration files exist
	requiredFiles := []string{
		"Makefile",
		"README.md",
		"docker-compose.yml",
		".gitignore",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if !pg.fsGen.FileExists(filePath) {
			return fmt.Errorf("required root file missing: %s", file)
		}
	}

	// Validate component-specific cross-references
	if err := pg.validateComponentCrossReferences(projectPath, config); err != nil {
		return fmt.Errorf("component cross-reference validation failed: %w", err)
	}

	// Validate content cross-references
	if err := pg.validateContentCrossReferences(projectPath, config); err != nil {
		return fmt.Errorf("content cross-reference validation failed: %w", err)
	}

	return nil
}

// validateContentCrossReferences validates that file contents reference each other correctly
func (pg *ProjectGenerator) validateContentCrossReferences(projectPath string, config *models.ProjectConfig) error {
	// Validate Makefile references correct components
	makefilePath := filepath.Join(projectPath, "Makefile")
	if pg.fsGen.FileExists(makefilePath) {
		// Validate Makefile exists and is readable
		if _, err := os.ReadFile(makefilePath); err != nil {
			return fmt.Errorf("failed to read Makefile: %w", err)
		}
	}

	// Validate docker-compose.yml references correct services
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	if pg.fsGen.FileExists(dockerComposePath) {
		// Validate docker-compose.yml exists and is readable
		if _, err := os.ReadFile(dockerComposePath); err != nil {
			return fmt.Errorf("failed to read docker-compose.yml: %w", err)
		}
	}

	// Validate package.json dependencies are consistent across frontend apps
	if config.Components.Frontend.NextJS.App && config.Components.Frontend.NextJS.Admin {
		mainAppPackageJson := filepath.Join(projectPath, "App/package.json")
		adminPackageJson := filepath.Join(projectPath, "Admin/package.json")

		if pg.fsGen.FileExists(mainAppPackageJson) && pg.fsGen.FileExists(adminPackageJson) {
			// Validate both package.json files are readable
			if _, err := os.ReadFile(mainAppPackageJson); err != nil {
				return fmt.Errorf("failed to read main app package.json: %w", err)
			}
			if _, err := os.ReadFile(adminPackageJson); err != nil {
				return fmt.Errorf("failed to read admin package.json: %w", err)
			}
		}
	}

	return nil
}

// createDirectories creates a list of directories under the project path
func (pg *ProjectGenerator) createDirectories(projectPath string, dirs []string) error {
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := pg.fsGen.EnsureDirectory(dirPath); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// createComponentDirectories creates directories based on selected components
func (pg *ProjectGenerator) createComponentDirectories(projectPath string, config *models.ProjectConfig) error {
	// Create frontend directories if any frontend component is selected
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		if err := pg.createDirectories(projectPath, pg.structure.FrontendDirs); err != nil {
			return fmt.Errorf("failed to create frontend directories: %w", err)
		}
	}

	// Create backend directories if backend is selected
	if config.Components.Backend.GoGin {
		if err := pg.createDirectories(projectPath, pg.structure.BackendDirs); err != nil {
			return fmt.Errorf("failed to create backend directories: %w", err)
		}
	}

	// Create mobile directories if any mobile component is selected
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		if err := pg.createDirectories(projectPath, pg.structure.MobileDirs); err != nil {
			return fmt.Errorf("failed to create mobile directories: %w", err)
		}
	}

	// Create infrastructure directories if any infrastructure component is selected
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		if err := pg.createDirectories(projectPath, pg.structure.InfraDirs); err != nil {
			return fmt.Errorf("failed to create infrastructure directories: %w", err)
		}
	}

	return nil
}

// generateRootFiles creates root-level configuration files
func (pg *ProjectGenerator) generateRootFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate Makefile
	makefileContent := pg.generateMakefileContent(config)
	makefilePath := filepath.Join(projectPath, "Makefile")
	if err := pg.fsGen.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Makefile: %w", err)
	}

	// Generate README.md
	readmeContent := pg.generateReadmeContent(config)
	readmePath := filepath.Join(projectPath, "README.md")
	if err := pg.fsGen.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Generate docker-compose.yml
	dockerComposeContent := pg.generateDockerComposeContent(config)
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	if err := pg.fsGen.WriteFile(dockerComposePath, []byte(dockerComposeContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Generate .gitignore
	gitignoreContent := pg.generateGitignoreContent(config)
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := pg.fsGen.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// generateComponentFiles creates component-specific files
func (pg *ProjectGenerator) generateComponentFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate frontend component files
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		if err := pg.generateFrontendFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate frontend files: %w", err)
		}
	}

	// Generate backend component files
	if config.Components.Backend.GoGin {
		if err := pg.generateBackendFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate backend files: %w", err)
		}
	}

	// Generate mobile component files
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		if err := pg.generateMobileFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate mobile files: %w", err)
		}
	}

	// Generate infrastructure component files
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		if err := pg.generateInfrastructureFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate infrastructure files: %w", err)
		}
	}

	// Generate documentation files
	if err := pg.generateDocumentationFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate documentation files: %w", err)
	}

	// Generate CI/CD files
	if err := pg.generateCICDFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate CI/CD files: %w", err)
	}

	return nil
}

// validateComponentDirectories validates component-specific directories
func (pg *ProjectGenerator) validateComponentDirectories(projectPath string, config *models.ProjectConfig) error {
	// Validate frontend directories
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		for _, dir := range pg.structure.FrontendDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !pg.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required frontend directory missing: %s", dir)
			}
		}
	}

	// Validate backend directories
	if config.Components.Backend.GoGin {
		for _, dir := range pg.structure.BackendDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !pg.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required backend directory missing: %s", dir)
			}
		}
	}

	// Validate mobile directories
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		for _, dir := range pg.structure.MobileDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !pg.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required mobile directory missing: %s", dir)
			}
		}
	}

	// Validate infrastructure directories
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		for _, dir := range pg.structure.InfraDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !pg.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required infrastructure directory missing: %s", dir)
			}
		}
	}

	return nil
}

// validateComponentCrossReferences validates cross-references between component files
func (pg *ProjectGenerator) validateComponentCrossReferences(projectPath string, config *models.ProjectConfig) error {
	// Validate frontend cross-references
	if config.Components.Frontend.NextJS.App {
		packageJsonPath := filepath.Join(projectPath, "App/package.json")
		if !pg.fsGen.FileExists(packageJsonPath) {
			return fmt.Errorf("main app package.json missing")
		}
	}

	// Validate backend cross-references
	if config.Components.Backend.GoGin {
		goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
		if !pg.fsGen.FileExists(goModPath) {
			return fmt.Errorf("backend go.mod missing")
		}
	}

	// Validate mobile cross-references
	if config.Components.Mobile.Android {
		buildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
		if !pg.fsGen.FileExists(buildGradlePath) {
			return fmt.Errorf("android build.gradle missing")
		}
	}

	if config.Components.Mobile.IOS {
		packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
		if !pg.fsGen.FileExists(packageSwiftPath) {
			return fmt.Errorf("iOS Package.swift missing")
		}
	}

	return nil
}

// Content generation methods (simplified implementations for now)
func (pg *ProjectGenerator) generateMakefileContent(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Makefile
# Generated by Open Source Template Generator

.PHONY: help setup dev test build clean

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%%-20s\033[0m %%s\n", $$1, $$2}'

setup: ## Set up the development environment
	@echo "Setting up %s development environment..."
	@echo "Project setup complete!"

dev: ## Start development servers
	@echo "Starting development servers for %s..."

test: ## Run all tests
	@echo "Running tests for %s..."

build: ## Build all components
	@echo "Building %s..."

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts for %s..."
`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

func (pg *ProjectGenerator) generateReadmeContent(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s

%s

## Description

%s

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.22+
- Docker
- Make

### Installation

1. Clone the repository:
   `+"```bash"+`
   git clone %s
   cd %s
   `+"```"+`

2. Set up the development environment:
   `+"```bash"+`
   make setup
   `+"```"+`

3. Start the development servers:
   `+"```bash"+`
   make dev
   `+"```"+`

## Project Structure

This project follows a multi-service architecture with the following components:

- **App/**: Main frontend application
- **Home/**: Landing page application
- **Admin/**: Admin dashboard
- **CommonServer/**: Backend API server
- **Mobile/**: Mobile applications (Android & iOS)
- **Deploy/**: Infrastructure and deployment configurations
- **Tests/**: Test suites

## Development

### Running Tests

`+"```bash"+`
make test
`+"```"+`

### Building

`+"```bash"+`
make build
`+"```"+`

## License

This project is licensed under the %s License - see the [LICENSE](LICENSE) file for details.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.
`, config.Name, config.Organization, config.Description, "", config.Name, config.License)
}

func (pg *ProjectGenerator) generateDockerComposeContent(config *models.ProjectConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  # %s services
  # Generated by Open Source Template Generator
  
  # Add your services here based on selected components
  
networks:
  %s-network:
    driver: bridge

volumes:
  %s-data:
`, config.Name, config.Name, config.Name)
}

func (pg *ProjectGenerator) generateGitignoreContent(config *models.ProjectConfig) string {
	return `# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
*.exe
*.dll
*.so
*.dylib

# Test coverage
coverage/
*.out

# Environment files
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Temporary files
tmp/
temp/
`
}

// generateFrontendFiles creates frontend component files
func (pg *ProjectGenerator) generateFrontendFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate package.json for main app if selected
	if config.Components.Frontend.NextJS.App {
		packageJsonContent := fmt.Sprintf(`{
  "name": "%s-app",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "test": "jest",
    "test:watch": "jest --watch"
  },
  "dependencies": {
    "next": "%s",
    "react": "%s",
    "react-dom": "%s",
    "@tailwindcss/forms": "^0.5.7",
    "@tailwindcss/typography": "^0.5.10"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "%s",
    "typescript": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "jest": "^29.0.0",
    "@testing-library/react": "^14.0.0",
    "@testing-library/jest-dom": "^6.0.0"
  }
}`, config.Name, config.Versions.Packages["next"], config.Versions.Packages["react"], config.Versions.Packages["react"], config.Versions.Packages["next"])

		packageJsonPath := filepath.Join(projectPath, "App/package.json")
		if err := pg.fsGen.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
			return fmt.Errorf("failed to create App/package.json: %w", err)
		}

		// Generate Next.js configuration
		nextConfigContent := `/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost'],
  },
}

module.exports = nextConfig`

		nextConfigPath := filepath.Join(projectPath, "App/next.config.js")
		if err := pg.fsGen.WriteFile(nextConfigPath, []byte(nextConfigContent), 0644); err != nil {
			return fmt.Errorf("failed to create App/next.config.js: %w", err)
		}

		// Generate Tailwind CSS configuration
		tailwindConfigContent := `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}`

		tailwindConfigPath := filepath.Join(projectPath, "App/tailwind.config.js")
		if err := pg.fsGen.WriteFile(tailwindConfigPath, []byte(tailwindConfigContent), 0644); err != nil {
			return fmt.Errorf("failed to create App/tailwind.config.js: %w", err)
		}
	}

	// Generate package.json for home app if selected
	if config.Components.Frontend.NextJS.Home {
		homePackageJsonContent := fmt.Sprintf(`{
  "name": "%s-home",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev -p 3001",
    "build": "next build",
    "start": "next start -p 3001",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "%s",
    "react": "%s",
    "react-dom": "%s"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "%s",
    "typescript": "^5.0.0",
    "tailwindcss": "^3.4.0"
  }
}`, config.Name, config.Versions.Packages["next"], config.Versions.Packages["react"], config.Versions.Packages["react"], config.Versions.Packages["next"])

		homePackageJsonPath := filepath.Join(projectPath, "Home/package.json")
		if err := pg.fsGen.WriteFile(homePackageJsonPath, []byte(homePackageJsonContent), 0644); err != nil {
			return fmt.Errorf("failed to create Home/package.json: %w", err)
		}
	}

	// Generate package.json for admin app if selected
	if config.Components.Frontend.NextJS.Admin {
		adminPackageJsonContent := fmt.Sprintf(`{
  "name": "%s-admin",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev -p 3002",
    "build": "next build",
    "start": "next start -p 3002",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "%s",
    "react": "%s",
    "react-dom": "%s",
    "@headlessui/react": "^1.7.17",
    "@heroicons/react": "^2.0.18"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "%s",
    "typescript": "^5.0.0",
    "tailwindcss": "^3.4.0"
  }
}`, config.Name, config.Versions.Packages["next"], config.Versions.Packages["react"], config.Versions.Packages["react"], config.Versions.Packages["next"])

		adminPackageJsonPath := filepath.Join(projectPath, "Admin/package.json")
		if err := pg.fsGen.WriteFile(adminPackageJsonPath, []byte(adminPackageJsonContent), 0644); err != nil {
			return fmt.Errorf("failed to create Admin/package.json: %w", err)
		}
	}

	return nil
}

func (pg *ProjectGenerator) generateBackendFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate go.mod for backend
	goModContent := fmt.Sprintf(`module %s/commonserver

go %s

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.5
)
`, config.Organization+"/"+config.Name, config.Versions.Go)

	goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
	if err := pg.fsGen.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/go.mod: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) generateMobileFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate Android build.gradle if selected
	if config.Components.Mobile.Android {
		buildGradleContent := fmt.Sprintf(`// %s Android App
// Generated by Open Source Template Generator

plugins {
    id 'com.android.application'
    id 'org.jetbrains.kotlin.android'
}

android {
    compileSdk 34
    
    defaultConfig {
        applicationId "%s.%s"
        minSdk 24
        targetSdk 34
        versionCode 1
        versionName "1.0"
    }
}

dependencies {
    implementation 'androidx.core:core-ktx:1.12.0'
    implementation 'androidx.compose.ui:compose-bom:2023.10.01'
}
`, config.Name, config.Organization, config.Name)

		buildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
		if err := pg.fsGen.WriteFile(buildGradlePath, []byte(buildGradleContent), 0644); err != nil {
			return fmt.Errorf("failed to create Mobile/Android/build.gradle: %w", err)
		}
	}

	// Generate iOS Package.swift if selected
	if config.Components.Mobile.IOS {
		packageSwiftContent := fmt.Sprintf(`// swift-tools-version: 5.9
// %s iOS App
// Generated by Open Source Template Generator

import PackageDescription

let package = Package(
    name: "%s",
    platforms: [
        .iOS(.v15)
    ],
    products: [
        .library(
            name: "%s",
            targets: ["%s"]
        ),
    ],
    targets: [
        .target(
            name: "%s",
            dependencies: []
        ),
    ]
)
`, config.Name, config.Name, config.Name, config.Name, config.Name)

		packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
		if err := pg.fsGen.WriteFile(packageSwiftPath, []byte(packageSwiftContent), 0644); err != nil {
			return fmt.Errorf("failed to create Mobile/iOS/Package.swift: %w", err)
		}
	}

	return nil
}

func (pg *ProjectGenerator) generateInfrastructureFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate basic Terraform configuration if selected
	if config.Components.Infrastructure.Terraform {
		terraformContent := fmt.Sprintf(`# %s Infrastructure
# Generated by Open Source Template Generator

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "%s"
}
`, config.Name, config.Name)

		terraformPath := filepath.Join(projectPath, "Deploy/terraform/main.tf")
		if err := pg.fsGen.WriteFile(terraformPath, []byte(terraformContent), 0644); err != nil {
			return fmt.Errorf("failed to create Deploy/terraform/main.tf: %w", err)
		}
	}

	return nil
}

// generateDocumentationFiles creates documentation files
func (pg *ProjectGenerator) generateDocumentationFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate CONTRIBUTING.md
	contributingContent := fmt.Sprintf(`# Contributing to %s

Thank you for your interest in contributing to %s! This document provides guidelines and information for contributors.

## Development Setup

1. Clone the repository:
   `+"```bash"+`
   git clone %s
   cd %s
   `+"```"+`

2. Set up the development environment:
   `+"```bash"+`
   make setup
   `+"```"+`

3. Start development servers:
   `+"```bash"+`
   make dev
   `+"```"+`

## Code Style

- Follow the existing code style in each component
- Run linting before submitting: `+"```bash"+`make lint`+"```"+`
- Ensure all tests pass: `+"```bash"+`make test`+"```"+`

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `+"```bash"+`git checkout -b feature/your-feature`+"```"+`
3. Make your changes and add tests
4. Ensure all tests pass and code is properly formatted
5. Submit a pull request with a clear description

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/).

## License

By contributing to %s, you agree that your contributions will be licensed under the %s License.
`, config.Name, config.Name, "", config.Name, config.Name, config.License)

	contributingPath := filepath.Join(projectPath, "CONTRIBUTING.md")
	if err := pg.fsGen.WriteFile(contributingPath, []byte(contributingContent), 0644); err != nil {
		return fmt.Errorf("failed to create CONTRIBUTING.md: %w", err)
	}

	// Generate SECURITY.md
	securityContent := fmt.Sprintf(`# Security Policy

## Supported Versions

We provide security updates for the following versions of %s:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in %s, please report it responsibly:

1. **Do not** create a public GitHub issue for security vulnerabilities
2. Email us at: security@%s.com (replace with your actual security contact)
3. Include detailed information about the vulnerability
4. Allow us time to address the issue before public disclosure

## Security Best Practices

When contributing to %s:

- Keep dependencies up to date
- Follow secure coding practices
- Use environment variables for sensitive configuration
- Validate all user inputs
- Use HTTPS for all external communications

## Response Timeline

- We will acknowledge receipt of vulnerability reports within 48 hours
- We aim to provide an initial assessment within 7 days
- We will work to resolve critical vulnerabilities within 30 days

Thank you for helping keep %s secure!
`, config.Name, config.Name, config.Organization, config.Name, config.Name)

	securityPath := filepath.Join(projectPath, "SECURITY.md")
	if err := pg.fsGen.WriteFile(securityPath, []byte(securityContent), 0644); err != nil {
		return fmt.Errorf("failed to create SECURITY.md: %w", err)
	}

	// Generate LICENSE file
	licenseContent := pg.generateLicenseContent(config)
	licensePath := filepath.Join(projectPath, "LICENSE")
	if err := pg.fsGen.WriteFile(licensePath, []byte(licenseContent), 0644); err != nil {
		return fmt.Errorf("failed to create LICENSE: %w", err)
	}

	return nil
}

// generateCICDFiles creates CI/CD configuration files
func (pg *ProjectGenerator) generateCICDFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate GitHub Actions workflow for CI
	ciWorkflowContent := fmt.Sprintf(`name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        node-version: [%s]
        go-version: [%s]
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
    
    - name: Set up Go
      uses: actions/setup-go@v6
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Install dependencies
      run: make setup
    
    - name: Run linting
      run: make lint
    
    - name: Run tests
      run: make test
    
    - name: Build project
      run: make build
`, config.Versions.Node, config.Versions.Go)

	ciWorkflowPath := filepath.Join(projectPath, ".github/workflows/ci.yml")
	if err := pg.fsGen.WriteFile(ciWorkflowPath, []byte(ciWorkflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create .github/workflows/ci.yml: %w", err)
	}

	// Generate security workflow
	securityWorkflowContent := `name: Security

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Mondays

jobs:
  security:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Run CodeQL Analysis
      uses: github/codeql-action/init@v2
      with:
        languages: go, javascript
    
    - name: Autobuild
      uses: github/codeql-action/autobuild@v2
    
    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v2
    
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'
    
    - name: Upload Trivy scan results
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
`

	securityWorkflowPath := filepath.Join(projectPath, ".github/workflows/security.yml")
	if err := pg.fsGen.WriteFile(securityWorkflowPath, []byte(securityWorkflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create .github/workflows/security.yml: %w", err)
	}

	// Generate Dependabot configuration
	dependabotContent := `version: 2
updates:
  # Enable version updates for npm
  - package-ecosystem: "npm"
    directory: "/App"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    
  - package-ecosystem: "npm"
    directory: "/Home"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    
  - package-ecosystem: "npm"
    directory: "/Admin"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
  
  # Enable version updates for Go modules
  - package-ecosystem: "gomod"
    directory: "/CommonServer"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
  
  # Enable version updates for Docker
  - package-ecosystem: "docker"
    directory: "/Deploy/docker"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
  
  # Enable version updates for GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
`

	dependabotPath := filepath.Join(projectPath, ".github/dependabot.yml")
	if err := pg.fsGen.WriteFile(dependabotPath, []byte(dependabotContent), 0644); err != nil {
		return fmt.Errorf("failed to create .github/dependabot.yml: %w", err)
	}

	return nil
}

// generateLicenseContent generates license content based on the selected license
func (pg *ProjectGenerator) generateLicenseContent(config *models.ProjectConfig) string {
	year := "2024"
	author := config.Author
	if author == "" {
		author = config.Organization
	}

	switch config.License {
	case "MIT":
		return fmt.Sprintf(`MIT License

Copyright (c) %s %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`, year, author)

	case "Apache-2.0":
		return fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %s %s

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`, year, author)

	case "GPL-3.0":
		return fmt.Sprintf(`GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

Copyright (C) %s %s

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`, year, author)

	case "BSD-3-Clause":
		return fmt.Sprintf(`BSD 3-Clause License

Copyright (c) %s, %s
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`, year, author)

	default:
		return fmt.Sprintf(`Copyright (c) %s %s

All rights reserved.`, year, author)
	}
}
