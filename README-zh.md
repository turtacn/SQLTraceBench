<div align="center">
  <img src="logo.png" alt="SQLTraceBench Logo" width="200" height="200">
  <h1>SQLTraceBench</h1>
  <p>
    <strong>一个面向跨数据库性能分析的通用、基于追踪驱动的基准测试框架</strong>
  </p>
  <p>
    <a href="https://github.com/turtacn/SQLTraceBench/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="许可证"></a>
    <a href="https://goreportcard.com/report/github.com/turtacn/SQLTraceBench"><img src="https://goreportcard.com/badge/github.com/turtacn/SQLTraceBench" alt="Go Report"></a>
    <a href="https://github.com/turtacn/SQLTraceBench/releases"><img src="https://img.shields.io/github/v/release/turtacn/SQLTraceBench" alt="最新版本"></a>
  </p>
  <p>
    <a href="README.md"><strong>English</strong></a> | <a href="docs/architecture.md"><strong>Architecure Design</strong></a>
  </p>
</div>

想彻底革新开发者和 DBA 对数据库进行基准测试和比较的方式。  
该框架能够帮助你在不同数据库技术之间 **回放真实的生产工作负载**，确保你的数据库选择 **基于数据，而不仅仅是厂商的宣传**。

## 核心目标

SQLTraceBench 是一个 **基于 SQL 追踪驱动** 的基准测试系统，能够读取已有的 **SQL Trace 与数据库模式**，为不同目标数据库生成等效的工作负载，并执行全面的性能对比。  
它提供了灵活的负载控制、并发度调节和数据生成能力，是数据库评估、迁移和性能调优的必备工具。

## 为什么选择 SQLTraceBench？

数据库基准测试一向困难重重。像 TPC-H/DS 这样的合成基准虽有价值，但往往无法反映你应用独特的查询模式与数据倾斜。  
SQLTraceBench 填补了这一空白：

* **真实工作负载**：直接使用生产环境的 SQL 追踪，能够高度还原应用的实际行为。  
* **跨数据库转换**：智能翻译 SQL 方言与数据库模式（例如从 StarRocks → ClickHouse，或 PostgreSQL → TiDB），实现真正的“同类对比”。  
* **可控回放**：不仅仅是重放。它会对查询做模板化、建模参数分布，并允许你调节 QPS、并发度和热点比例，从而模拟高峰流量或未来增长场景。  
* **可扩展设计**：插件机制强大，社区可以方便地扩展以支持新的数据库。  

## 关键特性

* **模式与 SQL 转换**：自动转换数据库模式与 SQL 查询，兼容不同数据库系统。  
* **SQL 模板化与参数化**：基于 AST 解析生成查询模板，并根据原始追踪建模查询参数分布。  
* **分布感知的数据合成**：分析源数据库的数据分布（基数、倾斜、频率），在目标数据库生成逼真且可扩展的合成数据集；规模可由你掌控。  
* **灵活的工作负载生成**：精细控制并发、QPS、查询混合比例和热点数据访问模式。  
* **性能验证**：对比基准测试的性能画像（QPS、延迟、扫描行数）与原始追踪，确保测试的真实性，并生成详细的偏差报告。  
* **命令行工具**：强大且易用的 CLI 工具 `sql_trace_bench`，可编排整个流程。  

## 快速开始

### 安装

请确保你已安装 Go（版本 1.21+）。  

```bash
go install github.com/turtacn/SQLTraceBench/cmd/sql_trace_bench@latest
````

### 基本用法

以下是一个运行基准测试的简单示例。假设你有一份来自 **StarRocks** 的 SQL 追踪，并希望在 **ClickHouse** 上测试等效的工作负载。

1. **准备配置文件 (`config.yaml`)：**

   ```yaml
   # config.yaml
   source:
     db_type: "starrocks"
     schema_file: "./schemas/starrocks_schema.yml"
     trace_file: "./traces/starrocks_trace.jsonl"
     # 数据合成需要连接源数据库
     data_source: 
       host: "starrocks-host"
       port: 9030
       user: "root"

   target:
     db_type: "clickhouse"
     host: "clickhouse-host"
     port: 9000
     user: "default"
     schema_setup_action: "drop-create" 

   # （可选）在运行工作负载前先在目标数据库生成合成数据
   data_synthesis:
     enabled: true
     # 为每张表定义扩展规则
     tables:
       - name: "users"
         # 生成源表行数的 10 倍
         scale_factor: 10.0
       - name: "network_security_log"
         # 或者直接指定生成固定行数
         target_rows: 100000000

   workload:
     concurrency: 64
     duration: "5d"
     qps_scale: 1.5  # 以源 QPS 的 1.5 倍回放
   ```

2. **运行基准测试：**

   ```bash
   sql_trace_bench run --config ./config.yaml
   ```

   此命令将依次完成：

   a. 解析 StarRocks 的模式文件与追踪文件。
   
   b. 转换模式以适配 ClickHouse。
   
   c. 翻译 SQL 查询为 ClickHouse 方言。
   
   d. 基于配置参数生成工作负载。
   
   e. 在 ClickHouse 实例执行基准测试并输出报告。

*(这是一个概念性示例，具体参数与配置详见文档。)*

## 贡献指南

我们正在构建一个对数据库性能与可靠性充满热情的开发者社区。
非常欢迎你的贡献！无论是增加新的数据库插件、改进 SQL 翻译逻辑，还是完善文档，都十分重要。

请阅读 [贡献指南](./CONTRIBUTING.md) 开始参与。

## 社区

欢迎加入社区频道，提问、分享想法并与其他用户交流：

* **GitHub Discussions**：用于提问与讨论。
* **Slack/Discord**：即将开放。

## 许可证

SQLTraceBench 使用 [Apache 2.0 许可证](./LICENSE)。
