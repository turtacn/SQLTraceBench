package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

type Field struct {
	Key   string
	Value interface{}
}

type Logger struct {
	entry *logrus.Entry
}

var globalLogger *Logger

func NewLogger(level string, format string, output io.Writer) *Logger {
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
	return &Logger{entry: logrus.NewEntry(logger)}
}

func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewLogger("info", "text", nil)
	}
	return globalLogger
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToMap(fields)).Debug(msg)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToMap(fields)).Info(msg)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToMap(fields)).Warn(msg)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToMap(fields)).Error(msg)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToMap(fields)).Fatal(msg)
}

func (l *Logger) WithFields(fields ...Field) *Logger {
	return &Logger{entry: l.entry.WithFields(fieldsToMap(fields))}
}

func (l *Logger) WithError(err *types.SQLTraceBenchError) *Logger {
	return &Logger{entry: l.entry.WithField("error", err)}
}

func fieldsToMap(fields []Field) map[string]interface{} {
	m := make(map[string]interface{})
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	return m
}

//Personal.AI order the ending
