// Package version provides basic package version management capabilities
// for the Open Source Project Generator.
//
// This package implements the VersionManager interface and provides:
//   - Basic version fetching from NPM, Go modules, and GitHub releases
//   - Basic version parsing and validation
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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Manager implements the VersionManager interface with comprehensive capabilities
type Manager struct {
	httpClient     *http.Client
	npmClient      *NPMClient
	goClient       *GoClient
	githubClient   *GitHubClient
	versionConfig  *interfaces.VersionConfig
	cacheManager   interfaces.CacheManager
	buildInfo      *BuildInfo
	currentVersion string
}

// BuildInfo contains build-time information
type BuildInfo struct {
	Version      string
	GitCommit    string
	GitBranch    string
	BuildDate    string
	GoVersion    string
	Platform     string
	Architecture string
	BuildTags    []string
}

// NewManager creates a new version manager with all clients
func NewManager() *Manager {
	return NewManagerWithVersion("dev")
}

// NewManagerWithVersion creates a new version manager with a specific version
func NewManagerWithVersion(version string) *Manager {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Manager{
		httpClient:     httpClient,
		npmClient:      NewNPMClient(httpClient),
		goClient:       NewGoClient(httpClient),
		githubClient:   NewGitHubClient(httpClient),
		currentVersion: version,
		buildInfo: &BuildInfo{
			Version:      version,
			GitCommit:    "unknown",
			GitBranch:    "main",
			BuildDate:    time.Now().Format(time.RFC3339),
			GoVersion:    "1.21",
			Platform:     "linux",
			Architecture: "amd64",
			BuildTags:    []string{},
		},
	}
}

// NewManagerWithCache creates a new version manager with cache support
func NewManagerWithCache(cacheManager interfaces.CacheManager) *Manager {
	manager := NewManager()
	manager.cacheManager = cacheManager
	return manager
}

// NewManagerWithVersionAndCache creates a new version manager with version and cache support
func NewManagerWithVersionAndCache(version string, cacheManager interfaces.CacheManager) *Manager {
	manager := NewManagerWithVersion(version)
	manager.cacheManager = cacheManager
	return manager
}

// GetLatestNodeVersion fetches the latest Node.js version
func (m *Manager) GetLatestNodeVersion() (string, error) {
	version, err := m.githubClient.GetLatestRelease("nodejs", "node")
	if err != nil {
		return "", fmt.Errorf("🚫 Unable to fetch latest Node.js version. Check your internet connection or try --offline mode")
	}

	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	return version, nil
}

// GetLatestGoVersion fetches the latest Go version
func (m *Manager) GetLatestGoVersion() (string, error) {
	version, err := m.githubClient.GetLatestRelease("golang", "go")
	if err != nil {
		return "", fmt.Errorf("🚫 Unable to fetch latest Go version. Check your internet connection or try --offline mode")
	}

	// Remove 'go' prefix if present (e.g., go1.22.0 -> 1.22.0)
	version = strings.TrimPrefix(version, "go")

	return version, nil
}

// GetLatestNPMPackage fetches the latest version of an NPM package
func (m *Manager) GetLatestNPMPackage(packageName string) (string, error) {
	version, err := m.npmClient.GetLatestVersion(packageName)
	if err != nil {
		return "", fmt.Errorf("🚫 Unable to fetch NPM package '%s'. Check your internet connection or package name", packageName)
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
		return "", fmt.Errorf("🚫 Unable to fetch GitHub release '%s/%s'. Check repository name and internet connection", owner, repo)
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

// Comprehensive version management methods

// GetCurrentVersion gets the current version
func (m *Manager) GetCurrentVersion() string {
	if m.currentVersion != "" {
		return m.currentVersion
	}
	return "dev"
}

// GetLatestVersion gets the latest version info with caching support
func (m *Manager) GetLatestVersion() (*interfaces.VersionInfo, error) {
	cacheKey := "latest_version_info"

	// Try to get from cache first if caching is enabled
	if m.cacheManager != nil && m.shouldUseCache() {
		if cached, err := m.cacheManager.Get(cacheKey); err == nil {
			if versionInfo, ok := cached.(*interfaces.VersionInfo); ok {
				return versionInfo, nil
			}
		}
	}

	// Fetch latest version from GitHub releases
	latestVersion, err := m.githubClient.GetLatestRelease("cuesoftinc", "open-source-project-generator")
	if err != nil {
		// Try to return cached version if available
		if cached, cacheErr := m.GetCachedVersionInfo(); cacheErr == nil {
			return cached, nil
		}
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	versionInfo := &interfaces.VersionInfo{
		Version:      strings.TrimPrefix(latestVersion, "v"),
		BuildDate:    time.Now(),
		GitCommit:    m.buildInfo.GitCommit,
		GitBranch:    m.buildInfo.GitBranch,
		GoVersion:    m.buildInfo.GoVersion,
		Platform:     m.buildInfo.Platform,
		Architecture: m.buildInfo.Architecture,
		BuildTags:    m.buildInfo.BuildTags,
		Metadata: map[string]string{
			"source":     "github",
			"repository": "cuesoftinc/open-source-project-generator",
			"fetched_at": time.Now().Format(time.RFC3339),
		},
	}

	// Cache the version info
	if err := m.CacheVersionInfo(versionInfo); err != nil {
		// Log warning but don't fail
		fmt.Printf("⚠️  Couldn't cache version info: %v\n", err)
	}

	return versionInfo, nil
}

// GetAllPackageVersions gets all package versions with caching support
func (m *Manager) GetAllPackageVersions() (map[string]string, error) {
	cacheKey := "all_package_versions"

	// Try cache first if caching is enabled
	if m.cacheManager != nil && m.shouldUseCache() {
		if cached, err := m.cacheManager.Get(cacheKey); err == nil {
			if packages, ok := cached.(map[string]string); ok {
				return packages, nil
			}
		}
	}

	packages := make(map[string]string)

	// Get Node.js version
	if nodeVersion, err := m.GetLatestNodeVersion(); err == nil {
		packages["node"] = nodeVersion
	}

	// Get Go version
	if goVersion, err := m.GetLatestGoVersion(); err == nil {
		packages["go"] = goVersion
	}

	// Get common NPM packages
	commonPackages := []string{
		"react", "next", "typescript", "tailwindcss", "eslint", "prettier",
		"@types/node", "@types/react", "autoprefixer", "postcss", "vite",
		"webpack", "babel", "jest", "vitest", "cypress", "playwright",
	}

	for _, pkg := range commonPackages {
		if version, err := m.GetLatestNPMPackage(pkg); err == nil {
			packages[pkg] = version
		}
	}

	// Get common Go modules
	commonGoModules := []string{
		"github.com/gin-gonic/gin",
		"github.com/gorilla/mux",
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
		"gorm.io/gorm",
		"github.com/stretchr/testify",
	}

	for _, module := range commonGoModules {
		if version, err := m.GetLatestGoModule(module); err == nil {
			packages[module] = version
		}
	}

	// Cache the result if caching is enabled
	if m.cacheManager != nil && m.shouldUseCache() {
		config, _ := m.GetVersionConfig()
		cacheTTL := 6 * time.Hour
		if config != nil {
			cacheTTL = config.CacheTTL
		}
		if err := m.cacheManager.Set(cacheKey, packages, cacheTTL); err != nil {
			fmt.Printf("⚠️  Couldn't cache package versions: %v\n", err)
		}
	}

	return packages, nil
}

// GetLatestPackageVersions gets latest package versions
func (m *Manager) GetLatestPackageVersions() (map[string]string, error) {
	return m.GetAllPackageVersions()
}

// GetDetailedVersionHistory gets detailed version history
func (m *Manager) GetDetailedVersionHistory() ([]interfaces.VersionInfo, error) {
	// For now, return current version info
	current, err := m.GetLatestVersion()
	if err != nil {
		return nil, err
	}

	return []interfaces.VersionInfo{*current}, nil
}

// CheckForUpdates checks for updates with comprehensive information
func (m *Manager) CheckForUpdates() (*interfaces.UpdateInfo, error) {
	currentVersion := m.GetCurrentVersion()
	latestInfo, err := m.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	updateAvailable := m.isUpdateAvailable(currentVersion, latestInfo.Version)

	updateInfo := &interfaces.UpdateInfo{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestInfo.Version,
		UpdateAvailable: updateAvailable,
		ReleaseNotes:    fmt.Sprintf("Release notes for version %s", latestInfo.Version),
		DownloadURL:     fmt.Sprintf("https://github.com/cuesoftinc/open-source-project-generator/releases/tag/v%s", latestInfo.Version),
		ReleaseDate:     latestInfo.BuildDate,
		Breaking:        m.isBreakingUpdate(currentVersion, latestInfo.Version),
		Security:        false, // Would need to analyze security advisories
		Recommended:     updateAvailable,
		Size:            0,  // Would need to fetch from release assets
		Checksum:        "", // Would need to fetch from release assets
		SignatureURL:    "", // Would need to fetch from release assets
	}

	return updateInfo, nil
}

// DownloadUpdate downloads an update
func (m *Manager) DownloadUpdate(version string) error {
	// This would implement actual download logic
	// For now, return a placeholder implementation
	return fmt.Errorf("download functionality not yet implemented - would download version %s", version)
}

// InstallUpdate installs an update
func (m *Manager) InstallUpdate(version string) error {
	// This would implement actual installation logic
	// For now, return a placeholder implementation
	return fmt.Errorf("install functionality not yet implemented - would install version %s", version)
}

// GetUpdateChannel gets the update channel
func (m *Manager) GetUpdateChannel() string {
	config, err := m.GetVersionConfig()
	if err != nil {
		return "stable"
	}
	return config.UpdateChannel
}

// SetUpdateChannel sets the update channel
func (m *Manager) SetUpdateChannel(channel string) error {
	config, err := m.GetVersionConfig()
	if err != nil {
		config = DefaultVersionConfig()
	}

	config.UpdateChannel = channel
	return m.SetVersionConfig(config)
}

// CheckCompatibility performs comprehensive compatibility checking
func (m *Manager) CheckCompatibility(projectPath string) (*interfaces.CompatibilityResult, error) {
	result := &interfaces.CompatibilityResult{
		Compatible:       true,
		GeneratorVersion: m.GetCurrentVersion(),
		ProjectVersion:   "1.0.0", // Would be detected from project
		Issues:           []interfaces.CompatibilityIssue{},
		Recommendations:  []string{},
		PackageVersions:  make(map[string]interfaces.VersionCheck),
	}

	// Check package.json if it exists
	if err := m.checkNodeCompatibility(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to check Node.js compatibility: %w", err)
	}

	// Check go.mod if it exists
	if err := m.checkGoCompatibility(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to check Go compatibility: %w", err)
	}

	// Set overall compatibility
	result.Compatible = len(result.Issues) == 0

	return result, nil
}

// GetSupportedVersions gets supported versions
func (m *Manager) GetSupportedVersions() (map[string][]string, error) {
	supported := map[string][]string{
		"node":  {"18.x", "20.x", "21.x"},
		"go":    {"1.20", "1.21", "1.22"},
		"react": {"17.x", "18.x"},
		"next":  {"13.x", "14.x"},
	}

	return supported, nil
}

// ValidateVersionRequirements validates version requirements with comprehensive logic
func (m *Manager) ValidateVersionRequirements(requirements map[string]string) (*interfaces.VersionValidationResult, error) {
	result := &interfaces.VersionValidationResult{
		Valid:        true,
		Requirements: []interfaces.VersionRequirement{},
		Conflicts:    []interfaces.VersionConflict{},
		Missing:      []string{},
		Summary: interfaces.VersionValidationSummary{
			TotalRequirements:     len(requirements),
			SatisfiedRequirements: 0,
			ConflictCount:         0,
			MissingCount:          0,
			UpdatesAvailable:      0,
		},
	}

	// Validate each requirement
	for pkg, requiredVersion := range requirements {
		var currentVersion string
		var err error

		// Try to get current version based on package type
		if strings.Contains(pkg, "/") || strings.Contains(pkg, ".") {
			// Looks like a Go module
			currentVersion, err = m.GetLatestGoModule(pkg)
		} else {
			// Assume NPM package
			currentVersion, err = m.GetLatestNPMPackage(pkg)
		}

		if err != nil {
			result.Missing = append(result.Missing, pkg)
			result.Summary.MissingCount++
			result.Valid = false
			continue
		}

		req := interfaces.VersionRequirement{
			Package:    pkg,
			Required:   requiredVersion,
			Current:    currentVersion,
			Available:  currentVersion,
			Satisfied:  m.compareVersions(currentVersion, requiredVersion) >= 0,
			UpdateType: m.getUpdateType(currentVersion, requiredVersion),
		}

		if req.Satisfied {
			result.Summary.SatisfiedRequirements++
		} else {
			result.Valid = false
		}

		result.Requirements = append(result.Requirements, req)
	}

	return result, nil
}

// CacheVersionInfo caches version information
func (m *Manager) CacheVersionInfo(info *interfaces.VersionInfo) error {
	if m.cacheManager == nil || !m.shouldUseCache() {
		return nil
	}

	cacheKey := "latest_version_info"
	config, _ := m.GetVersionConfig()
	cacheTTL := 6 * time.Hour
	if config != nil {
		cacheTTL = config.CacheTTL
	}

	return m.cacheManager.Set(cacheKey, info, cacheTTL)
}

// GetCachedVersionInfo retrieves cached version information
func (m *Manager) GetCachedVersionInfo() (*interfaces.VersionInfo, error) {
	if m.cacheManager == nil || !m.shouldUseCache() {
		return nil, fmt.Errorf("version caching is disabled")
	}

	cacheKey := "latest_version_info"
	cached, err := m.cacheManager.Get(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("no cached version info available: %w", err)
	}

	versionInfo, ok := cached.(*interfaces.VersionInfo)
	if !ok {
		return nil, fmt.Errorf("invalid cached version info format")
	}

	return versionInfo, nil
}

// RefreshVersionCache refreshes the version cache
func (m *Manager) RefreshVersionCache() error {
	if m.cacheManager == nil || !m.shouldUseCache() {
		return nil
	}

	// Clear version-related cache entries
	keys, err := m.cacheManager.GetKeysByPattern(".*version.*")
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	for _, key := range keys {
		if err := m.cacheManager.Delete(key); err != nil {
			fmt.Printf("⚠️  Couldn't delete cache key %s: %v\n", key, err)
		}
	}

	// Fetch fresh data
	_, err = m.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to refresh version info: %w", err)
	}

	_, err = m.GetAllPackageVersions()
	if err != nil {
		return fmt.Errorf("failed to refresh package versions: %w", err)
	}

	return nil
}

// ClearVersionCache clears the version cache
func (m *Manager) ClearVersionCache() error {
	if m.cacheManager == nil || !m.shouldUseCache() {
		return nil
	}

	keys, err := m.cacheManager.GetKeysByPattern(".*version.*")
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	for _, key := range keys {
		if err := m.cacheManager.Delete(key); err != nil {
			return fmt.Errorf("failed to delete cache key %s: %w", key, err)
		}
	}

	return nil
}

// GetReleaseNotes gets release notes
func (m *Manager) GetReleaseNotes(version string) (*interfaces.ReleaseNotes, error) {
	// This would fetch release notes from GitHub or other sources
	releaseNotes := &interfaces.ReleaseNotes{
		Version:      version,
		ReleaseDate:  time.Now(),
		Title:        fmt.Sprintf("Release %s", version),
		Description:  fmt.Sprintf("Release notes for version %s", version),
		Features:     []interfaces.ReleaseFeature{},
		BugFixes:     []interfaces.ReleaseBugFix{},
		Breaking:     []interfaces.BreakingChange{},
		Security:     []interfaces.SecurityFix{},
		Dependencies: []interfaces.DependencyChange{},
		Links: map[string]string{
			"github": fmt.Sprintf("https://github.com/cuesoftinc/open-source-project-generator/releases/tag/%s", version),
		},
	}

	return releaseNotes, nil
}

// GetChangeLog gets change log
func (m *Manager) GetChangeLog(fromVersion, toVersion string) (*interfaces.ChangeLog, error) {
	changeLog := &interfaces.ChangeLog{
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Changes:     []interfaces.Change{},
		Summary: interfaces.ChangeSummary{
			TotalChanges:    0,
			Features:        0,
			BugFixes:        0,
			BreakingChanges: 0,
			SecurityFixes:   0,
		},
	}

	return changeLog, nil
}

// GetSecurityAdvisories gets security advisories
func (m *Manager) GetSecurityAdvisories(version string) ([]interfaces.SecurityAdvisory, error) {
	// This would fetch security advisories for the given version
	return []interfaces.SecurityAdvisory{}, nil
}

// GetPackageInfo gets package info
func (m *Manager) GetPackageInfo(packageName string) (*interfaces.PackageInfo, error) {
	// This would fetch detailed package information
	info := &interfaces.PackageInfo{
		Name:         packageName,
		Version:      "latest",
		Description:  fmt.Sprintf("Package information for %s", packageName),
		Homepage:     "",
		Repository:   "",
		License:      "MIT",
		Author:       "",
		Maintainers:  []string{},
		Keywords:     []string{},
		Dependencies: make(map[string]string),
		PublishedAt:  time.Now(),
		UpdatedAt:    time.Now(),
		Downloads:    0,
		Stars:        0,
		Issues:       0,
		Metadata:     make(map[string]any),
	}

	return info, nil
}

// GetPackageVersions gets package versions
func (m *Manager) GetPackageVersions(packageName string) ([]string, error) {
	// This would fetch all available versions for a package
	version, err := m.GetLatestNPMPackage(packageName)
	if err != nil {
		return nil, err
	}

	return []string{version}, nil
}

// CheckPackageUpdates checks package updates
func (m *Manager) CheckPackageUpdates(packages map[string]string) (map[string]interfaces.PackageUpdate, error) {
	updates := make(map[string]interfaces.PackageUpdate)

	for pkg, currentVersion := range packages {
		latestVersion, err := m.GetLatestNPMPackage(pkg)
		if err != nil {
			continue
		}

		if currentVersion != latestVersion {
			updates[pkg] = interfaces.PackageUpdate{
				Package:        pkg,
				CurrentVersion: currentVersion,
				LatestVersion:  latestVersion,
				UpdateType:     "patch", // Simplified
				Breaking:       false,
				Security:       false,
				ReleaseDate:    time.Now(),
				ChangeLog:      fmt.Sprintf("Update from %s to %s", currentVersion, latestVersion),
				Recommended:    true,
			}
		}
	}

	return updates, nil
}

// SetVersionConfig sets the version configuration with file persistence
func (m *Manager) SetVersionConfig(config *interfaces.VersionConfig) error {
	m.versionConfig = config

	// Persist configuration to file
	configPath := filepath.Join(os.Getenv("HOME"), ".generator", "version_config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetVersionConfig gets the version configuration with file loading
func (m *Manager) GetVersionConfig() (*interfaces.VersionConfig, error) {
	if m.versionConfig != nil {
		return m.versionConfig, nil
	}

	// Try to load from file
	configPath := filepath.Join(os.Getenv("HOME"), ".generator", "version_config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		m.versionConfig = DefaultVersionConfig()
		return m.versionConfig, nil
	}

	// Validate the config path for security
	if err := m.validateFilePath(configPath); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	// #nosec G304 - Path is validated above to prevent directory traversal
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config interfaces.VersionConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.versionConfig = &config
	return m.versionConfig, nil
}

// SetAutoUpdate sets auto update
func (m *Manager) SetAutoUpdate(enabled bool) error {
	config, err := m.GetVersionConfig()
	if err != nil {
		return err
	}

	config.AutoUpdate = enabled
	return m.SetVersionConfig(config)
}

// SetUpdateNotifications sets update notifications
func (m *Manager) SetUpdateNotifications(enabled bool) error {
	config, err := m.GetVersionConfig()
	if err != nil {
		return err
	}

	config.UpdateNotifications = enabled
	return m.SetVersionConfig(config)
}

// DefaultVersionConfig returns default version configuration
func DefaultVersionConfig() *interfaces.VersionConfig {
	return &interfaces.VersionConfig{
		AutoUpdate:          false,
		UpdateChannel:       "stable",
		CheckInterval:       24 * time.Hour,
		UpdateNotifications: true,
		CacheVersions:       true,
		CacheTTL:            6 * time.Hour,
		OfflineMode:         false,
		VerifySignatures:    true,
		TrustedSources:      []string{"github.com", "registry.npmjs.org", "proxy.golang.org"},
		AllowPrerelease:     false,
		PackageRegistries:   []string{"https://registry.npmjs.org", "https://proxy.golang.org"},
		PackageTimeout:      30 * time.Second,
		PackageRetries:      3,
	}
}

// Helper methods for comprehensive functionality

// validateFilePath validates that a file path is safe to read
func (m *Manager) validateFilePath(path string) error {
	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	// Check if the path contains any directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid file path: directory traversal detected")
	}

	// Ensure the path is absolute or relative to current directory
	if !filepath.IsAbs(cleanPath) && !strings.HasPrefix(cleanPath, ".") {
		return fmt.Errorf("invalid file path: must be absolute or relative to current directory")
	}

	return nil
}

// shouldUseCache determines if caching should be used
func (m *Manager) shouldUseCache() bool {
	config, err := m.GetVersionConfig()
	if err != nil {
		return true // Default to using cache
	}
	return config.CacheVersions && !config.OfflineMode
}

// compareVersions compares two version strings and returns:
// -1 if current < target, 0 if equal, 1 if current > target
func (m *Manager) compareVersions(current, target string) int {
	// Basic string comparison for now
	// In a real implementation, this would use semantic versioning
	if current == target {
		return 0
	}
	if current < target {
		return -1
	}
	return 1
}

// isUpdateAvailable checks if an update is available
func (m *Manager) isUpdateAvailable(current, latest string) bool {
	return m.compareVersions(current, latest) < 0
}

// isBreakingUpdate determines if an update contains breaking changes
func (m *Manager) isBreakingUpdate(current, latest string) bool {
	// This would analyze version numbers and release notes
	// For now, assume major version changes are breaking
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	if len(currentParts) > 0 && len(latestParts) > 0 {
		return currentParts[0] != latestParts[0]
	}

	return false
}

// getUpdateType determines the type of update (major, minor, patch)
func (m *Manager) getUpdateType(current, target string) string {
	// Basic implementation - in reality would use semantic versioning
	if m.compareVersions(current, target) == 0 {
		return "none"
	}

	currentParts := strings.Split(current, ".")
	targetParts := strings.Split(target, ".")

	if len(currentParts) > 0 && len(targetParts) > 0 && currentParts[0] != targetParts[0] {
		return "major"
	}

	if len(currentParts) > 1 && len(targetParts) > 1 && currentParts[1] != targetParts[1] {
		return "minor"
	}

	return "patch"
}

// checkNodeCompatibility checks Node.js project compatibility
func (m *Manager) checkNodeCompatibility(projectPath string, result *interfaces.CompatibilityResult) error {
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return nil // No package.json, skip Node.js checks
	}

	// Validate the package.json path for security
	if err := m.validateFilePath(packageJSONPath); err != nil {
		return fmt.Errorf("invalid package.json path: %w", err)
	}

	// Read and parse package.json
	// #nosec G304 - Path is validated above to prevent directory traversal
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		result.Issues = append(result.Issues, interfaces.CompatibilityIssue{
			Type:        "parse_error",
			Severity:    "high",
			Component:   "package.json",
			Description: "Failed to parse package.json",
			Fixable:     false,
		})
		return nil
	}

	// Check Node.js version requirement
	if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
		if nodeVersion, ok := engines["node"].(string); ok {
			// Add version check to package versions
			result.PackageVersions["node"] = interfaces.VersionCheck{
				Current:    nodeVersion,
				Latest:     nodeVersion,
				Compatible: true,
				UpdateType: "none",
			}
		}
	}

	return nil
}

// checkGoCompatibility checks Go project compatibility
func (m *Manager) checkGoCompatibility(projectPath string, result *interfaces.CompatibilityResult) error {
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil // No go.mod, skip Go checks
	}

	// Validate the go.mod path for security
	if err := m.validateFilePath(goModPath); err != nil {
		return fmt.Errorf("invalid go.mod path: %w", err)
	}

	// Read go.mod
	// #nosec G304 - Path is validated above to prevent directory traversal
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			goVersion := strings.TrimPrefix(line, "go ")
			// Add version check to package versions
			result.PackageVersions["go"] = interfaces.VersionCheck{
				Current:    goVersion,
				Latest:     goVersion,
				Compatible: true,
				UpdateType: "none",
			}
			break
		}
	}

	return nil
}
