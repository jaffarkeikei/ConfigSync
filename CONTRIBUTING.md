# Contributing to ConfigSync

Thank you for your interest in contributing to ConfigSync! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to uphold our Code of Conduct. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Bugs

If you find a bug in the project:

1. Check if the bug has already been reported in the Issues section.
2. If not, create a new issue with a clear title and description.
3. Include steps to reproduce the bug, expected behavior, and actual behavior.
4. Include information about your environment (OS, Kubernetes version, etc.).

### Suggesting Enhancements

If you have an idea for an enhancement:

1. Check if the enhancement has already been suggested in the Issues section.
2. If not, create a new issue with a clear title and description.
3. Explain why this enhancement would be useful to most users.

### Contributing Code

1. Fork the repository.
2. Create a new branch for your feature or bugfix: `git checkout -b feature/your-feature-name` or `git checkout -b fix/your-bugfix-name`.
3. Make your changes.
4. Run the tests: `make test`.
5. Commit your changes with a clear and descriptive commit message.
6. Push your branch to your fork: `git push origin your-branch-name`.
7. Create a pull request to the main repository.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker
- Kubernetes cluster (local or remote)
- kubectl
- git

### Setting Up Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/config-sync.git
   cd config-sync
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run the tests:
   ```bash
   make test
   ```

## Project Structure

- `cmd/`: Contains the main application entry points
- `pkg/`: Contains the library code
  - `apis/`: API definitions and types
  - `controller/`: Kubernetes controllers
  - `git/`: Git operations
  - `validator/`: Configuration validation
- `deploy/`: Kubernetes deployment manifests
- `configs/`: Example configurations
- `scripts/`: Utility scripts
- `examples/`: Example usage

## Pull Request Process

1. Ensure all tests pass.
2. Update the documentation if necessary.
3. Make sure your code follows the project's style guidelines.
4. The PR will be merged once it receives approval from at least one maintainer.

## Coding Standards

- Follow the Go style guide: use `gofmt` and `golint`.
- Write tests for new features and bug fixes.
- Document your code with comments.
- Keep pull requests focused on a single change.

## License

By contributing to this project, you agree that your contributions will be licensed under the project's MIT License. 