package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`
	LogLevel    string `envconfig:"LOG_LEVEL" required:"true"`

	LogstashHost string `envconfig:"LOGSTASH_HOST" default:""`
	LogstashPort string `envconfig:"LOGSTASH_PORT" default:""`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load logger config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) GetLogLevel() logrus.Level {
	switch c.LogLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
