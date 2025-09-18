package template

import (
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MockTemplateEngine implements interfaces.TemplateEngine for testing
type MockTemplateEngine struct {
	processedTemplates map[string]*models.ProjectConfig
	processedDirs      map[string]string
	shouldError        bool
	errorMessage       string
}

// NewMockTemplateEngine creates a new mock template engine
func NewMockTemplateEngine() *MockTemplateEngine {
	return &MockTemplateEngine{
		processedTemplates: make(map[string]*models.ProjectConfig),
		processedDirs:      make(map[string]string),
		shouldError:        false,
	}
}

// ProcessTemplate processes a template file
func (m *MockTemplateEngine) ProcessTemplate(path string, config *models.ProjectConfig) ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}

	m.processedTemplates[path] = config
	return []byte("processed template content"), nil
}

// ProcessDirectory processes a template directory
func (m *MockTemplateEngine) ProcessDirectory(templatePath string, outputPath string, config *models.ProjectConfig) error {
	if m.shouldError {
		return fmt.Errorf("%s", m.errorMessage)
	}

	m.processedDirs[templatePath] = outputPath
	m.processedTemplates[templatePath] = config
	return nil
}

// LoadTemplate loads a template (not implemented for mock)
func (m *MockTemplateEngine) LoadTemplate(path string) (*template.Template, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return nil, fmt.Errorf("LoadTemplate not implemented in mock")
}

// RenderTemplate renders a template (not implemented for mock)
func (m *MockTemplateEngine) RenderTemplate(tmpl *template.Template, data any) ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return []byte("rendered content"), nil
}

// RegisterFunctions registers template functions (no-op for mock)
func (m *MockTemplateEngine) RegisterFunctions(funcMap template.FuncMap) {
	// No-op for mock
}

// SetError configures the mock to return errors
func (m *MockTemplateEngine) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}

// GetProcessedTemplates returns processed templates for testing
func (m *MockTemplateEngine) GetProcessedTemplates() map[string]*models.ProjectConfig {
	return m.processedTemplates
}

// GetProcessedDirs returns processed directories for testing
func (m *MockTemplateEngine) GetProcessedDirs() map[string]string {
	return m.processedDirs
}

func TestNewManager(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}
}

func TestManager_ListTemplates(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with empty filter
	templates, err := manager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	// Should return embedded templates
	if len(templates) == 0 {
		t.Error("Expected templates to be returned")
	}

	// Test with category filter
	templates, err = manager.ListTemplates(interfaces.TemplateFilter{
		Category: "backend",
	})
	if err != nil {
		t.Fatalf("ListTemplates with filter failed: %v", err)
	}

	// All returned templates should be backend
	for _, tmpl := range templates {
		if tmpl.Category != "backend" {
			t.Errorf("Expected backend template, got category: %s", tmpl.Category)
		}
	}
}

func TestManager_GetTemplateInfo(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing template
	info, err := manager.GetTemplateInfo("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateInfo failed: %v", err)
	}

	if info == nil {
		t.Fatal("Expected template info, got nil")
	}

	if info.Name != "go-gin" {
		t.Errorf("Expected name 'go-gin', got '%s'", info.Name)
	}

	// Test with non-existent template
	_, err = manager.GetTemplateInfo("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_SearchTemplates(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test search with query
	templates, err := manager.SearchTemplates("go")
	if err != nil {
		t.Fatalf("SearchTemplates failed: %v", err)
	}

	// Should find templates containing "go"
	found := false
	for _, tmpl := range templates {
		if strings.Contains(strings.ToLower(tmpl.Name), "go") ||
			strings.Contains(strings.ToLower(tmpl.Technology), "go") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find templates matching 'go'")
	}

	// Test empty search
	templates, err = manager.SearchTemplates("")
	if err != nil {
		t.Fatalf("SearchTemplates with empty query failed: %v", err)
	}

	// Should return all templates
	if len(templates) == 0 {
		t.Error("Expected templates for empty search")
	}
}

func TestManager_ValidateTemplate(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test with non-existent path
	result, err := manager.ValidateTemplate("/non/existent/path")
	if err != nil {
		t.Fatalf("ValidateTemplate failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result for non-existent path")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected issues for non-existent path")
	}

	// Test with existing directory
	result, err = manager.ValidateTemplate(tempDir)
	if err != nil {
		t.Fatalf("ValidateTemplate failed: %v", err)
	}

	// Should have some validation results
	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}
}

func TestManager_ValidateTemplateStructure(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	tests := []struct {
		name        string
		template    *interfaces.TemplateInfo
		expectError bool
	}{
		{
			name: "valid template",
			template: &interfaces.TemplateInfo{
				Name:     "test-template",
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: false,
		},
		{
			name: "missing name",
			template: &interfaces.TemplateInfo{
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: true,
		},
		{
			name: "missing category",
			template: &interfaces.TemplateInfo{
				Name:    "test-template",
				Version: "1.0.0",
			},
			expectError: true,
		},
		{
			name: "invalid category",
			template: &interfaces.TemplateInfo{
				Name:     "test-template",
				Category: "invalid",
				Version:  "1.0.0",
			},
			expectError: true,
		},
		{
			name: "invalid name format",
			template: &interfaces.TemplateInfo{
				Name:     "Test Template", // Should be kebab-case
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateTemplateStructure(tt.template)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestManager_ProcessTemplate(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		OutputPath:   "/test/output",
	}

	// Test processing existing template
	err := manager.ProcessTemplate("go-gin", config, "/test/output")
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	// Verify template was processed
	processedDirs := mockEngine.GetProcessedDirs()
	if len(processedDirs) == 0 {
		t.Error("Expected template to be processed")
	}

	// Test processing non-existent template
	err = manager.ProcessTemplate("non-existent", config, "/test/output")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_ProcessCustomTemplate(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		OutputPath:   "/test/output",
	}

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test processing custom template
	err := manager.ProcessCustomTemplate(tempDir, config, "/test/output")
	if err != nil {
		t.Fatalf("ProcessCustomTemplate failed: %v", err)
	}

	// Verify template was processed
	processedDirs := mockEngine.GetProcessedDirs()
	if _, exists := processedDirs[tempDir]; !exists {
		t.Error("Expected custom template to be processed")
	}
}

func TestManager_GetTemplateMetadata(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing template
	metadata, err := manager.GetTemplateMetadata("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateMetadata failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Expected metadata, got nil")
	}

	if metadata.Author == "" {
		t.Error("Expected author in metadata")
	}

	if metadata.License == "" {
		t.Error("Expected license in metadata")
	}

	// Test with non-existent template
	_, err = manager.GetTemplateMetadata("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_GetTemplateDependencies(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing embedded template
	deps, err := manager.GetTemplateDependencies("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateDependencies failed: %v", err)
	}

	// Dependencies should be a slice (may be empty since YAML parsing is not implemented yet)
	if deps == nil {
		t.Error("Expected dependencies slice, got nil")
	}

	// Test with non-existent template
	_, err = manager.GetTemplateDependencies("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_GetTemplateCompatibility(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing template
	compat, err := manager.GetTemplateCompatibility("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateCompatibility failed: %v", err)
	}

	if compat == nil {
		t.Fatal("Expected compatibility info, got nil")
	}

	if len(compat.SupportedPlatforms) == 0 {
		t.Error("Expected supported platforms")
	}

	// Test with non-existent template
	_, err = manager.GetTemplateCompatibility("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_GetTemplatesByCategory(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with valid category
	templates, err := manager.GetTemplatesByCategory("backend")
	if err != nil {
		t.Fatalf("GetTemplatesByCategory failed: %v", err)
	}

	// All templates should be backend
	for _, tmpl := range templates {
		if tmpl.Category != "backend" {
			t.Errorf("Expected backend template, got category: %s", tmpl.Category)
		}
	}

	// Test with invalid category
	templates, err = manager.GetTemplatesByCategory("invalid")
	if err != nil {
		t.Fatalf("GetTemplatesByCategory failed: %v", err)
	}

	// Should return empty slice
	if len(templates) > 0 {
		t.Error("Expected no templates for invalid category")
	}
}

func TestManager_GetTemplatesByTechnology(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with valid technology
	templates, err := manager.GetTemplatesByTechnology("Go")
	if err != nil {
		t.Fatalf("GetTemplatesByTechnology failed: %v", err)
	}

	// All templates should use Go technology
	for _, tmpl := range templates {
		if tmpl.Technology != "Go" {
			t.Errorf("Expected Go template, got technology: %s", tmpl.Technology)
		}
	}
}

func TestManager_ValidateTemplateMetadata(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	tests := []struct {
		name        string
		metadata    *interfaces.TemplateMetadata
		expectError bool
	}{
		{
			name: "valid metadata",
			metadata: &interfaces.TemplateMetadata{
				Author:  "Test Author",
				License: "MIT",
			},
			expectError: false,
		},
		{
			name: "missing author",
			metadata: &interfaces.TemplateMetadata{
				License: "MIT",
			},
			expectError: true,
		},
		{
			name: "missing license",
			metadata: &interfaces.TemplateMetadata{
				Author: "Test Author",
			},
			expectError: true,
		},
		{
			name:        "nil metadata",
			metadata:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateTemplateMetadata(tt.metadata)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestManager_PreviewTemplate(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	config := &models.ProjectConfig{
		Name:         "preview-project",
		Organization: "preview-org",
		OutputPath:   "/preview/output",
	}

	// Test with existing template
	preview, err := manager.PreviewTemplate("go-gin", config)
	if err != nil {
		t.Fatalf("PreviewTemplate failed: %v", err)
	}

	if preview == nil {
		t.Fatal("Expected preview, got nil")
	}

	if preview.TemplateName != "go-gin" {
		t.Errorf("Expected template name 'go-gin', got '%s'", preview.TemplateName)
	}

	if preview.OutputPath != config.OutputPath {
		t.Errorf("Expected output path '%s', got '%s'", config.OutputPath, preview.OutputPath)
	}

	// Should have variables from config
	if preview.Variables["Name"] != config.Name {
		t.Errorf("Expected Name variable '%s', got '%v'", config.Name, preview.Variables["Name"])
	}

	// Test with non-existent template
	_, err = manager.PreviewTemplate("non-existent", config)
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_GetTemplateVariables(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing template
	variables, err := manager.GetTemplateVariables("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateVariables failed: %v", err)
	}

	if variables == nil {
		t.Fatal("Expected variables, got nil")
	}

	// Should have default variables
	if _, exists := variables["Name"]; !exists {
		t.Error("Expected 'Name' variable")
	}

	if _, exists := variables["Organization"]; !exists {
		t.Error("Expected 'Organization' variable")
	}

	// Test with non-existent template
	_, err = manager.GetTemplateVariables("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_CacheOperations(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test cache template
	err := manager.CacheTemplate("go-gin")
	if err != nil {
		t.Fatalf("CacheTemplate failed: %v", err)
	}

	// Test get cached templates
	templates, err := manager.GetCachedTemplates()
	if err != nil {
		t.Fatalf("GetCachedTemplates failed: %v", err)
	}

	if len(templates) == 0 {
		t.Error("Expected cached templates")
	}

	// Test clear cache
	err = manager.ClearTemplateCache()
	if err != nil {
		t.Fatalf("ClearTemplateCache failed: %v", err)
	}

	// Test refresh cache
	err = manager.RefreshTemplateCache()
	if err != nil {
		t.Fatalf("RefreshTemplateCache failed: %v", err)
	}
}

func TestManager_GetTemplateLocation(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test with existing template
	location, err := manager.GetTemplateLocation("go-gin")
	if err != nil {
		t.Fatalf("GetTemplateLocation failed: %v", err)
	}

	if location == "" {
		t.Error("Expected template location")
	}

	// Should indicate embedded template
	if !strings.Contains(location, "embedded:") {
		t.Error("Expected embedded template location")
	}

	// Test with non-existent template
	_, err = manager.GetTemplateLocation("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestManager_TemplateManagement(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	// Test install template (should return not implemented)
	err := manager.InstallTemplate("https://github.com/test/template", "test-template")
	if err == nil {
		t.Error("Expected error for unimplemented InstallTemplate")
	}

	// Test uninstall template (should return not implemented)
	err = manager.UninstallTemplate("test-template")
	if err == nil {
		t.Error("Expected error for unimplemented UninstallTemplate")
	}

	// Test update template (should return not implemented)
	err = manager.UpdateTemplate("test-template")
	if err == nil {
		t.Error("Expected error for unimplemented UpdateTemplate")
	}
}

func TestManager_ErrorHandling(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	mockEngine.SetError(true, "mock error")
	manager := NewManager(mockEngine)

	config := &models.ProjectConfig{
		Name:       "test-project",
		OutputPath: "/test/output",
	}

	// Test ProcessTemplate with engine error
	err := manager.ProcessTemplate("go-gin", config, "/test/output")
	if err == nil {
		t.Error("Expected error from mock engine")
	}

	// Test ProcessCustomTemplate with engine error
	err = manager.ProcessCustomTemplate("/test/template", config, "/test/output")
	if err == nil {
		t.Error("Expected error from mock engine")
	}
}

// Helper functions for testing

func TestManager_HelperFunctions(t *testing.T) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine).(*Manager)

	// Test formatDisplayName
	displayName := manager.formatDisplayName("test-template-name")
	expected := "Test Template Name"
	if displayName != expected {
		t.Errorf("Expected display name '%s', got '%s'", expected, displayName)
	}

	// Test inferTechnology
	technology := manager.inferTechnology("go-gin-api")
	if technology != "Go" {
		t.Errorf("Expected technology 'Go', got '%s'", technology)
	}

	technology = manager.inferTechnology("nextjs-app")
	if technology != "Next.js" {
		t.Errorf("Expected technology 'Next.js', got '%s'", technology)
	}

	// Test inferTags
	tags := manager.inferTags("go-gin", "backend")
	if len(tags) == 0 {
		t.Error("Expected tags to be inferred")
	}

	// Should include category
	found := false
	for _, tag := range tags {
		if tag == "backend" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected category to be included in tags")
	}

	// Test isValidTemplateName
	if !manager.isValidTemplateName("valid-template-name") {
		t.Error("Expected valid template name to be valid")
	}

	if manager.isValidTemplateName("Invalid Template Name") {
		t.Error("Expected invalid template name to be invalid")
	}

	// Test contains helper
	slice := []string{"a", "b", "c"}
	if !manager.contains(slice, "b") {
		t.Error("Expected contains to find 'b'")
	}

	if manager.contains(slice, "d") {
		t.Error("Expected contains not to find 'd'")
	}
}

// Benchmark tests
func BenchmarkManager_ListTemplates(b *testing.B) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.ListTemplates(interfaces.TemplateFilter{})
		if err != nil {
			b.Fatalf("ListTemplates failed: %v", err)
		}
	}
}

func BenchmarkManager_GetTemplateInfo(b *testing.B) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GetTemplateInfo("go-gin")
		if err != nil {
			b.Fatalf("GetTemplateInfo failed: %v", err)
		}
	}
}

func BenchmarkManager_SearchTemplates(b *testing.B) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.SearchTemplates("go")
		if err != nil {
			b.Fatalf("SearchTemplates failed: %v", err)
		}
	}
}

func BenchmarkManager_ValidateTemplate(b *testing.B) {
	mockEngine := NewMockTemplateEngine()
	manager := NewManager(mockEngine)
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.ValidateTemplate(tempDir)
		if err != nil {
			b.Fatalf("ValidateTemplate failed: %v", err)
		}
	}
}
