# Wach

> **German for "awake" (pronounced "vakh")** — a minimal macOS menu bar app that keeps your Mac awake by subtly moving the mouse when you're away.

Wach is a lightweight, native macOS utility that prevents your Mac from going idle and keeps messaging apps (Slack, Teams, etc.) from switching your status to "Away." Unlike caffeine-style apps that only prevent sleep, Wach simulates **real user activity** by moving the mouse cursor — just enough to keep the system attentive, without disrupting your workflow.

## Features

- **🔵 Smart idle detection** — only moves the mouse when you're actually away (>60s of inactivity)
- **🌓 Automatic appearance adaptation** — icon switches between black (light mode) and white (dark mode) via macOS template rendering
- **👁 Visual state indicator** — open eye icon when active, closed eye when stopped
- **🪶 Native & lightweight** — uses CoreGraphics APIs directly, not bloated cross-platform libraries
- **🍎 Apple Silicon native** — built as a native `arm64` binary, optimized for M1/M2/M3/M4 Macs
- **🔌 No network access** — all processing is local, no data leaves your machine

## How it Works

```
┌──────────────────────────────────────────────────────┐
│  Every 10 seconds:                                   │
│                                                      │
│  1. Check idle time via CoreGraphics API             │
│  2. Idle > 60s? → Check if display is awake         │
│  3. Display awake? → Move mouse by 10px              │
│  4. Move direction alternates (prevents drift)       │
│                                                      │
│  If user activity detected → wait, do nothing        │
│  If display sleeping → skip (don't wake it)          │
└──────────────────────────────────────────────────────┘
```

### Why "Wach"?

"Wach" is German for **awake / alert / vigilant**. It's short, memorable, and describes exactly what the app does — it keeps your system awake while you're away from the keyboard.

## Requirements

- macOS 11.0 (Big Sur) or later
- Apple Silicon (M1/M2/M3/M4) or Intel
- [Accessibility permission](#granting-accessibility-access)

## Installation

### Option 1: Download the pre-built app

```bash
# Build from source (recommended):
git clone git@github.com:Smotherer007/wach.git
cd wach
make install   # builds + copies to /Applications/
```

### Option 2: Build and run directly

```bash
make run       # run directly from terminal (no .app bundle)
```

## First Launch

1. Open **Wach.app** from your Applications folder
2. **Right-click → Open** on first launch (app is not notarized)
3. Grant Accessibility permission when prompted (see below)

## Granting Accessibility Access

Wach needs Accessibility permission to move the mouse cursor. 

**System Settings → Privacy & Security → Accessibility → toggle Wach ON**

If Wach doesn't appear in the list, start the app once (right-click → Open), then check again.

## Usage

Click the **eye icon** in the menu bar (top-right, near the clock):

| Menu Item | Action |
|-----------|--------|
| **Start** | Begin mouse-mover loop (eye opens) |
| **Stop** | Pause mouse-mover loop (eye closes) |
| **About** | Version info and credits |
| **Quit** | Exit the app completely |

The app always starts automatically with the mover active. Use **Stop** if you need to step away without keeping the system awake.

## Development

### Project Structure

```
wach/
├── cmd/main.go                 # Entry point (systray menu bar app)
├── pkg/wach/
│   ├── wach.go                 # Core logic & lifecycle
│   ├── state.go                # Thread-safe state management
│   ├── cgo_darwin.go           # CoreGraphics C bindings
│   ├── alert_darwin.m          # Native macOS alert dialogs
│   ├── logger.go               # Logging (stdout or file)
│   ├── doc.go                  # Package documentation
│   └── wach_test.go            # 28 tests (unit + CGo bridge)
├── assets/icon/
│   ├── icon.go                 # Embedded tray icons (PNG byte arrays)
│   └── icon_test.go            # Icon validation tests
├── appInfo/
│   ├── Info.plist              # macOS bundle metadata
│   └── icon.icns               # App bundle icon
├── lib/systray/                # Patched systray library
│   └── systray_darwin.m        # 20pt icon size + template support
├── Makefile                    # build, install, run, test targets
├── go.mod / go.sum
└── README.md
```

### Commands

```bash
make install    # Build + copy to /Applications/
make build      # Build .app bundle in ./bin/
make run        # Run directly (no bundle)
make test       # Run all tests with race detector
make clean      # Remove build artifacts
```

### Technical Details

**Why not use robotgo?** The original [automatic-mouse-mover](https://github.com/prashantgupta24/automatic-mouse-mover) by prashantgupta24 used `robotgo` and `activity-tracker` — cross-platform libraries that pull in hundreds of dependencies (Windows/Linux APIs, OCR libraries, etc.). This fork replaces all of that with:

| Feature | Original | Wach |
|---------|----------|------|
| Mouse movement | robotgo (CGo + C++) | **CoreGraphics** (C) |
| Idle detection | activity-tracker | **CGEventSourceSecondsSinceLastEventType** |
| Dialogs | gosx-notifier | **NSAlert** (Objective-C) |
| Binary size | ~6.4 MB | **~4.2 MB** |
| Dependencies | ~40+ indirect | **~15 indirect** |

**Patched systray library:** The `getlantern/systray` library hardcodes icons at 16×16 points. This fork patches it to 20×20 for better visibility on modern Mac displays, and uses template images for automatic light/dark mode support.

### Icon Design

- **Open eye** — the app is active and monitoring idle time
- **Closed eye (arc)** — the app is paused
- Template rendering — macOS automatically adjusts icon color

## Credits

- Original concept by [prashantgupta24/automatic-mouse-mover](https://github.com/prashantgupta24/automatic-mouse-mover)
- systray library by [getlantern/systray](https://github.com/getlantern/systray)
- Built with Go, CoreGraphics, and native macOS APIs

## License

MIT
