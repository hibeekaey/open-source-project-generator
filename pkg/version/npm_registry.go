package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	// Implement npm security audit integration
	auditResult, err := r.performNPMAudit(packageName, version)
	if err != nil {
		// Log error but don't fail the entire operation
		// Security checking is supplementary to version checking
		return []models.SecurityIssue{}, nil
	}

	// Ensure we always return a non-nil slice
	if auditResult == nil {
		return []models.SecurityIssue{}, nil
	}

	return auditResult, nil
}

// performNPMAudit performs security audit using npm audit API
func (r *NPMRegistry) performNPMAudit(packageName, version string) ([]models.SecurityIssue, error) {
	// Use npm audit API to check for vulnerabilities
	// This integrates with the npm security advisory database
	auditURL := "https://registry.npmjs.org/-/npm/v1/security/audits"

	// Create audit request payload
	auditPayload := map[string]any{
		"name":    packageName,
		"version": version,
		"requires": map[string]string{
			packageName: version,
		},
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(auditPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal audit payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", auditURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create audit request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	// Use a shorter timeout for security checks to avoid blocking
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform npm audit: %w", err)
	}
	defer resp.Body.Close()

	// If audit API is unavailable, return empty result rather than failing
	if resp.StatusCode != http.StatusOK {
		return []models.SecurityIssue{}, nil
	}

	// Parse audit response
	var auditResponse NPMAuditResponse
	if err := json.NewDecoder(resp.Body).Decode(&auditResponse); err != nil {
		return nil, fmt.Errorf("failed to decode audit response: %w", err)
	}

	// Convert audit results to SecurityIssue format
	return r.convertAuditToSecurityIssues(auditResponse), nil
}

// NPMAuditResponse represents the response from npm audit API
type NPMAuditResponse struct {
	Advisories map[string]NPMAdvisory `json:"advisories"`
	Metadata   NPMAuditMetadata       `json:"metadata"`
}

// NPMAdvisory represents a security advisory from npm
type NPMAdvisory struct {
	ID                 int      `json:"id"`
	Title              string   `json:"title"`
	Severity           string   `json:"severity"`
	VulnerableVersions string   `json:"vulnerable_versions"`
	PatchedVersions    string   `json:"patched_versions"`
	Overview           string   `json:"overview"`
	References         []string `json:"references"`
	CreatedAt          string   `json:"created"`
	UpdatedAt          string   `json:"updated"`
	URL                string   `json:"url"`
}

// NPMAuditMetadata contains metadata about the audit
type NPMAuditMetadata struct {
	Vulnerabilities struct {
		Info     int `json:"info"`
		Low      int `json:"low"`
		Moderate int `json:"moderate"`
		High     int `json:"high"`
		Critical int `json:"critical"`
	} `json:"vulnerabilities"`
}

// convertAuditToSecurityIssues converts npm audit results to SecurityIssue format
func (r *NPMRegistry) convertAuditToSecurityIssues(auditResponse NPMAuditResponse) []models.SecurityIssue {
	var issues []models.SecurityIssue

	for _, advisory := range auditResponse.Advisories {
		// Parse created time
		createdAt, err := time.Parse("2006-01-02T15:04:05.000Z", advisory.CreatedAt)
		if err != nil {
			createdAt = time.Now() // Fallback to current time
		}

		// Map npm severity to our severity levels
		severity := r.mapNPMSeverity(advisory.Severity)

		issue := models.SecurityIssue{
			ID:          fmt.Sprintf("npm-%d", advisory.ID),
			Severity:    severity,
			Description: fmt.Sprintf("%s: %s", advisory.Title, advisory.Overview),
			FixedIn:     advisory.PatchedVersions,
			ReportedAt:  createdAt,
			URL:         advisory.URL,
		}

		issues = append(issues, issue)
	}

	return issues
}

// mapNPMSeverity maps npm severity levels to our standard severity levels
func (r *NPMRegistry) mapNPMSeverity(npmSeverity string) string {
	switch strings.ToLower(npmSeverity) {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "moderate":
		return "medium"
	case "low":
		return "low"
	case "info":
		return "low"
	default:
		return "medium" // Default to medium for unknown severities
	}
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
