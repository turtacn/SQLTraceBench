package types

import (
	"fmt"
	"time"
)

type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrInvalidInput
	ErrParseFailed
	ErrConversionFailed
	ErrDatabaseConnection
	ErrPluginNotFound
	ErrValidationFailed
	ErrExecutionFailed
	ErrReportGeneration
)

func (e ErrorCode) String() string {
	return [...]string{
		"ErrUnknown",
		"ErrInvalidInput",
		"ErrParseFailed",
		"ErrConversionFailed",
		"ErrDatabaseConnection",
		"ErrPluginNotFound",
		"ErrValidationFailed",
		"ErrExecutionFailed",
		"ErrReportGeneration",
	}[e]
}

type SQLTraceBenchError struct {
	Code      ErrorCode
	Message   string
	Details   string
	Component string
	Timestamp time.Time
	Cause     error
}

func (e *SQLTraceBenchError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (cause: %v)", e.Component, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Component, e.Code, e.Message)
}

func NewError(code ErrorCode, message string) *SQLTraceBenchError {
	return &SQLTraceBenchError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func WrapError(code ErrorCode, message string, cause error) *SQLTraceBenchError {
	return &SQLTraceBenchError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
	}
}

func IsErrorCode(err error, code ErrorCode) bool {
	if e, ok := err.(*SQLTraceBenchError); ok {
		return e.Code == code
	}
	return false
}

//Personal.AI order the ending
