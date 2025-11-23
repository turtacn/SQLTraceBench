# CLI Reference

## Installation
```bash
go install github.com/yourusername/sql-trace-bench@latest
```

## Global Flags
*   `--verbose, -v`: Enable verbose logging.
*   `--log-level`: Set log level (debug, info, warn, error) [default: info].
*   `--config, -c`: Path to configuration file.
*   `--help, -h`: Show help message.

## Commands

### `generate`
Generate SQL traces using a specified model.

**Usage**:
```bash
sql_trace_bench generate [flags]
```

**Flags**:
*   `--model, -m`: Model type (`markov`, `lstm`, `transformer`) [required].
*   `--count, -n`: Number of traces to generate [default: 1000].
*   `--output, -o`: Output directory [default: `./output`].
*   `--format`: Output format (`json`, `csv`, `parquet`) [default: `json`].
*   `--input, -i`: Input raw log file for training (if model not pre-trained).

**Examples**:
```bash
# Generate 1000 traces using Markov model
sql_trace_bench generate -m markov -n 1000

# Train on log file and generate
sql_trace_bench generate -m lstm -i ./logs/query.log -n 5000
```

### `validate`
Validate generated traces against original data.

**Usage**:
```bash
sql_trace_bench validate [flags]
```

**Flags**:
*   `--original`: Path to original traces [required].
*   `--generated`: Path to generated traces [required].
*   `--tests`: Comma-separated list of tests to run (e.g., `ks_test,chi_square`) [default: `all`].
*   `--report`: Generate HTML report [default: `true`].

**Examples**:
```bash
sql_trace_bench validate \
  --original ./data/original.json \
  --generated ./output/generated.json
```

### `benchmark run`
Run performance benchmarks.

**Usage**:
```bash
sql_trace_bench benchmark run [flags]
```

**Flags**:
*   `--config`: Benchmark config file [required].
*   `--output`: Output directory for reports.
*   `--prometheus`: Enable Prometheus metrics export [default: `true`].
*   `--port`: Port for Prometheus metrics [default: `9091`].

**Examples**:
```bash
sql_trace_bench benchmark run --config configs/benchmark.yaml
```

### `db migrate`
Manage database schema migrations.

**Usage**:
```bash
sql_trace_bench db migrate [flags]
```

**Flags**:
*   `--up`: Apply all up migrations.
*   `--down`: Rollback the last migration.
*   `--version`: Migrate to a specific version.

### `version`
Print the version number of SQL Trace Bench.

**Usage**:
```bash
sql_trace_bench version
```
