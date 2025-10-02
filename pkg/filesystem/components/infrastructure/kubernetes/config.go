package kubernetes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FileSystemOperations interface for file operations
type FileSystemOperations interface {
	WriteFile(path string, content []byte, perm os.FileMode) error
	EnsureDirectory(path string) error
}

// ConfigGenerator handles Kubernetes configuration generation
type ConfigGenerator struct {
	fsOps       FileSystemOperations
	manifestGen *ManifestGenerator
}

// NewConfigGenerator creates a new Kubernetes config generator
func NewConfigGenerator(fsOps FileSystemOperations) *ConfigGenerator {
	return &ConfigGenerator{
		fsOps:       fsOps,
		manifestGen: NewManifestGenerator(),
	}
}

// GenerateKubernetesFiles creates all Kubernetes configuration files
func (cg *ConfigGenerator) GenerateKubernetesFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate namespace.yaml
	namespaceContent := cg.manifestGen.GenerateNamespace(config)
	namespacePath := filepath.Join(projectPath, "Deploy/kubernetes/base/namespace.yaml")
	if err := cg.fsOps.WriteFile(namespacePath, []byte(namespaceContent), 0644); err != nil {
		return fmt.Errorf("failed to create namespace.yaml: %w", err)
	}

	// Generate configmap.yaml
	configMapContent := cg.manifestGen.GenerateConfigMap(config)
	configMapPath := filepath.Join(projectPath, "Deploy/kubernetes/base/configmap.yaml")
	if err := cg.fsOps.WriteFile(configMapPath, []byte(configMapContent), 0644); err != nil {
		return fmt.Errorf("failed to create configmap.yaml: %w", err)
	}

	// Generate secret.yaml
	secretContent := cg.manifestGen.GenerateSecret(config)
	secretPath := filepath.Join(projectPath, "Deploy/kubernetes/base/secret.yaml")
	if err := cg.fsOps.WriteFile(secretPath, []byte(secretContent), 0644); err != nil {
		return fmt.Errorf("failed to create secret.yaml: %w", err)
	}

	// Generate backend deployment
	backendDeploymentContent := cg.manifestGen.GenerateBackendDeployment(config)
	backendDeploymentPath := filepath.Join(projectPath, "Deploy/kubernetes/base/backend-deployment.yaml")
	if err := cg.fsOps.WriteFile(backendDeploymentPath, []byte(backendDeploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to create backend-deployment.yaml: %w", err)
	}

	// Generate frontend deployment
	frontendDeploymentContent := cg.manifestGen.GenerateFrontendDeployment(config)
	frontendDeploymentPath := filepath.Join(projectPath, "Deploy/kubernetes/base/frontend-deployment.yaml")
	if err := cg.fsOps.WriteFile(frontendDeploymentPath, []byte(frontendDeploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to create frontend-deployment.yaml: %w", err)
	}

	// Generate services
	servicesContent := cg.manifestGen.GenerateServices(config)
	servicesPath := filepath.Join(projectPath, "Deploy/kubernetes/base/services.yaml")
	if err := cg.fsOps.WriteFile(servicesPath, []byte(servicesContent), 0644); err != nil {
		return fmt.Errorf("failed to create services.yaml: %w", err)
	}

	// Generate ingress
	ingressContent := cg.manifestGen.GenerateIngress(config)
	ingressPath := filepath.Join(projectPath, "Deploy/kubernetes/base/ingress.yaml")
	if err := cg.fsOps.WriteFile(ingressPath, []byte(ingressContent), 0644); err != nil {
		return fmt.Errorf("failed to create ingress.yaml: %w", err)
	}

	return nil
}
