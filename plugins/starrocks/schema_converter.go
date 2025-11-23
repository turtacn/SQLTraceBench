package starrocks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// StarRocksConverter handles schema conversion from generic/MySQL to StarRocks.
type StarRocksConverter struct{}

// ConvertSchema converts the source schema to a StarRocks compatible schema.
func (c *StarRocksConverter) ConvertSchema(source *models.Schema) (*models.Schema, error) {
	newSchema := &models.Schema{
		Databases: make([]models.DatabaseSchema, len(source.Databases)),
	}

	for i, db := range source.Databases {
		newSchema.Databases[i] = c.convertDatabase(db)
	}

	return newSchema, nil
}

func (c *StarRocksConverter) convertDatabase(db models.DatabaseSchema) models.DatabaseSchema {
	newDb := models.DatabaseSchema{
		Name:   db.Name,
		Tables: make([]*models.TableSchema, len(db.Tables)),
	}
	for i, tbl := range db.Tables {
		newDb.Tables[i] = c.convertTable(tbl)
	}
	return newDb
}

func (c *StarRocksConverter) convertTable(tbl *models.TableSchema) *models.TableSchema {
	newTbl := &models.TableSchema{
		Name:    tbl.Name,
		PK:      tbl.PK,
		Indexes: tbl.Indexes, // Keeping indexes as is for now
	}

	// 1. Convert Columns
	newTbl.Columns = make([]*models.ColumnSchema, len(tbl.Columns))
	for i, col := range tbl.Columns {
		newTbl.Columns[i] = c.convertColumn(col)
	}

	// 2. Construct Engine and Distribution
	// Default to OLAP engine
	// Use DUPLICATE KEY model by default as requested
	// DISTRIBUTED BY HASH(pk)

	var distKeys []string
	if len(tbl.PK) > 0 {
		distKeys = tbl.PK
	} else if len(newTbl.Columns) > 0 {
		// Fallback to first column if no PK
		distKeys = []string{newTbl.Columns[0].Name}
	}

	// Construct the Engine string to be appended to CREATE TABLE
	// Format: OLAP DUPLICATE KEY(...) DISTRIBUTED BY HASH(...) BUCKETS 10 PROPERTIES(...)
	// We will simplify to just the Engine type and keys. The DDL generator will need to handle properties if needed,
	// or we put everything here.

	// Example result: "OLAP DUPLICATE KEY(`id`) DISTRIBUTED BY HASH(`id`)"

	engineParts := []string{"OLAP"}

	// Helper to quote identifiers
	quote := func(keys []string) string {
		quoted := make([]string, len(keys))
		for i, k := range keys {
			quoted[i] = fmt.Sprintf("`%s`", k)
		}
		return strings.Join(quoted, ", ")
	}

	// DUPLICATE KEY
	// Using PK columns for duplicate key (sorting)
	if len(distKeys) > 0 {
		engineParts = append(engineParts, fmt.Sprintf("DUPLICATE KEY(%s)", quote(distKeys)))
	}

	// DISTRIBUTED BY
	if len(distKeys) > 0 {
		engineParts = append(engineParts, fmt.Sprintf("DISTRIBUTED BY HASH(%s)", quote(distKeys)))
	} else {
		// Minimal fallback
		engineParts = append(engineParts, "DISTRIBUTED BY RANDOM")
	}

	newTbl.Engine = strings.Join(engineParts, " ")

	return newTbl
}

func (c *StarRocksConverter) convertColumn(col *models.ColumnSchema) *models.ColumnSchema {
	newCol := *col // shallow copy
	newCol.DataType = c.mapType(col.DataType)
	return &newCol
}

func (c *StarRocksConverter) mapType(dataType string) string {
	upperType := strings.ToUpper(dataType)

	// Regex for type with parameters, e.g., VARCHAR(255)
	re := regexp.MustCompile(`^([A-Z]+)(?:\((.*)\))?.*$`)
	matches := re.FindStringSubmatch(upperType)

	var baseType string
	var params string
	if len(matches) > 1 {
		baseType = matches[1]
	} else {
		baseType = upperType
	}
	if len(matches) > 2 {
		params = matches[2]
	}

	switch baseType {
	case "TINYINT", "SMALLINT", "BIGINT", "DECIMAL", "DATE", "DATETIME":
		// Keep as is, optionally preserving params
		if params != "" {
			return fmt.Sprintf("%s(%s)", baseType, params)
		}
		return baseType
	case "INT", "INTEGER", "MEDIUMINT":
		// StarRocks uses INT (4 bytes)
		// MEDIUMINT in MySQL is 3 bytes, StarRocks maps to INT
		return "INT"
	case "FLOAT":
		return "FLOAT"
	case "DOUBLE":
		return "DOUBLE"
	case "CHAR":
		if params != "" {
			return fmt.Sprintf("CHAR(%s)", params)
		}
		return "CHAR"
	case "VARCHAR":
		if params != "" {
			return fmt.Sprintf("VARCHAR(%s)", params)
		}
		return "VARCHAR(65533)" // Default max if unknown? Or just VARCHAR
	case "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT", "BLOB", "LONGBLOB":
		return "STRING"
	case "TIMESTAMP":
		return "DATETIME" // Map TIMESTAMP to DATETIME
	case "BOOL", "BOOLEAN":
		return "BOOLEAN"
	default:
		// Fallback for types we missed or are compatible
		if params != "" {
			return fmt.Sprintf("%s(%s)", baseType, params)
		}
		return baseType
	}
}
