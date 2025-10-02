package components

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FrontendGenerator handles frontend component file generation
type FrontendGenerator struct {
	fsOps FileSystemOperations
}

// FileSystemOperations interface for filesystem operations
type FileSystemOperations interface {
	WriteFile(path string, content []byte, perm os.FileMode) error
	EnsureDirectory(path string) error
	FileExists(path string) bool
}

// NewFrontendGenerator creates a new frontend generator
func NewFrontendGenerator(fsOps FileSystemOperations) *FrontendGenerator {
	return &FrontendGenerator{
		fsOps: fsOps,
	}
}

// GenerateFiles creates frontend component files based on configuration
func (fg *FrontendGenerator) GenerateFiles(projectPath string, config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Generate files for each selected frontend component
	if config.Components.Frontend.NextJS.App {
		if err := fg.generateAppFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate app files: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Home {
		if err := fg.generateHomeFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate home files: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Admin {
		if err := fg.generateAdminFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate admin files: %w", err)
		}
	}

	if config.Components.Frontend.NextJS.Shared {
		if err := fg.generateSharedFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate shared files: %w", err)
		}
	}

	return nil
}

// generateAppFiles creates main application files
func (fg *FrontendGenerator) generateAppFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate package.json for main app
	packageJsonContent := fg.generateAppPackageJson(config)
	packageJsonPath := filepath.Join(projectPath, "App/package.json")
	if err := fg.fsOps.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
		return fmt.Errorf("failed to create App/package.json: %w", err)
	}

	// Generate Next.js configuration
	nextConfigContent := fg.generateNextConfig()
	nextConfigPath := filepath.Join(projectPath, "App/next.config.js")
	if err := fg.fsOps.WriteFile(nextConfigPath, []byte(nextConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create App/next.config.js: %w", err)
	}

	// Generate Tailwind CSS configuration
	tailwindConfigContent := fg.generateTailwindConfig()
	tailwindConfigPath := filepath.Join(projectPath, "App/tailwind.config.js")
	if err := fg.fsOps.WriteFile(tailwindConfigPath, []byte(tailwindConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create App/tailwind.config.js: %w", err)
	}

	// Generate TypeScript configuration
	tsConfigContent := fg.generateTSConfig()
	tsConfigPath := filepath.Join(projectPath, "App/tsconfig.json")
	if err := fg.fsOps.WriteFile(tsConfigPath, []byte(tsConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create App/tsconfig.json: %w", err)
	}

	// Generate ESLint configuration
	eslintConfigContent := fg.generateESLintConfig()
	eslintConfigPath := filepath.Join(projectPath, "App/.eslintrc.json")
	if err := fg.fsOps.WriteFile(eslintConfigPath, []byte(eslintConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create App/.eslintrc.json: %w", err)
	}

	return nil
}

// generateHomeFiles creates home/landing page files
func (fg *FrontendGenerator) generateHomeFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate package.json for home app
	packageJsonContent := fg.generateHomePackageJson(config)
	packageJsonPath := filepath.Join(projectPath, "Home/package.json")
	if err := fg.fsOps.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
		return fmt.Errorf("failed to create Home/package.json: %w", err)
	}

	// Generate Next.js configuration for home
	nextConfigContent := fg.generateNextConfig()
	nextConfigPath := filepath.Join(projectPath, "Home/next.config.js")
	if err := fg.fsOps.WriteFile(nextConfigPath, []byte(nextConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create Home/next.config.js: %w", err)
	}

	return nil
}

// generateAdminFiles creates admin dashboard files
func (fg *FrontendGenerator) generateAdminFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate package.json for admin app
	packageJsonContent := fg.generateAdminPackageJson(config)
	packageJsonPath := filepath.Join(projectPath, "Admin/package.json")
	if err := fg.fsOps.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
		return fmt.Errorf("failed to create Admin/package.json: %w", err)
	}

	// Generate Next.js configuration for admin
	nextConfigContent := fg.generateNextConfig()
	nextConfigPath := filepath.Join(projectPath, "Admin/next.config.js")
	if err := fg.fsOps.WriteFile(nextConfigPath, []byte(nextConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create Admin/next.config.js: %w", err)
	}

	return nil
}

// generateSharedFiles creates shared component files
func (fg *FrontendGenerator) generateSharedFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate package.json for shared components
	packageJsonContent := fg.generateSharedPackageJson(config)
	packageJsonPath := filepath.Join(projectPath, "Shared/package.json")
	if err := fg.fsOps.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
		return fmt.Errorf("failed to create Shared/package.json: %w", err)
	}

	return nil
}

// generateAppPackageJson generates package.json content for main app
func (fg *FrontendGenerator) generateAppPackageJson(config *models.ProjectConfig) string {
	nextVersion := "14.0.0"
	reactVersion := "18.0.0"

	if config.Versions != nil && config.Versions.Packages != nil {
		if v, ok := config.Versions.Packages["next"]; ok {
			nextVersion = v
		}
		if v, ok := config.Versions.Packages["react"]; ok {
			reactVersion = v
		}
	}

	return fmt.Sprintf(`{
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
}`, config.Name, nextVersion, reactVersion, reactVersion, nextVersion)
}

// generateHomePackageJson generates package.json content for home app
func (fg *FrontendGenerator) generateHomePackageJson(config *models.ProjectConfig) string {
	nextVersion := "14.0.0"
	reactVersion := "18.0.0"

	if config.Versions != nil && config.Versions.Packages != nil {
		if v, ok := config.Versions.Packages["next"]; ok {
			nextVersion = v
		}
		if v, ok := config.Versions.Packages["react"]; ok {
			reactVersion = v
		}
	}

	return fmt.Sprintf(`{
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
}`, config.Name, nextVersion, reactVersion, reactVersion, nextVersion)
}

// generateAdminPackageJson generates package.json content for admin app
func (fg *FrontendGenerator) generateAdminPackageJson(config *models.ProjectConfig) string {
	nextVersion := "14.0.0"
	reactVersion := "18.0.0"

	if config.Versions != nil && config.Versions.Packages != nil {
		if v, ok := config.Versions.Packages["next"]; ok {
			nextVersion = v
		}
		if v, ok := config.Versions.Packages["react"]; ok {
			reactVersion = v
		}
	}

	return fmt.Sprintf(`{
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
}`, config.Name, nextVersion, reactVersion, reactVersion, nextVersion)
}

// generateSharedPackageJson generates package.json content for shared components
func (fg *FrontendGenerator) generateSharedPackageJson(config *models.ProjectConfig) string {
	reactVersion := "18.0.0"

	if config.Versions != nil && config.Versions.Packages != nil {
		if v, ok := config.Versions.Packages["react"]; ok {
			reactVersion = v
		}
	}

	return fmt.Sprintf(`{
  "name": "%s-shared",
  "version": "0.1.0",
  "private": true,
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch",
    "lint": "eslint src --ext .ts,.tsx",
    "test": "jest"
  },
  "peerDependencies": {
    "react": "%s",
    "react-dom": "%s"
  },
  "devDependencies": {
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "typescript": "^5.0.0",
    "eslint": "^8.0.0",
    "jest": "^29.0.0"
  }
}`, config.Name, reactVersion, reactVersion)
}

// generateNextConfig generates Next.js configuration
func (fg *FrontendGenerator) generateNextConfig() string {
	return `/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost'],
  },
}

module.exports = nextConfig`
}

// generateTailwindConfig generates Tailwind CSS configuration
func (fg *FrontendGenerator) generateTailwindConfig() string {
	return `/** @type {import('tailwindcss').Config} */
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
}

// generateTSConfig generates TypeScript configuration
func (fg *FrontendGenerator) generateTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "es6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}`
}

// generateESLintConfig generates ESLint configuration
func (fg *FrontendGenerator) generateESLintConfig() string {
	return `{
  "extends": ["next/core-web-vitals"],
  "rules": {
    "prefer-const": "error",
    "no-unused-vars": "warn",
    "no-console": "warn"
  }
}`
}
