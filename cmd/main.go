package main

import (
	"fmt"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/patrickweppelmann/wach/assets/icon"
	"github.com/patrickweppelmann/wach/pkg/wach"
	log "github.com/sirupsen/logrus"
)

var (
	mStart     *systray.MenuItem
	mStop      *systray.MenuItem
	mStatus    *systray.MenuItem
	mStats     *systray.MenuItem
	mIdle30    *systray.MenuItem
	mIdle60    *systray.MenuItem
	mIdle120   *systray.MenuItem
	mIdle300   *systray.MenuItem
	mMove1     *systray.MenuItem
	mMove5     *systray.MenuItem
	mMove10    *systray.MenuItem
	mMove20    *systray.MenuItem
	mLogin     *systray.MenuItem
	mBattery   *systray.MenuItem
	mSchedule  *systray.MenuItem
	mWorkdays  *systray.MenuItem
	mResetStat *systray.MenuItem
	mGitHub    *systray.MenuItem
	mAbout     *systray.MenuItem
	mQuit      *systray.MenuItem
)

func markIdle(secs int) {
	for _, item := range []*systray.MenuItem{mIdle30, mIdle60, mIdle120, mIdle300} {
		item.Uncheck()
	}
	switch secs {
	case 30:
		mIdle30.Check()
	case 60:
		mIdle60.Check()
	case 120:
		mIdle120.Check()
	case 300:
		mIdle300.Check()
	}
}

func markMove(px int) {
	for _, item := range []*systray.MenuItem{mMove1, mMove5, mMove10, mMove20} {
		item.Uncheck()
	}
	switch px {
	case 1:
		mMove1.Check()
	case 5:
		mMove5.Check()
	case 10:
		mMove10.Check()
	case 20:
		mMove20.Check()
	}
}

func updateMenuFromSettings(s wach.Settings) {
	markIdle(s.IdleSeconds)
	markMove(s.MovePixels)
	if s.StartAtLogin {
		mLogin.Check()
	} else {
		mLogin.Uncheck()
	}
	if s.BatterySave {
		mBattery.Check()
	} else {
		mBattery.Uncheck()
	}
	if s.ScheduleEnabled {
		mSchedule.Check()
	} else {
		mSchedule.Uncheck()
	}
	if s.ScheduleWorkdays {
		mWorkdays.Check()
	} else {
		mWorkdays.Uncheck()
	}
}

func statusLoop(w *wach.Wach) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		status := w.StatusText()
		mStatus.SetTitle(status)
		systray.SetTooltip(status)

		// Update stats display
		stats := w.GetStats()
		if stats.DailyMoves > 0 {
			total := time.Duration(stats.TotalDuration).Round(time.Second)
			mStats.SetTitle(fmt.Sprintf("Heute: %dx bewegt, %s aktiv", stats.DailyMoves, total))
		} else {
			mStats.SetTitle("Aktivitats-Log")
		}
	}
}

func applyLoginLaunch(enabled bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	launchDir := home + "/Library/LaunchAgents"
	plistPath := launchDir + "/com.patrickweppelmann.wach.plist"

	if enabled {
		os.MkdirAll(launchDir, 0755)
		content := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.patrickweppelmann.wach</string>
	<key>ProgramArguments</key>
	<array>
		<string>/Applications/Wach.app/Contents/MacOS/wach</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
</dict>
</plist>`
		os.WriteFile(plistPath, []byte(content), 0644)
	} else {
		os.Remove(plistPath)
	}
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		s := wach.LoadSettings()
		w := wach.GetInstance()

		systray.SetTitle("")
		systray.SetTooltip("Wach - Mac bleibt wach")
		systray.SetTemplateIcon(icon.IconOpen, icon.IconOpen)

		// Status (non-clickable)
		mStatus = systray.AddMenuItem("Wach", "Status")
		mStatus.Disable()
		systray.AddSeparator()

		// Start / Stop
		mStart = systray.AddMenuItem("Start", "Mausbewegung starten")
		mStop = systray.AddMenuItem("Stop", "Mausbewegung stoppen")

		systray.AddSeparator()

		// Idle time submenu
		mIdle := systray.AddMenuItem("Idle-Zeit", "Leerlauf bis Mausbewegung")
		mIdle30 = mIdle.AddSubMenuItem("30 Sekunden", "")
		mIdle60 = mIdle.AddSubMenuItem("60 Sekunden", "")
		mIdle120 = mIdle.AddSubMenuItem("2 Minuten", "")
		mIdle300 = mIdle.AddSubMenuItem("5 Minuten", "")

		// Move distance submenu
		mMove := systray.AddMenuItem("Bewegung", "Pixel pro Mausbewegung")
		mMove1 = mMove.AddSubMenuItem("1 Pixel", "kaum sichtbar")
		mMove5 = mMove.AddSubMenuItem("5 Pixel", "dezent")
		mMove10 = mMove.AddSubMenuItem("10 Pixel", "standard")
		mMove20 = mMove.AddSubMenuItem("20 Pixel", "deutlich")

		systray.AddSeparator()

		// Stats
		mStats = systray.AddMenuItem("Aktivitats-Log", "")
		mStats.Disable()
		mResetStat = systray.AddMenuItem("Log zurucksetzen", "")

		systray.AddSeparator()

		// Toggles
		mLogin = systray.AddMenuItem("Beim Anmelden starten", "LaunchAgent")
		mBattery = systray.AddMenuItem("Batterie sparen", "Pausieren bei <20% Akku")
		mSchedule = systray.AddMenuItem("Zeitplan aktiv", "Nur zu Arbeitszeiten")
		mWorkdays = systray.AddMenuItem("Nur Werktage (Mo-Fr)", "")

		systray.AddSeparator()

		// Info
		mGitHub = systray.AddMenuItem("GitHub", "Projektseite offnen")
		mAbout = systray.AddMenuItem("Uber Wach", "Info zur App")
		mQuit = systray.AddMenuItem("Beenden", "Wach beenden")

		// Apply settings to menu
		updateMenuFromSettings(s)

		// Start automatically
		applyLoginLaunch(s.StartAtLogin)

		w.UpdateSettings(s)
		w.Start()
		mStart.Disable()
		mStop.Enable()

		// Start status updater
		go statusLoop(w)

		for {
			select {
			case <-mStart.ClickedCh:
				w.Start()
				mStart.Disable()
				mStop.Enable()
				systray.SetTemplateIcon(icon.IconOpen, icon.IconOpen)

			case <-mStop.ClickedCh:
				mStart.Enable()
				mStop.Disable()
				w.Stop()
				systray.SetTemplateIcon(icon.IconClosed, icon.IconClosed)

			// Idle time
			case <-mIdle30.ClickedCh:
				s := w.GetSettings()
				s.IdleSeconds = 30
				w.UpdateSettings(s)
				markIdle(30)

			case <-mIdle60.ClickedCh:
				s := w.GetSettings()
				s.IdleSeconds = 60
				w.UpdateSettings(s)
				markIdle(60)

			case <-mIdle120.ClickedCh:
				s := w.GetSettings()
				s.IdleSeconds = 120
				w.UpdateSettings(s)
				markIdle(120)

			case <-mIdle300.ClickedCh:
				s := w.GetSettings()
				s.IdleSeconds = 300
				w.UpdateSettings(s)
				markIdle(300)

			// Move distance
			case <-mMove1.ClickedCh:
				s := w.GetSettings()
				s.MovePixels = 1
				w.UpdateSettings(s)
				markMove(1)

			case <-mMove5.ClickedCh:
				s := w.GetSettings()
				s.MovePixels = 5
				w.UpdateSettings(s)
				markMove(5)

			case <-mMove10.ClickedCh:
				s := w.GetSettings()
				s.MovePixels = 10
				w.UpdateSettings(s)
				markMove(10)

			case <-mMove20.ClickedCh:
				s := w.GetSettings()
				s.MovePixels = 20
				w.UpdateSettings(s)
				markMove(20)

			// Toggles
			case <-mLogin.ClickedCh:
				s := w.GetSettings()
				s.StartAtLogin = !s.StartAtLogin
				applyLoginLaunch(s.StartAtLogin)
				w.UpdateSettings(s)
				if s.StartAtLogin {
					mLogin.Check()
				} else {
					mLogin.Uncheck()
				}

			case <-mBattery.ClickedCh:
				s := w.GetSettings()
				s.BatterySave = !s.BatterySave
				w.UpdateSettings(s)
				if s.BatterySave {
					mBattery.Check()
				} else {
					mBattery.Uncheck()
				}

			case <-mSchedule.ClickedCh:
				s := w.GetSettings()
				s.ScheduleEnabled = !s.ScheduleEnabled
				w.UpdateSettings(s)
				if s.ScheduleEnabled {
					mSchedule.Check()
				} else {
					mSchedule.Uncheck()
				}

			case <-mWorkdays.ClickedCh:
				s := w.GetSettings()
				s.ScheduleWorkdays = !s.ScheduleWorkdays
				w.UpdateSettings(s)
				if s.ScheduleWorkdays {
					mWorkdays.Check()
				} else {
					mWorkdays.Uncheck()
				}

			// Stats
			case <-mResetStat.ClickedCh:
				w.ResetStats()

			// Info
			case <-mGitHub.ClickedCh:
				wach.OpenGitHub()

			case <-mAbout.ClickedCh:
				wach.ShowAbout()

			case <-mQuit.ClickedCh:
				log.Info("beende wach")
				w.Stop()
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Info("wach beendet")
}
