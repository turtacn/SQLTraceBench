package services

import (
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaService is responsible for schema-related operations, such as conversion.
type SchemaService struct {
	typeMappers map[string]TypeMapper
}

// TypeMapper defines the interface for mapping data types between different database dialects.
type TypeMapper func(sourceType string) (string, error)

// NewSchemaService creates a new SchemaService.
func NewSchemaService() *SchemaService {
	return &SchemaService{
		typeMappers: make(map[string]TypeMapper),
	}
}

// RegisterTypeMapper registers a type mapper for a specific target dialect.
func (s *SchemaService) RegisterTypeMapper(dialect string, mapper TypeMapper) {
	s.typeMappers[dialect] = mapper
}

// ConvertTo converts a DatabaseSchema to a target dialect.
// It iterates through the tables and columns, applying the appropriate type mapping.
func (s *SchemaService) ConvertTo(sourceSchema *models.DatabaseSchema, targetDialect string) (*models.DatabaseSchema, error) {
	mapper, ok := s.typeMappers[targetDialect]
	if !ok {
		return nil, fmt.Errorf("no type mapper registered for dialect: %s", targetDialect)
	}

	convertedSchema := &models.DatabaseSchema{
		Name:   sourceSchema.Name,
		Tables: make([]*models.TableSchema, 0, len(sourceSchema.Tables)),
	}

	for _, table := range sourceSchema.Tables {
		convertedTable := &models.TableSchema{
			Name:    table.Name,
			Columns: make([]*models.ColumnSchema, len(table.Columns)),
			PK:      table.PK,
			Indexes: table.Indexes,
			Engine:  table.Engine,
		}

		for i, col := range table.Columns {
			convertedType, err := mapper(col.DataType) // Use DataType
			if err != nil {
				return nil, fmt.Errorf("failed to convert type for column %s in table %s: %w", col.Name, table.Name, err)
			}
			convertedTable.Columns[i] = &models.ColumnSchema{
				Name:         col.Name,
				DataType:     convertedType, // Use DataType
				IsNullable:   col.IsNullable,
				IsPrimaryKey: col.IsPrimaryKey,
				Default:      col.Default,
			}
		}
		convertedSchema.Tables = append(convertedSchema.Tables, convertedTable)
	}

	return convertedSchema, nil
}

// MySQLToClickHouseTypeMapper is a simple type mapper for converting MySQL types to ClickHouse types.
func MySQLToClickHouseTypeMapper(sourceType string) (string, error) {
	sourceType = strings.ToLower(sourceType)
	switch {
	case strings.HasPrefix(sourceType, "int"):
		return "Int32", nil
	case strings.HasPrefix(sourceType, "bigint"):
		return "Int64", nil
	case strings.HasPrefix(sourceType, "varchar"):
		return "String", nil
	case strings.HasPrefix(sourceType, "text"):
		return "String", nil
	case strings.HasPrefix(sourceType, "datetime"):
		return "DateTime", nil
	case strings.HasPrefix(sourceType, "timestamp"):
		return "DateTime", nil
	default:
		return "String", nil // Default to String for unknown types.
	}
}