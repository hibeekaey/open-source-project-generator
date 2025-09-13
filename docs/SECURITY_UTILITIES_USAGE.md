# Security Utilities Usage Guide

This guide provides detailed instructions for using the security utilities in `pkg/security` and `pkg/models`. These utilities are designed to replace insecure patterns throughout the codebase and provide secure defaults for common operations.

## Table of Contents

1. [Quick Start](#quick-start)
2. [SecureRandom Interface](#securerandom-interface)
3. [SecureFileOperations Interface](#securefileoperations-interface)
4. [Security Configuration](#security-configuration)
5. [Integration Examples](#integration-examples)
6. [Migration Guide](#migration-guide)
7. [Troubleshooting](#troubleshooting)

## Quick Start

### Basic Usage

```go
import "your-project/pkg/security"

// Generate secure random data
randomID, err := security.GenerateSecureID("prefix")
randomBytes, err := security.GenerateBytes(32)
hexString, err := security.GenerateHexString(16)

// Secure file operations
err = security.WriteFileAtomic("config.json", data, 0600)
tempFile, err := security.CreateSecureTempFile("/tmp", "myapp-")
err = security.ValidatePath("/user/input/path", []string{"/allowed/dir"})
```

### Configuration

```go
import "your-project/pkg/models"

// Use secure defaults
securityConfig := models.DefaultSecurityConfig()
randomConfig := models.DefaultRandomConfig()

// Validate configuration
validator := models.NewSecurityValidator()
result := validator.ValidateCombinedConfig(securityConfig, randomConfig)
```

## SecureRandom Interface

The `SecureRandom` interface provides cryptographically secure random generation to replace predictable sources like timestamps and `math/rand`.

### Interface Definition

```go
type SecureRandom interface {
    GenerateRandomSuffix(length int) (string, error)
    GenerateSecureID(prefix string) (string, error)
    GenerateBytes(length int) ([]byte, error)
    GenerateHexString(length int) (string, error)
    GenerateBase64String(length int) (string, error)
    GenerateAlphanumeric(length int) (string, error)
}
```

### Creating SecureRandom Instances

```go
// Default configuration (16-character hex format)
secureRand := security.NewSecureRandom()

// Custom configuration
secureRand := security.NewSecureRandomWithConfig(
    32,        // suffix length
    "base64",  // format: "hex", "base64", or "alphanumeric"
)
```

### Method Usage

#### GenerateBytes - Core Random Generation

```go
// Generate cryptographically secure random bytes
randomBytes, err := secureRand.GenerateBytes(32)
if err != nil {
    return fmt.Errorf("entropy failure: %w", err)
}

// Use cases:
// - Cryptographic keys and IVs
// - Session tokens
// - CSRF tokens
// - Salt for password hashing
```

#### GenerateSecureID - Audit Trail IDs

```go
// Generate secure IDs for audit trails
auditID, err := secureRand.GenerateSecureID("audit")
// Result: "audit_a1b2c3d4e5f6g7h8"

sessionID, err := secureRand.GenerateSecureID("session")
// Result: "session_9i8j7k6l5m4n3o2p"

// Empty prefix for just the random part
randomID, err := secureRand.GenerateSecureID("")
// Result: "a1b2c3d4e5f6g7h8"
```

#### GenerateRandomSuffix - Temporary Files

```go
// Generate secure suffixes for temporary files
suffix, err := secureRand.GenerateRandomSuffix(16)
if err != nil {
    return fmt.Errorf("failed to generate temp file suffix: %w", err)
}

tempFileName := fmt.Sprintf("myapp_%s.tmp", suffix)
// Result: "myapp_a1b2c3d4e5f6g7h8.tmp"
```

#### Format-Specific Generation

```go
// Hex format (most compact)
hexString, err := secureRand.GenerateHexString(16)
// Result: "a1b2c3d4e5f6g7h8" (16 characters, 64 bits entropy)

// Base64 format (URL-safe)
base64String, err := secureRand.GenerateBase64String(16)
// Result: "Zm9vYmFyYmF6cXV4" (16 characters, ~96 bits entropy)

// Alphanumeric format (human-readable)
alphanumeric, err := secureRand.GenerateAlphanumeric(16)
// Result: "Kj8Nm2Qr5Xt9Wv3Z" (16 characters, ~95 bits entropy)
```

### Global Convenience Functions

For simple use cases, use the global convenience functions:

```go
// These use the default global instance
randomID, err := security.GenerateSecureID("prefix")
randomBytes, err := security.GenerateBytes(32)
hexString, err := security.GenerateHexString(16)
base64String, err := security.GenerateBase64String(16)
alphanumeric, err := security.GenerateAlphanumeric(16)
```

### Error Handling

Always handle errors from random generation functions:

```go
randomID, err := security.GenerateSecureID("audit")
if err != nil {
    // Entropy failure is a serious security issue
    log.Printf("SECURITY ALERT: Random generation failed: %v", err)
    return fmt.Errorf("security operation failed")
}
```

## SecureFileOperations Interface

The `SecureFileOperations` interface provides secure file operations to prevent race conditions, directory traversal, and other file-based attacks.

### Interface Definition

```go
type SecureFileOperations interface {
    WriteFileAtomic(filename string, data []byte, perm os.FileMode) error
    CreateSecureTempFile(dir, pattern string) (*os.File, error)
    ValidatePath(path string, allowedDirs []string) error
    SecureDelete(filename string) error
    EnsureSecurePermissions(path string, perm os.FileMode) error
}
```

### Creating SecureFileOperations Instances

```go
// Default configuration
secureFileOps := security.NewSecureFileOperations()

// Custom configuration
secureFileOps := security.NewSecureFileOperationsWithConfig(
    secureRandom,        // SecureRandom instance
    16,                  // temp file random length
    []string{"/tmp"},    // allowed temp directories
    true,                // enable path validation
)
```

### Method Usage

#### WriteFileAtomic - Atomic File Writes

```go
// Write data atomically with secure temporary file
data := []byte(`{"config": "value"}`)
err := secureFileOps.WriteFileAtomic("config.json", data, 0600)
if err != nil {
    return fmt.Errorf("failed to write config: %w", err)
}

// The operation is atomic - either succeeds completely or fails completely
// Temporary files use cryptographically secure random names
// Proper cleanup is guaranteed even on errors
```

#### CreateSecureTempFile - Secure Temporary Files

```go
// Create secure temporary file
tempFile, err := secureFileOps.CreateSecureTempFile("/tmp", "myapp-")
if err != nil {
    return fmt.Errorf("failed to create temp file: %w", err)
}

tempPath := tempFile.Name()
// Result: "/tmp/myapp-a1b2c3d4e5f6g7h8"

// Always clean up temporary files
defer func() {
    tempFile.Close()
    secureFileOps.SecureDelete(tempPath)
}()

// Use the temporary file...
```

#### ValidatePath - Path Validation

```go
// Validate user-provided paths
userPath := "/user/input/../../etc/passwd" // Malicious input
allowedDirs := []string{"/app/data", "/app/uploads"}

err := secureFileOps.ValidatePath(userPath, allowedDirs)
if err != nil {
    // Path validation failed - potential directory traversal
    log.Printf("Security violation: invalid path access attempt")
    return fmt.Errorf("access denied")
}

// Path is safe to use
data, err := os.ReadFile(userPath)
```

#### SecureDelete - Secure File Deletion

```go
// Securely delete sensitive files
err := secureFileOps.SecureDelete("sensitive-data.json")
if err != nil {
    return fmt.Errorf("failed to securely delete file: %w", err)
}

// File is overwritten with random data before deletion
// Prevents recovery of sensitive information
```

#### EnsureSecurePermissions - Permission Management

```go
// Set secure permissions on files/directories
err := secureFileOps.EnsureSecurePermissions("config.json", 0600)
if err != nil {
    return fmt.Errorf("failed to set secure permissions: %w", err)
}

// Common secure permission patterns:
// 0600 - Owner read/write only (sensitive files)
// 0644 - Owner read/write, group/other read (public files)
// 0700 - Owner full access only (sensitive directories)
// 0755 - Owner full access, group/other read/execute (public directories)
```

### Global Convenience Functions

```go
// These use the default global instance
err := security.WriteFileAtomic("file.json", data, 0600)
tempFile, err := security.CreateSecureTempFile("/tmp", "prefix-")
err := security.ValidatePath(path, allowedDirs)
err := security.SecureDelete("sensitive.txt")
err := security.EnsureSecurePermissions("file.txt", 0600)
```

## Security Configuration

### SecurityConfig Structure

```go
type SecurityConfig struct {
    TempFileRandomLength int         `yaml:"temp_file_random_length"`
    AllowedTempDirs      []string    `yaml:"allowed_temp_dirs"`
    FilePermissions      os.FileMode `yaml:"file_permissions"`
    EnablePathValidation bool        `yaml:"enable_path_validation"`
    MaxFileSize          int64       `yaml:"max_file_size"`
    SecureCleanup        bool        `yaml:"secure_cleanup"`
}
```

### RandomConfig Structure

```go
type RandomConfig struct {
    DefaultSuffixLength int    `yaml:"default_suffix_length"`
    IDFormat            string `yaml:"id_format"`
    MinEntropyBytes     int    `yaml:"min_entropy_bytes"`
    IDPrefixLength      int    `yaml:"id_prefix_length"`
    EnableEntropyCheck  bool   `yaml:"enable_entropy_check"`
}
```

### Configuration Usage

```go
// Create custom security configuration
securityConfig := &models.SecurityConfig{
    TempFileRandomLength: 32,                    // Longer random suffixes
    AllowedTempDirs:      []string{"/app/tmp"},  // Restricted temp directories
    FilePermissions:      0600,                  // Owner-only permissions
    EnablePathValidation: true,                  // Enable path validation
    MaxFileSize:          50 * 1024 * 1024,     // 50MB max file size
    SecureCleanup:        true,                  // Enable secure deletion
}

randomConfig := &models.RandomConfig{
    DefaultSuffixLength: 32,        // Longer random suffixes
    IDFormat:            "base64",  // Use base64 format
    MinEntropyBytes:     64,        // Higher entropy requirement
    IDPrefixLength:      8,         // Longer prefixes
    EnableEntropyCheck:  true,      // Enable entropy validation
}

// Validate configuration
validator := models.NewSecurityValidator()
result := validator.ValidateCombinedConfig(securityConfig, randomConfig)

if !result.Valid {
    for _, err := range result.Errors {
        log.Printf("Config error: %s", err.Message)
    }
    return fmt.Errorf("invalid security configuration")
}

// Handle warnings
for _, warning := range result.Warnings {
    log.Printf("Config warning: %s (impact: %s)", warning.Message, warning.Impact)
}
```

### Using Secure Defaults

```go
// Recommended: Use secure defaults
securityConfig := models.DefaultSecurityConfig()
randomConfig := models.DefaultRandomConfig()

// Customize only what you need
securityConfig.AllowedTempDirs = []string{"/app/tmp", "/var/app/tmp"}
randomConfig.IDFormat = "base64"

// Validate the configuration
validator := models.NewSecurityValidator()
result := validator.ValidateCombinedConfig(securityConfig, randomConfig)
```

## Integration Examples

### Replacing Insecure Storage Pattern

**Before (Insecure - from pkg/version/storage.go):**

```go
func (fs *FileStorage) Save(store *models.VersionStore) error {
    data, err := json.Marshal(store)
    if err != nil {
        return err
    }
    
    // INSECURE: Predictable temporary file name
    tempFile := fmt.Sprintf("%s.tmp.%d", fs.filePath, time.Now().UnixNano())
    
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return err
    }
    
    return os.Rename(tempFile, fs.filePath)
}
```

**After (Secure):**

```go
func (fs *FileStorage) Save(store *models.VersionStore) error {
    data, err := json.Marshal(store)
    if err != nil {
        return err
    }
    
    // SECURE: Atomic write with secure temporary file
    return security.WriteFileAtomic(fs.filePath, data, 0600)
}
```

### Replacing Insecure Audit ID Generation

**Before (Insecure - from pkg/reporting/audit.go):**

```go
func generateEventID() string {
    // INSECURE: Predictable timestamp-based ID
    return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}
```

**After (Secure):**

```go
func generateEventID() (string, error) {
    // SECURE: Cryptographically secure audit ID
    return security.GenerateSecureID("audit")
}
```

### Secure Configuration Management

```go
type ConfigManager struct {
    securityConfig *models.SecurityConfig
    randomConfig   *models.RandomConfig
    fileOps        security.SecureFileOperations
}

func NewConfigManager() (*ConfigManager, error) {
    // Load secure configuration
    securityConfig := models.DefaultSecurityConfig()
    randomConfig := models.DefaultRandomConfig()
    
    // Validate configuration
    validator := models.NewSecurityValidator()
    result := validator.ValidateCombinedConfig(securityConfig, randomConfig)
    if !result.Valid {
        return nil, fmt.Errorf("invalid security configuration")
    }
    
    // Create secure file operations with validated config
    fileOps := security.NewSecureFileOperationsWithConfig(
        security.NewSecureRandomWithConfig(
            randomConfig.DefaultSuffixLength,
            randomConfig.IDFormat,
        ),
        securityConfig.TempFileRandomLength,
        securityConfig.AllowedTempDirs,
        securityConfig.EnablePathValidation,
    )
    
    return &ConfigManager{
        securityConfig: securityConfig,
        randomConfig:   randomConfig,
        fileOps:        fileOps,
    }, nil
}

func (cm *ConfigManager) SaveConfig(config interface{}, path string) error {
    // Serialize configuration
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // Write atomically with secure permissions
    return cm.fileOps.WriteFileAtomic(path, data, cm.securityConfig.FilePermissions)
}
```

## Migration Guide

### Step 1: Identify Insecure Patterns

Search for these patterns in your codebase:

```bash
# Find timestamp-based random generation
grep -r "time\.Now()\.UnixNano()" .
grep -r "time\.Now()\.Unix()" .

# Find math/rand usage
grep -r "math/rand" .
grep -r "rand\.Int" .

# Find non-atomic file operations
grep -r "os\.WriteFile" .
grep -r "ioutil\.WriteFile" .

# Find predictable temporary files
grep -r "\.tmp\." .
grep -r "/tmp/" .
```

### Step 2: Replace Random Generation

```go
// Replace this pattern:
id := fmt.Sprintf("prefix_%d", time.Now().UnixNano())

// With this:
id, err := security.GenerateSecureID("prefix")
if err != nil {
    return fmt.Errorf("failed to generate secure ID: %w", err)
}
```

### Step 3: Replace File Operations

```go
// Replace this pattern:
tempFile := fmt.Sprintf("%s.tmp.%d", filename, time.Now().UnixNano())
os.WriteFile(tempFile, data, 0644)
os.Rename(tempFile, filename)

// With this:
err := security.WriteFileAtomic(filename, data, 0600)
if err != nil {
    return fmt.Errorf("failed to write file: %w", err)
}
```

### Step 4: Add Path Validation

```go
// Add validation for user-provided paths:
func handleUserFile(userPath string) error {
    allowedDirs := []string{"/app/data", "/app/uploads"}
    if err := security.ValidatePath(userPath, allowedDirs); err != nil {
        return fmt.Errorf("invalid file path")
    }
    
    // Proceed with file operation...
}
```

### Step 5: Update Error Handling

```go
// Replace detailed error messages:
return fmt.Errorf("failed to read file %s: %v", path, err)

// With generic messages and secure logging:
log.Printf("File operation failed for %s: %v", path, err)
return fmt.Errorf("file operation failed")
```

## Troubleshooting

### Common Issues

#### Entropy Failures

```go
// Problem: Random generation fails in low-entropy environments
randomID, err := security.GenerateSecureID("prefix")
if err != nil {
    // This indicates insufficient system entropy
    log.Printf("CRITICAL: Entropy failure - %v", err)
    
    // Don't fall back to insecure alternatives!
    // Instead, investigate the entropy source
    return fmt.Errorf("security operation failed")
}
```

**Solutions:**

- Check `/proc/sys/kernel/random/entropy_avail` on Linux
- Install `haveged` or `rng-tools` to improve entropy
- In containers, ensure host entropy is available
- Never fall back to `math/rand` or timestamps

#### Path Validation Failures

```go
// Problem: Legitimate paths being rejected
err := security.ValidatePath("/app/data/user/file.txt", []string{"/app/data"})
if err != nil {
    // Check if the path is actually within allowed directories
    log.Printf("Path validation failed: %v", err)
}
```

**Solutions:**

- Ensure allowed directories are absolute paths
- Check for symlinks that might escape allowed directories
- Use `filepath.Clean()` on input paths before validation
- Verify directory permissions and existence

#### Permission Issues

```go
// Problem: Cannot write files with secure permissions
err := security.WriteFileAtomic("config.json", data, 0600)
if err != nil {
    // Check if the directory is writable
    // Check if the file already exists with different ownership
}
```

**Solutions:**

- Ensure the application has write permissions to the directory
- Check file ownership and group membership
- Verify umask settings don't conflict with desired permissions
- Use `security.EnsureSecurePermissions()` to fix existing files

#### Configuration Validation Errors

```go
validator := models.NewSecurityValidator()
result := validator.ValidateSecurityConfig(config)

if !result.Valid {
    for _, err := range result.Errors {
        log.Printf("Config error: %s (field: %s, value: %v)", 
            err.Message, err.Field, err.Value)
    }
}
```

**Solutions:**

- Check minimum/maximum values for numeric fields
- Ensure directory paths exist and are accessible
- Verify enum values (like IDFormat) are valid
- Review security warnings and adjust configuration accordingly

### Performance Considerations

#### Random Generation Performance

```go
// For high-frequency operations, reuse SecureRandom instances
secureRand := security.NewSecureRandom()

// Generate many IDs efficiently
for i := 0; i < 1000; i++ {
    id, err := secureRand.GenerateSecureID("batch")
    // Process ID...
}
```

#### File Operation Performance

```go
// For many small files, consider batching
var batch []FileData
for _, file := range files {
    batch = append(batch, file)
    if len(batch) >= 100 {
        err := writeBatchAtomic(batch)
        batch = batch[:0]
    }
}
```

### Debugging

Enable detailed logging for security operations:

```go
import "log"

// Set up detailed logging
log.SetFlags(log.LstdFlags | log.Lshortfile)

// Log security operations
func debugSecurityOperation(operation string, details interface{}) {
    if os.Getenv("DEBUG_SECURITY") == "true" {
        log.Printf("SECURITY DEBUG: %s - %+v", operation, details)
    }
}
```

This comprehensive usage guide should help developers properly implement and migrate to the secure utilities while avoiding common pitfalls and security vulnerabilities.
