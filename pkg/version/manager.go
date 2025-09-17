// Package version provides basic package version management capabilities
// for the Open Source Project Generator.
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

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
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

// Enhanced version management methods

// GetCurrentVersion gets the current version
func (m *Manager) GetCurrentVersion() string {
	return "1.0.0" // Placeholder
}

// GetLatestVersion gets the latest version info
func (m *Manager) GetLatestVersion() (*interfaces.VersionInfo, error) {
	return nil, fmt.Errorf("GetLatestVersion implementation pending - will be implemented in task 8")
}

// GetAllPackageVersions gets all package versions
func (m *Manager) GetAllPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("GetAllPackageVersions implementation pending - will be implemented in task 8")
}

// GetLatestPackageVersions gets latest package versions
func (m *Manager) GetLatestPackageVersions() (map[string]string, error) {
	return nil, fmt.Errorf("GetLatestPackageVersions implementation pending - will be implemented in task 8")
}

// GetDetailedVersionHistory gets detailed version history
func (m *Manager) GetDetailedVersionHistory() ([]interfaces.VersionInfo, error) {
	return nil, fmt.Errorf("GetDetailedVersionHistory implementation pending - will be implemented in task 8")
}

// CheckForUpdates checks for updates
func (m *Manager) CheckForUpdates() (*interfaces.UpdateInfo, error) {
	return nil, fmt.Errorf("CheckForUpdates implementation pending - will be implemented in task 8")
}

// DownloadUpdate downloads an update
func (m *Manager) DownloadUpdate(version string) error {
	return fmt.Errorf("DownloadUpdate implementation pending - will be implemented in task 8")
}

// InstallUpdate installs an update
func (m *Manager) InstallUpdate(version string) error {
	return fmt.Errorf("InstallUpdate implementation pending - will be implemented in task 8")
}

// GetUpdateChannel gets the update channel
func (m *Manager) GetUpdateChannel() string {
	return "stable"
}

// SetUpdateChannel sets the update channel
func (m *Manager) SetUpdateChannel(channel string) error {
	return fmt.Errorf("SetUpdateChannel implementation pending - will be implemented in task 8")
}

// CheckCompatibility checks compatibility
func (m *Manager) CheckCompatibility(projectPath string) (*interfaces.CompatibilityResult, error) {
	return nil, fmt.Errorf("CheckCompatibility implementation pending - will be implemented in task 8")
}

// GetSupportedVersions gets supported versions
func (m *Manager) GetSupportedVersions() (map[string][]string, error) {
	return nil, fmt.Errorf("GetSupportedVersions implementation pending - will be implemented in task 8")
}

// ValidateVersionRequirements validates version requirements
func (m *Manager) ValidateVersionRequirements(requirements map[string]string) (*interfaces.VersionValidationResult, error) {
	return nil, fmt.Errorf("ValidateVersionRequirements implementation pending - will be implemented in task 8")
}

// CacheVersionInfo caches version info
func (m *Manager) CacheVersionInfo(info *interfaces.VersionInfo) error {
	return fmt.Errorf("CacheVersionInfo implementation pending - will be implemented in task 8")
}

// GetCachedVersionInfo gets cached version info
func (m *Manager) GetCachedVersionInfo() (*interfaces.VersionInfo, error) {
	return nil, fmt.Errorf("GetCachedVersionInfo implementation pending - will be implemented in task 8")
}

// RefreshVersionCache refreshes version cache
func (m *Manager) RefreshVersionCache() error {
	return fmt.Errorf("RefreshVersionCache implementation pending - will be implemented in task 8")
}

// ClearVersionCache clears version cache
func (m *Manager) ClearVersionCache() error {
	return fmt.Errorf("ClearVersionCache implementation pending - will be implemented in task 8")
}

// GetReleaseNotes gets release notes
func (m *Manager) GetReleaseNotes(version string) (*interfaces.ReleaseNotes, error) {
	return nil, fmt.Errorf("GetReleaseNotes implementation pending - will be implemented in task 8")
}

// GetChangeLog gets change log
func (m *Manager) GetChangeLog(fromVersion, toVersion string) (*interfaces.ChangeLog, error) {
	return nil, fmt.Errorf("GetChangeLog implementation pending - will be implemented in task 8")
}

// GetSecurityAdvisories gets security advisories
func (m *Manager) GetSecurityAdvisories(version string) ([]interfaces.SecurityAdvisory, error) {
	return nil, fmt.Errorf("GetSecurityAdvisories implementation pending - will be implemented in task 8")
}

// GetPackageInfo gets package info
func (m *Manager) GetPackageInfo(packageName string) (*interfaces.PackageInfo, error) {
	return nil, fmt.Errorf("GetPackageInfo implementation pending - will be implemented in task 8")
}

// GetPackageVersions gets package versions
func (m *Manager) GetPackageVersions(packageName string) ([]string, error) {
	return nil, fmt.Errorf("GetPackageVersions implementation pending - will be implemented in task 8")
}

// CheckPackageUpdates checks package updates
func (m *Manager) CheckPackageUpdates(packages map[string]string) (map[string]interfaces.PackageUpdate, error) {
	return nil, fmt.Errorf("CheckPackageUpdates implementation pending - will be implemented in task 8")
}

// SetVersionConfig sets version config
func (m *Manager) SetVersionConfig(config *interfaces.VersionConfig) error {
	return fmt.Errorf("SetVersionConfig implementation pending - will be implemented in task 8")
}

// GetVersionConfig gets version config
func (m *Manager) GetVersionConfig() (*interfaces.VersionConfig, error) {
	return nil, fmt.Errorf("GetVersionConfig implementation pending - will be implemented in task 8")
}

// SetAutoUpdate sets auto update
func (m *Manager) SetAutoUpdate(enabled bool) error {
	return fmt.Errorf("SetAutoUpdate implementation pending - will be implemented in task 8")
}

// SetUpdateNotifications sets update notifications
func (m *Manager) SetUpdateNotifications(enabled bool) error {
	return fmt.Errorf("SetUpdateNotifications implementation pending - will be implemented in task 8")
}
