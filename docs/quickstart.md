# Quick Start Guide

This guide will walk you through the essential commands to get started with SQLTraceBench. We will cover the end-to-end workflow, from generating a workload to running a benchmark and validating the results.

## 1. Prepare a Sample Trace File

First, you need a raw SQL trace file. The file should be in JSONL format, where each line is a JSON object representing a single query. Each object must contain at least a `query` field.

Create a file named `traces.jsonl` with the following content:

```json
{"timestamp": "2025-01-01T12:00:00Z", "query": "SELECT * FROM users WHERE id = 1", "latency": 120000000}
{"timestamp": "2025-01-01T12:00:01Z", "query": "SELECT * FROM products WHERE sku = 'abc'", "latency": 150000000}
{"timestamp": "2025-01-01T12:00:02Z", "query": "SELECT * FROM users WHERE id = 2", "latency": 110000000}
```

## 2. Generate a Workload

Next, use the `generate` command to create a workload file from your raw traces. This command analyzes the traces, extracts SQL templates, and generates a new set of queries that mimic the original workload's characteristics.

Run the following command:

```bash
./bin/sqltracebench generate --source-traces traces.jsonl --out workload.json --count 100
```

This will create a `workload.json` file containing 100 queries.

## 3. Run a Benchmark

Now you can use the `run` command to execute the generated workload against a target database. For this example, we will use the built-in `mock` plugin, which simulates query execution without needing a real database.

First, run the benchmark and save the results as a "base" run:

```bash
./bin/sqltracebench run --workload workload.json --out base_metrics.json --db mock
```

Next, run the benchmark again to get a "candidate" set of metrics for comparison:

```bash
./bin/sqltracebench run --workload workload.json --out candidate_metrics.json --db mock
```

You will now have two files, `base_metrics.json` and `candidate_metrics.json`, each containing the results of a benchmark run.

## 4. Validate the Results

Finally, use the `validate` command to compare the two benchmark runs. This command will generate an HTML report that highlights any performance deviations.

Run the following command:

```bash
./bin/sqltracebench validate --base base_metrics.json --candidate candidate_metrics.json --out ./report
```

This will create a `validation_report.html` file inside a new `report` directory. Open this file in your browser to see a detailed comparison of the two benchmark runs, including QPS and latency metrics.
