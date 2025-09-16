package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGitHubClient_GetLatestRelease(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		mockResponse   *GitHubRelease
		expectedStatus int
		expectedError  bool
		expectedResult string
	}{
		{
			name:  "successful request",
			owner: "golang",
			repo:  "go",
			mockResponse: &GitHubRelease{
				TagName:     "go1.22.0",
				Name:        "Go 1.22.0",
				Draft:       false,
				Prerelease:  false,
				PublishedAt: time.Now(),
				Body:        "Release notes...",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: "go1.22.0",
		},
		{
			name:           "repository not found",
			owner:          "nonexistent",
			repo:           "repo",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:  "empty tag name",
			owner: "example",
			repo:  "repo",
			mockResponse: &GitHubRelease{
				TagName:     "",
				Name:        "Release",
				Draft:       false,
				Prerelease:  false,
				PublishedAt: time.Now(),
			},
			expectedStatus: http.StatusOK,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
					t.Errorf("Expected Accept header to be application/vnd.github.v3+json, got %s", r.Header.Get("Accept"))
				}
				if r.Header.Get("User-Agent") != "open-source-project-generator/1.0" {
					t.Errorf("Expected User-Agent header to be open-source-project-generator/1.0, got %s", r.Header.Get("User-Agent"))
				}

				// Verify URL path
				expectedPath := "/repos/" + tt.owner + "/" + tt.repo + "/releases/latest"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK && tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewGitHubClient(server.Client())
			client.baseURL = server.URL

			// Execute test
			result, err := client.GetLatestRelease(tt.owner, tt.repo)

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result %s, got %s", tt.expectedResult, result)
				}
			}
		})
	}
}

func TestGitHubClient_GetReleases(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		limit          int
		mockResponse   []GitHubRelease
		expectedStatus int
		expectedError  bool
		expectedCount  int
	}{
		{
			name:  "successful request with limit",
			owner: "golang",
			repo:  "go",
			limit: 2,
			mockResponse: []GitHubRelease{
				{
					TagName:     "go1.22.0",
					Name:        "Go 1.22.0",
					Draft:       false,
					Prerelease:  false,
					PublishedAt: time.Now(),
				},
				{
					TagName:     "go1.21.0",
					Name:        "Go 1.21.0",
					Draft:       false,
					Prerelease:  false,
					PublishedAt: time.Now().Add(-24 * time.Hour),
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name:           "repository not found",
			owner:          "nonexistent",
			repo:           "repo",
			limit:          0,
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/repos/" + tt.owner + "/" + tt.repo + "/releases"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Check query parameters
				if tt.limit > 0 {
					perPage := r.URL.Query().Get("per_page")
					if perPage == "" {
						t.Errorf("Expected per_page query parameter when limit is set")
					}
				}

				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewGitHubClient(server.Client())
			client.baseURL = server.URL

			// Execute test
			result, err := client.GetReleases(tt.owner, tt.repo, tt.limit)

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("Expected %d releases, got %d", tt.expectedCount, len(result))
				}
			}
		})
	}
}

func TestGitHubClient_GetLatestStableRelease(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		mockResponse   []GitHubRelease
		expectedStatus int
		expectedError  bool
		expectedResult string
	}{
		{
			name:  "successful request with stable release",
			owner: "golang",
			repo:  "go",
			mockResponse: []GitHubRelease{
				{
					TagName:    "go1.22.0-rc1",
					Name:       "Go 1.22.0 RC1",
					Draft:      false,
					Prerelease: true, // This should be skipped
				},
				{
					TagName:    "go1.21.5",
					Name:       "Go 1.21.5",
					Draft:      false,
					Prerelease: false, // This should be returned
				},
				{
					TagName:    "go1.21.4",
					Name:       "Go 1.21.4",
					Draft:      false,
					Prerelease: false,
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: "go1.21.5",
		},
		{
			name:  "only prereleases available",
			owner: "example",
			repo:  "repo",
			mockResponse: []GitHubRelease{
				{
					TagName:    "v1.0.0-beta1",
					Name:       "Beta 1",
					Draft:      false,
					Prerelease: true,
				},
				{
					TagName:    "v1.0.0-alpha1",
					Name:       "Alpha 1",
					Draft:      false,
					Prerelease: true,
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  true,
		},
		{
			name:           "repository not found",
			owner:          "nonexistent",
			repo:           "repo",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewGitHubClient(server.Client())
			client.baseURL = server.URL

			// Execute test
			result, err := client.GetLatestStableRelease(tt.owner, tt.repo)

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result %s, got %s", tt.expectedResult, result)
				}
			}
		})
	}
}

func TestGitHubClient_SetAPIToken(t *testing.T) {
	client := NewGitHubClient(nil)
	token := "test-token"

	client.SetAPIToken(token)

	if client.apiToken != token {
		t.Errorf("Expected API token to be %s, got %s", token, client.apiToken)
	}
}

func TestGitHubClient_AuthorizationHeader(t *testing.T) {
	token := "test-token"

	// Create mock server that checks for authorization header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "token " + token
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %s, got %s", expectedAuth, authHeader)
		}

		// Return a valid response
		release := GitHubRelease{
			TagName: "v1.0.0",
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// Create client with token
	client := NewGitHubClient(server.Client())
	client.baseURL = server.URL
	client.SetAPIToken(token)

	// Make request
	_, err := client.GetLatestRelease("test", "repo")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestNewGitHubClient(t *testing.T) {
	t.Run("with custom http client", func(t *testing.T) {
		customClient := &http.Client{}
		client := NewGitHubClient(customClient)

		if client.httpClient != customClient {
			t.Errorf("Expected custom http client to be used")
		}
		if client.baseURL != "https://api.github.com" {
			t.Errorf("Expected baseURL to be https://api.github.com, got %s", client.baseURL)
		}
		if client.apiToken != "" {
			t.Errorf("Expected empty API token by default, got %s", client.apiToken)
		}
	})

	t.Run("with nil http client", func(t *testing.T) {
		client := NewGitHubClient(nil)

		if client.httpClient == nil {
			t.Errorf("Expected default http client to be created")
		}
		if client.baseURL != "https://api.github.com" {
			t.Errorf("Expected baseURL to be https://api.github.com, got %s", client.baseURL)
		}
	})
}
