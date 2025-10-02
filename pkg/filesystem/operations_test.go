package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func createTestConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project description",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "18.17.0",
			Go:   "1.22.0",
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.2.0",
			},
		},
		OutputPath:       "/tmp",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}
}

func TestNewFileSystemOperations(t *testing.T) {
	fso := NewFileSystemOperations()
	if fso == nil {
		t.Fatal("NewFileSystemOperations() returned nil")
	}

	if fso.dryRun {
		t.Fatal("NewFileSystemOperations() should not be in dry-run mode by default")
	}
}

func TestNewDryRunFileSystemOperations(t *testing.T) {
	fso := NewDryRunFileSystemOperations()
	if fso == nil {
		t.Fatal("NewDryRunFileSystemOperations() returned nil")
	}

	if !fso.dryRun {
		t.Fatal("NewDryRunFileSystemOperations() should be in dry-run mode")
	}
}

func TestFileSystemOperationsCreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid directory creation",
			path:        filepath.Join(tempDir, "test-dir"),
			expectError: false,
		},
		{
			name:        "nested directory creation",
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
			name:        "path traversal attack",
			path:        tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.CreateDirectory(tt.path)

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
				if !fso.FileExists(tt.path) {
					t.Errorf("CreateDirectory() did not create directory: %s", tt.path)
				}

				// Verify it's actually a directory
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("CreateDirectory() failed to stat created directory: %v", err)
				} else if !info.IsDir() {
					t.Errorf("CreateDirectory() created path is not a directory")
				}
			}
		})
	}
}

func TestFileSystemOperationsEnsureDirectory(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test directory first
	existingDir := filepath.Join(tempDir, "existing")
	if err := os.MkdirAll(existingDir, 0750); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test file to test error case
	testFile := filepath.Join(tempDir, "testfile")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "ensure existing directory",
			path:        existingDir,
			expectError: false,
		},
		{
			name:        "ensure new directory",
			path:        filepath.Join(tempDir, "new-dir"),
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "directory path cannot be empty",
		},
		{
			name:        "path is a file",
			path:        testFile,
			expectError: true,
			errorMsg:    "exists but is not a directory",
		},
		{
			name:        "path traversal attack",
			path:        tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.EnsureDirectory(tt.path)

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
				if !fso.FileExists(tt.path) {
					t.Errorf("EnsureDirectory() directory does not exist: %s", tt.path)
				}
			}
		})
	}
}

func TestFileSystemOperationsWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	testContent := []byte("test file content")

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
			content:     testContent,
			perm:        0644,
			expectError: false,
		},
		{
			name:        "nested file write",
			path:        filepath.Join(tempDir, "nested", "test.txt"),
			content:     testContent,
			perm:        0644,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			content:     testContent,
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
			name:        "invalid permissions",
			path:        filepath.Join(tempDir, "invalid-perm.txt"),
			content:     testContent,
			perm:        01000, // Invalid permission (exceeds 0777)
			expectError: true,
			errorMsg:    "invalid file permissions",
		},
		{
			name:        "path traversal attack",
			path:        tempDir + "/../malicious.txt",
			content:     testContent,
			perm:        0644,
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.WriteFile(tt.path, tt.content, tt.perm)

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

				// Verify file was created
				if !fso.FileExists(tt.path) {
					t.Errorf("WriteFile() did not create file: %s", tt.path)
				}

				// Verify file content
				content, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("WriteFile() failed to read created file: %v", err)
				} else if string(content) != string(tt.content) {
					t.Errorf("WriteFile() content mismatch: got %s, want %s", content, tt.content)
				}

				// Verify file permissions
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("WriteFile() failed to stat created file: %v", err)
				} else if info.Mode().Perm() != tt.perm {
					t.Errorf("WriteFile() permission mismatch: got %o, want %o", info.Mode().Perm(), tt.perm)
				}
			}
		})
	}
}

func TestFileSystemOperationsFileExists(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test directory
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0750); err != nil {
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
			path:     filepath.Join(tempDir, "nonexistent.txt"),
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
			result := fso.FileExists(tt.path)
			if result != tt.expected {
				t.Errorf("FileExists() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileSystemOperationsCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	srcContent := []byte("source file content")
	if err := os.WriteFile(srcFile, srcContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	tests := []struct {
		name        string
		srcPath     string
		destPath    string
		perm        os.FileMode
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid file copy",
			srcPath:     srcFile,
			destPath:    filepath.Join(tempDir, "dest.txt"),
			perm:        0644,
			expectError: false,
		},
		{
			name:        "copy to nested directory",
			srcPath:     srcFile,
			destPath:    filepath.Join(tempDir, "nested", "dest.txt"),
			perm:        0644,
			expectError: false,
		},
		{
			name:        "empty source path",
			srcPath:     "",
			destPath:    filepath.Join(tempDir, "dest.txt"),
			perm:        0644,
			expectError: true,
			errorMsg:    "source path cannot be empty",
		},
		{
			name:        "empty destination path",
			srcPath:     srcFile,
			destPath:    "",
			perm:        0644,
			expectError: true,
			errorMsg:    "destination path cannot be empty",
		},
		{
			name:        "non-existent source file",
			srcPath:     filepath.Join(tempDir, "nonexistent.txt"),
			destPath:    filepath.Join(tempDir, "dest.txt"),
			perm:        0644,
			expectError: true,
			errorMsg:    "source file does not exist",
		},
		{
			name:        "invalid permissions",
			srcPath:     srcFile,
			destPath:    filepath.Join(tempDir, "dest-invalid-perm.txt"),
			perm:        01000, // Invalid permission (exceeds 0777)
			expectError: true,
			errorMsg:    "invalid file permissions",
		},
		{
			name:        "path traversal in source",
			srcPath:     tempDir + "/../malicious.txt",
			destPath:    filepath.Join(tempDir, "dest.txt"),
			perm:        0644,
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "path traversal in destination",
			srcPath:     srcFile,
			destPath:    tempDir + "/../malicious.txt",
			perm:        0644,
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.CopyFile(tt.srcPath, tt.destPath, tt.perm)

			if tt.expectError {
				if err == nil {
					t.Errorf("CopyFile() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CopyFile() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CopyFile() unexpected error = %v", err)
				}

				// Verify destination file was created
				if !fso.FileExists(tt.destPath) {
					t.Errorf("CopyFile() did not create destination file: %s", tt.destPath)
				}

				// Verify file content
				destContent, err := os.ReadFile(tt.destPath)
				if err != nil {
					t.Errorf("CopyFile() failed to read destination file: %v", err)
				} else if string(destContent) != string(srcContent) {
					t.Errorf("CopyFile() content mismatch: got %s, want %s", destContent, srcContent)
				}

				// Verify file permissions
				info, err := os.Stat(tt.destPath)
				if err != nil {
					t.Errorf("CopyFile() failed to stat destination file: %v", err)
				} else if info.Mode().Perm() != tt.perm {
					t.Errorf("CopyFile() permission mismatch: got %o, want %o", info.Mode().Perm(), tt.perm)
				}
			}
		})
	}
}

func TestFileSystemOperationsCopyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0750); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create files in source directory
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create source subfile: %v", err)
	}

	tests := []struct {
		name        string
		srcDir      string
		destDir     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid directory copy",
			srcDir:      srcDir,
			destDir:     filepath.Join(tempDir, "dest"),
			expectError: false,
		},
		{
			name:        "empty source directory",
			srcDir:      "",
			destDir:     filepath.Join(tempDir, "dest"),
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
			srcDir:      filepath.Join(tempDir, "nonexistent"),
			destDir:     filepath.Join(tempDir, "dest"),
			expectError: true,
			errorMsg:    "source directory does not exist",
		},
		{
			name:        "path traversal in source",
			srcDir:      tempDir + "/../malicious",
			destDir:     filepath.Join(tempDir, "dest"),
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
			err := fso.CopyDirectory(tt.srcDir, tt.destDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("CopyDirectory() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CopyDirectory() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CopyDirectory() unexpected error = %v", err)
				}

				// Verify destination directory was created
				if !fso.FileExists(tt.destDir) {
					t.Errorf("CopyDirectory() did not create destination directory: %s", tt.destDir)
				}

				// Verify subdirectory was copied
				subDir := filepath.Join(tt.destDir, "subdir")
				if !fso.FileExists(subDir) {
					t.Errorf("CopyDirectory() did not copy subdirectory: %s", subDir)
				}

				// Verify files were copied
				file1 := filepath.Join(tt.destDir, "file1.txt")
				if !fso.FileExists(file1) {
					t.Errorf("CopyDirectory() did not copy file1.txt")
				}

				file2 := filepath.Join(tt.destDir, "subdir", "file2.txt")
				if !fso.FileExists(file2) {
					t.Errorf("CopyDirectory() did not copy subdir/file2.txt")
				}

				// Verify file contents
				content1, err := os.ReadFile(file1)
				if err != nil {
					t.Errorf("CopyDirectory() failed to read copied file1: %v", err)
				} else if string(content1) != "content1" {
					t.Errorf("CopyDirectory() file1 content mismatch: got %s, want content1", content1)
				}

				content2, err := os.ReadFile(file2)
				if err != nil {
					t.Errorf("CopyDirectory() failed to read copied file2: %v", err)
				} else if string(content2) != "content2" {
					t.Errorf("CopyDirectory() file2 content mismatch: got %s, want content2", content2)
				}
			}
		})
	}
}

func TestFileSystemOperationsCreateSymlink(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create target file
	targetFile := filepath.Join(tempDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
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
			link:        filepath.Join(tempDir, "link.txt"),
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
			link:        tempDir + "/../malicious.txt",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.CreateSymlink(tt.target, tt.link)

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
				if !fso.FileExists(tt.link) {
					t.Errorf("CreateSymlink() did not create symlink: %s", tt.link)
				}

				// Verify it's actually a symlink
				info, err := os.Lstat(tt.link)
				if err != nil {
					t.Errorf("CreateSymlink() failed to lstat symlink: %v", err)
				} else if info.Mode()&os.ModeSymlink == 0 {
					t.Errorf("CreateSymlink() created path is not a symlink")
				}

				// Verify symlink target
				target, err := os.Readlink(tt.link)
				if err != nil {
					t.Errorf("CreateSymlink() failed to read symlink target: %v", err)
				} else if target != tt.target {
					t.Errorf("CreateSymlink() target mismatch: got %s, want %s", target, tt.target)
				}
			}
		})
	}
}

func TestFileSystemOperationsValidateFileContent(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid file content validation",
			path:        testFile,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "non-existent file",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
			errorMsg:    "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.ValidateFileContent(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateFileContent() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateFileContent() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFileContent() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileSystemOperationsValidateDirectoryStructure(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create test directory structure
	testDirs := []string{"dir1", "dir2", "nested/dir3"}
	for _, dir := range testDirs {
		if err := os.MkdirAll(filepath.Join(tempDir, dir), 0750); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create a test file (not a directory)
	testFile := filepath.Join(tempDir, "notadir")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		basePath     string
		requiredDirs []string
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "valid directory structure",
			basePath:     tempDir,
			requiredDirs: []string{"dir1", "dir2"},
			expectError:  false,
		},
		{
			name:         "nested directory validation",
			basePath:     tempDir,
			requiredDirs: []string{"nested/dir3"},
			expectError:  false,
		},
		{
			name:         "empty base path",
			basePath:     "",
			requiredDirs: []string{"dir1"},
			expectError:  true,
			errorMsg:     "base path cannot be empty",
		},
		{
			name:         "empty required directories",
			basePath:     tempDir,
			requiredDirs: []string{},
			expectError:  true,
			errorMsg:     "required directories list cannot be empty",
		},
		{
			name:         "non-existent base directory",
			basePath:     filepath.Join(tempDir, "nonexistent"),
			requiredDirs: []string{"dir1"},
			expectError:  true,
			errorMsg:     "base directory does not exist",
		},
		{
			name:         "missing required directory",
			basePath:     tempDir,
			requiredDirs: []string{"dir1", "missing"},
			expectError:  true,
			errorMsg:     "required directory missing",
		},
		{
			name:         "path is file not directory",
			basePath:     tempDir,
			requiredDirs: []string{"notadir"},
			expectError:  true,
			errorMsg:     "exists but is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.ValidateDirectoryStructure(tt.basePath, tt.requiredDirs)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateDirectoryStructure() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateDirectoryStructure() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateDirectoryStructure() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileSystemOperationsValidateFileStructure(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "nested/file3.txt"}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
			t.Fatalf("Failed to create parent directory for %s: %v", file, err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create a test directory (not a file)
	testDir := filepath.Join(tempDir, "notafile")
	if err := os.MkdirAll(testDir, 0750); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name          string
		basePath      string
		requiredFiles []string
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "valid file structure",
			basePath:      tempDir,
			requiredFiles: []string{"file1.txt", "file2.txt"},
			expectError:   false,
		},
		{
			name:          "nested file validation",
			basePath:      tempDir,
			requiredFiles: []string{"nested/file3.txt"},
			expectError:   false,
		},
		{
			name:          "empty base path",
			basePath:      "",
			requiredFiles: []string{"file1.txt"},
			expectError:   true,
			errorMsg:      "base path cannot be empty",
		},
		{
			name:          "empty required files",
			basePath:      tempDir,
			requiredFiles: []string{},
			expectError:   true,
			errorMsg:      "required files list cannot be empty",
		},
		{
			name:          "non-existent base directory",
			basePath:      filepath.Join(tempDir, "nonexistent"),
			requiredFiles: []string{"file1.txt"},
			expectError:   true,
			errorMsg:      "base directory does not exist",
		},
		{
			name:          "missing required file",
			basePath:      tempDir,
			requiredFiles: []string{"file1.txt", "missing.txt"},
			expectError:   true,
			errorMsg:      "required file missing",
		},
		{
			name:          "path is directory not file",
			basePath:      tempDir,
			requiredFiles: []string{"notafile"},
			expectError:   true,
			errorMsg:      "exists but is a directory, expected file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.ValidateFileStructure(tt.basePath, tt.requiredFiles)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateFileStructure() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateFileStructure() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFileStructure() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileSystemOperationsCreateProjectRoot(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()
	config := createTestConfig()

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		outputPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project root creation",
			config:      config,
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
			config:      config,
			outputPath:  "",
			expectError: true,
			errorMsg:    "output path cannot be empty",
		},
		{
			name: "empty project name",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "test-org",
			},
			outputPath:  tempDir,
			expectError: true,
			errorMsg:    "project name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath, err := fso.CreateProjectRoot(tt.config, tt.outputPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("CreateProjectRoot() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("CreateProjectRoot() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CreateProjectRoot() unexpected error = %v", err)
				}

				// Verify project directory was created
				if !fso.FileExists(projectPath) {
					t.Errorf("CreateProjectRoot() did not create project directory: %s", projectPath)
				}

				// Verify returned path is correct
				expectedPath := filepath.Join(tt.outputPath, tt.config.Name)
				absExpectedPath, _ := filepath.Abs(expectedPath)
				if projectPath != absExpectedPath {
					t.Errorf("CreateProjectRoot() returned path = %s, want %s", projectPath, absExpectedPath)
				}
			}
		})
	}
}

func TestFileSystemOperationsValidateProjectRoot(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()
	config := createTestConfig()

	// Create a valid project directory
	validProjectPath := filepath.Join(tempDir, "valid-project")
	if err := os.MkdirAll(validProjectPath, 0750); err != nil {
		t.Fatalf("Failed to create valid project directory: %v", err)
	}

	// Create a file (not a directory)
	testFile := filepath.Join(tempDir, "notadir")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		projectPath string
		config      *models.ProjectConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project root",
			projectPath: validProjectPath,
			config:      config,
			expectError: false,
		},
		{
			name:        "empty project path",
			projectPath: "",
			config:      config,
			expectError: true,
			errorMsg:    "project path cannot be empty",
		},
		{
			name:        "nil config",
			projectPath: validProjectPath,
			config:      nil,
			expectError: true,
			errorMsg:    "project config cannot be nil",
		},
		{
			name:        "non-existent project directory",
			projectPath: filepath.Join(tempDir, "nonexistent"),
			config:      config,
			expectError: true,
			errorMsg:    "project directory does not exist",
		},
		{
			name:        "path is file not directory",
			projectPath: testFile,
			config:      config,
			expectError: true,
			errorMsg:    "exists but is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.ValidateProjectRoot(tt.projectPath, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateProjectRoot() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateProjectRoot() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateProjectRoot() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileSystemOperationsDryRunMode(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewDryRunFileSystemOperations()

	// Test that dry-run mode doesn't actually create files/directories
	testDir := filepath.Join(tempDir, "dry-run-test")
	testFile := filepath.Join(tempDir, "dry-run-file.txt")

	// These operations should not fail in dry-run mode
	if err := fso.CreateDirectory(testDir); err != nil {
		t.Errorf("CreateDirectory() in dry-run mode failed: %v", err)
	}

	if err := fso.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Errorf("WriteFile() in dry-run mode failed: %v", err)
	}

	// Verify that nothing was actually created
	if fso.FileExists(testDir) {
		t.Errorf("CreateDirectory() in dry-run mode actually created directory")
	}

	if fso.FileExists(testFile) {
		t.Errorf("WriteFile() in dry-run mode actually created file")
	}

	// Test IsDryRun method
	if !fso.IsDryRun() {
		t.Errorf("IsDryRun() should return true for dry-run operations")
	}

	// Test SetDryRun method
	fso.SetDryRun(false)
	if fso.IsDryRun() {
		t.Errorf("SetDryRun(false) should disable dry-run mode")
	}

	fso.SetDryRun(true)
	if !fso.IsDryRun() {
		t.Errorf("SetDryRun(true) should enable dry-run mode")
	}
}

func TestFileSystemOperationsGetFileInfo(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test directory
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0750); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
		checkIsDir  *bool
	}{
		{
			name:        "get file info for file",
			path:        testFile,
			expectError: false,
			checkIsDir:  &[]bool{false}[0],
		},
		{
			name:        "get file info for directory",
			path:        testDir,
			expectError: false,
			checkIsDir:  &[]bool{true}[0],
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "non-existent path",
			path:        filepath.Join(tempDir, "nonexistent"),
			expectError: true,
			errorMsg:    "failed to get file info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := fso.GetFileInfo(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("GetFileInfo() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("GetFileInfo() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("GetFileInfo() unexpected error = %v", err)
				}

				if info == nil {
					t.Errorf("GetFileInfo() returned nil info")
				}

				if tt.checkIsDir != nil {
					if info.IsDir() != *tt.checkIsDir {
						t.Errorf("GetFileInfo() IsDir() = %v, want %v", info.IsDir(), *tt.checkIsDir)
					}
				}
			}
		})
	}
}

func TestFileSystemOperationsRemoveFile(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "remove existing file",
			path:        testFile,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "path traversal attack",
			path:        tempDir + "/../malicious.txt",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "non-existent file",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
			errorMsg:    "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.RemoveFile(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("RemoveFile() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("RemoveFile() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("RemoveFile() unexpected error = %v", err)
				}

				// Verify file was removed
				if fso.FileExists(tt.path) {
					t.Errorf("RemoveFile() did not remove file: %s", tt.path)
				}
			}
		})
	}
}

func TestFileSystemOperationsRemoveDirectory(t *testing.T) {
	tempDir := t.TempDir()
	fso := NewFileSystemOperations()

	// Create a test directory with content
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(filepath.Join(testDir, "subdir"), 0750); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test file (not a directory)
	testFile := filepath.Join(tempDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "remove existing directory",
			path:        testDir,
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
			errorMsg:    "directory path cannot be empty",
		},
		{
			name:        "path traversal attack",
			path:        tempDir + "/../malicious",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "non-existent directory",
			path:        filepath.Join(tempDir, "nonexistent"),
			expectError: true,
			errorMsg:    "directory does not exist",
		},
		{
			name:        "path is file not directory",
			path:        testFile,
			expectError: true,
			errorMsg:    "exists but is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fso.RemoveDirectory(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("RemoveDirectory() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("RemoveDirectory() error = %v, expected to contain %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("RemoveDirectory() unexpected error = %v", err)
				}

				// Verify directory was removed
				if fso.FileExists(tt.path) {
					t.Errorf("RemoveDirectory() did not remove directory: %s", tt.path)
				}
			}
		})
	}
}
