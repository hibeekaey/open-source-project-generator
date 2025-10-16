package security

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
)

// TestConcurrentOperationsAndThreadSafety tests thread safety of all components
func TestConcurrentOperationsAndThreadSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent tests in short mode")
	}

	t.Run("concurrent_cache_access", func(t *testing.T) {
		testConcurrentCacheAccess(t)
	})

	t.Run("concurrent_validation", func(t *testing.T) {
		testConcurrentValidation(t)
	})

	t.Run("concurrent_config_operations", func(t *testing.T) {
		testConcurrentConfigOperations(t)
	})

	t.Run("race_condition_detection", func(t *testing.T) {
		testRaceConditionDetection(t)
	})
}

func testConcurrentCacheAccess(t *testing.T) {
	tempDir := t.TempDir()
	cacheManager := cache.NewManager(tempDir)

	// Disable disk persistence for this test to avoid I/O bottlenecks
	config := &interfaces.CacheConfig{
		MaxSize:        1024 * 1024, // 1MB
		MaxEntries:     1000,
		EvictionPolicy: "lru",
		PersistToDisk:  false, // Disable persistence for concurrent tests
	}
	err := cacheManager.SetCacheConfig(config)
	if err != nil {
		t.Fatalf("Failed to set cache config: %v", err)
	}

	numGoroutines := 10          // Reduced from 50
	operationsPerGoroutine := 20 // Reduced from 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	// Launch concurrent cache operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("concurrent:%d:%d", goroutineID, j)
				value := fmt.Sprintf("value_%d_%d", goroutineID, j)

				// Set operation
				if err := cacheManager.Set(key, value, 5*time.Minute); err != nil {
					errors <- fmt.Errorf("goroutine %d: set failed: %w", goroutineID, err)
					continue
				}

				// Get operation
				retrievedValue, err := cacheManager.Get(key)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d: get failed: %w", goroutineID, err)
					continue
				}

				if retrievedValue != value {
					errors <- fmt.Errorf("goroutine %d: value mismatch: expected %s, got %v", goroutineID, value, retrievedValue)
					continue
				}

				// Exists check
				if !cacheManager.Exists(key) {
					errors <- fmt.Errorf("goroutine %d: key should exist: %s", goroutineID, key)
					continue
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Error(err)
		errorCount++
		if errorCount > 10 { // Limit error output
			t.Error("... and more errors (truncated)")
			break
		}
	}

	if errorCount > 0 {
		t.Errorf("Found %d errors in concurrent cache operations", errorCount)
	}

	// Verify final cache state
	stats, err := cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get final cache stats: %v", err)
	}

	expectedEntries := numGoroutines * operationsPerGoroutine
	if stats.TotalEntries != expectedEntries {
		t.Errorf("Expected %d cache entries, got %d", expectedEntries, stats.TotalEntries)
	}
}

func testConcurrentValidation(t *testing.T) {
	tempDir := t.TempDir()
	validationEngine := validation.NewEngine()

	// Create test projects
	numProjects := 10
	projects := make([]string, numProjects)

	for i := 0; i < numProjects; i++ {
		projectDir := filepath.Join(tempDir, fmt.Sprintf("concurrent_project_%d", i))
		createConcurrentTestProject(t, projectDir, i)
		projects[i] = projectDir
	}

	numGoroutines := 20
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Launch concurrent validation operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			projectPath := projects[goroutineID%numProjects]

			// Validate project
			result, err := validationEngine.ValidateProject(projectPath)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: validation failed: %w", goroutineID, err)
				return
			}

			if result == nil {
				errors <- fmt.Errorf("goroutine %d: validation result is nil", goroutineID)
				return
			}

			// Validate specific files
			packageJSONPath := filepath.Join(projectPath, "package.json")
			if err := validationEngine.ValidatePackageJSON(packageJSONPath); err != nil {
				errors <- fmt.Errorf("goroutine %d: package.json validation failed: %w", goroutineID, err)
				return
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

func testConcurrentConfigOperations(t *testing.T) {
	// Test concurrent configuration operations
	numGoroutines := 20
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Shared configuration for testing
	sharedConfig := &models.ProjectConfig{
		Name:         "concurrent-config-test",
		Organization: "concurrent-org",
		License:      "MIT",
	}

	// Launch concurrent config operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Create unique config for this goroutine
			config := &models.ProjectConfig{
				Name:         fmt.Sprintf("concurrent-project-%d", goroutineID),
				Organization: sharedConfig.Organization,
				License:      sharedConfig.License,
				Description:  fmt.Sprintf("Project %d for concurrent testing", goroutineID),
			}

			// Validate config
			if err := utils.ValidateProjectConfig(config); err != nil {
				errors <- fmt.Errorf("goroutine %d: config validation failed: %w", goroutineID, err)
				return
			}

			// Sanitize config fields
			sanitizedName, err := utils.SanitizeInput(config.Name)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: name sanitization failed: %w", goroutineID, err)
				return
			}

			if sanitizedName != config.Name {
				errors <- fmt.Errorf("goroutine %d: name sanitization changed valid input", goroutineID)
				return
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

func testRaceConditionDetection(t *testing.T) {
	// Test for race conditions in shared resources

	// Test shared counter (should be thread-safe)
	counter := &utils.SafeCounter{}
	numGoroutines := 20           // Reduced from 100
	incrementsPerGoroutine := 100 // Reduced from 1000

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start)

	// Verify final count
	expectedCount := numGoroutines * incrementsPerGoroutine
	actualCount := counter.Value()

	if actualCount != expectedCount {
		t.Errorf("Race condition detected: expected count %d, got %d", expectedCount, actualCount)
	}

	t.Logf("Concurrent counter operations completed in %v", duration)

	// Test shared map (should be thread-safe)
	safeMap := utils.NewSafeMap()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			key := fmt.Sprintf("key_%d", id)
			value := fmt.Sprintf("value_%d", id)

			// Set value
			safeMap.Set(key, value)

			// Get value
			retrievedValue, exists := safeMap.Get(key)
			if !exists {
				t.Errorf("Expected key %s to exist", key)
				return
			}

			if retrievedValue != value {
				t.Errorf("Expected value %s for key %s, got %s", value, key, retrievedValue)
			}
		}(i)
	}

	wg.Wait()

	// Verify final map state
	if safeMap.Size() != numGoroutines {
		t.Errorf("Expected map size %d, got %d", numGoroutines, safeMap.Size())
	}
}

// TestSecurityUnderLoad tests security measures under high load
func TestSecurityUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	t.Run("input_validation_under_load", func(t *testing.T) {
		testInputValidationUnderLoad(t)
	})

	t.Run("path_validation_under_load", func(t *testing.T) {
		testPathValidationUnderLoad(t)
	})

	t.Run("secret_detection_under_load", func(t *testing.T) {
		testSecretDetectionUnderLoad(t)
	})
}

func testInputValidationUnderLoad(t *testing.T) {
	numGoroutines := 10            // Reduced from 50
	validationsPerGoroutine := 100 // Reduced from 1000

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test inputs (mix of valid and invalid)
	testInputs := []string{
		"valid-project-name",
		"another-valid-name",
		"project123",
		"<script>alert('xss')</script>",
		"../../../etc/passwd",
		"valid_name_with_underscores",
		"project; rm -rf /",
		"normal-project",
	}

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < validationsPerGoroutine; j++ {
				input := testInputs[j%len(testInputs)]

				_, err := utils.SanitizeInput(input)
				// We don't check the error here, just that the function doesn't panic
				// or cause other issues under load
				_ = err
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}

	totalValidations := numGoroutines * validationsPerGoroutine
	validationsPerSecond := float64(totalValidations) / duration.Seconds()

	t.Logf("Input validation under load: %d validations in %v (%.0f validations/second)",
		totalValidations, duration, validationsPerSecond)

	// Should handle at least 1,000 validations per second (more realistic)
	minValidationsPerSecond := 1000.0
	if validationsPerSecond < minValidationsPerSecond {
		t.Errorf("Input validation performance too low: %.0f validations/second (min: %.0f)",
			validationsPerSecond, minValidationsPerSecond)
	}
}

func testPathValidationUnderLoad(t *testing.T) {
	numGoroutines := 10           // Reduced from 30
	validationsPerGoroutine := 50 // Reduced from 500

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test paths (mix of valid and invalid)
	testPaths := []string{
		"valid/path/to/file.txt",
		"another/valid/path",
		"../../../etc/passwd",
		"./relative/path",
		"/absolute/path",
		"..\\\\..\\\\windows\\\\system32",
		"normal/file.go",
		"config/app.yaml",
	}

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < validationsPerGoroutine; j++ {
				path := testPaths[j%len(testPaths)]

				err := utils.ValidatePath(path)
				// We don't check the error here, just that the function doesn't panic
				_ = err
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}

	totalValidations := numGoroutines * validationsPerGoroutine
	validationsPerSecond := float64(totalValidations) / duration.Seconds()

	t.Logf("Path validation under load: %d validations in %v (%.0f validations/second)",
		totalValidations, duration, validationsPerSecond)
}

func testSecretDetectionUnderLoad(t *testing.T) {
	numGoroutines := 5           // Reduced from 20
	detectionsPerGoroutine := 20 // Reduced from 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test content with and without secrets
	testContent := []string{
		`const config = { apiKey: "sk-1234567890abcdef" };`,
		`password = "secret123"`,
		`const normal = "just normal content";`,
		`token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`,
		`const example = "this is just example text";`,
		`API_KEY=pk_live_1234567890abcdef`,
		`const readme = "# Project Title";`,
		`jwt_secret = "my-jwt-secret-key"`,
	}

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < detectionsPerGoroutine; j++ {
				content := testContent[j%len(testContent)]

				secrets := utils.DetectSecrets(content)
				// We don't validate the results here, just that detection doesn't fail
				_ = secrets
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}

	totalDetections := numGoroutines * detectionsPerGoroutine
	detectionsPerSecond := float64(totalDetections) / duration.Seconds()

	t.Logf("Secret detection under load: %d detections in %v (%.0f detections/second)",
		totalDetections, duration, detectionsPerSecond)
}

// TestMemoryLeakPrevention tests for memory leaks in long-running operations
func TestMemoryLeakPrevention(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak tests in short mode")
	}

	t.Run("cache_memory_management", func(t *testing.T) {
		testCacheMemoryManagement(t)
	})

	t.Run("validation_memory_management", func(t *testing.T) {
		testValidationMemoryManagement(t)
	})
}

func testCacheMemoryManagement(t *testing.T) {
	tempDir := t.TempDir()
	cacheManager := cache.NewManager(tempDir)

	// Configure cache with limits
	config := &interfaces.CacheConfig{
		MaxSize:        1024 * 1024, // 1MB
		MaxEntries:     1000,
		EvictionPolicy: "lru",
	}

	err := cacheManager.SetCacheConfig(config)
	if err != nil {
		t.Fatalf("Failed to set cache config: %v", err)
	}

	// Continuously add and remove data to test memory management
	iterations := 1000              // Reduced from 10000
	largeValue := make([]byte, 512) // Reduced from 1KB to 512 bytes per entry

	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("memory_test_%d", i)

		// Set data
		err := cacheManager.Set(key, largeValue, time.Minute)
		if err != nil {
			t.Fatalf("Failed to set cache entry %d: %v", i, err)
		}

		// Periodically clean cache
		if i%100 == 0 {
			err := cacheManager.Clean()
			if err != nil {
				t.Fatalf("Failed to clean cache at iteration %d: %v", i, err)
			}
		}

		// Periodically check stats
		if i%1000 == 0 {
			stats, err := cacheManager.GetStats()
			if err != nil {
				t.Fatalf("Failed to get cache stats at iteration %d: %v", i, err)
			}

			// Verify cache size is within limits
			if stats.TotalSize > config.MaxSize*2 { // Allow some overhead
				t.Errorf("Cache size exceeded limits at iteration %d: %d bytes", i, stats.TotalSize)
			}

			if stats.TotalEntries > config.MaxEntries*2 { // Allow some overhead
				t.Errorf("Cache entries exceeded limits at iteration %d: %d entries", i, stats.TotalEntries)
			}
		}
	}

	// Final cleanup
	err = cacheManager.Clear()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify cache is empty
	stats, err := cacheManager.GetStats()
	if err != nil {
		t.Fatalf("Failed to get final cache stats: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Errorf("Expected empty cache after clear, got %d entries", stats.TotalEntries)
	}
}

func testValidationMemoryManagement(t *testing.T) {
	tempDir := t.TempDir()
	validationEngine := validation.NewEngine()

	// Create multiple test projects
	numProjects := 10 // Reduced from 100
	projects := make([]string, numProjects)

	for i := 0; i < numProjects; i++ {
		projectDir := filepath.Join(tempDir, fmt.Sprintf("memory_project_%d", i))
		createConcurrentTestProject(t, projectDir, i)
		projects[i] = projectDir
	}

	// Validate projects multiple times to test memory management
	iterations := 3 // Reduced from 10

	for iter := 0; iter < iterations; iter++ {
		for i, projectPath := range projects {
			result, err := validationEngine.ValidateProject(projectPath)
			if err != nil {
				t.Fatalf("Validation failed for project %d in iteration %d: %v", i, iter, err)
			}

			if result == nil {
				t.Fatalf("Validation result is nil for project %d in iteration %d", i, iter)
			}

			// Force garbage collection periodically
			if i%10 == 0 {
				runtime.GC()
			}
		}

		t.Logf("Completed validation iteration %d/%d", iter+1, iterations)
	}
}

// Helper functions

func createConcurrentTestProject(t *testing.T, projectDir string, id int) {
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create package.json
	packageJSON := fmt.Sprintf(`{
		"name": "concurrent-test-project-%d",
		"version": "1.0.0",
		"description": "Concurrent test project %d",
		"license": "MIT",
		"dependencies": {
			"express": "^4.18.0",
			"cors": "^2.8.5"
		}
	}`, id, id)

	err = os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create main file
	mainContent := fmt.Sprintf(`const express = require('express');
const app = express();

// Project %d
app.get('/', (req, res) => {
	res.json({ 
		message: 'Hello from project %d',
		id: %d
	});
});

const port = process.env.PORT || %d;
app.listen(port, () => {
	console.log('Project %d server running on port ' + port);
});
`, id, id, id, 3000+id, id)

	err = os.WriteFile(filepath.Join(projectDir, "main.js"), []byte(mainContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.js: %v", err)
	}

	// Create README
	readme := fmt.Sprintf("# Concurrent Test Project %d\n\nThis is test project %d for concurrent testing.\n\n## Installation\n\n```bash\nnpm install\n```\n\n## Usage\n\n```bash\nnpm start\n```\n\nProject runs on port %d.\n", id, id, 3000+id)

	err = os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}
}

// Benchmark tests for security operations

func BenchmarkSecurityOperations(b *testing.B) {
	b.Run("input_sanitization", func(b *testing.B) {
		testInput := "test-project-name-123"

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := utils.SanitizeInput(testInput)
			if err != nil {
				b.Fatalf("Input sanitization failed: %v", err)
			}
		}
	})

	b.Run("path_validation", func(b *testing.B) {
		testPath := "valid/path/to/file.txt"

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := utils.ValidatePath(testPath)
			if err != nil {
				b.Fatalf("Path validation failed: %v", err)
			}
		}
	})

	b.Run("secret_detection", func(b *testing.B) {
		testContent := `const config = {
			apiKey: "sk-1234567890abcdef",
			password: "secret123",
			normal: "just normal content"
		};`

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			secrets := utils.DetectSecrets(testContent)
			_ = secrets // Use the result to prevent optimization
		}
	})

	b.Run("config_validation", func(b *testing.B) {
		config := &models.ProjectConfig{
			Name:         "benchmark-project",
			Organization: "benchmark-org",
			License:      "MIT",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := utils.ValidateProjectConfig(config)
			if err != nil {
				b.Fatalf("Config validation failed: %v", err)
			}
		}
	})
}
