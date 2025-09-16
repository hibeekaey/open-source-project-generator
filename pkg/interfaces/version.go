package interfaces

import "github.com/cuesoftinc/open-source-project-generator/pkg/models"

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

	// GetVersionHistory returns version history for a package
	GetVersionHistory(packageName string) ([]string, error)
}
