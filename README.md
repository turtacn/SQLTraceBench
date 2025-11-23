# SQL Trace Bench

[![Build Status](https://github.com/yourusername/sql-trace-bench/workflows/CI/badge.svg)](https://github.com/yourusername/sql-trace-bench/actions)
[![Coverage](https://codecov.io/gh/yourusername/sql-trace-bench/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/sql-trace-bench)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/sql-trace-bench)](https://goreportcard.com/report/github.com/yourusername/sql-trace-bench)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/github/v/release/yourusername/sql-trace-bench)](https://github.com/yourusername/sql-trace-bench/releases)

> **A high-performance SQL trace generation and validation framework powered by statistical analysis and machine learning.**

## 🚀 Features

-   **Multiple Generation Models**: Markov Chain, LSTM, Transformer.
-   **Statistical Validation**: 10+ validation tests (K-S, Chi-Square, Auto-correlation).
-   **Performance Benchmarking**: Compare models side-by-side.
-   **RESTful API**: Programmatic access with OpenAPI 3.0 spec.
-   **Real-time Monitoring**: Prometheus + Grafana integration.
-   **Rich Visualization**: Interactive HTML reports with charts.

## 📦 Quick Start

```bash
# Install
go install github.com/yourusername/sql-trace-bench@latest

# Generate 1000 traces
sql_trace_bench generate -m markov -n 1000

# Validate traces
sql_trace_bench validate \
  --original ./testdata/sample_traces.json \
  --generated ./output/generated.json

# View report
open ./output/validation_report.html
```

**[See Full Quickstart Guide →](docs/user_guide/quickstart.md)**

## 📚 Documentation

### User Guides
*   [Quickstart](docs/user_guide/quickstart.md)
*   [Advanced Usage](docs/user_guide/advanced_usage.md)
*   [Best Practices](docs/user_guide/best_practices.md)

### API Reference
*   [REST API](docs/api/rest_api.md)
*   [CLI Reference](docs/api/cli_reference.md)

### Architecture
*   [System Architecture](docs/architecture/system_architecture.md)
*   [Data Flow](docs/architecture/data_flow.md)

### Troubleshooting
*   [Common Issues](docs/troubleshooting/common_issues.md)

### Releases & Migration
*   [Release Notes v2.0.0](docs/releases/v2.0.0_release_notes.md)
*   [Migration Guide (v1 -> v2)](docs/migration/v1_to_v2_migration.md)
*   [Changelog](CHANGELOG.md)

## 🎯 Use Cases

*   **Performance Testing**: Generate realistic SQL workloads for database benchmarking.
*   **Load Testing**: Simulate production traffic patterns.
*   **Anomaly Detection**: Train models on normal traces, detect anomalies.
*   **Capacity Planning**: Predict system behavior under different loads.

## 📊 Performance

| Model       | Throughput | Memory  | P99 Latency |
| ----------- | ---------- | ------- | ----------- |
| Markov      | 1,245/sec  | 256 MB  | 45 ms       |
| LSTM        | 523/sec    | 512 MB  | 120 ms      |
| Transformer | 312/sec    | 1024 MB | 200 ms      |

## 🏗️ Architecture

```
┌─────────────┐
│     CLI     │
└──────┬──────┘
       │
┌──────▼──────┐
│ Application │
│    Layer    │
└──────┬──────┘
       │
┌──────▼──────┐     ┌─────────────┐
│   Domain    │────▶│Infrastructure│
│    Layer    │     │    Layer    │
└─────────────┘     └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Database   │
                    │  Prometheus │
                    │  Grafana    │
                    └─────────────┘
```

**[See Full Architecture →](docs/architecture/system_architecture.md)**

## 🛠️ Development

We welcome contributions! Please see our [Contributing Guide](docs/development/contributing.md) and [Development Guide](docs/development/development_guide.md).

## 📄 License
MIT License.
