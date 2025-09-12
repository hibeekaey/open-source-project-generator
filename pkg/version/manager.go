package version

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// Manager implements the VersionManager interface with enhanced functionality
type Manager struct {
	httpClient   *http.Client
	npmClient    *NPMClient
	goClient     *GoClient
	githubClient *GitHubClient
	cache        interfaces.VersionCache
	storage      interfaces.VersionStorage
	registries   map[string]interfaces.VersionRegistry
}

// NewManager creates a new version manager with all clients
func NewManager(cache interfaces.VersionCache) *Manager {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Manager{
		httpClient:   httpClient,
		npmClient:    NewNPMClient(httpClient),
		goClient:     NewGoClient(httpClient),
		githubClient: NewGitHubClient(httpClient),
		cache:        cache,
		registries:   make(map[string]interfaces.VersionRegistry),
	}
}

// NewManagerWithStorage creates a new version manager with storage and registry integration
func NewManagerWithStorage(cache interfaces.VersionCache, storage interfaces.VersionStorage) *Manager {
	manager := NewManager(cache)
	manager.storage = storage

	// Initialize registries
	manager.registries["npm"] = NewNPMRegistry(manager.npmClient)
	manager.registries["go"] = NewGoRegistry(manager.goClient)
	manager.registries["github"] = NewGitHubRegistry(manager.githubClient)

	return manager
}

// GetLatestNodeVersion fetches the latest Node.js version
func (m *Manager) GetLatestNodeVersion() (string, error) {
	cacheKey := "nodejs:latest"
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}

	version, err := m.githubClient.GetLatestRelease("nodejs", "node")
	if err != nil {
		return "", fmt.Errorf("failed to fetch Node.js version: %w", err)
	}

	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	if err := m.cache.Set(cacheKey, version); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache Node.js version: %v\n", err)
	}

	return version, nil
}

// GetLatestGoVersion fetches the latest Go version
func (m *Manager) GetLatestGoVersion() (string, error) {
	cacheKey := "golang:latest"
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}

	version, err := m.githubClient.GetLatestRelease("golang", "go")
	if err != nil {
		return "", fmt.Errorf("failed to fetch Go version: %w", err)
	}

	// Remove 'go' prefix if present (e.g., go1.22.0 -> 1.22.0)
	version = strings.TrimPrefix(version, "go")

	if err := m.cache.Set(cacheKey, version); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache Go version: %v\n", err)
	}

	return version, nil
}

// GetLatestNPMPackage fetches the latest version of an NPM package
func (m *Manager) GetLatestNPMPackage(packageName string) (string, error) {
	cacheKey := fmt.Sprintf("npm:%s", packageName)
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}

	version, err := m.npmClient.GetLatestVersion(packageName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch NPM package %s: %w", packageName, err)
	}

	if err := m.cache.Set(cacheKey, version); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache NPM package %s version: %v\n", packageName, err)
	}

	return version, nil
}

// GetLatestGoModule fetches the latest version of a Go module
func (m *Manager) GetLatestGoModule(moduleName string) (string, error) {
	cacheKey := fmt.Sprintf("gomod:%s", moduleName)
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}

	version, err := m.goClient.GetLatestVersion(moduleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Go module %s: %w", moduleName, err)
	}

	if err := m.cache.Set(cacheKey, version); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache Go module %s version: %v\n", moduleName, err)
	}

	return version, nil
}

// GetLatestGitHubRelease fetches the latest release version from GitHub
func (m *Manager) GetLatestGitHubRelease(owner, repo string) (string, error) {
	cacheKey := fmt.Sprintf("github:%s/%s", owner, repo)
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}

	version, err := m.githubClient.GetLatestRelease(owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitHub release %s/%s: %w", owner, repo, err)
	}

	if err := m.cache.Set(cacheKey, version); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache GitHub release %s/%s version: %v\n", owner, repo, err)
	}

	return version, nil
}

// UpdateVersionsConfig updates the version configuration with latest versions
func (m *Manager) UpdateVersionsConfig() (*models.VersionConfig, error) {
	config := &models.VersionConfig{
		Packages:  make(map[string]string),
		UpdatedAt: time.Now(),
	}

	// Fetch core language versions
	nodeVersion, err := m.GetLatestNodeVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get Node.js version: %w", err)
	}
	config.Node = nodeVersion

	goVersion, err := m.GetLatestGoVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go version: %w", err)
	}
	config.Go = goVersion

	// Fetch framework versions
	nextjsVersion, err := m.GetLatestNPMPackage("next")
	if err != nil {
		return nil, fmt.Errorf("failed to get Next.js version: %w", err)
	}
	config.NextJS = nextjsVersion

	reactVersion, err := m.GetLatestNPMPackage("react")
	if err != nil {
		return nil, fmt.Errorf("failed to get React version: %w", err)
	}
	config.React = reactVersion

	// Fetch mobile language versions
	kotlinVersion, err := m.GetLatestGitHubRelease("JetBrains", "kotlin")
	if err != nil {
		// Don't fail for optional versions
		fmt.Printf("Warning: failed to get Kotlin version: %v\n", err)
	} else {
		config.Kotlin = strings.TrimPrefix(kotlinVersion, "v")
	}

	swiftVersion, err := m.GetLatestGitHubRelease("apple", "swift")
	if err != nil {
		// Don't fail for optional versions
		fmt.Printf("Warning: failed to get Swift version: %v\n", err)
	} else {
		config.Swift = strings.TrimPrefix(swiftVersion, "swift-")
	}

	// Fetch common package versions
	commonPackages := map[string]string{
		"typescript":   "typescript",
		"tailwindcss":  "tailwindcss",
		"eslint":       "eslint",
		"prettier":     "prettier",
		"jest":         "jest",
		"@types/node":  "@types/node",
		"@types/react": "@types/react",
		"autoprefixer": "autoprefixer",
		"postcss":      "postcss",
	}

	for key, packageName := range commonPackages {
		version, err := m.GetLatestNPMPackage(packageName)
		if err != nil {
			fmt.Printf("Warning: failed to get %s version: %v\n", packageName, err)
			continue
		}
		config.Packages[key] = version
	}

	return config, nil
}

// CacheVersion caches a version for future use
func (m *Manager) CacheVersion(key, version string) error {
	return m.cache.Set(key, version)
}

// GetCachedVersion retrieves a cached version
func (m *Manager) GetCachedVersion(key string) (string, bool) {
	return m.cache.Get(key)
}

// CheckLatestVersions checks for latest versions across all registries
func (m *Manager) CheckLatestVersions() (*models.VersionReport, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	store, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load version store: %w", err)
	}

	report := &models.VersionReport{
		GeneratedAt:     time.Now(),
		Summary:         make(map[string]models.VersionSummary),
		Details:         make(map[string]*models.VersionInfo),
		Recommendations: make([]models.UpdateRecommendation, 0),
	}

	// Check languages
	languageSummary := models.VersionSummary{}
	for name, info := range store.Languages {
		languageSummary.Total++
		report.Details[name] = info

		latestInfo, err := m.getLatestVersionInfo(name, info)
		if err != nil {
			fmt.Printf("Warning: failed to check latest version for %s: %v\n", name, err)
			continue
		}

		if m.isVersionOutdated(info.CurrentVersion, latestInfo.LatestVersion) {
			languageSummary.Outdated++
			report.Recommendations = append(report.Recommendations, models.UpdateRecommendation{
				Name:               name,
				CurrentVersion:     info.CurrentVersion,
				RecommendedVersion: latestInfo.LatestVersion,
				Priority:           m.determinePriority(info, latestInfo),
				Reason:             m.generateUpdateReason(info, latestInfo),
				BreakingChange:     m.isBreakingChange(info.CurrentVersion, latestInfo.LatestVersion),
			})
		} else {
			languageSummary.Current++
		}

		if !latestInfo.IsSecure {
			languageSummary.Insecure++
			report.SecurityIssues++
		}
	}
	report.Summary["languages"] = languageSummary

	// Check frameworks
	frameworkSummary := models.VersionSummary{}
	for name, info := range store.Frameworks {
		frameworkSummary.Total++
		report.Details[name] = info

		latestInfo, err := m.getLatestVersionInfo(name, info)
		if err != nil {
			fmt.Printf("Warning: failed to check latest version for %s: %v\n", name, err)
			continue
		}

		if m.isVersionOutdated(info.CurrentVersion, latestInfo.LatestVersion) {
			frameworkSummary.Outdated++
			report.Recommendations = append(report.Recommendations, models.UpdateRecommendation{
				Name:               name,
				CurrentVersion:     info.CurrentVersion,
				RecommendedVersion: latestInfo.LatestVersion,
				Priority:           m.determinePriority(info, latestInfo),
				Reason:             m.generateUpdateReason(info, latestInfo),
				BreakingChange:     m.isBreakingChange(info.CurrentVersion, latestInfo.LatestVersion),
			})
		} else {
			frameworkSummary.Current++
		}

		if !latestInfo.IsSecure {
			frameworkSummary.Insecure++
			report.SecurityIssues++
		}
	}
	report.Summary["frameworks"] = frameworkSummary

	// Check packages
	packageSummary := models.VersionSummary{}
	for name, info := range store.Packages {
		packageSummary.Total++
		report.Details[name] = info

		latestInfo, err := m.getLatestVersionInfo(name, info)
		if err != nil {
			fmt.Printf("Warning: failed to check latest version for %s: %v\n", name, err)
			continue
		}

		if m.isVersionOutdated(info.CurrentVersion, latestInfo.LatestVersion) {
			packageSummary.Outdated++
			report.Recommendations = append(report.Recommendations, models.UpdateRecommendation{
				Name:               name,
				CurrentVersion:     info.CurrentVersion,
				RecommendedVersion: latestInfo.LatestVersion,
				Priority:           m.determinePriority(info, latestInfo),
				Reason:             m.generateUpdateReason(info, latestInfo),
				BreakingChange:     m.isBreakingChange(info.CurrentVersion, latestInfo.LatestVersion),
			})
		} else {
			packageSummary.Current++
		}

		if !latestInfo.IsSecure {
			packageSummary.Insecure++
			report.SecurityIssues++
		}
	}
	report.Summary["packages"] = packageSummary

	report.TotalPackages = languageSummary.Total + frameworkSummary.Total + packageSummary.Total
	report.OutdatedCount = languageSummary.Outdated + frameworkSummary.Outdated + packageSummary.Outdated
	report.LastUpdateCheck = time.Now()

	return report, nil
}

// UpdateVersionInfo updates version information for a specific package
func (m *Manager) UpdateVersionInfo(name string, targetVersion string, force bool) (*models.VersionUpdateResult, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	result := &models.VersionUpdateResult{
		UpdatedAt: time.Now(),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
		Metadata:  make(map[string]string),
	}

	// Get current version info
	currentInfo, err := m.storage.GetVersionInfo(name)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to get current version info: %v", err))
		return result, err
	}

	result.PreviousVersion = currentInfo.CurrentVersion

	// Validate target version if not forcing
	if !force {
		if m.isVersionOutdated(targetVersion, currentInfo.CurrentVersion) {
			result.Errors = append(result.Errors, "target version is older than current version")
			return result, fmt.Errorf("target version %s is older than current version %s", targetVersion, currentInfo.CurrentVersion)
		}

		// Check for breaking changes
		if m.isBreakingChange(currentInfo.CurrentVersion, targetVersion) {
			result.Warnings = append(result.Warnings, "target version may contain breaking changes")
		}
	}

	// Update version info
	updatedInfo := *currentInfo
	updatedInfo.PreviousVersion = currentInfo.CurrentVersion
	updatedInfo.CurrentVersion = targetVersion
	updatedInfo.UpdatedAt = time.Now()
	updatedInfo.UpdateSource = "manual"

	// Check security for new version
	if registry, exists := m.registries[m.getRegistryForPackage(name, &updatedInfo)]; exists {
		securityIssues, err := registry.CheckSecurity(name, targetVersion)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to check security: %v", err))
		} else {
			updatedInfo.SecurityIssues = securityIssues
			updatedInfo.IsSecure = len(securityIssues) == 0
		}
	}

	// Save updated info
	if err := m.storage.SetVersionInfo(name, &updatedInfo); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to save version info: %v", err))
		return result, err
	}

	result.Success = true
	result.NewVersion = targetVersion
	result.Metadata["update_source"] = "manual"

	return result, nil
}

// CompareVersions compares two version strings and returns comparison result
func (m *Manager) CompareVersions(version1, version2 string) (int, error) {
	// Use semver comparison logic
	semver1, err := ParseSemVer(version1)
	if err != nil {
		return 0, fmt.Errorf("invalid version format for %s: %w", version1, err)
	}

	semver2, err := ParseSemVer(version2)
	if err != nil {
		return 0, fmt.Errorf("invalid version format for %s: %w", version2, err)
	}

	return semver1.Compare(semver2), nil
}

// DetectVersionUpdates detects available updates for all packages
func (m *Manager) DetectVersionUpdates() (map[string]*models.VersionInfo, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	allVersions, err := m.storage.ListVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	updates := make(map[string]*models.VersionInfo)

	for name, info := range allVersions {
		latestInfo, err := m.getLatestVersionInfo(name, info)
		if err != nil {
			fmt.Printf("Warning: failed to check latest version for %s: %v\n", name, err)
			continue
		}

		if m.isVersionOutdated(info.CurrentVersion, latestInfo.LatestVersion) {
			updates[name] = latestInfo
		}
	}

	return updates, nil
}

// Helper methods

func (m *Manager) getLatestVersionInfo(name string, currentInfo *models.VersionInfo) (*models.VersionInfo, error) {
	registryName := m.getRegistryForPackage(name, currentInfo)
	registry, exists := m.registries[registryName]
	if !exists {
		return nil, fmt.Errorf("no registry found for package %s", name)
	}

	return registry.GetLatestVersion(name)
}

func (m *Manager) getRegistryForPackage(name string, info *models.VersionInfo) string {
	// Determine registry based on package name and type
	switch info.Language {
	case "javascript", "typescript", "nodejs":
		return "npm"
	case "go":
		return "go"
	default:
		// Check if it's a GitHub release
		if strings.Contains(name, "/") {
			return "github"
		}
		return "npm" // Default to npm
	}
}

func (m *Manager) isVersionOutdated(current, latest string) bool {
	comparison, err := m.CompareVersions(current, latest)
	if err != nil {
		return false
	}
	return comparison < 0
}

func (m *Manager) isBreakingChange(current, target string) bool {
	currentSemver, err := ParseSemVer(current)
	if err != nil {
		return false
	}

	targetSemver, err := ParseSemVer(target)
	if err != nil {
		return false
	}

	// Major version change is considered breaking
	return targetSemver.Major > currentSemver.Major
}

func (m *Manager) determinePriority(current, latest *models.VersionInfo) string {
	// Security issues get highest priority
	if len(latest.SecurityIssues) > 0 {
		for _, issue := range latest.SecurityIssues {
			if issue.Severity == "critical" {
				return "critical"
			}
		}
		return "high"
	}

	// Breaking changes get medium priority
	if m.isBreakingChange(current.CurrentVersion, latest.LatestVersion) {
		return "medium"
	}

	return "low"
}

func (m *Manager) generateUpdateReason(current, latest *models.VersionInfo) string {
	if len(latest.SecurityIssues) > 0 {
		return fmt.Sprintf("Security vulnerabilities found: %d issues", len(latest.SecurityIssues))
	}

	if m.isBreakingChange(current.CurrentVersion, latest.LatestVersion) {
		return "Major version update available (may contain breaking changes)"
	}

	return "Newer version available"
}

// GetVersionStore returns the current version store
func (m *Manager) GetVersionStore() (*models.VersionStore, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	return m.storage.Load()
}
