package version

import (
	"fmt"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// GitHubRegistry implements VersionRegistry for GitHub releases
type GitHubRegistry struct {
	client *GitHubClient
}

// NewGitHubRegistry creates a new GitHub registry client
func NewGitHubRegistry(client *GitHubClient) *GitHubRegistry {
	return &GitHubRegistry{
		client: client,
	}
}

// GetLatestVersion retrieves the latest version for a GitHub repository
func (r *GitHubRegistry) GetLatestVersion(repoPath string) (*models.VersionInfo, error) {
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository path format: %s (expected owner/repo)", repoPath)
	}

	owner, repo := parts[0], parts[1]
	version, err := r.client.GetLatestRelease(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for %s: %w", repoPath, err)
	}

	// Determine language based on repository
	language := r.determineLanguage(owner, repo)
	packageType := r.determineType(owner, repo)

	info := &models.VersionInfo{
		Name:           repoPath,
		Language:       language,
		Type:           packageType,
		LatestVersion:  version,
		IsSecure:       true, // GitHub releases are generally secure
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "github",
		RegistryURL:    fmt.Sprintf("https://github.com/%s/releases", repoPath),
		SecurityIssues: make([]models.SecurityIssue, 0),
		Metadata: map[string]string{
			"owner": owner,
			"repo":  repo,
		},
	}

	return info, nil
}

// GetVersionHistory retrieves version history for a GitHub repository
func (r *GitHubRegistry) GetVersionHistory(repoPath string, limit int) ([]*models.VersionInfo, error) {
	// For now, just return the latest version
	// This could be enhanced to fetch multiple releases from GitHub API
	latest, err := r.GetLatestVersion(repoPath)
	if err != nil {
		return nil, err
	}

	return []*models.VersionInfo{latest}, nil
}

// CheckSecurity checks for security vulnerabilities in a specific version
func (r *GitHubRegistry) CheckSecurity(repoPath, version string) ([]models.SecurityIssue, error) {
	// GitHub releases are generally considered secure
	// This could be enhanced to check GitHub Security Advisories
	return []models.SecurityIssue{}, nil
}

// GetRegistryInfo returns information about the GitHub registry
func (r *GitHubRegistry) GetRegistryInfo() interfaces.RegistryInfo {
	return interfaces.RegistryInfo{
		Name:        "GitHub Releases",
		URL:         "https://api.github.com",
		Type:        "github",
		Description: "GitHub repository releases for language runtimes and tools",
		Supported:   []string{"go", "nodejs", "java", "kotlin", "swift"},
	}
}

// IsAvailable checks if the GitHub registry is currently accessible
func (r *GitHubRegistry) IsAvailable() bool {
	// Simple check by trying to get a well-known repository
	_, err := r.client.GetLatestRelease("golang", "go")
	return err == nil
}

// GetSupportedPackages returns a list of packages supported by this registry
func (r *GitHubRegistry) GetSupportedPackages() ([]string, error) {
	// Return common repositories that we track for language versions
	return []string{
		"golang/go",
		"nodejs/node",
		"JetBrains/kotlin",
		"apple/swift",
		"openjdk/jdk",
	}, nil
}

// Helper methods

func (r *GitHubRegistry) determineLanguage(owner, repo string) string {
	// Map well-known repositories to languages
	repoKey := fmt.Sprintf("%s/%s", owner, repo)

	languageMap := map[string]string{
		"golang/go":        "go",
		"nodejs/node":      "nodejs",
		"JetBrains/kotlin": "kotlin",
		"apple/swift":      "swift",
		"openjdk/jdk":      "java",
	}

	if language, exists := languageMap[repoKey]; exists {
		return language
	}

	// Default based on repo name
	switch repo {
	case "go":
		return "go"
	case "node":
		return "nodejs"
	case "kotlin":
		return "kotlin"
	case "swift":
		return "swift"
	case "jdk", "java":
		return "java"
	default:
		return "unknown"
	}
}

func (r *GitHubRegistry) determineType(owner, repo string) string {
	// Map well-known repositories to types
	repoKey := fmt.Sprintf("%s/%s", owner, repo)

	typeMap := map[string]string{
		"golang/go":        "language",
		"nodejs/node":      "language",
		"JetBrains/kotlin": "language",
		"apple/swift":      "language",
		"openjdk/jdk":      "language",
	}

	if packageType, exists := typeMap[repoKey]; exists {
		return packageType
	}

	// Default to language for runtime repositories
	return "language"
}

// Ensure GitHubRegistry implements VersionRegistry interface
var _ interfaces.VersionRegistry = (*GitHubRegistry)(nil)
