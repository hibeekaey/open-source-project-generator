package audit

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRuleManager(t *testing.T) {
	rm := NewRuleManager()

	assert.NotNil(t, rm)
	rules := rm.GetRules()
	assert.NotEmpty(t, rules)

	// Should have default rules
	assert.True(t, len(rules) > 0)

	// Check that we have rules from different categories
	categories := make(map[string]bool)
	for _, rule := range rules {
		categories[rule.Category] = true
	}

	assert.True(t, categories[interfaces.AuditCategorySecurity])
	assert.True(t, categories[interfaces.AuditCategoryQuality])
	assert.True(t, categories[interfaces.AuditCategoryLicense])
	assert.True(t, categories[interfaces.AuditCategoryPerformance])
}

func TestNewRuleManagerWithRules(t *testing.T) {
	customRules := []interfaces.AuditRule{
		{
			ID:          "test-001",
			Name:        "Test Rule",
			Description: "A test rule",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
		},
	}

	rm := NewRuleManagerWithRules(customRules)

	assert.NotNil(t, rm)
	rules := rm.GetRules()
	assert.Len(t, rules, 1)
	assert.Equal(t, "test-001", rules[0].ID)
}

func TestRuleManager_GetRules(t *testing.T) {
	rm := NewRuleManager()

	rules1 := rm.GetRules()
	rules2 := rm.GetRules()

	// Should return copies, not the same slice
	assert.NotSame(t, &rules1, &rules2)
	assert.Equal(t, len(rules1), len(rules2))
}

func TestRuleManager_SetRules(t *testing.T) {
	rm := NewRuleManager()

	newRules := []interfaces.AuditRule{
		{
			ID:          "new-001",
			Name:        "New Rule",
			Description: "A new rule",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
		},
	}

	err := rm.SetRules(newRules)
	require.NoError(t, err)

	rules := rm.GetRules()
	assert.Len(t, rules, 1)
	assert.Equal(t, "new-001", rules[0].ID)
}

func TestRuleManager_SetRules_InvalidRule(t *testing.T) {
	rm := NewRuleManager()

	invalidRules := []interfaces.AuditRule{
		{
			ID:          "", // Invalid: empty ID
			Name:        "Invalid Rule",
			Description: "An invalid rule",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
		},
	}

	err := rm.SetRules(invalidRules)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rule ID is required")
}

func TestRuleManager_AddRule(t *testing.T) {
	rm := NewRuleManager()
	initialCount := len(rm.GetRules())

	newRule := interfaces.AuditRule{
		ID:          "test-add-001",
		Name:        "Test Add Rule",
		Description: "A test rule for adding",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     true,
	}

	err := rm.AddRule(newRule)
	require.NoError(t, err)

	rules := rm.GetRules()
	assert.Len(t, rules, initialCount+1)

	// Find the added rule
	found := false
	for _, rule := range rules {
		if rule.ID == "test-add-001" {
			found = true
			assert.Equal(t, newRule.Name, rule.Name)
			break
		}
	}
	assert.True(t, found)
}

func TestRuleManager_AddRule_Replace(t *testing.T) {
	rm := NewRuleManager()
	initialCount := len(rm.GetRules())

	// Add a rule
	rule1 := interfaces.AuditRule{
		ID:          "test-replace-001",
		Name:        "Original Rule",
		Description: "Original description",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     true,
	}

	err := rm.AddRule(rule1)
	require.NoError(t, err)

	// Replace with same ID
	rule2 := interfaces.AuditRule{
		ID:          "test-replace-001",
		Name:        "Updated Rule",
		Description: "Updated description",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityMedium,
		Enabled:     false,
	}

	err = rm.AddRule(rule2)
	require.NoError(t, err)

	rules := rm.GetRules()
	assert.Len(t, rules, initialCount+1) // Should not increase count

	// Find the updated rule
	found := false
	for _, rule := range rules {
		if rule.ID == "test-replace-001" {
			found = true
			assert.Equal(t, "Updated Rule", rule.Name)
			assert.Equal(t, "Updated description", rule.Description)
			assert.Equal(t, interfaces.AuditSeverityMedium, rule.Severity)
			assert.False(t, rule.Enabled)
			break
		}
	}
	assert.True(t, found)
}

func TestRuleManager_AddRule_Invalid(t *testing.T) {
	rm := NewRuleManager()

	invalidRule := interfaces.AuditRule{
		ID:          "test-invalid",
		Name:        "", // Invalid: empty name
		Description: "Invalid rule",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     true,
	}

	err := rm.AddRule(invalidRule)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rule name is required")
}

func TestRuleManager_RemoveRule(t *testing.T) {
	rm := NewRuleManager()

	// Add a rule to remove
	testRule := interfaces.AuditRule{
		ID:          "test-remove-001",
		Name:        "Test Remove Rule",
		Description: "A test rule for removal",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     true,
	}

	err := rm.AddRule(testRule)
	require.NoError(t, err)

	initialCount := len(rm.GetRules())

	// Remove the rule
	err = rm.RemoveRule("test-remove-001")
	require.NoError(t, err)

	rules := rm.GetRules()
	assert.Len(t, rules, initialCount-1)

	// Verify rule is gone
	for _, rule := range rules {
		assert.NotEqual(t, "test-remove-001", rule.ID)
	}
}

func TestRuleManager_RemoveRule_NotFound(t *testing.T) {
	rm := NewRuleManager()

	err := rm.RemoveRule("non-existent-rule")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "audit rule with ID non-existent-rule not found")
}

func TestRuleManager_GetRule(t *testing.T) {
	rm := NewRuleManager()

	// Get a default rule (we know security-001 exists)
	rule, err := rm.GetRule("security-001")
	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "security-001", rule.ID)
	assert.Equal(t, "No hardcoded secrets", rule.Name)
}

func TestRuleManager_GetRule_NotFound(t *testing.T) {
	rm := NewRuleManager()

	rule, err := rm.GetRule("non-existent-rule")
	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "audit rule with ID non-existent-rule not found")
}

func TestRuleManager_FilterRules(t *testing.T) {
	rm := NewRuleManager()

	// Filter by category
	securityRules := rm.FilterRules(RuleFilter{Category: interfaces.AuditCategorySecurity})
	assert.NotEmpty(t, securityRules)
	for _, rule := range securityRules {
		assert.Equal(t, interfaces.AuditCategorySecurity, rule.Category)
	}

	// Filter by severity
	criticalRules := rm.FilterRules(RuleFilter{Severity: interfaces.AuditSeverityCritical})
	assert.NotEmpty(t, criticalRules)
	for _, rule := range criticalRules {
		assert.Equal(t, interfaces.AuditSeverityCritical, rule.Severity)
	}

	// Filter by enabled status
	enabled := true
	enabledRules := rm.FilterRules(RuleFilter{Enabled: &enabled})
	assert.NotEmpty(t, enabledRules)
	for _, rule := range enabledRules {
		assert.True(t, rule.Enabled)
	}
}

func TestRuleManager_GetRulesByCategory(t *testing.T) {
	rm := NewRuleManager()

	securityRules := rm.GetRulesByCategory(interfaces.AuditCategorySecurity)
	assert.NotEmpty(t, securityRules)
	for _, rule := range securityRules {
		assert.Equal(t, interfaces.AuditCategorySecurity, rule.Category)
	}

	qualityRules := rm.GetRulesByCategory(interfaces.AuditCategoryQuality)
	assert.NotEmpty(t, qualityRules)
	for _, rule := range qualityRules {
		assert.Equal(t, interfaces.AuditCategoryQuality, rule.Category)
	}
}

func TestRuleManager_GetRulesByType(t *testing.T) {
	rm := NewRuleManager()

	securityRules := rm.GetRulesByType(interfaces.AuditCategorySecurity)
	assert.NotEmpty(t, securityRules)
	for _, rule := range securityRules {
		assert.Equal(t, interfaces.AuditCategorySecurity, rule.Type)
	}
}

func TestRuleManager_GetRulesBySeverity(t *testing.T) {
	rm := NewRuleManager()

	highRules := rm.GetRulesBySeverity(interfaces.AuditSeverityHigh)
	assert.NotEmpty(t, highRules)
	for _, rule := range highRules {
		assert.Equal(t, interfaces.AuditSeverityHigh, rule.Severity)
	}
}

func TestRuleManager_GetEnabledRules(t *testing.T) {
	rm := NewRuleManager()

	// Add a disabled rule
	disabledRule := interfaces.AuditRule{
		ID:          "test-disabled",
		Name:        "Disabled Rule",
		Description: "A disabled rule",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     false,
	}

	err := rm.AddRule(disabledRule)
	require.NoError(t, err)

	enabledRules := rm.GetEnabledRules()

	// All returned rules should be enabled
	for _, rule := range enabledRules {
		assert.True(t, rule.Enabled)
	}

	// The disabled rule should not be in the results
	for _, rule := range enabledRules {
		assert.NotEqual(t, "test-disabled", rule.ID)
	}
}

func TestRuleManager_EnableRule(t *testing.T) {
	rm := NewRuleManager()

	// Add a disabled rule
	disabledRule := interfaces.AuditRule{
		ID:          "test-enable",
		Name:        "Test Enable Rule",
		Description: "A rule to test enabling",
		Category:    interfaces.AuditCategorySecurity,
		Type:        interfaces.AuditCategorySecurity,
		Severity:    interfaces.AuditSeverityHigh,
		Enabled:     false,
	}

	err := rm.AddRule(disabledRule)
	require.NoError(t, err)

	// Enable the rule
	err = rm.EnableRule("test-enable")
	require.NoError(t, err)

	// Verify it's enabled
	rule, err := rm.GetRule("test-enable")
	require.NoError(t, err)
	assert.True(t, rule.Enabled)
}

func TestRuleManager_EnableRule_NotFound(t *testing.T) {
	rm := NewRuleManager()

	err := rm.EnableRule("non-existent-rule")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "audit rule with ID non-existent-rule not found")
}

func TestRuleManager_DisableRule(t *testing.T) {
	rm := NewRuleManager()

	// Disable an existing rule (we know security-001 exists and is enabled)
	err := rm.DisableRule("security-001")
	require.NoError(t, err)

	// Verify it's disabled
	rule, err := rm.GetRule("security-001")
	require.NoError(t, err)
	assert.False(t, rule.Enabled)
}

func TestRuleManager_DisableRule_NotFound(t *testing.T) {
	rm := NewRuleManager()

	err := rm.DisableRule("non-existent-rule")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "audit rule with ID non-existent-rule not found")
}

func TestRuleManager_ValidateRule(t *testing.T) {
	rm := NewRuleManager()

	tests := []struct {
		name    string
		rule    interfaces.AuditRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid rule",
			rule: interfaces.AuditRule{
				ID:          "test-valid",
				Name:        "Valid Rule",
				Description: "A valid rule",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			rule: interfaces.AuditRule{
				Name:        "Missing ID Rule",
				Description: "A rule missing ID",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
			},
			wantErr: true,
			errMsg:  "rule ID is required",
		},
		{
			name: "missing name",
			rule: interfaces.AuditRule{
				ID:          "test-missing-name",
				Description: "A rule missing name",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
			},
			wantErr: true,
			errMsg:  "rule name is required",
		},
		{
			name: "missing description",
			rule: interfaces.AuditRule{
				ID:       "test-missing-desc",
				Name:     "Missing Description Rule",
				Category: interfaces.AuditCategorySecurity,
				Type:     interfaces.AuditCategorySecurity,
				Severity: interfaces.AuditSeverityHigh,
				Enabled:  true,
			},
			wantErr: true,
			errMsg:  "rule description is required",
		},
		{
			name: "invalid category",
			rule: interfaces.AuditRule{
				ID:          "test-invalid-category",
				Name:        "Invalid Category Rule",
				Description: "A rule with invalid category",
				Category:    "invalid-category",
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
			},
			wantErr: true,
			errMsg:  "invalid category",
		},
		{
			name: "invalid type",
			rule: interfaces.AuditRule{
				ID:          "test-invalid-type",
				Name:        "Invalid Type Rule",
				Description: "A rule with invalid type",
				Category:    interfaces.AuditCategorySecurity,
				Type:        "invalid-type",
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
			},
			wantErr: true,
			errMsg:  "invalid type",
		},
		{
			name: "invalid severity",
			rule: interfaces.AuditRule{
				ID:          "test-invalid-severity",
				Name:        "Invalid Severity Rule",
				Description: "A rule with invalid severity",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    "invalid-severity",
				Enabled:     true,
			},
			wantErr: true,
			errMsg:  "invalid severity",
		},
		{
			name: "invalid regex pattern",
			rule: interfaces.AuditRule{
				ID:          "test-invalid-pattern",
				Name:        "Invalid Pattern Rule",
				Description: "A rule with invalid regex pattern",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
				Pattern:     "[invalid-regex",
			},
			wantErr: true,
			errMsg:  "invalid regex pattern",
		},
		{
			name: "invalid file type",
			rule: interfaces.AuditRule{
				ID:          "test-invalid-filetype",
				Name:        "Invalid File Type Rule",
				Description: "A rule with invalid file type",
				Category:    interfaces.AuditCategorySecurity,
				Type:        interfaces.AuditCategorySecurity,
				Severity:    interfaces.AuditSeverityHigh,
				Enabled:     true,
				FileTypes:   []string{"go"}, // Should start with dot
			},
			wantErr: true,
			errMsg:  "file type must start with dot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rm.ValidateRule(tt.rule)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRuleFilter_FileType(t *testing.T) {
	rm := NewRuleManager()

	// Filter by file type
	goRules := rm.FilterRules(RuleFilter{FileType: ".go"})

	// All returned rules should support .go files
	for _, rule := range goRules {
		found := false
		for _, fileType := range rule.FileTypes {
			if fileType == ".go" {
				found = true
				break
			}
		}
		// Only check if the rule has file types specified
		if len(rule.FileTypes) > 0 {
			assert.True(t, found, "Rule %s should support .go files", rule.ID)
		}
	}
}

func TestGetDefaultAuditRules(t *testing.T) {
	rules := getDefaultAuditRules()

	assert.NotEmpty(t, rules)

	// Check that we have rules from all major categories
	categories := make(map[string]int)
	for _, rule := range rules {
		categories[rule.Category]++
	}

	assert.Greater(t, categories[interfaces.AuditCategorySecurity], 0)
	assert.Greater(t, categories[interfaces.AuditCategoryQuality], 0)
	assert.Greater(t, categories[interfaces.AuditCategoryLicense], 0)
	assert.Greater(t, categories[interfaces.AuditCategoryPerformance], 0)

	// Verify all rules are valid
	rm := NewRuleManager()
	for _, rule := range rules {
		err := rm.ValidateRule(rule)
		assert.NoError(t, err, "Default rule %s should be valid", rule.ID)
	}

	// Check for specific expected rules
	ruleIDs := make(map[string]bool)
	for _, rule := range rules {
		ruleIDs[rule.ID] = true
	}

	expectedRules := []string{
		"security-001", "security-002", "security-003",
		"quality-001", "quality-002", "quality-003",
		"license-001", "license-002",
		"performance-001", "performance-002",
		"best-practices-001", "best-practices-002",
	}

	for _, expectedID := range expectedRules {
		assert.True(t, ruleIDs[expectedID], "Expected rule %s should exist", expectedID)
	}
}
