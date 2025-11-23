# P2: Parameter Modeling & Workload Generation

## 1. Overview

Phase 2 enhances the workload generation capabilities of SQLTraceBench by introducing advanced parameter modeling and temporal pattern recognition. This ensures that the generated workloads closely mimic the statistical properties of the original production traffic.

## 2. Key Components

### 2.1 Parameter Analyzer (`parameter_analyzer.go`)
- **Purpose**: Analyzes trace data to infer parameter types and calculate value frequencies.
- **Type Inference**: Automatically detects `INT`, `STRING`, and `DATETIME` types.
- **Stats**: Maintains frequency counts for all observed values.

### 2.2 Hotspot Detector (`hotspot_detector.go`)
- **Purpose**: Identifies "hot" parameter values that appear frequently.
- **Algorithm**: Selects the top N values that cumulatively account for a specific threshold (e.g., 5%) of the total occurrences, or simply the top frequent items.
- **Configuration**: `hotspot_threshold` (default 0.05).

### 2.3 Temporal Pattern Extractor (`temporal_pattern_extractor.go`)
- **Purpose**: Captures the query arrival rate over time.
- **Algorithm**: Buckets traces into time windows (e.g., 1 hour) and counts queries per window.
- **Output**: A probability distribution of query timestamps.

### 2.4 Enhanced Samplers
- **ZipfSampler (`zipf_sampler.go`)**:
  - Implements Zipfian distribution for parameter value selection.
  - **Hotspot Injection**: Supports overriding the natural tail with detected hotspot values (default 30% injection probability).
- **TemporalSampler (`temporal_sampler.go`)**:
  - Generates timestamps based on the extracted temporal pattern using weighted random sampling.

## 3. Workflow

1.  **Trace Parsing**: Raw SQL logs are parsed into `SQLTrace` objects.
2.  **Analysis**:
    - Parameters are extracted and typed.
    - Frequencies are counted.
    - Hotspots are identified.
    - Temporal patterns are extracted.
3.  **Generation**:
    - A template is selected (currently randomly from source).
    - Parameters are generated using `ZipfSampler` (with hotspots).
    - Timestamps are assigned using `TemporalSampler`.

## 4. Configuration

See `configs/generation.yaml` for tunable parameters.

```yaml
hotspot_threshold: 0.05
temporal_window: 1h
```
