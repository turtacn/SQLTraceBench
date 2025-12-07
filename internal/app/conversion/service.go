package conversion

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/conversion/schema"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

// ConvertRequest represents a request to convert a schema.
type ConvertRequest struct {
	SourceSchemaPath string
	TargetDBType     string
	OutputPath       string
	SourceDB         string // Optional, will be auto-detected if empty
}

// ConvertTraceRequest represents a request to convert traces.
type ConvertTraceRequest struct {
	SourcePath   string
	TargetDBType string
}

// ConversionResult holds the result of a trace conversion.
type ConversionResult struct {
	Traces    []models.SQLTrace
	Templates []models.SQLTemplate
}

// Service is the interface for the conversion service.
type Service interface {
	ConvertFromFile(ctx context.Context, req ConvertTraceRequest) (*ConversionResult, error)
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

// ConvertFromFile reads SQL traces from a file, optionally translates them, and converts them to templates.
func (s *DefaultService) ConvertFromFile(ctx context.Context, req ConvertTraceRequest) (*ConversionResult, error) {
	file, err := os.Open(req.SourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var traces []models.SQLTrace
	parser := parsers.NewStreamingTraceParser(0) // Default buffer size

	// Prepare plugin if translation is needed
	var plugin interface {
		TranslateQuery(string) (string, error)
	}
	if req.TargetDBType != "" {
		p, ok := s.pluginRegistry.Get(req.TargetDBType)
		if !ok {
			return nil, fmt.Errorf("plugin not found: %s", req.TargetDBType)
		}
		plugin = p
	}

	err = parser.Parse(file, func(trace models.SQLTrace) error {
		if plugin != nil {
			translated, err := plugin.TranslateQuery(trace.Query)
			if err == nil {
				trace.Query = translated
			}
		}
		traces = append(traces, trace)
		return nil
	})
	if err != nil {
		return nil, err
	}

	tc := models.TraceCollection{Traces: traces}
	tpls := s.templateSvc.ExtractTemplates(tc)

	for i := range tpls {
		tables, err := s.parser.ListTables(tpls[i].RawSQL)
		if err != nil {
			continue
		}
		_ = tables
	}

	return &ConversionResult{
		Traces:    traces,
		Templates: tpls,
	}, nil
}

// ConvertStreamingly processes traces line-by-line using the provided callback.
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
	content, err := os.ReadFile(req.SourceSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	schemaContent := string(content)

	sourceDB := req.SourceDB
	if sourceDB == "" {
		sourceDB = detectSQLDialect(schemaContent)
	}

	factory := schema.NewConverterFactory()
	converter, err := factory.GetConverter(sourceDB)
	if err != nil {
		return fmt.Errorf("unsupported source database: %s", sourceDB)
	}

	convertedDDL, err := converter.ConvertDDL(schemaContent, req.TargetDBType)
	if err != nil {
		return fmt.Errorf("failed to convert DDL: %w", err)
	}

	return os.WriteFile(req.OutputPath, []byte(convertedDDL), 0644)
}

func detectSQLDialect(ddl string) string {
	upperDDL := strings.ToUpper(ddl)
	if strings.Contains(upperDDL, "ENGINE=INNODB") || strings.Contains(upperDDL, "ENGINE=MYISAM") {
		return "mysql"
	}
	if strings.Contains(upperDDL, "JSONB") || strings.Contains(upperDDL, "SERIAL") || strings.Contains(ddl, "::") {
		return "postgres"
	}
	if strings.Contains(upperDDL, "SHARD_ROW_ID_BITS") || strings.Contains(ddl, "/*T![clustered") {
		return "tidb"
	}
	// Default to mysql as it is most common and standard DDL often looks like MySQL
	return "mysql"
}
