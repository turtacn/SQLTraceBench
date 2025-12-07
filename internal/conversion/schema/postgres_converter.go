package schema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

// PostgresConverter implements SchemaConverter for PostgreSQL.
type PostgresConverter struct {
}

// NewPostgresConverter creates a new PostgresConverter.
func NewPostgresConverter() *PostgresConverter {
	return &PostgresConverter{}
}

// ConvertDDL converts Postgres DDL to target DB format.
func (c *PostgresConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
	// Simple split by semicolon (naive approach, assumes no semicolon in strings/comments)
	// For robust implementation, a better scanner is needed.
	stmts := strings.Split(sourceDDL, ";")
	var sb strings.Builder

	for _, stmtStr := range stmts {
		stmtStr = strings.TrimSpace(stmtStr)
		if stmtStr == "" {
			continue
		}

		if strings.HasPrefix(strings.ToUpper(stmtStr), "CREATE TABLE") {
			tableSchema, err := c.parseCreateTable(stmtStr)
			if err != nil {
				utils.GetGlobalLogger().Warn("Failed to parse PG table", utils.Field{Key: "error", Value: err})
				continue
			}

			convertedTable, err := c.ConvertTable(tableSchema, targetDB)
			if err != nil {
				return "", err
			}

			sb.WriteString(c.generateCreateSQL(convertedTable))
			sb.WriteString("\n\n")
		}
	}

	return sb.String(), nil
}

func (c *PostgresConverter) parseCreateTable(sql string) (*models.TableSchema, error) {
	// Basic regex parsing logic similar to legacy service but adapted for PG quirks
	reName := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?["']?(\w+)["']?`)
	matches := reName.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not extract table name")
	}
	tableName := matches[1]

	start := strings.Index(sql, "(")
	end := strings.LastIndex(sql, ")")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("could not extract table body")
	}
	body := sql[start+1 : end]

	colsStr := splitWithBalance(body, ',')

	var columns []*models.ColumnSchema
	var pks []string

	for _, colStr := range colsStr {
		colStr = strings.TrimSpace(colStr)
		if colStr == "" {
			continue
		}

		upper := strings.ToUpper(colStr)

		// Constraint handling (PK, FK, etc.)
		if strings.HasPrefix(upper, "CONSTRAINT") || strings.HasPrefix(upper, "PRIMARY KEY") {
			if strings.Contains(upper, "PRIMARY KEY") {
				// Extract PKs
				rePK := regexp.MustCompile(`(?i)PRIMARY\s+KEY\s*\(([^)]+)\)`)
				pkMatches := rePK.FindStringSubmatch(colStr)
				if len(pkMatches) >= 2 {
					keys := strings.Split(pkMatches[1], ",")
					for _, k := range keys {
						k = strings.TrimSpace(k)
						k = strings.Trim(k, "\"'")
						pks = append(pks, k)
					}
				}
			}
			continue
		}

		parts := strings.Fields(colStr)
		if len(parts) < 2 {
			continue
		}
		name := strings.Trim(parts[0], "\"'")

		// Reconstruct type
		// PG types can be: INTEGER, VARCHAR(20), INTEGER[], JSONB, etc.
		// parts[1] is start of type, but type might have spaces (DOUBLE PRECISION)
		// We need to parse until constraints like NOT NULL, PRIMARY KEY, DEFAULT

		remaining := strings.Join(parts[1:], " ")
		typeEnd := len(remaining)

		constraints := []string{"NOT NULL", "NULL", "PRIMARY KEY", "DEFAULT", "REFERENCES", "CHECK", "UNIQUE"}
		lowerRemaining := strings.ToUpper(remaining)

		for _, cw := range constraints {
			idx := strings.Index(lowerRemaining, cw)
			if idx != -1 && idx < typeEnd {
				typeEnd = idx
			}
		}

		dataType := strings.TrimSpace(remaining[:typeEnd])
		isNullable := !strings.Contains(upper, "NOT NULL")
		isPrimaryKey := strings.Contains(upper, "PRIMARY KEY")

		if isPrimaryKey {
			pks = append(pks, name)
		}

		columns = append(columns, &models.ColumnSchema{
			Name:         name,
			DataType:     dataType,
			IsNullable:   isNullable,
			IsPrimaryKey: isPrimaryKey,
		})
	}

	return &models.TableSchema{
		Name:    tableName,
		Columns: columns,
		PK:      pks,
	}, nil
}

// ConvertTable converts a single table structure.
func (c *PostgresConverter) ConvertTable(sourceTable *models.TableSchema, targetDB string) (*models.TableSchema, error) {
	targetTable := &models.TableSchema{
		Name:    sourceTable.Name,
		PK:      sourceTable.PK,
		Indexes: make(map[string]*models.IndexSchema),
	}

	for _, col := range sourceTable.Columns {
		targetType := ""

		// Handle Array types
		if strings.HasSuffix(col.DataType, "[]") {
			baseType := strings.TrimSuffix(col.DataType, "[]")
			mappedBase, err := c.GetTypeMapping(baseType, targetDB)
			if err == nil {
				targetType = fmt.Sprintf("Array(%s)", mappedBase)
			} else {
				targetType = "Array(String)"
			}
		} else {
			var err error
			targetType, err = c.GetTypeMapping(col.DataType, targetDB)
			if err != nil {
				utils.GetGlobalLogger().Warn(fmt.Sprintf("Unknown type '%s' in table '%s', fallback to String", col.DataType, sourceTable.Name))
				targetType = "String"
			}
		}

		targetTable.Columns = append(targetTable.Columns, &models.ColumnSchema{
			Name:         col.Name,
			DataType:     targetType,
			IsNullable:   col.IsNullable,
			IsPrimaryKey: col.IsPrimaryKey,
		})
	}

	if targetDB == "clickhouse" {
		engine := "MergeTree()"
		orderBy := "tuple()"
		if len(sourceTable.PK) > 0 {
			orderBy = fmt.Sprintf("(%s)", strings.Join(sourceTable.PK, ", "))
		}
		engine += fmt.Sprintf(" ORDER BY %s", orderBy)
		targetTable.Engine = engine
	} else if targetDB == "starrocks" {
		targetTable.Engine = "OLAP"
	}

	return targetTable, nil
}

// GetTypeMapping gets the type mapping.
func (c *PostgresConverter) GetTypeMapping(sourceType string, targetDB string) (string, error) {
	baseType := getBaseType(sourceType)
	var mapping TypeMapping
	if targetDB == "clickhouse" {
		mapping = postgresToClickHouseTypeMap
	} else if targetDB == "starrocks" {
		mapping = postgresToStarRocksTypeMap
	} else {
		return "", fmt.Errorf("unsupported target database: %s", targetDB)
	}

	if val, ok := mapping[baseType]; ok {
		if baseType == "DECIMAL" || baseType == "NUMERIC" {
             // Extract params if available
             params := ""
             start := strings.Index(sourceType, "(")
             end := strings.LastIndex(sourceType, ")")
             if start != -1 && end != -1 {
                 params = sourceType[start : end+1]
             }
             return val + params, nil
        }
		return val, nil
	}
	return "", fmt.Errorf("type mapping not found for %s", sourceType)
}

func (c *PostgresConverter) generateCreateSQL(table *models.TableSchema) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", table.Name))
	for i, col := range table.Columns {
		sb.WriteString(fmt.Sprintf("    %s %s", col.Name, col.DataType))
		if i < len(table.Columns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	if table.Engine != "" {
		sb.WriteString(fmt.Sprintf(") ENGINE = %s;", table.Engine))
	} else {
		sb.WriteString(");")
	}
	return sb.String()
}

// Helper split function (same as in service.go, but duplicated here to avoid cyclic dep if moved to utils)
func splitWithBalance(s string, sep rune) []string {
	var parts []string
	var current strings.Builder
	balance := 0
	for _, r := range s {
		if r == '(' {
			balance++
		} else if r == ')' {
			balance--
		}
		if r == sep && balance == 0 {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
