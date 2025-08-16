# SQLTraceBench Architecture

## 1. Overview

SQLTraceBench is architected as a modular, extensible, and high-performance system for conducting trace-driven database benchmarks. Its design follows a layered and plugin-based approach, ensuring that core logic is decoupled from database-specific implementations. The system is orchestrated by a central CLI, which guides the user through a sequential pipeline from input analysis to final reporting.

The core principles guiding this architecture are:
* **Modularity and Decoupling**: Each stage of the pipeline (parsing, translation, execution) is an independent component. Database-specific logic is encapsulated within plugins, allowing for easy extension without modifying the core framework.
* **Extensibility**: The plugin interface is the cornerstone of the system, enabling first-party and third-party developers to add support for new databases with minimal effort.
* **Testability**: Clear interfaces between components facilitate unit and integration testing, ensuring the reliability of the transformation logic and workload execution.
* **Observability**: A centralized logging and metrics-gathering mechanism provides deep insights into the benchmark process, aiding in debugging and performance analysis.

### 1.1 DFX (Design for X) Analysis

| DFX Aspect | Problem Statement & Challenge | Solution in SQLTraceBench | Expected Outcome & Vision |
| :--- | :--- | :--- | :--- |
| **可扩展性 (Extensibility)** | The number of database systems is vast and growing. Hardcoding logic for each DB is unsustainable. | A plugin-based architecture. A `DatabasePlugin` interface defines a contract for all DB-specific operations (SQL dialect, schema mapping, driver). | Any developer can add support for a new database (e.g., Oracle, Snowflake) by implementing a single interface. The ecosystem can grow organically. |
| **可维护性 (Maintainability)** | Complex logic for SQL AST manipulation, parameter modeling, and workload generation can become a monolithic "big ball of mud". | A layered architecture separates concerns: a `pkg` layer for reusable components (AST, config, trace format), an `internal` layer for orchestration, and a `plugins` directory for extensions. | Clear boundaries between modules reduce cognitive load, simplify debugging, and allow parallel development on different system parts. |
| **可用性 (Usability)** | The process is complex (schema analysis, trace parsing, workload config). Users need a simple, guided experience. | A single, powerful CLI (`sql_trace_bench`) with clear commands (`run`, `validate`, `analyze`) and a comprehensive YAML configuration file. | Users can execute a complex cross-DB benchmark with a single command. The configuration file serves as a declarative, reproducible recipe for the experiment. |
| **可靠性 (Reliability)** | Network failures, incorrect DB credentials, or unsupported SQL features can cause failures. Benchmarks must be resilient. | Centralized error handling, connection retries, and validation steps at each stage. Plugins report their capabilities (e.g., supported functions). | The system provides clear, actionable error messages. Benchmarks can gracefully handle transient errors and produce partial results where possible. |
| **性能 (Performance)** | The workload generator must be capable of producing high QPS and managing thousands of concurrent connections without being the bottleneck itself. | The core workload generator is written in Go, leveraging goroutines for lightweight, highly concurrent execution. | The tool can saturate modern analytical databases and accurately simulate high-traffic production environments, ensuring the benchmark measures the DB, not the tool. |

## 2. System Workflow

The SQLTraceBench pipeline processes data through several distinct stages, orchestrated by the central command. Each stage produces an artifact that serves as the input for the next.

```mermaid
graph TD
    %% Main workflow for SQLTraceBench
    subgraph INPUT [输入层（Input Layer）]
        A1[SQL Trace<br/>(e.g., starrocks_trace.jsonl)]
        A2[DB Schema<br/>(e.g., starrocks_schema.yml)]
        A3[配置<br/>(config.yaml)]
    end

    subgraph CORE [核心处理流水线（Core Processing Pipeline）]
        B1[解析与加载（Parsing & Loading）] --> B2[SQL模板化（SQL Templating）]
        B2 --> B3[参数分布建模（Parameter Modeling）]
        B3 --> B4[跨数据库转换（Cross-DB Translation）]
        B4 --> B5[工作负载生成（Workload Generation）]
    end

    subgraph EXEC [执行与验证层（Execution & Validation Layer）]
        C1[负载执行器（Load Executor）] --> C2[性能指标收集（Metrics Collection）]
        C2 --> C3[偏差分析与报告（Deviation Analysis & Reporting）]
    end
    
    subgraph OUTPUT [输出层（Output Layer）]
        D1[Benchmark 报告<br/>(report.md/.json)]
    end
    
    A1 & A2 & A3 --> B1
    
    B5 --> C1
    
    C3 --> D1

    classDef input fill:#cde4ff,stroke:#5b97d1,stroke-width:2px;
    classDef core fill:#d5e8d4,stroke:#82b366,stroke-width:2px;
    classDef exec fill:#ffe6cc,stroke:#d79b00,stroke-width:2px;
    classDef output fill:#f8cecc,stroke:#b85450,stroke-width:2px;
    
    class A1,A2,A3 input;
    class B1,B2,B3,B4,B5 core;
    class C1,C2,C3 exec;
    class D1 output;

````

**流程阐述:**

1.  **输入层 (Input Layer)**: 用户提供三个核心输入：一份源数据库的 SQL 查询日志（`SQL Trace`）、一个定义了源数据库表结构的声明式文件（`DB Schema`），以及一个`config.yaml`文件，该文件指定了源/目标数据库类型、连接信息和工作负载参数。
2.  **解析与加载 (Parsing & Loading)**: 系统首先解析配置文件，并根据`source.db_type`加载对应的源数据库插件。该插件负责解析特定格式的 Trace 和 Schema 文件，将其转换为系统内部的统一表示。
3.  **SQL 模板化 (SQL Templating)**: 遍历所有查询，使用 AST（抽象语法树）解析器将 SQL 语句中的字面量（如 `'10.176.84.34'`, `1740672289`）替换为占位符（`?`），生成查询模板。同时，提取出的参数被保存下来用于下一步。
4.  **参数分布建模 (Parameter Modeling)**: 对上一步提取的所有参数进行统计分析，为每个参数（如 `tenant`, `recordTimestamp`）建立概率分布模型。此模型将用于在测试时生成符合原始数据特征的参数。
5.  **跨数据库转换 (Cross-DB Translation)**: 根据`target.db_type`加载目标数据库插件。此阶段执行两个关键转换：
      * **Schema 转换**: 目标插件将内部 Schema 表示转换为目标数据库兼容的 DDL 语句。
      * **SQL 转换**: 目标插件对每个 SQL 模板应用一系列转换规则，将其改写为目标数据库的方言。
6.  **工作负载生成 (Workload Generation)**: 结合转换后的 SQL 模板和参数模型，生成一个可执行的基准测试计划。此计划是一个包含大量具体 SQL 查询的序列，准备好被执行器运行。
7.  **负载执行器 (Load Executor)**: 根据`workload`配置（如并发数、时长），执行器启动多个工作 goroutine，连接到目标数据库，并从工作负载计划中抽取查询来执行。
8.  **性能指标收集 (Metrics Collection)**: 在执行期间，实时收集关键性能指标（KPIs），如查询延迟（平均、P95、P99）、QPS、扫描/返回行数等。
9.  **偏差分析与报告 (Deviation Analysis & Reporting)**: 测试结束后，将收集到的性能指标与从源 Trace 中分析出的基线指标进行对比，计算偏差，并生成一份详细的、人类可读的最终报告。

## 3. Core Components

```mermaid
graph TD
    %% Component Diagram
    CLI[CLI<br/>(Cobra)]

    subgraph Orchestration [编排层（Orchestration Layer）]
        Orchestrator[Orchestrator]
    end
    
    subgraph CoreServices [核心服务层（Core Services Layer）]
        TraceParser[Trace Parser]
        SchemaParser[Schema Parser]
        Templater[SQL Templater (AST)]
        Modeler[Parameter Modeler]
        WorkloadGen[Workload Generator]
        Executor[Load Executor]
        Reporter[Reporter]
    end
    
    subgraph PluginSystem [插件系统（Plugin System）]
        PluginManager[Plugin Manager]
        PluginInterface[DatabasePlugin Interface]
        StarRocksPlugin[StarRocks Plugin]
        ClickHousePlugin[ClickHouse Plugin]
        PostgresPlugin[PostgreSQL Plugin]
        MorePlugins[...]
    end
    
    subgraph Shared [共享基础库（Shared Libraries）]
        Config[Config Loader]
        Logger[Logging]
        InternalTypes[Internal Types]
        Errors[Error Handling]
    end

    CLI --> Orchestrator
    Orchestrator --> PluginManager
    Orchestrator --> TraceParser & SchemaParser & Templater & Modeler & WorkloadGen & Executor & Reporter
    
    PluginManager --> PluginInterface
    PluginInterface <.-> StarRocksPlugin & ClickHousePlugin & PostgresPlugin & MorePlugins
    
    CoreServices --> PluginInterface
    CoreServices --> Shared
    Orchestration --> Shared
    
    classDef cli fill:#f5f5f5,stroke:#333,stroke-width:2px;
    classDef orchestration fill:#e1d5e7,stroke:#9673a6,stroke-width:2px;
    classDef services fill:#d5e8d4,stroke:#82b366,stroke-width:2px;
    classDef plugins fill:#ffe6cc,stroke:#d79b00,stroke-width:2px;
    classDef shared fill:#cde4ff,stroke:#5b97d1,stroke-width:2px;
    
    class CLI cli;
    class Orchestrator orchestration;
    class TraceParser,SchemaParser,Templater,Modeler,WorkloadGen,Executor,Reporter services;
    class PluginManager,PluginInterface,StarRocksPlugin,ClickHousePlugin,PostgresPlugin,MorePlugins plugins;
    class Config,Logger,InternalTypes,Errors shared;
```

  * **CLI Engine (Cobra)**: `cmd/sql_trace_bench` - 负责解析命令行参数和子命令，是用户与系统的主要交互入口。它初始化配置和日志，并调用 `Orchestrator` 来执行核心逻辑。
  * **Orchestrator**: `internal/orchestrator` - 系统的“大脑”。它负责执行端到端的工作流程，按照顺序调用各个核心服务（解析、转换、执行等），并管理它们之间的状态和数据传递。
  * **Plugin System**: `pkg/plugins`
      * **Plugin Interface (`DatabasePlugin`)**: 定义了所有数据库插件必须实现的契约。这包括 SQL 方言转换、数据类型映射、Schema 转换以及提供数据库驱动等功能。
      * **Plugin Manager**: 负责根据配置动态加载和管理相应的源/目标数据库插件。
      * **Concrete Plugins** (`plugins/starrocks`, `plugins/clickhouse`): `DatabasePlugin` 接口的具体实现，封装了特定数据库的所有知识。
  * **Core Services**:
      * **Parsers (`pkg/trace`, `pkg/schema`)**: 负责将输入文件（Trace, Schema）解析成统一的内部数据结构。
      * **SQL Templater (`pkg/ast`)**: 封装了底层的 SQL AST 解析库，提供将具体查询抽象为模板的核心能力。
      * **Load Executor (`pkg/workload`)**: 高性能的负载执行引擎，管理并发的 worker 池，执行查询并收集原始性能数据。
  * **Shared Libraries**:
      * **`pkg/config`**: 使用 Viper 等库加载和验证 `config.yaml`。
      * **`internal/log`**: 提供全局统一的日志记录器实例 (Logger)。
      * **`pkg/types`**: 定义项目中的核心数据结构，如 `Query`, `Schema`, `Template`, `TraceEvent` 等。
      * **`pkg/errors`**: 集中定义项目中的自定义错误类型，便于统一处理。

## 4. Project Directory Structure

The project will follow the standard Go project layout to ensure clarity and scalability.

```
SQLTraceBench/
├── .github/
│   └── workflows/
│       └── go.yml           # CI/CD workflow
├── api/                     # OpenAPI specs, Protobuf definitions (if any)
├── assets/                  # Logos, etc.
├── cmd/
│   └── sql_trace_bench/     # Main application entry point
│       └── main.go
├── configs/                 # Example configuration files
│   └── sample_config.yaml
├── docs/
│   ├── architecture.md      # This file
│   └── assets/
│       └── sql_trace_bench-demo.gif
├── internal/
│   ├── orchestrator/        # Core workflow logic
│   └── log/                 # Internal logging setup
├── pkg/
│   ├── ast/                 # SQL AST parsing and templating
│   ├── config/              # Configuration loading
│   ├── errors/              # Custom error types
│   ├── plugins/             # Plugin interfaces and manager
│   ├── schema/              # Schema definition and parsing
│   ├── trace/               # Trace definition and parsing
│   ├── types/               # Core data types (enums, structs)
│   ├── datasynth/           # Data synthesis tools
│   ├── validation/          # Validation and reporting logic
│   └── workload/            # Workload execution engine
├── plugins/                 # Concrete plugin implementations
│   ├── clickhouse/
│   └── starrocks/
├── scripts/                 # Helper scripts (build, demo, etc.)
│   └── sql_trace_bench-demo
├── test/
│   └── e2e/                 # End-to-end tests
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── README-zh.md
```

## 5. Data Formats

### 5.1 SQL Trace Format (`.jsonl`)

A JSONL file where each line is a JSON object representing one query event. This format is simple to parse and stream.

```json
{"timestamp": "2025-08-15T10:00:01Z", "query": "SELECT * FROM users WHERE id = 123;", "execution_time_ms": 10.5, "rows_returned": 1}
{"timestamp": "2025-08-15T10:00:02Z", "query": "SELECT name FROM products WHERE category = 'electronics' LIMIT 10;", "execution_time_ms": 25.0, "rows_returned": 10}
```

### 5.2 DB Schema Format (`.yml`)

A YAML file provides a human-readable, database-agnostic way to define schemas.

```yaml
database: my_app
tables:
  - name: users
    columns:
      - name: id
        type: INT
        properties: [PRIMARY_KEY, NOT_NULL]
      - name: username
        type: VARCHAR(255)
      - name: created_at
        type: TIMESTAMP
  - name: products
    columns:
      - name: product_id
        type: BIGINT
        properties: [PRIMARY_KEY]
      - name: name
        type: VARCHAR(500)
      - name: category
        type: VARCHAR(100)
```

### 5.3 Configuration Template (`config.yaml`)

The main configuration file driving the tool.

```yaml
# config.yaml
source:
  # 'starrocks', 'clickhouse', 'postgresql', etc.
  db_type: "starrocks"
  schema_file: "./schemas/starrocks_schema.yml"
  trace_file: "./traces/starrocks_trace.jsonl"
  
  # Optional: connection details to source DB for live data analysis
  data_source: 
    host: "starrocks-host"
    port: 9030
    user: "user"
    password: "password"

target:
  db_type: "clickhouse"
  # Connection details for the target DB to run the benchmark
  host: "clickhouse-host"
  port: 9000
  user: "default"
  password: ""
  
  # Action to take on schema before benchmark: 'none', 'drop-create'
  schema_setup_action: "drop-create" 

workload:
  concurrency: 64
  duration: "5m"
  qps_scale: 1.0 # 1.0 means replay at original average QPS
  # More options like warmup period, query selection, etc.
  
validation:
  enabled: true
  # Tolerable deviation for key metrics
  kpi_deviation_threshold:
    qps: 0.10 # 10%
    latency_p95: 0.15 # 15%
```