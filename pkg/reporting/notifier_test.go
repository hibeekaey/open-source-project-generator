package reporting

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNotifier_NotifySecurityIssues(t *testing.T) {
	config := &NotificationConfig{
		Enabled:        true,
		Channels:       []NotificationChannel{NotificationChannelConsole},
		MinLevel:       NotificationLevelInfo,
		SecurityAlerts: true,
	}

	notifier := NewNotifier(config)

	// Create test security report
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_security_123",
		TotalIssues:    2,
		CriticalIssues: 1,
		HighIssues:     1,
		Issues: []models.SecurityIssueDetail{
			{
				PackageName:    "vulnerable-package",
				CurrentVersion: "1.0.0",
				SecurityIssue: models.SecurityIssue{
					ID:          "CVE-2023-1234",
					Severity:    "critical",
					Description: "Remote code execution vulnerability",
					FixedIn:     "1.1.0",
					ReportedAt:  time.Now(),
				},
				RecommendedFix: "1.1.0",
			},
		},
	}

	// Test notification
	err := notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Failed to send security notification: %v", err)
	}
}

func TestNotifier_NotifyVersionUpdates(t *testing.T) {
	config := &NotificationConfig{
		Enabled:      true,
		Channels:     []NotificationChannel{NotificationChannelConsole},
		MinLevel:     NotificationLevelInfo,
		UpdateAlerts: true,
	}

	notifier := NewNotifier(config)

	// Create test update report
	updateReport := &models.UpdateReport{
		GeneratedAt: time.Now(),
		ReportID:    "test_update_123",
		Type:        "version_update",
		Recommendations: []models.UpdateRecommendation{
			{
				Name:               "react",
				CurrentVersion:     "18.0.0",
				RecommendedVersion: "19.0.0",
				Priority:           "high",
				Reason:             "Security update available",
				BreakingChange:     false,
			},
			{
				Name:               "nextjs",
				CurrentVersion:     "14.0.0",
				RecommendedVersion: "15.0.0",
				Priority:           "medium",
				Reason:             "New features available",
				BreakingChange:     true,
			},
		},
	}

	// Test notification
	err := notifier.NotifyVersionUpdates(updateReport)
	if err != nil {
		t.Fatalf("Failed to send version update notification: %v", err)
	}
}

func TestNotifier_NotifyTemplateUpdates(t *testing.T) {
	config := &NotificationConfig{
		Enabled:  true,
		Channels: []NotificationChannel{NotificationChannelConsole},
		MinLevel: NotificationLevelInfo,
	}

	notifier := NewNotifier(config)

	// Create test template update report
	updateReport := &models.UpdateReport{
		GeneratedAt: time.Now(),
		ReportID:    "test_template_123",
		Type:        "template_update",
		TemplateUpdates: []models.TemplateUpdate{
			{
				TemplatePath: "templates/frontend/nextjs-app",
				Success:      true,
				UpdatedAt:    time.Now(),
				VersionChanges: map[string]string{
					"react": "18.0.0 -> 19.0.0",
				},
			},
			{
				TemplatePath: "templates/frontend/nextjs-home",
				Success:      false,
				UpdatedAt:    time.Now(),
				Error:        "Template file not found",
			},
		},
	}

	// Test notification
	err := notifier.NotifyTemplateUpdates(updateReport)
	if err != nil {
		t.Fatalf("Failed to send template update notification: %v", err)
	}
}

func TestNotifier_FileNotification(t *testing.T) {
	// Create temporary file for notifications
	tempFile, err := os.CreateTemp("", "test_notifications")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	config := &NotificationConfig{
		Enabled:        true,
		Channels:       []NotificationChannel{NotificationChannelFile},
		MinLevel:       NotificationLevelInfo,
		SecurityAlerts: true,
		OutputFile:     tempFile.Name(),
	}

	notifier := NewNotifier(config)

	// Create test security report
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_file_123",
		TotalIssues:    1,
		CriticalIssues: 1,
		HighIssues:     0,
		Issues: []models.SecurityIssueDetail{
			{
				PackageName:    "test-package",
				CurrentVersion: "1.0.0",
				SecurityIssue: models.SecurityIssue{
					ID:          "CVE-2023-5678",
					Severity:    "critical",
					Description: "Test vulnerability",
					ReportedAt:  time.Now(),
				},
			},
		},
	}

	// Test file notification
	err = notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Failed to send file notification: %v", err)
	}

	// Verify file content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read notification file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Security Vulnerabilities Detected") {
		t.Errorf("Expected notification title in file content, got: %s", contentStr)
	}

	if !strings.Contains(contentStr, "CRITICAL") {
		t.Errorf("Expected critical level in file content, got: %s", contentStr)
	}
}

func TestNotifier_MinLevelFiltering(t *testing.T) {
	config := &NotificationConfig{
		Enabled:        true,
		Channels:       []NotificationChannel{NotificationChannelConsole},
		MinLevel:       NotificationLevelCritical, // Only critical notifications
		SecurityAlerts: true,
	}

	notifier := NewNotifier(config)

	// Create test security report with low severity
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_filter_123",
		TotalIssues:    1,
		CriticalIssues: 0,
		HighIssues:     1, // This should be filtered out
		Issues: []models.SecurityIssueDetail{
			{
				PackageName:    "test-package",
				CurrentVersion: "1.0.0",
				SecurityIssue: models.SecurityIssue{
					ID:          "CVE-2023-9999",
					Severity:    "high", // Not critical, should be filtered
					Description: "Test vulnerability",
					ReportedAt:  time.Now(),
				},
			},
		},
	}

	// This should not send notification due to min level filtering
	err := notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Unexpected error from filtered notification: %v", err)
	}

	// Now test with critical issue
	securityReport.CriticalIssues = 1
	securityReport.Issues[0].SecurityIssue.Severity = "critical"

	// This should send notification
	err = notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Failed to send critical notification: %v", err)
	}
}

func TestNotifier_DisabledNotifications(t *testing.T) {
	config := &NotificationConfig{
		Enabled:        false, // Disabled
		Channels:       []NotificationChannel{NotificationChannelConsole},
		MinLevel:       NotificationLevelInfo,
		SecurityAlerts: true,
	}

	notifier := NewNotifier(config)

	// Create test security report
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_disabled_123",
		TotalIssues:    1,
		CriticalIssues: 1,
		HighIssues:     0,
	}

	// Should not send notification when disabled
	err := notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Unexpected error from disabled notifier: %v", err)
	}
}

func TestNotifier_DisabledSecurityAlerts(t *testing.T) {
	config := &NotificationConfig{
		Enabled:        true,
		Channels:       []NotificationChannel{NotificationChannelConsole},
		MinLevel:       NotificationLevelInfo,
		SecurityAlerts: false, // Security alerts disabled
	}

	notifier := NewNotifier(config)

	// Create test security report
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_no_security_123",
		TotalIssues:    1,
		CriticalIssues: 1,
		HighIssues:     0,
	}

	// Should not send notification when security alerts are disabled
	err := notifier.NotifySecurityIssues(securityReport)
	if err != nil {
		t.Fatalf("Unexpected error from disabled security alerts: %v", err)
	}
}

func TestNotifier_FormatSecurityMessage(t *testing.T) {
	config := GetDefaultConfig()
	notifier := NewNotifier(config)

	// Create test security report
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       "test_format_123",
		TotalIssues:    2,
		CriticalIssues: 1,
		HighIssues:     1,
		Issues: []models.SecurityIssueDetail{
			{
				PackageName:    "package1",
				CurrentVersion: "1.0.0",
				SecurityIssue: models.SecurityIssue{
					ID:          "CVE-2023-1111",
					Severity:    "critical",
					Description: "Critical vulnerability",
					ReportedAt:  time.Now(),
				},
			},
			{
				PackageName:    "package2",
				CurrentVersion: "2.0.0",
				SecurityIssue: models.SecurityIssue{
					ID:          "CVE-2023-2222",
					Severity:    "high",
					Description: "High severity vulnerability",
					ReportedAt:  time.Now(),
				},
			},
		},
	}

	message := notifier.formatSecurityMessage(securityReport)

	// Verify message content
	if !strings.Contains(message, "Total Issues: 2") {
		t.Errorf("Expected total issues count in message")
	}

	if !strings.Contains(message, "Critical: 1") {
		t.Errorf("Expected critical issues count in message")
	}

	if !strings.Contains(message, "High: 1") {
		t.Errorf("Expected high issues count in message")
	}

	if !strings.Contains(message, "package1") {
		t.Errorf("Expected package1 in affected packages")
	}

	if !strings.Contains(message, "package2") {
		t.Errorf("Expected package2 in affected packages")
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if !config.Enabled {
		t.Errorf("Expected default config to be enabled")
	}

	if !config.SecurityAlerts {
		t.Errorf("Expected default config to have security alerts enabled")
	}

	if !config.UpdateAlerts {
		t.Errorf("Expected default config to have update alerts enabled")
	}

	if config.MinLevel != NotificationLevelInfo {
		t.Errorf("Expected default min level to be info, got %s", config.MinLevel)
	}

	if len(config.Channels) != 1 || config.Channels[0] != NotificationChannelConsole {
		t.Errorf("Expected default channel to be console only")
	}
}
