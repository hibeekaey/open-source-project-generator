package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockVersionManager implements version fetching with mock responses
type MockVersionManager struct {
	nodeVersion    string
	goVersion      string
	npmPackages    map[string]string
	goModules      map[string]string
	githubReleases map[string]string
	shouldError    bool
	errorMessage   string
}

func NewMockVersionManager() *MockVersionManager {
	return &MockVersionManager{
		nodeVersion: "20.11.0",
		goVersion:   "1.22.0",
		npmPackages: map[string]string{
			"react":       "18.2.0",
			"next":        "14.0.4",
			"typescript":  "5.3.3",
			"tailwindcss": "3.4.0",
		},
		goModules: map[string]string{
			"github.com/gin-gonic/gin":    "v1.9.1",
			"gorm.io/gorm":                "v1.25.5",
			"github.com/stretchr/testify": "v1.8.4",
		},
		githubReleases: map[string]string{
			"kubernetes/kubernetes": "v1.29.0",
			"docker/compose":        "v2.24.0",
			"terraform":             "v1.6.6",
		},
	}
}

func (m *MockVersionManager) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

func (m *MockVersionManager) GetLatestNodeVersion() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMessage)
	}
	return m.nodeVersion, nil
}

func (m *MockVersionManager) GetLatestGoVersion() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMessage)
	}
	return m.goVersion, nil
}

func (m *MockVersionManager) GetLatestNPMPackage(packageName string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMessage)
	}
	if version, exists := m.npmPackages[packageName]; exists {
		return version, nil
	}
	return "", fmt.Errorf("package %s not found", packageName)
}

func (m *MockVersionManager) GetLatestGoModule(moduleName string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMessage)
	}
	if version, exists := m.goModules[moduleName]; exists {
		return version, nil
	}
	return "", fmt.Errorf("module %s not found", moduleName)
}

func (m *MockVersionManager) GetLatestGitHubRelease(repo string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMessage)
	}
	if version, exists := m.githubReleases[repo]; exists {
		return version, nil
	}
	return "", fmt.Errorf("repository %s not found", repo)
}

func TestMockVersionManager(t *testing.T) {
	manager := NewMockVersionManager()

	t.Run("successful version fetching", func(t *testing.T) {
		// Test Node.js version
		nodeVersion, err := manager.GetLatestNodeVersion()
		if err != nil {
			t.Errorf("GetLatestNodeVersion failed: %v", err)
		}
		if nodeVersion != "20.11.0" {
			t.Errorf("Expected Node version 20.11.0, got %s", nodeVersion)
		}

		// Test Go version
		goVersion, err := manager.GetLatestGoVersion()
		if err != nil {
			t.Errorf("GetLatestGoVersion failed: %v", err)
		}
		if goVersion != "1.22.0" {
			t.Errorf("Expected Go version 1.22.0, got %s", goVersion)
		}

		// Test NPM package
		reactVersion, err := manager.GetLatestNPMPackage("react")
		if err != nil {
			t.Errorf("GetLatestNPMPackage failed: %v", err)
		}
		if reactVersion != "18.2.0" {
			t.Errorf("Expected React version 18.2.0, got %s", reactVersion)
		}

		// Test Go module
		ginVersion, err := manager.GetLatestGoModule("github.com/gin-gonic/gin")
		if err != nil {
			t.Errorf("GetLatestGoModule failed: %v", err)
		}
		if ginVersion != "v1.9.1" {
			t.Errorf("Expected Gin version v1.9.1, got %s", ginVersion)
		}

		// Test GitHub release
		k8sVersion, err := manager.GetLatestGitHubRelease("kubernetes/kubernetes")
		if err != nil {
			t.Errorf("GetLatestGitHubRelease failed: %v", err)
		}
		if k8sVersion != "v1.29.0" {
			t.Errorf("Expected Kubernetes version v1.29.0, got %s", k8sVersion)
		}
	})

	t.Run("error handling", func(t *testing.T) {
		manager.SetError(true, "network timeout")

		_, err := manager.GetLatestNodeVersion()
		if err == nil {
			t.Error("Expected error for Node version fetch")
		}

		_, err = manager.GetLatestNPMPackage("react")
		if err == nil {
			t.Error("Expected error for NPM package fetch")
		}

		_, err = manager.GetLatestGoModule("github.com/gin-gonic/gin")
		if err == nil {
			t.Error("Expected error for Go module fetch")
		}
	})

	t.Run("non-existent packages", func(t *testing.T) {
		manager.SetError(false, "")

		_, err := manager.GetLatestNPMPackage("non-existent-package")
		if err == nil {
			t.Error("Expected error for non-existent NPM package")
		}

		_, err = manager.GetLatestGoModule("non-existent/module")
		if err == nil {
			t.Error("Expected error for non-existent Go module")
		}

		_, err = manager.GetLatestGitHubRelease("non-existent/repo")
		if err == nil {
			t.Error("Expected error for non-existent GitHub repository")
		}
	})
}

func TestHTTPMockServers(t *testing.T) {
	t.Run("NPM registry mock", func(t *testing.T) {
		// Create mock NPM registry server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/react" {
				response := map[string]interface{}{
					"dist-tags": map[string]string{
						"latest": "18.2.0",
					},
					"versions": map[string]interface{}{
						"18.2.0": map[string]interface{}{
							"version": "18.2.0",
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			} else {
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		// Test HTTP client against mock server
		resp, err := http.Get(server.URL + "/react")
		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var npmResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&npmResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		distTags, ok := npmResponse["dist-tags"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected dist-tags in response")
		}

		latest, ok := distTags["latest"].(string)
		if !ok || latest != "18.2.0" {
			t.Errorf("Expected latest version 18.2.0, got %v", latest)
		}
	})

	t.Run("GitHub API mock", func(t *testing.T) {
		// Create mock GitHub API server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/repos/kubernetes/kubernetes/releases/latest" {
				response := map[string]interface{}{
					"tag_name":     "v1.29.0",
					"name":         "v1.29.0",
					"published_at": time.Now().Format(time.RFC3339),
				}
				json.NewEncoder(w).Encode(response)
			} else {
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		// Test HTTP client against mock server
		resp, err := http.Get(server.URL + "/repos/kubernetes/kubernetes/releases/latest")
		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var githubResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&githubResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		tagName, ok := githubResponse["tag_name"].(string)
		if !ok || tagName != "v1.29.0" {
			t.Errorf("Expected tag_name v1.29.0, got %v", tagName)
		}
	})

	t.Run("Go proxy mock", func(t *testing.T) {
		// Create mock Go proxy server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/github.com/gin-gonic/gin/@latest" {
				response := map[string]interface{}{
					"Version": "v1.9.1",
					"Time":    time.Now().Format(time.RFC3339),
				}
				json.NewEncoder(w).Encode(response)
			} else {
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		// Test HTTP client against mock server
		resp, err := http.Get(server.URL + "/github.com/gin-gonic/gin/@latest")
		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var goResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&goResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		version, ok := goResponse["Version"].(string)
		if !ok || version != "v1.9.1" {
			t.Errorf("Expected Version v1.9.1, got %v", version)
		}
	})
}

func TestVersionManagerErrorScenarios(t *testing.T) {
	t.Run("network timeout simulation", func(t *testing.T) {
		// Create server that delays response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Simulate slow response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"dist-tags": {"latest": "1.0.0"}}`))
		}))
		defer server.Close()

		// Test with very short timeout
		client := &http.Client{
			Timeout: 10 * time.Millisecond, // Very short timeout
		}

		_, err := client.Get(server.URL)
		if err == nil {
			t.Error("Expected timeout error")
		}
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json {`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err == nil {
			t.Error("Expected JSON decode error")
		}
	})

	t.Run("HTTP error codes", func(t *testing.T) {
		errorCodes := []int{
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
		}

		for _, code := range errorCodes {
			t.Run(fmt.Sprintf("HTTP %d", code), func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(code)
					w.Write([]byte(fmt.Sprintf("HTTP %d Error", code)))
				}))
				defer server.Close()

				resp, err := http.Get(server.URL)
				if err != nil {
					t.Fatalf("HTTP request failed: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != code {
					t.Errorf("Expected status %d, got %d", code, resp.StatusCode)
				}
			})
		}
	})

	t.Run("rate limiting simulation", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			if requestCount <= 3 {
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded"))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"dist-tags": {"latest": "1.0.0"}}`))
			}
		}))
		defer server.Close()

		// Make multiple requests to trigger rate limiting
		for i := 0; i < 5; i++ {
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("HTTP request %d failed: %v", i, err)
			}

			if i < 3 {
				if resp.StatusCode != http.StatusTooManyRequests {
					t.Errorf("Request %d: expected rate limit, got status %d", i, resp.StatusCode)
				}
			} else {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Request %d: expected success after rate limit, got status %d", i, resp.StatusCode)
				}
			}
			resp.Body.Close()
		}
	})
}

func TestVersionManagerPerformance(t *testing.T) {
	t.Run("concurrent version fetching", func(t *testing.T) {
		manager := NewMockVersionManager()
		const numGoroutines = 100

		results := make(chan error, numGoroutines)
		start := time.Now()

		// Concurrent version fetches
		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := manager.GetLatestNodeVersion()
				results <- err
			}()
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent version fetch failed: %v", err)
			}
		}

		duration := time.Since(start)
		t.Logf("Completed %d concurrent version fetches in %v", numGoroutines, duration)

		// Should complete quickly with mock
		if duration > 100*time.Millisecond {
			t.Errorf("Concurrent version fetching too slow: %v", duration)
		}
	})

	t.Run("version fetching with caching", func(t *testing.T) {
		cache := NewMemoryCache(1 * time.Hour)

		// Pre-populate cache
		cache.Set("node-version", "20.11.0")
		cache.Set("npm:react", "18.2.0")
		cache.Set("go:github.com/gin-gonic/gin", "v1.9.1")

		start := time.Now()

		// Fetch from cache multiple times
		for i := 0; i < 1000; i++ {
			cache.Get("node-version")
			cache.Get("npm:react")
			cache.Get("go:github.com/gin-gonic/gin")
		}

		duration := time.Since(start)
		t.Logf("Completed 3000 cache lookups in %v", duration)

		// Cache lookups should be very fast
		if duration > 10*time.Millisecond {
			t.Errorf("Cache lookups too slow: %v", duration)
		}
	})
}
