package di

import (
	"courier/internal/infrastructure/telemetry"
	"go.uber.org/fx"
)

var TelemetryModule = fx.Provide(
	// Telemetry config
	telemetry.NewConfig,

	// Telemetry provider
	telemetry.NewProvider,
)
