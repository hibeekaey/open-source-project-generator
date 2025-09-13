package validation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// SetupEngine handles post-generation setup and verification
type SetupEngine struct {
	timeout time.Duration
}

// NewSetupEngine creates a new setup engine
func NewSetupEngine() *SetupEngine {
	return &SetupEngine{
		timeout: 5 * time.Minute, // Default timeout for setup operations
	}
}

// componentHandler defines the interface for component-specific operations
type componentHandler func(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error

// processComponents runs component handlers based on configuration
func (s *SetupEngine) processComponents(projectPath string, config *models.ProjectConfig, handlers map[string]componentHandler) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Process frontend components
	if config.Components.Frontend.MainApp || config.Components.Frontend.Home || config.Components.Frontend.Admin {
		if handler, exists := handlers["frontend"]; exists {
			if err := handler(projectPath, config, result); err != nil {
				return nil, fmt.Errorf("failed to process frontend components: %w", err)
			}
		}
	}

	// Process backend components
	if config.Components.Backend.API {
		if handler, exists := handlers["backend"]; exists {
			if err := handler(projectPath, config, result); err != nil {
				return nil, fmt.Errorf("failed to process backend components: %w", err)
			}
		}
	}

	// Process mobile components
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		if handler, exists := handlers["mobile"]; exists {
			if err := handler(projectPath, config, result); err != nil {
				return nil, fmt.Errorf("failed to process mobile components: %w", err)
			}
		}
	}

	// Process infrastructure components
	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		if handler, exists := handlers["infrastructure"]; exists {
			if err := handler(projectPath, config, result); err != nil {
				return nil, fmt.Errorf("failed to process infrastructure components: %w", err)
			}
		}
	}

	return result, nil
}

// SetupProject performs automated setup for a generated project
func (s *SetupEngine) SetupProject(projectPath string, config *models.ProjectConfig) (*models.ValidationResult, error) {
	handlers := map[string]componentHandler{
		"frontend":       s.setupFrontendComponents,
		"backend":        s.setupBackendComponents,
		"mobile":         s.setupMobileComponents,
		"infrastructure": s.setupInfrastructureComponents,
	}
	return s.processComponents(projectPath, config, handlers)
}

// VerifyProject verifies that a generated project can build and run
func (s *SetupEngine) VerifyProject(projectPath string, config *models.ProjectConfig) (*models.ValidationResult, error) {
	handlers := map[string]componentHandler{
		"frontend":       s.verifyFrontendComponents,
		"backend":        s.verifyBackendComponents,
		"mobile":         s.verifyMobileComponents,
		"infrastructure": s.verifyInfrastructureComponents,
	}
	return s.processComponents(projectPath, config, handlers)
}

// setupFrontendComponents sets up frontend applications
func (s *SetupEngine) setupFrontendComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	frontendDirs := []string{}

	if config.Components.Frontend.MainApp {
		frontendDirs = append(frontendDirs, "frontend")
	}
	if config.Components.Frontend.Home {
		frontendDirs = append(frontendDirs, "home")
	}
	if config.Components.Frontend.Admin {
		frontendDirs = append(frontendDirs, "admin")
	}

	for _, dir := range frontendDirs {
		frontendPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "FrontendSetup",
				Message: fmt.Sprintf("Frontend directory not found: %s", dir),
			})
			continue
		}

		// Install dependencies
		if err := s.runCommand(frontendPath, "npm", "install"); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "FrontendSetup",
				Tag:     "dependency_install",
				Value:   dir,
				Message: fmt.Sprintf("Failed to install dependencies for %s: %s", dir, err.Error()),
			})
			continue
		}

		// Run type checking if TypeScript is present
		tsConfigPath := filepath.Join(frontendPath, "tsconfig.json")
		if _, err := os.Stat(tsConfigPath); err == nil {
			if err := s.runCommand(frontendPath, "npx", "tsc", "--noEmit"); err != nil {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "FrontendSetup",
					Message: fmt.Sprintf("TypeScript type checking failed for %s: %s", dir, err.Error()),
				})
			}
		}
	}

	return nil
}

// setupBackendComponents sets up backend services
func (s *SetupEngine) setupBackendComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	backendPath := filepath.Join(projectPath, "backend")
	if _, err := os.Stat(backendPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "BackendSetup",
			Message: "Backend directory not found",
		})
		return nil
	}

	// Download Go dependencies
	if err := s.runCommand(backendPath, "go", "mod", "download"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "BackendSetup",
			Tag:     "dependency_download",
			Value:   "go mod download",
			Message: fmt.Sprintf("Failed to download Go dependencies: %s", err.Error()),
		})
		return nil
	}

	// Tidy Go modules
	if err := s.runCommand(backendPath, "go", "mod", "tidy"); err != nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "BackendSetup",
			Message: fmt.Sprintf("Go mod tidy failed: %s", err.Error()),
		})
	}

	return nil
}

// setupMobileComponents sets up mobile applications
func (s *SetupEngine) setupMobileComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	// Setup Android
	if config.Components.Mobile.Android {
		androidPath := filepath.Join(projectPath, "mobile", "android")
		if _, err := os.Stat(androidPath); err == nil {
			// Check if Gradle wrapper exists
			gradlewPath := filepath.Join(androidPath, "gradlew")
			if _, err := os.Stat(gradlewPath); err == nil {
				// Make gradlew executable
				if err := os.Chmod(gradlewPath, 0755); err != nil {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "AndroidSetup",
						Message: fmt.Sprintf("Failed to make gradlew executable: %s", err.Error()),
					})
				}
			}
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "AndroidSetup",
				Message: "Android directory not found",
			})
		}
	}

	// Setup iOS
	if config.Components.Mobile.IOS {
		iosPath := filepath.Join(projectPath, "mobile", "ios")
		if _, err := os.Stat(iosPath); err == nil {
			// Check if Podfile exists and install pods
			podfilePath := filepath.Join(iosPath, "Podfile")
			if _, err := os.Stat(podfilePath); err == nil {
				if err := s.runCommand(iosPath, "pod", "install"); err != nil {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "IOSSetup",
						Message: fmt.Sprintf("Pod install failed: %s", err.Error()),
					})
				}
			}
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "IOSSetup",
				Message: "iOS directory not found",
			})
		}
	}

	return nil
}

// setupInfrastructureComponents sets up infrastructure components
func (s *SetupEngine) setupInfrastructureComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	// Setup Terraform
	if config.Components.Infrastructure.Terraform {
		terraformFiles := []string{"main.tf", "variables.tf", "outputs.tf"}
		hasTerraform := false
		for _, file := range terraformFiles {
			if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
				hasTerraform = true
				break
			}
		}

		if hasTerraform {
			// Initialize Terraform
			if err := s.runCommand(projectPath, "terraform", "init"); err != nil {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "TerraformSetup",
					Message: fmt.Sprintf("Terraform init failed: %s", err.Error()),
				})
			}

			// Validate Terraform configuration
			if err := s.runCommand(projectPath, "terraform", "validate"); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "TerraformSetup",
					Tag:     "validation",
					Value:   "terraform validate",
					Message: fmt.Sprintf("Terraform validation failed: %s", err.Error()),
				})
			}
		}
	}

	return nil
}

// verifyFrontendComponents verifies frontend applications can build
func (s *SetupEngine) verifyFrontendComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	frontendDirs := []string{}

	if config.Components.Frontend.MainApp {
		frontendDirs = append(frontendDirs, "frontend")
	}
	if config.Components.Frontend.Home {
		frontendDirs = append(frontendDirs, "home")
	}
	if config.Components.Frontend.Admin {
		frontendDirs = append(frontendDirs, "admin")
	}

	for _, dir := range frontendDirs {
		frontendPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
			continue
		}

		// Check if build script exists
		packageJSONPath := filepath.Join(frontendPath, "package.json")
		if _, err := os.Stat(packageJSONPath); err == nil {
			// Try to build the project
			if err := s.runCommand(frontendPath, "npm", "run", "build"); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "FrontendVerification",
					Tag:     "build",
					Value:   dir,
					Message: fmt.Sprintf("Frontend build failed for %s: %s", dir, err.Error()),
				})
			}
		}
	}

	return nil
}

// verifyBackendComponents verifies backend services can build
func (s *SetupEngine) verifyBackendComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	backendPath := filepath.Join(projectPath, "backend")
	if _, err := os.Stat(backendPath); os.IsNotExist(err) {
		return nil
	}

	// Try to build the Go project
	if err := s.runCommand(backendPath, "go", "build", "."); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "BackendVerification",
			Tag:     "build",
			Value:   "go build",
			Message: fmt.Sprintf("Backend build failed: %s", err.Error()),
		})
	}

	// Run tests if they exist
	if err := s.runCommand(backendPath, "go", "test", "./..."); err != nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "BackendVerification",
			Message: fmt.Sprintf("Backend tests failed: %s", err.Error()),
		})
	}

	return nil
}

// verifyMobileComponents verifies mobile applications can build
func (s *SetupEngine) verifyMobileComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	// Verify Android
	if config.Components.Mobile.Android {
		androidPath := filepath.Join(projectPath, "mobile", "android")
		if _, err := os.Stat(androidPath); err == nil {
			gradlewPath := filepath.Join(androidPath, "gradlew")
			if _, err := os.Stat(gradlewPath); err == nil {
				// Try to build Android project
				if err := s.runCommand(androidPath, "./gradlew", "assembleDebug"); err != nil {
					result.Valid = false
					result.Errors = append(result.Errors, models.ValidationError{
						Field:   "AndroidVerification",
						Tag:     "build",
						Value:   "gradlew assembleDebug",
						Message: fmt.Sprintf("Android build failed: %s", err.Error()),
					})
				}
			}
		}
	}

	// Verify iOS
	if config.Components.Mobile.IOS {
		iosPath := filepath.Join(projectPath, "mobile", "ios")
		if _, err := os.Stat(iosPath); err == nil {
			// Look for .xcodeproj or .xcworkspace
			entries, err := os.ReadDir(iosPath)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() && (strings.HasSuffix(entry.Name(), ".xcodeproj") || strings.HasSuffix(entry.Name(), ".xcworkspace")) {
						// Try to build iOS project
						scheme := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
						if err := s.runCommand(iosPath, "xcodebuild", "-scheme", scheme, "-configuration", "Debug", "build"); err != nil {
							result.Valid = false
							result.Errors = append(result.Errors, models.ValidationError{
								Field:   "IOSVerification",
								Tag:     "build",
								Value:   "xcodebuild",
								Message: fmt.Sprintf("iOS build failed: %s", err.Error()),
							})
						}
						break
					}
				}
			}
		}
	}

	return nil
}

// verifyInfrastructureComponents verifies infrastructure components
func (s *SetupEngine) verifyInfrastructureComponents(projectPath string, config *models.ProjectConfig, result *models.ValidationResult) error {
	if config.Components.Infrastructure.Docker {
		s.verifyDockerComponents(projectPath, result)
	}

	if config.Components.Infrastructure.Kubernetes {
		s.verifyKubernetesComponents(projectPath, result)
	}

	return nil
}

// verifyDockerComponents verifies Docker-related components
func (s *SetupEngine) verifyDockerComponents(projectPath string, result *models.ValidationResult) {
	dockerFiles := []string{"Dockerfile", "docker-compose.yml", "docker-compose.yaml"}

	for _, file := range dockerFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			s.validateDockerFile(projectPath, file, result)
		}
	}
}

// validateDockerFile validates a specific Docker file
func (s *SetupEngine) validateDockerFile(projectPath, file string, result *models.ValidationResult) {
	if file == "Dockerfile" {
		s.validateDockerfile(projectPath, result)
	} else if strings.HasPrefix(file, "docker-compose") {
		s.validateDockerCompose(projectPath, file, result)
	}
}

// validateDockerfile attempts to build Docker image to validate Dockerfile
func (s *SetupEngine) validateDockerfile(projectPath string, result *models.ValidationResult) {
	if err := s.runCommand(projectPath, "docker", "build", "-t", "test-image", "."); err != nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "DockerVerification",
			Message: fmt.Sprintf("Docker build failed: %s", err.Error()),
		})
	}
}

// validateDockerCompose validates docker-compose file syntax
func (s *SetupEngine) validateDockerCompose(projectPath, file string, result *models.ValidationResult) {
	if err := s.runCommand(projectPath, "docker-compose", "-f", file, "config"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "DockerVerification",
			Tag:     "validation",
			Value:   file,
			Message: fmt.Sprintf("Docker Compose validation failed: %s", err.Error()),
		})
	}
}

// verifyKubernetesComponents verifies Kubernetes manifests
func (s *SetupEngine) verifyKubernetesComponents(projectPath string, result *models.ValidationResult) {
	k8sPath := filepath.Join(projectPath, "k8s")
	if _, err := os.Stat(k8sPath); err != nil {
		return // k8s directory doesn't exist
	}

	entries, err := os.ReadDir(k8sPath)
	if err != nil {
		return // Can't read directory
	}

	s.validateKubernetesManifests(k8sPath, entries, result)
}

// validateKubernetesManifests validates individual Kubernetes manifest files
func (s *SetupEngine) validateKubernetesManifests(k8sPath string, entries []os.DirEntry, result *models.ValidationResult) {
	for _, entry := range entries {
		if s.isKubernetesManifest(entry) {
			s.validateSingleK8sManifest(k8sPath, entry.Name(), result)
		}
	}
}

// isKubernetesManifest checks if a file is a Kubernetes manifest
func (s *SetupEngine) isKubernetesManifest(entry os.DirEntry) bool {
	return !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml"))
}

// validateSingleK8sManifest validates a single Kubernetes manifest file
func (s *SetupEngine) validateSingleK8sManifest(k8sPath, filename string, result *models.ValidationResult) {
	if err := s.runCommand(k8sPath, "kubectl", "apply", "--dry-run=client", "-f", filename); err != nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "KubernetesVerification",
			Message: fmt.Sprintf("Kubernetes manifest validation failed for %s: %s", filename, err.Error()),
		})
	}
}

// runCommand executes a command with timeout
func (s *SetupEngine) runCommand(workDir string, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workDir

	// Capture output for error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("command timed out after %v: %s %v", s.timeout, name, args)
		}
		return fmt.Errorf("command failed: %s %v\nOutput: %s", name, args, string(output))
	}

	return nil
}

// SetTimeout sets the timeout for command execution
func (s *SetupEngine) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}
