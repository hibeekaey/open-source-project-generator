package reporting

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// NotificationLevel represents the severity level of notifications
type NotificationLevel string

const (
	NotificationLevelInfo     NotificationLevel = "info"
	NotificationLevelWarning  NotificationLevel = "warning"
	NotificationLevelCritical NotificationLevel = "critical"
)

// NotificationChannel represents different notification delivery methods
type NotificationChannel string

const (
	NotificationChannelConsole NotificationChannel = "console"
	NotificationChannelFile    NotificationChannel = "file"
	NotificationChannelEmail   NotificationChannel = "email"
	NotificationChannelSlack   NotificationChannel = "slack"
)

// Notification represents a notification message
type Notification struct {
	ID        string            `json:"id"`
	Level     NotificationLevel `json:"level"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// NotificationConfig configures notification behavior
type NotificationConfig struct {
	Enabled        bool                  `json:"enabled"`
	Channels       []NotificationChannel `json:"channels"`
	MinLevel       NotificationLevel     `json:"min_level"`
	SecurityAlerts bool                  `json:"security_alerts"`
	UpdateAlerts   bool                  `json:"update_alerts"`
	OutputFile     string                `json:"output_file,omitempty"`
	EmailConfig    *EmailConfig          `json:"email_config,omitempty"`
	SlackConfig    *SlackConfig          `json:"slack_config,omitempty"`
}

// EmailConfig configures email notifications
type EmailConfig struct {
	SMTPServer string   `json:"smtp_server"`
	Port       int      `json:"port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	From       string   `json:"from"`
	To         []string `json:"to"`
}

// SlackConfig configures Slack notifications
type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
}

// Notifier handles sending notifications for version management events
type Notifier struct {
	config *NotificationConfig
	logger *log.Logger
}

// NewNotifier creates a new notification system
func NewNotifier(config *NotificationConfig) *Notifier {
	logger := log.New(os.Stdout, "[NOTIFIER] ", log.LstdFlags)

	return &Notifier{
		config: config,
		logger: logger,
	}
}

// NotifySecurityIssues sends notifications for security vulnerabilities
func (n *Notifier) NotifySecurityIssues(report *models.SecurityReport) error {
	if !n.config.Enabled || !n.config.SecurityAlerts {
		return nil
	}

	// Create notification based on severity
	level := NotificationLevelInfo
	if report.CriticalIssues > 0 {
		level = NotificationLevelCritical
	} else if report.HighIssues > 0 {
		level = NotificationLevelWarning
	}

	notification := &Notification{
		ID:        fmt.Sprintf("security_%s", report.ReportID),
		Level:     level,
		Title:     "Security Vulnerabilities Detected",
		Message:   n.formatSecurityMessage(report),
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"report_id":       report.ReportID,
			"total_issues":    fmt.Sprintf("%d", report.TotalIssues),
			"critical_issues": fmt.Sprintf("%d", report.CriticalIssues),
			"high_issues":     fmt.Sprintf("%d", report.HighIssues),
		},
	}

	return n.sendNotification(notification)
}

// NotifyVersionUpdates sends notifications for available version updates
func (n *Notifier) NotifyVersionUpdates(report *models.UpdateReport) error {
	if !n.config.Enabled || !n.config.UpdateAlerts {
		return nil
	}

	// Determine notification level based on recommendations
	level := NotificationLevelInfo
	criticalCount := 0
	highCount := 0

	for _, rec := range report.Recommendations {
		if rec.Priority == "critical" {
			criticalCount++
		} else if rec.Priority == "high" {
			highCount++
		}
	}

	if criticalCount > 0 {
		level = NotificationLevelCritical
	} else if highCount > 0 {
		level = NotificationLevelWarning
	}

	notification := &Notification{
		ID:        fmt.Sprintf("updates_%s", report.ReportID),
		Level:     level,
		Title:     "Version Updates Available",
		Message:   n.formatUpdateMessage(report),
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"report_id":        report.ReportID,
			"total_updates":    fmt.Sprintf("%d", len(report.Recommendations)),
			"critical_updates": fmt.Sprintf("%d", criticalCount),
			"high_updates":     fmt.Sprintf("%d", highCount),
		},
	}

	return n.sendNotification(notification)
}

// NotifyTemplateUpdates sends notifications for template update results
func (n *Notifier) NotifyTemplateUpdates(report *models.UpdateReport) error {
	if !n.config.Enabled {
		return nil
	}

	// Count successful and failed updates
	successCount := 0
	failureCount := 0
	for _, update := range report.TemplateUpdates {
		if update.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	level := NotificationLevelInfo
	if failureCount > 0 {
		level = NotificationLevelWarning
	}

	notification := &Notification{
		ID:        fmt.Sprintf("templates_%s", report.ReportID),
		Level:     level,
		Title:     "Template Updates Completed",
		Message:   n.formatTemplateUpdateMessage(report, successCount, failureCount),
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"report_id":     report.ReportID,
			"total_updates": fmt.Sprintf("%d", len(report.TemplateUpdates)),
			"successful":    fmt.Sprintf("%d", successCount),
			"failed":        fmt.Sprintf("%d", failureCount),
		},
	}

	return n.sendNotification(notification)
}

// sendNotification delivers a notification through configured channels
func (n *Notifier) sendNotification(notification *Notification) error {
	// Check if notification level meets minimum threshold
	if !n.shouldSendNotification(notification.Level) {
		return nil
	}

	var errors []string

	for _, channel := range n.config.Channels {
		if err := n.sendToChannel(notification, channel); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", channel, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification delivery failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// sendToChannel sends notification to a specific channel
func (n *Notifier) sendToChannel(notification *Notification, channel NotificationChannel) error {
	switch channel {
	case NotificationChannelConsole:
		return n.sendToConsole(notification)
	case NotificationChannelFile:
		return n.sendToFile(notification)
	case NotificationChannelEmail:
		return n.sendToEmail(notification)
	case NotificationChannelSlack:
		return n.sendToSlack(notification)
	default:
		return fmt.Errorf("unsupported notification channel: %s", channel)
	}
}

// sendToConsole outputs notification to console
func (n *Notifier) sendToConsole(notification *Notification) error {
	icon := "ℹ️"
	switch notification.Level {
	case NotificationLevelWarning:
		icon = "⚠️"
	case NotificationLevelCritical:
		icon = "🚨"
	}

	fmt.Printf("\n%s %s\n", icon, notification.Title)
	fmt.Printf("Time: %s\n", notification.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Level: %s\n", notification.Level)
	fmt.Printf("Message:\n%s\n", notification.Message)

	if len(notification.Metadata) > 0 {
		fmt.Printf("Details:\n")
		for key, value := range notification.Metadata {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	fmt.Println()

	return nil
}

// sendToFile writes notification to file
func (n *Notifier) sendToFile(notification *Notification) error {
	if n.config.OutputFile == "" {
		return fmt.Errorf("no output file configured for file notifications")
	}

	file, err := os.OpenFile(n.config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open notification file: %w", err)
	}
	defer file.Close()

	entry := fmt.Sprintf("[%s] %s - %s: %s\n",
		notification.Timestamp.Format("2006-01-02 15:04:05"),
		strings.ToUpper(string(notification.Level)),
		notification.Title,
		notification.Message)

	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("failed to write notification to file: %w", err)
	}

	return nil
}

// sendToEmail sends notification via email (placeholder implementation)
func (n *Notifier) sendToEmail(notification *Notification) error {
	if n.config.EmailConfig == nil {
		return fmt.Errorf("email configuration not provided")
	}

	// This is a placeholder implementation
	// In a real implementation, you would use an SMTP library to send emails
	n.logger.Printf("EMAIL NOTIFICATION: %s - %s", notification.Title, notification.Message)
	return nil
}

// sendToSlack sends notification to Slack (placeholder implementation)
func (n *Notifier) sendToSlack(notification *Notification) error {
	if n.config.SlackConfig == nil {
		return fmt.Errorf("Slack configuration not provided")
	}

	// This is a placeholder implementation
	// In a real implementation, you would use Slack's webhook API
	n.logger.Printf("SLACK NOTIFICATION: %s - %s", notification.Title, notification.Message)
	return nil
}

// shouldSendNotification checks if notification meets minimum level threshold
func (n *Notifier) shouldSendNotification(level NotificationLevel) bool {
	levelPriority := map[NotificationLevel]int{
		NotificationLevelInfo:     1,
		NotificationLevelWarning:  2,
		NotificationLevelCritical: 3,
	}

	return levelPriority[level] >= levelPriority[n.config.MinLevel]
}

// formatSecurityMessage formats security issue notification message
func (n *Notifier) formatSecurityMessage(report *models.SecurityReport) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("Security scan completed at %s\n\n",
		report.GeneratedAt.Format("2006-01-02 15:04:05")))

	message.WriteString(fmt.Sprintf("Total Issues: %d\n", report.TotalIssues))
	message.WriteString(fmt.Sprintf("Critical: %d\n", report.CriticalIssues))
	message.WriteString(fmt.Sprintf("High: %d\n", report.HighIssues))

	if len(report.Issues) > 0 {
		message.WriteString("\nAffected Packages:\n")
		for i, issue := range report.Issues {
			if i >= 5 { // Limit to first 5 issues in notification
				message.WriteString(fmt.Sprintf("... and %d more issues\n", len(report.Issues)-5))
				break
			}
			message.WriteString(fmt.Sprintf("- %s (%s): %s\n",
				issue.PackageName, issue.CurrentVersion, issue.SecurityIssue.Description))
		}
	}

	message.WriteString("\nPlease review the full security report and update affected packages.")

	return message.String()
}

// formatUpdateMessage formats version update notification message
func (n *Notifier) formatUpdateMessage(report *models.UpdateReport) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("Version check completed at %s\n\n",
		report.GeneratedAt.Format("2006-01-02 15:04:05")))

	message.WriteString(fmt.Sprintf("Total Updates Available: %d\n", len(report.Recommendations)))

	// Count by priority
	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, rec := range report.Recommendations {
		switch rec.Priority {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low":
			lowCount++
		}
	}

	if criticalCount > 0 {
		message.WriteString(fmt.Sprintf("Critical: %d\n", criticalCount))
	}
	if highCount > 0 {
		message.WriteString(fmt.Sprintf("High: %d\n", highCount))
	}
	if mediumCount > 0 {
		message.WriteString(fmt.Sprintf("Medium: %d\n", mediumCount))
	}
	if lowCount > 0 {
		message.WriteString(fmt.Sprintf("Low: %d\n", lowCount))
	}

	if len(report.Recommendations) > 0 {
		message.WriteString("\nTop Priority Updates:\n")
		count := 0
		for _, rec := range report.Recommendations {
			if rec.Priority == "critical" || rec.Priority == "high" {
				if count >= 5 { // Limit to first 5 high priority updates
					break
				}
				message.WriteString(fmt.Sprintf("- %s: %s → %s (%s)\n",
					rec.Name, rec.CurrentVersion, rec.RecommendedVersion, rec.Priority))
				count++
			}
		}
	}

	message.WriteString("\nRun 'generator versions update' to apply updates.")

	return message.String()
}

// formatTemplateUpdateMessage formats template update notification message
func (n *Notifier) formatTemplateUpdateMessage(report *models.UpdateReport, successCount, failureCount int) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("Template update completed at %s\n\n",
		report.GeneratedAt.Format("2006-01-02 15:04:05")))

	message.WriteString(fmt.Sprintf("Total Templates: %d\n", len(report.TemplateUpdates)))
	message.WriteString(fmt.Sprintf("Successful: %d\n", successCount))
	message.WriteString(fmt.Sprintf("Failed: %d\n", failureCount))

	if failureCount > 0 {
		message.WriteString("\nFailed Updates:\n")
		for _, update := range report.TemplateUpdates {
			if !update.Success {
				message.WriteString(fmt.Sprintf("- %s: %s\n", update.TemplatePath, update.Error))
			}
		}
	}

	return message.String()
}

// GetDefaultConfig returns a default notification configuration
func GetDefaultConfig() *NotificationConfig {
	return &NotificationConfig{
		Enabled:        true,
		Channels:       []NotificationChannel{NotificationChannelConsole},
		MinLevel:       NotificationLevelInfo,
		SecurityAlerts: true,
		UpdateAlerts:   true,
	}
}
