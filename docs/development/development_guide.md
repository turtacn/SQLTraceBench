# Development Guide

## Setup

### Prerequisites
*   Go 1.21+
*   Docker & Docker Compose
*   Make

### Installation
1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/sql-trace-bench.git
    cd sql-trace-bench
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Build the project:
    ```bash
    make build
    ```

## Running Tests
*   Run all tests:
    ```bash
    go test ./...
    ```
*   Run with coverage:
    ```bash
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
    ```

## Plugin Development

SQL Trace Bench uses HashiCorp's `go-plugin` system. To add a new database dialect:

1.  Create a new directory in `plugins/`.
2.  Implement the `DatabasePlugin` interface (defined in `pkg/proto`).
3.  Compile it as a standalone binary.
4.  Place the binary in `bin/` (or your configured plugin dir).

Example:
```go
// plugins/mysql/main.go
package main

import "github.com/hashicorp/go-plugin"

type MySQLPlugin struct{}

func (p *MySQLPlugin) Execute(query string) error {
    // Implementation
    return nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: shared.Handshake,
        Plugins: map[string]plugin.Plugin{
            "database": &MySQLPlugin{},
        },
    })
}
```

## Debugging
*   Use `dlv` for debugging.
*   Set `LOG_LEVEL=debug` environment variable.
*   Check `logs/app.log` for detailed runtime logs.
