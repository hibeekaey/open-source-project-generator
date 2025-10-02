package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// SecureStorage provides secure storage for configuration files
type SecureStorage struct {
	baseDir     string
	encryptor   *ConfigEncryptor
	logger      interfaces.Logger
	auditLogger *ConfigAuditLogger
	permissions fs.FileMode
}

// StorageOptions defines options for secure storage operations
type StorageOptions struct {
	Encrypt     bool        `json:"encrypt"`
	Permissions fs.FileMode `json:"permissions"`
	Backup      bool        `json:"backup"`
	Validate    bool        `json:"validate"`
}

// FileMetadata contains metadata about stored files
type FileMetadata struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	Permissions  string    `json:"permissions"`
	Checksum     string    `json:"checksum"`
	Encrypted    bool      `json:"encrypted"`
	Owner        string    `json:"owner"`
	LastAccessed time.Time `json:"last_accessed"`
}

// NewSecureStorage creates a new secure storage instance
func NewSecureStorage(baseDir string, encryptor *ConfigEncryptor, logger interfaces.Logger) (*SecureStorage, error) {
	// Ensure base directory exists with secure permissions
	if err := os.MkdirAll(baseDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Set secure permissions on the directory
	if err := os.Chmod(baseDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to set directory permissions: %w", err)
	}

	auditLogger := &ConfigAuditLogger{
		logFile: filepath.Join(baseDir, "storage_audit.log"),
		logger:  logger,
	}

	return &SecureStorage{
		baseDir:     baseDir,
		encryptor:   encryptor,
		logger:      logger,
		auditLogger: auditLogger,
		permissions: 0640, // Default secure permissions for files
	}, nil
}

// StoreFile securely stores a file with optional encryption
func (s *SecureStorage) StoreFile(filename string, data []byte, options *StorageOptions) error {
	if options == nil {
		options = &StorageOptions{
			Encrypt:     true,
			Permissions: s.permissions,
			Backup:      true,
			Validate:    true,
		}
	}

	// Validate filename
	if err := s.validateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	filePath := filepath.Join(s.baseDir, filename)

	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Create backup if file exists and backup is requested
	if options.Backup && s.fileExists(filePath) {
		if err := s.createBackup(filePath); err != nil {
			s.logger.WarnWithFields("Failed to create backup", map[string]interface{}{
				"file":  filePath,
				"error": err.Error(),
			})
		}
	}

	// Encrypt data if requested
	var finalData []byte
	var encrypted bool
	if options.Encrypt && s.encryptor != nil {
		encryptedData, err := s.encryptor.EncryptSensitiveFields(string(data))
		if err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}
		if encryptedStr, ok := encryptedData.(string); ok {
			finalData = []byte(encryptedStr)
			encrypted = true
		} else {
			finalData = data
		}
	} else {
		finalData = data
	}

	// Validate data if requested
	if options.Validate {
		if err := s.validateData(finalData); err != nil {
			return fmt.Errorf("data validation failed: %w", err)
		}
	}

	// Write file securely
	if err := s.writeFileSecurely(filePath, finalData, options.Permissions); err != nil {
		s.auditLogger.LogError("store_file", filename, err, map[string]interface{}{
			"size":      len(finalData),
			"encrypted": encrypted,
		})
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Calculate and store checksum
	checksum := s.calculateChecksum(finalData)

	// Log successful storage
	s.auditLogger.LogAction("store_file", filename, map[string]interface{}{
		"size":      len(finalData),
		"encrypted": encrypted,
		"checksum":  checksum,
	})

	return nil
}

// RetrieveFile securely retrieves a file with optional decryption
func (s *SecureStorage) RetrieveFile(filename string, decrypt bool) ([]byte, error) {
	// Validate filename
	if err := s.validateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	filePath := filepath.Join(s.baseDir, filename)

	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Check if file exists
	if !s.fileExists(filePath) {
		return nil, fmt.Errorf("file not found: %s", filename)
	}

	// Read file securely
	data, err := s.readFileSecurely(filePath)
	if err != nil {
		s.auditLogger.LogError("retrieve_file", filename, err, nil)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Decrypt data if requested and encryptor is available
	var finalData []byte
	if decrypt && s.encryptor != nil {
		decryptedData, err := s.encryptor.DecryptSensitiveFields(string(data))
		if err != nil {
			s.logger.WarnWithFields("Failed to decrypt data", map[string]interface{}{
				"file":  filename,
				"error": err.Error(),
			})
			finalData = data // Return original data if decryption fails
		} else {
			if decryptedStr, ok := decryptedData.(string); ok {
				finalData = []byte(decryptedStr)
			} else {
				finalData = data
			}
		}
	} else {
		finalData = data
	}

	// Log successful retrieval
	s.auditLogger.LogAction("retrieve_file", filename, map[string]interface{}{
		"size":      len(finalData),
		"decrypted": decrypt,
	})

	return finalData, nil
}

// DeleteFile securely deletes a file
func (s *SecureStorage) DeleteFile(filename string) error {
	// Validate filename
	if err := s.validateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	filePath := filepath.Join(s.baseDir, filename)

	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Check if file exists
	if !s.fileExists(filePath) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Create backup before deletion
	if err := s.createBackup(filePath); err != nil {
		s.logger.WarnWithFields("Failed to create backup before deletion", map[string]interface{}{
			"file":  filePath,
			"error": err.Error(),
		})
	}

	// Securely delete file (overwrite with random data first)
	if err := s.secureDelete(filePath); err != nil {
		s.auditLogger.LogError("delete_file", filename, err, nil)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Log successful deletion
	s.auditLogger.LogAction("delete_file", filename, nil)

	return nil
}

// ListFiles lists all files in the secure storage
func (s *SecureStorage) ListFiles() ([]*FileMetadata, error) {
	files, err := filepath.Glob(filepath.Join(s.baseDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var metadata []*FileMetadata
	for _, file := range files {
		meta, err := s.getFileMetadata(file)
		if err != nil {
			s.logger.WarnWithFields("Failed to get file metadata", map[string]interface{}{
				"file":  file,
				"error": err.Error(),
			})
			continue
		}
		metadata = append(metadata, meta)
	}

	return metadata, nil
}

// ValidateIntegrity validates the integrity of stored files
func (s *SecureStorage) ValidateIntegrity() (*IntegrityReport, error) {
	report := &IntegrityReport{
		TotalFiles:   0,
		ValidFiles:   0,
		CorruptFiles: []string{},
		MissingFiles: []string{},
		Errors:       []string{},
		CheckedAt:    time.Now(),
	}

	files, err := s.ListFiles()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Failed to list files: %v", err))
		return report, nil
	}

	report.TotalFiles = len(files)

	for _, file := range files {
		if err := s.validateFileIntegrity(file.Path); err != nil {
			report.CorruptFiles = append(report.CorruptFiles, file.Path)
			report.Errors = append(report.Errors, fmt.Sprintf("File %s: %v", file.Path, err))
		} else {
			report.ValidFiles++
		}
	}

	return report, nil
}

// IntegrityReport contains the results of integrity validation
type IntegrityReport struct {
	TotalFiles   int       `json:"total_files"`
	ValidFiles   int       `json:"valid_files"`
	CorruptFiles []string  `json:"corrupt_files"`
	MissingFiles []string  `json:"missing_files"`
	Errors       []string  `json:"errors"`
	CheckedAt    time.Time `json:"checked_at"`
}

// Helper methods

// validateFilename validates that a filename is safe
func (s *SecureStorage) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if len(filename) > 255 {
		return fmt.Errorf("filename too long: %d characters (max: 255)", len(filename))
	}

	// Check for dangerous characters
	dangerousChars := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains dangerous character: %s", char)
		}
	}

	// Check for reserved names (Windows)
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperFilename := strings.ToUpper(filename)
	for _, reserved := range reservedNames {
		if upperFilename == reserved || strings.HasPrefix(upperFilename, reserved+".") {
			return fmt.Errorf("filename uses reserved name: %s", reserved)
		}
	}

	return nil
}

// validateData validates that data is safe to store
func (s *SecureStorage) validateData(data []byte) error {
	// Check size limits
	maxSize := 100 * 1024 * 1024 // 100MB
	if len(data) > maxSize {
		return fmt.Errorf("data too large: %d bytes (max: %d)", len(data), maxSize)
	}

	// Check for null bytes (potential binary data)
	for i, b := range data {
		if b == 0 {
			return fmt.Errorf("null byte found at position %d", i)
		}
	}

	return nil
}

// writeFileSecurely writes a file with secure permissions
func (s *SecureStorage) writeFileSecurely(filePath string, data []byte, permissions fs.FileMode) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to temporary file first
	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, permissions); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Set secure permissions
	if err := os.Chmod(tempFile, permissions); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Atomic move to final location
	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to move file to final location: %w", err)
	}

	return nil
}

// readFileSecurely reads a file with security checks
func (s *SecureStorage) readFileSecurely(filePath string) ([]byte, error) {
	// Check file permissions
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Warn if file has overly permissive permissions
	if info.Mode().Perm() > 0644 {
		s.logger.WarnWithFields("File has permissive permissions", map[string]interface{}{
			"file":        filePath,
			"permissions": info.Mode().Perm().String(),
		})
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// fileExists checks if a file exists
func (s *SecureStorage) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// createBackup creates a backup of a file
func (s *SecureStorage) createBackup(filePath string) error {
	if !s.fileExists(filePath) {
		return nil // No file to backup
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup.%s", filePath, timestamp)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0640); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// secureDelete securely deletes a file by overwriting it with random data
func (s *SecureStorage) secureDelete(filePath string) error {
	// Get file size
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Open file for writing
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open file for overwriting: %w", err)
	}
	defer file.Close()

	// Overwrite with random data (3 passes)
	for pass := 0; pass < 3; pass++ {
		randomData := make([]byte, info.Size())
		if _, err := rand.Read(randomData); err != nil {
			return fmt.Errorf("failed to generate random data: %w", err)
		}

		if _, err := file.WriteAt(randomData, 0); err != nil {
			return fmt.Errorf("failed to overwrite file: %w", err)
		}

		if err := file.Sync(); err != nil {
			return fmt.Errorf("failed to sync file: %w", err)
		}
	}

	file.Close()

	// Finally delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

// calculateChecksum calculates SHA-256 checksum of data
func (s *SecureStorage) calculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// getFileMetadata gets metadata for a file
func (s *SecureStorage) getFileMetadata(filePath string) (*FileMetadata, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Read file to calculate checksum
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for checksum: %w", err)
	}

	checksum := s.calculateChecksum(data)

	// Check if file appears to be encrypted
	encrypted := s.isFileEncrypted(data)

	metadata := &FileMetadata{
		Path:         filePath,
		Size:         info.Size(),
		ModTime:      info.ModTime(),
		Permissions:  info.Mode().Perm().String(),
		Checksum:     checksum,
		Encrypted:    encrypted,
		Owner:        s.getFileOwner(info),
		LastAccessed: time.Now(), // Approximation
	}

	return metadata, nil
}

// isFileEncrypted checks if file data appears to be encrypted
func (s *SecureStorage) isFileEncrypted(data []byte) bool {
	// Simple heuristic: check for encrypted field prefix
	return strings.Contains(string(data), EncryptedFieldPrefix)
}

// getFileOwner gets the owner of a file (simplified)
func (s *SecureStorage) getFileOwner(info os.FileInfo) string {
	// This is a simplified implementation
	// In a real implementation, you would use syscalls to get the actual owner
	return "current_user"
}

// validateFileIntegrity validates the integrity of a single file
func (s *SecureStorage) validateFileIntegrity(filePath string) error {
	// Check if file exists
	if !s.fileExists(filePath) {
		return fmt.Errorf("file does not exist")
	}

	// Check file permissions
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if permissions are too permissive
	if info.Mode().Perm() > 0644 {
		return fmt.Errorf("file has overly permissive permissions: %s", info.Mode().Perm().String())
	}

	// Try to read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Validate data format (basic check)
	if len(data) == 0 {
		return fmt.Errorf("file is empty")
	}

	// Check for null bytes (shouldn't be in text config files)
	for i, b := range data {
		if b == 0 {
			return fmt.Errorf("null byte found at position %d", i)
		}
	}

	return nil
}
