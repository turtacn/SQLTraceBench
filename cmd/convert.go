package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"github.com/turtacn/SQLTraceBench/plugins"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert SQL schema or trace files",
		RunE:  runConvert,
	}
	tracePath    string // Deprecated alias
	sourcePath   string
	outputPath   string
	targetPlugin string
	convertMode  string
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&tracePath, "trace", "t", "", "Path to the raw SQL trace file (deprecated, use --source)")
	convertCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "Path to the source file (schema SQL or trace JSON)")
	convertCmd.Flags().StringVarP(&outputPath, "out", "o", "output.json", "Path to the output file")
	convertCmd.Flags().StringVar(&targetPlugin, "target", "", "Target database for SQL translation (e.g., clickhouse, starrocks)")
	convertCmd.Flags().StringVarP(&convertMode, "mode", "m", "auto", "Conversion mode: 'schema' or 'trace' (auto-detect by extension)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	logger := utils.GetGlobalLogger()

	// Handle deprecated flag
	if sourcePath == "" && tracePath != "" {
		sourcePath = tracePath
		logger.Warn("Using deprecated flag --trace, please use --source")
	}
	if sourcePath == "" {
		return fmt.Errorf("source file path is required")
	}

	// Detect mode
	if convertMode == "auto" {
		ext := strings.ToLower(filepath.Ext(sourcePath))
		if ext == ".sql" || ext == ".ddl" {
			convertMode = "schema"
		} else if ext == ".json" || ext == ".jsonl" {
			convertMode = "trace"
		} else {
			return fmt.Errorf("could not auto-detect mode from file extension %s, please specify --mode", ext)
		}
	}

	// Instantiate Service
	// We use the GlobalRegistry which was initialized in root.go
	parser := parsers.NewAntlrParser()
	svc := conversion.NewService(parser, plugin_registry.GlobalRegistry)

	if convertMode == "schema" {
		logger.Info("Starting Schema Conversion", utils.Field{Key: "source", Value: sourcePath}, utils.Field{Key: "target", Value: targetPlugin})
		req := conversion.ConvertRequest{
			SourceSchemaPath: sourcePath,
			TargetDBType:     targetPlugin,
			OutputPath:       outputPath,
		}
		if err := svc.ConvertSchemaFromFile(cmd.Context(), req); err != nil {
			return err
		}
		logger.Info("Schema Conversion Complete", utils.Field{Key: "output", Value: outputPath})
		return nil
	}

	// Trace Conversion Logic
	logger.Info("Starting Trace Conversion", utils.Field{Key: "source", Value: sourcePath})
	bufferSize := viper.GetInt("parser.buffer_size")
	if bufferSize == 0 {
		bufferSize = 1024 * 1024
	}

	isStreamOutput := len(outputPath) > 5 && outputPath[len(outputPath)-6:] == ".jsonl"

	if isStreamOutput {
		return runTraceStreamConversion(cmd, svc, logger, bufferSize)
	}

	return runTraceBatchConversion(cmd, svc, logger, bufferSize)
}

func runTraceStreamConversion(cmd *cobra.Command, svc conversion.Service, logger *utils.Logger, bufferSize int) error {
	var plugin plugins.Plugin
	if targetPlugin != "" {
		p, err := plugin_registry.GetPlugin(targetPlugin)
		if err != nil {
			return err
		}
		plugin = p
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	encoder := json.NewEncoder(outFile)

	count := 0
	err = svc.ConvertStreamingly(cmd.Context(), sourcePath, bufferSize, func(trace models.SQLTrace) error {
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
	logger.Info("Trace Conversion complete", utils.Field{Key: "total_processed", Value: count})
	return nil
}

func runTraceBatchConversion(cmd *cobra.Command, svc conversion.Service, logger *utils.Logger, bufferSize int) error {
	if targetPlugin != "" {
		// Manual pipeline with translation
		file, err := os.Open(sourcePath)
		if err != nil {
			return err
		}
		// defer file.Close() // Will close manually or in block if needed, but best to defer.

		// Use a closure to handle file closing safely if we were doing complex logic,
		// but here we just read it.
		// Wait, I should defer Close() immediately.
		defer file.Close()

		parser := parsers.NewStreamingTraceParser(bufferSize)
		var traces []models.SQLTrace

		p, err := plugin_registry.GetPlugin(targetPlugin)
		if err != nil {
			return err
		}

		count := 0
		err = parser.Parse(file, func(trace models.SQLTrace) error {
			translatedQuery, err := p.TranslateQuery(trace.Query)
			if err == nil {
				trace.Query = translatedQuery
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
		templateSvc := services.NewTemplateService() // Direct usage of domain service
		tpls := templateSvc.ExtractTemplates(tc)

		outFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		encoder := json.NewEncoder(outFile)
		encoder.SetIndent("", "  ")
		return encoder.Encode(tpls)
	}

	// No translation needed -> Use Service directly.
	tpls, err := svc.ConvertFromFile(cmd.Context(), sourcePath)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tpls)
}
