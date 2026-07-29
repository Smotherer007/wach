package wach

var localeEN = Locale{
	AppName:    "Wach",
	AppTagline: "Keep your Mac awake",

	MenuStart:    "Start",
	MenuStop:     "Stop",
	MenuIdleTime: "Idle Time",
	MenuMoveDist: "Move Distance",
	MenuStats:    "Activity Log",
	MenuResetStats: "Reset Log",
	MenuLogin:    "Launch at Login",
	MenuBattery:  "Battery Saver",
	MenuSchedule: "Schedule",
	MenuWorkdays: "Weekdays only (Mon-Fri)",
	MenuGitHub:   "GitHub",
	MenuAbout:    "About Wach",
	MenuQuit:     "Quit",

	Idle30:  "30 seconds",
	Idle60:  "60 seconds",
	Idle120: "2 minutes",
	Idle300: "5 minutes",

	Move1:  "1 pixel",
	Move5:  "5 pixels",
	Move10: "10 pixels",
	Move20: "20 pixels",

	StatusStopped:    "Wach - stopped",
	StatusActive:     "Wach - active",
	StatusNoMoveYet:  "no moves yet",
	StatusIdle:       "s idle",
	StatusTodayMoves: "%d moves today",
	StatsToday:       "Today: %d moves, %s active",

	TipStart:    "Start mouse mover",
	TipStop:     "Stop mouse mover",
	TipIdle:     "Idle time before mouse move",
	TipMove:     "Pixels per mouse move",
	TipLogin:    "LaunchAgent - start automatically at login",
	TipBattery:  "Pause when battery is below 20%",
	TipSchedule: "Only active during configured hours",
	TipWorkdays: "Schedule applies Mon-Fri only",
	TipGitHub:   "Open GitHub project page",
	TipAbout:    "About this app",
	TipQuit:     "Quit Wach",

	AboutTitle:  "About Wach",
	AboutMsg:    "Version %s\n\nA mouse mover for macOS (Apple Silicon).\nKeeps your Mac awake by moving the mouse when idle.\n\n(c) %s\n%s",
	AboutAuthor: "Patrick Weppelmann",
	AboutSource: "github.com/Smotherer007/wach",

	ErrNoPermissionTitle: "Accessibility Permission Required",
	ErrNoPermissionMsg:   "Could not move mouse (attempt %d/%d)\nPlease grant accessibility permission:\nSystem Settings > Privacy > Accessibility > Wach",
}
