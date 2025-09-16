package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNPMClient_GetLatestVersion(t *testing.T) {
	tests := []struct {
		name           string
		packageName    string
		mockResponse   NPMPackageResponse
		expectedStatus int
		expectedError  bool
		expectedResult string
	}{
		{
			name:        "successful request",
			packageName: "react",
			mockResponse: NPMPackageResponse{
				DistTags: struct {
					Latest string `json:"latest"`
				}{
					Latest: "18.2.0",
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: "18.2.0",
		},
		{
			name:           "package not found",
			packageName:    "nonexistent-package",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:        "empty latest version",
			packageName: "empty-package",
			mockResponse: NPMPackageResponse{
				DistTags: struct {
					Latest string `json:"latest"`
				}{
					Latest: "",
				},
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
				if r.Header.Get("Accept") != "application/json" {
					t.Errorf("Expected Accept header to be application/json, got %s", r.Header.Get("Accept"))
				}
				if r.Header.Get("User-Agent") != "open-source-template-generator/1.0" {
					t.Errorf("Expected User-Agent header to be open-source-template-generator/1.0, got %s", r.Header.Get("User-Agent"))
				}

				// Verify URL path
				expectedPath := "/" + tt.packageName
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewNPMClient(server.Client())
			client.baseURL = server.URL

			// Execute test
			result, err := client.GetLatestVersion(tt.packageName)

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

func TestNPMClient_GetVersionInfo(t *testing.T) {
	tests := []struct {
		name           string
		packageName    string
		version        string
		mockResponse   NPMVersionInfo
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "successful request",
			packageName: "react",
			version:     "18.2.0",
			mockResponse: NPMVersionInfo{
				Version: "18.2.0",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "version not found",
			packageName:    "react",
			version:        "999.999.999",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/" + tt.packageName + "/" + tt.version
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewNPMClient(server.Client())
			client.baseURL = server.URL

			// Execute test
			result, err := client.GetVersionInfo(tt.packageName, tt.version)

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

func TestNewNPMClient(t *testing.T) {
	t.Run("with custom http client", func(t *testing.T) {
		customClient := &http.Client{}
		client := NewNPMClient(customClient)

		if client.httpClient != customClient {
			t.Errorf("Expected custom http client to be used")
		}
		if client.baseURL != "https://registry.npmjs.org" {
			t.Errorf("Expected baseURL to be https://registry.npmjs.org, got %s", client.baseURL)
		}
	})

	t.Run("with nil http client", func(t *testing.T) {
		client := NewNPMClient(nil)

		if client.httpClient == nil {
			t.Errorf("Expected default http client to be created")
		}
		if client.baseURL != "https://registry.npmjs.org" {
			t.Errorf("Expected baseURL to be https://registry.npmjs.org, got %s", client.baseURL)
		}
	})
}
