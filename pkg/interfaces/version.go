package interfaces

import "github.com/open-source-template-generator/pkg/models"

// VersionCache defines the contract for version caching operations
type VersionCache interface {
	// Get retrieves a cached version by key
	Get(key string) (string, bool)

	// Set stores a version in the cache
	Set(key, version string) error

	// Delete removes a version from the cache
	Delete(key string) error

	// Clear removes all cached versions
	Clear() error

	// Keys returns all cached keys
	Keys() []string
}

// VersionManager defines the contract for package version management operations
type VersionManager interface {
	// GetLatestNodeVersion fetches the latest Node.js version
	GetLatestNodeVersion() (string, error)

	// GetLatestGoVersion fetches the latest Go version
	GetLatestGoVersion() (string, error)

	// GetLatestNPMPackage fetches the latest version of an NPM package
	GetLatestNPMPackage(packageName string) (string, error)

	// GetLatestGoModule fetches the latest version of a Go module
	GetLatestGoModule(moduleName string) (string, error)

	// UpdateVersionsConfig updates the version configuration with latest versions
	UpdateVersionsConfig() (*models.VersionConfig, error)

	// GetLatestGitHubRelease fetches the latest release version from GitHub
	GetLatestGitHubRelease(owner, repo string) (string, error)

	// CacheVersion caches a version for future use
	CacheVersion(key, version string) error

	// GetCachedVersion retrieves a cached version
	GetCachedVersion(key string) (string, bool)
}
