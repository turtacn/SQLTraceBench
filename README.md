<div align="center">
  <img src="logo.png" alt="SQLTraceBench Logo" width="200" height="200">
  <h1>SQLTraceBench</h1>
  <p>
    <strong>A Universal Trace-driven Benchmark Framework for Cross-Database Performance Analysis.</strong>
  </p>
  <p>
    <!--<a href="https://github.com/turtacn/SQLTraceBench/actions/workflows/go.yml"><img src="https://github.com/turtacn/SQLTraceBench/actions/workflows/go.yml/badge.svg" alt="Build Status"></a>-->
    <a href="https://github.com/turtacn/SQLTraceBench/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
    <a href="https://goreportcard.com/report/github.com/turtacn/SQLTraceBench"><img src="https://goreportcard.com/badge/github.com/turtacn/SQLTraceBench" alt="Go Report Card"></a>
    <a href="https://github.com/turtacn/SQLTraceBench/releases"><img src="https://img.shields.io/github/v/release/turtacn/SQLTraceBench" alt="Latest Release"></a>
  </p>
  <p>
    <a href="README-zh.md"><strong>简体中文</strong></a> | <a href="docs/architecture.md"><strong>架构设计</strong></a>
  </p>
</div>

Welcome to SQLTraceBench! We are on a mission to revolutionize how developers and DBAs benchmark and compare database systems. Our framework empowers you to replay real-world production workloads across different database technologies, ensuring your choice of database is backed by data, not just claims.

## Core Mission

SQLTraceBench is a trace-driven benchmark system designed to take an existing SQL trace and database schema, generate an equivalent workload for various target databases, and execute a comprehensive performance comparison. It provides flexible controls for load, concurrency, and data generation, making it an indispensable tool for database evaluation, migration, and performance tuning.

## Why SQLTraceBench?

Database benchmarking is notoriously difficult. Synthetic benchmarks like TPC-H/DS are valuable but often don't reflect the unique query patterns and data skew of your specific application. SQLTraceBench bridges this gap by:

* **Real-World Workloads**: Uses your actual production SQL traces, providing a highly accurate representation of your application's behavior.
* **Cross-Database Translation**: Intelligently translates SQL dialects and database schemas (e.g., from StarRocks to ClickHouse, or PostgreSQL to TiDB), enabling true "apples-to-apples" comparisons.
* **Controllable Replay**: Goes beyond simple replay. It templates queries, models data distributions, and allows you to adjust QPS, concurrency, and hotspot ratios to simulate various scenarios like peak traffic or future growth.
* **Extensible by Design**: A powerful plugin system allows the community to easily add support for new databases.

## Key Features

* **Schema & SQL Translation**: Automatically converts DB schemas and SQL queries between different database systems.
* **SQL Templating & Parameterization**: Utilizes AST parsing to create query templates and models the distribution of query parameters from the original trace.
* **Synthetic Data Generation**: Analyzes data distributions from a source database and generates realistic, scalable datasets for target systems.
* **Flexible Workload Generation**: Fine-grained control over concurrency, QPS, query mix, and hotspot data access patterns.
* **Performance Validation**: Compares the performance profile (QPS, latency, rows scanned) of the benchmark against the original trace to ensure fidelity and provides detailed deviation reports.
* **Command-Line Interface**: A powerful and easy-to-use CLI tool, `sql_trace_bench`, to orchestrate the entire workflow.

## Getting Started

### Installation

Ensure you have Go (version 1.21+) installed.

```bash
go install [github.com/turtacn/SQLTraceBench/cmd/sql_trace_bench@latest](https://github.com/turtacn/SQLTraceBench/cmd/sql_trace_bench@latest)
````

### Basic Usage

Here's a quick example of how to run a benchmark. Imagine you have a SQL trace from StarRocks and want to test an equivalent workload on ClickHouse.

1.  **Prepare your configuration file (`config.yaml`):**

    ```yaml
    # config.yaml
    source:
      db_type: "starrocks"
      schema_file: "./schemas/starrocks_schema.yml"
      trace_file: "./traces/starrocks_trace.jsonl"
      data_source: # Connection details to analyze data distribution
        host: "starrocks-host"
        port: 9030
        user: "root"

    target:
      db_type: "clickhouse"
      # Connection details for the target DB to run the benchmark
      host: "clickhouse-host"
      port: 9000
      user: "default"

    workload:
      concurrency: 64 # Number of parallel workers
      duration: "5m"  # How long to run the benchmark
      qps_scale: 1.5  # Replay at 1.5x the original average QPS
    ```

2.  **Run the benchmark:**

    ```bash
    sql_trace_bench run --config ./config.yaml
    ```

    This command will:
    a. Parse the StarRocks schema and trace file.
    b. Convert the schema to be ClickHouse-compatible.
    c. Translate the SQL queries into ClickHouse dialect.
    d. Generate a workload based on the specified parameters.
    e. Execute the benchmark against your ClickHouse instance and output a report.

*(This is a conceptual example. The exact flags and config options are detailed in our documentation.)*

## Contributing

We are building a community of developers passionate about database performance and reliability. Contributions are highly welcome! Whether it's adding a new database plugin, improving the SQL translation logic, or enhancing documentation, your help is appreciated.

Please read our [Contributing Guidelines](https://www.google.com/search?q=./CONTRIBUTING.md) to get started.

## Community

Join our community channels to ask questions, share your ideas, and connect with other users.

  * **GitHub Discussions**: For questions and discussions.
  * **Slack/Discord**: (Link to be added)

## License

SQLTraceBench is licensed under the [Apache 2.0 License](https://www.google.com/search?q=./LICENSE).