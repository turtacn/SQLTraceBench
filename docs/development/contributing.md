# Contributing to SQL Trace Bench

First off, thanks for taking the time to contribute!

## Code of Conduct
This project and everyone participating in it is governed by the [Code of Conduct](../../CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs
*   Check the issue tracker to ensure the bug hasn't already been reported.
*   Use the [Bug Report Template](../../.github/ISSUE_TEMPLATE/bug_report.md).
*   Include a reproduction script or detailed steps to reproduce.

### Suggesting Enhancements
*   Use the [Feature Request Template](../../.github/ISSUE_TEMPLATE/feature_request.md).
*   Clearly explain the rationale and use case.

### Pull Requests
1.  Fork the repo and create your branch from `main`.
2.  If you've added code that should be tested, add tests.
3.  If you've changed APIs, update the documentation.
4.  Ensure the test suite passes (`go test ./...`).
5.  Make sure your code lints (`golangci-lint run`).

## Style Guide
*   We follow standard Go formatting (`gofmt`).
*   Comments should be clear and concise.
*   Variable names should be descriptive.

## Development Workflow
See the [Development Guide](./development_guide.md) for details on setting up your environment.
