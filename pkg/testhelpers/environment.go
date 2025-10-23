package testhelpers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestEnvironment provides utilities for setting up test environments
type TestEnvironment struct {
	t            *testing.T
	TempDir      string
	MockTools    map[string]*MockTool
	OriginalPATH string
	MockBinDir   string
}

// MockTool represents a mocked external tool
type MockTool struct {
	Name       string
	Version    string
	Available  bool
	ExitCode   int
	Stdout     string
	Stderr     string
	ScriptPath string
}

// SetupTestEnv creates a new test environment with temporary directories
func SetupTestEnv(t *testing.T) *TestEnvironment {
	t.Helper()

	tempDir := t.TempDir()
	mockBinDir := filepath.Join(tempDir, "mock-bin")

	if err := os.MkdirAll(mockBinDir, 0755); err != nil {
		t.Fatalf("Failed to create mock bin directory: %v", err)
	}

	env := &TestEnvironment{
		t:            t,
		TempDir:      tempDir,
		MockTools:    make(map[string]*MockTool),
		OriginalPATH: os.Getenv("PATH"),
		MockBinDir:   mockBinDir,
	}

	return env
}

// Cleanup removes temporary directories and restores environment
func (te *TestEnvironment) Cleanup() {
	te.t.Helper()

	// Restore original PATH
	if te.OriginalPATH != "" {
		os.Setenv("PATH", te.OriginalPATH)
	}
}

// MockToolAvailable creates a mock tool that appears available
func (te *TestEnvironment) MockToolAvailable(name string, version string) error {
	te.t.Helper()

	mockTool := &MockTool{
		Name:      name,
		Version:   version,
		Available: true,
		ExitCode:  0,
		Stdout:    version,
	}

	// Create mock script
	scriptPath := filepath.Join(te.MockBinDir, name)

	// Create a simple shell script that outputs the version
	script := fmt.Sprintf(`#!/bin/sh
if [ "$1" = "--version" ] || [ "$1" = "version" ]; then
    echo "%s"
    exit 0
fi
echo "Mock %s executed"
exit 0
`, version, name)

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create mock script: %w", err)
	}

	mockTool.ScriptPath = scriptPath
	te.MockTools[name] = mockTool

	// Update PATH to include mock bin directory
	newPATH := te.MockBinDir + string(os.PathListSeparator) + os.Getenv("PATH")
	_ = os.Setenv("PATH", newPATH) // Ignore error in test helper

	return nil
}

// MockToolUnavailable ensures a tool is not available
func (te *TestEnvironment) MockToolUnavailable(name string) {
	te.t.Helper()

	mockTool := &MockTool{
		Name:      name,
		Available: false,
	}

	te.MockTools[name] = mockTool

	// Remove from mock bin if it exists
	scriptPath := filepath.Join(te.MockBinDir, name)
	_ = os.Remove(scriptPath) // Ignore error if file doesn't exist
}

// MockToolWithBehavior creates a mock tool with custom behavior
func (te *TestEnvironment) MockToolWithBehavior(name string, exitCode int, stdout, stderr string) error {
	te.t.Helper()

	mockTool := &MockTool{
		Name:      name,
		Available: true,
		ExitCode:  exitCode,
		Stdout:    stdout,
		Stderr:    stderr,
	}

	scriptPath := filepath.Join(te.MockBinDir, name)

	script := fmt.Sprintf(`#!/bin/sh
echo "%s"
echo "%s" >&2
exit %d
`, stdout, stderr, exitCode)

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create mock script: %w", err)
	}

	mockTool.ScriptPath = scriptPath
	te.MockTools[name] = mockTool

	// Update PATH
	newPATH := te.MockBinDir + string(os.PathListSeparator) + os.Getenv("PATH")
	_ = os.Setenv("PATH", newPATH) // Ignore error in test helper

	return nil
}

// CreateProjectStructure creates a basic project structure for testing
func (te *TestEnvironment) CreateProjectStructure(projectName string) (string, error) {
	te.t.Helper()

	projectDir := filepath.Join(te.TempDir, projectName)

	dirs := []string{
		filepath.Join(projectDir, "App"),
		filepath.Join(projectDir, "CommonServer"),
		filepath.Join(projectDir, "Mobile", "android"),
		filepath.Join(projectDir, "Mobile", "ios"),
		filepath.Join(projectDir, "Deploy"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return projectDir, nil
}

// CreateTestConfig creates a test project configuration
func (te *TestEnvironment) CreateTestConfig(projectName string, components []string) *models.ProjectConfig {
	te.t.Helper()

	config := &models.ProjectConfig{
		Name:        projectName,
		Description: "Test project",
		OutputDir:   filepath.Join(te.TempDir, projectName),
		Components:  make([]models.ComponentConfig, 0, len(components)),
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "http://localhost:8080",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     false,
		},
	}

	for _, compType := range components {
		comp := models.ComponentConfig{
			Type:    compType,
			Name:    fmt.Sprintf("%s-component", compType),
			Enabled: true,
			Config:  make(map[string]interface{}),
		}

		// Add component-specific config
		switch compType {
		case "nextjs":
			comp.Config["typescript"] = true
			comp.Config["tailwind"] = true
		case "go-backend":
			comp.Config["module"] = fmt.Sprintf("github.com/test/%s", projectName)
			comp.Config["framework"] = "gin"
		case "android":
			comp.Config["package"] = "com.test.app"
		case "ios":
			comp.Config["bundle_id"] = "com.test.app"
		}

		config.Components = append(config.Components, comp)
	}

	return config
}

// IsToolAvailable checks if a real tool is available on the system
func IsToolAvailable(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

// SkipIfToolNotAvailable skips the test if the specified tool is not available
func SkipIfToolNotAvailable(t *testing.T, toolName string) {
	t.Helper()

	if !IsToolAvailable(toolName) {
		t.Skipf("Skipping test: %s not available", toolName)
	}
}

// CreateMockNextJSProject creates a minimal Next.js project structure for testing
func (te *TestEnvironment) CreateMockNextJSProject(dir string) error {
	te.t.Helper()

	files := map[string]string{
		"package.json": `{
  "name": "test-app",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "16.0.0",
    "react": "^19.2.0",
    "react-dom": "^19.2.0"
  }
}`,
		"next.config.js": `/** @type {import('next').NextConfig} */
const nextConfig = {}
module.exports = nextConfig`,
		"app/page.tsx": `export default function Home() {
  return <main>Hello World</main>
}`,
		"app/layout.tsx": `export default function RootLayout({ children }: { children: React.ReactNode }) {
  return <html><body>{children}</body></html>
}`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(dir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// CreateMockGoProject creates a minimal Go project structure for testing
func (te *TestEnvironment) CreateMockGoProject(dir string, moduleName string) error {
	te.t.Helper()

	files := map[string]string{
		"go.mod": fmt.Sprintf(`module %s

go 1.25

require github.com/gin-gonic/gin v1.11.0
`, moduleName),
		"main.go": `package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})
	r.Run(":8080")
}`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(dir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// CreateMockAndroidProject creates a minimal Android project structure for testing
func (te *TestEnvironment) CreateMockAndroidProject(dir string) error {
	te.t.Helper()

	files := map[string]string{
		"settings.gradle": `rootProject.name = "TestApp"
include ':app'`,
		"build.gradle": `buildscript {
    repositories {
        google()
        mavenCentral()
    }
}`,
		"app/build.gradle": `plugins {
    id 'com.android.application'
    id 'org.jetbrains.kotlin.android'
}

android {
    namespace 'com.test.app'
    compileSdk 36
}`,
		"app/src/main/AndroidManifest.xml": `<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    <application android:label="TestApp">
    </application>
</manifest>`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(dir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// CreateMockIOSProject creates a minimal iOS project structure for testing
func (te *TestEnvironment) CreateMockIOSProject(dir string, appName string) error {
	te.t.Helper()

	files := map[string]string{
		"Package.swift": fmt.Sprintf(`// swift-tools-version:6.2
import PackageDescription

let package = Package(
    name: "%s",
    platforms: [.iOS(.v16)],
    products: [
        .library(name: "%s", targets: ["%s"])
    ],
    targets: [
        .target(name: "%s")
    ]
)`, appName, appName, appName, appName),
		fmt.Sprintf("%s/%sApp.swift", appName, appName): fmt.Sprintf(`import SwiftUI

@main
struct %sApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}`, appName),
		fmt.Sprintf("%s/ContentView.swift", appName): `import SwiftUI

struct ContentView: View {
    var body: some View {
        Text("Hello, World!")
    }
}`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(dir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteWithContext executes a command with context for testing
func ExecuteWithContext(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// AssertFileExists checks if a file exists and fails the test if it doesn't
func AssertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file does not exist: %s", path)
	}
}

// AssertFileContains checks if a file contains the expected content
func AssertFileContains(t *testing.T, path string, expected string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	if !contains(string(content), expected) {
		t.Errorf("File %s does not contain expected content: %s", path, expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
