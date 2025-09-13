# Security Coding Guide

This guide provides comprehensive documentation for secure coding patterns implemented in the core application, focusing on cryptographically secure operations, safe file handling, and security best practices.

## Table of Contents

1. [Overview](#overview)
2. [Secure Random Generation](#secure-random-generation)
3. [Secure File Operations](#secure-file-operations)
4. [Security Configuration](#security-configuration)
5. [Common Security Patterns](#common-security-patterns)
6. [Security Anti-Patterns](#security-anti-patterns)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

## Overview

The security utilities in this application are designed to address common security vulnerabilities in file operations and random number generation. The primary focus is on:

- **Cryptographically Secure Random Generation**: Using `crypto/rand` instead of predictable sources
- **Atomic File Operations**: Preventing race conditions and partial writes
- **Path Validation**: Preventing directory traversal attacks
- **Secure Temporary Files**: Using unpredictable names and proper permissions

### Key Security Principles

1. **Fail Securely**: When security operations fail, the system fails in a secure state
2. **Defense in Depth**: Multiple layers of security validation and protection
3. **Least Privilege**: Minimal permissions and access rights
4. **Secure by Default**: Safe defaults that don't require security expertise to use correctly

## Secure Random Generation

### Why Secure Random Generation Matters

**CRITICAL**: Never use `math/rand` or timestamp-based generation for security-sensitive operations. These sources are predictable and can be exploited by attackers.

```go
// ❌ INSECURE - Predictable and exploitable
func generateInsecureID() string {
    return fmt.Sprintf("id_%d", time.Now().UnixNano())
}

// ❌ INSECURE - math/rand is predictable
func generateWeakRandom() string {
    return fmt.Sprintf("%d", rand.Int63())
}

// ✅ SECURE - Uses cryptographically secure random generation
func generateSecureID() (string, error) {
    return security.GenerateSecureID("id")
}
```

### Using the SecureRandom Interface

The `SecureRandom` interface provides cryptographically secure random generation:

```go
import "your-project/pkg/security"

// Create a secure random generator
secureRand := security.NewSecureRandom()

// Generate different types of random data
hexString, err := secureRand.GenerateHexString(16)
base64String, err := secureRand.GenerateBase64String(16)
alphanumeric, err := secureRand.GenerateAlphanumeric(16)
randomBytes, err := secureRand.GenerateBytes(32)

// Generate secure IDs for audit trails
auditID, err := secureRand.GenerateSecureID("audit")
```

### Configuration Options

```go
// Custom configuration for different security requirements
secureRand := security.NewSecureRandomWithConfig(
    32,        // suffix length - longer is more secure
    "base64",  // format: "hex", "base64", or "alphanumeric"
)
```

### Security Considerations

- **Entropy Quality**: All functions use `crypto/rand` which provides cryptographically secure randomness
- **Error Handling**: Always check for errors - entropy failures can occur in low-entropy environments
- **Length Requirements**: Minimum 8 characters for basic security, 16+ recommended for high security
- **Format Selection**:
  - `hex`: Most compact, good for IDs
  - `base64`: URL-safe, good for tokens
  - `alphanumeric`: Human-readable, lower entropy density

## Secure File Operations

### Why Secure File Operations Matter

File operations are a common source of security vulnerabilities:

- **Race Conditions**: Predictable temporary file names allow attackers to create files before your application
- **Directory Traversal**: Unsanitized paths can allow access to unauthorized files
- **Information Disclosure**: Improper permissions can expose sensitive data
- **Partial Writes**: Non-atomic operations can leave files in inconsistent states

### Using SecureFileOperations

```go
import "your-project/pkg/security"

// Create secure file operations handler
secureFileOps := security.NewSecureFileOperations()

// Write files atomically with secure temporary files
data := []byte("sensitive configuration data")
err := secureFileOps.WriteFileAtomic("/path/to/config.json", data, 0600)

// Create secure temporary files
tempFile, err := secureFileOps.CreateSecureTempFile("/tmp", "myapp-")
defer tempFile.Close()

// Validate paths before use
err := secureFileOps.ValidatePath("/user/input/path", []string{"/allowed/dir1", "/allowed/dir2"})

// Securely delete sensitive files
err := secureFileOps.SecureDelete("/path/to/sensitive/file")
```

### Atomic File Operations

Atomic file operations prevent race conditions and ensure data consistency:

```go
// ❌ INSECURE - Non-atomic write, vulnerable to race conditions
func writeFileInsecure(filename string, data []byte) error {
    return os.WriteFile(filename, data, 0644)
}

// ✅ SECURE - Atomic write using secure temporary file
func writeFileSecure(filename string, data []byte) error {
    return security.WriteFileAtomic(filename, data, 0600)
}
```

### Path Validation

Always validate file paths to prevent directory traversal attacks:

```go
// ❌ INSECURE - No path validation
func readUserFile(userPath string) ([]byte, error) {
    return os.ReadFile(userPath) // Vulnerable to ../../../etc/passwd
}

// ✅ SECURE - Path validation prevents traversal
func readUserFileSecure(userPath string) ([]byte, error) {
    allowedDirs := []string{"/app/user-data", "/app/uploads"}
    if err := security.ValidatePath(userPath, allowedDirs); err != nil {
        return nil, fmt.Errorf("invalid path: %w", err)
    }
    return os.ReadFile(userPath)
}
```

### Secure Temporary Files

Temporary files must use unpredictable names and secure permissions:

```go
// ❌ INSECURE - Predictable temporary file name
func createTempFileInsecure() (*os.File, error) {
    filename := fmt.Sprintf("/tmp/myapp_%d", time.Now().UnixNano())
    return os.Create(filename) // Predictable name, wrong permissions
}

// ✅ SECURE - Cryptographically secure temporary file
func createTempFileSecure() (*os.File, error) {
    return security.CreateSecureTempFile("/tmp", "myapp-")
}
```

## Security Configuration

### Configuration Models

The security configuration system provides validation and secure defaults:

```go
import "your-project/pkg/models"

// Create security configuration with validation
securityConfig := &models.SecurityConfig{
    TempFileRandomLength: 16,
    AllowedTempDirs:      []string{"/tmp", "/var/tmp"},
    FilePermissions:      0600, // Owner read/write only
    EnablePathValidation: true,
    MaxFileSize:          10 * 1024 * 1024, // 10MB
    SecureCleanup:        true,
}

// Validate configuration
validator := models.NewSecurityValidator()
result := validator.ValidateSecurityConfig(securityConfig)

if !result.Valid {
    for _, err := range result.Errors {
        log.Printf("Security config error: %s", err.Message)
    }
}

for _, warning := range result.Warnings {
    log.Printf("Security config warning: %s", warning.Message)
}
```

### Random Generation Configuration

```go
randomConfig := &models.RandomConfig{
    DefaultSuffixLength: 16,
    IDFormat:            "hex",
    MinEntropyBytes:     32,
    IDPrefixLength:      4,
    EnableEntropyCheck:  true,
}

result := validator.ValidateRandomConfig(randomConfig)
```

### Using Secure Defaults

```go
// Use secure defaults - recommended for most applications
securityConfig := models.DefaultSecurityConfig()
randomConfig := models.DefaultRandomConfig()

// Validate combined configuration
result := validator.ValidateCombinedConfig(securityConfig, randomConfig)
```

## Common Security Patterns

### Pattern 1: Secure ID Generation for Audit Trails

```go
// ✅ SECURE - Cryptographically secure audit ID
func createAuditEvent(eventType string, userID string) (*AuditEvent, error) {
    eventID, err := security.GenerateSecureID("audit")
    if err != nil {
        return nil, fmt.Errorf("failed to generate secure audit ID: %w", err)
    }
    
    return &AuditEvent{
        ID:        eventID,
        Type:      eventType,
        UserID:    userID,
        Timestamp: time.Now().UTC(),
    }, nil
}
```

### Pattern 2: Secure Configuration File Updates

```go
// ✅ SECURE - Atomic configuration update with validation
func updateConfiguration(configPath string, newConfig *Config) error {
    // Validate the path
    allowedDirs := []string{"/app/config", "/etc/myapp"}
    if err := security.ValidatePath(configPath, allowedDirs); err != nil {
        return fmt.Errorf("invalid config path: %w", err)
    }
    
    // Serialize configuration
    data, err := json.Marshal(newConfig)
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // Write atomically with secure permissions
    return security.WriteFileAtomic(configPath, data, 0600)
}
```

### Pattern 3: Secure Session Token Generation

```go
// ✅ SECURE - Cryptographically secure session tokens
func generateSessionToken() (string, error) {
    // Generate 32 bytes of random data (256 bits)
    tokenBytes, err := security.GenerateBytes(32)
    if err != nil {
        return "", fmt.Errorf("failed to generate session token: %w", err)
    }
    
    // Encode as URL-safe base64
    return base64.URLEncoding.EncodeToString(tokenBytes), nil
}
```

### Pattern 4: Secure Temporary File Processing

```go
// ✅ SECURE - Complete secure temporary file workflow
func processUploadedFile(uploadData []byte) error {
    // Create secure temporary file
    tempFile, err := security.CreateSecureTempFile("", "upload-")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %w", err)
    }
    
    tempPath := tempFile.Name()
    
    // Ensure cleanup
    defer func() {
        tempFile.Close()
        security.SecureDelete(tempPath) // Secure deletion
    }()
    
    // Write data to temporary file
    if _, err := tempFile.Write(uploadData); err != nil {
        return fmt.Errorf("failed to write to temp file: %w", err)
    }
    
    // Process the file...
    return processFile(tempPath)
}
```

## Security Anti-Patterns

### Anti-Pattern 1: Predictable Random Generation

```go
// ❌ NEVER DO THIS - Predictable and exploitable
func generateInsecureToken() string {
    // Timestamp-based generation is predictable
    return fmt.Sprintf("token_%d", time.Now().UnixNano())
}

// ❌ NEVER DO THIS - math/rand is not cryptographically secure
func generateWeakToken() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("token_%d", rand.Int63())
}
```

### Anti-Pattern 2: Insecure File Operations

```go
// ❌ NEVER DO THIS - Race condition vulnerability
func saveConfigInsecure(config *Config) error {
    tempFile := fmt.Sprintf("/tmp/config_%d", time.Now().UnixNano())
    
    // Attacker can predict filename and create file first
    data, _ := json.Marshal(config)
    os.WriteFile(tempFile, data, 0644) // Wrong permissions too
    
    return os.Rename(tempFile, "/app/config.json")
}

// ❌ NEVER DO THIS - Directory traversal vulnerability
func readUserFileInsecure(filename string) ([]byte, error) {
    // No validation - vulnerable to ../../../etc/passwd
    return os.ReadFile(filename)
}
```

### Anti-Pattern 3: Information Disclosure in Errors

```go
// ❌ NEVER DO THIS - Leaks sensitive path information
func validateFileInsecure(path string) error {
    if strings.Contains(path, "..") {
        return fmt.Errorf("invalid path: %s contains directory traversal", path)
    }
    return nil
}

// ✅ DO THIS - Generic error message
func validateFileSecure(path string) error {
    if strings.Contains(path, "..") {
        return errors.New("invalid path: directory traversal detected")
    }
    return nil
}
```

## Best Practices

### 1. Always Use Cryptographically Secure Random Generation

- Use `crypto/rand` for all security-sensitive operations
- Never use `math/rand` or timestamp-based generation for security
- Always handle entropy errors appropriately

### 2. Implement Atomic File Operations

- Use secure temporary files with unpredictable names
- Perform atomic rename operations
- Set appropriate file permissions (0600 for sensitive files)

### 3. Validate All File Paths

- Sanitize user-provided paths
- Use allowlists for permitted directories
- Check for directory traversal patterns

### 4. Handle Errors Securely

- Don't leak sensitive information in error messages
- Log security events for audit purposes
- Fail securely when operations cannot be completed safely

### 5. Use Secure Defaults

- Enable path validation by default
- Use restrictive file permissions
- Enable secure cleanup of temporary files

### 6. Regular Security Audits

- Scan code for insecure patterns
- Validate security configurations
- Test for race conditions and edge cases

## Examples

### Complete Secure File Processing Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "your-project/pkg/security"
    "your-project/pkg/models"
)

type UserData struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func secureUserDataProcessing(userData *UserData, outputPath string) error {
    // 1. Generate secure ID for the user
    userID, err := security.GenerateSecureID("user")
    if err != nil {
        return fmt.Errorf("failed to generate secure user ID: %w", err)
    }
    userData.ID = userID
    
    // 2. Validate output path
    allowedDirs := []string{"/app/data/users", "/var/app/users"}
    if err := security.ValidatePath(outputPath, allowedDirs); err != nil {
        return fmt.Errorf("invalid output path: %w", err)
    }
    
    // 3. Serialize user data
    jsonData, err := json.Marshal(userData)
    if err != nil {
        return fmt.Errorf("failed to serialize user data: %w", err)
    }
    
    // 4. Write atomically with secure permissions
    if err := security.WriteFileAtomic(outputPath, jsonData, 0600); err != nil {
        return fmt.Errorf("failed to write user data: %w", err)
    }
    
    // 5. Generate audit event
    auditID, err := security.GenerateSecureID("audit")
    if err != nil {
        log.Printf("Warning: failed to generate audit ID: %v", err)
    } else {
        log.Printf("User data processed successfully - Audit ID: %s", auditID)
    }
    
    return nil
}

func main() {
    // Initialize with secure configuration
    securityConfig := models.DefaultSecurityConfig()
    randomConfig := models.DefaultRandomConfig()
    
    // Validate configuration
    validator := models.NewSecurityValidator()
    result := validator.ValidateCombinedConfig(securityConfig, randomConfig)
    
    if !result.Valid {
        log.Fatal("Invalid security configuration")
    }
    
    // Process user data securely
    userData := &UserData{
        Username: "john_doe",
        Email:    "john@example.com",
    }
    
    if err := secureUserDataProcessing(userData, "/app/data/users/john_doe.json"); err != nil {
        log.Fatalf("Failed to process user data: %v", err)
    }
    
    fmt.Println("User data processed successfully with secure patterns")
}
```

### Secure Configuration Management Example

```go
package config

import (
    "encoding/json"
    "fmt"
    "your-project/pkg/security"
    "your-project/pkg/models"
)

type AppConfig struct {
    DatabaseURL string `json:"database_url"`
    APIKey      string `json:"api_key"`
    Debug       bool   `json:"debug"`
}

type ConfigManager struct {
    secureFileOps security.SecureFileOperations
    configPath    string
    allowedDirs   []string
}

func NewConfigManager(configPath string) (*ConfigManager, error) {
    // Validate configuration path
    allowedDirs := []string{"/app/config", "/etc/myapp"}
    if err := security.ValidatePath(configPath, allowedDirs); err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }
    
    return &ConfigManager{
        secureFileOps: security.NewSecureFileOperations(),
        configPath:    configPath,
        allowedDirs:   allowedDirs,
    }, nil
}

func (cm *ConfigManager) LoadConfig() (*AppConfig, error) {
    // Validate path before reading
    if err := cm.secureFileOps.ValidatePath(cm.configPath, cm.allowedDirs); err != nil {
        return nil, fmt.Errorf("config path validation failed: %w", err)
    }
    
    // Read configuration file
    data, err := os.ReadFile(cm.configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config AppConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &config, nil
}

func (cm *ConfigManager) SaveConfig(config *AppConfig) error {
    // Generate secure backup filename
    backupSuffix, err := security.GenerateRandomSuffix(8)
    if err != nil {
        return fmt.Errorf("failed to generate backup suffix: %w", err)
    }
    
    backupPath := cm.configPath + ".backup." + backupSuffix
    
    // Serialize configuration
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // Create backup of existing config
    if _, err := os.Stat(cm.configPath); err == nil {
        if err := cm.secureFileOps.WriteFileAtomic(backupPath, data, 0600); err != nil {
            return fmt.Errorf("failed to create config backup: %w", err)
        }
    }
    
    // Write new configuration atomically
    if err := cm.secureFileOps.WriteFileAtomic(cm.configPath, data, 0600); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }
    
    return nil
}

func (cm *ConfigManager) RotateAPIKey(config *AppConfig) error {
    // Generate new secure API key
    newAPIKey, err := security.GenerateBase64String(32)
    if err != nil {
        return fmt.Errorf("failed to generate new API key: %w", err)
    }
    
    config.APIKey = newAPIKey
    
    // Save updated configuration
    return cm.SaveConfig(config)
}
```

This comprehensive security guide provides developers with the knowledge and examples needed to implement secure coding patterns throughout the application. Always refer to this guide when implementing security-sensitive operations.
