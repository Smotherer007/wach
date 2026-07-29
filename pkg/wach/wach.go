// Package wach provides a macOS-native mouse mover that keeps your Mac awake
// by simulating user activity when the system is idle.
//
// Wach (German for "awake") uses CoreGraphics APIs directly instead of
// third-party libraries, making it lightweight and optimized for Apple Silicon.
//
// Requirements:
//   - macOS 11.0+ (Big Sur or later)
//   - Accessibility permission:
//     System Settings > Privacy & Security > Accessibility > add Wach.app
//
// How it works:
//   Every 10 seconds, Wach checks the system idle time via CoreGraphics.
//   If the system has been idle for more than the configured threshold (default 60s),
//   and the display is awake, it moves the mouse cursor by a few pixels to simulate
//   user activity. The move direction alternates to keep the cursor from drifting.
package wach

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Package-level singleton
var global *Wach

// Default configuration constants
const (
	DefaultIdleThreshold = 60            // seconds of idleness before mouse move
	DefaultMovePixels    = 10            // pixels to move the cursor
	DefaultCheckInterval = 10            // seconds between idle checks
	ErrorAlertThreshold  = 10            // consecutive failures before showing alert
	ErrorAlertCooldown   = 24 * time.Hour // minimum time between error alerts
	IdleLogInterval      = 5             // log every Nth idle check when active
)

// Config holds the configurable parameters for the mouse mover.
type Config struct {
	IdleThreshold time.Duration
	MovePixels    int
	CheckInterval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		IdleThreshold: DefaultIdleThreshold * time.Second,
		MovePixels:    DefaultMovePixels,
		CheckInterval: DefaultCheckInterval * time.Second,
	}
}

// Wach is the main controller for the mouse-mover service.
type Wach struct {
	mu     sync.Mutex // protects start/stop lifecycle
	config Config
	quit   chan struct{}
	done   chan struct{}
	state  *state
}

// GetInstance returns the singleton Wach instance.
func GetInstance() *Wach {
	runtime.LockOSThread()
	return global
}

// Start begins the mouse-mover loop. Thread-safe. Idempotent.
func (w *Wach) Start() {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state.isRunning() {
		return
	}

	w.quit = make(chan struct{})
	w.done = make(chan struct{})
	w.state.setRunning(true)

	go w.run()
}

// Stop stops the mouse-mover loop and waits for completion. Thread-safe. Idempotent.
func (w *Wach) Stop() {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.state.isRunning() {
		return
	}

	w.quit <- struct{}{}
	<-w.done

	if w.state.logFile != nil {
		w.state.logFile.Close()
		w.state.logFile = nil
	}
}

// IsRunning returns whether the mover loop is active.
func (w *Wach) IsRunning() bool {
	if w == nil {
		return false
	}
	return w.state.isRunning()
}

// run is the main event loop (runs in a separate goroutine).
func (w *Wach) run() {
	logger := newLogger(w.state, false)
	movePixel := w.config.MovePixels
	idleThreshold := w.config.IdleThreshold

	w.state.setLastMoved(time.Now())
	logger.Infof("wach gestartet — check alle %v, idle nach %v, bewege um %dpx",
		w.config.CheckInterval, idleThreshold, movePixel)

	ticker := time.NewTicker(w.config.CheckInterval)
	defer ticker.Stop()

	errCount := 0
	activeCount := 0

	for {
		select {
		case <-ticker.C:
			idle := getIdleDuration()

			if idle < idleThreshold {
				activeCount++
				if activeCount%IdleLogInterval == 0 {
					logger.Debugf("benutzer aktiv (idle: %v) — kein Eingriff", idle.Round(time.Second))
				}
				errCount = 0
				continue
			}

			activeCount = 0

			if isDisplayAsleep() {
				logger.Debugf("display schläft — überspringe Mausbewegung")
				continue
			}

			moved := tryMoveMouse(movePixel)

			if moved {
				w.state.setLastMoved(time.Now())
				logger.Infof("maus bewegt um %+dpx (idle: %v)", movePixel, idle.Round(time.Second))
				movePixel = -movePixel
				errCount = 0
			} else {
				errCount++
				msg := fmt.Sprintf(
					"Maus konnte nicht bewegt werden (Versuch %d/%d)\n"+
						"Bitte Zuganglichkeit erlauben:\n"+
						"Systemeinstellungen > Datenschutz > Bedienungshilfen > Wach",
					errCount, ErrorAlertThreshold,
				)
				logger.Warn(msg)

				if errCount >= ErrorAlertThreshold {
					lastErr := w.state.getLastErrorTime()
					if time.Since(lastErr) > ErrorAlertCooldown {
						showAlert("Zugänglichkeit erforderlich", msg)
						w.state.setLastError(time.Now())
					}
					errCount = 0
				}
			}

		case <-w.quit:
			logger.Info("wach gestoppt")
			w.state.setRunning(false)
			close(w.done)
			return
		}
	}
}

func init() {
	global = &Wach{
		config: DefaultConfig(),
		state:  &state{},
	}
}
