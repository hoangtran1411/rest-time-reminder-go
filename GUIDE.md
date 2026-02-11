# User Guide - Rest Time Reminder

Welcome to **Rest Time Reminder**, a lightweight application designed to help you maintain healthy work habits by reminding you to take regular breaks.

---

## 🚀 Quick Start

1. **Download/Build**: Get the binary for your operating system.
2. **Configure**: Customize your break interval and sound in `config.yaml`.
3. **Run**: Start the application in your preferred mode (Console or Service).

---

## ⚙️ Configuration

The application uses `config.yaml` for its settings. You can copy the `config.example.yaml` file to get started.

### Reminder Settings
- `interval`: The time between reminders (e.g., `30m`, `1h`).
- `trigger_minutes`: (Optional) Specific minutes of the hour to trigger (e.g., `[0, 30]` for every half-hour).

### Audio Settings
- `enabled`: Set to `true` to hear a bell or custom sound.
- `file`: Path to a `.wav` file. If left blank, the application uses an internal "bell" sound.
- `volume`: Control the loudness from `0.0` to `1.0`.

### Desktop Notifications
- `desktop`: Enable/disable system-level pop-up notifications.
- `title` & `message`: Customize the text shown in the notification.

---

## 💻 Usage Modes

### 1. Console Mode (Interactive)
The simplest way to run the app. It will run in your terminal and log notifications as they happen.
```bash
./rest-time-reminder
```
*Press `Ctrl+C` to stop.*

### 2. Service Mode (Background)
Ideal for daily use. The application runs in the background even when the terminal is closed.

**Windows (PowerShell as Admin):**
- **Install**: `./rest-time-reminder.exe install`
- **Start**: `./rest-time-reminder.exe start`
- **Stop**: `./rest-time-reminder.exe stop`
- **Status**: `./rest-time-reminder.exe status`

**Linux:**
- **Install**: `sudo ./rest-time-reminder install`
- Use `systemctl` to manage the service after installation.

---

## 🛠️ Command Line Interface

You can override `config.yaml` settings using CLI flags:

- `-c, --config`: Specify a custom config file path.
- `-i, --interval`: Set a temporary interval (e.g., `45m`).
- `-s, --sound`: Use a specific sound file for this session.
- `-v, --verbose`: See more detailed logs for debugging.

**Commands:**
- `version`: Show current version and check for updates.
- `update`: Automatically upgrade to the latest release from GitHub.

---

## 💡 Best Practices for Breaks

- **The 20-20-20 Rule**: Every 20 minutes, look at something 20 feet away for 20 seconds to reduce eye strain.
- **Stand Up**: Use the reminder to stand up, stretch, or walk around for at least 2 minutes.
- **Hydrate**: Use break time to grab a glass of water.

---

## ❓ Troubleshooting

- **No Sound?**: Ensure your system volume is up and the `file` path in `config.yaml` is correct (or leave it empty for the default bell).
- **Service Fails to Start?**: Check the logs (usually in `reminder.log` or system logs). On Windows, make sure you ran the command prompt as an **Administrator**.
- **Update Error?**: Ensure you have an active internet connection to reach GitHub.

---
*Stay healthy and productive!*
