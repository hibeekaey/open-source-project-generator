// Package container provides a comprehensive dependency injection system for the application.
//
// The container manages service registration, creation, and lifecycle. It supports:
// - Service factory registration with automatic type checking
// - Singleton service caching for performance
// - Thread-safe service access with proper locking
// - Version information management for service creation
//
// Example usage:
//
//	container := NewContainer()
//	container.SetVersionInfo("1.0.0", "abc123", "2023-01-01")
//
//	// Register a service
//	container.RegisterService("myService", func() (interface{}, error) {
//		return &MyService{}, nil
//	})
//
//	// Get the service (cached after first creation)
//	service, err := container.GetService("myService")
//	if err != nil {
//		log.Fatal(err)
//	}
package container

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

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

// Container manages dependency injection for the application.
// It provides thread-safe service registration, creation, and caching.
// All services are treated as singletons and cached after first creation.
type Container struct {
	mu sync.RWMutex // Protects all fields below

	// factories stores service factory functions keyed by service name
	factories map[string]ServiceFactory

	// serviceCache stores created service instances for singleton behavior
	serviceCache map[string]interface{}

	// initialized tracks which services have been successfully initialized
	initialized map[string]bool

	// versionInfo stores application version information used by services
	versionInfo *VersionInfo
}

// ServiceFactory defines a factory function for creating services.
// The factory should return a service instance or an error if creation fails.
// Factories are called lazily when a service is first requested and the result
// is cached for subsequent requests (singleton pattern).
type ServiceFactory func() (interface{}, error)

// NewContainer creates a new dependency injection container with empty state.
// The container is ready to use immediately and is thread-safe.
// Services must be registered before they can be retrieved.
func NewContainer() *Container {
	return &Container{
		factories:    make(map[string]ServiceFactory),
		serviceCache: make(map[string]interface{}),
		initialized:  make(map[string]bool),
	}
}

// RegisterService registers a service factory with the container.
// The factory will be called lazily when the service is first requested.
// Service names must be unique and non-empty. The factory function cannot be nil.
// This method is thread-safe and can be called concurrently.
func (c *Container) RegisterService(name string, factory ServiceFactory) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "" {
		return fmt.Errorf("container: service name cannot be empty")
	}

	if factory == nil {
		return fmt.Errorf("container: factory cannot be nil for service '%s'", name)
	}

	c.factories[name] = factory
	return nil
}

// GetService retrieves a service from the container by name.
// If the service has not been created yet, its factory function is called.
// The created service is cached for future requests (singleton pattern).
// This method is thread-safe and can be called concurrently.
// Returns an error if the service is not registered or if creation fails.
func (c *Container) GetService(name string) (interface{}, error) {
	if name == "" {
		return nil, fmt.Errorf("container: service name cannot be empty")
	}

	// Check cache first for singleton services
	c.mu.RLock()
	if cached, exists := c.serviceCache[name]; exists {
		c.mu.RUnlock()
		return cached, nil
	}

	factory, exists := c.factories[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("container: service '%s' not registered", name)
	}

	service, err := factory()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create service '%s': %w", name, err)
	}

	// Cache the service for future use (all services are treated as singletons)
	c.mu.Lock()
	c.serviceCache[name] = service
	c.mu.Unlock()

	return service, nil
}

// RegisterAllServices registers all required services with their implementations
func (c *Container) RegisterAllServices() error {
	// Register CLI handlers
	if err := c.registerCLIHandlers(); err != nil {
		return fmt.Errorf("failed to register CLI handlers: %w", err)
	}

	// Register core services
	if err := c.registerCoreServices(); err != nil {
		return fmt.Errorf("failed to register core services: %w", err)
	}

	// Register infrastructure services
	if err := c.registerInfrastructureServices(); err != nil {
		return fmt.Errorf("failed to register infrastructure services: %w", err)
	}

	return nil
}

// registerCLIHandlers registers all CLI handlers in the dependency container
func (c *Container) registerCLIHandlers() error {
	// Register CLI interface
	if err := c.RegisterService("cli", func() (interface{}, error) {
		return c.createCLI()
	}); err != nil {
		return err
	}

	return nil
}

// registerCoreServices registers core services (generator, auditor, cache manager, template manager)
func (c *Container) registerCoreServices() error {
	// Register template engine
	if err := c.RegisterService("templateEngine", func() (interface{}, error) {
		return c.createTemplateEngine()
	}); err != nil {
		return err
	}

	// Register template manager
	if err := c.RegisterService("templateManager", func() (interface{}, error) {
		return c.createTemplateManager()
	}); err != nil {
		return err
	}

	// Register configuration manager
	if err := c.RegisterService("configManager", func() (interface{}, error) {
		return c.createConfigManager()
	}); err != nil {
		return err
	}

	// Register validation engine
	if err := c.RegisterService("validator", func() (interface{}, error) {
		return c.createValidator()
	}); err != nil {
		return err
	}

	// Register audit engine
	if err := c.RegisterService("auditEngine", func() (interface{}, error) {
		return c.createAuditEngine()
	}); err != nil {
		return err
	}

	// Register cache manager
	if err := c.RegisterService("cacheManager", func() (interface{}, error) {
		return c.createCacheManager()
	}); err != nil {
		return err
	}

	return nil
}

// registerInfrastructureServices registers infrastructure services (filesystem, security, version manager)
func (c *Container) registerInfrastructureServices() error {
	// Register filesystem generator
	if err := c.RegisterService("fsGenerator", func() (interface{}, error) {
		return c.createFileSystemGenerator()
	}); err != nil {
		return err
	}

	// Register security manager
	if err := c.RegisterService("securityManager", func() (interface{}, error) {
		return c.createSecurityManager()
	}); err != nil {
		return err
	}

	// Register version manager
	if err := c.RegisterService("versionManager", func() (interface{}, error) {
		return c.createVersionManager()
	}); err != nil {
		return err
	}

	// Register logger
	if err := c.RegisterService("logger", func() (interface{}, error) {
		return c.createLogger()
	}); err != nil {
		return err
	}

	// Register interactive UI
	if err := c.RegisterService("interactiveUI", func() (interface{}, error) {
		return c.createInteractiveUI()
	}); err != nil {
		return err
	}

	return nil
}

// getTypedService is a generic helper to get a service with type checking
func (c *Container) getTypedService(serviceName string, expectedType string) (interface{}, error) {
	service, err := c.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("container: failed to get %s service: %w", serviceName, err)
	}
	return service, nil
}

// Service getter methods with proper error handling
func (c *Container) GetCLI() (interfaces.CLIInterface, error) {
	service, err := c.getTypedService("cli", "CLIInterface")
	if err != nil {
		return nil, err
	}
	cli, ok := service.(interfaces.CLIInterface)
	if !ok {
		return nil, fmt.Errorf("container: service 'cli' has incorrect type, expected CLIInterface but got %T", service)
	}
	return cli, nil
}

func (c *Container) GetTemplateEngine() (interfaces.TemplateEngine, error) {
	service, err := c.getTypedService("templateEngine", "TemplateEngine")
	if err != nil {
		return nil, err
	}
	engine, ok := service.(interfaces.TemplateEngine)
	if !ok {
		return nil, fmt.Errorf("container: service 'templateEngine' has incorrect type, expected TemplateEngine but got %T", service)
	}
	return engine, nil
}

func (c *Container) GetTemplateManager() (interfaces.TemplateManager, error) {
	service, err := c.getTypedService("templateManager", "TemplateManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.TemplateManager)
	if !ok {
		return nil, fmt.Errorf("container: service 'templateManager' has incorrect type, expected TemplateManager but got %T", service)
	}
	return manager, nil
}

func (c *Container) GetConfigManager() (interfaces.ConfigManager, error) {
	service, err := c.getTypedService("configManager", "ConfigManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.ConfigManager)
	if !ok {
		return nil, fmt.Errorf("container: service 'configManager' has incorrect type, expected ConfigManager but got %T", service)
	}
	return manager, nil
}

func (c *Container) GetFileSystemGenerator() (interfaces.FileSystemGenerator, error) {
	service, err := c.getTypedService("fsGenerator", "FileSystemGenerator")
	if err != nil {
		return nil, err
	}
	generator, ok := service.(interfaces.FileSystemGenerator)
	if !ok {
		return nil, fmt.Errorf("container: service 'fsGenerator' has incorrect type, expected FileSystemGenerator but got %T", service)
	}
	return generator, nil
}

func (c *Container) GetVersionManager() (interfaces.VersionManager, error) {
	service, err := c.getTypedService("versionManager", "VersionManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.VersionManager)
	if !ok {
		return nil, fmt.Errorf("container: service 'versionManager' has incorrect type, expected VersionManager but got %T", service)
	}
	return manager, nil
}

func (c *Container) GetValidator() (interfaces.ValidationEngine, error) {
	service, err := c.getTypedService("validator", "ValidationEngine")
	if err != nil {
		return nil, err
	}
	validator, ok := service.(interfaces.ValidationEngine)
	if !ok {
		return nil, fmt.Errorf("container: service 'validator' has incorrect type, expected ValidationEngine but got %T", service)
	}
	return validator, nil
}

func (c *Container) GetAuditEngine() (interfaces.AuditEngine, error) {
	service, err := c.getTypedService("auditEngine", "AuditEngine")
	if err != nil {
		return nil, err
	}
	engine, ok := service.(interfaces.AuditEngine)
	if !ok {
		return nil, fmt.Errorf("container: service 'auditEngine' has incorrect type, expected AuditEngine but got %T", service)
	}
	return engine, nil
}

func (c *Container) GetCacheManager() (interfaces.CacheManager, error) {
	service, err := c.getTypedService("cacheManager", "CacheManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.CacheManager)
	if !ok {
		return nil, fmt.Errorf("container: service 'cacheManager' has incorrect type, expected CacheManager but got %T", service)
	}
	return manager, nil
}

func (c *Container) GetSecurityManager() (interfaces.SecurityManager, error) {
	service, err := c.getTypedService("securityManager", "SecurityManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.SecurityManager)
	if !ok {
		return nil, fmt.Errorf("container: service 'securityManager' has incorrect type, expected SecurityManager but got %T", service)
	}
	return manager, nil
}

func (c *Container) GetLogger() (interfaces.Logger, error) {
	service, err := c.getTypedService("logger", "Logger")
	if err != nil {
		return nil, err
	}
	logger, ok := service.(interfaces.Logger)
	if !ok {
		return nil, fmt.Errorf("container: service 'logger' has incorrect type, expected Logger but got %T", service)
	}
	return logger, nil
}

func (c *Container) GetInteractiveUI() (interfaces.InteractiveUIInterface, error) {
	service, err := c.getTypedService("interactiveUI", "InteractiveUIInterface")
	if err != nil {
		return nil, err
	}
	ui, ok := service.(interfaces.InteractiveUIInterface)
	if !ok {
		return nil, fmt.Errorf("container: service 'interactiveUI' has incorrect type, expected InteractiveUIInterface but got %T", service)
	}
	return ui, nil
}

// Service factory functions - these will create actual implementations
// For now, these return nil implementations that need to be replaced with actual implementations

func (c *Container) createCLI() (interfaces.CLIInterface, error) {
	// Reuse existing services from the container to avoid duplicate creation
	configManager, err := c.createConfigManager()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create config manager for CLI: %w", err)
	}

	validator, err := c.createValidator()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create validator for CLI: %w", err)
	}

	auditEngine, err := c.createAuditEngine()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create audit engine for CLI: %w", err)
	}

	cacheManager, err := c.createCacheManager()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create cache manager for CLI: %w", err)
	}

	templateManager, err := c.createTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create template manager for CLI: %w", err)
	}

	versionManager, err := c.createVersionManager()
	if err != nil {
		return nil, fmt.Errorf("container: failed to create version manager for CLI: %w", err)
	}

	// Create a simple logger (placeholder)
	logger := &SimpleLogger{}

	// Get version info
	versionInfo := c.GetVersionInfo()
	if versionInfo == nil {
		return nil, fmt.Errorf("container: version information not set, cannot create CLI service")
	}

	// Create CLI with all dependencies
	return cli.NewCLI(
		configManager,
		validator,
		templateManager,
		cacheManager,
		versionManager,
		auditEngine,
		logger,
		versionInfo.Version,
		versionInfo.GitCommit,
		versionInfo.BuildTime,
	), nil
}

func (c *Container) createTemplateEngine() (interfaces.TemplateEngine, error) {
	return template.NewEmbeddedEngine(), nil
}

func (c *Container) createTemplateManager() (interfaces.TemplateManager, error) {
	templateEngine, err := c.createTemplateEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create template engine: %w", err)
	}
	return template.NewManager(templateEngine), nil
}

func (c *Container) createConfigManager() (interfaces.ConfigManager, error) {
	return config.NewManager("", ""), nil
}

func (c *Container) createValidator() (interfaces.ValidationEngine, error) {
	return validation.NewEngine(), nil
}

func (c *Container) createAuditEngine() (interfaces.AuditEngine, error) {
	return audit.NewEngine(), nil
}

func (c *Container) createCacheManager() (interfaces.CacheManager, error) {
	// Initialize cache manager with default cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	cacheDir := filepath.Join(homeDir, ".generator", "cache")
	return cache.NewManager(cacheDir), nil
}

func (c *Container) createFileSystemGenerator() (interfaces.FileSystemGenerator, error) {
	return filesystem.NewGenerator(), nil
}

func (c *Container) createSecurityManager() (interfaces.SecurityManager, error) {
	// Initialize workspace directory for security manager
	workspaceDir, err := os.Getwd()
	if err != nil {
		workspaceDir = "."
	}
	return security.NewSecurityManager(workspaceDir), nil
}

func (c *Container) createVersionManager() (interfaces.VersionManager, error) {
	cacheManager, err := c.createCacheManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create cache manager: %w", err)
	}

	versionInfo := c.GetVersionInfo()
	if versionInfo == nil {
		return nil, fmt.Errorf("container: version information not set, cannot create version manager service")
	}

	return version.NewManagerWithVersionAndCache(versionInfo.Version, cacheManager), nil
}

func (c *Container) createLogger() (interfaces.Logger, error) {
	// Import the app package for logger - this creates a circular dependency
	// For now, we'll return an error and handle this in the next task
	return nil, fmt.Errorf("container: logger implementation requires refactoring to avoid circular dependency")
}

func (c *Container) createInteractiveUI() (interfaces.InteractiveUIInterface, error) {
	// TODO: Replace with actual InteractiveUIInterface implementation
	return nil, fmt.Errorf("container: InteractiveUIInterface implementation not yet available")
}

// Initialize initializes all registered services
func (c *Container) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Register all services first
	if err := c.RegisterAllServices(); err != nil {
		return fmt.Errorf("failed to register services: %w", err)
	}

	// Mark container as initialized
	c.initialized["container"] = true
	return nil
}

// IsInitialized checks if the container is initialized
func (c *Container) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized["container"]
}

// GetRegisteredServices returns a list of all registered service names
func (c *Container) GetRegisteredServices() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Pre-allocate slice with exact capacity to avoid reallocations
	services := make([]string, len(c.factories))
	i := 0
	for name := range c.factories {
		services[i] = name
		i++
	}
	return services
}

// ServiceExists checks if a service is registered
func (c *Container) ServiceExists(name string) bool {
	c.mu.RLock()
	_, exists := c.factories[name]
	c.mu.RUnlock()
	return exists
}

// VersionInfo stores application version information that can be used by services.
// This information is typically set during application startup and used by
// services that need to know the application version, git commit, or build time.
type VersionInfo struct {
	Version   string // Semantic version (e.g., "1.0.0")
	GitCommit string // Git commit hash (e.g., "abc123def")
	BuildTime string // Build timestamp (e.g., "2023-01-01T12:00:00Z")
}

// SetVersionInfo stores version information for use in service factories.
// This information is typically set during application startup and can be
// accessed by service factories that need version information.
// This method is thread-safe and can be called concurrently.
func (c *Container) SetVersionInfo(version, gitCommit, buildTime string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.versionInfo == nil {
		c.versionInfo = &VersionInfo{}
	}

	c.versionInfo.Version = version
	c.versionInfo.GitCommit = gitCommit
	c.versionInfo.BuildTime = buildTime
}

// GetVersionInfo returns stored version information
func (c *Container) GetVersionInfo() *VersionInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.versionInfo
}

// ClearServiceCache clears the service cache to free memory.
// This forces all services to be recreated on next access.
// Use this method carefully as it may cause service state to be lost.
// This method is thread-safe and can be called concurrently.
func (c *Container) ClearServiceCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.serviceCache = make(map[string]interface{})
}

// GetCacheSize returns the number of cached services.
// This can be useful for monitoring memory usage and debugging.
// This method is thread-safe and can be called concurrently.
func (c *Container) GetCacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.serviceCache)
}

// SimpleLogger is a basic logger implementation to avoid circular dependencies.
// It provides a minimal logging interface that writes to the standard log package.
// This logger is used as a fallback when the full logging system is not available
// or would create circular dependencies during container initialization.
type SimpleLogger struct{}

// logWithLevel is a helper method to reduce duplication in logging methods
func (l *SimpleLogger) logWithLevel(level string, msg string, args ...interface{}) {
	if level == "FATAL" {
		log.Fatalf("[%s] %s %v", level, msg, args)
	} else {
		log.Printf("[%s] %s %v", level, msg, args)
	}
}

// logWithFields is a helper method for structured logging
func (l *SimpleLogger) logWithFields(level string, msg string, fields map[string]interface{}) {
	if level == "FATAL" {
		log.Fatalf("[%s] %s %v", level, msg, fields)
	} else {
		log.Printf("[%s] %s %v", level, msg, fields)
	}
}

// Basic logging methods
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.logWithLevel("DEBUG", msg, args...)
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.logWithLevel("INFO", msg, args...)
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.logWithLevel("WARN", msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.logWithLevel("ERROR", msg, args...)
}

func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.logWithLevel("FATAL", msg, args...)
}

// Structured logging methods
func (l *SimpleLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("DEBUG", msg, fields)
}

func (l *SimpleLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("INFO", msg, fields)
}

func (l *SimpleLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("WARN", msg, fields)
}

func (l *SimpleLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("ERROR", msg, fields)
}

func (l *SimpleLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("FATAL", msg, fields)
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

// SimpleLoggerContext implements LoggerContext for the SimpleLogger.
// It provides structured logging capabilities by maintaining a set of fields
// that are included with each log message. This allows for contextual logging
// without the complexity of a full logging framework.
type SimpleLoggerContext struct {
	logger *SimpleLogger          // The underlying logger instance
	fields map[string]interface{} // Context fields to include with log messages
}

// logWithContext is a helper method to reduce duplication in context logging
func (c *SimpleLoggerContext) logWithContext(level string, msg string, args ...interface{}) {
	log.Printf("[%s] %s %v %v", level, msg, args, c.fields)
}

func (c *SimpleLoggerContext) Debug(msg string, args ...interface{}) {
	c.logWithContext("DEBUG", msg, args...)
}

func (c *SimpleLoggerContext) Info(msg string, args ...interface{}) {
	c.logWithContext("INFO", msg, args...)
}

func (c *SimpleLoggerContext) Warn(msg string, args ...interface{}) {
	c.logWithContext("WARN", msg, args...)
}

func (c *SimpleLoggerContext) Error(msg string, args ...interface{}) {
	c.logWithContext("ERROR", msg, args...)
}

func (c *SimpleLoggerContext) ErrorWithError(msg string, err error) {
	log.Printf("[ERROR] %s: %v %v", msg, err, c.fields)
}
