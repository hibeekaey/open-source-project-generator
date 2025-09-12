package version

import (
	"fmt"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// GoRegistry implements VersionRegistry for Go modules
type GoRegistry struct {
	client *GoClient
}

// NewGoRegistry creates a new Go registry client
func NewGoRegistry(client *GoClient) *GoRegistry {
	return &GoRegistry{
		client: client,
	}
}

// GetLatestVersion retrieves the latest version for a Go module
func (r *GoRegistry) GetLatestVersion(moduleName string) (*models.VersionInfo, error) {
	version, err := r.client.GetLatestVersion(moduleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for %s: %w", moduleName, err)
	}

	info := &models.VersionInfo{
		Name:           moduleName,
		Language:       "go",
		Type:           "package",
		LatestVersion:  version,
		IsSecure:       true, // Will be updated by security check
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "go",
		RegistryURL:    fmt.Sprintf("https://pkg.go.dev/%s", moduleName),
		SecurityIssues: make([]models.SecurityIssue, 0),
		Metadata:       make(map[string]string),
	}

	// Check for security issues
	securityIssues, err := r.CheckSecurity(moduleName, version)
	if err == nil {
		info.SecurityIssues = securityIssues
		info.IsSecure = len(securityIssues) == 0
	}

	return info, nil
}

// GetVersionHistory retrieves version history for a Go module
func (r *GoRegistry) GetVersionHistory(moduleName string, limit int) ([]*models.VersionInfo, error) {
	// For now, just return the latest version
	// This could be enhanced to fetch actual version history from Go proxy
	latest, err := r.GetLatestVersion(moduleName)
	if err != nil {
		return nil, err
	}

	return []*models.VersionInfo{latest}, nil
}

// CheckSecurity checks for security vulnerabilities in a specific version
func (r *GoRegistry) CheckSecurity(moduleName, version string) ([]models.SecurityIssue, error) {
	// This would integrate with Go vulnerability database
	// For now, return empty slice (no security issues found)
	// TODO: Implement actual security checking using govulncheck or similar
	return []models.SecurityIssue{}, nil
}

// GetRegistryInfo returns information about the Go registry
func (r *GoRegistry) GetRegistryInfo() interfaces.RegistryInfo {
	return interfaces.RegistryInfo{
		Name:        "Go Module Proxy",
		URL:         "https://proxy.golang.org",
		Type:        "go",
		Description: "Official Go module proxy for Go packages and modules",
		Supported:   []string{"go"},
	}
}

// IsAvailable checks if the Go registry is currently accessible
func (r *GoRegistry) IsAvailable() bool {
	// Simple check by trying to get a well-known module
	_, err := r.client.GetLatestVersion("github.com/gin-gonic/gin")
	return err == nil
}

// GetSupportedPackages returns a list of packages supported by this registry
func (r *GoRegistry) GetSupportedPackages() ([]string, error) {
	// Return common Go modules that we track
	return []string{
		"github.com/gin-gonic/gin",
		"github.com/gorilla/mux",
		"github.com/stretchr/testify",
		"gorm.io/gorm",
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
	}, nil
}

// Ensure GoRegistry implements VersionRegistry interface
var _ interfaces.VersionRegistry = (*GoRegistry)(nil)
