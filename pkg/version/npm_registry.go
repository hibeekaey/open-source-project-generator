package version

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
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
		Name:          packageName,
		Language:      constants.LanguageJavaScript,
		Type:          constants.FileTypePackage,
		LatestVersion: version,
		UpdatedAt:     time.Now(),
		CheckedAt:     time.Now(),
		UpdateSource:  constants.PackageManagerNPM,
		RegistryURL:   fmt.Sprintf("https://registry.npmjs.org/%s", packageName),
		Metadata:      make(map[string]string),
	}

	return info, nil
}

// GetVersionHistory retrieves version history for an NPM package
func (r *NPMRegistry) GetVersionHistory(packageName string, limit int) ([]*models.VersionInfo, error) {
	// For now, just return the latest version
	// This could be improved to fetch actual version history from NPM API
	latest, err := r.GetLatestVersion(packageName)
	if err != nil {
		return nil, err
	}

	return []*models.VersionInfo{latest}, nil
}
