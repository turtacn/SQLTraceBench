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
}

// NewMySQLConverter creates a new MySQLConverter.
func NewMySQLConverter() *MySQLConverter {
	return &MySQLConverter{}
}

// ConvertDDL converts MySQL DDL to target DB format.
func (c *MySQLConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
	// sqlparser only supports parsing a single statement at a time properly,
	// or we need to split it.
	// But since the sourceDDL might contain multiple statements (whole schema file),
	// we should split it.
	stmts := strings.Split(sourceDDL, ";")
	var sb strings.Builder

	for _, stmtStr := range stmts {
		stmtStr = strings.TrimSpace(stmtStr)
		if stmtStr == "" {
			continue
		}

		// Parse the statement
		stmt, err := sqlparser.Parse(stmtStr)
		if err != nil {
			// Fallback logic if parser fails
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
			if ddl.Action == sqlparser.CreateStr { // create table
				// Handle case where TableSpec is nil but err was nil (e.g. ENUM issues)
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
					// If fallback fails, just continue or log
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
			// Ignore other statements for now
		}
	}

	return sb.String(), nil
}

func (c *MySQLConverter) fallbackParseCreateTable(sql string) (*models.TableSchema, error) {
	// Simple regex parser for ENUM fallback
	// Extract table name
	sql = strings.TrimSpace(sql)
	parts := strings.SplitN(sql, "(", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid create table")
	}
	preamble := parts[0]
	nameParts := strings.Fields(preamble)
	tableName := nameParts[len(nameParts)-1] // rough guess

	// Extract body
	body := parts[1]
	body = strings.TrimSuffix(body, ")")
	body = strings.TrimSuffix(body, ";")

	// Split columns
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

		// Handle parameterized types
		if col.Type.Length != nil {
			// e.g. VARCHAR(255)
			// val := string(col.Type.Length.Val) // This is complicated in sqlparser AST
			// For simplicity, we just keep the base type and reconstruct if needed,
			// but converting to models.ColumnSchema stores raw DataType string usually.
			// Let's reconstruct the type string.
			// sqlparser doesn't easily give back the full string.
			// We will just use the base type for mapping or try to approximate.
		}

		// To get the full type string including length, we might need to buffer parsing.
		// Or simpler:
		// fullType := colType
		// This is a limitation of using the AST directly without a formatter.
		// However, for the purpose of this task, let's assume we can map the base type
		// and maybe we miss length for now unless we do extra work.
		// Wait, the prompt requirements said "supports 30+ common types".

		// Actually, let's look at how sqlparser stores types.
		// col.Type is ColumnType.
		// It has Length, Unsigned, Zerofill, etc.

		// Reconstruct full type string
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

		// Special handling for ENUM if values are available
		if len(col.Type.EnumValues) > 0 {
			// Construct ENUM string: ENUM('a','b')
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

		// Check explicit PK from parsing
		isPrimaryKey := false
		if col.Type.KeyOpt == 1 { // sqlparser.ColumnKeyPrimary
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
			// Store other indexes
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

	// Collect PKs from columns if not already in pk list
	for _, col := range columns {
		if col.IsPrimaryKey {
			// Check if already in pk list
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

	// Update columns IsPrimaryKey based on collected pk list
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

	// 1. Convert Columns
	for _, col := range sourceTable.Columns {
		targetType, err := c.GetTypeMapping(col.DataType, targetDB)
		if err != nil {
			// Fallback to String and warn
			utils.GetGlobalLogger().Warn(fmt.Sprintf("Unknown type '%s' in table '%s', fallback to String", col.DataType, sourceTable.Name))
			targetType = "String"
			if targetDB == "starrocks" {
				targetType = "VARCHAR(65533)" // StarRocks String equivalent
			}
		}

		// Handle ENUM special case for ClickHouse
		if strings.HasPrefix(strings.ToUpper(col.DataType), "ENUM") && targetDB == "clickhouse" {
			// Need to extract values. The parseTableSpec put the full string "ENUM('a','b')" in DataType.
			targetType = c.convertEnumToClickHouse(col.DataType)
		} else if strings.HasPrefix(strings.ToUpper(col.DataType), "ENUM") && targetDB == "starrocks" {
			targetType = "VARCHAR(65533)" // Fallback for StarRocks
		}

		targetTable.Columns = append(targetTable.Columns, &models.ColumnSchema{
			Name:         col.Name,
			DataType:     targetType,
			IsNullable:   col.IsNullable,
			IsPrimaryKey: col.IsPrimaryKey,
			Default:      col.Default,
		})
	}

	// 2. Convert Engine / Options
	if targetDB == "clickhouse" {
		engine := "MergeTree()"
		// ORDER BY
		orderBy := "tuple()"
		if len(sourceTable.PK) > 0 {
			orderBy = fmt.Sprintf("(%s)", strings.Join(sourceTable.PK, ", "))
		}
		engine += fmt.Sprintf(" ORDER BY %s", orderBy)
		targetTable.Engine = engine
	} else if targetDB == "starrocks" {
		// StarRocks DDL logic
		// DUPLICATE KEY or UNIQUE KEY based on PK?
		// For simplicity, using DUPLICATE KEY with primary key columns
		targetTable.Engine = "OLAP"
		// Keys logic would be handled in generateCreateSQL or here.
		// models.TableSchema has Engine string field which is raw.
		// We can construct it here.
	}

	return targetTable, nil
}

func (c *MySQLConverter) convertEnumToClickHouse(enumType string) string {
	// Input: ENUM('pending','shipped')
	// Output: Enum8('pending'=1, 'shipped'=2)

	start := strings.Index(enumType, "(")
	end := strings.LastIndex(enumType, ")")
	if start == -1 || end == -1 {
		return "String"
	}

	content := enumType[start+1 : end]
	// Split by comma, respecting quotes
	// Simple split for now
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
	baseType := getBaseType(sourceType)
	var mapping TypeMapping
	if targetDB == "clickhouse" {
		mapping = mysqlToClickHouseTypeMap
	} else if targetDB == "starrocks" {
		mapping = mysqlToStarRocksTypeMap
	} else {
		return "", fmt.Errorf("unsupported target database: %s", targetDB)
	}

	if val, ok := mapping[baseType]; ok {
		// Handle parameterized types
		if baseType == "DECIMAL" || baseType == "CHAR" {
			// Extract parameters from sourceType
			params := ""
			start := strings.Index(sourceType, "(")
			end := strings.LastIndex(sourceType, ")")
			if start != -1 && end != -1 {
				params = sourceType[start : end+1]
			}
			if targetDB == "clickhouse" {
				if baseType == "CHAR" {
					return "FixedString" + params, nil
				}
				return "Decimal" + params, nil
			}
			return val + params, nil
		}
		return val, nil
	}

	return "", fmt.Errorf("type mapping not found for %s", sourceType)
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
