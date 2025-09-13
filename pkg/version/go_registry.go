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
	// Implement Go vulnerability database integration
	vulnResult, err := r.performGoVulnCheck(moduleName, version)
	if err != nil {
		// Log error but don't fail the entire operation
		// Security checking is supplementary to version checking
		return []models.SecurityIssue{}, nil
	}

	// Ensure we always return a non-nil slice
	if vulnResult == nil {
		return []models.SecurityIssue{}, nil
	}

	return vulnResult, nil
}

// performGoVulnCheck performs vulnerability checking using Go vulnerability database
func (r *GoRegistry) performGoVulnCheck(moduleName, version string) ([]models.SecurityIssue, error) {
	// Use Go vulnerability database API to check for vulnerabilities
	// The Go vulnerability database is available at https://vuln.go.dev/
	vulnURL := "https://api.osv.dev/v1/query"

	// Create vulnerability query payload
	queryPayload := map[string]any{
		"package": map[string]any{
			"name":      moduleName,
			"ecosystem": "Go",
		},
		"version": version,
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vulnerability query: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", vulnURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create vulnerability request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	// Use a shorter timeout for security checks to avoid blocking
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vulnerability check: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// If vulnerability API is unavailable, return empty result rather than failing
	if resp.StatusCode != http.StatusOK {
		return []models.SecurityIssue{}, nil
	}

	// Parse vulnerability response
	var vulnResponse GoVulnResponse
	if err := json.NewDecoder(resp.Body).Decode(&vulnResponse); err != nil {
		return nil, fmt.Errorf("failed to decode vulnerability response: %w", err)
	}

	// Convert vulnerability results to SecurityIssue format
	return r.convertVulnToSecurityIssues(vulnResponse), nil
}

// GoVulnResponse represents the response from Go vulnerability database
type GoVulnResponse struct {
	Vulns []GoVulnerability `json:"vulns"`
}

// GoVulnerability represents a vulnerability from Go vulnerability database
type GoVulnerability struct {
	ID         string            `json:"id"`
	Summary    string            `json:"summary"`
	Details    string            `json:"details"`
	Aliases    []string          `json:"aliases"`
	Modified   string            `json:"modified"`
	Published  string            `json:"published"`
	Affected   []GoAffectedRange `json:"affected"`
	References []GoReference     `json:"references"`
	Severity   []GoSeverity      `json:"severity"`
}

// GoAffectedRange represents affected version ranges
type GoAffectedRange struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Ranges []struct {
		Type   string `json:"type"`
		Events []struct {
			Introduced string `json:"introduced,omitempty"`
			Fixed      string `json:"fixed,omitempty"`
		} `json:"events"`
	} `json:"ranges"`
}

// GoReference represents a reference link
type GoReference struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// GoSeverity represents severity information
type GoSeverity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

// convertVulnToSecurityIssues converts Go vulnerability results to SecurityIssue format
func (r *GoRegistry) convertVulnToSecurityIssues(vulnResponse GoVulnResponse) []models.SecurityIssue {
	var issues []models.SecurityIssue

	for _, vuln := range vulnResponse.Vulns {
		// Parse published time
		publishedAt, err := time.Parse("2006-01-02T15:04:05Z", vuln.Published)
		if err != nil {
			publishedAt = time.Now() // Fallback to current time
		}

		// Determine severity from CVSS score or use default
		severity := r.determineGoSeverity(vuln.Severity)

		// Find fixed version from affected ranges
		fixedIn := r.extractFixedVersion(vuln.Affected)

		// Get reference URL
		refURL := r.extractReferenceURL(vuln.References)

		issue := models.SecurityIssue{
			ID:          vuln.ID,
			Severity:    severity,
			Description: fmt.Sprintf("%s: %s", vuln.Summary, vuln.Details),
			FixedIn:     fixedIn,
			ReportedAt:  publishedAt,
			URL:         refURL,
		}

		issues = append(issues, issue)
	}

	return issues
}

// determineGoSeverity determines severity level from Go vulnerability data
func (r *GoRegistry) determineGoSeverity(severities []GoSeverity) string {
	for _, sev := range severities {
		if sev.Type == "CVSS_V3" {
			// Parse CVSS score to determine severity
			score := sev.Score
			if strings.Contains(score, "9.") || strings.Contains(score, "10.") {
				return "critical"
			} else if strings.Contains(score, "7.") || strings.Contains(score, "8.") {
				return "high"
			} else if strings.Contains(score, "4.") || strings.Contains(score, "5.") || strings.Contains(score, "6.") {
				return "medium"
			} else {
				return "low"
			}
		}
	}
	return "medium" // Default to medium for unknown severities
}

// extractFixedVersion extracts the fixed version from affected ranges
func (r *GoRegistry) extractFixedVersion(affected []GoAffectedRange) string {
	for _, aff := range affected {
		for _, rng := range aff.Ranges {
			for _, event := range rng.Events {
				if event.Fixed != "" {
					return event.Fixed
				}
			}
		}
	}
	return "" // No fixed version found
}

// extractReferenceURL extracts a reference URL from references
func (r *GoRegistry) extractReferenceURL(references []GoReference) string {
	for _, ref := range references {
		if ref.URL != "" {
			return ref.URL
		}
	}
	return "" // No URL found
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
