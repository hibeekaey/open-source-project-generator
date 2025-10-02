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

// Container manages comprehensive dependency injection for the application
type Container struct {
	mu sync.RWMutex

	// CLI components
	cli interfaces.CLIInterface

	// Core services
	templateEngine  interfaces.TemplateEngine
	templateManager interfaces.TemplateManager
	configManager   interfaces.ConfigManager
	fsGenerator     interfaces.FileSystemGenerator
	versionManager  interfaces.VersionManager
	validator       interfaces.ValidationEngine
	auditEngine     interfaces.AuditEngine
	cacheManager    interfaces.CacheManager

	// Infrastructure services
	securityManager interfaces.SecurityManager
	logger          interfaces.Logger
	interactiveUI   interfaces.InteractiveUIInterface

	// Service factories
	factories map[string]ServiceFactory

	// Initialization state
	initialized map[string]bool

	// Health monitoring
	componentHealth map[string]*ComponentHealth

	// Version information
	versionInfo *VersionInfo
}

// ServiceFactory defines a factory function for creating services
type ServiceFactory func() (interface{}, error)

// ServiceConfig defines configuration for service registration
type ServiceConfig struct {
	Name         string
	Factory      ServiceFactory
	Singleton    bool
	Dependencies []string
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{
		factories:       make(map[string]ServiceFactory),
		initialized:     make(map[string]bool),
		componentHealth: make(map[string]*ComponentHealth),
	}
}

// RegisterService registers a service with the container
func (c *Container) RegisterService(name string, factory ServiceFactory) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil for service %s", name)
	}

	c.factories[name] = factory
	return nil
}

// GetService retrieves a service from the container
func (c *Container) GetService(name string) (interface{}, error) {
	c.mu.RLock()
	factory, exists := c.factories[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not registered", name)
	}

	return factory()
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

// Service getter methods with proper error handling
func (c *Container) GetCLI() (interfaces.CLIInterface, error) {
	service, err := c.GetService("cli")
	if err != nil {
		return nil, err
	}
	cli, ok := service.(interfaces.CLIInterface)
	if !ok {
		return nil, fmt.Errorf("service cli is not of type CLIInterface")
	}
	return cli, nil
}

func (c *Container) GetTemplateEngine() (interfaces.TemplateEngine, error) {
	service, err := c.GetService("templateEngine")
	if err != nil {
		return nil, err
	}
	engine, ok := service.(interfaces.TemplateEngine)
	if !ok {
		return nil, fmt.Errorf("service templateEngine is not of type TemplateEngine")
	}
	return engine, nil
}

func (c *Container) GetTemplateManager() (interfaces.TemplateManager, error) {
	service, err := c.GetService("templateManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.TemplateManager)
	if !ok {
		return nil, fmt.Errorf("service templateManager is not of type TemplateManager")
	}
	return manager, nil
}

func (c *Container) GetConfigManager() (interfaces.ConfigManager, error) {
	service, err := c.GetService("configManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.ConfigManager)
	if !ok {
		return nil, fmt.Errorf("service configManager is not of type ConfigManager")
	}
	return manager, nil
}

func (c *Container) GetFileSystemGenerator() (interfaces.FileSystemGenerator, error) {
	service, err := c.GetService("fsGenerator")
	if err != nil {
		return nil, err
	}
	generator, ok := service.(interfaces.FileSystemGenerator)
	if !ok {
		return nil, fmt.Errorf("service fsGenerator is not of type FileSystemGenerator")
	}
	return generator, nil
}

func (c *Container) GetVersionManager() (interfaces.VersionManager, error) {
	service, err := c.GetService("versionManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.VersionManager)
	if !ok {
		return nil, fmt.Errorf("service versionManager is not of type VersionManager")
	}
	return manager, nil
}

func (c *Container) GetValidator() (interfaces.ValidationEngine, error) {
	service, err := c.GetService("validator")
	if err != nil {
		return nil, err
	}
	validator, ok := service.(interfaces.ValidationEngine)
	if !ok {
		return nil, fmt.Errorf("service validator is not of type ValidationEngine")
	}
	return validator, nil
}

func (c *Container) GetAuditEngine() (interfaces.AuditEngine, error) {
	service, err := c.GetService("auditEngine")
	if err != nil {
		return nil, err
	}
	engine, ok := service.(interfaces.AuditEngine)
	if !ok {
		return nil, fmt.Errorf("service auditEngine is not of type AuditEngine")
	}
	return engine, nil
}

func (c *Container) GetCacheManager() (interfaces.CacheManager, error) {
	service, err := c.GetService("cacheManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.CacheManager)
	if !ok {
		return nil, fmt.Errorf("service cacheManager is not of type CacheManager")
	}
	return manager, nil
}

func (c *Container) GetSecurityManager() (interfaces.SecurityManager, error) {
	service, err := c.GetService("securityManager")
	if err != nil {
		return nil, err
	}
	manager, ok := service.(interfaces.SecurityManager)
	if !ok {
		return nil, fmt.Errorf("service securityManager is not of type SecurityManager")
	}
	return manager, nil
}

func (c *Container) GetLogger() (interfaces.Logger, error) {
	service, err := c.GetService("logger")
	if err != nil {
		return nil, err
	}
	logger, ok := service.(interfaces.Logger)
	if !ok {
		return nil, fmt.Errorf("service logger is not of type Logger")
	}
	return logger, nil
}

func (c *Container) GetInteractiveUI() (interfaces.InteractiveUIInterface, error) {
	service, err := c.GetService("interactiveUI")
	if err != nil {
		return nil, err
	}
	ui, ok := service.(interfaces.InteractiveUIInterface)
	if !ok {
		return nil, fmt.Errorf("service interactiveUI is not of type InteractiveUIInterface")
	}
	return ui, nil
}

// Service factory functions - these will create actual implementations
// For now, these return nil implementations that need to be replaced with actual implementations

func (c *Container) createCLI() (interfaces.CLIInterface, error) {
	// For now, create a simple CLI without full dependency injection to avoid deadlock
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
	versionInfo := c.GetVersionInfo()
	if versionInfo == nil {
		return nil, fmt.Errorf("version information not set")
	}
	versionManager := version.NewManagerWithVersionAndCache(versionInfo.Version, cacheManager)

	// Create a simple logger (placeholder)
	logger := &SimpleLogger{}

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
		return nil, fmt.Errorf("version information not set")
	}

	return version.NewManagerWithVersionAndCache(versionInfo.Version, cacheManager), nil
}

func (c *Container) createLogger() (interfaces.Logger, error) {
	// Import the app package for logger - this creates a circular dependency
	// For now, we'll return an error and handle this in the next task
	return nil, fmt.Errorf("Logger implementation requires refactoring to avoid circular dependency")
}

func (c *Container) createInteractiveUI() (interfaces.InteractiveUIInterface, error) {
	// TODO: Replace with actual InteractiveUIInterface implementation
	return nil, fmt.Errorf("InteractiveUIInterface implementation not yet available")
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

	services := make([]string, 0, len(c.factories))
	for name := range c.factories {
		services = append(services, name)
	}
	return services
}

// ServiceExists checks if a service is registered
func (c *Container) ServiceExists(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.factories[name]
	return exists
}

// ComponentHealth represents the health status of a component
type ComponentHealth struct {
	Name         string                 `json:"name"`
	Status       string                 `json:"status"` // "healthy", "degraded", "failed", "initializing"
	LastCheck    time.Time              `json:"last_check"`
	Dependencies []string               `json:"dependencies"`
	Errors       []string               `json:"errors,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// SystemHealth represents the overall system health
type SystemHealth struct {
	Overall    string                      `json:"overall"`
	Components map[string]*ComponentHealth `json:"components"`
	Timestamp  time.Time                   `json:"timestamp"`
}

// InitializationOrder defines the order in which services should be initialized
var InitializationOrder = []string{
	"logger",          // First - needed for logging initialization
	"securityManager", // Second - needed for secure operations
	"configManager",   // Third - needed for configuration
	"cacheManager",    // Fourth - needed for caching
	"versionManager",  // Fifth - needed for version checks
	"templateEngine",  // Sixth - low-level template processing
	"templateManager", // Seventh - high-level template management
	"fsGenerator",     // Eighth - filesystem operations
	"validator",       // Ninth - validation capabilities
	"auditEngine",     // Tenth - auditing capabilities
	"interactiveUI",   // Eleventh - UI components
	"cli",             // Last - depends on most other services
}

// ServiceDependencies defines dependencies between services
var ServiceDependencies = map[string][]string{
	"logger":          {},
	"securityManager": {"logger"},
	"configManager":   {"logger", "securityManager"},
	"cacheManager":    {"logger", "configManager"},
	"versionManager":  {"logger", "configManager", "cacheManager"},
	"templateEngine":  {"logger", "securityManager"},
	"templateManager": {"logger", "templateEngine", "cacheManager"},
	"fsGenerator":     {"logger", "securityManager", "templateEngine"},
	"validator":       {"logger", "configManager"},
	"auditEngine":     {"logger", "validator", "securityManager"},
	"interactiveUI":   {"logger"},
	"cli":             {"logger", "configManager", "templateManager", "fsGenerator", "versionManager", "validator", "auditEngine", "cacheManager", "securityManager", "interactiveUI"},
}

// InitializeWithOrder initializes services in the proper dependency order
func (c *Container) InitializeWithOrder() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Register all services first
	if err := c.RegisterAllServices(); err != nil {
		return fmt.Errorf("failed to register services: %w", err)
	}

	// Initialize services in dependency order
	for _, serviceName := range InitializationOrder {
		if err := c.initializeService(serviceName); err != nil {
			return fmt.Errorf("failed to initialize service %s: %w", serviceName, err)
		}
	}

	// Mark container as initialized
	c.initialized["container"] = true
	return nil
}

// initializeService initializes a single service with dependency validation
func (c *Container) initializeService(serviceName string) error {
	// Check if service is already initialized
	if c.initialized[serviceName] {
		return nil
	}

	// Validate dependencies are initialized first
	dependencies := ServiceDependencies[serviceName]
	for _, dep := range dependencies {
		if !c.initialized[dep] {
			return fmt.Errorf("service %s depends on %s which is not initialized", serviceName, dep)
		}
	}

	// Mark as initializing
	c.setComponentStatus(serviceName, "initializing", nil)

	// Try to create the service to validate it can be initialized
	_, err := c.GetService(serviceName)
	if err != nil {
		c.setComponentStatus(serviceName, "failed", []string{err.Error()})

		// Check if graceful degradation is possible
		if c.canDegrade(serviceName) {
			c.setComponentStatus(serviceName, "degraded", []string{err.Error()})
			c.initialized[serviceName] = true // Mark as initialized but degraded
			return nil
		}

		return fmt.Errorf("failed to initialize service %s: %w", serviceName, err)
	}

	// Mark as successfully initialized
	c.setComponentStatus(serviceName, "healthy", nil)
	c.initialized[serviceName] = true
	return nil
}

// setComponentStatus sets the status of a component
func (c *Container) setComponentStatus(name, status string, errors []string) {
	if c.componentHealth == nil {
		c.componentHealth = make(map[string]*ComponentHealth)
	}

	health := &ComponentHealth{
		Name:         name,
		Status:       status,
		LastCheck:    time.Now(),
		Dependencies: ServiceDependencies[name],
		Errors:       errors,
		Metadata:     make(map[string]interface{}),
	}

	c.componentHealth[name] = health
}

// canDegrade checks if a service can operate in degraded mode
func (c *Container) canDegrade(serviceName string) bool {
	// Define which services can operate in degraded mode
	degradableServices := map[string]bool{
		"cacheManager":   true, // Can work without cache
		"versionManager": true, // Can work with cached/default versions
		"auditEngine":    true, // Can skip auditing
		"interactiveUI":  true, // Can fall back to basic prompts
	}

	return degradableServices[serviceName]
}

// GetSystemHealth returns the current health status of all components
func (c *Container) GetSystemHealth() *SystemHealth {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.componentHealth == nil {
		return &SystemHealth{
			Overall:    "unknown",
			Components: make(map[string]*ComponentHealth),
			Timestamp:  time.Now(),
		}
	}

	// Calculate overall health
	overall := "healthy"
	healthyCount := 0
	degradedCount := 0
	failedCount := 0

	for _, health := range c.componentHealth {
		switch health.Status {
		case "healthy":
			healthyCount++
		case "degraded":
			degradedCount++
			if overall == "healthy" {
				overall = "degraded"
			}
		case "failed":
			failedCount++
			overall = "unhealthy"
		}
	}

	return &SystemHealth{
		Overall:    overall,
		Components: c.componentHealth,
		Timestamp:  time.Now(),
	}
}

// ValidateInitialization validates that all critical services are properly initialized
func (c *Container) ValidateInitialization() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	criticalServices := []string{"logger", "configManager", "cli"}

	for _, service := range criticalServices {
		if !c.initialized[service] {
			return fmt.Errorf("critical service %s is not initialized", service)
		}

		if c.componentHealth != nil {
			if health, exists := c.componentHealth[service]; exists {
				if health.Status == "failed" {
					return fmt.Errorf("critical service %s has failed: %v", service, health.Errors)
				}
			}
		}
	}

	return nil
}

// PerformHealthCheck performs a health check on all components
func (c *Container) PerformHealthCheck() *SystemHealth {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update health status for each service
	for serviceName := range c.factories {
		c.checkServiceHealth(serviceName)
	}

	return c.GetSystemHealth()
}

// checkServiceHealth checks the health of a specific service
func (c *Container) checkServiceHealth(serviceName string) {
	var errors []string
	status := "healthy"

	// Try to get the service
	_, err := c.GetService(serviceName)
	if err != nil {
		errors = append(errors, err.Error())
		if c.canDegrade(serviceName) {
			status = "degraded"
		} else {
			status = "failed"
		}
	}

	c.setComponentStatus(serviceName, status, errors)
}

// GetComponentHealth returns the health status of a specific component
func (c *Container) GetComponentHealth(name string) (*ComponentHealth, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.componentHealth == nil {
		return nil, fmt.Errorf("health monitoring not initialized")
	}

	health, exists := c.componentHealth[name]
	if !exists {
		return nil, fmt.Errorf("component %s not found", name)
	}

	return health, nil
}

// RestartFailedComponents attempts to restart failed components
func (c *Container) RestartFailedComponents() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.componentHealth == nil {
		return fmt.Errorf("health monitoring not initialized")
	}

	var restartErrors []string

	for name, health := range c.componentHealth {
		if health.Status == "failed" {
			// Mark as not initialized to force restart
			c.initialized[name] = false

			// Try to reinitialize
			if err := c.initializeService(name); err != nil {
				restartErrors = append(restartErrors, fmt.Sprintf("%s: %v", name, err))
			}
		}
	}

	if len(restartErrors) > 0 {
		return fmt.Errorf("failed to restart components: %v", restartErrors)
	}

	return nil
}

// Add componentHealth field to Container struct
// This needs to be added to the Container struct definition
// Version information storage
type VersionInfo struct {
	Version   string
	GitCommit string
	BuildTime string
}

// SetVersionInfo stores version information for use in service factories
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

// SimpleLogger is a basic logger implementation to avoid circular dependencies
type SimpleLogger struct{}

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
