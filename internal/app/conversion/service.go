package conversion

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

// ConvertRequest represents a request to convert a schema.
type ConvertRequest struct {
	SourceSchemaPath string
	TargetDBType     string
	OutputPath       string
}

// Service is the interface for the conversion service.
type Service interface {
	ConvertFromFile(ctx context.Context, tracePath string) ([]models.SQLTemplate, error)
	ConvertSchemaFromFile(ctx context.Context, req ConvertRequest) error
	ConvertStreamingly(ctx context.Context, tracePath string, bufferSize int, callback func(models.SQLTrace) error) error
}

// DefaultService is the default implementation of the conversion service.
type DefaultService struct {
	templateSvc    *services.TemplateService
	parser         services.Parser
	pluginRegistry *plugin_registry.Registry
}

// NewService creates a new DefaultService.
// Note: If registry is not passed, we use the global one or create one.
// To keep signature compatible if possible, or we assume caller updates.
// For now, I'll use the GlobalRegistry if not passed via deps (but constructor here builds deps).
// Better to inject it.
func NewService(parser services.Parser, registry *plugin_registry.Registry) Service {
	if registry == nil {
		registry = plugin_registry.GlobalRegistry
	}
	return &DefaultService{
		templateSvc:    services.NewTemplateService(),
		parser:         parser,
		pluginRegistry: registry,
	}
}

// ConvertFromFile reads SQL traces from a file and converts them to templates.
func (s *DefaultService) ConvertFromFile(ctx context.Context, tracePath string) ([]models.SQLTemplate, error) {
	file, err := os.Open(tracePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var traces []models.SQLTrace
	parser := parsers.NewStreamingTraceParser(0) // Default buffer size

	err = parser.Parse(file, func(trace models.SQLTrace) error {
		traces = append(traces, trace)
		return nil
	})
	if err != nil {
		return nil, err
	}

	tc := models.TraceCollection{Traces: traces}
	tpls := s.templateSvc.ExtractTemplates(tc)

	// Use the parser to extract table names for each template.
	for i := range tpls {
		tables, err := s.parser.ListTables(tpls[i].RawSQL)
		if err != nil {
			// In a real application, we might want to handle this error more gracefully.
			continue
		}
		_ = tables // TODO: store the tables in the template
	}

	return tpls, nil
}

// ConvertStreamingly processes traces line-by-line using the provided callback.
// This prevents OOM errors by avoiding loading all traces into memory.
func (s *DefaultService) ConvertStreamingly(ctx context.Context, tracePath string, bufferSize int, callback func(models.SQLTrace) error) error {
	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}
	defer file.Close()

	parser := parsers.NewStreamingTraceParser(bufferSize)

	return parser.Parse(file, callback)
}

// ConvertSchemaFromFile reads a SQL schema file, converts it to the target dialect, and writes the result.
func (s *DefaultService) ConvertSchemaFromFile(ctx context.Context, req ConvertRequest) error {
	// 1. Parse Source
	srcSchema, err := parseDDLFile(req.SourceSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to parse source schema: %w", err)
	}

	// 2. Get Plugin from Registry
	p, ok := s.pluginRegistry.Get(req.TargetDBType)
	if !ok {
		// Try loading from global if not found in local ref (though they should be same usually)
		// Or try to lazy load?
		// For now, assume it's preloaded or we error.
		return fmt.Errorf("plugin not found for target db type: %s", req.TargetDBType)
	}

	// 3. Convert Schema
	tgtSchema, err := p.ConvertSchema(srcSchema)
	if err != nil {
		return fmt.Errorf("plugin failed to convert schema: %w", err)
	}

	// 4. Generate DDL String & Write
	// Note: We use the `CreateSQL` field if present (StarRocks plugin populates it),
	// otherwise we fallback to generic generation.
	ddl := s.generateDDL(tgtSchema)
	return os.WriteFile(req.OutputPath, []byte(ddl), 0644)
}

func parseDDLFile(path string) (*models.Schema, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sql := string(content)

	// Very simple regex parser for CREATE TABLE
	schema := &models.Schema{
		Databases: []models.DatabaseSchema{
			{
				Name:   "default", // Extracted from file or context
				Tables: []*models.TableSchema{},
			},
		},
	}

	// Split by semicolon
	stmts := strings.Split(sql, ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(stmt), "CREATE TABLE") {
			table, err := parseCreateTable(stmt)
			if err != nil {
				utils.GetGlobalLogger().Error("Failed to parse table", utils.Field{Key: "error", Value: err})
				continue
			}
			schema.Databases[0].Tables = append(schema.Databases[0].Tables, table)
		}
	}
	return schema, nil
}

func parseCreateTable(sql string) (*models.TableSchema, error) {
	// regex to get table name: CREATE TABLE `?name`?
	reName := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?[\x60"]?(\w+)[\x60"]?`)
	matches := reName.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not extract table name")
	}
	tableName := matches[1]

	// Extract body
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
		// Check for PRIMARY KEY (col, ...)
		if strings.HasPrefix(strings.ToUpper(colStr), "PRIMARY KEY") {
			// Extract keys
			rePK := regexp.MustCompile(`(?i)PRIMARY\s+KEY\s*\(([^)]+)\)`)
			pkMatches := rePK.FindStringSubmatch(colStr)
			if len(pkMatches) >= 2 {
				keys := strings.Split(pkMatches[1], ",")
				for _, k := range keys {
					k = strings.TrimSpace(k)
					k = strings.Trim(k, "`\"")
					pks = append(pks, k)
				}
			}
			continue
		}
		// Check for other keys like KEY, INDEX, CONSTRAINT
		upper := strings.ToUpper(colStr)
		if strings.HasPrefix(upper, "KEY") || strings.HasPrefix(upper, "INDEX") || strings.HasPrefix(upper, "CONSTRAINT") || strings.HasPrefix(upper, "UNIQUE") {
			continue
		}

		// Assume it's a column definition
		// format: name type [modifiers]
		parts := strings.Fields(colStr)
		if len(parts) < 2 {
			continue
		}
		name := strings.Trim(parts[0], "`\"")
		dataType := parts[1]

		isNullable := true
		isPrimaryKey := false

		// check modifiers
		upperStr := strings.ToUpper(colStr)
		if strings.Contains(upperStr, "NOT NULL") {
			isNullable = false
		}
		if strings.Contains(upperStr, "PRIMARY KEY") {
			isPrimaryKey = true
			pks = append(pks, name)
		}

		columns = append(columns, &models.ColumnSchema{
			Name:         name,
			DataType:     dataType,
			IsNullable:   isNullable,
			IsPrimaryKey: isPrimaryKey,
		})
	}

	// Mark IsPrimaryKey on columns if found in table-level PK definition
	for _, pk := range pks {
		for _, col := range columns {
			if col.Name == pk {
				col.IsPrimaryKey = true
			}
		}
	}

	return &models.TableSchema{
		Name:    tableName,
		Columns: columns,
		PK:      pks,
	}, nil
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

func (s *DefaultService) generateDDL(schema *models.Schema) string {
	var sb strings.Builder
	for _, db := range schema.Databases {
		for _, table := range db.Tables {
			// If CreateSQL is already populated by the plugin, use it.
			if table.CreateSQL != "" {
				sb.WriteString(table.CreateSQL)
				sb.WriteString("\n\n")
				continue
			}

			// Fallback generic generator (mostly for legacy/ClickHouse if not updated)
			sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", table.Name))
			for i, col := range table.Columns {
				sb.WriteString(fmt.Sprintf("    %s %s", col.Name, col.DataType))
				if i < len(table.Columns)-1 {
					sb.WriteString(",")
				}
				sb.WriteString("\n")
			}
			if table.Engine != "" {
				sb.WriteString(fmt.Sprintf(") ENGINE = %s;\n\n", table.Engine))
			} else {
				sb.WriteString(");\n\n")
			}
		}
	}
	return sb.String()
}
