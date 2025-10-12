package cmd

import "github.com/spf13/cobra"

var (
	schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Manage database schemas",
	}
)

func init() {
	rootCmd.AddCommand(schemaCmd)
}