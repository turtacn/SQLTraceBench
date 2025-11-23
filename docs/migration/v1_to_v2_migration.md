# Migration Guide: v1.x to v2.0

## Overview
This guide helps you migrate from SQL Trace Bench v1.x to v2.0.

## Breaking Changes Summary
1.  **Configuration file format changed**: The YAML structure for model and generation settings has been unified.
2.  **API endpoints renamed**: All endpoints are now versioned under `/api/v2/`.
3.  **CLI commands renamed**: Command structure has been simplified (e.g., `generate-trace` -> `generate`).
4.  **Database schema updated**: New tables for storing validation results and benchmark metrics.

## Step-by-Step Migration

### 1. Backup Data
```bash
# Backup v1 database
pg_dump -h localhost -U user trace_bench_v1 > backup_v1.sql

# Backup config files
cp -r configs/ configs_backup/
```

### 2. Update Configuration Files

**Old Format (v1.x)**:
```yaml
generator:
  model_type: "markov"
  parameters:
    order: 2
```

**New Format (v2.0)**:
```yaml
model:
  type: "markov"
  config:
    order: 2
```

**Automated Migration**:
We provide a utility to migrate your configs automatically:
```bash
sql_trace_bench migrate-config \
  --input configs/v1/generation.yaml \
  --output configs/v2/generation.yaml
```

### 3. Update CLI Commands

| v1.x Command      | v2.0 Command    |
| ----------------- | --------------- |
| `generate-trace`  | `generate`      |
| `verify-trace`    | `validate`      |
| `benchmark-model` | `benchmark run` |

### 4. Update API Clients

**Old Endpoint (v1.x)**: `POST /api/generate`
**New Endpoint (v2.0)**: `POST /api/v2/generation/traces`

### 5. Migrate Database Schema

```bash
sql_trace_bench db migrate \
  --from-version 1.5.0 \
  --to-version 2.0.0
```

### 6. Update Docker Compose

Ensure you are using the v2 image and have added the Prometheus/Grafana services if you want metrics.

```yaml
services:
  app:
    image: sql-trace-bench:v2.0.0
    environment:
      - ENABLE_VALIDATION=true
```

## Rollback Plan

If issues occur, rollback to v1.x:
1.  Restore database from `backup_v1.sql`.
2.  Restore configs from `configs_backup/`.
3.  Downgrade Docker image to `v1.5.0`.
