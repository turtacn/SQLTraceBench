# Phase 1: Streaming Trace Processing

## Overview

To address Out-Of-Memory (OOM) risks when processing large SQL trace files (e.g., 1GB+), we have introduced a **Streaming Trace Parser**. This replaces the previous method of loading the entire file into memory before processing.

## Architecture

### Streaming Trace Parser

The `StreamingTraceParser` reads JSONL (newline-delimited JSON) files line by line.

*   **Location**: `internal/infrastructure/parsers/streaming_trace_parser.go`
*   **Mechanism**: Uses `bufio.Scanner` to read the file.
*   **Buffer**: Configurable buffer size (default 1MB per line) to handle long SQL queries.
*   **Error Handling**: Skips malformed lines and logs errors without aborting the entire process.

### Conversion Service

The `ConversionService` has been refactored to use the `StreamingTraceParser`.

*   **Interface**: `Parse(reader io.Reader, callback func(models.SQLTrace) error)`
*   **Usage**: The service provides a callback function that processes each trace as it is parsed.

### CLI Adaptation

The `convert` command (`cmd/convert.go`) now supports streaming output.

*   **Streaming Mode**: Triggered when the output file extension is `.jsonl`.
*   **Behavior**: Reads traces streamingly, optionally translates them using a plugin, and writes them immediately to the output file. This ensures minimal memory footprint (< 500MB peak for 1GB input).
*   **Legacy Mode**: If the output is `.json` (default), it collects traces in memory to perform template extraction (aggregation). Note that this mode may still consume significant memory for large files until the Template Service is also refactored for streaming.

## Performance

*   **Memory**: The streaming approach maintains a stable memory profile regardless of input file size, determined mainly by the size of the largest single query and the batch size (if batched).
*   **Throughput**: Parsing speed is limited by I/O and JSON unmarshalling CPU time.

## Configuration

Buffer sizes and error thresholds can be configured in `configs/trace_parser.yaml`.
