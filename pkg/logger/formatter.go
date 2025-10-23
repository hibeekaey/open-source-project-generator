// Package logger provides structured logging with context support.
package logger

import (
	"fmt"
	"strings"
)

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Status symbols
const (
	SymbolSuccess = "✓"
	SymbolError   = "✗"
	SymbolWarning = "⚠"
	SymbolInfo    = "ℹ"
	SymbolArrow   = "→"
	SymbolBullet  = "•"
)

// Formatter provides user-friendly output formatting
type Formatter struct {
	enableColor bool
}

// NewFormatter creates a new formatter
func NewFormatter(enableColor bool) *Formatter {
	return &Formatter{
		enableColor: enableColor,
	}
}

// Success formats a success message
func (f *Formatter) Success(message string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s %s", ColorGreen, SymbolSuccess, ColorReset, message)
	}
	return fmt.Sprintf("%s %s", SymbolSuccess, message)
}

// Error formats an error message
func (f *Formatter) Error(message string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s %s", ColorRed, SymbolError, ColorReset, message)
	}
	return fmt.Sprintf("%s %s", SymbolError, message)
}

// Warning formats a warning message
func (f *Formatter) Warning(message string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s %s", ColorYellow, SymbolWarning, ColorReset, message)
	}
	return fmt.Sprintf("%s %s", SymbolWarning, message)
}

// Info formats an info message
func (f *Formatter) Info(message string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s %s", ColorBlue, SymbolInfo, ColorReset, message)
	}
	return fmt.Sprintf("%s %s", SymbolInfo, message)
}

// Step formats a step message
func (f *Formatter) Step(step int, total int, message string) string {
	if f.enableColor {
		return fmt.Sprintf("%s[%d/%d]%s %s", ColorBold, step, total, ColorReset, message)
	}
	return fmt.Sprintf("[%d/%d] %s", step, total, message)
}

// Header formats a header message
func (f *Formatter) Header(message string) string {
	separator := strings.Repeat("=", 70)
	if f.enableColor {
		return fmt.Sprintf("%s%s%s\n%s%s%s\n%s%s%s",
			ColorBold, separator, ColorReset,
			ColorBold, message, ColorReset,
			ColorBold, separator, ColorReset)
	}
	return fmt.Sprintf("%s\n%s\n%s", separator, message, separator)
}

// Section formats a section header
func (f *Formatter) Section(message string) string {
	separator := strings.Repeat("-", 50)
	if f.enableColor {
		return fmt.Sprintf("\n%s%s%s\n%s%s%s",
			ColorBold, message, ColorReset,
			ColorCyan, separator, ColorReset)
	}
	return fmt.Sprintf("\n%s\n%s", message, separator)
}

// Bullet formats a bullet point
func (f *Formatter) Bullet(message string) string {
	if f.enableColor {
		return fmt.Sprintf("  %s%s%s %s", ColorCyan, SymbolBullet, ColorReset, message)
	}
	return fmt.Sprintf("  %s %s", SymbolBullet, message)
}

// Arrow formats an arrow message
func (f *Formatter) Arrow(message string) string {
	if f.enableColor {
		return fmt.Sprintf("  %s%s%s %s", ColorBlue, SymbolArrow, ColorReset, message)
	}
	return fmt.Sprintf("  %s %s", SymbolArrow, message)
}

// Bold formats text in bold
func (f *Formatter) Bold(text string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s", ColorBold, text, ColorReset)
	}
	return text
}

// Colorize applies a color to text
func (f *Formatter) Colorize(text string, color string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s%s", color, text, ColorReset)
	}
	return text
}

// Duration formats a duration message
func (f *Formatter) Duration(message string, duration string) string {
	if f.enableColor {
		return fmt.Sprintf("%s %s(%s)%s", message, ColorCyan, duration, ColorReset)
	}
	return fmt.Sprintf("%s (%s)", message, duration)
}

// KeyValue formats a key-value pair
func (f *Formatter) KeyValue(key string, value string) string {
	if f.enableColor {
		return fmt.Sprintf("%s%s:%s %s", ColorBold, key, ColorReset, value)
	}
	return fmt.Sprintf("%s: %s", key, value)
}

// List formats a list of items
func (f *Formatter) List(items []string) string {
	var result strings.Builder
	for _, item := range items {
		result.WriteString(f.Bullet(item))
		result.WriteString("\n")
	}
	return result.String()
}

// Table formats a simple table
func (f *Formatter) Table(headers []string, rows [][]string) string {
	var result strings.Builder

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

	// Format header
	for i, header := range headers {
		if f.enableColor {
			result.WriteString(fmt.Sprintf("%s%-*s%s  ", ColorBold, colWidths[i], header, ColorReset))
		} else {
			result.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], header))
		}
	}
	result.WriteString("\n")

	// Format separator
	for _, width := range colWidths {
		result.WriteString(strings.Repeat("-", width))
		result.WriteString("  ")
	}
	result.WriteString("\n")

	// Format rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				result.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], cell))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// ProgressBar creates a simple progress bar
func (f *Formatter) ProgressBar(current int, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentStr := fmt.Sprintf("%.0f%%", percentage*100)

	if f.enableColor {
		return fmt.Sprintf("[%s%s%s] %s", ColorGreen, bar, ColorReset, percentStr)
	}
	return fmt.Sprintf("[%s] %s", bar, percentStr)
}

// Box creates a box around text
func (f *Formatter) Box(text string) string {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var result strings.Builder
	border := "+" + strings.Repeat("-", maxLen+2) + "+"

	result.WriteString(border + "\n")
	for _, line := range lines {
		padding := maxLen - len(line)
		result.WriteString(fmt.Sprintf("| %s%s |\n", line, strings.Repeat(" ", padding)))
	}
	result.WriteString(border)

	return result.String()
}
