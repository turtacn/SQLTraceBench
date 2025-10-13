package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/cmd"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

func TestMainCommand(t *testing.T) {
	// Redirect stdout to a buffer to capture the output.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the root command with the --version flag.
	rootCmd := cmd.GetRootCmd()
	rootCmd.SetArgs([]string{"--version"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	// Restore stdout.
	w.Close()
	os.Stdout = old

	// Read the output from the buffer.
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Assert that the output contains the version string.
	assert.Contains(t, buf.String(), types.Version)
}