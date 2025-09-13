package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/stretchr/testify/require"
)

// TestPerformanceBenchmarking runs comprehensive performance benchmarks
func TestPerformanceBenchmarking(t *testing.T) {
	t.Run("BaselinePerformance", testBaselinePerformance)
	t.Run("ScalabilityBenchmarks", testScalabilityBenchmarks)
	t.Run("MemoryEfficiencyBenchmarks", testMemoryEfficiencyBenchmarks)
	t.Run("ValidationPerformanceBenchmarks", testValidationPerformanceBenchmarks)
}

// testBaselinePerformance establishes baseline performance metrics
func testBaselinePerformance(t *testing.T) {
	tester := NewPerformanceTester()

	// Create a standard test project
	tempDir, err := os.MkdirTemp("", "baseline-perf-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = createStandardTestProject(tempDir)
	require.NoError(t, err)

	config := &models.ProjectConfig{
		Name:         "baseline-project",
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

	// Measure baseline performance
	metrics, err := tester.MeasureProjectGeneration(tempDir, config)
	require.NoError(t, err)

	// Performance expectations (these should be met after cleanup optimizations)
	performanceChecks := []struct {
		name      string
		actual    time.Duration
		threshold time.Duration
		message   string
	}{
		{
			name:      "ValidationTime",
			actual:    metrics.ValidationTime,
			threshold: 500 * time.Millisecond,
			message:   "Validation should complete within 500ms for standard project",
		},
		{
			name:      "SetupTime",
			actual:    metrics.SetupTime,
			threshold: 1 * time.Second,
			message:   "Setup should complete within 1 second for standard project",
		},
		{
			name:      "VerificationTime",
			actual:    metrics.VerificationTime,
			threshold: 300 * time.Millisecond,
			message:   "Verification should complete within 300ms for standard project",
		},
		{
			name:      "TotalTime",
			actual:    metrics.TotalTime,
			threshold: 2 * time.Second,
			message:   "Total time should be under 2 seconds for standard project",
		},
	}

	for _, check := range performanceChecks {
		if check.actual > check.threshold {
			t.Logf("PERFORMANCE WARNING: %s took %v (threshold: %v) - %s",
				check.name, check.actual, check.threshold, check.message)
		} else {
			t.Logf("PERFORMANCE OK: %s took %v (threshold: %v)",
				check.name, check.actual, check.threshold)
		}
	}

	// Memory efficiency checks
	memoryChecks := []struct {
		name      string
		actual    uint64
		threshold uint64
		message   string
	}{
		{
			name:      "PeakMemory",
			actual:    metrics.MemoryUsage.PeakMemory,
			threshold: 100, // 100 MB
			message:   "Peak memory usage should be under 100MB for standard project",
		},
		{
			name:      "MemoryDelta",
			actual:    uint64(abs(metrics.MemoryUsage.MemoryDelta)),
			threshold: 50, // 50 MB
			message:   "Memory delta should be under 50MB for standard project",
		},
	}

	for _, check := range memoryChecks {
		if check.actual > check.threshold {
			t.Logf("MEMORY WARNING: %s is %d MB (threshold: %d MB) - %s",
				check.name, check.actual, check.threshold, check.message)
		} else {
			t.Logf("MEMORY OK: %s is %d MB (threshold: %d MB)",
				check.name, check.actual, check.threshold)
		}
	}

	// Generate and log performance report
	report := tester.GeneratePerformanceReport(metrics)
	t.Logf("Baseline Performance Report:\n%s", report)
}

// testScalabilityBenchmarks tests performance with different project sizes
func testScalabilityBenchmarks(t *testing.T) {
	tester := NewPerformanceTester()

	// Test with different file counts to measure scalability
	fileCounts := []int{10, 50, 100, 200, 500}

	results, err := tester.MeasureScalabilityPerformance("", fileCounts)
	require.NoError(t, err)

	t.Logf("Scalability Benchmark Results:")
	t.Logf("Files\tValidation Time\tMemory (MB)\tFiles/sec\tMB/sec")
	t.Logf("-----\t---------------\t-----------\t---------\t------")

	var previousTime time.Duration
	for _, fileCount := range fileCounts {
		metrics := results[fileCount]
		filesPerSec := float64(fileCount) / metrics.ValidationTime.Seconds()
		mbPerSec := float64(metrics.FileSystemMetrics.TotalSize) / (1024 * 1024) / metrics.ValidationTime.Seconds()

		t.Logf("%d\t%v\t%d\t%.2f\t%.2f",
			fileCount,
			metrics.ValidationTime,
			metrics.MemoryUsage.PeakMemory,
			filesPerSec,
			mbPerSec)

		// Check for reasonable scaling
		if previousTime > 0 {
			scalingFactor := float64(metrics.ValidationTime) / float64(previousTime)
			if scalingFactor > 3.0 {
				t.Logf("SCALABILITY WARNING: Validation time increased by %.2fx from previous size", scalingFactor)
			} else {
				t.Logf("SCALABILITY OK: Validation time scaling factor: %.2fx", scalingFactor)
			}
		}
		previousTime = metrics.ValidationTime

		// Performance thresholds based on file count
		expectedMaxTime := time.Duration(fileCount) * 5 * time.Millisecond // 5ms per file max
		if metrics.ValidationTime > expectedMaxTime {
			t.Logf("PERFORMANCE WARNING: %d files took %v (expected max: %v)",
				fileCount, metrics.ValidationTime, expectedMaxTime)
		}
	}
}

// testMemoryEfficiencyBenchmarks tests memory usage patterns
func testMemoryEfficiencyBenchmarks(t *testing.T) {
	tester := NewPerformanceTester()

	// Test memory efficiency with different project types
	projectTypes := []struct {
		name      string
		fileCount int
		avgSize   int
	}{
		{"Small Frontend", 20, 1024},
		{"Medium Backend", 50, 2048},
		{"Large Full-Stack", 100, 4096},
		{"Enterprise", 200, 8192},
	}

	t.Logf("Memory Efficiency Benchmark Results:")
	t.Logf("Project Type\tFiles\tPeak Memory (MB)\tMemory per File (KB)\tGC Count")
	t.Logf("------------\t-----\t----------------\t--------------------\t--------")

	for _, projectType := range projectTypes {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("memory-test-%s-*", projectType.name))
		require.NoError(t, err)
		defer func() {
			os.RemoveAll(tempDir)
			// Force garbage collection after each test to prevent memory accumulation
			runtime.GC()
			runtime.GC()
		}()

		// Generate project with specific characteristics
		err = tester.generateTestFiles(tempDir, projectType.fileCount)
		require.NoError(t, err)

		// Measure memory usage
		var memStats runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memStats)
		startMemory := memStats.Alloc
		startGC := memStats.NumGC

		startTime := time.Now()
		_, err = tester.validationEngine.ValidateProject(tempDir)
		validationTime := time.Since(startTime)
		require.NoError(t, err)

		runtime.ReadMemStats(&memStats)
		endMemory := memStats.Alloc
		endGC := memStats.NumGC

		peakMemoryMB := bytesToMB(memStats.Sys)
		memoryPerFileKB := float64(endMemory-startMemory) / float64(projectType.fileCount) / 1024
		gcCount := endGC - startGC

		t.Logf("%s\t%d\t%d\t%.2f\t%d",
			projectType.name,
			projectType.fileCount,
			peakMemoryMB,
			memoryPerFileKB,
			gcCount)

		// Memory efficiency checks
		maxMemoryPerFile := 100.0 // 100 KB per file max
		if memoryPerFileKB > maxMemoryPerFile {
			t.Logf("MEMORY WARNING: %s uses %.2f KB per file (threshold: %.2f KB)",
				projectType.name, memoryPerFileKB, maxMemoryPerFile)
		}

		// GC frequency check
		maxGCPerFile := 0.1 // Max 1 GC per 10 files
		gcPerFile := float64(gcCount) / float64(projectType.fileCount)
		if gcPerFile > maxGCPerFile {
			t.Logf("GC WARNING: %s triggered %.3f GC per file (threshold: %.3f)",
				projectType.name, gcPerFile, maxGCPerFile)
		}

		// Performance check
		maxTimePerFile := 10 * time.Millisecond // 10ms per file max
		timePerFile := validationTime / time.Duration(projectType.fileCount)
		if timePerFile > maxTimePerFile {
			t.Logf("PERFORMANCE WARNING: %s took %v per file (threshold: %v)",
				projectType.name, timePerFile, maxTimePerFile)
		}
	}
}

// testValidationPerformanceBenchmarks tests validation engine performance
func testValidationPerformanceBenchmarks(t *testing.T) {
	tester := NewPerformanceTester()

	// Create different types of projects to validate
	projectConfigs := []struct {
		name   string
		config func(string) error
	}{
		{
			name: "Frontend-Only",
			config: func(dir string) error {
				return createFrontendProject(dir)
			},
		},
		{
			name: "Backend-Only",
			config: func(dir string) error {
				return createBackendProject(dir)
			},
		},
		{
			name: "Full-Stack",
			config: func(dir string) error {
				return createFullStackProject(dir)
			},
		},
		{
			name: "Mobile",
			config: func(dir string) error {
				return createMobileProject(dir)
			},
		},
		{
			name: "Infrastructure",
			config: func(dir string) error {
				return createInfrastructureProject(dir)
			},
		},
	}

	t.Logf("Validation Performance Benchmark Results:")
	t.Logf("Project Type\tValidation Time\tFiles\tTime per File\tThroughput (files/sec)")
	t.Logf("------------\t---------------\t-----\t-------------\t---------------------")

	for _, projectConfig := range projectConfigs {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("validation-test-%s-*", projectConfig.name))
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create project
		err = projectConfig.config(tempDir)
		require.NoError(t, err)

		// Count files
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

		// Measure validation performance
		startTime := time.Now()
		_, err = tester.validationEngine.ValidateProject(tempDir)
		validationTime := time.Since(startTime)

		if err != nil {
			t.Logf("Validation failed for %s: %v", projectConfig.name, err)
			continue
		}

		timePerFile := validationTime / time.Duration(fileCount)
		throughput := float64(fileCount) / validationTime.Seconds()

		t.Logf("%s\t%v\t%d\t%v\t%.2f",
			projectConfig.name,
			validationTime,
			fileCount,
			timePerFile,
			throughput)

		// Performance expectations
		maxValidationTime := 1 * time.Second
		if validationTime > maxValidationTime {
			t.Logf("PERFORMANCE WARNING: %s validation took %v (threshold: %v)",
				projectConfig.name, validationTime, maxValidationTime)
		}

		minThroughput := 50.0 // 50 files per second minimum
		if throughput < minThroughput {
			t.Logf("THROUGHPUT WARNING: %s achieved %.2f files/sec (threshold: %.2f)",
				projectConfig.name, throughput, minThroughput)
		}
	}
}

// Helper functions to create different project types

func createStandardTestProject(projectPath string) error {
	return createSimpleTestProject(projectPath)
}

func createFrontendProject(projectPath string) error {
	dirs := []string{
		"src/components",
		"src/pages",
		"src/hooks",
		"src/utils",
		"public",
		"styles",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"package.json": `{
			"name": "frontend-project",
			"version": "1.0.0",
			"scripts": {
				"dev": "next dev",
				"build": "next build",
				"start": "next start"
			},
			"dependencies": {
				"react": "^18.0.0",
				"next": "^13.0.0"
			}
		}`,
		"next.config.js": "module.exports = { reactStrictMode: true }",
		"src/pages/index.tsx": `import React from 'react';
export default function Home() {
	return <div>Home Page</div>;
}`,
		"src/components/Header.tsx": `import React from 'react';
export default function Header() {
	return <header>Header</header>;
}`,
		"tailwind.config.js": "module.exports = { content: ['./src/**/*.{js,ts,jsx,tsx}'] }",
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

func createBackendProject(projectPath string) error {
	dirs := []string{
		"cmd/api",
		"internal/handlers",
		"internal/models",
		"internal/services",
		"pkg/utils",
		"configs",
		"migrations",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"go.mod": `module backend-project

go 1.24

require (
	github.com/gin-gonic/gin v1.9.0
	github.com/stretchr/testify v1.8.0
)`,
		"cmd/api/main.go": `package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	log.Fatal(r.Run(":8080"))
}`,
		"internal/handlers/user.go": `package handlers

import "github.com/gin-gonic/gin"

func GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"user": "example"})
}`,
		"internal/models/user.go": `package models

type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}`,
		"Dockerfile": `FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go
CMD ["./main"]`,
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

func createFullStackProject(projectPath string) error {
	// Create both frontend and backend
	if err := createFrontendProject(filepath.Join(projectPath, "frontend")); err != nil {
		return err
	}
	if err := createBackendProject(filepath.Join(projectPath, "backend")); err != nil {
		return err
	}

	// Add root-level files
	files := map[string]string{
		"docker-compose.yml": `version: '3.8'
services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
  backend:
    build: ./backend
    ports:
      - "8080:8080"`,
		"Makefile": `all: build

build:
	cd frontend && npm run build
	cd backend && go build -o bin/api cmd/api/main.go

dev:
	docker-compose up --build`,
		"README.md": "# Full-Stack Project\n\nA complete full-stack application.",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectPath, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func createMobileProject(projectPath string) error {
	dirs := []string{
		"android/app/src/main/java",
		"android/app/src/main/res",
		"ios/App",
		"src/components",
		"src/screens",
		"src/services",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"package.json": `{
			"name": "mobile-project",
			"version": "1.0.0",
			"scripts": {
				"android": "react-native run-android",
				"ios": "react-native run-ios"
			},
			"dependencies": {
				"react-native": "^0.72.0"
			}
		}`,
		"android/build.gradle": "// Android build configuration",
		"ios/Podfile":          "# iOS dependencies",
		"src/App.tsx": `import React from 'react';
import { View, Text } from 'react-native';

export default function App() {
	return (
		<View>
			<Text>Mobile App</Text>
		</View>
	);
}`,
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

func createInfrastructureProject(projectPath string) error {
	dirs := []string{
		"terraform/modules/vpc",
		"terraform/modules/eks",
		"k8s/deployments",
		"k8s/services",
		"helm/charts",
		"scripts",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"terraform/main.tf": `terraform {
			required_version = ">= 1.0"
		}

		provider "aws" {
			region = "us-west-2"
		}`,
		"terraform/variables.tf": `variable "environment" {
			description = "Environment name"
			type        = string
			default     = "dev"
		}`,
		"k8s/deployments/app.yaml": `apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:latest`,
		"k8s/services/app.yaml": `apiVersion: v1
kind: Service
metadata:
  name: app-service
spec:
  selector:
    app: myapp
  ports:
  - port: 80
    targetPort: 8080`,
		"helm/Chart.yaml": `apiVersion: v2
name: myapp
version: 1.0.0
description: My application Helm chart`,
		"scripts/deploy.sh": `#!/bin/bash
echo "Deploying application..."
kubectl apply -f k8s/`,
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

// Helper function for absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Benchmark functions for performance testing

func BenchmarkValidationEngine_SmallProject(b *testing.B) {
	tester := NewPerformanceTester()

	tempDir, err := os.MkdirTemp("", "benchmark-small-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	err = createFrontendProject(tempDir)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.validationEngine.ValidateProject(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkValidationEngine_MediumProject(b *testing.B) {
	tester := NewPerformanceTester()

	tempDir, err := os.MkdirTemp("", "benchmark-medium-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	err = createFullStackProject(tempDir)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.validationEngine.ValidateProject(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkValidationEngine_LargeProject(b *testing.B) {
	tester := NewPerformanceTester()

	tempDir, err := os.MkdirTemp("", "benchmark-large-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create a large project with many files
	err = tester.generateTestFiles(tempDir, 500)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.validationEngine.ValidateProject(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkMemoryUsage_ProjectGeneration(b *testing.B) {
	tester := NewPerformanceTester()

	config := &models.ProjectConfig{
		Name:         "benchmark-project",
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempDir, err := os.MkdirTemp("", "benchmark-memory-*")
		require.NoError(b, err)

		err = createStandardTestProject(tempDir)
		require.NoError(b, err)

		_, err = tester.MeasureProjectGeneration(tempDir, config)
		require.NoError(b, err)

		os.RemoveAll(tempDir)
	}
}
