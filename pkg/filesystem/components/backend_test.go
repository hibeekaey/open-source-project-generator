package components

import (
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestBackendGenerator_GenerateFiles(t *testing.T) {
	tests := []struct {
		name          string
		config        *models.ProjectConfig
		expectedFiles []string
		expectedError bool
	}{
		{
			name: "Generate Go Gin backend files",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Backend: models.BackendComponents{
						GoGin: true,
					},
				},
				Versions: &models.VersionConfig{
					Go: "1.22",
				},
			},
			expectedFiles: []string{
				"testproject/CommonServer/go.mod",
				"testproject/CommonServer/main.go",
				"testproject/CommonServer/Dockerfile",
				"testproject/CommonServer/.env.example",
				"testproject/CommonServer/internal/controllers/health.go",
				"testproject/CommonServer/internal/middleware/cors.go",
				"testproject/CommonServer/internal/config/config.go",
				"testproject/CommonServer/Makefile",
			},
			expectedError: false,
		},
		{
			name: "No backend components selected",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Backend: models.BackendComponents{
						GoGin: false,
					},
				},
			},
			expectedFiles: []string{},
			expectedError: false,
		},
		{
			name:          "Nil config should return error",
			config:        nil,
			expectedFiles: []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := NewMockFileSystemOperations()
			bg := NewBackendGenerator(mockFS)

			err := bg.GenerateFiles("testproject", tt.config)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that expected files were created
			for _, expectedFile := range tt.expectedFiles {
				if !mockFS.FileExists(expectedFile) {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}

			// Verify go.mod content if backend is enabled
			if tt.config != nil && tt.config.Components.Backend.GoGin {
				goModPath := "testproject/CommonServer/go.mod"
				content := string(mockFS.GetFileContent(goModPath))

				expectedModule := tt.config.Organization + "/" + tt.config.Name + "/commonserver"
				if !strings.Contains(content, expectedModule) {
					t.Errorf("go.mod should contain module name %s", expectedModule)
				}
			}
		})
	}
}

func TestBackendGenerator_generateGoMod(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
		Versions: &models.VersionConfig{
			Go: "1.22",
		},
	}

	content := bg.generateGoMod(config)

	expectedElements := []string{
		"module testorg/testapp/commonserver",
		"go 1.22",
		"github.com/gin-gonic/gin",
		"github.com/joho/godotenv",
		"gorm.io/gorm",
		"gorm.io/driver/postgres",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("go.mod should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateMainGo(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
	}

	content := bg.generateMainGo(config)

	expectedElements := []string{
		"package main",
		"github.com/gin-gonic/gin",
		"github.com/joho/godotenv",
		"testorg/testapp/commonserver/internal/config",
		"testorg/testapp/commonserver/internal/controllers",
		"testorg/testapp/commonserver/internal/middleware",
		"func main()",
		"router := gin.Default()",
		"router.GET(\"/health\", controllers.HealthCheck)",
		"Starting testapp server",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("main.go should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateHealthController(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := bg.generateHealthController(config)

	expectedElements := []string{
		"package controllers",
		"github.com/gin-gonic/gin",
		"type HealthResponse struct",
		"type StatusResponse struct",
		"func HealthCheck(c *gin.Context)",
		"func Status(c *gin.Context)",
		"Service:   \"testapp\"",
		"testapp API is running",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("health controller should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateCORSMiddleware(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := bg.generateCORSMiddleware(config)

	expectedElements := []string{
		"package middleware",
		"github.com/gin-gonic/gin",
		"func CORS() gin.HandlerFunc",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"c.Request.Method == \"OPTIONS\"",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("CORS middleware should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateConfig(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := bg.generateConfig(config)

	expectedElements := []string{
		"package config",
		"type Config struct",
		"Environment string",
		"Port        string",
		"DatabaseURL string",
		"func Load() *Config",
		"func getEnv(key, fallback string) string",
		"os.Getenv(key)",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("config should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateEnvExample(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := bg.generateEnvExample(config)

	expectedElements := []string{
		"# testapp Backend Configuration",
		"ENVIRONMENT=development",
		"PORT=8080",
		"DATABASE_URL=postgres://username:password@localhost:5432/testapp_db",
		"JWT_SECRET=your-jwt-secret-key",
		"CORS_ORIGINS=http://localhost:3000",
		"LOG_LEVEL=info",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf(".env.example should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateDockerfile(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Go: "1.22",
		},
	}

	content := bg.generateDockerfile(config)

	expectedElements := []string{
		"# testapp Backend Dockerfile",
		"FROM golang:1.22-alpine AS builder",
		"WORKDIR /app",
		"COPY go.mod go.sum ./",
		"RUN go mod download",
		"RUN CGO_ENABLED=0 GOOS=linux go build",
		"FROM alpine:latest",
		"EXPOSE 8080",
		"HEALTHCHECK",
		"CMD [\"./main\"]",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Dockerfile should contain %s", element)
		}
	}
}

func TestBackendGenerator_generateMakefile(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := bg.generateMakefile(config)

	expectedElements := []string{
		"# testapp Backend Makefile",
		".PHONY:",
		"help:",
		"build:",
		"run:",
		"test:",
		"clean:",
		"docker-build:",
		"docker-run:",
		"go build -o bin/testapp main.go",
		"docker build -t testapp-backend .",
		"docker run -p 8080:8080 --env-file .env testapp-backend",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Makefile should contain %s", element)
		}
	}
}

func TestBackendGenerator_WithDefaultVersions(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	bg := NewBackendGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
		// No versions specified
	}

	content := bg.generateGoMod(config)

	// Should use default Go version
	if !strings.Contains(content, "go 1.22") {
		t.Errorf("Should use default Go version when none specified")
	}
}
