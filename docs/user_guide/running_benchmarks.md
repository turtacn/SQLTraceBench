# Running Benchmarks

## Prerequisites
- Go 1.20+
- Prometheus & Grafana (optional, for visualization)

## Quick Start

1. **Configure Benchmark**
   Edit `configs/benchmark.yaml` to define your models and test scenarios.

2. **Run Benchmark**
   ```bash
   ./scripts/run_benchmark.sh
   ```

3. **View Results**
   Open the generated HTML report in `benchmark_reports/`.

## CLI Usage

```bash
sql_trace_bench benchmark run --config <path> --output <dir> --prometheus
```

## Interpreting Results
- **Throughput**: Higher is better.
- **P99 Latency**: Lower is better. Indicates the 99th percentile of generation time.
