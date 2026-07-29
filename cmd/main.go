package main

import (
	"fmt"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/patrickweppelmann/wach/assets/icon"
	"github.com/patrickweppelmann/wach/pkg/wach"
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
	l := w.GetLocale()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		status := w.StatusText()
		mStatus.SetTitle(status)
		systray.SetTooltip(status)

		stats := w.GetStats()
		if stats.DailyMoves > 0 {
			total := time.Duration(stats.TotalDuration).Round(time.Second)
			mStats.SetTitle(fmt.Sprintf(l.StatsToday, stats.DailyMoves, total))
		} else {
			mStats.SetTitle(l.MenuStats)
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
		l := w.GetLocale()

		systray.SetTitle("")
		systray.SetTooltip(w.StatusText())
		systray.SetTemplateIcon(icon.IconOpen, icon.IconOpen)

		mStatus = systray.AddMenuItem("Wach", "")
		mStatus.Disable()
		systray.AddSeparator()

		mStart = systray.AddMenuItem(l.MenuStart, l.TipStart)
		mStop = systray.AddMenuItem(l.MenuStop, l.TipStop)

		systray.AddSeparator()

		mIdle := systray.AddMenuItem(l.MenuIdleTime, l.TipIdle)
		mIdle30 = mIdle.AddSubMenuItem(l.Idle30, "")
		mIdle60 = mIdle.AddSubMenuItem(l.Idle60, "")
		mIdle120 = mIdle.AddSubMenuItem(l.Idle120, "")
		mIdle300 = mIdle.AddSubMenuItem(l.Idle300, "")

		mMove := systray.AddMenuItem(l.MenuMoveDist, l.TipMove)
		mMove1 = mMove.AddSubMenuItem(l.Move1, "")
		mMove5 = mMove.AddSubMenuItem(l.Move5, "")
		mMove10 = mMove.AddSubMenuItem(l.Move10, "")
		mMove20 = mMove.AddSubMenuItem(l.Move20, "")

		systray.AddSeparator()

		mStats = systray.AddMenuItem(l.MenuStats, "")
		mStats.Disable()
		mResetStat = systray.AddMenuItem(l.MenuResetStats, "")

		systray.AddSeparator()

		mLogin = systray.AddMenuItem(l.MenuLogin, l.TipLogin)
		mBattery = systray.AddMenuItem(l.MenuBattery, l.TipBattery)
		mSchedule = systray.AddMenuItem(l.MenuSchedule, l.TipSchedule)
		mWorkdays = systray.AddMenuItem(l.MenuWorkdays, l.TipWorkdays)

		systray.AddSeparator()

		mGitHub = systray.AddMenuItem(l.MenuGitHub, l.TipGitHub)
		mAbout = systray.AddMenuItem(l.MenuAbout, l.TipAbout)
		mQuit = systray.AddMenuItem(l.MenuQuit, l.TipQuit)

		updateMenuFromSettings(s)
		applyLoginLaunch(s.StartAtLogin)

		w.UpdateSettings(s)
		w.Start()
		mStart.Disable()
		mStop.Enable()

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

			case <-mIdle30.ClickedCh:
				s := w.GetSettings(); s.IdleSeconds = 30
				w.UpdateSettings(s); markIdle(30)

			case <-mIdle60.ClickedCh:
				s := w.GetSettings(); s.IdleSeconds = 60
				w.UpdateSettings(s); markIdle(60)

			case <-mIdle120.ClickedCh:
				s := w.GetSettings(); s.IdleSeconds = 120
				w.UpdateSettings(s); markIdle(120)

			case <-mIdle300.ClickedCh:
				s := w.GetSettings(); s.IdleSeconds = 300
				w.UpdateSettings(s); markIdle(300)

			case <-mMove1.ClickedCh:
				s := w.GetSettings(); s.MovePixels = 1
				w.UpdateSettings(s); markMove(1)

			case <-mMove5.ClickedCh:
				s := w.GetSettings(); s.MovePixels = 5
				w.UpdateSettings(s); markMove(5)

			case <-mMove10.ClickedCh:
				s := w.GetSettings(); s.MovePixels = 10
				w.UpdateSettings(s); markMove(10)

			case <-mMove20.ClickedCh:
				s := w.GetSettings(); s.MovePixels = 20
				w.UpdateSettings(s); markMove(20)

			case <-mLogin.ClickedCh:
				s := w.GetSettings(); s.StartAtLogin = !s.StartAtLogin
				applyLoginLaunch(s.StartAtLogin)
				w.UpdateSettings(s)
				if s.StartAtLogin { mLogin.Check() } else { mLogin.Uncheck() }

			case <-mBattery.ClickedCh:
				s := w.GetSettings(); s.BatterySave = !s.BatterySave
				w.UpdateSettings(s)
				if s.BatterySave { mBattery.Check() } else { mBattery.Uncheck() }

			case <-mSchedule.ClickedCh:
				s := w.GetSettings(); s.ScheduleEnabled = !s.ScheduleEnabled
				w.UpdateSettings(s)
				if s.ScheduleEnabled { mSchedule.Check() } else { mSchedule.Uncheck() }

			case <-mWorkdays.ClickedCh:
				s := w.GetSettings(); s.ScheduleWorkdays = !s.ScheduleWorkdays
				w.UpdateSettings(s)
				if s.ScheduleWorkdays { mWorkdays.Check() } else { mWorkdays.Uncheck() }

			case <-mResetStat.ClickedCh:
				w.ResetStats()

			case <-mGitHub.ClickedCh:
				wach.OpenGitHub()

			case <-mAbout.ClickedCh:
				wach.ShowAbout()

			case <-mQuit.ClickedCh:
				w.Stop()
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {}
