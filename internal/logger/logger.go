package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func Setup() {
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	logger.SetOutput(os.Stdout)
}

func Close() {
	logger.Writer().Close()
}

func LogWithField(msg string, field ...interface{}) {
	Setup()
	logger.WithFields(log.Fields{
		"data": field,
	}).Info(msg)
	Close()
}

func Info(msg string, field interface{}) {
	Setup()
	if field != nil {
		logger.WithFields(log.Fields{
			"data": field,
		}).Info(msg)
	} else {
		logger.Info(msg)
	}
	Close()
}

func New() *log.Logger {
	Setup()
	return logger
}
