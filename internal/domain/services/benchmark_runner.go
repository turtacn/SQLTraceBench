package services

import (
    "context"
    "fmt"
    "sync"
    "time"
    "github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type BenchmarkConfig struct {
    TraceCount  int
    Concurrency int
    Timeout     time.Duration
}

type BenchmarkResult struct {
    ModelName        string
    Throughput       float64 // traces/sec
    AvgLatency       float64 // ms
    P95Latency       float64
    P99Latency       float64
    MemoryUsageMB    float64
    CPUUsagePercent  float64
    ValidationScore  float64 // 0-1
    ErrorRate        float64 // 0-1
}

// TraceGeneratorModel defines the interface for workload generation models in benchmarking.
type TraceGeneratorModel interface {
    Name() string
    // Generate generates a specific number of traces.
    Generate(ctx context.Context, count int) ([]models.SQLTrace, error)
}

type BenchmarkRunner struct {
    Models  []TraceGeneratorModel
    Analyzer *PerformanceAnalyzer
}

func NewBenchmarkRunner(models []TraceGeneratorModel) *BenchmarkRunner {
    return &BenchmarkRunner{
        Models:   models,
        Analyzer: &PerformanceAnalyzer{},
    }
}

func (r *BenchmarkRunner) Run(
    ctx context.Context,
    config BenchmarkConfig,
) ([]BenchmarkResult, error) {
    var results []BenchmarkResult
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, model := range r.Models {
        wg.Add(1)
        go func(m TraceGeneratorModel) {
            defer wg.Done()

            result := r.runSingleModel(ctx, m, config)

            mu.Lock()
            results = append(results, result)
            mu.Unlock()
        }(model)
    }

    wg.Wait()
    return results, nil
}

func (r *BenchmarkRunner) runSingleModel(
    ctx context.Context,
    model TraceGeneratorModel,
    config BenchmarkConfig,
) BenchmarkResult {
    startTime := time.Now()
    latencies := make([]float64, 0, config.TraceCount)
    var latenciesMu sync.Mutex
    errors := 0
    var errorsMu sync.Mutex

    // WORKER POOL IMPLEMENTATION

    // Channel to distribute work items (number of traces to generate)
    // We can distribute in batches to reduce channel overhead, e.g., batch size 1.
    // Given the requirement is "trace_count", we can send indices or just units.
    jobs := make(chan int, config.TraceCount)

    // Fill the job queue
    go func() {
        for i := 0; i < config.TraceCount; i++ {
            jobs <- 1
        }
        close(jobs)
    }()

    var wg sync.WaitGroup

    // Start workers
    concurrency := config.Concurrency
    if concurrency <= 0 {
        concurrency = 1
    }

    for w := 0; w < concurrency; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for range jobs {
                genStart := time.Now()
                // Generating 1 trace per job
                _, err := model.Generate(ctx, 1)
                latency := time.Since(genStart).Milliseconds()

                if err != nil {
                    errorsMu.Lock()
                    errors++
                    errorsMu.Unlock()
                } else {
                    latenciesMu.Lock()
                    latencies = append(latencies, float64(latency))
                    latenciesMu.Unlock()
                }
            }
        }()
    }

    wg.Wait()
    duration := time.Since(startTime).Seconds()

    throughput := float64(config.TraceCount) / duration
    avgLatency := r.Analyzer.CalculateAverage(latencies)
    p95 := r.Analyzer.CalculatePercentile(latencies, 0.95)
    p99 := r.Analyzer.CalculatePercentile(latencies, 0.99)

    memUsage, cpuUsage := r.Analyzer.GetResourceUsage()

    errorRate := 0.0
    if config.TraceCount > 0 {
        errorRate = float64(errors) / float64(config.TraceCount)
    }

    // Placeholder for validation logic
    validationScore := 0.5

    fmt.Printf("Model %s finished: Throughput=%.2f, P99=%.2f\n", model.Name(), throughput, p99)

    return BenchmarkResult{
        ModelName:       model.Name(),
        Throughput:      throughput,
        AvgLatency:      avgLatency,
        P95Latency:      p95,
        P99Latency:      p99,
        MemoryUsageMB:   memUsage,
        CPUUsagePercent: cpuUsage,
        ValidationScore: validationScore,
        ErrorRate:       errorRate,
    }
}
