// Package audit provides rule management functionality for the audit engine.
//
// This file contains the RuleManager struct and related functionality for
// managing audit rules including registration, filtering, and validation.
package audit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// RuleManager manages audit rules for the audit engine.
//
// It provides functionality for:
//   - Rule registration and management
//   - Rule filtering by category, type, and severity
//   - Rule validation and configuration
//   - Default rule definitions
type RuleManager struct {
	rules []interfaces.AuditRule
}

// NewRuleManager creates a new rule manager with default rules.
func NewRuleManager() *RuleManager {
	return &RuleManager{
		rules: getDefaultAuditRules(),
	}
}

// NewRuleManagerWithRules creates a new rule manager with custom rules.
func NewRuleManagerWithRules(rules []interfaces.AuditRule) *RuleManager {
	rm := &RuleManager{
		rules: make([]interfaces.AuditRule, len(rules)),
	}
	copy(rm.rules, rules)
	return rm
}

// GetRules returns a copy of all audit rules.
func (rm *RuleManager) GetRules() []interfaces.AuditRule {
	rules := make([]interfaces.AuditRule, len(rm.rules))
	copy(rules, rm.rules)
	return rules
}

// SetRules replaces all audit rules with the provided rules.
func (rm *RuleManager) SetRules(rules []interfaces.AuditRule) error {
	// Validate all rules before setting
	for _, rule := range rules {
		if err := rm.ValidateRule(rule); err != nil {
			return fmt.Errorf("invalid rule %s: %w", rule.ID, err)
		}
	}

	rm.rules = make([]interfaces.AuditRule, len(rules))
	copy(rm.rules, rules)
	return nil
}

// AddRule adds a new audit rule or replaces an existing one with the same ID.
func (rm *RuleManager) AddRule(rule interfaces.AuditRule) error {
	if err := rm.ValidateRule(rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	// Check if rule with same ID already exists
	for i, existingRule := range rm.rules {
		if existingRule.ID == rule.ID {
			rm.rules[i] = rule // Replace existing rule
			return nil
		}
	}

	// Add new rule
	rm.rules = append(rm.rules, rule)
	return nil
}

// RemoveRule removes an audit rule by ID.
func (rm *RuleManager) RemoveRule(ruleID string) error {
	for i, rule := range rm.rules {
		if rule.ID == ruleID {
			rm.rules = append(rm.rules[:i], rm.rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("audit rule with ID %s not found", ruleID)
}

// GetRule returns a specific rule by ID.
func (rm *RuleManager) GetRule(ruleID string) (*interfaces.AuditRule, error) {
	for _, rule := range rm.rules {
		if rule.ID == ruleID {
			ruleCopy := rule
			return &ruleCopy, nil
		}
	}
	return nil, fmt.Errorf("audit rule with ID %s not found", ruleID)
}

// FilterRules returns rules filtered by the specified criteria.
func (rm *RuleManager) FilterRules(filter RuleFilter) []interfaces.AuditRule {
	var filtered []interfaces.AuditRule

	for _, rule := range rm.rules {
		if rm.matchesFilter(rule, filter) {
			filtered = append(filtered, rule)
		}
	}

	return filtered
}

// GetRulesByCategory returns all rules for a specific category.
func (rm *RuleManager) GetRulesByCategory(category string) []interfaces.AuditRule {
	return rm.FilterRules(RuleFilter{Category: category})
}

// GetRulesByType returns all rules for a specific type.
func (rm *RuleManager) GetRulesByType(ruleType string) []interfaces.AuditRule {
	return rm.FilterRules(RuleFilter{Type: ruleType})
}

// GetRulesBySeverity returns all rules for a specific severity level.
func (rm *RuleManager) GetRulesBySeverity(severity string) []interfaces.AuditRule {
	return rm.FilterRules(RuleFilter{Severity: severity})
}

// GetEnabledRules returns only enabled rules.
func (rm *RuleManager) GetEnabledRules() []interfaces.AuditRule {
	enabled := true
	return rm.FilterRules(RuleFilter{Enabled: &enabled})
}

// EnableRule enables a rule by ID.
func (rm *RuleManager) EnableRule(ruleID string) error {
	for i, rule := range rm.rules {
		if rule.ID == ruleID {
			rm.rules[i].Enabled = true
			return nil
		}
	}
	return fmt.Errorf("audit rule with ID %s not found", ruleID)
}

// DisableRule disables a rule by ID.
func (rm *RuleManager) DisableRule(ruleID string) error {
	for i, rule := range rm.rules {
		if rule.ID == ruleID {
			rm.rules[i].Enabled = false
			return nil
		}
	}
	return fmt.Errorf("audit rule with ID %s not found", ruleID)
}

// ValidateRule validates an audit rule configuration.
func (rm *RuleManager) ValidateRule(rule interfaces.AuditRule) error {
	// Check required fields
	if rule.ID == "" {
		return fmt.Errorf("rule ID is required")
	}
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}
	if rule.Description == "" {
		return fmt.Errorf("rule description is required")
	}

	// Validate category
	if !rm.isValidCategory(rule.Category) {
		return fmt.Errorf("invalid category: %s", rule.Category)
	}

	// Validate type
	if !rm.isValidType(rule.Type) {
		return fmt.Errorf("invalid type: %s", rule.Type)
	}

	// Validate severity
	if !rm.isValidSeverity(rule.Severity) {
		return fmt.Errorf("invalid severity: %s", rule.Severity)
	}

	// Validate pattern if provided
	if rule.Pattern != "" {
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	// Validate file types if provided
	for _, fileType := range rule.FileTypes {
		if !strings.HasPrefix(fileType, ".") {
			return fmt.Errorf("file type must start with dot: %s", fileType)
		}
	}

	return nil
}

// RuleFilter defines criteria for filtering audit rules.
type RuleFilter struct {
	Category string // Filter by category
	Type     string // Filter by type
	Severity string // Filter by severity
	Enabled  *bool  // Filter by enabled status (nil = no filter)
	FileType string // Filter by supported file type
}

// matchesFilter checks if a rule matches the given filter criteria.
func (rm *RuleManager) matchesFilter(rule interfaces.AuditRule, filter RuleFilter) bool {
	// Check category filter
	if filter.Category != "" && rule.Category != filter.Category {
		return false
	}

	// Check type filter
	if filter.Type != "" && rule.Type != filter.Type {
		return false
	}

	// Check severity filter
	if filter.Severity != "" && rule.Severity != filter.Severity {
		return false
	}

	// Check enabled filter
	if filter.Enabled != nil && rule.Enabled != *filter.Enabled {
		return false
	}

	// Check file type filter
	if filter.FileType != "" {
		found := false
		for _, fileType := range rule.FileTypes {
			if fileType == filter.FileType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// isValidCategory checks if a category is valid.
func (rm *RuleManager) isValidCategory(category string) bool {
	validCategories := []string{
		interfaces.AuditCategorySecurity,
		interfaces.AuditCategoryQuality,
		interfaces.AuditCategoryLicense,
		interfaces.AuditCategoryPerformance,
		interfaces.AuditCategoryCompliance,
		interfaces.AuditCategoryBestPractices,
	}

	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

// isValidType checks if a type is valid.
func (rm *RuleManager) isValidType(ruleType string) bool {
	validTypes := []string{
		interfaces.AuditCategorySecurity,
		interfaces.AuditCategoryQuality,
		interfaces.AuditCategoryLicense,
		interfaces.AuditCategoryPerformance,
		interfaces.AuditCategoryCompliance,
		interfaces.AuditCategoryBestPractices,
	}

	for _, valid := range validTypes {
		if ruleType == valid {
			return true
		}
	}
	return false
}

// isValidSeverity checks if a severity is valid.
func (rm *RuleManager) isValidSeverity(severity string) bool {
	validSeverities := []string{
		interfaces.AuditSeverityCritical,
		interfaces.AuditSeverityHigh,
		interfaces.AuditSeverityMedium,
		interfaces.AuditSeverityLow,
		interfaces.AuditSeverityInfo,
	}

	for _, valid := range validSeverities {
		if severity == valid {
			return true
		}
	}
	return false
}

// getDefaultAuditRules returns the default set of audit rules.
func getDefaultAuditRules() []interfaces.AuditRule {
	return []interfaces.AuditRule{
		{
			ID:          "security-001",
			Name:        "No hardcoded secrets",
			Description: "Check for hardcoded secrets in source code",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityCritical,
			Enabled:     true,
			Pattern:     `(?i)(password|secret|key|token)\s*[:=]\s*["'][^"']+["']`,
			FileTypes:   []string{".go", ".js", ".ts", ".py", ".java", ".cs"},
		},
		{
			ID:          "security-002",
			Name:        "Secure dependencies",
			Description: "Check for known vulnerabilities in dependencies",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
		},
		{
			ID:          "security-003",
			Name:        "Secure configuration",
			Description: "Check for secure configuration practices",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
			FileTypes:   []string{".yaml", ".yml", ".json", ".env"},
		},
		{
			ID:          "quality-001",
			Name:        "Code complexity",
			Description: "Check for high cyclomatic complexity",
			Category:    interfaces.AuditCategoryQuality,
			Type:        interfaces.AuditCategoryQuality,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
			Config:      map[string]any{"max_complexity": 10},
		},
		{
			ID:          "quality-002",
			Name:        "Code duplication",
			Description: "Check for code duplication",
			Category:    interfaces.AuditCategoryQuality,
			Type:        interfaces.AuditCategoryQuality,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
			Config:      map[string]any{"min_lines": 5},
		},
		{
			ID:          "quality-003",
			Name:        "Test coverage",
			Description: "Check for adequate test coverage",
			Category:    interfaces.AuditCategoryQuality,
			Type:        interfaces.AuditCategoryQuality,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
			Config:      map[string]any{"min_coverage": 80.0},
		},
		{
			ID:          "license-001",
			Name:        "License compatibility",
			Description: "Check for license compatibility issues",
			Category:    interfaces.AuditCategoryLicense,
			Type:        interfaces.AuditCategoryLicense,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
		},
		{
			ID:          "license-002",
			Name:        "License presence",
			Description: "Check for presence of license file",
			Category:    interfaces.AuditCategoryLicense,
			Type:        interfaces.AuditCategoryLicense,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
		},
		{
			ID:          "performance-001",
			Name:        "Bundle size",
			Description: "Check for large bundle sizes",
			Category:    interfaces.AuditCategoryPerformance,
			Type:        interfaces.AuditCategoryPerformance,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
			Config:      map[string]any{"max_size_mb": 5},
		},
		{
			ID:          "performance-002",
			Name:        "Load time",
			Description: "Check for slow load times",
			Category:    interfaces.AuditCategoryPerformance,
			Type:        interfaces.AuditCategoryPerformance,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
			Config:      map[string]any{"max_load_time_ms": 2000},
		},
		{
			ID:          "best-practices-001",
			Name:        "Documentation",
			Description: "Check for adequate documentation",
			Category:    interfaces.AuditCategoryBestPractices,
			Type:        interfaces.AuditCategoryBestPractices,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
		},
		{
			ID:          "best-practices-002",
			Name:        "Code formatting",
			Description: "Check for consistent code formatting",
			Category:    interfaces.AuditCategoryBestPractices,
			Type:        interfaces.AuditCategoryBestPractices,
			Severity:    interfaces.AuditSeverityLow,
			Enabled:     true,
		},
	}
}
