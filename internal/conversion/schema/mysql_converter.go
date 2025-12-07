package schema

import (
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/xwb1989/sqlparser"
)

// MySQLConverter implements SchemaConverter for MySQL.
type MySQLConverter struct {
	typeMapper    *IntelligentTypeMapper
	analyzer      *TypeAnalyzer
	precision     *PrecisionHandler
	warnings      *WarningCollector
	ruleLoader    *MappingRuleLoader
    sourceDB      string
}

// NewMySQLConverter creates a new MySQLConverter.
func NewMySQLConverter() *MySQLConverter {
    return NewMySQLConverterWithSource("mysql")
}

// NewMySQLConverterWithSource creates a new MySQLConverter with specific source DB.
func NewMySQLConverterWithSource(sourceDB string) *MySQLConverter {
	ruleLoader, _ := NewMappingRuleLoader("configs/type_mapping_rules.yaml")
	analyzer := NewTypeAnalyzer()
	precision := NewPrecisionHandler("configs/precision_policy.yaml")
	warnings := NewWarningCollector()

	typeMapper := NewIntelligentTypeMapper(ruleLoader, analyzer, precision)

	c := &MySQLConverter{
		typeMapper: typeMapper,
		analyzer:   analyzer,
		precision:  precision,
		warnings:   warnings,
		ruleLoader: ruleLoader,
        sourceDB:   sourceDB,
	}

	return c
}

// ConvertDDL converts MySQL DDL to target DB format.
func (c *MySQLConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
	stmts := strings.Split(sourceDDL, ";")
	var sb strings.Builder

	for _, stmtStr := range stmts {
		stmtStr = strings.TrimSpace(stmtStr)
		if stmtStr == "" {
			continue
		}

		stmt, err := sqlparser.Parse(stmtStr)
		if err != nil {
			if strings.HasPrefix(strings.ToUpper(stmtStr), "CREATE TABLE") {
				table, err := c.fallbackParseCreateTable(stmtStr)
				if err == nil {
					convertedTable, err := c.ConvertTable(table, targetDB)
					if err == nil {
						sb.WriteString(c.generateCreateSQL(convertedTable))
						sb.WriteString("\n\n")
						continue
					}
				}
			}

			utils.GetGlobalLogger().Warn("Failed to parse statement in MySQL converter", utils.Field{Key: "statement", Value: stmtStr}, utils.Field{Key: "error", Value: err})
			continue
		}

		switch ddl := stmt.(type) {
		case *sqlparser.DDL:
			if ddl.Action == sqlparser.CreateStr {
				if ddl.TableSpec == nil {
					table, err := c.fallbackParseCreateTable(stmtStr)
					if err == nil {
						convertedTable, err := c.ConvertTable(table, targetDB)
						if err == nil {
							sb.WriteString(c.generateCreateSQL(convertedTable))
							sb.WriteString("\n\n")
							continue
						}
					}
					utils.GetGlobalLogger().Warn("Failed to parse table spec and fallback failed", utils.Field{Key: "statement", Value: stmtStr})
					continue
				}

				tableSchema, err := c.parseTableSpec(ddl.NewName.Name.String(), ddl.TableSpec)
				if err != nil {
					return "", err
				}

				convertedTable, err := c.ConvertTable(tableSchema, targetDB)
				if err != nil {
					return "", err
				}

				sb.WriteString(c.generateCreateSQL(convertedTable))
				sb.WriteString("\n\n")
			}
		default:
		}
	}

	return sb.String(), nil
}

func (c *MySQLConverter) fallbackParseCreateTable(sql string) (*models.TableSchema, error) {
	sql = strings.TrimSpace(sql)
	parts := strings.SplitN(sql, "(", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid create table")
	}
	preamble := parts[0]
	nameParts := strings.Fields(preamble)
	tableName := nameParts[len(nameParts)-1]

	body := parts[1]
	body = strings.TrimSuffix(body, ")")
	body = strings.TrimSuffix(body, ";")

	cols := SplitWithBalance(body, ',')
	var columns []*models.ColumnSchema

	for _, colStr := range cols {
		colStr = strings.TrimSpace(colStr)
		colParts := strings.Fields(colStr)
		if len(colParts) < 2 {
			continue
		}
		name := colParts[0]
		typeStr := strings.Join(colParts[1:], " ")

		columns = append(columns, &models.ColumnSchema{
			Name: name,
			DataType: typeStr,
			IsNullable: true,
		})
	}

	return &models.TableSchema{
		Name: tableName,
		Columns: columns,
	}, nil
}

func (c *MySQLConverter) parseTableSpec(tableName string, spec *sqlparser.TableSpec) (*models.TableSchema, error) {
	var columns []*models.ColumnSchema
	var pk []string
	indexes := make(map[string]*models.IndexSchema)

	for _, col := range spec.Columns {
		colName := col.Name.String()
		colType := col.Type.Type

		fullTypeStr := colType
		if col.Type.Length != nil {
			fullTypeStr += fmt.Sprintf("(%s)", col.Type.Length.Val)
		}
		if col.Type.Scale != nil {
			fullTypeStr += fmt.Sprintf(",%s", col.Type.Scale.Val)
		}
		if col.Type.Unsigned {
			fullTypeStr += " UNSIGNED"
		}

		if len(col.Type.EnumValues) > 0 {
			var vals []string
			for _, v := range col.Type.EnumValues {
				vals = append(vals, fmt.Sprintf("'%s'", v))
			}
			fullTypeStr = fmt.Sprintf("ENUM(%s)", strings.Join(vals, ","))
		}

		isNullable := true
		if col.Type.NotNull {
			isNullable = false
		}

		isPrimaryKey := false
		if col.Type.KeyOpt == 1 {
			isPrimaryKey = true
		}

		columns = append(columns, &models.ColumnSchema{
			Name:         colName,
			DataType:     fullTypeStr,
			IsNullable:   isNullable,
			IsPrimaryKey: isPrimaryKey,
		})
	}

	for _, idx := range spec.Indexes {
		if idx.Info.Primary {
			for _, col := range idx.Columns {
				pk = append(pk, col.Column.String())
			}
		} else {
			idxName := idx.Info.Name.String()
			var idxCols []string
			for _, col := range idx.Columns {
				idxCols = append(idxCols, col.Column.String())
			}
			indexes[idxName] = &models.IndexSchema{
				Name:    idxName,
				Columns: idxCols,
				IsUnique: idx.Info.Unique,
			}
		}
	}

	for _, col := range columns {
		if col.IsPrimaryKey {
			found := false
			for _, p := range pk {
				if strings.EqualFold(col.Name, p) {
					found = true
					break
				}
			}
			if !found {
				pk = append(pk, col.Name)
			}
		}
	}

	for _, p := range pk {
		for _, col := range columns {
			if strings.EqualFold(col.Name, p) {
				col.IsPrimaryKey = true
			}
		}
	}

	return &models.TableSchema{
		Name:    tableName,
		Columns: columns,
		PK:      pk,
		Indexes: indexes,
	}, nil
}

// ConvertTable converts a single table structure.
func (c *MySQLConverter) ConvertTable(sourceTable *models.TableSchema, targetDB string) (*models.TableSchema, error) {
	targetTable := &models.TableSchema{
		Name:    sourceTable.Name,
		PK:      sourceTable.PK,
		Indexes: make(map[string]*models.IndexSchema),
	}

    tableCtx := &TableContext{
        TableName: sourceTable.Name,
    }

	for _, col := range sourceTable.Columns {
        ctx := &TypeMappingContext{
            SourceType: col.DataType,
            SourceDB: c.sourceDB,
            TargetDB: targetDB,
            ColumnName: col.Name,
            IsPrimaryKey: col.IsPrimaryKey,
            IsNullable: col.IsNullable,
            DefaultValue: col.Default,
            TableContext: tableCtx,
        }

        result, err := c.typeMapper.MapType(ctx)
        if err != nil {
             utils.GetGlobalLogger().Warn(fmt.Sprintf("Failed to map type for column %s: %v", col.Name, err))
             result = &TypeMappingResult{TargetType: "String"}
        }

        for _, w := range result.Warnings {
            w.AffectedColumn = fmt.Sprintf("%s.%s", sourceTable.Name, col.Name)
            c.warnings.Add(w)
        }

        targetType := result.TargetType
		if strings.HasPrefix(strings.ToUpper(col.DataType), "ENUM") && targetDB == "clickhouse" {
			targetType = c.convertEnumToClickHouse(col.DataType)
		}

		targetTable.Columns = append(targetTable.Columns, &models.ColumnSchema{
			Name:         col.Name,
			DataType:     targetType,
			IsNullable:   col.IsNullable && !col.IsPrimaryKey,
			IsPrimaryKey: col.IsPrimaryKey,
			Default:      col.Default,
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
	} else if targetDB == "mock_plugin" {
		targetTable.Engine = "MockEngine"
	}

	return targetTable, nil
}

func (c *MySQLConverter) convertEnumToClickHouse(enumType string) string {
	start := strings.Index(enumType, "(")
	end := strings.LastIndex(enumType, ")")
	if start == -1 || end == -1 {
		return "String"
	}

	content := enumType[start+1 : end]
	parts := strings.Split(content, ",")
	var enumDefs []string
	for i, p := range parts {
		val := strings.TrimSpace(p)
		enumDefs = append(enumDefs, fmt.Sprintf("%s=%d", val, i+1))
	}

	if len(enumDefs) > 127 {
		return fmt.Sprintf("Enum16(%s)", strings.Join(enumDefs, ", "))
	}
	return fmt.Sprintf("Enum8(%s)", strings.Join(enumDefs, ", "))
}

// GetTypeMapping gets the type mapping.
func (c *MySQLConverter) GetTypeMapping(sourceType string, targetDB string) (string, error) {
    ctx := &TypeMappingContext{
        SourceType: sourceType,
        SourceDB: c.sourceDB,
        TargetDB: targetDB,
    }
    res, err := c.typeMapper.MapType(ctx)
    if err != nil {
        return "", err
    }
    return res.TargetType, nil
}

func (c *MySQLConverter) generateCreateSQL(table *models.TableSchema) string {
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

// GenerateWarningReport generates warnings.
func (c *MySQLConverter) GenerateWarningReport(format string) (string, error) {
    return c.warnings.GenerateReport(format)
}
