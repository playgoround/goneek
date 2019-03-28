package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func New() *log.Logger {
	lvlStr := os.Getenv("GONEEK_LOG")
	lvl, err := log.ParseLevel(lvlStr)
	if err != nil {
		lvl = log.InfoLevel
	}

	logger := log.New()
	logger.SetLevel(lvl)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}
