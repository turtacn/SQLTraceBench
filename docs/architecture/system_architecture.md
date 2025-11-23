# SQL Trace Bench - System Architecture

## 1. Introduction
SQL Trace Bench is a high-performance, extensible tool designed for generating, validating, and benchmarking SQL trace workloads. It employs a Clean Architecture approach to ensure separation of concerns, testability, and maintainability. The core domain logic is isolated from external frameworks and drivers, allowing for easy adaptation to new database engines and storage systems.

### Key Design Principles
*   **Domain-Driven Design (DDD)**: The system is modeled around the core business domains: Trace Generation, Validation, and Benchmarking.
*   **Clean Architecture**: Dependencies point inwards. The core use cases and entities have no knowledge of the outer layers (web, database, CLI).
*   **SOLID**: Adherence to SOLID principles to ensure code quality and extensibility.
*   **Plugin Architecture**: Database dialects and specific execution logic can be extended via HashiCorp's `go-plugin` system (Out-of-Process).

## 2. C4 Model

### 2.1 Context Diagram (Level 1)
This diagram shows the system in context with its users and external dependencies.

```mermaid
graph TD
    User[User/Tester]
    System[SQL Trace Bench]
    TargetDB[Target Database<br/>(PostgreSQL, MySQL, etc.)]
    Prometheus[Prometheus]
    Grafana[Grafana]

    User -->|Uses CLI/API| System
    System -->|Executes Queries| TargetDB
    System -->|Exports Metrics| Prometheus
    Prometheus -->|Scrapes Metrics| Grafana
    Grafana -->|Visualizes Data| User
```

### 2.2 Container Diagram (Level 2)
This diagram illustrates the high-level containers that make up the system.

```mermaid
graph TD
    CLI[CLI Application]
    API[REST API Server]

    subgraph "SQL Trace Bench System"
        AppLayer[Application Layer<br/>(Service Orchestration)]
        DomainLayer[Domain Layer<br/>(Core Logic & Models)]
        InfraLayer[Infrastructure Layer<br/>(Implementations)]
    end

    CLI -->|Calls| AppLayer
    API -->|Calls| AppLayer

    AppLayer -->|Uses| DomainLayer
    DomainLayer -->|Defines Interfaces| InfraLayer
    InfraLayer -->|Implements Interfaces| DomainLayer

    InfraLayer -->|Connects| DB[(Target Database)]
    InfraLayer -->|Exposes| Metrics(Prometheus Metrics)
```

### 2.3 Component Diagram (Level 3)
A deeper look into the components within the layers.

#### Domain Layer
*   **TraceGenerator**: Interface for trace generation strategies.
    *   `MarkovGenerator`: Generates traces based on Markov Chains.
    *   `LSTMGenerator`: Uses LSTM neural networks for sequence generation.
    *   `TransformerGenerator`: Uses Transformer models for advanced pattern learning.
*   **StatisticalValidator**: Service for validating generated traces against originals using statistical tests (KS, Chi-Square).
*   **BenchmarkRunner**: orchestrates performance benchmarks comparing different models.
*   **Models**: `Trace`, `Schema`, `ValidationResult`, `BenchmarkConfig`.

#### Application Layer
*   **GenerationService**: Coordinates the training and generation process.
*   **ValidationService**: Handles trace validation requests and report generation.
*   **BenchmarkService**: Manages benchmark execution and result aggregation.

#### Infrastructure Layer
*   **PostgresRepository**: Implementation of storage interfaces for PostgreSQL.
*   **S3Storage**: (Optional) Object storage implementation.
*   **PrometheusMetrics**: Implementation of the metrics interface.
*   **PluginSystem**: Manages loading and communication with external database plugins.

## 3. Data Flow

1.  **Workload Learning**:
    *   Raw SQL logs -> `Parser` -> `Trace Model` -> `ParameterExtractor` -> `Model Training` (Markov/LSTM).
2.  **Trace Generation**:
    *   User Request -> `GenerationService` -> `TraceGenerator` -> `Synthesized Traces` -> `Output Writer`.
3.  **Validation**:
    *   Original & Generated Traces -> `ValidationService` -> `StatisticalValidator` -> `Validation Report` (HTML/JSON).
4.  **Benchmarking**:
    *   `BenchmarkService` -> `Worker Pool` -> `DBExecutionService` -> `Target Database`.
    *   Metrics -> `PrometheusReporter`.

## 4. Technology Stack
*   **Language**: Go 1.21+
*   **Database**: PostgreSQL 15+ (for metadata and repository), plus target databases for benchmarking.
*   **Metrics**: Prometheus (instrumentation), Grafana (visualization).
*   **Parsing**: ANTLR v4.
*   **ML/Stats**: `gonum` for statistics, custom implementations for lightweight ML models.
*   **CLI**: Cobra.
*   **Containerization**: Docker, Docker Compose.

## 5. Design Patterns
*   **Repository Pattern**: Decouples domain logic from data access.
*   **Factory Pattern**: Used for creating specific Generator and Validator instances based on configuration.
*   **Strategy Pattern**: Allows swapping generation algorithms (Markov vs LSTM) and validation tests at runtime.
*   **Observer Pattern**: Used in the validation module to trigger events/alerts upon validation failures.
*   **Worker Pool**: Used in `BenchmarkRunner` to manage concurrent execution load.
