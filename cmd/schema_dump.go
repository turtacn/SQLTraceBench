package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
)

var (
	dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump a database schema to a file",
		RunE:  runDump,
	}
	dumpDsn    string
	dumpDriver string
	dumpOut    string
)

func init() {
	schemaCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().StringVar(&dumpDsn, "dsn", "", "Data Source Name for the database to dump")
	dumpCmd.Flags().StringVar(&dumpDriver, "driver", "mysql", "Database driver for the database to dump")
	dumpCmd.Flags().StringVarP(&dumpOut, "out", "o", "schema.json", "Path to the output schema file")
}

func runDump(cmd *cobra.Command, args []string) error {
	db, err := sql.Open(dumpDriver, dumpDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	extractor := database.NewMySqlSchemaExtractor(db)
	schema, err := extractor.ExtractSchema(context.Background())
	if err != nil {
		return err
	}

	file, err := os.Create(dumpOut)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(schema)
}