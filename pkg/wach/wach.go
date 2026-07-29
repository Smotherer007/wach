// Package wach provides a macOS-native mouse mover that keeps your Mac awake
// by simulating user activity when the system is idle.
//
// Wach (German for "awake") uses CoreGraphics APIs directly instead of
// third-party libraries, making it lightweight and optimized for Apple Silicon.
package wach

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var global *Wach

const (
	DefaultIdleThreshold = 60
	DefaultMovePixels    = 10
	DefaultCheckInterval = 10
	ErrorAlertThreshold  = 10
	ErrorAlertCooldown   = 24 * time.Hour
	IdleLogInterval      = 5
)

// Config holds the operational parameters for the mouse mover.
type Config struct {
	IdleThreshold time.Duration
	MovePixels    int
	CheckInterval time.Duration
}

// Wach is the main controller for the mouse-mover service.
type Wach struct {
	mu       sync.Mutex
	config   Config
	settings Settings
	locale   Locale
	quit     chan struct{}
	done     chan struct{}
	state    *state
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

	w.settings = LoadSettings()
	w.applySettings()

	w.quit = make(chan struct{})
	w.done = make(chan struct{})
	w.state.setRunning(true)

	go w.run()
}

func (w *Wach) applySettings() {
	w.config.IdleThreshold = time.Duration(w.settings.IdleSeconds) * time.Second
	w.config.MovePixels = w.settings.MovePixels
	w.config.CheckInterval = DefaultCheckInterval * time.Second
}

// Stop stops the mouse-mover loop and waits for completion.
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

// GetSettings returns the current settings.
func (w *Wach) GetSettings() Settings {
	return w.settings
}

// UpdateSettings persists new settings and applies them at next start.
func (w *Wach) UpdateSettings(s Settings) {
	w.settings = s
	SaveSettings(s)
	if w.state.isRunning() {
		w.applySettings()
	}
}

// GetLocale returns the current locale.
func (w *Wach) GetLocale() Locale {
	return w.locale
}

// GetStats returns current activity statistics.
func (w *Wach) GetStats() ActivityStats {
	return w.state.getStats()
}

// ResetStats resets activity counters.
func (w *Wach) ResetStats() {
	w.state.resetStats()
}

// StatusText returns a human-readable status line.
func (w *Wach) StatusText() string {
	l := w.locale
	if !w.state.isRunning() {
		return l.StatusStopped
	}
	stats := w.state.getStats()
	last := w.state.getLastMoved()
	if last.IsZero() {
		return l.StatusActive + " - " + l.StatusNoMoveYet
	}
	idle := getIdleDuration().Round(time.Second)
	return fmt.Sprintf("%s (%d%s, %s)",
		l.StatusActive, int(idle.Seconds()), l.StatusIdle,
		fmt.Sprintf(l.StatusTodayMoves, stats.DailyMoves))
}

// IsWithinSchedule checks whether the current time falls within the configured schedule.
func (w *Wach) IsWithinSchedule() bool {
	s := w.settings
	if !s.ScheduleEnabled {
		return true
	}

	now := time.Now()

	// Workdays check (Mon=1..Fri=5)
	if s.ScheduleWorkdays {
		weekday := now.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			return false
		}
	}

	// Time range check
	start, err1 := time.Parse("15:04", s.ScheduleStart)
	end, err2 := time.Parse("15:04", s.ScheduleEnd)
	if err1 != nil || err2 != nil {
		return true // invalid config = don't block
	}

	nowMin := now.Hour()*60 + now.Minute()
	startMin := start.Hour()*60 + start.Minute()
	endMin := end.Hour()*60 + end.Minute()

	return nowMin >= startMin && nowMin < endMin
}

// IsBatteryLow returns true if on battery and below 20%.
func (w *Wach) IsBatteryLow() bool {
	if !w.settings.BatterySave {
		return false
	}
	pct := getBatteryPercent()
	return pct >= 0 && pct < 20
}

// run is the main event loop.
func (w *Wach) run() {
	logger := newLogger(w.state, false)
	movePixel := w.config.MovePixels
	idleThreshold := w.config.IdleThreshold

	// Prevent the system from sleeping while wach is active
	assertionID := createWakeAssertion("Wach keeps your Mac awake")
	if assertionID != 0 {
		logger.Info("system-schlaf-assertion aktiv — rechner bleibt wach")
	} else {
		logger.Warn("konnte schlaf-assertion nicht erstellen")
	}
	defer releaseWakeAssertion(assertionID)

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
			// Schedule check
			if !w.IsWithinSchedule() {
				if activeCount == 0 {
					logger.Debug("ausserhalb des Zeitplans — pausiert")
				}
				activeCount++
				continue
			}

			// Battery check
			if w.IsBatteryLow() {
				if activeCount == 0 {
					logger.Debug("akku niedrig — pausiert (Batteriesparmodus)")
				}
				activeCount++
				continue
			}

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
				logger.Debugf("display schläft — überspringe")
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
				l := w.locale
				msg := fmt.Sprintf(l.ErrNoPermissionMsg, errCount, ErrorAlertThreshold)
				logger.Warn(msg)

				if errCount >= ErrorAlertThreshold {
					lastErr := w.state.getLastErrorTime()
					if time.Since(lastErr) > ErrorAlertCooldown {
						showAlert(l.ErrNoPermissionTitle, msg)
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
	lang := SystemLanguage()
	global = &Wach{
		config: Config{
			IdleThreshold: DefaultIdleThreshold * time.Second,
			MovePixels:    DefaultMovePixels,
			CheckInterval: DefaultCheckInterval * time.Second,
		},
		settings: DefaultSettings(),
		locale:   GetLocale(lang),
		state:    &state{},
	}
}
