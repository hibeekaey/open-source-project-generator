package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GoClient handles Go module proxy API interactions
type GoClient struct {
	httpClient *http.Client
	BaseURL    string
}

// GoModuleVersions represents the list of versions from the Go module proxy
type GoModuleVersions struct {
	Versions []string `json:"versions"`
}

// GoModuleInfo represents module information from the Go module proxy
type GoModuleInfo struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

// NewGoClient creates a new Go module proxy client
func NewGoClient(httpClient *http.Client) *GoClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &GoClient{
		httpClient: httpClient,
		BaseURL:    "https://proxy.golang.org",
	}
}

// GetLatestVersion fetches the latest version of a Go module
func (c *GoClient) GetLatestVersion(moduleName string) (string, error) {
	// Encode module name for URL (replace uppercase with !lowercase)
	encodedModule := encodeModuleName(moduleName)

	// First try to get the latest version directly
	url := fmt.Sprintf("%s/%s/@latest", c.BaseURL, encodedModule)

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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// If @latest fails, try to get the list of versions
		return c.getLatestFromVersionList(encodedModule)
	}

	var moduleInfo GoModuleInfo
	if err := json.NewDecoder(resp.Body).Decode(&moduleInfo); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return moduleInfo.Version, nil
}

// getLatestFromVersionList fetches all versions and returns the latest
func (c *GoClient) getLatestFromVersionList(encodedModule string) (string, error) {
	url := fmt.Sprintf("%s/%s/@v/list", c.BaseURL, encodedModule)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("go module proxy returned status %d for module %s", resp.StatusCode, encodedModule)
	}

	// The response is a plain text list of versions, one per line
	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	trimmedBody := strings.TrimSpace(string(body))
	if trimmedBody == "" {
		return "", fmt.Errorf("no versions found for module %s", encodedModule)
	}

	versions := strings.Split(trimmedBody, "\n")
	if len(versions) == 0 || (len(versions) == 1 && versions[0] == "") {
		return "", fmt.Errorf("no versions found for module %s", encodedModule)
	}

	// Return the last version in the list (typically the latest)
	return versions[len(versions)-1], nil
}

// GetVersionInfo fetches detailed information about a specific version
func (c *GoClient) GetVersionInfo(moduleName, version string) (*GoModuleInfo, error) {
	encodedModule := encodeModuleName(moduleName)
	url := fmt.Sprintf("%s/%s/@v/%s.info", c.BaseURL, encodedModule, version)

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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("go module proxy returned status %d for module %s@%s", resp.StatusCode, moduleName, version)
	}

	var moduleInfo GoModuleInfo
	if err := json.NewDecoder(resp.Body).Decode(&moduleInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &moduleInfo, nil
}

// encodeModuleName encodes a module name for use in Go module proxy URLs
// According to Go module proxy spec, uppercase letters are encoded as !lowercase
func encodeModuleName(moduleName string) string {
	var result strings.Builder
	for _, r := range moduleName {
		if r >= 'A' && r <= 'Z' {
			result.WriteRune('!')
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
