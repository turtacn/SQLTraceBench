package models

// DatabaseSchema represents the schema of an entire database.
type DatabaseSchema struct {
	Name   string                 `json:"name"`
	Tables map[string]*TableSchema `json:"tables"`
}

// TableSchema represents the schema of a single database table.
type TableSchema struct {
	Name    string                  `json:"name"`
	Columns []*ColumnSchema         `json:"columns"`
	PK      []string                `json:"pk"` // Primary Key columns
	Indexes map[string]*IndexSchema `json:"indexes"`
}

// ColumnSchema represents a single column in a database table.
type ColumnSchema struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"is_nullable"`
	Default    string `json:"default"`
}

// IndexSchema represents a database index.
type IndexSchema struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
}