package cmd

import (
	"github.com/spf13/cobra"
)

var (
	benchCmd = &cobra.Command{
		Use:   "bench",
		Short: "Run workload benchmark",
	}
	benchWorkloadFile string
	benchTargetType   string
)

func init() {
	benchCmd.Flags().StringVarP(&benchWorkloadFile, "workload", "w", "workload.yaml", "workload file")
	benchCmd.Flags().StringVarP(&benchTargetType, "type", "d", "mysql", "database driver: mysql|clickhouse|starrocks")
	rootCmd.AddCommand(benchCmd)
}
