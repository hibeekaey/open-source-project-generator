package orchestrator

import (
	"fmt"
	"sync"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// ExecutorRegistry manages bootstrap executors for different component types
type ExecutorRegistry struct {
	executors map[string]interfaces.BootstrapExecutorInterface
	mu        sync.RWMutex
	logger    *logger.Logger
}

// NewExecutorRegistry creates a new executor registry
func NewExecutorRegistry(log *logger.Logger) *ExecutorRegistry {
	return &ExecutorRegistry{
		executors: make(map[string]interfaces.BootstrapExecutorInterface),
		logger:    log,
	}
}

// Register adds an executor for a component type
func (er *ExecutorRegistry) Register(componentType string, executor interfaces.BootstrapExecutorInterface) {
	er.mu.Lock()
	defer er.mu.Unlock()

	er.executors[componentType] = executor
	er.logger.Debug(fmt.Sprintf("Registered executor for component type: %s", componentType))
}

// Get retrieves an executor for a component type
func (er *ExecutorRegistry) Get(componentType string) (interfaces.BootstrapExecutorInterface, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	executor, exists := er.executors[componentType]
	if !exists {
		return nil, fmt.Errorf("no executor registered for component type: %s", componentType)
	}

	return executor, nil
}

// Has checks if an executor is registered for a component type
func (er *ExecutorRegistry) Has(componentType string) bool {
	er.mu.RLock()
	defer er.mu.RUnlock()

	_, exists := er.executors[componentType]
	return exists
}

// GetDefaultFlags returns default flags for a component type
func (er *ExecutorRegistry) GetDefaultFlags(componentType string) []string {
	er.mu.RLock()
	defer er.mu.RUnlock()

	executor, exists := er.executors[componentType]
	if !exists {
		return []string{}
	}

	return executor.GetDefaultFlags(componentType)
}

// GetSupportedTypes returns all registered component types
func (er *ExecutorRegistry) GetSupportedTypes() []string {
	er.mu.RLock()
	defer er.mu.RUnlock()

	types := make([]string, 0, len(er.executors))
	for componentType := range er.executors {
		types = append(types, componentType)
	}

	return types
}
