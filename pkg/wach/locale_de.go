package wach

var localeDE = Locale{
	AppName:    "Wach",
	AppTagline: "Mac bleibt wach",

	MenuStart:    "Start",
	MenuStop:     "Stop",
	MenuIdleTime: "Idle-Zeit",
	MenuMoveDist: "Bewegung",
	MenuStats:    "Aktivitats-Log",
	MenuResetStats: "Log zurucksetzen",
	MenuLogin:    "Beim Anmelden starten",
	MenuBattery:  "Batterie sparen",
	MenuSchedule: "Zeitplan aktiv",
	MenuWorkdays: "Nur Werktage (Mo-Fr)",
	MenuGitHub:   "GitHub",
	MenuAbout:    "Uber Wach",
	MenuQuit:     "Beenden",

	Idle30:  "30 Sekunden",
	Idle60:  "60 Sekunden",
	Idle120: "2 Minuten",
	Idle300: "5 Minuten",

	Move1:  "1 Pixel",
	Move5:  "5 Pixel",
	Move10: "10 Pixel",
	Move20: "20 Pixel",

	StatusStopped:    "Wach - gestoppt",
	StatusActive:     "Wach - aktiv",
	StatusNoMoveYet:  "noch keine Bewegung",
	StatusIdle:       "s idle",
	StatusTodayMoves: "heute %dx bewegt",
	StatsToday:       "Heute: %dx bewegt, %s aktiv",

	TipStart:    "Mausbewegung starten",
	TipStop:     "Mausbewegung stoppen",
	TipIdle:     "Leerlauf bis Mausbewegung",
	TipMove:     "Pixel pro Mausbewegung",
	TipLogin:    "LaunchAgent - automatisch starten",
	TipBattery:  "Pausieren bei weniger als 20 Prozent Akku",
	TipSchedule: "Nur zu eingestellten Zeiten aktiv",
	TipWorkdays: "Zeitplan nur Mo-Fr",
	TipGitHub:   "Projektseite im Browser offnen",
	TipAbout:    "Info zur App",
	TipQuit:     "Wach beenden",

	AboutTitle:  "Uber Wach",
	AboutMsg:    "Version %s\n\nEin Mausbeweger fur macOS (Apple Silicon).\nHalt den Mac wach, indem er bei Inaktivitat die Maus bewegt.\n\n(c) %s\n%s",
	AboutAuthor: "Patrick Weppelmann",
	AboutSource: "github.com/Smotherer007/wach",

	ErrNoPermissionTitle: "Zuganglichkeit erforderlich",
	ErrNoPermissionMsg:   "Maus konnte nicht bewegt werden (Versuch %d/%d)\nBitte Zuganglichkeit erlauben:\nSystemeinstellungen > Datenschutz > Bedienungshilfen > Wach",
}
