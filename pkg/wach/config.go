package wach

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Settings represents user-configurable options persisted to disk.
type Settings struct {
	IdleSeconds  int  `json:"idle_seconds"`
	MovePixels   int  `json:"move_pixels"`
	StartAtLogin bool `json:"start_at_login"`
	BatterySave  bool `json:"battery_save"`

	ScheduleEnabled    bool   `json:"schedule_enabled"`
	ScheduleStart      string `json:"schedule_start"`      // "09:00"
	ScheduleEnd        string `json:"schedule_end"`        // "18:00"
	ScheduleWorkdays   bool   `json:"schedule_workdays"`   // only Mon-Fri
}

// ActivityStats tracks usage statistics.
type ActivityStats struct {
	TotalMoves    int64         `json:"total_moves"`
	TotalDuration time.Duration `json:"total_duration"` // cumulative active time
	SessionStart  time.Time     `json:"session_start"`
	DailyMoves    int           `json:"daily_moves"`
	DailyDate     string        `json:"daily_date"` // "YYYY-MM-DD"
}

func DefaultSettings() Settings {
	return Settings{
		IdleSeconds:  60,
		MovePixels:   10,
		StartAtLogin: false,
		BatterySave:  false,
	}
}

var settingsPath string

func init() {
	home, err := os.UserHomeDir()
	if err == nil {
		dir := filepath.Join(home, "Library", "Application Support", "wach")
		os.MkdirAll(dir, 0755)
		settingsPath = filepath.Join(dir, "settings.json")
	}
}

// LoadSettings reads settings from disk, returns defaults on error.
func LoadSettings() Settings {
	s := DefaultSettings()
	if settingsPath == "" {
		return s
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return s
	}
	json.Unmarshal(data, &s)
	// fill zero values with defaults
	if s.IdleSeconds <= 0 {
		s.IdleSeconds = DefaultIdleThreshold
	}
	if s.MovePixels <= 0 {
		s.MovePixels = DefaultMovePixels
	}
	return s
}

// SaveSettings writes settings to disk.
func SaveSettings(s Settings) error {
	if settingsPath == "" {
		return nil
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, data, 0644)
}
