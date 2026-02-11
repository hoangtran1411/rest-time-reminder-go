# Agent Instructions - Rest Time Reminder

This document provides essential context and rules for AI agents (like Antigravity) working on the **Rest Time Reminder** project.

## 🤖 Role & Context
You are an expert Go developer assisting with the maintenance and enhancement of a cross-platform reminder application. Your goal is to ensure the codebase remains clean, idiomatic (Effective Go), and follows the defined architecture.

## 📌 Project Overview
- **Name**: Rest Time Reminder (Go version)
- **Goal**: A lightweight service/CLI that plays a bell sound at regular intervals to remind users to take breaks.
- **Tech Stack**: Go 1.21+, Cobra (CLI), Viper (Config), Kardianos Service (OS Service), Beep (Audio).
- **Core Principles**: Clean Architecture, SOLID, High Testability.

## 📏 Critical Rules (MUST FOLLOW)

### 1. Go Style & Idioms
- Follow **Effective Go** and [Uber Go Style Guide](https://github.com/uber-go/guide).
- Use `gofmt` and `goimports` for formatting.
- **Error Handling**: Use `fmt.Errorf("context: %w", err)` for wrapping. Don't log and return the same error.
- **Internal Packages**: All core logic MUST reside in `internal/`. Entry points are in `cmd/`.

### 2. Project Structure
```text
internal/
  scheduler/    - Time-based scheduling logic
  audio/        - Audio playback (Beep)
  notification/ - Desktop notifications (Beeep)
  config/       - Configuration management (Viper)
  service/      - OS Service wrapper (Kardianos)
cmd/
  reminder/     - CLI entry point (Cobra)
```

### 3. Linting (Mandatory)
- Use `golangci-lint` (v2 configuration schema).
- Run `make lint` before finishing any task.
- Configuration is in `.golangci.yml`.

### 4. Service & CLI
- CLI commands use `spf13/cobra`.
- Config management uses `spf13/viper`.
- Support graceful shutdown for all modes (Console, Service).

## 🛠️ Development Workflow

### Useful Commands (Makefile)
- `make build`: Build the binary to `build/rest-time-reminder`.
- `make test`: Run all tests.
- `make lint`: Run golangci-lint.
- `make run`: Run the application in console mode using `go run`.
- `make clean`: Remove build artifacts.

### Service Installation
Service installation is managed by the application binary:
- `rest-time-reminder install`: Install as a system service.
- `rest-time-reminder start`: Start the installed service.
- `rest-time-reminder stop`: Stop the service.
- `rest-time-reminder uninstall`: Remove the service.

### Configuration
- Default config is in `config.yaml`.
- Agents should respect the `config.yaml` schema when suggesting changes.

## 📂 File Edit Guidelines
- **Look Before You Leap**: Always check existing implementations in `internal/` before adding new logic.
- **Atomic Edits**: Use `replace_file_content` or `multi_replace_file_content` for surgical edits.
- **Single Responsibility**: Keep functions small and focused on one task.

## 🧠 Brain Dump (Context)
- The app handles embedded audio (bell.wav) if no custom sound is provided.
- Service mode requires administrative privileges for installation/uninstallation.
- macOS support is planned but primarily tested on Windows and Linux.

---
*Last Updated: 2026-02-11*
