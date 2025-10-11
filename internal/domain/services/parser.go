// Package services contains the interfaces for the application's core services.
package services

// Parser is the interface for a SQL parser.
// It defines the contract for any parser implementation, whether it's regex-based or a full-fledged AST parser.
type Parser interface {
	// ListTables extracts the table names from a given SQL query.
	ListTables(sql string) ([]string, error)
}