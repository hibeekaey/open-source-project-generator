package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestSecureFileOperationsIntegration tests end-to-end secure file operations
func TestSecureFileOperationsIntegration(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("AtomicWriteOperations", func(t *testing.T) {
		testAtomicWriteOperations(t, tempDir)
	})

	t.Run("SecureTempFileCreation", func(t *testing.T) {
		testSecureTempFileCreation(t, tempDir)
	})

	t.Run("PathValidationSecurity", func(t *testing.T) {
		testPathValidationSecurity(t, tempDir)
	})

	t.Run("SecureDeleteOperations", func(t *testing.T) {
		testSecureDeleteOperations(t, tempDir)
	})

	t.Run("PermissionManagement", func(t *testing.T) {
		testPermissionManagement(t, tempDir)
	})
}

// testAtomicWriteOperations verifies that file writes are truly atomic
func testAtomicWriteOperations(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()
	testFile := filepath.Join(tempDir, "atomic_test.txt")

	// Test data
	testData := []byte("This is test data for atomic operations")

	// Write the file atomically
	err := secureOps.WriteFileAtomic(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("WriteFileAtomic failed: %v", err)
	}

	// Verify the file exists and has correct content
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("File should exist after atomic write")
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != string(testData) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testData, content)
	}

	// Verify no temporary files remain
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".tmp.") {
			t.Errorf("Temporary file should not remain: %s", file.Name())
		}
	}

	t.Logf("Atomic write operations verified successfully")
}

// testSecureTempFileCreation verifies secure temporary file creation
func testSecureTempFileCreation(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Create multiple temporary files to verify randomness
	tempFiles := make([]*os.File, 10)
	tempNames := make([]string, 10)

	for i := 0; i < 10; i++ {
		tempFile, err := secureOps.CreateSecureTempFile(tempDir, "test_")
		if err != nil {
			t.Fatalf("CreateSecureTempFile failed: %v", err)
		}

		tempFiles[i] = tempFile
		tempNames[i] = tempFile.Name()

		// Verify the file has secure permissions (0600)
		info, err := tempFile.Stat()
		if err != nil {
			t.Fatalf("Failed to stat temp file: %v", err)
		}

		expectedPerm := os.FileMode(0600)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("Expected permissions %v, got %v", expectedPerm, info.Mode().Perm())
		}
	}

	// Verify all temp file names are different (randomness check)
	nameSet := make(map[string]bool)
	for _, name := range tempNames {
		if nameSet[name] {
			t.Errorf("Duplicate temp file name detected: %s", name)
		}
		nameSet[name] = true

		// Verify the name contains the pattern and a random suffix
		if !strings.Contains(filepath.Base(name), "test_") {
			t.Errorf("Temp file name should contain pattern: %s", name)
		}

		// Verify the random suffix is not predictable (no timestamps)
		baseName := filepath.Base(name)
		if strings.Contains(baseName, "tmp.") {
			parts := strings.Split(baseName, "tmp.")
			if len(parts) > 1 && len(parts[1]) >= 10 {
				// Check if it looks like a timestamp
				t.Errorf("Temp file name may contain timestamp pattern: %s", name)
			}
		}
	}

	// Clean up
	for _, tempFile := range tempFiles {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	t.Logf("Secure temp file creation verified successfully")
}

// testPathValidationSecurity verifies path validation prevents security issues
func testPathValidationSecurity(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Test cases for path validation
	testCases := []struct {
		name        string
		path        string
		allowedDirs []string
		shouldFail  bool
	}{
		{
			name:        "ValidPath",
			path:        filepath.Join(tempDir, "valid.txt"),
			allowedDirs: []string{tempDir},
			shouldFail:  false,
		},
		{
			name:        "DirectoryTraversal",
			path:        filepath.Join(tempDir, "../../../etc/passwd"),
			allowedDirs: []string{tempDir},
			shouldFail:  true,
		},
		{
			name:        "RelativeTraversal",
			path:        "../../etc/passwd",
			allowedDirs: []string{tempDir},
			shouldFail:  true,
		},
		{
			name:        "SystemDirectory",
			path:        "/etc/passwd",
			allowedDirs: nil, // Use default dangerous path checks
			shouldFail:  true,
		},
		{
			name:        "OutsideAllowedDir",
			path:        "/tmp/outside.txt",
			allowedDirs: []string{tempDir},
			shouldFail:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := secureOps.ValidatePath(tc.path, tc.allowedDirs)

			if tc.shouldFail && err == nil {
				t.Errorf("Expected path validation to fail for: %s", tc.path)
			}

			if !tc.shouldFail && err != nil {
				t.Errorf("Expected path validation to succeed for: %s, got error: %v", tc.path, err)
			}
		})
	}

	t.Logf("Path validation security verified successfully")
}

// testSecureDeleteOperations verifies secure file deletion
func testSecureDeleteOperations(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Create a test file with sensitive content
	testFile := filepath.Join(tempDir, "sensitive.txt")
	sensitiveData := []byte("This is sensitive data that should be securely deleted")

	err := os.WriteFile(testFile, sensitiveData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file should exist before deletion")
	}

	// Securely delete the file
	err = secureOps.SecureDelete(testFile)
	if err != nil {
		t.Fatalf("SecureDelete failed: %v", err)
	}

	// Verify file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("File should not exist after secure deletion")
	}

	// Test deleting non-existent file (should not error)
	err = secureOps.SecureDelete(filepath.Join(tempDir, "nonexistent.txt"))
	if err != nil {
		t.Errorf("SecureDelete should not error on non-existent file: %v", err)
	}

	t.Logf("Secure delete operations verified successfully")
}

// testPermissionManagement verifies secure permission setting
func testPermissionManagement(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Create a test file
	testFile := filepath.Join(tempDir, "permissions.txt")
	err := os.WriteFile(testFile, []byte("test"), 0777)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set secure permissions
	securePerms := os.FileMode(0600)
	err = secureOps.EnsureSecurePermissions(testFile, securePerms)
	if err != nil {
		t.Fatalf("EnsureSecurePermissions failed: %v", err)
	}

	// Verify permissions were set correctly
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Mode().Perm() != securePerms {
		t.Errorf("Expected permissions %v, got %v", securePerms, info.Mode().Perm())
	}

	t.Logf("Permission management verified successfully")
}

// TestConcurrentSecureOperations tests concurrent access scenarios to verify race condition prevention
func TestConcurrentSecureOperations(t *testing.T) {
	tempDir := t.TempDir()
	secureOps := NewSecureFileOperations()

	t.Run("ConcurrentAtomicWrites", func(t *testing.T) {
		testConcurrentAtomicWrites(t, tempDir, secureOps)
	})

	t.Run("ConcurrentTempFileCreation", func(t *testing.T) {
		testConcurrentTempFileCreation(t, tempDir, secureOps)
	})

	t.Run("ConcurrentRandomGeneration", func(t *testing.T) {
		testConcurrentRandomGeneration(t)
	})
}

// testConcurrentAtomicWrites verifies atomic writes work correctly under concurrent access
func testConcurrentAtomicWrites(t *testing.T, tempDir string, secureOps SecureFileOperations) {
	testFile := filepath.Join(tempDir, "concurrent_writes.txt")
	numGoroutines := 10
	numWritesPerGoroutine := 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numWritesPerGoroutine)

	// Launch concurrent writers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numWritesPerGoroutine; j++ {
				data := []byte(fmt.Sprintf("Goroutine %d, Write %d", goroutineID, j))
				err := secureOps.WriteFileAtomic(testFile, data, 0644)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, write %d: %w", goroutineID, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent write error: %v", err)
	}

	// Verify the final file exists and is readable
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Final file should exist after concurrent writes")
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read final file: %v", err)
	}

	// The content should be from one of the writes (atomic operation)
	if len(content) == 0 {
		t.Errorf("File should not be empty after concurrent writes")
	}

	t.Logf("Concurrent atomic writes completed successfully")
}

// testConcurrentTempFileCreation verifies temp file creation works under concurrent access
func testConcurrentTempFileCreation(t *testing.T, tempDir string, secureOps SecureFileOperations) {
	numGoroutines := 20

	var wg sync.WaitGroup
	tempFiles := make(chan *os.File, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Launch concurrent temp file creators
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			tempFile, err := secureOps.CreateSecureTempFile(tempDir, fmt.Sprintf("concurrent_%d_", goroutineID))
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: %w", goroutineID, err)
				return
			}

			tempFiles <- tempFile
		}(i)
	}

	wg.Wait()
	close(tempFiles)
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent temp file creation error: %v", err)
	}

	// Collect and verify temp files
	createdFiles := make([]*os.File, 0, numGoroutines)
	fileNames := make(map[string]bool)

	for tempFile := range tempFiles {
		createdFiles = append(createdFiles, tempFile)

		// Verify unique names
		name := tempFile.Name()
		if fileNames[name] {
			t.Errorf("Duplicate temp file name: %s", name)
		}
		fileNames[name] = true
	}

	// Clean up
	for _, tempFile := range createdFiles {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	t.Logf("Concurrent temp file creation completed successfully, created %d unique files", len(createdFiles))
}

// testConcurrentRandomGeneration verifies random generation works under concurrent access
func testConcurrentRandomGeneration(t *testing.T) {
	secureRandom := NewSecureRandom()
	numGoroutines := 50
	numGenerationsPerGoroutine := 10

	var wg sync.WaitGroup
	randomValues := make(chan string, numGoroutines*numGenerationsPerGoroutine)
	errors := make(chan error, numGoroutines*numGenerationsPerGoroutine)

	// Launch concurrent random generators
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numGenerationsPerGoroutine; j++ {
				randomValue, err := secureRandom.GenerateRandomSuffix(16)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, generation %d: %w", goroutineID, j, err)
					continue
				}

				randomValues <- randomValue
			}
		}(i)
	}

	wg.Wait()
	close(randomValues)
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent random generation error: %v", err)
	}

	// Verify uniqueness of generated values
	valueSet := make(map[string]bool)
	duplicates := 0

	for value := range randomValues {
		if valueSet[value] {
			duplicates++
		}
		valueSet[value] = true
	}

	// With cryptographically secure random generation, duplicates should be extremely rare
	if duplicates > 0 {
		t.Errorf("Found %d duplicate random values out of %d generated", duplicates, len(valueSet))
	}

	t.Logf("Concurrent random generation completed successfully, generated %d unique values", len(valueSet))
}

// TestSecureOperationsPerformance adds performance benchmarks comparing secure vs insecure operations
func TestSecureOperationsPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("AtomicWritePerformance", func(t *testing.T) {
		benchmarkAtomicWritePerformance(t, tempDir)
	})

	t.Run("RandomGenerationPerformance", func(t *testing.T) {
		benchmarkRandomGenerationPerformance(t)
	})

	t.Run("TempFileCreationPerformance", func(t *testing.T) {
		benchmarkTempFileCreationPerformance(t, tempDir)
	})
}

// benchmarkAtomicWritePerformance compares atomic vs direct write performance
func benchmarkAtomicWritePerformance(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()
	testData := []byte("This is test data for performance benchmarking")
	iterations := 100

	// Benchmark atomic writes
	atomicFile := filepath.Join(tempDir, "atomic_perf.txt")
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		err := secureOps.WriteFileAtomic(atomicFile, testData, 0644)
		if err != nil {
			t.Fatalf("Atomic write failed: %v", err)
		}
	}

	atomicDuration := time.Since(startTime)

	// Benchmark direct writes
	directFile := filepath.Join(tempDir, "direct_perf.txt")
	startTime = time.Now()

	for i := 0; i < iterations; i++ {
		err := os.WriteFile(directFile, testData, 0644)
		if err != nil {
			t.Fatalf("Direct write failed: %v", err)
		}
	}

	directDuration := time.Since(startTime)

	// Calculate overhead
	overhead := float64(atomicDuration-directDuration) / float64(directDuration) * 100

	t.Logf("Atomic write performance:")
	t.Logf("  Atomic writes: %v (%v per operation)", atomicDuration, atomicDuration/time.Duration(iterations))
	t.Logf("  Direct writes: %v (%v per operation)", directDuration, directDuration/time.Duration(iterations))
	t.Logf("  Overhead: %.2f%%", overhead)

	// Atomic writes should not be more than 2000% slower than direct writes
	// This is expected due to the additional security operations (temp file creation, sync, rename)
	if overhead > 2000 {
		t.Errorf("Atomic write overhead too high: %.2f%%", overhead)
	}
}

// benchmarkRandomGenerationPerformance compares secure vs insecure random generation
func benchmarkRandomGenerationPerformance(t *testing.T) {
	secureRandom := NewSecureRandom()
	iterations := 1000
	length := 16

	// Benchmark secure random generation
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		_, err := secureRandom.GenerateRandomSuffix(length)
		if err != nil {
			t.Fatalf("Secure random generation failed: %v", err)
		}
	}

	secureDuration := time.Since(startTime)

	t.Logf("Random generation performance:")
	t.Logf("  Secure random: %v (%v per operation)", secureDuration, secureDuration/time.Duration(iterations))
	t.Logf("  Generated %d random strings of length %d", iterations, length)

	// Secure random generation should complete within reasonable time
	avgPerOperation := secureDuration / time.Duration(iterations)
	if avgPerOperation > time.Millisecond {
		t.Errorf("Secure random generation too slow: %v per operation", avgPerOperation)
	}
}

// benchmarkTempFileCreationPerformance compares secure vs standard temp file creation
func benchmarkTempFileCreationPerformance(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()
	iterations := 100

	// Benchmark secure temp file creation
	startTime := time.Now()
	secureFiles := make([]*os.File, iterations)

	for i := 0; i < iterations; i++ {
		tempFile, err := secureOps.CreateSecureTempFile(tempDir, "perf_")
		if err != nil {
			t.Fatalf("Secure temp file creation failed: %v", err)
		}
		secureFiles[i] = tempFile
	}

	secureDuration := time.Since(startTime)

	// Clean up secure files
	for _, file := range secureFiles {
		file.Close()
		os.Remove(file.Name())
	}

	// Benchmark standard temp file creation
	startTime = time.Now()
	standardFiles := make([]*os.File, iterations)

	for i := 0; i < iterations; i++ {
		tempFile, err := os.CreateTemp(tempDir, "perf_")
		if err != nil {
			t.Fatalf("Standard temp file creation failed: %v", err)
		}
		standardFiles[i] = tempFile
	}

	standardDuration := time.Since(startTime)

	// Clean up standard files
	for _, file := range standardFiles {
		file.Close()
		os.Remove(file.Name())
	}

	// Calculate overhead
	overhead := float64(secureDuration-standardDuration) / float64(standardDuration) * 100

	t.Logf("Temp file creation performance:")
	t.Logf("  Secure creation: %v (%v per operation)", secureDuration, secureDuration/time.Duration(iterations))
	t.Logf("  Standard creation: %v (%v per operation)", standardDuration, standardDuration/time.Duration(iterations))
	t.Logf("  Overhead: %.2f%%", overhead)

	// Secure temp file creation should not be more than 200% slower
	if overhead > 200 {
		t.Errorf("Secure temp file creation overhead too high: %.2f%%", overhead)
	}
}

// TestSecurityFocusedIntegrationSuite creates a comprehensive security-focused integration test suite
func TestSecurityFocusedIntegrationSuite(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("EndToEndSecureWorkflow", func(t *testing.T) {
		testEndToEndSecureWorkflow(t, tempDir)
	})

	t.Run("SecurityVulnerabilityPrevention", func(t *testing.T) {
		testSecurityVulnerabilityPrevention(t, tempDir)
	})

	t.Run("ErrorHandlingSecurityImplications", func(t *testing.T) {
		testErrorHandlingSecurityImplications(t, tempDir)
	})
}

// testEndToEndSecureWorkflow tests a complete secure workflow
func testEndToEndSecureWorkflow(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()
	secureRandom := NewSecureRandom()

	// Step 1: Generate secure random data
	randomData, err := secureRandom.GenerateBytes(1024)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	// Step 2: Create secure temp file
	tempFile, err := secureOps.CreateSecureTempFile(tempDir, "workflow_")
	if err != nil {
		t.Fatalf("Failed to create secure temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Step 3: Write data atomically
	finalPath := filepath.Join(tempDir, "final_data.bin")
	err = secureOps.WriteFileAtomic(finalPath, randomData, 0600)
	if err != nil {
		t.Fatalf("Failed to write file atomically: %v", err)
	}

	// Step 4: Verify file integrity
	readData, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatalf("Failed to read final file: %v", err)
	}

	if len(readData) != len(randomData) {
		t.Errorf("Data length mismatch: expected %d, got %d", len(randomData), len(readData))
	}

	// Step 5: Secure cleanup
	err = secureOps.SecureDelete(finalPath)
	if err != nil {
		t.Fatalf("Failed to securely delete file: %v", err)
	}

	// Clean up temp file if it still exists
	os.Remove(tempPath)

	t.Logf("End-to-end secure workflow completed successfully")
}

// testSecurityVulnerabilityPrevention tests prevention of common security vulnerabilities
func testSecurityVulnerabilityPrevention(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Test 1: Prevent directory traversal
	maliciousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"/etc/shadow",
		"C:\\Windows\\System32\\config\\SAM",
	}

	for _, path := range maliciousPaths {
		err := secureOps.ValidatePath(path, []string{tempDir})
		if err == nil {
			t.Errorf("Should have rejected malicious path: %s", path)
		}
	}

	// Test 2: Prevent predictable temp file names
	tempFiles := make([]*os.File, 100)
	for i := 0; i < 100; i++ {
		tempFile, err := secureOps.CreateSecureTempFile(tempDir, "security_test_")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tempFiles[i] = tempFile
	}

	// Analyze temp file names for patterns
	names := make([]string, len(tempFiles))
	for i, file := range tempFiles {
		names[i] = filepath.Base(file.Name())
		file.Close()
		os.Remove(file.Name())
	}

	// Check for timestamp patterns or sequential patterns
	for i, name := range names {
		for j, otherName := range names {
			if i != j && name == otherName {
				t.Errorf("Duplicate temp file names detected: %s", name)
			}
		}

		// Check for timestamp-like patterns (long sequences of digits)
		if strings.Contains(name, "tmp.") {
			parts := strings.Split(name, "tmp.")
			if len(parts) > 1 && len(parts[1]) >= 10 {
				// Could be a timestamp
				t.Errorf("Potential timestamp pattern in temp file name: %s", name)
			}
		}
	}

	t.Logf("Security vulnerability prevention tests passed")
}

// testErrorHandlingSecurityImplications tests security implications of error handling
func testErrorHandlingSecurityImplications(t *testing.T, tempDir string) {
	secureOps := NewSecureFileOperations()

	// Test 1: Ensure errors don't leak sensitive path information
	sensitiveDir := "/etc"
	err := secureOps.ValidatePath(filepath.Join(sensitiveDir, "passwd"), []string{tempDir})
	if err != nil {
		errorMsg := err.Error()
		// Error message should not contain the full sensitive path
		if strings.Contains(errorMsg, "/etc/passwd") {
			t.Errorf("Error message leaks sensitive path information: %s", errorMsg)
		}
	}

	// Test 2: Ensure failed operations clean up properly
	invalidDir := "/nonexistent/directory"
	err = secureOps.WriteFileAtomic(filepath.Join(invalidDir, "test.txt"), []byte("test"), 0644)
	if err == nil {
		t.Errorf("Should have failed to write to nonexistent directory")
	}

	// Verify no temp files were left behind in any accessible location
	files, _ := os.ReadDir(tempDir)
	for _, file := range files {
		if strings.Contains(file.Name(), ".tmp.") {
			t.Errorf("Temp file left behind after failed operation: %s", file.Name())
		}
	}

	t.Logf("Error handling security implications tests passed")
}
