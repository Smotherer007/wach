package wach

import (
	"os"
	"sync"
	"time"
)

// state holds all mutable state for the Wach instance.
type state struct {
	mu        sync.RWMutex
	running   bool
	lastMoved time.Time
	lastErr   time.Time
	logFile   *os.File
	stats     ActivityStats
}

func (s *state) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	s.lastMoved = time.Time{}
	s.lastErr = time.Time{}
}

func (s *state) isRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *state) setRunning(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = v
	if v {
		s.stats.SessionStart = time.Now()
	}
}

func (s *state) getLastMoved() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastMoved
}

func (s *state) setLastMoved(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastMoved = t
	s.stats.TotalMoves++
	s.stats.TotalDuration += 10 * time.Second // each check interval

	// Daily tracking
	today := t.Format("2006-01-02")
	if s.stats.DailyDate == today {
		s.stats.DailyMoves++
	} else {
		s.stats.DailyDate = today
		s.stats.DailyMoves = 1
	}
}

func (s *state) getLastErrorTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastErr
}

func (s *state) setLastError(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastErr = t
}

func (s *state) getStats() ActivityStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

func (s *state) resetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats = ActivityStats{}
}
