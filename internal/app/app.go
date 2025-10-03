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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/internal/container"
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
	// Dependency injection container
	container *container.Container

	// CLI interface
	cli interfaces.CLIInterface

	// Direct dependencies (used when container is nil)
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	templateManager interfaces.TemplateManager
	cacheManager    interfaces.CacheManager
	versionManager  interfaces.VersionManager
	auditEngine     interfaces.AuditEngine
	logger          interfaces.Logger
	securityManager interfaces.SecurityManager
	generator       interfaces.FileSystemGenerator
	templateEngine  interfaces.TemplateEngine

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
	// For now, create CLI directly to avoid container deadlock issues
	// This will be improved in the next task when we fix the dependency injection system

	// Create basic dependencies directly
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	auditEngine := audit.NewEngine()

	// Create cache manager
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	cacheDir := filepath.Join(homeDir, ".generator", "cache")
	cacheManager := cache.NewManager(cacheDir)

	// Create template engine and manager
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)

	// Create version manager
	versionManager := version.NewManagerWithVersionAndCache(appVersion, cacheManager)

	// Create a simple logger (unique instance per app)
	logger := &SimpleLogger{id: int(time.Now().UnixNano())}

	// Create security manager
	workingDir, _ := os.Getwd()
	securityManager := security.NewSecurityManager(workingDir)

	// Create filesystem generator
	generator := filesystem.NewGenerator()

	// Use the same template engine we created above
	// templateEngine is already created

	// Create CLI with all dependencies
	cliInstance := cli.NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		logger,
		appVersion,
		gitCommit,
		buildTime,
	)

	return &App{
		container:       nil, // No container for now
		cli:             cliInstance,
		configManager:   configManager,
		validator:       validator,
		templateManager: templateManager,
		cacheManager:    cacheManager,
		versionManager:  versionManager,
		auditEngine:     auditEngine,
		logger:          logger,
		securityManager: securityManager,
		generator:       generator,
		templateEngine:  templateEngine,
		version:         appVersion,
		gitCommit:       gitCommit,
		buildTime:       buildTime,
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
	if a.cli == nil {
		return fmt.Errorf("CLI not initialized")
	}
	return a.cli.Run(args)
}

// GetConfigManager returns the configuration manager instance
func (a *App) GetConfigManager() (interfaces.ConfigManager, error) {
	if a.container != nil {
		return a.container.GetConfigManager()
	}
	return a.configManager, nil
}

// GetValidator returns the validation engine instance
func (a *App) GetValidator() (interfaces.ValidationEngine, error) {
	if a.container != nil {
		return a.container.GetValidator()
	}
	return a.validator, nil
}

// GetTemplateManager returns the template manager instance
func (a *App) GetTemplateManager() (interfaces.TemplateManager, error) {
	if a.container != nil {
		return a.container.GetTemplateManager()
	}
	return a.templateManager, nil
}

// GetCacheManager returns the cache manager instance
func (a *App) GetCacheManager() (interfaces.CacheManager, error) {
	if a.container != nil {
		return a.container.GetCacheManager()
	}
	return a.cacheManager, nil
}

// GetVersionManager returns the version manager instance
func (a *App) GetVersionManager() (interfaces.VersionManager, error) {
	if a.container != nil {
		return a.container.GetVersionManager()
	}
	return a.versionManager, nil
}

// GetAuditEngine returns the audit engine instance
func (a *App) GetAuditEngine() (interfaces.AuditEngine, error) {
	if a.container != nil {
		return a.container.GetAuditEngine()
	}
	return a.auditEngine, nil
}

// GetSecurityManager returns the security manager instance
func (a *App) GetSecurityManager() (interfaces.SecurityManager, error) {
	if a.container != nil {
		return a.container.GetSecurityManager()
	}
	return a.securityManager, nil
}

// GetCLI returns the CLI interface instance
func (a *App) GetCLI() (interfaces.CLIInterface, error) {
	if a.container != nil {
		return a.container.GetCLI()
	}
	return a.cli, nil
}

// GetGenerator returns the filesystem generator instance
func (a *App) GetGenerator() (interfaces.FileSystemGenerator, error) {
	if a.container != nil {
		return a.container.GetFileSystemGenerator()
	}
	return a.generator, nil
}

// GetTemplateEngine returns the template engine instance
func (a *App) GetTemplateEngine() (interfaces.TemplateEngine, error) {
	if a.container != nil {
		return a.container.GetTemplateEngine()
	}
	return a.templateEngine, nil
}

// GetLogger returns the logger instance
func (a *App) GetLogger() (interfaces.Logger, error) {
	if a.container != nil {
		return a.container.GetLogger()
	}
	return a.logger, nil
}

// GetVersion returns the application version information
func (a *App) GetVersion() (version, gitCommit, buildTime string) {
	return a.version, a.gitCommit, a.buildTime
}

// GetSystemHealth returns basic system health information
func (a *App) GetSystemHealth() map[string]interface{} {
	return map[string]interface{}{
		"status":     "healthy",
		"services":   len(a.container.GetRegisteredServices()),
		"cache_size": a.container.GetCacheSize(),
	}
}

// PerformHealthCheck performs a basic health check
func (a *App) PerformHealthCheck() map[string]interface{} {
	return a.GetSystemHealth()
}

// RestartFailedComponents is a no-op since we removed the complex health system
func (a *App) RestartFailedComponents() error {
	// Simple implementation: clear cache to force service recreation
	a.container.ClearServiceCache()
	return nil
}

// GetContainer returns the dependency injection container
func (a *App) GetContainer() *container.Container {
	return a.container
}

// SimpleLogger is a basic logger implementation to avoid circular dependencies
type SimpleLogger struct {
	id int // Unique identifier to ensure different instances
}

// Basic logging methods
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, args)
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] %s %v", msg, args)
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] %s %v", msg, args)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, args)
}

func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	log.Fatalf("[FATAL] %s %v", msg, args)
}

// Structured logging methods
func (l *SimpleLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *SimpleLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (l *SimpleLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *SimpleLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	log.Fatalf("[FATAL] %s %v", msg, fields)
}

// Error logging with error objects
func (l *SimpleLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {
	log.Printf("[ERROR] %s: %v %v", msg, err, fields)
}

// Operation tracking
func (l *SimpleLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	return &interfaces.OperationContext{
		Operation: operation,
		StartTime: time.Now(),
		Fields:    fields,
	}
}

func (l *SimpleLogger) LogOperationStart(operation string, fields map[string]interface{}) {
	log.Printf("[INFO] Starting operation: %s %v", operation, fields)
}

func (l *SimpleLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
	log.Printf("[INFO] Operation completed: %s (duration: %v) %v", operation, duration, fields)
}

func (l *SimpleLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {
	log.Printf("[ERROR] Operation failed: %s: %v %v", operation, err, fields)
}

func (l *SimpleLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
	duration := time.Since(ctx.StartTime)
	log.Printf("[INFO] Operation completed: %s (duration: %v) %v", ctx.Operation, duration, additionalFields)
}

func (l *SimpleLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
	duration := time.Since(ctx.StartTime)
	log.Printf("[ERROR] Operation failed: %s (duration: %v): %v %v", ctx.Operation, duration, err, additionalFields)
}

// Performance logging
func (l *SimpleLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
	log.Printf("[INFO] Performance metrics for %s: %v", operation, metrics)
}

func (l *SimpleLogger) LogMemoryUsage(operation string) {
	log.Printf("[INFO] Memory usage for %s: (not implemented)", operation)
}

// Configuration methods
func (l *SimpleLogger) SetLevel(level int)        {}
func (l *SimpleLogger) GetLevel() int             { return 0 }
func (l *SimpleLogger) SetJSONOutput(enable bool) {}
func (l *SimpleLogger) SetCallerInfo(enable bool) {}
func (l *SimpleLogger) IsDebugEnabled() bool      { return true }
func (l *SimpleLogger) IsInfoEnabled() bool       { return true }

// Context methods
func (l *SimpleLogger) WithComponent(component string) interfaces.Logger {
	return l // Simple implementation, doesn't actually store context
}

func (l *SimpleLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &SimpleLoggerContext{logger: l, fields: fields}
}

// Log management
func (l *SimpleLogger) GetLogDir() string                                { return "" }
func (l *SimpleLogger) GetRecentEntries(limit int) []interfaces.LogEntry { return nil }
func (l *SimpleLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (l *SimpleLogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (l *SimpleLogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }

// Lifecycle
func (l *SimpleLogger) Close() error {
	return nil // Nothing to close for simple logger
}

// SimpleLoggerContext implements LoggerContext
type SimpleLoggerContext struct {
	logger *SimpleLogger
	fields map[string]interface{}
}

func (c *SimpleLoggerContext) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] %s %v %v", msg, args, c.fields)
}

func (c *SimpleLoggerContext) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] %s %v %v", msg, args, c.fields)
}

func (c *SimpleLoggerContext) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] %s %v %v", msg, args, c.fields)
}

func (c *SimpleLoggerContext) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] %s %v %v", msg, args, c.fields)
}

func (c *SimpleLoggerContext) ErrorWithError(msg string, err error) {
	log.Printf("[ERROR] %s: %v %v", msg, err, c.fields)
}
