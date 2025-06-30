package di

import (
	"order/internal/infrastructure/logger"

	"go.uber.org/fx"
)

var LoggerModule = fx.Provide(
	// Config
	logger.NewConfig,

	// Logstash
	logger.NewLogstash,

	// Logrus
	logger.NewLogrus,

	// Logger
	logger.NewLogger,
)
