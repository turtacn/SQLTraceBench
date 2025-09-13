# Contributing to SQLTraceBench

First off, thank you for considering contributing to SQLTraceBench! Your help is greatly appreciated. This project is built by developers for developers, and we welcome contributions of all kinds, from bug fixes to new features to documentation improvements.

This document provides guidelines for contributing to the project.

## Getting Started

1.  **Fork** the repository on GitHub.
2.  **Clone** your fork locally:
    ```bash
    git clone https://github.com/YOUR_USERNAME/SQLTraceBench.git
    cd SQLTraceBench
    ```
3.  **Create a new branch** for your changes:
    ```bash
    git checkout -b feature/my-amazing-feature
    ```

## Prerequisites

Before you can build and test SQLTraceBench, you'll need to have the following tools installed:

-   **Go**: Version 1.20 or later.
-   **Docker**: Required for running integration and end-to-end tests that depend on live databases.
-   **Make**: Used to run common development commands.

## Development Setup

Once you have the prerequisites, you can set up your local development environment:

1.  **Install Dependencies**: The project uses Go Modules to manage dependencies. They will be automatically downloaded when you build or test the project. To ensure everything is in sync, you can run:
    ```bash
    go mod tidy
    ```

2.  **Install Development Tools**: We use `golangci-lint` for linting. You can install it using the following command:
    ```bash
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    ```
    **Note**: Ensure that `$(go env GOPATH)/bin` is in your shell's `PATH`.

## How to Run Tests

The project has a growing test suite. To run all tests, use the `make` command:

```bash
make test
```

This will run all unit and integration tests in the project.

## How to Run the Linter

We use `golangci-lint` to enforce code quality and style. To run the linter, use the `make` command:

```bash
make lint
```

**Note**: There is a known issue with the linter's `typecheck` module in some environments. If you encounter persistent `typecheck` errors, you may need to investigate the `golangci-lint` cache or configuration. The CI pipeline currently bypasses this issue to prevent blocking development.

## Submitting Contributions

When you are ready to submit your changes, please follow these steps:

1.  **Commit** your changes with a clear and descriptive commit message:
    ```bash
    git commit -m "feat: Add some amazing feature"
    ```
2.  **Push** your changes to your fork on GitHub:
    ```bash
    git push origin feature/my-amazing-feature
    ```
3.  **Open a Pull Request** from your fork's branch to the `main` branch of the original SQLTraceBench repository.
4.  In your pull request description, please explain the changes you made and why. If your PR addresses an open issue, please reference it (e.g., `Closes #123`).

Thank you again for your contribution!
