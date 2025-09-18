package security

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
			"/etc/passwd", "/etc/shadow", "/etc/hosts", "/etc/sudoers",
			"/proc", "/sys", "/dev", "/boot", "/root",
			"C:\\Windows\\System32", "C:\\Windows\\SysWOW64",
		},
		requirePermCheck: true,
	}
}

// FileOperationResult contains the result of a secure file operation
type FileOperationResult struct {
	Success      bool      `json:"success"`
	FilePath     string    `json:"file_path"`
	Operation    string    `json:"operation"`
	BytesWritten int64     `json:"bytes_written,omitempty"`
	BytesRead    int64     `json:"bytes_read,omitempty"`
	Checksum     string    `json:"checksum,omitempty"`
	Permissions  string    `json:"permissions,omitempty"`
	Error        string    `json:"error,omitempty"`
	Warnings     []string  `json:"warnings,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// ValidateFilePath performs comprehensive file path validation
func (sfo *SecureFileOperations) ValidateFilePath(path string, operation string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in path: %s", path)
	}

	// Check against blocked paths
	for _, blockedPath := range sfo.blockedPaths {
		if strings.HasPrefix(strings.ToLower(cleanPath), strings.ToLower(blockedPath)) {
			return fmt.Errorf("access to blocked path denied: %s", path)
		}
	}

	// Validate against allowed base paths if specified
	if len(sfo.allowedBasePaths) > 0 {
		allowed := false
		for _, basePath := range sfo.allowedBasePaths {
			cleanBasePath := filepath.Clean(basePath)
			if strings.HasPrefix(cleanPath, cleanBasePath) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("path not within allowed directories: %s", path)
		}
	}

	// Check file extension for write operations
	if operation == "write" || operation == "create" {
		ext := strings.ToLower(filepath.Ext(cleanPath))
		if ext != "" { // Only check if there's an extension
			allowed := false
			for _, allowedExt := range sfo.allowedExts {
				if ext == allowedExt {
					allowed = true
					break
				}
			}
			if !allowed {
				return fmt.Errorf("file extension '%s' is not allowed", ext)
			}
		}
	}

	return nil
}

// SecureReadFile reads a file with security validation
func (sfo *SecureFileOperations) SecureReadFile(path string) (*FileOperationResult, []byte, error) {
	result := &FileOperationResult{
		FilePath:  path,
		Operation: "read",
		Timestamp: time.Now(),
	}

	// Validate path
	if err := sfo.ValidateFilePath(path, "read"); err != nil {
		result.Error = err.Error()
		return result, nil, err
	}

	// Check file exists and get info
	fileInfo, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Sprintf("failed to stat file: %v", err)
		return result, nil, err
	}

	// Check file size
	if fileInfo.Size() > sfo.maxFileSize {
		err := fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", fileInfo.Size(), sfo.maxFileSize)
		result.Error = err.Error()
		return result, nil, err
	}

	// Check permissions
	if sfo.requirePermCheck {
		if err := sfo.checkReadPermissions(path); err != nil {
			result.Error = fmt.Sprintf("permission check failed: %v", err)
			result.Warnings = append(result.Warnings, "File permissions may be too restrictive")
		}
	}

	// Read file
	data, err := os.ReadFile(path) // #nosec G304 - path is validated above
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file: %v", err)
		return result, nil, err
	}

	// Calculate checksum
	hash := sha256.Sum256(data)
	result.Checksum = fmt.Sprintf("%x", hash)
	result.BytesRead = int64(len(data))
	result.Permissions = fileInfo.Mode().String()
	result.Success = true

	return result, data, nil
}

// SecureWriteFile writes a file with security validation
func (sfo *SecureFileOperations) SecureWriteFile(path string, data []byte, perm os.FileMode) (*FileOperationResult, error) {
	result := &FileOperationResult{
		FilePath:  path,
		Operation: "write",
		Timestamp: time.Now(),
	}

	// Validate path
	if err := sfo.ValidateFilePath(path, "write"); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Check data size
	if int64(len(data)) > sfo.maxFileSize {
		err := fmt.Errorf("data size (%d bytes) exceeds maximum allowed size (%d bytes)", len(data), sfo.maxFileSize)
		result.Error = err.Error()
		return result, err
	}

	// Ensure secure permissions
	if perm > 0666 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Permissions %o may be too permissive, using 0644 instead", perm))
		perm = 0644
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		result.Error = fmt.Sprintf("failed to create directory: %v", err)
		return result, err
	}

	// Write file
	if err := os.WriteFile(path, data, perm); err != nil {
		result.Error = fmt.Sprintf("failed to write file: %v", err)
		return result, err
	}

	// Calculate checksum
	hash := sha256.Sum256(data)
	result.Checksum = fmt.Sprintf("%x", hash)
	result.BytesWritten = int64(len(data))
	result.Permissions = perm.String()
	result.Success = true

	return result, nil
}

// SecureCopyFile copies a file with security validation
func (sfo *SecureFileOperations) SecureCopyFile(srcPath, dstPath string) (*FileOperationResult, error) {
	result := &FileOperationResult{
		FilePath:  fmt.Sprintf("%s -> %s", srcPath, dstPath),
		Operation: "copy",
		Timestamp: time.Now(),
	}

	// Validate both paths
	if err := sfo.ValidateFilePath(srcPath, "read"); err != nil {
		result.Error = fmt.Sprintf("source path validation failed: %v", err)
		return result, err
	}

	if err := sfo.ValidateFilePath(dstPath, "write"); err != nil {
		result.Error = fmt.Sprintf("destination path validation failed: %v", err)
		return result, err
	}

	// Check source file
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to stat source file: %v", err)
		return result, err
	}

	if srcInfo.Size() > sfo.maxFileSize {
		err := fmt.Errorf("source file size (%d bytes) exceeds maximum allowed size (%d bytes)", srcInfo.Size(), sfo.maxFileSize)
		result.Error = err.Error()
		return result, err
	}

	// Open source file
	srcFile, err := os.Open(srcPath) // #nosec G304 - path is validated above
	if err != nil {
		result.Error = fmt.Sprintf("failed to open source file: %v", err)
		return result, err
	}
	defer func() {
		if closeErr := srcFile.Close(); closeErr != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to close source file: %v", closeErr))
		}
	}()

	// Create destination directory
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		result.Error = fmt.Sprintf("failed to create destination directory: %v", err)
		return result, err
	}

	// Create destination file with secure permissions
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // #nosec G304 - path is validated above
	if err != nil {
		result.Error = fmt.Sprintf("failed to create destination file: %v", err)
		return result, err
	}
	defer func() {
		if closeErr := dstFile.Close(); closeErr != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("failed to close destination file: %v", closeErr))
		}
	}()

	// Copy data with checksum calculation
	hash := sha256.New()
	multiWriter := io.MultiWriter(dstFile, hash)

	bytesWritten, err := io.Copy(multiWriter, srcFile)
	if err != nil {
		result.Error = fmt.Sprintf("failed to copy file data: %v", err)
		return result, err
	}

	result.Checksum = fmt.Sprintf("%x", hash.Sum(nil))
	result.BytesWritten = bytesWritten
	result.Permissions = "0644"
	result.Success = true

	return result, nil
}

// SecureCreateDirectory creates a directory with security validation
func (sfo *SecureFileOperations) SecureCreateDirectory(path string, perm os.FileMode) (*FileOperationResult, error) {
	result := &FileOperationResult{
		FilePath:  path,
		Operation: "mkdir",
		Timestamp: time.Now(),
	}

	// Validate path
	if err := sfo.ValidateFilePath(path, "write"); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Ensure secure permissions
	if perm > 0755 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Directory permissions %o may be too permissive, using 0755 instead", perm))
		perm = 0755
	}

	// Create directory
	if err := os.MkdirAll(path, perm); err != nil {
		result.Error = fmt.Sprintf("failed to create directory: %v", err)
		return result, err
	}

	result.Permissions = perm.String()
	result.Success = true

	return result, nil
}

// SecureRemoveFile removes a file with security validation
func (sfo *SecureFileOperations) SecureRemoveFile(path string) (*FileOperationResult, error) {
	result := &FileOperationResult{
		FilePath:  path,
		Operation: "remove",
		Timestamp: time.Now(),
	}

	// Validate path
	if err := sfo.ValidateFilePath(path, "write"); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "File does not exist")
		result.Success = true
		return result, nil
	}

	// Remove file
	if err := os.Remove(path); err != nil {
		result.Error = fmt.Sprintf("failed to remove file: %v", err)
		return result, err
	}

	result.Success = true
	return result, nil
}

// checkReadPermissions checks if the current process can read the file
func (sfo *SecureFileOperations) checkReadPermissions(path string) error {
	file, err := os.Open(path) // #nosec G304 - path is validated by caller
	if err != nil {
		return fmt.Errorf("cannot open file for reading: %w", err)
	}
	_ = file.Close()
	return nil
}

// GetFilePermissions returns detailed file permission information
func (sfo *SecureFileOperations) GetFilePermissions(path string) (map[string]interface{}, error) {
	if err := sfo.ValidateFilePath(path, "read"); err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	permissions := map[string]interface{}{
		"mode":          fileInfo.Mode().String(),
		"mode_octal":    fmt.Sprintf("%o", fileInfo.Mode().Perm()),
		"is_dir":        fileInfo.IsDir(),
		"size":          fileInfo.Size(),
		"mod_time":      fileInfo.ModTime(),
		"is_readable":   true, // We were able to stat it
		"is_writable":   false,
		"is_executable": fileInfo.Mode().Perm()&0111 != 0,
	}

	// Check write permissions by attempting to open for writing
	if !fileInfo.IsDir() {
		// #nosec G304 -- path is validated by caller
		if file, err := os.OpenFile(path, os.O_WRONLY, 0); err == nil {
			permissions["is_writable"] = true
			_ = file.Close()
		}
	} else {
		// For directories, check if we can create a temp file
		tempPath := filepath.Join(path, ".temp_perm_check")
		// #nosec G304 -- tempPath is constructed safely
		if file, err := os.Create(tempPath); err == nil {
			permissions["is_writable"] = true
			_ = file.Close()
			_ = os.Remove(tempPath)
		}
	}

	// Get system-specific information if available
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		permissions["uid"] = stat.Uid
		permissions["gid"] = stat.Gid
	}

	return permissions, nil
}

// ValidateFileIntegrity validates file integrity using checksums
func (sfo *SecureFileOperations) ValidateFileIntegrity(path string, expectedChecksum string) (*FileOperationResult, error) {
	result := &FileOperationResult{
		FilePath:  path,
		Operation: "integrity_check",
		Timestamp: time.Now(),
	}

	// Read file and calculate checksum
	fileResult, _, err := sfo.SecureReadFile(path)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file for integrity check: %v", err)
		return result, err
	}

	result.BytesRead = fileResult.BytesRead
	result.Checksum = fileResult.Checksum

	// Compare checksums
	if result.Checksum != expectedChecksum {
		result.Error = fmt.Sprintf("integrity check failed: expected %s, got %s", expectedChecksum, result.Checksum)
		result.Warnings = append(result.Warnings, "File may have been modified or corrupted")
		return result, fmt.Errorf("file integrity check failed")
	}

	result.Success = true
	return result, nil
}

// SetSecurityConfig updates security configuration
func (sfo *SecureFileOperations) SetSecurityConfig(config map[string]interface{}) error {
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
