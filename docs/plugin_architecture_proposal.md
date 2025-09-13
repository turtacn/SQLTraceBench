# Proposal: Migrating to an Out-of-Process Plugin System

## 1. Executive Summary

This document proposes a critical architectural enhancement for SQLTraceBench: migrating from the current Go native plugin system (`buildmode=plugin`) to a more robust, maintainable, and extensible out-of-process plugin system using [HashiCorp's `go-plugin` library](https://github.com/hashicorp/go-plugin).

This change is essential for achieving the project's goals of becoming a production-grade, cross-platform, and community-driven tool. The current plugin model is too brittle and restrictive for long-term growth and stability.

## 2. Current Architecture and Its Limitations

SQLTraceBench currently uses `go build -buildmode=plugin` to create database-specific plugins as shared object (`.so`) files. These are loaded dynamically by the main application at runtime.

While this is a native Go feature, it comes with severe limitations that make it unsuitable for this project:

-   **Extreme Brittleness**: The host application and all plugins **must** be compiled with the exact same version of Go. Furthermore, they must be compiled with the exact same versions of all dependencies listed in `go.mod`. Any mismatch, however small, will cause the application to fail at runtime when it tries to load the plugin. This makes the build and release process incredibly fragile.

-   **Platform Incompatibility**: This plugin model is only supported on Linux and macOS. It **does not work on Windows**, which immediately prevents a significant portion of the developer community from using or contributing to the project.

-   **Difficult Contribution Workflow**: New contributors wanting to add support for a new database would need to set up a development environment that perfectly mirrors the main project's environment. This creates a high barrier to entry and discourages community contributions.

-   **No Fault Tolerance**: A panic or crash within a plugin will immediately crash the entire main application, as they share the same process space.

## 3. Proposed Solution: Out-of-Process Plugins with `go-plugin`

The proposed solution is to adopt HashiCorp's `go-plugin` library. This library is the de-facto standard for building extensible Go applications and is battle-tested in major open-source projects like Terraform, Vault, and Nomad.

### How It Works

With `go-plugin`, each plugin is a standalone executable that the main SQLTraceBench application launches as a child process. Communication between the main application and the plugin happens securely and efficiently over a gRPC interface.

The main application discovers and starts the plugin executables, which then "handshake" over stdout before establishing the gRPC connection.

### Key Benefits

Adopting this model provides immediate and significant advantages:

1.  **Robust and Decoupled**: Plugins are completely independent processes. They can be built with different versions of Go and can even have conflicting dependencies with the host application. This eliminates the brittleness of the current system.

2.  **True Cross-Platform Support**: Since plugins are just executables communicating over a standard protocol, this system works flawlessly on Linux, macOS, and **Windows**.

3.  **Enhanced Security and Stability**: A crash in a plugin process will not take down the main SQLTraceBench application. The host can detect the crash and handle it gracefully.

4.  **Simplified Contribution**: A developer wanting to add a new database plugin can simply implement the required Go interface and build their plugin as a standard executable. They do not need to worry about the host application's complex build environment. This dramatically lowers the barrier to contribution.

5.  **Polyglot Potential**: While not an immediate goal, this architecture opens the door for plugins to be written in other languages that support gRPC, such as Python or Rust.

## 4. Implementation Sketch

Migrating to this new system would involve the following steps:

1.  **Define the gRPC Service**: The `DatabasePlugin` interface and its related interfaces (`SchemaConverter`, `QueryTranslator`) would be translated into a Protocol Buffers (`.proto`) service definition.

2.  **Refactor the Plugin Manager**: The `PluginManager` in the host application would be updated to use the `go-plugin` client library to discover and launch plugin executables.

3.  **Update Database Plugins**: Each existing plugin (e.g., for ClickHouse, StarRocks) would be refactored into a `main` package that uses the `go-plugin` server library to serve the gRPC service.

4.  **Adjust the Build Process**: The `Makefile` would be changed to build each plugin as a separate executable file (e.g., `sqltracebench-plugin-clickhouse`) instead of a shared object (`.so`) file.

## 5. Trade-Offs

-   **Performance Overhead**: Communication via gRPC has higher latency than an in-process function call. However, for the tasks performed by SQLTraceBench (e.g., schema conversion, file I/O, trace analysis), this overhead is negligible and is a small price to pay for the massive gains in stability, cross-platform support, and maintainability.

-   **Added Dependency**: This approach adds gRPC and Protocol Buffers to the project's dependency tree. However, these are standard, well-supported technologies in the cloud-native ecosystem.

## 6. Conclusion and Recommendation

The current `buildmode=plugin` architecture is a critical flaw that will inhibit the growth, stability, and adoption of SQLTraceBench.

It is strongly recommended that the project **prioritize migrating to the `go-plugin` out-of-process model**. This architectural change will provide a robust, cross-platform, and contributor-friendly foundation, directly enabling the project to achieve its ambition of becoming a production-grade, open-source standard for database benchmarking.
