package template

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestE2EVersionSubstitution(t *testing.T) {
	// Create a temporary directory structure mimicking real templates
	tempDir := t.TempDir()

	// Create frontend template structure
	frontendDir := filepath.Join(tempDir, "frontend", "nextjs-app")
	err := os.MkdirAll(frontendDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create frontend directory: %v", err)
	}

	// Create package.json template similar to the real one
	packageJSONTemplate := `{
  "name": "{{.Name}}-app",
  "version": "0.1.0",
  "private": true,
  "description": "{{.Description}} - Main Application",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}",
    "react": "{{.Versions.React}}",
    "next": "{{.Versions.NextJS}}"
  }
}`

	packageJSONPath := filepath.Join(frontendDir, "package.json.tmpl")
	err = os.WriteFile(packageJSONPath, []byte(packageJSONTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json template: %v", err)
	}

	// Create Dockerfile template
	dockerfileTemplate := `FROM {{nodeDockerImage .}} AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM {{nodeDockerImage .}} AS runner
WORKDIR /app
COPY --from=builder /app/.next/standalone ./
USER nextjs
EXPOSE 3000
CMD ["node", "server.js"]`

	dockerfilePath := filepath.Join(frontendDir, "Dockerfile.tmpl")
	err = os.WriteFile(dockerfilePath, []byte(dockerfileTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile template: %v", err)
	}

	// Create project configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project for version substitution",
		License:      "MIT",
		Versions: &models.VersionConfig{
			React:  "18.2.0",
			NextJS: "15.0.0",
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
				Description:  "Node.js 20 LTS for production stability",
			},
		},
	}

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Process templates
	engine := NewEngine()

	// Process package.json
	packageResult, err := engine.ProcessTemplate(packageJSONPath, config)
	if err != nil {
		t.Fatalf("Failed to process package.json template: %v", err)
	}

	// Write processed package.json
	outputPackageJSON := filepath.Join(outputDir, "package.json")
	err = os.WriteFile(outputPackageJSON, packageResult, 0644)
	if err != nil {
		t.Fatalf("Failed to write processed package.json: %v", err)
	}

	// Process Dockerfile
	dockerResult, err := engine.ProcessTemplate(dockerfilePath, config)
	if err != nil {
		t.Fatalf("Failed to process Dockerfile template: %v", err)
	}

	// Write processed Dockerfile
	outputDockerfile := filepath.Join(outputDir, "Dockerfile")
	err = os.WriteFile(outputDockerfile, dockerResult, 0644)
	if err != nil {
		t.Fatalf("Failed to write processed Dockerfile: %v", err)
	}

	// Verify package.json content
	var packageData map[string]interface{}
	err = json.Unmarshal(packageResult, &packageData)
	if err != nil {
		t.Fatalf("Failed to parse generated package.json: %v", err)
	}

	// Verify engines
	engines, ok := packageData["engines"].(map[string]interface{})
	if !ok {
		t.Fatal("engines field not found or not an object")
	}

	if engines["node"] != ">=20.0.0" {
		t.Errorf("Expected node engine >=20.0.0, got %v", engines["node"])
	}

	if engines["npm"] != ">=10.0.0" {
		t.Errorf("Expected npm engine >=10.0.0, got %v", engines["npm"])
	}

	// Verify dependencies
	dependencies, ok := packageData["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatal("dependencies field not found or not an object")
	}

	if dependencies["@types/node"] != "^20.17.0" {
		t.Errorf("Expected @types/node ^20.17.0, got %v", dependencies["@types/node"])
	}

	if dependencies["react"] != "18.2.0" {
		t.Errorf("Expected react 18.2.0, got %v", dependencies["react"])
	}

	if dependencies["next"] != "15.0.0" {
		t.Errorf("Expected next 15.0.0, got %v", dependencies["next"])
	}

	// Verify Dockerfile content
	dockerContent := string(dockerResult)
	expectedDockerLines := []string{
		"FROM node:20-alpine AS builder",
		"FROM node:20-alpine AS runner",
	}

	for _, expectedLine := range expectedDockerLines {
		if !containsLine(dockerContent, expectedLine) {
			t.Errorf("Expected Dockerfile to contain line: %s\nActual content:\n%s", expectedLine, dockerContent)
		}
	}

	t.Logf("Successfully processed templates with version substitution")
	t.Logf("Generated package.json: %s", string(packageResult))
	t.Logf("Generated Dockerfile: %s", dockerContent)
}

func TestE2EVersionConsistencyAcrossTemplates(t *testing.T) {
	// Create multiple frontend templates to test consistency
	tempDir := t.TempDir()

	templates := map[string]string{
		"nextjs-app/package.json.tmpl": `{
  "name": "{{.Name}}-app",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
		"nextjs-admin/package.json.tmpl": `{
  "name": "{{.Name}}-admin",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
		"shared-components/package.json.tmpl": `{
  "name": "@{{.Organization}}/{{.Name}}-ui",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "devDependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
	}

	// Create template files
	for templatePath, content := range templates {
		fullPath := filepath.Join(tempDir, templatePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", templatePath, err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template %s: %v", templatePath, err)
		}
	}

	// Create configuration
	config := &models.ProjectConfig{
		Name:         "consistency-test",
		Organization: "test-org",
		Description:  "Test project for version consistency",
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

	// Process all templates
	engine := NewEngine()
	results := make(map[string]map[string]interface{})

	for templatePath := range templates {
		fullPath := filepath.Join(tempDir, templatePath)
		result, err := engine.ProcessTemplate(fullPath, config)
		if err != nil {
			t.Fatalf("Failed to process template %s: %v", templatePath, err)
		}

		var data map[string]interface{}
		err = json.Unmarshal(result, &data)
		if err != nil {
			t.Fatalf("Failed to parse result for %s: %v", templatePath, err)
		}

		results[templatePath] = data
	}

	// Verify consistency across all templates
	expectedVersions := map[string]string{
		"node":        ">=20.0.0",
		"npm":         ">=10.0.0",
		"@types/node": "^20.17.0",
	}

	for templatePath, data := range results {
		// Check engines
		engines, ok := data["engines"].(map[string]interface{})
		if !ok {
			t.Errorf("Template %s: engines field not found", templatePath)
			continue
		}

		if engines["node"] != expectedVersions["node"] {
			t.Errorf("Template %s: expected node engine %s, got %v",
				templatePath, expectedVersions["node"], engines["node"])
		}

		if engines["npm"] != expectedVersions["npm"] {
			t.Errorf("Template %s: expected npm engine %s, got %v",
				templatePath, expectedVersions["npm"], engines["npm"])
		}

		// Check @types/node in dependencies or devDependencies
		var typesNodeVersion interface{}
		if deps, ok := data["dependencies"].(map[string]interface{}); ok {
			typesNodeVersion = deps["@types/node"]
		} else if devDeps, ok := data["devDependencies"].(map[string]interface{}); ok {
			typesNodeVersion = devDeps["@types/node"]
		}

		if typesNodeVersion != expectedVersions["@types/node"] {
			t.Errorf("Template %s: expected @types/node %s, got %v",
				templatePath, expectedVersions["@types/node"], typesNodeVersion)
		}
	}

	t.Logf("Successfully verified version consistency across %d templates", len(templates))
}

// Helper function to check if content contains a specific line
func containsLine(content, line string) bool {
	lines := splitLines(content)
	for _, l := range lines {
		if strings.TrimSpace(l) == strings.TrimSpace(line) {
			return true
		}
	}
	return false
}

// Helper function to split content into lines
func splitLines(content string) []string {
	return strings.Split(content, "\n")
}
