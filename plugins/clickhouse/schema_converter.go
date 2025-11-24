package clickhouse

import (
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaConverter is the interface for schema conversion logic.
type SchemaConverter interface {
	ConvertSchema(source *models.Schema) (*models.Schema, error)
}

type ClickHouseConverter struct{}

// NewSchemaConverter creates a new SchemaConverter.
func NewSchemaConverter() SchemaConverter {
	return &ClickHouseConverter{}
}

// ConvertSchema converts a source schema to a ClickHouse schema.
func (c *ClickHouseConverter) ConvertSchema(source *models.Schema) (*models.Schema, error) {
	target := &models.Schema{Databases: make([]models.DatabaseSchema, 0)}

	for _, db := range source.Databases {
		tgtDB := models.DatabaseSchema{Name: db.Name}
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
		chType, err := c.mapType(col.DataType)
		if err != nil {
			// Fallback to String for unknown types as per requirement implicitly or just log?
			// Requirement says "DEFAULT -> String (safe fallback)"
			chType = "String"
		}

		if col.IsPrimaryKey {
			orderByCols = append(orderByCols, col.Name)
		}
		tgtCols = append(tgtCols, &models.ColumnSchema{
			Name:         col.Name,
			DataType:     chType,
			IsNullable:   col.IsNullable,
			IsPrimaryKey: col.IsPrimaryKey,
			Default:      col.Default,
		})
	}

	// 2. Engine definition (Default: MergeTree)
	engineDef := "MergeTree()"
	if len(orderByCols) > 0 {
		engineDef += fmt.Sprintf(" ORDER BY (%s)", strings.Join(orderByCols, ", "))
	} else {
		engineDef += " ORDER BY tuple()"
	}

	tbl := &models.TableSchema{
		Name:    src.Name,
		Columns: tgtCols,
		Engine:  engineDef,
		PK:      src.PK,
		Indexes: src.Indexes,
	}

	// 3. Build CreateSQL
	tbl.CreateSQL = buildCHCreateSQL(tbl, orderByCols)
	return tbl
}

// mapType maps MySQL types to ClickHouse types.
func (c *ClickHouseConverter) mapType(mysqlType string) (string, error) {
	lowerType := strings.ToLower(mysqlType)
	baseType := lowerType
	if idx := strings.Index(lowerType, "("); idx != -1 {
		baseType = lowerType[:idx]
	}
	baseType = strings.TrimSpace(baseType)

	switch baseType {
	case "tinyint":
		return "Int8", nil
	case "int", "integer":
		return "Int32", nil
	case "bigint":
		return "Int64", nil
	case "varchar", "text":
		return "String", nil
	case "datetime":
		return "DateTime", nil
	case "decimal":
		// DECIMAL(p,s) -> Decimal128(s)
		return c.mapDecimal(lowerType), nil
	default:
		// Requirement: DEFAULT -> String (safe fallback)
		return "String", nil
	}
}

func (c *ClickHouseConverter) mapDecimal(fullType string) string {
	// fullType e.g. "decimal(10,2)" or "decimal"
	if !strings.Contains(fullType, "(") {
		return "Decimal128(18)" // Default if no scale provided
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
		return fmt.Sprintf("Decimal128(%s)", scale)
	}
	return "Decimal128(18)"
}

func buildCHCreateSQL(tbl *models.TableSchema, pkCols []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tbl.Name))

	for i, col := range tbl.Columns {
		sb.WriteString(fmt.Sprintf("  `%s` %s", col.Name, col.DataType))
		// Handle Default? Not strictly asked in buildCHCreateSQL example but good to have
		if col.Default != "" {
			// Basic handling, might need more complex parsing for expressions
			// For string types, ensure quotes if not present?
			// Let's keep it simple as per requirements.
			// sb.WriteString(fmt.Sprintf(" DEFAULT %s", col.Default))
		}
		if i < len(tbl.Columns)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	sb.WriteString(fmt.Sprintf(") ENGINE = %s", tbl.Engine))
	return sb.String()
}
