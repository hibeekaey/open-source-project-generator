// Package errors provides logger adapters for interface compatibility
package errors

import "fmt"

// LoggerAdapter adapts ErrorLogger to RecoveryLogger interface
type LoggerAdapter struct {
	logger *ErrorLogger
}

// NewLoggerAdapter creates a new logger adapter
func NewLoggerAdapter(logger *ErrorLogger) *LoggerAdapter {
	return &LoggerAdapter{logger: logger}
}

// Info logs an info message
func (la *LoggerAdapter) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	la.logger.Info(message, nil)
}

// Warn logs a warning message
func (la *LoggerAdapter) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	la.logger.Warn(message, nil)
}

// Error logs an error message
func (la *LoggerAdapter) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	la.logger.Error(message, nil)
}

// Debug logs a debug message
func (la *LoggerAdapter) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	la.logger.Debug(message, nil)
}
