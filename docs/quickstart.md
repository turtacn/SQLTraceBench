# Quick Start Guide

This guide will walk you through the essential commands to get started with SQLTraceBench. We will cover the end-to-end workflow, from generating a workload to running a benchmark and validating the results.

## Prerequisites

- Go 1.18 or higher installed
- Make sure `GOPATH/bin` is in your PATH.

## 1. Build the Project

First, build the main application and the plugins.

```bash
make build
```

This will create the following binaries in the `./bin` directory:
- `sqltracebench`: The main CLI tool.
- `clickhouse`: The ClickHouse database plugin.

## 2. Prepare Data

You need a source file for conversion. This can be a SQL schema file (`.sql`) or a trace file (`.json` or `.jsonl`).

### Option A: SQL Schema

Create a file named `schema.sql` with a simple table definition:

```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(255),
    created_at DATETIME
);
```

### Option B: SQL Traces

Create a file named `traces.json` with a list of query objects:

```json
[
    {"timestamp": 1600000000, "duration": 100, "query": "SELECT * FROM users WHERE id = 1"},
    {"timestamp": 1600000001, "duration": 120, "query": "SELECT * FROM users WHERE id = 2"}
]
```

## 3. Convert

Use the `convert` command to translate your source file into the target database dialect (e.g., ClickHouse).

```bash
# Convert Schema
./bin/sqltracebench convert \
  --source schema.sql \
  --target clickhouse \
  --out ch_schema.sql \
  --mode schema

# Check result
cat ch_schema.sql
```

## 4. Generate Workload

Use the `generate` command to create a synthetic workload based on your source traces. This requires a trace file.

```bash
./bin/sqltracebench generate \
  --source-traces traces.json \
  --out workload.json \
  --count 100
```

## 5. Run Benchmark

Run the benchmark using the generated workload against the target database plugin.

```bash
./bin/sqltracebench run \
  --workload workload.json \
  --db clickhouse \
  --out metrics.json
```

**Note:** Ensure the target database is accessible if the plugin requires a connection. For this quickstart, the ClickHouse plugin might try to connect to localhost:9000 by default.

## 6. Validate (Optional)

If you have results from two different runs (e.g., `metrics_v1.json` and `metrics_v2.json`), you can compare them.

```bash
./bin/sqltracebench validate \
  --base metrics_v1.json \
  --candidate metrics_v2.json \
  --out report_dir
```

## Running the Automated Quickstart Script

We provide a script that automates all these steps for verification.

```bash
./examples/quickstart.sh
```

This script will compile the project, create dummy data, and run through the conversion, generation, and execution steps.
