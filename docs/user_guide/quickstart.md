# Quickstart Guide

## Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM

## Step 1: Clone Repository
```bash
git clone https://github.com/yourusername/sql-trace-bench.git
cd sql-trace-bench
```

## Step 2: Start Services

```bash
docker-compose up -d
```

This will start:
*   PostgreSQL (port 5432)
*   Prometheus (port 9090)
*   Grafana (port 3000)
*   SQL Trace Bench API (port 8080)

## Step 3: Generate Traces

You can use the CLI tool directly if installed, or run it via Docker.

**Option A: Using Docker (Recommended)**
```bash
docker-compose exec app sql_trace_bench generate \
  -m markov \
  -n 1000 \
  -o /output
```

**Option B: Using Local Binary**
First, install the binary:
```bash
go install .
```
Then run:
```bash
sql_trace_bench generate -m markov -n 1000
```

## Step 4: Validate Traces

Compare the generated traces against a sample dataset to ensure quality.

```bash
sql_trace_bench validate \
  --original ./testdata/sample_traces.json \
  --generated ./output/generated_traces.json
```

## Step 5: View Results

*   **Validation Report**: Open `./output/validation_report.html` in your browser.
*   **Metrics**: Visit `http://localhost:3000` (Grafana) to see real-time metrics (default admin/admin).

## Next Steps

*   [Advanced Usage](./advanced_usage.md)
*   [API Reference](../api/rest_api.md)
*   [Best Practices](./best_practices.md)

## Troubleshooting

If you encounter issues, check [Common Issues](../troubleshooting/common_issues.md).
