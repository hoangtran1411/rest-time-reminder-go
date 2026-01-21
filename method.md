# RestTimeReminder - Core Concept Analysis

## 📋 Tổng Quan Project

**Mục đích:** Ứng dụng chạy nền (background service) nhắc nhở người dùng nghỉ ngơi định kỳ bằng cách phát âm thanh chuông.

**Công nghệ hiện tại:** C# Windows Service (.NET Framework 4.7.2)

---

## 🧠 Core Concepts

### 1. **Time-Based Trigger (Kích hoạt dựa trên thời gian)**

```
Logic: Kiểm tra thời gian hiện tại mỗi giây
       → Nếu phút = 00 hoặc phút = 30 → Phát âm thanh
       → Sleep 60 giây để tránh phát lặp
```

**Implementation hiện tại (C#):**
```csharp
while (true) {
    DateTime currentTime = DateTime.Now;
    if (currentTime.Minute == 30 || currentTime.Minute == 00) {
        SoundPlayer simpleSound = new SoundPlayer(@"C:\\bell.wav");
        simpleSound.Play();
        Thread.Sleep(60000); // Tránh phát lặp trong cùng 1 phút
    }
    Thread.Sleep(1000); // Check mỗi giây
}
```

**Trigger Rules:**
- Phát chuông khi phút = `00` (đầu giờ)
- Phát chuông khi phút = `30` (giữa giờ)
- Tần suất mặc định: Mỗi 30 phút

---

### 2. **Background Process (Tiến trình chạy nền)**

**Đặc điểm:**
- Chạy độc lập, không cần tương tác người dùng
- Tự động khởi động cùng hệ thống
- Chạy liên tục 24/7

**Kiểu chạy nền hiện tại:** Windows Service
- Yêu cầu `InstallUtil.exe` để cài đặt
- Chạy với quyền LocalSystem
- Quản lý qua Services.msc

---

### 3. **Audio Notification (Thông báo âm thanh)**

**Cách thức:**
- Sử dụng file âm thanh `.wav`
- Đường dẫn cố định: `C:\bell.wav`
- Phát âm thanh đơn giản (không loop)

---

## 🔄 Workflow

```
┌─────────────────────────────────────────────────────────┐
│                    SERVICE START                         │
└─────────────────────────┬───────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│               INFINITE LOOP                              │
│  ┌────────────────────────────────────────────────────┐ │
│  │  1. Lấy thời gian hiện tại                         │ │
│  │  2. Kiểm tra phút = 00 hoặc 30                     │ │
│  │     ├─ YES → Phát âm thanh → Sleep 60s             │ │
│  │     └─ NO  → Tiếp tục                              │ │
│  │  3. Sleep 1 giây                                   │ │
│  │  4. Lặp lại                                        │ │
│  └────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## 🚀 Phiên Bản Go - Kế Hoạch

### Ưu Điểm của Go cho Project này:

| Feature | Lợi ích |
|---------|---------|
| **Single Binary** | Không cần .NET runtime, dễ phân phối |
| **Cross-Platform** | Có thể chạy trên Linux/macOS ngoài Windows |
| **Goroutines** | Xử lý concurrent nhẹ hơn Thread |
| **time.Ticker** | API scheduling tốt hơn |
| **Small Footprint** | Memory usage thấp |
| **Easy Deployment** | Chỉ cần copy 1 file binary |

### Cải Tiến Đề Xuất:

#### 1. **Flexible Scheduling**
```go
// Thay vì hardcode 00/30, có thể config
type Config struct {
    Interval      time.Duration // Mặc định 30 phút
    SoundFile     string        // Đường dẫn file âm thanh
    TriggerMinutes []int        // Danh sách phút kích hoạt [0, 30]
}
```

#### 2. **Multiple Run Modes**
```go
// Hỗ trợ nhiều chế độ chạy
const (
    ModeConsole   = "console"   // Chạy foreground (debug)
    ModeService   = "service"   // Windows Service
    ModeDaemon    = "daemon"    // Linux systemd
    ModeTray      = "tray"      // System Tray với UI
)
```

#### 3. **Better Audio Support**
```go
// Hỗ trợ nhiều format audio
type AudioPlayer interface {
    Play(filePath string) error
    PlayFromEmbed() error  // Nhúng âm thanh vào binary
    Stop()
}
```

#### 4. **Configuration File**
```yaml
# config.yaml
reminder:
  interval: 30m  # 30 phút
  sound: "bell.wav"
  
schedule:
  type: "interval"  # hoặc "fixed_minutes"
  minutes: [0, 30]
  
notification:
  sound: true
  desktop: true  # Có thể thêm desktop notification
```

#### 5. **Graceful Shutdown**
```go
// Xử lý tắt service đúng cách
func (s *Service) Start() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Lắng nghe tín hiệu OS (SIGTERM, SIGINT)
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
    
    go func() {
        <-signals
        cancel() // Graceful shutdown
    }()
    
    s.run(ctx)
}
```

---

## 📁 Đề Xuất Cấu Trúc Project Go

```
rest-time-reminder-go/
├── cmd/
│   └── reminder/
│       └── main.go           # Entry point
├── internal/
│   ├── scheduler/
│   │   └── scheduler.go      # Time-based scheduling logic
│   ├── audio/
│   │   └── player.go         # Audio playback
│   ├── notification/
│   │   └── notifier.go       # Desktop notifications
│   ├── config/
│   │   └── config.go         # Configuration handling
│   └── service/
│       └── service.go        # Windows Service wrapper
├── assets/
│   └── bell.wav              # Embedded sound file
├── config.yaml               # Configuration file
├── go.mod
├── go.sum
└── README.md
```

---

## 🔧 Thư Viện Go Đề Xuất

| Thư viện | Mục đích |
|----------|----------|
| `github.com/robfig/cron/v3` | Cron-style scheduling |
| `github.com/faiface/beep` | Audio playback (cross-platform) |
| `github.com/gen2brain/beeep` | Desktop notifications |
| `github.com/kardianos/service` | Windows Service / Linux daemon |
| `github.com/spf13/viper` | Configuration management |
| `embed` (built-in Go 1.16+) | Nhúng file vào binary |

---

## 📝 Implementation Checklist

- [ ] **Phase 1: Core Logic**
  - [ ] Tạo scheduler với time.Ticker
  - [ ] Implement audio player
  - [ ] Basic console app chạy được

- [ ] **Phase 2: Configuration**
  - [ ] Config file support (YAML/JSON)
  - [ ] CLI flags
  - [ ] Environment variables

- [ ] **Phase 3: Service Mode**
  - [ ] Windows Service support
  - [ ] Linux systemd support
  - [ ] Install/uninstall commands

- [ ] **Phase 4: Enhanced Features**
  - [ ] Desktop notifications
  - [ ] System tray icon (optional)
  - [ ] Embed sound file
  - [ ] Graceful shutdown

- [ ] **Phase 5: Distribution**
  - [ ] Cross-compilation
  - [ ] Release automation
  - [ ] Installer script

---

## 💡 Quick Start Code (Go)

```go
package main

import (
    "log"
    "time"
)

func main() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        now := time.Now()
        minute := now.Minute()
        
        if minute == 0 || minute == 30 {
            playSound("bell.wav")
            time.Sleep(60 * time.Second) // Tránh phát lặp
        }
    }
}

func playSound(path string) {
    log.Printf("🔔 Playing sound at %s\n", time.Now().Format("15:04:05"))
    // TODO: Implement audio playback
}
```
 
---

## 🔗 Tài Liệu Tham Khảo

- [kardianos/service](https://github.com/kardianos/service) - Go Windows Service
- [faiface/beep](https://github.com/faiface/beep) - Audio for Go
- [Go embed](https://pkg.go.dev/embed) - Embed files in binary
- [robfig/cron](https://github.com/robfig/cron) - Cron scheduler for Go
