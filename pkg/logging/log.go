package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	name string
	*log.Entry
}

func GetLogger(name string) *Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)
	return &Logger{name, logger.WithField("logging", name)}
}
