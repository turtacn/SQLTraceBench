package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/app/workflow"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/internal/utils/terminal"
	"github.com/turtacn/SQLTraceBench/pkg/config"
	"github.com/turtacn/SQLTraceBench/pkg/types"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"gopkg.in/yaml.v3"
)

var (
	Version   = types.Version
	cfgFile   string
	pluginDir string
	noColor   bool
	verbose   bool
	autoYes   bool // For workflow run
	cfg       *types.Config
	rootCmd   = &cobra.Command{
		Use:     "sqltracebench",
		Short:   "SQL trace-based workload benchmark CLI",
		Version: Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Handle flags
			if noColor {
				terminal.ColorEnabled = false
			}

			// Load the configuration.
			var err error
			cfg, err = config.Load(cfgFile)
			if err != nil {
				return err
			}

			// Initialize the logger.
			// If verbose is on, maybe force Debug level?
			logLevel := cfg.Log.Level
			if verbose {
				logLevel = "debug"
			}
			logger := utils.NewLogger(logLevel, cfg.Log.Format, nil)
			utils.SetGlobalLogger(logger)

			// Load plugins
			if err := loadPlugins(); err != nil {
				return err
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cleanupPlugins()
		},
	}
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Run multi-phase benchmark workflow",
}

var workflowRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute full pipeline: convert -> generate -> run -> validate",
	RunE: func(cmd *cobra.Command, args []string) error {
		pipelineConfigFile, _ := cmd.Flags().GetString("config")

		// Load Pipeline Config
		var pipelineCfg workflow.WorkflowConfig
		data, err := os.ReadFile(pipelineConfigFile)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, &pipelineCfg); err != nil {
			return err
		}

		// Confirmation Step
		if !autoYes && terminal.IsTerminal() {
			fmt.Println(terminal.Info("Workflow Plan:"))
			fmt.Printf("  Target Plugin:    %s\n", pipelineCfg.TargetPlugin)
			fmt.Printf("  Input Traces:     %s\n", pipelineCfg.InputTracePath)
			fmt.Printf("  Generation Count: %d\n", pipelineCfg.Generation.Count)
			fmt.Printf("  Concurrency:      %d\n", pipelineCfg.Execution.Concurrency)
			fmt.Printf("  Output Dir:       %s\n", pipelineCfg.OutputDir)

			fmt.Print("\nDo you want to proceed? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println(terminal.Warning("Workflow cancelled by user."))
				return nil
			}
		}

		// Initialize Services
		parser := parsers.NewAntlrParser()
		registry := plugin_registry.GlobalRegistry

		convSvc := conversion.NewService(parser, registry)
		genSvc := generation.NewService()
		execSvc := execution.NewService(registry)
		valSvc := validation.NewService()

		// Initialize Manager
		mgr := workflow.NewManager(convSvc, genSvc, execSvc, valSvc)

		// Run
		return mgr.Run(cmd.Context(), pipelineCfg)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", types.DefaultConfigPath, "config file")
	rootCmd.PersistentFlags().StringVar(&pluginDir, "plugin-dir", "./bin", "Directory where plugins are located")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	workflowRunCmd.Flags().StringP("config", "c", "", "Pipeline config YAML")
	workflowRunCmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "Skip confirmation prompt")
	workflowRunCmd.MarkFlagRequired("config")
	workflowCmd.AddCommand(workflowRunCmd)
	rootCmd.AddCommand(workflowCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd returns the root command for testing purposes.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// loadPlugins loads plugins from the configured directory.
func loadPlugins() error {
	return plugin_registry.LoadPlugins(pluginDir)
}

func cleanupPlugins() {
	plugin_registry.ClosePlugins()
}
