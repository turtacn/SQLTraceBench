package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
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
	// This is a temporary solution to the import cycle.
	// In a real application, this logic would be handled in the app layer.
	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var traces []models.SQLTrace
	dec := json.NewDecoder(file)
	for dec.More() {
		var t models.SQLTrace
		if err := dec.Decode(&t); err != nil {
			return err
		}
		traces = append(traces, t)
	}

	if targetPlugin != "" {
		plugin, err := plugins.GetPlugin(targetPlugin)
		if err != nil {
			return err
		}
		for i := range traces {
			translatedQuery, err := plugin.TranslateQuery(traces[i].Query)
			if err != nil {
				continue
			}
			traces[i].Query = translatedQuery
		}
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