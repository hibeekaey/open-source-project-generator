package validation

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceTester_MeasureProjectGeneration(t *testing.T) {
	tester := NewPerformanceTester()

	// Create a test project
	tempDir, err := os.MkdirTemp("", "perf-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a simple project structure
	err = createSimpleTestProject(tempDir)
	require.NoError(t, err)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
		},
	}

	// Measure performance
	metrics, err := tester.MeasureProjectGeneration(tempDir, config)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	// Verify metrics are populated
	assert.Greater(t, metrics.ValidationTime, time.Duration(0))
	assert.Greater(t, metrics.SetupTime, time.Duration(0))
	assert.Greater(t, metrics.VerificationTime, time.Duration(0))
	assert.Greater(t, metrics.TotalTime, time.Duration(0))
	assert.GreaterOrEqual(t, metrics.MemoryUsage.StartMemory, uint64(0))
	assert.GreaterOrEqual(t, metrics.MemoryUsage.EndMemory, uint64(0))
	assert.Greater(t, metrics.FileSystemMetrics.FilesCreated, 0)
	assert.Greater(t, metrics.FileSystemMetrics.DirectoriesCreated, 0)
	assert.Greater(t, metrics.FileSystemMetrics.TotalSize, int64(0))
}

func TestPerformanceTester_MeasureValidationPerformance(t *testing.T) {
	tester := NewPerformanceTester()

	// Create multiple test projects
	projectPaths := []string{}
	for i := 0; i < 3; i++ {
		tempDir, err := os.MkdirTemp("", "perf-validation-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		err = createSimpleTestProject(tempDir)
		require.NoError(t, err)

		projectPaths = append(projectPaths, tempDir)
	}

	// Measure validation performance
	results, err := tester.MeasureValidationPerformance(projectPaths)
	require.NoError(t, err)
	require.NotNil(t, results)

	// Verify results for each project
	assert.Len(t, results, len(projectPaths))
	for projectName, metrics := range results {
		assert.NotEmpty(t, projectName)
		assert.Greater(t, metrics.ValidationTime, time.Duration(0))
		assert.Greater(t, metrics.TotalTime, time.Duration(0))
		assert.GreaterOrEqual(t, metrics.MemoryUsage.StartMemory, uint64(0))
		assert.Greater(t, metrics.FileSystemMetrics.FilesCreated, 0)
	}
}

func TestPerformanceTester_MeasureScalabilityPerformance(t *testing.T) {
	tester := NewPerformanceTester()

	// Test with different file counts
	fileCounts := []int{10, 50, 100}

	results, err := tester.MeasureScalabilityPerformance("", fileCounts)
	require.NoError(t, err)
	require.NotNil(t, results)

	// Verify results for each file count
	assert.Len(t, results, len(fileCounts))

	var previousTime time.Duration
	for _, fileCount := range fileCounts {
		metrics, exists := results[fileCount]
		require.True(t, exists)

		assert.Greater(t, metrics.ValidationTime, time.Duration(0))
		assert.Equal(t, fileCount, metrics.FileSystemMetrics.FilesCreated)

		// Validation time should generally increase with file count
		if previousTime > 0 {
			// Allow some variance, but generally expect longer times for more files
			assert.True(t, metrics.ValidationTime >= previousTime/2,
				"Validation time should scale with file count: %v vs %v for %d files",
				metrics.ValidationTime, previousTime, fileCount)
		}
		previousTime = metrics.ValidationTime
	}
}

func TestPerformanceTester_MeasureFileSystemMetrics(t *testing.T) {
	tester := NewPerformanceTester()

	// Create a test directory with known structure
	tempDir, err := os.MkdirTemp("", "fs-metrics-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	testFiles := map[string]string{
		"file1.txt":           "small content",
		"file2.txt":           "larger content with more text to make it bigger",
		"dir1/file3.txt":      "nested file content",
		"dir1/dir2/file4.txt": "deeply nested file",
	}

	expectedFiles := 0
	expectedDirs := 1 // tempDir itself
	totalSize := int64(0)

	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)

		expectedFiles++
		totalSize += int64(len(content))
	}

	// Count expected directories (dir1, dir1/dir2)
	expectedDirs += 2

	// Measure file system metrics
	metrics, err := tester.measureFileSystemMetrics(tempDir)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	assert.Equal(t, expectedFiles, metrics.FilesCreated)
	assert.Equal(t, expectedDirs, metrics.DirectoriesCreated)
	assert.Equal(t, totalSize, metrics.TotalSize)
	assert.Greater(t, metrics.LargestFile, int64(0))
	assert.Greater(t, metrics.SmallestFile, int64(0))
	assert.LessOrEqual(t, metrics.SmallestFile, metrics.LargestFile)
}

func TestPerformanceTester_GenerateTestFiles(t *testing.T) {
	tester := NewPerformanceTester()

	tests := []struct {
		name      string
		fileCount int
	}{
		{"small project", 10},
		{"medium project", 50},
		{"large project", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "generate-files-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			err = tester.generateTestFiles(tempDir, tt.fileCount)
			require.NoError(t, err)

			// Count generated files
			fileCount := 0
			err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					fileCount++
				}
				return nil
			})
			require.NoError(t, err)

			assert.Equal(t, tt.fileCount, fileCount)
		})
	}
}

func TestPerformanceTester_GeneratePerformanceReport(t *testing.T) {
	tester := NewPerformanceTester()

	metrics := &PerformanceMetrics{
		GenerationTime:   100 * time.Millisecond,
		ValidationTime:   50 * time.Millisecond,
		SetupTime:        30 * time.Millisecond,
		VerificationTime: 20 * time.Millisecond,
		TotalTime:        200 * time.Millisecond,
		MemoryUsage: MemoryUsage{
			StartMemory:    10,
			PeakMemory:     15,
			EndMemory:      12,
			MemoryDelta:    2,
			GCCount:        3,
			AllocatedBytes: 1024,
		},
		FileSystemMetrics: FileSystemMetrics{
			FilesCreated:       50,
			DirectoriesCreated: 10,
			TotalSize:          1024 * 1024, // 1 MB
			LargestFile:        50 * 1024,   // 50 KB
			SmallestFile:       1024,        // 1 KB
		},
	}

	report := tester.GeneratePerformanceReport(metrics)

	// Verify report contains expected sections
	assert.Contains(t, report, "Performance Report")
	assert.Contains(t, report, "Timing Metrics")
	assert.Contains(t, report, "Memory Usage")
	assert.Contains(t, report, "File System Metrics")
	assert.Contains(t, report, "Performance Summary")

	// Verify specific values are included
	assert.Contains(t, report, "100ms")             // Generation time
	assert.Contains(t, report, "Files Created: 50") // Files created
	assert.Contains(t, report, "10 MB")             // Start memory
	assert.Contains(t, report, "1048576 bytes")     // Total size
}

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected uint64
	}{
		{0, 0},
		{1024, 0},
		{1024 * 1024, 1},
		{1024 * 1024 * 5, 5},
		{1024*1024*10 + 512*1024, 10}, // 10.5 MB rounds down to 10
	}

	for _, tt := range tests {
		result := bytesToMB(tt.bytes)
		assert.Equal(t, tt.expected, result, "bytesToMB(%d) = %d, expected %d", tt.bytes, result, tt.expected)
	}
}

// Helper function to create a simple test project
func createSimpleTestProject(projectPath string) error {
	// Create directories
	dirs := []string{
		"frontend/src",
		"backend/internal",
		"docs",
		".github/workflows",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	// Create files
	files := map[string]string{
		"README.md": "# Test Project\n\nThis is a test project.",
		"Makefile":  "all:\n\techo 'Building...'\n",
		"frontend/package.json": `{
			"name": "test-frontend",
			"version": "1.0.0",
			"scripts": {
				"build": "echo 'Building frontend...'",
				"dev": "echo 'Starting dev server...'"
			},
			"dependencies": {}
		}`,
		"backend/go.mod": `module test-backend

go 1.24`,
		"backend/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"docker-compose.yml": `version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"`,
		".github/workflows/ci.yml": `name: CI
on:
  push:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectPath, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// Benchmark tests
func BenchmarkPerformanceTester_MeasureValidation(b *testing.B) {
	tester := NewPerformanceTester()

	// Create a test project
	tempDir, err := os.MkdirTemp("", "benchmark-validation-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	err = createSimpleTestProject(tempDir)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.MeasureValidationPerformance([]string{tempDir})
		require.NoError(b, err)
	}
}

func BenchmarkPerformanceTester_GenerateTestFiles(b *testing.B) {
	tester := NewPerformanceTester()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempDir, err := os.MkdirTemp("", "benchmark-generate-*")
		require.NoError(b, err)

		err = tester.generateTestFiles(tempDir, 100)
		require.NoError(b, err)

		os.RemoveAll(tempDir)
	}
}

func BenchmarkPerformanceTester_MeasureFileSystemMetrics(b *testing.B) {
	tester := NewPerformanceTester()

	// Create a test project once
	tempDir, err := os.MkdirTemp("", "benchmark-fs-metrics-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	err = tester.generateTestFiles(tempDir, 100)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.measureFileSystemMetrics(tempDir)
		require.NoError(b, err)
	}
}
