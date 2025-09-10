package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Fatal("NewGenerator() returned nil")
	}

	// Verify it's the correct type
	if _, ok := gen.(*Generator); !ok {
		t.Fatal("NewGenerator() did not return a *Generator")
	}
}

func TestNewDryRunGenerator(t *testing.T) {
	gen := NewDryRunGenerator()
	if gen == nil {
		t.Fatal("NewDryRunGenerator() returned nil")
	}

	// Verify it's the correct type and in dry-run mode
	if g, ok := gen.(*Generator); !ok {
		t.Fatal("NewDryRunGenerator() did not return a *Generator")
	} else if !g.dryRun {
		t.Fatal("NewDryRunGenerator() did not set dryRun to true")
	}
}

func TestCreateDirectory(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid directory path",
			path:        filepath.Join(tempDir, "test-dir"),
			expectError: false,
		},
		{
			name:        "nested directory path",
			path:        filepath.Join(tempDir, "nested", "deep", "directory"),
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "directory path cannot be empty",
		},
		{
			name:        "path traversal attempt",
			path:        tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.CreateDirectory(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("CreateDirectory() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CreateDirectory() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CreateDirectory() unexpected error = %v", err)
				}
				// Verify directory was created
				if !gen.FileExists(tt.path) {
					t.Errorf("CreateDirectory() did not create directory %s", tt.path)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		content     []byte
		perm        os.FileMode
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid file write",
			path:        filepath.Join(tempDir, "test.txt"),
			content:     []byte("test content"),
			perm:        0644,
			expectError: false,
		},
		{
			name:        "file in nested directory",
			path:        filepath.Join(tempDir, "nested", "test.txt"),
			content:     []byte("nested content"),
			perm:        0644,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			content:     []byte("content"),
			perm:        0644,
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "nil content",
			path:        filepath.Join(tempDir, "nil-content.txt"),
			content:     nil,
			perm:        0644,
			expectError: true,
			errorMsg:    "file content cannot be nil",
		},
		{
			name:        "path traversal attempt",
			path:        tempDir + "/../malicious.txt",
			content:     []byte("malicious"),
			perm:        0644,
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "invalid permissions",
			path:        filepath.Join(tempDir, "invalid-perm.txt"),
			content:     []byte("content"),
			perm:        01000,
			expectError: true,
			errorMsg:    "invalid file permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.WriteFile(tt.path, tt.content, tt.perm)

			if tt.expectError {
				if err == nil {
					t.Errorf("WriteFile() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("WriteFile() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("WriteFile() unexpected error = %v", err)
				}
				// Verify file was created with correct content
				if !gen.FileExists(tt.path) {
					t.Errorf("WriteFile() did not create file %s", tt.path)
				} else {
					content, err := os.ReadFile(tt.path)
					if err != nil {
						t.Errorf("Failed to read created file: %v", err)
					} else if string(content) != string(tt.content) {
						t.Errorf("WriteFile() content = %s, expected %s", content, tt.content)
					}
				}
			}
		})
	}
}

func TestCopyAssets(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	// Create source directory with test files
	srcDir := filepath.Join(tempDir, "source")
	if err := gen.CreateDirectory(srcDir); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create test files in source
	testFiles := map[string][]byte{
		"file1.txt":          []byte("content1"),
		"subdir/file2.txt":   []byte("content2"),
		"subdir/file3.bin":   []byte{0x00, 0x01, 0x02, 0x03},
		"empty-dir/.gitkeep": []byte(""),
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(srcDir, filePath)
		if err := gen.WriteFile(fullPath, content, 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	tests := []struct {
		name        string
		srcDir      string
		destDir     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid asset copy",
			srcDir:      srcDir,
			destDir:     filepath.Join(tempDir, "dest1"),
			expectError: false,
		},
		{
			name:        "empty source directory",
			srcDir:      "",
			destDir:     filepath.Join(tempDir, "dest2"),
			expectError: true,
			errorMsg:    "source directory cannot be empty",
		},
		{
			name:        "empty destination directory",
			srcDir:      srcDir,
			destDir:     "",
			expectError: true,
			errorMsg:    "destination directory cannot be empty",
		},
		{
			name:        "non-existent source directory",
			srcDir:      filepath.Join(tempDir, "non-existent"),
			destDir:     filepath.Join(tempDir, "dest3"),
			expectError: true,
			errorMsg:    "source directory does not exist",
		},
		{
			name:        "path traversal in source",
			srcDir:      tempDir + "/../malicious",
			destDir:     filepath.Join(tempDir, "dest4"),
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "path traversal in destination",
			srcDir:      srcDir,
			destDir:     tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.CopyAssets(tt.srcDir, tt.destDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("CopyAssets() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CopyAssets() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CopyAssets() unexpected error = %v", err)
				}
				// Verify all files were copied
				for filePath, expectedContent := range testFiles {
					destFilePath := filepath.Join(tt.destDir, filePath)
					if !gen.FileExists(destFilePath) {
						t.Errorf("CopyAssets() did not copy file %s", filePath)
					} else {
						content, err := os.ReadFile(destFilePath)
						if err != nil {
							t.Errorf("Failed to read copied file %s: %v", filePath, err)
						} else if string(content) != string(expectedContent) {
							t.Errorf("CopyAssets() file %s content = %s, expected %s", filePath, content, expectedContent)
						}
					}
				}
			}
		})
	}
}

func TestCreateSymlink(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	// Create a target file
	targetFile := filepath.Join(tempDir, "target.txt")
	if err := gen.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	tests := []struct {
		name        string
		target      string
		link        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid symlink creation",
			target:      targetFile,
			link:        filepath.Join(tempDir, "link.txt"),
			expectError: false,
		},
		{
			name:        "symlink in nested directory",
			target:      targetFile,
			link:        filepath.Join(tempDir, "nested", "link.txt"),
			expectError: false,
		},
		{
			name:        "empty target",
			target:      "",
			link:        filepath.Join(tempDir, "empty-target.txt"),
			expectError: true,
			errorMsg:    "symlink target cannot be empty",
		},
		{
			name:        "empty link path",
			target:      targetFile,
			link:        "",
			expectError: true,
			errorMsg:    "symlink path cannot be empty",
		},
		{
			name:        "path traversal in link",
			target:      targetFile,
			link:        tempDir + "/../malicious-link.txt",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.CreateSymlink(tt.target, tt.link)

			if tt.expectError {
				if err == nil {
					t.Errorf("CreateSymlink() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CreateSymlink() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CreateSymlink() unexpected error = %v", err)
				}
				// Verify symlink was created
				if !gen.FileExists(tt.link) {
					t.Errorf("CreateSymlink() did not create symlink %s", tt.link)
				} else {
					// Verify it's actually a symlink
					linkInfo, err := os.Lstat(tt.link)
					if err != nil {
						t.Errorf("Failed to stat symlink: %v", err)
					} else if linkInfo.Mode()&os.ModeSymlink == 0 {
						t.Errorf("CreateSymlink() created file is not a symlink")
					}
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "exists.txt")
	if err := gen.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test directory
	testDir := filepath.Join(tempDir, "exists-dir")
	if err := gen.CreateDirectory(testDir); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     testFile,
			expected: true,
		},
		{
			name:     "existing directory",
			path:     testDir,
			expected: true,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tempDir, "non-existent.txt"),
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.FileExists(tt.path)
			if result != tt.expected {
				t.Errorf("FileExists() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEnsureDirectory(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	// Create an existing directory
	existingDir := filepath.Join(tempDir, "existing")
	if err := gen.CreateDirectory(existingDir); err != nil {
		t.Fatalf("Failed to create existing directory: %v", err)
	}

	// Create a file that conflicts with directory creation
	conflictFile := filepath.Join(tempDir, "conflict.txt")
	if err := gen.WriteFile(conflictFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create conflict file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "new directory",
			path:        filepath.Join(tempDir, "new-dir"),
			expectError: false,
		},
		{
			name:        "existing directory",
			path:        existingDir,
			expectError: false,
		},
		{
			name:        "nested directory",
			path:        filepath.Join(tempDir, "nested", "deep", "dir"),
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "directory path cannot be empty",
		},
		{
			name:        "path traversal",
			path:        tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "conflict with existing file",
			path:        conflictFile,
			expectError: true,
			errorMsg:    "exists but is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.EnsureDirectory(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("EnsureDirectory() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("EnsureDirectory() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("EnsureDirectory() unexpected error = %v", err)
				}
				// Verify directory exists
				if !gen.FileExists(tt.path) {
					t.Errorf("EnsureDirectory() did not create directory %s", tt.path)
				}
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	// Create a valid project config
	validConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project description",
		License:      "MIT",
		OutputPath:   tempDir,
		GeneratedAt:  time.Now(),
	}

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		outputPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project creation",
			config:      validConfig,
			outputPath:  tempDir,
			expectError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			outputPath:  tempDir,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
		{
			name:        "empty output path",
			config:      validConfig,
			outputPath:  "",
			expectError: true,
			errorMsg:    "output path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.CreateProject(tt.config, tt.outputPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("CreateProject() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CreateProject() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CreateProject() unexpected error = %v", err)
				}
				// Verify project directory was created
				projectPath := filepath.Join(tt.outputPath, tt.config.Name)
				if !gen.FileExists(projectPath) {
					t.Errorf("CreateProject() did not create project directory %s", projectPath)
				}
			}
		})
	}
}

func TestDryRunMode(t *testing.T) {
	gen := NewDryRunGenerator()
	tempDir := t.TempDir()

	// Test that dry-run mode doesn't actually create files/directories
	testDir := filepath.Join(tempDir, "dry-run-test")
	testFile := filepath.Join(tempDir, "dry-run-file.txt")

	// These operations should not fail in dry-run mode
	if err := gen.CreateDirectory(testDir); err != nil {
		t.Errorf("CreateDirectory() in dry-run mode failed: %v", err)
	}

	if err := gen.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Errorf("WriteFile() in dry-run mode failed: %v", err)
	}

	if err := gen.EnsureDirectory(testDir); err != nil {
		t.Errorf("EnsureDirectory() in dry-run mode failed: %v", err)
	}

	// Verify that nothing was actually created
	if gen.FileExists(testDir) {
		t.Errorf("CreateDirectory() in dry-run mode actually created directory")
	}

	if gen.FileExists(testFile) {
		t.Errorf("WriteFile() in dry-run mode actually created file")
	}
}

func TestFileSystemGeneratorErrorHandling(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	t.Run("disk space simulation", func(t *testing.T) {
		// Create a very large file to test disk space handling
		largeContent := make([]byte, 1024*1024) // 1MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}

		largePath := filepath.Join(tempDir, "large-file.bin")
		err := gen.WriteFile(largePath, largeContent, 0644)
		if err != nil {
			t.Logf("Large file write result: %v", err)
			// This might fail on systems with limited space, which is expected
		}
	})

	t.Run("permission denied scenarios", func(t *testing.T) {
		// Create a directory with restricted permissions
		restrictedDir := filepath.Join(tempDir, "restricted")
		if err := gen.CreateDirectory(restrictedDir); err != nil {
			t.Fatalf("Failed to create restricted directory: %v", err)
		}

		// Make it read-only
		if err := os.Chmod(restrictedDir, 0444); err != nil {
			t.Fatalf("Failed to change permissions: %v", err)
		}
		defer os.Chmod(restrictedDir, 0755) // Restore for cleanup

		// Try to create a file in the restricted directory
		restrictedFile := filepath.Join(restrictedDir, "test.txt")
		err := gen.WriteFile(restrictedFile, []byte("content"), 0644)
		if err == nil {
			t.Error("Expected permission error when writing to restricted directory")
		}
	})

	t.Run("concurrent file operations", func(t *testing.T) {
		const numGoroutines = 50
		results := make(chan error, numGoroutines)

		// Concurrent file writes
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				filePath := filepath.Join(tempDir, fmt.Sprintf("concurrent-%d.txt", index))
				content := fmt.Sprintf("Content for file %d", index)
				results <- gen.WriteFile(filePath, []byte(content), 0644)
			}(i)
		}

		// Check results
		for i := 0; i < numGoroutines; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent file write %d failed: %v", i, err)
			}
		}

		// Verify all files were created
		for i := 0; i < numGoroutines; i++ {
			filePath := filepath.Join(tempDir, fmt.Sprintf("concurrent-%d.txt", i))
			if !gen.FileExists(filePath) {
				t.Errorf("Concurrent file %d was not created", i)
			}
		}
	})

	t.Run("symlink edge cases", func(t *testing.T) {
		// Test symlink to non-existent target
		nonExistentTarget := filepath.Join(tempDir, "non-existent-target.txt")
		symlinkPath := filepath.Join(tempDir, "broken-symlink.txt")

		err := gen.CreateSymlink(nonExistentTarget, symlinkPath)
		if err != nil {
			t.Logf("Broken symlink creation result: %v", err)
			// This might be allowed on some systems
		}

		// Test circular symlinks
		symlink1 := filepath.Join(tempDir, "symlink1.txt")
		symlink2 := filepath.Join(tempDir, "symlink2.txt")

		gen.CreateSymlink(symlink2, symlink1)
		err = gen.CreateSymlink(symlink1, symlink2)
		if err != nil {
			t.Logf("Circular symlink creation result: %v", err)
			// This might be detected by the OS
		}
	})

	t.Run("asset copying with special files", func(t *testing.T) {
		srcDir := filepath.Join(tempDir, "special-src")
		destDir := filepath.Join(tempDir, "special-dest")

		if err := gen.CreateDirectory(srcDir); err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Create files with special names
		specialFiles := []string{
			"file with spaces.txt",
			"file-with-unicode-测试.txt",
			"file.with.dots.txt",
			".hidden-file",
			"UPPERCASE.TXT",
		}

		for _, fileName := range specialFiles {
			filePath := filepath.Join(srcDir, fileName)
			content := fmt.Sprintf("Content of %s", fileName)
			if err := gen.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Errorf("Failed to create special file %s: %v", fileName, err)
				continue
			}
		}

		// Copy assets
		err := gen.CopyAssets(srcDir, destDir)
		if err != nil {
			t.Errorf("CopyAssets with special files failed: %v", err)
		}

		// Verify special files were copied
		for _, fileName := range specialFiles {
			destPath := filepath.Join(destDir, fileName)
			if !gen.FileExists(destPath) {
				t.Errorf("Special file %s was not copied", fileName)
			}
		}
	})

	t.Run("very deep directory structure", func(t *testing.T) {
		// Create a very deep directory structure
		deepPath := tempDir
		for i := 0; i < 100; i++ {
			deepPath = filepath.Join(deepPath, fmt.Sprintf("level-%d", i))
		}

		err := gen.CreateDirectory(deepPath)
		if err != nil {
			t.Logf("Deep directory creation result: %v", err)
			// This might fail on systems with path length limits
		} else {
			// Try to create a file in the deep directory
			deepFile := filepath.Join(deepPath, "deep-file.txt")
			err = gen.WriteFile(deepFile, []byte("deep content"), 0644)
			if err != nil {
				t.Logf("Deep file creation result: %v", err)
			}
		}
	})
}

func TestFileSystemGeneratorPerformance(t *testing.T) {
	gen := NewGenerator()
	tempDir := t.TempDir()

	t.Run("performance with many small files", func(t *testing.T) {
		const numFiles = 1000

		start := time.Now()

		for i := 0; i < numFiles; i++ {
			filePath := filepath.Join(tempDir, fmt.Sprintf("small-file-%d.txt", i))
			content := fmt.Sprintf("Content for file %d", i)

			err := gen.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				t.Errorf("Failed to create file %d: %v", i, err)
			}
		}

		duration := time.Since(start)
		t.Logf("Created %d small files in %v (avg: %v per file)",
			numFiles, duration, duration/time.Duration(numFiles))

		// Performance should be reasonable
		avgDuration := duration / time.Duration(numFiles)
		if avgDuration > 5*time.Millisecond {
			t.Errorf("File creation too slow: %v per file", avgDuration)
		}
	})

	t.Run("performance with large files", func(t *testing.T) {
		const fileSize = 10 * 1024 * 1024 // 10MB
		largeContent := make([]byte, fileSize)

		// Fill with some pattern
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}

		start := time.Now()

		largePath := filepath.Join(tempDir, "large-file.bin")
		err := gen.WriteFile(largePath, largeContent, 0644)

		duration := time.Since(start)

		if err != nil {
			t.Logf("Large file creation failed: %v", err)
		} else {
			t.Logf("Created %d MB file in %v", fileSize/(1024*1024), duration)

			// Verify file size
			if info, err := os.Stat(largePath); err == nil {
				if info.Size() != int64(fileSize) {
					t.Errorf("Expected file size %d, got %d", fileSize, info.Size())
				}
			}
		}
	})

	t.Run("directory creation performance", func(t *testing.T) {
		const numDirs = 500

		start := time.Now()

		for i := 0; i < numDirs; i++ {
			dirPath := filepath.Join(tempDir, fmt.Sprintf("perf-dir-%d", i))
			err := gen.CreateDirectory(dirPath)
			if err != nil {
				t.Errorf("Failed to create directory %d: %v", i, err)
			}
		}

		duration := time.Since(start)
		t.Logf("Created %d directories in %v (avg: %v per directory)",
			numDirs, duration, duration/time.Duration(numDirs))

		// Performance should be reasonable
		avgDuration := duration / time.Duration(numDirs)
		if avgDuration > 2*time.Millisecond {
			t.Errorf("Directory creation too slow: %v per directory", avgDuration)
		}
	})
}
