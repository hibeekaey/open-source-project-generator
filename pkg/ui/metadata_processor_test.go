package ui

import (
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestMetadataProcessor_ProcessMetadata(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Description:  "A test project for validation",
		Author:       "John Doe",
		Email:        "john@example.com",
		Organization: "Test Corp",
		License:      "MIT",
		Repository:   "https://github.com/user/test-project",
		GeneratedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	metadata, err := mp.ProcessMetadata(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test basic fields
	if metadata.ProjectName != "test-project" {
		t.Errorf("expected project name 'test-project', got %q", metadata.ProjectName)
	}

	if metadata.ProjectDescription != "A test project for validation" {
		t.Errorf("expected description to match, got %q", metadata.ProjectDescription)
	}

	// Test computed fields
	if metadata.ProjectSlug != "test-project" {
		t.Errorf("expected project slug 'test-project', got %q", metadata.ProjectSlug)
	}

	if metadata.ProjectTitle != "Test Project" {
		t.Errorf("expected project title 'Test Project', got %q", metadata.ProjectTitle)
	}

	if metadata.AuthorWithEmail != "John Doe <john@example.com>" {
		t.Errorf("expected author with email 'John Doe <john@example.com>', got %q", metadata.AuthorWithEmail)
	}

	// Test boolean flags
	if !metadata.HasRepository {
		t.Error("expected HasRepository to be true")
	}

	if !metadata.HasOrganization {
		t.Error("expected HasOrganization to be true")
	}

	// Test copyright
	if metadata.Copyright != "Copyright (c) 2024 Test Corp" {
		t.Errorf("expected copyright 'Copyright (c) 2024 Test Corp', got %q", metadata.Copyright)
	}

	// Test module name extraction
	if metadata.ModuleName != "github.com/user/test-project" {
		t.Errorf("expected module name 'github.com/user/test-project', got %q", metadata.ModuleName)
	}
}

func TestMetadataProcessor_ProcessMetadata_NilConfig(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	_, err := mp.ProcessMetadata(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestMetadataProcessor_createSlug(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic name",
			input:    "MyProject",
			expected: "myproject",
		},
		{
			name:     "name with spaces",
			input:    "My Awesome Project",
			expected: "my-awesome-project",
		},
		{
			name:     "name with underscores",
			input:    "my_awesome_project",
			expected: "my-awesome-project",
		},
		{
			name:     "name with mixed case and special chars",
			input:    "My-Awesome_Project 2024!",
			expected: "my-awesome-project-2024",
		},
		{
			name:     "name with multiple hyphens",
			input:    "my--awesome---project",
			expected: "my-awesome-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mp.createSlug(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMetadataProcessor_createTitle(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase",
			input:    "myAwesomeProject",
			expected: "My Awesome Project",
		},
		{
			name:     "PascalCase",
			input:    "MyAwesomeProject",
			expected: "My Awesome Project",
		},
		{
			name:     "kebab-case",
			input:    "my-awesome-project",
			expected: "My Awesome Project",
		},
		{
			name:     "snake_case",
			input:    "my_awesome_project",
			expected: "My Awesome Project",
		},
		{
			name:     "already spaced",
			input:    "My Awesome Project",
			expected: "My Awesome Project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mp.createTitle(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMetadataProcessor_extractModuleName(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	tests := []struct {
		name        string
		repository  string
		projectSlug string
		expected    string
	}{
		{
			name:        "GitHub HTTPS URL",
			repository:  "https://github.com/user/repo",
			projectSlug: "repo",
			expected:    "github.com/user/repo",
		},
		{
			name:        "GitHub HTTP URL",
			repository:  "http://github.com/user/repo",
			projectSlug: "repo",
			expected:    "github.com/user/repo",
		},
		{
			name:        "GitLab URL",
			repository:  "https://gitlab.com/user/repo",
			projectSlug: "repo",
			expected:    "gitlab.com/user/repo",
		},
		{
			name:        "Empty repository",
			repository:  "",
			projectSlug: "my-project",
			expected:    "my-project",
		},
		{
			name:        "Unknown host",
			repository:  "https://example.com/user/repo",
			projectSlug: "repo",
			expected:    "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mp.extractModuleName(tt.repository, tt.projectSlug)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMetadataProcessor_determineFileType(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Go file",
			filePath: "main.go",
			expected: "go",
		},
		{
			name:     "JavaScript file",
			filePath: "index.js",
			expected: "javascript",
		},
		{
			name:     "TypeScript file",
			filePath: "app.ts",
			expected: "typescript",
		},
		{
			name:     "Markdown file",
			filePath: "README.md",
			expected: "markdown",
		},
		{
			name:     "YAML file",
			filePath: "config.yml",
			expected: "yaml",
		},
		{
			name:     "Unknown file",
			filePath: "unknown.xyz",
			expected: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mp.determineFileType(tt.filePath)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMetadataProcessor_ValidateMetadata(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	// Test valid metadata
	validMetadata := &TemplateMetadata{
		ProjectName:   "test-project",
		Author:        "John Doe",
		ProjectSlug:   "test-project",
		CopyrightYear: "2024",
	}

	err := mp.ValidateMetadata(validMetadata)
	if err != nil {
		t.Errorf("expected no error for valid metadata, got: %v", err)
	}

	// Test nil metadata
	err = mp.ValidateMetadata(nil)
	if err == nil {
		t.Error("expected error for nil metadata")
	}

	// Test missing project name
	invalidMetadata := &TemplateMetadata{
		Author:        "John Doe",
		ProjectSlug:   "test-project",
		CopyrightYear: "2024",
	}

	err = mp.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("expected error for missing project name")
	}

	// Test missing author
	invalidMetadata = &TemplateMetadata{
		ProjectName:   "test-project",
		ProjectSlug:   "test-project",
		CopyrightYear: "2024",
	}

	err = mp.ValidateMetadata(invalidMetadata)
	if err == nil {
		t.Error("expected error for missing author")
	}
}

func TestMetadataProcessor_GetMetadataForTemplate(t *testing.T) {
	mockLogger := &MockLogger{}
	mp := NewMetadataProcessor(mockLogger)

	metadata := &TemplateMetadata{
		ProjectName:        "test-project",
		ProjectDescription: "A test project",
		Author:             "John Doe",
		Email:              "john@example.com",
		License:            "MIT",
	}

	templateData := mp.GetMetadataForTemplate(metadata)

	expectedKeys := []string{
		"Name", "Description", "Author", "Email", "License",
		"Slug", "Title", "Copyright", "Repository", "ModuleName",
	}

	for _, key := range expectedKeys {
		if _, exists := templateData[key]; !exists {
			t.Errorf("expected key %q to exist in template data", key)
		}
	}

	// Check specific values
	if templateData["Name"] != "test-project" {
		t.Errorf("expected Name to be 'test-project', got %v", templateData["Name"])
	}

	if templateData["Author"] != "John Doe" {
		t.Errorf("expected Author to be 'John Doe', got %v", templateData["Author"])
	}
}
