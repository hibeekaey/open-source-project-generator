package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NPMClient handles NPM registry API interactions
type NPMClient struct {
	httpClient *http.Client
	baseURL    string
}

// NPMPackageResponse represents the NPM registry response
type NPMPackageResponse struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]NPMVersionInfo `json:"versions"`
}

// NPMVersionInfo contains version-specific information
type NPMVersionInfo struct {
	Version string `json:"version"`
}

// NewNPMClient creates a new NPM registry client
func NewNPMClient(httpClient *http.Client) *NPMClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &NPMClient{
		httpClient: httpClient,
		baseURL:    "https://registry.npmjs.org",
	}
}

// SetBaseURL sets the base URL for the NPM client (for testing purposes)
func (c *NPMClient) SetBaseURL(url string) {
	c.baseURL = url
}

// GetLatestVersion fetches the latest version of an NPM package
func (c *NPMClient) GetLatestVersion(packageName string) (string, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, packageName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NPM registry returned status %d for package %s", resp.StatusCode, packageName)
	}

	var packageInfo NPMPackageResponse
	if err := json.NewDecoder(resp.Body).Decode(&packageInfo); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if packageInfo.DistTags.Latest == "" {
		return "", fmt.Errorf("no latest version found for package %s", packageName)
	}

	return packageInfo.DistTags.Latest, nil
}

// GetVersionInfo fetches detailed information about a specific version
func (c *NPMClient) GetVersionInfo(packageName, version string) (*NPMVersionInfo, error) {
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, packageName, version)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NPM registry returned status %d for package %s@%s", resp.StatusCode, packageName, version)
	}

	var versionInfo NPMVersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &versionInfo, nil
}
