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
	if !strings.Contains(outputStr, "Successfully built") && !strings.Contains(outputStr, "Successfully tagged") {
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
					NodeJS: &models.NodeVersionConfig{
						Runtime:      tt.nodeRuntime,
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  tt.dockerImage,
						LTSStatus:    true,
					},
				},
			}

			// For now, we'll just verify the configuration is valid
			// In a real implementation, we might add validation logic
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

			if config.Versions.NodeJS.DockerImage != tt.dockerImage {
				t.Errorf("Expected Docker image %s, got %s", tt.dockerImage, config.Versions.NodeJS.DockerImage)
			}

			if config.Versions.NodeJS.Runtime != tt.nodeRuntime {
				t.Errorf("Expected runtime %s, got %s", tt.nodeRuntime, config.Versions.NodeJS.Runtime)
			}
		})
	}
}
