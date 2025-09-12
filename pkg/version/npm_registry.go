package version

import (
	"fmt"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// NPMRegistry implements VersionRegistry for NPM packages
type NPMRegistry struct {
	client *NPMClient
}

// NewNPMRegistry creates a new NPM registry client
func NewNPMRegistry(client *NPMClient) *NPMRegistry {
	return &NPMRegistry{
		client: client,
	}
}

// GetLatestVersion retrieves the latest version for an NPM package
func (r *NPMRegistry) GetLatestVersion(packageName string) (*models.VersionInfo, error) {
	version, err := r.client.GetLatestVersion(packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for %s: %w", packageName, err)
	}

	info := &models.VersionInfo{
		Name:           packageName,
		Language:       "javascript",
		Type:           "package",
		LatestVersion:  version,
		IsSecure:       true, // Will be updated by security check
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "npm",
		RegistryURL:    fmt.Sprintf("https://registry.npmjs.org/%s", packageName),
		SecurityIssues: make([]models.SecurityIssue, 0),
		Metadata:       make(map[string]string),
	}

	// Check for security issues
	securityIssues, err := r.CheckSecurity(packageName, version)
	if err == nil {
		info.SecurityIssues = securityIssues
		info.IsSecure = len(securityIssues) == 0
	}

	return info, nil
}

// GetVersionHistory retrieves version history for an NPM package
func (r *NPMRegistry) GetVersionHistory(packageName string, limit int) ([]*models.VersionInfo, error) {
	// For now, just return the latest version
	// This could be enhanced to fetch actual version history from NPM API
	latest, err := r.GetLatestVersion(packageName)
	if err != nil {
		return nil, err
	}

	return []*models.VersionInfo{latest}, nil
}

// CheckSecurity checks for security vulnerabilities in a specific version
func (r *NPMRegistry) CheckSecurity(packageName, version string) ([]models.SecurityIssue, error) {
	// This would integrate with npm audit or security advisory databases
	// For now, return empty slice (no security issues found)
	// TODO: Implement actual security checking
	return []models.SecurityIssue{}, nil
}

// GetRegistryInfo returns information about the NPM registry
func (r *NPMRegistry) GetRegistryInfo() interfaces.RegistryInfo {
	return interfaces.RegistryInfo{
		Name:        "NPM Registry",
		URL:         "https://registry.npmjs.org",
		Type:        "npm",
		Description: "Official NPM package registry for JavaScript and TypeScript packages",
		Supported:   []string{"javascript", "typescript", "nodejs"},
	}
}

// IsAvailable checks if the NPM registry is currently accessible
func (r *NPMRegistry) IsAvailable() bool {
	// Simple check by trying to get a well-known package
	_, err := r.client.GetLatestVersion("react")
	return err == nil
}

// GetSupportedPackages returns a list of packages supported by this registry
func (r *NPMRegistry) GetSupportedPackages() ([]string, error) {
	// Return common packages that we track
	return []string{
		"react",
		"next",
		"typescript",
		"tailwindcss",
		"eslint",
		"prettier",
		"jest",
		"@types/node",
		"@types/react",
		"autoprefixer",
		"postcss",
	}, nil
}

// Ensure NPMRegistry implements VersionRegistry interface
var _ interfaces.VersionRegistry = (*NPMRegistry)(nil)
