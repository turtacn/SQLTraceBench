package clickhouse

import (
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaConverter is the interface for schema conversion logic.
type SchemaConverter interface {
	ConvertSchema(source *models.Schema, targetType string) (*models.Schema, error)
}

type ClickHouseConverter struct{}

// NewSchemaConverter creates a new SchemaConverter.
func NewSchemaConverter() SchemaConverter {
	return &ClickHouseConverter{}
}

// ConvertSchema converts a source schema to a ClickHouse schema.
func (c *ClickHouseConverter) ConvertSchema(source *models.Schema, targetType string) (*models.Schema, error) {
	target := &models.Schema{Databases: make([]models.DatabaseSchema, 0)}

	for _, db := range source.Databases {
		tgtDB := models.DatabaseSchema{Name: db.Name}
		// Assuming DatabaseSchema also updated to have slice of Tables
		tgtDB.Tables = make([]*models.TableSchema, 0, len(db.Tables))
		for _, tbl := range db.Tables {
			tgtTable := c.convertTable(tbl)
			tgtDB.Tables = append(tgtDB.Tables, tgtTable)
		}
		target.Databases = append(target.Databases, tgtDB)
	}
	return target, nil
}

func (c *ClickHouseConverter) convertTable(src *models.TableSchema) *models.TableSchema {
	var tgtCols []*models.ColumnSchema
	var orderByCols []string

	// 1. Column conversion
	for _, col := range src.Columns {
		chType := c.mapType(col.DataType)
		// Assuming we want to preserve PKs for ORDER BY
		if col.IsPrimaryKey {
			orderByCols = append(orderByCols, col.Name)
		}
		tgtCols = append(tgtCols, &models.ColumnSchema{
			Name:         col.Name,
			DataType:     chType,
			IsNullable:   col.IsNullable,
			IsPrimaryKey: col.IsPrimaryKey,
			Default:      col.Default, // You might need to adjust default values too
		})
	}

	// 2. Engine definition (Default: MergeTree)
	engineDef := "MergeTree()"
	if len(orderByCols) > 0 {
		engineDef += fmt.Sprintf(" ORDER BY (%s)", strings.Join(orderByCols, ", "))
	} else {
		engineDef += " ORDER BY tuple()"
	}

	return &models.TableSchema{
		Name:    src.Name,
		Columns: tgtCols,
		Engine:  engineDef,
		PK:      src.PK, // Keep original PKs metadata if needed, though Engine defines structure
		Indexes: src.Indexes,
	}
}

// mapType maps MySQL types to ClickHouse types.
func (c *ClickHouseConverter) mapType(mysqlType string) string {
	lowerType := strings.ToLower(mysqlType)
	// Handle parameterized types like decimal(10, 2)
	baseType := lowerType
	if idx := strings.Index(lowerType, "("); idx != -1 {
		baseType = lowerType[:idx]
	}

	switch baseType {
	case "tinyint":
		return "Int8"
	case "smallint":
		return "Int16"
	case "int", "integer", "mediumint":
		return "Int32"
	case "bigint":
		return "Int64"
	case "varchar", "char", "text", "mediumtext", "longtext":
		return "String"
	case "datetime", "timestamp":
		return "DateTime64" // Prefer DateTime64 for better precision if needed, or DateTime
	case "date":
		return "Date32"
	case "boolean", "bool":
		return "Bool"
	case "float":
		return "Float32"
	case "double":
		return "Float64"
	case "decimal":
		// Simple handling, just pass through or default.
		// For accurate mapping we need to extract precision/scale.
		// For now, if input was decimal(10,2), we return Decimal(10,2) or Decimal64(2)
		// The prompt example says: DECIMAL(10,2) -> Decimal64(2)
		return c.mapDecimal(lowerType)
	default:
		return "String" // Default fallback
	}
}

func (c *ClickHouseConverter) mapDecimal(fullType string) string {
	// fullType e.g. "decimal(10,2)" or "decimal"
	if !strings.Contains(fullType, "(") {
		return "Decimal128(18)" // Default
	}
	// Extract content inside parens
	start := strings.Index(fullType, "(")
	end := strings.LastIndex(fullType, ")")
	if start == -1 || end == -1 || end <= start {
		return "Decimal128(18)"
	}
	params := fullType[start+1 : end]
	parts := strings.Split(params, ",")
	if len(parts) == 2 {
		// scale is parts[1]
		scale := strings.TrimSpace(parts[1])
		// ClickHouse has Decimal32(S), Decimal64(S), Decimal128(S), Decimal256(S)
		// Or just Decimal(P, S) which is alias.
		// Prompt requested Decimal64(2) for input DECIMAL(10,2).
		// We can return Decimal64(scale).
		return fmt.Sprintf("Decimal64(%s)", scale)
	}
	return "Decimal128(18)"
}
