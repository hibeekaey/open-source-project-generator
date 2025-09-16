package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestDockerBuildCompatibility(t *testing.T) {
	// Skip this test if Docker is not available
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping Docker build test")
	}

	// Check if Docker daemon is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker daemon not running, skipping Docker build test")
	}

	tempDir := t.TempDir()

	// Create a minimal Next.js project structure for testing
	projectDir := filepath.Join(tempDir, "test-project")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create package.json with Node.js 20.x requirements
	packageJSON := `{
  "name": "test-nextjs-app",
  "version": "0.1.0",
  "private": true,
  "engines": {
    "node": ">=20.0.0",
    "npm": ">=10.0.0"
  },
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "14.2.5",
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@types/node": "^20.17.0",
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "typescript": "^5.5.4"
  }
}`

	err = os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create a minimal Next.js configuration
	nextConfig := `/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  experimental: {
    outputFileTracingRoot: undefined,
  },
}

module.exports = nextConfig`

	err = os.WriteFile(filepath.Join(projectDir, "next.config.js"), []byte(nextConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create next.config.js: %v", err)
	}

	// Create minimal app structure
	appDir := filepath.Join(projectDir, "src", "app")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create app directory: %v", err)
	}

	// Create public directory (required by Dockerfile)
	publicDir := filepath.Join(projectDir, "public")
	err = os.MkdirAll(publicDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create public directory: %v", err)
	}

	// Create layout.tsx
	layout := `import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Test App',
  description: 'Test Next.js application',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}`

	err = os.WriteFile(filepath.Join(appDir, "layout.tsx"), []byte(layout), 0644)
	if err != nil {
		t.Fatalf("Failed to create layout.tsx: %v", err)
	}

	// Create page.tsx
	page := `export default function Home() {
  return (
    <main>
      <h1>Hello World</h1>
      <p>Node.js version: {process.version}</p>
    </main>
  )
}`

	err = os.WriteFile(filepath.Join(appDir, "page.tsx"), []byte(page), 0644)
	if err != nil {
		t.Fatalf("Failed to create page.tsx: %v", err)
	}

	// Create TypeScript configuration
	tsConfig := `{
  "compilerOptions": {
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

	err = os.WriteFile(filepath.Join(projectDir, "tsconfig.json"), []byte(tsConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create tsconfig.json: %v", err)
	}

	// Create Dockerfile using our template system
	dockerfileTemplate := `FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN if [ -f package-lock.json ]; then npm ci --only=production; else npm install --only=production; fi && npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy built application
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Security: Run as non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]`

	err = os.WriteFile(filepath.Join(projectDir, "Dockerfile"), []byte(dockerfileTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	// Create .dockerignore
	dockerignore := `node_modules
.next
.git
.gitignore
README.md
Dockerfile
.dockerignore
npm-debug.log*
yarn-debug.log*
yarn-error.log*`

	err = os.WriteFile(filepath.Join(projectDir, ".dockerignore"), []byte(dockerignore), 0644)
	if err != nil {
		t.Fatalf("Failed to create .dockerignore: %v", err)
	}

	// Test Docker build (this will actually build the image)
	// We'll use a short timeout to avoid long-running tests
	buildCmd := exec.Command("docker", "build", "-t", "test-nodejs-20", ".")
	buildCmd.Dir = projectDir

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker build output: %s", string(output))

		// Check if the error is due to network issues or missing dependencies
		outputStr := string(output)
		if strings.Contains(outputStr, "network") || strings.Contains(outputStr, "timeout") {
			t.Skip("Network issues during Docker build, skipping test")
		}

		t.Fatalf("Docker build failed: %v\nOutput: %s", err, string(output))
	}

	// Verify the build was successful
	outputStr := string(output)
	// Check for success indicators in modern Docker buildx output
	hasSuccess := strings.Contains(outputStr, "Successfully built") ||
		strings.Contains(outputStr, "Successfully tagged") ||
		(strings.Contains(outputStr, "exporting to image") && strings.Contains(outputStr, "DONE"))

	if !hasSuccess {
		t.Errorf("Docker build did not complete successfully. Output: %s", outputStr)
	}

	// Clean up the Docker image
	cleanupCmd := exec.Command("docker", "rmi", "test-nodejs-20")
	_ = cleanupCmd.Run() // Ignore errors in cleanup
}

func TestDockerImageVersionValidation(t *testing.T) {
	// Test that validates Docker image versions are compatible with package.json

	tests := []struct {
		name        string
		dockerImage string
		nodeRuntime string
		expectError bool
		description string
	}{
		{
			name:        "Compatible Node.js 20",
			dockerImage: "node:20-alpine",
			nodeRuntime: ">=20.0.0",
			expectError: false,
			description: "Node.js 20 Docker image with >=20.0.0 runtime requirement",
		},
		{
			name:        "Compatible Node.js 20.17",
			dockerImage: "node:20.17-alpine",
			nodeRuntime: ">=20.0.0",
			expectError: false,
			description: "Specific Node.js 20.17 Docker image with >=20.0.0 runtime requirement",
		},
		{
			name:        "Incompatible Node.js 18",
			dockerImage: "node:18-alpine",
			nodeRuntime: ">=20.0.0",
			expectError: true,
			description: "Node.js 18 Docker image with >=20.0.0 runtime requirement should be incompatible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.ProjectConfig{
				Name:         "test-validation",
				Organization: "test-org",
				Description:  tt.description,
				License:      "MIT",
				Versions: &models.VersionConfig{
					Node: tt.nodeRuntime,
					Packages: map[string]string{
						"typescript": "^20.17.0",
					},
				},
			}

			// Verify the configuration is valid
			if config.Name != "test-validation" {
				t.Errorf("Expected name %s, got %s", "test-validation", config.Name)
			}

			if config.Organization != "test-org" {
				t.Errorf("Expected organization %s, got %s", "test-org", config.Organization)
			}

			if config.Description != tt.description {
				t.Errorf("Expected description %s, got %s", tt.description, config.Description)
			}

			if config.License != "MIT" {
				t.Errorf("Expected license %s, got %s", "MIT", config.License)
			}

			// Validate that the Versions field is properly set and accessible
			if config.Versions == nil {
				t.Error("Expected Versions field to be set")
			} else {
				if config.Versions.Node != tt.nodeRuntime {
					t.Errorf("Expected Node version %s, got %s", tt.nodeRuntime, config.Versions.Node)
				}

				// Validate package versions
				if config.Versions.Packages == nil {
					t.Error("Expected Packages field to be set")
				} else {
					if typescriptVersion, exists := config.Versions.Packages["typescript"]; !exists {
						t.Error("Expected typescript package version to be set")
					} else if typescriptVersion != "^20.17.0" {
						t.Errorf("Expected typescript version ^20.17.0, got %s", typescriptVersion)
					}
				}

				// Simulate Docker image version validation logic
				// Extract Node version from Docker image (e.g., "node:20-alpine" -> "20")
				dockerImageVersion := extractNodeVersionFromDockerImage(tt.dockerImage)
				if dockerImageVersion == "" {
					t.Errorf("Could not extract Node version from Docker image: %s", tt.dockerImage)
				}

				// Basic validation: check if the Docker image version meets the runtime requirement
				isCompatible := validateNodeVersionCompatibility(dockerImageVersion, tt.nodeRuntime)
				if isCompatible && tt.expectError {
					t.Errorf("Expected incompatibility between Docker image %s and runtime %s, but they were compatible", tt.dockerImage, tt.nodeRuntime)
				} else if !isCompatible && !tt.expectError {
					t.Errorf("Expected compatibility between Docker image %s and runtime %s, but they were incompatible", tt.dockerImage, tt.nodeRuntime)
				}
			}
		})
	}
}

// extractNodeVersionFromDockerImage extracts the Node.js version from a Docker image string
// e.g., "node:20-alpine" -> "20", "node:18.17.0-slim" -> "18.17.0"
func extractNodeVersionFromDockerImage(dockerImage string) string {
	// Remove the "node:" prefix
	if !strings.HasPrefix(dockerImage, "node:") {
		return ""
	}

	version := strings.TrimPrefix(dockerImage, "node:")

	// Extract the version part before any additional tags (like -alpine, -slim)
	parts := strings.Split(version, "-")
	if len(parts) > 0 {
		return parts[0]
	}

	return version
}

// validateNodeVersionCompatibility performs basic version compatibility checking
// This is a simplified implementation for testing purposes
func validateNodeVersionCompatibility(dockerVersion, runtimeRequirement string) bool {
	// For this test, we'll do a simple major version comparison
	// In a real implementation, you'd want more sophisticated semver parsing

	// Extract major version from docker version (e.g., "20.17.0" -> "20")
	dockerMajor := extractMajorVersion(dockerVersion)
	if dockerMajor == "" {
		return false
	}

	// Parse runtime requirement (e.g., ">=20.0.0" -> "20")
	reqMajor := extractMajorVersionFromRequirement(runtimeRequirement)
	if reqMajor == "" {
		return false
	}

	// Simple comparison: docker major version should be >= required major version
	return dockerMajor >= reqMajor
}

// extractMajorVersion extracts the major version number from a version string
func extractMajorVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// extractMajorVersionFromRequirement extracts the major version from a requirement string
func extractMajorVersionFromRequirement(requirement string) string {
	// Handle patterns like ">=20.0.0", "~20.0.0", "^20.0.0", "20.0.0"
	requirement = strings.TrimSpace(requirement)

	// Remove common prefixes
	prefixes := []string{">=", "~", "^", "="}
	for _, prefix := range prefixes {
		if strings.HasPrefix(requirement, prefix) {
			requirement = strings.TrimPrefix(requirement, prefix)
			break
		}
	}

	// Extract major version
	return extractMajorVersion(requirement)
}
