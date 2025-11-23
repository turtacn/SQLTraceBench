# Validation Reports Guide

## HTML Report

The HTML report provides a comprehensive view of the validation results.

### Sections

1.  **Summary**: Overall score and pass/fail status.
2.  **Distribution Tests**: Bar chart and table showing KS test results for each parameter.
    - **P-Value**: Indicates the probability that the generated distribution matches the original. Higher is better (closer to 1.0).
    - **Status**: PASS if P-Value > Threshold (0.05).
3.  **Temporal Tests**: (Planned) Heatmap of query frequency over time.
4.  **Query Type Tests**: (Planned) Comparison of query template mix.

## Prometheus Metrics

Enable Prometheus export by setting `prometheus_port` in `validation.yaml` or CLI args.

| Metric Name | Type | Description |
|---|---|---|
| `validation_pass_rate` | Gauge | Overall pass rate (0.0 to 1.0). |
| `validation_ks_pvalue` | GaugeVec | KS test p-value per parameter. |

### Example Query
```promql
validation_pass_rate < 0.8
```
Alert if pass rate drops below 80%.
