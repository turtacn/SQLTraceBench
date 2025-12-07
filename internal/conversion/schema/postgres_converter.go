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
    typeMapper    *IntelligentTypeMapper
	analyzer      *TypeAnalyzer
	precision     *PrecisionHandler
	warnings      *WarningCollector
	ruleLoader    *MappingRuleLoader
}

// NewPostgresConverter creates a new PostgresConverter.
func NewPostgresConverter() *PostgresConverter {
	ruleLoader, _ := NewMappingRuleLoader("configs/type_mapping_rules.yaml")
	analyzer := NewTypeAnalyzer()
	precision := NewPrecisionHandler("configs/precision_policy.yaml")
	warnings := NewWarningCollector()

	typeMapper := NewIntelligentTypeMapper(ruleLoader, analyzer, precision)

	c := &PostgresConverter{
		typeMapper: typeMapper,
		analyzer:   analyzer,
		precision:  precision,
		warnings:   warnings,
		ruleLoader: ruleLoader,
	}
	return c
}

// ConvertDDL converts Postgres DDL to target DB format.
func (c *PostgresConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
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

		if strings.HasPrefix(upper, "CONSTRAINT") || strings.HasPrefix(upper, "PRIMARY KEY") {
			if strings.Contains(upper, "PRIMARY KEY") {
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

    tableCtx := &TableContext{
        TableName: sourceTable.Name,
    }

	for _, col := range sourceTable.Columns {
		targetType := ""

        ctx := &TypeMappingContext{
            SourceType: col.DataType,
            SourceDB: "postgres",
            TargetDB: targetDB,
            ColumnName: col.Name,
            IsPrimaryKey: col.IsPrimaryKey,
            IsNullable: col.IsNullable,
            TableContext: tableCtx,
        }

        result, err := c.typeMapper.MapType(ctx)
        if err == nil {
            targetType = result.TargetType
            for _, w := range result.Warnings {
                w.AffectedColumn = fmt.Sprintf("%s.%s", sourceTable.Name, col.Name)
                c.warnings.Add(w)
            }
        } else {
            // Fallback for Arrays logic
            if strings.HasSuffix(col.DataType, "[]") {
                baseType := strings.TrimSuffix(col.DataType, "[]")
                ctx.SourceType = baseType
                res, err := c.typeMapper.MapType(ctx)
                if err == nil {
                    targetType = fmt.Sprintf("Array(%s)", res.TargetType)
                } else {
                    targetType = "Array(String)"
                }
            } else {
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
    ctx := &TypeMappingContext{
        SourceType: sourceType,
        SourceDB: "postgres",
        TargetDB: targetDB,
    }
    res, err := c.typeMapper.MapType(ctx)
    if err != nil {
        return "", err
    }
    return res.TargetType, nil
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
