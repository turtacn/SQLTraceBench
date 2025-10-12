package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Field is a key-value pair used for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is a wrapper around logrus.
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new Logger instance.
func NewLogger(level, format string, output io.Writer) *Logger {
	logger := logrus.New()
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	if output == nil {
		output = os.Stdout
	}
	logger.SetOutput(output)

	return &Logger{logger}
}

// Info logs a message at the info level with structured fields.
func (l *Logger) Info(msg string, fields ...Field) {
	l.WithFields(fieldsToLogrus(fields)).Info(msg)
}

// Error logs a message at the error level with structured fields.
func (l *Logger) Error(msg string, fields ...Field) {
	l.WithFields(fieldsToLogrus(fields)).Error(msg)
}

// fieldsToLogrus converts a slice of our custom Field type to logrus.Fields.
func fieldsToLogrus(fields []Field) logrus.Fields {
	f := make(logrus.Fields)
	for _, field := range fields {
		f[field.Key] = field.Value
	}
	return f
}

var globalLogger *Logger

// SetGlobalLogger sets the global logger instance.
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance.
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// Provide a default logger if none is set.
		globalLogger = NewLogger("info", "text", nil)
	}
	return globalLogger
}