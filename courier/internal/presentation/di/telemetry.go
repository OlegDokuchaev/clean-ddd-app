package di

import (
	"context"
	"courier/internal/infrastructure/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var TelemetryModule = fx.Options(
	fx.Invoke(setupTelemetryLifecycle),
)

func setupTelemetryLifecycle(lc fx.Lifecycle, tr *sdktrace.TracerProvider, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Println("Telemetry stopping...")

			if err := tr.Shutdown(ctx); err != nil {
				logger.Printf("error starting telemetry: %v", err)
				return err
			}

			logger.Println("Telemetry stopped")
			return nil
		},
	})
}
