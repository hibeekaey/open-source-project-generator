package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIntegrationManager(t *testing.T) {
	im := NewIntegrationManager("/test/project", true)
	assert.NotNil(t, im)
	assert.Equal(t, "/test/project", im.projectRoot)
	assert.True(t, im.verbose)
}

func TestGenerateDockerCompose(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		wantErr    bool
		validate   func(t *testing.T, compose string)
	}{
		{
			name: "nextjs and go-backend components",
			components: []*models.Component{
				{
					Type:        "nextjs",
					Name:        "web-app",
					Path:        "/project/App",
					GeneratedAt: time.Now(),
				},
				{
					Type:        "go-backend",
					Name:        "api-server",
					Path:        "/project/CommonServer",
					GeneratedAt: time.Now(),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, compose string) {
				assert.Contains(t, compose, "version: '3.8'")
				assert.Contains(t, compose, "services:")
				assert.Contains(t, compose, "web-app:")
				assert.Contains(t, compose, "api-server:")
				assert.Contains(t, compose, "networks:")
				assert.Contains(t, compose, "app-network:")
				assert.Contains(t, compose, "volumes:")
				assert.Contains(t, compose, "3000:3000")
				assert.Contains(t, compose, "8080:8080")
			},
		},
		{
			name: "single nextjs component",
			components: []*models.Component{
				{
					Type:        "nextjs",
					Name:        "frontend",
					Path:        "/project/App",
					GeneratedAt: time.Now(),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, compose string) {
				assert.Contains(t, compose, "frontend:")
				assert.Contains(t, compose, "NEXT_PUBLIC_API_URL")
				assert.Contains(t, compose, "npm run dev")
			},
		},
		{
			name: "mobile components excluded",
			components: []*models.Component{
				{
					Type:        "android",
					Name:        "mobile-android",
					Path:        "/project/Mobile/android",
					GeneratedAt: time.Now(),
				},
				{
					Type:        "ios",
					Name:        "mobile-ios",
					Path:        "/project/Mobile/ios",
					GeneratedAt: time.Now(),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, compose string) {
				// Mobile components should not generate services
				assert.NotContains(t, compose, "mobile-android:")
				assert.NotContains(t, compose, "mobile-ios:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			compose, err := im.GenerateDockerCompose(tt.components)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, compose)
				}
			}
		})
	}
}

func TestGenerateServiceDefinition(t *testing.T) {
	im := NewIntegrationManager("/test/project", false)

	tests := []struct {
		name      string
		component *models.Component
		wantErr   bool
		validate  func(t *testing.T, service string)
	}{
		{
			name: "nextjs service",
			component: &models.Component{
				Type: "nextjs",
				Name: "web-app",
				Path: "/project/App",
			},
			wantErr: false,
			validate: func(t *testing.T, service string) {
				assert.Contains(t, service, "web-app:")
				assert.Contains(t, service, "build:")
				assert.Contains(t, service, "ports:")
				assert.Contains(t, service, "3000:3000")
				assert.Contains(t, service, "NEXT_PUBLIC_API_URL")
				assert.Contains(t, service, "depends_on:")
				assert.Contains(t, service, "backend")
			},
		},
		{
			name: "go-backend service",
			component: &models.Component{
				Type: "go-backend",
				Name: "api-server",
				Path: "/project/CommonServer",
			},
			wantErr: false,
			validate: func(t *testing.T, service string) {
				assert.Contains(t, service, "api-server:")
				assert.Contains(t, service, "8080:8080")
				assert.Contains(t, service, "GO_ENV")
				assert.Contains(t, service, "go run main.go")
			},
		},
		{
			name: "android component returns empty",
			component: &models.Component{
				Type: "android",
				Name: "mobile-app",
				Path: "/project/Mobile/android",
			},
			wantErr: false,
			validate: func(t *testing.T, service string) {
				assert.Empty(t, service)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := im.generateServiceDefinition(tt.component)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, service)
				}
			}
		})
	}
}

func TestSanitizeServiceName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"WebApp", "webapp"},
		{"api_server", "api-server"},
		{"My Service", "my-service"},
		{"frontend-app", "frontend-app"},
		{"Backend_API_v2", "backend-api-v2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeServiceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrate(t *testing.T) {
	im := NewIntegrationManager("/test/project", false)
	ctx := context.Background()

	components := []*models.Component{
		{
			Type: "nextjs",
			Name: "web-app",
			Path: "/project/App",
		},
	}

	config := &models.IntegrationConfig{
		GenerateDockerCompose: true,
		GenerateScripts:       true,
		APIEndpoints: map[string]string{
			"backend": "http://localhost:8080",
		},
	}

	// Note: This will partially work as some methods are stubs
	err := im.Integrate(ctx, components, config)
	// We expect no error even with stub methods
	assert.NoError(t, err)
}

func TestConfigureEnvironment(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		config     *models.IntegrationConfig
		wantErr    bool
	}{
		{
			name: "nextjs and go-backend with API endpoints",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			config: &models.IntegrationConfig{
				APIEndpoints: map[string]string{
					"backend": "http://localhost:8080",
				},
				SharedEnvironment: map[string]string{
					"LOG_LEVEL": "debug",
				},
			},
			wantErr: false,
		},
		{
			name: "single component with shared environment",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "frontend",
					Path: "/project/App",
				},
			},
			config: &models.IntegrationConfig{
				SharedEnvironment: map[string]string{
					"NEXT_PUBLIC_APP_NAME": "MyApp",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			err := im.ConfigureEnvironment(tt.components, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateRootEnvFile(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		config     *models.IntegrationConfig
		validate   func(t *testing.T, content string)
	}{
		{
			name: "with API endpoints and shared environment",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			config: &models.IntegrationConfig{
				APIEndpoints: map[string]string{
					"backend": "http://localhost:8080",
				},
				SharedEnvironment: map[string]string{
					"LOG_LEVEL": "debug",
				},
			},
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "# Shared Environment Configuration")
				assert.Contains(t, content, "BACKEND_URL=http://localhost:8080")
				assert.Contains(t, content, "LOG_LEVEL=debug")
				assert.Contains(t, content, "NEXT_PUBLIC_API_URL")
				assert.Contains(t, content, "PORT=8080")
			},
		},
		{
			name: "minimal configuration",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "app",
					Path: "/project/App",
				},
			},
			config: &models.IntegrationConfig{},
			validate: func(t *testing.T, content string) {
				assert.Contains(t, content, "# Component Configuration")
				assert.Contains(t, content, "NEXT_PUBLIC_API_URL")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			content, err := im.generateRootEnvFile(tt.components, tt.config)

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, content)
			}
		})
	}
}

func TestConfigureComponentEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		component *models.Component
		config    *models.IntegrationConfig
		wantErr   bool
	}{
		{
			name: "nextjs component",
			component: &models.Component{
				Type: "nextjs",
				Name: "web-app",
				Path: "/project/App",
			},
			config: &models.IntegrationConfig{
				APIEndpoints: map[string]string{
					"backend": "http://api.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "go-backend component",
			component: &models.Component{
				Type: "go-backend",
				Name: "api-server",
				Path: "/project/CommonServer",
			},
			config: &models.IntegrationConfig{
				SharedEnvironment: map[string]string{
					"LOG_LEVEL": "info",
				},
			},
			wantErr: false,
		},
		{
			name: "unsupported component type",
			component: &models.Component{
				Type: "android",
				Name: "mobile-app",
				Path: "/project/Mobile/android",
			},
			config:  &models.IntegrationConfig{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			err := im.configureComponentEnvironment(tt.component, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateScripts(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		wantErr    bool
	}{
		{
			name: "nextjs and go-backend components",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			wantErr: false,
		},
		{
			name: "single component",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "frontend",
					Path: "/project/App",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			err := im.GenerateScripts(tt.components)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateBuildScript(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		validate   func(t *testing.T, script string)
	}{
		{
			name: "nextjs and go-backend",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			validate: func(t *testing.T, script string) {
				assert.Contains(t, script, "#!/bin/bash")
				assert.Contains(t, script, "npm install")
				assert.Contains(t, script, "npm run build")
				assert.Contains(t, script, "go mod download")
				assert.Contains(t, script, "go build")
			},
		},
		{
			name: "android and ios",
			components: []*models.Component{
				{
					Type: "android",
					Name: "mobile-android",
					Path: "/project/Mobile/android",
				},
				{
					Type: "ios",
					Name: "mobile-ios",
					Path: "/project/Mobile/ios",
				},
			},
			validate: func(t *testing.T, script string) {
				assert.Contains(t, script, "./gradlew build")
				assert.Contains(t, script, "xcodebuild")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			script, err := im.generateBuildScript(tt.components)

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, script)
			}
		})
	}
}

func TestGenerateDevScript(t *testing.T) {
	im := NewIntegrationManager("/test/project", false)

	components := []*models.Component{
		{
			Type: "nextjs",
			Name: "web-app",
			Path: "/project/App",
		},
		{
			Type: "go-backend",
			Name: "api-server",
			Path: "/project/CommonServer",
		},
	}

	script, err := im.generateDevScript(components)

	require.NoError(t, err)
	assert.Contains(t, script, "#!/bin/bash")
	assert.Contains(t, script, "npm run dev")
	assert.Contains(t, script, "go run main.go")
	assert.Contains(t, script, "cleanup()")
	assert.Contains(t, script, "trap cleanup")
}

func TestGenerateProdScript(t *testing.T) {
	im := NewIntegrationManager("/test/project", false)

	components := []*models.Component{
		{
			Type: "nextjs",
			Name: "web-app",
			Path: "/project/App",
		},
		{
			Type: "go-backend",
			Name: "api-server",
			Path: "/project/CommonServer",
		},
	}

	script, err := im.generateProdScript(components)

	require.NoError(t, err)
	assert.Contains(t, script, "#!/bin/bash")
	assert.Contains(t, script, "npm start")
	assert.Contains(t, script, "./bin/server")
	assert.Contains(t, script, "cleanup()")
}

func TestGenerateDockerScript(t *testing.T) {
	im := NewIntegrationManager("/test/project", false)

	components := []*models.Component{
		{
			Type: "nextjs",
			Name: "web-app",
			Path: "/project/App",
		},
	}

	script, err := im.generateDockerScript(components)

	require.NoError(t, err)
	assert.Contains(t, script, "#!/bin/bash")
	assert.Contains(t, script, "docker-compose")
	assert.Contains(t, script, "up|down|logs|build|restart")
	assert.Contains(t, script, "case $COMMAND in")
}

func TestGenerateDocumentation(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		config     *models.IntegrationConfig
		wantErr    bool
	}{
		{
			name: "full stack project",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			config: &models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			err := im.GenerateDocumentation(tt.components, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateMainReadme(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		config     *models.IntegrationConfig
		validate   func(t *testing.T, readme string)
	}{
		{
			name: "nextjs and go-backend with docker",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			config: &models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
			},
			validate: func(t *testing.T, readme string) {
				assert.Contains(t, readme, "# Project Overview")
				assert.Contains(t, readme, "## Table of Contents")
				assert.Contains(t, readme, "## Project Structure")
				assert.Contains(t, readme, "## Components")
				assert.Contains(t, readme, "## Getting Started")
				assert.Contains(t, readme, "## Development")
				assert.Contains(t, readme, "## Production")
				assert.Contains(t, readme, "## Docker")
				assert.Contains(t, readme, "web-app")
				assert.Contains(t, readme, "api-server")
				assert.Contains(t, readme, "docker-compose.yml")
				assert.Contains(t, readme, "Next.js")
				assert.Contains(t, readme, "Go backend")
			},
		},
		{
			name: "mobile components",
			components: []*models.Component{
				{
					Type: "android",
					Name: "mobile-android",
					Path: "/project/Mobile/android",
				},
				{
					Type: "ios",
					Name: "mobile-ios",
					Path: "/project/Mobile/ios",
				},
			},
			config: &models.IntegrationConfig{
				GenerateDockerCompose: false,
				GenerateScripts:       false,
			},
			validate: func(t *testing.T, readme string) {
				assert.Contains(t, readme, "Android")
				assert.Contains(t, readme, "iOS")
				assert.Contains(t, readme, "Kotlin")
				assert.Contains(t, readme, "Swift")
				assert.NotContains(t, readme, "docker-compose.yml")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			readme, err := im.generateMainReadme(tt.components, tt.config)

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, readme)
			}
		})
	}
}

func TestGenerateTroubleshootingGuide(t *testing.T) {
	tests := []struct {
		name       string
		components []*models.Component
		validate   func(t *testing.T, guide string)
	}{
		{
			name: "nextjs and go-backend",
			components: []*models.Component{
				{
					Type: "nextjs",
					Name: "web-app",
					Path: "/project/App",
				},
				{
					Type: "go-backend",
					Name: "api-server",
					Path: "/project/CommonServer",
				},
			},
			validate: func(t *testing.T, guide string) {
				assert.Contains(t, guide, "# Troubleshooting Guide")
				assert.Contains(t, guide, "## Common Issues")
				assert.Contains(t, guide, "Port Already in Use")
				assert.Contains(t, guide, "Environment Variables")
				assert.Contains(t, guide, "Next.js Build Errors")
				assert.Contains(t, guide, "Go Module Issues")
				assert.Contains(t, guide, "Docker Issues")
			},
		},
		{
			name: "mobile components",
			components: []*models.Component{
				{
					Type: "android",
					Name: "mobile-android",
					Path: "/project/Mobile/android",
				},
				{
					Type: "ios",
					Name: "mobile-ios",
					Path: "/project/Mobile/ios",
				},
			},
			validate: func(t *testing.T, guide string) {
				assert.Contains(t, guide, "Android Build Failures")
				assert.Contains(t, guide, "iOS Build Failures")
				assert.Contains(t, guide, "Gradle")
				assert.Contains(t, guide, "Xcode")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIntegrationManager("/test/project", false)
			guide, err := im.generateTroubleshootingGuide(tt.components)

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, guide)
			}
		})
	}
}
