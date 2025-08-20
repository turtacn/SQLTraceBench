package cmd

import (
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate workload YAML from templates",
	}
	genTemplateFile string
	genOut          string
)

func init() {
	generateCmd.Flags().StringVarP(&genTemplateFile, "templates", "f", "templates.yaml", "templates file")
	generateCmd.Flags().StringVarP(&genOut, "out", "o", "workload.yaml", "generated workload")
	rootCmd.AddCommand(generateCmd)
}
