# SQLTraceBench

SQLTraceBench is a tool for benchmarking database performance using trace-based workload generation.

## Features

- **Trace-Driven**: Generates workloads based on real SQL traces.
- **Multi-Database Support**: Architecture allows plugins for different databases (e.g., MySQL, ClickHouse, StarRocks).
- **Statistical Modeling**: Uses statistical models to synthesize realistic parameter values.
- **Extensible**: Easily add new database dialects and workload patterns.
- **Visual Reports**: Generates comprehensive HTML validation reports with interactive charts.

## Getting Started

Check out our **[Quick Start Guide](docs/quickstart.md)** to get up and running in minutes!

It covers:
* Building the project
* Preparing data
* Running conversion, generation, and benchmarking commands
* Using the automated `examples/quickstart.sh` script

## Validation Reports

SQLTraceBench generates detailed validation reports to compare your benchmark results against a baseline.

![Validation Report](docs/images/report_preview.png)

*Example of an HTML validation report showing QPS deviation and latency distribution.*

See [Report Interpretation Guide](docs/user_guide/report_interpretation.md) for details on how to read the reports.

## Development

1.  **Build**: `make build`
2.  **Test**: `make test`

## Architecture

See [docs/architecture/architecture.md](docs/architecture/architecture.md) for details on the system design.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
