# Precision Handling Policy

This document describes the strategies used to handle precision and scale during database schema conversion.

## DECIMAL Handling
Configuration in `configs/precision_policy.yaml`:

```yaml
decimal_policy:
  overflow_strategy: "WARN"  # Options: TRUNCATE, ERROR, WARN
  max_precision: 38          # Max precision for Decimal128 (default)
  prefer_decimal256: false   # Set true to use Decimal256 for P > 38
```

### Strategies
- **Exact Match**: If target supports exact precision (e.g. `Decimal(10,2)`), it is used.
- **Promotion**: If target supports higher precision bucket (e.g. `Decimal32` -> `Decimal64`), it is promoted.
- **Overflow**: If source precision > target max precision:
    - `WARN`: Generates a warning and uses max available or String.
    - `TRUNCATE`: Truncates precision (dangerous).
    - `ERROR`: Aborts conversion.

## TIMESTAMP Handling

```yaml
timestamp_policy:
  default_fractional_seconds: 3
  preserve_timezone: true
```

- **MySQL TIMESTAMP(N)**: Mapped to `DateTime64(N)`.
- **Postgres TIMESTAMPTZ**: Mapped to `DateTime64(N, 'UTC')` if supported, or `DateTime64` with warning.

## Floating Point

- **FLOAT/DOUBLE**:
    - `DOUBLE` -> `Float64` (Safe)
    - `DOUBLE` -> `Float32` (Unsafe, generates warning)

## String Length

- **VARCHAR(N)**:
    - If N <= threshold (e.g. 256) AND column is PK -> `FixedString(N)`
    - Otherwise -> `String`
