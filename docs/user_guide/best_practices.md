# Best Practices

## Data Preparation
*   **Sanitize Logs**: Before training, ensure your logs don't contain sensitive PII. The tool creates synthetic data, but the model files might retain patterns.
*   **Representative Sampling**: Ensure your input logs cover peak hours, off-peak hours, and special events (e.g., end-of-month processing).

## Model Selection
*   **Markov**: Best for simple workloads where query sequence is strictly probabilistic (e.g., web navigation). Fast training, low memory.
*   **LSTM**: Better for capturing long-term dependencies and complex user sessions. Slower training.
*   **Transformer**: Use for the highest fidelity when modeling complex SQL syntax and parameter correlations. Requires GPU for reasonable training time.

## Validation
*   **Don't skip validation**: Always validate a new model before using it for critical benchmarks.
*   **Check the tail**: Look at the P99 and P99.9 metrics in the validation report. Averages can hide poor tail performance modeling.

## Benchmarking
*   **Isolate the Environment**: Run the benchmark tool on a separate machine from the target database to avoid resource contention.
*   **Warmup**: Always configure a warmup period (e.g., 5 minutes) to allow database caches to populate.
*   **Clean State**: Reset the database to a known state before each run to ensure reproducibility.

## CI/CD Integration
*   Use the CLI in your CI pipeline to automatically validate that schema changes don't degrade performance.
*   Store benchmark results in a time-series database (like the built-in Prometheus exporter) to track performance trends over time.
