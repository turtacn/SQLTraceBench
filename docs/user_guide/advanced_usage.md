# Advanced Usage

## Customizing Models

### LSTM Configuration
You can fine-tune the LSTM model by creating a custom configuration file.

`configs/lstm_config.yaml`:
```yaml
model:
  type: lstm
  epochs: 50
  batch_size: 64
  hidden_layers: [128, 64]
  dropout: 0.2
```

Run with config:
```bash
sql_trace_bench generate -m lstm -c configs/lstm_config.yaml
```

## Workload Hotspots
To simulate realistic workloads, you often need to model "hotspots" (frequently accessed data).

```yaml
generation:
  hotspots:
    - field: "user_id"
      distribution: "zipf"
      s: 1.1
    - field: "status"
      values: ["active", "pending"]
      weights: [0.9, 0.1]
```

## Large-Scale Generation
For generating millions of traces, use stream mode to avoid memory issues.

```bash
# Output to a compressed stream
sql_trace_bench generate -n 10000000 --format json.gz
```

## Benchmark Scenarios

Define complex scenarios in `configs/benchmark.yaml`.

```yaml
scenarios:
  - name: "heavy_read"
    concurrency: 50
    duration: "10m"
    workload: "read_heavy_traces.json"

  - name: "burst_write"
    concurrency: 100
    duration: "5m"
    workload: "write_burst_traces.json"
```

## Plugin Development
See [Development Guide](../development/development_guide.md) to learn how to write custom database plugins.
