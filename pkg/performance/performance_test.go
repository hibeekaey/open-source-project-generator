package performance

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

// hasEmbeddedTemplates checks if embedded templates are available for testing
// by attempting to create an embedded template engine and checking if it can access templates
func hasEmbeddedTemplates() bool {
	// Try to create an embedded template engine
	engine := template.NewEmbeddedEngine()
	if engine == nil {
		return false
	}

	// Try to load a known template to verify embedded templates are accessible
	// We'll try to load a backend template that should exist
	templatePath := filepath.Join(constants.TemplateBaseDir, constants.TemplateBackend, "go-gin", "template.yaml")
	_, err := engine.LoadTemplate(templatePath)
	return err == nil
}

// TestLargeProjectGeneration tests performance with large project configurations
func TestLargeProjectGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	// Check if embedded templates are available
	if !hasEmbeddedTemplates() {
		t.Skip("Embedded templates not available, skipping performance test")
	}

	tempDir := t.TempDir()

	// Create large project configuration
	config := createLargeProjectConfig(tempDir)

	// Measure generation time
	start := time.Now()

	// Create embedded template engine and manager
	templateEngine := template.NewEmbeddedEngine()
	templateManager := template.NewManager(templateEngine)

	// Process template
	err := templateManager.ProcessTemplate("go-gin", config, config.OutputPath)
	if err != nil {
		t.Fatalf("Failed to process large project: %v", err)
	}

	duration := time.Since(start)

	// Performance assertions
	maxDuration := 30 * time.Second // Should complete within 30 seconds
	if duration > maxDuration {
		t.Errorf("Large project generation took too long: %v (max: %v)", duration, maxDuration)
	}

	t.Logf("Large project generation completed in %v", duration)

	// Verify output was created
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		t.Error("Expected output directory to be created")
	}
}

// TestConcurrentOperations tests performance under concurrent load
func TestConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tempDir := t.TempDir()

	// Number of concurrent operations
	numOperations := 10

	// Create cache manager for concurrent testing
	cacheManager := cache.NewManager(filepath.Join(tempDir, "cache"))

	// Test concurrent cache operations
	t.Run("concurrent_cache_operations", func(t *testing.T) {
		testConcurrentCacheOperations(t, cacheManager, numOperations)
	})

	// Test concurrent validation operations
	t.Run("concurrent_validation_operations", func(t *testing.T) {
		testConcurrentValidationOperations(t, tempDir, numOperations)
	})

	// Test concurrent version operations
	t.Run("concurrent_version_operations", func(t *testing.T) {
		testConcurrentVersionOperations(t, numOperations)
	})
}

func testConcurrentCacheOperations(t *testing.T, cacheManager interfaces.CacheManager, numOperations int) {
	var wg sync.WaitGroup
	errors := make(chan error, numOperations*3) // 3 operations per goroutine

	start := time.Now()

	// Launch concurrent operations
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			key := fmt.Sprintf("test-key-%d", id)
			value := fmt.Sprintf("test-value-%d", id)

			// Set operation
			if err := cacheManager.Set(key, value, 5*time.Minute); err != nil {
				errors <- fmt.Errorf("set operation failed for %s: %v", key, err)
				return
			}

			// Get operation
			if _, err := cacheManager.Get(key); err != nil {
				errors <- fmt.Errorf("get operation failed for %s: %v", key, err)
				return
			}

			// Delete operation
			if err := cacheManager.Delete(key); err != nil {
				errors <- fmt.Errorf("delete operation failed for %s: %v", key, err)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Performance assertions
	maxDuration := 5 * time.Second
	if duration > maxDuration {
		t.Errorf("Concurrent cache operations took too long: %v (max: %v)", duration, maxDuration)
	}

	t.Logf("Concurrent cache operations (%d) completed in %v", numOperations, duration)
}

func testConcurrentValidationOperations(t *testing.T, tempDir string, numOperations int) {
	// Create test projects for validation
	projects := make([]string, numOperations)
	for i := 0; i < numOperations; i++ {
		projectDir := filepath.Join(tempDir, fmt.Sprintf("project-%d", i))
		createTestProject(t, projectDir)
		projects[i] = projectDir
	}

	var wg sync.WaitGroup
	errors := make(chan error, numOperations)

	validationEngine := validation.NewEngine()

	start := time.Now()

	// Launch concurrent validations
	for i, projectDir := range projects {
		wg.Add(1)
		go func(id int, path string) {
			defer wg.Done()

			result, err := validationEngine.ValidateProject(path)
			if err != nil {
				errors <- fmt.Errorf("validation failed for project %d: %v", id, err)
				return
			}

			if result == nil {
				errors <- fmt.Errorf("validation result is nil for project %d", id)
				return
			}
		}(i, projectDir)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Performance assertions
	maxDuration := 10 * time.Second
	if duration > maxDuration {
		t.Errorf("Concurrent validation operations took too long: %v (max: %v)", duration, maxDuration)
	}

	t.Logf("Concurrent validation operations (%d) completed in %v", numOperations, duration)
}

func testConcurrentVersionOperations(t *testing.T, numOperations int) {
	var wg sync.WaitGroup
	errors := make(chan error, numOperations)

	versionManager := version.NewManager()

	start := time.Now()

	// Launch concurrent version checks
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Get current version (should be fast)
			version := versionManager.GetCurrentVersion()
			if version == "" {
				errors <- fmt.Errorf("version check failed for operation %d", id)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Performance assertions
	maxDuration := 2 * time.Second
	if duration > maxDuration {
		t.Errorf("Concurrent version operations took too long: %v (max: %v)", duration, maxDuration)
	}

	t.Logf("Concurrent version operations (%d) completed in %v", numOperations, duration)
}

// Helper functions

func createLargeProjectConfig(tempDir string) *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:         "large-performance-test-project",
		Organization: "performance-test-org",
		Description:  "A large project configuration for performance testing with many components and dependencies",
		License:      "MIT",
		OutputPath:   filepath.Join(tempDir, "large-project"),
		Components: models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Home:  true,
					Admin: true,
				},
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			},
		},
		Features: []string{
			"authentication", "authorization", "database", "cache", "logging",
			"monitoring", "testing", "documentation", "deployment", "security",
		},
	}
}

func createTestProject(t *testing.T, projectDir string) {
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create package.json
	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"dependencies": {
			"express": "^4.18.0",
			"cors": "^2.8.5"
		}
	}`

	err = os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create main.js
	mainJS := `const express = require('express');
const app = express();

app.get('/', (req, res) => {
	res.json({ message: 'Hello World' });
});

const port = process.env.PORT || 3000;
app.listen(port, () => {
	console.log('Server running on port ' + port);
});
`

	err = os.WriteFile(filepath.Join(projectDir, "main.js"), []byte(mainJS), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.js: %v", err)
	}

	// Create README.md
	readme := "# Test Project\n\nThis is a test project for performance testing.\n\n## Installation\n\n```bash\nnpm install\n```\n\n## Usage\n\n```bash\nnpm start\n```\n"

	err = os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}
}

// Benchmark tests

func BenchmarkCacheOperations(b *testing.B) {
	tempDir := b.TempDir()
	cacheManager := cache.NewManager(filepath.Join(tempDir, "bench-cache"))

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench-key-%d", i)
			value := fmt.Sprintf("bench-value-%d", i)
			err := cacheManager.Set(key, value, 5*time.Minute)
			if err != nil {
				b.Fatalf("Cache set failed: %v", err)
			}
		}
	})

	// Pre-populate for get benchmark
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("get-key-%d", i)
		value := fmt.Sprintf("get-value-%d", i)
		_ = cacheManager.Set(key, value, 5*time.Minute)
	}

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("get-key-%d", i%1000)
			_, err := cacheManager.Get(key)
			if err != nil {
				b.Fatalf("Cache get failed: %v", err)
			}
		}
	})
}

func BenchmarkValidation(b *testing.B) {
	tempDir := b.TempDir()
	validationEngine := validation.NewEngine()

	// Create test project
	projectDir := filepath.Join(tempDir, "bench-project")
	createTestProject(&testing.T{}, projectDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := validationEngine.ValidateProject(projectDir)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}

func BenchmarkVersionManager(b *testing.B) {
	versionManager := version.NewManager()

	b.Run("GetCurrentVersion", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			version := versionManager.GetCurrentVersion()
			if version == "" {
				b.Fatal("Version should not be empty")
			}
		}
	})
}
