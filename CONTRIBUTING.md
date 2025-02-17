# Contributing to Remote Certificate Exporter

First off, thank you for considering contributing to Remote Certificate Exporter! It's people like you that make this tool better for everyone.

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct: be respectful, constructive, and professional in all interactions.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include details about your configuration and environment

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* A clear and descriptive title
* A detailed description of the proposed enhancement
* Any possible drawbacks or alternatives you've considered
* If possible, a rough implementation sketch

### Pull Requests

1. Fork the repository
2. Create a new branch from `main` for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bug-fix
   ```

3. Make your changes, following our coding conventions:
   * Use Go formatting standards (`go fmt`)
   * Add tests for new features
   * Update documentation as needed
   * Follow existing code patterns and practices

4. Run the test suite:
   ```bash
   make test
   ```

5. Run the linter:
   ```bash
   make lint
   ```

6. Commit your changes using clear commit messages:
   ```bash
   git commit -m "A brief description of your changes"
   ```

7. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

8. Open a Pull Request

### Pull Request Guidelines

* Update the README.md with details of changes to the interface, if applicable
* Update the documentation with any new features or configuration options
* The PR should work for all supported Go versions
* Include unit tests for new features
* Follow the standard Go code formatting

## Development Setup

1. Install Go (1.19 or later recommended)
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Install development tools:
   ```bash
   make setup-dev
   ```

## Testing

* Run all tests:
  ```bash
  make test
  ```
* Run tests with coverage:
  ```bash
  make coverage
  ```

## Questions?

Feel free to open an issue with your question or reach out to the maintainers directly.

Thank you for contributing!
