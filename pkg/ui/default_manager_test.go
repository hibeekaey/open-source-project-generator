package ui

import (
	"testing"
)

func TestDefaultManager_GetDefaultsForProject(t *testing.T) {
	mockLogger := &MockLogger{}
	dm := NewDefaultManager(mockLogger)

	projectName := "test-project"
	defaults := dm.GetDefaultsForProject(projectName)

	if defaults == nil {
		t.Fatal("expected defaults to be non-nil")
	}

	// Check that we get some default license
	if defaults.License == "" {
		t.Error("expected default license to be set")
	}

	// The default license should be MIT (system default)
	if defaults.License != "MIT" {
		t.Errorf("expected default license to be 'MIT', got %q", defaults.License)
	}
}

func TestDefaultManager_ValidateDefaults(t *testing.T) {
	mockLogger := &MockLogger{}
	dm := NewDefaultManager(mockLogger)

	// Test with valid defaults
	dm.userDefaults = &UserDefaults{
		Author:  "John Doe",
		Email:   "john@example.com",
		License: "MIT",
	}

	err := dm.ValidateDefaults()
	if err != nil {
		t.Errorf("expected no error for valid defaults, got: %v", err)
	}

	// Test with invalid email
	dm.userDefaults.Email = "invalid-email"
	err = dm.ValidateDefaults()
	if err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestDefaultManager_GetDefaultSources(t *testing.T) {
	mockLogger := &MockLogger{}
	dm := NewDefaultManager(mockLogger)

	sources := dm.GetDefaultSources()

	expectedKeys := []string{"license", "author", "email", "organization"}
	for _, key := range expectedKeys {
		if _, exists := sources[key]; !exists {
			t.Errorf("expected source for %q to exist", key)
		}
	}

	// Check that license has system source
	if sources["license"].Source != DefaultSourceSystem {
		t.Errorf("expected license source to be %q, got %q", DefaultSourceSystem, sources["license"].Source)
	}
}

func TestDefaultManager_extractUsernameFromGitURL(t *testing.T) {
	mockLogger := &MockLogger{}
	dm := NewDefaultManager(mockLogger)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "GitHub SSH URL",
			url:      "git@github.com:user/repo.git",
			expected: "user",
		},
		{
			name:     "GitHub HTTPS URL",
			url:      "https://github.com/user/repo.git",
			expected: "user",
		},
		{
			name:     "GitHub HTTPS URL without .git",
			url:      "https://github.com/user/repo",
			expected: "user",
		},
		{
			name:     "Non-GitHub URL",
			url:      "https://gitlab.com/user/repo.git",
			expected: "",
		},
		{
			name:     "Invalid URL",
			url:      "not-a-url",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dm.extractUsernameFromGitURL(tt.url)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
