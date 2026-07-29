package main

import (
	"github.com/getlantern/systray"
	"github.com/patrickweppelmann/wach/assets/icon"
	"github.com/patrickweppelmann/wach/pkg/wach"
	log "github.com/sirupsen/logrus"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		systray.SetTitle("")
		systray.SetTooltip("Wach - Mac bleibt wach")

		mAbout := systray.AddMenuItem("Uber Wach", "Info zur App")
		systray.AddSeparator()

		mStart := systray.AddMenuItem("Start", "Mausbewegung starten")
		mStop := systray.AddMenuItem("Stop", "Mausbewegung stoppen")
		mStop.Disable()

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Beenden", "Wach beenden")

		w := wach.GetInstance()
		w.Start()
		mStart.Disable()
		mStop.Enable()
		// Offenes Auge = aktiv
		systray.SetTemplateIcon(icon.IconOpen, icon.IconOpen)

		for {
			select {
			case <-mStart.ClickedCh:
				log.Info("starte wach")
				w.Start()
				mStart.Disable()
				mStop.Enable()
				systray.SetTemplateIcon(icon.IconOpen, icon.IconOpen)

			case <-mStop.ClickedCh:
				log.Info("stoppe wach")
				mStart.Enable()
				mStop.Disable()
				w.Stop()
				// Geschlossenes Auge = gestoppt
				systray.SetTemplateIcon(icon.IconClosed, icon.IconClosed)

			case <-mQuit.ClickedCh:
				log.Info("beende wach")
				w.Stop()
				systray.Quit()
				return

			case <-mAbout.ClickedCh:
				wach.ShowAbout()
			}
		}
	}()
}

func onExit() {
	log.Info("wach beendet")
}
