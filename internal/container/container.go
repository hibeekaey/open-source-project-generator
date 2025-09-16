package container

import (
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Container manages basic dependency injection for the application
type Container struct {
	cli            interfaces.CLIInterface
	templateEngine interfaces.TemplateEngine
	configManager  interfaces.ConfigManager
	fsGenerator    interfaces.FileSystemGenerator
	versionManager interfaces.VersionManager
	validator      interfaces.ValidationEngine
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{}
}

// SetCLI sets the CLI interface implementation
func (c *Container) SetCLI(cli interfaces.CLIInterface) {
	c.cli = cli
}

// GetCLI returns the CLI interface implementation
func (c *Container) GetCLI() interfaces.CLIInterface {
	return c.cli
}

// SetTemplateEngine sets the template engine implementation
func (c *Container) SetTemplateEngine(engine interfaces.TemplateEngine) {
	c.templateEngine = engine
}

// GetTemplateEngine returns the template engine implementation
func (c *Container) GetTemplateEngine() interfaces.TemplateEngine {
	return c.templateEngine
}

// SetConfigManager sets the configuration manager implementation
func (c *Container) SetConfigManager(manager interfaces.ConfigManager) {
	c.configManager = manager
}

// GetConfigManager returns the configuration manager implementation
func (c *Container) GetConfigManager() interfaces.ConfigManager {
	return c.configManager
}

// SetFileSystemGenerator sets the file system generator implementation
func (c *Container) SetFileSystemGenerator(generator interfaces.FileSystemGenerator) {
	c.fsGenerator = generator
}

// GetFileSystemGenerator returns the file system generator implementation
func (c *Container) GetFileSystemGenerator() interfaces.FileSystemGenerator {
	return c.fsGenerator
}

// SetVersionManager sets the version manager implementation
func (c *Container) SetVersionManager(manager interfaces.VersionManager) {
	c.versionManager = manager
}

// GetVersionManager returns the version manager implementation
func (c *Container) GetVersionManager() interfaces.VersionManager {
	return c.versionManager
}

// SetValidator sets the validation engine implementation
func (c *Container) SetValidator(validator interfaces.ValidationEngine) {
	c.validator = validator
}

// GetValidator returns the validation engine implementation
func (c *Container) GetValidator() interfaces.ValidationEngine {
	return c.validator
}
