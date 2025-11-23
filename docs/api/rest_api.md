# REST API Reference

## Base URL
`http://localhost:8080/api/v2`

## Authentication
Currently, the v2.0.0 API does not require authentication for local deployments. For production deployments behind an ingress, standard Basic Auth or OAuth2 should be configured at the gateway level.

## Endpoints

### 1. Generation

#### Generate Traces
`POST /generation/traces`

Initiates a new trace generation job.

**Request Body**:
```json
{
  "model": "markov",
  "count": 1000,
  "config": {
    "order": 2,
    "seed": 42
  }
}
```

**Response**:
```json
{
  "job_id": "gen-20240101-123456",
  "status": "running",
  "trace_count": 1000
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v2/generation/traces \
  -H "Content-Type: application/json" \
  -d '{"model":"markov","count":1000}'
```

### 2. Validation

#### Validate Traces
`POST /validation/traces/{job_id}`

Triggers validation of a generated trace set against its training data.

**Parameters**:
*   `job_id` (path): The ID of the generation job.

**Response**:
```json
{
  "validation_score": 0.87,
  "tests_passed": 8,
  "report_url": "http://localhost:8080/reports/val-20240101.html"
}
```

### 3. Benchmarking

#### Run Benchmark
`POST /benchmark`

Starts a performance benchmark using specified models or trace files.

**Request Body**:
```json
{
  "models": ["markov", "lstm"],
  "trace_count": 5000,
  "concurrency": 10
}
```

**Response**:
```json
{
  "benchmark_id": "bench-20240101-999",
  "status": "started"
}
```

### 4. System

#### Health Check
`GET /health`

Returns the health status of the service.

**Response**:
```json
{
  "status": "ok",
  "version": "v2.0.0",
  "uptime": "2h 15m"
}
```

---

## OpenAPI Specification

For the full OpenAPI 3.0 specification, please refer to the [`rest_api.yaml`](./rest_api.yaml) file. Developers can import this spec into tools like Postman or Swagger UI.
