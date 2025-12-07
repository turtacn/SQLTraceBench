# ClickHouse Type Mapping Guide

## Overview
This guide details the type mapping rules from MySQL/PostgreSQL/TiDB to ClickHouse used by the conversion engine.

## Intelligent Mapping Engine

### Context-Aware Mapping
The mapping engine optimizes type selection based on the following context:
- **Column Name Semantics**: Identifies special column names like `email`, `ip`, `status`.
- **Primary Key/Index**: `VARCHAR` in primary keys is prioritized to use `FixedString` for performance.
- **Nullability**: Affects `Nullable()` wrapping.
- **Default Value**: Used to infer Enum values (if applicable).

### Precision Preservation Policy
1. **DECIMAL Precision**:
   - `DECIMAL(P≤9, S)` → `Decimal32(S)`
   - `DECIMAL(P≤18, S)` → `Decimal64(S)`
   - `DECIMAL(P≤38, S)` → `Decimal128(S)`
   - `DECIMAL(P>38, S)` → `Decimal256(S)` (or `String` depending on policy)

2. **TIMESTAMP Precision**:
   - `TIMESTAMP` → `DateTime` (Second precision)
   - `TIMESTAMP(3)` → `DateTime64(3)` (Millisecond)
   - `TIMESTAMP(6)` → `DateTime64(6)` (Microsecond)

## Type Mapping Table

| Source Type (MySQL) | ClickHouse Type | Notes |
|---------------------|-----------------|-------|
| `INT` | `Int32` | |
| `BIGINT` | `Int64` | |
| `BIGINT UNSIGNED` | `UInt64` | Check for overflow risk if mapped to Int64 |
| `VARCHAR(N)` | `String` or `FixedString(N)` | `FixedString` used for PKs with N <= 256 |
| `TEXT` | `String` | |
| `DECIMAL(P,S)` | `Decimal*(S)` | Automatically selects best precision |
| `TIMESTAMP` | `DateTime` or `DateTime64` | Preserves original precision |
| `JSON` | `String` | Suggest using `JSONExtract` functions |

*(See `configs/type_mapping_rules.yaml` for full rules)*

## Custom Mapping Rules

### Editing Rule File
Edit `configs/type_mapping_rules.yaml` to add or modify rules:

```yaml
custom_rules:
  mysql:clickhouse:
    GEOMETRY: "String"  # Store spatial types as String

context_rules:
  - name: "custom_uuid_column"
    priority: 100
    conditions:
      - field: "column_name"
        operator: "equals"
        value: "uuid"
    target_type: "UUID"
```

### Hot Reload
Changes to the rule file are automatically detected and applied (latency < 100ms) without restarting the service.

## Warning System

### Warning Levels
* **INFO**: Informational messages (e.g., optimizations applied).
* **WARNING**: Potential data issues (e.g., precision loss, overflow risk).
* **ERROR**: Critical issues requiring manual review.

### Generating Reports
The conversion tool collects warnings which can be exported to JSON or Markdown.

## Best Practices
1. **Primary Key Design**: Use `FixedString` for fixed-length keys to optimize performance.
2. **Precision**: Use the smallest `Decimal` type that fits your data.
3. **String Optimization**: Use `LowCardinality(String)` for low-cardinality columns like status or category.
4. **Timezones**: Explicitly handle timezones if needed, e.g., `DateTime('UTC')`.
