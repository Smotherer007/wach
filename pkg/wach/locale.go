package wach

// Locale holds all user-facing strings for one language.
type Locale struct {
	AppName   string
	AppTagline string

	// Menu
	MenuStart     string
	MenuStop      string
	MenuIdleTime  string
	MenuMoveDist  string
	MenuStats     string
	MenuResetStats string
	MenuLogin     string
	MenuBattery   string
	MenuSchedule  string
	MenuWorkdays  string
	MenuGitHub    string
	MenuAbout     string
	MenuQuit      string

	// Idle time options
	Idle30  string
	Idle60  string
	Idle120 string
	Idle300 string

	// Move distance options
	Move1  string
	Move5  string
	Move10 string
	Move20 string

	// Status
	StatusStopped    string
	StatusActive     string
	StatusNoMoveYet  string
	StatusIdle       string
	StatusTodayMoves string
	StatsToday       string

	// Tooltip descriptions
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

	// About dialog
	AboutTitle  string
	AboutMsg    string
	AboutSource string
	AboutAuthor string

	// Errors
	ErrNoPermissionTitle string
	ErrNoPermissionMsg   string
}

// locales maps language codes to their locale data.
var locales = map[string]Locale{
	"de": localeDE,
	"en": localeEN,
}

func GetLocale(lang string) Locale {
	if loc, ok := locales[lang]; ok {
		return loc
	}
	// Fallback to English
	return localeEN
}
