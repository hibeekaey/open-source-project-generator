package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// SensitiveFields defines which fields should be encrypted
var SensitiveFields = map[string]bool{
	"email":                true,
	"repository":           true,
	"author":               true,
	"default_email":        true,
	"default_author":       true,
	"default_organization": true,
	"api_key":              true,
	"token":                true,
	"password":             true,
	"secret":               true,
}

// EncryptedFieldPrefix is the prefix used to identify encrypted fields
const EncryptedFieldPrefix = "encrypted:"

// EncryptSensitiveFields encrypts sensitive fields in the configuration data
func (e *ConfigEncryptor) EncryptSensitiveFields(data interface{}) (interface{}, error) {
	if e.key == nil {
		return data, fmt.Errorf("encryption key not set")
	}

	return e.encryptValue(data), nil
}

// DecryptSensitiveFields decrypts sensitive fields in the configuration data
func (e *ConfigEncryptor) DecryptSensitiveFields(data interface{}) (interface{}, error) {
	if e.key == nil {
		return data, fmt.Errorf("encryption key not set")
	}

	decrypted, err := e.decryptValue(data)
	if err != nil {
		return data, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decrypted, nil
}

// encryptValue recursively encrypts sensitive values in the data structure
func (e *ConfigEncryptor) encryptValue(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if e.isSensitiveField(key) {
				if strValue, ok := value.(string); ok && strValue != "" {
					if encrypted, err := e.encryptString(strValue); err == nil {
						result[key] = EncryptedFieldPrefix + encrypted
					} else {
						if e.logger != nil {
							e.logger.WarnWithFields("Failed to encrypt field", map[string]interface{}{
								"field": key,
								"error": err.Error(),
							})
						}
						result[key] = value
					}
				} else {
					result[key] = value
				}
			} else {
				result[key] = e.encryptValue(value)
			}
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = e.encryptValue(item)
		}
		return result

	default:
		return data
	}
}

// decryptValue recursively decrypts encrypted values in the data structure
func (e *ConfigEncryptor) decryptValue(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if strValue, ok := value.(string); ok && strings.HasPrefix(strValue, EncryptedFieldPrefix) {
				encryptedData := strings.TrimPrefix(strValue, EncryptedFieldPrefix)
				if decrypted, err := e.decryptString(encryptedData); err == nil {
					result[key] = decrypted
				} else {
					if e.logger != nil {
						e.logger.WarnWithFields("Failed to decrypt field", map[string]interface{}{
							"field": key,
							"error": err.Error(),
						})
					}
					result[key] = value
				}
			} else {
				decryptedValue, err := e.decryptValue(value)
				if err != nil {
					return nil, err
				}
				result[key] = decryptedValue
			}
		}
		return result, nil

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			decryptedItem, err := e.decryptValue(item)
			if err != nil {
				return nil, err
			}
			result[i] = decryptedItem
		}
		return result, nil

	default:
		return data, nil
	}
}

// encryptString encrypts a string value using AES-GCM
func (e *ConfigEncryptor) encryptString(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptString decrypts a string value using AES-GCM
func (e *ConfigEncryptor) decryptString(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check minimum length
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// isSensitiveField checks if a field name is considered sensitive
func (e *ConfigEncryptor) isSensitiveField(fieldName string) bool {
	// Check exact match
	if SensitiveFields[strings.ToLower(fieldName)] {
		return true
	}

	// Check if field name contains sensitive keywords
	lowerField := strings.ToLower(fieldName)
	sensitiveKeywords := []string{"password", "secret", "token", "key", "credential"}
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerField, keyword) {
			return true
		}
	}

	return false
}

// EncryptConfiguration encrypts an entire configuration structure
func (e *ConfigEncryptor) EncryptConfiguration(config *SavedConfiguration) (*SavedConfiguration, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Create a copy to avoid modifying the original
	encrypted := *config

	// Encrypt user preferences if present
	if config.UserPreferences != nil {
		encryptedPrefs, err := e.encryptUserPreferences(config.UserPreferences)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt user preferences: %w", err)
		}
		encrypted.UserPreferences = encryptedPrefs
	}

	// Encrypt project configuration if present
	if config.ProjectConfig != nil {
		encryptedProject, err := e.encryptProjectConfig(config.ProjectConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt project config: %w", err)
		}
		encrypted.ProjectConfig = encryptedProject
	}

	return &encrypted, nil
}

// DecryptConfiguration decrypts an entire configuration structure
func (e *ConfigEncryptor) DecryptConfiguration(config *SavedConfiguration) (*SavedConfiguration, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Create a copy to avoid modifying the original
	decrypted := *config

	// Decrypt user preferences if present
	if config.UserPreferences != nil {
		decryptedPrefs, err := e.decryptUserPreferences(config.UserPreferences)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt user preferences: %w", err)
		}
		decrypted.UserPreferences = decryptedPrefs
	}

	// Decrypt project configuration if present
	if config.ProjectConfig != nil {
		decryptedProject, err := e.decryptProjectConfig(config.ProjectConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt project config: %w", err)
		}
		decrypted.ProjectConfig = decryptedProject
	}

	return &decrypted, nil
}

// encryptUserPreferences encrypts sensitive fields in user preferences
func (e *ConfigEncryptor) encryptUserPreferences(prefs *UserPreferences) (*UserPreferences, error) {
	encrypted := *prefs

	// Encrypt sensitive fields
	if prefs.DefaultEmail != "" {
		if encryptedEmail, err := e.encryptString(prefs.DefaultEmail); err == nil {
			encrypted.DefaultEmail = EncryptedFieldPrefix + encryptedEmail
		} else {
			return nil, fmt.Errorf("failed to encrypt default email: %w", err)
		}
	}

	if prefs.DefaultAuthor != "" {
		if encryptedAuthor, err := e.encryptString(prefs.DefaultAuthor); err == nil {
			encrypted.DefaultAuthor = EncryptedFieldPrefix + encryptedAuthor
		} else {
			return nil, fmt.Errorf("failed to encrypt default author: %w", err)
		}
	}

	if prefs.DefaultOrganization != "" {
		if encryptedOrg, err := e.encryptString(prefs.DefaultOrganization); err == nil {
			encrypted.DefaultOrganization = EncryptedFieldPrefix + encryptedOrg
		} else {
			return nil, fmt.Errorf("failed to encrypt default organization: %w", err)
		}
	}

	// Encrypt custom defaults
	if len(prefs.CustomDefaults) > 0 {
		encryptedDefaults := make(map[string]string)
		for key, value := range prefs.CustomDefaults {
			if e.isSensitiveField(key) && value != "" {
				if encryptedValue, err := e.encryptString(value); err == nil {
					encryptedDefaults[key] = EncryptedFieldPrefix + encryptedValue
				} else {
					return nil, fmt.Errorf("failed to encrypt custom default %s: %w", key, err)
				}
			} else {
				encryptedDefaults[key] = value
			}
		}
		encrypted.CustomDefaults = encryptedDefaults
	}

	return &encrypted, nil
}

// decryptUserPreferences decrypts sensitive fields in user preferences
func (e *ConfigEncryptor) decryptUserPreferences(prefs *UserPreferences) (*UserPreferences, error) {
	decrypted := *prefs

	// Decrypt sensitive fields
	if strings.HasPrefix(prefs.DefaultEmail, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(prefs.DefaultEmail, EncryptedFieldPrefix)
		if decryptedEmail, err := e.decryptString(encryptedData); err == nil {
			decrypted.DefaultEmail = decryptedEmail
		} else {
			return nil, fmt.Errorf("failed to decrypt default email: %w", err)
		}
	}

	if strings.HasPrefix(prefs.DefaultAuthor, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(prefs.DefaultAuthor, EncryptedFieldPrefix)
		if decryptedAuthor, err := e.decryptString(encryptedData); err == nil {
			decrypted.DefaultAuthor = decryptedAuthor
		} else {
			return nil, fmt.Errorf("failed to decrypt default author: %w", err)
		}
	}

	if strings.HasPrefix(prefs.DefaultOrganization, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(prefs.DefaultOrganization, EncryptedFieldPrefix)
		if decryptedOrg, err := e.decryptString(encryptedData); err == nil {
			decrypted.DefaultOrganization = decryptedOrg
		} else {
			return nil, fmt.Errorf("failed to decrypt default organization: %w", err)
		}
	}

	// Decrypt custom defaults
	if len(prefs.CustomDefaults) > 0 {
		decryptedDefaults := make(map[string]string)
		for key, value := range prefs.CustomDefaults {
			if strings.HasPrefix(value, EncryptedFieldPrefix) {
				encryptedData := strings.TrimPrefix(value, EncryptedFieldPrefix)
				if decryptedValue, err := e.decryptString(encryptedData); err == nil {
					decryptedDefaults[key] = decryptedValue
				} else {
					return nil, fmt.Errorf("failed to decrypt custom default %s: %w", key, err)
				}
			} else {
				decryptedDefaults[key] = value
			}
		}
		decrypted.CustomDefaults = decryptedDefaults
	}

	return &decrypted, nil
}

// encryptProjectConfig encrypts sensitive fields in project configuration
func (e *ConfigEncryptor) encryptProjectConfig(config *models.ProjectConfig) (*models.ProjectConfig, error) {
	encrypted := *config

	// Encrypt sensitive fields
	if config.Email != "" {
		if encryptedEmail, err := e.encryptString(config.Email); err == nil {
			encrypted.Email = EncryptedFieldPrefix + encryptedEmail
		} else {
			return nil, fmt.Errorf("failed to encrypt email: %w", err)
		}
	}

	if config.Author != "" {
		if encryptedAuthor, err := e.encryptString(config.Author); err == nil {
			encrypted.Author = EncryptedFieldPrefix + encryptedAuthor
		} else {
			return nil, fmt.Errorf("failed to encrypt author: %w", err)
		}
	}

	if config.Repository != "" {
		if encryptedRepo, err := e.encryptString(config.Repository); err == nil {
			encrypted.Repository = EncryptedFieldPrefix + encryptedRepo
		} else {
			return nil, fmt.Errorf("failed to encrypt repository: %w", err)
		}
	}

	return &encrypted, nil
}

// decryptProjectConfig decrypts sensitive fields in project configuration
func (e *ConfigEncryptor) decryptProjectConfig(config *models.ProjectConfig) (*models.ProjectConfig, error) {
	decrypted := *config

	// Decrypt sensitive fields
	if strings.HasPrefix(config.Email, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(config.Email, EncryptedFieldPrefix)
		if decryptedEmail, err := e.decryptString(encryptedData); err == nil {
			decrypted.Email = decryptedEmail
		} else {
			return nil, fmt.Errorf("failed to decrypt email: %w", err)
		}
	}

	if strings.HasPrefix(config.Author, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(config.Author, EncryptedFieldPrefix)
		if decryptedAuthor, err := e.decryptString(encryptedData); err == nil {
			decrypted.Author = decryptedAuthor
		} else {
			return nil, fmt.Errorf("failed to decrypt author: %w", err)
		}
	}

	if strings.HasPrefix(config.Repository, EncryptedFieldPrefix) {
		encryptedData := strings.TrimPrefix(config.Repository, EncryptedFieldPrefix)
		if decryptedRepo, err := e.decryptString(encryptedData); err == nil {
			decrypted.Repository = decryptedRepo
		} else {
			return nil, fmt.Errorf("failed to decrypt repository: %w", err)
		}
	}

	return &decrypted, nil
}

// ValidateEncryptionKey validates that the encryption key is properly set
func (e *ConfigEncryptor) ValidateEncryptionKey() error {
	if e.key == nil {
		return fmt.Errorf("encryption key is not set")
	}

	if len(e.key) != 32 {
		return fmt.Errorf("encryption key must be 32 bytes long, got %d", len(e.key))
	}

	// Test encryption/decryption with a sample string
	testString := "test-encryption-validation"
	encrypted, err := e.encryptString(testString)
	if err != nil {
		return fmt.Errorf("encryption test failed: %w", err)
	}

	decrypted, err := e.decryptString(encrypted)
	if err != nil {
		return fmt.Errorf("decryption test failed: %w", err)
	}

	if decrypted != testString {
		return fmt.Errorf("encryption/decryption test failed: expected %s, got %s", testString, decrypted)
	}

	return nil
}

// IsEncrypted checks if a value is encrypted
func (e *ConfigEncryptor) IsEncrypted(value string) bool {
	return strings.HasPrefix(value, EncryptedFieldPrefix)
}

// GetEncryptionInfo returns information about the encryption setup
func (e *ConfigEncryptor) GetEncryptionInfo() map[string]interface{} {
	info := map[string]interface{}{
		"algorithm":        "AES-256-GCM",
		"key_length":       len(e.key),
		"prefix":           EncryptedFieldPrefix,
		"sensitive_fields": getSensitiveFieldsList(),
	}

	if err := e.ValidateEncryptionKey(); err != nil {
		info["status"] = "invalid"
		info["error"] = err.Error()
	} else {
		info["status"] = "valid"
	}

	return info
}

// getSensitiveFieldsList returns a list of sensitive field names
func getSensitiveFieldsList() []string {
	var fields []string
	for field := range SensitiveFields {
		fields = append(fields, field)
	}
	return fields
}
