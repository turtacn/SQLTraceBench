# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-12-01

### Added
- **Statistical Validation**: Added K-S test, Chi-Square, and Auto-correlation validation tests.
- **Benchmarking**: Introduced `benchmark` command to compare model performance.
- **API**: Added RESTful API (`/api/v2`) with OpenAPI 3.0 specification.
- **Visualization**: Added HTML report generation with interactive charts.
- **Metrics**: Added Prometheus metrics exporter and Grafana dashboards.
- **Models**: Added Transformer-based trace generator.
- **Documentation**: Added comprehensive documentation (User Guide, API Reference, Architecture).

### Changed
- **CLI**: Renamed `generate-trace` to `generate`, `verify-trace` to `validate`.
- **Config**: Unified configuration format for all models.
- **Performance**: Improved Markov chain generation speed by 3x.

### Fixed
- Fixed race condition in concurrent generation.
- Fixed memory leak in LSTM model training.
- Fixed timestamp generation precision issues.

## [1.5.0] - 2024-10-15

### Added
- Initial support for LSTM models.
- Basic CSV export format.

### Changed
- Refactored database connection logic to support plugins.

## [1.0.0] - 2024-06-01

### Added
- Initial release.
- Markov chain trace generation.
- PostgreSQL support.
