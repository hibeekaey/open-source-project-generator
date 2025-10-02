package terraform

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

// ConfigGenerator handles Terraform configuration generation
type ConfigGenerator struct {
	fsOps        FileSystemOperations
	mainGen      *MainGenerator
	variablesGen *VariablesGenerator
	modulesGen   *ModulesGenerator
}

// NewConfigGenerator creates a new Terraform config generator
func NewConfigGenerator(fsOps FileSystemOperations) *ConfigGenerator {
	return &ConfigGenerator{
		fsOps:        fsOps,
		mainGen:      NewMainGenerator(),
		variablesGen: NewVariablesGenerator(),
		modulesGen:   NewModulesGenerator(),
	}
}

// GenerateTerraformFiles creates all Terraform configuration files
func (cg *ConfigGenerator) GenerateTerraformFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate main.tf
	mainTfContent := cg.mainGen.GenerateMain(config)
	mainTfPath := filepath.Join(projectPath, "Deploy/terraform/main.tf")
	if err := cg.fsOps.WriteFile(mainTfPath, []byte(mainTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.tf: %w", err)
	}

	// Generate variables.tf
	variablesTfContent := cg.variablesGen.GenerateVariables(config)
	variablesTfPath := filepath.Join(projectPath, "Deploy/terraform/variables.tf")
	if err := cg.fsOps.WriteFile(variablesTfPath, []byte(variablesTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create variables.tf: %w", err)
	}

	// Generate outputs.tf
	outputsTfContent := cg.variablesGen.GenerateOutputs(config)
	outputsTfPath := filepath.Join(projectPath, "Deploy/terraform/outputs.tf")
	if err := cg.fsOps.WriteFile(outputsTfPath, []byte(outputsTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create outputs.tf: %w", err)
	}

	// Generate terraform.tfvars.example
	tfVarsExampleContent := cg.variablesGen.GenerateTfVarsExample(config)
	tfVarsExamplePath := filepath.Join(projectPath, "Deploy/terraform/terraform.tfvars.example")
	if err := cg.fsOps.WriteFile(tfVarsExamplePath, []byte(tfVarsExampleContent), 0644); err != nil {
		return fmt.Errorf("failed to create terraform.tfvars.example: %w", err)
	}

	// Generate modules/vpc/main.tf
	vpcModuleContent := cg.modulesGen.GenerateVPCModule(config)
	vpcModulePath := filepath.Join(projectPath, "Deploy/terraform/modules/vpc/main.tf")
	if err := cg.fsOps.WriteFile(vpcModulePath, []byte(vpcModuleContent), 0644); err != nil {
		return fmt.Errorf("failed to create modules/vpc/main.tf: %w", err)
	}

	return nil
}
