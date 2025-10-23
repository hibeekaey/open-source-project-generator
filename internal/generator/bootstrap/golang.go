package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// GoExecutor handles Go backend project generation
type GoExecutor struct {
	*BaseExecutor
}

// NewGoExecutor creates a new Go executor
func NewGoExecutor() *GoExecutor {
	return &GoExecutor{
		BaseExecutor: NewBaseExecutor("go"),
	}
}

// Execute generates a Go backend project with Gin framework
func (ge *GoExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Get project configuration
	projectName, ok := spec.Config["name"].(string)
	if !ok || projectName == "" {
		return nil, fmt.Errorf("project name is required in config")
	}

	moduleName, ok := spec.Config["module"].(string)
	if !ok || moduleName == "" {
		return nil, fmt.Errorf("module name is required in config")
	}

	// Create project directory
	projectDir := filepath.Join(spec.TargetDir, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	result := &models.ExecutionResult{
		OutputDir: projectDir,
		ToolUsed:  "go",
	}

	// Step 1: Initialize Go module
	if err := ge.initGoModule(ctx, projectDir, moduleName); err != nil {
		result.Success = false
		result.Stderr = err.Error()
		return result, fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Step 2: Install Gin framework
	framework := "gin"
	if fw, ok := spec.Config["framework"].(string); ok && fw != "" {
		framework = fw
	}

	if err := ge.installDependencies(ctx, projectDir, framework); err != nil {
		result.Success = false
		result.Stderr = err.Error()
		return result, fmt.Errorf("failed to install dependencies: %w", err)
	}

	// Step 3: Generate basic server setup
	if err := ge.generateServerFiles(projectDir, moduleName, framework); err != nil {
		result.Success = false
		result.Stderr = err.Error()
		return result, fmt.Errorf("failed to generate server files: %w", err)
	}

	// Step 4: Run go mod tidy
	if err := ge.tidyModules(ctx, projectDir); err != nil {
		result.Success = false
		result.Stderr = err.Error()
		return result, fmt.Errorf("failed to tidy modules: %w", err)
	}

	result.Success = true
	result.Stdout = fmt.Sprintf("Successfully created Go backend project at %s", projectDir)
	return result, nil
}

// SupportsComponent checks if this executor supports the given component type
func (ge *GoExecutor) SupportsComponent(componentType string) bool {
	return componentType == "go-backend" || componentType == "go" || componentType == "backend"
}

// GetDefaultFlags returns default flags for Go generation
func (ge *GoExecutor) GetDefaultFlags(componentType string) []string {
	if !ge.SupportsComponent(componentType) {
		return []string{}
	}

	return []string{"mod", "init"}
}

// ValidateConfig validates component-specific configuration
func (ge *GoExecutor) ValidateConfig(config map[string]interface{}) error {
	// Validate name
	if name, ok := config["name"].(string); !ok || name == "" {
		return fmt.Errorf("name is required and must be a string")
	}

	// Validate module
	if module, ok := config["module"].(string); !ok || module == "" {
		return fmt.Errorf("module is required and must be a valid Go module path")
	}

	// Validate port if provided
	if port, exists := config["port"]; exists {
		switch v := port.(type) {
		case int:
			if v < 1 || v > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
		case float64:
			if v < 1 || v > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
		default:
			return fmt.Errorf("port must be a number")
		}
	}

	// Validate framework if provided
	if framework, exists := config["framework"]; exists {
		if fw, ok := framework.(string); ok {
			validFrameworks := map[string]bool{"gin": true, "echo": true, "fiber": true}
			if !validFrameworks[fw] {
				return fmt.Errorf("framework must be one of: gin, echo, fiber")
			}
		} else {
			return fmt.Errorf("framework must be a string")
		}
	}

	return nil
}

// initGoModule initializes a Go module
func (ge *GoExecutor) initGoModule(ctx context.Context, projectDir, moduleName string) error {
	execSpec := &BootstrapSpec{
		ComponentType: "go",
		TargetDir:     projectDir,
		Flags:         []string{"mod", "init", moduleName},
	}

	result, err := ge.BaseExecutor.Execute(ctx, execSpec)
	if err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("go mod init failed with exit code %d", result.ExitCode)
	}

	return nil
}

// installDependencies installs required Go dependencies
func (ge *GoExecutor) installDependencies(ctx context.Context, projectDir, framework string) error {
	var packages []string

	switch framework {
	case "gin":
		packages = []string{
			"github.com/gin-gonic/gin@v1.11.0",
			"github.com/gin-contrib/cors@v1.7.6",
		}
	case "echo":
		packages = []string{
			"github.com/labstack/echo/v4@v4.13.4",
			"github.com/labstack/echo/v4/middleware@v4.13.4",
		}
	case "fiber":
		packages = []string{
			"github.com/gofiber/fiber/v2@v2.52.9",
		}
	default:
		packages = []string{
			"github.com/gin-gonic/gin@v1.11.0",
			"github.com/gin-contrib/cors@v1.7.6",
		}
	}

	for _, pkg := range packages {
		execSpec := &BootstrapSpec{
			ComponentType: "go",
			TargetDir:     projectDir,
			Flags:         []string{"get", pkg},
		}

		result, err := ge.BaseExecutor.Execute(ctx, execSpec)
		if err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}

		if !result.Success {
			return fmt.Errorf("failed to install %s with exit code %d", pkg, result.ExitCode)
		}
	}

	return nil
}

// generateServerFiles creates basic server files
func (ge *GoExecutor) generateServerFiles(projectDir, moduleName, framework string) error {
	// Create main.go
	mainContent := ge.generateMainFile(moduleName, framework)
	mainPath := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	// Create .gitignore
	gitignoreContent := ge.generateGitignore()
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Create README.md
	readmeContent := ge.generateReadme(moduleName, framework)
	readmePath := filepath.Join(projectDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	return nil
}

// generateMainFile generates the main.go content based on framework
func (ge *GoExecutor) generateMainFile(moduleName, framework string) string {
	switch framework {
	case "gin":
		return fmt.Sprintf(`package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "%s",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to %s API",
			})
		})
	}

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %%s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %%v", err)
	}
}
`, moduleName, moduleName)

	default:
		return ge.generateMainFile(moduleName, "gin")
	}
}

// generateGitignore generates .gitignore content
func (ge *GoExecutor) generateGitignore() string {
	return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`
}

// generateReadme generates README.md content
func (ge *GoExecutor) generateReadme(moduleName, framework string) string {
	return fmt.Sprintf(`# %s

A Go backend service built with %s framework.

## Getting Started

### Prerequisites

- Go 1.25 or higher

### Installation

1. Install dependencies:
   `+"```bash"+`
   go mod download
   `+"```"+`

2. Run the server:
   `+"```bash"+`
   go run main.go
   `+"```"+`

The server will start on http://localhost:8080

### API Endpoints

- `+"`GET /health`"+` - Health check endpoint
- `+"`GET /api/v1/`"+` - API root endpoint

## Development

### Running Tests

`+"```bash"+`
go test ./...
`+"```"+`

### Building

`+"```bash"+`
go build -o server
`+"```"+`

## Project Structure

`+"```"+`
.
├── main.go          # Application entry point
├── go.mod           # Go module definition
└── go.sum           # Go module checksums
`+"```"+`

## License

MIT
`, moduleName, framework)
}

// tidyModules runs go mod tidy
func (ge *GoExecutor) tidyModules(ctx context.Context, projectDir string) error {
	execSpec := &BootstrapSpec{
		ComponentType: "go",
		TargetDir:     projectDir,
		Flags:         []string{"mod", "tidy"},
	}

	result, err := ge.BaseExecutor.Execute(ctx, execSpec)
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("go mod tidy failed with exit code %d", result.ExitCode)
	}

	return nil
}

// GetManualSteps returns manual steps required after Go generation
func (ge *GoExecutor) GetManualSteps(spec *BootstrapSpec) []string {
	return []string{
		"Navigate to the project directory",
		"Run 'go run main.go' to start the server",
		"The server will be available at http://localhost:8080",
		"Test the health endpoint: curl http://localhost:8080/health",
	}
}
