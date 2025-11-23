# Validation Framework Design

## Overview
The Validation Framework ensures that the generated SQL workloads statistically resemble the original workloads. It uses statistical tests to compare distributions, temporal patterns, and query types.

## Architecture

### Components

1.  **StatisticalValidator**: Core service that performs statistical tests (KS Test, Chi-Square, Jensen-Shannon).
2.  **ReportGenerator**: Aggregates test results into a structured `ValidationReport`.
3.  **HTMLReporter**: Renders the report into a human-readable HTML format with charts.
4.  **PrometheusReporter**: Exports validation metrics (pass rate, p-values) to Prometheus.
5.  **ValidationService**: Orchestrates the validation process (Load -> Validate -> Report -> Export).

## Statistical Tests

### Kolmogorov-Smirnov (KS) Test
Used for continuous distributions (e.g., parameter values).
- **Null Hypothesis**: The two samples are drawn from the same distribution.
- **Pass Condition**: P-Value > Threshold (default 0.05).

### Chi-Square Test
Used for categorical data (e.g., query templates).
- **Pass Condition**: P-Value > Threshold (default 0.01).

## Metrics
- `validation_ks_pvalue{parameter="..."}`: P-value for specific parameter distribution.
- `validation_pass_rate`: Percentage of passed tests.
