package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewGoRegistry(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{
			name:       "with custom http client",
			httpClient: &http.Client{Timeout: 10 * time.Second},
		},
		{
			name:       "with nil http client",
			httpClient: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goClient := NewGoClient(tt.httpClient)
			registry := NewGoRegistry(goClient)

			if registry == nil {
				t.Errorf("expected registry instance but got nil")
				return
			}

			if registry.client == nil {
				t.Errorf("expected go client to be set")
			}
		})
	}
}

func TestGoRegistry_GetLatestVersion(t *testing.T) {
	tests := []struct {
		name            string
		moduleName      string
		mockLatest      string
		mockVersions    string
		expectedError   bool
		expectedVersion string
	}{
		{
			name:       "successful @latest request",
			moduleName: "github.com/gin-gonic/gin",
			mockLatest: `{
				"Version": "v1.9.1",
				"Time": "2023-07-18T14:30:00Z"
			}`,
			expectedError:   false,
			expectedVersion: "v1.9.1",
		},
		{
			name:            "@latest fails, fallback to version list",
			moduleName:      "github.com/example/module",
			mockLatest:      "", // Will return 404
			mockVersions:    "v1.0.0\nv1.1.0\nv1.2.0",
			expectedError:   false,
			expectedVersion: "v1.2.0",
		},
		{
			name:          "both @latest and version list fail",
			moduleName:    "github.com/nonexistent/module",
			mockLatest:    "", // Will return 404
			mockVersions:  "", // Will return 404
			expectedError: true,
		},
		{
			name:          "empty version list",
			moduleName:    "github.com/empty/module",
			mockLatest:    "", // Will return 404
			mockVersions:  "", // Empty response
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "@latest") {
					if tt.mockLatest == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(tt.mockLatest))
					return
				}

				if strings.Contains(r.URL.Path, "@v/list") {
					if tt.mockVersions == "" && !tt.expectedError {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					w.Header().Set("Content-Type", "text/plain")
					w.Write([]byte(tt.mockVersions))
					return
				}

				// Mock version info endpoint
				if strings.Contains(r.URL.Path, ".info") {
					mockInfo := GoModuleInfo{
						Version: tt.expectedVersion,
						Time:    time.Now(),
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(mockInfo)
					return
				}

				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			goClient := NewGoClient(server.Client())
			// Override the base URL for testing
			goClient.BaseURL = server.URL
			registry := NewGoRegistry(goClient)

			versionInfo, err := registry.GetLatestVersion(tt.moduleName)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if versionInfo == nil {
				t.Errorf("expected version info but got nil")
				return
			}

			if versionInfo.LatestVersion != tt.expectedVersion {
				t.Errorf("expected latest version %s, got %s", tt.expectedVersion, versionInfo.LatestVersion)
			}

			if versionInfo.Name != tt.moduleName {
				t.Errorf("expected module name %s, got %s", tt.moduleName, versionInfo.Name)
			}

			if versionInfo.Language != "go" {
				t.Errorf("expected language go, got %s", versionInfo.Language)
			}
		})
	}
}

func TestGoRegistry_GetVersionHistory(t *testing.T) {
	mockVersions := "v1.0.0\nv1.1.0\nv1.2.0\nv2.0.0"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "@v/list") {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(mockVersions))
			return
		}

		if strings.Contains(r.URL.Path, ".info") {
			// Extract version from path
			pathParts := strings.Split(r.URL.Path, "/")
			versionPart := pathParts[len(pathParts)-1]
			version := strings.TrimSuffix(versionPart, ".info")

			mockInfo := GoModuleInfo{
				Version: version,
				Time:    time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockInfo)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	goClient := NewGoClient(server.Client())
	goClient.BaseURL = server.URL
	registry := NewGoRegistry(goClient)

	tests := []struct {
		name          string
		limit         int
		expectedCount int
	}{
		{
			name:          "no limit",
			limit:         0,
			expectedCount: 4,
		},
		{
			name:          "with limit",
			limit:         2,
			expectedCount: 2,
		},
		{
			name:          "limit larger than available",
			limit:         10,
			expectedCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions, err := registry.GetVersionHistory("github.com/example/module", tt.limit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(versions) != tt.expectedCount {
				t.Errorf("expected %d versions, got %d", tt.expectedCount, len(versions))
			}

			// Verify versions are sorted in descending order
			if len(versions) > 1 {
				for i := 0; i < len(versions)-1; i++ {
					v1, err1 := ParseSemVer(versions[i].CurrentVersion)
					v2, err2 := ParseSemVer(versions[i+1].CurrentVersion)
					if err1 == nil && err2 == nil {
						if v1.Compare(v2) <= 0 {
							t.Errorf("versions not sorted correctly: %s should be > %s",
								versions[i].CurrentVersion, versions[i+1].CurrentVersion)
						}
					}
				}
			}

			// Verify all versions have correct metadata
			for _, version := range versions {
				if version.Name != "github.com/example/module" {
					t.Errorf("expected module name github.com/example/module, got %s", version.Name)
				}
				if version.Language != "go" {
					t.Errorf("expected language go, got %s", version.Language)
				}
			}
		})
	}
}

func TestGoRegistry_CheckSecurity(t *testing.T) {
	// Test that CheckSecurity returns empty slice when no vulnerabilities found
	// The actual vulnerability API is external and may not be available in tests
	goClient := NewGoClient(nil)
	registry := NewGoRegistry(goClient)

	securityIssues, err := registry.CheckSecurity("github.com/example/module", "1.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty slice when no vulnerabilities or service unavailable
	if securityIssues == nil {
		t.Errorf("expected empty slice, got nil")
	}
}

func TestGoRegistry_GetRegistryInfo(t *testing.T) {
	goClient := NewGoClient(nil)
	registry := NewGoRegistry(goClient)

	info := registry.GetRegistryInfo()

	if info.Name != "Go Module Proxy" {
		t.Errorf("expected name 'Go Module Proxy', got %s", info.Name)
	}
	if info.Type != "go" {
		t.Errorf("expected type 'go', got %s", info.Type)
	}
	if info.URL != "https://proxy.golang.org" {
		t.Errorf("expected URL 'https://proxy.golang.org', got %s", info.URL)
	}
	if len(info.Supported) == 0 {
		t.Errorf("expected supported languages to be populated")
	}
}

func TestGoRegistry_IsAvailable(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse int
		expected       bool
	}{
		{
			name:           "registry available",
			serverResponse: http.StatusOK,
			expected:       true,
		},
		{
			name:           "registry unavailable",
			serverResponse: http.StatusInternalServerError,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverResponse)
			}))
			defer server.Close()

			goClient := NewGoClient(server.Client())
			goClient.BaseURL = server.URL
			registry := NewGoRegistry(goClient)

			available := registry.IsAvailable()
			if available != tt.expected {
				t.Errorf("expected availability %v, got %v", tt.expected, available)
			}
		})
	}
}

func TestGoRegistry_GetSupportedPackages(t *testing.T) {
	goClient := NewGoClient(nil)
	registry := NewGoRegistry(goClient)

	packages, err := registry.GetSupportedPackages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(packages) == 0 {
		t.Errorf("expected supported packages to be populated")
	}

	// Check for some common packages
	expectedPackages := []string{
		"github.com/gin-gonic/gin",
		"github.com/gorilla/mux",
		"gorm.io/gorm",
		"github.com/stretchr/testify",
	}
	for _, expected := range expectedPackages {
		found := false
		for _, pkg := range packages {
			if pkg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected package %s to be in supported packages", expected)
		}
	}
}

func TestGoRegistry_Caching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "@latest") {
			callCount++
			mockResponse := GoModuleInfo{
				Version: "v1.0.0",
				Time:    time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockResponse)
			return
		}

		if strings.Contains(r.URL.Path, ".info") {
			mockInfo := GoModuleInfo{
				Version: "v1.0.0",
				Time:    time.Now(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockInfo)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	goClient := NewGoClient(server.Client())
	goClient.BaseURL = server.URL
	registry := NewGoRegistry(goClient)

	// First call should hit the server
	_, err := registry.GetLatestVersion("github.com/example/module")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}

	// Second call should use cache
	_, err = registry.GetLatestVersion("github.com/example/module")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 server call (cached), got %d", callCount)
	}
}
