package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAutoFixManager(t *testing.T) {
	manager := NewAutoFixManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.fixStrategies)
	assert.False(t, manager.dryRun)
	assert.True(t, manager.backupEnabled)

	// Check that default strategies are loaded
	assert.Contains(t, manager.fixStrategies, "structure.readme.required")
	assert.Contains(t, manager.fixStrategies, "structure.license.required")
}

func TestAutoFixManager_SetDryRun(t *testing.T) {
	manager := NewAutoFixManager()

	// Test enabling dry run
	manager.SetDryRun(true)
	assert.True(t, manager.dryRun)

	// Test disabling dry run
	manager.SetDryRun(false)
	assert.False(t, manager.dryRun)
}

func TestAutoFixManager_SetBackupEnabled(t *testing.T) {
	manager := NewAutoFixManager()

	// Test disabling backup
	manager.SetBackupEnabled(false)
	assert.False(t, manager.backupEnabled)

	// Test enabling backup
	manager.SetBackupEnabled(true)
	assert.True(t, manager.backupEnabled)
}

func TestAutoFixManager_FixIssues(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	tests := []struct {
		name            string
		issues          []interfaces.ValidationIssue
		dryRun          bool
		expectedApplied int
		expectedSkipped int
	}{
		{
			name: "fix README issue in dry run",
			issues: []interfaces.ValidationIssue{
				{
					Type:    "error",
					Message: "README file is missing",
					File:    tempDir,
					Rule:    "structure.readme.required",
					Fixable: true,
				},
			},
			dryRun:          true,
			expectedApplied: 1,
			expectedSkipped: 0,
		},
		{
			name: "fix LICENSE issue in dry run",
			issues: []interfaces.ValidationIssue{
				{
					Type:    "error",
					Message: "LICENSE file is missing",
					File:    tempDir,
					Rule:    "structure.license.required",
					Fixable: true,
				},
			},
			dryRun:          true,
			expectedApplied: 1,
			expectedSkipped: 0,
		},
		{
			name: "skip non-fixable issue",
			issues: []interfaces.ValidationIssue{
				{
					Type:    "error",
					Message: "Some error",
					File:    tempDir,
					Rule:    "unknown.rule",
					Fixable: false,
				},
			},
			dryRun:          true,
			expectedApplied: 0,
			expectedSkipped: 0,
		},
		{
			name: "skip issue without strategy",
			issues: []interfaces.ValidationIssue{
				{
					Type:    "error",
					Message: "Some error",
					File:    tempDir,
					Rule:    "unknown.rule",
					Fixable: true,
				},
			},
			dryRun:          true,
			expectedApplied: 0,
			expectedSkipped: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager.SetDryRun(tt.dryRun)

			result, err := manager.FixIssues(tempDir, tt.issues)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedApplied, len(result.Applied))
			assert.Equal(t, tt.expectedSkipped, len(result.Skipped))
		})
	}
}

func TestAutoFixManager_PreviewFixes(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	issues := []interfaces.ValidationIssue{
		{
			Type:    "error",
			Message: "README file is missing",
			File:    tempDir,
			Rule:    "structure.readme.required",
			Fixable: true,
		},
		{
			Type:    "error",
			Message: "LICENSE file is missing",
			File:    tempDir,
			Rule:    "structure.license.required",
			Fixable: true,
		},
	}

	preview, err := manager.PreviewFixes(tempDir, issues)
	require.NoError(t, err)
	require.NotNil(t, preview)

	assert.Equal(t, 2, len(preview.Fixes))
	assert.Equal(t, 2, preview.Summary.AppliedFixes)
}

func TestAutoFixManager_GetFixableIssues(t *testing.T) {
	manager := NewAutoFixManager()

	issues := []interfaces.ValidationIssue{
		{
			Type:    "error",
			Message: "README file is missing",
			File:    "/test",
			Rule:    "structure.readme.required",
			Fixable: true,
		},
		{
			Type:    "error",
			Message: "LICENSE file is missing",
			File:    "/test",
			Rule:    "structure.license.required",
			Fixable: true,
		},
		{
			Type:    "error",
			Message: "Some non-fixable error",
			File:    "/test",
			Rule:    "unknown.rule",
			Fixable: false,
		},
		{
			Type:    "error",
			Message: "Fixable but no strategy",
			File:    "/test",
			Rule:    "no.strategy",
			Fixable: true,
		},
	}

	fixableIssues := manager.GetFixableIssues(issues)

	assert.Equal(t, 2, len(fixableIssues))
	assert.Equal(t, "structure.readme.required", fixableIssues[0].Rule)
	assert.Equal(t, "structure.license.required", fixableIssues[1].Rule)
}

func TestAutoFixManager_applyCreateFix(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	fix := interfaces.Fix{
		ID:          "test_create",
		Type:        "create_file",
		Description: "Create test file",
		File:        filepath.Join(tempDir, "test.txt"),
		Action:      interfaces.FixActionCreate,
		Content:     "Test content",
		Automatic:   true,
	}

	err := manager.applyCreateFix(fix)
	require.NoError(t, err)

	// Verify file was created
	content, err := os.ReadFile(fix.File)
	require.NoError(t, err)
	assert.Equal(t, "Test content", string(content))
}

func TestAutoFixManager_applyReplaceFix(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := "Line 1\nLine 2\nLine 3"
	err := os.WriteFile(testFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	fix := interfaces.Fix{
		ID:          "test_replace",
		Type:        "replace_line",
		Description: "Replace line 2",
		File:        testFile,
		Action:      interfaces.FixActionReplace,
		Content:     "New Line 2",
		Line:        2,
		Automatic:   true,
	}

	err = manager.applyReplaceFix(fix)
	require.NoError(t, err)

	// Verify file was modified
	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	expectedContent := "Line 1\nNew Line 2\nLine 3"
	assert.Equal(t, expectedContent, string(content))
}

func TestAutoFixManager_applyInsertFix(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := "Line 1\nLine 2"
	err := os.WriteFile(testFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	fix := interfaces.Fix{
		ID:          "test_insert",
		Type:        "insert_line",
		Description: "Insert line at position 2",
		File:        testFile,
		Action:      interfaces.FixActionInsert,
		Content:     "Inserted Line",
		Line:        2,
		Automatic:   true,
	}

	err = manager.applyInsertFix(fix)
	require.NoError(t, err)

	// Verify file was modified
	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	expectedContent := "Line 1\nInserted Line\nLine 2"
	assert.Equal(t, expectedContent, string(content))
}

func TestAutoFixManager_applyDeleteFix(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := "Line 1\nLine 2\nLine 3"
	err := os.WriteFile(testFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	fix := interfaces.Fix{
		ID:          "test_delete",
		Type:        "delete_line",
		Description: "Delete line 2",
		File:        testFile,
		Action:      interfaces.FixActionDelete,
		Line:        2,
		Automatic:   true,
	}

	err = manager.applyDeleteFix(fix)
	require.NoError(t, err)

	// Verify file was modified
	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	expectedContent := "Line 1\nLine 3"
	assert.Equal(t, expectedContent, string(content))
}

func TestAutoFixManager_applyRenameFix(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file
	oldFile := filepath.Join(tempDir, "old.txt")
	newFile := filepath.Join(tempDir, "new.txt")
	err := os.WriteFile(oldFile, []byte("test content"), 0644)
	require.NoError(t, err)

	fix := interfaces.Fix{
		ID:          "test_rename",
		Type:        "rename_file",
		Description: "Rename file",
		File:        oldFile,
		Action:      interfaces.FixActionRename,
		Content:     newFile,
		Automatic:   true,
	}

	err = manager.applyRenameFix(fix)
	require.NoError(t, err)

	// Verify old file doesn't exist and new file exists
	_, err = os.Stat(oldFile)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newFile)
	assert.NoError(t, err)
}

func TestAutoFixManager_createBackup(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := "original content"
	err := os.WriteFile(testFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	err = manager.createBackup(testFile)
	require.NoError(t, err)

	// Verify backup was created
	backupFile := testFile + ".backup"
	content, err := os.ReadFile(backupFile)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(content))
}

func TestAutoFixManager_createReadmeFix(t *testing.T) {
	manager := NewAutoFixManager()

	issue := interfaces.ValidationIssue{
		Type:    "error",
		Message: "README file is missing",
		File:    "/test/project",
		Rule:    "structure.readme.required",
		Fixable: true,
	}

	fix, err := manager.createReadmeFix(issue)
	require.NoError(t, err)
	require.NotNil(t, fix)

	assert.Equal(t, "create_file", fix.Type)
	assert.Equal(t, interfaces.FixActionCreate, fix.Action)
	assert.Contains(t, fix.File, "README.md")
	assert.Contains(t, fix.Content, "# Project Name")
	assert.True(t, fix.Automatic)
}

func TestAutoFixManager_createLicenseFix(t *testing.T) {
	manager := NewAutoFixManager()

	issue := interfaces.ValidationIssue{
		Type:    "error",
		Message: "LICENSE file is missing",
		File:    "/test/project",
		Rule:    "structure.license.required",
		Fixable: true,
	}

	fix, err := manager.createLicenseFix(issue)
	require.NoError(t, err)
	require.NotNil(t, fix)

	assert.Equal(t, "create_file", fix.Type)
	assert.Equal(t, interfaces.FixActionCreate, fix.Action)
	assert.Contains(t, fix.File, "LICENSE")
	assert.Contains(t, fix.Content, "MIT License")
	assert.True(t, fix.Automatic)
}

func TestAutoFixManager_fixNamingConventions(t *testing.T) {
	manager := NewAutoFixManager()

	tests := []struct {
		name        string
		issue       interfaces.ValidationIssue
		expectFix   bool
		expectedNew string
	}{
		{
			name: "file with spaces",
			issue: interfaces.ValidationIssue{
				Type:    "warning",
				Message: "File name contains spaces",
				File:    "/test/my file.txt",
				Rule:    "quality.naming.conventions",
				Fixable: true,
			},
			expectFix:   true,
			expectedNew: "/test/my_file.txt",
		},
		{
			name: "unsupported naming issue",
			issue: interfaces.ValidationIssue{
				Type:    "warning",
				Message: "Some other naming issue",
				File:    "/test/file.txt",
				Rule:    "quality.naming.conventions",
				Fixable: true,
			},
			expectFix: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fix, err := manager.fixNamingConventions(tt.issue)

			if tt.expectFix {
				require.NoError(t, err)
				require.NotNil(t, fix)
				assert.Equal(t, "rename_file", fix.Type)
				assert.Equal(t, interfaces.FixActionRename, fix.Action)
				assert.Equal(t, tt.expectedNew, fix.Content)
			} else {
				assert.Error(t, err)
				assert.Nil(t, fix)
			}
		})
	}
}

func TestAutoFixManager_createMissingFile(t *testing.T) {
	manager := NewAutoFixManager()

	tests := []struct {
		name            string
		filePath        string
		expectedContent string
	}{
		{
			name:            "markdown file",
			filePath:        "/test/README.md",
			expectedContent: "# README.md\n\nThis file was automatically generated.\n",
		},
		{
			name:            "text file",
			filePath:        "/test/notes.txt",
			expectedContent: "This file was automatically generated.\n",
		},
		{
			name:            "json file",
			filePath:        "/test/config.json",
			expectedContent: "{}\n",
		},
		{
			name:            "yaml file",
			filePath:        "/test/config.yaml",
			expectedContent: "# This file was automatically generated\n",
		},
		{
			name:            "unknown extension",
			filePath:        "/test/unknown.xyz",
			expectedContent: "# This file was automatically generated\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := interfaces.ValidationIssue{
				Type:    "error",
				Message: "Missing file",
				File:    tt.filePath,
				Rule:    "generic.create_missing_file",
				Fixable: true,
			}

			fix, err := manager.createMissingFile(issue)
			require.NoError(t, err)
			require.NotNil(t, fix)

			assert.Equal(t, "create_file", fix.Type)
			assert.Equal(t, interfaces.FixActionCreate, fix.Action)
			assert.Equal(t, tt.expectedContent, fix.Content)
			assert.False(t, fix.Automatic) // Generic file creation should require confirmation
		})
	}
}

func TestAutoFixManager_createFileChangePreview(t *testing.T) {
	manager := NewAutoFixManager()
	tempDir := t.TempDir()

	// Create test file for preview tests
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("Line 1\nLine 2\nLine 3"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name            string
		fix             interfaces.Fix
		expectedPreview string
	}{
		{
			name: "create fix",
			fix: interfaces.Fix{
				Action:  interfaces.FixActionCreate,
				Content: "New file\nwith content",
			},
			expectedPreview: "Create file with 2 lines",
		},
		{
			name: "replace fix",
			fix: interfaces.Fix{
				File:    testFile,
				Action:  interfaces.FixActionReplace,
				Content: "Replaced line",
				Line:    2,
			},
			expectedPreview: "Replace line 2: Replaced line",
		},
		{
			name: "insert fix",
			fix: interfaces.Fix{
				File:    testFile,
				Action:  interfaces.FixActionInsert,
				Content: "Inserted line",
				Line:    2,
			},
			expectedPreview: "Insert at line 2: Inserted line",
		},
		{
			name: "delete fix",
			fix: interfaces.Fix{
				File:   testFile,
				Action: interfaces.FixActionDelete,
				Line:   2,
			},
			expectedPreview: "Delete line 2",
		},
		{
			name: "rename fix",
			fix: interfaces.Fix{
				Action:  interfaces.FixActionRename,
				Content: "new_name.txt",
			},
			expectedPreview: "Rename to: new_name.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change, err := manager.createFileChangePreview(tt.fix)
			require.NoError(t, err)
			require.NotNil(t, change)

			assert.Equal(t, tt.expectedPreview, change.Preview)
			assert.Equal(t, tt.fix.Action, change.Action)
		})
	}
}

func TestAutoFixManager_findGenericStrategy(t *testing.T) {
	manager := NewAutoFixManager()

	tests := []struct {
		name        string
		issue       interfaces.ValidationIssue
		expectFound bool
	}{
		{
			name: "missing file issue",
			issue: interfaces.ValidationIssue{
				Message: "missing file detected",
			},
			expectFound: true,
		},
		{
			name: "naming issue with spaces",
			issue: interfaces.ValidationIssue{
				Message: "file name contains space characters",
			},
			expectFound: true,
		},
		{
			name: "unrecognized issue",
			issue: interfaces.ValidationIssue{
				Message: "some other issue",
			},
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, found := manager.findGenericStrategy(tt.issue)

			assert.Equal(t, tt.expectFound, found)
			if tt.expectFound {
				// The strategy name is the descriptive name, not the rule ID
				assert.NotEmpty(t, strategy.Name)
			}
		})
	}
}
