package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AuditLogEntry represents a single audit log entry
type AuditLogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	User      string                 `json:"user,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// AuditLogSummary provides summary statistics for audit logs
type AuditLogSummary struct {
	TotalEntries   int              `json:"total_entries"`
	SuccessCount   int              `json:"success_count"`
	ErrorCount     int              `json:"error_count"`
	ActionCounts   map[string]int   `json:"action_counts"`
	ResourceCounts map[string]int   `json:"resource_counts"`
	TimeRange      *AuditTimeRange  `json:"time_range"`
	TopErrors      []string         `json:"top_errors"`
	RecentActivity []*AuditLogEntry `json:"recent_activity"`
}

// AuditTimeRange represents a time range for audit logs
type AuditTimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// AuditLogFilter defines filters for querying audit logs
type AuditLogFilter struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Actions   []string   `json:"actions,omitempty"`
	Resources []string   `json:"resources,omitempty"`
	User      string     `json:"user,omitempty"`
	Success   *bool      `json:"success,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

// LogAction logs a configuration management action
func (a *ConfigAuditLogger) LogAction(action, resource string, details map[string]interface{}) {
	entry := &AuditLogEntry{
		Timestamp: time.Now(),
		Action:    action,
		Resource:  resource,
		Details:   details,
		Success:   true,
	}

	// Add user information if available
	if user := a.getCurrentUser(); user != "" {
		entry.User = user
	}

	if err := a.writeLogEntry(entry); err != nil && a.logger != nil {
		a.logger.ErrorWithFields("Failed to write audit log entry", map[string]interface{}{
			"action":   action,
			"resource": resource,
			"error":    err.Error(),
		})
	}
}

// LogError logs a configuration management error
func (a *ConfigAuditLogger) LogError(action, resource string, err error, details map[string]interface{}) {
	entry := &AuditLogEntry{
		Timestamp: time.Now(),
		Action:    action,
		Resource:  resource,
		Details:   details,
		Success:   false,
		Error:     err.Error(),
	}

	// Add user information if available
	if user := a.getCurrentUser(); user != "" {
		entry.User = user
	}

	if writeErr := a.writeLogEntry(entry); writeErr != nil && a.logger != nil {
		a.logger.ErrorWithFields("Failed to write audit log entry", map[string]interface{}{
			"action":   action,
			"resource": resource,
			"error":    writeErr.Error(),
		})
	}
}

// GetAuditLogs retrieves audit logs with optional filtering
func (a *ConfigAuditLogger) GetAuditLogs(filter *AuditLogFilter) ([]*AuditLogEntry, error) {
	if filter == nil {
		filter = &AuditLogFilter{
			Limit: 100,
		}
	}

	// Read all log entries
	entries, err := a.readLogEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to read log entries: %w", err)
	}

	// Apply filters
	filtered := a.applyFilters(entries, filter)

	// Apply pagination
	if filter.Offset > 0 || filter.Limit > 0 {
		filtered = a.applyPagination(filtered, filter.Offset, filter.Limit)
	}

	return filtered, nil
}

// GetAuditSummary generates a summary of audit log activity
func (a *ConfigAuditLogger) GetAuditSummary(filter *AuditLogFilter) (*AuditLogSummary, error) {
	entries, err := a.GetAuditLogs(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	summary := &AuditLogSummary{
		TotalEntries:   len(entries),
		ActionCounts:   make(map[string]int),
		ResourceCounts: make(map[string]int),
		TopErrors:      []string{},
		RecentActivity: []*AuditLogEntry{},
	}

	errorCounts := make(map[string]int)
	var timeRange *AuditTimeRange

	for _, entry := range entries {
		// Count successes and errors
		if entry.Success {
			summary.SuccessCount++
		} else {
			summary.ErrorCount++
			if entry.Error != "" {
				errorCounts[entry.Error]++
			}
		}

		// Count actions and resources
		summary.ActionCounts[entry.Action]++
		summary.ResourceCounts[entry.Resource]++

		// Track time range
		if timeRange == nil {
			timeRange = &AuditTimeRange{
				Start: entry.Timestamp,
				End:   entry.Timestamp,
			}
		} else {
			if entry.Timestamp.Before(timeRange.Start) {
				timeRange.Start = entry.Timestamp
			}
			if entry.Timestamp.After(timeRange.End) {
				timeRange.End = entry.Timestamp
			}
		}
	}

	summary.TimeRange = timeRange

	// Get top errors
	summary.TopErrors = a.getTopErrors(errorCounts, 5)

	// Get recent activity (last 10 entries)
	recentCount := 10
	if len(entries) < recentCount {
		recentCount = len(entries)
	}
	summary.RecentActivity = entries[:recentCount]

	return summary, nil
}

// ClearAuditLogs clears all audit log entries
func (a *ConfigAuditLogger) ClearAuditLogs() error {
	// Backup current log before clearing
	if err := a.backupLogFile(); err != nil {
		return fmt.Errorf("failed to backup log file: %w", err)
	}

	// Clear the log file
	if err := os.Truncate(a.logFile, 0); err != nil {
		return fmt.Errorf("failed to clear log file: %w", err)
	}

	// Log the clear action
	a.LogAction("clear_logs", "audit_log", map[string]interface{}{
		"cleared_at": time.Now(),
	})

	return nil
}

// RotateAuditLogs rotates the audit log file
func (a *ConfigAuditLogger) RotateAuditLogs() error {
	// Create rotated filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	rotatedFile := fmt.Sprintf("%s.%s", a.logFile, timestamp)

	// Move current log to rotated file
	if err := os.Rename(a.logFile, rotatedFile); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Log the rotation action
	a.LogAction("rotate_logs", "audit_log", map[string]interface{}{
		"rotated_to": rotatedFile,
		"rotated_at": time.Now(),
	})

	return nil
}

// ValidateAuditLog validates the integrity of the audit log
func (a *ConfigAuditLogger) ValidateAuditLog() (*AuditLogValidation, error) {
	validation := &AuditLogValidation{
		Valid:      true,
		Errors:     []string{},
		Warnings:   []string{},
		Statistics: make(map[string]interface{}),
	}

	// Check if log file exists
	if _, err := os.Stat(a.logFile); os.IsNotExist(err) {
		validation.Warnings = append(validation.Warnings, "Audit log file does not exist")
		return validation, nil
	}

	// Read and validate entries
	entries, err := a.readLogEntries()
	if err != nil {
		validation.Valid = false
		validation.Errors = append(validation.Errors, fmt.Sprintf("Failed to read log entries: %v", err))
		return validation, nil
	}

	// Validate each entry
	for i, entry := range entries {
		if err := a.validateLogEntry(entry); err != nil {
			validation.Errors = append(validation.Errors, fmt.Sprintf("Entry %d: %v", i, err))
			validation.Valid = false
		}
	}

	// Add statistics
	validation.Statistics["total_entries"] = len(entries)
	validation.Statistics["file_size"] = a.getLogFileSize()
	validation.Statistics["oldest_entry"] = a.getOldestEntryTime(entries)
	validation.Statistics["newest_entry"] = a.getNewestEntryTime(entries)

	return validation, nil
}

// AuditLogValidation represents the result of audit log validation
type AuditLogValidation struct {
	Valid      bool                   `json:"valid"`
	Errors     []string               `json:"errors"`
	Warnings   []string               `json:"warnings"`
	Statistics map[string]interface{} `json:"statistics"`
}

// writeLogEntry writes a single log entry to the audit log file
func (a *ConfigAuditLogger) writeLogEntry(entry *AuditLogEntry) error {
	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(a.logFile), 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file for appending
	file, err := os.OpenFile(a.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Marshal entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Write entry with newline
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

// readLogEntries reads all log entries from the audit log file
func (a *ConfigAuditLogger) readLogEntries() ([]*AuditLogEntry, error) {
	// Check if file exists
	if _, err := os.Stat(a.logFile); os.IsNotExist(err) {
		return []*AuditLogEntry{}, nil
	}

	// Read file content
	content, err := os.ReadFile(a.logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// Parse entries line by line
	var entries []*AuditLogEntry
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry AuditLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Log parsing error but continue
			if a.logger != nil {
				a.logger.WarnWithFields("Failed to parse audit log entry", map[string]interface{}{
					"line":  i + 1,
					"error": err.Error(),
				})
			}
			continue
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

// applyFilters applies filters to audit log entries
func (a *ConfigAuditLogger) applyFilters(entries []*AuditLogEntry, filter *AuditLogFilter) []*AuditLogEntry {
	var filtered []*AuditLogEntry

	for _, entry := range entries {
		// Time range filter
		if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
			continue
		}

		// Action filter
		if len(filter.Actions) > 0 {
			found := false
			for _, action := range filter.Actions {
				if entry.Action == action {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Resource filter
		if len(filter.Resources) > 0 {
			found := false
			for _, resource := range filter.Resources {
				if entry.Resource == resource {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// User filter
		if filter.User != "" && entry.User != filter.User {
			continue
		}

		// Success filter
		if filter.Success != nil && entry.Success != *filter.Success {
			continue
		}

		filtered = append(filtered, entry)
	}

	return filtered
}

// applyPagination applies pagination to audit log entries
func (a *ConfigAuditLogger) applyPagination(entries []*AuditLogEntry, offset, limit int) []*AuditLogEntry {
	if offset >= len(entries) {
		return []*AuditLogEntry{}
	}

	end := offset + limit
	if limit <= 0 || end > len(entries) {
		end = len(entries)
	}

	return entries[offset:end]
}

// getCurrentUser gets the current user for audit logging
func (a *ConfigAuditLogger) getCurrentUser() string {
	// Try to get user from environment
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

// backupLogFile creates a backup of the current log file
func (a *ConfigAuditLogger) backupLogFile() error {
	if _, err := os.Stat(a.logFile); os.IsNotExist(err) {
		return nil // No file to backup
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s.backup.%s", a.logFile, timestamp)

	content, err := os.ReadFile(a.logFile)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	if err := os.WriteFile(backupFile, content, 0640); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// getTopErrors returns the most frequent errors
func (a *ConfigAuditLogger) getTopErrors(errorCounts map[string]int, limit int) []string {
	type errorCount struct {
		error string
		count int
	}

	var errors []errorCount
	for err, count := range errorCounts {
		errors = append(errors, errorCount{error: err, count: count})
	}

	// Sort by count (descending)
	for i := 0; i < len(errors)-1; i++ {
		for j := i + 1; j < len(errors); j++ {
			if errors[j].count > errors[i].count {
				errors[i], errors[j] = errors[j], errors[i]
			}
		}
	}

	// Return top errors
	var topErrors []string
	for i := 0; i < limit && i < len(errors); i++ {
		topErrors = append(topErrors, errors[i].error)
	}

	return topErrors
}

// validateLogEntry validates a single log entry
func (a *ConfigAuditLogger) validateLogEntry(entry *AuditLogEntry) error {
	if entry.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if entry.Action == "" {
		return fmt.Errorf("action is required")
	}

	if entry.Resource == "" {
		return fmt.Errorf("resource is required")
	}

	return nil
}

// getLogFileSize returns the size of the log file
func (a *ConfigAuditLogger) getLogFileSize() int64 {
	if info, err := os.Stat(a.logFile); err == nil {
		return info.Size()
	}
	return 0
}

// getOldestEntryTime returns the timestamp of the oldest entry
func (a *ConfigAuditLogger) getOldestEntryTime(entries []*AuditLogEntry) *time.Time {
	if len(entries) == 0 {
		return nil
	}

	oldest := entries[0].Timestamp
	for _, entry := range entries {
		if entry.Timestamp.Before(oldest) {
			oldest = entry.Timestamp
		}
	}

	return &oldest
}

// getNewestEntryTime returns the timestamp of the newest entry
func (a *ConfigAuditLogger) getNewestEntryTime(entries []*AuditLogEntry) *time.Time {
	if len(entries) == 0 {
		return nil
	}

	newest := entries[0].Timestamp
	for _, entry := range entries {
		if entry.Timestamp.After(newest) {
			newest = entry.Timestamp
		}
	}

	return &newest
}
