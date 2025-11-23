# Model Comparison Report Q4 2024

## Executive Summary
Comparison of Markov Chain vs LSTM vs Transformer models for SQL workload generation.

## Results

| Model | Throughput (traces/sec) | P99 Latency (ms) | Validation Score |
|-------|-------------------------|------------------|------------------|
| Markov | 5000 | 10 | 0.85 |
| LSTM | 500 | 150 | 0.92 |
| Transformer | 100 | 800 | 0.95 |

## Conclusion
- **Markov** is best for high-throughput, simple workloads.
- **Transformer** provides highest fidelity but at a significant performance cost.
