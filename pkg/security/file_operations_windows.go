//go:build windows
// +build windows

package security

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SecureFileOperations provides secure file system operations with validation
type SecureFileOperations struct {
	allowedBasePaths []string
	maxFileSize      int64
	allowedExts      []string
	blockedPaths     []string
	requirePermCheck bool
}

// FileOperationResult represents the result of a file operation
type FileOperationResult struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	FilePath     string    `json:"file_path"`
	Operation    string    `json:"operation"`
	BytesWritten int64     `json:"bytes_written,omitempty"`
	BytesRead    int64     `json:"bytes_read,omitempty"`
	Checksum     string    `json:"checksum,omitempty"`
	Permissions  string    `json:"permissions,omitempty"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modified_time"`
	Error        string    `json:"error,omitempty"`
	Warnings     []string  `json:"warnings,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// NewSecureFileOperations creates a new secure file operations manager
func NewSecureFileOperations(allowedBasePaths []string) *SecureFileOperations {
	return &SecureFileOperations{
		allowedBasePaths: allowedBasePaths,
		maxFileSize:      100 * 1024 * 1024, // 100MB default
		allowedExts: []string{
			".go", ".js", ".ts", ".jsx", ".tsx", ".json", ".yaml", ".yml",
			".md", ".txt", ".html", ".css", ".scss", ".sql", ".sh", ".bat",
			".dockerfile", ".gitignore", ".env", ".toml", ".ini", ".conf", ".backup",
		},
		blockedPaths: []string{
			"/etc/passwd", "/etc/shadow", "/etc/hosts", "/etc/hostname",
			"/proc/", "/sys/", "/dev/", "/root/", "/home/", "/var/log/",
		},
		requirePermCheck: true,
	}
}

// ValidateFilePath performs comprehensive file path validation
func (sfo *SecureFileOperations) ValidateFilePath(path string, operation string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if path is within allowed base paths
	allowed := false
	for _, basePath := range sfo.allowedBasePaths {
		if strings.HasPrefix(absPath, basePath) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("path not in allowed base paths: %s", absPath)
	}

	// Check if path is blocked
	for _, blockedPath := range sfo.blockedPaths {
		if strings.Contains(absPath, blockedPath) {
			return fmt.Errorf("path is blocked: %s", absPath)
		}
	}

	return nil
}

// SecureReadFile reads a file with security validation
func (sfo *SecureFileOperations) SecureReadFile(path string) (*FileOperationResult, []byte, error) {
	if err := sfo.ValidateFilePath(path, "read"); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("path validation failed: %v", err),
			FilePath: path,
		}, nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to read file: %v", err),
			FilePath: path,
		}, nil, err
	}

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get file info: %v", err),
			FilePath: path,
		}, nil, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(path)

	return &FileOperationResult{
		Success:      true,
		Message:      "file read successfully",
		FilePath:     path,
		Operation:    "read",
		BytesRead:    int64(len(data)),
		Size:         fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime(),
		Permissions:  permissionsToString(permissions),
		Timestamp:    time.Now(),
	}, data, nil
}

// SecureWriteFile writes a file with security validation
func (sfo *SecureFileOperations) SecureWriteFile(path string, data []byte, perm os.FileMode) (*FileOperationResult, error) {
	if err := sfo.ValidateFilePath(path, "write"); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("path validation failed: %v", err),
			FilePath: path,
		}, err
	}

	// Check data size
	if int64(len(data)) > sfo.maxFileSize {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("data size exceeds maximum allowed size: %d > %d", len(data), sfo.maxFileSize),
			FilePath: path,
		}, fmt.Errorf("data too large")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	allowed := false
	for _, allowedExt := range sfo.allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("file extension not allowed: %s", ext),
			FilePath: path,
		}, fmt.Errorf("file extension not allowed")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to create directory: %v", err),
			FilePath: path,
		}, err
	}

	// Write file
	if err := os.WriteFile(path, data, perm); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to write file: %v", err),
			FilePath: path,
		}, err
	}

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get file info: %v", err),
			FilePath: path,
		}, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(path)

	return &FileOperationResult{
		Success:      true,
		Message:      "file written successfully",
		FilePath:     path,
		Operation:    "write",
		BytesWritten: int64(len(data)),
		Size:         fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime(),
		Permissions:  permissionsToString(permissions),
		Timestamp:    time.Now(),
	}, nil
}

// SecureCopyFile copies a file with security validation
func (sfo *SecureFileOperations) SecureCopyFile(srcPath, dstPath string) (*FileOperationResult, error) {
	if err := sfo.ValidateFilePath(srcPath, "read"); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("source path validation failed: %v", err),
			FilePath: srcPath,
		}, err
	}

	if err := sfo.ValidateFilePath(dstPath, "write"); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("destination path validation failed: %v", err),
			FilePath: dstPath,
		}, err
	}

	// Check source file size
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get source file info: %v", err),
			FilePath: srcPath,
		}, err
	}

	if fileInfo.Size() > sfo.maxFileSize {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("file size exceeds maximum allowed size: %d > %d", fileInfo.Size(), sfo.maxFileSize),
			FilePath: srcPath,
		}, fmt.Errorf("file too large")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(srcPath))
	allowed := false
	for _, allowedExt := range sfo.allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("file extension not allowed: %s", ext),
			FilePath: srcPath,
		}, fmt.Errorf("file extension not allowed")
	}

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to create destination directory: %v", err),
			FilePath: dstPath,
		}, err
	}

	// Perform the copy
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to open source file: %v", err),
			FilePath: srcPath,
		}, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to create destination file: %v", err),
			FilePath: dstPath,
		}, err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to copy file: %v", err),
			FilePath: srcPath,
		}, err
	}

	// Get destination file info
	dstFileInfo, err := dstFile.Stat()
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get destination file info: %v", err),
			FilePath: dstPath,
		}, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(dstPath)

	return &FileOperationResult{
		Success:      true,
		Message:      "file copied successfully",
		FilePath:     dstPath,
		Operation:    "copy",
		BytesWritten: dstFileInfo.Size(),
		Size:         dstFileInfo.Size(),
		ModifiedTime: dstFileInfo.ModTime(),
		Permissions:  permissionsToString(permissions),
		Timestamp:    time.Now(),
	}, nil
}

// SecureCreateDirectory creates a directory with security validation
func (sfo *SecureFileOperations) SecureCreateDirectory(path string, perm os.FileMode) (*FileOperationResult, error) {
	if err := sfo.ValidateFilePath(path, "write"); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("path validation failed: %v", err),
			FilePath: path,
		}, err
	}

	// Create directory
	if err := os.MkdirAll(path, perm); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to create directory: %v", err),
			FilePath: path,
		}, err
	}

	// Get directory info
	dirInfo, err := os.Stat(path)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get directory info: %v", err),
			FilePath: path,
		}, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(path)

	return &FileOperationResult{
		Success:      true,
		Message:      "directory created successfully",
		FilePath:     path,
		Operation:    "create_directory",
		Size:         dirInfo.Size(),
		ModifiedTime: dirInfo.ModTime(),
		Permissions:  permissionsToString(permissions),
		Timestamp:    time.Now(),
	}, nil
}

// ValidatePath validates if a path is safe to operate on
func (sfo *SecureFileOperations) ValidatePath(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if path is within allowed base paths
	allowed := false
	for _, basePath := range sfo.allowedBasePaths {
		if strings.HasPrefix(absPath, basePath) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("path not in allowed base paths: %s", absPath)
	}

	// Check if path is blocked
	for _, blockedPath := range sfo.blockedPaths {
		if strings.Contains(absPath, blockedPath) {
			return fmt.Errorf("path is blocked: %s", absPath)
		}
	}

	return nil
}

// GetFilePermissions gets file permissions (Windows version)
func (sfo *SecureFileOperations) GetFilePermissions(path string) (map[string]interface{}, error) {
	permissions := make(map[string]interface{})

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Basic permissions
	permissions["is_dir"] = fileInfo.IsDir()
	permissions["is_regular"] = fileInfo.Mode().IsRegular()
	permissions["mode"] = fileInfo.Mode().String()
	permissions["size"] = fileInfo.Size()
	permissions["mod_time"] = fileInfo.ModTime()

	// Check if writable (Windows doesn't have Unix-style UID/GID)
	if sfo.requirePermCheck {
		tempPath := filepath.Join(path, ".temp_perm_check")
		// #nosec G304 -- tempPath is constructed safely
		if file, err := os.Create(tempPath); err == nil {
			permissions["is_writable"] = true
			_ = file.Close()
			_ = os.Remove(tempPath)
		}
	}

	// Windows doesn't have Unix-style UID/GID, so we skip that part
	// The sys() method on Windows returns different types

	return permissions, nil
}

// ValidateFileIntegrity validates file integrity using checksums
func (sfo *SecureFileOperations) ValidateFileIntegrity(path string, expectedChecksum string) (*FileOperationResult, error) {
	if err := sfo.ValidatePath(path); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("path validation failed: %v", err),
			FilePath: path,
		}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to open file: %v", err),
			FilePath: path,
		}, err
	}
	defer file.Close()

	// Calculate checksum
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to calculate checksum: %v", err),
			FilePath: path,
		}, err
	}

	actualChecksum := fmt.Sprintf("%x", hasher.Sum(nil))
	success := actualChecksum == expectedChecksum

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get file info: %v", err),
			FilePath: path,
		}, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(path)

	return &FileOperationResult{
		Success:      success,
		Message:      fmt.Sprintf("checksum validation %s", map[bool]string{true: "passed", false: "failed"}[success]),
		FilePath:     path,
		Checksum:     actualChecksum,
		Permissions:  permissionsToString(permissions),
		Size:         fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime(),
	}, nil
}

// SecureCopy copies a file with validation
func (sfo *SecureFileOperations) SecureCopy(src, dst string) (*FileOperationResult, error) {
	if err := sfo.ValidatePath(src); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("source path validation failed: %v", err),
			FilePath: src,
		}, err
	}

	if err := sfo.ValidatePath(dst); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("destination path validation failed: %v", err),
			FilePath: dst,
		}, err
	}

	// Check file size
	fileInfo, err := os.Stat(src)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get source file info: %v", err),
			FilePath: src,
		}, err
	}

	if fileInfo.Size() > sfo.maxFileSize {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("file size exceeds maximum allowed size: %d > %d", fileInfo.Size(), sfo.maxFileSize),
			FilePath: src,
		}, fmt.Errorf("file too large")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(src))
	allowed := false
	for _, allowedExt := range sfo.allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("file extension not allowed: %s", ext),
			FilePath: src,
		}, fmt.Errorf("file extension not allowed")
	}

	// Perform the copy
	srcFile, err := os.Open(src)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to open source file: %v", err),
			FilePath: src,
		}, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to create destination file: %v", err),
			FilePath: dst,
		}, err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to copy file: %v", err),
			FilePath: src,
		}, err
	}

	// Get destination file info
	dstFileInfo, err := dstFile.Stat()
	if err != nil {
		return &FileOperationResult{
			Success:  false,
			Message:  fmt.Sprintf("failed to get destination file info: %v", err),
			FilePath: dst,
		}, err
	}

	// Get permissions
	permissions, _ := sfo.GetFilePermissions(dst)

	return &FileOperationResult{
		Success:      true,
		Message:      "file copied successfully",
		FilePath:     dst,
		Size:         dstFileInfo.Size(),
		ModifiedTime: dstFileInfo.ModTime(),
		Permissions:  permissionsToString(permissions),
	}, nil
}

// permissionsToString converts permissions map to JSON string
func permissionsToString(permissions map[string]interface{}) string {
	if permissions == nil {
		return ""
	}
	jsonBytes, err := json.Marshal(permissions)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// SetSecurityConfig sets security configuration
func (sfo *SecureFileOperations) SetSecurityConfig(config map[string]interface{}) error {
	// Update configuration based on provided config
	if maxSize, ok := config["max_file_size"].(int64); ok {
		sfo.maxFileSize = maxSize
	}
	if allowedExts, ok := config["allowed_extensions"].([]string); ok {
		sfo.allowedExts = allowedExts
	}
	if blockedPaths, ok := config["blocked_paths"].([]string); ok {
		sfo.blockedPaths = blockedPaths
	}
	if requirePermCheck, ok := config["require_permission_check"].(bool); ok {
		sfo.requirePermCheck = requirePermCheck
	}
	return nil
}
