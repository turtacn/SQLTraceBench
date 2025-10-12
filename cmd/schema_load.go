package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
)

var (
	loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Load a database schema from a file",
		RunE:  runLoad,
	}
	loadSchemaPath string
	loadDsn        string
	loadDriver     string
	loadTarget     string
)

func init() {
	schemaCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringVarP(&loadSchemaPath, "schema", "s", "schema.json", "Path to the schema file")
	loadCmd.Flags().StringVar(&loadDsn, "dsn", "", "Data Source Name for the database to load to")
	loadCmd.Flags().StringVar(&loadDriver, "driver", "mysql", "Database driver for the database to load to")
	loadCmd.Flags().StringVar(&loadTarget, "target", "mysql", "Target dialect for schema translation")
}

func runLoad(cmd *cobra.Command, args []string) error {
	// Read the schema file.
	file, err := os.Open(loadSchemaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var schema models.DatabaseSchema
	if err := json.NewDecoder(file).Decode(&schema); err != nil {
		return err
	}

	// Translate the schema.
	schemaSvc := services.NewSchemaService()
	schemaSvc.RegisterTypeMapper("clickhouse", services.MySQLToClickHouseTypeMapper)
	translatedSchema, err := schemaSvc.ConvertTo(&schema, loadTarget)
	if err != nil {
		return err
	}

	// Load the schema.
	db, err := sql.Open(loadDriver, loadDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	var loader services.SchemaLoader
	switch loadTarget {
	case "mysql":
		loader = database.NewMySqlSchemaLoader(db)
	case "clickhouse":
		loader = database.NewClickHouseSchemaLoader(db)
	default:
		return fmt.Errorf("unsupported schema load target: %s", loadTarget)
	}

	return loader.LoadSchema(context.Background(), translatedSchema)
}