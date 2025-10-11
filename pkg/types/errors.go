// Package types contains shared data structures and constants used across the application.
package types

import "fmt"

// ErrorCode defines a typed string for application-specific error codes.
type ErrorCode string

// Defines the set of standard error codes.
const (
	ErrInvalidInput ErrorCode = "invalid_input"
	ErrParseFailed  ErrorCode = "parse_failed"
	ErrNotFound     ErrorCode = "not_found"
	ErrInternal     ErrorCode = "internal"
)

// SQLTraceBenchError is a custom error type for the application.
// It includes a machine-readable error code, a human-readable message,
// and the underlying error that caused it.
type SQLTraceBenchError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error returns the string representation of the SQLTraceBenchError.
func (e *SQLTraceBenchError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a new SQLTraceBenchError.
func NewError(code ErrorCode, message string) *SQLTraceBenchError {
	return &SQLTraceBenchError{Code: code, Message: message}
}

// WrapError creates a new SQLTraceBenchError that wraps an existing error.
func WrapError(code ErrorCode, message string, cause error) *SQLTraceBenchError {
	return &SQLTraceBenchError{Code: code, Message: message, Cause: cause}
}