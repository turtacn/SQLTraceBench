// Package utils provides utility functions, such as logging.
package utils

import (
	"fmt"
	"log"
	"os"
)

// Field is a key-value pair used for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is a simplified logger that wraps the standard log.Logger.
type Logger struct {
	l *log.Logger
}

// NewLogger creates a new Logger instance.
// For this minimal implementation, level and format are ignored.
func NewLogger(level, format string, output *os.File) *Logger {
	return &Logger{l: log.New(os.Stdout, "", log.LstdFlags)}
}

// Info logs an informational message.
func (lg *Logger) Info(msg string, fields ...Field) {
	lg.l.Println("[INFO]", msg, fmtFields(fields...))
}

// Error logs an error message.
func (lg *Logger) Error(msg string, fields ...Field) {
	lg.l.Println("[ERROR]", msg, fmtFields(fields...))
}

// fmtFields formats a slice of Fields into a string for logging.
func fmtFields(fields ...Field) string {
	var s string
	for _, f := range fields {
		s += fmt.Sprintf(" %s=%v", f.Key, f.Value)
	}
	return s
}

var globalLogger *Logger

// SetGlobalLogger sets the global logger instance.
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance.
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// Provide a default logger if none is set
		globalLogger = NewLogger("info", "text", nil)
	}
	return globalLogger
}