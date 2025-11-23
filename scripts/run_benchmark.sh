#!/bin/bash
# scripts/run_benchmark.sh

set -e

echo "=========================================="
echo "  SQL Trace Bench - Benchmark Runner"
echo "=========================================="

# Config
CONFIG_FILE="${1:-configs/benchmark.yaml}"
OUTPUT_DIR="./benchmark_reports/$(date +%Y%m%d_%H%M%S)"
PROMETHEUS_PORT=9091

# 1. Create Output Directory
mkdir -p "$OUTPUT_DIR"

# 2. Build the tool
echo "Building sql_trace_bench..."
go build -o sql_trace_bench_bin cmd/sql_trace_bench/main.go

# 3. Execute Benchmark in Background
echo "Running benchmark with config: $CONFIG_FILE"

./sql_trace_bench_bin benchmark run \
    --config "$CONFIG_FILE" \
    --output "$OUTPUT_DIR" \
    --prometheus &
BENCH_PID=$!

# 4. Wait for Benchmark to finish its work and start serving metrics
# We need to detect when the benchmark is "done" but waiting for Ctrl+C.
# For simplicity, we sleep for a duration estimated or poll the report file.
echo "Waiting for benchmark to complete..."
sleep 10

# 5. Scrape Metrics
echo "Collecting Prometheus metrics..."
curl -s "http://localhost:9091/metrics" > "$OUTPUT_DIR/metrics.txt"

# 6. Kill the benchmark process
kill $BENCH_PID

echo "=========================================="
echo "Benchmark completed!"
echo "Report: $OUTPUT_DIR"
echo "Metrics: $OUTPUT_DIR/metrics.txt"
echo "=========================================="
