package version

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// Manager implements the VersionManager interface
type Manager struct {
	httpClient   *http.Client
	npmClient    *NPMClient
	goClient     *GoClient
	githubClient *GitHubClient
	cache        interfaces.VersionCache
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
	}
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
