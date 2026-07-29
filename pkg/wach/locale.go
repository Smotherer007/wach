package wach

// Locale holds all user-facing strings for one language.
type Locale struct {
	AppName   string
	AppTagline string

	MenuStart      string
	MenuStop       string
	MenuIdleTime   string
	MenuMoveDist   string
	MenuStats      string
	MenuResetStats string
	MenuLogin      string
	MenuBattery    string
	MenuSchedule   string
	MenuWorkdays   string
	MenuGitHub     string
	MenuAbout      string
	MenuQuit       string

	Idle30  string
	Idle60  string
	Idle120 string
	Idle300 string

	Move1  string
	Move5  string
	Move10 string
	Move20 string

	StatusStopped    string
	StatusActive     string
	StatusNoMoveYet  string
	StatusIdle       string
	StatusTodayMoves string
	StatsToday       string

	TipStart    string
	TipStop     string
	TipIdle     string
	TipMove     string
	TipLogin    string
	TipBattery  string
	TipSchedule string
	TipWorkdays string
	TipGitHub   string
	TipAbout    string
	TipQuit     string

	AboutTitle  string
	AboutMsg    string
	AboutSource string
	AboutAuthor string

	ErrNoPermissionTitle string
	ErrNoPermissionMsg   string
}

func GetLocale(lang string) Locale {
	if loc, ok := locales[lang]; ok {
		return loc
	}
	// Try parent language (e.g. "zh-Hans-CN" -> "zh-hans")
	if len(lang) >= 2 {
		parent := lang[:2]
		if loc, ok := locales[parent]; ok {
			return loc
		}
	}
	return locales["en"]
}
