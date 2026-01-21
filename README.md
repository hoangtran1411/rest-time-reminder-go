# 🔔 RestTimeReminder

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Windows%20|%20Linux%20|%20macOS-blue?style=for-the-badge)

**A lightweight, cross-platform reminder application that plays a bell sound at regular intervals to remind you to take breaks.**

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Configuration](#-configuration) • [Contributing](#-contributing)

</div>

---

## 📋 Overview

RestTimeReminder is a background service that helps you maintain healthy work habits by reminding you to take breaks. By default, it plays a pleasant bell sound every 30 minutes (at :00 and :30 of each hour).

### Why RestTimeReminder?

- **Health First**: Regular breaks reduce eye strain, prevent RSI, and improve focus
- **Lightweight**: Single binary, minimal resource usage (~5MB RAM)
- **Cross-Platform**: Works on Windows, Linux, and macOS
- **Flexible**: Run as console app, system service, or with system tray
- **Customizable**: Configure intervals, sounds, and notification methods

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🕐 **Flexible Scheduling** | Configure reminder intervals (default: 30 minutes) |
| 🔊 **Audio Notifications** | Play custom sound files or use embedded bell sound |
| 💻 **Multiple Run Modes** | Console, Windows Service, Linux daemon, or System Tray |
| ⚙️ **Configuration File** | YAML-based configuration for easy customization |
| 📦 **Single Binary** | No dependencies, easy deployment |
| 🔄 **Graceful Shutdown** | Properly handles system signals |
| 🖥️ **Desktop Notifications** | Optional desktop notification support |

---

## 📥 Installation

### Pre-built Binaries

Download the latest release from the [Releases](https://github.com/yourusername/rest-time-reminder-go/releases) page.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/rest-time-reminder-go.git
cd rest-time-reminder-go

# Build the binary
go build -o rest-time-reminder ./cmd/reminder

# Or install directly
go install ./cmd/reminder
```

### Requirements

- Go 1.21 or later (for building from source)
- Audio output device (for sound notifications)

---

## 🚀 Usage

### Console Mode (Default)

```bash
# Run with default settings (30-minute intervals)
./rest-time-reminder

# Run with custom interval
./rest-time-reminder --interval 45m

# Run with custom sound file
./rest-time-reminder --sound /path/to/custom-sound.wav
```

### Service Mode

#### Windows

```powershell
# Install as Windows Service
.\rest-time-reminder.exe install

# Start the service
.\rest-time-reminder.exe start

# Stop the service
.\rest-time-reminder.exe stop

# Uninstall the service
.\rest-time-reminder.exe uninstall
```

#### Linux (systemd)

```bash
# Install as systemd service
sudo ./rest-time-reminder install

# Start the service
sudo systemctl start rest-time-reminder

# Enable on boot
sudo systemctl enable rest-time-reminder
```

### Command Line Options

```
Usage: rest-time-reminder [OPTIONS] [COMMAND]

Commands:
  install     Install as system service
  uninstall   Remove system service
  start       Start the service
  stop        Stop the service
  status      Show service status

Options:
  -c, --config      Path to configuration file (default: config.yaml)
  -i, --interval    Reminder interval (e.g., 30m, 1h)
  -s, --sound       Path to custom sound file
  -m, --mode        Run mode: console, service, tray
  -v, --verbose     Enable verbose logging
  -h, --help        Show help message
  --version         Show version information
```

---

## ⚙️ Configuration

Create a `config.yaml` file in the same directory as the binary:

```yaml
# RestTimeReminder Configuration

reminder:
  # Interval between reminders (e.g., 30m, 1h, 45m)
  interval: 30m
  
  # Fixed trigger minutes (overrides interval if set)
  # trigger_minutes: [0, 30]

sound:
  # Enable/disable sound notifications
  enabled: true
  
  # Path to custom sound file (leave empty for embedded sound)
  file: ""
  
  # Sound volume (0.0 - 1.0)
  volume: 1.0

notification:
  # Enable/disable desktop notifications
  desktop: false
  
  # Notification title
  title: "Break Time!"
  
  # Notification message
  message: "Time to take a short break and rest your eyes."

logging:
  # Log level: debug, info, warn, error
  level: info
  
  # Log to file (leave empty for stdout only)
  file: ""
```

---

## 🏗️ Project Structure

```
rest-time-reminder-go/
├── cmd/
│   └── reminder/
│       └── main.go           # Application entry point
├── internal/
│   ├── scheduler/
│   │   └── scheduler.go      # Time-based scheduling logic
│   ├── audio/
│   │   └── player.go         # Audio playback functionality
│   ├── notification/
│   │   └── notifier.go       # Desktop notifications
│   ├── config/
│   │   └── config.go         # Configuration management
│   └── service/
│       └── service.go        # Windows/Linux service wrapper
├── assets/
│   └── bell.wav              # Embedded sound file
├── config.yaml               # Default configuration
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── CONTRIBUTING.md
```

---

## 🧰 Dependencies

| Library | Purpose |
|---------|---------|
| [kardianos/service](https://github.com/kardianos/service) | Cross-platform service management |
| [faiface/beep](https://github.com/faiface/beep) | Audio playback |
| [gen2brain/beeep](https://github.com/gen2brain/beeep) | Desktop notifications |
| [spf13/viper](https://github.com/spf13/viper) | Configuration management |
| [spf13/cobra](https://github.com/spf13/cobra) | CLI framework |

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by the original [RestTimeReminder](https://github.com/original/RestTimeReminder) C# project
- Thanks to all contributors who help improve this project

---

<div align="center">

**Made with ❤️ for healthier work habits**

⭐ Star this repo if you find it helpful!

</div>
