package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestDisabledTemplateFiltering(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test template structure with disabled files
	templateDir := filepath.Join(tempDir, "templates")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory structure
	if err := os.MkdirAll(filepath.Join(templateDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create test files
	testFiles := []string{
		"normal.txt.tmpl",
		"disabled.txt.tmpl.disabled",
		"subdir/normal.txt.tmpl",
		"subdir/disabled.txt.tmpl.disabled",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(templateDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test with regular engine
	engine := NewEngine()
	config := &models.ProjectConfig{
		Name:         "test",
		Organization: "test-org",
	}

	// Process templates
	err := engine.ProcessDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("Failed to process templates: %v", err)
	}

	// Check that disabled files were not generated
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl.disabled") {
			t.Errorf("Disabled template file was generated: %s", path)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk output directory: %v", err)
	}

	// Check that normal files were generated (without .tmpl extension)
	expectedFiles := []string{
		"normal.txt",
		"subdir/normal.txt",
	}

	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(outputDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected file was not generated: %s", expectedFile)
		}
	}
}
