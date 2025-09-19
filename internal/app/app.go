// Package app provides the core application logic for the Open Source Project Generator.
//
// This package implements the main application structure, CLI command handling,
// and orchestrates the interaction between different components like template
// processing, validation, and project generation.
//
// The application follows clean architecture principles with dependency injection
// to ensure testability and maintainability.
package app

import (
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/audit"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

// App represents the main application instance that orchestrates all CLI operations.
// It manages all components, CLI interface, and comprehensive functionality.
//
// The App struct serves as the central coordinator for:
//   - CLI command processing and routing
//   - Component initialization and dependency injection
//   - Project generation workflows with advanced options
//   - Comprehensive validation and auditing operations
//   - Configuration management with multiple sources
//   - Template management and processing
//   - Cache management and offline mode support
//   - Version management and update checking
//   - Security management and validation
//   - Logging and debugging capabilities
type App struct {
	// Core managers and engines
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	templateManager interfaces.TemplateManager
	cacheManager    interfaces.CacheManager
	versionManager  interfaces.VersionManager
	auditEngine     interfaces.AuditEngine
	securityManager interfaces.SecurityManager

	// CLI and generation components
	cli            interfaces.CLIInterface
	generator      interfaces.FileSystemGenerator
	templateEngine interfaces.TemplateEngine
	logger         interfaces.Logger

	// Version information
	version   string
	gitCommit string
	buildTime string
}

// NewApp creates a new application instance with all required dependencies.
//
// Parameters:
//   - appVersion: Application version string
//   - gitCommit: Git commit hash
//   - buildTime: Build timestamp
//
// Returns:
//   - *App: New application instance ready for use
//   - error: Any error that occurred during initialization
func NewApp(appVersion, gitCommit, buildTime string) (*App, error) {
	// Initialize logger first
	logger, err := NewLogger(LogLevelInfo, true)
	if err != nil {
		return nil, err
	}

	// Initialize workspace directory for security manager
	workspaceDir, err := os.Getwd()
	if err != nil {
		workspaceDir = "."
	}

	// Initialize cache manager with default cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	cacheDir := filepath.Join(homeDir, ".generator", "cache")
	cacheManager := cache.NewManager(cacheDir)

	// Initialize all core managers and engines
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	generator := filesystem.NewGenerator()
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManagerWithVersionAndCache(appVersion, cacheManager)
	auditEngine := audit.NewEngine()
	securityManager := security.NewSecurityManager(workspaceDir)

	// Initialize CLI with all dependencies
	cli := cli.NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		logger,
		appVersion,
	)

	return &App{
		// Core managers and engines
		configManager:   configManager,
		validator:       validator,
		templateManager: templateManager,
		cacheManager:    cacheManager,
		versionManager:  versionManager,
		auditEngine:     auditEngine,
		securityManager: securityManager,

		// CLI and generation components
		cli:            cli,
		generator:      generator,
		templateEngine: templateEngine,
		logger:         logger,

		// Version information
		version:   appVersion,
		gitCommit: gitCommit,
		buildTime: buildTime,
	}, nil
}

// Run starts the application and processes command-line arguments.
//
// Parameters:
//   - args: Command-line arguments (typically os.Args[1:])
//
// Returns:
//   - error: Any error that occurred during execution
func (a *App) Run(args []string) error {
	return a.cli.Run(args)
}

// GetConfigManager returns the configuration manager instance
func (a *App) GetConfigManager() interfaces.ConfigManager {
	return a.configManager
}

// GetValidator returns the validation engine instance
func (a *App) GetValidator() interfaces.ValidationEngine {
	return a.validator
}

// GetTemplateManager returns the template manager instance
func (a *App) GetTemplateManager() interfaces.TemplateManager {
	return a.templateManager
}

// GetCacheManager returns the cache manager instance
func (a *App) GetCacheManager() interfaces.CacheManager {
	return a.cacheManager
}

// GetVersionManager returns the version manager instance
func (a *App) GetVersionManager() interfaces.VersionManager {
	return a.versionManager
}

// GetAuditEngine returns the audit engine instance
func (a *App) GetAuditEngine() interfaces.AuditEngine {
	return a.auditEngine
}

// GetSecurityManager returns the security manager instance
func (a *App) GetSecurityManager() interfaces.SecurityManager {
	return a.securityManager
}

// GetCLI returns the CLI interface instance
func (a *App) GetCLI() interfaces.CLIInterface {
	return a.cli
}

// GetGenerator returns the filesystem generator instance
func (a *App) GetGenerator() interfaces.FileSystemGenerator {
	return a.generator
}

// GetTemplateEngine returns the template engine instance
func (a *App) GetTemplateEngine() interfaces.TemplateEngine {
	return a.templateEngine
}

// GetLogger returns the logger instance
func (a *App) GetLogger() interfaces.Logger {
	return a.logger
}

// GetVersion returns the application version information
func (a *App) GetVersion() (version, gitCommit, buildTime string) {
	return a.version, a.gitCommit, a.buildTime
}
