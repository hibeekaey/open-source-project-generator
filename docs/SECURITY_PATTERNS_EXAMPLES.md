# Security Patterns Examples

This document provides concrete examples of secure and insecure coding patterns, demonstrating the proper use of security utilities and common vulnerabilities to avoid.

## Table of Contents

1. [Random Generation Patterns](#random-generation-patterns)
2. [File Operation Patterns](#file-operation-patterns)
3. [Temporary File Patterns](#temporary-file-patterns)
4. [Path Validation Patterns](#path-validation-patterns)
5. [Configuration Management Patterns](#configuration-management-patterns)
6. [Audit Trail Patterns](#audit-trail-patterns)
7. [Error Handling Patterns](#error-handling-patterns)

## Random Generation Patterns

### ❌ Insecure: Timestamp-Based Random Generation

```go
// VULNERABILITY: Predictable random generation using timestamps
// This pattern was found in pkg/version/storage.go and pkg/reporting/audit.go

package main

import (
    "fmt"
    "time"
    "math/rand"
)

// ❌ NEVER DO THIS - Predictable and exploitable
func generateInsecureID() string {
    // Attackers can predict these IDs based on timing
    return fmt.Sprintf("id_%d", time.Now().UnixNano())
}

// ❌ NEVER DO THIS - math/rand is deterministic
func generateWeakToken() string {
    rand.Seed(time.Now().UnixNano()) // Predictable seed
    return fmt.Sprintf("token_%d", rand.Int63())
}

// ❌ NEVER DO THIS - Timestamp-based temporary file naming
func createInsecureTempFile() string {
    // Race condition vulnerability - attackers can predict and pre-create files
    return fmt.Sprintf("/tmp/myapp_%d.tmp", time.Now().UnixNano())
}

// SECURITY IMPACT:
// - Attackers can predict future IDs/filenames
// - Race condition attacks on temporary files
// - Session hijacking with predictable tokens
// - Audit trail manipulation
```

### ✅ Secure: Cryptographically Secure Random Generation

```go
// SOLUTION: Use cryptographically secure random generation

package main

import (
    "fmt"
    "your-project/pkg/security"
)

// ✅ SECURE - Cryptographically secure ID generation
func generateSecureID() (string, error) {
    return security.GenerateSecureID("id")
}

// ✅ SECURE - Cryptographically secure token generation
func generateSecureToken() (string, error) {
    // 32 bytes = 256 bits of entropy
    return security.GenerateBase64String(32)
}

// ✅ SECURE - Cryptographically secure temporary file naming
func createSecureTempFile() (*os.File, error) {
    return security.CreateSecureTempFile("/tmp", "myapp-")
}

// SECURITY BENEFITS:
// - Unpredictable IDs prevent timing attacks
// - Race condition protection for temp files
// - Cryptographically secure tokens prevent hijacking
// - Tamper-resistant audit trails
```

### Real-World Example: Session Token Generation

```go
// ❌ INSECURE PATTERN (commonly found in web applications)
func generateSessionTokenInsecure(userID int64) string {
    // Predictable based on user ID and timestamp
    return fmt.Sprintf("sess_%d_%d", userID, time.Now().Unix())
}

// ✅ SECURE PATTERN
func generateSessionTokenSecure(userID int64) (string, error) {
    // Generate cryptographically secure random token
    tokenBytes, err := security.GenerateBytes(32)
    if err != nil {
        return "", fmt.Errorf("failed to generate session token: %w", err)
    }
    
    // Encode as URL-safe base64
    token := base64.URLEncoding.EncodeToString(tokenBytes)
    
    // Optionally include non-sensitive prefix for debugging
    return fmt.Sprintf("sess_%s", token), nil
}
```

## File Operation Patterns

### ❌ Insecure: Non-Atomic File Operations

```go
// VULNERABILITY: Race conditions and partial writes

package main

import (
    "encoding/json"
    "os"
)

type Config struct {
    DatabaseURL string `json:"database_url"`
    APIKey      string `json:"api_key"`
}

// ❌ NEVER DO THIS - Non-atomic write operation
func saveConfigInsecure(config *Config, filename string) error {
    data, err := json.Marshal(config)
    if err != nil {
        return err
    }
    
    // VULNERABILITY: Direct write can be interrupted
    // - Partial writes leave file in inconsistent state
    // - Concurrent reads may see incomplete data
    // - System crash can corrupt the file
    return os.WriteFile(filename, data, 0644) // Also wrong permissions
}

// ❌ NEVER DO THIS - Predictable temporary file with race condition
func saveConfigWithTempInsecure(config *Config, filename string) error {
    data, err := json.Marshal(config)
    if err != nil {
        return err
    }
    
    // VULNERABILITY: Predictable temp file name
    tempFile := fmt.Sprintf("%s.tmp.%d", filename, time.Now().UnixNano())
    
    // Race condition: attacker can create this file first
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return err
    }
    
    // If rename fails, temp file is left behind
    return os.Rename(tempFile, filename)
}

// SECURITY IMPACT:
// - Data corruption from interrupted writes
// - Information disclosure from temp files
// - Race condition attacks
// - Inconsistent application state
```

### ✅ Secure: Atomic File Operations

```go
// SOLUTION: Use atomic file operations with secure temporary files

package main

import (
    "encoding/json"
    "fmt"
    "your-project/pkg/security"
)

// ✅ SECURE - Atomic write with secure temporary file
func saveConfigSecure(config *Config, filename string) error {
    data, err := json.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // Atomic write with secure temp file and proper permissions
    return security.WriteFileAtomic(filename, data, 0600)
}

// ✅ SECURE - Complete secure configuration management
func saveConfigWithValidation(config *Config, filename string) error {
    // 1. Validate the target path
    allowedDirs := []string{"/app/config", "/etc/myapp"}
    if err := security.ValidatePath(filename, allowedDirs); err != nil {
        return fmt.Errorf("invalid config path: %w", err)
    }
    
    // 2. Serialize configuration
    data, err := json.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // 3. Write atomically with secure permissions
    if err := security.WriteFileAtomic(filename, data, 0600); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }
    
    return nil
}

// SECURITY BENEFITS:
// - Atomic operations prevent corruption
// - Secure temp files prevent race conditions
// - Path validation prevents directory traversal
// - Proper permissions protect sensitive data
```

## Temporary File Patterns

### ❌ Insecure: Predictable Temporary Files

```go
// VULNERABILITY: Predictable temporary file creation (from pkg/version/storage.go)

package main

import (
    "fmt"
    "os"
    "time"
)

// ❌ NEVER DO THIS - Original vulnerable pattern from storage.go
func saveVersionStoreInsecure(data []byte, filePath string) error {
    // VULNERABILITY: Predictable temporary file name
    tempFile := fmt.Sprintf("%s.tmp.%d", filePath, time.Now().UnixNano())
    
    // Race condition: attacker can predict and create this file
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return err
    }
    
    // If this fails, temp file is left behind with sensitive data
    return os.Rename(tempFile, filePath)
}

// ❌ NEVER DO THIS - Using process ID (still predictable)
func createTempWithPID() string {
    // Still predictable - PID can be guessed
    return fmt.Sprintf("/tmp/app_%d_%d", os.Getpid(), time.Now().Unix())
}

// ATTACK SCENARIO:
// 1. Attacker monitors application behavior
// 2. Predicts next temporary file name
// 3. Creates file with malicious content or symlink
// 4. Application writes sensitive data to attacker-controlled file
// 5. Information disclosure or privilege escalation
```

### ✅ Secure: Cryptographically Secure Temporary Files

```go
// SOLUTION: Use cryptographically secure temporary file creation

package main

import (
    "fmt"
    "your-project/pkg/security"
)

// ✅ SECURE - Fixed version of storage.go pattern
func saveVersionStoreSecure(data []byte, filePath string) error {
    // Use atomic write with secure temporary file
    return security.WriteFileAtomic(filePath, data, 0600)
}

// ✅ SECURE - Manual secure temporary file creation
func createSecureTempFile(dir, prefix string) (*os.File, error) {
    return security.CreateSecureTempFile(dir, prefix)
}

// ✅ SECURE - Complete secure temporary file workflow
func processDataWithTempFile(inputData []byte) error {
    // Create secure temporary file
    tempFile, err := security.CreateSecureTempFile("", "processing-")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %w", err)
    }
    
    tempPath := tempFile.Name()
    
    // Ensure cleanup with secure deletion
    defer func() {
        tempFile.Close()
        security.SecureDelete(tempPath)
    }()
    
    // Process data using temporary file
    if _, err := tempFile.Write(inputData); err != nil {
        return fmt.Errorf("failed to write to temp file: %w", err)
    }
    
    // Additional processing...
    return processFile(tempPath)
}

// SECURITY BENEFITS:
// - Unpredictable file names prevent race conditions
// - Secure permissions protect temporary data
// - Automatic cleanup prevents data leakage
// - Atomic operations ensure consistency
```

## Path Validation Patterns

### ❌ Insecure: Missing Path Validation

```go
// VULNERABILITY: Directory traversal attacks

package main

import (
    "os"
    "path/filepath"
)

// ❌ NEVER DO THIS - No path validation
func readUserFileInsecure(userPath string) ([]byte, error) {
    // VULNERABILITY: Directory traversal
    // User can provide: "../../../etc/passwd"
    return os.ReadFile(userPath)
}

// ❌ NEVER DO THIS - Insufficient validation
func readUserFileWeakValidation(userPath string) ([]byte, error) {
    // Weak validation - can be bypassed
    if strings.Contains(userPath, "..") {
        return nil, errors.New("invalid path")
    }
    
    // Still vulnerable to: "/etc/passwd", symlink attacks, etc.
    return os.ReadFile(userPath)
}

// ❌ NEVER DO THIS - Information disclosure in errors
func validatePathInsecure(path string) error {
    if strings.Contains(path, "..") {
        // VULNERABILITY: Leaks the actual path in error message
        return fmt.Errorf("invalid path: %s contains directory traversal", path)
    }
    return nil
}

// ATTACK SCENARIOS:
// - ../../../etc/passwd (read system files)
// - ../../../home/user/.ssh/id_rsa (steal SSH keys)
// - ../../app/config/database.yml (steal credentials)
// - Symlink attacks to escape sandboxes
```

### ✅ Secure: Comprehensive Path Validation

```go
// SOLUTION: Comprehensive path validation with allowlists

package main

import (
    "fmt"
    "os"
    "your-project/pkg/security"
)

// ✅ SECURE - Path validation with allowlist
func readUserFileSecure(userPath string) ([]byte, error) {
    // Define allowed directories
    allowedDirs := []string{
        "/app/user-data",
        "/app/uploads",
        "/var/app/files",
    }
    
    // Validate path before use
    if err := security.ValidatePath(userPath, allowedDirs); err != nil {
        // Generic error message - no information disclosure
        return nil, fmt.Errorf("access denied: invalid file path")
    }
    
    return os.ReadFile(userPath)
}

// ✅ SECURE - Secure file upload handling
func handleFileUpload(filename string, data []byte) error {
    // Sanitize filename
    cleanName := filepath.Base(filename) // Remove any path components
    
    // Construct safe path
    uploadDir := "/app/uploads"
    safePath := filepath.Join(uploadDir, cleanName)
    
    // Validate the constructed path
    if err := security.ValidatePath(safePath, []string{uploadDir}); err != nil {
        return fmt.Errorf("invalid upload path")
    }
    
    // Write file securely
    return security.WriteFileAtomic(safePath, data, 0644)
}

// ✅ SECURE - Configuration file access with validation
func loadConfigFile(configName string) ([]byte, error) {
    // Allowlist of valid config names
    validConfigs := map[string]string{
        "app":      "/app/config/app.json",
        "database": "/app/config/database.json",
        "logging":  "/app/config/logging.json",
    }
    
    configPath, exists := validConfigs[configName]
    if !exists {
        return nil, fmt.Errorf("unknown configuration: %s", configName)
    }
    
    // Additional path validation
    allowedDirs := []string{"/app/config"}
    if err := security.ValidatePath(configPath, allowedDirs); err != nil {
        return nil, fmt.Errorf("configuration access denied")
    }
    
    return os.ReadFile(configPath)
}

// SECURITY BENEFITS:
// - Allowlist prevents access to unauthorized directories
// - Path sanitization removes dangerous components
// - Generic error messages prevent information disclosure
// - Multiple validation layers provide defense in depth
```

## Configuration Management Patterns

### ❌ Insecure: Configuration Handling

```go
// VULNERABILITY: Insecure configuration management

package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type DatabaseConfig struct {
    Host     string `json:"host"`
    Password string `json:"password"`
    APIKey   string `json:"api_key"`
}

// ❌ NEVER DO THIS - Insecure configuration save
func saveConfigInsecure(config *DatabaseConfig, path string) error {
    data, _ := json.Marshal(config)
    
    // VULNERABILITIES:
    // 1. No path validation - directory traversal possible
    // 2. World-readable permissions expose secrets
    // 3. Non-atomic write can corrupt config
    // 4. No backup in case of failure
    return os.WriteFile(path, data, 0644) // Wrong permissions!
}

// ❌ NEVER DO THIS - Hardcoded sensitive paths
func loadProductionConfig() (*DatabaseConfig, error) {
    // Hardcoded path - inflexible and potentially insecure
    return loadConfigFromFile("/etc/myapp/production.json")
}

// ❌ NEVER DO THIS - Logging sensitive configuration
func debugConfig(config *DatabaseConfig) {
    // VULNERABILITY: Logs sensitive data
    fmt.Printf("Config: %+v\n", config) // Exposes password and API key
}

// SECURITY IMPACT:
// - Credential exposure through file permissions
// - Configuration corruption from non-atomic writes
// - Directory traversal attacks
// - Sensitive data in logs
```

### ✅ Secure: Configuration Management

```go
// SOLUTION: Secure configuration management

package main

import (
    "encoding/json"
    "fmt"
    "strings"
    "your-project/pkg/security"
)

type DatabaseConfig struct {
    Host     string `json:"host"`
    Password string `json:"password"`
    APIKey   string `json:"api_key"`
}

// ✅ SECURE - Comprehensive secure configuration save
func saveConfigSecure(config *DatabaseConfig, path string) error {
    // 1. Validate path
    allowedDirs := []string{"/app/config", "/etc/myapp"}
    if err := security.ValidatePath(path, allowedDirs); err != nil {
        return fmt.Errorf("invalid config path")
    }
    
    // 2. Create backup with secure naming
    backupSuffix, err := security.GenerateRandomSuffix(8)
    if err != nil {
        return fmt.Errorf("failed to generate backup suffix: %w", err)
    }
    
    backupPath := path + ".backup." + backupSuffix
    
    // 3. Serialize configuration
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // 4. Create backup if original exists
    if _, err := os.Stat(path); err == nil {
        originalData, err := os.ReadFile(path)
        if err == nil {
            security.WriteFileAtomic(backupPath, originalData, 0600)
        }
    }
    
    // 5. Write new config atomically with secure permissions
    return security.WriteFileAtomic(path, data, 0600) // Owner only
}

// ✅ SECURE - Flexible configuration loading with validation
func loadConfigSecure(configPath string) (*DatabaseConfig, error) {
    // Validate configuration path
    allowedDirs := []string{"/app/config", "/etc/myapp"}
    if err := security.ValidatePath(configPath, allowedDirs); err != nil {
        return nil, fmt.Errorf("invalid config path")
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config DatabaseConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &config, nil
}

// ✅ SECURE - Safe configuration logging (redacts sensitive fields)
func debugConfigSecure(config *DatabaseConfig) {
    // Create safe copy for logging
    safeConfig := struct {
        Host     string `json:"host"`
        Password string `json:"password"`
        APIKey   string `json:"api_key"`
    }{
        Host:     config.Host,
        Password: redactSensitive(config.Password),
        APIKey:   redactSensitive(config.APIKey),
    }
    
    fmt.Printf("Config: %+v\n", safeConfig)
}

func redactSensitive(value string) string {
    if len(value) <= 4 {
        return "***"
    }
    return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// ✅ SECURE - Configuration rotation with secure key generation
func rotateAPIKey(config *DatabaseConfig) error {
    // Generate new secure API key
    newAPIKey, err := security.GenerateBase64String(32)
    if err != nil {
        return fmt.Errorf("failed to generate new API key: %w", err)
    }
    
    config.APIKey = newAPIKey
    return nil
}

// SECURITY BENEFITS:
// - Secure file permissions protect credentials
// - Path validation prevents directory traversal
// - Atomic operations prevent corruption
// - Backup creation enables recovery
// - Safe logging prevents credential exposure
// - Secure key rotation maintains security
```

## Audit Trail Patterns

### ❌ Insecure: Predictable Audit IDs

```go
// VULNERABILITY: Predictable audit trail IDs (from pkg/reporting/audit.go)

package main

import (
    "fmt"
    "time"
)

type AuditEvent struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    UserID    string    `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details"`
}

// ❌ NEVER DO THIS - Original vulnerable pattern from audit.go
func generateEventIDInsecure() string {
    // VULNERABILITY: Predictable audit IDs
    return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

func createAuditEventInsecure(eventType, userID, details string) *AuditEvent {
    return &AuditEvent{
        ID:        generateEventIDInsecure(), // Predictable!
        Type:      eventType,
        UserID:    userID,
        Timestamp: time.Now(),
        Details:   details,
    }
}

// SECURITY IMPACT:
// - Attackers can predict future audit IDs
// - Audit trail manipulation becomes possible
// - Compliance violations (audit IDs must be tamper-resistant)
// - Forensic analysis becomes unreliable
```

### ✅ Secure: Cryptographically Secure Audit IDs

```go
// SOLUTION: Cryptographically secure audit trail

package main

import (
    "fmt"
    "time"
    "your-project/pkg/security"
)

// ✅ SECURE - Fixed audit ID generation
func generateEventIDSecure() (string, error) {
    return security.GenerateSecureID("audit")
}

func createAuditEventSecure(eventType, userID, details string) (*AuditEvent, error) {
    // Generate cryptographically secure audit ID
    auditID, err := generateEventIDSecure()
    if err != nil {
        return nil, fmt.Errorf("failed to generate secure audit ID: %w", err)
    }
    
    return &AuditEvent{
        ID:        auditID,
        Type:      eventType,
        UserID:    userID,
        Timestamp: time.Now().UTC(), // Always use UTC for consistency
        Details:   details,
    }, nil
}

// ✅ SECURE - Complete audit trail with integrity protection
func logSecurityEvent(eventType, userID, details string) error {
    // Create audit event with secure ID
    event, err := createAuditEventSecure(eventType, userID, details)
    if err != nil {
        return fmt.Errorf("failed to create audit event: %w", err)
    }
    
    // Serialize audit event
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to serialize audit event: %w", err)
    }
    
    // Write to secure audit log with atomic operation
    auditLogPath := "/var/log/myapp/audit.log"
    
    // Validate audit log path
    allowedDirs := []string{"/var/log/myapp", "/app/logs"}
    if err := security.ValidatePath(auditLogPath, allowedDirs); err != nil {
        return fmt.Errorf("invalid audit log path")
    }
    
    // Append to audit log atomically
    return appendToAuditLogSecure(auditLogPath, data)
}

func appendToAuditLogSecure(logPath string, data []byte) error {
    // Read existing log
    existingData, err := os.ReadFile(logPath)
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to read existing audit log: %w", err)
    }
    
    // Append new entry with newline
    newData := append(existingData, data...)
    newData = append(newData, '\n')
    
    // Write atomically with secure permissions
    return security.WriteFileAtomic(logPath, newData, 0600)
}

// SECURITY BENEFITS:
// - Unpredictable audit IDs prevent manipulation
// - Atomic log writes prevent corruption
// - Secure permissions protect audit data
// - Tamper-resistant audit trail for compliance
```

## Error Handling Patterns

### ❌ Insecure: Information Disclosure in Errors

```go
// VULNERABILITY: Error messages that leak sensitive information

package main

import (
    "fmt"
    "os"
)

// ❌ NEVER DO THIS - Leaks sensitive path information
func readFileInsecure(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // VULNERABILITY: Exposes internal file paths
        return nil, fmt.Errorf("failed to read file %s: %v", path, err)
    }
    return data, nil
}

// ❌ NEVER DO THIS - Leaks system information
func validateUserInsecure(username string) error {
    if username == "admin" {
        // VULNERABILITY: Confirms existence of admin user
        return fmt.Errorf("admin user cannot be modified")
    }
    
    userFile := fmt.Sprintf("/etc/users/%s.json", username)
    if _, err := os.Stat(userFile); err != nil {
        // VULNERABILITY: Leaks internal file structure
        return fmt.Errorf("user file %s not found: %v", userFile, err)
    }
    
    return nil
}

// ❌ NEVER DO THIS - Database errors leak schema information
func queryUserInsecure(userID int) (*User, error) {
    query := "SELECT * FROM users WHERE id = ?"
    row := db.QueryRow(query, userID)
    
    var user User
    if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
        // VULNERABILITY: Exposes database schema and query details
        return nil, fmt.Errorf("database query failed: %s with ID %d: %v", query, userID, err)
    }
    
    return &user, nil
}

// SECURITY IMPACT:
// - Information disclosure aids attackers
// - Reveals internal system structure
// - Exposes file paths and database schema
// - Helps with reconnaissance for further attacks
```

### ✅ Secure: Safe Error Handling

```go
// SOLUTION: Generic error messages with secure logging

package main

import (
    "fmt"
    "log"
    "your-project/pkg/security"
)

// ✅ SECURE - Generic error messages with secure logging
func readFileSecure(path string) ([]byte, error) {
    // Validate path first
    allowedDirs := []string{"/app/data", "/app/config"}
    if err := security.ValidatePath(path, allowedDirs); err != nil {
        // Log detailed error securely (not returned to user)
        log.Printf("Path validation failed for %s: %v", path, err)
        // Return generic error to user
        return nil, fmt.Errorf("access denied")
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        // Log detailed error securely
        log.Printf("File read failed for %s: %v", path, err)
        // Return generic error to user
        return nil, fmt.Errorf("file operation failed")
    }
    
    return data, nil
}

// ✅ SECURE - Safe user validation with audit logging
func validateUserSecure(username string) error {
    // Generate audit ID for this validation attempt
    auditID, err := security.GenerateSecureID("validation")
    if err != nil {
        log.Printf("Failed to generate audit ID: %v", err)
    }
    
    // Log validation attempt with audit ID
    log.Printf("User validation attempt - Audit ID: %s, Username: %s", auditID, username)
    
    // Perform validation without revealing system details
    if !isValidUser(username) {
        // Log detailed information securely
        log.Printf("User validation failed - Audit ID: %s, Username: %s", auditID, username)
        // Return generic error
        return fmt.Errorf("user validation failed")
    }
    
    log.Printf("User validation successful - Audit ID: %s", auditID)
    return nil
}

// ✅ SECURE - Database error handling with audit trail
func queryUserSecure(userID int) (*User, error) {
    // Generate audit ID for database operation
    auditID, err := security.GenerateSecureID("db_query")
    if err != nil {
        log.Printf("Failed to generate audit ID for database query: %v", err)
    }
    
    // Log query attempt
    log.Printf("Database query attempt - Audit ID: %s, UserID: %d", auditID, userID)
    
    query := "SELECT * FROM users WHERE id = ?"
    row := db.QueryRow(query, userID)
    
    var user User
    if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
        // Log detailed error securely with audit ID
        log.Printf("Database query failed - Audit ID: %s, UserID: %d, Error: %v", auditID, userID, err)
        
        // Return generic error to user
        return nil, fmt.Errorf("user lookup failed")
    }
    
    log.Printf("Database query successful - Audit ID: %s", auditID)
    return &user, nil
}

// ✅ SECURE - Structured error handling with security context
type SecurityError struct {
    AuditID   string
    Operation string
    UserID    string
    Message   string // Generic message for user
    Details   string // Detailed message for logs
}

func (se *SecurityError) Error() string {
    return se.Message
}

func (se *SecurityError) LogSecurely() {
    log.Printf("Security Error - Audit ID: %s, Operation: %s, User: %s, Details: %s",
        se.AuditID, se.Operation, se.UserID, se.Details)
}

func performSecureOperation(userID, operation string) error {
    auditID, _ := security.GenerateSecureID("operation")
    
    // Simulate operation that might fail
    if err := doSomethingRisky(); err != nil {
        secErr := &SecurityError{
            AuditID:   auditID,
            Operation: operation,
            UserID:    userID,
            Message:   "Operation failed", // Generic for user
            Details:   fmt.Sprintf("Detailed error: %v", err), // Detailed for logs
        }
        
        secErr.LogSecurely()
        return secErr
    }
    
    return nil
}

// SECURITY BENEFITS:
// - Generic error messages prevent information disclosure
// - Detailed logging enables debugging and forensics
// - Audit IDs enable correlation of events
// - Structured error handling ensures consistency
// - Security context preserved for analysis
```

## Summary

These examples demonstrate the critical importance of:

1. **Cryptographically Secure Random Generation**: Always use `crypto/rand` for security-sensitive operations
2. **Atomic File Operations**: Prevent race conditions and ensure data consistency
3. **Path Validation**: Prevent directory traversal and unauthorized file access
4. **Secure Error Handling**: Avoid information disclosure while maintaining auditability
5. **Defense in Depth**: Multiple layers of security validation and protection

The patterns shown here address real vulnerabilities found in the codebase and provide secure alternatives that maintain functionality while eliminating security risks. Always prefer the secure patterns and avoid the anti-patterns, even if they seem more convenient or performant.
