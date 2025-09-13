package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
	apiToken   string
}

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(httpClient *http.Client) *GitHubClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &GitHubClient{
		httpClient: httpClient,
		baseURL:    "https://api.github.com",
		// API token can be set via environment variable if needed for higher rate limits
		apiToken: "", // Will be set from environment if available
	}
}

// SetAPIToken sets the GitHub API token for authenticated requests
func (c *GitHubClient) SetAPIToken(token string) {
	c.apiToken = token
}

// GetLatestRelease fetches the latest release version from GitHub
func (c *GitHubClient) GetLatestRelease(owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", c.baseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	// Add authorization header if token is available
	if c.apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.apiToken))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("repository %s/%s not found or has no releases", owner, repo)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no tag name found in latest release for %s/%s", owner, repo)
	}

	return release.TagName, nil
}

// GetReleases fetches all releases for a repository
func (c *GitHubClient) GetReleases(owner, repo string, limit int) ([]GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases", c.baseURL, owner, repo)
	if limit > 0 {
		url = fmt.Sprintf("%s?per_page=%d", url, limit)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "open-source-template-generator/1.0")

	// Add authorization header if token is available
	if c.apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.apiToken))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return releases, nil
}

// GetLatestStableRelease fetches the latest non-prerelease version
func (c *GitHubClient) GetLatestStableRelease(owner, repo string) (string, error) {
	releases, err := c.GetReleases(owner, repo, 10) // Get first 10 releases
	if err != nil {
		return "", err
	}

	for _, release := range releases {
		if !release.Draft && !release.Prerelease && release.TagName != "" {
			return release.TagName, nil
		}
	}

	return "", fmt.Errorf("no stable releases found for %s/%s", owner, repo)
}
