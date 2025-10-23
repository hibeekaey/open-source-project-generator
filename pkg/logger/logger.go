// Package logger provides structured logging with context support.
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level represents the log level.
type Level int

const (
	// DebugLevel is for debug messages.
	DebugLevel Level = iota
	// InfoLevel is for informational messages.
	InfoLevel
	// WarnLevel is for warning messages.
	WarnLevel
	// ErrorLevel is for error messages.
	ErrorLevel
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging functionality.
type Logger struct {
	mu          sync.Mutex
	level       Level
	output      io.Writer
	fileOutput  io.Writer
	context     map[string]interface{}
	enableColor bool
	enableFile  bool
	formatter   *Formatter
}

// NewLogger creates a new logger with default settings.
func NewLogger() *Logger {
	return &Logger{
		level:       InfoLevel,
		output:      os.Stdout,
		context:     make(map[string]interface{}),
		enableColor: true,
		enableFile:  false,
		formatter:   NewFormatter(true),
	}
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output writer for console logging.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// EnableFileLogging enables logging to a file.
func (l *Logger) EnableFileLogging(logFilePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.fileOutput = file
	l.enableFile = true

	return nil
}

// DisableFileLogging disables logging to a file.
func (l *Logger) DisableFileLogging() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fileOutput != nil {
		if closer, ok := l.fileOutput.(io.Closer); ok {
			closer.Close()
		}
	}

	l.fileOutput = nil
	l.enableFile = false
}

// EnableColor enables colored output for console logging.
func (l *Logger) EnableColor(enable bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enableColor = enable
	l.formatter = NewFormatter(enable)
}

// WithContext returns a new logger with additional context.
func (l *Logger) WithContext(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newContext := make(map[string]interface{})
	for k, v := range l.context {
		newContext[k] = v
	}
	newContext[key] = value

	return &Logger{
		level:       l.level,
		output:      l.output,
		fileOutput:  l.fileOutput,
		context:     newContext,
		enableColor: l.enableColor,
		enableFile:  l.enableFile,
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message, nil)
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...), nil)
}

// Info logs an informational message.
func (l *Logger) Info(message string) {
	l.log(InfoLevel, message, nil)
}

// Infof logs a formatted informational message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
}

// Warn logs a warning message.
func (l *Logger) Warn(message string) {
	l.log(WarnLevel, message, nil)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...), nil)
}

// Error logs an error message.
func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message, nil)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
}

// ErrorWithErr logs an error message with an error object.
func (l *Logger) ErrorWithErr(message string, err error) {
	l.log(ErrorLevel, message, map[string]interface{}{"error": err.Error()})
}

// log is the internal logging function.
func (l *Logger) log(level Level, message string, extra map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if we should log this level
	if level < l.level {
		return
	}

	// Build log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := level.String()

	// Format console output
	consoleMsg := l.formatConsoleMessage(timestamp, levelStr, message, extra)

	// Format file output
	fileMsg := l.formatFileMessage(timestamp, levelStr, message, extra)

	// Write to console
	if l.output != nil {
		fmt.Fprintln(l.output, consoleMsg)
	}

	// Write to file
	if l.enableFile && l.fileOutput != nil {
		fmt.Fprintln(l.fileOutput, fileMsg)
	}
}

// formatConsoleMessage formats a log message for console output.
func (l *Logger) formatConsoleMessage(timestamp, level, message string, extra map[string]interface{}) string {
	var colorCode string
	var resetCode string

	if l.enableColor {
		resetCode = "\033[0m"
		switch level {
		case "DEBUG":
			colorCode = "\033[36m" // Cyan
		case "INFO":
			colorCode = "\033[32m" // Green
		case "WARN":
			colorCode = "\033[33m" // Yellow
		case "ERROR":
			colorCode = "\033[31m" // Red
		}
	}

	// Build message
	msg := fmt.Sprintf("%s[%s]%s %s %s", colorCode, level, resetCode, timestamp, message)

	// Add context
	if len(l.context) > 0 {
		msg += " |"
		for k, v := range l.context {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	// Add extra fields
	if len(extra) > 0 {
		msg += " |"
		for k, v := range extra {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	return msg
}

// formatFileMessage formats a log message for file output.
func (l *Logger) formatFileMessage(timestamp, level, message string, extra map[string]interface{}) string {
	// Build message without colors
	msg := fmt.Sprintf("[%s] %s %s", level, timestamp, message)

	// Add context
	if len(l.context) > 0 {
		msg += " |"
		for k, v := range l.context {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	// Add extra fields
	if len(extra) > 0 {
		msg += " |"
		for k, v := range extra {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	return msg
}

// GetWriter returns the output writer for the logger
func (l *Logger) GetWriter() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.output
}

// Close closes the logger and any open file handles.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fileOutput != nil {
		if closer, ok := l.fileOutput.(io.Closer); ok {
			return closer.Close()
		}
	}

	return nil
}

// GetFormatter returns the formatter for this logger
func (l *Logger) GetFormatter() *Formatter {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.formatter
}

// Success logs a success message with formatting
func (l *Logger) Success(message string) {
	l.mu.Lock()
	formatted := l.formatter.Success(message)
	l.mu.Unlock()
	l.Info(formatted)
}

// PrintHeader prints a formatted header
func (l *Logger) PrintHeader(message string) {
	l.mu.Lock()
	formatted := l.formatter.Header(message)
	l.mu.Unlock()
	if l.output != nil {
		fmt.Fprintln(l.output, formatted)
	}
}

// PrintSection prints a formatted section header
func (l *Logger) PrintSection(message string) {
	l.mu.Lock()
	formatted := l.formatter.Section(message)
	l.mu.Unlock()
	if l.output != nil {
		fmt.Fprintln(l.output, formatted)
	}
}

// PrintBullet prints a formatted bullet point
func (l *Logger) PrintBullet(message string) {
	l.mu.Lock()
	formatted := l.formatter.Bullet(message)
	l.mu.Unlock()
	if l.output != nil {
		fmt.Fprintln(l.output, formatted)
	}
}

// PrintKeyValue prints a formatted key-value pair
func (l *Logger) PrintKeyValue(key string, value string) {
	l.mu.Lock()
	formatted := l.formatter.KeyValue(key, value)
	l.mu.Unlock()
	if l.output != nil {
		fmt.Fprintln(l.output, formatted)
	}
}

// Global logger instance
var defaultLogger = NewLogger()

// SetDefaultLevel sets the level for the default logger.
func SetDefaultLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetDefaultOutput sets the output for the default logger.
func SetDefaultOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// EnableDefaultFileLogging enables file logging for the default logger.
func EnableDefaultFileLogging(logFilePath string) error {
	return defaultLogger.EnableFileLogging(logFilePath)
}

// DisableDefaultFileLogging disables file logging for the default logger.
func DisableDefaultFileLogging() {
	defaultLogger.DisableFileLogging()
}

// EnableDefaultColor enables colored output for the default logger.
func EnableDefaultColor(enable bool) {
	defaultLogger.EnableColor(enable)
}

// Debug logs a debug message using the default logger.
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Debugf logs a formatted debug message using the default logger.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info logs an informational message using the default logger.
func Info(message string) {
	defaultLogger.Info(message)
}

// Infof logs a formatted informational message using the default logger.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warn logs a warning message using the default logger.
func Warn(message string) {
	defaultLogger.Warn(message)
}

// Warnf logs a formatted warning message using the default logger.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Error logs an error message using the default logger.
func Error(message string) {
	defaultLogger.Error(message)
}

// Errorf logs a formatted error message using the default logger.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// ErrorWithErr logs an error message with an error object using the default logger.
func ErrorWithErr(message string, err error) {
	defaultLogger.ErrorWithErr(message, err)
}

// WithContext returns a new logger with additional context using the default logger.
func WithContext(key string, value interface{}) *Logger {
	return defaultLogger.WithContext(key, value)
}
