package models

import (
	"github.com/turtacn/SQLTraceBench/pkg/types"
	"regexp"
	"strings"
)

type DatabaseSchema struct {
	DatabaseName     string
	TableDefinitions []TableDefinition
}

type TableDefinition struct {
	Name        string
	Columns     []ColumnDefinition
	PrimaryKeys []string
	Indexes     []IndexDefinition
	Partition   string
	Engine      string
}

type ColumnDefinition struct {
	Name     string
	Type     string
	Length   int
	Nullable bool
	Default  interface{}
}

type IndexDefinition struct {
	Name    string
	Type    string
	Columns []string
	Unique  bool
}

func (s *DatabaseSchema) ConvertTo(target types.DatabaseType) (*DatabaseSchema, error) {
	// Minimal conversion for MVP: just clone for now
	return &DatabaseSchema{DatabaseName: s.DatabaseName, TableDefinitions: s.TableDefinitions}, nil
}

func (s *DatabaseSchema) Validate() *types.SQLTraceBenchError {
	if len(s.TableDefinitions) == 0 {
		return types.NewError(types.ErrInvalidInput, "schema must contain at least one table")
	}
	return nil
}

func (s *DatabaseSchema) ExtractFromSQL(sqlText string) (*DatabaseSchema, error) {
	// Very basic extraction for MVP - just match CREATE TABLE statements
	schema := &DatabaseSchema{}
	re := regexp.MustCompile(`(?i)CREATE TABLE\s+([^\s(]+)\s*\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(sqlText, -1)

	for _, match := range matches {
		table := TableDefinition{Name: strings.Trim(match[1], "`'\"")}
		columnsText := match[2]
		cols := parseColumns(columnsText)
		table.Columns = cols
		schema.TableDefinitions = append(schema.TableDefinitions, table)
	}
	return schema, nil
}

func parseColumns(text string) []ColumnDefinition {
	// Simplified parsing
	parts := strings.Split(text, ",")
	var cols []ColumnDefinition
	for _, part := range parts {
		fields := strings.Fields(strings.TrimSpace(part))
		if len(fields) >= 2 {
			cols = append(cols, ColumnDefinition{
				Name: strings.Trim(fields[0], "`'\""),
				Type: strings.ToUpper(fields[1]),
			})
		}
	}
	return cols
}

//Personal.AI order the ending
