package generators

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ConfigurationGenerator handles configuration file generation
type ConfigurationGenerator struct {
	fsOps FileSystemOperationsInterface
}

// NewConfigurationGenerator creates a new configuration generator
func NewConfigurationGenerator(fsOps FileSystemOperationsInterface) *ConfigurationGenerator {
	return &ConfigurationGenerator{
		fsOps: fsOps,
	}
}

// GenerateFrontendFiles creates frontend component files
func (cg *ConfigurationGenerator) GenerateFrontendFiles(projectPath string, config *models.ProjectConfig) error {
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
		if err := cg.fsOps.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644); err != nil {
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
		if err := cg.fsOps.WriteFile(nextConfigPath, []byte(nextConfigContent), 0644); err != nil {
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
		if err := cg.fsOps.WriteFile(tailwindConfigPath, []byte(tailwindConfigContent), 0644); err != nil {
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
		if err := cg.fsOps.WriteFile(homePackageJsonPath, []byte(homePackageJsonContent), 0644); err != nil {
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
		if err := cg.fsOps.WriteFile(adminPackageJsonPath, []byte(adminPackageJsonContent), 0644); err != nil {
			return fmt.Errorf("failed to create Admin/package.json: %w", err)
		}
	}

	return nil
}

// GenerateBackendFiles creates backend configuration files
func (cg *ConfigurationGenerator) GenerateBackendFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate go.mod for backend
	goModContent := fmt.Sprintf(`module %s/commonserver

go %s

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.5
)
`, config.Organization+"/"+config.Name, config.Versions.Go)

	goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
	if err := cg.fsOps.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/go.mod: %w", err)
	}

	return nil
}

// GenerateMobileFiles creates mobile configuration files
func (cg *ConfigurationGenerator) GenerateMobileFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate Android build.gradle if selected
	if config.Components.Mobile.Android {
		buildGradleContent := fmt.Sprintf(`// %s Android App
// Generated by Open Source Project Generator

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
		if err := cg.fsOps.WriteFile(buildGradlePath, []byte(buildGradleContent), 0644); err != nil {
			return fmt.Errorf("failed to create Mobile/Android/build.gradle: %w", err)
		}
	}

	// Generate iOS Package.swift if selected
	if config.Components.Mobile.IOS {
		packageSwiftContent := fmt.Sprintf(`// swift-tools-version: 5.9
// %s iOS App
// Generated by Open Source Project Generator

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
		if err := cg.fsOps.WriteFile(packageSwiftPath, []byte(packageSwiftContent), 0644); err != nil {
			return fmt.Errorf("failed to create Mobile/iOS/Package.swift: %w", err)
		}
	}

	return nil
}

// GenerateInfrastructureFiles creates infrastructure configuration files
func (cg *ConfigurationGenerator) GenerateInfrastructureFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate basic Terraform configuration if selected
	if config.Components.Infrastructure.Terraform {
		terraformContent := fmt.Sprintf(`# %s Infrastructure
# Generated by Open Source Project Generator

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
		if err := cg.fsOps.WriteFile(terraformPath, []byte(terraformContent), 0644); err != nil {
			return fmt.Errorf("failed to create Deploy/terraform/main.tf: %w", err)
		}
	}

	return nil
}
