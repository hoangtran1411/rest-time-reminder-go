---
trigger: always_on
---

# Go Style Guide - Rest Time Reminder

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is a **cross-platform reminder application (Service/CLI)** built with:
- **Cobra** for CLI commands
- **Kardianos Service** for OS service management
- **Internal packages** for modular business logic

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint` (v2.8.0+) `run ./...` before commit.
- **Linting Configuration**: MUST use `golangci-lint` v2 configuration schema (v2.8.x+). 
  - Top-level `version: "2"` is mandatory.
  - Use kebab-case for all linter settings.
  - Exclusions move to `linters: exclusions: rules`.
- Adhere to [Effective Go](https://go.dev/doc/effective_go).
- Core logic in `internal/` packages. CLI entry points in `cmd/`.

## Project Structure

```
rest-time-reminder-go/
  cmd/
    reminder/
      main.go         - Application entry point
  internal/
    scheduler/        - Time-based scheduling logic
    audio/            - Audio playback functionality
    notification/     - Desktop notifications
    config/           - Configuration management
    service/          - Windows/Linux service wrapper
  assets/             - Embedded resources (e.g. sounds)
  config.yaml         - Default configuration
```

## Error Handling

- Wrap errors: `fmt.Errorf("context: %w", err)`
- Guard clauses for fail-fast
- Do not log and return the same error
- One responsibility per layer

## Service & CLI Integration

- Use `spf13/cobra` for command structure.
- Use `spf13/viper` for configuration.
- Graceful shutdown handling for service mode.
- Log to file when running as service, stdout when console.

## Testing & Linting

- Table-driven tests with `t.Run`
- Target high coverage for `internal/` packages like `scheduler` and `config`.
- Use `go test ./...` and `golangci-lint run`.

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [golangci-lint](https://github.com/golangci/golangci-lint)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.
