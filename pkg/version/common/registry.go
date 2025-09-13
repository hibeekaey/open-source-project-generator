package common

import (
	"fmt"
	"sort"
)

// RegistryClient defines the interface for registry clients
type RegistryClient interface {
	GetVersions(packageName string) ([]string, error)
	GetPackageInfo(packageName string) (PackageInfo, error)
}

// PackageInfo represents information about a package
type PackageInfo struct {
	Name        string
	Description string
	Homepage    string
	Repository  string
}

// GetLatestVersion retrieves the latest version from a registry
func GetLatestVersion(client RegistryClient, packageName string) (string, error) {
	versions, err := client.GetVersions(packageName)
	if err != nil {
		return "", fmt.Errorf("failed to get versions for %s: %w", packageName, err)
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found for %s", packageName)
	}

	// Sort versions to get the latest (assuming semantic versioning)
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j]) > 0
	})

	return versions[0], nil
}

// GetVersionHistory retrieves version history for a package
func GetVersionHistory(client RegistryClient, packageName string, limit int) ([]string, error) {
	versions, err := client.GetVersions(packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get version history for %s: %w", packageName, err)
	}

	// Sort versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j]) > 0
	})

	// Limit results if requested
	if limit > 0 && len(versions) > limit {
		versions = versions[:limit]
	}

	return versions, nil
}

// compareVersions compares two version strings
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	// This is a simplified version comparison
	// In a real implementation, you would use a proper semver library
	if v1 == v2 {
		return 0
	}
	if v1 > v2 {
		return 1
	}
	return -1
}
