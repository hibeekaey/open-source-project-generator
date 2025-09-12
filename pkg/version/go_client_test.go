package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGoClient_GetLatestVersion(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		mockLatest     *GoModuleInfo
		mockVersions   []string
		latestStatus   int
		versionsStatus int
		expectedError  bool
		expectedResult string
	}{
		{
			name:       "successful @latest request",
			moduleName: "github.com/gin-gonic/gin",
			mockLatest: &GoModuleInfo{
				Version: "v1.9.1",
				Time:    time.Now(),
			},
			latestStatus:   http.StatusOK,
			expectedError:  false,
			expectedResult: "v1.9.1",
		},
		{
			name:           "@latest fails, fallback to version list",
			moduleName:     "github.com/example/module",
			mockVersions:   []string{"v1.0.0", "v1.1.0", "v1.2.0"},
			latestStatus:   http.StatusNotFound,
			versionsStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: "v1.2.0",
		},
		{
			name:           "both @latest and version list fail",
			moduleName:     "github.com/nonexistent/module",
			latestStatus:   http.StatusNotFound,
			versionsStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "empty version list",
			moduleName:     "github.com/empty/module",
			mockVersions:   []string{},
			latestStatus:   http.StatusNotFound,
			versionsStatus: http.StatusOK,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				encodedModule := encodeModuleName(tt.moduleName)

				if strings.HasSuffix(r.URL.Path, "/@latest") {
					w.WriteHeader(tt.latestStatus)
					if tt.latestStatus == http.StatusOK && tt.mockLatest != nil {
						json.NewEncoder(w).Encode(tt.mockLatest)
					}
				} else if strings.HasSuffix(r.URL.Path, "/@v/list") {
					w.WriteHeader(tt.versionsStatus)
					if tt.versionsStatus == http.StatusOK {
						w.Write([]byte(strings.Join(tt.mockVersions, "\n")))
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}

				// Verify encoded module name in path
				if !strings.Contains(r.URL.Path, encodedModule) {
					t.Errorf("Expected encoded module %s in path %s", encodedModule, r.URL.Path)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewGoClient(server.Client())
			client.BaseURL = server.URL

			// Execute test
			result, err := client.GetLatestVersion(tt.moduleName)

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

func TestGoClient_GetVersionInfo(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		version        string
		mockResponse   *GoModuleInfo
		expectedStatus int
		expectedError  bool
	}{
		{
			name:       "successful request",
			moduleName: "github.com/gin-gonic/gin",
			version:    "v1.9.1",
			mockResponse: &GoModuleInfo{
				Version: "v1.9.1",
				Time:    time.Now(),
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "version not found",
			moduleName:     "github.com/gin-gonic/gin",
			version:        "v999.999.999",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				encodedModule := encodeModuleName(tt.moduleName)
				expectedPath := "/" + encodedModule + "/@v/" + tt.version + ".info"

				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK && tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewGoClient(server.Client())
			client.BaseURL = server.URL

			// Execute test
			result, err := client.GetVersionInfo(tt.moduleName, tt.version)

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.Version != tt.mockResponse.Version {
					t.Errorf("Expected version %s, got %s", tt.mockResponse.Version, result.Version)
				}
			}
		})
	}
}

func TestEncodeModuleName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "github.com/gin-gonic/gin",
			expected: "github.com/gin-gonic/gin",
		},
		{
			input:    "github.com/Azure/azure-sdk-for-go",
			expected: "github.com/!azure/azure-sdk-for-go",
		},
		{
			input:    "github.com/GoogleCloudPlatform/functions-framework-go",
			expected: "github.com/!google!cloud!platform/functions-framework-go",
		},
		{
			input:    "example.com/MyModule",
			expected: "example.com/!my!module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := encodeModuleName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewGoClient(t *testing.T) {
	t.Run("with custom http client", func(t *testing.T) {
		customClient := &http.Client{}
		client := NewGoClient(customClient)

		if client.httpClient != customClient {
			t.Errorf("Expected custom http client to be used")
		}
		if client.BaseURL != "https://proxy.golang.org" {
			t.Errorf("Expected baseURL to be https://proxy.golang.org, got %s", client.BaseURL)
		}
	})

	t.Run("with nil http client", func(t *testing.T) {
		client := NewGoClient(nil)

		if client.httpClient == nil {
			t.Errorf("Expected default http client to be created")
		}
		if client.BaseURL != "https://proxy.golang.org" {
			t.Errorf("Expected baseURL to be https://proxy.golang.org, got %s", client.BaseURL)
		}
	})
}
