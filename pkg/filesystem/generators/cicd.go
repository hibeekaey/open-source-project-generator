package generators

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// CICDGenerator handles CI/CD configuration file generation
type CICDGenerator struct {
	fsOps FileSystemOperationsInterface
}

// NewCICDGenerator creates a new CI/CD generator
func NewCICDGenerator(fsOps FileSystemOperationsInterface) *CICDGenerator {
	return &CICDGenerator{
		fsOps: fsOps,
	}
}

// GenerateCICDFiles creates CI/CD configuration files
func (cg *CICDGenerator) GenerateCICDFiles(projectPath string, config *models.ProjectConfig) error {
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
	if err := cg.fsOps.WriteFile(ciWorkflowPath, []byte(ciWorkflowContent), 0644); err != nil {
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
	if err := cg.fsOps.WriteFile(securityWorkflowPath, []byte(securityWorkflowContent), 0644); err != nil {
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
	if err := cg.fsOps.WriteFile(dependabotPath, []byte(dependabotContent), 0644); err != nil {
		return fmt.Errorf("failed to create .github/dependabot.yml: %w", err)
	}

	return nil
}
