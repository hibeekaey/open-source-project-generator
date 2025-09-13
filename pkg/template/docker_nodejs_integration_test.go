package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestDockerNodeJSVersionConsistency(t *testing.T) {
	// Test that Docker configurations use consistent Node.js 20.x versions
	// and are compatible with package.json requirements

	tempDir := t.TempDir()

	// Create a frontend Dockerfile template
	dockerfileContent := `FROM {{nodeDockerImage .}} AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force
COPY . .
RUN npm run build

FROM {{nodeDockerImage .}} AS runner
WORKDIR /app
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"
USER nextjs
EXPOSE 3000
CMD ["node", "server.js"]`

	dockerfilePath := filepath.Join(tempDir, "Dockerfile.tmpl")
	err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile template: %v", err)
	}

	// Create a package.json template
	packageJSONContent := `{
  "name": "{{.Name}}",
  "version": "0.1.0",
  "private": true,
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "next": "14.2.5",
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@types/node": "{{nodeTypesVersion .}}",
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "typescript": "^5.5.4"
  }
}`

	packageJSONPath := filepath.Join(tempDir, "package.json.tmpl")
	err = os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json template: %v", err)
	}

	// Create test configuration with Node.js 20.x
	config := &models.ProjectConfig{
		Name:         "test-nodejs-docker",
		Organization: "test-org",
		Description:  "Test Node.js Docker consistency",
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

	// Process Dockerfile template
	dockerfileResult, err := engine.ProcessTemplate(dockerfilePath, config)
	if err != nil {
		t.Fatalf("Failed to process Dockerfile template: %v", err)
	}

	// Process package.json template
	packageJSONResult, err := engine.ProcessTemplate(packageJSONPath, config)
	if err != nil {
		t.Fatalf("Failed to process package.json template: %v", err)
	}

	dockerfileStr := string(dockerfileResult)
	packageJSONStr := string(packageJSONResult)

	// Verify Docker image uses Node.js 20
	expectedDockerImages := []string{
		"FROM node:20-alpine AS builder",
		"FROM node:20-alpine AS runner",
	}

	for _, expected := range expectedDockerImages {
		if !strings.Contains(dockerfileStr, expected) {
			t.Errorf("Expected Docker image not found: %s\nDockerfile: %s", expected, dockerfileStr)
		}
	}

	// Verify package.json uses consistent Node.js 20.x versions
	expectedPackageJSONContent := []string{
		`"node": ">=20.0.0"`,
		`"npm": ">=10.0.0"`,
		`"@types/node": "^20.17.0"`,
	}

	for _, expected := range expectedPackageJSONContent {
		if !strings.Contains(packageJSONStr, expected) {
			t.Errorf("Expected package.json content not found: %s\npackage.json: %s", expected, packageJSONStr)
		}
	}

	// Verify version consistency between Docker and package.json
	// Both should reference Node.js 20.x
	if !strings.Contains(dockerfileStr, "node:20-alpine") {
		t.Error("Dockerfile should use node:20-alpine base image")
	}

	if !strings.Contains(packageJSONStr, ">=20.0.0") {
		t.Error("package.json should require Node.js >=20.0.0")
	}
}

func TestDockerSecurityConfiguration(t *testing.T) {
	// Test that security configuration allows Node.js 20 images

	tempDir := t.TempDir()

	// Create security configuration template
	securityContent := `images:
  base_images:
    allowed:
      - "alpine:*"
      - "node:*-alpine"
      - "golang:*-alpine"
      - "postgres:*-alpine"
      - "redis:*-alpine"`

	securityPath := filepath.Join(tempDir, "security.yml.tmpl")
	err := os.WriteFile(securityPath, []byte(securityContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create security template: %v", err)
	}

	config := &models.ProjectConfig{
		Name:         "test-security",
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

	resultStr := string(result)

	// Verify that Node.js alpine images are allowed
	if !strings.Contains(resultStr, `"node:*-alpine"`) {
		t.Error("Security configuration should allow node:*-alpine images")
	}
}

func TestDockerComposeNodeJSConsistency(t *testing.T) {
	// Test that docker-compose configurations use consistent Node.js versions

	tempDir := t.TempDir()

	// Create docker-compose template that references the frontend Dockerfile
	composeContent := `version: '3.8'
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
    networks:
      - {{.Name}}-network
networks:
  {{.Name}}-network:
    driver: bridge`

	composePath := filepath.Join(tempDir, "docker-compose.yml.tmpl")
	err := os.WriteFile(composePath, []byte(composeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create docker-compose template: %v", err)
	}

	config := &models.ProjectConfig{
		Name:         "test-compose",
		Organization: "test-org",
		Description:  "Test docker-compose configuration",
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

	resultStr := string(result)

	// Verify that the compose file references the correct Dockerfile
	if !strings.Contains(resultStr, "frontend.Dockerfile.tmpl") {
		t.Error("docker-compose should reference frontend.Dockerfile.tmpl")
	}

	// Verify network configuration
	if !strings.Contains(resultStr, "test-compose-network") {
		t.Error("docker-compose should create project-specific network")
	}
}
