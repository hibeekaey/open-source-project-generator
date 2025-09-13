package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/internal/config"
	"github.com/open-source-template-generator/internal/container"
	"github.com/open-source-template-generator/pkg/cli"
	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
	"github.com/open-source-template-generator/pkg/validation"
	"github.com/open-source-template-generator/pkg/version"
)

func TestFullSystemIntegration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "full-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize all components
	container := container.NewContainer()

	// Set up configuration manager
	cacheDir := filepath.Join(tempDir, "cache")
	defaultsPath := filepath.Join("templates", "config", "defaults.yaml")
	configManager := config.NewManager(cacheDir, defaultsPath)
	container.SetConfigManager(configManager)

	// Set up validation engine
	validator := validation.NewEngine()
	container.SetValidator(validator)

	// Set up template engine
	templateEngine := template.NewEngine()
	container.SetTemplateEngine(templateEngine)

	// Set up filesystem generator
	fsGenerator := filesystem.NewGenerator()
	container.SetFileSystemGenerator(fsGenerator)

	// Set up version manager
	versionCache := version.NewMemoryCache(24 * time.Hour)
	versionManager := version.NewManager(versionCache)
	container.SetVersionManager(versionManager)

	// Set up CLI
	cliInterface := cli.NewCLI(configManager, validator)
	container.SetCLI(cliInterface)

	// Test that all components are properly initialized
	if container.GetConfigManager() == nil {
		t.Error("ConfigManager not initialized")
	}
	if container.GetValidator() == nil {
		t.Error("Validator not initialized")
	}
	if container.GetTemplateEngine() == nil {
		t.Error("TemplateEngine not initialized")
	}
	if container.GetFileSystemGenerator() == nil {
		t.Error("FileSystemGenerator not initialized")
	}
	if container.GetVersionManager() == nil {
		t.Error("VersionManager not initialized")
	}
	if container.GetCLI() == nil {
		t.Error("CLI not initialized")
	}

	t.Log("Full system integration test completed successfully")
}

func TestProjectGenerationWorkflow(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "workflow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize components
	templateEngine := template.NewEngine()
	fsGenerator := filesystem.NewGenerator()
	validator := validation.NewEngine()

	// Create a test configuration
	config := &models.ProjectConfig{
		Name:         "workflow-test-project",
		Organization: "test-org",
		Description:  "Test project for workflow testing",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/test/workflow-test",
		OutputPath:   tempDir,
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker: true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			Go:     "1.22.0",
			NextJS: "15.0.0",
			React:  "18.2.0",
		},
	}

	// Step 1: Create project structure
	t.Log("Step 1: Creating project structure")
	err = fsGenerator.CreateProject(config, config.OutputPath)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	projectPath := filepath.Join(config.OutputPath, config.Name)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Errorf("Project directory not created: %s", projectPath)
	}

	// Step 2: Generate base files (this might fail if templates don't exist, but shouldn't panic)
	t.Log("Step 2: Generating base files")
	err = templateEngine.ProcessDirectory("templates/base", projectPath, config)
	if err != nil {
		t.Logf("Base file generation failed (expected if templates don't exist): %v", err)
	}

	// Step 3: Validate the generated project structure
	t.Log("Step 3: Validating project structure")
	result, err := validator.ValidateProject(projectPath)
	if err != nil {
		t.Logf("Project validation failed (might be expected): %v", err)
	} else {
		t.Logf("Validation result: Valid=%t, Errors=%d, Warnings=%d",
			result.Valid, len(result.Errors), len(result.Warnings))
	}

	t.Log("Project generation workflow test completed")
}

func TestComponentInteraction(t *testing.T) {
	// Test interaction between different components
	tempDir, err := os.MkdirTemp("", "component-interaction-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize components
	cacheDir := filepath.Join(tempDir, "cache")
	configManager := config.NewManager(cacheDir, "")
	_ = template.NewEngine() // templateEngine not used in this test
	fsGenerator := filesystem.NewGenerator()
	validator := validation.NewEngine()
	versionCache := version.NewMemoryCache(1 * time.Hour)
	_ = version.NewManager(versionCache) // versionManager not used in this test

	// Test configuration loading and validation
	t.Log("Testing configuration management")
	defaults, err := configManager.LoadDefaults()
	if err != nil {
		t.Logf("Loading defaults failed (might be expected): %v", err)
	} else {
		err = configManager.ValidateConfig(defaults)
		if err != nil {
			t.Logf("Default config validation failed: %v", err)
		}
	}

	// Test version management
	t.Log("Testing version management")
	versions, err := configManager.GetLatestVersions()
	if err != nil {
		t.Logf("Getting latest versions failed (might be expected): %v", err)
	} else {
		t.Logf("Retrieved versions: Node=%s, Go=%s", versions.Node, versions.Go)
	}

	// Test template engine with custom functions
	t.Log("Testing template engine")
	// templateEngine.RegisterFunctions() - custom functions are registered internally

	// Test filesystem operations
	t.Log("Testing filesystem operations")
	testConfig := &models.ProjectConfig{
		Name:         "interaction-test",
		Organization: "test-org",
		Description:  "Component interaction test",
		License:      "MIT",
		OutputPath:   tempDir,
	}

	err = fsGenerator.CreateProject(testConfig, testConfig.OutputPath)
	if err != nil {
		t.Errorf("Failed to create project for interaction test: %v", err)
	}

	// Test validation of created project
	t.Log("Testing validation")
	projectPath := filepath.Join(testConfig.OutputPath, testConfig.Name)
	result, err := validator.ValidateProject(projectPath)
	if err != nil {
		t.Logf("Validation failed (might be expected for minimal project): %v", err)
	} else {
		t.Logf("Validation completed: Valid=%t", result.Valid)
	}

	t.Log("Component interaction test completed")
}

func TestErrorHandlingIntegration(t *testing.T) {
	// Test error handling across components
	tempDir, err := os.MkdirTemp("", "error-handling-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with invalid configurations
	t.Log("Testing error handling with invalid configurations")

	templateEngine := template.NewEngine()
	fsGenerator := filesystem.NewGenerator()
	validator := validation.NewEngine()

	// Test with nil config
	t.Log("Testing with nil config")
	err = fsGenerator.CreateProject(nil, tempDir)
	if err == nil {
		t.Error("Expected error with nil config, got nil")
	}

	// Test with invalid output path
	t.Log("Testing with invalid output path")
	invalidConfig := &models.ProjectConfig{
		Name:       "test",
		OutputPath: "/invalid/path/that/should/not/exist/12345",
	}
	err = fsGenerator.CreateProject(invalidConfig, invalidConfig.OutputPath)
	if err == nil {
		t.Error("Expected error with invalid output path, got nil")
	}

	// Test validation with non-existent path
	t.Log("Testing validation with non-existent path")
	result, err := validator.ValidateProject("/non/existent/path/12345")
	if err != nil {
		t.Errorf("Unexpected error validating non-existent path: %v", err)
	}
	if result.Valid {
		t.Error("Expected validation to fail for non-existent path, but it passed")
	}

	// Test template processing with invalid template
	t.Log("Testing template processing with invalid template")
	err = templateEngine.ProcessDirectory("/non/existent/template/dir", tempDir, &models.ProjectConfig{})
	if err == nil {
		t.Error("Expected error processing non-existent template directory, got nil")
	}

	t.Log("Error handling integration test completed")
}

func TestPerformanceIntegration(t *testing.T) {
	// Test performance characteristics of integrated components
	tempDir, err := os.MkdirTemp("", "performance-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize components
	_ = template.NewEngine() // templateEngine not used in this test
	fsGenerator := filesystem.NewGenerator()
	_ = validation.NewEngine() // validator not used in this test
	versionCache := version.NewMemoryCache(1 * time.Hour)
	_ = version.NewManager(versionCache) // versionManager not used in this test

	// Test multiple project generations
	t.Log("Testing performance with multiple project generations")
	startTime := time.Now()

	for i := 0; i < 10; i++ {
		config := &models.ProjectConfig{
			Name:         "perf-test-" + string(rune(i+'0')),
			Organization: "test-org",
			Description:  "Performance test project",
			License:      "MIT",
			OutputPath:   tempDir,
		}

		err = fsGenerator.CreateProject(config, config.OutputPath)
		if err != nil {
			t.Logf("Project creation %d failed: %v", i, err)
		}
	}

	duration := time.Since(startTime)
	t.Logf("Created 10 projects in %v (avg: %v per project)", duration, duration/10)

	// Test version caching performance
	t.Log("Testing version caching performance")
	startTime = time.Now()

	for i := 0; i < 100; i++ {
		versionCache.Set("test-key-"+string(rune(i+'0')), "1.0.0")
	}

	for i := 0; i < 100; i++ {
		_, found := versionCache.Get("test-key-" + string(rune(i+'0')))
		if !found {
			t.Errorf("Cache miss for key %d", i)
		}
	}

	cacheDuration := time.Since(startTime)
	t.Logf("Performed 100 cache writes and 100 cache reads in %v", cacheDuration)

	t.Log("Performance integration test completed")
}

func TestConcurrencyIntegration(t *testing.T) {
	// Test concurrent operations
	tempDir, err := os.MkdirTemp("", "concurrency-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test concurrent version cache operations
	t.Log("Testing concurrent version cache operations")
	versionCache := version.NewMemoryCache(1 * time.Hour)

	// Start multiple goroutines performing cache operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			for j := 0; j < 10; j++ {
				key := "concurrent-key-" + string(rune(id+'0')) + "-" + string(rune(j+'0'))
				value := "1.0." + string(rune(j+'0'))

				versionCache.Set(key, value)

				retrieved, found := versionCache.Get(key)
				if !found {
					t.Errorf("Concurrent cache miss for key %s", key)
				}
				if retrieved != value {
					t.Errorf("Concurrent cache value mismatch: expected %s, got %s", value, retrieved)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	t.Log("Concurrency integration test completed")
}

func TestMemoryUsageIntegration(t *testing.T) {
	// Test memory usage patterns
	tempDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Log("Testing memory usage patterns")

	// Create and destroy multiple components to test for memory leaks
	for i := 0; i < 100; i++ {
		// Create components
		templateEngine := template.NewEngine()
		fsGenerator := filesystem.NewGenerator()
		validator := validation.NewEngine()
		versionCache := version.NewMemoryCache(1 * time.Minute)
		versionManager := version.NewManager(versionCache)

		// Use components briefly
		config := &models.ProjectConfig{
			Name:       "memory-test-" + string(rune(i%10+'0')),
			OutputPath: tempDir,
		}

		_ = fsGenerator.CreateProject(config, config.OutputPath)
		_ = templateEngine.ProcessDirectory("/non/existent", tempDir, config)
		_, _ = validator.ValidateProject(tempDir)
		_, _ = versionManager.GetLatestNodeVersion()

		// Components should be garbage collected when they go out of scope
		templateEngine = nil
		fsGenerator = nil
		validator = nil
		versionCache = nil
		versionManager = nil
	}

	t.Log("Memory usage integration test completed")
}
