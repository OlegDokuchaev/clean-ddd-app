package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogrus(config *Config, logstashHook logrus.Hook) *logrus.Logger {
	log := logrus.New()

	log.Hooks.Add(logstashHook)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	log.SetLevel(config.GetLogLevel())

	return log
}
