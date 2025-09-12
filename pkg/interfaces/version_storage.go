package interfaces

import (
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// VersionStorage defines the contract for persisting version information
type VersionStorage interface {
	// Load loads the complete version store from storage
	Load() (*models.VersionStore, error)

	// Save saves the complete version store to storage
	Save(store *models.VersionStore) error

	// GetVersionInfo retrieves version information for a specific package
	GetVersionInfo(name string) (*models.VersionInfo, error)

	// SetVersionInfo stores version information for a specific package
	SetVersionInfo(name string, info *models.VersionInfo) error

	// DeleteVersionInfo removes version information for a specific package
	DeleteVersionInfo(name string) error

	// ListVersions returns all stored version information
	ListVersions() (map[string]*models.VersionInfo, error)

	// Query searches for version information based on criteria
	Query(query *models.VersionQuery) (map[string]*models.VersionInfo, error)

	// Backup creates a backup of the current version store
	Backup() error

	// Restore restores the version store from a backup
	Restore(backupPath string) error

	// GetLastUpdated returns the timestamp of the last update
	GetLastUpdated() (time.Time, error)

	// SetLastUpdated updates the last updated timestamp
	SetLastUpdated(timestamp time.Time) error
}

// VersionRegistry defines the contract for querying external package registries
type VersionRegistry interface {
	// GetLatestVersion retrieves the latest version for a package
	GetLatestVersion(packageName string) (*models.VersionInfo, error)

	// GetVersionHistory retrieves version history for a package
	GetVersionHistory(packageName string, limit int) ([]*models.VersionInfo, error)

	// CheckSecurity checks for security vulnerabilities in a specific version
	CheckSecurity(packageName, version string) ([]models.SecurityIssue, error)

	// GetRegistryInfo returns information about the registry
	GetRegistryInfo() RegistryInfo

	// IsAvailable checks if the registry is currently accessible
	IsAvailable() bool

	// GetSupportedPackages returns a list of packages supported by this registry
	GetSupportedPackages() ([]string, error)
}

// RegistryInfo provides metadata about a version registry
type RegistryInfo struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Supported   []string `json:"supported_languages"`
}

// TemplateUpdater defines the contract for updating template files with new versions
type TemplateUpdater interface {
	// UpdateTemplate updates a single template file with new version information
	UpdateTemplate(templatePath string, versions map[string]*models.VersionInfo) error

	// UpdateAllTemplates updates all templates with new version information
	UpdateAllTemplates(versions map[string]*models.VersionInfo) error

	// ValidateTemplate validates that a template can be updated
	ValidateTemplate(templatePath string) error

	// GetAffectedTemplates returns templates that would be affected by version changes
	GetAffectedTemplates(versions map[string]*models.VersionInfo) ([]string, error)

	// BackupTemplates creates backups of templates before updating
	BackupTemplates(templatePaths []string) error

	// RestoreTemplates restores templates from backup
	RestoreTemplates(templatePaths []string) error
}
