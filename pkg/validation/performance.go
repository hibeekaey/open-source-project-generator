package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// PerformanceMetrics holds performance measurement data
type PerformanceMetrics struct {
	GenerationTime    time.Duration     `json:"generation_time"`
	ValidationTime    time.Duration     `json:"validation_time"`
	SetupTime         time.Duration     `json:"setup_time"`
	VerificationTime  time.Duration     `json:"verification_time"`
	TotalTime         time.Duration     `json:"total_time"`
	MemoryUsage       MemoryUsage       `json:"memory_usage"`
	FileSystemMetrics FileSystemMetrics `json:"filesystem_metrics"`
}

// MemoryUsage holds memory usage statistics
type MemoryUsage struct {
	StartMemory    uint64 `json:"start_memory_mb"`
	PeakMemory     uint64 `json:"peak_memory_mb"`
	EndMemory      uint64 `json:"end_memory_mb"`
	MemoryDelta    int64  `json:"memory_delta_mb"`
	GCCount        uint32 `json:"gc_count"`
	AllocatedBytes uint64 `json:"allocated_bytes"`
}

// FileSystemMetrics holds file system operation statistics
type FileSystemMetrics struct {
	FilesCreated       int   `json:"files_created"`
	DirectoriesCreated int   `json:"directories_created"`
	TotalSize          int64 `json:"total_size_bytes"`
	LargestFile        int64 `json:"largest_file_bytes"`
	SmallestFile       int64 `json:"smallest_file_bytes"`
}

// PerformanceTester provides performance testing capabilities
type PerformanceTester struct {
	validationEngine interfaces.ValidationEngine
	setupEngine      *SetupEngine
}

// NewPerformanceTester creates a new performance tester
func NewPerformanceTester() *PerformanceTester {
	return &PerformanceTester{
		validationEngine: NewEngine(),
		setupEngine:      NewSetupEngine(),
	}
}

// MeasureProjectGeneration measures the performance of project generation and validation
func (pt *PerformanceTester) MeasureProjectGeneration(projectPath string, config *models.ProjectConfig) (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{}

	// Record initial memory state
	var startMemStats, endMemStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&startMemStats)

	totalStartTime := time.Now()

	// Measure validation time
	validationStartTime := time.Now()
	validationResult, err := pt.validationEngine.ValidateProject(projectPath)
	metrics.ValidationTime = time.Since(validationStartTime)

	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if !validationResult.Valid {
		return metrics, fmt.Errorf("project validation failed with %d errors", len(validationResult.Errors))
	}

	// Measure setup time
	setupStartTime := time.Now()
	setupResult, err := pt.setupEngine.SetupProject(projectPath, config)
	metrics.SetupTime = time.Since(setupStartTime)

	if err != nil {
		return nil, fmt.Errorf("setup failed: %w", err)
	}

	if !setupResult.Valid {
		return metrics, fmt.Errorf("project setup failed with %d errors", len(setupResult.Errors))
	}

	// Measure verification time
	verificationStartTime := time.Now()
	verificationResult, err := pt.setupEngine.VerifyProject(projectPath, config)
	metrics.VerificationTime = time.Since(verificationStartTime)

	if err != nil {
		return nil, fmt.Errorf("verification failed: %w", err)
	}

	if !verificationResult.Valid {
		return metrics, fmt.Errorf("project verification failed with %d errors", len(verificationResult.Errors))
	}

	// Calculate total time
	metrics.TotalTime = time.Since(totalStartTime)

	// Record final memory state
	runtime.GC()
	runtime.ReadMemStats(&endMemStats)

	// Calculate memory metrics
	metrics.MemoryUsage = MemoryUsage{
		StartMemory:    bytesToMB(startMemStats.Alloc),
		EndMemory:      bytesToMB(endMemStats.Alloc),
		PeakMemory:     bytesToMB(endMemStats.Sys),
		MemoryDelta:    int64(bytesToMB(endMemStats.Alloc)) - int64(bytesToMB(startMemStats.Alloc)),
		GCCount:        endMemStats.NumGC - startMemStats.NumGC,
		AllocatedBytes: endMemStats.TotalAlloc - startMemStats.TotalAlloc,
	}

	// Measure file system metrics
	fsMetrics, err := pt.measureFileSystemMetrics(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to measure filesystem metrics: %w", err)
	}
	metrics.FileSystemMetrics = *fsMetrics

	return metrics, nil
}

// MeasureValidationPerformance measures validation performance across different project types
func (pt *PerformanceTester) MeasureValidationPerformance(projectPaths []string) (map[string]*PerformanceMetrics, error) {
	results := make(map[string]*PerformanceMetrics)

	for _, projectPath := range projectPaths {
		projectName := filepath.Base(projectPath)

		var memStats runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memStats)
		startMemory := memStats.Alloc

		startTime := time.Now()
		_, err := pt.validationEngine.ValidateProject(projectPath)
		validationTime := time.Since(startTime)

		runtime.ReadMemStats(&memStats)
		endMemory := memStats.Alloc

		if err != nil {
			return nil, fmt.Errorf("validation failed for %s: %w", projectName, err)
		}

		fsMetrics, err := pt.measureFileSystemMetrics(projectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to measure filesystem metrics for %s: %w", projectName, err)
		}

		results[projectName] = &PerformanceMetrics{
			ValidationTime: validationTime,
			TotalTime:      validationTime,
			MemoryUsage: MemoryUsage{
				StartMemory: bytesToMB(startMemory),
				EndMemory:   bytesToMB(endMemory),
				MemoryDelta: int64(bytesToMB(endMemory)) - int64(bytesToMB(startMemory)),
			},
			FileSystemMetrics: *fsMetrics,
		}
	}

	return results, nil
}

// MeasureScalabilityPerformance measures performance with different project sizes
func (pt *PerformanceTester) MeasureScalabilityPerformance(baseProjectPath string, fileCounts []int) (map[int]*PerformanceMetrics, error) {
	results := make(map[int]*PerformanceMetrics)

	for _, fileCount := range fileCounts {
		// Create a temporary project with specified number of files
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("perf-test-%d-files-*", fileCount))
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Generate test files
		if generateErr := pt.generateTestFiles(tempDir, fileCount); generateErr != nil {
			return nil, fmt.Errorf("failed to generate test files: %w", generateErr)
		}

		var memStats runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memStats)
		startMemory := memStats.Alloc

		startTime := time.Now()
		_, err = pt.validationEngine.ValidateProject(tempDir)
		validationTime := time.Since(startTime)

		runtime.ReadMemStats(&memStats)
		endMemory := memStats.Alloc

		if err != nil {
			return nil, fmt.Errorf("validation failed for %d files: %w", fileCount, err)
		}

		fsMetrics, err := pt.measureFileSystemMetrics(tempDir)
		if err != nil {
			return nil, fmt.Errorf("failed to measure filesystem metrics for %d files: %w", fileCount, err)
		}

		results[fileCount] = &PerformanceMetrics{
			ValidationTime: validationTime,
			TotalTime:      validationTime,
			MemoryUsage: MemoryUsage{
				StartMemory: bytesToMB(startMemory),
				EndMemory:   bytesToMB(endMemory),
				MemoryDelta: int64(bytesToMB(endMemory)) - int64(bytesToMB(startMemory)),
			},
			FileSystemMetrics: *fsMetrics,
		}
	}

	return results, nil
}

// measureFileSystemMetrics calculates file system statistics for a project
func (pt *PerformanceTester) measureFileSystemMetrics(projectPath string) (*FileSystemMetrics, error) {
	metrics := &FileSystemMetrics{
		SmallestFile: int64(^uint64(0) >> 1), // Max int64
	}

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			metrics.DirectoriesCreated++
		} else {
			metrics.FilesCreated++
			size := info.Size()
			metrics.TotalSize += size

			if size > metrics.LargestFile {
				metrics.LargestFile = size
			}
			if size < metrics.SmallestFile {
				metrics.SmallestFile = size
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// If no files were found, reset smallest file to 0
	if metrics.FilesCreated == 0 {
		metrics.SmallestFile = 0
	}

	return metrics, nil
}

// generateTestFiles creates test files for scalability testing
func (pt *PerformanceTester) generateTestFiles(projectPath string, fileCount int) error {
	// Create basic project structure
	dirs := []string{
		"frontend/src/components",
		"frontend/src/pages",
		"backend/internal/handlers",
		"backend/internal/models",
		"mobile/android/app/src/main",
		"mobile/ios/App",
		"deploy/k8s",
		"docs",
		".github/workflows",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	// Create basic required files
	requiredFiles := map[string]string{
		"README.md": "# Test Project\n\nThis is a test project for performance testing.",
		"Makefile":  "all:\n\techo 'Building project...'\n",
		"frontend/package.json": `{
			"name": "test-frontend",
			"version": "1.0.0",
			"scripts": {
				"dev": "next dev",
				"build": "next build"
			}
		}`,
		"backend/go.mod": "module test-backend\n\ngo 1.24\n",
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
	}

	filesCreated := 0
	for filePath, content := range requiredFiles {
		fullPath := filepath.Join(projectPath, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
		filesCreated++
	}

	// Generate additional files to reach target count
	for filesCreated < fileCount {
		// Generate different types of files
		fileTypes := []struct {
			dir       string
			extension string
			content   string
		}{
			{"frontend/src/components", ".tsx", "import React from 'react';\n\nexport default function Component() {\n  return <div>Component</div>;\n}"},
			{"frontend/src/pages", ".tsx", "import React from 'react';\n\nexport default function Page() {\n  return <div>Page</div>;\n}"},
			{"backend/internal/handlers", ".go", "package handlers\n\nimport \"net/http\"\n\nfunc Handler(w http.ResponseWriter, r *http.Request) {\n\tw.Write([]byte(\"Hello\"))\n}"},
			{"backend/internal/models", ".go", "package models\n\ntype Model struct {\n\tID   int    `json:\"id\"`\n\tName string `json:\"name\"`\n}"},
			{"deploy/k8s", ".yaml", "apiVersion: v1\nkind: Service\nmetadata:\n  name: test-service\nspec:\n  selector:\n    app: test\n  ports:\n  - port: 80"},
			{"docs", ".md", "# Documentation\n\nThis is documentation for the test project."},
		}

		for _, fileType := range fileTypes {
			if filesCreated >= fileCount {
				break
			}

			fileName := fmt.Sprintf("file_%d%s", filesCreated, fileType.extension)
			filePath := filepath.Join(projectPath, fileType.dir, fileName)

			if err := os.WriteFile(filePath, []byte(fileType.content), 0644); err != nil {
				return err
			}
			filesCreated++
		}
	}

	return nil
}

// GeneratePerformanceReport creates a human-readable performance report
func (pt *PerformanceTester) GeneratePerformanceReport(metrics *PerformanceMetrics) string {
	report := fmt.Sprintf(`Performance Report
==================

Timing Metrics:
- Generation Time: %v
- Validation Time: %v
- Setup Time: %v
- Verification Time: %v
- Total Time: %v

Memory Usage:
- Start Memory: %d MB
- Peak Memory: %d MB
- End Memory: %d MB
- Memory Delta: %+d MB
- GC Count: %d
- Total Allocated: %d bytes

File System Metrics:
- Files Created: %d
- Directories Created: %d
- Total Size: %d bytes (%.2f MB)
- Largest File: %d bytes
- Smallest File: %d bytes

Performance Summary:
- Files per second: %.2f
- MB processed per second: %.2f
- Memory efficiency: %.2f MB per file
`,
		metrics.GenerationTime,
		metrics.ValidationTime,
		metrics.SetupTime,
		metrics.VerificationTime,
		metrics.TotalTime,
		metrics.MemoryUsage.StartMemory,
		metrics.MemoryUsage.PeakMemory,
		metrics.MemoryUsage.EndMemory,
		metrics.MemoryUsage.MemoryDelta,
		metrics.MemoryUsage.GCCount,
		metrics.MemoryUsage.AllocatedBytes,
		metrics.FileSystemMetrics.FilesCreated,
		metrics.FileSystemMetrics.DirectoriesCreated,
		metrics.FileSystemMetrics.TotalSize,
		float64(metrics.FileSystemMetrics.TotalSize)/(1024*1024),
		metrics.FileSystemMetrics.LargestFile,
		metrics.FileSystemMetrics.SmallestFile,
		float64(metrics.FileSystemMetrics.FilesCreated)/metrics.TotalTime.Seconds(),
		float64(metrics.FileSystemMetrics.TotalSize)/(1024*1024)/metrics.TotalTime.Seconds(),
		float64(metrics.MemoryUsage.PeakMemory)/float64(metrics.FileSystemMetrics.FilesCreated),
	)

	return report
}

// bytesToMB converts bytes to megabytes
func bytesToMB(bytes uint64) uint64 {
	return bytes / (1024 * 1024)
}
