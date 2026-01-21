# Contributing to RestTimeReminder

First off, thank you for considering contributing to RestTimeReminder! 🎉

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating a bug report, please check if the issue already exists. When creating a bug report, include:

- **Clear title** describing the issue
- **Steps to reproduce** the behavior
- **Expected behavior** vs **actual behavior**
- **System information** (OS, Go version, etc.)
- **Logs or error messages** if available

### 💡 Suggesting Features

We love feature suggestions! When suggesting a feature:

- Use a **clear and descriptive title**
- Provide a **detailed description** of the proposed feature
- Explain **why this feature would be useful**
- Include **mockups or examples** if applicable

### 🔧 Code Contributions

1. **Look for good first issues** - Issues labeled `good first issue` are great for newcomers
2. **Comment on the issue** - Let others know you're working on it
3. **Follow the coding standards** - See below

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- A code editor with Go support (VS Code with Go extension recommended)

### Getting Started

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/rest-time-reminder-go.git
cd rest-time-reminder-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Build the application
go build -o rest-time-reminder ./cmd/reminder

# Run in development mode
go run ./cmd/reminder
```

### Project Structure

```
├── cmd/reminder/     # Application entry point
├── internal/         # Private application code
│   ├── scheduler/    # Scheduling logic
│   ├── audio/        # Audio playback
│   ├── notification/ # Desktop notifications
│   ├── config/       # Configuration handling
│   └── service/      # Service management
├── assets/           # Static assets (sounds, etc.)
└── docs/             # Documentation
```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Keep functions small and focused
- Write meaningful comments for exported functions

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests with coverage
go test -cover ./...
```

### File Organization

- One package per directory
- Keep related functions together
- Use meaningful file names (e.g., `scheduler.go`, `player.go`)

### Error Handling

```go
// ✅ Good - Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to play sound: %w", err)
}

// ❌ Bad - Returning raw error
if err != nil {
    return err
}
```

### Logging

```go
// Use structured logging with log/slog
slog.Info("reminder triggered",
    "time", time.Now(),
    "interval", config.Interval,
)
```

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(scheduler): add support for cron expressions
fix(audio): handle missing sound file gracefully
docs(readme): update installation instructions
test(scheduler): add unit tests for interval parsing
```

## Pull Request Process

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards

3. **Write/update tests** for your changes

4. **Run all checks**:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ```

5. **Commit your changes** using conventional commits

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request** with:
   - Clear title and description
   - Reference to related issues (if any)
   - Screenshots/demos for UI changes

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated (if needed)
- [ ] Tests added/updated
- [ ] All tests pass locally
- [ ] No new warnings introduced

## ❓ Questions?

Feel free to open an issue with your question or reach out to the maintainers.

---

Thank you for contributing! 🙌
