# Data Flow Architecture

This document describes the flow of data through the SQL Trace Bench system, covering the three main pipelines: Learning & Generation, Validation, and Benchmarking.

## 1. Learning & Generation Pipeline

This pipeline transforms raw SQL logs into a statistical model and then synthesizes new trace workloads.

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Parser
    participant Extractor as ParameterExtractor
    participant Trainer as ModelTrainer
    participant Generator as TraceGenerator
    participant Storage

    User->>CLI: generate -m markov -i raw.log
    CLI->>Parser: Parse SQL Logs
    Parser->>Extractor: AST / Tokens
    Extractor->>Extractor: Identify Templates & Params
    Extractor->>Trainer: Structured Trace Data
    Trainer->>Trainer: Train Model (Markov/LSTM)
    Trainer->>Storage: Save Model State

    User->>CLI: (Continue Generation)
    CLI->>Generator: Generate(count=1000)
    Generator->>Storage: Load Model
    loop Generation
        Generator->>Generator: Sample Next Query
        Generator->>Generator: Fill Parameters
    end
    Generator->>Storage: Write Generated Traces (JSON/CSV)
```

### Data Structures
*   **Raw Log**: Plain text SQL queries.
*   **Trace**: Structured object with `Timestamp`, `QueryTemplate`, `Parameters`, `Duration`.
*   **Model Artifact**: Serialized probability matrices (Markov) or weights (LSTM).

## 2. Validation Pipeline

This pipeline compares the statistical properties of the generated traces against the original workload to ensure fidelity.

```mermaid
flowchart LR
    Original[Original Traces] --> Validator
    Generated[Generated Traces] --> Validator

    subgraph Validator [Validation Service]
        Stats[Statistical Analyzer]
        Tests[Hypothesis Tests]

        Stats -->|Distributions| Tests
        Tests -->|K-S Test| Score1[Score]
        Tests -->|Chi-Square| Score2[Score]
    end

    Score1 --> Report[Validation Report]
    Score2 --> Report
    Report -->|HTML/JSON| Output
```

### Key Metrics
*   **Query Type Distribution**: Frequency of SELECT, INSERT, UPDATE, DELETE.
*   **Parameter Distribution**: Value distributions for extracted parameters.
*   **Temporal Patterns**: Arrival rate and inter-arrival times.
*   **Validation Score**: A composite score (0.0 - 1.0) indicating similarity.

## 3. Benchmarking Pipeline

This pipeline executes the generated traces against a target database to measure performance.

```mermaid
graph TD
    Input[Generated Traces] --> Runner[Benchmark Runner]
    Config[Benchmark Config] --> Runner

    subgraph Execution [Execution Engine]
        Pool[Worker Pool]
        Plugin[DB Plugin]
    end

    Runner --> Pool
    Pool -->|Execute| Plugin
    Plugin -->|SQL| TargetDB[(Target Database)]

    TargetDB -->|Result/Error| Plugin
    Plugin -->|Latency/Status| Pool

    Pool -->|Metrics| Reporter[Prometheus Reporter]
    Reporter -->|Expose| HTTP[HTTP :9091/metrics]
```

### Metric Flow
1.  **Latency Recording**: Each execution measures start and end time.
2.  **Aggregation**: Metrics are aggregated by Query Template ID.
3.  **Export**: Prometheus scrapes the metrics endpoint.
4.  **Visualization**: Grafana dashboards query Prometheus to display QPS, P99, Error Rate.
