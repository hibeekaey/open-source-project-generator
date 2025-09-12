package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// AuditEvent represents a single audit log entry
type AuditEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	EventType string                 `json:"event_type"`
	User      string                 `json:"user,omitempty"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	OldValue  string                 `json:"old_value,omitempty"`
	NewValue  string                 `json:"new_value,omitempty"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AuditTrail manages audit logging for version management operations
type AuditTrail struct {
	logFile string
}

// NewAuditTrail creates a new audit trail logger
func NewAuditTrail(logFile string) *AuditTrail {
	return &AuditTrail{
		logFile: logFile,
	}
}

// LogVersionUpdate logs a version update operation
func (a *AuditTrail) LogVersionUpdate(packageName, oldVersion, newVersion string, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "version_update",
		Action:    "update_version",
		Resource:  packageName,
		OldValue:  oldVersion,
		NewValue:  newVersion,
		Success:   success,
		Metadata: map[string]interface{}{
			"package_name": packageName,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogTemplateUpdate logs a template update operation
func (a *AuditTrail) LogTemplateUpdate(templatePath string, versionChanges map[string]string, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "template_update",
		Action:    "update_template",
		Resource:  templatePath,
		Success:   success,
		Metadata: map[string]interface{}{
			"template_path":   templatePath,
			"version_changes": versionChanges,
			"changes_count":   len(versionChanges),
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogSecurityScan logs a security scan operation
func (a *AuditTrail) LogSecurityScan(packagesScanned int, issuesFound int, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "security_scan",
		Action:    "scan_packages",
		Resource:  "security_database",
		Success:   success,
		Metadata: map[string]interface{}{
			"packages_scanned": packagesScanned,
			"issues_found":     issuesFound,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogVersionCheck logs a version check operation
func (a *AuditTrail) LogVersionCheck(packagesChecked int, updatesFound int, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "version_check",
		Action:    "check_versions",
		Resource:  "version_registries",
		Success:   success,
		Metadata: map[string]interface{}{
			"packages_checked": packagesChecked,
			"updates_found":    updatesFound,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogConfigurationChange logs configuration changes
func (a *AuditTrail) LogConfigurationChange(configKey, oldValue, newValue string, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "configuration_change",
		Action:    "update_config",
		Resource:  configKey,
		OldValue:  oldValue,
		NewValue:  newValue,
		Success:   success,
		Metadata: map[string]interface{}{
			"config_key": configKey,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogReportGeneration logs report generation events
func (a *AuditTrail) LogReportGeneration(reportType, reportID string, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "report_generation",
		Action:    "generate_report",
		Resource:  reportType,
		Success:   success,
		Metadata: map[string]interface{}{
			"report_type": reportType,
			"report_id":   reportID,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// LogNotification logs notification events
func (a *AuditTrail) LogNotification(notificationType, channel string, success bool, err error) error {
	event := AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now(),
		EventType: "notification",
		Action:    "send_notification",
		Resource:  channel,
		Success:   success,
		Metadata: map[string]interface{}{
			"notification_type": notificationType,
			"channel":           channel,
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	return a.writeEvent(event)
}

// GetAuditHistory retrieves audit events within a time range
func (a *AuditTrail) GetAuditHistory(since, until time.Time, eventType string) ([]AuditEvent, error) {
	var events []AuditEvent

	if _, err := os.Stat(a.logFile); os.IsNotExist(err) {
		return events, nil // No audit log exists yet
	}

	file, err := os.Open(a.logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var event AuditEvent
		if err := decoder.Decode(&event); err != nil {
			continue // Skip malformed entries
		}

		// Filter by time range
		if event.Timestamp.Before(since) || event.Timestamp.After(until) {
			continue
		}

		// Filter by event type if specified
		if eventType != "" && event.EventType != eventType {
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

// GetAuditSummary provides summary statistics for audit events
func (a *AuditTrail) GetAuditSummary(since time.Time) (*models.AuditSummary, error) {
	events, err := a.GetAuditHistory(since, time.Now(), "")
	if err != nil {
		return nil, fmt.Errorf("failed to get audit history: %w", err)
	}

	summary := &models.AuditSummary{
		Period:      fmt.Sprintf("%s to %s", since.Format("2006-01-02"), time.Now().Format("2006-01-02")),
		TotalEvents: len(events),
		EventTypes:  make(map[string]int),
		Actions:     make(map[string]int),
		SuccessRate: 0.0,
	}

	successCount := 0
	for _, event := range events {
		summary.EventTypes[event.EventType]++
		summary.Actions[event.Action]++
		if event.Success {
			successCount++
		}
	}

	if len(events) > 0 {
		summary.SuccessRate = float64(successCount) / float64(len(events)) * 100
	}

	return summary, nil
}

// writeEvent writes an audit event to the log file
func (a *AuditTrail) writeEvent(event AuditEvent) error {
	// Ensure directory exists
	dir := filepath.Dir(a.logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create audit log directory: %w", err)
	}

	// Open file for appending
	file, err := os.OpenFile(a.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit log file: %w", err)
	}
	defer file.Close()

	// Write event as JSON line
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(event); err != nil {
		return fmt.Errorf("failed to write audit event: %w", err)
	}

	return nil
}

// generateEventID creates a unique event identifier
func (a *AuditTrail) generateEventID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

// RotateLog rotates the audit log file if it exceeds size limit
func (a *AuditTrail) RotateLog(maxSizeBytes int64) error {
	if _, err := os.Stat(a.logFile); os.IsNotExist(err) {
		return nil // No log file to rotate
	}

	info, err := os.Stat(a.logFile)
	if err != nil {
		return fmt.Errorf("failed to stat audit log: %w", err)
	}

	if info.Size() < maxSizeBytes {
		return nil // File is not large enough to rotate
	}

	// Create rotated filename with timestamp
	rotatedFile := fmt.Sprintf("%s.%s", a.logFile, time.Now().Format("20060102-150405"))

	if err := os.Rename(a.logFile, rotatedFile); err != nil {
		return fmt.Errorf("failed to rotate audit log: %w", err)
	}

	return nil
}

// CleanupOldLogs removes audit logs older than the specified duration
func (a *AuditTrail) CleanupOldLogs(maxAge time.Duration) error {
	dir := filepath.Dir(a.logFile)
	baseName := filepath.Base(a.logFile)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read audit log directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if this is a rotated log file
		if !strings.HasPrefix(entry.Name(), baseName+".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			logPath := filepath.Join(dir, entry.Name())
			if err := os.Remove(logPath); err != nil {
				fmt.Printf("Warning: failed to remove old audit log %s: %v\n", logPath, err)
			}
		}
	}

	return nil
}
