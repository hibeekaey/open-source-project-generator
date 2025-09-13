package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestDockerTemplateValidation(t *testing.T) {
	// Test that validates the Docker template generates correct configurations
	// without actually building the Docker image

	tempDir := t.TempDir()

	// Create the actual frontend Dockerfile template
	dockerfileTemplate := `FROM {{nodeDockerImage .}} AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY yarn.lock* ./

# Install dependencies
RUN if [ -f package-lock.json ]; then npm ci --only=production; else npm install --only=production; fi && npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM {{nodeDockerImage .}} AS runner

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

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/api/health || exit 1

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]`

	dockerfilePath := filepath.Join(tempDir, "Dockerfile.tmpl")
	err := os.WriteFile(dockerfilePath, []byte(dockerfileTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile template: %v", err)
	}

	// Test with Node.js 20.x configuration
	config := &models.ProjectConfig{
		Name:         "test-docker-validation",
		Organization: "test-org",
		Description:  "Test Docker template validation",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	engine := NewEngine()
	result, err := engine.ProcessTemplate(dockerfilePath, config)
	if err != nil {
		t.Fatalf("Failed to process Dockerfile template: %v", err)
	}

	dockerfileContent := string(result)

	// Validate Docker image versions
	expectedImages := []string{
		"FROM node:20-alpine AS builder",
		"FROM node:20-alpine AS runner",
	}

	for _, expected := range expectedImages {
		if !strings.Contains(dockerfileContent, expected) {
			t.Errorf("Expected Docker image not found: %s\nGenerated Dockerfile:\n%s", expected, dockerfileContent)
		}
	}

	// Validate security best practices
	securityChecks := []string{
		"RUN addgroup --system --gid 1001 nodejs",
		"RUN adduser --system --uid 1001 nextjs",
		"USER nextjs",
		"ENV NODE_ENV=production",
		"HEALTHCHECK",
	}

	for _, check := range securityChecks {
		if !strings.Contains(dockerfileContent, check) {
			t.Errorf("Security best practice not found: %s\nGenerated Dockerfile:\n%s", check, dockerfileContent)
		}
	}

	// Validate npm installation logic
	if !strings.Contains(dockerfileContent, "if [ -f package-lock.json ]") {
		t.Error("Dockerfile should handle both npm ci and npm install cases")
	}

	// Validate multi-stage build
	if !strings.Contains(dockerfileContent, "AS builder") || !strings.Contains(dockerfileContent, "AS runner") {
		t.Error("Dockerfile should use multi-stage build")
	}
}

func TestDockerComposeTemplateValidation(t *testing.T) {
	// Test docker-compose template validation

	tempDir := t.TempDir()

	// Create docker-compose template
	composeTemplate := `version: '3.8'

services:
  app:
    build:
      context: ./App
      dockerfile: ../templates/infrastructure/docker/frontend.Dockerfile.tmpl
      target: runner
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=http://api:8080
    depends_on:
      - api
    networks:
      - {{.Name}}-network
    restart: unless-stopped
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  {{.Name}}-network:
    driver: bridge`

	composePath := filepath.Join(tempDir, "docker-compose.yml.tmpl")
	err := os.WriteFile(composePath, []byte(composeTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create docker-compose template: %v", err)
	}

	config := &models.ProjectConfig{
		Name:         "test-compose-validation",
		Organization: "test-org",
		Description:  "Test docker-compose validation",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	engine := NewEngine()
	result, err := engine.ProcessTemplate(composePath, config)
	if err != nil {
		t.Fatalf("Failed to process docker-compose template: %v", err)
	}

	composeContent := string(result)

	// Validate compose configuration
	expectedContent := []string{
		"version: '3.8'",
		"dockerfile: ../templates/infrastructure/docker/frontend.Dockerfile.tmpl",
		"NODE_ENV=production",
		"test-compose-validation-network",
		"healthcheck:",
		"deploy:",
		"resources:",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(composeContent, expected) {
			t.Errorf("Expected docker-compose content not found: %s\nGenerated compose:\n%s", expected, composeContent)
		}
	}
}

func TestDockerSecurityTemplateValidation(t *testing.T) {
	// Test security configuration template

	tempDir := t.TempDir()

	// Use the actual security template content
	securityTemplate := `# Docker Security Configuration
images:
  base_images:
    allowed:
      - "alpine:*"
      - "node:*-alpine"
      - "golang:*-alpine"
      - "postgres:*-alpine"
      - "redis:*-alpine"
    
  vulnerability_scanning:
    enabled: true
    fail_on: "HIGH"
    ignore_unfixed: false

runtime:
  apparmor:
    enabled: true
    profiles:
      - docker-default
      
  seccomp:
    enabled: true
    profile: runtime/default`

	securityPath := filepath.Join(tempDir, "security.yml.tmpl")
	err := os.WriteFile(securityPath, []byte(securityTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create security template: %v", err)
	}

	config := &models.ProjectConfig{
		Name:         "test-security-validation",
		Organization: "test-org",
		Description:  "Test security configuration",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	engine := NewEngine()
	result, err := engine.ProcessTemplate(securityPath, config)
	if err != nil {
		t.Fatalf("Failed to process security template: %v", err)
	}

	securityContent := string(result)

	// Validate security configuration
	securityChecks := []string{
		`"node:*-alpine"`,
		"vulnerability_scanning:",
		"enabled: true",
		"fail_on: \"HIGH\"",
		"apparmor:",
		"seccomp:",
		"runtime/default",
	}

	for _, check := range securityChecks {
		if !strings.Contains(securityContent, check) {
			t.Errorf("Security configuration not found: %s\nGenerated security config:\n%s", check, securityContent)
		}
	}

	// Ensure Node.js 20 alpine images are allowed
	if !strings.Contains(securityContent, "node:*-alpine") {
		t.Error("Security configuration should allow Node.js alpine images")
	}
}
