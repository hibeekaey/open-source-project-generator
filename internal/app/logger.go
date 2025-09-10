package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger provides structured logging for the application
type Logger struct {
	level       LogLevel
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errLogger   *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, logToFile bool) (*Logger, error) {
	logger := &Logger{
		level: level,
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Add file logging if requested
	if logToFile {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		logDir := filepath.Join(homeDir, ".cache", "template-generator", "logs")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		logFile := filepath.Join(logDir, fmt.Sprintf("generator-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.logFile = file
		writers = append(writers, file)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create loggers for different levels
	logger.infoLogger = log.New(multiWriter, "INFO  ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.warnLogger = log.New(multiWriter, "WARN  ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.errLogger = log.New(multiWriter, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.debugLogger = log.New(multiWriter, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)

	return logger, nil
}

// Close closes the logger and any open file handles
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.debugLogger.Printf(msg, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.infoLogger.Printf(msg, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.warnLogger.Printf(msg, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.errLogger.Printf(msg, args...)
	}
}

// LogLevelFromString converts a string to LogLevel
func LogLevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}
