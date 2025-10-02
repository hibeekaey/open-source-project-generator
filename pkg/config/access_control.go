package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// AccessLevel defines the level of access to configuration operations
type AccessLevel int

const (
	AccessLevelNone AccessLevel = iota
	AccessLevelRead
	AccessLevelWrite
	AccessLevelAdmin
)

// AccessPolicy defines access control policies for configuration operations
type AccessPolicy struct {
	User        string      `json:"user"`
	Resource    string      `json:"resource"`
	Action      string      `json:"action"`
	AccessLevel AccessLevel `json:"access_level"`
	Conditions  []string    `json:"conditions,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
}

// AccessController manages configuration access control
type AccessController struct {
	policies    []*AccessPolicy
	logger      interfaces.Logger
	auditLogger *ConfigAuditLogger
}

// SecurityContext contains security information for configuration operations
type SecurityContext struct {
	User      string            `json:"user"`
	IPAddress string            `json:"ip_address,omitempty"`
	UserAgent string            `json:"user_agent,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// SecureConfigManager wraps the configuration manager with security features
type SecureConfigManager struct {
	manager           *Manager
	accessController  *AccessController
	securityContext   *SecurityContext
	encryptionEnabled bool
	auditEnabled      bool
}

// NewAccessController creates a new access controller
func NewAccessController(logger interfaces.Logger, auditLogger *ConfigAuditLogger) *AccessController {
	return &AccessController{
		policies:    []*AccessPolicy{},
		logger:      logger,
		auditLogger: auditLogger,
	}
}

// NewSecureConfigManager creates a new secure configuration manager
func NewSecureConfigManager(manager *Manager, logger interfaces.Logger) (*SecureConfigManager, error) {
	auditLogger := &ConfigAuditLogger{
		logFile: filepath.Join(manager.configDir, "security_audit.log"),
		logger:  logger,
	}

	accessController := NewAccessController(logger, auditLogger)

	// Set up default policies
	if err := accessController.setupDefaultPolicies(); err != nil {
		return nil, fmt.Errorf("failed to setup default policies: %w", err)
	}

	return &SecureConfigManager{
		manager:           manager,
		accessController:  accessController,
		encryptionEnabled: true,
		auditEnabled:      true,
	}, nil
}

// SetSecurityContext sets the security context for operations
func (s *SecureConfigManager) SetSecurityContext(ctx *SecurityContext) {
	s.securityContext = ctx
}

// ExportConfiguration exports a configuration with security checks
func (s *SecureConfigManager) ExportConfiguration(name string, options *ConfigExportOptions) ([]byte, error) {
	// Check access permissions
	if err := s.checkAccess("export", name); err != nil {
		s.auditSecurityEvent("export_denied", name, err)
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Enable encryption for sensitive data by default
	if options == nil {
		options = &ConfigExportOptions{
			Format:         "yaml",
			IncludeMeta:    true,
			EncryptSecrets: s.encryptionEnabled,
		}
	} else if s.encryptionEnabled {
		options.EncryptSecrets = true
	}

	// Perform export
	data, err := s.manager.ExportConfiguration(name, options)
	if err != nil {
		s.auditSecurityEvent("export_failed", name, err)
		return nil, err
	}

	// Audit successful export
	s.auditSecurityEvent("export_success", name, nil)

	return data, nil
}

// ImportConfiguration imports a configuration with security checks
func (s *SecureConfigManager) ImportConfiguration(name string, data []byte, options *ConfigImportOptions) error {
	// Check access permissions
	if err := s.checkAccess("import", name); err != nil {
		s.auditSecurityEvent("import_denied", name, err)
		return fmt.Errorf("access denied: %w", err)
	}

	// Validate import data security
	if err := s.validateImportSecurity(data); err != nil {
		s.auditSecurityEvent("import_security_violation", name, err)
		return fmt.Errorf("security validation failed: %w", err)
	}

	// Enable schema validation by default
	if options == nil {
		options = &ConfigImportOptions{
			Format:         "auto",
			ValidateSchema: true,
			MergeStrategy:  "replace",
			Transform:      true,
		}
	} else {
		options.ValidateSchema = true
	}

	// Perform import
	err := s.manager.ImportConfiguration(name, data, options)
	if err != nil {
		s.auditSecurityEvent("import_failed", name, err)
		return err
	}

	// Audit successful import
	s.auditSecurityEvent("import_success", name, nil)

	return nil
}

// DeleteConfiguration deletes a configuration with security checks
func (s *SecureConfigManager) DeleteConfiguration(name string) error {
	// Check access permissions
	if err := s.checkAccess("delete", name); err != nil {
		s.auditSecurityEvent("delete_denied", name, err)
		return fmt.Errorf("access denied: %w", err)
	}

	// Perform deletion
	err := s.manager.DeleteConfiguration(name)
	if err != nil {
		s.auditSecurityEvent("delete_failed", name, err)
		return err
	}

	// Audit successful deletion
	s.auditSecurityEvent("delete_success", name, nil)

	return nil
}

// ListConfigurations lists configurations with security filtering
func (s *SecureConfigManager) ListConfigurations(options *ConfigListOptions) ([]*ConfigInfo, error) {
	// Check access permissions
	if err := s.checkAccess("list", "*"); err != nil {
		s.auditSecurityEvent("list_denied", "*", err)
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Get configurations
	configs, err := s.manager.ListConfigurations(options)
	if err != nil {
		return nil, err
	}

	// Filter configurations based on access permissions
	var filteredConfigs []*ConfigInfo
	for _, config := range configs {
		if err := s.checkAccess("read", config.Name); err == nil {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	return filteredConfigs, nil
}

// checkAccess checks if the current user has access to perform an action
func (s *SecureConfigManager) checkAccess(action, resource string) error {
	if s.securityContext == nil {
		return fmt.Errorf("no security context set")
	}

	user := s.securityContext.User
	if user == "" {
		user = s.getCurrentUser()
	}

	// Check policies
	for _, policy := range s.accessController.policies {
		if s.matchesPolicy(policy, user, action, resource) {
			// Check if policy has expired
			if policy.ExpiresAt != nil && time.Now().After(*policy.ExpiresAt) {
				continue
			}

			// Check conditions
			if len(policy.Conditions) > 0 {
				if err := s.checkConditions(policy.Conditions); err != nil {
					continue
				}
			}

			// Check access level
			requiredLevel := s.getRequiredAccessLevel(action)
			if policy.AccessLevel >= requiredLevel {
				return nil
			}
		}
	}

	return fmt.Errorf("insufficient permissions for action '%s' on resource '%s'", action, resource)
}

// matchesPolicy checks if a policy matches the given parameters
func (s *SecureConfigManager) matchesPolicy(policy *AccessPolicy, user, action, resource string) bool {
	// Check user
	if policy.User != "*" && policy.User != user {
		return false
	}

	// Check action
	if policy.Action != "*" && policy.Action != action {
		return false
	}

	// Check resource
	if policy.Resource != "*" && policy.Resource != resource {
		// Check for wildcard patterns
		if !s.matchesPattern(policy.Resource, resource) {
			return false
		}
	}

	return true
}

// matchesPattern checks if a resource matches a pattern
func (s *SecureConfigManager) matchesPattern(pattern, resource string) bool {
	// Simple wildcard matching
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(resource, prefix)
	}

	return pattern == resource
}

// getRequiredAccessLevel returns the required access level for an action
func (s *SecureConfigManager) getRequiredAccessLevel(action string) AccessLevel {
	switch action {
	case "read", "list":
		return AccessLevelRead
	case "export", "import", "save", "create":
		return AccessLevelWrite
	case "delete", "admin":
		return AccessLevelAdmin
	default:
		return AccessLevelAdmin
	}
}

// checkConditions checks if security conditions are met
func (s *SecureConfigManager) checkConditions(conditions []string) error {
	for _, condition := range conditions {
		if err := s.evaluateCondition(condition); err != nil {
			return err
		}
	}
	return nil
}

// evaluateCondition evaluates a single security condition
func (s *SecureConfigManager) evaluateCondition(condition string) error {
	// Parse condition (simplified implementation)
	parts := strings.SplitN(condition, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid condition format: %s", condition)
	}

	key, expectedValue := parts[0], parts[1]

	switch key {
	case "ip_range":
		return s.checkIPRange(expectedValue)
	case "time_range":
		return s.checkTimeRange(expectedValue)
	case "user_agent":
		return s.checkUserAgent(expectedValue)
	default:
		return fmt.Errorf("unknown condition: %s", key)
	}
}

// checkIPRange checks if the current IP is in the allowed range
func (s *SecureConfigManager) checkIPRange(allowedRange string) error {
	if s.securityContext == nil || s.securityContext.IPAddress == "" {
		return fmt.Errorf("no IP address in security context")
	}

	// Simplified IP range check (in production, use proper CIDR matching)
	if allowedRange == "*" || allowedRange == s.securityContext.IPAddress {
		return nil
	}

	return fmt.Errorf("IP address %s not in allowed range %s", s.securityContext.IPAddress, allowedRange)
}

// checkTimeRange checks if the current time is in the allowed range
func (s *SecureConfigManager) checkTimeRange(timeRange string) error {
	// Simplified time range check (format: "09:00-17:00")
	now := time.Now()
	currentHour := now.Hour()

	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid time range format: %s", timeRange)
	}

	// Parse start and end hours (simplified)
	var startHour, endHour int
	if _, err := fmt.Sscanf(parts[0], "%d:", &startHour); err != nil {
		return fmt.Errorf("invalid start time: %s", parts[0])
	}
	if _, err := fmt.Sscanf(parts[1], "%d:", &endHour); err != nil {
		return fmt.Errorf("invalid end time: %s", parts[1])
	}

	if currentHour >= startHour && currentHour <= endHour {
		return nil
	}

	return fmt.Errorf("current time %d:00 not in allowed range %s", currentHour, timeRange)
}

// checkUserAgent checks if the user agent is allowed
func (s *SecureConfigManager) checkUserAgent(allowedPattern string) error {
	if s.securityContext == nil || s.securityContext.UserAgent == "" {
		return fmt.Errorf("no user agent in security context")
	}

	if allowedPattern == "*" || strings.Contains(s.securityContext.UserAgent, allowedPattern) {
		return nil
	}

	return fmt.Errorf("user agent not allowed")
}

// validateImportSecurity validates the security of import data
func (s *SecureConfigManager) validateImportSecurity(data []byte) error {
	// Check for potentially dangerous content
	dataStr := string(data)

	// Check for script injection attempts
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"eval(",
		"exec(",
		"system(",
		"shell_exec(",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(dataStr), pattern) {
			return fmt.Errorf("potentially dangerous content detected: %s", pattern)
		}
	}

	// Check file size
	maxSize := 10 * 1024 * 1024 // 10MB
	if len(data) > maxSize {
		return fmt.Errorf("import data too large: %d bytes (max: %d)", len(data), maxSize)
	}

	return nil
}

// auditSecurityEvent logs a security event
func (s *SecureConfigManager) auditSecurityEvent(event, resource string, err error) {
	if !s.auditEnabled {
		return
	}

	details := map[string]interface{}{
		"event":    event,
		"resource": resource,
	}

	if s.securityContext != nil {
		details["user"] = s.securityContext.User
		details["ip_address"] = s.securityContext.IPAddress
		details["user_agent"] = s.securityContext.UserAgent
		details["session_id"] = s.securityContext.SessionID
	}

	if err != nil {
		s.accessController.auditLogger.LogError("security_event", resource, err, details)
	} else {
		s.accessController.auditLogger.LogAction("security_event", resource, details)
	}
}

// getCurrentUser gets the current user from the environment
func (s *SecureConfigManager) getCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

// setupDefaultPolicies sets up default access control policies
func (a *AccessController) setupDefaultPolicies() error {
	// Default policy: current user has full access
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}
	if currentUser == "" {
		currentUser = "*" // Allow all users if we can't determine current user
	}

	defaultPolicies := []*AccessPolicy{
		{
			User:        currentUser,
			Resource:    "*",
			Action:      "*",
			AccessLevel: AccessLevelAdmin,
		},
		{
			User:        "*",
			Resource:    "*",
			Action:      "read",
			AccessLevel: AccessLevelRead,
		},
		{
			User:        "*",
			Resource:    "*",
			Action:      "list",
			AccessLevel: AccessLevelRead,
		},
	}

	a.policies = append(a.policies, defaultPolicies...)
	return nil
}

// AddPolicy adds a new access control policy
func (a *AccessController) AddPolicy(policy *AccessPolicy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	if policy.User == "" || policy.Resource == "" || policy.Action == "" {
		return fmt.Errorf("policy must have user, resource, and action")
	}

	a.policies = append(a.policies, policy)

	if a.auditLogger != nil {
		a.auditLogger.LogAction("add_policy", "access_control", map[string]interface{}{
			"user":         policy.User,
			"resource":     policy.Resource,
			"action":       policy.Action,
			"access_level": policy.AccessLevel,
		})
	}

	return nil
}

// RemovePolicy removes an access control policy
func (a *AccessController) RemovePolicy(user, resource, action string) error {
	var newPolicies []*AccessPolicy

	removed := false
	for _, policy := range a.policies {
		if policy.User == user && policy.Resource == resource && policy.Action == action {
			removed = true
			continue
		}
		newPolicies = append(newPolicies, policy)
	}

	if !removed {
		return fmt.Errorf("policy not found")
	}

	a.policies = newPolicies

	if a.auditLogger != nil {
		a.auditLogger.LogAction("remove_policy", "access_control", map[string]interface{}{
			"user":     user,
			"resource": resource,
			"action":   action,
		})
	}

	return nil
}

// ListPolicies returns all access control policies
func (a *AccessController) ListPolicies() []*AccessPolicy {
	// Return a copy to prevent modification
	policies := make([]*AccessPolicy, len(a.policies))
	copy(policies, a.policies)
	return policies
}

// ValidatePolicy validates an access control policy
func (a *AccessController) ValidatePolicy(policy *AccessPolicy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	if policy.User == "" {
		return fmt.Errorf("policy user cannot be empty")
	}

	if policy.Resource == "" {
		return fmt.Errorf("policy resource cannot be empty")
	}

	if policy.Action == "" {
		return fmt.Errorf("policy action cannot be empty")
	}

	if policy.AccessLevel < AccessLevelNone || policy.AccessLevel > AccessLevelAdmin {
		return fmt.Errorf("invalid access level: %d", policy.AccessLevel)
	}

	// Validate conditions
	for _, condition := range policy.Conditions {
		if !strings.Contains(condition, "=") {
			return fmt.Errorf("invalid condition format: %s", condition)
		}
	}

	return nil
}
