package components

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// BackendGenerator handles backend component file generation
type BackendGenerator struct {
	fsOps FileSystemOperations
}

// NewBackendGenerator creates a new backend generator
func NewBackendGenerator(fsOps FileSystemOperations) *BackendGenerator {
	return &BackendGenerator{
		fsOps: fsOps,
	}
}

// GenerateFiles creates backend component files based on configuration
func (bg *BackendGenerator) GenerateFiles(projectPath string, config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Generate Go Gin backend files if selected
	if config.Components.Backend.GoGin {
		if err := bg.generateGoGinFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Go Gin files: %w", err)
		}
	}

	return nil
}

// generateGoGinFiles creates Go Gin backend files
func (bg *BackendGenerator) generateGoGinFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate go.mod
	goModContent := bg.generateGoMod(config)
	goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
	if err := bg.fsOps.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/go.mod: %w", err)
	}

	// Generate main.go
	mainGoContent := bg.generateMainGo(config)
	mainGoPath := filepath.Join(projectPath, "CommonServer/main.go")
	if err := bg.fsOps.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/main.go: %w", err)
	}

	// Generate Dockerfile
	dockerfileContent := bg.generateDockerfile(config)
	dockerfilePath := filepath.Join(projectPath, "CommonServer/Dockerfile")
	if err := bg.fsOps.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/Dockerfile: %w", err)
	}

	// Generate .env.example
	envExampleContent := bg.generateEnvExample(config)
	envExamplePath := filepath.Join(projectPath, "CommonServer/.env.example")
	if err := bg.fsOps.WriteFile(envExamplePath, []byte(envExampleContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/.env.example: %w", err)
	}

	// Generate basic controller
	controllerContent := bg.generateHealthController(config)
	controllerPath := filepath.Join(projectPath, "CommonServer/internal/controllers/health.go")
	if err := bg.fsOps.WriteFile(controllerPath, []byte(controllerContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/internal/controllers/health.go: %w", err)
	}

	// Generate basic middleware
	middlewareContent := bg.generateCORSMiddleware(config)
	middlewarePath := filepath.Join(projectPath, "CommonServer/internal/middleware/cors.go")
	if err := bg.fsOps.WriteFile(middlewarePath, []byte(middlewareContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/internal/middleware/cors.go: %w", err)
	}

	// Generate config package
	configContent := bg.generateConfig(config)
	configPath := filepath.Join(projectPath, "CommonServer/internal/config/config.go")
	if err := bg.fsOps.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/internal/config/config.go: %w", err)
	}

	// Generate Makefile
	makefileContent := bg.generateMakefile(config)
	makefilePath := filepath.Join(projectPath, "CommonServer/Makefile")
	if err := bg.fsOps.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to create CommonServer/Makefile: %w", err)
	}

	return nil
}

// generateGoMod generates go.mod content
func (bg *BackendGenerator) generateGoMod(config *models.ProjectConfig) string {
	goVersion := "1.22"
	if config.Versions != nil && config.Versions.Go != "" {
		goVersion = config.Versions.Go
	}

	return fmt.Sprintf(`module %s/commonserver

go %s

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.5.1
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`, config.Organization+"/"+config.Name, goVersion)
}

// generateMainGo generates main.go content
func (bg *BackendGenerator) generateMainGo(config *models.ProjectConfig) string {
	return fmt.Sprintf(`package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"%s/commonserver/internal/config"
	"%s/commonserver/internal/controllers"
	"%s/commonserver/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", controllers.HealthCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", controllers.Status)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting %s server on port %%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
`, config.Organization+"/"+config.Name, config.Organization+"/"+config.Name, config.Organization+"/"+config.Name, config.Name)
}

// generateHealthController generates health controller content
func (bg *BackendGenerator) generateHealthController(config *models.ProjectConfig) string {
	return fmt.Sprintf(`package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `+"`json:\"status\"`"+`
	Timestamp time.Time `+"`json:\"timestamp\"`"+`
	Service   string    `+"`json:\"service\"`"+`
	Version   string    `+"`json:\"version\"`"+`
}

// StatusResponse represents the status response
type StatusResponse struct {
	Message string `+"`json:\"message\"`"+`
	Service string `+"`json:\"service\"`"+`
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "%s",
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// Status handles status requests
func Status(c *gin.Context) {
	response := StatusResponse{
		Message: "%s API is running",
		Service: "%s",
	}

	c.JSON(http.StatusOK, response)
}
`, config.Name, config.Name, config.Name)
}

// generateCORSMiddleware generates CORS middleware content
func (bg *BackendGenerator) generateCORSMiddleware(config *models.ProjectConfig) string {
	return `package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}`
}

// generateConfig generates config package content
func (bg *BackendGenerator) generateConfig(config *models.ProjectConfig) string {
	return `package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Environment string
	Port        string
	DatabaseURL string
}

// Load loads the application configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}`
}

// generateEnvExample generates .env.example content
func (bg *BackendGenerator) generateEnvExample(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Backend Configuration

# Environment (development, staging, production)
ENVIRONMENT=development

# Server port
PORT=8080

# Database configuration
DATABASE_URL=postgres://username:password@localhost:5432/%s_db?sslmode=disable

# JWT Secret (generate a secure random string for production)
JWT_SECRET=your-jwt-secret-key

# CORS origins (comma-separated list)
CORS_ORIGINS=http://localhost:3000,http://localhost:3001,http://localhost:3002

# Log level (debug, info, warn, error)
LOG_LEVEL=info
`, config.Name, config.Name)
}

// generateDockerfile generates Dockerfile content
func (bg *BackendGenerator) generateDockerfile(config *models.ProjectConfig) string {
	goVersion := "1.22"
	if config.Versions != nil && config.Versions.Go != "" {
		goVersion = config.Versions.Go
	}

	return fmt.Sprintf(`# %s Backend Dockerfile
FROM golang:%s-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
`, config.Name, goVersion)
}

// generateMakefile generates Makefile content for backend
func (bg *BackendGenerator) generateMakefile(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Backend Makefile

.PHONY: help build run test clean dev docker-build docker-run

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%%-20s\033[0m %%s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/%s main.go

run: ## Run the application
	go run main.go

dev: ## Run in development mode with hot reload
	go run main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

docker-build: ## Build Docker image
	docker build -t %s-backend .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env %s-backend

docker-dev: ## Run with docker-compose for development
	docker-compose up --build

migrate-up: ## Run database migrations up
	@echo "Running database migrations..."
	# Add your migration command here

migrate-down: ## Run database migrations down
	@echo "Rolling back database migrations..."
	# Add your rollback command here

.DEFAULT_GOAL := help
`, config.Name, config.Name, config.Name, config.Name)
}
