package models

// Schema represents the entire database schema, containing multiple databases.
type Schema struct {
	Databases []DatabaseSchema `json:"databases"`
}

// DatabaseSchema represents the schema of an entire database.
type DatabaseSchema struct {
	Name   string         `json:"name"`
	Tables []*TableSchema `json:"tables"`
}

// TableSchema represents the schema of a single database table.
type TableSchema struct {
	Name    string                  `json:"name"`
	Columns []*ColumnSchema         `json:"columns"`
	PK      []string                `json:"pk"` // Primary Key columns
	Indexes map[string]*IndexSchema `json:"indexes"`
	Engine  string                  `json:"engine,omitempty"` // e.g., "MergeTree() ORDER BY ..."
}

// ColumnSchema represents a single column in a database table.
type ColumnSchema struct {
	Name         string `json:"name"`
	DataType     string `json:"data_type"`
	IsNullable   bool   `json:"is_nullable"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	Default      string `json:"default"`
}

// IndexSchema represents a database index.
type IndexSchema struct {
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
}
