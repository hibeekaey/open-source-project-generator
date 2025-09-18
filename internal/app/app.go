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
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

// App represents the main application instance that orchestrates all CLI operations.
// It manages the basic components, CLI interface, and error handling.
//
// The App struct serves as the central coordinator for:
//   - CLI command processing and routing
//   - Basic component initialization
//   - Project generation workflows
//   - Basic validation operations
//   - Configuration management
type App struct {
	configManager  interfaces.ConfigManager
	validator      interfaces.ValidationEngine
	cli            interfaces.CLIInterface
	generator      interfaces.FileSystemGenerator
	templateEngine interfaces.TemplateEngine
	versionManager interfaces.VersionManager
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

	// Initialize all components
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	generator := filesystem.NewGenerator()
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)
	versionManager := version.NewManager()
	auditEngine := audit.NewEngine()

	// Initialize cache manager with default cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	cacheDir := filepath.Join(homeDir, ".generator", "cache")
	cacheManager := cache.NewManager(cacheDir)

	// Initialize CLI with all dependencies including logger
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
		configManager:  configManager,
		validator:      validator,
		cli:            cli,
		generator:      generator,
		templateEngine: templateEngine,
		versionManager: versionManager,
		logger:         logger,
		version:        appVersion,
		gitCommit:      gitCommit,
		buildTime:      buildTime,
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
