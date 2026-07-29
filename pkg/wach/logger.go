package wach

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	logDir      = "wach"
	logFileName = "wach.log"
)

// newLogger creates a configured logger for the wach service.
// If writeToFile is true, logs are written to ~/Library/Logs/wach/wach.log.
func newLogger(s *state, writeToFile bool) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	if writeToFile {
		home, err := os.UserHomeDir()
		if err == nil {
			logPath := filepath.Join(home, "Library", "Logs", logDir)
			if err := os.MkdirAll(logPath, 0755); err == nil {
				f, err := os.OpenFile(
					filepath.Join(logPath, logFileName),
					os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
				)
				if err == nil {
					logger.SetOutput(f)
					s.logFile = f
				}
			}
		}
	}

	return logger
}

// init configures the default logrus settings.
func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
