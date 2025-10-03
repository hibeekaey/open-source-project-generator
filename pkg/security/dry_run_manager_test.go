package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDryRunManager_SetEnabled(t *testing.T) {
	drm := NewDryRunManager()

	// Test default state
	assert.False(t, drm.IsEnabled())

	// Test enabling
	drm.SetEnabled(true)
	assert.True(t, drm.IsEnabled())
	assert.Empty(t, drm.GetOperations()) // Should reset operations when enabling

	// Add some operations
	drm.RecordFileWrite("/test/file.txt", []byte("content"), false)
	assert.Len(t, drm.GetOperations(), 1)

	// Enable again - should reset operations
	drm.SetEnabled(true)
	assert.Empty(t, drm.GetOperations())

	// Test disabling
	drm.SetEnabled(false)
	assert.False(t, drm.IsEnabled())
}

func TestDryRunManager_RecordFileWrite(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	tests := []struct {
		name           string
		setupFile      bool
		overwrite      bool
		expectedImpact string
		expectedDesc   string
	}{
		{
			name:           "create new file",
			setupFile:      false,
			overwrite:      false,
			expectedImpact: "safe",
			expectedDesc:   "Create file:",
		},
		{
			name:           "overwrite existing file",
			setupFile:      true,
			overwrite:      true,
			expectedImpact: "destructive",
			expectedDesc:   "Overwrite file:",
		},
		{
			name:           "skip existing file",
			setupFile:      true,
			overwrite:      false,
			expectedImpact: "warning",
			expectedDesc:   "File exists, would skip:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupFile {
				err := os.WriteFile(testFile, []byte("existing"), 0644)
				require.NoError(t, err)
			} else {
				_ = os.Remove(testFile) // Ensure file doesn't exist
			}

			drm.Clear()

			// Record operation
			data := []byte("test content")
			drm.RecordFileWrite(testFile, data, tt.overwrite)

			// Verify
			ops := drm.GetOperations()
			require.Len(t, ops, 1)

			op := ops[0]
			assert.Equal(t, "file_write", op.Type)
			assert.Contains(t, op.Description, tt.expectedDesc)
			assert.Equal(t, tt.expectedImpact, op.Impact)
			assert.Equal(t, testFile, op.Path)
			assert.Equal(t, int64(len(data)), op.Size)
			assert.Equal(t, len(data), op.Details["size"])
			assert.Equal(t, tt.overwrite, op.Details["overwrite"])
		})
	}
}

func TestDryRunManager_RecordFileDelete(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.txt")
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")

	// Create existing file
	content := []byte("test content")
	err := os.WriteFile(existingFile, content, 0644)
	require.NoError(t, err)

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Test deleting existing file
	drm.RecordFileDelete(existingFile)
	ops := drm.GetOperations()
	require.Len(t, ops, 1)

	op := ops[0]
	assert.Equal(t, "file_delete", op.Type)
	assert.Contains(t, op.Description, "Delete file:")
	assert.Equal(t, "destructive", op.Impact)
	assert.Equal(t, existingFile, op.Path)
	assert.Equal(t, int64(len(content)), op.Size)
	assert.True(t, op.Details["exists"].(bool))

	// Test deleting non-existent file
	drm.Clear()
	drm.RecordFileDelete(nonExistentFile)
	ops = drm.GetOperations()
	require.Len(t, ops, 1)

	op = ops[0]
	assert.Equal(t, "file_delete", op.Type)
	assert.Contains(t, op.Description, "File does not exist")
	assert.Equal(t, "safe", op.Impact)
	assert.Equal(t, nonExistentFile, op.Path)
	assert.Equal(t, int64(0), op.Size)
	assert.False(t, op.Details["exists"].(bool))
}

func TestDryRunManager_RecordDirectoryCreate(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := filepath.Join(tempDir, "existing")
	newDir := filepath.Join(tempDir, "new")

	// Create existing directory
	err := os.MkdirAll(existingDir, 0755)
	require.NoError(t, err)

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Test creating new directory
	drm.RecordDirectoryCreate(newDir)
	ops := drm.GetOperations()
	require.Len(t, ops, 1)

	op := ops[0]
	assert.Equal(t, "directory_create", op.Type)
	assert.Contains(t, op.Description, "Create directory:")
	assert.Equal(t, "safe", op.Impact)
	assert.Equal(t, newDir, op.Path)
	assert.False(t, op.Details["exists"].(bool))

	// Test creating existing directory
	drm.Clear()
	drm.RecordDirectoryCreate(existingDir)
	ops = drm.GetOperations()
	require.Len(t, ops, 1)

	op = ops[0]
	assert.Equal(t, "directory_create", op.Type)
	assert.Contains(t, op.Description, "Directory exists, would skip:")
	assert.Equal(t, "warning", op.Impact)
	assert.Equal(t, existingDir, op.Path)
	assert.True(t, op.Details["exists"].(bool))
}

func TestDryRunManager_RecordDirectoryDelete(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := filepath.Join(tempDir, "existing")
	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	// Create existing directory with files
	err := os.MkdirAll(existingDir, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(existingDir, "file1.txt")
	file2 := filepath.Join(existingDir, "subdir", "file2.txt")

	err = os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Dir(file2), 0755)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Test deleting existing directory
	drm.RecordDirectoryDelete(existingDir)
	ops := drm.GetOperations()
	require.Len(t, ops, 1)

	op := ops[0]
	assert.Equal(t, "directory_delete", op.Type)
	assert.Contains(t, op.Description, "Delete directory:")
	assert.Contains(t, op.Description, "2 files") // Should count files
	assert.Equal(t, "destructive", op.Impact)
	assert.Equal(t, existingDir, op.Path)
	assert.True(t, op.Details["exists"].(bool))
	assert.Equal(t, 2, op.Details["file_count"])
	assert.Greater(t, op.Details["total_size"].(int64), int64(0))

	// Test deleting non-existent directory
	drm.Clear()
	drm.RecordDirectoryDelete(nonExistentDir)
	ops = drm.GetOperations()
	require.Len(t, ops, 1)

	op = ops[0]
	assert.Equal(t, "directory_delete", op.Type)
	assert.Contains(t, op.Description, "Directory does not exist")
	assert.Equal(t, "safe", op.Impact)
	assert.Equal(t, nonExistentDir, op.Path)
	assert.False(t, op.Details["exists"].(bool))
	assert.Equal(t, 0, op.Details["file_count"])
}

func TestDryRunManager_RecordFileCopy(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")
	nonExistentSrc := filepath.Join(tempDir, "nonexistent.txt")

	// Create source file
	content := []byte("test content")
	err := os.WriteFile(srcFile, content, 0644)
	require.NoError(t, err)

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	tests := []struct {
		name           string
		srcPath        string
		dstPath        string
		setupDst       bool
		overwrite      bool
		expectedImpact string
		expectedDesc   string
	}{
		{
			name:           "copy to new destination",
			srcPath:        srcFile,
			dstPath:        dstFile,
			setupDst:       false,
			overwrite:      false,
			expectedImpact: "safe",
			expectedDesc:   "Copy file:",
		},
		{
			name:           "copy with overwrite",
			srcPath:        srcFile,
			dstPath:        dstFile,
			setupDst:       true,
			overwrite:      true,
			expectedImpact: "destructive",
			expectedDesc:   "Copy and overwrite:",
		},
		{
			name:           "copy without overwrite (existing dst)",
			srcPath:        srcFile,
			dstPath:        dstFile,
			setupDst:       true,
			overwrite:      false,
			expectedImpact: "warning",
			expectedDesc:   "Destination exists, would skip:",
		},
		{
			name:           "copy non-existent source",
			srcPath:        nonExistentSrc,
			dstPath:        dstFile,
			setupDst:       false,
			overwrite:      false,
			expectedImpact: "warning",
			expectedDesc:   "Source file does not exist:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup destination if needed
			if tt.setupDst {
				err := os.WriteFile(tt.dstPath, []byte("existing"), 0644)
				require.NoError(t, err)
			} else {
				_ = os.Remove(tt.dstPath)
			}

			drm.Clear()

			// Record operation
			drm.RecordFileCopy(tt.srcPath, tt.dstPath, tt.overwrite)

			// Verify
			ops := drm.GetOperations()
			require.Len(t, ops, 1)

			op := ops[0]
			assert.Equal(t, "file_copy", op.Type)
			assert.Contains(t, op.Description, tt.expectedDesc)
			assert.Equal(t, tt.expectedImpact, op.Impact)
			assert.Contains(t, op.Path, tt.srcPath)
			assert.Contains(t, op.Path, tt.dstPath)
			assert.Equal(t, tt.srcPath, op.Details["source"])
			assert.Equal(t, tt.dstPath, op.Details["destination"])
			assert.Equal(t, tt.overwrite, op.Details["overwrite"])
		})
	}
}

func TestDryRunManager_RecordTemplateProcess(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "template.tmpl")
	outputPath := filepath.Join(tempDir, "output.txt")

	drm := NewDryRunManager()
	drm.SetEnabled(true)

	variables := map[string]interface{}{
		"name":    "test",
		"version": "1.0.0",
	}

	// Test processing to new output
	drm.RecordTemplateProcess(templatePath, outputPath, variables)
	ops := drm.GetOperations()
	require.Len(t, ops, 1)

	op := ops[0]
	assert.Equal(t, "template_process", op.Type)
	assert.Contains(t, op.Description, "Process template:")
	assert.Equal(t, "safe", op.Impact)
	assert.Equal(t, outputPath, op.Path)
	assert.Equal(t, templatePath, op.Details["template"])
	assert.Equal(t, outputPath, op.Details["output"])
	assert.Equal(t, variables, op.Details["variables"])
	assert.False(t, op.Details["output_exists"].(bool))

	// Test processing to existing output
	err := os.WriteFile(outputPath, []byte("existing"), 0644)
	require.NoError(t, err)

	drm.Clear()
	drm.RecordTemplateProcess(templatePath, outputPath, variables)
	ops = drm.GetOperations()
	require.Len(t, ops, 1)

	op = ops[0]
	assert.Equal(t, "template_process", op.Type)
	assert.Contains(t, op.Description, "Process template and overwrite:")
	assert.Equal(t, "destructive", op.Impact)
	assert.True(t, op.Details["output_exists"].(bool))
}

func TestDryRunManager_RecordCustomOperation(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	details := map[string]interface{}{
		"custom_field": "custom_value",
		"count":        42,
	}

	// Test valid impact
	drm.RecordCustomOperation("custom_op", "Custom operation description", "/test/path", details, "warning")
	ops := drm.GetOperations()
	require.Len(t, ops, 1)

	op := ops[0]
	assert.Equal(t, "custom_op", op.Type)
	assert.Equal(t, "Custom operation description", op.Description)
	assert.Equal(t, "/test/path", op.Path)
	assert.Equal(t, "warning", op.Impact)
	assert.Equal(t, details, op.Details)

	// Test invalid impact (should default to warning)
	drm.Clear()
	drm.RecordCustomOperation("custom_op", "Description", "/path", details, "invalid")
	ops = drm.GetOperations()
	require.Len(t, ops, 1)
	assert.Equal(t, "warning", ops[0].Impact)
}

func TestDryRunManager_GetOperationsByImpact(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Add operations with different impacts
	drm.RecordCustomOperation("safe_op", "Safe operation", "/path1", nil, "safe")
	drm.RecordCustomOperation("warning_op", "Warning operation", "/path2", nil, "warning")
	drm.RecordCustomOperation("destructive_op", "Destructive operation", "/path3", nil, "destructive")
	drm.RecordCustomOperation("another_safe_op", "Another safe operation", "/path4", nil, "safe")

	// Test filtering by impact
	safeOps := drm.GetOperationsByImpact("safe")
	assert.Len(t, safeOps, 2)

	warningOps := drm.GetOperationsByImpact("warning")
	assert.Len(t, warningOps, 1)

	destructiveOps := drm.GetOperationsByImpact("destructive")
	assert.Len(t, destructiveOps, 1)

	nonExistentOps := drm.GetOperationsByImpact("nonexistent")
	assert.Len(t, nonExistentOps, 0)
}

func TestDryRunManager_GetSummary(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Create a temporary file to test overwrite scenario
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "file2.txt")
	err := os.WriteFile(existingFile, []byte("existing content"), 0644)
	require.NoError(t, err)

	// Add various operations
	drm.RecordFileWrite("/file1.txt", []byte("content1"), false)                // safe
	drm.RecordFileWrite(existingFile, []byte("content2"), true)                 // destructive (overwrite)
	drm.RecordDirectoryCreate("/newdir")                                        // safe
	drm.RecordCustomOperation("warning_op", "Warning", "/path", nil, "warning") // warning

	summary := drm.GetSummary()

	assert.Equal(t, 4, summary["total_operations"])
	assert.Equal(t, 2, summary["safe_operations"])
	assert.Equal(t, 1, summary["warning_operations"])
	assert.Equal(t, 1, summary["destructive_operations"])
	assert.Equal(t, 4, summary["total_files_affected"]) // 4 unique paths
	assert.Greater(t, summary["total_size_affected"].(int64), int64(0))

	opsByType := summary["operations_by_type"].(map[string]int)
	assert.Equal(t, 2, opsByType["file_write"])
	assert.Equal(t, 1, opsByType["directory_create"])
	assert.Equal(t, 1, opsByType["warning_op"])
}

func TestDryRunManager_GenerateReport(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Add some operations
	drm.RecordFileWrite("/file.txt", []byte("content"), false)
	drm.RecordCustomOperation("warning_op", "Warning operation", "/path", nil, "warning")
	drm.RecordCustomOperation("destructive_op", "Destructive operation", "/path2", nil, "destructive")

	// Test text report
	textReport, err := drm.GenerateReport("text")
	assert.NoError(t, err)
	assert.Contains(t, textReport, "DRY RUN REPORT")
	assert.Contains(t, textReport, "Total Operations: 3")
	assert.Contains(t, textReport, "DESTRUCTIVE OPERATIONS")
	assert.Contains(t, textReport, "WARNING OPERATIONS")
	assert.Contains(t, textReport, "SAFE OPERATIONS")

	// Test JSON report
	jsonReport, err := drm.GenerateReport("json")
	assert.NoError(t, err)
	assert.Contains(t, jsonReport, "summary")
	assert.Contains(t, jsonReport, "operations")
	assert.Contains(t, jsonReport, "timestamp")

	// Test unsupported format
	_, err = drm.GenerateReport("unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported report format")
}

func TestDryRunManager_Clear(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Add operations
	drm.RecordFileWrite("/file.txt", []byte("content"), false)
	drm.RecordDirectoryCreate("/dir")

	assert.Len(t, drm.GetOperations(), 2)

	// Clear operations
	drm.Clear()
	assert.Len(t, drm.GetOperations(), 0)
}

func TestDryRunManager_HasDestructiveOperations(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Initially no destructive operations
	assert.False(t, drm.HasDestructiveOperations())

	// Add safe operation
	drm.RecordFileWrite("/file.txt", []byte("content"), false)
	assert.False(t, drm.HasDestructiveOperations())

	// Add destructive operation
	drm.RecordCustomOperation("destructive_op", "Destructive", "/path", nil, "destructive")
	assert.True(t, drm.HasDestructiveOperations())
}

func TestDryRunManager_HasWarningOperations(t *testing.T) {
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Initially no warning operations
	assert.False(t, drm.HasWarningOperations())

	// Add safe operation
	drm.RecordFileWrite("/file.txt", []byte("content"), false)
	assert.False(t, drm.HasWarningOperations())

	// Add warning operation
	drm.RecordCustomOperation("warning_op", "Warning", "/path", nil, "warning")
	assert.True(t, drm.HasWarningOperations())
}

func TestDryRunManager_DisabledMode(t *testing.T) {
	drm := NewDryRunManager()
	// Keep disabled (default state)

	// Try to record operations while disabled
	drm.RecordFileWrite("/file.txt", []byte("content"), false)
	drm.RecordDirectoryCreate("/dir")
	drm.RecordFileCopy("/src", "/dst", false)

	// Should not record any operations
	assert.Empty(t, drm.GetOperations())
	assert.False(t, drm.HasDestructiveOperations())
	assert.False(t, drm.HasWarningOperations())

	summary := drm.GetSummary()
	assert.Equal(t, 0, summary["total_operations"])
}

func TestDryRunManager_Integration(t *testing.T) {
	tempDir := t.TempDir()
	drm := NewDryRunManager()
	drm.SetEnabled(true)

	// Create necessary source files for the test
	templatesDir := filepath.Join(tempDir, "templates")
	staticDir := filepath.Join(tempDir, "static")
	err := os.MkdirAll(templatesDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(staticDir, 0755)
	require.NoError(t, err)

	// Create template files
	err = os.WriteFile(filepath.Join(templatesDir, "main.go.tmpl"), []byte("package main"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(templatesDir, "README.md.tmpl"), []byte("# {{.project_name}}"), 0644)
	require.NoError(t, err)

	// Create static file
	err = os.WriteFile(filepath.Join(staticDir, "Makefile"), []byte("all:\n\techo 'build'"), 0644)
	require.NoError(t, err)

	// Simulate a complete project generation workflow

	// 1. Create project directory
	projectDir := filepath.Join(tempDir, "my-project")
	drm.RecordDirectoryCreate(projectDir)

	// 2. Process templates
	templateVars := map[string]interface{}{
		"project_name": "my-project",
		"author":       "test-author",
	}

	drm.RecordTemplateProcess(
		filepath.Join(tempDir, "templates", "main.go.tmpl"),
		filepath.Join(projectDir, "main.go"),
		templateVars,
	)

	drm.RecordTemplateProcess(
		filepath.Join(tempDir, "templates", "README.md.tmpl"),
		filepath.Join(projectDir, "README.md"),
		templateVars,
	)

	// 3. Copy static files
	drm.RecordFileCopy(
		filepath.Join(tempDir, "static", "Makefile"),
		filepath.Join(projectDir, "Makefile"),
		false,
	)

	// 4. Create additional directories
	drm.RecordDirectoryCreate(filepath.Join(projectDir, "cmd"))
	drm.RecordDirectoryCreate(filepath.Join(projectDir, "pkg"))

	// Verify the complete workflow
	ops := drm.GetOperations()
	assert.Len(t, ops, 6)

	summary := drm.GetSummary()
	assert.Equal(t, 6, summary["total_operations"])
	assert.Equal(t, 6, summary["safe_operations"]) // All should be safe operations
	assert.Equal(t, 0, summary["warning_operations"])
	assert.Equal(t, 0, summary["destructive_operations"])

	// Generate and verify report
	report, err := drm.GenerateReport("text")
	assert.NoError(t, err)
	assert.Contains(t, report, "Total Operations: 6")
	assert.Contains(t, report, "Safe Operations: 6")
	assert.Contains(t, report, "SAFE OPERATIONS")

	// Verify no destructive or warning operations
	assert.False(t, drm.HasDestructiveOperations())
	assert.False(t, drm.HasWarningOperations())
}
