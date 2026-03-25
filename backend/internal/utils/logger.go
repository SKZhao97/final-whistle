package utils

import (
	"log"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// DebugLevel logs everything
	DebugLevel LogLevel = iota
	// InfoLevel logs info, warnings, and errors
	InfoLevel
	// WarnLevel logs warnings and errors
	WarnLevel
	// ErrorLevel logs only errors
	ErrorLevel
)

var (
	currentLevel = InfoLevel
	logger       = log.New(os.Stdout, "", log.LstdFlags)
)

// SetLogLevel sets the current logging level
func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DebugLevel
	case "info":
		currentLevel = InfoLevel
	case "warn", "warning":
		currentLevel = WarnLevel
	case "error":
		currentLevel = ErrorLevel
	default:
		currentLevel = InfoLevel
		logger.Printf("Unknown log level '%s', defaulting to INFO", level)
	}
}

// Debug logs debug messages
func Debug(format string, v ...interface{}) {
	if currentLevel <= DebugLevel {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs info messages
func Info(format string, v ...interface{}) {
	if currentLevel <= InfoLevel {
		logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs warning messages
func Warn(format string, v ...interface{}) {
	if currentLevel <= WarnLevel {
		logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	if currentLevel <= ErrorLevel {
		logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs fatal messages and exits
func Fatal(format string, v ...interface{}) {
	logger.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

// WithFields creates a structured log entry (simplified version)
func WithFields(fields map[string]interface{}) *StructuredLogger {
	return &StructuredLogger{fields: fields}
}

// StructuredLogger provides basic structured logging
type StructuredLogger struct {
	fields map[string]interface{}
}

// Debug logs with fields
func (sl *StructuredLogger) Debug(msg string) {
	if currentLevel <= DebugLevel {
		logger.Printf("[DEBUG] %s %v", msg, sl.fields)
	}
}

// Info logs with fields
func (sl *StructuredLogger) Info(msg string) {
	if currentLevel <= InfoLevel {
		logger.Printf("[INFO] %s %v", msg, sl.fields)
	}
}

// Warn logs with fields
func (sl *StructuredLogger) Warn(msg string) {
	if currentLevel <= WarnLevel {
		logger.Printf("[WARN] %s %v", msg, sl.fields)
	}
}

// Error logs with fields
func (sl *StructuredLogger) Error(msg string) {
	if currentLevel <= ErrorLevel {
		logger.Printf("[ERROR] %s %v", msg, sl.fields)
	}
}

// GetLogger returns the underlying logger
func GetLogger() *log.Logger {
	return logger
}

// SetLogger sets a custom logger
func SetLogger(l *log.Logger) {
	logger = l
}

// LogLevelString returns string representation of current log level
func LogLevelString() string {
	switch currentLevel {
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
