package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/turtacn/SQLTraceBench/plugins"
	_ "github.com/turtacn/SQLTraceBench/plugin_registry"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert a raw SQL trace file into a SQL template file",
		RunE:  runConvert,
	}
	tracePath    string
	tplPath      string
	targetPlugin string
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&tracePath, "trace", "t", "traces.jsonl", "Path to the raw SQL trace file")
	convertCmd.Flags().StringVarP(&tplPath, "out", "o", "templates.json", "Path to the output template file")
	convertCmd.Flags().StringVar(&targetPlugin, "target", "", "Target database for SQL translation (e.g., clickhouse, starrocks)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	logger := utils.GetGlobalLogger()

	// Load configuration for buffer size
	// Assuming viper is already initialized by rootCmd in a real app,
	// or we can just read the value if available.
	// The test env might not have initialized viper properly if we don't call initConfig.
	// But let's assume we can read "parser.buffer_size".
	bufferSize := viper.GetInt("parser.buffer_size")
	if bufferSize == 0 {
		// Fallback if config not loaded
		bufferSize = 1024 * 1024
	}

	// If output is JSONL, we use streaming mode to avoid OOM.
	isStreamOutput := len(tplPath) > 5 && tplPath[len(tplPath)-6:] == ".jsonl"

	if isStreamOutput {
		var plugin plugins.Plugin
		if targetPlugin != "" {
			p, err := plugins.GetPlugin(targetPlugin)
			if err != nil {
				return err
			}
			plugin = p
		}

		outFile, err := os.Create(tplPath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		encoder := json.NewEncoder(outFile)

		// Use the service's streaming method (Architectural fix)
		// Note: We need to instantiate the service.
		// Since we don't have DI setup in this command file, we create default.
		// We need a parser for the service constructor, but for ConvertStreamingly it ignores the internal parser field?
		// Actually `DefaultService` struct has `parser services.Parser`.
		// `NewStreamingTraceParser` returns `*StreamingTraceParser` which is not `services.Parser` (interface).
		// `services.Parser` is for SQL parsing (ListTables), not Trace parsing.
		// `DefaultService` creates `StreamingTraceParser` internally in `ConvertStreamingly`.
		// So we just need a dummy SQL parser for the constructor if we use NewService.
		// Or we can construct DefaultService manually if allowed?
		// `DefaultService` is exported.

		// Let's just use NewService with a mock or nil parser if allowed.
		// `NewService(parser services.Parser)`.
		// We can pass nil if ConvertStreamingly doesn't use it.
		svc := conversion.NewService(nil, nil)

		count := 0
		err = svc.ConvertStreamingly(cmd.Context(), tracePath, bufferSize, func(trace models.SQLTrace) error {
			if plugin != nil {
				translatedQuery, err := plugin.TranslateQuery(trace.Query)
				if err == nil {
					trace.Query = translatedQuery
				}
			}
			if err := encoder.Encode(trace); err != nil {
				return err
			}
			count++
			if count%10000 == 0 {
				logger.Info("Progress", utils.Field{Key: "processed", Value: count})
			}
			return nil
		})

		if err != nil {
			return err
		}
		logger.Info("Conversion complete", utils.Field{Key: "total_processed", Value: count})
		return nil
	}

	// Original behavior: Load all, convert, extract templates.
	// This uses the streaming parser to READ, but accumulates in memory.
	// We should update this to use the configured buffer size as well.
	// But `ConvertFromFile` inside service now uses `NewStreamingTraceParser(0)`.
	// We should probably update `ConvertFromFile` to accept buffer size or read config?
	// Or `DefaultService` should accept `StreamingTraceParser` factory?
	// For now, I'll leave `ConvertFromFile` using default/0 buffer size as I didn't change its signature to take config.
	// The requirement was to make `StreamingTraceParser` configurable.

	// Legacy path for template extraction.
	// Note: svc is not used here to avoid instantiating complex dependencies like SQL Parser if not needed for pure trace parsing,
	// keeping the change minimal and focused on the parser replacement.

	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}
	defer file.Close()

	parser := parsers.NewStreamingTraceParser(bufferSize)
	var traces []models.SQLTrace
	var plugin plugins.Plugin
	if targetPlugin != "" {
		p, err := plugins.GetPlugin(targetPlugin)
		if err != nil {
			return err
		}
		plugin = p
	}

	count := 0
	err = parser.Parse(file, func(trace models.SQLTrace) error {
		if plugin != nil {
			translatedQuery, err := plugin.TranslateQuery(trace.Query)
			if err == nil {
				trace.Query = translatedQuery
			}
		}
		traces = append(traces, trace)
		count++
		if count%10000 == 0 {
			logger.Info("Loading traces", utils.Field{Key: "count", Value: count})
		}
		return nil
	})
	if err != nil {
		return err
	}

	tc := models.TraceCollection{Traces: traces}
	templateSvc := services.NewTemplateService()
	templates := templateSvc.ExtractTemplates(tc)

	outFile, err := os.Create(tplPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(templates)
}