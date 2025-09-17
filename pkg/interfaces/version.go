package interfaces

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// VersionManager defines the contract for comprehensive version management operations
type VersionManager interface {
	// Basic version operations
	GetLatestNodeVersion() (string, error)
	GetLatestGoVersion() (string, error)
	GetLatestNPMPackage(packageName string) (string, error)
	GetLatestGoModule(moduleName string) (string, error)
	UpdateVersionsConfig() (*models.VersionConfig, error)
	GetLatestGitHubRelease(owner, repo string) (string, error)
	GetVersionHistory(packageName string) ([]string, error)

	// Enhanced version information
	GetCurrentVersion() string
	GetLatestVersion() (*VersionInfo, error)
	GetAllPackageVersions() (map[string]string, error)
	GetLatestPackageVersions() (map[string]string, error)
	GetDetailedVersionHistory() ([]VersionInfo, error)

	// Update management
	CheckForUpdates() (*UpdateInfo, error)
	DownloadUpdate(version string) error
	InstallUpdate(version string) error
	GetUpdateChannel() string
	SetUpdateChannel(channel string) error

	// Version compatibility
	CheckCompatibility(projectPath string) (*CompatibilityResult, error)
	GetSupportedVersions() (map[string][]string, error)
	ValidateVersionRequirements(requirements map[string]string) (*VersionValidationResult, error)

	// Version caching
	CacheVersionInfo(info *VersionInfo) error
	GetCachedVersionInfo() (*VersionInfo, error)
	RefreshVersionCache() error
	ClearVersionCache() error

	// Release information
	GetReleaseNotes(version string) (*ReleaseNotes, error)
	GetChangeLog(fromVersion, toVersion string) (*ChangeLog, error)
	GetSecurityAdvisories(version string) ([]SecurityAdvisory, error)

	// Package management
	GetPackageInfo(packageName string) (*PackageInfo, error)
	GetPackageVersions(packageName string) ([]string, error)
	CheckPackageUpdates(packages map[string]string) (map[string]PackageUpdate, error)

	// Version configuration
	SetVersionConfig(config *VersionConfig) error
	GetVersionConfig() (*VersionConfig, error)
	SetAutoUpdate(enabled bool) error
	SetUpdateNotifications(enabled bool) error
}

// Enhanced version types and structures

// VersionInfo contains detailed version information
type VersionInfo struct {
	Version      string            `json:"version"`
	BuildDate    time.Time         `json:"build_date"`
	GitCommit    string            `json:"git_commit"`
	GitBranch    string            `json:"git_branch"`
	GoVersion    string            `json:"go_version"`
	Platform     string            `json:"platform"`
	Architecture string            `json:"architecture"`
	BuildTags    []string          `json:"build_tags"`
	Metadata     map[string]string `json:"metadata"`
}

// VersionConfig defines configuration for version management
type VersionConfig struct {
	// Update settings
	AutoUpdate          bool          `json:"auto_update"`
	UpdateChannel       string        `json:"update_channel"` // stable, beta, alpha
	CheckInterval       time.Duration `json:"check_interval"`
	UpdateNotifications bool          `json:"update_notifications"`

	// Cache settings
	CacheVersions bool          `json:"cache_versions"`
	CacheTTL      time.Duration `json:"cache_ttl"`
	OfflineMode   bool          `json:"offline_mode"`

	// Security settings
	VerifySignatures bool     `json:"verify_signatures"`
	TrustedSources   []string `json:"trusted_sources"`
	AllowPrerelease  bool     `json:"allow_prerelease"`

	// Package settings
	PackageRegistries []string      `json:"package_registries"`
	PackageTimeout    time.Duration `json:"package_timeout"`
	PackageRetries    int           `json:"package_retries"`
}

// VersionValidationResult contains version validation results
type VersionValidationResult struct {
	Valid        bool                     `json:"valid"`
	Requirements []VersionRequirement     `json:"requirements"`
	Conflicts    []VersionConflict        `json:"conflicts"`
	Missing      []string                 `json:"missing"`
	Summary      VersionValidationSummary `json:"summary"`
}

// VersionRequirement represents a version requirement check
type VersionRequirement struct {
	Package    string `json:"package"`
	Required   string `json:"required"`
	Current    string `json:"current"`
	Available  string `json:"available"`
	Satisfied  bool   `json:"satisfied"`
	UpdateType string `json:"update_type"` // major, minor, patch
}

// VersionConflict represents a version conflict
type VersionConflict struct {
	Package1   string `json:"package1"`
	Version1   string `json:"version1"`
	Package2   string `json:"package2"`
	Version2   string `json:"version2"`
	Reason     string `json:"reason"`
	Severity   string `json:"severity"`
	Resolution string `json:"resolution"`
}

// VersionValidationSummary contains version validation statistics
type VersionValidationSummary struct {
	TotalRequirements     int `json:"total_requirements"`
	SatisfiedRequirements int `json:"satisfied_requirements"`
	ConflictCount         int `json:"conflict_count"`
	MissingCount          int `json:"missing_count"`
	UpdatesAvailable      int `json:"updates_available"`
}

// ReleaseNotes contains release notes for a version
type ReleaseNotes struct {
	Version      string             `json:"version"`
	ReleaseDate  time.Time          `json:"release_date"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Features     []ReleaseFeature   `json:"features"`
	BugFixes     []ReleaseBugFix    `json:"bug_fixes"`
	Breaking     []BreakingChange   `json:"breaking"`
	Security     []SecurityFix      `json:"security"`
	Dependencies []DependencyChange `json:"dependencies"`
	Migration    *MigrationGuide    `json:"migration,omitempty"`
	Links        map[string]string  `json:"links"`
}

// ReleaseFeature represents a new feature in a release
type ReleaseFeature struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Impact      string   `json:"impact"` // major, minor, patch
	Components  []string `json:"components"`
	PRNumber    int      `json:"pr_number,omitempty"`
}

// ReleaseBugFix represents a bug fix in a release
type ReleaseBugFix struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Components  []string `json:"components"`
	IssueNumber int      `json:"issue_number,omitempty"`
	PRNumber    int      `json:"pr_number,omitempty"`
}

// BreakingChange represents a breaking change in a release
type BreakingChange struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Component   string         `json:"component"`
	Migration   *MigrationStep `json:"migration,omitempty"`
	Workaround  string         `json:"workaround,omitempty"`
}

// SecurityFix represents a security fix in a release
type SecurityFix struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	CVEID       string   `json:"cve_id,omitempty"`
	CVSS        float64  `json:"cvss,omitempty"`
	Components  []string `json:"components"`
}

// DependencyChange represents a dependency change in a release
type DependencyChange struct {
	Name       string `json:"name"`
	Type       string `json:"type"` // added, updated, removed
	OldVersion string `json:"old_version,omitempty"`
	NewVersion string `json:"new_version,omitempty"`
	Reason     string `json:"reason"`
	Breaking   bool   `json:"breaking"`
}

// MigrationGuide contains migration information for a release
type MigrationGuide struct {
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Steps         []MigrationStep `json:"steps"`
	Automated     bool            `json:"automated"`
	EstimatedTime string          `json:"estimated_time"`
}

// MigrationStep represents a step in a migration guide
type MigrationStep struct {
	Order       int    `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Manual      bool   `json:"manual"`
	Required    bool   `json:"required"`
}

// ChangeLog contains changes between versions
type ChangeLog struct {
	FromVersion string        `json:"from_version"`
	ToVersion   string        `json:"to_version"`
	Changes     []Change      `json:"changes"`
	Summary     ChangeSummary `json:"summary"`
}

// Change represents a single change in the changelog
type Change struct {
	Type        string    `json:"type"` // feature, bugfix, breaking, security
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Component   string    `json:"component"`
	Impact      string    `json:"impact"`
	Date        time.Time `json:"date"`
	Author      string    `json:"author"`
	PRNumber    int       `json:"pr_number,omitempty"`
	IssueNumber int       `json:"issue_number,omitempty"`
}

// ChangeSummary contains summary statistics for changes
type ChangeSummary struct {
	TotalChanges    int `json:"total_changes"`
	Features        int `json:"features"`
	BugFixes        int `json:"bug_fixes"`
	BreakingChanges int `json:"breaking_changes"`
	SecurityFixes   int `json:"security_fixes"`
}

// SecurityAdvisory contains security advisory information
type SecurityAdvisory struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Severity         string    `json:"severity"`
	CVEID            string    `json:"cve_id"`
	CVSS             float64   `json:"cvss"`
	PublishedAt      time.Time `json:"published_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	AffectedVersions []string  `json:"affected_versions"`
	FixedInVersion   string    `json:"fixed_in_version"`
	Workaround       string    `json:"workaround,omitempty"`
	References       []string  `json:"references"`
}

// PackageInfo contains information about a package
type PackageInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Homepage     string            `json:"homepage"`
	Repository   string            `json:"repository"`
	License      string            `json:"license"`
	Author       string            `json:"author"`
	Maintainers  []string          `json:"maintainers"`
	Keywords     []string          `json:"keywords"`
	Dependencies map[string]string `json:"dependencies"`
	PublishedAt  time.Time         `json:"published_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Downloads    int64             `json:"downloads"`
	Stars        int               `json:"stars"`
	Issues       int               `json:"issues"`
	Metadata     map[string]any    `json:"metadata"`
}

// PackageUpdate contains information about a package update
type PackageUpdate struct {
	Package        string    `json:"package"`
	CurrentVersion string    `json:"current_version"`
	LatestVersion  string    `json:"latest_version"`
	UpdateType     string    `json:"update_type"` // major, minor, patch
	Breaking       bool      `json:"breaking"`
	Security       bool      `json:"security"`
	ReleaseDate    time.Time `json:"release_date"`
	ChangeLog      string    `json:"changelog"`
	Recommended    bool      `json:"recommended"`
}

// UpdateChannel defines update channels
const (
	UpdateChannelStable  = "stable"
	UpdateChannelBeta    = "beta"
	UpdateChannelAlpha   = "alpha"
	UpdateChannelNightly = "nightly"
)

// UpdateType defines types of updates
const (
	UpdateTypeMajor = "major"
	UpdateTypeMinor = "minor"
	UpdateTypePatch = "patch"
)

// VersionSeverity defines severity levels for version issues
const (
	VersionSeverityCritical = "critical"
	VersionSeverityHigh     = "high"
	VersionSeverityMedium   = "medium"
	VersionSeverityLow      = "low"
)

// DefaultVersionConfig returns default version configuration
func DefaultVersionConfig() *VersionConfig {
	return &VersionConfig{
		AutoUpdate:          false,
		UpdateChannel:       UpdateChannelStable,
		CheckInterval:       24 * time.Hour,
		UpdateNotifications: true,
		CacheVersions:       true,
		CacheTTL:            6 * time.Hour,
		OfflineMode:         false,
		VerifySignatures:    true,
		TrustedSources:      []string{"github.com", "registry.npmjs.org", "proxy.golang.org"},
		AllowPrerelease:     false,
		PackageRegistries:   []string{"https://registry.npmjs.org", "https://proxy.golang.org"},
		PackageTimeout:      30 * time.Second,
		PackageRetries:      3,
	}
}
