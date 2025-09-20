# Claude Code Session Extender

A Go-based automation tool that extends your Claude Code usage beyond the standard 5-hour session limits without additional costs.

## 🎯 Purpose

Claude Code has a built-in 5-hour session limit. This tool automatically manages session cycling to provide continuous availability by strategically starting new sessions before you begin work, effectively doubling your usable Claude Code time.

## 🧠 How It Works

### The Strategy
- **Pre-emptive Sessions**: Starts Claude Code sessions 3 hours before your typical work hours
- **Overlapping Windows**: Creates multiple fresh 5-hour sessions during your work period
- **Automatic Cycling**: Handles session cleanup and restart scheduling
- **Continuous Coverage**: Ensures Claude Code is always available when you need it

### Example Scenario
If you work from **7:00 AM - 12:00 PM** (5 hours):
- 🕐 **4:02 AM**: Session #1 starts (fresh 5-hour limit)
- 🕘 **7:00 AM**: You begin work with a fresh session
- 🕘 **9:02 AM**: Session #2 starts (overlapping fresh limit)
- 📈 **Result**: You get **2 full sessions** instead of just 1!

This effectively **doubles your available Claude Code time**.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Claude Code CLI installed and accessible via `claude` command

### Installation
1. Clone or download this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Basic Usage

#### Use Default Schedule (Optimized for 7 AM start)
```bash
go run main.go
```
**Default times**: 4:02, 9:02, 14:02, 19:02, 0:02

#### Custom Schedule
```bash
go run main.go "2 6,11,16,21,2 * * *"  # For 9 AM - 2 PM work schedule
```

#### With Options
```bash
# Send initial message at startup
go run main.go --run-at-start

# Disable response printing
go run main.go --print-responses=false

# Both options with custom schedule
go run main.go --run-at-start --print-responses=false "2 8,13,18,23,4 * * *"
```

## ⚙️ Configuration

### Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--run-at-start` | `false` | Send message immediately at startup |
| `--print-responses` | `true` | Print Claude's responses to console |

### Cron Schedule Format
Uses standard cron syntax: `minute hour day month weekday`

```
# Examples:
"2 4,9,14,19,0 * * *"    # Every 5 hours starting at 4:02 AM
"*/30 * * * *"           # Every 30 minutes
"0 9 * * 1-5"            # 9 AM on weekdays only
```

## 📅 Schedule Customization

### For Different Work Hours

The key principle: **Start sessions 3 hours before your work begins**

| Your Work Hours | Recommended Schedule | Command |
|----------------|---------------------|---------|
| 6:00 AM - 11:00 AM | 3:02, 8:02, 13:02, 18:02, 23:02 | `"2 3,8,13,18,23 * * *"` |
| 7:00 AM - 12:00 PM | 4:02, 9:02, 14:02, 19:02, 0:02 | Default schedule |
| 9:00 AM - 2:00 PM | 6:02, 11:02, 16:02, 21:02, 2:02 | `"2 6,11,16,21,2 * * *"` |
| 10:00 AM - 3:00 PM | 7:02, 12:02, 17:02, 22:02, 3:02 | `"2 7,12,17,22,3 * * *"` |

### Why These Times Work
- **3-hour buffer**: Ensures sessions are fresh when you start work
- **5-hour intervals**: Maximizes overlap during your work period
- **24-hour coverage**: Provides continuous availability throughout the day

## 📊 Expected Output

When running, you'll see:
```
╔══════════════════════════════════════════════════════════════════════╗
║                    Claude Code Session Extender                      ║
╠══════════════════════════════════════════════════════════════════════╣
║  Purpose: Automatically extend Claude Code usage beyond 5-hour       ║
║           session limits without additional costs                    ║
╚══════════════════════════════════════════════════════════════════════╝

Using default schedule: 4:02 AM and every 5 hours after
Times: 4:02, 9:02, 14:02, 19:02, 0:02
Starting Claude Code automation with schedule: 2 4,9,14,19,0 * * *
Will wait for first scheduled time
Press Ctrl+C to stop

[2025-09-20 09:02:00] Starting new Claude Code session...
Message 'hi' sent successfully
Claude's response:
---
[stderr] Tip: You can launch Claude Code with just `claude`
---
Claude Code session ended
[2025-09-20 09:02:12] Session completed
```

## 🔧 Troubleshooting

### Common Issues

**"command not found: claude"**
- Ensure Claude Code CLI is installed and in your PATH
- Try running `claude --version` to verify installation

**"No response received (timeout)"**
- This is normal - Claude Code may not always provide immediate responses
- The session is still being created and will be available for use

**Sessions not starting at expected times**
- Verify your cron expression using online cron validators
- Check system time zone settings

### Debugging
Enable response printing to see what Claude Code outputs:
```bash
go run main.go --print-responses=true
```

## 🏗️ Building

### Create Executable
```bash
# For current platform
go build -o claude-extender main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o claude-extender.exe main.go

# For Linux
GOOS=linux GOARCH=amd64 go build -o claude-extender main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o claude-extender main.go
```

### Run as Service (Linux/macOS)
Create a systemd service or launchd daemon to run continuously in the background.

## 📝 License

This project is provided as-is for educational and personal use. Please respect Claude Code's terms of service.

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

## ⚠️ Disclaimer

This tool automates interaction with Claude Code CLI. Ensure compliance with Anthropic's terms of service. Use responsibly and ethically.