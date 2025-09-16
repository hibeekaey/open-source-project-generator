// Package version provides basic package version management capabilities
// for the Open Source Template Generator.
//
// This package implements the VersionManager interface and provides:
//   - Basic version fetching from NPM, Go modules, and GitHub releases
//   - Simple version parsing and validation
//   - Integration with template processing for automatic injection
//
// Supported Registries:
//   - NPM Registry: Node.js packages and frameworks (React, Next.js, etc.)
//   - Go Module Registry: Go packages and modules
//   - GitHub Releases: GitHub-hosted projects and tools
//
// Usage:
//
//	manager := version.NewManager()
//
//	// Get latest version
//	nodeVersion, err := manager.GetLatestNodeVersion()
//	if err != nil {
//	    log.Fatal(err)
//	}
package version

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// Manager implements the VersionManager interface
type Manager struct {
	httpClient   *http.Client
	npmClient    *NPMClient
	goClient     *GoClient
	githubClient *GitHubClient
}

// NewManager creates a new version manager with all clients
func NewManager() *Manager {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Manager{
		httpClient:   httpClient,
		npmClient:    NewNPMClient(httpClient),
		goClient:     NewGoClient(httpClient),
		githubClient: NewGitHubClient(httpClient),
	}
}

// GetLatestNodeVersion fetches the latest Node.js version
func (m *Manager) GetLatestNodeVersion() (string, error) {
	version, err := m.githubClient.GetLatestRelease("nodejs", "node")
	if err != nil {
		return "", fmt.Errorf("failed to fetch Node.js version: %w", err)
	}

	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	return version, nil
}

// GetLatestGoVersion fetches the latest Go version
func (m *Manager) GetLatestGoVersion() (string, error) {
	version, err := m.githubClient.GetLatestRelease("golang", "go")
	if err != nil {
		return "", fmt.Errorf("failed to fetch Go version: %w", err)
	}

	// Remove 'go' prefix if present (e.g., go1.22.0 -> 1.22.0)
	version = strings.TrimPrefix(version, "go")

	return version, nil
}

// GetLatestNPMPackage fetches the latest version of an NPM package
func (m *Manager) GetLatestNPMPackage(packageName string) (string, error) {
	version, err := m.npmClient.GetLatestVersion(packageName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch NPM package %s: %w", packageName, err)
	}

	return version, nil
}

// GetLatestGoModule fetches the latest version of a Go module
func (m *Manager) GetLatestGoModule(moduleName string) (string, error) {
	version, err := m.goClient.GetLatestVersion(moduleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Go module %s: %w", moduleName, err)
	}

	return version, nil
}

// GetLatestGitHubRelease fetches the latest release version from GitHub
func (m *Manager) GetLatestGitHubRelease(owner, repo string) (string, error) {
	version, err := m.githubClient.GetLatestRelease(owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitHub release %s/%s: %w", owner, repo, err)
	}

	return version, nil
}

// UpdateVersionsConfig creates a version configuration with latest versions
func (m *Manager) UpdateVersionsConfig() (*models.VersionConfig, error) {
	config := &models.VersionConfig{
		Packages: make(map[string]string),
	}

	// Get Node.js version
	if nodeVersion, err := m.GetLatestNodeVersion(); err == nil {
		config.Node = nodeVersion
	}

	// Get Go version
	if goVersion, err := m.GetLatestGoVersion(); err == nil {
		config.Go = goVersion
	}

	// Get common NPM packages
	commonPackages := []string{
		"react", "next", "typescript", "tailwindcss", "eslint", "prettier",
		"@types/node", "@types/react", "autoprefixer", "postcss",
	}

	for _, pkg := range commonPackages {
		if version, err := m.GetLatestNPMPackage(pkg); err == nil {
			config.Packages[pkg] = version
		}
	}

	return config, nil
}

// GetVersionHistory returns version history for a package (simplified implementation)
func (m *Manager) GetVersionHistory(packageName string) ([]string, error) {
	// For simplicity, just return the latest version
	version, err := m.GetLatestNPMPackage(packageName)
	if err != nil {
		return nil, err
	}
	return []string{version}, nil
}
