# SQLTraceBench API Documentation

## Overview

SQLTraceBench provides a comprehensive set of APIs for SQL trace analysis, workload generation, and benchmark execution. This document describes all available interfaces and their usage.

## Table of Contents

- [Core Service APIs](#core-service-apis)
- [Data Models](#data-models)
- [Plugin Development APIs](#plugin-development-apis)
- [Configuration Reference](#configuration-reference)
- [Error Codes](#error-codes)
- [Examples](#examples)

## Core Service APIs

### ConversionService Interface

The ConversionService handles SQL trace and schema conversion between different database systems.

```go
type ConversionService interface {
    ConvertTraces(ctx context.Context, request *ConversionRequest) (*ConversionResult, error)
    ConvertSchema(ctx context.Context, schema *models.DatabaseSchema, targetDB types.DatabaseType) (*models.DatabaseSchema, error)
    ValidateConversion(ctx context.Context, result *ConversionResult) (*ValidationReport, error)
}
````

#### ConvertTraces

Converts SQL traces from source database format to target database format.

**Parameters:**

* `ctx`: Context for request lifecycle management
* `request`: ConversionRequest containing source traces and target database information

**Request Format:**

```go
type ConversionRequest struct {
    SourceDB     types.DatabaseType     `json:"source_db"`
    TargetDB     types.DatabaseType     `json:"target_db"`
    TracesFile   string                 `json:"traces_file"`
    SchemaFile   string                 `json:"schema_file,omitempty"`
    OutputPath   string                 `json:"output_path"`
    Options      ConversionOptions      `json:"options,omitempty"`
}

type ConversionOptions struct {
    PreserveComments bool   `json:"preserve_comments"`
    OptimizeQueries  bool   `json:"optimize_queries"`
    BatchSize        int    `json:"batch_size"`
    ParallelWorkers  int    `json:"parallel_workers"`
}
```

**Response Format:**

```go
type ConversionResult struct {
    ConvertedTraces    []models.SQLTrace      `json:"converted_traces"`
    ConvertedSchema    *models.DatabaseSchema `json:"converted_schema,omitempty"`
    ConversionStats    ConversionStats        `json:"conversion_stats"`
    Warnings          []string               `json:"warnings,omitempty"`
    OutputFiles       []string               `json:"output_files"`
}

type ConversionStats struct {
    TotalQueries      int64         `json:"total_queries"`
    ConvertedQueries  int64         `json:"converted_queries"`
    FailedQueries     int64         `json:"failed_queries"`
    ProcessingTime    time.Duration `json:"processing_time"`
    ConversionRate    float64       `json:"conversion_rate"`
}
```

**Usage Example:**

```go
conversionService := conversion.NewDefaultConversionService(
    templateService,
    schemaService,
    pluginRegistry,
)

request := &ConversionRequest{
    SourceDB:   types.DatabaseTypeStarRocks,
    TargetDB:   types.DatabaseTypeClickHouse,
    TracesFile: "input/starrocks_traces.jsonl",
    SchemaFile: "input/schema.sql",
    OutputPath: "output/",
    Options: ConversionOptions{
        PreserveComments: true,
        OptimizeQueries:  true,
        BatchSize:        1000,
        ParallelWorkers:  4,
    },
}

result, err := conversionService.ConvertTraces(ctx, request)
if err != nil {
    return fmt.Errorf("conversion failed: %w", err)
}

fmt.Printf("Converted %d queries successfully\n", result.ConversionStats.ConvertedQueries)
```

### GenerationService Interface

The GenerationService generates benchmark workloads based on SQL templates and parameter models.

```go
type GenerationService interface {
    GenerateWorkload(ctx context.Context, request *GenerationRequest) (*GenerationResult, error)
    OptimizeWorkload(ctx context.Context, workload *models.BenchmarkWorkload, criteria OptimizationCriteria) (*models.BenchmarkWorkload, error)
    EstimatePerformance(ctx context.Context, workload *models.BenchmarkWorkload) (*PerformanceEstimate, error)
}
```

#### GenerateWorkload

Generates a benchmark workload from SQL templates and parameter models.

**Parameters:**

* `ctx`: Context for request lifecycle management
* `request`: GenerationRequest containing templates and generation parameters

**Request Format:**

```go
type GenerationRequest struct {
    TemplateFile     string              `json:"template_file"`
    ParameterModel   string              `json:"parameter_model"`
    WorkloadConfig   WorkloadConfig      `json:"workload_config"`
    OutputPath       string              `json:"output_path"`
    GenerationMode   GenerationMode      `json:"generation_mode"`
}

type WorkloadConfig struct {
    TargetQPS        float64       `json:"target_qps"`
    Duration         time.Duration `json:"duration"`
    Concurrency      int           `json:"concurrency"`
    HotspotRatio     float64       `json:"hotspot_ratio"`
    QueryDistribution map[string]float64 `json:"query_distribution"`
}

type GenerationMode string

const (
    GenerationModeTime        GenerationMode = "time_based"
    GenerationModeCount       GenerationMode = "count_based"
    GenerationModeDistribution GenerationMode = "distribution_based"
)
```

**Response Format:**

```go
type GenerationResult struct {
    GeneratedWorkload *models.BenchmarkWorkload `json:"generated_workload"`
    GenerationStats   GenerationStats           `json:"generation_stats"`
    OutputFiles       []string                  `json:"output_files"`
    Recommendations   []string                  `json:"recommendations,omitempty"`
}

type GenerationStats struct {
    TotalQueries       int64         `json:"total_queries"`
    UniqueTemplates    int           `json:"unique_templates"`
    ParametersGenerated int64        `json:"parameters_generated"`
    GenerationTime     time.Duration `json:"generation_time"`
    EstimatedQPS       float64       `json:"estimated_qps"`
}
```

### ExecutionService Interface

The ExecutionService executes benchmark workloads against target databases.

```go
type ExecutionService interface {
    ExecuteWorkload(ctx context.Context, request *ExecutionRequest) (*ExecutionResult, error)
    MonitorExecution(ctx context.Context, executionID string) (<-chan ExecutionProgress, error)
    StopExecution(ctx context.Context, executionID string) error
}
```

#### ExecuteWorkload

Executes a benchmark workload against the specified database.

**Request Format:**

```go
type ExecutionRequest struct {
    WorkloadFile     string           `json:"workload_file"`
    DatabaseConfig   DatabaseConfig   `json:"database_config"`
    ExecutionConfig  ExecutionConfig  `json:"execution_config"`
    OutputPath       string           `json:"output_path"`
}

type DatabaseConfig struct {
    DatabaseType types.DatabaseType `json:"database_type"`
    Host         string             `json:"host"`
    Port         int                `json:"port"`
    Database     string             `json:"database"`
    Username     string             `json:"username"`
    Password     string             `json:"password"`
    MaxConns     int                `json:"max_connections"`
    Timeout      time.Duration      `json:"timeout"`
}

type ExecutionConfig struct {
    WarmupDuration   time.Duration `json:"warmup_duration"`
    ReportInterval   time.Duration `json:"report_interval"`
    EnableMonitoring bool          `json:"enable_monitoring"`
    FailureThreshold float64       `json:"failure_threshold"`
    RetryAttempts    int           `json:"retry_attempts"`
}
```

**Response Format:**

```go
type ExecutionResult struct {
    ExecutionID      string              `json:"execution_id"`
    WorkloadStats    WorkloadStats       `json:"workload_stats"`
    PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
    ErrorSummary     ErrorSummary        `json:"error_summary,omitempty"`
    ReportFiles      []string            `json:"report_files"`
}

type WorkloadStats struct {
    TotalQueries     int64         `json:"total_queries"`
    SuccessfulQueries int64        `json:"successful_queries"`
    FailedQueries    int64         `json:"failed_queries"`
    ExecutionTime    time.Duration `json:"execution_time"`
    ActualQPS        float64       `json:"actual_qps"`
}
```

### ValidationService Interface

The ValidationService validates benchmark results and generates analysis reports.

```go
type ValidationService interface {
    ValidateResults(ctx context.Context, request *ValidationRequest) (*ValidationResult, error)
    CompareMetrics(ctx context.Context, original, synthetic *PerformanceMetrics) (*ComparisonResult, error)
    GenerateReport(ctx context.Context, validation *ValidationResult, format ReportFormat) ([]byte, error)
}
```

## Data Models

### SQLTrace

Represents a single SQL trace entry from database logs.

```go
type SQLTrace struct {
    Timestamp       time.Time         `json:"timestamp"`
    QueryText       string            `json:"query_text"`
    QueryHash       string            `json:"query_hash"`
    ExecutionTime   time.Duration     `json:"execution_time"`
    RowsReturned    int64            `json:"rows_returned,omitempty"`
    RowsScanned     int64            `json:"rows_scanned,omitempty"`
    DatabaseName    string           `json:"database_name,omitempty"`
    Username        string           `json:"username,omitempty"`
    ClientIP        string           `json:"client_ip,omitempty"`
    QueryType       types.QueryType  `json:"query_type"`
    Tables          []string         `json:"tables,omitempty"`
    Metadata        map[string]interface{} `json:"metadata,omitempty"`
}
```

### SQLTemplate

Represents a parameterized SQL template extracted from traces.

```go
type SQLTemplate struct {
    TemplateID      string                `json:"template_id"`
    OriginalQuery   string                `json:"original_query"`
    TemplateQuery   string                `json:"template_query"`
    Parameters      []Parameter           `json:"parameters"`
    QueryType       types.QueryType       `json:"query_type"`
    Tables          []string              `json:"tables"`
    Frequency       int64                 `json:"frequency"`
    CreatedAt       time.Time             `json:"created_at"`
    UpdatedAt       time.Time             `json:"updated_at"`
    Stats           TemplateStats         `json:"stats,omitempty"`
}

type Parameter struct {
    Name         string                `json:"name"`
    Type         types.ParameterType   `json:"type"`
    Position     int                   `json:"position"`
    Distribution Distribution          `json:"distribution,omitempty"`
}

type TemplateStats struct {
    AvgExecutionTime time.Duration `json:"avg_execution_time"`
    MaxExecutionTime time.Duration `json:"max_execution_time"`
    MinExecutionTime time.Duration `json:"min_execution_time"`
    TotalExecutions  int64         `json:"total_executions"`
    ErrorRate        float64       `json:"error_rate"`
}
```

### BenchmarkWorkload

Represents a complete benchmark workload configuration.

```go
type BenchmarkWorkload struct {
    WorkloadID       string            `json:"workload_id"`
    Name             string            `json:"name"`
    Description      string            `json:"description,omitempty"`
    Queries          []WorkloadQuery   `json:"queries"`
    ExecutionConfig  ExecutionConfig   `json:"execution_config"`
    Schedule         QuerySchedule     `json:"schedule,omitempty"`
    CreatedAt        time.Time         `json:"created_at"`
    Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type WorkloadQuery struct {
    QueryID          string              `json:"query_id"`
    QueryText        string              `json:"query_text"`
    TemplateID       string              `json:"template_id,omitempty"`
    ExecutionTime    time.Duration       `json:"execution_time,omitempty"`
    Weight           float64             `json:"weight"`
    Parameters       map[string]interface{} `json:"parameters,omitempty"`
    ExpectedResult   *QueryResult        `json:"expected_result,omitempty"`
}
```

## Plugin Development APIs

### DatabasePlugin Interface

Base interface for all database plugins.

```go
type DatabasePlugin interface {
    GetName() string
    GetVersion() string
    ValidateConnection() error
    GetSchemaConverter() SchemaConverter
    GetQueryTranslator() QueryTranslator
    GetDataGenerator() DataGenerator
    GetBenchmarkExecutor() BenchmarkExecutor
}
```

### SchemaConverter Interface

Handles database schema conversion between different systems.

```go
type SchemaConverter interface {
    ConvertSchema(sourceSchema *models.DatabaseSchema, targetDB types.DatabaseType) (*models.DatabaseSchema, error)
    ValidateSchema(schema *models.DatabaseSchema) error
    GetSupportedTypes() []string
}
```

### QueryTranslator Interface

Handles SQL query translation between different database dialects.

```go
type QueryTranslator interface {
    TranslateQuery(query string, targetDB types.DatabaseType) (string, error)
    ExtractTemplate(query string) (*models.SQLTemplate, error)
    ValidateQuery(query string) error
}
```

### Plugin Registration

Register custom database plugins with the system.

```go
// Register a new plugin
func RegisterPlugin(name string, factory PluginFactory) error {
    return pluginRegistry.Register(name, factory)
}

// PluginFactory creates plugin instances
type PluginFactory func(config map[string]interface{}) (DatabasePlugin, error)

// Example plugin registration
func init() {
    RegisterPlugin("custom-db", func(config map[string]interface{}) (DatabasePlugin, error) {
        return NewCustomDBPlugin(config)
    })
}
```

## Configuration Reference

### Application Configuration

```yaml
# Application settings
app:
  version: "1.0.0"
  log_level: "info"
  log_format: "json"
  temp_dir: "/tmp/sqltrace"
  
# Database connections
databases:
  starrocks:
    host: "localhost"
    port: 9030
    username: "root"
    password: ""
    max_connections: 100
    timeout: "30s"
    
  clickhouse:
    host: "localhost"
    port: 9000
    username: "default"
    password: ""
    database: "default"
    compression: "lz4"

# Benchmark settings
benchmark:
  default_qps: 100
  max_concurrency: 64
  default_duration: "5m"
  report_path: "./reports"
  
# Plugin settings  
plugins:
  directory: "./plugins"
  auto_load: true
  timeout: "10s"
```

### Environment Variables

| Variable              | Description                       | Default         |
| --------------------- | --------------------------------- | --------------- |
| `SQLTRACE_CONFIG`     | Configuration file path           | `./config.yaml` |
| `SQLTRACE_LOG_LEVEL`  | Log level (debug/info/warn/error) | `info`          |
| `SQLTRACE_TEMP_DIR`   | Temporary directory               | `/tmp/sqltrace` |
| `SQLTRACE_PLUGIN_DIR` | Plugin directory                  | `./plugins`     |
| `SQLTRACE_REPORT_DIR` | Report output directory           | `./reports`     |

## Error Codes

| Code | Name                  | Description                |
| ---- | --------------------- | -------------------------- |
| 1000 | ErrUnknown            | Unknown error occurred     |
| 1001 | ErrInvalidInput       | Invalid input parameters   |
| 1002 | ErrParseFailed        | SQL parsing failed         |
| 1003 | ErrConversionFailed   | Query conversion failed    |
| 1004 | ErrDatabaseConnection | Database connection error  |
| 1005 | ErrPluginNotFound     | Plugin not found           |
| 1006 | ErrValidationFailed   | Validation failed          |
| 1007 | ErrExecutionFailed    | Benchmark execution failed |
| 1008 | ErrReportGeneration   | Report generation failed   |

### Error Handling Example

```go
if err != nil {
    if sqlError, ok := err.(*types.SQLTraceBenchError); ok {
        switch sqlError.Code {
        case types.ErrDatabaseConnection:
            log.Error("Database connection failed", "details", sqlError.Details)
            // Handle connection error
        case types.ErrValidationFailed:
            log.Warn("Validation failed", "message", sqlError.Message)
            // Handle validation error
        default:
            log.Error("Unexpected error", "error", sqlError)
        }
    }
}
```

## Examples

### Complete Workflow Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/turtacn/SQLTraceBench/internal/app/conversion"
    "github.com/turtacn/SQLTraceBench/internal/app/generation"
    "github.com/turtacn/SQLTraceBench/internal/app/execution"
    "github.com/turtacn/SQLTraceBench/internal/app/validation"
    "github.com/turtacn/SQLTraceBench/pkg/types"
)

func main() {
    ctx := context.Background()
    
    // Step 1: Convert traces
    conversionService := conversion.NewDefaultConversionService(...)
    conversionReq := &conversion.ConversionRequest{
        SourceDB:   types.DatabaseTypeStarRocks,
        TargetDB:   types.DatabaseTypeClickHouse,
        TracesFile: "input/traces.jsonl",
        OutputPath: "output/converted/",
    }
    
    conversionResult, err := conversionService.ConvertTraces(ctx, conversionReq)
    if err != nil {
        fmt.Printf("Conversion failed: %v\n", err)
        return
    }
    
    // Step 2: Generate workload
    generationService := generation.NewDefaultGenerationService(...)
    generationReq := &generation.GenerationRequest{
        TemplateFile: conversionResult.OutputFiles[0],
        WorkloadConfig: generation.WorkloadConfig{
            TargetQPS:   100,
            Duration:    5 * time.Minute,
            Concurrency: 10,
        },
        OutputPath: "output/workload/",
    }
    
    generationResult, err := generationService.GenerateWorkload(ctx, generationReq)
    if err != nil {
        fmt.Printf("Generation failed: %v\n", err)
        return
    }
    
    // Step 3: Execute benchmark
    executionService := execution.NewDefaultExecutionService(...)
    executionReq := &execution.ExecutionRequest{
        WorkloadFile: generationResult.OutputFiles[0],
        DatabaseConfig: execution.DatabaseConfig{
            DatabaseType: types.DatabaseTypeClickHouse,
            Host:         "localhost",
            Port:         9000,
            Database:     "test",
        },
        OutputPath: "output/results/",
    }
    
    executionResult, err := executionService.ExecuteWorkload(ctx, executionReq)
    if err != nil {
        fmt.Printf("Execution failed: %v\n", err)
        return
    }
    
    // Step 4: Validate results
    validationService := validation.NewDefaultValidationService(...)
    validationReq := &validation.ValidationRequest{
        OriginalTraces:  "input/traces.jsonl",
        SyntheticResults: executionResult.ReportFiles[0],
        ValidationRules: validation.DefaultRules(),
    }
    
    validationResult, err := validationService.ValidateResults(ctx, validationReq)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }
    
    fmt.Printf("Benchmark completed successfully!\n")
    fmt.Printf("QPS Deviation: %.2f%%\n", validationResult.QPSDeviation*100)
    fmt.Printf("Latency Deviation: %.2f%%\n", validationResult.LatencyDeviation*100)
}
```

### Custom Plugin Development Example

```go
package myplugin

import (
    "github.com/turtacn/SQLTraceBench/pkg/plugins"
    "github.com/turtacn/SQLTraceBench/pkg/types"
)

// MyDatabasePlugin implements DatabasePlugin interface
type MyDatabasePlugin struct {
    config *MyDBConfig
    // ... other fields
}

func NewMyDatabasePlugin(config *MyDBConfig) (*MyDatabasePlugin, error) {
    return &MyDatabasePlugin{
        config: config,
    }, nil
}

func (p *MyDatabasePlugin) GetName() string {
    return "my-database"
}

func (p *MyDatabasePlugin) GetVersion() string {
    return "1.0.0"
}

func (p *MyDatabasePlugin) ValidateConnection() error {
    // Implement connection validation
    return nil
}

func (p *MyDatabasePlugin) GetSchemaConverter() plugins.SchemaConverter {
    return NewMySchemaConverter()
}

// Register the plugin
func init() {
    plugins.RegisterPlugin("my-database", func(config map[string]interface{}) (plugins.DatabasePlugin, error) {
        // Parse config and create plugin instance
        myConfig, err := parseMyConfig(config)
        if err != nil {
            return nil, err
        }
        return NewMyDatabasePlugin(myConfig)
    })
}
```

## Best Practices

### Performance Optimization

1. **Batch Processing**: Process traces in batches for better memory usage
2. **Parallel Execution**: Use multiple workers for CPU-intensive operations
3. **Connection Pooling**: Reuse database connections for better performance
4. **Memory Management**: Stream large files instead of loading entirely into memory

### Error Handling

1. **Use Structured Errors**: Always use SQLTraceBenchError for consistent error handling
2. **Context Cancellation**: Respect context cancellation in long-running operations
3. **Retry Logic**: Implement exponential backoff for transient failures
4. **Resource Cleanup**: Always clean up resources in defer statements

### Plugin Development

1. **Interface Compliance**: Ensure your plugin implements all required interfaces
2. **Configuration Validation**: Validate plugin configuration during initialization
3. **Error Reporting**: Provide detailed error messages with context
4. **Documentation**: Document plugin-specific features and limitations