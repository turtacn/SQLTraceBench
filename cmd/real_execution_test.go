package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealExecution(t *testing.T) {
	// This test requires Docker and docker-compose to be installed.
	// It will be skipped if Docker is not available.
	if os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_MACHINE_NAME") == "" {
		t.Skip("Skipping test; Docker not available")
	}

	// Setup: Start the containers
	// In a real CI environment, this would be handled by a script.
	// For this test, we'll assume the containers are already running.
	// You can start them with `docker-compose up -d`.

	// Setup: Connect to the MySQL database
	dsn := "root:root@tcp(127.0.0.1:3306)/test?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Wait for the database to be ready.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, db.PingContext(ctx))

	// Setup: Create a table and insert some data
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY,
			name VARCHAR(255)
		)
	`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob') ON DUPLICATE KEY UPDATE name=VALUES(name)")
	require.NoError(t, err)

	// Setup: Create a temporary directory for the test files
	tmpDir, err := ioutil.TempDir("", "real_exec_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Setup: Create a dummy workload file
	workloadPath := filepath.Join(tmpDir, "workload.json")
	workloadContent := fmt.Sprintf(`
		{
			"queries": [
				{"query": "SELECT * FROM users WHERE id = ?", "args": [1]},
				{"query": "SELECT * FROM users WHERE id = ?", "args": [2]}
			]
		}
	`)
	ioutil.WriteFile(workloadPath, []byte(workloadContent), 0644)

	// Test the run command with the real executor
	metricsPath := filepath.Join(tmpDir, "metrics.json")
	rootCmd.SetArgs([]string{
		"run",
		"--executor", "real",
		"--driver", "mysql",
		"--dsn", dsn,
		"--workload", workloadPath,
		"--out", metricsPath,
	})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, metricsPath)

	// Verify the metrics file
	// In a real test, we would unmarshal the JSON and check the values.
	// For this smoke test, we'll just check that the file is not empty.
	metricsData, err := ioutil.ReadFile(metricsPath)
	require.NoError(t, err)
	assert.NotEmpty(t, metricsData)
}