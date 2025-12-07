package schema

import (
	"fmt"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaConverter defines the interface for converting schema DDLs between databases.
type SchemaConverter interface {
	// ConvertDDL converts DDL statements to the target database format.
	ConvertDDL(sourceDDL string, targetDB string) (string, error)

	// ConvertTable converts a single table structure.
	ConvertTable(sourceTable *models.TableSchema, targetDB string) (*models.TableSchema, error)

	// GetTypeMapping gets the type mapping (source type -> target type).
	GetTypeMapping(sourceType string, targetDB string) (string, error)
}

// ConverterFactory creates and manages SchemaConverters.
type ConverterFactory struct {
	converters map[string]SchemaConverter
}

// NewConverterFactory creates a new ConverterFactory with registered converters.
func NewConverterFactory() *ConverterFactory {
	f := &ConverterFactory{
		converters: make(map[string]SchemaConverter),
	}
	// Converters will be registered here or externally.
	// For circular dependency reasons, we might need to register them after creation
	// or in the init() of the specific converter files.
	// However, simple way is to instantiate them here if they are in the same package.
	// Since we are in package schema, and the converters will also be in package schema,
	// we can instantiate them here.

	f.converters["mysql"] = NewMySQLConverter()
	f.converters["postgres"] = NewPostgresConverter()
	f.converters["postgresql"] = f.converters["postgres"] // Alias
	f.converters["tidb"] = NewTiDBConverter()

	return f
}

// GetConverter returns the converter for the specified source database.
func (f *ConverterFactory) GetConverter(sourceDB string) (SchemaConverter, error) {
	if converter, ok := f.converters[sourceDB]; ok {
		return converter, nil
	}
	return nil, fmt.Errorf("unsupported source database: %s", sourceDB)
}
