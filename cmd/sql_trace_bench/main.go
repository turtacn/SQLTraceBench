package main

import (
	"os"

	"github.com/turtacn/SQLTraceBench/cmd"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

func main() {
	utils.SetGlobalLogger(utils.NewLogger("info", "text", nil))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
