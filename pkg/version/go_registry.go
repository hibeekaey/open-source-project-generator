package version

import (
	"fmt"
	"time"

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
		Name:          moduleName,
		Language:      "go",
		Type:          "package",
		LatestVersion: version,
		UpdatedAt:     time.Now(),
		CheckedAt:     time.Now(),
		UpdateSource:  "go",
		RegistryURL:   fmt.Sprintf("https://pkg.go.dev/%s", moduleName),
		Metadata:      make(map[string]string),
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
