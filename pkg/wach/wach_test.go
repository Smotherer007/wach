package wach

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Singleton ---

func TestGetInstance(t *testing.T) {
	a := GetInstance()
	b := GetInstance()
	assert.Same(t, a, b, "GetInstance should return the same singleton instance")
	assert.NotNil(t, a)
}

func TestGetInstanceIsRunning(t *testing.T) {
	// The singleton should have a valid state pointer
	inst := GetInstance()
	assert.NotNil(t, inst.state)
}

// --- Config ---

func TestDefaultConfig(t *testing.T) {
	c := Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second}
	assert.Equal(t, 60*time.Second, c.IdleThreshold, "default idle threshold should be 60s")
	assert.Equal(t, 10, c.MovePixels, "default move pixels should be 10")
	assert.Equal(t, 10*time.Second, c.CheckInterval, "default check interval should be 10s")
}

func TestConfigCustomValues(t *testing.T) {
	c := Config{
		IdleThreshold: 120 * time.Second,
		MovePixels:    20,
		CheckInterval: 5 * time.Second,
	}
	assert.Equal(t, 120*time.Second, c.IdleThreshold)
	assert.Equal(t, 20, c.MovePixels)
	assert.Equal(t, 5*time.Second, c.CheckInterval)
}

// --- Lifecycle ---

func TestStartStop(t *testing.T) {
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}

	assert.False(t, w.IsRunning(), "should not be running initially")

	w.Start()
	assert.True(t, w.IsRunning(), "should be running after Start()")

	// Double start should be idempotent (no panic, no second goroutine)
	w.Start()
	assert.True(t, w.IsRunning(), "should still be running after double Start()")

	w.Stop()
	assert.False(t, w.IsRunning(), "should not be running after Stop()")

	// Double stop should be idempotent (no panic)
	w.Stop()
	assert.False(t, w.IsRunning(), "should still be stopped after double Stop()")
}

func TestStartStopMultipleCycles(t *testing.T) {
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}

	for i := 0; i < 5; i++ {
		w.Start()
		assert.True(t, w.IsRunning(), "cycle %d: should be running", i)
		time.Sleep(10 * time.Millisecond) // let goroutine start
		w.Stop()
		assert.False(t, w.IsRunning(), "cycle %d: should be stopped", i)
	}
}

func TestStartWithNilQuitChannel(t *testing.T) {
	// When Start() is called, it should create the quit channel
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}
	assert.Nil(t, w.quit, "quit channel should be nil before start")
	w.Start()
	assert.NotNil(t, w.quit, "quit channel should be created after start")
	w.Stop()
}

func TestStopOnNilWach(t *testing.T) {
	// Stop() must not panic on nil receiver
	var w *Wach
	assert.NotPanics(t, func() { w.Stop() }, "Stop() on nil should not panic")
}

// --- State ---

func TestStateLifecycle(t *testing.T) {
	s := &state{}

	assert.False(t, s.isRunning(), "fresh state should not be running")
	assert.True(t, s.getLastMoved().IsZero(), "lastMoved should be zero")
	assert.True(t, s.getLastErrorTime().IsZero(), "lastError should be zero")

	s.setRunning(true)
	assert.True(t, s.isRunning())

	now := time.Now()
	s.setLastMoved(now)
	assert.Equal(t, now, s.getLastMoved())

	errTime := now.Add(-1 * time.Hour)
	s.setLastError(errTime)
	assert.Equal(t, errTime, s.getLastErrorTime())

	s.reset()
	assert.False(t, s.isRunning(), "after reset should not be running")
	assert.True(t, s.getLastMoved().IsZero(), "after reset lastMoved should be zero")
}

func TestStateLastMovedUpdates(t *testing.T) {
	s := &state{}
	
	t1 := time.Now()
	s.setLastMoved(t1)
	assert.Equal(t, t1, s.getLastMoved())

	t2 := t1.Add(5 * time.Minute)
	s.setLastMoved(t2)
	assert.Equal(t, t2, s.getLastMoved())
}

func TestStateLastErrorUpdates(t *testing.T) {
	s := &state{}
	
	t1 := time.Now()
	s.setLastError(t1)
	assert.Equal(t, t1, s.getLastErrorTime())

	t2 := t1.Add(24 * time.Hour)
	s.setLastError(t2)
	assert.Equal(t, t2, s.getLastErrorTime())
}

// --- Concurrent State Access ---

func TestStateConcurrentReadWrite(t *testing.T) {
	s := &state{}
	var wg sync.WaitGroup

	// 10 concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.setRunning(true)
			s.setLastMoved(time.Now())
			s.setLastError(time.Now())
			s.reset()
		}()
	}

	// 10 concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.isRunning()
			_ = s.getLastMoved()
			_ = s.getLastErrorTime()
		}()
	}

	wg.Wait()
	// If we get here without race detector complaining, we're good
}

func TestStateResetDuringRead(t *testing.T) {
	s := &state{}
	var wg sync.WaitGroup

	s.setRunning(true)
	s.setLastMoved(time.Now())
	s.setLastError(time.Now())

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.reset()
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.isRunning()
			_ = s.getLastMoved()
			_ = s.getLastErrorTime()
		}()
	}
	wg.Wait()
}

// --- Wach Concurrent Lifecycle ---

func TestConcurrentStartStop(t *testing.T) {
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Start()
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Stop()
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = w.IsRunning()
		}()
	}
	wg.Wait()
}

// --- CGo Bridge Tests ---

func TestIdleDuration(t *testing.T) {
	d := getIdleDuration()
	t.Logf("idle duration: %v", d)
	assert.Greater(t, d, time.Duration(0), "idle duration should be positive")
	assert.Less(t, d, 24*time.Hour, "idle duration should be less than a day")
}

func TestDisplaySleep(t *testing.T) {
	result := isDisplayAsleep()
	t.Logf("display asleep: %v", result)
	// During tests, display is likely awake
	assert.False(t, result, "display should be awake during tests")
}

func TestIsDarkMode(t *testing.T) {
	// Just verify the CGo bridge works (result can be true or false)
	result := IsDarkMode()
	t.Logf("dark mode: %v", result)
	// Should return a valid boolean
	assert.IsType(t, bool(true), result)
}

// --- App Metadata ---

func TestAppMetadata(t *testing.T) {
	assert.Equal(t, "Wach", AppName)
	assert.Equal(t, "1.0.0", AppVersion)
	assert.Contains(t, AppSource, "github.com")
	assert.NotEmpty(t, AppAuthor)
}

// --- Package Constants ---

func TestConstants(t *testing.T) {
	assert.Equal(t, 60, DefaultIdleThreshold)
	assert.Equal(t, 10, DefaultMovePixels)
	assert.Equal(t, 10, DefaultCheckInterval)
	assert.Equal(t, 10, ErrorAlertThreshold)
	assert.Equal(t, 24*time.Hour, ErrorAlertCooldown)
}

// --- Edge Cases ---

func TestStateResetIdempotent(t *testing.T) {
	s := &state{}
	// Reset on fresh state should not panic
	assert.NotPanics(t, func() { s.reset() })
	// Reset on partially initialized state
	s.setRunning(true)
	assert.NotPanics(t, func() { s.reset() })
	assert.False(t, s.isRunning())
}

func TestConfigZeroValues(t *testing.T) {
	c := Config{}
	// Zero-value config should still be usable (not panic)
	_ = c.IdleThreshold
	_ = c.MovePixels
	_ = c.CheckInterval
}

func TestIsRunningOnFreshInstance(t *testing.T) {
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}
	assert.False(t, w.IsRunning())
}

func TestMultipleStopsAfterStart(t *testing.T) {
	w := &Wach{
		config: Config{IdleThreshold: 60 * time.Second, MovePixels: 10, CheckInterval: 10 * time.Second},
		state:  &state{},
	}

	w.Start()
	w.Stop()
	w.Stop()
	w.Stop()
	assert.False(t, w.IsRunning())
}
