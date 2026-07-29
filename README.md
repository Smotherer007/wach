# Wach

> German for "awake" (pronounced "vakh") — a minimal macOS menu bar app that keeps your Mac awake by subtly moving the mouse when you're away.

Wach is a lightweight macOS utility that prevents your Mac from going idle and keeps messaging apps (Slack, Teams, etc.) from switching your status to "Away." Unlike caffeine-style apps that only prevent sleep, Wach simulates real user activity by moving the mouse cursor just enough to keep the system attentive, without disrupting your workflow.

## Features

- Smart idle detection — only moves the mouse when you are actually away (after 60s of inactivity)
- Automatic appearance adaptation — icon adapts between light and dark mode via macOS template rendering
- Visual state indicator — open eye icon when active, closed eye when stopped
- Native and lightweight — uses CoreGraphics APIs directly instead of bloated cross-platform libraries
- Apple Silicon native — built as a native arm64 binary, optimized for M-series Macs
- No network access — all processing is local, no data leaves your machine

## How it Works

Every 10 seconds, Wach checks the system idle time using the CoreGraphics API. If the system has been idle for more than 60 seconds and the display is awake, it moves the mouse cursor by 10 pixels (alternating direction to prevent drift). If there is user activity or the display is sleeping, no action is taken.

### Why "Wach"?

"Wach" is German for awake, alert, or vigilant. It is short, memorable, and describes exactly what the app does — it keeps your system awake while you are away from the keyboard.

## Requirements

- macOS 11.0 (Big Sur) or later
- Apple Silicon (M1/M2/M3/M4) or Intel
- Accessibility permission (see below)

## Installation

```bash
git clone git@github.com:Smotherer007/wach.git
cd wach
make install
```

This builds the app and copies it to /Applications/. You can also run `make run` to start directly from the terminal without creating an app bundle.

## First Launch

1. Open Wach.app from your Applications folder
2. Right-click and choose Open on first launch (the app is not notarized)
3. Grant Accessibility permission when prompted

## Granting Accessibility Access

Wach needs Accessibility permission to move the mouse cursor.

Open System Settings, go to Privacy and Security, select Accessibility, and toggle Wach on. If Wach does not appear in the list, start the app once (right-click, Open), then check again.

## Usage

Click the eye icon in the menu bar (top-right, near the clock). The menu offers:

- **Start** — begin the mouse-mover loop (eye opens)
- **Stop** — pause the mouse-mover loop (eye closes)
- **About** — version info and credits
- **Quit** — exit the app completely

The app starts with the mover active by default. Use Stop if you need to step away without keeping the system awake.

## Development

### Project Structure

```
wach/
  cmd/main.go                 Entry point (systray menu bar app)
  pkg/wach/
    wach.go                   Core logic and lifecycle
    state.go                  Thread-safe state management
    cgo_darwin.go             CoreGraphics C bindings
    alert_darwin.m            Native macOS alert dialogs
    logger.go                 Logging (stdout or file)
    doc.go                    Package documentation
    wach_test.go              28 tests (unit and CGo bridge)
  assets/icon/
    icon.go                   Embedded tray icons (PNG byte arrays)
    icon_test.go              Icon validation tests
  appInfo/
    Info.plist                macOS bundle metadata
    icon.icns                 App bundle icon
  lib/systray/                Patched systray library
  Makefile                    build, install, run, test targets
  go.mod / go.sum
  README.md
```

### Commands

```bash
make install    build and copy to /Applications/
make build      build app bundle in ./bin/
make run        run directly from terminal (no bundle)
make test       run all tests with race detector
make clean      remove build artifacts
```

### Technical Details

The original automatic-mouse-mover by prashantgupta24 used robotgo and activity-tracker, which are cross-platform libraries that pull in hundreds of dependencies (Windows and Linux APIs, OCR libraries, and more). This version replaces all of that with native macOS APIs:

| Feature | Original | Wach |
|---------|----------|------|
| Mouse movement | robotgo (CGo + C++) | CoreGraphics (C) |
| Idle detection | activity-tracker | CGEventSourceSecondsSinceLastEventType |
| Dialogs | gosx-notifier | NSAlert (Objective-C) |
| Binary size | ~6.4 MB | ~4.2 MB |
| Dependencies | ~40+ indirect | ~15 indirect |

The getlantern/systray library hardcodes icons at 16x16 points. This fork patches it to 20x20 for better visibility on modern Mac displays and uses template images for automatic light and dark mode support.

## Credits

Original concept by prashantgupta24 (github.com/prashantgupta24/automatic-mouse-mover). Built with Go, CoreGraphics, and native macOS APIs.

## License

MIT
