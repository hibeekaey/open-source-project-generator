// Package utils provides utility functions for the CLI interface.
//
// This module contains formatting utilities for CLI output including
// data formatting, machine-readable output, and display helpers.
package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Formatter provides output formatting utilities for the CLI.
//
// The Formatter provides methods for:
//   - Machine-readable output formatting (JSON, YAML)
//   - Human-readable display formatting
//   - Data structure formatting
//   - Error message formatting
type Formatter struct{}

// NewFormatter creates a new formatter instance.
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatMachineReadable formats data in machine-readable format
func (f *Formatter) FormatMachineReadable(data interface{}, format string) (string, error) {
	switch strings.ToLower(format) {
	case "json":
		return f.formatJSON(data)
	case "yaml", "yml":
		return f.formatYAML(data)
	default:
		return "", fmt.Errorf("unsupported output format: %s (supported: json, yaml)", format)
	}
}

// formatJSON formats data as JSON
func (f *Formatter) formatJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonData), nil
}

// formatYAML formats data as YAML
func (f *Formatter) formatYAML(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(yamlData), nil
}

// FormatBytes formats byte count as human-readable string
func (f *Formatter) FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats duration as human-readable string
func (f *Formatter) FormatDuration(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%.0fms", float64(duration.Nanoseconds())/1e6)
	}
	if duration < time.Minute {
		return fmt.Sprintf("%.1fs", duration.Seconds())
	}
	if duration < time.Hour {
		return fmt.Sprintf("%.1fm", duration.Minutes())
	}
	return fmt.Sprintf("%.1fh", duration.Hours())
}

// FormatList formats a list of items with bullets
func (f *Formatter) FormatList(items []string, bullet string) string {
	if len(items) == 0 {
		return ""
	}

	if bullet == "" {
		bullet = "•"
	}

	var result strings.Builder
	for i, item := range items {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("%s %s", bullet, item))
	}
	return result.String()
}

// FormatTable formats data as a simple table
func (f *Formatter) FormatTable(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var result strings.Builder

	// Write headers
	for i, header := range headers {
		if i > 0 {
			result.WriteString(" | ")
		}
		result.WriteString(f.padRight(header, colWidths[i]))
	}
	result.WriteString("\n")

	// Write separator
	for i := range headers {
		if i > 0 {
			result.WriteString("-+-")
		}
		result.WriteString(strings.Repeat("-", colWidths[i]))
	}
	result.WriteString("\n")

	// Write rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				result.WriteString(" | ")
			}
			if i < len(colWidths) {
				result.WriteString(f.padRight(cell, colWidths[i]))
			} else {
				result.WriteString(cell)
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// padRight pads a string to the right with spaces
func (f *Formatter) padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// FormatKeyValue formats key-value pairs
func (f *Formatter) FormatKeyValue(pairs map[string]string, separator string) string {
	if len(pairs) == 0 {
		return ""
	}

	if separator == "" {
		separator = ": "
	}

	var result strings.Builder
	first := true
	for key, value := range pairs {
		if !first {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("%s%s%s", key, separator, value))
		first = false
	}
	return result.String()
}

// FormatProgress formats a progress indicator
func (f *Formatter) FormatProgress(current, total int, width int) string {
	if total <= 0 {
		return ""
	}

	if width <= 0 {
		width = 20
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))

	var result strings.Builder
	result.WriteString("[")

	for i := 0; i < width; i++ {
		if i < filled {
			result.WriteString("=")
		} else if i == filled && current < total {
			result.WriteString(">")
		} else {
			result.WriteString(" ")
		}
	}

	result.WriteString(fmt.Sprintf("] %d/%d (%.1f%%)", current, total, percentage*100))
	return result.String()
}

// FormatError formats an error message with context
func (f *Formatter) FormatError(err error, context string) string {
	if err == nil {
		return ""
	}

	if context == "" {
		return err.Error()
	}

	return fmt.Sprintf("%s: %s", context, err.Error())
}

// FormatSuccess formats a success message
func (f *Formatter) FormatSuccess(message string, details map[string]interface{}) string {
	var result strings.Builder
	result.WriteString(message)

	if len(details) > 0 {
		result.WriteString("\n")
		for key, value := range details {
			result.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// TruncateString truncates a string to a maximum length with ellipsis
func (f *Formatter) TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return s[:maxLength]
	}

	return s[:maxLength-3] + "..."
}

// WrapText wraps text to a specified width
func (f *Formatter) WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// IndentText indents each line of text with the specified prefix
func (f *Formatter) IndentText(text string, indent string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
